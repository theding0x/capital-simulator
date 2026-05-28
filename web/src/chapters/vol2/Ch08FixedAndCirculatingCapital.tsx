import { useEffect, useState } from "react";
import type { CapitalComponent, FixedCapitalItem, SinkingFundForItem } from "../../types";
import "./Ch08FixedAndCirculatingCapital.css";

const BASE = "/api/v1";

async function fetchComponents(): Promise<CapitalComponent[]> {
  const res = await fetch(`${BASE}/capital-components`);
  if (!res.ok) throw new Error(`${res.status}`);
  return res.json();
}

async function fetchFixedItem(id: string): Promise<FixedCapitalItem | null> {
  const res = await fetch(`${BASE}/fixed-items/${id}`);
  if (res.status === 404) return null;
  if (!res.ok) throw new Error(`${res.status}`);
  return res.json();
}

async function fetchSinkingFund(itemId: string): Promise<SinkingFundForItem | null> {
  const res = await fetch(`${BASE}/fixed-items/${itemId}/sinking-fund`);
  if (res.status === 404) return null;
  if (!res.ok) throw new Error(`${res.status}`);
  return res.json();
}

function pence(n: number): string {
  return `£${(n / 100).toFixed(2)}`;
}

function roleLabel(role: string): string {
  const labels: Record<string, string> = {
    instrument_of_labour: "Instrument of Labour",
    raw_material: "Raw Material",
    auxiliary_material: "Auxiliary Material",
    labour_power: "Labour Power",
  };
  return labels[role] ?? role;
}

function ComponentCard({ component }: { component: CapitalComponent }) {
  return (
    <div className="ch08-card">
      <span className={`ch08-card__kind ch08-card__kind--${component.kind}`}>
        {component.kind}
      </span>
      <div className="ch08-card__role">{roleLabel(component.role)}</div>
      <div className="ch08-card__pence">{pence(component.pence)}</div>
    </div>
  );
}

function FixedItemCard({ item }: { item: FixedCapitalItem }) {
  const [sf, setSF] = useState<SinkingFundForItem | null>(null);

  useEffect(() => {
    fetchSinkingFund(item.id).then(setSF).catch(() => setSF(null));
  }, [item.id]);

  const pct =
    sf && sf.target_pence > 0
      ? Math.min(100, (sf.accumulated_pence / sf.target_pence) * 100)
      : 0;

  const yearsRemaining = item.service_life_minutes / 525_600;

  return (
    <div className="ch08-card">
      <span className="ch08-card__kind ch08-card__kind--fixed">Fixed Item</span>
      <div className="ch08-card__role">{item.description || "Unnamed instrument"}</div>
      <div className="ch08-card__pence">{pence(item.pence_purchased)}</div>
      <div className="ch08-sf">
        Service life: {yearsRemaining.toFixed(1)} yr · {item.wear_model}
      </div>
      {sf && (
        <div className="ch08-sf">
          Sinking fund: {pct.toFixed(1)}% of {pence(sf.target_pence)}
          <div className="ch08-sf__bar">
            <div className="ch08-sf__fill" style={{ width: `${pct}%` }} />
          </div>
        </div>
      )}
    </div>
  );
}

const SEED_ITEM_IDS = [
  "5eed00000000000008f1",
  "5eed00000000000008f2",
];

export function Ch08FixedAndCirculatingCapital() {
  const [components, setComponents] = useState<CapitalComponent[]>([]);
  const [fixedItems, setFixedItems] = useState<FixedCapitalItem[]>([]);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetchComponents()
      .then(setComponents)
      .catch((e) => setError(String(e)));

    Promise.all(SEED_ITEM_IDS.map(fetchFixedItem))
      .then((items) => setFixedItems(items.filter((i): i is FixedCapitalItem => i !== null)))
      .catch(() => {});
  }, []);

  const fixed = components.filter((c) => c.kind === "fixed");
  const circulating = components.filter((c) => c.kind === "circulating");

  return (
    <div>
      {error && <p style={{ color: "var(--red, #f38ba8)" }}>{error}</p>}

      <h3>Fixed Capital</h3>
      <p style={{ color: "var(--ink-muted)", fontSize: "0.85rem" }}>
        Instruments of labour — value transfers piecemeal to the product through
        wear over the item's service life.
      </p>
      {fixed.length === 0 && fixedItems.length === 0 ? (
        <p className="ch08-empty">No fixed capital components recorded.</p>
      ) : (
        <div className="ch08-grid">
          {fixed.map((c) => (
            <ComponentCard key={c.id} component={c} />
          ))}
          {fixedItems.map((item) => (
            <FixedItemCard key={item.id} item={item} />
          ))}
        </div>
      )}

      <h3 style={{ marginTop: "2rem" }}>Circulating Capital</h3>
      <p style={{ color: "var(--ink-muted)", fontSize: "0.85rem" }}>
        Raw materials, auxiliary materials, and labour-power — value passes
        entirely into the product in each cycle.
      </p>
      {circulating.length === 0 ? (
        <p className="ch08-empty">No circulating capital components recorded.</p>
      ) : (
        <div className="ch08-grid">
          {circulating.map((c) => (
            <ComponentCard key={c.id} component={c} />
          ))}
        </div>
      )}
    </div>
  );
}
