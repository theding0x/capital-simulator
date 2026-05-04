# Architecture

Capital Simulator is a microservices simulation of an economy as described in Karl Marx's *Capital, Volume I*. The architecture is intentionally modular: each Marxist economic category (the commodity, the agent, the market, the production process, the circulation of capital) gets its own service so that chapter-by-chapter additions stay localized.

## Topology

```
                ┌──────────────┐
        browser │     web      │  React + Vite + TS, served by nginx
                └──────┬───────┘
                       │  /api/*
                       ▼
                ┌──────────────┐
                │ api-gateway  │  external HTTP entrypoint, fans out
                └──┬─┬─┬─┬─────┘
                   │ │ │ │
       ┌───────────┘ │ │ └────────────┐
       ▼             ▼ ▼              ▼
┌────────────┐ ┌──────────┐ ┌──────────────────┐
│ commodity- │ │  agent-  │ │ market-service   │
│  service   │ │ service  │ │                  │
└─────┬──────┘ └────┬─────┘ └────────┬─────────┘
      │             │                │
      └─────────────┼────────────────┘
                    │
            ┌───────▼────────────┐
            │ simulation-engine  │  drives ticks
            └───────┬────────────┘
                    │
   ┌────────────────┴───────────────┐
   ▼                                ▼
┌────────┐                       ┌───────┐
│ MySQL  │  durable state        │ Redis │  hot caches & tick state
└────────┘                       └───────┘
```

## Services

| Service              | Port  | Marxist role                                                            | Persistence       |
|----------------------|-------|-------------------------------------------------------------------------|-------------------|
| `api-gateway`        | 8080  | External entrypoint; fans out to domain services.                       | —                 |
| `commodity-service`  | 8081  | Use-value, exchange-value, value (Ch. 1).                               | MySQL             |
| `agent-service`      | 8082  | Workers, capitalists, and other class-bearers (Ch. 4+).                 | MySQL             |
| `market-service`     | 8083  | Exchange and circulation; C-M-C and M-C-M' (Ch. 2-3).                   | MySQL + Redis     |
| `simulation-engine`  | 8084  | Time-step orchestrator; advances the economy one period at a time.      | MySQL + Redis     |

All Go services share a single root `go.mod` (`github.com/theding0x/capital-simulator`). Cross-cutting concerns live under `pkg/`:

- `pkg/log` — structured logging via `log/slog`.
- `pkg/httpx` — HTTP server scaffolding with `/healthz`, `/readyz`, and graceful shutdown.
- `pkg/mysql` — MySQL driver, connection config, and `Migrate` helper (goose v3). **Live as of Ch. 1 (refactored from MongoDB); migrations added Ch. 3.**
- `pkg/redis` — Redis connection config (driver wired in by a later chapter).

## Data flow per simulation tick (target shape)

1. `simulation-engine` advances the clock and tells `agent-service` to act.
2. Agents form intentions to produce / sell / buy and notify `commodity-service` and `market-service`.
3. `market-service` matches trades and updates prices.
4. State changes are persisted to MySQL; hot tick state is cached in Redis.
5. `api-gateway` exposes a read view of the world to the React UI.

This is the *target*; the initial scaffold ships health endpoints only.

## Roadmap (chapter-driven)

Each chapter of *Capital* turns into a feature branch and PR. Approximate mapping:

| Chapter   | Status      | Concepts                                                 | Primary services              |
|-----------|-------------|----------------------------------------------------------|-------------------------------|
| Ch. 1     | ✅ Done     | Commodity, use-value, value, exchange-value, value-forms, fetishism | commodity-service |
| Ch. 2     | ✅ Done     | Exchange, owners, offers, barter ratio, universal equivalent, money-form, prices | market-service |
| Ch. 3     | ✅ Done     | Money: hoarding, means of payment, world money            | market-service                |
| Ch. 4     | ✅ Done     | Money → capital; M—C—M′, class positions, surplus-value  | agent-service                 |
| Ch. 5     | ✅ Done     | Contradictions in the general formula; value conservation proof | agent-service |
| Ch. 6-7   | Pending     | Labour-process, valorization, surplus-value              | agent-service, simulation-eng |
| Ch. 8-9   | Pending     | Constant/variable capital, rate of surplus-value         | commodity, simulation-eng     |
| Ch. 10    | Pending     | The working day                                          | agent-service, simulation-eng |
| Ch. 11+   | Pending     | Cooperation, machinery, wages, accumulation              | all                           |

