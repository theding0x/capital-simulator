import { useMemo } from "react";
import type { AbodeReadout } from "../types";
import { clamp, formatBP, formatPence } from "./animation";

/** The industrial reserve army as a reservoir of units; a pressure bar widths to
 *  reserve_army_pressure_bp. The disposable reserve made physical. */
export function ReserveArmy({ abode }: { abode: AbodeReadout }) {
  const count = abode.reserve_army_count;
  const pressure = abode.reserve_army_pressure_bp / 10000;
  const dots = Math.min(140, Math.round(count / 2.6));
  const cells = useMemo(() => Array.from({ length: 140 }), []);
  return (
    <section className="reserve">
      <div className="reserve-head">
        <div className="reserve-eyebrow">The industrial reserve army</div>
        <div className="reserve-meta">
          <span className="num">{count}</span> at the gates · wage pressure{" "}
          <span className="num">{formatBP(abode.reserve_army_pressure_bp)}</span> · paid wage{" "}
          <span className="num gold">{formatPence(abode.wage_pence)}</span>
        </div>
      </div>
      <div className="reservoir" aria-hidden="true">
        {cells.map((_, i) => (
          <span
            className={"unit" + (i < dots ? " on" : "")}
            key={i}
            style={{ transitionDelay: `${(i % 20) * 8}ms` }}
          ></span>
        ))}
      </div>
      <div className="reserve-fill" style={{ width: `${clamp(pressure * 100, 4, 100)}%` }}></div>
      <div className="reserve-foot">
        A lever of accumulation — the disposable reserve that presses the wage below the value of
        labour-power.
      </div>
    </section>
  );
}
