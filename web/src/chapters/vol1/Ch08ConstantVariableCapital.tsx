import { useState } from "react";
import type { FormEvent } from "react";
import { api } from "../../api";
import type {
  ConstantCapitalInput,
  DecomposeCapitalResult,
  CapitalCompositionResult,
} from "../../types";
import "./Ch08ConstantVariableCapital.css";

interface Ch08Props {
  onSharedChanged: () => void;
}

function minutesToHours(m: number): string {
  const h = Math.floor(m / 60);
  const min = m % 60;
  return min === 0 ? `${h}h` : `${h}h ${min}m`;
}

const COMPARE_INDUSTRIES = [
  {
    id: "spinning",
    label: "1871 Spinning Mill",
    tag: "High c/v",
    c: 1440,
    v: 360,
    s: 360,
  },
  {
    id: "tailoring",
    label: "Tailoring (handicraft)",
    tag: "Low c/v",
    c: 120,
    v: 360,
    s: 360,
  },
] as const;

export function Ch08ConstantVariableCapital({ onSharedChanged: _onSharedChanged }: Ch08Props) {
  return (
    <>
      <CapitalInsight />
      <DecomposeCapitalPanel />
      <CapitalCompositionPanel />
      <IndustryComparisonPanel />
      <Coda />
    </>
  );
}

function CapitalInsight() {
  return (
    <section className="v1-ch08-insight">
      <h2 className="v1-ch08-insight-h2">Constant transfers. Variable expands.</h2>
      <p className="v1-ch08-insight-prose">
        Constant capital transfers its value to the product unchanged; variable
        capital reproduces its own value <em>and</em> expands it (v → v + s).
        The composition <code>c/v</code> measures how much dead labour each
        living labourer animates.
      </p>
    </section>
  );
}

function CvsBar({
  c,
  v,
  s,
  label,
  total,
}: {
  c: number;
  v: number;
  s: number;
  label: string;
  total?: number;
}) {
  const sum = total ?? c + v + s;
  const pct = (n: number) => (sum > 0 ? `${(n / sum) * 100}%` : "0%");
  return (
    <div className="v1-ch08-decomp" role="img" aria-label={`c + v + s decomposition: ${label}`}>
      <span className="v1-ch08-decomp-label">{label}</span>
      <div className="v1-ch08-decomp-track">
        <div
          className="v1-ch08-decomp-seg v1-ch08-decomp-seg--c"
          style={{ flexBasis: pct(c) }}
          title={`c: ${minutesToHours(c)}`}
        >
          c
        </div>
        <div
          className="v1-ch08-decomp-seg v1-ch08-decomp-seg--v"
          style={{ flexBasis: pct(v) }}
          title={`v: ${minutesToHours(v)}`}
        >
          v
        </div>
        <div
          className="v1-ch08-decomp-seg v1-ch08-decomp-seg--s"
          style={{ flexBasis: pct(s) }}
          title={`s: ${minutesToHours(s)}`}
        >
          s
        </div>
      </div>
      <span className="v1-ch08-decomp-total">{minutesToHours(sum)}</span>
    </div>
  );
}

function CvsLegend({ c, v, s }: { c: number; v: number; s: number }) {
  return (
    <p className="v1-ch08-legend">
      <span>
        <span className="v1-ch08-legend-swatch" style={{ background: "var(--lead-hover)" }} />
        c — constant ({minutesToHours(c)})
      </span>
      <span>
        <span className="v1-ch08-legend-swatch" style={{ background: "var(--lead)" }} />
        v — variable ({minutesToHours(v)})
      </span>
      <span>
        <span className="v1-ch08-legend-swatch" style={{ background: "var(--gold-bright)" }} />
        s — surplus ({minutesToHours(s)})
      </span>
    </p>
  );
}

function CompositionKpi({ c, v }: { c: number; v: number }) {
  const ratio = v > 0 ? (c / v).toFixed(2) : "∞";
  return (
    <div className="v1-ch08-composition">
      <span className="v1-ch08-composition-label">c / v</span>
      <span className="v1-ch08-composition-value">{ratio}</span>
      <span className="v1-ch08-composition-ratio">
        ({minutesToHours(c)} dead / {minutesToHours(v)} living)
      </span>
    </div>
  );
}

