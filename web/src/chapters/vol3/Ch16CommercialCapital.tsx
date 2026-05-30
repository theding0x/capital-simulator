import { useEffect, useState } from "react";
import { financeApi } from "../../api";
import type { CommercialCapitalFunction, CommercialCapitalResponse } from "../../types";
import "./Ch16CommercialCapital.css";

function pence(n: number): string {
  const pounds = n / 100;
  return `£${pounds.toLocaleString("en-GB", { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`;
}

function num(n: number): string {
  return n.toLocaleString("en-GB");
}

interface FunctionMeta {
  value: CommercialCapitalFunction;
  label: string;
  description: string;
}

// The four functions of commercial capital — §1–§4. None produces surplus-value.
const FUNCTIONS: FunctionMeta[] = [
  {
    value: 0,
    label: "Realisation",
    description:
      "§1: The merchant converts C′ into M′ for the industrial capitalist, realising surplus-value already produced. M—C—M′.",
  },
  {
    value: 1,
    label: "Storage",
    description:
      "§2: The merchant holds stock, bearing warehousing costs while awaiting a favourable price. No surplus-value produced.",
  },
  {
    value: 2,
    label: "Transport",
    description:
      "§3: The merchant moves commodities from production to consumption. Use-value is added for the buyer; no surplus-value for the merchant.",
  },
  {
    value: 3,
    label: "Distribution",
    description:
      "§4: Breaking bulk, re-packaging and allied activities. Again, no surplus-value is generated.",
  },
];

const FUNCTION_LABELS: Record<CommercialCapitalFunction, string> = {
  0: "Realisation",
  1: "Storage",
  2: "Transport",
  3: "Distribution",
};

// Marx §1 linen-merchant fixture: £3,000 advanced, 30,000 yards at 10p/yard.
const LINEN_MERCHANT_MONEY_ADVANCED = 300000;
const LINEN_MERCHANT_YARDS = 30000;
const LINEN_MERCHANT_PRICE_PER_YARD = 10;

// merchantProfit computes the profit on a M—C—M′ circuit (client-side preview).
function merchantProfit(moneyAdvanced: number, moneyReturned: number): number {
  return moneyReturned - moneyAdvanced;
}

