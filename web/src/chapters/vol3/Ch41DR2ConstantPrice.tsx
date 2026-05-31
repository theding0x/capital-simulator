import { useEffect, useState } from "react";
import { financeApi } from "../../api";
import type {
  CapitalProductivityCase,
  DifferentialRentIIOutcome,
  IntensiveExtensiveComparison,
} from "../../types";
import "./Ch41DR2ConstantPrice.css";

const DEFAULT_PARCEL_ID = "5eed000000000000003700";

function poundDisplay(bp: number): string {
  return `£${(bp / 100).toFixed(2)}`;
}

const CASE_LABELS: Record<CapitalProductivityCase, string> = {
  1: "1 — Average-only",
  2: "2 — Proportional",
  3: "3 — Decreasing-rate",
  4: "4 — Increasing-return",
};

export function Ch41DR2ConstantPrice() {
  const [outcomes, setOutcomes] = useState<DifferentialRentIIOutcome[]>([]);
  const [comparisons, setComparisons] = useState<IntensiveExtensiveComparison[]>([]);
  const [listError, setListError] = useState<string | null>(null);

  // DR2 outcome form
  const [outcomeCase, setOutcomeCase] = useState<CapitalProductivityCase>(2);
  const [tranche1ParcelId, setTranche1ParcelId] = useState(DEFAULT_PARCEL_ID);
  const [tranche1Capital, setTranche1Capital] = useState(250);
  const [tranche1Output, setTranche1Output] = useState(400);
  const [tranche1Surplus, setTranche1Surplus] = useState(900);
  const [tranche2ParcelId, setTranche2ParcelId] = useState(DEFAULT_PARCEL_ID);
  const [tranche2Capital, setTranche2Capital] = useState(250);
  const [tranche2Output, setTranche2Output] = useState(400);
  const [tranche2Surplus, setTranche2Surplus] = useState(900);
  const [outcomeError, setOutcomeError] = useState<string | null>(null);
  const [lastOutcome, setLastOutcome] = useState<DifferentialRentIIOutcome | null>(null);

  // Intensive vs extensive form
  const [ieIntensive, setIeIntensive] = useState(1800);
  const [ieExtensive, setIeExtensive] = useState(900);
  const [ieTotalCapital, setIeTotalCapital] = useState(500);
  const [ieTotalRental, setIeTotalRental] = useState(1800);
  const [ieError, setIeError] = useState<string | null>(null);
  const [lastIe, setLastIe] = useState<IntensiveExtensiveComparison | null>(null);

  async function refreshLists() {
    try {
      const [outs, comps] = await Promise.all([
        financeApi.listDR2Outcomes(),
        financeApi.listIntensiveExtensive(),
      ]);
      setOutcomes(outs);
      setComparisons(comps);
    } catch (e) {
      setListError(String(e));
    }
  }

  useEffect(() => {
    refreshLists().catch(() => {});
  }, []);

  async function handleCreateOutcome() {
    setOutcomeError(null);
    setLastOutcome(null);
    try {
      const result = await financeApi.createDR2Outcome({
        case: outcomeCase,
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

  async function handleCreateComparison() {
    setIeError(null);
    setLastIe(null);
    try {
      const result = await financeApi.createIntensiveExtensive({
        intensive_rent_per_acre_bp: ieIntensive,
        extensive_rent_per_acre_bp: ieExtensive,
        total_capital_bp: ieTotalCapital,
        total_rental_same_bp: ieTotalRental,
      });
      setLastIe(result);
      await refreshLists();
    } catch (e) {
      setIeError(String(e));
    }
  }

  return (
    <div className="v3-ch41">
      <section className="v3-ch41-intro">
        <p>
          Chapter&nbsp;41 examines DR&nbsp;II under a{" "}
          <strong>constant regulating price of production</strong> — soil&nbsp;A remains the
          worst soil and sets the market price throughout. Under this condition, every
          additional investment of capital on the better soils (B, C, D) always raises the
          absolute rent per acre, regardless of whether the productivity of the added capital
          is proportional to the first investment (Case&nbsp;II), falls short of it
          (Case&nbsp;III), or exceeds it (Case&nbsp;IV).
        </p>
        <p>
          The critical consequence is that{" "}
          <strong>concentrating additional capital on a smaller area</strong> raises rent per
          acre more than dispersing the same total capital across a larger area at equal
          tranche size (intensive vs extensive comparison). The market price is regulated by
          the worst conditions; any surplus on better soil flows to the landowner as rent —
          whether the added capital is more or less productive than the first.
        </p>
      </section>

      {/* Section 1: DR II outcomes */}
      <section className="v3-ch41-section" aria-label="DR II outcome by case">
        <h4>1. DR&nbsp;II outcomes (constant price)</h4>
        <p className="v3-ch41-subintro">
          Select the productivity case for the second investment tranche and enter two tranche
          rows. Capital and surplus-profit are in basis-points (100&nbsp;bp = £1). The seed
          fixture uses Case&nbsp;II (proportional): two £2.50 tranches each yielding
          £9&nbsp;surplus-profit → total £18 at a rate of 36&nbsp;000&nbsp;bp.
        </p>
        <div className="v3-ch41-form">
          <label>
            Productivity case
            <select
              value={outcomeCase}
              onChange={(e) => setOutcomeCase(Number(e.target.value) as CapitalProductivityCase)}
            >
              <option value={1}>{CASE_LABELS[1]}</option>
              <option value={2}>{CASE_LABELS[2]}</option>
              <option value={3}>{CASE_LABELS[3]}</option>
              <option value={4}>{CASE_LABELS[4]}</option>
            </select>
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
              <span className="v3-ch41-hint">100&nbsp;bp = £1. Default £2.50 = 250&nbsp;bp.</span>
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
              Surplus-profit (basis-points)
              <input
                type="number"
                min={0}
                value={tranche1Surplus}
                onChange={(e) => setTranche1Surplus(Math.max(0, Number(e.target.value)))}
              />
              <span className="v3-ch41-hint">Default £9 = 900&nbsp;bp.</span>
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
              Surplus-profit (basis-points)
              <input
                type="number"
                min={0}
                value={tranche2Surplus}
                onChange={(e) => setTranche2Surplus(Math.max(0, Number(e.target.value)))}
              />
            </label>
          </fieldset>
        </div>
        <button className="v3-ch41-btn" onClick={handleCreateOutcome}>
          Compute DR&nbsp;II outcome
        </button>
        {outcomeError && <p className="v3-ch41-error">{outcomeError}</p>}
        {lastOutcome && (
          <div className="v3-ch41-preview">
            <strong>
              Saved — {CASE_LABELS[lastOutcome.case]}, total rental{" "}
              {poundDisplay(lastOutcome.total_rental_bp)}
            </strong>
            <span>
              Additional capital: {poundDisplay(lastOutcome.additional_capital_bp)}, new rent/acre:{" "}
              {poundDisplay(lastOutcome.new_rent_per_acre_bp)}, rate of surplus-profit:{" "}
              {(lastOutcome.rate_of_surplus_profit_bp / 100).toFixed(0)}%
            </span>
          </div>
        )}
        {outcomes.length > 0 && (
          <div className="v3-ch41-table-wrap">
            <table className="v3-ch41-table">
              <thead>
                <tr>
                  <th>ID</th>
                  <th>Case</th>
                  <th>Addl. capital (£)</th>
                  <th>New rent/acre (£)</th>
                  <th>Total rental (£)</th>
                  <th>Rate of surplus-profit (%)</th>
                </tr>
              </thead>
              <tbody>
                {outcomes.map((o) => (
                  <tr key={o.id}>
                    <td title={o.id}>{o.id.slice(0, 8)}&hellip;</td>
                    <td>{CASE_LABELS[o.case]}</td>
                    <td>{(o.additional_capital_bp / 100).toFixed(2)}</td>
                    <td>{(o.new_rent_per_acre_bp / 100).toFixed(2)}</td>
                    <td>{(o.total_rental_bp / 100).toFixed(2)}</td>
                    <td>{(o.rate_of_surplus_profit_bp / 100).toFixed(0)}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
        {listError && <p className="v3-ch41-error">{listError}</p>}
      </section>

      {/* Section 2: Intensive vs extensive */}
      <section className="v3-ch41-section" aria-label="Intensive vs extensive comparison">
        <h4>2. Intensive vs extensive comparison</h4>
        <p className="v3-ch41-subintro">
          Concentrating additional capital on the same (better) parcel raises rent per acre
          more than spreading the same total capital across a wider area. Enter rent-per-acre
          figures for each mode (basis-points, 100&nbsp;bp = £1) and submit to record the
          comparison. The seed fixture: intensive £18/acre (1800&nbsp;bp) vs extensive
          £9/acre (900&nbsp;bp) at equal total capital £5 (500&nbsp;bp) — same total rental
          of £18.
        </p>
        <div className="v3-ch41-form">
          <label>
            Intensive rent per acre (basis-points)
            <input
              type="number"
              min={0}
              value={ieIntensive}
              onChange={(e) => setIeIntensive(Math.max(0, Number(e.target.value)))}
            />
            <span className="v3-ch41-hint">Default £18 = 1800&nbsp;bp.</span>
          </label>
          <label>
            Extensive rent per acre (basis-points)
            <input
              type="number"
              min={0}
              value={ieExtensive}
              onChange={(e) => setIeExtensive(Math.max(0, Number(e.target.value)))}
            />
            <span className="v3-ch41-hint">Default £9 = 900&nbsp;bp.</span>
          </label>
          <label>
            Total capital (basis-points)
            <input
              type="number"
              min={1}
              value={ieTotalCapital}
              onChange={(e) => setIeTotalCapital(Math.max(1, Number(e.target.value)))}
            />
          </label>
          <label>
            Total rental — same in both modes (basis-points)
            <input
              type="number"
              min={0}
              value={ieTotalRental}
              onChange={(e) => setIeTotalRental(Math.max(0, Number(e.target.value)))}
            />
            <span className="v3-ch41-hint">Default £18 = 1800&nbsp;bp.</span>
          </label>
        </div>
        <button className="v3-ch41-btn" onClick={handleCreateComparison}>
          Save comparison
        </button>
        {ieError && <p className="v3-ch41-error">{ieError}</p>}
        {lastIe && (
          <div className="v3-ch41-preview">
            Saved &mdash; intensive {poundDisplay(lastIe.intensive_rent_per_acre_bp)}/acre vs
            extensive {poundDisplay(lastIe.extensive_rent_per_acre_bp)}/acre; total rental{" "}
            {poundDisplay(lastIe.total_rental_same_bp)}
          </div>
        )}
        {comparisons.length > 0 && (
          <div className="v3-ch41-table-wrap">
            <table className="v3-ch41-table">
              <thead>
                <tr>
                  <th>ID</th>
                  <th>Intensive rent/acre (£)</th>
                  <th>Extensive rent/acre (£)</th>
                  <th>Total capital (£)</th>
                  <th>Total rental (£)</th>
                </tr>
              </thead>
              <tbody>
                {comparisons.map((c) => (
                  <tr key={c.id}>
                    <td title={c.id}>{c.id.slice(0, 8)}&hellip;</td>
                    <td>{(c.intensive_rent_per_acre_bp / 100).toFixed(2)}</td>
                    <td>{(c.extensive_rent_per_acre_bp / 100).toFixed(2)}</td>
                    <td>{(c.total_capital_bp / 100).toFixed(2)}</td>
                    <td>{(c.total_rental_same_bp / 100).toFixed(2)}</td>
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
