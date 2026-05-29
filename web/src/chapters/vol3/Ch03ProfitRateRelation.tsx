import { useEffect, useState } from "react";
import { financeApi } from "../../api";
import type {
  ProfitRateComparisonResponse,
  VariationAnalysisResponse,
  VariationCase,
} from "../../types";
import "./Ch03ProfitRateRelation.css";

// Rates arrive in basis points (10000 bp = 100%); render as percent.
function pct(bp: number): string {
  return `${(bp / 100).toLocaleString("en-GB", { maximumFractionDigits: 2 })}%`;
}

// Capital magnitudes are LabourMinutes; Marx denominates them in £, so render
// them as whole pounds.
function money(n: number): string {
  return `£${n.toLocaleString("en-GB")}`;
}

// Client-side mirrors of the Go integer formulas, so the live slider matches
// the server to the last basis point. Inputs are non-negative, so truncating
// division is Math.floor.
function profitRateBp(c: number, v: number, sRateBp: number): number {
  const total = c + v;
  return total > 0 ? Math.floor((sRateBp * v) / total) : 0;
}
function compositionBp(c: number, v: number): number {
  const total = c + v;
  return total > 0 ? Math.floor((v * 10000) / total) : 0;
}

const CASE_LABEL: Record<VariationCase, string> = {
  s_constant_vc_variable: "§I — s′ constant, value-composition varies",
  s_variable_vc_constant: "§II — value-composition constant, s′ varies",
  both_variable: "§II.3 — both s′ and value-composition vary",
};

interface Decomp {
  c: number;
  v: number;
  s_rate: number;
}

interface VariationPreset {
  label: string;
  note: string;
  initial: Decomp;
  changed: Decomp;
}

// Marx's worked variations from Ch. 3, §I.4 and §II.1.
const VARIATION_PRESETS: VariationPreset[] = [
  {
    label: "s′ constant · v/C falls",
    note: "80c+20v → 100c+20v: the normal case of modern industry — p′ falls.",
    initial: { c: 80, v: 20, s_rate: 10000 },
    changed: { c: 100, v: 20, s_rate: 10000 },
  },
  {
    label: "s′ constant · v/C rises",
    note: "80c+20v → 60c+20v: a lower composition lifts p′.",
    initial: { c: 80, v: 20, s_rate: 10000 },
    changed: { c: 60, v: 20, s_rate: 10000 },
  },
  {
    label: "v/C constant · s′ halves",
    note: "80c+20v at 100% → 160c+40v at 50%: p′ falls in exact proportion to s′.",
    initial: { c: 80, v: 20, s_rate: 10000 },
    changed: { c: 160, v: 40, s_rate: 5000 },
  },
];

interface ComparePreset {
  label: string;
  note: string;
  first: Decomp;
  second: Decomp;
}

const COMPARE_PRESETS: ComparePreset[] = [
  {
    label: "Equal s′ (Marx's 15,000)",
    note: "12,000c+3,000v vs 13,000c+2,000v, both at s′ = 100%: rates stand as v/C.",
    first: { c: 12000, v: 3000, s_rate: 10000 },
    second: { c: 13000, v: 2000, s_rate: 10000 },
  },
  {
    label: "Unequal s′",
    note: "80c+20v at 100% vs 160c+40v at 50%: same v/C, but rates stand as s′.",
    first: { c: 80, v: 20, s_rate: 10000 },
    second: { c: 160, v: 40, s_rate: 5000 },
  },
];

