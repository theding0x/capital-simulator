# Atlas — Design Handoff Brief (for Claude Design)

- **Date:** 2026-06-05
- **For:** a design-focused session/pass to reimagine the Atlas landing experience.
- **Status:** functionally complete (3 slices shipped, faithful, live), **visually underwhelming** — the mandate here is a real design pass.
- **Branch:** `feature/atlas-observatory` (unmerged). All Atlas code is under `web/src/atlas/`.
- **Mandate (all four, per the client):** ① aesthetic & atmosphere, ② motion & the descent, ③ information design, ④ **rethink the concept** — not just restyle. You have latitude to reimagine the experience, within the guardrails in §7.

> **Boot it live before designing.** Static screenshots miss the point (orbit motion, the descent, the law responding to a lever). `docker compose up --build -d`, open `http://localhost:5173/`, then `POST http://localhost:8080/v1/engine/start`. Reference stills: `atlas-abode.png`, `atlas-levers.png` (repo root).

---

## 1. What Atlas is (and must remain)

Capital Simulator models an economy **chapter-by-chapter on Marx's *Capital***. Atlas is its landing page: not a control panel, but **the whole circuit of capital — `M—C(Lp+Mp)…P…C′—M′` — in motion**, at the scale of many capitals. Three layers, already built:

1. **The surface — a field of capital-orbits.** Each capital is a ring of three coexisting arcs (M money / P production / C′ commodity), with value-dots travelling the circumference (one lap = one turnover, lingering in production). The field is centred on the **average rate of profit `p̄′`** as a centre of gravity. Capitals **spiral outward** as they accumulate.
2. **The hidden abode** — cross a threshold ("No admittance except on business," Vol. I Ch. 6) from the value-form *surface* down into *production*, where the class relation is laid bare: the **social working day** (necessary vs surplus labour = the rate of exploitation s/v), **living labour** (Σv), the **industrial reserve army** + wage pressure, **surplus extraction** (Σs), and an **immiseration trend** (the General Law of Accumulation, Ch. 25, running in real time).
3. **The levers** — working day, wage, accumulation rate α — perturb the live law and watch it respond.

**The point:** the shipped chapters define *a system in motion*, and Atlas should make you *feel* the circuit, the antagonism, and the law of motion — alienation, accumulation as dead labour towering over living. It's an "Observatory."

**Non-negotiable fidelity (semantic invariants — do not break these in any redesign):**
- The circuit order `M → P → C′ → M′` and the three coexisting arcs (Vol. II *Nebeneinander*).
- s/v is the working day's *necessary | surplus* split; surplus labour is *unpaid* labour.
- Rising organic composition (c/v) → growing reserve army → wage pressed below the value of labour-power → rising s/v → more accumulation. The loop must read as a **loop**.
- Gold = surplus/value, red = the brand/capital, lead-blue = necessary/constant. (Colour carries meaning — see §6.)
- Money is integer pence; rates are basis points; labour is minutes. No fake data — design against the real contract (§5).

---

## 2. Current state (what's built)

| Layer | Components (`web/src/atlas/`) | Notes |
|---|---|---|
| Shell | `Atlas.tsx` | topbar + left rail (`VitalSigns`) + body + bottom transport (`TickHeartbeat`); `descended` toggles surface↔abode |
| Surface | `CircuitField.tsx`, `Orbit.tsx`, `animation.ts` | SVG rings; rAF dot pacing (`pacedAngle`, lingers in P); CSS-scale growth vs a load-time reference |
| Abode | `Abode.tsx`, `GeneralLawTrend.tsx` | card stack: working-day bar, 2×2 stat grid, a small dual-line sparkline |
| Levers | `Levers.tsx` | three `<input type=range>` sliders, POST on change |
| Data | `useSnapshot.ts` | polls `GET /v1/observatory/snapshot` every 2s, holds last-good on error |
| Styles | `atlas.css` + `web/src/index.css` (`:root` tokens) | |

Design spec + build plans for full context: `docs/superpowers/specs/2026-06-04-atlas-general-law-design.md`, `docs/superpowers/plans/2026-06-04-atlas-general-law-slice{1,2}.md`, `docs/superpowers/plans/2026-06-05-atlas-general-law-slice3.md`.

---

## 3. Honest critique — why it underwhelms

**Aesthetic & atmosphere.** It reads as a *generic dark dashboard*, not an observatory. Flat panels on a near-black ground; almost no depth, light, texture, or focal drama. **It doesn't even use its own brand system** — the abode/levers CSS hardcodes hex (`#141821`, `#c8a240`, `#4a5a8a`) instead of the `:root` tokens, and uses **none** of the brand's display serif (Playfair Display) or mono figures (IBM Plex Mono with tabular-nums). The result is generic where the brand is distinctive.

