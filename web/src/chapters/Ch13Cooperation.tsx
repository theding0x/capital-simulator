// Drop-in replacement for web/src/chapters/Ch13Cooperation.tsx in
// theding0x/capital-simulator. Redesigns the Co-operation page:
//
//   - Worker roster becomes a proper data-table row list instead of the
//     3-column checkbox+"sup" grid (which had alignment issues — the
//     bottom-borders ran across gutters and "sup" read like a label, not
//     a control).
//   - A live composite-working-day stat row sits between the form and the
//     roster: labourers × per-worker day = collective labour-minutes. This
//     is the central claim of Ch. 13 §3 and was previously buried in the
//     post-assembly compute step.
//   - "supervisory" becomes a gold OVERSEER chip on selected rows only.
//   - Inline-expand for assembled co-operations to show their roster,
//     mirroring Ch. 02's reveal pattern.
//   - All form labels in the command row are single-line and inputs share
//     the same baseline (the original "Capitalist (single command)" label
//     wrapped, dropping the Name input out of alignment).
//
// Co-located CSS lives in Ch13Cooperation.css (imported below). It only
// adds a few classes scoped to .ch13-*; the existing data-table, card,
// form-grid, and reveal-panel rules from index.css carry the rest.

import { useEffect, useMemo, useState } from "react";
import type { FormEvent } from "react";
import { api } from "../api";
import { RelativeSurplusBridge } from "../components/RelativeSurplusBridge";
import { fmtHoursLong } from "../format";
import { usePounds } from "../CurrencyContext";
import type {
  Agent,
  CollectiveWorkingDayResult,
  Cooperation,
  CooperationMember,
  MinimumCapitalResult,
} from "../types";
import "./Ch13Cooperation.css";

interface Ch13Props {
  onSharedChanged: () => void;
}

export function Ch13Cooperation({ onSharedChanged: _unused }: Ch13Props) {
  return (
    <>
      <CooperationLedgerPanel />
      <MinimumCapitalPanel />
    </>
  );
}

// ── Cooperation ledger ────────────────────────────────────────────────────────

function CooperationLedgerPanel() {
  const [capitalists, setCapitalists] = useState<Agent[]>([]);
  const [workers, setWorkers] = useState<Agent[]>([]);
  const [cooperations, setCooperations] = useState<Cooperation[]>([]);
  const [refreshTick, setRefreshTick] = useState(0);

  useEffect(() => {
    (async () => {
      try {
        const [caps, ws, coops] = await Promise.all([
          api.listAgents("capitalist"),
          api.listAgents("worker"),
          api.listCooperations(),
        ]);
        setCapitalists(caps);
        setWorkers(ws);
        setCooperations(coops);
      } catch {
        // child panels handle their own errors
      }
    })();
  }, [refreshTick]);

  const refresh = () => setRefreshTick((n) => n + 1);

  return (
    <>
      <CreateCooperationPanel
        capitalists={capitalists}
        workers={workers}
        onCreated={refresh}
      />
      <CooperationsList
        cooperations={cooperations}
        workers={workers}
        capitalists={capitalists}
      />
    </>
  );
}

// ── Assemble panel ────────────────────────────────────────────────────────────

interface Selection {
  // worker_id -> { sup }
  [id: string]: { sup: boolean };
}

