import { useEffect, useState } from "react";
import { simpleReproductionApi } from "../../api";
import type {
  SimpleReproductionScheme,
  DepartmentalCapital,
  BalanceCheckResult,
  ReproductionTick,
} from "../../types";
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

// ── canonical fixture (fallback) ──────────────────────────────────────────────

const DEPT_I_FALLBACK = { c: 4000, v: 1000, s: 1000, total: 6000 };
const DEPT_II_FALLBACK = { c: 2000, v: 500, s: 500, total: 3000 };

// ── helpers ───────────────────────────────────────────────────────────────────

function fmt(n: number): string {
  return Math.round(n).toLocaleString();
}

// Backend stores pence × 1 000; abstract units = pence / 1 000.
function toAbstract(pence: number): number {
  return Math.round(pence / 1000);
}

function deptToAbstract(d: DepartmentalCapital) {
  return {
    c: toAbstract(d.constant_pence),
    v: toAbstract(d.variable_pence),
    s: toAbstract(d.surplus_pence),
    total: toAbstract(d.total_pence),
  };
}

// ── Widget 1 — Canonical Two-Department Table ─────────────────────────────────

interface DepartmentTableProps {
  scheme: SimpleReproductionScheme | null;
}

function DepartmentTable({ scheme }: DepartmentTableProps) {
  const dI   = scheme?.department_i   ? deptToAbstract(scheme.department_i)   : DEPT_I_FALLBACK;
  const dII  = scheme?.department_ii  ? deptToAbstract(scheme.department_ii)  : DEPT_II_FALLBACK;
  const period = scheme?.period ?? "1865";
  const ivs = dI.v + dI.s;
  const balanced = ivs === dII.c;
  return (
    <div className="ch20-widget">
      <p className="ch20-section-title">The Canonical Reproduction Scheme — {period}</p>
      <p className="ch20-subtitle">
        Values in abstract units. Backend stores pence ×&thinsp;1&thinsp;000.
        {!scheme && " (displaying local fixture — connect to backend to load seed data)"}
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
              <td className="col-c num">{fmt(dI.c)}</td>
              <td className="col-v num">{fmt(dI.v)}</td>
              <td className="col-s num">{fmt(dI.s)}</td>
              <td className="col-total num">{fmt(dI.total)}</td>
            </tr>
            <tr>
              <td className="row-dept-ii">II — Articles of Consumption</td>
              <td className="col-c num">{fmt(dII.c)}</td>
              <td className="col-v num">{fmt(dII.v)}</td>
              <td className="col-s num">{fmt(dII.s)}</td>
              <td className="col-total num">{fmt(dII.total)}</td>
            </tr>
          </tbody>
        </table>
      </div>
      <div className={`ch20-balance-badge ${balanced ? "balanced" : "imbalanced"}`}>
        <span className="ch20-balance-eq">
          I(v + s) = {fmt(ivs)}&nbsp;&nbsp;{balanced ? "==" : "≠"}&nbsp;&nbsp;II(c) = {fmt(dII.c)}
        </span>
        <span className="ch20-balance-verdict">
          {balanced ? "✓ Balanced — simple reproduction holds" : "✗ Imbalanced"}
        </span>
      </div>
    </div>
  );
}

// ── Widget 2 — Three-Exchange Flow Diagram (SVG) ──────────────────────────────

interface ExchangeFlowProps {
  scheme: SimpleReproductionScheme | null;
}

