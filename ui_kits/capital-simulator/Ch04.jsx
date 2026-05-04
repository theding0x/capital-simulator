// Ch. 04 — The General Formula for Capital. Mirrors web/src/chapters/Ch04Capital.tsx.
// Note: Ch. 04 in the source uses lighter form styling (no `.primary` button, plain `<label>` text)
// — we preserve that contrast deliberately.

function penceToGBP(pence) { return `£${(pence / 100).toFixed(2)}`; }

function CreateAgentPanel({ onCreate }) {
  const [name, setName] = React.useState("");
  const [cls, setCls] = React.useState("capitalist");
  const [balance, setBalance] = React.useState(10000);
  function submit(e) {
    e.preventDefault();
    if (!name.trim()) return;
    onCreate({
      id: "ag_" + Math.random().toString(36).slice(2, 8),
      name: name.trim(), class: cls, money_balance: balance, hoarding: false,
    });
    setName("");
  }
  return (
    <Card title="Create Agent">
      <form className="form-grid" onSubmit={submit}>
        <label><span>Name</span>
          <input value={name} onChange={(e) => setName(e.target.value)} required />
        </label>
        <label><span>Class</span>
          <select value={cls} onChange={(e) => setCls(e.target.value)}>
            <option value="capitalist">Capitalist</option>
            <option value="worker">Worker</option>
            <option value="miser">Miser</option>
          </select>
        </label>
        <label><span>Initial balance (pence)</span>
          <input type="number" min={0} value={balance} onChange={(e) => setBalance(Number(e.target.value))} />
        </label>
        <div className="span2 form-actions">
          <button type="submit">Create</button>
        </div>
      </form>
    </Card>
  );
}

function AgentClassSection({ title, agents, circuits, setCircuits, onHoard }) {
  if (agents.length === 0) return null;
  return (
    <Card title={title}>
      <div style={{ display: "flex", flexDirection: "column", gap: "1rem" }}>
        {agents.map((a) => (
          <AgentCard key={a.id} agent={a} circuits={circuits.filter((c) => c.agent_id === a.id)} setCircuits={setCircuits} onHoard={onHoard} />
        ))}
      </div>
    </Card>
  );
}

function AgentCard({ agent: a, circuits, setCircuits, onHoard }) {
  const [show, setShow] = React.useState(false);
  return (
    <div style={{
      padding: "1rem 1.1rem",
      border: "1px solid var(--border-subtle)",
      borderRadius: "var(--radius-sm)",
      background: "rgba(0,0,0,0.18)",
    }}>
      <div style={{ display: "flex", alignItems: "baseline", gap: "1rem", flexWrap: "wrap" }}>
        <span style={{ fontFamily: "var(--font-display)", fontSize: "1.05rem", fontWeight: 600 }}>{a.name}</span>
        <span className="snlt-cell" style={{ fontFamily: "var(--font-mono)" }}>{penceToGBP(a.money_balance)}</span>
        {a.hoarding && (
          <span style={{
            fontFamily: "var(--font-mono)", fontSize: "0.62rem", letterSpacing: "0.18em",
            textTransform: "uppercase", color: "var(--gold-bright)",
            padding: "0.15rem 0.5rem", border: "1px solid var(--gold-border)", borderRadius: "var(--radius-sm)",
            background: "var(--gold-bg)",
          }}>hoarding</span>
        )}
      </div>
      <div className="form-actions" style={{ marginTop: "0.6rem" }}>
        <button onClick={() => setShow((s) => !s)}>{show ? "Hide circuits" : "Show circuits"}</button>
        {a.class === "miser" && !a.hoarding && <button onClick={() => onHoard(a.id)}>Hoard</button>}
      </div>
      {show && (
        <div style={{ marginTop: "1rem" }}>
          <CircuitTable circuits={circuits} />
          {a.class !== "miser" && (
            <CreateCircuitForm agent={a} onCreate={(c) => setCircuits((xs) => [...xs, { ...c, id: "cc_" + Math.random().toString(36).slice(2, 8), agent_id: a.id, surplus_value: c.m_returned - c.m_advanced }])} />
          )}
        </div>
      )}
    </div>
  );
}

