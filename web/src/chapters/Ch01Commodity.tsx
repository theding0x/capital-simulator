import { useState, useMemo } from "react";
import type { FormEvent } from "react";
import { api } from "../api";
import type { Commodity, ExchangeRatio, SocialRelations, UpdateCommodityInput, CreateCommodityInput } from "../types";
import { fmtMinutes, fmtQty } from "../format";

function emptyDraft(): CreateCommodityInput {
  return {
    name: "",
    use_value: { description: "", unit: "" },
    concrete_labour: { kind: "", description: "" },
    snlt_per_unit: 0,
  };
}

interface Ch01Props {
  commodities: Commodity[];
  onSharedChanged: () => void;
}

export function Ch01Commodity({ commodities, onSharedChanged }: Ch01Props) {
  return (
    <>
      <NewCommodityForm onCreated={onSharedChanged} />
      <section className="card">
        <h2>Commodities</h2>
        {commodities.length === 0 ? (
          <div className="empty-state">
            <p>No commodities registered yet.</p>
            <p className="small">Use the form above to add the first one.</p>
          </div>
        ) : (
          <CommodityTable commodities={commodities} onChanged={onSharedChanged} />
        )}
      </section>
      <ExchangeRatioPanel commodities={commodities} />
    </>
  );
}

function NewCommodityForm({ onCreated }: { onCreated: () => void }) {
  const [draft, setDraft] = useState<CreateCommodityInput>(emptyDraft);
  const [busy, setBusy] = useState(false);
  const [err, setErr] = useState<string | null>(null);

  async function submit(e: FormEvent) {
    e.preventDefault();
    setBusy(true);
    setErr(null);
    try {
      await api.createCommodity({
        ...draft,
        snlt_per_unit: Number(draft.snlt_per_unit) || 0,
      });
      setDraft(emptyDraft());
      onCreated();
    } catch (e) {
      setErr(e instanceof Error ? e.message : String(e));
    } finally {
      setBusy(false);
    }
  }

  return (
    <section className="card">
      <h2>Register a commodity</h2>
      <form className="form-grid" onSubmit={submit}>
        <label>
          <span>Name</span>
          <input
            value={draft.name}
            onChange={(e) => setDraft({ ...draft, name: e.target.value })}
            placeholder="linen"
            required
          />
        </label>
        <label>
          <span>Unit of measure</span>
          <input
            value={draft.use_value.unit}
            onChange={(e) =>
              setDraft({ ...draft, use_value: { ...draft.use_value, unit: e.target.value } })
            }
            placeholder="yards"
            required
          />
        </label>
        <label className="span2">
          <span>Use-value description</span>
          <input
            value={draft.use_value.description}
            onChange={(e) =>
              setDraft({ ...draft, use_value: { ...draft.use_value, description: e.target.value } })
            }
            placeholder="linen for clothing"
            required
          />
        </label>
        <label>
          <span>Concrete labour kind</span>
          <input
            value={draft.concrete_labour.kind}
            onChange={(e) =>
              setDraft({ ...draft, concrete_labour: { ...draft.concrete_labour, kind: e.target.value } })
            }
            placeholder="weaving"
            required
          />
        </label>
        <label>
          <span>SNLT (minutes per unit)</span>
          <input
            type="number"
            min={0}
            value={draft.snlt_per_unit}
            onChange={(e) => setDraft({ ...draft, snlt_per_unit: Number(e.target.value) })}
            placeholder="30"
            required
          />
        </label>
        <div className="span2 form-actions">
          <button type="submit" className="primary" disabled={busy}>
            {busy ? "Saving…" : "Register"}
          </button>
          {err && <span className="error">{err}</span>}
        </div>
      </form>
    </section>
  );
}

