import { useState } from "react";
import type { FormEvent } from "react";
import { api } from "../api";
import { usePounds } from "../CurrencyContext";
import type { GeneralLawScenarioResult, OrganicCompositionResult, LabourDemandResult, ReserveArmyResult } from "../types";

export function Ch25GeneralLaw() {
  const fmt = usePounds();

  // Organic composition calculator
  const [ocConstant, setOcConstant] = useState(8000);
  const [ocVariable, setOcVariable] = useState(2000);
  const [ocResult, setOcResult] = useState<OrganicCompositionResult | null>(null);
  const [ocLoading, setOcLoading] = useState(false);
  const [ocError, setOcError] = useState("");

  // Labour demand calculator
  const [ldTotalCapital, setLdTotalCapital] = useState(10000);
  const [ldOcRatio, setLdOcRatio] = useState(0.8);
  const [ldWage, setLdWage] = useState(200);
  const [ldResult, setLdResult] = useState<LabourDemandResult | null>(null);
  const [ldLoading, setLdLoading] = useState(false);
  const [ldError, setLdError] = useState("");

  // Reserve army calculator
  const [raSupply, setRaSupply] = useState(50);
  const [raDemanded, setRaDemanded] = useState(40);
  const [raResult, setRaResult] = useState<ReserveArmyResult | null>(null);
  const [raLoading, setRaLoading] = useState(false);
  const [raError, setRaError] = useState("");

  // Scenario simulation
  const [scenName, setScenName] = useState("§1 Unchanged composition");
  const [scenConstant, setScenConstant] = useState(8000);
  const [scenVariable, setScenVariable] = useState(2000);
  const [scenSurplusRate, setScenSurplusRate] = useState(1.0);
  const [scenAccumRate, setScenAccumRate] = useState(1.0);
  const [scenProductivity, setScenProductivity] = useState(0.0);
  const [scenWage, setScenWage] = useState(200);
  const [scenSupply, setScenSupply] = useState(50);
  const [scenPeriods, setScenPeriods] = useState(5);
  const [scenResult, setScenResult] = useState<GeneralLawScenarioResult | null>(null);
  const [scenLoading, setScenLoading] = useState(false);
  const [scenError, setScenError] = useState("");

  function applySection1() {
    setScenName("§1 Unchanged composition — demand rises with capital");
    setScenConstant(8000);
    setScenVariable(2000);
    setScenSurplusRate(1.0);
    setScenAccumRate(1.0);
    setScenProductivity(0.0);
    setScenWage(200);
    setScenSupply(50);
    setScenPeriods(5);
    setScenResult(null);
    setScenError("");
  }

  function applySection2() {
    setScenName("§2 Rising OC — reserve army swells as machinery displaces labour");
    setScenConstant(8000);
    setScenVariable(2000);
    setScenSurplusRate(1.0);
    setScenAccumRate(1.0);
    setScenProductivity(0.05);
    setScenWage(200);
    setScenSupply(50);
    setScenPeriods(10);
    setScenResult(null);
    setScenError("");
  }

  async function handleOcSubmit(e: FormEvent) {
    e.preventDefault();
    setOcError("");
    setOcResult(null);
    setOcLoading(true);
    try {
      const r = await api.computeOrganicComposition({
        constant_capital: ocConstant,
        variable_capital: ocVariable,
      });
      setOcResult(r);
    } catch (err) {
      setOcError(err instanceof Error ? err.message : String(err));
    } finally {
      setOcLoading(false);
    }
  }

  async function handleLdSubmit(e: FormEvent) {
    e.preventDefault();
    setLdError("");
    setLdResult(null);
    setLdLoading(true);
    try {
      const r = await api.computeLabourDemand({
        total_capital: ldTotalCapital,
        organic_composition_ratio: ldOcRatio,
        wage_pence: ldWage,
      });
      setLdResult(r);
    } catch (err) {
      setLdError(err instanceof Error ? err.message : String(err));
    } finally {
      setLdLoading(false);
    }
  }

  async function handleRaSubmit(e: FormEvent) {
    e.preventDefault();
    setRaError("");
    setRaResult(null);
    setRaLoading(true);
    try {
      const r = await api.computeReserveArmy({
        worker_supply: raSupply,
        workers_demanded: raDemanded,
      });
      setRaResult(r);
    } catch (err) {
      setRaError(err instanceof Error ? err.message : String(err));
    } finally {
      setRaLoading(false);
    }
  }

  async function handleScenSubmit(e: FormEvent) {
    e.preventDefault();
    setScenError("");
    setScenResult(null);
    setScenLoading(true);
    try {
      const r = await api.createGeneralLawScenario({
        name: scenName,
        constant_capital: scenConstant,
        variable_capital: scenVariable,
        surplus_rate: scenSurplusRate,
        accumulation_rate: scenAccumRate,
        productivity_growth: scenProductivity,
        wage_pence: scenWage,
        worker_supply: scenSupply,
        periods: scenPeriods,
      });
      setScenResult(r);
    } catch (err) {
      setScenError(err instanceof Error ? err.message : String(err));
    } finally {
      setScenLoading(false);
    }
  }

  return (
    <div className="chapter-panel">
      <h2>Ch. 25 — The General Law of Capitalist Accumulation</h2>
      <p className="chapter-blurb">
        As capital accumulates and the organic composition rises (more machinery relative to labour),
        the demand for workers grows more slowly than total capital — and may fall absolutely.
        The industrial reserve army expands, depressing wages and enabling renewed accumulation.
      </p>

      {/* Organic Composition */}
      <section className="sim-section">
        <h3>Organic Composition of Capital (§2)</h3>
        <p>c/(c+v) — the ratio of constant to total capital.</p>
        <form onSubmit={handleOcSubmit} className="sim-form">
          <label>
            Constant Capital (£)
            <input type="number" value={ocConstant} onChange={e => setOcConstant(Number(e.target.value))} min={0} />
          </label>
          <label>
            Variable Capital (£)
            <input type="number" value={ocVariable} onChange={e => setOcVariable(Number(e.target.value))} min={1} />
          </label>
          <button type="submit" disabled={ocLoading}>{ocLoading ? "Computing…" : "Compute"}</button>
        </form>
        {ocError && <p className="error">{ocError}</p>}
        {ocResult && (
          <table className="result-table">
            <tbody>
              <tr><th>Constant Capital</th><td>{fmt(ocResult.constant_capital)}</td></tr>
              <tr><th>Variable Capital</th><td>{fmt(ocResult.variable_capital)}</td></tr>
              <tr><th>Organic Composition Ratio</th><td>{ocResult.ratio.toFixed(4)}</td></tr>
            </tbody>
          </table>
        )}
      </section>

      {/* Labour Demand */}
      <section className="sim-section">
        <h3>Labour Demand (§1)</h3>
        <p>Workers absorbed = total_capital × (1 − OC) ÷ wage</p>
        <form onSubmit={handleLdSubmit} className="sim-form">
          <label>
            Total Capital (£)
            <input type="number" value={ldTotalCapital} onChange={e => setLdTotalCapital(Number(e.target.value))} min={1} />
          </label>
          <label>
            Organic Composition Ratio
            <input type="number" step="0.01" value={ldOcRatio} onChange={e => setLdOcRatio(Number(e.target.value))} min={0} max={0.99} />
          </label>
          <label>
            Wage (pence)
            <input type="number" value={ldWage} onChange={e => setLdWage(Number(e.target.value))} min={1} />
          </label>
          <button type="submit" disabled={ldLoading}>{ldLoading ? "Computing…" : "Compute"}</button>
        </form>
        {ldError && <p className="error">{ldError}</p>}
        {ldResult && (
          <table className="result-table">
            <tbody>
              <tr><th>Workers Demanded</th><td>{ldResult.workers.toLocaleString()}</td></tr>
            </tbody>
          </table>
        )}
      </section>

      {/* Reserve Army */}
      <section className="sim-section">
        <h3>Industrial Reserve Army (§3)</h3>
        <p>Relative surplus population = supply − demanded.</p>
        <form onSubmit={handleRaSubmit} className="sim-form">
          <label>
            Worker Supply
            <input type="number" value={raSupply} onChange={e => setRaSupply(Number(e.target.value))} min={0} />
          </label>
          <label>
            Workers Demanded
            <input type="number" value={raDemanded} onChange={e => setRaDemanded(Number(e.target.value))} min={0} />
          </label>
          <button type="submit" disabled={raLoading}>{raLoading ? "Computing…" : "Compute"}</button>
        </form>
        {raError && <p className="error">{raError}</p>}
        {raResult && (
          <table className="result-table">
            <tbody>
              <tr><th>Reserve Army Size</th><td>{raResult.reserve_army_size.toLocaleString()}</td></tr>
              <tr><th>Relative Proportion</th><td>{(raResult.relative_proportion * 100).toFixed(1)}%</td></tr>
            </tbody>
          </table>
        )}
      </section>

      {/* Scenario Simulation */}
      <section className="sim-section">
        <h3>Multi-Period Simulation</h3>
        <p>Simulate accumulation over time with fixed or rising organic composition.</p>
        <div className="preset-row">
          <button type="button" onClick={applySection1}>§1 Unchanged OC</button>
          <button type="button" onClick={applySection2}>§2 Rising OC</button>
        </div>
        <form onSubmit={handleScenSubmit} className="sim-form">
          <label>
            Scenario Name
            <input type="text" value={scenName} onChange={e => setScenName(e.target.value)} />
          </label>
          <label>
            Constant Capital (£)
            <input type="number" value={scenConstant} onChange={e => setScenConstant(Number(e.target.value))} min={0} />
          </label>
          <label>
            Variable Capital (£)
            <input type="number" value={scenVariable} onChange={e => setScenVariable(Number(e.target.value))} min={1} />
          </label>
          <label>
            Surplus Rate (s′)
            <input type="number" step="0.01" value={scenSurplusRate} onChange={e => setScenSurplusRate(Number(e.target.value))} min={0} />
          </label>
          <label>
            Accumulation Rate
            <input type="number" step="0.01" value={scenAccumRate} onChange={e => setScenAccumRate(Number(e.target.value))} min={0} max={1} />
          </label>
          <label>
            Productivity Growth (per period)
            <input type="number" step="0.01" value={scenProductivity} onChange={e => setScenProductivity(Number(e.target.value))} min={0} />
          </label>
          <label>
            Wage (pence)
            <input type="number" value={scenWage} onChange={e => setScenWage(Number(e.target.value))} min={1} />
          </label>
          <label>
            Worker Supply
            <input type="number" value={scenSupply} onChange={e => setScenSupply(Number(e.target.value))} min={1} />
          </label>
          <label>
            Periods
            <input type="number" value={scenPeriods} onChange={e => setScenPeriods(Number(e.target.value))} min={1} />
          </label>
          <button type="submit" disabled={scenLoading}>{scenLoading ? "Running…" : "Run Simulation"}</button>
        </form>
        {scenError && <p className="error">{scenError}</p>}
        {scenResult && (
          <div>
            <p><strong>Scenario:</strong> {scenResult.name}</p>
            <table className="result-table">
              <thead>
                <tr>
                  <th>Period</th>
                  <th>C (£)</th>
                  <th>V (£)</th>
                  <th>OC</th>
                  <th>Workers</th>
                  <th>Reserve Army</th>
                  <th>Relative %</th>
                </tr>
              </thead>
              <tbody>
                {scenResult.series.map(s => (
                  <tr key={s.period}>
                    <td>{s.period}</td>
                    <td>{fmt(s.constant_capital)}</td>
                    <td>{fmt(s.variable_capital)}</td>
                    <td>{s.organic_composition.toFixed(4)}</td>
                    <td>{s.workers.toLocaleString()}</td>
                    <td>{s.reserve_army_size.toLocaleString()}</td>
                    <td>{(s.relative_proportion * 100).toFixed(1)}%</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </section>
    </div>
  );
}
