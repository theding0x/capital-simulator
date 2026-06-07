# Atlas Observatory — Redesign Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Rebuild the Atlas landing experience to match the handed-back Claude Design pass — **one continuous vertical world** you descend through: a luminous `<canvas>` orrery surface above, the hidden abode of production below, the reserve army beneath — with a gilded threshold gate, the gold-mote alienation loop, a promoted immiseration chart, and a reserve-army reservoir.

**Architecture:** Frontend-only. The mock engine in the prototype is replaced by the live `useSnapshot()` poll + the real engine/lever endpoints. The SVG orbits (`CircuitField`/`Orbit`) become an **imperative TS canvas controller** (`AtlasSurface`) fed snapshots each poll, owning its own `requestAnimationFrame`. `Atlas.tsx` becomes the scrolling world (descent camera, scroll-derived depth, gate, descent tint). **No backend changes** — the existing `ObservatorySnapshot`/`AbodeReadout`/`law_series` + `GET /v1/engine/ticks` already supply every field.

**Tech Stack:** React 18 + Vite + TypeScript, plain CSS on `web/src/index.css :root` tokens, Canvas 2D + `requestAnimationFrame`, no router beyond the existing hash `Root.tsx`, no new dependencies.

**Design source (in-repo reference — copy from it):** `design/design_handoff_atlas_observatory/` — `README.md` (the spec) and `prototype/{app.jsx,surface.js,components.jsx,atlas.css}` (the working reference). The prototype's `data.js` is a **mock — do not port it**; the real data comes from `useSnapshot.ts`.

**Branch:** `feature/atlas-observatory` (continues the v2 work).

---

## Adaptations (prototype → production — apply everywhere)

1. **Mock engine → live data.** Drop `window.AtlasEngine`/`makeEngine`/`eng.step`/`setInterval` stepping. `Atlas.tsx` consumes `useSnapshot()` (polls `GET /v1/observatory/snapshot` every 2s, holds last-good). The law advances **server-side** (the scheduler's `GeneralLawTicker`, ~5s); the page just renders successive snapshots.
2. **Helpers.** The prototype's `AtlasEngine.{lerp,clamp,hash,formatPence,formatBP,formatMinutes}` map to `web/src/atlas/animation.ts`. `formatPence/formatBP/formatMinutes` already exist; **add `clamp`, `lerp`, `hash`**.
3. **Canvas controller.** `surface.js` (an IIFE attaching `window.AtlasSurface`) becomes `web/src/atlas/surface.ts` — a normally-exported, typed `AtlasSurface` class importing `{ clamp, lerp, hash }` from `./animation`. Keep all draw logic verbatim; type it; replace the `'Meslo LG'` font literal with `'IBM Plex Mono'`; add a `setSpeed`.
4. **Speed** is a **cosmetic surface-tempo multiplier** (orbit drift + dot laps), not a law accelerator (the law is server-paced). Transport play/pause calls the real `api.startEngine()/stopEngine()`.
5. **Lever ranges track the real backend** (seed `surplus_rate_base_bp` 10000, `base_wage_pence` 2500, `accumulation_rate_bp` 5000; clamps α∈[0,10000], wage≥1, surplus∈[0,100000]) — **not** the mock's 200–4000. Use: working-day `min 2000 max 40000 step 500`, wage `min 500 max 6000 step 100`, α `min 0 max 10000 step 250`. Slider state is local (seeded from the snapshot once), POSTs on change via `api.setObservatoryLevers` (Slice-3 behavior); the law's response shows on the next poll.
6. **Fonts.** `index.css :root` has the **color/radius** tokens but **not** `--font-display/--font-sans/--font-mono`. Define those three on `.atlas {}` in the new `atlas.css` (Playfair Display / IBM Plex Sans / IBM Plex Mono — already `@import`ed by `index.css`). Keep the repo's IBM Plex Mono (the design-system's Meslo substitution does not apply here).
7. **Transport ECG** reads `api.listEngineTicks(60)` → `entities_advanced` (the existing `TickHeartbeat` already does this; keep that fetch).

---

## File Structure

**Create:**
- `web/src/atlas/surface.ts` — the `AtlasSurface` canvas controller (port of `surface.js`).
- `web/src/atlas/ImmiserationChart.tsx` — the promoted SVG chart (port of `components.jsx` `<ImmiserationChart>`); replaces `GeneralLawTrend.tsx`.
- `web/src/atlas/ReserveArmy.tsx` — the reservoir + pressure (port of `<ReserveArmy>`).

