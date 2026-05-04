// Ch. 03 — Money. Mirrors web/src/chapters/Ch03Money.tsx with mock state.
// Sub-panels: CirculationPanel, CircuitPanel, HoardPanel, PaymentObligationsPanel, WorldMoneyPanel.

function CirculationPanel() {
  const [sum, setSum] = React.useState(120);
  const [velocity, setVelocity] = React.useState(4);
  const [result, setResult] = React.useState(null);
  function compute(e) {
    e.preventDefault();
    if (!velocity) return;
    setResult(Number((sum / velocity).toFixed(2)));
  }
  return (
    <Card
      title="Quantity of circulating medium"
      description='"The quantity of money as circulating medium = the sum of the prices of the commodities / the number of moves made by coins." (Ch. 3 §2B)'
    >
      <form className="form-grid" onSubmit={compute}>
        <label><span>Sum of prices</span>
          <input type="number" min={1} value={sum} onChange={(e) => setSum(Number(e.target.value))} required />
        </label>
        <label><span>Velocity (moves per coin)</span>
          <input type="number" min={1} value={velocity} onChange={(e) => setVelocity(Number(e.target.value))} required />
        </label>
        <div className="span2 form-actions">
          <button type="submit" className="primary">Compute</button>
        </div>
      </form>
      {result !== null && (
        <p className="reveal-note" style={{ marginTop: "0.75rem" }}>
          Money required: <strong>{result}</strong> units (M = {sum} / {velocity})
        </p>
      )}
    </Card>
  );
}

