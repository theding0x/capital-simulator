import { useEffect, useMemo, useState } from "react";
import type { FormEvent } from "react";
import { api } from "../api";
import type {
  Commodity,
  ComputePriceInput,
  Exchange,
  MoneyCommodity,
  Offer,
  Owner,
  Price,
  UniversalEquivalent,
} from "../types";
import { fmtMinutes, fmtQty } from "../format";

interface Ch02Props {
  commodities: Commodity[];
  owners: Owner[];
  onSharedChanged: () => void;
}

export function Ch02Exchange({ commodities, owners, onSharedChanged }: Ch02Props) {
  const [offers, setOffers] = useState<Offer[]>([]);
  const [exchanges, setExchanges] = useState<Exchange[]>([]);
  const [universalEquivalent, setUniversalEquivalent] = useState<UniversalEquivalent | null>(null);
  const [moneyCommodity, setMoneyCommodity] = useState<MoneyCommodity | null>(null);
  const [prices, setPrices] = useState<Price[]>([]);

  async function refreshLocal() {
    const [ofs, exs, ps] = await Promise.allSettled([
      api.listOffers(),
      api.listExchanges(),
      api.listPrices(),
    ]);
    if (ofs.status === "fulfilled") setOffers(ofs.value);
    if (exs.status === "fulfilled") setExchanges(exs.value);
    if (ps.status === "fulfilled") setPrices(ps.value);
    try {
      setUniversalEquivalent(await api.getUniversalEquivalent());
    } catch {
      setUniversalEquivalent(null);
    }
    try {
      setMoneyCommodity(await api.getMoneyCommodity());
    } catch {
      setMoneyCommodity(null);
    }
  }

  useEffect(() => { refreshLocal(); }, []);

  const commodityById = useMemo(() => {
    const m = new Map<string, Commodity>();
    for (const c of commodities) m.set(c.id, c);
    return m;
  }, [commodities]);

  const ownerById = useMemo(() => {
    const m = new Map<string, Owner>();
    for (const o of owners) m.set(o.id, o);
    return m;
  }, [owners]);

  return (
    <>
      <OwnersPanel owners={owners} onOwnerCreated={onSharedChanged} />
      <OffersPanel
        owners={owners}
        offers={offers}
        commodities={commodities}
        commodityById={commodityById}
        ownerById={ownerById}
        onChanged={refreshLocal}
      />
      <ExchangesPanel
        owners={owners}
        exchanges={exchanges}
        commodities={commodities}
        commodityById={commodityById}
        ownerById={ownerById}
        onChanged={refreshLocal}
      />
      <MoneyPanel
        commodities={commodities}
        commodityById={commodityById}
        universalEquivalent={universalEquivalent}
        moneyCommodity={moneyCommodity}
        prices={prices}
        onChanged={refreshLocal}
      />
    </>
  );
}