**Modify (rewrite):**
- `web/src/atlas/animation.ts` — add `clamp`, `lerp`, `hash`.
- `web/src/atlas/Atlas.tsx` — the continuous world (descent camera, depth, gate, loop, instruments wiring).
- `web/src/atlas/Abode.tsx` — `WorkingDay` hero + `StatRow` + chart-card + `ReserveArmy` + `Levers`.
- `web/src/atlas/Levers.tsx` — prototype markup, **real backend ranges**, self-POST.
- `web/src/atlas/VitalSigns.tsx` — `.vitals` markup.
- `web/src/atlas/TickHeartbeat.tsx` — export `Transport` (play/pause, speed, ECG, reduced-motion toggle, turn).
- `web/src/atlas/atlas.css` — replace wholesale with the prototype CSS + the `.atlas` font-var header.

**Delete:**
- `web/src/atlas/CircuitField.tsx`, `web/src/atlas/Orbit.tsx` (replaced by `surface.ts`), `web/src/atlas/GeneralLawTrend.tsx` (replaced by `ImmiserationChart.tsx`).

**Unchanged:** `web/src/atlas/useSnapshot.ts`, `web/src/Root.tsx`, `web/src/types.ts`, `web/src/api.ts`, `web/src/index.css`.

---

# GROUP A — Engine layer (math + canvas)

## Task A1: animation helpers

**Files:** Modify `web/src/atlas/animation.ts`

- [ ] **Step 1: Append the three helpers** (keep all existing exports):

```ts
/** Clamp v to [lo, hi]. */
export function clamp(v: number, lo: number, hi: number): number {
  return v < lo ? lo : v > hi ? hi : v;
}

/** Linear interpolation from a to b by t (t in [0,1]). */
export function lerp(a: number, b: number, t: number): number {
  return a + (b - a) * t;
}

/** Deterministic 32-bit hash of a string — stable per-id pseudo-randomness. */
export function hash(s: string): number {
  let h = 0;
  for (let i = 0; i < s.length; i++) h = (h * 31 + s.charCodeAt(i)) >>> 0;
  return h;
}
```

- [ ] **Step 2: Typecheck** — `cd web && npm run lint 2>&1 | tail` → no errors.
- [ ] **Step 3: Commit**

```bash
git add web/src/atlas/animation.ts
git commit --no-gpg-sign -m "feat(atlas): add clamp/lerp/hash helpers for the canvas orrery"
```

## Task A2: the `AtlasSurface` canvas controller

**Files:** Create `web/src/atlas/surface.ts`

- [ ] **Step 1: Port `surface.js` → `surface.ts`.** Copy the body of `design/design_handoff_atlas_observatory/prototype/surface.js` and transform:
  - Remove the IIFE wrapper and `window.AtlasSurface = ...`; remove `const { lerp, clamp, hash } = window.AtlasEngine;`.
  - Add at top:
    ```ts
    import type { ObservatorySnapshot } from "../types";
    import { clamp, lerp, hash } from "./animation";
    ```
  - Define types and export the class:
    ```ts
    type FieldCapital = ObservatorySnapshot["capitals"][number];
    interface Body {
      id: string; cap?: FieldCapital;
      theta: number; orbitSpeed: number; refTotal: number; ang: number;
      dots: number[]; growth: number; dist: number; ringR: number;
    }
    interface Mote { x: number; y: number; vy: number; life: number; r: number; }
    const GOLD: [number, number, number] = [200, 162, 64];
    // …GOLD_HI, RED, LEAD, BONE, rgba(), pacedFrac() exactly as in the prototype…
    export class AtlasSurface { /* …prototype class body… */ }
    ```
  - Type the fields: `private canvas: HTMLCanvasElement; private ctx: CanvasRenderingContext2D; private snap: ObservatorySnapshot | null = null; private bodies = new Map<string, Body>(); private motes: Mote[] = [];` and the numeric scalars (`depth`, `t`, `refMax`, `_raf`, `_last`, `_spawnAcc`, `_surplusPrev`, `_pulse`, `_dpr`, `_w`, `_h`). `setSnapshot(snap: ObservatorySnapshot)`, `setDepth(d: number)`, `setReduced(v: boolean)`, `setRunning(v: boolean)`.
  - **Add a speed control:** add `private speed = 1;` and `setSpeed(v: number) { this.speed = v > 0 ? v : 1; }`. In `_frame`, multiply motion by speed: orbit drift `b.theta += b.orbitSpeed * dt * this.speed * (this.running ? 1 : 0);` and dot laps `b.dots[i] = (b.dots[i] + dt * lap * this.speed) % 1;`.
  - Replace the halted-label font `"9px 'Meslo LG', monospace"` → `"9px 'IBM Plex Mono', monospace"`.
  - Keep **all** other draw logic byte-for-byte (ground, core/`p̄′`, guide rings, bodies, three arcs, value-dots, motes, vignette).

