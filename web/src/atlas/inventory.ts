/** Atlas — Total Representation. The complete 106-chapter census of Capital,
 * each entry assigned a representation tier (T1 live · T2 panel · T3 gloss ·
 * T4 historical) and, where live, a data source.
 *
 * status:  done    = already represented in the Observatory
 *          partial = partially represented
 *          new     = net-new
 */

export type Tier = "T1" | "T2" | "T3" | "T4";
export type InventoryStatus = "done" | "partial" | "new";
export type InventoryNode =
  | "M"
  | "M-C"
  | "P"
  | "C′"
  | "C-M′"
  | "M′"
  | "ΔM"
  | "whole"
  | "historical";
export type ServiceKey =
  | "commodity"
  | "agent"
  | "market"
  | "sim"
  | "finance"
  | "none";

export interface InventoryEntry {
  volume: 1 | 2 | 3;
  part: string;
  ch: number;
  title: string;
  nodes: InventoryNode[];
  status: InventoryStatus;
  tiers: Tier[];
  src: ServiceKey;
  note: string;
}

// Circuit-moment ordering for the spine's horizontal axis.
export const CIRCUIT_NODES: InventoryNode[] = [
  "M", "M-C", "P", "C′", "C-M′", "M′", "ΔM", "whole", "historical",
];

// Service metadata.
export const SERVICES: Record<ServiceKey, { name: string; port: number | null }> = {
  commodity: { name: "commodity-service", port: 8081 },
  agent:     { name: "agent-service",     port: 8082 },
  market:    { name: "market-service",    port: 8083 },
  sim:       { name: "simulation-engine", port: 8084 },
  finance:   { name: "finance-service",   port: 8085 },
  none:      { name: "—",                 port: null  },
};

function R(
  volume: 1 | 2 | 3,
  part: string,
  ch: number,
  title: string,
  nodes: InventoryNode[],
  status: InventoryStatus,
  tiers: Tier[],
  src: ServiceKey,
  note: string
): InventoryEntry {
  return { volume, part, ch, title, nodes, status, tiers, src, note };
}

