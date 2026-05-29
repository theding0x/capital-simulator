import { useEffect, useState } from "react";
import type { FormEvent } from "react";
import { commodityCircuitApi } from "../../api";
import { usePounds } from "../../CurrencyContext";
import type {
  CommodityCircuit,
  CreateCommodityCircuitInput,
  CommodityCircuitAggregate,
} from "../../types";
import "./Ch03CircuitCommodityCapital.css";

const DEFAULT_FORM: CreateCommodityCircuitInput = {
  constant_pence: 37200,
  variable_pence: 5000,
  surplus_pence: 7800,
  pounds_total: 10000,
  mode: "simple",
};

function SocialCapitalCoda() {
  return (
    <aside className="v2-ch03-coda">
      <p className="v2-ch03-coda-quote">
        “Form III, C′…C′, alone presupposes the cycle as the cycle of social
        capital, not the cycle of an individual capital in isolation.
        Each capitalist's commodity-product is at the same time another
        capitalist's means of production — the means of production enter the
        opening C′ already as the output of other circuits, and the closing
        C′ enters someone else's circuit as input. The form makes the
        reproduction of the social total product visible.”
        <span className="v2-ch03-coda-cite">
          — Marx, Capital Vol. II, Ch. 3 (paraphrasing §IV)
        </span>
      </p>
    </aside>
  );
}

