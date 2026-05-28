import { useEffect, useState } from "react";
import type { EconomistAttribution, KnownEconomistError } from "../../types";
import { api } from "../../api";
import "./Ch19FormerPresentations.css";

// Vol. II Ch. 19 — Former Presentations of the Subject.
//
// Marx surveys how Quesnay, Smith, and Ricardo each attempted to analyse the
// circulation of the total annual product and failed — each at a different level
// of abstraction.
//
//   Quesnay (1758) — Tableau Économique: grasped that reproduction is a social
//     circuit, not an individual transaction. Mapped flows between classes.
//     Errors: limits production to agriculture; lacks a value-theory foundation.
//
//   Smith (1776) — revenue dogma: resolves the entire annual product into wages,
//     profit, and rent, omitting constant-capital replacement entirely. Three
//     consequences: no c term → no balancing equation; no replacement fund for
//     means of production; infinite regress when he tries to fix it.
//
//   Ricardo (1817): sees Smith's individual-capital mistake but never generalises
//     to total social capital. Cannot derive I(v+s) = II(c).
//
// This chapter introduces no new simulation logic. It documents the intellectual
// genealogy that makes Chs. 20–21 (the two-department reproduction schemes)
// intelligible.

// ── error metadata ────────────────────────────────────────────────────────────

const ERROR_DESC: Partial<
  Record<KnownEconomistError, { label: string; desc: string; refutation: string }>
> = {
  // Ch. 10–11 errors (displayed in Ch. 10 and Ch. 11 panels; repeated here for
  // completeness when the full attribution list is shown)
  error_smith_conflation: {
    label: "Fixed/Circulating ≠ Constant/Variable",
    desc: "Smith conflates fixed/circulating (productive-capital form) with constant/variable (value-expansion role).",
    refutation: "Refuted by Ch. 8 KindForRole invariant.",
  },
  error_smith_circulation_capital_conflation: {
    label: "Money + Commodity Capital ≠ 'Circulating Capital'",
    desc: "Smith groups circuit-forms M and C′ under 'circulating capital', collapsing the circuit-form distinction.",
    refutation: "Refuted by Ch. 2–4: circuit-form distinction is prior and independent.",
  },
  error_smith_revenue_in_capital: {
    label: "Revenue Flows ≠ Capital",
    desc: "Smith imports wages, profit, and rent (revenue exiting C′—M′) back into the capital circuit.",
    refutation: "Refuted by Ch. 2 RevenueCircuit invariant.",
  },
  error_ricardo_durability_collapse: {
    label: "Fixed/Circulating = Durability",
    desc: "Ricardo reduces the distinction to physical durability, not role in the labour-process.",
    refutation: "Refuted by Ch. 8: kind is determined by role, not material durability.",
  },
  error_ricardo_conflation: {
    label: "Fixed/Circulating ≠ Constant/Variable (Ricardo)",
    desc: "Ricardo repeats Smith's fixed/circulating ↔ constant/variable conflation.",
    refutation: "Refuted by Ch. 8 KindForRole invariant.",
  },
  error_ricardo_fixed_capital_price_explanation: {
    label: "Fixed-Capital Intensity → Price Effect (mis-located cause)",
    desc: "Ricardo locates the price effect in the fixed/circulating ratio, not the constant/variable ratio.",
    refutation: "Full resolution deferred to Vol. III (prices of production).",
  },
  error_ricardo_no_aggregate_turnover: {
    label: "No Aggregate-Turnover Formula",
    desc: "Ricardo lacks Ch. 9's mean-average turnover and therefore mis-states periodicity of capital returns.",
    refutation: "Refuted by Ch. 9 aggregate-turnover scheme.",
  },
  error_ricardo_no_value_revolution: {
    label: "Value Treated as Static",
    desc: "Ricardo cannot account for a mid-turnover change in the value of means of production.",
    refutation: "Refuted by Ch. 4 ValueRevolutionEvent.",
  },
  // Ch. 19 errors
  error_quesnay_productive_sterile_division: {
    label: "Only Agriculture is Productive",
    desc: "Quesnay divides labour into a 'productive class' (agriculture) and a 'sterile class' (manufacture). Correct in locating surplus in production; wrong to limit it to agriculture.",
    refutation:
      "Refuted by Vol. I Ch. 7: manufacture creates value and surplus-value through the labour-process regardless of sector.",
  },
  error_quesnay_missing_value_theory: {
    label: "No Labour-Theory-of-Value Foundation",
    desc: "The Tableau maps the social circuit without a theory of why the net product arises. Its flows prefigure the two-department scheme (Chs. 20–21) but cannot be made rigorous without the value categories of Vol. I.",
    refutation:
      "Completed by Marx's Vol. I value analysis: constant capital c, variable capital v, surplus-value s replace the produit net.",
  },
  error_smith_revenue_dogma: {
    label: "Revenue Dogma: Annual Product = v + p + r",
    desc: "Smith resolves the entire value of the annual product into wages, profit, and rent. There is no constant-capital term. The equation I(v+s) = II(c) is inexpressible — one side of the balance has no referent.",
    refutation:
      "Refuted by the two-department scheme (Chs. 20–21): Department I must reproduce constant capital, not resolve it into revenue.",
  },
  error_smith_missing_constant_replacement: {
    label: "No Replacement Fund for Means of Production",
    desc: "If the entire product = wages + profit + rent, there is no fund for replacing the means of production consumed during the year. The economy cannot physically reproduce its productive apparatus.",
    refutation:
      "Department I of the reproduction scheme (Chs. 20–21) is precisely the department that produces means of production consumed as constant capital across both departments.",
  },
  error_smith_labour_regress: {
    label: "Infinite Regress of Constant Capital into Revenue",
    desc: "Smith tries to reduce every input to wages, profit, and rent paid to produce it — those inputs in turn to earlier wages, profit, and rent, and so on. The regress never terminates.",
    refutation:
      "Refuted by Marx's synchronic value analysis: constant capital is a present, ongoing component of every act of production, not a historical residue.",
  },
  error_ricardo_incomplete_reproduction: {
    label: "Sees Smith's Error, Cannot Fix It",
    desc: "Ricardo observes (at the individual-capital level) that replacing worn machinery is a cost, not revenue. He never applies this to total social capital and cannot derive I(v+s) = II(c).",
    refutation:
      "Completed by Marx in Chs. 20–21: the two-department scheme shows how constant-capital replacement at the social level operates through inter-department exchange.",
  },
};

