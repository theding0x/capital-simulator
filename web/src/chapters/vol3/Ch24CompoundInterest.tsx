import { useEffect, useState } from "react";
import { financeApi } from "../../api";
import type { CompoundInterestSchedule } from "../../types";
import "./Ch24CompoundInterest.css";

// FetishForm constants mirror credit.FetishForm.
const FETISH_FORMS: { value: number; label: string }[] = [
  { value: 1, label: "Simple Interest" },
  { value: 2, label: "Compound Interest" },
  { value: 3, label: "M—M' (pure fetish)" },
];

function bpToPercent(bp: number): string {
  return `${(bp / 100).toFixed(2)}%`;
}

function fetishLabel(form: number): string {
  return FETISH_FORMS.find((f) => f.value === form)?.label ?? String(form);
}

// Seed exemplars from Capital Vol. III Ch. 24.
const SEED_EXEMPLARS = [
  {
    label: "Dr. Price — 1 penny at 5% compound for 100 years",
    principal: 1,
    rate_bp: 500,
    period_years: 100,
  },
  {
    label: "Pitt's Sinking Fund — £250,000 at 5% compound for 10 periods",
    principal: 250000,
    rate_bp: 500,
    period_years: 10,
  },
];

// computePreview mirrors the Go integer compound iteration using Math.floor.
// v_{k+1} = floor(v_k * (10000 + rateBP) / 10000) (simplified, no round-half-up)
function computePreview(principal: number, rateBP: number, periodYears: number): number | null {
  if (principal <= 0 || rateBP <= 0 || periodYears <= 0) return null;
  let v = principal;
  for (let i = 0; i < periodYears; i++) {
    v = Math.floor((v * (10000 + rateBP)) / 10000);
  }
  return v;
}

