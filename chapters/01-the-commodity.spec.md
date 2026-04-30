---
chapter: 01
title: "The Commodity"
status: implemented
primary_service: commodity-service
---

## Concepts → types

| Marx term | Go identifier | Kind | Package | Notes |
|---|---|---|---|---|
| Commodity | `Commodity` | struct | `commodity` | Central domain object; carries use-value, concrete labour, and SNLT |
| Use-value | `UseValue` | struct | `commodity` | Description + unit; satisfies human wants |
| Exchange-value | `ExchangeRatio` | struct | `commodity` | Quantitative proportion between two commodities |
| Value (magnitude) | `LabourMinutes` | named int64 | `commodity` | Canonical unit for value-magnitude; duration of abstract labour |
| Socially necessary labour-time (SNLT) | `SNLTPerUnit` | field on `Commodity` | `commodity` | `LabourMinutes` per unit; set at creation; recalculated by `ProductivityChange` |
| Abstract labour | `AsAbstractLabour` | func | `commodity` | Reduces concrete labour to a `LabourMinutes` magnitude |
| Concrete labour | `ConcreteLabour` | struct | `commodity` | Kind + description; qualitatively distinct; does not directly express value |
| Value of a quantity | `Commodity.Value` | method | `commodity` | Returns `LabourMinutes` for a given `Quantity` |
| Productivity change | `ProductivityChange` | struct | `commodity` | Applies a factor to `SNLTPerUnit`; models the inverse law from §1 |
| Simple (accidental) form of value | `SimpleForm` | struct | `commodity` | x A = y B; relative + equivalent poles |
| Expanded form of value | `ExpandedForm` | struct | `commodity` | z A = u B or = v C … |
| General form of value | `GeneralForm` | struct | `commodity` | All commodities express value in one equivalent |
| Money-form | `MoneyForm` | struct | `commodity` | General form where the equivalent is socially fixed as the money-commodity |
| Value-form kind | `ValueFormKind` | named string + consts | `commodity` | `KindSimple`, `KindExpanded`, `KindGeneral`, `KindMoney` |
| Commodity fetishism / social relations | `SocialRelations` | struct | `commodity` | Reveals congealed labour-time and labour-relations hidden behind exchange |
| Labour relation (pairwise) | `LabourRelation` | struct | `commodity` | Subject ↔ counterpart; the social relation expressed as a thing-relation |
| Quantity | `Quantity` | named float64 | `commodity` | Number of units of a commodity entering an exchange |

## Fixtures

- **§1** `"1 quarter corn = x cwt. iron"` → `ExchangeRatioBetween(corn, iron, 1)` returns `CommonValue == corn.Value(1)` and `QuoteQty == corn.Value(1) / iron.SNLTPerUnit`
- **§1** `"power-looms reduced by one-half the labour required to weave"` → `ProductivityChange{Factor: 2.0}.Apply(handLoom)` yields `SNLTPerUnit == handLoom.SNLTPerUnit / 2`
- **§2** `"if 10 yards of linen = W, the coat = 2W"` → `coat.SNLTPerUnit == 2 * linen.SNLTPerUnit` in the standard two-commodity fixture
- **§2** `"20 yds of linen must have the same value as one coat"` → `linen.Value(20) == coat.Value(1)` when coat SNLT is double linen SNLT
- **§3A** `"20 yards of linen = 1 coat"` → `SimpleFormOf(linen, coat)` returns `SimpleForm` where `RelativeQty == 20`, `EquivalentQty == 1`, `CommonValue == linen.Value(20)`
- **§3B** `"20 yards of linen = 1 coat or = 10 lbs. tea or = 40 lbs. coffee"` → `ExpandedFormOf(linen, population)` returns one `SimpleForm` per equivalent, all sharing equal `CommonValue`
- **§3C** `"1 coat = … 20 yards of linen"` → `GeneralFormOf(linen, population)` returns `GeneralForm` where every `Relatives[i].CommonValue` is equal
- **§3D** `"2 oz. gold = … all commodities"` → `MoneyFormOf(gold, population)` returns `MoneyForm` with `MoneyCommodity == gold` and `GeneralForm.Relatives` covering the full population
- **§4** `"a coat is worth 20 yards of linen … a direct reflection of social labour-time"` → `SocialRelationsOf(coat, []Commodity{linen})` has `LabourRelations[0].LabourTime == linen.Value(20)` and `LabourPerUnit == coat.SNLTPerUnit`

## Invariants

- `c.Value(qty) == LabourMinutes(math.Round(float64(c.SNLTPerUnit) * float64(qty)))` for all qty ≥ 0 [§1]
- `ExchangeRatioBetween(a, b, q).CommonValue == a.Value(q)` for any valid pair [§1]
- `ProductivityChange{Factor: f}.Apply(c)` → result `SNLTPerUnit == LabourMinutes(math.Round(float64(c.SNLTPerUnit) / f))` for all f > 0 [§1]
- For any valid `SimpleForm sf`: `sf.CommonValue == sf.Relative.Value(sf.RelativeQty)` and `sf.CommonValue == sf.Equivalent.Value(sf.EquivalentQty)` [§3A]
- For any valid `ExpandedForm ef`: all `ef.Equivalents[i].CommonValue` are equal and equal `ef.CommonValue` [§3B]
- `GeneralFormOf(e, pop)` → every `sf` in `Relatives` satisfies `sf.Equivalent.ID == e.ID` [§3C]
- `MoneyFormOf(m, pop).MoneyCommodity.ID == MoneyFormOf(m, pop).GeneralForm.Equivalent.ID` [§3D]
- `SocialRelationsOf(s, pop).LabourPerUnit == s.SNLTPerUnit` regardless of population size [§4]

## Scope

### This chapter builds
- Services: `commodity-service` (port 8081) — full CRUD plus value, value-form, social-relations, and exchange-ratio endpoints
- New domain types: `Commodity`, `UseValue`, `ConcreteLabour`, `LabourMinutes`, `Quantity`, `ExchangeRatio`, `ProductivityChange`, `SimpleForm`, `ExpandedForm`, `GeneralForm`, `MoneyForm`, `ValueFormKind`, `SocialRelations`, `LabourRelation`
- New HTTP endpoints: `POST /v1/commodities`, `GET /v1/commodities`, `GET /v1/commodities/{id}`, `PATCH /v1/commodities/{id}`, `DELETE /v1/commodities/{id}`, `POST /v1/commodities/{id}/value`, `GET /v1/commodities/{id}/value-form`, `GET /v1/commodities/{id}/social-relations`, `POST /v1/exchange-ratio`
- React: commodity CRUD panel, exchange-ratio calculator, Reveal toggle (fetishism critique) in `App.tsx`; `types.ts` mirrors all Go structs; `api.ts` client for all nine endpoints

### Explicitly deferred to later chapters
- Money as a medium of circulation and means of payment (Ch. 2–3; `market-service`)
- The transformation of money into capital; buying/selling of labour-power (Ch. 4+; `agent-service`)
- Absolute and relative surplus-value (Ch. 7–15)
- The general formula for capital M–C–M′ (Ch. 4)
