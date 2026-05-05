---
chapter: 29
title: "Genesis of the Capitalist Farmer"
status: proposed
primary_service: simulation-engine
---

## Concepts → types

| Marx term | Go identifier | Kind | Package | Notes |
|---|---|---|---|---|
| tenant form | `TenantForm` | type | `simulation` | `Name string` — `"bailiff"`, `"metayer"`, `"capitalist-farmer"`; stages of the farmer's evolution |
| farm tenure | `FarmTenure` | type | `simulation` | `Form TenantForm; LeasePeriodYears int64; RentPence int64; CapitalAdvancedPence int64` — one tenure arrangement |
| money depreciation effect | `MoneyDepreciation` | type | `simulation` | `NominalRentPence int64; RealRentPence int64; DepreciationFactor float64` — how 16th-century inflation benefited the farmer: fixed nominal rent shrank in real terms |
| farming surplus | `FarmingSurplus` | type | `simulation` | `Revenue Pence; NominalRent Pence; WageCosts Pence; Profit Pence` — the capitalist farmer's appropriated surplus |
| compute real rent | `ComputeRealRent` | func | `simulation` | `func(nominalRent Pence, depreciation MoneyDepreciation) Pence` — deflate nominal by depreciation factor |

## Fixtures

- **§ 15th-century metayer** `"He advances one part of the agricultural stock, the landlord the other. The two divide the total product in proportions determined by contract"` → `FarmTenure{Form:"metayer", LeasePeriodYears:1, RentPence:0, CapitalAdvancedPence:halfOfStock}` with product split 50/50
- **§ long-lease effect** `"contracts for farms ran for a long time, often for 99 years. The progressive fall in the value of the precious metals ... lowered wages. A portion of the latter was now added to profits"` → `FarmTenure{Form:"capitalist-farmer", LeasePeriodYears:99, RentPence:fixedNominal}` with `ComputeRealRent` showing rent shrinking over the period
- **§ depreciation factor** `"continuous rise in the price of corn ... swelled the money capital of the farmer without any action on his part"` — if prices doubled, `MoneyDepreciation{DepreciationFactor:0.5}` means real rent halved while nominal held fixed
- **§ farmer surplus** Harrison's estimate: `"if he have not six or seven years rent lying by him, fifty or a hundred pounds, yet will the farmer think his gains very small"` → `FarmingSurplus.Profit` accumulated over the long lease represents the capitalist farmer's primitive accumulation

## Invariants

- `ComputeRealRent(nominal, md).Pence == int64(float64(nominal.Pence) * md.DepreciationFactor)` — deflation is multiplicative
- `FarmingSurplus.Profit == FarmingSurplus.Revenue - FarmingSurplus.NominalRent - FarmingSurplus.WageCosts` — the farmer's appropriated surplus
- `TenantForm` progression: `"bailiff"` → `"metayer"` → `"capitalist-farmer"` — the sequence is fixed by historical development; a simulation scenario must respect this ordering

## Scope

### This chapter builds
- Services: `simulation-engine`
- New domain types:
  - `TenantForm` — `string` enum: `"bailiff"`, `"metayer"`, `"capitalist-farmer"`
  - `FarmTenure` — lease terms, form, rent, capital advanced
  - `MoneyDepreciation` — captures the 16th-century silver depreciation effect on fixed rents
  - `FarmingSurplus` — revenue, rent, wages, profit; the farmer's primitive accumulation
- New functions:
  - `ComputeRealRent(nominal Pence, md MoneyDepreciation) Pence` — pure; deflates nominal rent
- New HTTP endpoints:
  - `POST /v1/historical-stages/{id}/farm-tenures` — add a farm tenure record to a stage
  - `POST /v1/farm-tenures/real-rent` — stateless; compute deflated real rent
  - `GET /v1/historical-stages/{id}/farm-tenures` — retrieve all tenures and computed surplus for the stage
- Migration: `00007_ch29_farm_tenures.sql` — `farm_tenures` table
- React: add to the Ch. 26 primitive accumulation panel; a "Genesis of the Capitalist Farmer" sub-section showing the three stages of tenure with rent levels; slider for depreciation factor showing how real rent falls as nominal is fixed; cumulative farmer surplus over a 99-year lease

### Explicitly deferred to later chapters
- Ground rent as a category — the farmer pays rent to the landlord but the theory of ground rent belongs to Book III (Vol. III of Capital)
- French equivalent (régisseur/fermier) — contextual; can be added as records with `Jurisdiction:"France"`, not new types
