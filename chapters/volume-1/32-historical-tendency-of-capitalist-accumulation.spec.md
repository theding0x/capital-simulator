---
chapter: 32
title: "Historical Tendency of Capitalist Accumulation"
status: proposed
primary_service: simulation-engine
---

## Concepts → types

| Marx term | Go identifier | Kind | Package | Notes |
|---|---|---|---|---|
| petty property | `PettyProperty` | type | `simulation` | `Producers int64; CapitalPerProducerPence Pence` — self-owned, labour-based property; precondition for petty industry |
| capitalist private property | `CapitalistPrivateProperty` | type | `simulation` | `Firms int64; TotalCapitalPence Pence; WageLabourers int64` — property resting on exploitation of others' unpaid labour |
| centralisation | `CentralisationStep` | type | `simulation` | `FirmsAbsorbed int64; CapitalConcentratedPence Pence` — one step in the expropriation of many capitalists by few |
| negation | `Negation` | type | `simulation` | `Stage string; Description string` — dialectical stage: `"petty-property"`, `"capitalist-expropriation"`, `"socialised-property"` |
| negation of negation | `NegationOfNegation` | func | `simulation` | `func(stages []Negation) Negation` — returns the third stage (socialised property) as the dialectical resolution |
| accumulation trajectory | `AccumulationTrajectory` | type | `simulation` | `Periods []CentralisationStep; FinalFirms int64; FinalCapitalPence Pence; ReserveArmySize int64` — long-run outcome of the general law |

## Fixtures

- **§ first negation** `"self-earned private property ... is supplanted by capitalistic private property"` → `Negation{Stage:"capitalist-expropriation", Description:"many producers expropriated by few capitalists"}`
- **§ negation of negation** `"capitalist production begets ... its own negation ... gives him individual property based on cooperation and possession in common"` → `NegationOfNegation([petty, capitalist, socialised]).Stage == "socialised-property"`
- **§ centralisation spiral** as total capital doubles while organic composition rises: `AccumulationTrajectory` shows declining `FinalFirms` and rising `FinalCapitalPence`; the few usurp the many's capital
- **§ the knell** `"monopoly of capital becomes a fetter ... The expropriators are expropriated"` — `AccumulationTrajectory.FinalFirms == 1` marks the logical terminus; in simulation this is an asymptotic limit, not a reachable state

## Invariants

- `NegationOfNegation(stages)` requires exactly three stages in the historical sequence: petty → capitalist → socialised
- `trajectory.FinalFirms < trajectory.Periods[0].FirmsAbsorbed` — centralisation always reduces the number of capitals
- `trajectory.FinalCapitalPence > sum(step.CapitalConcentratedPence)` — concentration of ownership does not destroy value; it redistributes it
- Each `CentralisationStep.CapitalConcentratedPence == totalCapital * (1 - 1/e^(step.FirmsAbsorbed/initialFirms))` is an approximation; the exact form is left to the implementation

## Scope

### This chapter builds
- Services: `simulation-engine`
- New domain types:
  - `PettyProperty` — producers, capital per producer; the first stage
  - `CapitalistPrivateProperty` — firms, total capital, wage labourers; the second stage
  - `CentralisationStep` — one step of capital concentration: firms absorbed, value centralised
  - `Negation` — named dialectical stage with description
  - `AccumulationTrajectory` — long-run series of centralisation steps and final state
- New functions:
  - `NegationOfNegation(stages []Negation) Negation` — returns the third dialectical stage
  - `RunCentralisation(initial CapitalistPrivateProperty, steps int64) AccumulationTrajectory` — simulate progressive absorption of firms
- New HTTP endpoints:
  - `POST /v1/accumulation/centralisation` — run centralisation simulation; returns `AccumulationTrajectory`
  - `GET /v1/accumulation/negation-of-negation` — stateless; return the three Negation stages as structured data
- React: add a "Ch. 32 — Historical Tendency" panel; timeline showing the three stages (petty property → capitalist property → socialised property); animated chart of firms declining and capital concentrating; endpoint of the simulation series links back to Ch. 25's reserve army model

### Explicitly deferred to later chapters
- The specific mechanism of the transition to socialised property — Marx gestures at it but Capital Vol. I does not model it; beyond scope
- Rate of profit and its tendency to fall — mentioned as connected to centralisation but Book III topic
