import { useEffect, useMemo, useState } from "react";
import type { FormEvent } from "react";
import { api } from "../api";
import { RelativeSurplusBridge } from "../components/RelativeSurplusBridge";
import { fmtHoursLong as minutesToHours, fmtPounds as poundsFromPence } from "../format";
import type {
  Agent,
  CollectiveWorkingDayResult,
  Cooperation,
  CooperationMember,
  MinimumCapitalResult,
} from "../types";

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
        // panels handle their own errors
      }
    })();
  }, [refreshTick]);

  function refresh() {
    setRefreshTick((n) => n + 1);
  }

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
  const [memberIds, setMemberIds] = useState<string[]>([]);
  const [supervisorIds, setSupervisorIds] = useState<string[]>([]);
  const [workingDay, setWorkingDay] = useState(720);
  const [err, setErr] = useState<string | null>(null);

  useEffect(() => {
    if (!capitalistId && capitalists.length > 0) setCapitalistId(capitalists[0].id);
  }, [capitalists, capitalistId]);

  function toggleMember(id: string) {
    setMemberIds((cur) =>
      cur.includes(id) ? cur.filter((x) => x !== id) : [...cur, id]
    );
    setSupervisorIds((cur) => cur.filter((x) => x !== id || memberIds.includes(id)));
  }

  function toggleSupervisor(id: string) {
    setSupervisorIds((cur) =>
      cur.includes(id) ? cur.filter((x) => x !== id) : [...cur, id]
    );
  }

  async function submit(e: FormEvent) {
    e.preventDefault();
    setErr(null);
    if (!capitalistId) {
      setErr("Select a capitalist.");
      return;
    }
    if (memberIds.length === 0) {
      setErr("Add at least one worker.");
      return;
    }
    const members: CooperationMember[] = memberIds.map((id) => ({
      worker_id: id,
      supervisory: supervisorIds.includes(id),
      working_day_minutes: workingDay,
    }));
    try {
      await api.createCooperation({
        name: name.trim() || "cooperation",
        capitalist_id: capitalistId,
        members,
      });
      setMemberIds([]);
      setSupervisorIds([]);
      onCreated();
    } catch (e) {
      setErr(e instanceof Error ? e.message : String(e));
    }
  }

  return (
    <section className="card">
      <h2>Assemble a Co-operation</h2>
      <p className="description">
        Ch. 13 §1: "Capitalist production only really begins ... when each individual capital
        employs simultaneously a comparatively large number of labourers." Pick one capitalist
        and the workers they will assemble under a single command.
      </p>
      <form className="form-grid" onSubmit={submit}>
        <label>
          <span>Name</span>
          <input value={name} onChange={(e) => setName(e.target.value)} />
        </label>
        <label>
          <span>Capitalist</span>
          <select
            value={capitalistId}
            onChange={(e) => setCapitalistId(e.target.value)}
          >
            <option value="">— select —</option>
            {capitalists.map((c) => (
              <option key={c.id} value={c.id}>
                {c.name}
              </option>
            ))}
          </select>
        </label>
        <label>
          <span>Working Day per Worker (min)</span>
          <input
            type="number"
            min={1}
            max={1440}
            value={workingDay}
            onChange={(e) => setWorkingDay(Number(e.target.value))}
          />
        </label>
        <div className="span2">
          <p className="small muted" style={{ marginTop: 0 }}>
            Select workers (and optionally mark some as supervisory — Marx's "directing authority"
            who is still a wage-labourer, not a capitalist).
          </p>
          <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr 1fr", gap: "0.25rem" }}>
            {workers.length === 0 && (
              <p className="small muted">No worker-class agents yet — create some in Ch. 04.</p>
            )}
            {workers.map((w) => {
              const checked = memberIds.includes(w.id);
              const sup = supervisorIds.includes(w.id);
              return (
                <div key={w.id} style={{ display: "flex", alignItems: "center", gap: "0.4rem" }}>
                  <input
                    type="checkbox"
                    checked={checked}
                    onChange={() => toggleMember(w.id)}
                  />
                  <span style={{ flex: 1 }}>{w.name}</span>
                  <label style={{ display: "flex", alignItems: "center", gap: "0.2rem" }}>
                    <input
                      type="checkbox"
                      disabled={!checked}
                      checked={sup}
                      onChange={() => toggleSupervisor(w.id)}
                    />
                    <span className="small muted">sup</span>
                  </label>
                </div>
              );
            })}
          </div>
        </div>
        <div className="form-actions span2">
          <button type="submit" className="primary">
            Assemble
          </button>
          {err && <span className="error">{err}</span>}
        </div>
      </form>
    </section>
  );
}

