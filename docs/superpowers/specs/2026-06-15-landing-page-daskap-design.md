# Landing page for daskap.io — design

**Date:** 2026-06-15
**Status:** Approved (brainstorming) — pending implementation plan

## Goal

A public landing page served at `https://daskap.io` that presents the contents
of `LETTER.md` (Aaron's first-person letter introducing the project) and a
prominent **Open App** button near the top that opens `https://app.daskap.io`
in a new browser tab.

`daskap.io` is a separate domain from the existing app, which runs as the `web`
Cloud Run service at `app.daskap.io`. The landing page is its own static site.

## Decisions (locked during brainstorming)

- **Hosting:** static HTML served by nginx in a small Docker image, deployed as
  its own Cloud Run service named `landing`, mirroring the existing `web/` →
  nginx → Cloud Run pattern. No React, no backend.
- **Content fidelity:** verbatim. `LETTER.md` is rendered exactly as written;
  only typography/formatting is applied.
- **Layout:** "Gilded Gate Hero" (option B) — a full hero gate with the primary
  CTA first, then scroll into the letter.
- **Source of truth:** `LETTER.md` stays canonical. `index.html` is generated
  from it at Docker build time so the rendered text can never drift.

## Architecture & file layout

New top-level `landing/` directory:

```
landing/
  Dockerfile        # Node build stage renders LETTER.md -> dist/index.html, then nginx serves dist
  nginx.conf        # static serve + fallback; NO /api proxy (no backend)
  template.html     # page shell: <head>, top bar, hero gate, letter container, footer
  gen.mjs           # build script: LETTER.md + template.html -> dist/index.html (markdown-it)
  package.json      # single dependency: markdown-it
  styles.css        # brand tokens + layout, copied/adapted from web/src/index.css
  public/           # favicons + webmanifest reused from web/public
```

### Content pipeline

`gen.mjs` (a ~30-line Node script run in the Docker build stage):

1. Reads `LETTER.md` from the repo root.
2. Converts markdown to HTML with `markdown-it` (default options; `#`/`##`
   become `<h1>`/`<h2>`).
3. Injects the rendered HTML into `template.html` at a `<!-- LETTER -->`
   placeholder, writes `dist/index.html`, copies `styles.css` and `public/`
   into `dist/`.
4. Asserts the output is non-empty and contains the sign-off "Solidarity" —
   fails the build otherwise.

Because the build reads `./LETTER.md`, the Docker **build context is the repo
root** (`.`), and the Dockerfile path is `landing/Dockerfile`. (The existing
`web` service uses context `./web`; `landing` differs because it reaches the
root-level `LETTER.md`.)

### Dockerfile shape

Mirrors `web/Dockerfile`:

```
FROM node:20-alpine AS build
WORKDIR /src
COPY landing/package.json landing/package-lock.json* ./
RUN npm ci || npm install
COPY landing/ ./
COPY LETTER.md /LETTER.md
RUN node gen.mjs            # writes /src/dist

FROM nginx:1.27-alpine
COPY --from=build /src/dist /usr/share/nginx/html
COPY landing/nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
```

`nginx.conf` is a plain static server (root `/usr/share/nginx/html`, `index
index.html`, `try_files $uri /index.html`). No `/api/` proxy block — the landing
page has no backend.

### Deploy wiring (`.github/workflows/deploy.yml`)

- Add a `landing` filter to the `detect` job: `landing: - 'landing/**'` and
  `- 'LETTER.md'` (so editing the letter redeploys the page).
- Add a row to the `ALL` matrix:
  `{"name":"landing","dockerfile":"landing/Dockerfile","context":".","port":80}`.
- The `daskap.io` domain mapping to the `landing` Cloud Run service is a
  one-time `gcloud run domain-mappings create` step, performed out-of-band like
  other runtime config (per the workflow's existing convention).

## Page design (Gilded Gate Hero)

**Top bar** — sticky, slim. `DasKap.io` wordmark left (gilt `.io`). `Open App ↗`
red button right, implemented as
`<a href="https://app.daskap.io" target="_blank" rel="noopener">`.

**Hero gate** — first viewport. Centered `DasKap.io` in Playfair Display, a
one-line tagline ("Capital, in motion — a simulation in Marx's own words"), a
prominent red **Open App ↗** CTA (same href/target), and a quiet "↓ read the
letter" scroll hint. Radial vignette `--surface` → `--bg`, gilt hairline rule
beneath.

**Letter body** — single centered reading column (~640px). Playfair Display for
the section headings ("The AI", "The Experiment", "The Process", "The Results",
"The Call", "The Human"); IBM Plex Sans for prose at `line-height: 1.65`. The
opening "Hello Comrade," and sign-off "Solidarity, / Aaron" render verbatim.
Gilt rules between major sections.

**Footer** — small muted line: link to the GitHub repo (the letter invites
issues) and a `daskap.io` mark.

**Visual system** — reuse existing tokens verbatim: `--bg #07080a`, `--ink
#e8e2d5`, `--red #c0392b`, `--gold #9d7a2a` / `--gold-bright #c8a240`, the
grain-texture overlay, and the `fadeUp` section animation. Webfonts: Playfair
Display, IBM Plex Sans, IBM Plex Mono (same Google Fonts import as the app).

**Responsive** — under ~700px the column goes full-width with horizontal
padding; hero type scales down; top bar stays sticky.

## Testing & verification

No backend → no unit tests. Verification:

1. `node gen.mjs` produces `dist/index.html` containing the letter's first line
   and all six section headings; the build-time assertion (contains
   "Solidarity") guards against an empty render.
2. `docker build -f landing/Dockerfile .` succeeds; the container serves `200`
   at `/` locally.
3. Playwright check: load the page, assert the hero `Open App` anchor `href` is
   `https://app.daskap.io` with `target="_blank"`, assert the letter text is
   present, and capture desktop + mobile screenshots.

## Out of scope

- Light copy-editing of the letter (decision: verbatim).
- Analytics, cookie banners, SEO beyond basic `<meta>` description/title.
- Any backend, forms, or `/api` integration.
- DNS/domain-mapping automation (one-time manual `gcloud` step).
