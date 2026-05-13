import { useEffect, useMemo, useState } from "react";
import type { FormEvent } from "react";
import { api } from "../api";
import { RelativeSurplusBridge } from "../components/RelativeSurplusBridge";
import { fmtHoursLong as minutesToHours, fmtPounds as poundsFromPence } from "../format";
import type {
  Agent,
  DetailRole,
  Manufacture,
  ManufactureForm,
  ManufactureOrigin,
  ProportionalGroupSizeResult,
  ManufactureMinimumCapitalResult,
  SkillLevel,
} from "../types";

interface Ch14Props {
  onSharedChanged: () => void;
}

const FORMS: ManufactureForm[] = ["heterogeneous", "serial"];
const ORIGINS: ManufactureOrigin[] = ["combination", "splitting"];
const SKILLS: SkillLevel[] = ["skilled", "unskilled"];

export function Ch14Manufacture({ onSharedChanged: _ }: Ch14Props) {
  const [capitalists, setCapitalists] = useState<Agent[]>([]);
  const [manufactures, setManufactures] = useState<Manufacture[]>([]);
  const [refreshTick, setRefreshTick] = useState(0);

  useEffect(() => {
    (async () => {
      try {
        const [caps, mfs] = await Promise.all([
          api.listAgents("capitalist"),
          api.listManufactures(),
        ]);
        setCapitalists(caps);
        setManufactures(mfs);
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
      <CreateManufacturePanel capitalists={capitalists} onCreated={refresh} />
      <ManufacturesListPanel manufactures={manufactures} onChanged={refresh} />
    </>
  );
}

// ── Create ────────────────────────────────────────────────────────────────────

interface RoleDraft {
  name: string;
  skill_level: SkillLevel;
  output_rate_per_hour: number;
  head_count: number;
  tool_name: string;
}

const SEED_NEEDLE: RoleDraft[] = [
  { name: "wire-draw", skill_level: "skilled", output_rate_per_hour: 100, head_count: 1, tool_name: "draw-plate" },
  { name: "wire-straighten", skill_level: "unskilled", output_rate_per_hour: 100, head_count: 1, tool_name: "straightening-frame" },
  { name: "wire-cut", skill_level: "unskilled", output_rate_per_hour: 100, head_count: 1, tool_name: "shears" },
  { name: "wire-point", skill_level: "skilled", output_rate_per_hour: 100, head_count: 1, tool_name: "grinder" },
];

const SEED_TYPE: RoleDraft[] = [
  { name: "founder", skill_level: "skilled", output_rate_per_hour: 2000, head_count: 4, tool_name: "casting-mould" },
  { name: "breaker", skill_level: "unskilled", output_rate_per_hour: 4000, head_count: 2, tool_name: "breaking-iron" },
  { name: "rubber", skill_level: "unskilled", output_rate_per_hour: 8000, head_count: 1, tool_name: "polishing-stone" },
];

function CreateManufacturePanel({
  capitalists,
  onCreated,
}: {
  capitalists: Agent[];
  onCreated: () => void;
}) {
  const [name, setName] = useState("Type-foundry");
  const [capitalistId, setCapitalistId] = useState("");
  const [form, setForm] = useState<ManufactureForm>("serial");
  const [origin, setOrigin] = useState<ManufactureOrigin>("splitting");
  const [workingDay, setWorkingDay] = useState(720);
  const [roles, setRoles] = useState<RoleDraft[]>(SEED_TYPE);
  const [err, setErr] = useState<string | null>(null);

  useEffect(() => {
    if (!capitalistId && capitalists.length > 0) setCapitalistId(capitalists[0].id);
  }, [capitalists, capitalistId]);

  function updateRole(i: number, patch: Partial<RoleDraft>) {
    setRoles((rs) => rs.map((r, idx) => (idx === i ? { ...r, ...patch } : r)));
  }
  function addRole() {
    setRoles((rs) => [
      ...rs,
      { name: `op-${rs.length + 1}`, skill_level: "unskilled", output_rate_per_hour: 100, head_count: 1, tool_name: "" },
    ]);
  }
  function removeRole(i: number) {
    setRoles((rs) => rs.filter((_, idx) => idx !== i));
  }

  async function submit(e: FormEvent) {
    e.preventDefault();
    setErr(null);
    if (!capitalistId) {
      setErr("Select a capitalist.");
      return;
    }
    if (roles.length === 0) {
      setErr("Add at least one detail role.");
      return;
    }
    try {
      await api.createManufacture({
        name: name.trim() || "manufacture",
        capitalist_id: capitalistId,
        form,
        origin,
        individual_working_day_minutes: workingDay,
        roles: roles.map<DetailRole>((r) => ({
          name: r.name.trim(),
          skill_level: r.skill_level,
          output_rate_per_hour: r.output_rate_per_hour,
          head_count: r.head_count,
          tool_name: r.tool_name.trim(),
        })),
      });
      onCreated();
    } catch (e) {
      setErr(e instanceof Error ? e.message : String(e));
    }
  }

  return (
    <section className="card">
      <h2>Establish a Manufacture</h2>
      <p className="description">
        Ch. 14 §1: a Manufacture is a Co-operation reorganised around a fixed division of detail
        labour. Pick its two-fold origin — combination of independent handicrafts (carriage shop)
        or splitting of one handicraft into successive operations (needle, type) — and its form:
        heterogeneous (independent parts assembled) or serial (one product through ordered stages).
        Detail labour produces no commodities; only the collective product enters exchange (§4).
      </p>
      <div style={{ marginBottom: "0.75rem", display: "flex", gap: "0.5rem", flexWrap: "wrap" }}>
        <button type="button" onClick={() => { setName("Type-foundry"); setForm("serial"); setOrigin("splitting"); setRoles(SEED_TYPE); }}>
          §2 type-manufacture (4 / 2 / 1)
        </button>
        <button type="button" onClick={() => { setName("Needle works"); setForm("serial"); setOrigin("splitting"); setRoles(SEED_NEEDLE); }}>
          §2 needle wire (serial)
        </button>
      </div>
      <form className="form-grid" onSubmit={submit}>
        <label>
          <span>Name</span>
          <input value={name} onChange={(e) => setName(e.target.value)} />
        </label>
        <label>
          <span>Capitalist</span>
          <select value={capitalistId} onChange={(e) => setCapitalistId(e.target.value)}>
            <option value="">— select —</option>
            {capitalists.map((c) => (
              <option key={c.id} value={c.id}>{c.name}</option>
            ))}
          </select>
        </label>
        <label>
          <span>Form</span>
          <select value={form} onChange={(e) => setForm(e.target.value as ManufactureForm)}>
            {FORMS.map((f) => <option key={f} value={f}>{f}</option>)}
          </select>
        </label>
        <label>
          <span>Origin</span>
          <select value={origin} onChange={(e) => setOrigin(e.target.value as ManufactureOrigin)}>
            {ORIGINS.map((o) => <option key={o} value={o}>{o}</option>)}
          </select>
        </label>
        <label>
          <span>Individual Working-Day (min)</span>
          <input
            type="number" min={1} max={1440}
            value={workingDay}
            onChange={(e) => setWorkingDay(Number(e.target.value))}
          />
        </label>
        <div className="span2">
          <p className="small muted" style={{ marginTop: 0, marginBottom: "0.4rem" }}>
            Detail roles — each is a fractional function; workers are "annexed" to exactly one.
          </p>
          <table className="data-table">
            <thead>
              <tr>
                <th>Name</th>
                <th>Skill</th>
                <th>Output/hr</th>
                <th>Headcount</th>
                <th>Tool</th>
                <th></th>
              </tr>
            </thead>
            <tbody>
              {roles.map((r, i) => (
                <tr key={i}>
                  <td><input value={r.name} onChange={(e) => updateRole(i, { name: e.target.value })} /></td>
                  <td>
                    <select value={r.skill_level} onChange={(e) => updateRole(i, { skill_level: e.target.value as SkillLevel })}>
                      {SKILLS.map((s) => <option key={s} value={s}>{s}</option>)}
                    </select>
                  </td>
                  <td>
                    <input
                      type="number" min={1}
                      value={r.output_rate_per_hour}
                      onChange={(e) => updateRole(i, { output_rate_per_hour: Number(e.target.value) })}
                    />
                  </td>
                  <td>
                    <input
                      type="number" min={1}
                      value={r.head_count}
                      onChange={(e) => updateRole(i, { head_count: Number(e.target.value) })}
                    />
                  </td>
                  <td><input value={r.tool_name} onChange={(e) => updateRole(i, { tool_name: e.target.value })} /></td>
                  <td><button type="button" onClick={() => removeRole(i)}>×</button></td>
                </tr>
              ))}
            </tbody>
          </table>
          <button type="button" onClick={addRole} style={{ marginTop: "0.5rem" }}>+ Add role</button>
        </div>
        <div className="form-actions span2">
          <button type="submit" className="primary">Establish</button>
          {err && <span className="error">{err}</span>}
        </div>
      </form>
    </section>
  );
}

// ── List + per-manufacture controls ───────────────────────────────────────────

function ManufacturesListPanel({
  manufactures,
  onChanged,
}: {
  manufactures: Manufacture[];
  onChanged: () => void;
}) {
  const [selectedId, setSelectedId] = useState<string | null>(null);

  useEffect(() => {
    if (manufactures.length === 0) {
      setSelectedId(null);
    } else if (!selectedId || !manufactures.find((m) => m.id === selectedId)) {
      setSelectedId(manufactures[0].id);
    }
  }, [manufactures, selectedId]);

  const selected = useMemo(
    () => manufactures.find((m) => m.id === selectedId) ?? null,
    [manufactures, selectedId]
  );

  return (
    <section className="card">
      <h2>Established Manufactures</h2>
      <p className="description">
        Each row shows the collective labourer (§4) and its productive power factor — the
        cooperation bonus compounded with the division-of-labour bonus that grows with the number
        of detail roles. Value still scales linearly with workers; the bonus is in use-values.
      </p>
      {manufactures.length === 0 ? (
        <p className="small muted">No manufactures yet.</p>
      ) : (
        <table className="data-table">
          <thead>
            <tr>
              <th>Name</th>
              <th>Form</th>
              <th>Origin</th>
              <th>Roles</th>
              <th>Workers</th>
              <th>Output Factor</th>
              <th>Productive Power</th>
              <th></th>
            </tr>
          </thead>
          <tbody>
            {manufactures.map((m) => (
              <tr key={m.id} style={selectedId === m.id ? { background: "var(--surface-2, #f4f4ee)" } : undefined}>
                <td>{m.name || m.id.slice(0, 8)}</td>
                <td>{m.form}</td>
                <td>{m.origin}</td>
                <td>{m.roles.length}</td>
                <td>{m.total_workers}</td>
                <td>{m.productive_power_factor.toFixed(3)}×</td>
                <td>{minutesToHours(m.productive_power_minutes)}</td>
                <td><button type="button" onClick={() => setSelectedId(m.id)}>Inspect</button></td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
      {selected && (
        <div style={{ marginTop: "1rem", display: "grid", gap: "1rem" }}>
          <HierarchyTable m={selected} />
          <ProportionalGroupPanel m={selected} />
          <ScalePanel m={selected} onScaled={onChanged} />
          <MinimumCapitalPanel m={selected} />
          <RelativeSurplusBridge
            source="manufacture"
            sourceId={selected.id}
            factor={selected.productive_power_factor}
            sourceLabel={`Manufacture ${selected.name || selected.id.slice(0, 8)}`}
          />
        </div>
      )}
    </section>
  );
}

function HierarchyTable({ m }: { m: Manufacture }) {
  return (
    <div style={{ padding: "0.75rem", background: "var(--surface-2, #f4f4ee)" }}>
      <h3 style={{ marginTop: 0 }}>Labour Hierarchy — {m.name}</h3>
      <p className="small muted">
        §5: manufacture begets a skill hierarchy; unskilled labour-power costs less because
        apprenticeship cost vanishes.
      </p>
      <table className="data-table">
        <thead>
          <tr>
            <th>Role</th>
            <th>Skill</th>
            <th>Output / hr</th>
            <th>Headcount</th>
            <th>Tool</th>
          </tr>
        </thead>
        <tbody>
          {m.hierarchy.map((r) => (
            <tr key={r.name}>
              <td>{r.name}</td>
              <td>{r.skill_level}</td>
              <td>{r.output_rate_per_hour.toLocaleString()}</td>
              <td>{r.head_count}</td>
              <td>{r.tool_name || "—"}</td>
            </tr>
          ))}
        </tbody>
      </table>
      <p className="small" style={{ marginTop: "0.6rem" }}>
        Collective labourer: <strong>{m.collective_labourer.total_workers}</strong> workers; output
        per {minutesToHours(m.collective_labourer.period_minutes)} =
        {" "}<strong>{m.collective_labourer.output_per_period_minutes.toLocaleString()}</strong>{" "}
        units (bottleneck-bounded). Paralysed if any role goes absent:{" "}
        <strong>{m.collective_labourer.is_paralysed_if_absent_one ? "yes" : "no"}</strong>.
      </p>
    </div>
  );
}

function ProportionalGroupPanel({ m }: { m: Manufacture }) {
  const [target, setTarget] = useState<number>(() => {
    const rates = m.roles.map((r) => r.output_rate_per_hour);
    return rates.length ? Math.max(...rates) : 100;
  });
  const [result, setResult] = useState<ProportionalGroupSizeResult | null>(null);
  const [err, setErr] = useState<string | null>(null);

  useEffect(() => {
    const rates = m.roles.map((r) => r.output_rate_per_hour);
    setTarget(rates.length ? Math.max(...rates) : 100);
    setResult(null);
    setErr(null);
  }, [m.id, m.roles]);

  async function submit(e: FormEvent) {
    e.preventDefault();
    setErr(null);
    setResult(null);
    try {
      const r = await api.proportionalGroupSize(m.id, target);
      setResult(r);
    } catch (e) {
      setErr(e instanceof Error ? e.message : String(e));
    }
  }

  return (
    <div style={{ padding: "0.75rem", background: "var(--surface-2, #f4f4ee)" }}>
      <h3 style={{ marginTop: 0 }}>Proportional Group Size — {m.name}</h3>
      <p className="small muted">
        §2: required headcount per role so partial products flow without bottlenecks. Marx:
        "the founder casts 2,000 type an hour, the breaker breaks up 4,000, and the rubber polishes
        8,000" → 4 / 2 / 1 at target = 8,000.
      </p>
      <form className="form-grid" onSubmit={submit}>
        <label>
          <span>Target output / hour</span>
          <input type="number" min={1} value={target} onChange={(e) => setTarget(Number(e.target.value))} />
        </label>
        <div className="form-actions">
          <button type="submit" className="primary">Compute</button>
        </div>
      </form>
      {err && <p className="error">{err}</p>}
      {result && (
        <table className="data-table" style={{ marginTop: "0.6rem" }}>
          <thead><tr><th>Role</th><th>Headcount required</th></tr></thead>
          <tbody>
            {Object.entries(result.headcounts).map(([role, hc]) => (
              <tr key={role}><td>{role}</td><td>{hc}</td></tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  );
}

function ScalePanel({ m, onScaled }: { m: Manufacture; onScaled: () => void }) {
  const [k, setK] = useState(2);
  const [preview, setPreview] = useState<Manufacture | null>(null);
  const [err, setErr] = useState<string | null>(null);

  async function submit(e: FormEvent) {
    e.preventDefault();
    setErr(null);
    setPreview(null);
    try {
      const r = await api.scaleManufacture(m.id, k);
      setPreview(r);
      onScaled();
    } catch (e) {
      setErr(e instanceof Error ? e.message : String(e));
    }
  }

  return (
    <div style={{ padding: "0.75rem", background: "var(--surface-2, #f4f4ee)" }}>
      <h3 style={{ marginTop: 0 }}>Scale — {m.name}</h3>
      <p className="small muted">
        §3: once the proportion is set, scale extends only by integer multiples of every group.
        Tripling the manufacture triples every role's headcount and triples collective output.
      </p>
      <form className="form-grid" onSubmit={submit}>
        <label>
          <span>Multiplier (k ≥ 1)</span>
          <input type="number" min={1} value={k} onChange={(e) => setK(Number(e.target.value))} />
        </label>
        <div className="form-actions">
          <button type="submit" className="primary">Scale (preview)</button>
        </div>
      </form>
      {err && <p className="error">{err}</p>}
      {preview && (
        <p className="small" style={{ marginTop: "0.5rem" }}>
          Scaled preview: <strong>{preview.total_workers}</strong> workers (was {m.total_workers});
          productive power = <strong>{minutesToHours(preview.productive_power_minutes)}</strong>;
          output factor {preview.productive_power_factor.toFixed(3)}×.
        </p>
      )}
    </div>
  );
}

function MinimumCapitalPanel({ m }: { m: Manufacture }) {
  const [factor, setFactor] = useState(0.1);
  const [result, setResult] = useState<ManufactureMinimumCapitalResult | null>(null);
  const [err, setErr] = useState<string | null>(null);

  useEffect(() => {
    setResult(null);
    setErr(null);
  }, [m.id]);

  async function compute(e: FormEvent) {
    e.preventDefault();
    setErr(null);
    try {
      const r = await api.getManufactureMinimumCapital(m.id, factor);
      setResult(r);
    } catch (e) {
      setErr(e instanceof Error ? e.message : String(e));
    }
  }

  return (
    <div style={{ padding: "0.75rem", background: "var(--surface-2, #f4f4ee)" }}>
      <h3 style={{ marginTop: 0 }}>Minimum Capital — {m.name}</h3>
      <p className="small muted">
        §5: "the minimum amount of capital ... must keep increasing" because raw material is
        consumed at the faster collective rate. Manufacture's minimum always exceeds the
        simple-cooperation minimum of the same headcount.
      </p>
      <form className="form-grid" onSubmit={compute}>
        <label>
          <span>Raw-material cost factor (pence / labour-minute)</span>
          <input
            type="number" step="0.001" min={0}
            value={factor}
            onChange={(e) => setFactor(Number(e.target.value))}
          />
        </label>
        <div className="form-actions">
          <button type="submit" className="primary">Compute</button>
        </div>
      </form>
      {err && <p className="error">{err}</p>}
      {result && (
        <p className="small" style={{ marginTop: "0.5rem" }}>
          Manufacture minimum: <strong>{poundsFromPence(result.minimum_capital_pence)}</strong>{" "}
          ({result.minimum_capital_pence.toLocaleString()} pence). Simple-cooperation baseline:{" "}
          {poundsFromPence(result.cooperation_baseline_pence)}.
        </p>
      )}
    </div>
  );
}
