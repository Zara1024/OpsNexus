import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('../src/views/app/app_quick_release.vue', import.meta.url), 'utf8')

function testQuickReleaseServiceColumnKeepsEnoughWidth() {
  assert.match(
    source,
    /<el-table-column\s+prop="applications"\s+label="发布服务"\s+min-width="240">/,
    'expected the quick release service column to reserve enough width for service chips'
  )
}

function testQuickReleaseServiceNameUsesHighContrastText() {
  assert.match(
    source,
    /\.service-name\s*\{[\s\S]*color:\s*(?:var\(--text-primary\)|#e2e8f0);/s,
    'expected quick release service names to use a high-contrast light text color'
  )
}

function testQuickReleaseServiceChipUsesDarkReadableSurface() {
  assert.match(
    source,
    /\.service-item\s*\{[\s\S]*border:\s*1px solid [^;]+;[\s\S]*background:\s*linear-gradient\([^;]+\);/s,
    'expected quick release service chips to use a bordered dark surface instead of the old flat light-theme tint'
  )

  assert.match(
    source,
    /\.environment-wrapper\s*\{[\s\S]*background:\s*rgba\(15,\s*23,\s*42,\s*0\.85\);/s,
    'expected quick release environment badges to use a dark readable backing surface'
  )
}

async function main() {
  testQuickReleaseServiceColumnKeepsEnoughWidth()
  testQuickReleaseServiceNameUsesHighContrastText()
  testQuickReleaseServiceChipUsesDarkReadableSurface()
  console.log('quick release service contrast tests passed')
}

main().catch((error) => {
  console.error(error)
  process.exitCode = 1
})
