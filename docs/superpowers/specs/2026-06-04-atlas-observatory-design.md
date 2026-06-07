# Atlas — The Whole Circuit in Motion (Design Spec)

- **Date:** 2026-06-04
- **Status:** Approved (brainstorm) — pending implementation plan
- **Topic:** A new full-page visualization that shows the simulated economy as a system *in motion* — the circuit `M—C(Lp+Mp)…P…C'—M'` as a living, self-expanding loop, at the scale of many capitals.

## 1. Motivation

All three volumes of *Capital* are implemented as backend domain logic and 106 chapter "control panels." The app reads as a set of static forms, which undersells what it models: **value in motion** — the self-expanding circuit. This initiative shifts focus from configuration to **visualization of the whole circuit in motion**.

What already exists and is reused:

- A **real tick scheduler** in `simulation-engine` (`internal/engine/scheduler.go`): periodic passes advance domain entities via registered tickers (`FactoryTicker`, `ReproductionTicker`, `PiecePriceTicker`). Controlled by `POST /v1/engine/start|stop`; observed by `GET /v1/engine/status` and `GET /v1/engine/ticks`.
- `IndustrialCapital` (`internal/circulation/industrial_capital.go`) — Vol. II Ch. 4's synthesis: each capital references a money / productive / commodity circuit and carries a live **`StageDistribution`** (`money_pence + production_pence + commodity_pence == total_pence`) plus an **`OrganicComposition`** (c, v). This is the per-capital animation source.
- Aggregate endpoints across services (industrial-capital supply/demand aggregate, reproduction social-aggregate, finance general rate of profit, etc.).

The current `web/src/components/CircuitSpine.tsx` + `TurnoverPlayer.tsx` are **decorative** — they cycle a highlighted node on a timer and drive no data. They are superseded by Atlas on the landing view.

## 2. Decided design spine

| Decision | Choice |
|---|---|
| Hero motif | **The faithful orbit** — a ring of three *coexisting* arcs (M·P·C) sized to the live `StageDistribution`; value circulates through them; the ring spirals outward as ΔM compounds. |
| Scope | **Many circuits at once** — a field of orbits; the **average rate of profit p̄′** is the centre of gravity that visibly forms from their spread. |
| Motion source | **Real engine + continuous animation** — the page controls the real scheduler and polls real magnitudes, interpolating between ticks at 60fps. |
| Placement | **New landing page ("Atlas")**; the 106 chapter panels demote to a drill-down library behind a lightweight route. |
| Interactivity | **Watch + a few global levers** — rate of surplus-value, accumulation rate, organic-composition drift. |
| Layout | **Observatory** — field centre-stage; a slim rail of levers + vital-signs; a bottom bar with transport + the tick heartbeat; click an orbit → inspector. |
| Data transport | **A new snapshot endpoint** assembles the whole field server-side; the client is a pure renderer. |

## 3. Fidelity rationale (why the faithful orbit)

The chosen spine — "value in self-expanding motion" — has a precise textual basis, and the orbit honors both halves:

- **Vol. I, Ch. 4** — value as the *automatic subject*: it "is constantly changing from one form into the other… assumes an automatically active character," expanding spontaneously. → a **circular movement that returns to its origin enlarged** (the ring; the spiral).
- **Vol. II, Chs. 1–4** — **Nebeneinander** (coexistence): an industrial capital is *not* at one stage at a time; "one part exists as money-capital… another as productive capital… a third as commodity-capital," the three circuits carried on side by side. → the **three arcs are present at once**, sized to `StageDistribution`. This is the correction over a single bead hopping stages (which depicts succession only).
- **Accumulation as spiral** — simple reproduction is a closed circle; reproduction with accumulation "becomes a spiral." → the ring **steps outward each completed turn**.

The c+v+s ledger (the production/composition reading) and the named-metamorphosis lane (the circulation reading) are **secondary readings** surfaced in the inspector — faithful to where they belong rather than overstated as the hero.

## 4. Architecture

A net-new Atlas landing page; one new backend read-model; everything else reused.

- **`simulation-engine`** (owns field data + scheduler):
  - New `GET /v1/observatory/snapshot` → `internal/transport/httpapi/observatory_handler.go`, registered in `routes.go` and constructed in `handler.go`. A **projection / read-model** over existing data — **no schema migration** for the endpoint itself.
  - A focused store read (`internal/store`) returning field capitals with their latest `StageDistribution` + `OrganicComposition` (added to the `Store` interface, `memory.go`, `mysql.go`).
  - Engine controls reused unchanged.
- **`api-gateway`**: explicit per-path proxy lines for `/v1/observatory/*` (and verify `/v1/engine/*`, `/v1/accumulation/*` are present — a missing line is a silent 502 that CI does not catch).
- **`web/`**: new `web/src/atlas/` module; a minimal route split (Atlas = `/`, existing dashboard = `/chapters`).

## 5. The snapshot contract

`GET /v1/observatory/snapshot`:

```json
{
  "tick": 142,
  "running": true,
  "interval_ms": 5000,
  "capitals": [
    {
      "id": "…",
      "total_pence": 0,
      "money_pence": 0,
      "production_pence": 0,
      "commodity_pence": 0,
      "constant_pence": 0,
      "variable_pence": 0,
      "surplus_pence": 0,
      "turnover_number": 0,
      "status": "active"
    }
  ],
  "aggregate": {
    "total_social_capital_pence": 0,
    "avg_rate_of_profit_bp": 0,
    "rate_of_surplus_value_bp": 0,
    "dept_i_pence": 0,
    "dept_ii_pence": 0,
    "accumulation_pence": 0
  }
}
```

