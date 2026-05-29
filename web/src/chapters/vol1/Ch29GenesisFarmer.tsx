import { useCallback, useEffect, useMemo, useState } from "react";
import type { FormEvent } from "react";
import { api } from "../../api";
import type {
  FarmTenure,
  FarmTenureInput,
  HistoricalStage,
  RealRentResult,
  TenantForm,
} from "../../types";
import "./Ch29GenesisFarmer.css";

interface TenurePreset {
  id: string;
  label: string;
  input: FarmTenureInput;
}

const TENURE_PRESETS: TenurePreset[] = [
  {
    id: "bailiff-1300",
    label: "Bailiff (≈1300)",
    input: {
      form: "bailiff",
      lease_period_years: 1,
      rent_pence: 0,
      capital_advanced_pence: 0,
      revenue_pence: 400,
      wage_costs_pence: 320,
    },
  },
  {
    id: "metayer-1450",
    label: "Métayer (1450)",
    input: {
      form: "metayer",
      lease_period_years: 7,
      rent_pence: 600,
      capital_advanced_pence: 1200,
      revenue_pence: 2000,
      wage_costs_pence: 800,
    },
  },
  {
    id: "capitalist-1650",
    label: "Capitalist Farmer (1650)",
    input: {
      form: "capitalist-farmer",
      lease_period_years: 99,
      rent_pence: 2000,
      capital_advanced_pence: 50000,
      revenue_pence: 60000,
      wage_costs_pence: 12000,
    },
  },
];

const TENURE_FORMS: { value: TenantForm; label: string }[] = [
  { value: "bailiff", label: "Bailiff" },
  { value: "metayer", label: "Métayer" },
  { value: "capitalist-farmer", label: "Capitalist farmer" },
];

function fmtPence(p: number): string {
  const pounds = p / 100;
  return `£${pounds.toLocaleString(undefined, { minimumFractionDigits: 0, maximumFractionDigits: 2 })}`;
}

function Stat({
  label,
  value,
  tone,
}: {
  label: string;
  value: string;
  tone?: "red" | "gold" | "lead" | "muted";
}) {
  const color =
    tone === "red"   ? "var(--red)" :
    tone === "gold"  ? "var(--gold-bright)" :
    tone === "lead"  ? "var(--lead-hover)" :
    tone === "muted" ? "var(--ink-muted)" :
    "var(--ink)";
  return (
    <div className="ch32-stat">
      <div className="ch32-stat-label">{label}</div>
      <div className="ch32-stat-value mono" style={{ color }}>{value}</div>
    </div>
  );
}

const STEPS: {
  era: string;
  figure: string;
  desc: string;
  farmer: number;
  labour: number;
  final?: boolean;
}[] = [
  {
    era: "≈ 1300",
    figure: "Bailiff",
    desc: "Holder of office under the lord. Paid in kind from the demesne; tied to the manor; no surplus of his own to dispose of.",
    farmer: 0.08,
    labour: 0.95,
  },
  {
    era: "1400s",
    figure: "Métayer · half-share tenant",
    desc: "Lord supplies seed, cattle, implements; tenant supplies labour and management; product divided. The tenant begins to retain a residual.",
    farmer: 0.20,
    labour: 0.78,
  },
  {
    era: "1450–1550",
    figure: "Yeoman · freeholder",
    desc: "Owner-occupier of a small estate. Holds independent of the lord, but exposed to enclosure, taxation and the price revolution.",
    farmer: 0.38,
    labour: 0.58,
  },
  {
    era: "1550–1650",
    figure: "Lease-holder of the Tudor type",
    desc: "Long lease at a customary money rent. Through the price revolution the real value of the rent collapses; the tenant pockets the difference. The agricultural revolution is in part a windfall.",
    farmer: 0.62,
    labour: 0.34,
  },
  {
    era: "1650–1750",
    figure: "Capitalist farmer",
    desc: "Employs wage-labour on a multi-hundred-acre lease. Pays a rack-rent, but the rate of surplus-value extracted from the labourers leaves a profit characteristic of the new mode of production.",
    farmer: 0.90,
    labour: 0.12,
    final: true,
  },
];

