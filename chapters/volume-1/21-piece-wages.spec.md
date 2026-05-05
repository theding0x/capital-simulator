---
chapter: 21
title: "Piece-Wages"
status: proposed
primary_service: agent-service
---

## Concepts → types

| Marx term | Go identifier | Kind | Package | Notes |
|---|---|---|---|---|
| piece-wage | `PieceWage` | type | `agent` | `AgentID AgentID; PricePence int64; NormalOutput int64; QualityThreshold float64` — price per accepted piece |
| piece price | `PiecePrice` | type | `agent` | `Pence int64` — `ComputePiecePrice(dailyWage, normalOutput)` |
| piece value | `PieceValue` | type | `agent` | `Pence int64` — `ComputePieceValue(dayValueProduct, normalOutput)`; always > piece price |
| normal output | `NormalOutput` | type | `agent` | `int64` — number of pieces socially expected in one working day |
| quality outcome | `QualityOutcome` | type | `agent` | `string` — `"accepted"` or `"rejected"`; rejected pieces earn no pay |
| sub-contract | `SubContract` | type | `agent` | `HeadLabourerID AgentID; AssistantIDs []AgentID; PieceRatePence int64; AssistantRatePence int64` — sweating: head labourer keeps spread |
| compute piece price | `ComputePiecePrice` | func | `agent` | `func(dailyWage Pence, normalOutput int64) Pence` — `dailyWage / normalOutput` |
| compute piece value | `ComputePieceValue` | func | `agent` | `func(dayValueProduct Pence, normalOutput int64) Pence` — `dayValueProduct / normalOutput`; exceeds piece price because value > wage |
| piece session | `PieceSession` | type | `agent` | `WorkingSession` extended with `PiecesProduced int64; QualityOutcome QualityOutcome`; compute actual earnings |
| compute session earnings | `ComputePieceEarnings` | func | `agent` | `func(s PieceSession, pw PieceWage) Pence` — `acceptedPieces * pw.PricePence` |

## Fixtures

- **§ 12-hour day, 24 pieces** `"ordinary working-day 12 hours ... 24 pieces; value of 24 pieces = 6s.; labourer receives 1½d. per piece; earns 3s. in 12 hours"` → `ComputePiecePrice(Pence(36), 24) == 1` (rounded; exact is 36/24=1.5d → use int pence: `ComputePiecePrice(36,24) == 1`; note 24×1=24 < 36, so store Pence as halfpence or keep numerator/denominator; alternatively use 2-halfpenny representation)
  - Simpler: let pence be in halfpennies (halfpenny = smallest unit): `ComputePiecePrice(72_halfpennies, 24) == 3_halfpennies`; `3×24 == 72` check passes
- **§ doubled productivity, 48 pieces** `"doubled productiveness ... 48 pieces; piece-wage falls from 1½d. to ¾d.; 48×¾d. = 3s."` → `ComputePiecePrice(72_halfpennies, 48) == 1` halfpenny (≈¾d.); `1×48 == 48_halfpennies == 24 pence == 2s.` — wait, Marx says 3s.; this means price per piece rounds to maintain total: 72/48 = 1.5 halfpennies per piece → same fractional issue; use `int64` farthings (¼d.) unit: `ComputePiecePrice(144_farthings, 24) == 6_farthings` (1½d.); `6×24==144`; `ComputePiecePrice(144_farthings, 48) == 3_farthings` (¾d.); `3×48==144`
- **§ piece value > piece price** `"value of 24 pieces = 6s."` while wage for 24 pieces = 3s. → `ComputePieceValue(144_farthings_value, 24) == 6` farthings (= 1½d. per piece value); `ComputePiecePrice(72_farthings_wage, 24) == 3` farthings; value > price: `6 > 3`
- **§ sub-contract sweating** `"gain of ... middlemen comes entirely from the difference between the labour-price which the capitalist pays, and the part ... they actually allow to reach the labourer"` → `SubContract{HeadLabourerID:X, AssistantIDs:[A,B], PieceRatePence:6, AssistantRatePence:4}`; surplus to head = `6 - 4 == 2` pence per piece

## Invariants

- `ComputePiecePrice(dailyWage, normalOutput) * normalOutput == dailyWage` — a worker producing exactly normal output earns exactly the daily wage (integer arithmetic using farthings to avoid rounding)
- `ComputePieceValue(dayValueProduct, normalOutput) > ComputePiecePrice(dailyWage, normalOutput)` whenever `dayValueProduct > dailyWage` — piece value always exceeds piece price because the day value product includes surplus-value
- `ComputePieceEarnings(s, pw).Pence == countAccepted(s) * pw.PricePence` — rejected pieces earn zero; quality control is enforced by the wage form itself
- `SubContract` spread: head labourer's net = `(PieceRatePence - AssistantRatePence) * piecesProduced` — positive whenever PieceRatePence > AssistantRatePence

## Scope

### This chapter builds
- Services: `agent-service`
- New domain types:
  - `PieceWage` — `AgentID AgentID; PricePence int64; NormalOutput int64`; persisted wage contract
  - `QualityOutcome` — `"accepted"` | `"rejected"`
  - `PieceSession` — extends working session concept with `PiecesProduced int64; QualityOutcome QualityOutcome`
  - `SubContract` — `HeadLabourerID AgentID; AssistantIDs []AgentID; PieceRatePence int64; AssistantRatePence int64`; models the sweating system
- New functions:
  - `ComputePiecePrice(dailyWage, normalOutput int64) int64` — exact integer division; use smallest monetary unit (farthings) to avoid fractional pence
  - `ComputePieceValue(dayValueProduct, normalOutput int64) int64` — analogous; always exceeds price
  - `ComputePieceEarnings(s PieceSession, pw PieceWage) int64` — multiply accepted pieces by piece price
- New HTTP endpoints:
  - `POST /v1/agents/{id}/piece-wages` — register a piece-wage contract
  - `GET /v1/agents/{id}/piece-wages` — retrieve active piece-wage contract
  - `POST /v1/piece-price` — stateless; compute piece price and piece value from inputs
  - `POST /v1/sub-contracts` — create a sub-contract record
  - `GET /v1/sub-contracts/{id}` — retrieve sub-contract with earnings summary
- Migration: `00007_ch21_piece_wages.sql` — `piece_wages` and `sub_contracts` tables
- React: add a "Ch. 21 — Piece-Wages" panel; inputs for daily wage, normal output, pieces produced, and quality pass/fail; display piece price, piece value, actual earnings, and a comparison showing value extracted vs. wage paid

### Explicitly deferred to later chapters
- Dynamic piece-price reduction when productivity rises — described in the chapter but requires a simulation tick with changing technology; deferred to accumulation chapters
- The "domestic labour" system built on piece-wages — mentioned as a consequence but requires a household/family model not yet in scope
- Trade union resistance to piece-price cuts — historical examples given; out of scope for this simulation