function CommodityTable({
  commodities,
  onChanged,
}: {
  commodities: Commodity[];
  onChanged: () => void;
}) {
  return (
    <table className="commodity-table">
      <thead>
        <tr>
          <th>Name</th>
          <th>Use-value</th>
          <th>Concrete labour</th>
          <th>SNLT / unit</th>
          <th></th>
        </tr>
      </thead>
      <tbody>
        {commodities.map((c) => (
          <CommodityRow key={c.id} commodity={c} onChanged={onChanged} />
        ))}
      </tbody>
    </table>
  );
}

function CommodityRow({
  commodity: c,
  onChanged,
}: {
  commodity: Commodity;
  onChanged: () => void;
}) {
  const [editing, setEditing] = useState(false);
  const [revealed, setRevealed] = useState(false);
  const [draft, setDraft] = useState<UpdateCommodityInput>({});
  const [busy, setBusy] = useState(false);
  const [err, setErr] = useState<string | null>(null);
  const [relations, setRelations] = useState<SocialRelations | null>(null);

  async function save() {
    setBusy(true);
    setErr(null);
    try {
      await api.updateCommodity(c.id, draft);
      setEditing(false);
      setDraft({});
      onChanged();
    } catch (e) {
      setErr(e instanceof Error ? e.message : String(e));
    } finally {
      setBusy(false);
    }
  }

  async function remove() {
    if (!confirm(`Delete "${c.name}"?`)) return;
    setBusy(true);
    try {
      await api.deleteCommodity(c.id);
      onChanged();
    } catch (e) {
      setErr(e instanceof Error ? e.message : String(e));
    } finally {
      setBusy(false);
    }
  }

  async function reveal() {
    if (revealed) {
      setRevealed(false);
      return;
    }
    try {
      const r = await api.socialRelations(c.id);
      setRelations(r);
      setRevealed(true);
    } catch (e) {
      setErr(e instanceof Error ? e.message : String(e));
    }
  }

  return (
    <>
      <tr>
        <td className="name-cell">{c.name}</td>
        <td>
          <div>{c.use_value.description}</div>
          <div className="muted small">per {c.use_value.unit}</div>
        </td>
        <td className="muted">{c.concrete_labour.kind}</td>
        <td className="snlt-cell">{fmtMinutes(c.snlt_per_unit)}</td>
        <td className="actions">
          <button className="secondary" onClick={() => setEditing((v) => !v)}>
            {editing ? "Cancel" : "Edit"}
          </button>
          <button className="reveal-btn" onClick={reveal}>
            {revealed ? "Hide" : "Reveal"}
          </button>
          <button className="danger" onClick={remove} disabled={busy}>
            Delete
          </button>
        </td>
      </tr>
      {editing && (
        <tr className="edit-row">
          <td colSpan={5}>
            <div className="form-grid">
              <label>
                <span>Name</span>
                <input
                  defaultValue={c.name}
                  onChange={(e) => setDraft((d) => ({ ...d, name: e.target.value }))}
                />
              </label>
              <label>
                <span>SNLT (minutes / unit)</span>
                <input
                  type="number"
                  min={0}
                  defaultValue={c.snlt_per_unit}
                  onChange={(e) =>
                    setDraft((d) => ({ ...d, snlt_per_unit: Number(e.target.value) }))
                  }
                />
              </label>
              <label>
                <span>Use-value description</span>
                <input
                  defaultValue={c.use_value.description}
                  onChange={(e) =>
                    setDraft((d) => ({ ...d, use_value_description: e.target.value }))
                  }
                />
              </label>
              <label>
                <span>Unit</span>
                <input
                  defaultValue={c.use_value.unit}
                  onChange={(e) =>
                    setDraft((d) => ({ ...d, use_value_unit: e.target.value }))
                  }
                />
              </label>
              <label>
                <span>Concrete labour kind</span>
                <input
                  defaultValue={c.concrete_labour.kind}
                  onChange={(e) =>
                    setDraft((d) => ({ ...d, concrete_labour_kind: e.target.value }))
                  }
                />
              </label>
              <div className="span2 form-actions">
                <button className="primary" onClick={save} disabled={busy}>
                  {busy ? "Saving…" : "Save changes"}
                </button>
                {err && <span className="error">{err}</span>}
              </div>
            </div>
          </td>
        </tr>
      )}
      {revealed && relations && (
        <tr>
          <td colSpan={5} className="reveal">
            <div className="reveal-panel">
              <h3>Social relations of {c.name}</h3>
              <p className="reveal-note">{relations.note}</p>
              <p className="labour-statement">
                <strong>{fmtMinutes(relations.labour_per_unit)}</strong> of abstract human labour
                are congealed in one {c.use_value.unit} of {c.name}, produced by{" "}
                {relations.concrete_labour.kind}.
              </p>
              {relations.labour_relations.length > 0 && (
                <ul className="reveal-relations">
                  {relations.labour_relations.map((lr) => (
                    <li key={lr.counterpart.id}>
                      1&nbsp;{c.use_value.unit} {c.name} ≡{" "}
                      <strong>{fmtQty(lr.counterpart_qty)}</strong>{" "}
                      {lr.counterpart.use_value.unit} {lr.counterpart.name} — both expressing{" "}
                      <strong>{fmtMinutes(lr.labour_time)}</strong> of abstract human labour (
                      {c.concrete_labour.kind} ↔ {lr.counterpart.concrete_labour.kind})
                    </li>
                  ))}
                </ul>
              )}
            </div>
          </td>
        </tr>
      )}
    </>
  );
}

