import { useId, useState } from "react";

interface GlossProps {
  /** Short uppercase label beside the ❡ button. */
  label: string;
  /** Verbatim Marx quote; rendered Playfair italic. */
  quote: string;
  /** Mono citation line. Empty string hides the citation. */
  citation: string;
  /** Optional circuit-form badge shown above the quote when open. */
  form?: string;
}

/** A single gold ❡ marker that expands an accessible gloss (a Marx quote +
 *  citation). The recurring T3 device that carries the qualitative load across
 *  the descent. CSS drives the max-height reveal; reduced-motion drops it. */
export function Gloss({ label, quote, citation, form }: GlossProps) {
  const [open, setOpen] = useState(false);
  const revealId = useId();
  return (
    <div className="gloss">
      <button
        className="gloss-btn"
        aria-expanded={open}
        aria-controls={revealId}
        onClick={() => setOpen((v) => !v)}
      >
        <span className="gloss-icon" aria-hidden="true">
          &#10081;
        </span>
        <span className="gloss-label">{label}</span>
      </button>
      <div
        id={revealId}
        role="region"
        aria-label={label}
        className={"gloss-reveal" + (open ? " open" : "")}
      >
        {form && <span className="gloss-form">{form}</span>}
        <p className="gloss-quote">
          <em>&ldquo;{quote}&rdquo;</em>
          {citation && <span className="gloss-cit">{citation}</span>}
        </p>
      </div>
    </div>
  );
}
