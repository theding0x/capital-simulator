---
chapter: 10
title: "The Working-Day"
status: proposed
primary_service: agent-service
---

## Concepts → types

| Marx term | Go identifier | Kind | Package | Notes |
|---|---|---|---|---|
| working day | `WorkingDay` | type | `agent` | Struct encoding necessary + surplus labour durations |
| necessary labour time | `NecessaryLabourMinutes` | type (`int64`) | `agent` | AB segment; time to reproduce labour-power value |
| surplus labour time | `SurplusLabourMinutes` | type (`int64`) | `agent` | BC segment; extension beyond AB |
| working day length | `TotalMinutes() int64` | method on `WorkingDay` | `agent` | Returns `NecessaryLabourMinutes + SurplusLabourMinutes` |
| rate of surplus-value | `RateOfSurplusValue() float64` | method on `WorkingDay` | `agent` | `SurplusLabourMinutes / NecessaryLabourMinutes` |
| normal working day | `NormalWorkingDay` | type | `agent` | A `WorkingDay` whose total is bounded by a statutory limit |
| statutory limit | `StatutoryLimitMinutes` | type (`int64`) | `agent` | Maximum legal daily working time |
| working-day constraint | `WorkingDayConstraint` | type | `agent` | Pairs a `StatutoryLimitMinutes` with the jurisdiction/epoch it applies to |
| relay system | `RelaySchedule` | type | `agent` | Two alternating worker sets covering a 24-hour period |
| relay set | `RelaySet` | type | `agent` | One half of a `RelaySchedule`; has its own `WorkingDay` |
| day shift / night shift | `ShiftKind` | type (`string`) | `agent` | Enum: `"day"` / `"night"` |
| full-timer / half-timer | `TimerClass` | type (`string`) | `agent` | Enum: `"full"` / `"half"`; maps to adult vs. child workers |
| ten-hours bill | `FactoryAct` | type | `agent` | Named historical limit: year + `StatutoryLimitMinutes` |
| overwork (nibbling) | `Overwork` | type | `agent` | Minutes stolen above statutory limit; used in simulation ticks |
| physical maximum | `PhysicalMaxMinutes` | const `int64` | `agent` | 1440 (24 h × 60); hard ceiling |
| labourer's right | `ErrWorkingDayExceedsPhysicalMax` | sentinel error | `agent` | Returned when `TotalMinutes() > PhysicalMaxMinutes` |
| exceeds statutory limit | `ErrWorkingDayExceedsStatutoryLimit` | sentinel error | `agent` | Returned when limit is set and `TotalMinutes() > StatutoryLimitMinutes` |

## Fixtures

- **§1** `working-day I: A–B = 6 h, B–C = 1 h → total 7 h; rate 1/6 = 16.67%` → `WorkingDay{NecessaryLabourMinutes: 360, SurplusLabourMinutes: 60}.RateOfSurplusValue() ≈ 0.1667`
- **§1** `working-day II: A–B = 6 h, B–C = 3 h → total 9 h; rate 3/6 = 50%` → `WorkingDay{NecessaryLabourMinutes: 360, SurplusLabourMinutes: 180}.RateOfSurplusValue() == 0.50`
- **§1** `working-day III: A–B = 6 h, B–C = 6 h → total 12 h; rate 6/6 = 100%` → `WorkingDay{NecessaryLabourMinutes: 360, SurplusLabourMinutes: 360}.RateOfSurplusValue() == 1.00`
- **§1** `rate of surplus-value = 100% does not fix working-day length — 8, 10, 12 h all satisfy it` → constructing `WorkingDay{480, 480}`, `WorkingDay{600, 600}`, `WorkingDay{720, 720}` each returns `RateOfSurplusValue() == 1.00`
- **§2** `6 h necessary + 6 h surplus = 36 h surplus per week (6 days × 6 h)` → `WorkingDay{360, 360}.TotalMinutes() * 6 / 60 == 72` hours total per week; `SurplusLabourMinutes * 6 / 60 == 36`
- **§2** `Wallachian corvée: 56 corvée days / 140 working days → rate 56/84 = 66.67%` → `WorkingDay{NecessaryLabourMinutes: 84 * 480, SurplusLabourMinutes: 56 * 480}.RateOfSurplusValue() ≈ 0.6667`
- **§3** `Factory Inspector: 5 min/day × 300 working days = 1500 min per year` → `Overwork{MinutesPerDay: 5}.AnnualMinutes(300) == 1500`
- **§4** `relay system: two sets alternating so that total machinery absorption = 24 h` → `RelaySchedule` with two `RelaySet` entries whose `ShiftKind` differ and whose combined coverage spans `PhysicalMaxMinutes`
- **§6** `Factory Act 1833: children 9–13 limited to 8 h/day; young persons 13–18 limited to 12 h/day` → `FactoryAct{Year: 1833, ChildLimitMinutes: 480, YoungPersonLimitMinutes: 720}`
- **§6** `Factory Act 1847 (Ten Hours' Bill): young persons and women limited to 10 h/day from 1 May 1848` → `FactoryAct{Year: 1847, YoungPersonLimitMinutes: 600}`
- **§6** `Factory Act 1850: net 10.5 h weekdays + 7.5 h Saturday = 60 h/week` → `NormalWorkingDay` with `StatutoryLimitMinutes: 630` for weekdays

