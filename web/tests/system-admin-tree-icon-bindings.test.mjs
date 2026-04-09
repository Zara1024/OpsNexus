import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('../src/views/system/Admin.vue', import.meta.url), 'utf8')

function testTreeIconsUseRegisteredComponentNames() {
  assert.match(
    source,
    /<component\s+:is="!data\.parentId && expandedKeys\.includes\(node\.key\) \? 'FolderOpened' : 'Folder'"/,
    'expected the admin department tree to switch icons by registered component name strings'
  )

  assert.doesNotMatch(
    source,
    /<component\s+:is="!data\.parentId && expandedKeys\.includes\(node\.key\) \? FolderOpened : Folder"/,
    'expected the admin department tree to avoid referencing unscoped Folder icon variables in the template expression'
  )
}

async function main() {
  testTreeIconsUseRegisteredComponentNames()
  console.log('system admin tree icon binding tests passed')
}

main().catch((error) => {
  console.error(error)
  process.exitCode = 1
})