function CircuitTable({ circuits }) {
  if (circuits.length === 0) return <p className="muted small">No circuits yet.</p>;
  return (
    <table className="commodity-table">
      <thead><tr><th>Type</th><th>M (advanced)</th><th>C (commodity)</th><th>M′ (returned)</th><th>∆M (surplus)</th></tr></thead>
      <tbody>
        {circuits.map((c) => (
          <tr key={c.id}>
            <td className="muted small">{c.circuit_type}</td>
            <td className="snlt-cell">{penceToGBP(c.m_advanced)}</td>
            <td className="muted small">{c.commodity_id.slice(0, 8)}…</td>
            <td className="snlt-cell">{penceToGBP(c.m_returned)}</td>
            <td className="snlt-cell" style={{ color: c.surplus_value > 0 ? "#9ec07a" : c.surplus_value < 0 ? "#d44040" : "var(--ink-muted)" }}>
              {c.surplus_value > 0 ? "+" : ""}{penceToGBP(c.surplus_value)}
            </td>
          </tr>
        ))}
      </tbody>
    </table>
  );
}

function CreateCircuitForm({ agent, onCreate }) {
  const [commodityId, setCommodityId] = React.useState("");
  const [mAdvanced, setMAdvanced] = React.useState(10000);
  const [mReturned, setMReturned] = React.useState(11000);
  const [circuitType, setCircuitType] = React.useState(agent.class === "worker" ? "C-M-C" : "M-C-M-prime");
  function submit(e) {
    e.preventDefault();
    if (!commodityId.trim()) return;
    onCreate({ commodity_id: commodityId.trim(), m_advanced: mAdvanced, m_returned: mReturned, circuit_type: circuitType });
    setCommodityId("");
  }
  return (
    <form className="form-grid" onSubmit={submit} style={{ marginTop: "1rem" }}>
      <label><span>Commodity ID</span>
        <input value={commodityId} onChange={(e) => setCommodityId(e.target.value)} placeholder="commodity hex ID" required />
      </label>
      <label><span>M advanced (pence)</span>
        <input type="number" min={1} value={mAdvanced} onChange={(e) => setMAdvanced(Number(e.target.value))} />
      </label>
      <label><span>M′ returned (pence)</span>
        <input type="number" min={0} value={mReturned} onChange={(e) => setMReturned(Number(e.target.value))} />
      </label>
      {agent.class !== "worker" && (
        <label><span>Circuit type</span>
          <select value={circuitType} onChange={(e) => setCircuitType(e.target.value)}>
            <option value="M-C-M-prime">M—C—M′ (capital)</option>
            <option value="C-M-C">C—M—C (worker)</option>
          </select>
        </label>
      )}
      <div className="span2 form-actions">
        <button type="submit">Record circuit</button>
      </div>
    </form>
  );
}

function Ch04CapitalView({ ch04 }) {
  const { agents, circuits } = ch04;
  const setAgents = (fn) => ch04.set("agents", fn);
  const setCircuits = (fn) => ch04.set("circuits", fn);
  function onHoard(id) { setAgents((xs) => xs.map((a) => a.id === id ? { ...a, hoarding: true } : a)); }
  return (
    <>
      <CreateAgentPanel onCreate={(a) => setAgents((xs) => [...xs, a])} />
      <AgentClassSection title="Capitalists" agents={agents.filter((a) => a.class === "capitalist")} circuits={circuits} setCircuits={setCircuits} onHoard={onHoard} />
      <AgentClassSection title="Workers"     agents={agents.filter((a) => a.class === "worker")}     circuits={circuits} setCircuits={setCircuits} onHoard={onHoard} />
      <AgentClassSection title="Misers"      agents={agents.filter((a) => a.class === "miser")}      circuits={circuits} setCircuits={setCircuits} onHoard={onHoard} />
    </>
  );
}

Object.assign(window, { CreateAgentPanel, AgentClassSection, AgentCard, CircuitTable, CreateCircuitForm, Ch04CapitalView, penceToGBP });
