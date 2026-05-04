---
chapter: 05
title: "Contradictions in the General Formula of Capital"
status: proposed
primary_service: agent-service
---

## Concepts → types

| Marx term | Go identifier | Kind | Package | Notes |
|---|---|---|---|---|
| Money-owner / embryo capitalist | `Agent` | type | `agent` | Bears money-capital M before the circuit begins |
| Agent identity | `AgentID` | type | `agent` | 96-bit hex string; constructed with `NewAgentID()` |
| Capitalist agent | `AgentKindCapitalist` | const | `agent` | Enumerated `AgentKind` value |
| Simple commodity owner | `AgentKindOwner` | const | `agent` | Buys and sells equivalents; cannot originate surplus-value |
| Money in circuit | `MoneyAmount int64` | type | `agent` | Expressed in the smallest monetary unit (pence/cents) |
| Value magnitude | `LabourMinutes int64` | type | `agent` | Canonical unit carried from Ch. 1; no new type needed |
| Circuit M-C-M′ | `Circuit` | type | `agent` | Holds `M`, `C` (commodity ID or description), `MPrime` fields |
| Surplus-value (Δm) | `SurplusValue` | func | `agent` | `MPrime - M`; returns 0 for equivalent exchange, negative for loss |
| Exchange of equivalents | `ExchangeEquivalents` | func | `agent` | Models bilateral swap; asserts total value is conserved |
| Non-equivalent exchange | `ExchangeNonEquivalents` | func | `agent` | Redistributes value between A and B without creating new value |
| Total social value | `TotalValue` | func | `agent` | Sum of all agent holdings; invariant across any exchange |
| Seller role | `RoleSeller` | const | `agent` | Same agent becomes buyer in the next step |
| Buyer role | `RoleBuyer` | const | `agent` | Mirrors `RoleSeller`; agents cycle roles |
| Merchants' capital | `MerchantsCapital` | type | `agent` | M-C-M′ through sphere of circulation only; surplus-value still unexplained |
| Money-lending / usurer's capital | `UsurersCapital` | type | `agent` | M-M′ (degenerate form); models `M` → `MPrime` without a commodity intermediary |

## Fixtures

Marx's own numbers and examples from the text. These become test case names and values verbatim.

- **§1** `"A man who has plenty of wine and no corn treats with a man who has plenty of corn and no wine; an exchange takes place between them of corn to the value of 50, for wine of the same value."` → `ExchangeEquivalents(wine=50, corn=50)` yields `SurplusValue == 0` for both parties; `TotalValue` before and after is 100.

- **§2** `"Suppose then, that by some inexplicable privilege, the seller is enabled to sell his commodities above their value, what is worth 100 for 110, in which case the price is nominally raised 10%."` → `ExchangeNonEquivalents(sellerValue=100, price=110)`: seller gains 10, buyer loses 10, `TotalValue` remains 210 across both parties.

- **§3** `"A sells wine worth £40 to B, and obtains from him in exchange corn to the value of £50. A has converted his £40 into £50 ... The value in circulation has not increased by one iota ... total value of £90."` → Before: A holds 40, B holds 50, total 90. After: A holds 50, B holds 40, total still 90. `SurplusValue` from circuit perspective is 0.

- **§4** `"If equivalents are exchanged, no surplus-value results, and if non-equivalents are exchanged, still no surplus-value. Circulation, or the exchange of commodities, begets no value."` → Property test: for any pair `(x, y)` exchanged (whether `x == y` or `x != y`), `TotalValue(before) == TotalValue(after)`.

- **§5** `"general and nominal rise of prices has the same effect as if the values had been expressed in weight of silver instead of in weight of gold"` → Scaling all `MoneyAmount` values by factor k leaves `SurplusValue == 0` and all value ratios unchanged.

- **§6** `"M-M, money exchanged for more money, a form that is incompatible with the nature of money"` → `UsurersCapital{M: 100, MPrime: 110}`: `SurplusValue()` returns 10, but the source of that 10 cannot be located within the circuit — verified by absence of any intermediate commodity field.

## Invariants

Mathematical or logical laws the chapter establishes that tests must enforce:

- `TotalValue(allAgents, before) == TotalValue(allAgents, after)` for any `ExchangeEquivalents` or `ExchangeNonEquivalents` call [§3, §18-footnote: "L'échange … ne change rien … à la somme des valeurs sociales"]
- `SurplusValue(circuit) == circuit.MPrime - circuit.M` [definition of Δm; §49]
- `ExchangeEquivalents(x, x).SurplusValue() == 0` [§1, §3, §7-footnote: "Exchange … not a means of self-enrichment"]
- `ExchangeNonEquivalents(a, b).SurplusValueA() + ExchangeNonEquivalents(a, b).SurplusValueB() == 0` (zero-sum redistribution) [§3, §31]
- `MerchantsCapital.SurplusValue()` is non-zero only if one counterparty is systematically deceived — the model must flag `origin == "redistribution"` not `"creation"` [§37]
- A `Circuit` where every sub-exchange is at value must produce `MPrime == M`; `SurplusValue == 0` [§33: "Circulation … begets no value"]

## Scope

### This chapter builds
- **Services:** agent-service
- **New domain types:**
  - `AgentID` — 96-bit hex identity for an economic agent, with `NewAgentID()` constructor
  - `AgentKind` — enum discriminating `Capitalist | Owner`
  - `Agent` — struct carrying `AgentID`, `AgentKind`, `MoneyAmount`, `LabourMinutes`
  - `Circuit` — value-object representing M→C→M′; exposes `SurplusValue() int64`
  - `MerchantsCapital` — Circuit subtype operating purely within circulation
  - `UsurersCapital` — degenerate M→M′ circuit (no commodity intermediate)
  - `ExchangeResult` — result of bilateral exchange: value held by each party, zero-sum redistribution proof
- **New HTTP endpoints:**
  - `POST /v1/agents` — create a new agent (capitalist or owner) with initial money holdings
  - `GET /v1/agents` — list all agents and their current `MoneyAmount`
  - `GET /v1/agents/{id}` — retrieve a single agent
  - `POST /v1/circuits` — record a M-C-M′ circuit attempt; returns `SurplusValue` and `origin` tag
  - `POST /v1/exchanges` — simulate a bilateral commodity exchange; returns `ExchangeResult` proving value conservation
- **React:** Add an "Agents" panel to `App.tsx` listing agents with their money holdings; add a "Circuit" form to submit M and M′ values and display the computed `SurplusValue` with an explanatory note that circulation cannot be its source.

### Explicitly deferred to later chapters
- Labour-power as a commodity — the resolution of the contradiction (where surplus-value actually comes from) is Ch. 6; this chapter only proves circulation *cannot* be the source.
- Worker agents (`AgentKindWorker`) — introduced in Ch. 6 when labour-power enters the market.
- Valorization process and the production of surplus-value — Ch. 7 (`simulation-engine`).
- Landowner agents and ground-rent — Volume III; out of scope for Volume I.
- Persistence layer (MySQL store for agents) — can be scaffolded now but fully exercised once agents participate in production circuits from Ch. 6 onward.
- `simulation-engine` integration — the engine coordinates multi-agent circuits; deferred until Ch. 7.
