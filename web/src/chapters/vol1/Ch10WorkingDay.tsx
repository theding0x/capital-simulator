import { useState } from "react";
import type { FormEvent } from "react";
import { api } from "../../api";
import type {
  WorkingDayResponse,
  ValidateWorkingDayResponse,
  RelaySchedule,
} from "../../types";
import "./Ch10WorkingDay.css";

interface Ch10Props {
  onSharedChanged: () => void;
}

const PHYSICAL_MAX_MIN = 24 * 60;

function minutesToHours(m: number): string {
  const h = Math.floor(m / 60);
  const min = m % 60;
  return min === 0 ? `${h}h` : `${h}h ${min}m`;
}

export function Ch10WorkingDay({ onSharedChanged: _onSharedChanged }: Ch10Props) {
  return (
    <>
      <WorkingDayInsight />
      <ValidatePanel />
      <CreatePanel />
      <RelaySchedulePanel />
      <Coda />
    </>
  );
}

function WorkingDayInsight() {
  return (
    <section className="v1-ch10-insight">
      <h2 className="v1-ch10-insight-h2">A — B — C</h2>
      <p className="v1-ch10-insight-prose">
        The working day is bounded only physically (24 h) and morally. Between
        A and B the worker reproduces their own labour-power; between B and C
        they produce surplus for capital. The whole class struggle of the
        chapter is over where, between B and the physical maximum, the
        statutory cap is drawn.
      </p>
    </section>
  );
}

function AbcBar({
  necessary,
  surplus,
  statutoryLimit,
}: {
  necessary: number;
  surplus: number;
  statutoryLimit?: number | null;
}) {
  const pct = (n: number) => `${(n / PHYSICAL_MAX_MIN) * 100}%`;
  return (
    <div className="v1-ch10-abc">
      <div className="v1-ch10-abc-ticks" aria-hidden="true">
        {Array.from({ length: 24 }).map((_, i) => (
          <span key={i} className="v1-ch10-abc-tick">
            {i % 3 === 0 ? i : ""}
          </span>
        ))}
      </div>
      <div
        className="v1-ch10-abc-track"
        role="img"
        aria-label={`A–B–C working day: necessary ${minutesToHours(necessary)}, surplus ${minutesToHours(surplus)}`}
      >
        <div
          className="v1-ch10-abc-seg v1-ch10-abc-seg--necessary"
          style={{ width: pct(necessary) }}
          title={`A → B (necessary): ${minutesToHours(necessary)}`}
        >
          A → B
        </div>
        <div
          className="v1-ch10-abc-seg v1-ch10-abc-seg--surplus"
          style={{ left: pct(necessary), width: pct(surplus) }}
          title={`B → C (surplus): ${minutesToHours(surplus)}`}
        >
          B → C
        </div>
        {statutoryLimit != null && statutoryLimit > 0 && statutoryLimit <= PHYSICAL_MAX_MIN && (
          <div className="v1-ch10-abc-cap" style={{ left: pct(statutoryLimit) }}>
            <span className="v1-ch10-abc-cap-label">
              Statutory cap · {minutesToHours(statutoryLimit)}
            </span>
          </div>
        )}
      </div>
      <p className="v1-ch10-legend">
        <span>
          <span className="v1-ch10-legend-swatch" style={{ background: "var(--lead)" }} />
          A → B necessary ({minutesToHours(necessary)})
        </span>
        <span>
          <span className="v1-ch10-legend-swatch" style={{ background: "var(--gold-bright)" }} />
          B → C surplus ({minutesToHours(surplus)})
        </span>
        <span>
          <span className="v1-ch10-legend-swatch" style={{ background: "var(--red)" }} />
          Statutory cap
        </span>
      </p>
    </div>
  );
}

function RateChip({ v, s, valid }: { v: number; s: number; valid?: boolean }) {
  const rate = v > 0 ? ((s / v) * 100).toFixed(1) : "∞";
  return (
    <div className={"v1-ch10-rate" + (valid === false ? " v1-ch10-rate--invalid" : "")}>
      <span className="v1-ch10-rate-label">s / v</span>
      <span className="v1-ch10-rate-value">{rate}%</span>
    </div>
  );
}

