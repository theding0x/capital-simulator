---
chapter: 14
title: "Division of Labour and Manufacture"
status: proposed
primary_service: agent-service
---

## Concepts → types

| Marx term | Go identifier | Kind | Package | Notes |
|---|---|---|---|---|
| manufacture | `Manufacture` | type | `services/agent-service/internal/agent` | A `Cooperation` that has been reorganised around a fixed division of detail labour; holds a set of `DetailRole`s and their assigned workers |
| detail labourer | `DetailRole` | type | `services/agent-service/internal/agent` | A named fractional function within a manufacture; each worker is annexed to exactly one role |
| heterogeneous manufacture | `HeterogeneousManufacture` | const | `services/agent-service/internal/agent` | `ManufactureForm` variant — product assembled from independently producible partial products (e.g. carriage, watch) |
| serial manufacture | `SerialManufacture` | const | `services/agent-service/internal/agent` | `ManufactureForm` variant — product passes through connected sequential stages (e.g. needle wire through 72+ hands) |
| manufacture form | `ManufactureForm` | type | `services/agent-service/internal/agent` | `string` enum: `"heterogeneous"` or `"serial"` |
| collective labourer | `CollectiveLabourer` | type | `services/agent-service/internal/agent` | The organic whole formed by all `DetailRole` workers in a `Manufacture`; its productive power exceeds the sum of its parts |
| detail worker output (partial product) | `PartialProduct` | type | `services/agent-service/internal/agent` | `LabourMinutes` produced by one `DetailRole` per period; not a commodity on its own |
| hierarchy of labour-powers | `LabourHierarchy` | type | `services/agent-service/internal/agent` | Ordered slice of `DetailRole` by skill level; determines relative value of labour-power per role |
| skilled / unskilled distinction | `SkillLevel` | type | `services/agent-service/internal/agent` | `string` enum: `"skilled"` or `"unskilled"`; unskilled roles have lower `LabourPowerValue` |
| value of labour-power (per role) | `RoleLabourPowerValue` | func | `services/agent-service/internal/agent` | Returns `LabourMinutes` SNLT to reproduce a worker assigned to a given `DetailRole`; falls as skill requirements fall |
| social productive power of manufacture | `ManufactureProductivePower` | func | `services/agent-service/internal/agent` | Collective output in `LabourMinutes` of use-value; exceeds `CollectiveWorkingDay` of the same workers without division of labour |
| proportionality law (iron law) | `ProportionalGroupSize` | func | `services/agent-service/internal/agent` | Given a target output rate, returns the required number of workers per `DetailRole` so that partial products flow without bottlenecks |
| minimum capital for manufacture | `ManufactureMinimumCapital` | func | `services/agent-service/internal/agent` | `MinimumCapital` variant: must cover wages for all roles in their required proportions, plus raw material consumed at the faster collective rate |
| division of labour in society (contrast) | `SocialDivisionOfLabour` | const | `services/agent-service/internal/agent` | Marker constant used in commentary/tests to distinguish market-mediated social division (anarchic, mediated by commodity exchange) from the internally planned workshop division |
| two-fold origin | `ManufactureOrigin` | type | `services/agent-service/internal/agent` | `string` enum: `"combination"` (diverse handicrafts united) or `"splitting"` (one handicraft subdivided) |
| tool specialisation | `SpecialisedTool` | type | `services/agent-service/internal/agent` | A tool adapted to a single `DetailRole`; count in Birmingham hammer example drives the differentiation of instruments |

## Fixtures

- **§1 (two-fold origin — combination)** A carriage manufacture assembles wheelwrights, harness-makers, locksmiths, etc. under one capitalist; over time each loses general craft ability and performs only one operation. → `Manufacture{Origin: "combination", Roles: [wheelwright, harness-maker, locksmith, ...]}` with `ManufactureOrigin == "combination"`; `len(Roles) >= 2`

- **§1 (two-fold origin — splitting)** A needle-maker guild performs 20 sequential operations; English needle manufacture assigns each to a separate detail labourer. → `Manufacture{Origin: "splitting", Roles: [wire-draw, wire-straighten, wire-cut, wire-point, ...]}` with `len(Roles) == 20` (or more after further subdivision); `ManufactureOrigin == "splitting"`

