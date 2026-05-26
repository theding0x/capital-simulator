import { useEffect, useState } from "react";
import type {
  WorkingPeriod,
  NaturalWorkingPeriodConstraint,
} from "../../types";
import { api } from "../../api";
import "./Ch12WorkingPeriod.css";

// Vol. II Ch. 12 — The Working Period.
//
// Key result: in a *connected* working period (e.g. a locomotive taking 84
// working-days), every day's circulating-capital outlay is tied up until the
// final day. The effective capital-advance multiplier = WorkingDayCount, so
// a capitalist must advance 84x a single day's outlay before any return.
// In a *discrete* process (e.g. yarn spun daily), the multiplier = 1.
//
// Natural constraints (annual crops, multi-year timber, livestock gestation)
// set inviolable floors — ErrPrematureDelivery blocks any period below them.
//
// Credit-financing (the London builder mortgage example) is the capitalist's
// response to the capital-locking problem; interest charges are resolved in
// Vol. III Part Five.

// Seed fixtures visible on fresh DB (migration 00045_v2_ch12_seed.sql):
//   SpinningMill1871 — discrete, 1 day, multiplier = 1
//   LocomotiveWorks   — connected, 84 days, multiplier = 84

function multiplier(wp: WorkingPeriod): number {
  return wp.mode === "connected" ? wp.working_day_count : 1;
}

function WorkingPeriodCard({ wp }: { wp: WorkingPeriod }) {
  const m = multiplier(wp);
  const isConnected = wp.mode === "connected";

  return (
    <div className={`ch12-period-card${isConnected ? " connected" : ""}`}>
      <h3>{wp.industrial_capital_id}</h3>
      <div className="ch12-stat-row">
        <span className="ch12-stat-label">Mode</span>
        <span className="ch12-stat-value">{wp.mode}</span>
      </div>
      <div className="ch12-stat-row">
        <span className="ch12-stat-label">Working days</span>
        <span className="ch12-stat-value">{wp.working_day_count}</span>
      </div>
      <div className="ch12-stat-row">
        <span className="ch12-stat-label">Day length</span>
        <span className="ch12-stat-value">{wp.working_day_length_minutes} min</span>
      </div>
      <div className="ch12-stat-row">
        <span className="ch12-stat-label">Capital-advance multiplier</span>
        <span className="ch12-stat-value accent">
          {m}
          {isConnected && <span className="ch12-multiplier-badge">connected</span>}
        </span>
      </div>
      {wp.interruptions && wp.interruptions.length > 0 && (
        <div className="ch12-stat-row">
          <span className="ch12-stat-label">Interruptions</span>
          <span className="ch12-stat-value">{wp.interruptions.length}</span>
        </div>
      )}
      {wp.shortenings && wp.shortenings.length > 0 && (
        <div className="ch12-stat-row">
          <span className="ch12-stat-label">Shortenings</span>
          <span className="ch12-stat-value">{wp.shortenings.length}</span>
        </div>
      )}
      {wp.credit_financing && (
        <div className="ch12-stat-row">
          <span className="ch12-stat-label">Borrowed capital</span>
          <span className="ch12-stat-value accent">
            {(wp.credit_financing.borrowed_capital_pence / 100).toFixed(2)} GBP
          </span>
        </div>
      )}
    </div>
  );
}

function ConstraintRow({ nc }: { nc: NaturalWorkingPeriodConstraint }) {
  return (
    <div className="ch12-constraint-row">
      <span>{nc.commodity_id}</span>
      <span className="ch12-constraint-days">
        &ge; {nc.min_working_days} days
      </span>
      <span className="ch12-constraint-reason">{nc.reason}</span>
    </div>
  );
}

export function Ch12WorkingPeriod() {
  const [periods, setPeriods] = useState<WorkingPeriod[] | null>(null);
  const [constraints, setConstraints] = useState<
    NaturalWorkingPeriodConstraint[] | null
  >(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    Promise.all([api.listWorkingPeriods(), api.listNaturalConstraints()])
      .then(([wps, ncs]) => {
        setPeriods(wps);
        setConstraints(ncs);
      })
      .catch((e: unknown) =>
        setError(e instanceof Error ? e.message : "failed to load")
      );
  }, []);

  if (error) return <p className="ch12-error">Error: {error}</p>;
  if (!periods || !constraints) return <p className="ch12-loading">Loading...</p>;

  const discrete = periods.filter((p) => p.mode === "discrete");
  const connected = periods.filter((p) => p.mode === "connected");

  return (
    <div className="ch12-panel">
      {/* SpinningMill vs LocomotiveWorks side-by-side */}
      {periods.length > 0 && (
        <section className="ch12-comparison">
          {discrete.length > 0 ? (
            <WorkingPeriodCard wp={discrete[0]} />
          ) : (
            <div className="ch12-period-card">
              <h3>Discrete (e.g. SpinningMill1871)</h3>
              <p className="ch12-empty">No discrete periods recorded.</p>
            </div>
          )}
          {connected.length > 0 ? (
            <WorkingPeriodCard wp={connected[0]} />
          ) : (
            <div className="ch12-period-card connected">
              <h3>Connected (e.g. LocomotiveWorks1867)</h3>
              <p className="ch12-empty">No connected periods recorded.</p>
            </div>
          )}
        </section>
      )}

      {periods.length === 0 && (
        <p className="ch12-empty">
          No working periods recorded. The seed migration populates
          SpinningMill1871 (discrete) and LocomotiveWorks1867 (connected).
        </p>
      )}

      {/* Natural-constraint reference table */}
      <section className="ch12-constraints-section">
        <h3>Natural Working-Period Constraints</h3>
        {constraints.length === 0 ? (
          <p className="ch12-empty">No constraints registered.</p>
        ) : (
          <div className="ch12-constraint-list">
            {constraints.map((nc) => (
              <ConstraintRow key={nc.id} nc={nc} />
            ))}
          </div>
        )}
      </section>

      {/* All shortenings across every period */}
      {periods.some((p) => p.shortenings && p.shortenings.length > 0) && (
        <section className="ch12-events-section">
          <h3>Working-Period Shortenings</h3>
          <div className="ch12-event-list">
            {periods.flatMap((p) =>
              (p.shortenings ?? []).map((s) => (
                <div key={s.id} className="ch12-event-item">
                  <span className="ch12-event-kind">{s.kind}</span>
                  <span>
                    {s.old_working_day_count} days &rarr;{" "}
                    {s.new_working_day_count} days
                    {s.additional_fixed_capital_pence > 0 &&
                      ` (£${(s.additional_fixed_capital_pence / 100).toFixed(2)} fixed capital)`}
                    {s.additional_labour_count > 0 &&
                      ` (+${s.additional_labour_count} workers)`}
                  </span>
                </div>
              ))
            )}
          </div>
        </section>
      )}

      {/* Vol. III forward-reference */}
      <section className="ch12-deferred-section">
        <h3>Vol. III &mdash; Interest on Borrowed Capital</h3>
        <div className="ch12-deferred-box">
          <p className="ch12-deferred-label">Deferred to Vol. III Part Five</p>
          <p>
            Credit-financing (the London builder mortgage) defers the
            capital-locking problem: the capitalist borrows against an unfinished
            building rather than tying up own reserves for 200+ days.
            MortgageProviderID is a forward-reference placeholder.
          </p>
          <p>
            Interest charges and the distinction between money-dealing capital
            and interest-bearing capital are resolved in Vol. III Part Five
            (finance-service).
          </p>
        </div>
      </section>
    </div>
  );
}
