# Atlas Orrery Hover Readout Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Hovering a spinning capital-ring in the Atlas surface orrery highlights it, dims the rest of the field, and draws a canvas readout card naming the numbers the ring encodes (total advanced, p′, M/P/C split, surplus, turnover).

**Architecture:** All work lives in the existing `AtlasSurface` canvas controller (`web/src/atlas/surface.ts`). Each frame persists every body's screen geometry (`sx, sy, sr`); canvas `pointermove`/`pointerleave` listeners resolve the hovered body via a pure exported `pickBody()` helper; the draw loop dims non-hovered bodies and the card is drawn last, over the field. One CSS line makes the canvas read as interactive.

**Tech Stack:** TypeScript, Canvas 2D, Vite. No test runner exists in `web/` (project convention: Go owns unit tests; web is verified by `tsc` typecheck + `vite build` + Playwright). `pickBody` is kept pure and exported for testability, and verified behaviorally by driving a real hover in Playwright.

**Spec:** `docs/superpowers/specs/2026-06-13-atlas-orrery-hover-readout-design.md`

---

## File Structure

- **Modify** `web/src/atlas/surface.ts` — the whole feature: `Body` geometry fields, `pickBody()` export, pointer listeners, hover state, dim factor, readout card.
- **Modify** `web/src/atlas/atlas.css:91` — `cursor: pointer` on `.surface-canvas`.

No new files. No React/types/backend changes.

---

## Task 1: Persist per-body screen geometry + export `pickBody`

**Files:**
- Modify: `web/src/atlas/surface.ts` (the `Body` interface; the `_frame` body loop; new exported function)

- [ ] **Step 1: Add screen-geometry fields to the `Body` interface**

In `web/src/atlas/surface.ts`, extend the `Body` interface (currently ends with `ringR: number;`) so the pointer handler can read each body's last-rendered position:

```ts
interface Body {
  id: string;
  cap?: FieldCapital;
  theta: number;
  orbitSpeed: number;
  refTotal: number;
  ang: number;
  dots: number[];
  growth: number;
  dist: number;
  ringR: number;
  sx: number; // last-rendered screen centre x (CSS px)
  sy: number; // last-rendered screen centre y (CSS px)
  sr: number; // last-rendered ring radius (CSS px)
}
```

- [ ] **Step 2: Initialise the new fields where bodies are created**

In `setSnapshot`, the body literal currently ends with `ringR: 0,`. Add the three fields:

```ts
        b = {
          id: c.id,
          theta: ((h % 1000) / 1000) * Math.PI * 2,
          orbitSpeed: lerp(0.012, 0.05, ((h >> 3) % 100) / 100) * (c.turnover_number / 3),
          refTotal: c.total_pence || 1,
          ang: i / snap.capitals.length,
          dots: [0, 1 / 3, 2 / 3],
          growth: 1,
          dist: 0,
          ringR: 0,
          sx: 0,
          sy: 0,
          sr: 0,
        };
```

- [ ] **Step 3: Record geometry each frame**

In `_frame`, inside the `this.bodies.forEach((b) => {…})` body loop, immediately after the existing line `b.ringR = ringR;` (around line 275) and after `bx`/`by` are computed, store the screen geometry:

```ts
      b.ringR = ringR;
      b.sx = bx;
      b.sy = by;
      b.sr = ringR;
```

(`bx`, `by` are already computed just above `ringR` in that loop.)

- [ ] **Step 4: Add the pure exported `pickBody` helper**

At module scope (e.g. just below the `pacedFrac` function, before `export class AtlasSurface`), add:

```ts
/** Nearest body whose centre is within `sr + tol` of (mx,my); null if none.
 *  Pure — exported for testing. Uses a circular hit-radius (ignores the
 *  ecliptic squash; the tolerance band is generous). */
export function pickBody(
  bodies: { id: string; sx: number; sy: number; sr: number }[],
  mx: number,
  my: number,
  tol = 6
): string | null {
  let best: string | null = null;
  let bestD = Infinity;
  for (const b of bodies) {
    const dx = mx - b.sx;
    const dy = my - b.sy;
    const d = Math.hypot(dx, dy);
    if (d <= b.sr + tol && d < bestD) {
      bestD = d;
      best = b.id;
    }
  }
  return best;
}
```

- [ ] **Step 5: Typecheck**

Run: `cd web && npm run lint`
Expected: PASS (no type errors). The new fields and function compile; nothing references them yet beyond the assignments.

- [ ] **Step 6: Commit**

