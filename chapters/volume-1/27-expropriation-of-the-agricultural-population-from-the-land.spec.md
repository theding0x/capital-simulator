---
chapter: 27
title: "Expropriation of the Agricultural Population from the Land"
status: proposed
primary_service: simulation-engine
---

## Concepts → types

| Marx term | Go identifier | Kind | Package | Notes |
|---|---|---|---|---|
| expropriation episode | `ExpropriationEpisode` | type | `simulation` | `Period string; Method string; AcresSeized int64; PeasantsDispossessed int64; NewProletarians int64` — one act of agricultural expropriation |
| enclosure | `Enclosure` | type | `simulation` | `CommonAcres int64; PrivateAcres int64; DisplacedFamilies int64` — conversion of common land to private property |
| agrarian transition | `AgrianTransition` | type | `simulation` | `From string; To string; Episodes []ExpropriationEpisode; TotalProletarians int64` — series of episodes producing a free labour force |

## Fixtures

- **§ Duchess of Sutherland** `"15,000 inhabitants ... 3,000 families ... systematically hunted and rooted out ... 794,000 acres ... 6,000 acres on the sea-shore — 2 acres per family"` → `ExpropriationEpisode{Period:"1814-1820", Method:"Highland clearing", AcresSeized:794000, PeasantsDispossessed:15000, NewProletarians:15000}`
- **§ enclosure statistics** `"between 1801 and 1831 ... 3,511,770 acres of common land ... presented to the landlords"` → `Enclosure{CommonAcres:3511770, PrivateAcres:3511770, DisplacedFamilies:3511770/avgFamilyHolding}`
- **§ sheep displacing men** `"24 farms ... melted into three farms"` → `AgrianTransition` records farm consolidation as a proxy for displacement

## Invariants

- `episode.NewProletarians <= episode.PeasantsDispossessed` — dispossession is the upper bound on new proletarians; some may emigrate or die
- `enclosure.PrivateAcres == enclosure.CommonAcres` — what was common becomes private; no acres created or destroyed
- `transition.TotalProletarians == sum(episode.NewProletarians for episode in transition.Episodes)` — cumulative expropriation

## Scope

### This chapter builds
- Services: `simulation-engine`
- New domain types:
  - `ExpropriationEpisode` — one historical act: period, method, scale, outcome in proletarians
  - `Enclosure` — the parliamentary/legal form of expropriation: common → private acres, families displaced
  - `AgrianTransition` — a named series of episodes; feeds into `HistoricalStage` from Ch. 26
- New HTTP endpoints:
  - `POST /v1/historical-stages/{id}/expropriation-episodes` — add an episode to a stage
  - `GET /v1/historical-stages/{id}/expropriation-episodes` — list episodes and aggregate totals
  - `POST /v1/enclosures` — record an enclosure act
- Migration: `00005_ch27_expropriation_episodes.sql` — `expropriation_episodes`, `enclosures` tables
- React: extend the Ch. 26 primitive accumulation panel with a sub-section for agrarian episodes; table of episodes with acres seized and peasants dispossessed; running total of proletarians created; final count feeds into Ch. 23/24 simulation as initial labour supply

### Explicitly deferred to later chapters
- Reformation and church lands (mentioned as cause) — contextual history; no new domain type needed
- Bloody legislation against vagabonds created by expropriation — modelled in Ch. 28
- Sheep-walks, deer forests: specific uses of enclosed land — contextual; land use type is a string field, not a domain model
