import { useState } from "react";
import "./Ch18RoleOfMoneyCapital.css";

// Vol. II Ch. 18 — The Role of Money-Capital in Reproduction.
//
// Marx's key insight: money-capital does not create value, but it is
// indispensable for keeping the circuit in motion. The two departments
// must exchange with each other, and money is the medium.
//
// The chapter analyses:
//   1. How the total circulating money-stock is apportioned across the
//      two departments, wage-rotation fund, and idle hoards.
//   2. The quantity theory correction: M × V = effective circulating value.
//   3. The wage-rotation fund: money that must be advanced before output is
//      sold, because wages are paid weekly/monthly, not annually.
//   4. The inter-department exchange flows: I(v+s)→II and II(c)→I.
//
// Seed fixtures (00055_v2_ch18_seed.sql, period "1865"):
//   Total M = £1 000; Dept I = 40%, Dept II = 35%, Wages = 20%, Idle = 5%
//   V = 2× per year (20000 BP); effective = £2 000
//   Both departments pay weekly (52 cycles); each fund = £104
//   I(v+s)→II = £208; II(c)→I = £208

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

const SEED_APPORTIONMENT = {
  total: 10000000,
  deptI: 4000000,
  deptII: 3500000,
  wage: 2000000,
  idle: 500000,
  period: "1865",
};

const APPORT_BARS: { label: string; pence: number; css: string }[] = [
  { label: "Dept I (means of production)",  pence: SEED_APPORTIONMENT.deptI,  css: "bar-dept-i"  },
  { label: "Dept II (articles of consumption)", pence: SEED_APPORTIONMENT.deptII, css: "bar-dept-ii" },
  { label: "Wage-rotation fund",            pence: SEED_APPORTIONMENT.wage,   css: "bar-wage"   },
  { label: "Idle hoard",                    pence: SEED_APPORTIONMENT.idle,   css: "bar-idle"   },
];

