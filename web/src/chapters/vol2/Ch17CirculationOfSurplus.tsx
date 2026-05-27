import { useEffect, useState } from "react";
import type { SurplusCirculation, RealisationPuzzle } from "../../types";
import { surplusCirculationApi } from "../../api";
import "./Ch17CirculationOfSurplus.css";

// Vol. II Ch. 17 — The Circulation of Surplus-Value.
//
// The central puzzle: the capitalist class collectively spends money to
// realise surplus-value as M'. But where does that money come from?
// Marx identifies five sources — none of which create new value, each of
// which merely recirculates existing money or defers realisation to a
// future period.
//
// Seed fixtures (00053_v2_ch17_seed.sql):
//   Marx 1871 — 5eed0000000000001701: total surplus £5,000,000 p.a.
//     - revenue_expenditure: £4,000,000 (80% consumed)
//     - accumulated_capital:  £950,000  (19% re-capitalised)
//     - gold_replacement:      £50,000  ( 1% gold reproduction)

// ── helpers ──────────────────────────────────────────────────────────────────

function penceToPounds(p: number): string {
  const pounds = p / 240;
  if (pounds >= 1_000_000) return "£" + (pounds / 1_000_000).toFixed(2) + "m";
  if (pounds >= 1_000) return "£" + (pounds / 1_000).toFixed(1) + "k";
  return "£" + pounds.toFixed(0);
}

function basisPointsPct(bp: number): string {
  return (bp / 100).toFixed(2) + "%";
}

const SOURCES: { key: string; label: string }[] = [
  { key: "revenue_expenditure",  label: "Revenue expenditure (capitalists consume)" },
  { key: "accumulated_capital",  label: "Accumulated capital (re-capitalised)" },
  { key: "gold_replacement",     label: "Gold reproduction (money-commodity)" },
  { key: "credit_deferred",      label: "Credit / deferred realisation" },
  { key: "foreign_trade",        label: "Foreign trade (external market)" },
];

// ── Widget 1 — Realisation puzzle ────────────────────────────────────────────

