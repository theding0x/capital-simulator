import { useEffect, useState } from "react";
import type {
  NaturalProcessIndustry,
  ProductionLabourGap,
  WorkingPeriod,
} from "../../types";
import { api } from "../../api";
import "./Ch13TimeOfProduction.css";

// Vol. II Ch. 13 — The Time of Production.
//
// Key result: production-time = labour-time + gap.
// Only labour-time creates surplus-value. The gap is latent production-time
// during which capital is locked in a natural process — fermentation, growth,
// chemical transformation, seasonal rest — and no new value is being added.
//
// Marx's six fixtures (seed migration 00047_v2_ch13_seed.sql):
//   Corn farm 1871      — 280-day production, 10 labour-days (3.6%)
//   Winery 1871         — 365-day production, 10 labour-days (2.7%)
//   Cattle farm 1871    — 730-day production,  5 labour-days (0.7%)
//   Timber 1871         — 18 250-day production, 30 labour-days (0.16%)
//   Porcelain kiln 1871 — 30-day production,   3 labour-days (10%)
//   Fallow land 1871    — 180-day production,   5 labour-days (2.8%)

// Human-readable nanosecond formatter.
function fmtNanos(ns: number): string {
  const days = ns / 86_400_000_000_000;
  if (days >= 365 * 10) return `${(days / 365).toFixed(0)} yr`;
  if (days >= 365) return `${(days / 365).toFixed(1)} yr`;
  if (days >= 1) return `${Math.round(days)} d`;
  const hours = ns / 3_600_000_000_000;
  return `${hours.toFixed(1)} h`;
}

// Names for the six seed working periods (keyed by industrial_capital_id).
const PERIOD_NAMES: Record<string, string> = {
  "corn-farm-1871-ic":   "Corn farm (annual crop)",
  "winery-1871-ic":      "Winery (fermentation)",
  "cattle-farm-1871-ic": "Cattle farm (maturation)",
  "timber-1871-ic":      "Timber (50-year rotation)",
  "porcelain-1871-ic":   "Porcelain kiln (firing)",
  "fallow-farm-1871-ic": "Fallow land (seasonal rest)",
};

// ---- GapBar ---------------------------------------------------------------

function GapBar({ gap }: { gap: ProductionLabourGap }) {
  const labourFraction = gap.labour_time_nanos / gap.production_time_nanos;
  const labourPct = (labourFraction * 100).toFixed(1);
  const gapPct = ((1 - labourFraction) * 100).toFixed(1);

  return (
    <div className="ch13-gap-bar-container">
      <div className="ch13-gap-bar">
        <div
          className="ch13-gap-bar-labour"
          style={{ width: `${labourFraction * 100}%` }}
          title={`Labour: ${fmtNanos(gap.labour_time_nanos)}`}
        />
        <div
          className="ch13-gap-bar-gap"
          style={{ width: `${(1 - labourFraction) * 100}%` }}
          title={`Gap: ${fmtNanos(gap.gap_nanos)}`}
        />
      </div>
      <div className="ch13-gap-bar-legend">
        <span className="ch13-labour-label">
          Labour&nbsp;{labourPct}%&nbsp;({fmtNanos(gap.labour_time_nanos)})
        </span>
        <span className="ch13-gap-label">
          Gap&nbsp;{gapPct}%&nbsp;({fmtNanos(gap.gap_nanos)})
        </span>
      </div>
    </div>
  );
}

// ---- BenchmarkRow ---------------------------------------------------------

function BenchmarkRow({ b }: { b: NaturalProcessIndustry }) {
  const labourFraction = b.ratio_basis_points / 10_000;
  const labourPct = (labourFraction * 100).toFixed(2);

  return (
    <div className="ch13-benchmark-row">
      <span className="ch13-bench-name">{b.industry_name}</span>
      <span className="ch13-bench-days">
        {b.typical_production_days_count.toLocaleString()}&nbsp;d
      </span>
      <span className="ch13-bench-labour">
        {b.typical_labour_days_count}&nbsp;d
      </span>
      <div className="ch13-bench-bar">
        <div
          className="ch13-bench-bar-fill"
          style={{ width: `${labourFraction * 100}%` }}
        />
      </div>
      <span className="ch13-bench-ratio">{labourPct}%</span>
    </div>
  );
}

// ---- Panel ----------------------------------------------------------------

