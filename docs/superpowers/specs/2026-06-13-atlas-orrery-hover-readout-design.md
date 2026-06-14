# Atlas orrery — hover readout on the spinning circuits

Date: 2026-06-13
Scope: `web/src/atlas/surface.ts` (+ a tiny CSS hint, + one pure-helper unit test)

## Goal

Make the spinning capital-rings in the Atlas **surface** orrery informational on
mouseover. Hovering a ring highlights it (dimming the rest of the field) and
draws a canvas readout card naming the numbers the ring already encodes.

This is a frontend-only enhancement: no React component, no `types.ts`, no
backend, no new endpoints. All rendering stays inside the existing
`AtlasSurface` canvas controller.

## Decisions (locked with the user)

- **Content**: full circuit readout — total advanced capital, rate of profit,
  the M/P/C split, surplus, turnover.
- **Presentation**: canvas-drawn card (not HTML/React tooltip), in the orrery
  aesthetic.
- **Hover feedback**: highlight the hovered ring + dim the others.
- **Identity**: no header line — the card is pure numbers. The ring under the
  cursor is its own identity. (`FieldCapital` has only a hex `id`, no name.)

## Current state

`AtlasSurface._frame(dt)` computes, per body, a screen centre (`bx, by`), a
ring radius (`ringR`), stroke width (`sw`), and the three arc fractions
(`fm, fp, fc`) from `money_pence / production_pence / commodity_pence`. There
is **no pointer handling on the canvas today** — `surface.ts` only listens to a
`ResizeObserver`.

`FieldCapital` (web/src/types.ts) carries:
`id, total_pence, money_pence, production_pence, commodity_pence,
cost_price_pence, surplus_pence, status, turnover_number`.

Reusable helpers in `web/src/atlas/animation.ts`:
`formatPence(pence) → "£N"`, `formatBP(bp) → "N.N%"`, `clamp`.

## Design

### 1. Pointer plumbing + hit-testing

- Persist each body's rendered screen geometry onto `Body`: add `sx`, `sy`,
  `sr` (set every frame from `bx`, `by`, `ringR`).
- In the constructor, register `pointermove` and `pointerleave` listeners on the
  canvas. `pointermove` converts `clientX/Y` to canvas CSS px via
  `getBoundingClientRect()` and stores `this._mx`, `this._my`.
- A pure helper resolves the hovered body from the last frame's stored geometry:

  ```ts
  // exported for unit testing
  export function pickBody(
    bodies: { id: string; sx: number; sy: number; sr: number }[],
    mx: number,
    my: number,
    tol = 6
  ): string | null
  ```

  Returns the id of the nearest body whose centre distance ≤ `sr + tol`, or
  `null`. Circular hit-radius (ignores the ecliptic squash — close enough; the
  band is generous). One-frame-stale geometry is imperceptible at 60fps and
  avoids restructuring the draw loop.
- `pointermove` calls `pickBody` against the current bodies and sets
  `this.hoveredId`. `pointerleave` clears it. Hover works while paused or under
  reduced-motion (rings are still drawn).
- `stop()` removes the listeners (alongside the existing `ResizeObserver`
  teardown).

### 2. Highlight + dim

In the body draw loop, derive a per-body factor from `this.hoveredId`:

- nothing hovered → factor `1` (unchanged behaviour).
- this body hovered → factor `1` + a thin gold outer glow stroke just outside
  `ringR`.
- another body hovered → factor `~0.3`, multiplied into the existing `dim`
  term so arcs, track, and value-dots all fade together.

### 3. The readout card

After the body loop (drawn last so it sits above the field), if `hoveredId`
resolves to a body, draw a rounded-rect card anchored to that ring:

- Width ~195px; positioned to the ring's right (`sx + sr + 14`), flipped to the
  left when it would overflow the right edge; vertically clamped within the
  canvas.
- Background: semi-opaque dark fill (`rgba(12,12,16,0.9)`) + 1px hairline
  border; subtle so it reads over the bloom.
- Rows (mono, ~11px), label left / value right:
  - **Total advanced** `C` — `formatPence(total_pence)`
  - **Rate of profit** `p′` — `formatBP(round(surplus_pence / total_pence × 10000))`
  - **M · P · C′** — three values `formatPence(money_pence)` /
    `production_pence` / `commodity_pence`, each tinted to its arc colour
    (`GOLD` / `RED` / `LEAD`) so the card maps onto the ring.
  - **Surplus** `Σs` — `formatPence(surplus_pence)`
  - **Turnover** `n` — `turnover_number.toFixed(1)`
- If `status === "halted"`, a small muted `halted` tag in the card's top-right
  corner (the only non-numeric element).

### 4. CSS hint

`web/src/atlas/atlas.css`: add `cursor: pointer` to `.surface-canvas` (or set
`canvas.style.cursor` from the controller when `hoveredId` is non-null) so the
rings read as interactive.

## Out of scope

- Touch/tap readout — the orrery is the desktop showpiece; pointer hover only.
- Currency conversion — the card shows £ via `formatPence`, matching the
  existing helpers and the rest of the app's pence basis.
- Any change to React, `types.ts`, the gateway, or services.

## Testing

- **Unit**: `surface.test.ts` (new, or extend if present) covering `pickBody`:
  hit inside a ring returns its id; a click in empty space returns `null`;
  overlapping rings return the nearer centre; tolerance band is respected.
- **Live (Playwright)**: boot the stack, load the Atlas surface, hover a ring,
  confirm (a) the readout card renders with the expected rows and (b) the other
  rings visibly dim. Per project convention, a changed panel is not done until
  driven in the browser.

## Done when

`cd web && npm run lint && npm run build` passes · `pickBody` unit test passes ·
the hover readout + dim verified live in Playwright on the surface stratum.
