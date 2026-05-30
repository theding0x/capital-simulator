import { useEffect, useState } from "react";
import { financeApi } from "../../api";
import type { ProductionSphereResponse } from "../../types";
import "./Ch08CompositionDifferences.css";

// Rates arrive in basis points (10000 bp = 100%); render as percent.
function pct(bp: number): string {
  return `${(bp / 100).toLocaleString("en-GB", { maximumFractionDigits: 2 })}%`;
}

function money(n: number): string {
  return `£${n.toLocaleString("en-GB")}`;
}

// Client-side mirror of the Go round-half-up formula:
//   s = floor(sRate * v / 10000)
//   p' = round(s * 10000 / C)   [Math.round gives half-up for positive values]
function clientSurplusValue(sRate: number, v: number): number {
  return Math.floor((sRate * v) / 10000);
}
function clientProfitRate(s: number, c: number): number {
  return c > 0 ? Math.round((s * 10000) / c) : 0;
}
function clientVariablePercent(v: number, c: number): number {
  return c > 0 ? Math.floor((v * 100) / c) : 0;
}
function clientLabourPowerIndex(v: number, c: number): number {
  return c > 0 ? Math.floor((v * 3600 * 100) / c) : 0;
}

// Marx's international comparison presets.
const PRESETS = [
  { name: "Sphere A (600c+100v)", c: 700, v: 100, sRate: 10000 },
  { name: "Sphere B (100c+600v)", c: 700, v: 600, sRate: 10000 },
  { name: "European capital (84c+16v)", c: 100, v: 16, sRate: 10000 },
  { name: "Asian capital (16c+84v)", c: 100, v: 84, sRate: 2500 },
];

