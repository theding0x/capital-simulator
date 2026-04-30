#!/usr/bin/env bash
# End-of-chapter precheck for capital-simulator.
#
# Verifies branch name, chapter HTML, architecture.md status row, and that
# there's a real diff to ship. Exits non-zero on any failure with a clear
# message naming the offending check.

set -u
shopt -s nullglob

cd "$(git rev-parse --show-toplevel 2>/dev/null)" || {
  echo "FAIL: not inside a git repository" >&2
  exit 2
}

fail() {
  echo "FAIL: $*" >&2
  exit 1
}

ok() {
  echo "ok: $*"
}

branch=$(git rev-parse --abbrev-ref HEAD)
if [[ ! $branch =~ ^chapter-([0-9]{2})-([a-z0-9-]+)$ ]]; then
  fail "branch '$branch' does not match chapter-NN-<slug>"
fi
chapter_num="${BASH_REMATCH[1]}"
slug="${BASH_REMATCH[2]}"
ok "branch $branch (chapter $chapter_num, slug $slug)"

html="chapters/${chapter_num}-${slug}.html"
if [[ ! -f $html ]]; then
  fail "chapter HTML missing: $html"
fi
size=$(wc -c <"$html" | tr -d '[:space:]')
# The chapter-scaffold stub is well under 1KB; real Marx content is 50KB+.
# 5KB is a generous threshold that catches the stub without false-flagging
# a short chapter.
if (( size < 5000 )); then
  fail "chapter HTML at $html is only $size bytes — looks like the placeholder stub. Paste the real chapter text before opening the PR."
fi
ok "chapter HTML present ($size bytes)"

arch="docs/architecture.md"
if [[ ! -f $arch ]]; then
  fail "missing $arch"
fi

# Find the row for this chapter in the roadmap table. Tolerate "Ch. N",
# "Ch. NN", or ranges like "Ch. 2-3".
chapter_n_no_pad=$((10#$chapter_num))
row=$(grep -E "^\| Ch\. ${chapter_n_no_pad}([ -]|$|\|)" "$arch" || true)
if [[ -z $row ]]; then
  # Try matching ranges that cover this chapter (e.g. "Ch. 2-3" for ch 02).
  row=$(awk -F'|' -v n="$chapter_n_no_pad" '
    /^\| Ch\./ {
      gsub(/^[ ]*Ch\. */, "", $2);
      split($2, parts, "-");
      lo = parts[1] + 0;
      hi = (parts[2] != "" ? parts[2] + 0 : lo);
      if (n >= lo && n <= hi) { print $0; exit }
    }' "$arch")
fi
if [[ -z $row ]]; then
  fail "no roadmap row in $arch covers chapter $chapter_n_no_pad"
fi
status_col=$(awk -F'|' '{print $3}' <<<"$row")
if [[ $status_col != *Done* && $status_col != *✅* ]]; then
  fail "roadmap row for chapter $chapter_n_no_pad in $arch is not marked Done — got '$status_col'"
fi
ok "architecture.md roadmap row marked Done"

# Verify diff against origin/main exists. If origin/main isn't fetched,
# fall back to local main.
base=""
if git rev-parse --verify --quiet origin/main >/dev/null; then
  base="origin/main"
elif git rev-parse --verify --quiet main >/dev/null; then
  base="main"
else
  fail "cannot find origin/main or main to diff against"
fi

ahead=$(git rev-list --count "$base..HEAD")
if (( ahead == 0 )); then
  fail "branch is not ahead of $base — nothing to ship"
fi
ok "branch is $ahead commit(s) ahead of $base"

changed=$(git diff --name-only "$base..HEAD" | wc -l | tr -d '[:space:]')
if (( changed == 0 )); then
  fail "no file diff against $base"
fi
ok "$changed file(s) changed vs $base"

echo
echo "all checks passed — ready to open PR for chapter $chapter_n_no_pad"
