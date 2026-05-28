import { useEffect, useMemo, useState } from "react";
import type { FormEvent } from "react";
import { api } from "../../api";
import type { Agent, CapitalCircuit } from "../../types";
import { usePounds } from "../../CurrencyContext";
import "./Ch04Capital.css";

interface Ch04Props {
  onSharedChanged: () => void;
}

export function Ch04Capital({ onSharedChanged: _onSharedChanged }: Ch04Props) {
  const [agents, setAgents] = useState<Agent[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);

  async function refreshAgents() {
    try {
      const list = await api.listAgents();
      setAgents(list);
      setError(null);
    } catch (e) {
      setError(e instanceof Error ? e.message : String(e));
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    refreshAgents();
  }, []);

  const capitalists = agents.filter((a) => a.class === "capitalist");
  const misers = agents.filter((a) => a.class === "miser");
  const workers = agents.filter((a) => a.class === "worker");
  const owners = agents.filter((a) => a.class === "owner");

  return (
    <>
      <InsightBox />
      <CreateAgentPanel onCreated={refreshAgents} />
      {error && <p className="v1-ch04-error">{error}</p>}
      {loading ? (
        <p className="v1-ch04-loading">Loading agents…</p>
      ) : (
        <>
          <AgentClassSection
            title="Active Capitalists"
            agents={capitalists}
            circuit="capital"
            emptyHint="No capitalist agents yet. Create one above to record an M—C—M′ circuit."
            onChanged={refreshAgents}
          />
          <AgentClassSection
            title="Active Misers"
            agents={misers}
            circuit="capital"
            emptyHint="No misers yet. Misers hoard money instead of advancing it as capital."
            onChanged={refreshAgents}
          />
          <CollapsibleSection
            title="Workers"
            agents={workers}
            circuit="simple"
            defaultOpen={false}
            emptyHint="No workers yet."
            onChanged={refreshAgents}
          />
          <CollapsibleSection
            title="Owners"
            agents={owners}
            circuit="simple"
            defaultOpen={false}
            emptyHint="No owners yet."
            onChanged={refreshAgents}
          />
        </>
      )}
      <Coda />
    </>
  );
}

function InsightBox() {
  return (
    <section className="v1-ch04-insight">
      <h2 className="v1-ch04-insight-h2">
        Capital is value that returns to itself with an increment.
      </h2>
      <div className="v1-ch04-insight-formulas">
        <div className="v1-ch04-formula v1-ch04-formula--surplus">
          <span className="v1-ch04-formula-tag">Capital</span>
          <span className="v1-ch04-formula-code">
            M—C—M
            <span className="v1-ch04-formula-prime">′</span>
          </span>
          <span className="v1-ch04-formula-gloss">
            Money advanced, commodity bought, money returned <em>with surplus</em>.
            The circuit closes only if M′ &gt; M.
          </span>
        </div>
        <div className="v1-ch04-formula v1-ch04-formula--simple">
          <span className="v1-ch04-formula-tag">Simple circulation</span>
          <span className="v1-ch04-formula-code">C—M—C</span>
          <span className="v1-ch04-formula-gloss">
            Commodity sold for money, money spent on a different commodity.
            Value is exchanged for use-value and the movement ends.
          </span>
        </div>
      </div>
      <p className="v1-ch04-insight-prose">
        The two circuits share the same elements but reverse their order — and
        with that reversal, the purpose changes. In C—M—C the start and end
        points are use-values; money is a fleeting mediator. In M—C—M′ money is
        both start and end; the commodity is the mediator; the <em>destination
        is exchange-value itself</em>, and only its quantitative expansion
        justifies the movement.
      </p>
    </section>
  );
}

function CreateAgentPanel({ onCreated }: { onCreated: () => void }) {
  const [name, setName] = useState("");
  const [cls, setCls] = useState<"capitalist" | "worker" | "miser" | "owner">("capitalist");
  const [balance, setBalance] = useState(10000);
  const [err, setErr] = useState<string | null>(null);

  async function submit(e: FormEvent) {
    e.preventDefault();
    setErr(null);
    try {
      await api.createAgent({ name, class: cls, money_balance: balance });
      setName("");
      onCreated();
    } catch (e) {
      setErr(e instanceof Error ? e.message : String(e));
    }
  }

  return (
    <section className="card">
      <h2>Create Agent</h2>
      <form className="form-grid" onSubmit={submit}>
        <label>
          <span>Name</span>
          <input value={name} onChange={(e) => setName(e.target.value)} required />
        </label>
        <label>
          <span>Class</span>
          <select value={cls} onChange={(e) => setCls(e.target.value as typeof cls)}>
            <option value="capitalist">Capitalist</option>
            <option value="worker">Worker</option>
            <option value="miser">Miser</option>
            <option value="owner">Owner</option>
          </select>
        </label>
        <label>
          <span>Initial balance (pence)</span>
          <input
            type="number"
            value={balance}
            onChange={(e) => setBalance(Number(e.target.value))}
            min={0}
          />
        </label>
        <div className="form-actions span2">
          <button type="submit" className="primary">Create</button>
          {err && <span className="error">{err}</span>}
        </div>
      </form>
    </section>
  );
}

