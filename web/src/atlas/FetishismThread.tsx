import { useState } from "react";

interface FetishNode {
  id: string;
  /** Volume·chapter reference, e.g. "I · Ch. 1". */
  ref: string;
  title: string;
  /** The circuit-form the mystification takes. */
  form: string;
  quote: string;
  citation: string;
}

/** The four great mystifications down the descent — the qualitative spine of
 *  the fetishism thread. Verbatim glosses from the source text. */
const NODES: FetishNode[] = [
  {
    id: "ch1",
    ref: "I · Ch. 1",
    title: "Commodity fetishism",
    form: "C — C",
    quote:
      "A definite social relation between men assumes, in their eyes, the fantastic form of a relation between things.",
    citation: "Capital I · Ch. 1 §4",
  },
  {
    id: "ch24",
    ref: "III · Ch. 24",
    title: "The capital fetish",
    form: "M — M′",
    quote:
      "Money breeding money — M…M′ — the production process vanished entirely. The most externalised, most fetish-like form of capital.",
    citation: "Capital III · Ch. 24",
  },
  {
    id: "ch48",
    ref: "III · Ch. 48",
    title: "The trinity formula",
    form: "the trinity",
    quote:
      "Capital→interest, Land→rent, Labour→wages: the ideological synthesis — an enchanted, perverted, topsy-turvy world — each term an impossible combination.",
    citation: "Capital III · Ch. 48",
  },
  {
    id: "ch50",
    ref: "III · Ch. 50",
    title: "Illusions of competition",
    form: "competition",
    quote:
      "The cost-price fetish: c+v appears as the commodity's real value, and profit a mere addition made in exchange.",
    citation: "Capital III · Ch. 50",
  },
];

/** The fetishism thread — a single connective T3 motif binding the four great
 *  mystifications across the whole descent. A gold ❡ at each node; selecting one
 *  reveals its gloss. Distinct from the standalone Gloss marker by its
 *  one-active-at-a-time selector pattern. */
export function FetishismThread() {
  const [active, setActive] = useState<string>("ch1");
  const node = NODES.find((n) => n.id === active) ?? NODES[0];
  return (
    <section className="fetish-thread" aria-label="The fetishism thread">
      <div className="ft-eyebrow">
        The fetishism thread &mdash; one connective gloss across the descent
        &nbsp;&middot;&nbsp; Ch.1 &rarr; 24 &rarr; 48 &rarr; 50
      </div>
      <div className="ft-nodes" role="group" aria-label="The four great mystifications">
        {NODES.map((n) => {
          const on = n.id === active;
          return (
            <button
              key={n.id}
              className={"ft-node" + (on ? " active" : "")}
              aria-pressed={on}
              aria-label={`${n.ref} — ${n.title}`}
              onClick={() => setActive(n.id)}
            >
              <span className="ft-icon" aria-hidden="true">
                &#10081;
              </span>
              <span className="ft-ref">{n.ref}</span>
              <span className="ft-title">{n.title}</span>
            </button>
          );
        })}
      </div>
      <div className="ft-gloss-area" aria-live="polite" aria-atomic="true">
        <span className="ft-form">{node.form}</span>
        <p className="ft-quote">
          <em>&ldquo;{node.quote}&rdquo;</em>
          <span className="ft-cit">{node.citation}</span>
        </p>
      </div>
      <p className="ft-cap">
        A single recurring device &mdash; a gold &#10081; that opens a gloss &mdash; binds
        the four great mystifications down the descent, carrying much of the qualitative load.
      </p>
    </section>
  );
}
