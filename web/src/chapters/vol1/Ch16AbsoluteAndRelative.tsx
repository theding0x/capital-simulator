import { useCallback, useEffect, useState } from "react";
import type { FormEvent } from "react";
import { api } from "../../api";
import { fmtHoursLong } from "../../format";
import type {
  AbsoluteSurplusValueResult,
  RelativeSurplusValueResult,
  SurplusValueRateResult,
} from "../../types";

// Ch.16 is the synthesis chapter: the same working-day partition drives
// both branches. The panel keeps one set of inputs (working day +
// necessary labour) and computes both surplus-value forms side-by-side.
export function Ch16AbsoluteAndRelative() {
  return (
    <>
      <SurplusValueCalculator />
      <MillCritiquePanel />
    </>
  );
}

function SurplusValueCalculator() {
  const [workingDay, setWorkingDay] = useState(720);
  const [necessary, setNecessary] = useState(600);
  const [extension, setExtension] = useState(60);
  const [productivity, setProductivity] = useState(2.0);
  const [abs, setAbs] = useState<AbsoluteSurplusValueResult | null>(null);
  const [rel, setRel] = useState<RelativeSurplusValueResult | null>(null);
  const [err, setErr] = useState<string | null>(null);

  const recompute = useCallback(async () => {
    setErr(null);
    try {
      const [a, r] = await Promise.all([
        api.computeAbsoluteSurplusValue({
          working_day_minutes: workingDay,
          necessary_labour_minutes: necessary,
          extension_minutes: extension,
        }),
        api.computeRelativeSurplusValue({
          working_day_minutes: workingDay,
          necessary_labour_minutes: necessary,
          productivity_factor: productivity,
        }),
      ]);
      setAbs(a);
      setRel(r);
    } catch (e) {
      setErr(e instanceof Error ? e.message : String(e));
    }
  }, [workingDay, necessary, extension, productivity]);

  useEffect(() => {
    void recompute();
  }, [recompute]);

  function submit(e: FormEvent) {
    e.preventDefault();
    void recompute();
  }

  return (
    <section className="card">
      <h2>Working-day partition (§1)</h2>
      <p className="description">
        Ch. 16 §1: surplus-value comes in two analytically distinct forms
        — absolute (prolong the working day, hold necessary labour fixed)
        and relative (raise productivity, hold the working day fixed) —
        but the magnitude type is the same. The Origin tag is the only
        thing that distinguishes them. Configure a starting partition and
        watch both mechanisms produce surplus from the same input.
      </p>
      <form className="form-grid" onSubmit={submit}>
        <label>
          <span>Working day (min)</span>
          <input
            type="number"
            min={60}
            max={1440}
            value={workingDay}
            onChange={(e) => setWorkingDay(Number(e.target.value))}
          />
        </label>
        <label>
          <span>Necessary labour (min)</span>
          <input
            type="number"
            min={1}
            max={workingDay - 1}
            value={necessary}
            onChange={(e) => setNecessary(Number(e.target.value))}
          />
        </label>
        <label>
          <span>Extension (min) — absolute path</span>
          <input
            type="number"
            min={1}
            value={extension}
            onChange={(e) => setExtension(Number(e.target.value))}
          />
        </label>
        <label>
          <span>Productivity factor — relative path</span>
          <input
            type="number"
            min={1.01}
            step={0.1}
            value={productivity}
            onChange={(e) => setProductivity(Number(e.target.value))}
          />
        </label>
      </form>
      {err && <p className="error">{err}</p>}
      <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: "1rem", marginTop: "1rem" }}>
        <ResultCard
          heading="Absolute — prolongation"
          subtitle="Working day grows; necessary labour fixed."
          oldWD={abs?.old_working_day}
          newWD={abs?.new_working_day}
          surplus={abs?.absolute_surplus_value.labour_minutes ?? 0}
          origin={abs?.absolute_surplus_value.origin ?? "absolute"}
        />
        <ResultCard
          heading="Relative — productivity"
          subtitle="Necessary labour shrinks; working day fixed."
          oldWD={rel?.old_working_day}
          newWD={rel?.new_working_day}
          surplus={rel?.relative_surplus_value.labour_minutes ?? 0}
          origin={rel?.relative_surplus_value.origin ?? "relative"}
        />
      </div>
      <p className="small muted" style={{ marginTop: "1rem" }}>
        §1, closing: "From one standpoint, any distinction between
        absolute and relative surplus-value appears illusory. Relative
        surplus-value is absolute … Absolute surplus-value is relative."
        The two share a magnitude type; the Origin tag is the line
        between them.
      </p>
    </section>
  );
}

interface ProductionWorkingDay {
  total: number;
  necessary_labour: number;
  surplus_labour: number;
}

