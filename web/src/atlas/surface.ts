/* Atlas — the Surface. A lush canvas orrery of capital-orbits drifting around
 * the average rate of profit p̄′ as a luminous centre of gravity.
 *
 * Each capital is a celestial body: a ring carrying the three coexisting arcs of
 * its circuit (M money · P production · C′ commodity), with value-dots travelling
 * the circumference — one lap per turnover, lingering in production. As a capital
 * accumulates it both GROWS and SPIRALS OUTWARD from the centre.
 *
 * The alienation loop: when surplus (Σs) is pumped up from the abode below, gold
 * motes rise from the lower edge and are drawn into the orbits, thickening them.
 *
 * One rAF loop, additive bloom via layered strokes (cheap). Honours reduced-motion.
 * Typed port of the design prototype's surface.js.
 */
import type { ObservatorySnapshot } from "../types";
import { clamp, lerp, hash, formatPence, formatBP } from "./animation";

type FieldCapital = ObservatorySnapshot["capitals"][number];

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

interface Mote {
  x: number;
  y: number;
  vy: number;
  life: number;
  r: number;
}

const GOLD: [number, number, number] = [200, 162, 64]; // surplus / value / money arc
const GOLD_HI: [number, number, number] = [232, 198, 96];
const RED: [number, number, number] = [192, 57, 43]; // capital / production arc
const LEAD: [number, number, number] = [74, 90, 138]; // necessary / constant / commodity arc
const BONE: [number, number, number] = [244, 236, 216]; // value-dots
const rgba = (c: [number, number, number], a: number): string =>
  `rgba(${c[0]},${c[1]},${c[2]},${a})`;

// paced angle: lap-progress p∈[0,1) → fraction around the circle over arcs
// [fm,fp,fc], lingering in P by `lingerP`× (the dot visibly slows in production).
function pacedFrac(
  p: number,
  fm: number,
  fp: number,
  fc: number,
  lingerP: number
): number {
  const wm = fm,
    wp = fp * lingerP,
    wc = fc;
  const wsum = wm + wp + wc || 1;
  const tm = wm / wsum,
    tp = wp / wsum;
  const aM = fm,
    aP = fm + fp;
  const q = ((p % 1) + 1) % 1;
  if (q < tm) return (q / tm) * aM;
  if (q < tm + tp) return aM + ((q - tm) / tp) * (aP - aM);
  return aP + ((q - tm - tp) / (1 - tm - tp)) * (1 - aP);
}

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

export class AtlasSurface {
  private canvas: HTMLCanvasElement;
  private ctx: CanvasRenderingContext2D;
  private snap: ObservatorySnapshot | null = null;
  private bodies = new Map<string, Body>(); // id → render state
  private motes: Mote[] = []; // rising gold motes (the loop)
  private depth = 0; // 0 surface … 1 abode (recede/parallax)
  private reduced = false;
  private running = true;
  private speed = 1;
  private refMax = 0; // field reference total at load
  private t = 0; // seconds
  private _raf = 0;
  private _last = 0;
  private _spawnAcc = 0;
  private _surplusPrev = 0;
  private _pulse = 0; // recent Σs growth → mote intensity & ring thickening
  private _dpr = 1;
  private _w = 0;
  private _h = 0;
  private _ro: ResizeObserver | null = null;
  private hoveredId: string | null = null;
  private _onMove: ((e: PointerEvent) => void) | null = null;
  private _onLeave: (() => void) | null = null;

