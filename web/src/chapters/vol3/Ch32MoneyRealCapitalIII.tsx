import { useEffect, useState } from "react";
import { financeApi } from "../../api";
import type { CapitalRelease } from "../../types";
import "./Ch32MoneyRealCapitalIII.css";

const CAUSES: { value: number; label: string }[] = [
  { value: 1, label: "1 — Input (raw-material) price fall" },
  { value: 2, label: "2 — Merchants' idle money" },
  { value: 3, label: "3 — Rentier-class savings" },
  { value: 4, label: "4 — Revenue passing through money form" },
];

function causeLabel(cause: number): string {
  return CAUSES.find((c) => c.value === cause)?.label ?? `Cause ${cause}`;
}

function penceToPounds(pence: number): string {
  return `£${(pence / 100).toLocaleString("en-GB", { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`;
}

function bpToPct(bp: number): string {
  return `${(bp / 100).toLocaleString("en-GB", { maximumFractionDigits: 2 })}%`;
}

// Mirror of (CapitalRelease).AccumulationVerdict() from the Go domain.
function accumulationVerdict(reinvested: boolean): string {
  return reinvested
    ? "reflects real accumulation — the released capital is reinvested in production"
    : "overstates real accumulation — the released capital swells loan capital but is not reinvested";
}