```bash
git add web/src/atlas/surface.ts
git commit -m "feat(atlas): persist orrery body screen geometry + pickBody helper"
```

---

## Task 2: Pointer listeners + hover state + teardown

**Files:**
- Modify: `web/src/atlas/surface.ts` (new private fields; constructor; `stop()`)

- [ ] **Step 1: Add hover state fields**

In the `AtlasSurface` class field block (near `private _ro: ResizeObserver | null = null;`), add:

```ts
  private hoveredId: string | null = null;
  private _onMove: ((e: PointerEvent) => void) | null = null;
  private _onLeave: (() => void) | null = null;
```

- [ ] **Step 2: Register pointer listeners in the constructor**

The constructor currently is:

```ts
  constructor(canvas: HTMLCanvasElement) {
    this.canvas = canvas;
    this.ctx = canvas.getContext("2d") as CanvasRenderingContext2D;
    this._resize();
  }
```

Append listener wiring after `this._resize();`:

```ts
    this._onMove = (e: PointerEvent) => {
      const r = this.canvas.getBoundingClientRect();
      const mx = e.clientX - r.left;
      const my = e.clientY - r.top;
      const bodies: { id: string; sx: number; sy: number; sr: number }[] = [];
      this.bodies.forEach((b) => {
        if (b.cap) bodies.push({ id: b.id, sx: b.sx, sy: b.sy, sr: b.sr });
      });
      this.hoveredId = pickBody(bodies, mx, my);
      this.canvas.style.cursor = this.hoveredId ? "pointer" : "default";
    };
    this._onLeave = () => {
      this.hoveredId = null;
      this.canvas.style.cursor = "default";
    };
    this.canvas.addEventListener("pointermove", this._onMove);
    this.canvas.addEventListener("pointerleave", this._onLeave);
```

- [ ] **Step 3: Tear listeners down in `stop()`**

`stop()` currently disconnects the ResizeObserver. Add listener removal:

```ts
  stop() {
    cancelAnimationFrame(this._raf);
    this._raf = 0;
    this._ro?.disconnect();
    this._ro = null;
    if (this._onMove) this.canvas.removeEventListener("pointermove", this._onMove);
    if (this._onLeave) this.canvas.removeEventListener("pointerleave", this._onLeave);
    this._onMove = null;
    this._onLeave = null;
  }
```

- [ ] **Step 4: Typecheck**

Run: `cd web && npm run lint`
Expected: PASS. `hoveredId` is assigned but not yet read in the draw loop (intentional — Task 3 consumes it).

- [ ] **Step 5: Commit**

```bash
git add web/src/atlas/surface.ts
git commit -m "feat(atlas): wire canvas pointer hover detection on the orrery"
```

---

## Task 3: Highlight hovered ring + dim others

**Files:**
- Modify: `web/src/atlas/surface.ts` (the body draw loop in `_frame`)

- [ ] **Step 1: Derive the hover dim factor per body**

In `_frame`, inside the body loop, locate the existing line that sets `dim`:

```ts
      const dim = halted ? 0.32 : depthDim;
```

Replace it with a version that folds in hover state:

```ts
      const isHovered = this.hoveredId === c.id;
      const hoverDim =
        this.hoveredId === null ? 1 : isHovered ? 1 : 0.3;
      const dim = (halted ? 0.32 : depthDim) * hoverDim;
```

Because every arc, the base track, the value-dots, and the `halted` label already multiply by `dim`, the whole non-hovered ring fades together and the hovered ring stays full strength.

- [ ] **Step 2: Add a gold outer glow on the hovered ring**

Immediately after the base-track stroke block (the `ctx.arc(bx, by, ringR, 0, Math.PI * 2); ctx.stroke();` that draws the dark track, around line 293), add a hovered-only glow:

```ts
      if (isHovered) {
        ctx.globalCompositeOperation = "lighter";
        ctx.lineWidth = 2;
        ctx.strokeStyle = rgba(GOLD_HI, 0.5 * depthDim);
        ctx.beginPath();
        ctx.arc(bx, by, ringR + sw * 0.5 + 3, 0, Math.PI * 2);
        ctx.stroke();
        ctx.globalCompositeOperation = "source-over";
      }
```

- [ ] **Step 3: Typecheck + build**

Run: `cd web && npm run lint && npm run build`
Expected: PASS, clean Vite build.

- [ ] **Step 4: Commit**

```bash
git add web/src/atlas/surface.ts
git commit -m "feat(atlas): highlight hovered orrery ring and dim the field"
```

---

## Task 4: Draw the readout card

