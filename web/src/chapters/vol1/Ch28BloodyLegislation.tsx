import { useState } from "react";
import "./Ch28BloodyLegislation.css";

type Offence = "idle" | "sturdy" | "rep2" | "rep3";

const STATUTES: {
  cite: string;
  year: string;
  name: string;
  clause: string;
  marks: boolean[];
  punish: string;
}[] = [
  {
    cite: "27 Henry VIII, c. 25",
    year: "1535",
    name: "Act for Punishment of Sturdy Vagabonds and Beggars",
    clause:
      "Idle persons taken begging shall be whipped at the cart's tail till their bodies be bloody, and on the second offence have the upper part of the right ear cut off.",
    marks: [false, true, false, true, false],
    punish: "Whipping · ear-cropping",
  },
  {
    cite: "1 Edward VI, c. 3",
    year: "1547",
    name: "Statute of Slaves",
    clause:
      "A vagabond refusing to work shall be adjudged the slave of any person informing against him, for two years; if he run away, branded with S, slave for life; if again, hanged as a felon.",
    marks: [false, false, true, true, true],
    punish: "Slavery · branding · gallows",
  },
  {
    cite: "14 Elizabeth, c. 5",
    year: "1572",
    name: "Act for the Punishment of Vagabonds",
    clause:
      "Vagabonds above 14 to be grievously whipped and burnt through the gristle of the right ear with a hot iron, unless some person take them into service for one year.",
    marks: [false, true, true, true, false],
    punish: "Whipping · branding · ear-bored",
  },
  {
    cite: "39 Elizabeth, c. 4",
    year: "1597",
    name: "Act for the Punishment of Rogues, Vagabonds and Sturdy Beggars",
    clause:
      "Incorrigible and dangerous rogues to be banished from the realm, or condemned perpetually to the galleys.",
    marks: [false, true, false, true, false],
    punish: "Banishment · galleys",
  },
  {
    cite: "1 James I, c. 7",
    year: "1604",
    name: "Act for Continuing the Statutes of Vagabonds",
    clause:
      "Rogues to be branded with letter R on the left shoulder; if found at large after, executed as a felon without benefit of clergy.",
    marks: [false, false, true, false, true],
    punish: "Branding · gallows",
  },
  {
    cite: "13 &amp; 14 Charles II, c. 12",
    year: "1662",
    name: "Act of Settlement",
    clause:
      "Any person likely to become chargeable to the parish may within forty days be returned to the parish of his last legal settlement.",
    marks: [true, false, false, false, false],
    punish: "Removal · fine",
  },
  {
    cite: "31 George II, c. 5",
    year: "1757",
    name: "Vagrant Act",
    clause:
      "Vagrants to be publicly whipped, then committed to the house of correction with hard labour not exceeding six months.",
    marks: [false, true, false, true, false],
    punish: "Whipping · hard labour",
  },
];

const MARK_LABELS = ["fine", "whip", "brand", "bond·galley·ear", "gallows"];

const REIGNS: {
  id: string;
  label: string;
  punish: Record<Offence, string>;
}[] = [
  {
    id: "henry",
    label: "Henry VIII (1535)",
    punish: { idle: "fine", sturdy: "whip · ear-bored", rep2: "whip · ear-cropped", rep3: "hanged" },
  },
  {
    id: "edward",
    label: "Edward VI (1547)",
    punish: { idle: "whip", sturdy: "2-yr slave", rep2: "brand S · life slave", rep3: "hanged" },
  },
  {
    id: "elizabeth",
    label: "Elizabeth (1572)",
    punish: { idle: "whip", sturdy: "ear-bored · whip", rep2: "felon · 2-yr service", rep3: "hanged" },
  },
  {
    id: "james",
    label: "James I (1604)",
    punish: { idle: "whip · returned", sturdy: "brand R", rep2: "felon", rep3: "hanged" },
  },
];

const OFFENCE_PILLS: { id: Offence; tag: string; label: string }[] = [
  { id: "idle",   tag: "I",   label: "Idle · first offence" },
  { id: "sturdy", tag: "II",  label: "Sturdy · refusing work" },
  { id: "rep2",   tag: "III", label: "Second offence" },
  { id: "rep3",   tag: "IV",  label: "Third offence" },
];

