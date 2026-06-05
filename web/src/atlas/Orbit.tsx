import { useEffect, useRef } from "react";
import type { FieldCapital } from "../types";
import { arcFractions, lapRateFor, pacedAngle, targetScale } from "./animation";

interface OrbitProps {
  capital: FieldCapital;
  maxTotal: number;
  /** Animation tempo multiplier (1 = base). */
  speed: number;
}

const GOLD = "#c8a240";
const RED = "#c0392b";
const LEAD = "#4a5a8a";
const BASE_R = 60; // base ring radius in svg units; growth applied via CSS scale
const PAD = 12;
const DOT_COUNT = 3;

/** One faithful orbit: a STILL ring of M/P/C arcs; value-dots travel it, lapping
 *  once per turnover and lingering in production; the whole orbit scales toward
 *  its growing magnitude. */
export function Orbit({ capital, maxTotal, speed }: OrbitProps) {
  const size = (BASE_R + PAD) * 2;
  const c = size / 2;
  const sw = 11;
  const [fm, fp] = arcFractions(capital);
  const fc = 1 - fm - fp;
  const gap = 1.5;
  const lm = Math.max(0, fm * 100 - gap);
  const lp = Math.max(0, fp * 100 - gap);
  const lc = Math.max(0, fc * 100 - gap);

  const wrapRef = useRef<HTMLDivElement>(null);
  const dotRefs = useRef<(SVGCircleElement | null)[]>([]);
  // Keep latest values in a ref so the rAF loop reads fresh data without restart.
  const live = useRef({ fm, fp, fc, lap: lapRateFor(capital.turnover_number, speed) });
  live.current = { fm, fp, fc, lap: lapRateFor(capital.turnover_number, speed) };

  // Growth: set the wrapper's CSS scale toward the target (CSS transition eases it).
  useEffect(() => {
    if (wrapRef.current) {
      wrapRef.current.style.transform = `scale(${targetScale(capital.total_pence, maxTotal).toFixed(3)})`;
    }
  }, [capital.total_pence, maxTotal]);

  // Dots travel the ring; rAF updates positions imperatively (no re-render).
  useEffect(() => {
    let raf = 0;
    let prev = performance.now();
    const phase = [0, 1 / DOT_COUNT, 2 / DOT_COUNT]; // evenly spaced
    const tick = (now: number) => {
      const dt = Math.min(0.05, (now - prev) / 1000);
      prev = now;
      const cur = live.current;
      for (let i = 0; i < DOT_COUNT; i++) {
        phase[i] = (phase[i] + dt * cur.lap) % 1;
        const a = pacedAngle(phase[i], cur.fm, cur.fp, cur.fc) - Math.PI / 2; // 0 at top
        const el = dotRefs.current[i];
        if (el) {
          el.setAttribute("cx", (c + BASE_R * Math.cos(a)).toFixed(2));
          el.setAttribute("cy", (c + BASE_R * Math.sin(a)).toFixed(2));
        }
      }
      raf = requestAnimationFrame(tick);
    };
    raf = requestAnimationFrame(tick);
    return () => cancelAnimationFrame(raf);
  }, [c]);

  const halted = capital.status === "halted";

  return (
    <div
      ref={wrapRef}
      className={"atlas-orbit" + (halted ? " halted" : "")}
      style={{ transition: "transform 1.6s ease-out", willChange: "transform" }}
      title={`${capital.id.slice(0, 8)} · total ${capital.total_pence} · turnover ${capital.turnover_number}× · ${capital.status}`}
    >
      <svg width={size} height={size} viewBox={`0 0 ${size} ${size}`}>
        <circle cx={c} cy={c} r={BASE_R} fill="none" stroke="#161922" strokeWidth={sw} />
        <g transform={`rotate(-90 ${c} ${c})`}>
          <circle cx={c} cy={c} r={BASE_R} fill="none" stroke={GOLD} strokeWidth={sw}
            pathLength={100} strokeDasharray={`${lm} ${100 - lm}`} strokeDashoffset={0} />
          <circle cx={c} cy={c} r={BASE_R} fill="none" stroke={RED} strokeWidth={sw}
            pathLength={100} strokeDasharray={`${lp} ${100 - lp}`} strokeDashoffset={-(fm * 100)} />
          <circle cx={c} cy={c} r={BASE_R} fill="none" stroke={LEAD} strokeWidth={sw}
            pathLength={100} strokeDasharray={`${lc} ${100 - lc}`} strokeDashoffset={-((fm + fp) * 100)} />
        </g>
        {Array.from({ length: DOT_COUNT }).map((_, i) => (
          <circle key={i} ref={(el) => { dotRefs.current[i] = el; }}
            cx={c} cy={c - BASE_R} r={3.2} fill="#f4ecd8"
            style={{ filter: "drop-shadow(0 0 3px rgba(244,236,216,.8))" }} />
        ))}
      </svg>
    </div>
  );
}
