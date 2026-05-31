import { useEffect, useState } from "react";
import { financeApi } from "../../api";
import type {
  RentEmergenceMechanism,
  RentOnWorstSoil,
  SoilImprovementCapital,
  RegulatingPriceShift,
} from "../../types";
import "./Ch44WorstSoilRent.css";

const MECHANISM_LABELS: Record<RentEmergenceMechanism, string> = {
  1: "Mechanism 1 — better soil's sub-average investment becomes new price-regulator",
  2: "Mechanism 2 — additional capital raises productivity on A, lowering regulating price",
  3: "Mechanism 3 — additional capital yields decreasing returns, raising A's average price",
};

const GRADE_LABELS: Record<number, string> = {
  1: "A",
  2: "B",
  3: "C",
  4: "D",
};

function formatPounds(bp: number): string {
  return `£${(bp / 10000).toFixed(4).replace(/\.?0+$/, "")}`;
}

export function Ch44WorstSoilRent() {
  const [worstSoilRents, setWorstSoilRents] = useState<RentOnWorstSoil[]>([]);
  const [improvements, setImprovements] = useState<SoilImprovementCapital[]>([]);
  const [priceShifts, setPriceShifts] = useState<RegulatingPriceShift[]>([]);
  const [listError, setListError] = useState<string | null>(null);

  // Rent on worst soil form
  const [wsParceldId, setWsParcelId] = useState("5eed000000000000003700");
  const [wsMechanism, setWsMechanism] = useState<RentEmergenceMechanism>(1);
  const [wsOldPrice, setWsOldPrice] = useState(30000);
  const [wsNewPrice, setWsNewPrice] = useState(35000);
  const [wsError, setWsError] = useState<string | null>(null);
  const [lastWs, setLastWs] = useState<RentOnWorstSoil | null>(null);

  // Soil improvement capital form
  const [siParcelId, setSiParcelId] = useState("5eed000000000000003700");
  const [siCapital, setSiCapital] = useState(3000);
  const [siGain, setSiGain] = useState(500);
  const [siAmortYears, setSiAmortYears] = useState(15);
  const [siAmortised, setSiAmortised] = useState(0);
  const [siError, setSiError] = useState<string | null>(null);
  const [lastSi, setLastSi] = useState<SoilImprovementCapital | null>(null);

  // Regulating price shift form
  const [rpsFromGrade, setRpsFromGrade] = useState(1);
  const [rpsToGrade, setRpsToGrade] = useState(2);
  const [rpsOldPrice, setRpsOldPrice] = useState(30000);
  const [rpsNewPrice, setRpsNewPrice] = useState(35000);
  const [rpsError, setRpsError] = useState<string | null>(null);
  const [lastRps, setLastRps] = useState<RegulatingPriceShift | null>(null);

  async function refreshLists() {
    try {
      const [rents, imps, shifts] = await Promise.all([
        financeApi.listRentOnWorstSoils(),
        financeApi.listSoilImprovementCapitals(),
        financeApi.listRegulatingPriceShifts(),
      ]);
      setWorstSoilRents(rents);
      setImprovements(imps);
      setPriceShifts(shifts);
    } catch (e) {
      setListError(String(e));
    }
  }

  useEffect(() => {
    refreshLists().catch(() => {});
  }, []);

  async function handleCreateWorstSoilRent() {
    setWsError(null);
    setLastWs(null);
    try {
      const result = await financeApi.createRentOnWorstSoil({
        parcel_id: wsParceldId,
        mechanism: wsMechanism,
        old_price_of_production_bp: wsOldPrice,
        new_price_of_production_bp: wsNewPrice,
      });
      setLastWs(result);
      await refreshLists();
    } catch (e) {
      setWsError(String(e));
    }
  }

  async function handleCreateSoilImprovement() {
    setSiError(null);
    setLastSi(null);
    try {
      const result = await financeApi.createSoilImprovementCapital({
        parcel_id: siParcelId,
        capital_invested_bp: siCapital,
        productivity_gain_bp: siGain,
        amortisation_years: siAmortYears,
        amortised_years: siAmortised,
      });
      setLastSi(result);
      await refreshLists();
    } catch (e) {
      setSiError(String(e));
    }
  }

  async function handleCreatePriceShift() {
    setRpsError(null);
    setLastRps(null);
    try {
      const result = await financeApi.createRegulatingPriceShift({
        from_grade: rpsFromGrade,
        to_grade: rpsToGrade,
        old_price_bp: rpsOldPrice,
        new_price_bp: rpsNewPrice,
      });
      setLastRps(result);
      await refreshLists();
    } catch (e) {
      setRpsError(String(e));
    }
  }

  return (
    <div className="v3-ch44">
      <section className="v3-ch44-intro">
        <p>
          Chapter&nbsp;44 addresses a paradox: how can soil&nbsp;A — defined as the worst,
          rent-free, regulating soil — ever yield a rent? Marx identifies three mechanisms
          by which differential rent can emerge even on the hitherto rentless worst soil.
        </p>
        <p>
          <strong>Mechanism&nbsp;1</strong>: an additional investment on a{" "}
          <em>better</em> soil yields below the average return, becoming the new
          price-regulator and opening a gap above soil&nbsp;A's price of production.{" "}
          <strong>Mechanism&nbsp;2</strong>: increasing productivity of additional capital
          on A itself lowers the regulating price, so A now yields a surplus-profit.{" "}
          <strong>Mechanism&nbsp;3</strong>: decreasing productivity of additional capital
          raises A's average price of production until it exceeds the old regulating price
          and a rent appears. Permanent soil improvements such as drainage yield rent as
          interest on the incorporated capital until the term is amortised, whereupon the
          improvement dissolves into pure differential rent.
        </p>
        <p className="v3-ch44-note">
          Formula: rent&nbsp;per&nbsp;acre&nbsp;= max(0,&nbsp;new&nbsp;&minus;&nbsp;old price of
          production). Price scale: £1&nbsp;=&nbsp;10,000&nbsp;bp.
        </p>
      </section>

      {/* Section 1 — Rent on worst soil */}
      <section className="v3-ch44-section" aria-label="Rent on worst soil">
        <h4>1. Rent on worst soil</h4>
        <p className="v3-ch44-subintro">
          Record the emergence of rent on soil&nbsp;A. Choose the mechanism, enter the old
          and new prices of production (basis points, £1&nbsp;=&nbsp;10,000&nbsp;bp). The
          server computes <code>rent_per_acre_bp&nbsp;= max(0,&nbsp;new&nbsp;&minus;&nbsp;old)</code>.
          Seed fixture: mechanism&nbsp;1, old&nbsp;=&nbsp;£3 (30,000&nbsp;bp),
          new&nbsp;=&nbsp;£3.50 (35,000&nbsp;bp), rent&nbsp;=&nbsp;£0.50.
        </p>
        <div className="v3-ch44-form">
          <label>
            Parcel ID
            <input
              value={wsParceldId}
              onChange={(e) => setWsParcelId(e.target.value)}
              placeholder="5eed000000000000003700"
            />
          </label>
          <label>
            Mechanism
            <select
              value={wsMechanism}
              onChange={(e) => setWsMechanism(Number(e.target.value) as RentEmergenceMechanism)}
            >
              <option value={1}>{MECHANISM_LABELS[1]}</option>
              <option value={2}>{MECHANISM_LABELS[2]}</option>
              <option value={3}>{MECHANISM_LABELS[3]}</option>
            </select>
          </label>
          <label>
            Old price of production (bp)
            <input
              type="number"
              min={0}
              value={wsOldPrice}
              onChange={(e) => setWsOldPrice(Math.max(0, Number(e.target.value)))}
            />
            <span className="v3-ch44-hint">{formatPounds(wsOldPrice)}</span>
          </label>
          <label>
            New price of production (bp)
            <input
              type="number"
              min={0}
              value={wsNewPrice}
              onChange={(e) => setWsNewPrice(Math.max(0, Number(e.target.value)))}
            />
            <span className="v3-ch44-hint">{formatPounds(wsNewPrice)}</span>
          </label>
        </div>
        <button className="v3-ch44-btn" onClick={handleCreateWorstSoilRent}>
          Record rent on worst soil
        </button>
        {wsError && <p className="v3-ch44-error">{wsError}</p>}
        {lastWs && (
          <div className="v3-ch44-preview">
            <strong>
              Saved — {MECHANISM_LABELS[lastWs.mechanism]}: old&nbsp;
              {formatPounds(lastWs.old_price_of_production_bp)}, new&nbsp;
              {formatPounds(lastWs.new_price_of_production_bp)}, rent&nbsp;
              {formatPounds(lastWs.rent_per_acre_bp)}/acre
            </strong>
            <span className="v3-ch44-hint">ID: {lastWs.id}</span>
          </div>
        )}
        {worstSoilRents.length > 0 && (
          <div className="v3-ch44-table-wrap">
            <table className="v3-ch44-table">
              <thead>
                <tr>
                  <th>ID</th>
                  <th>Mechanism</th>
                  <th>Old price</th>
                  <th>New price</th>
                  <th>Rent/acre</th>
                </tr>
              </thead>
              <tbody>
                {worstSoilRents.map((r) => (
                  <tr key={r.id}>
                    <td title={r.id}>{r.id.slice(0, 8)}&hellip;</td>
                    <td>{MECHANISM_LABELS[r.mechanism]}</td>
                    <td>{formatPounds(r.old_price_of_production_bp)}</td>
                    <td>{formatPounds(r.new_price_of_production_bp)}</td>
                    <td>{formatPounds(r.rent_per_acre_bp)}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
        {listError && <p className="v3-ch44-error">{listError}</p>}
      </section>

      {/* Section 2 — Soil improvement capital */}
      <section className="v3-ch44-section" aria-label="Soil improvement capital">
        <h4>2. Soil improvement capital</h4>
        <p className="v3-ch44-subintro">
          Permanent improvements (drainage, clearing) incorporate capital into the soil.
          The capitalist charges interest on this capital as rent until it is amortised.
          Once <code>amortised_years &ge; amortisation_years</code> the improvement
          dissolves into the land's natural fertility and the extra rent becomes pure DR.
          Seed fixture: £0.30 invested, £0.05 productivity gain, 15-year term.
        </p>
        <div className="v3-ch44-form">
          <label>
            Parcel ID
            <input
              value={siParcelId}
              onChange={(e) => setSiParcelId(e.target.value)}
              placeholder="5eed000000000000003700"
            />
          </label>
          <label>
            Capital invested (bp)
            <input
              type="number"
              min={0}
              value={siCapital}
              onChange={(e) => setSiCapital(Math.max(0, Number(e.target.value)))}
            />
            <span className="v3-ch44-hint">{formatPounds(siCapital)}</span>
          </label>
          <label>
            Productivity gain (bp)
            <input
              type="number"
              min={0}
              value={siGain}
              onChange={(e) => setSiGain(Math.max(0, Number(e.target.value)))}
            />
            <span className="v3-ch44-hint">{formatPounds(siGain)}</span>
          </label>
          <label>
            Amortisation years
            <input
              type="number"
              min={1}
              value={siAmortYears}
              onChange={(e) => setSiAmortYears(Math.max(1, Number(e.target.value)))}
            />
          </label>
          <label>
            Years already amortised
            <input
              type="number"
              min={0}
              value={siAmortised}
              onChange={(e) => setSiAmortised(Math.max(0, Number(e.target.value)))}
            />
          </label>
        </div>
        <button className="v3-ch44-btn" onClick={handleCreateSoilImprovement}>
          Record soil improvement
        </button>
        {siError && <p className="v3-ch44-error">{siError}</p>}
        {lastSi && (
          <div className="v3-ch44-preview">
            <strong>
              Saved — capital&nbsp;{formatPounds(lastSi.capital_invested_bp)}, gain&nbsp;
              {formatPounds(lastSi.productivity_gain_bp)}, term&nbsp;{lastSi.amortisation_years}&nbsp;yr
              {lastSi.amortised_years >= lastSi.amortisation_years
                ? " — fully amortised (dissolved into DR)"
                : ` — ${lastSi.amortised_years}/${lastSi.amortisation_years} yr`}
            </strong>
            <span className="v3-ch44-hint">ID: {lastSi.id}</span>
          </div>
        )}
        {improvements.length > 0 && (
          <div className="v3-ch44-table-wrap">
            <table className="v3-ch44-table">
              <thead>
                <tr>
                  <th>ID</th>
                  <th>Capital</th>
                  <th>Gain</th>
                  <th>Term (yr)</th>
                  <th>Amortised (yr)</th>
                  <th>Status</th>
                </tr>
              </thead>
              <tbody>
                {improvements.map((i) => (
                  <tr key={i.id}>
                    <td title={i.id}>{i.id.slice(0, 8)}&hellip;</td>
                    <td>{formatPounds(i.capital_invested_bp)}</td>
                    <td>{formatPounds(i.productivity_gain_bp)}</td>
                    <td>{i.amortisation_years}</td>
                    <td>{i.amortised_years}</td>
                    <td>
                      {i.amortised_years >= i.amortisation_years
                        ? "Dissolved into DR"
                        : "Active"}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </section>

      {/* Section 3 — Regulating price shifts */}
      <section className="v3-ch44-section" aria-label="Regulating price shifts">
        <h4>3. Regulating price shifts</h4>
        <p className="v3-ch44-subintro">
          A shift in which soil grade acts as the regulating (price-setting) soil changes
          the baseline for all differential rents. Record the from-grade and to-grade and
          the old and new regulating prices. Seed fixture: grade&nbsp;A&nbsp;&rarr;&nbsp;B,
          old&nbsp;£3&nbsp;&rarr;&nbsp;new&nbsp;£3.50.
        </p>
        <div className="v3-ch44-form">
          <label>
            From grade (regulating soil becomes non-regulating)
            <select
              value={rpsFromGrade}
              onChange={(e) => setRpsFromGrade(Number(e.target.value))}
            >
              {[1, 2, 3, 4].map((g) => (
                <option key={g} value={g}>
                  {GRADE_LABELS[g]}
                </option>
              ))}
            </select>
          </label>
          <label>
            To grade (new regulating soil)
            <select
              value={rpsToGrade}
              onChange={(e) => setRpsToGrade(Number(e.target.value))}
            >
              {[1, 2, 3, 4].map((g) => (
                <option key={g} value={g}>
                  {GRADE_LABELS[g]}
                </option>
              ))}
            </select>
          </label>
          <label>
            Old regulating price (bp)
            <input
              type="number"
              min={0}
              value={rpsOldPrice}
              onChange={(e) => setRpsOldPrice(Math.max(0, Number(e.target.value)))}
            />
            <span className="v3-ch44-hint">{formatPounds(rpsOldPrice)}</span>
          </label>
          <label>
            New regulating price (bp)
            <input
              type="number"
              min={0}
              value={rpsNewPrice}
              onChange={(e) => setRpsNewPrice(Math.max(0, Number(e.target.value)))}
            />
            <span className="v3-ch44-hint">{formatPounds(rpsNewPrice)}</span>
          </label>
        </div>
        <button className="v3-ch44-btn" onClick={handleCreatePriceShift}>
          Record regulating price shift
        </button>
        {rpsError && <p className="v3-ch44-error">{rpsError}</p>}
        {lastRps && (
          <div className="v3-ch44-preview">
            <strong>
              Saved — grade&nbsp;{GRADE_LABELS[lastRps.from_grade]}&nbsp;&rarr;&nbsp;
              {GRADE_LABELS[lastRps.to_grade]}: {formatPounds(lastRps.old_price_bp)}&nbsp;&rarr;&nbsp;
              {formatPounds(lastRps.new_price_bp)}
            </strong>
            <span className="v3-ch44-hint">ID: {lastRps.id}</span>
          </div>
        )}
        {priceShifts.length > 0 && (
          <div className="v3-ch44-table-wrap">
            <table className="v3-ch44-table">
              <thead>
                <tr>
                  <th>ID</th>
                  <th>From grade</th>
                  <th>To grade</th>
                  <th>Old price</th>
                  <th>New price</th>
                </tr>
              </thead>
              <tbody>
                {priceShifts.map((s) => (
                  <tr key={s.id}>
                    <td title={s.id}>{s.id.slice(0, 8)}&hellip;</td>
                    <td>{GRADE_LABELS[s.from_grade] ?? s.from_grade}</td>
                    <td>{GRADE_LABELS[s.to_grade] ?? s.to_grade}</td>
                    <td>{formatPounds(s.old_price_bp)}</td>
                    <td>{formatPounds(s.new_price_bp)}</td>
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
