import { useState } from "react";
import "./Ch20SimpleReproduction.css";

// Vol. II Ch. 20 — Simple Reproduction.
//
// Marx derives the two-department reproduction scheme for total social capital:
//   Department I  — produces means of production (c+v+s)
//   Department II — produces articles of consumption (c+v+s)
//
// The keystone equation: I(v + s) == II(c)
//
// If this holds, society can reproduce itself on the same scale each year.
// No accumulation — next-period magnitudes are identical to current-period.
//
// Marx's canonical fixture (period "1865"):
//   Dept I:  4 000c + 1 000v + 1 000s = 6 000   I(v+s) = 2 000
//   Dept II: 2 000c +   500v +   500s = 3 000   II(c)  = 2 000  ✓
//
// Three exchange flows per cycle:
//   1. Intra-Dept-I  — I replaces its own c internally; no money crosses boundary
//   2. I(v+s) → II   — Dept I workers/capitalists buy consumption goods
//   3. II(c)  → I    — Dept II buys means of production from Dept I
//
// Chapter resolves the Ch. 17 realisation puzzle and refutes Smith's revenue
// dogma from Ch. 19.

// ── canonical fixture ─────────────────────────────────────────────────────────

const DEPT_I = { c: 4000, v: 1000, s: 1000, total: 6000 };
const DEPT_II = { c: 2000, v: 500, s: 500, total: 3000 };

// ── helpers ───────────────────────────────────────────────────────────────────

function fmt(n: number): string {
  return n.toLocaleString();
}

// ── Widget 1 — Canonical Two-Department Table ─────────────────────────────────

function DepartmentTable() {
  const ivs = DEPT_I.v + DEPT_I.s;
  const balanced = ivs === DEPT_II.c;
  return (
    <div className="ch20-widget">
      <p className="ch20-section-title">The Canonical Reproduction Scheme — 1865</p>
      <p className="ch20-subtitle">
        Values in abstract units. Backend stores pence ×&thinsp;1&thinsp;000.
      </p>
      <div className="ch20-table-wrap">
        <table className="ch20-table">
          <thead>
            <tr>
              <th>Department</th>
              <th className="col-c">c</th>
              <th className="col-v">v</th>
              <th className="col-s">s</th>
              <th className="col-total">Total</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td className="row-dept-i">I — Means of Production</td>
              <td className="col-c num">{fmt(DEPT_I.c)}</td>
              <td className="col-v num">{fmt(DEPT_I.v)}</td>
              <td className="col-s num">{fmt(DEPT_I.s)}</td>
              <td className="col-total num">{fmt(DEPT_I.total)}</td>
            </tr>
            <tr>
              <td className="row-dept-ii">II — Articles of Consumption</td>
              <td className="col-c num">{fmt(DEPT_II.c)}</td>
              <td className="col-v num">{fmt(DEPT_II.v)}</td>
              <td className="col-s num">{fmt(DEPT_II.s)}</td>
              <td className="col-total num">{fmt(DEPT_II.total)}</td>
            </tr>
          </tbody>
        </table>
      </div>
      <div className={`ch20-balance-badge ${balanced ? "balanced" : "imbalanced"}`}>
        <span className="ch20-balance-eq">
          I(v + s) = {fmt(ivs)}&nbsp;&nbsp;{balanced ? "==" : "≠"}&nbsp;&nbsp;II(c) = {fmt(DEPT_II.c)}
        </span>
        <span className="ch20-balance-verdict">
          {balanced ? "✓ Balanced — simple reproduction holds" : "✗ Imbalanced"}
        </span>
      </div>
    </div>
  );
}

// ── Widget 2 — Three-Exchange Flow Diagram (SVG) ──────────────────────────────