function PuzzleBox({ puzzles }: { puzzles: RealisationPuzzle[] }) {
  return (
    <div className="ch17-puzzle-box">
      <p className="ch17-section-title">The Realisation Puzzle</p>
      <blockquote className="ch17-quote">
        "Where does the money come from to realise the aggregate surplus-value?
        The capitalists cannot conjure it out of nothing."
        <cite>— Capital, Vol. II, Ch. 17</cite>
      </blockquote>
      {puzzles.length > 0 && (
        <ul className="ch17-puzzle-list">
          {puzzles.map((p) => (
            <li key={p.surplus_circulation_id} className="ch17-puzzle-item">
              <span className="ch17-puzzle-statement">{p.puzzle_statement}</span>
              {p.resolved_in_chapter > 0 && (
                <span className="ch17-puzzle-resolution">
                  → Resolved in Ch. {p.resolved_in_chapter}
                </span>
              )}
              {p.resolved_in_chapter === 0 && (
                <span className="ch17-puzzle-unresolved">→ Formally deferred</span>
              )}
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}

// ── Widget 2 — Five sources flow ─────────────────────────────────────────────

const FLOW_BARS: { label: string; pct: number; css: string }[] = [
  { label: "Revenue expenditure",     pct: 80, css: "bar-revenue" },
  { label: "Re-capitalised (acc.)",   pct: 19, css: "bar-capital" },
  { label: "Gold reproduction",       pct:  1, css: "bar-gold"    },
];

function FlowDiagram() {
  return (
    <div className="ch17-flow-box">
      <p className="ch17-section-title">
        Surplus-Value Flow — Marx's 1871 Illustration
      </p>
      <p className="ch17-flow-subtitle">
        £5 m surplus-value decomposed by realisation source
      </p>
      <div className="ch17-bars">
        {FLOW_BARS.map((b) => (
          <div key={b.label} className="ch17-bar-row">
            <span className="ch17-bar-label">{b.label}</span>
            <div className="ch17-bar-track">
              <div
                className={`ch17-bar-fill ${b.css}`}
                style={{ width: `${b.pct}%` }}
              />
            </div>
            <span className="ch17-bar-pct">{b.pct}%</span>
          </div>
        ))}
      </div>
      <p className="ch17-flow-note">
        Credit and foreign trade defer rather than originate the money-fund.
        Chapters 20–21 resolve the aggregate reproduction puzzle formally.
      </p>
    </div>
  );
}

// ── Widget 3 — Live surplus circulations ────────────────────────────────────

interface SCRowProps {
  sc: SurplusCirculation;
}

function SCRow({ sc }: SCRowProps) {
  const realisedPct =
    sc.total_surplus_pence > 0
      ? Math.round((sc.realised_surplus_pence / sc.total_surplus_pence) * 100)
      : 0;

  return (
    <div className="ch17-sc-row">
      <div className="ch17-sc-header">
        <span className="ch17-sc-period">{sc.period}</span>
        <span className="ch17-sc-id">…{sc.id.slice(-6)}</span>
      </div>
      <div className="ch17-sc-amounts">
        <div className="ch17-sc-amount">
          <span className="ch17-sc-amount-label">Total surplus</span>
          <span>{penceToPounds(sc.total_surplus_pence)}</span>
        </div>
        <div className="ch17-sc-amount">
          <span className="ch17-sc-amount-label">Realised</span>
          <span className="ch17-amount-realised">
            {penceToPounds(sc.realised_surplus_pence)}
          </span>
        </div>
        <div className="ch17-sc-amount">
          <span className="ch17-sc-amount-label">Unrealised</span>
          <span className="ch17-amount-unrealised">
            {penceToPounds(sc.unrealised_surplus_pence)}
          </span>
        </div>
      </div>
      <div className="ch17-sc-progress-track">
        <div
          className="ch17-sc-progress-fill"
          style={{ width: `${realisedPct}%` }}
        />
      </div>
      <div className="ch17-sc-pct-label">{realisedPct}% realised</div>
      {sc.realisation_sources.length > 0 && (
        <ul className="ch17-source-list">
          {sc.realisation_sources.map((e, i) => {
            const srcLabel =
              SOURCES.find((s) => s.key === e.source)?.label ?? e.source;
            return (
              <li key={i} className="ch17-source-item">
                <span className="ch17-source-label">{srcLabel}</span>
                <span>{penceToPounds(e.pence)}</span>
                {e.deferred_realisation_pence > 0 && (
                  <span className="ch17-source-deferred">
                    +{penceToPounds(e.deferred_realisation_pence)} deferred
                  </span>
                )}
              </li>
            );
          })}
        </ul>
      )}
    </div>
  );
}

// ── Widget 4 — Social aggregate ──────────────────────────────────────────────

const SEED_PERIOD = "1871";

function SocialAggregateWidget() {
  const [data, setData] = useState<{
    total_advanced_pence: number;
    total_annual_output_pence: number;
    total_annual_surplus_pence: number;
    department_i_share_bp: number;
    department_ii_share_bp: number;
  } | null>(null);
  const [err, setErr] = useState<string | null>(null);

  useEffect(() => {
    surplusCirculationApi
      .socialAggregate(SEED_PERIOD)
      .then(setData)
      .catch(() => setErr("aggregate unavailable"));
  }, []);

  if (err) return <p className="ch17-error">{err}</p>;
  if (!data) return <p className="ch17-loading">Loading aggregate…</p>;

  return (
    <div className="ch17-aggregate-box">
      <p className="ch17-section-title">
        Social Capital Aggregate — {SEED_PERIOD}
      </p>
      <div className="ch17-aggregate-grid">
        <div className="ch17-agg-cell">
          <span className="ch17-agg-label">Total capital advanced</span>
          <span className="ch17-agg-value">
            {penceToPounds(data.total_advanced_pence)}
          </span>
        </div>
        <div className="ch17-agg-cell">
          <span className="ch17-agg-label">Annual output</span>
          <span className="ch17-agg-value">
            {penceToPounds(data.total_annual_output_pence)}
          </span>
        </div>
        <div className="ch17-agg-cell">
          <span className="ch17-agg-label">Annual surplus-value</span>
          <span className="ch17-agg-value ch17-agg-surplus">
            {penceToPounds(data.total_annual_surplus_pence)}
          </span>
        </div>
        <div className="ch17-agg-cell">
          <span className="ch17-agg-label">Dept I share</span>
          <span className="ch17-agg-value">
            {basisPointsPct(data.department_i_share_bp)}
          </span>
        </div>
        <div className="ch17-agg-cell">
          <span className="ch17-agg-label">Dept II share</span>
          <span className="ch17-agg-value">
            {basisPointsPct(data.department_ii_share_bp)}
          </span>
        </div>
      </div>
      <p className="ch17-agg-note">
        Dept I / II breakdown is the seed placeholder from the 1871 illustration;
        the full two-department model is formalised in Chs. 20–21.
      </p>
    </div>
  );
}

// ── Root ─────────────────────────────────────────────────────────────────────

export function Ch17CirculationOfSurplus() {
  const [circulations, setCirculations] = useState<SurplusCirculation[]>([]);
  const [puzzles, setPuzzles] = useState<RealisationPuzzle[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    Promise.all([
      surplusCirculationApi.list(),
      surplusCirculationApi.realisationPuzzles(),
    ])
      .then(([listResp, puzzleResp]) => {
        setCirculations(listResp.items ?? []);
        setPuzzles(puzzleResp.items ?? []);
      })
      .catch(() => setError("Could not load surplus circulations"))
      .finally(() => setLoading(false));
  }, []);

  return (
    <div className="ch17-root">
      <PuzzleBox puzzles={puzzles} />
      <FlowDiagram />
      <SocialAggregateWidget />

      <div className="ch17-circulations-section">
        <p className="ch17-section-title">Surplus Circulations</p>
        {loading && <p className="ch17-loading">Loading…</p>}
        {error && <p className="ch17-error">{error}</p>}
        {!loading && !error && circulations.length === 0 && (
          <p className="ch17-empty">No surplus circulations recorded.</p>
        )}
        {circulations.map((sc) => (
          <SCRow key={sc.id} sc={sc} />
        ))}
      </div>
    </div>
  );
}
