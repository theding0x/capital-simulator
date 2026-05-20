import { useEffect, useState } from "react";
import type { FormEvent } from "react";
import { api } from "../../api";
import type {
  LabourWorker,
  LabourCapitalist,
  RunLabourProcessResult,
  RawMaterial,
  Instrument,
} from "../../types";

interface Ch07Props {
  onSharedChanged: () => void;
}

function minutesToHours(m: number): string {
  const h = Math.floor(m / 60);
  const min = m % 60;
  return min === 0 ? `${h}h` : `${h}h ${min}m`;
}

export function Ch07LabourProcess({ onSharedChanged: _onSharedChanged }: Ch07Props) {
  const [workers, setWorkers] = useState<LabourWorker[]>([]);
  const [capitalists, setCapitalists] = useState<LabourCapitalist[]>([]);
  const [result, setResult] = useState<RunLabourProcessResult | null>(null);
  const [loadErr, setLoadErr] = useState<string | null>(null);

  async function refresh() {
    try {
      const [ws, cs] = await Promise.all([
        api.listLabourWorkers(),
        api.listLabourCapitalists(),
      ]);
      setWorkers(ws);
      setCapitalists(cs);
    } catch (e) {
      setLoadErr(e instanceof Error ? e.message : String(e));
    }
  }

  useEffect(() => { refresh(); }, []);

  return (
    <>
      {loadErr && <p className="error">{loadErr}</p>}
      <RunLabourProcessPanel
        workers={workers}
        capitalists={capitalists}
        onResult={setResult}
      />
      {result && <ValorizationResultPanel result={result} />}
    </>
  );
}

interface RunPanelProps {
  workers: LabourWorker[];
  capitalists: LabourCapitalist[];
  onResult: (r: RunLabourProcessResult) => void;
}

function RunLabourProcessPanel({ workers, capitalists, onResult }: RunPanelProps) {
  const [workerID, setWorkerID] = useState("");
  const [capitalistID, setCapitalistID] = useState("");
  const [duration, setDuration] = useState(720);
  const [productKind, setProductKind] = useState("yarn");
  const [productQty, setProductQty] = useState(10);
  const [rawMaterials, setRawMaterials] = useState<RawMaterial[]>([
    { commodity_id: "cotton", quantity: 10, snlt_per_unit: 120 },
  ]);
  const [instruments, setInstruments] = useState<Instrument[]>([
    { commodity_id: "spindle", wear_per_run: 240 },
  ]);
  const [err, setErr] = useState<string | null>(null);

  async function submit(e: FormEvent) {
    e.preventDefault();
    setErr(null);
    try {
      const r = await api.runLabourProcess({
        worker_id: workerID,
        capitalist_id: capitalistID,
        means_of_production: { raw_materials: rawMaterials, instruments },
        duration_minutes: duration,
        product_kind: productKind,
        product_quantity: productQty,
      });
      onResult(r);
    } catch (e) {
      setErr(e instanceof Error ? e.message : String(e));
    }
  }

  function updateRawMaterial(i: number, field: keyof RawMaterial, value: string | number) {
    setRawMaterials((prev) =>
      prev.map((rm, idx) => (idx === i ? { ...rm, [field]: value } : rm))
    );
  }

  function updateInstrument(i: number, field: keyof Instrument, value: string | number) {
    setInstruments((prev) =>
      prev.map((inst, idx) => (idx === i ? { ...inst, [field]: value } : inst))
    );
  }

  return (
    <section className="card">
      <h2>Run Labour Process</h2>
      <p className="description">
        The capitalist sets the worker to work on raw materials with instruments of labour.
        The result is a new product whose value exceeds the value of labour-power advanced.
      </p>
      <form className="form-grid" onSubmit={submit}>
        <label>
          <span>Worker</span>
          <select value={workerID} onChange={(e) => setWorkerID(e.target.value)} required>
            <option value="">Select a worker…</option>
            {workers.map((w) => (
              <option key={w.id} value={w.id}>
                {w.id.slice(0, 8)}… (capacity {minutesToHours(w.labour_power.capacity_minutes_per_day)},
                repro cost {minutesToHours(w.labour_power_value_minutes)})
              </option>
            ))}
          </select>
        </label>
        <label>
          <span>Capitalist</span>
          <select value={capitalistID} onChange={(e) => setCapitalistID(e.target.value)} required>
            <option value="">Select a capitalist…</option>
            {capitalists.map((c) => (
              <option key={c.id} value={c.id}>
                {c.id.slice(0, 8)}…
              </option>
            ))}
          </select>
        </label>
        <label>
          <span>Working Day Duration (minutes)</span>
          <input
            type="number"
            min={1}
            value={duration}
            onChange={(e) => setDuration(Number(e.target.value))}
          />
        </label>

        <div className="span2">
          <fieldset>
            <legend>Raw Materials</legend>
            {rawMaterials.map((rm, i) => (
              <div key={i} className="item-row">
                <input
                  placeholder="commodity_id"
                  value={rm.commodity_id}
                  onChange={(e) => updateRawMaterial(i, "commodity_id", e.target.value)}
                />
                <input
                  type="number"
                  placeholder="qty"
                  value={rm.quantity}
                  onChange={(e) => updateRawMaterial(i, "quantity", Number(e.target.value))}
                />
                <input
                  type="number"
                  placeholder="snlt_per_unit (min)"
                  value={rm.snlt_per_unit}
                  onChange={(e) => updateRawMaterial(i, "snlt_per_unit", Number(e.target.value))}
                />
                <button
                  type="button"
                  onClick={() => setRawMaterials((prev) => prev.filter((_, idx) => idx !== i))}
                >
                  ×
                </button>
              </div>
            ))}
            <button
              type="button"
              onClick={() =>
                setRawMaterials((prev) => [
                  ...prev,
                  { commodity_id: "", quantity: 1, snlt_per_unit: 0 },
                ])
              }
            >
              + Add Raw Material
            </button>
          </fieldset>
        </div>

        <div className="span2">
          <fieldset>
            <legend>Instruments of Labour</legend>
            {instruments.map((inst, i) => (
              <div key={i} className="item-row">
                <input
                  placeholder="commodity_id"
                  value={inst.commodity_id}
                  onChange={(e) => updateInstrument(i, "commodity_id", e.target.value)}
                />
                <input
                  type="number"
                  placeholder="wear_per_run (min)"
                  value={inst.wear_per_run}
                  onChange={(e) => updateInstrument(i, "wear_per_run", Number(e.target.value))}
                />
                <button
                  type="button"
                  onClick={() => setInstruments((prev) => prev.filter((_, idx) => idx !== i))}
                >
                  ×
                </button>
              </div>
            ))}
            <button
              type="button"
              onClick={() =>
                setInstruments((prev) => [...prev, { commodity_id: "", wear_per_run: 0 }])
              }
            >
              + Add Instrument
            </button>
          </fieldset>
        </div>

        <label>
          <span>Product Kind</span>
          <input
            value={productKind}
            onChange={(e) => setProductKind(e.target.value)}
            placeholder="yarn"
          />
        </label>
        <label>
          <span>Product Quantity</span>
          <input
            type="number"
            min={1}
            value={productQty}
            onChange={(e) => setProductQty(Number(e.target.value))}
          />
        </label>

        <div className="form-actions span2">
          <button type="submit" className="primary">Run Labour Process</button>
          {err && <span className="error">{err}</span>}
        </div>
      </form>
    </section>
  );
}

