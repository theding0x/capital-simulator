#!/usr/bin/env python3
"""Fetch all Capital Vol. I chapters from marxists.org.

Starts at ch01.htm, follows "Next:" links to the end of the volume, and
saves each chapter as raw HTML into the red-vault staging directory at
`marx-engels/1867/capital-volume-i/raw-html/NN-<slug>.html`. The polished
markdown in the sibling `texts/` directory is hand-tuned and is NOT
written by this script — staged HTML is raw material for whatever
conversion pipeline you run separately.

Usage:
    python3 scripts/fetch_chapters.py

Set the env var `RED_VAULT_DIR` to point at a different vault root if
needed (defaults to the maintainer's local clone path). Existing files
are skipped. Politely rate-limited to 1 request/second.
"""
from __future__ import annotations

import os
import re
import sys
import time
import urllib.request
from html.parser import HTMLParser
from pathlib import Path

DEFAULT_RED_VAULT_DIR = Path(
    "/mnt/c/Users/AaronHulse/Repositories/github.com/theding0x/red-vault"
)
VAULT_DIR = Path(os.environ.get("RED_VAULT_DIR", str(DEFAULT_RED_VAULT_DIR)))
OUTPUT_DIR = VAULT_DIR / "marx-engels" / "1867" / "capital-volume-i" / "raw-html"
START_URL = "https://www.marxists.org/archive/marx/works/1867-c1/ch01.htm"
DELAY_S = 1.0


# ---------------------------------------------------------------------------
# Fetching
# ---------------------------------------------------------------------------

def fetch(url: str) -> str:
    req = urllib.request.Request(
        url,
        headers={"User-Agent": "Mozilla/5.0 (research crawler; capital-simulator project)"},
    )
    with urllib.request.urlopen(req, timeout=30) as resp:
        charset = "windows-1252"
        ct = resp.headers.get_content_charset()
        if ct:
            charset = ct
        return resp.read().decode(charset, errors="replace")


# ---------------------------------------------------------------------------
# HTML parsing
# ---------------------------------------------------------------------------

class _MetaParser(HTMLParser):
    """Extracts the chapter h3 title and the next-chapter href.

    Live marxists.org pages use <h3> for the chapter title, e.g.:
        <h3>Chapter Seven: The Labour-Process ...</h3>
    The next-chapter link is a relative href, e.g.:
        <a href="ch08.htm">Next: Chapter Eight: ...</a>
    """

    BASE = "https://www.marxists.org/archive/marx/works/1867-c1/"

    def __init__(self) -> None:
        super().__init__()
        self.title: str = ""
        self.next_url: str = ""
        self._in_h3 = False
        self._buf: list[str] = []
        self._pending_href: str = ""

    def handle_starttag(self, tag: str, attrs: list[tuple[str, str | None]]) -> None:
        if tag == "h3" and not self.title:
            # Only grab the first h3 — that's the chapter title
            self._in_h3 = True
            self._buf = []
        if tag == "a":
            attr = dict(attrs)
            href = attr.get("href", "")
            # Only interested in same-directory chapter links (ch08.htm, not ../foo/ch01.htm)
            if re.fullmatch(r"ch\d+\.htm", href):
                self._pending_href = href
            else:
                self._pending_href = ""
        elif tag != "h3":
            self._pending_href = ""

    def handle_endtag(self, tag: str) -> None:
        if tag == "h3" and self._in_h3:
            self._in_h3 = False
            self.title = " ".join(self._buf).strip()
            self._buf = []

    def handle_data(self, data: str) -> None:
        if self._in_h3:
            self._buf.append(data)
        # Capture next-chapter link by its visible text starting with "Next:"
        if self._pending_href and data.strip().startswith("Next:"):
            self.next_url = self.BASE + self._pending_href
            self._pending_href = ""


def parse_page(html: str) -> tuple[str, str]:
    """Return (raw_title, next_url) extracted from the page HTML."""
    p = _MetaParser()
    p.feed(html)
    return p.title, p.next_url


# ---------------------------------------------------------------------------
# Slug / filename helpers
# ---------------------------------------------------------------------------

def slugify(title: str) -> str:
    # Strip "Chapter <anything>: " prefix (handles "Chapter Twenty-One: ...")
    if ":" in title:
        title = title.split(":", 1)[1]
    # Remove footnote references like [1], [2]
    title = re.sub(r"\s*\[\d+\]", "", title)
    # Lowercase, replace runs of non-alphanumeric chars with hyphens
    slug = re.sub(r"[^a-z0-9]+", "-", title.lower())
    return slug.strip("-")


def chapter_num(url: str) -> int | None:
    m = re.search(r"ch(\d+)\.htm", url)
    return int(m.group(1)) if m else None


# ---------------------------------------------------------------------------
# Main loop
# ---------------------------------------------------------------------------

def main() -> None:
    OUTPUT_DIR.mkdir(parents=True, exist_ok=True)
    print(f"Output dir: {OUTPUT_DIR}")
    url: str = START_URL
    visited: set[str] = set()
    saved = 0
    skipped = 0

    while url:
        if url in visited:
            print(f"Loop detected at {url}, stopping.", file=sys.stderr)
            break
        visited.add(url)

        num = chapter_num(url)
        if num is None:
            print(f"Cannot determine chapter number from {url}, stopping.", file=sys.stderr)
            break

        print(f"Ch.{num:02d}  {url}")

        try:
            html = fetch(url)
        except Exception as exc:
            print(f"  ERROR fetching: {exc}", file=sys.stderr)
            break

        raw_title, next_url = parse_page(html)

        if not raw_title:
            print("  WARNING: no h2 title found; using url-derived name")
            slug = f"ch{num:02d}"
        else:
            slug = slugify(raw_title)
            subtitle = re.sub(r"^Chapter\s+\w+\s*[:\-]\s*", "", raw_title, flags=re.IGNORECASE)
            print(f"  Title: {subtitle}")
            print(f"  Slug:  {slug}")

        out_path = OUTPUT_DIR / f"{num:02d}-{slug}.html"

        existing = list(OUTPUT_DIR.glob(f"{num:02d}-*.html"))
        if existing:
            print(f"  Skipped (exists): {existing[0].name}")
            skipped += 1
        else:
            out_path.write_text(html, encoding="utf-8")
            print(f"  Saved:  {out_path.name}")
            saved += 1

        url = next_url
        if url:
            time.sleep(DELAY_S)

    print(f"\nDone. Saved: {saved}, skipped: {skipped}.")


if __name__ == "__main__":
    main()
