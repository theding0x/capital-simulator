import { useEffect, useState } from "react";
import { financeApi } from "../../api";
import type { PriceFluctuationAnalysisResponse } from "../../types";
import "./Ch06PriceFluctuation.css";

// Rates arrive in basis points (10000 bp = 100%); render as percent.
function pct(bp: number): string {
  return `${(bp / 100).toLocaleString("en-GB", { maximumFractionDigits: 2 })}%`;
}

// Capital magnitudes are LabourMinutes; Marx denominates them in £.
function money(n: number): string {
  return `£${n.toLocaleString("en-GB")}`;
}

// Client-side mirror of the Go profitRateBP: s / (c + v) in basis points, with
// truncating (floor) division to match the server to the last basis point.
function profitRateBp(s: number, c: number, v: number): number {
  const total = c + v;
  return total > 0 ? Math.floor((s * 10000) / total) : 0;
}

interface FluctuationPreset {
  label: string;
  note: string;
  fixed: number;
  material: number;
  factor: number;
  v: number;
  s: number;
}

// The three seeded fluctuations, drawn from Ch. 6's §I illustrations. All start
// from c = £400, v = £100, s = £100 (p' = 20%) and move only the raw-material price.
const PRESETS: FluctuationPreset[] = [
  {
    label: "Cotton doubles (§I)",
    note: "£320 fixed + £80 raw material. Doubling the price lifts c to £480; p′ falls 20% → 17.24%.",
    fixed: 320,
    material: 80,
    factor: 200,
    v: 100,
    s: 100,
  },
  {
    label: "Cotton halves (§I)",
    note: "£20 fixed + £380 raw material. Halving the price drops c to £210; p′ rises 20% → 32.26%.",
    fixed: 20,
    material: 380,
    factor: 50,
    v: 100,
    s: 100,
  },
  {
    label: "Cotton falls a quarter (§I)",
    note: "£40 fixed + £360 raw material. A 25% fall drops c to £310; p′ rises 20% → 24.39%.",
    fixed: 40,
    material: 360,
    factor: 75,
    v: 100,
    s: 100,
  },
];

