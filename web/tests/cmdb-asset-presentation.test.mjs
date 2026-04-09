import assert from 'node:assert/strict'

import {
  buildAssetAddress,
  buildAssetAccountLabel,
  buildAssetConnectivityTag,
  createAssetFormSections
} from '../src/utils/cmdbAssetPresentation.mjs'

function testBuildAssetAddress() {
  assert.equal(buildAssetAddress({ type: 'host', sshIp: '10.0.0.8' }), '10.0.0.8')
  assert.equal(buildAssetAddress({ type: 'device', address: '10.0.0.9' }), '10.0.0.9')
  assert.equal(
    buildAssetAddress({ type: 'database', address: 'db.internal:3306' }),
    'db.internal:3306'
  )
}

function testBuildAssetAccountLabel() {
  assert.equal(buildAssetAccountLabel({ accountName: 'root' }), 'root')
}

function testBuildAssetAccountLabelFallbacks() {
  assert.equal(buildAssetAccountLabel({ account: { name: 'admin' } }), 'admin')
  assert.equal(buildAssetAccountLabel({ owner: 'ops-team' }), 'ops-team')
}

function testBuildAssetConnectivityTag() {
  assert.deepEqual(buildAssetConnectivityTag({ reachable: true }), {
    text: '成功',
    type: 'success'
  })
  const firstTag = buildAssetConnectivityTag({ reachable: true })
  const secondTag = buildAssetConnectivityTag({ reachable: true })
  assert.notStrictEqual(firstTag, secondTag, 'connectivity helper should return a fresh tag each time')
}

function testBuildAssetConnectivityTagVariants() {
  assert.deepEqual(buildAssetConnectivityTag({ connected: true }), { text: '成功', type: 'success' })
  assert.deepEqual(buildAssetConnectivityTag({ online: true }), { text: '成功', type: 'success' })
  assert.deepEqual(buildAssetConnectivityTag({ reachable: false }), { text: '失败', type: 'danger' })
  assert.deepEqual(buildAssetConnectivityTag({ reachable: 'false' }), { text: '失败', type: 'danger' })
}

const ASSET_TYPES = ['host', 'device', 'database']

function assertBaseSection(section, type) {
  assert.ok(section, `${type} should produce at least one section`)
  assert.equal(section.title, '基本设置', `${type} first section must be the generic basic section`)
}

function testCreateAssetFormSections() {
  ASSET_TYPES.forEach((type) => {
    const sections = createAssetFormSections(type)
    assert.ok(Array.isArray(sections))
    assert.equal(
      sections.length,
      1,
      `${type} should only return the minimal base section in this task`
    )
    const [baseSection] = sections
    assertBaseSection(baseSection, type)
    assert.equal(baseSection.assetType, type, `${type} section should preserve the requested assetType`)
  })
}

function testBuildAssetAddressFallbacks() {
  assert.equal(buildAssetAddress({ type: 'host', privateIp: '10.1.1.5' }), '10.1.1.5')
  assert.equal(buildAssetAddress({ type: 'host', publicIp: '10.1.1.6' }), '10.1.1.6')
  assert.equal(buildAssetAddress({ type: 'host', address: '10.1.1.7' }), '10.1.1.7')
  assert.equal(buildAssetAddress({ type: 'host', ip: '10.1.1.8' }), '10.1.1.8')
}

async function main() {
  testBuildAssetAddress()
  testBuildAssetAccountLabel()
  testBuildAssetAccountLabelFallbacks()
  testBuildAssetConnectivityTag()
  testBuildAssetConnectivityTagVariants()
  testCreateAssetFormSections()
  testBuildAssetAddressFallbacks()
  console.log('cmdb asset presentation helpers tests passed')
}

main().catch((error) => {
  console.error(error)
  process.exitCode = 1
})
