import { useState } from "react";
import type { FormEvent } from "react";
import { api } from "../../api";
import type {
  WorkingDayResponse,
  ValidateWorkingDayResponse,
  RelaySchedule,
} from "../../types";

interface Ch10Props {
  onSharedChanged: () => void;
}

function minutesToHours(m: number): string {
  const h = Math.floor(m / 60);
  const min = m % 60;
  return min === 0 ? `${h}h` : `${h}h ${min}m`;
}

export function Ch10WorkingDay({ onSharedChanged: _onSharedChanged }: Ch10Props) {
  return (
    <>
      <ValidatePanel />
      <CreatePanel />
      <RelaySchedulePanel />
    </>
  );
}

function ValidatePanel() {
  const [nl, setNl] = useState(360);
  const [sl, setSl] = useState(360);
  const [limit, setLimit] = useState<number | "">("");
  const [result, setResult] = useState<ValidateWorkingDayResponse | null>(null);
  const [err, setErr] = useState<string | null>(null);

  function loadFixture(necessary: number, surplus: number, lim?: number) {
    setNl(necessary);
    setSl(surplus);
    setLimit(lim ?? "");
    setResult(null);
  }

  async function submit(e: FormEvent) {
    e.preventDefault();
    setErr(null);
    try {
      const input = {
        necessary_labour_minutes: nl,
        surplus_labour_minutes: sl,
        ...(limit !== "" ? { statutory_limit_minutes: Number(limit) } : {}),
      };
      const r = await api.validateWorkingDay(input);
      setResult(r);
    } catch (e) {
      setErr(e instanceof Error ? e.message : String(e));
    }
  }

  return (
    <section className="card">
      <h2>Working-Day Validator</h2>
      <p className="description">
        Validate a working day against the physical maximum (24 h) and an optional statutory
        limit. Computes rate of surplus-value from the A–B / B–C segments.
      </p>
      <div style={{ marginBottom: "0.75rem", display: "flex", gap: "0.5rem", flexWrap: "wrap" }}>
        <button type="button" onClick={() => loadFixture(360, 60)}>
          §1 WD I (6h+1h)
        </button>
        <button type="button" onClick={() => loadFixture(360, 180)}>
          §1 WD II (6h+3h)
        </button>
        <button type="button" onClick={() => loadFixture(360, 360)}>
          §1 WD III (6h+6h)
        </button>
        <button type="button" onClick={() => loadFixture(315, 315, 630)}>
          Factory Act 1850 (10.5 h)
        </button>
        <button type="button" onClick={() => loadFixture(84 * 480, 56 * 480)}>
          Wallachian Corvée
        </button>
      </div>
      <form className="form-grid" onSubmit={submit}>
        <label>
          <span>Necessary Labour A–B (min)</span>
          <input type="number" min={1} value={nl} onChange={(e) => setNl(Number(e.target.value))} />
        </label>
        <label>
          <span>Surplus Labour B–C (min)</span>
          <input type="number" min={0} value={sl} onChange={(e) => setSl(Number(e.target.value))} />
        </label>
        <label>
          <span>Statutory Limit (min, optional)</span>
          <input
            type="number"
            min={1}
            value={limit}
            placeholder="e.g. 630"
            onChange={(e) => setLimit(e.target.value === "" ? "" : Number(e.target.value))}
          />
        </label>
        <div className="form-actions span2">
          <button type="submit" className="primary">Validate</button>
          {err && <span className="error">{err}</span>}
        </div>
      </form>
      {result && (
        <table className="data-table">
          <tbody>
            <tr><td>Total</td><td>{minutesToHours(result.total_minutes)}</td></tr>
            <tr>
              <td>Rate of Surplus-Value (s/v)</td>
              <td><strong>{(result.rate_of_surplus_value * 100).toFixed(2)}%</strong></td>
            </tr>
            <tr>
              <td>Valid</td>
              <td>{result.valid ? "Yes" : <span className="error">No — {result.error}</span>}</td>
            </tr>
          </tbody>
        </table>
      )}
    </section>
  );
}

