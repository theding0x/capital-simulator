import { useEffect, useState } from "react";
import { financeApi } from "../../api";
import type { ProfitRateAnalysisResponse } from "../../types";
import "./Ch02RateOfProfit.css";

// Rates arrive from the API in basis points (10000 bp = 100%); render as percent.
function pct(bp: number): string {
  return `${(bp / 100).toLocaleString("en-GB", { maximumFractionDigits: 2 })}%`;
}

// Capital magnitudes are LabourMinutes; Marx denominates them in £ (6s. = one
// social working-day), so we render them as whole pounds.
function money(n: number): string {
  return `£${n.toLocaleString("en-GB")}`;
}

interface DecompForm {
  constant: number;
  variable: number;
  surplus: number;
}

interface Preset {
  label: string;
  values: DecompForm;
}

// Capital III, Ch. 2 — Marx's worked figures and his limiting cases.
const PRESETS: Preset[] = [
  {
    label: "§ 400c + 100v + 100s",
    values: { constant: 400, variable: 100, surplus: 100 },
  },
  {
    label: "v = C (no constant capital)",
    values: { constant: 0, variable: 100, surplus: 100 },
  },
  {
    label: "High composition 900c + 100v",
    values: { constant: 900, variable: 100, surplus: 100 },
  },
];