function ExchangeFlowDiagram({ scheme }: ExchangeFlowProps) {
  const dI  = scheme?.department_i  ? deptToAbstract(scheme.department_i)  : DEPT_I_FALLBACK;
  const dII = scheme?.department_ii ? deptToAbstract(scheme.department_ii) : DEPT_II_FALLBACK;
  const ivs = dI.v + dI.s;
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
          <text x="107" y="114" textAnchor="middle" className="ch20-svg-label-sub">{fmt(dI.c)}c + {fmt(dI.v)}v + {fmt(dI.s)}s</text>
          <text x="107" y="129" textAnchor="middle" className="ch20-svg-label-sub">Total = {fmt(dI.total)}</text>

          {/* Dept II box */}
          <rect x="365" y="70" width="175" height="75" rx="6" className="ch20-svg-box-ii" />
          <text x="453" y="99" textAnchor="middle" className="ch20-svg-label-main">Dept II</text>
          <text x="453" y="114" textAnchor="middle" className="ch20-svg-label-sub">{fmt(dII.c)}c + {fmt(dII.v)}v + {fmt(dII.s)}s</text>
          <text x="453" y="129" textAnchor="middle" className="ch20-svg-label-sub">Total = {fmt(dII.total)}</text>

          {/* Flow ① — Intra-Dept-I self-loop */}
          <path d="M 75 70 C 40 25, 145 25, 145 70" stroke="#a78bfa" strokeWidth="2" fill="none" markerEnd="url(#arr-self)" />
          <text x="107" y="19" textAnchor="middle" className="ch20-svg-flow-label ch20-purple">① Intra-I: c = {fmt(dI.c)}</text>

          {/* Flow ② — I(v+s) → II */}
          <path d="M 195 90 L 365 90" stroke="#818cf8" strokeWidth="2" fill="none" markerEnd="url(#arr-i)" />
          <text x="280" y="80" textAnchor="middle" className="ch20-svg-flow-label ch20-indigo">② I(v+s) → II: {fmt(ivs)}</text>

          {/* Flow ③ — II(c) → I */}
          <path d="M 365 125 L 195 125" stroke="#34d399" strokeWidth="2" fill="none" markerEnd="url(#arr-ii)" />
          <text x="280" y="150" textAnchor="middle" className="ch20-svg-flow-label ch20-green">③ II(c) → I: {fmt(dII.c)}</text>

          {/* Keystone equation */}
          <text x="280" y="195" textAnchor="middle" className="ch20-svg-eq-label">
            I(v+s) = {fmt(ivs)} {ivs === dII.c ? "==" : "≠"} II(c) = {fmt(dII.c)} {ivs === dII.c ? "✓" : "✗"}
          </text>
        </svg>
      </div>
      <div className="ch20-exchange-list">
        <div className="ch20-exchange-row ch20-purple">
          <span className="ch20-exch-num">①</span>
          <span className="ch20-exch-body">
            <strong>Intra-Dept-I replacement</strong> — Dept I buys MP from itself (I.c = {fmt(dI.c)}).
            No money crosses the departmental boundary.
          </span>
        </div>
        <div className="ch20-exchange-row ch20-indigo">
          <span className="ch20-exch-num">②</span>
          <span className="ch20-exch-body">
            <strong>I(v+s) → II</strong> — Dept I workers (v) and capitalists (s) purchase
            articles of consumption. Money flows I→II; goods flow II→I. Value = {fmt(ivs)}.
          </span>
        </div>
        <div className="ch20-exchange-row ch20-green">
          <span className="ch20-exch-num">③</span>
          <span className="ch20-exch-body">
            <strong>II(c) → I</strong> — Dept II purchases replacement MP from Dept I.
            Money flows II→I; goods flow I→II. Value = {fmt(dII.c)}.
          </span>
        </div>
      </div>
    </div>
  );
}

// ── Widget 3 — Tick Demonstrator ──────────────────────────────────────────────

interface TickDemonstratorProps {
  scheme: SimpleReproductionScheme | null;
  onTick: (updated: SimpleReproductionScheme) => void;
}