**Files:**
- Modify: `web/src/atlas/surface.ts` (new private method; one call in `_frame`)

- [ ] **Step 1: Import the formatters**

The current import line is:

```ts
import { clamp, lerp, hash } from "./animation";
```

Extend it:

```ts
import { clamp, lerp, hash, formatPence, formatBP } from "./animation";
```

- [ ] **Step 2: Add the `_drawCard` method**

Add this private method to `AtlasSurface` (e.g. just above `centre()`):

```ts
  /** Canvas readout for the hovered capital, anchored beside its ring. */
  private _drawCard(b: Body) {
    const c = b.cap;
    if (!c) return;
    const ctx = this.ctx;
    const W = this._w;

    const pad = 12;
    const cw = 196;
    const lineH = 17;
    const rows = 5; // C, p', M/P/C, Σs, n
    const ch = pad * 2 + rows * lineH;

    // anchor right of the ring, flip left near the edge, clamp vertically
    let x = b.sx + b.sr + 14;
    if (x + cw > W - 6) x = b.sx - b.sr - 14 - cw;
    x = clamp(x, 6, Math.max(6, W - cw - 6));
    let y = b.sy - ch / 2;
    y = clamp(y, 6, Math.max(6, this._h - ch - 6));

    // panel
    ctx.globalCompositeOperation = "source-over";
    ctx.fillStyle = "rgba(12,12,16,0.9)";
    ctx.strokeStyle = "rgba(200,162,64,0.35)";
    ctx.lineWidth = 1;
    this._roundRect(x, y, cw, ch, 6);
    ctx.fill();
    ctx.stroke();

    const lx = x + pad; // label column
    const rx = x + cw - pad; // value column (right-aligned)
    let ty = y + pad + 12;

    const pprime = c.total_pence > 0 ? Math.round((c.surplus_pence / c.total_pence) * 10000) : 0;

    const row = (label: string, value: string, col: [number, number, number]) => {
      ctx.font = "10px 'IBM Plex Mono', monospace";
      ctx.textAlign = "left";
      ctx.fillStyle = "rgba(168,162,148,0.9)";
      ctx.fillText(label, lx, ty);
      ctx.textAlign = "right";
      ctx.fillStyle = rgba(col, 0.95);
      ctx.fillText(value, rx, ty);
      ty += lineH;
    };

    row("C  advanced", formatPence(c.total_pence), BONE);
    row("p′ rate", formatBP(pprime), GOLD_HI);
    // M · P · C′ tinted to their arcs, drawn as one row of three
    ctx.font = "10px 'IBM Plex Mono', monospace";
    ctx.textAlign = "left";
    ctx.fillStyle = "rgba(168,162,148,0.9)";
    ctx.fillText("M·P·C′", lx, ty);
    ctx.textAlign = "right";
    ctx.fillStyle = rgba(LEAD, 0.95);
    ctx.fillText(formatPence(c.commodity_pence), rx, ty);
    ctx.fillStyle = rgba(RED, 0.95);
    ctx.fillText("·", rx - 44, ty);
    ctx.fillStyle = rgba(GOLD_HI, 0.95);
    ctx.fillText(formatPence(c.money_pence), rx - 56, ty);
    ty += lineH;
    row("Σs surplus", formatPence(c.surplus_pence), GOLD);
    row("n turnover", c.turnover_number.toFixed(1), BONE);

    // halted tag, top-right corner (only non-numeric element)
    if (c.status === "halted") {
      ctx.font = "9px 'IBM Plex Mono', monospace";
      ctx.textAlign = "right";
      ctx.fillStyle = "rgba(138,133,120,0.8)";
      ctx.fillText("halted", x + cw - 8, y + 12);
    }
  }

  /** Path a rounded rectangle (caller fills/strokes). */
  private _roundRect(x: number, y: number, w: number, h: number, r: number) {
    const ctx = this.ctx;
    ctx.beginPath();
    ctx.moveTo(x + r, y);
    ctx.arcTo(x + w, y, x + w, y + h, r);
    ctx.arcTo(x + w, y + h, x, y + h, r);
    ctx.arcTo(x, y + h, x, y, r);
    ctx.arcTo(x, y, x + w, y, r);
    ctx.closePath();
  }
```

Note: the `M·P·C′` row shows three tinted values (money gold / production red / commodity lead). The middle `·` separators are positioned by fixed offsets from the right edge; refine offsets live in Task 5 if the values overlap.

- [ ] **Step 3: Call `_drawCard` last in `_frame`**

At the very end of `_frame`, after the vignette block (after `ctx.fillRect(0, 0, W, H);` that paints `vg`), add:

