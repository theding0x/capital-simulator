---
name: capital-simulator-design
description: Use this skill to generate well-branded interfaces and assets for Capital Simulator — a microservice simulation of an economy modeled chapter-by-chapter on Marx's *Capital, Vol. I*. Contains essential design guidelines (palette, type, spacing, voice), webfont references, and a UI kit with React components for prototyping the dashboard.
user-invocable: true
---

Read the README.md file within this skill, and explore the other available files (colors_and_type.css, ui_kits/capital-simulator/, preview/, assets/).

If creating visual artifacts (slides, mocks, throwaway prototypes, etc), copy assets out and create static HTML files for the user to view. The product surface is dark, scholarly-journal-as-terminal: bone ink on near-black, Playfair Display + IBM Plex Sans + Meslo LG (locally-bundled mono, see `fonts/`), hairline rules instead of card containers, antique-gold reveal moments, and quoted Marx in italics for chrome. No emoji, no SVG icons (Unicode glyphs only), no shadows, no gradients.

If working on production code, copy assets and read the rules in README.md and colors_and_type.css to become an expert in designing with this system.

If the user invokes this skill without any other guidance, ask them what they want to build or design (a new chapter panel? a slide? a marketing page? something else entirely?), ask some clarifying questions about scope and audience, and act as an expert designer who outputs HTML artifacts or production-ready code, depending on the need.