function CreatePanel() {
  const [nl, setNl] = useState(360);
  const [sl, setSl] = useState(360);
  const [limit, setLimit] = useState<number | "">("");
  const [result, setResult] = useState<WorkingDayResponse | null>(null);
  const [err, setErr] = useState<string | null>(null);

  async function submit(e: FormEvent) {
    e.preventDefault();
    setErr(null);
    try {
      const input = {
        necessary_labour_minutes: nl,
        surplus_labour_minutes: sl,
        ...(limit !== "" ? { statutory_limit_minutes: Number(limit) } : {}),
      };
      const r = await api.createWorkingDay(input);
      setResult(r);
    } catch (e) {
      setErr(e instanceof Error ? e.message : String(e));
    }
  }

  return (
    <section className="card">
      <h2>Create Working-Day Record</h2>
      <p className="description">
        Persist a working-day record. Returns the stored ID alongside computed metrics.
      </p>
      <form className="form-grid" onSubmit={submit}>
        <label>
          <span>Necessary Labour A–B (min)</span>
          <input type="number" min={1} value={nl} onChange={(e) => setNl(Number(e.target.value))} />
        </label>
        <label>
          <span>Surplus Labour B–C (min)</span>
          <input type="number" min={0} value={sl} onChange={(e) => setSl(Number(e.target.value))} />
        </label>
        <label>
          <span>Statutory Limit (min, optional)</span>
          <input
            type="number"
            min={1}
            value={limit}
            placeholder="e.g. 630"
            onChange={(e) => setLimit(e.target.value === "" ? "" : Number(e.target.value))}
          />
        </label>
        <div className="form-actions span2">
          <button type="submit" className="primary">Create</button>
          {err && <span className="error">{err}</span>}
        </div>
      </form>
      {result && (
        <div>
          <p className="small muted">Saved: <span className="monospace">{result.working_day.id}</span></p>
          <table className="data-table">
            <tbody>
              <tr><td>Total</td><td>{minutesToHours(result.total_minutes)}</td></tr>
              <tr>
                <td>Rate of Surplus-Value (s/v)</td>
                <td><strong>{(result.rate_of_surplus_value * 100).toFixed(2)}%</strong></td>
              </tr>
            </tbody>
          </table>
        </div>
      )}
    </section>
  );
}

function RelaySchedulePanel() {
  const [nl0, setNl0] = useState(360);
  const [sl0, setSl0] = useState(360);
  const [nl1, setNl1] = useState(360);
  const [sl1, setSl1] = useState(0);
  const [result, setResult] = useState<RelaySchedule | null>(null);
  const [err, setErr] = useState<string | null>(null);

  async function submit(e: FormEvent) {
    e.preventDefault();
    setErr(null);
    try {
      const r = await api.createRelaySchedule({
        sets: [
          { shift_kind: "day", necessary_labour_minutes: nl0, surplus_labour_minutes: sl0, worker_ids: [] },
          { shift_kind: "night", necessary_labour_minutes: nl1, surplus_labour_minutes: sl1, worker_ids: [] },
        ],
      });
      setResult(r);
    } catch (e) {
      setErr(e instanceof Error ? e.message : String(e));
    }
  }

  const combinedMin = nl0 + sl0 + nl1 + sl1;

  return (
    <section className="card">
      <h2>Relay Schedule</h2>
      <p className="description">
        A relay system alternates two worker sets to keep machinery running. Combined shifts must
        not exceed 24 h (1440 min).
      </p>
      <form className="form-grid" onSubmit={submit}>
        <label>
          <span>Day shift — Necessary (min)</span>
          <input type="number" min={1} value={nl0} onChange={(e) => setNl0(Number(e.target.value))} />
        </label>
        <label>
          <span>Day shift — Surplus (min)</span>
          <input type="number" min={0} value={sl0} onChange={(e) => setSl0(Number(e.target.value))} />
        </label>
        <label>
          <span>Night shift — Necessary (min)</span>
          <input type="number" min={1} value={nl1} onChange={(e) => setNl1(Number(e.target.value))} />
        </label>
        <label>
          <span>Night shift — Surplus (min)</span>
          <input type="number" min={0} value={sl1} onChange={(e) => setSl1(Number(e.target.value))} />
        </label>
        <div className="form-actions span2">
          <span className="small muted">
            Combined: {minutesToHours(combinedMin)} / 24h max
            {combinedMin > 1440 && <span className="error"> — exceeds physical max</span>}
          </span>
          <button type="submit" className="primary">Create Relay Schedule</button>
          {err && <span className="error">{err}</span>}
        </div>
      </form>
      {result && (
        <div>
          <p className="small muted">Saved: <span className="monospace">{result.id}</span></p>
          <table className="data-table">
            <thead>
              <tr><th>Shift</th><th>Necessary</th><th>Surplus</th><th>Total</th></tr>
            </thead>
            <tbody>
              {result.sets.map((s, i) => (
                <tr key={i}>
                  <td>{s.shift_kind}</td>
                  <td>{minutesToHours(s.working_day.necessary_labour_minutes)}</td>
                  <td>{minutesToHours(s.working_day.surplus_labour_minutes)}</td>
                  <td>
                    {minutesToHours(
                      s.working_day.necessary_labour_minutes + s.working_day.surplus_labour_minutes
                    )}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </section>
  );
}
