import { useEffect, useState } from "react";
import { financeApi } from "../../api";
import type { CommercialProfitResponse } from "../../types";
import "./Ch17CommercialProfit.css";

function pence(n: number): string {
  const pounds = n / 100;
  return `£${pounds.toLocaleString("en-GB", { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`;
}

function bp(n: number): string {
  return `${(n / 100).toFixed(2)}%`;
}

// Marx Ch. 17 fixture: productive=£7,200 (720000p), commercial=£1,800 (180000p),
// surplus=£1,440 (144000p) → unadjusted 20%, adjusted 16%.
const FIXTURE = {
  productive: 720000,
  commercial: 180000,
  surplus: 144000,
  unadjusted_bp: 2000,
  adjusted_bp: 1600,
};

// Commercial clerk fixture: necessary 14,400 min / workday 28,800 min → unpaid 14,400 min.
const CLERK = {
  necessary: 14400,
  workday: 28800,
  unpaid: 14400,
};

export function Ch17CommercialProfit() {
  const [items, setItems] = useState<CommercialProfitResponse[]>([]);
  const [listError, setListError] = useState<string | null>(null);

  // Create-form state.
  const [productiveCapital, setProductiveCapital] = useState(720000);
  const [commercialCapital, setCommercialCapital] = useState(180000);
  const [surplusValue, setSurplusValue] = useState(144000);
  const [createError, setCreateError] = useState<string | null>(null);

  async function refreshItems() {
    try {
      const data = await financeApi.listCommercialProfits();
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
      await financeApi.createCommercialProfit({
        productive_capital: productiveCapital,
        commercial_capital: commercialCapital,
        surplus_value: surplusValue,
      });
      await refreshItems();
    } catch (e) {
      setCreateError(String(e));
    }
  }

  // Client-side preview of both rates (mirrors the Go round-half-up idiom).
  function previewRate(sv: number, denom: number): number {
    if (denom <= 0) return 0;
    return Math.round((sv * 10000 + denom / 2) / denom);
  }
  const previewUnadj = previewRate(surplusValue, productiveCapital);
  const previewAdj = previewRate(surplusValue, productiveCapital + commercialCapital);

  return (
    <div className="v3-ch17">
      <section className="v3-ch17-intro">
        <p>
          Commercial profit is not new value — it is a portion of the surplus-value
          transferred from the industrial circuit. The merchant advances{" "}
          <strong>productive capital</strong> (borrowed from the industrial sphere)
          alongside their own <strong>commercial capital</strong>. Once both are
          included in the denominator, the general rate of profit falls:{" "}
          <em>adjusted p′ &lt; unadjusted p′</em> whenever commercial capital &gt; 0.
        </p>
      </section>

      {/* Marx Ch. 17 fixture */}
      <section className="v3-ch17-section" aria-label="Marx Ch. 17 fixture">
        <h4>Marx Ch. 17 fixture</h4>
        <p className="v3-ch17-subintro">
          Productive capital {pence(FIXTURE.productive)} · Commercial capital{" "}
          {pence(FIXTURE.commercial)} · Surplus-value {pence(FIXTURE.surplus)}
        </p>
        <div className="v3-ch17-fixture-grid">
          <div className="v3-ch17-fixture-cell">
            <span className="v3-ch17-fixture-label">Unadjusted p′</span>
            <span className="v3-ch17-fixture-value">{bp(FIXTURE.unadjusted_bp)}</span>
          </div>
          <div className="v3-ch17-fixture-cell v3-ch17-fixture-cell--adj">
            <span className="v3-ch17-fixture-label">Adjusted p′</span>
            <span className="v3-ch17-fixture-value">{bp(FIXTURE.adjusted_bp)}</span>
          </div>
          <div className="v3-ch17-fixture-cell v3-ch17-fixture-cell--gap">
            <span className="v3-ch17-fixture-label">Rate dilution</span>
            <span className="v3-ch17-fixture-value">
              −{bp(FIXTURE.unadjusted_bp - FIXTURE.adjusted_bp)}
            </span>
          </div>
        </div>
      </section>

      {/* Commercial clerk — unpaid labour panel */}
      <section className="v3-ch17-section" aria-label="Commercial clerk">
        <h4>Commercial clerk — unpaid labour creates no new value</h4>
        <p className="v3-ch17-subintro">
          The commercial wage-worker performs unpaid labour for the merchant, enabling
          a larger profit share. But no new value is produced: NewValueCreated = 0.
        </p>
        <div className="v3-ch17-fixture-grid">
          <div className="v3-ch17-fixture-cell">
            <span className="v3-ch17-fixture-label">Necessary labour</span>
            <span className="v3-ch17-fixture-value">{CLERK.necessary.toLocaleString("en-GB")} min</span>
          </div>
          <div className="v3-ch17-fixture-cell">
            <span className="v3-ch17-fixture-label">Workday</span>
            <span className="v3-ch17-fixture-value">{CLERK.workday.toLocaleString("en-GB")} min</span>
          </div>
          <div className="v3-ch17-fixture-cell v3-ch17-fixture-cell--adj">
            <span className="v3-ch17-fixture-label">Unpaid labour</span>
            <span className="v3-ch17-fixture-value">{CLERK.unpaid.toLocaleString("en-GB")} min</span>
          </div>
          <div className="v3-ch17-fixture-cell v3-ch17-fixture-cell--zero">
            <span className="v3-ch17-fixture-label">New value created</span>
            <span className="v3-ch17-fixture-value">0</span>
          </div>
        </div>
      </section>

      {/* Create commercial profit */}
      <section className="v3-ch17-section" aria-label="Record commercial profit">
        <h4>Record a commercial-profit computation</h4>
        <div className="v3-ch17-form">
          <label>
            Productive capital (pence)
            <input
              type="number"
              min={1}
              value={productiveCapital}
              onChange={(e) => setProductiveCapital(Number(e.target.value))}
            />
          </label>
          <label>
            Commercial capital (pence)
            <input
              type="number"
              min={1}
              value={commercialCapital}
              onChange={(e) => setCommercialCapital(Number(e.target.value))}
            />
          </label>
          <label>
            Surplus-value (pence)
            <input
              type="number"
              min={0}
              value={surplusValue}
              onChange={(e) => setSurplusValue(Number(e.target.value))}
            />
          </label>
          <div className="v3-ch17-preview-pair">
            <div className="v3-ch17-preview">
              <span className="v3-ch17-preview-label">Unadjusted p′</span>
              <span className="v3-ch17-preview-figure">{bp(previewUnadj)}</span>
            </div>
            <div className="v3-ch17-preview v3-ch17-preview--adj">
              <span className="v3-ch17-preview-label">Adjusted p′</span>
              <span className="v3-ch17-preview-figure">{bp(previewAdj)}</span>
            </div>
          </div>
          <button type="button" className="v3-ch17-btn" onClick={handleCreate}>
            Record
          </button>
        </div>
        {createError && <p className="v3-ch17-error">{createError}</p>}
      </section>

      {/* List */}
      {listError && <p className="v3-ch17-error">{listError}</p>}
      {items.length > 0 && (
        <section className="v3-ch17-section" aria-label="Commercial profit records">
          <h4>Commercial-profit records ({items.length})</h4>
          <table className="v3-ch17-table">
            <thead>
              <tr>
                <th>Productive C</th>
                <th>Commercial C</th>
                <th>Surplus-value</th>
                <th>Unadj. p′</th>
                <th>Adj. p′</th>
                <th>New value</th>
              </tr>
            </thead>
            <tbody>
              {items.map((cp) => (
                <tr key={cp.id}>
                  <td className="v3-ch17-money">{pence(cp.productive_capital)}</td>
                  <td className="v3-ch17-money">{pence(cp.commercial_capital)}</td>
                  <td className="v3-ch17-money">{pence(cp.surplus_value)}</td>
                  <td className="v3-ch17-rate">{bp(cp.unadjusted_rate_bp)}</td>
                  <td className="v3-ch17-rate v3-ch17-rate--adj">{bp(cp.adjusted_rate_bp)}</td>
                  <td className="v3-ch17-zero">{cp.new_value_created}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </section>
      )}
    </div>
  );
}
