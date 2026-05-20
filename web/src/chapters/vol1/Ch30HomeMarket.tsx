import { useEffect, useState } from "react";
import type { FormEvent } from "react";
import { api } from "../../api";
import { usePounds } from "../../CurrencyContext";
import type {
  HistoricalStage,
  DomesticIndustry,
  MarketFormation,
  HomeMarketSize,
} from "../../types";

export function Ch30HomeMarket() {
  const fmt = usePounds();

  const [stages, setStages] = useState<HistoricalStage[]>([]);
  const [selectedStageId, setSelectedStageId] = useState("");
  const [homeMarket, setHomeMarket] = useState<HomeMarketSize | null>(null);
  const [homeMarketError, setHomeMarketError] = useState("");
  const [industries, setIndustries] = useState<DomesticIndustry[]>([]);

  // create domestic industry form
  const [diName, setDiName] = useState("Westphalian flax spinning");
  const [diHouseholds, setDiHouseholds] = useState("2000");
  const [diOutput, setDiOutput] = useState("48000");
  const [diError, setDiError] = useState("");
  const [diCreating, setDiCreating] = useState(false);

  // market formation (stateless) form
  const [mfName, setMfName] = useState("Westphalian flax spinning");
  const [mfHouseholds, setMfHouseholds] = useState("2000");
  const [mfOutput, setMfOutput] = useState("48000");
  const [mfExpropriated, setMfExpropriated] = useState("1500");
  const [mfResult, setMfResult] = useState<MarketFormation | null>(null);
  const [mfError, setMfError] = useState("");
  const [mfComputing, setMfComputing] = useState(false);

  // Mirabeau comparison
  const [reunieWorkers, setReunie] = useState("200");
  const [reunieOutput, setReuniedOutput] = useState("12000");
  const [separeeProducers, setSeparee] = useState("200");
  const [separeeOutput, setSepareeOutput] = useState("8000");

  useEffect(() => {
    api.listHistoricalStages().then((list) => {
      setStages(list);
      if (list.length > 0) {
        setSelectedStageId(list[0].id);
      }
    });
  }, []);

  useEffect(() => {
    if (!selectedStageId) return;
    loadHomeMarket(selectedStageId);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [selectedStageId]);

  async function loadHomeMarket(id: string) {
    setHomeMarketError("");
    try {
      const size = await api.getHomeMarket(id);
      setHomeMarket(size);
    } catch (e) {
      setHomeMarketError(String(e));
    }
  }

  async function handleCreateDomesticIndustry(e: FormEvent) {
    e.preventDefault();
    if (!selectedStageId) return;
    setDiError("");
    setDiCreating(true);
    try {
      const di = await api.createDomesticIndustry(selectedStageId, {
        name: diName,
        households_engaged: parseInt(diHouseholds, 10),
        annual_output_pence: parseInt(diOutput, 10),
      });
      setIndustries((prev) => [...prev, di]);
      await loadHomeMarket(selectedStageId);
    } catch (e) {
      setDiError(String(e));
    } finally {
      setDiCreating(false);
    }
  }

  async function handleComputeFormation(e: FormEvent) {
    e.preventDefault();
    setMfError("");
    setMfComputing(true);
    try {
      const result = await api.computeMarketFormation({
        name: mfName,
        households_engaged: parseInt(mfHouseholds, 10),
        annual_output_pence: parseInt(mfOutput, 10),
        expropriated: parseInt(mfExpropriated, 10),
      });
      setMfResult(result);
    } catch (e) {
      setMfError(String(e));
    } finally {
      setMfComputing(false);
    }
  }

  const reunieOut = parseInt(reunieOutput, 10) || 0;
  const separeeOut = parseInt(separeeOutput, 10) || 0;

  return (
    <div className="chapter-panel">
      <section className="panel-section">
        <h2>Home Market Size</h2>
        <p className="muted small">
          Select a historical stage to see the aggregate home market created by
          proletarianisation of its domestic industries.
        </p>
        <label>
          Historical Stage
          <select
            value={selectedStageId}
            onChange={(e) => setSelectedStageId(e.target.value)}
          >
            {stages.map((s) => (
              <option key={s.id} value={s.id}>
                {s.name}
              </option>
            ))}
          </select>
        </label>
        {homeMarketError && <p className="error-msg">{homeMarketError}</p>}
        {homeMarket && (
          <dl className="result-list">
            <dt>Commoditised Output</dt>
            <dd>{fmt(homeMarket.commoditised_output_pence)}</dd>
            <dt>Wage Labourers Created</dt>
            <dd>{homeMarket.wage_labourers.toLocaleString()}</dd>
          </dl>
        )}
      </section>

      <section className="panel-section">
        <h2>Register Domestic Industry</h2>
        <p className="muted small">
          Record a pre-capitalist household industry for the selected stage.
          Each household registered will be counted as a potential wage worker
          once expropriated, and its annual output becomes a commodity.
        </p>
        <form onSubmit={handleCreateDomesticIndustry} className="form-grid">
          <label>
            Name
            <input
              value={diName}
              onChange={(e) => setDiName(e.target.value)}
              required
            />
          </label>
          <label>
            Households Engaged
            <input
              type="number"
              min="1"
              value={diHouseholds}
              onChange={(e) => setDiHouseholds(e.target.value)}
              required
            />
          </label>
          <label>
            Annual Output (halfpennies)
            <input
              type="number"
              min="1"
              value={diOutput}
              onChange={(e) => setDiOutput(e.target.value)}
              required
            />
          </label>
          {diError && <p className="error-msg">{diError}</p>}
          <button type="submit" className="btn-primary" disabled={diCreating || !selectedStageId}>
            {diCreating ? "Registering…" : "Register Industry"}
          </button>
        </form>
        {industries.length > 0 && (
          <ul className="result-list" style={{ marginTop: "0.75rem" }}>
            {industries.map((d) => (
              <li key={d.id}>
                <strong>{d.name}</strong> — {d.households_engaged.toLocaleString()} households,{" "}
                {fmt(d.annual_output_pence)} annual output
              </li>
            ))}
          </ul>
        )}
      </section>

      <section className="panel-section">
        <h2>Home Market Formation (Stateless)</h2>
        <p className="muted small">
          Compute how the expropriation of a domestic industry creates new wage
          workers and commoditises their output. The proletarianised are the new
          wage workers: every displaced household must now buy what it once made.
        </p>
        <form onSubmit={handleComputeFormation} className="form-grid">
          <label>
            Industry Name
            <input
              value={mfName}
              onChange={(e) => setMfName(e.target.value)}
              required
            />
          </label>
          <label>
            Households Engaged
            <input
              type="number"
              min="1"
              value={mfHouseholds}
              onChange={(e) => setMfHouseholds(e.target.value)}
              required
            />
          </label>
          <label>
            Annual Output (halfpennies)
            <input
              type="number"
              min="1"
              value={mfOutput}
              onChange={(e) => setMfOutput(e.target.value)}
              required
            />
          </label>
          <label>
            Households Expropriated
            <input
              type="number"
              min="0"
              value={mfExpropriated}
              onChange={(e) => setMfExpropriated(e.target.value)}
              required
            />
          </label>
          {mfError && <p className="error-msg">{mfError}</p>}
          <button type="submit" className="btn-primary" disabled={mfComputing}>
            {mfComputing ? "Computing…" : "Compute Formation"}
          </button>
        </form>
        {mfResult && (
          <dl className="result-list">
            <dt>Period</dt>
            <dd>{mfResult.period}</dd>
            <dt>Labourers Proletarianised</dt>
            <dd>{mfResult.labourers_proletarianised.toLocaleString()}</dd>
            <dt>New Wage Labourers</dt>
            <dd>{mfResult.new_wage_labourers.toLocaleString()}</dd>
            <dt>Market Size</dt>
            <dd>{fmt(mfResult.market_size_pence)}</dd>
            <dt>Domestic Industry Destroyed</dt>
            <dd
              style={{
                color: mfResult.domestic_industry_destroyed
                  ? "var(--color-error, #c0392b)"
                  : "inherit",
              }}
            >
              {mfResult.domestic_industry_destroyed ? "Yes" : "No"}
            </dd>
          </dl>
        )}
      </section>

      <section className="panel-section">
        <h2>Manufacture Réunie vs. Manufacture Séparée</h2>
        <p className="muted small">
          Mirabeau's distinction: the "grand manufactory" (réunie) concentrates
          workers under one direction and outproduces the scattered discrete
          workshop (séparée) at equal headcount. His preference for séparée is
          political, not economic.
        </p>
        <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: "1rem" }}>
          <div>
            <h3 style={{ marginBottom: "0.5rem" }}>Réunie (concentrated)</h3>
            <label>
              Workers
              <input
                type="number"
                min="1"
                value={reunieWorkers}
                onChange={(e) => setReunie(e.target.value)}
              />
            </label>
            <label style={{ marginTop: "0.5rem" }}>
              Output (halfpennies)
              <input
                type="number"
                min="0"
                value={reunieOutput}
                onChange={(e) => setReuniedOutput(e.target.value)}
              />
            </label>
          </div>
          <div>
            <h3 style={{ marginBottom: "0.5rem" }}>Séparée (scattered)</h3>
            <label>
              Independent Producers
              <input
                type="number"
                min="1"
                value={separeeProducers}
                onChange={(e) => setSeparee(e.target.value)}
              />
            </label>
            <label style={{ marginTop: "0.5rem" }}>
              Output (halfpennies)
              <input
                type="number"
                min="0"
                value={separeeOutput}
                onChange={(e) => setSepareeOutput(e.target.value)}
              />
            </label>
          </div>
        </div>
        {reunieOut > 0 && separeeOut > 0 && (
          <p className="muted small" style={{ marginTop: "0.75rem" }}>
            Output ratio (réunie / séparée):{" "}
            <strong
              style={{
                color:
                  reunieOut > separeeOut
                    ? "var(--color-success, #27ae60)"
                    : "var(--color-error, #c0392b)",
              }}
            >
              {(reunieOut / separeeOut).toFixed(2)}×
            </strong>
            {reunieOut > separeeOut
              ? " — concentration outproduces dispersal"
              : " — dispersal equals or exceeds concentration (unusual)"}
          </p>
        )}
      </section>
    </div>
  );
}
