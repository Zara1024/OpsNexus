import assert from 'node:assert/strict'
import fs from 'node:fs'
import path from 'node:path'

const globalCssPath = path.resolve('web/src/assets/css/global.css')
const globalCss = fs.readFileSync(globalCssPath, 'utf8')

function testDevicePageUsesSharedPageShellSelectors() {
  assert.match(
    globalCss,
    /\.cmdb-host-management,\s*\r?\n\s*\.cmdb-device-management,/,
    'expected cmdb-device-management to share the same page-shell selector group as cmdb-host-management'
  )
}

function testDevicePageUsesSharedTableHeaderSelectors() {
  assert.match(
    globalCss,
    /\.cmdb-host-management \.el-table__header,\s*\r?\n\s*\.cmdb-device-management \.el-table__header/,
    'expected cmdb-device-management to share the same themed table-header selector group as cmdb-host-management'
  )

  assert.match(
    globalCss,
    /\.cmdb-host-management \.el-table__header-wrapper th\.el-table__cell,\s*\r?\n\s*\.cmdb-device-management \.el-table__header-wrapper th\.el-table__cell/,
    'expected cmdb-device-management to share the same themed table-header-cell selector group as cmdb-host-management'
  )
}

function testDevicePageUsesSharedCardSelectors() {
  assert.match(
    globalCss,
    /\.cmdb-host-management \.host-card,\s*\r?\n\s*\.cmdb-device-management \.device-card/,
    'expected cmdb-device-management to share the same card selector group as cmdb-host-management'
  )
}

async function main() {
  testDevicePageUsesSharedPageShellSelectors()
  testDevicePageUsesSharedTableHeaderSelectors()
  testDevicePageUsesSharedCardSelectors()
  console.log('cmdb device style alignment tests passed')
}

main().catch((error) => {
  console.error(error)
  process.exitCode = 1
})