- [ ] **Step 2: Typecheck** — `cd web && npm run lint 2>&1 | tail`. Resolve any `noUnusedLocals`/strict-null issues (e.g. guard `getContext("2d")` non-null with `as CanvasRenderingContext2D`). Expected: clean.
- [ ] **Step 3: Commit**

```bash
git add web/src/atlas/surface.ts
git commit --no-gpg-sign -m "feat(atlas): AtlasSurface — the canvas orrery controller (port of the design prototype)"
```

---

# GROUP B — DOM layer (instruments + abode)

## Task B1: `VitalSigns` (rail)

**Files:** Modify `web/src/atlas/VitalSigns.tsx`

- [ ] **Step 1: Replace contents** — port `components.jsx` `VitalSigns` to TS:

```tsx
import type { AggregateVitals } from "../types";
import { formatPence, formatBP } from "./animation";

/** Rail vital signs: total social capital, p̄′ (gold), Σs, Σc, capitals in motion. */
export function VitalSigns({ vitals, count }: { vitals: AggregateVitals; count: number }) {
  const rows = [
    { k: "Total social capital", v: formatPence(vitals.total_social_capital_pence), gold: false },
    { k: "Average rate of profit · p̄′", v: formatBP(vitals.avg_rate_of_profit_bp), gold: true },
    { k: "Surplus-value · Σs", v: formatPence(vitals.surplus_pence), gold: false },
    { k: "Cost-price · Σc (c+v)", v: formatPence(vitals.cost_price_pence), gold: false },
    { k: "Capitals in motion", v: String(count), gold: false },
  ];
  return (
    <div className="vitals">
      <div className="vitals-title">Vital signs</div>
      {rows.map((r) => (
        <div className="vital" key={r.k}>
          <div className="vital-k">{r.k}</div>
          <div className={"vital-v" + (r.gold ? " gold" : "")}>{r.v}</div>
        </div>
      ))}
    </div>
  );
}
```

(Signature changes from the shipped `{ vitals, capitalCount }` to `{ vitals, count }` — `Atlas.tsx` in Group C passes `count`.)

- [ ] **Step 2: Typecheck** — `cd web && npm run lint 2>&1 | tail` (Atlas.tsx may still reference the old shape until Group C; that error is expected and clears in C2). Verify VitalSigns.tsx itself has no type error.
- [ ] **Step 3: Commit**

```bash
git add web/src/atlas/VitalSigns.tsx
git commit --no-gpg-sign -m "feat(atlas): restyle VitalSigns to the redesign rail"
```

## Task B2: `ImmiserationChart` (promoted) — replaces `GeneralLawTrend`

**Files:** Create `web/src/atlas/ImmiserationChart.tsx`; Delete `web/src/atlas/GeneralLawTrend.tsx`

- [ ] **Step 1: Create the chart** — port `components.jsx` `<ImmiserationChart>` verbatim into a TS component:

```tsx
import { useMemo } from "react";
import type { GeneralLawTrendPoint } from "../types";
import { formatPence, formatBP } from "./animation";

/** The immiseration trend: s/v rising (gold), wage falling (red dashed), reserve
 *  army swelling (lead area). The wage/s-v crossing IS the immiseration story. */
export function ImmiserationChart({ series }: { series: GeneralLawTrendPoint[] }) {
  const W = 680, H = 210, padL = 16, padR = 16, padT = 18, padB = 26;
  const geom = useMemo(() => {
    if (!series || series.length < 2) return null;
    // …paste the prototype's geom body verbatim (sv/wg/ra arrays, x(), norm(), ny(),
    //   line(), area(), returning svPath/wgPath/raArea/last/n/svDot/wgDot)…
  }, [series]);
  if (!geom) return <p className="chart-empty">The law has not yet run — start the engine.</p>;
  const innerH = H - padT - padB;
  return (/* …paste the prototype's <div className="chart"> … </div> verbatim… */);
}
```

Type the dot tuples as `[number, number]`. The legend uses `formatBP(geom.last.rate_of_exploitation_bp)` and `formatPence(geom.last.wage_pence)`.

