import { useState } from "react";
import "./Ch26PrimitiveAccumulation.css";

const MYTHS = [
  {
    id: "diligent",
    idyll:
      "There were two sorts of people; one, the diligent, intelligent, and above all frugal élite; the other, lazy rascals spending their substance in riotous living.",
    record:
      "The expropriators of the Tudor period were neither frugal nor diligent; they were the great feudal lords who, by force, drove the peasantry from the land in defiance of their tenures.",
    year: "1485–1547",
    evidence: "Henry VII · suppression of retainers · 1488",
  },
  {
    id: "saved",
    idyll:
      "The former sort accumulated wealth; the latter sort had at last nothing to sell except their own skins.",
    record:
      "The peasants did not sell their skins. The skins were taken — common land enclosed by statute, monastic estates seized by Reformation, freeholds dissolved into sheep-walks.",
    year: "1530–1601",
    evidence: "Dissolution of monasteries · 1536 · 1539",
  },
  {
    id: "moral",
    idyll:
      "From this original sin dates the poverty of the great majority, who, despite all their labour, have up to now nothing to sell but themselves, and the wealth of the few, that increases though they have long ceased to work.",
    record:
      "It was not a sin and there was nothing original about it. The wealth was extracted by direct violence over four hundred years, codified by Parliament, paid for by the public purse, and policed by the gallows.",
    year: "1450–1850",
    evidence: "4,000+ Enclosure Acts · ≈ 7,000,000 acres",
  },
  {
    id: "trickle",
    idyll:
      "Such insipid childishness is every day preached to us in the defence of property.",
    record:
      "The conquest of the New World and the looting of the East Indies furnished the silver, the slaves and the markets without which the European workshop could not have absorbed its own dispossessed.",
    year: "1492–1750",
    evidence: "≈ 16,000 tons Potosí silver · 12.5 M Africans enslaved",
  },
];

type PillId = "diligent" | "saved" | "moral" | "trickle";

const PILLS: { id: PillId; tag: string; label: string }[] = [
  { id: "diligent", tag: "I",   label: "Feudal expropriation" },
  { id: "saved",    tag: "II",  label: "Reformation seizure" },
  { id: "moral",    tag: "III", label: "Parliamentary enclosure" },
  { id: "trickle",  tag: "IV",  label: "Colonial extraction" },
];

export function Ch26PrimitiveAccumulation() {
  const [selectedId, setSelectedId] = useState<PillId>("diligent");

  return (
    <>
      <section className="card">
        <h2>The Question Posed</h2>
        <p className="description">
          Capital presupposes wage-labour; wage-labour presupposes capital. Each is the precondition
          of the other. The circular argument has to be broken somewhere. The breaking-point — the
          prior process by which the producer was separated from the means of production — is what
          Marx calls primitive accumulation.
        </p>
        <p className="note" style={{ maxWidth: "62ch", marginTop: "0.5rem" }}>
          Political economy explains the origin of capital by reference to a moral fable. The fable
          answers a chronological question with a character study. This chapter substitutes for the
          fable an account drawn from acts of Parliament, parish registers, ship manifests and the
          gallows.
        </p>
      </section>

      <section className="card">
        <h2>Idyll &#x27F6; Record</h2>
        <p className="description">
          The left column reproduces the explanation political economy offers for the origin of the
          wage-relation. The right column reproduces what the archives show.
        </p>
        <ul className="ch26-myths">
          {MYTHS.map((m) => (
            <li key={m.id} className="ch26-myth">
              <div className="ch26-myth-side ch26-myth-side--idyll">
                <div className="ch26-myth-tag">&#167; idyll &middot; the fable</div>
                <p className="ch26-myth-claim">&#8220;{m.idyll}&#8221;</p>
                <div className="ch26-myth-evidence muted small">attributed to political economy</div>
              </div>
              <div className="ch26-myth-divider">vs.</div>
              <div className="ch26-myth-side ch26-myth-side--record">
                <div className="ch26-myth-tag">&#167; record &middot; {m.year}</div>
                <p className="ch26-myth-claim">{m.record}</p>
                <div className="ch26-myth-evidence">{m.evidence}</div>
              </div>
            </li>
          ))}
        </ul>
      </section>

      <section className="card">
        <h2>Filter the Archive</h2>
        <p className="description">
          Each row of the table below names one episode of primitive accumulation. The selected row
          is the seed Marx returns to across Chapters 27–31.
        </p>
        <div className="ch17-presets" style={{ marginBottom: 0 }}>
          {PILLS.map((p) => (
            <button
              key={p.id}
              type="button"
              className={"ch17-preset" + (selectedId === p.id ? " ch17-preset--active" : "")}
              onClick={() => setSelectedId(p.id)}
            >
              <span className="ch17-preset-tag">{p.tag}</span>
              <span className="ch17-preset-label">{p.label}</span>
            </button>
          ))}
        </div>

        <table className="commodity-table" style={{ marginTop: "1.25rem" }}>
          <thead>
            <tr>
              <th style={{ width: "8rem" }}>Episode</th>
              <th>Mechanism</th>
              <th style={{ textAlign: "right" }}>Period</th>
              <th style={{ textAlign: "right" }}>Read in</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td className="snlt-cell">I</td>
              <td>Feudal lords disperse retainers; common rights overridden by force.</td>
              <td className="snlt-cell" style={{ textAlign: "right" }}>1450–1547</td>
              <td className="snlt-cell" style={{ textAlign: "right" }}>Ch. 27</td>
            </tr>
            <tr>
              <td className="snlt-cell">II</td>
              <td>Bloody legislation forces vagabonds into wage discipline.</td>
              <td className="snlt-cell" style={{ textAlign: "right" }}>1530–1601</td>
              <td className="snlt-cell" style={{ textAlign: "right" }}>Ch. 28</td>
            </tr>
            <tr>
              <td className="snlt-cell">III</td>
              <td>Tenant farmers re-form as a capitalist class on long leases.</td>
              <td className="snlt-cell" style={{ textAlign: "right" }}>1500–1750</td>
              <td className="snlt-cell" style={{ textAlign: "right" }}>Ch. 29</td>
            </tr>
            <tr>
              <td className="snlt-cell">IV</td>
              <td>Rural surplus becomes industrial proletariat &amp; home market.</td>
              <td className="snlt-cell" style={{ textAlign: "right" }}>1600–1800</td>
              <td className="snlt-cell" style={{ textAlign: "right" }}>Ch. 30</td>
            </tr>
            <tr>
              <td className="snlt-cell">V</td>
              <td>Colonial system, public debt, taxes, protection found industrial capital.</td>
              <td className="snlt-cell" style={{ textAlign: "right" }}>1492–1800</td>
              <td className="snlt-cell" style={{ textAlign: "right" }}>Ch. 31</td>
            </tr>
          </tbody>
        </table>

        <div className="ch26-coda">
          <p>
            &#8220;In actual history, it is notorious that conquest, enslavement, robbery, murder,
            briefly force, play the great part. In the tender annals of political economy, the
            idyllic reigns from time immemorial.&#8221;
          </p>
        </div>
      </section>
    </>
  );
}
