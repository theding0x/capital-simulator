import { useEffect, useState } from "react";
import { financeApi } from "../../api";
import type {
  EconomyAnalysisResponse,
  EconomyKind,
  SharedConditionResponse,
  WasteReductionResponse,
} from "../../types";
import "./Ch05ConstantCapitalEconomy.css";

// Rates arrive in basis points (10000 bp = 100%); render as percent.
function pct(bp: number): string {
  return `${(bp / 100).toLocaleString("en-GB", { maximumFractionDigits: 2 })}%`;
}

// Capital magnitudes are LabourMinutes; Marx denominates them in £.
function money(n: number): string {
  return `£${n.toLocaleString("en-GB")}`;
}

// Client-side mirror of the Go profitRateBP: s / (c + v) in basis points, with
// truncating (floor) division to match the server to the last basis point.
function profitRateBp(s: number, c: number, v: number): number {
  const total = c + v;
  return total > 0 ? Math.floor((s * 10000) / total) : 0;
}

const KIND_LABELS: Record<EconomyKind, string> = {
  shared_conditions: "Shared conditions (boiler, building, power)",
  waste_reduction: "Waste reclaimed",
  workday_extension: "Workday extension",
  improved_machinery: "Improved machinery",
  raw_material_saving: "Raw-material saving",
};

const KINDS = Object.keys(KIND_LABELS) as EconomyKind[];

interface EconomyPreset {
  label: string;
  note: string;
  kind: EconomyKind;
  c: number;
  v: number;
  s: number;
  saving: number;
}

// The three seeded economies, drawn from Ch. 5's worked illustrations.
const PRESETS: EconomyPreset[] = [
  {
    label: "Two enterprises (§I)",
    note: "11,000c + 1,000v + 1,000s. Economising £2,000 of constant capital lifts p′ from 8⅓% to 10%.",
    kind: "raw_material_saving",
    c: 11000,
    v: 1000,
    s: 1000,
    saving: 2000,
  },
  {
    label: "Waste cotton (§IV)",
    note: "2,000c + 1,000v + 1,000s. Reclaiming 40 lb of waste at £5 (= £200) raises p′ from 33.33% to 35.71%.",
    kind: "waste_reduction",
    c: 2000,
    v: 1000,
    s: 1000,
    saving: 200,
  },
  {
    label: "Communal boiler (§IV)",
    note: "1,200c + 600v + 600s. Sharing a boiler-house with 11 others saves £100, raising p′ from 33.33% to 35.29%.",
    kind: "shared_conditions",
    c: 1200,
    v: 600,
    s: 600,
    saving: 100,
  },
];