function OwnersPanel({
  owners,
  onOwnerCreated,
}: {
  owners: Owner[];
  onOwnerCreated: () => void;
}) {
  const [name, setName] = useState("");
  const [busy, setBusy] = useState(false);
  const [err, setErr] = useState<string | null>(null);

  async function submit(e: FormEvent) {
    e.preventDefault();
    setBusy(true);
    setErr(null);
    try {
      await api.createOwner({ name });
      setName("");
      onOwnerCreated();
    } catch (e) {
      setErr(e instanceof Error ? e.message : String(e));
    } finally {
      setBusy(false);
    }
  }

  return (
    <section className="card">
      <h2>Owners</h2>
      <p className="description">
        "Commodities cannot go to market of their own account. We must have recourse to their
        guardians, who are also their owners." (Ch. 2)
      </p>
      <form className="form-grid" onSubmit={submit}>
        <label className="span2">
          <span>Owner name</span>
          <input
            value={name}
            onChange={(e) => setName(e.target.value)}
            placeholder="weaver"
            required
          />
        </label>
        <div className="span2 form-actions">
          <button type="submit" className="primary" disabled={busy}>
            {busy ? "Saving…" : "Register owner"}
          </button>
          {err && <span className="error">{err}</span>}
        </div>
      </form>
      {owners.length > 0 && (
        <table className="commodity-table" style={{ marginTop: "1rem" }}>
          <thead>
            <tr>
              <th>Name</th>
              <th>ID</th>
            </tr>
          </thead>
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
    </section>
  );
}

function OffersPanel({
  owners,
  offers,
  commodities,
  commodityById,
  ownerById,
  onChanged,
}: {
  owners: Owner[];
  offers: Offer[];
  commodities: Commodity[];
  commodityById: Map<string, Commodity>;
  ownerById: Map<string, Owner>;
  onChanged: () => void;
}) {
  const [ownerID, setOwnerID] = useState("");
  const [commodityID, setCommodityID] = useState("");
  const [qty, setQty] = useState(1);
  const [seeksKind, setSeeksKind] = useState("");
  const [busy, setBusy] = useState(false);
  const [err, setErr] = useState<string | null>(null);

  async function submit(e: FormEvent) {
    e.preventDefault();
    setBusy(true);
    setErr(null);
    try {
      await api.createOffer({
        owner_id: ownerID,
        commodity_id: commodityID,
        quantity: qty,
        seeks_kind: seeksKind,
      });
      setQty(1);
      setSeeksKind("");
      onChanged();
    } catch (e) {
      setErr(e instanceof Error ? e.message : String(e));
    } finally {
      setBusy(false);
    }
  }

  if (owners.length === 0 || commodities.length === 0) return null;

  return (
    <section className="card">
      <h2>Market offers</h2>
      <p className="description">
        Each owner brings their commodity to market as a non-use-value for themselves, seeking
        another's use-value in return.
      </p>
      <form className="form-grid" onSubmit={submit}>
        <label>
          <span>Owner</span>
          <select value={ownerID} onChange={(e) => setOwnerID(e.target.value)} required>
            <option value="">—</option>
            {owners.map((o) => (
              <option key={o.id} value={o.id}>
                {o.name}
              </option>
            ))}
          </select>
        </label>
        <label>
          <span>Commodity offered</span>
          <select value={commodityID} onChange={(e) => setCommodityID(e.target.value)} required>
            <option value="">—</option>
            {commodities.map((c) => (
              <option key={c.id} value={c.id}>
                {c.name}
              </option>
            ))}
          </select>
        </label>
        <label>
          <span>Quantity</span>
          <input
            type="number"
            min={0}
            step="any"
            value={qty}
            onChange={(e) => setQty(Number(e.target.value))}
            required
          />
        </label>
        <label>
          <span>Seeks (kind)</span>
          <input
            value={seeksKind}
            onChange={(e) => setSeeksKind(e.target.value)}
            placeholder="coat"
            required
          />
        </label>
        <div className="span2 form-actions">
          <button type="submit" className="primary" disabled={busy}>
            {busy ? "Posting…" : "Post offer"}
          </button>
          {err && <span className="error">{err}</span>}
        </div>
      </form>
      {offers.length > 0 && (
        <table className="commodity-table" style={{ marginTop: "1rem" }}>
          <thead>
            <tr>
              <th>Owner</th>
              <th>Offers</th>
              <th>Seeks</th>
              <th></th>
            </tr>
          </thead>
          <tbody>
            {offers.map((o) => {
              const owner = ownerById.get(o.owner_id);
              const commodity = commodityById.get(o.commodity_id);
              return (
                <tr key={o.id}>
                  <td>{owner?.name ?? o.owner_id.slice(0, 8)}</td>
                  <td>
                    {fmtQty(o.quantity)} {commodity?.use_value.unit ?? ""}{" "}
                    {commodity?.name ?? o.commodity_id.slice(0, 8)}
                  </td>
                  <td className="muted">{o.seeks_kind}</td>
                  <td>
                    <button
                      className="danger"
                      onClick={async () => {
                        await api.deleteOffer(o.id);
                        onChanged();
                      }}
                    >
                      Withdraw
                    </button>
                  </td>
                </tr>
              );
            })}
          </tbody>
        </table>
      )}
    </section>
  );
}

function ExchangesPanel({
  owners,
  exchanges,
  commodities,
  commodityById,
  ownerById,
  onChanged,
}: {
  owners: Owner[];
  exchanges: Exchange[];
  commodities: Commodity[];
  commodityById: Map<string, Commodity>;
  ownerById: Map<string, Owner>;
  onChanged: () => void;
}) {
  const [giverID, setGiverID] = useState("");
  const [receiverID, setReceiverID] = useState("");
  const [giverCommodityID, setGiverCommodityID] = useState("");
  const [giverQty, setGiverQty] = useState(1);
  const [receiverCommodityID, setReceiverCommodityID] = useState("");
  const [receiverQty, setReceiverQty] = useState(1);
  const [realisedValue, setRealisedValue] = useState(0);
  const [busy, setBusy] = useState(false);
  const [err, setErr] = useState<string | null>(null);

  async function submit(e: FormEvent) {
    e.preventDefault();
    setBusy(true);
    setErr(null);
    try {
      await api.createExchange({
        giver_id: giverID,
        receiver_id: receiverID,
        giver_commodity_id: giverCommodityID,
        giver_qty: giverQty,
        receiver_commodity_id: receiverCommodityID,
        receiver_qty: receiverQty,
        realised_value: realisedValue,
      });
      onChanged();
    } catch (e) {
      setErr(e instanceof Error ? e.message : String(e));
    } finally {
      setBusy(false);
    }
  }

  if (owners.length < 2 || commodities.length < 1) return null;

  return (
    <section className="card">
      <h2>Record an exchange</h2>
      <p className="description">
        "Commodities must be realised as values before they can be realised as use-values." (Ch. 2)
      </p>
      <form className="form-grid" onSubmit={submit}>
        <label>
          <span>Giver</span>
          <select value={giverID} onChange={(e) => setGiverID(e.target.value)} required>
            <option value="">—</option>
            {owners.map((o) => (
              <option key={o.id} value={o.id}>
                {o.name}
              </option>
            ))}
          </select>
        </label>
        <label>
          <span>Giver commodity</span>
          <select
            value={giverCommodityID}
            onChange={(e) => setGiverCommodityID(e.target.value)}
            required
          >
            <option value="">—</option>
            {commodities.map((c) => (
              <option key={c.id} value={c.id}>
                {c.name}
              </option>
            ))}
          </select>
        </label>
        <label>
          <span>Giver qty</span>
          <input
            type="number"
            min={0}
            step="any"
            value={giverQty}
            onChange={(e) => setGiverQty(Number(e.target.value))}
            required
          />
        </label>
        <label>
          <span>Receiver</span>
          <select value={receiverID} onChange={(e) => setReceiverID(e.target.value)} required>
            <option value="">—</option>
            {owners.map((o) => (
              <option key={o.id} value={o.id}>
                {o.name}
              </option>
            ))}
          </select>
        </label>
        <label>
          <span>Receiver commodity</span>
          <select
            value={receiverCommodityID}
            onChange={(e) => setReceiverCommodityID(e.target.value)}
            required
          >
            <option value="">—</option>
            {commodities.map((c) => (
              <option key={c.id} value={c.id}>
                {c.name}
              </option>
            ))}
          </select>
        </label>
        <label>
          <span>Receiver qty</span>
          <input
            type="number"
            min={0}
            step="any"
            value={receiverQty}
            onChange={(e) => setReceiverQty(Number(e.target.value))}
            required
          />
        </label>
        <label>
          <span>Realised value (labour-minutes)</span>
          <input
            type="number"
            min={1}
            value={realisedValue}
            onChange={(e) => setRealisedValue(Number(e.target.value))}
            required
          />
        </label>
        <div className="span2 form-actions">
          <button type="submit" className="primary" disabled={busy}>
            {busy ? "Recording…" : "Record exchange"}
          </button>
          {err && <span className="error">{err}</span>}
        </div>
      </form>
      {exchanges.length > 0 && (
        <table className="commodity-table" style={{ marginTop: "1rem" }}>
          <thead>
            <tr>
              <th>Giver</th>
              <th>Gives</th>
              <th>Receiver</th>
              <th>Receives</th>
              <th>Value</th>
            </tr>
          </thead>
          <tbody>
            {exchanges.map((e) => {
              const giver = ownerById.get(e.giver_id);
              const receiver = ownerById.get(e.receiver_id);
              const giverC = commodityById.get(e.giver_commodity_id);
              const receiverC = commodityById.get(e.receiver_commodity_id);
              return (
                <tr key={e.id}>
                  <td>{giver?.name ?? e.giver_id.slice(0, 8)}</td>
                  <td>
                    {fmtQty(e.giver_qty)} {giverC?.name ?? e.giver_commodity_id.slice(0, 8)}
                  </td>
                  <td>{receiver?.name ?? e.receiver_id.slice(0, 8)}</td>
                  <td>
                    {fmtQty(e.receiver_qty)}{" "}
                    {receiverC?.name ?? e.receiver_commodity_id.slice(0, 8)}
                  </td>
                  <td className="snlt-cell">{fmtMinutes(e.realised_value)}</td>
                </tr>
              );
            })}
          </tbody>
        </table>
      )}
    </section>
  );
}

function MoneyPanel({
  commodities,
  commodityById,
  universalEquivalent,
  moneyCommodity,
  prices,
  onChanged,
}: {
  commodities: Commodity[];
  commodityById: Map<string, Commodity>;
  universalEquivalent: UniversalEquivalent | null;
  moneyCommodity: MoneyCommodity | null;
  prices: Price[];
  onChanged: () => void;
}) {
  const [ueID, setUeID] = useState("");
  const [ueBusy, setUeBusy] = useState(false);
  const [ueErr, setUeErr] = useState<string | null>(null);
  const [mcBusy, setMcBusy] = useState(false);
  const [mcErr, setMcErr] = useState<string | null>(null);
  const [priceInput, setPriceInput] = useState<ComputePriceInput>({
    commodity_id: "",
    money_commodity_id: "",
    commodity_snlt: 0,
    money_snlt: 0,
    unit_qty: 1,
  });
  const [priceBusy, setPriceBusy] = useState(false);
  const [priceErr, setPriceErr] = useState<string | null>(null);

  async function electUE(e: FormEvent) {
    e.preventDefault();
    setUeBusy(true);
    setUeErr(null);
    try {
      await api.setUniversalEquivalent(ueID);
      onChanged();
    } catch (e) {
      setUeErr(e instanceof Error ? e.message : String(e));
    } finally {
      setUeBusy(false);
    }
  }

  async function crystallise() {
    if (!universalEquivalent) return;
    setMcBusy(true);
    setMcErr(null);
    try {
      await api.setMoneyCommodity(universalEquivalent.commodity_id);
      onChanged();
    } catch (e) {
      setMcErr(e instanceof Error ? e.message : String(e));
    } finally {
      setMcBusy(false);
    }
  }

  async function computePrice(e: FormEvent) {
    e.preventDefault();
    setPriceBusy(true);
    setPriceErr(null);
    try {
      await api.computePrice(priceInput);
      onChanged();
    } catch (e) {
      setPriceErr(e instanceof Error ? e.message : String(e));
    } finally {
      setPriceBusy(false);
    }
  }

  if (commodities.length === 0) return null;

  const ueComm = universalEquivalent ? commodityById.get(universalEquivalent.commodity_id) : null;
  const mcComm = moneyCommodity ? commodityById.get(moneyCommodity.commodity_id) : null;

  return (
    <section className="card">
      <h2>Universal equivalent &amp; money</h2>
      <p className="description">
        "The social action of all other commodities sets apart the particular commodity in which
        they all represent their values. Thereby it becomes — money." (Ch. 2)
      </p>
      <div style={{ marginBottom: "1.5rem" }}>
        <h3 style={{ marginBottom: "0.5rem" }}>Elect universal equivalent</h3>
        {universalEquivalent && (
          <p className="reveal-note">
            Current: <strong>{ueComm?.name ?? universalEquivalent.commodity_id}</strong>
            {moneyCommodity ? (
              <>
                {" "}
                &rarr; crystallised as money (
                <strong>{mcComm?.name ?? moneyCommodity.commodity_id}</strong>)
              </>
            ) : (
              <> (not yet crystallised into money)</>
            )}
          </p>
        )}
        <form className="form-grid" onSubmit={electUE}>
          <label className="span2">
            <span>Elect as universal equivalent</span>
            <select value={ueID} onChange={(e) => setUeID(e.target.value)} required>
              <option value="">—</option>
              {commodities.map((c) => (
                <option key={c.id} value={c.id}>
                  {c.name}
                </option>
              ))}
            </select>
          </label>
          <div className="span2 form-actions">
            <button type="submit" className="primary" disabled={ueBusy}>
              {ueBusy ? "Setting…" : "Set universal equivalent"}
            </button>
            {!moneyCommodity && universalEquivalent && (
              <button
                type="button"
                className="secondary"
                disabled={mcBusy}
                onClick={crystallise}
              >
                {mcBusy ? "Crystallising…" : "Crystallise into money"}
              </button>
            )}
            {ueErr && <span className="error">{ueErr}</span>}
            {mcErr && <span className="error">{mcErr}</span>}
          </div>
        </form>
      </div>
      <div>
        <h3 style={{ marginBottom: "0.5rem" }}>Compute price</h3>
        <p className="description small">
          Price = commodity SNLT &times; qty &divide; money SNLT &mdash; halving money SNLT doubles
          the price. (Petty, fn. 12)
        </p>
        <form className="form-grid" onSubmit={computePrice}>
          <label>
            <span>Commodity</span>
            <select
              value={priceInput.commodity_id}
              onChange={(e) => setPriceInput({ ...priceInput, commodity_id: e.target.value })}
              required
            >
              <option value="">—</option>
              {commodities.map((c) => (
                <option key={c.id} value={c.id}>
                  {c.name}
                </option>
              ))}
            </select>
          </label>
          <label>
            <span>Commodity SNLT (min/unit)</span>
            <input
              type="number"
              min={0}
              value={priceInput.commodity_snlt}
              onChange={(e) =>
                setPriceInput({ ...priceInput, commodity_snlt: Number(e.target.value) })
              }
              required
            />
          </label>
          <label>
            <span>Money commodity</span>
            <select
              value={priceInput.money_commodity_id}
              onChange={(e) =>
                setPriceInput({ ...priceInput, money_commodity_id: e.target.value })
              }
              required
            >
              <option value="">—</option>
              {commodities.map((c) => (
                <option key={c.id} value={c.id}>
                  {c.name}
                </option>
              ))}
            </select>
          </label>
          <label>
            <span>Money SNLT (min/unit)</span>
            <input
              type="number"
              min={1}
              value={priceInput.money_snlt}
              onChange={(e) =>
                setPriceInput({ ...priceInput, money_snlt: Number(e.target.value) })
              }
              required
            />
          </label>
          <label>
            <span>Unit qty</span>
            <input
              type="number"
              min={1}
              value={priceInput.unit_qty}
              onChange={(e) => setPriceInput({ ...priceInput, unit_qty: Number(e.target.value) })}
              required
            />
          </label>
          <div className="span2 form-actions">
            <button type="submit" className="primary" disabled={priceBusy}>
              {priceBusy ? "Computing…" : "Compute & store price"}
            </button>
            {priceErr && <span className="error">{priceErr}</span>}
          </div>
        </form>
        {prices.length > 0 && (
          <table className="commodity-table" style={{ marginTop: "1rem" }}>
            <thead>
              <tr>
                <th>Commodity</th>
                <th>Money commodity</th>
                <th>Price (units)</th>
              </tr>
            </thead>
            <tbody>
              {prices.map((p) => {
                const comm = commodityById.get(p.commodity_id);
                const money = commodityById.get(p.money_commodity_id);
                return (
                  <tr key={p.commodity_id}>
                    <td>{comm?.name ?? p.commodity_id.slice(0, 10)}</td>
                    <td className="muted">{money?.name ?? p.money_commodity_id.slice(0, 10)}</td>
                    <td className="snlt-cell">{p.amount}</td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        )}
      </div>
    </section>
  );
}
