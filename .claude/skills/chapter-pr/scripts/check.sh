#!/usr/bin/env bash
# End-of-chapter precheck for capital-simulator.
#
# Verifies branch name, architecture.md status row IN THE RIGHT VOLUME'S
# section, and that there's a real diff to ship. Exits non-zero on any
# failure with a clear message naming the offending check. Chapter source
# text and spec live in the red-vault Obsidian vault, not in this repo,
# so they are not validated here.

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

# Roman numeral mapping for the roadmap section header.
roman_for() {
  case "$1" in
    1) echo "I" ;;
    2) echo "II" ;;
    3) echo "III" ;;
    *) echo "?" ;;
  esac
}

branch=$(git rev-parse --abbrev-ref HEAD)
volume_num=""
if [[ $branch =~ ^chapter-([0-9]{2})-([a-z0-9-]+)$ ]]; then
  # Legacy pre-foundation-Phase-1 branch form. Default to Vol. I.
  chapter_num="${BASH_REMATCH[1]}"
  volume_num=1
elif [[ $branch =~ ^volume-([0-9]+)/chapter-([0-9]+)(-([a-z0-9-]+))?$ ]]; then
  volume_num="${BASH_REMATCH[1]}"
  chapter_num=$(printf '%02d' "${BASH_REMATCH[2]}")
else
  fail "branch '$branch' does not match chapter-NN-<slug> or volume-N/chapter-NN"
fi
volume_roman=$(roman_for "$volume_num")
if [[ $volume_roman == "?" ]]; then
  fail "volume $volume_num is not in {1,2,3}"
fi
ok "branch $branch (Volume $volume_roman, chapter $chapter_num)"

arch="docs/architecture.md"
if [[ ! -f $arch ]]; then
  fail "missing $arch"
fi

# Extract just the Volume X Roadmap section so the chapter-row search
# can't match a row in the wrong volume's table.
section=$(awk -v vol="$volume_roman" '
  $0 ~ "^### Volume " vol " Roadmap" { in_sec = 1; next }
  in_sec && /^### / { in_sec = 0 }
  in_sec { print }
' "$arch")
if [[ -z $section ]]; then
  fail "no '### Volume $volume_roman Roadmap' section found in $arch"
fi

# Find the row for this chapter inside the volume section. Tolerate
# "Ch. N", "Ch. NN", or ranges like "Ch. 2-3" / "Ch. 8-9".
chapter_n_no_pad=$((10#$chapter_num))
row=$(grep -E "^\| Ch\. ${chapter_n_no_pad}([ -]|$|\|)" <<<"$section" || true)
if [[ -z $row ]]; then
  # Try matching ranges that cover this chapter (e.g. "Ch. 2-3" for ch 02).
  row=$(awk -F'|' -v n="$chapter_n_no_pad" '
    /^\| Ch\./ {
      gsub(/^[ ]*Ch\. */, "", $2);
      split($2, parts, "-");
      lo = parts[1] + 0;
      hi = (parts[2] != "" ? parts[2] + 0 : lo);
      if (n >= lo && n <= hi) { print $0; exit }
    }' <<<"$section")
fi
if [[ -z $row ]]; then
  fail "no roadmap row in Volume $volume_roman section of $arch covers chapter $chapter_n_no_pad"
fi
status_col=$(awk -F'|' '{print $3}' <<<"$row")
if [[ $status_col != *Done* && $status_col != *✅* ]]; then
  fail "Volume $volume_roman roadmap row for chapter $chapter_n_no_pad in $arch is not marked Done — got '$status_col'"
fi
ok "architecture.md Volume $volume_roman roadmap row marked Done"

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
echo "all checks passed — ready to open PR for Vol. $volume_roman chapter $chapter_n_no_pad"
