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
