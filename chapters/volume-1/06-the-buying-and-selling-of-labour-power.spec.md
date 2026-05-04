---
chapter: 06
title: "The Buying and Selling of Labour-Power"
status: proposed
primary_service: agent-service
---

## Concepts → types

| Marx term | Go identifier | Kind | Package | Notes |
|---|---|---|---|---|
| labour-power / capacity for labour | `LabourPower` | type | `agent` | Aggregate of mental and physical capabilities; what the worker sells |
| worker / labourer | `Worker` | type | `agent` | Agent who owns and sells labour-power as a commodity |
| capitalist / owner of money | `Capitalist` | type | `agent` | Agent who purchases labour-power; possesses money-capital |
| agent (base) | `Agent` | type | `agent` | Shared fields for any economic agent (ID, kind, created/updated) |
| agent ID | `AgentID` | type | `agent` | 96-bit hex ID; `NewAgentID()` constructor mirrors `commodity.NewID()` |
| agent kind | `AgentKind` | type | `agent` | Enum: `AgentKindWorker`, `AgentKindCapitalist` |
| value of labour-power | `LabourPowerValue` | type | `agent` | Value in `LabourMinutes`; determined by SNLT of means of subsistence |
| means of subsistence | `SubsistenceBasket` | type | `agent` | Set of commodities required to maintain and reproduce the worker |
| daily labour-power value | `DailyValue() LabourMinutes` | method | `agent` | Value of one day's labour-power; sum of SNLT of subsistence basket |
| selling price of labour-power | `WageMinutes LabourMinutes` | field | `agent` | Agreed daily wage expressed in labour-minutes; may equal or diverge from value |
| labour contract / sale period | `ContractDays int64` | field | `agent` | Duration for which labour-power is sold; must be > 0 and finite |
| minimum subsistence value | `MinimumValue() LabourMinutes` | method | `agent` | Floor of labour-power value; physically indispensable means only |
| labour-power as commodity | `LabourPowerOffering` | type | `agent` | A worker's labour-power offered for sale: owner, duration, asking wage |
| purchase of labour-power | `LabourPowerPurchase` | type | `agent` | Completed transaction: buyer, seller, agreed wage, contract period |
| purchase ID | `PurchaseID` | type | `agent` | 96-bit hex ID; `NewPurchaseID()` constructor |
| free labourer (double sense) | `IsFreeLabourer(w Worker) bool` | func | `agent` | True when worker owns their labour-power and lacks other commodities to sell |
| reproduction of labour-power | `ReproductionCost() LabourMinutes` | method | `agent` | Total SNLT to reproduce labour-power; equals `DailyValue()` under normal conditions |

## Fixtures

- **§1** `"the change must, therefore, take place in the commodity bought by the first act, M-C, but not in its value, for equivalents are exchanged"` → purchasing labour-power at its value (`WageMinutes == DailyValue()`) is an exchange of equivalents; surplus arises in use, not in the price

- **§2** `"By labour-power or capacity for labour is to be understood the aggregate of those mental and physical capabilities existing in a human being"` → `LabourPower` holds `CapacityMinutesPerDay LabourMinutes`; must be > 0 to be valid

- **§3** `"he must have it at his disposal, must be the untrammelled owner of his capacity for labour"` → `IsFreeLabourer` returns false when `Worker.OwnsLabourPower == false`; a worker cannot offer labour-power they do not own

- **§3** `"he were to sell it rump and stump, once for all, he would be converting himself from a free man into a slave"` → `ContractDays` must be finite and > 0; creating a `LabourPowerOffering` with `ContractDays <= 0` returns `ErrInvalidContract`

- **§4** `"the labourer ... must be obliged to offer for sale as a commodity that very labour-power"` → `IsFreeLabourer` requires `Worker.OwnsCommoditiesToSell == false`; a worker with independent commodity-wealth is not a free labourer in Marx's sense

- **§5** `"If half a day's average social labour is incorporated in three shillings, then three shillings is the price corresponding to the value of a day's labour-power"` → with `SNLTPerDay = 480` minutes (8 h) and a half-day subsistence cost, `DailyValue()` returns `240` (`LabourMinutes(240)`); `WageMinutes == 240` is the fair price

