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
import "./Ch07LabourProcess.css";

interface Ch07Props {
  onSharedChanged: () => void;
}

function minutesToHours(m: number): string {
  const h = Math.floor(m / 60);
  const min = m % 60;
  return min === 0 ? `${h}h` : `${h}h ${min}m`;
}

interface FixturePreset {
  id: string;
  label: string;
  duration: number;
  productKind: string;
  productQty: number;
  rawMaterials: RawMaterial[];
  instruments: Instrument[];
}

const FIXTURES: FixturePreset[] = [
  {
    id: "yarn-1871",
    label: "Yarn / cotton (Ch.7 §2)",
    duration: 720,
    productKind: "yarn",
    productQty: 10,
    rawMaterials: [{ commodity_id: "cotton", quantity: 10, snlt_per_unit: 120 }],
    instruments: [{ commodity_id: "spindle", wear_per_run: 240 }],
  },
  {
    id: "yarn-1871-twelve",
    label: "Yarn / 12-hour day (s/v = 100%)",
    duration: 720,
    productKind: "yarn",
    productQty: 20,
    rawMaterials: [{ commodity_id: "cotton", quantity: 20, snlt_per_unit: 60 }],
    instruments: [{ commodity_id: "spindle", wear_per_run: 240 }],
  },
];

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
      <TwoFacesInsight />
      {loadErr && <p className="error">{loadErr}</p>}
      <RunLabourProcessPanel
        workers={workers}
        capitalists={capitalists}
        onResult={setResult}
      />
      {result && <ValorizationResultPanel result={result} />}
      <Coda />
    </>
  );
}

