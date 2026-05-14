import { useState } from "react";
import type { FormEvent } from "react";
import { api } from "../api";
import type {
  ProductionWorkingDay,
  ShortenWorkingDayResponse,
  ProductionRateResult,
  ExtraSurplusValueResult,
} from "../types";

interface Ch12Props {
  onSharedChanged: () => void;
}

export function Ch12RelativeSurplusValue({ onSharedChanged: _unused }: Ch12Props) {
  return (
    <>
      <WorkingDayPanel />
      <ShortenWorkingDayPanel />
      <ExtraSurplusValuePanel />
    </>
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

  const necessaryPct = result ? ((result.necessary_labour / result.total) * 100).toFixed(1) : null;
  const surplusPct = result ? ((result.surplus_labour / result.total) * 100).toFixed(1) : null;

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
          <div style={{ marginTop: "1rem" }}>
            <div style={{ display: "flex", height: "1.5rem", borderRadius: "4px", overflow: "hidden", fontSize: "0.75rem" }}>
              <div
                style={{
                  width: `${necessaryPct}%`,
                  background: "var(--color-accent)",
                  display: "flex",
                  alignItems: "center",
                  justifyContent: "center",
                  color: "#fff",
                  fontWeight: 600,
                }}
              >
                Necessary {necessaryPct}%
              </div>
              <div
                style={{
                  width: `${surplusPct}%`,
                  background: "var(--color-muted)",
                  display: "flex",
                  alignItems: "center",
                  justifyContent: "center",
                  color: "#fff",
                  fontWeight: 600,
                }}
              >
                Surplus {surplusPct}%
              </div>
            </div>
          </div>
          <table className="data-table" style={{ marginTop: "0.75rem" }}>
            <tbody>
              <tr><td>Total</td><td>{result.total} min ({(result.total / 60).toFixed(1)} h)</td></tr>
              <tr><td>Necessary Labour</td><td>{result.necessary_labour} min ({(result.necessary_labour / 60).toFixed(1)} h)</td></tr>
              <tr><td>Surplus Labour</td><td>{result.surplus_labour} min ({(result.surplus_labour / 60).toFixed(1)} h)</td></tr>
              {rateResult && (
                <tr>
                  <td><strong>Rate of Surplus-Value s/v</strong></td>
                  <td><strong>{(rateResult.rate * 100).toFixed(1)}%</strong></td>
                </tr>
              )}
            </tbody>
          </table>
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
      <div style={{ marginBottom: "0.75rem", display: "flex", gap: "0.5rem", flexWrap: "wrap" }}>
        {SHORTEN_FIXTURES.map((f) => (
          <button key={f.label} type="button" onClick={() => loadFixture(f)}>{f.label}</button>
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
        <table className="data-table">
          <tbody>
            <tr><td>New Total</td><td>{result.working_day.total} min</td></tr>
            <tr><td>New Necessary Labour</td><td>{result.working_day.necessary_labour} min</td></tr>
            <tr><td>New Surplus Labour</td><td>{result.working_day.surplus_labour} min</td></tr>
            <tr>
              <td><strong>Relative Surplus-Value (delta)</strong></td>
              <td><strong>+{result.relative_surplus_value} min</strong></td>
            </tr>
            <tr>
              <td>New Rate s/v</td>
              <td>
                {result.working_day.necessary_labour > 0
                  ? ((result.working_day.surplus_labour / result.working_day.necessary_labour) * 100).toFixed(1) + "%"
                  : "—"}
              </td>
            </tr>
          </tbody>
        </table>
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
      <div style={{ marginBottom: "0.75rem", display: "flex", gap: "0.5rem", flexWrap: "wrap" }}>
        {EXTRA_FIXTURES.map((f) => (
          <button key={f.label} type="button" onClick={() => loadFixture(f)}>{f.label}</button>
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