function TickDemonstrator({ scheme, onTick }: TickDemonstratorProps) {
  const [ticking, setTicking] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [workers, setWorkers] = useState(0);
  const [basket, setBasket] = useState(0); // abstract units per worker
  const [lastTick, setLastTick] = useState<ReproductionTick | null>(null);
  const dI  = scheme?.department_i  ? deptToAbstract(scheme.department_i)  : DEPT_I_FALLBACK;
  const dII = scheme?.department_ii ? deptToAbstract(scheme.department_ii) : DEPT_II_FALLBACK;
  const tickCount = scheme?.tick_count ?? 0;

  async function handleTick() {
    if (!scheme) return;
    setTicking(true);
    setError(null);
    try {
      const tick = await simpleReproductionApi.advanceTick(scheme.id, {
        worker_pool_size: workers,
        subsistence_basket_pence: basket * 1000,
      });
      setLastTick(tick);
      const updated = await simpleReproductionApi.getScheme(scheme.id);
      onTick(updated);
    } catch {
      setError("Failed to advance tick — is the backend running?");
    } finally {
      setTicking(false);
    }
  }

  return (
    <div className="ch20-widget">
      <p className="ch20-section-title">Simple Reproduction Tick — Identical Magnitudes</p>
      <p className="ch20-subtitle">
        No accumulation: each tick the scheme repeats with the exact same magnitudes.
        {scheme
          ? ` Tick ${tickCount} recorded on the backend.`
          : " (local simulation — connect to backend to persist ticks)"}
      </p>
      <div className="ch20-tick-display">
        <div className="ch20-tick-period">
          Period: <strong>{scheme?.period ?? "1865"}</strong>
          {tickCount > 0 && (
            <span className="ch20-tick-badge">
              {" "}+ {tickCount} cycle{tickCount !== 1 ? "s" : ""}
            </span>
          )}
        </div>
        <div className="ch20-tick-row">
          <span className="dept-i-tag">I</span>
          <span className="ch20-tick-formula">
            {fmt(dI.c)}c + {fmt(dI.v)}v + {fmt(dI.s)}s = {fmt(dI.total)}
          </span>
        </div>
        <div className="ch20-tick-row">
          <span className="dept-ii-tag">II</span>
          <span className="ch20-tick-formula">
            {fmt(dII.c)}c + {fmt(dII.v)}v + {fmt(dII.s)}s = {fmt(dII.total)}
          </span>
        </div>
        <div className="ch20-tick-note">Magnitudes unchanged — no surplus converted to capital.</div>
      </div>
      <div className="ch20-labour-inputs">
        <label>
          Worker pool
          <input
            type="number"
            min={0}
            value={workers}
            onChange={(e) => setWorkers(Math.max(0, Math.floor(Number(e.target.value))))}
          />
        </label>
        <label>
          Subsistence / worker
          <input
            type="number"
            min={0}
            value={basket}
            onChange={(e) => setBasket(Math.max(0, Math.floor(Number(e.target.value))))}
          />
        </label>
      </div>
      <button
        className="ch20-tick-btn"
        onClick={handleTick}
        disabled={!scheme || ticking}
      >
        {ticking ? "Advancing…" : "Advance one reproduction cycle →"}
      </button>
      {error && <p className="ch20-tick-error">{error}</p>}
      {lastTick && lastTick.worker_pool_size > 0 && (
        <p className={`ch20-labour-status ${lastTick.subsistence_covered ? "ok" : "fail"}`}>
          Labour-force: wage bill {fmt(toAbstract(lastTick.wage_bill_pence))} vs required{" "}
          {fmt(lastTick.worker_pool_size * toAbstract(lastTick.subsistence_basket_pence))} (
          {lastTick.worker_pool_size} worker{lastTick.worker_pool_size !== 1 ? "s" : ""}) —{" "}
          {lastTick.subsistence_covered ? "✓ reproduced" : "✗ labour-force shortfall"}
        </p>
      )}
      {tickCount > 0 && (
        <p className="ch20-tick-confirm">
          After {tickCount} cycle{tickCount !== 1 ? "s" : ""}: Dept I still {fmt(dI.total)}, Dept II
          still {fmt(dII.total)}. Simple reproduction holds.
        </p>
      )}
    </div>
  );
}

// ── Widget 4 — Imbalance Toy ──────────────────────────────────────────────────

interface ImbalanceToyProps {
  scheme: SimpleReproductionScheme | null;
  balance: BalanceCheckResult | null;
}

