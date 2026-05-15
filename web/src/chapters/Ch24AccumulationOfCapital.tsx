import { useState } from "react";
import type { FormEvent } from "react";
import { api } from "../api";
import { usePounds } from "../CurrencyContext";
import type { AccumulationPeriod, ExtendedReproductionResult, SplitSurplusResult } from "../types";

export function Ch24AccumulationOfCapital() {
  const fmt = usePounds();

  // Extended reproduction inputs
  const [constantCapital, setConstantCapital] = useState(8000);
  const [variableCapital, setVariableCapital] = useState(2000);
  const [surplusRate, setSurplusRate] = useState(1.0);
  const [accumRate, setAccumRate] = useState(1.0);
  const [compositionRatio, setCompositionRatio] = useState(0.8);
  const [periods, setPeriods] = useState(3);
  const [result, setResult] = useState<ExtendedReproductionResult | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  // Split-surplus inputs
  const [splitSurplus, setSplitSurplus] = useState(2000);
  const [splitAccumRate, setSplitAccumRate] = useState(1.0);
  const [splitCompositionRatio, setSplitCompositionRatio] = useState(0.8);
  const [splitResult, setSplitResult] = useState<SplitSurplusResult | null>(null);
  const [splitLoading, setSplitLoading] = useState(false);
  const [splitError, setSplitError] = useState("");

  function applySpinnerExample() {
    setConstantCapital(8000);
    setVariableCapital(2000);
    setSurplusRate(1.0);
    setAccumRate(1.0);
    setCompositionRatio(0.8);
    setPeriods(3);
    setResult(null);
    setError("");
  }

  function applyPartialAccumulationExample() {
    setConstantCapital(8000);
    setVariableCapital(2000);
    setSurplusRate(1.0);
    setAccumRate(0.5);
    setCompositionRatio(0.8);
    setPeriods(3);
    setResult(null);
    setError("");
  }

  async function handleExtendedSubmit(e: FormEvent) {
    e.preventDefault();
    setError("");
    setResult(null);
    setLoading(true);
    try {
      const r = await api.runExtendedReproduction({
        constant_capital: constantCapital,
        variable_capital: variableCapital,
        surplus_rate: surplusRate,
        accum_rate: accumRate,
        composition_ratio: compositionRatio,
        periods,
      });
      setResult(r);
    } catch (err) {
      setError(err instanceof Error ? err.message : String(err));
    } finally {
      setLoading(false);
    }
  }

  async function handleSplitSubmit(e: FormEvent) {
    e.preventDefault();
    setSplitError("");
    setSplitResult(null);
    setSplitLoading(true);
    try {
      const r = await api.splitSurplus({
        surplus: splitSurplus,
        accum_rate: splitAccumRate,
        composition_ratio: splitCompositionRatio,
      });
      setSplitResult(r);
    } catch (err) {
      setSplitError(err instanceof Error ? err.message : String(err));
    } finally {
      setSplitLoading(false);
    }
  }

  const totalNewCapital = (c: AccumulationPeriod) => c.new_constant + c.new_variable;

  return (
    <div className="ch-panel">
      {/* Extended Reproduction Section */}
      <section className="ch-section">
        <h2 className="ch-section-title">Extended Reproduction — Genealogy of Accumulated Capital</h2>
        <p className="ch-description">
          In extended reproduction a fraction of each period&apos;s surplus-value is
          converted back into capital — new constant capital (machinery, raw materials)
          and new variable capital (wages). Each layer of accumulated capital in turn
          begets further surplus-value: the &ldquo;Abraham begat Isaac&rdquo; sequence.
        </p>

        <div className="preset-row">
          <button type="button" className="preset-btn" onClick={applySpinnerExample}>
            § Spinner §1 — £8,000c + £2,000v, s′=100%, full accumulation
          </button>
          <button type="button" className="preset-btn" onClick={applyPartialAccumulationExample}>
            § Partial §3 — accumulation rate 50%
          </button>
        </div>

        <form className="ch-form" onSubmit={handleExtendedSubmit}>
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
              Accumulation rate (0–1)
              <input
                type="number"
                min={0}
                max={1}
                step={0.01}
                className="form-input"
                value={accumRate}
                onChange={(e) => setAccumRate(Number(e.target.value))}
                required
              />
            </label>
          </div>
          <div className="form-row">
            <label className="form-label">
              Composition ratio (c-share of new capital)
              <input
                type="number"
                min={0}
                max={1}
                step={0.01}
                className="form-input"
                value={compositionRatio}
                onChange={(e) => setCompositionRatio(Number(e.target.value))}
                required
              />
            </label>
            <label className="form-label">
              Periods: {periods}
              <input
                type="range"
                min={1}
                max={10}
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
          <h2 className="ch-section-title">Results — Accumulated Capital Genealogy</h2>
          <p className="ch-description">
            Each row shows the capital layer entering that period, the surplus it
            produces, and how that surplus is divided. The next period&apos;s capital is
            the &ldquo;new capital&rdquo; column — the portion of surplus converted back
            into means of production and labour-power.
          </p>

          <table className="ch-table">
            <thead>
              <tr>
                <th>Period</th>
                <th>c entering</th>
                <th>v entering</th>
                <th>s produced</th>
                <th>New c</th>
                <th>New v</th>
                <th>New capital</th>
                <th>Revenue consumed</th>
              </tr>
            </thead>
            <tbody>
              {result.cycles.map((c) => (
                <tr key={c.period}>
                  <td>{c.period}</td>
                  <td>{fmt(c.constant_capital)}</td>
                  <td>{fmt(c.variable_capital)}</td>
                  <td>{fmt(c.surplus_produced)}</td>
                  <td>{fmt(c.new_constant)}</td>
                  <td>{fmt(c.new_variable)}</td>
                  <td>{fmt(totalNewCapital(c))}</td>
                  <td>{fmt(c.revenue)}</td>
                </tr>
              ))}
            </tbody>
          </table>

          {result.accum_rate === 1.0 && (
            <p className="ch-description" style={{ marginTop: "1rem" }}>
              Under full accumulation the entire surplus-value is converted to capital
              — revenue = 0 each period. Capital grows by compounding the accumulated
              fraction, not the total stock.
            </p>
          )}
        </section>
      )}

      {/* Split Surplus Section */}
      <section className="ch-section">
        <h2 className="ch-section-title">Split Surplus — Stateless Calculator</h2>
        <p className="ch-description">
          Given a surplus-value amount, an accumulation rate, and the organic
          composition ratio, compute the division into new constant capital, new
          variable capital, and revenue consumed by the capitalist.
        </p>

        <form className="ch-form" onSubmit={handleSplitSubmit}>
          <div className="form-row">
            <label className="form-label">
              Surplus-value (pence)
              <input
                type="number"
                min={0}
                className="form-input"
                value={splitSurplus}
                onChange={(e) => setSplitSurplus(Number(e.target.value))}
                required
              />
            </label>
            <label className="form-label">
              Accumulation rate (0–1)
              <input
                type="number"
                min={0}
                max={1}
                step={0.01}
                className="form-input"
                value={splitAccumRate}
                onChange={(e) => setSplitAccumRate(Number(e.target.value))}
                required
              />
            </label>
            <label className="form-label">
              Composition ratio (c-share)
              <input
                type="number"
                min={0}
                max={1}
                step={0.01}
                className="form-input"
                value={splitCompositionRatio}
                onChange={(e) => setSplitCompositionRatio(Number(e.target.value))}
                required
              />
            </label>
          </div>
          <button type="submit" className="ch-btn" disabled={splitLoading}>
            {splitLoading ? "Computing..." : "Split surplus"}
          </button>
        </form>

        {splitError && <p className="ch-error">{splitError}</p>}

        {splitResult && (
          <div className="rates-grid" style={{ marginTop: "1rem" }}>
            <div className="rates-card">
              <div className="rates-label">New constant capital</div>
              <div className="rates-value">{fmt(splitResult.new_constant)}</div>
              <div className="rates-note">
                {(splitResult.composition_ratio * 100).toFixed(0)}% of accumulated portion
              </div>
            </div>
            <div className="rates-card">
              <div className="rates-label">New variable capital</div>
              <div className="rates-value">{fmt(splitResult.new_variable)}</div>
              <div className="rates-note">
                {((1 - splitResult.composition_ratio) * 100).toFixed(0)}% of accumulated portion
              </div>
            </div>
            <div className="rates-card">
              <div className="rates-label">Revenue consumed</div>
              <div className="rates-value">{fmt(splitResult.revenue)}</div>
              <div className="rates-note">
                {((1 - splitResult.accum_rate) * 100).toFixed(0)}% of total surplus — capitalist&apos;s fund of consumption
              </div>
            </div>
          </div>
        )}
      </section>
    </div>
  );
}
