import { useEffect, useState } from "react";
import { financeApi } from "../../api";
import type { RateOfInterest, InterestRateAnalysis } from "../../types";
import "./Ch22RateOfInterest.css";

// IndustrialCyclePhase constants mirror credit.IndustrialCyclePhase.
const PHASES: { value: number; label: string }[] = [
  { value: 1, label: "Inactivity" },
  { value: 2, label: "Revival" },
  { value: 3, label: "Prosperity" },
  { value: 4, label: "Overproduction" },
  { value: 5, label: "Crisis" },
];

function bpToPercent(bp: number): string {
  return `${(bp / 100).toFixed(2)}%`;
}

function phaseName(phase: number): string {
  return PHASES.find((p) => p.value === phase)?.label ?? String(phase);
}

function phaseClass(phase: number): string {
  if (phase === 5) return "v3-ch22-phase-tag v3-ch22-phase-tag--crisis";
  if (phase === 3) return "v3-ch22-phase-tag v3-ch22-phase-tag--prosperity";
  return "v3-ch22-phase-tag";
}

// Seed fixtures from Capital Vol. III Ch. 22.
const SEED_EXEMPLARS = [
  {
    label: "Canonical — avg 20%, interest ¼ of profit",
    rate_bp: 500,
    average_profit_rate_bp: 2000,
    cycle_phase: 3,
    enterprise_bp: 1500,
  },
  {
    label: "Prosperity 1843 — interest 1.5%",
    rate_bp: 150,
    average_profit_rate_bp: 2000,
    cycle_phase: 3,
    enterprise_bp: 1850,
  },
  {
    label: "Crisis 1847 — interest 8%+",
    rate_bp: 800,
    average_profit_rate_bp: 2000,
    cycle_phase: 5,
    enterprise_bp: 1200,
  },
];