function ExchangeFlowDiagram() {
  return (
    <div className="ch20-widget">
      <p className="ch20-section-title">Three Exchange Flows per Reproduction Cycle</p>
      <p className="ch20-subtitle">
        The three flows that close the circuit each year. Money returns to its starting point.
      </p>
      <div className="ch20-flow-wrap">
        <svg
          viewBox="0 0 560 210"
          className="ch20-flow-svg"
          aria-label="Three-exchange flow diagram"
        >
          <defs>
            <marker id="arr-i" markerWidth="8" markerHeight="8" refX="6" refY="3" orient="auto">
              <path d="M0,0 L0,6 L8,3 z" fill="#818cf8" />
            </marker>
            <marker id="arr-ii" markerWidth="8" markerHeight="8" refX="6" refY="3" orient="auto">
              <path d="M0,0 L0,6 L8,3 z" fill="#34d399" />
            </marker>
            <marker id="arr-self" markerWidth="8" markerHeight="8" refX="6" refY="3" orient="auto">
              <path d="M0,0 L0,6 L8,3 z" fill="#a78bfa" />
            </marker>
          </defs>

          {/* Dept I box */}
          <rect x="20" y="70" width="175" height="75" rx="6" className="ch20-svg-box-i" />
          <text x="107" y="99" textAnchor="middle" className="ch20-svg-label-main">Dept I</text>
          <text x="107" y="114" textAnchor="middle" className="ch20-svg-label-sub">4000c + 1000v + 1000s</text>
          <text x="107" y="129" textAnchor="middle" className="ch20-svg-label-sub">Total = 6000</text>

          {/* Dept II box */}
          <rect x="365" y="70" width="175" height="75" rx="6" className="ch20-svg-box-ii" />
          <text x="453" y="99" textAnchor="middle" className="ch20-svg-label-main">Dept II</text>
          <text x="453" y="114" textAnchor="middle" className="ch20-svg-label-sub">2000c + 500v + 500s</text>
          <text x="453" y="129" textAnchor="middle" className="ch20-svg-label-sub">Total = 3000</text>

          {/* Flow ① — Intra-Dept-I self-loop */}
          <path d="M 75 70 C 40 25, 145 25, 145 70" stroke="#a78bfa" strokeWidth="2" fill="none" markerEnd="url(#arr-self)" />
          <text x="107" y="19" textAnchor="middle" className="ch20-svg-flow-label ch20-purple">① Intra-I: c = 4000</text>

          {/* Flow ② — I(v+s) → II */}
          <path d="M 195 90 L 365 90" stroke="#818cf8" strokeWidth="2" fill="none" markerEnd="url(#arr-i)" />
          <text x="280" y="80" textAnchor="middle" className="ch20-svg-flow-label ch20-indigo">② I(v+s) → II: 2000</text>

          {/* Flow ③ — II(c) → I */}
          <path d="M 365 125 L 195 125" stroke="#34d399" strokeWidth="2" fill="none" markerEnd="url(#arr-ii)" />
          <text x="280" y="150" textAnchor="middle" className="ch20-svg-flow-label ch20-green">③ II(c) → I: 2000</text>

          {/* Keystone equation */}
          <text x="280" y="195" textAnchor="middle" className="ch20-svg-eq-label">
            I(v+s) = 2000 == II(c) = 2000 ✓
          </text>
        </svg>
      </div>
      <div className="ch20-exchange-list">
        <div className="ch20-exchange-row ch20-purple">
          <span className="ch20-exch-num">①</span>
          <span className="ch20-exch-body">
            <strong>Intra-Dept-I replacement</strong> — Dept I buys MP from itself (I.c = 4&thinsp;000).
            No money crosses the departmental boundary.
          </span>
        </div>
        <div className="ch20-exchange-row ch20-indigo">
          <span className="ch20-exch-num">②</span>
          <span className="ch20-exch-body">
            <strong>I(v+s) → II</strong> — Dept I workers (v) and capitalists (s) purchase
            articles of consumption. Money flows I→II; goods flow II→I. Value = 2&thinsp;000.
          </span>
        </div>
        <div className="ch20-exchange-row ch20-green">
          <span className="ch20-exch-num">③</span>
          <span className="ch20-exch-body">
            <strong>II(c) → I</strong> — Dept II purchases replacement MP from Dept I.
            Money flows II→I; goods flow I→II. Value = 2&thinsp;000.
          </span>
        </div>
      </div>
    </div>
  );
}

// ── Widget 3 — Tick Demonstrator ──────────────────────────────────────────────

function TickDemonstrator() {
  const [tick, setTick] = useState(0);
  return (
    <div className="ch20-widget">
      <p className="ch20-section-title">Simple Reproduction Tick — Identical Magnitudes</p>
      <p className="ch20-subtitle">
        No accumulation: each tick the scheme repeats with the exact same magnitudes.
      </p>
      <div className="ch20-tick-display">
        <div className="ch20-tick-period">
          Period: <strong>1865 + {tick}</strong>
        </div>
        <div className="ch20-tick-row">
          <span className="dept-i-tag">I</span>
          <span className="ch20-tick-formula">
            {fmt(DEPT_I.c)}c + {fmt(DEPT_I.v)}v + {fmt(DEPT_I.s)}s = {fmt(DEPT_I.total)}
          </span>
        </div>
        <div className="ch20-tick-row">
          <span className="dept-ii-tag">II</span>
          <span className="ch20-tick-formula">
            {fmt(DEPT_II.c)}c + {fmt(DEPT_II.v)}v + {fmt(DEPT_II.s)}s = {fmt(DEPT_II.total)}
          </span>
        </div>
        <div className="ch20-tick-note">Magnitudes unchanged — no surplus converted to capital.</div>
      </div>
      <button className="ch20-tick-btn" onClick={() => setTick((t) => t + 1)}>
        Advance one reproduction cycle →
      </button>
      {tick > 0 && (
        <p className="ch20-tick-confirm">
          After {tick} cycle{tick !== 1 ? "s" : ""}: Dept I still {fmt(DEPT_I.total)}, Dept II
          still {fmt(DEPT_II.total)}. Simple reproduction holds.
        </p>
      )}
    </div>
  );
}

