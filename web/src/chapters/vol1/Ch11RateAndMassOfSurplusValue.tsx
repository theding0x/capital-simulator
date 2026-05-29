import { useState } from "react";
import type { FormEvent } from "react";
import { api } from "../../api";
import type { SurplusValueSnapshot, SurplusLimitsResponse } from "../../types";
import "./Ch11RateAndMassOfSurplusValue.css";

interface Ch11Props {
  onSharedChanged: () => void;
}

export function Ch11RateAndMassOfSurplusValue({ onSharedChanged: _unused }: Ch11Props) {
  return (
    <>
      <CompensationInsight />
      <MassCalculatorPanel />
      <CompensationPanel />
      <LimitsPanel />
      <ExhaustionCoda />
    </>
  );
}

function CompensationInsight() {
  return (
    <section className="v1-ch11-insight">
      <h2 className="v1-ch11-insight-h2">S = rate × V (and S = v × rate × n)</h2>
      <div className="v1-ch11-formulas">
        <div className="v1-ch11-formula">
          <span className="v1-ch11-formula-tag">Rate formula</span>
          <span className="v1-ch11-formula-eqn">S = (s / v) × V</span>
        </div>
        <div className="v1-ch11-formula">
          <span className="v1-ch11-formula-tag">Worker formula</span>
          <span className="v1-ch11-formula-eqn">S = v × (s / v) × n</span>
        </div>
      </div>
      <p className="v1-ch11-insight-prose">
        The two formulations agree when <code>V == v × n</code>. The
        compensation law shows that raising the rate of exploitation can
        offset a fall in the variable capital advanced — the mass S stays
        constant as the path changes.
      </p>
    </section>
  );
}

interface CompensationCase {
  id: string;
  label: string;
  rate: number;
  V: number;
}

const COMPENSATION_CASES: CompensationCase[] = [
  { id: "low-rate-large-v", label: "Case A — rate 100%, V = 300", rate: 1.0, V: 300 },
  { id: "high-rate-small-v", label: "Case B — rate 200%, V = 150", rate: 2.0, V: 150 },
];

function CompensationPanel() {
  const total = Math.max(
    ...COMPENSATION_CASES.map((c) => c.V * (1 + c.rate)),
  );
  return (
    <section className="card">
      <h2>Compensation Law</h2>
      <p className="description">
        Two paths to the same surplus. Each bar shows variable capital (v,
        lead) followed by surplus-value (s, gold). Total length is{" "}
        <code>V + S = V × (1 + rate)</code>; the gold area is{" "}
        <code>S = rate × V</code>.
      </p>
      <div className="v1-ch11-compensation">
        {COMPENSATION_CASES.map((c, idx) => {
          const S = c.rate * c.V;
          const span = c.V + S;
          const pct = (n: number) => `${(n / total) * 100}%`;
          return (
            <div key={c.id} className="v1-ch11-case">
              <span className="v1-ch11-case-tag">
                {idx === 0 ? "Path A" : "Path B"}
              </span>
              <h3 className="v1-ch11-case-h3">{c.label}</h3>
              <div
                className="v1-ch11-case-track"
                role="img"
                aria-label={`v=${c.V}, S=${S}`}
              >
                <div
                  className="v1-ch11-case-seg v1-ch11-case-seg--v"
                  style={{ width: pct(c.V) }}
                  title={`V = ${c.V}`}
                >
                  V
                </div>
                <div
                  className="v1-ch11-case-seg v1-ch11-case-seg--s"
                  style={{ width: pct(S) }}
                  title={`S = ${S}`}
                >
                  S
                </div>
              </div>
              <span className="v1-ch11-case-meta">
                S = rate × V = <strong>{S}</strong>{" "}
                <span className="v1-ch11-case-meta-muted">
                  · V + S = {span}
                </span>
              </span>
            </div>
          );
        })}
      </div>
      <div className="v1-ch11-invariant">
        <span className="v1-ch11-invariant-label">Invariant S</span>
        <span className="v1-ch11-invariant-value">
          {COMPENSATION_CASES[0].rate * COMPENSATION_CASES[0].V}
        </span>
      </div>
    </section>
  );
}

