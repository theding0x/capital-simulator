import { useEffect, useState } from "react";
import { turnoverTimeApi } from "../../api";
import type { TurnoverTime, ActiveFractionResponse, Perishability } from "../../types";
import "./Ch05TimeOfCirculation.css";

const NANOS_PER_DAY = 86_400_000_000_000;
const NANOS_PER_HOUR = 3_600_000_000_000;

function nanosToLabel(nanos: number): string {
  if (nanos === 0) return "0 s";
  const days = Math.floor(nanos / NANOS_PER_DAY);
  const hours = Math.floor((nanos % NANOS_PER_DAY) / NANOS_PER_HOUR);
  if (days >= 1) return hours > 0 ? `${days}d ${hours}h` : `${days}d`;
  return `${hours}h`;
}

// --- TurnoverBar: stacked-bar of production sub-spans + circulation sub-spans ---
function TurnoverBar({ tt }: { tt: TurnoverTime }) {
  const prod = tt.production;
  const circ = tt.circulation;
  const total =
    prod.labour_time_nanos +
    prod.labour_interruption_nanos +
    prod.latent_nanos +
    prod.natural_process_nanos +
    circ.selling_time_nanos +
    circ.buying_time_nanos;

  if (total === 0) {
    return <p className="tt-empty">No time recorded yet.</p>;
  }

  function pct(n: number) {
    return ((n / total) * 100).toFixed(1) + "%";
  }

  const segments = [
    { label: "Labour time", nanos: prod.labour_time_nanos, cls: "seg-labour" },
    { label: "Interruptions", nanos: prod.labour_interruption_nanos, cls: "seg-interrupt" },
    { label: "Latent MP", nanos: prod.latent_nanos, cls: "seg-latent" },
    { label: "Natural process", nanos: prod.natural_process_nanos, cls: "seg-natural" },
    { label: "Selling (C—M)", nanos: circ.selling_time_nanos, cls: "seg-selling" },
    { label: "Buying (M—C)", nanos: circ.buying_time_nanos, cls: "seg-buying" },
  ].filter((s) => s.nanos > 0);

  return (
    <div className="tt-bar-wrap">
      <div className="tt-bar">
        {segments.map((s) => (
          <div
            key={s.label}
            className={`tt-seg ${s.cls}`}
            style={{ width: pct(s.nanos) }}
            title={`${s.label}: ${nanosToLabel(s.nanos)}`}
          />
        ))}
      </div>
      <div className="tt-legend">
        {segments.map((s) => (
          <span key={s.label} className="tt-legend-item">
            <span className={`tt-swatch ${s.cls}`} />
            {s.label} ({nanosToLabel(s.nanos)})
          </span>
        ))}
      </div>
    </div>
  );
}

// --- Natural process examples (grain / wine / leather) ---
const NATURAL_EXAMPLES = [
  {
    kind: "Grain ripening",
    durationDays: 90,
    quote: "grain sown in the fields",
  },
  {
    kind: "Wine fermentation",
    durationDays: 14,
    quote: "wine fermenting in the cellar",
  },
  {
    kind: "Leather tanning",
    durationDays: 30,
    quote: "tanneries exposed to chemical processes",
  },
];

// --- Perishability cliff widget ---
function PerishabilityCliff({ perishabilities }: { perishabilities: Perishability[] }) {
  const [elapsed, setElapsed] = useState(0);
  const STEP = NANOS_PER_HOUR * 6;

  function advance() {
    setElapsed((e) => e + STEP);
  }

  function reset() {
    setElapsed(0);
  }

  return (
    <div className="perishability-cliff">
      <div className="cliff-controls">
        <button onClick={advance} className="btn-secondary">+6 h</button>
        <button onClick={reset} className="btn-secondary">Reset</button>
        <span className="cliff-elapsed">Elapsed: {nanosToLabel(elapsed)}</span>
      </div>
      <div className="cliff-rows">
        {perishabilities.map((p) => {
          const spoiled = p.window_nanos > 0 && elapsed > p.window_nanos;
          const pct = p.window_nanos > 0 ? Math.min((elapsed / p.window_nanos) * 100, 100) : 0;
          return (
            <div key={p.id} className={`cliff-row ${spoiled ? "spoiled" : ""}`}>
              <span className="cliff-commodity">{p.commodity_id}</span>
              <span className="cliff-window">{nanosToLabel(p.window_nanos)} window</span>
              <div className="cliff-bar-wrap">
                <div className="cliff-bar-fill" style={{ width: `${pct.toFixed(1)}%` }} />
              </div>
              <span className="cliff-status">{spoiled ? "SPOILED" : "ok"}</span>
            </div>
          );
        })}
        {perishabilities.length === 0 && (
          <p className="cliff-empty muted small">No perishability records loaded yet.</p>
        )}
      </div>
    </div>
  );
}

