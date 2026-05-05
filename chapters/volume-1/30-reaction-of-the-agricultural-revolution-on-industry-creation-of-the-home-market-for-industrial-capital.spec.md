---
chapter: 30
title: "Reaction of the Agricultural Revolution on Industry. Creation of the Home-Market for Industrial Capital"
status: proposed
primary_service: simulation-engine
---

## Concepts → types

| Marx term | Go identifier | Kind | Package | Notes |
|---|---|---|---|---|
| home market formation | `HomeMarketFormation` | type | `simulation` | `Period string; LabourersProletarianised int64; DomesticIndustryDestroyed bool; NewWageLabourers int64; MarketSizePence Pence` — how agrarian expropriation creates an internal market |
| domestic industry | `DomesticIndustry` | type | `simulation` | `Name string; HouseholdsEngaged int64; AnnualOutputPence Pence` — pre-capitalist household production (spinning, weaving, etc.) |
| manufacture réunie | `ManufactureReunie` | type | `simulation` | `Workers int64; OutputPence Pence` — concentrated workshop replacing scattered domestic producers |
| manufacture separee | `ManufactureSeparee` | type | `simulation` | `IndependentProducers int64; OutputPence Pence` — discrete scattered workshops (Mirabeau's preferred form) |
| home market size | `HomeMarketSize` | type | `simulation` | `CommoditisedOutputPence Pence; WageLabourers int64` — size of the internal market created by proletarianisation |
| compute market formation | `ComputeMarketFormation` | func | `simulation` | `func(di DomesticIndustry, expropriated int64) HomeMarketFormation` — converts domestic producers into wage workers and subsistence goods into commodities |

## Fixtures

- **§ Westphalian flax** `"a part of the Westphalian peasants, who at the time of Frederick II all span flax, forcibly expropriated ... At the same time arise large establishments for flax-spinning and weaving"` → `DomesticIndustry{Name:"Westphalian flax spinning", HouseholdsEngaged:N}` → `ComputeMarketFormation(di, expropriated)` returns `HomeMarketFormation` where output is now sold as commodities rather than self-consumed
- **§ Mirabeau's distinction** `"grand manufactories ... hundreds of men under a director"` vs. `"discrete workshop ... no one will become rich, but many labourers will be comfortable"` → `ManufactureReunie{Workers:200}` vs. `ManufactureSeparee{IndependentProducers:200}` with equal output but different surplus distribution
- **§ peasant → commodity buyer** `"peasant family produced the means of subsistence and the raw materials ... transformed into commodities"` → `HomeMarketSize.CommoditisedOutputPence` equals what was previously self-produced; new wage workers must now buy what they once made

## Invariants

- `formation.NewWageLabourers == formation.LabourersProletarianised` — the proletarianised are the new wage workers
- `formation.MarketSizePence >= prevDomesticOutput` — commoditisation of formerly self-produced goods expands the market at least as large as domestic output was
- `ManufactureReunie.OutputPence > ManufactureSeparee.OutputPence` when workers == producers — concentration enables greater output; Mirabeau's preference for discrete shops is a political not an economic argument

## Scope

### This chapter builds
- Services: `simulation-engine`
- New domain types:
  - `DomesticIndustry` — `Name string; HouseholdsEngaged int64; AnnualOutputPence Pence`; pre-capitalist household production
  - `ManufactureReunie` — concentrated workshop: workers, output
  - `ManufactureSeparee` — scattered workshop: independent producers, output
  - `HomeMarketFormation` — the quantitative result of proletarianisation: new wage workers, market size
  - `HomeMarketSize` — commoditised output and wage-labourer count
- New functions:
  - `ComputeMarketFormation(di DomesticIndustry, expropriated int64) HomeMarketFormation` — pure; converts domestic producers to wage workers
- New HTTP endpoints:
  - `POST /v1/historical-stages/{id}/domestic-industries` — register a domestic industry for a stage
  - `POST /v1/market-formation` — stateless; compute home market formation from domestic industry + expropriation count
  - `GET /v1/historical-stages/{id}/home-market` — retrieve aggregate home market size
- Migration: `00008_ch30_domestic_industries.sql` — `domestic_industries` table
- React: add to the Ch. 26 primitive accumulation panel; a "Home Market Formation" sub-section; inputs for domestic industry (name, households, output); expropriation count from Ch. 27; display the resulting market size and wage-labour supply; contrast `ManufactureReunie` vs. `ManufactureSeparee` side-by-side

### Explicitly deferred to later chapters
- Modern industry's complete destruction of rural domestic industry — the chapter notes that manufacture only partially achieves this; complete separation belongs to the machinery and large industry chapter (Ch. 15)
- Factory Acts as limit on industrial capital's expansion — covered in Ch. 10 and Ch. 15
