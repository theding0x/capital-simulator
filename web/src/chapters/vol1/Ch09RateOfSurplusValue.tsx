import { useEffect, useState } from "react";
import type { FormEvent } from "react";
import { api } from "../../api";
import type {
  RateOfSurplusValueResult,
  ProductionAccountResult,
  ProductPartsResult,
  WorkingDayRevisionResult,
} from "../../types";
import "./Ch09RateOfSurplusValue.css";

interface Ch09Props {
  onSharedChanged: () => void;
}

function minutesToHours(m: number): string {
  const h = Math.floor(m / 60);
  const min = m % 60;
  return min === 0 ? `${h}h` : `${h}h ${min}m`;
}

function DayBar({ v, s }: { v: number; s: number }) {
  const total = v + s;
  const pct = (n: number) => (total > 0 ? `${(n / total) * 100}%` : "0%");
  return (
    <>
      <div className="v1-ch09-day" role="img" aria-label="Working-day partition">
        <span className="v1-ch09-day-label">{minutesToHours(total)} day</span>
        <div className="v1-ch09-day-track">
          <div
            className="v1-ch09-day-seg v1-ch09-day-seg--necessary"
            style={{ flexBasis: pct(v) }}
            title={`Necessary labour (v): ${minutesToHours(v)}`}
          >
            v
          </div>
          <div
            className="v1-ch09-day-seg v1-ch09-day-seg--surplus"
            style={{ flexBasis: pct(s) }}
            title={`Surplus labour (s): ${minutesToHours(s)}`}
          >
            s
          </div>
        </div>
        <span className="v1-ch09-day-total">{minutesToHours(total)}</span>
      </div>
      <p className="v1-ch09-legend">
        <span>
          <span className="v1-ch09-legend-swatch" style={{ background: "var(--lead)" }} />
          v — necessary ({minutesToHours(v)})
        </span>
        <span>
          <span className="v1-ch09-legend-swatch" style={{ background: "var(--gold-bright)" }} />
          s — surplus ({minutesToHours(s)})
        </span>
      </p>
    </>
  );
}

function RateChip({ v, s }: { v: number; s: number }) {
  const rate = v > 0 ? ((s / v) * 100).toFixed(1) : "∞";
  const fraction = v + s > 0 ? ((s / (v + s)) * 100).toFixed(1) : "—";
  return (
    <div className="v1-ch09-rate">
      <span className="v1-ch09-rate-label">s / v (rate of exploitation)</span>
      <span className="v1-ch09-rate-value">{rate}%</span>
      <span className="v1-ch09-rate-meta">· surplus-produce {fraction}%</span>
    </div>
  );
}

function RateInsight() {
  return (
    <section className="v1-ch09-insight">
      <h2 className="v1-ch09-insight-h2">s/v is not the rate of profit</h2>
      <p className="v1-ch09-insight-prose">
        The rate of surplus-value compares surplus labour to <em>necessary</em>
        labour — not to the total capital advanced. Marx insists this is the
        rate of <strong>exploitation</strong> of labour-power by capital, and
        every variation in the working day, the wage, or productivity shows up
        first in this single number.
      </p>
    </section>
  );
}

function RateCoda() {
  return (
    <aside className="v1-ch09-coda">
      <p className="v1-ch09-coda-quote">
        “The rate of surplus-value is therefore an exact expression for the
        degree of exploitation of labour-power by capital, or of the labourer
        by the capitalist.”
        <span className="v1-ch09-coda-cite">
          — Marx, Capital Vol. I, Ch. 9 §1
        </span>
      </p>
    </aside>
  );
}

export function Ch09RateOfSurplusValue({ onSharedChanged: _onSharedChanged }: Ch09Props) {
  const [accounts, setAccounts] = useState<ProductionAccountResult[]>([]);

  async function refreshAccounts() {
    try {
      const list = await api.listProductionAccounts();
      setAccounts(list);
    } catch (_) {
      // non-fatal if listing fails
    }
  }

  useEffect(() => { refreshAccounts(); }, []);

  return (
    <>
      <RateInsight />
      <RateProbePanel />
      <ProductPartsPanel />
      <SeniorLastHourPanel />
      <RecordAccountPanel onCreated={refreshAccounts} />
      {accounts.length > 0 && <AccountListPanel accounts={accounts} />}
      <RateCoda />
    </>
  );
}