- **§2 (type manufacture proportionality)** `"there are four founders and two breakers to one rubber: the founder casts 2,000 type an hour, the breaker breaks up 4,000, and the rubber polishes 8,000"` → `ProportionalGroupSize(roles=[founder@2000/hr, breaker@4000/hr, rubber@8000/hr], targetOutput=8000)` returns `{founder: 4, breaker: 2, rubber: 1}`; product of each group's rate × headcount equals `8000`

- **§2 (serial manufacture — needle wire)** Wire passes through 72 different detail workmen. → `Manufacture{Form: "serial", Roles: [r1..r72]}` with `len(Roles) == 72`; `ManufactureProductivePower` is strictly greater than the same 72 workers performing all 72 ops sequentially in isolation

- **§3 (glass furnace group)** A glass "hole" consists of 5 detail workers (bottlemaker, blower, gatherer, putter-up, taker-in); the whole group is paralysed if any one is absent. → `CollectiveLabourer{Roles: [bottlemaker, blower, gatherer, putter-up, taker-in]}` has `IsParalysed(absentRoles=1) == true`; for `absentRoles == 0`, output is positive

- **§3 (scaling — multiples of the group)** `"when once the most fitting proportion has been experimentally established ... that scale can be extended only by employing a multiple of each particular group"` → `ScaleManufacture(m, multiplier=3)` triples every role's headcount and triples collective output; non-integer multiples are invalid

- **§5 (minimum capital grows)** `"the minimum amount of capital ... must keep increasing"` because raw material consumption grows in the same ratio as productive power. → `ManufactureMinimumCapital(roles, dailyWagePerRole, rawMaterialCostFactor)` exceeds `MinimumCapital(n=totalWorkers, dailyWage)` by `rawMaterialCostFactor * ManufactureProductivePower / CooperationProductivePower`

- **§5 (labour-power value falls with skill)** Manufacture begets a class of unskilled labourers for whom cost of apprenticeship vanishes; value of their labour-power falls. → `RoleLabourPowerValue(role{SkillLevel: "unskilled"})` < `RoleLabourPowerValue(role{SkillLevel: "skilled"})` for otherwise identical roles; the gap is apprenticeship cost expressed as `LabourMinutes`

## Invariants

- `PartialProduct` of any single `DetailRole` is not a commodity; only the finished output of the `CollectiveLabourer` is. [§4] — enforced structurally: `PartialProduct` has no `ExchangeValue` field
- `ManufactureProductivePower(m) > CollectiveWorkingDay(m.TotalWorkers, m.IndividualWorkingDay)` for any manufacture with `len(Roles) > 1` [§2, §3]
- `ProportionalGroupSize` must satisfy: for every `DetailRole` r, `r.OutputRate * r.Headcount == targetOutputRate` [§2]
- `ScaleManufacture(m, k)` for integer `k >= 1` must produce `ManufactureProductivePower == k * original power` [§3]
- `RoleLabourPowerValue(unskilled) <= RoleLabourPowerValue(skilled)` for all roles; equality only if skill distinction carries no apprenticeship cost [§2]
- `ManufactureMinimumCapital(m) >= MinimumCapital(m.TotalWorkers, averageDailyWage)` — manufacture always requires more capital than simple co-operation of the same headcount [§5]
- Detail labour within manufacture is governed by the `"a priori"` plan of the capitalist; social division of labour is governed by the `"a posteriori"` law of the market. The two must never be conflated in domain logic. [§4]

## Scope

