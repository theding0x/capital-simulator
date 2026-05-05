---
chapter: 08
title: "Constant Capital and Variable Capital"
status: proposed
primary_service: commodity-service
---

## Concepts → types

| Marx term | Go identifier | Kind | Package | Notes |
|---|---|---|---|---|
| constant capital | `ConstantCapital` | type | `commodity` | Portion of capital represented by means of production; value is preserved (transferred), not created |
| variable capital | `VariableCapital` | type | `commodity` | Portion of capital represented by labour-power; reproduces its own value and produces surplus-value |
| constant capital ID | `ConstantCapitalID` | type | `commodity` | 96-bit hex ID, `NewConstantCapitalID() ConstantCapitalID` |
| variable capital ID | `VariableCapitalID` | type | `commodity` | 96-bit hex ID, `NewVariableCapitalID() VariableCapitalID` |
| value transferred | `TransferredValue` | func | `commodity` | SNLT transferred from means of production to product; `TransferredValue(c ConstantCapital, fraction Quantity) LabourMinutes` |
| wear and tear / depreciation | `WearFraction` | type | `commodity` | `float64` in (0,1] representing the portion of instrument's value transferred per period |
| value added by labour | `NewValue` | func | `commodity` | Labour-power's value-creating act; `NewValue(labourTime LabourMinutes) LabourMinutes` |
| product value decomposition | `ProductValue` | type | `commodity` | `c + v + s`: constant-capital transferred value, variable-capital reproduced value, surplus-value |
| decompose product value | `DecomposeProductValue` | func | `commodity` | Returns `ProductValue` from `ConstantCapital`, `VariableCapital`, and surplus labour-time |
| instruments of labour (machines) | `Instrument` | type | `commodity` | A means of production that transfers value by fractions (depreciation) over its service life |
| raw material / auxiliary material | `RawMaterial` | type | `commodity` | A means of production consumed wholly in one production cycle; transfers full value at once |
| capital composition c:v | `CapitalComposition` | type | `commodity` | Ratio `ConstantValue : VariableValue`; struct with `Constant` and `Variable LabourMinutes` |
| surplus-value | `SurplusValue` | type | `commodity` | `LabourMinutes` alias to distinguish surplus from total product value |

## Fixtures

- **§1** `"By the 6 hours' spinning, the value of the raw material preserved and transferred to the product is six times as great as before, although the new value added by the labour of the spinner to each pound of the very same raw material is one-sixth what it was formerly."` → `TransferredValue` scales with quantity processed; `NewValue` scales with time worked, not quantity.

- **§1** `"Suppose its use-value in the labour-process to last only six days. Then, on the average, it loses each day one-sixth of its use-value, and therefore parts with one-sixth of its value to the daily product."` → `Instrument{ServiceLifeDays: 6}` yields `WearFraction(1.0/6.0)` per day; `TransferredValue` for one day = `instrument.OriginalValue / 6`.

- **§1** `"the 36 lbs. of cotton absorb only the same amount of labour as formerly did the 6 lbs."` → `ConstantCapital` value is fixed by the SNLT at time of production, not by the productivity of the labour that processes it.

- **§1** `"Suppose a machine to be worth £1,000, and to wear out in 1,000 days. Then one thousandth part of the value of the machine is daily transferred to the day's product."` → `Instrument{OriginalValue: 100_000 /* pence */, ServiceLifeDays: 1000}` → daily `TransferredValue` = `100` pence (= 1 LabourMinutes equivalent at prevailing SNLT).

- **§2** `"That part of capital ... which is represented by the means of production ... does not, in the process of production, undergo any quantitative alteration of value. I therefore call it the constant part of capital, or, more shortly, constant capital."` → `ConstantCapital.TransferredValue(qty) == originalValue` (no increment); `DecomposeProductValue.Constant == sum of all transferred values`.

- **§2** `"That part of capital, represented by labour-power, does, in the process of production, undergo an alteration of value. It both reproduces the equivalent of its own value, and also produces an excess, a surplus-value ... I therefore call it the variable part of capital, or, shortly, variable capital."` → `DecomposeProductValue.Variable + DecomposeProductValue.Surplus == totalLabourTime * SNLT_wage_rate`.

- **§2** `"The surplus of the total value of the product, over the sum of the values of its constituent factors, is the surplus of the expanded capital over the capital originally advanced."` → `ProductValue.Total() == c + v + s`; `SurplusValue == ProductValue.Total() - (constant_advanced + variable_advanced)`.

## Invariants

- `productValue.Constant == sum(instrument.TransferredValue(fraction) for each instrument) + sum(rawMaterial.Value(qtyConsumed) for each rawMaterial)` [§1–2]
- `productValue.Total() == productValue.Constant + productValue.Variable + productValue.Surplus` [§2]
- `ConstantCapital.transferredValue <= ConstantCapital.originalValue` — means of production cannot add more value than they possessed before entering the process [§1]
- `instrument.TransferredValue(fraction) == LabourMinutes(math.Round(float64(instrument.OriginalValue) * float64(fraction)))` where `fraction = 1.0 / float64(instrument.ServiceLifeDays)` [§1]
- For any change in `ConstantCapital` value (e.g. cotton price doubles), `NewValue(labourTime)` is unaffected — the two magnitudes are independent [§1]

## Scope

### This chapter builds
- Services: `commodity-service`
- New domain types:
  - `ConstantCapitalID` — 96-bit hex ID for constant-capital records
  - `VariableCapitalID` — 96-bit hex ID for variable-capital records
  - `ConstantCapital` — means of production with `OriginalValue LabourMinutes`, `Kind` (`"instrument"` / `"raw_material"` / `"auxiliary"`), and `ServiceLifeDays int64` (0 = consumed in one cycle)
  - `VariableCapital` — labour-power with `WageValue LabourMinutes` (value of labour-power) and `WorkingDay LabourMinutes` (actual hours worked)
  - `Instrument` — convenience alias / constructor for `ConstantCapital` where `ServiceLifeDays > 0`
  - `RawMaterial` — convenience alias / constructor for `ConstantCapital` where `ServiceLifeDays == 0`
  - `WearFraction` — `float64` fraction of an instrument's value transferred per period
  - `ProductValue` — struct `{ Constant, Variable, Surplus LabourMinutes }` with `Total() LabourMinutes` method
  - `CapitalComposition` — struct `{ Constant, Variable LabourMinutes }` with `Ratio() float64` method
- New HTTP endpoints:
  - `POST /v1/capital/decompose` — accepts `ConstantCapital[]`, `VariableCapital`, returns `ProductValue`
  - `GET /v1/capital/composition` — accepts query params `constant_value` and `variable_value`, returns `CapitalComposition` and ratio
- React: add "Ch. 08 — Constant & Variable Capital" panel; inputs for means-of-production list (kind, original value, service life) and labour-power (wage value, working day); output table showing c / v / s decomposition and capital composition ratio

### Explicitly deferred to later chapters
- Rate of surplus-value (`s/v`) — Ch. 9 introduces this as a distinct ratio and formula; Ch. 8 only names surplus-value
- Rate of profit (`s/(c+v)`) — Ch. 10+ (working day, accumulation context)
- Organic composition of capital as a driver of falling rate of profit — Ch. 25 (general law of capitalist accumulation)
- Production of relative surplus-value (lengthening the working day vs. intensification) — Ch. 10–16
- Simulation-engine tick integration (advancing the clock, periodic depreciation deductions) — deferred until simulation-engine is activated (Ch. 6-7 work completes)
- Moral depreciation (machinery losing value due to new invention) — mentioned in §2 footnote; full treatment in Ch. 15 (machinery and modern industry)
