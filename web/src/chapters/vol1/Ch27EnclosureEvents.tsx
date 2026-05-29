import { useState } from "react";
import type { FormEvent } from "react";
import "./Ch27EnclosureEvents.css";

interface RegistryEntry {
  id: string;
  parish: string;
  acres: number;
  holding: number;
  hh: number;
  common: number;
  displaced: number;
  year: number;
}

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

const TIMELINE = [
  { year: "1485",      name: "Henry VII forbids retainers",              meta: "feudal levies dispersed",                      major: true,  acres: 18000 },
  { year: "1489",      name: "Statute against engrossing",               meta: "first parliamentary check (ineffectual)",       major: false, acres: 12000 },
  { year: "1517",      name: "Wolsey's enclosure commissions",           meta: "264 verdicts returned",                        major: false, acres: 26000 },
  { year: "1536–39",   name: "Dissolution of monasteries",              meta: "monastic land transferred to gentry",           major: true,  acres: 320000 },
  { year: "1549",      name: "Kett's Rebellion",                         meta: "16,000 yeomen in arms · suppressed",           major: false, acres: 14000 },
  { year: "1601",      name: "Elizabethan Poor Law",                     meta: "parish responsibility codified",               major: false, acres: 8000 },
  { year: "1649",      name: "Diggers at St George's Hill",              meta: "common right asserted · dispersed",            major: false, acres: 5000 },
  { year: "1688",      name: "Glorious Revolution — landlords victorious", meta: "Crown lands granted on enormous scale",     major: true,  acres: 240000 },
  { year: "1710–1810", name: "Parliamentary Enclosure Acts",             meta: "≈ 4,000 acts · 6.8 M acres",                  major: true,  acres: 6800000 },
  { year: "1746–1820", name: "Highland Clearances",                      meta: "clans transformed into hunting-grounds",       major: true,  acres: 794000 },
];

const maxAcres = Math.max(...TIMELINE.map((t) => t.acres));

const INITIAL_REGISTRY: RegistryEntry[] = [
  { id: "raunds-1798", parish: "Raunds, Northants.",   acres: 2680, holding: 6, hh: 5, common: 720, displaced: 2233, year: 1798 },
  { id: "otmoor-1815", parish: "Otmoor, Oxfordshire",  acres: 4000, holding: 4, hh: 6, common: 600, displaced: 6000, year: 1815 },
];

