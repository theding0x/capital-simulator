# Atlas Observatory — Total Representation
## Design handoff for Claude Design

**Date:** 2026-06-06
**Branch:** `feature/atlas-observatory`
**Prepared by:** Claude (Opus) — brainstorming/handoff pass
**Audience:** Claude Design (returns a complete design plan)

> This is a **handoff brief**, not an implementation plan. Its job is to give
> Claude Design everything needed to return a *complete design plan* for an
> expanded Atlas Observatory: information architecture, navigation model,
> representational vocabulary, a concept→representation mapping for **every one of
> the 106 chapters of *Capital***, the snapshot-schema additions implied, the
> motion/visual language, accessibility, and a phased build order.
>
> **Relationship to the prior handoff.** `2026-06-05-atlas-design-handoff.md` was
> an *aesthetic redesign* of the existing 3-slice Observatory ("make it beautiful,
> rethink the concept"). **This brief is a scope expansion**: from a dozen
> represented concepts to *all* of them. The two are complementary — the visual
> language from that pass should carry into this one. Read both.

---

## 1. Premise & deliverable

The Atlas Observatory today is a single live visualization: an orrery of
circulating industrial capitals (surface) above the "hidden abode of production"
(the general law of accumulation). It represents perhaps a dozen of *Capital*'s
concepts.

**The goal:** the Observatory graduates into the **primary navigable interface to
the whole of *Capital*** — every major *and* minor concept and argument across
all three volumes (all 106 chapters, every registry row already `done`) gets
**some** form of representation. The standalone per-chapter "Chapters" page is
absorbed into this world (its fate is an open question for Claude Design — see §8).

**What Claude Design returns:** a *complete* design plan for the whole
Observatory at once, covering:

1. Information architecture & navigation model (the Circuit × Depth spine, §3).
2. The representational vocabulary — the visual/interaction language for each
   representation **tier** (§4).
3. A **concept → representation mapping** for every chapter and its key
   sub-chapter concepts (the inventory in §6 is the backbone; each entry must be
   assigned a tier and, where live, a data source).
4. The **snapshot schema additions** implied by the live representations (§5).
5. Motion & visual language, consistent with the existing orrery/abode aesthetic
   and the `capital-simulator-design` skill.
6. Accessibility (reduced-motion parity, keyboard nav, semantic structure).
7. A phased build order (what ships first, what each phase unlocks).

---

## 2. What exists today (the starting point)

**Frontend — `web/src/atlas/`:**

| File | Role |
|---|---|
| `Atlas.tsx` | The world shell: topbar, rail, `<main class="stage">` with a vertical `world` (surface zone → gate → abode zone), footer console. Animated scroll descent (inOutCubic). Hash nav (`#/`, `#/chapters`). |
| `surface.ts` | Canvas controller (`AtlasSurface`). Renders the orrery: each capital a ring of three coexisting arcs — **M** money / **P** production / **C′** commodity — spiraling outward as it accumulates. Centre = average rate of profit `p̄′` (centre of gravity). |
| `Abode.tsx` | The hidden abode: working-day hero (necessary vs surplus, s/v), demoted stat tiles (Σv, Σs, organic composition c/v), the immiseration chart, the reserve army, the levers. |
| `ImmiserationChart.tsx` | The general law in motion — time-series of wage / exploitation / reserve army / organic composition. |
| `ReserveArmy.tsx` | The industrial reserve army reservoir. |
| `Levers.tsx` | Three live perturbation sliders: working day (s′), wage (value of labour-power), accumulation rate α. |
| `VitalSigns.tsx` | Aggregate vitals in the rail (avg rate of profit, count). |
| `TickHeartbeat.tsx` | Footer transport: run/pause, speed, reduced-motion. |
| `useSnapshot.ts` | Polls `GET /v1/observatory/snapshot?advance=N` with `X-Atlas-Session`. |
| `session.ts` / `prefs.ts` | Per-session id; currency/speed/reduced prefs in localStorage. |
| `animation.ts` | Easing + format helpers (`formatPence`, `formatBP`, `formatMinutes`). |
| `atlas.css` | The full visual system (~17KB): gate, descent, gilded palette. |

**Backend — `services/simulation-engine/`:**

- `internal/simulation/abode.go` — `AbodeState` (total social capital: c/v, value
  of labour-power, worker supply, law parameters) + `AdvanceGeneralLaw` (one
  period of the Ch.25 general law) + `Readout` (instantaneous class state) +
  `ApplyLevers`.
- `internal/observatory/` — per-session **ephemeral, in-memory** runs. No MySQL
  writes. Seeded once at boot; advanced on poll.
- `internal/transport/httpapi/observatory_handler.go` — `GET
  /v1/observatory/snapshot`. Returns `{ tick, running, interval_ms, capitals[],
  aggregate, abode }`.

**Concepts currently represented (the dozen):** the circuit M—C…P…C′—M′ (orrery
rings), the general formula (Ch.4), the Ch.6 gate, working day & rate of
surplus-value (Ch.7–10), constant/variable capital & organic composition (Ch.8),
the general law / reserve army / immiseration (Ch.25), average rate of profit
(Vol.III Ch.9), turnover number (Vol.II Ch.7) on each capital.

---

## 3. The spine — Circuit × Depth

The world is organized on two axes that already exist in the codebase
(`web/src/chapters/registry.ts`: every chapter carries `volume` and
`circuitNode[]`).

```
        M    M-C    P    C'   C-M'   M'   ΔM        ← circuit moments (horizontal)
      ┌----+------+----+----+------+----+----┐
Vol I |  money|buy Lp|PROD|comm| sell |ret |s/v |   surface: the orrery
 prod |       | +Mp  |●●●●|    |      |    |    |   (many capitals in motion)
      ├----+------+----+----+------+----+----┤
Vol II| metamorphoses · turnover · reproduction  |  mid stratum (circulation)
 circ |       |      |    |    |      |    |    |
      ├----+------+----+----+------+----+----┤
VolIII| profit → avg p̄' → TRPF → merchant →      |  deep stratum (totality)
total | interest · rent · trinity · classes      |  (the abode and beyond)
      └------------------------------------------┘
  descend = go deeper in explanation;  pan = move along the circuit
```

- **Horizontal axis — circuit moments:** `M → M-C → P → C′ → C-M′ → M′ → ΔM`,
  plus two cross-cutting registers from the registry's `CircuitNode` union:
  `whole` (chapters that read the entire motion) and `historical` (the conditions
  that produced the circuit — primitive accumulation, genesis, secular
  tendencies). The existing orrery already renders M/P/C′ arcs per capital.
- **Vertical axis — depth of explanation = volume:** Vol I (production, where
  surplus is pumped in **P**) → Vol II (circulation & turnover, how value moves
  through **M-C** and **C′-M′**) → Vol III (totality & distribution, how
  surplus-value distributes as profit/rent/interest). This is Marx's own framing
  (see `CLAUDE.md` → "Volumes"): each volume is a *depth of explanation on the
  same circuit*, not a separate exhibit.
- **The existing descent generalizes:** the current surface→gate→abode fall is
  the first leg. Descending continues *through* the abode into Vol II and Vol III
  strata. The Ch.6 "No admittance except on business" gate stays as the threshold
  between the noisy sphere of circulation (surface) and production (abode).
- **Navigation:** pan along the circuit; descend to go deeper. Any **cell**
  (circuit-moment × volume) opens to its Part(s) → chapters → concepts. The
  `circuitNode` + `volume` + `part` fields in the registry are the addressing
  scheme — Claude Design should treat them as the routing keys.

---

## 4. Representation tiers (the "tiered" decision, made concrete)

Every chapter, and every key sub-chapter concept, maps to **at least one** tier.
The quantitative core gets live tiers; qualitative/historical/polemical material
gets annotation/historical tiers. Claude Design defines the visual+interaction
language for each tier and assigns every inventory entry (§6) to one.

- **T1 — Live instrument.** Data-driven, animated, reads the live snapshot.
  Continuously in motion or responsive to the levers. *Examples already built:*
  working day, p̄′, reserve army, immiseration series, turnover number.
- **T2 — Live-derived panel.** Computed or queried-on-demand values surfaced when
  a cell/concept is opened — tables, gauges, small multiples fed by an endpoint
  but not necessarily animating every tick. *Examples:* prices of production,
  differential-rent tables, fixed-capital wear & sinking funds, value-form ladder.
- **T3 — Annotation / marker / gloss.** Qualitative concepts pinned onto the
  relevant circuit-cell as a marker that opens a short prose gloss with a Marx
  citation. *Examples:* commodity fetishism, the fetish character of
  interest-bearing capital (M…M′), Illusions Created by Competition, the trinity
  formula's mystification.
- **T4 — Historical stratum.** The "conditions that made the circuit" — a layer
  behind/beneath the living world. Primitive accumulation, expropriation, bloody
  legislation, genesis chapters, colonisation, pre-capitalist relationships,
  historical tendency. Narrative/temporal rather than live-numeric.

**Rule of thumb:** numeric/lawful → T1 or T2; critique/qualitative → T3;
historical/genetic/secular → T4. A chapter may carry several (e.g. Ch.25 is T1
for the general law *and* T4 for its historical illustration).

---

## 5. Data sources & backend

**Preserve the ephemeral, zero-persistence, per-session ethos.** The Observatory
must not write MySQL; live data is aggregated into the per-session in-memory run
and surfaced through (an extended) `GET /v1/observatory/snapshot`. Extend the
snapshot to pull from the services that already model each chapter:

| Service | Port | Models (relevant to live tiers) |
|---|---|---|
| commodity-service | 8081 | Vol I — commodity, value, **value-forms** (simple→money), fetishism, c+v decomposition, SNLT / productivity |
| agent-service | 8082 | Vol I — workers, capitalists, labour-process, **wages** (time/piece/national), cooperation, manufacture |
| market-service | 8083 | Vol I — exchange, money, prices; Vol II — **circulation phases**, costs of circulation |
| simulation-engine | 8084 | Vol I — production tick, **machinery**, **reproduction**, **accumulation**, primitive accumulation; Vol II — **turnover** & reproduction schemes; **hosts the observatory** |
| finance-service | 8085 | Vol III — **profit**, **average rate of profit**, **prices of production**, **rent**, **interest**, credit, **fictitious capital**, the **trinity formula** |

Claude Design should specify, per live (T1/T2) representation, which service +
endpoint feeds it and what new fields the snapshot DTO must carry. The current
snapshot DTO (`observatory_handler.go`) is the template: `capitals[]`,
`aggregate`, `abode` — expect new sibling sections per volume/stratum (e.g.
`circulation`, `distribution`).

---

## 6. The concept inventory (the backbone)

The complete 106-chapter census follows, grouped by **volume → Part →
chapter**, with each chapter's `circuitNode` and its **major + minor concepts and
arguments** (mined from the red-vault specs' `Concepts → types` tables and the
chapter arguments). Claude Design must assign every concept a **tier** (§4) and,
where T1/T2, a **data source** (§5).

Legend: ★ = already represented in today's Observatory · ◐ = partially · (blank) =
net-new. Suggested tiers (T1–T4) are starting points, not final.

<!-- INVENTORY:BEGIN -->

## Volume I — The Process of Production of Capital

### Part I — Commodities & Money

- **Ch. 1 — The Commodity** `[whole]` ◐ — use-value vs exchange-value; **value** & socially-necessary labour-time (SNLT); abstract vs concrete labour; the **value-form ladder** (simple → expanded → general → money-form); **commodity fetishism** (social relations of labour appearing as relations between things); the inverse productivity↔value law. → *value-form ladder T2; fetishism T3; value magnitude T1*
- **Ch. 2 — Exchange** `[M-C, C-M-prime]` — the exchange act; commodity-owners as personifications; direct barter ratio; the **social act** that sets apart a universal equivalent; money **crystallising** out of exchange; price as value-expression; C-M-C sale/purchase legs. → *T2 + T3 (money's genesis)*
- **Ch. 3 — Money, or the Circulation of Commodities** `[M, M-prime]` ◐ — C-M-C circuit; the **functions of money** — measure of value, medium of circulation (M = ΣP ÷ velocity), **hoarding**, means of payment (creditor↔debtor), world money. → *money-required gauge T2; functions-of-money T3*

### Part II — The Transformation of Money into Capital

- **Ch. 4 — The General Formula for Capital** `[whole]` ★ — C-M-C vs **M-C-M′**; surplus-value ΔM = M′ − M; capital as **self-expanding value**; the miser vs the rational capitalist; the abridged M-M′. → *★ orrery rings; ΔM T1*
- **Ch. 5 — Contradictions in the General Formula** `[whole]` — exchange of equivalents / non-equivalents both **conserve total social value**; circulation begets no value; zero-sum redistribution; merchants' & usurer's capital as the still-unexplained ΔM; the contradiction (surplus arises neither in nor outside circulation). → *value-conservation demo T2; the contradiction T3*
- **Ch. 6 — The Sale and Purchase of Labour-Power** `[M-C]` ★ — labour-power as a **commodity**; its value = SNLT of the **subsistence basket**; the **double-free labourer**; finite contract; subsistence/minimum floor; sale is an exchange of equivalents (surplus deferred to *use*); reproduction cost. → *★ the gate; subsistence basket T2; double-freedom T3*

### Part III — The Production of Absolute Surplus-Value

- **Ch. 7 — The Labour-Process and Valorization** `[P]` ★ — the labour-process (purposeful activity, subject of labour, instruments); means of production; **constant-capital value-transfer**; value added by living labour; **necessary vs surplus labour**; the working-day partition; valorization; surplus-value as **unpaid labour**; product value = c + v + s. → *★ working-day T1; c+v+s composition T2*
- **Ch. 8 — Constant Capital and Variable Capital** `[P]` ★ — **c (value transferred) vs v (value created)**; depreciation/wear of instruments; raw vs auxiliary material; product-value decomposition c + v + s; capital composition c:v. → *★ organic composition; c+v+s T2*
- **Ch. 9 — The Rate of Surplus-Value** `[P, delta-M]` ★ — **rate s/v = degree of exploitation**; value-product (v + s, c excluded); capital advanced C = c+v vs expanded C′; necessary/surplus labour; surplus-produce fraction; the s/C-vs-s/v distinction. → *★ rate of exploitation T1*
- **Ch. 10 — The Working-Day** `[P]` ★ — the **contested length** (necessary vs surplus segments); statutory limits & the **Factory Acts** (1833 / 1847 Ten Hours' Bill / 1850); the **relay/shift system** (day/night); overwork ("nibbling"); the physical 24h maximum; the corvée comparison. → *★ working-day T1; Factory Acts T4; relay system T2*

### Part III/IV bridge — Rate & Mass, then Relative Surplus-Value

- **Ch. 11 — The Rate and Mass of Surplus-Value** `[P, delta-M]` — mass **S = (s/v)·V = P·(a′/a)·n**; number of labourers; the **compensation law** (rate ↑ ↔ workers ↓ holds S constant, within the 24h limit); minimum capital to *become* a capitalist. → *S mass T1/T2; compensation law T2*

### Part IV — The Production of Relative Surplus-Value

- **Ch. 12 — The Concept of Relative Surplus-Value** `[P]` — **absolute vs relative** surplus-value; shortening necessary labour by raising productivity; value of labour-power falling with cheaper necessaries; individual vs social value; **extra surplus-value** (the temporary innovator's bonus); the inverse value↔productivity law. → *absolute/relative split T2; extra surplus-value T2*
- **Ch. 13 — Co-operation** `[P]` — simple co-operation; the **collective working-day**; scale; **collective productive power > the sum of isolated labours** (the co-operation bonus); average social labour (deviations cancel); supervision/directing authority; minimum capital; the collective power *appearing as a property of capital*. → *collective power T2; appears-as-capital T3*
- **Ch. 14 — Division of Labour and Manufacture** `[P]` — manufacture (heterogeneous vs serial); the **detail labourer**; the **collective labourer**; partial products (not yet commodities); the **hierarchy of labour-powers** & deskilling; the proportionality law (group sizes); rising minimum capital; two-fold origin (combination/splitting); tool specialisation; **workshop's a-priori plan vs the market's a-posteriori social division of labour**. → *manufacture/proportionality T2; deskilling T3; the two DOLs T3*
- **Ch. 15 — Machinery and Modern Industry** `[P]` ◐ — the machine (motor + transmission + tool); **wear & tear value-transfer**; **moral depreciation**; productive power; the factory & prime mover; **labour displaced → the reserve army**; intensification of labour; the **composition shift** (c rises, v falls); the industrial cycle (prosperity / over-production / crisis / stagnation). → *◐ machinery feeds the abode's organic composition; machine wear T2; labour-displaced → reserve army ★; industrial cycle T2/T4*

### Part V — Absolute and Relative Surplus-Value

- **Ch. 16 — Absolute and Relative Surplus-Value** `[P, delta-M]` ★ — the unity & distinction of absolute/relative; **formal vs real subsumption** of labour under capital; productive labour (in the capitalist sense); the **rate of profit vs rate of surplus-value** (Mill critique). → *★ working-day split; formal/real subsumption T3; profit-vs-surplus rates T2*
- **Ch. 17 — Changes of Magnitude in the Price of Labour-Power and Surplus-Value** `[P, delta-M]` ◐ — the three variable factors (working-day **length**, **intensity**, **productiveness**); the laws relating value of labour-power and surplus-value (constant daily value; inverse relation); **reserve-army wage compression**. → *magnitude scenarios T2; ◐ reserve-army compression feeds the abode wage*
- **Ch. 18 — Various Formulae for the Rate of Surplus-Value** `[delta-M]` — Formula I (s/v, unbounded), II (s ÷ working-day, always < 100%), III (unpaid ÷ paid); **how Formula II conceals exploitation**. → *formula comparison T2; the exploitation-hiding formula T3*

### Part VI — Wages

- **Ch. 19 — The Transformation of the Value of Labour-Power into Wages** `[M-C, delta-M]` — the **wage-form**; wages *appear* to pay for all labour (the ideological inversion); the true paid/unpaid decomposition beneath the form; hourly wage / "price of labour"; real wage falling while the form conceals it. → *paid/unpaid T1; the wage-illusion T3*
- **Ch. 20 — Time-Wages** `[M-C]` — daily value of labour-power → the **price of the working-hour** (exact fraction); nominal vs real wage; lengthening the day **lowers the hourly price** even at constant pay; overtime as continued unpaid extraction. → *hourly-price T2; longer-day-lower-price T3*
- **Ch. 21 — Piece-Wages** `[M-C]` — piece-price = daily wage ÷ normal output; **piece value > piece price**; quality control built into the wage-form (rejected pieces unpaid); the **sweating / sub-contract** system (middlemen keep the spread); dynamic piece-price reduction as productivity rises. → *piece-price/value T2; sweating T3*
- **Ch. 22 — National Differences in Wages** `[M-C]` — national **intensity** of labour; the standardised wage (reduced to a uniform day); **relative labour price** (the England paradox — highest nominal, lowest *relative* price); spindle ratios as a productivity proxy; the world-market law of value. → *wage-comparison table T2; the England paradox T3*

### Part VII — The Accumulation of Capital

- **Ch. 23 — Simple Reproduction** `[whole]` — capital stock c/v; the surplus-value fund (revenue vs accumulated); the **repayment period** (the original capital dissolves into consumed surplus, C ÷ S); variable capital *is the worker's own product* advanced back; the **reproduction of the class relation**. → *cycle series T2; class-relation reproduction T3*
- **Ch. 24 — The Transformation of Surplus-Value into Capital** `[whole, delta-M]` ★ — **accumulation**; the accumulation rate α (consume vs reinvest); additional capital split per composition; revenue; **compound growth** ("Abraham begat Isaac"); accumulation as self-expanding spiral; abstinence-theory critique. → *★ the α lever; compound-growth spiral T1*
- **Ch. 25 — The General Law of Capitalist Accumulation** `[whole, delta-M, historical]` ★★ — value vs technical vs **organic composition**; the **rising organic composition**; **concentration & centralisation**; the **industrial reserve army** and its three strata (floating / latent / stagnant); labour-demand from variable capital; **the general law** (accumulation of wealth = accumulation of misery); the empirical illustrations (England 1846–66, Ireland). → *★★ the abode's core — immiseration, reserve army, organic composition; strata T2; concentration/centralisation T2; the illustrations T4*

### Part VIII — So-Called Primitive Accumulation `[historical]` — the T4 stratum

This whole Part is the **historical stratum (T4)** — the conditions that produced the circuit. Suggest a single continuous "genesis" layer *beneath* the world, walkable as a timeline.

- **Ch. 26 — The Secret of Primitive Accumulation** `[historical]` — the "original sin" / circularity riddle; the **separation of the producer from the means of production**; the two kinds of commodity-possessor (money-owners + free labourers); **capital as a social relation, not a thing**. → *T4; "capital is a social relation" gloss T3*
- **Ch. 27 — Expropriation of the Agricultural Population** `[historical]` — **enclosure** of the commons; the Highland Clearances (Duchess of Sutherland); dispossession → proletariat; sheep displacing men. → *T4 expropriation timeline*
- **Ch. 28 — Bloody Legislation against the Expropriated** `[historical]` — **wage statutes** (Statute of Labourers 1349 — a maximum, never a minimum); **vagrancy laws** (whipping, branding); coercing the dispossessed into wage-labour; the shift from legal coercion to the "dull compulsion of economic relations." → *T4*
- **Ch. 29 — Genesis of the Capitalist Farmer** `[historical]` — the tenant-form progression (bailiff → métayer → capitalist-farmer); long leases + money depreciation enriching the farmer; the farming surplus. → *T4*
- **Ch. 30 — Reaction of the Agricultural Revolution on Industry; the Home Market** `[historical]` — destruction of rural domestic industry; proletarianisation **creating the home market**; subsistence goods becoming commodities; manufacture réunie vs séparée. → *T4; market-formation T2*
- **Ch. 31 — Genesis of the Industrial Capitalist** `[historical]` — the **sources of capital**: colonial plunder, the slave trade, the national debt, taxation, protectionism, usury/commerce; capital "dripping from head to foot with blood and dirt"; the Bank of England. → *T4 sources breakdown*
- **Ch. 32 — The Historical Tendency of Capitalist Accumulation** `[historical]` — petty property → capitalist property → **socialised property**; the **negation of the negation**; the centralisation spiral; "the expropriators are expropriated"; the knell of capitalist private property. → *T4; dialectical-stages T3*
- **Ch. 33 — The Modern Theory of Colonisation** `[historical]` — the colonies reveal **capital as a social relation** (Peel's fiasco — left without a servant); free land undermining wage-dependence; Wakefield's "sufficient price" of land / systematic colonisation. → *T4; Peel's fiasco T3*

## Volume II — The Process of Circulation of Capital

*Depth stratum 2. The mid-layer of the world: how value moves through M-C and C′-M′ between two functionings of P.*

### Part I — The Metamorphoses of Capital and Their Circuits

- **Ch. 1 — The Circuit of Money-Capital** `[whole, M, M-prime]` ★ — **M—C(L+MP) … P … C′—M′**; the six moments; the purchase phase split into **M—L** and **M—MP** (the proportional law); the productive state; commodity-capital C′ = C + c; realisation; surplus arises only in P; form-change preserves magnitude; partial/failed realisation. → *★ the orrery circuit; six-moment timeline T1*
- **Ch. 2 — The Circuit of Productive Capital** `[whole, P]` — **P … C′—M′—C … P** (production read as *reproduction*); simple / extended / mixed reproduction; the revenue circuit c—m—c (*outside* the capital circuit); **latent / hoarded money-capital** + the minimum capitalisation increment (£1/spindle); the capitalisation step; reserve fund & draws. → *reproduction lens T2; latent-capital gauge T2*
- **Ch. 3 — The Circuit of Commodity-Capital** `[whole, C-prime, C-M-prime]` — **C′—M′—C … P … C′**; opens already carrying C + c; **aliquot decomposition** (every pound = c + v + s); successive partial sales; the **branching** of realised money into revenue (m) vs reproduction-fund (M); the **external-C presupposition** (C is another capital's C′ — the social intertwining); material adequacy of accumulation; the social-capital lens. → *aliquot decomposition T2; intertwining T3; social-capital lens → reproduction schemes*
- **Ch. 4 — The Three Formulas of the Circuit** `[whole]` ★★ — **industrial capital as the unity of M…M′, P…P, C′…C′, all three coexisting** (the *Nebeneinander* — exactly the orrery's three arcs); capital parts simultaneously at every stage; the stage distribution; **stagnation propagation** (a block at one stage jams the others — the germ of crisis); value-revolution (money set free / tied up); the metamorphosis interlock (one's M—C is another's C—M); supply structurally exceeds demand by the surplus-value; natural / money / credit economy; the sinking fund. → *★★ the orrery's three coexisting arcs; stage-distribution T1; stagnation toy T1; supply-demand gap T2*
- **Ch. 5 — The Time of Circulation** `[M-C, C-M-prime]` — total circuit time = **production time + circulation time**; production time is broader than labour-time (natural processes; latent productive capital held in readiness); circulation time = selling time (C—M) + buying time (M—C); the perishability window; spatial separation of markets; the tendency toward zero circulation time. → *time decomposition T2; feeds orrery pacing*
- **Ch. 6 — The Costs of Circulation** `[M-C, C-M-prime]` — **pure circulation costs** (value-preserving *faux frais*) vs value-adding costs; circulation agents (buying/selling/book-keeping labour); money as faux frais (gold tied up in circulation); commodity-supply/stock & storage costs; **transport as genuinely value-adding**; the transport tariff. → *faux-frais T2; transport-adds-value T3*

### Part II — The Turnover of Capital

- **Ch. 7 — The Turnover Time and the Number of Turnovers** `[whole]` ★ — turnover = the periodically-repeated circuit; the **number of turnovers n = T ÷ t** (the year as unit); the Form I vs Form II lens. → *★ turnover number already on each orrery capital*
- **Ch. 8 — Fixed Capital and Circulating Capital** `[whole]` ◐ — the **fixed/circulating distinction** within productive capital (orthogonal to constant/variable); components by role; fixed-capital items (machine/building/beast of toil); **wear & tear**; causes of wear (use / natural forces / moral depreciation); the **sinking fund**; repairs (running vs capitalised); subcomponents with their own lifetimes; the circulating-capital cycle; sinking-fund reinvestment (intensive/extensive). → *fixed/circulating T2; ◐ wear feeds the composition; sinking fund T2*
- **Ch. 9 — The Aggregate Turnover of Advanced Capital** `[whole]` — aggregate turnover of a mixed (fixed + circulating) capital; per-component contribution; the multi-year **lifetime cycle of fixed-capital reproduction** as the material basis of the **crisis-cycle**; mean term of turnover. → *aggregate turnover T2; fixed-capital crisis cycle T2/T4*
- **Ch. 10 — Theories of Fixed & Circulating Capital: Physiocrats & Adam Smith** `[whole, historical]` — Quesnay's *avances primitives / annuelles*; **Smith's conflation** of fixed/circulating with constant/variable; his misclassification of money- and commodity-capital; revenue-flow lodged inside capital-flow. → *doctrine-critique T3*
- **Ch. 11 — Theories of Fixed & Circulating Capital: Ricardo** `[whole, historical]` — Ricardo's collapse of fixed/circulating onto mere durability; conflation with constant/variable; differential fixed-capital intensity explaining divergence from labour-value pricing; his missing aggregate-turnover and value-revolution concepts. → *doctrine-critique T3*
- **Ch. 12 — The Working Period** `[P]` — the **working period** (connected working-days for one finished product); discrete vs connected; the advance multiplier (circulating capital tied up); interruption with deterioration; shortening; the natural floor; speculative-build / credit financing. → *working period T2*
- **Ch. 13 — The Time of Production** `[P]` — **production time** as the umbrella over labour-time; natural processes; the production-time/labour-time gap; latent productive capital while nature works; branches where natural-process time dominates (wine, forestry, agriculture). → *production time T2*
- **Ch. 14 — The Time of Circulation** `[M-C, C-M-prime]` — **circulation time** with market-distance and communication-lag; the **annual-rate-of-surplus-value penalty** from long circulation; capital set free / tied up by speed change; circulation-speed improvements (railways, telegraph). → *circulation time T2; annual-rate penalty T2*
- **Ch. 15 — The Effects of a Change of Prices** `[whole]` — a **price-change mid-turnover**; the affected element (MP / labour-power / output); revalued capital; realised-value delta; inventory revaluation; surplus realised against the *original* advance; the three cases on a fall and the three on a rise; **speculation**; the compound price-×-speed interaction. → *price-revolution T2; speculation T3*
- **Ch. 16 — The Turnover of Variable Capital** `[P, delta-M]` ◐ — the **annual rate of surplus-value** (vs the per-turnover rate); the capitalist-A-vs-B contrast; "the same advance produces n× the surplus per year"; the asymmetry (the working period bounds n but not s/v). → *annual rate of surplus-value T2; ◐ ties to the orrery turnover*

### Part III — The Reproduction and Circulation of the Total Social Capital

- **Ch. 17 — The Circulation of Surplus-Value** `[delta-M, whole]` — aggregate surplus-value circulation; **where the money to realise surplus comes from** (the realisation puzzle); the social-capital aggregate snapshot; the question Ch. 20–21 answers. → *the realisation-money puzzle T3*
- **Ch. 18 — The Role of Money-Capital in Reproduction** `[whole]` — aggregate money-supply apportionment; per-department money-reserves; **the two Departments (I & II) enter here**; the circulating money-mass (M·V = P·Q at the social level); money velocity; the wage-rotation fund; inter-department settlement. → *money supply T2*
- **Ch. 19 — Former Presentations of the Subject** `[whole, historical]` — **Quesnay's *Tableau Économique* (1758)** and its errors; **Smith's dogma** that revenue resolves entirely into wages + profit + rent (the missing constant-capital replacement, the infinite regress of c into past labour); Ricardo's improvements; the Tooke–Wilson debate. → *doctrine-history T3; Smith's dogma T3*
- **Ch. 20 — Simple Reproduction** `[whole]` ★ — the **two Departments** (I = means of production, II = means of consumption); departmental capital c + v + s; **the three great exchanges** (the keystone: I(v+s) ↔ II(c)); the inter-department flow; the closed money loop; the reproduction tick; Dept-I c replaced internally. → *★ the two-department reproduction scheme — a prime new T1 viz*
- **Ch. 21 — Accumulation and Reproduction on an Extended Scale** `[whole, delta-M]` ★ — the **extended-reproduction scheme** (the keystone aggregate); the accumulation rate; the reinvestment step; year-over-year growth; the **composition shift** (c grows relative to v over cycles); the money-supply requirement; the Department-I-leads asymmetry. → *extended-reproduction growth T1; ◐ feeds accumulation & the general law*

## Volume III — The Process of Capitalist Production as a Whole

*Depth stratum 3. The deepest layer: how surplus-value distributes as profit, rent, interest — the totality, where the abode already lives. This is where the bulk of the net-new representation lands.*

### Part I — The Transformation of Surplus-Value into Profit

- **Ch. 1 — Cost-Price and Profit** `[delta-M]` ◐ — **cost-price k = c + v**; **profit as the mystified form of surplus-value** (C = k + s rewritten as k + p, concealing the origin); fixed vs circulating share entering the cost-price; selling below value yet still pocketing a profit. → *cost-price/profit-form T2; the mystification T3*
- **Ch. 2 — The Rate of Profit** `[delta-M]` ★ — **p′ = s ÷ (c+v) = s ÷ C**; p′ vs s′ (p′ always < s′ when c > 0); the mystification degree (the gap between s′ and p′). → *★ p′ is the orrery's centre of gravity; p′-vs-s′ T2; mystification T3*
- **Ch. 3 — The Relation of the Rate of Profit to the Rate of Surplus-Value** `[delta-M]` — **p′ = s′·(v/C)**; the composition ratio v/C; the variation cases (s′ and v/C constant or varying); the profit-rate comparison. → *interactive p′ = s′·(v/C) T2*
- **Ch. 4 — The Effect of the Turnover on the Rate of Profit** `[M-C, C-M-prime, delta-M]` ◐ — the **annual rate of profit p′ = s′·n·(v/C)**; the annual rate of surplus-value S′ = s′·n; the turnover count n. → *annual p′ T2; ◐ turnover*
- **Ch. 5 — Economy in the Employment of Constant Capital** `[P, delta-M]` — **economy in c raising p′** without touching s or v; shared production conditions; waste reduction; the workday-extension spreading fixed-c; improved machinery. → *economy-in-c T2*
- **Ch. 6 — The Effect of Price Fluctuation** `[M-C, delta-M]` — raw-material price changes rippling into p′; **capital release** (a price fall frees tied-up capital) vs **capital tie-up** (a price rise); supply regulation. → *price-fluctuation T2; release/tie-up T2*
- **Ch. 7 — Supplementary Remarks** `[P, delta-M]` — profit's **surface appearance** (it *seems* to spring from business acumen, not exploitation); the organic-composition effect (same s, different c/v → different p′); the Rodbertus refutation (a nominal magnitude change leaves the rate unchanged); the Part I summary. → *surface-appearance T3; composition effect T2*

### Part II — Conversion of Profit into Average Profit

- **Ch. 8 — Different Compositions of Capitals in Different Branches** `[delta-M, whole]` ◐ — **production spheres/branches** with different organic compositions; technical vs organic composition; each sphere's individual p′ diverging *before* equalisation; the European-vs-Asian contrast (high comp/high s′ vs low comp/low s′). → *multi-sphere p′ scatter T2*
- **Ch. 9 — Formation of a General Rate of Profit; Prices of Production** `[delta-M, whole]` ★ — **the general/average rate p̄′ = ΣS ÷ ΣC**; average profit; **price of production = k + average profit**; the value→price transformation; value-price deviation (high-comp sells below value, low-comp above); the two conservation laws (Σprices = Σvalues, Σprofit = Σsurplus). → *★ p̄′ is the orrery's centre of gravity; prices of production T2; the transformation T2/T3*
- **Ch. 10 — Equalisation of the General Rate through Competition** `[C-M-prime, delta-M]` ★ — **market value vs market price**; **surplus-profit** (a below-average individual value); **capital flow between spheres** (the equalisation mechanism); convergence to the general rate. → *★ capital flow IS the orrery's centre-of-gravity dynamics; market value T2; surplus-profit T2*
- **Ch. 11 — Effects of General Wage Fluctuations on Prices of Production** `[delta-M]` — a general wage rise/fall; the average-composition price unchanged; lower-comp price rises while higher-comp falls (and the reverse); the general rate moving inversely to wages. → *wage→price effect T2*
- **Ch. 12 — Supplementary Remarks** `[delta-M]` — the two causes of a price-of-production change (general-rate vs individual-value); capitalist **compensation grounds** (turnover/risk/skill = a *share in* aggregate surplus, not value creation); the Part II summary. → *compensation-grounds myth T3; summary T2*

### Part III — The Law of the Tendential Fall in the Rate of Profit

- **Ch. 13 — The Law As Such** `[delta-M, whole]` ★ — **rising organic composition + constant s′ → falling p′**; the mass of profit (which can rise *as* the rate falls — the **rate–mass contradiction**); absolute over-accumulation. → *★ TRPF is core to the totality; trajectory T1; rate-mass contradiction T1*
- **Ch. 14 — Counteracting Influences** `[delta-M]` — the **six forces** that retard the fall: intensified exploitation, depressed wages, cheaper constant capital, relative overpopulation, foreign trade, stock capital; the law is *retarded, not abolished*. → *counteracting-forces toggles T2*
- **Ch. 15 — Internal Contradictions of the Law** `[delta-M, whole]` ★ — falling-rate/rising-mass; over-production / under-consumption; over-accumulation / relative overpopulation; concentration & centralisation; **crisis as the violent resolution** (devalues constant capital, restores p′); the Part III summary. → *★ crisis restores p′ T1; concentration/centralisation T2; the contradictions T3*

### Part IV — Merchant's Capital

- **Ch. 16 — Commercial Capital** `[C-M-prime, delta-M]` — **merchant's capital** as a functionally independent form (M—C—M′) that mediates C′→M′ for industrial capital; it **creates no new value**, only appropriates a share of already-produced surplus; distinct from industrial and usurer's capital. → *merchant circuit T2; creates-no-value T3*
- **Ch. 17 — Commercial Profit** `[C-M-prime, delta-M]` — commercial profit as a **deduction from social surplus-value**; the adjusted general rate (counting commercial capital in the denominator lowers p̄′); the exploitation of commercial workers (their unpaid labour swells merchant profit without creating value); the wholesale–retail spread. → *adjusted general rate T2; commercial-worker exploitation T3*
- **Ch. 18 — The Turnover of Merchant's Capital; Prices** `[M, C-M-prime, delta-M]` — merchant turnover; the **per-unit markup inversely proportional to turnover**; annual merchant profit = the general rate on the capital advanced. → *turnover→markup T2*
- **Ch. 19 — Money-Dealing Capital** `[M, delta-M]` — the technical operations of money circulation (receipts, payments, currency exchange, safekeeping, bookkeeping); M—M′; creates no value; the cash reserve; participates in equalising the general rate. → *money-dealing T2*
- **Ch. 20 — Historical Facts about Merchant's Capital** `[historical, delta-M]` — merchant's capital as the **oldest independent form** (pre-capitalist M—C—M′); profit via unequal exchange (Venice, Genoa, the Dutch carrying trade); its **subordination** to industrial capital as capitalism matures (the inverse relation). → *T4 historical; subordination T2*

### Part V — Interest, Credit & Fictitious Capital

- **Ch. 21 — Interest-Bearing Capital** `[M, M-prime]` ★ — money loaned to return M + ΔM (interest); **loanable capital as a commodity *sui generis*** (its use-value = the capacity to produce profit); the money-capitalist vs the functioning capitalist; interest as a portion of average profit; the M—M′ circuit. → *interest/loan T2; the M—M′ form T3*
- **Ch. 22 — Division of Profit; the Rate of Interest** `[M, M-prime]` — the rate of interest has **no natural law** (set by supply/demand in the money market); maximum = average profit, minimum ≈ 0; the industrial cycle driving it; profit of enterprise = profit − interest. → *interest rate T2; cycle-driven rate T2*
- **Ch. 23 — Interest and Profit of Enterprise** `[M, M-prime]` — total profit splits **qualitatively** into interest (to the owner) and profit of enterprise (to the active capitalist); the ideological inversion (interest *seems* the fruit of ownership, profit of enterprise the wage of management); wages of superintendence. → *profit division T2; the inversion T3*
- **Ch. 24 — Externalisation of the Relations of Capital (Fetish Capital)** `[M, M-prime]` ★ — interest-bearing capital as the **most fetishised form**: M—M′, money breeding money, production vanished entirely (pairs with Ch. 1's commodity fetishism as the deepest mystification); compound-interest fantasies (Dr Price, Pitt's sinking fund); the qualitative accumulation limit set by living labour. → *★ the capital fetish T3; compound-interest absurdity T2*
- **Ch. 25 — Credit and Fictitious Capital** `[delta-M]` ★ — the credit system growing from money-as-means-of-payment; **bills of exchange** (commercial money); bank credit (concentrating loanable capital); **fictitious capital** (capitalised income streams — bonds = consumed capital, stocks = titles to future surplus); capitalisation = income ÷ rate; the 1847 crisis. → *fictitious capital / capitalisation T2; fictitious-vs-real T3*
- **Ch. 26 — Accumulation of Money-Capital; its Influence on the Interest Rate** `[M, M-prime, delta-M]` — accumulation of **loanable** money-capital vs **real** accumulation (they diverge through the cycle); idle industrial capital pooling post-crisis; the real-vs-money-accumulation gap. → *loanable-vs-real T2; cycle T2*
- **Ch. 27 — The Role of Credit in Capitalist Production** `[delta-M]` ★ — credit's double role; the **joint-stock company** (capital as directly social capital; ownership split from management); the **cooperative factory** (the antithesis transcended within capitalism); circulation-cost saving; credit as lever of accumulation *and* of swindle/crisis. → *joint-stock socialisation T3; cooperative factory T3*
- **Ch. 28 — Medium of Circulation and Capital (Tooke & Fullarton)** `[M, delta-M]` — the **currency-vs-capital distinction**; the reserve fund; the revenue-expenditure circuit vs the capital-transfer circuit. → *doctrine T3; currency-vs-capital T2*
- **Ch. 29 — Component Parts of Bank Capital** `[delta-M]` — the components of bank capital; fictitious-capital valuation; deposits; the **double-counting** of the same money as capital. → *bank-capital decomposition T2*
- **Ch. 30 — Money-Capital and Real Capital, I** `[M, delta-M]` — real-capital accumulation vs commercial credit; accumulation correspondence; credit limits; the return flow. → *T2*
- **Ch. 31 — Money-Capital and Real Capital, II** `[M, delta-M]` — floating capital; bill-brokers; rediscount; loan-capital velocity; circulation-reserve saving. → *T2*
- **Ch. 32 — Money-Capital and Real Capital, III** `[M, delta-M]` — revenue→loan-capital conversion; capital release; means-of-payment demand; the **rentier**; loan-capital overstatement. → *T2*
- **Ch. 33 — The Medium of Circulation in the Credit System** `[M, delta-M]` — currency velocity; **clearing-house settlement** (millions settled with little actual money); credit contraction; the effective money supply. → *clearing house / velocity T2*
- **Ch. 34 — The Currency Principle and the Bank Act of 1844** `[M, delta-M, historical]` — the Bank Act regime; the note-issue constraint; the Bank's two departments; the currency-principle critique; the gold drain. → *Bank-Act regime T2; the 1844 critique T3; T4*
- **Ch. 35 — Precious Metal and the Rate of Exchange** `[M]` — the gold reserve; the rate of exchange; the balance of payments; gold flow; the world-money function. → *exchange rates / gold flow T2*
- **Ch. 36 — Pre-Capitalist Relationships** `[historical]` — **usurer's capital**; pre-capitalist credit forms; usury's development stages and subordination; historical interest rates. → *T4 historical; usury T3*

### Part VI — Transformation of Surplus-Profit into Ground-Rent

- **Ch. 37 — Introduction** `[delta-M]` — the land parcel; the **triad** of agricultural capitalist (farmer) / landowner / worker; ground-rent; the rent forms. → *the rent triad T3*
- **Ch. 38 — General Remarks on Differential Rent** `[delta-M]` — **surplus-profit converted into rent**; the natural monopoly of land; differential rent in general; **capitalised rent = rent ÷ interest = the price of land**. → *differential rent / capitalised land-price T2*
- **Ch. 39 — Differential Rent I** `[delta-M]` — **soil grades** (fertility); DR I from differences in fertility and location at equal capital; the location factor; the worst soil regulates the price. → *DR I table T2*
- **Ch. 40 — Differential Rent II** `[delta-M]` — **successive capital investments** on the same land; DR II from the differential productivity of successive doses. → *DR II T2*
- **Ch. 41 — DR II, Constant Price of Production** `[delta-M]` — the capital-productivity case; intensive vs extensive comparison under a constant regulating price. → *T2*
- **Ch. 42 — DR II, Falling Price of Production** `[delta-M]` — the soil-exclusion event; rent outcomes under a falling regulating price; normal capital per acre. → *T2*
- **Ch. 43 — DR II, Rising Price of Production** `[delta-M]` — rent sequences under a rising price; Engels's worked tables. → *T2*
- **Ch. 44 — Differential Rent Also on the Worst Soil** `[delta-M]` — rent emerging even on the worst (no-rent) soil via improvement capital and regulating-price shifts. → *T2*
- **Ch. 45 — Absolute Ground-Rent** `[delta-M]` ★ — **absolute rent** from the *lower organic composition* of agriculture (value > price of production); the value–price gap; landed property as the barrier that captures it; the absolute-rent limit. → *absolute rent T2; the monopoly barrier T3*
- **Ch. 46 — Building-Site Rent, Rent in Mining, the Price of Land** `[delta-M]` — building-site rent; mining rent; **monopoly-price rent**; land-price scenarios (capitalised rent). → *land price T2; monopoly rent T3*
- **Ch. 47 — Genesis of Capitalist Ground-Rent** `[historical, delta-M]` — historical rent forms (**labour-rent → rent-in-kind → money-rent → capitalist rent**); the rent-form stages; small-peasant production; the historical transition. → *T4 historical rent genesis*

### Part VII — The Revenues and Their Sources

- **Ch. 48 — The Trinity Formula** `[whole, delta-M]` ★★ — the **trinity formula** (Capital → interest, Land → rent, Labour → wages) as the ideological synthesis of all bourgeois economics; each term an "impossible combination"; **all revenues derive from surplus-labour alone**; the deepening mystification; the **realm of necessity vs the realm of freedom** (shortening the working day). → *★★ the culminating mystification — pairs with all the fetishisms; the three revenue streams T1/T2; realm of freedom T3*
- **Ch. 49 — Concerning the Analysis of the Process of Production** `[whole, P, delta-M]` — the **annual value composition c + v + s**; only v + s is newly created; the **revenue limit** (distributable = v + s; c must be reproduced); **Smith's "adding-up" dogma demolished** (it omits c); the surplus-value partition (profit / interest / rent). → *annual value composition T2; Smith's dogma T3*
- **Ch. 50 — Illusions Created by Competition** `[whole, delta-M]` ★ — competition's **systematic illusions**; the cost-price fetish (c+v appears as "real value," profit as a mere addition in exchange); the annual revenue account; the three illusions (profit-from-exchange, interest-from-capital-as-thing, wages-pay-labour). → *★ a gallery of fetishes T3; cost-price fetish T2*
- **Ch. 51 — Distribution Relations and Production Relations** `[whole]` — distribution relations as the **other side of production relations** (both historically transitory); the two senses of distribution (of conditions = premise; of product = result); capitalist production's two features; the **productive-forces / social-form contradiction → crisis and the historical limit**. → *distribution = production T3; the contradiction T3*
- **Ch. 52 — Classes** `[historical, whole]` ★ — the **three great classes** (wage-labourers/wages, capitalists/profit, landowners/rent); class defined by the **relation to the means of production** (not by income source); the polarising tendency (concentration, expropriation); the unanswered "what constitutes a class?" (*Capital* breaks off here). → *★ the three classes as the culmination; class tendencies T2; the unfinished ending T3/T4*

---

### Synthesis — the highest-value net-new representations

The inventory above shows the Observatory today covers the production/general-law core
plus the average rate of profit. The biggest gaps — and the most rewarding net-new
**T1/T2** vizzes for Claude Design to prioritise — are:

1. **The reproduction schemes** (Vol. II Ch. 20–21) — the two-Department flow is a
   natural living diagram and the missing bridge between the surface field and the abode.
2. **The transformation & prices of production** (Vol. III Ch. 9–10) — the orrery's
   centre of gravity (p̄′) *is* this; surface the equalisation/capital-flow dynamics that
   currently sit only as a static number.
3. **The TRPF in motion** (Vol. III Ch. 13–15) — the rising-composition → falling-p′
   trajectory, the rate-mass contradiction, and crisis-as-restoration: a deep-stratum
   counterpart to the abode's immiseration chart.
4. **The interest / fictitious-capital layer** (Vol. III Ch. 21–25) — M…M′, capitalisation,
   the fetish; the deepest, most "external" form, and a strong candidate for the floor of
   the world.
5. **The rent surface** (Vol. III Ch. 37–47) — differential & absolute rent as a literal
   terrain/landscape overlay.
6. **The fetishism thread** (Ch. 1 → Ch. 24 → Ch. 48 → Ch. 50) — a recurring **T3**
   motif binding the whole descent: commodity fetish → capital fetish → trinity → competition
   illusions. A single connective device for these would carry much of the "minor concepts"
   load gracefully.
<!-- INVENTORY:END -->
<!-- INVENTORY:END -->

---

## 7. Constraints / non-negotiables

- **Stack (fixed):** React 18 + Vite + TypeScript. Canvas for the orrery.
- **No web router.** Navigation is hash-based today (`#/`, `#/chapters`); a deeper
  world may need richer in-app routing, but **react-router is a project
  anti-pattern** — if routing must grow, propose a lightweight hash/state scheme.
- **Wire types mirror Go.** Every new snapshot field appears in `web/src/types.ts`
  as a `snake_case` mirror of the Go DTO. Keep them in sync (`sync-types` skill).
- **Accessibility parity.** Reduced-motion must yield a fully usable static
  experience (the current world already branches on `prefers-reduced-motion`).
  Keyboard navigation and semantic structure for the new IA.
- **Brand & aesthetic.** Follow the `capital-simulator-design` skill (palette,
  type, spacing, voice) and the existing `atlas.css` gilded/observatory language.
  Colour carries meaning: gold = surplus/value, red = capital/brand, lead-blue =
  necessary/constant.
- **The gate motif stays.** Ch.6 "No admittance except on business" as the
  production threshold.
- **Ephemerality.** Per-session in-memory runs; currency/speed/reduced prefs in
  localStorage; no Atlas MySQL writes.
- **Money/units contract.** Money is integer pence; rates are basis points;
  labour is minutes. No fake data — design against the real contract.
- **Performance.** The orrery is canvas; at the scale of many capitals + new
  strata, Claude Design should flag where WebGL or virtualization is warranted.

---

## 8. Open questions handed to Claude Design

1. **Fate of the standalone Chapters page** (`#/chapters`, the
   `registry.ts`-driven index). Fold entirely into the world, or keep as a flat
   fallback/index (a11y + deep-link target)?
2. **Canvas vs WebGL** for the orrery once strata + many capitals coexist.
3. **How deep can navigation go** before it genuinely needs routing rather than
   scroll-depth + pan?
4. **Cross-cutting registers.** How `whole` and `historical` chapters surface —
   as their own stratum/region, or as overlays on the circuit grid?
5. **Density management.** 106 chapters is a lot of surface area; what is the
   progressive-disclosure strategy (overview → cell → chapter → concept)?

---

## 9. Reference materials

- Prior aesthetic handoff: `docs/superpowers/specs/2026-06-05-atlas-design-handoff.md`
- Atlas design lineage: `2026-06-04-atlas-observatory-design.md`,
  `2026-06-04-atlas-general-law-design.md`,
  `2026-06-05-atlas-ephemeral-session-runs-design.md`
- Chapter registry (routing keys): `web/src/chapters/registry.ts`
- Current Observatory: `web/src/atlas/*`
- Abode domain & general law: `services/simulation-engine/internal/simulation/abode.go`
- Snapshot DTO: `services/simulation-engine/internal/transport/httpapi/observatory_handler.go`
- Fixed/circulating capital: `services/simulation-engine/internal/composition/composition.go`
- Architecture + chapter roadmap: `docs/architecture.md`
- Project framing (the three volumes as depth): `CLAUDE.md` → "Volumes"
- Chapter specs (concept tables, fixtures): red-vault Obsidian vault,
  `marx-engels/<year>/capital-volume-<roman>/specs/NN-<slug>.spec.md`
- Design system: `capital-simulator-design` skill
