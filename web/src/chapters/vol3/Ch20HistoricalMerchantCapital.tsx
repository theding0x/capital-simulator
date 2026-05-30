import { useEffect, useState } from "react";
import { financeApi } from "../../api";
import type {
  DevelopmentStage,
  HistoricalMerchantCapitalResponse,
} from "../../types";
import "./Ch20HistoricalMerchantCapital.css";

function bpLabel(bp: number): string {
  return `${(bp / 100).toFixed(2)}%`;
}

const STAGE_LABELS: Record<DevelopmentStage, string> = {
  1: "Pre-capitalist",
  2: "Transition",
  3: "Subordinated",
};

const PROFIT_SOURCE_OPTIONS = [
  { value: "unequal_exchange",       label: "Unequal exchange" },
  { value: "monopoly_trade",         label: "Monopoly trade" },
  { value: "carrying_monopoly",      label: "Carrying monopoly" },
  { value: "industrial_subordination", label: "Industrial subordination" },
];

const ALL_STAGES: DevelopmentStage[] = [1, 2, 3];

// Seed fixtures illustrating the inverse-development law (invariant 2):
// SubordinationIndex rises with Stage.
const SEED_EXEMPLARS = [
  { name: "Venice carrying trade",              stage: 1 as DevelopmentStage, profit_source: "unequal_exchange",       wage_labour: false, subordination_index: 500  },
  { name: "Genoa merchant republic",            stage: 1 as DevelopmentStage, profit_source: "monopoly_trade",         wage_labour: false, subordination_index: 800  },
  { name: "Dutch carrying trade",               stage: 2 as DevelopmentStage, profit_source: "carrying_monopoly",      wage_labour: false, subordination_index: 3000 },
  { name: "Subordination under industrial capital", stage: 3 as DevelopmentStage, profit_source: "industrial_subordination", wage_labour: true, subordination_index: 7000 },
];

function stageLabel(s: DevelopmentStage): string {
  return STAGE_LABELS[s] ?? String(s);
}

function profitSourceLabel(src: string): string {
  return PROFIT_SOURCE_OPTIONS.find((o) => o.value === src)?.label ?? src;
}

function stageClass(s: DevelopmentStage): string {
  if (s === 1) return "v3-ch20-stage--pre";
  if (s === 2) return "v3-ch20-stage--trans";
  return "v3-ch20-stage--sub";
}