export function Ch06PriceFluctuation() {
  const [fixed, setFixed] = useState(320);
  const [material, setMaterial] = useState(80);
  const [factor, setFactor] = useState(200);
  const [v, setV] = useState(100);
  const [s, setS] = useState(100);

  // §II — the capital tied up in a stock of raw material, released by a fall and
  // absorbed by a rise. The two are one magnitude with opposite sign.
  const [stockMonths, setStockMonths] = useState(3);
  const [stockUnits, setStockUnits] = useState(1000);
  const [oldPrice, setOldPrice] = useState(10);
  const [newPrice, setNewPrice] = useState(7);

  // Only the circulating raw-material component is re-priced; the fixed component
  // and v and s stay put. p' = s / (c + v) before and after the movement.
  const newMaterial = Math.floor((material * factor) / 100);
  const oldC = fixed + material;
  const newC = fixed + newMaterial;
  const oldP = profitRateBp(s, oldC, v);
  const newP = profitRateBp(s, newC, v);
  const kind: "rise" | "fall" = factor > 100 ? "rise" : "fall";

  // §II release (price fall, positive) / tie-up (price rise, negative).
  const released = stockMonths * stockUnits * (oldPrice - newPrice);

  const [result, setResult] = useState<PriceFluctuationAnalysisResponse | null>(null);
  const [recorded, setRecorded] = useState<PriceFluctuationAnalysisResponse[]>([]);
  const [error, setError] = useState<string | null>(null);

  async function refreshRecorded() {
    try {
      setRecorded(await financeApi.listPriceFluctuationAnalyses());
    } catch {
      // the seed list is best-effort; the panel still works without it
    }
  }

  useEffect(() => {
    refreshRecorded();
  }, []);

  function applyPreset(p: FluctuationPreset) {
    setFixed(p.fixed);
    setMaterial(p.material);
    setFactor(p.factor);
    setV(p.v);
    setS(p.s);
  }

  async function record() {
    setError(null);
    try {
      const a = await financeApi.createPriceFluctuationAnalysis({
        fixed,
        original_material: material,
        price_factor: factor,
        variable: v,
        surplus_value: s,
      });
      setResult(a);
      await refreshRecorded();
    } catch (e) {
      setError(String(e));
    }
  }

  return (
    <div className="v3-ch06">
      <section className="v3-ch06-intro">
        <p>
          The rate of profit <strong>p′ = s / (c + v)</strong> moves with the price
          of the raw materials that make up the circulating part of the constant
          capital — though that price-movement adds not a farthing of surplus-value.
          The surplus <strong>s</strong> is wrung from living labour and is
          indifferent to what the cotton cost; but it is <em>measured</em> against
          the whole capital advanced.
        </p>
        <p className="v3-ch06-equation">price of cotton ↑ → c ↑ → p′ ↓</p>
        <p>
          A rise in the price of the raw material swells <strong>c</strong> and the
          rate of profit falls; a fall shrinks <strong>c</strong> and the rate rises.
          The fixed capital is not re-priced by the day's market, so only the
          circulating raw-material component fluctuates.
        </p>
      </section>

      <section className="v3-ch06-calc" aria-label="Price fluctuation in the raw material">
        <h4>Re-price the raw material, watch p′ move</h4>

        <div className="v3-ch06-controls">
          <label className="v3-ch06-control">
            <span>Fixed capital (machinery, buildings) — {money(fixed)}</span>
            <input type="range" min={0} max={1000} step={10} value={fixed}
              onChange={(e) => setFixed(Number(e.target.value))} />
          </label>
          <label className="v3-ch06-control">
            <span>Raw-material capital (circulating) — {money(material)}</span>
            <input type="range" min={0} max={1000} step={10} value={material}
              onChange={(e) => setMaterial(Number(e.target.value))} />
          </label>
          <label className="v3-ch06-control">
            <span>Raw-material price — {factor}% of original ({factor > 100 ? "risen" : factor < 100 ? "fallen" : "unchanged"})</span>
            <input type="range" min={25} max={300} step={5} value={factor}
              onChange={(e) => setFactor(Number(e.target.value))} />
          </label>
          <label className="v3-ch06-control">
            <span>Variable capital v — {money(v)}</span>
            <input type="range" min={0} max={1000} step={10} value={v}
              onChange={(e) => setV(Number(e.target.value))} />
          </label>
          <label className="v3-ch06-control">
            <span>Surplus-value s — {money(s)}</span>
            <input type="range" min={0} max={1000} step={10} value={s}
              onChange={(e) => setS(Number(e.target.value))} />
          </label>
        </div>

        <div className="v3-ch06-readout">
          <div className="v3-ch06-stat">
            <span className="v3-ch06-stat-label">c before → after</span>
            <span className="v3-ch06-stat-value">{money(oldC)} → {money(newC)}</span>
          </div>
          <div className="v3-ch06-stat">
            <span className="v3-ch06-stat-label">Old p′ = s/(c+v)</span>
            <span className="v3-ch06-stat-value">{pct(oldP)}</span>
          </div>
          <div className={`v3-ch06-stat v3-ch06-stat-headline v3-ch06-${kind}`}>
            <span className="v3-ch06-stat-label">New p′ ({kind})</span>
            <span className="v3-ch06-stat-value">{pct(newP)}</span>
          </div>
          <div className="v3-ch06-stat">
            <span className="v3-ch06-stat-label">Movement</span>
            <span className="v3-ch06-stat-value">{newP >= oldP ? "+" : ""}{pct(newP - oldP)}</span>
          </div>
        </div>

        <p className="v3-ch06-calc-note">
          Re-pricing the raw material to <strong>{factor}%</strong> of its original
          level moves the circulating component to <strong>{money(newMaterial)}</strong>,
          so <strong>c</strong> goes from {money(oldC)} to {money(newC)} and the rate
          of profit {newP >= oldP ? "rises" : "falls"} from{" "}
          <strong>{pct(oldP)}</strong> to <strong>{pct(newP)}</strong> — with{" "}
          <strong>s</strong> and <strong>v</strong> untouched.
        </p>

        <button type="button" className="v3-ch06-record" onClick={record}>
          Record this fluctuation
        </button>
      </section>

      <section className="v3-ch06-release" aria-label="Capital release and tie-up">
        <h4>§II — capital released and tied up</h4>
        <p className="v3-ch06-sub-title">
          A stock of raw material ties up capital. When the price falls, the capital
          advanced to hold the stock is <em>released</em>; when it rises, more must
          be <em>tied up</em>. The two are one magnitude with opposite sign.
        </p>
        <div className="v3-ch06-controls">
          <label className="v3-ch06-control">
            <span>Months of stock held — {stockMonths}</span>
            <input type="range" min={1} max={12} step={1} value={stockMonths}
              onChange={(e) => setStockMonths(Number(e.target.value))} />
          </label>
          <label className="v3-ch06-control">
            <span>Units consumed per month — {stockUnits.toLocaleString("en-GB")}</span>
            <input type="range" min={0} max={5000} step={100} value={stockUnits}
              onChange={(e) => setStockUnits(Number(e.target.value))} />
          </label>
          <label className="v3-ch06-control">
            <span>Old price per unit — {money(oldPrice)}</span>
            <input type="range" min={0} max={20} step={1} value={oldPrice}
              onChange={(e) => setOldPrice(Number(e.target.value))} />
          </label>
          <label className="v3-ch06-control">
            <span>New price per unit — {money(newPrice)}</span>
            <input type="range" min={0} max={20} step={1} value={newPrice}
              onChange={(e) => setNewPrice(Number(e.target.value))} />
          </label>
        </div>
        <div className={`v3-ch06-release-card ${released >= 0 ? "v3-ch06-fall" : "v3-ch06-rise"}`}>
          <span className="v3-ch06-stat-label">
            {released >= 0 ? "Capital released (price fell)" : "Capital tied up (price rose)"}
          </span>
          <span className="v3-ch06-stat-value">{money(Math.abs(released))}</span>
        </div>
      </section>

      <section className="v3-ch06-presets-section" aria-label="Marx's worked fluctuations">
        <h4>Marx's worked fluctuations</h4>
        <div className="v3-ch06-presets">
          {PRESETS.map((p) => (
            <button key={p.label} type="button" className="v3-ch06-preset" onClick={() => applyPreset(p)}
              title={p.note}>
              {p.label}
            </button>
          ))}
        </div>
        {result && (
          <div className="v3-ch06-result">
            <div className="v3-ch06-shift">
              <span className="v3-ch06-shift-old">{pct(result.old_profit_rate)}</span>
              <span className="v3-ch06-shift-arrow">→</span>
              <span className={`v3-ch06-shift-new v3-ch06-${result.kind}`}>{pct(result.new_profit_rate)}</span>
              <span className="v3-ch06-shift-word">p′ ({result.kind})</span>
            </div>
            <p className="v3-ch06-result-detail">
              {money(result.effect.fixed_capital)} fixed + {money(result.effect.original_material_capital)} raw material,
              re-priced to {result.effect.price_factor}%, with {money(result.effect.variable_capital)}v +{" "}
              {money(result.effect.surplus_value)}s.
            </p>
          </div>
        )}
      </section>

      {error && <p className="v3-ch06-error">{error}</p>}

      {recorded.length > 0 && (
        <section className="v3-ch06-recorded" aria-label="Recorded fluctuations">
          <h4>Recorded fluctuations</h4>
          <table className="v3-ch06-table">
            <thead>
              <tr>
                <th>Direction</th>
                <th>fixed</th>
                <th>material</th>
                <th>price</th>
                <th>v</th>
                <th>s</th>
                <th>old p′</th>
                <th>new p′</th>
              </tr>
            </thead>
            <tbody>
              {recorded.map((a) => (
                <tr key={a.id}>
                  <td>{a.kind}</td>
                  <td>{money(a.effect.fixed_capital)}</td>
                  <td>{money(a.effect.original_material_capital)}</td>
                  <td>{a.effect.price_factor}%</td>
                  <td>{money(a.effect.variable_capital)}</td>
                  <td>{money(a.effect.surplus_value)}</td>
                  <td>{pct(a.old_profit_rate)}</td>
                  <td>{pct(a.new_profit_rate)}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </section>
      )}
    </div>
  );
}