- [ ] **Step 2: Delete the old trend** — `git rm web/src/atlas/GeneralLawTrend.tsx` (Abode.tsx will import the new chart in B5).
- [ ] **Step 3: Typecheck** — `cd web && npm run lint 2>&1 | tail` (Abode.tsx import error expected until B5).
- [ ] **Step 4: Commit**

```bash
git add web/src/atlas/ImmiserationChart.tsx
git rm web/src/atlas/GeneralLawTrend.tsx
git commit --no-gpg-sign -m "feat(atlas): promote the immiseration trend to a full chart"
```

## Task B3: `ReserveArmy`

**Files:** Create `web/src/atlas/ReserveArmy.tsx`

- [ ] **Step 1: Port `<ReserveArmy>`** from `components.jsx`:

```tsx
import { useMemo } from "react";
import type { AbodeReadout } from "../types";
import { clamp, formatBP, formatPence } from "./animation";

/** The industrial reserve army as a reservoir of units; a pressure bar widths to
 *  reserve_army_pressure_bp. The disposable reserve made physical. */
export function ReserveArmy({ abode }: { abode: AbodeReadout }) {
  const count = abode.reserve_army_count;
  const pressure = abode.reserve_army_pressure_bp / 10000;
  const dots = Math.min(140, Math.round(count / 2.6));
  const cells = useMemo(() => Array.from({ length: 140 }), []);
  return (/* …paste the prototype <section className="reserve">…</section> verbatim,
            using {`${clamp(pressure*100,4,100)}%`} for the .reserve-fill width… */);
}
```

- [ ] **Step 2: Typecheck** — `cd web && npm run lint 2>&1 | tail`.
- [ ] **Step 3: Commit**

```bash
git add web/src/atlas/ReserveArmy.tsx
git commit --no-gpg-sign -m "feat(atlas): the industrial reserve army reservoir"
```

## Task B4: `Levers` (restyle + real ranges)

**Files:** Modify `web/src/atlas/Levers.tsx`

- [ ] **Step 1: Rewrite** with the prototype's `Lever`/`Levers` markup but **real backend ranges** and **self-POST** (local slider state, seeded once):

```tsx
import { useState } from "react";
import type { AbodeReadout, LeverUpdate } from "../types";
import { api } from "../api";
import { formatBP, formatPence } from "./animation";

function Lever(props: {
  label: string; sub: string; value: number; display: string;
  min: number; max: number; step: number; testId: string; onChange: (v: number) => void;
}) {
  return (
    <label className="lever">
      <div className="lever-top">
        <span className="lever-label">{props.label}</span>
        <span className="lever-value">{props.display}</span>
      </div>
      <input type="range" min={props.min} max={props.max} step={props.step} value={props.value}
        data-testid={props.testId}
        onChange={(e) => props.onChange(Number(e.target.value))} />
      <span className="lever-sub">{props.sub}</span>
    </label>
  );
}

/** The levers — perturb the live AbodeState; the law ripples over the next polls.
 *  Slider state is local (seeded once); ranges track the real backend, not the
 *  design mock. */
export function Levers({ abode }: { abode: AbodeReadout }) {
  const [surplus, setSurplus] = useState(abode.surplus_rate_base_bp);
  const [wage, setWage] = useState(abode.base_wage_pence);
  const [accum, setAccum] = useState(abode.accumulation_rate_bp);
  const push = (u: LeverUpdate) => {
    void api.setObservatoryLevers(u).catch(() => { /* next poll reflects truth */ });
  };
  return (
    <section className="levers-card" data-testid="levers">
      <div className="levers-eyebrow">The levers — perturb the law</div>
      <div className="levers">
        <Lever label="Working day · rate of surplus-value" sub="lengthen the unpaid hours"
          value={surplus} display={formatBP(surplus)} min={2000} max={40000} step={500} testId="lever-workingday"
          onChange={(v) => { setSurplus(v); push({ surplus_rate_base_bp: v }); }} />
        <Lever label="Wage · value of labour-power" sub="press it toward subsistence"
          value={wage} display={formatPence(wage)} min={500} max={6000} step={100} testId="lever-wage"
          onChange={(v) => { setWage(v); push({ base_wage_pence: v }); }} />
        <Lever label="Accumulation rate · α" sub="reinvest surplus as machinery"
          value={accum} display={formatBP(accum)} min={0} max={10000} step={250} testId="lever-accumulation"
          onChange={(v) => { setAccum(v); push({ accumulation_rate_bp: v }); }} />
      </div>
    </section>
  );
}
```