function DecomposeCapitalPanel() {
  const [constants, setConstants] = useState<ConstantCapitalInput[]>([
    { original_value: 1200, kind: "raw_material", service_life_days: 0 },
    { original_value: 240, kind: "instrument", service_life_days: 1 },
  ]);
  const [wageValue, setWageValue] = useState(360);
  const [workingDay, setWorkingDay] = useState(720);
  const [result, setResult] = useState<DecomposeCapitalResult | null>(null);
  const [err, setErr] = useState<string | null>(null);

  async function submit(e: FormEvent) {
    e.preventDefault();
    setErr(null);
    try {
      const r = await api.decomposeCapital({
        constant_capitals: constants,
        variable_capital: { wage_value: wageValue, working_day: workingDay },
      });
      setResult(r);
    } catch (e) {
      setErr(e instanceof Error ? e.message : String(e));
    }
  }

  function updateConstant(i: number, field: keyof ConstantCapitalInput, value: string | number) {
    setConstants((prev) =>
      prev.map((c, idx) => (idx === i ? { ...c, [field]: value } : c))
    );
  }

  return (
    <section className="card">
      <h2>Decompose Product Value</h2>
      <p className="description">
        Enter means of production (constant capital) and variable capital to decompose the
        product value into c + v + s — the transferred value, reproduced value, and surplus.
      </p>
      <form className="form-grid" onSubmit={submit}>
        <div className="span2">
          <fieldset>
            <legend>Constant Capital (Means of Production)</legend>
            {constants.map((c, i) => (
              <div key={i} className="item-row">
                <select
                  value={c.kind}
                  onChange={(e) => updateConstant(i, "kind", e.target.value)}
                >
                  <option value="raw_material">Raw Material</option>
                  <option value="instrument">Instrument</option>
                  <option value="auxiliary">Auxiliary</option>
                </select>
                <input
                  type="number"
                  placeholder="original value (min)"
                  min={0}
                  value={c.original_value}
                  onChange={(e) => updateConstant(i, "original_value", Number(e.target.value))}
                />
                <input
                  type="number"
                  placeholder="service life (days, 0=consumed)"
                  min={0}
                  value={c.service_life_days}
                  onChange={(e) => updateConstant(i, "service_life_days", Number(e.target.value))}
                />
                <button
                  type="button"
                  onClick={() => setConstants((prev) => prev.filter((_, idx) => idx !== i))}
                >
                  ×
                </button>
              </div>
            ))}
            <button
              type="button"
              onClick={() =>
                setConstants((prev) => [
                  ...prev,
                  { original_value: 0, kind: "raw_material", service_life_days: 0 },
                ])
              }
            >
              + Add Constant Capital
            </button>
          </fieldset>
        </div>

        <label>
          <span>Wage Value (necessary labour, min)</span>
          <input
            type="number"
            min={0}
            value={wageValue}
            onChange={(e) => setWageValue(Number(e.target.value))}
          />
        </label>
        <label>
          <span>Working Day (total duration, min)</span>
          <input
            type="number"
            min={1}
            value={workingDay}
            onChange={(e) => setWorkingDay(Number(e.target.value))}
          />
        </label>

        <div className="form-actions span2">
          <button type="submit" className="primary">Decompose</button>
          {err && <span className="error">{err}</span>}
        </div>
      </form>

      {result && (
        <>
          <h3 className="v1-ch08-result-h3">c + v + s</h3>
          <CvsBar
            c={result.product_value.constant}
            v={result.product_value.variable}
            s={result.product_value.surplus}
            label="Product value"
          />
          <CvsLegend
            c={result.product_value.constant}
            v={result.product_value.variable}
            s={result.product_value.surplus}
          />
          <CompositionKpi
            c={result.product_value.constant}
            v={result.product_value.variable}
          />
        </>
      )}
    </section>
  );
}

function CapitalCompositionPanel() {
  const [constantValue, setConstantValue] = useState(410);
  const [variableValue, setVariableValue] = useState(90);
  const [result, setResult] = useState<CapitalCompositionResult | null>(null);
  const [err, setErr] = useState<string | null>(null);

  async function submit(e: FormEvent) {
    e.preventDefault();
    setErr(null);
    try {
      const r = await api.getCapitalComposition(constantValue, variableValue);
      setResult(r);
    } catch (e) {
      setErr(e instanceof Error ? e.message : String(e));
    }
  }

  return (
    <section className="card">
      <h2>Capital Composition</h2>
      <p className="description">
        Compute c/v — the technical and value composition of capital, expressing the ratio of
        dead labour (means of production) to living labour (labour-power).
      </p>
      <form className="form-grid" onSubmit={submit}>
        <label>
          <span>Constant Capital c (min)</span>
          <input
            type="number"
            min={0}
            value={constantValue}
            onChange={(e) => setConstantValue(Number(e.target.value))}
          />
        </label>
        <label>
          <span>Variable Capital v (min)</span>
          <input
            type="number"
            min={0}
            value={variableValue}
            onChange={(e) => setVariableValue(Number(e.target.value))}
          />
        </label>
        <div className="form-actions span2">
          <button type="submit" className="primary">Compute</button>
          {err && <span className="error">{err}</span>}
        </div>
      </form>

      {result && (
        <>
          <h3 className="v1-ch08-result-h3">Composition</h3>
          <CvsBar c={result.constant} v={result.variable} s={0} label="c + v" />
          <CompositionKpi c={result.constant} v={result.variable} />
        </>
      )}
    </section>
  );
}

function IndustryComparisonPanel() {
  const max = Math.max(
    ...COMPARE_INDUSTRIES.map((it) => it.c + it.v + it.s),
  );
  return (
    <section className="card">
      <h2>High vs Low Organic Composition</h2>
      <p className="description">
        Two industries with the same wage bill (v) and the same surplus rate
        (s/v = 100%) produce strikingly different aggregates — the spinning
        mill animates 12× as much dead labour as the tailor.
      </p>
      <div className="v1-ch08-compare">
        {COMPARE_INDUSTRIES.map((it) => (
          <div key={it.id} className="v1-ch08-compare-card">
            <span className="v1-ch08-compare-tag">{it.tag}</span>
            <h3 className="v1-ch08-compare-h3">{it.label}</h3>
            <CvsBar
              c={it.c}
              v={it.v}
              s={it.s}
              label="c + v + s"
              total={max}
            />
            <span className="v1-ch08-compare-meta">
              c / v = <strong>{(it.c / it.v).toFixed(2)}</strong>{" "}
              <span className="v1-ch08-compare-meta-muted">
                · s / v = {(it.s / it.v).toFixed(2)}
              </span>
            </span>
          </div>
        ))}
      </div>
    </section>
  );
}

function Coda() {
  return (
    <aside className="v1-ch08-coda">
      <p className="v1-ch08-coda-quote">
        “That part of capital which is converted into means of production …
        does not, in the process of production, undergo any quantitative
        alteration of value. I therefore call it the <em>constant</em> part of
        capital.”
        <span className="v1-ch08-coda-cite">
          — Marx, Capital Vol. I, Ch. 8 §1
        </span>
      </p>
    </aside>
  );
}
