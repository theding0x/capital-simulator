import { useEffect, useState } from "react";
import type { FormEvent } from "react";
import { api } from "../api";
import type {
  RateOfSurplusValueResult,
  ProductionAccountResult,
} from "../types";

interface Ch09Props {
  onSharedChanged: () => void;
}

function minutesToHours(m: number): string {
  const h = Math.floor(m / 60);
  const min = m % 60;
  return min === 0 ? `${h}h` : `${h}h ${min}m`;
}

export function Ch09RateOfSurplusValue({ onSharedChanged: _onSharedChanged }: Ch09Props) {
  const [accounts, setAccounts] = useState<ProductionAccountResult[]>([]);

  async function refreshAccounts() {
    try {
      const list = await api.listProductionAccounts();
      setAccounts(list);
    } catch (_) {
      // non-fatal if listing fails
    }
  }

  useEffect(() => { refreshAccounts(); }, []);

  return (
    <>
      <RateProbePanel />
      <RecordAccountPanel onCreated={refreshAccounts} />
      {accounts.length > 0 && <AccountListPanel accounts={accounts} />}
    </>
  );
}

function RateProbePanel() {
  const [surplus, setSurplus] = useState(90);
  const [variable, setVariable] = useState(90);
  const [result, setResult] = useState<RateOfSurplusValueResult | null>(null);
  const [err, setErr] = useState<string | null>(null);

  async function submit(e: FormEvent) {
    e.preventDefault();
    setErr(null);
    try {
      const r = await api.computeRateOfSurplusValue({ surplus, variable });
      setResult(r);
    } catch (e) {
      setErr(e instanceof Error ? e.message : String(e));
    }
  }

  function loadFixture(s: number, v: number) {
    setSurplus(s);
    setVariable(v);
    setResult(null);
  }

  return (
    <section className="card">
      <h2>Rate of Surplus-Value Probe</h2>
      <p className="description">
        Compute s/v — the rate of surplus-value (degree of exploitation) from surplus and
        variable capital magnitudes in labour-minutes.
      </p>
      <div style={{ marginBottom: "0.75rem", display: "flex", gap: "0.5rem", flexWrap: "wrap" }}>
        <button type="button" onClick={() => loadFixture(80, 52)}>
          1871 Spinning Mill (s=80, v=52)
        </button>
        <button type="button" onClick={() => loadFixture(211, 210)}>
          Jacob's Wheat 1815 (s=211, v=210)
        </button>
        <button type="button" onClick={() => loadFixture(360, 360)}>
          Cotton Spinner (s=360, v=360)
        </button>
      </div>
      <form className="form-grid" onSubmit={submit}>
        <label>
          <span>Surplus Value s (min)</span>
          <input
            type="number"
            min={0}
            value={surplus}
            onChange={(e) => setSurplus(Number(e.target.value))}
          />
        </label>
        <label>
          <span>Variable Capital v (min)</span>
          <input
            type="number"
            min={1}
            value={variable}
            onChange={(e) => setVariable(Number(e.target.value))}
          />
        </label>
        <div className="form-actions span2">
          <button type="submit" className="primary">Compute Rate</button>
          {err && <span className="error">{err}</span>}
        </div>
      </form>

      {result && (
        <table className="data-table">
          <tbody>
            <tr>
              <td>Surplus Value (s)</td>
              <td>{minutesToHours(result.surplus)}</td>
            </tr>
            <tr>
              <td>Variable Capital (v)</td>
              <td>{minutesToHours(result.variable)}</td>
            </tr>
            <tr>
              <td>Rate of Surplus-Value (s/v)</td>
              <td><strong>{(result.rate * 100).toFixed(2)}%</strong></td>
            </tr>
            <tr>
              <td>Value Product (v + s)</td>
              <td>{minutesToHours(result.value_product)}</td>
            </tr>
            <tr>
              <td>Surplus-Produce Fraction (s/(v+s))</td>
              <td>{(result.surplus_produce_fraction * 100).toFixed(1)}%</td>
            </tr>
          </tbody>
        </table>
      )}
    </section>
  );
}

interface RecordAccountPanelProps {
  onCreated: () => void;
}

function RecordAccountPanel({ onCreated }: RecordAccountPanelProps) {
  const [constant, setConstant] = useState(410);
  const [variable, setVariable] = useState(90);
  const [surplus, setSurplus] = useState(90);
  const [result, setResult] = useState<ProductionAccountResult | null>(null);
  const [err, setErr] = useState<string | null>(null);

  async function submit(e: FormEvent) {
    e.preventDefault();
    setErr(null);
    try {
      const r = await api.createProductionAccount({ constant, variable, surplus });
      setResult(r);
      onCreated();
    } catch (e) {
      setErr(e instanceof Error ? e.message : String(e));
    }
  }

  return (
    <section className="card">
      <h2>Record Production Account</h2>
      <p className="description">
        Persist a production account recording c, v, s for a single production run.
        Rate of surplus-value is computed and stored automatically.
      </p>
      <form className="form-grid" onSubmit={submit}>
        <label>
          <span>Constant Capital c (min)</span>
          <input
            type="number"
            min={0}
            value={constant}
            onChange={(e) => setConstant(Number(e.target.value))}
          />
        </label>
        <label>
          <span>Variable Capital v (min)</span>
          <input
            type="number"
            min={1}
            value={variable}
            onChange={(e) => setVariable(Number(e.target.value))}
          />
        </label>
        <label>
          <span>Surplus Value s (min)</span>
          <input
            type="number"
            min={0}
            value={surplus}
            onChange={(e) => setSurplus(Number(e.target.value))}
          />
        </label>
        <div className="form-actions span2">
          <button type="submit" className="primary">Record Account</button>
          {err && <span className="error">{err}</span>}
        </div>
      </form>

      {result && (
        <div>
          <p className="small muted">Saved: <span className="monospace">{result.id}</span></p>
          <table className="data-table">
            <tbody>
              <tr><td>c</td><td>{minutesToHours(result.constant)}</td></tr>
              <tr><td>v</td><td>{minutesToHours(result.variable)}</td></tr>
              <tr><td>s</td><td>{minutesToHours(result.surplus)}</td></tr>
              <tr>
                <td>Rate (s/v)</td>
                <td><strong>{(result.rate_of_surplus_value * 100).toFixed(2)}%</strong></td>
              </tr>
              <tr>
                <td>Value Product (v + s)</td>
                <td>{minutesToHours(result.value_product)}</td>
              </tr>
              <tr>
                <td>Expanded Capital (c + v + s)</td>
                <td>{minutesToHours(result.expanded_capital)}</td>
              </tr>
            </tbody>
          </table>
        </div>
      )}
    </section>
  );
}

function AccountListPanel({ accounts }: { accounts: ProductionAccountResult[] }) {
  return (
    <section className="card">
      <h2>Recorded Production Accounts</h2>
      <table className="data-table">
        <thead>
          <tr>
            <th>ID</th>
            <th>c</th>
            <th>v</th>
            <th>s</th>
            <th>Rate (s/v)</th>
            <th>Expanded Capital</th>
          </tr>
        </thead>
        <tbody>
          {accounts.map((a) => (
            <tr key={a.id}>
              <td className="monospace">{a.id.slice(0, 8)}…</td>
              <td>{minutesToHours(a.constant)}</td>
              <td>{minutesToHours(a.variable)}</td>
              <td>{minutesToHours(a.surplus)}</td>
              <td>{(a.rate_of_surplus_value * 100).toFixed(2)}%</td>
              <td>{minutesToHours(a.expanded_capital)}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </section>
  );
}
