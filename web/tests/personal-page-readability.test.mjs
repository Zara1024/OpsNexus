import assert from 'node:assert/strict'
import fs from 'node:fs'
import path from 'node:path'

const personalVuePath = path.resolve('web/src/views/system/Personal.vue')
const personalVue = fs.readFileSync(personalVuePath, 'utf8')

function testProfileValueUsesHighContrastForeground() {
  assert.match(
    personalVue,
    /\.profile-value\s*\{[\s\S]*color:\s*var\(--text-primary\)/,
    'expected profile-value text to use high-contrast foreground color'
  )
}

function testPersonalPageLabelsUseThemeTextColor() {
  assert.match(
    personalVue,
    /\.profile-form\s*:deep\(\.el-form-item__label\),[\s\S]*\.edit-form\s*:deep\(\.el-form-item__label\)\s*\{[\s\S]*color:\s*var\(--text-secondary\)/,
    'expected profile and edit form labels to use readable theme text color'
  )
}

function testPersonalPageInputsUseReadableWrapperAndTextColors() {
  assert.match(
    personalVue,
    /\.personal-management\s*:deep\(\.el-input__wrapper\),[\s\S]*\.personal-management\s*:deep\(\.el-textarea__inner\)\s*\{[\s\S]*background:\s*rgba\(15,\s*23,\s*42,\s*0\.78\)/,
    'expected personal page inputs to use a dark readable input surface'
  )

  assert.match(
    personalVue,
    /\.personal-management\s*:deep\(\.el-input__inner\),[\s\S]*\.personal-management\s*:deep\(\.el-textarea__inner\)\s*\{[\s\S]*color:\s*var\(--text-primary\)/,
    'expected personal page input text to use readable foreground color'
  )

  assert.match(
    personalVue,
    /\.personal-management\s*:deep\(\.el-input__inner::placeholder\),[\s\S]*\.personal-management\s*:deep\(\.el-textarea__inner::placeholder\)\s*\{[\s\S]*color:\s*var\(--text-muted\)/,
    'expected personal page placeholders to use readable muted text color'
  )
}

function testPersonalTabsUseReadableInactiveAndActiveColors() {
  assert.match(
    personalVue,
    /\.personal-tabs\s*:deep\(\.el-tabs__item\)\s*\{[\s\S]*color:\s*var\(--text-secondary\)/,
    'expected personal tab items to use readable inactive color'
  )

  assert.match(
    personalVue,
    /\.personal-tabs\s*:deep\(\.el-tabs__item\.is-active\)\s*\{[\s\S]*color:\s*var\(--text-primary\)/,
    'expected active personal tab item to use readable active color'
  )
}

async function main() {
  testProfileValueUsesHighContrastForeground()
  testPersonalPageLabelsUseThemeTextColor()
  testPersonalPageInputsUseReadableWrapperAndTextColors()
  testPersonalTabsUseReadableInactiveAndActiveColors()
  console.log('personal page readability tests passed')
}

main().catch((error) => {
  console.error(error)
  process.exitCode = 1
})
