import { useEffect, useState } from "react";
import { financeApi } from "../../api";
import type { MoneyDealingCapitalResponse, MoneyOperationKind } from "../../types";
import "./Ch19MoneyDealingCapital.css";

function pence(n: number): string {
  const pounds = n / 100;
  return `£${pounds.toLocaleString("en-GB", { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`;
}

function bp(n: number): string {
  return `${(n / 100).toFixed(2)}%`;
}

// Marx Ch. 19 fixture: capital=£5,000 (500000p), general rate=15% (1500 bp).
// profit = roundHalfUp(500000*1500, 10000) = £750 (75000p).
// NewValueCreated = 0. Operations: Receipt, Payment, Exchange (kinds 1,2,3).
const FIXTURE = {
  money_advanced: 500000,
  general_rate_bp: 1500,
  money_dealing_profit: 75000,
  new_value_created: 0,
};

const KIND_LABELS: Record<number, string> = {
  1: "Receipt",
  2: "Payment",
  3: "Exchange",
  4: "Safekeeping",
  5: "Bookkeeping",
};

const ALL_KINDS: MoneyOperationKind[] = [1, 2, 3, 4, 5];

// Client-side round-half-up preview of money_dealing_profit.
function previewProfit(money: number, rateBP: number): number {
  if (money <= 0 || rateBP <= 0) return 0;
  return Math.floor((money * rateBP + 5000) / 10000);
}

function kindLabel(k: MoneyOperationKind): string {
  return KIND_LABELS[k] ?? String(k);
}

function kindsList(kinds: MoneyOperationKind[]): string {
  if (kinds.length === 0) return "—";
  return kinds.map((k) => kindLabel(k)).join(", ");
}

