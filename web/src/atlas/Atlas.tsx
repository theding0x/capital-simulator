import { useCallback, useEffect, useRef, useState } from "react";
import "./atlas.css";
import { useSnapshot } from "./useSnapshot";
import { AtlasSurface } from "./surface";
import { VitalSigns } from "./VitalSigns";
import { Abode } from "./Abode";
import { Circulation } from "./Circulation";
import { Totality } from "./Totality";
import { Transport } from "./TickHeartbeat";
import { clamp, formatBP } from "./animation";
import { loadPrefs, savePrefs } from "./prefs";
import { SpineMatrix } from "./SpineMatrix";
import { AuthBar } from "../components/AuthBar";
import { Gloss } from "./Gloss";
import { GLOSS_MAP } from "./glossContent";
import { GenesisTimeline } from "./GenesisTimeline";

const STRATA = [
  { key: "surface", label: "The surface · the circuit", short: "Surface" },
  { key: "abode", label: "Vol I · the abode of production", short: "Vol I" },
  { key: "circulation", label: "Vol II · circulation & reproduction", short: "Vol II" },
  { key: "totality", label: "Vol III · the totality", short: "Vol III" },
  { key: "genesis", label: "The genesis · primitive accumulation", short: "Genesis" },
] as const;

/** Animated scrollTop tween (inOutCubic) on the stage container. */
function useAnimatedScroll(stageRef: React.RefObject<HTMLDivElement | null>) {
  const raf = useRef(0);
  return useCallback(
    (to: number, ms: number, done?: () => void) => {
      const el = stageRef.current;
      if (!el) return;
      cancelAnimationFrame(raf.current);
      const from = el.scrollTop,
        d = to - from,
        t0 = performance.now();
      const ease = (x: number) => (x < 0.5 ? 4 * x * x * x : 1 - Math.pow(-2 * x + 2, 3) / 2);
      if (ms <= 0) {
        el.scrollTop = to;
        done?.();
        return;
      }
      const step = (now: number) => {
        const p = Math.min(1, (now - t0) / ms);
        el.scrollTop = from + d * ease(p);
        if (p < 1) raf.current = requestAnimationFrame(step);
        else done?.();
      };
      raf.current = requestAnimationFrame(step);
    },
    [stageRef]
  );
}

interface AtlasProps {
  /** Pre-open the SpineMatrix index overlay on mount, and optionally deep-link to a chapter. */
  initialChapterId?: string;
}

/** The Observatory — one continuous vertical world descending through four strata:
 *  surface (orrery), abode (Vol I), circulation (Vol II), totality (Vol III). */
