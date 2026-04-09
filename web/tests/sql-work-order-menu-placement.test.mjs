import assert from 'node:assert/strict'

import routes from '../src/router/cmdb.js'
import { normalizePlatformMenu } from '../src/utils/platformMenu.js'

function normalizeUrl(value) {
  return String(value || '')
    .trim()
    .replace(/^\/+/, '')
}

function testSqlWorkOrderRouteBelongsToWorkCenter() {
  const route = routes.find((entry) => entry.path === '/cmdb/sql-work-orders')
  assert(route, 'SQL work order route must exist')
  assert.equal(route.meta?.sTitle, '工单中心')
  assert.equal(route.meta?.tTitle, 'SQL工单')
}

function testNormalizePlatformMenuMovesSqlWorkOrderUnderWorkCenter() {
  const normalized = normalizePlatformMenu([
    {
      id: 80,
      menuName: '资产管理',
      menuSvoList: [
        { id: 1, menuName: '主机', url: '/cmdb/ecs' },
        { id: 2, menuName: '数据库', url: '/cmdb/db' },
        { id: 3, menuName: 'SQL工单', url: '/cmdb/sql-work-orders' }
      ]
    },
    {
      id: 262,
      menuName: '工单中心',
      menuSvoList: [
        { id: 4, menuName: '工单列表', url: '/work/orders' }
      ]
    }
  ])

  const assetMenu = normalized.find((item) => item.menuName === '资产管理')
  const workMenu = normalized.find((item) => item.menuName === '工单中心')

  assert(assetMenu, '资产管理 should exist after normalization')
  assert(workMenu, '工单中心 should exist after normalization')

  const assetPaths = (assetMenu.menuSvoList || []).map((item) => normalizeUrl(item.url))
  const workPaths = (workMenu.menuSvoList || []).map((item) => normalizeUrl(item.url))

  assert.equal(
    assetPaths.includes('cmdb/sql-work-orders'),
    false,
    'SQL工单 should no longer stay under 资产管理'
  )
  assert.equal(
    workPaths.filter((item) => item === 'cmdb/sql-work-orders').length,
    1,
    'SQL工单 should appear exactly once under 工单中心'
  )
}

async function main() {
  testSqlWorkOrderRouteBelongsToWorkCenter()
  testNormalizePlatformMenuMovesSqlWorkOrderUnderWorkCenter()
  console.log('sql work order menu placement tests passed')
}

main().catch((error) => {
  console.error(error)
  process.exitCode = 1
})
