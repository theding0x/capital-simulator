import { useEffect, useState } from "react";
import { financeApi } from "../../api";
import type { RealCapitalAccumulation } from "../../types";
import "./Ch30MoneyRealCapital.css";

const PHASES: { value: number; label: string }[] = [
  { value: 1, label: "1 — Inactivity (post-crisis stagnation)" },
  { value: 2, label: "2 — Revival" },
  { value: 3, label: "3 — Prosperity" },
  { value: 4, label: "4 — Overproduction" },
  { value: 5, label: "5 — Crisis" },
];

function phaseLabel(phase: number): string {
  return PHASES.find((p) => p.value === phase)?.label ?? `Phase ${phase}`;
}

function bpToPercent(bp: number): string {
  const sign = bp > 0 ? "+" : "";
  return `${sign}${(bp / 100).toFixed(1)}%`;
}

function penceToMillions(pence: number): string {
  return `£${(pence / 100 / 1_000_000).toFixed(2)}m`;
}

// Mirror of (RealCapitalAccumulation).CorrespondenceVerdict() from the Go domain.
function correspondenceVerdict(a: RealCapitalAccumulation): string {
  if (a.commercial_credit_contracts) {
    return a.cause_is_money_scarcity
      ? "commercial credit contracting — attributed to a scarcity of money-capital"
      : "commercial credit contracting — cause is blocked returns and lost confidence, not a money shortage";
  }
  return a.accumulation_correspondence.corresponds
    ? "money-capital and real-capital accumulation correspond"
    : "money-capital and real-capital accumulation diverge";
}

