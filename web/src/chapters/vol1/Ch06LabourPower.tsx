import { useEffect, useState } from "react";
import type { FormEvent } from "react";
import { api } from "../../api";
import type {
  LabourWorker,
  LabourCapitalist,
  LabourPowerOffering,
  LabourPowerPurchase,
} from "../../types";

interface Ch06Props {
  onSharedChanged: () => void;
}

function minutesToHours(m: number): string {
  const h = Math.floor(m / 60);
  const min = m % 60;
  return min === 0 ? `${h}h` : `${h}h ${min}m`;
}

export function Ch06LabourPower({ onSharedChanged: _onSharedChanged }: Ch06Props) {
  const [workers, setWorkers] = useState<LabourWorker[]>([]);
  const [capitalists, setCapitalists] = useState<LabourCapitalist[]>([]);
  const [offerings, setOfferings] = useState<LabourPowerOffering[]>([]);
  const [purchases, setPurchases] = useState<LabourPowerPurchase[]>([]);
  const [loadErr, setLoadErr] = useState<string | null>(null);

  async function refresh() {
    try {
      const [ws, cs, os, ps] = await Promise.all([
        api.listLabourWorkers(),
        api.listLabourCapitalists(),
        api.listOfferings(),
        api.listLabourPurchases(),
      ]);
      setWorkers(ws);
      setCapitalists(cs);
      setOfferings(os);
      setPurchases(ps);
    } catch (e) {
      setLoadErr(e instanceof Error ? e.message : String(e));
    }
  }

  useEffect(() => { refresh(); }, []);

  return (
    <>
      {loadErr && <p className="error">{loadErr}</p>}
      <RegisterWorkerPanel onCreated={refresh} />
      <RegisterCapitalistPanel onCreated={refresh} />
      <WorkerListPanel workers={workers} />
      <CapitalistListPanel capitalists={capitalists} />
      <PostOfferingPanel workers={workers} onCreated={refresh} />
      <ActiveOfferingsPanel offerings={offerings} workers={workers} />
      <PurchasePanel
        workers={workers}
        capitalists={capitalists}
        onCreated={refresh}
      />
      <PurchaseHistoryPanel purchases={purchases} workers={workers} capitalists={capitalists} />
    </>
  );
}

function RegisterWorkerPanel({ onCreated }: { onCreated: () => void }) {
  const [capacity, setCapacity] = useState(480);
  const [ownsLP, setOwnsLP] = useState(true);
  const [ownsCommodities, setOwnsCommodities] = useState(false);
  const [err, setErr] = useState<string | null>(null);

  async function submit(e: FormEvent) {
    e.preventDefault();
    setErr(null);
    try {
      await api.createLabourWorker({
        owns_labour_power: ownsLP,
        owns_commodities_to_sell: ownsCommodities,
        capacity_minutes_per_day: capacity,
      });
      onCreated();
    } catch (e) {
      setErr(e instanceof Error ? e.message : String(e));
    }
  }

  return (
    <section className="card">
      <h2>Register Worker</h2>
      <p className="description">
        A labourer who owns their capacity for labour and is obliged to sell it.
      </p>
      <form className="form-grid" onSubmit={submit}>
        <label>
          <span>Capacity (minutes/day)</span>
          <input
            type="number"
            value={capacity}
            onChange={(e) => setCapacity(Number(e.target.value))}
            min={1}
          />
        </label>
        <label>
          <span>Owns labour-power</span>
          <input type="checkbox" checked={ownsLP} onChange={(e) => setOwnsLP(e.target.checked)} />
        </label>
        <label>
          <span>Owns commodities to sell</span>
          <input
            type="checkbox"
            checked={ownsCommodities}
            onChange={(e) => setOwnsCommodities(e.target.checked)}
          />
        </label>
        <div className="form-actions span2">
          <button type="submit" className="primary">Register</button>
          {err && <span className="error">{err}</span>}
        </div>
      </form>
    </section>
  );
}

function RegisterCapitalistPanel({ onCreated }: { onCreated: () => void }) {
  const [moneyCapital, setMoneyCapital] = useState(480);
  const [err, setErr] = useState<string | null>(null);

  async function submit(e: FormEvent) {
    e.preventDefault();
    setErr(null);
    try {
      await api.createLabourCapitalist({ money_capital: moneyCapital });
      onCreated();
    } catch (e) {
      setErr(e instanceof Error ? e.message : String(e));
    }
  }

  return (
    <section className="card">
      <h2>Register Capitalist</h2>
      <p className="description">
        The owner of money-capital who purchases labour-power as a commodity.
      </p>
      <form className="form-grid" onSubmit={submit}>
        <label>
          <span>Money capital (labour-minutes)</span>
          <input
            type="number"
            value={moneyCapital}
            onChange={(e) => setMoneyCapital(Number(e.target.value))}
            min={0}
          />
        </label>
        <div className="form-actions span2">
          <button type="submit" className="primary">Register</button>
          {err && <span className="error">{err}</span>}
        </div>
      </form>
    </section>
  );
}

