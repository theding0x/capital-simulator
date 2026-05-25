# Architecture

Capital Simulator is a microservices simulation of an economy as described in Karl Marx's *Capital* — all three volumes. The architecture is intentionally modular: each Marxist economic category (the commodity, the agent, the market, the production process, the circulation of capital, the distribution of profit) gets its own service so that chapter-by-chapter additions stay localized. The application models capital as **value in motion** — the circuit **M—C(Lp+Mp)…P…C'—M'**. Vol. I zooms into **P** (production), Vol. II into the circulation phases and turnover, Vol. III into the totality and the distribution of surplus-value.

Textual references for Vol. I use the Moore–Aveling English translation (1887), as digitised and hosted by the [Marxists Internet Archive](https://www.marxists.org/archive/marx/works/1867-c1/). Vol. II (1885) and Vol. III (1894) texts live in the [red-vault Obsidian vault](../CLAUDE.md#what-this-is) mirrored from the same source. See the [README's Sources and acknowledgements section](../README.md#sources-and-acknowledgements) for credit details.

## Topology

```
                ┌──────────────┐
        browser │     web      │  React + Vite + TS, served by nginx
                └──────┬───────┘
                       │  /api/*
                       ▼
                ┌──────────────┐
                │ api-gateway  │  external HTTP entrypoint, fans out
                └──┬─┬─┬─┬─┬───┘
                   │ │ │ │ │
       ┌───────────┘ │ │ │ └────────────────────┐
       │   ┌─────────┘ │ └────────────────┐     │
       ▼   ▼           ▼                  ▼     ▼
┌────────────┐ ┌──────────┐ ┌──────────────────┐ ┌────────────────┐
│ commodity- │ │  agent-  │ │ market-service   │ │  finance-      │
│  service   │ │ service  │ │                  │ │   service      │
└─────┬──────┘ └────┬─────┘ └────────┬─────────┘ └────────┬───────┘
      │             │                │                    │
      └─────────────┼────────────────┘                    │
                    │                                     │
            ┌───────▼────────────┐                        │
            │ simulation-engine  │  drives ticks          │
            └───────┬────────────┘                        │
                    │              ◀────────── reads ─────┘
   ┌────────────────┴───────────────┐
   ▼                                ▼
┌────────┐                       ┌───────┐
│ MySQL  │  durable state        │ Redis │  hot caches & tick state
└────────┘                       └───────┘
```

## Services

| Service              | Port  | Marxist role                                                                             | Persistence       |
|----------------------|-------|------------------------------------------------------------------------------------------|-------------------|
| `api-gateway`        | 8080  | External entrypoint; fans out to domain services.                                        | —                 |
| `commodity-service`  | 8081  | Vol. I — commodity, value, value-forms, c+v decomposition.                               | MySQL             |
| `agent-service`      | 8082  | Vol. I — workers, capitalists, labour-process, wages, cooperation, manufacture.          | MySQL             |
| `market-service`     | 8083  | Vol. I — exchange, money, prices; Vol. II — circulation phases.                          | MySQL + Redis     |
| `simulation-engine`  | 8084  | Vol. I — production tick, machinery, reproduction; Vol. II — turnover; Vol. III — avg-rate-of-profit. | MySQL + Redis     |
| `finance-service`    | 8085  | Vol. III — profit, prices of production, rent, interest, credit, fictitious capital, the trinity formula. | MySQL             |

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

Each chapter of *Capital* turns into a feature branch and PR. Branches are named `volume-X/chapter-Y` where X ∈ {1, 2, 3}. Approximate mapping per volume:

### Volume I Roadmap — The Process of Production of Capital

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
| Ch. 21    | ✅ Done     | Piece-wages; PieceWageID, PieceWage, PieceSession, QualityOutcome, SubContractID, SubContract; ComputePiecePrice (farthings integer arithmetic), ComputePieceValue, ComputePieceEarnings, SubContractSpread; POST/GET /v1/agents/{id}/piece-wages, POST /v1/piece-price, POST/GET /v1/sub-contracts[/{id}] | agent-service |
| Ch. 22    | ✅ Done     | National differences in wages; NationalIntensity, DayWage, StandardisedWage, RelativeLabourPrice, SpindleRatio; StandardiseWage, ComputeRelativePrice; POST/GET /v1/intensities, POST /v1/wages, GET /v1/wages/{country}/standardised, GET /v1/comparisons | agent-service |
| Ch. 23    | ✅ Done     | Simple reproduction; CapitalStock, SurplusValueFund, ReproductionCycle, RepaymentPeriod | simulation-engine |
| Ch. 24    | ✅ Done     | Accumulation of capital; AdditionalCapital, Accumulation, SplitSurplus, RunExtendedReproduction | simulation-engine |
| Ch. 25    | ✅ Done     | General law of capitalist accumulation; ValueComposition, OrganicComposition, LabourDemand, IndustrialReserveArmy, GeneralLawScenario, RunGeneralLaw | simulation-engine |
| Ch. 26    | ✅ Done     | Secret of primitive accumulation; HistoricalStage, PrimitiveAccumulation, ProducerSeparation, SeparationFromStage, SeedCapitalStock | simulation-engine |
| Ch. 27    | ✅ Done     | Expropriation of agricultural population; EnclosureEvent, displacement timeline, key English enclosure waves | simulation-engine |
| Ch. 28    | ✅ Done     | Bloody legislation; WageStatute, VagrancyLaw, StatutoryWage, LabourDisciplineRegime, statutory-vs-market wage comparison | simulation-engine |
| Ch. 29    | ✅ Done     | Genesis of the capitalist farmer; TenantForm (bailiff/metayer/capitalist-farmer), FarmTenure, MoneyDepreciation, FarmingSurplus; ComputeRealRent, ComputeFarmingSurplus; currency depreciation calculator | simulation-engine |
| Ch. 30    | ✅ Done     | Reaction of the agricultural revolution on industry; DomesticIndustry, MarketFormation, HomeMarketSize; ComputeHomeMarketSize, ComputeMarketFormation; Mirabeau réunie/séparée comparison | simulation-engine |
| Ch. 31    | ✅ Done     | Genesis of the industrial capitalist; CapitalOrigin, ColonialTransfer, NationalDebt, ProtectionSystem, IndustrialCapitalGenesis; ComputeGenesis; POST capital-origins/colonial-transfers/national-debts, GET genesis | simulation-engine |
| Ch. 32    | ✅ Done     | Historical tendency of capitalist accumulation; PettyProperty, CapitalistPrivateProperty, CentralisationStep, Negation, AccumulationTrajectory; NegationOfNegation, RunCentralisation; POST /v1/accumulation/centralisation, GET negation-of-negation, persisted trajectories | simulation-engine |
| Ch. 33    | ✅ Done     | Modern theory of colonisation; ColonialLabourMarket, SufficientPrice, WageWorkerIndependence, SystematicColonisation; ColonialLabourRegulation; POST /v1/colonial-markets, regulate, independence | simulation-engine |

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

### Ch. 21 — what was built

`agent-service` adds the `piece_wage.go` domain file implementing Marx's Ch. 21 analysis of the piece-wage form.

- **Domain types.** `PieceWageID`, `PieceWage{AgentID, PricePence, NormalOutput int64}` (persisted wage contract); `QualityOutcome` ("accepted" | "rejected"); `PieceSession{PiecesProduced, QualityOutcome}` (transient session input); `SubContractID`, `SubContract{HeadLabourerID, AssistantIDs []AgentID, PieceRatePence, AssistantRatePence int64}` (sweating system).
- **ComputePiecePrice(dailyWage, normalOutput).** Integer division in farthings (¼d.): `ComputePiecePrice(144, 24) == 6`, `ComputePiecePrice(144, 48) == 3`. Invariant: `price × normalOutput == dailyWage`.
- **ComputePieceValue(dayValueProduct, normalOutput).** Analogous to piece price but over the day's full value product; always exceeds piece price because `dayValueProduct > dailyWage`.
- **ComputePieceEarnings(s, pw).** Returns `piecesProduced × pricePence` if accepted, 0 if rejected. Quality control is enforced by the wage form itself.
- **SubContractSpread(sc).** Returns `PieceRatePence - AssistantRatePence` — the per-piece gain the middleman retains.
- **Store.** `PieceWageStore` interface (one contract per agent, `ErrAlreadyExists` on duplicate); `Memory` and `MySQL` implementations. Migration `00018_ch21_piece_wages.sql`; seed `00019_ch21_seed.sql` with Thomas Hobson fixture (6 farthings/piece, 24 pieces/day) and sweating sub-contract.
- **HTTP endpoints.** `POST/GET /v1/agents/{id}/piece-wages` (register/retrieve contract); `POST /v1/piece-price` (stateless probe returning piece price, piece value, actual earnings); `POST/GET /v1/sub-contracts[/{id}]` (create/retrieve sub-contract with spread).
- **React UI.** "Ch. 21 — Piece-Wages" panel: inputs for daily wage and day value product (in farthings), normal output, pieces produced, and quality outcome; live piece price/value preview; result cards showing the price vs. value gap (the surplus made visible) and actual session earnings.

### Ch. 22 — what was built

`agent-service` adds the national wage domain from Capital Vol. I, Ch. 22:

- **Domain types.** `CountryCode` (ISO 3166-1 alpha-2 string), `NationalIntensity{CountryCode, Factor float64}` (intensity relative to international average), `DayWage{CountryCode, NominalPence, WorkingDayMinutes int64}` (domestic wage), `StandardisedWage{CountryCode, Amount int64}` (wage reduced to uniform working day), `RelativeLabourPrice{CountryCode, Ratio float64}` (wage as fraction of value produced), `SpindleRatio{CountryCode, SpindlesPerWorker int64}` (Redgrave's 1866 productivity proxy), `WageComparison` (aggregate output).
- **Pure functions.** `StandardiseWage(w, referenceDayMinutes)` — integer arithmetic, `Amount = NominalPence * ref / WorkingDayMinutes`. `ComputeRelativePrice(w, ni)` — `Ratio = 1.0 / Factor`, encoding Marx's inversion: England's higher nominal wage paradoxically implies a lower relative labour price to the capitalist because of higher intensity (Cowell 1833; Redgrave 1866).
- **Store.** `NationalWageStore` interface; `Memory` and `MySQL` implementations. `national_intensities` uses `UPSERT` (country records are updated, not duplicated). `day_wages` is one-per-country (`ErrAlreadyExists` on duplicate). `spindle_ratios` is seed-only read via `ListSpindleRatios`.
- **Migration.** `00020_ch22_national_wages.sql` adds `national_intensities`, `day_wages`, and `spindle_ratios` tables with `country_code CHAR(2)` PKs. `00021_ch22_seed.sql` inserts GB/FR/DE/RU/AT/BE/CH intensity records, GB/FR/DE day wages, and all seven Redgrave spindle entries. Down deletes only the seeded rows.
- **HTTP endpoints.** `POST /v1/intensities` (upsert; idempotent), `GET /v1/intensities` (list), `POST /v1/wages` (register; 201 + body), `GET /v1/wages/{country}/standardised?reference_day_minutes=600` (standardised wage for one country), `GET /v1/comparisons[?reference_day_minutes=600]` (full `WageComparison` with standardised wages and relative prices computed on the fly).
- **api-gateway.** Reverse-proxies `/v1/intensities`, `/v1/wages`, `/v1/comparisons` (and `/{rest...}` forms) to agent-service.
- **React UI.** "Ch. 22 — National Wages" panel: table of countries with nominal wage, standardised wage (600-min ref), relative price, and spindle count; GB paradox highlighted (highest nominal, lowest relative price); register-intensity and register-day-wage forms.

### Ch. 23 — what was built

`simulation-engine` adds the simple-reproduction domain from Capital Vol. I, Ch. 23:

- **CapitalStock.** `CapitalStock{ConstantCapital, VariableCapital Pence}` — value composition of capital at a point in time.
- **SurplusValueFund.** `SurplusValueFund{Total, Revenue, Accumulated Pence}` — the split of surplus-value produced in one cycle. Under simple reproduction `Revenue == Total` and `Accumulated == 0` every period.
- **ReproductionCycle.** `ReproductionCycle{Period int64; Capital CapitalStock; SurplusRate float64; Fund SurplusValueFund}` — one circuit of production. Capital stock does not grow across cycles.
- **OriginalCapital, IndividualConsumption, ProductiveConsumption.** `Pence`-typed wrappers distinguishing the historically given advance, personal capitalist expenditure, and productive consumption of means of production and labour-power.
- **RunSimpleReproduction(initial CapitalStock, surplusRate float64, periods int64) []ReproductionCycle.** Pure function producing a cycle-by-cycle snapshot. Each cycle's surplus-value is `math.Round(v × s′)`; under simple reproduction the fund has `Accumulated == 0` in every period.
- **RepaymentPeriod(capital, annualRevenue Pence) int64.** Integer ceiling division: the number of periods before the total consumed surplus-value equals or exceeds the original capital advanced. Encodes Marx's §repayment law: `RepaymentPeriod(1000, 200) == 5`; `RepaymentPeriod(1000, 100) == 10`.
- **Migration.** `00004_ch23_reproduction_cycles.sql` adds the `reproduction_cycles` table. `00005_ch23_seed.sql` inserts Marx's core §example (£800c + £200v, s′=100%, 5 periods; seed ID `5eed000000000000002301`) so the dashboard comes up populated.
- **HTTP endpoints (stateless).** `POST /v1/reproductions/simple` — runs the simulation, returns the cycle array and pre-computed repayment period. `POST /v1/reproductions/repayment-period` — stateless ceiling-division probe.
- **api-gateway.** Reverse-proxies `/v1/reproductions/simple` and `/v1/reproductions/repayment-period` to simulation-engine.
- **React UI.** "Ch. 23 — Simple Reproduction" panel: inputs for c, v, s′, period slider (1–20); preset button for Marx's £1,000/£200 example; result cards (original capital, annual surplus-value, repayment period); cycle table with running total of consumed surplus-value vs. original capital; rows highlighting green once the repayment period is reached.

### Ch. 24 — what was built

`simulation-engine` adds the extended-reproduction and accumulation domain from Capital Vol. I, Ch. 24:

- **AdditionalCapital.** `AdditionalCapital{Constant, Variable Pence}` — the portion of surplus-value converted into new means of production and new labour-power.
- **Accumulation.** `Accumulation{Period int64; CapitalStock CapitalStock; SurplusProduced Pence; NewConstant, NewVariable, Revenue Pence}` — one period's record of the extended-reproduction genealogy.
- **AccumulationRate / CompositionRatio.** Named `float64` types for the fraction of surplus-value reinvested (0.0 = simple reproduction; 1.0 = full accumulation) and the fraction of new capital that is constant capital (the organic composition).
- **SplitSurplus(surplus Pence, accumRate, compositionRatio float64) AdditionalCapital.** Pure function. Divides surplus-value per the accumulation rate and organic composition: `SplitSurplus(2000, 1.0, 0.8) == {1600, 400}`; `SplitSurplus(2000, 0.5, 0.8) == {800, 200}` with Revenue=1000.
- **RunExtendedReproduction(initial CapitalStock, surplusRate, accumRate, compositionRatio float64, periods int64) []Accumulation.** Traces the genealogy of accumulated capital. Each period's CapitalStock is the AdditionalCapital from the prior step — the "Abraham begat Isaac" sequence: `RunExtendedReproduction({8000,2000}, 1.0, 1.0, 0.8, 3)` produces surpluses of £2,000 → £400 → £80.
- **Migration.** `00006_ch24_accumulation_scenarios.sql` adds the `accumulation_scenarios` table. `00007_ch24_seed.sql` inserts Marx's spinner §1 (full accumulation, £8,000c + £2,000v) and §3 (partial accumulation, rate=0.5) as named scenarios with seed IDs `5eed000000000000002401`/`5eed000000000000002402`.
- **HTTP endpoints (stateless).** `POST /v1/reproductions/extended` — runs the genealogy simulation, returns the cycle array. `POST /v1/reproductions/split-surplus` — stateless; divides a single surplus-value into new constant capital, new variable capital, and revenue.
- **api-gateway.** Reverse-proxies `/v1/reproductions/extended` and `/v1/reproductions/split-surplus` to simulation-engine.
- **React UI.** "Ch. 24 — Accumulation of Capital" panel: extended-reproduction form (c, v, s′, accumulation rate, composition ratio, period slider); preset buttons for spinner §1 and partial-accumulation §3; cycle table showing the genealogy across periods; stateless split-surplus calculator with result cards for new constant capital, new variable capital, and capitalist revenue.

### Ch. 25 — what was built

`simulation-engine` adds the general-law domain from Capital Vol. I, Ch. 25:

- **ValueComposition / OrganicComposition.** `ValueComposition{ConstantCapital, VariableCapital Pence}` — the technical composition at a point in time. `OrganicComposition{Ratio float64}` — c/(c+v), always in [0, 1). `ComputeOrganicComposition(vc) OrganicComposition` — pure function.
- **LabourDemand.** `LabourDemand{Workers int64}`. `ComputeLabourDemand(totalCapital, oc, wagePence) LabourDemand` — variable-capital share divided by wage: `(totalCapital × (1−oc.Ratio)) / wagePence`.
- **IndustrialReserveArmy.** `IndustrialReserveArmy{Size int64; RelativeProportion float64}`. `ComputeReserveArmy(supply, demand) IndustrialReserveArmy` — size = max(supply − demanded, 0); proportion = size / demanded.
- **GeneralLawScenario / RunGeneralLaw.** `GeneralLawScenario` persists the initial conditions; `RunGeneralLaw(s) []GeneralLawSnapshot` simulates accumulation over `s.Periods` periods. Each period: record OC, compute workers absorbed and reserve army, then accumulate surplus via `SplitSurplus`, and advance OC by `productivityGrowth` for the next period.
- **Migration.** `00008_ch25_general_law.sql` adds `general_law_scenarios`. `00009_ch25_seed.sql` inserts §1 (unchanged OC, `productivity_growth=0.0`) and §2 (rising OC, `productivity_growth=0.05`) with seed IDs `5eed000000000000002501`/`5eed000000000000002502`.
- **HTTP endpoints.** `POST /v1/accumulation/organic-composition` (stateless OC calculation); `POST /v1/accumulation/labour-demand` (stateless workers absorbed); `POST /v1/accumulation/reserve-army` (stateless reserve army); `POST /v1/accumulation/scenarios` (persist scenario + run, returns 201 + Location); `GET /v1/accumulation/scenarios/{id}` (retrieve with series).
- **api-gateway.** Reverse-proxies `/v1/accumulation` and `/v1/accumulation/{rest...}` to simulation-engine.
- **React UI.** "Ch. 25 — General Law" panel: organic composition calculator; labour demand probe; reserve army calculator; multi-period scenario simulator with §1/§2 preset buttons; series table showing C, V, OC, workers, reserve army and relative proportion across periods.

### Part VII API prefix decision

Marx's Part VII ("The Accumulation of Capital") covers Ch. 23–25 as a single conceptual arc. Our HTTP surface currently uses two prefixes:

- `POST /v1/reproductions/simple` and `/v1/reproductions/repayment-period` (Ch. 23)
- `POST /v1/reproductions/extended` and `/v1/reproductions/split-surplus` (Ch. 24)
- `POST /v1/accumulation/organic-composition`, `/labour-demand`, `/reserve-army`, `/scenarios`, `GET /v1/accumulation/scenarios/{id}` (Ch. 25)

**Decision:** The canonical Part VII prefix is `/v1/accumulation/*`. Simple reproduction (Ch. 23) is the zero-accumulation limit case and belongs under that tree. The `/v1/reproductions/*` routes are legacy; they will be unified under `/v1/accumulation/*` when Vol. II reproduction-schemes work makes a coordinated route migration worthwhile. Until then the split is kept to avoid breaking the React UI without a migration plan.

### Ch. 26 — what was built

`simulation-engine` adds the historical-genesis layer from Capital Vol. I, Ch. 26 — the layer the rest of Part VII has been treating as given.

- **PrimitiveAccumulation.** `PrimitiveAccumulation{Period, Method string; LabourersExpropriated int64; CapitalFormed Pence}` — one historical episode. `Validate()` enforces `CapitalFormed > 0` (primitive accumulation always results in positive capital).
- **ProducerSeparation.** `ProducerSeparation{PreCapitalistWorkers, DisplacedWorkers, FreeProletarians int64}` — the conservation invariant `FreeProletarians == DisplacedWorkers` enforced by `Validate()`. Every dispossessed producer reappears on the labour market as a free wage-labourer.
- **HistoricalStage.** `HistoricalStage{ID, Name, Description, PrimitiveAccumulations, CreatedAt}` — a named epoch (canonically "England 15th–18th c.") grouping episodes. Helpers: `TotalCapitalFormed()`, `TotalLabourersExpropriated()`.
- **SeparationFromStage(h, preCapitalistWorkers) ProducerSeparation.** Projects a stage onto a separation, clamping `DisplacedWorkers` to the population and setting `FreeProletarians` equal to it.
- **SeedCapitalStock(h, compositionRatio) CapitalStock.** Splits the stage's aggregate `CapitalFormed` into constant and variable capital so that `c + v == total` — the bridge into the Ch. 23/24 reproduction loop.
- **Migration.** `00010_ch26_historical_stages.sql` adds `historical_stages` (case-insensitive unique name) and `primitive_accumulations` (FK CASCADE on stage_id). `00011_ch26_seed.sql` inserts the canonical England 15th–18th c. fixture: enclosure of common lands (£80k), expropriation of church property (£30k), colonial plunder (£90k). Seed IDs `5eed00000000000000260[1-4]`.
- **HTTP endpoints.** `POST /v1/historical-stages` (create stage + episodes in a single transaction, returns 201 + Location; 409 on duplicate name); `GET /v1/historical-stages` (list); `POST /v1/historical-stages/{id}/seed-scenario` (derive starting `CapitalStock` and `ProducerSeparation`, run `RunSimpleReproduction` over them, return cycles inline).
- **api-gateway.** Reverse-proxies `/v1/historical-stages` and `/v1/historical-stages/{rest...}` to simulation-engine.
- **React UI.** "Ch. 26 — Primitive Accumulation" panel: stages table with totals; create-stage form with editable episode rows and an England 15th–18th c. preset; seed-scenario form (pre-capitalist workers, organic composition, surplus rate, periods) that derives the starting capital, projects the separation, and runs a simple-reproduction series. The panel makes explicit the chapter's thesis: *capital is not a thing but a social relation — these numbers are the record of expropriation*.

### Ch. 29 — what was built

`simulation-engine` adds the farm-tenure domain from Capital Vol. I, Ch. 29 — tracing how the capitalist farmer emerged from the metayer system as currency depreciation and long leases transformed nominal rent burdens into windfall profits.

- **TenantForm.** String enum: `"bailiff"` (manages for absentee landlord), `"metayer"` (shares stock and product), `"capitalist-farmer"` (advances own capital, pays money rent).
- **FarmTenure.** `FarmTenure{ID, HistoricalStageID, Form, LeasePeriodYears, RentPence, CapitalAdvancedPence, RevenuePence, WageCostsPence, CreatedAt}` — one tenure arrangement. `Validate()` enforces non-empty `HistoricalStageID` and `Form`, `LeasePeriodYears ≥ 1`, all monetary fields `≥ 0`.
- **MoneyDepreciation.** `MoneyDepreciation{NominalRentPence, RealRentPence, DepreciationFactor float64}` — tracks how currency fall reduces the real rent burden. `ComputeRealRent(nominal, md)` = `Pence(float64(nominal) * md.DepreciationFactor)`.
- **FarmingSurplus.** `FarmingSurplus{Revenue, NominalRent, WageCosts, Profit Pence}` — `ComputeFarmingSurplus` returns `Profit = Revenue − NominalRent − WageCosts`. This is Harrison's "£50 or £100" accumulated per long lease.
- **Migration.** `00016_ch29_farm_tenures.sql` adds the `farm_tenures` table. `00017_ch29_seed.sql` inserts a 15th-c. metayer (`5eed000000000000002901`, 1-yr lease, revenue=1000, wages=200) and a 99-yr capitalist lease (`5eed000000000000002902`, rent=1200, revenue=3000, wages=600), both linked to the seeded England 15th–18th c. stage.
- **HTTP endpoints.** `POST /v1/historical-stages/{id}/farm-tenures` (record a tenure; returns 201 + Location); `GET /v1/historical-stages/{id}/farm-tenures` (list tenures for a stage, inline surplus computed); `POST /v1/farm-tenures/real-rent` (stateless — compute real rent given nominal rent and depreciation factor).
- **api-gateway.** Adds `/v1/farm-tenures` and `/v1/farm-tenures/{rest...}` proxy rules (historical-stages sub-paths already covered by the Ch.26 rules).
- **React UI.** Ch. 29 panel: stage selector; add-tenure form with 15th-c. Metayer and 99-yr Capitalist Lease presets; tenure records table with inline profit (green if positive) and cumulative surplus footer; currency depreciation calculator (nominal rent input + depreciation factor slider with live preview; API call returns nominal → real comparison).

### Ch. 28 — what was built

`simulation-engine` adds the state-coercion domain from Capital Vol. I, Ch. 28 — the extra-economic force bridging expropriation and the mature wage-labour relation.

- **WageStatute.** `WageStatute{ID, HistoricalStageID, Period, Jurisdiction, MaxWagePence, MinWagePence, EnforcementPenalty, CreatedAt}` — a legal ceiling or floor on wages. `Validate()` enforces `MaxWagePence ≥ MinWagePence ≥ 0`, non-empty `Period`, `Jurisdiction`, and `EnforcementPenalty`.
- **VagrancyLaw.** `VagrancyLaw{ID, HistoricalStageID, Period, Jurisdiction, Punishment, TargetPopulation, CreatedAt}` — a statute converting the dispossessed into compulsory wage-labourers.
- **StatutoryWage.** Pure struct; `ComputeStatutoryWage(acted, market int64)` computes `Deviation = (acted − market) / market`. Negative deviation = statutory cap suppresses wages below market.
- **LabourDisciplineRegime.** Computed aggregate; `AssembleRegime(stagePeriod, wageStatutes, vagrancyLaws)` deduplicates enforcement mechanisms across all records. Empty regime = post-1825 market compulsion.
- **Migration.** `00014_ch28_wage_statutes.sql` adds `wage_statutes` and `vagrancy_laws` tables. `00015_ch28_seed.sql` seeds the Statute of Labourers 1349, George II tailors statute 1740s, Henry VIII vagrancy act 1530, and 39 Elizabeth 1597 — all linked to the seeded England 15th–18th c. historical stage.
- **HTTP endpoints.** `POST /v1/historical-stages/{id}/wage-statutes`; `POST /v1/historical-stages/{id}/vagrancy-laws`; `GET /v1/historical-stages/{id}/labour-discipline` (returns assembled regime); `POST /v1/statutory-wages/compare` (stateless).
- **api-gateway.** Adds `/v1/statutory-wages` and `/v1/statutory-wages/{rest...}` proxy rules (historical-stages sub-paths already covered).
- **React UI.** Ch. 28 panel: stage selector; add-wage-statute form with 1349/1740s presets; add-vagrancy-law form with 1530/1597 presets; labour discipline regime table; statutory-vs-market deviation calculator.

### Ch. 27 — what was built

`simulation-engine` adds the enclosure-event domain from Capital Vol. I, Ch. 27 — the paradigm case of how capital expropriates the agricultural population from the land.

- **EnclosureEvent.** `EnclosureEvent{ID, Period, AcresEnclosed, PopulationDisplaced, Beneficiary, CreatedAt}` — one wave of agrarian expropriation. `Validate()` enforces `AcresEnclosed > 0` (land must actually be seized), `PopulationDisplaced ≥ 0`, and non-empty `Period` and `Beneficiary`.
- **Migration.** `00012_ch27_enclosure_events.sql` adds the `enclosure_events` table. `00013_ch27_seed.sql` inserts three canonical English waves: 15th-c. gentry enclosures (500,000 acres, 40,000 displaced), Parliamentary enclosures 1760–1820 (3,000,000 acres, 150,000 displaced), and the Sutherland clearances 1814–1820 (794,000 acres, 15,000 displaced — the specific figure Marx cites in his footnote).
- **HTTP endpoints.** `POST /v1/enclosure-events` (record a new wave; returns 201 + Location); `GET /v1/enclosure-events` (list all waves ordered by creation time).
- **api-gateway.** Reverse-proxies `/v1/enclosure-events` and `/v1/enclosure-events/{rest...}` to simulation-engine.
- **React UI.** "Ch. 27 — Expropriation of the Agricultural Population" panel: preset buttons for the three canonical English waves; record form (period, acres, population, beneficiary); displacement timeline table with running totals of acres seized and population expelled; summary note connecting each expelled producer to their reappearance as a "free" wage-labourer.

### Ch. 20 — what was built

`agent-service` adds the `time_wage.go` domain file implementing Marx's Ch. 20 analysis of the time-wage form.

- **Domain types.** `DailyLabourPowerValue{Pence int64}`, `HourlyPriceOfLabour{Numerator, Denominator int64}` with `AsFloat() float64`, `WorkingDayHours{Hours int64}`, `OvertimeHours{Hours int64}`, `OvertimeRatePence{Pence int64}`, `WagePeriod` ("daily" | "weekly"), `NominalWage{Pence int64}`, `WorkingSession` (persisted composite).
- **ComputeHourlyPrice(v, h).** Returns the exact rational fraction — no rounding. Encodes the invariant that lengthening the working day lowers the hourly price even when the daily value is unchanged.
- **ComputeSessionWage(s).** Integer arithmetic: `(daily_pence / hours) * hours + overtime_hours * overtime_rate`. The floor-and-multiply captures Marx's observation that the capitalist computes time-wages in whole integers.
- **Store.** `TimeWageStore` interface; `Memory` and `MySQL` implementations. Migration `00016_ch20_time_wage_sessions.sql`; seed `00017_ch20_seed.sql` with two Thomas Hobson sessions (12h standard, 15h + 3 OT at 4d).
- **HTTP endpoints.** `POST /v1/time-wages/hourly-price` (stateless probe), `POST /v1/time-wages/sessions` (create session, returns 201 + Location), `GET /v1/time-wages/sessions/{id}` (retrieve with derived hourly price and nominal wage).
- **api-gateway.** Reverse-proxies `/v1/time-wages` and `/v1/time-wages/{rest...}` to agent-service.
- **React UI.** "Ch. 20 — Time-Wages" panel: inputs for daily value, normal hours, overtime hours, overtime rate, wage period; live hourly price preview; result cards showing the exact fraction + decimal hourly price and the nominal wage breakdown (normal pay + overtime).

### Ch. 31 — what was built

`simulation-engine` adds the `industrial_capitalist.go` domain file implementing Marx's Ch. 31 analysis of the genesis of the industrial capitalist through primitive accumulation.

- **Domain types.** `CapitalOriginID`, `CapitalOrigin{Source, AmountPence, Period}`, `ColonialTransferID`, `ColonialTransfer{From, To, ValuePence, Method}`, `NationalDebtID`, `NationalDebt{AmountPence, InterestRateBps, CreditorClass}`, `ProtectionSystemID`, `ProtectionSystem{TariffRateBps, Beneficiary, PeriodStart, PeriodEnd}`, `IndustrialCapitalGenesis`.
- **ComputeGenesis(stageID, origins, transfers, debts, systems).** Assembles all four slices and computes `TotalCapitalFormedPence = sum(Origins.AmountPence) + sum(ColonialTransfers.ValuePence)`. National debts and protection systems are structural levers; they are returned for display but excluded from the capital total.
- **Store.** Four new store interfaces (`CapitalOriginStore`, `ColonialTransferStore`, `NationalDebtStore`, `ProtectionSystemStore`) with `Memory` and `MySQL` implementations. `ProtectionSystemStore` is read-only (seed-only per spec). Migrations `00020_ch31_capital_origins.sql` and `00021_ch31_seed.sql`: Liverpool slave trade series (1730/1751/1760/1770/1792), Bank of England founding (£1.2m at 8%), colonial plunder from the Americas, English manufacturers protection system (30% tariff, 17th–19th c.).
- **HTTP endpoints.** `POST /v1/historical-stages/{id}/capital-origins` (201), `POST /v1/historical-stages/{id}/colonial-transfers` (201), `POST /v1/historical-stages/{id}/national-debts` (201), `GET /v1/historical-stages/{id}/genesis` (200 — full genesis summary). Protection systems are seed-only; no POST endpoint.
- **React UI.** "Ch. 31 — Genesis of the Industrial Capitalist" panel: stage selector with live genesis summary (total capital formed, counts per category, seeded protection systems); forms to register capital origins, record colonial transfers, and record national debts; results lists with currency formatting.

### Ch. 32 — what was built

`simulation-engine` adds `historical_tendency.go` to model the dialectical telos of Part VIII: capitalist private property begetting its own negation through centralisation.

- **Domain types.** `PettyProperty{Producers, CapitalPerProducerPence}`, `CapitalistPrivateProperty{Firms, TotalCapitalPence, WageLabourers}`, `CentralisationStep{StepIndex, FirmsAbsorbed, CapitalConcentratedPence}`, `Negation{Stage, Description}` (stage ∈ `petty-property` | `capitalist-expropriation` | `socialised-property`), `AccumulationTrajectory{ID, Name, InitialFirms, InitialCapitalPence, Steps, FinalFirms, FinalCapitalPence, ReserveArmySize, CreatedAt}`.
- **NegationOfNegation(stages []Negation) (Negation, error).** Validates that the three stages appear in the canonical order and returns the third (socialised property). `DialecticalSequence()` returns the canonical three-stage list with glosses drawn from Ch. 32.
- **RunCentralisation(initial, steps, absorptionRate, capitalGrowthRate).** Each step absorbs ⌊firms · absorptionRate⌋ firms; capital optionally grows by `capitalGrowthRate` per step. `CapitalConcentratedPence` per step is apportioned against the *initial* capital pool so the invariant Σ(concentrated) ≤ initial ≤ final holds — concentration of ownership redistributes value rather than destroying it. The simulation floors at one absorption per step and never reaches `FinalFirms == 0`; Marx's "knell of capitalist private property" is treated as an asymptote.
- **Store.** New `AccumulationTrajectoryStore` interface with `Memory` (case-insensitive name uniqueness mirrors the MySQL constraint) and `MySQL` (transactional header + steps insert) implementations. Migrations `00022_ch32_accumulation_trajectories.sql` (header + steps tables, FK with `ON DELETE CASCADE`) and `00023_ch32_seed.sql`: Lancashire cotton 1820–1880 (1000 mills → 217), English banking 1810–1900 (650 country banks → 110 joint-stocks), American railroads 1860–1900 (420 lines → 70 trunks).
- **HTTP endpoints.** `POST /v1/accumulation/centralisation` (200 stateless or 201 when `persist: true` with a name); `GET /v1/accumulation/negation-of-negation` (200 — three dialectical stages with the resolved third negation); `GET /v1/accumulation/trajectories` (200 — list, returns `[]` not `null` when empty); `GET /v1/accumulation/trajectories/{id}` (200 / 404). api-gateway already proxies `/v1/accumulation/{rest...}` — no new proxy rules needed.
- **React UI.** "Ch. 32 — The Historical Tendency of Capitalist Accumulation" panel: dialectical-stage card (petty → capitalist → socialised), seeded trajectory selector with absorption-step table (firms remaining bar + concentration intensity bar), and a stateless centralisation simulator form (firms, capital, labourers, steps, absorption rate, capital growth) that renders the resulting trajectory inline.

### Ch. 33 — what was built

`simulation-engine` adds `colonisation.go` to read Wakefield's colonial theory against the grain: in a colony with free land the wage relation cannot reproduce itself; capital meets a producer who is still owner of his conditions of labour, and Mr. Peel's £50,000 plus 300 servants dissolve at Swan River. Wakefield's "sufficient price" is the state's artificial floor on virgin land that postpones the labourer's exit from the wage market and thereby manufactures wage-dependence by other means.

- **Domain types.** `ColonialLabourMarket{ID, Colony, FreeLabourers, AnnualWagePence, LandAvailable, WakefieldSchemeApplied, IndependenceYears, SurplusLabourExtractable, CreatedAt}`, `SufficientPrice{PricePerAcrePence, DesiredAcres}` with `Ransom() = PricePerAcrePence * DesiredAcres`, `WageWorkerIndependence{WorkerID, AnnualSavingsPence, YearsWorked, SavingsPence, TargetRansomPence, BecameLandowner}`, `SystematicColonisation{Colony, SufficientPrice, LandSalesCompleted, LandFundPence, ImportedLabourers}`.
- **ColonialLabourRegulation(market, scheme).** Pure function that applies a Wakefield scheme to a colonial market: preserves `FreeLabourers`, sets `IndependenceYears = ceil(ransom / annualSavings)` (assuming half the colonial wage is savings), flips `WakefieldSchemeApplied = true`, and sets `SurplusLabourExtractable = (IndependenceYears > 0)`. `RecordLandSale(scheme)` enforces Wakefield's self-financing invariant: each completed sale deposits the ransom into `LandFundPence` and imports one fresh labourer.
- **Store.** New `ColonialLabourMarketStore` interface with `Memory` and `MySQL` implementations (case-insensitive colony uniqueness, partial-update via `SELECT ... FOR UPDATE` inside a transaction). Migrations `00024_ch33_colonial_markets.sql` (single in-row table — unregulated and regulated states co-resident so the dashboard can compare without joins) and `00025_ch33_seed.sql`: Swan River 1829 (Peel's fiasco, no scheme), Saint Domingo 1700 (Spanish settlers, no scheme), South Australia 1836 (Wakefield testbed, scheme applied, 50-year dependence), Upper Canada 1830s (scheme applied, 8-year dependence).
- **HTTP endpoints.** `POST /v1/colonial-markets` (201 + Location, 409 on duplicate colony, 400 on validation failure), `GET /v1/colonial-markets` (200 — list, `[]` not `null`), `GET /v1/colonial-markets/{id}` (200/404), `POST /v1/colonial-markets/{id}/regulate` (200 — applies sufficient price, persists the regulated state), `POST /v1/colonial-markets/{id}/independence` (200 — per-worker projection with optional `years_limit`). api-gateway proxies `/v1/colonial-markets` and `/v1/colonial-markets/{rest...}` to simulation-engine.
- **React UI.** "Ch. 33 — The Modern Theory of Colonisation" panel: a four-card grid with the seeded colony list (free / wage-dependent pill), Wakefield's remedy controls (price per acre + desired acres → before/after comparison), a colony registration form, and a per-worker independence projection (years to ransom at the colonial wage).
- **Vol. I closure.** This is the final chapter of Volume I. The simulation now spans every stage of the dialectic from §1 (the commodity) through §33 (the negation of the negation read through Wakefield's confession).

### Vol. II Ch. 3 — what was built

`simulation-engine` adds the third lens on the industrial circuit: C′—M′—C(Lp+Mp)…P…C′. Unlike Form I (M…M′) and Form II (P…P), this form opens and closes in commodity-form. Surplus-value is already contained in the opening C′ — which is what makes it specifically C′ and not bare C. The chapter also introduces the social-capital aggregate view: summing individual circuits to read the total commodity-capital of society.

- **Domain types.** `CommodityCircuit{ID, AgentID, Initial, Mode, IsFirstInvestment, SocialCapitalLens, PartialSales, MPSources, Terminal, CreatedAt}`. `OpeningCommodityCapital{ConstantPence, VariablePence, SurplusPence, PoundsTotal}` with `.Total() = c+v+s` and `.ValueOriginalPence() = c+v`. `CommodityAugmented{…, CapitalisedPence, ClosedAt}` with `.IsExtended(openingTotal) bool`. `ValueComposition{ConstantShare, VariableShare, SurplusShare}` with `.Total() == realisedPence`. `SuccessivePartialSale{Quantity, RealisedPence, Decomposition}`. `MeansOfProductionSource{SourceKind}` — one of `producer_circuit`, `merchant_holding`, `import`.
- **Domain errors.** `ErrNotCommodityCapital` — SurplusPence ≤ 0 on a non-first-investment circuit. `ErrInsufficientMaterialBase` — capitalisation attempted with no MP sources linked. `ErrUnsourcedMeansOfProduction` — unrecognised source kind.
- **Pure functions.** `ComputeDecomposition(quantity int64, realisedPence Pence) ValueComposition` — aliquot decomposition from the opening c/v ratios; SurplusShare absorbs integer rounding so the three components always sum exactly to `realisedPence`. `EnsureMaterialAdequacy(amount Pence) error` — guards capitalisation against an empty MP-source list.
- **Store.** `CommodityCircuitStore` interface with `CreateCommodityCircuit`, `GetCommodityCircuit`, `ListCommodityCircuits`, `RecordPartialSale`, `LinkMPSource`, `CloseCommodityCircuit`. Memory and MySQL implementations. Terminal state stored inline as nullable columns in `commodity_circuits`. Migrations `00029_v2_ch03_commodity_circuits.sql` and `00030_v2_ch03_seed.sql`: SpinningMill1871Simple (closed, three partial sales summing to £500 C′) and SpinningMill1871Extended (closed, surplus fully capitalised, P′(n+1) = £500).
- **HTTP endpoints.** `POST /v1/commodity-circuits` (201 + Location), `GET /v1/commodity-circuits` (200 list, `[]` not null), `GET /v1/commodity-circuits/aggregate` (200 social-capital totals), `GET /v1/commodity-circuits/{id}` (200/404), `POST /v1/commodity-circuits/{id}/partial-sales` (201 — server computes aliquot decomposition), `POST /v1/commodity-circuits/{id}/mp-source` (201), `POST /v1/commodity-circuits/{id}/close` (200 with `is_extended` flag). api-gateway proxies `/v1/commodity-circuits` and `/v1/commodity-circuits/{rest...}` to simulation-engine.
- **React UI.** "Vol. II Ch. 3 — The Circuit of Commodity-Capital" panel: circuit diagram C′—M′—C(Lp+Mp)…P…C′, social-capital aggregate card, circuit list + create form (c/v/s pence, pounds total, mode), detail panel with opening C′ balance sheet, successive-partial-sales table (qty / realised / c-share / v-share / s-share), MP-source linker, and terminal C′ card (colour-coded simple vs extended).

### Vol. II Ch. 1 — what was built

`simulation-engine` models the circuit of money-capital from Capital Vol. II, Ch. 1:

- **Domain.** `circulation.MoneyCircuit` — the full M—C(Lp+Mp)…P…C′—M′ circuit as a state-machine with six moments (`M`, `M-C`, `P`, `C-prime`, `C-M-prime`, `M-prime`). Sub-structs: `MoneyAdvance` (M), `PurchasePhase` (M—C, with `LabourLeg` and `MeansLeg`), `ProductiveState` (P, c+v decomposition), `CommodityCapital` (C′, original + surplus), `Realisation` (M′, realised pence + surplus realised). Invariant enforcement: phase-order transitions, labour+means total must equal advance (within tolerance), magnitude-preservation check on realisation.

- **Store.** `MoneyCircuitStore` interface with `CreateMoneyCircuit`, `GetMoneyCircuit`, `ListMoneyCircuits` (agent\_id + moment filters), `RecordPurchase`, `RecordProductive`, `RecordCommodity`, `RecordRealisation`. Memory and MySQL implementations. Migrations `00031_v2_ch01_money_circuits.sql` and `00032_v2_ch01_seed.sql`: SpinningMill1871 at 156% rate (advance=£422) and 100% rate (advance=£436).

- **HTTP API.** `POST /v1/money-circuits`, `GET /v1/money-circuits`, `GET /v1/money-circuits/{id}`, `POST /v1/money-circuits/{id}/purchase`, `POST /v1/money-circuits/{id}/produce`, `POST /v1/money-circuits/{id}/commodity`, `POST /v1/money-circuits/{id}/realise`. Validation errors map to 400; `ErrNotFound` maps to 404.

- **React UI.** "Vol. II Ch. 1 — The Circuit of Money-Capital" panel: animated circuit diagram (active node highlighted per current moment), circuit list + create form (advance pence, optional agent ID), detail panel with moment-progress chips, balance sheet (M → C decomposition → P → C′ → M′/ΔM), and step-by-step phase forms (each form appears only when the circuit is at the preceding moment).

### Vol. II Ch. 5 — what was built

`market-service` gains a new `internal/circulation` package modelling the time dimension of capital's circuit: the period during which capital sits in money-form or commodity-form without producing surplus-value.

- **Domain types.** `TurnoverTime{ID, IndustrialCapitalID, Production, Circulation, CirculationOpen}` with `.TotalNanos()` and `.ActiveFractionBasisPoints()` (10000 × production / total). `ProductionTime{LabourTimeNanos, LabourInterruptionNanos, LatentNanos, NaturalProcessNanos}`. `CirculationTime{SellingTimeNanos, BuyingTimeNanos}`. `SellingPhase{…, Outcome SellingOutcome}` (sold/spoiled/partial/pending). `BuyingPhase{…, MarketLocation}`. `NaturalProcessSpan{Process NaturalProcessKind, DurationNanos}` (ripening/fermentation/tanning/drying/other). `LatentProductiveCapital{Pence, HeldAt, EnteredProductionAt *time.Time}`. `Perishability{CommodityID, WindowNanos}` with `.Spoiled(elapsed) bool`. `MarketSeparation{SellingMarketID, BuyingMarketID}`. `CirculationCompressionTarget{TargetCirculationNanos}`.
- **Sentinel errors.** `ErrConcurrentProductionAndCirculation`, `ErrSellingPhaseAlreadyOpen`, `ErrBuyingPhaseAlreadyOpen`, `ErrNoOpenSellingPhase`, `ErrNoOpenBuyingPhase`.
- **Store.** `TurnoverTimeStore` interface embedded in `Store`. Memory and MySQL implementations. Migrations `00003_v2_ch05_circulation_time.sql` and `00004_v2_ch05_seed.sql`: SpinningMill1871 (1-week production + 1-week circulation), grain/wine/leather natural-process spans, milk (2d) and beer (7d) perishability rows, £50 cotton stockpile latent capital.
- **HTTP endpoints.** `POST /v1/turnover-time` (201), `GET /v1/turnover-time` (200 list), `GET /v1/turnover-time/{id}` (200/404), `POST /v1/turnover-time/{id}/labour-time` (200), `POST /v1/turnover-time/{id}/labour-interruption` (200), `POST /v1/turnover-time/{id}/latent-mp` (201), `POST /v1/turnover-time/{id}/natural-process` (201), `POST /v1/turnover-time/{id}/selling-phase` (open: 201, close: 200), `POST /v1/turnover-time/{id}/buying-phase` (open: 201, close: 200), `GET /v1/turnover-time/{id}/active-fraction` (200), `POST /v1/perishability` (201), `POST /v1/market-separation` (201).
- **React UI.** "Vol. II Ch. 5 — The Time of Circulation" panel: (1) list of turnover records, (2) stacked-bar of production sub-spans + circulation sub-spans with active-fraction badge, (3) zero-circulation toggle showing rise in active fraction toward 100%, (4) natural-process examples (grain 90d / wine 14d / leather 30d), (5) perishability cliff widget counting elapsed time against commodity spoilage windows.

### Vol. II Ch. 5 — what was built

`market-service` gains a new `internal/circulation` package modelling the time dimension of capital's circuit: the period during which capital sits in money-form or commodity-form without producing surplus-value.

- **Domain types.** `TurnoverTime{ID, IndustrialCapitalID, Production, Circulation, CirculationOpen}` with `.TotalNanos()` and `.ActiveFractionBasisPoints()` (10000 × production / total). `ProductionTime{LabourTimeNanos, LabourInterruptionNanos, LatentNanos, NaturalProcessNanos}`. `CirculationTime{SellingTimeNanos, BuyingTimeNanos}`. `SellingPhase{…, Outcome SellingOutcome}` (sold/spoiled/partial/pending). `BuyingPhase{…, MarketLocation}`. `NaturalProcessSpan{Process NaturalProcessKind, DurationNanos}` (ripening/fermentation/tanning/drying/other). `LatentProductiveCapital{Pence, HeldAt, EnteredProductionAt *time.Time}`. `Perishability{CommodityID, WindowNanos}` with `.Spoiled(elapsed) bool`. `MarketSeparation{SellingMarketID, BuyingMarketID}`. `CirculationCompressionTarget{TargetCirculationNanos}`.
- **Sentinel errors.** `ErrConcurrentProductionAndCirculation`, `ErrSellingPhaseAlreadyOpen`, `ErrBuyingPhaseAlreadyOpen`, `ErrNoOpenSellingPhase`, `ErrNoOpenBuyingPhase`.
- **Store.** `TurnoverTimeStore` interface embedded in `Store`. Memory and MySQL implementations. Migrations `00003_v2_ch05_circulation_time.sql` and `00004_v2_ch05_seed.sql`: SpinningMill1871 (1-week production + 1-week circulation), grain/wine/leather natural-process spans, milk (2d) and beer (7d) perishability rows, £50 cotton stockpile latent capital.
- **HTTP endpoints.** `POST /v1/turnover-time` (201), `GET /v1/turnover-time` (200 list), `GET /v1/turnover-time/{id}` (200/404), `POST /v1/turnover-time/{id}/labour-time` (200), `POST /v1/turnover-time/{id}/labour-interruption` (200), `POST /v1/turnover-time/{id}/latent-mp` (201), `POST /v1/turnover-time/{id}/natural-process` (201), `POST /v1/turnover-time/{id}/selling-phase` (open: 201, close: 200), `POST /v1/turnover-time/{id}/buying-phase` (open: 201, close: 200), `GET /v1/turnover-time/{id}/active-fraction` (200), `POST /v1/perishability` (201), `POST /v1/market-separation` (201).
- **React UI.** "Vol. II Ch. 5 — The Time of Circulation" panel: (1) list of turnover records, (2) stacked-bar of production sub-spans + circulation sub-spans with active-fraction badge, (3) zero-circulation toggle showing rise in active fraction toward 100%, (4) natural-process examples (grain 90d / wine 14d / leather 30d), (5) perishability cliff widget counting elapsed time against commodity spoilage windows.

### Vol. II Ch. 6 — what was built

`market-service` gains a new `internal/costs` package modelling the costs that capital bears during its circulation phase — the deductions from surplus-value that arise not from production but from the metamorphoses of commodity-form and money-form.

- **Domain types.** `CirculationCost{ID, IndustrialCapitalID, Kind, Nature, Pence, IncurredAt}` — faux frais classified as `value_preserving` or `value_adding`. `CirculationCostKind` enum: `purchase_sale_time`, `book_keeping`, `money_as_faux_frais`, `supply_formation`, `commodity_supply`, `transportation`. `NatureOf(kind)` returns `value_adding` only for transportation. `CirculationAgent{…, LabourMinutesNecessary, LabourMinutesSurplus}` with `.ValueContribution() == 0` (circulation labour adds no value) and `.WageEffectivePence()`. `MoneyAsFauxFrais{Pence, AnnualReplacementPence}` — §I(c) national gold reserve wearing down as coin. `CommoditySupply{IsVoluntary bool, StorageCosts []StorageCost}` — voluntary reserve vs involuntary glut. `StorageCost{BuildingPence, LabourPence, PreservationLabourPence}` with `.Total()`. `TransportTariff{BasePencePerTonMile, FragilityMultiplierBasisPoints, BreakageRiskMultiplierBasisPoints}` with `.EffectiveRatePencePerTonMile()` applying the max multiplier in basis points. `TransportLeg{LabourCostPence, MeansOfTransportPence, SurplusPence}` with `.AddedValue()` — the value transport adds to the commodity.
- **Store.** `CirculationCostStore` interface embedded in `Store`. Memory and MySQL implementations (`mysql_costs.go`). Migrations `00005_v2_ch06_circulation_costs.sql` (7 tables) and `00006_v2_ch06_seed.sql`: Spinning Mill 1871 seller-agent, book-keeping cost 1200p, national gold reserve 1M/10K pence, voluntary + involuntary commodity supply with storage costs, Birmingham glass transport tariffs (historic 1×, modern 3× breakage-risk), two transport legs Birmingham→London.
- **HTTP endpoints.** `POST /v1/circulation-costs` (201), `GET /v1/circulation-costs` (200 list with filters), `GET /v1/circulation-costs/{id}` (200/404), `GET /v1/circulation-costs/aggregate` (200), `GET /v1/circulation-costs/system-faux-frais` (200), `POST /v1/circulation-agents` (201), `POST /v1/money-as-faux-frais` (201), `POST /v1/commodity-supplies` (201 + Location), `POST /v1/commodity-supplies/{id}/storage-cost` (201/404), `POST /v1/transport-tariffs` (201, upsert by commodity), `POST /v1/transport-legs` (201 + Location).
- **React UI.** "Vol. II Ch. 6 — The Costs of Circulation" panel: (1) circulation costs list with nature pills (faux frais / value-adding), (2) 3-card aggregate (value-preserving / value-adding / total), (3) system faux-frais money-reserve block, (4) transport tariff cards showing effective rate with multiplier, (5) transport leg cards showing route, weight, distance, and added value.

### Vol. II Ch. 7 — what was built

`simulation-engine` adds `internal/turnover/turnover.go` — Marx's analysis of the *turnover time* (u) and the *number of turnovers* (n = U/u, where U = 525,600 min = 1 year). Turnover is always measured on the productive or money circuit form; the commodity-capital form (Form III) is explicitly rejected for this analysis.

- **Domain types.** `Turnover{ID, IndustrialCapitalID, Lens, TurnoverTimeMinutes, Cycles}`. `TurnoverCycle{ID, TurnoverID, StartedAt, EndedAt, AdvancePence, ReturnedPence, ProductionMinutes, CirculationMinutes}` with `.Validate()` (invariant: `returned >= advance`). `TurnoverNumber{ID, TurnoverID, YearReferenceMinutes, TurnoverTimeMinutes, BasisPoints}` where `BasisPoints = 10000 × YearReferenceMinutes / TurnoverTimeMinutes` (integer, no float64). `CircuitLens` (`"money"` | `"productive"`) — `"commodity"` is rejected with `ErrLensNotApplicableForTurnover`.
- **Sentinel errors.** `ErrLensNotApplicableForTurnover`, `ErrCycleNotComplete`, `ErrNoCompletedCycles`.
- **Pure functions.** `ValidateLens(l)`, `ComputeTurnoverNumber(turnoverID, avgMinutes)`, `AverageTurnoverTime(cycles)`.
- **Store.** `TurnoverStore` interface with `CreateTurnover`, `GetTurnover`, `ListTurnovers`, `RecordCycle`, `RecomputeNumber`, `GetNumber`. Memory and MySQL implementations. `RecomputeNumber` averages `production_minutes + circulation_minutes` across all recorded cycles and upserts `turnover_numbers`.
- **Migrations.** `00035_v2_ch07_turnover.sql` (tables: `turnovers`, `turnover_cycles`, `turnover_numbers` with `UNIQUE KEY uq_tn_turnover`). `00036_v2_ch07_seed.sql`: SpinningMill (productive, 10,080 min, BP=521,428 ≈ 52/year), LocomotiveMaker (money, 777,600 min, BP=6,759 ≈ 2/3/year), WineMaturer (productive, 2,102,400 min, BP=2,500 = 1/4/year).
- **HTTP.** `POST /v1/turnovers`, `GET /v1/turnovers`, `GET /v1/turnovers/{id}`, `POST /v1/turnovers/{id}/cycles`, `POST /v1/turnovers/{id}/recompute-number`, `GET /v1/turnovers/{id}/number`. Gateway proxies `/v1/turnovers` and `/v1/turnovers/{rest...}` to simulation-engine.
- **React UI.** "Vol. II Ch. 7 — The Turnover Time" panel: one card per turnover showing lens, turnover time (formatted as weeks/years), cycle count, n expressed as `BasisPoints / 10000` per year.

### Vol. II Ch. 4 — what was built

`simulation-engine` adds `internal/circulation/industrial_capital.go` to synthesise the three circuit-forms into a single object: `IndustrialCapital`. Where Chs. 1–3 each presented one uninterrupted circuit form in isolation, Ch. 4 shows that every real industrial capital is permanently divided across all three stages simultaneously — Marx's *Nebeneinander* ("co-existence"). A stoppage at any one stage ripples as disorder into the co-existence of the other two.

- **Domain types.** `IndustrialCapital{ID, AgentID, MoneyCircuitID, ProductiveCircuitID, CommodityCircuitID, TotalPence, EconomyMode, StagnationToleranceTicks, Status, LatestDistribution, OpenBlocks, CreatedAt, UpdatedAt}`. `CapitalPart{ID, IndustrialCapitalID, Pence, Stage, EnteredStageAt}`. `StageDistribution{ID, IndustrialCapitalID, At, MoneyPence, ProductionPence, CommodityPence}` with `.Total()` and `.Validate(expectedTotal)`. `StageBlock{ID, IndustrialCapitalID, Stage, Reason, OpenedAt, ClosedAt}` with `.IsOpen()`. `ValueRevolutionEvent{ID, IndustrialCapitalID, Affecting, DirectionPence, OccurredAt}`. `MoneyCapitalSetFree` / `MoneyCapitalTiedUp` records emitted by `ApplyValueRevolution`. `MetamorphosisInterlock{ID, BuyerIndustrialCapitalID, SellerIndustrialCapitalID, SellerMPSourceOrigin, Pence, OccurredAt}`. `SupplyDemandImbalance{ID, IndustrialCapitalID, Period, DemandPence, SupplyPence, ExcessPence}` computed by `ComputeSupplyDemandImbalance(oc, surplusPence)`. `AggregateSupplyDemandImbalance` — class-level totals. `SinkingFund{ID, IndustrialCapitalID, FixedCapitalPence, LifetimeYears, AnnualPaymentPence, AccumulatedPence}` with `.Validate()` (invariant: `annualPayment * lifetime == fixedCapital`) and `.Tick()`.
- **Sentinel errors.** `ErrNonCapitalCounterparty`, `ErrCreditWithoutMoneyEconomy`, `ErrSinkingFundMismatch`, `ErrNonZeroDirectionPence`.
- **Pure functions.** `ApplyValueRevolution(evt)` — routes the event to `MoneyCapitalSetFree` (prices fell) or `MoneyCapitalTiedUp` (prices rose). `ComputeSupplyDemandImbalance(ic, period, oc, surplusPence)` — demand = c+v, supply = c+v+s, excess = s. `OrganicComposition.RatioBasisPoints()` — 10000 × c / v (integer, no float64).
- **Store.** `IndustrialCapitalStore` interface with `CreateIndustrialCapital`, `GetIndustrialCapital`, `ListIndustrialCapitals`, `RecordCapitalPart`, `SnapshotStageDistribution`, `OpenStageBlock`, `CloseStageBlock`, `RecordValueRevolution`, `RecordInterlock`, `RecordSupplyDemand`, `GetSupplyDemand`, `AggregateSupplyDemand`, `SetSinkingFund`, `TickSinkingFund`. Memory and MySQL implementations. Migrations `00033_v2_ch04_industrial_capital.sql` (eight tables) and `00034_v2_ch04_seed.sql`: SpinningMill1871 seed capital (id `5eed000000000000000401`) with an opening stage distribution, a supply-demand imbalance, and a sinking fund.
- **HTTP endpoints.** `POST /v1/industrial-capitals` (201 + Location), `GET /v1/industrial-capitals` (200 list), `GET /v1/industrial-capitals/{id}` (200/404, includes latest distribution + open blocks), `POST /v1/industrial-capitals/{id}/parts` (201), `POST /v1/industrial-capitals/{id}/snapshot` (201), `POST /v1/industrial-capitals/{id}/blocks` (201), `POST /v1/industrial-capitals/{id}/blocks/{block_id}/close` (200), `POST /v1/industrial-capitals/{id}/value-revolution` (200), `POST /v1/industrial-capitals/{id}/interlocks` (201), `POST /v1/industrial-capitals/{id}/supply-demand` (201), `GET /v1/industrial-capitals/{id}/supply-demand` (200), `GET /v1/industrial-capitals/supply-demand/aggregate` (200), `POST /v1/industrial-capitals/{id}/sinking-fund` (201), `POST /v1/industrial-capitals/{id}/sinking-fund/tick` (200), `GET /v1/historical-stages/{id}/genesis` (200). api-gateway proxies `/v1/industrial-capitals`, `/v1/industrial-capitals/{rest...}`, and `/v1/historical-stages/{rest...}` to simulation-engine.
- **React UI.** "Vol. II Ch. 4 — The Three Formulas of the Circuit" panel: the three circuit formulas (I, II, III) displayed as monospace code blocks, a card list of industrial capitals with colour-coded stage status (M=blue, P=green, C′=orange, blocked=red), a proportional M/P/C′ stage-distribution bar, distribution detail (pence per stage), and open stagnation-block badges.

### Volume II Roadmap — The Process of Circulation of Capital

Spec sweep pending; titles below are drawn from the vault filenames and will be refined as each chapter spec is authored at `marx-engels/1885/capital-volume-ii/specs/`. Primary-service column is a planning guess.

| Chapter   | Status      | Concepts                                                                                          | Primary services                  |
|-----------|-------------|---------------------------------------------------------------------------------------------------|-----------------------------------|
| Ch. 1     | ✅ Done     | The Circuit of Money-Capital — M—C…P…C'—M' as the money-form of the circuit                       | simulation-engine                 |
| Ch. 2     | ✅ Done     | The Circuit of Productive Capital — P…C'—M'—C…P as the production-form of the circuit             | simulation-engine                 |
| Ch. 3     | ✅ Done     | The Circuit of Commodity-Capital — C'—M'—C…P…C' as the commodity-form of the circuit              | simulation-engine                 |
| Ch. 4     | ✅ Done     | The Three Formulas of the Circuit — interruption, continuity, the unity of all three forms        | simulation-engine                 |
| Ch. 5     | ✅ Done     | The Time of Circulation — selling time, buying time, capital tied up in metamorphosis              | market-service                    |
| Ch. 6     | ✅ Done     | The Costs of Circulation — faux frais (purchase/sale time, book-keeping, money reserve) vs value-adding transport costs | market-service |
| Ch. 7     | ✅ Done     | The Turnover Time and the Number of Turnovers — turnover time, number of turnovers, basis-points    | simulation-engine                 |
| Ch. 8     | ⏳ Pending  | Fixed Capital and Circulating Capital                                                               | simulation-engine                 |
| Ch. 9     | ⏳ Pending  | The Aggregate Turnover of Advanced Capital — cycles of turnover                                     | simulation-engine                 |
| Ch. 10    | ⏳ Pending  | Theories of Fixed and Circulating Capital — Physiocrats and Adam Smith                              | simulation-engine                 |
| Ch. 11    | ⏳ Pending  | Theories of Fixed and Circulating Capital — Ricardo                                                 | simulation-engine                 |
| Ch. 12    | ⏳ Pending  | The Working Period                                                                                  | simulation-engine                 |
| Ch. 13    | ⏳ Pending  | The Time of Production                                                                              | simulation-engine                 |
| Ch. 14    | ⏳ Pending  | The Time of Circulation                                                                             | simulation-engine, market-service |
| Ch. 15    | ⏳ Pending  | The Effects of a Change of Prices on capital advanced and released                                  | simulation-engine                 |
| Ch. 16    | ⏳ Pending  | The Turnover of Variable Capital — annual rate of surplus-value                                     | simulation-engine                 |
| Ch. 17    | ⏳ Pending  | The Circulation of Surplus-Value (introductory)                                                     | simulation-engine                 |
| Ch. 18    | ⏳ Pending  | The Role of Money-Capital in Reproduction                                                           | simulation-engine                 |
| Ch. 19    | ⏳ Pending  | Former Presentations of the Subject — Quesnay, Smith, and others                                    | simulation-engine                 |
| Ch. 20    | ⏳ Pending  | Simple Reproduction (Vol. II treatment — Departments I & II)                                        | simulation-engine                 |
| Ch. 21    | ⏳ Pending  | Accumulation and Reproduction on an Extended Scale                                                  | simulation-engine                 |

### Volume III Roadmap — The Process of Capitalist Production as a Whole

Spec sweep pending; titles below are drawn from the vault filenames at `marx-engels/1894/capital-volume-iii/texts/` and will be refined as each chapter spec is authored. Primary-service column is a planning guess; most rows are `finance-service` with `simulation-engine` for cross-capital aggregations (average rate of profit, formation of prices of production).

| Chapter   | Status      | Concepts                                                                                          | Primary services                  |
|-----------|-------------|---------------------------------------------------------------------------------------------------|-----------------------------------|
| Ch. 1     | ⏳ Pending  | Cost-Price and Profit (k + p) — surplus-value reappears as profit on total capital                | finance-service                   |
| Ch. 2     | ⏳ Pending  | The Rate of Profit — p′ = s / (c + v)                                                              | finance-service                   |
| Ch. 3     | ⏳ Pending  | Relation of the Rate of Profit to the Rate of Surplus-Value                                        | finance-service                   |
| Ch. 4     | ⏳ Pending  | The Effect of the Turnover on the Rate of Profit                                                   | finance-service, simulation-engine|
| Ch. 5     | ⏳ Pending  | Economy in the Employment of Constant Capital                                                       | finance-service                   |
| Ch. 6     | ⏳ Pending  | The Effect of Price Fluctuation on the Rate of Profit                                              | finance-service                   |
| Ch. 7     | ⏳ Pending  | Supplementary Remarks (on the rate of profit)                                                       | finance-service                   |
| Ch. 8     | ⏳ Pending  | Different Compositions of Capitals in Different Branches → differences in rates of profit          | finance-service                   |
| Ch. 9     | ⏳ Pending  | Formation of a General Rate of Profit; transformation of values into prices of production          | finance-service, simulation-engine|
| Ch. 10    | ⏳ Pending  | Equalisation of the General Rate of Profit through Competition; market prices and market values    | finance-service, market-service   |
| Ch. 11    | ⏳ Pending  | Effects of General Wage Fluctuations on Prices of Production                                       | finance-service                   |
| Ch. 12    | ⏳ Pending  | Supplementary Remarks (on prices of production)                                                     | finance-service                   |
| Ch. 13    | ⏳ Pending  | The Law of the Tendential Fall in the Rate of Profit — As Such                                     | finance-service, simulation-engine|
| Ch. 14    | ⏳ Pending  | Counteracting Influences                                                                            | finance-service                   |
| Ch. 15    | ⏳ Pending  | Exposition of the Internal Contradictions of the Law                                                | finance-service, simulation-engine|
| Ch. 16    | ⏳ Pending  | Commercial Capital                                                                                  | finance-service                   |
| Ch. 17    | ⏳ Pending  | Commercial Profit                                                                                   | finance-service                   |
| Ch. 18    | ⏳ Pending  | The Turnover of Merchant's Capital; Prices                                                          | finance-service                   |
| Ch. 19    | ⏳ Pending  | Money-Dealing Capital                                                                               | finance-service                   |
| Ch. 20    | ⏳ Pending  | Historical Facts about Merchant's Capital                                                           | finance-service                   |
| Ch. 21    | ⏳ Pending  | Interest-Bearing Capital                                                                            | finance-service                   |
| Ch. 22    | ⏳ Pending  | Division of Profit; Rate of Interest; Natural Rate of Interest                                     | finance-service                   |
| Ch. 23    | ⏳ Pending  | Interest and Profit of Enterprise                                                                   | finance-service                   |
| Ch. 24    | ⏳ Pending  | Externalisation of the Relations of Capital in the Form of Interest-Bearing Capital                | finance-service                   |
| Ch. 25    | ⏳ Pending  | Credit and Fictitious Capital                                                                       | finance-service                   |
| Ch. 26    | ⏳ Pending  | Accumulation of Money-Capital — its influence on the interest rate                                  | finance-service                   |
| Ch. 27    | ⏳ Pending  | The Role of Credit in Capitalist Production                                                         | finance-service                   |
| Ch. 28    | ⏳ Pending  | Medium of Circulation and Capital — views of Tooke and Fullarton                                    | finance-service                   |
| Ch. 29    | ⏳ Pending  | Component Parts of Bank Capital                                                                     | finance-service                   |
| Ch. 30    | ⏳ Pending  | Money-Capital and Real Capital, I                                                                   | finance-service                   |
| Ch. 31    | ⏳ Pending  | Money-Capital and Real Capital, II                                                                  | finance-service                   |
| Ch. 32    | ⏳ Pending  | Money-Capital and Real Capital, III                                                                 | finance-service                   |
| Ch. 33    | ⏳ Pending  | The Medium of Circulation in the Credit System                                                      | finance-service                   |
| Ch. 34    | ⏳ Pending  | The Currency Principle and the English Bank Legislation of 1844                                     | finance-service                   |
| Ch. 35    | ⏳ Pending  | Precious Metal and Rate of Exchange                                                                 | finance-service                   |
| Ch. 36    | ⏳ Pending  | Pre-Capitalist Relationships                                                                        | finance-service                   |
| Ch. 37    | ⏳ Pending  | Introduction (to Ground-Rent)                                                                       | finance-service                   |
| Ch. 38    | ⏳ Pending  | General Remarks (on differential rent)                                                              | finance-service                   |
| Ch. 39    | ⏳ Pending  | First Form of Differential Rent — Differential Rent I                                              | finance-service                   |
| Ch. 40    | ⏳ Pending  | Second Form of Differential Rent — Differential Rent II                                            | finance-service                   |
| Ch. 41    | ⏳ Pending  | Differential Rent II — Constant Price of Production                                                | finance-service                   |
| Ch. 42    | ⏳ Pending  | Differential Rent II — Falling Price of Production                                                 | finance-service                   |
| Ch. 43    | ⏳ Pending  | Differential Rent II — Rising Price of Production                                                  | finance-service                   |
| Ch. 44    | ⏳ Pending  | Differential Rent Also on the Worst Cultivated Soil                                                 | finance-service                   |
| Ch. 45    | ⏳ Pending  | Absolute Ground-Rent                                                                                | finance-service                   |
| Ch. 46    | ⏳ Pending  | Building Site Rent — Rent in Mining — Price of Land                                                | finance-service                   |
| Ch. 47    | ⏳ Pending  | Genesis of Capitalist Ground-Rent                                                                   | finance-service                   |
| Ch. 48    | ⏳ Pending  | The Trinity Formula — capital→profit, land→rent, labour→wages                                      | finance-service                   |
| Ch. 49    | ⏳ Pending  | Concerning the Analysis of the Process of Production                                                | finance-service, simulation-engine|
| Ch. 50    | ⏳ Pending  | Illusions Created by Competition                                                                    | finance-service                   |
| Ch. 51    | ⏳ Pending  | Distribution Relations and Production Relations                                                     | finance-service                   |
| Ch. 52    | ⏳ Pending  | Classes (unfinished — Marx died mid-chapter)                                                        | finance-service                   |
