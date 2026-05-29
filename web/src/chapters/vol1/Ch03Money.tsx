import { useEffect, useState } from "react";
import type { FormEvent } from "react";
import { api } from "../../api";
import type {
  Circuit,
  CircuitLeg,
  CreateCircuitInput,
  CreateHoardInput,
  CreatePaymentObligationInput,
  CreateWorldMoneyTransferInput,
  Hoard,
  Owner,
  PaymentObligation,
  WorldMoneyTransfer,
} from "../../types";
import { fmtMinutes } from "../../format";
import "./Ch03Money.css";

function MoneyFunctionsInsight() {
  return (
    <section className="v1-ch03-insight">
      <h2 className="v1-ch03-insight-h2">Three functions of money</h2>
      <p className="v1-ch03-insight-prose">
        Money is the universal equivalent, but it does its work in three
        distinct functions. The panels below probe each in turn — circulation
        as C—M—C, payment obligations as credit, and world money in
        cross-currency transfers.
      </p>
      <div className="v1-ch03-functions">
        <div className="v1-ch03-function">
          <span className="v1-ch03-function-tag">§2</span>
          <span className="v1-ch03-function-name">Medium of circulation</span>
          <span className="v1-ch03-function-gloss">
            C—M—C · money mediates the exchange of one commodity for another.
          </span>
        </div>
        <div className="v1-ch03-function">
          <span className="v1-ch03-function-tag">§3b</span>
          <span className="v1-ch03-function-name">Means of payment</span>
          <span className="v1-ch03-function-gloss">
            Deferred settlement · credit precedes the actual transfer of money.
          </span>
        </div>
        <div className="v1-ch03-function">
          <span className="v1-ch03-function-tag">§3c</span>
          <span className="v1-ch03-function-name">World money</span>
          <span className="v1-ch03-function-gloss">
            Beyond domestic territory · the universal equivalent in its
            unmediated, bullion form.
          </span>
        </div>
      </div>
    </section>
  );
}

function MoneyUniversalCoda() {
  return (
    <aside className="v1-ch03-coda">
      <p className="v1-ch03-coda-quote">
        “Money is the universal commodity, the medium of all sale and
        purchase, the measure of all values, the universal equivalent —
        because all commodities relate to it as the form in which their own
        value is realised.”
        <span className="v1-ch03-coda-cite">
          — Marx, Capital Vol. I, Ch. 3
        </span>
      </p>
    </aside>
  );
}

interface Ch03Props {
  owners: Owner[];
  onSharedChanged: () => void;
}

export function Ch03Money({ owners, onSharedChanged: _onSharedChanged }: Ch03Props) {
  const [circuits, setCircuits] = useState<Circuit[]>([]);
  const [hoards, setHoards] = useState<Hoard[]>([]);
  const [obligations, setObligations] = useState<PaymentObligation[]>([]);
  const [transfers, setTransfers] = useState<WorldMoneyTransfer[]>([]);

  async function refreshLocal() {
    const [cs, hs, obs, ts] = await Promise.allSettled([
      api.listCircuits(),
      api.listHoards(),
      api.listPaymentObligations(),
      api.listWorldMoneyTransfers(),
    ]);
    if (cs.status === "fulfilled") setCircuits(cs.value);
    if (hs.status === "fulfilled") setHoards(hs.value);
    if (obs.status === "fulfilled") setObligations(obs.value);
    if (ts.status === "fulfilled") setTransfers(ts.value);
  }

  useEffect(() => { refreshLocal(); }, []);

  return (
    <>
      <MoneyFunctionsInsight />
      <CirculationPanel />
      <CircuitPanel owners={owners} circuits={circuits} onChanged={refreshLocal} />
      <HoardPanel owners={owners} hoards={hoards} onChanged={refreshLocal} />
      <PaymentObligationsPanel owners={owners} obligations={obligations} onChanged={refreshLocal} />
      <WorldMoneyPanel owners={owners} transfers={transfers} onChanged={refreshLocal} />
      <MoneyUniversalCoda />
    </>
  );
}

