---
name: chapter-spec
description: Generate a chapters/volume-1/NN-<slug>.spec.md file from a chapter HTML — strips the raw Marx text to clean prose, then spawns a dedicated subagent to analyze concepts, extract test fixtures, derive invariants, and map scope to Go types and HTTP routes. Use this whenever the user asks to generate a chapter spec, convert a chapter HTML, "prep" a chapter for implementation, or says something like "create the spec for chapter N". Works both retrospectively (chapter already implemented — maps existing types) and prospectively (chapter not yet implemented — proposes what to build).
---

# chapter-spec

Chapter HTML files are 100–200KB of raw Marx prose. The spec is a
~100-line companion file that distils what the chapter means for code:
concepts → types, Marx's own numbers → test fixtures, economic laws →
invariants, and section scope. Generating it well requires reading the
full chapter, which is too large for the main context — so this skill
delegates to a subagent.

## Inputs to gather

Before starting, confirm:

1. **Chapter number and slug.** Infer from the branch name
   (`chapter-NN-<slug>`) or the chapters/volume-1/ directory listing. If ambiguous,
   ask once.
2. **HTML path.** `chapters/volume-1/NN-<slug>.html`. Must exist.
3. **Implementation status.** Is the chapter already implemented (types
   exist in `services/*/internal/`)? If yes, the spec will map them. If
   no, the spec will propose them. Ask if not obvious from context.

## Steps

### 1. Strip the HTML

Run the strip script to convert the chapter HTML to clean prose:

```bash
python3 scripts/strip_chapter_html.py chapters/volume-1/NN-<slug>.html
```

Capture the output. It will be 15–30KB — small enough to pass to a
subagent. Do not read the raw HTML directly.

### 2. Collect codebase context

Read the following (small, fast) — these go into the subagent prompt:

- `services/<svc>/internal/<domain>/` — all `.go` files for the chapter's
  primary service (to list existing or analogous types)
- `services/<svc>/internal/transport/httpapi/routes.go` — existing routes
- `docs/architecture.md` — the roadmap row for this chapter

For a prospective spec (not yet implemented), read the equivalent files
from the closest analogous service (commodity-service is always the
reference implementation).

### 3. Spawn the spec subagent

Use the Agent tool (general-purpose) with this prompt, substituting the
actual stripped prose and codebase snippets inline:

---

**Subagent prompt template:**

```
You are generating a chapter spec for the capital-simulator project.
The project is a Go + React microservices simulation of Marx's Capital,
Vol. I, implemented chapter by chapter.

## Your task

Produce a `chapters/volume-1/NN-<slug>.spec.md` file in exactly the format
shown in the ## Output format section below. Do not add prose outside
that format. Return only the markdown content of the spec file.

## Chapter

Chapter NN: <Title>
Implementation status: <already implemented | not yet implemented>
Primary service: <svc>-service (port <port>)

## Existing Go types (for reference)

<paste the relevant .go file contents here>

## Existing HTTP routes

<paste routes.go content here>

## Architecture roadmap row

<paste the relevant row from docs/architecture.md>

## Stripped chapter prose

<paste the full output of strip_chapter_html.py here>

## Output format

Produce exactly this structure (fill in the brackets):

---
chapter: NN
title: "<Marx's chapter title>"
status: <implemented | proposed>
primary_service: <svc>-service
---

## Concepts → types

| Marx term | Go identifier | Kind | Package | Notes |
|---|---|---|---|---|
| [term] | [Type or FuncName] | type/func/method | [pkg] | [one-line note] |

Include every concept the chapter names that maps to something in code.
For a prospective spec, propose identifiers that follow repo conventions
(PascalCase types, camelCase methods, snake_case JSON tags).

## Fixtures

Marx's own numbers and examples from the text. These become test case
names and values verbatim — copy the wording from the text exactly.
Format:

- **§N** `<exact quote or equation from text>` → `<what it asserts in code>`

Include at least one fixture per section. Prefer equations and
quantitative comparisons over prose.

## Invariants

Mathematical or logical laws the chapter establishes that tests must
enforce. State each as a code-checkable assertion:

- `<expression> == <expression>` or `<condition>` [cite §N]

## Scope

### This chapter builds
- Services: [list]
- New domain types: [list with one-line descriptions]
- New HTTP endpoints: [METHOD /path — description]
- React: [what UI changes, if any]

### Explicitly deferred to later chapters
- [thing left out] — [why / which chapter picks it up]

---
End of subagent prompt template.
```

### 4. Write the spec file

Take the subagent's output and write it to
`chapters/volume-1/NN-<slug>.spec.md`.

Verify it:
- Has frontmatter (chapter, title, status, primary_service)
- Has all four sections (Concepts, Fixtures, Invariants, Scope)
- Fixtures cite section numbers (§1, §2, etc.)
- Invariants are checkable expressions, not prose descriptions
- Scope lists deferred items explicitly

If any section is thin (< 2 entries), run a quick pass yourself — the
subagent may have missed examples in dense prose sections.

### 5. Don't commit

Leave the spec uncommitted. The user should review it — the subagent
reads prose correctly but can mis-map Go types for prospective specs
where naming conventions are ambiguous. The spec is a starting point,
not a final artifact.

## What this skill does NOT do

- **Implement the chapter.** The spec is planning material. Chapter
  implementation is separate.
- **Replace the HTML.** The HTML stays as the canonical authoritative
  source. The spec is a code-facing view of it.
- **Run `make vet test build`.** That's the user's job after
  implementing from the spec.

## Anti-patterns

- Don't read the HTML directly — it will consume most of the context
  window and the strip script produces much cleaner output.
- Don't invent Go types not warranted by the chapter text. If the chapter
  doesn't introduce a concept, it shouldn't appear in the spec.
- Don't put implementation details (function bodies, algorithms) in the
  spec — just identifiers and signatures. The spec is a map, not the
  territory.
