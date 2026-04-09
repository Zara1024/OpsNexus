import fs from 'node:fs/promises';
import path from 'node:path';
import process from 'node:process';
import { pathToFileURL } from 'node:url';

const playwrightEntry =
  process.env.PLAYWRIGHT_ENTRY ||
  'C:\\Users\\zq\\.codex\\skills\\opsnexus-browser-ui-adjuster\\scripts\\node_modules\\playwright\\index.mjs';
const { chromium } = await import(pathToFileURL(playwrightEntry).href);

const PAGE_CASES = [
  { name: 'system-menu', route: '/system/menu', selectors: ['.page-header', '.page-toolbar', '.section-card', '.el-table', '.el-pagination'] },
  { name: 'quick-release', route: '/app/quick-release', selectors: ['.page-header', '.page-toolbar', '.section-card', '.el-table', '.el-pagination'] },
  { name: 'monitor-recording', route: '/monitor/recording', selectors: ['.page-header', '.page-toolbar', '.section-card', '.el-table', '.el-pagination'] },
  { name: 'knowledge-base', route: '/knowledge/base', selectors: ['.page-header', '.page-toolbar', '.section-card', '.el-table', '.el-pagination'] },
  { name: 'task-ansible', route: '/task/ansible', selectors: ['.page-header', '.page-toolbar', '.section-card', '.el-table', '.el-pagination'] },
  { name: 'work-orders', route: '/work/orders', selectors: ['.page-header', '.page-toolbar', '.section-card', '.el-table', '.el-pagination'] },
  { name: 'ai-assistant', route: '/ai/assistant', selectors: ['.page-header', '.section-card', '.chat-shell', '.sidebar-card'] },
  { name: 'ai-diagnosis', route: '/ai/diagnosis', selectors: ['.page-header', '.section-card', '.el-table'] },
  { name: 'dashboard', route: '/dashboard', selectors: ['.page-header', '.section-card', '.chart-card', '.tools-card'] }
];

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

function asInt(value, fallback) {
  const parsed = Number.parseInt(value ?? '', 10);
  return Number.isFinite(parsed) ? parsed : fallback;
}

async function ensureDir(dir) {
  await fs.mkdir(dir, { recursive: true });
}

async function readAuthEntries(stateFile) {
  if (!stateFile) {
    return [];
  }

  const payload = JSON.parse(await fs.readFile(stateFile, 'utf8'));
  const origins = Array.isArray(payload.origins) ? payload.origins : [];
  const origin = origins.find((item) => Array.isArray(item.localStorage) && item.localStorage.length) || origins[0];
  return Array.isArray(origin?.localStorage) ? origin.localStorage : [];
}

async function waitForStable(page, waitMs) {
  await page.waitForLoadState('networkidle');
  await page.waitForTimeout(waitMs);
}

async function bootstrapLocalApp(page, baseUrl, authEntries, waitMs) {
  await page.goto(`${baseUrl}/`, { waitUntil: 'networkidle' });
  await waitForStable(page, waitMs);

  if (!authEntries.length) {
    return;
  }

  await page.evaluate((entries) => {
    window.localStorage.clear();
    for (const entry of entries) {
      window.localStorage.setItem(entry.name, entry.value);
    }
  }, authEntries);

  await page.goto(`${baseUrl}/`, { waitUntil: 'networkidle' });
  await waitForStable(page, waitMs);
}

async function navigateToRoute(page, route, waitMs) {
  await page.evaluate((targetRoute) => {
    window.history.pushState({}, '', targetRoute);
    window.dispatchEvent(new PopStateEvent('popstate'));
  }, route);
  await page.waitForTimeout(waitMs);
}

async function collectSelectors(page, selectors) {
  return page.evaluate((inputSelectors) => {
    const collectBox = (element) => {
      if (!element) return null;
      const rect = element.getBoundingClientRect();
      return {
        x: Math.round(rect.x),
        y: Math.round(rect.y),
        width: Math.round(rect.width),
        height: Math.round(rect.height),
        right: Math.round(rect.right),
        bottom: Math.round(rect.bottom)
      };
    };

    return inputSelectors.map((selector) => {
      const nodes = Array.from(document.querySelectorAll(selector)).slice(0, 3);
      return {
        selector,
        count: nodes.length,
        matches: nodes.map((node) => ({
          text: (node.innerText || node.textContent || '').replace(/\s+/g, ' ').trim().slice(0, 180),
          box: collectBox(node)
        }))
      };
    });
  }, selectors);
}

async function capturePage(page, outputDir, baseUrl, pageCase, waitMs) {
  await navigateToRoute(page, pageCase.route, waitMs);
  await page.evaluate(() => window.scrollTo(0, 0));
  await page.waitForTimeout(200);

  const screenshotPath = path.join(outputDir, `${pageCase.name}.png`);
  await page.screenshot({ path: screenshotPath, fullPage: true });

  return {
    route: pageCase.route,
    screenshotPath,
    selectors: await collectSelectors(page, pageCase.selectors)
  };
}