export function Ch30MoneyRealCapital() {
  const [records, setRecords] = useState<RealCapitalAccumulation[]>([]);
  const [listError, setListError] = useState<string | null>(null);

  const [name, setName] = useState("Crisis 1847");
  const [phase, setPhase] = useState(5);
  const [moneyGrowthBP, setMoneyGrowthBP] = useState(800);
  const [realGrowthBP, setRealGrowthBP] = useState(-600);
  const [corresponds, setCorresponds] = useState(false);
  const [explanation, setExplanation] = useState(
    "loanable money piles up in banks while productive capital and stocks contract",
  );
  const [reserveCapital, setReserveCapital] = useState(500_000);
  const [creditLimitDesc, setCreditLimitDesc] = useState("Bank of England banking reserve");
  const [contracts, setContracts] = useState(true);
  const [moneyScarcity, setMoneyScarcity] = useState(false);
  const [description, setDescription] = useState("blocked returns and lost confidence");
  const [createError, setCreateError] = useState<string | null>(null);
  const [formError, setFormError] = useState<string | null>(null);

  async function refresh() {
    try {
      const items = await financeApi.listRealCapitalAccumulations();
      setRecords(items);
    } catch (e) {
      setListError(String(e));
    }
  }

  useEffect(() => {
    refresh().catch(() => {});
  }, []);

  function validate(): boolean {
    if (reserveCapital < 0) {
      setFormError("Reserve capital must be >= 0.");
      return false;
    }
    setFormError(null);
    return true;
  }

  async function handleCreate() {
    if (!validate()) return;
    setCreateError(null);
    try {
      await financeApi.createRealCapitalAccumulation({
        name,
        phase,
        money_capital_growth_bp: moneyGrowthBP,
        real_capital_growth_bp: realGrowthBP,
        corresponds,
        correspondence_explanation: explanation,
        reserve_capital: reserveCapital,
        credit_limit_description: creditLimitDesc,
        commercial_credit_contracts: contracts,
        cause_is_money_scarcity: moneyScarcity,
        description,
      });
      await refresh();
    } catch (e) {
      setCreateError(String(e));
    }
  }

  // Live preview of the verdict for the current form state.
  const previewVerdict = correspondenceVerdict({
    id: "",
    name,
    phase,
    accumulation_correspondence: {
      money_capital_growth_bp: moneyGrowthBP,
      real_capital_growth_bp: realGrowthBP,
      corresponds,
      explanation,
    },
    credit_limit: { reserve_capital: reserveCapital, description: creditLimitDesc },
    commercial_credit_contracts: contracts,
    cause_is_money_scarcity: moneyScarcity,
    description,
    created_at: "",
  });

  return (
    <div className="v3-ch30">
      <section className="v3-ch30-intro">
        <p>
          Part V's closing question: does the accumulation of <strong>money-capital</strong>{" "}
          (loanable funds) correspond to the accumulation of <strong>real capital</strong>{" "}
          (productive capital plus commodities)? Marx's answer is{" "}
          <em>not necessarily</em>. The two can diverge sharply — loanable money can
          pile up in the banks at the very moment productive capital and saleable
          stocks are contracting.
        </p>
        <p>
          <strong>Commercial credit</strong> — the bills of exchange that producers
          draw on one another as a commodity passes down the line of production
          (spinner → weaver → manufacturer → exporter) — is the{" "}
          <strong>foundation of the credit system</strong>. It is tied directly to
          the rhythm of reproduction. In a crisis, commercial credit contracts not
          because money is scarce but because <em>returns are blocked</em> (the
          C→M sale never completes) and mutual confidence collapses. The ultimate
          cause of real crises remains the restricted consumption of the masses set
          against the limitless drive of production.
        </p>
      </section>

      {/* Seed exemplar card — Crisis 1847 */}
      <section className="v3-ch30-section" aria-label="Crisis 1847 exemplar">
        <h4>Marx's exemplar: the Crisis of 1847</h4>
        <div className="v3-ch30-seed-card">
          <div className="v3-ch30-seed-row">
            <span className="v3-ch30-seed-label">Phase</span>
            <span className="v3-ch30-seed-val">Crisis</span>
          </div>
          <div className="v3-ch30-seed-row">
            <span className="v3-ch30-seed-label">Money-capital growth</span>
            <span className="v3-ch30-seed-val">+8.0% — loanable funds swell in the banks</span>
          </div>
          <div className="v3-ch30-seed-row">
            <span className="v3-ch30-seed-label">Real-capital growth</span>
            <span className="v3-ch30-seed-val">−6.0% — productive capital and stocks contract</span>
          </div>
          <div className="v3-ch30-seed-row">
            <span className="v3-ch30-seed-label">Commercial credit</span>
            <span className="v3-ch30-seed-val">Contracting</span>
          </div>
          <div className="v3-ch30-seed-row">
            <span className="v3-ch30-seed-label">Cause</span>
            <span className="v3-ch30-seed-val">Blocked returns + lost confidence — not money scarcity</span>
          </div>
          <div className="v3-ch30-seed-verdict">
            Verdict: money- and real-capital accumulation <strong>diverge</strong>;
            credit contracts on a collapse of confidence, not a coin shortage.
          </div>
        </div>
      </section>

      {/* Create form */}
      <section className="v3-ch30-section" aria-label="Record an accumulation observation">
        <h4>Record an accumulation observation</h4>
        <p className="v3-ch30-subintro">
          Record the cycle phase and the growth of money- vs. real-capital. The two
          rates may diverge freely. The verdict below mirrors the server's
          classification.
        </p>
        <div className="v3-ch30-form">
          <label>
            Name
            <input value={name} onChange={(e) => setName(e.target.value)} />
          </label>
          <label>
            Industrial cycle phase
            <select value={phase} onChange={(e) => setPhase(Number(e.target.value))}>
              {PHASES.map((p) => (
                <option key={p.value} value={p.value}>
                  {p.label}
                </option>
              ))}
            </select>
          </label>
          <label>
            Money-capital growth (bp; may be negative)
            <input
              type="number"
              value={moneyGrowthBP}
              onChange={(e) => setMoneyGrowthBP(Number(e.target.value))}
            />
            <span className="v3-ch30-hint">{bpToPercent(moneyGrowthBP)}</span>
          </label>
          <label>
            Real-capital growth (bp; may be negative)
            <input
              type="number"
              value={realGrowthBP}
              onChange={(e) => setRealGrowthBP(Number(e.target.value))}
            />
            <span className="v3-ch30-hint">{bpToPercent(realGrowthBP)}</span>
          </label>
          <label>
            Reserve capital — credit limit (pence, &ge; 0)
            <input
              type="number"
              min={0}
              value={reserveCapital}
              onChange={(e) => setReserveCapital(Number(e.target.value))}
            />
            <span className="v3-ch30-hint">{penceToMillions(reserveCapital)} reserve</span>
          </label>
          <label>
            Credit-limit description
            <input value={creditLimitDesc} onChange={(e) => setCreditLimitDesc(e.target.value)} />
          </label>
          <label>
            Correspondence explanation
            <input value={explanation} onChange={(e) => setExplanation(e.target.value)} />
          </label>
          <label className="v3-ch30-check">
            <input
              type="checkbox"
              checked={corresponds}
              onChange={(e) => setCorresponds(e.target.checked)}
            />
            Money- and real-capital accumulation correspond
          </label>
          <label className="v3-ch30-check">
            <input
              type="checkbox"
              checked={contracts}
              onChange={(e) => setContracts(e.target.checked)}
            />
            Commercial credit is contracting
          </label>
          <label className="v3-ch30-check">
            <input
              type="checkbox"
              checked={moneyScarcity}
              onChange={(e) => setMoneyScarcity(e.target.checked)}
            />
            Cause is a scarcity of money (Marx: usually false)
          </label>
          <label>
            Description
            <input value={description} onChange={(e) => setDescription(e.target.value)} />
          </label>
        </div>
        <div className="v3-ch30-preview">Verdict: {previewVerdict}</div>
        {formError && <p className="v3-ch30-error">{formError}</p>}
        <button className="v3-ch30-btn" onClick={handleCreate}>
          Save observation
        </button>
        {createError && <p className="v3-ch30-error">{createError}</p>}
      </section>

      {/* Stored records */}
      <section className="v3-ch30-section" aria-label="Stored accumulation observations">
        <h4>Stored observations</h4>
        {listError && <p className="v3-ch30-error">{listError}</p>}
        <div className="v3-ch30-cards">
          {records.length === 0 && (
            <p className="v3-ch30-subintro">No observations stored yet.</p>
          )}
          {records.map((a) => (
            <div className="v3-ch30-card" key={a.id}>
              <div className="v3-ch30-card-name">{a.name}</div>
              <div className="v3-ch30-card-rows">
                <div className="v3-ch30-card-row">
                  <span>Phase</span>
                  <span>{phaseLabel(a.phase)}</span>
                </div>
                <div className="v3-ch30-card-row">
                  <span>Money-capital growth</span>
                  <span>{bpToPercent(a.accumulation_correspondence.money_capital_growth_bp)}</span>
                </div>
                <div className="v3-ch30-card-row">
                  <span>Real-capital growth</span>
                  <span>{bpToPercent(a.accumulation_correspondence.real_capital_growth_bp)}</span>
                </div>
                <div className="v3-ch30-card-row">
                  <span>Reserve (credit limit)</span>
                  <span>{penceToMillions(a.credit_limit.reserve_capital)}</span>
                </div>
                <div className="v3-ch30-card-row">
                  <span>Commercial credit</span>
                  <span>{a.commercial_credit_contracts ? "contracting" : "flowing"}</span>
                </div>
              </div>
              <div className="v3-ch30-verdict">{correspondenceVerdict(a)}</div>
              {a.description && <div className="v3-ch30-desc">{a.description}</div>}
            </div>
          ))}
        </div>
      </section>
    </div>
  );
}
