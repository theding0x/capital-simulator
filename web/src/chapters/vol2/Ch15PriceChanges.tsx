import { useState } from "react";
import type {
  PriceRevolutionRecord,
  ValueRevolutionAffecting,
} from "../../types";
import { priceRevolutionApi } from "../../api";
import "./Ch15PriceChanges.css";

// Vol. II Ch. 15 — The Effects of a Change of Prices.
//
// A change in the prices of the elements of productive capital affects the
// circuit in three ways:
//   §1 — Price fall/rise on means of production → capital freed or tied up.
//   §2 — Price change on inventory held → paper gain/loss (unrealised).
//   §3 — Price fall after production but before sale → loss of surplus-value.
//
// Seed fixtures (00049_v2_ch15_seed.sql):
//   A — Spinning Mill 1871: cotton −15% → capital freed £72 (7200p).
//   B — Manchester Spinner 1857 crisis: output −10% → loss of £60 (6000p).
//   C — Cotton stock: 5000 lbs × (10p−12p) = −£100 paper loss.
//   D — Speculative hold: manchester-cotton-1857, anticipated price rise.

function pence(p: number): string {
  if (Math.abs(p) >= 240) {
    const pounds = p / 240;
    return `£${pounds % 1 === 0 ? pounds : pounds.toFixed(2)}`;
  }
  return `${p}d`;
}

function signedPence(p: number): string {
  return (p >= 0 ? "+" : "") + pence(p);
}

function affectingLabel(a: ValueRevolutionAffecting): string {
  switch (a) {
    case "means_of_production": return "Means of Production";
    case "labour_power":        return "Labour Power";
    case "output":              return "Output (C′)";
  }
}

function bpLabel(bp: number): string {
  return bp === 0 ? "0 bp" : `${bp > 0 ? "+" : ""}${bp} bp (${(bp / 100).toFixed(1)}%)`;
}

// ── EventCard ────────────────────────────────────────────────────────────────

