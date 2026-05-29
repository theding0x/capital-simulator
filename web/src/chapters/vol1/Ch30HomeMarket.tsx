import { useState } from "react";
import "./Ch30HomeMarket.css";

function Stat({
  label,
  value,
  tone,
}: {
  label: string;
  value: string;
  tone?: "red" | "gold" | "lead" | "muted";
}) {
  const color =
    tone === "red"   ? "var(--red)" :
    tone === "gold"  ? "var(--gold-bright)" :
    tone === "lead"  ? "var(--lead-hover)" :
    tone === "muted" ? "var(--ink-muted)" :
    "var(--ink)";
  return (
    <div className="ch32-stat">
      <div className="ch32-stat-label">{label}</div>
      <div className="ch32-stat-value mono" style={{ color }}>{value}</div>
    </div>
  );
}

export function Ch30HomeMarket() {
  const [released,      setReleased] = useState(50000);
  const [absorb,        setAbsorb]   = useState(0.62);
  const [annualSubsist, setSubsist]  = useState(720);
  const [factorySpend,  setSpend]    = useState(540);
  const [householdSize, setHH]       = useState(4.5);

  const households   = Math.round(released / Math.max(householdSize, 1));
  const proletarians = Math.round(released * absorb);
  const lostSubsist  = households * annualSubsist;
  const newDemand    = households * factorySpend;
  const sourceShare  = newDemand / Math.max(lostSubsist, 1);

  return (
    <>
      <section className="card">
        <h2>One Process, Two Inputs</h2>
        <p className="description">
          The expropriation that supplied the towns with proletarians simultaneously created the
          market into which those proletarians&#8217; product could be sold. The peasant who used
          to spin her own yarn now buys factory cloth out of a wage paid by the spinning-mill.
        </p>
        <form className="form-grid" onSubmit={(e) => e.preventDefault()}>
          <label>
            <span>Persons released from the soil</span>
            <input type="number" min={1} value={released} onChange={(e) => setReleased(+e.target.value)} />
          </label>
          <label>
            <span>Household size</span>
            <input type="number" min={1} step={0.1} value={householdSize} onChange={(e) => setHH(+e.target.value)} />
          </label>
          <label>
            <span>Share absorbed by industry</span>
            <input type="number" min={0} max={1} step={0.05} value={absorb} onChange={(e) => setAbsorb(+e.target.value)} />
          </label>
          <label>
            <span>Annual subsistence value forgone (&#163; / household)</span>
            <input type="number" min={0} value={annualSubsist} onChange={(e) => setSubsist(+e.target.value)} />
          </label>
          <label className="span2">
            <span>Annual industrial spend per household (&#163;)</span>
            <input type="number" min={0} value={factorySpend} onChange={(e) => setSpend(+e.target.value)} />
          </label>
        </form>
      </section>

      <section className="card">
        <h2>The Fork</h2>
        <p className="description">
          The single rural pool feeds both halves of the industrial mode at once. Each arrow carries
          a quantity computed from the inputs above.
        </p>
        <div className="ch30-diagram">
          <div className="ch30-pool ch30-pool--source">
            <div className="ch30-pool-tag">released · rural surplus</div>
            <div className="ch30-pool-value mono">{released.toLocaleString()}</div>
            <p className="ch30-pool-sub">
              {households.toLocaleString()} households · expropriated by enclosure, Reformation
              seizure or parliamentary act
            </p>
          </div>
          <div className="ch30-fork">
            <span>supply</span>
            <span className="ch30-fork-arrow">&#x27F6;</span>
            <span>demand</span>
          </div>
          <div className="ch30-destinations">
            <div className="ch30-pool">
              <div className="ch30-pool-tag">industrial proletariat</div>
              <div className="ch30-pool-value mono">{proletarians.toLocaleString()}</div>
              <p className="ch30-pool-sub">
                absorbed into the manufactories; selling labour-power at the going wage
              </p>
            </div>
            <div className="ch30-pool ch30-pool--demand">
              <div className="ch30-pool-tag">home market · annual demand</div>
              <div className="ch30-pool-value mono">&#163;{newDemand.toLocaleString()}</div>
              <p className="ch30-pool-sub">
                cloth, boots, candles, beer formerly produced at home, now purchased from manufacture
              </p>
            </div>
          </div>
        </div>

        <div className="ch30-stats">
          <Stat label="Households expropriated" value={households.toLocaleString()}              tone="red" />
          <Stat label="Proletarianised"          value={proletarians.toLocaleString()}           tone="lead" />
          <Stat label="Subsistence extinguished" value={`£${lostSubsist.toLocaleString()}/yr`}  tone="muted" />
          <Stat label="Demand transferred"       value={`£${newDemand.toLocaleString()}/yr`}    tone="gold" />
        </div>

        <p className="note" style={{ marginTop: "1.25rem", maxWidth: "62ch" }}>
          &#8220;The expropriation and expulsion of the agricultural population, intermittent but
          renewed again and again, supplied the town industries with a mass of proletarians entirely
          unconnected with the corporate guilds and unfettered by them; a fortunate circumstance
          that makes old A. Anderson believe in direct intervention of providence.&#8221;
        </p>
      </section>

      <section className="card">
        <h2>Self-Provision &#x27F6; Commodity</h2>
        <p className="description">
          The labourer&#8217;s diet, cloth and tools do not vanish: they change hands. What was
          previously a use-value produced by the household becomes a commodity, bought back from
          the master who sets the household to work.
        </p>
        <table className="commodity-table">
          <thead>
            <tr>
              <th>Article</th>
              <th>Pre-expropriation source</th>
              <th>Post-expropriation source</th>
              <th style={{ textAlign: "right" }}>Annual value (&#163;)</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td className="snlt-cell">Yarn &amp; cloth</td>
              <td>Household spinning · loom</td>
              <td>Manchester spinning-mill</td>
              <td className="snlt-cell" style={{ textAlign: "right" }}>&#163;{Math.round(annualSubsist * 0.32).toLocaleString()}</td>
            </tr>
            <tr>
              <td className="snlt-cell">Beer</td>
              <td>Domestic brewing</td>
              <td>Commercial brewery</td>
              <td className="snlt-cell" style={{ textAlign: "right" }}>&#163;{Math.round(annualSubsist * 0.18).toLocaleString()}</td>
            </tr>
            <tr>
              <td className="snlt-cell">Boots</td>
              <td>Village cordwainer</td>
              <td>Northampton factory</td>
              <td className="snlt-cell" style={{ textAlign: "right" }}>&#163;{Math.round(annualSubsist * 0.12).toLocaleString()}</td>
            </tr>
            <tr>
              <td className="snlt-cell">Candles · soap</td>
              <td>Tallow rendered at home</td>
              <td>Provincial chandler</td>
              <td className="snlt-cell" style={{ textAlign: "right" }}>&#163;{Math.round(annualSubsist * 0.10).toLocaleString()}</td>
            </tr>
            <tr>
              <td className="snlt-cell">Implements</td>
              <td>Smith on commission</td>
              <td>Wholesale ironmonger</td>
              <td className="snlt-cell" style={{ textAlign: "right" }}>&#163;{Math.round(annualSubsist * 0.15).toLocaleString()}</td>
            </tr>
            <tr>
              <td className="snlt-cell">Misc.</td>
              <td>Barter · common right</td>
              <td>Cash market</td>
              <td className="snlt-cell" style={{ textAlign: "right" }}>&#163;{Math.round(annualSubsist * 0.13).toLocaleString()}</td>
            </tr>
          </tbody>
        </table>

        <p className="note" style={{ marginTop: "1rem", maxWidth: "62ch" }}>
          Each row is the same use-value before and after. The transfer is what makes the
          manufactory&#8217;s product saleable in the first place: a market is created by the very
          dispossession that supplies its workforce. The proportion of demand recovered by industry
          at our default inputs is&#160;
          <span className="mono" style={{ color: "var(--gold-bright)" }}>
            {(sourceShare * 100).toFixed(0)}%
          </span>
          &#160;of the value extinguished &#8212; the residual is monetised savings, taxation, and
          a permanent reduction in the standard of living.
        </p>
      </section>

      <section className="card">
        <h2>What Lloyd &amp; Anderson Saw</h2>
        <p className="description">
          The contemporary witnesses Marx quotes were not, in the main, hostile to the process. They
          observed the same fork and found it providential.
        </p>
        <ul className="ch33-list">
          <li className="ch33-row">
            <div className="ch33-row-main">
              <div className="ch33-row-name">A. Anderson · 1764</div>
              <div className="ch33-row-meta muted small">
                Origin of Commerce &#8212; &#8220;direct interposition of providence&#8221;
              </div>
            </div>
            <span className="ch33-pill ch33-pill--free">apologetic</span>
          </li>
          <li className="ch33-row">
            <div className="ch33-row-main">
              <div className="ch33-row-name">Sir F. M. Eden · 1797</div>
              <div className="ch33-row-meta muted small">
                State of the Poor &#8212; &#8220;increased manufactures must absorb the rural surplus&#8221;
              </div>
            </div>
            <span className="ch33-pill ch33-pill--free">apologetic</span>
          </li>
          <li className="ch33-row">
            <div className="ch33-row-main">
              <div className="ch33-row-name">Sir T. More · 1516</div>
              <div className="ch33-row-meta muted small">
                Utopia &#8212; &#8220;sheep devour men&#8221;
              </div>
            </div>
            <span className="ch33-pill ch33-pill--dep">opposed</span>
          </li>
          <li className="ch33-row">
            <div className="ch33-row-main">
              <div className="ch33-row-name">Dr. Price · 1780</div>
              <div className="ch33-row-meta muted small">
                Observations &#8212; &#8220;the kingdom is wasting away&#8221;
              </div>
            </div>
            <span className="ch33-pill ch33-pill--dep">opposed</span>
          </li>
        </ul>
      </section>

      <aside className="v1-ch30-coda">
        <p className="v1-ch30-coda-quote">
          &#8220;The expropriation and expulsion of a part of the agricultural population not only
          set free for industrial capital the labourers, their means of subsistence, and material
          for labour; it also created the home market.&#8221;
          <span className="v1-ch30-coda-cite">&#8212; Marx, Capital Vol. I, Ch. 30</span>
        </p>
      </aside>
    </>
  );
}