export function Ch29GenesisFarmer() {
  const [wage1500,  setWage1500]  = useState(4);
  const [wage1750,  setWage1750]  = useState(8);
  const [price1500, setPrice1500] = useState(100);
  const [price1750, setPrice1750] = useState(420);
  const [rentReal,  setRentReal]  = useState(0.35);

  const wageRise   = wage1750 / Math.max(wage1500, 1);
  const priceRise  = price1750 / Math.max(price1500, 1);
  const realWage   = wageRise / priceRise;
  const fatten     = (1 - realWage) * 100;
  const farmerGain = (priceRise - 1) * rentReal * 100;

  return (
    <>
      <DepreciationInsight />
      <FarmTenuresPanel />
      <RealRentPanel />
      <section className="card">
        <h2>The Twin Movement</h2>
        <p className="description">
          A class is not born ex nihilo. The capitalist farmer is the bailiff transformed by three
          centuries of intervening events: the price revolution, the enclosures, the dispossession
          of the labourer from the soil. Each enriched him, and the same events impoverished the
          people he came to employ.
        </p>
        <ol className="ch29-genealogy">
          {STEPS.map((s) => (
            <li key={s.era} className={"ch29-step" + (s.final ? " ch29-step--final" : "")}>
              <div className="ch29-step-era">{s.era}</div>
              <div className="ch29-step-figure">{s.figure}</div>
              <div>
                <p className="ch29-step-note">{s.desc}</p>
                <div className="ch29-wealth-bar">
                  <div className="ch29-wealth-cell ch29-wealth-cell--gain">
                    <span style={{ width: "5.5rem" }}>farmer&#8217;s share</span>
                    <div className="ch29-wealth-meter">
                      <span style={{ width: (s.farmer * 100) + "%" }} />
                    </div>
                    <span className="mono" style={{ width: "2.5rem", textAlign: "right" }}>
                      {Math.round(s.farmer * 100)}%
                    </span>
                  </div>
                  <div className="ch29-wealth-cell ch29-wealth-cell--loss">
                    <span style={{ width: "5.5rem" }}>labour&#8217;s claim</span>
                    <div className="ch29-wealth-meter">
                      <span style={{ width: (s.labour * 100) + "%" }} />
                    </div>
                    <span className="mono" style={{ width: "2.5rem", textAlign: "right" }}>
                      {Math.round(s.labour * 100)}%
                    </span>
                  </div>
                </div>
              </div>
            </li>
          ))}
        </ol>
      </section>

      <section className="card">
        <h2>The Price Revolution as Windfall</h2>
        <p className="description">
          Long leases at fixed money rents, struck before the 16th-century inflation, left tenants
          paying historic prices for produce sold at current ones. The chapter cites Thorold
          Rogers&#8217;s wage and price series; the inputs below default to those figures.
        </p>
        <form className="form-grid" onSubmit={(e) => e.preventDefault()}>
          <label>
            <span>Daily money wage · 1500 (pence)</span>
            <input type="number" min={1} value={wage1500} onChange={(e) => setWage1500(+e.target.value)} />
          </label>
          <label>
            <span>Daily money wage · 1750 (pence)</span>
            <input type="number" min={1} value={wage1750} onChange={(e) => setWage1750(+e.target.value)} />
          </label>
          <label>
            <span>Food-price index · 1500</span>
            <input type="number" min={1} value={price1500} onChange={(e) => setPrice1500(+e.target.value)} />
          </label>
          <label>
            <span>Food-price index · 1750</span>
            <input type="number" min={1} value={price1750} onChange={(e) => setPrice1750(+e.target.value)} />
          </label>
          <label className="span2">
            <span>Share of price-rise NOT captured by rent</span>
            <input
              type="number"
              step={0.05}
              min={0}
              max={1}
              value={rentReal}
              onChange={(e) => setRentReal(+e.target.value)}
            />
          </label>
        </form>

        <div className="ch32-stats" style={{ marginTop: "1.5rem" }}>
          <Stat label="Money-wage rise"      value={`× ${wageRise.toFixed(2)}`} />
          <Stat label="Food-price rise"      value={`× ${priceRise.toFixed(2)}`} tone="red" />
          <Stat label="Real wage (1500 = 1)" value={realWage.toFixed(2)} tone="red" />
          <Stat label="Tenant windfall"      value={`+${farmerGain.toFixed(0)}%`} tone="gold" />
        </div>

        <p className="note" style={{ marginTop: "1.25rem", maxWidth: "62ch" }}>
          On these inputs the real wage of the agricultural labourer falls to
          <span className="mono"> {realWage.toFixed(2)} </span>
          of its 1500 level &#8212; a loss of
          <span className="mono"> {fatten.toFixed(0)}% </span>
          &#8212; while the difference is transferred, at the going lease, to the tenant farmer.
          The capitalist farmer is, to a first approximation, the difference between two price
          series.
        </p>
      </section>

      <section className="card">
        <h2>Comparative Wealth</h2>
        <p className="description">
          The same data, drawn as a single trajectory. The farmer&#8217;s share is what remains
          after fixed rent and wages have been paid out of the product; the labourer&#8217;s share
          is the day-wage purchasing power.
        </p>
        <table className="commodity-table">
          <thead>
            <tr>
              <th>Era</th>
              <th>Figure</th>
              <th style={{ textAlign: "right" }}>Farmer&#8217;s share</th>
              <th style={{ textAlign: "right" }}>Labour&#8217;s claim</th>
              <th>Mechanism in play</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td className="snlt-cell">1300</td>
              <td>Bailiff</td>
              <td className="snlt-cell" style={{ textAlign: "right" }}>8%</td>
              <td className="snlt-cell" style={{ textAlign: "right" }}>95%</td>
              <td className="muted small">manorial obligation</td>
            </tr>
            <tr>
              <td className="snlt-cell">1400s</td>
              <td>M&#233;tayer</td>
              <td className="snlt-cell" style={{ textAlign: "right" }}>20%</td>
              <td className="snlt-cell" style={{ textAlign: "right" }}>78%</td>
              <td className="muted small">share-cropping</td>
            </tr>
            <tr>
              <td className="snlt-cell">1450–1550</td>
              <td>Yeoman</td>
              <td className="snlt-cell" style={{ textAlign: "right" }}>38%</td>
              <td className="snlt-cell" style={{ textAlign: "right" }}>58%</td>
              <td className="muted small">freehold · vulnerable</td>
            </tr>
            <tr>
              <td className="snlt-cell">1550–1650</td>
              <td>Lease-holder</td>
              <td className="snlt-cell" style={{ textAlign: "right" }}>62%</td>
              <td className="snlt-cell" style={{ textAlign: "right" }}>34%</td>
              <td className="muted small">price-revolution windfall</td>
            </tr>
            <tr>
              <td className="snlt-cell">1650–1750</td>
              <td style={{ color: "var(--gold-bright)" }}>Capitalist farmer</td>
              <td className="snlt-cell" style={{ textAlign: "right", color: "var(--gold-bright)" }}>90%</td>
              <td className="snlt-cell" style={{ textAlign: "right", color: "var(--red)" }}>12%</td>
              <td className="muted small">wage-labour at rack-rent</td>
            </tr>
          </tbody>
        </table>
      </section>
      <HarrisonCoda />
    </>
  );
}