function CreateCooperationPanel({
  capitalists,
  workers,
  onCreated,
}: {
  capitalists: Agent[];
  workers: Agent[];
  onCreated: () => void;
}) {
  const [name, setName] = useState("Spinning floor");
  const [capitalistId, setCapitalistId] = useState("");
  const [workingDay, setWorkingDay] = useState(720);
  const [selection, setSelection] = useState<Selection>({});
  const [err, setErr] = useState<string | null>(null);

  useEffect(() => {
    if (!capitalistId && capitalists.length > 0) setCapitalistId(capitalists[0].id);
  }, [capitalists, capitalistId]);

  const memberCount = Object.keys(selection).length;
  const supCount = Object.values(selection).filter((s) => s.sup).length;
  const labourers = memberCount - supCount;
  const compositeDay = memberCount * workingDay;

  function toggle(id: string) {
    setSelection((cur) => {
      const next = { ...cur };
      if (id in next) delete next[id];
      else next[id] = { sup: false };
      return next;
    });
  }

  function toggleSup(id: string) {
    setSelection((cur) => {
      if (!(id in cur)) return cur;
      return { ...cur, [id]: { sup: !cur[id].sup } };
    });
  }

  function selectAll() {
    const next: Selection = {};
    workers.forEach((w) => { next[w.id] = { sup: false }; });
    setSelection(next);
  }

  function clear() {
    setSelection({});
  }

  async function submit(e: FormEvent) {
    e.preventDefault();
    setErr(null);
    if (!capitalistId) {
      setErr("Select a capitalist.");
      return;
    }
    if (memberCount < 2) {
      setErr('Select at least 2 workers — "a comparatively large number" of labourers.');
      return;
    }
    const members: CooperationMember[] = Object.entries(selection).map(
      ([worker_id, { sup }]) => ({
        worker_id,
        supervisory: sup,
        working_day_minutes: workingDay,
      })
    );
    try {
      await api.createCooperation({
        name: name.trim() || "cooperation",
        capitalist_id: capitalistId,
        members,
      });
      clear();
      onCreated();
    } catch (e) {
      setErr(e instanceof Error ? e.message : String(e));
    }
  }

  const capitalist = capitalists.find((c) => c.id === capitalistId);

  return (
    <section className="card">
      <h2>Assemble a Co-operation</h2>
      <p className="description">
        Ch. 13 §1: "Capitalist production only really begins ... when each individual capital
        employs simultaneously a comparatively large number of labourers." Pick one capitalist
        and the workers they will assemble under a single command.
      </p>

      <form onSubmit={submit}>
        {/* Command row — name + capitalist + per-worker day, baselines aligned. */}
        <div className="form-grid ch13-command-row">
          <label>
            <span>Co-operation name</span>
            <input
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="Spinning floor"
              required
            />
          </label>
          <label>
            <span>Capitalist</span>
            <select
              value={capitalistId}
              onChange={(e) => setCapitalistId(e.target.value)}
              required
            >
              <option value="">— select —</option>
              {capitalists.map((c) => (
                <option key={c.id} value={c.id}>{c.name}</option>
              ))}
            </select>
          </label>
          <label>
            <span>Day per worker (min)</span>
            <input
              type="number"
              min={1}
              max={1440}
              value={workingDay}
              onChange={(e) => setWorkingDay(Number(e.target.value))}
            />
          </label>
        </div>

        {/* Live composite-working-day readout. */}
        <CompositeSummary
          labourers={labourers}
          supervisors={supCount}
          perWorkerDay={workingDay}
          compositeDay={compositeDay}
          totalMembers={memberCount}
        />

        {/* Roster toolbar. */}
        <div className="ch13-roster-toolbar">
          <div>
            <h3 className="ch13-h3">Workers · {workers.length} available</h3>
            <p className="small muted ch13-roster-help">
              Click a worker to assign them to {capitalist?.name ?? "this capitalist"}. Mark
              some as <span className="ch13-overseer-inline">overseer</span> — Marx's "directing
              authority", still a wage-labourer.
            </p>
          </div>
          <div className="ch13-roster-actions">
            <button type="button" className="secondary" onClick={selectAll} disabled={workers.length === 0}>
              Select all
            </button>
            <button type="button" className="secondary" onClick={clear} disabled={memberCount === 0}>
              Clear
            </button>
          </div>
        </div>

        <RosterTable
          workers={workers}
          selection={selection}
          onToggle={toggle}
          onToggleSup={toggleSup}
        />

        <div className="ch13-submit-row">
          <span className="small muted">
            {memberCount < 2
              ? `Select at least 2 workers — "a comparatively large number" of labourers.`
              : `Ready to assemble ${memberCount} workers under ${capitalist?.name ?? "—"} for a composite working day of ${fmtHoursLong(compositeDay)}.`}
          </span>
          <div className="form-actions">
            <button type="submit" className="primary" disabled={memberCount < 2 || !capitalistId}>
              Assemble co-operation
            </button>
            {err && <span className="error">{err}</span>}
          </div>
        </div>
      </form>
    </section>
  );
}