type CircuitKind = "capital" | "simple";

function AgentClassSection({
  title,
  agents,
  circuit,
  emptyHint,
  onChanged,
}: {
  title: string;
  agents: Agent[];
  circuit: CircuitKind;
  emptyHint: string;
  onChanged: () => void;
}) {
  return (
    <section className="card">
      <h2>{title}</h2>
      {agents.length === 0 ? (
        <p className="v1-ch04-empty">{emptyHint}</p>
      ) : (
        <div className="item-list">
          {agents.map((a) => (
            <AgentCard key={a.id} agent={a} circuit={circuit} onChanged={onChanged} />
          ))}
        </div>
      )}
    </section>
  );
}

function CollapsibleSection({
  title,
  agents,
  circuit,
  defaultOpen,
  emptyHint,
  onChanged,
}: {
  title: string;
  agents: Agent[];
  circuit: CircuitKind;
  defaultOpen: boolean;
  emptyHint: string;
  onChanged: () => void;
}) {
  const [open, setOpen] = useState(defaultOpen);
  return (
    <section className="card">
      <button
        type="button"
        className="v1-ch04-section-toggle"
        aria-expanded={open}
        onClick={() => setOpen((v) => !v)}
      >
        <span className="v1-ch04-section-toggle-caret">{open ? "▾" : "▸"}</span>
        {title}
        <span className="v1-ch04-section-count">[{agents.length}]</span>
      </button>
      {open && (
        agents.length === 0 ? (
          <p className="v1-ch04-empty">{emptyHint}</p>
        ) : (
          <div className="item-list">
            {agents.map((a) => (
              <AgentCard key={a.id} agent={a} circuit={circuit} onChanged={onChanged} />
            ))}
          </div>
        )
      )}
    </section>
  );
}

function AgentCard({
  agent: a,
  circuit,
  onChanged,
}: {
  agent: Agent;
  circuit: CircuitKind;
  onChanged: () => void;
}) {
  const fmt = usePounds();
  const [circuits, setCircuits] = useState<CapitalCircuit[]>([]);
  const [showCircuits, setShowCircuits] = useState(false);
  const [err, setErr] = useState<string | null>(null);
  const [loadingCircuits, setLoadingCircuits] = useState(false);

  async function loadCircuits() {
    setErr(null);
    setLoadingCircuits(true);
    try {
      const list = await api.listAgentCircuits(a.id);
      setCircuits(list);
      setShowCircuits(true);
    } catch (e) {
      setErr(e instanceof Error ? e.message : String(e));
    } finally {
      setLoadingCircuits(false);
    }
  }

  async function hoard() {
    setErr(null);
    try {
      await api.hoardAgent(a.id);
      onChanged();
    } catch (e) {
      setErr(e instanceof Error ? e.message : String(e));
    }
  }

  return (
    <div className="item-card">
      <div className="item-header">
        <span className="item-name">{a.name}</span>
        <span className="item-meta">{fmt(a.money_balance)}</span>
        <span
          className={
            circuit === "capital"
              ? "v1-ch04-circuit-tag v1-ch04-circuit-tag--capital"
              : "v1-ch04-circuit-tag v1-ch04-circuit-tag--simple"
          }
          title={circuit === "capital" ? "M—C—M′ circuit" : "C—M—C circuit"}
        >
          {circuit === "capital" ? "M—C—M′" : "C—M—C"}
        </span>
        {a.hoarding && <span className="item-tag">hoarding</span>}
      </div>
      <div className="item-actions">
        <button onClick={loadCircuits} type="button" disabled={loadingCircuits}>
          {loadingCircuits ? "Loading…" : showCircuits ? "Refresh circuits" : "Show circuits"}
        </button>
        {a.class === "miser" && !a.hoarding && (
          <button onClick={hoard} type="button">
            Hoard
          </button>
        )}
      </div>
      {err && <span className="error">{err}</span>}
      {showCircuits && (
        <>
          <CircuitSparklines circuits={circuits} />
          {a.class !== "miser" && (
            <CreateCircuitForm
              agentId={a.id}
              agentClass={a.class}
              onCreated={() => {
                loadCircuits();
                onChanged();
              }}
            />
          )}
        </>
      )}
    </div>
  );
}

