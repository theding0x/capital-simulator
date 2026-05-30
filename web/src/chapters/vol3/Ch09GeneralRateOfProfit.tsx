import { useEffect, useState } from "react";
import { financeApi } from "../../api";
import type {
  PriceOfProductionResponse,
  SocialCapitalAggregateResponse,
} from "../../types";
import "./Ch09GeneralRateOfProfit.css";

// Rates arrive in basis points (10000 bp = 100%); render as percent.
function pct(bp: number): string {
  return `${(bp / 100).toLocaleString("en-GB", { maximumFractionDigits: 2 })}%`;
}

function money(n: number): string {
  return `£${n.toLocaleString("en-GB")}`;
}

// Client-side mirror of the Go round-half-up formula for the general rate:
//   p̄′ = Math.round(sumS * 10000 / sumC)  [half-up for positive values]
function clientGeneralRate(sumS: number, sumC: number): number {
  return sumC > 0 ? Math.round((sumS * 10000) / sumC) : 0;
}

function clientSurplusValue(sRate: number, v: number): number {
  return Math.floor((sRate * v) / 10000);
}

// The five canonical spheres from Marx's Ch. 9 worked example.
const DEFAULT_SPHERES = [
  { name: "Sphere I (80c+20v)",   c: 100, v: 20, s_rate: 10000 },
  { name: "Sphere II (70c+30v)",  c: 100, v: 30, s_rate: 10000 },
  { name: "Sphere III (60c+40v)", c: 100, v: 40, s_rate: 10000 },
  { name: "Sphere IV (85c+15v)",  c: 100, v: 15, s_rate: 10000 },
  { name: "Sphere V (95c+5v)",    c: 100, v:  5, s_rate: 10000 },
];

type SphereRow = { name: string; c: number; v: number; s_rate: number };

