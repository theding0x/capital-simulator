import { useState } from "react";
import type { FormEvent } from "react";
import { api } from "../../api";
import type { RatesOfSurplusValueResult } from "../../types";

type Preset = {
  id: string;
  label: string;
  necessary: number;
  surplus: number;
  caption: string;
};

const PRESETS: Preset[] = [
  {
    id: "p1",
    label: "§ Formula I — 6h / 6h = 100%",
    necessary: 360,
    surplus: 360,
    caption:
      "6 hours surplus-labour / 6 hours necessary labour = 100%. Formula I is unbounded above 1.0.",
  },
  {
    id: "p2",
    label: "§ Formula II bounded",
    necessary: 360,
    surplus: 359,
    caption:
      "Surplus = 359 min, necessary = 360 min: Formula II = 49.9%, well below 100% even when surplus nearly equals necessary.",
  },
  {
    id: "p3",
    label: "§ Agricultural labourer — 300%",
    necessary: 180,
    surplus: 540,
    caption:
      "\"The labourer gets only 1/4, the capitalist ... 3/4\": NL = 3 h, SL = 9 h. Rate = 300%.",
  },
];

function pct(r: number) {
  return (r * 100).toFixed(2) + "%";
}

export function Ch18RatesOfSurplusValue() {
  const [necessary, setNecessary] = useState(360);
  const [surplus, setSurplus] = useState(360);
  const [result, setResult] = useState<RatesOfSurplusValueResult | null>(null);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);
  const [caption, setCaption] = useState("");

  function applyPreset(p: Preset) {
    setNecessary(p.necessary);
    setSurplus(p.surplus);
    setResult(null);
    setError("");
    setCaption(p.caption);
  }

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setError("");
    setResult(null);
    setLoading(true);
    try {
      const r = await api.computeRatesOfSurplusValue({
        necessary_labour_minutes: necessary,
        surplus_labour_minutes: surplus,
      });
      setResult(r);
    } catch (err) {
      setError(err instanceof Error ? err.message : String(err));
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="ch-panel">
      <section className="ch-section">
        <h2 className="ch-section-title">Three Formulae</h2>
        <p className="ch-description">
          Marx shows that the rate of surplus-value admits three equivalent
          expressions. Formula&nbsp;I (s/v) is unbounded; Formula&nbsp;II
          (s/(s+v)) is always below 100%; Formula&nbsp;III (unpaid/paid) is
          identical to Formula&nbsp;I in magnitude.
        </p>

        <div className="preset-row">
          {PRESETS.map((p) => (
            <button
              key={p.id}
              type="button"
              className="preset-btn"
              onClick={() => applyPreset(p)}
            >
              {p.label}
            </button>
          ))}
        </div>

        {caption && <p className="ch-caption italic">{caption}</p>}

        <form className="ch-form" onSubmit={handleSubmit}>
          <div className="form-row">
            <label className="form-label">
              Necessary labour (min)
              <input
                type="number"
                min={1}
                className="form-input"
                value={necessary}
                onChange={(e) => setNecessary(Number(e.target.value))}
                required
              />
            </label>
            <label className="form-label">
              Surplus labour (min)
              <input
                type="number"
                min={0}
                className="form-input"
                value={surplus}
                onChange={(e) => setSurplus(Number(e.target.value))}
                required
              />
            </label>
          </div>
          <button type="submit" className="ch-btn" disabled={loading}>
            {loading ? "Computing..." : "Compute rates"}
          </button>
        </form>

        {error && <p className="ch-error">{error}</p>}

        {result && (
          <div className="rates-grid">
            <div className="rates-card">
              <div className="rates-label">Formula I &mdash; s/v</div>
              <div className="rates-value">{pct(result.formula_i)}</div>
              <div className="rates-note">
                Surplus labour / necessary labour. Unbounded above 100%.
              </div>
            </div>
            <div className="rates-card">
              <div className="rates-label">Formula II &mdash; s/(s+v)</div>
              <div className="rates-value">{pct(result.formula_ii)}</div>
              <div className="rates-note">
                Surplus labour / working day. Always strictly below 100%.
              </div>
            </div>
            <div className="rates-card">
              <div className="rates-label">Formula III &mdash; unpaid/paid</div>
              <div className="rates-value">{pct(result.formula_iii)}</div>
              <div className="rates-note">
                Unpaid labour / paid labour. Identical to Formula I.
              </div>
            </div>
            <div className="rates-summary">
              <table className="ch-table">
                <tbody>
                  <tr>
                    <td>Necessary labour</td>
                    <td>{result.necessary_labour_minutes} min</td>
                  </tr>
                  <tr>
                    <td>Surplus labour</td>
                    <td>{result.surplus_labour_minutes} min</td>
                  </tr>
                  <tr>
                    <td>Working day</td>
                    <td>{result.working_day_minutes} min</td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        )}
      </section>
    </div>
  );
}
