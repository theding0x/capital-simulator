import { useEffect, useState } from "react";
import type { FormEvent } from "react";
import { api } from "../api";
import { usePounds } from "../CurrencyContext";
import type {
  HistoricalStage,
  PrimitiveAccumulation,
  SeedScenarioResult,
} from "../types";

interface EpisodeDraft {
  period: string;
  method: string;
  labourers_expropriated: string;
  capital_formed: string;
}

const englandPreset: EpisodeDraft[] = [
  { period: "1450–1640", method: "enclosure of common lands", labourers_expropriated: "40000", capital_formed: "80000" },
  { period: "1500–1640", method: "expropriation of church property", labourers_expropriated: "15000", capital_formed: "30000" },
  { period: "1500–1800", method: "colonial plunder and slave trade", labourers_expropriated: "5000", capital_formed: "90000" },
];

const emptyEpisode: EpisodeDraft = { period: "", method: "", labourers_expropriated: "0", capital_formed: "0" };

export function Ch26PrimitiveAccumulation() {
  const fmt = usePounds();

  // existing stages
  const [stages, setStages] = useState<HistoricalStage[]>([]);
  const [listError, setListError] = useState("");
  const [loadingStages, setLoadingStages] = useState(false);

  // create-stage form
  const [stageName, setStageName] = useState("England 15th–18th c.");
  const [stageDescription, setStageDescription] = useState(
    "Agrarian expropriation creating a free proletariat — enclosure, statutory wage-fixing, and colonial plunder together found capitalist production.",
  );
  const [episodes, setEpisodes] = useState<EpisodeDraft[]>(englandPreset);
  const [createError, setCreateError] = useState("");
  const [creating, setCreating] = useState(false);

  // seed-scenario form
  const [seedStageId, setSeedStageId] = useState<string>("");
  const [preCapitalistWorkers, setPreCapitalistWorkers] = useState(100000);
  const [compositionRatio, setCompositionRatio] = useState(0.8);
  const [surplusRate, setSurplusRate] = useState(1.0);
  const [periods, setPeriods] = useState(3);
  const [seedResult, setSeedResult] = useState<SeedScenarioResult | null>(null);
  const [seedError, setSeedError] = useState("");
  const [seeding, setSeeding] = useState(false);

  async function refreshStages() {
    setListError("");
    setLoadingStages(true);
    try {
      const list = await api.listHistoricalStages();
      setStages(list);
      if (!seedStageId && list.length > 0) {
        setSeedStageId(list[0].id);
      }
    } catch (err) {
      setListError(err instanceof Error ? err.message : String(err));
    } finally {
      setLoadingStages(false);
    }
  }

  useEffect(() => {
    refreshStages();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  function updateEpisode(i: number, patch: Partial<EpisodeDraft>) {
    setEpisodes((prev) => prev.map((e, idx) => (idx === i ? { ...e, ...patch } : e)));
  }

  function addEpisode() {
    setEpisodes((prev) => [...prev, { ...emptyEpisode }]);
  }

  function removeEpisode(i: number) {
    setEpisodes((prev) => prev.filter((_, idx) => idx !== i));
  }

  function applyEnglandPreset() {
    setStageName("England 15th–18th c.");
    setStageDescription(
      "Agrarian expropriation creating a free proletariat — enclosure, statutory wage-fixing, and colonial plunder together found capitalist production.",
    );
    setEpisodes(englandPreset.map((e) => ({ ...e })));
    setCreateError("");
  }

  async function handleCreate(e: FormEvent) {
    e.preventDefault();
    setCreateError("");
    setCreating(true);
    try {
      const payload: PrimitiveAccumulation[] = episodes
        .filter((ep) => ep.period.trim() !== "" && ep.method.trim() !== "")
        .map((ep) => ({
          period: ep.period,
          method: ep.method,
          labourers_expropriated: Number(ep.labourers_expropriated) || 0,
          capital_formed: Number(ep.capital_formed) || 0,
        }));
      if (payload.length === 0) {
        throw new Error("at least one episode with period + method is required");
      }
      const created = await api.createHistoricalStage({
        name: stageName,
        description: stageDescription,
        primitive_accumulations: payload,
      });
      setSeedStageId(created.id);
      await refreshStages();
    } catch (err) {
      setCreateError(err instanceof Error ? err.message : String(err));
    } finally {
      setCreating(false);
    }
  }

  async function handleSeed(e: FormEvent) {
    e.preventDefault();
    setSeedError("");
    setSeedResult(null);
    setSeeding(true);
    try {
      if (!seedStageId) {
        throw new Error("pick a stage first");
      }
      const r = await api.seedScenarioFromStage(seedStageId, {
        pre_capitalist_workers: preCapitalistWorkers,
        composition_ratio: compositionRatio,
        surplus_rate: surplusRate,
        periods,
      });
      setSeedResult(r);
    } catch (err) {
      setSeedError(err instanceof Error ? err.message : String(err));
    } finally {
      setSeeding(false);
    }
  }

  return (
    <div className="chapter-panel">
      <h2>Ch. 26 — The Secret of Primitive Accumulation</h2>
      <p className="chapter-blurb">
        Capitalist production presupposes a class of free wage-labourers and a concentrated stock of capital — but neither
        is given by nature. They are the product of a historical process that separates producers from the means of
        production: <em>primitive accumulation</em>. Capital is not a thing but a social relation; these numbers are the
        record of expropriation that founds it.
      </p>

      {/* Existing stages */}
      <section className="sim-section">
        <h3>Historical Stages</h3>
        <p>Named epochs from which a starting capital stock and free-proletarian supply are derived.</p>
        {listError && <p className="error">{listError}</p>}
        {loadingStages ? (
          <p className="muted">Loading…</p>
        ) : stages.length === 0 ? (
          <p className="muted">No stages yet — create one below.</p>
        ) : (
          <table className="result-table">
            <thead>
              <tr>
                <th>Name</th>
                <th>Episodes</th>
                <th>Total capital formed</th>
                <th>Total labourers dispossessed</th>
              </tr>
            </thead>
            <tbody>
              {stages.map((s) => (
                <tr key={s.id}>
                  <td>
                    <strong>{s.name}</strong>
                    {s.description && <div className="muted small">{s.description}</div>}
                  </td>
                  <td>{s.primitive_accumulations.length}</td>
                  <td>{fmt(s.total_capital_formed)}</td>
                  <td>{s.total_labourers_expropriated.toLocaleString()}</td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </section>

      {/* Create stage */}
      <section className="sim-section">
        <h3>Define a Historical Stage</h3>
        <p>
          Stage = a name + description + one or more primitive-accumulation episodes (period, method, labourers dispossessed,
          capital formed). Every dispossessed producer becomes a free proletarian; capital formed is always positive.
        </p>
        <div className="preset-row">
          <button type="button" onClick={applyEnglandPreset}>England 15th–18th c. preset</button>
        </div>
        <form onSubmit={handleCreate} className="sim-form">
          <label>
            Stage name
            <input type="text" value={stageName} onChange={(e) => setStageName(e.target.value)} required />
          </label>
          <label>
            Description
            <textarea value={stageDescription} onChange={(e) => setStageDescription(e.target.value)} rows={2} />
          </label>

          <table className="result-table">
            <thead>
              <tr>
                <th>Period</th>
                <th>Method</th>
                <th>Labourers expropriated</th>
                <th>Capital formed (£)</th>
                <th></th>
              </tr>
            </thead>
            <tbody>
              {episodes.map((ep, i) => (
                <tr key={i}>
                  <td>
                    <input type="text" value={ep.period} onChange={(e) => updateEpisode(i, { period: e.target.value })} />
                  </td>
                  <td>
                    <input type="text" value={ep.method} onChange={(e) => updateEpisode(i, { method: e.target.value })} />
                  </td>
                  <td>
                    <input
                      type="number"
                      min={0}
                      value={ep.labourers_expropriated}
                      onChange={(e) => updateEpisode(i, { labourers_expropriated: e.target.value })}
                    />
                  </td>
                  <td>
                    <input
                      type="number"
                      min={1}
                      value={ep.capital_formed}
                      onChange={(e) => updateEpisode(i, { capital_formed: e.target.value })}
                    />
                  </td>
                  <td>
                    <button type="button" onClick={() => removeEpisode(i)} disabled={episodes.length <= 1}>
                      remove
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
          <div className="preset-row">
            <button type="button" onClick={addEpisode}>+ add episode</button>
            <button type="submit" disabled={creating}>{creating ? "Saving…" : "Create stage"}</button>
          </div>
        </form>
        {createError && <p className="error">{createError}</p>}
      </section>

      {/* Seed scenario from stage */}
      <section className="sim-section">
        <h3>Seed a Reproduction Scenario</h3>
        <p>
          Project a stage onto a starting <code>CapitalStock</code> (via organic composition) and a{" "}
          <code>ProducerSeparation</code> (free proletarians = dispossessed), then run simple reproduction.
        </p>
        <form onSubmit={handleSeed} className="sim-form">
          <label>
            Stage
            <select value={seedStageId} onChange={(e) => setSeedStageId(e.target.value)} disabled={stages.length === 0}>
              {stages.map((s) => (
                <option key={s.id} value={s.id}>
                  {s.name}
                </option>
              ))}
            </select>
          </label>
          <label>
            Pre-capitalist workers
            <input
              type="number"
              min={0}
              value={preCapitalistWorkers}
              onChange={(e) => setPreCapitalistWorkers(Number(e.target.value))}
            />
          </label>
          <label>
            Organic composition ratio (c / c+v)
            <input
              type="number"
              step="0.01"
              min={0}
              max={0.99}
              value={compositionRatio}
              onChange={(e) => setCompositionRatio(Number(e.target.value))}
            />
          </label>
          <label>
            Surplus-value rate (s′)
            <input
              type="number"
              step="0.01"
              min={0}
              value={surplusRate}
              onChange={(e) => setSurplusRate(Number(e.target.value))}
            />
          </label>
          <label>
            Periods
            <input type="number" min={1} value={periods} onChange={(e) => setPeriods(Number(e.target.value))} />
          </label>
          <button type="submit" disabled={seeding || stages.length === 0}>
            {seeding ? "Running…" : "Seed scenario"}
          </button>
        </form>
        {seedError && <p className="error">{seedError}</p>}
        {seedResult && (
          <>
            <table className="result-table">
              <tbody>
                <tr>
                  <th>Stage</th>
                  <td>{seedResult.stage_name}</td>
                </tr>
                <tr>
                  <th>Seeded constant capital</th>
                  <td>{fmt(seedResult.seeded_capital.constant_capital)}</td>
                </tr>
                <tr>
                  <th>Seeded variable capital</th>
                  <td>{fmt(seedResult.seeded_capital.variable_capital)}</td>
                </tr>
                <tr>
                  <th>Pre-capitalist workers</th>
                  <td>{seedResult.separation.pre_capitalist_workers.toLocaleString()}</td>
                </tr>
                <tr>
                  <th>Displaced workers</th>
                  <td>{seedResult.separation.displaced_workers.toLocaleString()}</td>
                </tr>
                <tr>
                  <th>Free proletarians</th>
                  <td>{seedResult.separation.free_proletarians.toLocaleString()}</td>
                </tr>
              </tbody>
            </table>
            <h4>Reproduction series</h4>
            <table className="result-table">
              <thead>
                <tr>
                  <th>Period</th>
                  <th>Constant</th>
                  <th>Variable</th>
                  <th>Surplus produced</th>
                  <th>Revenue</th>
                </tr>
              </thead>
              <tbody>
                {seedResult.cycles.map((c) => (
                  <tr key={c.period}>
                    <td>{c.period}</td>
                    <td>{fmt(c.capital.constant_capital)}</td>
                    <td>{fmt(c.capital.variable_capital)}</td>
                    <td>{fmt(c.fund.total)}</td>
                    <td>{fmt(c.fund.revenue)}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </>
        )}
      </section>
    </div>
  );
}
