import { useEffect, useMemo, useState } from "react";
import type { FormEvent } from "react";
import { api } from "../api";
import { usePounds } from "../CurrencyContext";
import type {
  ColonialLabourMarket,
  CreateColonialLabourMarketInput,
  WageWorkerIndependence,
} from "../types";
import "./Ch33ModernColonisation.css";

const DEFAULT_NEW_MARKET: CreateColonialLabourMarketInput = {
  colony: "",
  free_labourers: 300,
  annual_wage_pence: 4800,
  land_available: true,
};

const DEFAULT_PRICE_PENCE = 240;
const DEFAULT_ACRES = 50;

export function Ch33ModernColonisation() {
  const fmt = usePounds();

  const [markets, setMarkets] = useState<ColonialLabourMarket[]>([]);
  const [listError, setListError] = useState("");
  const [selectedId, setSelectedId] = useState("");

  const [newForm, setNewForm] = useState<CreateColonialLabourMarketInput>(DEFAULT_NEW_MARKET);
  const [createError, setCreateError] = useState("");
  const [creating, setCreating] = useState(false);

  const [pricePence, setPricePence] = useState(DEFAULT_PRICE_PENCE);
  const [acres, setAcres] = useState(DEFAULT_ACRES);
  const [regulating, setRegulating] = useState(false);
  const [regulateError, setRegulateError] = useState("");

  const [workerId, setWorkerId] = useState("settler-001");
  const [independence, setIndependence] = useState<WageWorkerIndependence | null>(null);
  const [independenceError, setIndependenceError] = useState("");
  const [computing, setComputing] = useState(false);

  useEffect(() => {
    refresh();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  async function refresh() {
    setListError("");
    try {
      const list = await api.listColonialMarkets();
      setMarkets(list);
      if (list.length > 0 && !selectedId) {
        setSelectedId(list[0].id);
      }
    } catch (e) {
      setListError(String(e));
    }
  }

  const selected = useMemo(
    () => markets.find((m) => m.id === selectedId) ?? null,
    [markets, selectedId],
  );

  async function handleCreate(e: FormEvent) {
    e.preventDefault();
    setCreateError("");
    setCreating(true);
    try {
      const created = await api.createColonialMarket(newForm);
      setNewForm(DEFAULT_NEW_MARKET);
      await refresh();
      setSelectedId(created.id);
    } catch (e) {
      setCreateError(String(e));
    } finally {
      setCreating(false);
    }
  }

  async function handleRegulate() {
    if (!selected) return;
    setRegulateError("");
    setRegulating(true);
    try {
      await api.regulateColonialMarket(selected.id, {
        sufficient_price: {
          price_per_acre_pence: pricePence,
          desired_acres: acres,
        },
      });
      await refresh();
    } catch (e) {
      setRegulateError(String(e));
    } finally {
      setRegulating(false);
    }
  }

  async function handleIndependence() {
    if (!selected) return;
    setIndependenceError("");
    setComputing(true);
    try {
      const result = await api.computeIndependence(selected.id, {
        worker_id: workerId,
        sufficient_price: {
          price_per_acre_pence: pricePence,
          desired_acres: acres,
        },
      });
      setIndependence(result);
    } catch (e) {
      setIndependenceError(String(e));
    } finally {
      setComputing(false);
    }
  }

  return (
    <div className="ch33-grid">
      <section className="ch33-card">
        <h2>The colonial labour market</h2>
        <p className="ch33-quote">
          “Mr. Peel took with him from England to Swan River, West Australia,
          means of subsistence and of production to the amount of £50,000.
          Mr. Peel had the foresight to bring with him, besides, 300 persons
          of the working class … Once arrived at his destination,
          ‘Mr. Peel was left without a servant to make his bed.’”
        </p>
        {listError && <div className="ch33-error">{listError}</div>}
        <ul className="ch33-list">
          {markets.map((m) => (
            <li
              key={m.id}
              className={`ch33-row${m.id === selectedId ? " selected" : ""}`}
              onClick={() => setSelectedId(m.id)}
            >
              <div>
                <div className="ch33-row-title">{m.colony}</div>
                <div className="ch33-row-meta">
                  {m.free_labourers.toLocaleString()} labourers ·{" "}
                  {fmt(m.annual_wage_pence)}/yr ·{" "}
                  {m.wakefield_scheme_applied
                    ? `${m.independence_years}-yr dependence`
                    : "no scheme"}
                </div>
              </div>
              <span
                className={`ch33-pill ${m.surplus_labour_extractable ? "extract" : "free"}`}
              >
                {m.surplus_labour_extractable ? "wage-dependent" : "free"}
              </span>
            </li>
          ))}
        </ul>
      </section>

      <section className="ch33-card">
        <h2>Wakefield's remedy</h2>
        {selected ? (
          <>
            <p className="ch33-quote">
              “The price of the soil imposed by the State must … be a
              ‘sufficient price’ — so high as to prevent the labourers from
              becoming independent landowners until others had followed to
              take their place.”
            </p>
            <div className="ch33-form-grid">
              <div className="ch33-field">
                <label htmlFor="price-pence">Price per acre (pence)</label>
                <input
                  id="price-pence"
                  type="number"
                  min={0}
                  value={pricePence}
                  onChange={(e) => setPricePence(Number(e.target.value))}
                />
              </div>
              <div className="ch33-field">
                <label htmlFor="acres">Desired acres</label>
                <input
                  id="acres"
                  type="number"
                  min={0}
                  value={acres}
                  onChange={(e) => setAcres(Number(e.target.value))}
                />
              </div>
            </div>
            <div className="ch33-actions">
              <button
                type="button"
                disabled={regulating || !selected}
                onClick={handleRegulate}
              >
                {regulating ? "Applying…" : "Apply sufficient price"}
              </button>
            </div>
            {regulateError && <div className="ch33-error">{regulateError}</div>}
            <div className="ch33-comparison">
              <div className="ch33-side">
                <h3>Without a scheme</h3>
                <p className="ch33-stat">
                  {selected.wakefield_scheme_applied ? "—" : "Free settler"}
                </p>
                <p className="ch33-row-meta">
                  Labour walks off into independent production.
                </p>
              </div>
              <div className="ch33-side">
                <h3>With a sufficient price</h3>
                <p className="ch33-stat">
                  {selected.wakefield_scheme_applied
                    ? `${selected.independence_years} yrs dependence`
                    : "—"}
                </p>
                <p className="ch33-row-meta">
                  Wage relation reproduced; land fund imports replacements.
                </p>
              </div>
            </div>
          </>
        ) : (
          <p className="ch33-row-meta">Select a colony to apply Wakefield's scheme.</p>
        )}
      </section>

      <section className="ch33-card ch33-form">
        <h2>Register a colony</h2>
        <form onSubmit={handleCreate}>
          <div className="ch33-form-grid">
            <div className="ch33-field">
              <label htmlFor="new-colony">Colony</label>
              <input
                id="new-colony"
                type="text"
                value={newForm.colony}
                onChange={(e) => setNewForm({ ...newForm, colony: e.target.value })}
                placeholder="e.g. Swan River, 1829"
                required
              />
            </div>
            <div className="ch33-field">
              <label htmlFor="new-free">Free labourers</label>
              <input
                id="new-free"
                type="number"
                min={0}
                value={newForm.free_labourers}
                onChange={(e) =>
                  setNewForm({ ...newForm, free_labourers: Number(e.target.value) })
                }
              />
            </div>
            <div className="ch33-field">
              <label htmlFor="new-wage">Annual wage (pence)</label>
              <input
                id="new-wage"
                type="number"
                min={1}
                value={newForm.annual_wage_pence}
                onChange={(e) =>
                  setNewForm({ ...newForm, annual_wage_pence: Number(e.target.value) })
                }
              />
            </div>
            <div className="ch33-field">
              <label htmlFor="new-land">Land available</label>
              <input
                id="new-land"
                type="checkbox"
                checked={newForm.land_available}
                onChange={(e) =>
                  setNewForm({ ...newForm, land_available: e.target.checked })
                }
              />
            </div>
          </div>
          <div className="ch33-actions">
            <button type="submit" disabled={creating}>
              {creating ? "Registering…" : "Register colony"}
            </button>
          </div>
          {createError && <div className="ch33-error">{createError}</div>}
        </form>
      </section>

      <section className="ch33-card ch33-independence">
        <h2>Per-worker independence projection</h2>
        <p className="ch33-quote">
          “No labourer would be able to procure land until he had worked for
          money … all immigrant labourers, working for a time for wages and
          in combination, would produce capital for the employment of more
          labourers.”
        </p>
        <div className="ch33-form-grid">
          <div className="ch33-field">
            <label htmlFor="worker-id">Worker ID</label>
            <input
              id="worker-id"
              type="text"
              value={workerId}
              onChange={(e) => setWorkerId(e.target.value)}
            />
          </div>
        </div>
        <div className="ch33-actions">
          <button
            type="button"
            disabled={computing || !selected}
            onClick={handleIndependence}
          >
            {computing ? "Computing…" : "Project independence"}
          </button>
        </div>
        {independenceError && <div className="ch33-error">{independenceError}</div>}
        {independence && (
          <div className="ch33-readout">
            At <strong>{fmt(independence.annual_savings_pence)}/yr</strong> in
            savings the labourer reaches the ransom of{" "}
            <strong>{fmt(independence.target_ransom_pence)}</strong> after{" "}
            <strong>{independence.years_worked}</strong> years.{" "}
            {independence.became_landowner
              ? "Independence: achieved."
              : "Independence: not reached within the limit."}
          </div>
        )}
      </section>
    </div>
  );
}
