import { useEffect, useMemo, useState } from "react";
import type { FormEvent } from "react";
import { api } from "../../api";
import { RelativeSurplusBridge } from "../../components/RelativeSurplusBridge";
import { fmtCompact as compactNumber, fmtLabourMinutes as poundsFromLabourMinutes } from "../../format";
import type {
  CreateFactoryInput,
  Factory,
  FactoryTick,
  FactoryTickResult,
  Machine,
  MachineEntryInput,
  PrimeMoverKind,
} from "../../types";

interface Ch15Props {
  onSharedChanged: () => void;
}

const PRIME_MOVERS: PrimeMoverKind[] = ["steam", "water", "electric", "animal"];

export function Ch15Machinery({ onSharedChanged: _ }: Ch15Props) {
  const [machines, setMachines] = useState<Machine[]>([]);
  const [factories, setFactories] = useState<Factory[]>([]);
  const [refreshTick, setRefreshTick] = useState(0);
  const [lastTick, setLastTick] = useState<FactoryTickResult | null>(null);

  useEffect(() => {
    (async () => {
      try {
        const [ms, fs] = await Promise.all([api.listMachines(), api.listFactories()]);
        setMachines(ms);
        setFactories(fs);
      } catch {
        // panels surface their own errors
      }
    })();
  }, [refreshTick]);

  function refresh() {
    setRefreshTick((n) => n + 1);
  }

  return (
    <>
      <RegisterMachinePanel onCreated={refresh} />
      <MachineRegistryPanel machines={machines} />
      <AssembleFactoryPanel machines={machines} onCreated={refresh} />
      <FactoryFloorPanel
        factories={factories}
        lastTick={lastTick}
        onTick={async (id) => {
          try {
            const r = await api.tickFactory(id);
            setLastTick(r);
            refresh();
          } catch {
            // ignore
          }
        }}
      />
      <CapitalCompositionPanel />
    </>
  );
}

// ── Register a machine ────────────────────────────────────────────────────────