export function Ch03CircuitCommodityCapital() {
  const fmt = usePounds();

  const [circuits, setCircuits] = useState<CommodityCircuit[]>([]);
  const [selectedId, setSelectedId] = useState("");
  const [listError, setListError] = useState("");
  const [aggregate, setAggregate] = useState<CommodityCircuitAggregate | null>(null);

  const [form, setForm] = useState<CreateCommodityCircuitInput>(DEFAULT_FORM);
  const [creating, setCreating] = useState(false);
  const [createError, setCreateError] = useState("");

  const [saleQty, setSaleQty] = useState(7440);
  const [saleRealised, setSaleRealised] = useState(37200);
  const [saleError, setSaleError] = useState("");

  const [mpKind, setMpKind] = useState<"producer_circuit" | "merchant_holding" | "import">("import");
  const [mpError, setMpError] = useState("");

  useEffect(() => {
    refresh();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  async function refresh() {
    setListError("");
    try {
      const [list, agg] = await Promise.all([
        commodityCircuitApi.list(),
        commodityCircuitApi.aggregate(),
      ]);
      setCircuits(list);
      setAggregate(agg);
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
      const created = await commodityCircuitApi.create(form);
      setForm(DEFAULT_FORM);
      await refresh();
      setSelectedId(created.id);
    } catch (e) {
      setCreateError(String(e));
    } finally {
      setCreating(false);
    }
  }

  async function handleSale(e: FormEvent) {
    e.preventDefault();
    if (!selectedId) return;
    setSaleError("");
    try {
      await commodityCircuitApi.recordPartialSale(selectedId, saleQty, saleRealised);
      await refresh();
    } catch (e) {
      setSaleError(String(e));
    }
  }

  async function handleLinkMP(e: FormEvent) {
    e.preventDefault();
    if (!selectedId) return;
    setMpError("");
    try {
      await commodityCircuitApi.linkMPSource(selectedId, mpKind);
      await refresh();
    } catch (e) {
      setMpError(String(e));
    }
  }

  return (
    <div>
      {/* Circuit diagram */}
      <div className="cc-card" style={{ marginBottom: "1.5rem" }}>
        <h3>The Circuit of Commodity-Capital — C&prime;—M&prime;—C(Lp+Mp) … P … C&prime;</h3>
        <div className="cc-diagram">
          <span className="node active">C&prime;</span>
          <span className="arrow">—</span>
          <span className="node">M&prime;</span>
          <span className="arrow">—</span>
          <span className="node">C(Lp+Mp)</span>
          <span className="arrow">…</span>
          <span className="node">P</span>
          <span className="arrow">…</span>
          <span className="node active">C&prime;</span>
        </div>
        <p style={{ fontSize: "0.8rem", color: "var(--ink-muted)", margin: 0 }}>
          The third form of the capital circuit. Starts and ends in commodity-form.
          Unlike Forms I and II, surplus-value is already present in the opening C&prime; (c+v+s).
          Exposes the reproduction of social capital and the external origin of Means of Production.
        </p>
      </div>

      {/* Aggregate (social capital lens) */}
      {aggregate && aggregate.circuit_count > 0 && (
        <div className="cc-card cc-aggregate" style={{ marginBottom: "1.5rem" }}>
          <h3>Social Capital Aggregate</h3>
          <dl className="cc-balance">
            <dt>Total circuits</dt>
            <dd>{aggregate.circuit_count}</dd>
            <dt>Closed circuits</dt>
            <dd>{aggregate.closed_count}</dd>
            <dt>Total opening C&prime;</dt>
            <dd>{fmt(aggregate.total_initial_pence)}</dd>
            <dt>Total terminal C&prime;</dt>
            <dd>{fmt(aggregate.total_terminal_pence)}</dd>
          </dl>
        </div>
      )}

      <div className="cc-grid">
        {/* Left: list + create form */}
        <div>
          <h3 style={{ margin: "0 0 0.5rem", fontSize: "0.9rem" }}>Commodity Circuits</h3>
          {listError && <p className="cc-error">{listError}</p>}
          <ul className="cc-list">
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
                {fmt(c.initial.total)} — {c.mode}
                {c.terminal && <span className="cc-closed-badge">closed</span>}
              </li>
            ))}
            {circuits.length === 0 && (
              <li style={{ color: "var(--ink-muted)", cursor: "default" }}>No circuits yet</li>
            )}
          </ul>

          <form className="cc-form" onSubmit={handleCreate}>
            <div className="cc-form-row">
              <label>
                C (constant, pence)
                <input
                  type="number"
                  value={form.constant_pence}
                  onChange={(e) => setForm({ ...form, constant_pence: Number(e.target.value) })}
                  min={0}
                />
              </label>
              <label>
                V (variable, pence)
                <input
                  type="number"
                  value={form.variable_pence}
                  onChange={(e) => setForm({ ...form, variable_pence: Number(e.target.value) })}
                  min={0}
                />
              </label>
            </div>
            <div className="cc-form-row">
              <label>
                S (surplus, pence)
                <input
                  type="number"
                  value={form.surplus_pence}
                  onChange={(e) => setForm({ ...form, surplus_pence: Number(e.target.value) })}
                  min={0}
                />
              </label>
              <label>
                Pounds total (lbs)
                <input
                  type="number"
                  value={form.pounds_total}
                  onChange={(e) => setForm({ ...form, pounds_total: Number(e.target.value) })}
                  min={1}
                />
              </label>
            </div>
            <div className="cc-form-row">
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
            </div>
            <button type="submit" disabled={creating}>
              {creating ? "Opening…" : "Open C′ circuit"}
            </button>
            {createError && <p className="cc-error">{createError}</p>}
          </form>
        </div>

        {/* Right: selected circuit detail */}
        {selected && (
          <div>
            <div className="cc-card" style={{ marginBottom: "1rem" }}>
              <h3>Opening C&prime; = c + v + s</h3>
              <dl className="cc-balance">
                <dt>Constant (c)</dt>
                <dd>{fmt(selected.initial.constant_pence)}</dd>
                <dt>Variable (v)</dt>
                <dd>{fmt(selected.initial.variable_pence)}</dd>
                <dt>Capital-value (c+v)</dt>
                <dd>{fmt(selected.initial.value_original_pence)}</dd>
                <dt>Surplus-value (s)</dt>
                <dd>{fmt(selected.initial.surplus_pence)}</dd>
                <dt>Total C&prime;</dt>
                <dd><strong>{fmt(selected.initial.total)}</strong></dd>
                <dt>Output (lbs)</dt>
                <dd>{selected.initial.pounds_total.toLocaleString()}</dd>
              </dl>
            </div>

            {/* Partial sales */}
            <div className="cc-card" style={{ marginBottom: "1rem" }}>
              <h3>Successive Partial Sales</h3>
              {selected.partial_sales.length > 0 ? (
                <table className="cc-sales-table">
                  <thead>
                    <tr>
                      <th>Qty (lbs)</th>
                      <th>Realised</th>
                      <th>c-share</th>
                      <th>v-share</th>
                      <th>s-share</th>
                    </tr>
                  </thead>
                  <tbody>
                    {selected.partial_sales.map((s, i) => (
                      <tr key={i}>
                        <td>{s.quantity.toLocaleString()}</td>
                        <td>{fmt(s.realised_pence)}</td>
                        <td>{fmt(s.decomposition.constant_share)}</td>
                        <td>{fmt(s.decomposition.variable_share)}</td>
                        <td>{fmt(s.decomposition.surplus_share)}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              ) : (
                <p style={{ fontSize: "0.8rem", color: "var(--ink-muted)", margin: 0 }}>
                  No sales recorded yet.
                </p>
              )}
              {!selected.terminal && (
                <form className="cc-form" style={{ marginTop: "0.75rem" }} onSubmit={handleSale}>
                  <div className="cc-form-row">
                    <label>
                      Qty (lbs)
                      <input
                        type="number"
                        value={saleQty}
                        onChange={(e) => setSaleQty(Number(e.target.value))}
                        min={1}
                      />
                    </label>
                    <label>
                      Realised (pence)
                      <input
                        type="number"
                        value={saleRealised}
                        onChange={(e) => setSaleRealised(Number(e.target.value))}
                        min={1}
                      />
                    </label>
                    <button type="submit">Record sale</button>
                  </div>
                  {saleError && <p className="cc-error">{saleError}</p>}
                </form>
              )}
            </div>

            {/* MP Sources */}
            <div className="cc-card" style={{ marginBottom: "1rem" }}>
              <h3>Means of Production Sources</h3>
              {selected.mp_sources.length > 0 ? (
                <ul className="cc-mp-list">
                  {selected.mp_sources.map((src, i) => (
                    <li key={i}>
                      <span className="cc-mp-kind">{src.source_kind.replace(/_/g, " ")}</span>
                      {src.source_commodity_circuit_id && (
                        <span className="cc-mp-id">{src.source_commodity_circuit_id.slice(0, 8)}&hellip;</span>
                      )}
                    </li>
                  ))}
                </ul>
              ) : (
                <p style={{ fontSize: "0.8rem", color: "var(--ink-muted)", margin: 0 }}>
                  No MP sources linked.
                </p>
              )}
              {!selected.terminal && (
                <form className="cc-form" style={{ marginTop: "0.75rem" }} onSubmit={handleLinkMP}>
                  <div className="cc-form-row">
                    <label>
                      Source kind
                      <select
                        value={mpKind}
                        onChange={(e) => setMpKind(e.target.value as typeof mpKind)}
                      >
                        <option value="import">Import</option>
                        <option value="producer_circuit">Producer circuit</option>
                        <option value="merchant_holding">Merchant holding</option>
                      </select>
                    </label>
                    <button type="submit">Link</button>
                  </div>
                  {mpError && <p className="cc-error">{mpError}</p>}
                </form>
              )}
            </div>

            {/* Terminal state */}
            {selected.terminal && (
              <div className={`cc-card cc-terminal ${selected.terminal.is_extended ? "extended" : "simple"}`}>
                <h3>
                  Terminal C&prime; &mdash;{" "}
                  {selected.terminal.is_extended ? "Extended reproduction" : "Simple reproduction"}
                </h3>
                <dl className="cc-balance">
                  <dt>Constant (c)</dt>
                  <dd>{fmt(selected.terminal.constant_pence)}</dd>
                  <dt>Variable (v)</dt>
                  <dd>{fmt(selected.terminal.variable_pence)}</dd>
                  <dt>Surplus (s)</dt>
                  <dd>{fmt(selected.terminal.surplus_pence)}</dd>
                  <dt>Capitalised (m)</dt>
                  <dd>{fmt(selected.terminal.capitalised_pence)}</dd>
                  <dt>Total</dt>
                  <dd><strong>{fmt(selected.terminal.total)}</strong></dd>
                </dl>
              </div>
            )}
          </div>
        )}
      </div>
      <SocialCapitalCoda />
    </div>
  );
}