// ── Composite summary ────────────────────────────────────────────────────────

function CompositeSummary({
  labourers,
  supervisors,
  perWorkerDay,
  compositeDay,
  totalMembers,
}: {
  labourers: number;
  supervisors: number;
  perWorkerDay: number;
  compositeDay: number;
  totalMembers: number;
}) {
  return (
    <div className="ch13-summary">
      <Stat label="Labourers" value={labourers} unit="under one command" />
      <Stat
        label="Overseers"
        value={supervisors}
        unit={supervisors === 1 ? "directing authority" : "directing authorities"}
        tone="gold"
      />
      <Stat
        label="Composite working day"
        value={compositeDay}
        unit={
          totalMembers
            ? `= ${totalMembers} × ${perWorkerDay} min  (${fmtHoursLong(compositeDay)})`
            : "—"
        }
        tone="lead"
        big
      />
      {totalMembers >= 2 && (
        <p className="ch13-gloss">
          <span aria-hidden>›</span>{" "}
          The composite labour-day of {totalMembers} co-operating workers yields more than
          {" "}the sum of {totalMembers} isolated labour-days. "The special productive power of
          the combined working day is, under all circumstances, the social productive power
          of labour, or the productive power of social labour." (Ch. 13 §3)
        </p>
      )}
    </div>
  );
}

function Stat({
  label, value, unit, tone, big,
}: {
  label: string;
  value: number;
  unit: string;
  tone?: "gold" | "lead";
  big?: boolean;
}) {
  const cls = [
    "ch13-stat-value",
    tone === "gold" && "ch13-stat-value--gold",
    tone === "lead" && "ch13-stat-value--lead",
    big && "ch13-stat-value--big",
  ].filter(Boolean).join(" ");
  return (
    <div className="ch13-stat">
      <div className="ch13-stat-label">{label}</div>
      <div className={cls}>{value}</div>
      <div className="ch13-stat-unit">{unit}</div>
    </div>
  );
}

// ── Roster table ─────────────────────────────────────────────────────────────

function RosterTable({
  workers, selection, onToggle, onToggleSup,
}: {
  workers: Agent[];
  selection: Selection;
  onToggle: (id: string) => void;
  onToggleSup: (id: string) => void;
}) {
  if (workers.length === 0) {
    return (
      <p className="small muted ch13-roster-empty">
        No worker-class agents yet — create some in Ch. 04.
      </p>
    );
  }
  return (
    <table className="data-table ch13-roster">
      <thead>
        <tr>
          <th className="ch13-roster-dot" />
          <th>Worker</th>
          <th>Role</th>
          <th className="ch13-roster-chip" />
        </tr>
      </thead>
      <tbody>
        {workers.map((w) => {
          const sel = selection[w.id];
          const included = !!sel;
          const sup = !!sel?.sup;
          return (
            <tr
              key={w.id}
              className={included ? "ch13-row ch13-row--on" : "ch13-row"}
              onClick={() => onToggle(w.id)}
            >
              <td className="ch13-roster-dot">
                <SelectionDot on={included} />
              </td>
              <td className="ch13-roster-name">{w.name}</td>
              <td className="ch13-roster-role">
                {included ? (
                  sup
                    ? <span className="ch13-role-overseer">overseer</span>
                    : <span className="muted">labourer</span>
                ) : (
                  <span className="ch13-role-empty">—</span>
                )}
              </td>
              <td
                className="ch13-roster-chip"
                onClick={(e) => e.stopPropagation()}
              >
                {included && (
                  <button
                    type="button"
                    className={sup ? "ch13-chip ch13-chip--on" : "ch13-chip"}
                    onClick={() => onToggleSup(w.id)}
                    title='Marx&apos;s "directing authority" — still a wage-labourer, not a capitalist'
                  >
                    Overseer
                  </button>
                )}
              </td>
            </tr>
          );
        })}
      </tbody>
    </table>
  );
}