function RegisterMachinePanel({ onCreated }: { onCreated: () => void }) {
  const [name, setName] = useState("Needle-machine");
  const [motor, setMotor] = useState("steam piston");
  const [transmission, setTransmission] = useState("shaft & belt");
  const [tool, setTool] = useState("needle die");
  const [machineValue, setMachineValue] = useState(600_000);
  const [lifespan, setLifespan] = useState(1000);
  const [productivePower, setProductivePower] = useState(145_000);
  const [handLabour, setHandLabour] = useState(0);
  const [err, setErr] = useState<string | null>(null);

  function seedNeedleMachine() {
    setName("Needle-machine");
    setMotor("steam piston");
    setTransmission("shaft & belt");
    setTool("needle die");
    setMachineValue(600_000);
    setLifespan(1000);
    setProductivePower(145_000);
    setHandLabour(0);
  }

  function seedSteamPlough() {
    setName("Steam-plough");
    setMotor("steam piston");
    setTransmission("traction cable");
    setTool("plough share");
    setMachineValue(720_000);
    setLifespan(500);
    setProductivePower(1);
    setHandLabour(60 * 66); // 66 men × 60 min = 3,960 lab-min displaced per ploughing hour
  }

  async function submit(e: FormEvent) {
    e.preventDefault();
    setErr(null);
    try {
      await api.createMachine({
        name: name.trim() || "machine",
        motor_mechanism: motor.trim(),
        transmitting_mechanism: transmission.trim(),
        working_tool: tool.trim(),
        machine_value: machineValue,
        lifespan_days: lifespan,
        productive_power: productivePower,
        hand_labour_per_unit: handLabour,
      });
      onCreated();
    } catch (e) {
      setErr(e instanceof Error ? e.message : String(e));
    }
  }

  return (
    <section className="card">
      <h2>Register a Machine</h2>
      <p className="description">
        Ch. 15 §1: a Machine is the three-part instrument — motor mechanism,
        transmitting mechanism, working tool — driven by a Prime Mover. Its
        MachineValue is the labour crystallised in it at production; it
        transfers that value to product only by bits, at the rate of wear and
        tear (§2). Set ProductivePower to the daily use-value output (one
        needle-machine: 145,000/day).
      </p>
      <div style={{ marginBottom: "0.75rem", display: "flex", gap: "0.5rem", flexWrap: "wrap" }}>
        <button type="button" onClick={seedNeedleMachine}>
          §8A needle-machine (145,000/day)
        </button>
        <button type="button" onClick={seedSteamPlough}>
          §2 steam-plough (66 men/hr displaced)
        </button>
      </div>
      <form className="form-grid" onSubmit={submit}>
        <label>
          <span>Name</span>
          <input value={name} onChange={(e) => setName(e.target.value)} />
        </label>
        <label>
          <span>Motor mechanism</span>
          <input value={motor} onChange={(e) => setMotor(e.target.value)} />
        </label>
        <label>
          <span>Transmitting mechanism</span>
          <input value={transmission} onChange={(e) => setTransmission(e.target.value)} />
        </label>
        <label>
          <span>Working tool</span>
          <input value={tool} onChange={(e) => setTool(e.target.value)} />
        </label>
        <label>
          <span>Machine value (lab-min)</span>
          <input type="number" min={1} value={machineValue} onChange={(e) => setMachineValue(Number(e.target.value))} />
        </label>
        <label>
          <span>Lifespan (working days)</span>
          <input type="number" min={1} value={lifespan} onChange={(e) => setLifespan(Number(e.target.value))} />
        </label>
        <label>
          <span>Productive power (units/day)</span>
          <input type="number" min={1} value={productivePower} onChange={(e) => setProductivePower(Number(e.target.value))} />
        </label>
        <label>
          <span>Hand labour displaced (lab-min/unit)</span>
          <input type="number" min={0} value={handLabour} onChange={(e) => setHandLabour(Number(e.target.value))} />
        </label>
        <button type="submit" className="primary">Register</button>
        {err && <span className="error">{err}</span>}
      </form>
    </section>
  );
}

// ── Machine registry ──────────────────────────────────────────────────────────

function MachineRegistryPanel({ machines }: { machines: Machine[] }) {
  if (machines.length === 0) {
    return (
      <section className="card">
        <h2>Machine registry</h2>
        <p className="muted">No machines yet. Register one above.</p>
      </section>
    );
  }
  return (
    <section className="card">
      <h2>Machine registry</h2>
      <p className="description">
        Daily wear and tear is <code>MachineValue ÷ LifespanDays</code> (§2);
        per-unit value is wear divided by ProductivePower. Labour displaced is
        the hand labour replaced for the day&apos;s output — the §2 advantage
        of the machine over hand labour.
      </p>
      <table className="data-table">
        <thead>
          <tr>
            <th>Name</th>
            <th>Value</th>
            <th>Lifespan</th>
            <th>Productive power</th>
            <th>Daily wear</th>
            <th>Per-unit wear</th>
            <th>Labour displaced/day</th>
            <th>Accumulated wear</th>
          </tr>
        </thead>
        <tbody>
          {machines.map((m) => (
            <tr key={m.id}>
              <td>{m.name}</td>
              <td>{compactNumber(m.machine_value)} lab-min</td>
              <td>{compactNumber(m.lifespan_days)} d</td>
              <td>{compactNumber(m.productive_power)}/d</td>
              <td>{poundsFromLabourMinutes(m.daily_wear_and_tear)}</td>
              <td>{m.value_transferred_per_unit} lab-min</td>
              <td>{compactNumber(m.daily_labour_displaced)} lab-min</td>
              <td>{compactNumber(m.accumulated_wear)} lab-min</td>
            </tr>
          ))}
        </tbody>
      </table>
    </section>
  );
}

// ── Assemble a factory ────────────────────────────────────────────────────────

