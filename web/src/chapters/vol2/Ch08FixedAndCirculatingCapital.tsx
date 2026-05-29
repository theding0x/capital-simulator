import { useEffect, useState } from "react";
import { compositionApi } from "../../api";
import type {
  CapitalComponent,
  FixedCapitalItem,
  SinkingFundForItem,
} from "../../types";
import "./Ch08FixedAndCirculatingCapital.css";

// Seed item IDs the panel pre-loads as illustrations; the backend doesn't
// currently list fixed items, so the panel falls back to these IDs and
// silently skips any that 404. Real stages register their own fixed items
// via compositionApi.registerFixedItem.
const SEED_ITEM_IDS = [
  "5eed00000000000008f1",
  "5eed00000000000008f2",
];

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
    compositionApi
      .getSinkingFund(item.id)
      .then(setSF)
      .catch(() => setSF(null));
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

export function Ch08FixedAndCirculatingCapital() {
  const [components, setComponents] = useState<CapitalComponent[]>([]);
  const [fixedItems, setFixedItems] = useState<FixedCapitalItem[]>([]);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    compositionApi
      .listComponents()
      .then(setComponents)
      .catch((e) => setError(String(e)));

    Promise.all(
      SEED_ITEM_IDS.map((id) =>
        compositionApi.getFixedItem(id).catch(() => null),
      ),
    )
      .then((items) =>
        setFixedItems(items.filter((i): i is FixedCapitalItem => i !== null)),
      )
      .catch(() => {});
  }, []);

  const fixed = components.filter((c) => c.kind === "fixed");
  const circulating = components.filter((c) => c.kind === "circulating");

  return (
    <>
      <KindForRoleInsight />
      <RailwaySubcomponentsDemo />
      <div>
        {error && <p className="error">{error}</p>}

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
      <RicardoCoda />
    </>
  );
}

function KindForRoleInsight() {
  return (
    <section className="v2-ch08-insight">
      <h2 className="v2-ch08-insight-h2">KindForRole — only instruments are fixed</h2>
      <p className="v2-ch08-insight-prose">
        The fixed/circulating distinction is determined by the role the
        component plays in the labour-process, <em>not</em> by its physical
        durability. Only instruments of labour transfer value piecemeal; raw
        material, auxiliary material, and labour-power circulate through each
        cycle whole. (Ricardo's durability-based classification — modelled in
        the engine as <code>error_ricardo_durability_collapse</code> — is what
        this rule rejects.)
      </p>
      <div className="v2-ch08-kind-rules">
        <div className="v2-ch08-rule v2-ch08-rule--fixed">
          <span className="v2-ch08-rule-role">instrument_of_labour</span>
          <span>
            <span className="v2-ch08-rule-arrow">→</span>
            <span className="v2-ch08-rule-kind">fixed</span>
          </span>
        </div>
        <div className="v2-ch08-rule v2-ch08-rule--circulating">
          <span className="v2-ch08-rule-role">raw_material</span>
          <span>
            <span className="v2-ch08-rule-arrow">→</span>
            <span className="v2-ch08-rule-kind">circulating</span>
          </span>
        </div>
        <div className="v2-ch08-rule v2-ch08-rule--circulating">
          <span className="v2-ch08-rule-role">auxiliary_material</span>
          <span>
            <span className="v2-ch08-rule-arrow">→</span>
            <span className="v2-ch08-rule-kind">circulating</span>
          </span>
        </div>
        <div className="v2-ch08-rule v2-ch08-rule--circulating">
          <span className="v2-ch08-rule-role">labour_power</span>
          <span>
            <span className="v2-ch08-rule-arrow">→</span>
            <span className="v2-ch08-rule-kind">circulating</span>
          </span>
        </div>
      </div>
    </section>
  );
}

interface RailwayPart {
  name: string;
  years: number;
  cls: string;
}

const RAILWAY_PARTS: RailwayPart[] = [
  { name: "rails", years: 2, cls: "v2-ch08-railway-seg--rails" },
  { name: "sleepers", years: 8, cls: "v2-ch08-railway-seg--sleepers" },
  { name: "earthworks", years: 80, cls: "v2-ch08-railway-seg--earthworks" },
];

function RailwaySubcomponentsDemo() {
  const total = RAILWAY_PARTS.reduce((acc, p) => acc + p.years, 0);
  return (
    <section className="v2-ch08-railway">
      <h3 className="v2-ch08-railway-h3">A single fixed item, three service lives</h3>
      <p className="v2-ch08-railway-prose">
        A railway is one instrument of labour but its subcomponents wear at
        wildly different rates: rails replace every ~2 years, sleepers ~8,
        earthworks ~80. The bar widths are proportional to the part's service
        life — long-life subcomponents dwarf short-life ones over the asset's
        full life cycle.
      </p>
      <div className="v2-ch08-railway-track" role="img" aria-label="Railway subcomponent service lives">
        {RAILWAY_PARTS.map((p) => (
          <div
            key={p.name}
            className={`v2-ch08-railway-seg ${p.cls}`}
            style={{ flexBasis: `${(p.years / total) * 100}%` }}
            title={`${p.name}: ~${p.years} yr`}
          >
            {p.years}y
          </div>
        ))}
      </div>
      <p className="v2-ch08-railway-legend">
        {RAILWAY_PARTS.map((p) => (
          <span key={p.name}>
            <span
              className="v2-ch08-railway-legend-swatch"
              style={{
                background:
                  p.name === "rails"
                    ? "var(--red)"
                    : p.name === "sleepers"
                    ? "var(--gold-bright)"
                    : "var(--lead)",
              }}
            />
            {p.name} · ~{p.years} yr
          </span>
        ))}
      </p>
    </section>
  );
}

function RicardoCoda() {
  return (
    <aside className="v2-ch08-coda">
      <p className="v2-ch08-coda-quote">
        Ricardo collapses fixed/circulating into a question of relative
        durability and so loses the role-based distinction that grounds value
        transfer. Chapters 10 and 11 walk through that collapse; the
        sinking-fund-reinvestment forward link picks up here when the engine's
        full reproduction schemes land.
      </p>
      <span className="v2-ch08-coda-cite">
        — forward link to Ch.10–11 (Ricardo's collapse)
      </span>
    </aside>
  );
}
