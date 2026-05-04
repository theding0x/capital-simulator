# Fonts

This folder holds the locally-bundled mono font (Meslo LG, user-supplied) plus documentation of the Google-Fonts imports used for display + sans.

## Imports (already in colors_and_type.css)

```css
/* Display + sans — Google Fonts CDN */
@import url('https://fonts.googleapis.com/css2?family=Playfair+Display:ital,wght@0,400;0,600;0,700;1,400;1,600&family=IBM+Plex+Sans:ital,wght@0,400;0,500;0,600;1,400&display=swap');

/* Mono — local, four @font-face declarations for 400/700 × normal/italic */
```

## Families & weights

| Family             | Role                                   | Weights / styles | Source |
|--------------------|----------------------------------------|------------------|--------|
| Playfair Display   | Display — chapter titles, equations    | 400, 600, 700; italic 400, 600 | Google Fonts |
| IBM Plex Sans      | Body, descriptions, italic notes       | 400, 500, 600; italic 400 | Google Fonts |
| Meslo LG (M variant) | Labels, numerals, button text, code  | 400, 700; italic 400, 700 | Local TTF |

## Fallback stacks (set in colors_and_type.css)

- Display: `'Playfair Display', 'Iowan Old Style', Georgia, serif`
- Sans:    `'IBM Plex Sans', ui-sans-serif, system-ui, sans-serif`
- Mono:    `'Meslo LG', ui-monospace, SFMono-Regular, Menlo, Consolas, monospace`

## Meslo LG — substitution note (FLAGGED)

> **The live product at `web/src/index.css` uses IBM Plex Mono for monospace.** This design system substitutes **Meslo LG M** (a Menlo derivative, Nerd Font edition) per user upload.
>
> **Visual difference to be aware of:**
> - Meslo is a more rigid, terminal-style monospace (closer to Menlo / classic Apple coding fonts) — straighter strokes, less humanist warmth.
> - IBM Plex Mono is more humanist, with subtle slab-serif influence, a single-storey `a`, and softer terminals.
> - The two fonts have similar x-height and width so layouts won't reflow significantly, but **labels and numerals will read more "code editor," less "publication."**
>
> If you want the upstream's exact look, ship IBM Plex Mono instead (revert the `@font-face` block + the `--font-mono` stack in `colors_and_type.css`). Otherwise, embrace Meslo as a deliberate brand decision — it leans the terminal-side of the "scholarly journal × terminal" duality further toward terminal, which can read as more developer-native.

## File inventory

The user uploaded the full Meslo Nerd Font family (S/M/L width variants, DZ "dotted-zero" cuts, plus Mono and Windows-Compatible variants — 96 TTF files total). The four files actually wired up in `colors_and_type.css` are the **LG M Mono** cuts:

- `Meslo_LG_M_Regular_Nerd_Font_Complete_Mono.ttf` — 400 normal
- `Meslo_LG_M_Italic_Nerd_Font_Complete_Mono.ttf` — 400 italic
- `Meslo_LG_M_Bold_Nerd_Font_Complete_Mono.ttf` — 700 normal
- `Meslo_LG_M_Bold_Italic_Nerd_Font_Complete_Mono.ttf` — 700 italic

**Why these four?** "M" is the medium line-spacing cut (S = compact, L = airy), and the "Mono" suffix means strict fixed-width — Nerd Font icons are squeezed to one cell, which keeps tabular numbers and aligned columns honest in the commodity table.

The other 92 TTFs are kept on disk for reference but are not loaded by any CSS. Delete them if file size matters.