function CircuitPanel({ owners, circuits, setCircuits }) {
  const empty = (kind) => ({ kind, commodity_id: "", money_id: "", owner_id: "", value: 0 });
  const [sale, setSale] = React.useState(empty("sale"));
  const [purchase, setPurchase] = React.useState(empty("purchase"));
  function submit(e) {
    e.preventDefault();
    setCircuits((xs) => [...xs, { id: "ci_" + Math.random().toString(36).slice(2, 8), sale_leg: sale, purchase_leg: purchase }]);
    setSale(empty("sale")); setPurchase(empty("purchase"));
  }
  return (
    <Card
      title="C-M-C Circuits"
      description='"The act C-M-C — sell in order to buy — a commodity is converted into money, then the money is re-converted into a commodity." (Ch. 3 §2A)'
    >
      <form className="form-grid" onSubmit={submit}>
        <h3 className="span2" style={{ marginBottom: 0 }}>Sale leg (C → M)</h3>
        <label><span>Owner</span>
          <select value={sale.owner_id} onChange={(e) => setSale({ ...sale, owner_id: e.target.value })} required>
            <option value="">—</option>{owners.map((o) => <option key={o.id} value={o.id}>{o.name}</option>)}
          </select>
        </label>
        <label><span>Commodity ID sold</span>
          <input value={sale.commodity_id} onChange={(e) => setSale({ ...sale, commodity_id: e.target.value })} placeholder="commodity hex ID" required />
        </label>
        <label><span>Money commodity ID</span>
          <input value={sale.money_id} onChange={(e) => setSale({ ...sale, money_id: e.target.value })} placeholder="money hex ID" required />
        </label>
        <label><span>Value (labour-minutes)</span>
          <input type="number" min={1} value={sale.value} onChange={(e) => setSale({ ...sale, value: Number(e.target.value) })} required />
        </label>
        <h3 className="span2" style={{ marginBottom: 0 }}>Purchase leg (M → C′)</h3>
        <label><span>Owner</span>
          <select value={purchase.owner_id} onChange={(e) => setPurchase({ ...purchase, owner_id: e.target.value })} required>
            <option value="">—</option>{owners.map((o) => <option key={o.id} value={o.id}>{o.name}</option>)}
          </select>
        </label>
        <label><span>Commodity ID bought</span>
          <input value={purchase.commodity_id} onChange={(e) => setPurchase({ ...purchase, commodity_id: e.target.value })} placeholder="commodity hex ID" required />
        </label>
        <label><span>Money commodity ID</span>
          <input value={purchase.money_id} onChange={(e) => setPurchase({ ...purchase, money_id: e.target.value })} placeholder="money hex ID" required />
        </label>
        <label><span>Value (labour-minutes)</span>
          <input type="number" min={1} value={purchase.value} onChange={(e) => setPurchase({ ...purchase, value: Number(e.target.value) })} required />
        </label>
        <div className="span2 form-actions">
          <button type="submit" className="primary">Record circuit</button>
        </div>
      </form>
      {circuits.length > 0 && (
        <table className="commodity-table" style={{ marginTop: "1rem" }}>
          <thead><tr><th>ID</th><th>Sells</th><th>Buys</th><th>Value</th></tr></thead>
          <tbody>
            {circuits.map((c) => (
              <tr key={c.id}>
                <td className="muted small">{c.id}…</td>
                <td className="muted small">{(c.sale_leg.commodity_id || "").slice(0, 8)}</td>
                <td className="muted small">{(c.purchase_leg.commodity_id || "").slice(0, 8)}</td>
                <td className="snlt-cell">{fmtMinutes(c.sale_leg.value)}</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </Card>
  );
}

function HoardPanel({ owners, hoards, setHoards }) {
  const [ownerID, setOwnerID] = React.useState("");
  const [amount, setAmount] = React.useState(0);
  if (owners.length === 0) return null;
  const ownerById = new Map(owners.map((o) => [o.id, o]));
  function submit(e) {
    e.preventDefault();
    if (!ownerID || !amount) return;
    setHoards((xs) => [...xs, { id: "h_" + Math.random().toString(36).slice(2, 8), owner_id: ownerID, amount }]);
    setAmount(0);
  }
  return (
    <Card
      title="Hoards"
      description='"To accumulate money it is necessary to sell, to keep on selling, and to abstain from buying." (Ch. 3 §3A)'
    >
      <form className="form-grid" onSubmit={submit}>
        <label><span>Owner</span>
          <select value={ownerID} onChange={(e) => setOwnerID(e.target.value)} required>
            <option value="">—</option>{owners.map((o) => <option key={o.id} value={o.id}>{o.name}</option>)}
          </select>
        </label>
        <label><span>Amount (units)</span>
          <input type="number" min={1} value={amount} onChange={(e) => setAmount(Number(e.target.value))} required />
        </label>
        <div className="span2 form-actions">
          <button type="submit" className="primary">Hoard</button>
        </div>
      </form>
      {hoards.length > 0 && (
        <table className="commodity-table" style={{ marginTop: "1rem" }}>
          <thead><tr><th>Owner</th><th>Amount</th></tr></thead>
          <tbody>
            {hoards.map((h) => (
              <tr key={h.id}>
                <td>{ownerById.get(h.owner_id)?.name ?? h.owner_id.slice(0, 8)}</td>
                <td className="snlt-cell">{h.amount}</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </Card>
  );
}

function PaymentObligationsPanel({ owners, obligations, setObligations }) {
  const [d, setD] = React.useState({ creditor: "", debtor: "", amount: 0, due: "" });
  if (owners.length === 0) return null;
  const ownerById = new Map(owners.map((o) => [o.id, o]));
  function submit(e) {
    e.preventDefault();
    setObligations((xs) => [...xs, {
      id: "po_" + Math.random().toString(36).slice(2, 8),
      creditor_id: d.creditor, debtor_id: d.debtor, amount: d.amount,
      due_at: new Date(d.due).toISOString(), paid_at: null,
    }]);
    setD({ creditor: "", debtor: "", amount: 0, due: "" });
  }
  function settle(id) { setObligations((xs) => xs.map((o) => o.id === id ? { ...o, paid_at: new Date().toISOString() } : o)); }
  return (
    <Card
      title="Payment obligations"
      description='"The vendor becomes a creditor, the purchaser a debtor. Since money now functions as a means of payment after the commodity has changed hands, it becomes a debtor obligation." (Ch. 3 §3B)'
    >
      <form className="form-grid" onSubmit={submit}>
        <label><span>Creditor (vendor)</span>
          <select value={d.creditor} onChange={(e) => setD({ ...d, creditor: e.target.value })} required>
            <option value="">—</option>{owners.map((o) => <option key={o.id} value={o.id}>{o.name}</option>)}
          </select>
        </label>
        <label><span>Debtor (purchaser)</span>
          <select value={d.debtor} onChange={(e) => setD({ ...d, debtor: e.target.value })} required>
            <option value="">—</option>{owners.map((o) => <option key={o.id} value={o.id}>{o.name}</option>)}
          </select>
        </label>
        <label><span>Amount (units)</span>
          <input type="number" min={1} value={d.amount} onChange={(e) => setD({ ...d, amount: Number(e.target.value) })} required />
        </label>
        <label><span>Due date</span>
          <input type="date" value={d.due} onChange={(e) => setD({ ...d, due: e.target.value })} required />
        </label>
        <div className="span2 form-actions">
          <button type="submit" className="primary">Create obligation</button>
        </div>
      </form>
      {obligations.length > 0 && (
        <table className="commodity-table" style={{ marginTop: "1rem" }}>
          <thead><tr><th>Creditor</th><th>Debtor</th><th>Amount</th><th>Due</th><th>Status</th><th></th></tr></thead>
          <tbody>
            {obligations.map((ob) => (
              <tr key={ob.id}>
                <td>{ownerById.get(ob.creditor_id)?.name ?? ob.creditor_id.slice(0, 8)}</td>
                <td>{ownerById.get(ob.debtor_id)?.name ?? ob.debtor_id.slice(0, 8)}</td>
                <td className="snlt-cell">{ob.amount}</td>
                <td className="muted small">{ob.due_at.slice(0, 10)}</td>
                <td>{ob.paid_at ? <span className="muted">settled</span> : <span className="small">outstanding</span>}</td>
                <td className="actions">{!ob.paid_at && <button className="secondary" onClick={() => settle(ob.id)}>Settle</button>}</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </Card>
  );
}

function WorldMoneyPanel({ owners, transfers, setTransfers }) {
  const [d, setD] = React.useState({ sender: "", receiver: "", gold: 0 });
  if (owners.length === 0) return null;
  const ownerById = new Map(owners.map((o) => [o.id, o]));
  function submit(e) {
    e.preventDefault();
    setTransfers((xs) => [...xs, { id: "wm_" + Math.random().toString(36).slice(2, 8), sender_id: d.sender, receiver_id: d.receiver, gold_mg: d.gold }]);
    setD({ sender: "", receiver: "", gold: 0 });
  }
  return (
    <Card
      title="World money transfers"
      description='"On the world market, money strips off its national uniform and appears in its original form as bullion." (Ch. 3 §3C)'
    >
      <form className="form-grid" onSubmit={submit}>
        <label><span>Sender</span>
          <select value={d.sender} onChange={(e) => setD({ ...d, sender: e.target.value })} required>
            <option value="">—</option>{owners.map((o) => <option key={o.id} value={o.id}>{o.name}</option>)}
          </select>
        </label>
        <label><span>Receiver</span>
          <select value={d.receiver} onChange={(e) => setD({ ...d, receiver: e.target.value })} required>
            <option value="">—</option>{owners.map((o) => <option key={o.id} value={o.id}>{o.name}</option>)}
          </select>
        </label>
        <label><span>Gold weight (mg)</span>
          <input type="number" min={1} value={d.gold} onChange={(e) => setD({ ...d, gold: Number(e.target.value) })} required />
        </label>
        <div className="span2 form-actions">
          <button type="submit" className="primary">Transfer bullion</button>
        </div>
      </form>
      {transfers.length > 0 && (
        <table className="commodity-table" style={{ marginTop: "1rem" }}>
          <thead><tr><th>Sender</th><th>Receiver</th><th>Gold (mg)</th></tr></thead>
          <tbody>
            {transfers.map((t) => (
              <tr key={t.id}>
                <td>{ownerById.get(t.sender_id)?.name}</td>
                <td>{ownerById.get(t.receiver_id)?.name}</td>
                <td className="snlt-cell">{t.gold_mg}</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </Card>
  );
}

function Ch03MoneyView({ ch02, ch03 }) {
  return (
    <>
      <CirculationPanel />
      <CircuitPanel owners={ch02.owners} circuits={ch03.circuits} setCircuits={(fn) => ch03.set("circuits", fn)} />
      <HoardPanel owners={ch02.owners} hoards={ch03.hoards} setHoards={(fn) => ch03.set("hoards", fn)} />
      <PaymentObligationsPanel owners={ch02.owners} obligations={ch03.obligations} setObligations={(fn) => ch03.set("obligations", fn)} />
      <WorldMoneyPanel owners={ch02.owners} transfers={ch03.transfers} setTransfers={(fn) => ch03.set("transfers", fn)} />
    </>
  );
}

Object.assign(window, { CirculationPanel, CircuitPanel, HoardPanel, PaymentObligationsPanel, WorldMoneyPanel, Ch03MoneyView });
