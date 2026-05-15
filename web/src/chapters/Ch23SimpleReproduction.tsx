import { useState } from "react";
import type { FormEvent } from "react";
import { api } from "../api";
import { usePounds } from "../CurrencyContext";
import type { ReproductionCycleItem, SimpleReproductionResult } from "../types";

function runningTotal(cycles: ReproductionCycleItem[], upTo: number): number {
  return cycles.slice(0, upTo).reduce((sum, c) => sum + c.revenue, 0);
}

export function Ch23SimpleReproduction() {
  const fmt = usePounds();
  const [constantCapital, setConstantCapital] = useState(800);
  const [variableCapital, setVariableCapital] = useState(200);
  const [surplusRate, setSurplusRate] = useState(1.0);
  const [periods, setPeriods] = useState(10);
  const [result, setResult] = useState<SimpleReproductionResult | null>(null);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setError("");
    setResult(null);
    setLoading(true);
    try {
      const r = await api.runSimpleReproduction({
        constant_capital: constantCapital,
        variable_capital: variableCapital,
        surplus_rate: surplusRate,
        periods,
      });
      setResult(r);
    } catch (err) {
      setError(err instanceof Error ? err.message : String(err));
    } finally {
      setLoading(false);
    }
  }

  function applyMarxExample() {
    setConstantCapital(800);
    setVariableCapital(200);
    setSurplusRate(1.0);
    setPeriods(10);
    setResult(null);
    setError("");
  }

  const originalCapital = result
    ? result.constant_capital + result.variable_capital
    : constantCapital + variableCapital;

  return (
    <div className="ch-panel">
      <section className="ch-section">
        <h2 className="ch-section-title">Simple Reproduction</h2>
        <p className="ch-description">
          In simple reproduction the entire surplus-value is consumed by the
          capitalist as revenue each cycle. The capital stock is reproduced
          without growth. After a calculable number of periods the original
          capital is replaced entirely by past unpaid labour.
        </p>

        <div className="preset-row">
          <button type="button" className="preset-btn" onClick={applyMarxExample}>
            § Core example — £800c + £200v, s′=100%
          </button>
        </div>

        <form className="ch-form" onSubmit={handleSubmit}>
          <div className="form-row">
            <label className="form-label">
              Constant capital c (pence)
              <input
                type="number"
                min={0}
                className="form-input"
                value={constantCapital}
                onChange={(e) => setConstantCapital(Number(e.target.value))}
                required
              />
            </label>
            <label className="form-label">
              Variable capital v (pence)
              <input
                type="number"
                min={1}
                className="form-input"
                value={variableCapital}
                onChange={(e) => setVariableCapital(Number(e.target.value))}
                required
              />
            </label>
          </div>
          <div className="form-row">
            <label className="form-label">
              Rate of surplus-value s′
              <input
                type="number"
                min={0}
                step={0.01}
                className="form-input"
                value={surplusRate}
                onChange={(e) => setSurplusRate(Number(e.target.value))}
                required
              />
            </label>
            <label className="form-label">
              Periods: {periods}
              <input
                type="range"
                min={1}
                max={20}
                className="form-input"
                value={periods}
                onChange={(e) => setPeriods(Number(e.target.value))}
              />
            </label>
          </div>
          <button type="submit" className="ch-btn" disabled={loading}>
            {loading ? "Running..." : "Run simulation"}
          </button>
        </form>

        {error && <p className="ch-error">{error}</p>}
      </section>

      {result && (
        <section className="ch-section">
          <h2 className="ch-section-title">Results</h2>

          <div className="rates-grid">
            <div className="rates-card">
              <div className="rates-label">Original capital</div>
              <div className="rates-value">
                {fmt(result.constant_capital + result.variable_capital)}
              </div>
              <div className="rates-note">
                c = {fmt(result.constant_capital)}&ensp;|&ensp;v ={" "}
                {fmt(result.variable_capital)}
              </div>
            </div>
            <div className="rates-card">
              <div className="rates-label">Annual surplus-value</div>
              <div className="rates-value">
                {fmt(result.cycles[0]?.surplus_total ?? 0)}
              </div>
              <div className="rates-note">s′ = {(result.surplus_rate * 100).toFixed(0)}%</div>
            </div>
            <div className="rates-card">
              <div className="rates-label">Repayment period</div>
              <div className="rates-value">{result.repayment_period} periods</div>
              <div className="rates-note">
                After this, original capital = consumed surplus-value
              </div>
            </div>
          </div>

          <table className="ch-table" style={{ marginTop: "1.5rem" }}>
            <thead>
              <tr>
                <th>Period</th>
                <th>c</th>
                <th>v</th>
                <th>s</th>
                <th>Revenue consumed</th>
                <th>Running total consumed</th>
                <th>vs. original capital</th>
              </tr>
            </thead>
            <tbody>
              {result.cycles.map((c) => {
                const total = runningTotal(result.cycles, c.period);
                const repaid = total >= originalCapital;
                return (
                  <tr
                    key={c.period}
                    style={repaid ? { background: "var(--color-success-bg, #e6f4ea)" } : undefined}
                  >
                    <td>{c.period}</td>
                    <td>{fmt(c.constant_capital)}</td>
                    <td>{fmt(c.variable_capital)}</td>
                    <td>{fmt(c.surplus_total)}</td>
                    <td>{fmt(c.revenue)}</td>
                    <td>{fmt(total)}</td>
                    <td>
                      {repaid
                        ? `≥ ${fmt(originalCapital)} ✓`
                        : `${fmt(originalCapital - total)} remaining`}
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>

          <p className="ch-description" style={{ marginTop: "1rem" }}>
            Variable capital in each cycle is drawn from the surplus-value of
            prior cycles — not from the original owner&apos;s fund. After period{" "}
            {result.repayment_period} the capitalist&apos;s claim to the
            original advance is dissolved entirely by past unpaid labour.
          </p>
        </section>
      )}
    </div>
  );
}
