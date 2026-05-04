function emptyDraft() {
  return {
    name: "", unit: "", use_value_desc: "", concrete_labour: "", snlt: 0,
  };
}

function CommodityForm({ onCreate }) {
  const [d, setD] = React.useState(emptyDraft());
  const [busy, setBusy] = React.useState(false);

  function submit(e) {
    e.preventDefault();
    if (!d.name || !d.unit || !d.snlt) return;
    setBusy(true);
    setTimeout(() => {
      onCreate({
        id: "c_" + Math.random().toString(36).slice(2, 10),
        name: d.name,
        use_value: { description: d.use_value_desc || `${d.name} for use`, unit: d.unit },
        concrete_labour: { kind: d.concrete_labour || "labour" },
        snlt_per_unit: Number(d.snlt) || 0,
      });
      setD(emptyDraft());
      setBusy(false);
    }, 250);
  }

  return (
    <Card title="Register a commodity">
      <form className="form-grid" onSubmit={submit}>
        <label><span>Name</span>
          <input value={d.name} onChange={(e) => setD({ ...d, name: e.target.value })} placeholder="linen" required />
        </label>
        <label><span>Unit of measure</span>
          <input value={d.unit} onChange={(e) => setD({ ...d, unit: e.target.value })} placeholder="yards" required />
        </label>
        <label className="span2"><span>Use-value description</span>
          <input value={d.use_value_desc} onChange={(e) => setD({ ...d, use_value_desc: e.target.value })} placeholder="linen for clothing" />
        </label>
        <label><span>Concrete labour kind</span>
          <input value={d.concrete_labour} onChange={(e) => setD({ ...d, concrete_labour: e.target.value })} placeholder="weaving" />
        </label>
        <label><span>SNLT (minutes per unit)</span>
          <input type="number" min={0} value={d.snlt} onChange={(e) => setD({ ...d, snlt: Number(e.target.value) })} placeholder="30" required />
        </label>
        <div className="span2 form-actions">
          <button type="submit" className="primary" disabled={busy}>{busy ? "Saving…" : "Register"}</button>
        </div>
      </form>
    </Card>
  );
}

function fmtMinutes(m) { return `${m.toLocaleString()} min`; }
function fmtQty(q) {
  return Number.isInteger(q) ? String(q) : q.toFixed(2).replace(/0+$/, "").replace(/\.$/, "");
}

function CommodityTable({ commodities, onDelete, onReveal, revealed }) {
  if (commodities.length === 0) {
    return (
      <Card title="Commodities">
        <div className="empty-state">
          <p>No commodities registered yet.</p>
          <p className="small">Use the form above to add the first one.</p>
        </div>
      </Card>
    );
  }
  return (
    <Card title="Commodities">
      <table className="commodity-table">
        <thead>
          <tr><th>Name</th><th>Use-value</th><th>Concrete labour</th><th>SNLT / unit</th><th></th></tr>
        </thead>
        <tbody>
          {commodities.map((c) => (
            <React.Fragment key={c.id}>
              <tr>
                <td className="name-cell">{c.name}</td>
                <td>
                  <div>{c.use_value.description}</div>
                  <div className="muted small">per {c.use_value.unit}</div>
                </td>
                <td className="muted">{c.concrete_labour.kind}</td>
                <td className="snlt-cell">{fmtMinutes(c.snlt_per_unit)}</td>
                <td className="actions">
                  <button className="secondary">Edit</button>
                  <button className="reveal-btn" onClick={() => onReveal(c.id)}>
                    {revealed === c.id ? "Hide" : "Reveal"}
                  </button>
                  <button className="danger" onClick={() => onDelete(c.id)}>Delete</button>
                </td>
              </tr>
              {revealed === c.id && (
                <tr>
                  <td colSpan={5} className="reveal">
                    <RevealPanel commodity={c} all={commodities} />
                  </td>
                </tr>
              )}
            </React.Fragment>
          ))}
        </tbody>
      </table>
    </Card>
  );
}

function RevealPanel({ commodity: c, all }) {
  const others = all.filter((x) => x.id !== c.id);
  return (
    <div className="reveal-panel">
      <h3>Social relations of {c.name}</h3>
      <p className="reveal-note">
        The exchange of commodities is, behind the back of the producers, an exchange of human labour-times taking the form of a thing.
      </p>
      <p className="labour-statement">
        <strong>{fmtMinutes(c.snlt_per_unit)}</strong> of abstract human labour are congealed in one {c.use_value.unit} of {c.name}, produced by {c.concrete_labour.kind}.
      </p>
      {others.length > 0 && (
        <ul className="reveal-relations">
          {others.map((o) => {
            const ratio = c.snlt_per_unit / o.snlt_per_unit;
            return (
              <li key={o.id}>
                1&nbsp;{c.use_value.unit} {c.name} ≡ <strong>{fmtQty(ratio)}</strong> {o.use_value.unit} {o.name} — both expressing <strong>{fmtMinutes(c.snlt_per_unit)}</strong> of abstract human labour ({c.concrete_labour.kind} ↔ {o.concrete_labour.kind})
              </li>
            );
          })}
        </ul>
      )}
    </div>
  );
}

function ExchangeRatio({ commodities }) {
  const [baseId, setBaseId] = React.useState("");
  const [quoteId, setQuoteId] = React.useState("");
  const [qty, setQty] = React.useState(1);
  const [result, setResult] = React.useState(null);

  if (commodities.length < 2) return null;
  const can = baseId && quoteId && baseId !== quoteId && qty > 0;

  function compute() {
    const base = commodities.find((c) => c.id === baseId);
    const quote = commodities.find((c) => c.id === quoteId);
    const common = base.snlt_per_unit * qty;
    const quoteQty = common / quote.snlt_per_unit;
    setResult({ base, quote, base_qty: qty, quote_qty: quoteQty, common_value: common });
  }

  return (
    <Card
      title="Exchange ratio"
      description='"1 quarter corn = x cwt iron" — pick two commodities to derive the equation of values the simple form of value expresses.'
    >
      <div className="form-grid">
        <label><span>Base commodity</span>
          <select value={baseId} onChange={(e) => setBaseId(e.target.value)}>
            <option value="">—</option>
            {commodities.map((c) => <option key={c.id} value={c.id}>{c.name}</option>)}
          </select>
        </label>
        <label><span>Quantity</span>
          <input type="number" min={0} step="any" value={qty} onChange={(e) => setQty(Number(e.target.value))} />
        </label>
        <label><span>Quote commodity</span>
          <select value={quoteId} onChange={(e) => setQuoteId(e.target.value)}>
            <option value="">—</option>
            {commodities.map((c) => <option key={c.id} value={c.id}>{c.name}</option>)}
          </select>
        </label>
        <div className="span2 form-actions">
          <button className="primary" onClick={compute} disabled={!can}>Compute ratio</button>
        </div>
      </div>
      {result && (
        <div className="ratio-result">
          <div className="ratio-equation">
            <span><span className="eq-qty">{fmtQty(result.base_qty)}</span> <span className="eq-unit">{result.base.use_value.unit} {result.base.name}</span></span>
            <span className="eq-sign">=</span>
            <span><span className="eq-qty">{fmtQty(result.quote_qty)}</span> <span className="eq-unit">{result.quote.use_value.unit} {result.quote.name}</span></span>
          </div>
          <p className="ratio-common">common value: {fmtMinutes(result.common_value)} of abstract human labour</p>
        </div>
      )}
    </Card>
  );
}

Object.assign(window, { CommodityForm, CommodityTable, RevealPanel, ExchangeRatio, fmtMinutes, fmtQty });