function ExhaustionCoda() {
  return (
    <aside className="v1-ch11-coda">
      <p className="v1-ch11-coda-quote">
        The compensation between rate and variable capital meets a hard
        physiological boundary: the workman's exhaustion fixes a floor under
        v and a ceiling over s. The mass of surplus-value cannot be sustained
        indefinitely by raising the rate alone.
        <span className="v1-ch11-coda-cite">
          — closing note, Capital Vol. I, Ch. 11 §1
        </span>
      </p>
    </aside>
  );
}

type Fixture = {
  label: string;
  surplusLabour: number;
  necessaryLabour: number;
  variableCapital?: number;
  labourPowerValue?: number;
  workerCount?: number;
};

const FIXTURES: Fixture[] = [
  {
    label: "§1 Rate 100%, V=3s (1 worker)",
    surplusLabour: 6,
    necessaryLabour: 6,
    variableCapital: 3,
  },
  {
    label: "§1 100 workers, v=3s, rate 100%",
    surplusLabour: 6,
    necessaryLabour: 6,
    labourPowerValue: 3,
    workerCount: 100,
  },
  {
    label: "§1 Compensation: rate 200%, V=150s",
    surplusLabour: 12,
    necessaryLabour: 6,
    variableCapital: 150,
  },
  {
    label: "§1 V=1500s, rate 100%",
    surplusLabour: 6,
    necessaryLabour: 6,
    variableCapital: 1500,
  },
  {
    label: "§1 V=300s, rate 200%",
    surplusLabour: 12,
    necessaryLabour: 6,
    variableCapital: 300,
  },
];

function MassCalculatorPanel() {
  const [sl, setSl] = useState(6);
  const [nl, setNl] = useState(6);
  const [vc, setVc] = useState<number | "">("");
  const [lpv, setLpv] = useState<number | "">("");
  const [wc, setWc] = useState<number | "">("");
  const [result, setResult] = useState<SurplusValueSnapshot | null>(null);
  const [err, setErr] = useState<string | null>(null);

  function loadFixture(f: Fixture) {
    setSl(f.surplusLabour);
    setNl(f.necessaryLabour);
    setVc(f.variableCapital ?? "");
    setLpv(f.labourPowerValue ?? "");
    setWc(f.workerCount ?? "");
    setResult(null);
    setErr(null);
  }

  async function submit(e: FormEvent) {
    e.preventDefault();
    setErr(null);
    try {
      const input = {
        surplus_labour: sl,
        necessary_labour: nl,
        ...(vc !== "" ? { variable_capital: Number(vc) } : {}),
        ...(lpv !== "" ? { labour_power_value: Number(lpv) } : {}),
        ...(wc !== "" ? { worker_count: Number(wc) } : {}),
      };
      const r = await api.computeSurplusMass(input);
      setResult(r);
    } catch (e) {
      setErr(e instanceof Error ? e.message : String(e));
    }
  }

  const rate = nl > 0 ? ((sl / nl) * 100).toFixed(1) : "—";

  return (
    <section className="card">
      <h2>Mass of Surplus-Value</h2>
      <p className="description">
        Computes S via two formulas: S = (s/v) × V [rate formula] and S = v × (s/v) × n
        [worker formula]. When V = v × n both must agree — the compensation law shows that
        raising the rate can offset a fall in the number of workers.
      </p>
      <div className="v1-ch11-presets">
        <span className="v1-ch11-presets-label">Marx fixture</span>
        {FIXTURES.map((f) => (
          <button
            key={f.label}
            type="button"
            className="v1-ch11-preset-button"
            onClick={() => loadFixture(f)}
          >
            {f.label}
          </button>
        ))}
      </div>
      <form className="form-grid" onSubmit={submit}>
        <label>
          <span>Surplus Labour a′ (min)</span>
          <input type="number" min={0} value={sl} onChange={(e) => setSl(Number(e.target.value))} />
        </label>
        <label>
          <span>Necessary Labour a (min)</span>
          <input type="number" min={1} value={nl} onChange={(e) => setNl(Number(e.target.value))} />
        </label>
        <label>
          <span>Rate s/v</span>
          <input type="text" readOnly value={`${rate}%`} />
        </label>
        <label>
          <span>Variable Capital V (optional)</span>
          <input
            type="number"
            min={0}
            value={vc}
            placeholder="e.g. 300"
            onChange={(e) => setVc(e.target.value === "" ? "" : Number(e.target.value))}
          />
        </label>
        <label>
          <span>Labour-Power Value v per worker (optional)</span>
          <input
            type="number"
            min={0}
            value={lpv}
            placeholder="e.g. 3"
            onChange={(e) => setLpv(e.target.value === "" ? "" : Number(e.target.value))}
          />
        </label>
        <label>
          <span>Worker Count n (optional)</span>
          <input
            type="number"
            min={1}
            value={wc}
            placeholder="e.g. 100"
            onChange={(e) => setWc(e.target.value === "" ? "" : Number(e.target.value))}
          />
        </label>
        <div className="form-actions span2">
          <button type="submit" className="primary">Compute</button>
          {err && <span className="error">{err}</span>}
        </div>
      </form>
      {result && (
        <table className="data-table">
          <tbody>
            <tr>
              <td>Rate (a′/a)</td>
              <td><strong>{(result.rate.surplus_labour / result.rate.necessary_labour * 100).toFixed(1)}%</strong></td>
            </tr>
            <tr><td>Variable Capital V</td><td>{result.variable_capital}</td></tr>
            <tr><td>Worker Count n</td><td>{result.worker_count ?? "—"}</td></tr>
            <tr>
              <td>Mass by Rate formula  S = (s/v)×V</td>
              <td>{result.mass_by_rate}</td>
            </tr>
            <tr>
              <td>Mass by Workers formula  S = v×(s/v)×n</td>
              <td>{result.mass_by_workers ?? "—"}</td>
            </tr>
            <tr>
              <td><strong>Total Surplus-Value S</strong></td>
              <td><strong>{result.mass}</strong></td>
            </tr>
            {result.mass_by_workers !== undefined && (
              <tr>
                <td>Formulas agree?</td>
                <td>
                  {result.mass_by_rate === result.mass_by_workers
                    ? "Yes"
                    : <span className="error">No — inputs inconsistent (V ≠ v×n)</span>}
                </td>
              </tr>
            )}
          </tbody>
        </table>
      )}
    </section>
  );
}