function RateProbePanel() {
  const [surplus, setSurplus] = useState(90);
  const [variable, setVariable] = useState(90);
  const [result, setResult] = useState<RateOfSurplusValueResult | null>(null);
  const [err, setErr] = useState<string | null>(null);

  async function submit(e: FormEvent) {
    e.preventDefault();
    setErr(null);
    try {
      const r = await api.computeRateOfSurplusValue({ surplus, variable });
      setResult(r);
    } catch (e) {
      setErr(e instanceof Error ? e.message : String(e));
    }
  }

  function loadFixture(s: number, v: number) {
    setSurplus(s);
    setVariable(v);
    setResult(null);
  }

  return (
    <section className="card">
      <h2>Rate of Surplus-Value Probe</h2>
      <p className="description">
        Compute s/v — the rate of surplus-value (degree of exploitation) from surplus and
        variable capital magnitudes in labour-minutes.
      </p>
      <div className="v1-ch09-presets">
        <span className="v1-ch09-presets-label">Marx fixture</span>
        <button
          type="button"
          className="v1-ch09-preset-button"
          onClick={() => loadFixture(80, 52)}
        >
          §1 Spinning Mill (s=80, v=52)
        </button>
        <button
          type="button"
          className="v1-ch09-preset-button"
          onClick={() => loadFixture(211, 210)}
        >
          Jacob's Wheat 1815 (s=211, v=210)
        </button>
        <button
          type="button"
          className="v1-ch09-preset-button"
          onClick={() => loadFixture(360, 360)}
        >
          Cotton Spinner (s=360, v=360)
        </button>
      </div>
      <form className="form-grid" onSubmit={submit}>
        <label>
          <span>Surplus Value s (min)</span>
          <input
            type="number"
            min={0}
            value={surplus}
            onChange={(e) => setSurplus(Number(e.target.value))}
          />
        </label>
        <label>
          <span>Variable Capital v (min)</span>
          <input
            type="number"
            min={1}
            value={variable}
            onChange={(e) => setVariable(Number(e.target.value))}
          />
        </label>
        <div className="form-actions span2">
          <button type="submit" className="primary">Compute Rate</button>
          {err && <span className="error">{err}</span>}
        </div>
      </form>

      {result && (
        <>
          <DayBar v={result.variable} s={result.surplus} />
          <RateChip v={result.variable} s={result.surplus} />
        </>
      )}
    </section>
  );
}

// Marx's yarn, §2: 20 lbs. worth 30s. = 24s. c + 3s. v + 3s. s, at 1s = 120 minutes.
const SHILLING_MINUTES = 120;

function formatPounds(n: number): string {
  return `${Number(n.toFixed(4))} lbs.`;
}

