import { useEffect, useState } from "react";
import { financeApi } from "../../api";
import type {
  PriceOfProductionChangeResponse,
  CompensationGroundResponse,
  PartIISummaryResponse,
} from "../../types";
import "./Ch12SupplementaryRemarks.css";

type Cause = "general_rate_change" | "individual_value_change" | "both";

function pct(bp: number): string {
  return `${(bp / 100).toLocaleString("en-GB", { maximumFractionDigits: 2 })}%`;
}

const CAUSE_LABELS: Record<Cause, string> = {
  general_rate_change: "The general rate of profit changed (own value unchanged)",
  individual_value_change: "The commodity's own value changed (general rate unchanged)",
  both: "Both the general rate and the commodity's own value changed",
};

export function Ch12SupplementaryRemarks() {
  // --- Remark 1: the two causes of a change in the price of production -------
  const [cause, setCause] = useState<Cause>("general_rate_change");
  const [oldPrice, setOldPrice] = useState(122);
  const [newPrice, setNewPrice] = useState(130);
  const [change, setChange] = useState<PriceOfProductionChangeResponse | null>(null);
  const [changeError, setChangeError] = useState<string | null>(null);

  // --- Remark 2: compensation grounds are not value creation ----------------
  const [reason, setReason] = useState("slow turnover");
  const [ground, setGround] = useState<CompensationGroundResponse | null>(null);
  const [groundError, setGroundError] = useState<string | null>(null);

  // --- Remark 3: Part II conservation summary -------------------------------
  const [summary, setSummary] = useState<PartIISummaryResponse | null>(null);

  useEffect(() => {
    financeApi.getPartIISummary().then(setSummary).catch(() => {});
    financeApi
      .createCompensationGround({ reason: "slow turnover", appears_to_be: "extra value creation" })
      .then(setGround)
      .catch(() => {});
  }, []);

  async function computeChange() {
    setChangeError(null);
    try {
      const r = await financeApi.createPriceOfProductionChange({
        sphere_name: "cotton",
        cause,
        old_price: oldPrice,
        new_price: newPrice,
      });
      setChange(r);
    } catch (e) {
      setChangeError(String(e));
    }
  }

  async function bustMyth() {
    setGroundError(null);
    try {
      const r = await financeApi.createCompensationGround({
        reason,
        appears_to_be: "extra value creation",
      });
      setGround(r);
    } catch (e) {
      setGroundError(String(e));
    }
  }

  return (
    <div className="v3-ch12">
      <section className="v3-ch12-intro">
        <p>
          Two closing remarks on prices of production. First, a commodity's price
          of production changes for <strong>only two reasons</strong>: the{" "}
          <strong>general rate of profit</strong> changes, or the commodity's{" "}
          <strong>own value</strong> changes. Second, a capitalist's claimed{" "}
          <strong>grounds</strong> for a higher price — slow turnover, risk,
          skill — create no new value; they are only a claim on a larger share of
          the aggregate surplus-value. Across the whole social capital, total
          price of production equals total value.
        </p>
      </section>

      <section className="v3-ch12-section" aria-label="Two causes of price change">
        <h4>The two causes of a change in price of production</h4>
        <div className="v3-ch12-form">
          <fieldset className="v3-ch12-causes">
            <legend>Cause</legend>
            {(Object.keys(CAUSE_LABELS) as Cause[]).map((c) => (
              <label key={c} className="v3-ch12-radio">
                <input
                  type="radio"
                  name="cause"
                  value={c}
                  checked={cause === c}
                  onChange={() => setCause(c)}
                />
                {CAUSE_LABELS[c]}
              </label>
            ))}
          </fieldset>
          <label>
            Old price
            <input
              type="number" min={0} value={oldPrice}
              onChange={(e) => setOldPrice(Number(e.target.value))}
            />
          </label>
          <label>
            New price
            <input
              type="number" min={0} value={newPrice}
              onChange={(e) => setNewPrice(Number(e.target.value))}
            />
          </label>
          <button type="button" className="v3-ch12-btn" onClick={computeChange}>
            Determine
          </button>
        </div>
        {changeError && <p className="v3-ch12-error">{changeError}</p>}
        {change && (
          <div className="v3-ch12-flags">
            <span className={`v3-ch12-flag ${change.rate_changed ? "on" : "off"}`}>
              General rate changed: {change.rate_changed ? "yes" : "no"}
            </span>
            <span className={`v3-ch12-flag ${change.value_changed ? "on" : "off"}`}>
              Own value changed: {change.value_changed ? "yes" : "no"}
            </span>
            <span className={`v3-ch12-flag ${change.price_changed ? "on" : "off"}`}>
              Price of production changed: {change.price_changed ? "yes" : "no"}
            </span>
          </div>
        )}
      </section>

      <section className="v3-ch12-section" aria-label="Compensation grounds">
        <h4>Compensation grounds — myth-buster</h4>
        <div className="v3-ch12-form">
          <label className="v3-ch12-grow">
            Claimed reason for a higher price
            <input
              type="text" value={reason}
              onChange={(e) => setReason(e.target.value)}
              placeholder="e.g. slow turnover, risk, skill"
            />
          </label>
          <button type="button" className="v3-ch12-btn" onClick={bustMyth}>
            Examine
          </button>
        </div>
        {groundError && <p className="v3-ch12-error">{groundError}</p>}
        {ground && (
          <div className="v3-ch12-myth">
            <div className="v3-ch12-myth-row">
              <span className="v3-ch12-myth-label">Reason</span>
              <span className="v3-ch12-myth-value">{ground.reason}</span>
            </div>
            <div className="v3-ch12-myth-row">
              <span className="v3-ch12-myth-label">Appears to be</span>
              <span className="v3-ch12-myth-appears">{ground.appears_to_be}</span>
            </div>
            <div className="v3-ch12-myth-row">
              <span className="v3-ch12-myth-label">Actually is</span>
              <span className="v3-ch12-myth-actual">{ground.actually_is}</span>
            </div>
            <div className="v3-ch12-myth-row">
              <span className="v3-ch12-myth-label">Creates value?</span>
              <span className={`v3-ch12-tag ${ground.creates_value ? "v3-ch12-tag-yes" : "v3-ch12-tag-no"}`}>
                {ground.creates_value ? "yes" : "no"}
              </span>
            </div>
          </div>
        )}
      </section>

      {summary && (
        <section className="v3-ch12-section" aria-label="Part II summary">
          <h4>Part II — conservation of value</h4>
          <div className={`v3-ch12-banner ${summary.value_conserved ? "ok" : "bad"}`}>
            {summary.value_conserved
              ? "Σ prices of production = Σ values — value is conserved across the social capital."
              : "Σ prices of production ≠ Σ values."}
          </div>
          <div className="v3-ch12-summary-grid">
            <div className="v3-ch12-stat">
              <span className="v3-ch12-stat-label">General rate of profit</span>
              <span className="v3-ch12-stat-value">{pct(summary.general_rate)}</span>
            </div>
            <div className="v3-ch12-stat">
              <span className="v3-ch12-stat-label">Σ values</span>
              <span className="v3-ch12-stat-value">{summary.total_values}</span>
            </div>
            <div className="v3-ch12-stat">
              <span className="v3-ch12-stat-label">Σ prices of production</span>
              <span className="v3-ch12-stat-value">{summary.total_prices_of_production}</span>
            </div>
            <div className="v3-ch12-stat">
              <span className="v3-ch12-stat-label">Σ surplus-value</span>
              <span className="v3-ch12-stat-value">{summary.average_composition_check.total_surplus_value}</span>
            </div>
            <div className="v3-ch12-stat">
              <span className="v3-ch12-stat-label">Σ average profit</span>
              <span className="v3-ch12-stat-value">{summary.average_composition_check.total_average_profit}</span>
            </div>
            <div className="v3-ch12-stat">
              <span className="v3-ch12-stat-label">s = average profit?</span>
              <span className="v3-ch12-stat-value">
                {summary.average_composition_check.surplus_value_equals_average_profit ? "yes" : "no"}
              </span>
            </div>
          </div>
          <div className="v3-ch12-counts">
            Spheres: {summary.sphere_count} · Prices of production:{" "}
            {summary.prices_of_production_count} · Wage-effect analyses:{" "}
            {summary.wage_effect_count}
          </div>
        </section>
      )}
    </div>
  );
}