function CooperationsList({
  cooperations,
  workers,
  capitalists,
}: {
  cooperations: Cooperation[];
  workers: Agent[];
  capitalists: Agent[];
}) {
  const [selectedId, setSelectedId] = useState<string | null>(null);
  const [cwd, setCwd] = useState<CollectiveWorkingDayResult | null>(null);
  const [err, setErr] = useState<string | null>(null);

  const workerName = useMemo(() => {
    const map = new Map(workers.map((w) => [w.id, w.name]));
    return (id: string) => map.get(id) ?? id.slice(0, 8);
  }, [workers]);

  const capitalistName = useMemo(() => {
    const map = new Map(capitalists.map((c) => [c.id, c.name]));
    return (id: string) => map.get(id) ?? id.slice(0, 8);
  }, [capitalists]);

  async function compute(id: string) {
    setSelectedId(id);
    setErr(null);
    try {
      const r = await api.computeCollectiveWorkingDay(id);
      setCwd(r);
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
        <p className="small muted">No cooperations yet.</p>
      ) : (
        <table className="data-table">
          <thead>
            <tr>
              <th>Name</th>
              <th>Capitalist</th>
              <th>Size (n)</th>
              <th>Collective Working-Day</th>
              <th>Avg Social Labour</th>
              <th>Output Factor</th>
              <th>Supervisors</th>
              <th></th>
            </tr>
          </thead>
          <tbody>
            {cooperations.map((c) => (
              <tr key={c.id}>
                <td>{c.name || c.id.slice(0, 8)}</td>
                <td>{capitalistName(c.capitalist_id)}</td>
                <td>{c.size}</td>
                <td>{minutesToHours(c.collective_working_day_minutes)}</td>
                <td>{minutesToHours(c.average_social_labour_minutes)}</td>
                <td>{c.collective_power_factor.toFixed(3)}×</td>
                <td>
                  {!c.supervisors || c.supervisors.length === 0
                    ? "—"
                    : c.supervisors.map((s) => workerName(s.worker_id)).join(", ")}
                </td>
                <td>
                  <button type="button" onClick={() => compute(c.id)}>
                    Compute
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
      {err && <p className="error">{err}</p>}
      {cwd && selectedId && (
        <div style={{ marginTop: "1rem", padding: "0.75rem", background: "var(--surface-2, #f4f4ee)" }}>
          <h3 style={{ marginTop: 0 }}>Computed</h3>
          <p className="small">
            Collective working-day: <strong>{minutesToHours(cwd.collective_working_day_minutes)}</strong>{" "}
            ({cwd.collective_working_day_minutes} min) — {cwd.size} workers × {minutesToHours(cwd.individual_working_day_minutes)}.
          </p>
          <p className="small">
            Cooperative output (use-values): <strong>{cwd.collective_output_use_value_units.toFixed(0)}</strong>{" "}
            min-equivalent (factor {cwd.collective_power_factor.toFixed(3)}×). Value remains {cwd.collective_working_day_minutes} min.
          </p>
          <p className="small muted">
            Cooperative origin appears as a property of <em>{cwd.cooperative_origin}</em>, not labour.
            The bonus costs the capitalist nothing — "because, on the other hand, the workman
            himself does not develop it before his labour belongs to capital." (§13)
          </p>
          <RelativeSurplusBridge
            source="cooperation"
            sourceId={selectedId}
            factor={cwd.collective_power_factor}
            sourceLabel={`Cooperation of ${cwd.size}`}
          />
        </div>
      )}
    </section>
  );
}

// ── Minimum capital ───────────────────────────────────────────────────────────

function MinimumCapitalPanel() {
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
      <div style={{ marginBottom: "0.75rem", display: "flex", gap: "0.5rem", flexWrap: "wrap" }}>
        <button type="button" onClick={() => { setWorkerCount(300); setWagePence(72); }}>
          §8 (300 @ 6 s.)
        </button>
        <button type="button" onClick={() => { setWorkerCount(10); setWagePence(72); }}>
          §8 (10 @ 6 s.)
        </button>
        <button type="button" onClick={() => { setWorkerCount(1200); setWagePence(72); }}>
          §2 (1,200 @ 6 s.)
        </button>
      </div>
      <form className="form-grid" onSubmit={submit}>
        <label>
          <span>Worker Count (n)</span>
          <input
            type="number"
            min={1}
            value={workerCount}
            onChange={(e) => setWorkerCount(Number(e.target.value))}
          />
        </label>
        <label>
          <span>Daily Wage (pence)</span>
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
        <div style={{ marginTop: "1rem" }}>
          <p>
            Minimum capital: <strong>{poundsFromPence(result.minimum_capital_pence)}</strong>{" "}
            ({result.minimum_capital_pence} pence) for {result.worker_count} workers at{" "}
            {poundsFromPence(result.daily_wage_pence)}/day.
          </p>
        </div>
      )}
    </section>
  );
}