function EventCard({ rec }: { rec: PriceRevolutionRecord }) {
  const ev = rec.event;
  const delta = rec.revalued_capital?.delta_pence ?? rec.realised_delta?.delta_pence ?? 0;
  const positive = delta >= 0;

  return (
    <div className={`ch15-event-card ${positive ? "ch15-positive" : "ch15-negative"}`}>
      <div className="ch15-event-header">
        <span className="ch15-affecting">{affectingLabel(ev.affecting)}</span>
        <span className={`ch15-delta ${positive ? "ch15-gain" : "ch15-loss"}`}>
          {signedPence(delta)}
        </span>
      </div>

      <div className="ch15-event-meta">
        <span>{new Date(ev.occurred_at).toLocaleDateString("en-GB", { year: "numeric", month: "short", day: "numeric" })}</span>
        <span className="ch15-bp">{bpLabel(ev.basis_points_change)}</span>
      </div>

      {rec.revalued_capital && (
        <div className="ch15-detail-row">
          <span>Capital advance:</span>
          <span>
            {pence(rec.revalued_capital.old_advance_pence)}
            &nbsp;→&nbsp;
            <strong>{pence(rec.revalued_capital.new_advance_pence)}</strong>
          </span>
        </div>
      )}

      {rec.realised_delta && (
        <div className="ch15-detail-row">
          <span>Realisation:</span>
          <span>
            expected {pence(rec.realised_delta.expected_realisation_pence)},
            actual {pence(rec.realised_delta.actual_realisation_pence)}
          </span>
        </div>
      )}

      {rec.fall_case && (
        <div className="ch15-case-tag ch15-case-fall">
          Fall case: <code>{rec.fall_case}</code>
        </div>
      )}
      {rec.rise_case && (
        <div className="ch15-case-tag ch15-case-rise">
          Rise case: <code>{rec.rise_case}</code>
        </div>
      )}

      {rec.inventory_revaluations.length > 0 && (
        <table className="ch15-inv-table">
          <thead>
            <tr>
              <th>Commodity</th><th>Qty</th><th>Old</th><th>New</th><th>Δ</th>
            </tr>
          </thead>
          <tbody>
            {rec.inventory_revaluations.map((ir) => (
              <tr key={ir.id}>
                <td>{ir.commodity_id}</td>
                <td>{ir.quantity_held}</td>
                <td>{pence(ir.old_unit_pence)}</td>
                <td>{pence(ir.new_unit_pence)}</td>
                <td className={ir.delta_pence < 0 ? "ch15-loss" : "ch15-gain"}>
                  {signedPence(ir.delta_pence)}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  );
}

// ── CreateForm ───────────────────────────────────────────────────────────────

type FormState = {
  industrial_capital_id: string;
  affecting: ValueRevolutionAffecting;
  direction_pence: string;
  basis_points_change: string;
  old_advance_pence: string;
  expected_realisation_pence: string;
  actual_realisation_pence: string;
};

const BLANK: FormState = {
  industrial_capital_id: "",
  affecting: "means_of_production",
  direction_pence: "",
  basis_points_change: "",
  old_advance_pence: "",
  expected_realisation_pence: "",
  actual_realisation_pence: "",
};

// Fixture presets matching seed IDs.
const FIXTURE_A: FormState = {
  industrial_capital_id: "5eed0000000000001501",
  affecting: "means_of_production",
  direction_pence: "-7200",
  basis_points_change: "-1500",
  old_advance_pence: "48000",
  expected_realisation_pence: "",
  actual_realisation_pence: "",
};

const FIXTURE_B: FormState = {
  industrial_capital_id: "5eed0000000000001504",
  affecting: "output",
  direction_pence: "-6000",
  basis_points_change: "0",
  old_advance_pence: "",
  expected_realisation_pence: "60000",
  actual_realisation_pence: "54000",
};

function CreateForm({
  onCreated,
}: {
  onCreated: (rec: PriceRevolutionRecord) => void;
}) {
  const [form, setForm] = useState<FormState>(BLANK);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  const set =
    (k: keyof FormState) =>
    (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) =>
      setForm((f) => ({ ...f, [k]: e.target.value }));

  async function submit(e: React.FormEvent) {
    e.preventDefault();
    setError(null);
    setLoading(true);
    try {
      const rec = await priceRevolutionApi.create({
        industrial_capital_id: form.industrial_capital_id,
        affecting: form.affecting,
        direction_pence: parseInt(form.direction_pence, 10) || 0,
        basis_points_change: parseInt(form.basis_points_change, 10) || 0,
        old_advance_pence: form.old_advance_pence
          ? parseInt(form.old_advance_pence, 10)
          : undefined,
        expected_realisation_pence: form.expected_realisation_pence
          ? parseInt(form.expected_realisation_pence, 10)
          : undefined,
        actual_realisation_pence: form.actual_realisation_pence
          ? parseInt(form.actual_realisation_pence, 10)
          : undefined,
      });
      onCreated(rec);
      setForm(BLANK);
    } catch (err) {
      setError(err instanceof Error ? err.message : String(err));
    } finally {
      setLoading(false);
    }
  }

  return (
    <form className="ch15-form" onSubmit={submit}>
      <div className="ch15-fixture-row">
        <button
          type="button"
          className="ch15-fixture-btn"
          onClick={() => setForm(FIXTURE_A)}
        >
          § A — Spinning Mill 1871 (cotton −15%)
        </button>
        <button
          type="button"
          className="ch15-fixture-btn"
          onClick={() => setForm(FIXTURE_B)}
        >
          § B — Manchester Spinner 1857 Crisis
        </button>
      </div>

      <label>
        Industrial Capital ID
        <input
          value={form.industrial_capital_id}
          onChange={set("industrial_capital_id")}
          required
        />
      </label>

      <label>
        Affecting
        <select value={form.affecting} onChange={set("affecting")}>
          <option value="means_of_production">Means of Production</option>
          <option value="labour_power">Labour Power</option>
          <option value="output">Output (C′)</option>
        </select>
      </label>

      <label>
        Price change (pence)
        <input
          type="number"
          value={form.direction_pence}
          onChange={set("direction_pence")}
          required
          placeholder="e.g. −7200"
        />
      </label>

      <label>
        Basis points change
        <input
          type="number"
          value={form.basis_points_change}
          onChange={set("basis_points_change")}
          placeholder="e.g. −1500 = −15%"
        />
      </label>

      {(form.affecting === "means_of_production" ||
        form.affecting === "labour_power") && (
        <label>
          Old advance (pence) — derives revalued capital
          <input
            type="number"
            value={form.old_advance_pence}
            onChange={set("old_advance_pence")}
            placeholder="e.g. 48000"
          />
        </label>
      )}

      {form.affecting === "output" && (
        <>
          <label>
            Expected realisation (pence)
            <input
              type="number"
              value={form.expected_realisation_pence}
              onChange={set("expected_realisation_pence")}
              placeholder="e.g. 60000"
            />
          </label>
          <label>
            Actual realisation (pence)
            <input
              type="number"
              value={form.actual_realisation_pence}
              onChange={set("actual_realisation_pence")}
              placeholder="e.g. 54000"
            />
          </label>
        </>
      )}

      {error && <p className="ch15-error">{error}</p>}
      <button type="submit" className="ch15-submit" disabled={loading}>
        {loading ? "Recording…" : "Record Price Revolution"}
      </button>
    </form>
  );
}

// ── Main panel ───────────────────────────────────────────────────────────────

export function Ch15PriceChanges() {
  const [records, setRecords] = useState<PriceRevolutionRecord[]>([]);

  function onCreated(rec: PriceRevolutionRecord) {
    setRecords((rs) => [rec, ...rs]);
  }

  return (
    <div className="ch15-root">
      <blockquote className="ch15-quote">
        "If the value of the elements of production rises in price, a part of
        the capital previously employed as money-capital cannot be replaced after
        the circuit is completed; hence it disappears from the circuit."
        <cite> — Marx, Capital Vol. II, Ch. 15</cite>
      </blockquote>

      <div className="ch15-layout">
        <section className="ch15-left">
          <h3>Record a Price Revolution</h3>
          <p className="ch15-hint">
            A "price revolution" is any change in market prices of the elements
            of productive capital — raw materials, labour power — or in the
            output price. Each event either frees or ties up additional capital
            and redistributes surplus-value between producer and buyer.
          </p>
          <CreateForm onCreated={onCreated} />
        </section>

        <section className="ch15-right">
          <h3>Price Revolution Events</h3>
          {records.length === 0 ? (
            <p className="ch15-empty">
              No events recorded this session. Use a fixture button or fill the
              form to create one.
            </p>
          ) : (
            <div className="ch15-event-list">
              {records.map((r) => (
                <EventCard key={r.event.id} rec={r} />
              ))}
            </div>
          )}
        </section>
      </div>

      <section className="ch15-cases">
        <h3>Marx's Classification of Price Cases</h3>
        <div className="ch15-cases-grid">
          <div className="ch15-case-block ch15-case-fall-block">
            <h4>Price Fall on MP / LP</h4>
            <dl>
              <dt>
                <code>same_scale</code>
              </dt>
              <dd>
                Production continues at same scale. Freed capital hoarded as
                money.
              </dd>
              <dt>
                <code>expanded_scale</code>
              </dt>
              <dd>
                Freed capital re-employed — expanded reproduction at lower
                cost.
              </dd>
              <dt>
                <code>additional_stock</code>
              </dt>
              <dd>
                Capital held as surplus inventory ("reserve capital").
              </dd>
            </dl>
          </div>
          <div className="ch15-case-block ch15-case-rise-block">
            <h4>Price Rise on MP / LP</h4>
            <dl>
              <dt>
                <code>partial_replacement</code>
              </dt>
              <dd>
                Capital advanced cannot fully replace consumed elements —
                contraction.
              </dd>
              <dt>
                <code>reduced_scale</code>
              </dt>
              <dd>Production shrinks to match the reduced real capital.</dd>
              <dt>
                <code>credit_bridge</code>
              </dt>
              <dd>Additional money-capital borrowed to bridge the gap.</dd>
            </dl>
          </div>
        </div>
      </section>
    </div>
  );
}

export default Ch15PriceChanges;
