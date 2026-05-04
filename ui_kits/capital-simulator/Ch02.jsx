// Ch. 02 — Exchange. Mirrors web/src/chapters/Ch02Exchange.tsx with mock state.
// Sub-panels: OwnersPanel, OffersPanel, ExchangesPanel, MoneyPanel.

function OwnersPanel({ owners, setOwners }) {
  const [name, setName] = React.useState("");
  function submit(e) {
    e.preventDefault();
    if (!name.trim()) return;
    setOwners((xs) => [...xs, { id: "o_" + Math.random().toString(36).slice(2, 10), name: name.trim() }]);
    setName("");
  }
  return (
    <Card
      title="Owners"
      description='"Commodities cannot go to market of their own account. We must have recourse to their guardians, who are also their owners." (Ch. 2)'
    >
      <form className="form-grid" onSubmit={submit}>
        <label className="span2"><span>Owner name</span>
          <input value={name} onChange={(e) => setName(e.target.value)} placeholder="weaver" required />
        </label>
        <div className="span2 form-actions">
          <button type="submit" className="primary">Register owner</button>
        </div>
      </form>
      {owners.length > 0 && (
        <table className="commodity-table" style={{ marginTop: "1rem" }}>
          <thead><tr><th>Name</th><th>ID</th></tr></thead>
          <tbody>
            {owners.map((o) => (
              <tr key={o.id}>
                <td>{o.name}</td>
                <td className="muted small">{o.id.slice(0, 10)}…</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </Card>
  );
}

function OffersPanel({ owners, commodities, offers, setOffers }) {
  const [ownerID, setOwnerID] = React.useState("");
  const [commodityID, setCommodityID] = React.useState("");
  const [qty, setQty] = React.useState(1);
  const [seeksKind, setSeeksKind] = React.useState("");
  if (owners.length === 0 || commodities.length === 0) return null;

  const ownerById = new Map(owners.map((o) => [o.id, o]));
  const commodityById = new Map(commodities.map((c) => [c.id, c]));

  function submit(e) {
    e.preventDefault();
    if (!ownerID || !commodityID || !qty || !seeksKind) return;
    setOffers((xs) => [...xs, {
      id: "off_" + Math.random().toString(36).slice(2, 10),
      owner_id: ownerID, commodity_id: commodityID, quantity: qty, seeks_kind: seeksKind,
    }]);
    setQty(1); setSeeksKind("");
  }

  return (
    <Card
      title="Market offers"
      description="Each owner brings their commodity to market as a non-use-value for themselves, seeking another's use-value in return."
    >
      <form className="form-grid" onSubmit={submit}>
        <label><span>Owner</span>
          <select value={ownerID} onChange={(e) => setOwnerID(e.target.value)} required>
            <option value="">—</option>
            {owners.map((o) => <option key={o.id} value={o.id}>{o.name}</option>)}
          </select>
        </label>
        <label><span>Commodity offered</span>
          <select value={commodityID} onChange={(e) => setCommodityID(e.target.value)} required>
            <option value="">—</option>
            {commodities.map((c) => <option key={c.id} value={c.id}>{c.name}</option>)}
          </select>
        </label>
        <label><span>Quantity</span>
          <input type="number" min={0} step="any" value={qty} onChange={(e) => setQty(Number(e.target.value))} required />
        </label>
        <label><span>Seeks (kind)</span>
          <input value={seeksKind} onChange={(e) => setSeeksKind(e.target.value)} placeholder="coat" required />
        </label>
        <div className="span2 form-actions">
          <button type="submit" className="primary">Post offer</button>
        </div>
      </form>
      {offers.length > 0 && (
        <table className="commodity-table" style={{ marginTop: "1rem" }}>
          <thead><tr><th>Owner</th><th>Offers</th><th>Seeks</th><th></th></tr></thead>
          <tbody>
            {offers.map((o) => {
              const owner = ownerById.get(o.owner_id);
              const c = commodityById.get(o.commodity_id);
              return (
                <tr key={o.id}>
                  <td>{owner?.name ?? o.owner_id.slice(0, 8)}</td>
                  <td>{fmtQty(o.quantity)} {c?.use_value.unit ?? ""} {c?.name ?? ""}</td>
                  <td className="muted">{o.seeks_kind}</td>
                  <td className="actions"><button className="danger" onClick={() => setOffers((xs) => xs.filter((x) => x.id !== o.id))}>Withdraw</button></td>
                </tr>
              );
            })}
          </tbody>
        </table>
      )}
    </Card>
  );
}

function ExchangesPanel({ owners, commodities, exchanges, setExchanges }) {
  const [d, setD] = React.useState({
    giverID: "", receiverID: "", giverCommodityID: "", giverQty: 1,
    receiverCommodityID: "", receiverQty: 1, realisedValue: 0,
  });
  if (owners.length < 2 || commodities.length < 1) return null;

  const ownerById = new Map(owners.map((o) => [o.id, o]));
  const commodityById = new Map(commodities.map((c) => [c.id, c]));

  function submit(e) {
    e.preventDefault();
    setExchanges((xs) => [...xs, {
      id: "ex_" + Math.random().toString(36).slice(2, 10),
      giver_id: d.giverID, receiver_id: d.receiverID,
      giver_commodity_id: d.giverCommodityID, giver_qty: d.giverQty,
      receiver_commodity_id: d.receiverCommodityID, receiver_qty: d.receiverQty,
      realised_value: d.realisedValue,
    }]);
  }

  return (
    <Card
      title="Record an exchange"
      description='"Commodities must be realised as values before they can be realised as use-values." (Ch. 2)'
    >
      <form className="form-grid" onSubmit={submit}>
        <label><span>Giver</span>
          <select value={d.giverID} onChange={(e) => setD({ ...d, giverID: e.target.value })} required>
            <option value="">—</option>{owners.map((o) => <option key={o.id} value={o.id}>{o.name}</option>)}
          </select>
        </label>
        <label><span>Giver commodity</span>
          <select value={d.giverCommodityID} onChange={(e) => setD({ ...d, giverCommodityID: e.target.value })} required>
            <option value="">—</option>{commodities.map((c) => <option key={c.id} value={c.id}>{c.name}</option>)}
          </select>
        </label>
        <label><span>Giver qty</span>
          <input type="number" min={0} step="any" value={d.giverQty} onChange={(e) => setD({ ...d, giverQty: Number(e.target.value) })} required />
        </label>
        <label><span>Receiver</span>
          <select value={d.receiverID} onChange={(e) => setD({ ...d, receiverID: e.target.value })} required>
            <option value="">—</option>{owners.map((o) => <option key={o.id} value={o.id}>{o.name}</option>)}
          </select>
        </label>
        <label><span>Receiver commodity</span>
          <select value={d.receiverCommodityID} onChange={(e) => setD({ ...d, receiverCommodityID: e.target.value })} required>
            <option value="">—</option>{commodities.map((c) => <option key={c.id} value={c.id}>{c.name}</option>)}
          </select>
        </label>
        <label><span>Receiver qty</span>
          <input type="number" min={0} step="any" value={d.receiverQty} onChange={(e) => setD({ ...d, receiverQty: Number(e.target.value) })} required />
        </label>
        <label className="span2"><span>Realised value (labour-minutes)</span>
          <input type="number" min={1} value={d.realisedValue} onChange={(e) => setD({ ...d, realisedValue: Number(e.target.value) })} required />
        </label>
        <div className="span2 form-actions">
          <button type="submit" className="primary">Record exchange</button>
        </div>
      </form>
      {exchanges.length > 0 && (
        <table className="commodity-table" style={{ marginTop: "1rem" }}>
          <thead><tr><th>Giver</th><th>Gives</th><th>Receiver</th><th>Receives</th><th>Value</th></tr></thead>
          <tbody>
            {exchanges.map((e) => {
              const gC = commodityById.get(e.giver_commodity_id);
              const rC = commodityById.get(e.receiver_commodity_id);
              return (
                <tr key={e.id}>
                  <td>{ownerById.get(e.giver_id)?.name}</td>
                  <td>{fmtQty(e.giver_qty)} {gC?.name ?? ""}</td>
                  <td>{ownerById.get(e.receiver_id)?.name}</td>
                  <td>{fmtQty(e.receiver_qty)} {rC?.name ?? ""}</td>
                  <td className="snlt-cell">{fmtMinutes(e.realised_value)}</td>
                </tr>
              );
            })}
          </tbody>
        </table>
      )}
    </Card>
  );
}

function MoneyPanel({ commodities, money, setMoney }) {
  const { universalEquivalent, moneyCommodity, prices } = money;
  const [ueID, setUeID] = React.useState("");
  const [pi, setPi] = React.useState({ commodity_id: "", money_commodity_id: "", commodity_snlt: 0, money_snlt: 0, unit_qty: 1 });
  if (commodities.length === 0) return null;

  const commodityById = new Map(commodities.map((c) => [c.id, c]));
  const ueComm = universalEquivalent ? commodityById.get(universalEquivalent.commodity_id) : null;
  const mcComm = moneyCommodity ? commodityById.get(moneyCommodity.commodity_id) : null;

  function electUE(e) { e.preventDefault(); if (!ueID) return; setMoney({ ...money, universalEquivalent: { commodity_id: ueID } }); }
  function crystallise() { if (!universalEquivalent) return; setMoney({ ...money, moneyCommodity: { commodity_id: universalEquivalent.commodity_id } }); }
  function computePrice(e) {
    e.preventDefault();
    if (!pi.commodity_id || !pi.money_commodity_id || !pi.money_snlt) return;
    const amount = (pi.commodity_snlt * pi.unit_qty) / pi.money_snlt;
    const next = prices.filter((p) => p.commodity_id !== pi.commodity_id || p.money_commodity_id !== pi.money_commodity_id);
    setMoney({ ...money, prices: [...next, { commodity_id: pi.commodity_id, money_commodity_id: pi.money_commodity_id, amount: Number(amount.toFixed(3)) }] });
  }

  return (
    <Card
      title="Universal equivalent & money"
      description='"The social action of all other commodities sets apart the particular commodity in which they all represent their values. Thereby it becomes — money." (Ch. 2)'
    >
      <div style={{ marginBottom: "1.5rem" }}>
        <h3 style={{ marginBottom: "0.5rem" }}>Elect universal equivalent</h3>
        {universalEquivalent && (
          <p className="reveal-note">
            Current: <strong>{ueComm?.name ?? universalEquivalent.commodity_id}</strong>
            {moneyCommodity ? <> &rarr; crystallised as money (<strong>{mcComm?.name}</strong>)</> : <> (not yet crystallised into money)</>}
          </p>
        )}
        <form className="form-grid" onSubmit={electUE}>
          <label className="span2"><span>Elect as universal equivalent</span>
            <select value={ueID} onChange={(e) => setUeID(e.target.value)} required>
              <option value="">—</option>{commodities.map((c) => <option key={c.id} value={c.id}>{c.name}</option>)}
            </select>
          </label>
          <div className="span2 form-actions">
            <button type="submit" className="primary">Set universal equivalent</button>
            {!moneyCommodity && universalEquivalent && (
              <button type="button" className="secondary" onClick={crystallise}>Crystallise into money</button>
            )}
          </div>
        </form>
      </div>
      <div>
        <h3 style={{ marginBottom: "0.5rem" }}>Compute price</h3>
        <p className="description small">
          Price = commodity SNLT × qty ÷ money SNLT — halving money SNLT doubles the price. (Petty, fn. 12)
        </p>
        <form className="form-grid" onSubmit={computePrice}>
          <label><span>Commodity</span>
            <select value={pi.commodity_id} onChange={(e) => setPi({ ...pi, commodity_id: e.target.value, commodity_snlt: commodityById.get(e.target.value)?.snlt_per_unit ?? 0 })} required>
              <option value="">—</option>{commodities.map((c) => <option key={c.id} value={c.id}>{c.name}</option>)}
            </select>
          </label>
          <label><span>Commodity SNLT (min/unit)</span>
            <input type="number" min={0} value={pi.commodity_snlt} onChange={(e) => setPi({ ...pi, commodity_snlt: Number(e.target.value) })} required />
          </label>
          <label><span>Money commodity</span>
            <select value={pi.money_commodity_id} onChange={(e) => setPi({ ...pi, money_commodity_id: e.target.value, money_snlt: commodityById.get(e.target.value)?.snlt_per_unit ?? 0 })} required>
              <option value="">—</option>{commodities.map((c) => <option key={c.id} value={c.id}>{c.name}</option>)}
            </select>
          </label>
          <label><span>Money SNLT (min/unit)</span>
            <input type="number" min={1} value={pi.money_snlt} onChange={(e) => setPi({ ...pi, money_snlt: Number(e.target.value) })} required />
          </label>
          <label><span>Unit qty</span>
            <input type="number" min={1} value={pi.unit_qty} onChange={(e) => setPi({ ...pi, unit_qty: Number(e.target.value) })} required />
          </label>
          <div className="span2 form-actions">
            <button type="submit" className="primary">Compute & store price</button>
          </div>
        </form>
        {prices.length > 0 && (
          <table className="commodity-table" style={{ marginTop: "1rem" }}>
            <thead><tr><th>Commodity</th><th>Money commodity</th><th>Price (units)</th></tr></thead>
            <tbody>
              {prices.map((p, i) => (
                <tr key={i}>
                  <td>{commodityById.get(p.commodity_id)?.name}</td>
                  <td className="muted">{commodityById.get(p.money_commodity_id)?.name}</td>
                  <td className="snlt-cell">{p.amount}</td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
    </Card>
  );
}

function Ch02ExchangeView({ commodities, ch02 }) {
  return (
    <>
      <OwnersPanel owners={ch02.owners} setOwners={(fn) => ch02.set("owners", fn)} />
      <OffersPanel owners={ch02.owners} commodities={commodities} offers={ch02.offers} setOffers={(fn) => ch02.set("offers", fn)} />
      <ExchangesPanel owners={ch02.owners} commodities={commodities} exchanges={ch02.exchanges} setExchanges={(fn) => ch02.set("exchanges", fn)} />
      <MoneyPanel commodities={commodities} money={ch02.money} setMoney={(v) => ch02.set("money", () => v)} />
    </>
  );
}

Object.assign(window, { OwnersPanel, OffersPanel, ExchangesPanel, MoneyPanel, Ch02ExchangeView });
