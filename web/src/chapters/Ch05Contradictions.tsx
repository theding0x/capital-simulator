import { useState } from "react";
import type { FormEvent } from "react";
import { api } from "../api";
import type { CircuitProof, ExchangeSimulation } from "../types";

interface Ch05Props {
  onSharedChanged: () => void;
}

function penceToGBP(pence: number): string {
  return `£${(pence / 100).toFixed(2)}`;
}

export function Ch05Contradictions({ onSharedChanged: _onSharedChanged }: Ch05Props) {
  return (
    <>
      <CircuitProbePanel />
      <ExchangeSimulationPanel />
    </>
  );
}

function CircuitProbePanel() {
  const [m, setM] = useState(10000);
  const [commodityId, setCommodityId] = useState("cotton");
  const [mPrime, setMPrime] = useState(10000);
  const [result, setResult] = useState<CircuitProof | null>(null);
  const [err, setErr] = useState<string | null>(null);

  async function submit(e: FormEvent) {
    e.preventDefault();
    setErr(null);
    setResult(null);
    try {
      const r = await api.computeCircuit({
        m,
        commodity_id: commodityId || undefined,
        m_prime: mPrime,
      });
      setResult(r);
    } catch (e) {
      setErr(e instanceof Error ? e.message : String(e));
    }
  }

  return (
    <section className="card">
      <h2>Circuit Probe</h2>
      <p className="muted small">
        Test whether M—C—M′ can produce surplus-value through circulation alone.
        Leave commodity blank for M—M′ (usurer&rsquo;s capital).
      </p>
      <form className="form-grid" onSubmit={submit}>
        <label>
          M advanced (pence)
          <input
            type="number"
            value={m}
            onChange={(e) => setM(Number(e.target.value))}
            min={1}
          />
        </label>
        <label>
          Commodity ID (blank = M—M′)
          <input value={commodityId} onChange={(e) => setCommodityId(e.target.value)} />
        </label>
        <label>
          M′ returned (pence)
          <input
            type="number"
            value={mPrime}
            onChange={(e) => setMPrime(Number(e.target.value))}
            min={0}
          />
        </label>
        <button type="submit">Compute</button>
        {err && <span className="error">{err}</span>}
      </form>
      {result && (
        <div className="item-card">
          <div className="item-header">
            <span className="item-name">
              {penceToGBP(result.m)}
              {result.commodity_id ? ` → C (${result.commodity_id.slice(0, 8)}…) →` : " →"}
              {" "}{penceToGBP(result.m_prime)}
            </span>
            <span className={`item-tag${result.origin === "redistribution" ? " negative" : ""}`}>
              {result.origin}
            </span>
          </div>
          <p className="small muted">
            ∆M = {result.surplus_value >= 0 ? "+" : ""}
            {penceToGBP(result.surplus_value)}.{" "}
            {result.origin === "redistribution"
              ? "Value was redistributed between parties, not created."
              : "Exchange at value: no surplus arose from circulation."}
          </p>
        </div>
      )}
    </section>
  );
}

function ExchangeSimulationPanel() {
  const [aValue, setAValue] = useState(5000);
  const [bValue, setBValue] = useState(5000);
  const [result, setResult] = useState<ExchangeSimulation | null>(null);
  const [err, setErr] = useState<string | null>(null);

  async function submit(e: FormEvent) {
    e.preventDefault();
    setErr(null);
    setResult(null);
    try {
      const r = await api.simulateExchange({ a_value: aValue, b_value: bValue });
      setResult(r);
    } catch (e) {
      setErr(e instanceof Error ? e.message : String(e));
    }
  }

  const totalBefore = result ? result.a_before + result.b_before : 0;
  const totalAfter = result ? result.a_after + result.b_after : 0;

  return (
    <section className="card">
      <h2>Exchange Simulation</h2>
      <p className="muted small">
        Prove that bilateral exchange conserves total social value.
      </p>
      <form className="form-grid" onSubmit={submit}>
        <label>
          A holds (pence)
          <input
            type="number"
            value={aValue}
            onChange={(e) => setAValue(Number(e.target.value))}
            min={0}
          />
        </label>
        <label>
          B holds (pence)
          <input
            type="number"
            value={bValue}
            onChange={(e) => setBValue(Number(e.target.value))}
            min={0}
          />
        </label>
        <button type="submit">Simulate</button>
        {err && <span className="error">{err}</span>}
      </form>
      {result && (
        <>
          <table className="data-table">
            <thead>
              <tr>
                <th>Party</th>
                <th>Before</th>
                <th>After</th>
                <th>∆</th>
              </tr>
            </thead>
            <tbody>
              <tr>
                <td>A</td>
                <td>{penceToGBP(result.a_before)}</td>
                <td>{penceToGBP(result.a_after)}</td>
                <td
                  className={
                    result.a_after > result.a_before
                      ? "positive"
                      : result.a_after < result.a_before
                      ? "negative"
                      : ""
                  }
                >
                  {result.a_after >= result.a_before ? "+" : ""}
                  {penceToGBP(result.a_after - result.a_before)}
                </td>
              </tr>
              <tr>
                <td>B</td>
                <td>{penceToGBP(result.b_before)}</td>
                <td>{penceToGBP(result.b_after)}</td>
                <td
                  className={
                    result.b_after > result.b_before
                      ? "positive"
                      : result.b_after < result.b_before
                      ? "negative"
                      : ""
                  }
                >
                  {result.b_after >= result.b_before ? "+" : ""}
                  {penceToGBP(result.b_after - result.b_before)}
                </td>
              </tr>
              <tr>
                <td>
                  <strong>Total</strong>
                </td>
                <td>{penceToGBP(totalBefore)}</td>
                <td>{penceToGBP(totalAfter)}</td>
                <td>{penceToGBP(totalAfter - totalBefore)}</td>
              </tr>
            </tbody>
          </table>
          <p className="small muted">
            Origin: <strong>{result.origin}</strong>.{" "}
            {result.origin === "redistribution"
              ? "A's gain is exactly B's loss. No new value was created."
              : "Values are equal: exchange conserves value perfectly."}
          </p>
        </>
      )}
    </section>
  );
}
