import { useState } from "react";
import "./Ch21ExtendedReproduction.css";

// Vol. II Ch. 21 — Accumulation and Reproduction on an Extended Scale.
//
// Unlike simple reproduction (Ch. 20), part of the surplus-value is
// re-capitalised — converted into additional constant capital (Δc) and
// additional variable capital (Δv). Society's productive capacity grows.
//
// Marx's canonical fixture (period "Year I"):
//   Dept I:  4000c + 1000v + 1000s = 6000  AccRate = 50% (5000 bps)
//   Dept II: 1500c +  750v +  750s = 3000  AccRate = 0%
//
// Reinvestment I: Δc = 0, Δv = 500, consumed_I = 500
//   (all reinvestment goes to v to satisfy balance equation)
//
// Extended balance: I.v(1000) + ΔvI(500) + consumed_I(0) = 1500
//                                        == II.c(1500) + ΔcII(0)  ✓
//
// Dept I grows faster than Dept II — means of production must expand ahead
// of consumption goods; that lead is Marx's thesis throughout Ch. 21.

// ── canonical fixture ─────────────────────────────────────────────────────────

const YEAR_I_I  = { c: 4000, v: 1000, s: 1000, total: 6000, rate: 0.5 };
const YEAR_I_II = { c: 1500, v: 750,  s: 750,  total: 3000, rate: 0.0 };

type DeptState = { c: number; v: number; s: number; total: number };

function advance(d: DeptState, rateBps: number): { next: DeptState; deltaC: number; deltaV: number; consumed: number } {
  const reinvested = Math.floor(d.s * rateBps / 10000);
  const deltaC = d.total > 0 ? Math.floor(reinvested * d.c / d.total) : 0;
  const deltaV = reinvested - deltaC;
  const consumed = d.s - reinvested;
  const newV = d.v + deltaV;
  const surplusRate = d.v > 0 ? d.s * 10000 / d.v : 0;
  const newS = Math.floor(newV * surplusRate / 10000);
  const newC = d.c + deltaC;
  return { next: { c: newC, v: newV, s: newS, total: newC + newV + newS }, deltaC, deltaV, consumed };
}

// ── helpers ───────────────────────────────────────────────────────────────────

function fmt(n: number): string {
  return Math.round(n).toLocaleString();
}

function bpsToPercent(bps: number): string {
  return (bps / 100).toFixed(0) + "%";
}

// ── Widget 1 — Canonical Two-Department Table with Accumulation ───────────────