### This chapter builds
- Services: `agent-service`
- New domain types (all in `services/agent-service/internal/agent`):
  - `ManufactureForm` — `string` named type; values `"heterogeneous"` or `"serial"`
  - `ManufactureOrigin` — `string` named type; values `"combination"` or `"splitting"`
  - `SkillLevel` — `string` named type; values `"skilled"` or `"unskilled"`
  - `DetailRole` — named fractional function within a manufacture; fields: `Name string`, `SkillLevel SkillLevel`, `OutputRatePerHour int64`, `HeadCount int`, `ToolName string`
  - `LabourHierarchy` — `[]DetailRole` ordered by skill level descending; type with `Len()`, `At(i)`, `Unskilled()` helper
  - `PartialProduct` — `LabourMinutes` alias representing fractional output of one `DetailRole` per period; deliberately not embeddable in a commodity or exchange
  - `CollectiveLabourer` — aggregate of all `DetailRole` workers in a `Manufacture`; methods: `TotalWorkers() int`, `IsParalysed(absentRoles int) bool`, `OutputPerPeriod(periodMinutes int64) LabourMinutes`
  - `Manufacture` — extends `Cooperation` (Ch. 13); fields: `Form ManufactureForm`, `Origin ManufactureOrigin`, `Roles []DetailRole`; holds reference to a `Capitalist` owner
- New pure functions:
  - `ProportionalGroupSize(roles []DetailRole, targetOutputRate int64) (map[string]int, error)` — returns required headcount per role; returns `ErrBottleneck` if no integer solution exists
  - `ManufactureProductivePower(m Manufacture, periodMinutes int64) LabourMinutes` — collective output above simple co-operation baseline
  - `RoleLabourPowerValue(role DetailRole, apprenticeshipMinutes LabourMinutes) LabourMinutes` — SNLT to reproduce a worker in that role
  - `ManufactureMinimumCapital(m Manufacture, rawMaterialCostFactor float64) Pence` — minimum capital the capitalist must hold before the manufacture can operate
  - `ScaleManufacture(m Manufacture, multiplier int) (Manufacture, error)` — returns a new `Manufacture` with all role headcounts scaled; `ErrInvalidMultiplier` for `multiplier < 1`
- New HTTP endpoints (agent-service, routed via api-gateway):
  - `POST /v1/manufactures` — create a manufacture (capitalist ID, form, origin, list of detail roles with headcounts and skill levels)
  - `GET /v1/manufactures/{id}` — retrieve a manufacture with its collective labourer and hierarchy
  - `GET /v1/manufactures` — list all manufactures; optional `?capitalist_id=` and `?form=` filters
  - `POST /v1/manufactures/{id}/proportional-group-size` — body: `{"target_output_rate": N}`; returns required headcount per role or `ErrBottleneck`
  - `POST /v1/manufactures/{id}/scale` — body: `{"multiplier": K}`; returns scaled manufacture
  - `GET /v1/manufactures/{id}/collective-labourer` — returns `CollectiveLabourer` summary: total workers, paralysis flag (absent=0), output per 8-hour period
  - `GET /v1/manufactures/{id}/minimum-capital?raw_material_cost_factor=F` — returns `Pence` minimum capital
- Persistence: new `manufactures` and `detail_roles` tables via migration `00006_ch14_manufacture.sql` in `services/agent-service/internal/store/migrations/`; `Manufacture` store methods added to the existing store interface
- React: "Ch. 14 — Division of Labour and Manufacture" panel; form to create a manufacture (select capitalist, choose form and origin, add detail roles with name/skill/output-rate/headcount); display collective labourer summary, hierarchy table, proportional group size calculator, scale control, and minimum capital readout; note on the panel that detail labour produces no commodities — only the collective product enters exchange

### Explicitly deferred to later chapters
- Machinery as successor to manufacture — Ch. 15 (Machinery and Modern Industry); manufacture's tool specialisation is named here but the machine as congealed collective labour is Ch. 15's subject
- Effect of manufacture's reduced `RoleLabourPowerValue` on the aggregate rate of surplus-value — reopened in Ch. 16 (Absolute and Relative Surplus Value)
- Wage forms shaped by the skill hierarchy (time-wages vs piece-wages) — Ch. 20–21
- The reserve army created by deskilling — Ch. 25 (General Law of Capitalist Accumulation)
- Social division of labour at the level of whole branches of industry and world market — Ch. 3 (market-service) models commodity exchange but the societal DOL analysis remains descriptive commentary in this chapter; no new market-service endpoints needed
- Colonial markets as condition of possibility for manufacture's scale — historical note only; not modelled as a service interaction until accumulation chapters