// ── genealogy data ────────────────────────────────────────────────────────────

interface GenealogyNode {
  theorist: string;
  year: number;
  work: string;
  anticipates: string;
  errorCount: number;
}

const GENEALOGY: GenealogyNode[] = [
  {
    theorist: "Quesnay",
    year: 1758,
    work: "Tableau Économique",
    anticipates: "Reproduction is a social circuit",
    errorCount: 2,
  },
  {
    theorist: "Smith",
    year: 1776,
    work: "Wealth of Nations",
    anticipates: "Fixed/circulating distinction (imperfectly)",
    errorCount: 6,
  },
  {
    theorist: "Ricardo",
    year: 1817,
    work: "Principles of Political Economy",
    anticipates: "Individual-capital replacement cost",
    errorCount: 6,
  },
  {
    theorist: "Marx",
    year: 1885,
    work: "Capital Vol. II (Ch. 20–21)",
    anticipates: "Two-department scheme: I(v+s) = II(c)",
    errorCount: 0,
  },
];

// ── components ────────────────────────────────────────────────────────────────

function GenealogyDiagram() {
  return (
    <section className="ch19-genealogy">
      <h3>Theoretical Genealogy</h3>
      <div className="ch19-genealogy-row">
        {GENEALOGY.map((node, i) => (
          <div key={node.theorist} className="ch19-genealogy-step">
            <div
              className={`ch19-genealogy-node${node.theorist === "Marx" ? " ch19-genealogy-node--marx" : ""}`}
            >
              <span className="ch19-genealogy-theorist">{node.theorist}</span>
              <span className="ch19-genealogy-year">{node.year}</span>
              <span className="ch19-genealogy-work">{node.work}</span>
              <span className="ch19-genealogy-anticipates">
                ↗ {node.anticipates}
              </span>
              {node.errorCount > 0 && (
                <span className="ch19-genealogy-errors">
                  {node.errorCount} theoretical error
                  {node.errorCount !== 1 ? "s" : ""}
                </span>
              )}
            </div>
            {i < GENEALOGY.length - 1 && (
              <span className="ch19-genealogy-arrow">→</span>
            )}
          </div>
        ))}
      </div>
    </section>
  );
}

