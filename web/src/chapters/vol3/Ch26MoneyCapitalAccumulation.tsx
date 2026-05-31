import { useEffect, useState } from "react";
import { financeApi } from "../../api";
import type { MoneyCapitalAccumulation } from "../../types";
import "./Ch26MoneyCapitalAccumulation.css";

// Industrial cycle phase labels mirror credit.IndustrialCyclePhase.
const PHASE_LABELS: Record<number, string> = {
  1: "Inactivity",
  2: "Revival",
  3: "Prosperity",
  4: "Overproduction",
  5: "Crisis",
};

function phaseLabel(phase: number): string {
  return PHASE_LABELS[phase] ?? String(phase);
}

function bpToPercent(bp: number): string {
  return `${(bp / 100).toFixed(2)}%`;
}

// Seed exemplars from Capital Vol. III Ch. 26.
const ACCUMULATION_EXEMPLARS = [
  {
    label: "Post-1847 Inactivity — loanable money piles up, rate depressed to 1%",
    phase: 1,
    interest_rate_bp: 100,
    abundant: true,
  },
  {
    label: "Revival 1850s — demand rising, supply ample, rate moderate at 3%",
    phase: 2,
    interest_rate_bp: 300,
    abundant: true,
  },
  {
    label: "Crisis Peak 1857 — acute demand for means of payment, rate spikes to 10%",
    phase: 5,
    interest_rate_bp: 1000,
    abundant: false,
  },
];