**Motion & the descent.** The orbits are static rings with three small dots; nothing conveys *value in motion*, accumulation as a spiral, or the menace of dead labour. Growth is a quiet CSS `scale()`. The surface→abode "descent" — the conceptual heart — is a 0.5s `translateY` fade; it should feel like *crossing a threshold into a hidden place*. When a lever moves the law, the response is a number changing, not a felt consequence.

**Information design.** The abode is a **stack of equal-weight cards** with no hierarchy — the social working day (the hero) sits at the same visual altitude as a stat tile. The immiseration trend (the entire point of "the General Law in motion") is a tiny afterthought sparkline. Σ-values, rates, counts, and minutes all compete; nothing guides the eye to the *story* (exploitation widening, the reserve army swelling, the wage falling). The left rail's `VitalSigns` is a plain key/value list.

**Concept.** The surface-field + descend-into-abode framing is sound and faithful, but the execution makes it feel like two disconnected screens rather than one continuous world. The link between the abode (where surplus is pumped from labour) and the surface (where it returns as towering capital) — the alienation loop — is asserted in copy, never *shown*. There's room to reimagine this as a single, vertically-continuous space, or a richer spatial metaphor, rather than a toggle.

---

## 4. Design direction (north star — provocations, not specs)

Aim: someone lands on this and feels they're watching *a living economy and its law of motion*, beautiful and a little ominous. Treat the four axes together.

**① Aesthetic & atmosphere**
- Make it an *observatory / orrery*: depth, a luminous centre of gravity (`p̄′`), capitals as celestial bodies whose light/heat encodes their state. Use the brand's Playfair Display for gravitas and Plex Mono (tabular-nums) for all figures — the "ledger meets cosmos" tension is the brand.
- Materiality: consider a sense of light from the centre, vignette, fine grain/paper texture under the data, gold leaf for surplus. Earn the near-black ground with contrast and glow rather than flatness.

**② Motion & the descent**
- Accumulation should *spiral*, legibly — the orbit grows and the dots' lap-rate reads as turnover. Higher-turnover capitals visibly churn faster.
- The **descent** is the signature moment: reimagine it as physically crossing a threshold — the surface parting/receding, the camera (or the world) sinking into production, the "No admittance" notice as a real gate. It should be ~1s of intent, not a fade.
- The **alienation loop**: when the abode pumps surplus (Σs), show it *rising to the surface as capital* (gold flowing up, orbits thickening). When a lever moves, the law's response should ripple visibly over the next passes.
- Respect reduced-motion preferences.

**③ Information design**
- Give the abode a clear hero (the **social working day** as s/v) and a deliberate hierarchy beneath it; demote stat tiles. Make the **immiseration trend** a real, prominent chart — the reserve army swelling, the wage falling, s/v widening over time is the *narrative*, not a footnote.
- Encode quantity in form, not just number: the working-day bar, the reserve-army "reservoir," composition as the weight of dead vs living labour. Numbers in Plex Mono, tabular, with units.
- The left rail (`VitalSigns`) and bottom transport (`TickHeartbeat` ECG) should feel like instruments on one console.

**④ Concept (you may rethink)**
- Consider collapsing the surface/abode *toggle* into one continuous vertical world you descend through (circulation above, production below, the reserve army beneath that), so the loop is spatial and always legible.
- Or push the orrery metaphor further (the field as a system with real gravitation toward `p̄′`).
- Keep the fidelity invariants (§1) intact whatever the framing.

---

## 5. The live data contract (design against this — it's all real)

`GET /v1/observatory/snapshot` (polled every 2s via `useSnapshot.ts`) returns:

```ts
interface ObservatorySnapshot {
  tick: number; running: boolean; interval_ms: number;
  capitals: FieldCapital[];      // the surface field
  aggregate: AggregateVitals;    // the rail's vital signs
  abode: AbodeReadout;           // the hidden abode + immiseration series
}

interface FieldCapital {
  id: string; status: string;        // status "active" | "halted"
  total_pence: number;               // grows each pass (the spiral)
  money_pence: number; production_pence: number; commodity_pence: number; // the M/P/C arcs
  cost_price_pence: number;          // c+v
  surplus_pence: number;             // s
  turnover_number: number;           // laps/period → dot lap-rate
}

interface AggregateVitals {
  total_social_capital_pence: number; cost_price_pence: number;
  surplus_pence: number; avg_rate_of_profit_bp: number; // p̄′ in basis points
}

interface AbodeReadout {
  total_variable_pence: number;        // Σv (paid labour)
  total_surplus_pence: number;         // Σs (unpaid labour)
  rate_of_exploitation_bp: number;     // s/v (effective)
  necessary_labour_minutes: number;    // ┐ the social working day (sum = 600 = 10h)
  surplus_labour_minutes: number;      // ┘
  organic_composition_bp: number;      // c/v
  reserve_army_count: number;
  reserve_army_pressure_bp: number;
  employed_count: number;
  wage_pence: number;                  // paid wage (compressed below value)
  surplus_rate_base_bp: number;        // ┐ live LEVER positions
  base_wage_pence: number;             // │ (slider init values)
  accumulation_rate_bp: number;        // ┘
  law_series: GeneralLawTrendPoint[];  // immiseration time-series (ascending, up to 60 pts)
}

interface GeneralLawTrendPoint {
  period: number; wage_pence: number; rate_of_exploitation_bp: number;
  reserve_army_count: number; organic_composition_bp: number;
}
```

