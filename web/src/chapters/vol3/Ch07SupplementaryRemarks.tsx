import { useEffect, useState } from "react";
import { financeApi } from "../../api";
import type {
  CompositionEffectAnalysisResponse,
  MagnitudeChangeAnalysisResponse,
  MagnitudeChangeKind,
  PartISummaryResponse,
} from "../../types";
import "./Ch07SupplementaryRemarks.css";

// Rates arrive in basis points (10000 bp = 100%); render as percent.
function pct(bp: number): string {
  return `${(bp / 100).toLocaleString("en-GB", { maximumFractionDigits: 2 })}%`;
}

function money(n: number): string {
  return `£${n.toLocaleString("en-GB")}`;
}

// Client-side mirrors of the Go basis-point formulae, truncating to match the
// server to the last basis point.
function profitRateBp(s: number, c: number, v: number): number {
  const total = c + v;
  return total > 0 ? Math.floor((s * 10000) / total) : 0;
}
function magnitudeRateBp(profit: number, capital: number): number {
  return capital > 0 ? Math.floor((profit * 10000) / capital) : 0;
}

const MAGNITUDE_LABELS: Record<MagnitudeChangeKind, string> = {
  money_value_change: "Money-value change (nominal — gold rises/falls)",
  proportional_change: "Proportional change (same composition)",
  organic_change: "Organic change (composition shifts)",
};

const MAGNITUDE_KINDS = Object.keys(MAGNITUDE_LABELS) as MagnitudeChangeKind[];