### Ch. 1 — what was built

`commodity-service` now models all four sections of Capital Vol. I, Ch. 1:

- **§1 Two factors of a commodity.** `Commodity`, `UseValue`, `LabourMinutes`, `ProductivityChange`, with the inverse-proportionality law between productivity and value enforced by tests.
- **§2 Dual character of labour.** Each `Commodity` carries a `ConcreteLabour`; `AsAbstractLabour` makes the reduction to homogeneous human labour explicit at every value-computation site.
- **§3 The form of value.** `SimpleFormOf`, `ExpandedFormOf`, `GeneralFormOf`, `MoneyFormOf` — each derivable as a view over a population of commodities, with a chosen money-commodity for the money-form.
- **§4 The fetishism of commodities.** `SocialRelationsOf` and the `/v1/commodities/{id}/social-relations` endpoint surface the labour relations that exchange-value normally hides.

MySQL is wired up via `pkg/mysql`; the `commodities` table has a unique case-insensitive index on `name` via `utf8mb4_unicode_ci` collation. The api-gateway reverse-proxies `/v1/commodities/*` and `/v1/exchange-ratio` to commodity-service, and the React dashboard exposes full CRUD plus a "Reveal" toggle that renders the fetishism critique inline.

### Ch. 2 — what was built

`market-service` models the exchange process from Capital Vol. I, Ch. 2:

- **Owners.** `Owner`, `OwnerID`, `NewOwnerID` — the commodity owner as economic subject; "guardians" who bring commodities to market. CRUD endpoints at `/v1/owners`.
- **Offers.** `Offer`, `OfferID` — a trade intention (owner + commodity + quantity + seeks-kind). Validated so an owner cannot seek the same commodity they offer (`ErrOfferInvalid`). Endpoints: `POST/GET /v1/offers`, `DELETE /v1/offers/{id}`.
- **Exchanges.** `Exchange`, `ExchangeID`, `RealisedValue` — a completed bilateral transfer. Enforces `ErrSelfExchange`; records the labour-time value confirmed by the act. Endpoints: `POST/GET /v1/exchanges`, `GET /v1/exchanges/{id}`.
- **BarterRatio.** `BarterRatio` — the direct proportion x use-value A = y use-value B, validated structurally in the domain layer.
- **Universal equivalent.** `UniversalEquivalent` — the commodity set apart by social act. Idempotent: calling `SetUniversalEquivalent` with the same commodity ID is a no-op. Endpoint: `POST/GET /v1/universal-equivalent`.
- **Money-commodity.** `MoneyCommodity` — the crystallised universal equivalent. "Money is a crystal formed of necessity in the course of the exchanges." Endpoint: `POST/GET /v1/money-commodity`.
- **Prices.** `Price`, `PriceAmount`, `ComputePrice` — value expressed as a quantity of the money-commodity. Encodes the Petty law (fn. 12): halving money SNLT doubles the price. Endpoints: `POST /v1/prices`, `GET /v1/prices`, `GET /v1/prices/{commodityID}`.

The api-gateway reverse-proxies all market routes to market-service. The React UI adds a "Ch. 02 — Exchange" section with owner registration, offer board, exchange recorder, universal-equivalent election, money crystallisation, and price computation.

### Ch. 3 — what was built

`market-service` extends Ch. 2 to model the three functions of money from Capital Vol. I, Ch. 3:

- **Circuits C—M—C.** `Circuit`, `CircuitLeg` — commodity sold for money, money spent on a different commodity; value realised and spent.
- **Hoarding.** `Hoard` — the miser's withdrawal of money from circulation; tracks the gold quantity withheld.
- **Means of payment / credit.** `PaymentObligation` — a deferred payment (creditor, debtor, due date); settled via `POST /v1/payment-obligations/{id}/settle`.
- **World money.** `WorldMoneyTransfer` — cross-border gold movement; records sender, receiver, and gold milligrams.
- **Quantity of money in circulation.** `GET /v1/circulation/money-required?sum_of_prices=&velocity=` — implements Marx's formula M = (ΣP) / V.

The React UI adds a "Ch. 03 — Money" section with circuit recorder, hoard panel, payment-obligation tracker, and world-money transfer log.

### Ch. 4 — what was built

`agent-service` replaces the placeholder stub to model Capital Vol. I, Ch. 4 — the general formula for capital:

- **Agent.** `Agent`, `ID`, `Class` (`Capitalist`/`Worker`/`Miser`), `Pence` — economic subjects with class positions and money balances (in pence). CRUD endpoints at `/v1/agents` with optional `?class=` filter.
- **CapitalCircuit.** `CapitalCircuit`, `CircuitType` (`C-M-C` / `M-C-M-prime`) — a recorded movement of money through commodities. `SurplusValue = MReturned − MAdvanced` is enforced as a domain invariant. Endpoints: `POST/GET /v1/agents/{id}/circuits`.
- **Reinvest.** `POST /v1/agents/{id}/reinvest` — Capitalist-only: the full current balance becomes M-advanced in a new circuit; balance is updated atomically in the same MySQL transaction.
- **Hoard.** `POST /v1/agents/{id}/hoard` — Miser-only: sets the `hoarding` flag; the agent withdraws from circulation.
- **Class restrictions.** Workers are barred from M—C—M′ circuits (`ErrWrongClass`); Misers and Workers cannot Reinvest or Advance (`ErrNotCapitalist`); only Misers can Hoard.
- **MySQL store.** `CreateCircuit` is fully atomic: `SELECT money_balance FOR UPDATE` → check ≥ MAdvanced → INSERT circuit → `UPDATE money_balance += SurplusValue`, all inside one transaction.

The api-gateway reverse-proxies `/v1/agents` and `/v1/agents/{rest...}` to agent-service. The React UI adds a "Ch. 04 — The General Formula for Capital" panel with agent creation, per-class grouping, £-denominated balances, circuit history table (Type / M / C / M′ / ∆M), and Hoard action for Misers.

### Ch. 5 — what was built

`agent-service` extends Ch. 4 to prove the central contradiction of Capital Vol. I, Ch. 5 — that circulation alone cannot be the source of surplus-value:

- **Owner class.** `Owner` added to `Class` enum — a simple commodity owner who buys and sells equivalents but cannot originate surplus-value (cannot do M—C—M′ circuits).
- **LabourMinutes.** `labour_minutes int64` field added to `Agent`, tracking the abstract-labour magnitude the agent commands; persisted via migration `00004_ch05_labour_minutes.sql`.
- **ExchangeResult.** Bilateral exchange outcome with before/after values for both parties and `origin` tag. `TotalBefore() == TotalAfter()` is the invariant.
- **MerchantsCapital.** M-C-M' operating purely in circulation; `Origin()` returns `"redistribution"` if `SurplusValue != 0`.
- **UsurersCapital.** Degenerate M-M' circuit (no commodity); the source of surplus cannot be located within the circuit. `Origin()` follows the same logic as `MerchantsCapital`.
- **ExchangeEquivalents / ExchangeNonEquivalents / TotalValue.** Pure functions proving value conservation for any bilateral exchange.
- **POST /v1/circuits** — stateless circuit probe; returns `surplus_value` and `origin` tag.
- **POST /v1/exchange-simulations** — stateless exchange simulator; returns full `ExchangeResult` with conservation proof.

The React UI adds a "Ch. 05 — Contradictions in the General Formula" panel with a circuit probe form and an exchange simulation table.
