import { useState } from "react";
import type { FormEvent } from "react";
import { api } from "../../api";
import type {
  ProductionWorkingDay,
  ShortenWorkingDayResponse,
  ProductionRateResult,
  ExtraSurplusValueResult,
} from "../../types";
import "./Ch12RelativeSurplusValue.css";

interface Ch12Props {
  onSharedChanged: () => void;
}

export function Ch12RelativeSurplusValue({ onSharedChanged: _unused }: Ch12Props) {
  return (
    <>
      <SequenceInsight />
      <WorkingDayPanel />
      <ShortenWorkingDayPanel />
      <ExtraSurplusValuePanel />
      <RelativeCoda />
    </>
  );
}

function SequenceInsight() {
  return (
    <section className="v1-ch12-insight">
      <h2 className="v1-ch12-insight-h2">The mechanism of relative surplus-value</h2>
      <div className="v1-ch12-sequence">
        <div className="v1-ch12-step">
          <span className="v1-ch12-step-tag">Productivity</span>
          <span className="v1-ch12-step-eqn">π ↑</span>
          <span className="v1-ch12-step-arrow">→</span>
        </div>
        <div className="v1-ch12-step">
          <span className="v1-ch12-step-tag">Value of LP</span>
          <span className="v1-ch12-step-eqn">v ↓</span>
          <span className="v1-ch12-step-arrow">→</span>
        </div>
        <div className="v1-ch12-step">
          <span className="v1-ch12-step-tag">Necessary labour</span>
          <span className="v1-ch12-step-eqn">a ↓</span>
          <span className="v1-ch12-step-arrow">→</span>
        </div>
        <div className="v1-ch12-step">
          <span className="v1-ch12-step-tag">Surplus labour</span>
          <span className="v1-ch12-step-eqn">a′ ↑</span>
        </div>
      </div>
      <p className="v1-ch12-insight-prose">
        With the working day held constant, a productivity rise in the
        wage-goods industries cheapens labour-power. Necessary labour contracts,
        and surplus labour expands by the same amount — relative surplus-value.
      </p>
    </section>
  );
}

function DayBar({
  necessary,
  surplus,
  label,
  total,
}: {
  necessary: number;
  surplus: number;
  label: string;
  total?: number;
}) {
  const sum = total ?? necessary + surplus;
  const pct = (n: number) => (sum > 0 ? `${(n / sum) * 100}%` : "0%");
  return (
    <div className="v1-ch12-bar" role="img" aria-label={`${label}: necessary ${necessary}, surplus ${surplus}`}>
      <span className="v1-ch12-bar-label">{label}</span>
      <div className="v1-ch12-bar-track">
        <div
          className="v1-ch12-bar-seg v1-ch12-bar-seg--necessary"
          style={{ flexBasis: pct(necessary) }}
          title={`Necessary: ${necessary} min`}
        >
          v
        </div>
        <div
          className="v1-ch12-bar-seg v1-ch12-bar-seg--surplus"
          style={{ flexBasis: pct(surplus) }}
          title={`Surplus: ${surplus} min`}
        >
          s
        </div>
      </div>
      <span className="v1-ch12-bar-total">{sum} min</span>
    </div>
  );
}

function RelativeCoda() {
  return (
    <aside className="v1-ch12-coda">
      <p className="v1-ch12-coda-quote">
        “I call that surplus-value which is produced by the curtailment of the
        necessary labour-time, and by the corresponding alteration in the
        respective lengths of the two components of the working day, relative
        surplus-value.”
        <span className="v1-ch12-coda-cite">
          — Marx, Capital Vol. I, Ch. 12 §1
        </span>
      </p>
    </aside>
  );
}

// ── Working Day Panel ──────────────────────────────────────────────────────────

