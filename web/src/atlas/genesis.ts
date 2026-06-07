/** Atlas — Total Representation · R2-4. The T4 genesis stratum: the historical
 *  bedrock beneath the living world. Primitive accumulation (Vol I Ch.26–33) and
 *  its pre-capitalist antecedents (Vol III Ch.20, 36, 47) as a walkable timeline.
 *
 *  Static narrative — NOT part of the live run (per the data contract). Each
 *  episode links to its chapter, opened in the shared ChapterDrawer.
 *
 *  "So-called primitive accumulation … is nothing else than the historical
 *  process of divorcing the producer from the means of production." — Ch.26
 */

export type GenesisKind =
  | "antecedent" // pre-capitalist forms (Vol III historical)
  | "expropriation" // the land torn from the people
  | "legislation" // the bloody discipline of the dispossessed
  | "colonial" // plunder of the world market
  | "credit" // the national debt as a lever
  | "tendency"; // the negation of the negation

export interface GenesisEpisode {
  id: string;
  /** Axis label — a year ("1349") or an era ("Antiquity", "The tendency"). */
  era: string;
  title: string;
  kind: GenesisKind;
  /** One-line gloss in Marx's register. */
  blurb: string;
  citation: string;
  /** The chapter whose panel this episode opens. */
  volume: 1 | 2 | 3;
  ch: number;
  /** Historical weight, 1 (slight) … 5 (epochal) — drives the node's size. */
  magnitude: number;
}

export const GENESIS_EPISODES: GenesisEpisode[] = [
  {
    id: "g-merchant",
    era: "Antiquity",
    title: "Merchant's capital",
    kind: "antecedent",
    blurb: "M—C—M′ older than capital itself — commercial wealth in the pores of ancient society.",
    citation: "Capital III · Ch. 20",
    volume: 3,
    ch: 20,
    magnitude: 2,
  },
  {
    id: "g-usury",
    era: "Antiquity",
    title: "Usurer's capital",
    kind: "antecedent",
    blurb: "Usury — the antediluvian form of interest-bearing capital — bleeds the small producer.",
    citation: "Capital III · Ch. 36",
    volume: 3,
    ch: 36,
    magnitude: 2,
  },
  {
    id: "g-money-rent",
    era: "Medieval",
    title: "Labour-rent → money-rent",
    kind: "antecedent",
    blurb: "Ground-rent passes from labour-services to rent in kind to money — dissolving feudal bonds.",
    citation: "Capital III · Ch. 47",
    volume: 3,
    ch: 47,
    magnitude: 3,
  },
  {
    id: "g-1349",
    era: "1349",
    title: "Statute of Labourers",
    kind: "legislation",
    blurb: "A legal wage-maximum — the first bloody discipline of the newly dispossessed.",
    citation: "Capital I · Ch. 28",
    volume: 1,
    ch: 28,
    magnitude: 3,
  },
  {
    id: "g-enclosure",
    era: "1470s",
    title: "Enclosure of the commons",
    kind: "expropriation",
    blurb: "The agricultural people forcibly torn from the soil; arable land turned to sheep-walks.",
    citation: "Capital I · Ch. 27",
    volume: 1,
    ch: 27,
    magnitude: 5,
  },
  {
    id: "g-reformation",
    era: "1530s",
    title: "Reformation plunder",
    kind: "expropriation",
    blurb: "Church estates seized and squandered; their hereditary tenants flung into the proletariat.",
    citation: "Capital I · Ch. 27",
    volume: 1,
    ch: 27,
    magnitude: 4,
  },
  {
    id: "g-bloody",
    era: "1547",
    title: "Bloody legislation",
    kind: "legislation",
    blurb: "Vagabonds whipped, branded, enslaved for the crime of their own expropriation.",
    citation: "Capital I · Ch. 28",
    volume: 1,
    ch: 28,
    magnitude: 4,
  },
  {
    id: "g-colonial",
    era: "1600",
    title: "The colonial system",
    kind: "colonial",
    blurb: "Chartered companies; the world market made a hunting-ground for plunder and enslavement.",
    citation: "Capital I · Ch. 31",
    volume: 1,
    ch: 31,
    magnitude: 4,
  },
  {
    id: "g-debt",
    era: "1694",
    title: "National debt & the Bank",
    kind: "credit",
    blurb: "Public credit becomes the lever of accumulation; the state itself mortgaged to capital.",
    citation: "Capital I · Ch. 31",
    volume: 1,
    ch: 31,
    magnitude: 4,
  },
  {
    id: "g-bengal",
    era: "1757",
    title: "Plunder of Bengal",
    kind: "colonial",
    blurb: "Treasures captured outside Europe flow back as capital — and light the dawn of industry.",
    citation: "Capital I · Ch. 31",
    volume: 1,
    ch: 31,
    magnitude: 5,
  },
  {
    id: "g-clearances",
    era: "1814",
    title: "The Highland Clearances",
    kind: "expropriation",
    blurb: "The Duchess of Sutherland clears 15,000 from the land for sheep-runs and deer-forests.",
    citation: "Capital I · Ch. 27",
    volume: 1,
    ch: 27,
    magnitude: 5,
  },
  {
    id: "g-peel",
    era: "1829",
    title: "Peel's fiasco (Swan River)",
    kind: "colonial",
    blurb: "Mr Peel's labourers desert him: capital is revealed as a social relation, not a thing.",
    citation: "Capital I · Ch. 33",
    volume: 1,
    ch: 33,
    magnitude: 3,
  },
  {
    id: "g-poorlaw",
    era: "1834",
    title: "The New Poor Law",
    kind: "legislation",
    blurb: "The workhouse — relief made a terror, the reserve army disciplined into the factory.",
    citation: "Capital I · Ch. 28",
    volume: 1,
    ch: 28,
    magnitude: 3,
  },
  {
    id: "g-tendency",
    era: "The tendency",
    title: "The expropriators are expropriated",
    kind: "tendency",
    blurb: "Centralisation reaches its limit; the knell of capitalist private property sounds.",
    citation: "Capital I · Ch. 32",
    volume: 1,
    ch: 32,
    magnitude: 5,
  },
];
