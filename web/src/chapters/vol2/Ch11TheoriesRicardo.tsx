import { useEffect, useState } from "react";
import type { EconomistAttribution, KnownEconomistError } from "../../types";
import { api } from "../../api";
import "./Ch11TheoriesRicardo.css";

// Vol. II Ch. 11 — Theories of Fixed and Circulating Capital: Ricardo.

// Ricardo's four-fold durability table (Principles, 1817).
// He distinguishes fixed from circulating capital purely by physical longevity.
const RICARDO_DURABILITY_TABLE: { term: string; duration: string }[] = [
  { term: "plough, farm buildings, implements", duration: "many years → fixed" },
  { term: "seeds, manure, wages-goods", duration: "one season → circulating" },
  { term: "machinery in manufacture", duration: "many years → fixed" },
  { term: "raw materials, wage-fund", duration: "one period → circulating" },
];

// Marx's Ch. 8 role-derives-kind map.
const MARX_ROLE_MAP: { role: string; kind: string; correct: boolean }[] = [
  { role: "instrument of labour (plough, machine)", kind: "Fixed Capital", correct: true },
  { role: "raw material (cotton, iron ore)", kind: "Circulating Capital", correct: true },
  { role: "auxiliary material (coal, oil)", kind: "Circulating Capital", correct: true },
  { role: "labour-power (variable capital)", kind: "Circulating Capital — ≠ variable", correct: false },
];

const ERROR_DESC: Record<
  KnownEconomistError,
  { label: string; desc: string; refutation: string } | undefined
> = {
  error_smith_conflation: undefined,
  error_smith_circulation_capital_conflation: undefined,
  error_smith_revenue_in_capital: undefined,
  error_ricardo_durability_collapse: {
    label: "Durability ≠ Form of Circulation",
    desc: "Ricardo treats fixed/circulating as purely a matter of physical durability: tools last many periods (fixed); wages-goods are consumed at once (circulating). This collapses the form-of-circulation invariant — kind is determined by the role a capital component plays in the labour-process, not by how long it physically lasts.",
    refutation:
      "Refuted by Ch. 8 KindForRole invariant: a building that serves as a warehouse and a building that serves as a machine-room are both durable, but their circuit-form differs by role.",
  },
  error_ricardo_conflation: {
    label: "Fixed ≠ Constant; Circulating ≠ Variable",
    desc: "Ricardo repeats Smith's conflation: instruments of labour (fixed capital) are treated as equivalent to constant capital, and the wage-fund (circulating capital) as equivalent to variable capital. Fixed capital is always constant capital, but the converse does not hold — raw materials are constant and circulating.",
    refutation:
      "Refuted by Ch. 8 KindForRole invariant: the constant/variable distinction turns on whether capital-value expands in the labour-process; the fixed/circulating distinction turns on how value is transferred to the product.",
  },
  error_ricardo_fixed_capital_price_explanation: {
    label: "Organic Composition → Prices (Mis-located Cause)",
    desc: "Ricardo correctly observes that when the wage rate changes, industries with higher fixed-capital intensity show larger price movements than those with lower intensity. He attributes this to the fixed/circulating ratio. The actual cause is the constant/variable ratio and the extraction of surplus-value — the full resolution requires the average rate of profit.",
    refutation:
      "Partially anticipated by Ch. 8 (constant/variable distinction) and Ch. 9 (aggregate turnover). Full resolution deferred to Vol. III Part Two: prices of production, average rate of profit.",
  },
  error_ricardo_no_aggregate_turnover: {
    label: "No Aggregate-Turnover Formula",
    desc: "Ricardo lacks the Ch. 9 aggregate-turnover formula. Without the mean-average turnover number across heterogeneous capital components, he cannot state the periodicity of capital returns and so cannot distinguish how fixed vs. circulating portions of the same capital contribute to the annual rate of profit.",
    refutation:
      "Refuted constructively by Ch. 9 AggregateTurnover: AnnualTurnedOverPence may exceed AdvancedTotalPence once the mean-average turnover exceeds unity — a result Ricardo's framework cannot produce.",
  },
  error_ricardo_no_value_revolution: {
    label: "Value Is Not Static within a Turnover",
    desc: "Ricardo treats the value of means of production as fixed for the duration of a turnover. He cannot account for a change in the value of raw materials mid-turnover that alters MoneyCapitalTiedUp without altering the physical volume of output.",
    refutation:
      "Refuted by Ch. 4 ValueRevolutionEvent: a DirectionPence > 0 event on a Means-of-Production component increases MoneyCapitalTiedUp independently of the labour-process — invisible to Ricardo's static value assumption.",
  },
};

