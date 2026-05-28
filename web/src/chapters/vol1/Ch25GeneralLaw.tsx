import { useState } from "react";
import type { FormEvent } from "react";
import { api } from "../../api";
import { usePounds } from "../../CurrencyContext";
import type { GeneralLawScenarioResult, OrganicCompositionResult, LabourDemandResult, ReserveArmyResult } from "../../types";
import { CapitalCompositionForm } from "../CapitalCompositionForm";
import type { CapitalCompositionValues } from "../CapitalCompositionForm";

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
  const [scenCv, setScenCv] = useState<CapitalCompositionValues>({ constant: 8000, variable: 2000, surplusRate: 1.0, periods: 5 });
  const [scenAccumRate, setScenAccumRate] = useState(1.0);
  const [scenProductivity, setScenProductivity] = useState(0.0);
  const [scenWage, setScenWage] = useState(200);
  const [scenSupply, setScenSupply] = useState(50);
  const [scenResult, setScenResult] = useState<GeneralLawScenarioResult | null>(null);
  const [scenLoading, setScenLoading] = useState(false);
  const [scenError, setScenError] = useState("");

  function applySection1() {
    setScenName("§1 Unchanged composition — demand rises with capital");
    setScenCv({ constant: 8000, variable: 2000, surplusRate: 1.0, periods: 5 });
    setScenAccumRate(1.0);
    setScenProductivity(0.0);
    setScenWage(200);
    setScenSupply(50);
    setScenResult(null);
    setScenError("");
  }

  function applySection2() {
    setScenName("§2 Rising OC — reserve army swells as machinery displaces labour");
    setScenCv({ constant: 8000, variable: 2000, surplusRate: 1.0, periods: 10 });
    setScenAccumRate(1.0);
    setScenProductivity(0.05);
    setScenWage(200);
    setScenSupply(50);
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
        constant_capital: scenCv.constant,
        variable_capital: scenCv.variable,
        surplus_rate: scenCv.surplusRate,
        accumulation_rate: scenAccumRate,
        productivity_growth: scenProductivity,
        wage_pence: scenWage,
        worker_supply: scenSupply,
        periods: scenCv.periods,
      });
      setScenResult(r);
    } catch (err) {
      setScenError(err instanceof Error ? err.message : String(err));
    } finally {
      setScenLoading(false);
    }
  }

  return (
    <div className="ch-panel">
      <p className="ch-description">
        As capital accumulates and the organic composition rises (more machinery
        relative to labour), the demand for workers grows more slowly than total
        capital — and may fall absolutely. The industrial reserve army expands,
        depressing wages and enabling renewed accumulation.
      </p>

      {/* Organic Composition */}
      <section className="ch-section">
        <h2 className="ch-section-title">Organic Composition of Capital (§2)</h2>
        <p className="ch-description">
          c/(c+v) — the ratio of constant to total capital.
        </p>
        <form onSubmit={handleOcSubmit} className="ch-form">
          <div className="form-row">
            <label className="form-label">
              Constant Capital (£)
              <input
                type="number"
                className="form-input"
                value={ocConstant}
                onChange={e => setOcConstant(Number(e.target.value))}
                min={0}
              />
            </label>
            <label className="form-label">
              Variable Capital (£)
              <input
                type="number"
                className="form-input"
                value={ocVariable}
                onChange={e => setOcVariable(Number(e.target.value))}
                min={1}
              />
            </label>
          </div>
          <button type="submit" className="ch-btn" disabled={ocLoading}>
            {ocLoading ? "Computing…" : "Compute"}
          </button>
        </form>
        {ocError && <p className="ch-error">{ocError}</p>}
        {ocResult && (
          <table className="ch-table" style={{ marginTop: "1.5rem" }}>
            <thead>
              <tr>
                <th>Metric</th>
                <th>Value</th>
              </tr>
            </thead>
            <tbody>
              <tr><td>Constant Capital</td><td>{fmt(ocResult.constant_capital)}</td></tr>
              <tr><td>Variable Capital</td><td>{fmt(ocResult.variable_capital)}</td></tr>
              <tr><td>Organic Composition Ratio</td><td>{ocResult.ratio.toFixed(4)}</td></tr>
            </tbody>
          </table>
        )}
      </section>

      {/* Labour Demand */}
      <section className="ch-section">
        <h2 className="ch-section-title">Labour Demand (§1)</h2>
        <p className="ch-description">
          Workers absorbed = total_capital &times; (1 &minus; OC) &divide; wage
        </p>
        <form onSubmit={handleLdSubmit} className="ch-form">
          <div className="form-row">
            <label className="form-label">
              Total Capital (£)
              <input
                type="number"
                className="form-input"
                value={ldTotalCapital}
                onChange={e => setLdTotalCapital(Number(e.target.value))}
                min={1}
              />
            </label>
            <label className="form-label">
              Organic Composition Ratio
              <input
                type="number"
                className="form-input"
                step="0.01"
                value={ldOcRatio}
                onChange={e => setLdOcRatio(Number(e.target.value))}
                min={0}
                max={0.99}
              />
            </label>
            <label className="form-label">
              Wage (pence)
              <input
                type="number"
                className="form-input"
                value={ldWage}
                onChange={e => setLdWage(Number(e.target.value))}
                min={1}
              />
            </label>
          </div>
          <button type="submit" className="ch-btn" disabled={ldLoading}>
            {ldLoading ? "Computing…" : "Compute"}
          </button>
        </form>
        {ldError && <p className="ch-error">{ldError}</p>}
        {ldResult && (
          <table className="ch-table" style={{ marginTop: "1.5rem" }}>
            <thead>
              <tr>
                <th>Metric</th>
                <th>Value</th>
              </tr>
            </thead>
            <tbody>
              <tr><td>Workers Demanded</td><td>{ldResult.workers.toLocaleString()}</td></tr>
            </tbody>
          </table>
        )}
      </section>

      {/* Reserve Army */}
      <section className="ch-section">
        <h2 className="ch-section-title">Industrial Reserve Army (§3)</h2>
        <p className="ch-description">
          Relative surplus population = supply &minus; demanded.
        </p>
        <form onSubmit={handleRaSubmit} className="ch-form">
          <div className="form-row">
            <label className="form-label">
              Worker Supply
              <input
                type="number"
                className="form-input"
                value={raSupply}
                onChange={e => setRaSupply(Number(e.target.value))}
                min={0}
              />
            </label>
            <label className="form-label">
              Workers Demanded
              <input
                type="number"
                className="form-input"
                value={raDemanded}
                onChange={e => setRaDemanded(Number(e.target.value))}
                min={0}
              />
            </label>
          </div>
          <button type="submit" className="ch-btn" disabled={raLoading}>
            {raLoading ? "Computing…" : "Compute"}
          </button>
        </form>
        {raError && <p className="ch-error">{raError}</p>}
        {raResult && (
          <div className="rates-grid" style={{ marginTop: "1rem" }}>
            <div className="rates-card">
              <div className="rates-label">Reserve Army Size</div>
              <div className="rates-value">{raResult.reserve_army_size.toLocaleString()}</div>
            </div>
            <div className="rates-card">
              <div className="rates-label">Relative Proportion</div>
              <div className="rates-value">{(raResult.relative_proportion * 100).toFixed(1)}%</div>
            </div>
          </div>
        )}
      </section>

      {/* Scenario Simulation */}
      <section className="ch-section">
        <h2 className="ch-section-title">Multi-Period Simulation</h2>
        <p className="ch-description">
          Simulate accumulation over time with fixed or rising organic composition.
        </p>
        <div className="preset-row">
          <button type="button" className="preset-btn" onClick={applySection1}>§1 Unchanged OC</button>
          <button type="button" className="preset-btn" onClick={applySection2}>§2 Rising OC</button>
        </div>
        <form onSubmit={handleScenSubmit} className="ch-form">
          <div className="form-row">
            <label className="form-label">
              Scenario Name
              <input
                type="text"
                className="form-input"
                value={scenName}
                onChange={e => setScenName(e.target.value)}
              />
            </label>
          </div>
          <CapitalCompositionForm
            values={scenCv}
            presets={[]}
            onChange={(next) => { setScenCv(next); setScenResult(null); setScenError(""); }}
          />
          <div className="form-row">
            <label className="form-label">
              Accumulation Rate
              <input
                type="number"
                className="form-input"
                step="0.01"
                value={scenAccumRate}
                onChange={e => setScenAccumRate(Number(e.target.value))}
                min={0}
                max={1}
              />
            </label>
            <label className="form-label">
              Productivity Growth (per period)
              <input
                type="number"
                className="form-input"
                step="0.01"
                value={scenProductivity}
                onChange={e => setScenProductivity(Number(e.target.value))}
                min={0}
              />
            </label>
            <label className="form-label">
              Wage (pence)
              <input
                type="number"
                className="form-input"
                value={scenWage}
                onChange={e => setScenWage(Number(e.target.value))}
                min={1}
              />
            </label>
            <label className="form-label">
              Worker Supply
              <input
                type="number"
                className="form-input"
                value={scenSupply}
                onChange={e => setScenSupply(Number(e.target.value))}
                min={1}
              />
            </label>
          </div>
          <button type="submit" className="ch-btn" disabled={scenLoading}>
            {scenLoading ? "Running…" : "Run Simulation"}
          </button>
        </form>
        {scenError && <p className="ch-error">{scenError}</p>}
        {scenResult && (
          <div>
            <p className="ch-description" style={{ marginTop: "1rem" }}>
              <strong style={{ color: "var(--ink)", fontStyle: "normal" }}>Scenario:</strong>{" "}
              {scenResult.name}
            </p>
            <table className="ch-table">
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
