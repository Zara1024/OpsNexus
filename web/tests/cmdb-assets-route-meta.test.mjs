import assert from 'node:assert/strict'

import routes from '../src/router/cmdb.js'
import { normalizePlatformMenu } from '../src/utils/platformMenu.js'

const ASSET_CHILD_PATHS = ['cmdb/ecs', 'cmdb/device', 'cmdb/db']

function assertRouteTitle(path, expectedTitle) {
  const route = routes.find((entry) => entry.path === path)
  assert(route, `route ${path} must exist`)
  assert.equal(
    route.meta?.tTitle,
    expectedTitle,
    `route ${path} should declare title ${expectedTitle}`
  )
}

function testRoutesExposeAssetTitles() {
  assertRouteTitle('/cmdb/ecs', '主机')
  assertRouteTitle('/cmdb/db', '数据库')
  assertRouteTitle('/cmdb/device', '网络设备')
}

function testDeviceRouteKeepsLegacySwitchAlias() {
  const route = routes.find((entry) => entry.path === '/cmdb/device')
  assert(route, 'route /cmdb/device must exist')

  const aliases = Array.isArray(route.alias)
    ? route.alias
    : route.alias
      ? [route.alias]
      : []

  assert(
    aliases.includes('/cmdb/switch'),
    'route /cmdb/device should keep /cmdb/switch as a legacy alias'
  )
}

function buildAssetMenuInput() {
  return [
    {
      id: 1,
      menuName: '资产管理',
      menuSvoList: [
        { id: 10, menuName: '主机', url: '/cmdb/ecs' },
        { id: 11, menuName: '网络设备', url: '/cmdb/device' },
        { id: 12, menuName: '数据库', url: '/cmdb/db' }
      ]
    }
  ]
}

function normalizeUrl(value) {
  return String(value || '')
    .trim()
    .replace(/^\/+/, '')
}

function testNormalizePlatformMenuAssets() {
  const normalized = normalizePlatformMenu(buildAssetMenuInput())
  const assetMenu = normalized.find((menu) => menu.menuName === '资产管理')
  assert(assetMenu, '资产管理 root menu must survive normalization')
  const children = Array.isArray(assetMenu.menuSvoList) ? assetMenu.menuSvoList : []
  const childNames = children.map((child) => child.menuName)
  assert(childNames.includes('主机'), '主机 entry should exist under 资产管理')
  assert(childNames.includes('网络设备'), '网络设备 entry should exist under 资产管理')
  assert(childNames.includes('数据库'), '数据库 entry should exist under 资产管理')

  const normalizedPaths = children.map((child) => normalizeUrl(child.url))
  ASSET_CHILD_PATHS.forEach((path) => {
    const matches = normalizedPaths.filter((entry) => entry === path)
    assert.equal(
      matches.length,
      1,
      `there must be exactly one ${path} entry under 资产管理`
    )
  })
}

async function main() {
  testRoutesExposeAssetTitles()
  testDeviceRouteKeepsLegacySwitchAlias()
  testNormalizePlatformMenuAssets()
  console.log('cmdb assets route meta tests passed')
}

main().catch((error) => {
  console.error(error)
  process.exitCode = 1
})