function CircuitSparklines({ circuits }: { circuits: CapitalCircuit[] }) {
  const fmt = usePounds();
  const maxValue = useMemo(() => {
    let m = 0;
    for (const c of circuits) {
      if (c.m_advanced > m) m = c.m_advanced;
      if (c.m_returned > m) m = c.m_returned;
    }
    return m;
  }, [circuits]);

  if (circuits.length === 0) {
    return <p className="v1-ch04-sparkline-empty">No circuits yet.</p>;
  }

  return (
    <div className="v1-ch04-sparkline">
      {circuits.map((c) => {
        const advancedPct = maxValue > 0 ? (c.m_advanced / maxValue) * 100 : 0;
        const returnedPct = maxValue > 0 ? (c.m_returned / maxValue) * 100 : 0;
        const surplus = c.surplus_value;
        const deltaSign = surplus > 0 ? "+" : "";
        const deltaCls =
          surplus > 0
            ? "v1-ch04-sparkline-value--gold"
            : surplus < 0
            ? "v1-ch04-sparkline-value--red"
            : "v1-ch04-sparkline-value--muted";
        return (
          <CircuitRow
            key={c.id}
            label={c.circuit_type === "M-C-M-prime" ? "M→M′" : "C→C"}
            advancedPct={advancedPct}
            returnedPct={returnedPct}
            advanced={c.m_advanced}
            returned={c.m_returned}
            surplus={surplus}
            deltaSign={deltaSign}
            deltaCls={deltaCls}
            fmt={fmt}
          />
        );
      })}
    </div>
  );
}

function CircuitRow({
  label,
  advancedPct,
  returnedPct,
  advanced,
  returned,
  surplus,
  deltaSign,
  deltaCls,
  fmt,
}: {
  label: string;
  advancedPct: number;
  returnedPct: number;
  advanced: number;
  returned: number;
  surplus: number;
  deltaSign: string;
  deltaCls: string;
  fmt: (n: number) => string;
}) {
  const sharedPct = Math.min(advancedPct, returnedPct);
  const overflowPct = Math.max(0, returnedPct - advancedPct);
  const shortfallPct = Math.max(0, advancedPct - returnedPct);
  return (
    <>
      <span className="v1-ch04-sparkline-label">{label}</span>
      <div className="v1-ch04-sparkline-track" title={`M ${fmt(advanced)} → M′ ${fmt(returned)}`}>
        <div
          className="v1-ch04-sparkline-bar v1-ch04-sparkline-bar--lead"
          style={{ width: `${sharedPct}%` }}
        />
        {overflowPct > 0 && (
          <div
            className="v1-ch04-sparkline-bar v1-ch04-sparkline-bar--surplus"
            style={{ left: `${sharedPct}%`, width: `${overflowPct}%` }}
          />
        )}
        {shortfallPct > 0 && (
          <div
            className="v1-ch04-sparkline-bar v1-ch04-sparkline-bar--shortfall"
            style={{ left: `${returnedPct}%`, width: `${shortfallPct}%` }}
          />
        )}
      </div>
      <span className={`v1-ch04-sparkline-value ${deltaCls}`}>
        {deltaSign}
        {fmt(surplus)}
      </span>
    </>
  );
}

function CreateCircuitForm({
  agentId,
  agentClass,
  onCreated,
}: {
  agentId: string;
  agentClass: "capitalist" | "worker" | "miser" | "owner";
  onCreated: () => void;
}) {
  const [commodityId, setCommodityId] = useState("");
  const [mAdvanced, setMAdvanced] = useState(10000);
  const [mReturned, setMReturned] = useState(11000);
  const [circuitType, setCircuitType] = useState<"C-M-C" | "M-C-M-prime">(
    agentClass === "worker" || agentClass === "owner" ? "C-M-C" : "M-C-M-prime"
  );
  const [err, setErr] = useState<string | null>(null);

  async function submit(e: FormEvent) {
    e.preventDefault();
    setErr(null);
    try {
      await api.createAgentCircuit(agentId, {
        m_advanced: mAdvanced,
        commodity_id: commodityId,
        m_returned: mReturned,
        circuit_type: circuitType,
      });
      onCreated();
    } catch (e) {
      setErr(e instanceof Error ? e.message : String(e));
    }
  }

  return (
    <form className="form-grid" onSubmit={submit}>
      <label>
        <span>Commodity ID</span>
        <input value={commodityId} onChange={(e) => setCommodityId(e.target.value)} required />
      </label>
      <label>
        <span>M advanced (pence)</span>
        <input
          type="number"
          value={mAdvanced}
          onChange={(e) => setMAdvanced(Number(e.target.value))}
          min={1}
        />
      </label>
      <label>
        <span>M′ returned (pence)</span>
        <input
          type="number"
          value={mReturned}
          onChange={(e) => setMReturned(Number(e.target.value))}
          min={0}
        />
      </label>
      {agentClass !== "worker" && agentClass !== "owner" && (
        <label>
          <span>Circuit type</span>
          <select
            value={circuitType}
            onChange={(e) => setCircuitType(e.target.value as typeof circuitType)}
          >
            <option value="M-C-M-prime">M—C—M′ (capital)</option>
            <option value="C-M-C">C—M—C (worker)</option>
          </select>
        </label>
      )}
      <div className="form-actions span2">
        <button type="submit" className="primary">Record circuit</button>
        {err && <span className="error">{err}</span>}
      </div>
    </form>
  );
}

function Coda() {
  return (
    <aside className="v1-ch04-coda">
      <p className="v1-ch04-coda-quote">
        “In simple circulation… in the circulation M—C—M′, both… have one and
        the same destination: exchange-value.”
        <span className="v1-ch04-coda-cite">
          — Marx, Capital Vol. I, Ch. 4
        </span>
      </p>
    </aside>
  );
}
