# Atlas — The General Law in Motion (Design Spec)

- **Date:** 2026-06-04
- **Status:** Approved (brainstorm) — pending implementation plan
- **Builds on:** the shipped Phase 1 "living field" (`2026-06-04-atlas-observatory-design.md`), on branch `feature/atlas-observatory`. **Supersedes** that spec's old Phase 2 (inspector) and Phase 3 (levers) with a more faithful design.

## 1. Motivation

Phase 1 shipped the surface — a field of capital-orbits polling a live snapshot — but three critiques (from watching it run) exposed that it shows the *value-form* while missing what's actually in motion:

1. **Nothing expands as ticks accumulate.** Orbit size = a *static* seeded `total_pence`; the engine's tickers advance unrelated entities, so capitals never grow. The promised "spiral of accumulation" was never implemented — there is no accumulation mechanic.
2. **The orbit spin was arbitrary.** Rotation came from a hash of the capital's id. It must have a concrete source.
3. **No alienation / class struggle.** The view shows M/P/C, c/v/s, p̄′ — precisely the surface that, in Marx, *masks* the social relation (commodity fetishism). The class relation — surplus as unpaid labour, the reserve army, the antagonism — is absent.

This design fixes all three by reconceiving the Atlas as **the General Law of Capitalist Accumulation (Vol. I Ch. 25) in motion**: capital spirals outward by pumping surplus from living labour, which swells the reserve army, depresses wages, intensifies exploitation, and drives further accumulation. The antagonism becomes the engine.

## 2. Decided design spine

| Decision | Choice |
|---|---|
| Issue #1 — growth | Capitals **actually accumulate**: each turnover a share of surplus is capitalised → `total_pence` grows → the orbit spirals. Driven by the General Law. |
| Issue #2 — motion | The **ring is static** (M/P/C are fixed coexistence positions); **value-dots travel the circumference**; one lap = **one turnover** (speed ← real turnover); the dot **lingers in production (P)** and darts through the buy/sell acts (time-of-production vs time-of-circulation, Vol. II). |
| Issue #3 — class relation | The working class is a **co-equal pole in the simulation**, **revealed in the view** via the **hidden abode** (cross the threshold from the value-circuit surface into production). |
| Reveal scope | **Aggregate abode** — one threshold for the whole field; descend into the production relation beneath them all (the class as a class). |
| Abode hero | **The social working day**: necessary (paid, v) vs surplus (unpaid, s) labour — their ratio is the **rate of exploitation s/v**, widening as the law intensifies. |
| Driving law | **The General Law of Accumulation (Ch. 25)** — the full feedback loop, run as a backend ticker. |
| Levers | **Working day · wage level · accumulation rate** — perturb the live loop. |

## 3. Fidelity rationale

- **The hidden abode** — Vol. I Ch. 6 fin.: *"Let us therefore… leave this noisy sphere, where everything takes place on the surface and in view of all men, and follow them into the hidden abode of production, on whose threshold there hangs the notice 'No admittance except on business.'"* The surface (circulation, the value-forms, fetishism) hides the abode (production, the class relation). The view literalises the descent.
- **The General Law** — Vol. I Ch. 25: accumulation raises the organic composition, which relatively expels labour into the **industrial reserve army**, which presses wages down and raises the rate of surplus-value — the law of **immiseration** and concentration. This is the law of motion of capital *and* the working class together.
- **The working day & exploitation** — Vol. I Chs. 7–10: value produced = v + s; the rate of surplus-value s/v *is* the rate of exploitation; surplus labour is unpaid labour.
- **Turnover & time-in-stage** — Vol. II Chs. 7, 12–14: the circuit's stages take real time (time of production vs time of circulation); the capital tied up in a stage is proportional to the time spent there — so a dot lapping once per turnover and lingering in P is faithful.
- **Accumulation as spiral** — Vol. I Ch. 24 / Vol. II Ch. 21: capitalised surplus enlarges the circuit; the closed circle of simple reproduction becomes the spiral.
- **Alienation** — Vol. I Chs. 23–25: accumulated capital is **dead labour dominating living labour**; the worker produces the alien power that dominates them. Surfaced as: the very surplus extracted from labour (the abode) is what makes capital (the surface orbits) grow and tower over it.

