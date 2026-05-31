import { useEffect, useState } from "react";
import { financeApi } from "../../api";
import type {
  DifferentialRentITable,
  DifferentialRentITableWithEntries,
  LocationRentFactor,
  SoilGrade,
} from "../../types";
import "./Ch39DifferentialRentI.css";

const GRADE_LABELS: Record<SoilGrade, string> = { 1: "A", 2: "B", 3: "C", 4: "D" };

// Marx Table I defaults: capital 250 each, outputs 100/200/300/400 qtr, reg. price 300 bp/qtr
const DEFAULT_OUTPUTS: Record<SoilGrade, number> = { 1: 100, 2: 200, 3: 300, 4: 400 };
const CAPITAL_BP = 25000; // £250 = 25000 bp
const PRICE_BP = 300;     // £3/qtr = 300 bp

function poundDisplay(bp: number): string {
  return `£${(bp / 100).toFixed(2)}`;
}

export function Ch39DifferentialRentI() {
  const [tables, setTables] = useState<DifferentialRentITable[]>([]);
  const [locationFactors, setLocationFactors] = useState<LocationRentFactor[]>([]);
  const [listError, setListError] = useState<string | null>(null);

  // DR I table form
  const [description, setDescription] = useState("Marx Table I — Ch. 39");
  const [outputs, setOutputs] = useState<Record<SoilGrade, number>>({ ...DEFAULT_OUTPUTS });
  const [tableError, setTableError] = useState<string | null>(null);
  const [lastTable, setLastTable] = useState<DifferentialRentITableWithEntries | null>(null);

  // Location rent factor form
  const [lrfParcelId, setLrfParcelId] = useState("");
  const [lrfTransportBp, setLrfTransportBp] = useState(500);
  const [lrfRentEquivBp, setLrfRentEquivBp] = useState(500);
  const [lrfError, setLrfError] = useState<string | null>(null);
  const [lastLrf, setLastLrf] = useState<LocationRentFactor | null>(null);

  async function refreshLists() {
    try {
      const [tbls, lrfs] = await Promise.all([
        financeApi.listDR1Tables(),
        financeApi.listLocationRentFactors(),
      ]);
      setTables(tbls);
      setLocationFactors(lrfs);
    } catch (e) {
      setListError(String(e));
    }
  }

  useEffect(() => {
    refreshLists().catch(() => {});
  }, []);

  async function handleCreateTable() {
    setTableError(null);
    setLastTable(null);
    try {
      const grades: SoilGrade[] = [1, 2, 3, 4];
      const result = await financeApi.createDR1Table({
        description,
        entries: grades.map((g) => ({
          grade: g,
          capital_advanced_bp: CAPITAL_BP,
          output_quarters: outputs[g],
          price_of_production_bp: PRICE_BP,
        })),
      });
      setLastTable(result);
      await refreshLists();
    } catch (e) {
      setTableError(String(e));
    }
  }

  async function handleCreateLocationRentFactor() {
    setLrfError(null);
    setLastLrf(null);
    try {
      const result = await financeApi.createLocationRentFactor({
        parcel_id: lrfParcelId,
        transport_cost_bp: lrfTransportBp,
        rent_equivalent_bp: lrfRentEquivBp,
      });
      setLastLrf(result);
      await refreshLists();
    } catch (e) {
      setLrfError(String(e));
    }
  }

  const grades: SoilGrade[] = [1, 2, 3, 4];

  return (
    <div className="v3-ch39">
      <section className="v3-ch39-intro">
        <p>
          Chapter&nbsp;39 introduces <strong>Differential Rent&nbsp;I</strong> — the
          first form of differential rent, arising from the simultaneous cultivation of
          soils of <strong>unequal fertility</strong> with equal applications of capital.
          The market price of agricultural produce is regulated by the{" "}
          <strong>worst (least fertile) soil</strong> that must be cultivated to meet
          social demand. That worst soil yields no rent. Every better soil, producing
          the same output with the same capital outlay, yields a{" "}
          <strong>surplus-product</strong> — output above the regulating price — which
          landed property converts into <strong>differential rent</strong>.
        </p>
        <p>
          Marx's Table&nbsp;I shows four soil grades A–D each receiving £250 capital at
          £3/qtr regulating price. Grade&nbsp;A (worst) yields 100&nbsp;qtr = no rent;
          B yields 200&nbsp;qtr = £150 rent; C yields 300&nbsp;qtr = £300 rent;
          D yields 400&nbsp;qtr = £450 rent — total rental £900. Differential rent I
          also arises from <strong>location</strong>: soils nearer to market have lower
          transport costs, which reduces their individual price of production below the
          regulating price, yielding a rent equivalent to the transport saving.
        </p>
      </section>

      {/* Section 1: Build a DR I table */}
      <section className="v3-ch39-section" aria-label="Build a Differential Rent I table">
        <h4>1. Build a Differential Rent&nbsp;I table</h4>
        <p className="v3-ch39-subintro">
          Pre-filled with Marx Table&nbsp;I (capital £250 each, regulating price
          £3/qtr). Edit the output quantities for grades B–D then submit.
        </p>
        <div className="v3-ch39-form">
          <label>
            Description
            <input
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              placeholder="e.g. Marx Table I — Ch. 39"
            />
          </label>
          {grades.map((g) => (
            <label key={g}>
              Soil&nbsp;{GRADE_LABELS[g]} output (quarters)
              <input
                type="number"
                min={1}
                value={outputs[g]}
                onChange={(e) =>
                  setOutputs((prev) => ({ ...prev, [g]: Math.max(1, Number(e.target.value)) }))
                }
              />
              <span className="v3-ch39-hint">
                Capital £250, reg. price £3/qtr
                {g === 1 ? " — regulating (worst) soil, no rent" : ""}
              </span>
            </label>
          ))}
        </div>
        <button className="v3-ch39-btn" onClick={handleCreateTable}>
          Save DR&nbsp;I table
        </button>
        {tableError && <p className="v3-ch39-error">{tableError}</p>}
        {lastTable && (
          <div className="v3-ch39-preview">
            <strong>Saved — {lastTable.table.description}</strong>
            <span>
              {" "}Total rental: {poundDisplay(lastTable.table.total_rental_bp)},
              regulating grade: {GRADE_LABELS[lastTable.table.regulating_grade]}
            </span>
            <div className="v3-ch39-table-wrap">
              <table className="v3-ch39-table">
                <thead>
                  <tr>
                    <th>Soil</th>
                    <th>Capital (£)</th>
                    <th>Output (qtr)</th>
                    <th>Reg. Price (£/qtr)</th>
                    <th>Surplus-Product (qtr)</th>
                    <th>Money-Rent (£)</th>
                  </tr>
                </thead>
                <tbody>
                  {lastTable.entries.map((e) => (
                    <tr key={e.id}>
                      <td>{GRADE_LABELS[e.grade]}</td>
                      <td>{(e.capital_advanced_bp / 100).toFixed(0)}</td>
                      <td>{e.output_quarters}</td>
                      <td>{(e.price_of_production_bp / 100).toFixed(2)}</td>
                      <td>{e.surplus_product_qtrs}</td>
                      <td>{(e.money_rent_bp / 100).toFixed(2)}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
        )}
      </section>

      {/* Section 2: Saved tables */}
      <section className="v3-ch39-section" aria-label="Saved Differential Rent I tables">
        <h4>Saved DR&nbsp;I tables</h4>
        {listError && <p className="v3-ch39-error">{listError}</p>}
        {tables.length === 0 ? (
          <p className="v3-ch39-subintro">No tables saved yet.</p>
        ) : (
          <div className="v3-ch39-table-wrap">
            <table className="v3-ch39-table">
              <thead>
                <tr>
                  <th>ID</th>
                  <th>Description</th>
                  <th>Total Rental</th>
                  <th>Regulating Grade</th>
                </tr>
              </thead>
              <tbody>
                {tables.map((t) => (
                  <tr key={t.id}>
                    <td title={t.id}>{t.id.slice(0, 8)}&hellip;</td>
                    <td>{t.description}</td>
                    <td>{poundDisplay(t.total_rental_bp)}</td>
                    <td>{GRADE_LABELS[t.regulating_grade]}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </section>

      {/* Section 3: Location rent factors */}
      <section className="v3-ch39-section" aria-label="Location rent factors">
        <h4>3. Location rent factors</h4>
        <p className="v3-ch39-subintro">
          A parcel closer to market saves transport costs. That saving = the
          individual price of production below the regulating price = a rent
          equivalent. Record the transport-cost saving and the resulting rent
          equivalent for a given parcel.
        </p>
        <div className="v3-ch39-form">
          <label>
            Parcel ID
            <input
              value={lrfParcelId}
              onChange={(e) => setLrfParcelId(e.target.value)}
              placeholder="paste a parcel ID from Ch. 37"
            />
          </label>
          <label>
            Transport cost saving (basis-points)
            <input
              type="number"
              min={0}
              value={lrfTransportBp}
              onChange={(e) => setLrfTransportBp(Math.max(0, Number(e.target.value)))}
            />
            <span className="v3-ch39-hint">100 bp = £1</span>
          </label>
          <label>
            Rent equivalent (basis-points)
            <input
              type="number"
              min={0}
              value={lrfRentEquivBp}
              onChange={(e) => setLrfRentEquivBp(Math.max(0, Number(e.target.value)))}
            />
          </label>
        </div>
        <button className="v3-ch39-btn" onClick={handleCreateLocationRentFactor}>
          Save location rent factor
        </button>
        {lrfError && <p className="v3-ch39-error">{lrfError}</p>}
        {lastLrf && (
          <div className="v3-ch39-preview">
            Saved &mdash; parcel {lastLrf.parcel_id.slice(0, 8)}&hellip;,
            transport saving {poundDisplay(lastLrf.transport_cost_bp)},
            rent equiv. {poundDisplay(lastLrf.rent_equivalent_bp)}
          </div>
        )}
        {locationFactors.length > 0 && (
          <div className="v3-ch39-table-wrap" style={{ marginTop: "1rem" }}>
            <table className="v3-ch39-table">
              <thead>
                <tr>
                  <th>ID</th>
                  <th>Parcel</th>
                  <th>Transport Saving</th>
                  <th>Rent Equivalent</th>
                </tr>
              </thead>
              <tbody>
                {locationFactors.map((r) => (
                  <tr key={r.id}>
                    <td title={r.id}>{r.id.slice(0, 8)}&hellip;</td>
                    <td title={r.parcel_id}>{r.parcel_id.slice(0, 8)}&hellip;</td>
                    <td>{poundDisplay(r.transport_cost_bp)}</td>
                    <td>{poundDisplay(r.rent_equivalent_bp)}</td>
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