function ProductPartsPanel() {
  const [constant, setConstant] = useState(24 * SHILLING_MINUTES);
  const [variable, setVariable] = useState(3 * SHILLING_MINUTES);
  const [surplus, setSurplus] = useState(3 * SHILLING_MINUTES);
  const [total, setTotal] = useState(20);
  const [result, setResult] = useState<ProductPartsResult | null>(null);
  const [err, setErr] = useState<string | null>(null);

  async function submit(e: FormEvent) {
    e.preventDefault();
    setErr(null);
    try {
      setResult(
        await api.apportionProduct({
          constant,
          variable,
          surplus,
          total_product: total,
        })
      );
    } catch (e) {
      setResult(null);
      setErr(e instanceof Error ? e.message : String(e));
    }
  }

  function loadYarn() {
    setConstant(24 * SHILLING_MINUTES);
    setVariable(3 * SHILLING_MINUTES);
    setSurplus(3 * SHILLING_MINUTES);
    setTotal(20);
    setResult(null);
    setErr(null);
  }

  const rows: { label: string; value: number | undefined }[] = [
    { label: "Constant part (c) — cotton and spindle-waste in disguise", value: result?.parts.constant },
    { label: "Variable part (v) — the equivalent of the wages", value: result?.parts.variable },
    { label: "Surplus part (s) — the surplus-value", value: result?.parts.surplus },
  ];

  return (
    <section className="card">
      <h2>Parts of the Product</h2>
      <p className="description">
        §2: each component of the value of the product is contained in a
        corresponding proportional part of the product itself. Enter c, v and s
        (labour-minutes) and the total product.
      </p>
      <div className="v1-ch09-presets">
        <span className="v1-ch09-presets-label">Marx fixture</span>
        <button type="button" className="v1-ch09-preset-button" onClick={loadYarn}>
          §2 Yarn (24s. + 3s. + 3s.; 20 lbs.)
        </button>
      </div>
      <form className="form-grid" onSubmit={submit}>
        <label>
          <span>Constant Capital c (min)</span>
          <input type="number" min={0} value={constant} onChange={(e) => setConstant(Number(e.target.value))} />
        </label>
        <label>
          <span>Variable Capital v (min)</span>
          <input type="number" min={0} value={variable} onChange={(e) => setVariable(Number(e.target.value))} />
        </label>
        <label>
          <span>Surplus Value s (min)</span>
          <input type="number" min={0} value={surplus} onChange={(e) => setSurplus(Number(e.target.value))} />
        </label>
        <label>
          <span>Total product (lbs.)</span>
          <input type="number" min={0} step="any" value={total} onChange={(e) => setTotal(Number(e.target.value))} />
        </label>
        <div className="form-actions span2">
          <button type="submit" className="primary">Apportion Product</button>
          {err && <span className="error">{err}</span>}
        </div>
      </form>

      {result && (
        <table className="data-table">
          <thead>
            <tr>
              <th>Part of the product</th>
              <th>Quantity</th>
            </tr>
          </thead>
          <tbody>
            {rows.map((r) => (
              <tr key={r.label}>
                <td>{r.label}</td>
                <td>{r.value === undefined ? "—" : formatPounds(r.value)}</td>
              </tr>
            ))}
            <tr>
              <td><strong>Whole product</strong></td>
              <td><strong>{formatPounds(result.total_product)}</strong></td>
            </tr>
          </tbody>
        </table>
      )}
    </section>
  );
}

// Senior's "last hour", §3: wages and net profit are equal on his own figures,
// so 5¾ hours are necessary labour and 5¾ hours surplus labour in 11½ hours.
const SENIOR_NECESSARY_MINUTES = 345;

const SENIOR_PRESETS: { label: string; day: number }[] = [
  { label: "11½ h — Senior's day (100%)", day: 690 },
  { label: "13 h — lengthened (126 2/23%)", day: 780 },
  { label: "10½ h — cut by one hour (82 14/23%)", day: 630 },
];

function SeniorLastHourPanel() {
  const [necessary, setNecessary] = useState(SENIOR_NECESSARY_MINUTES);
  const [day, setDay] = useState(690);
  const [result, setResult] = useState<WorkingDayRevisionResult | null>(null);
  const [err, setErr] = useState<string | null>(null);

  async function run(necessaryLabour: number, workingDay: number) {
    setErr(null);
    try {
      setResult(await api.reviseWorkingDay({ necessary_labour: necessaryLabour, working_day: workingDay }));
    } catch (e) {
      setResult(null);
      setErr(e instanceof Error ? e.message : String(e));
    }
  }

  function submit(e: FormEvent) {
    e.preventDefault();
    void run(necessary, day);
  }

  function loadPreset(workingDay: number) {
    setNecessary(SENIOR_NECESSARY_MINUTES);
    setDay(workingDay);
    void run(SENIOR_NECESSARY_MINUTES, workingDay);
  }

  return (
    <section className="card">
      <h2>Senior's “Last Hour”</h2>
      <p className="description">
        §3: Senior claimed the whole net profit comes from the last hour of an
        11½-hour day. Hold the necessary labour fixed and re-cut the day: only
        the surplus labour moves, and the rate of surplus-value with it.
      </p>
      <div className="v1-ch09-presets">
        <span className="v1-ch09-presets-label">Marx fixture</span>
        {SENIOR_PRESETS.map((p) => (
          <button
            key={p.day}
            type="button"
            className="v1-ch09-preset-button"
            onClick={() => loadPreset(p.day)}
          >
            {p.label}
          </button>
        ))}
      </div>
      <form className="form-grid" onSubmit={submit}>
        <label>
          <span>Necessary labour (min, held fixed)</span>
          <input type="number" min={1} value={necessary} onChange={(e) => setNecessary(Number(e.target.value))} />
        </label>
        <label>
          <span>Working day (min)</span>
          <input type="number" min={0} max={1440} value={day} onChange={(e) => setDay(Number(e.target.value))} />
        </label>
        <div className="form-actions span2">
          <button type="submit" className="primary">Re-cut the Day</button>
          {err && <span className="error">{err}</span>}
        </div>
      </form>

      {result && (
        <>
          <DayBar v={result.necessary_labour} s={result.surplus_labour} />
          <div className="v1-ch09-rate">
            <span className="v1-ch09-rate-label">s / v (rate of surplus-value)</span>
            <span className="v1-ch09-rate-value">{(result.rate_basis_points / 100).toFixed(2)}%</span>
            <span className="v1-ch09-rate-meta">
              · surplus labour {minutesToHours(result.surplus_labour)} of {minutesToHours(result.working_day)}
            </span>
          </div>
          <p className="small muted">
            Marx's text says the day is cut “from 11½ to 10” hours, but his 4¾ hours
            and 82 14/23% only hold for 10½ hours, a cut of one hour.
          </p>
        </>
      )}
    </section>
  );
}

