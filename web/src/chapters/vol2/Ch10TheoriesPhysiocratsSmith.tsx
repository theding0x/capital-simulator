import { useEffect, useState } from "react";
import type { EconomistAttribution, KnownEconomistError } from "../../types";
import { api } from "../../api";
import "./Ch10TheoriesPhysiocratsSmith.css";

// Vol. II Ch. 10 — Theories of Fixed and Circulating Capital:
// The Physiocrats and Adam Smith.

const ANTICIPATES_LABEL: Record<string, string> = {
  fixed_capital_item: "Fixed Capital",
  circulating: "Circulating Capital",
};

const ERROR_DESC: Partial<
  Record<KnownEconomistError, { label: string; desc: string; refutation: string }>
> = {
  error_smith_conflation: {
    label: "Fixed/Circulating ≠ Constant/Variable",
    desc: "Smith conflates the fixed/circulating distinction (which belongs to productive capital) with the constant/variable distinction (which turns on value-expansion). A fixed machine is constant capital, but not all constant capital is fixed.",
    refutation:
      "Refuted by Ch. 8 KindForRole invariant: kind is determined by role in the labour-process, independent of how value is transferred.",
  },
  error_smith_circulation_capital_conflation: {
    label: "Money + Commodity Capital ≠ 'Circulating Capital'",
    desc: "Smith groups money-capital and commodity-capital under 'circulating capital' because they circulate. But these are circuit-forms of capital (M and C′), not sub-divisions of productive capital's fixed/circulating pair.",
    refutation:
      "Refuted by Ch. 2–4: the circuit-form distinction (M, P, C′) is prior to and independent of the productive-capital distinction.",
  },
  error_smith_revenue_in_capital: {
    label: "Revenue Flows ≠ Capital",
    desc: "Smith treats wages, profit, and rent (revenue flows exiting the circuit) as elements of 'circulating capital', importing them back into the capital circuit. Revenue exits at C′—M′ and is not re-advanced as M—C(Lp+Mp).",
    refutation:
      "Refuted by Ch. 2 RevenueCircuit invariant: revenue rows do not contribute to ConstantPence or VariablePence in the capital circuit.",
  },
};

// Verbatim chapter-text quotes (Marx, Vol. II, Ch. 10).
const QUOTES: { text: string; citation: string }[] = [
  {
    text: "Quesnay's Tableau Économique shows how a definite net product circulates … he distinguishes avances primitives (capital invested once for several years, i.e. fixed capital) from avances annuelles (capital invested anew each year, i.e. circulating capital).",
    citation: "Marx, Capital Vol. II, Ch. 10 (Quesnay's anticipation)",
  },
  {
    text: "Adam Smith … confuses the distinctions between fixed and circulating capital with the distinction between constant and variable capital, and thus begets endless confusion.",
    citation: "Marx, Capital Vol. II, Ch. 10 (Smith's first error)",
  },
  {
    text: "According to Adam Smith, the circulating capital consists of money, provisions, materials, and finished work. … He thus mixes up the distinction between fixed and circulating capital with the distinction between capital and revenue.",
    citation: "Marx, Capital Vol. II, Ch. 10 (Smith's third error)",
  },
];

function AnticipationsView({
  attributions,
}: {
  attributions: EconomistAttribution[];
}) {
  const quesnay = attributions.filter((a) => a.theorist === "Quesnay");
  return (
    <section>
      <p className="ch10-col-header">Quesnay — Tableau Économique (1758)</p>
      <div className="ch10-concept-list">
        {quesnay.map((a) => (
          <div key={a.id} className="ch10-concept-row">
            <span className="ch10-concept-term">{a.concept}</span>
            <span className="ch10-arrow">→</span>
            <span className="ch10-concept-target">
              {ANTICIPATES_LABEL[a.anticipates] ?? a.anticipates}
            </span>
          </div>
        ))}
      </div>
    </section>
  );
}

function SmithErrorsView({
  attributions,
}: {
  attributions: EconomistAttribution[];
}) {
  const smith = attributions.find((a) => a.theorist === "Smith");
  if (!smith || smith.errors.length === 0) return null;
  return (
    <section className="ch10-errors-section">
      <h3>Where Smith Goes Wrong</h3>
      <div className="ch10-error-list">
        {smith.errors.map((e) => {
          const meta = ERROR_DESC[e];
          if (!meta) return null;
          return (
            <div key={e} className="ch10-error-item">
              <code className="ch10-error-code">{e}</code>
              <p className="ch10-error-desc">
                <strong>{meta.label}.</strong> {meta.desc}
              </p>
              <p className="ch10-error-refutation">{meta.refutation}</p>
            </div>
          );
        })}
      </div>
    </section>
  );
}

export function Ch10TheoriesPhysiocratsSmith() {
  const [attributions, setAttributions] = useState<EconomistAttribution[]>([]);
  const [loading, setLoading] = useState(true);
  const [fetchError, setFetchError] = useState<string | null>(null);

  useEffect(() => {
    api
      .listEconomistAttributions()
      .then(setAttributions)
      .catch((e: unknown) =>
        setFetchError(
          e instanceof Error ? e.message : "Failed to load attributions"
        )
      )
      .finally(() => setLoading(false));
  }, []);

  if (loading) return <p className="ch10-loading">Loading…</p>;
  if (fetchError) return <p className="ch10-error">Error: {fetchError}</p>;

  return (
    <div className="ch10-panel">
      <div className="ch10-anticipations">
        <AnticipationsView attributions={attributions} />
        <section>
          <p className="ch10-col-header">Marx — Ch. 8 Distinction</p>
          <div className="ch10-concept-list">
            <div className="ch10-concept-row">
              <span className="ch10-concept-term">avances primitives</span>
              <span className="ch10-arrow">≈</span>
              <span className="ch10-concept-target">
                Fixed Capital (instrument of labour)
              </span>
            </div>
            <div className="ch10-concept-row">
              <span className="ch10-concept-term">avances annuelles</span>
              <span className="ch10-arrow">≈</span>
              <span className="ch10-concept-target">
                Circulating Capital (raw + auxiliary + labour-power)
              </span>
            </div>
          </div>
        </section>
      </div>

      <SmithErrorsView attributions={attributions} />

      <section className="ch10-quotes">
        <h3>From the Chapter</h3>
        {QUOTES.map((q, i) => (
          <blockquote key={i} className="ch10-quote">
            <p>"{q.text}"</p>
            <cite>— {q.citation}</cite>
          </blockquote>
        ))}
      </section>
    </div>
  );
}