function AssembleFactoryPanel({
  machines,
  onCreated,
}: {
  machines: Machine[];
  onCreated: () => void;
}) {
  const [name, setName] = useState("Needle Works");
  const [primeMover, setPrimeMover] = useState<PrimeMoverKind>("steam");
  const [horsepower, setHorsepower] = useState(30);
  const [intensityFactor, setIntensityFactor] = useState(1.0);
  const [picked, setPicked] = useState<Set<string>>(new Set());
  const [err, setErr] = useState<string | null>(null);

  function toggle(id: string) {
    setPicked((prev) => {
      const next = new Set(prev);
      if (next.has(id)) next.delete(id);
      else next.add(id);
      return next;
    });
  }

  async function submit(e: FormEvent) {
    e.preventDefault();
    setErr(null);
    if (picked.size === 0) {
      setErr("Pick at least one machine.");
      return;
    }
    const entries: MachineEntryInput[] = Array.from(picked).map((id) => ({ id }));
    const input: CreateFactoryInput = {
      name: name.trim() || "factory",
      prime_mover: { kind: primeMover, horsepower },
      machines: entries,
      intensity_factor: intensityFactor,
    };
    try {
      await api.createFactory(input);
      setPicked(new Set());
      onCreated();
    } catch (e) {
      setErr(e instanceof Error ? e.message : String(e));
    }
  }

  return (
    <section className="card">
      <h2>Assemble a Factory</h2>
      <p className="description">
        Ch. 15 §1c: a Factory is an organised system of machines driven by one
        Prime Mover. Pick machines from the registry above and choose the
        motive power. The factory&apos;s daily value transfer is the sum of
        wear across its machines; its productive power is the sum of theirs.
      </p>
      <form className="form-grid" onSubmit={submit}>
        <label>
          <span>Name</span>
          <input value={name} onChange={(e) => setName(e.target.value)} />
        </label>
        <label>
          <span>Prime mover</span>
          <select value={primeMover} onChange={(e) => setPrimeMover(e.target.value as PrimeMoverKind)}>
            {PRIME_MOVERS.map((k) => <option key={k} value={k}>{k}</option>)}
          </select>
        </label>
        <label>
          <span>Horsepower</span>
          <input type="number" min={1} value={horsepower} onChange={(e) => setHorsepower(Number(e.target.value))} />
        </label>
        <label>
          <span>Intensity factor (§3C, ≥ 1.0)</span>
          <input
            type="number"
            min={1}
            step={0.1}
            value={intensityFactor}
            onChange={(e) => setIntensityFactor(Number(e.target.value))}
          />
        </label>
        <fieldset style={{ gridColumn: "1 / -1" }}>
          <legend>Machines</legend>
          {machines.length === 0 ? (
            <p className="muted">Register a machine first.</p>
          ) : (
            machines.map((m) => (
              <label key={m.id} style={{ display: "block" }}>
                <input
                  type="checkbox"
                  checked={picked.has(m.id)}
                  onChange={() => toggle(m.id)}
                />{" "}
                {m.name} &mdash; {compactNumber(m.productive_power)}/d
              </label>
            ))
          )}
        </fieldset>
        <button type="submit" className="primary">Assemble</button>
        {err && <span className="error">{err}</span>}
      </form>
    </section>
  );
}

// ── Factory floor ─────────────────────────────────────────────────────────────

