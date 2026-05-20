import { useEffect, useState } from "react";
import type { FormEvent } from "react";
import { api } from "../../api";
import type { FarmTenure, HistoricalStage, TenantForm } from "../../types";

const TENANT_FORMS: { value: TenantForm; label: string }[] = [
  { value: "bailiff", label: "Bailiff — manages for absentee landlord" },
  { value: "metayer", label: "Metayer — shares stock and product with landlord" },
  { value: "capitalist-farmer", label: "Capitalist Farmer — advances own capital, pays rent" },
];

const tenurePresets = [
  {
    label: "15th-c. Metayer",
    form: "metayer" as TenantForm,
    lease_period_years: 1,
    rent_pence: 0,
    capital_advanced_pence: 500,
    revenue_pence: 1000,
    wage_costs_pence: 200,
  },
  {
    label: "99-yr Capitalist Lease",
    form: "capitalist-farmer" as TenantForm,
    lease_period_years: 99,
    rent_pence: 1200,
    capital_advanced_pence: 0,
    revenue_pence: 3000,
    wage_costs_pence: 600,
  },
];

export function Ch29GenesisFarmer() {
  const [stages, setStages] = useState<HistoricalStage[]>([]);
  const [selectedStageId, setSelectedStageId] = useState("");
  const [tenures, setTenures] = useState<FarmTenure[]>([]);
  const [loadingTenures, setLoadingTenures] = useState(false);
  const [tenuresError, setTenuresError] = useState("");

  const [form, setForm] = useState<TenantForm>("metayer");
  const [leaseYears, setLeaseYears] = useState("1");
  const [rentPence, setRentPence] = useState("0");
  const [capitalPence, setCapitalPence] = useState("0");
  const [revenuePence, setRevenuePence] = useState("0");
  const [wagePence, setWagePence] = useState("0");
  const [createError, setCreateError] = useState("");
  const [creating, setCreating] = useState(false);

  const [nominalRent, setNominalRent] = useState("1200");
  const [depFactor, setDepFactor] = useState("0.5");
  const [realRentResult, setRealRentResult] = useState<{ nominal: number; real: number; factor: number } | null>(null);
  const [realRentError, setRealRentError] = useState("");

  useEffect(() => {
    api.listHistoricalStages().then((list) => {
      setStages(list);
      if (list.length > 0) setSelectedStageId(list[0].id);
    });
  }, []);

  useEffect(() => {
    if (!selectedStageId) return;
    loadTenures(selectedStageId);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [selectedStageId]);

  async function loadTenures(id: string) {
    setTenuresError("");
    setLoadingTenures(true);
    try {
      setTenures(await api.listFarmTenures(id));
    } catch (err) {
      setTenuresError(err instanceof Error ? err.message : String(err));
    } finally {
      setLoadingTenures(false);
    }
  }

  function applyPreset(p: (typeof tenurePresets)[0]) {
    setForm(p.form);
    setLeaseYears(String(p.lease_period_years));
    setRentPence(String(p.rent_pence));
    setCapitalPence(String(p.capital_advanced_pence));
    setRevenuePence(String(p.revenue_pence));
    setWagePence(String(p.wage_costs_pence));
  }

  async function handleCreateTenure(e: FormEvent) {
    e.preventDefault();
    setCreateError("");
    setCreating(true);
    try {
      await api.createFarmTenure(selectedStageId, {
        form,
        lease_period_years: parseInt(leaseYears, 10),
        rent_pence: parseInt(rentPence, 10),
        capital_advanced_pence: parseInt(capitalPence, 10),
        revenue_pence: parseInt(revenuePence, 10),
        wage_costs_pence: parseInt(wagePence, 10),
      });
      await loadTenures(selectedStageId);
      setLeaseYears("1"); setRentPence("0"); setCapitalPence("0");
      setRevenuePence("0"); setWagePence("0");
    } catch (err) {
      setCreateError(err instanceof Error ? err.message : String(err));
    } finally {
      setCreating(false);
    }
  }

  async function handleComputeRealRent(e: FormEvent) {
    e.preventDefault();
    setRealRentError("");
    try {
      const result = await api.computeRealRent(
        parseInt(nominalRent, 10),
        parseFloat(depFactor),
      );
      setRealRentResult({
        nominal: result.nominal_rent_pence,
        real: result.real_rent_pence,
        factor: result.depreciation_factor,
      });
    } catch (err) {
      setRealRentError(err instanceof Error ? err.message : String(err));
    }
  }

  const cumulativeProfit = tenures.reduce((sum, t) => sum + t.surplus.profit, 0);

  return (
    <div className="chapter-panel">
      <section className="panel-section">
        <h2>Historical Stage</h2>
        <p className="muted small">
          Farm tenures are attached to a historical stage. Select the stage to view
          or record tenure arrangements.
        </p>
        {stages.length === 0 && (
          <p className="muted small">No historical stages found — create one in Ch. 26 first.</p>
        )}
        {stages.length > 0 && (
          <label>
            Stage
            <select value={selectedStageId} onChange={(e) => setSelectedStageId(e.target.value)}>
              {stages.map((s) => (
                <option key={s.id} value={s.id}>{s.name}</option>
              ))}
            </select>
          </label>
        )}
      </section>

      {selectedStageId && (
        <section className="panel-section">
          <h2>Add Farm Tenure</h2>
          <p className="muted small">
            Record a tenure arrangement — bailiff, metayer, or capitalist farmer.
            Revenue and wage costs determine the appropriated surplus.
          </p>
          <div className="preset-buttons">
            {tenurePresets.map((p) => (
              <button
                key={p.label}
                type="button"
                className="btn-secondary btn-sm"
                onClick={() => applyPreset(p)}
              >
                {p.label}
              </button>
            ))}
          </div>
          <form onSubmit={handleCreateTenure} className="form-grid">
            <label>
              Tenant Form
              <select value={form} onChange={(e) => setForm(e.target.value as TenantForm)}>
                {TENANT_FORMS.map((f) => (
                  <option key={f.value} value={f.value}>{f.label}</option>
                ))}
              </select>
            </label>
            <label>Lease Period (years)<input type="number" min="1" value={leaseYears} onChange={(e) => setLeaseYears(e.target.value)} required /></label>
            <label>Nominal Rent (½d)<input type="number" min="0" value={rentPence} onChange={(e) => setRentPence(e.target.value)} required /></label>
            <label>Capital Advanced (½d)<input type="number" min="0" value={capitalPence} onChange={(e) => setCapitalPence(e.target.value)} required /></label>
            <label>Revenue (½d)<input type="number" min="0" value={revenuePence} onChange={(e) => setRevenuePence(e.target.value)} required /></label>
            <label>Wage Costs (½d)<input type="number" min="0" value={wagePence} onChange={(e) => setWagePence(e.target.value)} required /></label>
            {createError && <p className="error-msg">{createError}</p>}
            <button type="submit" className="btn-primary" disabled={creating || !selectedStageId}>
              {creating ? "Recording…" : "Record Tenure"}
            </button>
          </form>
        </section>
      )}

      {selectedStageId && (
        <section className="panel-section">
          <h2>Tenure Records</h2>
          {tenuresError && <p className="error-msg">{tenuresError}</p>}
          {loadingTenures && <p className="muted">Loading…</p>}
          {!loadingTenures && tenures.length === 0 && (
            <p className="muted small">No tenures recorded for this stage yet.</p>
          )}
          {tenures.length > 0 && (
            <>
              <table className="data-table">
                <thead>
                  <tr>
                    <th>Form</th>
                    <th>Lease (yrs)</th>
                    <th>Rent (½d)</th>
                    <th>Revenue (½d)</th>
                    <th>Wages (½d)</th>
                    <th>Profit (½d)</th>
                  </tr>
                </thead>
                <tbody>
                  {tenures.map((t) => (
                    <tr key={t.id}>
                      <td>{t.form}</td>
                      <td>{t.lease_period_years}</td>
                      <td>{t.rent_pence}</td>
                      <td>{t.revenue_pence}</td>
                      <td>{t.wage_costs_pence}</td>
                      <td style={{ color: t.surplus.profit > 0 ? "var(--color-success, #27ae60)" : "inherit" }}>
                        {t.surplus.profit}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
              <p className="muted small" style={{ marginTop: "0.5rem" }}>
                Cumulative farmer surplus:{" "}
                <strong>{cumulativeProfit} ½d</strong> — Harrison's farmer
                accumulated "fifty or a hundred pounds" over a long lease.
              </p>
            </>
          )}
        </section>
      )}

      <section className="panel-section">
        <h2>Currency Depreciation Calculator</h2>
        <p className="muted small">
          "The progressive fall in the value of the precious metals ... lowered wages.
          A portion of the latter was now added to profits." — With a fixed nominal rent
          and depreciating currency, the farmer's real rent burden fell while corn
          prices rose. A depreciation factor of 0.5 means the currency halved in value.
        </p>
        <form onSubmit={handleComputeRealRent} className="form-grid">
          <label>
            Nominal Rent (½d)
            <input type="number" min="0" value={nominalRent} onChange={(e) => setNominalRent(e.target.value)} required />
          </label>
          <label>
            Depreciation Factor (0–1)
            <input
              type="range" min="0.01" max="1" step="0.01"
              value={depFactor}
              onChange={(e) => setDepFactor(e.target.value)}
              style={{ width: "100%" }}
            />
            <span className="muted small">{depFactor} — real rent = {Math.round(parseInt(nominalRent || "0", 10) * parseFloat(depFactor || "1"))} ½d</span>
          </label>
          {realRentError && <p className="error-msg">{realRentError}</p>}
          <button type="submit" className="btn-primary">Compute Real Rent</button>
        </form>
        {realRentResult && (
          <p className="muted small" style={{ marginTop: "0.5rem" }}>
            Nominal: <strong>{realRentResult.nominal} ½d</strong> &rarr; Real:{" "}
            <strong style={{ color: "var(--color-success, #27ae60)" }}>{realRentResult.real} ½d</strong>{" "}
            (factor {realRentResult.factor}) — rent fell by{" "}
            <strong>{Math.round((1 - realRentResult.factor) * 100)}%</strong> in real terms.
          </p>
        )}
      </section>
    </div>
  );
}