export function Ch19MoneyDealingCapital() {
  const [items, setItems] = useState<MoneyDealingCapitalResponse[]>([]);
  const [listError, setListError] = useState<string | null>(null);

  // Create-form state.
  const [moneyAdvanced, setMoneyAdvanced] = useState(500000);
  const [generalRateBP, setGeneralRateBP] = useState(1500);
  const [selectedKinds, setSelectedKinds] = useState<Set<MoneyOperationKind>>(
    new Set([1, 2, 3] as MoneyOperationKind[]),
  );
  const [createError, setCreateError] = useState<string | null>(null);

  async function refreshItems() {
    try {
      const data = await financeApi.listMoneyDealingCapitals();
      setItems(data);
    } catch (e) {
      setListError(String(e));
    }
  }

  useEffect(() => {
    refreshItems().catch(() => {});
  }, []);

  function toggleKind(k: MoneyOperationKind) {
    setSelectedKinds((prev) => {
      const next = new Set(prev);
      if (next.has(k)) {
        next.delete(k);
      } else {
        next.add(k);
      }
      return next;
    });
  }

  async function handleCreate() {
    setCreateError(null);
    const kinds = ALL_KINDS.filter((k) => selectedKinds.has(k));
    try {
      await financeApi.createMoneyDealingCapital({
        money_advanced: moneyAdvanced,
        general_rate_bp: generalRateBP,
        operation_kinds: kinds,
      });
      await refreshItems();
    } catch (e) {
      setCreateError(String(e));
    }
  }

  const profitPreview = previewProfit(moneyAdvanced, generalRateBP);

  return (
    <div className="v3-ch19">
      <section className="v3-ch19-intro">
        <p>
          Money-dealing capital specialises in the purely technical operations of monetary
          circulation — receipt, payment, currency exchange, safekeeping, and bookkeeping — that
          arise necessarily from the circulation of commodities and wages. Like commercial capital
          it participates in forming the general rate of profit, earning{" "}
          <strong>MoneyDealingProfit = roundHalfUp(capital × rate)</strong>. But the M—M′ circuit
          creates <strong>no new value</strong>: <em>NewValueCreated == 0</em> always.
        </p>
      </section>

      {/* Marx Ch. 19 fixture */}
      <section className="v3-ch19-section" aria-label="Marx Ch. 19 fixture">
        <h4>Marx Ch. 19 fixture — M—M′, no value created</h4>
        <p className="v3-ch19-subintro">
          Capital {pence(FIXTURE.money_advanced)} · General rate{" "}
          {bp(FIXTURE.general_rate_bp)} · Operations: Receipt, Payment, Exchange
        </p>
        <div className="v3-ch19-fixture-grid">
          <div className="v3-ch19-fixture-cell">
            <span className="v3-ch19-fixture-label">Money advanced (M)</span>
            <span className="v3-ch19-fixture-value">{pence(FIXTURE.money_advanced)}</span>
          </div>
          <div className="v3-ch19-fixture-cell v3-ch19-fixture-cell--profit">
            <span className="v3-ch19-fixture-label">Money-dealing profit</span>
            <span className="v3-ch19-fixture-value">{pence(FIXTURE.money_dealing_profit)}</span>
          </div>
          <div className="v3-ch19-fixture-cell v3-ch19-fixture-cell--zero">
            <span className="v3-ch19-fixture-label">New value created</span>
            <span className="v3-ch19-fixture-value">£0.00</span>
          </div>
        </div>
      </section>

      {/* Create form */}
      <section className="v3-ch19-section" aria-label="Record money-dealing capital">
        <h4>Record a money-dealing capital advance</h4>
        <div className="v3-ch19-form">
          <label>
            Money advanced (pence)
            <input
              type="number"
              min={1}
              value={moneyAdvanced}
              onChange={(e) => setMoneyAdvanced(Number(e.target.value))}
            />
          </label>
          <label>
            General rate (bp)
            <input
              type="number"
              min={1}
              value={generalRateBP}
              onChange={(e) => setGeneralRateBP(Number(e.target.value))}
            />
          </label>
          <fieldset className="v3-ch19-kinds-fieldset">
            <legend>Operation kinds</legend>
            <div className="v3-ch19-kinds-grid">
              {ALL_KINDS.map((k) => (
                <label key={k} className="v3-ch19-kind-label">
                  <input
                    type="checkbox"
                    checked={selectedKinds.has(k)}
                    onChange={() => toggleKind(k)}
                  />
                  {kindLabel(k)}
                </label>
              ))}
            </div>
          </fieldset>
          <div className="v3-ch19-preview-pair">
            <div className="v3-ch19-preview v3-ch19-preview--profit">
              <span className="v3-ch19-preview-label">Profit (M′ − M)</span>
              <span className="v3-ch19-preview-figure">{pence(profitPreview)}</span>
            </div>
            <div className="v3-ch19-preview v3-ch19-preview--zero">
              <span className="v3-ch19-preview-label">New value created</span>
              <span className="v3-ch19-preview-figure">£0.00</span>
            </div>
          </div>
          <button type="button" className="v3-ch19-btn" onClick={handleCreate}>
            Record
          </button>
        </div>
        {createError && <p className="v3-ch19-error">{createError}</p>}
      </section>

      {/* List */}
      {listError && <p className="v3-ch19-error">{listError}</p>}
      {items.length > 0 && (
        <section className="v3-ch19-section" aria-label="Money-dealing capital records">
          <h4>Money-dealing capital records ({items.length})</h4>
          <table className="v3-ch19-table">
            <thead>
              <tr>
                <th>Money advanced</th>
                <th>Rate</th>
                <th>Operations</th>
                <th>Profit (M′ − M)</th>
                <th>New value</th>
              </tr>
            </thead>
            <tbody>
              {items.map((md) => (
                <tr key={md.id}>
                  <td className="v3-ch19-money">{pence(md.money_advanced)}</td>
                  <td className="v3-ch19-rate">{bp(md.general_rate_bp)}</td>
                  <td className="v3-ch19-ops">{kindsList(md.operation_kinds)}</td>
                  <td className="v3-ch19-money v3-ch19-profit">{pence(md.money_dealing_profit)}</td>
                  <td className="v3-ch19-zero">£0.00</td>
                </tr>
              ))}
            </tbody>
          </table>
        </section>
      )}
    </div>
  );
}
