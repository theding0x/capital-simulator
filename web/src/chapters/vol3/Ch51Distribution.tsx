import { useEffect, useState } from "react";
import { financeApi } from "../../api";
import type {
  DistributionRelation,
  HistoricalTransience,
  ProductionRelationBasis,
  ProductiveForceContradiction,
} from "../../types";
import "./Ch51Distribution.css";

export function Ch51Distribution() {
  const [bases, setBases] = useState<ProductionRelationBasis[]>([]);
  const [relations, setRelations] = useState<DistributionRelation[]>([]);
  const [transiences, setTransiences] = useState<HistoricalTransience[]>([]);
  const [contradictions, setContradictions] = useState<ProductiveForceContradiction[]>([]);
  const [listError, setListError] = useState<string | null>(null);

  const [force, setForce] = useState("socialised large-scale industry, credit, joint-stock companies");
  const [social, setSocial] = useState("private ownership of the means of production");
  const [nature, setNature] = useState("socialised production constrained by private appropriation");
  const [signal, setSignal] = useState("commercial and industrial crises");
  const [formError, setFormError] = useState<string | null>(null);
  const [result, setResult] = useState<ProductiveForceContradiction | null>(null);

  async function refreshLists() {
    try {
      const [b, r, t, c] = await Promise.all([
        financeApi.listProductionRelationBases(),
        financeApi.listDistributionRelations(),
        financeApi.listHistoricalTransiences(),
        financeApi.listProductiveForceContradictions(),
      ]);
      setBases(b);
      setRelations(r);
      setTransiences(t);
      setContradictions(c);
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
      const r = await financeApi.createProductiveForceContradiction({
        productive_force: force,
        social_form: social,
        contradiction_nature: nature,
        crisis_signal: signal,
      });
      setResult(r);
      await refreshLists();
    } catch (e) {
      setFormError(String(e));
    }
  }

  return (
    <div className="v3-ch51">
      <section className="v3-ch51-intro">
        <p>
          Chapter&nbsp;51 establishes that <strong>distribution relations</strong> (wages, profit,
          rent) are not a separate sphere but the &ldquo;other side&rdquo; of{" "}
          <strong>production relations</strong>, sharing their historically transitory character. The
          bourgeois error treats distribution as adjustable while treating production relations
          (private property, wage-labour) as natural and eternal.
        </p>
        <p>
          Each distribution form rests on a specific production relation and corresponds to a
          component of the new value: wages → variable capital (v); profit and rent → portions of
          surplus-value (s). A <strong>crisis</strong> arises when the development of the productive
          forces contradicts the existing social form — the historical limit of the mode of
          production.
        </p>
        {listError && <p className="v3-ch51-error">{listError}</p>}
      </section>

      <section className="v3-ch51-section" aria-label="Production relation bases">
        <h4>1. Production-relation bases</h4>
        {bases.length > 0 && (
          <div className="v3-ch51-table-wrap">
            <table className="v3-ch51-table">
              <thead>
                <tr>
                  <th>Name</th>
                  <th>Characteristic feature</th>
                  <th>Historically specific</th>
                </tr>
              </thead>
              <tbody>
                {bases.map((b) => (
                  <tr key={b.id}>
                    <td>{b.name}</td>
                    <td>{b.characteristic_feature}</td>
                    <td>{b.is_historically_specific ? "yes" : "no"}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </section>

      <section className="v3-ch51-section" aria-label="Distribution relations">
        <h4>2. Distribution relations</h4>
        {relations.length > 0 && (
          <div className="v3-ch51-table-wrap">
            <table className="v3-ch51-table">
              <thead>
                <tr>
                  <th>Form</th>
                  <th>Apparent character</th>
                  <th>Real character</th>
                  <th>Surplus form</th>
                </tr>
              </thead>
              <tbody>
                {relations.map((r) => (
                  <tr key={r.id}>
                    <td>{r.form}</td>
                    <td>{r.apparent_character}</td>
                    <td>{r.real_character}</td>
                    <td>{r.corresponds_to_surplus_form}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </section>

      <section className="v3-ch51-section" aria-label="Historical transience">
        <h4>3. Historical transience</h4>
        {transiences.length > 0 && (
          <div className="v3-ch51-table-wrap">
            <table className="v3-ch51-table">
              <thead>
                <tr>
                  <th>Mode of production</th>
                  <th>Precondition</th>
                  <th>Succeeding form</th>
                </tr>
              </thead>
              <tbody>
                {transiences.map((t) => (
                  <tr key={t.id}>
                    <td>{t.mode_of_production}</td>
                    <td>{t.precondition}</td>
                    <td>{t.succeeding_form}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </section>

      <section className="v3-ch51-section" aria-label="Productive-force contradiction">
        <h4>4. Productive-force contradiction</h4>
        <p className="v3-ch51-subintro">
          Record a tension between the developing productive forces and the existing social form —
          the contradiction that signals the historical limit of the mode of production.
        </p>
        <div className="v3-ch51-form">
          <label>
            Productive force
            <input value={force} onChange={(e) => setForce(e.target.value)} />
          </label>
          <label>
            Social form
            <input value={social} onChange={(e) => setSocial(e.target.value)} />
          </label>
          <label>
            Contradiction nature
            <input value={nature} onChange={(e) => setNature(e.target.value)} />
          </label>
          <label>
            Crisis signal
            <input value={signal} onChange={(e) => setSignal(e.target.value)} />
          </label>
        </div>
        <button className="v3-ch51-btn" onClick={handleCreate}>
          Record contradiction
        </button>
        {formError && <p className="v3-ch51-error">{formError}</p>}
        {result && (
          <div className="v3-ch51-preview">
            <strong>Saved — {result.productive_force} vs {result.social_form}</strong>
            <span className="v3-ch51-hint">crisis signal: {result.crisis_signal}</span>
          </div>
        )}
        {contradictions.length > 0 && (
          <div className="v3-ch51-table-wrap">
            <table className="v3-ch51-table">
              <thead>
                <tr>
                  <th>Productive force</th>
                  <th>Social form</th>
                  <th>Contradiction</th>
                  <th>Crisis signal</th>
                </tr>
              </thead>
              <tbody>
                {contradictions.map((c) => (
                  <tr key={c.id}>
                    <td>{c.productive_force}</td>
                    <td>{c.social_form}</td>
                    <td>{c.contradiction_nature}</td>
                    <td>{c.crisis_signal}</td>
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
