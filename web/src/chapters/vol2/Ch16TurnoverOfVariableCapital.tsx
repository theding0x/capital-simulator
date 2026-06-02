import { useEffect, useState } from "react";
import type {
  AnnualSurplusRate,
  VariableCapitalReproduction,
} from "../../types";
import { annualSurplusRateApi } from "../../api";
import "./Ch16TurnoverOfVariableCapital.css";

// Vol. II Ch. 16 — The Turnover of Variable Capital.
//
// The annual rate of surplus-value depends not only on the per-circuit rate
// s/v but on how many times that circuit repeats in a year (n):
//
//   annual_S = n × (s/v)
//
// Seed fixtures (00051_v2_ch16_seed.sql):
//   A  — 5eed0000000000001601: £500, n=10, s/v=100% → annual=1000%
//   B  — 5eed0000000000001602: £5000, n=1, s/v=100% → annual=100%
//   Sm — 5eed0000000000001603: Spinning Mill 1871, n=52, s/v=156% → 8112%
//   Lw — 5eed0000000000001605: Locomotive Works, n≈4.3, s/v=80% → ≈348%

// ── helpers ─────────────────────────────────────────────────────────────────

function bpToPercent(bp: number): string {
  return (bp / 100).toFixed(bp % 100 === 0 ? 0 : 2) + "%";
}

function penceToPounds(p: number): string {
  const pounds = p / 240;
  return "£" + (pounds % 1 === 0 ? pounds.toFixed(0) : pounds.toFixed(2));
}

// ── Widget 1 — Capitalist A vs B ────────────────────────────────────────────

const SEED_SM = "5eed0000000000001603";
const SEED_LW = "5eed0000000000001605";

interface CapitalistData {
  label: string;
  advancePounds: string;
  turnoverWeeks: string;
  n: string;
  pcRate: string;
  annualRate: string;
  annualSurplusPounds: string;
}

const STATIC_A: CapitalistData = {
  label: "Capitalist A",
  advancePounds: "£500",
  turnoverWeeks: "5 weeks",
  n: "10",
  pcRate: "100%",
  annualRate: "1,000%",
  annualSurplusPounds: "£5,000",
};

const STATIC_B: CapitalistData = {
  label: "Capitalist B",
  advancePounds: "£5,000",
  turnoverWeeks: "50 weeks",
  n: "1",
  pcRate: "100%",
  annualRate: "100%",
  annualSurplusPounds: "£5,000",
};

function CapitalistCard({ data }: { data: CapitalistData }) {
  return (
    <div className="ch16-capital-card">
      <div className="ch16-capital-card-header">
        <span className="ch16-capital-label">{data.label}</span>
        <span className="ch16-annual-rate">{data.annualRate}</span>
      </div>
      <div className="ch16-detail-row">
        <span>Variable capital advanced</span>
        <span>{data.advancePounds}</span>
      </div>
      <div className="ch16-detail-row">
        <span>Turnover period</span>
        <span>{data.turnoverWeeks}</span>
      </div>
      <div className="ch16-detail-row">
        <span>Turnovers per year (n)</span>
        <span>{data.n}</span>
      </div>
      <div className="ch16-detail-row">
        <span>Per-circuit s/v</span>
        <span>{data.pcRate}</span>
      </div>
      <div className="ch16-detail-row">
        <span>Annual surplus produced</span>
        <span>{data.annualSurplusPounds}</span>
      </div>
    </div>
  );
}

// ── Widget 2 — Live formula slider ──────────────────────────────────────────

function FormulaSlider() {
  const [n, setN] = useState(10);
  const [pcRate, setPcRate] = useState(100);
  const annual = n * pcRate;

  return (
    <div className="ch16-formula-box">
      <p className="ch16-section-title">Live formula: Annual S = n × (s/v)</p>
      <div className="ch16-formula-display">
        {n} × {pcRate}% = <strong>{annual.toLocaleString()}%</strong>
      </div>
      <div className="ch16-slider-grid">
        <div className="ch16-slider-row">
          <label className="ch16-slider-label">
            Number of turnovers (n):{" "}
            <span className="ch16-slider-value">{n}</span>
          </label>
          <input
            type="range"
            min={1}
            max={52}
            value={n}
            onChange={(e) => setN(Number(e.target.value))}
          />
        </div>
        <div className="ch16-slider-row">
          <label className="ch16-slider-label">
            Per-circuit rate s/v:{" "}
            <span className="ch16-slider-value">{pcRate}%</span>
          </label>
          <input
            type="range"
            min={1}
            max={200}
            value={pcRate}
            onChange={(e) => setPcRate(Number(e.target.value))}
          />
        </div>
      </div>
      <p className="ch16-quadruple-note">
        Doubling both n and s/v quadruples the annual rate (
        {n * 2} × {pcRate * 2}% = {(n * 2 * pcRate * 2).toLocaleString()}%).
      </p>
    </div>
  );
}