## Invariants

- `wd.TotalMinutes() == int64(wd.NecessaryLabourMinutes) + int64(wd.SurplusLabourMinutes)` [§1]
- `wd.TotalMinutes() <= PhysicalMaxMinutes` — working day cannot exceed 24 hours [§1]
- `wd.NecessaryLabourMinutes > 0` — minimum is reproduction time; zero necessary labour is not a valid capitalist working day [§1]
- `wd.RateOfSurplusValue() == float64(wd.SurplusLabourMinutes) / float64(wd.NecessaryLabourMinutes)` [§1]
- `Validate(wd, constraint) == ErrWorkingDayExceedsStatutoryLimit` when `wd.TotalMinutes() > constraint.StatutoryLimitMinutes` [§6]
- For any `RelaySchedule` rs: `rs.Sets[0].WorkingDay.TotalMinutes() + rs.Sets[1].WorkingDay.TotalMinutes() <= PhysicalMaxMinutes` [§4]

## Scope

### This chapter builds
- Services: `agent-service`
- New domain types:
  - `WorkingDay` — core struct: `NecessaryLabourMinutes`, `SurplusLabourMinutes`; methods `TotalMinutes()`, `RateOfSurplusValue()`, `Validate()`
  - `NecessaryLabourMinutes` — named `int64` for the AB segment
  - `SurplusLabourMinutes` — named `int64` for the BC segment
  - `StatutoryLimitMinutes` — named `int64` for a legal cap
  - `WorkingDayConstraint` — pairs a `StatutoryLimitMinutes` with a label string
  - `NormalWorkingDay` — a `WorkingDay` + `WorkingDayConstraint`; validated together
  - `FactoryAct` — historical act: `Year int`, `ChildLimitMinutes`, `YoungPersonLimitMinutes`, `AdultLimitMinutes`
  - `RelaySchedule` — two `RelaySet` entries covering complementary shifts
  - `RelaySet` — one shift: `ShiftKind`, `WorkingDay`, `WorkerIDs []ID`
  - `ShiftKind` — enum: `"day"` / `"night"`
  - `TimerClass` — enum: `"full"` / `"half"` (adult vs. child worker)
  - `Overwork` — `MinutesPerDay int64`; method `AnnualMinutes(workingDays int) int64`
  - Sentinel errors: `ErrWorkingDayExceedsPhysicalMax`, `ErrWorkingDayExceedsStatutoryLimit`
- New HTTP endpoints:
  - `POST /v1/working-days` — create and validate a `WorkingDay` (returns rate of surplus-value, total minutes, validation result)
  - `GET /v1/working-days/{id}` — retrieve a stored `WorkingDay`
  - `POST /v1/working-days/validate` — stateless: accepts `WorkingDay` + optional `StatutoryLimitMinutes`, returns validation errors and computed fields
  - `POST /v1/relay-schedules` — create a `RelaySchedule` assigning workers to shift sets
  - `GET /v1/relay-schedules/{id}` — retrieve a `RelaySchedule`
- New store files: extend `store.go` interface with `CreateWorkingDay`, `GetWorkingDay`, `CreateRelaySchedule`, `GetRelaySchedule`; add memory and MySQL implementations; add migration `00005_ch10_working_day.sql`
- React: "Ch. 10 — The Working-Day" panel; working-day builder with inputs for necessary/surplus minutes showing live rate of surplus-value; relay schedule visualiser showing day/night shift assignment

### Explicitly deferred to later chapters
- `cooperation` and division of labour within the working day — Ch. 11 (Cooperation)
- `relative surplus-value` (shortening necessary labour by raising productivity) — Ch. 12 (Division of Labour and Manufacture)
- `intensity of labour` as a further dimension of surplus extraction — Ch. 15 (Machinery and Modern Industry)
- Simulation-engine tick integration (advancing the economy clock by one `WorkingDay` at a time) — when simulation-engine is activated (Ch. 6-7 scope)
- Worker health/depreciation modelling (Marx's "premature exhaustion and death" argument) — no health state on `Agent` yet
- Multi-agent `WorkingDay` aggregation across a whole factory population — deferred to accumulation chapters (Ch. 23+)