interface SeriesPoint {
  id: string;
  label: string;
  necessary: number;
  surplus: number;
}

const SERIES_I: SeriesPoint[] = [
  { id: "wd1", label: "WD I", necessary: 360, surplus: 60 },
  { id: "wd2", label: "WD II", necessary: 360, surplus: 180 },
  { id: "wd3", label: "WD III", necessary: 360, surplus: 360 },
];

function RateSeriesChart({ visited }: { visited: Set<string> }) {
  const points = SERIES_I.filter((p) => visited.has(p.id));
  if (points.length === 0) return null;
  const maxLen = Math.max(...points.map((p) => p.necessary + p.surplus));
  return (
    <div className="v1-ch10-series">
      <h3 className="v1-ch10-series-h3">Rate vs working-day length</h3>
      {points.map((p) => {
        const len = p.necessary + p.surplus;
        const rate = (p.surplus / p.necessary) * 100;
        const nPct = `${(p.necessary / maxLen) * 100}%`;
        const sPct = `${(p.surplus / maxLen) * 100}%`;
        return (
          <div key={p.id} className="v1-ch10-series-row">
            <span className="v1-ch10-series-label">{p.label} · {minutesToHours(len)}</span>
            <div className="v1-ch10-series-track">
              <div
                className="v1-ch10-abc-seg v1-ch10-abc-seg--necessary"
                style={{ position: "static", width: nPct }}
              />
              <div
                className="v1-ch10-abc-seg v1-ch10-abc-seg--surplus"
                style={{ position: "static", width: sPct }}
              />
            </div>
            <span className="v1-ch10-series-rate">{rate.toFixed(0)}%</span>
          </div>
        );
      })}
    </div>
  );
}

function ValidatePanel() {
  const [nl, setNl] = useState(360);
  const [sl, setSl] = useState(360);
  const [limit, setLimit] = useState<number | "">("");
  const [result, setResult] = useState<ValidateWorkingDayResponse | null>(null);
  const [err, setErr] = useState<string | null>(null);
  const [seriesVisited, setSeriesVisited] = useState<Set<string>>(new Set());

  function loadFixture(necessary: number, surplus: number, lim?: number, seriesId?: string) {
    setNl(necessary);
    setSl(surplus);
    setLimit(lim ?? "");
    setResult(null);
    if (seriesId) {
      setSeriesVisited((prev) => {
        const next = new Set(prev);
        next.add(seriesId);
        return next;
      });
    }
  }

  async function submit(e: FormEvent) {
    e.preventDefault();
    setErr(null);
    try {
      const input = {
        necessary_labour_minutes: nl,
        surplus_labour_minutes: sl,
        ...(limit !== "" ? { statutory_limit_minutes: Number(limit) } : {}),
      };
      const r = await api.validateWorkingDay(input);
      setResult(r);
    } catch (e) {
      setErr(e instanceof Error ? e.message : String(e));
    }
  }

  const statutoryLimitMin = typeof limit === "number" ? limit : null;

  return (
    <section className="card">
      <h2>Working-Day Validator</h2>
      <p className="description">
        Validate a working day against the physical maximum (24 h) and an optional statutory
        limit. Click any §1 fixture to step through Marx's working-day series.
      </p>
      <div className="v1-ch10-presets">
        <span className="v1-ch10-presets-label">Marx fixture</span>
        <button
          type="button"
          className="v1-ch10-preset-button"
          onClick={() => loadFixture(360, 60, undefined, "wd1")}
        >
          §1 WD I (6h+1h)
        </button>
        <button
          type="button"
          className="v1-ch10-preset-button"
          onClick={() => loadFixture(360, 180, undefined, "wd2")}
        >
          §1 WD II (6h+3h)
        </button>
        <button
          type="button"
          className="v1-ch10-preset-button"
          onClick={() => loadFixture(360, 360, undefined, "wd3")}
        >
          §1 WD III (6h+6h)
        </button>
        <button
          type="button"
          className="v1-ch10-preset-button"
          onClick={() => loadFixture(315, 315, 630)}
        >
          Factory Act 1850 (10.5h)
        </button>
        <button
          type="button"
          className="v1-ch10-preset-button"
          onClick={() => loadFixture(84 * 480, 56 * 480)}
        >
          Wallachian Corvée
        </button>
      </div>

      <AbcBar necessary={nl} surplus={sl} statutoryLimit={statutoryLimitMin} />

      <form className="form-grid" onSubmit={submit}>
        <label>
          <span>Necessary Labour A–B (min)</span>
          <input type="number" min={1} value={nl} onChange={(e) => setNl(Number(e.target.value))} />
        </label>
        <label>
          <span>Surplus Labour B–C (min)</span>
          <input type="number" min={0} value={sl} onChange={(e) => setSl(Number(e.target.value))} />
        </label>
        <label>
          <span>Statutory Limit (min, optional)</span>
          <input
            type="number"
            min={1}
            value={limit}
            placeholder="e.g. 630"
            onChange={(e) => setLimit(e.target.value === "" ? "" : Number(e.target.value))}
          />
        </label>
        <div className="form-actions span2">
          <button type="submit" className="primary">Validate</button>
          {err && <span className="error">{err}</span>}
        </div>
      </form>

      {result && (
        <>
          <RateChip v={nl} s={sl} valid={result.valid} />
          {!result.valid && (
            <p className="small error" style={{ marginTop: "0.4rem" }}>
              Invalid — {result.error}
            </p>
          )}
        </>
      )}

      <RateSeriesChart visited={seriesVisited} />
    </section>
  );
}

