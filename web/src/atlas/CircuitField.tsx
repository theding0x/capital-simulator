import type { ObservatorySnapshot } from "../types";
import { Orbit } from "./Orbit";
import { formatBP } from "./animation";

interface CircuitFieldProps {
  snapshot: ObservatorySnapshot;
  speed: number;
}

/** The field of orbits, with the average rate of profit as centre of gravity. */
export function CircuitField({ snapshot, speed }: CircuitFieldProps) {
  const maxTotal = snapshot.capitals.reduce((m, c) => Math.max(m, c.total_pence), 0);
  return (
    <div className="atlas-field-wrap">
      <div className="atlas-field-centre">
        <div className="pbar">{formatBP(snapshot.aggregate.avg_rate_of_profit_bp)}</div>
        <div className="lbl">p̄′ · centre of gravity</div>
      </div>
      <div className="atlas-field" data-testid="atlas-field">
        {snapshot.capitals.map((cap) => (
          <Orbit key={cap.id} capital={cap} maxTotal={maxTotal} speed={speed} />
        ))}
      </div>
    </div>
  );
}
