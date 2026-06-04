import { useEffect, useState } from "react";
import { financeApi } from "../../api";
import type {
  AgriculturalComposition,
  ValuePriceGap,
  AbsoluteRent,
  AbsoluteRentLimit,
} from "../../types";
import "./Ch45AbsoluteRent.css";

function formatBp(bp: number): string {
  return `${(bp / 100).toFixed(2)}%`;
}

export function Ch45AbsoluteRent() {
  const [compositions, setCompositions] = useState<AgriculturalComposition[]>([]);
  const [valueGaps, setValueGaps] = useState<ValuePriceGap[]>([]);
  const [absoluteRents, setAbsoluteRents] = useState<AbsoluteRent[]>([]);
  const [rentLimits, setRentLimits] = useState<AbsoluteRentLimit[]>([]);
  const [listError, setListError] = useState<string | null>(null);

  // Agricultural composition form
  const [acConstant, setAcConstant] = useState(60);
  const [acVariable, setAcVariable] = useState(40);
  const [acSocialAvg, setAcSocialAvg] = useState(40000);
  const [acError, setAcError] = useState<string | null>(null);
  const [lastAc, setLastAc] = useState<AgriculturalComposition | null>(null);

  // Value-price gap form
  const [vpProductId, setVpProductId] = useState("agricultural-output");
  const [vpValue, setVpValue] = useState(140);
  const [vpPrice, setVpPrice] = useState(120);
  const [vpError, setVpError] = useState<string | null>(null);
  const [lastVp, setLastVp] = useState<ValuePriceGap | null>(null);

  // Absolute rent form
  const [arParcelId, setArParcelId] = useState("5eed000000000000003700");
  const [arGapId, setArGapId] = useState("5eed000000000000004501");
  const [arSurplus, setArSurplus] = useState(40);
  const [arAvgProfit, setArAvgProfit] = useState(20);
  const [arError, setArError] = useState<string | null>(null);
  const [lastAr, setLastAr] = useState<AbsoluteRent | null>(null);

  // Rent limit form
  const [rlMax, setRlMax] = useState(20);
  const [rlActual, setRlActual] = useState(20);
  const [rlError, setRlError] = useState<string | null>(null);
  const [lastRl, setLastRl] = useState<AbsoluteRentLimit | null>(null);

  async function refreshLists() {
    try {
      const [comps, gaps, rents, limits] = await Promise.all([
        financeApi.listAgriculturalCompositions(),
        financeApi.listValuePriceGaps(),
        financeApi.listAbsoluteRents(),
        financeApi.listAbsoluteRentLimits(),
      ]);
      setCompositions(comps);
      setValueGaps(gaps);
      setAbsoluteRents(rents);
      setRentLimits(limits);
    } catch (e) {
      setListError(String(e));
    }
  }

  useEffect(() => {
    refreshLists().catch(() => {});
  }, []);

  async function handleCreateComposition() {
    setAcError(null);
    setLastAc(null);
    try {
      const result = await financeApi.createAgriculturalComposition({
        constant_capital_bp: acConstant,
        variable_capital_bp: acVariable,
        social_average_ratio_bp: acSocialAvg,
      });
      setLastAc(result);
      await refreshLists();
    } catch (e) {
      setAcError(String(e));
    }
  }

  async function handleCreateValueGap() {
    setVpError(null);
    setLastVp(null);
    try {
      const result = await financeApi.createValuePriceGap({
        product_id: vpProductId,
        value_labour_minutes: vpValue,
        price_of_production_labour_minutes: vpPrice,
      });
      setLastVp(result);
      await refreshLists();
    } catch (e) {
      setVpError(String(e));
    }
  }

  async function handleCreateAbsoluteRent() {
    setArError(null);
    setLastAr(null);
    try {
      const result = await financeApi.createAbsoluteRent({
        parcel_id: arParcelId,
        value_price_gap_id: arGapId,
        surplus_value_labour_minutes: arSurplus,
        average_profit_labour_minutes: arAvgProfit,
      });
      setLastAr(result);
      await refreshLists();
    } catch (e) {
      setArError(String(e));
    }
  }

  async function handleCreateRentLimit() {
    setRlError(null);
    setLastRl(null);
    try {
      const result = await financeApi.createAbsoluteRentLimit({
        max_rent_bp: rlMax,
        actual_rent_bp: rlActual,
      });
      setLastRl(result);
      await refreshLists();
    } catch (e) {
      setRlError(String(e));
    }
  }

  return (
    <div className="v3-ch45">
      <section className="v3-ch45-intro">
        <p>
          Chapter&nbsp;45 establishes that even the worst soil — which yields no differential
          rent — nonetheless always yields a rent: absolute ground-rent. The source lies in
          agriculture's characteristically <em>lower organic composition of capital</em> (more
          variable relative to constant) than the social average. Because agricultural products
          are sold near their value rather than at their price of production, they contain a
          surplus-value above the average profit. Landed property captures this excess by blocking
          capital from entering land without paying rent.
        </p>
        <p>
          The <strong>magnitude</strong> of absolute rent is bounded by the value–price-of-production
          gap: <code>AR = max(0, S&nbsp;&minus;&nbsp;average&nbsp;profit)</code>. If the organic
          composition of agriculture were to rise to the social average, the gap would close and
          absolute rent would vanish — only differential rent would remain. The two forms of
          ground-rent (DR + AR) coexist: DR explains surplus-profit above the average; AR explains
          why even the worst land commands a rent at all.
        </p>
        <p className="v3-ch45-note">
          Rates in basis points (bp). 10,000&nbsp;bp&nbsp;= 100%. Labour values in minutes.
          Ratio = C/V &times;&nbsp;10,000 (bp).
        </p>
      </section>

      {/* Section 1 — Agricultural compositions */}
      <section className="v3-ch45-section" aria-label="Agricultural compositions">
        <h4>1. Agricultural compositions</h4>
        <p className="v3-ch45-subintro">
          Record the organic composition of agricultural capital and compare it with the social
          average. The server computes <code>composition_ratio_bp = roundHalfUp(C&times;10000,&nbsp;V)</code>{" "}
          and sets <code>is_below_average</code> when the ratio falls below{" "}
          <code>social_average_ratio_bp</code>. A below-average composition is the necessary
          condition for absolute rent to exist. Seed fixture: C&nbsp;=&nbsp;60, V&nbsp;=&nbsp;40,
          ratio&nbsp;=&nbsp;15,000&nbsp;bp (1.5), social average&nbsp;=&nbsp;40,000&nbsp;bp (4.0) — below average.
        </p>
        <div className="v3-ch45-form">
          <label>
            Constant capital (bp)
            <input
              type="number"
              min={0}
              value={acConstant}
              onChange={(e) => setAcConstant(Math.max(0, Number(e.target.value)))}
            />
          </label>
          <label>
            Variable capital (bp)
            <input
              type="number"
              min={1}
              value={acVariable}
              onChange={(e) => setAcVariable(Math.max(1, Number(e.target.value)))}
            />
          </label>
          <label>
            Social average ratio (bp)
            <input
              type="number"
              min={0}
              value={acSocialAvg}
              onChange={(e) => setAcSocialAvg(Math.max(0, Number(e.target.value)))}
            />
            <span className="v3-ch45-hint">{formatBp(acSocialAvg)} C/V ratio</span>
          </label>
        </div>
        <button className="v3-ch45-btn" onClick={handleCreateComposition}>
          Record composition
        </button>
        {acError && <p className="v3-ch45-error">{acError}</p>}
        {lastAc && (
          <div className="v3-ch45-preview">
            <strong>
              Saved — C&nbsp;{lastAc.constant_capital_bp}&nbsp;V&nbsp;{lastAc.variable_capital_bp},
              ratio&nbsp;{lastAc.composition_ratio_bp}&nbsp;bp,
              below average:&nbsp;{lastAc.is_below_average ? "yes" : "no"}
            </strong>
            <span className="v3-ch45-hint">ID: {lastAc.id}</span>
          </div>
        )}
        {compositions.length > 0 && (
          <div className="v3-ch45-table-wrap">
            <table className="v3-ch45-table">
              <thead>
                <tr>
                  <th>ID</th>
                  <th>C (bp)</th>
                  <th>V (bp)</th>
                  <th>Ratio (bp)</th>
                  <th>Social avg (bp)</th>
                  <th>Below average?</th>
                </tr>
              </thead>
              <tbody>
                {compositions.map((c) => (
                  <tr key={c.id}>
                    <td title={c.id}>{c.id.slice(0, 8)}&hellip;</td>
                    <td>{c.constant_capital_bp}</td>
                    <td>{c.variable_capital_bp}</td>
                    <td>{c.composition_ratio_bp}</td>
                    <td>{c.social_average_ratio_bp}</td>
                    <td>{c.is_below_average ? "yes" : "no"}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
        {listError && <p className="v3-ch45-error">{listError}</p>}
      </section>

      {/* Section 2 — Value-price gaps */}
      <section className="v3-ch45-section" aria-label="Value-price gaps">
        <h4>2. Value-price gaps</h4>
        <p className="v3-ch45-subintro">
          Record the gap between the value of an agricultural product (in labour-minutes) and its
          price of production. The server computes{" "}
          <code>gap_labour_minutes&nbsp;= value&nbsp;&minus;&nbsp;price</code>. A positive gap
          is the material basis for absolute rent. Seed fixture: value&nbsp;=&nbsp;140&nbsp;min,
          price&nbsp;=&nbsp;120&nbsp;min, gap&nbsp;=&nbsp;20&nbsp;min.
        </p>
        <div className="v3-ch45-form">
          <label>
            Product ID
            <input
              value={vpProductId}
              onChange={(e) => setVpProductId(e.target.value)}
              placeholder="agricultural-output"
            />
          </label>
          <label>
            Value (labour-minutes)
            <input
              type="number"
              min={0}
              value={vpValue}
              onChange={(e) => setVpValue(Math.max(0, Number(e.target.value)))}
            />
          </label>
          <label>
            Price of production (labour-minutes)
            <input
              type="number"
              min={0}
              value={vpPrice}
              onChange={(e) => setVpPrice(Math.max(0, Number(e.target.value)))}
            />
            <span className="v3-ch45-hint">gap = {vpValue - vpPrice} min</span>
          </label>
        </div>
        <button className="v3-ch45-btn" onClick={handleCreateValueGap}>
          Record value-price gap
        </button>
        {vpError && <p className="v3-ch45-error">{vpError}</p>}
        {lastVp && (
          <div className="v3-ch45-preview">
            <strong>
              Saved — {lastVp.product_id}: value&nbsp;{lastVp.value_labour_minutes}&nbsp;min,
              price&nbsp;{lastVp.price_of_production_labour_minutes}&nbsp;min,
              gap&nbsp;{lastVp.gap_labour_minutes}&nbsp;min
            </strong>
            <span className="v3-ch45-hint">ID: {lastVp.id}</span>
          </div>
        )}
        {valueGaps.length > 0 && (
          <div className="v3-ch45-table-wrap">
            <table className="v3-ch45-table">
              <thead>
                <tr>
                  <th>ID</th>
                  <th>Product</th>
                  <th>Value (min)</th>
                  <th>Price (min)</th>
                  <th>Gap (min)</th>
                </tr>
              </thead>
              <tbody>
                {valueGaps.map((g) => (
                  <tr key={g.id}>
                    <td title={g.id}>{g.id.slice(0, 8)}&hellip;</td>
                    <td>{g.product_id}</td>
                    <td>{g.value_labour_minutes}</td>
                    <td>{g.price_of_production_labour_minutes}</td>
                    <td>{g.gap_labour_minutes}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </section>

      {/* Section 3 — Absolute rents */}
      <section className="v3-ch45-section" aria-label="Absolute rents">
        <h4>3. Absolute rents</h4>
        <p className="v3-ch45-subintro">
          Compute the absolute rent for a parcel given the surplus-value and average profit
          (both in basis points). The server computes{" "}
          <code>absolute_rent_labour_minutes = max(0, surplus_value_labour_minutes&nbsp;&minus;&nbsp;average_profit_labour_minutes)</code>.
          Seed fixture: parcel <code>4502</code>, S&nbsp;=&nbsp;40&nbsp;bp,
          average&nbsp;profit&nbsp;=&nbsp;20&nbsp;bp, absolute&nbsp;rent&nbsp;=&nbsp;20&nbsp;bp.
        </p>
        <div className="v3-ch45-form">
          <label>
            Parcel ID
            <input
              value={arParcelId}
              onChange={(e) => setArParcelId(e.target.value)}
              placeholder="5eed000000000000003700"
            />
          </label>
          <label>
            Value-price gap ID
            <input
              value={arGapId}
              onChange={(e) => setArGapId(e.target.value)}
              placeholder="5eed000000000000004501"
            />
          </label>
          <label>
            Surplus value (bp)
            <input
              type="number"
              min={0}
              value={arSurplus}
              onChange={(e) => setArSurplus(Math.max(0, Number(e.target.value)))}
            />
          </label>
          <label>
            Average profit (bp)
            <input
              type="number"
              min={0}
              value={arAvgProfit}
              onChange={(e) => setArAvgProfit(Math.max(0, Number(e.target.value)))}
            />
            <span className="v3-ch45-hint">
              absolute rent = max(0,&nbsp;{arSurplus}&nbsp;&minus;&nbsp;{arAvgProfit})&nbsp;=&nbsp;
              {Math.max(0, arSurplus - arAvgProfit)}&nbsp;bp
            </span>
          </label>
        </div>
        <button className="v3-ch45-btn" onClick={handleCreateAbsoluteRent}>
          Record absolute rent
        </button>
        {arError && <p className="v3-ch45-error">{arError}</p>}
        {lastAr && (
          <div className="v3-ch45-preview">
            <strong>
              Saved — parcel&nbsp;{lastAr.parcel_id.slice(0, 8)}&hellip;,
              S&nbsp;=&nbsp;{lastAr.surplus_value_labour_minutes}&nbsp;LM,
              avg&nbsp;profit&nbsp;=&nbsp;{lastAr.average_profit_labour_minutes}&nbsp;LM,
              absolute&nbsp;rent&nbsp;=&nbsp;{lastAr.absolute_rent_labour_minutes}&nbsp;LM
            </strong>
            <span className="v3-ch45-hint">ID: {lastAr.id}</span>
          </div>
        )}
        {absoluteRents.length > 0 && (
          <div className="v3-ch45-table-wrap">
            <table className="v3-ch45-table">
              <thead>
                <tr>
                  <th>ID</th>
                  <th>Parcel</th>
                  <th>Surplus value (bp)</th>
                  <th>Avg profit (bp)</th>
                  <th>Absolute rent (bp)</th>
                </tr>
              </thead>
              <tbody>
                {absoluteRents.map((r) => (
                  <tr key={r.id}>
                    <td title={r.id}>{r.id.slice(0, 8)}&hellip;</td>
                    <td title={r.parcel_id}>{r.parcel_id.slice(0, 8)}&hellip;</td>
                    <td>{r.surplus_value_labour_minutes}</td>
                    <td>{r.average_profit_labour_minutes}</td>
                    <td>{r.absolute_rent_labour_minutes}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </section>

      {/* Section 4 — Rent limits */}
      <section className="v3-ch45-section" aria-label="Rent limits">
        <h4>4. Rent limits</h4>
        <p className="v3-ch45-subintro">
          The value–price-of-production gap sets an upper bound on absolute rent. Record the
          maximum rent (= the gap) and the actual rent charged. The server flags{" "}
          <code>is_at_limit</code> when <code>actual_rent_bp &ge; max_rent_bp</code>.
          If agriculture's composition rose to the social average the gap would close and
          the maximum would fall to zero — absolute rent would vanish. Seed fixture:
          max&nbsp;=&nbsp;20&nbsp;bp, actual&nbsp;=&nbsp;20&nbsp;bp, at limit.
        </p>
        <div className="v3-ch45-form">
          <label>
            Maximum rent (bp)
            <input
              type="number"
              min={0}
              value={rlMax}
              onChange={(e) => setRlMax(Math.max(0, Number(e.target.value)))}
            />
          </label>
          <label>
            Actual rent (bp)
            <input
              type="number"
              min={0}
              value={rlActual}
              onChange={(e) => setRlActual(Math.max(0, Number(e.target.value)))}
            />
            <span className="v3-ch45-hint">
              at limit: {rlActual >= rlMax ? "yes" : "no"}
            </span>
          </label>
        </div>
        <button className="v3-ch45-btn" onClick={handleCreateRentLimit}>
          Record rent limit
        </button>
        {rlError && <p className="v3-ch45-error">{rlError}</p>}
        {lastRl && (
          <div className="v3-ch45-preview">
            <strong>
              Saved — max&nbsp;{lastRl.max_rent_bp}&nbsp;bp,
              actual&nbsp;{lastRl.actual_rent_bp}&nbsp;bp,
              at limit:&nbsp;{lastRl.is_at_limit ? "yes" : "no"}
            </strong>
            <span className="v3-ch45-hint">ID: {lastRl.id}</span>
          </div>
        )}
        {rentLimits.length > 0 && (
          <div className="v3-ch45-table-wrap">
            <table className="v3-ch45-table">
              <thead>
                <tr>
                  <th>ID</th>
                  <th>Max rent (bp)</th>
                  <th>Actual rent (bp)</th>
                  <th>At limit?</th>
                </tr>
              </thead>
              <tbody>
                {rentLimits.map((l) => (
                  <tr key={l.id}>
                    <td title={l.id}>{l.id.slice(0, 8)}&hellip;</td>
                    <td>{l.max_rent_bp}</td>
                    <td>{l.actual_rent_bp}</td>
                    <td>{l.is_at_limit ? "yes" : "no"}</td>
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