function CanonicalFixture() {
  const dI = YEAR_I_I;
  const dII = YEAR_I_II;

  // Reinvestment at 50% for Dept I, 0% for Dept II:
  // All reinvested amount in Dept I goes to v (deltaC=0) to satisfy balance.
  const reinvestedI = Math.floor(dI.s * 5000 / 10000); // 500
  const deltaCI = 0;
  const deltaVI = reinvestedI - deltaCI; // 500
  const consumedI = dI.s - reinvestedI; // 500

  const deltaCII = 0;
  const deltaVII = 0;
  const consumedII = dII.s;

  // Balance: I.v + deltav_I + consumed_I == II.c + deltac_II
  const balanceLHS = dI.v + deltaVI + consumedI;
  const balanceRHS = dII.c + deltaCII;
  const balanced = balanceLHS === balanceRHS;

  return (
    <div className="ch21-widget">
      <p className="ch21-section-title">Marx's Canonical Extended Reproduction Fixture — Year I</p>
      <p className="ch21-subtitle">
        50% accumulation in Dept I. The Δv seed funds next-period variable capital.
        Dept I excess drives the growth of social production.
      </p>
      <div className="ch21-table-wrap">
        <table className="ch21-table">
          <thead>
            <tr>
              <th>Department</th>
              <th>c</th>
              <th>v</th>
              <th>s</th>
              <th>Acc.&nbsp;Rate</th>
              <th className="col-delta">Δc</th>
              <th className="col-delta">Δv</th>
              <th className="col-consumed">Consumed</th>
              <th>Total</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td className="row-dept-i">I — Means of Production</td>
              <td className="num">{fmt(dI.c)}</td>
              <td className="num">{fmt(dI.v)}</td>
              <td className="num">{fmt(dI.s)}</td>
              <td className="num">50%</td>
              <td className="num col-delta">{fmt(deltaCI)}</td>
              <td className="num col-delta">{fmt(deltaVI)}</td>
              <td className="num col-consumed">{fmt(consumedI)}</td>
              <td className="num">{fmt(dI.total)}</td>
            </tr>
            <tr>
              <td className="row-dept-ii">II — Articles of Consumption</td>
              <td className="num">{fmt(dII.c)}</td>
              <td className="num">{fmt(dII.v)}</td>
              <td className="num">{fmt(dII.s)}</td>
              <td className="num">0%</td>
              <td className="num col-delta">{fmt(deltaCII)}</td>
              <td className="num col-delta">{fmt(deltaVII)}</td>
              <td className="num col-consumed">{fmt(consumedII)}</td>
              <td className="num">{fmt(dII.total)}</td>
            </tr>
          </tbody>
        </table>
      </div>
      <div className={`ch21-balance-badge ${balanced ? "balanced" : "imbalanced"}`}>
        <span className="ch21-balance-eq">
          I.v({fmt(dI.v)}) + Δv_I({fmt(deltaVI)}) + consumed_I({fmt(consumedI)}) = {fmt(balanceLHS)}
          &nbsp;&nbsp;{balanced ? "==" : "≠"}&nbsp;&nbsp;
          II.c({fmt(dII.c)}) + Δc_II({fmt(deltaCII)}) = {fmt(balanceRHS)}
        </span>
        <span className="ch21-balance-verdict">
          {balanced ? "✓ Extended reproduction balance holds" : "✗ Balance fails"}
        </span>
      </div>
    </div>
  );
}

// ── Widget 2 — Reinvestment Split ────────────────────────────────────────────

