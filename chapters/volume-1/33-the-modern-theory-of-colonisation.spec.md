---
chapter: 33
title: "The Modern Theory of Colonisation"
status: proposed
primary_service: simulation-engine
---

## Concepts → types

| Marx term | Go identifier | Kind | Package | Notes |
|---|---|---|---|---|
| colonial labour market | `ColonialLabourMarket` | type | `simulation` | `Colony string; FreeLabourers int64; WagePence Pence; LandAvailable bool` — the colonial context where free land undermines wage dependence |
| sufficient price of land | `SufficientPrice` | type | `simulation` | `PricePerAcre Pence; MinimumWageYears int64` — Wakefield's artificial land price designed to delay independence |
| wage-worker independence | `WageWorkerIndependence` | type | `simulation` | `WorkerID string; YearsWorked int64; SavingsPence Pence; BecameLandowner bool` — whether a worker has saved enough to buy land and exit the labour market |
| systematic colonisation | `SystematicColonisation` | type | `simulation` | `Colony string; SufficientPrice SufficientPrice; ImportedLabourers int64; LandFundPence Pence` — Wakefield's scheme |
| labour market regulation | `ColonialLabourRegulation` | func | `simulation` | `func(market ColonialLabourMarket, scheme SystematicColonisation) ColonialLabourMarket` — apply the sufficient price scheme to a colonial labour market and return the regulated state |

## Fixtures

- **§ Peel's fiasco** `"Mr. Peel took with him ... means of subsistence and of production to the amount of £50,000 ... 300 persons of the working class ... 'Mr. Peel was left without a servant to make his bed'"` → `ColonialLabourMarket{Colony:"Swan River", FreeLabourers:300, LandAvailable:true}` → without a `SufficientPrice`, all 300 become independent; Peel's £50,000 capital cannot self-expand
- **§ sufficient price** Wakefield's remedy: set `SufficientPrice{PricePerAcre:artificiallyHigh, MinimumWageYears:N}` so workers must work for N years before saving enough to buy land; `ColonialLabourRegulation(market, scheme)` returns a market where `FreeLabourers` remain wage-workers for at least N periods
- **§ independence threshold** at normal colonial wages: `WageWorkerIndependence{YearsWorked:5, SavingsPence:5*annualWage}` — after 5 years a worker can buy land and exit → `BecameLandowner:true`; with sufficient price 10× higher: `MinimumWageYears:50` → worker remains dependent for 50 years

## Invariants

- `WageWorkerIndependence.BecameLandowner == (savings >= sufficientPrice.PricePerAcre * desiredAcres)` — exit from labour market is determined by savings relative to land price
- `ColonialLabourRegulation(market, scheme).FreeLabourers == market.FreeLabourers` — the total number of labourers does not change; only their independence status changes
- `scheme.LandFundPence == sum(price paid for land) used to import replacement labourers` — Wakefield's self-financing scheme: land revenue funds fresh labour imports; `LandFundPence` grows with each land sale

## Scope

### This chapter builds
- Services: `simulation-engine`
- New domain types:
  - `ColonialLabourMarket` — colony, free labourers, wage, land availability
  - `SufficientPrice` — Wakefield's artificial land price: price per acre, minimum wage-years required to save enough
  - `WageWorkerIndependence` — per-worker record: years worked, savings, landowner status
  - `SystematicColonisation` — the full scheme: colony, price, imported labourers, land fund
- New functions:
  - `ColonialLabourRegulation(market ColonialLabourMarket, scheme SystematicColonisation) ColonialLabourMarket` — apply the scheme; return the regulated market showing how many remain wage-dependent
- New HTTP endpoints:
  - `POST /v1/colonial-markets` — create a colonial labour market scenario
  - `POST /v1/colonial-markets/{id}/regulate` — apply a `SystematicColonisation` scheme; return regulated state
  - `POST /v1/colonial-markets/{id}/independence` — compute how many periods a worker needs before becoming a landowner
- Migration: `00010_ch33_colonial_markets.sql` — `colonial_markets`, `sufficient_prices` tables
- React: add a "Ch. 33 — Colonisation Theory" panel; inputs for colony name, free labourers, wage, land price; "without Wakefield scheme" vs. "with Wakefield scheme" comparison; show how raising the land price extends worker dependence; closing the full Vol. I arc: from primitive accumulation (Ch. 26) to the reproduction of capitalist relations in the colonies

### Explicitly deferred to later chapters
- Actual colonial history beyond Wakefield's theory — contextual; add as records, not new domain types
- Imperialism and monopoly capital — Book III and Lenin's extension; beyond Vol. I scope
- The "modern" colonial period post-1870 — outside Marx's frame of reference in this chapter
