---
chapter: 15
title: "Machinery and Modern Industry"
status: proposed
primary_service: simulation-engine
---

## Concepts → types

| Marx term | Go identifier | Kind | Package | Notes |
|---|---|---|---|---|
| Machine (motor + transmission + tool) | `Machine` | type | `services/simulation-engine/internal/machinery` | Three-part structure: MotorMechanism, TransmittingMechanism, WorkingTool |
| Machine ID | `MachineID` | type | `services/simulation-engine/internal/machinery` | 96-bit hex; NewMachineID() follows commodity.NewID() pattern |
| Total value of machine | `MachineValue` | type alias LabourMinutes | `services/simulation-engine/internal/machinery` | Labour crystallised in the machine at time of production |
| Lifetime of machine (working days) | `LifespanDays` | type int64 | `services/simulation-engine/internal/machinery` | Duration over which the machine transfers its value |
| Daily value transfer (wear and tear) | `DailyWearAndTear` | func | `services/simulation-engine/internal/machinery` | MachineValue / LifespanDays |
| Value transferred per unit of product | `ValueTransferredPerUnit` | func | `services/simulation-engine/internal/machinery` | DailyWearAndTear / UnitsProducedPerDay |
| Material wear and tear | `MaterialWear` | type | `services/simulation-engine/internal/machinery` | Wear proportional to use; tracked in LabourMinutes |
| Moral depreciation | `MoralDepreciation` | type | `services/simulation-engine/internal/machinery` | Value lost because cheaper/better machines are produced |
| Productive power of machine | `ProductivePower` | type int64 | `services/simulation-engine/internal/machinery` | Units of output the machine produces per working day |
| Factory (system of machines) | `Factory` | type | `services/simulation-engine/internal/machinery` | Aggregate of Machine records sharing a PrimeMover |
| Prime mover | `PrimeMover` | type | `services/simulation-engine/internal/machinery` | Steam engine, water wheel, etc. |
| Production tick / simulation period | `Tick` | type | `services/simulation-engine/internal/engine` | One time-step; advances wear, value transfer, output |
| Labour displaced by machine | `LabourDisplaced` | func | `services/simulation-engine/internal/machinery` | LabourMinutes formerly required to produce same output by hand |
| Intensification of labour | `IntensityFactor` | type float64 | `services/simulation-engine/internal/machinery` | Multiplier on effective labour per hour when day is shortened |
| Variable capital | `VariableCapital` | type alias int64 (pence) | `services/simulation-engine/internal/machinery` | Capital portion laid out in labour-power |
| Constant capital | `ConstantCapital` | type alias int64 (pence) | `services/simulation-engine/internal/machinery` | Capital portion in means of production |
| Industrial cycle phase | `CyclePhase` | type (enum) | `services/simulation-engine/internal/engine` | Prosperity / Overproduction / Crisis / Stagnation |

## Fixtures

- **§2** `"the machine ... enters into the value-begetting process only by bits. It never adds more value than it loses, on an average, by wear and tear"` → `DailyWearAndTear(machineValue=600_000, lifespanDays=1000) == 600` (LabourMinutes per day)

- **§2** `"2½ operatives spin weekly 365⅝ lbs. of yarn ... 366 lbs. of cotton absorb only 150 hours' labour"` → machine-spun 366 lbs = 9,000 minutes (150h × 60); hand-spun 366 lbs = 1,620,000 minutes (27,000h × 60); ratio = 180× labour saving

- **§2** `"a machine working 16 hours daily for 7½ years covers as long a working period as the same machine working only 8 hours daily for 15 years. But in the first case the value would be reproduced twice as quickly"` → `TotalValueTransferred(hoursPerDay=16, days=2737) == TotalValueTransferred(hoursPerDay=8, days=5475)`; surplus extracted twice as fast in the 16h case

- **§2** `"a steam-plough does as much work in one hour at a cost of three-pence, as 66 men at a cost of 15 shillings"` → `LabourDisplaced` = 66 man-hours; machine cost per hour = 3 pence; displaced wage bill = 180 pence; `LabourDisplaced > MachineCostInLabour`

- **§3B** `"It loses exchange-value, either by machines of the same sort being produced cheaper than it, or by better machines entering into competition with it"` → `MoralDepreciation(m).Value > 0` when a new machine with lower `MachineValue` enters the pool at the same `ProductivePower`