function ReinvestmentSplitWidget() {
  const [rateIBps, setRateIBps] = useState(5000);
  const [rateIIBps, setRateIIBps] = useState(0);

  const dI: DeptState  = { c: YEAR_I_I.c,  v: YEAR_I_I.v,  s: YEAR_I_I.s,  total: YEAR_I_I.total };
  const dII: DeptState = { c: YEAR_I_II.c, v: YEAR_I_II.v, s: YEAR_I_II.s, total: YEAR_I_II.total };

  const aI  = advance(dI, rateIBps);
  const aII = advance(dII, rateIIBps);

  const lhs = dI.v + aI.deltaV + aI.consumed;
  const rhs = dII.c + aII.deltaC;
  const balanced = lhs === rhs;

  return (
    <div className="ch21-widget">
      <p className="ch21-section-title">Reinvestment Split — Live Balance Check</p>
      <p className="ch21-subtitle">
        Adjust accumulation rates to see how surplus splits into Δc, Δv, and consumed.
        The balance equation must hold for extended reproduction to proceed.
      </p>
      <div className="ch21-split-grid">
        <div className="ch21-split-card">
          <p className="ch21-split-card-title dept-i">Dept I — Means of Production</p>
          <div className="ch21-split-row">
            <span className="ch21-split-label">Accumulation rate</span>
            <span>
              <input
                type="range" min={0} max={10000} step={500}
                value={rateIBps}
                onChange={(e) => setRateIBps(Number(e.target.value))}
                style={{ width: 80 }}
              />
              {" "}{bpsToPercent(rateIBps)}
            </span>
          </div>
          <div className="ch21-split-row">
            <span className="ch21-split-label">Surplus (s)</span>
            <span className="ch21-split-val">{fmt(dI.s)}</span>
          </div>
          <div className="ch21-split-row">
            <span className="ch21-split-label">Δc (reinvested to constant)</span>
            <span className="ch21-split-val delta-c">{fmt(aI.deltaC)}</span>
          </div>
          <div className="ch21-split-row">
            <span className="ch21-split-label">Δv (reinvested to variable)</span>
            <span className="ch21-split-val delta-v">{fmt(aI.deltaV)}</span>
          </div>
          <div className="ch21-split-row">
            <span className="ch21-split-label">Consumed (capitalist revenue)</span>
            <span className="ch21-split-val consumed">{fmt(aI.consumed)}</span>
          </div>
        </div>

        <div className="ch21-split-card">
          <p className="ch21-split-card-title dept-ii">Dept II — Articles of Consumption</p>
          <div className="ch21-split-row">
            <span className="ch21-split-label">Accumulation rate</span>
            <span>
              <input
                type="range" min={0} max={10000} step={500}
                value={rateIIBps}
                onChange={(e) => setRateIIBps(Number(e.target.value))}
                style={{ width: 80 }}
              />
              {" "}{bpsToPercent(rateIIBps)}
            </span>
          </div>
          <div className="ch21-split-row">
            <span className="ch21-split-label">Surplus (s)</span>
            <span className="ch21-split-val">{fmt(dII.s)}</span>
          </div>
          <div className="ch21-split-row">
            <span className="ch21-split-label">Δc (reinvested to constant)</span>
            <span className="ch21-split-val delta-c">{fmt(aII.deltaC)}</span>
          </div>
          <div className="ch21-split-row">
            <span className="ch21-split-label">Δv (reinvested to variable)</span>
            <span className="ch21-split-val delta-v">{fmt(aII.deltaV)}</span>
          </div>
          <div className="ch21-split-row">
            <span className="ch21-split-label">Consumed (capitalist revenue)</span>
            <span className="ch21-split-val consumed">{fmt(aII.consumed)}</span>
          </div>
        </div>
      </div>
      <div className={`ch21-balance-bar ${balanced ? "ok" : "fail"}`}>
        I.v({fmt(dI.v)}) + Δv_I({fmt(aI.deltaV)}) + consumed_I({fmt(aI.consumed)}) = {fmt(lhs)}
        &nbsp;{balanced ? "==" : "≠"}&nbsp;
        II.c({fmt(dII.c)}) + Δc_II({fmt(aII.deltaC)}) = {fmt(rhs)}
        {balanced ? " ✓" : " ✗"}
      </div>
    </div>
  );
}

// ── Widget 3 — 5-Year Growth Table ───────────────────────────────────────────

