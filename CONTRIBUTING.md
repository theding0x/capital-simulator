# Contributing

Capital Simulator turns each chapter of Marx's *Capital, Volume I* into
running code: domain types, HTTP endpoints, MySQL tables, and React
panels, with Marx's own numerical examples preserved as test fixtures
and seed rows. The implementation is AI-assisted and reviewed by one
person. That review can catch software defects — type errors, broken
tests, malformed JSON — but it has trouble catching the more
interesting kind of defect: a category that quietly drifts from what
Marx actually wrote.

**That is where you come in.** If you have read *Capital* carefully —
in any edition, in any language, with any interpretive tradition — and
you spot a place where this code misrepresents the argument, the most
useful thing you can do is open an issue.

You do not need to write Go. You do not need to write TypeScript. You
do not need to run the project. You can contribute meaningfully with
nothing but a copy of *Capital* and a GitHub account.

---

## A note before you read further: chapter numbering

This project uses the **Moore–Aveling English chapter numbering**, the
one digitised and served by the [Marxists Internet
Archive](https://www.marxists.org/archive/marx/works/1867-c1/) (MIA)
and reproduced in most older English editions. The **Penguin / Fowkes**
translation renumbers some chapters and combines others. The
differences begin around Part III and become significant later.

If a chapter number in your edition doesn't match ours, please cite
either the chapter **title** or a marxists.org URL so we can be sure
we're talking about the same passage. A bare "Ch. 17" without a title
is the most common source of cross-edition confusion.

---

## Three ways to contribute

### 1. Interpretation issues (the kind most needed)

If you read a chapter's implementation summary in
[`docs/architecture.md`](docs/architecture.md), or play with a
dashboard panel, or read the test fixtures in
`services/<svc>/internal/<pkg>/*_test.go`, and you find:

- a Marx concept the code names but mis-defines,
- a Marx concept the code defines correctly but tests against the
  wrong textual example,
- a Marx concept the code conflates with another concept (e.g.
  treating use-value as if it were utility, or treating constant
  capital as if it included variable),
- an invariant the code enforces that Marx does *not* claim (an
  artefact of trying to make the model executable),
- an invariant Marx *does* claim that the code silently breaks,
- a translation of Marx's English-edition prose into a Go type that
  loses something the German original keeps clearer,

— open an issue. A good interpretation issue looks roughly like this:

> **Title:** Ch. 8 — `ConstantCapital.WearFractionFor` treats
> auxiliary materials as fully consumed in one cycle, but Marx (Ch. 8
> §1, marxists.org/.../ch08.htm) classes some auxiliaries (lubricating
> oil, fuel for slow-burning furnaces) as transferring value across
> multiple cycles.
>
> **Where in the code:**
> `services/commodity-service/internal/commodity/constant_capital.go`,
> the `auxiliary` branch of `WearFractionFor`.
>
> **What Marx writes:** [one-paragraph paraphrase or short quote with
> citation].
>
> **What the code does:** [one-paragraph description of the current
> behaviour].
>
> **Suggested correction:** [optional — if you have one in mind. If you
> don't, that's fine; the issue itself is the contribution].

The maintainer will read the issue, look at the passage, look at the
code, and either: (a) accept it and open a fix branch, (b) push back
with a counter-reading, or (c) tag it as a deliberate interpretive
choice with a comment in `docs/architecture.md` explaining the
trade-off.

**You are not expected to be diplomatic about errors.** A blunt "this
is wrong because Marx says X" is more useful than a hedged "I wonder
if perhaps it might be worth considering whether X."

### 2. Fixture and numerical corrections

The project preserves Marx's worked numerical examples as test fixtures
and as seed rows in MySQL. Examples already in the tree:

- "20 yards linen = 1 coat" (Ch. 1 §3)
- "£10,000 capital = £8,000c + £2,000v, s′ = 100%" (Ch. 24 §1, the spinner)
- Redgrave's 1866 spindle ratios across seven European countries (Ch. 22)
- Burke's five-man platoon (Ch. 13 §3)
- The Caslon type-foundry 4/2/1 ratio (Ch. 14 §2)
- Thomas Hobson's piece-wage at 6 farthings / 24 pieces (Ch. 21)

These numbers are easy to get wrong: transposed digits, wrong currency
unit, rounding from a translator's footnote rather than Marx's
original. If a fixture in `*_test.go` or a seed row in
`internal/store/migrations/NNNNN_chNN_seed.sql` does not match what
Marx actually wrote in your edition, open an issue:

> **Title:** Ch. 22 seed — Belgium intensity factor is 0.89, but Marx
> (Ch. 22 §IV table, marxists.org/.../ch22.htm) gives 0.82.
>
> **Where:** `services/agent-service/internal/store/migrations/00021_ch22_seed.sql`,
> line 5.
>
> **Source:** [citation, ideally a marxists.org URL or page reference].

Numerical corrections are usually unambiguous, get fast turnaround,
and ship as a new append-only migration that supersedes the seed row.
They are extremely welcome.

### 3. Code contributions

If you want to write code — fix a bug, add a chapter the roadmap marks
as Pending, extend a panel, or improve the test suite — the full
contributor playbook is in [`CLAUDE.md`](CLAUDE.md). The short form:

- Branch off `main` as `volume-X/chapter-Y` (for chapter work) or any
  descriptive branch name (for bug fixes and refactors).
- Run `make vet test build` before pushing. The CI run is the same.
- For frontend changes: `cd web && npm run lint && npm run build`.
- Sign your commits if you can (`commit.gpgsign=true`). If you can't,
  unsigned PRs are still accepted.
- Open the PR with a description of *what you changed* and *why*. For
  chapter work the body should reference the chapter section and the
  Marx example you used as a test fixture.
- Migrations are append-only. Never edit an existing `.sql` file — add
  a new numbered file.
- Domain types use concrete Go types only. No `interface{}` / `any` in
  struct fields.

The maintainer reviews every PR personally. Expect comments. Expect
questions about whether your design captures Marx's argument or just a
software engineer's version of it.

---

## How disagreements are resolved

Marx's argument is contested. Two scholars can read the same paragraph
and reach genuinely different formal models of it — and both can be
defensible readings of the text. The project does not try to be
neutral; it tries to be one *coherent* reading, with the trade-offs
documented.

When an interpretation issue gets pushback rather than acceptance, the
options are:

- **Maintainer concedes.** Issue accepted, fix ships,
  `docs/architecture.md` updated.
- **Contributor concedes.** Issue closed with a comment that records
  the alternative reading for future readers searching the issue
  tracker.
- **Neither concedes, both readings are defensible.** Issue stays
  open and gets labelled `interpretive-choice`. A note is added to
  `docs/architecture.md` explaining which reading the code uses and
  which it doesn't. The issue becomes a permanent footnote rather
  than a bug.

The third path is the most common for genuinely contested questions,
and it is fine. The project is more useful when its choices are
visible than when they pretend to be obvious.

---

## Citation conventions

The repo's source-of-truth text is the Moore–Aveling translation of
*Capital, Vol. I* (1887), as digitised, transcribed, and hosted by the
**Marxists Internet Archive**:

`https://www.marxists.org/archive/marx/works/1867-c1/`

Every chapter implementation, test fixture, and seed row in this
repository ultimately traces back to that URL. Crediting MIA when you
quote *Capital* is not a legal requirement — Marx's text and the
Moore–Aveling translation are both public domain — but it is the right
thing to do. The archive's transcription work is what made this
project possible.

When citing Marx, please include at least one of:

1. **A marxists.org URL** (preferred — unambiguous, no edition drift).
2. **A chapter + section heading** (e.g. "Ch. 15 §3c, Intensification
   of Labour").
3. **An edition + page number** (e.g. "Penguin/Fowkes p. 534" or
   "Moscow/Progress p. 422" or "MEW 23, S. 446").

The first is best; the second is fine on its own; the third is fine
combined with one of the others. A bare page number from an unnamed
edition is hard to verify.

For secondary literature (Rubin, Heinrich, Postone, Harvey, Bidet,
Banaji, and others) standard academic citation is welcome but never
required. Your reading of Marx is enough.

---

## Code of conduct

Three rules, applied to issue threads and PR reviews:

1. **No condescension about programming skill.** If a Marx scholar
   opens an issue with a screenshot of a dashboard panel and the
   sentence "this is wrong," that issue is a contribution, not a
   support ticket. Treat it as such.
2. **No condescension about Marx scholarship.** If a contributor
   makes an interpretive claim the maintainer disagrees with, the
   reply engages the claim. "Read Heinrich" or "that's not what Marx
   means" without elaboration are not replies.
3. **No Marx wars in the issue tracker.** Disagreements about Marx's
   politics, his legacy, the merits of value-form theory vs. the
   traditional reading, the relationship of Marx to Engels, etc. are
   genuinely interesting and genuinely off-topic. Take them
   elsewhere; this tracker is for the simulator.

Beyond those three: the
[Contributor Covenant](https://www.contributor-covenant.org/version/2/1/code_of_conduct/)
v2.1 is the default. Personal attacks, harassment, and discrimination
are not tolerated regardless of how much Marx the attacker has read.

---

## What happens after you open an issue

The maintainer will typically:

1. Acknowledge the issue within a few days.
2. Read the cited passage and the implicated code.
3. Either open a fix branch, push back with a counter-reading, or tag
   the issue as an interpretive choice.

Issues are not closed silently. If an issue stalls, ping it — the
maintainer is one person and the queue grows.

---

## Thank you

A formal model of Marx's argument is only useful to the extent that
the model corresponds to the argument. The corresponding part is the
part this maintainer cannot reliably verify alone. Every interpretation
issue, every fixture correction, every blunt "this is wrong" makes the
project more useful than it was before.

And a separate, prior thank-you to the volunteers of the **Marxists
Internet Archive** for digitising, transcribing, proof-reading,
indexing, and hosting *Capital* (along with the rest of the Marx /
Engels corpus, and a great deal else) for thirty-plus years and
counting. If this project is useful to you, please consider
[donating to MIA](https://www.marxists.org/admin/intro/general/donate.htm).