export function Ch27EnclosureEvents() {
  const [parish,        setParish]        = useState("Sutton, Notts.");
  const [acresEnclosed, setAcresEnclosed] = useState(1200);
  const [holdingAcres,  setHoldingAcres]  = useState(8);
  const [householdSize, setHouseholdSize] = useState(5);
  const [commonAcres,   setCommonAcres]   = useState(400);
  const [registry, setRegistry] = useState<RegistryEntry[]>(INITIAL_REGISTRY);

  const displaced   = Math.round((acresEnclosed / Math.max(holdingAcres, 0.01)) * householdSize);
  const commonsLost = Math.max(0, Math.round(commonAcres * (acresEnclosed / Math.max(acresEnclosed + commonAcres, 1))));

  function handleRecord(e: FormEvent) {
    e.preventDefault();
    if (!parish.trim()) return;
    const id =
      parish.toLowerCase().replace(/[^a-z0-9]+/g, "-").replace(/^-|-$/g, "") +
      "-" +
      new Date().getFullYear();
    setRegistry((r) => [
      ...r,
      { id, parish, acres: acresEnclosed, holding: holdingAcres, hh: householdSize, common: commonAcres, displaced, year: 1820 },
    ]);
  }

  return (
    <>
      <section className="card">
        <h2>Four Phases</h2>
        <p className="description">
          The expropriation runs through four overlapping mechanisms. The same compulsion produces
          different juridical forms in each century.
        </p>
        <table className="commodity-table">
          <thead>
            <tr>
              <th style={{ width: "5rem" }}>&#167;</th>
              <th>Phase</th>
              <th>Instrument</th>
              <th style={{ textAlign: "right" }}>Century</th>
            </tr>
          </thead>
          <tbody>
            <tr><td className="snlt-cell">I</td>   <td>Disbanding of feudal retinues</td>          <td>Royal proclamation</td>                     <td className="snlt-cell" style={{ textAlign: "right" }}>15th</td></tr>
            <tr><td className="snlt-cell">II</td>  <td>Seizure of Church &amp; Crown estates</td>  <td>Reformation statutes</td>                   <td className="snlt-cell" style={{ textAlign: "right" }}>16th</td></tr>
            <tr><td className="snlt-cell">III</td> <td>Clearing of common rights</td>              <td>Glorious Revolution grants</td>              <td className="snlt-cell" style={{ textAlign: "right" }}>17th</td></tr>
            <tr><td className="snlt-cell">IV</td>  <td>Parliamentary enclosure of commons</td>     <td>Private Acts &amp; Highland Clearances</td>  <td className="snlt-cell" style={{ textAlign: "right" }}>18–19th</td></tr>
          </tbody>
        </table>
      </section>

      <section className="card">
        <h2>A Chronicle of Dispossession</h2>
        <p className="description">
          The chapter&#8217;s footnotes turned into a timeline. The dot is gold for events Marx
          treats as decisive turning-points; the meter at right shows acreage involved, on a
          log-ish scale.
        </p>
        <div className="ch27-timeline">
          {TIMELINE.map((t) => {
            const pct = Math.min(100, (Math.log10(t.acres + 1) / Math.log10(maxAcres + 1)) * 100);
            return (
              <div key={t.year} className={"ch27-event" + (t.major ? " ch27-event--major" : "")}>
                <div className="ch27-event-year">{t.year}</div>
                <div>
                  <div className="ch27-event-name">{t.name}</div>
                  <div className="muted small" style={{ marginTop: "0.2rem" }}>{t.meta}</div>
                </div>
                <div className="ch27-event-meta">
                  <div className="ch27-acreage-bar" title={`${t.acres.toLocaleString()} acres`}>
                    <span style={{ width: pct + "%" }} />
                  </div>
                  <div style={{ marginTop: "0.3rem" }}>{t.acres.toLocaleString()} ac.</div>
                </div>
              </div>
            );
          })}
        </div>
      </section>

      <section className="card">
        <h2>Register an Enclosure</h2>
        <p className="description">
          Compute the displacement implied by a hypothetical enclosure. Displaced = acres &#247;
          avg. holding &#215; household size. The figure is what the parish loses; the new
          sheep-walk or wheat-field is what the landlord gains.
        </p>
        <form className="form-grid" onSubmit={handleRecord}>
          <label className="span2">
            <span>Parish / common</span>
            <input type="text" value={parish} onChange={(e) => setParish(e.target.value)} />
          </label>
          <label>
            <span>Acres enclosed</span>
            <input type="number" min={1} value={acresEnclosed} onChange={(e) => setAcresEnclosed(+e.target.value)} />
          </label>
          <label>
            <span>Average holding (acres)</span>
            <input type="number" min={0.5} step={0.5} value={holdingAcres} onChange={(e) => setHoldingAcres(+e.target.value)} />
          </label>
          <label>
            <span>Household size</span>
            <input type="number" min={1} value={householdSize} onChange={(e) => setHouseholdSize(+e.target.value)} />
          </label>
          <label>
            <span>Residual common (acres)</span>
            <input type="number" min={0} value={commonAcres} onChange={(e) => setCommonAcres(+e.target.value)} />
          </label>
          <label className="span2">
            <span>Computed</span>
            <input
              readOnly
              value={`${displaced.toLocaleString()} persons displaced · ${commonsLost.toLocaleString()} ac. common rights extinguished`}
              style={{ color: "var(--red)" }}
            />
          </label>
          <div className="form-actions span2">
            <button type="submit" className="primary">Record enclosure</button>
          </div>
        </form>

        <div className="ch32-stats" style={{ marginTop: "1.5rem" }}>
          <Stat label="Records on file"     value={registry.length.toLocaleString()} />
          <Stat label="Total acres"         value={registry.reduce((a, r) => a + r.acres, 0).toLocaleString()} tone="lead" />
          <Stat label="Total displaced"     value={registry.reduce((a, r) => a + r.displaced, 0).toLocaleString()} tone="red" />
          <Stat label="Mean household lost" value={householdSize.toFixed(1)} tone="muted" />
        </div>
      </section>

      <section className="card">
        <h2>The Registry</h2>
        <p className="description">
          Each row is one enclosure on file. The total at the bottom is the cumulative number of
          peasants released — into vagrancy in the first instance, into the wage-labour market in
          the second.
        </p>
        <table className="commodity-table">
          <thead>
            <tr>
              <th>Parish</th>
              <th style={{ textAlign: "right" }}>Year</th>
              <th style={{ textAlign: "right" }}>Acres</th>
              <th style={{ textAlign: "right" }}>Avg holding</th>
              <th style={{ textAlign: "right" }}>Displaced</th>
            </tr>
          </thead>
          <tbody>
            {registry.map((r) => (
              <tr key={r.id}>
                <td>{r.parish}</td>
                <td className="snlt-cell" style={{ textAlign: "right" }}>{r.year}</td>
                <td className="snlt-cell" style={{ textAlign: "right" }}>{r.acres.toLocaleString()}</td>
                <td className="snlt-cell" style={{ textAlign: "right" }}>{r.holding} ac.</td>
                <td className="snlt-cell" style={{ textAlign: "right", color: "var(--red)" }}>{r.displaced.toLocaleString()}</td>
              </tr>
            ))}
            <tr>
              <td className="muted small">total</td>
              <td></td>
              <td className="snlt-cell" style={{ textAlign: "right" }}>{registry.reduce((a, r) => a + r.acres, 0).toLocaleString()}</td>
              <td></td>
              <td className="snlt-cell" style={{ textAlign: "right", color: "var(--red)" }}>{registry.reduce((a, r) => a + r.displaced, 0).toLocaleString()}</td>
            </tr>
          </tbody>
        </table>
      </section>

      <aside className="v1-ch27-coda">
        <p className="v1-ch27-coda-quote">
          &#8220;The transformation of arable land into sheep-walks was the watchword of the time in
          England. The robbery of the common lands, the usurpation of feudal and clan property, its
          transformation into modern private property under circumstances of reckless terrorism, were
          just so many idyllic methods of primitive accumulation.&#8221;
          <span className="v1-ch27-coda-cite">&#8212; Marx, Capital Vol. I, Ch. 27</span>
        </p>
      </aside>
    </>
  );
}
