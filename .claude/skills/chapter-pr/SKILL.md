---
name: chapter-pr
description: Manage the per-chapter PR for capital-simulator across its three phases — (A) open a draft PR upfront populated from the chapter spec, (B) sync the PR description after commits land so it reflects what actually shipped, (C) wait for GitHub Actions, fix failures, and mark the PR ready for merge. Auto-detects the phase from `gh pr view` for the current branch, and the volume from the branch name. Use when the user says "open the draft PR", "sync the chapter PR", "update the PR description", "mark this PR ready", "wrap up the chapter", "is this chapter shippable", or any other phase of the chapter PR flow. Branch convention is `volume-X/chapter-Y` (no slug). Do not use to scaffold a new chapter (use chapter-scaffold) or to load the chapter spec (use chapter-spec).
---

# chapter-pr

The PR-management skill for the chapter workflow described in CLAUDE.md. The
workflow opens a **draft** PR upfront from the spec (step 2), iterates on
commits and re-syncs the description (step 5), then waits on CI and marks
the PR ready (steps 6–8). This skill handles all three phases and resolves
the volume from the branch name (`volume-X/chapter-Y`).

## Volume → vault path

The chapter spec lives in the red-vault at a volume-keyed path:

| Volume | Year | Vault prefix                                |
|--------|------|---------------------------------------------|
| I      | 1867 | `marx-engels/1867/capital-volume-i/`        |
| II     | 1885 | `marx-engels/1885/capital-volume-ii/`       |
| III    | 1894 | `marx-engels/1894/capital-volume-iii/`      |

Resolve the volume from the branch's `volume-X/...` segment. The spec
filepath is `<prefix>specs/NN-<slug>.spec.md`.

## Phases at a glance

| Phase | Workflow step | Trigger condition                                              | What this skill does                                                                                                       |
|-------|---------------|----------------------------------------------------------------|----------------------------------------------------------------------------------------------------------------------------|
| A     | step 2        | branch exists, no PR yet for the branch                        | Open a **draft** PR populated from the vault spec + the PR template                                                        |
| B     | step 5        | draft PR exists; new commits since the description last synced | Re-derive Services touched + Summary from the actual diff; `gh pr edit --body`                                             |
| C     | steps 6–8     | implementation done; awaiting / chasing CI                     | Run the precheck, poll `gh pr checks`, fix failures, then `gh pr ready`                                                    |

## Detect the phase and volume first

Before acting, run:

```bash
gh pr view --json number,isDraft,state,headRefName 2>/dev/null
```

Parse the branch name `volume-X/chapter-Y` to extract:

- **X** → the volume (1, 2, or 3) → resolves the vault prefix from the
  table above
- **Y** → the chapter number → zero-padded for the spec filename match

Then:

- No PR found → Phase A.
- PR exists, `isDraft: true` → Phase B if the user wants to sync the description, Phase C if they want to ship.
- PR exists, `isDraft: false` → Phase C (already marked ready; only thing left is the merge, which the user owns).

Announce which phase + volume you're entering before doing anything else.

## When the user invokes this

| User says                                            | Phase |
|------------------------------------------------------|-------|
| "Open the draft PR for chapter N"                    | A     |
| "Kick off chapter N's PR"                            | A     |
| "Sync the chapter PR" / "Update the PR description"  | B     |
| "Mark this PR ready" / "Wrap up the chapter"         | C     |
| "Is this chapter shippable?" / "Are checks green?"   | C     |

If ambiguous, infer from the phase detection above.

## Phase A — Open the draft PR

Prereqs: branch `volume-X/chapter-Y` exists; the vault spec
`<prefix>specs/NN-<slug>.spec.md` exists (where `NN` is the zero-padded
chapter number and `<prefix>` is resolved from the volume). Specs are
pre-authored — if one is genuinely missing for the chapter, surface that
as a gap to the user rather than fabricating one.

### Steps

1. Verify branch name matches `^volume-[0-9]+/chapter-[0-9]+$`. Bail if
   not. Extract volume (X) and chapter (Y).
2. Resolve `<prefix>` from the table above. Locate the chapter spec by
   listing the vault specs directory and matching the `NN-` prefix:
   ```
   mcp__obsidian__obsidian_list_files_in_dir
     dirpath: <prefix>specs
   ```
   Then fetch it via `mcp__obsidian__obsidian_get_file_contents`. If the
   matching file is absent, stop and tell the user (Vol. II / III specs
   are not all authored yet — see Phase 5 of the multi-volume plan).
3. Read the spec frontmatter (`title`, `primary_service`) and the
   **Scope → This chapter builds** section to extract planned domain
   types, endpoints, and UI work.
4. Push the branch to origin (`git push -u origin volume-X/chapter-Y`).
   GitHub needs a remote branch to attach the PR to. If there are no
   commits yet, make an empty marker commit first:
   `git commit --allow-empty -m "chore(vol-X/chN): kick off chapter N branch"`.
5. Run `gh pr create --draft --base main --title ... --body ...` with
   the body template below.

### PR body template (Phase A)

Match `.github/pull_request_template.md` sections. All `How I tested`
boxes start unchecked.

