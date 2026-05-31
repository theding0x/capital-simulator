import { useEffect, useState } from "react";
import { financeApi } from "../../api";
import type { ProfitDivision } from "../../types";
import "./Ch23ProfitOfEnterprise.css";

// CapitalOwnershipForm constants mirror credit.CapitalOwnershipForm.
const OWNERSHIP_FORMS: { value: number; label: string }[] = [
  { value: 1, label: "Owner-Operator (own capital, no external lender)" },
  { value: 2, label: "Separated (borrowed capital, distinct money-capitalist)" },
];

function bpToPercent(bp: number): string {
  return `${(bp / 100).toFixed(2)}%`;
}

function ownershipLabel(form: number): string {
  return OWNERSHIP_FORMS.find((f) => f.value === form)?.label ?? String(form);
}

function ownershipClass(form: number): string {
  if (form === 2) return "v3-ch23-ownership-tag v3-ch23-ownership-tag--separated";
  return "v3-ch23-ownership-tag";
}

// Seed fixtures from Capital Vol. III Ch. 23.
const SEED_EXEMPLARS = [
  {
    label: "Borrowed capital — interest 5%, enterprise 15%",
    total_profit_bp: 2000,
    interest_bp: 500,
    profit_of_enterprise_bp: 1500,
    ownership_form: 2,
  },
  {
    label: "Own capital — no external lender, enterprise 20%",
    total_profit_bp: 2000,
    interest_bp: 0,
    profit_of_enterprise_bp: 2000,
    ownership_form: 1,
  },
];

