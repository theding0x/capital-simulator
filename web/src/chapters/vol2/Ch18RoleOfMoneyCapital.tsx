import { useEffect, useState } from "react";
import { moneyCapitalApi } from "../../api";
import type { MoneySupplyApportionment, InterDepartmentSettlement } from "../../types";
import "./Ch18RoleOfMoneyCapital.css";

// Vol. II Ch. 18 — The Role of Money-Capital in Reproduction.
//
// Marx's key insight: money-capital does not create value, but it is
// indispensable for keeping the circuit in motion. The two departments
// must exchange with each other, and money is the medium.

// ── helpers ──────────────────────────────────────────────────────────────────

function penceToPounds(p: number): string {
  const pounds = p / 240;
  if (pounds >= 1_000_000) return "£" + (pounds / 1_000_000).toFixed(2) + "m";
  if (pounds >= 1_000) return "£" + (pounds / 1_000).toFixed(1) + "k";
  return "£" + pounds.toFixed(0);
}

function bpToTurnovers(bp: number): string {
  return (bp / 10000).toFixed(2) + "×/yr";
}

// ── Widget 1 — Apportionment Pie (percentage bars) ───────────────────────────

function ApportionmentPie() {
  const [appt, setAppt] = useState<MoneySupplyApportionment | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    moneyCapitalApi
      .listApportionments()
      .then((list) => setAppt(list.length > 0 ? list[0] : null))
      .catch(() => setAppt(null))
      .finally(() => setLoading(false));
  }, []);

  if (loading) return <div className="ch18-widget"><p className="ch18-loading">Loading…</p></div>;

  if (!appt) {
    return (
      <div className="ch18-widget">
        <p className="ch18-section-title">Money-Stock Apportionment</p>
        <p className="ch18-empty">No apportionment data. Seed the database to populate this panel.</p>
      </div>
    );
  }

  const total = appt.total_social_money_pence;
  const bars: { label: string; pence: number; css: string }[] = [
    { label: "Dept I (means of production)",       pence: appt.dept_i_reserve_pence,  css: "bar-dept-i"  },
    { label: "Dept II (articles of consumption)",  pence: appt.dept_ii_reserve_pence, css: "bar-dept-ii" },
    { label: "Wage-rotation fund",                 pence: appt.wage_rotation_pence,    css: "bar-wage"    },
    { label: "Idle hoard",                         pence: appt.idle_hoard_pence,            css: "bar-idle"    },
  ];

  return (
    <div className="ch18-widget">
      <p className="ch18-section-title">
        Money-Stock Apportionment — {appt.period}
      </p>
      <p className="ch18-widget-subtitle">
        Total M = {penceToPounds(total)} distributed across departments and funds
      </p>
      <div className="ch18-bars">
        {bars.map((b) => {
          const pct = Math.round((b.pence / total) * 100);
          return (
            <div key={b.label} className="ch18-bar-row">
              <span className="ch18-bar-label">{b.label}</span>
              <div className="ch18-bar-track">
                <div className={`ch18-bar-fill ${b.css}`} style={{ width: `${pct}%` }} />
              </div>
              <span className="ch18-bar-meta">{penceToPounds(b.pence)} ({pct}%)</span>
            </div>
          );
        })}
      </div>
      <p className="ch18-note">
        Money-capital does not produce value — it is the lubricant, not the
        engine. The idle hoard represents latent money-capital not yet
        re-entered into productive circulation.
      </p>
    </div>
  );
}

// ── Widget 2 — M×V Calculator ────────────────────────────────────────────────

function MVCalculator() {
  const [stockStr, setStockStr] = useState("1000");
  const [velocityStr, setVelocityStr] = useState("2");

  const stock = Math.max(0, Math.round(parseFloat(stockStr) || 0));
  const velocity = Math.max(0, parseFloat(velocityStr) || 0);
  const velocityBP = Math.round(velocity * 10000);
  const effectivePounds = stock * velocity;

  return (
    <div className="ch18-widget">
      <p className="ch18-section-title">M × V Calculator</p>
      <p className="ch18-widget-subtitle">
        The same stock of money can serve a larger total circulation if it
        turns over more frequently.
      </p>
      <div className="ch18-mv-form">
        <label className="ch18-mv-label">
          Money stock M (£)
          <input type="number" min="0" className="ch18-mv-input" value={stockStr}
            onChange={(e) => setStockStr(e.target.value)} />
        </label>
        <span className="ch18-mv-times">×</span>
        <label className="ch18-mv-label">
          Velocity V (turnovers/year)
          <input type="number" min="0" step="0.5" className="ch18-mv-input" value={velocityStr}
            onChange={(e) => setVelocityStr(e.target.value)} />
        </label>
        <span className="ch18-mv-equals">=</span>
        <div className="ch18-mv-result">
          <span className="ch18-mv-result-label">Effective circulating value</span>
          <span className="ch18-mv-result-value">£{effectivePounds.toFixed(0)}</span>
          <span className="ch18-mv-result-sub">
            ({bpToTurnovers(velocityBP)} at basis points {velocityBP})
          </span>
        </div>
      </div>
      <p className="ch18-note">
        Marx's correction to the Quantity Theory: the velocity of circulation
        (V) means a smaller M suffices for the same aggregate exchange.
      </p>
    </div>
  );
}

// ── Widget 3 — Wage Rotation Toy ─────────────────────────────────────────────