// ── Widget 3 — Industry comparison (LocoWorks vs SpinMill) ──────────────────

interface IndustryData {
  name: string;
  year: string;
  nApprox: string;
  pcRate: string;
  annualRateFallback: string;
}

const LOCO_DATA: IndustryData = {
  name: "Locomotive Works",
  year: "1871",
  nApprox: "≈4.3",
  pcRate: "80%",
  annualRateFallback: "≈348%",
};

const SPIN_DATA: IndustryData = {
  name: "Spinning Mill",
  year: "1871",
  nApprox: "52",
  pcRate: "156%",
  annualRateFallback: "8,112%",
};

function IndustryCard({
  data,
  seedId,
  allRates,
}: {
  data: IndustryData;
  seedId: string;
  allRates: AnnualSurplusRate[];
}) {
  const found = allRates.find((r) => r.variable_capital_advance_id === seedId);
  const displayRate = found
    ? bpToPercent(found.annual_basis_points)
    : data.annualRateFallback;

  return (
    <div className="ch16-industry-card">
      <div className="ch16-industry-name">{data.name}</div>
      <div className="ch16-industry-rate">{displayRate}</div>
      <div className="ch16-detail-row">
        <span>Turnovers per year</span>
        <span>{data.nApprox}</span>
      </div>
      <div className="ch16-detail-row">
        <span>Per-circuit s/v</span>
        <span>{data.pcRate}</span>
      </div>
    </div>
  );
}

// ── Widget 4 — Interruption effect ──────────────────────────────────────────

function InterruptionTable() {
  // Capitalist A: base n=10, s/v=100%; show effect of interruptions reducing n.
  const baseN = 10;
  const pcRate = 100;
  const scenarios = [
    { n: baseN, label: "No interruption (base)" },
    { n: 8, label: "2-week interruption" },
    { n: 6, label: "4-week interruption" },
    { n: 4, label: "6-week interruption" },
  ];

  return (
    <table className="ch16-interruption-table">
      <thead>
        <tr>
          <th>Scenario</th>
          <th>n</th>
          <th>Annual rate</th>
          <th>Lost surplus</th>
        </tr>
      </thead>
      <tbody>
        {scenarios.map(({ n, label }) => {
          const annual = n * pcRate;
          const lost = (baseN - n) * pcRate;
          return (
            <tr key={n} className={n === baseN ? "ch16-base" : undefined}>
              <td>{label}</td>
              <td>{n}</td>
              <td>{annual}%</td>
              <td>{lost > 0 ? `−${lost}%` : "—"}</td>
            </tr>
          );
        })}
      </tbody>
    </table>
  );
}

// ── Widget 5 — Spinning Mill 1871 ────────────────────────────────────────────

function SpinMillCard({ repro }: { repro: VariableCapitalReproduction | null }) {
  return (
    <div className="ch16-spinmill-card">
      <p className="ch16-section-title">Spinning Mill 1871 — full calculation</p>
      <div className="ch16-spinmill-calc">
        <span>n = <span className="ch16-val">52</span></span>
        <span>×</span>
        <span>s/v = <span className="ch16-val">156%</span></span>
        <span>=</span>
        <span className="ch16-spinmill-result">8,112%</span>
      </div>
      {repro && (
        <>
          <div className="ch16-detail-row">
            <span>Total wages paid in year</span>
            <span>{penceToPounds(repro.total_variable_paid_annual_pence)}</span>
          </div>
          <div className="ch16-detail-row">
            <span>Total surplus produced in year</span>
            <span>{penceToPounds(repro.total_surplus_annual_pence)}</span>
          </div>
        </>
      )}
    </div>
  );
}

// ── Root component ────────────────────────────────────────────────────────────

