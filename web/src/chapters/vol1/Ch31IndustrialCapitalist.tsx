import { useState } from "react";
import "./Ch31IndustrialCapitalist.css";

const PILLARS: {
  id: string;
  name: string;
  idx: string;
  gloss: string;
  pence: number;
}[] = [
  {
    id: "colony",
    name: "Colonial system",
    idx: "I",
    gloss: "Conquest, plunder, slave trade. The treasures captured outside Europe by undisguised looting, enslavement and murder, floated back to the mother country and were turned there into capital.",
    pence: 230,
  },
  {
    id: "debt",
    name: "Public debt",
    idx: "II",
    gloss: "The State's deficit-financing endows banks with bullion, and the bonds with which the State pays the rentier circulate as money — capital fallen from the sky.",
    pence: 180,
  },
  {
    id: "tax",
    name: "Modern taxation",
    idx: "III",
    gloss: "Taxes on necessaries fall on the labourer; they finance the public debt and provide a continuous compulsion to wage-labour by raising the floor of subsistence.",
    pence: 120,
  },
  {
    id: "tariff",
    name: "Protection",
    idx: "IV",
    gloss: "The protectionist system was an artificial means for manufacturing manufacturers, expropriating independent labourers, and forcibly accelerating capital concentration.",
    pence: 95,
  },
];

