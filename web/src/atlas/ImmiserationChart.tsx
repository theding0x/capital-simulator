import { useMemo } from "react";
import type { GeneralLawTrendPoint } from "../types";
import { formatPence, formatBP, seriesArea, seriesLine } from "./animation";

/** The immiseration trend: s/v rising (gold), wage falling (red dashed), reserve
 *  army swelling (lead area). The wage/s-v crossing IS the immiseration story. */
export function ImmiserationChart({ series }: { series: GeneralLawTrendPoint[] }) {
  const W = 680,
    H = 210,
    padL = 16,
    padR = 16,
    padT = 18,
    padB = 26;
  const geom = useMemo(() => {
    if (!series || series.length < 2) return null;
    const sv = series.map((p) => p.rate_of_exploitation_bp);
    const wg = series.map((p) => p.wage_pence);
    const ra = series.map((p) => p.reserve_army_count);
    const n = series.length;
    const innerW = W - padL - padR,
      innerH = H - padT - padB;
    const x = (i: number) => padL + (n > 1 ? (i / (n - 1)) * innerW : 0);
    const norm = (arr: number[]) => {
      const lo = Math.min(...arr),
        hi = Math.max(...arr),
        s = hi - lo || 1;
      return (v: number) => (v - lo) / s;
    };
    const ny = (arr: number[]) => {
      const f = norm(arr);
      return (v: number) => padT + innerH * (1 - f(v));
    };
    const yS = ny(sv),
      yW = ny(wg),
      yR = ny(ra);
    return {
      svPath: seriesLine(sv, x, yS),
      wgPath: seriesLine(wg, x, yW),
      raArea: seriesArea(ra, x, yR, padT + innerH),
      last: series[n - 1],
      n,
      svDot: [x(n - 1), yS(sv[n - 1])] as [number, number],
      wgDot: [x(n - 1), yW(wg[n - 1])] as [number, number],
    };
  }, [series]);

  if (!geom) return <p className="chart-empty">The law has not yet run — start the engine.</p>;
  const innerH = H - padT - padB;
  return (
    <div className="chart">
      <svg
        viewBox={`0 0 ${W} ${H}`}
        width="100%"
        preserveAspectRatio="none"
        role="img"
        aria-label="immiseration trend: exploitation rising, wage falling, reserve army swelling"
      >
        <defs>
          <linearGradient id="raFill" x1="0" y1="0" x2="0" y2="1">
            <stop offset="0%" stopColor="rgba(74,90,138,0.26)" />
            <stop offset="100%" stopColor="rgba(74,90,138,0)" />
          </linearGradient>
        </defs>
        {[0.25, 0.5, 0.75].map((f) => (
          <line
            key={f}
            x1={padL}
            x2={W - padR}
            y1={padT + innerH * f}
            y2={padT + innerH * f}
            stroke="rgba(255,255,255,0.05)"
            strokeWidth="1"
          />
        ))}
        <path d={geom.raArea} fill="url(#raFill)" stroke="none" />
        <path d={geom.wgPath} fill="none" stroke="#c0392b" strokeWidth="2" strokeDasharray="4 3" />
        <path d={geom.svPath} fill="none" stroke="#c8a240" strokeWidth="2.4" />
        <circle cx={geom.svDot[0]} cy={geom.svDot[1]} r="3.4" fill="#e8c660" />
        <circle cx={geom.wgDot[0]} cy={geom.wgDot[1]} r="3" fill="#d44030" />
      </svg>
      <div className="chart-legend">
        <span className="lg gold">s/v {formatBP(geom.last.rate_of_exploitation_bp)} ↑</span>
        <span className="lg red">wage {formatPence(geom.last.wage_pence)} ↓</span>
        <span className="lg lead">reserve army {geom.last.reserve_army_count} ↑</span>
        <span className="lg dim">period {geom.last.period}</span>
      </div>
    </div>
  );
}