function ExchangeRatioPanel({ commodities }: { commodities: Commodity[] }) {
  const [baseId, setBaseId] = useState<string>("");
  const [quoteId, setQuoteId] = useState<string>("");
  const [qty, setQty] = useState<number>(1);
  const [result, setResult] = useState<ExchangeRatio | null>(null);
  const [err, setErr] = useState<string | null>(null);

  const canCompute = useMemo(
    () => baseId && quoteId && baseId !== quoteId && qty > 0,
    [baseId, quoteId, qty]
  );

  async function compute() {
    setErr(null);
    try {
      const r = await api.exchangeRatio(baseId, quoteId, qty);
      setResult(r);
    } catch (e) {
      setErr(e instanceof Error ? e.message : String(e));
    }
  }

  if (commodities.length < 2) return null;

  return (
    <section className="card">
      <h2>Exchange ratio</h2>
      <p className="description">
        "1 quarter corn = x cwt iron" — pick two commodities to derive the equation of values the
        simple form of value expresses.
      </p>
      <div className="form-grid">
        <label>
          <span>Base commodity</span>
          <select value={baseId} onChange={(e) => setBaseId(e.target.value)}>
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
          />
        </label>
        <label>
          <span>Quote commodity</span>
          <select value={quoteId} onChange={(e) => setQuoteId(e.target.value)}>
            <option value="">—</option>
            {commodities.map((c) => (
              <option key={c.id} value={c.id}>
                {c.name}
              </option>
            ))}
          </select>
        </label>
        <div className="span2 form-actions">
          <button className="primary" onClick={compute} disabled={!canCompute}>
            Compute ratio
          </button>
          {err && <span className="error">{err}</span>}
        </div>
      </div>
      {result && (
        <div className="ratio-result">
          <div className="ratio-equation">
            <span>
              <span className="eq-qty">{fmtQty(result.base_qty)}</span>{" "}
              <span className="eq-unit">
                {result.base.use_value.unit} {result.base.name}
              </span>
            </span>
            <span className="eq-sign">=</span>
            <span>
              <span className="eq-qty">{fmtQty(result.quote_qty)}</span>{" "}
              <span className="eq-unit">
                {result.quote.use_value.unit} {result.quote.name}
              </span>
            </span>
          </div>
          <p className="ratio-common">
            common value: {fmtMinutes(result.common_value)} of abstract human labour
          </p>
        </div>
      )}
    </section>
  );
}
