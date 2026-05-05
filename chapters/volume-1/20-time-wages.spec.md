---
chapter: 20
title: "Time-Wages"
status: proposed
primary_service: agent-service
---

## Concepts → types

| Marx term | Go identifier | Kind | Package | Notes |
|---|---|---|---|---|
| daily value of labour-power | `DailyLabourPowerValue` | type | `agent` | `Pence int64` — monetary expression of necessary labour-time |
| price of the working-hour | `HourlyPriceOfLabour` | type | `agent` | `Numerator int64; Denominator int64` — rational fraction: daily value / working-day hours; kept exact to avoid rounding |
| working day | `WorkingDayHours` | type | `agent` | `Hours int64` — total hours worked in the period |
| nominal wage | `NominalWage` | type | `agent` | `Pence int64` — actual money received for the period |
| overtime hours | `OvertimeHours` | type | `agent` | `Hours int64` — hours beyond the normal working day |
| overtime rate | `OvertimeRatePence` | type | `agent` | `Pence int64` — extra pay per overtime hour (often > normal hourly rate) |
| wage period | `WagePeriod` | type | `agent` | `string` — `"daily"` or `"weekly"` |
| working session | `WorkingSession` | type | `agent` | Holds `AgentID AgentID; DailyLabourPowerValue DailyLabourPowerValue; WorkingDayHours WorkingDayHours; OvertimeHours OvertimeHours; OvertimeRatePence OvertimeRatePence; WagePeriod WagePeriod` |
| compute hourly price | `ComputeHourlyPrice` | func | `agent` | `func(dailyValue DailyLabourPowerValue, hours WorkingDayHours) HourlyPriceOfLabour` — rational fraction; never rounds |
| compute session wage | `ComputeSessionWage` | func | `agent` | `func(s WorkingSession) NominalWage` — total pay including overtime |

## Fixtures

- **§ unit measure** `"daily value of labour-power 3s., working-day 12 hours → price of 1 hour = 3/12 = 3d."` → `ComputeHourlyPrice(DailyLabourPowerValue{36}, WorkingDayHours{12}) == HourlyPriceOfLabour{Numerator:36, Denominator:12}` which simplifies to 3 pence per hour
- **§ 10-hour day** `"if the working-day is 10 hours ... price of working-hour is 3⅗d."` → `ComputeHourlyPrice(DailyLabourPowerValue{36}, WorkingDayHours{10}) == HourlyPriceOfLabour{36, 10}` (18/5 pence, not rounded)
- **§ 15-hour day** `"rises to 15 hours"` → `ComputeHourlyPrice(DailyLabourPowerValue{36}, WorkingDayHours{15}) == HourlyPriceOfLabour{36, 15}` (12/5 pence)
- **§ nominal wage stable, price falls** with 10-hour day at 3s. nominal → price is 3⅗d.; extend to 12 hours keeping 3s. nominal → price falls to 3d. — hourly price drops even though total pay is unchanged
- **§ overtime** extra pay at 4d./hour when normal is 3d. — capitalist still extracts unpaid labour at overtime rate: `4d. = value-product of 2/3 hour` vs. full hour worked; `ComputeSessionWage` reflects the nominal amount only

## Invariants

- `ComputeHourlyPrice(v, h).Numerator == v.Pence && ComputeHourlyPrice(v, h).Denominator == h.Hours` — stored as exact fraction, not a float; no rounding at this layer
- Nominal hourly price falls when working-day lengthens without a change in daily value: `ComputeHourlyPrice(v, WorkingDayHours{h2}).AsFloat() < ComputeHourlyPrice(v, WorkingDayHours{h1}).AsFloat()` when `h2 > h1`
- `ComputeSessionWage(s).Pence == (s.DailyLabourPowerValue.Pence / s.WorkingDayHours.Hours) * s.WorkingDayHours.Hours + s.OvertimeHours.Hours * s.OvertimeRatePence.Pence` — integer arithmetic; overtime paid at specified rate
- Overtime hours add to nominal earnings but do not change the hourly price of normal hours

## Scope

### This chapter builds
- Services: `agent-service`
- New domain types:
  - `DailyLabourPowerValue` — `Pence int64`; daily money value of the agent's labour-power
  - `HourlyPriceOfLabour` — `Numerator, Denominator int64`; exact rational fraction; add `AsFloat() float64` helper
  - `WorkingDayHours` — `Hours int64`; total hours in a given working period
  - `OvertimeHours` — `Hours int64`; extension beyond the normal day
  - `OvertimeRatePence` — `Pence int64`; per-hour pay for overtime
  - `WagePeriod` — `string`; `"daily"` or `"weekly"`
  - `WorkingSession` — composite record grouping all time-wage inputs; persisted
- New functions:
  - `ComputeHourlyPrice(v DailyLabourPowerValue, h WorkingDayHours) HourlyPriceOfLabour` — pure; exact fraction
  - `ComputeSessionWage(s WorkingSession) NominalWage` — total nominal earnings including overtime
- New HTTP endpoints:
  - `POST /v1/time-wages/hourly-price` — stateless; returns `HourlyPriceOfLabour` as `{numerator, denominator, as_float}`
  - `POST /v1/time-wages/sessions` — create a working session record
  - `GET /v1/time-wages/sessions/{id}` — retrieve session with computed nominal wage
- Migration: `00006_ch20_time_wage_sessions.sql` — `time_wage_sessions` table
- React: add a "Ch. 20 — Time-Wages" panel; inputs for daily value, normal hours, overtime hours, overtime rate; display hourly price (as fraction and decimal) and total nominal wage; highlight that a longer day lowers the price per hour even when total pay holds steady

### Explicitly deferred to later chapters
- Piece-wages — the other fundamental wage form; deferred to Ch. 21
- Legal limitation of the working day as a mechanism to fix the normal working day — discussed in Ch. 10 (working day) and touched here but no new model required
- Minimum-wage / subsistence floor enforcement — mentioned but requires a labour-market model not yet in scope
