---
chapter: 25
title: "The General Law of Capitalist Accumulation"
status: proposed
primary_service: simulation-engine
---

## Concepts → types

| Marx term | Go identifier | Kind | Package | Notes |
|---|---|---|---|---|
| value composition of capital | `ValueComposition` | type | `simulation` | `ConstantCapital Pence; VariableCapital Pence` — monetary ratio c:v |
| technical composition of capital | `TechnicalComposition` | type | `simulation` | `MeansOfProductionUnits int64; LabourPowerUnits int64` — material ratio of machines to workers |
| organic composition of capital | `OrganicComposition` | type | `simulation` | `Ratio float64` — value composition determined by technical composition; `c / (c + v)` |
| concentration of capital | `Concentration` | type | `simulation` | `TotalCapital Pence; Firms int64` — average capital per firm as accumulation proceeds |
| centralisation of capital | `Centralisation` | type | `simulation` | `AcquiredCapital Pence; FirmsAbsorbed int64` — merger/takeover of many small capitals by few |
| relative surplus population | `RelativeSurplusPopulation` | type | `simulation` | `Floating int64; Latent int64; Stagnant int64` — three strata of the reserve army |
| industrial reserve army | `IndustrialReserveArmy` | type | `simulation` | `Size int64; RelativeProportion float64` — ratio of reserve to active labour army |
| demand for labour | `LabourDemand` | type | `simulation` | `Workers int64` — derived from capital stock and organic composition |
| compute organic composition | `ComputeOrganicComposition` | func | `simulation` | `func(vc ValueComposition) OrganicComposition` — `vc.ConstantCapital / (vc.ConstantCapital + vc.VariableCapital)` |
| compute labour demand | `ComputeLabourDemand` | func | `simulation` | `func(totalCapital Pence, oc OrganicComposition) LabourDemand` — demand = variable capital / wage per worker |
| compute reserve army | `ComputeReserveArmy` | func | `simulation` | `func(supply int64, demand LabourDemand) IndustrialReserveArmy` — supply minus demand |

## Fixtures

- **§1 composition unchanged** `"accumulation doubles total capital while composition remains same"` → if `vc == ValueComposition{8000,2000}` and capital doubles to 20,000, then `ComputeLabourDemand(20000, ComputeOrganicComposition(vc))` should return double the original demand — additional workers absorbed, reserve army shrinks
- **§2 rising organic composition** as machinery displaces labour, `OrganicComposition` rises: `ComputeOrganicComposition(ValueComposition{9000,1000}).Ratio == 0.9 > ComputeOrganicComposition(ValueComposition{8000,2000}).Ratio == 0.8` — with same total capital, demand for labour falls
- **§3 reserve army** `"the greater the social wealth, the functioning capital ... the greater also is the industrial reserve army"` → `ComputeReserveArmy(supply, demand)` grows as total capital grows and organic composition rises simultaneously
- **§4 strata** three forms: floating (factory workers discharged by machinery); latent (agricultural population ready to enter industry); stagnant (irregular casual labour) → `RelativeSurplusPopulation{Floating:f, Latent:l, Stagnant:s}` where f+l+s == total reserve army

## Invariants

- `ComputeOrganicComposition(vc).Ratio == float64(vc.ConstantCapital) / float64(vc.ConstantCapital+vc.VariableCapital)` — always in [0, 1)
- `ComputeLabourDemand(total, oc).Workers == int64(float64(total) * (1-oc.Ratio) / wagePence)` — variable portion divided by wage
- `IndustrialReserveArmy.Size == workerSupply - ComputeLabourDemand(total, oc).Workers` — positive when supply exceeds demand; can be 0 but not negative (workers don't disappear into negative employment)
- `IndustrialReserveArmy.RelativeProportion == float64(Size) / float64(ComputeLabourDemand.Workers)` — reserve relative to active army; rises as organic composition rises

## Scope

### This chapter builds
- Services: `simulation-engine`
- New domain types:
  - `ValueComposition` — `ConstantCapital, VariableCapital Pence`; monetary c:v
  - `TechnicalComposition` — `MeansOfProductionUnits, LabourPowerUnits int64`; material ratio
  - `OrganicComposition` — `Ratio float64`; c/(c+v); determined by technical composition
  - `RelativeSurplusPopulation` — `Floating, Latent, Stagnant int64`; three strata
  - `IndustrialReserveArmy` — `Size int64; RelativeProportion float64`; total reserve and its ratio to active army
  - `Concentration` — average capital per firm at a given period
  - `Centralisation` — capital absorbed by takeover vs. accumulation by surplus-value
- New functions:
  - `ComputeOrganicComposition(vc ValueComposition) OrganicComposition` — pure
  - `ComputeLabourDemand(totalCapital Pence, oc OrganicComposition, wagePence int64) LabourDemand` — pure
  - `ComputeReserveArmy(supply int64, demand LabourDemand) IndustrialReserveArmy` — pure
- New HTTP endpoints:
  - `POST /v1/accumulation/organic-composition` — stateless; compute OC from value composition
  - `POST /v1/accumulation/labour-demand` — stateless; compute demand from capital and OC
  - `POST /v1/accumulation/reserve-army` — stateless; compute reserve army size and proportion
  - `POST /v1/accumulation/scenarios` — persist a named accumulation scenario (composition + growth trajectory)
  - `GET /v1/accumulation/scenarios/{id}` — retrieve scenario with computed series
- Migration: `00003_ch25_accumulation_scenarios.sql` — `accumulation_scenarios` table
- React: add a "Ch. 25 — General Law" panel; inputs for starting capital composition, accumulation rate, productivity growth; charts showing organic composition rising over time, labour demand vs. supply, and reserve army size; highlight the inverse relationship between capital growth and labour absorption when OC rises

### Explicitly deferred to later chapters
- Empirical illustrations (sections A-F: England 1846-66, agricultural proletariat, Ireland) — historical data, not simulation logic
- Credit system as mechanism of centralisation — deferred to Book III
- Pauperism as distinct from relative surplus population — described qualitatively; not a separate domain type
- Differential rent as affected by agricultural revolution — deferred to Book III