export function Ch03ProfitRateRelation() {
  // §0 — the live formula slider.
  const [c, setC] = useState(80);
  const [v, setV] = useState(20);
  const [sRatePct, setSRatePct] = useState(100);
  const sRateBp = Math.round(sRatePct * 100);
  const total = c + v;
  const compBp = compositionBp(c, v);
  const pPrimeBp = profitRateBp(c, v, sRateBp);

  // §I/§II — the persisted variation analyser.
  const [variation, setVariation] = useState<VariationAnalysisResponse | null>(null);
  const [comparison, setComparison] = useState<ProfitRateComparisonResponse | null>(null);
  const [saved, setSaved] = useState<VariationAnalysisResponse[]>([]);
  const [error, setError] = useState<string | null>(null);

  async function refreshSaved() {
    try {
      setSaved(await financeApi.listVariations());
    } catch {
      // the seed list is best-effort; the panel still works without it
    }
  }

  useEffect(() => {
    refreshSaved();
  }, []);

  async function analyse(p: VariationPreset) {
    setError(null);
    try {
      const a = await financeApi.analyseVariation({ initial: p.initial, changed: p.changed });
      setVariation(a);
      await refreshSaved();
    } catch (e) {
      setError(String(e));
    }
  }

  async function compare(p: ComparePreset) {
    setError(null);
    try {
      setComparison(await financeApi.compareProfitRates({ first: p.first, second: p.second }));
    } catch (e) {
      setError(String(e));
    }
  }

  const movement =
    variation == null
      ? ""
      : variation.new_profit_rate > variation.old_profit_rate
        ? "rose"
        : variation.new_profit_rate < variation.old_profit_rate
          ? "fell"
          : "held steady";

  return (
    <div className="v3-ch03">
      <section className="v3-ch03-intro">
        <p>
          The rate of profit is the rate of surplus-value seen through the total
          capital. Substituting <strong>s = s′v</strong> into{" "}
          <strong>p′ = s/C</strong> gives Marx's master equation:
        </p>
        <p className="v3-ch03-equation">p′ = s′ · (v / C) = s′ · v / (c + v)</p>
        <p>
          The rate of profit resolves into two independent factors — the{" "}
          <strong>rate of surplus-value s′</strong> and the{" "}
          <strong>value-composition v/C</strong>. The chapter is the systematic
          variation of the two.
        </p>
      </section>

      <section className="v3-ch03-slider" aria-label="Live profit-rate formula">
        <h4>p′ = s′ · (v / C)</h4>
        <div className="v3-ch03-controls">
          <label className="v3-ch03-control">
            <span>Constant capital c — {money(c)}</span>
            <input type="range" min={0} max={400} step={10} value={c}
              onChange={(e) => setC(Number(e.target.value))} />
          </label>
          <label className="v3-ch03-control">
            <span>Variable capital v — {money(v)}</span>
            <input type="range" min={0} max={200} step={5} value={v}
              onChange={(e) => setV(Number(e.target.value))} />
          </label>
          <label className="v3-ch03-control">
            <span>Rate of surplus-value s′ — {sRatePct}%</span>
            <input type="range" min={0} max={300} step={5} value={sRatePct}
              onChange={(e) => setSRatePct(Number(e.target.value))} />
          </label>
        </div>
        <div className="v3-ch03-readout">
          <div className="v3-ch03-stat">
            <span className="v3-ch03-stat-label">Total capital C = c + v</span>
            <span className="v3-ch03-stat-value">{money(total)}</span>
          </div>
          <div className="v3-ch03-stat">
            <span className="v3-ch03-stat-label">Value-composition v / C</span>
            <span className="v3-ch03-stat-value">{pct(compBp)}</span>
          </div>
          <div className="v3-ch03-stat v3-ch03-stat-headline">
            <span className="v3-ch03-stat-label">Rate of profit p′</span>
            <span className="v3-ch03-stat-value">{pct(pPrimeBp)}</span>
          </div>
        </div>
        <p className="v3-ch03-slider-note">
          p′ is the rate of surplus-value <strong>{pct(sRateBp)}</strong> scaled
          by the variable share of capital <strong>{pct(compBp)}</strong>.
        </p>
      </section>

      <section className="v3-ch03-variation" aria-label="Variation analyser">
        <h4>How a change shifts p′</h4>
        <div className="v3-ch03-presets">
          {VARIATION_PRESETS.map((p) => (
            <button key={p.label} type="button" className="v3-ch03-preset" onClick={() => analyse(p)}>
              {p.label}
            </button>
          ))}
        </div>
        {variation && (
          <div className="v3-ch03-variation-result">
            <p className="v3-ch03-case">{CASE_LABEL[variation.case]}</p>
            <div className="v3-ch03-shift">
              <span className="v3-ch03-shift-old">{pct(variation.old_profit_rate)}</span>
              <span className="v3-ch03-shift-arrow">→</span>
              <span className={`v3-ch03-shift-new v3-ch03-${movement}`}>
                {pct(variation.new_profit_rate)}
              </span>
              <span className="v3-ch03-shift-word">p′ {movement}</span>
            </div>
            <p className="v3-ch03-shift-detail">
              Value-composition v/C: {pct(variation.old_composition_bp)} →{" "}
              {pct(variation.new_composition_bp)}.
              {variation.proportional_change && (
                <strong> p′ moved in exact proportion to s′.</strong>
              )}
            </p>
          </div>
        )}
      </section>

      <section className="v3-ch03-compare" aria-label="Two-capital comparison">
        <h4>Two capitals side by side</h4>
        <div className="v3-ch03-presets">
          {COMPARE_PRESETS.map((p) => (
            <button key={p.label} type="button" className="v3-ch03-preset" onClick={() => compare(p)}>
              {p.label}
            </button>
          ))}
        </div>
        {comparison && (
          <div className="v3-ch03-compare-result">
            <div className="v3-ch03-compare-cards">
              <div className="v3-ch03-compare-card">
                <span className="v3-ch03-compare-name">Capital 1</span>
                <span className="v3-ch03-compare-rate">{pct(comparison.first_profit_rate)}</span>
                <span className="v3-ch03-compare-frac">v/C = {pct(comparison.first_composition.basis_points)}</span>
              </div>
              <div className="v3-ch03-compare-card">
                <span className="v3-ch03-compare-name">Capital 2</span>
                <span className="v3-ch03-compare-rate">{pct(comparison.second_profit_rate)}</span>
                <span className="v3-ch03-compare-frac">v/C = {pct(comparison.second_composition.basis_points)}</span>
              </div>
            </div>
            <p className="v3-ch03-compare-note">
              {comparison.equal_surplus_rate
                ? "With s′ equal, the rates of profit stand to each other as the value-compositions v/C."
                : "With s′ unequal, the difference in rates reflects the difference in s′ as well as v/C."}
            </p>
          </div>
        )}
      </section>

      {error && <p className="v3-ch03-error">{error}</p>}

      {saved.length > 0 && (
        <section className="v3-ch03-saved" aria-label="Recorded variations">
          <h4>Recorded variation analyses</h4>
          <table className="v3-ch03-table">
            <thead>
              <tr>
                <th>case</th>
                <th>initial (c+v, s′)</th>
                <th>changed (c+v, s′)</th>
                <th>p′ old → new</th>
                <th>v/C old → new</th>
                <th>∝ s′</th>
              </tr>
            </thead>
            <tbody>
              {saved.map((a) => (
                <tr key={a.id}>
                  <td>{a.case}</td>
                  <td>{money(a.initial.c + a.initial.v)}, {pct(a.initial.s_rate)}</td>
                  <td>{money(a.changed.c + a.changed.v)}, {pct(a.changed.s_rate)}</td>
                  <td>{pct(a.old_profit_rate)} → {pct(a.new_profit_rate)}</td>
                  <td>{pct(a.old_composition_bp)} → {pct(a.new_composition_bp)}</td>
                  <td>{a.proportional_change ? "yes" : "—"}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </section>
      )}
    </div>
  );
}