// --- Zero-circulation toggle ---
function ZeroCirculationToggle({ tt }: { tt: TurnoverTime }) {
  const [zeroed, setZeroed] = useState(false);

  const prodNanos =
    tt.production.labour_time_nanos +
    tt.production.labour_interruption_nanos +
    tt.production.latent_nanos +
    tt.production.natural_process_nanos;

  const circNanos = zeroed
    ? 0
    : tt.circulation.selling_time_nanos + tt.circulation.buying_time_nanos;

  const total = prodNanos + circNanos;
  const basisPts = total > 0 ? Math.round((10000 * prodNanos) / total) : 0;

  return (
    <div className="zero-circ">
      <label className="zero-circ-toggle">
        <input
          type="checkbox"
          checked={zeroed}
          onChange={(e) => setZeroed(e.target.checked)}
        />
        Collapse circulation time to zero (payment-on-delivery ideal)
      </label>
      <div className="zero-circ-result">
        <span className="zero-circ-label">Active fraction:</span>
        <span className="zero-circ-value">{(basisPts / 100).toFixed(1)}%</span>
        <span className="zero-circ-bp muted small"> ({basisPts} bp)</span>
      </div>
      <p className="muted small">
        §: "If this payment is made in his own means of production, the time of circulation approaches zero."
      </p>
    </div>
  );
}

// --- Main panel ---
export function Ch05TimeOfCirculation() {
  const [turnoverTimes, setTurnoverTimes] = useState<TurnoverTime[]>([]);
  const [selectedId, setSelectedId] = useState("");
  const [activeFraction, setActiveFraction] = useState<ActiveFractionResponse | null>(null);
  const [perishabilities, setPerishabilities] = useState<Perishability[]>([]);
  const [listError, setListError] = useState("");

  useEffect(() => {
    refresh();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  async function refresh() {
    setListError("");
    try {
      const list = await turnoverTimeApi.list();
      const tts = Array.isArray(list) ? list : [];
      setTurnoverTimes(tts);
      if (tts.length > 0 && !selectedId) {
        setSelectedId(tts[0].id);
        loadActiveFraction(tts[0].id);
      }
      // Seed perishability display rows (milk 2d, beer 7d).
      const milkP = await turnoverTimeApi.setPerishability("milk", 2 * NANOS_PER_DAY).catch(() => null);
      const beerP = await turnoverTimeApi.setPerishability("beer", 7 * NANOS_PER_DAY).catch(() => null);
      setPerishabilities([milkP, beerP].filter(Boolean) as Perishability[]);
    } catch (e) {
      setListError(String(e));
    }
  }

  function loadActiveFraction(id: string) {
    turnoverTimeApi.getActiveFraction(id)
      .then(setActiveFraction)
      .catch(() => setActiveFraction(null));
  }

  function handleSelect(id: string) {
    setSelectedId(id);
    loadActiveFraction(id);
  }

  const selected = turnoverTimes.find((t) => t.id === selectedId) ?? null;

  return (
    <div className="ch05-panel">
      {/* ── TurnoverTime list ── */}
      <section className="ch05-section">
        <h2 className="ch05-h2">Turnover Records</h2>
        {listError && <p className="error-msg">{listError}</p>}
        {turnoverTimes.length === 0 && !listError && (
          <p className="muted small">No turnover records. Seed the database to populate SpinningMill1871.</p>
        )}
        <div className="tt-list">
          {turnoverTimes.map((tt) => {
            const ttTotal =
              tt.production.labour_time_nanos +
              tt.production.labour_interruption_nanos +
              tt.production.latent_nanos +
              tt.production.natural_process_nanos +
              tt.circulation.selling_time_nanos +
              tt.circulation.buying_time_nanos;
            return (
              <button
                key={tt.id}
                className={`tt-list-item ${tt.id === selectedId ? "selected" : ""}`}
                onClick={() => handleSelect(tt.id)}
              >
                <span className="tt-id">{tt.id.slice(0, 10)}…</span>
                <span className="tt-total">{nanosToLabel(ttTotal)}</span>
              </button>
            );
          })}
        </div>
      </section>

      {/* ── Stacked bar for selected TurnoverTime ── */}
      {selected && (
        <section className="ch05-section">
          <h2 className="ch05-h2">Time Decomposition</h2>
          <TurnoverBar tt={selected} />

          {activeFraction && (
            <div className="active-fraction-badge">
              Active fraction:{" "}
              <strong>{(activeFraction.active_fraction_basis_points / 100).toFixed(1)}%</strong>
              <span className="muted small"> ({activeFraction.active_fraction_basis_points} bp)</span>
            </div>
          )}

          <ZeroCirculationToggle tt={selected} />
        </section>
      )}

      {/* ── Natural process examples ── */}
      <section className="ch05-section">
        <h2 className="ch05-h2">Natural-Process Spans</h2>
        <p className="muted small">
          §: "Grain sown, wine fermenting in the cellar, tanneries exposed to chemical
          processes — labour is absent; no surplus-value is created."
        </p>
        <div className="np-grid">
          {NATURAL_EXAMPLES.map((ex) => (
            <div key={ex.kind} className="np-card">
              <div className="np-kind">{ex.kind}</div>
              <div className="np-duration">{ex.durationDays}d</div>
              <div className="np-quote muted small">§ {ex.quote}</div>
            </div>
          ))}
        </div>
      </section>

      {/* ── Perishability cliff ── */}
      <section className="ch05-section">
        <h2 className="ch05-h2">Perishability Cliff</h2>
        <p className="muted small">
          §: "The more perishable a commodity is … the less is it suited to be an object of
          capitalist production." Advance the clock to see the spoilage threshold.
        </p>
        <PerishabilityCliff perishabilities={perishabilities} />
      </section>
    </div>
  );
}
