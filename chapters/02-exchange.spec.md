plu---
chapter: 02
title: "Exchange"
status: proposed
primary_service: market-service
---

## Concepts → types

| Marx term | Go identifier | Kind | Package | Notes |
|---|---|---|---|---|
| Exchange / trade act | `Exchange` | struct | `market` | A completed bilateral transfer: two parties, two commodities, the quantities surrendered by each |
| Party to exchange (commodity owner) | `Owner` | struct | `market` | Represents a person qua owner; carries an ID and a name; personification of a commodity relation |
| Owner ID | `OwnerID` | named string | `market` | 96-bit hex from crypto/rand, same pattern as `commodity.ID` |
| Owner ID constructor | `NewOwnerID` | func | `market` | Returns a fresh `OwnerID` |
| Offer (one side of a trade intention) | `Offer` | struct | `market` | The commodity and quantity an owner brings to market and the commodity-kind they seek in return |
| Exchange ID | `ExchangeID` | named string | `market` | 96-bit hex from crypto/rand; identifies a completed exchange record |
| Exchange ID constructor | `NewExchangeID` | func | `market` | Returns a fresh `ExchangeID` |
| Value realised in exchange | `RealisedValue` | named int64 (`LabourMinutes`) | `market` | The value magnitude confirmed by the act of exchange; mirrors `commodity.LabourMinutes` |
| Direct barter (x use-value A = y use-value B) | `BarterRatio` | struct | `market` | Captures the quantitative proportion in pre-money direct exchange; does not yet assume a universal equivalent |
| Universal equivalent | `UniversalEquivalent` | struct | `market` | The commodity that the social act of all others has set apart to express value; precursor to the money-form |
| Social act (setting apart the equivalent) | `SetUniversalEquivalent` | func | `market` | Given a population of commodities and a chosen commodity, records it as the universal equivalent |
| Money-commodity | `MoneyCommodity` | struct | `market` | The universal equivalent once crystallised; carries the commodity-service ID of the elected commodity |
| C-M-C circuit leg | `CircuitLeg` | struct | `market` | One step in C-M-C: either C→M (sale) or M→C (purchase); tagged with `LegKind` |
| Circuit leg kind | `LegKind` | named string + consts | `market` | `KindSale` (`"sale"`) or `KindPurchase` (`"purchase"`) |
| Price (value expressed in money) | `Price` | struct | `market` | Quantity of the money-commodity that expresses the value of one unit of a given commodity |
| Price magnitude | `PriceAmount` | named int64 | `market` | Quantity of money-commodity units; integer to avoid float precision drift |

## Fixtures

- **direct barter** `"x Commodity A = y Commodity B"` → `BarterRatio{CommodityA: linecID, QtyA: 20, CommodityB: coatID, QtyB: 1}` asserts `RealisedValue` of both sides equal when SNLT is matched
- **non-use-value for owner** `"All commodities are non-use-values for their owners, and use-values for their non-owners"` → an `Offer` where `OfferedCommodityID == owner.OwnedCommodityID` and `SeeksCommodityKind != owner.OwnedCommodityID` is valid; an `Offer` where owner seeks their own commodity is rejected with `ErrOfferInvalid`
- **social act / universal equivalent** `"The social action therefore of all other commodities, sets apart the particular commodity in which they all represent their values"` → `SetUniversalEquivalent(pop, gold)` returns `UniversalEquivalent{CommodityID: goldID}` and every other commodity can compute a `BarterRatio` against it
- **money crystallises from exchange** `"Money is a crystal formed of necessity in the course of the exchanges"` → `MoneyCommodity` created from a `UniversalEquivalent` carries the same `CommodityID` and a non-zero `CreatedAt`
- **price as value-expression** `"If by reason of new or more easier mines a man can procure two ounces of silver as easily as he formerly did one, the corn will be as cheap at ten shillings the bushel as it was before at five shillings"` → doubling `SNLTPerUnit` of silver halves `Price.Amount` of corn expressed in silver (Petty, fn. 12)
- **C→M sale leg** `"commodities must be realised as values before they can be realised as use-values"` → `CircuitLeg{Kind: KindSale, CommodityID: linecID, MoneyID: goldID}` completes only when `Exchange` records equal `RealisedValue` on both sides

