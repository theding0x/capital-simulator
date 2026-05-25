import { useState } from "react";
import "./Ch29GenesisFarmer.css";

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

const STEPS: {
  era: string;
  figure: string;
  desc: string;
  farmer: number;
  labour: number;
  final?: boolean;
}[] = [
  {
    era: "≈ 1300",
    figure: "Bailiff",
    desc: "Holder of office under the lord. Paid in kind from the demesne; tied to the manor; no surplus of his own to dispose of.",
    farmer: 0.08,
    labour: 0.95,
  },
  {
    era: "1400s",
    figure: "Métayer · half-share tenant",
    desc: "Lord supplies seed, cattle, implements; tenant supplies labour and management; product divided. The tenant begins to retain a residual.",
    farmer: 0.20,
    labour: 0.78,
  },
  {
    era: "1450–1550",
    figure: "Yeoman · freeholder",
    desc: "Owner-occupier of a small estate. Holds independent of the lord, but exposed to enclosure, taxation and the price revolution.",
    farmer: 0.38,
    labour: 0.58,
  },
  {
    era: "1550–1650",
    figure: "Lease-holder of the Tudor type",
    desc: "Long lease at a customary money rent. Through the price revolution the real value of the rent collapses; the tenant pockets the difference. The agricultural revolution is in part a windfall.",
    farmer: 0.62,
    labour: 0.34,
  },
  {
    era: "1650–1750",
    figure: "Capitalist farmer",
    desc: "Employs wage-labour on a multi-hundred-acre lease. Pays a rack-rent, but the rate of surplus-value extracted from the labourers leaves a profit characteristic of the new mode of production.",
    farmer: 0.90,
    labour: 0.12,
    final: true,
  },
];