const FREQ_OPTIONS: { label: string; value: number }[] = [
  { label: "Weekly (52)", value: 52 },
  { label: "Monthly (12)", value: 12 },
  { label: "Annual (1)", value: 1 },
];

function WageRotationToy() {
  const [fundStr, setFundStr] = useState("1000");
  const [freq, setFreq] = useState(52);

  const fundPounds = Math.max(0, parseFloat(fundStr) || 0);
  const tiedUpPerCycle = freq > 0 ? fundPounds / freq : 0;

  return (
    <div className="ch18-widget">
      <p className="ch18-section-title">Wage Rotation Fund</p>
      <p className="ch18-widget-subtitle">
        Wages must be paid before output is sold. The capitalist must hold a
        money reserve proportional to the wage cycle.
      </p>
      <div className="ch18-wr-form">
        <label className="ch18-wr-label">
          Annual wage fund (£)
          <input type="number" min="0" className="ch18-wr-input" value={fundStr}
            onChange={(e) => setFundStr(e.target.value)} />
        </label>
        <label className="ch18-wr-label">
          Payment cycle
          <select className="ch18-wr-select" value={freq}
            onChange={(e) => setFreq(Number(e.target.value))}>
            {FREQ_OPTIONS.map((o) => (
              <option key={o.value} value={o.value}>{o.label}</option>
            ))}
          </select>
        </label>
      </div>
      <div className="ch18-wr-result">
        <span className="ch18-wr-result-label">Money tied up per cycle</span>
        <span className="ch18-wr-result-value">£{tiedUpPerCycle.toFixed(2)}</span>
        <span className="ch18-wr-result-sub">{freq} payment{freq !== 1 ? "s" : ""} per year</span>
      </div>
      <p className="ch18-note">
        The more frequently wages are paid, the smaller the advance required
        at any one time.
      </p>
    </div>
  );
}

// ── Widget 4 — Two-Department Flow Diagram ────────────────────────────────────

function TwoDeptFlowDiagram() {
  const [settlement, setSettlement] = useState<InterDepartmentSettlement | null>(null);

  useEffect(() => {
    moneyCapitalApi
      .listInterDepartmentSettlements()
      .then((list) => setSettlement(list.length > 0 ? list[0] : null))
      .catch(() => setSettlement(null));
  }, []);

  const flowLabel = settlement
    ? penceToPounds(settlement.settled_pence)
    : "£208";

  return (
    <div className="ch18-widget">
      <p className="ch18-section-title">Two-Department Exchange Flows</p>
      <p className="ch18-widget-subtitle">
        Simple reproduction requires that I(v+s) = II(c) — the exchange
        between departments must balance at the aggregate level.
      </p>
      <div className="ch18-flow-svg-wrapper">
        <svg className="ch18-flow-svg" viewBox="0 0 560 200" xmlns="http://www.w3.org/2000/svg"
          aria-label="Two-department flow diagram">
          <rect x="20" y="60" width="180" height="80" rx="8" className="ch18-svg-box-i" />
          <text x="110" y="93" textAnchor="middle" className="ch18-svg-label-main">Dept I</text>
          <text x="110" y="113" textAnchor="middle" className="ch18-svg-label-sub">Means of Production</text>
          <text x="110" y="130" textAnchor="middle" className="ch18-svg-label-sub">c + v + s</text>
          <rect x="360" y="60" width="180" height="80" rx="8" className="ch18-svg-box-ii" />
          <text x="450" y="93" textAnchor="middle" className="ch18-svg-label-main">Dept II</text>
          <text x="450" y="113" textAnchor="middle" className="ch18-svg-label-sub">Articles of Consumption</text>
          <text x="450" y="130" textAnchor="middle" className="ch18-svg-label-sub">c + v + s</text>
          <defs>
            <marker id="arrowhead" markerWidth="8" markerHeight="6" refX="8" refY="3" orient="auto">
              <polygon points="0 0, 8 3, 0 6" className="ch18-svg-arrow-marker" />
            </marker>
          </defs>
          <line x1="200" y1="85" x2="360" y2="85" className="ch18-svg-arrow-i" markerEnd="url(#arrowhead)" />
          <text x="280" y="78" textAnchor="middle" className="ch18-svg-flow-label">I(v+s) → IIc  {flowLabel}</text>
          <line x1="360" y1="115" x2="200" y2="115" className="ch18-svg-arrow-ii" markerEnd="url(#arrowhead)" />
          <text x="280" y="132" textAnchor="middle" className="ch18-svg-flow-label">IIc → I (for MP)  {flowLabel}</text>
          <text x="280" y="175" textAnchor="middle" className="ch18-svg-eq-label">
            Equilibrium: I(v+s) = IIc
          </text>
        </svg>
      </div>
      <p className="ch18-note">
        {settlement
          ? `Period ${settlement.period}: I(v+s)→IIc = ${flowLabel}; IIc→I = ${flowLabel}. Equal flows — simple reproduction.`
          : "Seed fixtures (1865): I(v+s)→IIc = £208; IIc→I = £208. Equal flows — simple reproduction."}
      </p>
    </div>
  );
}

// ── Root ─────────────────────────────────────────────────────────────────────

export function Ch18RoleOfMoneyCapital() {
  return (
    <div className="ch18-root">
      <ApportionmentPie />
      <MVCalculator />
      <WageRotationToy />
      <TwoDeptFlowDiagram />
    </div>
  );
}