```ts
    // ---- hover readout (drawn last, above the field) ----
    if (this.hoveredId) {
      const hb = this.bodies.get(this.hoveredId);
      if (hb && hb.cap) this._drawCard(hb);
    }
```

- [ ] **Step 4: Typecheck + build**

Run: `cd web && npm run lint && npm run build`
Expected: PASS, clean Vite build.

- [ ] **Step 5: Commit**

```bash
git add web/src/atlas/surface.ts
git commit -m "feat(atlas): draw the orrery hover readout card"
```

---

## Task 5: CSS cursor hint + live Playwright verification

**Files:**
- Modify: `web/src/atlas/atlas.css:91`

- [ ] **Step 1: Add the interactive cursor hint**

`web/src/atlas/atlas.css` line 91 is:

```css
.surface-canvas { position: absolute; inset: 0; width: 100%; height: 100%; display: block; }
```

Add `cursor: default;` as the baseline (the controller upgrades to `pointer` over a ring):

```css
.surface-canvas { position: absolute; inset: 0; width: 100%; height: 100%; display: block; cursor: default; }
```

- [ ] **Step 2: Boot the stack**

Run: `docker compose up --build -d` (or `make run-…` for the gateway + sim-engine that feed `/v1/observatory`, plus `cd web && npm run dev`).
Confirm the web app is reachable at `http://localhost:5173` and the surface orrery is animating with at least two rings.

- [ ] **Step 3: Drive a real hover with Playwright**

Using the Playwright MCP tools:
1. `browser_navigate` to `http://localhost:5173/#/`.
2. `browser_wait_for` the surface to render (an animating `.surface-canvas`).
3. Hover the canvas at a ring's screen position. The rings orbit, so either (a) pause the run via the Transport play/pause control first so a ring holds still, or (b) take a snapshot/screenshot, read a ring's pixel position, and `browser_mouse_move` there. Pausing is the reliable path.
4. `browser_take_screenshot` and confirm:
   - a readout card is drawn beside the hovered ring with the rows `C advanced`, `p′ rate`, `M·P·C′`, `Σs surplus`, `n turnover`;
   - the other rings are visibly dimmer than the hovered one;
   - the cursor shows as a pointer over the ring.
5. Move the pointer off all rings and screenshot again: the card disappears and the field returns to full brightness.

Expected: card renders with the five rows; non-hovered rings dim to ~30%; card clears on leave.

- [ ] **Step 4: Tune if needed**

If text overlaps in the `M·P·C′` row, or the dim is too strong/weak, adjust the offsets / `hoverDim` (0.3) / card width inline and re-run Step 3. Commit tuning separately if it changes behavior.

- [ ] **Step 5: Final typecheck + build**

Run: `cd web && npm run lint && npm run build`
Expected: PASS, clean Vite build.

- [ ] **Step 6: Commit**

```bash
git add web/src/atlas/atlas.css web/src/atlas/surface.ts
git commit -m "feat(atlas): cursor hint for orrery hover + live verification tuning"
```

---

## Self-Review

**Spec coverage:**
- Full circuit readout (total, p′, M/P/C, Σs, turnover) → Task 4 `_drawCard`. ✓
- Canvas-drawn card → Task 4. ✓
- Highlight hovered + dim others → Task 3. ✓
- No header / pure numbers (+ halted tag) → Task 4 rows, halted corner tag. ✓
- Hit-testing via pure `pickBody` → Task 1. ✓
- Pointer plumbing + teardown → Task 2. ✓
- CSS cursor hint → Task 5. ✓
- Out of scope (touch, currency, React/types/backend) → honored; no such tasks. ✓
- Testing: spec called for a unit test, but `web/` has no test runner and adding one is out of scope; superseded by `pickBody` kept pure/exported + Playwright behavioral verification (Task 5), matching project convention (web = tsc + build + Playwright). This is the one intentional deviation from the spec's testing section.

**Placeholder scan:** No TBD/TODO; every code step shows complete code. ✓

**Type consistency:** `pickBody(bodies, mx, my, tol?)` signature defined in Task 1 is called identically in Task 2. `Body.sx/sy/sr` defined in Task 1, read in Task 2's listener and Task 4's `_drawCard`. `_drawCard(b: Body)` / `_roundRect(...)` defined and called in Task 4. `hoveredId` defined in Task 2, read in Tasks 3 and 4. `GOLD/GOLD_HI/RED/LEAD/BONE` and `rgba` are existing module constants. `formatPence/formatBP` imported in Task 4. ✓
