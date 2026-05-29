import { useEffect, useState } from "react";
import { financeApi } from "../../api";
import type { TurnoverAnalysisResponse } from "../../types";
import "./Ch04TurnoverEffect.css";

// Rates arrive in basis points (10000 bp = 100%); render as percent.
function pct(bp: number): string {
  return `${(bp / 100).toLocaleString("en-GB", { maximumFractionDigits: 2 })}%`;
}

// Capital magnitudes are LabourMinutes; Marx denominates them in £, so render
// them as whole pounds.
function money(n: number): string {
  return `£${n.toLocaleString("en-GB")}`;
}

// Client-side mirrors of the Go integer formulas, so the live slider matches the
// server to the last basis point. Inputs are non-negative, so truncating
// division is Math.floor.
function annualProfitBp(c: number, v: number, sRateBp: number, n: number): number {
  return c > 0 ? Math.floor((sRateBp * n * v) / c) : 0;
}
function singleTurnoverBp(c: number, v: number, sRateBp: number): number {
  return c > 0 ? Math.floor((sRateBp * v) / c) : 0;
}

interface CapitalPreset {
  label: string;
  note: string;
  c: number;
  v: number;
  s_rate: number;
  n: number;
}

// Marx's §main worked capitals from Ch. 4.
const CAPITAL_PRESETS: CapitalPreset[] = [
  {
    label: "Capital A — n = 2",
    note: "80c+20v = 100C, s′ = 100%, turned over twice → annual p′ = 40%.",
    c: 100,
    v: 20,
    s_rate: 10000,
    n: 2,
  },
  {
    label: "Capital B — n = 1",
    note: "160c+40v = 200C, s′ = 100%, turned over once → annual p′ = 20%. Equal composition to A, half the rate.",
    c: 200,
    v: 40,
    s_rate: 10000,
    n: 1,
  },
  {
    label: "Capital I — n = 10",
    note: "£10,000 fixed + £500c + £500v = £11,000C, s′ = 100%, ten turnovers → annual p′ = 45 5/11%.",
    c: 11000,
    v: 500,
    s_rate: 10000,
    n: 10,
  },
];

