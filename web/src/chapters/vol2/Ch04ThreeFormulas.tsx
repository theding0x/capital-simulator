import { useEffect, useState } from "react";
import type { IndustrialCapital } from "../../types";
import "./Ch04ThreeFormulas.css";

const BASE = "/api/v1/industrial-capitals";

async function fetchCapitals(): Promise<IndustrialCapital[]> {
  const res = await fetch(BASE);
  if (!res.ok) throw new Error(`${res.status}`);
  return res.json();
}

function pence(n: number): string {
  return `£${(n / 100).toFixed(2)}`;
}

function StageBar({ ic }: { ic: IndustrialCapital }) {
  const dist = ic.latest_distribution;
  if (!dist) return <div className="ch04-stage-bar ch04-stage-bar--empty">No distribution yet</div>;
  const total = dist.money_pence + dist.production_pence + dist.commodity_pence;
  if (total === 0) return null;
  const pct = (v: number) => `${((v / total) * 100).toFixed(1)}%`;
  return (
    <div className="ch04-stage-bar" title="M / P / C distribution">
      <div className="ch04-stage-bar__segment ch04-stage-bar__segment--money" style={{ width: pct(dist.money_pence) }}>
        <span>M</span>
      </div>
      <div className="ch04-stage-bar__segment ch04-stage-bar__segment--prod" style={{ width: pct(dist.production_pence) }}>
        <span>P</span>
      </div>
      <div className="ch04-stage-bar__segment ch04-stage-bar__segment--commodity" style={{ width: pct(dist.commodity_pence) }}>
        <span>C′</span>
      </div>
    </div>
  );
}

export function Ch04ThreeFormulas() {
  const [capitals, setCapitals] = useState<IndustrialCapital[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchCapitals()
      .then(setCapitals)
      .catch((e) => setError(String(e)))
      .finally(() => setLoading(false));
  }, []);

  if (loading) return <div className="ch04-panel ch04-loading">Loading industrial capitals…</div>;
  if (error) return <div className="ch04-panel ch04-error">Error: {error}</div>;

  return (
    <div className="ch04-panel">
      <div className="ch04-header">
        <h2>Ch. 4 — The Three Formulas of the Circuit</h2>
        <p className="ch04-subtitle">
          Every industrial capital occupies M, P, and C simultaneously —{" "}
          <em>Nebeneinander</em> (co-existence).
        </p>
      </div>

      <div className="ch04-formulas">
        <div className="ch04-formula">
          <span className="ch04-formula__label">I</span>
          <code>M—C(Lp+Mp)…P…C′—M′</code>
        </div>
        <div className="ch04-formula">
          <span className="ch04-formula__label">II</span>
          <code>P…C′—M′—C(Lp+Mp)…P</code>
        </div>
        <div className="ch04-formula">
          <span className="ch04-formula__label">III</span>
          <code>C′—M′—C(Lp+Mp)…P…C′</code>
        </div>
      </div>

      <p className="v2-ch04-nebeneinander">
        Marx insists the three formulas are not stages a capital passes
        through one after another — they are simultaneous moments of a
        single capital in motion. <em>Nebeneinander</em>, not Nacheinander.
      </p>

      {capitals.length === 0 ? (
        <p className="ch04-empty">No industrial capitals recorded yet.</p>
      ) : (
        <div className="ch04-list">
          {capitals.map((ic) => (
            <div key={ic.id} className={`ch04-card ch04-card--${ic.status}`}>
              <div className="ch04-card__header">
                <span className="ch04-card__agent">{ic.agent_id || ic.id.slice(0, 8)}</span>
                <span className={`ch04-card__status ch04-card__status--${ic.status}`}>
                  {ic.status}
                </span>
                <span className="ch04-card__mode">{ic.economy_mode}</span>
              </div>

              <div className="ch04-card__total">Total: {pence(ic.total_pence)}</div>

              <StageBar ic={ic} />

              {ic.latest_distribution && (
                <div className="ch04-card__dist">
                  <span>M: {pence(ic.latest_distribution.money_pence)}</span>
                  <span>P: {pence(ic.latest_distribution.production_pence)}</span>
                  <span>C′: {pence(ic.latest_distribution.commodity_pence)}</span>
                </div>
              )}

              {ic.open_blocks.length > 0 && (
                <div className="ch04-card__blocks">
                  {ic.open_blocks.map((b) => (
                    <span key={b.id} className="ch04-block-badge">
                      ⚠ {b.stage}: {b.reason.replace(/_/g, " ")}
                    </span>
                  ))}
                </div>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
