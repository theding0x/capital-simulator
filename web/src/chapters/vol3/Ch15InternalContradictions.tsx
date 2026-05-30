import { useEffect, useState } from "react";
import { financeApi } from "../../api";
import type {
  ContradictionKind,
  CrisisResponse,
  InternalContradictionResponse,
  PartIIISummaryResponse,
} from "../../types";
import "./Ch15InternalContradictions.css";

const BASIS_POINTS = 10000;

function pct(bp: number): string {
  return `${(bp / 100).toLocaleString("en-GB", { maximumFractionDigits: 2 })}%`;
}

function num(n: number): string {
  return n.toLocaleString("en-GB");
}

interface KindMeta {
  kind: ContradictionKind;
  section: string;
  label: string;
  description: string;
}

// The four faces of the law's internal contradiction, §I–§IV, in Marx's order.
const KINDS: KindMeta[] = [
  {
    kind: "falling_rate_rising_mass",
    section: "§I",
    label: "Falling rate, rising mass",
    description:
      "The falling rate and the rising mass of profit are two sides of the same movement: both express the tendential fall in the rate of profit.",
  },
  {
    kind: "overproduction_underconsumption",
    section: "§II",
    label: "Over-production / under-consumption",
    description:
      "Rising productivity expands the mass of commodities beyond the limits of consumption, producing periodic crises of over-production.",
  },
  {
    kind: "overaccumulation_relative_overpopulation",
    section: "§III",
    label: "Over-accumulation / relative over-population",
    description:
      "Excess capital cannot be deployed at the average rate of profit; the surplus population grows as the organic composition rises.",
  },
  {
    kind: "concentration_centralisation",
    section: "§IV",
    label: "Concentration & centralisation",
    description:
      "The fall in the rate of profit accelerates the concentration and centralisation of capital: small capitals are squeezed out and swallowed by large ones.",
  },
];

// restoredProfitRate mirrors tendency.restoredProfitRate (round-half-up): writing
// down a fraction w of constant capital shrinks the denominator and lifts p′.
function restoredProfitRate(pre: number, writedown: number): number {
  const denom = BASIS_POINTS - writedown * 100;
  if (denom <= 0) return pre;
  const numer = pre * BASIS_POINTS;
  return Math.floor((numer + Math.floor(denom / 2)) / denom);
}

// concentrationPath mirrors tendency.CapitalConcentration.Path / FinalC:
// each period C ← C + ⌊C·s′/10000⌋, with C₀ included as element [0].
function concentrationPath(initialC: number, sRate: number, periods: number): number[] {
  const path: number[] = [initialC];
  let capital = initialC;
  for (let i = 0; i < periods; i++) {
    capital += Math.floor((capital * sRate) / BASIS_POINTS);
    path.push(capital);
  }
  return path;
}