// Widget 6 — Price-adjusted annual rate (Ch. 15 → Ch. 16 bridge, issue #222).
// Client-side mirror of valorisation.ComputeAdjustedAnnualRate.
function PriceAdjustedRate() {
  const [v, setV] = useState(1000);
  const [s, setS] = useState(1000);
  const [n, setN] = useState(1);
  const [pricePct, setPricePct] = useState(0);

  const priceBP = Math.round(pricePct * 100);
  const nBP = Math.round(n * 10000);
  const valueAdded = v + s;
  const adjustedSurplus = s + Math.trunc((valueAdded * priceBP) / 10000);
  const basePerCircuit = v > 0 ? Math.trunc((10000 * s) / v) : 0;
  const adjPerCircuit = v > 0 ? Math.trunc((10000 * adjustedSurplus) / v) : 0;
  const baseAnnual = Math.trunc((nBP * basePerCircuit) / 10000);
  const adjAnnual = Math.trunc((nBP * adjPerCircuit) / 10000);

  return (
    <div className="ch16-price-adjusted">
      <p>
        A Ch. 15 price revolution on output re-prices the value added (v + s);
        because the advance is fixed, the swing accrues to surplus against the
        original advance, lifting or cutting the annual rate of surplus-value.
      </p>
      <div className="ch16-adj-inputs">
        <label>
          v (pence)
          <input type="number" min={1} value={v}
            onChange={(e) => setV(Math.max(1, Math.floor(Number(e.target.value))))} />
        </label>
        <label>
          s (pence)
          <input type="number" min={0} value={s}
            onChange={(e) => setS(Math.max(0, Math.floor(Number(e.target.value))))} />
        </label>
        <label>
          turnovers / yr
          <input type="number" min={1} value={n}
            onChange={(e) => setN(Math.max(1, Math.floor(Number(e.target.value))))} />
        </label>
        <label>
          output price Δ %
          <input type="number" value={pricePct}
            onChange={(e) => setPricePct(Number(e.target.value))} />
        </label>
      </div>
      <div className="ch16-adj-result">
        <span>Base annual rate: <strong>{bpToPercent(baseAnnual)}</strong></span>
        <span>Price-adjusted annual rate: <strong>{bpToPercent(adjAnnual)}</strong></span>
      </div>
    </div>
  );
}

export function Ch16TurnoverOfVariableCapital() {
  const [annualRates1871, setAnnualRates1871] = useState<AnnualSurplusRate[]>([]);
  const [spinMillRepro, setSpinMillRepro] =
    useState<VariableCapitalReproduction | null>(null);

  useEffect(() => {
    annualSurplusRateApi
      .listAnnualRates(undefined, "1871")
      .then((res) => setAnnualRates1871(res.items))
      .catch(() => undefined);

    annualSurplusRateApi
      .getReproduction(SEED_SM)
      .then(setSpinMillRepro)
      .catch(() => undefined);
  }, []);

  return (
    <div className="ch16-root">
      <blockquote className="ch16-quote">
        "If we consider the annual rate of surplus-value, it is evident that it
        depends not only on the per-circuit rate of surplus-value but also on
        the number of times the capital undergoes this circuit in the course of
        a year. If s/v = 100% and n = 10, then annual S/V = 1,000%."
        <cite>— Capital, Vol. II, Ch. 16</cite>
      </blockquote>

      {/* Widget 1 — A vs B */}
      <section>
        <p className="ch16-section-title">
          §2 — Capitalist A vs B: same annual surplus, different annual rates
        </p>
        <div className="ch16-ab-grid">
          <CapitalistCard data={STATIC_A} />
          <CapitalistCard data={STATIC_B} />
        </div>
      </section>

      {/* Widget 2 — Formula slider */}
      <section>
        <FormulaSlider />
      </section>

      {/* Widget 3 — Industry comparison */}
      <section>
        <p className="ch16-section-title">
          Industry comparison — Locomotive Works vs Spinning Mill (1871)
        </p>
        <div className="ch16-industry-grid">
          <IndustryCard
            data={LOCO_DATA}
            seedId={SEED_LW}
            allRates={annualRates1871}
          />
          <IndustryCard
            data={SPIN_DATA}
            seedId={SEED_SM}
            allRates={annualRates1871}
          />
        </div>
      </section>

      {/* Widget 4 — Interruption effect */}
      <section>
        <p className="ch16-section-title">
          Effect of interruptions on Capitalist A (base: n=10, s/v=100%)
        </p>
        <InterruptionTable />
      </section>

      {/* Widget 5 — SpinMill 1871 */}
      <section>
        <SpinMillCard repro={spinMillRepro} />
      </section>

      {/* Widget 6 — Price-adjusted annual rate (issue #222) */}
      <section>
        <p className="ch16-section-title">
          Price-adjusted annual rate — Ch. 15 value revolution → Ch. 16
        </p>
        <PriceAdjustedRate />
      </section>
    </div>
  );
}

export default Ch16TurnoverOfVariableCapital;