export function Ch28BloodyLegislation() {
  const [offence, setOffence] = useState<Offence>("sturdy");

  return (
    <>
      <section className="card">
        <h2>The Problem of the Expropriated</h2>
        <p className="description">
          Driven from the soil but not yet absorbed by the still-young manufactories, the
          dispossessed appeared as beggars, robbers and vagabonds. Two and a half centuries of
          legislation — under five reigns — converted them into wage-labourers by force of the
          statute book.
        </p>
        <p className="note" style={{ maxWidth: "62ch", marginTop: "0.5rem" }}>
          The chapter reads as a litany of statutes. Each one is a tool. Together they form the
          juridical infrastructure of the wage-relation: a class compelled by the gallows to sell
          its labour-power at the going price.
        </p>
      </section>

      <section className="card">
        <h2>Statute Registry</h2>
        <p className="description">
          One row per act. The severity meter marks which categories of punishment the statute
          prescribes; the right-hand column quotes the prescribed sentence in the chapter&#8217;s
          language.
        </p>
        <ul className="ch28-statutes">
          {STATUTES.map((s) => (
            <li key={s.cite} className="ch28-statute">
              <div className="ch28-statute-cite">
                <div dangerouslySetInnerHTML={{ __html: s.cite }} />
                <div>{s.year}</div>
              </div>
              <div>
                <div className="ch28-statute-name">{s.name}</div>
                <p className="ch28-statute-clause">&#8220;{s.clause}&#8221;</p>
              </div>
              <div className="ch28-severity">
                <div className="ch28-severity-label">severity</div>
                <div className="ch28-severity-marks" title={MARK_LABELS.join(" · ")}>
                  {s.marks.map((on, i) => (
                    <span key={i} className={on ? (i === 4 ? "on gold" : "on") : ""} />
                  ))}
                </div>
                <div className="ch28-severity-punish">{s.punish}</div>
              </div>
            </li>
          ))}
        </ul>
      </section>

      <section className="card">
        <h2>Classify a Vagabond</h2>
        <p className="description">
          The same person, classified by successive statutes. Pick the offence; read across the
          reigns. The trajectory is not a softening — it is the consolidation of a labour
          discipline that will not, by 1800, require the gallows to be enforced.
        </p>
        <div className="ch17-presets" style={{ marginBottom: "1rem" }}>
          {OFFENCE_PILLS.map((p) => (
            <button
              key={p.id}
              type="button"
              className={"ch17-preset" + (offence === p.id ? " ch17-preset--active" : "")}
              onClick={() => setOffence(p.id)}
            >
              <span className="ch17-preset-tag">{p.tag}</span>
              <span className="ch17-preset-label">{p.label}</span>
            </button>
          ))}
        </div>

        <table className="commodity-table" style={{ marginTop: "0.5rem" }}>
          <thead>
            <tr>
              <th>Reign</th>
              <th>Citation</th>
              <th style={{ textAlign: "right" }}>Punishment</th>
            </tr>
          </thead>
          <tbody>
            {REIGNS.map((r, i) => {
              const p = r.punish[offence];
              const death = /hang|gallows|felon/.test(p);
              return (
                <tr key={r.id}>
                  <td>{r.label}</td>
                  <td className="snlt-cell">{STATUTES[i]?.cite.replace(/&amp;/g, "&") ?? "—"}</td>
                  <td
                    className="snlt-cell"
                    style={{ textAlign: "right", color: death ? "var(--red)" : "var(--ink)" }}
                  >
                    {p}
                  </td>
                </tr>
              );
            })}
          </tbody>
        </table>

        <p className="note" style={{ marginTop: "1.25rem", maxWidth: "62ch" }}>
          &#8220;Thus were the agricultural people, first forcibly expropriated from the soil,
          driven from their homes, turned into vagabonds, and then whipped, branded, tortured by
          laws grotesquely terrible, into the discipline necessary for the wage system.&#8221;
        </p>
      </section>

      <section className="card">
        <h2>Maximum of Wages</h2>
        <p className="description">
          The same statute book that hanged the vagabond also fixed a legal ceiling on the wage.
          The combination is what makes the period coherent: violence at the bottom, legislation
          at the top.
        </p>
        <table className="commodity-table">
          <thead>
            <tr>
              <th>Statute</th>
              <th>Reign</th>
              <th>Provision</th>
              <th style={{ textAlign: "right" }}>Penalty for excess</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td className="snlt-cell">23 Edw. III</td>
              <td>1349 &#8212; Ordinance of Labourers</td>
              <td>Compels work at pre-plague wages; refusal punished by imprisonment.</td>
              <td className="snlt-cell" style={{ textAlign: "right", color: "var(--red)" }}>imprisonment</td>
            </tr>
            <tr>
              <td className="snlt-cell">5 Eliz. c. 4</td>
              <td>1563 &#8212; Statute of Artificers</td>
              <td>Justices to assess maximum wages annually for every trade and county.</td>
              <td className="snlt-cell" style={{ textAlign: "right", color: "var(--red)" }}>10 days&#8217; imprisonment</td>
            </tr>
            <tr>
              <td className="snlt-cell">7 Geo. I c. 13</td>
              <td>1720 &#8212; Tailors</td>
              <td>Combinations to raise wages punished by imprisonment with hard labour.</td>
              <td className="snlt-cell" style={{ textAlign: "right", color: "var(--red)" }}>2 mo. hard labour</td>
            </tr>
            <tr>
              <td className="snlt-cell">39 Geo. III c. 81</td>
              <td>1799 &#8212; Combination Acts</td>
              <td>All workmen&#8217;s combinations made criminal; repealed only in 1824.</td>
              <td className="snlt-cell" style={{ textAlign: "right", color: "var(--red)" }}>3 mo. gaol</td>
            </tr>
          </tbody>
        </table>
      </section>
    </>
  );
}
