#!/usr/bin/env python3
"""Diff Go domain structs against web/src/types.ts.

Walks services/*/internal/**/*.go for structs with at least one `json:` tag,
parses web/src/types.ts for interface declarations, and reports drift in
field names, types, and optionality.

Regex-based, intentionally simple. Re-verify any flagged mismatch by hand
before changing code - false positives on unusual struct embedding are
expected.

Exit codes:
  0  no drift
  1  drift found
  2  could not locate repo root or files
"""

from __future__ import annotations

import re
import sys
from dataclasses import dataclass, field
from pathlib import Path


GO_TS_TYPE_MAP = {
    "string": "string",
    "bool": "boolean",
    "int": "number",
    "int8": "number",
    "int16": "number",
    "int32": "number",
    "int64": "number",
    "uint": "number",
    "uint8": "number",
    "uint16": "number",
    "uint32": "number",
    "uint64": "number",
    "float32": "number",
    "float64": "number",
    "time.Time": "string",
}


@dataclass
class GoField:
    name: str
    go_type: str
    optional: bool
    file: Path
    line: int


@dataclass
class GoStruct:
    name: str
    file: Path
    line: int
    fields: list[GoField] = field(default_factory=list)


@dataclass
class TSField:
    name: str
    ts_type: str
    optional: bool


@dataclass
class TSInterface:
    name: str
    line: int
    fields: list[TSField] = field(default_factory=list)


def find_repo_root(start: Path) -> Path | None:
    cur = start.resolve()
    for _ in range(8):
        if (cur / "go.mod").exists() and (cur / "web").exists():
            return cur
        if cur.parent == cur:
            return None
        cur = cur.parent
    return None


# Match the fields inside a `type Foo struct { ... }` block. Crude but works
# for the conventional formatting in this repo.
STRUCT_RE = re.compile(r"type\s+(\w+)\s+struct\s*\{([^}]*)\}", re.MULTILINE | re.DOTALL)
JSON_TAG_RE = re.compile(r'json:"([^"]*)"')


def parse_go_file(path: Path) -> list[GoStruct]:
    text = path.read_text(encoding="utf-8", errors="replace")
    out: list[GoStruct] = []
    for m in STRUCT_RE.finditer(text):
        name = m.group(1)
        body = m.group(2)
        struct_line = text[: m.start()].count("\n") + 1
        gs = GoStruct(name=name, file=path, line=struct_line)
        for raw_line in body.splitlines():
            line = raw_line.strip()
            if not line or line.startswith("//"):
                continue
            tag_match = JSON_TAG_RE.search(line)
            if not tag_match:
                continue
            tag_parts = tag_match.group(1).split(",")
            json_name = tag_parts[0]
            if json_name == "-" or json_name == "":
                continue
            optional = "omitempty" in tag_parts[1:]
            # Strip the tag and trailing backtick to inspect the type column.
            head = line.split("`", 1)[0].strip()
            tokens = head.split()
            if len(tokens) < 2:
                continue
            field_name = tokens[0]
            go_type = " ".join(tokens[1:])
            if go_type.startswith("*"):
                optional = True
                go_type = go_type[1:]
            gs.fields.append(
                GoField(
                    name=json_name,
                    go_type=go_type,
                    optional=optional,
                    file=path,
                    line=struct_line + body[: body.find(raw_line)].count("\n"),
                )
            )
            del field_name  # currently unused
        if gs.fields:
            out.append(gs)
    return out


# `export interface Foo { ... }` — body up to the matching closing brace.
TS_INTERFACE_RE = re.compile(
    r"export\s+interface\s+(\w+)\s*\{([^}]*)\}", re.MULTILINE | re.DOTALL
)


def parse_ts_file(path: Path) -> list[TSInterface]:
    text = path.read_text(encoding="utf-8", errors="replace")
    out: list[TSInterface] = []
    for m in TS_INTERFACE_RE.finditer(text):
        name = m.group(1)
        body = m.group(2)
        line_no = text[: m.start()].count("\n") + 1
        ti = TSInterface(name=name, line=line_no)
        for raw_line in body.splitlines():
            line = raw_line.strip()
            if not line or line.startswith("//"):
                continue
            # Strip inline // comment before further parsing.
            comment_idx = line.find("//")
            if comment_idx != -1:
                line = line[:comment_idx].strip()
            line = line.rstrip(";").rstrip(",")
            if not line or ":" not in line:
                continue
            head, ts_type = line.split(":", 1)
            head = head.strip()
            ts_type = ts_type.strip()
            optional = head.endswith("?")
            field_name = head.rstrip("?").strip()
            ti.fields.append(
                TSField(name=field_name, ts_type=ts_type, optional=optional)
            )
        if ti.fields:
            out.append(ti)
    return out


