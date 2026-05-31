import { useEffect, useState } from "react";
import { financeApi } from "../../api";
import type { GroundRent, LandParcel } from "../../types";
import "./Ch37Introduction.css";

function fertilityLabel(grade: number): string {
  switch (grade) {
    case 1: return "A (worst — regulating)";
    case 2: return "B";
    case 3: return "C";
    case 4: return "D (best)";
    default: return `Grade ${grade}`;
  }
}

function rentFormLabel(form: number): string {
  switch (form) {
    case 1: return "Differential";
    case 2: return "Absolute";
    case 3: return "Monopoly";
    default: return `Form ${form}`;
  }
}

export function Ch37Introduction() {
  const [parcels, setParcels] = useState<LandParcel[]>([]);
  const [rents, setRents] = useState<GroundRent[]>([]);
  const [listError, setListError] = useState<string | null>(null);

  // Parcel form state
  const [ownerId, setOwnerId] = useState("Duke of Sutherland");
  const [fertilityGrade, setFertilityGrade] = useState(4);
  const [areaHectares, setAreaHectares] = useState(800000);
  const [location, setLocation] = useState("Sutherland, Scotland");
  const [parcelCreateError, setParcelCreateError] = useState<string | null>(null);
  const [lastParcel, setLastParcel] = useState<LandParcel | null>(null);

  // Rent form state
  const [parcelId, setParcelId] = useState("");
  const [capitalistId, setCapitalistId] = useState("");
  const [landOwnerId, setLandOwnerId] = useState("");
  const [rentForm, setRentForm] = useState(1);
  const [amountLabourMinutes, setAmountLabourMinutes] = useState(600);
  const [periodYears, setPeriodYears] = useState(1);
  const [rentCreateError, setRentCreateError] = useState<string | null>(null);
  const [lastRent, setLastRent] = useState<GroundRent | null>(null);

  async function refreshLists() {
    try {
      const [ps, rs] = await Promise.all([
        financeApi.listLandParcels(),
        financeApi.listGroundRents(),
      ]);
      setParcels(ps);
      setRents(rs);
    } catch (e) {
      setListError(String(e));
    }
  }

  useEffect(() => {
    refreshLists().catch(() => {});
  }, []);

  async function handleCreateParcel() {
    setParcelCreateError(null);
    setLastParcel(null);
    try {
      const result = await financeApi.createLandParcel({
        owner_id: ownerId,
        fertility_grade: fertilityGrade,
        area_hectares: areaHectares,
        location,
      });
      setLastParcel(result);
      await refreshLists();
    } catch (e) {
      setParcelCreateError(String(e));
    }
  }

  async function handleCreateRent() {
    setRentCreateError(null);
    setLastRent(null);
    try {
      const result = await financeApi.createGroundRent({
        parcel_id: parcelId,
        capitalist_id: capitalistId,
        land_owner_id: landOwnerId,
        form: rentForm,
        amount_labour_minutes: amountLabourMinutes,
        period_years: periodYears,
      });
      setLastRent(result);
      await refreshLists();
    } catch (e) {
      setRentCreateError(String(e));
    }
  }

  return (
    <div className="v3-ch37">
      <section className="v3-ch37-intro">
        <p>
          Part&nbsp;VI opens with the <strong>transformation of surplus-profit
          into ground-rent</strong>. Ground-rent is the portion of surplus-value
          that falls to the <strong>landowner</strong> by virtue of the monopoly
          of land. It presupposes three classes whose revenues correspond to the
          three great factors of production: the <strong>landowner</strong>{" "}
          (rent), the <strong>agricultural capitalist</strong> (profit), and the{" "}
          <strong>wage-labourer</strong> (wages). The capitalist farms the land
          for a term; at the end of the lease the landlord appropriates all
          improvements; the worker produces the surplus out of which both profit
          and rent are drawn.
        </p>
        <p>
          Marx distinguishes <strong>three forms</strong> of ground-rent:
          (1)&nbsp;<strong>differential rent</strong> — arises from differences
          in fertility or location among parcels of equal capital investment;
          (2)&nbsp;<strong>absolute rent</strong> — arises from the monopoly of
          landownership itself, even on the worst soil; (3)&nbsp;
          <strong>monopoly rent</strong> — arises where a commodity commands a
          price exceeding its price of production and value by pure scarcity.
          All three forms are ultimately deductions from surplus-value; none
          create new value.
        </p>
      </section>

      {/* Record a land parcel */}
      <section className="v3-ch37-section" aria-label="Record a land parcel">
        <h4>Record a land parcel</h4>
        <p className="v3-ch37-subintro">
          Enter a parcel of agricultural land. Fertility grade 1 = worst
          (regulating) soil; grade 4 = best.
        </p>
        <div className="v3-ch37-form">
          <label>
            Owner ID
            <input
              value={ownerId}
              onChange={(e) => setOwnerId(e.target.value)}
            />
            <span className="v3-ch37-hint">e.g. &ldquo;Duke of Sutherland&rdquo;</span>
          </label>
          <label>
            Fertility grade
            <select
              value={fertilityGrade}
              onChange={(e) => setFertilityGrade(Number(e.target.value))}
            >
              <option value={1}>1 — A (worst, regulating)</option>
              <option value={2}>2 — B</option>
              <option value={3}>3 — C</option>
              <option value={4}>4 — D (best)</option>
            </select>
          </label>
          <label>
            Area (hectares)
            <input
              type="number"
              min={0}
              value={areaHectares}
              onChange={(e) => setAreaHectares(Math.max(0, Number(e.target.value)))}
            />
          </label>
          <label>
            Location
            <input
              value={location}
              onChange={(e) => setLocation(e.target.value)}
            />
          </label>
        </div>
        <button className="v3-ch37-btn" onClick={handleCreateParcel}>
          Save land parcel
        </button>
        {parcelCreateError && <p className="v3-ch37-error">{parcelCreateError}</p>}
        {lastParcel && (
          <div className="v3-ch37-preview">
            Saved &mdash; {lastParcel.location} ({fertilityLabel(lastParcel.fertility_grade)},&nbsp;
            {lastParcel.area_hectares.toLocaleString()}&nbsp;ha)
          </div>
        )}
      </section>

      {/* Parcel list */}
      <section className="v3-ch37-section" aria-label="Stored land parcels">
        <h4>Land parcels</h4>
        {listError && <p className="v3-ch37-error">{listError}</p>}
        {parcels.length === 0 ? (
          <p className="v3-ch37-subintro">No parcels stored yet.</p>
        ) : (
          <div className="v3-ch37-table-wrap">
            <table className="v3-ch37-table">
              <thead>
                <tr>
                  <th>ID</th>
                  <th>Owner</th>
                  <th>Fertility</th>
                  <th>Area (ha)</th>
                  <th>Location</th>
                </tr>
              </thead>
              <tbody>
                {parcels.map((p) => (
                  <tr key={p.id}>
                    <td title={p.id}>{p.id.slice(0, 8)}&hellip;</td>
                    <td>{p.owner_id}</td>
                    <td>{fertilityLabel(p.fertility_grade)}</td>
                    <td>{p.area_hectares.toLocaleString()}</td>
                    <td>{p.location}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </section>

      {/* Record a ground-rent payment */}
      <section className="v3-ch37-section" aria-label="Record a ground-rent payment">
        <h4>Record a ground-rent payment</h4>
        <p className="v3-ch37-subintro">
          Record a rent payment from an agricultural capitalist to a landowner
          for use of a registered parcel.
        </p>
        <div className="v3-ch37-form">
          <label>
            Parcel ID
            <input
              value={parcelId}
              onChange={(e) => setParcelId(e.target.value)}
              placeholder="paste a parcel ID from the table above"
            />
          </label>
          <label>
            Capitalist ID
            <input
              value={capitalistId}
              onChange={(e) => setCapitalistId(e.target.value)}
              placeholder="e.g. Norfolk farmer"
            />
          </label>
          <label>
            Land owner ID
            <input
              value={landOwnerId}
              onChange={(e) => setLandOwnerId(e.target.value)}
              placeholder="e.g. Duke of Sutherland"
            />
          </label>
          <label>
            Rent form
            <select
              value={rentForm}
              onChange={(e) => setRentForm(Number(e.target.value))}
            >
              <option value={1}>1 — Differential</option>
              <option value={2}>2 — Absolute</option>
              <option value={3}>3 — Monopoly</option>
            </select>
          </label>
          <label>
            Amount (labour-minutes)
            <input
              type="number"
              min={0}
              value={amountLabourMinutes}
              onChange={(e) => setAmountLabourMinutes(Math.max(0, Number(e.target.value)))}
            />
            <span className="v3-ch37-hint">
              {(amountLabourMinutes / 60).toFixed(1)}&nbsp;hours of surplus-value redirected to the landowner
            </span>
          </label>
          <label>
            Period (years)
            <input
              type="number"
              min={1}
              value={periodYears}
              onChange={(e) => setPeriodYears(Math.max(1, Number(e.target.value)))}
            />
          </label>
        </div>
        <button className="v3-ch37-btn" onClick={handleCreateRent}>
          Save ground-rent
        </button>
        {rentCreateError && <p className="v3-ch37-error">{rentCreateError}</p>}
        {lastRent && (
          <div className="v3-ch37-preview">
            Saved &mdash; {rentFormLabel(lastRent.form)} rent,&nbsp;
            {lastRent.amount_labour_minutes.toLocaleString()}&nbsp;min over{" "}
            {lastRent.period_years}&nbsp;yr(s)
          </div>
        )}
      </section>

      {/* Rents list */}
      <section className="v3-ch37-section" aria-label="Stored ground-rent records">
        <h4>Ground-rent records</h4>
        {rents.length === 0 ? (
          <p className="v3-ch37-subintro">No rents stored yet.</p>
        ) : (
          <div className="v3-ch37-table-wrap">
            <table className="v3-ch37-table">
              <thead>
                <tr>
                  <th>ID</th>
                  <th>Parcel</th>
                  <th>Capitalist</th>
                  <th>Landowner</th>
                  <th>Form</th>
                  <th>Labour-min</th>
                  <th>Years</th>
                </tr>
              </thead>
              <tbody>
                {rents.map((r) => (
                  <tr key={r.id}>
                    <td title={r.id}>{r.id.slice(0, 8)}&hellip;</td>
                    <td title={r.parcel_id}>{r.parcel_id.slice(0, 8)}&hellip;</td>
                    <td>{r.capitalist_id}</td>
                    <td>{r.land_owner_id}</td>
                    <td>{rentFormLabel(r.form)}</td>
                    <td>{r.amount_labour_minutes.toLocaleString()}</td>
                    <td>{r.period_years}</td>
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