async function clickByText(page, text) {
  const button = page
    .locator('button, .el-button, a')
    .filter({ hasText: text })
    .first();
  await button.waitFor({ state: 'visible', timeout: 10000 });
  await button.click();
}

async function waitForOverlaySettled(page, selector) {
  await page.locator(selector).waitFor({ state: 'visible', timeout: 10000 });
  await page.waitForTimeout(800);
}

async function captureDialogs(page, outputDir, baseUrl, waitMs) {
  const results = {};

  await navigateToRoute(page, '/task/ansible', waitMs);
  await clickByText(page, '新增任务');
  await waitForOverlaySettled(page, '.el-dialog');
  results.taskAnsibleCreate = {
    route: '/task/ansible',
    screenshotPath: path.join(outputDir, 'task-ansible-create.png'),
    selectors: await collectSelectors(page, ['.el-dialog', '.el-dialog__body', '.el-dialog__footer', '.task-form-block'])
  };
  await page.screenshot({ path: results.taskAnsibleCreate.screenshotPath, fullPage: true });
  await page.keyboard.press('Escape');

  await navigateToRoute(page, '/work/orders', waitMs);
  await clickByText(page, '详情');
  await waitForOverlaySettled(page, '.el-dialog');
  results.workOrdersDetail = {
    route: '/work/orders',
    screenshotPath: path.join(outputDir, 'work-orders-detail.png'),
    selectors: await collectSelectors(page, ['.el-dialog', '.el-dialog__body', '.el-dialog__footer', '.detail-items'])
  };
  await page.screenshot({ path: results.workOrdersDetail.screenshotPath, fullPage: true });
  await page.keyboard.press('Escape');

  await navigateToRoute(page, '/knowledge/base', waitMs);
  await clickByText(page, '新建文章');
  await waitForOverlaySettled(page, '.el-dialog');
  results.knowledgeCreate = {
    route: '/knowledge/base',
    screenshotPath: path.join(outputDir, 'knowledge-create.png'),
    selectors: await collectSelectors(page, ['.el-dialog', '.el-dialog__body', '.el-dialog__footer'])
  };
  await page.screenshot({ path: results.knowledgeCreate.screenshotPath, fullPage: true });
  await page.keyboard.press('Escape');

  await navigateToRoute(page, '/system/menu', waitMs);
  await clickByText(page, '编辑');
  await waitForOverlaySettled(page, '.el-dialog');
  results.systemMenuEdit = {
    route: '/system/menu',
    screenshotPath: path.join(outputDir, 'system-menu-edit.png'),
    selectors: await collectSelectors(page, ['.el-dialog', '.el-dialog__body', '.el-dialog__footer'])
  };
  await page.screenshot({ path: results.systemMenuEdit.screenshotPath, fullPage: true });
  await page.keyboard.press('Escape');

  await navigateToRoute(page, '/monitor/recording', waitMs);
  await clickByText(page, '详情');
  await waitForOverlaySettled(page, '.el-drawer');
  results.recordingDetail = {
    route: '/monitor/recording',
    screenshotPath: path.join(outputDir, 'recording-detail.png'),
    selectors: await collectSelectors(page, ['.el-drawer', '.el-drawer__body', '.ops-drawer__section', '.player-box'])
  };
  await page.screenshot({ path: results.recordingDetail.screenshotPath, fullPage: true });

  return results;
}

async function main() {
  const args = parseArgs(process.argv.slice(2));
  const baseUrl = (args['base-url'] || 'http://127.0.0.1:8080').replace(/\/+$/, '');
  const stateFile = args['state-file'] ? path.resolve(args['state-file']) : null;
  const injectAuthFrom = path.resolve(args['inject-auth-from'] || stateFile || 'tmp/opsnexus-browser-artifacts/local-ui-verify/opsnexus-storage-state.json');
  const outputDir = path.resolve(args['output-dir'] || 'tmp/ui-round1-verify');
  const waitMs = asInt(args['wait-ms'], 1800);
  const headless = String(args.headless || 'true') !== 'false';
  const authEntries = await readAuthEntries(injectAuthFrom);

  await ensureDir(outputDir);

  const browser = await chromium.launch({ headless });
  const context = await browser.newContext({
    viewport: { width: 1440, height: 960 }
  });
  const page = await context.newPage();

  const result = {
    baseUrl,
    stateFile,
    injectAuthFrom,
    outputDir,
    capturedAt: new Date().toISOString(),
    pages: [],
    dialogs: {}
  };

  try {
    await bootstrapLocalApp(page, baseUrl, authEntries, waitMs);

    for (const pageCase of PAGE_CASES) {
      result.pages.push(await capturePage(page, outputDir, baseUrl, pageCase, waitMs));
    }

    const dialogDir = path.join(outputDir, 'dialogs');
    await ensureDir(dialogDir);
    result.dialogs = await captureDialogs(page, dialogDir, baseUrl, waitMs);

    const summaryPath = path.join(outputDir, 'summary.json');
    await fs.writeFile(summaryPath, `${JSON.stringify(result, null, 2)}\n`, 'utf8');
    console.log(JSON.stringify({ ok: true, summaryPath, outputDir }, null, 2));
  } finally {
    await browser.close();
  }
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
