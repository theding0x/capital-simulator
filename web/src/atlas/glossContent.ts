/** Atlas — Total Representation · R2-2. The T3 fetishism gloss content.
 *
 * The qualitative "minor concepts" load carried by a single recurring device —
 * a gold ❡ marker that reveals a Marx quote + citation. Keyed by chapter id
 * (`v{volume}-ch{nn}`, zero-padded), matching SpineMatrix's chapterId().
 *
 * The four great mystifications form the fetishism thread:
 *   Ch.1 commodity fetishism → Ch.24 the capital fetish (M…M′)
 *   → Ch.48 the trinity formula → Ch.50 illusions of competition.
 */

export interface GlossContent {
  /** Short label beside the ❡ marker. */
  label: string;
  /** Verbatim Marx quote; rendered Playfair italic. */
  quote: string;
  /** Mono citation line, e.g. "Capital I · Ch. 1 §4". Empty hides the line. */
  citation: string;
  /** Optional circuit-form badge shown above the quote when open. */
  form?: string;
}

/** Shown when a T3 chapter has no explicit gloss and no fallback note. */
export const GLOSS_FALLBACK: GlossContent = {
  label: "gloss",
  quote: "See the chapter text for the qualitative argument at this node.",
  citation: "",
};

export const GLOSS_MAP: Record<string, GlossContent> = {
  "v1-ch01": {
    label: "commodity fetishism",
    quote:
      "A definite social relation between men assumes, in their eyes, the fantastic form of a relation between things.",
    citation: "Capital I · Ch. 1 §4",
    form: "C — C",
  },
  "v1-ch06": {
    label: "the double freedom",
    quote:
      "The labourer is free in the double sense: free to dispose of his labour-power as his own commodity, and free of all the means of production needed to realise it.",
    citation: "Capital I · Ch. 6",
    form: "M — C(Lp)",
  },
  "v3-ch24": {
    label: "the capital fetish",
    quote:
      "In interest-bearing capital — money breeding money, M…M′ — the production process has vanished entirely; it is the most externalised, most fetish-like form of capital.",
    citation: "Capital III · Ch. 24",
    form: "M — M′",
  },
  "v3-ch48": {
    label: "the trinity formula",
    quote:
      "Capital→interest, Land→rent, Labour→wages: the ideological synthesis — an enchanted, perverted, topsy-turvy world — each term an impossible combination of incommensurables.",
    citation: "Capital III · Ch. 48",
    form: "the trinity",
  },
  "v3-ch50": {
    label: "illusions of competition",
    quote:
      "The cost-price fetish: c+v appears as the commodity's own real value, and profit a mere addition made in exchange — the surface where competition reverses every relation.",
    citation: "Capital III · Ch. 50",
    form: "competition",
  },
};
