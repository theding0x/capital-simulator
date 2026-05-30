import { useEffect, useState } from "react";
import { financeApi } from "../../api";
import type { WageEffectAnalysisResponse, SphereWageOutcomeResponse } from "../../types";
import "./Ch11WageFluctuations.css";

function pct(bp: number): string {
  return `${(bp / 100).toLocaleString("en-GB", { maximumFractionDigits: 2 })}%`;
}

function price(centi: number): string {
  return `£${(centi / 100).toLocaleString("en-GB", { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`;
}

function dirTag(dir: string): string {
  if (dir === "rises") return "rises";
  if (dir === "falls") return "falls";
  return "unchanged";
}

const DEFAULT_SPHERES = [
  { name: "average", c: 80, v: 20 },
  { name: "lower",   c: 50, v: 50 },
  { name: "higher",  c: 92, v:  8 },
];

const DEFAULT_BASE = { base_constant: 80, base_variable: 20, s_rate: 10000 };

export function Ch11WageFluctuations() {
  const [result, setResult] = useState<WageEffectAnalysisResponse | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [wageFactor, setWageFactor] = useState(125);
  const [base, setBase] = useState(DEFAULT_BASE);

  // On mount, POST pre-seeded Marx fixtures to populate the panel.
  useEffect(() => {
    financeApi
      .createWageEffectAnalysis({
        ...DEFAULT_BASE,
        wage_factor: 125,
        spheres: DEFAULT_SPHERES,
      })
      .then(setResult)
      .catch(() => {});
  }, []);

  async function compute() {
    setError(null);
    try {
      const r = await financeApi.createWageEffectAnalysis({
        ...base,
        wage_factor: wageFactor,
        spheres: DEFAULT_SPHERES,
      });
      setResult(r);
    } catch (e) {
      setError(String(e));
    }
  }

  return (
    <div className="v3-ch11">
      <section className="v3-ch11-intro">
        <p>
          A general rise in wages does <strong>not</strong> raise all prices.
          It lowers the <strong>general rate of profit</strong> and shifts
          prices of production in opposite directions by organic composition:
          higher-composition (capital-intensive) spheres see prices{" "}
          <strong>fall</strong>; lower-composition (labour-intensive) spheres
          see prices <strong>rise</strong>; the average-composition sphere is{" "}
          <strong>unchanged</strong>. Total price is conserved (Σ changes = 0).
        </p>
      </section>

      {error && <p className="v3-ch11-error">{error}</p>}

      <section className="v3-ch11-section" aria-label="Wage-change controls">
        <h4>Wage change</h4>
        <div className="v3-ch11-form">
          <label>
            Base c (average composition)
            <input
              type="number" min={0} value={base.base_constant}
              onChange={(e) => setBase({ ...base, base_constant: Number(e.target.value) })}
            />
          </label>
          <label>
            Base v (average composition)
            <input
              type="number" min={0} value={base.base_variable}
              onChange={(e) => setBase({ ...base, base_variable: Number(e.target.value) })}
            />
          </label>
          <label>
            s′ (basis points, 10000 = 100%)
            <input
              type="number" min={0} value={base.s_rate}
              onChange={(e) => setBase({ ...base, s_rate: Number(e.target.value) })}
            />
          </label>
          <label className="v3-ch11-slider-label">
            Wage factor: {wageFactor}% ({wageFactor > 100 ? "rise" : wageFactor < 100 ? "fall" : "no change"})
            <input
              type="range" min={50} max={150} step={1} value={wageFactor}
              onChange={(e) => setWageFactor(Number(e.target.value))}
            />
          </label>
          <button type="button" className="v3-ch11-btn" onClick={compute}>
            Compute
          </button>
        </div>
      </section>

      {result && (
        <section className="v3-ch11-section" aria-label="General rate comparison">
          <h4>General rate of profit</h4>
          <div className="v3-ch11-rates">
            <div className="v3-ch11-rate-cell">
              <span className="v3-ch11-rate-label">Before</span>
              <span className="v3-ch11-rate-value">{pct(result.old_general_rate)}</span>
            </div>
            <div className={`v3-ch11-rate-arrow ${result.new_general_rate < result.old_general_rate ? "falls" : result.new_general_rate > result.old_general_rate ? "rises" : ""}`}>
              {result.new_general_rate < result.old_general_rate ? "↓" : result.new_general_rate > result.old_general_rate ? "↑" : "–"}
            </div>
            <div className="v3-ch11-rate-cell">
              <span className="v3-ch11-rate-label">After</span>
              <span className="v3-ch11-rate-value">{pct(result.new_general_rate)}</span>
            </div>
            <div className="v3-ch11-rate-kind">
              {result.kind === "rise" ? "Wage rise — general rate falls" : result.kind === "fall" ? "Wage fall — general rate rises" : "No change"}
            </div>
          </div>
        </section>
      )}

      {result && (
        <section className="v3-ch11-section" aria-label="Sphere outcomes">
          <h4>Price of production by sphere</h4>
          <table className="v3-ch11-table">
            <thead>
              <tr>
                <th>Sphere</th>
                <th>Composition</th>
                <th>Old price</th>
                <th>New price</th>
                <th>Change</th>
                <th>Direction</th>
              </tr>
            </thead>
            <tbody>
              {[result.average_outcome, ...result.outcomes.filter((o) => o.sphere_name !== "average")].map(
                (o: SphereWageOutcomeResponse) => (
                  <tr key={o.sphere_name} className={`v3-ch11-row-${o.direction}`}>
                    <td>{o.sphere_name}</td>
                    <td>{o.composition}</td>
                    <td>{price(o.old_price)}</td>
                    <td>{price(o.new_price)}</td>
                    <td className={o.price_change > 0 ? "v3-ch11-positive" : o.price_change < 0 ? "v3-ch11-negative" : ""}>
                      {o.price_change > 0 ? "+" : ""}{price(o.price_change)}
                    </td>
                    <td>
                      <span className={`v3-ch11-tag v3-ch11-tag-${dirTag(o.direction)}`}>
                        {o.direction}
                      </span>
                    </td>
                  </tr>
                )
              )}
            </tbody>
          </table>
        </section>
      )}
    </div>
  );
}
