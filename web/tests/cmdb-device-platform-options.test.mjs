import assert from 'node:assert/strict'

import { getCmdbDevicePlatformOptions } from '../src/utils/cmdbAssetPresentation.mjs'

function testDefaultPlatformOptions() {
  const options = getCmdbDevicePlatformOptions()

  assert.deepEqual(
    options.map((item) => item.value),
    ['H3C', 'Cisco', 'Huawei', 'Juniper', 'Ruijie', '其他']
  )
}

function testLegacyPlatformValueIsPreservedAsTemporaryOption() {
  const options = getCmdbDevicePlatformOptions('Aruba')

  assert.deepEqual(
    options.map((item) => item.value),
    ['H3C', 'Cisco', 'Huawei', 'Juniper', 'Ruijie', '其他', 'Aruba']
  )

  assert.equal(options.at(-1)?.label, '历史值：Aruba')
  assert.equal(options.at(-1)?.legacy, true)
}

async function main() {
  testDefaultPlatformOptions()
  testLegacyPlatformValueIsPreservedAsTemporaryOption()
  console.log('cmdb device platform option tests passed')
}

main().catch((error) => {
  console.error(error)
  process.exitCode = 1
})
