---
name: chapter-pr
description: Run the end-of-chapter checklist for capital-simulator and open the PR. Verifies the branch name matches chapter-NN-<slug>, the chapters/volume-1/NN-<slug>.html exists with real content, the docs/architecture.md roadmap row is flipped to Done, and that the diff against main is non-trivial. Then drafts a multi-paragraph conventional commit body referencing the chapter and runs gh pr create. Use when the user says "open the chapter PR", "ship this chapter", "we're done with chapter N", or asks to wrap up / finalize a Capital chapter. Do not use to scaffold a new chapter (use chapter-scaffold).
---

# chapter-pr

The closeout step in the chapter workflow described in CLAUDE.md. By the
time this runs, the user has:

1. Implemented the chapter's domain logic
2. Pasted the real Marx text into `chapters/volume-1/NN-<slug>.html`
3. Run `make vet test build` and `cd web && npm run lint && npm run build`
   locally (the sandbox can't)

This skill does the last-mile verification, then opens the PR.

## When the user invokes this

- "Open the chapter PR"
- "Ship this chapter"
- "We're done with Ch. 3, push it"
- "Wrap up the chapter"

If unsure whether the chapter is actually done, ask: "Have you run
`make vet test build` and the web build? Both need to pass before this -
I can't verify them here." Wait for confirmation before proceeding.

## Steps

### 1. Run the precheck script

```bash
bash .claude/skills/chapter-pr/scripts/check.sh
```

It verifies:

- Current branch matches `chapter-NN-<slug>`
- `chapters/volume-1/NN-<slug>.html` exists and is larger than the placeholder
  stub (real Marx content is ~50-100KB; the stub from `chapter-scaffold`
  is well under 1KB)
- `docs/architecture.md` roadmap row for this chapter shows status
  containing `Done` or a checkmark, not `In progress` / `Next` /
  `Pending`
- The diff against `origin/main` is non-empty
- There is at least one chapter-relevant commit on the branch beyond
  what's on `main`

If any check fails, the script prints which one and exits non-zero. Stop
and surface the failure to the user; don't try to "fix" by editing the
roadmap or HTML on their behalf - those are content decisions.

Before declaring the precheck "passed," also confirm by eye:

- **Seed migration shipped.** If the chapter introduces a new domain
  type (new table, or new fields on an existing one), there must be a
  `services/<svc>/internal/store/migrations/NNNNN_chNN_seed.sql` that
  inserts Marx-faithful exemplars and has a complete `-- +goose Down`
  that DELETEs every seeded id. The dashboard must come up populated on
  a fresh MySQL volume. If the seed is missing, write it before opening
  the PR — empty panels on first boot are a regression (CLAUDE.md
  Conventions → Seeds).

### 2. Draft the commit + PR body

If commits aren't already pushed (which is the common case in this
sandbox - signed commits fail here, so the user typically does the final
commit themselves on their machine), draft a conventional commit body
the user can paste:

```
feat(<svc>): implement Capital Vol. I, Ch. NN - <Title>

<2-4 sentence summary of what was built and the economic concept>

Highlights:
- §1 ...
- §2 ...
- §3 ...

Refs: chapters/volume-1/NN-<slug>.html
```

Use the section headings from Marx's chapter (visible in the chapter
HTML <h2>/<h3>) to drive the bullet list. The `<svc>` scope is the
primary service touched (e.g. `commodity` for Ch. 1, `market` for
Ch. 2-3).

Look at the existing `9560923 feat(commodity): implement Capital Vol. I,
Ch. 1 - The Commodity` commit (`git log --format=%B -1 9560923` if it
still exists, otherwise reference the structure above) for tone.

### 3. Open the PR

Once commits are pushed, run:

```bash
gh pr create --base main --title "<conventional commit title>" --body "$(cat <<'EOF'
## Chapter

Vol. I, Ch. NN - <Title>

## Summary

<2-4 sentence economic argument summary>

## Services touched

- [x] <service>
- [x] web (UI)
<...>

## Chapter HTML

chapters/volume-1/NN-<slug>.html

## How I tested

- [x] `make vet`
- [x] `make test`
- [x] `make build`
- [ ] `docker compose up --build` smoke
- [ ] Other:

## Notes for review

<anything notable - tricky tests, intentional patterns, open questions>

🤖 Generated with [Claude Code](https://claude.com/claude-code)
EOF
)"
```

Match the body to the repo's `.github/pull_request_template.md` -
sections and ordering. Tick the boxes the user actually ran; leave
unchecked the ones they explicitly said they skipped.

### 4. Don't do destructive things

- Don't `git push --force`. The chapter branches are fresh and shouldn't
  need it; if they do, surface the situation to the user.
- Don't commit or push from the sandbox without explicit user
  confirmation - signed commits fail locally and pushing unsigned
  commits would break the repo's signing convention.
- Don't merge the PR. Let the user merge after review.

## What this skill does NOT do

- **Run the Go or web test suites.** The sandbox can't. Trust the user's
  confirmation that they passed locally.
- **Edit the chapter HTML or architecture.md.** Those are content
  decisions; this skill verifies them, doesn't write them.
- **Resolve merge conflicts on main.** If the precheck reveals the
  branch is behind, tell the user and let them rebase.

## Anti-patterns

- Don't open a PR with the placeholder chapter HTML from
  `chapter-scaffold` - the precheck guards this, but if the user
  insists, they should know the chapter HTML is part of the merge
  artifact (CLAUDE.md `Done is when`).
- Don't bundle multiple chapters in one PR. One chapter, one branch,
  one PR (CLAUDE.md chapter workflow).
- Don't rewrite the conventional commit format. The PR template fills
  from the commit body; deviating from `feat(<svc>): ...` breaks that.
