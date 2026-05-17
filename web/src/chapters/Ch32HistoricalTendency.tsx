import { useEffect, useMemo, useState } from "react";
import type { FormEvent } from "react";
import { api } from "../api";
import { usePounds } from "../CurrencyContext";
import type {
  AccumulationTrajectory,
  Negation,
  RunCentralisationInput,
} from "../types";
import "./Ch32HistoricalTendency.css";

const STAGE_LABEL: Record<string, string> = {
  "petty-property": "Petty property",
  "capitalist-expropriation": "Capitalist private property",
  "socialised-property": "Socialised property",
};

export function Ch32HistoricalTendency() {
  const fmt = usePounds();

  const [stages, setStages] = useState<Negation[]>([]);
  const [stagesError, setStagesError] = useState("");

  const [trajectories, setTrajectories] = useState<AccumulationTrajectory[]>([]);
  const [trajectoryError, setTrajectoryError] = useState("");
  const [selectedTrajectoryId, setSelectedTrajectoryId] = useState("");

  const [simForm, setSimForm] = useState<RunCentralisationInput>({
    firms: 1000,
    total_capital_pence: 480_000_000,
    wage_labourers: 50_000,
    steps: 6,
    absorption_rate: 0.2,
    capital_growth_rate: 0.05,
  });
  const [simResult, setSimResult] = useState<AccumulationTrajectory | null>(null);
  const [simError, setSimError] = useState("");
  const [simRunning, setSimRunning] = useState(false);

  useEffect(() => {
    api
      .getNegationOfNegation()
      .then((r) => setStages(r.stages))
      .catch((e) => setStagesError(String(e)));
    refreshTrajectories();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  async function refreshTrajectories() {
    setTrajectoryError("");
    try {
      const list = await api.listAccumulationTrajectories();
      setTrajectories(list);
      if (list.length > 0 && !selectedTrajectoryId) {
        setSelectedTrajectoryId(list[0].id);
      }
    } catch (e) {
      setTrajectoryError(String(e));
    }
  }

  async function handleRunSim(e: FormEvent) {
    e.preventDefault();
    setSimError("");
    setSimRunning(true);
    try {
      const result = await api.runCentralisation(simForm);
      setSimResult(result);
    } catch (e) {
      setSimError(String(e));
    } finally {
      setSimRunning(false);
    }
  }

  function updateForm<K extends keyof RunCentralisationInput>(
    key: K,
    value: RunCentralisationInput[K],
  ) {
    setSimForm((prev) => ({ ...prev, [key]: value }));
  }

  const selectedTrajectory = useMemo(
    () => trajectories.find((t) => t.id === selectedTrajectoryId) ?? null,
    [trajectories, selectedTrajectoryId],
  );

  return (
    <div className="ch32-grid">
      <section className="ch32-card ch32-stages">
        <h2>The dialectic of expropriation</h2>
        <p className="muted small">
          The negation of the negation. Marx names three stages: petty property
          based on the labourer&apos;s own use of the means of labour;
          capitalist private property resting on the exploitation of others;
          and a final socialised property in which the expropriators are
          themselves expropriated.
        </p>
        {stagesError && <p className="error">{stagesError}</p>}
        <ol className="ch32-stage-list">
          {stages.map((stage, i) => (
            <li key={stage.stage}>
              <span className="ch32-stage-index">§{i + 1}</span>
              <div>
                <div className="ch32-stage-name">
                  {STAGE_LABEL[stage.stage] ?? stage.stage}
                </div>
                <div className="ch32-stage-desc small">{stage.description}</div>
              </div>
            </li>
          ))}
        </ol>
      </section>

      <section className="ch32-card ch32-trajectories">
        <h2>Long-run centralisation</h2>
        <p className="muted small">
          Each trajectory is one historical series of capitalists absorbing
          capitalists. The number of firms declines; the average capital per
          surviving firm grows. The reserve army grows in step.
        </p>
        {trajectoryError && <p className="error">{trajectoryError}</p>}
        {trajectories.length === 0 ? (
          <p className="muted small">No persisted trajectories yet.</p>
        ) : (
          <>
            <label className="ch32-field">
              <span>Trajectory</span>
              <select
                value={selectedTrajectoryId}
                onChange={(e) => setSelectedTrajectoryId(e.target.value)}
              >
                {trajectories.map((t) => (
                  <option key={t.id} value={t.id}>
                    {t.name}
                  </option>
                ))}
              </select>
            </label>
            {selectedTrajectory && (
              <TrajectoryView
                trajectory={selectedTrajectory}
                fmt={fmt}
              />
            )}
          </>
        )}
      </section>

      <section className="ch32-card ch32-sim">
        <h2>Run your own centralisation</h2>
        <p className="muted small">
          Stateless simulation. Provide a starting population of firms and how
          aggressively competitors absorb one another each step. The trajectory
          approaches monopoly asymptotically, never reaching it.
        </p>
        <form onSubmit={handleRunSim} className="ch32-form">
          <label className="ch32-field">
            <span>Firms</span>
            <input
              type="number"
              min={2}
              value={simForm.firms}
              onChange={(e) =>
                updateForm("firms", parseInt(e.target.value || "0", 10))
              }
            />
          </label>
          <label className="ch32-field">
            <span>Total capital (pence)</span>
            <input
              type="number"
              min={1}
              value={simForm.total_capital_pence}
              onChange={(e) =>
                updateForm(
                  "total_capital_pence",
                  parseInt(e.target.value || "0", 10),
                )
              }
            />
          </label>
          <label className="ch32-field">
            <span>Wage labourers</span>
            <input
              type="number"
              min={0}
              value={simForm.wage_labourers}
              onChange={(e) =>
                updateForm(
                  "wage_labourers",
                  parseInt(e.target.value || "0", 10),
                )
              }
            />
          </label>
          <label className="ch32-field">
            <span>Steps</span>
            <input
              type="number"
              min={1}
              value={simForm.steps}
              onChange={(e) =>
                updateForm("steps", parseInt(e.target.value || "0", 10))
              }
            />
          </label>
          <label className="ch32-field">
            <span>Absorption rate (0&ndash;1]</span>
            <input
              type="number"
              step="0.01"
              min={0.01}
              max={1}
              value={simForm.absorption_rate}
              onChange={(e) =>
                updateForm(
                  "absorption_rate",
                  parseFloat(e.target.value || "0"),
                )
              }
            />
          </label>
          <label className="ch32-field">
            <span>Capital growth per step</span>
            <input
              type="number"
              step="0.01"
              min={0}
              value={simForm.capital_growth_rate}
              onChange={(e) =>
                updateForm(
                  "capital_growth_rate",
                  parseFloat(e.target.value || "0"),
                )
              }
            />
          </label>
          <button type="submit" disabled={simRunning} className="ch32-run">
            {simRunning ? "Running..." : "Run centralisation"}
          </button>
        </form>
        {simError && <p className="error">{simError}</p>}
        {simResult && <TrajectoryView trajectory={simResult} fmt={fmt} />}
      </section>
    </div>
  );
}

interface TrajectoryViewProps {
  trajectory: AccumulationTrajectory;
  fmt: (pence: number) => string;
}

function TrajectoryView({ trajectory, fmt }: TrajectoryViewProps) {
  const maxConcentrated = Math.max(
    1,
    ...trajectory.steps.map((s) => s.capital_concentrated_pence),
  );

  let runningFirms = trajectory.initial_firms;

  return (
    <div className="ch32-traj">
      <header className="ch32-traj-header">
        <div>
          <div className="muted small">Initial firms</div>
          <div className="ch32-stat">{trajectory.initial_firms}</div>
        </div>
        <div>
          <div className="muted small">Final firms</div>
          <div className="ch32-stat ch32-stat-final">
            {trajectory.final_firms}
          </div>
        </div>
        <div>
          <div className="muted small">Final capital</div>
          <div className="ch32-stat">{fmt(trajectory.final_capital_pence)}</div>
        </div>
        <div>
          <div className="muted small">Reserve army</div>
          <div className="ch32-stat">{trajectory.reserve_army_size}</div>
        </div>
      </header>
      {trajectory.steps.length === 0 ? (
        <p className="muted small">No absorption steps recorded.</p>
      ) : (
        <table className="ch32-steps">
          <thead>
            <tr>
              <th>Step</th>
              <th>Firms absorbed</th>
              <th>Firms remaining</th>
              <th>Capital concentrated</th>
              <th>Intensity</th>
            </tr>
          </thead>
          <tbody>
            {trajectory.steps.map((s) => {
              runningFirms = Math.max(0, runningFirms - s.firms_absorbed);
              const intensity =
                (s.capital_concentrated_pence / maxConcentrated) * 100;
              const firmsPct =
                (runningFirms / Math.max(1, trajectory.initial_firms)) * 100;
              return (
                <tr key={s.step_index}>
                  <td>{s.step_index}</td>
                  <td>{s.firms_absorbed}</td>
                  <td>
                    <div className="ch32-bar ch32-bar-firms">
                      <div
                        className="ch32-bar-fill"
                        style={{ width: `${firmsPct}%` }}
                      />
                      <span className="ch32-bar-label">{runningFirms}</span>
                    </div>
                  </td>
                  <td>{fmt(s.capital_concentrated_pence)}</td>
                  <td>
                    <div className="ch32-bar ch32-bar-capital">
                      <div
                        className="ch32-bar-fill"
                        style={{ width: `${intensity}%` }}
                      />
                    </div>
                  </td>
                </tr>
              );
            })}
          </tbody>
        </table>
      )}
    </div>
  );
}
