import { useEffect, useState } from "react";
import { financeApi } from "../../api";
import type {
  CounteractingForceKind,
  CounteractingScenarioResponse,
} from "../../types";
import "./Ch14CounteractingForces.css";

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

interface ForceMeta {
  kind: CounteractingForceKind;
  section: string;
  label: string;
  description: string;
  excluded?: boolean;
}

// The six counteracting forces, §I–§VI, in Marx's order.
const FORCES: ForceMeta[] = [
  {
    kind: "intensified_exploitation",
    section: "§I",
    label: "Intensified exploitation",
    description: "More surplus-value is extracted from a given labour-force by raising the rate of exploitation above the average.",
  },
  {
    kind: "depressed_wages",
    section: "§II",
    label: "Depression of wages below value",
    description: "Wages are forced below the value of labour-power, lifting surplus-value without any rise in productivity.",
  },
  {
    kind: "cheaper_constant_capital",
    section: "§III",
    label: "Cheapening of constant capital",
    description: "Rising productivity cheapens the elements of constant capital, so c grows in mass more slowly than in value.",
  },
  {
    kind: "relative_overpopulation",
    section: "§IV",
    label: "Relative over-population",
    description: "The reserve army of labour keeps wages low and supplies cheap labour to new, low-composition branches.",
  },
  {
    kind: "foreign_trade",
    section: "§V",
    label: "Foreign trade",
    description: "Foreign trade cheapens inputs and the means of subsistence and raises the rate of surplus-value.",
  },
  {
    kind: "stock_capital",
    section: "§VI",
    label: "Increase of stock capital",
    description: "Interest-bearing (stock) capital yields interest, not average profit — it does not enter the industrial path.",
    excluded: true,
  },
];

// The five industrial forces are checked by default; stock capital is not, since
// it is excluded from the general rate and adds no uplift.
const DEFAULT_CHECKED: CounteractingForceKind[] = [
  "intensified_exploitation",
  "depressed_wages",
  "cheaper_constant_capital",
  "relative_overpopulation",
  "foreign_trade",
];