// ── Widget 4 — Imbalance Toy ──────────────────────────────────────────────────

function ImbalanceToy() {
  const [ivs, setIvs] = useState(DEPT_I.v + DEPT_I.s);
  const [iic, setIic] = useState(DEPT_II.c);
  const deficit = ivs - iic;
  const balanced = deficit === 0;

  return (
    <div className="ch20-widget">
      <p className="ch20-section-title">Imbalance Toy — Break the Keystone Equation</p>
      <p className="ch20-subtitle">
        Adjust I(v+s) or II(c) to see what happens when the equation fails.
      </p>
      <div className="ch20-imbal-form">
        <label className="ch20-imbal-label">
          I(v + s)
          <input
            type="number"
            className="ch20-imbal-input"
            value={ivs}
            min={0}
            step={100}
            onChange={(e) => setIvs(Number(e.target.value))}
          />
        </label>
        <span className="ch20-imbal-op">{balanced ? "==" : "≠"}</span>
        <label className="ch20-imbal-label">
          II(c)
          <input
            type="number"
            className="ch20-imbal-input"
            value={iic}
            min={0}
            step={100}
            onChange={(e) => setIic(Number(e.target.value))}
          />
        </label>
      </div>
      <div className={`ch20-imbal-result ${balanced ? "ok" : "fail"}`}>
        {balanced ? (
          <>✓ Balanced — simple reproduction is possible.</>
        ) : deficit > 0 ? (
          <>
            ✗ Surplus of {fmt(deficit)}: Dept I overproduces MP. Dept II cannot absorb all
            I(v+s). Realisation fails — cf. Ch.&nbsp;17 realisation puzzle.
          </>
        ) : (
          <>
            ✗ Deficit of {fmt(Math.abs(deficit))}: II.c exceeds what Dept I produces as v+s.
            Dept II cannot replace worn-out MP. Production contracts next period.
          </>
        )}
      </div>
    </div>
  );
}

// ── Widget 5 — Money Closed Loop ──────────────────────────────────────────────

function MoneyClosedLoop() {
  const advanced = DEPT_II.c;
  const returnedV = DEPT_I.v;
  const returnedS = DEPT_I.s;
  const returned = returnedV + returnedS;
  const net = advanced - returned;
  return (
    <div className="ch20-widget">
      <p className="ch20-section-title">Money Closed Loop — Verification</p>
      <p className="ch20-subtitle">
        Money advanced by Dept II returns to Dept II within the year, closing the circuit.
      </p>
      <div className="ch20-loop-table">
        <div className="ch20-loop-row">
          <span>Dept II advances money to buy MP from Dept I (II.c)</span>
          <span className="ch20-loop-val out">−{fmt(advanced)}</span>
        </div>
        <div className="ch20-loop-row">
          <span>Returned via Dept I workers purchasing consumption goods (I.v → II)</span>
          <span className="ch20-loop-val in">+{fmt(returnedV)}</span>
        </div>
        <div className="ch20-loop-row">
          <span>Returned via Dept I capitalists purchasing consumption goods (I.s → II)</span>
          <span className="ch20-loop-val in">+{fmt(returnedS)}</span>
        </div>
        <div className="ch20-loop-row total">
          <span>Net flow back to Dept II</span>
          <span className={`ch20-loop-val ${net === 0 ? "zero" : "nonzero"}`}>
            {net === 0 ? "0 — loop closed ✓" : fmt(net)}
          </span>
        </div>
      </div>
      <p className="ch20-note">
        The money circuit is self-closing within the year. No new money supply is required
        for simple reproduction to continue indefinitely.
      </p>
    </div>
  );
}

// ── Root export ───────────────────────────────────────────────────────────────

export function Ch20SimpleReproduction() {
  return (
    <div className="ch20-root">
      <DepartmentTable />
      <ExchangeFlowDiagram />
      <TickDemonstrator />
      <ImbalanceToy />
      <MoneyClosedLoop />
    </div>
  );
}
