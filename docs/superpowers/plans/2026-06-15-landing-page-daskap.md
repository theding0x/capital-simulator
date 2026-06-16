# DasKap.io Landing Page Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Ship a static landing page at `daskap.io` that renders `LETTER.md` verbatim under a "Gilded Gate Hero" with an "Open App" button linking to `app.daskap.io`, deployed as its own nginx Cloud Run service.

**Architecture:** A new top-level `landing/` directory. At Docker build time a tiny Node script (`gen.mjs`) renders the root `LETTER.md` through `markdown-it` into a templated `index.html`; nginx serves the static `dist/`. `LETTER.md` stays the single source of truth. Wired into the existing `deploy.yml` matrix as service `landing`.

**Tech Stack:** Node 20 (build only) + markdown-it · nginx 1.27-alpine · Docker · GitHub Actions / Cloud Run. No React, no backend.

**Reference:** Design spec at `docs/superpowers/specs/2026-06-15-landing-page-daskap-design.md`. Brand tokens mirror `web/src/index.css`.

---

## File structure

| File | Responsibility |
|------|----------------|
| `landing/package.json` | Declares the single build dep (`markdown-it`) and `build` script. |
| `landing/gen.mjs` | Build script: read `LETTER.md` → markdown→HTML → inject into template → write `dist/`. Fails build if render is empty. |
| `landing/template.html` | Page shell: head, sticky top bar, hero gate, `<!-- LETTER -->` placeholder, footer. |
| `landing/styles.css` | Brand tokens + all layout (top bar, hero, letter column, footer, responsive). |
| `landing/nginx.conf` | Plain static server. No `/api` proxy. |
| `landing/Dockerfile` | Node build stage → nginx runtime. Build context = repo root (to read `LETTER.md`). |
| `landing/.dockerignore` | Exclude `dist/`, `node_modules/`. |
| `landing/public/` | Favicons + webmanifest copied from `web/public/`. |
| `.github/workflows/deploy.yml` | Add `landing` change-filter + matrix row. |

---

## Task 1: Generator scaffold renders the letter

**Files:**
- Create: `landing/package.json`
- Create: `landing/template.html`
- Create: `landing/gen.mjs`
- Create: `landing/public/` (copied favicons)

- [ ] **Step 1: Create `landing/package.json`**

```json
{
  "name": "daskap-landing",
  "private": true,
  "version": "1.0.0",
  "type": "module",
  "scripts": {
    "build": "node gen.mjs"
  },
  "dependencies": {
    "markdown-it": "^14.1.0"
  }
}
```

- [ ] **Step 2: Create `landing/template.html`** (styles.css is added in Task 2; the `<link>` is already wired here)

```html
<!doctype html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <link rel="icon" href="/favicon.ico" sizes="any" />
    <link rel="icon" type="image/png" sizes="32x32" href="/favicon-32.png" />
    <link rel="apple-touch-icon" sizes="180x180" href="/apple-touch-icon.png" />
    <link rel="manifest" href="/site.webmanifest" />
    <meta name="theme-color" content="#07080a" />
    <meta name="description" content="DasKap.io — an AI experiment simulating Marx's Capital in motion, in his own words." />
    <meta property="og:title" content="DasKap.io" />
    <meta property="og:description" content="A simulation of Capital in motion, in Marx's own words." />
    <meta property="og:type" content="website" />
    <title>DasKap.io</title>
    <link rel="preconnect" href="https://fonts.googleapis.com" />
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin />
    <link href="https://fonts.googleapis.com/css2?family=Playfair+Display:ital,wght@0,400;0,600;0,700;1,400;1,600&family=IBM+Plex+Mono:wght@400;500&family=IBM+Plex+Sans:ital,wght@0,400;0,500;0,600;1,400&display=swap" rel="stylesheet" />
    <link rel="stylesheet" href="/styles.css" />
  </head>
  <body>
    <header class="topbar">
      <a class="wordmark" href="/">Das<span class="gilt">Kap.io</span></a>
      <a class="btn-app" href="https://app.daskap.io" target="_blank" rel="noopener">Open App ↗</a>
    </header>

    <section class="hero">
      <h1 class="hero-title">DasKap<span class="gilt">.io</span></h1>
      <p class="hero-tagline">Capital, in motion — a simulation in Marx's own words.</p>
      <a class="btn-app btn-app-lg" href="https://app.daskap.io" target="_blank" rel="noopener">Open App ↗</a>
      <a class="scroll-hint" href="#letter">↓ read the letter</a>
    </section>

    <main id="letter" class="letter">
      <!-- LETTER -->
    </main>

    <footer class="footer">
      <a href="https://github.com/theding0x/capital-simulator" target="_blank" rel="noopener">Source on GitHub</a>
      <span class="dot">·</span>
      <span>daskap.io</span>
    </footer>
  </body>
</html>
```