- [ ] **Step 2: Typecheck** — `cd web && npm run lint 2>&1 | tail`.
- [ ] **Step 3: Commit**

```bash
git add web/src/atlas/Levers.tsx
git commit --no-gpg-sign -m "feat(atlas): restyle the levers (design markup, real backend ranges)"
```

## Task B5: `Abode` (compose the hero + tiles + chart + reserve + levers)

**Files:** Modify `web/src/atlas/Abode.tsx`

- [ ] **Step 1: Rewrite `Abode.tsx`** — port `components.jsx` `WorkingDay`, `StatRow`, and `Abode` (local `WorkingDay`/`StatRow`, imported `ImmiserationChart`/`ReserveArmy`/`Levers`):

```tsx
import type { AbodeReadout } from "../types";
import { formatPence, formatBP, formatMinutes } from "./animation";
import { ImmiserationChart } from "./ImmiserationChart";
import { ReserveArmy } from "./ReserveArmy";
import { Levers } from "./Levers";

function WorkingDay({ abode }: { abode: AbodeReadout }) { /* …prototype WorkingDay; keep data-testid="workingday" on .wd-bar… */ }
function StatRow({ abode }: { abode: AbodeReadout }) { /* …prototype StatRow… */ }

/** The hidden abode: working-day hero, demoted tiles, the immiseration chart
 *  (promoted), the reserve army, the levers. */
export function Abode({ abode }: { abode: AbodeReadout }) {
  return (
    <div className="abode-inner" data-testid="abode">
      <WorkingDay abode={abode} />
      <StatRow abode={abode} />
      <section className="chart-card">
        <div className="chart-eyebrow">The general law in motion</div>
        <p className="chart-gloss">Accumulation widens the rate of exploitation, swells the reserve army, and presses the wage to its floor — the immiseration of the producer, run in real time.</p>
        <ImmiserationChart series={abode.law_series} />
      </section>
      <ReserveArmy abode={abode} />
      <Levers abode={abode} />
    </div>
  );
}
```

Keep `data-testid="abode"` on the container and `data-testid="workingday"` on the working-day bar (E2E parity). The export becomes a **named** `Abode` (Atlas.tsx imports `{ Abode }`).

- [ ] **Step 2: Typecheck** — `cd web && npm run lint 2>&1 | tail`.
- [ ] **Step 3: Commit**

```bash
git add web/src/atlas/Abode.tsx
git commit --no-gpg-sign -m "feat(atlas): rebuild the abode — working-day hero, chart, reserve, levers"
```

## Task B6: `Transport` (console)

**Files:** Modify `web/src/atlas/TickHeartbeat.tsx`

- [ ] **Step 1: Rewrite** as a controlled `Transport` that self-fetches the ECG ticks:

```tsx
import { useEffect, useMemo, useRef, useState } from "react";
import { api } from "../api";
import type { EngineTick } from "../types";

interface TransportProps {
  tick: number; running: boolean; onToggle: () => void;
  speed: number; onSpeed: (s: number) => void;
  reduced: boolean; onReduced: (v: boolean) => void;
}

/** The console: play/pause, ×1/2/5/10 speed, an ECG of recent engine ticks, a
 *  reduced-motion toggle, and the turn counter. */
export function Transport({ tick, running, onToggle, speed, onSpeed, reduced, onReduced }: TransportProps) {
  const [ticks, setTicks] = useState<EngineTick[]>([]);
  const timer = useRef<number | null>(null);
  useEffect(() => {
    let active = true;
    const poll = async () => { try { const t = await api.listEngineTicks(60); if (active) setTicks(t); } catch { /* keep last */ } };
    void poll();
    timer.current = window.setInterval(() => void poll(), 2000);
    return () => { active = false; if (timer.current !== null) window.clearInterval(timer.current); };
  }, []);
  const pts = useMemo(() => {
    if (!ticks.length) return "0,14 240,14";
    const ordered = [...ticks].reverse(); // oldest → newest
    const max = Math.max(1, ...ordered.map((t) => t.entities_advanced));
    const n = ordered.length;
    return ordered.map((t, i) => `${((i / Math.max(1, n - 1)) * 240).toFixed(1)},${(26 - (t.entities_advanced / max) * 24).toFixed(1)}`).join(" ");
  }, [ticks]);
  return (/* …prototype <div className="transport"> … </div> verbatim, wired to these props + {pts} for the ECG polyline… */);
}
```

