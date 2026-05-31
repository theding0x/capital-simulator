import { useEffect, useState } from "react";
import { financeApi } from "../../api";
import type { FloatingCapital } from "../../types";
import "./Ch31MoneyRealCapitalII.css";

const SOURCES: { value: number; label: string }[] = [
  { value: 1, label: "1 — Simple transformation of idle money (inverse to accumulation)" },
  { value: 2, label: "2 — Conversion of capital/revenue (accompanies accumulation)" },
];

function sourceLabel(source: number): string {
  return SOURCES.find((s) => s.value === source)?.label ?? `Source ${source}`;
}

function penceToPounds(pence: number): string {
  return `£${(pence / 100).toLocaleString("en-GB", { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`;
}

// Mirror of (FloatingCapital).RelationToAccumulation() from the Go domain.
function relationToAccumulation(source: number): string {
  switch (source) {
    case 1:
      return "inverse to real accumulation — idle money transformed into loan capital swells when production stagnates";
    case 2:
      return "accompanies real accumulation — capital and revenue converted into money on the way to fresh investment";
    default:
      return "unknown source";
  }
}

export function Ch31MoneyRealCapitalII() {
  const [records, setRecords] = useState<FloatingCapital[]>([]);
  const [listError, setListError] = useState<string | null>(null);

  // Velocity calculator state.
  const [velMoney, setVelMoney] = useState(2000);
  const [velTimes, setVelTimes] = useState(5);
  const loanCapitalCreated = velTimes >= 1 ? velMoney * velTimes : 0;

  // Floating-capital form state.
  const [name, setName] = useState("Ipswich farmers' surplus deposits");
  const [amount, setAmount] = useState(2_000_000);
  const [source, setSource] = useState(2);
  const [reserveSaving, setReserveSaving] = useState(50_000);
  const [reserveSavingDesc, setReserveSavingDesc] = useState("country-bank till economies");
  const [description, setDescription] = useState(
    "rural surplus rediscounted by Lombard Street bill-brokers to the manufacturing centres",
  );
  const [createError, setCreateError] = useState<string | null>(null);
  const [formError, setFormError] = useState<string | null>(null);

  async function refresh() {
    try {
      const items = await financeApi.listFloatingCapitals();
      setRecords(items);
    } catch (e) {
      setListError(String(e));
    }
  }

  useEffect(() => {
    refresh().catch(() => {});
  }, []);

  function validate(): boolean {
    if (amount < 0) {
      setFormError("Amount must be >= 0.");
      return false;
    }
    if (velTimes < 1) {
      setFormError("Times lent must be >= 1.");
      return false;
    }
    if (reserveSaving < 0) {
      setFormError("Reserve saving must be >= 0.");
      return false;
    }
    setFormError(null);
    return true;
  }

  async function handleCreate() {
    if (!validate()) return;
    setCreateError(null);
    try {
      await financeApi.createFloatingCapital({
        name,
        amount,
        source,
        velocity_money_amount: velMoney,
        velocity_times_lent: velTimes,
        reserve_saving_amount: reserveSaving,
        reserve_saving_description: reserveSavingDesc,
        description,
      });
      await refresh();
    } catch (e) {
      setCreateError(String(e));
    }
  }

  const previewRelation = relationToAccumulation(source);

  return (
    <div className="v3-ch31">
      <section className="v3-ch31-intro">
        <p>
          Marx names two sources of <strong>loanable money-capital</strong>. The
          first is the <strong>simple transformation of money into loan capital</strong>:
          industrial capital that has fallen idle — money lying fallow after a
          crisis — plus the economies in the reserve of means of circulation that
          banking concentration pools. This source moves <em>inversely</em> to real
          accumulation: it swells precisely when production stagnates and capital
          cannot find productive employment.
        </p>
        <p>
          The second is the <strong>transformation of capital or revenue into money</strong>{" "}
          that then becomes loan capital. This <em>accompanies</em> real accumulation:
          as the reproduction process expands, a growing stream of surplus-value and
          revenue is momentarily converted into money on its way to fresh investment,
          and a part of it is lent.
        </p>
        <p>
          <strong>Bill-brokers</strong> on <strong>Lombard Street</strong> redistribute
          the surplus deposits of agricultural districts — money idle in country banks
          — to the manufacturing centres by <strong>rediscounting</strong> bills. And
          the mass of loanable capital is <em>not</em> the quantity of currency: the
          same <strong>£20 lent five times over performs £100</strong> of loan-capital
          service. Velocity multiplies a fixed money mass into a far larger mass of
          loan capital.
        </p>
      </section>

      {/* Loan-capital velocity calculator */}
      <section className="v3-ch31-section" aria-label="Loan-capital velocity calculator">
        <h4>Loan-capital velocity: the same £20 lent five times = £100</h4>
        <p className="v3-ch31-subintro">
          A fixed stock of money, lent and re-lent, performs a multiple of itself as
          loan capital. loan capital created = money amount × times lent.
        </p>
        <div className="v3-ch31-calc">
          <label>
            Money amount (pence)
            <input
              type="number"
              min={0}
              value={velMoney}
              onChange={(e) => setVelMoney(Number(e.target.value))}
            />
            <span className="v3-ch31-hint">{penceToPounds(velMoney)}</span>
          </label>
          <span className="v3-ch31-calc-op">×</span>
          <label>
            Times lent (&ge; 1)
            <input
              type="number"
              min={1}
              value={velTimes}
              onChange={(e) => setVelTimes(Number(e.target.value))}
            />
          </label>
          <span className="v3-ch31-calc-op">=</span>
          <div className="v3-ch31-calc-result">
            <span className="v3-ch31-calc-label">Loan capital created</span>
            <span className="v3-ch31-calc-val">{penceToPounds(loanCapitalCreated)}</span>
          </div>
        </div>
      </section>

      {/* Exemplar cards */}
      <section className="v3-ch31-section" aria-label="Marx's exemplars">
        <h4>Marx's exemplars: the two sources</h4>
        <div className="v3-ch31-exemplars">
          <div className="v3-ch31-exemplar">
            <div className="v3-ch31-exemplar-name">Ipswich farmers' surplus deposits</div>
            <div className="v3-ch31-exemplar-row">
              <span>Source</span>
              <span>Conversion of capital/revenue (2)</span>
            </div>
            <div className="v3-ch31-exemplar-row">
              <span>Relation</span>
              <span>Accompanies real accumulation</span>
            </div>
            <div className="v3-ch31-exemplar-note">
              Rural surplus rediscounted by Lombard Street bill-brokers to the
              manufacturing centres.
            </div>
          </div>
          <div className="v3-ch31-exemplar">
            <div className="v3-ch31-exemplar-name">Post-crisis idle industrial capital</div>
            <div className="v3-ch31-exemplar-row">
              <span>Source</span>
              <span>Simple transformation of money (1)</span>
            </div>
            <div className="v3-ch31-exemplar-row">
              <span>Relation</span>
              <span>Inverse to real accumulation</span>
            </div>
            <div className="v3-ch31-exemplar-note">
              Money lying fallow after a crisis, simply transformed into loan capital —
              it swells as production stagnates.
            </div>
          </div>
        </div>
      </section>

      {/* Create form */}
      <section className="v3-ch31-section" aria-label="Record a parcel of floating capital">
        <h4>Record a parcel of floating capital</h4>
        <p className="v3-ch31-subintro">
          Record the amount, its source, the velocity that multiplies it into loan
          capital, and any economy in the circulation reserve. The relation-to-accumulation
          verdict below mirrors the server's classification.
        </p>
        <div className="v3-ch31-form">
          <label>
            Name
            <input value={name} onChange={(e) => setName(e.target.value)} />
          </label>
          <label>
            Amount (pence, &ge; 0)
            <input
              type="number"
              min={0}
              value={amount}
              onChange={(e) => setAmount(Number(e.target.value))}
            />
            <span className="v3-ch31-hint">{penceToPounds(amount)}</span>
          </label>
          <label>
            Source
            <select value={source} onChange={(e) => setSource(Number(e.target.value))}>
              {SOURCES.map((s) => (
                <option key={s.value} value={s.value}>
                  {s.label}
                </option>
              ))}
            </select>
          </label>
          <label>
            Circulation reserve saving (pence, &ge; 0)
            <input
              type="number"
              min={0}
              value={reserveSaving}
              onChange={(e) => setReserveSaving(Number(e.target.value))}
            />
            <span className="v3-ch31-hint">{penceToPounds(reserveSaving)}</span>
          </label>
          <label>
            Reserve-saving description
            <input value={reserveSavingDesc} onChange={(e) => setReserveSavingDesc(e.target.value)} />
          </label>
          <label>
            Description
            <input value={description} onChange={(e) => setDescription(e.target.value)} />
          </label>
        </div>
        <div className="v3-ch31-preview">
          Velocity: {penceToPounds(velMoney)} × {velTimes} = {penceToPounds(loanCapitalCreated)} of loan capital
        </div>
        <div className="v3-ch31-preview">Relation: {previewRelation}</div>
        {formError && <p className="v3-ch31-error">{formError}</p>}
        <button className="v3-ch31-btn" onClick={handleCreate}>
          Save parcel
        </button>
        {createError && <p className="v3-ch31-error">{createError}</p>}
      </section>

      {/* Stored records */}
      <section className="v3-ch31-section" aria-label="Stored floating-capital parcels">
        <h4>Stored parcels</h4>
        {listError && <p className="v3-ch31-error">{listError}</p>}
        <div className="v3-ch31-cards">
          {records.length === 0 && (
            <p className="v3-ch31-subintro">No parcels stored yet.</p>
          )}
          {records.map((f) => (
            <div className="v3-ch31-card" key={f.id}>
              <div className="v3-ch31-card-name">{f.name}</div>
              <div className="v3-ch31-card-rows">
                <div className="v3-ch31-card-row">
                  <span>Amount</span>
                  <span>{penceToPounds(f.amount)}</span>
                </div>
                <div className="v3-ch31-card-row">
                  <span>Source</span>
                  <span>{sourceLabel(f.source)}</span>
                </div>
                <div className="v3-ch31-card-row">
                  <span>Velocity</span>
                  <span>
                    {penceToPounds(f.loan_capital_velocity.money_amount)} ×{" "}
                    {f.loan_capital_velocity.times_lent} ={" "}
                    {penceToPounds(f.loan_capital_velocity.loan_capital_created)}
                  </span>
                </div>
                <div className="v3-ch31-card-row">
                  <span>Reserve saving</span>
                  <span>{penceToPounds(f.circulation_reserve_saving.amount)}</span>
                </div>
              </div>
              <div className="v3-ch31-verdict">{relationToAccumulation(f.source)}</div>
              {f.description && <div className="v3-ch31-desc">{f.description}</div>}
            </div>
          ))}
        </div>
      </section>
    </div>
  );
}