function TwoFacesInsight() {
  return (
    <section className="v1-ch07-insight">
      <h2 className="v1-ch07-insight-h2">One process, two faces</h2>
      <div className="v1-ch07-insight-faces">
        <div className="v1-ch07-face">
          <span className="v1-ch07-face-tag">Labour-process</span>
          <h3 className="v1-ch07-face-h3">Production of use-values</h3>
          <span className="v1-ch07-face-gloss">
            Purposeful activity, raw material, and instruments combine to make
            something useful. A process between man and nature.
          </span>
        </div>
        <div className="v1-ch07-face">
          <span className="v1-ch07-face-tag">Valorization-process</span>
          <h3 className="v1-ch07-face-h3">Production of surplus-value</h3>
          <span className="v1-ch07-face-gloss">
            The same activity produces new value beyond the value of
            labour-power advanced. A process between capital and the worker.
          </span>
        </div>
      </div>
      <p className="v1-ch07-insight-prose">
        Marx insists they are one and the same act, taken from two angles. The
        bars below split that single act first by time (necessary vs surplus
        labour) and then by value (c + v + s).
      </p>
    </section>
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

  function applyPreset(p: FixturePreset) {
    setDuration(p.duration);
    setProductKind(p.productKind);
    setProductQty(p.productQty);
    setRawMaterials(p.rawMaterials.map((rm) => ({ ...rm })));
    setInstruments(p.instruments.map((i) => ({ ...i })));
  }

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
      <div className="v1-ch07-presets">
        <span className="v1-ch07-presets-label">Marx fixture</span>
        {FIXTURES.map((p) => (
          <button
            key={p.id}
            type="button"
            className="v1-ch07-preset-button"
            onClick={() => applyPreset(p)}
          >
            {p.label}
          </button>
        ))}
      </div>
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
  const c = Math.max(v.product_value - lp.duration, 0);
  const necessary = v.necessary_labour;
  const surplus = v.surplus_labour;
  const dayTotal = lp.duration;
  const valueTotal = v.product_value;

  const surplusRate =
    necessary > 0 ? ((surplus / necessary) * 100).toFixed(1) : "∞";

  const dayPct = (n: number) => (dayTotal > 0 ? `${(n / dayTotal) * 100}%` : "0%");
  const valPct = (n: number) => (valueTotal > 0 ? `${(n / valueTotal) * 100}%` : "0%");

  return (
    <section className="card">
      <h2>Valorization Result</h2>
      <p className="small muted">Process ID: <span className="monospace">{lp.id}</span></p>

      <h3 className="v1-ch07-result-h3">Working day</h3>
      <div className="v1-ch07-partition" role="img" aria-label="Working-day partition">
        <span className="v1-ch07-partition-label">{minutesToHours(dayTotal)} day</span>
        <div className="v1-ch07-partition-track">
          <div
            className="v1-ch07-partition-seg v1-ch07-partition-seg--necessary"
            style={{ flexBasis: dayPct(necessary) }}
            title={`Necessary labour: ${minutesToHours(necessary)}`}
          >
            v
          </div>
          <div
            className="v1-ch07-partition-seg v1-ch07-partition-seg--surplus"
            style={{ flexBasis: dayPct(surplus) }}
            title={`Surplus labour: ${minutesToHours(surplus)}`}
          >
            s
          </div>
        </div>
        <span className="v1-ch07-partition-total">{minutesToHours(dayTotal)}</span>
      </div>
      <p className="v1-ch07-legend">
        <span>
          <span className="v1-ch07-legend-swatch" style={{ background: "var(--lead)" }} />
          Necessary labour ({minutesToHours(necessary)})
        </span>
        <span>
          <span className="v1-ch07-legend-swatch" style={{ background: "var(--gold-bright)" }} />
          Surplus labour ({minutesToHours(surplus)})
        </span>
      </p>

      <h3 className="v1-ch07-result-h3">Value composition</h3>
      <div className="v1-ch07-partition" role="img" aria-label="Value composition c + v + s">
        <span className="v1-ch07-partition-label">{minutesToHours(valueTotal)} value</span>
        <div className="v1-ch07-partition-track">
          <div
            className="v1-ch07-partition-seg v1-ch07-partition-seg--constant"
            style={{ flexBasis: valPct(c) }}
            title={`c (transferred): ${minutesToHours(c)}`}
          >
            c
          </div>
          <div
            className="v1-ch07-partition-seg v1-ch07-partition-seg--necessary"
            style={{ flexBasis: valPct(necessary) }}
            title={`v (necessary): ${minutesToHours(necessary)}`}
          >
            v
          </div>
          <div
            className="v1-ch07-partition-seg v1-ch07-partition-seg--surplus"
            style={{ flexBasis: valPct(surplus) }}
            title={`s (surplus): ${minutesToHours(surplus)}`}
          >
            s
          </div>
        </div>
        <span className="v1-ch07-partition-total">{minutesToHours(valueTotal)}</span>
      </div>
      <p className="v1-ch07-legend">
        <span>
          <span className="v1-ch07-legend-swatch" style={{ background: "var(--lead-hover)" }} />
          c — transferred ({minutesToHours(c)})
        </span>
        <span>
          <span className="v1-ch07-legend-swatch" style={{ background: "var(--lead)" }} />
          v — necessary ({minutesToHours(necessary)})
        </span>
        <span>
          <span className="v1-ch07-legend-swatch" style={{ background: "var(--gold-bright)" }} />
          s — surplus ({minutesToHours(surplus)})
        </span>
      </p>

      <div className="v1-ch07-surplus-rate">
        <span className="v1-ch07-surplus-rate-label">Rate of surplus-value (s/v)</span>
        <span className="v1-ch07-surplus-rate-value">{surplusRate}%</span>
      </div>

      <p className="small" style={{ marginTop: "1rem" }}>
        <strong>Product:</strong> {lp.product_quantity} × {product.commodity_kind} —
        total value {minutesToHours(product.total_value)}
      </p>
    </section>
  );
}

function Coda() {
  return (
    <aside className="v1-ch07-coda">
      <p className="v1-ch07-coda-quote">
        “By turning his money into commodities that serve as the material
        elements of a new product, and as factors in the labour-process … the
        capitalist drinks up the labour-power of others, transforming his money
        not merely into commodities, but into capital.”
        <span className="v1-ch07-coda-cite">
          — Marx, Capital Vol. I, Ch. 7
        </span>
      </p>
    </aside>
  );
}
