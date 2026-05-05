---
chapter: 22
title: "National Differences of Wages"
status: proposed
primary_service: agent-service
---

## Concepts → types

| Marx term | Go identifier | Kind | Package | Notes |
|---|---|---|---|---|
| country code | `CountryCode` | type | `agent` | `string` — ISO 3166-1 alpha-2, e.g. `"GB"`, `"FR"`, `"DE"` |
| national intensity | `NationalIntensity` | type | `agent` | `CountryCode CountryCode; Factor float64` — intensity of labour relative to international average (1.0 = average) |
| day wage | `DayWage` | type | `agent` | `CountryCode CountryCode; NominalPence int64; WorkingDayMinutes int64` — as paid in the domestic market |
| standardised wage | `StandardisedWage` | type | `agent` | `CountryCode CountryCode; Amount int64` — day wage reduced to a uniform working day for comparison |
| relative labour price | `RelativeLabourPrice` | type | `agent` | `CountryCode CountryCode; Ratio float64` — price of labour relative to value produced; inversely related to nominal wage in high-productivity nations |
| spindle ratio | `SpindleRatio` | type | `agent` | `CountryCode CountryCode; SpindlesPerWorker int64` — proxy for productivity/intensity; from Redgrave's 1866 data |
| wage comparison | `WageComparison` | type | `agent` | `Countries []CountryCode; DayWages []DayWage; StandardisedWages []StandardisedWage; RelativePrices []RelativeLabourPrice` |
| standardise wage | `StandardiseWage` | func | `agent` | `func(w DayWage, referenceDayMinutes int64) StandardisedWage` — reduce to common working-day length |
| compute relative price | `ComputeRelativePrice` | func | `agent` | `func(w DayWage, ni NationalIntensity) RelativeLabourPrice` — price of labour as fraction of value produced; higher intensity → lower relative price |

## Fixtures

- **§ spindle table** Redgrave's 1866 data: France 1 person per 14 spindles; Russia 28; Prussia 37; Bavaria 46; Austria 49; Belgium 50; Saxony 50; Switzerland 55; Smaller German states 55; Great Britain 74 → `SpindleRatio{CountryCode:"GB", SpindlesPerWorker:74}`, `SpindleRatio{CountryCode:"FR", SpindlesPerWorker:14}`, etc.
- **§ standardised wage** `"average day-wage for the same trades, in different countries, to a uniform working-day"` → if GB works 600 minutes and FR works 720 minutes at the same nominal wage, `StandardiseWage(DayWage{"FR", 36, 720}, 600) == StandardisedWage{"FR", 30}` — FR nominal looks equal but standardised is lower
- **§ relative price inversion** `"in England wages are virtually lower to the capitalist, though higher to the operative than on the Continent"` → `ComputeRelativePrice(DayWage{"GB", 48, 600}, NationalIntensity{"GB", 2.0}).Ratio < ComputeRelativePrice(DayWage{"DE", 24, 720}, NationalIntensity{"DE", 1.0}).Ratio` — GB nominal wage is higher, but relative labour price is lower because intensity is higher
- **§ factory average** `"average of spindles per factory, England: 12,600; France: 1,500; Prussia: 1,500"` — stored as reference data; not modelled as a computed quantity

## Invariants

- `StandardiseWage(w, ref).Amount == w.NominalPence * ref / w.WorkingDayMinutes` — standardisation is a proportional reduction to the reference day length
- `ComputeRelativePrice(w, ni).Ratio` is inversely related to `ni.Factor`: higher national intensity → lower relative price of labour despite potentially higher nominal wage
- `SpindleRatio.SpindlesPerWorker > 0` — must be positive; GB 74 > all continental entries in Redgrave's table

## Scope

### This chapter builds
- Services: `agent-service`
- New domain types:
  - `CountryCode` — `string`; country identifier
  - `NationalIntensity` — `CountryCode; Factor float64`; intensity of national labour relative to international average
  - `DayWage` — `CountryCode; NominalPence int64; WorkingDayMinutes int64`; wage as paid domestically
  - `StandardisedWage` — `CountryCode; Amount int64`; wage reduced to common working-day basis
  - `RelativeLabourPrice` — `CountryCode; Ratio float64`; price of labour relative to value it creates
  - `SpindleRatio` — `CountryCode; SpindlesPerWorker int64`; productivity proxy from historical data
  - `WageComparison` — aggregate output grouping all per-country metrics
- New functions:
  - `StandardiseWage(w DayWage, referenceDayMinutes int64) StandardisedWage` — pure; proportional scaling
  - `ComputeRelativePrice(w DayWage, ni NationalIntensity) RelativeLabourPrice` — captures the inversion
- New HTTP endpoints:
  - `POST /v1/intensities` — register or update a national intensity record
  - `GET /v1/intensities` — list all national intensities
  - `POST /v1/wages` — register a day-wage record for a country
  - `GET /v1/wages/{country}/standardised` — return standardised wage for a given reference day length (query param)
  - `GET /v1/comparisons` — return a full `WageComparison` across all registered countries
- Migration: `00008_ch22_national_wages.sql` — `national_intensities`, `day_wages`, `spindle_ratios` tables
- React: add a "Ch. 22 — National Wages" panel; table of countries with nominal wage, standardised wage, relative price, and spindle ratio; highlight the UK paradox (highest nominal, lowest relative price)

### Explicitly deferred to later chapters
- World market value and international law of value — mentioned in the chapter ("the more productive national labour reckons as the more intense on the world market") but requires a multi-country market simulation not yet in scope
- Tax deductions and "State expenses" (Carey's argument) — referenced but outside the simulation's economic model
- Absolute international price levels — the chapter focuses on ratios; absolute currency conversion is deferred