function ImbalanceToy({ scheme, balance }: ImbalanceToyProps) {
  const initIvs = scheme?.department_i
    ? toAbstract(scheme.department_i.variable_pence + scheme.department_i.surplus_pence)
    : DEPT_I_FALLBACK.v + DEPT_I_FALLBACK.s;
  const initIIc = scheme?.department_ii
    ? toAbstract(scheme.department_ii.constant_pence)
    : DEPT_II_FALLBACK.c;

  const [ivs, setIvs] = useState(initIvs);
  const [iic, setIic] = useState(initIIc);
  const deficit = ivs - iic;
  const balanced = deficit === 0;

  // Mirror backend balance check result when sliders are at seed values.
  const showBackendBadge = balance !== null && ivs === initIvs && iic === initIIc;

  return (
    <div className="ch20-widget">
      <p className="ch20-section-title">Imbalance Toy — Break the Keystone Equation</p>
      <p className="ch20-subtitle">
        Adjust I(v+s) or II(c) to see what happens when the equation fails.
        {showBackendBadge &&
          ` Backend confirms: ${balance.is_balanced ? "balanced" : "imbalanced"}.`}
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

interface MoneyClosedLoopProps {
  scheme: SimpleReproductionScheme | null;
}

function MoneyClosedLoop({ scheme }: MoneyClosedLoopProps) {
  const dI  = scheme?.department_i  ? deptToAbstract(scheme.department_i)  : DEPT_I_FALLBACK;
  const dII = scheme?.department_ii ? deptToAbstract(scheme.department_ii) : DEPT_II_FALLBACK;
  const advanced = dII.c;
  const returnedV = dI.v;
  const returnedS = dI.s;
  const returned = returnedV + returnedS;
  const net = advanced - returned;
  return (
    <div className="ch20-widget">
      <p className="ch20-section-title">Money Closed Loop — Verification</p>
      <p className="ch20-subtitle">
        Money advanced by Dept II returns to Dept II within the year, closing the circuit.
        {scheme && <> Period: <strong>{scheme.period}</strong>.</>}
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

// ── Widget — Annualised Settlement (Ch. 14 turnover) ──────────────────────────

interface AnnualisedSettlementProps {
  balance: BalanceCheckResult | null;
}

function AnnualisedSettlement({ balance }: AnnualisedSettlementProps) {
  if (!balance) return null;
  const turns = (bp: number) =>
    (bp / 10000).toLocaleString(undefined, { maximumFractionDigits: 2 });
  return (
    <div className="ch20-widget">
      <p className="ch20-section-title">Annualised Settlement — Turnover (Ch. 14)</p>
      <p>
        When the two departments turn over at different speeds, the keystone
        exchange settles on annual flows rather than raw per-period magnitudes —
        a scheme balanced each period can be unbalanced over the year.
      </p>
      <table className="ch20-annual-table">
        <thead>
          <tr>
            <th scope="col"></th>
            <th scope="col">Dept I — I(v+s)</th>
            <th scope="col">Dept II — II(c)</th>
          </tr>
        </thead>
        <tbody>
          <tr>
            <th scope="row">Turnover / yr</th>
            <td>{turns(balance.dept_i_turnover_bp)}×</td>
            <td>{turns(balance.dept_ii_turnover_bp)}×</td>
          </tr>
          <tr>
            <th scope="row">Raw (per period)</th>
            <td>{fmt(toAbstract(balance.dept_i_v_plus_s_pence))}</td>
            <td>{fmt(toAbstract(balance.dept_ii_constant_pence))}</td>
          </tr>
          <tr>
            <th scope="row">Annualised</th>
            <td>{fmt(toAbstract(balance.dept_i_v_plus_s_annualised_pence))}</td>
            <td>{fmt(toAbstract(balance.dept_ii_constant_annualised_pence))}</td>
          </tr>
        </tbody>
      </table>
      <div className={`ch20-balance-badge ${balance.annualised_balanced ? "balanced" : "imbalanced"}`}>
        <span className="ch20-balance-verdict">
          {balance.annualised_balanced
            ? "✓ Annually balanced"
            : `✗ Annual deficit of ${fmt(toAbstract(Math.abs(balance.annualised_deficit_pence)))}`}
        </span>
      </div>
    </div>
  );
}

// ── Widget — Build from Annual Rates (Ch. 16 → Ch. 20, issue #222) ────────────

function BuildFromAnnualRates() {
  const [rateI, setRateI] = useState(100);
  const [rateII, setRateII] = useState(100);
  const [result, setResult] = useState<SimpleReproductionScheme | null>(null);
  const [building, setBuilding] = useState(false);
  const [error, setError] = useState<string | null>(null);

  async function handleBuild() {
    setBuilding(true);
    setError(null);
    try {
      const scheme = await simpleReproductionApi.buildFromAnnualRates({
        period: "from-annual-rates",
        departments: [
          { department: "department_i", constant_pence: 4000000, variable_pence: 1000000, annual_surplus_bp: Math.round(rateI * 100) },
          { department: "department_ii", constant_pence: 2000000, variable_pence: 500000, annual_surplus_bp: Math.round(rateII * 100) },
        ],
      });
      setResult(scheme);
    } catch {
      setError("Failed to build — is the backend running?");
    } finally {
      setBuilding(false);
    }
  }

  return (
    <div className="ch20-widget">
      <p className="ch20-section-title">Build from Annual Rates (Ch. 16 → Ch. 20)</p>
      <p>
        An alternative entry point: derive the scheme from Ch. 16 annual rates of
        surplus-value instead of entering surplus directly. Each department's
        surplus is v × annual-rate, so the two layers cannot drift.
      </p>
      <div className="ch20-labour-inputs">
        <label>
          Dept I annual rate %
          <input type="number" min={0} value={rateI}
            onChange={(e) => setRateI(Math.max(0, Number(e.target.value)))} />
        </label>
        <label>
          Dept II annual rate %
          <input type="number" min={0} value={rateII}
            onChange={(e) => setRateII(Math.max(0, Number(e.target.value)))} />
        </label>
      </div>
      <button className="ch20-tick-btn" onClick={handleBuild} disabled={building}>
        {building ? "Building…" : "Build scheme from annual rates →"}
      </button>
      {error && <p className="ch20-tick-error">{error}</p>}
      {result && (
        <div className={`ch20-balance-badge ${result.is_balanced ? "balanced" : "imbalanced"}`}>
          <span className="ch20-balance-eq">
            I(v + s) = {fmt(toAbstract((result.department_i?.variable_pence ?? 0) + (result.department_i?.surplus_pence ?? 0)))}
            &nbsp;&nbsp;{result.is_balanced ? "==" : "≠"}&nbsp;&nbsp;
            II(c) = {fmt(toAbstract(result.department_ii?.constant_pence ?? 0))}
          </span>
          <span className="ch20-balance-verdict">
            {result.is_balanced
              ? "✓ Balanced — built from annual rates"
              : "✗ Imbalanced at these rates"}
          </span>
        </div>
      )}
    </div>
  );
}

// ── Root export ───────────────────────────────────────────────────────────────

export function Ch20SimpleReproduction() {
  const [scheme, setScheme] = useState<SimpleReproductionScheme | null>(null);
  const [balance, setBalance] = useState<BalanceCheckResult | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    simpleReproductionApi
      .listSchemes()
      .then(async (list) => {
        if (list.length === 0) return;
        const s = list[0];
        setScheme(s);
        simpleReproductionApi
          .balanceCheck(s.id)
          .then(setBalance)
          .catch(() => null);
      })
      .catch(() => null)
      .finally(() => setLoading(false));
  }, []);

  function handleTick(updated: SimpleReproductionScheme) {
    setScheme(updated);
  }

  if (loading) {
    return (
      <div className="ch20-root">
        <div className="ch20-widget">
          <p className="ch20-loading">Loading…</p>
        </div>
      </div>
    );
  }

  if (!scheme) {
    return (
      <div className="ch20-root">
        <div className="ch20-widget">
          <p className="ch20-section-title">Vol. II Ch. 20 — Simple Reproduction</p>
          <p className="ch20-empty">
            No seed data found. Run the database migrations to populate this panel.
          </p>
        </div>
      </div>
    );
  }

  return (
    <div className="ch20-root">
      <DepartmentTable scheme={scheme} />
      <ExchangeFlowDiagram scheme={scheme} />
      <TickDemonstrator scheme={scheme} onTick={handleTick} />
      <ImbalanceToy scheme={scheme} balance={balance} />
      <AnnualisedSettlement balance={balance} />
      <BuildFromAnnualRates />
      <MoneyClosedLoop scheme={scheme} />
    </div>
  );
}