- **§6** `"variable capital of £3,000 ... 100 workmen at £30 a year ... machinery costing £1,500 displaces 50 men; variable capital becomes £1,500"` → `VariableCapital` halves; `ConstantCapital` rises by £1,500; total capital £6,000 unchanged; headcount drops 100 → 50

- **§8A** `"a single needle-machine makes 145,000 in a working-day of 11 hours. One woman superintends four such machines producing near upon 600,000 needles a day"` → `ProductivePower` per machine = 145,000 units/day; four machines under single supervision ≈ 580,000 units/day

## Invariants

- `DailyWearAndTear(m) == m.MachineValue / m.LifespanDays` [§2] — a machine adds no more value per day than it loses; any tick crediting more is a domain error

- `ValueTransferredPerUnit(m, unitsPerDay) == DailyWearAndTear(m) / unitsPerDay` [§2] — per-unit value transfer falls inversely as productive power rises; doubling output halves the value added per unit

- `LabourDisplaced(m) > LabourEmbodiedInMachine(m)` [§2] — "there is always a difference of labour saved in favour of the machine"; if false, the machine offers no advantage over hand labour

- `VariableCapital + ConstantCapital == TotalCapital` [§6] — mechanisation converts variable into constant capital, but total capital is conserved in the substitution period; only the composition shifts

- `MoralDepreciation(m).Value >= 0` [§3B] — moral depreciation is non-negative; competition from cheaper successors never raises a machine's value

## Scope

### This chapter builds
- Services: `simulation-engine` (primary — machinery domain, Factory, Tick loop); `commodity-service` (extended — `ValueTransferredPerUnit` links machine wear to commodity SNLT reduction)
- New domain types:
  - `Machine` — three-part instrument (motor, transmission, tool) with value, lifespan, productive power
  - `MachineID` — 96-bit hex identifier
  - `LifespanDays` — working-day lifetime for amortisation
  - `DailyWearAndTear` — value slice transferred to product per tick
  - `ValueTransferredPerUnit` — per-commodity value addition from machine wear
  - `MaterialWear` — use-proportional degradation accumulator
  - `MoralDepreciation` — market-driven value loss, independent of use
  - `ProductivePower` — units of output per working day
  - `Factory` — aggregate of machines sharing a prime mover; owns the tick
  - `PrimeMover` — power source (steam engine, water wheel, etc.)
  - `LabourDisplaced` — labour-minutes the machine replaces per day
  - `IntensityFactor` — labour density multiplier when working day is shortened (§3C)
  - `VariableCapital` / `ConstantCapital` — capital-composition pair
  - `CyclePhase` — industrial-cycle enum (Prosperity / Overproduction / Crisis / Stagnation)
  - `Tick` — one simulation time-step
- New HTTP endpoints:
  - `POST /v1/machines` — register a machine with value, lifespan, productive power, and prime mover
  - `GET /v1/machines` — list all machines
  - `GET /v1/machines/{id}` — fetch a single machine record
  - `POST /v1/factories` — assemble machines under a shared prime mover
  - `GET /v1/factories/{id}` — fetch factory state (machines, productive power, accumulated wear)
  - `POST /v1/factories/{id}/tick` — advance one simulation period; returns value transferred and output units
  - `GET /v1/machines/{id}/wear` — return accumulated MaterialWear and MoralDepreciation
- React: "Ch. 15 — Machinery and Modern Industry" panel: machine registry, factory view with tick controls, per-tick value-transfer vs. hand-labour baseline, capital-composition gauge (constant vs. variable)

### Explicitly deferred to later chapters
- Reserve army of labour — displaced workers from §§5–7 form the industrial reserve army; explicit subject of Ch. 25 (General Law of Capitalist Accumulation)
- Wages under machinery — systematic treatment of how machinery alters wage rates belongs to Ch. 17–20 (Wages)
- Agriculture and machinery — §10 (Modern Industry and Agriculture) introduces ground-rent implications; deferred to Ch. 24–25
- Factory Acts / legal working-day enforcement — §9 connects to Ch. 10 (The Working Day); maximum-hours simulation already partially scoped there
- International division of labour / colonial raw-material extraction — §7 references India, Ireland, United States; cross-border production chains deferred to accumulation and world-market chapters
