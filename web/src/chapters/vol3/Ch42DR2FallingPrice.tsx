import { useEffect, useState } from "react";
import { financeApi } from "../../api";
import type {
  FallingPriceVariant,
  FallingPriceDR2Outcome,
  NormalCapitalPerAcre,
  SoilExclusionEvent,
} from "../../types";
import "./Ch42DR2FallingPrice.css";

const DEFAULT_PARCEL_ID = "5eed000000000000003700";

function poundDisplay(bp: number): string {
  return `£${(bp / 100).toFixed(2)}`;
}

const GRADE_LABELS: Record<number, string> = {
  1: "A (grade 1 — worst)",
  2: "B (grade 2)",
  3: "C (grade 3)",
  4: "D (grade 4 — best)",
};

const VARIANT_LABELS: Record<FallingPriceVariant, string> = {
  1: "1 — Constant productivity",
  2: "2 — Decreasing productivity",
  3: "3 — Rising productivity",
};

export function Ch42DR2FallingPrice() {
  const [exclusionEvents, setExclusionEvents] = useState<SoilExclusionEvent[]>([]);
  const [outcomes, setOutcomes] = useState<FallingPriceDR2Outcome[]>([]);
  const [normalCapitals, setNormalCapitals] = useState<NormalCapitalPerAcre[]>([]);
  const [listError, setListError] = useState<string | null>(null);

  // Soil exclusion event form
  const [excludedGrade, setExcludedGrade] = useState<number>(1);
  const [newRegulatingGrade, setNewRegulatingGrade] = useState<number>(2);
  const [newPriceBp, setNewPriceBp] = useState<number>(540);
  const [exclusionError, setExclusionError] = useState<string | null>(null);
  const [lastExclusion, setLastExclusion] = useState<SoilExclusionEvent | null>(null);

  // Falling-price DR2 outcome form
  const [variant, setVariant] = useState<FallingPriceVariant>(1);
  const [exclusionEventId, setExclusionEventId] = useState<string>("");
  const [tranche1ParcelId, setTranche1ParcelId] = useState(DEFAULT_PARCEL_ID);
  const [tranche1Capital, setTranche1Capital] = useState(250);
  const [tranche1Output, setTranche1Output] = useState(400);
  const [tranche1Surplus, setTranche1Surplus] = useState(125);
  const [tranche2ParcelId, setTranche2ParcelId] = useState(DEFAULT_PARCEL_ID);
  const [tranche2Capital, setTranche2Capital] = useState(250);
  const [tranche2Output, setTranche2Output] = useState(400);
  const [tranche2Surplus, setTranche2Surplus] = useState(125);
  const [outcomeError, setOutcomeError] = useState<string | null>(null);
  const [lastOutcome, setLastOutcome] = useState<FallingPriceDR2Outcome | null>(null);

  // Normal capital per acre form
  const [periodLabel, setPeriodLabel] = useState("England pre-1848");
  const [capitalPerAcreBp, setCapitalPerAcreBp] = useState(800);
  const [normalCapitalError, setNormalCapitalError] = useState<string | null>(null);
  const [lastNormalCapital, setLastNormalCapital] = useState<NormalCapitalPerAcre | null>(null);

  async function refreshLists() {
    try {
      const [events, outs, normals] = await Promise.all([
        financeApi.listSoilExclusionEvents(),
        financeApi.listFallingPriceDR2Outcomes(),
        financeApi.listNormalCapitalPerAcre(),
      ]);
      setExclusionEvents(events);
      setOutcomes(outs);
      setNormalCapitals(normals);
    } catch (e) {
      setListError(String(e));
    }
  }

  useEffect(() => {
    refreshLists().catch(() => {});
  }, []);

  async function handleCreateExclusionEvent() {
    setExclusionError(null);
    setLastExclusion(null);
    try {
      const result = await financeApi.createSoilExclusionEvent({
        excluded_grade: excludedGrade,
        new_regulating_grade: newRegulatingGrade,
        new_price_of_production_bp: newPriceBp,
      });
      setLastExclusion(result);
      await refreshLists();
    } catch (e) {
      setExclusionError(String(e));
    }
  }

  async function handleCreateOutcome() {
    setOutcomeError(null);
    setLastOutcome(null);
    try {
      const result = await financeApi.createFallingPriceDR2Outcome({
        variant,
        exclusion_event_id: exclusionEventId,
        investments: [
          {
            parcel_id: tranche1ParcelId,
            tranche_number: 1,
            capital_bp: tranche1Capital,
            output_quarters: tranche1Output,
            surplus_profit_bp: tranche1Surplus,
          },
          {
            parcel_id: tranche2ParcelId,
            tranche_number: 2,
            capital_bp: tranche2Capital,
            output_quarters: tranche2Output,
            surplus_profit_bp: tranche2Surplus,
          },
        ],
      });
      setLastOutcome(result);
      await refreshLists();
    } catch (e) {
      setOutcomeError(String(e));
    }
  }

  async function handleCreateNormalCapital() {
    setNormalCapitalError(null);
    setLastNormalCapital(null);
    try {
      const result = await financeApi.createNormalCapitalPerAcre({
        period_label: periodLabel,
        capital_per_acre_bp: capitalPerAcreBp,
      });
      setLastNormalCapital(result);
      await refreshLists();
    } catch (e) {
      setNormalCapitalError(String(e));
    }
  }

  return (
    <div className="v3-ch42">
      <section className="v3-ch42-intro">
        <p>
          Chapter&nbsp;42 examines DR&nbsp;II under a{" "}
          <strong>falling price of production</strong>. When additional capital invested on
          better soils raises their output sufficiently, the worst soil (A) can be driven out
          of cultivation altogether. A better soil then becomes the new regulating soil and the
          market price of production <em>falls</em>.
        </p>
        <p>
          The key result is the{" "}
          <strong>compensation law</strong>: even though the price per quarter falls, the
          grain-rent (in quarters) can double, so that the total money-rent is preserved or
          even rises. Marx&apos;s Table&nbsp;IVd shows that when soil&nbsp;A is excluded and
          soil&nbsp;B becomes the regulator, the grain-rent doubles while the price halves —
          leaving the money-rent at £18 unchanged. Additionally, as cultivation intensifies,
          the normal capital per acre rises (Marx notes England moved from roughly £8 to
          £12&nbsp;/&nbsp;acre between the Napoleonic wars and the mid-nineteenth century).
        </p>
      </section>

      {/* Section 1: Soil exclusion events */}
      <section className="v3-ch42-section" aria-label="Soil exclusion events">
        <h4>1. Soil exclusion events</h4>
        <p className="v3-ch42-subintro">
          Record when additional investment on better soils drives the worst grade out of
          cultivation and establishes a new regulating soil at a lower price of production.
          The new regulating grade must be a higher grade (lower number = worse) than the
          excluded grade. Default price 540&nbsp;bp = 45s/qtr (price falls from 60s to 45s
          when B replaces A as regulator).
        </p>
        <div className="v3-ch42-form">
          <label>
            Excluded grade (driven out of cultivation)
            <select
              value={excludedGrade}
              onChange={(e) => setExcludedGrade(Number(e.target.value))}
            >
              <option value={1}>{GRADE_LABELS[1]}</option>
              <option value={2}>{GRADE_LABELS[2]}</option>
              <option value={3}>{GRADE_LABELS[3]}</option>
              <option value={4}>{GRADE_LABELS[4]}</option>
            </select>
          </label>
          <label>
            New regulating grade (must be a higher grade than excluded)
            <select
              value={newRegulatingGrade}
              onChange={(e) => setNewRegulatingGrade(Number(e.target.value))}
            >
              <option value={1}>{GRADE_LABELS[1]}</option>
              <option value={2}>{GRADE_LABELS[2]}</option>
              <option value={3}>{GRADE_LABELS[3]}</option>
              <option value={4}>{GRADE_LABELS[4]}</option>
            </select>
          </label>
          <label>
            New price of production (basis-points)
            <input
              type="number"
              min={1}
              value={newPriceBp}
              onChange={(e) => setNewPriceBp(Math.max(1, Number(e.target.value)))}
            />
            <span className="v3-ch42-hint">
              100&nbsp;bp = £1. Default 540&nbsp;bp = 45s (Table&nbsp;IVd: A→B, price falls
              from 60s to 45s).
            </span>
          </label>
        </div>
        <button className="v3-ch42-btn" onClick={handleCreateExclusionEvent}>
          Record exclusion event
        </button>
        {exclusionError && <p className="v3-ch42-error">{exclusionError}</p>}
        {lastExclusion && (
          <div className="v3-ch42-preview">
            <strong>
              Saved — grade {lastExclusion.excluded_grade} excluded; new regulator grade{" "}
              {lastExclusion.new_regulating_grade} at{" "}
              {poundDisplay(lastExclusion.new_price_of_production_bp)}/qtr
            </strong>
            <span className="v3-ch42-hint">ID: {lastExclusion.id}</span>
          </div>
        )}
        {exclusionEvents.length > 0 && (
          <div className="v3-ch42-table-wrap">
            <table className="v3-ch42-table">
              <thead>
                <tr>
                  <th>ID</th>
                  <th>Excluded grade</th>
                  <th>New regulating grade</th>
                  <th>New price (£/qtr)</th>
                </tr>
              </thead>
              <tbody>
                {exclusionEvents.map((ev) => (
                  <tr key={ev.id}>
                    <td title={ev.id}>{ev.id.slice(0, 8)}&hellip;</td>
                    <td>{ev.excluded_grade}</td>
                    <td>{ev.new_regulating_grade}</td>
                    <td>{(ev.new_price_of_production_bp / 100).toFixed(2)}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
        {listError && <p className="v3-ch42-error">{listError}</p>}
      </section>

      {/* Section 2: Falling-price DR II outcomes */}
      <section className="v3-ch42-section" aria-label="Falling-price DR II outcomes">
        <h4>2. Falling-price DR&nbsp;II outcomes</h4>
        <p className="v3-ch42-subintro">
          Compute the aggregate rental after the exclusion event. Paste an exclusion event ID
          from the table above. Capital and surplus-profit are in basis-points (100&nbsp;bp =
          £1). The surplus_profit_bp field here represents the grain-rent in 1/100-qtr units
          (grain-rent in hundredths of a quarter). The seed fixture uses variant&nbsp;1
          (constant productivity) at the Table&nbsp;IVd compensation scenario: total
          money-rent £18, grain-rent doubles.
        </p>
        <div className="v3-ch42-form">
          <label>
            Productivity variant
            <select
              value={variant}
              onChange={(e) => setVariant(Number(e.target.value) as FallingPriceVariant)}
            >
              <option value={1}>{VARIANT_LABELS[1]}</option>
              <option value={2}>{VARIANT_LABELS[2]}</option>
              <option value={3}>{VARIANT_LABELS[3]}</option>
            </select>
          </label>
          <label>
            Exclusion event ID
            <input
              value={exclusionEventId}
              onChange={(e) => setExclusionEventId(e.target.value)}
              placeholder="Paste an ID from the exclusion events table above"
            />
          </label>
          <fieldset style={{ border: "1px solid #ddd", borderRadius: 4, padding: "0.75rem" }}>
            <legend style={{ fontSize: "0.8125rem", fontWeight: 600 }}>Tranche 1</legend>
            <label>
              Parcel ID
              <input
                value={tranche1ParcelId}
                onChange={(e) => setTranche1ParcelId(e.target.value)}
                placeholder={DEFAULT_PARCEL_ID}
              />
            </label>
            <label>
              Capital (basis-points)
              <input
                type="number"
                min={1}
                value={tranche1Capital}
                onChange={(e) => setTranche1Capital(Math.max(1, Number(e.target.value)))}
              />
              <span className="v3-ch42-hint">100&nbsp;bp = £1. Default £2.50 = 250&nbsp;bp.</span>
            </label>
            <label>
              Output (quarters)
              <input
                type="number"
                min={1}
                value={tranche1Output}
                onChange={(e) => setTranche1Output(Math.max(1, Number(e.target.value)))}
              />
            </label>
            <label>
              Grain-rent (1/100 qtr)
              <input
                type="number"
                min={0}
                value={tranche1Surplus}
                onChange={(e) => setTranche1Surplus(Math.max(0, Number(e.target.value)))}
              />
              <span className="v3-ch42-hint">
                Grain-rent in hundredths of a quarter. Default 125 = 1.25&nbsp;qtrs.
              </span>
            </label>
          </fieldset>
          <fieldset style={{ border: "1px solid #ddd", borderRadius: 4, padding: "0.75rem" }}>
            <legend style={{ fontSize: "0.8125rem", fontWeight: 600 }}>Tranche 2</legend>
            <label>
              Parcel ID
              <input
                value={tranche2ParcelId}
                onChange={(e) => setTranche2ParcelId(e.target.value)}
                placeholder={DEFAULT_PARCEL_ID}
              />
            </label>
            <label>
              Capital (basis-points)
              <input
                type="number"
                min={1}
                value={tranche2Capital}
                onChange={(e) => setTranche2Capital(Math.max(1, Number(e.target.value)))}
              />
            </label>
            <label>
              Output (quarters)
              <input
                type="number"
                min={1}
                value={tranche2Output}
                onChange={(e) => setTranche2Output(Math.max(1, Number(e.target.value)))}
              />
            </label>
            <label>
              Grain-rent (1/100 qtr)
              <input
                type="number"
                min={0}
                value={tranche2Surplus}
                onChange={(e) => setTranche2Surplus(Math.max(0, Number(e.target.value)))}
              />
            </label>
          </fieldset>
        </div>
        <button className="v3-ch42-btn" onClick={handleCreateOutcome}>
          Compute falling-price DR&nbsp;II outcome
        </button>
        {outcomeError && <p className="v3-ch42-error">{outcomeError}</p>}
        {lastOutcome && (
          <div className="v3-ch42-preview">
            <strong>
              Saved — {VARIANT_LABELS[lastOutcome.variant]}, money-rent{" "}
              {poundDisplay(lastOutcome.total_money_rental_bp)}
            </strong>
            <span>
              Grain-rent: {(lastOutcome.total_grain_rent_qtrs / 100).toFixed(2)}&nbsp;qtrs,
              total output: {lastOutcome.total_output_quarters}&nbsp;qtrs, capital:{" "}
              {poundDisplay(lastOutcome.total_capital_bp)}
            </span>
          </div>
        )}
        {outcomes.length > 0 && (
          <div className="v3-ch42-table-wrap">
            <table className="v3-ch42-table">
              <thead>
                <tr>
                  <th>ID</th>
                  <th>Variant</th>
                  <th>Grain-rent (qtrs)</th>
                  <th>Money-rent (£)</th>
                  <th>Capital (£)</th>
                </tr>
              </thead>
              <tbody>
                {outcomes.map((o) => (
                  <tr key={o.id}>
                    <td title={o.id}>{o.id.slice(0, 8)}&hellip;</td>
                    <td>{VARIANT_LABELS[o.variant]}</td>
                    <td>{(o.total_grain_rent_qtrs / 100).toFixed(2)}</td>
                    <td>{(o.total_money_rental_bp / 100).toFixed(2)}</td>
                    <td>{(o.total_capital_bp / 100).toFixed(2)}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </section>

      {/* Section 3: Normal capital per acre */}
      <section className="v3-ch42-section" aria-label="Normal capital per acre">
        <h4>3. Normal capital per acre</h4>
        <p className="v3-ch42-subintro">
          As cultivation intensifies, the socially normal capital invested per acre rises.
          Marx records that England moved from roughly £8/acre before 1848 to £12/acre
          after, reflecting the general advance of agricultural technique. Record period
          benchmarks here. Capital per acre is in basis-points (100&nbsp;bp = £1).
        </p>
        <div className="v3-ch42-form">
          <label>
            Period label
            <input
              value={periodLabel}
              onChange={(e) => setPeriodLabel(e.target.value)}
              placeholder="England pre-1848"
            />
          </label>
          <label>
            Capital per acre (basis-points)
            <input
              type="number"
              min={1}
              value={capitalPerAcreBp}
              onChange={(e) => setCapitalPerAcreBp(Math.max(1, Number(e.target.value)))}
            />
            <span className="v3-ch42-hint">
              100&nbsp;bp = £1. Default 800&nbsp;bp = £8/acre (pre-1848 England).
            </span>
          </label>
        </div>
        <button className="v3-ch42-btn" onClick={handleCreateNormalCapital}>
          Save normal capital benchmark
        </button>
        {normalCapitalError && <p className="v3-ch42-error">{normalCapitalError}</p>}
        {lastNormalCapital && (
          <div className="v3-ch42-preview">
            Saved &mdash; {lastNormalCapital.period_label}:{" "}
            {poundDisplay(lastNormalCapital.capital_per_acre_bp)}/acre
          </div>
        )}
        {normalCapitals.length > 0 && (
          <div className="v3-ch42-table-wrap">
            <table className="v3-ch42-table">
              <thead>
                <tr>
                  <th>ID</th>
                  <th>Period</th>
                  <th>Capital/acre (£)</th>
                </tr>
              </thead>
              <tbody>
                {normalCapitals.map((n) => (
                  <tr key={n.id}>
                    <td title={n.id}>{n.id.slice(0, 8)}&hellip;</td>
                    <td>{n.period_label}</td>
                    <td>{(n.capital_per_acre_bp / 100).toFixed(2)}</td>
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
