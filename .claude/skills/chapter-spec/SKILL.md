---
name: chapter-spec
description: Load the chapter spec for a given chapter from the red-vault Obsidian vault. Specs are fully authored ahead of time and live at `marx-engels/1867/capital-volume-i/specs/NN-<slug>.spec.md` in the vault. Use this whenever the user asks to read, fetch, or inspect a chapter spec, or says "load the spec for chapter N", "what's in the spec for chapter N", or "prep me to implement chapter N".
---

# chapter-spec

Chapter specs (concepts → Go types, Marx's fixtures, invariants, scope)
are authored and committed in the red-vault Obsidian vault, not in this
repo. This skill is a thin wrapper around `mcp__obsidian__obsidian_get_file_contents`
that names where to look and what to do with the result.

## Inputs to gather

1. **Chapter number.** Infer from the current branch (`volume-X/chapter-Y`)
   if you're already on a chapter branch. Otherwise ask once.
2. **Volume.** Defaults to 1. We only have Vol. I specs right now.

## Steps

### 1. Resolve the slug

The vault file is named with the chapter number + slug, e.g.
`17-changes-of-magnitude-in-the-price-of-labour-power-and-in-surplus-value.spec.md`.
List the specs directory once to pick the matching file:

```
mcp__obsidian__obsidian_list_files_in_dir
  dirpath: marx-engels/1867/capital-volume-i/specs
```

Match by the `NN-` prefix.

### 2. Fetch the spec

```
mcp__obsidian__obsidian_get_file_contents
  filepath: marx-engels/1867/capital-volume-i/specs/NN-<slug>.spec.md
```

The spec is small (~100 lines) — fine to keep in context for the
implementation work that follows.

### 3. Present what matters

When handing the spec back to the user, lead with the four sections in
this order:

1. **Concepts → types** (the Go identifiers to introduce)
2. **Fixtures** (Marx's numbers for tests)
3. **Invariants** (laws the tests must enforce)
4. **Scope** (what this chapter builds vs. defers)

Don't reprint the spec wholesale — summarize the planned domain types
and endpoints, and quote fixtures verbatim (their wording becomes test
case names).

## What this skill does NOT do

- **Author or rewrite specs.** Specs live in the vault and are authored
  by the user out-of-band. If the user asks to edit a spec, write the
  change with `mcp__obsidian__obsidian_patch_content` against the same
  vault path — but treat that as an explicit content edit, not a
  generation step.
- **Read the chapter source text.** The Marx prose lives at
  `marx-engels/1867/capital-volume-i/texts/NN-<slug>.md`. Pull it only
  when the spec is genuinely insufficient for the implementation
  decision at hand — the spec is the code-facing view, the text is the
  audit trail.
- **Implement the chapter.** Use the spec as planning material; the
  implementation work is a separate task.

## Anti-patterns

- Don't look for `chapters/volume-1/*.spec.md` in the repo — those were
  migrated out. The vault is the source of truth.
- Don't generate a spec from the chapter text. Specs are pre-authored.
  If the vault is missing a spec the user expects, surface that as a gap
  rather than fabricating one.
