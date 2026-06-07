import { useId, useMemo } from "react";
import type { DistSeriesPoint } from "../types";
import { formatBP, formatPence, seriesArea, seriesLine } from "./animation";

interface TRPFChartProps {
  series: DistSeriesPoint[];
}

const W = 680, H = 200, PAD_L = 16, PAD_R = 16, PAD_T = 18, PAD_B = 22;

/** SVG chart of the tendential fall of the rate of profit — p̄′ line over mass area. */
export function TRPFChart({ series }: TRPFChartProps) {
  const uid = useId();
  const gradId = `massFill-${uid.replace(/:/g, "")}`;

  const geom = useMemo(() => {
    if (!series || series.length < 2) return null;
    const n = series.length;
    const innerW = W - PAD_L - PAD_R;
    const innerH = H - PAD_T - PAD_B;
    const x = (i: number) => PAD_L + (i / (n - 1)) * innerW;
    const pp = series.map((p) => p.pprime_bp);
    const mass = series.map((p) => p.profit_mass_pence);
    const ppMax = Math.max(...pp, 1);
    const massMax = Math.max(...mass, 1);
    const yPP = (v: number) => PAD_T + innerH * (1 - v / (ppMax * 1.1));
    const yM = (v: number) => PAD_T + innerH * (1 - v / (massMax * 1.1));
    const crises = series
      .map((p, i) => (p.crisis ? x(i) : null))
      .filter((v): v is number => v !== null);
    const last = series[n - 1];
    return {
      ppPath: seriesLine(pp, x, yPP),
      massArea: seriesArea(mass, x, yM, PAD_T + innerH),
      crises,
      last,
      ppDot: [x(n - 1), yPP(pp[n - 1])] as [number, number],
      innerH,
    };
  }, [series]);

  if (!geom) {
    return (
      <p className="chart-empty">
        The totality has not yet run — start the engine.
      </p>
    );
  }

  return (
    <div className="chart">
      <svg
        viewBox={`0 0 ${W} ${H}`}
        width="100%"
        preserveAspectRatio="none"
        role="img"
        aria-label="tendential fall of the rate of profit; mass of profit rising; crises restore the rate"
      >
        <defs>
          <linearGradient id={gradId} x1="0" y1="0" x2="0" y2="1">
            <stop offset="0%" stopColor="rgba(74,90,138,0.24)" />
            <stop offset="100%" stopColor="rgba(74,90,138,0)" />
          </linearGradient>
        </defs>
        {[0.25, 0.5, 0.75].map((f) => (
          <line
            key={f}
            x1={PAD_L}
            x2={W - PAD_R}
            y1={PAD_T + geom.innerH * f}
            y2={PAD_T + geom.innerH * f}
            stroke="rgba(255,255,255,0.05)"
          />
        ))}
        {geom.crises.map((cx, i) => (
          <line
            key={i}
            x1={cx}
            x2={cx}
            y1={PAD_T}
            y2={PAD_T + geom.innerH}
            stroke="rgba(192,57,43,0.4)"
            strokeWidth="1.5"
          />
        ))}
        <path d={geom.massArea} fill={`url(#${gradId})`} stroke="none" />
        <path d={geom.ppPath} fill="none" stroke="#c8a240" strokeWidth="2.4" />
        <circle cx={geom.ppDot[0]} cy={geom.ppDot[1]} r="3.4" fill="#e8c660" />
      </svg>
      <div className="chart-legend">
        <span className="lg gold">p&#x0305;&#x2032; {formatBP(geom.last.pprime_bp)}</span>
        <span className="lg lead">mass of profit {formatPence(geom.last.profit_mass_pence)} &uarr;</span>
        <span className="lg red">&#x2502; crisis &mdash; c devalued, rate restored</span>
        <span className="lg dim">period {geom.last.period}</span>
      </div>
    </div>
  );
}