export function Ch14CounteractingForces() {
  const [label, setLabel] = useState("All five industrial forces — law still holds");
  const [sRatePct, setSRatePct] = useState(100); // s' as a percentage
  const [steps, setSteps] = useState<StepRow[]>(
    DEFAULT_STEPS.map(([c, v]) => ({ c, v })),
  );
  const [checked, setChecked] = useState<Set<CounteractingForceKind>>(
    new Set(DEFAULT_CHECKED),
  );
  const [scenario, setScenario] = useState<CounteractingScenarioResponse | null>(null);
  const [scenarioError, setScenarioError] = useState<string | null>(null);
  const [stored, setStored] = useState<CounteractingScenarioResponse[]>([]);

  useEffect(() => {
    financeApi
      .listCounteractingScenarios()
      .then((items) => {
        setStored(items);
        if (items.length > 0 && scenario === null) {
          setScenario(items[0]);
        }
      })
      .catch(() => {});
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  function toggleForce(kind: CounteractingForceKind) {
    setChecked((prev) => {
      const next = new Set(prev);
      if (next.has(kind)) {
        next.delete(kind);
      } else {
        next.add(kind);
      }
      return next;
    });
  }

  function setStepField(i: number, field: keyof StepRow, value: number) {
    setSteps((prev) => prev.map((s, j) => (j === i ? { ...s, [field]: value } : s)));
  }

  function addStep() {
    setSteps((prev) => [...prev, { c: prev.length ? prev[prev.length - 1].c + 100 : 100, v: 100 }]);
  }

  function removeStep(i: number) {
    setSteps((prev) => prev.filter((_, j) => j !== i));
  }

  async function buildScenario() {
    setScenarioError(null);
    try {
      const s = await financeApi.createCounteractingScenario({
        label,
        force_kinds: FORCES.filter((f) => checked.has(f.kind)).map((f) => f.kind),
        steps: steps.map((s) => [s.c, s.v] as [number, number]),
        surplus_value_rate: Math.round(sRatePct * 100),
      });
      setScenario(s);
      const items = await financeApi.listCounteractingScenarios();
      setStored(items);
    } catch (e) {
      setScenarioError(String(e));
    }
  }

  const lawHolds = scenario
    ? scenario.final_profit_rate < scenario.initial_profit_rate
    : false;

  return (
    <div className="v3-ch14">
      <section className="v3-ch14-intro">
        <p>
          The same influences that drive the general rate of profit{" "}
          <strong>down</strong> also call forth <strong>counter-effects</strong>{" "}
          that hamper, retard, and partly paralyse the fall. Marx enumerates six:
          intensified exploitation, the depression of wages below value, the
          cheapening of constant capital, relative over-population, foreign trade,
          and the increase of stock capital. They <strong>counteract</strong> the
          law without abolishing it — the rate still falls, only more slowly. The
          sixth, stock (interest-bearing) capital, is <strong>excluded from the
          general rate</strong> entirely and so adds nothing to the industrial
          path.
        </p>
      </section>

      <section className="v3-ch14-section" aria-label="Counteracting forces">
        <h4>The six counteracting forces</h4>
        <ul className="v3-ch14-forces">
          {FORCES.map((f) => (
            <li key={f.kind}>
              <label className={`v3-ch14-force ${checked.has(f.kind) ? "checked" : ""}`}>
                <input
                  type="checkbox"
                  checked={checked.has(f.kind)}
                  onChange={() => toggleForce(f.kind)}
                />
                <span className="v3-ch14-force-text">
                  <span className="v3-ch14-force-label">
                    <span className="v3-ch14-force-section">{f.section} </span>
                    {f.label}
                    {f.excluded && (
                      <span className="v3-ch14-force-excluded">
                        {" "}
                        (excluded from general rate)
                      </span>
                    )}
                  </span>
                  <span className="v3-ch14-force-desc">{f.description}</span>
                </span>
              </label>
            </li>
          ))}
        </ul>
      </section>

      <section className="v3-ch14-section" aria-label="Base trajectory">
        <h4>Apply the forces to a rising-composition trajectory</h4>
        <div className="v3-ch14-form">
          <label className="v3-ch14-grow">
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

        <table className="v3-ch14-steps">
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
                    className="v3-ch14-link"
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

        <div className="v3-ch14-actions">
          <button type="button" className="v3-ch14-link" onClick={addStep}>
            + add period
          </button>
          <button
            type="button"
            className="v3-ch14-btn"
            onClick={buildScenario}
            disabled={checked.size === 0}
          >
            Build scenario
          </button>
        </div>
        {scenarioError && <p className="v3-ch14-error">{scenarioError}</p>}

        {scenario && (
          <div className="v3-ch14-result">
            <table className="v3-ch14-periods">
              <thead>
                <tr>
                  <th>Period</th>
                  <th>c</th>
                  <th>v</th>
                  <th>Base p′</th>
                  <th>Uplift</th>
                  <th>Modified p′</th>
                </tr>
              </thead>
              <tbody>
                {scenario.base_trajectory.periods.map((p, i) => {
                  const base = p.profit_rate;
                  const modified = scenario.modified_rates[i] ?? base;
                  return (
                    <tr key={i}>
                      <td>{i + 1}</td>
                      <td>{p.constant_capital}</td>
                      <td>{p.variable_capital}</td>
                      <td className="v3-ch14-rate-base">{pct(base)}</td>
                      <td className="v3-ch14-uplift">+{pct(modified - base)}</td>
                      <td className="v3-ch14-rate">{pct(modified)}</td>
                    </tr>
                  );
                })}
              </tbody>
            </table>

            <div className={`v3-ch14-summary ${lawHolds ? "holds" : ""}`}>
              <span className="v3-ch14-summary-figure">
                {pct(scenario.initial_profit_rate)} → {pct(scenario.final_profit_rate)}
              </span>
              <span>
                {lawHolds
                  ? "The law is retarded, not abolished: even with every counteracting force active, the final rate of profit remains below the initial rate."
                  : "The final rate did not fall below the initial rate."}
              </span>
            </div>
            <p className="v3-ch14-note">
              Each force adds a bounded per-period uplift to <code>p′</code>, so the
              modified path lies above the bare law of Ch. 13. Yet the rising
              composition still dominates — the counter-tendencies{" "}
              <strong>slow</strong> the fall without reversing it.
            </p>
          </div>
        )}
      </section>

      {stored.length > 0 && (
        <section className="v3-ch14-section" aria-label="Stored scenarios">
          <h4>Stored scenarios</h4>
          <ul className="v3-ch14-stored">
            {stored.map((s) => (
              <li key={s.id}>
                <button
                  type="button"
                  className="v3-ch14-stored-item"
                  onClick={() => setScenario(s)}
                >
                  <span className="v3-ch14-stored-label">{s.label}</span>
                  <span className="v3-ch14-stored-rates">
                    {s.modified_rates.map((r) => pct(r)).join(" → ")}
                  </span>
                </button>
              </li>
            ))}
          </ul>
        </section>
      )}
    </div>
  );
}
