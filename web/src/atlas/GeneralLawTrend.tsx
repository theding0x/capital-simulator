import type { GeneralLawTrendPoint } from "../types";
import { sparklinePath, formatBP } from "./animation";

interface TrendProps {
  series: GeneralLawTrendPoint[];
}

const W = 260;
const H = 64;
const GOLD = "#c8a240";
const RED = "#c0392b";

/** The immiseration trend: rate of exploitation (s/v) rising and the paid wage
 *  falling across the recorded periods of the general law. */
export function GeneralLawTrend({ series }: TrendProps) {
  if (series.length < 2) {
    return <p className="abode-hint">The law has not yet run — start the engine.</p>;
  }
  const exploitation = series.map((p) => p.rate_of_exploitation_bp);
  const wages = series.map((p) => p.wage_pence);
  const latest = series[series.length - 1];
  return (
    <div className="abode-trend">
      <svg width={W} height={H} viewBox={`0 0 ${W} ${H}`} role="img"
        aria-label="immiseration trend: rising exploitation, falling wage">
        <path d={sparklinePath(exploitation, W, H)} fill="none" stroke={GOLD} strokeWidth={2} />
        <path d={sparklinePath(wages, W, H)} fill="none" stroke={RED} strokeWidth={2}
          strokeDasharray="3 2" />
      </svg>
      <div className="abode-trend-legend">
        <span style={{ color: GOLD }}>s/v {formatBP(latest.rate_of_exploitation_bp)}</span>
        <span style={{ color: RED }}>wage ↓</span>
        <span className="abode-hint">period {latest.period}</span>
      </div>
    </div>
  );
}
