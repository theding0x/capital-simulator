import { useEffect, useState } from "react";
import { financeApi } from "../../api";
import type { MerchantTurnoverResponse, TurnoverEffectAnalysisResponse } from "../../types";
import "./Ch18MerchantTurnover.css";

function pence(n: number): string {
  const pounds = n / 100;
  return `£${pounds.toLocaleString("en-GB", { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`;
}

function bp(n: number): string {
  return `${(n / 100).toFixed(2)}%`;
}

// Marx Ch. 18 fixture: capital=£1,000 (100000p), general rate=15% (1500 bp).
// annual profit = roundHalfUp(100000*1500, 10000) = £150 (15000p).
// turnover=1: markup = 15000 / 300 = 50p per unit.
// turnover=5: markup = 15000 / 1500 = 10p per unit.
const FIXTURE = {
  money_advanced: 100000,
  general_rate_bp: 1500,
  annual_profit: 15000,
  turnover_1_markup: 50,
  turnover_5_markup: 10,
};

export function Ch18MerchantTurnover() {
  const [items, setItems] = useState<MerchantTurnoverResponse[]>([]);
  const [listError, setListError] = useState<string | null>(null);

  // Create-form state.
  const [moneyAdvanced, setMoneyAdvanced] = useState(100000);
  const [generalRateBP, setGeneralRateBP] = useState(1500);
  const [turnoverCount, setTurnoverCount] = useState(1);
  const [unitsPerTurnover, setUnitsPerTurnover] = useState(300);
  const [createError, setCreateError] = useState<string | null>(null);

  // Turnover-effect compare panel state.
  const [effectMoney, setEffectMoney] = useState(100000);
  const [effectRate, setEffectRate] = useState(1500);
  const [effectUnits, setEffectUnits] = useState(300);
  const [effectCountsStr, setEffectCountsStr] = useState("1,2,3,4,5");
  const [effectResult, setEffectResult] = useState<TurnoverEffectAnalysisResponse | null>(null);
  const [effectError, setEffectError] = useState<string | null>(null);

  async function refreshItems() {
    try {
      const data = await financeApi.listMerchantTurnovers();
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
      await financeApi.createMerchantTurnover({
        money_advanced: moneyAdvanced,
        general_rate_bp: generalRateBP,
        turnover_count: turnoverCount,
        units_per_turnover: unitsPerTurnover,
      });
      await refreshItems();
    } catch (e) {
      setCreateError(String(e));
    }
  }

  async function handleComputeEffect() {
    setEffectError(null);
    setEffectResult(null);
    const counts = effectCountsStr
      .split(",")
      .map((s) => parseInt(s.trim(), 10))
      .filter((n) => !isNaN(n));
    try {
      const result = await financeApi.computeTurnoverEffect({
        money_advanced: effectMoney,
        general_rate_bp: effectRate,
        units_per_turnover: effectUnits,
        turnover_counts: counts,
      });
      setEffectResult(result);
    } catch (e) {
      setEffectError(String(e));
    }
  }

  // Client-side preview: annual profit (round-half-up) and per-unit markup (integer div).
  function previewAnnual(money: number, rateBP: number): number {
    if (money <= 0 || rateBP <= 0) return 0;
    return Math.floor((money * rateBP + 5000) / 10000);
  }
  const annual = previewAnnual(moneyAdvanced, generalRateBP);
  const totalUnits = turnoverCount * unitsPerTurnover;
  const markupPreview = totalUnits > 0 ? Math.floor(annual / totalUnits) : 0;

  return (
    <div className="v3-ch18">
      <section className="v3-ch18-intro">
        <p>
          The annual profit on merchant's capital depends on the money advanced and the general
          rate of profit — <strong>not</strong> on how many times the capital turns over or how
          many units are sold. Faster turnover lowers the per-unit markup but leaves annual
          profit unchanged: <em>profit ∝ capital, markup ∝ 1/turnover</em>.
        </p>
      </section>

      {/* Marx Ch. 18 fixture */}
      <section className="v3-ch18-section" aria-label="Marx Ch. 18 fixture">
        <h4>Marx Ch. 18 fixture</h4>
        <p className="v3-ch18-subintro">
          Capital {pence(FIXTURE.money_advanced)} · General rate {bp(FIXTURE.general_rate_bp)} ·
          Annual profit {pence(FIXTURE.annual_profit)} · Units per turnover 300
        </p>
        <div className="v3-ch18-fixture-grid">
          <div className="v3-ch18-fixture-cell">
            <span className="v3-ch18-fixture-label">Turnover ×1</span>
            <span className="v3-ch18-fixture-value">{pence(FIXTURE.turnover_1_markup)} / unit</span>
          </div>
          <div className="v3-ch18-fixture-cell v3-ch18-fixture-cell--fast">
            <span className="v3-ch18-fixture-label">Turnover ×5</span>
            <span className="v3-ch18-fixture-value">{pence(FIXTURE.turnover_5_markup)} / unit</span>
          </div>
          <div className="v3-ch18-fixture-cell v3-ch18-fixture-cell--profit">
            <span className="v3-ch18-fixture-label">Annual profit (both)</span>
            <span className="v3-ch18-fixture-value">{pence(FIXTURE.annual_profit)}</span>
          </div>
        </div>
      </section>

      {/* Turnover-effect compare panel */}
      <section className="v3-ch18-section" aria-label="Turnover effect analysis">
        <h4>Turnover-effect analysis — markup shrinks as turnover rises</h4>
        <div className="v3-ch18-form">
          <label>
            Money advanced (pence)
            <input
              type="number"
              min={1}
              value={effectMoney}
              onChange={(e) => setEffectMoney(Number(e.target.value))}
            />
          </label>
          <label>
            General rate (bp)
            <input
              type="number"
              min={1}
              value={effectRate}
              onChange={(e) => setEffectRate(Number(e.target.value))}
            />
          </label>
          <label>
            Units per turnover
            <input
              type="number"
              min={1}
              value={effectUnits}
              onChange={(e) => setEffectUnits(Number(e.target.value))}
            />
          </label>
          <label>
            Turnover counts (comma-separated)
            <input
              type="text"
              value={effectCountsStr}
              onChange={(e) => setEffectCountsStr(e.target.value)}
            />
          </label>
          <button type="button" className="v3-ch18-btn" onClick={handleComputeEffect}>
            Compare
          </button>
        </div>
        {effectError && <p className="v3-ch18-error">{effectError}</p>}
        {effectResult && (
          <div className="v3-ch18-effect-result">
            <p className="v3-ch18-subintro">
              Annual profit: <strong>{pence(effectResult.annual_merchant_profit)}</strong> (invariant across all turnover counts)
            </p>
            <table className="v3-ch18-table">
              <thead>
                <tr>
                  <th>Turnover ×</th>
                  <th>Total units</th>
                  <th>Markup / unit</th>
                </tr>
              </thead>
              <tbody>
                {effectResult.points.map((pt) => (
                  <tr key={pt.turnover_count}>
                    <td className="v3-ch18-num">{pt.turnover_count}</td>
                    <td className="v3-ch18-num">{pt.total_units.toLocaleString("en-GB")}</td>
                    <td className="v3-ch18-money">{pence(pt.markup_per_unit)}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </section>

      {/* Create merchant turnover */}
      <section className="v3-ch18-section" aria-label="Record merchant turnover">
        <h4>Record a merchant-turnover computation</h4>
        <div className="v3-ch18-form">
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
          <label>
            Turnover count
            <input
              type="number"
              min={1}
              value={turnoverCount}
              onChange={(e) => setTurnoverCount(Number(e.target.value))}
            />
          </label>
          <label>
            Units per turnover
            <input
              type="number"
              min={1}
              value={unitsPerTurnover}
              onChange={(e) => setUnitsPerTurnover(Number(e.target.value))}
            />
          </label>
          <div className="v3-ch18-preview-pair">
            <div className="v3-ch18-preview">
              <span className="v3-ch18-preview-label">Annual profit</span>
              <span className="v3-ch18-preview-figure">{pence(annual)}</span>
            </div>
            <div className="v3-ch18-preview v3-ch18-preview--markup">
              <span className="v3-ch18-preview-label">Markup / unit</span>
              <span className="v3-ch18-preview-figure">{pence(markupPreview)}</span>
            </div>
          </div>
          <button type="button" className="v3-ch18-btn" onClick={handleCreate}>
            Record
          </button>
        </div>
        {createError && <p className="v3-ch18-error">{createError}</p>}
      </section>

      {/* List */}
      {listError && <p className="v3-ch18-error">{listError}</p>}
      {items.length > 0 && (
        <section className="v3-ch18-section" aria-label="Merchant turnover records">
          <h4>Merchant-turnover records ({items.length})</h4>
          <table className="v3-ch18-table">
            <thead>
              <tr>
                <th>Money advanced</th>
                <th>Rate</th>
                <th>Turnover ×</th>
                <th>Units / turn</th>
                <th>Annual profit</th>
                <th>Markup / unit</th>
              </tr>
            </thead>
            <tbody>
              {items.map((mt) => (
                <tr key={mt.id}>
                  <td className="v3-ch18-money">{pence(mt.money_advanced)}</td>
                  <td className="v3-ch18-rate">{bp(mt.general_rate_bp)}</td>
                  <td className="v3-ch18-num">{mt.turnover_count}</td>
                  <td className="v3-ch18-num">{mt.units_per_turnover.toLocaleString("en-GB")}</td>
                  <td className="v3-ch18-money v3-ch18-profit">{pence(mt.annual_merchant_profit)}</td>
                  <td className="v3-ch18-money">{pence(mt.markup_per_unit)}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </section>
      )}
    </div>
  );
}
