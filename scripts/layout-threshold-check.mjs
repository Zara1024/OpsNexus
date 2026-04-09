import { createRequire } from 'node:module';
import path from 'node:path';
import process from 'node:process';

const require = createRequire(import.meta.url);
const { chromium } = require('playwright');

function parseArgs(argv) {
  const parsed = {};
  for (let index = 0; index < argv.length; index += 1) {
    const token = argv[index];
    if (!token.startsWith('--')) continue;
    const key = token.slice(2);
    const next = argv[index + 1];
    if (next && !next.startsWith('--')) {
      parsed[key] = next;
      index += 1;
    } else {
      parsed[key] = 'true';
    }
  }
  return parsed;
}

function toInteger(value, fallback) {
  const parsed = Number.parseInt(value ?? '', 10);
  return Number.isFinite(parsed) ? parsed : fallback;
}

async function collectMetrics(page) {
  return page.evaluate(() => {
    const textOf = (element) => {
      if (!element) return '';
      return (element.innerText || element.textContent || '').replace(/\s+/g, ' ').trim();
    };
    const selectorOf = (element) => {
      if (!element || !element.tagName) return '';
      const tag = element.tagName.toLowerCase();
      const id = element.id ? `#${element.id}` : '';
      const classes = Array.from(element.classList || []).slice(0, 3).map(name => `.${name}`).join('');
      return `${tag}${id}${classes}`;
    };
    const isVisible = (element) => {
      const style = window.getComputedStyle(element);
      const rect = element.getBoundingClientRect();
      return style.display !== 'none' && style.visibility !== 'hidden' && Number(style.opacity) !== 0 && rect.width > 0 && rect.height > 0;
    };

    const visibleElements = Array.from(document.body.querySelectorAll('*')).filter(isVisible);
    const wideScrollContainers = visibleElements.map((element) => {
      const overflow = Math.round(element.scrollWidth - element.clientWidth);
      if (overflow <= 48) return null;
      return {
        selector: selectorOf(element),
        text: textOf(element).slice(0, 180),
        clientWidth: element.clientWidth,
        scrollWidth: element.scrollWidth,
        overflow
      };
    }).filter(Boolean).sort((a, b) => b.overflow - a.overflow).slice(0, 10);

    const clippedTextElements = visibleElements.filter((element) => element.children.length === 0).map((element) => {
      const overflow = Math.round(element.scrollWidth - element.clientWidth);
      if (overflow <= 24) return null;
      const text = textOf(element);
      if (!text || text.length < 8) return null;
      return {
        selector: selectorOf(element),
        text: text.slice(0, 180),
        clientWidth: element.clientWidth,
        scrollWidth: element.scrollWidth,
        overflow
      };
    }).filter(Boolean).sort((a, b) => b.overflow - a.overflow).slice(0, 10);

    return {
      pageTitle: textOf(document.querySelector('.page-title, .ops-shell__route-title, h1')),
      wideScrollContainers,
      clippedTextElements
    };
  });
}

async function main() {
  const args = parseArgs(process.argv.slice(2));
  const route = args.route || '/dashboard';
  const baseUrl = (args['base-url'] || 'http://localhost:8080').replace(/\/+$/, '');
  const stateFile = path.resolve(args['state-file'] || 'tmp/opsnexus-browser-artifacts/local-ui-verify/opsnexus-storage-state.json');
  const maxWideOverflow = toInteger(args['max-wide-overflow'], 300);
  const maxTextOverflow = toInteger(args['max-text-overflow'], 80);
  const waitMs = toInteger(args['wait-ms'], 1800);
  const headless = String(args.headless || 'true') !== 'false';

  const browser = await chromium.launch({ headless });
  const context = await browser.newContext({
    storageState: stateFile,
    viewport: { width: 1280, height: 900 }
  });
  const page = await context.newPage();

  try {
    await page.goto(`${baseUrl}${route}`, { waitUntil: 'networkidle' });
    await page.waitForTimeout(waitMs);

    const metrics = await collectMetrics(page);
    const maxWide = Math.max(0, ...metrics.wideScrollContainers.map(item => item.overflow || 0));
    const maxText = Math.max(0, ...metrics.clippedTextElements.map(item => item.overflow || 0));
    const result = {
      route,
      finalUrl: page.url(),
      pageTitle: metrics.pageTitle,
      maxWideOverflow: maxWide,
      maxTextOverflow: maxText,
      maxWideOverflowAllowed: maxWideOverflow,
      maxTextOverflowAllowed: maxTextOverflow,
      wideScrollContainers: metrics.wideScrollContainers,
      clippedTextElements: metrics.clippedTextElements
    };

    console.log(JSON.stringify(result, null, 2));

    if (maxWide > maxWideOverflow || maxText > maxTextOverflow) {
      process.exitCode = 1;
    }
  } finally {
    await browser.close();
  }
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