function CreatePanel() {
  const [nl, setNl] = useState(360);
  const [sl, setSl] = useState(360);
  const [limit, setLimit] = useState<number | "">("");
  const [result, setResult] = useState<WorkingDayResponse | null>(null);
  const [err, setErr] = useState<string | null>(null);

  async function submit(e: FormEvent) {
    e.preventDefault();
    setErr(null);
    try {
      const input = {
        necessary_labour_minutes: nl,
        surplus_labour_minutes: sl,
        ...(limit !== "" ? { statutory_limit_minutes: Number(limit) } : {}),
      };
      const r = await api.createWorkingDay(input);
      setResult(r);
    } catch (e) {
      setErr(e instanceof Error ? e.message : String(e));
    }
  }

  return (
    <section className="card">
      <h2>Create Working-Day Record</h2>
      <p className="description">
        Persist a working-day record. Returns the stored ID alongside computed metrics.
      </p>
      <form className="form-grid" onSubmit={submit}>
        <label>
          <span>Necessary Labour A–B (min)</span>
          <input type="number" min={1} value={nl} onChange={(e) => setNl(Number(e.target.value))} />
        </label>
        <label>
          <span>Surplus Labour B–C (min)</span>
          <input type="number" min={0} value={sl} onChange={(e) => setSl(Number(e.target.value))} />
        </label>
        <label>
          <span>Statutory Limit (min, optional)</span>
          <input
            type="number"
            min={1}
            value={limit}
            placeholder="e.g. 630"
            onChange={(e) => setLimit(e.target.value === "" ? "" : Number(e.target.value))}
          />
        </label>
        <div className="form-actions span2">
          <button type="submit" className="primary">Create</button>
          {err && <span className="error">{err}</span>}
        </div>
      </form>
      {result && (
        <div>
          <p className="small muted">Saved: <span className="monospace">{result.working_day.id}</span></p>
          <RateChip v={nl} s={sl} />
        </div>
      )}
    </section>
  );
}