- [ ] **Step 3: Create `landing/gen.mjs`**

```js
import { readFileSync, writeFileSync, mkdirSync, cpSync, existsSync } from 'node:fs';
import { dirname, join } from 'node:path';
import { fileURLToPath } from 'node:url';
import MarkdownIt from 'markdown-it';

const here = dirname(fileURLToPath(import.meta.url));

// In Docker the build copies LETTER.md to /LETTER.md; locally it sits at ../LETTER.md.
const letterPath = existsSync('/LETTER.md') ? '/LETTER.md' : join(here, '..', 'LETTER.md');
const markdown = readFileSync(letterPath, 'utf8');

const md = new MarkdownIt({ html: false, linkify: true, typographer: true });
const letterHtml = md.render(markdown);

if (!letterHtml.includes('Solidarity')) {
  console.error('gen.mjs: rendered letter is missing "Solidarity" — aborting build.');
  process.exit(1);
}

const template = readFileSync(join(here, 'template.html'), 'utf8');
const out = template.replace('<!-- LETTER -->', letterHtml);

const dist = join(here, 'dist');
mkdirSync(dist, { recursive: true });
writeFileSync(join(dist, 'index.html'), out, 'utf8');
cpSync(join(here, 'styles.css'), join(dist, 'styles.css'));
if (existsSync(join(here, 'public'))) {
  cpSync(join(here, 'public'), dist, { recursive: true });
}

console.log(`gen.mjs: wrote dist/index.html (${out.length} bytes)`);
```

- [ ] **Step 4: Copy favicons into `landing/public/`**