function ValorizationResultPanel({ result }: { result: RunLabourProcessResult }) {
  const { labour_process: lp, product, valorization: v } = result;
  const surplusRate =
    v.necessary_labour > 0
      ? ((v.surplus_value / v.necessary_labour) * 100).toFixed(1)
      : "∞";

  return (
    <section className="card">
      <h2>Valorization Result</h2>
      <p className="small muted">Process ID: <span className="monospace">{lp.id}</span></p>
      <table className="data-table">
        <tbody>
          <tr>
            <td>Working Day</td>
            <td>{minutesToHours(lp.duration)}</td>
          </tr>
          <tr>
            <td>Necessary Labour</td>
            <td>{minutesToHours(v.necessary_labour)}</td>
          </tr>
          <tr>
            <td>Surplus Labour</td>
            <td>{minutesToHours(v.surplus_labour)}</td>
          </tr>
          <tr>
            <td>Transferred Value (constant capital)</td>
            <td>{minutesToHours(v.product_value - lp.duration)}</td>
          </tr>
          <tr>
            <td>Value Added (living labour)</td>
            <td>{minutesToHours(lp.duration)}</td>
          </tr>
          <tr>
            <td>Total Product Value</td>
            <td><strong>{minutesToHours(v.product_value)}</strong></td>
          </tr>
          <tr>
            <td>Surplus Value</td>
            <td><strong>{minutesToHours(v.surplus_value)}</strong></td>
          </tr>
          <tr>
            <td>Rate of Surplus Value (s/v)</td>
            <td>{surplusRate}%</td>
          </tr>
        </tbody>
      </table>
      <p className="small">
        <strong>Product:</strong> {lp.product_quantity} × {product.commodity_kind} —
        total value {minutesToHours(product.total_value)}
      </p>
    </section>
  );
}