export function Ch04TurnoverEffect() {
  // §0 — the live formula slider.
  const [c, setC] = useState(100);
  const [v, setV] = useState(20);
  const [sRatePct, setSRatePct] = useState(100);
  const [n, setN] = useState(2);
  const sRateBp = Math.round(sRatePct * 100);
  const single = singleTurnoverBp(c, v, sRateBp);
  const annual = annualProfitBp(c, v, sRateBp, n);
  const annualSurplus = sRateBp * n;
  const annualWages = v * n;

  // The persisted turnover analyser.
  const [result, setResult] = useState<TurnoverAnalysisResponse | null>(null);
  const [saved, setSaved] = useState<TurnoverAnalysisResponse[]>([]);
  const [error, setError] = useState<string | null>(null);

  async function refreshSaved() {
    try {
      setSaved(await financeApi.listTurnoverAnalyses());
    } catch {
      // the seed list is best-effort; the panel still works without it
    }
  }

  useEffect(() => {
    refreshSaved();
  }, []);

  async function analyse(p: CapitalPreset) {
    setError(null);
    try {
      const a = await financeApi.createTurnoverAnalysis({ c: p.c, v: p.v, s_rate: p.s_rate, n: p.n });
      setResult(a);
      await refreshSaved();
    } catch (e) {
      setError(String(e));
    }
  }

  return (
    <div className="v3-ch04">
      <section className="v3-ch04-intro">
        <p>
          Chapters 2–3 measured the rate of profit on a single turnover. But a
          variable capital yields its surplus-value not once a year — once for
          every time it turns over. With the annual rate of surplus-value{" "}
          <strong>S′ = s′·n</strong>, the rate of profit corrected for turnover
          becomes:
        </p>
        <p className="v3-ch04-equation">p′ = s′ · n · (v / C)</p>
        <p>
          Two capitals of equal composition and equal s′ then have rates of
          profit <em>related inversely as their periods of turnover</em>: the one
          that turns over twice a year earns twice the annual rate of the one
          that turns over once.
        </p>
      </section>

      <section className="v3-ch04-slider" aria-label="Live annual-profit-rate formula">
        <h4>p′ = s′ · n · (v / C)</h4>
        <div className="v3-ch04-controls">
          <label className="v3-ch04-control">
            <span>Total capital C — {money(c)}</span>
            <input type="range" min={0} max={1000} step={10} value={c}
              onChange={(e) => setC(Number(e.target.value))} />
          </label>
          <label className="v3-ch04-control">
            <span>Variable capital advanced v — {money(v)}</span>
            <input type="range" min={0} max={400} step={5} value={v}
              onChange={(e) => setV(Number(e.target.value))} />
          </label>
          <label className="v3-ch04-control">
            <span>Rate of surplus-value s′ — {sRatePct}%</span>
            <input type="range" min={0} max={300} step={5} value={sRatePct}
              onChange={(e) => setSRatePct(Number(e.target.value))} />
          </label>
          <label className="v3-ch04-control">
            <span>Turnovers per year n — {n}</span>
            <input type="range" min={0} max={12} step={1} value={n}
              onChange={(e) => setN(Number(e.target.value))} />
          </label>
        </div>
        <div className="v3-ch04-readout">
          <div className="v3-ch04-stat">
            <span className="v3-ch04-stat-label">Single-turnover p′ = s′·(v/C)</span>
            <span className="v3-ch04-stat-value">{pct(single)}</span>
          </div>
          <div className="v3-ch04-stat v3-ch04-stat-headline">
            <span className="v3-ch04-stat-label">Annual p′ = s′·n·(v/C)</span>
            <span className="v3-ch04-stat-value">{pct(annual)}</span>
          </div>
          <div className="v3-ch04-stat">
            <span className="v3-ch04-stat-label">Annual s.-value rate S′ = s′·n</span>
            <span className="v3-ch04-stat-value">{pct(annualSurplus)}</span>
          </div>
          <div className="v3-ch04-stat">
            <span className="v3-ch04-stat-label">Annual wage bill vn</span>
            <span className="v3-ch04-stat-value">{money(annualWages)}</span>
          </div>
        </div>
        <p className="v3-ch04-slider-note">
          Turning over <strong>{n}×</strong> a year multiplies the single-turnover
          rate <strong>{pct(single)}</strong> into the annual rate{" "}
          <strong>{pct(annual)}</strong>. The capitalist's books show only the
          year's wages <strong>{money(annualWages)}</strong> (vn), never the
          variable capital <strong>{money(v)}</strong> (v) the formula needs.
        </p>
      </section>

      <section className="v3-ch04-presets-section" aria-label="Marx's worked capitals">
        <h4>Marx's worked capitals</h4>
        <div className="v3-ch04-presets">
          {CAPITAL_PRESETS.map((p) => (
            <button key={p.label} type="button" className="v3-ch04-preset" onClick={() => analyse(p)}>
              {p.label}
            </button>
          ))}
        </div>
        {result && (
          <div className="v3-ch04-result">
            <div className="v3-ch04-shift">
              <span className="v3-ch04-shift-single">{pct(result.single_turnover_profit_rate)}</span>
              <span className="v3-ch04-shift-arrow">×{result.n} →</span>
              <span className="v3-ch04-shift-annual">{pct(result.annual_profit_rate.basis_points)}</span>
              <span className="v3-ch04-shift-word">annual p′</span>
            </div>
            <p className="v3-ch04-result-detail">
              {money(result.c)} total capital, {money(result.v)} variable, s′ ={" "}
              {pct(result.s_rate)} over {result.n} turnover{result.n === 1 ? "" : "s"}. Annual
              rate of surplus-value S′ = {pct(result.annual_surplus_value_rate)}; year's wage
              bill vn = {money(result.annual_wages)}.
            </p>
          </div>
        )}
      </section>

      <section className="v3-ch04-manchester" aria-label="Manchester cotton spinnery">
        <h4>From the living experience of Manchester</h4>
        <p>
          Marx's real-world illustration — the 10,000-spindle cotton spinnery of
          April 1871: total capital <strong>£12,500</strong>, variable capital
          advanced <strong>£318</strong>, rate of surplus-value{" "}
          <strong>153 11/13%</strong>, the variable capital turned over{" "}
          <strong>≈ 8½ times</strong> a year (annual wage bill £2,704). The annual
          rate of profit comes to <strong>≈ 33.27%</strong> and the annual rate of
          surplus-value to a startling <strong>≈ 1,307%</strong> — "such a rate is
          by no means a rarity" in times of greatest prosperity.
        </p>
        <p className="v3-ch04-manchester-note">
          The model counts whole turnovers, so Marx's fractional 8½ is shown here
          as his own figures rather than recomputed.
        </p>
      </section>

      {error && <p className="v3-ch04-error">{error}</p>}

      {saved.length > 0 && (
        <section className="v3-ch04-saved" aria-label="Recorded turnover analyses">
          <h4>Recorded turnover analyses</h4>
          <table className="v3-ch04-table">
            <thead>
              <tr>
                <th>C</th>
                <th>v</th>
                <th>s′</th>
                <th>n</th>
                <th>single p′</th>
                <th>annual p′</th>
                <th>S′</th>
              </tr>
            </thead>
            <tbody>
              {saved.map((a) => (
                <tr key={a.id}>
                  <td>{money(a.c)}</td>
                  <td>{money(a.v)}</td>
                  <td>{pct(a.s_rate)}</td>
                  <td>{a.n}</td>
                  <td>{pct(a.single_turnover_profit_rate)}</td>
                  <td>{pct(a.annual_profit_rate.basis_points)}</td>
                  <td>{pct(a.annual_surplus_value_rate)}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </section>
      )}
    </div>
  );
}
