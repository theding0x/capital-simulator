import { useState } from "react";
import type { AbodeReadout, LeverUpdate } from "../types";
import { api } from "../api";
import { formatBP, formatPence } from "./animation";

interface LeversProps {
  abode: AbodeReadout;
}

/** The levers — perturb the live General Law and watch it respond over the next
 *  passes. The working day sets the rate of surplus-value (necessary↔surplus);
 *  the wage is the value of labour-power; α is the share of surplus
 *  re-accumulated. Slider state is local (seeded once); the law's *response*
 *  shows in the abode readout above, not in the slider position. */
export function Levers({ abode }: LeversProps) {
  const [surplus, setSurplus] = useState(abode.surplus_rate_base_bp);
  const [wage, setWage] = useState(abode.base_wage_pence);
  const [accum, setAccum] = useState(abode.accumulation_rate_bp);

  const push = (u: LeverUpdate) => {
    void api.setObservatoryLevers(u).catch(() => {
      /* the 2s poll reflects the actual persisted state */
    });
  };

  return (
    <section className="abode-card levers" data-testid="levers">
      <div className="abode-card-k">The levers — perturb the law</div>

      <label className="lever">
        <span>Working day · rate of surplus-value <b>{formatBP(surplus)}</b></span>
        <input
          type="range" min={2000} max={40000} step={500} value={surplus}
          data-testid="lever-workingday"
          onChange={(e) => {
            const v = Number(e.target.value);
            setSurplus(v);
            push({ surplus_rate_base_bp: v });
          }}
        />
      </label>

      <label className="lever">
        <span>Wage · value of labour-power <b>{formatPence(wage)}</b></span>
        <input
          type="range" min={500} max={6000} step={100} value={wage}
          data-testid="lever-wage"
          onChange={(e) => {
            const v = Number(e.target.value);
            setWage(v);
            push({ base_wage_pence: v });
          }}
        />
      </label>

      <label className="lever">
        <span>Accumulation rate · α <b>{formatBP(accum)}</b></span>
        <input
          type="range" min={0} max={10000} step={250} value={accum}
          data-testid="lever-accumulation"
          onChange={(e) => {
            const v = Number(e.target.value);
            setAccum(v);
            push({ accumulation_rate_bp: v });
          }}
        />
      </label>
    </section>
  );
}
