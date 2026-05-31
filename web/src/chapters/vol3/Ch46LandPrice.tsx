import { useEffect, useState } from "react";
import { financeApi } from "../../api";
import type {
  BuildingSiteRent,
  MiningRent,
  MonopolyPriceRent,
  LandPriceScenario,
  LandPriceScenarioKind,
} from "../../types";
import "./Ch46LandPrice.css";

const KIND_LABELS: Record<LandPriceScenarioKind, string> = {
  1: "1 — Falling interest rate",
  2: "2 — Capital improvements",
  3: "3 — Rent rises with product",
  4: "4 — Rent rises more capital",
  5: "5 — Rent rises price falls",
};

export function Ch46LandPrice() {
  const [buildingRents, setBuildingRents] = useState<BuildingSiteRent[]>([]);
  const [miningRents, setMiningRents] = useState<MiningRent[]>([]);
  const [monopolyRents, setMonopolyRents] = useState<MonopolyPriceRent[]>([]);
  const [landScenarios, setLandScenarios] = useState<LandPriceScenario[]>([]);
  const [listError, setListError] = useState<string | null>(null);

  // Building-site rent form
  const [bParcelId, setBParcelId] = useState("5eed000000000000003700");
  const [bGrade, setBGrade] = useState(5);
  const [bRent, setBRent] = useState(300);
  const [bLease, setBLease] = useState(99);
  const [bError, setBError] = useState<string | null>(null);
  const [lastB, setLastB] = useState<BuildingSiteRent | null>(null);

  // Mining rent form
  const [mParcelId, setMParcelId] = useState("5eed000000000000003700");
  const [mGrade, setMGrade] = useState(1);
  const [mRent, setMRent] = useState(0);
  const [mSurplus, setMSurplus] = useState(0);
  const [mError, setMError] = useState<string | null>(null);
  const [lastM, setLastM] = useState<MiningRent | null>(null);

  // Monopoly-price rent form
  const [pParcelId, setPParcelId] = useState("5eed000000000000003700");
  const [pKind, setPKind] = useState("burgundy_grand_cru");
  const [pValue, setPValue] = useState(200);
  const [pSale, setPSale] = useState(400);
  const [pError, setPError] = useState<string | null>(null);
  const [lastP, setLastP] = useState<MonopolyPriceRent | null>(null);

  // Land-price scenario form
  const [sKind, setSKind] = useState<LandPriceScenarioKind>(1);
  const [sParcelId, setSParcelId] = useState("5eed000000000000003700");
  const [sOldRent, setSoldRent] = useState(10);
  const [sNewRent, setSNewRent] = useState(10);
  const [sOldRate, setSoldRate] = useState(500);
  const [sNewRate, setSNewRate] = useState(250);
  const [sOldPrice, setSoldPrice] = useState(200);
  const [sError, setSError] = useState<string | null>(null);
  const [lastS, setLastS] = useState<LandPriceScenario | null>(null);

  async function refreshLists() {
    try {
      const [br, mr, pr, ls] = await Promise.all([
        financeApi.listBuildingSiteRents(),
        financeApi.listMiningRents(),
        financeApi.listMonopolyPriceRents(),
        financeApi.listLandPriceScenarios(),
      ]);
      setBuildingRents(br);
      setMiningRents(mr);
      setMonopolyRents(pr);
      setLandScenarios(ls);
    } catch (e) {
      setListError(String(e));
    }
  }

  useEffect(() => {
    refreshLists().catch(() => {});
  }, []);

  async function handleCreateBuilding() {
    setBError(null);
    setLastB(null);
    try {
      const result = await financeApi.createBuildingSiteRent({
        parcel_id: bParcelId,
        location_grade: bGrade,
        annual_rent_labour_minutes: bRent,
        lease_term_years: bLease,
      });
      setLastB(result);
      await refreshLists();
    } catch (e) {
      setBError(String(e));
    }
  }

  async function handleCreateMining() {
    setMError(null);
    setLastM(null);
    try {
      const result = await financeApi.createMiningRent({
        parcel_id: mParcelId,
        ore_grade: mGrade,
        annual_rent_labour_minutes: mRent,
        surplus_profit_bp: mSurplus,
      });
      setLastM(result);
      await refreshLists();
    } catch (e) {
      setMError(String(e));
    }
  }

  async function handleCreateMonopoly() {
    setPError(null);
    setLastP(null);
    try {
      const result = await financeApi.createMonopolyPriceRent({
        parcel_id: pParcelId,
        product_kind: pKind,
        value_labour_minutes: pValue,
        monopoly_sale_price_labour_minutes: pSale,
      });
      setLastP(result);
      await refreshLists();
    } catch (e) {
      setPError(String(e));
    }
  }

  async function handleCreateScenario() {
    setSError(null);
    setLastS(null);
    try {
      const result = await financeApi.createLandPriceScenario({
        kind: sKind,
        parcel_id: sParcelId,
        old_annual_rent_bp: sOldRent,
        new_annual_rent_bp: sNewRent,
        old_interest_rate_bp: sOldRate,
        new_interest_rate_bp: sNewRate,
        old_land_price_bp: sOldPrice,
      });
      setLastS(result);
      await refreshLists();
    } catch (e) {
      setSError(String(e));
    }
  }

  return (
    <div className="v3-ch46">
      <section className="v3-ch46-intro">
        <p>
          Chapter&nbsp;46 extends the rent framework beyond agriculture. <strong>Building-site
          rent</strong> is determined by location rather than fertility — the landlord is entirely
          passive; the tenant-capitalist does all the improvement. The land under a building is
          valued by its anticipated rent capitalised at the prevailing interest rate, not by what
          it might grow.
        </p>
        <p>
          <strong>Rent in mining</strong> follows the same logic as agricultural differential rent:
          productivity differences between mines set the rent. The <em>worst</em> working mine
          (ore grade&nbsp;1) yields <em>zero</em> rent — it sets the regulating price; all better
          mines yield a positive differential rent.
        </p>
        <p>
          <strong>Monopoly-price rent</strong> is the purest form: the sale price is fixed by
          buyer demand and willingness to pay, entirely above value. The whole excess of sale price
          over value becomes rent — no surplus-profit mediates it.
        </p>
        <p>
          Finally, the <strong>price of land</strong> is capitalised rent:{" "}
          <code>P&nbsp;= annual&nbsp;rent &divide; interest&nbsp;rate</code>. A fall in the
          interest rate alone raises land prices without any change in the rent itself — a pure
          monetary effect.
        </p>
        <p className="v3-ch46-note">
          Rates in basis points (bp). 10,000&nbsp;bp&nbsp;= 100%. Labour values in minutes.
          Land-price bp: £1&nbsp;= 10,000&nbsp;bp (so £10 @ 5%&nbsp;= 200&nbsp;bp, £10 @ 2.5%&nbsp;= 400&nbsp;bp).
        </p>
      </section>

      {/* Section 1 — Building-site rents */}
      <section className="v3-ch46-section" aria-label="Building-site rents">
        <h4>1. Building-site rents</h4>
        <p className="v3-ch46-subintro">
          Record a building-site lease. Location grade (higher = better location) determines
          the rent; the landlord's ownership is purely passive. Seed fixture: London West End
          parcel, grade&nbsp;5, £300/yr rent (= 300&nbsp;labour-minutes), 99-year lease.
        </p>
        <div className="v3-ch46-form">
          <label>
            Parcel ID
            <input value={bParcelId} onChange={(e) => setBParcelId(e.target.value)} />
          </label>
          <label>
            Location grade
            <input
              type="number"
              min={1}
              value={bGrade}
              onChange={(e) => setBGrade(Math.max(1, Number(e.target.value)))}
            />
          </label>
          <label>
            Annual rent (labour-minutes)
            <input
              type="number"
              min={0}
              value={bRent}
              onChange={(e) => setBRent(Math.max(0, Number(e.target.value)))}
            />
          </label>
          <label>
            Lease term (years)
            <input
              type="number"
              min={1}
              value={bLease}
              onChange={(e) => setBLease(Math.max(1, Number(e.target.value)))}
            />
          </label>
        </div>
        <button className="v3-ch46-btn" onClick={handleCreateBuilding}>
          Record building-site rent
        </button>
        {bError && <p className="v3-ch46-error">{bError}</p>}
        {lastB && (
          <div className="v3-ch46-preview">
            <strong>
              Saved — parcel&nbsp;{lastB.parcel_id.slice(0, 8)}&hellip;,
              grade&nbsp;{lastB.location_grade},
              rent&nbsp;{lastB.annual_rent_labour_minutes}&nbsp;min/yr,
              lease&nbsp;{lastB.lease_term_years}&nbsp;yr
            </strong>
            <span className="v3-ch46-hint">ID: {lastB.id}</span>
          </div>
        )}
        {buildingRents.length > 0 && (
          <div className="v3-ch46-table-wrap">
            <table className="v3-ch46-table">
              <thead>
                <tr>
                  <th>ID</th>
                  <th>Parcel</th>
                  <th>Grade</th>
                  <th>Annual rent (min)</th>
                  <th>Lease (yr)</th>
                </tr>
              </thead>
              <tbody>
                {buildingRents.map((r) => (
                  <tr key={r.id}>
                    <td title={r.id}>{r.id.slice(0, 8)}&hellip;</td>
                    <td title={r.parcel_id}>{r.parcel_id.slice(0, 8)}&hellip;</td>
                    <td>{r.location_grade}</td>
                    <td>{r.annual_rent_labour_minutes}</td>
                    <td>{r.lease_term_years}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
        {listError && <p className="v3-ch46-error">{listError}</p>}
      </section>

      {/* Section 2 — Mining rents */}
      <section className="v3-ch46-section" aria-label="Mining rents">
        <h4>2. Mining rents</h4>
        <p className="v3-ch46-subintro">
          Mining rent follows differential rent: the worst working mine (ore grade&nbsp;1) sets
          the regulating price and yields <strong>zero rent</strong>. Higher-grade mines yield
          positive surplus profit, which the landowner captures as rent. Seed fixture: Scottish
          coal worst mine, grade&nbsp;1, rent&nbsp;=&nbsp;0, surplus&nbsp;profit&nbsp;=&nbsp;0.
        </p>
        <div className="v3-ch46-form">
          <label>
            Parcel ID
            <input value={mParcelId} onChange={(e) => setMParcelId(e.target.value)} />
          </label>
          <label>
            Ore grade
            <input
              type="number"
              min={1}
              value={mGrade}
              onChange={(e) => setMGrade(Math.max(1, Number(e.target.value)))}
            />
            <span className="v3-ch46-hint">grade&nbsp;1 = worst (zero rent)</span>
          </label>
          <label>
            Annual rent (labour-minutes)
            <input
              type="number"
              min={0}
              value={mRent}
              onChange={(e) => setMRent(Math.max(0, Number(e.target.value)))}
            />
          </label>
          <label>
            Surplus profit (bp)
            <input
              type="number"
              min={0}
              value={mSurplus}
              onChange={(e) => setMSurplus(Math.max(0, Number(e.target.value)))}
            />
          </label>
        </div>
        <button className="v3-ch46-btn" onClick={handleCreateMining}>
          Record mining rent
        </button>
        {mError && <p className="v3-ch46-error">{mError}</p>}
        {lastM && (
          <div className="v3-ch46-preview">
            <strong>
              Saved — parcel&nbsp;{lastM.parcel_id.slice(0, 8)}&hellip;,
              grade&nbsp;{lastM.ore_grade},
              rent&nbsp;{lastM.annual_rent_labour_minutes}&nbsp;min/yr,
              surplus&nbsp;{lastM.surplus_profit_bp}&nbsp;bp
            </strong>
            <span className="v3-ch46-hint">ID: {lastM.id}</span>
          </div>
        )}
        {miningRents.length > 0 && (
          <div className="v3-ch46-table-wrap">
            <table className="v3-ch46-table">
              <thead>
                <tr>
                  <th>ID</th>
                  <th>Parcel</th>
                  <th>Ore grade</th>
                  <th>Annual rent (min)</th>
                  <th>Surplus profit (bp)</th>
                </tr>
              </thead>
              <tbody>
                {miningRents.map((r) => (
                  <tr key={r.id}>
                    <td title={r.id}>{r.id.slice(0, 8)}&hellip;</td>
                    <td title={r.parcel_id}>{r.parcel_id.slice(0, 8)}&hellip;</td>
                    <td>{r.ore_grade}</td>
                    <td>{r.annual_rent_labour_minutes}</td>
                    <td>{r.surplus_profit_bp}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </section>

      {/* Section 3 — Monopoly-price rents */}
      <section className="v3-ch46-section" aria-label="Monopoly-price rents">
        <h4>3. Monopoly-price rents</h4>
        <p className="v3-ch46-subintro">
          In monopoly-price rent the sale price is determined purely by buyer demand — it exceeds
          value without any mediation by the rate of profit. The server computes{" "}
          <code>monopoly_rent = max(0, sale&nbsp;&minus;&nbsp;value)</code>. Seed fixture:
          Burgundy grand cru — value&nbsp;200&nbsp;min, sale&nbsp;400&nbsp;min, rent&nbsp;200&nbsp;min.
        </p>
        <div className="v3-ch46-form">
          <label>
            Parcel ID
            <input value={pParcelId} onChange={(e) => setPParcelId(e.target.value)} />
          </label>
          <label>
            Product kind
            <input
              value={pKind}
              onChange={(e) => setPKind(e.target.value)}
              placeholder="burgundy_grand_cru"
            />
          </label>
          <label>
            Value (labour-minutes)
            <input
              type="number"
              min={0}
              value={pValue}
              onChange={(e) => setPValue(Math.max(0, Number(e.target.value)))}
            />
          </label>
          <label>
            Monopoly sale price (labour-minutes)
            <input
              type="number"
              min={0}
              value={pSale}
              onChange={(e) => setPSale(Math.max(0, Number(e.target.value)))}
            />
            <span className="v3-ch46-hint">
              monopoly rent = max(0,&nbsp;{pSale}&nbsp;&minus;&nbsp;{pValue})&nbsp;=&nbsp;
              {Math.max(0, pSale - pValue)}&nbsp;min
            </span>
          </label>
        </div>
        <button className="v3-ch46-btn" onClick={handleCreateMonopoly}>
          Record monopoly-price rent
        </button>
        {pError && <p className="v3-ch46-error">{pError}</p>}
        {lastP && (
          <div className="v3-ch46-preview">
            <strong>
              Saved — {lastP.product_kind}: value&nbsp;{lastP.value_labour_minutes}&nbsp;min,
              sale&nbsp;{lastP.monopoly_sale_price_labour_minutes}&nbsp;min,
              rent&nbsp;{lastP.monopoly_rent_labour_minutes}&nbsp;min
            </strong>
            <span className="v3-ch46-hint">ID: {lastP.id}</span>
          </div>
        )}
        {monopolyRents.length > 0 && (
          <div className="v3-ch46-table-wrap">
            <table className="v3-ch46-table">
              <thead>
                <tr>
                  <th>ID</th>
                  <th>Parcel</th>
                  <th>Product</th>
                  <th>Value (min)</th>
                  <th>Sale (min)</th>
                  <th>Rent (min)</th>
                </tr>
              </thead>
              <tbody>
                {monopolyRents.map((r) => (
                  <tr key={r.id}>
                    <td title={r.id}>{r.id.slice(0, 8)}&hellip;</td>
                    <td title={r.parcel_id}>{r.parcel_id.slice(0, 8)}&hellip;</td>
                    <td>{r.product_kind}</td>
                    <td>{r.value_labour_minutes}</td>
                    <td>{r.monopoly_sale_price_labour_minutes}</td>
                    <td>{r.monopoly_rent_labour_minutes}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </section>

      {/* Section 4 — Land-price scenarios */}
      <section className="v3-ch46-section" aria-label="Land-price scenarios">
        <h4>4. Land-price scenarios</h4>
        <p className="v3-ch46-subintro">
          The price of land is capitalised rent:{" "}
          <code>P&nbsp;= annual&nbsp;rent &divide; interest&nbsp;rate</code>. The server computes{" "}
          <code>new_land_price_bp&nbsp;=&nbsp;roundHalfUp(new_annual_rent_bp&nbsp;&times;&nbsp;10000,&nbsp;new_interest_rate_bp)</code>.
          Seed fixture: £10/yr rent at 5% (500&nbsp;bp) → price £200 (200&nbsp;bp); rate halves to
          2.5% (250&nbsp;bp) → price £400 (400&nbsp;bp) with rent unchanged.
        </p>
        <div className="v3-ch46-form">
          <label>
            Scenario kind
            <select
              value={sKind}
              onChange={(e) => setSKind(Number(e.target.value) as LandPriceScenarioKind)}
            >
              {(Object.keys(KIND_LABELS) as unknown as LandPriceScenarioKind[]).map((k) => (
                <option key={k} value={k}>{KIND_LABELS[k]}</option>
              ))}
            </select>
          </label>
          <label>
            Parcel ID
            <input value={sParcelId} onChange={(e) => setSParcelId(e.target.value)} />
          </label>
          <label>
            Old annual rent (bp)
            <input
              type="number"
              min={0}
              value={sOldRent}
              onChange={(e) => setSoldRent(Math.max(0, Number(e.target.value)))}
            />
          </label>
          <label>
            New annual rent (bp)
            <input
              type="number"
              min={0}
              value={sNewRent}
              onChange={(e) => setSNewRent(Math.max(0, Number(e.target.value)))}
            />
          </label>
          <label>
            Old interest rate (bp)
            <input
              type="number"
              min={1}
              value={sOldRate}
              onChange={(e) => setSoldRate(Math.max(1, Number(e.target.value)))}
            />
          </label>
          <label>
            New interest rate (bp)
            <input
              type="number"
              min={1}
              value={sNewRate}
              onChange={(e) => setSNewRate(Math.max(1, Number(e.target.value)))}
            />
          </label>
          <label>
            Old land price (bp)
            <input
              type="number"
              min={0}
              value={sOldPrice}
              onChange={(e) => setSoldPrice(Math.max(0, Number(e.target.value)))}
            />
          </label>
        </div>
        <button className="v3-ch46-btn" onClick={handleCreateScenario}>
          Record land-price scenario
        </button>
        {sError && <p className="v3-ch46-error">{sError}</p>}
        {lastS && (
          <div className="v3-ch46-preview">
            <strong>
              Saved — kind&nbsp;{lastS.kind}, parcel&nbsp;{lastS.parcel_id.slice(0, 8)}&hellip;,
              old&nbsp;price&nbsp;{lastS.old_land_price_bp}&nbsp;bp,
              new&nbsp;price&nbsp;{lastS.new_land_price_bp}&nbsp;bp
            </strong>
            <span className="v3-ch46-hint">ID: {lastS.id}</span>
          </div>
        )}
        {landScenarios.length > 0 && (
          <div className="v3-ch46-table-wrap">
            <table className="v3-ch46-table">
              <thead>
                <tr>
                  <th>ID</th>
                  <th>Kind</th>
                  <th>Parcel</th>
                  <th>Old rent (bp)</th>
                  <th>New rent (bp)</th>
                  <th>Old rate (bp)</th>
                  <th>New rate (bp)</th>
                  <th>Old price (bp)</th>
                  <th>New price (bp)</th>
                </tr>
              </thead>
              <tbody>
                {landScenarios.map((s) => (
                  <tr key={s.id}>
                    <td title={s.id}>{s.id.slice(0, 8)}&hellip;</td>
                    <td>{s.kind}</td>
                    <td title={s.parcel_id}>{s.parcel_id.slice(0, 8)}&hellip;</td>
                    <td>{s.old_annual_rent_bp}</td>
                    <td>{s.new_annual_rent_bp}</td>
                    <td>{s.old_interest_rate_bp}</td>
                    <td>{s.new_interest_rate_bp}</td>
                    <td>{s.old_land_price_bp}</td>
                    <td>{s.new_land_price_bp}</td>
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