export function Ch05ConstantCapitalEconomy() {
  const [kind, setKind] = useState<EconomyKind>("raw_material_saving");
  const [c, setC] = useState(11000);
  const [v, setV] = useState(1000);
  const [s, setS] = useState(1000);
  const [saving, setSaving] = useState(2000);

  // §IV calculators that derive a saving for two of the kinds.
  const [shared, setShared] = useState<Omit<SharedConditionResponse, "kind">>({
    solo_cost: 2400,
    shared_cost: 1200,
    participant_count: 12,
  });
  const [waste, setWaste] = useState<Omit<WasteReductionResponse, "kind">>({
    original_waste: 100,
    recovered_waste: 40,
    price_per_unit: 5,
  });

  // The economy realised in constant capital, derived per kind. shared_conditions
  // and waste_reduction compute it from their own inputs; the rest use the slider.
  const derivedSaving =
    kind === "shared_conditions"
      ? shared.participant_count > 0
        ? Math.max(0, Math.floor((shared.solo_cost - shared.shared_cost) / shared.participant_count))
        : 0
      : kind === "waste_reduction"
        ? waste.recovered_waste * waste.price_per_unit
        : saving;

  // A saving can never exceed the constant capital it economises.
  const effSaving = Math.min(derivedSaving, c);
  const oldP = profitRateBp(s, c, v);
  const newP = profitRateBp(s, c - effSaving, v);
  const gain = newP - oldP;

  const [result, setResult] = useState<EconomyAnalysisResponse | null>(null);
  const [recorded, setRecorded] = useState<EconomyAnalysisResponse[]>([]);
  const [error, setError] = useState<string | null>(null);

  async function refreshRecorded() {
    try {
      setRecorded(await financeApi.listEconomyAnalyses());
    } catch {
      // the seed list is best-effort; the panel still works without it
    }
  }

  useEffect(() => {
    refreshRecorded();
  }, []);

  function applyPreset(p: EconomyPreset) {
    setKind(p.kind);
    setC(p.c);
    setV(p.v);
    setS(p.s);
    setSaving(p.saving);
  }

  async function record() {
    setError(null);
    try {
      const a = await financeApi.createEconomyAnalysis({
        kind,
        constant: c,
        variable: v,
        surplus_value: s,
        saving: effSaving,
      });
      setResult(a);
      await refreshRecorded();
    } catch (e) {
      setError(String(e));
    }
  }

  return (
    <div className="v3-ch05">
      <section className="v3-ch05-intro">
        <p>
          The rate of profit <strong>p′ = s / (c + v)</strong> carries the constant
          capital in its denominator alone. Leave the surplus-value{" "}
          <strong>s</strong> and the variable capital <strong>v</strong> untouched,
          and every economy in <strong>c</strong> shrinks the denominator and lifts
          the rate — without a farthing of new surplus-value.
        </p>
        <p className="v3-ch05-equation">p′ = s / (c − saving + v) ↑</p>
        <p>
          Marx surveys the forms this economy takes: the same fixed capital made to
          serve a longer working-day, the conditions of production shared in common,
          and the waste of raw material reclaimed and returned to production.
        </p>
      </section>

      <section className="v3-ch05-calc" aria-label="Economy in constant capital">
        <h4>Economise constant capital, watch p′ rise</h4>

        <label className="v3-ch05-control">
          <span>Form of economy</span>
          <select value={kind} onChange={(e) => setKind(e.target.value as EconomyKind)}>
            {KINDS.map((k) => (
              <option key={k} value={k}>
                {KIND_LABELS[k]}
              </option>
            ))}
          </select>
        </label>

        <div className="v3-ch05-controls">
          <label className="v3-ch05-control">
            <span>Constant capital c — {money(c)}</span>
            <input type="range" min={0} max={20000} step={100} value={c}
              onChange={(e) => setC(Number(e.target.value))} />
          </label>
          <label className="v3-ch05-control">
            <span>Variable capital v — {money(v)}</span>
            <input type="range" min={0} max={5000} step={50} value={v}
              onChange={(e) => setV(Number(e.target.value))} />
          </label>
          <label className="v3-ch05-control">
            <span>Surplus-value s — {money(s)}</span>
            <input type="range" min={0} max={5000} step={50} value={s}
              onChange={(e) => setS(Number(e.target.value))} />
          </label>
        </div>

        {kind === "shared_conditions" ? (
          <div className="v3-ch05-sub" aria-label="Shared-condition calculator">
            <p className="v3-ch05-sub-title">
              Conditions of production shared in common — the saving is{" "}
              <em>(solo − shared) / participants</em>.
            </p>
            <div className="v3-ch05-controls">
              <label className="v3-ch05-control">
                <span>Solo cost (all alone) — {money(shared.solo_cost)}</span>
                <input type="range" min={0} max={10000} step={100} value={shared.solo_cost}
                  onChange={(e) => setShared({ ...shared, solo_cost: Number(e.target.value) })} />
              </label>
              <label className="v3-ch05-control">
                <span>Shared facility cost — {money(shared.shared_cost)}</span>
                <input type="range" min={0} max={10000} step={100} value={shared.shared_cost}
                  onChange={(e) => setShared({ ...shared, shared_cost: Number(e.target.value) })} />
              </label>
              <label className="v3-ch05-control">
                <span>Participants — {shared.participant_count}</span>
                <input type="range" min={1} max={24} step={1} value={shared.participant_count}
                  onChange={(e) => setShared({ ...shared, participant_count: Number(e.target.value) })} />
              </label>
            </div>
          </div>
        ) : kind === "waste_reduction" ? (
          <div className="v3-ch05-sub" aria-label="Waste-reduction calculator">
            <p className="v3-ch05-sub-title">
              Waste reclaimed and returned to production — the saving is{" "}
              <em>recovered × price</em>.
            </p>
            <div className="v3-ch05-controls">
              <label className="v3-ch05-control">
                <span>Waste recovered (units) — {waste.recovered_waste}</span>
                <input type="range" min={0} max={waste.original_waste} step={1} value={waste.recovered_waste}
                  onChange={(e) => setWaste({ ...waste, recovered_waste: Number(e.target.value) })} />
              </label>
              <label className="v3-ch05-control">
                <span>Price per unit — {money(waste.price_per_unit)}</span>
                <input type="range" min={0} max={20} step={1} value={waste.price_per_unit}
                  onChange={(e) => setWaste({ ...waste, price_per_unit: Number(e.target.value) })} />
              </label>
            </div>
          </div>
        ) : (
          <label className="v3-ch05-control">
            <span>Saving in constant capital — {money(saving)}</span>
            <input type="range" min={0} max={c} step={50} value={saving}
              onChange={(e) => setSaving(Number(e.target.value))} />
          </label>
        )}

        <div className="v3-ch05-readout">
          <div className="v3-ch05-stat">
            <span className="v3-ch05-stat-label">Old p′ = s/(c+v)</span>
            <span className="v3-ch05-stat-value">{pct(oldP)}</span>
          </div>
          <div className="v3-ch05-stat v3-ch05-stat-headline">
            <span className="v3-ch05-stat-label">New p′ = s/(c−saving+v)</span>
            <span className="v3-ch05-stat-value">{pct(newP)}</span>
          </div>
          <div className="v3-ch05-stat">
            <span className="v3-ch05-stat-label">Gain</span>
            <span className="v3-ch05-stat-value">+{pct(gain)}</span>
          </div>
          <div className="v3-ch05-stat">
            <span className="v3-ch05-stat-label">Economy in c</span>
            <span className="v3-ch05-stat-value">{money(effSaving)}</span>
          </div>
        </div>

        <p className="v3-ch05-calc-note">
          Economising <strong>{money(effSaving)}</strong> of constant capital, the
          surplus-value <strong>{money(s)}</strong> and the variable capital{" "}
          <strong>{money(v)}</strong> unchanged, raises the rate of profit from{" "}
          <strong>{pct(oldP)}</strong> to <strong>{pct(newP)}</strong>.
        </p>

        <button type="button" className="v3-ch05-record" onClick={record}>
          Record this economy
        </button>
      </section>

      <section className="v3-ch05-presets-section" aria-label="Marx's worked economies">
        <h4>Marx's worked economies</h4>
        <div className="v3-ch05-presets">
          {PRESETS.map((p) => (
            <button key={p.label} type="button" className="v3-ch05-preset" onClick={() => applyPreset(p)}
              title={p.note}>
              {p.label}
            </button>
          ))}
        </div>
        {result && (
          <div className="v3-ch05-result">
            <div className="v3-ch05-shift">
              <span className="v3-ch05-shift-old">{pct(result.old_profit_rate)}</span>
              <span className="v3-ch05-shift-arrow">→</span>
              <span className="v3-ch05-shift-new">{pct(result.new_profit_rate)}</span>
              <span className="v3-ch05-shift-word">p′ (+{pct(result.profit_rate_gain)})</span>
            </div>
            <p className="v3-ch05-result-detail">
              {KIND_LABELS[result.economy.kind]}: {money(result.economy.constant_capital)}c +{" "}
              {money(result.economy.variable_capital)}v + {money(result.economy.surplus_value)}s, economising{" "}
              {money(result.economy.saving)} of constant capital.
            </p>
          </div>
        )}
      </section>

      {error && <p className="v3-ch05-error">{error}</p>}

      {recorded.length > 0 && (
        <section className="v3-ch05-recorded" aria-label="Recorded economies">
          <h4>Recorded economies</h4>
          <table className="v3-ch05-table">
            <thead>
              <tr>
                <th>Kind</th>
                <th>c</th>
                <th>v</th>
                <th>s</th>
                <th>saving</th>
                <th>old p′</th>
                <th>new p′</th>
                <th>gain</th>
              </tr>
            </thead>
            <tbody>
              {recorded.map((a) => (
                <tr key={a.id}>
                  <td>{KIND_LABELS[a.economy.kind]}</td>
                  <td>{money(a.economy.constant_capital)}</td>
                  <td>{money(a.economy.variable_capital)}</td>
                  <td>{money(a.economy.surplus_value)}</td>
                  <td>{money(a.economy.saving)}</td>
                  <td>{pct(a.old_profit_rate)}</td>
                  <td>{pct(a.new_profit_rate)}</td>
                  <td>+{pct(a.profit_rate_gain)}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </section>
      )}
    </div>
  );
}