- **§5** `"The minimum limit of the value of labour-power is determined by the value of the commodities, without the daily supply of which the labourer cannot renew his vital energy"` → `MinimumValue()` equals the sum of SNLT for physically indispensable items only; it is ≤ `DailyValue()`

- **§5** `"The value of labour-power resolves itself into the value of a definite quantity of the means of subsistence"` → `DailyValue() == SubsistenceBasket.TotalSNLT()` when all items are at their values

## Invariants

- `DailyValue() == SubsistenceBasket.TotalSNLT()` — labour-power value is fully reducible to subsistence SNLT [§5]
- `MinimumValue() <= DailyValue()` — the physical minimum is a lower bound on the full daily value [§5]
- `ContractDays > 0` on any valid `LabourPowerOffering` — the sale is always for a finite period [§3]
- `WageMinutes >= 0` — a wage expressed in labour-minutes cannot be negative [§5]
- `IsFreeLabourer(w) == (w.OwnsLabourPower && !w.OwnsCommoditiesToSell)` — both conditions of "double freedom" must hold simultaneously [§3, §4]
- A `LabourPowerPurchase` with `WageMinutes == DailyValue()` is an exchange of equivalents; surplus-value is not yet realised here but only in the production process (explicitly deferred to Ch. 7) [§1]

## Scope

### This chapter builds
- Services: agent-service
- New domain types:
  - `AgentID` — 96-bit hex ID type with `NewAgentID()` constructor
  - `AgentKind` — string enum (`worker`, `capitalist`)
  - `Agent` — base struct with ID, kind, created/updated timestamps
  - `Worker` — embeds `Agent`; adds `OwnsLabourPower bool`, `OwnsCommoditiesToSell bool`, `LabourPower LabourPower`
  - `Capitalist` — embeds `Agent`; adds `MoneyCapital LabourMinutes` (capital expressed in value terms)
  - `LabourPower` — `CapacityMinutesPerDay LabourMinutes`; the use-value the worker sells
  - `SubsistenceBasket` — slice of `SubsistenceItem{Name string; SNLTMinutes LabourMinutes}`; `TotalSNLT()` method
  - `LabourPowerValue` — `DailyValue() LabourMinutes`, `MinimumValue() LabourMinutes`, computed from a `SubsistenceBasket`
  - `LabourPowerOffering` — `OwnerID AgentID`, `CapacityMinutesPerDay LabourMinutes`, `ContractDays int64`, `AskingWage LabourMinutes`
  - `PurchaseID` — 96-bit hex ID with `NewPurchaseID()` constructor
  - `LabourPowerPurchase` — `ID PurchaseID`, `SellerID AgentID`, `BuyerID AgentID`, `WageMinutes LabourMinutes`, `ContractDays int64`, `CreatedAt time.Time`
- New HTTP endpoints:
  - `POST /v1/workers` — register a new worker agent
  - `GET /v1/workers/{id}` — retrieve worker by ID
  - `POST /v1/capitalists` — register a new capitalist agent
  - `GET /v1/capitalists/{id}` — retrieve capitalist by ID
  - `POST /v1/labour-power/offerings` — worker posts labour-power for sale
  - `GET /v1/labour-power/offerings` — list available offerings
  - `POST /v1/labour-power/purchases` — capitalist purchases an offering; creates `LabourPowerPurchase`
  - `GET /v1/labour-power/purchases/{id}` — retrieve a completed purchase
- React: add AgentPanel component surfacing worker and capitalist registration, active offerings list, and purchase history; add `Agent`, `Worker`, `Capitalist`, `LabourPowerOffering`, `LabourPowerPurchase` to `types.ts` and corresponding fetch calls to `api.ts`

### Explicitly deferred to later chapters
- Surplus-value production — labour-power is consumed in the production process (Ch. 7); only the sale is modelled here
- Concrete labour-process mechanics (`LabourProcess`, `ProductOfLabour`) — Ch. 7 (valorization and the labour process)
- Working-day length and overtime — Ch. 10 (the working day)
- Reproduction of labour-power across generations / historical and moral element of subsistence basket — mentioned in the text but not quantified until Ch. 17+ (wages)
- Market price fluctuations of labour-power vs. its value — Ch. 6 establishes the value anchor; price dynamics belong to market-service (Ch. 2–3 work)
- Money-form of the wage (`WageShillings` etc.) — deferred until the money chapter work in market-service is complete