function RelaySchedulePanel() {
  const [nl0, setNl0] = useState(360);
  const [sl0, setSl0] = useState(360);
  const [nl1, setNl1] = useState(360);
  const [sl1, setSl1] = useState(0);
  const [result, setResult] = useState<RelaySchedule | null>(null);
  const [err, setErr] = useState<string | null>(null);

  async function submit(e: FormEvent) {
    e.preventDefault();
    setErr(null);
    try {
      const r = await api.createRelaySchedule({
        sets: [
          { shift_kind: "day", necessary_labour_minutes: nl0, surplus_labour_minutes: sl0, worker_ids: [] },
          { shift_kind: "night", necessary_labour_minutes: nl1, surplus_labour_minutes: sl1, worker_ids: [] },
        ],
      });
      setResult(r);
    } catch (e) {
      setErr(e instanceof Error ? e.message : String(e));
    }
  }

  const day = nl0 + sl0;
  const night = nl1 + sl1;
  const combined = day + night;
  const max = PHYSICAL_MAX_MIN;
  const pct = (n: number) => `${(n / max) * 100}%`;

  return (
    <section className="card">
      <h2>Relay Schedule</h2>
      <p className="description">
        A relay system alternates two worker sets to keep machinery running. Combined shifts must
        not exceed 24 h (1440 min).
      </p>
      <form className="form-grid" onSubmit={submit}>
        <label>
          <span>Day shift — Necessary (min)</span>
          <input type="number" min={1} value={nl0} onChange={(e) => setNl0(Number(e.target.value))} />
        </label>
        <label>
          <span>Day shift — Surplus (min)</span>
          <input type="number" min={0} value={sl0} onChange={(e) => setSl0(Number(e.target.value))} />
        </label>
        <label>
          <span>Night shift — Necessary (min)</span>
          <input type="number" min={1} value={nl1} onChange={(e) => setNl1(Number(e.target.value))} />
        </label>
        <label>
          <span>Night shift — Surplus (min)</span>
          <input type="number" min={0} value={sl1} onChange={(e) => setSl1(Number(e.target.value))} />
        </label>
        <div className="form-actions span2">
          <button type="submit" className="primary">Create Relay Schedule</button>
          {err && <span className="error">{err}</span>}
        </div>
      </form>

      <div className="v1-ch10-relay-stack">
        <div className="v1-ch10-relay-row">
          <span className="v1-ch10-relay-label">Day shift</span>
          <div className="v1-ch10-relay-track">
            <div
              className="v1-ch10-abc-seg v1-ch10-abc-seg--necessary"
              style={{ position: "static", width: pct(nl0) }}
              title={`Day necessary: ${minutesToHours(nl0)}`}
            />
            <div
              className="v1-ch10-abc-seg v1-ch10-abc-seg--surplus"
              style={{ position: "static", width: pct(sl0) }}
              title={`Day surplus: ${minutesToHours(sl0)}`}
            />
          </div>
          <span className="v1-ch10-relay-total">{minutesToHours(day)}</span>
        </div>
        <div className="v1-ch10-relay-row">
          <span className="v1-ch10-relay-label">Night shift</span>
          <div className="v1-ch10-relay-track">
            <div
              className="v1-ch10-abc-seg v1-ch10-abc-seg--necessary"
              style={{ position: "static", width: pct(nl1) }}
              title={`Night necessary: ${minutesToHours(nl1)}`}
            />
            <div
              className="v1-ch10-abc-seg v1-ch10-abc-seg--surplus"
              style={{ position: "static", width: pct(sl1) }}
              title={`Night surplus: ${minutesToHours(sl1)}`}
            />
          </div>
          <span className="v1-ch10-relay-total">{minutesToHours(night)}</span>
        </div>
        <div className="v1-ch10-relay-row">
          <span className="v1-ch10-relay-label">Combined</span>
          <div className="v1-ch10-relay-track">
            <div
              className="v1-ch10-abc-seg v1-ch10-abc-seg--necessary"
              style={{ position: "static", width: pct(Math.min(combined, max)) }}
              title={`Combined: ${minutesToHours(combined)} of ${minutesToHours(max)}`}
            />
          </div>
          <span
            className={
              "v1-ch10-relay-total" + (combined > max ? " v1-ch10-relay-total--over" : "")
            }
          >
            {minutesToHours(combined)} / 24h
          </span>
        </div>
      </div>

      {result && (
        <p className="small muted" style={{ marginTop: "0.85rem" }}>
          Saved: <span className="monospace">{result.id}</span>
        </p>
      )}
    </section>
  );
}

function Coda() {
  return (
    <aside className="v1-ch10-coda">
      <p className="v1-ch10-coda-quote">
        “Capital is dead labour which, vampire-like, lives only by sucking
        living labour, and lives the more, the more labour it sucks.”
        <span className="v1-ch10-coda-cite">
          — Marx, Capital Vol. I, Ch. 10 §1
        </span>
      </p>
    </aside>
  );
}
