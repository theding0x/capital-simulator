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

// Capital magnitudes are LabourMinutes-typed; Marx denominates them in £. Part I
// keeps the same integer, so one stored unit prints as £1 (this is not the
// pkg/money MELT, where £1 = 2,400 LM). Rates are basis points, so the scale
// never enters them.
function money(n: number): string {
  return `£${n.toLocaleString("en-GB")}`;
}

// Client-side mirror of the Go profitRateBP: s / (c + v) in basis points, with
// truncating (floor) division. The server rounds half-up (#410, migration 00108),
// so this can read one basis point lower when the fraction is .5 or more. None of
// the preset economies below is affected.
function profitRateBp(s: number, c: number, v: number): number {
  const total = c + v;
  return total > 0 ? Math.floor((s * 10000) / total) : 0;
}

const KIND_LABELS: Record<EconomyKind, string> = {
  shared_conditions: "Shared conditions (shaft, building, power)",
  waste_reduction: "Waste reclaimed",
  workday_extension: "Workday extension",
  improved_machinery: "Improved machinery",
  raw_material_saving: "Raw-material saving",
};

const KINDS = Object.keys(KIND_LABELS) as EconomyKind[];

type SharedInputs = Omit<SharedConditionResponse, "kind">;
type WasteInputs = Omit<WasteReductionResponse, "kind">;

// §III: the saving each participant realises from a condition used in common.
function sharedSaving(x: SharedInputs): number {
  if (x.participant_count <= 0) return 0;
  return Math.max(0, Math.floor((x.solo_cost - x.shared_cost) / x.participant_count));
}

// §IV: the value of the reclaimed waste no longer bought afresh.
function wasteSaving(x: WasteInputs): number {
  return x.recovered_waste * x.price_per_unit;
}

// A preset for one of the two kinds that derive their saving carries that
// calculator's inputs and no saving of its own, so the two cannot disagree; the
// other kinds carry the saving directly. The union makes a waste_reduction preset
// without waste inputs (or a shared_conditions one without shared inputs) a type
// error rather than a preset that silently reuses stale calculator state.
interface PresetBase {
  label: string;
  note: string;
  c: number;
  v: number;
  s: number;
}

type EconomyPreset =
  | (PresetBase & { kind: "waste_reduction"; waste: WasteInputs })
  | (PresetBase & { kind: "shared_conditions"; shared: SharedInputs })
  | (PresetBase & {
      kind: Exclude<EconomyKind, "waste_reduction" | "shared_conditions">;
      saving: number;
    });

// The rag prices Marx cites are per ton: about £4 4s. thirty years ago, £44 now. Only the
// £44 is held (whole pounds); the 10 tons are illustrative. Saving = 10 x 44 = 440.
const RAGS_WASTE: WasteInputs = { original_waste: 10, recovered_waste: 10, price_per_unit: 44 };

// The single shaft Marx quotes saves 60% of the shafting other firms need (80% in some
// cases). The £500 is illustrative, so this mill's shafting is £200 and the saving
// £300, one participant.
const SHAFT_SHARED: SharedInputs = { solo_cost: 500, shared_cost: 200, participant_count: 1 };

const RAGS: EconomyPreset = {
  label: "Woollen rags (§IV)",
  note: "Marx cites a factory report: woollen rags about £4 4s. a ton thirty years ago, £44 a ton now. Illustrative base 2,000c + 1,000v + 1,000s. Reclaiming 10 tons (= £440) raises p′ from 33.33% to 39.06%.",
  kind: "waste_reduction",
  c: 2000,
  v: 1000,
  s: 1000,
  waste: RAGS_WASTE,
};

const SHAFT: EconomyPreset = {
  label: "Single shaft (§III)",
  note: "Marx quotes a cotton spinner: one shaft drives 29,000 doubling spindles, saving 60% of the shafting other firms need (80% in some cases). Illustrative base 1,200c + 600v + 600s with £500 of shafting. The £300 saved raises p′ from 33.33% to 40%.",
  kind: "shared_conditions",
  c: 1200,
  v: 600,
  s: 600,
  shared: SHAFT_SHARED,
};

const WORKDAY = {
  label: "Working-day 12 to 18 hours (§I)",
  note: "Marx: the fixed capital is the same whether it serves 12 or 16 hours, and at 18 instead of 12 its value is reproduced in ⅔ of the time. Illustrative: 3,600 of fixed capital, 1,200v, 1,200s. The charge on the same 12 working hours falls by a third (£1,200): p′ rises from 25% to 33.33%.",
  kind: "workday_extension",
  c: 3600,
  v: 1200,
  s: 1200,
  saving: 1200,
} satisfies EconomyPreset;

const TWO_ENTERPRISES: EconomyPreset = {
  label: "Two enterprises (Ch. 7)",
  note: "Marx's Ch. 7 illustration, borrowed here (not Ch. 5): s = v = £1,000, c = £9,000 or £11,000. Economising the £2,000 between them lifts p′ from 8⅓% to 10%.",
  kind: "raw_material_saving",
  c: 11000,
  v: 1000,
  s: 1000,
  saving: 2000,
};

// The four seeded economies: Ch. 5's own three illustrations, then Ch. 7's pair.
const PRESETS: EconomyPreset[] = [RAGS, SHAFT, WORKDAY, TWO_ENTERPRISES];

export function Ch05ConstantCapitalEconomy() {
  // The panel opens on Ch. 5's own argument (the working-day), not Ch. 7's pair.
  const [kind, setKind] = useState<EconomyKind>(WORKDAY.kind);
  const [c, setC] = useState(WORKDAY.c);
  const [v, setV] = useState(WORKDAY.v);
  const [s, setS] = useState(WORKDAY.s);
  const [saving, setSaving] = useState(WORKDAY.saving);

  // §III and §IV calculators that derive a saving for two of the kinds.
  const [shared, setShared] = useState<SharedInputs>(SHAFT_SHARED);
  const [waste, setWaste] = useState<WasteInputs>(RAGS_WASTE);

  // The economy realised in constant capital, derived per kind. shared_conditions
  // and waste_reduction compute it from their own inputs; the rest use the slider.
  const derivedSaving =
    kind === "shared_conditions"
      ? sharedSaving(shared)
      : kind === "waste_reduction"
        ? wasteSaving(waste)
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
    switch (p.kind) {
      case "waste_reduction":
        setWaste(p.waste);
        setSaving(wasteSaving(p.waste));
        break;
      case "shared_conditions":
        setShared(p.shared);
        setSaving(sharedSaving(p.shared));
        break;
      default:
        setSaving(p.saving);
    }
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
                <span>Waste recovered (tons) — {waste.recovered_waste}</span>
                <input type="range" min={0} max={waste.original_waste} step={1} value={waste.recovered_waste}
                  onChange={(e) => setWaste({ ...waste, recovered_waste: Number(e.target.value) })} />
              </label>
              <label className="v3-ch05-control">
                <span>Price per ton — {money(waste.price_per_unit)}</span>
                <input type="range" min={0} max={60} step={1} value={waste.price_per_unit}
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

      <section className="v3-ch05-presets-section" aria-label="Worked economies">
        <h4>Worked economies</h4>
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