export function Ch08CompositionDifferences() {
  const [spheres, setSpheres] = useState<ProductionSphereResponse[]>([]);
  const [error, setError] = useState<string | null>(null);

  // Builder state.
  const [name, setName] = useState("New sphere");
  const [c, setC] = useState(700);
  const [v, setV] = useState(100);
  const [sRate, setSRate] = useState(10000);

  const s = clientSurplusValue(sRate, v);
  const profitRate = clientProfitRate(s, c);
  const varPct = clientVariablePercent(v, c);
  const lpi = clientLabourPowerIndex(v, c);
  const constCapital = c - v;

  async function refresh() {
    try {
      const items = await financeApi.listProductionSpheres();
      setSpheres(items);
    } catch {
      // best-effort; the table still renders what was loaded
    }
  }

  useEffect(() => {
    refresh();
  }, []);

  async function recordSphere() {
    setError(null);
    try {
      await financeApi.createProductionSphere({ name, c, v, s_rate: sRate });
      await refresh();
    } catch (e) {
      setError(String(e));
    }
  }

  function applyPreset(p: typeof PRESETS[0]) {
    setName(p.name);
    setC(p.c);
    setV(p.v);
    setSRate(p.sRate);
  }

  return (
    <div className="v3-ch08">
      <section className="v3-ch08-intro">
        <p>
          With the <strong>same rate of surplus-value</strong> in every branch
          but different organic compositions — different v/C splits — the{" "}
          <strong>individual rate of profit differs</strong>. Living labour, and
          therefore surplus-value, is set in motion only by the variable capital.
          Sphere B (100c + 600v) sets six times as much living labour in motion
          per £700 as Sphere A (600c + 100v), and shows a rate of profit six
          times as high. The equalisation of these divergent individual rates
          into a general (average) rate is deferred to Ch. 9.
        </p>
      </section>

      <section className="v3-ch08-calc" aria-label="Production-sphere builder">
        <h4>Sphere builder — p′ = s / C</h4>
        <div className="v3-ch08-presets">
          {PRESETS.map((p) => (
            <button
              key={p.name}
              type="button"
              className="v3-ch08-preset"
              onClick={() => applyPreset(p)}
            >
              {p.name}
            </button>
          ))}
        </div>
        <div className="v3-ch08-controls">
          <label className="v3-ch08-control">
            <span>Branch name</span>
            <input
              type="text"
              value={name}
              maxLength={80}
              onChange={(e) => setName(e.target.value)}
              className="v3-ch08-text-input"
            />
          </label>
          <label className="v3-ch08-control">
            <span>Total capital C — {money(c)}</span>
            <input
              type="range"
              min={0}
              max={2000}
              step={10}
              value={c}
              onChange={(e) => setC(Number(e.target.value))}
            />
          </label>
          <label className="v3-ch08-control">
            <span>Variable capital v — {money(Math.min(v, c))}</span>
            <input
              type="range"
              min={0}
              max={c}
              step={10}
              value={Math.min(v, c)}
              onChange={(e) => setV(Number(e.target.value))}
            />
          </label>
          <label className="v3-ch08-control">
            <span>Rate of surplus-value s′ — {pct(sRate)}</span>
            <input
              type="range"
              min={0}
              max={20000}
              step={100}
              value={sRate}
              onChange={(e) => setSRate(Number(e.target.value))}
            />
          </label>
        </div>
        <div className="v3-ch08-readout">
          <div className="v3-ch08-stat">
            <span className="v3-ch08-stat-label">Constant capital c</span>
            <span className="v3-ch08-stat-value">{money(constCapital)}</span>
          </div>
          <div className="v3-ch08-stat">
            <span className="v3-ch08-stat-label">Surplus-value s</span>
            <span className="v3-ch08-stat-value">{money(s)}</span>
          </div>
          <div className="v3-ch08-stat v3-ch08-stat-headline">
            <span className="v3-ch08-stat-label">Individual p′ = s/C</span>
            <span className="v3-ch08-stat-value">{pct(profitRate)}</span>
          </div>
          <div className="v3-ch08-stat">
            <span className="v3-ch08-stat-label">v/C ratio</span>
            <span className="v3-ch08-stat-value">{varPct}%</span>
          </div>
          <div className="v3-ch08-stat">
            <span className="v3-ch08-stat-label">Labour-power index</span>
            <span className="v3-ch08-stat-value">{lpi.toLocaleString("en-GB")}</span>
          </div>
        </div>
        <button type="button" className="v3-ch08-record" onClick={recordSphere}>
          Record sphere
        </button>
      </section>

      <section className="v3-ch08-compare" aria-label="International comparison">
        <h4>International comparison — Asian vs European capital</h4>
        <p className="v3-ch08-compare-note">
          The Asian capital (16c + 84v, s′ = 25%) shows a{" "}
          <strong>higher</strong> individual p′ (21%) than the European capital
          (84c + 16v, s′ = 100%, p′ = 16%), despite its lower rate of
          surplus-value — because its larger v/C sets so much more living labour
          in motion per £100 of total capital that the surplus-value generated
          outweighs the disadvantage in s′. The profit-form completely conceals
          this mechanism.
        </p>
        <div className="v3-ch08-intl-grid">
          {PRESETS.slice(2).map((p) => {
            const sv = clientSurplusValue(p.sRate, p.v);
            const pr = clientProfitRate(sv, p.c);
            const vp = clientVariablePercent(p.v, p.c);
            return (
              <div key={p.name} className="v3-ch08-intl-card">
                <span className="v3-ch08-intl-name">{p.name}</span>
                <span className="v3-ch08-intl-rate">{pct(pr)}</span>
                <span className="v3-ch08-intl-detail">
                  {money(p.c - p.v)}c + {money(p.v)}v · s′ {pct(p.sRate)} · v/C {vp}%
                </span>
              </div>
            );
          })}
        </div>
      </section>

      {error && <p className="v3-ch08-error">{error}</p>}

      {spheres.length > 0 && (
        <section className="v3-ch08-recorded" aria-label="Recorded spheres">
          <h4>Recorded production spheres</h4>
          <table className="v3-ch08-table">
            <thead>
              <tr>
                <th>Branch</th>
                <th>C</th>
                <th>c</th>
                <th>v</th>
                <th>c:v</th>
                <th>s′</th>
                <th>s</th>
                <th>p′</th>
              </tr>
            </thead>
            <tbody>
              {spheres.map((sp) => (
                <tr key={sp.id}>
                  <td>{sp.name}</td>
                  <td>{money(sp.c)}</td>
                  <td>{money(sp.constant_capital)}</td>
                  <td>{money(sp.v)}</td>
                  <td>{sp.constant_capital}:{sp.v}</td>
                  <td>{pct(sp.s_rate)}</td>
                  <td>{money(sp.surplus_value)}</td>
                  <td>{pct(sp.individual_profit_rate)}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </section>
      )}
    </div>
  );
}