export function Ch24CompoundInterest() {
  const [items, setItems] = useState<CompoundInterestSchedule[]>([]);
  const [listError, setListError] = useState<string | null>(null);

  // Create-form state.
  const [principal, setPrincipal] = useState(1);
  const [rateBP, setRateBP] = useState(500);
  const [periodYears, setPeriodYears] = useState(100);
  const [fetishForm, setFetishForm] = useState(2);
  const [createError, setCreateError] = useState<string | null>(null);

  // Live-preview using Math.floor (approximation of integer iteration).
  const previewFinal = computePreview(principal, rateBP, periodYears);

  async function refreshItems() {
    try {
      const data = await financeApi.listCompoundInterestSchedules();
      setItems(data);
    } catch (e) {
      setListError(String(e));
    }
  }

  useEffect(() => {
    refreshItems().catch(() => {});
  }, []);

  async function handleCreate() {
    setCreateError(null);
    try {
      await financeApi.createCompoundInterestSchedule({
        principal,
        rate_bp: rateBP,
        period_years: periodYears,
        fetish_form: fetishForm,
      });
      await refreshItems();
    } catch (e) {
      setCreateError(String(e));
    }
  }

  return (
    <div className="v3-ch24">
      <section className="v3-ch24-intro">
        <p>
          Interest-bearing capital in the form <strong>M—M'</strong> is the most fetishised
          form of capital: money appears to breed money automatically, erasing all trace of
          production and surplus-labour. The compound-interest formula{" "}
          <code>s = c(1+i)^n</code> captures this fetish mathematically.
        </p>
        <p>
          Dr. Price's fantasy — 1 penny at 5% compound over 1776 years exceeding the gold
          value of 150 million earths — illustrates the absurdity of mistaking arithmetic
          for economic law. The qualitative limit is the working day: surplus-value is
          bounded by living labour, not by a mathematical formula.
        </p>
        <p>
          Invariant: <code>new_value_created == 0</code> — the fetish form conceals but
          does not eliminate the zero-value-creation character of money as such.
        </p>
      </section>

      {/* Canonical Marx exemplars */}
      <section className="v3-ch24-section" aria-label="Compound interest exemplars">
        <h4>Marx's exemplars (Ch. 24)</h4>
        <p className="v3-ch24-subintro">
          Dr. Price and Pitt's sinking fund — compound interest as fetish illustration.
        </p>
        <div className="v3-ch24-cards">
          {SEED_EXEMPLARS.map((ex) => (
            <div className="v3-ch24-card" key={ex.label}>
              <div className="v3-ch24-card-label">{ex.label}</div>
              <div className="v3-ch24-meta">
                <span>Principal: {ex.principal}</span>
                <span>Rate: {bpToPercent(ex.rate_bp)}</span>
                <span>Periods: {ex.period_years}</span>
              </div>
            </div>
          ))}
        </div>
      </section>

      {/* Calculator */}
      <section className="v3-ch24-section" aria-label="Compound interest calculator">
        <h4>Compute a schedule</h4>
        <p className="v3-ch24-subintro">
          FinalValue is computed by integer iteration on the server (no float64).
          new_value_created is always 0.
        </p>
        <div className="v3-ch24-form">
          <label>
            Principal (smallest currency unit, &gt; 0)
            <input
              type="number"
              min={1}
              value={principal}
              onChange={(e) => setPrincipal(Number(e.target.value))}
            />
          </label>
          <label>
            Rate (basis points, &gt; 0)
            <input
              type="number"
              min={1}
              value={rateBP}
              onChange={(e) => setRateBP(Number(e.target.value))}
            />
            <span className="v3-ch24-hint">
              100 bp = 1%. Dr. Price's 5% = 500 bp.
            </span>
          </label>
          <label>
            Periods (years, &gt; 0)
            <input
              type="number"
              min={1}
              value={periodYears}
              onChange={(e) => setPeriodYears(Number(e.target.value))}
            />
          </label>
          <label>
            Fetish form
            <select value={fetishForm} onChange={(e) => setFetishForm(Number(e.target.value))}>
              {FETISH_FORMS.map((f) => (
                <option key={f.value} value={f.value}>
                  {f.label}
                </option>
              ))}
            </select>
          </label>
        </div>

        {previewFinal !== null && (
          <div className="v3-ch24-preview" aria-label="Preview">
            <div className="v3-ch24-preview-row">
              <span className="v3-ch24-preview-label">Principal</span>
              <span className="v3-ch24-preview-value">{principal}</span>
            </div>
            <div className="v3-ch24-preview-row">
              <span className="v3-ch24-preview-label">Rate</span>
              <span className="v3-ch24-preview-value">{bpToPercent(rateBP)}</span>
            </div>
            <div className="v3-ch24-preview-row">
              <span className="v3-ch24-preview-label">Periods</span>
              <span className="v3-ch24-preview-value">{periodYears}</span>
            </div>
            <div className="v3-ch24-preview-row">
              <span className="v3-ch24-preview-label">Final Value (preview)</span>
              <span className="v3-ch24-preview-value v3-ch24-preview-final">{previewFinal}</span>
            </div>
            <div className="v3-ch24-preview-row">
              <span className="v3-ch24-preview-label">New Value Created</span>
              <span className="v3-ch24-preview-value">0 (fetish form)</span>
            </div>
          </div>
        )}

        <button
          className="v3-ch24-btn"
          onClick={handleCreate}
          disabled={principal <= 0 || rateBP <= 0 || periodYears <= 0}
        >
          Compute &amp; save schedule
        </button>
        {createError && <p className="v3-ch24-error">{createError}</p>}
      </section>

      {/* Stored schedules */}
      <section className="v3-ch24-section" aria-label="Stored schedules">
        <h4>Stored compound-interest schedules</h4>
        {listError && <p className="v3-ch24-error">{listError}</p>}
        <div className="v3-ch24-cards">
          {items.length === 0 && (
            <p className="v3-ch24-subintro">No schedules stored yet.</p>
          )}
          {items.map((s) => (
            <div className="v3-ch24-card v3-ch24-card--stored" key={s.id}>
              <div className="v3-ch24-card-label">
                Principal {s.principal} @ {bpToPercent(s.rate_bp)} × {s.period_years} periods
              </div>
              <div className="v3-ch24-final">Final Value: {s.final_value}</div>
              <div className="v3-ch24-meta">
                <span className="v3-ch24-fetish-tag">{fetishLabel(s.fetish_form)}</span>
                <span>new_value_created: {s.new_value_created}</span>
              </div>
            </div>
          ))}
        </div>
      </section>
    </div>
  );
}
