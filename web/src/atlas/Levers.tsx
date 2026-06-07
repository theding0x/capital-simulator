import { useState } from "react";
import type { AbodeReadout, LeverUpdate } from "../types";
import { api } from "../api";
import { formatBP, formatPence } from "./animation";

function Lever(props: {
  label: string;
  sub: string;
  value: number;
  display: string;
  min: number;
  max: number;
  step: number;
  testId: string;
  onChange: (v: number) => void;
}) {
  return (
    <label className="lever">
      <div className="lever-top">
        <span className="lever-label">{props.label}</span>
        <span className="lever-value">{props.display}</span>
      </div>
      <input
        type="range"
        min={props.min}
        max={props.max}
        step={props.step}
        value={props.value}
        data-testid={props.testId}
        onChange={(e) => props.onChange(Number(e.target.value))}
      />
      <span className="lever-sub">{props.sub}</span>
    </label>
  );
}

/** The levers — perturb the live AbodeState; the law ripples over the next polls.
 *  Slider state is local (seeded once); ranges track the real backend, not the
 *  design mock. The law's *response* shows in the abode readout above, not in the
 *  slider position. */
export function Levers({ abode }: { abode: AbodeReadout }) {
  const [surplus, setSurplus] = useState(abode.surplus_rate_base_bp);
  const [wage, setWage] = useState(abode.base_wage_pence);
  const [accum, setAccum] = useState(abode.accumulation_rate_bp);
  const push = (u: LeverUpdate) => {
    void api.setObservatoryLevers(u).catch(() => {
      /* next poll reflects truth */
    });
  };
  return (
    <section className="levers-card" data-testid="levers">
      <div className="levers-eyebrow">The levers — perturb the law</div>
      <div className="levers">
        <Lever
          label="Working day · rate of surplus-value"
          sub="lengthen the unpaid hours"
          value={surplus}
          display={formatBP(surplus)}
          min={2000}
          max={40000}
          step={500}
          testId="lever-workingday"
          onChange={(v) => {
            setSurplus(v);
            push({ surplus_rate_base_bp: v });
          }}
        />
        <Lever
          label="Wage · value of labour-power"
          sub="press it toward subsistence"
          value={wage}
          display={formatPence(wage)}
          min={500}
          max={6000}
          step={100}
          testId="lever-wage"
          onChange={(v) => {
            setWage(v);
            push({ base_wage_pence: v });
          }}
        />
        <Lever
          label="Accumulation rate · α"
          sub="reinvest surplus as machinery"
          value={accum}
          display={formatBP(accum)}
          min={0}
          max={10000}
          step={250}
          testId="lever-accumulation"
          onChange={(v) => {
            setAccum(v);
            push({ accumulation_rate_bp: v });
          }}
        />
      </div>
    </section>
  );
}