export const INVENTORY: InventoryEntry[] = [
  /* ─────────────── VOLUME I — Production ─────────────── */
  R(1, "I — Commodities & Money", 1, "The Commodity", ["whole"], "partial", ["T2","T3","T1"], "commodity",
    "Value-form ladder simple→money (T2); commodity-fetishism gloss (T3); live value magnitude vs SNLT (T1)."),
  R(1, "I — Commodities & Money", 2, "Exchange", ["M-C","C-M′"], "new", ["T2","T3"], "market",
    "Barter-ratio panel; money crystallising out of exchange as a genesis gloss."),
  R(1, "I — Commodities & Money", 3, "Money, or the Circulation of Commodities", ["M","M′"], "partial", ["T2","T3"], "market",
    "Money-required gauge (M=ΣP÷velocity); functions-of-money gloss (measure, medium, hoard, payment, world)."),

  R(1, "II — Money into Capital", 4, "The General Formula for Capital", ["whole"], "done", ["T1"], "sim",
    "★ The orrery rings — C-M-C vs M-C-M′; ΔM surfaced as the surplus the circuit returns."),
  R(1, "II — Money into Capital", 5, "Contradictions in the General Formula", ["whole"], "new", ["T2","T3"], "market",
    "Value-conservation demo (equivalents/non-equivalents both conserve ΣV); the contradiction as a gloss."),
  R(1, "II — Money into Capital", 6, "The Sale and Purchase of Labour-Power", ["M-C"], "done", ["T1","T2","T3"], "agent",
    "★ The threshold gate; subsistence-basket panel (value of LP = SNLT of the basket); double-freedom gloss."),

  R(1, "III — Absolute Surplus-Value", 7, "The Labour-Process and Valorization", ["P"], "done", ["T1","T2"], "sim",
    "★ The working-day hero; c+v+s composition panel; surplus-value as unpaid labour."),
  R(1, "III — Absolute Surplus-Value", 8, "Constant Capital and Variable Capital", ["P"], "done", ["T1","T2"], "sim",
    "★ Organic composition c/v in the abode; c+v+s decomposition panel."),
  R(1, "III — Absolute Surplus-Value", 9, "The Rate of Surplus-Value", ["P","ΔM"], "done", ["T1"], "sim",
    "★ The rate of exploitation s/v — the abode's headline figure."),
  R(1, "III — Absolute Surplus-Value", 10, "The Working-Day", ["P"], "done", ["T1","T2","T4"], "sim",
    "★ The contested necessary|surplus split; relay-system panel; Factory Acts on the genesis timeline (T4)."),

  R(1, "III/IV — Rate, Mass & Relative SV", 11, "The Rate and Mass of Surplus-Value", ["P","ΔM"], "new", ["T1","T2"], "sim",
    "Mass S=(s/v)·V live; compensation-law panel (rate↑↔workers↓ holds S, within the 24h limit)."),

  R(1, "IV — Relative Surplus-Value", 12, "The Concept of Relative Surplus-Value", ["P"], "new", ["T2"], "sim",
    "Absolute/relative split; extra-surplus-value (the innovator's temporary bonus) panel."),
  R(1, "IV — Relative Surplus-Value", 13, "Co-operation", ["P"], "new", ["T2","T3"], "agent",
    "Collective-power panel (sum > parts); the bonus appearing as a property of capital (gloss)."),
  R(1, "IV — Relative Surplus-Value", 14, "Division of Labour and Manufacture", ["P"], "new", ["T2","T3"], "agent",
    "Manufacture / proportionality panel; deskilling & the two divisions of labour as glosses."),
  R(1, "IV — Relative Surplus-Value", 15, "Machinery and Modern Industry", ["P"], "partial", ["T1","T2","T4"], "sim",
    "◐ Machine wear feeds organic composition; labour-displaced→reserve army (T1 ★); industrial-cycle panel/timeline."),

  R(1, "V — Absolute & Relative SV", 16, "Absolute and Relative Surplus-Value", ["P","ΔM"], "done", ["T1","T2","T3"], "sim",
    "★ The working-day split; formal vs real subsumption gloss; profit-vs-surplus rates panel."),
  R(1, "V — Absolute & Relative SV", 17, "Changes of Magnitude in Price of LP & SV", ["P","ΔM"], "partial", ["T1","T2"], "sim",
    "Magnitude scenarios (length·intensity·productiveness); ◐ reserve-army wage compression feeds the abode."),
  R(1, "V — Absolute & Relative SV", 18, "Various Formulae for the Rate of Surplus-Value", ["ΔM"], "new", ["T2","T3"], "sim",
    "Formula comparison (I s/v · II s/day · III unpaid/paid); how Formula II conceals exploitation (gloss)."),

  R(1, "VI — Wages", 19, "The Transformation of Value of LP into Wages", ["M-C","ΔM"], "new", ["T1","T3"], "agent",
    "Live paid|unpaid bar; the wage-illusion (wages appear to pay for all labour) gloss."),
  R(1, "VI — Wages", 20, "Time-Wages", ["M-C"], "new", ["T2","T3"], "agent",
    "Hourly-price panel; lengthening the day lowers the hourly price even at constant pay (gloss)."),
  R(1, "VI — Wages", 21, "Piece-Wages", ["M-C"], "new", ["T2","T3"], "agent",
    "Piece-price vs piece-value panel; the sweating / sub-contract system gloss."),
  R(1, "VI — Wages", 22, "National Differences in Wages", ["M-C"], "new", ["T2","T3"], "agent",
    "Wage-comparison table (standardised day); the England paradox (highest nominal, lowest relative) gloss."),

  R(1, "VII — Accumulation", 23, "Simple Reproduction", ["whole"], "new", ["T2","T3"], "sim",
    "Cycle series (revenue vs accumulated; repayment period C÷S); reproduction of the class relation (gloss)."),
  R(1, "VII — Accumulation", 24, "The Transformation of Surplus-Value into Capital", ["whole","ΔM"], "done", ["T1"], "sim",
    "★ The accumulation lever α; the compound-growth spiral (capital self-expanding)."),
  R(1, "VII — Accumulation", 25, "The General Law of Capitalist Accumulation", ["whole","ΔM","historical"], "done", ["T1","T2","T4"], "sim",
    "★★ The abode's core — immiseration, reserve army, organic composition (T1); strata & centralisation (T2); England/Ireland illustrations (T4)."),

  R(1, "VIII — Primitive Accumulation", 26, "The Secret of Primitive Accumulation", ["historical"], "new", ["T4","T3"], "none",
    "Genesis-stratum opener; 'capital is a social relation, not a thing' gloss."),
  R(1, "VIII — Primitive Accumulation", 27, "Expropriation of the Agricultural Population", ["historical"], "new", ["T4"], "none",
    "Enclosure / Clearances on the expropriation timeline (already designed: Ch.27 page)."),
  R(1, "VIII — Primitive Accumulation", 28, "Bloody Legislation against the Expropriated", ["historical"], "new", ["T4"], "none",
    "Statute registry with severity meters on the timeline (already designed: Ch.28 page)."),
  R(1, "VIII — Primitive Accumulation", 29, "Genesis of the Capitalist Farmer", ["historical"], "new", ["T4"], "none",
    "Bailiff→capitalist-farmer genealogy; real-wage drift (already designed: Ch.29 page)."),
  R(1, "VIII — Primitive Accumulation", 30, "Reaction on Industry; the Home Market", ["historical"], "new", ["T4","T2"], "none",
    "Rural pool → industrial supply + home-market demand fork; market-formation panel (Ch.30 page)."),
  R(1, "VIII — Primitive Accumulation", 31, "Genesis of the Industrial Capitalist", ["historical"], "new", ["T4"], "none",
    "Four-pillar sources stack (colony·debt·tax·tariff); blood-and-dirt coda (Ch.31 page)."),
  R(1, "VIII — Primitive Accumulation", 32, "The Historical Tendency of Capitalist Accumulation", ["historical"], "new", ["T4","T3"], "none",
    "Dialectical stages → negation of the negation; the expropriators expropriated (gloss)."),
  R(1, "VIII — Primitive Accumulation", 33, "The Modern Theory of Colonisation", ["historical"], "new", ["T4","T3"], "none",
    "Peel's fiasco — capital as a social relation revealed in the colonies (gloss)."),

  /* ─────────────── VOLUME II — Circulation ─────────────── */
  R(2, "I — Metamorphoses & Circuits", 1, "The Circuit of Money-Capital", ["whole","M","M′"], "done", ["T1"], "sim",
    "★ The orrery circuit; six-moment timeline M—C(L+MP)…P…C′—M′."),
  R(2, "I — Metamorphoses & Circuits", 2, "The Circuit of Productive Capital", ["whole","P"], "new", ["T2"], "sim",
    "Reproduction-lens panel; latent / hoarded money-capital gauge."),
  R(2, "I — Metamorphoses & Circuits", 3, "The Circuit of Commodity-Capital", ["whole","C′","C-M′"], "new", ["T2","T3"], "sim",
    "Aliquot decomposition (every £ = c+v+s) panel; the social intertwining (C is another's C′) gloss."),
  R(2, "I — Metamorphoses & Circuits", 4, "The Three Formulas of the Circuit", ["whole"], "done", ["T1","T2"], "sim",
    "★★ The orrery's three coexisting arcs (the Nebeneinander); stage-distribution + a stagnation toy (T1); supply-demand gap (T2)."),
  R(2, "I — Metamorphoses & Circuits", 5, "The Time of Circulation", ["M-C","C-M′"], "new", ["T2"], "market",
    "Time decomposition (production + circulation time); feeds orrery pacing."),
  R(2, "I — Metamorphoses & Circuits", 6, "The Costs of Circulation", ["M-C","C-M′"], "new", ["T2","T3"], "market",
    "Faux-frais panel (pure vs value-adding); transport genuinely adds value (gloss)."),

  R(2, "II — Turnover", 7, "Turnover Time and the Number of Turnovers", ["whole"], "done", ["T1"], "sim",
    "★ The turnover number n=T÷t already paces each orrery capital."),
  R(2, "II — Turnover", 8, "Fixed Capital and Circulating Capital", ["whole"], "partial", ["T2"], "sim",
    "Fixed/circulating panel (orthogonal to c/v); wear & the sinking fund. ◐ wear feeds composition."),
  R(2, "II — Turnover", 9, "The Aggregate Turnover of Advanced Capital", ["whole"], "new", ["T2","T4"], "sim",
    "Aggregate turnover of mixed capital; the fixed-capital lifetime as the material basis of the crisis-cycle (T2/T4)."),
  R(2, "II — Turnover", 10, "Theories of Fixed & Circulating Capital: Physiocrats & Smith", ["whole","historical"], "new", ["T3"], "none",
    "Doctrine-critique gloss (Smith's conflation of fixed/circulating with constant/variable)."),
  R(2, "II — Turnover", 11, "Theories of Fixed & Circulating Capital: Ricardo", ["whole","historical"], "new", ["T3"], "none",
    "Doctrine-critique gloss (Ricardo's collapse onto mere durability)."),
  R(2, "II — Turnover", 12, "The Working Period", ["P"], "new", ["T2"], "sim",
    "Working-period panel (connected working-days; the advance multiplier)."),
  R(2, "II — Turnover", 13, "The Time of Production", ["P"], "new", ["T2"], "sim",
    "Production-time panel (natural-process time; latent productive capital)."),
  R(2, "II — Turnover", 14, "The Time of Circulation", ["M-C","C-M′"], "new", ["T2"], "market",
    "Circulation-time panel with the annual-rate-of-surplus-value penalty from long circulation."),
  R(2, "II — Turnover", 15, "The Effects of a Change of Prices", ["whole"], "new", ["T2","T3"], "market",
    "Price-revolution panel (capital set free / tied up mid-turnover); speculation gloss."),
  R(2, "II — Turnover", 16, "The Turnover of Variable Capital", ["P","ΔM"], "partial", ["T2"], "sim",
    "Annual rate of surplus-value panel (A-vs-B; n× per year). ◐ ties to the orrery turnover."),

  R(2, "III — Total Social Capital", 17, "The Circulation of Surplus-Value", ["ΔM","whole"], "new", ["T3"], "sim",
    "The realisation-money puzzle (where the money to realise surplus comes from) — gloss prefacing Ch.20–21."),
  R(2, "III — Total Social Capital", 18, "The Role of Money-Capital in Reproduction", ["whole"], "new", ["T2"], "market",
    "Money-supply panel; the two Departments enter; the wage-rotation fund and velocity."),
  R(2, "III — Total Social Capital", 19, "Former Presentations of the Subject", ["whole","historical"], "new", ["T3"], "none",
    "Doctrine-history gloss: Quesnay's Tableau; Smith's dogma (the missing constant-capital replacement)."),
  R(2, "III — Total Social Capital", 20, "Simple Reproduction", ["whole"], "done", ["T1"], "sim",
    "★ The two-Department reproduction scheme — I(v+s)↔II(c), the three great exchanges — a prime new live diagram."),
  R(2, "III — Total Social Capital", 21, "Accumulation and Reproduction on an Extended Scale", ["whole","ΔM"], "done", ["T1"], "sim",
    "★ Extended-reproduction growth (Dept-I-leads asymmetry; the composition shift over cycles)."),

  /* ─────────────── VOLUME III — The Totality ─────────────── */
  R(3, "I — Surplus-Value into Profit", 1, "Cost-Price and Profit", ["ΔM"], "partial", ["T2","T3"], "finance",
    "Cost-price k=c+v panel; profit as the mystified form of surplus-value (gloss)."),
  R(3, "I — Surplus-Value into Profit", 2, "The Rate of Profit", ["ΔM"], "done", ["T1","T2","T3"], "finance",
    "★ p′=s÷C is the orrery's centre of gravity; p′-vs-s′ panel; the mystification gap (gloss)."),
  R(3, "I — Surplus-Value into Profit", 3, "Relation of the Rate of Profit to the Rate of Surplus-Value", ["ΔM"], "new", ["T2"], "finance",
    "Interactive p′ = s′·(v/C) panel — vary composition, watch the rate move."),
  R(3, "I — Surplus-Value into Profit", 4, "The Effect of the Turnover on the Rate of Profit", ["M-C","C-M′","ΔM"], "partial", ["T2"], "finance",
    "Annual rate of profit p′=s′·n·(v/C) panel. ◐ turnover."),
  R(3, "I — Surplus-Value into Profit", 5, "Economy in the Employment of Constant Capital", ["P","ΔM"], "new", ["T2"], "finance",
    "Economy-in-c panel (raising p′ without touching s or v)."),
  R(3, "I — Surplus-Value into Profit", 6, "The Effect of Price Fluctuation", ["M-C","ΔM"], "new", ["T2"], "finance",
    "Price-fluctuation panel; capital release (fall) vs tie-up (rise)."),
  R(3, "I — Surplus-Value into Profit", 7, "Supplementary Remarks", ["P","ΔM"], "new", ["T2","T3"], "finance",
    "Composition-effect panel; profit's surface appearance (it seems to spring from acumen) — gloss."),

  R(3, "II — Average Profit", 8, "Different Compositions of Capitals in Different Branches", ["ΔM","whole"], "partial", ["T2"], "finance",
    "Multi-sphere p′ scatter (branches with different organic compositions diverge before equalisation)."),
  R(3, "II — Average Profit", 9, "Formation of a General Rate; Prices of Production", ["ΔM","whole"], "done", ["T1","T2","T3"], "finance",
    "★ p̄′=ΣS÷ΣC is the centre of gravity; prices-of-production panel; the value→price transformation (T2/T3)."),
  R(3, "II — Average Profit", 10, "Equalisation of the General Rate through Competition", ["C-M′","ΔM"], "done", ["T1","T2"], "finance",
    "★ Capital flow between spheres IS the orrery's centre-of-gravity dynamics; market value & surplus-profit (T2)."),
  R(3, "II — Average Profit", 11, "Effects of General Wage Fluctuations on Prices of Production", ["ΔM"], "new", ["T2"], "finance",
    "Wage→price effect panel (the general rate moves inversely to wages)."),
  R(3, "II — Average Profit", 12, "Supplementary Remarks", ["ΔM"], "new", ["T2","T3"], "finance",
    "Compensation-grounds (turnover/risk/skill = a share in aggregate surplus, not value creation) — myth gloss."),

  R(3, "III — Tendential Fall (TRPF)", 13, "The Law As Such", ["ΔM","whole"], "done", ["T1"], "finance",
    "★ Rising composition + constant s′ → falling p′; the rate–mass contradiction (mass rises as the rate falls) — live trajectory."),
  R(3, "III — Tendential Fall (TRPF)", 14, "Counteracting Influences", ["ΔM"], "new", ["T2"], "finance",
    "Six counteracting-force toggles (the law is retarded, not abolished)."),
  R(3, "III — Tendential Fall (TRPF)", 15, "Internal Contradictions of the Law", ["ΔM","whole"], "done", ["T1","T2","T3"], "finance",
    "★ Crisis as the violent resolution — devalues c, restores p′ (T1); concentration/centralisation (T2); the contradictions (T3)."),

  R(3, "IV — Merchant's Capital", 16, "Commercial Capital", ["C-M′","ΔM"], "new", ["T2","T3"], "finance",
    "Merchant-circuit panel (M—C—M′ mediating C′→M′); it creates no new value (gloss)."),
  R(3, "IV — Merchant's Capital", 17, "Commercial Profit", ["C-M′","ΔM"], "new", ["T2","T3"], "finance",
    "Adjusted general rate (merchant capital in the denominator lowers p̄′); commercial-worker exploitation (gloss)."),
  R(3, "IV — Merchant's Capital", 18, "The Turnover of Merchant's Capital; Prices", ["M","C-M′","ΔM"], "new", ["T2"], "finance",
    "Turnover→markup panel (per-unit markup inversely proportional to turnover)."),
  R(3, "IV — Merchant's Capital", 19, "Money-Dealing Capital", ["M","ΔM"], "new", ["T2"], "finance",
    "Money-dealing panel (the technical operations of money circulation; M—M′; creates no value)."),
  R(3, "IV — Merchant's Capital", 20, "Historical Facts about Merchant's Capital", ["historical","ΔM"], "new", ["T4","T2"], "finance",
    "The oldest independent form on the genesis timeline (T4); its subordination as capitalism matures (T2)."),

  R(3, "V — Interest, Credit & Fictitious Capital", 21, "Interest-Bearing Capital", ["M","M′"], "done", ["T2","T3"], "finance",
    "★ Loan→M+ΔM panel; loanable capital as a commodity sui generis; the M—M′ form (gloss)."),
  R(3, "V — Interest, Credit & Fictitious Capital", 22, "Division of Profit; the Rate of Interest", ["M","M′"], "new", ["T2"], "finance",
    "Interest-rate panel (no natural law; set by the money market; cycle-driven; ceiling = average profit)."),
  R(3, "V — Interest, Credit & Fictitious Capital", 23, "Interest and Profit of Enterprise", ["M","M′"], "new", ["T2","T3"], "finance",
    "Profit-division panel (interest vs profit-of-enterprise); the ownership/management inversion (gloss)."),
  R(3, "V — Interest, Credit & Fictitious Capital", 24, "Externalisation of the Relations of Capital (Fetish Capital)", ["M","M′"], "done", ["T3","T2"], "finance",
    "★ The capital fetish — M—M′, money breeding money (gloss); the compound-interest absurdity (Dr Price) panel."),
  R(3, "V — Interest, Credit & Fictitious Capital", 25, "Credit and Fictitious Capital", ["ΔM"], "done", ["T2","T3"], "finance",
    "★ Fictitious capital / capitalisation = income÷rate (T2); fictitious-vs-real (gloss); the 1847 crisis."),
  R(3, "V — Interest, Credit & Fictitious Capital", 26, "Accumulation of Money-Capital; Influence on Interest", ["M","M′","ΔM"], "new", ["T2"], "finance",
    "Loanable-vs-real accumulation panel (they diverge through the cycle)."),
  R(3, "V — Interest, Credit & Fictitious Capital", 27, "The Role of Credit in Capitalist Production", ["ΔM"], "done", ["T3"], "finance",
    "★ Joint-stock socialisation & the cooperative factory — the antithesis transcended within capitalism (glosses)."),
  R(3, "V — Interest, Credit & Fictitious Capital", 28, "Medium of Circulation and Capital (Tooke & Fullarton)", ["M","ΔM"], "new", ["T2","T3"], "finance",
    "Currency-vs-capital panel; doctrine gloss."),
  R(3, "V — Interest, Credit & Fictitious Capital", 29, "Component Parts of Bank Capital", ["ΔM"], "new", ["T2"], "finance",
    "Bank-capital decomposition (the double-counting of the same money as capital)."),
  R(3, "V — Interest, Credit & Fictitious Capital", 30, "Money-Capital and Real Capital, I", ["M","ΔM"], "new", ["T2"], "finance",
    "Real-capital accumulation vs commercial credit; the return flow."),
  R(3, "V — Interest, Credit & Fictitious Capital", 31, "Money-Capital and Real Capital, II", ["M","ΔM"], "new", ["T2"], "finance",
    "Floating capital; bill-brokers; rediscount; loan-capital velocity."),
  R(3, "V — Interest, Credit & Fictitious Capital", 32, "Money-Capital and Real Capital, III", ["M","ΔM"], "new", ["T2"], "finance",
    "Revenue→loan-capital conversion; the rentier; loan-capital overstatement."),
  R(3, "V — Interest, Credit & Fictitious Capital", 33, "The Medium of Circulation in the Credit System", ["M","ΔM"], "new", ["T2"], "finance",
    "Clearing-house settlement / velocity panel (millions settled with little actual money)."),
  R(3, "V — Interest, Credit & Fictitious Capital", 34, "The Currency Principle and the Bank Act of 1844", ["M","ΔM","historical"], "new", ["T2","T3","T4"], "finance",
    "Bank-Act regime panel (T2); the currency-principle critique (T3); the 1844 episode on the timeline (T4)."),
  R(3, "V — Interest, Credit & Fictitious Capital", 35, "Precious Metal and the Rate of Exchange", ["M"], "new", ["T2"], "finance",
    "Exchange-rates / gold-flow panel (balance of payments; world-money function)."),
  R(3, "V — Interest, Credit & Fictitious Capital", 36, "Pre-Capitalist Relationships", ["historical"], "new", ["T4","T3"], "finance",
    "Usurer's capital on the genesis timeline (T4); usury's subordination (gloss)."),

  R(3, "VI — Ground-Rent", 37, "Introduction", ["ΔM"], "new", ["T3"], "finance",
    "The rent triad (farmer / landowner / worker) — orienting gloss for the rent surface."),
  R(3, "VI — Ground-Rent", 38, "General Remarks on Differential Rent", ["ΔM"], "new", ["T2"], "finance",
    "Differential-rent panel; capitalised rent = rent÷interest = the price of land."),
  R(3, "VI — Ground-Rent", 39, "Differential Rent I", ["ΔM"], "new", ["T2"], "finance",
    "DR I table — soil grades (fertility + location) at equal capital; the worst soil regulates price."),
  R(3, "VI — Ground-Rent", 40, "Differential Rent II", ["ΔM"], "new", ["T2"], "finance",
    "DR II — successive capital doses on the same land; differential productivity of doses."),
  R(3, "VI — Ground-Rent", 41, "DR II — Constant Price of Production", ["ΔM"], "new", ["T2"], "finance",
    "The capital-productivity case; intensive vs extensive under a constant regulating price."),
  R(3, "VI — Ground-Rent", 42, "DR II — Falling Price of Production", ["ΔM"], "new", ["T2"], "finance",
    "Soil-exclusion event; rent outcomes under a falling regulating price."),
  R(3, "VI — Ground-Rent", 43, "DR II — Rising Price of Production", ["ΔM"], "new", ["T2"], "finance",
    "Rent sequences under a rising price (Engels's worked tables)."),
  R(3, "VI — Ground-Rent", 44, "Differential Rent Also on the Worst Soil", ["ΔM"], "new", ["T2"], "finance",
    "Rent emerging on the no-rent soil via improvement capital and regulating-price shifts."),
  R(3, "VI — Ground-Rent", 45, "Absolute Ground-Rent", ["ΔM"], "done", ["T2","T3"], "finance",
    "★ Absolute rent from agriculture's lower organic composition (value > price of production); the monopoly barrier (gloss)."),
  R(3, "VI — Ground-Rent", 46, "Building-Site Rent, Rent in Mining, the Price of Land", ["ΔM"], "new", ["T2","T3"], "finance",
    "Land-price panel (capitalised rent); monopoly-price rent (gloss)."),
  R(3, "VI — Ground-Rent", 47, "Genesis of Capitalist Ground-Rent", ["historical","ΔM"], "new", ["T4"], "finance",
    "Labour-rent → rent-in-kind → money-rent → capitalist rent on the genesis timeline."),

  R(3, "VII — Revenues & Their Sources", 48, "The Trinity Formula", ["whole","ΔM"], "done", ["T1","T2","T3"], "finance",
    "★★ The culminating mystification (Capital→interest, Land→rent, Labour→wages); the three revenue streams (T1/T2); realm of freedom (T3)."),
  R(3, "VII — Revenues & Their Sources", 49, "Concerning the Analysis of the Process of Production", ["whole","P","ΔM"], "new", ["T2","T3"], "finance",
    "Annual value composition c+v+s panel; Smith's adding-up dogma demolished (it omits c) — gloss."),
  R(3, "VII — Revenues & Their Sources", 50, "Illusions Created by Competition", ["whole","ΔM"], "done", ["T3","T2"], "finance",
    "★ A gallery of fetishes (profit-from-exchange · interest-from-a-thing · wages-pay-labour); the cost-price fetish panel."),
  R(3, "VII — Revenues & Their Sources", 51, "Distribution Relations and Production Relations", ["whole"], "new", ["T3"], "finance",
    "Distribution = the other side of production (both transitory); the productive-forces/social-form contradiction (gloss)."),
  R(3, "VII — Revenues & Their Sources", 52, "Classes", ["historical","whole"], "done", ["T2","T3"], "finance",
    "★ The three great classes by relation to the means of production (T2); the polarising tendency; the unfinished ending (T3/T4)."),
];