function WorkerListPanel({ workers }: { workers: LabourWorker[] }) {
  if (workers.length === 0) return null;
  return (
    <section className="card">
      <h2>Workers</h2>
      <div className="item-list">
        {workers.map((w) => (
          <div key={w.id} className="item-card">
            <div className="item-header">
              <span className="item-name monospace small">{w.id.slice(0, 8)}&hellip;</span>
              <span className="item-meta">
                {minutesToHours(w.labour_power.capacity_minutes_per_day)} / day
              </span>
              {w.owns_labour_power && !w.owns_commodities_to_sell && (
                <span className="item-tag">free labourer</span>
              )}
            </div>
            <p className="small muted">
              Owns LP: {w.owns_labour_power ? "yes" : "no"} &middot;
              Owns commodities: {w.owns_commodities_to_sell ? "yes" : "no"}
            </p>
          </div>
        ))}
      </div>
    </section>
  );
}

function CapitalistListPanel({ capitalists }: { capitalists: LabourCapitalist[] }) {
  if (capitalists.length === 0) return null;
  return (
    <section className="card">
      <h2>Capitalists</h2>
      <div className="item-list">
        {capitalists.map((c) => (
          <div key={c.id} className="item-card">
            <div className="item-header">
              <span className="item-name monospace small">{c.id.slice(0, 8)}&hellip;</span>
              <span className="item-meta">{minutesToHours(c.money_capital)} capital</span>
            </div>
          </div>
        ))}
      </div>
    </section>
  );
}

function PostOfferingPanel({
  workers,
  onCreated,
}: {
  workers: LabourWorker[];
  onCreated: () => void;
}) {
  const [ownerID, setOwnerID] = useState("");
  const [capacity, setCapacity] = useState(480);
  const [contractDays, setContractDays] = useState(5);
  const [askingWage, setAskingWage] = useState(240);
  const [err, setErr] = useState<string | null>(null);

  async function submit(e: FormEvent) {
    e.preventDefault();
    setErr(null);
    try {
      await api.createOffering({
        owner_id: ownerID,
        capacity_minutes_per_day: capacity,
        contract_days: contractDays,
        asking_wage: askingWage,
      });
      onCreated();
    } catch (e) {
      setErr(e instanceof Error ? e.message : String(e));
    }
  }

  return (
    <section className="card">
      <h2>Post Labour-Power for Sale</h2>
      <p className="description">
        A worker offers their capacity for labour for a finite period at an asking wage.
      </p>
      <form className="form-grid" onSubmit={submit}>
        <label>
          <span>Worker</span>
          <select value={ownerID} onChange={(e) => setOwnerID(e.target.value)} required>
            <option value="">Select a worker…</option>
            {workers.map((w) => (
              <option key={w.id} value={w.id}>
                {w.id.slice(0, 8)}… ({minutesToHours(w.labour_power.capacity_minutes_per_day)}/day)
              </option>
            ))}
          </select>
        </label>
        <label>
          <span>Capacity (min/day)</span>
          <input
            type="number"
            value={capacity}
            onChange={(e) => setCapacity(Number(e.target.value))}
            min={1}
          />
        </label>
        <label>
          <span>Contract (days)</span>
          <input
            type="number"
            value={contractDays}
            onChange={(e) => setContractDays(Number(e.target.value))}
            min={1}
          />
        </label>
        <label>
          <span>Asking wage (labour-min/day)</span>
          <input
            type="number"
            value={askingWage}
            onChange={(e) => setAskingWage(Number(e.target.value))}
            min={0}
          />
        </label>
        <div className="form-actions span2">
          <button type="submit" className="primary">Post offering</button>
          {err && <span className="error">{err}</span>}
        </div>
      </form>
    </section>
  );
}

