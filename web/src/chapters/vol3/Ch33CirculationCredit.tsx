import { useEffect, useState } from "react";
import { financeApi } from "../../api";
import type { ClearingHouseSettlement } from "../../types";
import "./Ch33CirculationCredit.css";

function penceToPounds(pence: number): string {
  return `£${(pence / 100).toLocaleString("en-GB", {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  })}`;
}

function moneySavedPct(totalClaims: number, moneyUsed: number): string {
  if (totalClaims <= 0) return "0.0%";
  return `${(((totalClaims - moneyUsed) / totalClaims) * 100).toFixed(1)}%`;
}

export function Ch33CirculationCredit() {
  const [records, setRecords] = useState<ClearingHouseSettlement[]>([]);
  const [listError, setListError] = useState<string | null>(null);

  // Velocity calculator state
  const [actualNotes, setActualNotes] = useState(50_000);
  const [velocity, setVelocity] = useState(5);

  // Settlement form state
  const [name, setName] = useState("London Clearing House 1839");
  const [description, setDescription] = useState(
    "mutual claims among London banks netted daily; only the balance settled in coin",
  );
  const [totalClaims, setTotalClaims] = useState(100_000_000);
  const [moneyUsed, setMoneyUsed] = useState(20_000_000);
  const [createError, setCreateError] = useState<string | null>(null);
  const [formError, setFormError] = useState<string | null>(null);
  const [lastCreated, setLastCreated] = useState<ClearingHouseSettlement | null>(null);

  async function refresh() {
    try {
      const items = await financeApi.listClearingHouseSettlements();
      setRecords(items);
    } catch (e) {
      setListError(String(e));
    }
  }

  useEffect(() => {
    refresh().catch(() => {});
  }, []);

  function validate(): boolean {
    if (totalClaims <= 0) {
      setFormError("Total claims must be > 0.");
      return false;
    }
    if (moneyUsed < 0) {
      setFormError("Money used must be >= 0.");
      return false;
    }
    if (moneyUsed > totalClaims) {
      setFormError("Money used cannot exceed total claims.");
      return false;
    }
    setFormError(null);
    return true;
  }

  async function handleCreate() {
    if (!validate()) return;
    setCreateError(null);
    setLastCreated(null);
    try {
      const result = await financeApi.createClearingHouseSettlement({
        name,
        description,
        total_claims: totalClaims,
        money_used: moneyUsed,
      });
      setLastCreated(result);
      await refresh();
    } catch (e) {
      setCreateError(String(e));
    }
  }

  const effectiveMoneySupply = actualNotes * velocity;

  return (
    <div className="v3-ch33">
      <section className="v3-ch33-intro">
        <p>
          Credit <strong>governs currency velocity</strong>: the same notes circulate
          far faster than if each commodity transaction required its own coin. A
          bill of exchange passes through many hands before maturity — each transfer
          is a payment accomplished <em>without</em> money. The mass of actual
          currency in circulation is therefore a <strong>fraction</strong> of the
          value of commodities being circulated.
        </p>
        <p>
          The <strong>London Clearing House</strong> makes this most visible. Banks
          present mutual claims against one another every afternoon; the claims are
          netted and only the <em>net balance</em> is settled in coin. In Marx's
          time, claims of £3–4 million were regularly settled with less than
          £200,000 of actual money — a saving of over 90 per cent.
        </p>
        <p>
          The conclusion is that <strong>money as medium of circulation is
          minimised by credit</strong>. What looks like a vast monetary economy
          rests on a thin substrate of actual coin, replaced almost entirely by
          mutual cancellation of claims.
        </p>
      </section>

      {/* Currency-velocity calculator */}
      <section className="v3-ch33-section" aria-label="Currency velocity calculator">
        <h4>Currency-velocity calculator</h4>
        <p className="v3-ch33-subintro">
          EffectiveMoneySupply = ActualNotes &times; Velocity. Marx's fixture:
          £500 of notes circulating 5 times = £2,500 of purchasing power.
        </p>
        <div className="v3-ch33-calc">
          <label>
            Actual notes in circulation (pence)
            <input
              type="number"
              min={0}
              value={actualNotes}
              onChange={(e) => setActualNotes(Number(e.target.value))}
            />
            <span className="v3-ch33-hint">{penceToPounds(actualNotes)}</span>
          </label>
          <label>
            Velocity (times per period)
            <input
              type="number"
              min={1}
              value={velocity}
              onChange={(e) => setVelocity(Math.max(1, Number(e.target.value)))}
            />
          </label>
        </div>
        <div className="v3-ch33-result">
          Effective money supply: {penceToPounds(effectiveMoneySupply)}
        </div>
        <p className="v3-ch33-fixture">
          Marx's fixture: £500 &times; 5 = £2,500
          {actualNotes === 50_000 && velocity === 5
            ? " — current inputs match the fixture."
            : "."}
        </p>
      </section>

      {/* Clearing-house settlement form */}
      <section className="v3-ch33-section" aria-label="Record a clearing-house settlement">
        <h4>Record a clearing-house settlement</h4>
        <p className="v3-ch33-subintro">
          Enter the gross mutual claims presented and the money actually used to
          settle the net balance. The server computes net_balance = total_claims
          &minus; money_used and enforces money_used &le; total_claims.
        </p>
        <div className="v3-ch33-form">
          <label>
            Name
            <input value={name} onChange={(e) => setName(e.target.value)} />
          </label>
          <label>
            Description
            <input
              value={description}
              onChange={(e) => setDescription(e.target.value)}
            />
          </label>
          <label>
            Total claims (pence, &gt; 0)
            <input
              type="number"
              min={1}
              value={totalClaims}
              onChange={(e) => setTotalClaims(Number(e.target.value))}
            />
            <span className="v3-ch33-hint">{penceToPounds(totalClaims)}</span>
          </label>
          <label>
            Money used to settle (pence, &ge; 0)
            <input
              type="number"
              min={0}
              value={moneyUsed}
              onChange={(e) => setMoneyUsed(Number(e.target.value))}
            />
            <span className="v3-ch33-hint">{penceToPounds(moneyUsed)}</span>
          </label>
        </div>
        <div className="v3-ch33-preview">
          Money saved: {moneySavedPct(totalClaims, moneyUsed)} of claims settled
          without coin
        </div>
        {formError && <p className="v3-ch33-error">{formError}</p>}
        <button className="v3-ch33-btn" onClick={handleCreate}>
          Save settlement
        </button>
        {createError && <p className="v3-ch33-error">{createError}</p>}
        {lastCreated && (
          <div className="v3-ch33-preview">
            Saved — net balance: {penceToPounds(lastCreated.net_balance)} &nbsp;|&nbsp;
            money saved: {moneySavedPct(lastCreated.total_claims, lastCreated.money_used)}
          </div>
        )}
      </section>

      {/* Stored settlements */}
      <section className="v3-ch33-section" aria-label="Stored clearing-house settlements">
        <h4>Stored settlements</h4>
        {listError && <p className="v3-ch33-error">{listError}</p>}
        <div className="v3-ch33-cards">
          {records.length === 0 && (
            <p className="v3-ch33-subintro">No settlements stored yet.</p>
          )}
          {records.map((s) => (
            <div className="v3-ch33-card" key={s.id}>
              <div className="v3-ch33-card-name">{s.name}</div>
              <div className="v3-ch33-card-rows">
                <div className="v3-ch33-card-row">
                  <span>Total claims</span>
                  <span>{penceToPounds(s.total_claims)}</span>
                </div>
                <div className="v3-ch33-card-row">
                  <span>Money used</span>
                  <span>{penceToPounds(s.money_used)}</span>
                </div>
                <div className="v3-ch33-card-row">
                  <span>Net balance</span>
                  <span>{penceToPounds(s.net_balance)}</span>
                </div>
              </div>
              <div className="v3-ch33-saving">
                {moneySavedPct(s.total_claims, s.money_used)} of claims cancelled
                without coin
              </div>
              {s.description && (
                <div className="v3-ch33-desc">{s.description}</div>
              )}
            </div>
          ))}
        </div>
      </section>
    </div>
  );
}
