---
chapter: 31
title: "Genesis of the Industrial Capitalist"
status: proposed
primary_service: simulation-engine
---

## Concepts → types

| Marx term | Go identifier | Kind | Package | Notes |
|---|---|---|---|---|
| capital origin | `CapitalOrigin` | type | `simulation` | `Source string; AmountPence Pence; Period string` — `Source` is one of: `"usury"`, `"commerce"`, `"colonial-plunder"`, `"national-debt"`, `"taxation"`, `"guild-master-accumulation"` |
| colonial transfer | `ColonialTransfer` | type | `simulation` | `From string; To string; ValuePence Pence; Method string` — looted wealth flowing from periphery to metropole |
| national debt | `NationalDebt` | type | `simulation` | `AmountPence Pence; InterestRateBps int64; CreditorClass string` — public debt as a lever of primitive accumulation |
| system of protection | `ProtectionSystem` | type | `simulation` | `TariffRateBps int64; Beneficiary string; PeriodStart string; PeriodEnd string` — state-backed capital formation |
| industrial capital genesis | `IndustrialCapitalGenesis` | type | `simulation` | `Origins []CapitalOrigin; ColonialTransfers []ColonialTransfer; NationalDebts []NationalDebt; ProtectionSystems []ProtectionSystem; TotalCapitalFormedPence Pence` — aggregate primitive accumulation producing the industrial capitalist |

## Fixtures

- **§ colonial plunder** `"discovery of gold and silver in America, the extirpation, enslavement and entombment in mines of the aboriginal population"` → `ColonialTransfer{From:"Americas", To:"England/Spain/Portugal", ValuePence:X, Method:"colonial-plunder"}`
- **§ Bank of England** `"Bank of England began with lending its money to the Government at 8%"` → `NationalDebt{AmountPence:bankFoundingCapital, InterestRateBps:800, CreditorClass:"private-bankers"}` (8% = 800 basis points)
- **§ slave trade** `"Liverpool employed in the slave-trade, in 1730, 15 ships; 1751, 53; 1760, 74; 1770, 96; 1792, 132"` → series of `ColonialTransfer` records with `Method:"slave-trade"` and growing `ValuePence`
- **§ protectionism** `"the system of protection was an artificial means of manufacturing manufacturers"` → `ProtectionSystem{TariffRateBps:highTariff, Beneficiary:"English manufacturers", PeriodStart:"17th c", PeriodEnd:"19th c"}`

## Invariants

- `genesis.TotalCapitalFormedPence == sum(origin.AmountPence for origin in genesis.Origins) + sum(transfer.ValuePence for transfer in genesis.ColonialTransfers)` — total is the sum of all sources
- `ColonialTransfer.ValuePence > 0` — always a positive transfer from periphery to metropole (wealth extraction)
- `NationalDebt.InterestRateBps > 0` — public debt always extracts interest; there is no interest-free state debt modelled here

## Scope

### This chapter builds
- Services: `simulation-engine`
- New domain types:
  - `CapitalOrigin` — source, amount, period; enumerated sources include usury, commerce, colonial plunder, national debt, taxation, guild accumulation
  - `ColonialTransfer` — from/to geography, value, method
  - `NationalDebt` — amount, interest rate in basis points, creditor class
  - `ProtectionSystem` — tariff rate in basis points, beneficiary, period
  - `IndustrialCapitalGenesis` — aggregate record of all primitive accumulation sources producing the industrial capitalist
- New HTTP endpoints:
  - `POST /v1/historical-stages/{id}/capital-origins` — add a capital origin to a stage
  - `POST /v1/historical-stages/{id}/colonial-transfers` — add a colonial transfer
  - `POST /v1/historical-stages/{id}/national-debts` — add a national debt record
  - `GET /v1/historical-stages/{id}/genesis` — retrieve full `IndustrialCapitalGenesis` summary for the stage
- Migration: `00009_ch31_capital_origins.sql` — `capital_origins`, `colonial_transfers`, `national_debts`, `protection_systems` tables
- React: complete the Ch. 26 primitive accumulation panel with a "Genesis of Industrial Capital" breakdown; pie chart showing the share of total capital formed by each source (usury, colonial plunder, national debt, etc.); total capital feeds into Ch. 23/24 simulation as starting capital stock

### Explicitly deferred to later chapters
- Credit system as an ongoing mechanism of accumulation (not just primitive) — Book III
- International finance and national debt as instrument of imperialism — beyond Vol. I scope