export function Ch32MoneyRealCapitalIII() {
  const [records, setRecords] = useState<CapitalRelease[]>([]);
  const [listError, setListError] = useState<string | null>(null);

  // Capital-release form state.
  const [name, setName] = useState("Capital released by raw-material price fall");
  const [releasedAmount, setReleasedAmount] = useState(5_000_000);
  const [cause, setCause] = useState(1);
  const [reinvested, setReinvested] = useState(false);
  const [overstatementPct, setOverstatementPct] = useState(2500);
  const [overstatementDesc, setOverstatementDesc] = useState(
    "freed money lent out but not reinvested",
  );
  const [description, setDescription] = useState(
    "raw-material prices fall; the same reproduction ties up less money; the difference is lent out",
  );
  const [createError, setCreateError] = useState<string | null>(null);
  const [formError, setFormError] = useState<string | null>(null);

  async function refresh() {
    try {
      const items = await financeApi.listCapitalReleases();
      setRecords(items);
    } catch (e) {
      setListError(String(e));
    }
  }

  useEffect(() => {
    refresh().catch(() => {});
  }, []);

  function validate(): boolean {
    if (releasedAmount < 0) {
      setFormError("Released amount must be >= 0.");
      return false;
    }
    if (overstatementPct < 0) {
      setFormError("Overstatement must be >= 0.");
      return false;
    }
    setFormError(null);
    return true;
  }

  async function handleCreate() {
    if (!validate()) return;
    setCreateError(null);
    try {
      await financeApi.createCapitalRelease({
        name,
        released_amount: releasedAmount,
        cause,
        reinvested,
        overstatement_pct: overstatementPct,
        overstatement_description: overstatementDesc,
        description,
      });
      await refresh();
    } catch (e) {
      setCreateError(String(e));
    }
  }

  const previewVerdict = accumulationVerdict(reinvested);

  return (
    <div className="v3-ch32">
      <section className="v3-ch32-intro">
        <p>
          The conclusion of the trilogy: accumulated loanable money-capital{" "}
          <strong>always overstates real accumulation</strong>. Revenue intended
          for consumption passes temporarily through a money form and appears as
          loan capital, though it is <em>no new productive capital</em>. Capital{" "}
          <strong>released</strong> by falling input prices, merchants' idle money,
          and the savings of the growing <strong>rentier class</strong> all swell
          loan capital independently of any real expansion.
        </p>
        <p>
          A release represents real accumulation <em>only when the freed capital is
          actually reinvested</em> in production. Lent out but idle, it merely
          inflates the apparent stock of capital while no new productive capital
          stands behind it.
        </p>
        <p>
          In a <strong>crisis</strong> the demand for loan capital is a demand for{" "}
          <strong>means of payment</strong> — convertibility into gold — <em>not</em>{" "}
          for productive capital. Banking theory that confuses the two collapses on
          exactly this distinction.
        </p>
      </section>

      {/* Exemplar cards */}
      <section className="v3-ch32-section" aria-label="Marx's exemplars">
        <h4>Marx's exemplars: loan capital that overstates real accumulation</h4>
        <div className="v3-ch32-exemplars">
          <div className="v3-ch32-exemplar">
            <div className="v3-ch32-exemplar-name">Capital released by price fall</div>
            <div className="v3-ch32-exemplar-row">
              <span>Cause</span>
              <span>Input price fall (1)</span>
            </div>
            <div className="v3-ch32-exemplar-row">
              <span>Verdict</span>
              <span>Overstates until reinvested</span>
            </div>
            <div className="v3-ch32-exemplar-note">
              Raw-material prices fall; the same reproduction ties up less money; the
              difference is lent out.
            </div>
          </div>
          <div className="v3-ch32-exemplar">
            <div className="v3-ch32-exemplar-name">Rentier-class savings</div>
            <div className="v3-ch32-exemplar-row">
              <span>Cause</span>
              <span>Rentier savings (3)</span>
            </div>
            <div className="v3-ch32-exemplar-row">
              <span>Verdict</span>
              <span>Pure overstatement</span>
            </div>
            <div className="v3-ch32-exemplar-note">
              The growing rentier class lives off interest, swelling loan capital
              independently of real expansion.
            </div>
          </div>
          <div className="v3-ch32-exemplar">
            <div className="v3-ch32-exemplar-name">Crisis 1857: means of payment</div>
            <div className="v3-ch32-exemplar-row">
              <span>Phase</span>
              <span>Crisis (5)</span>
            </div>
            <div className="v3-ch32-exemplar-row">
              <span>Demand</span>
              <span>Convertibility, not capital</span>
            </div>
            <div className="v3-ch32-exemplar-note">
              The cry is for convertibility into gold, not for expansion of productive
              capacity.
            </div>
          </div>
        </div>
      </section>

      {/* Create form */}
      <section className="v3-ch32-section" aria-label="Record a capital release">
        <h4>Record a capital release</h4>
        <p className="v3-ch32-subintro">
          Record the amount freed, its cause, whether it was reinvested, and the gap
          between apparent and real accumulation. The accumulation verdict below
          mirrors the server's classification.
        </p>
        <div className="v3-ch32-form">
          <label>
            Name
            <input value={name} onChange={(e) => setName(e.target.value)} />
          </label>
          <label>
            Released amount (pence, &ge; 0)
            <input
              type="number"
              min={0}
              value={releasedAmount}
              onChange={(e) => setReleasedAmount(Number(e.target.value))}
            />
            <span className="v3-ch32-hint">{penceToPounds(releasedAmount)}</span>
          </label>
          <label>
            Cause
            <select value={cause} onChange={(e) => setCause(Number(e.target.value))}>
              {CAUSES.map((c) => (
                <option key={c.value} value={c.value}>
                  {c.label}
                </option>
              ))}
            </select>
          </label>
          <label className="v3-ch32-check">
            <input
              type="checkbox"
              checked={reinvested}
              onChange={(e) => setReinvested(e.target.checked)}
            />
            Reinvested in production?
          </label>
          <label>
            Overstatement (basis points, &ge; 0)
            <input
              type="number"
              min={0}
              value={overstatementPct}
              onChange={(e) => setOverstatementPct(Number(e.target.value))}
            />
            <span className="v3-ch32-hint">{bpToPct(overstatementPct)}</span>
          </label>
          <label>
            Overstatement description
            <input value={overstatementDesc} onChange={(e) => setOverstatementDesc(e.target.value)} />
          </label>
          <label>
            Description
            <input value={description} onChange={(e) => setDescription(e.target.value)} />
          </label>
        </div>
        <div className="v3-ch32-preview">Verdict: {previewVerdict}</div>
        {formError && <p className="v3-ch32-error">{formError}</p>}
        <button className="v3-ch32-btn" onClick={handleCreate}>
          Save release
        </button>
        {createError && <p className="v3-ch32-error">{createError}</p>}
      </section>

      {/* Stored records */}
      <section className="v3-ch32-section" aria-label="Stored capital releases">
        <h4>Stored releases</h4>
        {listError && <p className="v3-ch32-error">{listError}</p>}
        <div className="v3-ch32-cards">
          {records.length === 0 && (
            <p className="v3-ch32-subintro">No releases stored yet.</p>
          )}
          {records.map((c) => (
            <div className="v3-ch32-card" key={c.id}>
              <div className="v3-ch32-card-name">{c.name}</div>
              <div className="v3-ch32-card-rows">
                <div className="v3-ch32-card-row">
                  <span>Released</span>
                  <span>{penceToPounds(c.released_amount)}</span>
                </div>
                <div className="v3-ch32-card-row">
                  <span>Cause</span>
                  <span>{causeLabel(c.cause)}</span>
                </div>
                <div className="v3-ch32-card-row">
                  <span>Reinvested</span>
                  <span>{c.reinvested ? "yes" : "no"}</span>
                </div>
                <div className="v3-ch32-card-row">
                  <span>Overstatement</span>
                  <span>{bpToPct(c.loan_capital_overstatement.overstatement_pct)}</span>
                </div>
              </div>
              <div className="v3-ch32-verdict">{accumulationVerdict(c.reinvested)}</div>
              {c.description && <div className="v3-ch32-desc">{c.description}</div>}
            </div>
          ))}
        </div>
      </section>
    </div>
  );
}