export function Ch23ProfitOfEnterprise() {
  const [items, setItems] = useState<ProfitDivision[]>([]);
  const [listError, setListError] = useState<string | null>(null);

  // Create-form state.
  const [totalBP, setTotalBP] = useState(2000);
  const [interestBP, setInterestBP] = useState(500);
  const [ownershipForm, setOwnershipForm] = useState(2);
  const [labourMinutes, setLabourMinutes] = useState(0);
  const [mystificationIndex, setMystificationIndex] = useState(5000);
  const [createError, setCreateError] = useState<string | null>(null);

  // Live-preview: profit of enterprise = total - interest (integer subtraction).
  const previewEnterprise =
    totalBP > 0 && interestBP >= 0 && interestBP < totalBP
      ? Math.floor(totalBP - interestBP)
      : null;

  async function refreshItems() {
    try {
      const data = await financeApi.listProfitDivisions();
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
      await financeApi.createProfitDivision({
        total_profit_bp: totalBP,
        interest_bp: interestBP,
        ownership_form: ownershipForm,
        wages_labour_minutes: labourMinutes,
        mystification_index: mystificationIndex,
      });
      await refreshItems();
    } catch (e) {
      setCreateError(String(e));
    }
  }

  return (
    <div className="v3-ch23">
      <section className="v3-ch23-intro">
        <p>
          The total profit of capital (= average profit) is divided into two qualitatively
          distinct categories: <strong>interest</strong> (paid to the money-capitalist as
          owner) and <strong>profit of enterprise</strong> (retained by the functioning
          capitalist as active user). This division appears as if interest is the natural
          fruit of capital-ownership — an ideological inversion that conceals the
          surplus-value relation.
        </p>
        <p>
          Invariant: <code>total_profit_bp = interest_bp + profit_of_enterprise_bp</code>
          . Interest is always a portion, strictly less than the whole.
        </p>
      </section>

      {/* Canonical Marx exemplars */}
      <section className="v3-ch23-section" aria-label="Profit division exemplars">
        <h4>Marx's exemplars (Ch. 23)</h4>
        <p className="v3-ch23-subintro">
          ProfitOfEnterprise = TotalProfit &minus; Interest (basis points).
        </p>
        <div className="v3-ch23-cards">
          {SEED_EXEMPLARS.map((ex) => (
            <div key={ex.label} className="v3-ch23-card">
              <div className="v3-ch23-card-label">{ex.label}</div>
              <div className="v3-ch23-meta">
                <span>Total profit: {bpToPercent(ex.total_profit_bp)}</span>
                <span>Interest: {bpToPercent(ex.interest_bp)}</span>
              </div>
              <div className="v3-ch23-enterprise">
                Enterprise: {bpToPercent(ex.profit_of_enterprise_bp)}
              </div>
              <div className="v3-ch23-meta" style={{ marginTop: "0.3rem" }}>
                <span className={ownershipClass(ex.ownership_form)}>
                  {ownershipLabel(ex.ownership_form)}
                </span>
              </div>
            </div>
          ))}
        </div>
      </section>

      {/* Record profit division form */}
      <section className="v3-ch23-section" aria-label="Record a profit division">
        <h4>Record a profit division</h4>
        <div className="v3-ch23-form">
          <label>
            Total profit (basis points)
            <input
              type="number"
              min={1}
              value={totalBP}
              onChange={(e) => setTotalBP(Number(e.target.value))}
            />
            <span className="v3-ch23-hint">{bpToPercent(totalBP)}</span>
          </label>
          <label>
            Interest (basis points)
            <input
              type="number"
              min={0}
              value={interestBP}
              onChange={(e) => setInterestBP(Number(e.target.value))}
            />
            <span className="v3-ch23-hint">{bpToPercent(interestBP)}</span>
          </label>
          <label>
            Capital ownership form
            <select
              value={ownershipForm}
              onChange={(e) => setOwnershipForm(Number(e.target.value))}
            >
              {OWNERSHIP_FORMS.map((f) => (
                <option key={f.value} value={f.value}>
                  {f.label}
                </option>
              ))}
            </select>
          </label>
          <label>
            Wages of superintendence (labour minutes)
            <input
              type="number"
              min={0}
              value={labourMinutes}
              onChange={(e) => setLabourMinutes(Number(e.target.value))}
            />
          </label>
          <label>
            Mystification index (basis points, 0–10000)
            <input
              type="number"
              min={0}
              max={10000}
              value={mystificationIndex}
              onChange={(e) => setMystificationIndex(Number(e.target.value))}
            />
            <span className="v3-ch23-hint">{bpToPercent(mystificationIndex)}</span>
          </label>
        </div>

        {/* Live preview */}
        {previewEnterprise !== null && (
          <div className="v3-ch23-preview">
            <div className="v3-ch23-preview-row">
              <span className="v3-ch23-preview-label">Total profit</span>
              <span className="v3-ch23-preview-value">{bpToPercent(totalBP)}</span>
            </div>
            <div className="v3-ch23-preview-row">
              <span className="v3-ch23-preview-label">Interest</span>
              <span className="v3-ch23-preview-value">{bpToPercent(interestBP)}</span>
            </div>
            <div className="v3-ch23-preview-row">
              <span className="v3-ch23-preview-label">Profit of enterprise</span>
              <span className="v3-ch23-preview-value v3-ch23-preview-enterprise">
                {bpToPercent(previewEnterprise)}
              </span>
            </div>
          </div>
        )}

        <button
          type="button"
          className="v3-ch23-btn"
          disabled={totalBP <= 0 || interestBP < 0 || interestBP >= totalBP}
          onClick={handleCreate}
        >
          Record division
        </button>
        {createError && <p className="v3-ch23-error">{createError}</p>}
      </section>

      {/* Stored records */}
      {listError && <p className="v3-ch23-error">{listError}</p>}
      {items.length > 0 && (
        <section className="v3-ch23-section" aria-label="Stored profit divisions">
          <h4>Recorded profit divisions ({items.length})</h4>
          <div className="v3-ch23-cards">
            {items.map((pd) => (
              <div key={pd.id} className="v3-ch23-card v3-ch23-card--stored">
                <div className="v3-ch23-card-label">{pd.id}</div>
                <div className="v3-ch23-meta">
                  <span>Total: {bpToPercent(pd.total_profit_bp)}</span>
                  <span>Interest: {bpToPercent(pd.interest_bp)}</span>
                  <span>
                    Enterprise:{" "}
                    {bpToPercent(
                      Math.floor(pd.total_profit_bp - pd.interest_bp),
                    )}
                  </span>
                </div>
                <div className="v3-ch23-meta" style={{ marginTop: "0.3rem" }}>
                  <span className={ownershipClass(pd.ownership_form)}>
                    {ownershipLabel(pd.ownership_form)}
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