function CirculationPanel() {
  const [sum, setSum] = useState(0);
  const [velocity, setVelocity] = useState(1);
  const [result, setResult] = useState<number | null>(null);
  const [err, setErr] = useState<string | null>(null);

  async function compute(e: FormEvent) {
    e.preventDefault();
    setErr(null);
    try {
      const r = await api.getMoneyRequired(sum, velocity);
      setResult(r.money_required);
    } catch (e) {
      setErr(e instanceof Error ? e.message : String(e));
    }
  }

  return (
    <section className="card">
      <h2>Quantity of circulating medium</h2>
      <p className="description">
        "The quantity of money as circulating medium = the sum of the prices of the commodities /
        the number of moves made by coins." (Ch. 3 §2B)
      </p>
      <form className="form-grid" onSubmit={compute}>
        <label>
          <span>Sum of prices</span>
          <input
            type="number"
            min={1}
            value={sum}
            onChange={(e) => setSum(Number(e.target.value))}
            required
          />
        </label>
        <label>
          <span>Velocity (moves per coin)</span>
          <input
            type="number"
            min={1}
            value={velocity}
            onChange={(e) => setVelocity(Number(e.target.value))}
            required
          />
        </label>
        <div className="span2 form-actions">
          <button type="submit" className="primary">
            Compute
          </button>
          {err && <span className="error">{err}</span>}
        </div>
      </form>
      {result !== null && (
        <p className="reveal-note" style={{ marginTop: "0.75rem" }}>
          Money required: <strong>{result}</strong> units (M = {sum} / {velocity})
        </p>
      )}
    </section>
  );
}

