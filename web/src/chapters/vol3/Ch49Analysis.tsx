import { useState } from "react";
import { financeApi } from "../../api";
import type {
  AnnualCompositionResponse,
  SurplusValuePartition,
} from "../../types";
import "./Ch49Analysis.css";

export function Ch49Analysis() {
  // Annual-composition builder (defaults = Ch. 49 seed fixture).
  const [constant, setConstant] = useState(6000);
  const [variable, setVariable] = useState(1500);
  const [surplus, setSurplus] = useState(1500);
  const [compResult, setCompResult] = useState<AnnualCompositionResponse | null>(null);
  const [compError, setCompError] = useState<string | null>(null);

  // Surplus-value-partition builder.
  const [compId, setCompId] = useState("5eed000000000000004900");
  const [profitEnt, setProfitEnt] = useState(900);
  const [interest, setInterest] = useState(300);
  const [groundRent, setGroundRent] = useState(300);
  const [residual, setResidual] = useState(0);
  const [partResult, setPartResult] = useState<SurplusValuePartition | null>(null);
  const [partError, setPartError] = useState<string | null>(null);

  const totalValue = constant + variable + surplus;
  const newValue = variable + surplus;
  const partTotal = profitEnt + interest + groundRent + residual;

  async function handleCreateComposition() {
    setCompError(null);
    setCompResult(null);
    try {
      const r = await financeApi.createAnnualComposition({
        constant_capital_labour_minutes: constant,
        variable_capital_labour_minutes: variable,
        surplus_value_labour_minutes: surplus,
      });
      setCompResult(r);
      setCompId(r.composition.id);
    } catch (e) {
      setCompError(String(e));
    }
  }

  async function handleCreatePartition() {
    setPartError(null);
    setPartResult(null);
    try {
      const r = await financeApi.createSurplusPartition({
        composition_id: compId,
        profit_enterprise_lm: profitEnt,
        interest_lm: interest,
        ground_rent_lm: groundRent,
        residual_surplus_lm: residual,
      });
      setPartResult(r);
    } catch (e) {
      setPartError(String(e));
    }
  }

  return (
    <div className="v3-ch49">
      <section className="v3-ch49-intro">
        <p>
          Chapter&nbsp;49 analyses the value of the annual product and demolishes Adam Smith&apos;s
          &ldquo;adding-up&rdquo; doctrine. The value of any commodity resolves into{" "}
          <strong>c&nbsp;+&nbsp;v&nbsp;+&nbsp;s</strong>: constant capital transferred from the means
          of production, variable capital reproduced as wages, and surplus-value newly created.
        </p>
        <p>
          Only <strong>v&nbsp;+&nbsp;s</strong> is newly created value; the constant capital c is
          merely transferred and must be reproduced, not consumed as revenue. Smith&apos;s error is
          to resolve the whole product into wages&nbsp;+&nbsp;profit&nbsp;+&nbsp;rent, omitting c — as
          if constant capital needed no replacement. At the level of the total social capital this is
          impossible: a portion of the annual product must always return to production.
        </p>
        <p>
          The surplus-value s is the <strong>limit</strong> from which profit (of enterprise +
          interest) and rent are drawn; together they cannot exceed it. Distributable revenue is
          therefore capped at the new value v&nbsp;+&nbsp;s.
        </p>
        <p className="v3-ch49-note">
          Magnitudes in labour-minutes. Seed fixture: c=6000, v=1500, s=1500 → total 9000, new value
          3000; partition s=1500 → profit of enterprise 900 + interest 300 + ground-rent 300.
        </p>
      </section>

      <section className="v3-ch49-section" aria-label="Annual value composition">
        <h4>1. Annual value composition</h4>
        <p className="v3-ch49-subintro">
          Enter the c, v, s of the annual product. The server computes the total value (c+v+s) and
          new value (v+s), and derives the revenue limit: maximum distributable revenue = v+s, with c
          to be reproduced.
        </p>
        <div className="v3-ch49-form">
          <label>
            Constant capital c (lm)
            <input
              type="number"
              min={0}
              value={constant}
              onChange={(e) => setConstant(Math.max(0, Number(e.target.value)))}
            />
          </label>
          <label>
            Variable capital v (lm)
            <input
              type="number"
              min={0}
              value={variable}
              onChange={(e) => setVariable(Math.max(0, Number(e.target.value)))}
            />
          </label>
          <label>
            Surplus-value s (lm)
            <input
              type="number"
              min={0}
              value={surplus}
              onChange={(e) => setSurplus(Math.max(0, Number(e.target.value)))}
            />
          </label>
        </div>
        <p className="v3-ch49-totals">
          total value (c+v+s) = <strong>{totalValue}</strong>&nbsp;lm&nbsp;·&nbsp; new value (v+s) ={" "}
          <strong>{newValue}</strong>&nbsp;lm
        </p>
        <button className="v3-ch49-btn" onClick={handleCreateComposition}>
          Compute composition
        </button>
        {compError && <p className="v3-ch49-error">{compError}</p>}
        {compResult && (
          <div className="v3-ch49-preview">
            <strong>
              Saved — total {compResult.composition.total_value_labour_minutes}&nbsp;lm, new value{" "}
              {compResult.composition.new_value_labour_minutes}&nbsp;lm; revenue limit max{" "}
              {compResult.revenue_limit.max_distributable_revenue_lm}&nbsp;lm, constant replacement{" "}
              {compResult.revenue_limit.constant_capital_replacement}&nbsp;lm
            </strong>
            <span className="v3-ch49-hint">Composition ID: {compResult.composition.id}</span>
          </div>
        )}
      </section>

      <section className="v3-ch49-section" aria-label="Surplus-value partition">
        <h4>2. Surplus-value partition</h4>
        <p className="v3-ch49-subintro">
          Divide the surplus-value into profit of enterprise, interest, ground-rent, and any
          residual. The server rejects (422) a partition whose total exceeds the surplus-value the
          referenced composition produced — Smith&apos;s error made impossible.
        </p>
        <div className="v3-ch49-form">
          <label>
            Composition ID
            <input value={compId} onChange={(e) => setCompId(e.target.value)} />
          </label>
          <label>
            Profit of enterprise (lm)
            <input
              type="number"
              min={0}
              value={profitEnt}
              onChange={(e) => setProfitEnt(Math.max(0, Number(e.target.value)))}
            />
          </label>
          <label>
            Interest (lm)
            <input
              type="number"
              min={0}
              value={interest}
              onChange={(e) => setInterest(Math.max(0, Number(e.target.value)))}
            />
          </label>
          <label>
            Ground-rent (lm)
            <input
              type="number"
              min={0}
              value={groundRent}
              onChange={(e) => setGroundRent(Math.max(0, Number(e.target.value)))}
            />
          </label>
          <label>
            Residual surplus (lm)
            <input
              type="number"
              min={0}
              value={residual}
              onChange={(e) => setResidual(Math.max(0, Number(e.target.value)))}
            />
            <span className="v3-ch49-hint">partition total = {partTotal}&nbsp;lm</span>
          </label>
        </div>
        <button className="v3-ch49-btn" onClick={handleCreatePartition}>
          Partition surplus-value
        </button>
        {partError && <p className="v3-ch49-error">{partError}</p>}
        {partResult && (
          <div className="v3-ch49-preview">
            <strong>
              Saved — total surplus-value {partResult.total_surplus_value_lm}&nbsp;lm (profit{" "}
              {partResult.profit_enterprise_lm} + interest {partResult.interest_lm} + rent{" "}
              {partResult.ground_rent_lm} + residual {partResult.residual_surplus_lm})
            </strong>
            <span className="v3-ch49-hint">Partition ID: {partResult.id}</span>
          </div>
        )}
      </section>
    </div>
  );
}
