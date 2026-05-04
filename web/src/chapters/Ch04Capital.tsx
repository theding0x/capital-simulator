import { useEffect, useState } from "react";
import type { FormEvent } from "react";
import { api } from "../api";
import type { Agent, CapitalCircuit } from "../types";

interface Ch04Props {
  onSharedChanged: () => void;
}

function penceToGBP(pence: number): string {
  return `£${(pence / 100).toFixed(2)}`;
}

export function Ch04Capital({ onSharedChanged: _onSharedChanged }: Ch04Props) {
  const [agents, setAgents] = useState<Agent[]>([]);
  const [error, setError] = useState<string | null>(null);

  async function refreshAgents() {
    try {
      const list = await api.listAgents();
      setAgents(list);
    } catch (e) {
      setError(e instanceof Error ? e.message : String(e));
    }
  }

  useEffect(() => {
    refreshAgents();
  }, []);

  const capitalists = agents.filter((a) => a.class === "capitalist");
  const workers = agents.filter((a) => a.class === "worker");
  const misers = agents.filter((a) => a.class === "miser");
  const owners = agents.filter((a) => a.class === "owner");

  return (
    <>
      <CreateAgentPanel onCreated={refreshAgents} />
      {error && <p className="error">{error}</p>}
      <AgentClassSection title="Capitalists" agents={capitalists} onChanged={refreshAgents} />
      <AgentClassSection title="Workers" agents={workers} onChanged={refreshAgents} />
      <AgentClassSection title="Misers" agents={misers} onChanged={refreshAgents} />
      <AgentClassSection title="Owners" agents={owners} onChanged={refreshAgents} />
    </>
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
          Name
          <input value={name} onChange={(e) => setName(e.target.value)} required />
        </label>
        <label>
          Class
          <select value={cls} onChange={(e) => setCls(e.target.value as typeof cls)}>
            <option value="capitalist">Capitalist</option>
            <option value="worker">Worker</option>
            <option value="miser">Miser</option>
            <option value="owner">Owner</option>
          </select>
        </label>
        <label>
          Initial balance (pence)
          <input
            type="number"
            value={balance}
            onChange={(e) => setBalance(Number(e.target.value))}
            min={0}
          />
        </label>
        <button type="submit">Create</button>
        {err && <span className="error">{err}</span>}
      </form>
    </section>
  );
}

function AgentClassSection({
  title,
  agents,
  onChanged,
}: {
  title: string;
  agents: Agent[];
  onChanged: () => void;
}) {
  if (agents.length === 0) return null;
  return (
    <section className="card">
      <h2>{title}</h2>
      <div className="item-list">
        {agents.map((a) => (
          <AgentCard key={a.id} agent={a} onChanged={onChanged} />
        ))}
      </div>
    </section>
  );
}

function AgentCard({ agent: a, onChanged }: { agent: Agent; onChanged: () => void }) {
  const [circuits, setCircuits] = useState<CapitalCircuit[]>([]);
  const [showCircuits, setShowCircuits] = useState(false);
  const [err, setErr] = useState<string | null>(null);

  async function loadCircuits() {
    try {
      const list = await api.listAgentCircuits(a.id);
      setCircuits(list);
      setShowCircuits(true);
    } catch (e) {
      setErr(e instanceof Error ? e.message : String(e));
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
        <span className="item-meta">{penceToGBP(a.money_balance)}</span>
        {a.hoarding && <span className="item-tag">hoarding</span>}
      </div>
      <div className="item-actions">
        <button onClick={loadCircuits} type="button">
          {showCircuits ? "Refresh circuits" : "Show circuits"}
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
          <CircuitTable circuits={circuits} />
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

function CircuitTable({ circuits }: { circuits: CapitalCircuit[] }) {
  if (circuits.length === 0) return <p className="muted small">No circuits yet.</p>;
  return (
    <table className="data-table">
      <thead>
        <tr>
          <th>Type</th>
          <th>M (advanced)</th>
          <th>C (commodity)</th>
          <th>M′ (returned)</th>
          <th>∆M (surplus)</th>
        </tr>
      </thead>
      <tbody>
        {circuits.map((c) => (
          <tr key={c.id}>
            <td>{c.circuit_type}</td>
            <td>{penceToGBP(c.m_advanced)}</td>
            <td className="monospace small">{c.commodity_id.slice(0, 8)}&hellip;</td>
            <td>{penceToGBP(c.m_returned)}</td>
            <td
              className={
                c.surplus_value > 0 ? "positive" : c.surplus_value < 0 ? "negative" : ""
              }
            >
              {c.surplus_value > 0 ? "+" : ""}
              {penceToGBP(c.surplus_value)}
            </td>
          </tr>
        ))}
      </tbody>
    </table>
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
        Commodity ID
        <input value={commodityId} onChange={(e) => setCommodityId(e.target.value)} required />
      </label>
      <label>
        M advanced (pence)
        <input
          type="number"
          value={mAdvanced}
          onChange={(e) => setMAdvanced(Number(e.target.value))}
          min={1}
        />
      </label>
      <label>
        M′ returned (pence)
        <input
          type="number"
          value={mReturned}
          onChange={(e) => setMReturned(Number(e.target.value))}
          min={0}
        />
      </label>
      {agentClass !== "worker" && agentClass !== "owner" && (
        <label>
          Circuit type
          <select
            value={circuitType}
            onChange={(e) => setCircuitType(e.target.value as typeof circuitType)}
          >
            <option value="M-C-M-prime">M—C—M′ (capital)</option>
            <option value="C-M-C">C—M—C (worker)</option>
          </select>
        </label>
      )}
      <button type="submit">Record circuit</button>
      {err && <span className="error">{err}</span>}
    </form>
  );
}