function LimitsPanel() {
  const [lpv, setLpv] = useState<number | "">("");
  const [wc, setWc] = useState<number | "">("");
  const [result, setResult] = useState<SurplusLimitsResponse | null>(null);
  const [err, setErr] = useState<string | null>(null);

  async function submit(e: FormEvent) {
    e.preventDefault();
    setErr(null);
    try {
      const r = await api.getSurplusLimits(
        lpv !== "" ? Number(lpv) : undefined,
        wc !== "" ? Number(wc) : undefined,
      );
      setResult(r);
    } catch (e) {
      setErr(e instanceof Error ? e.message : String(e));
    }
  }

  return (
    <section className="card">
      <h2>Surplus-Value Limits</h2>
      <p className="description">
        Returns the absolute physical limit of the working day (24 h = 1440 min) and,
        optionally, the minimum variable capital required to employ n workers at daily
        reproduction cost v.
      </p>
      <form className="form-grid" onSubmit={submit}>
        <label>
          <span>Labour-Power Value v (optional)</span>
          <input
            type="number"
            min={1}
            value={lpv}
            placeholder="e.g. 3"
            onChange={(e) => setLpv(e.target.value === "" ? "" : Number(e.target.value))}
          />
        </label>
        <label>
          <span>Worker Count n (optional)</span>
          <input
            type="number"
            min={1}
            value={wc}
            placeholder="e.g. 2"
            onChange={(e) => setWc(e.target.value === "" ? "" : Number(e.target.value))}
          />
        </label>
        <div className="form-actions span2">
          <button type="submit" className="primary">Get Limits</button>
          {err && <span className="error">{err}</span>}
        </div>
      </form>
      {result && (
        <table className="data-table">
          <tbody>
            <tr>
              <td>Absolute Workday Limit</td>
              <td>{result.absolute_workday_limit} min ({result.absolute_workday_limit / 60} h)</td>
            </tr>
            {result.minimum_capital !== undefined && (
              <>
                <tr>
                  <td>Labour-Power Value v</td>
                  <td>{result.labour_power_value}</td>
                </tr>
                <tr>
                  <td>Worker Count n</td>
                  <td>{result.worker_count}</td>
                </tr>
                <tr>
                  <td><strong>Minimum Capital V_min = v × n</strong></td>
                  <td><strong>{result.minimum_capital}</strong></td>
                </tr>
              </>
            )}
          </tbody>
        </table>
      )}
    </section>
  );
}