```markdown
## Chapter

Vol. X, Ch. Y — <Title from spec>

## Summary (planned)

<2–4 sentences drawn from the spec describing the economic argument and
what this chapter will add to the simulation. Mark as "(planned)" until
Phase B refreshes it.>

## Services touched

- [x] <primary_service from spec frontmatter>
- [ ] api-gateway
- [ ] commodity-service
- [ ] agent-service
- [ ] market-service
- [ ] simulation-engine
- [ ] finance-service
- [ ] web (UI)
- [ ] pkg/* (shared)
- [ ] deploy/k8s
- [ ] docker-compose.yml

## How I tested

- [ ] `make vet`
- [ ] `make test`
- [ ] `make build`
- [ ] `docker compose up --build` smoke
- [ ] Other (describe):

## Planned changes (from spec)

**New domain types**
- `Type1` — short description
- `Type2` — short description

**New HTTP endpoints**
- `POST /v1/...` — description
- `GET /v1/...` — description

**UI**
- "Ch. Y — <Title>" panel under `web/src/chapters/vol<X>/`

## Notes for review

<Anything the spec flagged as deferred to a later chapter, or design
tension worth raising.>

🤖 Generated with [Claude Code](https://claude.com/claude-code)
```

## Phase B — Sync the PR description

Triggered when commits land on the branch and the draft PR's planned
bullets are now out of date.

### Steps

1. `gh pr view --json number,body,headRefName` to get the current PR
   number and body. Re-extract volume + chapter from `headRefName`.
2. `git diff --name-only origin/main..HEAD` to list every file that
   changed since branching.
3. Group the file list by service to derive **Services touched**:
   - `services/<svc>/...` → that service
   - `web/...` → `web (UI)`
   - `pkg/...` → `pkg/* (shared)`
   - `deploy/k8s/...` → `deploy/k8s`
   - `docker-compose.yml` → `docker-compose.yml`
4. Re-read the relevant commit message bodies (`git log origin/main..HEAD --format=%B`)
   to derive the **Summary** bullets — the conventional commit body is
   the canonical source.
5. Replace **Planned changes (from spec)** with a refreshed **Summary**
   section reflecting what actually shipped. Keep the **Notes for review**
   section, appending anything new the user wants flagged for reviewers.
6. Tick **How I tested** boxes only for what the user has explicitly
   confirmed they ran locally. Leave the rest unchecked — CI in Phase C
   covers the rest.
7. `gh pr edit <number> --body "$(cat <<'EOF' ... EOF)"` with the new body.

**Do not flip the PR out of draft in Phase B.** Mark-ready belongs to Phase C.

## Phase C — Wait for CI, fix, mark ready

### Steps

1. Run the precheck script (it covers branch name, architecture roadmap
   row marked Done in the **right Volume's table**, branch ahead of
   main, real diff):

   ```bash
   sed -i 's/\r$//' .claude/skills/chapter-pr/scripts/check.sh   # WSL CRLF guard
   bash .claude/skills/chapter-pr/scripts/check.sh 2>&1
   ```

   If a check fails, surface it to the user. Roadmap-row flips are
   content decisions — confirm before editing on the user's behalf.

2. Confirm by eye that the **seed migration shipped** for any new
   domain type the chapter introduced
   (`services/<svc>/internal/store/migrations/NNNNN_v<X>_chNN_seed.sql`,
   with a `-- +goose Down` deleting every seeded id). Empty panels on
   first boot are a regression — see CLAUDE.md → Seeds.

3. Run `gh pr checks <number>` to inspect CI status:

   - **In progress** — tell the user the status. Don't block on it; the
     user can re-invoke this skill once checks complete. If asked to
     wait, use a `ScheduleWakeup` (60–270s) or `Monitor` rather than
     a blind sleep loop.
   - **Failed** — fetch the failing log:
     ```bash
     gh run view <run-id> --log-failed
     ```
     Surface the failure to the user, diagnose the root cause, fix in
     code, commit (signed), push. Loop back to step 3. **Do not** "fix"
     by skipping the check, deleting the test, or using `--no-verify`.
   - **All passing** — run `gh pr ready <number>` to flip the PR out of
     draft, then tell the user: "PR #N is ready to merge into main."

4. Do not merge the PR yourself. The user merges after review.

## Don't do destructive things

- Don't `git push --force`. Chapter branches shouldn't need rewriting;
  if they do, surface to the user.
- Don't skip pre-commit hooks (`--no-verify`) or signing
  (`--no-gpg-sign`).
- Don't merge the PR.
- Don't "fix" CI failures by deleting the failing check, stubbing the
  test, or pinning around it — fix the underlying issue.

## What this skill does NOT do

- **Run the Go or web test suites locally.** The sandbox can't. CI in
  Phase C covers it.
- **Edit the chapter spec or source text.** Those live in the red-vault
  and the user owns them. If a spec needs an update mid-PR, the user
  edits it in Obsidian.
- **Write the architecture.md roadmap row.** Content decision.
- **Resolve merge conflicts on main.** If the branch is behind, tell
  the user and let them rebase.
- **Scaffold a new chapter** — use `chapter-scaffold`.
- **Load the chapter spec** — use `chapter-spec`.

## Anti-patterns

- Don't open the PR as non-draft in Phase A. The whole point of the
  workflow is a paper trail spec → planned → shipped.
- Don't bundle multiple chapters in one PR. One chapter, one branch,
  one PR.
- Don't rewrite the conventional commit format. The PR description in
  Phase B fills from the commit body; deviating from `feat(<svc>): ...`
  breaks that.
- Don't reference `chapters/volume-1/` — that directory was retired.
  Spec + chapter text live in the red-vault.
- Don't omit the `v<X>` token from the seed-migration check. Migration
  filenames in foundation Phase 1+ carry `NNNNN_v<X>_ch<NN>_<slug>.sql`.
- Don't validate the roadmap row against just `^| Ch. N` — Vol. I,
  Vol. II, and Vol. III each have a row for Ch. 1. The precheck must
  scope its search to the volume's section (`### Volume X Roadmap`).
