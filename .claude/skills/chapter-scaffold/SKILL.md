---
name: chapter-scaffold
description: Set up a new chapter branch in the capital-simulator repo — creates the `volume-X/chapter-Y` branch off main (X ∈ {1,2,3}, Y = chapter number, no slug suffix), flips the roadmap row in `docs/architecture.md` for the right Vol. X section, and stubs a domain test file with a Marx textual fixture. Use this whenever the user says they're starting a new chapter, mentions "Vol. II Ch. 5" or "chapter N of volume X", or references the chapter workflow in CLAUDE.md. Do not use for service scaffolding (use store-scaffold) or for opening the end-of-chapter PR (use chapter-pr).
---

# chapter-scaffold

Sets up the workspace at the *start* of a new Capital chapter. The actual
domain implementation is chapter-specific — that's the user's job — but the
boilerplate around it is the same every time, and getting it wrong (branch
name, slug, roadmap row format) creates rework.

## When the user invokes this

They'll say something like:

- "Let's start Vol. II Ch. 1"
- "Kick off volume 2 chapter 5"
- "Set up Vol. III Ch. 13 — the law of the tendential fall"

**Gather both volume and chapter number.** If either isn't given, ask once.
Chapter numbers reset per volume, so "chapter 5" is ambiguous without the
volume.

## Volume → vault path + folder map

| Volume | Year | Vault prefix                                | Frontend folder                  |
|--------|------|---------------------------------------------|----------------------------------|
| I      | 1867 | `marx-engels/1867/capital-volume-i/`        | `web/src/chapters/vol1/`         |
| II     | 1885 | `marx-engels/1885/capital-volume-ii/`       | `web/src/chapters/vol2/`         |
| III    | 1894 | `marx-engels/1894/capital-volume-iii/`      | `web/src/chapters/vol3/`         |

The canonical slug for each chapter lives at the vault filename:
`<prefix>texts/NN-<slug>.md` and `<prefix>specs/NN-<slug>.spec.md`.
Match those exactly — don't invent a slug.

## Steps

Do these in order. Don't batch the git operations — if branch creation
fails (uncommitted changes, branch already exists) you want to stop before
touching anything else.

### 1. Verify clean tree and on main

```bash
git status --porcelain   # must be empty
git rev-parse --abbrev-ref HEAD   # must be "main"
```

If either check fails, stop and tell the user. Don't stash or auto-commit
their work.

### 2. Pull main, create the branch

Branch convention (CLAUDE.md): `volume-X/chapter-Y` where X ∈ {1,2,3} and
Y is the chapter number — **no slug suffix**. Slug lives in the vault
filenames and the component filename; the branch name stays terse.

```bash
git pull --ff-only origin main
git checkout -b volume-X/chapter-Y
```

If you hit a `.git/index.lock` permissions error (sandbox limit noted in
CLAUDE.md), `mv` the lock file out of the way:

```bash
mv .git/index.lock .git/index.lock.bak
```

then retry the git command.

### 3. Update the roadmap in docs/architecture.md

The roadmap now has three sub-tables: `### Volume I Roadmap`,
`### Volume II Roadmap`, `### Volume III Roadmap`. Find the row for **this
chapter, in the matching Volume X section** and flip its Status from
`⏳ Pending` (or `Pending`) to `In progress`. Don't flip a row in the
wrong volume — Vol. I Ch. 5, Vol. II Ch. 5, and Vol. III Ch. 5 are three
different things.

For chapters that touch a brand-new service for the first time, also note
this in the Status column ("In progress · scaffolds X-service").

### 4. Stub a test file with a Marx textual fixture

The repo convention is that every chapter's tests use Marx's own examples
as fixtures: `20 yards linen = 1 coat`, `1 quarter corn = x cwt iron`,
the circuit `M—C…P…C'—M'`, etc. This grounds the simulation in the text
and gives the test names a reading-list quality. Drop a `*_test.go` file
in the relevant `internal/<domain>/` package that imports the package and
has a single skipped test naming the chapter:

```go
package <domain>

import "testing"

// TestVolXChNN_<slug> is the textual anchor for Capital Vol. X, Ch. NN.
// Replace t.Skip with the chapter's first textual example as a fixture
// (e.g. "M — C(Lp+Mp) … P … C′ — M′" for Vol. II Ch. 1 §1).
func TestVolXChNN_<slug>(t *testing.T) {
	t.Parallel()
	t.Skip("TODO: implement Vol. X Ch. NN")
}
```

Skipped tests serve as a TODO list visible to `go test ./...` rather than
buried in a comment.

### 5. Don't commit

Leave the working tree dirty. The user wants to see the diff before any
commit, and signed commits fail in the sandbox anyway. Tell them what was
changed and what to do next:

> Branch `volume-X/chapter-Y` is set up. Roadmap row flipped to In progress
> in the Vol. X table, test stub in `services/<svc>/internal/<domain>/`.
> The chapter spec lives in the red-vault at
> `<vault-prefix>specs/NN-<slug>.spec.md` — pull it via `chapter-spec` for
> the implementation plan. Next: implement the domain types, register the
> React panel at `web/src/chapters/vol<V>/ChNN<Title>.tsx`, add a schema
> migration AND a paired seed migration
> (`NNNNN_v<V>_chNN_seed.sql`) with Marx-faithful exemplars so the
> dashboard comes up populated. Run `make vet test build` before
> committing — I can't run the Go toolchain here.

The seed migration is part of the chapter's deliverable, not a polish
step. `chapter-pr` will check for it before opening the PR (CLAUDE.md
Conventions → Seeds).

## What this skill does NOT do

- **Scaffold a new service.** If the chapter needs a brand-new service
  (e.g. finance-service for Vol. III), call out that service-scaffolding
  is a separate concern and ask whether to proceed manually.
- **Generate domain types.** What goes in `commodity.go`, `value.go`, etc.
  is the heart of the chapter and is judgment-laden — it's Marx's
  argument turned into types. Don't pre-populate it.
- **Open a PR.** That's `chapter-pr`'s job, run at the end of the chapter.
- **Author the chapter spec.** Specs are pre-authored in the vault. Use
  `chapter-spec` to read one, not to write one. If the spec doesn't exist
  yet (Vol. II and III specs are not all written), surface that as a gap
  rather than fabricating one.

## Anti-patterns

- Don't create per-service `go.mod` files — the repo is intentionally a
  single module.
- Don't add a TypeScript router or new web routes here — one `App.tsx`
  composes panels until a chapter forces otherwise (CLAUDE.md).
- Don't try to drop a chapter HTML or spec into `chapters/volume-1/` —
  that directory was retired. Source text and specs live in the red-vault.
- Don't use the old branch convention `chapter-NN-<slug>`. CLAUDE.md
  specifies `volume-X/chapter-Y` (no slug suffix) for all volumes.
- Don't land a chapter component at `web/src/chapters/Ch<NN>...tsx` (the
  flat layout was retired in foundation Phase 1). Land it under
  `web/src/chapters/vol<V>/`.
- Don't omit the `v<V>` token from migration filenames. Goose tracks by
  the integer prefix only, but the token is what disambiguates Vol. I
  Ch. 5 migrations from Vol. II Ch. 5 migrations during review.