export function Ch31IndustrialCapitalist() {
  const [vals, setVals] = useState<number[]>(PILLARS.map((p) => p.pence));
  const total  = vals.reduce((a, b) => a + b, 0);
  const widths = vals.map((v) => (v / Math.max(total, 1)) * 100);

  function setAt(i: number, n: number) {
    setVals((vs) => vs.map((v, j) => (j === i ? Math.max(0, n) : v)));
  }

  return (
    <>
      <section className="card">
        <h2>The Four Pillars</h2>
        <p className="description">
          Marx enumerates four sources of the founding capital of the manufactories. Each is a
          separate institution; together they form the financial counterpart of the expropriation
          in the countryside. The values below are in millions of pounds and roughly trace
          17th&#8211;18th century England.
        </p>
        <div className="ch31-pillars">
          {PILLARS.map((p, i) => (
            <div key={p.id} className="ch31-pillar">
              <div className="ch31-pillar-idx">{p.idx}</div>
              <div className="ch31-pillar-name">{p.name}</div>
              <p className="ch31-pillar-gloss">{p.gloss}</p>
              <label style={{ marginTop: "0.5rem" }}>
                <span style={{ fontSize: "0.55rem", letterSpacing: "0.2em", textTransform: "uppercase", color: "var(--ink-dim)" }}>
                  &#163; millions contributed
                </span>
                <input
                  type="number"
                  min={0}
                  value={vals[i]}
                  onChange={(e) => setAt(i, +e.target.value)}
                />
              </label>
              <div className="ch31-pillar-value">&#163;{vals[i].toLocaleString()} M</div>
            </div>
          ))}
        </div>

        <div style={{ marginTop: "1.5rem" }}>
          <div
            className="muted small"
            style={{ fontFamily: "'IBM Plex Mono', monospace", letterSpacing: "0.2em", textTransform: "uppercase", fontSize: "0.6rem", color: "var(--ink-dim)" }}
          >
            composition · &#163;{total.toLocaleString()} M founded
          </div>
          <div className="ch31-stack">
            <div className="ch31-stack-band ch31-stack-band--colony" style={{ width: widths[0] + "%" }}>
              colony · {widths[0].toFixed(0)}%
            </div>
            <div className="ch31-stack-band ch31-stack-band--debt" style={{ width: widths[1] + "%" }}>
              debt · {widths[1].toFixed(0)}%
            </div>
            <div className="ch31-stack-band ch31-stack-band--tax" style={{ width: widths[2] + "%" }}>
              tax · {widths[2].toFixed(0)}%
            </div>
            <div className="ch31-stack-band ch31-stack-band--tariff" style={{ width: widths[3] + "%" }}>
              tariff · {widths[3].toFixed(0)}%
            </div>
          </div>
        </div>
      </section>

      <section className="card">
        <h2>Colonial Receipts &#8212; East Indies, West Indies, Africa</h2>
        <p className="description">
          The chapter cites figures that, by today&#8217;s accounting, are conservative. The
          right-hand column converts each line into the share of national capital formation it
          represented over the period.
        </p>
        <table className="commodity-table">
          <thead>
            <tr>
              <th>Source</th>
              <th>Period</th>
              <th>Mechanism</th>
              <th style={{ textAlign: "right" }}>&#163; extracted</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td className="snlt-cell">Spanish America</td>
              <td>1503–1660</td>
              <td>Silver remitted via Seville</td>
              <td className="snlt-cell" style={{ textAlign: "right" }}>&#8776; 180 M</td>
            </tr>
            <tr>
              <td className="snlt-cell">East India Company</td>
              <td>1757–1815</td>
              <td>Bengal land revenue + bullion drain</td>
              <td className="snlt-cell" style={{ textAlign: "right" }}>&#8776; 120 M</td>
            </tr>
            <tr>
              <td className="snlt-cell">West Indian plantations</td>
              <td>1660–1807</td>
              <td>Sugar &amp; rum on enslaved labour</td>
              <td className="snlt-cell" style={{ textAlign: "right" }}>&#8776; 90 M</td>
            </tr>
            <tr>
              <td className="snlt-cell">Trans-Atlantic slave trade</td>
              <td>1680–1807</td>
              <td>Liverpool · Bristol · London</td>
              <td className="snlt-cell" style={{ textAlign: "right" }}>&#8776; 60 M</td>
            </tr>
            <tr>
              <td className="snlt-cell">Royal African Company</td>
              <td>1672–1752</td>
              <td>Forts at Cape Coast · Gambia</td>
              <td className="snlt-cell" style={{ textAlign: "right" }}>&#8776; 15 M</td>
            </tr>
          </tbody>
        </table>
      </section>

      <section className="card">
        <h2>Public Debt, in Sketch</h2>
        <p className="description">
          The Bank of England (1694) is the moment the State&#8217;s deficit becomes the
          rentier&#8217;s annuity. The bond is itself a commodity, and its existence multiplies the
          money capital available to industry.
        </p>
        <ul className="ch33-list">
          <li className="ch33-row">
            <div className="ch33-row-main">
              <div className="ch33-row-name">Bank of England · 1694</div>
              <div className="ch33-row-meta muted small">
                &#163;1.2 M lent to the Crown at 8% · banknotes issued against the debt &middot; convertible into specie
              </div>
            </div>
            <span className="ch33-pill ch33-pill--dep">&#163; 1.2 M</span>
          </li>
          <li className="ch33-row">
            <div className="ch33-row-main">
              <div className="ch33-row-name">Funded debt · post-Anne</div>
              <div className="ch33-row-meta muted small">
                South Sea Company · East India bonds · Exchequer bills &middot; the rentier class is born
              </div>
            </div>
            <span className="ch33-pill ch33-pill--dep">&#163; 50 M</span>
          </li>
          <li className="ch33-row">
            <div className="ch33-row-main">
              <div className="ch33-row-name">Napoleonic wars</div>
              <div className="ch33-row-meta muted small">
                Debt rises from &#163;243 M (1793) to &#163;834 M (1815) &middot; financed by indirect taxes on necessaries
              </div>
            </div>
            <span className="ch33-pill ch33-pill--dep">&#163; 834 M</span>
          </li>
          <li className="ch33-row">
            <div className="ch33-row-main">
              <div className="ch33-row-name">Stamp Act &middot; Tea Act</div>
              <div className="ch33-row-meta muted small">
                Colonial taxation extends the metropolitan fiscal apparatus across the Atlantic &middot; revolt in 1776
              </div>
            </div>
            <span className="ch33-pill ch33-pill--free">contested</span>
          </li>
        </ul>
      </section>

      <section className="card">
        <h2>Protection: Tariff as Compulsion</h2>
        <p className="description">
          The protectionist tariff is the mirror image of the colonial system: instead of importing
          the surplus of an enslaved population, it manufactures domestic capitalists by sealing off
          competitors. Marx treats it as a tool of capital concentration, not as a national virtue.
        </p>
        <table className="commodity-table">
          <thead>
            <tr>
              <th>Measure</th>
              <th>Year</th>
              <th>Effect on the labourer</th>
              <th style={{ textAlign: "right" }}>Effect on capital</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td className="snlt-cell">Navigation Acts</td>
              <td>1651&#183;1660</td>
              <td>Higher freight passed to consumer</td>
              <td className="snlt-cell" style={{ textAlign: "right", color: "var(--gold-bright)" }}>Merchant marine</td>
            </tr>
            <tr>
              <td className="snlt-cell">Calico Acts</td>
              <td>1700&#183;1721</td>
              <td>Indian cloth banned · domestic price rises</td>
              <td className="snlt-cell" style={{ textAlign: "right", color: "var(--gold-bright)" }}>Manchester emerges</td>
            </tr>
            <tr>
              <td className="snlt-cell">Corn Laws</td>
              <td>1815</td>
              <td>Bread price floored above world market</td>
              <td className="snlt-cell" style={{ textAlign: "right", color: "var(--gold-bright)" }}>Landed capital</td>
            </tr>
            <tr>
              <td className="snlt-cell">Combination Acts</td>
              <td>1799&#183;1800</td>
              <td>Wage-bargaining criminalised</td>
              <td className="snlt-cell" style={{ textAlign: "right", color: "var(--gold-bright)" }}>Labour discipline</td>
            </tr>
          </tbody>
        </table>
      </section>

      <div className="ch31-coda">
        <p>
          &#8220;If money, according to Augier, &#8216;comes into the world with a congenital
          blood-stain on one cheek,&#8217;{" "}
          <span className="blood">
            capital comes dripping from head to foot, from every pore, with blood and dirt.
          </span>&#8221;
        </p>
      </div>
    </>
  );
}