export function Ch20HistoricalMerchantCapital() {
  const [items, setItems] = useState<HistoricalMerchantCapitalResponse[]>([]);
  const [listError, setListError] = useState<string | null>(null);

  // Create-form state.
  const [name, setName] = useState("");
  const [stage, setStage] = useState<DevelopmentStage>(1);
  const [profitSource, setProfitSource] = useState("unequal_exchange");
  const [wageLabour, setWageLabour] = useState(false);
  const [subordinationIndex, setSubordinationIndex] = useState(500);
  const [createError, setCreateError] = useState<string | null>(null);

  // Client-side guard: pre-capitalist stage cannot employ wage-labour.
  const wageLabourForbidden = stage === 1 && wageLabour;

  async function refreshItems() {
    try {
      const data = await financeApi.listHistoricalMerchantCapitals();
      setItems(data);
    } catch (e) {
      setListError(String(e));
    }
  }

  useEffect(() => {
    refreshItems().catch(() => {});
  }, []);

  // Reset wage_labour to false when stage switches to pre-capitalist.
  function handleStageChange(s: DevelopmentStage) {
    setStage(s);
    if (s === 1) setWageLabour(false);
  }

  async function handleCreate() {
    setCreateError(null);
    try {
      await financeApi.createHistoricalMerchantCapital({
        name,
        stage,
        profit_source: profitSource,
        wage_labour: wageLabour,
        subordination_index: subordinationIndex,
      });
      await refreshItems();
      setName("");
    } catch (e) {
      setCreateError(String(e));
    }
  }

  // Sort items by stage ascending for the timeline view.
  const sorted = [...items].sort((a, b) =>
    a.stage !== b.stage ? a.stage - b.stage : a.created_at.localeCompare(b.created_at),
  );

  return (
    <div className="v3-ch20">
      <section className="v3-ch20-intro">
        <p>
          Merchant's capital is historically the oldest free form of capital — it predates
          the capitalist mode of production. Marx traces three stages: (1){" "}
          <strong>pre-capitalist</strong> forms (Venice, Genoa) that appropriate surplus
          through unequal exchange, <em>never</em> wage-labour; (2) a{" "}
          <strong>transition</strong> stage as industrial capital expands; (3){" "}
          <strong>subordination</strong> under industrial capital, where the merchant
          merely facilitates commodity circulation and participates in the general rate of
          profit. The <em>inverse-development law</em>: as industrial capitalism advances,
          the SubordinationIndex of merchant capital rises.
        </p>
      </section>

      {/* Inverse-development law illustration using seed fixtures */}
      <section className="v3-ch20-section" aria-label="Inverse-development law">
        <h4>Inverse-development law — historical exemplars</h4>
        <p className="v3-ch20-subintro">
          SubordinationIndex (bp) rises as the stage advances — Marx's concrete forms.
        </p>
        <div className="v3-ch20-timeline">
          {SEED_EXEMPLARS.map((ex) => (
            <div key={ex.name} className={`v3-ch20-timeline-item ${stageClass(ex.stage)}`}>
              <div className="v3-ch20-timeline-stage">{stageLabel(ex.stage)}</div>
              <div className="v3-ch20-timeline-name">{ex.name}</div>
              <div className="v3-ch20-timeline-source">{profitSourceLabel(ex.profit_source)}</div>
              <div className="v3-ch20-timeline-index">
                <span className="v3-ch20-index-label">Subordination</span>
                <span className="v3-ch20-index-value">{bpLabel(ex.subordination_index)}</span>
              </div>
              <div className="v3-ch20-timeline-wage">
                {ex.wage_labour ? "Employs wage-labour" : "No wage-labour"}
              </div>
            </div>
          ))}
        </div>
      </section>

      {/* Create form */}
      <section className="v3-ch20-section" aria-label="Record historical merchant capital">
        <h4>Record a historical form of merchant capital</h4>
        <div className="v3-ch20-form">
          <label>
            Name
            <input
              type="text"
              value={name}
              placeholder="e.g. Venice carrying trade"
              onChange={(e) => setName(e.target.value)}
            />
          </label>
          <label>
            Development stage
            <select
              value={stage}
              onChange={(e) => handleStageChange(Number(e.target.value) as DevelopmentStage)}
            >
              {ALL_STAGES.map((s) => (
                <option key={s} value={s}>
                  {s} — {stageLabel(s)}
                </option>
              ))}
            </select>
          </label>
          <label>
            Profit source
            <select value={profitSource} onChange={(e) => setProfitSource(e.target.value)}>
              {PROFIT_SOURCE_OPTIONS.map((o) => (
                <option key={o.value} value={o.value}>
                  {o.label}
                </option>
              ))}
            </select>
          </label>
          <label className={`v3-ch20-wage-label ${stage === 1 ? "v3-ch20-wage-label--disabled" : ""}`}>
            <input
              type="checkbox"
              checked={wageLabour}
              disabled={stage === 1}
              onChange={(e) => setWageLabour(e.target.checked)}
            />
            Employs wage-labour
          </label>
          {stage === 1 && (
            <p className="v3-ch20-hint">
              Pre-capitalist merchant capital cannot employ wage-labour (invariant 1).
            </p>
          )}
          <label>
            Subordination index (bp, 0–10000)
            <input
              type="number"
              min={0}
              max={10000}
              value={subordinationIndex}
              onChange={(e) => setSubordinationIndex(Number(e.target.value))}
            />
            <span className="v3-ch20-bp-hint">{bpLabel(subordinationIndex)} of total social capital</span>
          </label>
          <button
            type="button"
            className="v3-ch20-btn"
            disabled={!name || wageLabourForbidden}
            onClick={handleCreate}
          >
            Record
          </button>
        </div>
        {createError && <p className="v3-ch20-error">{createError}</p>}
      </section>

      {/* Timeline of stored records */}
      {listError && <p className="v3-ch20-error">{listError}</p>}
      {sorted.length > 0 && (
        <section className="v3-ch20-section" aria-label="Historical merchant capital records">
          <h4>Historical forms ({sorted.length}) — ordered by stage</h4>
          <div className="v3-ch20-timeline v3-ch20-timeline--live">
            {sorted.map((hm) => (
              <div key={hm.id} className={`v3-ch20-timeline-item ${stageClass(hm.stage)}`}>
                <div className="v3-ch20-timeline-stage">{stageLabel(hm.stage)}</div>
                <div className="v3-ch20-timeline-name">{hm.name}</div>
                <div className="v3-ch20-timeline-source">{profitSourceLabel(hm.profit_source)}</div>
                <div className="v3-ch20-timeline-index">
                  <span className="v3-ch20-index-label">Subordination</span>
                  <span className="v3-ch20-index-value">{bpLabel(hm.subordination_index)}</span>
                </div>
                <div className="v3-ch20-timeline-wage">
                  {hm.wage_labour ? "Employs wage-labour" : "No wage-labour"}
                </div>
              </div>
            ))}
          </div>
        </section>
      )}
    </div>
  );
}
