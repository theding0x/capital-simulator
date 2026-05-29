---
name: issue-to-completion
description: Use when picking up an open GitHub issue to work through from selection to a merged PR — finding the oldest actionable issue, triaging it, implementing the fix, and driving CI to green.
---

# issue-to-completion

Workflow for taking the oldest unassigned GitHub issue through to a merged PR. Three phases: **Pick** → **Implement** → **Ship**.

## Phase detection

Announce the phase before acting.

| Condition | Phase |
|-----------|-------|
| No branch or PR exists for the issue | Pick |
| Branch exists, work in progress | Implement |
| Implementation done, PR open | Ship |

---

## Phase 1 — Pick the issue

### Find the oldest actionable issue

```bash
gh issue list --json number,title,createdAt,assignees,labels --limit 50
```

Sort by `createdAt` ascending. For each candidate from oldest to newest, check:

```bash
gh issue view <number> --json body,linkedPullRequests,assignees
```

**Skip the issue and try the next if:**

- It has an assignee
- It has a linked PR (already in flight)
- Its body contains any of: "tracking only", "out of scope", "do not implement now", "hold off until", "deferred", "blocked by"

The first issue that passes all checks is the target. If every issue is skipped, surface the list to the user and ask which to tackle.

### Read the full issue

Before writing a line of code, understand:

- What specifically needs to change (file, line numbers if given)
- Acceptance criteria (explicit or implicit from the description)
- Which service/package is affected
- Whether tests are called out

If the issue is ambiguous about **what** to change (not just **why**), stop and ask the user before branching.

### Branch

Naming: `issue/<number>-<slug>` where slug is the first 3–5 words of the title, lowercased, hyphenated, alphanumeric only.

```bash
git checkout -b issue/<number>-<slug>
git push -u origin issue/<number>-<slug>
```

---

## Phase 2 — Implement

Follow all conventions in CLAUDE.md. Relevant reminders:

- **Go imports**: stdlib / third-party / local — three groups, blank lines between.
- **JSON tags**: `snake_case` only.
- **IDs**: `crypto/rand` hex via `NewXxxID()`. No `google/uuid`.
- **Errors**: store sentinels `ErrNotFound` / `ErrAlreadyExists`; HTTP layer maps via `errors.Is`.
- **Tests**: `t.Parallel()` first; use Marx's textual examples as fixtures.
- **Migrations**: append-only — add a new numbered file, never edit an existing one.
- **No inline DDL**: all `CREATE`/`ALTER` in `.sql` files under `migrations/`.

### Common issue types

| Issue label / signal | What to change |
|----------------------|----------------|
| Bug + proposed fix in body | Implement the suggested fix; add the test cases listed in the issue |
| `refactor` / duplicate types | Consolidate; update all callers; add/adjust tests to confirm behaviour is unchanged |
| `docs` / comment / docstring | Edit the comment in the target file only |
| `chore` / gitignore / tooling | Config files only; no domain logic |
| `a11y` / accessibility | React component + CSS; check with Playwright snapshot |

### Verify before committing

```bash
make vet test build
```

If any file under `web/` changed:
```bash
cd web && npm run lint && npm run build
```

If any new UI panel or interactive component was added or changed:
- Start the dev server
- Open the feature in a browser
- Verify the golden path and any edge cases named in the issue

---

## Phase 3 — Ship

### Open the PR

```bash
gh pr create --base main \
  --title "<type>(<scope>): <short description>" \
  --body "$(cat <<'EOF'
## Issue

Closes #<number>

## Summary

<2–3 sentences: what the issue identified, what this PR does about it>

## Changes

- `path/to/file.go` — description of change
- `path/to/file_test.go` — tests added/updated

## How I tested

- [ ] `make vet`
- [ ] `make test`
- [ ] `make build`
- [ ] `cd web && npm run lint && npm run build` (if web changed)
- [ ] Playwright / browser check (if UI changed)

🤖 Generated with [Claude Code](https://claude.com/claude-code)
EOF
)"
```

Conventional commit type in the title must match the issue: `fix` for bugs, `refactor` for consolidation, `docs` for comments, `chore` for tooling, `feat` for new behaviour.

### CI loop

```bash
gh pr checks <number>
```

- **In progress** — report status; use `ScheduleWakeup` (60–270 s) rather than a sleep loop if waiting.
- **Failed** — fetch the log: `gh run view <run-id> --log-failed`. Diagnose, fix, commit, push. Loop back.
- **All passing** — `gh pr ready <number>`. Tell the user: "PR #N is ready to merge into main."

Do not merge the PR. The user merges after review.

---

## What NOT to do

- Don't work a tracking-only issue — those exist to record decisions, not trigger code changes.
- Don't bundle multiple issues in one PR.
- Don't create a new migration unless the issue explicitly requires schema changes.
- Don't use chapter branch naming (`volume-X/chapter-Y`) — issue branches use `issue/<number>-<slug>`.
- Don't skip pre-commit hooks (`--no-verify`) or signing to fix a failing commit.
- Don't mark the PR ready before all CI checks pass.
- Don't add features, refactor surrounding code, or clean up unrelated files — fix exactly what the issue describes.