function SelectionDot({ on }: { on: boolean }) {
  return (
    <span className={on ? "ch13-dot ch13-dot--on" : "ch13-dot"} aria-hidden>
      {on && <span className="ch13-dot-check">✓</span>}
    </span>
  );
}

// ── Assembled co-operations list ─────────────────────────────────────────────

function CooperationsList({
  cooperations, workers, capitalists,
}: {
  cooperations: Cooperation[];
  workers: Agent[];
  capitalists: Agent[];
}) {
  const [expandedId, setExpandedId] = useState<string | null>(null);
  const [computed, setComputed] = useState<CollectiveWorkingDayResult | null>(null);
  const [err, setErr] = useState<string | null>(null);

  const workerName = useMemo(() => {
    const map = new Map(workers.map((w) => [w.id, w.name]));
    return (id: string) => map.get(id) ?? id.slice(0, 8);
  }, [workers]);

  const capitalistName = useMemo(() => {
    const map = new Map(capitalists.map((c) => [c.id, c.name]));
    return (id: string) => map.get(id) ?? id.slice(0, 8);
  }, [capitalists]);

  async function expand(c: Cooperation) {
    if (expandedId === c.id) {
      setExpandedId(null);
      setComputed(null);
      return;
    }
    setExpandedId(c.id);
    setComputed(null);
    setErr(null);
    try {
      const r = await api.computeCollectiveWorkingDay(c.id);
      setComputed(r);
    } catch (e) {
      setErr(e instanceof Error ? e.message : String(e));
    }
  }

  return (
    <section className="card">
      <h2>Assembled Co-operations</h2>
      <p className="description">
        Each row aggregates n × individual working-day. The collective output factor (≥ 1.0)
        captures the qualitative bonus from §5 — "the new power that arises from the fusion
        of many forces into one single force." Value (§2) still scales linearly.
      </p>

      {cooperations.length === 0 ? (
        <div className="empty-state">
          <p>No co-operations assembled yet.</p>
          <p className="small">Assemble at least two workers under a single capitalist above.</p>
        </div>
      ) : (
        <table className="data-table">
          <thead>
            <tr>
              <th>Name</th>
              <th>Capitalist</th>
              <th>Size (n)</th>
              <th>Collective day</th>
              <th>Avg social labour</th>
              <th>Output factor</th>
              <th>Supervisors</th>
              <th />
            </tr>
          </thead>
          <tbody>
            {cooperations.map((c) => {
              const open = expandedId === c.id;
              return (
                <FragmentRow
                  key={c.id}
                  open={open}
                  onExpand={() => expand(c)}
                  cooperation={c}
                  computed={open ? computed : null}
                  err={open ? err : null}
                  workerName={workerName}
                  capitalistName={capitalistName}
                />
              );
            })}
          </tbody>
        </table>
      )}
    </section>
  );
}