function ApportionmentPie() {
  const total = SEED_APPORTIONMENT.total;
  return (
    <div className="ch18-widget">
      <p className="ch18-section-title">
        Money-Stock Apportionment — {SEED_APPORTIONMENT.period}
      </p>
      <p className="ch18-widget-subtitle">
        Total M = {penceToPounds(total)} distributed across departments and funds
      </p>
      <div className="ch18-bars">
        {APPORT_BARS.map((b) => {
          const pct = Math.round((b.pence / total) * 100);
          return (
            <div key={b.label} className="ch18-bar-row">
              <span className="ch18-bar-label">{b.label}</span>
              <div className="ch18-bar-track">
                <div
                  className={`ch18-bar-fill ${b.css}`}
                  style={{ width: `${pct}%` }}
                />
              </div>
              <span className="ch18-bar-meta">
                {penceToPounds(b.pence)} ({pct}%)
              </span>
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
  // velocity expressed in turnovers/year; convert to basis points for display
  const velocityBP = Math.round(velocity * 10000);
  // effective = stock * velocityBP / 10000 (integer maths matches Go)
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
          <input
            type="number"
            min="0"
            className="ch18-mv-input"
            value={stockStr}
            onChange={(e) => setStockStr(e.target.value)}
          />
        </label>
        <span className="ch18-mv-times">×</span>
        <label className="ch18-mv-label">
          Velocity V (turnovers/year)
          <input
            type="number"
            min="0"
            step="0.5"
            className="ch18-mv-input"
            value={velocityStr}
            onChange={(e) => setVelocityStr(e.target.value)}
          />
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
        Velocity is expressed here as turnovers per year; the Go store records
        it in basis points (10 000 = 1×/yr).
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
  // The rotation fund ties up fund / freq for one wage cycle
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
          <input
            type="number"
            min="0"
            className="ch18-wr-input"
            value={fundStr}
            onChange={(e) => setFundStr(e.target.value)}
          />
        </label>
        <label className="ch18-wr-label">
          Payment cycle
          <select
            className="ch18-wr-select"
            value={freq}
            onChange={(e) => setFreq(Number(e.target.value))}
          >
            {FREQ_OPTIONS.map((o) => (
              <option key={o.value} value={o.value}>
                {o.label}
              </option>
            ))}
          </select>
        </label>
      </div>
      <div className="ch18-wr-result">
        <span className="ch18-wr-result-label">
          Money tied up per cycle
        </span>
        <span className="ch18-wr-result-value">
          £{tiedUpPerCycle.toFixed(2)}
        </span>
        <span className="ch18-wr-result-sub">
          {freq} payment{freq !== 1 ? "s" : ""} per year
        </span>
      </div>
      <p className="ch18-note">
        The more frequently wages are paid, the smaller the advance required
        at any one time — but the capitalist must still hold the fund until
        the next turnover completes.
      </p>
    </div>
  );
}

// ── Widget 4 — Two-Department Flow Diagram (static SVG) ──────────────────────

function TwoDeptFlowDiagram() {
  return (
    <div className="ch18-widget">
      <p className="ch18-section-title">Two-Department Exchange Flows</p>
      <p className="ch18-widget-subtitle">
        Simple reproduction requires that I(v+s) = II(c) — the exchange
        between departments must balance at the aggregate level.
      </p>
      <div className="ch18-flow-svg-wrapper">
        <svg
          className="ch18-flow-svg"
          viewBox="0 0 560 200"
          xmlns="http://www.w3.org/2000/svg"
          aria-label="Two-department flow diagram"
        >
          {/* Department I box */}
          <rect x="20" y="60" width="180" height="80" rx="8" className="ch18-svg-box-i" />
          <text x="110" y="93" textAnchor="middle" className="ch18-svg-label-main">Dept I</text>
          <text x="110" y="113" textAnchor="middle" className="ch18-svg-label-sub">Means of Production</text>
          <text x="110" y="130" textAnchor="middle" className="ch18-svg-label-sub">c + v + s</text>

          {/* Department II box */}
          <rect x="360" y="60" width="180" height="80" rx="8" className="ch18-svg-box-ii" />
          <text x="450" y="93" textAnchor="middle" className="ch18-svg-label-main">Dept II</text>
          <text x="450" y="113" textAnchor="middle" className="ch18-svg-label-sub">Articles of Consumption</text>
          <text x="450" y="130" textAnchor="middle" className="ch18-svg-label-sub">c + v + s</text>

          {/* Arrow I → II (top): I(v+s) → IIc */}
          <defs>
            <marker id="arrowhead" markerWidth="8" markerHeight="6" refX="8" refY="3" orient="auto">
              <polygon points="0 0, 8 3, 0 6" className="ch18-svg-arrow-marker" />
            </marker>
          </defs>
          <line
            x1="200" y1="85" x2="360" y2="85"
            className="ch18-svg-arrow-i"
            markerEnd="url(#arrowhead)"
          />
          <text x="280" y="78" textAnchor="middle" className="ch18-svg-flow-label">I(v+s) → IIc</text>

          {/* Arrow II → I (bottom): IIc → I for MP */}
          <line
            x1="360" y1="115" x2="200" y2="115"
            className="ch18-svg-arrow-ii"
            markerEnd="url(#arrowhead)"
          />
          <text x="280" y="132" textAnchor="middle" className="ch18-svg-flow-label">IIc → I (for MP)</text>

          {/* Equilibrium condition */}
          <text x="280" y="175" textAnchor="middle" className="ch18-svg-eq-label">
            Equilibrium: I(v+s) = IIc
          </text>
        </svg>
      </div>
      <p className="ch18-note">
        Seed fixtures (1865): I(v+s)→IIc = £208; IIc→I = £208. The flows are
        equal — simple reproduction. In extended reproduction (Ch. 21) I(v+s)
        exceeds IIc and the surplus is re-invested in Dept I.
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
