import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('../src/views/system/Menu.vue', import.meta.url), 'utf8')

function testDeleteActionUsesMenuDeletePermission() {
  assert.match(
    source,
    /<el-button\s+type="danger"\s+link\s+@click="handleMenuDelete\(row\)"\s+v-authority="\['base:menu:delete'\]">/s,
    'expected the menu delete action to require base:menu:delete'
  )

  assert.doesNotMatch(
    source,
    /v-authority="\['base:admin:delete'\]"/,
    'expected the menu page to avoid reusing the admin delete permission on menu deletion'
  )
}

async function main() {
  testDeleteActionUsesMenuDeletePermission()
  console.log('system menu authority binding tests passed')
}

main().catch((error) => {
  console.error(error)
  process.exitCode = 1
})
