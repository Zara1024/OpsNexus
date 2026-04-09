import assert from 'node:assert/strict'
import fs from 'node:fs'

const hostPageSource = fs.readFileSync(new URL('../src/views/cmdb/cmdbHost.vue', import.meta.url), 'utf8')
const groupComponentSource = fs.readFileSync(new URL('../src/views/cmdb/Host/CmdbGroup.vue', import.meta.url), 'utf8')

function testHostPageKeepsRoomForTheGroupPanelHeader() {
  assert.match(
    hostPageSource,
    /\.group-tree-section\s*\{[\s\S]*width:\s*280px;/s,
    'expected the host page to keep the group sidebar wide enough for the group header card'
  )
}

function testGroupHeaderUsesDedicatedTitleClasses() {
  assert.match(
    groupComponentSource,
    /<h3 class="group-card__title">资产分组<\/h3>/,
    'expected the group header to use a dedicated title class instead of the global .title class'
  )

  assert.match(
    groupComponentSource,
    /class="group-card__subtitle">Asset Groups<\/span>/,
    'expected the group header to use a dedicated subtitle class'
  )
}

function testGroupHeaderAllowsTheTitleBlockToShrinkSafely() {
  assert.match(
    groupComponentSource,
    /\.title-wrapper\s*\{[\s\S]*flex:\s*1;[\s\S]*min-width:\s*0;/s,
    'expected the group header title wrapper to own the remaining width and shrink safely'
  )

  assert.match(
    groupComponentSource,
    /\.title-content\s*\{[\s\S]*flex:\s*1;[\s\S]*min-width:\s*0;/s,
    'expected the group header title content to shrink safely inside the title wrapper'
  )
}

async function main() {
  testHostPageKeepsRoomForTheGroupPanelHeader()
  testGroupHeaderUsesDedicatedTitleClasses()
  testGroupHeaderAllowsTheTitleBlockToShrinkSafely()
  console.log('cmdb group header layout tests passed')
}

main().catch((error) => {
  console.error(error)
  process.exitCode = 1
})
