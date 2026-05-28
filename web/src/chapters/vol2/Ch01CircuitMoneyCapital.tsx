import { useEffect, useState } from "react";
import type { FormEvent } from "react";
import { moneyCircuitApi } from "../../api";
import { usePounds } from "../../CurrencyContext";
import type { MoneyCircuit, CreateMoneyCircuitInput } from "../../types";
import "./Ch01CircuitMoneyCapital.css";

const MOMENTS = ["M", "M-C", "P", "C-prime", "C-M-prime", "M-prime"] as const;

const MOMENT_LABELS: Record<string, string> = {
  M: "M",
  "M-C": "M—C",
  P: "P",
  "C-prime": "C′",
  "C-M-prime": "C′—M′",
  "M-prime": "M′",
};

const DEFAULT_CREATE: CreateMoneyCircuitInput = {
  amount: 42200,
};

export function Ch01CircuitMoneyCapital() {
  const fmt = usePounds();

  const [circuits, setCircuits] = useState<MoneyCircuit[]>([]);
  const [selectedId, setSelectedId] = useState("");
  const [listError, setListError] = useState("");

  // create form
  const [createForm, setCreateForm] = useState<CreateMoneyCircuitInput>(DEFAULT_CREATE);
  const [creating, setCreating] = useState(false);
  const [createError, setCreateError] = useState("");

  // purchase form
  const [labourAmount, setLabourAmount] = useState(26540);
  const [labourHours, setLabourHours] = useState(1440);
  const [meansAmount, setMeansAmount] = useState(15660);
  const [meansHours, setMeansHours] = useState(1440);
  const [purchaseError, setPurchaseError] = useState("");

  // produce form
  const [constantPence, setConstantPence] = useState(15660);
  const [variablePence, setVariablePence] = useState(26540);
  const [produceError, setProduceError] = useState("");

  // commodity form
  const [valueOriginal, setValueOriginal] = useState(42200);
  const [valueSurplus, setValueSurplus] = useState(27200);
  const [commodityError, setCommodityError] = useState("");

  // realise form
  const [realisedPence, setRealisedPence] = useState(69400);
  const [realiseError, setRealiseError] = useState("");

  useEffect(() => {
    refresh();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  async function refresh() {
    setListError("");
    try {
      const list = await moneyCircuitApi.list();
      setCircuits(list);
      if (list.length > 0 && !selectedId) setSelectedId(list[0].id);
    } catch (e) {
      setListError(String(e));
    }
  }

  const selected = circuits.find((c) => c.id === selectedId) ?? null;

  async function handleCreate(e: FormEvent) {
    e.preventDefault();
    setCreateError("");
    setCreating(true);
    try {
      const mc = await moneyCircuitApi.create(createForm);
      setCreateForm(DEFAULT_CREATE);
      await refresh();
      setSelectedId(mc.id);
    } catch (e) {
      setCreateError(String(e));
    } finally {
      setCreating(false);
    }
  }

  async function handlePurchase(e: FormEvent) {
    e.preventDefault();
    if (!selectedId) return;
    setPurchaseError("");
    try {
      await moneyCircuitApi.recordPurchase(selectedId, {
        labour_amount: labourAmount,
        labour_power_hours: labourHours,
        means_amount: meansAmount,
        means_capacity_hours: meansHours,
      });
      await refresh();
    } catch (e) {
      setPurchaseError(String(e));
    }
  }

  async function handleProduce(e: FormEvent) {
    e.preventDefault();
    if (!selectedId) return;
    setProduceError("");
    try {
      await moneyCircuitApi.recordProduce(selectedId, {
        constant_pence: constantPence,
        variable_pence: variablePence,
      });
      await refresh();
    } catch (e) {
      setProduceError(String(e));
    }
  }

  async function handleCommodity(e: FormEvent) {
    e.preventDefault();
    if (!selectedId) return;
    setCommodityError("");
    try {
      await moneyCircuitApi.recordCommodity(selectedId, {
        value_original: valueOriginal,
        value_surplus: valueSurplus,
      });
      await refresh();
    } catch (e) {
      setCommodityError(String(e));
    }
  }

  async function handleRealise(e: FormEvent) {
    e.preventDefault();
    if (!selectedId) return;
    setRealiseError("");
    try {
      await moneyCircuitApi.recordRealise(selectedId, {
        realised_pence: realisedPence,
      });
      await refresh();
    } catch (e) {
      setRealiseError(String(e));
    }
  }

  const currentMoment = selected?.moment ?? null;

  return (
    <div>
      {/* Circuit diagram */}
      <div className="mc-card" style={{ marginBottom: "1.5rem" }}>
        <h3>The Circuit of Money-Capital — M—C(Lp+Mp)…P…C′—M′</h3>
        <div className="mc-diagram">
          {(["M", "M-C", "P", "C-prime", "C-M-prime", "M-prime"] as const).map((m, i, arr) => (
            <span key={m} className="mc-diagram-group">
              <span className={`node ${currentMoment === m ? "active" : ""}`}>
                {MOMENT_LABELS[m]}
              </span>
              {i < arr.length - 1 && (
                <span className="arrow">{m === "M-C" || m === "C-prime" ? "…" : "—"}</span>
              )}
            </span>
          ))}
        </div>
        <p style={{ fontSize: "0.8rem", color: "var(--ink-muted)", margin: 0 }}>
          Money-capital advanced (M) buys labour-power and means of production (M—C). Production
          consumes them (…P…). Output is commodity-capital bearing surplus-value (C′). Sale
          realises the expanded money-capital (C′—M′). The highlighted node shows the current
          moment of the selected circuit.
        </p>
      </div>

      <div className="mc-grid">
        {/* Left: list + create */}
        <div>
          <h3 style={{ margin: "0 0 0.5rem", fontSize: "0.9rem" }}>Money Circuits</h3>
          {listError && <p className="mc-error">{listError}</p>}
          <ul className="mc-list">
            {circuits.map((c) => (
              <li
                key={c.id}
                className={c.id === selectedId ? "selected" : ""}
                role="button"
                tabIndex={0}
                aria-pressed={c.id === selectedId}
                onClick={() => setSelectedId(c.id)}
                onKeyDown={(e) => { if (e.key === 'Enter' || e.key === ' ') { e.preventDefault(); setSelectedId(c.id); } }}
              >
                <span className="mc-list-moment">{MOMENT_LABELS[c.moment]}</span>
                <span>{fmt(c.advance.amount)}</span>
              </li>
            ))}
            {circuits.length === 0 && (
              <li style={{ color: "var(--ink-muted)", cursor: "default" }}>
                No circuits yet
              </li>
            )}
          </ul>

          <form className="mc-form" onSubmit={handleCreate}>
            <label>
              Advance (pence)
              <input
                type="number"
                value={createForm.amount}
                onChange={(e) => setCreateForm({ ...createForm, amount: Number(e.target.value) })}
                min={1}
              />
            </label>
            <label>
              Agent ID (optional)
              <input
                type="text"
                value={createForm.agent_id ?? ""}
                onChange={(e) =>
                  setCreateForm({ ...createForm, agent_id: e.target.value || undefined })
                }
                placeholder="agent id"
              />
            </label>
            <button type="submit" disabled={creating}>
              {creating ? "Advancing…" : "Advance M"}
            </button>
            {createError && <p className="mc-error">{createError}</p>}
          </form>
        </div>

        {/* Right: selected circuit detail + phase forms */}
        {selected && (
          <div>
            {/* Moment progress */}
            <div className="mc-card" style={{ marginBottom: "1rem" }}>
              <h3>Current moment</h3>
              <div className="mc-moment-row">
                {MOMENTS.map((m) => (
                  <span
                    key={m}
                    className={`mc-moment-chip ${selected.moment === m ? "current" : ""} ${
                      momentIndex(selected.moment) > momentIndex(m) ? "past" : ""
                    }`}
                  >
                    {MOMENT_LABELS[m]}
                  </span>
                ))}
              </div>
            </div>

            {/* Balance summary */}
            <div className="mc-card" style={{ marginBottom: "1rem" }}>
              <h3>Balance sheet</h3>
              <dl className="mc-balance">
                <dt>Money advanced (M)</dt>
                <dd>{fmt(selected.advance.amount)}</dd>
                {selected.purchase && (
                  <>
                    <dt>Purchase total (C)</dt>
                    <dd>{fmt(selected.purchase.total)}</dd>
                    <dt className="sub">— Labour-power (Lp)</dt>
                    <dd className="sub">{fmt(selected.purchase.labour_leg.amount)}</dd>
                    <dt className="sub">— Means of production (Mp)</dt>
                    <dd className="sub">{fmt(selected.purchase.means_leg.amount)}</dd>
                  </>
                )}
                {selected.productive && (
                  <>
                    <dt>Productive capital (c+v)</dt>
                    <dd>{fmt(selected.productive.total_advance)}</dd>
                  </>
                )}
                {selected.commodity && (
                  <>
                    <dt>Commodity-capital (C′ = c+v+s)</dt>
                    <dd>{fmt(selected.commodity.total)}</dd>
                    <dt className="sub">— Surplus-value (s)</dt>
                    <dd className="sub">{fmt(selected.commodity.value_surplus)}</dd>
                  </>
                )}
                {selected.realisation && (
                  <>
                    <dt>Realised (M′)</dt>
                    <dd>
                      <strong>{fmt(selected.realisation.realised_pence)}</strong>
                    </dd>
                    <dt>ΔM (surplus realised)</dt>
                    <dd>
                      <strong>{fmt(selected.realisation.surplus_realised)}</strong>
                    </dd>
                  </>
                )}
              </dl>
            </div>

            {/* Phase forms — only show the next available action */}
            {selected.moment === "M" && (
              <div className="mc-card">
                <h3>Record purchase phase — M—C(Lp+Mp)</h3>
                <form className="mc-form" onSubmit={handlePurchase}>
                  <div className="mc-form-row">
                    <label>
                      Labour amount (p)
                      <input
                        type="number"
                        value={labourAmount}
                        onChange={(e) => setLabourAmount(Number(e.target.value))}
                        min={1}
                      />
                    </label>
                    <label>
                      Labour-power hours
                      <input
                        type="number"
                        value={labourHours}
                        onChange={(e) => setLabourHours(Number(e.target.value))}
                        min={1}
                      />
                    </label>
                  </div>
                  <div className="mc-form-row">
                    <label>
                      Means amount (p)
                      <input
                        type="number"
                        value={meansAmount}
                        onChange={(e) => setMeansAmount(Number(e.target.value))}
                        min={1}
                      />
                    </label>
                    <label>
                      Means capacity hours
                      <input
                        type="number"
                        value={meansHours}
                        onChange={(e) => setMeansHours(Number(e.target.value))}
                        min={1}
                      />
                    </label>
                  </div>
                  <button type="submit">Record M—C</button>
                  {purchaseError && <p className="mc-error">{purchaseError}</p>}
                </form>
              </div>
            )}

            {selected.moment === "M-C" && (
              <div className="mc-card">
                <h3>Enter production — …P…</h3>
                <form className="mc-form" onSubmit={handleProduce}>
                  <div className="mc-form-row">
                    <label>
                      Constant capital (p)
                      <input
                        type="number"
                        value={constantPence}
                        onChange={(e) => setConstantPence(Number(e.target.value))}
                        min={0}
                      />
                    </label>
                    <label>
                      Variable capital (p)
                      <input
                        type="number"
                        value={variablePence}
                        onChange={(e) => setVariablePence(Number(e.target.value))}
                        min={0}
                      />
                    </label>
                  </div>
                  <button type="submit">Enter production</button>
                  {produceError && <p className="mc-error">{produceError}</p>}
                </form>
              </div>
            )}

            {selected.moment === "P" && (
              <div className="mc-card">
                <h3>Exit production — C′ commodity-capital</h3>
                <form className="mc-form" onSubmit={handleCommodity}>
                  <div className="mc-form-row">
                    <label>
                      Original value (p)
                      <input
                        type="number"
                        value={valueOriginal}
                        onChange={(e) => setValueOriginal(Number(e.target.value))}
                        min={1}
                      />
                    </label>
                    <label>
                      Surplus-value (p)
                      <input
                        type="number"
                        value={valueSurplus}
                        onChange={(e) => setValueSurplus(Number(e.target.value))}
                        min={0}
                      />
                    </label>
                  </div>
                  <button type="submit">Attach C′</button>
                  {commodityError && <p className="mc-error">{commodityError}</p>}
                </form>
              </div>
            )}

            {selected.moment === "C-prime" && (
              <div className="mc-card">
                <h3>Realise — C′—M′</h3>
                <form className="mc-form" onSubmit={handleRealise}>
                  <label>
                    Realised pence (M′)
                    <input
                      type="number"
                      value={realisedPence}
                      onChange={(e) => setRealisedPence(Number(e.target.value))}
                      min={1}
                    />
                  </label>
                  <button type="submit">Realise M′</button>
                  {realiseError && <p className="mc-error">{realiseError}</p>}
                </form>
              </div>
            )}

            {(selected.moment === "C-M-prime" || selected.moment === "M-prime") && (
              <div className="mc-card">
                <h3>Circuit complete — M′</h3>
                <p style={{ fontSize: "0.875rem", color: "var(--ink-muted)" }}>
                  The circuit M—C(Lp+Mp)…P…C′—M′ has completed its full rotation. The expanded
                  money-capital M′ is available for a new circuit.
                </p>
                {selected.realisation && (
                  <p style={{ fontWeight: 600, fontSize: "1rem" }}>
                    M′ = {fmt(selected.realisation.realised_pence)} &nbsp;|&nbsp; ΔM ={" "}
                    {fmt(selected.realisation.surplus_realised)}
                  </p>
                )}
              </div>
            )}
          </div>
        )}
      </div>
    </div>
  );
}

function momentIndex(m: string): number {
  return ["M", "M-C", "P", "C-prime", "C-M-prime", "M-prime"].indexOf(m);
}