  constructor(canvas: HTMLCanvasElement) {
    this.canvas = canvas;
    this.ctx = canvas.getContext("2d") as CanvasRenderingContext2D;
    this._resize(); // measure once up front; the observer keeps it current
    this._onMove = (e: PointerEvent) => {
      const mx = e.offsetX;
      const my = e.offsetY;
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
  }

  setReduced(v: boolean) {
    this.reduced = v;
  }
  setDepth(d: number) {
    this.depth = clamp(d, 0, 1);
  }
  setSpeed(v: number) {
    this.speed = v > 0 ? v : 1;
  }

  setSnapshot(snap: ObservatorySnapshot) {
    this.snap = snap;
    const max = snap.capitals.reduce((m, c) => Math.max(m, c.total_pence), 0);
    if (this.refMax === 0 && max > 0) this.refMax = max;
    // surplus growth → loop pulse
    const s = snap.aggregate.surplus_pence;
    if (this._surplusPrev > 0) {
      const growth = (s - this._surplusPrev) / Math.max(1, this._surplusPrev);
      this._pulse = clamp(this._pulse * 0.6 + growth * 6, 0, 1.4);
    }
    this._surplusPrev = s;

    snap.capitals.forEach((c, i) => {
      let b = this.bodies.get(c.id);
      if (!b) {
        const h = hash(c.id);
        b = {
          id: c.id,
          theta: ((h % 1000) / 1000) * Math.PI * 2, // orbital angle around p̄′
          orbitSpeed: lerp(0.012, 0.05, ((h >> 3) % 100) / 100) * (c.turnover_number / 3),
          refTotal: c.total_pence || 1,
          ang: i / snap.capitals.length, // seat angle (kept stable)
          dots: [0, 1 / 3, 2 / 3],
          growth: 1,
          dist: 0,
          ringR: 0,
          sx: 0,
          sy: 0,
          sr: 0,
        };
        this.bodies.set(c.id, b);
      }
      b.cap = c;
    });
    // paint immediately so the orrery is always composed, even where rAF is
    // throttled (hidden tabs / capture contexts); the rAF loop adds motion.
    try {
      this._frame(1 / 60);
    } catch {
      /* canvas not ready yet */
    }
  }

  start() {
    if (this._raf) return;
    // Measure on element resize only — keeps getBoundingClientRect (a forced
    // synchronous layout reflow) out of the per-frame hot path.
    if (typeof ResizeObserver !== "undefined" && !this._ro) {
      this._ro = new ResizeObserver(() => this._resize());
      this._ro.observe(this.canvas);
    }
    this._last = performance.now();
    const loop = (now: number) => {
      const dt = Math.min(0.05, (now - this._last) / 1000);
      this._last = now;
      if (this.running) this._frame(dt);
      this._raf = requestAnimationFrame(loop);
    };
    this._raf = requestAnimationFrame(loop);
  }
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
  setRunning(v: boolean) {
    this.running = v;
  }

  private _resize() {
    const dpr = Math.min(2, window.devicePixelRatio || 1);
    const r = this.canvas.getBoundingClientRect();
    const w = Math.max(1, Math.round(r.width * dpr));
    const h = Math.max(1, Math.round(r.height * dpr));
    if (this.canvas.width !== w || this.canvas.height !== h) {
      this.canvas.width = w;
      this.canvas.height = h;
    }
    this._dpr = dpr;
    this._w = r.width;
    this._h = r.height;
  }

  private _frame(dt: number) {
    const ctx = this.ctx,
      W = this._w,
      H = this._h,
      dpr = this._dpr;
    ctx.setTransform(dpr, 0, 0, dpr, 0, 0);
    this.t += dt;
    this._pulse = Math.max(0, this._pulse - dt * 0.35);

    const cx = W / 2;
    const cy = H * 0.46;
    const reduced = this.reduced;

    // ---- ground: warm radial well toward the centre of gravity ----
    ctx.clearRect(0, 0, W, H);
    const depthDim = 1 - this.depth * 0.55;
    const g = ctx.createRadialGradient(cx, cy, 0, cx, cy, Math.max(W, H) * 0.72);
    g.addColorStop(0, `rgba(26,24,18,${0.95 * depthDim})`);
    g.addColorStop(0.32, `rgba(16,16,22,${0.9 * depthDim})`);
    g.addColorStop(1, "rgba(7,8,10,0)");
    ctx.fillStyle = g;
    ctx.fillRect(0, 0, W, H);

    if (!this.snap) return;

    // ---- field scale: hold the orrery in frame as it accumulates ----
    let maxGrowth = 1;
    this.bodies.forEach((b) => {
      if (b.cap) maxGrowth = Math.max(maxGrowth, b.cap.total_pence / b.refTotal);
    });
    const zoom = Math.pow(maxGrowth, 0.32);
    const baseDist = Math.min(W, H) * 0.16;
    const distSpan = Math.min(W, H) * 0.32;

    // ---- centre of gravity: luminous gold core (p̄′) ----
    ctx.globalCompositeOperation = "lighter";
    const coreR = Math.min(W, H) * (0.07 + 0.016 * Math.sin(this.t * 1.3) + this._pulse * 0.03);
    const core = ctx.createRadialGradient(cx, cy, 0, cx, cy, coreR * 6.5);
    const coreA = 0.7 * depthDim;
    core.addColorStop(0, rgba(GOLD_HI, 1.0 * coreA));
    core.addColorStop(0.18, rgba(GOLD, 0.6 * coreA));
    core.addColorStop(0.5, rgba([120, 80, 30], 0.28 * coreA));
    core.addColorStop(1, rgba(GOLD, 0));
    ctx.fillStyle = core;
    ctx.beginPath();
    ctx.arc(cx, cy, coreR * 6.5, 0, Math.PI * 2);
    ctx.fill();
    // dense core
    ctx.fillStyle = rgba(GOLD_HI, 0.6 * depthDim);
    ctx.beginPath();
    ctx.arc(cx, cy, coreR * 0.42, 0, Math.PI * 2);
    ctx.fill();
    ctx.globalCompositeOperation = "source-over";

    // ---- faint gravitation guide rings ----
    ctx.strokeStyle = rgba(GOLD, 0.05 * depthDim);
    ctx.lineWidth = 1;
    for (let i = 1; i <= 3; i++) {
      ctx.beginPath();
      ctx.arc(cx, cy, baseDist + (distSpan / 3) * i * 0.8, 0, Math.PI * 2);
      ctx.stroke();
    }

    // ---- bodies ----
    const lingerP = 2.4;
    this.bodies.forEach((b) => {
      const c = b.cap;
      if (!c) return;
      const growth = c.total_pence / b.refTotal;
      b.growth = lerp(b.growth || growth, growth, 0.08);
      // accumulation → spiral outward + enlarge (relative to field zoom)
      const norm = clamp(Math.sqrt(b.growth) / zoom, 0.18, 1.25);
      const targetDist = baseDist + distSpan * norm;
      b.dist = lerp(b.dist || targetDist, targetDist, 0.06);
      if (!reduced) b.theta += b.orbitSpeed * dt * this.speed * (this.running ? 1 : 0);

      const bx = cx + Math.cos(b.theta) * b.dist;
      const by = cy + Math.sin(b.theta) * b.dist * 0.72; // ecliptic squash — keep the field a horizontal band

      const ringR = clamp(
        ((Math.min(W, H) * 0.062) * Math.sqrt(b.growth)) / Math.pow(zoom, 0.5),
        18,
        Math.min(W, H) * 0.17
      );
      b.ringR = ringR;
      b.sx = bx;
      b.sy = by;
      b.sr = ringR;

      const halted = c.status === "halted";
      const sum = c.money_pence + c.production_pence + c.commodity_pence || 1;
      const fm = c.money_pence / sum,
        fp = c.production_pence / sum,
        fc = c.commodity_pence / sum;

      // ring thickness swells with surplus share + loop pulse (dead labour towering)
      const surplusShare = clamp(c.surplus_pence / Math.max(1, c.total_pence), 0, 0.6);
      const sw = clamp(ringR * (0.16 + surplusShare * 0.5 + this._pulse * 0.12), 3, ringR * 0.7);
      const isHovered = this.hoveredId === c.id;
      const hoverDim =
        this.hoveredId === null ? 1 : isHovered ? 1 : 0.3;
      const dim = (halted ? 0.32 : depthDim) * hoverDim;

      // base track
      ctx.lineWidth = sw;
      ctx.lineCap = "butt";
      ctx.strokeStyle = rgba([22, 25, 34], 0.9 * dim);
      ctx.beginPath();
      ctx.arc(bx, by, ringR, 0, Math.PI * 2);
      ctx.stroke();

      if (isHovered) {
        ctx.globalCompositeOperation = "lighter";
        ctx.lineWidth = 2;
        ctx.strokeStyle = rgba(GOLD_HI, 0.5 * depthDim);
        ctx.beginPath();
        ctx.arc(bx, by, ringR + sw * 0.5 + 3, 0, Math.PI * 2);
        ctx.stroke();
        ctx.globalCompositeOperation = "source-over";
      }

      // three arcs (additive bloom = fat translucent + crisp)
      const arcs = [
        { col: GOLD, a0: 0, a1: fm },
        { col: RED, a0: fm, a1: fm + fp },
        { col: LEAD, a0: fm + fp, a1: 1 },
      ];
      const TAU = Math.PI * 2,
        gap = 0.012;
      for (const seg of arcs) {
        const s = -Math.PI / 2 + (seg.a0 + gap) * TAU;
        const e = -Math.PI / 2 + (seg.a1 - gap) * TAU;
        if (e <= s) continue;
        ctx.globalCompositeOperation = "lighter";
        ctx.lineWidth = sw * 2.6;
        ctx.strokeStyle = rgba(seg.col, 0.16 * dim);
        ctx.beginPath();
        ctx.arc(bx, by, ringR, s, e);
        ctx.stroke();
        ctx.lineWidth = sw * 1.5;
        ctx.strokeStyle = rgba(seg.col, 0.14 * dim);
        ctx.beginPath();
        ctx.arc(bx, by, ringR, s, e);
        ctx.stroke();
        ctx.globalCompositeOperation = "source-over";
        ctx.lineWidth = sw;
        ctx.strokeStyle = rgba(seg.col, 0.95 * dim);
        ctx.beginPath();
        ctx.arc(bx, by, ringR, s, e);
        ctx.stroke();
      }

      // value-dots travelling the circuit
      const lap = (c.turnover_number > 0 ? c.turnover_number : 1) / 8;
      for (let i = 0; i < b.dots.length; i++) {
        if (!reduced && this.running) b.dots[i] = (b.dots[i] + dt * lap * this.speed) % 1;
        const frac = pacedFrac(b.dots[i], fm, fp, fc, lingerP);
        const a = -Math.PI / 2 + frac * TAU;
        const dx = bx + Math.cos(a) * ringR;
        const dy = by + Math.sin(a) * ringR;
        ctx.globalCompositeOperation = "lighter";
        ctx.fillStyle = rgba(BONE, 0.5 * dim);
        ctx.beginPath();
        ctx.arc(dx, dy, sw * 0.9, 0, Math.PI * 2);
        ctx.fill();
        ctx.globalCompositeOperation = "source-over";
        ctx.fillStyle = rgba(BONE, 0.95 * dim);
        ctx.beginPath();
        ctx.arc(dx, dy, Math.max(1.6, sw * 0.42), 0, Math.PI * 2);
        ctx.fill();
      }

      if (halted) {
        ctx.fillStyle = rgba([138, 133, 120], 0.5);
        ctx.font = "9px 'IBM Plex Mono', monospace";
        ctx.textAlign = "center";
        ctx.fillText("halted", bx, by + ringR + 12);
      }
    });

    // ---- the alienation loop: gold rising from the abode below ----
    const spawnRate = (0.6 + this._pulse * 14) * (reduced ? 0 : 1);
    this._spawnAcc += spawnRate * dt;
    while (this._spawnAcc >= 1) {
      this._spawnAcc -= 1;
      this.motes.push({
        x: cx + (Math.random() - 0.5) * W * 0.7,
        y: H + 8,
        vy: -(40 + Math.random() * 60),
        life: 1,
        r: 1.2 + Math.random() * 1.8,
      });
    }
    ctx.globalCompositeOperation = "lighter";
    for (const m of this.motes) {
      m.y += m.vy * dt;
      // gentle pull toward the centre of gravity
      m.x += (cx - m.x) * 0.4 * dt;
      m.life -= dt * 0.32;
      const a = clamp(m.life, 0, 1) * 0.7 * depthDim;
      ctx.fillStyle = rgba(GOLD_HI, a);
      ctx.beginPath();
      ctx.arc(m.x, m.y, m.r, 0, Math.PI * 2);
      ctx.fill();
    }
    ctx.globalCompositeOperation = "source-over";
    this.motes = this.motes.filter((m) => m.life > 0 && m.y > -20);

    // ---- vignette ----
    const vg = ctx.createRadialGradient(cx, cy, Math.min(W, H) * 0.3, cx, cy, Math.max(W, H) * 0.7);
    vg.addColorStop(0, "rgba(7,8,10,0)");
    vg.addColorStop(1, `rgba(5,5,7,${0.55 + this.depth * 0.3})`);
    ctx.fillStyle = vg;
    ctx.fillRect(0, 0, W, H);

    // ---- hover readout (drawn last, above the field) ----
    if (this.hoveredId) {
      const hb = this.bodies.get(this.hoveredId);
      if (hb && hb.cap) this._drawCard(hb);
    }
  }

  /** Canvas readout for the hovered capital, anchored beside its ring. */
  private _drawCard(b: Body) {
    const c = b.cap;
    if (!c) return;
    const ctx = this.ctx;
    const W = this._w;

    const pad = 12;
    const cw = 196;
    const lineH = 17;
    const rows = 7; // C, p', M, P, C', Σs, n
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
    // the three coexisting arcs of the circuit, each tinted to its ring colour,
    // one per row so values never collide at large magnitudes
    row("M  money", formatPence(c.money_pence), GOLD_HI);
    row("P  production", formatPence(c.production_pence), RED);
    row("C′ commodity", formatPence(c.commodity_pence), LEAD);
    row("Σs surplus", formatPence(c.surplus_pence), GOLD);
    row("n  turnover", c.turnover_number.toFixed(1), BONE);

    // halted tag, top-right corner (only non-numeric element)
    if (c.status === "halted") {
      ctx.font = "9px 'IBM Plex Mono', monospace";
      ctx.textAlign = "right";
      ctx.fillStyle = "rgba(138,133,120,0.8)";
      ctx.fillText("halted", x + cw - 8, y + 12);
    }
    ctx.textAlign = "left"; // restore default so later text draws aren't right-aligned
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

  centre() {
    return { x: this._w / 2, y: this._h * 0.46 };
  }
}