function ActiveOfferingsPanel({
  offerings,
  workers,
}: {
  offerings: LabourPowerOffering[];
  workers: LabourWorker[];
}) {
  if (offerings.length === 0) return null;
  const workerMap = new Map(workers.map((w) => [w.id, w]));
  return (
    <section className="card">
      <h2>Active Offerings</h2>
      <table className="data-table">
        <thead>
          <tr>
            <th>Worker</th>
            <th>Capacity</th>
            <th>Days</th>
            <th>Asking wage / day</th>
          </tr>
        </thead>
        <tbody>
          {offerings.map((o) => {
            const w = workerMap.get(o.owner_id);
            return (
              <tr key={o.id}>
                <td className="monospace small">
                  {w ? `${o.owner_id.slice(0, 8)}…` : o.owner_id.slice(0, 8) + "…"}
                </td>
                <td>{minutesToHours(o.capacity_minutes_per_day)}</td>
                <td>{o.contract_days}</td>
                <td>{minutesToHours(o.asking_wage)}</td>
              </tr>
            );
          })}
        </tbody>
      </table>
    </section>
  );
}

function PurchasePanel({
  workers,
  capitalists,
  onCreated,
}: {
  workers: LabourWorker[];
  capitalists: LabourCapitalist[];
  onCreated: () => void;
}) {
  const [sellerID, setSellerID] = useState("");
  const [buyerID, setBuyerID] = useState("");
  const [wageMinutes, setWageMinutes] = useState(240);
  const [contractDays, setContractDays] = useState(5);
  const [err, setErr] = useState<string | null>(null);

  async function submit(e: FormEvent) {
    e.preventDefault();
    setErr(null);
    try {
      await api.createLabourPurchase({
        seller_id: sellerID,
        buyer_id: buyerID,
        wage_minutes: wageMinutes,
        contract_days: contractDays,
      });
      onCreated();
    } catch (e) {
      setErr(e instanceof Error ? e.message : String(e));
    }
  }

  return (
    <section className="card">
      <h2>Purchase Labour-Power</h2>
      <p className="description">
        The capitalist buys labour-power as a commodity. The price equals its value when
        wage = daily subsistence cost; surplus arises only in production (deferred to Ch. 7).
      </p>
      <form className="form-grid" onSubmit={submit}>
        <label>
          <span>Seller (worker)</span>
          <select value={sellerID} onChange={(e) => setSellerID(e.target.value)} required>
            <option value="">Select a worker…</option>
            {workers.map((w) => (
              <option key={w.id} value={w.id}>{w.id.slice(0, 8)}…</option>
            ))}
          </select>
        </label>
        <label>
          <span>Buyer (capitalist)</span>
          <select value={buyerID} onChange={(e) => setBuyerID(e.target.value)} required>
            <option value="">Select a capitalist…</option>
            {capitalists.map((c) => (
              <option key={c.id} value={c.id}>{c.id.slice(0, 8)}…</option>
            ))}
          </select>
        </label>
        <label>
          <span>Wage (labour-min/day)</span>
          <input
            type="number"
            value={wageMinutes}
            onChange={(e) => setWageMinutes(Number(e.target.value))}
            min={0}
          />
        </label>
        <label>
          <span>Contract (days)</span>
          <input
            type="number"
            value={contractDays}
            onChange={(e) => setContractDays(Number(e.target.value))}
            min={1}
          />
        </label>
        <div className="form-actions span2">
          <button type="submit" className="primary">Purchase</button>
          {err && <span className="error">{err}</span>}
        </div>
      </form>
    </section>
  );
}

function PurchaseHistoryPanel({
  purchases,
  workers,
  capitalists,
}: {
  purchases: LabourPowerPurchase[];
  workers: LabourWorker[];
  capitalists: LabourCapitalist[];
}) {
  if (purchases.length === 0) return null;
  const workerMap = new Map(workers.map((w) => [w.id, w]));
  const capitalistMap = new Map(capitalists.map((c) => [c.id, c]));

  return (
    <section className="card">
      <h2>Purchase History</h2>
      <table className="data-table">
        <thead>
          <tr>
            <th>Seller</th>
            <th>Buyer</th>
            <th>Wage / day</th>
            <th>Days</th>
            <th>Total wage</th>
          </tr>
        </thead>
        <tbody>
          {purchases.map((p) => {
            const seller = workerMap.get(p.seller_id);
            const buyer = capitalistMap.get(p.buyer_id);
            return (
              <tr key={p.id}>
                <td className="monospace small">
                  {seller ? p.seller_id.slice(0, 8) + "…" : p.seller_id.slice(0, 8) + "…"}
                </td>
                <td className="monospace small">
                  {buyer ? p.buyer_id.slice(0, 8) + "…" : p.buyer_id.slice(0, 8) + "…"}
                </td>
                <td>{minutesToHours(p.wage_minutes)}</td>
                <td>{p.contract_days}</td>
                <td>{minutesToHours(p.wage_minutes * p.contract_days)}</td>
              </tr>
            );
          })}
        </tbody>
      </table>
    </section>
  );
}