export function Ch22RateOfInterest() {
  const [items, setItems] = useState<RateOfInterest[]>([]);
  const [listError, setListError] = useState<string | null>(null);

  // Create-form state.
  const [rateBP, setRateBP] = useState(500);
  const [avgProfitBP, setAvgProfitBP] = useState(2000);
  const [cyclePhase, setCyclePhase] = useState(3);
  const [period, setPeriod] = useState("");
  const [createError, setCreateError] = useState<string | null>(null);

  // Compute-form state.
  const [computeAvg, setComputeAvg] = useState(2000);
  const [computeInt, setComputeInt] = useState(500);
  const [analysis, setAnalysis] = useState<InterestRateAnalysis | null>(null);
  const [computeError, setComputeError] = useState<string | null>(null);

  // Live-preview: profit of enterprise = avg - interest (integer math, no rounding needed).
  const previewEnterprise =
    rateBP >= 0 && avgProfitBP > rateBP ? avgProfitBP - rateBP : null;

  async function refreshItems() {
    try {
      const data = await financeApi.listRatesOfInterest();
      setItems(data);
    } catch (e) {
      setListError(String(e));
    }
  }

  useEffect(() => {
    refreshItems().catch(() => {});
  }, []);

  async function handleCreate() {
    setCreateError(null);
    try {
      await financeApi.createRateOfInterest({
        rate_bp: rateBP,
        average_profit_rate_bp: avgProfitBP,
        cycle_phase: cyclePhase,
        period,
      });
      await refreshItems();
    } catch (e) {
      setCreateError(String(e));
    }
  }

  async function handleCompute() {
    setComputeError(null);
    setAnalysis(null);
    try {
      const result = await financeApi.computeInterestRateAnalysis({
        average_profit_rate_bp: computeAvg,
        interest_rate_bp: computeInt,
      });
      setAnalysis(result);
    } catch (e) {
      setComputeError(String(e));
    }
  }

  return (
    <div className="v3-ch22">
      <section className="v3-ch22-intro">
        <p>
          The rate of interest has no "natural" or law-governed level — it is determined
          empirically by supply and demand in the money market. Its{" "}
          <strong>maximum limit</strong> is the average profit rate (above which no
          borrower could profit); its <strong>minimum</strong> approaches but never
          reaches zero. Interest is a portion of the average profit; the residual is the{" "}
          <strong>profit of enterprise</strong> retained by the functioning capitalist.
        </p>
      </section>

      {/* Canonical Marx exemplars */}
      <section className="v3-ch22-section" aria-label="Rate of interest exemplars">
        <h4>Marx's exemplars (Ch. 22)</h4>
        <p className="v3-ch22-subintro">
          ProfitOfEnterprise = AverageProfitRate &minus; InterestRate (basis points).
        </p>
        <div className="v3-ch22-cards">
          {SEED_EXEMPLARS.map((ex) => (
            <div key={ex.label} className="v3-ch22-card">
              <div className="v3-ch22-card-label">{ex.label}</div>
              <div className="v3-ch22-meta">
                <span>Avg profit: {bpToPercent(ex.average_profit_rate_bp)}</span>
                <span>Interest: {bpToPercent(ex.rate_bp)}</span>
              </div>
              <div className="v3-ch22-enterprise">
                Enterprise: {bpToPercent(ex.enterprise_bp)}
              </div>
              <div className="v3-ch22-meta" style={{ marginTop: "0.3rem" }}>
                <span className={phaseClass(ex.cycle_phase)}>
                  {phaseName(ex.cycle_phase)}
                </span>
              </div>
            </div>
          ))}
        </div>
      </section>

      {/* Compute form */}
      <section className="v3-ch22-section" aria-label="Compute profit of enterprise">
        <h4>Compute profit of enterprise</h4>
        <div className="v3-ch22-form">
          <label>
            Average profit rate (basis points)
            <input
              type="number"
              min={1}
              value={computeAvg}
              onChange={(e) => setComputeAvg(Number(e.target.value))}
            />
            <span className="v3-ch22-hint">{bpToPercent(computeAvg)}</span>
          </label>
          <label>
            Interest rate (basis points)
            <input
              type="number"
              min={0}
              value={computeInt}
              onChange={(e) => setComputeInt(Number(e.target.value))}
            />
            <span className="v3-ch22-hint">{bpToPercent(computeInt)}</span>
          </label>
        </div>
        <button
          type="button"
          className="v3-ch22-btn"
          disabled={computeAvg <= 0 || computeInt < 0 || computeInt >= computeAvg}
          onClick={handleCompute}
        >
          Compute
        </button>
        {computeError && <p className="v3-ch22-error">{computeError}</p>}
        {analysis && (
          <div className="v3-ch22-preview" style={{ marginTop: "0.75rem" }}>
            <div className="v3-ch22-preview-row">
              <span className="v3-ch22-preview-label">Average profit rate</span>
              <span className="v3-ch22-preview-value">
                {bpToPercent(analysis.average_profit_rate_bp)}
              </span>
            </div>
            <div className="v3-ch22-preview-row">
              <span className="v3-ch22-preview-label">Interest rate</span>
              <span className="v3-ch22-preview-value">
                {bpToPercent(analysis.interest_rate_bp)}
              </span>
            </div>
            <div className="v3-ch22-preview-row">
              <span className="v3-ch22-preview-label">Profit of enterprise</span>
              <span className="v3-ch22-preview-value v3-ch22-preview-enterprise">
                {bpToPercent(analysis.profit_of_enterprise_bp)}
              </span>
            </div>
          </div>
        )}
      </section>

      {/* Record observation form */}
      <section className="v3-ch22-section" aria-label="Record an interest-rate observation">
        <h4>Record an interest-rate observation</h4>
        <div className="v3-ch22-form">
          <label>
            Interest rate (basis points)
            <input
              type="number"
              min={0}
              value={rateBP}
              onChange={(e) => setRateBP(Number(e.target.value))}
            />
            <span className="v3-ch22-hint">{bpToPercent(rateBP)}</span>
          </label>
          <label>
            Average profit rate (basis points)
            <input
              type="number"
              min={1}
              value={avgProfitBP}
              onChange={(e) => setAvgProfitBP(Number(e.target.value))}
            />
            <span className="v3-ch22-hint">{bpToPercent(avgProfitBP)}</span>
          </label>
          <label>
            Industrial cycle phase
            <select
              value={cyclePhase}
              onChange={(e) => setCyclePhase(Number(e.target.value))}
            >
              {PHASES.map((p) => (
                <option key={p.value} value={p.value}>
                  {p.label}
                </option>
              ))}
            </select>
          </label>
          <label>
            Period
            <input
              type="text"
              value={period}
              placeholder="e.g. 1843, 1847 crisis"
              onChange={(e) => setPeriod(e.target.value)}
            />
          </label>
        </div>

        {/* Live preview */}
        {previewEnterprise !== null && (
          <div className="v3-ch22-preview">
            <div className="v3-ch22-preview-row">
              <span className="v3-ch22-preview-label">Profit of enterprise</span>
              <span className="v3-ch22-preview-value v3-ch22-preview-enterprise">
                {bpToPercent(previewEnterprise)}
              </span>
            </div>
          </div>
        )}

        <button
          type="button"
          className="v3-ch22-btn"
          disabled={rateBP < 0 || avgProfitBP <= 0 || rateBP >= avgProfitBP}
          onClick={handleCreate}
        >
          Record observation
        </button>
        {createError && <p className="v3-ch22-error">{createError}</p>}
      </section>

      {/* Stored observations */}
      {listError && <p className="v3-ch22-error">{listError}</p>}
      {items.length > 0 && (
        <section className="v3-ch22-section" aria-label="Stored interest-rate observations">
          <h4>Recorded observations ({items.length})</h4>
          <div className="v3-ch22-cards">
            {items.map((r) => (
              <div key={r.id} className="v3-ch22-card v3-ch22-card--stored">
                <div className="v3-ch22-card-label">{r.period || r.id}</div>
                <div className="v3-ch22-meta">
                  <span>Avg profit: {bpToPercent(r.average_profit_rate_bp)}</span>
                  <span>Interest: {bpToPercent(r.rate_bp)}</span>
                  <span>
                    Enterprise:{" "}
                    {bpToPercent(
                      Math.floor(r.average_profit_rate_bp - r.rate_bp),
                    )}
                  </span>
                </div>
                <div className="v3-ch22-meta" style={{ marginTop: "0.3rem" }}>
                  <span className={phaseClass(r.cycle_phase)}>
                    {phaseName(r.cycle_phase)}
                  </span>
                </div>
              </div>
            ))}
          </div>
        </section>
      )}
    </div>
  );
}
