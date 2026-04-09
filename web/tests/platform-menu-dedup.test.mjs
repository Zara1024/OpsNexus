import assert from 'node:assert/strict'

import { normalizePlatformMenu } from '../src/utils/platformMenu.js'

function buildDuplicateK8sMenuInput() {
  return [
    {
      id: 81,
      menuName: '容器管理',
      url: '',
      menuSvoList: [
        { id: 82, menuName: '集群管理', url: 'k8s/list', icon: 'Menu' },
        { id: 260, menuName: '监控仪表板', url: 'k8s/monitoring', icon: 'Monitor' },
        { id: 260, menuName: '监控仪表板', url: 'k8s/monitoring', icon: 'Monitor' }
      ]
    },
    {
      id: 81,
      menuName: '容器管理',
      url: '',
      menuSvoList: [
        { id: 83, menuName: '节点管理', url: 'k8s/node', icon: 'Help' }
      ]
    }
  ]
}

function testNormalizePlatformMenuDeduplicatesDuplicateRootsAndChildren() {
  const normalized = normalizePlatformMenu(buildDuplicateK8sMenuInput())
  const k8sRootMenus = normalized.filter((item) => item.menuName === '容器管理')

  assert.equal(k8sRootMenus.length, 1, 'expected duplicate 容器管理 roots to collapse into one menu')

  const children = Array.isArray(k8sRootMenus[0]?.menuSvoList) ? k8sRootMenus[0].menuSvoList : []
  const monitoringMenus = children.filter((item) => item.url === 'k8s/monitoring')
  const nodeMenus = children.filter((item) => item.url === 'k8s/node')
  const clusterMenus = children.filter((item) => item.url === 'k8s/list')

  assert.equal(monitoringMenus.length, 1, 'expected duplicate 监控仪表板 children to collapse into one menu')
  assert.equal(nodeMenus.length, 1, 'expected 节点管理 child to remain after deduplication')
  assert.equal(clusterMenus.length, 1, 'expected 集群管理 child to remain after deduplication')
}

function testNormalizePlatformMenuCanonicalizesLegacyMenuUrls() {
  const normalized = normalizePlatformMenu([
    {
      id: 80,
      menuName: '资产管理',
      url: '',
      menuSvoList: [
        { id: 88, menuName: '网络设备', url: 'cmdb/switch', icon: 'Shop' }
      ]
    },
    {
      id: 109,
      menuName: '服务管理',
      url: '',
      menuSvoList: [
        { id: 242, menuName: '应用视图', url: 'app/view', icon: 'Help' },
        { id: 110, menuName: '应用列表', url: 'app/application', icon: 'Menu' }
      ]
    },
    {
      id: 4,
      menuName: '系统管理',
      url: '',
      menuSvoList: [
        { id: 215, menuName: '系统配置', url: 'system/config', icon: 'List' }
      ]
    },
    {
      id: 84,
      menuName: '配置中心',
      url: '',
      menuSvoList: [
        { id: 86, menuName: '通用凭据', url: 'config/accountauth', icon: 'User' }
      ]
    },
    {
      id: 44,
      menuName: '监控告警',
      url: '',
      menuSvoList: [
        { id: 213, menuName: '域名监控', url: 'monitor/domain', icon: 'Monitor' }
      ]
    },
    {
      id: 44,
      menuName: '监控告警',
      url: '',
      menuSvoList: [
        { id: 261, menuName: '监控深化', url: 'monitor/https', icon: 'Monitor' }
      ]
    },
    {
      id: 44,
      menuName: '监控告警',
      url: '',
      menuSvoList: [
        { id: 216, menuName: '故障管理', url: 'monitor/incident', icon: 'Help' }
      ]
    },
    {
      id: 44,
      menuName: '监控告警',
      url: '',
      menuSvoList: [
        { id: 256, menuName: '告警中心', url: 'monitor/alert-center', icon: 'Bell' }
      ]
    }
  ])

  const assetMenu = normalized.find((item) => item.menuName === '资产管理')
  const serviceMenu = normalized.find((item) => item.menuName === '服务管理')
  const configMenu = normalized.find((item) => item.menuName === '配置中心')
  const alertMenu = normalized.find((item) => item.menuName === '监控告警')

  const assetDeviceMenus = (assetMenu?.menuSvoList || []).filter((item) => item.url === 'cmdb/device')
  const legacyAssetMenus = (assetMenu?.menuSvoList || []).filter((item) => item.url === 'cmdb/switch')
  const appMenus = (serviceMenu?.menuSvoList || []).filter((item) => item.url === 'app/application')
  const legacyAppMenus = (serviceMenu?.menuSvoList || []).filter((item) => item.url === 'app/view')
  const configLandingMenus = (configMenu?.menuSvoList || []).filter((item) => item.url === 'config/accountauth')
  const systemMenu = normalized.find((item) => item.menuName === '系统管理')
  const alertWorkbenchMenus = (alertMenu?.menuSvoList || []).filter((item) => item.url === 'monitor/https')
  const alertCenterMenus = (alertMenu?.menuSvoList || []).filter((item) => item.url === 'monitor/alert-center')
  const legacyConfigMenus = normalized
    .flatMap((item) => Array.isArray(item.menuSvoList) ? item.menuSvoList : [])
    .filter((item) => item.url === 'system/config')
  const legacyDomainMenus = normalized
    .flatMap((item) => Array.isArray(item.menuSvoList) ? item.menuSvoList : [])
    .filter((item) => item.url === 'monitor/domain')
  const legacyIncidentMenus = normalized
    .flatMap((item) => Array.isArray(item.menuSvoList) ? item.menuSvoList : [])
    .filter((item) => item.url === 'monitor/incident')
  const duplicatedConfigMenusUnderSystem = (systemMenu?.menuSvoList || []).filter((item) => item.url === 'config/accountauth')

  assert.equal(assetDeviceMenus.length, 1, 'expected legacy cmdb/switch entry to collapse into cmdb/device')
  assert.equal(legacyAssetMenus.length, 0, 'expected legacy cmdb/switch entry to be removed after normalization')
  assert.equal(appMenus.length, 1, 'expected legacy app/view entry to collapse into app/application')
  assert.equal(legacyAppMenus.length, 0, 'expected legacy app/view entry to be removed after normalization')
  assert.equal(configLandingMenus.length, 1, 'expected legacy system/config entry to collapse into config/accountauth')
  assert.equal(legacyConfigMenus.length, 0, 'expected legacy system/config entry to be removed after normalization')
  assert.equal(duplicatedConfigMenusUnderSystem.length, 0, 'expected legacy system/config entry not to remain under 系统管理')
  assert.equal(alertWorkbenchMenus.length, 1, 'expected legacy monitor/domain entry to collapse into monitor/https')
  assert.equal(legacyDomainMenus.length, 0, 'expected legacy monitor/domain entry to be removed after normalization')
  assert.equal(alertCenterMenus.length, 1, 'expected legacy monitor/incident entry to collapse into monitor/alert-center')
  assert.equal(legacyIncidentMenus.length, 0, 'expected legacy monitor/incident entry to be removed after normalization')
}

async function main() {
  testNormalizePlatformMenuDeduplicatesDuplicateRootsAndChildren()
  testNormalizePlatformMenuCanonicalizesLegacyMenuUrls()
  console.log('platform menu dedup tests passed')
}

main().catch((error) => {
  console.error(error)
  process.exitCode = 1
})
