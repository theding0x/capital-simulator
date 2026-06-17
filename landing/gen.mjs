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

if (!letterHtml.includes('Comrade')) {
  console.error('gen.mjs: rendered letter is missing "Comrade" — aborting build.');
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
