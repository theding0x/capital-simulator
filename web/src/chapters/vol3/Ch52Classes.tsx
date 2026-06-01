import { useEffect, useState } from "react";
import { financeApi } from "../../api";
import type {
  ClassAmbiguity,
  SocialClass,
  SocialClassDetail,
} from "../../types";
import "./Ch52Classes.css";

export function Ch52Classes() {
  const [classes, setClasses] = useState<SocialClass[]>([]);
  const [ambiguities, setAmbiguities] = useState<ClassAmbiguity[]>([]);
  const [detail, setDetail] = useState<SocialClassDetail | null>(null);
  const [listError, setListError] = useState<string | null>(null);

  async function refreshLists() {
    try {
      const [cls, amb] = await Promise.all([
        financeApi.listSocialClasses(),
        financeApi.listClassAmbiguities(),
      ]);
      setClasses(cls);
      setAmbiguities(amb);
    } catch (e) {
      setListError(String(e));
    }
  }

  useEffect(() => {
    refreshLists().catch(() => {});
  }, []);

  async function handleSelect(id: string) {
    try {
      setDetail(await financeApi.getSocialClass(id));
    } catch (e) {
      setListError(String(e));
    }
  }

  return (
    <div className="v3-ch52">
      <section className="v3-ch52-intro">
        <p>
          Chapter&nbsp;52 is the final, <strong>unfinished</strong> chapter of <em>Capital</em> — the
          manuscript breaks off after two pages. Marx identifies the three great classes of capitalist
          society: <strong>wage-labourers</strong> (wages), <strong>capitalists</strong> (profit),
          and <strong>landowners</strong> (ground-rent).
        </p>
        <p>
          He then poses the defining question — <em>what constitutes a class?</em> — and notes that a
          shared source of income is not enough (physicians and officials have defined incomes yet
          form no class). The answer he never wrote: a class is defined by its{" "}
          <strong>relation to the means of production</strong>. The officials ambiguity below is
          seeded with <code>manuscript_breaks_off</code> to mark that incompleteness faithfully.
        </p>
        {listError && <p className="v3-ch52-error">{listError}</p>}
      </section>

      <section className="v3-ch52-section" aria-label="The three great classes">
        <h4>1. The three great classes</h4>
        {classes.length > 0 && (
          <div className="v3-ch52-table-wrap">
            <table className="v3-ch52-table">
              <thead>
                <tr>
                  <th>Class</th>
                  <th>Relation to means of production</th>
                  <th>Revenue form</th>
                  <th>Historical precondition</th>
                  <th></th>
                </tr>
              </thead>
              <tbody>
                {classes.map((c) => (
                  <tr key={c.id}>
                    <td>{c.name}</td>
                    <td>{c.relation_to_means_of_production}</td>
                    <td>{c.primary_revenue_form}</td>
                    <td>{c.historical_precondition}</td>
                    <td>
                      <button className="v3-ch52-link" onClick={() => handleSelect(c.id)}>
                        view detail
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </section>

      {detail && (
        <section className="v3-ch52-section" aria-label="Class detail">
          <h4>2. {detail.class.name} — income sources &amp; tendencies</h4>
          {detail.income_sources.length > 0 && (
            <div className="v3-ch52-table-wrap">
              <table className="v3-ch52-table">
                <thead>
                  <tr>
                    <th>Revenue form</th>
                    <th>Surplus-value component (bp)</th>
                  </tr>
                </thead>
                <tbody>
                  {detail.income_sources.map((s) => (
                    <tr key={s.id}>
                      <td>{s.revenue_form}</td>
                      <td>{s.surplus_value_component_bp} {s.surplus_value_component_bp === 0 ? "(= v, not surplus-value)" : "(portion of s)"}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
          {detail.tendencies.length > 0 && (
            <div className="v3-ch52-table-wrap">
              <table className="v3-ch52-table">
                <thead>
                  <tr>
                    <th>Tendency</th>
                    <th>Direction</th>
                  </tr>
                </thead>
                <tbody>
                  {detail.tendencies.map((t) => (
                    <tr key={t.id}>
                      <td>{t.description}</td>
                      <td>{t.direction}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </section>
      )}

      <section className="v3-ch52-section" aria-label="Class ambiguities">
        <h4>3. What constitutes a class? (the unanswered question)</h4>
        <p className="v3-ch52-subintro">
          Intermediate strata blur the pure class lines. Marx&apos;s answer to why they form no class
          of their own was never written — the manuscript ends here, and with it Volume&nbsp;III.
        </p>
        {ambiguities.length > 0 && (
          <div className="v3-ch52-table-wrap">
            <table className="v3-ch52-table">
              <thead>
                <tr>
                  <th>Intermediate group</th>
                  <th>Income source</th>
                  <th>Why not a class?</th>
                </tr>
              </thead>
              <tbody>
                {ambiguities.map((a) => (
                  <tr key={a.id}>
                    <td>{a.intermediate_group}</td>
                    <td>{a.income_source}</td>
                    <td><em>{a.why_not_a_class}</em></td>
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