export function Ch26MoneyCapitalAccumulation() {
  const [accumulations, setAccumulations] = useState<MoneyCapitalAccumulation[]>([]);
  const [listError, setListError] = useState<string | null>(null);

  // Create form state.
  const [name, setName] = useState("Post-1847 Inactivity");
  const [phase, setPhase] = useState(1);
  const [interestRateBP, setInterestRateBP] = useState(100);
  const [period, setPeriod] = useState("post-1847 crisis");
  const [loanableAmount, setLoanableAmount] = useState(5000000);
  const [loanableAbundant, setLoanableAbundant] = useState(true);
  const [idleAmount, setIdleAmount] = useState(2000000);
  const [idleDescription, setIdleDescription] = useState("idle cotton-mill capital");
  const [gapBP, setGapBP] = useState(500);
  const [gapDescription, setGapDescription] = useState("money-capital accumulates faster than productive capital");
  const [description, setDescription] = useState("");
  const [createError, setCreateError] = useState<string | null>(null);

  async function refresh() {
    try {
      const items = await financeApi.listMoneyCapitalAccumulations();
      setAccumulations(items);
    } catch (e) {
      setListError(String(e));
    }
  }

  useEffect(() => {
    refresh().catch(() => {});
  }, []);

  async function handleCreate() {
    setCreateError(null);
    try {
      await financeApi.createMoneyCapitalAccumulation({
        name,
        phase,
        interest_rate_bp: interestRateBP,
        period,
        loanable_amount: loanableAmount,
        loanable_abundant: loanableAbundant,
        idle_amount: idleAmount,
        idle_description: idleDescription,
        gap_bp: gapBP,
        gap_description: gapDescription,
        description,
      });
      await refresh();
    } catch (e) {
      setCreateError(String(e));
    }
  }

  return (
    <div className="v3-ch26">
      <section className="v3-ch26-intro">
        <p>
          As money-capital accumulates during post-crisis <strong>inactivity</strong>, loanable
          funds become abundant while demand for them is low, driving the interest rate down.
          When <strong>revival</strong> begins, idle industrial capital seeks loanable outlets
          and demand rises. By <strong>prosperity</strong>, credit is stretched and the rate
          climbs above the average. In <strong>overproduction</strong> the rate peaks; in{" "}
          <strong>crisis</strong> it spikes sharply as every capitalist scrambles for means
          of payment.
        </p>
        <p>
          The accumulation of money-capital thus moves <em>inversely</em> with the demand for
          loanable capital across the cycle. The interest rate is its barometer — not a law,
          but an empirical reflex of supply and demand in the money market.
        </p>
      </section>

      {/* Seed exemplars */}
      <section className="v3-ch26-section" aria-label="Accumulation exemplars">
        <h4>Marx's exemplars (Ch. 26)</h4>
        <div className="v3-ch26-cards">
          {ACCUMULATION_EXEMPLARS.map((ex) => (
            <div className="v3-ch26-card" key={ex.label}>
              <div className="v3-ch26-card-label">{ex.label}</div>
              <div className="v3-ch26-meta">
                <span className="v3-ch26-phase-tag">{phaseLabel(ex.phase)}</span>
                <span>Rate: {bpToPercent(ex.interest_rate_bp)}</span>
                <span className={ex.abundant ? "v3-ch26-abundant" : "v3-ch26-scarce"}>
                  {ex.abundant ? "supply abundant" : "supply scarce"}
                </span>
              </div>
            </div>
          ))}
        </div>
      </section>

      {/* Create form */}
      <section className="v3-ch26-section" aria-label="Record a cycle observation">
        <h4>Record a cycle observation</h4>
        <p className="v3-ch26-subintro">
          Phase, interest rate, and loanable-capital supply are combined to classify
          the accumulation dynamic server-side.
        </p>
        <div className="v3-ch26-form">
          <label>
            Name
            <input value={name} onChange={(e) => setName(e.target.value)} />
          </label>
          <label>
            Phase
            <select value={phase} onChange={(e) => setPhase(Number(e.target.value))}>
              <option value={1}>1 — Inactivity</option>
              <option value={2}>2 — Revival</option>
              <option value={3}>3 — Prosperity</option>
              <option value={4}>4 — Overproduction</option>
              <option value={5}>5 — Crisis</option>
            </select>
          </label>
          <label>
            Interest rate (basis points, &gt;= 0)
            <input
              type="number"
              min={0}
              value={interestRateBP}
              onChange={(e) => setInterestRateBP(Number(e.target.value))}
            />
            <span className="v3-ch26-hint">100 bp = 1%. 10% = 1000 bp.</span>
          </label>
          <label>
            Period (e.g. "1847 crisis")
            <input value={period} onChange={(e) => setPeriod(e.target.value)} />
          </label>
          <label>
            Loanable capital supply (pence, &gt;= 0)
            <input
              type="number"
              min={0}
              value={loanableAmount}
              onChange={(e) => setLoanableAmount(Number(e.target.value))}
            />
          </label>
          <label>
            Supply abundant?
            <select
              value={loanableAbundant ? "true" : "false"}
              onChange={(e) => setLoanableAbundant(e.target.value === "true")}
            >
              <option value="true">Yes — supply exceeds demand</option>
              <option value="false">No — supply scarce</option>
            </select>
          </label>
          <label>
            Idle industrial capital (pence, &gt;= 0)
            <input
              type="number"
              min={0}
              value={idleAmount}
              onChange={(e) => setIdleAmount(Number(e.target.value))}
            />
          </label>
          <label>
            Idle capital description
            <input value={idleDescription} onChange={(e) => setIdleDescription(e.target.value)} />
          </label>
          <label>
            Accumulation gap (bp, may be negative)
            <input
              type="number"
              value={gapBP}
              onChange={(e) => setGapBP(Number(e.target.value))}
            />
            <span className="v3-ch26-hint">
              Positive: money-capital grows faster. Negative: productive capital outpaces supply.
            </span>
          </label>
          <label>
            Gap description
            <input value={gapDescription} onChange={(e) => setGapDescription(e.target.value)} />
          </label>
          <label>
            Description
            <input value={description} onChange={(e) => setDescription(e.target.value)} />
          </label>
        </div>
        <button
          className="v3-ch26-btn"
          onClick={handleCreate}
          disabled={interestRateBP < 0 || idleAmount < 0 || loanableAmount < 0}
        >
          Save observation
        </button>
        {createError && <p className="v3-ch26-error">{createError}</p>}
      </section>

      {/* Stored observations */}
      <section className="v3-ch26-section" aria-label="Stored accumulation observations">
        <h4>Stored observations</h4>
        {listError && <p className="v3-ch26-error">{listError}</p>}
        <div className="v3-ch26-cards">
          {accumulations.length === 0 && (
            <p className="v3-ch26-subintro">No observations stored yet.</p>
          )}
          {accumulations.map((a) => (
            <div className="v3-ch26-card v3-ch26-card--stored" key={a.id}>
              <div className="v3-ch26-card-label">{a.name}</div>
              <div className="v3-ch26-meta">
                <span className="v3-ch26-phase-tag">
                  {phaseLabel(a.cycle_observation.phase)}
                </span>
                <span>Rate: {bpToPercent(a.cycle_observation.interest_rate_bp)}</span>
                <span
                  className={
                    a.loanable_supply.abundant ? "v3-ch26-abundant" : "v3-ch26-scarce"
                  }
                >
                  {a.loanable_supply.abundant ? "supply abundant" : "supply scarce"}
                </span>
                <span>Idle: {a.idle_capital.amount}</span>
                <span>Gap: {a.gap.gap_bp} bp</span>
              </div>
              {a.description && (
                <div className="v3-ch26-classify">{a.description}</div>
              )}
            </div>
          ))}
        </div>
      </section>
    </div>
  );
}