export function Ch29GenesisFarmer() {
  const [wage1500,  setWage1500]  = useState(4);
  const [wage1750,  setWage1750]  = useState(8);
  const [price1500, setPrice1500] = useState(100);
  const [price1750, setPrice1750] = useState(420);
  const [rentReal,  setRentReal]  = useState(0.35);

  const wageRise   = wage1750 / Math.max(wage1500, 1);
  const priceRise  = price1750 / Math.max(price1500, 1);
  const realWage   = wageRise / priceRise;
  const fatten     = (1 - realWage) * 100;
  const farmerGain = (priceRise - 1) * rentReal * 100;

  return (
    <>
      <section className="card">
        <h2>The Twin Movement</h2>
        <p className="description">
          A class is not born ex nihilo. The capitalist farmer is the bailiff transformed by three
          centuries of intervening events: the price revolution, the enclosures, the dispossession
          of the labourer from the soil. Each enriched him, and the same events impoverished the
          people he came to employ.
        </p>
        <ol className="ch29-genealogy">
          {STEPS.map((s) => (
            <li key={s.era} className={"ch29-step" + (s.final ? " ch29-step--final" : "")}>
              <div className="ch29-step-era">{s.era}</div>
              <div className="ch29-step-figure">{s.figure}</div>
              <div>
                <p className="ch29-step-note">{s.desc}</p>
                <div className="ch29-wealth-bar">
                  <div className="ch29-wealth-cell ch29-wealth-cell--gain">
                    <span style={{ width: "5.5rem" }}>farmer&#8217;s share</span>
                    <div className="ch29-wealth-meter">
                      <span style={{ width: (s.farmer * 100) + "%" }} />
                    </div>
                    <span className="mono" style={{ width: "2.5rem", textAlign: "right" }}>
                      {Math.round(s.farmer * 100)}%
                    </span>
                  </div>
                  <div className="ch29-wealth-cell ch29-wealth-cell--loss">
                    <span style={{ width: "5.5rem" }}>labour&#8217;s claim</span>
                    <div className="ch29-wealth-meter">
                      <span style={{ width: (s.labour * 100) + "%" }} />
                    </div>
                    <span className="mono" style={{ width: "2.5rem", textAlign: "right" }}>
                      {Math.round(s.labour * 100)}%
                    </span>
                  </div>
                </div>
              </div>
            </li>
          ))}
        </ol>
      </section>

      <section className="card">
        <h2>The Price Revolution as Windfall</h2>
        <p className="description">
          Long leases at fixed money rents, struck before the 16th-century inflation, left tenants
          paying historic prices for produce sold at current ones. The chapter cites Thorold
          Rogers&#8217;s wage and price series; the inputs below default to those figures.
        </p>
        <form className="form-grid" onSubmit={(e) => e.preventDefault()}>
          <label>
            <span>Daily money wage · 1500 (pence)</span>
            <input type="number" min={1} value={wage1500} onChange={(e) => setWage1500(+e.target.value)} />
          </label>
          <label>
            <span>Daily money wage · 1750 (pence)</span>
            <input type="number" min={1} value={wage1750} onChange={(e) => setWage1750(+e.target.value)} />
          </label>
          <label>
            <span>Food-price index · 1500</span>
            <input type="number" min={1} value={price1500} onChange={(e) => setPrice1500(+e.target.value)} />
          </label>
          <label>
            <span>Food-price index · 1750</span>
            <input type="number" min={1} value={price1750} onChange={(e) => setPrice1750(+e.target.value)} />
          </label>
          <label className="span2">
            <span>Share of price-rise NOT captured by rent</span>
            <input
              type="number"
              step={0.05}
              min={0}
              max={1}
              value={rentReal}
              onChange={(e) => setRentReal(+e.target.value)}
            />
          </label>
        </form>

        <div className="ch32-stats" style={{ marginTop: "1.5rem" }}>
          <Stat label="Money-wage rise"      value={`× ${wageRise.toFixed(2)}`} />
          <Stat label="Food-price rise"      value={`× ${priceRise.toFixed(2)}`} tone="red" />
          <Stat label="Real wage (1500 = 1)" value={realWage.toFixed(2)} tone="red" />
          <Stat label="Tenant windfall"      value={`+${farmerGain.toFixed(0)}%`} tone="gold" />
        </div>

        <p className="note" style={{ marginTop: "1.25rem", maxWidth: "62ch" }}>
          On these inputs the real wage of the agricultural labourer falls to
          <span className="mono"> {realWage.toFixed(2)} </span>
          of its 1500 level &#8212; a loss of
          <span className="mono"> {fatten.toFixed(0)}% </span>
          &#8212; while the difference is transferred, at the going lease, to the tenant farmer.
          The capitalist farmer is, to a first approximation, the difference between two price
          series.
        </p>
      </section>

      <section className="card">
        <h2>Comparative Wealth</h2>
        <p className="description">
          The same data, drawn as a single trajectory. The farmer&#8217;s share is what remains
          after fixed rent and wages have been paid out of the product; the labourer&#8217;s share
          is the day-wage purchasing power.
        </p>
        <table className="commodity-table">
          <thead>
            <tr>
              <th>Era</th>
              <th>Figure</th>
              <th style={{ textAlign: "right" }}>Farmer&#8217;s share</th>
              <th style={{ textAlign: "right" }}>Labour&#8217;s claim</th>
              <th>Mechanism in play</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td className="snlt-cell">1300</td>
              <td>Bailiff</td>
              <td className="snlt-cell" style={{ textAlign: "right" }}>8%</td>
              <td className="snlt-cell" style={{ textAlign: "right" }}>95%</td>
              <td className="muted small">manorial obligation</td>
            </tr>
            <tr>
              <td className="snlt-cell">1400s</td>
              <td>M&#233;tayer</td>
              <td className="snlt-cell" style={{ textAlign: "right" }}>20%</td>
              <td className="snlt-cell" style={{ textAlign: "right" }}>78%</td>
              <td className="muted small">share-cropping</td>
            </tr>
            <tr>
              <td className="snlt-cell">1450–1550</td>
              <td>Yeoman</td>
              <td className="snlt-cell" style={{ textAlign: "right" }}>38%</td>
              <td className="snlt-cell" style={{ textAlign: "right" }}>58%</td>
              <td className="muted small">freehold · vulnerable</td>
            </tr>
            <tr>
              <td className="snlt-cell">1550–1650</td>
              <td>Lease-holder</td>
              <td className="snlt-cell" style={{ textAlign: "right" }}>62%</td>
              <td className="snlt-cell" style={{ textAlign: "right" }}>34%</td>
              <td className="muted small">price-revolution windfall</td>
            </tr>
            <tr>
              <td className="snlt-cell">1650–1750</td>
              <td style={{ color: "var(--gold-bright)" }}>Capitalist farmer</td>
              <td className="snlt-cell" style={{ textAlign: "right", color: "var(--gold-bright)" }}>90%</td>
              <td className="snlt-cell" style={{ textAlign: "right", color: "var(--red)" }}>12%</td>
              <td className="muted small">wage-labour at rack-rent</td>
            </tr>
          </tbody>
        </table>
      </section>
    </>
  );
}