interface RecordAccountPanelProps {
  onCreated: () => void;
}

function RecordAccountPanel({ onCreated }: RecordAccountPanelProps) {
  const [constant, setConstant] = useState(410);
  const [variable, setVariable] = useState(90);
  const [surplus, setSurplus] = useState(90);
  const [result, setResult] = useState<ProductionAccountResult | null>(null);
  const [err, setErr] = useState<string | null>(null);

  async function submit(e: FormEvent) {
    e.preventDefault();
    setErr(null);
    try {
      const r = await api.createProductionAccount({ constant, variable, surplus });
      setResult(r);
      onCreated();
    } catch (e) {
      setErr(e instanceof Error ? e.message : String(e));
    }
  }

  return (
    <section className="card">
      <h2>Record Production Account</h2>
      <p className="description">
        Persist a production account recording c, v, s for a single production run.
        Rate of surplus-value is computed and stored automatically.
      </p>
      <form className="form-grid" onSubmit={submit}>
        <label>
          <span>Constant Capital c (min)</span>
          <input
            type="number"
            min={0}
            value={constant}
            onChange={(e) => setConstant(Number(e.target.value))}
          />
        </label>
        <label>
          <span>Variable Capital v (min)</span>
          <input
            type="number"
            min={1}
            value={variable}
            onChange={(e) => setVariable(Number(e.target.value))}
          />
        </label>
        <label>
          <span>Surplus Value s (min)</span>
          <input
            type="number"
            min={0}
            value={surplus}
            onChange={(e) => setSurplus(Number(e.target.value))}
          />
        </label>
        <div className="form-actions span2">
          <button type="submit" className="primary">Record Account</button>
          {err && <span className="error">{err}</span>}
        </div>
      </form>

      {result && (
        <div>
          <p className="small muted">Saved: <span className="monospace">{result.id}</span></p>
          <DayBar v={result.variable} s={result.surplus} />
          <RateChip v={result.variable} s={result.surplus} />
          <p className="small muted" style={{ marginTop: "0.5rem" }}>
            c {minutesToHours(result.constant)} · v + s = {minutesToHours(result.value_product)} ·
            c + v + s = {minutesToHours(result.expanded_capital)}
          </p>
        </div>
      )}
    </section>
  );
}

function AccountListPanel({ accounts }: { accounts: ProductionAccountResult[] }) {
  return (
    <section className="card">
      <h2>Recorded Production Accounts</h2>
      <table className="data-table">
        <thead>
          <tr>
            <th>ID</th>
            <th>c</th>
            <th>v</th>
            <th>s</th>
            <th>Rate (s/v)</th>
            <th>Expanded Capital</th>
          </tr>
        </thead>
        <tbody>
          {accounts.map((a) => (
            <tr key={a.id}>
              <td className="monospace">{a.id.slice(0, 8)}…</td>
              <td>{minutesToHours(a.constant)}</td>
              <td>{minutesToHours(a.variable)}</td>
              <td>{minutesToHours(a.surplus)}</td>
              <td>{(a.rate_of_surplus_value * 100).toFixed(2)}%</td>
              <td>{minutesToHours(a.expanded_capital)}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </section>
  );
}
