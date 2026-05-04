---
chapter: 04
title: "The General Formula for Capital"
status: proposed
primary_service: agent-service
---

## Concepts → types

| Marx term | Go identifier | Kind | Package | Notes |
|---|---|---|---|---|
| Economic agent (bearer of class relation) | `Agent` | type | `agent` | Holds identity, class, and money balance |
| Agent identity | `ID` | type | `agent` | `type ID string`; 96-bit hex from `crypto/rand` via `NewID()` |
| Class position | `Class` | type | `agent` | `type Class string`; constants `Capitalist`, `Worker`, `Miser` |
| Money balance (pennies) | `Pence` | type | `agent` | `type Pence int64`; £100 = 10000, £110 = 11000 |
| Simple commodity circuit C—M—C | `CircuitCMC` | const | `agent` | `CircuitType = "C-M-C"` — sell to buy; ends in use-value |
| Capital circuit M—C—M′ | `CircuitMCM` | const | `agent` | `CircuitType = "M-C-M-prime"` — buy to sell dearer; ends in more money |
| Circuit discriminator | `CircuitType` | type | `agent` | `type CircuitType string`; identifies which form governs a transaction |
| Single execution of M—C—M′ | `CapitalCircuit` | type | `agent` | Records MAdvanced, CommodityID, MReturned, SurplusValue; SurplusValue computed on creation |
| Surplus-value (∆M, observed) | `SurplusValue Pence` | field | `agent` | `MReturned - MAdvanced`; positive ⇒ valorisation; source deferred to Ch. 5 |
| Money-as-capital (advancing, not spending) | `Advance` | method | `agent.Agent` | Deducts MAdvanced from MoneyBalance, returns updated Agent |
| Realise circuit (complete M—C—M′) | `Realise` | method | `agent.Agent` | Adds MReturned to MoneyBalance; surplus may be reinvested or hoarded |
| Reinvest (capitalist behaviour) | `Reinvest` | method | `agent.Agent` | Immediately advances full balance into new circuit; only valid for `Capitalist` class |
| Hoard (miser behaviour) | `Hoard` | method | `agent.Agent` | Withdraws balance from circulation; only valid for `Miser` class; returns `ErrNotCapitalist` if called on `Capitalist` |
| Agent store contract | `Store` | type | `store` | Interface: Create/Get/List/Update/Delete + ListByClass |
| In-memory store | `Memory` | type | `store` | Test implementation of `Store` |
| MySQL store | `MySQL` | type | `store` | Production implementation; `CREATE TABLE IF NOT EXISTS agents` on `NewMySQL` |
| Update payload | `Update` | type | `store` | Partial patch: optional Name, MoneyBalance fields |
| Circuit store contract | `CircuitStore` | type | `store` | Interface: CreateCircuit/GetCircuit/ListCircuits(agentID) |

## Fixtures

- **§1** `"The circulation of commodities is the starting-point of capital."` → an `Agent` with `Class = Capitalist` and `MoneyBalance = 10000` (£100) begins a `CapitalCircuit` with `MAdvanced = 10000`; the circuit is recorded and the agent's balance decremented to `0`.
- **§8** `"I purchase 2,000 lbs. of cotton for £100, and resell the 2,000 lbs. of cotton for £110"` → `CapitalCircuit{MAdvanced: 10000, MReturned: 11000}` yields `SurplusValue = 1000` (£10); `circuit.SurplusValue == circuit.MReturned - circuit.MAdvanced`.
- **§10** `"The circuit M—C—M would be absurd … if the intention were to exchange … £100 for £100"` → a `CapitalCircuit` where `MAdvanced == MReturned` has `SurplusValue == 0`; the circuit is valid to record but signals no valorisation.
- **§14** `"The miser is merely a capitalist gone mad; the capitalist is a rational miser."` → an `Agent` with `Class = Miser` calling `Hoard()` succeeds and returns balance unchanged; the same agent calling `Reinvest()` returns `ErrNotCapitalist`; an `Agent` with `Class = Capitalist` calling `Reinvest()` succeeds and calling `Hoard()` returns `ErrNotCapitalist`.
- **§15** `"M—M′, money which begets money"` → after `Realise`, agent balance becomes `MAdvanced + SurplusValue`; a second `CapitalCircuit` created immediately after uses the full new balance as `MAdvanced`, demonstrating the self-expanding circuit with no upper bound.
- **§6** `"In the circulation C—M—C … consumption … is its end and aim"` → a `CapitalCircuit` with `CircuitType = CircuitCMC` records `MAdvanced` and `MReturned` equal (net zero money gain); an `Agent` with `Class = Worker` may only initiate `CircuitCMC` circuits; attempting `CircuitMCM` returns `ErrWrongClass`.