export function Ch09GeneralRateOfProfit() {
  const [rows, setRows] = useState<SphereRow[]>(DEFAULT_SPHERES);
  const [result, setResult] = useState<SocialCapitalAggregateResponse | null>(null);
  const [seeded, setSeeded] = useState<PriceOfProductionResponse[]>([]);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    financeApi.listPricesOfProduction().then(setSeeded).catch(() => {});
  }, []);

  async function compute() {
    setError(null);
    try {
      const agg = await financeApi.computeSocialAggregate({ spheres: rows });
      setResult(agg);
    } catch (e) {
      setError(String(e));
    }
  }

  function updateRow(i: number, field: keyof SphereRow, value: string | number) {
    setRows((prev) => {
      const next = [...prev];
      next[i] = { ...next[i], [field]: typeof value === "string" ? value : Number(value) };
      return next;
    });
  }

  // Derive client-side preview values for the table.
  const sumS = rows.reduce((acc, r) => acc + clientSurplusValue(r.s_rate, r.v), 0);
  const sumC = rows.reduce((acc, r) => acc + r.c, 0);
  const previewRate = clientGeneralRate(sumS, sumC);

  return (
    <div className="v3-ch09">
      <section className="v3-ch09-intro">
        <p>
          Competition of capitals equalises the divergent individual rates of
          profit (Ch. 8) into a single{" "}
          <strong>general (average) rate</strong> p̄′ = ΣS / ΣC — total social
          surplus-value divided by total social capital, capital-weighted. Each
          capital then draws <strong>average profit</strong> = its C × p̄′,
          regardless of how much surplus-value it itself produced. Commodities
          no longer sell at their <em>values</em> but at their{" "}
          <strong>prices of production</strong> = cost-price + average profit.
        </p>
        <p>
          The transfer is zero-sum: higher-composition branches (more c, less
          v → produce <em>less</em> surplus-value) receive more than they
          produced — price &gt; value, deviation positive. Lower-composition
          branches (more v) are net donors: price &lt; value, deviation
          negative.
        </p>
      </section>

      <section className="v3-ch09-builder" aria-label="Five-sphere editor">
        <h4>Sphere inputs — edit name, C, v, s′</h4>
        <table className="v3-ch09-input-table">
          <thead>
            <tr>
              <th>Sphere</th>
              <th>C</th>
              <th>v</th>
              <th>s′ (bp)</th>
              <th>s (preview)</th>
            </tr>
          </thead>
          <tbody>
            {rows.map((row, i) => (
              <tr key={i}>
                <td>
                  <input
                    className="v3-ch09-name-input"
                    type="text"
                    value={row.name}
                    maxLength={60}
                    onChange={(e) => updateRow(i, "name", e.target.value)}
                  />
                </td>
                <td>
                  <input
                    className="v3-ch09-num-input"
                    type="number"
                    min={0}
                    value={row.c}
                    onChange={(e) => updateRow(i, "c", e.target.value)}
                  />
                </td>
                <td>
                  <input
                    className="v3-ch09-num-input"
                    type="number"
                    min={0}
                    max={row.c}
                    value={row.v}
                    onChange={(e) => updateRow(i, "v", e.target.value)}
                  />
                </td>
                <td>
                  <input
                    className="v3-ch09-num-input"
                    type="number"
                    min={0}
                    value={row.s_rate}
                    onChange={(e) => updateRow(i, "s_rate", e.target.value)}
                  />
                </td>
                <td className="v3-ch09-preview">
                  {money(clientSurplusValue(row.s_rate, row.v))}
                </td>
              </tr>
            ))}
          </tbody>
          <tfoot>
            <tr>
              <th>Totals</th>
              <td>{money(sumC)}</td>
              <td></td>
              <td></td>
              <td>{money(sumS)}</td>
            </tr>
          </tfoot>
        </table>
        <div className="v3-ch09-preview-rate">
          Preview p̄′ ≈ <strong>{pct(previewRate)}</strong>
        </div>
        <button type="button" className="v3-ch09-compute-btn" onClick={compute}>
          Compute general rate
        </button>
      </section>

      {error && <p className="v3-ch09-error">{error}</p>}

      {result && (
        <section className="v3-ch09-result" aria-label="Computation result">
          <h4>General rate and prices of production</h4>

          <div className="v3-ch09-rate-banner">
            <div className="v3-ch09-rate-stat">
              <span className="v3-ch09-stat-label">General rate p̄′</span>
              <span className="v3-ch09-stat-value v3-ch09-headline">
                {pct(result.weighted_general_rate)}
              </span>
            </div>
            <div className="v3-ch09-rate-stat">
              <span className="v3-ch09-stat-label">ΣS (total surplus-value)</span>
              <span className="v3-ch09-stat-value">{money(result.sum_surplus_values)}</span>
            </div>
            <div className="v3-ch09-rate-stat">
              <span className="v3-ch09-stat-label">ΣC (total capital)</span>
              <span className="v3-ch09-stat-value">{money(result.sum_total_capitals)}</span>
            </div>
            <div className="v3-ch09-rate-stat">
              <span className="v3-ch09-stat-label">Σ avg profits</span>
              <span className="v3-ch09-stat-value">{money(result.sum_average_profits)}</span>
            </div>
            <div className="v3-ch09-rate-stat">
              <span className="v3-ch09-stat-label">Avg v/C</span>
              <span className="v3-ch09-stat-value">{result.average_variable_percent}%</span>
            </div>
          </div>

          <table className="v3-ch09-pop-table">
            <thead>
              <tr>
                <th>Sphere</th>
                <th>k (cost-price)</th>
                <th>Value</th>
                <th>Avg profit</th>
                <th>Price of production</th>
                <th>Deviation</th>
                <th>Composition</th>
              </tr>
            </thead>
            <tbody>
              {result.prices_of_production.map((pop, i) => (
                <tr key={i} className={`v3-ch09-row-${pop.composition || "unknown"}`}>
                  <td>{pop.sphere_name || result.spheres[i]?.name || `Sphere ${i + 1}`}</td>
                  <td>{money(pop.cost_price)}</td>
                  <td>{money(pop.commodity_value)}</td>
                  <td>{money(pop.average_profit)}</td>
                  <td>{money(pop.price)}</td>
                  <td className={pop.deviation > 0 ? "v3-ch09-positive" : pop.deviation < 0 ? "v3-ch09-negative" : ""}>
                    {pop.deviation > 0 ? "+" : ""}{pop.deviation}
                  </td>
                  <td>
                    <span className={`v3-ch09-tag v3-ch09-tag-${pop.composition || "unknown"}`}>
                      {pop.composition || "—"}
                    </span>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>

          <div className={`v3-ch09-conservation ${result.value_conserved ? "v3-ch09-conserved" : "v3-ch09-not-conserved"}`}>
            {result.value_conserved ? (
              <>
                Conservation holds: Σ prices = Σ values = {money(result.total_values)};
                Σ deviations = 0
              </>
            ) : (
              <>
                Conservation violated: Σ prices = {money(result.total_prices_of_production)},
                Σ values = {money(result.total_values)}
              </>
            )}
          </div>
        </section>
      )}

      {seeded.length > 0 && (
        <section className="v3-ch09-seeded" aria-label="Seeded prices of production">
          <h4>Recorded prices of production</h4>
          <table className="v3-ch09-pop-table">
            <thead>
              <tr>
                <th>Sphere</th>
                <th>k</th>
                <th>p̄′</th>
                <th>Value</th>
                <th>Avg profit</th>
                <th>Price</th>
                <th>Deviation</th>
                <th>Composition</th>
              </tr>
            </thead>
            <tbody>
              {seeded.map((p) => (
                <tr key={p.id} className={`v3-ch09-row-${p.composition || "unknown"}`}>
                  <td>{p.sphere_name}</td>
                  <td>{money(p.cost_price)}</td>
                  <td>{pct(p.general_rate)}</td>
                  <td>{money(p.commodity_value)}</td>
                  <td>{money(p.average_profit)}</td>
                  <td>{money(p.price)}</td>
                  <td className={p.deviation > 0 ? "v3-ch09-positive" : p.deviation < 0 ? "v3-ch09-negative" : ""}>
                    {p.deviation > 0 ? "+" : ""}{p.deviation}
                  </td>
                  <td>
                    <span className={`v3-ch09-tag v3-ch09-tag-${p.composition || "unknown"}`}>
                      {p.composition || "—"}
                    </span>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </section>
      )}
    </div>
  );
}
