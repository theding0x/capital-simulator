import { useEffect, useState } from "react";
import type {
  AggregateTurnover,
  ComponentTurnoverContribution,
  CrisisCyclePhaseRecord,
} from "../../types";
import "./Ch09AggregateTurnover.css";

const BASE = "/api/v1";

async function fetchAggregates(): Promise<AggregateTurnover[]> {
  const res = await fetch(`${BASE}/aggregate-turnovers`);
  if (!res.ok) throw new Error(`${res.status}`);
  return res.json();
}

function pence(n: number): string {
  return `£${(n / 100).toLocaleString(undefined, { minimumFractionDigits: 0 })}`;
}

function bpToMultiple(bp: number): string {
  return (bp / 10_000).toFixed(4).replace(/\.?0+$/, "");
}

function minutesToYears(m: number): string {
  return (m / 525_600).toFixed(1);
}

const PHASE_LABEL: Record<string, string> = {
  depression: "Depression",
  medium_activity: "Medium Activity",
  precipitancy: "Precipitancy",
  crisis: "Crisis",
};

function ContributionRow({ c }: { c: ComponentTurnoverContribution }) {
  const multiple = c.turnover_number_basis_points / 10_000;
  return (
    <tr className="ch09-contrib-row">
      <td>
        <span className={`ch09-kind ch09-kind--${c.kind}`}>{c.kind}</span>
      </td>
      <td>{pence(c.advanced_pence)}</td>
      <td>{multiple.toFixed(2)}×/yr</td>
      <td>{pence(c.annual_return_pence)}</td>
      <td>
        <span className={`ch09-diff ch09-diff--${c.difference_kind}`}>
          {c.difference_kind}
        </span>
      </td>
    </tr>
  );
}

function CrisisPhaseTag({ p }: { p: CrisisCyclePhaseRecord }) {
  return (
    <span className={`ch09-phase ch09-phase--${p.phase}`}>
      {PHASE_LABEL[p.phase] ?? p.phase}
    </span>
  );
}

function AggregateCard({ at }: { at: AggregateTurnover }) {
  const meanTermYears =
    at.annual_turned_over_pence > 0
      ? (at.advanced_total_pence / at.annual_turned_over_pence).toFixed(2)
      : "—";

  return (
    <div className="ch09-card">
      <div className="ch09-card__header">
        <span className="ch09-card__id" title={at.id}>
          {at.industrial_capital_id || at.id.slice(0, 12)}
        </span>
        <span className="ch09-card__n">
          n&nbsp;=&nbsp;{bpToMultiple(at.aggregate_number_basis_points)}
          <span className="ch09-card__bp">
            &nbsp;({at.aggregate_number_basis_points}&nbsp;bp)
          </span>
        </span>
      </div>

      <div className="ch09-card__totals">
        <div className="ch09-stat">
          <span className="ch09-stat__label">Advanced</span>
          <span className="ch09-stat__value">{pence(at.advanced_total_pence)}</span>
        </div>
        <div className="ch09-stat">
          <span className="ch09-stat__label">Annual Turnover</span>
          <span className="ch09-stat__value">{pence(at.annual_turned_over_pence)}</span>
        </div>
        <div className="ch09-stat">
          <span className="ch09-stat__label">Mean Term</span>
          <span className="ch09-stat__value">{meanTermYears}&nbsp;yr</span>
        </div>
      </div>

      {at.contributions.length > 0 && (
        <table className="ch09-table">
          <thead>
            <tr>
              <th>Kind</th>
              <th>Advanced</th>
              <th>Turns/yr</th>
              <th>Annual Return</th>
              <th>Difference</th>
            </tr>
          </thead>
          <tbody>
            {at.contributions.map((c) => (
              <ContributionRow key={c.id} c={c} />
            ))}
          </tbody>
        </table>
      )}

      {at.lifetime_cycle && (
        <div className="ch09-cycle">
          <div className="ch09-cycle__header">
            Lifetime Cycle —&nbsp;
            {minutesToYears(at.lifetime_cycle.duration_minutes)}&nbsp;yr
            {at.lifetime_cycle.completed_cycles > 0 && (
              <span>&nbsp;·&nbsp;{at.lifetime_cycle.completed_cycles} completed</span>
            )}
          </div>
          {at.latest_crisis_phase && (
            <div className="ch09-cycle__latest">
              Latest phase:&nbsp;
              <CrisisPhaseTag p={at.latest_crisis_phase} />
              {at.latest_crisis_phase.notes && (
                <span className="ch09-cycle__notes">
                  &nbsp;—&nbsp;{at.latest_crisis_phase.notes}
                </span>
              )}
            </div>
          )}
          {at.lifetime_cycle.crisis_phases &&
            at.lifetime_cycle.crisis_phases.length > 1 && (
              <div className="ch09-cycle__history">
                {at.lifetime_cycle.crisis_phases.map((p) => (
                  <CrisisPhaseTag key={p.id} p={p} />
                ))}
              </div>
            )}
        </div>
      )}
    </div>
  );
}

export function Ch09AggregateTurnover() {
  const [aggregates, setAggregates] = useState<AggregateTurnover[]>([]);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetchAggregates()
      .then(setAggregates)
      .catch((e: unknown) =>
        setError(e instanceof Error ? e.message : String(e))
      );
  }, []);

  return (
    <div className="ch09">
      <p className="ch09-epigraph">
        "The aggregate turnover of the advanced productive capital is the mean
        average of the turnover of its different component parts." — Marx, Vol. II Ch. 9
      </p>

      {error && <p className="ch09-error">Could not load data: {error}</p>}

      {aggregates.length === 0 && !error && (
        <p className="ch09-empty">No aggregate turnovers recorded yet.</p>
      )}

      <div className="ch09-grid">
        {aggregates.map((at) => (
          <AggregateCard key={at.id} at={at} />
        ))}
      </div>
    </div>
  );
}
