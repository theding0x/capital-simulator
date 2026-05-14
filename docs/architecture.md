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
| Ch. 6     | ✅ Done     | Labour-power as commodity; workers, capitalists, labour-power value, wage, subsistence basket; buying and selling of labour-power | agent-service |
| Ch. 7     | ✅ Done     | Labour-process, valorization, surplus-value production   | agent-service, simulation-eng |
| Ch. 8-9   | ✅ Done     | Constant/variable capital, rate of surplus-value         | commodity-service             |
| Ch. 10    | ✅ Done     | Working-day segments (necessary/surplus), relay schedules, statutory limits, Factory Acts, overwork | agent-service |
| Ch. 11    | ✅ Done     | Rate and mass of surplus-value; SurplusValueRate, MassByRate, MassByWorkers, compensation law | simulation-engine |
| Ch. 12    | ✅ Done     | Relative surplus-value; WorkingDay, ShortenNecessaryLabour, RateOfSurplusValue, ExtraSurplusValue, ApplyProductivityToSNLT | simulation-engine |
| Ch. 13    | ✅ Done     | Co-operation; collective working-day, average social labour, scale of co-operation, minimum capital, supervision, cooperative productive power as property of capital | agent-service |
| Ch. 14    | ✅ Done     | Division of labour and manufacture; detail roles, collective labourer, two-fold origin (combination / splitting), heterogeneous vs serial form, proportional group size, scaling by integer multiples, skill hierarchy, manufacture minimum capital | agent-service |
| Ch. 15    | ✅ Done     | Machinery and modern industry; Machine (motor / transmission / tool), MachineValue, LifespanDays, DailyWearAndTear, ValueTransferredPerUnit, MaterialWear, MoralDepreciation, ProductivePower, Factory + PrimeMover, Tick loop, LabourDisplaced, IntensityFactor, capital-composition (V/C) conservation, CyclePhase enum | simulation-engine |
| Ch. 16    | ✅ Done     | Absolute and relative surplus-value; SurplusValue (Origin tag), AbsoluteSurplusValue, RelativeSurplusValue, ProlongWorkingDay, ReduceNecessaryLabour, RateSurplusValue, RateOfProfit, formal vs. real subjection, Mill critique | simulation-engine |
| Ch. 17    | ✅ Done     | Changes of magnitude in price of labour-power and surplus-value; WorkingDay × NecessaryLabour × LabourIntensity × LabourProductivity → LabourScenario / ScenarioOutcome; §1 inverse-relation law, §2 intensity-scales-value, §3 length-of-day shifts | agent-service |
| Ch. 18    | ✅ Done     | Various formulae for the rate of surplus-value; FormulaI (s/v), FormulaII (s/(s+v)), FormulaIII (unpaid/paid) | simulation-engine |
| Ch. 19    | ✅ Done     | Transformation of value of labour-power into wages; WageFormID, LabourPowerValue, Wage, WageAppearance, LabourDecomposition, WageForm; HourlyWage, Decompose (paid/unpaid split), Appearance (ideological inversion); POST /v1/wage-forms, GET /v1/wage-forms/{agentID} | agent-service |
| Ch. 20    | ✅ Done     | Time-wages; WorkingSessionID, DailyLabourPowerValue, HourlyPriceOfLabour (exact fraction + AsFloat), WorkingDayHours, OvertimeHours, OvertimeRatePence, WagePeriod, NominalWage, WorkingSession; ComputeHourlyPrice (exact rational fraction), ComputeSessionWage (integer arithmetic); POST /v1/time-wages/hourly-price, POST /v1/time-wages/sessions, GET /v1/time-wages/sessions/{id} | agent-service |
| Ch. 21-22 | Pending     | Wages (piece-wage, national differences) | agent-service |
| Ch. 23+   | Pending     | Accumulation of capital | all |

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

### Ch. 6 — what was built

`agent-service` implements Capital Vol. I, Ch. 6 — labour-power as a commodity bought and sold on the market:

- **LabourMinutes type.** `type LabourMinutes int64` — canonical value-magnitude unit for all Ch. 6 domain objects (separate from Ch. 4's pence-denominated money balances).
- **Worker.** A labourer who owns their capacity for labour (`OwnsLabourPower`) but not the means of production (`OwnsCommoditiesToSell = false`); the "double freedom" of the free labourer. `IsFreeLabourer()` encodes this condition.
- **Capitalist.** The owner of money-capital (`MoneyCapital LabourMinutes`) who appears on the market to purchase labour-power. Both `Worker` and `Capitalist` embed `LabourAgent` (base) with stable `AgentID` identifiers.
- **SubsistenceBasket / LabourPowerValue.** `SubsistenceBasket` is a slice of `SubsistenceItem` (name, SNLT in labour-minutes, essential flag). `LabourPowerValue` wraps a basket and exposes `DailyValue()`, `MinimumValue()` (essential items only), and `ReproductionCost()` (= `DailyValue()`).
- **LabourPowerOffering.** A worker posts their capacity for sale for a finite `ContractDays` at an `AskingWage`; validated to prevent perpetual or zero-duration contracts (`ErrInvalidContract`).
- **LabourPowerPurchase.** Records the act of buying labour-power: seller (Worker), buyer (Capitalist), `WageMinutes` per day, and `ContractDays`. Server validates that both parties exist in their respective roles.
- **MySQL migration.** `00005_ch06_labour_power.sql` adds four tables: `labour_workers`, `labour_capitalists`, `labour_power_offerings`, `labour_power_purchases`.
- **Endpoints.** `POST/GET /v1/workers`, `POST/GET /v1/capitalists`, `POST/GET /v1/labour-power/offerings`, `POST/GET/GET(id) /v1/labour-power/purchases`.

The React UI adds a "Ch. 06 — The Sale and Purchase of Labour-Power" panel with worker registration, capitalist registration, offering posting, and purchase recording forms, plus live lists of all entities.

### Ch. 7 — what was built

`agent-service` and `simulation-engine` model Capital Vol. I, Ch. 7 — the labour-process as the unity of the labour-process proper and the valorization process:

- **LabourProcess.** `LabourProcess`, `LabourProcessID`, `NewLabourProcessID` — one purposeful act of production tying `WorkerID`, `CapitalistID`, `MeansOfProduction`, and `Duration` together. `Validate()` enforces §1.d: zero-duration runs and missing parties are rejected.
- **MeansOfProduction.** `RawMaterial` (commodity reference + quantity + SNLT per unit) and `Instrument` (commodity reference + wear per run) — Marx's three factors of §1.c.
- **ValorizationProcess.** Wraps a `LabourProcess`; exposes `NecessaryLabour()`, `SurplusLabour()`, `SurplusValue()`, and `ProductValue()`. All seven invariants from the spec are test-covered.
- **Pure functions.** `TransferredValue(MeansOfProduction)` (constant capital), `ValueAdded(duration)` (living labour, identity for uniform skill), `SurplusLabour(wd, nl) = wd - nl`.
- **Worker extension.** `Worker` gains `LabourPowerValueMinutes` — the daily reproduction cost of labour-power, snapshotted into each `LabourProcess` record at run-time.
- **Store.** `LabourProcessStore` interface; `Memory` and `MySQL` implementations. Migration `00006` adds `labour_processes` table (means stored as JSON) and `labour_power_value_minutes` column on `labour_workers`.
- **HTTP.** `POST /v1/labour-processes` (run a process; returns product + full valorization summary), `GET /v1/labour-processes/{id}` (fetch a recorded run). Proxied through api-gateway.
- **ProductionRun.** `simulation-engine/engine.ProductionRun` type introduced; full tick scheduler deferred to Ch. 10+.
- **React UI.** Ch. 07 panel: worker/capitalist picker, means-of-production builder (raw materials + instruments), working-day duration input, valorization result card (necessary / surplus labour breakdown, rate of surplus value, product total value).

### Ch. 8–9 — what was built

`commodity-service` models Capital Vol. I, Ch. 8 (Constant and Variable Capital) and Ch. 9 (The Rate of Surplus-Value):

- **ConstantCapital.** `ConstantCapital` struct (OriginalValue, Kind, ServiceLifeDays) with `Validate()`; `ConstantKind` enum (`instrument`, `raw_material`, `auxiliary`). `WearFractionFor` returns the fraction of value transferred per production cycle (1/ServiceLifeDays for instruments, 1.0 for materials consumed in one cycle). `TransferredValue` computes the LabourMinutes transferred using `math.Round`.
- **VariableCapital.** `VariableCapital` struct (WageValue, WorkingDay) with `Validate()` and `SurplusLabourFrom()`.
- **ProductValue.** `ProductValue` (c/v/s decomposition) with `Total()`; `DecomposeProductValue` accumulates constant-capital transfers across a list of inputs plus one variable-capital input.
- **CapitalComposition.** `CapitalComposition` with `Ratio()` (c/v).
- **ProductionAccount.** `ProductionAccount` records the scalar c/v/s for a completed production run. Methods: `ValueProduct()` (v+s), `ExpandedCapital()` (c+v+s), `Rate()` (s/v). `ComputeRate` returns `ErrDivisionByZero` when v=0. `SurplusProduceRatio` = s/(v+s).
- **MySQL migration.** `00002_ch08_production_accounts.sql` adds the `production_accounts` table.
- **Endpoints (stateless, Ch.8).** `POST /v1/capital/decompose` — decomposes product value given constant and variable capital inputs. `GET /v1/capital/composition?constant_value=&variable_value=` — returns c/v ratio.
- **Endpoints (persistent, Ch.9).** `POST /v1/production-accounts` — records a production account; rejects zero variable capital. `GET /v1/production-accounts` — lists all accounts. `GET /v1/production-accounts/{id}` — fetches one. `POST /v1/rate-of-surplus-value` — stateless s/v probe.
- **React UI.** Ch. 08 panel: Decompose Capital form (dynamic constant-capital list with kind/value/service-life, variable capital block) and Capital Composition form. Ch. 09 panel: Rate of Surplus-Value probe with fixture buttons (1871 Spinning Mill, Jacob's Wheat 1815, Cotton Spinner), Record Production Account form, and list of recorded accounts.

### Ch. 13 — what was built

`agent-service` models simple co-operation from Capital Vol. I, Ch. 13:

- **Cooperation.** `Cooperation` struct (ID, Name, CapitalistID, Members, CreatedAt) — a named pool of `Worker` agents assembled by one `Capitalist` agent doing the same or same-kind work.
- **CooperationMember & Supervisor.** Each member has `WorkerID`, a `Supervisory bool` flag, and `WorkingDayMinutes`. Supervisors are wage-labourers with a directing role, not capitalists. `Supervisors()` returns the directing-authority subset.
- **CooperationSize.** `int` count of workers. `CooperationMinSize = 5` codifies the Burke footnote (§3): five workers cancel individual deviations.
- **Pure functions.** `CollectiveWorkingDay(n, d) = n × d` (§2: value addition is strictly additive). `AverageSocialLabour(collective, n) = collective / n` (§3: the social law reduces back to the per-worker average). `CollectiveProductivePower(n, d)` returns a use-value output factor ≥ 1.0 (§5: the qualitative cooperation bonus). `MinimumCapital(n, dailyWage) = n × dailyWage` (§8: capital-minimum constraint). `CooperativeOrigin()` returns `"capital"` (§13: the social productive power appears as a property of capital).
- **Store.** `CooperationStore` interface with `Memory` and `MySQL` implementations. Migration `00010_ch13_cooperations.sql` adds the `cooperations` table (members stored as JSON; indexed by `capitalist_id`).
- **HTTP endpoints.** `POST /v1/cooperations` (assemble), `GET /v1/cooperations` (list, optional `?capitalist_id=`), `GET /v1/cooperations/{id}` (fetch with computed aggregates), `POST /v1/cooperations/{id}/collective-working-day` (computed n × d plus output factor), `POST /v1/cooperations/{id}/average-social-labour` (reduces collective back to per-worker), `POST /v1/cooperations/minimum-capital` (stateless probe for n × wage).
- **api-gateway.** Reverse-proxies `/v1/cooperations` and `/v1/cooperations/{rest...}` to agent-service.
- **React UI.** Ch. 13 panel with: assembly form (pick capitalist + workers, mark some as supervisory, set working-day), ledger of assembled cooperations with computed Collective Working-Day / Average Social Labour / Output Factor / Supervisors columns, and a Minimum Capital probe with Marx's §8 fixture buttons (300 @ 6 s., 10 @ 6 s., 1,200 @ 6 s.).

### Ch. 14 — what was built

`agent-service` adds the Manufacture domain from Capital Vol. I, Ch. 14:

- **Manufacture.** `Manufacture` struct (ID, Name, CapitalistID, Form, Origin, IndividualWorkingDayMinutes, Roles, CreatedAt) — a `Cooperation` reorganised around a fixed division of detail labour, owned by one capitalist. `ManufactureForm` enumerates `"heterogeneous"` and `"serial"`; `ManufactureOrigin` enumerates `"combination"` (diverse handicrafts united) and `"splitting"` (one handicraft subdivided).
- **DetailRole, SkillLevel, SpecialisedTool.** Each role names a fractional function (`Name`, `SkillLevel`, `OutputRatePerHour`, `HeadCount`, `ToolName`). `SkillLevel` is `"skilled"` or `"unskilled"`. `PartialProduct` is a `LabourMinutes` alias deliberately without an exchange-value — detail labour produces no commodities (§4).
- **CollectiveLabourer & LabourHierarchy.** `CollectiveLabourer.TotalWorkers`, `IsParalysed(absentRoles)`, `OutputPerPeriod(periodMinutes)` (bottleneck-bounded by the slowest role-group). `LabourHierarchy` is the roles ordered skilled-first with an `Unskilled()` helper (§5).
- **Pure functions.** `ProportionalGroupSize(roles, target)` returns headcount per role so rate × headcount equals the target — Marx's "iron law" (§2); returns `ErrBottleneck` if no integer solution exists. `ManufactureProductivePower(m, period)` exceeds the simple-cooperation baseline whenever `len(Roles) > 1`. `RoleLabourPowerValue(role, apprenticeshipMinutes)` returns subsistence + apprenticeship for skilled roles, subsistence only for unskilled (§5). `ManufactureMinimumCapital(m, rawMaterialCostFactor)` always exceeds `MinimumCapital(totalWorkers, averageWage)` (§5). `ScaleManufacture(m, k)` multiplies every role headcount by `k`; `ErrInvalidMultiplier` for `k < 1` (§3).
- **Store.** `ManufactureStore` interface with `Memory` and `MySQL` implementations. Migration `00012_ch14_manufacture.sql` adds the `manufactures` and `detail_roles` tables.
- **HTTP endpoints.** `POST /v1/manufactures` (create), `GET /v1/manufactures` (list, optional `?capitalist_id=` / `?form=`), `GET /v1/manufactures/{id}` (fetch with collective-labourer / hierarchy / productive-power summary), `POST /v1/manufactures/{id}/proportional-group-size` (Marx's iron law), `POST /v1/manufactures/{id}/scale` (integer-multiple preview), `GET /v1/manufactures/{id}/collective-labourer`, `GET /v1/manufactures/{id}/minimum-capital?raw_material_cost_factor=F`.
- **api-gateway.** Reverse-proxies `/v1/manufactures` and `/v1/manufactures/{rest...}` to agent-service.
- **React UI.** Ch. 14 panel with: establish-manufacture form (capitalist, form, origin, working-day, detail-role table with skill / output-rate / headcount / tool) and §2 seed buttons (type-foundry 4/2/1, needle wire serial); manufactures table with output factor and productive power; inspector that lists the labour hierarchy, runs `ProportionalGroupSize` for a target output rate, previews integer scaling, and shows manufacture minimum capital against the cooperation baseline. Note on the panel that detail labour produces no commodities — only the collective product enters exchange.

### Ch. 15 — what was built

`simulation-engine` gains its first persistent domain — the machinery / factory loop from Capital Vol. I, Ch. 15:

- **Machine.** `Machine` struct (ID, Name, MotorMechanism, TransmittingMechanism, WorkingTool, MachineValue, LifespanDays, ProductivePower, HandLabourPerUnit, AccumulatedWear, AccumulatedDepreciation, CreatedAt) — Marx's three-part instrument of §1 (motor mechanism, transmitting mechanism, working tool). `MachineID` is a 96-bit hex identifier via `NewMachineID()`.
- **Pure value-transfer functions.** `DailyWearAndTear(m) = MachineValue / LifespanDays` (§2 "by bits"). `ValueTransferredPerUnit(m, unitsPerDay) = DailyWearAndTear(m) / unitsPerDay` (§2 inverse-proportionality). `TotalValueTransferred(m, hoursPerDay, days)` caps at MachineValue once the lifespan is exhausted; `RateOfValueReproduction(m, hoursPerDay)` doubles when daily hours double (§2 lifecycle equivalence). `LabourSavingRatio(handMinutes, machineMinutes)` returns the 180× spinning fixture. `LabourDisplaced(handLabourPerUnit, unitsProduced)` is the hand labour the machine replaces per period. `ComputeMoralDepreciation(existing, rival)` is non-negative — competition from cheaper successors never raises a machine's value (§3B).
- **Factory.** `Factory` struct (ID, Name, PrimeMover, Machines, TickCount, CreatedAt) — "an organised system of machines, to which motion is communicated by the transmitting mechanism from a central automaton" (§1c). `PrimeMover` is a `(Kind, Horsepower)` pair with kinds `steam`, `water`, `electric`, `animal`. `TotalProductivePower`, `DailyValueTransfer`, and `RunTick` return one period's aggregate use-value output, value transfer, and hand labour saved.
- **Capital composition.** `CapitalComposition{Variable, Constant}` (in pence) with `Total()`, `Ratio()`, and `Substitute(wageBillRemoved, machineCost)` — encodes the §6 invariant that mechanisation conserves total capital and shifts the variable / constant composition. `IntensityFactor` (clamped ≥ 1.0) and `EffectiveLabour(nominal, factor)` cover §3C intensification of labour.
- **Engine package extensions.** `engine.Tick` (one simulation period) and `engine.CyclePhase` enum (`prosperity`, `overproduction`, `crisis`, `stagnation`) cover the §7 industrial-cycle reference.
- **Store.** First simulation-engine persistence layer. `MachineStore` and `FactoryStore` interfaces with `Memory` and `MySQL` implementations. `AdvanceTick` is fully atomic: `SELECT … FOR UPDATE` → re-read machines → accumulate per-machine wear → bump `tick_count` → append a `factory_ticks` row, all in one transaction. Migration `00001_ch15_machinery.sql` adds `machines`, `factories`, `factory_machines`, and `factory_ticks` tables.
- **HTTP endpoints.** `POST /v1/machines` (register), `GET /v1/machines` (list with computed daily wear / per-unit wear / labour displaced), `GET /v1/machines/{id}` (fetch), `GET /v1/machines/{id}/wear` (accumulated wear / depreciation / remaining value), `POST /v1/factories` (assemble — supports both bare-ID and inline machine entries), `GET /v1/factories` (list), `GET /v1/factories/{id}` (with computed totals), `POST /v1/factories/{id}/tick` (advance one period; returns the updated factory and a `TickResult` of value transferred, units produced, and hand labour saved).
- **api-gateway.** Reverse-proxies `/v1/machines`, `/v1/machines/{rest...}`, `/v1/factories`, and `/v1/factories/{rest...}` to simulation-engine.
- **React UI.** "Ch. 15 — Machinery and Modern Industry" panel with: a machine-registration form (seed buttons for the §8A needle-machine and §2 steam-plough), a machine registry table with daily wear / per-unit wear / labour displaced columns, a factory-assembly form that lets the user pick registered machines and a prime mover, a factory floor with one-click tick advance and a "last tick" result card (units produced, value transferred, hand labour saved), and a §6 capital-composition probe that flags whether the V → C substitution conserves total capital.

### Part IV cohesion refactor (post-Ch. 15)

Following an audit of Ch. 12–15 (`docs/plans/part-iv-cohesion-review.md`), Part IV was reworked to make Marx's narrative arc visible in the wiring rather than implicit in four isolated kingdoms:

- **One `LabourMinutes` per service.** `services/simulation-engine/internal/labour/labour.go` is the canonical home of `LabourMinutes` and `ProductivityFactor` within sim-engine. The `surplus`, `production`, and `machinery` packages now declare those types as Go aliases (`type LabourMinutes = labour.LabourMinutes`) rather than re-defining them — values flow between Ch. 11, Ch. 12, and Ch. 15 functions without explicit casts.
- **Productivity bridge.** Each of Ch. 13/14/15 now emits a `productivity_factor` field in its existing JSON response (Cooperation, Manufacture, Factory). Two new endpoints close the loop on Marx's claim that simple co-operation, manufacture, and machinery are the three forms by which capital realises relative surplus-value:
  - `POST /v1/production/relative-surplus` — accepts `{working_day_total, current_lpv, productivity_factor}` and returns the shortened working day and the relative-SV delta.
  - `POST /v1/production/relative-surplus-from-productivity` — accepts `{working_day_total, current_lpv, source, source_id}` where `source` is one of `cooperation`, `manufacture`, `factory`. Simulation-engine fetches the productivity factor (HTTP-to-agent-service for cooperation/manufacture, local store for factory) and pipes it through the same calculation.
  - `services/simulation-engine/internal/productivity` houses the fetcher; the bridge handler is tested via a stub that hard-codes the factor so the test stays hermetic.
- **React bridge.** A shared `RelativeSurplusBridge` component lives in `web/src/components/` and renders inside the Ch.13 cooperation inspector, the Ch.14 manufacture inspector, and the Ch.15 factory floor. Each instance prefills the source id and factor; the user adjusts the working day and current LPV and watches the relative-SV gain.
- **`engine.Tick` wired through.** `AdvanceTick` now returns `engine.Tick`; the `factory_ticks` table is read via the new `GET /v1/factories/{id}/ticks?limit=N` endpoint; the Ch.15 React panel shows the last 10 ticks for each factory as a strip table. The orphan type from the Ch.15 spec is now the canonical wire shape.
- **`CyclePhase` deleted.** The industrial-cycle enum was orphan in the original Ch.15 spec; deferred to a later volume.
- **`IntensityFactor` wired.** Factories carry an `intensity_factor` column (migration `00003_ch15_factory_intensity.sql`); the tick path multiplies `HandLabourSaved` by it via `machinery.EffectiveLabour`, encoding the §3C intensification claim.
- **Shared React format helpers.** `minutesToHours` and `poundsFromPence` (previously byte-identical in Ch.13 and Ch.14) and Ch.15's `compactNumber` / `poundsFromLabourMinutes` all live in `web/src/format.ts` as `fmtHoursLong`, `fmtPounds`, `fmtCompact`, and `fmtLabourMinutes`. Chapter files import via aliases for minimal churn.
- **Agent-service handler tests.** `cooperation_handler_test.go` and `manufacture_handler_test.go` mirror the structure of the sim-engine handler tests (Burke five-man platoon, Caslon type-foundry 4/2/1 proportionality, iron-law headcount probe, integer scaling, manufacture-minimum vs cooperation-baseline). Brings the section to parity.
- **Agent-service handler constructor.** `httpapi.New` now takes a single `AgentStore` composite interface instead of seven positional stores. `cmd/agent-service/main.go`'s `agentStore` is a one-line alias to keep the bootstrap signature stable.

### Ch. 16 — what was built

`simulation-engine` adds the §1 synthesis of Chs. 11–15 into the `surplus` package — the two analytically distinct forms of surplus-value, treated as a single magnitude type with a tag:

- **SurplusValue + Origin.** `SurplusValue` is a `(LabourMinutes, Origin)` pair; `Origin` is the enum `"absolute" | "relative"`. The magnitude type is shared; the tag is the only thing that distinguishes the two mechanisms. `AbsoluteSurplusValue` and `RelativeSurplusValue` are zero-overhead wrappers that embed `SurplusValue`.
- **WorkingDay.** `WorkingDay{Total, NecessaryLabour, SurplusLabour}` with the partition invariant `Total == NecessaryLabour + SurplusLabour`. Fields use bare `labour.LabourMinutes` per the chapter spec.
- **ProlongWorkingDay(wd, extraMinutes).** §1 mechanism for absolute SV: raises `Total` and `SurplusLabour` by `extraMinutes`, holds `NecessaryLabour` fixed. Returns the new partition plus an `AbsoluteSurplusValue` carrying the gain.
- **ReduceNecessaryLabour(wd, factor).** §1 mechanism for relative SV: applies a `ProductivityFactor` (re-exported from the `labour` package) to `NecessaryLabour`, holds `Total` fixed. Returns the new partition plus a `RelativeSurplusValue` carrying the reduction. Guards against driving `NecessaryLabour` to zero (the wage relation must survive).
- **RateSurplusValue / RateOfProfit.** Two functions surfacing Marx's §1 correction of Mill. `RateSurplusValue = surplusLabour / necessaryLabour`; `RateOfProfit = surplusValue / totalCapitalAdvanced`. The two return different values for the same surplus magnitude whenever constant capital is positive.
- **ProductiveLabour / CollectiveLabourer / SubjectionKind.** Labelled placeholders for the §1 supporting cast: productive labour is labour that produces surplus-value for capital (not merely useful labour); the collective labourer is the agent-group the analysis treats as a unit (cooperation chapter computes its productive power); `SubjectionKind` is the `"formal" | "real"` enum distinguishing capital's formal seizure of an existing labour-process from its real revolutionisation of the technical process.
- **HTTP endpoints (stateless).** `POST /v1/surplus-value/absolute` (compute ASV by prolongation), `POST /v1/surplus-value/relative` (compute RSV by productivity), `GET /v1/surplus-value/rate?surplus_labour=…&necessary_labour=…&total_capital=…[&surplus_value=…]` (returns both rates plus a `mill_critique_holds` boolean).
- **api-gateway.** Reverse-proxies `/v1/surplus-value/absolute`, `/v1/surplus-value/relative`, and `/v1/surplus-value/rate` to simulation-engine.
- **React UI.** "Ch. 16 — Absolute and Relative Surplus-Value" panel with: a working-day configurator (Total + NecessaryLabour numeric inputs) driving a side-by-side calculator that runs both the absolute (extension) and relative (productivity) paths simultaneously; a Mill-critique sub-panel that lets the user enter `(s, v, c+v)` and surfaces the two rates plus a "Mill critique holds?" indicator.

### Ch. 18 — what was built

`simulation-engine` adds a new `simulation` package implementing Marx's three equivalent formulae for the rate of surplus-value from Capital Vol. I, Ch. 18:

- **NecessaryLabour / SurplusLabour / WorkingDay.** `NecessaryLabour{Minutes}`, `SurplusLabour{Minutes}`, and `WorkingDay{TotalMinutes}` — typed wrappers over `labour.LabourMinutes` that make the ch.18 partition explicit. `WorkingDay.TotalMinutes` is always derived as `NecessaryLabour.Minutes + SurplusLabour.Minutes` by `ComputeRates`, enforcing the partition invariant.
- **VariableCapital / SurplusValue / ValueOfProduct.** Money-pence proxies for the same labour-time ratios; defined per spec but deferred from active use — absolute pricing belongs to Chs. 2-3.
- **RateScenario / RateResult.** `RateScenario{NecessaryLabour, SurplusLabour}` is the input; `RateResult{FormulaI, FormulaII, FormulaIII float64}` is the output.
- **FormulaI(s, v) float64.** s/v — the rate of exploitation. Unbounded above 1.0 (English agricultural labourer: 540 min / 180 min = 3.0 = 300%).
- **FormulaII(s, wd) float64.** s/(s+v) — surplus labour as a fraction of the working day. Always strictly less than 1.0, because surplus-labour is always less than the full working day.
- **FormulaIII(unpaid, paid LabourMinutes) float64.** Unpaid / paid — identical in magnitude to Formula I; only the labels differ.
- **ComputeRates(RateScenario) RateResult.** Pure function; applies all three formulae in one call.
- **HTTP endpoint (stateless).** `POST /v1/surplus-value/rates` — accepts `{necessary_labour_minutes, surplus_labour_minutes}`, returns `{formula_i, formula_ii, formula_iii, working_day_minutes, ...}`.
- **api-gateway.** Reverse-proxies `/v1/surplus-value/rates` to simulation-engine.
- **React UI.** "Ch. 18 — Different Formulae for the Rate of Surplus-Value" panel with: two numeric inputs (necessary / surplus labour in minutes); three preset buttons (§ Formula I 100%, § Formula II bounded, § agricultural labourer 300%); side-by-side result cards for each formula with a note on its key property.

### Ch. 20 — what was built

`agent-service` adds the `time_wage.go` domain file implementing Marx's Ch. 20 analysis of the time-wage form.

- **Domain types.** `DailyLabourPowerValue{Pence int64}`, `HourlyPriceOfLabour{Numerator, Denominator int64}` with `AsFloat() float64`, `WorkingDayHours{Hours int64}`, `OvertimeHours{Hours int64}`, `OvertimeRatePence{Pence int64}`, `WagePeriod` ("daily" | "weekly"), `NominalWage{Pence int64}`, `WorkingSession` (persisted composite).
- **ComputeHourlyPrice(v, h).** Returns the exact rational fraction — no rounding. Encodes the invariant that lengthening the working day lowers the hourly price even when the daily value is unchanged.
- **ComputeSessionWage(s).** Integer arithmetic: `(daily_pence / hours) * hours + overtime_hours * overtime_rate`. The floor-and-multiply captures Marx's observation that the capitalist computes time-wages in whole integers.
- **Store.** `TimeWageStore` interface; `Memory` and `MySQL` implementations. Migration `00016_ch20_time_wage_sessions.sql`; seed `00017_ch20_seed.sql` with two Thomas Hobson sessions (12h standard, 15h + 3 OT at 4d).
- **HTTP endpoints.** `POST /v1/time-wages/hourly-price` (stateless probe), `POST /v1/time-wages/sessions` (create session, returns 201 + Location), `GET /v1/time-wages/sessions/{id}` (retrieve with derived hourly price and nominal wage).
- **api-gateway.** Reverse-proxies `/v1/time-wages` and `/v1/time-wages/{rest...}` to agent-service.
- **React UI.** "Ch. 20 — Time-Wages" panel: inputs for daily value, normal hours, overtime hours, overtime rate, wage period; live hourly price preview; result cards showing the exact fraction + decimal hourly price and the nominal wage breakdown (normal pay + overtime).