export function Ch15InternalContradictions() {
  // §I–§IV taxonomy + persisted contradictions.
  const [contradictions, setContradictions] = useState<InternalContradictionResponse[]>([]);
  const [contradictionError, setContradictionError] = useState<string | null>(null);

  // §II/§III crisis events.
  const [writedown, setWritedown] = useState(30);
  const [preCrisis, setPreCrisis] = useState(1000);
  const [crises, setCrises] = useState<CrisisResponse[]>([]);
  const [crisisError, setCrisisError] = useState<string | null>(null);

  // §IV concentration calculator (client-side; value object not persisted).
  const [initialC, setInitialC] = useState(1_000_000);
  const [concSRate, setConcSRate] = useState(2000);
  const [periods, setPeriods] = useState(3);

  // §IV centralisation calculator (client-side).
  const [absorbedCapitals, setAbsorbedCapitals] = useState(10);
  const [absorbedSize, setAbsorbedSize] = useState(100_000);

  // Part III summary.
  const [summary, setSummary] = useState<PartIIISummaryResponse | null>(null);

  async function refreshCrises() {
    const items = await financeApi.listCrises();
    setCrises(items);
  }

  async function refreshSummary() {
    const s = await financeApi.getPartIIISummary();
    setSummary(s);
  }

  useEffect(() => {
    refreshCrises().catch(() => {});
    refreshSummary().catch(() => {});
  }, []);

  async function recordContradiction(kind: ContradictionKind) {
    setContradictionError(null);
    try {
      const c = await financeApi.createContradiction({ kind });
      setContradictions((prev) => [c, ...prev]);
      await refreshSummary();
    } catch (e) {
      setContradictionError(String(e));
    }
  }

  async function recordCrisis() {
    setCrisisError(null);
    try {
      await financeApi.createCrisis({
        constant_capital_writedown: writedown,
        pre_crisis_profit_rate: preCrisis,
      });
      await refreshCrises();
      await refreshSummary();
    } catch (e) {
      setCrisisError(String(e));
    }
  }

  const previewPost = restoredProfitRate(preCrisis, writedown);
  const path = concentrationPath(initialC, concSRate, periods);
  const finalC = path[path.length - 1];
  const resultingCapital = absorbedCapitals * absorbedSize;

  return (
    <div className="v3-ch15">
      <section className="v3-ch15-intro">
        <p>
          The law of the falling rate of profit is no smooth curve — it is a{" "}
          <strong>contradiction in motion</strong>. The same rising composition
          that lowers the rate also swells the mass of profit (§I), floods the
          market with commodities that cannot be realised (§II), piles up capital
          that finds no profitable employment (§III), and drives the{" "}
          <strong>concentration and centralisation</strong> of capital into ever
          fewer hands (§IV). The contradiction is resolved only violently, in{" "}
          <strong>crisis</strong>: a mass devaluation of capital that, by writing
          down <code>c</code>, momentarily restores the rate of profit and clears
          the ground for a fresh cycle of accumulation.
        </p>
      </section>

      <section className="v3-ch15-section" aria-label="Contradiction taxonomy">
        <h4>The four internal contradictions (§I–§IV)</h4>
        <p className="v3-ch15-subintro">
          Each is a <strong>coexistent</strong> pairing — not a sequence of
          causes, but one tendency seen from four sides. Record any face below to
          add it to the Part III ledger.
        </p>
        <ul className="v3-ch15-kinds">
          {KINDS.map((k) => (
            <li key={k.kind} className="v3-ch15-kind">
              <div className="v3-ch15-kind-text">
                <span className="v3-ch15-kind-label">
                  <span className="v3-ch15-kind-section">{k.section} </span>
                  {k.label}
                  <span className="v3-ch15-kind-coexist"> · coexistent</span>
                </span>
                <span className="v3-ch15-kind-desc">{k.description}</span>
              </div>
              <button
                type="button"
                className="v3-ch15-btn"
                onClick={() => recordContradiction(k.kind)}
              >
                Record
              </button>
            </li>
          ))}
        </ul>
        {contradictionError && <p className="v3-ch15-error">{contradictionError}</p>}
        {contradictions.length > 0 && (
          <ul className="v3-ch15-recorded">
            {contradictions.map((c) => (
              <li key={c.id}>
                <span className="v3-ch15-recorded-kind">{c.kind}</span>
                <span className="v3-ch15-recorded-flag">
                  {c.is_coexistent ? "coexistent" : "—"}
                </span>
              </li>
            ))}
          </ul>
        )}
      </section>

      <section className="v3-ch15-section" aria-label="Crisis events">
        <h4>Crisis: violent resolution restores the rate</h4>
        <p className="v3-ch15-subintro">
          A crisis writes down a percentage of constant capital. Shrinking the
          denominator of <code>p′ = s/(c+v)</code> lifts the rate — the harsher
          the devaluation, the greater the restoration.
        </p>
        <div className="v3-ch15-form">
          <label>
            Constant-capital writedown (%)
            <input
              type="number"
              min={1}
              max={99}
              value={writedown}
              onChange={(e) => setWritedown(Number(e.target.value))}
            />
          </label>
          <label>
            Pre-crisis rate of profit (bp)
            <input
              type="number"
              min={0}
              value={preCrisis}
              onChange={(e) => setPreCrisis(Number(e.target.value))}
            />
          </label>
          <div className="v3-ch15-preview">
            <span className="v3-ch15-preview-label">Restored p′</span>
            <span className="v3-ch15-preview-figure">
              {pct(preCrisis)} → {pct(previewPost)}
            </span>
          </div>
          <button type="button" className="v3-ch15-btn" onClick={recordCrisis}>
            Record crisis
          </button>
        </div>
        {crisisError && <p className="v3-ch15-error">{crisisError}</p>}

        {crises.length > 0 && (
          <table className="v3-ch15-table">
            <thead>
              <tr>
                <th>Writedown</th>
                <th>Pre-crisis p′</th>
                <th>Post-crisis p′</th>
                <th>Δ</th>
              </tr>
            </thead>
            <tbody>
              {crises.map((c) => (
                <tr key={c.id}>
                  <td>{c.constant_capital_writedown}%</td>
                  <td className="v3-ch15-rate-base">{pct(c.pre_crisis_profit_rate)}</td>
                  <td className="v3-ch15-rate">{pct(c.post_crisis_profit_rate)}</td>
                  <td className="v3-ch15-uplift">
                    +{pct(c.post_crisis_profit_rate - c.pre_crisis_profit_rate)}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </section>

      <section className="v3-ch15-section" aria-label="Concentration calculator">
        <h4>§IV — Concentration: a single capital grows on its own surplus</h4>
        <p className="v3-ch15-subintro">
          Concentration <strong>creates</strong> new capital: each period the
          capital accumulates a fraction <code>s′</code> of itself,{" "}
          <code>C ← C + C·s′</code>, compounding period on period.
        </p>
        <div className="v3-ch15-form">
          <label>
            Initial capital C
            <input
              type="number"
              min={0}
              value={initialC}
              onChange={(e) => setInitialC(Number(e.target.value))}
            />
          </label>
          <label>
            Accumulation rate s′ (bp)
            <input
              type="number"
              min={0}
              value={concSRate}
              onChange={(e) => setConcSRate(Number(e.target.value))}
            />
          </label>
          <label>
            Periods
            <input
              type="number"
              min={1}
              value={periods}
              onChange={(e) => setPeriods(Number(e.target.value))}
            />
          </label>
          <div className="v3-ch15-preview">
            <span className="v3-ch15-preview-label">Final C</span>
            <span className="v3-ch15-preview-figure">{num(finalC)}</span>
          </div>
        </div>
        <table className="v3-ch15-table">
          <thead>
            <tr>
              <th>Period</th>
              <th>Capital</th>
            </tr>
          </thead>
          <tbody>
            {path.map((c, i) => (
              <tr key={i}>
                <td>{i === 0 ? "initial" : i}</td>
                <td className="v3-ch15-rate">{num(c)}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </section>

      <section className="v3-ch15-section" aria-label="Centralisation calculator">
        <h4>§IV — Centralisation: many capitals fused into one</h4>
        <p className="v3-ch15-subintro">
          Centralisation creates <strong>no new</strong> capital — it merely
          redistributes the existing mass. Absorbing many capitals yields one of
          their combined size, <code>AbsorbedCapitals · AbsorbedSize</code>.
        </p>
        <div className="v3-ch15-form">
          <label>
            Capitals absorbed
            <input
              type="number"
              min={0}
              value={absorbedCapitals}
              onChange={(e) => setAbsorbedCapitals(Number(e.target.value))}
            />
          </label>
          <label>
            Average size each
            <input
              type="number"
              min={0}
              value={absorbedSize}
              onChange={(e) => setAbsorbedSize(Number(e.target.value))}
            />
          </label>
          <div className="v3-ch15-preview">
            <span className="v3-ch15-preview-label">Resulting capital</span>
            <span className="v3-ch15-preview-figure">{num(resultingCapital)}</span>
          </div>
        </div>
      </section>

      {summary && (
        <section className="v3-ch15-section v3-ch15-summary" aria-label="Part III summary">
          <h4>Part III Complete — The Law of the Tendential Fall in the Rate of Profit</h4>
          <p className="v3-ch15-subintro">
            Chapters 13–15 close Part III: the law as such, its counteracting
            influences, and the internal contradictions that drive it to crisis.{" "}
            <strong>The real barrier of capitalist production is capital itself.</strong>
          </p>
          <div className="v3-ch15-counts">
            <div className="v3-ch15-count">
              <span className="v3-ch15-count-figure">{summary.trajectory_count}</span>
              <span className="v3-ch15-count-label">Trajectories (Ch. 13)</span>
            </div>
            <div className="v3-ch15-count">
              <span className="v3-ch15-count-figure">{summary.scenario_count}</span>
              <span className="v3-ch15-count-label">Scenarios (Ch. 14)</span>
            </div>
            <div className="v3-ch15-count">
              <span className="v3-ch15-count-figure">{summary.crisis_count}</span>
              <span className="v3-ch15-count-label">Crises (Ch. 15)</span>
            </div>
            <div className="v3-ch15-count">
              <span className="v3-ch15-count-figure">{summary.contradiction_count}</span>
              <span className="v3-ch15-count-label">Contradictions (Ch. 15)</span>
            </div>
          </div>
          <ul className="v3-ch15-latest">
            {summary.latest_trajectory && (
              <li>
                <span className="v3-ch15-latest-label">Latest trajectory</span>
                <span>{summary.latest_trajectory.label}</span>
              </li>
            )}
            {summary.latest_scenario && (
              <li>
                <span className="v3-ch15-latest-label">Latest scenario</span>
                <span>{summary.latest_scenario.label}</span>
              </li>
            )}
            {summary.latest_crisis && (
              <li>
                <span className="v3-ch15-latest-label">Latest crisis</span>
                <span>
                  {summary.latest_crisis.constant_capital_writedown}% writedown ·{" "}
                  {pct(summary.latest_crisis.pre_crisis_profit_rate)} →{" "}
                  {pct(summary.latest_crisis.post_crisis_profit_rate)}
                </span>
              </li>
            )}
            {summary.latest_contradiction && (
              <li>
                <span className="v3-ch15-latest-label">Latest contradiction</span>
                <span>{summary.latest_contradiction.kind}</span>
              </li>
            )}
          </ul>
        </section>
      )}
    </div>
  );
}