function SmithRevenueDogmaView({
  attributions,
}: {
  attributions: EconomistAttribution[];
}) {
  const smith = attributions.find((a) => a.theorist === "Smith");
  const ch19ErrorCodes: KnownEconomistError[] = [
    "error_smith_revenue_dogma",
    "error_smith_missing_constant_replacement",
    "error_smith_labour_regress",
  ];
  const errors = smith
    ? smith.errors.filter((e) => ch19ErrorCodes.includes(e))
    : ch19ErrorCodes;

  return (
    <section className="ch19-smith-dogma">
      <h3>Smith's Revenue Dogma</h3>
      <div className="ch19-dogma-equation">
        <div className="ch19-equation-row">
          <span className="ch19-eq-label">Smith:</span>
          <span className="ch19-eq-formula">
            Annual Product = <em>v</em> + <em>p</em> + <em>r</em>
          </span>
          <span className="ch19-eq-verdict ch19-eq-verdict--wrong">✗</span>
        </div>
        <div className="ch19-equation-row">
          <span className="ch19-eq-label">Marx:</span>
          <span className="ch19-eq-formula">
            Annual Product = <em>c</em> + <em>v</em> + <em>s</em>
          </span>
          <span className="ch19-eq-verdict ch19-eq-verdict--right">✓</span>
        </div>
        <p className="ch19-dogma-note">
          Smith's formula has no <em>c</em> term. Without constant capital there
          is no replacement fund for means of production, and the two-department
          balance I(<em>v</em>+<em>s</em>) = II(<em>c</em>) is inexpressible.
        </p>
      </div>
      <div className="ch19-error-list">
        {errors.map((e) => {
          const meta = ERROR_DESC[e];
          if (!meta) return null;
          return (
            <div key={e} className="ch19-error-item">
              <code className="ch19-error-code">{e}</code>
              <p className="ch19-error-desc">
                <strong>{meta.label}.</strong> {meta.desc}
              </p>
              <p className="ch19-error-refutation">{meta.refutation}</p>
            </div>
          );
        })}
      </div>
    </section>
  );
}

function HistoricalErrorsView({
  attributions,
}: {
  attributions: EconomistAttribution[];
}) {
  const ch19AttrIds = ["5eed000000000000001901", "5eed000000000000001902"];
  const ch19Attrs = attributions.filter((a) => ch19AttrIds.includes(a.id));

  return (
    <section className="ch19-historical-errors">
      <h3>Quesnay and Ricardo: Reproduction Errors</h3>
      {ch19Attrs.map((attr) => (
        <div key={attr.id} className="ch19-attr-block">
          <h4>
            {attr.theorist} ({attr.edition_year}) — {attr.concept}
          </h4>
          <p className="ch19-attr-anticipates">
            Anticipates: <code>{attr.anticipates}</code>
          </p>
          <div className="ch19-error-list">
            {attr.errors.map((e) => {
              const meta = ERROR_DESC[e];
              if (!meta) return <code key={e}>{e}</code>;
              return (
                <div key={e} className="ch19-error-item">
                  <code className="ch19-error-code">{e}</code>
                  <p className="ch19-error-desc">
                    <strong>{meta.label}.</strong> {meta.desc}
                  </p>
                  <p className="ch19-error-refutation">{meta.refutation}</p>
                </div>
              );
            })}
          </div>
        </div>
      ))}
    </section>
  );
}

const QUOTES: { text: string; citation: string }[] = [
  {
    text: "Quesnay's Tableau Économique … shows, in a few broad outlines, how the result of annual production, distributed among the different classes of society, circulates through consumption, and thus reproduces its starting-point, the annual production.",
    citation: "Marx, Capital Vol. II, Ch. 19",
  },
  {
    text: "Adam Smith resolves the entire value of the annual social product into revenue — that is, into wages, profit, and ground-rent — and thus entirely omits the constant portion of value which must replace the constant capital that has been consumed.",
    citation: "Marx, Capital Vol. II, Ch. 19",
  },
  {
    text: "Ricardo's merit consists not in his having examined the replacement of the constant capital, but in his having recognised the necessity of doing so … yet he never succeeds in making it the basis of a systematic analysis.",
    citation: "Marx, Capital Vol. II, Ch. 19 (paraphrase)",
  },
];

// ── main export ───────────────────────────────────────────────────────────────

export function Ch19FormerPresentations() {
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

  if (loading) return <p className="ch19-loading">Loading…</p>;
  if (fetchError) return <p className="ch19-error">Error: {fetchError}</p>;

  return (
    <div className="ch19-panel">
      <GenealogyDiagram />
      <SmithRevenueDogmaView attributions={attributions} />
      <HistoricalErrorsView attributions={attributions} />
      <section className="ch19-quotes">
        <h3>From the Chapter</h3>
        {QUOTES.map((q, i) => (
          <blockquote key={i} className="ch19-quote">
            <p>"{q.text}"</p>
            <cite>— {q.citation}</cite>
          </blockquote>
        ))}
      </section>
    </div>
  );
}
