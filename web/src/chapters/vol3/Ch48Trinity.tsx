import { useEffect, useState } from "react";
import { financeApi } from "../../api";
import type {
  RevenueFetishForm,
  RevenueSourceKind,
  TrinityFormulaResponse,
} from "../../types";
import "./Ch48Trinity.css";

const SOURCE_LABELS: Record<RevenueSourceKind, string> = {
  1: "Capital — interest + profit of enterprise",
  2: "Land — ground-rent",
  3: "Labour — wages",
};

export function Ch48Trinity() {
  const [fetishForms, setFetishForms] = useState<RevenueFetishForm[]>([]);
  const [listError, setListError] = useState<string | null>(null);

  // Trinity-formula builder fields (defaults = Ch. 48 seed fixture).
  const [capApparent, setCapApparent] = useState(20);
  const [capActual, setCapActual] = useState(17); // profit of enterprise 12 + interest 5 (rent already deducted)
  const [landApparent, setLandApparent] = useState(3);
  const [landActual, setLandActual] = useState(3);
  const [labApparent, setLabApparent] = useState(80); // wages = variable capital
  const [formError, setFormError] = useState<string | null>(null);
  const [result, setResult] = useState<TrinityFormulaResponse | null>(null);

  const sumApparent = capApparent + landApparent + labApparent;
  const sumSurplus = capActual + landActual; // profit + rent = the single surplus-value s (labour's actual is 0)
  const newValue = sumSurplus + labApparent; // v + s — all the living labour added
  const fetishOvershoot = sumApparent - newValue; // the rent counted twice — the measure of the fetish

  async function refreshLists() {
    try {
      const fs = await financeApi.listRevenueFetishForms();
      setFetishForms(fs);
    } catch (e) {
      setListError(String(e));
    }
  }

  useEffect(() => {
    refreshLists().catch(() => {});
  }, []);

  async function handleCreate() {
    setFormError(null);
    setResult(null);
    try {
      const r = await financeApi.createTrinityFormula({
        capital_stream: {
          apparent_revenue_bp: capApparent,
          actual_source_bp: capActual,
          is_fetishised: true,
        },
        land_stream: {
          apparent_revenue_bp: landApparent,
          actual_source_bp: landActual,
          is_fetishised: true,
        },
        labour_stream: {
          apparent_revenue_bp: labApparent,
          actual_source_bp: 0,
          is_fetishised: true,
        },
      });
      setResult(r);
      await refreshLists();
    } catch (e) {
      setFormError(String(e));
    }
  }

  return (
    <div className="v3-ch48">
      <section className="v3-ch48-intro">
        <p>
          Chapter&nbsp;48 opens Part&nbsp;VII with the <strong>Trinity Formula</strong> — the
          ideological synthesis in which all bourgeois economics culminates:{" "}
          <em>Capital&nbsp;→&nbsp;interest, Land&nbsp;→&nbsp;ground-rent, Labour&nbsp;→&nbsp;wages</em>.
          The formula presents three independent, co-ordinate sources of revenue, each appearing to
          yield its income by its own nature.
        </p>
        <p>
          Marx shows each pairing is a fetish — a &ldquo;prima facie impossible combination.&rdquo;
          Land has no value yet is paired with rent; capital-interest makes value appear to beget
          value (£100 becomes £110); labour-wages is irrational, since labour is the very substance
          of value, not a commodity priced against itself.
        </p>
        <p>
          The real relation: capital pumps surplus-labour from wage-workers. The whole of
          surplus-value (profit&nbsp;+&nbsp;interest&nbsp;+&nbsp;rent) derives from this single
          source and is merely <strong>partitioned</strong> among the three claimants. Landed
          property and wage-labour generate no value of their own — they only claim shares of
          surplus-value already produced. Wages alone are different in kind: they recover the{" "}
          <strong>variable capital</strong> (v), not surplus-value.
        </p>
        <p className="v3-ch48-note">
          Magnitudes in basis points (bp). The seed fixture: average profit&nbsp;20, of which
          ground-rent&nbsp;3 is paid to the landlord, leaving the capitalist&nbsp;17 (profit of
          enterprise&nbsp;12 + interest&nbsp;5). Profit&nbsp;17 + rent&nbsp;3 = the single
          surplus-value&nbsp;s&nbsp;=&nbsp;20; wages&nbsp;80 = variable capital&nbsp;v. Apparent
          revenue&nbsp;103 overshoots newly-created value&nbsp;v&nbsp;+&nbsp;s&nbsp;=&nbsp;100 by
          exactly the rent&nbsp;3 — that overshoot is the fetish (rent appears as income with no
          value behind it), not extra wealth.
        </p>
      </section>

      <section className="v3-ch48-section" aria-label="Trinity formula builder">
        <h4>1. Bind a trinity formula</h4>
        <p className="v3-ch48-subintro">
          Each stream carries an <em>apparent</em> revenue (its fetishised yield) and an{" "}
          <em>actual</em> source (the portion of surplus-value it really represents). The labour
          stream&apos;s actual source is fixed at&nbsp;0 — wages are variable capital, not
          surplus-value. The server computes the totals and binds the three streams.
        </p>
        <div className="v3-ch48-form">
          <label>
            Capital — apparent (bp)
            <input
              type="number"
              min={0}
              value={capApparent}
              onChange={(e) => setCapApparent(Math.max(0, Number(e.target.value)))}
            />
          </label>
          <label>
            Capital — actual s (bp)
            <input
              type="number"
              min={0}
              value={capActual}
              onChange={(e) => setCapActual(Math.max(0, Number(e.target.value)))}
            />
          </label>
          <label>
            Land — apparent (bp)
            <input
              type="number"
              min={0}
              value={landApparent}
              onChange={(e) => setLandApparent(Math.max(0, Number(e.target.value)))}
            />
          </label>
          <label>
            Land — actual s (bp)
            <input
              type="number"
              min={0}
              value={landActual}
              onChange={(e) => setLandActual(Math.max(0, Number(e.target.value)))}
            />
          </label>
          <label>
            Labour — wages / apparent (bp)
            <input
              type="number"
              min={0}
              value={labApparent}
              onChange={(e) => setLabApparent(Math.max(0, Number(e.target.value)))}
            />
            <span className="v3-ch48-hint">actual surplus-value contributed = 0 (wages are v)</span>
          </label>
        </div>
        <p className="v3-ch48-totals">
          apparent revenue&nbsp;= <strong>{sumApparent}</strong>&nbsp;bp&nbsp;·&nbsp; surplus-value
          behind it (profit&nbsp;+&nbsp;rent)&nbsp;= <strong>{sumSurplus}</strong>&nbsp;bp&nbsp;·&nbsp;
          new value (s&nbsp;+&nbsp;v)&nbsp;= <strong>{newValue}</strong>&nbsp;bp&nbsp;·&nbsp; fetish
          overshoot (should&nbsp;=&nbsp;rent)&nbsp;= <strong>{fetishOvershoot}</strong>&nbsp;bp
        </p>
        <button className="v3-ch48-btn" onClick={handleCreate}>
          Bind trinity formula
        </button>
        {formError && <p className="v3-ch48-error">{formError}</p>}
        {result && (
          <>
            <div className="v3-ch48-preview">
              <strong>
                Bound — apparent revenue {result.formula.total_apparent_revenue_bp}&nbsp;bp,
                surplus-value {result.formula.total_surplus_value_bp}&nbsp;bp, fetish overshoot{" "}
                {result.fetish_overshoot_bp}&nbsp;bp
              </strong>
              <span className="v3-ch48-hint">Formula ID: {result.formula.id}</span>
            </div>
            <div className="v3-ch48-table-wrap">
              <table className="v3-ch48-table">
                <thead>
                  <tr>
                    <th>Source</th>
                    <th>Apparent revenue (bp)</th>
                    <th>Actual surplus-value (bp)</th>
                    <th>Fetishised</th>
                  </tr>
                </thead>
                <tbody>
                  {result.streams.map((s) => (
                    <tr key={s.id}>
                      <td>{SOURCE_LABELS[s.source]}</td>
                      <td>{s.apparent_revenue_bp}</td>
                      <td>{s.actual_source_bp}</td>
                      <td>{s.is_fetishised ? "yes" : "no"}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </>
        )}
      </section>

      <section className="v3-ch48-section" aria-label="Revenue fetish forms">
        <h4>2. Revenue fetish forms</h4>
        <p className="v3-ch48-subintro">
          Each surface formula is contrasted with the real relation it conceals, and with the kind
          of mystification at work. The seed records all three branches of the trinity.
        </p>
        {listError && <p className="v3-ch48-error">{listError}</p>}
        {fetishForms.length > 0 && (
          <div className="v3-ch48-table-wrap">
            <table className="v3-ch48-table">
              <thead>
                <tr>
                  <th>Source</th>
                  <th>Surface formula</th>
                  <th>Real relation</th>
                  <th>Mystification</th>
                </tr>
              </thead>
              <tbody>
                {fetishForms.map((f) => (
                  <tr key={f.id}>
                    <td>{SOURCE_LABELS[f.source]}</td>
                    <td>{f.surface_formula}</td>
                    <td>{f.real_relation}</td>
                    <td>{f.mystification_kind}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </section>
    </div>
  );
}
