import { useEffect, useState } from "react";
import { financeApi } from "../../api";
import type { CurrencyObservation } from "../../types";
import "./Ch28MediumOfCirculation.css";

function penceToMillions(pence: number): string {
  return `£${(pence / 100 / 1_000_000).toFixed(1)}m`;
}

function bpToPercent(bp: number): string {
  return `${(bp / 100).toFixed(1)}%`;
}

function coinShareBP(obs: CurrencyObservation): number {
  if (obs.total_currency_outstanding === 0) return 0;
  return Math.floor((obs.coin_function_amount * 10000) / obs.total_currency_outstanding);
}

function classifyCurrency(obs: CurrencyObservation): string {
  const bp = coinShareBP(obs);
  if (bp >= 7500) return "Coin function dominant (revenue expenditure / retail)";
  if (bp >= 2500) return "Mixed: both revenue expenditure and capital transfer";
  return "Capital-transfer function dominant (productive-capital movement)";
}

export function Ch28MediumOfCirculation() {
  const [observations, setObservations] = useState<CurrencyObservation[]>([]);
  const [listError, setListError] = useState<string | null>(null);

  const [name, setName] = useState("Bank of England 1844");
  const [total, setTotal] = useState(2_000_000_000);
  const [coinAmount, setCoinAmount] = useState(1_400_000_000);
  const [capitalAmount, setCapitalAmount] = useState(400_000_000);
  const [reserveAmount, setReserveAmount] = useState(800_000_000);
  const [reserveDesc, setReserveDesc] = useState("Banking Department reserve 1844");
  const [description, setDescription] = useState("");
  const [createError, setCreateError] = useState<string | null>(null);
  const [splitError, setSplitError] = useState<string | null>(null);

  async function refresh() {
    try {
      const items = await financeApi.listCurrencyObservations();
      setObservations(items);
    } catch (e) {
      setListError(String(e));
    }
  }

  useEffect(() => {
    refresh().catch(() => {});
  }, []);

  function validate(): boolean {
    if (total <= 0) {
      setSplitError("Total currency outstanding must be > 0.");
      return false;
    }
    if (coinAmount + capitalAmount > total) {
      setSplitError("Coin + capital-transfer amounts must not exceed total outstanding.");
      return false;
    }
    if (reserveAmount < 0) {
      setSplitError("Reserve fund amount must be >= 0.");
      return false;
    }
    setSplitError(null);
    return true;
  }

  async function handleCreate() {
    if (!validate()) return;
    setCreateError(null);
    try {
      await financeApi.createCurrencyObservation({
        name,
        total_currency_outstanding: total,
        coin_function_amount: coinAmount,
        capital_transfer_amount: capitalAmount,
        reserve_amount: reserveAmount,
        reserve_description: reserveDesc,
        description,
      });
      await refresh();
    } catch (e) {
      setCreateError(String(e));
    }
  }

  const coinPct = total > 0 ? ((coinAmount / total) * 100).toFixed(1) : "—";
  const capitalPct = total > 0 ? ((capitalAmount / total) * 100).toFixed(1) : "—";

  return (
    <div className="v3-ch28">
      <section className="v3-ch28-intro">
        <p>
          Tooke and Fullarton conflate three distinct categories that Marx separates:
        </p>
        <ol className="v3-ch28-list">
          <li>
            <strong>Coin function</strong> — money as medium of circulation for{" "}
            <em>revenue expenditure</em>: wages paid to workers, retail purchases,
            consumer spending. The volume of notes required here is governed by the
            velocity and price level of retail trade.
          </li>
          <li>
            <strong>Capital-transfer function</strong> — money as the vehicle for
            transferring <em>productive capital</em> between capitalists. These
            transfers are largely settled by bills of exchange and book credits;
            they do not require a net addition to the note circulation.
          </li>
          <li>
            <strong>Interest-bearing capital</strong> — the distinct M—M′ circuit,
            where money begets money without passing through production at all.
          </li>
        </ol>
        <p>
          By collapsing (1) and (2), Tooke and Fullarton cannot explain why the
          quantity of notes in circulation moves with the rhythm of retail trade
          rather than with the volume of inter-capitalist capital transfers.
          The Bank of England 1844 data — £20m notes outstanding, £8m reserve —
          illustrates the split.
        </p>
      </section>

      {/* Seed exemplar card — Bank of England 1844 */}
      <section className="v3-ch28-section" aria-label="Bank of England 1844 exemplar">
        <h4>Marx's exemplar: Bank of England 1844</h4>
        <div className="v3-ch28-seed-card">
          <div className="v3-ch28-seed-row">
            <span className="v3-ch28-seed-label">Notes outstanding</span>
            <span className="v3-ch28-seed-val">£20m (2,000,000,000 pence)</span>
          </div>
          <div className="v3-ch28-seed-row">
            <span className="v3-ch28-seed-label">Banking Department reserve</span>
            <span className="v3-ch28-seed-val">£8m (800,000,000 pence)</span>
          </div>
          <div className="v3-ch28-seed-row">
            <span className="v3-ch28-seed-label">Coin function (retail/wages)</span>
            <span className="v3-ch28-seed-val">~£14m — 70% of total (7000 bp)</span>
          </div>
          <div className="v3-ch28-seed-row">
            <span className="v3-ch28-seed-label">Capital-transfer function</span>
            <span className="v3-ch28-seed-val">~£4m — settled via bills/book credit</span>
          </div>
          <div className="v3-ch28-seed-verdict">
            Verdict: predominantly <strong>coin function</strong> — currency serves
            mainly revenue expenditure.
          </div>
        </div>
      </section>

      {/* Create form */}
      <section className="v3-ch28-section" aria-label="Record a currency observation">
        <h4>Record a currency observation</h4>
        <p className="v3-ch28-subintro">
          Enter total notes outstanding and classify the split between coin function
          (retail/revenue) and capital-transfer function. The server validates that
          coin + capital-transfer ≤ total.
        </p>
        <div className="v3-ch28-form">
          <label>
            Name
            <input value={name} onChange={(e) => setName(e.target.value)} />
          </label>
          <label>
            Total currency outstanding (pence)
            <input
              type="number"
              min={1}
              value={total}
              onChange={(e) => setTotal(Number(e.target.value))}
            />
            <span className="v3-ch28-hint">{penceToMillions(total)}</span>
          </label>
          <label>
            Coin function amount — revenue expenditure (pence)
            <input
              type="number"
              min={0}
              value={coinAmount}
              onChange={(e) => setCoinAmount(Number(e.target.value))}
            />
            <span className="v3-ch28-hint">{penceToMillions(coinAmount)} — {coinPct}% of total</span>
          </label>
          <label>
            Capital-transfer amount — productive-capital movement (pence)
            <input
              type="number"
              min={0}
              value={capitalAmount}
              onChange={(e) => setCapitalAmount(Number(e.target.value))}
            />
            <span className="v3-ch28-hint">{penceToMillions(capitalAmount)} — {capitalPct}% of total</span>
          </label>
          <label>
            Reserve fund amount (pence, &ge; 0)
            <input
              type="number"
              min={0}
              value={reserveAmount}
              onChange={(e) => setReserveAmount(Number(e.target.value))}
            />
            <span className="v3-ch28-hint">{penceToMillions(reserveAmount)} held against demand</span>
          </label>
          <label>
            Reserve fund description
            <input value={reserveDesc} onChange={(e) => setReserveDesc(e.target.value)} />
          </label>
          <label>
            Description
            <input value={description} onChange={(e) => setDescription(e.target.value)} />
          </label>
        </div>
        {splitError && <p className="v3-ch28-error">{splitError}</p>}
        <button className="v3-ch28-btn" onClick={handleCreate}>
          Save observation
        </button>
        {createError && <p className="v3-ch28-error">{createError}</p>}
      </section>

      {/* Stored observations */}
      <section className="v3-ch28-section" aria-label="Stored currency observations">
        <h4>Stored observations</h4>
        {listError && <p className="v3-ch28-error">{listError}</p>}
        <div className="v3-ch28-cards">
          {observations.length === 0 && (
            <p className="v3-ch28-subintro">No observations stored yet.</p>
          )}
          {observations.map((o) => (
            <div className="v3-ch28-card" key={o.id}>
              <div className="v3-ch28-card-name">{o.name}</div>
              <div className="v3-ch28-card-rows">
                <div className="v3-ch28-card-row">
                  <span>Total outstanding</span>
                  <span>{penceToMillions(o.total_currency_outstanding)}</span>
                </div>
                <div className="v3-ch28-card-row">
                  <span>Coin function</span>
                  <span>
                    {penceToMillions(o.coin_function_amount)}
                    {" — "}
                    {bpToPercent(coinShareBP(o))} of total
                  </span>
                </div>
                <div className="v3-ch28-card-row">
                  <span>Capital-transfer</span>
                  <span>{penceToMillions(o.capital_transfer_amount)}</span>
                </div>
                <div className="v3-ch28-card-row">
                  <span>Reserve</span>
                  <span>
                    {penceToMillions(o.reserve_fund.amount)}
                    {o.reserve_fund.description ? ` — ${o.reserve_fund.description}` : ""}
                  </span>
                </div>
              </div>
              <div className="v3-ch28-verdict">{classifyCurrency(o)}</div>
              {o.description && (
                <div className="v3-ch28-desc">{o.description}</div>
              )}
            </div>
          ))}
        </div>
      </section>
    </div>
  );
}