function WorkingDayPanel() {
  const [total, setTotal] = useState(720);
  const [lpv, setLpv] = useState(600);
  const [result, setResult] = useState<ProductionWorkingDay | null>(null);
  const [rateResult, setRateResult] = useState<ProductionRateResult | null>(null);
  const [err, setErr] = useState<string | null>(null);

  async function submit(e: FormEvent) {
    e.preventDefault();
    setErr(null);
    try {
      const wd = await api.recordWorkingDay({ total, labour_power_value: lpv });
      setResult(wd);
      const rate = await api.getProductionRate(wd.necessary_labour, wd.surplus_labour);
      setRateResult(rate);
    } catch (e) {
      setErr(e instanceof Error ? e.message : String(e));
    }
  }

  return (
    <section className="card">
      <h2>Working-Day Split</h2>
      <p className="description">
        Records the split of the working day into necessary labour (reproducing labour-power value)
        and surplus labour (producing surplus-value for capital). §1: "The working day is thus not
        a constant, but a variable quantity."
      </p>
      <form className="form-grid" onSubmit={submit}>
        <label>
          <span>Total Working Day (min)</span>
          <input
            type="number"
            min={1}
            value={total}
            onChange={(e) => setTotal(Number(e.target.value))}
          />
        </label>
        <label>
          <span>Labour-Power Value (min)</span>
          <input
            type="number"
            min={1}
            value={lpv}
            onChange={(e) => setLpv(Number(e.target.value))}
          />
        </label>
        <div className="form-actions span2">
          <button type="submit" className="primary">Record</button>
          {err && <span className="error">{err}</span>}
        </div>
      </form>
      {result && (
        <>
          <DayBar
            necessary={result.necessary_labour}
            surplus={result.surplus_labour}
            label="Working day"
          />
          <p className="v1-ch12-legend">
            <span>
              <span className="v1-ch12-legend-swatch" style={{ background: "var(--lead)" }} />
              v — necessary ({result.necessary_labour} min)
            </span>
            <span>
              <span className="v1-ch12-legend-swatch" style={{ background: "var(--gold-bright)" }} />
              s — surplus ({result.surplus_labour} min)
            </span>
          </p>
          {rateResult && (
            <div className="v1-ch12-rate">
              <span className="v1-ch12-rate-label">s / v</span>
              <span className="v1-ch12-rate-value">{(rateResult.rate * 100).toFixed(1)}%</span>
            </div>
          )}
        </>
      )}
    </section>
  );
}

// ── Shorten Necessary Labour Panel ────────────────────────────────────────────

const SHORTEN_FIXTURES = [
  { label: "§1 Surplus ½→⅓ (9h→8h NL)", total: 720, nl: 540, sl: 180, newLpv: 480 },
  { label: "§1 LPV 5s→3s (600→360 min NL)", total: 720, nl: 600, sl: 120, newLpv: 360 },
];

function ShortenWorkingDayPanel() {
  const [total, setTotal] = useState(720);
  const [nl, setNl] = useState(600);
  const [sl, setSl] = useState(120);
  const [newLpv, setNewLpv] = useState(540);
  const [result, setResult] = useState<ShortenWorkingDayResponse | null>(null);
  const [err, setErr] = useState<string | null>(null);

  function loadFixture(f: typeof SHORTEN_FIXTURES[number]) {
    setTotal(f.total); setNl(f.nl); setSl(f.sl); setNewLpv(f.newLpv);
    setResult(null); setErr(null);
  }

  async function submit(e: FormEvent) {
    e.preventDefault();
    setErr(null);
    try {
      const r = await api.shortenWorkingDay({
        total,
        necessary_labour: nl,
        surplus_labour: sl,
        new_labour_power_value: newLpv,
      });
      setResult(r);
    } catch (e) {
      setErr(e instanceof Error ? e.message : String(e));
    }
  }

  return (
    <section className="card">
      <h2>Shorten Necessary Labour (Relative Surplus-Value)</h2>
      <p className="description">
        A productivity rise reduces the value of labour-power, contracting necessary labour and
        expanding surplus labour within the same working day. The gain in surplus labour is the
        relative surplus-value. §1: "I call that surplus-value which is produced by curtailment
        of the necessary labour-time … relative surplus-value."
      </p>
      <div className="v1-ch12-presets">
        <span className="v1-ch12-presets-label">Marx fixture</span>
        {SHORTEN_FIXTURES.map((f) => (
          <button
            key={f.label}
            type="button"
            className="v1-ch12-preset-button"
            onClick={() => loadFixture(f)}
          >
            {f.label}
          </button>
        ))}
      </div>
      <form className="form-grid" onSubmit={submit}>
        <label>
          <span>Total (min)</span>
          <input type="number" min={1} value={total} onChange={(e) => setTotal(Number(e.target.value))} />
        </label>
        <label>
          <span>Current Necessary Labour (min)</span>
          <input type="number" min={1} value={nl} onChange={(e) => setNl(Number(e.target.value))} />
        </label>
        <label>
          <span>Current Surplus Labour (min)</span>
          <input type="number" min={0} value={sl} onChange={(e) => setSl(Number(e.target.value))} />
        </label>
        <label>
          <span>New Labour-Power Value (min)</span>
          <input type="number" min={1} value={newLpv} onChange={(e) => setNewLpv(Number(e.target.value))} />
        </label>
        <div className="form-actions span2">
          <button type="submit" className="primary">Shorten</button>
          {err && <span className="error">{err}</span>}
        </div>
      </form>
      {result && (
        <>
          <DayBar necessary={nl} surplus={sl} label="Before" total={total} />
          <DayBar
            necessary={result.working_day.necessary_labour}
            surplus={result.working_day.surplus_labour}
            label="After"
            total={total}
          />
          <span className="v1-ch12-delta">
            ∆ surplus labour = +{result.relative_surplus_value} min (relative surplus-value)
          </span>
          <div className="v1-ch12-rate">
            <span className="v1-ch12-rate-label">New s / v</span>
            <span className="v1-ch12-rate-value">
              {result.working_day.necessary_labour > 0
                ? ((result.working_day.surplus_labour / result.working_day.necessary_labour) * 100).toFixed(1) + "%"
                : "∞"}
            </span>
          </div>
        </>
      )}
    </section>
  );
}

