import { useEffect, useState } from "react";
import { financeApi } from "../../api";
import type { GoldReserve, RateOfExchange } from "../../types";
import "./Ch35PreciousMetal.css";

function penceToPounds(pence: number): string {
  return `£${(pence / 100).toLocaleString("en-GB", {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  })}`;
}

function deviationLabel(bp: number): { text: string; cls: string } {
  if (bp > 0) return { text: `+${bp} bp`, cls: "v3-ch35-deviation-pos" };
  if (bp < 0) return { text: `−${Math.abs(bp)} bp`, cls: "v3-ch35-deviation-neg" };
  return { text: "0 bp (at par)", cls: "v3-ch35-deviation-zero" };
}

export function Ch35PreciousMetal() {
  const [reserves, setReserves] = useState<GoldReserve[]>([]);
  const [reserveListError, setReserveListError] = useState<string | null>(null);
  const [rates, setRates] = useState<RateOfExchange[]>([]);
  const [rateListError, setRateListError] = useState<string | null>(null);

  // Gold Reserve form state
  const [reserveName, setReserveName] = useState("Bank of England Bullion Reserve 1857");
  const [reserveAmount, setReserveAmount] = useState(800_000_000);
  const [reserveDescription, setReserveDescription] = useState(
    "BoE gold reserve at the peak of the 1857 drain; pivot of the credit system",
  );
  const [reserveCreateError, setReserveCreateError] = useState<string | null>(null);
  const [lastReserve, setLastReserve] = useState<GoldReserve | null>(null);

  // Rate of Exchange form state
  const [rateName, setRateName] = useState("Sterling on Hamburg 1847");
  const [deviationBp, setDeviationBp] = useState(-250);
  const [rateDescription, setRateDescription] = useState(
    "Balance-of-payments deficit → pound below par; gold drain to Hamburg",
  );
  const [rateCreateError, setRateCreateError] = useState<string | null>(null);
  const [lastRate, setLastRate] = useState<RateOfExchange | null>(null);

  async function refreshReserves() {
    try {
      const items = await financeApi.listGoldReserves();
      setReserves(items);
    } catch (e) {
      setReserveListError(String(e));
    }
  }

  async function refreshRates() {
    try {
      const items = await financeApi.listRatesOfExchange();
      setRates(items);
    } catch (e) {
      setRateListError(String(e));
    }
  }

  useEffect(() => {
    refreshReserves().catch(() => {});
    refreshRates().catch(() => {});
  }, []);

  async function handleCreateReserve() {
    setReserveCreateError(null);
    setLastReserve(null);
    try {
      const result = await financeApi.createGoldReserve({
        name: reserveName,
        amount: reserveAmount,
        description: reserveDescription,
      });
      setLastReserve(result);
      await refreshReserves();
    } catch (e) {
      setReserveCreateError(String(e));
    }
  }

  async function handleCreateRate() {
    setRateCreateError(null);
    setLastRate(null);
    try {
      const result = await financeApi.createRateOfExchange({
        name: rateName,
        deviation_bp: deviationBp,
        description: rateDescription,
      });
      setLastRate(result);
      await refreshRates();
    } catch (e) {
      setRateCreateError(String(e));
    }
  }

  return (
    <div className="v3-ch35">
      <section className="v3-ch35-intro">
        <p>
          In the world market, <strong>gold functions as world money</strong> — not
          as a national coin but as the universal commodity in its bullion form.
          Every nation-state must maintain a <strong>reserve fund</strong> of gold
          that serves a dual purpose: it is the reserve for domestic payments and the
          pivot of the international credit system. When the reserve falls below a
          critical level, the central bank raises its discount rate to attract gold
          inflows — demonstrating that the rate of interest is the instrument by which
          the gold reserve is defended, not an expression of the &ldquo;price of
          money.&rdquo;
        </p>
        <p>
          The <strong>Bank of England gold reserve</strong> occupied precisely this
          pivot position throughout the mid-nineteenth century. In 1847 and again in
          1857, a drain on the bullion reserve — not a shortage of notes — was the
          trigger for the rate of interest to be driven to 10&ndash;12&nbsp;%. The
          reserve had to be rebuilt before confidence returned. Marx shows that the
          Currency School&rsquo;s mechanical note-issue rule was irrelevant to this
          process: what mattered was gold, not paper.
        </p>
        <p>
          The <strong>rate of exchange</strong> measures the deviation of a
          currency&rsquo;s purchasing power from <em>mint parity</em> — the ratio
          of the gold contents of two national coins. A persistent{" "}
          <strong>balance-of-payments deficit</strong> (outflows for debt
          service, commodity imports, and capital remittances exceeding inflows)
          depresses the exchange below par, making it cheaper to ship gold
          than to buy foreign bills. The rate of exchange is therefore the{" "}
          <strong>thermometer of the balance of payments</strong>, not of the
          balance of trade: a favourable trade balance may coexist with an
          adverse payments balance once interest, dividends, and loan
          repayments are included.
        </p>
        <p>
          The <strong>gold drains to India and China</strong> in the 1850s
          illustrated this mechanism. The East India trade produced a
          persistent payments drain: Britain exported manufactured goods but
          re-exported silver and gold to settle the trade surplus of the
          sub-continent. The result was a persistent sub-par exchange on
          Calcutta and Canton and a corresponding reduction in the BoE bullion
          reserve — quite independently of any excess note circulation.
        </p>
      </section>

      {/* Gold Reserve form */}
      <section className="v3-ch35-section" aria-label="Record a gold reserve">
        <h4>Record a gold reserve</h4>
        <p className="v3-ch35-subintro">
          Enter a central-bank gold stock observation. Amount is in pence
          (divide by 100 for&nbsp;&pound;).
        </p>
        <div className="v3-ch35-form">
          <label>
            Name
            <input value={reserveName} onChange={(e) => setReserveName(e.target.value)} />
          </label>
          <label>
            Amount (pence, &ge; 0)
            <input
              type="number"
              min={0}
              value={reserveAmount}
              onChange={(e) => setReserveAmount(Math.max(0, Number(e.target.value)))}
            />
            <span className="v3-ch35-hint">{penceToPounds(reserveAmount)}</span>
          </label>
          <label>
            Description
            <input
              value={reserveDescription}
              onChange={(e) => setReserveDescription(e.target.value)}
            />
          </label>
        </div>
        <button className="v3-ch35-btn" onClick={handleCreateReserve}>
          Save reserve
        </button>
        {reserveCreateError && <p className="v3-ch35-error">{reserveCreateError}</p>}
        {lastReserve && (
          <div className="v3-ch35-preview">
            Saved &mdash; {lastReserve.name}: {penceToPounds(lastReserve.amount)}
          </div>
        )}
      </section>

      {/* Gold reserves list */}
      <section className="v3-ch35-section" aria-label="Stored gold reserves">
        <h4>Stored gold reserves</h4>
        {reserveListError && <p className="v3-ch35-error">{reserveListError}</p>}
        {reserves.length === 0 ? (
          <p className="v3-ch35-subintro">No reserves stored yet.</p>
        ) : (
          <div className="v3-ch35-table-wrap">
            <table className="v3-ch35-table">
              <thead>
                <tr>
                  <th>Name</th>
                  <th>Amount</th>
                  <th>Description</th>
                  <th>Recorded</th>
                </tr>
              </thead>
              <tbody>
                {reserves.map((r) => (
                  <tr key={r.id}>
                    <td>{r.name}</td>
                    <td>{penceToPounds(r.amount)}</td>
                    <td>{r.description}</td>
                    <td>{new Date(r.created_at).toLocaleDateString("en-GB")}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </section>

      {/* Rate of Exchange form */}
      <section className="v3-ch35-section" aria-label="Record a rate of exchange">
        <h4>Record a rate of exchange</h4>
        <p className="v3-ch35-subintro">
          Deviation from mint parity in basis points (signed integer). Positive =
          above par (inflow); negative = below par (outflow / gold drain).
        </p>
        <div className="v3-ch35-form">
          <label>
            Name
            <input value={rateName} onChange={(e) => setRateName(e.target.value)} />
          </label>
          <label>
            Deviation from par (basis points, signed)
            <input
              type="number"
              value={deviationBp}
              onChange={(e) => setDeviationBp(Number(e.target.value))}
            />
            <span className="v3-ch35-hint">
              {deviationBp === 0
                ? "At mint parity"
                : deviationBp > 0
                ? `${deviationBp} bp above par — gold flows in`
                : `${Math.abs(deviationBp)} bp below par — gold drains out`}
            </span>
          </label>
          <label>
            Description
            <input
              value={rateDescription}
              onChange={(e) => setRateDescription(e.target.value)}
            />
          </label>
        </div>
        <button className="v3-ch35-btn" onClick={handleCreateRate}>
          Save rate
        </button>
        {rateCreateError && <p className="v3-ch35-error">{rateCreateError}</p>}
        {lastRate && (
          <div className="v3-ch35-preview">
            Saved &mdash; {lastRate.name}:{" "}
            {(() => {
              const d = deviationLabel(lastRate.deviation_bp);
              return <span className={d.cls}>{d.text}</span>;
            })()}
          </div>
        )}
      </section>

      {/* Rates of exchange list */}
      <section className="v3-ch35-section" aria-label="Stored rates of exchange">
        <h4>Stored rates of exchange</h4>
        {rateListError && <p className="v3-ch35-error">{rateListError}</p>}
        {rates.length === 0 ? (
          <p className="v3-ch35-subintro">No rates stored yet.</p>
        ) : (
          <div className="v3-ch35-table-wrap">
            <table className="v3-ch35-table">
              <thead>
                <tr>
                  <th>Name</th>
                  <th>Deviation from par</th>
                  <th>Description</th>
                  <th>Recorded</th>
                </tr>
              </thead>
              <tbody>
                {rates.map((r) => {
                  const d = deviationLabel(r.deviation_bp);
                  return (
                    <tr key={r.id}>
                      <td>{r.name}</td>
                      <td>
                        <span className={d.cls}>{d.text}</span>
                      </td>
                      <td>{r.description}</td>
                      <td>{new Date(r.created_at).toLocaleDateString("en-GB")}</td>
                    </tr>
                  );
                })}
              </tbody>
            </table>
          </div>
        )}
      </section>
    </div>
  );
}