def normalize_go_type(go_type: str) -> str:
    """Reduce a Go type expression to a TS-comparable token. Returns 'unknown'
    when the mapping isn't safe to assert."""
    t = go_type.strip()
    if t.startswith("[]"):
        inner = normalize_go_type(t[2:])
        return f"{inner}[]" if inner != "unknown" else "unknown"
    if t.startswith("map["):
        # map[K]V → Record<K, V> — only handle string keys.
        end = t.find("]")
        if end == -1:
            return "unknown"
        key = t[4:end]
        val = normalize_go_type(t[end + 1 :])
        if key == "string" and val != "unknown":
            return f"Record<string,{val}>"
        return "unknown"
    if t in GO_TS_TYPE_MAP:
        return GO_TS_TYPE_MAP[t]
    # Named ID types and labour magnitudes that the repo aliases over a
    # primitive. Heuristic: a type ending in ID or a known alias maps to
    # string/number respectively.
    if t.endswith("ID") or t.endswith(".ID"):
        return "string"
    if t.endswith("LabourMinutes") or t.endswith("Quantity"):
        return "number"
    # Otherwise treat as a nested struct reference; we can't resolve it
    # without a real parser, so return the bare name and let the comparison
    # flag any TS-side surprise.
    return t.rsplit(".", 1)[-1]


def normalize_ts_type(ts_type: str) -> str:
    return ts_type.replace(" ", "")


def report(go_structs: list[GoStruct], ts_ifaces: list[TSInterface]) -> int:
    ts_by_name = {ti.name: ti for ti in ts_ifaces}
    drift_count = 0
    seen_ts: set[str] = set()

    # Compare each Go struct to the TS interface of the same name. Structs
    # without a TS counterpart are *not* errors — many internal Go structs
    # are server-only.
    for gs in go_structs:
        ti = ts_by_name.get(gs.name)
        if ti is None:
            continue
        seen_ts.add(gs.name)
        problems: list[str] = []

        go_by_name = {f.name: f for f in gs.fields}
        ts_by_field = {f.name: f for f in ti.fields}

        for name, gf in go_by_name.items():
            if name not in ts_by_field:
                problems.append(
                    f"  - missing in TS: {name} ({gf.go_type})"
                )
                continue
            tf = ts_by_field[name]
            want = normalize_go_type(gf.go_type)
            got = normalize_ts_type(tf.ts_type)
            if want != "unknown" and want != got:
                problems.append(
                    f"  - type mismatch: {name}  Go={gf.go_type}→{want}  TS={tf.ts_type}"
                )
            if gf.optional != tf.optional:
                problems.append(
                    f"  - optionality: {name}  Go optional={gf.optional}  TS optional={tf.optional}"
                )

        for name in ts_by_field:
            if name not in go_by_name:
                problems.append(
                    f"  - extra in TS: {name} ({ts_by_field[name].ts_type})"
                )

        if problems:
            drift_count += 1
            print(f"{gs.name} ({gs.file}:{gs.line})")
            print(f"  ts: web/src/types.ts:{ti.line}")
            for p in problems:
                print(p)
            print()

    # Surface TS interfaces that have no Go-side struct of the same name.
    # These might be intentional (e.g. UI-only request shapes) but the user
    # should at least see them once.
    for name, ti in ts_by_name.items():
        if name in seen_ts:
            continue
        if name in {
            "CreateCommodityInput",
            "UpdateCommodityInput",
            "ValueResponse",
            "SocialRelations",
        }:
            # Known wrapper shapes, named distinctly on purpose.
            continue
        print(f"{name}: TS interface has no Go struct of the same name")
        print(f"  ts: web/src/types.ts:{ti.line}")
        print()

    if drift_count == 0:
        print("no drift between Go structs and web/src/types.ts")
        return 0
    print(f"drift found in {drift_count} struct(s)")
    return 1


def main() -> int:
    here = Path(__file__).resolve()
    root = find_repo_root(here)
    if root is None:
        print("could not locate repo root (need go.mod and web/)", file=sys.stderr)
        return 2

    go_files: list[Path] = []
    for svc_dir in (root / "services").glob("*/internal"):
        go_files.extend(svc_dir.rglob("*.go"))
    go_files = [p for p in go_files if not p.name.endswith("_test.go")]

    ts_path = root / "web" / "src" / "types.ts"
    if not ts_path.exists():
        print(f"missing {ts_path}", file=sys.stderr)
        return 2

    go_structs: list[GoStruct] = []
    for gp in go_files:
        go_structs.extend(parse_go_file(gp))

    ts_ifaces = parse_ts_file(ts_path)

    return report(go_structs, ts_ifaces)


if __name__ == "__main__":
    sys.exit(main())
