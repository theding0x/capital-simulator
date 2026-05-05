---
chapter: 03
title: "Money, Or the Circulation of Commodities"
status: implemented
primary_service: market-service
---

## Concepts → types

| Marx term | Go identifier | Kind | Package | Notes |
|---|---|---|---|---|
| C—M—C circuit | `Circuit` | type | `market` | Two-leg commodity → money → commodity movement; `sale_leg` + `purchase_leg` |
| Circuit leg (sale / purchase) | `CircuitLeg` | type | `market` | One step: `KindSale` (C→M) or `KindPurchase` (M→C) |
| Leg kind | `LegKind` | type | `market` | `"sale"` / `"purchase"` string enum |
| Hoarding | `Hoard` | type | `market` | Withdrawal of gold from circulation; tracks `owner_id` + `amount` |
| Means of payment / credit | `PaymentObligation` | type | `market` | Deferred payment: `creditor_id`, `debtor_id`, `amount`, `due_at`, `paid_at` |
| World money | `WorldMoneyTransfer` | type | `market` | Cross-border gold movement; records `sender_id`, `receiver_id`, `gold_mg` |
| Quantity of money in circulation | `MoneyRequired` (query) | func | `httpapi` | `GET /v1/circulation/money-required` implements M = ΣP / V |
| Measure of value (§1) | `MoneyCommodity` | type | `market` | Gold as standard of price; carried forward from Ch.02 |
| Price | `Price`, `PriceAmount`, `ComputePrice` | type/func | `market` | Commodity value expressed in money-commodity units (Petty law) |
| Universal equivalent | `UniversalEquivalent` | type | `market` | Socially elected money-form; carried forward from Ch.02 |

## Fixtures

- **§2** `"one quarter of wheat — £2 — 20 yards of linen — £2 — 1 Bible — £2 — 4 gallons of brandy — £2"` → `Circuit` with `sale_leg` (linen → £2) and `purchase_leg` (£2 → Bible); `Circuit.sale_leg.value == Circuit.purchase_leg.value`

- **§2** `"the quantity of money functioning as the circulating medium is equal to the sum of the prices of the commodities divided by the number of moves made by coins of the same denomination"` → `GET /v1/circulation/money-required?sum_of_prices=8&velocity=4` returns `money_required = 2`

- **§3A** `"In order that gold may be held as money, and made to form a hoard, it must be prevented from circulating"` → `POST /v1/hoards` with `owner_id` + `amount`; a `Hoard` record represents gold withdrawn from the `Circuit` flow

- **§3B** `"The vendor becomes a creditor, the purchaser becomes a debtor"` → `PaymentObligation{creditor_id, debtor_id, amount, due_at}` created via `POST /v1/payment-obligations`; settled (i.e. `paid_at` set) via `POST /v1/payment-obligations/{id}/settle`

- **§3C** `"Money of the world serves as the universal medium of payment, as the universal means of purchasing, and as the universally recognised embodiment of all wealth"` → `WorldMoneyTransfer{sender_id, receiver_id, gold_mg}` recorded via `POST /v1/world-money-transfers`

## Invariants

- `Circuit.sale_leg.value == Circuit.purchase_leg.value` — C—M—C is a circuit of equivalents; the same value-magnitude leaves in commodity-form and returns in commodity-form [§2]
- `money_required = sum_of_prices / velocity` (integer division, `velocity > 0`) — Marx's circulation law; `ComputeMoneyRequired` must reject `velocity <= 0` [§2]
- `PaymentObligation.paid_at == nil` until `settle` is called; once settled the obligation is closed and cannot be re-settled [§3B]

## Scope

### This chapter builds
- Services: market-service
- New domain types:
  - `Circuit` — two-leg C—M—C movement with `sale_leg` and `purchase_leg` of type `CircuitLeg`
  - `CircuitLeg` — one step in the circuit: `LegKind`, `CommodityID`, `MoneyID`, `OwnerID`, `Value`
  - `LegKind` — `"sale"` | `"purchase"` string constant pair (`KindSale`, `KindPurchase`)
  - `Hoard` — miser's gold withdrawal: `HoardID`, `OwnerID`, `Amount`, `CreatedAt`
  - `PaymentObligation` — deferred payment: `ObligationID`, `CreditorID`, `DebtorID`, `Amount`, `DueAt`, `PaidAt`
  - `WorldMoneyTransfer` — cross-border bullion: `TransferID`, `SenderID`, `ReceiverID`, `GoldMg`, `CreatedAt`
- New HTTP endpoints:
  - `POST /v1/circuits` — record a C—M—C circuit
  - `GET /v1/circuits` — list circuits
  - `POST /v1/hoards` — hoarder withdraws gold from circulation
  - `GET /v1/hoards` — list hoards
  - `POST /v1/payment-obligations` — creditor extends deferred payment to debtor
  - `GET /v1/payment-obligations` — list obligations
  - `POST /v1/payment-obligations/{id}/settle` — debtor pays; sets `paid_at`
  - `POST /v1/world-money-transfers` — record cross-border gold movement
  - `GET /v1/world-money-transfers` — list transfers
  - `GET /v1/circulation/money-required` — query param `sum_of_prices` + `velocity`; returns `money_required`
- React: "Ch. 03 — Money" panel with circuit recorder, hoard panel, payment-obligation tracker (with settle action), world-money transfer log, and money-required calculator
- Migration: `00002_ch03_money.sql` — adds `circuits`, `hoards`, `payment_obligations`, `world_money_transfers` tables (the existing migration in the repo currently only contains `market_config` and `prices`; Ch.03 store methods are defined in the frontend types but not yet in the Go `Store` interface or `Memory`/`MySQL` implementations)

### Explicitly deferred to later chapters
- M—C—M′ (capital circuit producing surplus-value) — Ch.04; `CircuitType` and `SurplusValue` belong to `agent-service`, not `market-service`
- Credit-money / bank-notes — Ch.03 mentions these but Marx explicitly defers full credit analysis: "Money based upon credit implies conditions, which, from our standpoint of the simple circulation of commodities, are as yet totally unknown to us"
- Coin as token / paper money — §2C covers symbols of value; simulation defers to a later chapter that models state-issued currency
- Gateway proxying of Ch.03 routes (`/v1/circuits`, `/v1/hoards`, `/v1/payment-obligations`, `/v1/world-money-transfers`, `/v1/circulation`) — not yet wired in `api-gateway/cmd/api-gateway/main.go`