function DepreciationInsight() {
  return (
    <section className="v1-ch29-insight">
      <h2 className="v1-ch29-insight-h2">Long leases against a depreciating currency</h2>
      <p className="v1-ch29-insight-prose">
        Between 1500 and 1650 the silver content of the English coinage fell and
        New-World bullion drove a sustained rise in commodity prices. A tenant
        farmer who had locked in a nominal money rent on a long lease watched
        his real obligation evaporate while his harvest sold for ever more.
        The <em>Real-Rent</em> probe below deflates a nominal rent by a chosen
        depreciation factor; the <em>Recorded Tenures</em> panel ties tenure
        type to its computed profit on the live backend.
      </p>
    </section>
  );
}

function FarmTenuresPanel() {
  const [stages, setStages] = useState<HistoricalStage[]>([]);
  const [selectedStage, setSelectedStage] = useState<string>("");
  const [tenures, setTenures] = useState<FarmTenure[]>([]);
  const [draft, setDraft] = useState<FarmTenureInput>(TENURE_PRESETS[2].input);
  const [err, setErr] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    (async () => {
      try {
        const list = await api.listHistoricalStages();
        setStages(list);
        if (list.length > 0 && !selectedStage) {
          setSelectedStage(list[0].id);
        }
      } catch (e) {
        setErr(e instanceof Error ? e.message : String(e));
      }
    })();
  }, [selectedStage]);

  const refreshTenures = useCallback(async (stageId: string) => {
    if (!stageId) return;
    setLoading(true);
    try {
      const list = await api.listFarmTenures(stageId);
      setTenures(list);
      setErr(null);
    } catch (e) {
      setErr(e instanceof Error ? e.message : String(e));
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    void refreshTenures(selectedStage);
  }, [selectedStage, refreshTenures]);

  async function submit(e: FormEvent) {
    e.preventDefault();
    if (!selectedStage) {
      setErr("Pick a historical stage first.");
      return;
    }
    setErr(null);
    try {
      await api.createFarmTenure(selectedStage, draft);
      await refreshTenures(selectedStage);
    } catch (e) {
      setErr(e instanceof Error ? e.message : String(e));
    }
  }

  const finalForm: TenantForm = "capitalist-farmer";

  return (
    <section className="card">
      <h2>Recorded Tenures</h2>
      <p className="description">
        Persist farm-tenure records against a chosen historical stage. The
        backend computes the surplus
        (<code>profit = revenue − rent − wages</code>); each new tenure
        appears as a card with its form and profit. The progression climbs as
        the form shifts from bailiff to capitalist farmer.
      </p>

      <div className="v1-ch29-presets">
        <span className="v1-ch29-presets-label">Stage</span>
        <select
          value={selectedStage}
          onChange={(e) => setSelectedStage(e.target.value)}
        >
          <option value="">— pick a historical stage —</option>
          {stages.map((s) => (
            <option key={s.id} value={s.id}>{s.name}</option>
          ))}
        </select>
        <span className="v1-ch29-presets-label" style={{ marginLeft: "0.4rem" }}>
          Preset
        </span>
        {TENURE_PRESETS.map((p) => (
          <button
            key={p.id}
            type="button"
            className="v1-ch29-preset-button"
            onClick={() => setDraft(p.input)}
          >
            {p.label}
          </button>
        ))}
      </div>

      <form className="form-grid" onSubmit={submit}>
        <label>
          <span>Form</span>
          <select
            value={draft.form}
            onChange={(e) =>
              setDraft({ ...draft, form: e.target.value as TenantForm })
            }
          >
            {TENURE_FORMS.map((f) => (
              <option key={f.value} value={f.value}>{f.label}</option>
            ))}
          </select>
        </label>
        <label>
          <span>Lease period (years)</span>
          <input
            type="number"
            min={1}
            value={draft.lease_period_years}
            onChange={(e) =>
              setDraft({ ...draft, lease_period_years: Number(e.target.value) })
            }
          />
        </label>
        <label>
          <span>Rent (pence)</span>
          <input
            type="number"
            min={0}
            value={draft.rent_pence}
            onChange={(e) =>
              setDraft({ ...draft, rent_pence: Number(e.target.value) })
            }
          />
        </label>
        <label>
          <span>Capital advanced (pence)</span>
          <input
            type="number"
            min={0}
            value={draft.capital_advanced_pence}
            onChange={(e) =>
              setDraft({
                ...draft,
                capital_advanced_pence: Number(e.target.value),
              })
            }
          />
        </label>
        <label>
          <span>Revenue (pence)</span>
          <input
            type="number"
            min={0}
            value={draft.revenue_pence}
            onChange={(e) =>
              setDraft({ ...draft, revenue_pence: Number(e.target.value) })
            }
          />
        </label>
        <label>
          <span>Wage costs (pence)</span>
          <input
            type="number"
            min={0}
            value={draft.wage_costs_pence}
            onChange={(e) =>
              setDraft({ ...draft, wage_costs_pence: Number(e.target.value) })
            }
          />
        </label>
        <div className="form-actions span2">
          <button type="submit" className="primary" disabled={!selectedStage}>
            Record tenure
          </button>
          {err && <span className="error">{err}</span>}
        </div>
      </form>

      {loading && <p className="small muted">Loading tenures…</p>}
      {!loading && tenures.length === 0 && selectedStage && (
        <p className="small muted">No tenures recorded yet on this stage.</p>
      )}
      {tenures.length > 0 && (
        <div className="v1-ch29-tenures">
          {tenures.map((t) => {
            const isFinal = t.form === finalForm;
            const negative = t.surplus.profit < 0;
            return (
              <div
                key={t.id}
                className={
                  "v1-ch29-tenure" +
                  (isFinal ? " v1-ch29-tenure--final" : "")
                }
              >
                <span className="v1-ch29-tenure-tag">{t.form}</span>
                <h3 className="v1-ch29-tenure-h3">
                  {t.lease_period_years}-yr lease
                </h3>
                <span className="v1-ch29-tenure-meta">
                  rent {fmtPence(t.rent_pence)} · capital {fmtPence(t.capital_advanced_pence)} ·
                  revenue {fmtPence(t.revenue_pence)} · wages {fmtPence(t.wage_costs_pence)}
                </span>
                <span
                  className={
                    "v1-ch29-tenure-profit" +
                    (negative ? " v1-ch29-tenure-profit--negative" : "")
                  }
                >
                  profit {fmtPence(t.surplus.profit)}
                </span>
              </div>
            );
          })}
        </div>
      )}
    </section>
  );
}

