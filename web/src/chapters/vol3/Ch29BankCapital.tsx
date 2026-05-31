import { useEffect, useState } from "react";
import { financeApi } from "../../api";
import type {
  BankCapital,
  BankCapitalComponent,
  BankCapitalComponentKind,
  FictitiousCapitalValuation,
} from "../../types";
import "./Ch29BankCapital.css";

function penceToPounds(pence: number): string {
  return `£${(pence / 100).toLocaleString(undefined, { maximumFractionDigits: 2 })}`;
}

function bpToPercent(bp: number): string {
  return `${(bp / 100).toFixed(2)}%`;
}

// roundHalfUp mirrors the Go domain: roundHalfUp(numer, denom).
function roundHalfUp(numer: number, denom: number): number {
  return Math.floor((numer + Math.floor(denom / 2)) / denom);
}

const KIND_LABELS: Record<BankCapitalComponentKind, string> = {
  1: "Cash (gold / notes)",
  2: "Bill of exchange",
  3: "Government bond",
  4: "Stock",
  5: "Mortgage",
};

export function Ch29BankCapital() {
  const [bankCapitals, setBankCapitals] = useState<BankCapital[]>([]);
  const [valuations, setValuations] = useState<FictitiousCapitalValuation[]>([]);
  const [listError, setListError] = useState<string | null>(null);

  // ── Fictitious-capital valuation calculator ──
  const [valName, setValName] = useState("Government Bond @5%");
  const [nominalValue, setNominalValue] = useState(10000);
  const [annualIncome, setAnnualIncome] = useState(500);
  const [interestRateBP, setInterestRateBP] = useState(500);
  const [dividendRateBP, setDividendRateBP] = useState(0);
  const [valDescription, setValDescription] = useState("");
  const [valError, setValError] = useState<string | null>(null);

  // ── Bank-capital form ──
  const [bankName, setBankName] = useState("Scottish Bank");
  const [cashAmount, setCashAmount] = useState(300_000_000);
  const [securitiesAmount, setSecuritiesAmount] = useState(2_400_000_000);
  const [components, setComponents] = useState<BankCapitalComponent[]>([
    { amount: 300_000_000, kind: 1, description: "gold and notes in till" },
    { amount: 1_400_000_000, kind: 2, description: "discounted commercial paper" },
    { amount: 1_000_000_000, kind: 3, description: "government consols" },
  ]);
  const [bankDescription, setBankDescription] = useState("");
  const [bankError, setBankError] = useState<string | null>(null);

  async function refresh() {
    try {
      const [bcs, vals] = await Promise.all([
        financeApi.listBankCapitals(),
        financeApi.listFictitiousCapitalValuations(),
      ]);
      setBankCapitals(bcs);
      setValuations(vals);
    } catch (e) {
      setListError(String(e));
    }
  }

  useEffect(() => {
    refresh().catch(() => {});
  }, []);

  // Live preview of the capitalisation (matches the server's derived market value).
  const previewMarketValue =
    interestRateBP > 0 ? roundHalfUp(annualIncome * 10000, interestRateBP) : null;

  async function handleCreateValuation() {
    if (interestRateBP <= 0) {
      setValError("Interest rate must be > 0.");
      return;
    }
    setValError(null);
    try {
      await financeApi.createFictitiousCapitalValuation({
        name: valName,
        nominal_value: nominalValue,
        annual_income: annualIncome,
        interest_rate_bp: interestRateBP,
        dividend_rate_bp: dividendRateBP,
        description: valDescription,
      });
      await refresh();
    } catch (e) {
      setValError(String(e));
    }
  }

  const computedTotal = cashAmount + securitiesAmount;

  function updateComponent(i: number, patch: Partial<BankCapitalComponent>) {
    setComponents((cs) => cs.map((c, idx) => (idx === i ? { ...c, ...patch } : c)));
  }

  function addComponent() {
    setComponents((cs) => [...cs, { amount: 0, kind: 1, description: "" }]);
  }

  function removeComponent(i: number) {
    setComponents((cs) => cs.filter((_, idx) => idx !== i));
  }

  async function handleCreateBankCapital() {
    setBankError(null);
    try {
      await financeApi.createBankCapital({
        name: bankName,
        cash_amount: cashAmount,
        securities_amount: securitiesAmount,
        components,
        description: bankDescription,
      });
      await refresh();
    } catch (e) {
      setBankError(String(e));
    }
  }

  return (
    <div className="v3-ch29">
      <section className="v3-ch29-intro">
        <p>
          Bank capital splits into two material parts: <strong>cash</strong> — the
          gold coin and bank-notes actually held in the till — and{" "}
          <strong>securities</strong>, which divide into commercial paper (bills of
          exchange) and public securities (government bonds, stocks, mortgages).
        </p>
        <p>
          The greater part of bank capital is <em>fictitious</em>. Deposits are mere
          book entries; bonds and stocks are titles to a future stream of
          surplus-value, not the underlying real capital. A security's market value
          is found by capitalising the income it yields at the prevailing rate of
          interest:
        </p>
        <p className="v3-ch29-formula">
          market value = annual income × 10000 / interest-rate-bp
        </p>
        <p>
          Because the rate of interest sits in the denominator, paper values move{" "}
          <strong>inversely</strong> with the rate: a fall in the rate inflates them,
          a rise deflates them, while the real capital (where it exists at all) is
          untouched.
        </p>
      </section>

      {/* Fictitious-capital valuation calculator */}
      <section className="v3-ch29-section" aria-label="Capitalise a security's income">
        <h4>Fictitious-capital calculator</h4>
        <p className="v3-ch29-subintro">
          Enter an annual income and the prevailing rate of interest to capitalise a
          security. £5 income at 5% is worth £100; halve the rate and the paper value
          doubles; double the rate and it halves.
        </p>
        <div className="v3-ch29-form">
          <label>
            Name
            <input value={valName} onChange={(e) => setValName(e.target.value)} />
          </label>
          <label>
            Nominal value (pence)
            <input
              type="number"
              min={0}
              value={nominalValue}
              onChange={(e) => setNominalValue(Number(e.target.value))}
            />
            <span className="v3-ch29-hint">{penceToPounds(nominalValue)} face / par</span>
          </label>
          <label>
            Annual income (pence)
            <input
              type="number"
              min={0}
              value={annualIncome}
              onChange={(e) => setAnnualIncome(Number(e.target.value))}
            />
            <span className="v3-ch29-hint">{penceToPounds(annualIncome)} yearly yield</span>
          </label>
          <label>
            Interest rate (bp, &gt; 0)
            <input
              type="number"
              min={1}
              value={interestRateBP}
              onChange={(e) => setInterestRateBP(Number(e.target.value))}
            />
            <span className="v3-ch29-hint">{bpToPercent(interestRateBP)} prevailing rate</span>
          </label>
          <label>
            Dividend rate (bp, optional — for stocks)
            <input
              type="number"
              min={0}
              value={dividendRateBP}
              onChange={(e) => setDividendRateBP(Number(e.target.value))}
            />
          </label>
          <label>
            Description
            <input value={valDescription} onChange={(e) => setValDescription(e.target.value)} />
          </label>
        </div>
        <div className="v3-ch29-result">
          Market value:{" "}
          <strong>
            {previewMarketValue === null ? "—" : penceToPounds(previewMarketValue)}
          </strong>
          {previewMarketValue !== null && (
            <span className="v3-ch29-result-note">
              {" "}
              = {penceToPounds(annualIncome)} × 10000 / {interestRateBP} bp
            </span>
          )}
        </div>
        {valError && <p className="v3-ch29-error">{valError}</p>}
        <button className="v3-ch29-btn" onClick={handleCreateValuation}>
          Save valuation
        </button>
      </section>

      {/* Inverse-relation illustration */}
      <section className="v3-ch29-section" aria-label="Inverse relation to the interest rate">
        <h4>Inverse relation to the rate of interest</h4>
        <p className="v3-ch29-subintro">
          The same £5 income capitalised at three rates — the paper value falls as the
          rate rises, though the underlying income never changes.
        </p>
        <div className="v3-ch29-inverse">
          {[250, 500, 1000].map((bp) => (
            <div className="v3-ch29-inverse-row" key={bp}>
              <span className="v3-ch29-inverse-rate">{bpToPercent(bp)}</span>
              <span className="v3-ch29-inverse-arrow">→</span>
              <span className="v3-ch29-inverse-val">{penceToPounds(roundHalfUp(500 * 10000, bp))}</span>
            </div>
          ))}
        </div>
      </section>

      {/* Seed exemplar card — Scottish bank */}
      <section className="v3-ch29-section" aria-label="Scottish bank exemplar">
        <h4>Marx's exemplar: a Scottish bank</h4>
        <div className="v3-ch29-seed-card">
          <div className="v3-ch29-seed-row">
            <span className="v3-ch29-seed-label">Cash (gold + notes)</span>
            <span className="v3-ch29-seed-val">£3,000,000 (300,000,000 pence)</span>
          </div>
          <div className="v3-ch29-seed-row">
            <span className="v3-ch29-seed-label">Securities (bills + public)</span>
            <span className="v3-ch29-seed-val">£24,000,000 (2,400,000,000 pence)</span>
          </div>
          <div className="v3-ch29-seed-row">
            <span className="v3-ch29-seed-label">Total bank capital</span>
            <span className="v3-ch29-seed-val">£27,000,000 (2,700,000,000 pence)</span>
          </div>
          <div className="v3-ch29-seed-verdict">
            Cash is barely a ninth of the whole — the greater part of bank capital is
            fictitious.
          </div>
        </div>
      </section>

      {/* Bank-capital form */}
      <section className="v3-ch29-section" aria-label="Record a bank capital decomposition">
        <h4>Record a bank capital</h4>
        <p className="v3-ch29-subintro">
          Split a bank's capital into cash and securities; the total is computed as
          cash + securities. Add component parts to break the holdings down by kind.
        </p>
        <div className="v3-ch29-form">
          <label>
            Name
            <input value={bankName} onChange={(e) => setBankName(e.target.value)} />
          </label>
          <label>
            Cash amount (pence)
            <input
              type="number"
              min={0}
              value={cashAmount}
              onChange={(e) => setCashAmount(Number(e.target.value))}
            />
            <span className="v3-ch29-hint">{penceToPounds(cashAmount)} gold + notes</span>
          </label>
          <label>
            Securities amount (pence)
            <input
              type="number"
              min={0}
              value={securitiesAmount}
              onChange={(e) => setSecuritiesAmount(Number(e.target.value))}
            />
            <span className="v3-ch29-hint">{penceToPounds(securitiesAmount)} bills + public securities</span>
          </label>
          <label>
            Description
            <input value={bankDescription} onChange={(e) => setBankDescription(e.target.value)} />
          </label>
        </div>

        <div className="v3-ch29-components">
          <div className="v3-ch29-components-head">
            <span>Component parts</span>
            <button className="v3-ch29-btn-small" onClick={addComponent}>
              + Add
            </button>
          </div>
          {components.map((c, i) => (
            <div className="v3-ch29-component-row" key={i}>
              <input
                className="v3-ch29-component-amount"
                type="number"
                min={0}
                value={c.amount}
                onChange={(e) => updateComponent(i, { amount: Number(e.target.value) })}
              />
              <select
                className="v3-ch29-component-kind"
                value={c.kind}
                onChange={(e) =>
                  updateComponent(i, { kind: Number(e.target.value) as BankCapitalComponentKind })
                }
              >
                {([1, 2, 3, 4, 5] as BankCapitalComponentKind[]).map((k) => (
                  <option key={k} value={k}>
                    {KIND_LABELS[k]}
                  </option>
                ))}
              </select>
              <input
                className="v3-ch29-component-desc"
                value={c.description}
                placeholder="description"
                onChange={(e) => updateComponent(i, { description: e.target.value })}
              />
              <button
                className="v3-ch29-btn-small v3-ch29-btn-remove"
                onClick={() => removeComponent(i)}
                aria-label="Remove component"
              >
                ×
              </button>
            </div>
          ))}
        </div>

        <div className="v3-ch29-total">
          Total capital: <strong>{penceToPounds(computedTotal)}</strong>
          <span className="v3-ch29-result-note"> = cash + securities</span>
        </div>
        {bankError && <p className="v3-ch29-error">{bankError}</p>}
        <button className="v3-ch29-btn" onClick={handleCreateBankCapital}>
          Save bank capital
        </button>
      </section>

      {/* Stored bank capitals */}
      <section className="v3-ch29-section" aria-label="Stored bank capitals">
        <h4>Stored bank capitals</h4>
        {listError && <p className="v3-ch29-error">{listError}</p>}
        <div className="v3-ch29-cards">
          {bankCapitals.length === 0 && (
            <p className="v3-ch29-subintro">No bank capitals stored yet.</p>
          )}
          {bankCapitals.map((bc) => (
            <div className="v3-ch29-card" key={bc.id}>
              <div className="v3-ch29-card-name">{bc.name}</div>
              <div className="v3-ch29-card-rows">
                <div className="v3-ch29-card-row">
                  <span>Cash</span>
                  <span>{penceToPounds(bc.cash_amount)}</span>
                </div>
                <div className="v3-ch29-card-row">
                  <span>Securities</span>
                  <span>{penceToPounds(bc.securities_amount)}</span>
                </div>
                <div className="v3-ch29-card-row v3-ch29-card-total">
                  <span>Total capital</span>
                  <span>{penceToPounds(bc.total_capital)}</span>
                </div>
              </div>
              {bc.components.length > 0 && (
                <ul className="v3-ch29-card-components">
                  {bc.components.map((c, i) => (
                    <li key={i}>
                      {KIND_LABELS[c.kind]}: {penceToPounds(c.amount)}
                      {c.description ? ` — ${c.description}` : ""}
                    </li>
                  ))}
                </ul>
              )}
              {bc.description && <div className="v3-ch29-desc">{bc.description}</div>}
            </div>
          ))}
        </div>
      </section>

      {/* Stored valuations */}
      <section className="v3-ch29-section" aria-label="Stored fictitious-capital valuations">
        <h4>Stored valuations</h4>
        <div className="v3-ch29-cards">
          {valuations.length === 0 && (
            <p className="v3-ch29-subintro">No valuations stored yet.</p>
          )}
          {valuations.map((v) => (
            <div className="v3-ch29-card" key={v.id}>
              <div className="v3-ch29-card-name">{v.name}</div>
              <div className="v3-ch29-card-rows">
                <div className="v3-ch29-card-row">
                  <span>Annual income</span>
                  <span>{penceToPounds(v.annual_income)}</span>
                </div>
                <div className="v3-ch29-card-row">
                  <span>Interest rate</span>
                  <span>{bpToPercent(v.interest_rate_bp)}</span>
                </div>
                {v.dividend_rate_bp > 0 && (
                  <div className="v3-ch29-card-row">
                    <span>Dividend rate</span>
                    <span>{bpToPercent(v.dividend_rate_bp)}</span>
                  </div>
                )}
                <div className="v3-ch29-card-row v3-ch29-card-total">
                  <span>Market value</span>
                  <span>{penceToPounds(v.market_value)}</span>
                </div>
              </div>
              {v.description && <div className="v3-ch29-desc">{v.description}</div>}
            </div>
          ))}
        </div>
      </section>
    </div>
  );
}