function DurabilityComparisonView() {
  return (
    <div className="ch11-comparison">
      <section>
        <p className="ch11-col-header">Ricardo (1817) — Physical Durability</p>
        <div className="ch11-concept-list">
          {RICARDO_DURABILITY_TABLE.map((row) => (
            <div key={row.term} className="ch11-concept-row">
              <span className="ch11-concept-term">{row.term}</span>
              <span className="ch11-arrow">→</span>
              <span className="ch11-concept-target conflation">{row.duration}</span>
            </div>
          ))}
        </div>
      </section>
      <section>
        <p className="ch11-col-header">Marx — Ch. 8 Role-Derives-Kind</p>
        <div className="ch11-concept-list">
          {MARX_ROLE_MAP.map((row) => (
            <div key={row.role} className="ch11-concept-row">
              <span className="ch11-concept-term">{row.role}</span>
              <span className="ch11-arrow">→</span>
              <span
                className={
                  "ch11-concept-target" + (row.correct ? "" : " conflation")
                }
              >
                {row.kind}
              </span>
            </div>
          ))}
        </div>
      </section>
    </div>
  );
}

function RicardoErrorsView({
  attributions,
}: {
  attributions: EconomistAttribution[];
}) {
  const ricardo = attributions.find((a) => a.theorist === "Ricardo");
  if (!ricardo || ricardo.errors.length === 0) return null;
  return (
    <section className="ch11-errors-section">
      <h3>Where Ricardo Goes Wrong</h3>
      <div className="ch11-error-list">
        {ricardo.errors.map((e) => {
          const meta = ERROR_DESC[e];
          if (!meta) return null;
          return (
            <div key={e} className="ch11-error-item">
              <code className="ch11-error-code">{e}</code>
              <p className="ch11-error-desc">
                <strong>{meta.label}.</strong> {meta.desc}
              </p>
              <p className="ch11-error-refutation">{meta.refutation}</p>
            </div>
          );
        })}
      </div>
    </section>
  );
}

function DeferredSection() {
  return (
    <section className="ch11-deferred-section">
      <h3>Coming in Vol. III</h3>
      <div className="ch11-deferred-box">
        <p className="ch11-deferred-label">Deferred — Vol. III Part Two</p>
        <p>
          Ricardo's observation that differential fixed-capital intensity produces
          differential price movements is <em>correct as a phenomenon</em>.
          Resolving <em>why</em> requires the average rate of profit, the
          formation of prices of production, and the redistribution of
          surplus-value across branches of production.
        </p>
        <p>
          These concepts belong to Vol. III Part Two (Ch. 8–12). The{" "}
          <code>error_ricardo_fixed_capital_price_explanation</code> entry will
          be linked to the prices-of-production panel when that chapter lands.
        </p>
      </div>
    </section>
  );
}

export function Ch11TheoriesRicardo() {
  const [attributions, setAttributions] = useState<EconomistAttribution[]>([]);
  const [loading, setLoading] = useState(true);
  const [fetchError, setFetchError] = useState<string | null>(null);

  useEffect(() => {
    api
      .listEconomistAttributions("Ricardo")
      .then(setAttributions)
      .catch((e: unknown) =>
        setFetchError(
          e instanceof Error ? e.message : "Failed to load attributions"
        )
      )
      .finally(() => setLoading(false));
  }, []);

  if (loading) return <p className="ch11-loading">Loading…</p>;
  if (fetchError) return <p className="ch11-error">Error: {fetchError}</p>;

  return (
    <div className="ch11-panel">
      <DurabilityComparisonView />
      <RicardoErrorsView attributions={attributions} />
      <DeferredSection />
    </div>
  );
}