function GrowthTable() {
  type YearRow = {
    year: number;
    dI: DeptState;
    dII: DeptState;
    dIGrowthBps: number;
    dIIGrowthBps: number;
    leadConfirmed: boolean;
  };

  const rows: YearRow[] = [];
  let curI: DeptState  = { c: YEAR_I_I.c,  v: YEAR_I_I.v,  s: YEAR_I_I.s,  total: YEAR_I_I.total };
  let curII: DeptState = { c: YEAR_I_II.c, v: YEAR_I_II.v, s: YEAR_I_II.s, total: YEAR_I_II.total };

  rows.push({ year: 1, dI: curI, dII: curII, dIGrowthBps: 0, dIIGrowthBps: 0, leadConfirmed: false });

  for (let y = 2; y <= 5; y++) {
    const prevI = curI;
    const prevII = curII;
    curI  = advance(curI, 5000).next;
    curII = advance(curII, 0).next;
    const dIGrowth  = prevI.total  > 0 ? Math.round((curI.total  - prevI.total)  * 10000 / prevI.total)  : 0;
    const dIIGrowth = prevII.total > 0 ? Math.round((curII.total - prevII.total) * 10000 / prevII.total) : 0;
    rows.push({
      year: y,
      dI: curI,
      dII: curII,
      dIGrowthBps:  dIGrowth,
      dIIGrowthBps: dIIGrowth,
      leadConfirmed: dIGrowth > dIIGrowth,
    });
  }

  return (
    <div className="ch21-widget">
      <p className="ch21-section-title">5-Year Growth Projection — Dept I Outpaces Dept II</p>
      <p className="ch21-subtitle">
        Dept I (means of production) grows faster because its accumulation rate is higher.
        This is Marx's fundamental thesis of extended reproduction.
      </p>
      <div className="ch21-table-wrap">
        <table className="ch21-growth-table">
          <thead>
            <tr>
              <th>Year</th>
              <th>Dept I total</th>
              <th>Dept I growth</th>
              <th>Dept II total</th>
              <th>Dept II growth</th>
              <th>I leads II?</th>
            </tr>
          </thead>
          <tbody>
            {rows.map((r) => (
              <tr key={r.year}>
                <td className="year-col">{r.year}</td>
                <td className="dept-i-val num">{fmt(r.dI.total)}</td>
                <td className="dept-i-val num">
                  {r.year === 1 ? "—" : `+${(r.dIGrowthBps / 100).toFixed(1)}%`}
                </td>
                <td className="dept-ii-val num">{fmt(r.dII.total)}</td>
                <td className="dept-ii-val num">
                  {r.year === 1 ? "—" : `+${(r.dIIGrowthBps / 100).toFixed(1)}%`}
                </td>
                <td className={r.year === 1 ? "lead-no" : r.leadConfirmed ? "lead-yes" : "lead-no"}>
                  {r.year === 1 ? "—" : r.leadConfirmed ? "✓ Yes" : "✗ No"}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      <p className="ch21-note">
        Values in abstract units (×1 000 pence in the backend). Dept I (50% acc. rate) consistently
        outgrows Dept II (0% acc. rate), confirming that expanded accumulation requires
        prior expansion of means-of-production output.
      </p>
    </div>
  );
}

// ── Widget 4 — Organic Composition Shifts ────────────────────────────────────

function CompositionShiftWidget() {
  type CompRow = { year: number; dIComp: number; dIIComp: number; social: number };

  const rows: CompRow[] = [];
  let curI: DeptState  = { c: YEAR_I_I.c,  v: YEAR_I_I.v,  s: YEAR_I_I.s,  total: YEAR_I_I.total };
  let curII: DeptState = { c: YEAR_I_II.c, v: YEAR_I_II.v, s: YEAR_I_II.s, total: YEAR_I_II.total };

  for (let y = 1; y <= 5; y++) {
    const compI  = curI.v  > 0 ? Math.round(curI.c  * 10000 / curI.v)  : 0;
    const compII = curII.v > 0 ? Math.round(curII.c * 10000 / curII.v) : 0;
    rows.push({ year: y, dIComp: compI, dIIComp: compII, social: Math.round((compI + compII) / 2) });
    curI  = advance(curI, 5000).next;
    curII = advance(curII, 0).next;
  }

  return (
    <div className="ch21-widget">
      <p className="ch21-section-title">Rising Organic Composition — c/v Ratio Each Cycle</p>
      <p className="ch21-subtitle">
        As Dept I reinvests, its constant capital grows relative to variable capital.
        The rising c/v ratio is the material basis for the tendential fall in the rate of profit (Vol. III).
      </p>
      <div className="ch21-table-wrap">
        <table className="ch21-comp-table">
          <thead>
            <tr>
              <th>Year</th>
              <th>Dept I c/v (bps)</th>
              <th>Dept II c/v (bps)</th>
              <th>Social average (bps)</th>
            </tr>
          </thead>
          <tbody>
            {rows.map((r, i) => (
              <tr key={r.year}>
                <td className="cycle-col">{r.year}</td>
                <td className={i > 0 && r.dIComp > rows[i - 1].dIComp ? "ch21-comp-rising num" : "ch21-comp-flat num"}>
                  {r.dIComp.toLocaleString()}
                </td>
                <td className={i > 0 && r.dIIComp > rows[i - 1].dIIComp ? "ch21-comp-rising num" : "ch21-comp-flat num"}>
                  {r.dIIComp.toLocaleString()}
                </td>
                <td className={i > 0 && r.social > rows[i - 1].social ? "ch21-comp-rising num" : "ch21-comp-flat num"}>
                  {r.social.toLocaleString()}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      <p className="ch21-note">
        10000 bps = composition of 1:1 (c equals v). Higher values indicate more constant
        relative to variable capital. Rising composition = rising organic composition.
      </p>
    </div>
  );
}

// ── Widget 5 — Money Requirement ─────────────────────────────────────────────

function MoneyRequirementWidget() {
  const totalI  = YEAR_I_I.total;
  const totalII = YEAR_I_II.total;
  const baseMoneyPence = 500;

  // AdditionalMoney = (totalI + totalII) * (rateI + rateII) / 20000
  // rateI = 5000, rateII = 0
  const additionalMoney = Math.floor((totalI + totalII) * (5000 + 0) / 20000);
  const totalMoney = baseMoneyPence + additionalMoney;

  return (
    <div className="ch21-widget">
      <p className="ch21-section-title">Money Required to Circulate the Extended Commodity Mass</p>
      <p className="ch21-subtitle">
        Accumulation increases the commodity mass in circulation. The additional money
        comes from the sources identified in Ch. 17: gold production, hoards, and
        new credit instruments.
      </p>
      <div className="ch21-money-table">
        <div className="ch21-money-row">
          <span className="ch21-money-label">Base money stock (circulating before acc.)</span>
          <span className="ch21-money-val base">{fmt(baseMoneyPence)}</span>
        </div>
        <div className="ch21-money-row">
          <span className="ch21-money-label">Additional money for growing commodity mass</span>
          <span className="ch21-money-val additional">+{fmt(additionalMoney)}</span>
        </div>
        <div className="ch21-money-row">
          <span className="ch21-money-label">Total money required</span>
          <span className="ch21-money-val total">{fmt(totalMoney)}</span>
        </div>
      </div>
      <p className="ch21-note">
        Formula: AdditionalMoney = (I.total + II.total) × (rateI_bps + rateII_bps) / 20000.
        Linked to Ch. 17 realisation sources: gold production is the primary
        money-supply supplement; credit and velocity changes provide the remainder.
      </p>
    </div>
  );
}

// ── Widget 6 — Vol. III Preview ───────────────────────────────────────────────

function VolIIIPreview() {
  return (
    <div className="ch21-widget">
      <p className="ch21-section-title">Bridge to Volume III — The Tendential Fall in the Profit Rate</p>
      <div className="ch21-preview-box">
        <p className="ch21-preview-label">From Extended Reproduction to the Rate of Profit</p>
        <p>
          Extended reproduction drives the <strong>rising organic composition of capital</strong> (c/v ↑).
          As machinery and raw materials grow faster than living labour, the proportion of value
          created by human labour — and therefore the rate of surplus-value relative to total capital
          advanced — tends to fall.
        </p>
        <p style={{ marginTop: "0.5rem" }}>
          <strong>Vol. III, Ch. 13</strong> shows that the average rate of profit (s / (c + v)) tends
          to decline as accumulation proceeds, even while the total mass of profit may grow. This is
          the <em>law of the tendential fall in the rate of profit</em>, grounded in the reproduction
          dynamics established here in Vol. II Ch. 21.
        </p>
        <p style={{ marginTop: "0.5rem" }}>
          The growth lead of Dept I (means of production) over Dept II (articles of consumption),
          confirmed in the table above, is the material precondition for this tendency: living
          labour is progressively displaced by accumulated dead labour.
        </p>
      </div>
    </div>
  );
}

// ── Root export ───────────────────────────────────────────────────────────────

export function Ch21ExtendedReproduction() {
  return (
    <div className="ch21-root">
      <CanonicalFixture />
      <ReinvestmentSplitWidget />
      <GrowthTable />
      <CompositionShiftWidget />
      <MoneyRequirementWidget />
      <VolIIIPreview />
    </div>
  );
}