export function Ch13TimeOfProduction() {
  const [periods, setPeriods] = useState<WorkingPeriod[] | null>(null);
  const [gaps, setGaps] = useState<Map<string, ProductionLabourGap>>(
    new Map()
  );
  const [benchmarks, setBenchmarks] = useState<NaturalProcessIndustry[] | null>(
    null
  );
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let cancelled = false;

    async function load() {
      try {
        const [wps, benchs] = await Promise.all([
          api.listWorkingPeriods({ mode: "connected" }),
          api.listIndustryBenchmarks(),
        ]);

        const ch13Periods = wps.filter((p) =>
          Object.prototype.hasOwnProperty.call(PERIOD_NAMES, p.industrial_capital_id)
        );

        const gapResults = await Promise.all(
          ch13Periods.map((p) =>
            api
              .getProductionLabourGap(p.id)
              .then((g): [string, ProductionLabourGap] => [p.id, g])
              .catch(() => null)
          )
        );

        if (cancelled) return;

        const m = new Map<string, ProductionLabourGap>();
        for (const r of gapResults) {
          if (r) m.set(r[0], r[1]);
        }

        setPeriods(ch13Periods);
        setBenchmarks(benchs);
        setGaps(m);
      } catch (e) {
        if (!cancelled)
          setError(e instanceof Error ? e.message : "failed to load");
      }
    }

    load();
    return () => {
      cancelled = true;
    };
  }, []);

  if (error) return <p className="ch13-error">Error: {error}</p>;
  if (!periods || !benchmarks) return <p className="ch13-loading">Loading…</p>;

  return (
    <div className="ch13-panel">
      {/* Key-result insight */}
      <div className="ch13-insight-box">
        <p className="ch13-insight-formula">
          Production-time&nbsp;=&nbsp;Labour-time&nbsp;+&nbsp;
          <span className="ch13-gap-accent">Gap</span>
        </p>
        <p className="ch13-insight-note">
          Only labour-time creates surplus-value. The gap is latent
          production-time during which capital is locked in a natural process —
          fermentation, crop growth, animal maturation, chemical change, seasonal
          rest — and no value is being added. Capital immobilised in the gap
          earns nothing; the longer the gap, the lower the annual rate of
          surplus-value for a given capital advance.
        </p>
      </div>

      {/* Production-labour gap bars for each Ch13 working period */}
      {periods.length > 0 && (
        <section className="ch13-periods-section">
          <h3>Production-Labour Gaps</h3>
          <div className="ch13-periods-grid">
            {periods.map((p) => {
              const gap = gaps.get(p.id);
              const name =
                PERIOD_NAMES[p.industrial_capital_id] ??
                p.industrial_capital_id;

              return (
                <div key={p.id} className="ch13-period-card">
                  <div className="ch13-period-name">{name}</div>
                  <div className="ch13-period-stat">
                    <span className="ch13-stat-label">Production period</span>
                    <span className="ch13-stat-value">
                      {p.working_day_count}&nbsp;days
                    </span>
                  </div>
                  {gap ? (
                    <GapBar gap={gap} />
                  ) : (
                    <p className="ch13-empty-gap">No gap recorded.</p>
                  )}
                </div>
              );
            })}
          </div>
        </section>
      )}

      {periods.length === 0 && (
        <p className="ch13-empty">
          No connected working periods found. The seed migration populates six
          natural-process fixtures (corn, wine, cattle, timber, porcelain,
          fallow).
        </p>
      )}

      {/* Industry benchmarks table */}
      <section className="ch13-benchmarks-section">
        <h3>Industry Benchmarks — Labour/Production Ratio</h3>
        <p className="ch13-bench-subtitle">
          ratio&nbsp;=&nbsp;labour-days&nbsp;&divide;&nbsp;production-days
          &nbsp;&times;&nbsp;100&nbsp;(basis-points&nbsp;/&nbsp;100)
        </p>
        <div className="ch13-benchmark-header">
          <span>Industry</span>
          <span>Production</span>
          <span>Labour</span>
          <span>Labour fraction</span>
          <span>Ratio</span>
        </div>
        {benchmarks.length === 0 ? (
          <p className="ch13-empty">No benchmarks recorded.</p>
        ) : (
          <div className="ch13-benchmark-list">
            {benchmarks.map((b) => (
              <BenchmarkRow key={b.id} b={b} />
            ))}
          </div>
        )}
      </section>
    </div>
  );
}
