# Ch.25 `RunGeneralLaw`: organic composition by mechanism, not additive bump

- **Date:** 2026-06-03
- **Issue:** [#89](https://github.com/theding0x/capital-simulator/issues/89)
- **Service:** simulation-engine · package `internal/simulation`
- **Status:** Approved (Direction 1)

## Problem

`RunGeneralLaw` (general_law.go) advances the organic composition between
periods by adding a scalar to the *aggregate* ratio, then hard-clamping:

```go
sv := ProduceSurplus(v, s.SurplusRate)
newOCRatio := oc.Ratio + s.ProductivityGrowth
if newOCRatio >= 1.0 {
    newOCRatio = 0.99
}
additional := SplitSurplus(sv, s.AccumulationRate, newOCRatio)
```

This is a placeholder, not Marx's §25 mechanism. The "rising organic
composition" is modelled as an additive bump to the next split ratio, bounded
only by a **fictional** `0.99` ceiling — there is no mechanism by which OC
asymptotes there. The clamp also masks the §25 question of what happens to the
labour-demand path as OC approaches full automation (where the reserve army
balloons).

## Decision — stock-vs-flow composition (Direction 1, faithful to §25.1)

Replace the aggregate-bump with a **marginal (flow) composition** that the
*new* capital is invested at each period, and let the **aggregate** organic
composition emerge as the size-weighted average of old stock + accumulated
flows — which is exactly what `ComputeOrganicComposition(c, v) = c/(c+v)`
already computes once `c` and `v` carry the accreted flows.

### Mechanism

Carry a `marginal` composition across periods:

1. Initialise `marginal` to the starting aggregate OC: `c₀ / (c₀ + v₀)`.
2. Each period, productivity growth raises the marginal composition by closing
   a fraction of the remaining gap to full automation:

   ```
   marginal += ProductivityGrowth × (1 − marginal)
   ```

3. Accumulate the period's surplus at `marginal` via the existing
   `SplitSurplus`; `c` and `v` grow, and the next period's
   `ComputeOrganicComposition(c, v)` is the size-weighted average.

`ProductivityGrowth` is interpreted as the **fraction of the remaining
labour-share converted to machinery per period** (a rate in `[0, 1)`); the
implementation bounds the effective fraction to that domain so the marginal
composition can never overshoot 1.

### Why this is correct

- **The fictional clamp disappears.** `marginal += g(1−marginal)` is a convex
  combination of `marginal` and `1` for `g ∈ [0,1]`, so the marginal composition
  asymptotes to — but never reaches — full automation by construction. The
  aggregate trails strictly below it. No `0.99` ceiling.
- **It is the Marx mechanism.** New capital embodies a higher composition than
  the existing stock; the stock ratio compounds upward as high-composition
  flows accrete onto it (§25.1, the stock-vs-flow split).
- **The reserve army balloons naturally** as OC climbs toward its asymptote,
  rather than OC freezing at `0.99` and the labour path going flat.
- **Grounded in the existing labour/surplus packages** (`ProduceSurplus` /
  `SplitSurplus`) — exactly what the issue's acceptance criterion asks for.

### Behaviour preserved

- `ProductivityGrowth = 0` ⇒ `marginal` stays at `c₀/(c₀+v₀)` ⇒ every flow added
  at the starting composition ⇒ aggregate OC flat (unchanged-composition path).
- `ProductivityGrowth > 0` ⇒ each flow added above the current aggregate ⇒
  aggregate OC strictly rises each period, bounded below 1.
- `OrganicComposition == c/(c+v)` per-snapshot invariant unchanged.
- Variable capital stays positive every period (`marginal < 1`), so labour
  demand keeps growing slower than total capital — labour share falls.

## Scope / blast radius

- `services/simulation-engine/internal/simulation/general_law.go` — the
  `RunGeneralLaw` loop body and its doc comment; remove the additive-bump /
  clamp lines.
- `services/simulation-engine/internal/simulation/general_law_test.go` — new
  tests for the mechanism (below). Existing tests stay green.

The HTTP response series is recomputed on the fly from stored scenario
parameters (`toGeneralLawScenarioResponse → RunGeneralLaw`); struct and JSON
shapes are unchanged. Therefore:

- **No** migration (no schema change; scenarios store inputs, not the series).
- **No** handler, routes, gateway, `types.ts`, `api.ts`, or web-panel change.
- **No** seed change.

## Test plan (TDD — these are written first)

1. **Clamp removed / bounded asymptote.** A long run (e.g. 40 periods) with high
   `ProductivityGrowth` keeps every snapshot's OC strictly `< 1` and never lands
   exactly on `0.99`; OC increments shrink period over period (asymptotic), and
   `VariableCapital > 0` every period.
2. **Stock-vs-flow / decelerating rise.** With `ProductivityGrowth > 0`, each
   period's aggregate OC is strictly less than the next period's, and the
   per-period increase decelerates — confirming a weighted-average approach
   rather than a fixed additive step.
3. **Reserve army monotonic growth.** With rising OC and fixed worker supply,
   `ReserveArmySize` is non-decreasing and ends strictly above where it started
   (the §25 content: labour demand path as OC climbs).
4. Existing tests (`UnchangedComposition`, `RisingComposition`,
   `LabourShareFalls`, `OCInvariant`, `ZeroPeriods`) remain green unchanged.

## Out of scope

- **Direction 2** (wiring `ProductivityGrowth` through agent-service Ch.16
  `ReduceNecessaryLabour` / Ch.12 productivity-factor). Heavier cross-service
  coupling that edges into ADR-0001's *deferred* cross-node value-flow scope;
  not needed to close this gap faithfully.
- The MELT / value↔money bridge of **ADR-0001** — a different "composition gap"
  (cross-node value flow), untouched here.
