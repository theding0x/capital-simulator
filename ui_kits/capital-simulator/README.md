# Capital Simulator — UI kit

Pixel-faithful recreation of the dashboard at `web/src/` from theding0x/capital-simulator.

## Files

- `index.html` — interactive demo. Two-column shell (sidebar + chapter body) with a working Ch. 01 "Commodity" panel: register a new commodity, see it appear in the table, click **Reveal** to open the gold panel, click **Compute ratio** to derive an exchange equation. Sidebar lets you switch chapters; pending chapters show the placeholder.
- `Topbar.jsx` — sticky 42px topbar with wordmark + chapter breadcrumb.
- `Sidebar.jsx` — 188px chapter list grouped by Part, with done / active / pending row states.
- `ChapterHeader.jsx` — eyebrow + Playfair title + italic Marx quote.
- `Card.jsx` — the system's distinctive band-card (no fill, no shadow, hairline rule).
- `CommodityForm.jsx` — register-a-commodity form (Ch. 01).
- `CommodityTable.jsx` — commodity table with hover + Reveal/Edit/Delete actions.
- `RevealPanel.jsx` — gold-tinted social-relations panel.
- `ExchangeRatio.jsx` — exchange-ratio computer + equation result block.
- `chapters.js` — chapter registry (mirrors `web/src/chapters/registry.ts`).

## Coverage notes

- **Ch. 01 (The Commodity)** is fully wired in the demo.
- **Ch. 02–04** are styled with the same components but stubbed out — clicking them shows the chapter header and a "Coming up" note rather than the live forms (the kit's purpose is component coverage, not full chapter recreation).
- **Pending chapters (05–33)** show the same dashed placeholder the live product uses.

## What's missing

Per CLAUDE.md the live product only ships `done` chapters 01–04. The kit covers all the **components** those chapters use; it does not recreate Ch. 02's market-offer flow, Ch. 03's circuit-of-money panel, or Ch. 04's M-C-M-prime agent ledger in working form. Those would reuse `<Card>`, `<CommodityForm>`-style grids, and `<CommodityTable>`-style tables — the visual vocabulary is already covered.
