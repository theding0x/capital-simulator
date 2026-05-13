# Capital Simulator Design System

A design system for **Capital Simulator** — a microservice simulation of an economy modeled chapter-by-chapter on Karl Marx's *Capital, Vol. I*.

## What the product is

Capital Simulator is a developer-facing dashboard that lets you instantiate the economic concepts of Marx's *Capital* one chapter at a time. Each chapter of the source text is implemented as its own branch + PR in the repo. The web UI exposes the resulting domain models — commodities, owners, offers, exchanges, the universal equivalent, money, capital circuits — through small, form-driven panels.

The **vibe** is a 19th-century scholarly journal rendered as a dark-mode terminal: serif display type for headings (Playfair Display), monospaced labels and numerals (Meslo LG, replacing the upstream's IBM Plex Mono — see `fonts/README.md`), italicized notes pulled directly from Marx, and a restrained palette of bone-white ink on near-black surfaces, accented with Bordeaux red and antique gold.

## Sources

This system was reverse-engineered from a single source — the user's own codebase:

- **GitHub:** `theding0x/capital-simulator` (private). The web client lives at `web/`. Marx's source text and chapter specs live in the maintainer's `red-vault` Obsidian vault, accessed via the `obsidian` MCP server (not in this repo). Project memory in `CLAUDE.md`.
- **Stack (per CLAUDE.md):** Go 1.25 monorepo · React 18 + Vite + TypeScript · MySQL 8 · Redis · Docker · k8s.
- **Web entry:** `web/src/App.tsx`, styled by `web/src/index.css` (the canonical token + component sheet).

No Figma file. No marketing site. There is one product surface: the React dashboard.

---

## Content fundamentals

The product borrows Marx's voice for chrome and uses a clipped, technical voice for everything the developer actually has to interact with. The split is deliberate.

**Voice — Marx-quoted (chrome).** Chapter headers and panel descriptions quote *Capital* directly, in italics, no attribution:
- *"The wealth of societies in which the capitalist mode of production prevails appears as an immense collection of commodities."*
- *"Commodities cannot go to market and make exchanges of their own account."*
- *"Money is a crystal formed of necessity in the course of exchanges."*

**Voice — terse imperative (UI).** Buttons and labels are low-temperature, lowercase placeholders, no marketing fluff:
- Buttons: `Register`, `Compute ratio`, `Set universal equivalent`, `Crystallise into money`, `Withdraw`, `Reveal`, `Hide`.
- Form labels: `Name`, `Use-value description`, `Concrete labour kind`, `SNLT (minutes per unit)`.
- Placeholders are domain examples from Marx: `linen`, `yards`, `weaving`, `coat`, `30`.

**Casing.** Buttons and small chrome labels are `UPPERCASE` with `0.05–0.22em` letter-spacing — typewriter caps, not display caps. Headings (chapter titles, modal titles) are mixed-case Playfair. Body copy is sentence case.

**Person.** No "I" or "you" anywhere. The product addresses the user the way a textbook does: third-person, declarative. ("Each owner brings their commodity to market…")

**Numerals.** Always tabular, always monospaced, always with their unit attached: `30 min`, `1.5 yards`, `120 labour-minutes`. Currency in the agent layer is pence, divided by 100 for £.

**No emoji. No exclamation marks. No "Welcome!"** The closest the system gets to ornament is a `✓` next to completed chapters in the sidebar and a `›` next to the active one — both rendered as text glyphs, not icons.

**Vibe.** Scholarly. Low-stakes. The code is doing serious historical work and the UI gets out of its way. When the system has to address an empty state, it simply says `No commodities registered yet.` — no illustrations, no encouragement.

---

## Visual foundations

**Palette.** Five colors do all the work.
- **Bg / surface stack:** `#07080a` (page) → `#0f1014` (surface, sticky chrome) → `#15171d` (raised) → `#1a1d25` (hover). Pure black is avoided; everything has a slight blue cast.
- **Ink:** `#e8e2d5` (bone — primary text), `#8a8578` (muted — secondary), `#3a3830` (dim — tertiary, dividers, placeholders).
- **Lead** `#4a5a8a` — primary action color (a slate-blue, never saturated). Used for buttons, links, numeric values.
- **Red** `#c0392b` — destructive only. Bordeaux, never fire-engine.
- **Gold** `#9d7a2a` (base) / `#c8a240` (bright) — reserved for the "Reveal" affordance and chapter-completion checkmarks. Conceptually marks moments where the simulation reveals the social relation behind the commodity. Use sparingly.

**Typography.** Three families, no overlap.
- **Playfair Display** (serif, italic-capable) — chapter titles, commodity names in tables, large equations. Display-only.
- **IBM Plex Sans** — body copy, descriptions, italic chapter quotes, table data.
- **Meslo LG (M variant, Nerd Font cut)** — every uppercase label, every tabular numeral, every code-like value (IDs, amounts, button text). **Substitution:** the upstream codebase uses IBM Plex Mono. We've swapped Meslo per user-supplied font upload — it reads more terminal/code-editor and less humanist than Plex Mono. See `fonts/README.md` for trade-offs and how to revert.
- Base size **15px**, line-height **1.65**. Hero title `clamp(2rem, 4.5vw, 3.2rem)` with `-0.02em` tracking. Small caps labels run `0.6–0.68rem` with `0.16–0.22em` tracking.

**Spacing.** Implicit 4px grid. Cards have `2rem 0` vertical padding, separated by `1px solid rgba(255,255,255,0.06)` horizontal rules — never card shadows or borders. Form grids use `1rem 2rem` gap. No elaborate spacing scale; the rhythm comes from the rules.

**Backgrounds.** A single fixed SVG fractal-noise grain at 3.2% opacity, on `body::before`, z-index 999. No gradients. No imagery. No illustrations. The grain gives the whole product the feel of a printed page.

**Borders.** Hairline only. `--border #222530` (1px) for input bottoms, `--rule rgba(255,255,255,0.06)` for card dividers, `--border-subtle #181b23` for table rows. Borders are dividers, never containers. There are *no* boxed cards.

**Radii.** Almost none. `--radius-sm 4px` for buttons, `--radius 8px` for the rare placeholder block, `--radius-lg 12px` defined but not used. Inputs have `border-radius: 0` — they are bottom-borderlines, not boxes.

**Shadows.** None. The system uses surface-color steps (bg → surface → raised → hover) and hairline rules instead of elevation shadows. This is intentional and load-bearing — adding a shadow breaks the "printed journal" feel.

**Transparency & blur.** No backdrop blur. Colored fills use rgba alpha for tinted regions: `--gold-bg rgba(157,122,42,0.08)`, `--lead-dim rgba(74,90,138,0.14)`, `--red-dim rgba(192,57,43,0.14)`. These are flat tints, not glassmorphism.

**Animation.** Two motion patterns, both subtle.
- **`fadeUp`** — `0.55s cubic-bezier(0.22, 1, 0.36, 1)`, 14px translateY. Cards stagger in at `0.05 / 0.12 / 0.19 / 0.26s`.
- **`revealSlide`** — `0.38s cubic-bezier(0.22, 1, 0.36, 1)`, animates `clip-path: inset(0 0 100% 0)` to `inset(0 0 0% 0)`. Used only for the Reveal panel — feels like a vellum slip unfolding.
- **No bounces. No spring. No hover lift.** Hover changes `background-color` and/or `border-color` on a 0.12s linear transition. That's it.

**Interaction states.**
- **Hover (button):** background → `--surface-raised`, border → `--ink-muted`. Primary buttons additionally swap fill from `--lead-dim` to solid `--lead`, color → white.
- **Hover (table row):** background → `--surface-raised`, with `border-radius` rounded only on the first/last cells.
- **Hover (sidebar item):** background → `--surface-hover`, no color change.
- **Press:** no dedicated press state; the system relies on the immediate hover.
- **Focus (input):** bottom border darkens to `--ink-muted`. No focus ring.
- **Disabled:** `opacity: 0.35`, `cursor: not-allowed`.

**Layout rules.**
- **Two-column app shell:** 188px sidebar + 1fr content. A 42px sticky topbar spans both columns.
- **Content max-width:** 860px reading column inside a 2.5rem-padded chapter body.
- **Sidebar items:** 0.38rem × 1rem padding, no rounding, full-bleed selection (the active row gets `--surface-raised` background, no inset).
- **Forms:** 2-column grid with `.span2` for full-width fields. Labels stack above inputs. All inputs are bottom-bordered, transparent fills.

**Iconography.** Effectively none. The product uses **text glyphs** as icons:
- `›` for active sidebar item (chevron right, 0.75rem)
- `✓` for completed chapter
- `—` (em-dash) for empty `<option>` placeholders
- `≡` for value equivalence in revealed labour relations

There are no SVG icons, no icon font, no Lucide / Heroicons usage. If you need to add an icon, prefer a Unicode glyph in Meslo LG first; only reach for an SVG if the glyph genuinely doesn't exist. Document any addition here.

**Imagery.** None. The system has zero photography, zero illustration, zero stock art. The only visual texture is the SVG grain.

**Cards.** A "card" in this system is **not a container** — it is a horizontal band with `padding: 2rem 0` and a hairline bottom rule. No background fill, no border, no shadow, no radius. The `<h2>` inside is rendered as a 0.65rem uppercase mono label, not a heading. This is the most distinctive choice in the whole system.

---

## Iconography

See the `Iconography` and `Glyphs in use` cards in the design-system tab. Summary:

- **No icon font, no SVG icon library.** Confirmed by reading the codebase — the only `<svg>` references in `index.css` are the inline data-URIs for grain texture and the select-chevron.
- **Glyph icons (in Meslo LG):** `›` `✓` `—` `≡` `→` `&middot;` `&mdash;`. Meslo's Nerd Font edition also ships an icon range (Powerline, Devicons, Material) but the design system does not currently use any — sticking with plain Unicode keeps things portable.
- **Logo:** there is no logo asset. The product wordmark is the literal string **"Capital Simulator"** rendered in Playfair Display 1rem 700 in `--gold-bright`. Treated like a typeset book title, not a logo.
- **If a future surface needs icons:** substitute Lucide via CDN at `1.5px` stroke and document the substitution. Do not add an icon library without flagging here.

No emoji is used or permitted (per CLAUDE.md: *"Do not include emojis unless the user does first."*).

---

## Index — what's in this folder

```
README.md                  — this file
SKILL.md                   — Claude Code-compatible skill manifest
colors_and_type.css        — CSS variables for colors, type, semantic styles
fonts/                     — webfont references (Google Fonts; no local TTFs)
assets/                    — wordmark, grain texture, select chevron data-URIs
preview/                   — design-system tab preview cards (HTML)
ui_kits/
  capital-simulator/       — UI kit for the dashboard
    README.md              — what's in the kit
    index.html             — interactive demo: sidebar + chapter shell + Ch.01 panels
    Sidebar.jsx
    ChapterHeader.jsx
    Card.jsx
    CommodityForm.jsx
    CommodityTable.jsx
    RevealPanel.jsx
    ExchangeRatio.jsx
    Topbar.jsx
    chapters.js            — chapter registry (mirrors web/src/chapters/registry.ts)
```

There are no slide templates — the source repo did not ship any.