## Invariants

- `Exchange.RealisedValue == commodity.Value(Exchange.QtyA)` and `Exchange.RealisedValue == commodity.Value(Exchange.QtyB)` when both commodities have their canonical SNLT [direct barter §]
- `BarterRatio{QtyA: qa, QtyB: qb}` is valid iff `RealisedValue(A, qa) == RealisedValue(B, qb)` within rounding; mismatched values must return `ErrValueMismatch`
- `SetUniversalEquivalent` is idempotent: calling it twice with the same commodity returns the same `UniversalEquivalent` without error
- `Price{Amount: p, MoneyCommodityID: m}` satisfies `p * moneySNLT == commoditySNLT * unitQty`; halving money SNLT doubles `Price.Amount` [Petty fn. 12]
- `CircuitLeg` of `KindSale` followed by `KindPurchase` with equal `RealisedValue` completes one full C-M-C circuit; the net value transferred is zero [exchange §]
- An `Owner` cannot appear as both giver and receiver in the same `Exchange`; such a record must be rejected with `ErrSelfExchange`

## Scope

### This chapter builds
- Services: `market-service` (port 8083) — domain types, store layer, and HTTP endpoints for exchange, offers, the universal equivalent, and prices
- New domain types:
  - `Owner`, `OwnerID`, `NewOwnerID` — commodity-owner as economic subject
  - `Offer` — a trade intention brought to market
  - `Exchange`, `ExchangeID`, `NewExchangeID` — a completed bilateral transfer with `RealisedValue`
  - `BarterRatio` — direct x-use-value-A = y-use-value-B proportion
  - `UniversalEquivalent` — socially elected single equivalent commodity
  - `MoneyCommodity` — the crystallised universal equivalent (money-form in motion)
  - `Price`, `PriceAmount` — value expressed as a quantity of the money-commodity
  - `CircuitLeg`, `LegKind` — one step of C-M-C; `KindSale` / `KindPurchase`
  - `RealisedValue` — value magnitude confirmed by the act of exchange
- New HTTP endpoints:
  - `POST /v1/owners` — register a commodity owner
  - `GET /v1/owners` — list owners
  - `GET /v1/owners/{id}` — fetch owner
  - `POST /v1/offers` — submit an offer (owner + commodity + qty + sought kind)
  - `GET /v1/offers` — list open offers
  - `DELETE /v1/offers/{id}` — withdraw an offer
  - `POST /v1/exchanges` — record a completed bilateral exchange; validates equal `RealisedValue`
  - `GET /v1/exchanges` — list exchanges
  - `GET /v1/exchanges/{id}` — fetch exchange detail
  - `POST /v1/universal-equivalent` — social act: elect a commodity as universal equivalent
  - `GET /v1/universal-equivalent` — fetch the current universal equivalent
  - `POST /v1/money-commodity` — crystallise the universal equivalent into money
  - `GET /v1/money-commodity` — fetch the active money-commodity
  - `POST /v1/prices` — compute and store a price (commodity valued in money-commodity units)
  - `GET /v1/prices/{commodityID}` — retrieve the current price of a commodity
- React: add an Exchange panel to `App.tsx`; owner registration, offer board, exchange recorder, and a price display that shows value expressed in the money-commodity; extend `types.ts` with all new Go structs; extend `api.ts` with client calls for all new endpoints

### Explicitly deferred to later chapters
- Hoarding (money withdrawn from circulation as store of value) — Ch. 3; `market-service`
- Money as means of payment (credit, deferred settlement) — Ch. 3; `market-service`
- The general formula for capital M-C-M' and surplus-value — Ch. 4; `agent-service` + `market-service`
- The buying and selling of labour-power — Ch. 4+; `agent-service`
- Price oscillations around value (supply/demand dynamics) — deferred until `simulation-engine` drives ticks
- Redis hot-cache of price state — deferred until `simulation-engine` tick loop is wired
