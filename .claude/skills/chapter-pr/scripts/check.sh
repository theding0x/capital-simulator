#!/usr/bin/env bash
# End-of-chapter precheck for capital-simulator.
#
# Verifies branch name, architecture.md status row, and that there's a
# real diff to ship. Exits non-zero on any failure with a clear message
# naming the offending check. Chapter source text and spec live in the
# red-vault Obsidian vault, not in this repo, so they are not validated
# here.

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
if [[ $branch =~ ^chapter-([0-9]{2})-([a-z0-9-]+)$ ]]; then
  chapter_num="${BASH_REMATCH[1]}"
elif [[ $branch =~ ^volume-[0-9]+/chapter-([0-9]+)(-([a-z0-9-]+))?$ ]]; then
  chapter_num=$(printf '%02d' "${BASH_REMATCH[1]}")
else
  fail "branch '$branch' does not match chapter-NN-<slug> or volume-N/chapter-NN"
fi
ok "branch $branch (chapter $chapter_num)"

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
