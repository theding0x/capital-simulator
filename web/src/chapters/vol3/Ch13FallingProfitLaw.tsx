import { useEffect, useState } from "react";
import { financeApi } from "../../api";
import type {
  CompositionTrajectoryResponse,
  RateMassContradictionResponse,
} from "../../types";
import "./Ch13FallingProfitLaw.css";

function pct(bp: number): string {
  return `${(bp / 100).toLocaleString("en-GB", { maximumFractionDigits: 2 })}%`;
}

// The §law five-period series: s' = 100%, v = 100 held constant, c rises 50→400.
const DEFAULT_STEPS: [number, number][] = [
  [50, 100],
  [100, 100],
  [200, 100],
  [300, 100],
  [400, 100],
];

interface StepRow {
  c: number;
  v: number;
}

export function Ch13FallingProfitLaw() {
  // --- Trajectory builder ---------------------------------------------------
  const [label, setLabel] = useState("The Law As Such — rising composition");
  const [sRatePct, setSRatePct] = useState(100); // s' as a percentage
  const [steps, setSteps] = useState<StepRow[]>(
    DEFAULT_STEPS.map(([c, v]) => ({ c, v })),
  );
  const [trajectory, setTrajectory] = useState<CompositionTrajectoryResponse | null>(null);
  const [trajectoryError, setTrajectoryError] = useState<string | null>(null);

  // --- Stored trajectories --------------------------------------------------
  const [stored, setStored] = useState<CompositionTrajectoryResponse[]>([]);

  // --- Rate-vs-mass contradiction -------------------------------------------
  const [oldC, setOldC] = useState(1_000_000);
  const [oldRatePct, setOldRatePct] = useState(20); // p' as a percentage
  const [newC, setNewC] = useState(2_000_000);
  const [newRatePct, setNewRatePct] = useState(10);
  const [rateMass, setRateMass] = useState<RateMassContradictionResponse | null>(null);
  const [rateMassError, setRateMassError] = useState<string | null>(null);

  useEffect(() => {
    financeApi
      .listCompositionTrajectories()
      .then((items) => {
        setStored(items);
        if (items.length > 0 && trajectory === null) {
          setTrajectory(items[0]);
        }
      })
      .catch(() => {});
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  function setStepField(i: number, field: keyof StepRow, value: number) {
    setSteps((prev) => prev.map((s, j) => (j === i ? { ...s, [field]: value } : s)));
  }

  function addStep() {
    setSteps((prev) => [...prev, { c: prev.length ? prev[prev.length - 1].c + 100 : 100, v: 100 }]);
  }

  function removeStep(i: number) {
    setSteps((prev) => prev.filter((_, j) => j !== i));
  }

  async function buildTrajectory() {
    setTrajectoryError(null);
    try {
      const t = await financeApi.createCompositionTrajectory({
        label,
        surplus_value_rate: Math.round(sRatePct * 100),
        steps: steps.map((s) => [s.c, s.v] as [number, number]),
      });
      setTrajectory(t);
      const items = await financeApi.listCompositionTrajectories();
      setStored(items);
    } catch (e) {
      setTrajectoryError(String(e));
    }
  }

  async function computeRateMass() {
    setRateMassError(null);
    try {
      const r = await financeApi.computeRateMassContradiction({
        old_c: oldC,
        old_rate: Math.round(oldRatePct * 100),
        new_c: newC,
        new_rate: Math.round(newRatePct * 100),
      });
      setRateMass(r);
    } catch (e) {
      setRateMassError(String(e));
    }
  }

  const maxRate = trajectory
    ? Math.max(...trajectory.periods.map((p) => p.profit_rate), 1)
    : 1;

  return (
    <div className="v3-ch13">
      <section className="v3-ch13-intro">
        <p>
          The <strong>law as such</strong>: with the growing productivity of
          labour the <strong>organic composition of capital rises</strong> —
          constant capital <code>c</code> grows relative to variable capital{" "}
          <code>v</code>. Variable capital is the sole source of surplus-value, so
          with the rate of surplus-value <code>s′</code> held constant the
          surplus <code>s = s′·v</code> no longer keeps pace with the total
          capital <code>C = c + v</code>. The rate of profit{" "}
          <code>p′ = s / (c + v)</code> therefore <strong>falls</strong>. This is
          not a conjunctural fluctuation but a structural tendency of capitalist
          accumulation.
        </p>
      </section>

      <section className="v3-ch13-section" aria-label="Trajectory builder">
        <h4>Rising composition → falling rate of profit</h4>
        <div className="v3-ch13-form">
          <label className="v3-ch13-grow">
            Label
            <input
              type="text"
              value={label}
              onChange={(e) => setLabel(e.target.value)}
            />
          </label>
          <label>
            Rate of surplus-value s′ (%)
            <input
              type="number"
              min={0}
              value={sRatePct}
              onChange={(e) => setSRatePct(Number(e.target.value))}
            />
          </label>
        </div>

        <table className="v3-ch13-steps">
          <thead>
            <tr>
              <th>Period</th>
              <th>c</th>
              <th>v</th>
              <th />
            </tr>
          </thead>
          <tbody>
            {steps.map((s, i) => (
              <tr key={i}>
                <td>{i + 1}</td>
                <td>
                  <input
                    type="number"
                    min={0}
                    value={s.c}
                    onChange={(e) => setStepField(i, "c", Number(e.target.value))}
                  />
                </td>
                <td>
                  <input
                    type="number"
                    min={0}
                    value={s.v}
                    onChange={(e) => setStepField(i, "v", Number(e.target.value))}
                  />
                </td>
                <td>
                  <button
                    type="button"
                    className="v3-ch13-link"
                    onClick={() => removeStep(i)}
                    disabled={steps.length <= 1}
                  >
                    remove
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>

        <div className="v3-ch13-actions">
          <button type="button" className="v3-ch13-link" onClick={addStep}>
            + add period
          </button>
          <button type="button" className="v3-ch13-btn" onClick={buildTrajectory}>
            Build trajectory
          </button>
        </div>
        {trajectoryError && <p className="v3-ch13-error">{trajectoryError}</p>}

        {trajectory && (
          <div className="v3-ch13-result">
            <table className="v3-ch13-periods">
              <thead>
                <tr>
                  <th>Period</th>
                  <th>c</th>
                  <th>v</th>
                  <th>s</th>
                  <th>C = c+v</th>
                  <th>p′</th>
                  <th>Rate of profit</th>
                </tr>
              </thead>
              <tbody>
                {trajectory.periods.map((p, i) => (
                  <tr key={i}>
                    <td>{i + 1}</td>
                    <td>{p.constant_capital}</td>
                    <td>{p.variable_capital}</td>
                    <td>{p.surplus_value}</td>
                    <td>{p.constant_capital + p.variable_capital}</td>
                    <td className="v3-ch13-rate">{pct(p.profit_rate)}</td>
                    <td className="v3-ch13-barcell">
                      <span
                        className="v3-ch13-bar"
                        style={{ width: `${(p.profit_rate / maxRate) * 100}%` }}
                      />
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
            <p className="v3-ch13-note">
              s′ is constant at {pct(trajectory.surplus_value_rate)}, yet the rate
              of profit falls monotonically as the composition rises — the
              numerator <code>s</code> is unchanged while the denominator{" "}
              <code>C</code> grows.
            </p>
          </div>
        )}
      </section>

      {stored.length > 0 && (
        <section className="v3-ch13-section" aria-label="Stored trajectories">
          <h4>Stored trajectories</h4>
          <ul className="v3-ch13-stored">
            {stored.map((t) => (
              <li key={t.id}>
                <button
                  type="button"
                  className="v3-ch13-stored-item"
                  onClick={() => setTrajectory(t)}
                >
                  <span className="v3-ch13-stored-label">{t.label}</span>
                  <span className="v3-ch13-stored-rates">
                    {t.profit_rates.map((r) => pct(r)).join(" → ")}
                  </span>
                </button>
              </li>
            ))}
          </ul>
        </section>
      )}

      <section className="v3-ch13-section" aria-label="Rate versus mass">
        <h4>The rate falls — the mass need not</h4>
        <p className="v3-ch13-subintro">
          As the rate of profit falls the total capital grows, so the absolute{" "}
          <strong>mass of profit</strong> (<code>C·p′</code>) can stay constant or
          even rise. Capital must grow to sustain accumulation despite the falling
          rate.
        </p>
        <div className="v3-ch13-form">
          <label>
            Old C
            <input
              type="number"
              min={0}
              value={oldC}
              onChange={(e) => setOldC(Number(e.target.value))}
            />
          </label>
          <label>
            Old p′ (%)
            <input
              type="number"
              min={0}
              value={oldRatePct}
              onChange={(e) => setOldRatePct(Number(e.target.value))}
            />
          </label>
          <label>
            New C
            <input
              type="number"
              min={0}
              value={newC}
              onChange={(e) => setNewC(Number(e.target.value))}
            />
          </label>
          <label>
            New p′ (%)
            <input
              type="number"
              min={0}
              value={newRatePct}
              onChange={(e) => setNewRatePct(Number(e.target.value))}
            />
          </label>
          <button type="button" className="v3-ch13-btn" onClick={computeRateMass}>
            Compute
          </button>
        </div>
        {rateMassError && <p className="v3-ch13-error">{rateMassError}</p>}

        {rateMass && (
          <div className="v3-ch13-masses">
            <div className="v3-ch13-mass-col">
              <span className="v3-ch13-mass-head">Old</span>
              <span className="v3-ch13-mass-line">C = {rateMass.old_c.toLocaleString("en-GB")}</span>
              <span className="v3-ch13-mass-line">p′ = {pct(rateMass.old_rate)}</span>
              <span className="v3-ch13-mass-value">
                mass = {rateMass.old_mass.toLocaleString("en-GB")}
              </span>
            </div>
            <div className="v3-ch13-mass-arrow">→</div>
            <div className="v3-ch13-mass-col">
              <span className="v3-ch13-mass-head">New</span>
              <span className="v3-ch13-mass-line">C = {rateMass.new_c.toLocaleString("en-GB")}</span>
              <span className="v3-ch13-mass-line">p′ = {pct(rateMass.new_rate)}</span>
              <span className="v3-ch13-mass-value">
                mass = {rateMass.new_mass.toLocaleString("en-GB")}
              </span>
            </div>
            <div
              className={`v3-ch13-mass-change ${
                rateMass.mass_change > 0
                  ? "up"
                  : rateMass.mass_change < 0
                    ? "down"
                    : "flat"
              }`}
            >
              Δ mass = {rateMass.mass_change > 0 ? "+" : ""}
              {rateMass.mass_change.toLocaleString("en-GB")}
              <span className="v3-ch13-mass-gloss">
                {rateMass.mass_change > 0
                  ? "the rate fell, yet the mass of profit grew"
                  : rateMass.mass_change < 0
                    ? "the mass of profit fell"
                    : "the rate fell, yet the mass of profit was preserved"}
              </span>
            </div>
          </div>
        )}
      </section>
    </div>
  );
}
