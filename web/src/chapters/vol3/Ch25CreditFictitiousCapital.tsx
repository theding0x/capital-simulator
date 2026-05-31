import { useEffect, useState } from "react";
import { financeApi } from "../../api";
import type { BillOfExchange, FictitiousCapital } from "../../types";
import "./Ch25CreditFictitiousCapital.css";

// FictitiousKind constants mirror credit.FictitiousKind.
const FICTITIOUS_KINDS: { value: number; label: string }[] = [
  { value: 1, label: "Bill of exchange" },
  { value: 2, label: "Bank note" },
  { value: 3, label: "Public debt (consol)" },
  { value: 4, label: "Share (stock)" },
];

const KIND_PUBLIC_DEBT = 3;

function bpToPercent(bp: number): string {
  return `${(bp / 100).toFixed(2)}%`;
}

function kindLabel(kind: number): string {
  return FICTITIOUS_KINDS.find((k) => k.value === kind)?.label ?? String(kind);
}

// capitalisedPreview mirrors the Go roundHalfUp(annual_income*10000, rate_bp).
// Math.round rounds half up for positive operands, matching the server.
function capitalisedPreview(annualIncome: number, rateBP: number): number | null {
  if (annualIncome <= 0 || rateBP <= 0) return null;
  return Math.round((annualIncome * 10000) / rateBP);
}

// Seed exemplars from Capital Vol. III Ch. 25.
const FICTITIOUS_EXEMPLARS = [
  {
    label: "UK Government Consol — £5 income @ 5% → £100, capital consumed",
    kind: 3,
    annual_income: 500,
    rate_bp: 500,
  },
  {
    label: "Railway Share — £10 dividend @ 5% → £200, real capital exists",
    kind: 4,
    annual_income: 1000,
    rate_bp: 500,
  },
];

