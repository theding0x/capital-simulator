---
chapter: 26
title: "The Secret of Primitive Accumulation"
status: proposed
primary_service: simulation-engine
---

## Concepts → types

| Marx term | Go identifier | Kind | Package | Notes |
|---|---|---|---|---|
| primitive accumulation | `PrimitiveAccumulation` | type | `simulation` | `Period string; Method string; LabourersExpropriated int64; CapitalFormed Pence` — historical starting point of capitalist production |
| separation of producer from means | `ProducerSeparation` | type | `simulation` | `PreCapitalistWorkers int64; DisplacedWorkers int64; FreeProletarians int64` — the structural divide that founds wage labour |
| historical stage | `HistoricalStage` | type | `simulation` | `Name string; Description string; PrimitiveAccumulations []PrimitiveAccumulation` — epoch record for tracing the genesis of capital |

## Fixtures

- **§ the riddle** `"capitalistic production presupposes the pre-existence of considerable masses of capital ... in the hands of producers of commodities"` — this circularity is resolved by primitive accumulation; the simulation represents it as an initial `HistoricalStage` seeding the starting `CapitalStock` for Ch. 23 and Ch. 24 scenarios
- **§ two kinds of commodity possessors** `"owners of money, means of production, means of subsistence ... on the other hand, free labourers, sellers of their own labour power"` → `ProducerSeparation{PreCapitalistWorkers:N, DisplacedWorkers:D, FreeProletarians:D}` where displaced workers become the supply side of the labour market
- **§ England as the classic form** — `HistoricalStage{Name:"England 15th-18th c", Description:"agrarian expropriation creating free proletariat"}` is the canonical example seeded as reference data

## Invariants

- `separation.FreeProletarians == separation.DisplacedWorkers` — the dispossessed become the proletariat; the count is conserved
- `PrimitiveAccumulation.CapitalFormed > 0` — primitive accumulation always results in a positive initial capital stock, however it was obtained
- A `HistoricalStage` must precede the first `ReproductionCycle` in any simulation scenario — primitive accumulation is the logical prerequisite

## Scope

### This chapter builds
- Services: `simulation-engine`
- New domain types:
  - `PrimitiveAccumulation` — records a historical episode: period, method, workers expropriated, capital formed
  - `ProducerSeparation` — the structural pre-condition: pre-capitalist workers, displaced, free proletarians
  - `HistoricalStage` — named epoch grouping primitive accumulation episodes; seeds simulation scenarios
- New HTTP endpoints:
  - `POST /v1/historical-stages` — create a named historical stage with primitive accumulation records
  - `GET /v1/historical-stages` — list all historical stages
  - `POST /v1/historical-stages/{id}/seed-scenario` — use a historical stage to seed an accumulation scenario (links to Ch. 23/24 scenarios)
- Migration: `00004_ch26_historical_stages.sql` — `historical_stages`, `primitive_accumulations` tables
- React: add a "Ch. 26 — Primitive Accumulation" panel; form to define a historical starting point (period, method, initial capital, initial labour supply); button to seed a new reproduction scenario from it; the key insight displayed: "capital is not a thing but a social relation — these numbers are the record of expropriation"

### Explicitly deferred to later chapters
- Specific historical methods (enclosures, conquest, colonial plunder) — described in Ch. 26 by enumeration but modelled episode-by-episode in Ch. 27-31
- Colonial system, national debt, taxation, protectionism — all listed in Ch. 26 as "chief momenta of primitive accumulation" but each is the subject of its own chapter (27-31)