function RealRentPanel() {
  const [nominal, setNominal] = useState(2400);
  const [factor, setFactor] = useState(0.45);
  const [result, setResult] = useState<RealRentResult | null>(null);
  const [err, setErr] = useState<string | null>(null);

  const recompute = useCallback(async () => {
    setErr(null);
    try {
      const r = await api.computeRealRent(nominal, factor);
      setResult(r);
    } catch (e) {
      setErr(e instanceof Error ? e.message : String(e));
      setResult(null);
    }
  }, [nominal, factor]);

  useEffect(() => {
    void recompute();
  }, [recompute]);

  const windfallPence = useMemo(
    () => (result ? result.nominal_rent_pence - result.real_rent_pence : 0),
    [result],
  );

  return (
    <section className="card">
      <h2>Real-Rent Probe (§1 — currency depreciation)</h2>
      <p className="description">
        Deflates a nominal rent by a depreciation factor in (0, 1]. A factor of
        0.5 means the currency halved in real value; the tenant's effective
        rent burden is half the nominal figure, and the difference is windfall.
      </p>
      <form
        className="form-grid"
        onSubmit={(e) => {
          e.preventDefault();
          void recompute();
        }}
      >
        <label>
          <span>Nominal rent (pence)</span>
          <input
            type="number"
            min={0}
            value={nominal}
            onChange={(e) => setNominal(Number(e.target.value))}
          />
        </label>
        <label>
          <span>Depreciation factor (0–1)</span>
          <input
            type="number"
            min={0.01}
            max={1}
            step={0.05}
            value={factor}
            onChange={(e) => setFactor(Number(e.target.value))}
          />
        </label>
      </form>
      {err && <p className="error">{err}</p>}
      {result && (
        <div className="v1-ch29-rent-chip">
          <span className="v1-ch29-rent-chip-label">Real rent</span>
          <span className="v1-ch29-rent-chip-value">
            {fmtPence(result.real_rent_pence)}
          </span>
          <span className="v1-ch29-rent-chip-meta">
            · windfall {fmtPence(windfallPence)} (×{result.depreciation_factor.toFixed(2)})
          </span>
        </div>
      )}
    </section>
  );
}

function HarrisonCoda() {
  return (
    <aside className="v1-ch29-coda">
      <p className="v1-ch29-coda-quote">
        “Whereas in old time, a man would have been thought well to do that
        could keep £50 or £100 in store, now in our days £200 or £300 is
        scarcely accounted a great matter.”
        <span className="v1-ch29-coda-cite">
          — Harrison, quoted by Marx, Capital Vol. I, Ch. 29
        </span>
      </p>
    </aside>
  );
}