export function Ch02RateOfProfit() {
  const [form, setForm] = useState<DecompForm>(PRESETS[0].values);
  const [analysis, setAnalysis] = useState<ProfitRateAnalysisResponse | null>(null);
  const [saved, setSaved] = useState<ProfitRateAnalysisResponse[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [busy, setBusy] = useState(false);

  async function refreshSaved() {
    try {
      setSaved(await financeApi.listProfitRates());
    } catch {
      // the seed list is best-effort; the form still works without it
    }
  }

  useEffect(() => {
    refreshSaved();
  }, []);

  async function compute() {
    setError(null);
    setBusy(true);
    try {
      const a = await financeApi.computeProfitRate({
        constant: form.constant,
        variable: form.variable,
        surplus_value: form.surplus,
      });
      setAnalysis(a);
      await refreshSaved();
    } catch (e) {
      setError(String(e));
    } finally {
      setBusy(false);
    }
  }

  // The rate of profit occupies this fraction of the rate of surplus-value; the
  // remainder is what the profit-form conceals (the mystification).
  const profitShare =
    analysis && analysis.surplus_value_rate > 0
      ? (analysis.profit_rate.basis_points / analysis.surplus_value_rate) * 100
      : 0;
  const surplus = analysis?.profit_rate.surplus_value ?? 0;
  const totalCapital = analysis?.profit_rate.total_capital ?? 0;

  return (
    <div className="v3-ch02">
      <section className="v3-ch02-intro">
        <p>
          The same surplus-value <strong>s</strong> can be measured two ways.
          Against the variable capital that alone produced it, it gives the{" "}
          <strong>rate of surplus-value</strong> s′ = s/v — the degree of
          exploitation. Against the <em>total</em> capital advanced, it gives the{" "}
          <strong>rate of profit</strong> p′ = s/C = s/(c+v). Because{" "}
          <strong>C ≥ v</strong>, the rate of profit is always the smaller of the
          two, and the gap between them measures the profit-form's{" "}
          <em>mystification</em> of where the surplus comes from.
        </p>
      </section>

      <section className="v3-ch02-form" aria-label="Capital decomposition">
        <div className="v3-ch02-presets">
          {PRESETS.map((p) => (
            <button
              key={p.label}
              type="button"
              className="v3-ch02-preset"
              onClick={() => setForm(p.values)}
            >
              {p.label}
            </button>
          ))}
        </div>

        <div className="v3-ch02-fields">
          <label className="v3-ch02-field">
            Constant capital c
            <input
              type="number"
              min={0}
              value={form.constant}
              onChange={(e) => setForm({ ...form, constant: Number(e.target.value) })}
            />
          </label>
          <label className="v3-ch02-field">
            Variable capital v
            <input
              type="number"
              min={0}
              value={form.variable}
              onChange={(e) => setForm({ ...form, variable: Number(e.target.value) })}
            />
          </label>
          <label className="v3-ch02-field">
            Surplus-value s
            <input
              type="number"
              min={0}
              value={form.surplus}
              onChange={(e) => setForm({ ...form, surplus: Number(e.target.value) })}
            />
          </label>
        </div>

        <button type="button" className="v3-ch02-compute" onClick={compute} disabled={busy}>
          {busy ? "Computing…" : "Compute the two rates"}
        </button>
        {error && <p className="v3-ch02-error">{error}</p>}
      </section>

      {analysis && (
        <section className="v3-ch02-result" aria-label="Rate comparison">
          <div className="v3-ch02-rates">
            <div className="v3-ch02-rate v3-ch02-rate-surplus">
              <span className="v3-ch02-rate-label">Rate of surplus-value s′ = s / v</span>
              <span className="v3-ch02-rate-value">{pct(analysis.surplus_value_rate)}</span>
              <span className="v3-ch02-rate-frac">
                {money(surplus)} / {money(analysis.variable_capital)}
              </span>
            </div>
            <div className="v3-ch02-rate v3-ch02-rate-profit">
              <span className="v3-ch02-rate-label">Rate of profit p′ = s / C</span>
              <span className="v3-ch02-rate-value">{pct(analysis.profit_rate.basis_points)}</span>
              <span className="v3-ch02-rate-frac">
                {money(surplus)} / {money(totalCapital)}
              </span>
            </div>
          </div>

          <p className="v3-ch02-explain">
            The total capital is{" "}
            <strong>
              C = c + v = {money(analysis.constant_capital)} +{" "}
              {money(analysis.variable_capital)} = {money(totalCapital)}
            </strong>
            . The surplus is divided by the larger denominator, so p′ falls below s′
            {analysis.constant_capital === 0 ? (
              <> — except here, where c = 0 and the two rates coincide.</>
            ) : (
              <>.</>
            )}
          </p>

          <div className="v3-ch02-myst" aria-label="Mystification">
            <div className="v3-ch02-myst-track">
              <div
                className="v3-ch02-myst-profit"
                style={{ width: `${profitShare}%` }}
                title={`Rate of profit: ${pct(analysis.profit_rate.basis_points)}`}
              >
                p′
              </div>
              <div
                className="v3-ch02-myst-hidden"
                style={{ width: `${100 - profitShare}%` }}
                title={`Concealed: ${pct(analysis.mystification)}`}
              >
                concealed
              </div>
            </div>
            <p className="v3-ch02-myst-caption">
              Of the rate of surplus-value, the profit-form shows only the rate of
              profit and conceals <strong>{pct(analysis.mystification)}</strong> —
              the mystification degree (s′ − p′) / s′.
            </p>
          </div>
        </section>
      )}

      {saved.length > 0 && (
        <section className="v3-ch02-saved" aria-label="Recorded analyses">
          <h4>Recorded rate-of-profit analyses</h4>
          <table className="v3-ch02-table">
            <thead>
              <tr>
                <th>c</th>
                <th>v</th>
                <th>s</th>
                <th>C</th>
                <th>p′</th>
                <th>s′</th>
                <th>concealed</th>
              </tr>
            </thead>
            <tbody>
              {saved.map((a) => (
                <tr key={a.id}>
                  <td>{money(a.constant_capital)}</td>
                  <td>{money(a.variable_capital)}</td>
                  <td>{money(a.profit_rate.surplus_value)}</td>
                  <td>{money(a.profit_rate.total_capital)}</td>
                  <td>{pct(a.profit_rate.basis_points)}</td>
                  <td>{pct(a.surplus_value_rate)}</td>
                  <td>{pct(a.mystification)}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </section>
      )}
    </div>
  );
}
