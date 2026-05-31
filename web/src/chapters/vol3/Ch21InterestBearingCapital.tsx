import { useEffect, useState } from "react";
import { financeApi } from "../../api";
import type { InterestBearingCapital } from "../../types";
import "./Ch21InterestBearingCapital.css";

function pence(p: number): string {
  return `£${(p / 100).toFixed(2)}`;
}

function bpToPercent(bp: number): string {
  return `${(bp / 100).toFixed(2)}%`;
}

// Seed fixtures illustrating the M—M' circuit (Marx, Vol. III Ch. 21).
const SEED_EXEMPLARS = [
  {
    label: "£1,000 at 5% — canonical Marx fixture",
    money_advanced: 100000,
    interest_rate_bp: 500,
    loan_term_days: 365,
    interest_earned: 5000,
    money_returned: 105000,
  },
  {
    label: "£5,000 at 15% for 180 days",
    money_advanced: 500000,
    interest_rate_bp: 1500,
    loan_term_days: 180,
    interest_earned: 75000,
    money_returned: 575000,
  },
];

export function Ch21InterestBearingCapital() {
  const [items, setItems] = useState<InterestBearingCapital[]>([]);
  const [listError, setListError] = useState<string | null>(null);

  // Create-form state.
  const [moneyAdvanced, setMoneyAdvanced] = useState(100000);
  const [interestRateBP, setInterestRateBP] = useState(500);
  const [loanTermDays, setLoanTermDays] = useState(365);
  const [createError, setCreateError] = useState<string | null>(null);

  // Live-preview derived values.
  const previewInterest = Math.round((moneyAdvanced * interestRateBP + 5000) / 10000);
  const previewReturned = moneyAdvanced + previewInterest;

  async function refreshItems() {
    try {
      const data = await financeApi.listInterestBearingCapitals();
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
      await financeApi.createInterestBearingCapital({
        money_advanced: moneyAdvanced,
        interest_rate_bp: interestRateBP,
        loan_term_days: loanTermDays,
      });
      await refreshItems();
    } catch (e) {
      setCreateError(String(e));
    }
  }

  return (
    <div className="v3-ch21">
      <section className="v3-ch21-intro">
        <p>
          Interest-bearing capital is money-capital advanced to a functioning capitalist and
          returned with interest: the <strong>M—M′</strong> circuit. Interest is a portion of
          the average profit earned on the borrowed capital, paid by the borrower (functioning
          capitalist) to the lender (money capitalist). The money-capitalist creates{" "}
          <em>no new value</em> — interest is a redistribution of already-produced
          surplus-value. The functioning capitalist retains the{" "}
          <strong>profit of enterprise</strong> (average profit minus interest).
        </p>
      </section>

      {/* Canonical Marx exemplars */}
      <section className="v3-ch21-section" aria-label="M—M′ circuit exemplars">
        <h4>M—M′ circuit — Marx's exemplars</h4>
        <p className="v3-ch21-subintro">
          InterestEarned = roundHalfUp(MoneyAdvanced &times; InterestRateBP, 10000).
          NewValueCreated = 0 always.
        </p>
        <div className="v3-ch21-cards">
          {SEED_EXEMPLARS.map((ex) => (
            <div key={ex.label} className="v3-ch21-card">
              <div className="v3-ch21-card-label">{ex.label}</div>
              <div className="v3-ch21-circuit">
                <span className="v3-ch21-node v3-ch21-node--m">
                  M {pence(ex.money_advanced)}
                </span>
                <span className="v3-ch21-arrow">&#8594;</span>
                <span className="v3-ch21-node v3-ch21-node--mprime">
                  M′ {pence(ex.money_returned)}
                </span>
              </div>
              <div className="v3-ch21-meta">
                <span>Rate: {bpToPercent(ex.interest_rate_bp)}</span>
                <span>Term: {ex.loan_term_days} days</span>
                <span>Interest: {pence(ex.interest_earned)}</span>
              </div>
            </div>
          ))}
        </div>
      </section>

      {/* Create form */}
      <section className="v3-ch21-section" aria-label="Record an interest-bearing loan">
        <h4>Record an interest-bearing loan</h4>
        <div className="v3-ch21-form">
          <label>
            Money advanced (pence)
            <input
              type="number"
              min={1}
              value={moneyAdvanced}
              onChange={(e) => setMoneyAdvanced(Number(e.target.value))}
            />
            <span className="v3-ch21-hint">{pence(moneyAdvanced)}</span>
          </label>
          <label>
            Interest rate (basis points)
            <input
              type="number"
              min={1}
              value={interestRateBP}
              onChange={(e) => setInterestRateBP(Number(e.target.value))}
            />
            <span className="v3-ch21-hint">{bpToPercent(interestRateBP)}</span>
          </label>
          <label>
            Loan term (days)
            <input
              type="number"
              min={1}
              value={loanTermDays}
              onChange={(e) => setLoanTermDays(Number(e.target.value))}
            />
          </label>
        </div>

        {/* Live preview */}
        <div className="v3-ch21-preview">
          <div className="v3-ch21-preview-row">
            <span className="v3-ch21-preview-label">Interest earned</span>
            <span className="v3-ch21-preview-value">{pence(previewInterest)}</span>
          </div>
          <div className="v3-ch21-preview-row">
            <span className="v3-ch21-preview-label">Money returned (M′)</span>
            <span className="v3-ch21-preview-value">{pence(previewReturned)}</span>
          </div>
          <div className="v3-ch21-preview-row">
            <span className="v3-ch21-preview-label">New value created</span>
            <span className="v3-ch21-preview-value v3-ch21-preview-zero">£0.00</span>
          </div>
        </div>

        <button
          type="button"
          className="v3-ch21-btn"
          disabled={moneyAdvanced <= 0 || interestRateBP <= 0 || loanTermDays <= 0}
          onClick={handleCreate}
        >
          Record loan
        </button>
        {createError && <p className="v3-ch21-error">{createError}</p>}
      </section>

      {/* List of stored loans */}
      {listError && <p className="v3-ch21-error">{listError}</p>}
      {items.length > 0 && (
        <section className="v3-ch21-section" aria-label="Stored interest-bearing loans">
          <h4>Recorded loans ({items.length})</h4>
          <div className="v3-ch21-cards">
            {items.map((ibc) => (
              <div key={ibc.id} className="v3-ch21-card v3-ch21-card--stored">
                <div className="v3-ch21-circuit">
                  <span className="v3-ch21-node v3-ch21-node--m">
                    M {pence(ibc.money_advanced)}
                  </span>
                  <span className="v3-ch21-arrow">&#8594;</span>
                  <span className="v3-ch21-node v3-ch21-node--mprime">
                    M′ {pence(ibc.money_returned)}
                  </span>
                </div>
                <div className="v3-ch21-meta">
                  <span>Rate: {bpToPercent(ibc.interest_rate_bp)}</span>
                  <span>Term: {ibc.loan_term_days} days</span>
                  <span>Interest: {pence(ibc.interest_earned)}</span>
                </div>
                <div className="v3-ch21-no-value">
                  New value created: {pence(ibc.new_value_created)}
                </div>
              </div>
            ))}
          </div>
        </section>
      )}
    </div>
  );
}
