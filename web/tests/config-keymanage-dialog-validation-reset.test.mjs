import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('../src/views/configcenter/KeyManage.vue', import.meta.url), 'utf8')

function testKeyDialogDestroysStaleValidationStateOnClose() {
  assert.match(
    source,
    /<el-dialog[^>]*:title="dialogTitle"[^>]*destroy-on-close|<el-dialog[^>]*destroy-on-close[^>]*:title="dialogTitle"/,
    'expected the key management dialog to destroy its inner form state on close so stale validation does not survive reopen'
  )
}

async function main() {
  testKeyDialogDestroysStaleValidationStateOnClose()
  console.log('config keymanage dialog validation reset tests passed')
}

main().catch((error) => {
  console.error(error)
  process.exitCode = 1
})
