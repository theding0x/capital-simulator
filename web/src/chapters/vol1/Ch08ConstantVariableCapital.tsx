import { useState } from "react";
import type { FormEvent } from "react";
import { api } from "../../api";
import type {
  ConstantCapitalInput,
  DecomposeCapitalResult,
  CapitalCompositionResult,
} from "../../types";

interface Ch08Props {
  onSharedChanged: () => void;
}

function minutesToHours(m: number): string {
  const h = Math.floor(m / 60);
  const min = m % 60;
  return min === 0 ? `${h}h` : `${h}h ${min}m`;
}

export function Ch08ConstantVariableCapital({ onSharedChanged: _onSharedChanged }: Ch08Props) {
  return (
    <>
      <DecomposeCapitalPanel />
      <CapitalCompositionPanel />
    </>
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
        <table className="data-table">
          <tbody>
            <tr>
              <td>Constant Capital (c) — transferred</td>
              <td>{minutesToHours(result.product_value.constant)}</td>
            </tr>
            <tr>
              <td>Variable Capital (v) — reproduced</td>
              <td>{minutesToHours(result.product_value.variable)}</td>
            </tr>
            <tr>
              <td>Surplus Value (s)</td>
              <td><strong>{minutesToHours(result.product_value.surplus)}</strong></td>
            </tr>
            <tr>
              <td>Total Product Value (c + v + s)</td>
              <td>
                <strong>
                  {minutesToHours(
                    result.product_value.constant +
                    result.product_value.variable +
                    result.product_value.surplus
                  )}
                </strong>
              </td>
            </tr>
            <tr>
              <td>Capital Composition (c/v)</td>
              <td>{result.composition_ratio.toFixed(2)}</td>
            </tr>
          </tbody>
        </table>
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
        <table className="data-table">
          <tbody>
            <tr>
              <td>Constant Capital (c)</td>
              <td>{minutesToHours(result.constant)}</td>
            </tr>
            <tr>
              <td>Variable Capital (v)</td>
              <td>{minutesToHours(result.variable)}</td>
            </tr>
            <tr>
              <td>Composition Ratio (c/v)</td>
              <td><strong>{result.composition_ratio.toFixed(2)}</strong></td>
            </tr>
          </tbody>
        </table>
      )}
    </section>
  );
}