- JSON tags `snake_case`; value-magnitudes in the existing `Pence` unit; rates in basis points (`_bp`), integer (no `float64`).
- **p̄′ = ΣS/ΣC** computed server-side from the field's own capitals (the average rate *is* the field's ΣS/ΣC), keeping it single-service and avoiding a cross-service hop.
- Arrays never `null` (empty slice).

**Details to pin during planning (flagged, not hand-waved):**
- Source of per-capital `surplus_pence`: derive `s = s′·v` from a rate of surplus-value, or read `SupplyDemandImbalance.ExcessPence`.
- Source of `turnover_number`: the turnover store vs. a sensible default when absent.
- Dept I/II split source: reproduction scheme stores vs. classification of capitals.

## 6. Frontend components (`web/src/atlas/`)

- `Atlas.tsx` — Observatory shell (the landing).
- `Orbit.tsx` — one **faithful orbit** (3 arcs sized to M/P/C, value-flow, spiral growth). Reused by field and inspector.
- `CircuitField.tsx` — the field of orbits arranged around the p̄′ centre of gravity.
- `VitalSigns.tsx` — rail readouts: total social capital, p̄′, s′, Dept I/II, accumulation.
- `TickHeartbeat.tsx` — engine status + ECG of recent ticks + transport (play / pause / speed).
- `Levers.tsx` — global levers (Phase 3; absent in Phase 1).
- `Inspector.tsx` — click-orbit drill-in (Phase 2).
- `useSnapshot.ts` — polling hook (snapshot → tween targets + tick number).
- `animation.ts` — rAF tween loop (pure, unit-testable easing + angle advance).
- `api.ts` / `types.ts` — add `getObservatorySnapshot`, engine controls, and mirror types (`ObservatorySnapshot`, `FieldCapital`, `AggregateVitals`, `EngineStatus`, `EngineTick`).

## 7. Data flow (real engine + continuous animation)

1. `useSnapshot` polls `/v1/observatory/snapshot` every ~2s; stores the response as **tween targets** plus the tick number.
2. `animation.ts` runs at 60fps: each orbit's arc fractions ease toward target M/P/C; ring radius eases toward target total (the spiral); angular phase advances continuously at a rate derived from `turnover_number`. Motion stays fluid between the ~2s data refreshes; ΔM compounding reads as the radius easing upward.
3. Transport: Play → `POST /v1/engine/start`; Pause → `POST /v1/engine/stop`; running state + tick from `GET /v1/engine/status`.
   - **Honesty note:** backend tick cadence is fixed by `SIM_TICK_INTERVAL`; the ×1/×5 speed control governs **animation tempo** in v1, not the engine clock. A backend interval-setter is a possible later addition.
4. Heartbeat: poll `GET /v1/engine/ticks?limit=60`; plot `entities_advanced` as the ECG.

## 8. Error handling

- Snapshot poll failure → hold last-good field, show a subtle "stale / reconnecting" indicator, never blank (matches the existing "panels handle their own error display" ethos).
- Engine paused → orbits hold position; transport shows paused.
- Scheduler unavailable (engine handlers 500) → degraded banner.
- Empty field → friendly empty state (should not occur given seeds).

## 9. Seeds & fidelity

Phase 1 ships a **seed migration of several branches with differing organic compositions** (Marx-faithful — e.g. the five spheres of Vol. III Ch. 9, or the high/low-composition pair of Ch. 8), so the field comes up rich on a fresh MySQL volume and the average rate visibly *forms* from real spread. Seed IDs follow the `5eed00000000000000…` convention with a `-- +goose Down` deleting only those IDs.

## 10. Testing

- **Go:** `observatory_handler_test.go` (httptest + in-memory store): 200 well-formed snapshot; **p̄′ = ΣS/ΣC correct on a Marx composition fixture**; arrays never `null`; store-aggregation unit test with named Marx fixtures. `t.Parallel()` throughout.
- **Web:** `npm run lint` (tsc) + `npm run build`; pure tween/format functions unit-tested; **Playwright on the booted stack** (`/` renders orbits; Play flips engine status; ECG advances) — honoring the always-boot-for-E2E rule.

## 11. Phased build (each phase = one shippable PR)

1. **Living field (watch-only)** — snapshot endpoint + gateway proxy + branch seed; Atlas with field + vital-signs + transport + heartbeat; router split; continuous animation. *Delivers "the whole in motion" on its own.*
2. **Inspector** — click-orbit drill-in (single faithful orbit + c+v+s ledger + turnover + status + "open in chapters" bridge to the demoted library).
3. **Levers** — wire global levers to real `POST /v1/accumulation/*` perturbations; the field responds over subsequent snapshots.

## 12. Open decisions (defaulted unless overridden)

- **Router:** a minimal hand-rolled 2-view switch (no `react-router` dependency — honors the project's anti-router note; upgradeable later) rather than pulling in a router library.
- **Branch / PR:** `feature/atlas-observatory` + a normal PR. This is a feature initiative outside the Marx chapter cadence, so the chapter workflow (branch `volume-X/chapter-Y`, roadmap row, chapter spec) does not apply; the build/test/PR discipline still does.

## 13. Out of scope (YAGNI for now)

- Server-push streaming (SSE/WebSocket) — revisit after the snapshot shape is proven.
- A backend tick-interval setter — animation tempo covers v1.
- Per-capital historical time-series / scrubbing — the live field is the v1 goal.
- Any change to the 106 chapter panels beyond moving them behind the `/chapters` route.