function CircuitPanel({
  owners,
  circuits,
  onChanged,
}: {
  owners: Owner[];
  circuits: Circuit[];
  onChanged: () => void;
}) {
  const emptySaleLeg = (): CircuitLeg => ({
    kind: "sale",
    commodity_id: "",
    money_id: "",
    owner_id: "",
    value: 0,
  });
  const emptyPurchaseLeg = (): CircuitLeg => ({
    kind: "purchase",
    commodity_id: "",
    money_id: "",
    owner_id: "",
    value: 0,
  });
  const [sale, setSale] = useState<CircuitLeg>(emptySaleLeg);
  const [purchase, setPurchase] = useState<CircuitLeg>(emptyPurchaseLeg);
  const [busy, setBusy] = useState(false);
  const [err, setErr] = useState<string | null>(null);

  async function submit(e: FormEvent) {
    e.preventDefault();
    setBusy(true);
    setErr(null);
    try {
      const input: CreateCircuitInput = { sale_leg: sale, purchase_leg: purchase };
      await api.createCircuit(input);
      setSale(emptySaleLeg());
      setPurchase(emptyPurchaseLeg());
      onChanged();
    } catch (e) {
      setErr(e instanceof Error ? e.message : String(e));
    } finally {
      setBusy(false);
    }
  }

  return (
    <section className="card">
      <h2>C-M-C Circuits</h2>
      <p className="description">
        "The act C-M-C — sell in order to buy — a commodity is converted into money, then the money
        is re-converted into a commodity." (Ch. 3 §2A)
      </p>
      <form className="form-grid" onSubmit={submit}>
        <h3 className="span2" style={{ marginBottom: 0 }}>
          Sale leg (C → M)
        </h3>
        <label>
          <span>Owner</span>
          <select
            value={sale.owner_id}
            onChange={(e) => setSale({ ...sale, owner_id: e.target.value })}
            required
          >
            <option value="">—</option>
            {owners.map((o) => (
              <option key={o.id} value={o.id}>
                {o.name}
              </option>
            ))}
          </select>
        </label>
        <label>
          <span>Commodity ID sold</span>
          <input
            value={sale.commodity_id}
            onChange={(e) => setSale({ ...sale, commodity_id: e.target.value })}
            placeholder="commodity hex ID"
            required
          />
        </label>
        <label>
          <span>Money commodity ID</span>
          <input
            value={sale.money_id}
            onChange={(e) => setSale({ ...sale, money_id: e.target.value })}
            placeholder="money hex ID"
            required
          />
        </label>
        <label>
          <span>Value (labour-minutes)</span>
          <input
            type="number"
            min={1}
            value={sale.value}
            onChange={(e) => setSale({ ...sale, value: Number(e.target.value) })}
            required
          />
        </label>
        <h3 className="span2" style={{ marginBottom: 0 }}>
          Purchase leg (M → C′)
        </h3>
        <label>
          <span>Owner</span>
          <select
            value={purchase.owner_id}
            onChange={(e) => setPurchase({ ...purchase, owner_id: e.target.value })}
            required
          >
            <option value="">—</option>
            {owners.map((o) => (
              <option key={o.id} value={o.id}>
                {o.name}
              </option>
            ))}
          </select>
        </label>
        <label>
          <span>Commodity ID bought</span>
          <input
            value={purchase.commodity_id}
            onChange={(e) => setPurchase({ ...purchase, commodity_id: e.target.value })}
            placeholder="commodity hex ID"
            required
          />
        </label>
        <label>
          <span>Money commodity ID</span>
          <input
            value={purchase.money_id}
            onChange={(e) => setPurchase({ ...purchase, money_id: e.target.value })}
            placeholder="money hex ID"
            required
          />
        </label>
        <label>
          <span>Value (labour-minutes)</span>
          <input
            type="number"
            min={1}
            value={purchase.value}
            onChange={(e) => setPurchase({ ...purchase, value: Number(e.target.value) })}
            required
          />
        </label>
        <div className="span2 form-actions">
          <button type="submit" className="primary" disabled={busy}>
            {busy ? "Recording…" : "Record circuit"}
          </button>
          {err && <span className="error">{err}</span>}
        </div>
      </form>
      {circuits.length > 0 && (
        <table className="commodity-table" style={{ marginTop: "1rem" }}>
          <thead>
            <tr>
              <th>ID</th>
              <th>Sells</th>
              <th>Buys</th>
              <th>Value</th>
            </tr>
          </thead>
          <tbody>
            {circuits.map((c) => (
              <tr key={c.id}>
                <td className="muted small">{c.id.slice(0, 10)}…</td>
                <td className="muted small">{c.sale_leg.commodity_id.slice(0, 8)}</td>
                <td className="muted small">{c.purchase_leg.commodity_id.slice(0, 8)}</td>
                <td className="snlt-cell">{fmtMinutes(c.sale_leg.value)}</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </section>
  );
}

function HoardPanel({
  owners,
  hoards,
  onChanged,
}: {
  owners: Owner[];
  hoards: Hoard[];
  onChanged: () => void;
}) {
  const [ownerID, setOwnerID] = useState("");
  const [amount, setAmount] = useState(0);
  const [busy, setBusy] = useState(false);
  const [err, setErr] = useState<string | null>(null);

  async function submit(e: FormEvent) {
    e.preventDefault();
    setBusy(true);
    setErr(null);
    try {
      const input: CreateHoardInput = { owner_id: ownerID, amount };
      await api.createHoard(input);
      setAmount(0);
      onChanged();
    } catch (e) {
      setErr(e instanceof Error ? e.message : String(e));
    } finally {
      setBusy(false);
    }
  }

  if (owners.length === 0) return null;

  const ownerById = new Map(owners.map((o) => [o.id, o]));

  return (
    <section className="card">
      <h2>Hoards</h2>
      <p className="description">
        "To accumulate money it is necessary to sell, to keep on selling, and to abstain from
        buying." (Ch. 3 §3A)
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
          <span>Amount (units)</span>
          <input
            type="number"
            min={1}
            value={amount}
            onChange={(e) => setAmount(Number(e.target.value))}
            required
          />
        </label>
        <div className="span2 form-actions">
          <button type="submit" className="primary" disabled={busy}>
            {busy ? "Saving…" : "Hoard"}
          </button>
          {err && <span className="error">{err}</span>}
        </div>
      </form>
      {hoards.length > 0 && (
        <table className="commodity-table" style={{ marginTop: "1rem" }}>
          <thead>
            <tr>
              <th>Owner</th>
              <th>Amount</th>
            </tr>
          </thead>
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
    </section>
  );
}

function PaymentObligationsPanel({
  owners,
  obligations,
  onChanged,
}: {
  owners: Owner[];
  obligations: PaymentObligation[];
  onChanged: () => void;
}) {
  const [creditorID, setCreditorID] = useState("");
  const [debtorID, setDebtorID] = useState("");
  const [amount, setAmount] = useState(0);
  const [dueAt, setDueAt] = useState("");
  const [busy, setBusy] = useState(false);
  const [err, setErr] = useState<string | null>(null);

  async function submit(e: FormEvent) {
    e.preventDefault();
    setBusy(true);
    setErr(null);
    try {
      const input: CreatePaymentObligationInput = {
        creditor_id: creditorID,
        debtor_id: debtorID,
        amount,
        due_at: new Date(dueAt).toISOString(),
      };
      await api.createPaymentObligation(input);
      setAmount(0);
      setDueAt("");
      onChanged();
    } catch (e) {
      setErr(e instanceof Error ? e.message : String(e));
    } finally {
      setBusy(false);
    }
  }

  async function settle(id: string) {
    try {
      await api.settlePaymentObligation(id);
      onChanged();
    } catch (e) {
      alert(e instanceof Error ? e.message : String(e));
    }
  }

  if (owners.length === 0) return null;

  const ownerById = new Map(owners.map((o) => [o.id, o]));

  return (
    <section className="card">
      <h2>Payment obligations</h2>
      <p className="description">
        "The vendor becomes a creditor, the purchaser a debtor. Since money now functions as a
        means of payment after the commodity has changed hands, it becomes a debtor obligation."
        (Ch. 3 §3B)
      </p>
      <form className="form-grid" onSubmit={submit}>
        <label>
          <span>Creditor (vendor)</span>
          <select value={creditorID} onChange={(e) => setCreditorID(e.target.value)} required>
            <option value="">—</option>
            {owners.map((o) => (
              <option key={o.id} value={o.id}>
                {o.name}
              </option>
            ))}
          </select>
        </label>
        <label>
          <span>Debtor (purchaser)</span>
          <select value={debtorID} onChange={(e) => setDebtorID(e.target.value)} required>
            <option value="">—</option>
            {owners.map((o) => (
              <option key={o.id} value={o.id}>
                {o.name}
              </option>
            ))}
          </select>
        </label>
        <label>
          <span>Amount (units)</span>
          <input
            type="number"
            min={1}
            value={amount}
            onChange={(e) => setAmount(Number(e.target.value))}
            required
          />
        </label>
        <label>
          <span>Due date</span>
          <input
            type="date"
            value={dueAt}
            onChange={(e) => setDueAt(e.target.value)}
            required
          />
        </label>
        <div className="span2 form-actions">
          <button type="submit" className="primary" disabled={busy}>
            {busy ? "Creating…" : "Create obligation"}
          </button>
          {err && <span className="error">{err}</span>}
        </div>
      </form>
      {obligations.length > 0 && (
        <table className="commodity-table" style={{ marginTop: "1rem" }}>
          <thead>
            <tr>
              <th>Creditor</th>
              <th>Debtor</th>
              <th>Amount</th>
              <th>Due</th>
              <th>Status</th>
              <th></th>
            </tr>
          </thead>
          <tbody>
            {obligations.map((ob) => (
              <tr key={ob.id}>
                <td>{ownerById.get(ob.creditor_id)?.name ?? ob.creditor_id.slice(0, 8)}</td>
                <td>{ownerById.get(ob.debtor_id)?.name ?? ob.debtor_id.slice(0, 8)}</td>
                <td className="snlt-cell">{ob.amount}</td>
                <td className="muted small">{ob.due_at.slice(0, 10)}</td>
                <td>
                  {ob.paid_at ? (
                    <span className="muted">settled</span>
                  ) : (
                    <span className="small">outstanding</span>
                  )}
                </td>
                <td>
                  {!ob.paid_at && (
                    <button className="secondary" onClick={() => settle(ob.id)}>
                      Settle
                    </button>
                  )}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </section>
  );
}

function WorldMoneyPanel({
  owners,
  transfers,
  onChanged,
}: {
  owners: Owner[];
  transfers: WorldMoneyTransfer[];
  onChanged: () => void;
}) {
  const [senderID, setSenderID] = useState("");
  const [receiverID, setReceiverID] = useState("");
  const [goldMg, setGoldMg] = useState(0);
  const [busy, setBusy] = useState(false);
  const [err, setErr] = useState<string | null>(null);

  async function submit(e: FormEvent) {
    e.preventDefault();
    setBusy(true);
    setErr(null);
    try {
      const input: CreateWorldMoneyTransferInput = {
        sender_id: senderID,
        receiver_id: receiverID,
        gold_mg: goldMg,
      };
      await api.createWorldMoneyTransfer(input);
      setGoldMg(0);
      onChanged();
    } catch (e) {
      setErr(e instanceof Error ? e.message : String(e));
    } finally {
      setBusy(false);
    }
  }

  if (owners.length === 0) return null;

  const ownerById = new Map(owners.map((o) => [o.id, o]));

  return (
    <section className="card">
      <h2>World money transfers</h2>
      <p className="description">
        "On the world market, money strips off its national uniform and appears in its original form
        as bullion." (Ch. 3 §3C)
      </p>
      <form className="form-grid" onSubmit={submit}>
        <label>
          <span>Sender</span>
          <select value={senderID} onChange={(e) => setSenderID(e.target.value)} required>
            <option value="">—</option>
            {owners.map((o) => (
              <option key={o.id} value={o.id}>
                {o.name}
              </option>
            ))}
          </select>
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
          <span>Gold weight (mg)</span>
          <input
            type="number"
            min={1}
            value={goldMg}
            onChange={(e) => setGoldMg(Number(e.target.value))}
            required
          />
        </label>
        <div className="span2 form-actions">
          <button type="submit" className="primary" disabled={busy}>
            {busy ? "Transferring…" : "Transfer bullion"}
          </button>
          {err && <span className="error">{err}</span>}
        </div>
      </form>
      {transfers.length > 0 && (
        <table className="commodity-table" style={{ marginTop: "1rem" }}>
          <thead>
            <tr>
              <th>Sender</th>
              <th>Receiver</th>
              <th>Gold (mg)</th>
            </tr>
          </thead>
          <tbody>
            {transfers.map((t) => (
              <tr key={t.id}>
                <td>{ownerById.get(t.sender_id)?.name ?? t.sender_id.slice(0, 8)}</td>
                <td>{ownerById.get(t.receiver_id)?.name ?? t.receiver_id.slice(0, 8)}</td>
                <td className="snlt-cell">{t.gold_mg}</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </section>
  );
}