## 4. Architecture

Built on the shipped Phase 1 (`simulation-engine` snapshot read-model + `web/src/atlas/`). Three backend additions, three frontend additions.

### Backend (`simulation-engine`)
- **The General-Law ticker** — a new `engine.Ticker` registered in the scheduler (`internal/engine/scheduler.go`) that runs one period of the law per pass:
  - **Accumulate:** each `IndustrialCapital` capitalises `α · s` of its surplus into `total_pence` (α = accumulation rate); its `StageDistribution` rescales to the new magnitude (the spiral).
  - **Run the law:** recompute the aggregate using existing `internal/simulation/general_law.go` — `ComputeOrganicComposition`, `ComputeLabourDemand`, `ComputeReserveArmy` (and/or advance a `GeneralLawScenario` via `RunGeneralLaw`). Rising composition → falling relative labour demand → growing reserve army.
  - **Press wages:** `agent.ReserveArmyPressure.CompressPence` (from `reserve_army_pressure.go`) depresses v as the reserve army grows → s/v rises.
  - **Record** one aggregate `GeneralLawSnapshot`-style row per period (the immiseration time-series).
- **Snapshot v2** — extend `GET /v1/observatory/snapshot` (Phase 1) with: per-capital `turnover_number` and a stage-timing split (for dot pacing); and an `abode` block (aggregate class state). Per-capital `variable_pence`/`constant_pence` are carried so s/v is computable (today only c+v survives — see §11).
- **Lever params** — a small settable parameter set (working-day necessary/surplus split, wage level, accumulation rate α) the ticker reads, via `POST /v1/observatory/levers` (or reuse `POST /v1/accumulation/*` + `POST /v1/production/working-day`). Effects appear in subsequent snapshots.

### Frontend (`web/src/atlas/`)
- **Surface motion + growth** — update `Orbit.tsx`: keep the static ring; drive dot travel by `turnover_number`; pace dots with **variable per-arc speed** (linger in P); ease orbit radius toward the growing `total_pence` (smooth spiral between snapshots).
- **Threshold** — a "descend into production" control + an animated transition (surface ↕ abode), framed by *"No admittance except on business."*
- **The abode** — new components: `WorkingDay` (necessary|surplus bar = s/v), `SurplusExtraction` (gold rising out to the surface), `LivingLabour` (Σv employed), `ReserveArmy` (reservoir + wage-pressure), `GeneralLawTrend` (immiseration sparkline).
- **Levers** — working-day / wage / accumulation-rate controls that POST and let you watch the law respond.

## 5. The General-Law mechanic (per tick)

Grounded in `general_law.go` + `reserve_army_pressure.go` + `working_day.go`:

```
each scheduler pass (one period):
  s′      = surplus-rate from the social working day (necessary vs surplus minutes),
            intensified by reserve-army pressure
  for each capital:
      s        = v · s′
      total   += α · s                 # accumulation → spiral (issue #1)
      rescale StageDistribution to new total
  C        = Σ total                    # total social capital
  oc       = ComputeOrganicComposition(C-stock)   # rises with accumulation
  demand   = ComputeLabourDemand(C, oc, wage)     # falls relatively
  reserve  = ComputeReserveArmy(labourSupply, demand)   # grows
  pressure = ReserveArmyPressure(reserve)               # → CompressPence(wage)
  record GeneralLawSnapshot{ C, oc, reserve, wage, s′, … }   # immiseration series
```

The loop closes: more accumulation → higher composition → bigger reserve army → lower wages → higher s′ → more surplus → more accumulation.

## 6. Snapshot v2 contract (additions to Phase 1)