export function Ch16CommercialCapital() {
  const [items, setItems] = useState<CommercialCapitalResponse[]>([]);
  const [listError, setListError] = useState<string | null>(null);

  // Create-form state.
  const [moneyAdvanced, setMoneyAdvanced] = useState(300000);
  const [commodityDescription, setCommodityDescription] = useState(
    "Linen merchant: £3,000 advanced, 30,000 yards at 10p/yard",
  );
  const [fn, setFn] = useState<CommercialCapitalFunction>(0);
  const [createError, setCreateError] = useState<string | null>(null);

  // MerchantCircuit profit calculator (client-side; value object not persisted).
  const [circuitAdvanced, setCircuitAdvanced] = useState(300000);
  const [circuitValue, setCircuitValue] = useState(300000);
  const [circuitReturned, setCircuitReturned] = useState(330000);

  async function refreshItems() {
    try {
      const data = await financeApi.listCommercialCapitals();
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
      await financeApi.createCommercialCapital({
        money_advanced: moneyAdvanced,
        commodity_description: commodityDescription,
        function: fn,
      });
      await refreshItems();
    } catch (e) {
      setCreateError(String(e));
    }
  }

  const previewProfit = merchantProfit(circuitAdvanced, circuitReturned);
  const linenTotal = LINEN_MERCHANT_YARDS * LINEN_MERCHANT_PRICE_PER_YARD;

  return (
    <div className="v3-ch16">
      <section className="v3-ch16-intro">
        <p>
          Commercial capital is money-capital and commodity-capital separated off
          from the industrial circuit and made independent as a special form.{" "}
          <strong>It does not produce surplus-value</strong> — it appropriates a
          share of the surplus-value already produced in the industrial process by
          performing the function of realisation (C′→M′). The linen merchant
          advances £3,000, buys 30,000 yards, and sells them: the profit is a
          deduction from the industrial surplus-value, not a creation of new value.
        </p>
      </section>

      {/* Marx §1 fixture */}
      <section className="v3-ch16-section" aria-label="Linen merchant fixture">
        <h4>Marx §1 fixture — the linen merchant</h4>
        <p className="v3-ch16-subintro">
          £3,000 advanced · 30,000 yards at 10p/yard · total value{" "}
          <code>{pence(linenTotal)}</code> = <code>{num(linenTotal)} pence</code>.
          SurplusValueProduced = 0 in all cases.
        </p>
        <div className="v3-ch16-fixture-grid">
          <div className="v3-ch16-fixture-cell">
            <span className="v3-ch16-fixture-label">Money advanced</span>
            <span className="v3-ch16-fixture-value">{pence(LINEN_MERCHANT_MONEY_ADVANCED)}</span>
          </div>
          <div className="v3-ch16-fixture-cell">
            <span className="v3-ch16-fixture-label">Yards bought</span>
            <span className="v3-ch16-fixture-value">{num(LINEN_MERCHANT_YARDS)}</span>
          </div>
          <div className="v3-ch16-fixture-cell">
            <span className="v3-ch16-fixture-label">Price / yard</span>
            <span className="v3-ch16-fixture-value">{LINEN_MERCHANT_PRICE_PER_YARD}p</span>
          </div>
          <div className="v3-ch16-fixture-cell">
            <span className="v3-ch16-fixture-label">Total value</span>
            <span className="v3-ch16-fixture-value">{pence(linenTotal)}</span>
          </div>
        </div>
      </section>

      {/* Four functions taxonomy */}
      <section className="v3-ch16-section" aria-label="Four functions of commercial capital">
        <h4>The four functions of commercial capital (§1–§4)</h4>
        <p className="v3-ch16-subintro">
          All four produce <strong>zero</strong> surplus-value. The merchant's
          profit is always a redistribution of industrial surplus-value.
        </p>
        <ul className="v3-ch16-functions">
          {FUNCTIONS.map((f) => (
            <li key={f.value} className="v3-ch16-function">
              <span className="v3-ch16-function-label">
                §{f.value + 1} {f.label}
                <span className="v3-ch16-function-zero"> · s=0</span>
              </span>
              <span className="v3-ch16-function-desc">{f.description}</span>
            </li>
          ))}
        </ul>
      </section>

      {/* MerchantCircuit profit calculator */}
      <section className="v3-ch16-section" aria-label="Merchant circuit calculator">
        <h4>M—C—M′ profit calculator</h4>
        <p className="v3-ch16-subintro">
          Profit = M′ − M. The excess is a deduction from the industrial
          surplus-value, not new value created by the merchant.
        </p>
        <div className="v3-ch16-form">
          <label>
            Money advanced M (pence)
            <input
              type="number"
              min={1}
              value={circuitAdvanced}
              onChange={(e) => setCircuitAdvanced(Number(e.target.value))}
            />
          </label>
          <label>
            Commodity value C (pence)
            <input
              type="number"
              min={0}
              value={circuitValue}
              onChange={(e) => setCircuitValue(Number(e.target.value))}
            />
          </label>
          <label>
            Money returned M′ (pence)
            <input
              type="number"
              min={0}
              value={circuitReturned}
              onChange={(e) => setCircuitReturned(Number(e.target.value))}
            />
          </label>
          <div className="v3-ch16-preview">
            <span className="v3-ch16-preview-label">Merchant profit M′−M</span>
            <span className="v3-ch16-preview-figure">{pence(previewProfit)}</span>
          </div>
        </div>
      </section>

      {/* Create commercial capital */}
      <section className="v3-ch16-section" aria-label="Record commercial capital advance">
        <h4>Record a commercial-capital advance</h4>
        <div className="v3-ch16-form">
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
            Commodity description
            <input
              type="text"
              value={commodityDescription}
              onChange={(e) => setCommodityDescription(e.target.value)}
              className="v3-ch16-text-input"
            />
          </label>
          <label>
            Function
            <select
              value={fn}
              onChange={(e) => setFn(Number(e.target.value) as CommercialCapitalFunction)}
            >
              {FUNCTIONS.map((f) => (
                <option key={f.value} value={f.value}>
                  {f.label}
                </option>
              ))}
            </select>
          </label>
          <button type="button" className="v3-ch16-btn" onClick={handleCreate}>
            Record
          </button>
        </div>
        {createError && <p className="v3-ch16-error">{createError}</p>}
      </section>

      {/* List */}
      {listError && <p className="v3-ch16-error">{listError}</p>}
      {items.length > 0 && (
        <section className="v3-ch16-section" aria-label="Commercial capital records">
          <h4>Commercial capital advances ({items.length})</h4>
          <table className="v3-ch16-table">
            <thead>
              <tr>
                <th>Money advanced</th>
                <th>Commodity</th>
                <th>Function</th>
                <th>s produced</th>
              </tr>
            </thead>
            <tbody>
              {items.map((cc) => (
                <tr key={cc.id}>
                  <td className="v3-ch16-money">{pence(cc.money_advanced)}</td>
                  <td>{cc.commodity_description}</td>
                  <td>{FUNCTION_LABELS[cc.function]}</td>
                  <td className="v3-ch16-zero">{cc.surplus_value_produced}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </section>
      )}
    </div>
  );
}
