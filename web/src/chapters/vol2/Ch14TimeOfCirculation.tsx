import { useEffect, useState } from "react";
import type {
  CirculationSpeedImprovement,
  DistanceLagRelation,
} from "../../types";
import { circulationTimeV2Api } from "../../api";
import "./Ch14TimeOfCirculation.css";

// Vol. II Ch. 14 — The Time of Circulation (refined).
//
// Key result: selling-time and buying-time each decompose into four sub-spans.
// The 1870 England–India telegraph collapsed a ~250-day communication lag to
// 1 day — liberating capital previously tied up in the time of circulation.
//
// Seed fixtures (00008_v2_ch14_seed.sql):
//   Pre-telegraph (1850):  365-day selling phase, comm-lag 250d, penalty 5000 bp
//   Post-telegraph (1870):  30-day selling phase, comm-lag 1d
//   DistanceLagRelations:  London-Manchester 320km, London-Bombay 9000km x2
//   SpeedImprovements:     railway 1840, steamship 1845, telegraph 1870, telephone 1876

function fmtMin(m: number): string {
  const days = m / 1440;
  if (days >= 365) return `${(days / 365).toFixed(1)} yr`;
  if (days >= 1) return `${Math.round(days)} d`;
  return `${m} min`;
}

// ---- EnglandIndiaSlider ----------------------------------------------------

const PRE = { comm: 360000, search: 72000, bargain: 72000, legal: 21600, total: 525600 };
const POST = { comm: 1440, search: 14400, bargain: 14400, legal: 12960, total: 43200 };

function lerp(a: number, b: number, t: number) {
  return Math.round(a + (b - a) * t);
}

function EnglandIndiaSlider() {
  const [t, setT] = useState(0);
  const comm = lerp(PRE.comm, POST.comm, t);
  const search = lerp(PRE.search, POST.search, t);
  const bargain = lerp(PRE.bargain, POST.bargain, t);
  const legal = lerp(PRE.legal, POST.legal, t);
  const total = lerp(PRE.total, POST.total, t);
  const pct = (v: number) => `${((v / total) * 100).toFixed(1)}%`;

  return (
    <div className="ch14-slider-wrapper">
      <div className="ch14-era-label">
        <span>Pre-telegraph (1850) — 365-day phase</span>
        <span>Post-telegraph (1870) — 30-day phase</span>
      </div>
      <input type="range" className="ch14-slider" min={0} max={100}
        value={Math.round(t * 100)} onChange={(e) => setT(Number(e.target.value) / 100)} />
      <div className="ch14-bars">
        {[
          { label: "Comm. lag", value: comm, cls: "ch14-bar-comm" },
          { label: "Buyer search", value: search, cls: "ch14-bar-search" },
          { label: "Bargaining", value: bargain, cls: "ch14-bar-bargain" },
          { label: "Legal settle.", value: legal, cls: "ch14-bar-legal" },
        ].map(({ label, value, cls }) => (
          <div className="ch14-bar-group" key={label}>
            <label>{label}</label>
            <div className="ch14-bar-outer">
              <div className={`ch14-bar-inner ${cls}`} style={{ width: pct(value) }} />
            </div>
            <div className="ch14-bar-val">{fmtMin(value)}</div>
          </div>
        ))}
      </div>
      <div style={{ fontSize: "0.75rem", color: "#64748b" }}>
        Total selling time: <strong>{fmtMin(total)}</strong>
      </div>
    </div>
  );
}

// ---- PhaseBreakdownWidget --------------------------------------------------

