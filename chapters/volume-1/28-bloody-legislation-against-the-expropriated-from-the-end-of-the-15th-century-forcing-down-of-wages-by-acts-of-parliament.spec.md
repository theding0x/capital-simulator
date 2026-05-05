---
chapter: 28
title: "Bloody Legislation Against the Expropriated, from the End of the 15th Century. Forcing Down of Wages by Acts of Parliament"
status: proposed
primary_service: simulation-engine
---

## Concepts → types

| Marx term | Go identifier | Kind | Package | Notes |
|---|---|---|---|---|
| wage statute | `WageStatute` | type | `simulation` | `Period string; Jurisdiction string; MaxWagePence int64; MinWagePence int64; EnforcementPenalty string` — legal ceiling or floor on wages |
| vagrancy law | `VagrancyLaw` | type | `simulation` | `Period string; Jurisdiction string; Punishment string; TargetPopulation string` — legal coercion converting dispossessed into wage workers |
| statutory wage | `StatutoryWage` | type | `simulation` | `ActedWagePence int64; MarketWagePence int64; Deviation float64` — comparison of legally imposed vs. market wage |
| labour discipline regime | `LabourDisciplineRegime` | type | `simulation` | `Period string; Mechanisms []string; WageStatutes []WageStatute; VagrancyLaws []VagrancyLaw` — aggregate view of state coercion in a given epoch |

## Fixtures

- **§ Statute of Labourers 1349** `"maximum of wages is dictated by the State, but on no account a minimum"` → `WageStatute{Period:"1349", Jurisdiction:"England", MaxWagePence:X, MinWagePence:0, EnforcementPenalty:"imprisonment"}`
- **§ statute 8 George II** `"forbad a higher day's wage than 2s. 7½d. for journeymen tailors in and around London"` → `WageStatute{Period:"1740s", Jurisdiction:"London tailors", MaxWagePence:31, MinWagePence:0}` (2s 7½d = 31 halfpennies)
- **§ Henry VIII vagrancy** `"whipping and imprisonment for sturdy vagabonds ... tied to the cart-tail and whipped until the blood streams ... swear an oath to go back to their birthplace"` → `VagrancyLaw{Period:"1530", Jurisdiction:"England", Punishment:"whipping+imprisonment", TargetPopulation:"sturdy vagabonds"}`
- **§ statutory vs. market** Under normal conditions after 1813 repeal: `StatutoryWage{ActedWagePence:0, MarketWagePence:marketRate, Deviation:0.0}` — market wage prevails; before repeal: deviation can be negative (statutory cap below market)

## Invariants

- `statute.MaxWagePence >= statute.MinWagePence` — legal bounds are consistent; until 1796 min was always 0
- `StatutoryWage.Deviation == float64(ActedWage-MarketWage) / float64(MarketWage)` — negative deviation means statutory cap is below market (wage suppression)
- A `LabourDisciplineRegime` without any `WageStatutes` and without any `VagrancyLaws` represents the post-1825 norm where "dull compulsion of economic relations" replaces direct legal coercion

## Scope

### This chapter builds
- Services: `simulation-engine`
- New domain types:
  - `WageStatute` — a legal maximum or minimum wage: period, jurisdiction, pence, penalty
  - `VagrancyLaw` — a law coercing the dispossessed into wage labour: period, jurisdiction, punishment, target
  - `StatutoryWage` — comparison struct: imposed vs. market wage and the deviation
  - `LabourDisciplineRegime` — epoch-level aggregate of coercive mechanisms
- New HTTP endpoints:
  - `POST /v1/historical-stages/{id}/wage-statutes` — add a wage statute to a historical stage
  - `POST /v1/historical-stages/{id}/vagrancy-laws` — add a vagrancy law
  - `GET /v1/historical-stages/{id}/labour-discipline` — retrieve the full `LabourDisciplineRegime` for the stage
  - `POST /v1/statutory-wages/compare` — stateless; compute deviation between statutory and market wage
- Migration: `00006_ch28_wage_statutes.sql` — `wage_statutes`, `vagrancy_laws` tables
- React: extend the Ch. 26 panel with a "State Coercion" sub-section; table of wage statutes and vagrancy laws per period; show how the regime transitions from legal coercion (pre-1825) to market compulsion (post-1825)

### Explicitly deferred to later chapters
- Trades Union legislation (1825 repeal, 1871 act) — mentioned but belongs to the history of labour organisation, not a domain model
- French and continental equivalents (Chapelier law 1791) — contextual; add as additional `VagrancyLaw` / `WageStatute` records, not new types
