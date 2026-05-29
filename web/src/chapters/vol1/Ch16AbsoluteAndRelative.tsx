import { useCallback, useEffect, useState } from "react";
import type { FormEvent } from "react";
import { api } from "../../api";
import { fmtHoursLong } from "../../format";
import type {
  AbsoluteSurplusValueResult,
  RelativeSurplusValueResult,
  SurplusValueRateResult,
} from "../../types";
import "./Ch16AbsoluteAndRelative.css";

export function Ch16AbsoluteAndRelative() {
  return (
    <>
      <SeparationInsight />
      <SurplusValueCalculator />
      <MillCritiquePanel />
    </>
  );
}

function SeparationInsight() {
  return (
    <section className="v1-ch16-insight">
      <h2 className="v1-ch16-insight-h2">Same magnitude, two Origins</h2>
      <p className="v1-ch16-insight-prose">
        The two surplus-value forms produce the same magnitude type: a quantity
        of <em>surplus labour</em> embodied in product. What distinguishes them
        is the <strong>Origin</strong>: absolute surplus prolongs the working
        day; relative surplus shrinks necessary labour. The bars below split
        the same starting partition along the two paths.
      </p>
    </section>
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
        Configure a starting partition and watch both mechanisms produce
        surplus from the same input.
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
      <div className="v1-ch16-results">
        <ResultCard
          flavour="absolute"
          heading="Absolute — prolongation"
          subtitle="Working day grows; necessary labour fixed."
          oldWD={abs?.old_working_day}
          newWD={abs?.new_working_day}
          surplus={abs?.absolute_surplus_value.labour_minutes ?? 0}
          origin={abs?.absolute_surplus_value.origin ?? "absolute"}
        />
        <ResultCard
          flavour="relative"
          heading="Relative — productivity"
          subtitle="Necessary labour shrinks; working day fixed."
          oldWD={rel?.old_working_day}
          newWD={rel?.new_working_day}
          surplus={rel?.relative_surplus_value.labour_minutes ?? 0}
          origin={rel?.relative_surplus_value.origin ?? "relative"}
        />
      </div>
      <p className="v1-ch16-closing-prose">
        §1, closing: “From one standpoint, any distinction between
        absolute and relative surplus-value appears illusory. Relative
        surplus-value is absolute … Absolute surplus-value is relative.”
        The two share a magnitude type; the Origin tag is the line between them.
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
  flavour,
  heading,
  subtitle,
  oldWD,
  newWD,
  surplus,
  origin,
}: {
  flavour: "absolute" | "relative";
  heading: string;
  subtitle: string;
  oldWD?: ProductionWorkingDay;
  newWD?: ProductionWorkingDay;
  surplus: number;
  origin: string;
}) {
  if (!oldWD || !newWD) {
    return (
      <div className="v1-ch16-result">
        <span className={`v1-ch16-result-tag v1-ch16-result-tag--${flavour}`}>
          {flavour}
        </span>
        <h3 className="v1-ch16-result-h3">{heading}</h3>
        <p className="v1-ch16-result-sub">{subtitle}</p>
        <p className="v1-ch16-result-sub">—</p>
      </div>
    );
  }
  const max = Math.max(oldWD.total, newWD.total, 1);
  const pct = (n: number) => `${(n / max) * 100}%`;
  return (
    <div className="v1-ch16-result">
      <span className={`v1-ch16-result-tag v1-ch16-result-tag--${flavour}`}>
        {flavour}
      </span>
      <h3 className="v1-ch16-result-h3">{heading}</h3>
      <p className="v1-ch16-result-sub">{subtitle}</p>
      <div className="v1-ch16-bars">
        <div className="v1-ch16-bar-row">
          <span className="v1-ch16-bar-label">Before</span>
          <div className="v1-ch16-bar-track">
            <div
              className="v1-ch16-bar-seg--necessary"
              style={{ flexBasis: pct(oldWD.necessary_labour) }}
              title={`v: ${fmtHoursLong(oldWD.necessary_labour)}`}
            />
            <div
              className="v1-ch16-bar-seg--surplus"
              style={{ flexBasis: pct(oldWD.surplus_labour) }}
              title={`s: ${fmtHoursLong(oldWD.surplus_labour)}`}
            />
          </div>
          <span className="v1-ch16-bar-total">{fmtHoursLong(oldWD.total)}</span>
        </div>
        <div className="v1-ch16-bar-row">
          <span className="v1-ch16-bar-label">After</span>
          <div className="v1-ch16-bar-track">
            <div
              className="v1-ch16-bar-seg--necessary"
              style={{ flexBasis: pct(newWD.necessary_labour) }}
              title={`v: ${fmtHoursLong(newWD.necessary_labour)}`}
            />
            <div
              className="v1-ch16-bar-seg--surplus"
              style={{ flexBasis: pct(newWD.surplus_labour) }}
              title={`s: ${fmtHoursLong(newWD.surplus_labour)}`}
            />
          </div>
          <span className="v1-ch16-bar-total">{fmtHoursLong(newWD.total)}</span>
        </div>
      </div>
      <div className="v1-ch16-gain">
        <span className="v1-ch16-gain-label">Gain</span>
        <span className="v1-ch16-gain-value">+{surplus} min</span>
        <span className="v1-ch16-origin-tag">Origin: {origin}</span>
      </div>
    </div>
  );
}

function MillCritiquePanel() {
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
        will be 20:500, i.e., 4% and not 20%." Same numerator, different
        denominators. The Mill error is to collapse them.
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
        <>
          <div className="v1-ch16-mill">
            <div className="v1-ch16-mill-row">
              <span className="v1-ch16-mill-label">Rate of surplus-value</span>
              <span className="v1-ch16-mill-formula">
                s / v · {result.surplus_labour_minutes} / {result.necessary_labour_minutes}
              </span>
              <span className="v1-ch16-mill-value">
                {(result.rate_of_surplus_value * 100).toFixed(2)}%
              </span>
            </div>
            <div className="v1-ch16-mill-row">
              <span className="v1-ch16-mill-label">Rate of profit</span>
              <span className="v1-ch16-mill-formula">
                s / (c+v) · {result.surplus_value_minutes} / {result.total_capital_advanced}
              </span>
              <span className="v1-ch16-mill-value">
                {(result.rate_of_profit * 100).toFixed(2)}%
              </span>
            </div>
          </div>
          <span
            className={
              "v1-ch16-badge " +
              (result.mill_critique_holds
                ? "v1-ch16-badge--holds"
                : "v1-ch16-badge--breaks")
            }
          >
            {result.mill_critique_holds
              ? "Mill critique holds — rate of profit < rate of surplus-value"
              : "Mill error — rates collapsed"}
          </span>
        </>
      )}
    </section>
  );
}
