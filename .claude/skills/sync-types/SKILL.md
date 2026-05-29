---
name: sync-types
description: Diff Go domain structs in services/*/internal/ against TypeScript interfaces in web/src/types.ts to find drift in field names, JSON tags, types, and optionality. Use this whenever the user asks to check or sync the wire types, mentions Go/TS type drift, edits a Go struct that's exposed over HTTP, edits web/src/types.ts, or before opening a chapter PR. Also use proactively after editing any file under services/*/internal/<domain>/ that contains json: tags - the React types are hand-maintained and silently drift.
---

# sync-types

The Go services and the React app share a wire format but no code
generation. `web/src/types.ts` is hand-written to mirror the Go structs in
`services/*/internal/<domain>/`. When the Go side changes a field name, a
JSON tag, or a type and TS isn't updated, the bug is silent: requests still
serialize, fields just become `undefined` in the UI. This skill catches
that drift.

## When to run

- The user explicitly asks to "sync types", "check the wire types",
  "diff Go and TS", or similar.
- After any edit under `services/*/internal/<domain>/` that touches a
  `struct` with `json:` tags.
- Before running `chapter-pr` - drift in the wire types is a common
  end-of-chapter regression.
- When the user adds a new HTTP endpoint that returns a new shape.

If the user is just editing internal Go (no `json:` tags involved), skip
this - there's nothing to sync.

## How to do the diff

There are two viable approaches. Pick based on scale:

### Small diff (one or two structs touched in this session)

Read the Go struct(s) and `web/src/types.ts` directly and compare by eye.
This is faster than scripting for tiny changes. Look for:

1. **Field name match.** The TS interface field name must equal the Go
   `json:` tag, not the Go field name. Go uses PascalCase; the wire and
   TS are snake_case. `Name string \`json:"name"\`` → TS `name`.
2. **Type match.** Map Go → TS as:
   - `string` → `string`
   - `int*`, `float*`, `commodity.LabourMinutes`, `commodity.Quantity`
     → `number`
   - `bool` → `boolean`
   - `time.Time` → `string` (RFC3339, JSON serializes that way)
   - named ID types (`commodity.ID`) → `string`
   - structs → nested interface
   - slices → `T[]`
   - maps `map[string]T` → `Record<string, T>` or a typed shape
3. **Optionality match.** Go `*T` with `json:",omitempty"` → TS `field?: T`.
   Go `T` with no omitempty → TS required `field: T`.
4. **Tag presence.** A Go struct field without a `json:` tag still
   serializes - using its Go field name (PascalCase). That's almost
   always a bug; flag it.

### Larger diff (whole-repo audit, or several services involved)

Run `scripts/sync_check.py` from this skill directory:

```bash
python3 .claude/skills/sync-types/scripts/sync_check.py
```

It walks `services/*/internal/**/*.go`, extracts every struct with at
least one `json:` tag, parses `web/src/types.ts`, and prints a report of
mismatches grouped by struct. Exit code is non-zero if drift is found,
which makes it usable in `make` later.

The script is intentionally simple regex-based parsing, not a full Go AST
- it'll have false positives on unusual struct embedding. If a flagged
mismatch looks wrong, hand-verify it.

## Reporting

When you find drift, print a compact report. For each mismatched type,
show:

```
commodity.Commodity (services/commodity-service/internal/commodity/commodity.go:50)
  ts: web/src/types.ts:14
  - missing in TS: snlt_per_unit (number)
  - extra in TS:   legacy_field (string)
  - type mismatch: created_at  Go=time.Time→string  TS=number
```

Then ask the user whether to update `types.ts`, the Go struct, or leave
it (sometimes drift is intentional - e.g. a server-only field).

## What this skill does NOT do

- **Auto-fix the drift.** Direction of the fix (Go ↔ TS) is a judgment
  call. Surface it; let the user decide.
- **Sync request/response shapes that aren't in `types.ts`.** Some HTTP
  shapes live only inside `internal/transport/httpapi/handler.go` (e.g.
  `valueRequest`). If they're not in `types.ts`, they're not part of the
  wire contract this skill polices - mention them in passing if they
  changed, but don't fail on them.
- **Check the Go side against itself.** `bson:` and `json:` tags
  occasionally diverge (e.g. `_id` vs `id`); that's intentional and out
  of scope.

## Anti-patterns

- Don't suggest generating types from Go (e.g. tygo, openapi). The repo
  intentionally hand-maintains the TS side - it's a small surface and
  codegen would be more friction than the drift it prevents (CLAUDE.md
  anti-patterns).
- Don't rewrite `types.ts` to match a *transient* state of the Go code
  mid-edit. Wait until the Go side is settled.
