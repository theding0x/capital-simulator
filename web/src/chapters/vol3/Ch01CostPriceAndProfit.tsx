import { useEffect, useState } from "react";
import { financeApi } from "../../api";
import type {
  CostPriceResponse,
  ProfitFormResponse,
} from "../../types";
import "./Ch01CostPriceAndProfit.css";

// Magnitudes are LabourMinutes; Marx denominates them in £ (6s. = one social
// working-day), so we render them as whole pounds, matching his worked example.
function money(n: number): string {
  return `£${n.toLocaleString("en-GB")}`;
}

interface OutlayForm {
  constant: number;
  variable: number;
  fixedWear: number;
  fixedAdvanced: number;
  surplus: number;
}

interface Preset {
  label: string;
  values: OutlayForm;
}

// Capital III, Ch. 1, §1 — Marx's worked examples.
const PRESETS: Preset[] = [
  {
    label: "§1 — 400c + 100v + 100s",
    values: { constant: 400, variable: 100, fixedWear: 0, fixedAdvanced: 0, surplus: 100 },
  },
  {
    label: "§1 — £1,200 fixed, £20 wear",
    values: { constant: 400, variable: 100, fixedWear: 20, fixedAdvanced: 1200, surplus: 100 },
  },
];

export function Ch01CostPriceAndProfit() {
  const [form, setForm] = useState<OutlayForm>(PRESETS[1].values);
  const [costPrice, setCostPrice] = useState<CostPriceResponse | null>(null);
  const [profitForm, setProfitForm] = useState<ProfitFormResponse | null>(null);
  const [saved, setSaved] = useState<CostPriceResponse[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [busy, setBusy] = useState(false);

  async function refreshSaved() {
    try {
      setSaved(await financeApi.listCostPrices());
    } catch {
      // the seed list is best-effort; the form still works without it
    }
  }

  useEffect(() => {
    refreshSaved();
  }, []);

  async function compute() {
    setError(null);
    setBusy(true);
    try {
      const cp = await financeApi.computeCostPrice({
        constant: form.constant,
        variable: form.variable,
        fixed_wear_and_tear: form.fixedWear,
        fixed_advanced: form.fixedAdvanced,
      });
      setCostPrice(cp);
      const pf = await financeApi.getProfitForm({
        cost_price: cp.k,
        surplus_value: form.surplus,
      });
      setProfitForm(pf);
      await refreshSaved();
    } catch (e) {
      setError(String(e));
    } finally {
      setBusy(false);
    }
  }

  const fixedPct =
    costPrice && costPrice.k > 0 ? (costPrice.fixed_component / costPrice.k) * 100 : 0;

  return (
    <div className="v3-ch01">
      <section className="v3-ch01-intro">
        <p>
          Of a commodity's value <strong>C = c + v + s</strong>, the part that
          merely replaces the consumed capital — <strong>c + v</strong> — is what
          the commodity costs the <em>capitalist</em>. Marx names it the
          cost-price <strong>k</strong>, so the value formula becomes{" "}
          <strong>C = k + s</strong>. Re-attributed from variable capital to the
          whole outlay, the surplus-value then appears as <strong>profit</strong>:{" "}
          <strong>C = k + p</strong>.
        </p>
      </section>

      <section className="v3-ch01-form" aria-label="Capital outlay">
        <div className="v3-ch01-presets">
          {PRESETS.map((p) => (
            <button
              key={p.label}
              type="button"
              className="v3-ch01-preset"
              onClick={() => setForm(p.values)}
            >
              {p.label}
            </button>
          ))}
        </div>

        <div className="v3-ch01-fields">
          <label className="v3-ch01-field">
            Constant capital c
            <input
              type="number"
              min={0}
              value={form.constant}
              onChange={(e) => setForm({ ...form, constant: Number(e.target.value) })}
            />
          </label>
          <label className="v3-ch01-field">
            Variable capital v
            <input
              type="number"
              min={0}
              value={form.variable}
              onChange={(e) => setForm({ ...form, variable: Number(e.target.value) })}
            />
          </label>
          <label className="v3-ch01-field">
            of which fixed wear-and-tear
            <input
              type="number"
              min={0}
              value={form.fixedWear}
              onChange={(e) => setForm({ ...form, fixedWear: Number(e.target.value) })}
            />
          </label>
          <label className="v3-ch01-field">
            fixed capital advanced
            <input
              type="number"
              min={0}
              value={form.fixedAdvanced}
              onChange={(e) => setForm({ ...form, fixedAdvanced: Number(e.target.value) })}
            />
          </label>
          <label className="v3-ch01-field">
            Surplus-value s
            <input
              type="number"
              min={0}
              value={form.surplus}
              onChange={(e) => setForm({ ...form, surplus: Number(e.target.value) })}
            />
          </label>
        </div>

        <button type="button" className="v3-ch01-compute" onClick={compute} disabled={busy}>
          {busy ? "Computing…" : "Compute cost-price & profit"}
        </button>
        {error && <p className="v3-ch01-error">{error}</p>}
      </section>

      {costPrice && (
        <section className="v3-ch01-result" aria-label="Cost-price">
          <h3>
            Cost-price k = c + v = {money(costPrice.k)}
          </h3>
          <div
            className="v3-ch01-bar"
            role="img"
            aria-label={`Cost-price split: fixed component ${costPrice.fixed_component}, circulating component ${costPrice.circulating_component}`}
          >
            {costPrice.fixed_component > 0 && (
              <div
                className="v3-ch01-bar-fixed"
                style={{ width: `${fixedPct}%` }}
                title={`Fixed component (depreciation): ${money(costPrice.fixed_component)}`}
              >
                fixed {money(costPrice.fixed_component)}
              </div>
            )}
            <div
              className="v3-ch01-bar-circulating"
              style={{ width: `${100 - fixedPct}%` }}
              title={`Circulating component: ${money(costPrice.circulating_component)}`}
            >
              circulating {money(costPrice.circulating_component)}
            </div>
          </div>
          {costPrice.outlay.fixed_advanced > 0 && (
            <p className="v3-ch01-note">
              Only {money(costPrice.fixed_component)} of the{" "}
              {money(costPrice.outlay.fixed_advanced)} fixed capital advanced
              enters the cost-price — the rest remains as productive capital.
            </p>
          )}
        </section>
      )}

      {profitForm && (
        <section className="v3-ch01-result" aria-label="Profit-form">
          <div className="v3-ch01-formula">
            <span>C = k + s</span>
            <span className="v3-ch01-arrow">→</span>
            <span>C = k + p</span>
            <span className="v3-ch01-eq">
              {money(profitForm.cost_price)} + {money(profitForm.profit)} ={" "}
              {money(profitForm.commodity_value)}
            </span>
          </div>
          <p className="v3-ch01-mystify">
            Profit <strong>p = {money(profitForm.profit)}</strong> equals the
            surplus-value <strong>s = {money(profitForm.surplus_value)}</strong>,
            but it now appears as an offspring of the <em>total</em> capital
            advanced rather than of variable capital alone.{" "}
            {profitForm.mystifies_origin && (
              <span className="v3-ch01-flag">This is the mystification.</span>
            )}
          </p>

          {profitForm.selling_price_scenarios.length > 0 && (
            <>
              <h4>A profit even below value</h4>
              <p className="v3-ch01-note">
                The minimal limit of the selling-price is the cost-price{" "}
                {money(profitForm.cost_price)}. Any price above it yields a
                profit — including prices below the value{" "}
                {money(profitForm.commodity_value)}.
              </p>
              <table className="v3-ch01-table">
                <thead>
                  <tr>
                    <th>Selling price</th>
                    <th>Profit</th>
                    <th>Relation to value</th>
                  </tr>
                </thead>
                <tbody>
                  {profitForm.selling_price_scenarios.map((s) => (
                    <tr key={s.selling_price}>
                      <td>{money(s.selling_price)}</td>
                      <td>{money(s.amount)}</td>
                      <td>
                        {s.below_value ? (
                          <span className="v3-ch01-below">below value</span>
                        ) : (
                          <span className="v3-ch01-atvalue">at value</span>
                        )}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </>
          )}
        </section>
      )}

      {saved.length > 0 && (
        <section className="v3-ch01-saved" aria-label="Saved cost-prices">
          <h4>Recorded cost-prices</h4>
          <table className="v3-ch01-table">
            <thead>
              <tr>
                <th>c</th>
                <th>v</th>
                <th>fixed</th>
                <th>circulating</th>
                <th>k</th>
              </tr>
            </thead>
            <tbody>
              {saved.map((cp) => (
                <tr key={cp.id}>
                  <td>{money(cp.outlay.constant)}</td>
                  <td>{money(cp.outlay.variable)}</td>
                  <td>{money(cp.fixed_component)}</td>
                  <td>{money(cp.circulating_component)}</td>
                  <td>{money(cp.k)}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </section>
      )}
    </div>
  );
}