**Controls (already wired):**
- `POST /v1/engine/start` · `/stop` — run/pause the law. `GET /v1/engine/ticks` feeds the bottom ECG.
- `POST /v1/observatory/levers` with any subset of `{ surplus_rate_base_bp, base_wage_pence, accumulation_rate_bp }` — perturb the live law; effects appear in the next snapshots. (`api.setObservatoryLevers` in `web/src/api.ts`.)

**Live dynamics to design around (what actually happens over ~30s with the engine running):** capitals' `total_pence` climbs (orbits spiral out); s/v rises from ~120% toward several hundred %; the reserve army swells; c/v climbs steeply (machinery displacing labour); the paid wage falls to its floor; `law_series` lengthens. The economy *immiserates* — the design should let you watch it happen.

---

## 6. Brand system (fixed — compose with it, don't replace it)

Tokens live in `web/src/index.css` `:root`. Fonts via Google Fonts (already imported): **Playfair Display** (display serif — gravitas), **IBM Plex Mono** (all figures, `font-variant-numeric: tabular-nums`), **IBM Plex Sans** (body).

```
grounds   --bg #07080a   --surface #0f1014   --surface-raised #15171d   --surface-hover #1a1d25
borders   --border #222530   --border-subtle #181b23   --rule rgba(255,255,255,.06)
ink       --ink #e8e2d5 (warm paper)   --ink-muted #8a8578   --ink-dim #3a3830
red       --red #c0392b  (capital / the brand)   --red-hover #d44030   --red-dim rgba(192,57,43,.14)
gold      --gold #9d7a2a   --gold-bright #c8a240  (surplus / value)   --gold-bg / --gold-border
lead      --lead #4a5a8a  (necessary labour / constant capital)   --lead-hover #5a6a9e
radius    --radius-sm 4   --radius 8   --radius-lg 12
```

**Voice:** sober, materialist, a touch literary (Marx's own register — "the hidden abode," "dead labour dominating living"). Use the full brand system for the full Atlas surface — the project's `capital-simulator-design` skill carries the complete palette/type/spacing/voice and a React UI kit; lean on it.

**Semantic colour (keep consistent):** gold → surplus/value/dead-capital; lead-blue → necessary labour/variable→living; red → the brand accent and danger/pressure. The working-day bar already uses lead (necessary) | gold (surplus) — that mapping is correct and worth keeping.

---

## 7. Constraints & guardrails (what must not break)

- **Stack (fixed):** React 18 + Vite + TypeScript, **plain CSS** (no Tailwind/CSS-in-JS), **no router** (one `Atlas.tsx`, hash links to `#/chapters`). SVG + `requestAnimationFrame` for motion (no heavy 3D/canvas libs unless clearly justified — perf matters; the page polls every 2s and animates continuously).
- **Data contract (fixed):** design only with fields in §5. New visuals must not require backend changes unless you call them out explicitly as a separate ask (the backend is faithful and tested; changes there reopen the Marx modeling).
- **Fidelity (fixed):** the §1 semantic invariants.
- **Performance:** smooth on a laptop with the engine running; respect `prefers-reduced-motion`.
- **Free to change:** all composition, layout, hierarchy, motion design, the descent, typography scale, the abode/surface spatial relationship, and `atlas.css` wholesale. The `*.tsx` are restructurable as long as they consume the same snapshot.

---

## 8. Deliverable

Whatever form suits the design pass — annotated mockups, a restyle of `atlas.css` against the real tokens, a motion/transition spec, and/or a reworked component structure. Most valuable:
1. A **direction** for the four axes (atmosphere, motion, info hierarchy, concept) — ideally one strong point of view, not a menu.
2. The **descent** designed as the signature moment.
3. The **abode** re-laid-out with a clear hero + the immiseration trend promoted to a real chart.
4. A first concrete styling pass that actually uses Playfair Display, Plex Mono figures, and the `:root` tokens (fixing the off-brand drift noted in §3).

Then it can come back to engineering to build against this same data contract.
