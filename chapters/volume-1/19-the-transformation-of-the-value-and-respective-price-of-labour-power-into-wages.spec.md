---
chapter: 19
title: "The Transformation of the Value (and Respective Price) of Labour-Power into Wages"
status: proposed
primary_service: agent-service
---

## Concepts → types

| Marx term | Go identifier | Kind | Package | Notes |
|---|---|---|---|---|
| value of labour-power | `LabourPowerValue` | type | `agent` | `DailyPence int64; NecessaryMinutes LabourMinutes` — what reproduction actually costs |
| wage | `Wage` | type | `agent` | `DailyPence int64; WorkingDayHours int64` — money paid, which appears to cover all hours |
| wage appearance | `WageAppearance` | type | `agent` | `PaidHours int64 = WorkingDayHours; UnpaidHours int64 = 0` — ideological surface: wages look like pay for all labour |
| labour decomposition | `LabourDecomposition` | type | `agent` | `PaidMinutes LabourMinutes; UnpaidMinutes LabourMinutes` — true split beneath the wage form |
| hourly wage | `HourlyWage` | func | `agent` | `func(w Wage) Pence` — `DailyPence / WorkingDayHours`; used by workers and economists to reason about price of labour |
| decompose | `Decompose` | func | `agent` | `func(w Wage, lpv LabourPowerValue) LabourDecomposition` — reveals paid vs. unpaid minutes |
| appearance | `Appearance` | func | `agent` | `func(w Wage) WageAppearance` — models the ideological inversion: all hours appear paid |
| wage form | `WageForm` | type | `agent` | Holds `AgentID AgentID; Wage Wage; LabourPowerValue LabourPowerValue`; persisted record |

## Fixtures

- **§ core example** `"daily value of labour-power is 3s., the value of the product of 6 working-hours; working-day is 12 hours"` → `LabourPowerValue{DailyPence:36, NecessaryMinutes:360}; Wage{DailyPence:36, WorkingDayHours:12}`
- **§ hourly wage** from above: `HourlyWage(Wage{36, 12}) == 3` pence per hour
- **§ decompose** `Decompose(Wage{36,12}, LabourPowerValue{36,360}) == LabourDecomposition{PaidMinutes:360, UnpaidMinutes:360}` — 6 hours paid, 6 hours unpaid despite 12-hour wage payment
- **§ appearance** `Appearance(Wage{36,12}) == WageAppearance{PaidHours:12, UnpaidHours:0}` — ideological form shows zero unpaid hours
- **§ price falls** `"value of labour-power may vary from 3 to 4 shillings or from 3 to 2 shillings"` → `Wage{DailyPence:24, WorkingDayHours:12}` with same `LabourPowerValue{36,360}` yields `HourlyWage == 2` pence; real wage falls while nominal form conceals it

## Invariants

- `Decompose(w, lpv).PaidMinutes == lpv.NecessaryMinutes` — the paid minutes equal the necessary labour reproducing the value of labour-power, not the full working day
- `Decompose(w, lpv).PaidMinutes + Decompose(w, lpv).UnpaidMinutes == int64(w.WorkingDayHours) * 60` — decomposition covers the full working day
- `Appearance(w).PaidHours == w.WorkingDayHours && Appearance(w).UnpaidHours == 0` — wage form always presents all hours as paid; this is the ideological inversion Ch.19 critiques
- `HourlyWage(w) == w.DailyPence / w.WorkingDayHours` — the price of labour as it appears in everyday accounting

## Scope

### This chapter builds
- Services: `agent-service`
- New domain types:
  - `LabourPowerValue` — `DailyPence int64; NecessaryMinutes LabourMinutes`; what reproduction costs vs. how long it takes
  - `Wage` — `DailyPence int64; WorkingDayHours int64`; the money form paid per day
  - `WageAppearance` — surface-level ideological representation: PaidHours = full day, UnpaidHours = 0
  - `LabourDecomposition` — true split: `PaidMinutes`, `UnpaidMinutes`
  - `WageForm` — persisted record linking agent, wage, and labour-power value
- New functions:
  - `HourlyWage(w Wage) Pence` — daily pence divided by working day hours
  - `Decompose(w Wage, lpv LabourPowerValue) LabourDecomposition` — reveals the true paid/unpaid split
  - `Appearance(w Wage) WageAppearance` — models the ideological inversion
- New HTTP endpoints:
  - `POST /v1/wage-forms` — create a WageForm for an agent
  - `GET /v1/wage-forms/{agentID}` — retrieve wage form with decomposition
- Migration: `00005_ch19_wage_forms.sql` — `wage_forms` table
- React: add a "Ch. 19 — Wage Form" panel; inputs for daily wage pence, working day hours, and necessary minutes; display both the ideological appearance (all paid) and the true decomposition (paid vs. unpaid minutes)

### Explicitly deferred to later chapters
- Time-wages and piece-wages as specific wage forms — deferred to Ch. 20 and Ch. 21 respectively
- National differences in wages — deferred to Ch. 22
- Market-price fluctuations of labour-power above or below its value — mentioned but not modelled; belongs to a future market-service extension