(`EngineTick` has `entities_advanced` — the shipped `TickHeartbeat` already used it. Any prior default export is dropped — `Atlas.tsx` imports `{ Transport }`.)

- [ ] **Step 2: Typecheck** — `cd web && npm run lint 2>&1 | tail`.
- [ ] **Step 3: Commit**

```bash
git add web/src/atlas/TickHeartbeat.tsx
git commit --no-gpg-sign -m "feat(atlas): console transport — play/pause, speed, ECG, reduced-motion"
```

---

# GROUP C — Integration

## Task C1: `atlas.css` (replace wholesale)

**Files:** Modify `web/src/atlas/atlas.css`

- [ ] **Step 1: Replace the entire file** with `design/design_handoff_atlas_observatory/prototype/atlas.css`, with two edits:
  - In the `.atlas { … }` block, **add the font vars** (the prototype assumes them; `index.css :root` does not define them):
    ```css
    .atlas {
      --font-display: 'Playfair Display', Georgia, serif;
      --font-sans: 'IBM Plex Sans', ui-sans-serif, system-ui, sans-serif;
      --font-mono: 'IBM Plex Mono', ui-monospace, monospace;
      /* …keep the prototype's --rail-w / --topbar-h / --console-h and the rest… */
    }
    ```
  - Drop the prototype's global `html, body { … } #root { … }` rule (let `index.css` own the page); keep everything else verbatim (it uses the real `:root` color/radius tokens).

- [ ] **Step 2:** (CSS has no types; build verified in C2.)
- [ ] **Step 3: Commit**

```bash
git add web/src/atlas/atlas.css
git commit --no-gpg-sign -m "feat(atlas): redesign stylesheet — continuous world, gate, abode, console"
```

## Task C2: `Atlas.tsx` (the continuous world) + cleanup

**Files:** Modify `web/src/atlas/Atlas.tsx`; Delete `web/src/atlas/CircuitField.tsx`, `web/src/atlas/Orbit.tsx`

- [ ] **Step 1: Rewrite `Atlas.tsx`** — port `app.jsx`'s structure, wired to live data. Full file:

```tsx
import { useCallback, useEffect, useRef, useState } from "react";
import "./atlas.css";
import { useSnapshot } from "./useSnapshot";
import { AtlasSurface } from "./surface";
import { VitalSigns } from "./VitalSigns";
import { Abode } from "./Abode";
import { Transport } from "./TickHeartbeat";
import { api } from "../api";
import { clamp, formatBP } from "./animation";

/** Animated scrollTop tween (inOutCubic) on the stage container. */
function useAnimatedScroll(stageRef: React.RefObject<HTMLDivElement | null>) {
  const raf = useRef(0);
  return useCallback((to: number, ms: number, done?: () => void) => {
    const el = stageRef.current; if (!el) return;
    cancelAnimationFrame(raf.current);
    const from = el.scrollTop, d = to - from, t0 = performance.now();
    const ease = (x: number) => (x < 0.5 ? 4 * x * x * x : 1 - Math.pow(-2 * x + 2, 3) / 2);
    if (ms <= 0) { el.scrollTop = to; done?.(); return; }
    const step = (now: number) => {
      const p = Math.min(1, (now - t0) / ms);
      el.scrollTop = from + d * ease(p);
      if (p < 1) raf.current = requestAnimationFrame(step);
      else done?.();
    };
    raf.current = requestAnimationFrame(step);
  }, [stageRef]);
}

/** The Observatory — one continuous vertical world: the orrery surface above,
 *  the hidden abode below; descending is literal travel through a gilded gate. */
export default function Atlas() {
  const { snapshot } = useSnapshot();
  const prefersReduced = typeof window !== "undefined" &&
    window.matchMedia && window.matchMedia("(prefers-reduced-motion: reduce)").matches;

  const [running, setRunning] = useState(false);
  const [speed, setSpeed] = useState(1);
  const [reduced, setReduced] = useState(!!prefersReduced);
  const [depth, setDepth] = useState(0);

  const surfaceRef = useRef<AtlasSurface | null>(null);
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const stageRef = useRef<HTMLDivElement>(null);
  const surfaceZoneRef = useRef<HTMLDivElement>(null);
  const animateScroll = useAnimatedScroll(stageRef);

  // instantiate the canvas controller once
  useEffect(() => {
    if (!canvasRef.current) return;
    const surf = new AtlasSurface(canvasRef.current);
    surfaceRef.current = surf;
    surf.setReduced(!!prefersReduced);
    surf.start();
    return () => surf.stop();
  }, []); // eslint-disable-line react-hooks/exhaustive-deps

  // feed the controller each new snapshot
  useEffect(() => { if (snapshot) surfaceRef.current?.setSnapshot(snapshot); }, [snapshot]);
  useEffect(() => { surfaceRef.current?.setReduced(reduced); }, [reduced]);
  useEffect(() => { surfaceRef.current?.setSpeed(speed); }, [speed]);
  // reflect the server's run state on first load
  useEffect(() => { if (snapshot) setRunning(snapshot.running); }, [snapshot?.running]); // eslint-disable-line react-hooks/exhaustive-deps

  const depthRaf = useRef(0);
  const onScroll = useCallback(() => {
    cancelAnimationFrame(depthRaf.current);
    depthRaf.current = requestAnimationFrame(() => {
      const el = stageRef.current; if (!el) return;
      const h = surfaceZoneRef.current?.offsetHeight || window.innerHeight;
      const d = clamp(el.scrollTop / (h * 0.9), 0, 1);
      setDepth(d);
      surfaceRef.current?.setDepth(d);
      surfaceRef.current?.setRunning(running && d < 0.99);
    });
  }, [running]);

  const descend = () => animateScroll(surfaceZoneRef.current?.offsetHeight || window.innerHeight, reduced ? 0 : 1150);
  const ascend = () => animateScroll(0, reduced ? 0 : 1000);

  const toggleRun = () => {
    const next = !running;
    setRunning(next);
    surfaceRef.current?.setRunning(next);
    void (next ? api.startEngine() : api.stopEngine()).catch(() => { /* poll reflects truth */ });
  };

  const descended = depth > 0.5;

  return (
    <div className="atlas">
      <header className="topbar">
        <span className="brand"><span className="brand-mark">C</span> Capital Simulator <span className="brand-sub">· Atlas</span></span>
        <nav className="nav"><a className="active" href="#/">Atlas</a><a href="#/chapters">Chapters</a></nav>
      </header>

      <aside className="rail">
        {snapshot && <VitalSigns vitals={snapshot.aggregate} count={snapshot.capitals.length} />}
        <button className={"descend-btn" + (descended ? " open" : "")} data-testid="threshold" onClick={descended ? ascend : descend}>
          <span className="db-arrow">{descended ? "↑" : "↓"}</span>
          <span className="db-label">{descended ? "Ascend to the surface" : "Descend into production"}</span>
        </button>
        <p className="rail-foot">An observatory of the circuit of capital — <span className="i">M—C…P…C′—M′</span> — at the scale of many capitals.</p>
      </aside>

      <main className="stage" ref={stageRef} onScroll={onScroll}>
        <div className="world">
          <section className="zone-surface" ref={surfaceZoneRef}>
            <canvas className="surface-canvas" ref={canvasRef}></canvas>
            {snapshot && (
              <div className="centre-label" aria-hidden="true">
                <div className="centre-pbar">{formatBP(snapshot.aggregate.avg_rate_of_profit_bp)}</div>
                <div className="centre-lbl">p̄′ · centre of gravity</div>
              </div>
            )}
            <div className="surface-caption">
              <h1 className="surface-title">The circuit of capital</h1>
              <p className="surface-sub">Each capital a ring of three coexisting arcs — <span className="g">M</span> money, <span className="r">P</span> production, <span className="l">C′</span> commodity — value travelling its circumference. They spiral outward as they accumulate.</p>
            </div>
            <div className="gate" style={{ ["--cross"]: clamp((depth - 0.4) / 0.5, 0, 1) } as React.CSSProperties}>
              <div className="gate-half left"></div>
              <div className="gate-half right"></div>
              <button className="gate-notice" onClick={descend}>
                <span className="gate-no">No admittance except on business</span>
                <span className="gate-cite">Capital, Vol. I · Ch. 6 — leave the noisy sphere; descend ↓</span>
              </button>
            </div>
          </section>

          <section className={"zone-abode" + (descended ? " lit" : "")}>
            <div className="abode-seam"></div>
            <div className="abode-head">
              <div className="abode-eyebrow">The hidden abode of production</div>
              <h2 className="abode-title">Where surplus is pumped from living labour</h2>
            </div>
            {snapshot && <Abode abode={snapshot.abode} />}
          </section>
        </div>
      </main>

      <footer className="console">
        <Transport tick={snapshot?.tick ?? 0} running={running} onToggle={toggleRun}
          speed={speed} onSpeed={setSpeed} reduced={reduced} onReduced={setReduced} />
      </footer>

      <div className="descent-tint" style={{ opacity: depth * 0.5 }}></div>
    </div>
  );
}
```

