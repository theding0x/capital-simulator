import type { AggregateVitals } from "../types";
import { formatBP, formatPence } from "./animation";

/** Aggregate readouts for the rail. */
export function VitalSigns({ vitals, capitalCount }: { vitals: AggregateVitals; capitalCount: number }) {
  return (
    <div>
      <div className="atlas-vital">
        <div className="k">Total social capital</div>
        <div className="v">{formatPence(vitals.total_social_capital_pence)}</div>
      </div>
      <div className="atlas-vital">
        <div className="k">Average rate of profit · p̄′</div>
        <div className="v gold">{formatBP(vitals.avg_rate_of_profit_bp)}</div>
      </div>
      <div className="atlas-vital">
        <div className="k">Surplus-value · ΣS</div>
        <div className="v">{formatPence(vitals.surplus_pence)}</div>
      </div>
      <div className="atlas-vital">
        <div className="k">Cost-price · ΣC (c+v)</div>
        <div className="v">{formatPence(vitals.cost_price_pence)}</div>
      </div>
      <div className="atlas-vital">
        <div className="k">Capitals in motion</div>
        <div className="v">{capitalCount}</div>
      </div>
    </div>
  );
}
