import type { FieldCapital } from "../types";
import { arcFractions, orbitRadius, spinSeconds } from "./animation";

interface OrbitProps {
  capital: FieldCapital;
  maxTotal: number;
  /** Animation tempo multiplier (1 = base, 5 = 5× faster). */
  speed: number;
}

const GOLD = "#c8a240";
const RED = "#c0392b";
const LEAD = "#4a5a8a";

/** One faithful orbit: three coexisting arcs (M/P/C) with value-dots circulating. */
export function Orbit({ capital, maxTotal, speed }: OrbitProps) {
  const r = orbitRadius(capital.total_pence, maxTotal);
  const pad = 10;
  const size = (r + pad) * 2;
  const c = size / 2;
  const sw = Math.max(6, Math.round(r * 0.16));
  const [fm, fp] = arcFractions(capital);
  const fc = 1 - fm - fp;
  const gap = 1.5; // pathLength units between arcs

  // pathLength=100; lay arcs M, P, C clockwise from top with small gaps.
  const lm = Math.max(0, fm * 100 - gap);
  const lp = Math.max(0, fp * 100 - gap);
  const lc = Math.max(0, fc * 100 - gap);
  const offM = 0;
  const offP = -(fm * 100);
  const offC = -((fm + fp) * 100);

  const dur = Number((spinSeconds(capital.id) / Math.max(1, speed)).toFixed(2));
  const halted = capital.status === "halted";

  return (
    <div
      className={"atlas-orbit" + (halted ? " halted" : "")}
      title={`${capital.id.slice(0, 8)} · total ${capital.total_pence} · ${capital.status}`}
    >
      <svg width={size} height={size} viewBox={`0 0 ${size} ${size}`}>
        <circle cx={c} cy={c} r={r} fill="none" stroke="#161922" strokeWidth={sw} />
        <g transform={`rotate(-90 ${c} ${c})`}>
          <circle cx={c} cy={c} r={r} fill="none" stroke={GOLD} strokeWidth={sw}
            pathLength={100} strokeDasharray={`${lm} ${100 - lm}`} strokeDashoffset={offM} />
          <circle cx={c} cy={c} r={r} fill="none" stroke={RED} strokeWidth={sw}
            pathLength={100} strokeDasharray={`${lp} ${100 - lp}`} strokeDashoffset={offP} />
          <circle cx={c} cy={c} r={r} fill="none" stroke={LEAD} strokeWidth={sw}
            pathLength={100} strokeDasharray={`${lc} ${100 - lc}`} strokeDashoffset={offC} />
        </g>
        {[0, 1, 2, 3].map((i) => (
          <g key={i} className="atlas-flow"
            style={{ animationDuration: `${dur}s`, animationDelay: `-${(dur / 4) * i}s` }}>
            <circle cx={c} cy={c - r} r={Math.max(2, sw * 0.28)} fill="#f4ecd8" />
          </g>
        ))}
      </svg>
    </div>
  );
}
