import { useEffect, useState } from "react";
import { financeApi } from "../../api";
import type {
  DifferentialRentIITableWithInvestments,
  LeaseTerm,
  SuccessiveInvestment,
} from "../../types";
import "./Ch40DifferentialRentII.css";

const DEFAULT_PARCEL_ID = "5eed000000000000003700";

function poundDisplay(bp: number): string {
  return `£${(bp / 100).toFixed(2)}`;
}

export function Ch40DifferentialRentII() {
  const [investments, setInvestments] = useState<SuccessiveInvestment[]>([]);
  const [leaseTerms, setLeaseTerms] = useState<LeaseTerm[]>([]);
  const [listError, setListError] = useState<string | null>(null);

  // Successive investment form
  const [siParcelId, setSiParcelId] = useState(DEFAULT_PARCEL_ID);
  const [siTranche, setSiTranche] = useState(1);
  const [siCapital, setSiCapital] = useState(250);
  const [siOutput, setSiOutput] = useState(400);
  const [siSurplus, setSiSurplus] = useState(900);
  const [siError, setSiError] = useState<string | null>(null);
  const [lastSi, setLastSi] = useState<SuccessiveInvestment | null>(null);

  // DR II aggregate
  const [dr2ParcelId, setDr2ParcelId] = useState(DEFAULT_PARCEL_ID);
  const [dr2Result, setDr2Result] = useState<DifferentialRentIITableWithInvestments | null>(null);
  const [dr2Error, setDr2Error] = useState<string | null>(null);

  // Lease term form
  const [ltParcelId, setLtParcelId] = useState(DEFAULT_PARCEL_ID);
  const [ltStartYear, setLtStartYear] = useState(1867);
  const [ltDuration, setLtDuration] = useState(21);
  const [ltFixedRent, setLtFixedRent] = useState(900);
  const [ltError, setLtError] = useState<string | null>(null);
  const [lastLt, setLastLt] = useState<LeaseTerm | null>(null);

  async function refreshLists() {
    try {
      const [invs, lts] = await Promise.all([
        financeApi.listSuccessiveInvestments(),
        financeApi.listLeaseTerms(),
      ]);
      setInvestments(invs);
      setLeaseTerms(lts);
    } catch (e) {
      setListError(String(e));
    }
  }

  useEffect(() => {
    refreshLists().catch(() => {});
  }, []);

  async function handleCreateInvestment() {
    setSiError(null);
    setLastSi(null);
    try {
      const result = await financeApi.createSuccessiveInvestment({
        parcel_id: siParcelId,
        tranche_number: siTranche,
        capital_bp: siCapital,
        output_quarters: siOutput,
        surplus_profit_bp: siSurplus,
      });
      setLastSi(result);
      await refreshLists();
      try {
        const agg = await financeApi.getDR2Table(siParcelId);
        setDr2Result(agg);
        setDr2ParcelId(siParcelId);
      } catch {
        // aggregate fetch is best-effort after create
      }
    } catch (e) {
      setSiError(String(e));
    }
  }

  async function handleFetchDR2Table() {
    setDr2Error(null);
    setDr2Result(null);
    try {
      const result = await financeApi.getDR2Table(dr2ParcelId);
      setDr2Result(result);
    } catch (e) {
      setDr2Error(String(e));
    }
  }

  async function handleCreateLeaseTerm() {
    setLtError(null);
    setLastLt(null);
    try {
      const result = await financeApi.createLeaseTerm({
        parcel_id: ltParcelId,
        start_year: ltStartYear,
        duration_years: ltDuration,
        fixed_rent_bp: ltFixedRent,
      });
      setLastLt(result);
      await refreshLists();
    } catch (e) {
      setLtError(String(e));
    }
  }

  return (
    <div className="v3-ch40">
      <section className="v3-ch40-intro">
        <p>
          Chapter&nbsp;40 introduces <strong>Differential Rent&nbsp;II</strong> — the
          second form of differential rent, arising from{" "}
          <strong>successive capital investments on the same plot</strong> (intensive
          cultivation), as distinct from DR&nbsp;I's extensive cultivation across soils of
          unequal fertility. The underlying law is identical: the market price is regulated
          by the worst conditions; any better outcome on the same soil yields a
          surplus-profit. DR&nbsp;II therefore{" "}
          <strong>presupposes DR&nbsp;I</strong> — the differential between soil types is
          the baseline from which successive tranches of capital are evaluated.
        </p>
        <p>
          Concentrating more capital on <em>better</em> soil raises rent per acre beyond
          what DR&nbsp;I alone would yield. The decisive feature of DR&nbsp;II is the{" "}
          <strong>lease</strong>: during its term, any newly-formed surplus-profit remains
          with the capitalist farmer; only at renewal does the landlord claim it as
          increased rent. The duration of the lease is therefore the principal battlefield
          between landlord and tenant.
        </p>
      </section>

      {/* Section 1: Successive investments */}
      <section className="v3-ch40-section" aria-label="Record a successive investment">
        <h4>1. Successive capital investments</h4>
        <p className="v3-ch40-subintro">
          Each tranche represents an additional application of capital to the same parcel.
          Capital, surplus-profit and rent are entered in basis-points (100&nbsp;bp = £1). The
          seed fixture uses two £2.50 (250&nbsp;bp) tranches on parcel&nbsp;3700 each yielding £9
          (900&nbsp;bp) surplus-profit (total £18).
        </p>
        <div className="v3-ch40-form">
          <label>
            Parcel ID
            <input
              value={siParcelId}
              onChange={(e) => setSiParcelId(e.target.value)}
              placeholder={DEFAULT_PARCEL_ID}
            />
          </label>
          <label>
            Tranche number
            <input
              type="number"
              min={1}
              value={siTranche}
              onChange={(e) => setSiTranche(Math.max(1, Number(e.target.value)))}
            />
            <span className="v3-ch40-hint">1 = first application, 2 = second, etc.</span>
          </label>
          <label>
            Capital advanced (basis-points)
            <input
              type="number"
              min={1}
              value={siCapital}
              onChange={(e) => setSiCapital(Math.max(1, Number(e.target.value)))}
            />
            <span className="v3-ch40-hint">100&nbsp;bp = £1. Default £2.50 = 250&nbsp;bp.</span>
          </label>
          <label>
            Output (quarters)
            <input
              type="number"
              min={1}
              value={siOutput}
              onChange={(e) => setSiOutput(Math.max(1, Number(e.target.value)))}
            />
          </label>
          <label>
            Surplus-profit (basis-points)
            <input
              type="number"
              min={0}
              value={siSurplus}
              onChange={(e) => setSiSurplus(Math.max(0, Number(e.target.value)))}
            />
            <span className="v3-ch40-hint">100&nbsp;bp = £1. Default £9 = 900&nbsp;bp.</span>
          </label>
        </div>
        <button className="v3-ch40-btn" onClick={handleCreateInvestment}>
          Save investment
        </button>
        {siError && <p className="v3-ch40-error">{siError}</p>}
        {lastSi && (
          <div className="v3-ch40-preview">
            <strong>
              Saved — tranche&nbsp;{lastSi.tranche_number} on parcel&nbsp;
              {lastSi.parcel_id.slice(0, 8)}&hellip;
            </strong>
            <span>
              Capital: {poundDisplay(lastSi.capital_bp)}, output:{" "}
              {lastSi.output_quarters}&nbsp;qtr, surplus-profit:{" "}
              {poundDisplay(lastSi.surplus_profit_bp)}
            </span>
          </div>
        )}
        {investments.length > 0 && (
          <div className="v3-ch40-table-wrap">
            <table className="v3-ch40-table">
              <thead>
                <tr>
                  <th>ID</th>
                  <th>Parcel</th>
                  <th>Tranche</th>
                  <th>Capital (£)</th>
                  <th>Output (qtr)</th>
                  <th>Surplus-profit (£)</th>
                </tr>
              </thead>
              <tbody>
                {investments.map((inv) => (
                  <tr key={inv.id}>
                    <td title={inv.id}>{inv.id.slice(0, 8)}&hellip;</td>
                    <td title={inv.parcel_id}>{inv.parcel_id.slice(0, 8)}&hellip;</td>
                    <td>{inv.tranche_number}</td>
                    <td>{(inv.capital_bp / 100).toFixed(2)}</td>
                    <td>{inv.output_quarters}</td>
                    <td>{(inv.surplus_profit_bp / 100).toFixed(2)}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
        {listError && <p className="v3-ch40-error">{listError}</p>}
      </section>

      {/* Section 2: DR II aggregate */}
      <section className="v3-ch40-section" aria-label="DR II aggregate table">
        <h4>2. Differential Rent&nbsp;II aggregate</h4>
        <p className="v3-ch40-subintro">
          The aggregate sums all tranches on a parcel: total capital, total output, total
          rent (= sum of surplus-profits across all tranches), and rent per acre
          (= total rent &divide; acres; parcels with no recorded acreage default to 1 acre
          for display purposes).
        </p>
        <div className="v3-ch40-form">
          <label>
            Parcel ID
            <input
              value={dr2ParcelId}
              onChange={(e) => setDr2ParcelId(e.target.value)}
              placeholder={DEFAULT_PARCEL_ID}
            />
          </label>
        </div>
        <button className="v3-ch40-btn" onClick={handleFetchDR2Table}>
          Compute for parcel
        </button>
        {dr2Error && <p className="v3-ch40-error">{dr2Error}</p>}
        {dr2Result && (
          <>
            <dl className="v3-ch40-aggregate">
              <dt>Parcel</dt>
              <dd title={dr2Result.table.parcel_id}>
                {dr2Result.table.parcel_id.slice(0, 8)}&hellip;
              </dd>
              <dt>Total capital</dt>
              <dd>{poundDisplay(dr2Result.table.total_capital_bp)}</dd>
              <dt>Total output</dt>
              <dd>{dr2Result.table.total_output_quarters}&nbsp;qtr</dd>
              <dt>Total rent</dt>
              <dd>{poundDisplay(dr2Result.table.total_rent_bp)}</dd>
              <dt>Rent per acre</dt>
              <dd>{poundDisplay(dr2Result.table.rent_per_acre_bp)}</dd>
            </dl>
            {dr2Result.investments.length > 0 && (
              <div className="v3-ch40-table-wrap">
                <table className="v3-ch40-table">
                  <thead>
                    <tr>
                      <th>Tranche</th>
                      <th>Capital (£)</th>
                      <th>Output (qtr)</th>
                      <th>Surplus-profit (£)</th>
                    </tr>
                  </thead>
                  <tbody>
                    {dr2Result.investments.map((inv) => (
                      <tr key={inv.id}>
                        <td>{inv.tranche_number}</td>
                        <td>{(inv.capital_bp / 100).toFixed(2)}</td>
                        <td>{inv.output_quarters}</td>
                        <td>{(inv.surplus_profit_bp / 100).toFixed(2)}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}
          </>
        )}
      </section>

      {/* Section 3: Lease terms */}
      <section className="v3-ch40-section" aria-label="Lease terms">
        <h4>3. Lease terms</h4>
        <p className="v3-ch40-subintro">
          The lease fixes the rent for its duration. Any surplus-profit formed by new
          investments during the lease accrues to the tenant, not the landlord. At renewal
          the landlord raises the fixed rent to capture it. Marx uses 21-year leases as
          the English norm (1867 baseline).
        </p>
        <div className="v3-ch40-form">
          <label>
            Parcel ID
            <input
              value={ltParcelId}
              onChange={(e) => setLtParcelId(e.target.value)}
              placeholder={DEFAULT_PARCEL_ID}
            />
          </label>
          <label>
            Start year
            <input
              type="number"
              min={1600}
              value={ltStartYear}
              onChange={(e) => setLtStartYear(Math.max(1600, Number(e.target.value)))}
            />
          </label>
          <label>
            Duration (years)
            <input
              type="number"
              min={1}
              value={ltDuration}
              onChange={(e) => setLtDuration(Math.max(1, Number(e.target.value)))}
            />
          </label>
          <label>
            Fixed rent (basis-points)
            <input
              type="number"
              min={0}
              value={ltFixedRent}
              onChange={(e) => setLtFixedRent(Math.max(0, Number(e.target.value)))}
            />
            <span className="v3-ch40-hint">100&nbsp;bp = £1. Default £9 = 900&nbsp;bp.</span>
          </label>
        </div>
        <button className="v3-ch40-btn" onClick={handleCreateLeaseTerm}>
          Save lease term
        </button>
        {ltError && <p className="v3-ch40-error">{ltError}</p>}
        {lastLt && (
          <div className="v3-ch40-preview">
            Saved &mdash; parcel {lastLt.parcel_id.slice(0, 8)}&hellip;, start{" "}
            {lastLt.start_year}, {lastLt.duration_years}&nbsp;yr, fixed rent{" "}
            {poundDisplay(lastLt.fixed_rent_bp)}
          </div>
        )}
        {leaseTerms.length > 0 && (
          <div className="v3-ch40-table-wrap">
            <table className="v3-ch40-table">
              <thead>
                <tr>
                  <th>ID</th>
                  <th>Parcel</th>
                  <th>Start year</th>
                  <th>Duration (yr)</th>
                  <th>Fixed rent (£)</th>
                </tr>
              </thead>
              <tbody>
                {leaseTerms.map((lt) => (
                  <tr key={lt.id}>
                    <td title={lt.id}>{lt.id.slice(0, 8)}&hellip;</td>
                    <td title={lt.parcel_id}>{lt.parcel_id.slice(0, 8)}&hellip;</td>
                    <td>{lt.start_year}</td>
                    <td>{lt.duration_years}</td>
                    <td>{(lt.fixed_rent_bp / 100).toFixed(2)}</td>
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