function PhaseBreakdownWidget() {
  const segs = [
    { label: "Search 5d", minutes: 7200, color: "#f59e0b" },
    { label: "Bargain 5d", minutes: 7200, color: "#3b82f6" },
    { label: "Comm-lag 15d", minutes: 21600, color: "#ef4444" },
    { label: "Legal 5d", minutes: 7200, color: "#10b981" },
  ];
  const total = segs.reduce((s, x) => s + x.minutes, 0);
  return (
    <div>
      <div style={{ fontSize: "0.8rem", color: "#475569", marginBottom: "0.5rem" }}>
        Typical 30-day selling phase sub-span breakdown
      </div>
      <div className="ch14-breakdown">
        {segs.map(({ label, minutes, color }) => (
          <div key={label} className="ch14-seg"
            style={{ flex: minutes, background: color, padding: "0 4px" }}
            title={`${label} — ${fmtMin(minutes)}`}>
            {minutes / total > 0.15 ? label : ""}
          </div>
        ))}
      </div>
      <div style={{ display: "flex", gap: "0.75rem", marginTop: "0.5rem", flexWrap: "wrap" }}>
        {segs.map(({ label, color, minutes }) => (
          <div key={label} style={{ display: "flex", alignItems: "center", gap: "4px", fontSize: "0.72rem" }}>
            <span style={{ width: 10, height: 10, borderRadius: 2, background: color, display: "inline-block" }} />
            {label} ({fmtMin(minutes)})
          </div>
        ))}
      </div>
    </div>
  );
}

// ---- SpeedTimeline ---------------------------------------------------------

const SEED_IMPROVEMENTS: CirculationSpeedImprovement[] = [
  { id: "s1", market_id: "london-manchester", old_lag_minutes_per_km: 10, new_lag_minutes_per_km: 3, reason: "railway", effective_at: "1840-01-01T00:00:00Z", created_at: "" },
  { id: "s2", market_id: "london-bombay", old_lag_minutes_per_km: 90, new_lag_minutes_per_km: 58, reason: "steamship", effective_at: "1845-01-01T00:00:00Z", created_at: "" },
  { id: "s3", market_id: "london-bombay", old_lag_minutes_per_km: 58, new_lag_minutes_per_km: 1, reason: "telegraph", effective_at: "1870-01-01T00:00:00Z", created_at: "" },
  { id: "s4", market_id: "london-manchester", old_lag_minutes_per_km: 3, new_lag_minutes_per_km: 0, reason: "telephone", effective_at: "1876-01-01T00:00:00Z", created_at: "" },
];

const REASON_LABEL: Record<string, string> = {
  railway: "Railway", steamship: "Steamship", telegraph: "Telegraph", telephone: "Telephone", other: "Other",
};

function SpeedTimeline({ items }: { items: CirculationSpeedImprovement[] }) {
  const sorted = [...items].sort(
    (a, b) => new Date(a.effective_at).getTime() - new Date(b.effective_at).getTime()
  );
  return (
    <div className="ch14-timeline">
      {sorted.map((csi) => (
        <div key={csi.id} className="ch14-event">
          <div className="ch14-event-year">
            {new Date(csi.effective_at).getFullYear()} — {REASON_LABEL[csi.reason] ?? csi.reason}
          </div>
          <div className="ch14-event-desc">{csi.market_id.replace("-", " → ")}</div>
          <div className="ch14-event-lag">
            {csi.old_lag_minutes_per_km} min/km → {csi.new_lag_minutes_per_km} min/km
            {" "}(−{csi.old_lag_minutes_per_km - csi.new_lag_minutes_per_km} min/km)
          </div>
        </div>
      ))}
    </div>
  );
}

// ---- PenaltyCalculator -----------------------------------------------------