export function Ch07SupplementaryRemarks() {
  // §main — organic composition: two enterprises, equal s and v, different c.
  const [sA, setSA] = useState(1000);
  const [vA, setVA] = useState(1000);
  const [cA, setCA] = useState(9000);
  const [cB, setCB] = useState(11000);

  const pA = profitRateBp(sA, cA, vA);
  const pB = profitRateBp(sA, cB, vA); // equal s and v: the composition alone differs
  const rateDiff = Math.abs(pA - pB);

  // §Rodbertus — capital-magnitude change.
  const [kind, setKind] = useState<MagnitudeChangeKind>("money_value_change");
  const [origCapital, setOrigCapital] = useState(100);
  const [origProfit, setOrigProfit] = useState(20);
  const [factor, setFactor] = useState(200);

  const oldRate = magnitudeRateBp(origProfit, origCapital);
  const newRate =
    kind === "organic_change"
      ? magnitudeRateBp(origProfit, Math.floor((origCapital * factor) / 100))
      : oldRate;
  const rateUnchanged = newRate === oldRate;

  const [compositions, setCompositions] = useState<CompositionEffectAnalysisResponse[]>([]);
  const [magnitudes, setMagnitudes] = useState<MagnitudeChangeAnalysisResponse[]>([]);
  const [summary, setSummary] = useState<PartISummaryResponse | null>(null);
  const [error, setError] = useState<string | null>(null);

  async function refresh() {
    try {
      const [c, m, s] = await Promise.all([
        financeApi.listCompositionEffects(),
        financeApi.listMagnitudeChanges(),
        financeApi.getPartISummary(),
      ]);
      setCompositions(c);
      setMagnitudes(m);
      setSummary(s);
    } catch {
      // best-effort; the panel still works without the seed list
    }
  }

  useEffect(() => {
    refresh();
  }, []);

  async function recordComposition() {
    setError(null);
    try {
      await financeApi.computeCompositionEffect({ s_a: sA, v_a: vA, c_a: cA, s_b: sA, v_b: vA, c_b: cB });
      await refresh();
    } catch (e) {
      setError(String(e));
    }
  }

  async function recordMagnitude() {
    setError(null);
    try {
      await financeApi.computeMagnitudeChange({
        kind,
        original_capital: origCapital,
        original_profit: origProfit,
        factor,
      });
      await refresh();
    } catch (e) {
      setError(String(e));
    }
  }

  return (
    <div className="v3-ch07">
      <section className="v3-ch07-intro">
        <p>
          The supplementary remarks close Part I. Two things are settled. First,
          with the same surplus-value on the same variable capital, the{" "}
          <strong>organic composition</strong> alone sets the rate of profit: the
          capital carrying more constant capital shows the lower rate. Second,{" "}
          <strong>Rodbertus's error</strong> — a mere change in the money-value of
          capital, or a proportional change at unchanged composition, moves the rate
          of profit not at all; only a change in the composition does.
        </p>
      </section>

      <section className="v3-ch07-calc" aria-label="Organic composition comparison">
        <h4>Organic composition — same s and v, different c</h4>
        <div className="v3-ch07-controls">
          <label className="v3-ch07-control">
            <span>Surplus-value s (both) — {money(sA)}</span>
            <input type="range" min={0} max={5000} step={100} value={sA}
              onChange={(e) => setSA(Number(e.target.value))} />
          </label>
          <label className="v3-ch07-control">
            <span>Variable capital v (both) — {money(vA)}</span>
            <input type="range" min={0} max={5000} step={100} value={vA}
              onChange={(e) => setVA(Number(e.target.value))} />
          </label>
          <label className="v3-ch07-control">
            <span>Enterprise A — constant capital c — {money(cA)}</span>
            <input type="range" min={0} max={20000} step={100} value={cA}
              onChange={(e) => setCA(Number(e.target.value))} />
          </label>
          <label className="v3-ch07-control">
            <span>Enterprise B — constant capital c — {money(cB)}</span>
            <input type="range" min={0} max={20000} step={100} value={cB}
              onChange={(e) => setCB(Number(e.target.value))} />
          </label>
        </div>
        <div className="v3-ch07-versus">
          <div className="v3-ch07-firm">
            <span className="v3-ch07-firm-name">Enterprise A</span>
            <span className="v3-ch07-firm-rate">{pct(pA)}</span>
            <span className="v3-ch07-firm-detail">{money(cA)}c + {money(vA)}v + {money(sA)}s</span>
          </div>
          <div className="v3-ch07-firm">
            <span className="v3-ch07-firm-name">Enterprise B</span>
            <span className="v3-ch07-firm-rate">{pct(pB)}</span>
            <span className="v3-ch07-firm-detail">{money(cB)}c + {money(vA)}v + {money(sA)}s</span>
          </div>
          <div className="v3-ch07-firm v3-ch07-firm-diff">
            <span className="v3-ch07-firm-name">Difference</span>
            <span className="v3-ch07-firm-rate">{pct(rateDiff)}</span>
            <span className="v3-ch07-firm-detail">the composition alone</span>
          </div>
        </div>
        <button type="button" className="v3-ch07-record" onClick={recordComposition}>
          Record this comparison
        </button>
      </section>

      <section className="v3-ch07-calc" aria-label="Rodbertus refutation">
        <h4>Refutation of Rodbertus — does the rate move?</h4>
        <label className="v3-ch07-control">
          <span>Kind of magnitude change</span>
          <select value={kind} onChange={(e) => setKind(e.target.value as MagnitudeChangeKind)}>
            {MAGNITUDE_KINDS.map((k) => (
              <option key={k} value={k}>{MAGNITUDE_LABELS[k]}</option>
            ))}
          </select>
        </label>
        <div className="v3-ch07-controls">
          <label className="v3-ch07-control">
            <span>Original capital C — {money(origCapital)}</span>
            <input type="range" min={0} max={5000} step={50} value={origCapital}
              onChange={(e) => setOrigCapital(Number(e.target.value))} />
          </label>
          <label className="v3-ch07-control">
            <span>Original profit p — {money(origProfit)}</span>
            <input type="range" min={0} max={2000} step={20} value={origProfit}
              onChange={(e) => setOrigProfit(Number(e.target.value))} />
          </label>
          <label className="v3-ch07-control">
            <span>Change factor — {factor}% of original</span>
            <input type="range" min={25} max={400} step={5} value={factor}
              onChange={(e) => setFactor(Number(e.target.value))} />
          </label>
        </div>
        <div className="v3-ch07-readout">
          <div className="v3-ch07-stat">
            <span className="v3-ch07-stat-label">Old rate p/C</span>
            <span className="v3-ch07-stat-value">{pct(oldRate)}</span>
          </div>
          <div className="v3-ch07-stat v3-ch07-stat-headline">
            <span className="v3-ch07-stat-label">New rate</span>
            <span className="v3-ch07-stat-value">{pct(newRate)}</span>
          </div>
          <div className={`v3-ch07-verdict ${rateUnchanged ? "v3-ch07-unchanged" : "v3-ch07-moved"}`}>
            {rateUnchanged ? "Rate unchanged — Rodbertus refuted" : "Rate moved — an organic change"}
          </div>
        </div>
        <button type="button" className="v3-ch07-record" onClick={recordMagnitude}>
          Record this change
        </button>
      </section>

      {error && <p className="v3-ch07-error">{error}</p>}

      {summary && (
        <section className="v3-ch07-summary" aria-label="Part I summary">
          <h4>Part I — the determinants of the rate of profit</h4>
          <ul className="v3-ch07-determinants">
            {summary.determinants.map((d) => (
              <li key={d}>{d}</li>
            ))}
          </ul>
          <p className="small muted">{summary.analyses.length} rate-of-profit analyses recorded across Part I.</p>
        </section>
      )}

      {compositions.length > 0 && (
        <section className="v3-ch07-recorded" aria-label="Recorded comparisons">
          <h4>Recorded composition comparisons</h4>
          <table className="v3-ch07-table">
            <thead>
              <tr>
                <th>A: c/v/s</th>
                <th>B: c/v/s</th>
                <th>p′ A</th>
                <th>p′ B</th>
                <th>difference</th>
              </tr>
            </thead>
            <tbody>
              {compositions.map((a) => (
                <tr key={a.id}>
                  <td>{money(a.effect.c_a)}/{money(a.effect.v_a)}/{money(a.effect.s_a)}</td>
                  <td>{money(a.effect.c_b)}/{money(a.effect.v_b)}/{money(a.effect.s_b)}</td>
                  <td>{pct(a.profit_rate_a)}</td>
                  <td>{pct(a.profit_rate_b)}</td>
                  <td>{pct(a.rate_difference)}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </section>
      )}

      {magnitudes.length > 0 && (
        <section className="v3-ch07-recorded" aria-label="Recorded magnitude changes">
          <h4>Recorded magnitude changes</h4>
          <table className="v3-ch07-table">
            <thead>
              <tr>
                <th>Kind</th>
                <th>C</th>
                <th>p</th>
                <th>factor</th>
                <th>old</th>
                <th>new</th>
                <th>rate</th>
              </tr>
            </thead>
            <tbody>
              {magnitudes.map((a) => (
                <tr key={a.id}>
                  <td>{MAGNITUDE_LABELS[a.change.kind]}</td>
                  <td>{money(a.change.original_capital)}</td>
                  <td>{money(a.change.original_profit)}</td>
                  <td>{a.change.factor}%</td>
                  <td>{pct(a.old_rate)}</td>
                  <td>{pct(a.new_rate)}</td>
                  <td>{a.rate_unchanged ? "unchanged" : "moved"}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </section>
      )}
    </div>
  );
}
