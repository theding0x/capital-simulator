import { useEffect, useState } from "react";
import type { FormEvent } from "react";
import { api } from "../api";
import type { EnclosureEvent } from "../types";

const presets = [
  {
    period: "1450–1540",
    acres_enclosed: 500000,
    population_displaced: 40000,
    beneficiary: "landed gentry",
  },
  {
    period: "1760–1820",
    acres_enclosed: 3000000,
    population_displaced: 150000,
    beneficiary: "improving landlords (Parliamentary enclosures)",
  },
  {
    period: "1814–1820",
    acres_enclosed: 794000,
    population_displaced: 15000,
    beneficiary: "Countess of Sutherland",
  },
];

export function Ch27EnclosureEvents() {
  const [events, setEvents] = useState<EnclosureEvent[]>([]);
  const [listError, setListError] = useState("");
  const [loading, setLoading] = useState(false);

  const [period, setPeriod] = useState("1814–1820");
  const [acres, setAcres] = useState("794000");
  const [population, setPopulation] = useState("15000");
  const [beneficiary, setBeneficiary] = useState("Countess of Sutherland");
  const [createError, setCreateError] = useState("");
  const [creating, setCreating] = useState(false);

  async function refreshEvents() {
    setListError("");
    setLoading(true);
    try {
      const list = await api.listEnclosureEvents();
      setEvents(list);
    } catch (err) {
      setListError(err instanceof Error ? err.message : String(err));
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    refreshEvents();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  function applyPreset(p: (typeof presets)[0]) {
    setPeriod(p.period);
    setAcres(String(p.acres_enclosed));
    setPopulation(String(p.population_displaced));
    setBeneficiary(p.beneficiary);
  }

  async function handleCreate(e: FormEvent) {
    e.preventDefault();
    setCreateError("");
    setCreating(true);
    try {
      await api.createEnclosureEvent({
        period,
        acres_enclosed: parseInt(acres, 10),
        population_displaced: parseInt(population, 10),
        beneficiary,
      });
      await refreshEvents();
      setPeriod("");
      setAcres("");
      setPopulation("0");
      setBeneficiary("");
    } catch (err) {
      setCreateError(err instanceof Error ? err.message : String(err));
    } finally {
      setCreating(false);
    }
  }

  const totalAcres = events.reduce((s, e) => s + e.acres_enclosed, 0);
  const totalDisplaced = events.reduce((s, e) => s + e.population_displaced, 0);

  return (
    <div className="chapter-panel">
      <section className="panel-section">
        <h2>Record an Enclosure Event</h2>
        <p className="muted small">
          Each row is one wave of expropriation — land seized from peasants and commons,
          population driven off to become "free" wage-labourers.
        </p>

        <div className="preset-buttons">
          {presets.map((p) => (
            <button
              key={p.period}
              type="button"
              className="btn-secondary btn-sm"
              onClick={() => applyPreset(p)}
            >
              {p.period}
            </button>
          ))}
        </div>

        <form onSubmit={handleCreate} className="form-grid">
          <label>
            Period
            <input
              value={period}
              onChange={(e) => setPeriod(e.target.value)}
              placeholder="e.g. 1814–1820"
              required
            />
          </label>
          <label>
            Acres Enclosed
            <input
              type="number"
              min="1"
              value={acres}
              onChange={(e) => setAcres(e.target.value)}
              required
            />
          </label>
          <label>
            Population Displaced
            <input
              type="number"
              min="0"
              value={population}
              onChange={(e) => setPopulation(e.target.value)}
              required
            />
          </label>
          <label>
            Beneficiary
            <input
              value={beneficiary}
              onChange={(e) => setBeneficiary(e.target.value)}
              placeholder="e.g. Countess of Sutherland"
              required
            />
          </label>
          {createError && <p className="error-msg">{createError}</p>}
          <button type="submit" className="btn-primary" disabled={creating}>
            {creating ? "Recording…" : "Record"}
          </button>
        </form>
      </section>

      <section className="panel-section">
        <h2>Displacement Timeline</h2>
        {listError && <p className="error-msg">{listError}</p>}
        {loading && <p className="muted">Loading…</p>}
        {!loading && events.length === 0 && (
          <p className="muted small">No enclosure events recorded yet.</p>
        )}
        {events.length > 0 && (
          <>
            <table className="data-table">
              <thead>
                <tr>
                  <th>Period</th>
                  <th>Acres Enclosed</th>
                  <th>Population Displaced</th>
                  <th>Beneficiary</th>
                </tr>
              </thead>
              <tbody>
                {events.map((ev) => (
                  <tr key={ev.id}>
                    <td>{ev.period}</td>
                    <td>{ev.acres_enclosed.toLocaleString()}</td>
                    <td>{ev.population_displaced.toLocaleString()}</td>
                    <td>{ev.beneficiary}</td>
                  </tr>
                ))}
              </tbody>
              <tfoot>
                <tr>
                  <th>Total</th>
                  <th>{totalAcres.toLocaleString()}</th>
                  <th>{totalDisplaced.toLocaleString()}</th>
                  <th></th>
                </tr>
              </tfoot>
            </table>
            <p className="muted small" style={{ marginTop: "0.5rem" }}>
              {totalDisplaced.toLocaleString()} producers expropriated across{" "}
              {totalAcres.toLocaleString()} acres — every one reappears on the labour
              market as a "free" wage-labourer.
            </p>
          </>
        )}
      </section>
    </div>
  );
}