function FactoryFloorPanel({
  factories,
  lastTick,
  onTick,
}: {
  factories: Factory[];
  lastTick: FactoryTickResult | null;
  onTick: (id: string) => void | Promise<void>;
}) {
  if (factories.length === 0) {
    return (
      <section className="card">
        <h2>Factory floor</h2>
        <p className="muted">No factories yet. Assemble one above.</p>
      </section>
    );
  }
  return (
    <section className="card">
      <h2>Factory floor</h2>
      <p className="description">
        One tick advances the working day: each machine&apos;s wear is
        accumulated, its product is added to the day&apos;s output, and the
        hand labour the factory would otherwise have required is recorded
        alongside the value the machinery contributed. The two together drive
        the §2 invariant <code>LabourDisplaced &gt; LabourEmbodiedInMachine</code>.
      </p>
      <table className="data-table">
        <thead>
          <tr>
            <th>Name</th>
            <th>Prime mover</th>
            <th>Machines</th>
            <th>Productive power</th>
            <th>Daily value transfer</th>
            <th>Ticks</th>
            <th />
          </tr>
        </thead>
        <tbody>
          {factories.map((f) => (
            <tr key={f.id}>
              <td>{f.name}</td>
              <td>
                {f.prime_mover.kind} ({f.prime_mover.horsepower} hp)
              </td>
              <td>{f.machines.length}</td>
              <td>{compactNumber(f.total_productive_power)}/d</td>
              <td>{poundsFromLabourMinutes(f.daily_value_transfer)}</td>
              <td>{f.tick_count}</td>
              <td>
                <button type="button" onClick={() => void onTick(f.id)}>
                  Tick →
                </button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
      {lastTick && <LastTickCard r={lastTick} />}
      {factories.map((f) => (
        <TickHistoryStrip key={f.id} factoryId={f.id} factoryName={f.name} tickCount={f.tick_count} />
      ))}
      {factories.length > 0 && (
        <RelativeSurplusBridge
          source="factory"
          sourceId={factories[0].id}
          factor={factories[0].productivity_factor ?? 1.0}
          sourceLabel={`Factory ${factories[0].name || factories[0].id.slice(0, 8)}`}
        />
      )}
    </section>
  );
}

function TickHistoryStrip({
  factoryId,
  factoryName,
  tickCount,
}: {
  factoryId: string;
  factoryName: string;
  tickCount: number;
}) {
  const [ticks, setTicks] = useState<FactoryTick[]>([]);
  const [err, setErr] = useState<string | null>(null);

  useEffect(() => {
    (async () => {
      try {
        const list = await api.listFactoryTicks(factoryId, 10);
        setTicks(list);
      } catch (e) {
        setErr(e instanceof Error ? e.message : String(e));
      }
    })();
  }, [factoryId, tickCount]);

  if (ticks.length === 0 && !err) return null;
  return (
    <div className="small" style={{ marginTop: "0.75rem" }}>
      <strong>{factoryName}</strong> — last {ticks.length} tick{ticks.length === 1 ? "" : "s"}
      {err && <span className="error" style={{ marginLeft: "0.5rem" }}>{err}</span>}
      {ticks.length > 0 && (
        <table className="data-table" style={{ marginTop: "0.25rem" }}>
          <thead>
            <tr>
              <th>#</th>
              <th>Units</th>
              <th>Value transferred</th>
              <th>Hand labour saved</th>
            </tr>
          </thead>
          <tbody>
            {ticks.map((t) => (
              <tr key={t.sequence}>
                <td>{t.sequence}</td>
                <td>{compactNumber(t.units_produced)}</td>
                <td>{poundsFromLabourMinutes(t.value_transferred)}</td>
                <td>{compactNumber(t.hand_labour_saved)} lab-min</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  );
}

function LastTickCard({ r }: { r: FactoryTickResult }) {
  return (
    <div className="small" style={{ marginTop: "0.75rem", padding: "0.5rem 0.75rem", border: "1px solid var(--ink-muted)" }}>
      <strong>{r.factory.name}</strong> — tick {r.factory.tick_count}
      <div>Units produced: {compactNumber(r.units_produced)}</div>
      <div>Value transferred: {poundsFromLabourMinutes(r.value_transferred)}</div>
      <div>Hand labour saved: {compactNumber(r.hand_labour_saved)} lab-min</div>
    </div>
  );
}

// ── Capital composition probe ─────────────────────────────────────────────────

function CapitalCompositionPanel() {
  // §6 fixture: V=£3000 (100×£30), C=£3000 → V=£1500 (50×£30), C=£4500;
  // total £6000 conserved.
  const [annualWagePounds, setAnnualWagePounds] = useState(30);
  const [workersBefore, setWorkersBefore] = useState(100);
  const [workersDisplaced, setWorkersDisplaced] = useState(50);
  const [machineCostPounds, setMachineCostPounds] = useState(1500);
  const [otherConstantPounds, setOtherConstantPounds] = useState(3000);

  const before = useMemo(() => {
    const v = workersBefore * annualWagePounds;
    const c = otherConstantPounds;
    return { variable: v, constant: c, total: v + c };
  }, [workersBefore, annualWagePounds, otherConstantPounds]);

  const after = useMemo(() => {
    const v = (workersBefore - workersDisplaced) * annualWagePounds;
    const c = otherConstantPounds + machineCostPounds;
    return { variable: v, constant: c, total: v + c };
  }, [workersBefore, workersDisplaced, annualWagePounds, otherConstantPounds, machineCostPounds]);

  const conserved =
    before.total === after.total - (machineCostPounds - workersDisplaced * annualWagePounds);
  const totalDiff = after.total - before.total;

  return (
    <section className="card">
      <h2>Capital composition (§6)</h2>
      <p className="description">
        Mechanisation converts variable into constant capital. In Marx&apos;s
        §6 example (£3,000 variable / 100 workmen at £30) a £1,500 machine
        displaces 50 men; variable falls to £1,500, constant rises by £1,500,
        total capital unchanged. Edit the inputs to test the conservation
        invariant against your own numbers.
      </p>
      <form className="form-grid" onSubmit={(e) => e.preventDefault()}>
        <label>
          <span>Annual wage (£/worker)</span>
          <input type="number" min={1} value={annualWagePounds} onChange={(e) => setAnnualWagePounds(Number(e.target.value))} />
        </label>
        <label>
          <span>Workers (before)</span>
          <input type="number" min={1} value={workersBefore} onChange={(e) => setWorkersBefore(Number(e.target.value))} />
        </label>
        <label>
          <span>Workers displaced</span>
          <input type="number" min={0} max={workersBefore} value={workersDisplaced} onChange={(e) => setWorkersDisplaced(Number(e.target.value))} />
        </label>
        <label>
          <span>Machine cost (£)</span>
          <input type="number" min={0} value={machineCostPounds} onChange={(e) => setMachineCostPounds(Number(e.target.value))} />
        </label>
        <label>
          <span>Other constant capital (£)</span>
          <input type="number" min={0} value={otherConstantPounds} onChange={(e) => setOtherConstantPounds(Number(e.target.value))} />
        </label>
      </form>
      <div className="small" style={{ marginTop: "0.75rem", padding: "0.5rem 0.75rem", border: "1px solid var(--ink-muted)" }}>
        <div>
          <strong>Before:</strong> V=£{before.variable.toLocaleString()},
          C=£{before.constant.toLocaleString()},
          Total=£{before.total.toLocaleString()},
          c/v={before.variable === 0 ? "—" : (before.constant / before.variable).toFixed(2)}
        </div>
        <div>
          <strong>After:</strong> V=£{after.variable.toLocaleString()},
          C=£{after.constant.toLocaleString()},
          Total=£{after.total.toLocaleString()},
          c/v={after.variable === 0 ? "—" : (after.constant / after.variable).toFixed(2)}
        </div>
        <div style={{ marginTop: "0.5rem" }}>
          {totalDiff === 0 ? (
            <span style={{ color: "#1f7a3a" }}>Total capital conserved ({conserved ? "machine cost = displaced wage bill" : "by chance"}).</span>
          ) : (
            <span className="error">Total moved by £{totalDiff.toLocaleString()} — outside the pure §6 substitution period.</span>
          )}
        </div>
      </div>
    </section>
  );
}