(The `{ ["--cross"]: … } as React.CSSProperties` sets a CSS custom property in a typed style object.)

- [ ] **Step 2: Delete the SVG orbit components** — `git rm web/src/atlas/CircuitField.tsx web/src/atlas/Orbit.tsx`. Confirm nothing else imports them: `grep -rn "CircuitField\|from \"./Orbit\"\|from \"./CircuitField\"" web/src` → no matches remain.
- [ ] **Step 3: Typecheck + build** — `cd web && npm run lint && npm run build 2>&1 | tail -25` → success. (`animation.ts`'s now-unused `spinSeconds`/`orbitRadius`/`arcFractions`/`pacedAngle`/`targetScale`/`sparklinePath`/`ease` remain exported, so they won't trip `noUnusedLocals`; leave them.)
- [ ] **Step 4: Commit**

```bash
git add web/src/atlas/Atlas.tsx
git rm web/src/atlas/CircuitField.tsx web/src/atlas/Orbit.tsx
git commit --no-gpg-sign -m "feat(atlas): the continuous observatory world — descent, depth, gate, alienation loop"
```

---

# GROUP D — End-to-end acceptance

## Task D1: Boot + drive the redesigned world (Playwright MCP)

**Files:** none.

- [ ] **Step 1: Boot** — `docker compose down -v && docker compose up --build -d`; wait ~3 min. Start the engine and confirm the snapshot is intact (the contract is unchanged from Slice 3):
  ```bash
  curl -s -X POST http://localhost:8080/v1/engine/start >/dev/null
  curl -s http://localhost:8080/v1/observatory/snapshot | python3 -c "import sys,json; d=json.load(sys.stdin); print('capitals', len(d['capitals']), 'abode s/v', d['abode']['rate_of_exploitation_bp'])"
  ```
- [ ] **Step 2: Playwright MCP** against `http://localhost:5173/`:
  1. `browser_navigate` → confirm the **canvas** `.surface-canvas`, the `p̄′` centre label, the caption, and the threshold gate render.
  2. `browser_evaluate` → assert `.surface-canvas` exists with non-zero `width`/`height`; sample its pixels twice ~1s apart (`canvas.getContext('2d')` is the app's, so instead capture `toDataURL()` of a small region, or just confirm the controller's rAF by reading two `browser_take_screenshot`s) to confirm motion when not reduced.
  3. **Descend:** `browser_click` `[data-testid="threshold"]` → wait ~1.3s → assert `[data-testid="abode"]` is scrolled into view and the working-day hero (`[data-testid="workingday"]`), the chart (`.chart svg`), the reservoir (`.reservoir`), and `[data-testid="levers"]` are present.
  4. **Lever:** set `[data-testid="lever-workingday"]` high via the native-setter `input` event (as in the Slice-3 E2E); wait ~6s; assert the working-day surplus segment widened / s/v rose in the readout.
  5. `browser_take_screenshot` (full page) → `atlas-redesign-abode.png`; ascend (`[data-testid="threshold"]`) and screenshot the surface → `atlas-redesign-surface.png`.
  6. Toggle the console **motion** button → assert no error and the canvas holds (reduced); toggle back.
- [ ] **Step 3: Tear down + final check** — `docker compose down`; `cd web && npm run lint && npm run build`; `make vet test build` (backend unchanged — should stay green). Expected: all PASS.

---

# Done criteria

- Atlas renders as **one continuous vertical world**: a live `<canvas>` orrery (capitals as bodies orbiting `p̄′`, three coloured arcs, paced value-dots, spiraling outward, gold motes on Σs growth) above, the hidden abode below, joined by a gilded threshold gate that opens as you cross.
- **Descent** is an ~1.15s `inOutCubic` camera travel; scroll-derived `depth` drives the gate `--cross`, the descent tint, and the surface recede; the surface rAF pauses below the fold.
- The **abode** leads with the working-day hero (s/v), demotes the stat tiles, **promotes the immiseration chart**, shows the reserve-army reservoir, and the levers POST to perturb the live law.
- Built only on the existing data contract; `prefers-reduced-motion` + a manual toggle honored.
- `npm run lint && npm run build` pass; `make vet test build` stays green; Playwright confirms the canvas animates, the descent works, and a lever bends the law. Old `CircuitField.tsx`/`Orbit.tsx`/`GeneralLawTrend.tsx` removed.