function ResultCard({
  heading,
  subtitle,
  oldWD,
  newWD,
  surplus,
  origin,
}: {
  heading: string;
  subtitle: string;
  oldWD?: ProductionWorkingDay;
  newWD?: ProductionWorkingDay;
  surplus: number;
  origin: string;
}) {
  return (
    <div style={{ padding: "0.75rem", border: "1px solid var(--ink-muted)" }}>
      <h3 style={{ marginTop: 0, marginBottom: "0.25rem" }}>{heading}</h3>
      <p className="small muted" style={{ marginTop: 0 }}>{subtitle}</p>
      {oldWD && newWD ? (
        <table className="data-table">
          <tbody>
            <tr>
              <td>Total</td>
              <td>{fmtHoursLong(oldWD.total)} → <strong>{fmtHoursLong(newWD.total)}</strong></td>
            </tr>
            <tr>
              <td>Necessary</td>
              <td>{fmtHoursLong(oldWD.necessary_labour)} → <strong>{fmtHoursLong(newWD.necessary_labour)}</strong></td>
            </tr>
            <tr>
              <td>Surplus</td>
              <td>{fmtHoursLong(oldWD.surplus_labour)} → <strong>{fmtHoursLong(newWD.surplus_labour)}</strong></td>
            </tr>
            <tr>
              <td>Gain</td>
              <td><strong>{surplus} min</strong> <span className="muted">({origin})</span></td>
            </tr>
          </tbody>
        </table>
      ) : (
        <p className="small muted">—</p>
      )}
    </div>
  );
}

function MillCritiquePanel() {
  // §1 fixture: £400 constant + £100 variable + £20 surplus.
  const [surplusLabour, setSurplusLabour] = useState(20);
  const [necessaryLabour, setNecessaryLabour] = useState(100);
  const [totalCapital, setTotalCapital] = useState(500);
  const [result, setResult] = useState<SurplusValueRateResult | null>(null);
  const [err, setErr] = useState<string | null>(null);

  const recompute = useCallback(async () => {
    setErr(null);
    try {
      const r = await api.getSurplusValueRate({
        surplus_labour: surplusLabour,
        necessary_labour: necessaryLabour,
        total_capital: totalCapital,
      });
      setResult(r);
    } catch (e) {
      setErr(e instanceof Error ? e.message : String(e));
    }
  }, [surplusLabour, necessaryLabour, totalCapital]);

  useEffect(() => {
    void recompute();
  }, [recompute]);

  return (
    <section className="card">
      <h2>Rate of surplus-value vs rate of profit (§1, Mill critique)</h2>
      <p className="description">
        Marx, §1: "if the rate of surplus-value be 20%, the rate of profit
        will be 20:500, i.e., 4% and not 20%." The two rates use the same
        numerator (surplus-value) but different denominators — necessary
        labour vs total capital advanced. The Mill error is to collapse
        them.
      </p>
      <form className="form-grid" onSubmit={(e) => { e.preventDefault(); void recompute(); }}>
        <label>
          <span>Surplus labour (= surplus-value)</span>
          <input type="number" min={1} value={surplusLabour} onChange={(e) => setSurplusLabour(Number(e.target.value))} />
        </label>
        <label>
          <span>Necessary labour (= variable capital)</span>
          <input type="number" min={1} value={necessaryLabour} onChange={(e) => setNecessaryLabour(Number(e.target.value))} />
        </label>
        <label>
          <span>Total capital advanced (c + v)</span>
          <input type="number" min={1} value={totalCapital} onChange={(e) => setTotalCapital(Number(e.target.value))} />
        </label>
      </form>
      {err && <p className="error">{err}</p>}
      {result && (
        <table className="data-table" style={{ marginTop: "0.75rem" }}>
          <thead>
            <tr>
              <th>Quantity</th>
              <th>Formula</th>
              <th>Value</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td>Rate of surplus-value (s/v)</td>
              <td>{result.surplus_labour_minutes} / {result.necessary_labour_minutes}</td>
              <td><strong>{(result.rate_of_surplus_value * 100).toFixed(2)}%</strong></td>
            </tr>
            <tr>
              <td>Rate of profit (s / (c+v))</td>
              <td>{result.surplus_value_minutes} / {result.total_capital_advanced}</td>
              <td><strong>{(result.rate_of_profit * 100).toFixed(2)}%</strong></td>
            </tr>
            <tr>
              <td>Mill critique holds?</td>
              <td className="small muted">RateOfProfit &lt; RateSurplusValue</td>
              <td>{result.mill_critique_holds ? "yes" : "no"}</td>
            </tr>
          </tbody>
        </table>
      )}
    </section>
  );
}
