# ADR 0001 — Value composition across the circuit: the LabourMinutes ↔ Pence bridge

- **Status:** Accepted (decision); implementation **Deferred**
- **Date:** 2026-06-03
- **Issue:** [#369](https://github.com/theding0x/capital-simulator/issues/369) (umbrella) — related: #216, #222, #370

## Context

The simulator models capital as *value in motion* — the circuit
**M—C(Lp+Mp)…P…C′—M′** — chapter by chapter. It is faithful **node-by-node**,
but value does not currently **compose across nodes** in the running system:

1. **No value↔money bridge.** Vol. I production measures value in
   `LabourMinutes` (labour-substance — faithful to Marx). Vol. II/III
   circulation and distribution measure it in `Pence` (money-form — also
   faithful at the M/M′ ends). But there is **no conversion** between them. For
   example `simulation-engine/internal/surplus.SurplusValueMass` (LabourMinutes)
   and `circulation.CommodityCapital.ValueSurplus` (Pence) are the *same*
   economic surplus `s` in two units, with no path between them.
2. **finance-service is decoupled.** Every Vol. III endpoint's `c`, `v`, `s` are
   **caller-supplied** inputs, never sourced from the `P` node that produces
   surplus-value. The law *L-Mp2* ("ΔM = s, and s originates in P") therefore
   holds *locally* in each service, but value never actually flows P→ΔM
   end-to-end.
3. **`Pence` is package-local.** It is defined per package with differing
   £-scales (100 / 240 / ×1000); see `docs/architecture.md` →
   *"Money scale across packages"* and issue #370.

Because each calculator is internally consistent and the seeds are mutually
consistent, runtime conservation checks (Σprice = Σvalue, etc.) all pass. This
is a **composition gap**, not a wrong-number bug — but it is the gap between
`CLAUDE.md`'s framing ("value in motion … not a separate exhibit") and the
implementation. It is the headline structural finding of the whole-circuit
`/capital-reflection`.

## Decision

1. **Canonical value↔money bridge — the MELT.** Adopt a single *Monetary
   Expression of Labour Time*: a rate of **Pence per LabourMinute**, grounded in
   Marx's "value of money" (the labour-time embodied in the money-commodity,
   gold). Every conversion between the labour-substance unit (`LabourMinutes`)
   and the money-form (`Pence`) goes through the MELT. The MELT, together with a
   **single canonical £→pence scale**, lives in one shared package (e.g.
   `pkg/money`) that owns both the conversion and the pence-scale normalisation
   documented for #370. No service re-implements either.

2. **Downstream nodes consume upstream magnitudes.** Rather than re-deriving
   `c`, `v`, `s` at each node, downstream nodes (circulation, then
   distribution/finance) **consume** the magnitudes produced at the `P` node —
   via domain events and/or shared-store reads. This extends the cross-service
   circuit-leg precedent established in #216 (the market→agent webhook) and the
   Vol. II integration in #222.

3. **One source of surplus.** `s` is produced once, in `P`, in `LabourMinutes`;
   it is converted to `Pence` via the MELT at the M/M′ boundary; finance **reads**
   it rather than accepting a free `c/v/s`. No node independently re-derives the
   surplus.

## Scope of this decision

This ADR records the architectural **direction**. The implementation —
a `pkg/money` MELT + canonical scale, the event/store wiring from `P` to the
circulation and distribution nodes, and the migration of finance endpoints to
**source** `c/v/s` rather than accept them — is large and remains **Deferred**,
tracked by the umbrella issue #369 and its children (#216, #222). It is *not*
implemented in this change; this ADR closes the "what is the direction?"
question so future phases can build against a fixed contract.

## Consequences

- **Defuses #370.** A single canonical £→pence scale plus the MELT replaces the
  per-package pence scales for any *cross-node* flow. Existing intra-package
  fixtures stay as they are (migrations are append-only); conversion happens at
  package boundaries.
- **Gives finance a real source for `c/v/s`** — the `P` node — so *L-Mp2* holds
  end-to-end rather than only locally.
- **Realises "value in motion":** the circuit *flows* P→ΔM instead of each node
  re-deriving its own magnitudes from caller input.
- **Cost:** a shared money package becomes a coupling point that every
  value-bearing service depends on. It must be kept free of domain logic beyond
  conversion, and the MELT must be a single well-known value (configuration or a
  documented constant), not duplicated per service.
