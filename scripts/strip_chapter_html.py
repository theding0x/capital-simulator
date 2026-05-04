#!/usr/bin/env python3
"""Strip a capital-simulator chapter HTML to clean prose for LLM consumption.

Usage:
    python3 scripts/strip_chapter_html.py chapters/volume-1/01-the-commodity.html

Output is UTF-8 markdown-flavoured plain text on stdout. Sections become
## headings, paragraphs are preserved, value-form tables are rendered as
indented code blocks, and footnote references are dropped.

The output is intentionally compact — around 15-25 KB for a typical chapter —
so it fits comfortably in an LLM context without truncation.
"""
from __future__ import annotations

import re
import sys
from html.parser import HTMLParser
from pathlib import Path


# Tags whose text content we want to keep but not promote.
_INLINE_KEEP = {"em", "i", "b", "strong", "span", "sub"}
# Tags that produce block-level output.
_BLOCK_TAGS = {"p", "h3", "h4", "h5", "h6", "td", "th", "hr", "br", "tr"}


class ChapterParser(HTMLParser):
    def __init__(self) -> None:
        super().__init__()
        self._blocks: list[str] = []          # completed output blocks
        self._buf: list[str] = []             # text accumulator for current block
        self._tag_stack: list[str] = []
        self._in_footnotes = False            # True once we hit the h6 footnote block
        self._in_enote = False               # True inside sup.enote (drop content)
        self._table_row: list[str] = []       # cells for current table row
        self._table_rows: list[list[str]] = []
        self._in_table = False
        self._current_block_tag: str = ""
        self._current_class: str = ""

    # ------------------------------------------------------------------ #
    # HTML events                                                          #
    # ------------------------------------------------------------------ #

    def handle_starttag(self, tag: str, attrs: list[tuple[str, str | None]]) -> None:
        attr_dict = dict(attrs)
        cls = attr_dict.get("class") or ""

        self._tag_stack.append(tag)

        if tag == "sup" and "enote" in cls:
            self._in_enote = True
            return

        if tag == "h6":
            # h6 marks the start of the footnote block at the bottom — skip everything after.
            self._in_footnotes = True
            return

        if tag == "table":
            self._in_table = True
            self._table_rows = []
            return

        if tag == "tr":
            self._table_row = []
            return

        if tag in {"td", "th"}:
            self._buf = []
            return

        if tag in {"h3", "h4", "h5", "p"}:
            self._buf = []
            self._current_block_tag = tag
            self._current_class = cls
            return

        if tag == "hr":
            return  # decorative

        if tag == "br":
            if self._buf:
                self._buf.append(" ")
            return

    def handle_endtag(self, tag: str) -> None:
        if self._tag_stack and self._tag_stack[-1] == tag:
            self._tag_stack.pop()

        if tag == "sup":
            self._in_enote = False
            return

        if self._in_footnotes or self._in_enote:
            return

        if tag in {"td", "th"}:
            cell = _clean(" ".join(self._buf))
            if cell:
                self._table_row.append(cell)
            self._buf = []
            return

        if tag == "tr":
            if self._table_row:
                self._table_rows.append(self._table_row)
            self._table_row = []
            return

        if tag == "table":
            self._in_table = False
            rendered = _render_table(self._table_rows)
            if rendered:
                self._blocks.append(rendered)
            self._table_rows = []
            return

        if tag in {"h3", "h4", "h5", "p"}:
            text = _clean(" ".join(self._buf))
            if not text:
                self._buf = []
                return

            if tag == "h3":
                self._blocks.append(f"## {text}")
            elif tag == "h4":
                self._blocks.append(f"### {text}")
            elif tag == "h5":
                self._blocks.append(f"#### {text}")
            elif tag == "p":
                cls = self._current_class
                if "quoteb" in cls or "indentb" in cls:
                    self._blocks.append("> " + text)
                else:
                    self._blocks.append(text)

            self._buf = []
            self._current_block_tag = ""
            self._current_class = ""

    def handle_data(self, data: str) -> None:
        if self._in_footnotes or self._in_enote:
            return
        if self._in_table and self._tag_stack and self._tag_stack[-1] in {"td", "th"}:
            self._buf.append(data)
            return
        if self._current_block_tag or (self._tag_stack and self._tag_stack[-1] in _INLINE_KEEP):
            self._buf.append(data)

    def handle_entityref(self, name: str) -> None:
        _ENTITIES = {
            "nbsp": " ", "amp": "&", "lt": "<", "gt": ">",
            "quot": '"', "apos": "'", "mdash": "—", "ndash": "–",
            "rsquo": "'", "lsquo": "'", "rdquo": '"', "ldquo": '"',
            "hellip": "...", "middot": "·", "frac12": "½",
        }
        self.handle_data(_ENTITIES.get(name, ""))

    def handle_charref(self, name: str) -> None:
        try:
            code = int(name[1:], 16) if name.startswith("x") else int(name)
            self.handle_data(chr(code))
        except (ValueError, OverflowError):
            pass

    # ------------------------------------------------------------------ #
    # Result                                                               #
    # ------------------------------------------------------------------ #

    def text(self) -> str:
        out: list[str] = []
        prev_blank = False
        for block in self._blocks:
            stripped = block.strip()
            if not stripped:
                if not prev_blank:
                    out.append("")
                prev_blank = True
            else:
                out.append(stripped)
                prev_blank = False
        return "\n\n".join(out)


# ------------------------------------------------------------------ #
# Helpers                                                              #
# ------------------------------------------------------------------ #

def _clean(s: str) -> str:
    """Normalise whitespace and strip trailing footnote cruft."""
    s = re.sub(r"\s+", " ", s).strip()
    return s


def _render_table(rows: list[list[str]]) -> str:
    """Render a value-form table as an indented code block."""
    if not rows:
        return ""
    lines = ["```"]
    for row in rows:
        lines.append("  " + "  |  ".join(cell for cell in row if cell))
    lines.append("```")
    return "\n".join(lines)


# ------------------------------------------------------------------ #
# Entry point                                                          #
# ------------------------------------------------------------------ #

def strip(path: Path) -> str:
    parser = ChapterParser()
    parser.feed(path.read_text(encoding="utf-8", errors="replace"))
    return parser.text()


def main() -> None:
    if len(sys.argv) < 2:
        print(f"usage: {sys.argv[0]} <chapter.html>", file=sys.stderr)
        sys.exit(1)
    path = Path(sys.argv[1])
    if not path.exists():
        print(f"not found: {path}", file=sys.stderr)
        sys.exit(2)
    print(strip(path))


if __name__ == "__main__":
    main()