function PenaltyCalculator() {
  const [ctPct, setCtPct] = useState(50);
  const actual = Math.round(10000 * (1 - ctPct / 100));
  const penalty = Math.max(0, 10000 - actual);
  return (
    <div className="ch14-calc">
      <div className="ch14-calc-input">
        <label>Circulation time as % of annual period: {ctPct}%</label>
        <input type="range" min={0} max={80} value={ctPct}
          onChange={(e) => setCtPct(Number(e.target.value))} />
        <div style={{ fontSize: "0.72rem", color: "#475569" }}>
          Ideal rate: <strong>10 000 bp</strong> (zero CT)<br />
          Actual rate: <strong>{actual} bp</strong>
        </div>
      </div>
      <div className="ch14-penalty-display">
        <div className="ch14-penalty-bp">{penalty.toLocaleString()}</div>
        <div className="ch14-penalty-label">basis-point penalty vs zero-CT ideal</div>
        <div style={{ fontSize: "0.7rem", color: "#64748b", marginTop: "0.5rem" }}>
          Full annual-rate calculation deferred to Ch. 16
        </div>
      </div>
    </div>
  );
}

// ---- DistanceLagTable ------------------------------------------------------

function DistanceLagTable({ rows }: { rows: DistanceLagRelation[] }) {
  if (rows.length === 0) {
    return <div style={{ fontSize: "0.8rem", color: "#64748b" }}>No distance-lag relations recorded yet.</div>;
  }
  return (
    <table className="ch14-table">
      <thead>
        <tr><th>Market</th><th>Distance (km)</th><th>Lag (min/km)</th><th>Max lag</th></tr>
      </thead>
      <tbody>
        {rows.map((r) => (
          <tr key={r.id}>
            <td>{r.market_id}</td>
            <td>{(r.distance_meters / 1000).toLocaleString()}</td>
            <td>{r.lag_minutes_per_km}</td>
            <td>{fmtMin((r.lag_minutes_per_km * r.distance_meters) / 1000)}</td>
          </tr>
        ))}
      </tbody>
    </table>
  );
}

// ---- Main panel ------------------------------------------------------------

export default function Ch14TimeOfCirculation() {
  const [lagRelations, setLagRelations] = useState<DistanceLagRelation[]>([]);
  const [improvements, setImprovements] = useState<CirculationSpeedImprovement[]>(SEED_IMPROVEMENTS);

  useEffect(() => {
    circulationTimeV2Api.listDistanceLagRelations()
      .then((r) => { if (r.items?.length) setLagRelations(r.items); })
      .catch(() => {});
    circulationTimeV2Api.listSpeedImprovements()
      .then((r) => { if (r.items?.length) setImprovements(r.items); })
      .catch(() => {});
  }, []);

  return (
    <div className="ch14-root">
      <div className="ch14-section">
        <h3>§1 — England–India: Circulation-time collapse (telegraph slider)</h3>
        <p style={{ fontSize: "0.8rem", color: "#475569", marginBottom: "0.75rem" }}>
          Before 1870 the London–Bombay selling phase required ~365 days; 250 of those days were
          pure communication lag (ship post). The 1870 telegraph reduced that lag to 1 day.{" "}
          <em>§: "the communication required from 12 to 18 months — now a few hours."</em>
        </p>
        <EnglandIndiaSlider />
      </div>

      <div className="ch14-section">
        <h3>§2 — 30-day SellingPhase sub-span breakdown</h3>
        <PhaseBreakdownWidget />
      </div>

      <div className="ch14-section">
        <h3>§3 — Speed-improvement timeline</h3>
        <p style={{ fontSize: "0.8rem", color: "#475569", marginBottom: "0.75rem" }}>
          Each event reduces LagMinutesPerKm, shortening the communication sub-span and releasing tied-up capital.
        </p>
        <SpeedTimeline items={improvements} />
      </div>

      <div className="ch14-section">
        <h3>§4 — AnnualSurplusPenalty calculator</h3>
        <p style={{ fontSize: "0.8rem", color: "#475569", marginBottom: "0.75rem" }}>
          Drag to see how increasing circulation time penalises the annual rate of surplus-value
          relative to the zero-CT ideal.
        </p>
        <PenaltyCalculator />
      </div>

      <div className="ch14-section">
        <h3>Distance-lag relations</h3>
        <DistanceLagTable rows={lagRelations} />
      </div>
    </div>
  );
}
