---
chapter: 12
title: "The Concept of Relative Surplus Value"
status: proposed
primary_service: simulation-engine
---

## Concepts → types

| Marx term | Go identifier | Kind | Package | Notes |
|---|---|---|---|---|
| absolute surplus-value | `AbsoluteSurplusValue` | type | `simulation-engine/internal/production` | Surplus produced by prolonging the working day; `LabourMinutes` alias |
| relative surplus-value | `RelativeSurplusValue` | type | `simulation-engine/internal/production` | Surplus produced by shortening necessary labour via productivity increase; `LabourMinutes` alias |
| working day | `WorkingDay` | type | `simulation-engine/internal/production` | Total length of the working day; composed of `NecessaryLabour` + `SurplusLabour` |
| necessary labour-time | `NecessaryLabour` | type | `simulation-engine/internal/production` | `LabourMinutes` — portion reproducing value of labour-power |
| surplus labour-time | `SurplusLabour` | type | `simulation-engine/internal/production` | `LabourMinutes` — portion producing surplus-value for capital |
| value of labour-power | `LabourPowerValue` | type | `simulation-engine/internal/production` | `LabourMinutes` — socially necessary to reproduce the worker; determines `NecessaryLabour` |
| productiveness of labour | `ProductivityFactor` | type | `simulation-engine/internal/production` | `float64 > 0`; mirrors `ProductivityChange.Factor` in commodity package |
| individual value | `IndividualValue` | type | `simulation-engine/internal/production` | `LabourMinutes` — cost under one producer's conditions |
| social value | `SocialValue` | type | `simulation-engine/internal/production` | `LabourMinutes` — cost under average social conditions; the real value |
| extra surplus-value | `ExtraSurplusValue` | func | `simulation-engine/internal/production` | `func(IndividualValue, SocialValue, Quantity) LabourMinutes` — temporary gain from sub-social individual value |
| rate of surplus-value | `RateOfSurplusValue` | func | `simulation-engine/internal/production` | `func(SurplusLabour, NecessaryLabour) float64` — ratio s/v |
| shortening necessary labour | `ShortenNecessaryLabour` | func | `simulation-engine/internal/production` | `func(WorkingDay, LabourPowerValue) WorkingDay` — recomputes split after productivity rise |
| inverse law (value ∝ 1/productivity) | `ApplyProductivityToSNLT` | func | `simulation-engine/internal/production` | wraps `ProductivityChange.Apply`; asserts value falls as productivity rises |

## Fixtures

- **§1** `"let the whole line a c ... represent a working day of 12 hours; the portion a b 10 hours of necessary labour, and the portion b c 2 hours of surplus-labour"` → `WorkingDay{Total: 720, NecessaryLabour: 600, SurplusLabour: 120}` (all `LabourMinutes`)

- **§1** `"if now ... we move the point b to b', b c becomes b' c; the surplus-labour increases by one half, from 2 hours to 3 hours, although the working day remains as before at 12 hours"` → `ShortenNecessaryLabour(wd, 540)` returns `WorkingDay{Total: 720, NecessaryLabour: 540, SurplusLabour: 180}`

- **§1** `"one working-hour be embodied in sixpence, and the value of a day's labour-power be five shillings, the labourer must work 10 hours a day"` → `LabourPowerValue(300)` / `sixpence(10)` == `NecessaryLabour(600)`; surplus = `120`

- **§1** `"if, in consequence of increased productiveness, the value of the necessaries of life fall ... the value of a day's labour-power be thereby reduced from five shillings to three, the surplus-value increases from one shilling to three"` → reducing `LabourPowerValue` from `300` to `180` LabourMinutes changes `NecessaryLabour` from `600` to `360`, `SurplusLabour` from `120` to `360`

- **§1** `"let some one capitalist contrive to double the productiveness of labour, and to produce in the working day of 12 hours, 24 instead of 12 such articles ... the individual value of these articles is now below their social value"` → `IndividualValue(30)` vs `SocialValue(60)`; `ExtraSurplusValue(IndividualValue(30), SocialValue(60), Quantity(24))` reflects 3-shilling daily extra gain

- **§1** `"the capitalist who applies the improved method ... sells his commodity at its social value of one shilling, he sells it for threepence above its individual value, and thus realises an extra surplus-value of threepence"` → `ExtraSurplusValue` per article = `SocialValue − IndividualValue` when sold at social value

## Invariants

- `WorkingDay.NecessaryLabour + WorkingDay.SurplusLabour == WorkingDay.Total` [§1]
- `ShortenNecessaryLabour(wd, newLPV).SurplusLabour > wd.SurplusLabour` when `newLPV < wd.NecessaryLabour` [§1]
- `RateOfSurplusValue(sl, nl)` increases strictly as `ProductivityFactor` rises (holding `WorkingDay.Total` constant) [§1]
- `ExtraSurplusValue(iv, sv, qty) == 0` when `iv == sv` [§1]

## Scope

### This chapter builds
- Services: `simulation-engine`
- New domain types:
  - `WorkingDay` — split of total day into `NecessaryLabour` + `SurplusLabour`, all `LabourMinutes`
  - `AbsoluteSurplusValue` — `LabourMinutes` named type; surplus from extending the day
  - `RelativeSurplusValue` — `LabourMinutes` named type; surplus from contracting necessary labour
  - `LabourPowerValue` — `LabourMinutes` named type; SNLT to reproduce the worker
  - `IndividualValue` — `LabourMinutes` under one producer's conditions
  - `SocialValue` — `LabourMinutes` under average social conditions
  - `ProductivityFactor` — `float64`; multiplicative scalar (same semantics as `ProductivityChange.Factor` in commodity package, lifted into production domain)
- New HTTP endpoints:
  - `POST /v1/production/working-day` — record working-day split from total + labour-power value; returns `WorkingDay`
  - `POST /v1/production/working-day/shorten` — apply a new `LabourPowerValue`; returns updated `WorkingDay` and `RelativeSurplusValue` delta
  - `GET /v1/production/rate-of-surplus-value?necessary=&surplus=` — stateless rate computation
  - `POST /v1/production/extra-surplus-value` — stateless extra-surplus-value probe (individual value, social value, quantity)
- React: "Ch. 12 — Relative Surplus Value" panel with working-day split bar (necessary vs surplus), rate-of-surplus-value readout, and extra-surplus-value probe form

### Explicitly deferred to later chapters
- Co-operation as productivity mechanism (Ch. 13) — named in Ch. 12 closing but Marx defers analysis
- Machinery as the dominant technical mode of relative surplus-value (Ch. 15) — explicitly deferred by Marx at chapter end
- Full wage determination and reserve army dynamics — competition footnotes belong to wages theory (Ch. 17-20)
- Per-agent `ProductivityFactor` on `Agent` — requires a production domain on agents; unlocked by Ch. 13-14
