# Ch. 13 Co-operation — design patch

Drop-in redesign of the Co-operation page for `theding0x/capital-simulator`.
Fixes the clunky 3-column checkbox roster (whose bottom-border inputs ran
across gutters and whose "sup" toggle read as a label, not a control) and
foregrounds the live composite working day — the central claim of Ch. 13 §3.

## What's here

```
web/src/chapters/Ch13Cooperation.tsx   ← replaces existing file
web/src/chapters/Ch13Cooperation.css   ← new file, imported by the .tsx
```

## What changed (vs. the existing Ch13Cooperation.tsx)

- **`MinimumCapitalPanel`** — unchanged behavior, copy-paste preserved.
- **`CooperationLedgerPanel`** — same data wiring (lists capitalists, workers,
  cooperations; refresh tick on create); presentation rewritten.
- **`CreateCooperationPanel`** — the substantive redesign:
  - Command row (name / capitalist / day per worker) uses a 1.4fr 1fr 0.8fr
    grid with single-line labels and `align-items: end`, so the three
    inputs share a baseline regardless of label length. (Original labels
    like "Capitalist (single command)" wrapped and dropped the Name input
    out of alignment.)
  - The 3-column checkbox+"sup" worker grid is replaced by a `.data-table`
    roster. Each row: selection dot, name (Playfair), role label, and an
    OVERSEER chip that appears only on selected rows — so "supervisory"
    is now a state, not a stray label.
  - "Select all" / "Clear" affordances.
  - A live **CompositeSummary** between the form and the roster shows
    labourers × per-worker day = composite labour-day, with the §3
    productive-power gloss when ≥ 2 workers are selected. This was
    previously buried in the post-assembly `Compute` step.
  - "Assemble" is disabled until ≥ 2 workers and a capitalist are picked,
    with a sentence-form explanation in the submit row.
- **`CooperationsList`** — table column set is the same; rows are now
  click-to-expand and the computed result renders inline as a
  `.reveal-panel` (matching Ch. 02's pattern), with `RelativeSurplusBridge`
  preserved inside the expanded row.

## How to apply

Copy both files over the equivalents in your repo:

```sh
# from repo root
cp /path/to/Ch13Cooperation.tsx web/src/chapters/Ch13Cooperation.tsx
cp /path/to/Ch13Cooperation.css web/src/chapters/Ch13Cooperation.css
```

No other files need to change. The CSS imports as a side-effect from the
.tsx (`import "./Ch13Cooperation.css";`) — make sure your bundler is set
up for CSS-in-JS imports (Vite is, by default).

## Behavior preserved

- All API calls go through the same `api.*` methods.
- `members[]` are still constructed as
  `{ worker_id, supervisory, working_day_minutes }` per the
  `CreateCooperationInput` contract.
- `RelativeSurplusBridge` is still invoked with `source: "cooperation"`,
  `sourceId`, and the factor from `computeCollectiveWorkingDay`.

## Lifted markup

The redesign relies on these existing classes from `web/src/index.css`,
unchanged:

- `.card`, `.card > h2`, `.card > .description`
- `.form-grid`, `.form-grid label > span`, `.form-actions`, `.span2`
- `.data-table` and its descendants
- `.reveal-panel`, `.reveal-note`, `.labour-statement`
- `.empty-state`, `.muted`, `.small`, `.error`

All new selectors are scoped `.ch13-*`.