// ── Extra Surplus-Value Panel ──────────────────────────────────────────────────

const EXTRA_FIXTURES = [
  { label: "§1 Double productivity (iv=30, sv=60, qty=24)", iv: 30, sv: 60, qty: 24 },
  { label: "§1 Per-article gain (qty=1)", iv: 30, sv: 60, qty: 1 },
];

function ExtraSurplusValuePanel() {
  const [iv, setIv] = useState(30);
  const [sv, setSv] = useState(60);
  const [qty, setQty] = useState(24);
  const [result, setResult] = useState<ExtraSurplusValueResult | null>(null);
  const [err, setErr] = useState<string | null>(null);

  function loadFixture(f: typeof EXTRA_FIXTURES[number]) {
    setIv(f.iv); setSv(f.sv); setQty(f.qty);
    setResult(null); setErr(null);
  }

  async function submit(e: FormEvent) {
    e.preventDefault();
    setErr(null);
    try {
      const r = await api.computeExtraSurplusValue({ individual_value: iv, social_value: sv, quantity: qty });
      setResult(r);
    } catch (e) {
      setErr(e instanceof Error ? e.message : String(e));
    }
  }

  return (
    <section className="card">
      <h2>Extra Surplus-Value Probe</h2>
      <p className="description">
        When a capitalist's individual value falls below the social value, they gain an extra
        surplus-value by selling at the higher social value. This temporary advantage drives
        generalisation of the new productivity. §1: "The capitalist … realises an extra
        surplus-value."
      </p>
      <div className="v1-ch12-presets">
        <span className="v1-ch12-presets-label">Marx fixture</span>
        {EXTRA_FIXTURES.map((f) => (
          <button
            key={f.label}
            type="button"
            className="v1-ch12-preset-button"
            onClick={() => loadFixture(f)}
          >
            {f.label}
          </button>
        ))}
      </div>
      <form className="form-grid" onSubmit={submit}>
        <label>
          <span>Individual Value (min per unit)</span>
          <input type="number" min={1} value={iv} onChange={(e) => setIv(Number(e.target.value))} />
        </label>
        <label>
          <span>Social Value (min per unit)</span>
          <input type="number" min={1} value={sv} onChange={(e) => setSv(Number(e.target.value))} />
        </label>
        <label>
          <span>Quantity produced</span>
          <input type="number" min={1} value={qty} onChange={(e) => setQty(Number(e.target.value))} />
        </label>
        <div className="form-actions span2">
          <button type="submit" className="primary">Compute</button>
          {err && <span className="error">{err}</span>}
        </div>
      </form>
      {result && (
        <table className="data-table">
          <tbody>
            <tr><td>Individual Value</td><td>{result.individual_value} min/unit</td></tr>
            <tr><td>Social Value</td><td>{result.social_value} min/unit</td></tr>
            <tr><td>Quantity</td><td>{result.quantity}</td></tr>
            <tr><td>Per-Unit Extra Gain</td><td>{result.per_unit} min</td></tr>
            <tr>
              <td><strong>Total Extra Surplus-Value</strong></td>
              <td><strong>{result.extra_surplus_value} min</strong></td>
            </tr>
            {result.extra_surplus_value === 0 && (
              <tr>
                <td colSpan={2} className="muted">
                  No extra surplus-value — individual value ≥ social value.
                </td>
              </tr>
            )}
          </tbody>
        </table>
      )}
    </section>
  );
}
