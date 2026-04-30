---
name: chapter-scaffold
description: Set up a new chapter branch in the capital-simulator repo - creates the chapter-NN-<slug> branch off main, drops the placeholder chapters/NN-<slug>.html, adds the roadmap row to docs/architecture.md, and stubs a domain test file with a Marx textual fixture. Use this whenever the user says they're starting a new chapter, mentions "chapter N", asks to scaffold/begin/kick off a Capital chapter, or references the chapter workflow described in CLAUDE.md. Do not use for service scaffolding (use store-scaffold) or for opening the end-of-chapter PR (use chapter-pr).
---

# chapter-scaffold

Sets up the workspace at the *start* of a new Capital chapter. The actual
domain implementation is chapter-specific - that's the user's job - but the
boilerplate around it is the same every time, and getting it wrong (branch
name, chapter HTML slug, roadmap row format) creates rework.

## When the user invokes this

They'll say something like:

- "Let's start chapter 2"
- "Kick off chapter-04 - money into capital"
- "Set up the next chapter, exchange and money"

If the chapter number and slug aren't both given, ask once. Slug is
kebab-case and matches Marx's section heading without articles -
`02-the-process-of-exchange`, `04-the-general-formula-for-capital`. Look at
existing entries in `docs/architecture.md` (the roadmap table) before
inventing one - the chapter range is already mapped there.

## Steps

Do these in order. Don't batch the git operations - if branch creation
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

```bash
git pull --ff-only origin main
git checkout -b chapter-NN-<slug>
```

If you hit a `.git/index.lock` permissions error (sandbox limit noted in
CLAUDE.md), `mv` the lock file out of the way:

```bash
mv .git/index.lock .git/index.lock.bak
```

then retry the git command.

### 3. Create the chapter HTML placeholder

The user will paste the real ~100KB Marx text in later. For now write a
small stub at `chapters/NN-<slug>.html` so the path exists and the PR
template's "Chapter HTML" field has a target:

```html
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <title>Capital Vol. I, Ch. NN - <Title></title>
  </head>
  <body>
    <h1>Capital Vol. I, Ch. NN - <Title></h1>
    <p><em>Chapter source text to be pasted here by the user.</em></p>
  </body>
</html>
```

Replace `NN` and `<Title>` with the chapter number and Marx's chapter
heading (e.g. "The Process of Exchange").

### 4. Update the roadmap in docs/architecture.md

Find the row for this chapter in the "Roadmap (chapter-driven)" table and
flip its Status from `Pending` (or `Next`) to `In progress`. Don't add a
new row unless the chapter genuinely isn't in the roadmap - the table is
canonical.

For chapters that touch a brand-new service for the first time, also note
this in the Status column ("In progress · scaffolds X-service").

### 5. Stub a test file with a Marx textual fixture

The repo convention is that every chapter's tests use Marx's own examples
as fixtures: `20 yards linen = 1 coat`, `1 quarter corn = x cwt iron`,
etc. This grounds the simulation in the text and gives the test names a
reading-list quality. Drop a `*_test.go` file in the relevant
`internal/<domain>/` package that imports the package and has a single
skipped test naming the chapter:

```go
package <domain>

import "testing"

// TestChapterNN_<slug> is the textual anchor for Capital Vol. I, Ch. NN.
// Replace t.Skip with the chapter's first textual example as a fixture
// (e.g. "20 yards of linen = 1 coat" for Ch. 1 §3).
func TestChapterNN_<slug>(t *testing.T) {
	t.Parallel()
	t.Skip("TODO: implement chapter NN")
}
```

Skipped tests serve as a TODO list visible to `go test ./...` rather than
buried in a comment.

### 6. Don't commit

Leave the working tree dirty. The user wants to see the diff before any
commit, and signed commits fail in the sandbox anyway. Tell them what was
changed and what to do next:

> Branch `chapter-NN-<slug>` is set up. Roadmap row flipped to In progress,
> chapter HTML stub at `chapters/NN-<slug>.html`, test stub in
> `services/<svc>/internal/<domain>/`. Next: implement the domain types.
> Run `make vet test build` before committing - I can't run the Go toolchain
> here.

## What this skill does NOT do

- **Scaffold a new service.** If the chapter needs a brand-new service
  (e.g. agent-service for Ch. 4), call out that service-scaffolding is a
  separate concern and ask whether to proceed manually.
- **Generate domain types.** What goes in `commodity.go`, `value.go`, etc.
  is the heart of the chapter and is judgment-laden - it's Marx's
  argument turned into types. Don't pre-populate it.
- **Open a PR.** That's `chapter-pr`'s job, run at the end of the chapter.

## Anti-patterns

- Don't create per-service `go.mod` files - the repo is intentionally a
  single module.
- Don't add a TypeScript router or new web routes here - one `App.tsx`
  composes panels until a chapter forces otherwise (CLAUDE.md).
- Don't generate fake/lorem text for the chapter HTML. Leave the stub
  obviously unfilled so the user remembers to paste the real text.