function FragmentRow({
  cooperation: c, open, onExpand, computed, err, workerName, capitalistName,
}: {
  cooperation: Cooperation;
  open: boolean;
  onExpand: () => void;
  computed: CollectiveWorkingDayResult | null;
  err: string | null;
  workerName: (id: string) => string;
  capitalistName: (id: string) => string;
}) {
  return (
    <>
      <tr onClick={onExpand} className="ch13-coop-row">
        <td className="ch13-coop-name">{c.name || c.id.slice(0, 8)}</td>
        <td>{capitalistName(c.capitalist_id)}</td>
        <td className="ch13-coop-num">{c.size}</td>
        <td className="ch13-coop-num">{fmtHoursLong(c.collective_working_day_minutes)}</td>
        <td className="ch13-coop-num">{fmtHoursLong(c.average_social_labour_minutes)}</td>
        <td className="ch13-coop-num">{c.collective_power_factor.toFixed(3)}×</td>
        <td>
          {!c.supervisors || c.supervisors.length === 0
            ? <span className="muted">—</span>
            : c.supervisors.map((s) => workerName(s.worker_id)).join(", ")}
        </td>
        <td className="ch13-coop-actions">
          <button type="button" className="secondary" onClick={(e) => { e.stopPropagation(); onExpand(); }}>
            {open ? "Hide" : "Compute"}
          </button>
        </td>
      </tr>
      {open && (
        <tr>
          <td colSpan={8} className="reveal">
            <div className="reveal-panel">
              <h3>Collective working day · {c.name || c.id.slice(0, 8)}</h3>
              {err && <p className="error">{err}</p>}
              {!computed && !err && <p className="reveal-note">Computing…</p>}
              {computed && (
                <>
                  <p className="labour-statement">
                    {computed.size} workers × {fmtHoursLong(computed.individual_working_day_minutes)}{" "}
                    = <strong>{fmtHoursLong(computed.collective_working_day_minutes)}</strong>{" "}
                    ({computed.collective_working_day_minutes} min).
                  </p>
                  <p className="reveal-note">
                    Cooperative output (use-values):{" "}
                    <strong>{computed.collective_output_use_value_units.toFixed(0)}</strong>{" "}
                    min-equivalent (factor {computed.collective_power_factor.toFixed(3)}×).
                    Value remains {computed.collective_working_day_minutes} min — the bonus
                    appears as a property of <em>{computed.cooperative_origin}</em>, not labour,
                    and costs the capitalist nothing (§13).
                  </p>
                  <RelativeSurplusBridge
                    source="cooperation"
                    sourceId={c.id}
                    factor={computed.collective_power_factor}
                    sourceLabel={`Cooperation of ${computed.size}`}
                  />
                </>
              )}
            </div>
          </td>
        </tr>
      )}
    </>
  );
}

// ── Minimum capital panel (unchanged from previous version) ──────────────────

function MinimumCapitalPanel() {
  const fmt = usePounds();
  const [workerCount, setWorkerCount] = useState(300);
  const [wagePence, setWagePence] = useState(72);
  const [result, setResult] = useState<MinimumCapitalResult | null>(null);
  const [err, setErr] = useState<string | null>(null);

  async function submit(e: FormEvent) {
    e.preventDefault();
    setErr(null);
    try {
      const r = await api.computeMinimumCapital({
        worker_count: workerCount,
        daily_wage_pence: wagePence,
      });
      setResult(r);
    } catch (e) {
      setErr(e instanceof Error ? e.message : String(e));
    }
  }

  return (
    <section className="card">
      <h2>Minimum Capital for Co-operation</h2>
      <p className="description">
        Ch. 13 §8: "The payment of 300 workmen at once ... requires a greater outlay of capital
        than does the payment of a smaller number of men, week by week." The cooperation
        constraint is, before anything else, a capital-minimum constraint.
      </p>
      <div className="ch13-preset-row">
        <button type="button" onClick={() => { setWorkerCount(300); setWagePence(72); }}>§8 (300 @ 6 s.)</button>
        <button type="button" onClick={() => { setWorkerCount(10); setWagePence(72); }}>§8 (10 @ 6 s.)</button>
        <button type="button" onClick={() => { setWorkerCount(1200); setWagePence(72); }}>§2 (1,200 @ 6 s.)</button>
      </div>
      <form className="form-grid" onSubmit={submit}>
        <label>
          <span>Worker count (n)</span>
          <input
            type="number"
            min={1}
            value={workerCount}
            onChange={(e) => setWorkerCount(Number(e.target.value))}
          />
        </label>
        <label>
          <span>Daily wage (pence)</span>
          <input
            type="number"
            min={1}
            value={wagePence}
            onChange={(e) => setWagePence(Number(e.target.value))}
          />
        </label>
        <div className="form-actions span2">
          <button type="submit" className="primary">Compute</button>
          {err && <span className="error">{err}</span>}
        </div>
      </form>
      {result && (
        <p className="ch13-min-capital-result">
          Minimum capital: <strong>{fmt(result.minimum_capital_pence)}</strong>{" "}
          ({result.minimum_capital_pence} pence) for {result.worker_count} workers at{" "}
          {fmt(result.daily_wage_pence)}/day.
        </p>
      )}
    </section>
  );
}
