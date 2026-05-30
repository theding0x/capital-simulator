import { useEffect, useState } from "react";
import { financeApi } from "../../api";
import type {
  EqualisationResponse,
  MarketValueResponse,
  SurplusProfitResponse,
} from "../../types";
import "./Ch10Equalisation.css";

function pct(bp: number): string {
  return `${(bp / 100).toLocaleString("en-GB", { maximumFractionDigits: 2 })}%`;
}

function money(n: number): string {
  return `£${n.toLocaleString("en-GB")}`;
}

function directionArrow(dir: string): string {
  if (dir === "inflow") return "→";
  if (dir === "outflow") return "←";
  return "–";
}

// Pre-seed inputs matching Marx's fixtures.
const COTTON_FLOW = { sphere_name: "cotton", current_sphere_rate: 3000, general_rate: 2200, market_price: 0, market_value: 0 };
const FIRM_SURPLUS = { firm_name: "High-productivity firm", individual_value: 60, market_value: 100, output_qty: 1000, general_rate: 2200 };
const COTTON_MV = { sphere_name: "cotton", bulk_condition_value: 100, best_condition_value: 80, worst_condition_value: 130 };

export function Ch10Equalisation() {
  const [equalisation, setEqualisation] = useState<EqualisationResponse | null>(null);
  const [surplusProfit, setSurplusProfit] = useState<SurplusProfitResponse | null>(null);
  const [marketValue, setMarketValue] = useState<MarketValueResponse | null>(null);
  const [error, setError] = useState<string | null>(null);

  // Capital-flow form state.
  const [flowInput, setFlowInput] = useState(COTTON_FLOW);
  // Surplus-profit form state.
  const [spInput, setSpInput] = useState(FIRM_SURPLUS);
  // Market-value form state.
  const [mvInput, setMvInput] = useState(COTTON_MV);

  // On mount, POST the pre-seed fixtures so the panel is populated.
  useEffect(() => {
    Promise.all([
      financeApi.createCapitalFlow(COTTON_FLOW),
      financeApi.computeSurplusProfit(FIRM_SURPLUS),
      financeApi.createMarketValue(COTTON_MV),
    ])
      .then(([eq, sp, mv]) => {
        setEqualisation(eq);
        setSurplusProfit(sp);
        setMarketValue(mv);
      })
      .catch(() => {});
  }, []);

  async function computeFlow() {
    setError(null);
    try {
      const eq = await financeApi.createCapitalFlow(flowInput);
      setEqualisation(eq);
    } catch (e) {
      setError(String(e));
    }
  }

  async function computeSurplus() {
    setError(null);
    try {
      const sp = await financeApi.computeSurplusProfit(spInput);
      setSurplusProfit(sp);
    } catch (e) {
      setError(String(e));
    }
  }

  async function computeMarketVal() {
    setError(null);
    try {
      const mv = await financeApi.createMarketValue(mvInput);
      setMarketValue(mv);
    } catch (e) {
      setError(String(e));
    }
  }

  return (
    <div className="v3-ch10">
      <section className="v3-ch10-intro">
        <p>
          Competition is the <em>mechanism</em> that forms the general rate.
          Capital flows <strong>into</strong> branches whose rate is above the
          general rate and <strong>out of</strong> those below it; supply shifts
          move market prices until each branch's rate converges. Within a
          branch, the <strong>market value</strong> is set by the conditions
          under which the <em>bulk</em> of the output is produced. A firm
          producing below the market value (higher individual productivity)
          sells at the market value and pockets the difference as{" "}
          <strong>surplus-profit</strong>.
        </p>
      </section>

      {error && <p className="v3-ch10-error">{error}</p>}

      {/* Capital-flow visualiser */}
      <section className="v3-ch10-section" aria-label="Capital-flow visualiser">
        <h4>Capital-flow direction</h4>
        <div className="v3-ch10-form">
          <label>
            Sphere
            <input type="text" value={flowInput.sphere_name}
              onChange={(e) => setFlowInput({ ...flowInput, sphere_name: e.target.value })} />
          </label>
          <label>
            Current rate (bp)
            <input type="number" min={0} value={flowInput.current_sphere_rate}
              onChange={(e) => setFlowInput({ ...flowInput, current_sphere_rate: Number(e.target.value) })} />
          </label>
          <label>
            General rate (bp)
            <input type="number" min={0} value={flowInput.general_rate}
              onChange={(e) => setFlowInput({ ...flowInput, general_rate: Number(e.target.value) })} />
          </label>
          <label>
            Market price
            <input type="number" min={0} value={flowInput.market_price}
              onChange={(e) => setFlowInput({ ...flowInput, market_price: Number(e.target.value) })} />
          </label>
          <label>
            Market value
            <input type="number" min={0} value={flowInput.market_value}
              onChange={(e) => setFlowInput({ ...flowInput, market_value: Number(e.target.value) })} />
          </label>
          <button type="button" className="v3-ch10-btn" onClick={computeFlow}>
            Compute flow
          </button>
        </div>

        {equalisation && (
          <div className={`v3-ch10-flow-result v3-ch10-flow-${equalisation.direction}`}>
            <span className="v3-ch10-flow-arrow">{directionArrow(equalisation.direction)}</span>
            <span className="v3-ch10-flow-label">{equalisation.direction.toUpperCase()}</span>
            <span className="v3-ch10-flow-detail">
              {equalisation.sphere_name}: {pct(equalisation.initial_rate)} vs general{" "}
              {pct(equalisation.target_rate)}
              {equalisation.is_converging
                ? " — converging"
                : " — already at general rate (stable)"}
            </span>
          </div>
        )}
      </section>

      {/* Surplus-profit calculator */}
      <section className="v3-ch10-section" aria-label="Surplus-profit calculator">
        <h4>Surplus-profit calculator</h4>
        <div className="v3-ch10-form">
          <label>
            Firm name
            <input type="text" value={spInput.firm_name}
              onChange={(e) => setSpInput({ ...spInput, firm_name: e.target.value })} />
          </label>
          <label>
            Individual value
            <input type="number" min={0} value={spInput.individual_value}
              onChange={(e) => setSpInput({ ...spInput, individual_value: Number(e.target.value) })} />
          </label>
          <label>
            Market value
            <input type="number" min={0} value={spInput.market_value}
              onChange={(e) => setSpInput({ ...spInput, market_value: Number(e.target.value) })} />
          </label>
          <label>
            Output qty
            <input type="number" min={0} value={spInput.output_qty}
              onChange={(e) => setSpInput({ ...spInput, output_qty: Number(e.target.value) })} />
          </label>
          <label>
            General rate (bp)
            <input type="number" min={0} value={spInput.general_rate}
              onChange={(e) => setSpInput({ ...spInput, general_rate: Number(e.target.value) })} />
          </label>
          <button type="button" className="v3-ch10-btn" onClick={computeSurplus}>
            Compute surplus-profit
          </button>
        </div>

        {surplusProfit && (
          <div className={`v3-ch10-sp-result ${surplusProfit.amount >= 0 ? "v3-ch10-sp-positive" : "v3-ch10-sp-negative"}`}>
            <span className="v3-ch10-sp-label">{surplusProfit.firm_name}</span>
            <span className="v3-ch10-sp-amount">
              surplus-profit = {money(surplusProfit.amount)}
            </span>
            <span className="v3-ch10-sp-formula">
              ({money(surplusProfit.market_value)} &minus; {money(surplusProfit.individual_value)})
              &times; {surplusProfit.output_qty.toLocaleString("en-GB")} units
            </span>
          </div>
        )}
      </section>

      {/* Market-value panel */}
      <section className="v3-ch10-section" aria-label="Market-value panel">
        <h4>Market value (bulk-condition rule)</h4>
        <div className="v3-ch10-form">
          <label>
            Sphere / commodity
            <input type="text" value={mvInput.sphere_name}
              onChange={(e) => setMvInput({ ...mvInput, sphere_name: e.target.value })} />
          </label>
          <label>
            Bulk condition value
            <input type="number" min={0} value={mvInput.bulk_condition_value}
              onChange={(e) => setMvInput({ ...mvInput, bulk_condition_value: Number(e.target.value) })} />
          </label>
          <label>
            Best condition value
            <input type="number" min={0} value={mvInput.best_condition_value}
              onChange={(e) => setMvInput({ ...mvInput, best_condition_value: Number(e.target.value) })} />
          </label>
          <label>
            Worst condition value
            <input type="number" min={0} value={mvInput.worst_condition_value}
              onChange={(e) => setMvInput({ ...mvInput, worst_condition_value: Number(e.target.value) })} />
          </label>
          <button type="button" className="v3-ch10-btn" onClick={computeMarketVal}>
            Set market value
          </button>
        </div>

        {marketValue && (
          <div className="v3-ch10-mv-result">
            <span className="v3-ch10-mv-label">{marketValue.sphere_name}</span>
            <span className="v3-ch10-mv-value">
              Market value = {money(marketValue.value)} (bulk condition)
            </span>
            <span className="v3-ch10-mv-range">
              best {money(marketValue.best_condition_value)} — bulk{" "}
              {money(marketValue.bulk_condition_value)} — worst{" "}
              {money(marketValue.worst_condition_value)}
            </span>
          </div>
        )}
      </section>
    </div>
  );
}
