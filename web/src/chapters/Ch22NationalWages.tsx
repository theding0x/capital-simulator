import { useCallback, useEffect, useState } from "react";
import { api } from "../api";
import type { WageComparison } from "../types";
import { useCurrency, usePounds } from "../CurrencyContext";
import "./Ch22NationalWages.css";

export function Ch22NationalWages() {
  const { modern } = useCurrency();
  const fmt = usePounds();
  const [comparison, setComparison] = useState<WageComparison | null>(null);
  const [loadError, setLoadError] = useState<string | null>(null);

  // intensity form
  const [intensityCountry, setIntensityCountry] = useState("GB");
  const [intensityFactor, setIntensityFactor] = useState(2.0);
  const [intensityError, setIntensityError] = useState<string | null>(null);

  // day-wage form
  const [wageCountry, setWageCountry] = useState("FR");
  const [wagePence, setWagePence] = useState(36);
  const [wageMinutes, setWageMinutes] = useState(720);
  const [wageError, setWageError] = useState<string | null>(null);

  const loadComparison = useCallback(async () => {
    try {
      const cmp = await api.getWageComparison(600);
      setComparison(cmp);
      setLoadError(null);
    } catch (err) {
      setLoadError(err instanceof Error ? err.message : String(err));
    }
  }, []);

  useEffect(() => { loadComparison(); }, [loadComparison]);

  async function handleRegisterIntensity(e: React.FormEvent) {
    e.preventDefault();
    setIntensityError(null);
    try {
      await api.registerIntensity({ country_code: intensityCountry.toUpperCase(), factor: intensityFactor });
      loadComparison();
    } catch (err) {
      setIntensityError(err instanceof Error ? err.message : String(err));
    }
  }

  async function handleRegisterDayWage(e: React.FormEvent) {
    e.preventDefault();
    setWageError(null);
    try {
      await api.registerDayWage({
        country_code: wageCountry.toUpperCase(),
        nominal_pence: wagePence,
        working_day_minutes: wageMinutes,
      });
      loadComparison();
    } catch (err) {
      setWageError(err instanceof Error ? err.message : String(err));
    }
  }

  // Build lookup maps for the table
  const stdMap = new Map(comparison?.standardised_wages.map((s) => [s.country_code, s.amount]) ?? []);
  const relMap = new Map(comparison?.relative_prices.map((r) => [r.country_code, r.ratio]) ?? []);
  const spindleMap = new Map(comparison?.spindle_ratios.map((sr) => [sr.country_code, sr.spindles_per_worker]) ?? []);

  return (
    <div className="ch22-national-wages">
      <section className="ch22-section">
        <h2>Cross-National Wage Comparison</h2>
        <p className="ch22-explainer">
          Nominal wages differ across nations, but comparing them requires first reducing to a
          uniform working day, then adjusting for national intensity and productivity. England's
          higher nominal wage paradoxically represents a <em>lower</em> relative labour price
          to the capitalist — "wages are virtually lower to the capitalist, though higher to the
          operative than on the Continent" (Cowell, 1833; Redgrave, 1866).
        </p>

        {loadError && <p className="ch22-error">{loadError}</p>}

        {comparison && comparison.day_wages.length > 0 ? (
          <div className="ch22-table-wrap">
            <table className="ch22-table">
              <thead>
                <tr>
                  <th>Country</th>
                  <th>{modern ? "Nominal Wage (2025 £)" : "Nominal Wage (d.)"}</th>
                  <th>Std. Wage (600 min)</th>
                  <th>Relative Price</th>
                  <th>Spindles/Worker</th>
                </tr>
              </thead>
              <tbody>
                {comparison.day_wages.map((dw) => {
                  const isGB = dw.country_code === "GB";
                  const std = stdMap.get(dw.country_code);
                  const rel = relMap.get(dw.country_code);
                  const spindles = spindleMap.get(dw.country_code);
                  return (
                    <tr key={dw.country_code} className={isGB ? "ch22-row-gb" : undefined}>
                      <td>
                        {dw.country_code}
                        {isGB && <span className="ch22-paradox-tag"> &#9733; paradox</span>}
                      </td>
                      <td>{modern ? fmt(dw.nominal_pence) : `${dw.nominal_pence}d.`}</td>
                      <td>{std !== undefined ? (modern ? fmt(std) : `${std}d.`) : "—"}</td>
                      <td>{rel !== undefined ? rel.toFixed(3) : "—"}</td>
                      <td>{spindles !== undefined ? spindles : "—"}</td>
                    </tr>
                  );
                })}
              </tbody>
            </table>
          </div>
        ) : (
          !loadError && <p className="ch22-empty">No wage data registered yet. Use the forms below.</p>
        )}
      </section>

      <section className="ch22-section ch22-forms-row">
        <div className="ch22-form-card">
          <h3>Register National Intensity</h3>
          <form onSubmit={handleRegisterIntensity} className="ch22-form">
            <label>
              Country code (ISO)
              <input
                type="text"
                maxLength={2}
                value={intensityCountry}
                onChange={(e) => setIntensityCountry(e.target.value)}
                required
              />
            </label>
            <label>
              Intensity factor
              <input
                type="number"
                step="0.01"
                min={0.01}
                value={intensityFactor}
                onChange={(e) => setIntensityFactor(Number(e.target.value))}
                required
              />
            </label>
            {intensityError && <p className="ch22-error">{intensityError}</p>}
            <button type="submit">Register</button>
          </form>
        </div>

        <div className="ch22-form-card">
          <h3>Register Day Wage</h3>
          <form onSubmit={handleRegisterDayWage} className="ch22-form">
            <label>
              Country code (ISO)
              <input
                type="text"
                maxLength={2}
                value={wageCountry}
                onChange={(e) => setWageCountry(e.target.value)}
                required
              />
            </label>
            <label>
              Nominal wage (pence)
              <input
                type="number"
                min={1}
                value={wagePence}
                onChange={(e) => setWagePence(Number(e.target.value))}
                required
              />
            </label>
            <label>
              Working day (minutes)
              <input
                type="number"
                min={1}
                value={wageMinutes}
                onChange={(e) => setWageMinutes(Number(e.target.value))}
                required
              />
            </label>
            {wageError && <p className="ch22-error">{wageError}</p>}
            <button type="submit">Register</button>
          </form>
        </div>
      </section>
    </div>
  );
}
