import { useEffect, useState } from "react";
import type { FormEvent } from "react";
import { api } from "../../api";
import { usePounds } from "../../CurrencyContext";
import type {
  ProductiveCircuit,
  CreateProductiveCircuitInput,
} from "../../types";
import "./Ch02CircuitProductiveCapital.css";

type ReproductionView = "simple" | "extended";

const DEFAULT_FORM: CreateProductiveCircuitInput = {
  constant_pence: 37200,
  variable_pence: 5000,
  mode: "simple",
  min_capitalisation_increment_pence: 100,
};

function FormIIAsymmetryInsight() {
  return (
    <section className="v2-ch02-insight">
      <h2 className="v2-ch02-insight-h2">
        Form II — production is both start and end
      </h2>
      <p className="v2-ch02-insight-prose">
        Form I (M…M′) makes profit-making appear to be the goal of the whole
        process — money in, more money out, with production a means. Form II
        (P…P) re-reads the same circuit starting from production. Now
        money-capital and commodity-capital are only transitional moments;
        the continuation of production itself becomes the goal, and the
        capitalist's surplus appears as the means by which that continuation
        is enlarged. Neither form is wrong — together they show that the
        circuit is the unity of all three.
      </p>
    </section>
  );
}

export function Ch02CircuitProductiveCapital() {
  const fmt = usePounds();

  const [circuits, setCircuits] = useState<ProductiveCircuit[]>([]);
  const [selectedId, setSelectedId] = useState("");
  const [listError, setListError] = useState("");

  const [form, setForm] = useState<CreateProductiveCircuitInput>(DEFAULT_FORM);
  const [creating, setCreating] = useState(false);
  const [createError, setCreateError] = useState("");

  const [view, setView] = useState<ReproductionView>("simple");

  const [depositAmount, setDepositAmount] = useState(7800);
  const [depositError, setDepositError] = useState("");

  const [drawAmount, setDrawAmount] = useState(5000);
  const [drawReason, setDrawReason] = useState<"delayed_realisation" | "price_rise_in_mp">("delayed_realisation");
  const [drawError, setDrawError] = useState("");

  useEffect(() => {
    refresh();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  async function refresh() {
    setListError("");
    try {
      const list = await api.listProductiveCircuits();
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
      const created = await api.createProductiveCircuit(form);
      setForm(DEFAULT_FORM);
      await refresh();
      setSelectedId(created.id);
    } catch (e) {
      setCreateError(String(e));
    } finally {
      setCreating(false);
    }
  }

  async function handleDeposit(e: FormEvent) {
    e.preventDefault();
    if (!selectedId) return;
    setDepositError("");
    try {
      await api.depositReserve(selectedId, depositAmount);
      await refresh();
    } catch (e) {
      setDepositError(String(e));
    }
  }

  async function handleDraw(e: FormEvent) {
    e.preventDefault();
    if (!selectedId) return;
    setDrawError("");
    try {
      await api.drawReserve(selectedId, drawAmount, drawReason);
      await refresh();
    } catch (e) {
      setDrawError(String(e));
    }
  }

  const lmc = selected?.latent_money_capital;
  const pct = lmc && lmc.threshold > 0
    ? Math.min(100, Math.round((lmc.accumulated / lmc.threshold) * 100))
    : 0;

  return (
    <div>
      <FormIIAsymmetryInsight />
      {/* Circuit diagram — P…C'—M'—C…P */}
      <div className="pc-card" style={{ marginBottom: "1.5rem" }}>
        <h3>The Circuit of Productive Capital — P … C&prime;—M&prime;—C … P</h3>
        <div className="pc-diagram">
          <span className="node active">P</span>
          <span className="arrow">…</span>
          <span className="node">C&prime;</span>
          <span className="arrow">—</span>
          <span className="node">M&prime;</span>
          <span className="arrow">—</span>
          <span className="node">C(Lp+Mp)</span>
          <span className="arrow">…</span>
          <span className="node active">P</span>
        </div>
        <p style={{ fontSize: "0.8rem", color: "var(--ink-muted)", margin: 0 }}>
          Starts and ends in production. Circulation is the connecting link between two functionings of
          productive capital. Compare Ch. 1&apos;s M-form: both describe the same circuit from a different
          starting-point.
        </p>
      </div>

      <div className="pc-grid">
        {/* Left column: circuit list + create form */}
        <div>
          <h3 style={{ margin: "0 0 0.5rem", fontSize: "0.9rem" }}>Productive Circuits</h3>
          {listError && <p className="pc-error">{listError}</p>}
          <ul className="pc-list">
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
                {fmt(c.constant_pence + c.variable_pence)} — {c.mode}
              </li>
            ))}
            {circuits.length === 0 && (
              <li style={{ color: "var(--ink-muted)", cursor: "default" }}>No circuits yet</li>
            )}
          </ul>

          <form className="pc-form" onSubmit={handleCreate}>
            <div className="pc-form-row">
              <label>
                C (constant, pence)
                <input
                  type="number"
                  value={form.constant_pence}
                  onChange={(e) => setForm({ ...form, constant_pence: Number(e.target.value) })}
                  min={1}
                />
              </label>
              <label>
                V (variable, pence)
                <input
                  type="number"
                  value={form.variable_pence}
                  onChange={(e) => setForm({ ...form, variable_pence: Number(e.target.value) })}
                  min={1}
                />
              </label>
            </div>
            <div className="pc-form-row">
              <label>
                Mode
                <select
                  value={form.mode}
                  onChange={(e) =>
                    setForm({ ...form, mode: e.target.value as "simple" | "extended" | "mixed" })
                  }
                >
                  <option value="simple">Simple reproduction</option>
                  <option value="extended">Extended reproduction</option>
                  <option value="mixed">Mixed</option>
                </select>
              </label>
              <label>
                Min capitalisation (p)
                <input
                  type="number"
                  value={form.min_capitalisation_increment_pence}
                  onChange={(e) =>
                    setForm({ ...form, min_capitalisation_increment_pence: Number(e.target.value) })
                  }
                  min={1}
                />
              </label>
            </div>
            <button type="submit" disabled={creating}>
              {creating ? "Creating…" : "Open circuit"}
            </button>
            {createError && <p className="pc-error">{createError}</p>}
          </form>
        </div>

        {/* Right column: selected circuit detail */}
        {selected && (
          <div>
            <div className="pc-toggle">
              <button
                className={view === "simple" ? "active" : ""}
                onClick={() => setView("simple")}
              >
                Simple P…P
              </button>
              <button
                className={view === "extended" ? "active" : ""}
                onClick={() => setView("extended")}
              >
                Extended P…P&prime;
              </button>
            </div>

            <div className="pc-card" style={{ marginBottom: "1rem" }}>
              <h3>Balance sheet — circuit n</h3>
              <dl className="pc-balance">
                <dt>Constant capital (C)</dt>
                <dd>{fmt(selected.constant_pence)}</dd>
                <dt>Variable capital (V)</dt>
                <dd>{fmt(selected.variable_pence)}</dd>
                <dt>Total advance (C+V)</dt>
                <dd>
                  <strong>{fmt(selected.total_advance)}</strong>
                </dd>
              </dl>
            </div>

            {view === "extended" && selected.capitalisation_steps.length > 0 && (
              <div className="pc-card" style={{ marginBottom: "1rem" }}>
                <h3>Balance sheet — circuit n+1 (after capitalisation)</h3>
                {(() => {
                  const totalInjected = selected.capitalisation_steps.reduce(
                    (acc, s) => acc + s.amount_injected,
                    0,
                  );
                  return (
                    <dl className="pc-balance">
                      <dt>Opening advance</dt>
                      <dd>{fmt(selected.total_advance)}</dd>
                      <dt>+ capitalised m</dt>
                      <dd>+ {fmt(totalInjected)}</dd>
                      <dt>P&prime; (n+1)</dt>
                      <dd>
                        <strong>{fmt(selected.total_advance + totalInjected)}</strong>
                      </dd>
                    </dl>
                  );
                })()}
              </div>
            )}

            {/* Latent capital gauge */}
            <div className="pc-card" style={{ marginBottom: "1rem" }}>
              <h3>Latent money-capital (hoarded m)</h3>
              <div className="latent-gauge">
                <div className="latent-gauge-track">
                  <div
                    className={`latent-gauge-fill ${lmc?.is_armed ? "armed" : "not-armed"}`}
                    style={{ width: `${pct}%` }}
                  />
                </div>
                <p className="latent-gauge-threshold">
                  {fmt(lmc?.accumulated ?? 0)} / {fmt(lmc?.threshold ?? 0)} threshold
                  {lmc?.is_armed
                    ? " — armed, ready to capitalise"
                    : " — hoarding, not yet functional as capital"}
                </p>
              </div>
            </div>

            {/* Reserve fund */}
            <div className="pc-card">
              <h3>Reserve fund — {fmt(selected.reserve_fund.balance)}</h3>
              <form className="pc-form" onSubmit={handleDeposit}>
                <div className="pc-form-row">
                  <label>
                    Deposit (pence)
                    <input
                      type="number"
                      value={depositAmount}
                      onChange={(e) => setDepositAmount(Number(e.target.value))}
                      min={1}
                    />
                  </label>
                  <button type="submit">Deposit</button>
                </div>
                {depositError && <p className="pc-error">{depositError}</p>}
              </form>
              <form className="pc-form" onSubmit={handleDraw}>
                <div className="pc-form-row">
                  <label>
                    Draw (pence)
                    <input
                      type="number"
                      value={drawAmount}
                      onChange={(e) => setDrawAmount(Number(e.target.value))}
                      min={1}
                    />
                  </label>
                  <label>
                    Reason
                    <select
                      value={drawReason}
                      onChange={(e) => setDrawReason(e.target.value as typeof drawReason)}
                    >
                      <option value="delayed_realisation">Delayed realisation</option>
                      <option value="price_rise_in_mp">Price rise in MP</option>
                    </select>
                  </label>
                  <button type="submit">Draw</button>
                </div>
                {drawError && <p className="pc-error">{drawError}</p>}
              </form>

              {selected.reserve_draws.length > 0 && (
                <div className="reserve-history">
                  <h4>Draw history</h4>
                  <ul>
                    {selected.reserve_draws.map((d, i) => (
                      <li key={i}>
                        <span>{d.reason.replace(/_/g, " ")}</span>
                        <span>{fmt(d.drawn_pence)}</span>
                      </li>
                    ))}
                  </ul>
                </div>
              )}
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
