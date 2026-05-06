import { useState } from "react";
import type { FormEvent } from "react";
import { api } from "../api";
import type { SurplusValueSnapshot, SurplusLimitsResponse } from "../types";

interface Ch11Props {
  onSharedChanged: () => void;
}

export function Ch11RateAndMassOfSurplusValue({ onSharedChanged: _unused }: Ch11Props) {
  return (
    <>
      <MassCalculatorPanel />
      <LimitsPanel />
    </>
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
      <div style={{ marginBottom: "0.75rem", display: "flex", gap: "0.5rem", flexWrap: "wrap" }}>
        {FIXTURES.map((f) => (
          <button key={f.label} type="button" onClick={() => loadFixture(f)}>
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
            {result.mass_by_rate > 0 && result.mass_by_workers !== undefined && result.mass_by_workers > 0 && (
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
