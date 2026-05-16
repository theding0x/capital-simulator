import { useEffect, useState } from "react";
import type { FormEvent } from "react";
import { api } from "../api";
import type { HistoricalStage, LabourDisciplineRegime } from "../types";

const wagePresets = [
  {
    period: "1349",
    jurisdiction: "England",
    max_wage_pence: 24,
    min_wage_pence: 0,
    enforcement_penalty: "imprisonment",
  },
  {
    period: "1740s",
    jurisdiction: "London tailors",
    max_wage_pence: 31,
    min_wage_pence: 0,
    enforcement_penalty: "fine and imprisonment",
  },
];

const vagrancyPresets = [
  {
    period: "1530",
    jurisdiction: "England",
    punishment: "whipping+imprisonment",
    target_population: "sturdy vagabonds",
  },
  {
    period: "1597",
    jurisdiction: "England",
    punishment: "whipping, branding, execution for repeat offenders",
    target_population: "rogues, vagabonds, and sturdy beggars",
  },
];

export function Ch28BloodyLegislation() {
  const [stages, setStages] = useState<HistoricalStage[]>([]);
  const [selectedStageId, setSelectedStageId] = useState("");
  const [regime, setRegime] = useState<LabourDisciplineRegime | null>(null);
  const [regimeError, setRegimeError] = useState("");
  const [loadingRegime, setLoadingRegime] = useState(false);

  const [wsPeriod, setWsPeriod] = useState("");
  const [wsJurisdiction, setWsJurisdiction] = useState("");
  const [wsMax, setWsMax] = useState("0");
  const [wsMin, setWsMin] = useState("0");
  const [wsPenalty, setWsPenalty] = useState("");
  const [wsError, setWsError] = useState("");
  const [wsCreating, setWsCreating] = useState(false);

  const [vlPeriod, setVlPeriod] = useState("");
  const [vlJurisdiction, setVlJurisdiction] = useState("");
  const [vlPunishment, setVlPunishment] = useState("");
  const [vlTarget, setVlTarget] = useState("");
  const [vlError, setVlError] = useState("");
  const [vlCreating, setVlCreating] = useState(false);

  const [acted, setActed] = useState("31");
  const [market, setMarket] = useState("35");
  const [deviation, setDeviation] = useState<number | null>(null);
  const [compareError, setCompareError] = useState("");

  useEffect(() => {
    api.listHistoricalStages().then((list) => {
      setStages(list);
      if (list.length > 0) setSelectedStageId(list[0].id);
    });
  }, []);

  useEffect(() => {
    if (!selectedStageId) return;
    loadRegime(selectedStageId);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [selectedStageId]);

  async function loadRegime(id: string) {
    setRegimeError("");
    setLoadingRegime(true);
    try {
      setRegime(await api.getLabourDiscipline(id));
    } catch (err) {
      setRegimeError(err instanceof Error ? err.message : String(err));
    } finally {
      setLoadingRegime(false);
    }
  }

  function applyWagePreset(p: (typeof wagePresets)[0]) {
    setWsPeriod(p.period);
    setWsJurisdiction(p.jurisdiction);
    setWsMax(String(p.max_wage_pence));
    setWsMin(String(p.min_wage_pence));
    setWsPenalty(p.enforcement_penalty);
  }

  function applyVagrancyPreset(p: (typeof vagrancyPresets)[0]) {
    setVlPeriod(p.period);
    setVlJurisdiction(p.jurisdiction);
    setVlPunishment(p.punishment);
    setVlTarget(p.target_population);
  }

  async function handleCreateWageStatute(e: FormEvent) {
    e.preventDefault();
    setWsError("");
    setWsCreating(true);
    try {
      await api.createWageStatute(selectedStageId, {
        period: wsPeriod,
        jurisdiction: wsJurisdiction,
        max_wage_pence: parseInt(wsMax, 10),
        min_wage_pence: parseInt(wsMin, 10),
        enforcement_penalty: wsPenalty,
      });
      await loadRegime(selectedStageId);
      setWsPeriod(""); setWsJurisdiction(""); setWsMax("0"); setWsMin("0"); setWsPenalty("");
    } catch (err) {
      setWsError(err instanceof Error ? err.message : String(err));
    } finally {
      setWsCreating(false);
    }
  }

  async function handleCreateVagrancyLaw(e: FormEvent) {
    e.preventDefault();
    setVlError("");
    setVlCreating(true);
    try {
      await api.createVagrancyLaw(selectedStageId, {
        period: vlPeriod,
        jurisdiction: vlJurisdiction,
        punishment: vlPunishment,
        target_population: vlTarget,
      });
      await loadRegime(selectedStageId);
      setVlPeriod(""); setVlJurisdiction(""); setVlPunishment(""); setVlTarget("");
    } catch (err) {
      setVlError(err instanceof Error ? err.message : String(err));
    } finally {
      setVlCreating(false);
    }
  }

  async function handleCompare(e: FormEvent) {
    e.preventDefault();
    setCompareError("");
    try {
      const result = await api.compareStatutoryWage(
        parseInt(acted, 10),
        parseInt(market, 10)
      );
      setDeviation(result.deviation);
    } catch (err) {
      setCompareError(err instanceof Error ? err.message : String(err));
    }
  }

  return (
    <div className="chapter-panel">
      <section className="panel-section">
        <h2>Historical Stage</h2>
        <p className="muted small">
          Wage statutes and vagrancy laws are attached to a historical stage.
          Select the stage to view or add coercive legislation.
        </p>
        {stages.length === 0 && (
          <p className="muted small">No historical stages found — create one in Ch. 26 first.</p>
        )}
        {stages.length > 0 && (
          <label>
            Stage
            <select value={selectedStageId} onChange={(e) => setSelectedStageId(e.target.value)}>
              {stages.map((s) => (
                <option key={s.id} value={s.id}>{s.name}</option>
              ))}
            </select>
          </label>
        )}
      </section>

      {selectedStageId && (
        <section className="panel-section">
          <h2>Add Wage Statute</h2>
          <p className="muted small">
            Every statute that fixed a maximum or minimum wage — "the maximum of wages is dictated
            by the State, but on no account a minimum."
          </p>
          <div className="preset-buttons">
            {wagePresets.map((p) => (
              <button
                key={p.period + p.jurisdiction}
                type="button"
                className="btn-secondary btn-sm"
                onClick={() => applyWagePreset(p)}
              >
                {p.period} ({p.jurisdiction})
              </button>
            ))}
          </div>
          <form onSubmit={handleCreateWageStatute} className="form-grid">
            <label>Period<input value={wsPeriod} onChange={(e) => setWsPeriod(e.target.value)} required /></label>
            <label>Jurisdiction<input value={wsJurisdiction} onChange={(e) => setWsJurisdiction(e.target.value)} required /></label>
            <label>Max Wage (halfpennies)<input type="number" min="0" value={wsMax} onChange={(e) => setWsMax(e.target.value)} required /></label>
            <label>Min Wage (halfpennies)<input type="number" min="0" value={wsMin} onChange={(e) => setWsMin(e.target.value)} required /></label>
            <label>Enforcement Penalty<input value={wsPenalty} onChange={(e) => setWsPenalty(e.target.value)} required /></label>
            {wsError && <p className="error-msg">{wsError}</p>}
            <button type="submit" className="btn-primary" disabled={wsCreating}>
              {wsCreating ? "Recording…" : "Record Statute"}
            </button>
          </form>
        </section>
      )}

      {selectedStageId && (
        <section className="panel-section">
          <h2>Add Vagrancy Law</h2>
          <p className="muted small">
            Laws converting the dispossessed into compulsory wage-labourers by
            criminalising the condition of being landless and without an employer.
          </p>
          <div className="preset-buttons">
            {vagrancyPresets.map((p) => (
              <button
                key={p.period}
                type="button"
                className="btn-secondary btn-sm"
                onClick={() => applyVagrancyPreset(p)}
              >
                {p.period}
              </button>
            ))}
          </div>
          <form onSubmit={handleCreateVagrancyLaw} className="form-grid">
            <label>Period<input value={vlPeriod} onChange={(e) => setVlPeriod(e.target.value)} required /></label>
            <label>Jurisdiction<input value={vlJurisdiction} onChange={(e) => setVlJurisdiction(e.target.value)} required /></label>
            <label>Punishment<input value={vlPunishment} onChange={(e) => setVlPunishment(e.target.value)} required /></label>
            <label>Target Population<input value={vlTarget} onChange={(e) => setVlTarget(e.target.value)} required /></label>
            {vlError && <p className="error-msg">{vlError}</p>}
            <button type="submit" className="btn-primary" disabled={vlCreating}>
              {vlCreating ? "Recording…" : "Record Law"}
            </button>
          </form>
        </section>
      )}

      {selectedStageId && (
        <section className="panel-section">
          <h2>Labour Discipline Regime</h2>
          {regimeError && <p className="error-msg">{regimeError}</p>}
          {loadingRegime && <p className="muted">Loading…</p>}
          {regime && (
            <>
              <p className="muted small">Stage: <strong>{regime.period}</strong></p>
              {regime.mechanisms.length > 0 ? (
                <p className="muted small">Enforcement: {regime.mechanisms.join(" · ")}</p>
              ) : (
                <p className="muted small">No coercive statutes — market compulsion operates alone.</p>
              )}
              {regime.wage_statutes.length > 0 && (
                <>
                  <h3 style={{ marginTop: "1rem" }}>Wage Statutes</h3>
                  <table className="data-table">
                    <thead>
                      <tr><th>Period</th><th>Jurisdiction</th><th>Max (½d)</th><th>Min (½d)</th><th>Penalty</th></tr>
                    </thead>
                    <tbody>
                      {regime.wage_statutes.map((ws) => (
                        <tr key={ws.id}>
                          <td>{ws.period}</td><td>{ws.jurisdiction}</td>
                          <td>{ws.max_wage_pence}</td><td>{ws.min_wage_pence}</td>
                          <td>{ws.enforcement_penalty}</td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </>
              )}
              {regime.vagrancy_laws.length > 0 && (
                <>
                  <h3 style={{ marginTop: "1rem" }}>Vagrancy Laws</h3>
                  <table className="data-table">
                    <thead>
                      <tr><th>Period</th><th>Jurisdiction</th><th>Punishment</th><th>Target</th></tr>
                    </thead>
                    <tbody>
                      {regime.vagrancy_laws.map((vl) => (
                        <tr key={vl.id}>
                          <td>{vl.period}</td><td>{vl.jurisdiction}</td>
                          <td>{vl.punishment}</td><td>{vl.target_population}</td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </>
              )}
            </>
          )}
        </section>
      )}

      <section className="panel-section">
        <h2>Statutory vs. Market Wage</h2>
        <p className="muted small">
          Compute the deviation between a legally imposed wage cap and the prevailing
          market wage. Negative = state suppresses wages below market rate.
        </p>
        <form onSubmit={handleCompare} className="form-grid">
          <label>Acted Wage (halfpennies)<input type="number" min="0" value={acted} onChange={(e) => setActed(e.target.value)} required /></label>
          <label>Market Wage (halfpennies)<input type="number" min="1" value={market} onChange={(e) => setMarket(e.target.value)} required /></label>
          {compareError && <p className="error-msg">{compareError}</p>}
          <button type="submit" className="btn-primary">Compare</button>
        </form>
        {deviation !== null && (
          <p className="muted small" style={{ marginTop: "0.5rem" }}>
            Deviation:{" "}
            <strong style={{ color: deviation < 0 ? "var(--color-error, #c0392b)" : "inherit" }}>
              {(deviation * 100).toFixed(2)}%
            </strong>
            {deviation < 0
              ? " — statutory cap suppresses wages below market rate"
              : deviation === 0
              ? " — no suppression"
              : " — statute sets a floor above market rate"}
          </p>
        )}
      </section>
    </div>
  );
}