export function Ch25CreditFictitiousCapital() {
  const [bills, setBills] = useState<BillOfExchange[]>([]);
  const [fictitious, setFictitious] = useState<FictitiousCapital[]>([]);
  const [listError, setListError] = useState<string | null>(null);

  // Bill-of-exchange create form.
  const [drawer, setDrawer] = useState("Manchester Spinner");
  const [drawee, setDrawee] = useState("London Merchant");
  const [faceValue, setFaceValue] = useState(13212346000);
  const [maturityDays, setMaturityDays] = useState(90);
  const [billError, setBillError] = useState<string | null>(null);

  // Fictitious-capital create form.
  const [kind, setKind] = useState(KIND_PUBLIC_DEBT);
  const [name, setName] = useState("UK Government Consol");
  const [annualIncome, setAnnualIncome] = useState(500);
  const [rateBP, setRateBP] = useState(500);
  const [fcError, setFcError] = useState<string | null>(null);

  const previewCapitalised = capitalisedPreview(annualIncome, rateBP);
  const previewRealCapital = kind !== KIND_PUBLIC_DEBT;

  async function refresh() {
    try {
      const [b, f] = await Promise.all([
        financeApi.listBillsOfExchange(),
        financeApi.listFictitiousCapitals(),
      ]);
      setBills(b);
      setFictitious(f);
    } catch (e) {
      setListError(String(e));
    }
  }

  useEffect(() => {
    refresh().catch(() => {});
  }, []);

  async function handleCreateBill() {
    setBillError(null);
    try {
      await financeApi.createBillOfExchange({
        drawer,
        drawee,
        face_value: faceValue,
        maturity_days: maturityDays,
        description: "",
      });
      await refresh();
    } catch (e) {
      setBillError(String(e));
    }
  }

  async function handleCreateFictitious() {
    setFcError(null);
    try {
      await financeApi.createFictitiousCapital({
        kind,
        name,
        annual_income: annualIncome,
        rate_bp: rateBP,
        description: "",
      });
      await refresh();
    } catch (e) {
      setFcError(String(e));
    }
  }

  return (
    <div className="v3-ch25">
      <section className="v3-ch25-intro">
        <p>
          The credit system grows out of money's function as a{" "}
          <strong>means of payment</strong>. A <strong>bill of exchange</strong> —
          a promissory note with a fixed maturity — circulates as commercial money
          long before its due date, settling chains of debt without coin. In 1839,
          such bills made up the majority of British currency.
        </p>
        <p>
          <strong>Fictitious capital</strong> arises when an income stream is
          capitalised at the prevailing rate of interest:{" "}
          <code>capitalised_value = annual_income × 10000 ÷ rate_bp</code>. A
          government consol of £5 a year at 5% trades as £100 of "capital" — yet
          the original £100 loan was spent long ago. The capital is gone; only the
          claim on future taxes remains. A share is a title to future
          surplus-value, not the track and locomotives themselves.
        </p>
        <p>
          Invariant: for public debt (kind = 3),{" "}
          <code>real_capital_exists = false</code> — the value has been consumed.
        </p>
      </section>

      {/* Fictitious-capital exemplars */}
      <section className="v3-ch25-section" aria-label="Fictitious capital exemplars">
        <h4>Marx's exemplars (Ch. 25)</h4>
        <div className="v3-ch25-cards">
          {FICTITIOUS_EXEMPLARS.map((ex) => (
            <div className="v3-ch25-card" key={ex.label}>
              <div className="v3-ch25-card-label">{ex.label}</div>
              <div className="v3-ch25-meta">
                <span className="v3-ch25-kind-tag">{kindLabel(ex.kind)}</span>
                <span>Income: {ex.annual_income}</span>
                <span>Rate: {bpToPercent(ex.rate_bp)}</span>
              </div>
            </div>
          ))}
        </div>
      </section>

      {/* Fictitious-capital calculator */}
      <section className="v3-ch25-section" aria-label="Capitalise an income stream">
        <h4>Capitalise an income stream</h4>
        <p className="v3-ch25-subintro">
          capitalised_value and real_capital_exists are derived on the server.
        </p>
        <div className="v3-ch25-form">
          <label>
            Kind
            <select value={kind} onChange={(e) => setKind(Number(e.target.value))}>
              {FICTITIOUS_KINDS.map((k) => (
                <option key={k.value} value={k.value}>
                  {k.label}
                </option>
              ))}
            </select>
          </label>
          <label>
            Name
            <input value={name} onChange={(e) => setName(e.target.value)} />
          </label>
          <label>
            Annual income (pence, &gt; 0)
            <input
              type="number"
              min={1}
              value={annualIncome}
              onChange={(e) => setAnnualIncome(Number(e.target.value))}
            />
          </label>
          <label>
            Rate (basis points, &gt; 0)
            <input
              type="number"
              min={1}
              value={rateBP}
              onChange={(e) => setRateBP(Number(e.target.value))}
            />
            <span className="v3-ch25-hint">100 bp = 1%. 5% = 500 bp.</span>
          </label>
        </div>

        {previewCapitalised !== null && (
          <div className="v3-ch25-preview" aria-label="Preview">
            <div className="v3-ch25-preview-row">
              <span className="v3-ch25-preview-label">Capitalised value (preview)</span>
              <span className="v3-ch25-preview-value v3-ch25-preview-final">
                {previewCapitalised}
              </span>
            </div>
            <div className="v3-ch25-preview-row">
              <span className="v3-ch25-preview-label">Real capital exists</span>
              <span className="v3-ch25-preview-value">
                {previewRealCapital ? "true" : "false (consumed)"}
              </span>
            </div>
          </div>
        )}

        <button
          className="v3-ch25-btn"
          onClick={handleCreateFictitious}
          disabled={annualIncome <= 0 || rateBP <= 0}
        >
          Capitalise &amp; save
        </button>
        {fcError && <p className="v3-ch25-error">{fcError}</p>}
      </section>

      {/* Stored fictitious capital */}
      <section className="v3-ch25-section" aria-label="Stored fictitious capital">
        <h4>Stored fictitious capital</h4>
        {listError && <p className="v3-ch25-error">{listError}</p>}
        <div className="v3-ch25-cards">
          {fictitious.length === 0 && (
            <p className="v3-ch25-subintro">No fictitious capital stored yet.</p>
          )}
          {fictitious.map((fc) => (
            <div className="v3-ch25-card v3-ch25-card--stored" key={fc.id}>
              <div className="v3-ch25-card-label">{fc.name}</div>
              <div className="v3-ch25-final">
                Capitalised: {fc.capitalised_value} ({bpToPercent(fc.rate_bp)})
              </div>
              <div className="v3-ch25-meta">
                <span className="v3-ch25-kind-tag">{kindLabel(fc.kind)}</span>
                <span>Income: {fc.annual_income}</span>
                <span
                  className={
                    fc.real_capital_exists ? "v3-ch25-real" : "v3-ch25-fictitious"
                  }
                >
                  {fc.real_capital_exists ? "real capital" : "consumed"}
                </span>
              </div>
            </div>
          ))}
        </div>
      </section>

      {/* Bill-of-exchange form */}
      <section className="v3-ch25-section" aria-label="Draw a bill of exchange">
        <h4>Draw a bill of exchange</h4>
        <p className="v3-ch25-subintro">
          Commercial credit circulating as money. face_value and maturity_days
          must be &gt; 0.
        </p>
        <div className="v3-ch25-form">
          <label>
            Drawer (issuer)
            <input value={drawer} onChange={(e) => setDrawer(e.target.value)} />
          </label>
          <label>
            Drawee (debtor)
            <input value={drawee} onChange={(e) => setDrawee(e.target.value)} />
          </label>
          <label>
            Face value (pence, &gt; 0)
            <input
              type="number"
              min={1}
              value={faceValue}
              onChange={(e) => setFaceValue(Number(e.target.value))}
            />
          </label>
          <label>
            Maturity (days, &gt; 0)
            <input
              type="number"
              min={1}
              value={maturityDays}
              onChange={(e) => setMaturityDays(Number(e.target.value))}
            />
          </label>
        </div>
        <button
          className="v3-ch25-btn"
          onClick={handleCreateBill}
          disabled={faceValue <= 0 || maturityDays <= 0}
        >
          Draw &amp; save bill
        </button>
        {billError && <p className="v3-ch25-error">{billError}</p>}
      </section>

      {/* Stored bills */}
      <section className="v3-ch25-section" aria-label="Stored bills of exchange">
        <h4>Stored bills of exchange</h4>
        <div className="v3-ch25-cards">
          {bills.length === 0 && (
            <p className="v3-ch25-subintro">No bills stored yet.</p>
          )}
          {bills.map((b) => (
            <div className="v3-ch25-card v3-ch25-card--stored" key={b.id}>
              <div className="v3-ch25-card-label">
                {b.drawer} → {b.drawee}
              </div>
              <div className="v3-ch25-final">Face: {b.face_value}</div>
              <div className="v3-ch25-meta">
                <span>Maturity: {b.maturity_days} days</span>
              </div>
            </div>
          ))}
        </div>
      </section>
    </div>
  );
}
