import { useEffect, useState } from "react";
import { financeApi } from "../../api";
import type {
  AnnualRevenueAccount,
  CompetitionIllusion,
  CostPriceFetish,
  ValueAppearanceContrast,
} from "../../types";
import "./Ch50Illusions.css";

export function Ch50Illusions() {
  const [illusions, setIllusions] = useState<CompetitionIllusion[]>([]);
  const [contrasts, setContrasts] = useState<ValueAppearanceContrast[]>([]);
  const [listError, setListError] = useState<string | null>(null);

  // Cost-price fetish builder (individual capitalist; defaults = seed).
  const [c, setC] = useState(70);
  const [v, setV] = useState(30);
  const [profit, setProfit] = useState(20);
  const [fetishResult, setFetishResult] = useState<CostPriceFetish | null>(null);
  const [fetishError, setFetishError] = useState<string | null>(null);

  // Annual revenue account builder (society-wide; defaults = seed).
  const [accC, setAccC] = useState(6000);
  const [wages, setWages] = useState(1500);
  const [profitRent, setProfitRent] = useState(1500);
  const [accResult, setAccResult] = useState<AnnualRevenueAccount | null>(null);
  const [accError, setAccError] = useState<string | null>(null);

  const costPrice = c + v;
  const actualValue = costPrice + profit;
  const newValue = wages + profitRent;
  const totalProduct = accC + newValue;

  async function refreshLists() {
    try {
      const [ills, cons] = await Promise.all([
        financeApi.listCompetitionIllusions(),
        financeApi.listValueAppearanceContrasts(),
      ]);
      setIllusions(ills);
      setContrasts(cons);
    } catch (e) {
      setListError(String(e));
    }
  }

  useEffect(() => {
    refreshLists().catch(() => {});
  }, []);

  async function handleCreateFetish() {
    setFetishError(null);
    setFetishResult(null);
    try {
      const r = await financeApi.createCostPriceFetish({
        constant_capital_bp: c,
        variable_capital_bp: v,
        profit_bp: profit,
      });
      setFetishResult(r);
    } catch (e) {
      setFetishError(String(e));
    }
  }

  async function handleCreateAccount() {
    setAccError(null);
    setAccResult(null);
    try {
      const r = await financeApi.createAnnualRevenueAccount({
        constant_capital_lm: accC,
        wages_share_lm: wages,
        profit_and_rent_share_lm: profitRent,
      });
      setAccResult(r);
    } catch (e) {
      setAccError(String(e));
    }
  }

  return (
    <div className="v3-ch50">
      <section className="v3-ch50-intro">
        <p>
          Chapter&nbsp;50 examines how the sphere of <strong>competition</strong> — market prices,
          buying and selling — generates systematic illusions that hide the origin of value and
          surplus-value in production.
        </p>
        <p>
          To the individual capitalist, the <strong>cost-price</strong> (c&nbsp;+&nbsp;v) appears as
          the commodity&apos;s &ldquo;real value&rdquo;, and the profit added above it appears to
          arise in the <em>sale</em> — not in the production process where unpaid surplus-labour
          actually created it. Competition makes profit of enterprise look like wages of supervision,
          and interest look like a property of capital as a thing.
        </p>
        <p>
          At the level of society, the annual revenue is the new value{" "}
          <strong>v&nbsp;+&nbsp;s</strong> (wages + profit-and-rent); the constant capital c appears
          to reproduce itself automatically through exchange between departments, obscuring that past
          labour created it.
        </p>
        <p className="v3-ch50-note">
          Cost-price amounts in basis points; social aggregates in labour-minutes.
        </p>
      </section>

      <section className="v3-ch50-section" aria-label="Competition illusions">
        <h4>1. Competition illusions</h4>
        <p className="v3-ch50-subintro">
          Each illusion contrasts its surface appearance with the real relation it conceals, and
          names the sphere of competition that produces it.
        </p>
        {listError && <p className="v3-ch50-error">{listError}</p>}
        {illusions.length > 0 && (
          <div className="v3-ch50-table-wrap">
            <table className="v3-ch50-table">
              <thead>
                <tr>
                  <th>Name</th>
                  <th>Surface appearance</th>
                  <th>Real relation</th>
                  <th>Produced by</th>
                </tr>
              </thead>
              <tbody>
                {illusions.map((i) => (
                  <tr key={i.id}>
                    <td>{i.name}</td>
                    <td>{i.surface_appearance}</td>
                    <td>{i.real_relation}</td>
                    <td>{i.produced_by}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </section>

      <section className="v3-ch50-section" aria-label="Cost-price fetish">
        <h4>2. Cost-price fetish</h4>
        <p className="v3-ch50-subintro">
          The individual capitalist&apos;s view: cost-price = c&nbsp;+&nbsp;v appears as the
          commodity&apos;s real value; profit appears as an addition made in exchange. The server
          computes cost-price and actual value (= cost-price&nbsp;+&nbsp;profit).
        </p>
        <div className="v3-ch50-form">
          <label>
            Constant capital c (bp)
            <input type="number" min={0} value={c} onChange={(e) => setC(Math.max(0, Number(e.target.value)))} />
          </label>
          <label>
            Variable capital v (bp)
            <input type="number" min={0} value={v} onChange={(e) => setV(Math.max(0, Number(e.target.value)))} />
          </label>
          <label>
            Profit (bp)
            <input type="number" min={0} value={profit} onChange={(e) => setProfit(Math.max(0, Number(e.target.value)))} />
          </label>
        </div>
        <p className="v3-ch50-totals">
          cost-price (c+v) = <strong>{costPrice}</strong>&nbsp;bp&nbsp;·&nbsp; actual value
          (c+v+s) = <strong>{actualValue}</strong>&nbsp;bp
        </p>
        <button className="v3-ch50-btn" onClick={handleCreateFetish}>
          Record cost-price fetish
        </button>
        {fetishError && <p className="v3-ch50-error">{fetishError}</p>}
        {fetishResult && (
          <div className="v3-ch50-preview">
            <strong>
              Saved — cost-price {fetishResult.cost_price_bp}&nbsp;bp, actual value{" "}
              {fetishResult.actual_value_bp}&nbsp;bp
            </strong>
            <span className="v3-ch50-hint">ID: {fetishResult.id}</span>
          </div>
        )}
      </section>

      <section className="v3-ch50-section" aria-label="Annual revenue account">
        <h4>3. Annual revenue account (society-wide)</h4>
        <p className="v3-ch50-subintro">
          New value (v&nbsp;+&nbsp;s) is the total distributable revenue; the constant capital c is
          reproduced, not distributed. The server computes new value and total social product (c +
          new value).
        </p>
        <div className="v3-ch50-form">
          <label>
            Constant capital c (lm)
            <input type="number" min={0} value={accC} onChange={(e) => setAccC(Math.max(0, Number(e.target.value)))} />
          </label>
          <label>
            Wages share v (lm)
            <input type="number" min={0} value={wages} onChange={(e) => setWages(Math.max(0, Number(e.target.value)))} />
          </label>
          <label>
            Profit &amp; rent share s (lm)
            <input type="number" min={0} value={profitRent} onChange={(e) => setProfitRent(Math.max(0, Number(e.target.value)))} />
          </label>
        </div>
        <p className="v3-ch50-totals">
          new value (v+s) = <strong>{newValue}</strong>&nbsp;lm&nbsp;·&nbsp; total social product
          (c + new) = <strong>{totalProduct}</strong>&nbsp;lm
        </p>
        <button className="v3-ch50-btn" onClick={handleCreateAccount}>
          Record annual account
        </button>
        {accError && <p className="v3-ch50-error">{accError}</p>}
        {accResult && (
          <div className="v3-ch50-preview">
            <strong>
              Saved — new value {accResult.new_value_created_lm}&nbsp;lm, total social product{" "}
              {accResult.total_social_product_lm}&nbsp;lm
            </strong>
            <span className="v3-ch50-hint">ID: {accResult.id}</span>
          </div>
        )}
      </section>

      <section className="v3-ch50-section" aria-label="Value appearance contrasts">
        <h4>4. Surface vs reality</h4>
        <p className="v3-ch50-subintro">
          For selected illusions, what market agents believe against what analysis reveals.
        </p>
        {contrasts.length > 0 && (
          <div className="v3-ch50-table-wrap">
            <table className="v3-ch50-table">
              <thead>
                <tr>
                  <th>Surface (what agents believe)</th>
                  <th>Reality (what analysis reveals)</th>
                </tr>
              </thead>
              <tbody>
                {contrasts.map((c) => (
                  <tr key={c.id}>
                    <td>{c.surface}</td>
                    <td>{c.reality}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </section>
    </div>
  );
}