export default function Atlas({ initialChapterId }: AtlasProps = {}) {
  const prefersReduced =
    typeof window !== "undefined" &&
    window.matchMedia &&
    window.matchMedia("(prefers-reduced-motion: reduce)").matches;

  const [running, setRunning] = useState(true);
  const [speed, setSpeed] = useState(() => loadPrefs().speed);
  const [reduced, setReduced] = useState(() => loadPrefs().reduced || !!prefersReduced);
  const [depth, setDepth] = useState(0);
  const [level, setLevel] = useState(0);
  // SpineMatrix index overlay: open when initialChapterId is set, or via rail button.
  const [showIndex, setShowIndex] = useState(!!initialChapterId);

  // Advance the session's in-memory run by `speed` periods per poll while running.
  const { snapshot } = useSnapshot(running ? speed : 0);

  // Persist UI preferences across reloads (the run itself does not persist).
  useEffect(() => {
    savePrefs({ speed, reduced });
  }, [speed, reduced]);

  const surfaceRef = useRef<AtlasSurface | null>(null);
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const stageRef = useRef<HTMLDivElement>(null);

  // One ref per stratum: surface / abode / circulation / totality.
  const surfaceZoneRef = useRef<HTMLDivElement>(null);
  const abodeZoneRef = useRef<HTMLDivElement>(null);
  const circulationZoneRef = useRef<HTMLDivElement>(null);
  const totalityZoneRef = useRef<HTMLDivElement>(null);
  const genesisZoneRef = useRef<HTMLDivElement>(null);
  const zoneRefs = [
    surfaceZoneRef,
    abodeZoneRef,
    circulationZoneRef,
    totalityZoneRef,
    genesisZoneRef,
  ];

  const animateScroll = useAnimatedScroll(stageRef);

  // Instantiate the canvas controller once.
  useEffect(() => {
    if (!canvasRef.current) return;
    const surf = new AtlasSurface(canvasRef.current);
    surfaceRef.current = surf;
    surf.setReduced(!!prefersReduced);
    surf.start();
    return () => surf.stop();
  }, []); // eslint-disable-line react-hooks/exhaustive-deps

  // Feed the controller each new snapshot.
  useEffect(() => {
    if (snapshot) surfaceRef.current?.setSnapshot(snapshot);
  }, [snapshot]);
  useEffect(() => {
    surfaceRef.current?.setReduced(reduced);
  }, [reduced]);
  useEffect(() => {
    surfaceRef.current?.setSpeed(speed);
  }, [speed]);

  const goTo = useCallback(
    (i: number) => {
      const idx = clamp(i, 0, STRATA.length - 1);
      const ref = zoneRefs[idx];
      const target = idx === 0 ? 0 : (ref.current?.offsetTop ?? 0);
      animateScroll(target, reduced ? 0 : 1100);
    },
    // eslint-disable-next-line react-hooks/exhaustive-deps
    [animateScroll, reduced]
  );

  const descend = () => goTo(level + 1);
  const ascend = () => goTo(level - 1);

  const depthRaf = useRef(0);
  const onScroll = useCallback(() => {
    cancelAnimationFrame(depthRaf.current);
    depthRaf.current = requestAnimationFrame(() => {
      const el = stageRef.current;
      if (!el) return;
      const h = surfaceZoneRef.current?.offsetHeight || window.innerHeight;
      const d = clamp(el.scrollTop / (h * 0.9), 0, 1);
      setDepth(d);
      surfaceRef.current?.setDepth(d);
      surfaceRef.current?.setRunning(running && d < 0.99);

      // Determine which stratum we are in.
      const mid = el.scrollTop + window.innerHeight * 0.4;
      let lv = 0;
      zoneRefs.forEach((r, i) => {
        if (r.current && r.current.offsetTop <= mid) lv = i;
      });
      setLevel(lv);
    });
  }, [running]); // eslint-disable-line react-hooks/exhaustive-deps

  const toggleRun = () => {
    const next = !running;
    setRunning(next);
    surfaceRef.current?.setRunning(next);
  };

  const atBottom = level >= STRATA.length - 1;

  return (
    <div className="atlas">
      <header className="topbar">
        <span className="brand">
          <span className="brand-mark">C</span> Capital Simulator{" "}
          <span className="brand-sub">· Atlas</span>
        </span>
        <nav className="nav">
          <a className="active" href="#/">
            Atlas
          </a>
          <button
            className={"nav-btn" + (showIndex ? " active" : "")}
            onClick={() => setShowIndex((v) => !v)}
            aria-expanded={showIndex}
            aria-controls="spine-index-overlay"
            aria-label="Open the chapter index — SpineMatrix"
          >
            Index
          </button>
          <a href="#/chapters">Chapters</a>
        </nav>
        <AuthBar />
      </header>

      <aside className="rail">
        {snapshot && <VitalSigns vitals={snapshot.aggregate} count={snapshot.capitals.length} />}
        <nav className="strata-nav" aria-label="Strata">
          {STRATA.map((s, i) => (
            <button
              key={s.key}
              className={"strata-item" + (level === i ? " on" : "")}
              onClick={() => goTo(i)}
            >
              <span className="si-depth">{i === 0 ? "▲" : "▼"}</span>
              <span className="si-lbl">{s.label}</span>
            </button>
          ))}
        </nav>
        <button
          className={"descend-btn" + (level > 0 ? " open" : "")}
          data-testid="threshold"
          onClick={atBottom ? ascend : descend}
        >
          <span className="db-arrow">{atBottom ? "↑" : "↓"}</span>
          <span className="db-label">
            {atBottom ? "Ascend toward the surface" : "Descend — explain deeper"}
          </span>
        </button>
        <button
          className={"spine-rail-btn" + (showIndex ? " active" : "")}
          onClick={() => setShowIndex((v) => !v)}
          aria-expanded={showIndex}
          aria-label="Open the chapter index — SpineMatrix"
        >
          &#9776; Chapter index
        </button>
        <p className="rail-foot">
          The same circuit <span className="i">M—C…P…C′—M′</span>, explained ever deeper
          — production, then circulation, then the totality.
        </p>
      </aside>

      <main className="stage" ref={stageRef} onScroll={onScroll}>
        <div className="world">
          {/* ── SURFACE ── */}
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
              <p className="surface-sub">
                Each capital a ring of three coexisting arcs — <span className="g">M</span> money,{" "}
                <span className="r">P</span> production, <span className="l">C′</span> commodity —
                value travelling its circumference. They spiral outward as they accumulate.
              </p>
              <div className="surface-gloss">
                <Gloss
                  label={GLOSS_MAP["v1-ch01"].label}
                  quote={GLOSS_MAP["v1-ch01"].quote}
                  citation={GLOSS_MAP["v1-ch01"].citation}
                  form={GLOSS_MAP["v1-ch01"].form}
                />
              </div>
            </div>
            <div
              className="gate"
              style={{ ["--cross"]: clamp((depth - 0.4) / 0.5, 0, 1) } as React.CSSProperties}
            >
              <div className="gate-half left"></div>
              <div className="gate-half right"></div>
              <button className="gate-notice" onClick={descend}>
                <span className="gate-no">No admittance except on business</span>
                <span className="gate-cite">Capital, Vol. I · Ch. 6 — leave the noisy sphere; descend ↓</span>
              </button>
            </div>
          </section>

          {/* ── ABODE (Vol I) ── */}
          <section className="zone-abode" ref={abodeZoneRef}>
            <div className="abode-seam"></div>
            <div className="abode-head">
              <div className="abode-eyebrow">Vol I · the hidden abode of production</div>
              <h2 className="abode-title">Where surplus is pumped from living labour</h2>
            </div>
            {snapshot && <Abode abode={snapshot.abode} />}
          </section>

          {/* ── CIRCULATION (Vol II) ── */}
          <section className="zone-circulation" ref={circulationZoneRef}>
            <div className="strat-seam"></div>
            <div className="strat-head">
              <div className="ab-eyebrow">Vol II · the process of circulation</div>
              <h2 className="strat-title">How the surplus comes back — reproduction of the whole</h2>
            </div>
            {snapshot?.circulation && (
              <Circulation circulation={snapshot.circulation} reduced={reduced} />
            )}
          </section>

          {/* ── TOTALITY (Vol III) ── */}
          <section className="zone-totality" ref={totalityZoneRef}>
            <div className="strat-seam lead"></div>
            <div className="strat-head">
              <div className="ab-eyebrow lead">Vol III · capitalist production as a whole</div>
              <h2 className="strat-title">The totality — where surplus distributes, and the rate falls</h2>
            </div>
            {snapshot?.distribution && (
              <Totality distribution={snapshot.distribution} />
            )}
          </section>

          {/* ── GENESIS (the historical floor) ── */}
          <section className="zone-genesis" ref={genesisZoneRef}>
            <div className="strat-seam red"></div>
            <div className="strat-head">
              <div className="ab-eyebrow red">Vol I &middot; Ch. 26&ndash;33 &mdash; so-called primitive accumulation</div>
              <h2 className="strat-title">The genesis &mdash; the bedrock the living world rests on</h2>
            </div>
            <GenesisTimeline />
          </section>
        </div>
      </main>

      <footer className="console">
        <Transport
          tick={snapshot?.tick ?? 0}
          running={running}
          onToggle={toggleRun}
          speed={speed}
          onSpeed={setSpeed}
          reduced={reduced}
          onReduced={setReduced}
        />
      </footer>

      <div className="descent-tint" style={{ opacity: depth * 0.5 }}></div>

      {/* SpineMatrix index overlay */}
      {showIndex && (
        <div
          id="spine-index-overlay"
          role="dialog"
          aria-modal="true"
          aria-label="Chapter index — the spine matrix"
          style={{
            position: "fixed",
            inset: 0,
            zIndex: 45,
            background: "rgba(7,8,10,0.88)",
            overflowY: "auto",
            padding: "calc(var(--topbar-h, 50px) + 16px) 24px 80px",
          }}
        >
          <div style={{ maxWidth: 860, margin: "0 auto" }}>
            <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center", marginBottom: 18 }}>
              <div>
                <div style={{ fontFamily: "var(--font-mono)", fontSize: "0.58rem", letterSpacing: "0.18em", textTransform: "uppercase", color: "var(--ink-dim)", marginBottom: 6 }}>
                  Atlas · Total Representation
                </div>
                <h2 style={{ fontFamily: "var(--font-display)", fontWeight: 700, fontSize: "1.5rem", margin: 0, color: "var(--ink)" }}>
                  The spine matrix
                </h2>
              </div>
              <button
                onClick={() => setShowIndex(false)}
                aria-label="Close chapter index"
                style={{
                  border: "1px solid var(--border)",
                  borderRadius: "var(--radius-sm, 4px)",
                  background: "transparent",
                  color: "var(--ink-muted)",
                  fontFamily: "var(--font-mono)",
                  fontSize: "1.1rem",
                  cursor: "pointer",
                  width: 36,
                  height: 36,
                  display: "grid",
                  placeItems: "center",
                }}
              >
                &times;
              </button>
            </div>
            <SpineMatrix initialChapterId={initialChapterId} />
          </div>
        </div>
      )}
    </div>
  );
}
