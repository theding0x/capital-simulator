import { useEffect, useState } from "react";
import { financeApi } from "../../api";
import type {
  PriceOfProductionSurplusProfit,
  MonopolisedNaturalForce,
  DifferentialRentRecord,
  CapitalisedRentPrice,
} from "../../types";
import "./Ch38GeneralRemarks.css";

export function Ch38GeneralRemarks() {
  const [surplusProfits, setSurplusProfits] = useState<PriceOfProductionSurplusProfit[]>([]);
  const [naturalForces, setNaturalForces] = useState<MonopolisedNaturalForce[]>([]);
  const [diffRents, setDiffRents] = useState<DifferentialRentRecord[]>([]);
  const [capitalisedPrices, setCapitalisedPrices] = useState<CapitalisedRentPrice[]>([]);
  const [listError, setListError] = useState<string | null>(null);

  // Surplus-profit form
  const [spCapitalId, setSpCapitalId] = useState("Waterfall Mill, Lancashire");
  const [spGeneralBp, setSpGeneralBp] = useState(10000);
  const [spIndividualBp, setSpIndividualBp] = useState(8000);
  const [spError, setSpError] = useState<string | null>(null);
  const [lastSp, setLastSp] = useState<PriceOfProductionSurplusProfit | null>(null);

  // Monopolised natural force form
  const [mnfKind, setMnfKind] = useState("waterfall");
  const [mnfIsMonopolisable, setMnfIsMonopolisable] = useState(true);
  const [mnfParcelId, setMnfParcelId] = useState("");
  const [mnfError, setMnfError] = useState<string | null>(null);
  const [lastMnf, setLastMnf] = useState<MonopolisedNaturalForce | null>(null);

  // Differential rent form
  const [drParcelId, setDrParcelId] = useState("");
  const [drSurplusProfitBp, setDrSurplusProfitBp] = useState(2000);
  const [drAnnualRentLm, setDrAnnualRentLm] = useState(600);
  const [drError, setDrError] = useState<string | null>(null);
  const [lastDr, setLastDr] = useState<DifferentialRentRecord | null>(null);

  // Capitalised rent price form
  const [crpParcelId, setCrpParcelId] = useState("");
  const [crpAnnualRentLm, setCrpAnnualRentLm] = useState(600);
  const [crpInterestRateBp, setCrpInterestRateBp] = useState(500);
  const [crpError, setCrpError] = useState<string | null>(null);
  const [lastCrp, setLastCrp] = useState<CapitalisedRentPrice | null>(null);

  async function refreshLists() {
    try {
      const [sp, mnf, dr, crp] = await Promise.all([
        financeApi.listPoPSurplusProfits(),
        financeApi.listMonopolisedNaturalForces(),
        financeApi.listDifferentialRents(),
        financeApi.listCapitalisedRentPrices(),
      ]);
      setSurplusProfits(sp);
      setNaturalForces(mnf);
      setDiffRents(dr);
      setCapitalisedPrices(crp);
    } catch (e) {
      setListError(String(e));
    }
  }

  useEffect(() => {
    refreshLists().catch(() => {});
  }, []);

  async function handleCreateSurplusProfit() {
    setSpError(null);
    setLastSp(null);
    try {
      const result = await financeApi.createPoPSurplusProfit({
        capital_id: spCapitalId,
        general_price_of_production_bp: spGeneralBp,
        individual_price_of_production_bp: spIndividualBp,
      });
      setLastSp(result);
      await refreshLists();
    } catch (e) {
      setSpError(String(e));
    }
  }

  async function handleCreateNaturalForce() {
    setMnfError(null);
    setLastMnf(null);
    try {
      const result = await financeApi.createMonopolisedNaturalForce({
        kind: mnfKind,
        is_monopolisable: mnfIsMonopolisable,
        parcel_id: mnfParcelId,
      });
      setLastMnf(result);
      await refreshLists();
    } catch (e) {
      setMnfError(String(e));
    }
  }

  async function handleCreateDifferentialRent() {
    setDrError(null);
    setLastDr(null);
    try {
      const result = await financeApi.createDifferentialRent({
        parcel_id: drParcelId,
        surplus_profit_bp: drSurplusProfitBp,
        annual_rent_labour_minutes: drAnnualRentLm,
      });
      setLastDr(result);
      await refreshLists();
    } catch (e) {
      setDrError(String(e));
    }
  }

  async function handleCreateCapitalisedRentPrice() {
    setCrpError(null);
    setLastCrp(null);
    try {
      const result = await financeApi.createCapitalisedRentPrice({
        parcel_id: crpParcelId,
        annual_rent_labour_minutes: crpAnnualRentLm,
        interest_rate_bp: crpInterestRateBp,
      });
      setLastCrp(result);
      await refreshLists();
    } catch (e) {
      setCrpError(String(e));
    }
  }

  return (
    <div className="v3-ch38">
      <section className="v3-ch38-intro">
        <p>
          Chapter&nbsp;38 introduces the <strong>general conditions of
          differential rent</strong> through the waterfall analogy. A producer
          who commands a <strong>free natural force</strong> — a waterfall
          driving a mill — has a lower <em>individual</em> price of production
          than the market-regulating <em>general</em> price. The gap is a{" "}
          <strong>surplus-profit</strong>: extra value over and above the
          average profit. As long as the capitalist monopolises the natural
          force that surplus-profit accrues to him. Once a landowner holds
          title to the soil or waterfall, he converts that surplus-profit into{" "}
          <strong>differential rent</strong> — a permanent deduction from
          the capitalist's profit. The <strong>price of land</strong> is simply
          capitalised rent: at 5&nbsp;% interest, £10/yr rent = £200 purchase
          price (or 96&nbsp;000 labour-minutes at the canonical rate).
        </p>
        <p>
          Differential rent requires only two things: (1) an individual price
          of production below the general (market-regulating) price, and
          (2) landed property converting that gap into rent. It is therefore
          entirely consistent with the law of value — no extra value is
          created; the surplus-profit is merely redirected from capitalist
          to landowner.
        </p>
      </section>

      {/* Section 1: Surplus-Profit */}
      <section className="v3-ch38-section" aria-label="Record a surplus-profit">
        <h4>1. Surplus-profit (price-of-production gap)</h4>
        <p className="v3-ch38-subintro">
          Record a capital whose individual price of production is below the
          general (market-regulating) price. The surplus-profit is the gap
          = general&nbsp;&minus;&nbsp;individual, server-computed.
        </p>
        <div className="v3-ch38-form">
          <label>
            Capital ID
            <input
              value={spCapitalId}
              onChange={(e) => setSpCapitalId(e.target.value)}
              placeholder="e.g. Waterfall Mill, Lancashire"
            />
          </label>
          <label>
            General price of production (basis-points)
            <input
              type="number"
              min={0}
              value={spGeneralBp}
              onChange={(e) => setSpGeneralBp(Math.max(0, Number(e.target.value)))}
            />
            <span className="v3-ch38-hint">Market-regulating price; set by the worst (costliest) producer</span>
          </label>
          <label>
            Individual price of production (basis-points)
            <input
              type="number"
              min={0}
              value={spIndividualBp}
              onChange={(e) => setSpIndividualBp(Math.max(0, Number(e.target.value)))}
            />
            <span className="v3-ch38-hint">This capital's own cost; must be &lt; general to yield a surplus-profit</span>
          </label>
        </div>
        <button className="v3-ch38-btn" onClick={handleCreateSurplusProfit}>
          Save surplus-profit record
        </button>
        {spError && <p className="v3-ch38-error">{spError}</p>}
        {lastSp && (
          <div className="v3-ch38-preview">
            Saved &mdash; surplus-profit: {lastSp.surplus_profit_bp.toLocaleString()}&nbsp;bp
            (general {lastSp.general_price_of_production_bp.toLocaleString()} &minus; individual{" "}
            {lastSp.individual_price_of_production_bp.toLocaleString()})
          </div>
        )}
      </section>

      {/* Surplus-profit list */}
      <section className="v3-ch38-section" aria-label="Stored surplus-profit records">
        <h4>Surplus-profit records</h4>
        {listError && <p className="v3-ch38-error">{listError}</p>}
        {surplusProfits.length === 0 ? (
          <p className="v3-ch38-subintro">No records stored yet.</p>
        ) : (
          <div className="v3-ch38-table-wrap">
            <table className="v3-ch38-table">
              <thead>
                <tr>
                  <th>ID</th>
                  <th>Capital</th>
                  <th>General (bp)</th>
                  <th>Individual (bp)</th>
                  <th>Surplus-profit (bp)</th>
                </tr>
              </thead>
              <tbody>
                {surplusProfits.map((r) => (
                  <tr key={r.id}>
                    <td title={r.id}>{r.id.slice(0, 8)}&hellip;</td>
                    <td>{r.capital_id}</td>
                    <td>{r.general_price_of_production_bp.toLocaleString()}</td>
                    <td>{r.individual_price_of_production_bp.toLocaleString()}</td>
                    <td>{r.surplus_profit_bp.toLocaleString()}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </section>

      {/* Section 2: Monopolised Natural Forces */}
      <section className="v3-ch38-section" aria-label="Record a monopolised natural force">
        <h4>2. Monopolised natural force</h4>
        <p className="v3-ch38-subintro">
          Register a natural force (waterfall, fertile soil, mineral spring)
          that is monopolisable by a landowner and underlies a surplus-profit.
        </p>
        <div className="v3-ch38-form">
          <label>
            Kind
            <input
              value={mnfKind}
              onChange={(e) => setMnfKind(e.target.value)}
              placeholder="e.g. waterfall"
            />
          </label>
          <label>
            <span>Is monopolisable</span>
            <input
              type="checkbox"
              checked={mnfIsMonopolisable}
              onChange={(e) => setMnfIsMonopolisable(e.target.checked)}
              style={{ alignSelf: "flex-start", width: "auto" }}
            />
          </label>
          <label>
            Parcel ID
            <input
              value={mnfParcelId}
              onChange={(e) => setMnfParcelId(e.target.value)}
              placeholder="paste a parcel ID from Ch. 37"
            />
          </label>
        </div>
        <button className="v3-ch38-btn" onClick={handleCreateNaturalForce}>
          Save natural force
        </button>
        {mnfError && <p className="v3-ch38-error">{mnfError}</p>}
        {lastMnf && (
          <div className="v3-ch38-preview">
            Saved &mdash; {lastMnf.kind}{lastMnf.is_monopolisable ? " (monopolisable)" : " (not monopolisable)"}
          </div>
        )}
      </section>

      {/* Natural forces list */}
      <section className="v3-ch38-section" aria-label="Stored monopolised natural forces">
        <h4>Monopolised natural forces</h4>
        {naturalForces.length === 0 ? (
          <p className="v3-ch38-subintro">No natural forces stored yet.</p>
        ) : (
          <div className="v3-ch38-table-wrap">
            <table className="v3-ch38-table">
              <thead>
                <tr>
                  <th>ID</th>
                  <th>Kind</th>
                  <th>Monopolisable</th>
                  <th>Parcel</th>
                </tr>
              </thead>
              <tbody>
                {naturalForces.map((r) => (
                  <tr key={r.id}>
                    <td title={r.id}>{r.id.slice(0, 8)}&hellip;</td>
                    <td>{r.kind}</td>
                    <td>{r.is_monopolisable ? "Yes" : "No"}</td>
                    <td title={r.parcel_id}>{r.parcel_id.slice(0, 8)}&hellip;</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </section>

      {/* Section 3: Differential Rent */}
      <section className="v3-ch38-section" aria-label="Record a differential rent">
        <h4>3. Differential rent</h4>
        <p className="v3-ch38-subintro">
          Record the conversion of a surplus-profit into differential rent on a
          specific parcel. The annual rent in labour-minutes is the portion of
          surplus-profit appropriated by the landowner.
        </p>
        <div className="v3-ch38-form">
          <label>
            Parcel ID
            <input
              value={drParcelId}
              onChange={(e) => setDrParcelId(e.target.value)}
              placeholder="paste a parcel ID"
            />
          </label>
          <label>
            Surplus-profit (basis-points)
            <input
              type="number"
              min={0}
              value={drSurplusProfitBp}
              onChange={(e) => setDrSurplusProfitBp(Math.max(0, Number(e.target.value)))}
            />
          </label>
          <label>
            Annual rent (labour-minutes)
            <input
              type="number"
              min={0}
              value={drAnnualRentLm}
              onChange={(e) => setDrAnnualRentLm(Math.max(0, Number(e.target.value)))}
            />
            <span className="v3-ch38-hint">
              {(drAnnualRentLm / 60).toFixed(1)}&nbsp;hours; surplus-profit the landowner captures per year
            </span>
          </label>
        </div>
        <button className="v3-ch38-btn" onClick={handleCreateDifferentialRent}>
          Save differential rent
        </button>
        {drError && <p className="v3-ch38-error">{drError}</p>}
        {lastDr && (
          <div className="v3-ch38-preview">
            Saved &mdash; {lastDr.annual_rent_labour_minutes.toLocaleString()}&nbsp;min/yr on parcel{" "}
            {lastDr.parcel_id.slice(0, 8)}&hellip;
          </div>
        )}
      </section>

      {/* Differential rents list */}
      <section className="v3-ch38-section" aria-label="Stored differential rent records">
        <h4>Differential rent records</h4>
        {diffRents.length === 0 ? (
          <p className="v3-ch38-subintro">No differential rents stored yet.</p>
        ) : (
          <div className="v3-ch38-table-wrap">
            <table className="v3-ch38-table">
              <thead>
                <tr>
                  <th>ID</th>
                  <th>Parcel</th>
                  <th>Surplus-profit (bp)</th>
                  <th>Annual rent (LM)</th>
                </tr>
              </thead>
              <tbody>
                {diffRents.map((r) => (
                  <tr key={r.id}>
                    <td title={r.id}>{r.id.slice(0, 8)}&hellip;</td>
                    <td title={r.parcel_id}>{r.parcel_id.slice(0, 8)}&hellip;</td>
                    <td>{r.surplus_profit_bp.toLocaleString()}</td>
                    <td>{r.annual_rent_labour_minutes.toLocaleString()}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </section>

      {/* Section 4: Capitalised Rent Price */}
      <section className="v3-ch38-section" aria-label="Record a capitalised rent price">
        <h4>4. Capitalised rent price (price of land)</h4>
        <p className="v3-ch38-subintro">
          The purchase price of land is capitalised rent: annual rent divided
          by the interest rate. Example: £10/yr at 5&nbsp;% = £200 = 96&nbsp;000 LM.
          Formula: = roundHalfUp(rent&nbsp;&times;&nbsp;10000,&nbsp;rate_bp).
        </p>
        <div className="v3-ch38-form">
          <label>
            Parcel ID
            <input
              value={crpParcelId}
              onChange={(e) => setCrpParcelId(e.target.value)}
              placeholder="paste a parcel ID"
            />
          </label>
          <label>
            Annual rent (labour-minutes)
            <input
              type="number"
              min={0}
              value={crpAnnualRentLm}
              onChange={(e) => setCrpAnnualRentLm(Math.max(0, Number(e.target.value)))}
            />
          </label>
          <label>
            Interest rate (basis-points)
            <input
              type="number"
              min={1}
              value={crpInterestRateBp}
              onChange={(e) => setCrpInterestRateBp(Math.max(1, Number(e.target.value)))}
            />
            <span className="v3-ch38-hint">500 bp = 5&nbsp;%</span>
          </label>
        </div>
        <button className="v3-ch38-btn" onClick={handleCreateCapitalisedRentPrice}>
          Save capitalised rent price
        </button>
        {crpError && <p className="v3-ch38-error">{crpError}</p>}
        {lastCrp && (
          <div className="v3-ch38-preview">
            Saved &mdash; capitalised price:{" "}
            {lastCrp.capitalised_price_labour_minutes.toLocaleString()}&nbsp;LM
          </div>
        )}
      </section>

      {/* Capitalised prices list */}
      <section className="v3-ch38-section" aria-label="Stored capitalised rent prices">
        <h4>Capitalised rent prices</h4>
        {capitalisedPrices.length === 0 ? (
          <p className="v3-ch38-subintro">No capitalised prices stored yet.</p>
        ) : (
          <div className="v3-ch38-table-wrap">
            <table className="v3-ch38-table">
              <thead>
                <tr>
                  <th>ID</th>
                  <th>Parcel</th>
                  <th>Annual rent (LM)</th>
                  <th>Interest rate (bp)</th>
                  <th>Capitalised price (LM)</th>
                </tr>
              </thead>
              <tbody>
                {capitalisedPrices.map((r) => (
                  <tr key={r.id}>
                    <td title={r.id}>{r.id.slice(0, 8)}&hellip;</td>
                    <td title={r.parcel_id}>{r.parcel_id.slice(0, 8)}&hellip;</td>
                    <td>{r.annual_rent_labour_minutes.toLocaleString()}</td>
                    <td>{r.interest_rate_bp.toLocaleString()}</td>
                    <td>{r.capitalised_price_labour_minutes.toLocaleString()}</td>
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