Run:
```bash
mkdir -p landing/public
cp web/public/favicon.ico web/public/favicon-16.png web/public/favicon-32.png \
   web/public/apple-touch-icon.png web/public/site.webmanifest landing/public/ 2>/dev/null || true
ls landing/public/
```
Expected: lists the favicon files that exist in `web/public/`. (If some are missing it's non-fatal — `gen.mjs` copies whatever `public/` contains.)

- [ ] **Step 5: Add a placeholder `styles.css` so the generator runs** (real styling lands in Task 2)

Create `landing/styles.css`:
```css
/* styles added in Task 2 */
```

- [ ] **Step 6: Create `landing/.gitignore`**

```
dist
node_modules
```

- [ ] **Step 7: Install deps and run the generator**

Run:
```bash
cd landing && npm install && node gen.mjs && cd ..
```
Expected: prints `gen.mjs: wrote dist/index.html (NNNNN bytes)` with a non-trivial byte count.

- [ ] **Step 8: Verify the rendered output contains the letter**

Run:
```bash
grep -c "Hello Comrade" landing/dist/index.html && \
grep -c "Solidarity" landing/dist/index.html && \
grep -c "The Experiment" landing/dist/index.html
```
Expected: three lines, each `1` or greater.

- [ ] **Step 9: Commit**

```bash
git add landing/package.json landing/package-lock.json landing/gen.mjs \
  landing/template.html landing/styles.css landing/public landing/.gitignore
git commit -m "feat(landing): generator renders LETTER.md into templated index.html"
```

---

## Task 2: Gilded Gate Hero styling

**Files:**
- Modify: `landing/styles.css` (replace the placeholder)

- [ ] **Step 1: Write the full `landing/styles.css`**

```css
:root {
  --bg:             #07080a;
  --surface:        #0f1014;
  --surface-raised: #15171d;
  --border:         #222530;
  --rule:           rgba(255, 255, 255, 0.06);
  --ink:            #e8e2d5;
  --ink-muted:      #8a8578;
  --ink-dim:        #3a3830;
  --red:            #c0392b;
  --red-hover:      #d44030;
  --gold:           #9d7a2a;
  --gold-bright:    #c8a240;
  --gold-border:    rgba(157, 122, 42, 0.22);
  --maxw:           640px;
}

* { box-sizing: border-box; margin: 0; padding: 0; }
html { scroll-behavior: smooth; }
body {
  font-family: 'IBM Plex Sans', ui-sans-serif, system-ui, sans-serif;
  background: var(--bg);
  color: var(--ink);
  line-height: 1.65;
  font-size: 16px;
  -webkit-font-smoothing: antialiased;
  min-height: 100vh;
}

/* grain texture (matches web/src/index.css) */
body::before {
  content: '';
  position: fixed;
  inset: 0;
  pointer-events: none;
  z-index: 999;
  opacity: 0.032;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='300' height='300'%3E%3Cfilter id='n'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.75' numOctaves='4' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='300' height='300' filter='url(%23n)' opacity='1'/%3E%3C/svg%3E");
  background-repeat: repeat;
}

.gilt { color: var(--gold-bright); }

/* ── top bar ─────────────────────────────────────────────── */
.topbar {
  position: sticky;
  top: 0;
  z-index: 50;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.85rem 1.5rem;
  background: rgba(7, 8, 10, 0.82);
  backdrop-filter: blur(8px);
  border-bottom: 1px solid var(--rule);
}
.wordmark {
  font-family: 'Playfair Display', serif;
  font-weight: 700;
  font-size: 1.05rem;
  letter-spacing: 0.3px;
  color: var(--ink);
  text-decoration: none;
}

.btn-app {
  font-family: 'IBM Plex Mono', monospace;
  font-size: 0.78rem;
  letter-spacing: 0.4px;
  background: var(--red);
  color: #fff;
  text-decoration: none;
  padding: 0.5rem 0.9rem;
  border-radius: 4px;
  white-space: nowrap;
  transition: background 0.18s ease, transform 0.18s ease;
}
.btn-app:hover { background: var(--red-hover); transform: translateY(-1px); }

/* ── hero gate ───────────────────────────────────────────── */
.hero {
  min-height: min(82vh, 720px);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  text-align: center;
  padding: 3rem 1.5rem 3.5rem;
  background: radial-gradient(120% 75% at 50% 0%, var(--surface-raised), var(--bg) 70%);
  border-bottom: 1px solid var(--gold-border);
}
.hero-title {
  font-family: 'Playfair Display', serif;
  font-weight: 700;
  font-size: clamp(3rem, 9vw, 5.5rem);
  line-height: 1;
  letter-spacing: -0.5px;
}
.hero-tagline {
  font-family: 'Playfair Display', serif;
  font-style: italic;
  color: var(--ink-muted);
  font-size: clamp(0.95rem, 2.4vw, 1.15rem);
  margin: 1.1rem 0 2rem;
  max-width: 34ch;
}
.btn-app-lg { font-size: 0.95rem; padding: 0.8rem 1.6rem; border-radius: 6px; }
.scroll-hint {
  margin-top: 2.5rem;
  color: var(--ink-dim);
  font-size: 0.78rem;
  text-decoration: none;
  font-family: 'IBM Plex Mono', monospace;
  transition: color 0.18s ease;
}
.scroll-hint:hover { color: var(--ink-muted); }

/* ── letter ──────────────────────────────────────────────── */
.letter {
  max-width: var(--maxw);
  margin: 0 auto;
  padding: 4rem 1.5rem 2rem;
  animation: fadeUp 0.7s cubic-bezier(0.22, 1, 0.36, 1) both;
}
.letter h1,
.letter h2 {
  font-family: 'Playfair Display', serif;
  font-weight: 700;
  line-height: 1.2;
  margin: 2.6rem 0 1rem;
  padding-top: 1.6rem;
  border-top: 1px solid var(--gold-border);
}
.letter h1 { font-size: 1.9rem; color: var(--ink); }
.letter h2 { font-size: 1.4rem; color: var(--gold-bright); border-top-color: var(--rule); }
.letter > h1:first-child,
.letter > h2:first-child { border-top: none; padding-top: 0; margin-top: 0; }
.letter p { margin: 0 0 1.25rem; }
.letter em { color: var(--ink); font-style: italic; }
.letter strong { color: var(--ink); font-weight: 600; }
.letter a {
  color: var(--gold-bright);
  text-decoration: underline;
  text-underline-offset: 2px;
}
.letter a:hover { color: var(--gold); }

/* ── footer ──────────────────────────────────────────────── */
.footer {
  max-width: var(--maxw);
  margin: 2rem auto 0;
  padding: 2.5rem 1.5rem 4rem;
  border-top: 1px solid var(--rule);
  color: var(--ink-muted);
  font-size: 0.82rem;
  font-family: 'IBM Plex Mono', monospace;
  display: flex;
  align-items: center;
  gap: 0.6rem;
  justify-content: center;
}
.footer a { color: var(--ink-muted); text-decoration: none; }
.footer a:hover { color: var(--gold-bright); }
.footer .dot { color: var(--ink-dim); }

@keyframes fadeUp {
  from { opacity: 0; transform: translateY(14px); }
  to   { opacity: 1; transform: translateY(0); }
}

@media (max-width: 700px) {
  .topbar { padding: 0.7rem 1rem; }
  .letter { padding: 2.5rem 1.25rem 1.5rem; }
  .hero { min-height: 70vh; }
}
```

- [ ] **Step 2: Re-run the generator**

Run:
```bash
cd landing && node gen.mjs && cd ..
```
Expected: `gen.mjs: wrote dist/index.html (...)`.

- [ ] **Step 3: Serve `dist/` and screenshot with Playwright**

Run (background a static server):
```bash
cd landing/dist && python3 -m http.server 8088 >/tmp/landing-http.log 2>&1 &
sleep 1 && curl -s -o /dev/null -w "%{http_code}\n" http://localhost:8088/
```
Expected: `200`.

Then drive Playwright (MCP `browser_navigate` to `http://localhost:8088/`):
- Resize to 1280×900, screenshot — confirm hero gate fills the first viewport, red "Open App" CTA centered, letter below.
- Resize to 390×800, screenshot — confirm column is full-width with padding, top bar sticky, hero type scaled down.

- [ ] **Step 4: Verify the Open App links**

Using Playwright MCP `browser_evaluate` on `http://localhost:8088/`:
```js
() => [...document.querySelectorAll('a.btn-app')].map(a => ({
  href: a.href, target: a.target, rel: a.rel
}))
```
Expected: two entries, both `href: "https://app.daskap.io/"` (or without trailing slash), `target: "_blank"`, `rel` containing `noopener`.

- [ ] **Step 5: Stop the server and commit**

```bash
pkill -f "http.server 8088" || true
git add landing/styles.css
git commit -m "feat(landing): gilded gate hero styling, brand tokens, responsive letter column"
```

---

## Task 3: Dockerize and serve via nginx

**Files:**
- Create: `landing/nginx.conf`
- Create: `landing/Dockerfile`
- Create: `landing/.dockerignore`

- [ ] **Step 1: Create `landing/nginx.conf`**

```nginx
server {
  listen 80;
  server_name _;

  root /usr/share/nginx/html;
  index index.html;

  location / {
    try_files $uri $uri/ /index.html;
  }
}
```

- [ ] **Step 2: Create `landing/.dockerignore`**

```
dist
node_modules
npm-debug.log
```

- [ ] **Step 3: Create `landing/Dockerfile`** (build context is the repo root, so paths are repo-relative)

```dockerfile
# syntax=docker/dockerfile:1.7
# Build context: repo root (needs LETTER.md). Build with:
#   docker build -f landing/Dockerfile -t daskap-landing .
FROM node:20-alpine AS build
WORKDIR /src
COPY landing/package.json landing/package-lock.json* ./
RUN if [ -f package-lock.json ]; then npm ci; else npm install; fi
COPY landing/ ./
COPY LETTER.md /LETTER.md
RUN node gen.mjs

FROM nginx:1.27-alpine
COPY --from=build /src/dist /usr/share/nginx/html
COPY landing/nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
```

- [ ] **Step 4: Build the image**

Run (from repo root):
```bash
docker build -f landing/Dockerfile -t daskap-landing .
```
Expected: build succeeds; the `RUN node gen.mjs` layer logs `gen.mjs: wrote dist/index.html`.

- [ ] **Step 5: Run the container and verify it serves**

Run:
```bash
docker run -d --name daskap-landing-test -p 8089:80 daskap-landing
sleep 1
curl -s -o /dev/null -w "%{http_code}\n" http://localhost:8089/
curl -s http://localhost:8089/ | grep -c 'href="https://app.daskap.io"'
```
Expected: `200`, then a count `>= 2` (top-bar + hero CTA).

- [ ] **Step 6: Tear down and commit**

```bash
docker rm -f daskap-landing-test
git add landing/Dockerfile landing/nginx.conf landing/.dockerignore
git commit -m "feat(landing): nginx Dockerfile renders LETTER.md at build, serves static dist"
```

---

## Task 4: Wire into the deploy workflow

**Files:**
- Modify: `.github/workflows/deploy.yml` (`detect` filters ~line 97 and `ALL` matrix ~line 116)

- [ ] **Step 1: Add the `landing` change-filter**

In the `dorny/paths-filter` `filters:` block, after the `web:` entry, add:
```yaml
            landing:
              - 'landing/**'
              - 'LETTER.md'
```

- [ ] **Step 2: Add the `landing` matrix row**

In the `ALL='[ ... ]'` heredoc, append a `landing` object after `web` (and add the trailing comma to the `web` line). The two lines become:
```json
            {"name":"web",               "dockerfile":"web/Dockerfile",                        "context":"./web", "port":80},
            {"name":"landing",           "dockerfile":"landing/Dockerfile",                    "context":".",     "port":80}
```

- [ ] **Step 3: Validate the workflow YAML and the embedded matrix JSON**

Run:
```bash
python3 -c "import yaml; yaml.safe_load(open('.github/workflows/deploy.yml')); print('yaml ok')"
```
Expected: `yaml ok`.

Then sanity-check the matrix JSON literal still parses:
```bash
python3 - <<'PY'
import re, json
src = open('.github/workflows/deploy.yml').read()
m = re.search(r"ALL='(\[.*?\])'", src, re.S)
json.loads(m.group(1))
print('matrix json ok')
PY
```
Expected: `matrix json ok`.

- [ ] **Step 4: Commit**

```bash
git add .github/workflows/deploy.yml
git commit -m "ci(landing): build and deploy landing service on landing/ or LETTER.md changes"
```

---

## Task 5: Final verification and handoff notes

- [ ] **Step 1: Clean rebuild from scratch**

Run:
```bash
rm -rf landing/dist landing/node_modules
docker build -f landing/Dockerfile -t daskap-landing .
docker run -d --name daskap-landing-final -p 8090:80 daskap-landing
sleep 1
curl -s -o /dev/null -w "%{http_code}\n" http://localhost:8090/
docker rm -f daskap-landing-final
```
Expected: `200`.

- [ ] **Step 2: Surface the one-time domain mapping to the user**

The `daskap.io` → `landing` Cloud Run domain mapping is performed out-of-band, once, by an operator with `gcloud` (documentation only — no repo change):
```bash
gcloud run domain-mappings create --service landing --domain daskap.io \
  --region "$GCP_REGION" --project "$GCP_PROJECT_ID"
```
DNS records must then be added at the registrar per the command's output.

- [ ] **Step 3: Push the branch and open a PR (only if the user asks)**

```bash
git push -u origin feature/landing-page
gh pr create --base main --title "feat: daskap.io landing page" \
  --body "Static gilded-gate landing page rendering LETTER.md verbatim with an Open App CTA to app.daskap.io. New 'landing' Cloud Run service wired into deploy.yml."
```

---

## Self-review

- **Spec coverage:** hosting (Task 3 Dockerfile/nginx + Task 4 deploy) ✓; verbatim content via build-time render (Task 1) ✓; Gilded Gate Hero layout (Task 2) ✓; Open App new-tab CTA (Task 1 template + Task 2 verification) ✓; LETTER.md as source of truth (Task 1 gen.mjs) ✓; deploy wiring + domain-mapping note (Tasks 4–5) ✓; Playwright verification (Task 2) ✓.
- **Placeholders:** none — every file's full contents are inline; the Task 1 `styles.css` stub is intentional and replaced wholesale in Task 2.
- **Type/name consistency:** `<!-- LETTER -->` placeholder matches `template.html` ↔ `gen.mjs`; class names `btn-app`/`hero`/`letter`/`wordmark` consistent between `template.html` and `styles.css`; service name `landing`, port `80`, context `.` consistent across Dockerfile, deploy matrix, and filter.
