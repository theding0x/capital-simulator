import { useCallback, useEffect, useRef, useState } from "react";
import "./atlas.css";
import { useSnapshot } from "./useSnapshot";
import { AtlasSurface } from "./surface";
import { VitalSigns } from "./VitalSigns";
import { Abode } from "./Abode";
import { Transport } from "./TickHeartbeat";
import { clamp, formatBP } from "./animation";
import { loadPrefs, savePrefs } from "./prefs";

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

/** The Observatory — one continuous vertical world: the orrery surface above,
 *  the hidden abode below; descending is literal travel through a gilded gate. */
export default function Atlas() {
  const prefersReduced =
    typeof window !== "undefined" &&
    window.matchMedia &&
    window.matchMedia("(prefers-reduced-motion: reduce)").matches;

  const [running, setRunning] = useState(true);
  const [speed, setSpeed] = useState(() => loadPrefs().speed);
  const [reduced, setReduced] = useState(() => loadPrefs().reduced || !!prefersReduced);
  const [depth, setDepth] = useState(0);

  // Advance the session's in-memory run by `speed` periods per poll while running.
  const { snapshot } = useSnapshot(running ? speed : 0);

  // Persist UI preferences across reloads (the run itself does not persist).
  useEffect(() => {
    savePrefs({ speed, reduced });
  }, [speed, reduced]);

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
  useEffect(() => {
    if (snapshot) surfaceRef.current?.setSnapshot(snapshot);
  }, [snapshot]);
  useEffect(() => {
    surfaceRef.current?.setReduced(reduced);
  }, [reduced]);
  useEffect(() => {
    surfaceRef.current?.setSpeed(speed);
  }, [speed]);

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
    });
  }, [running]);

  const descend = () =>
    animateScroll(surfaceZoneRef.current?.offsetHeight || window.innerHeight, reduced ? 0 : 1150);
  const ascend = () => animateScroll(0, reduced ? 0 : 1000);

  const toggleRun = () => {
    const next = !running;
    setRunning(next);
    surfaceRef.current?.setRunning(next);
  };

  const descended = depth > 0.5;

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
          <a href="#/chapters">Chapters</a>
        </nav>
      </header>

      <aside className="rail">
        {snapshot && <VitalSigns vitals={snapshot.aggregate} count={snapshot.capitals.length} />}
        <button
          className={"descend-btn" + (descended ? " open" : "")}
          data-testid="threshold"
          onClick={descended ? ascend : descend}
        >
          <span className="db-arrow">{descended ? "↑" : "↓"}</span>
          <span className="db-label">
            {descended ? "Ascend to the surface" : "Descend into production"}
          </span>
        </button>
        <p className="rail-foot">
          An observatory of the circuit of capital — <span className="i">M—C…P…C′—M′</span> — at the
          scale of many capitals.
        </p>
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
              <p className="surface-sub">
                Each capital a ring of three coexisting arcs — <span className="g">M</span> money,{" "}
                <span className="r">P</span> production, <span className="l">C′</span> commodity —
                value travelling its circumference. They spiral outward as they accumulate.
              </p>
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
    </div>
  );
}