## Invariants

- `circuit.SurplusValue == circuit.MReturned - circuit.MAdvanced` [§8] — surplus-value is always the arithmetic difference; never stored independently of those two fields.
- `agent.MoneyBalance >= 0` [§1, §3] — an agent cannot advance more money than their balance; `Advance` returns `ErrInsufficientFunds` if `MAdvanced > MoneyBalance`.
- For any completed `CapitalCircuit`, `circuit.MAdvanced > 0` [§5] — advancing zero is meaningless circulation; `CreateCircuit` rejects `MAdvanced <= 0` with a validation error.
- After `Realise(circuit)`, `agent.MoneyBalance == previousBalance + circuit.MReturned` [§15] — realisation always credits the full returned sum, including any surplus-value, to the agent's balance.
- A `Miser` agent's `MoneyBalance` never decreases via normal circuit operations [§10] — `Hoard()` is the only valid money operation for a miser and it is a no-op on balance; `Advance` and `Reinvest` return `ErrNotCapitalist`.

## Scope

### This chapter builds
- Services: `agent-service` (primary — full domain, store, HTTP); `market-service` (touched — receives circuit phase records representing the buy and sell legs, but only as a stub endpoint; full market logic deferred)
- New domain types:
  - `Agent` — bearer of a class position with a money balance; the conscious representative of capital or labour
  - `ID` — 96-bit hex agent identity, generated by `NewID()`
  - `Class` — named string type; constants `Capitalist`, `Worker`, `Miser`
  - `Pence` — named `int64` money amount in pennies; canonical unit for all money in the simulation
  - `CircuitType` — named string type; constants `CircuitCMC` and `CircuitMCM`
  - `CapitalCircuit` — record of a single M—C—M′ execution: `MAdvanced`, `CommodityID`, `MReturned`, `SurplusValue` (computed), `CircuitType`, `AgentID`, timestamps
  - `Update` — partial-patch payload for `PATCH /v1/agents/{id}`
- New HTTP endpoints:
  - `POST /v1/agents` — create a new agent (name, class, initial money balance)
  - `GET /v1/agents` — list all agents; supports `?class=` query param
  - `GET /v1/agents/{id}` — get a single agent by ID
  - `PATCH /v1/agents/{id}` — update agent name or balance
  - `DELETE /v1/agents/{id}` — remove agent
  - `POST /v1/agents/{id}/circuits` — record a new capital circuit execution (MAdvanced, CommodityID, MReturned); computes SurplusValue and updates agent balance atomically
  - `GET /v1/agents/{id}/circuits` — list all circuits for an agent
  - `POST /v1/agents/{id}/reinvest` — capitalist-only: advance full current balance into a new circuit immediately
  - `POST /v1/agents/{id}/hoard` — miser-only: withdraw balance from circulation (idempotent no-op on balance, sets a `Hoarding bool` flag)
- React: add `Agent` and `CapitalCircuit` to `web/src/types.ts` mirroring the Go structs; add a Ch. 04 panel in the chapter shell that lists agents by class, shows each agent's balance in £ (converting from Pence), and displays a table of circuits with M→C→M′ and ∆M columns

### Explicitly deferred to later chapters
- Why surplus-value exists (the source of ∆M) — Ch. 4 observes it from market outcomes only; the production of surplus-value through unpaid labour-time is explained in Ch. 5–7
- Labour-power as a distinct commodity — Ch. 6; `Worker` agents in Ch. 4 are created and hold balances but do not yet sell labour-power
- Exploitation rate / rate of surplus-value (`s/v`) — Ch. 7; requires variable capital and surplus-value broken into paid/unpaid portions
- Merchant capital and interest-bearing capital as distinct circuit variants — Ch. 4 mentions the M—M′ abridged form but does not model it separately; deferred to later parts
- Full market-service integration — market-service receives stub calls for circuit legs in Ch. 4; the order book, price discovery, and exchange mechanics come in Ch. 2–3 follow-on work
- Multi-agent interaction (one agent selling to another through a market) — Ch. 4 records circuits on a single agent; coordinated multi-agent exchange deferred until market-service is complete
- Landed property and landowner class — Part VI/VII of Capital; `Class` type reserves space but no `Landowner` constant is added yet
