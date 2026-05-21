---
name: chapter-spec
description: Load the chapter spec for a given chapter from the red-vault Obsidian vault. Specs live at `marx-engels/<year>/capital-volume-<roman>/specs/NN-<slug>.spec.md` — Vol. I (1867), Vol. II (1885), Vol. III (1894). Use this whenever the user asks to read, fetch, or inspect a chapter spec, or says "load the spec for Vol. X chapter N", "what's in the spec for chapter N", or "prep me to implement Vol. II Ch. 5".
---

# chapter-spec

Chapter specs (concepts → Go types, Marx's fixtures, invariants, scope)
are authored and committed in the red-vault Obsidian vault, not in this
repo. This skill is a thin wrapper around `mcp__obsidian__obsidian_get_file_contents`
that names where to look and what to do with the result.

## Volume → vault path

Capital has three volumes, each at its own vault path keyed by
publication year:

| Volume | Year | Vault prefix                                |
|--------|------|---------------------------------------------|
| I      | 1867 | `marx-engels/1867/capital-volume-i/`        |
| II     | 1885 | `marx-engels/1885/capital-volume-ii/`       |
| III    | 1894 | `marx-engels/1894/capital-volume-iii/`      |

The spec lives at `<prefix>specs/NN-<slug>.spec.md`; the source text
lives at `<prefix>texts/NN-<slug>.md`.

## Inputs to gather

1. **Volume** (1, 2, or 3). Required, but inferable:
   - From the current branch when it matches `volume-X/chapter-Y` (the
     canonical chapter-branch convention).
   - From the user's phrasing ("Vol. II Ch. 5"). If they just say
     "chapter 5", ask which volume — Ch. 5 exists in all three.
2. **Chapter number.** Infer from the branch (`chapter-Y` segment) or ask.

## Steps

### 1. Resolve the vault path

Map the volume number to the path prefix from the table above. The two
working paths are:

```
<prefix>specs/      # NN-<slug>.spec.md      — code-facing spec (this skill)
<prefix>texts/      # NN-<slug>.md           — Marx's prose (audit trail)
```

### 2. Find the slug

The vault filename is the chapter number + slug, e.g.
`17-changes-of-magnitude-in-the-price-of-labour-power-and-in-surplus-value.spec.md`
(Vol. I) or `01-the-circuit-of-money-capital.spec.md` (Vol. II). List
the specs directory once and match by the `NN-` prefix:

```
mcp__obsidian__obsidian_list_files_in_dir
  dirpath: <prefix>specs
```

### 3. Fetch the spec

```
mcp__obsidian__obsidian_get_file_contents
  filepath: <prefix>specs/NN-<slug>.spec.md
```

The spec is small (~100 lines) — fine to keep in context for the
implementation work that follows.

### 4. Present what matters

When handing the spec back to the user, lead with the four sections in
this order:

1. **Concepts → types** (the Go identifiers to introduce)
2. **Fixtures** (Marx's numbers for tests)
3. **Invariants** (laws the tests must enforce)
4. **Scope** (what this chapter builds vs. defers)

Don't reprint the spec wholesale — summarize the planned domain types
and endpoints, and quote fixtures verbatim (their wording becomes test
case names).

## Spec coverage today

- **Vol. I**: all 33 specs authored at `marx-engels/1867/capital-volume-i/specs/`.
- **Vol. II**: specs **not yet authored** (texts exist; the spec sweep
  is Phase 5 of the multi-volume foundation work). If asked for a
  Vol. II spec, list the texts directory and surface the gap rather
  than fabricating one.
- **Vol. III**: specs **not yet authored** (same as Vol. II).

## What this skill does NOT do

- **Author or rewrite specs.** Specs live in the vault and are authored
  by the user out-of-band. If the user asks to edit a spec, write the
  change with `mcp__obsidian__obsidian_patch_content` against the same
  vault path — but treat that as an explicit content edit, not a
  generation step.
- **Read the chapter source text.** The Marx prose lives at
  `<prefix>texts/NN-<slug>.md`. Pull it only when the spec is genuinely
  insufficient for the implementation decision at hand — the spec is
  the code-facing view, the text is the audit trail.
- **Implement the chapter.** Use the spec as planning material; the
  implementation work is a separate task.

## Anti-patterns

- Don't look for `chapters/volume-1/*.spec.md` in the repo — those were
  migrated out. The vault is the source of truth.
- Don't generate a spec from the chapter text. Specs are pre-authored.
  If the vault is missing a spec the user expects, surface that as a gap
  rather than fabricating one.
- Don't default to Vol. I when the branch or the user's phrasing names
  a different volume.