```json
{
  "tick": 142, "running": true, "interval_ms": 5000,
  "capitals": [ { "...phase 1 fields...",
                  "variable_pence": 0, "constant_pence": 0,
                  "turnover_number": 0,
                  "production_time_share_bp": 0 } ],
  "aggregate": { "...phase 1 fields..." },
  "abode": {
    "total_variable_pence": 0,            // Σv = wages = paid labour
    "total_surplus_pence": 0,             // Σs = unpaid labour
    "rate_of_exploitation_bp": 0,         // s/v
    "necessary_labour_minutes": 0,
    "surplus_labour_minutes": 0,
    "organic_composition_bp": 0,          // c/v
    "reserve_army_count": 0,
    "reserve_army_pressure_bp": 0,
    "employed_count": 0,
    "law_series": [ { "period": 0, "wage_pence": 0, "rate_of_exploitation_bp": 0,
                      "reserve_army_count": 0, "organic_composition_bp": 0 } ]
  }
}
```
All `snake_case`, `Pence`/minutes/basis-points integers (no `float64`); arrays never null.

## 7. Data flow

The General-Law ticker runs server-side each ~5s, mutating the field (capitals grow) and recording the law series. The frontend polls snapshot v2 every ~2s: surface orbits ease toward new magnitudes (spiral) and pace dots by turnover (linger in P); the abode reflects the live class state; the immiseration trend appends. Levers POST → ticker params → the law visibly responds over the next passes.

## 8. Error handling

Inherits Phase 1 (hold last-good on poll failure; degraded banner on scheduler 500). The General-Law ticker advances as many capitals as it can and reports a joined error rather than aborting a pass (matching the existing `Ticker` contract). Accumulation guards: no negative capital; bounded growth per pass; reserve-army pressure clamped.

## 9. Testing

- **Go:** General-Law ticker unit tests on a Marx-fixture economy — accumulation grows capital, rising composition grows the reserve army, pressure raises s/v (the loop's direction), with `t.Parallel()`. Snapshot-v2 handler test: abode fields present, s/v = round(10000·Σs/Σv), arrays non-null. Reuse existing `general_law.go` tests as fixtures.
- **Web:** `npm run lint` + `build`; pure pacing/tween/format functions unit-test-ready; **Playwright on the booted stack** — orbits visibly grow over ~30s, a dot lingers in P, crossing the threshold reveals the abode with a live s/v, a lever moves the law.

## 10. Build slices (each a shippable PR on `feature/atlas-observatory`)

1. **Real growth + corrected motion (#1 + #2)** — the accumulation/growth ticker (capitals spiral); snapshot adds `turnover_number` + stage-timing + per-capital v; `Orbit.tsx` paces dots by turnover with linger-in-P and eases the growing radius. *The surface finally lives.*
2. **The hidden abode (#3)** — the full General-Law loop (composition → reserve army → wage pressure → s/v) + the `abode` snapshot block; the threshold transition and abode components (working day, living labour, reserve army, surplus extraction, immiseration trend).
3. **The levers** — working-day / wage / accumulation-rate perturbation wired to the ticker; the law responds live.

## 11. Open details to pin during planning

- **Per-capital v:** today only `demand = c+v` is stored. Carry organic composition (c, v) per capital — extend the seed + a store read (or a small `industrial_capital_composition` record). Needed for s/v and for accumulation raising composition.
- **Dot pacing source:** per-arc timing from the working-period / time-of-circulation stores (Vol. II Ch. 12–14), or the faithful proxy "time-in-stage ∝ value-in-stage" weighted to over-linger in P. Pick one.
- **Ticker persistence:** whether the General-Law ticker mutates the persisted `IndustrialCapital.total_pence` (sim evolves across restarts) or maintains a live overlay reset each boot. Recommend mutating persisted state on a fresh-seeded volume.
- **Aggregate labour supply & wage baseline:** source the initial labour supply and wage (agent-service labour-power / national-wage, or a seeded General-Law scenario).
- **Accumulation rate α default** and growth clamp per pass.

## 12. Out of scope (YAGNI)

- Per-capital abode (descend into one factory) — aggregate only for now.
- Server-push streaming; backend tick-interval setter.
- Re-opening the chapter panels beyond the existing `#/chapters` route.

## 13. Relationship to Phase 1

Phase 1 (the watch-only field) stays; this design makes it *move and mean*. The old inspector (Phase 2) and free levers (Phase 3) are replaced by the abode + General-Law levers here. All work continues on `feature/atlas-observatory`.
