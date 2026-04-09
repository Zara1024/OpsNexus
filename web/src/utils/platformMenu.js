const TITLES = Object.freeze({
  dashboard: '\u4eea\u8868\u76d8',
  asset: '\u8d44\u4ea7\u7ba1\u7406',
  host: '\u4e3b\u673a',
  device: '\u7f51\u7edc\u8bbe\u5907',
  database: '\u6570\u636e\u5e93',
  sqlWorkOrder: 'SQL\u5de5\u5355',
  k8s: '\u5bb9\u5668\u7ba1\u7406',
  cluster: '\u96c6\u7fa4\u7ba1\u7406',
  node: '\u8282\u70b9\u7ba1\u7406',
  namespace: '\u547d\u540d\u7a7a\u95f4',
  workload: '\u5de5\u4f5c\u8d1f\u8f7d',
  network: '\u7f51\u7edc\u7ba1\u7406',
  storage: '\u5b58\u50a8\u7ba1\u7406',
  k8sConfig: '\u914d\u7f6e\u7ba1\u7406',
  k8sMonitoring: '\u76d1\u63a7\u4eea\u8868\u677f',
  service: '\u670d\u52a1\u7ba1\u7406',
  application: '\u5e94\u7528\u5217\u8868',
  quickRelease: '\u5feb\u901f\u53d1\u5e03',
  task: '\u4efb\u52a1\u4e2d\u5fc3',
  taskJob: '\u4efb\u52a1\u4f5c\u4e1a',
  taskTemplate: '\u4efb\u52a1\u6a21\u677f',
  taskAnsible: 'Ansible\u4efb\u52a1',
  taskConfig: '\u914d\u7f6e\u7ba1\u7406',
  aiWorkspace: 'AI \u667a\u80fd\u8fd0\u7ef4\u52a9\u624b',
  assistant: '\u52a9\u624b\u5de5\u4f5c\u53f0',
  diagnosis: '\u8bca\u65ad\u5206\u6790\u53f0',
  tools: '\u5de5\u5177\u5217\u8868',
  agent: 'Agent \u5217\u8868',
  work: '\u5de5\u5355\u4e2d\u5fc3',
  workOrders: '\u5de5\u5355\u5217\u8868',
  workApply: '\u5de5\u5355\u7533\u8bf7',
  knowledge: '\u77e5\u8bc6\u5e93',
  knowledgeArticle: '\u77e5\u8bc6\u6587\u7ae0',
  alert: '\u76d1\u63a7\u544a\u8b66',
  alertCenter: '\u544a\u8b66\u4e2d\u5fc3',
  alertNotify: '\u544a\u8b66\u63a8\u9001',
  alertHistory: '\u544a\u8b66\u5386\u53f2',
  alertWorkbench: '\u76d1\u63a7\u6df1\u5316',
  audit: '\u64cd\u4f5c\u5ba1\u8ba1',
  loginLog: '\u767b\u5f55\u65e5\u5fd7',
  operatorLog: '\u64cd\u4f5c\u65e5\u5fd7',
  dbLog: '\u6570\u636e\u65e5\u5fd7',
  terminalAudit: '\u7ec8\u7aef\u5ba1\u8ba1',
  config: '\u914d\u7f6e\u4e2d\u5fc3',
  ecsKey: '\u4e3b\u673a\u51ed\u636e',
  accountAuth: '\u901a\u7528\u51ed\u636e',
  keyManage: '\u5bc6\u94a5\u7ba1\u7406',
  ldap: 'LDAP \u96c6\u6210',
  search: '\u5168\u5c40\u641c\u7d22',
  system: '\u7cfb\u7edf\u7ba1\u7406',
  admin: '\u7528\u6237\u4fe1\u606f',
  role: '\u89d2\u8272\u4fe1\u606f',
  menu: '\u83dc\u5355\u4fe1\u606f',
  post: '\u5c97\u4f4d\u4fe1\u606f',
  dept: '\u90e8\u95e8\u4fe1\u606f'
})

const MENU_TITLE_BY_ID = Object.freeze({
  4: TITLES.system,
  72: TITLES.dashboard,
  80: TITLES.asset,
  81: TITLES.k8s,
  84: TITLES.config,
  97: TITLES.task,
  101: TITLES.aiWorkspace,
  109: TITLES.service,
  247: TITLES.terminalAudit,
  254: TITLES.search,
  255: TITLES.search,
  256: TITLES.alertCenter,
  257: TITLES.alertNotify,
  258: TITLES.alertHistory,
  259: TITLES.ldap,
  260: TITLES.k8sMonitoring,
  261: TITLES.alertWorkbench,
  262: TITLES.work,
  263: TITLES.workOrders,
  264: TITLES.workApply,
  265: TITLES.knowledge,
  266: TITLES.knowledgeArticle,
  267: TITLES.diagnosis,
  268: TITLES.sqlWorkOrder,
  269: TITLES.taskConfig,
  99100: TITLES.audit,
  99101: TITLES.loginLog,
  99102: TITLES.operatorLog,
  99103: TITLES.dbLog,
  99104: TITLES.terminalAudit,
  99200: TITLES.alert,
  99201: TITLES.alertCenter,
  99202: TITLES.alertNotify,
  99203: TITLES.alertHistory,
  99204: TITLES.alertWorkbench,
  99888: TITLES.assistant,
  99989: TITLES.k8sMonitoring,
  99990: TITLES.workApply,
  99991: TITLES.aiWorkspace,
  99992: TITLES.diagnosis,
  99993: TITLES.sqlWorkOrder,
  99994: TITLES.knowledge,
  99995: TITLES.knowledgeArticle,
  99996: TITLES.work,
  99997: TITLES.workOrders,
  99998: TITLES.ldap,
  99999: TITLES.taskConfig
})

const MENU_TITLE_BY_URL = Object.freeze({
  dashboard: TITLES.dashboard,
  'cmdb/ecs': TITLES.host,
  'cmdb/switch': TITLES.device,
  'cmdb/device': TITLES.device,
  'cmdb/db': TITLES.database,
  'cmdb/sql-work-orders': TITLES.sqlWorkOrder,
  'k8s/list': TITLES.cluster,
  'k8s/node': TITLES.node,
  'k8s/namespace': TITLES.namespace,
  'k8s/workload': TITLES.workload,
  'k8s/network': TITLES.network,
  'k8s/config': TITLES.k8sConfig,
  'k8s/storage': TITLES.storage,
  'k8s/monitoring': TITLES.k8sMonitoring,
  'app/view': TITLES.application,
  'app/application': TITLES.application,
  'app/quick-release': TITLES.quickRelease,
  'task/job': TITLES.taskJob,
  'task/template': TITLES.taskTemplate,
  'task/ansible': TITLES.taskAnsible,
  'task/config': TITLES.taskConfig,
  'ops/tools': TITLES.tools,
  'ops/agent': TITLES.agent,
  'work/orders': TITLES.workOrders,
  'work/apply': TITLES.workApply,
  'knowledge/base': TITLES.knowledgeArticle,
  'ai/assistant': TITLES.assistant,
  'ai/diagnosis': TITLES.diagnosis,
  'monitor/loginlog': TITLES.loginLog,
  'monitor/operator': TITLES.operatorLog,
  'monitor/dblog': TITLES.dbLog,
  'monitor/recording': TITLES.terminalAudit,
  'monitor/alert-center': TITLES.alertCenter,
  'monitor/alert-notify': TITLES.alertNotify,
  'monitor/alert-history': TITLES.alertHistory,
  'monitor/https': TITLES.alertWorkbench,
  'config/ecskey': TITLES.ecsKey,
  'config/accountauth': TITLES.accountAuth,
  'config/keymanage': TITLES.keyManage,
  'config/ldap': TITLES.ldap,
  'search/global': TITLES.search,
  'system/admin': TITLES.admin,
  'system/role': TITLES.role,
  'system/menu': TITLES.menu,
  'system/post': TITLES.post,
  'system/dept': TITLES.dept
})

const ROOT_MENU_ORDER = Object.freeze({
  [TITLES.dashboard]: 0,
  [TITLES.asset]: 10,
  [TITLES.k8s]: 20,
  [TITLES.service]: 30,
  [TITLES.task]: 40,
  [TITLES.aiWorkspace]: 50,
  [TITLES.work]: 60,
  [TITLES.knowledge]: 70,
  [TITLES.alert]: 80,
  [TITLES.audit]: 90,
  [TITLES.config]: 100,
  [TITLES.search]: 110,
  [TITLES.system]: 120
})

const ROOT_CHILDREN = Object.freeze({
  [TITLES.asset]: [
    { id: 99980, menuName: TITLES.host, url: 'cmdb/ecs', icon: 'Monitor' },
    { id: 99981, menuName: TITLES.device, url: 'cmdb/device', icon: 'Grid' },
    { id: 99982, menuName: TITLES.database, url: 'cmdb/db', icon: 'Coin' }
  ],
  [TITLES.k8s]: [
    { id: 260, menuName: TITLES.k8sMonitoring, url: 'k8s/monitoring', icon: 'Monitor' }
  ],
  [TITLES.task]: [
    { id: 99999, menuName: TITLES.taskConfig, url: 'task/config', icon: 'Setting' }
  ],
  [TITLES.aiWorkspace]: [
    { id: 10101, menuName: TITLES.tools, url: 'ops/tools', icon: 'Tools' },
    { id: 10102, menuName: TITLES.agent, url: 'ops/agent', icon: 'Cpu' },
    { id: 99888, menuName: TITLES.assistant, url: 'ai/assistant', icon: 'ChatDotRound' },
    { id: 99992, menuName: TITLES.diagnosis, url: 'ai/diagnosis', icon: 'MagicStick' }
  ],
  [TITLES.work]: [
    { id: 99997, menuName: TITLES.workOrders, url: 'work/orders', icon: 'Tickets' },
    { id: 99990, menuName: TITLES.workApply, url: 'work/apply', icon: 'CirclePlus' },
    { id: 99993, menuName: TITLES.sqlWorkOrder, url: 'cmdb/sql-work-orders', icon: 'Tickets' }
  ],
  [TITLES.knowledge]: [
    { id: 99995, menuName: TITLES.knowledgeArticle, url: 'knowledge/base', icon: 'Reading' }
  ],
  [TITLES.alert]: [
    { id: 99201, menuName: TITLES.alertCenter, url: 'monitor/alert-center', icon: 'Bell' },
    { id: 99202, menuName: TITLES.alertNotify, url: 'monitor/alert-notify', icon: 'Promotion' },
    { id: 99203, menuName: TITLES.alertHistory, url: 'monitor/alert-history', icon: 'Histogram' },
    { id: 99204, menuName: TITLES.alertWorkbench, url: 'monitor/https', icon: 'Monitor' }
  ],
  [TITLES.audit]: [
    { id: 99101, menuName: TITLES.loginLog, url: 'monitor/loginlog', icon: 'User' },
    { id: 99102, menuName: TITLES.operatorLog, url: 'monitor/operator', icon: 'Document' },
    { id: 99103, menuName: TITLES.dbLog, url: 'monitor/dblog', icon: 'Collection' },
    { id: 99104, menuName: TITLES.terminalAudit, url: 'monitor/recording', icon: 'Coin' }
  ],
  [TITLES.config]: [
    { id: 10084, menuName: TITLES.ecsKey, url: 'config/ecskey', icon: 'Setting' },
    { id: 10085, menuName: TITLES.accountAuth, url: 'config/accountauth', icon: 'User' },
    { id: 10086, menuName: TITLES.keyManage, url: 'config/keymanage', icon: 'Key' },
    { id: 99998, menuName: TITLES.ldap, url: 'config/ldap', icon: 'Lock' }
  ]
})

const ROOT_FALLBACKS = Object.freeze([
  { id: 99991, menuName: TITLES.aiWorkspace, icon: 'MagicStick', url: '', menuSvoList: ROOT_CHILDREN[TITLES.aiWorkspace] },
  { id: 99996, menuName: TITLES.work, icon: 'Tickets', url: '', menuSvoList: ROOT_CHILDREN[TITLES.work] },
  { id: 99994, menuName: TITLES.knowledge, icon: 'Reading', url: '', menuSvoList: ROOT_CHILDREN[TITLES.knowledge] },
  { id: 99200, menuName: TITLES.alert, icon: 'BellFilled', url: '', menuSvoList: ROOT_CHILDREN[TITLES.alert] },
  { id: 99100, menuName: TITLES.audit, icon: 'Document', url: '', menuSvoList: ROOT_CHILDREN[TITLES.audit] }
])

const ROOT_BY_CHILD_URL = Object.freeze({
  'cmdb/sql-work-orders': TITLES.work,
  'config/ecskey': TITLES.config,
  'config/accountauth': TITLES.config,
  'config/keymanage': TITLES.config,
  'config/ldap': TITLES.config
})

const ROOT_PREFIX_MAP = Object.freeze([
  { title: TITLES.asset, prefixes: ['cmdb/'] },
  { title: TITLES.k8s, prefixes: ['k8s/'] },
  { title: TITLES.service, prefixes: ['app/'] },
  { title: TITLES.task, prefixes: ['task/'] },
  { title: TITLES.aiWorkspace, prefixes: ['ops/', 'ai/'] },
  { title: TITLES.work, prefixes: ['work/'] },
  { title: TITLES.knowledge, prefixes: ['knowledge/'] },
  { title: TITLES.alert, prefixes: ['monitor/alert-', 'monitor/https'] },
  { title: TITLES.audit, prefixes: ['monitor/loginlog', 'monitor/operator', 'monitor/dblog', 'monitor/recording'] },
  { title: TITLES.config, prefixes: ['config/'] },
  { title: TITLES.search, prefixes: ['search/'] },
  { title: TITLES.system, prefixes: ['system/'] }
])

const LEGACY_URL_ALIASES = Object.freeze({
  'cmdb/switch': 'cmdb/device',
  'app/view': 'app/application',
  'system/config': 'config/accountauth',
  'monitor/domain': 'monitor/https',
  'monitor/incident': 'monitor/alert-center'
})

function normalizeUrl(value) {
  const normalized = String(value || '').trim().replace(/^\/+/, '')
  return LEGACY_URL_ALIASES[normalized] || normalized
}

function menuDedupKey(item = {}) {
  const url = normalizeUrl(item.url)
  if (url) {
    return `url:${url}`
  }

  const menuName = String(item.menuName || '').trim()
  if (menuName) {
    return `menu:${menuName}`
  }

  const menuId = Number(item.id || item.Id || 0)
  if (menuId) {
    return `id:${menuId}`
  }

  return 'unknown'
}

function mergeMenuEntries(existing, incoming) {
  const merged = {
    ...existing,
    id: existing.id ?? incoming.id,
    icon: existing.icon || incoming.icon,
    menuName: existing.menuName || incoming.menuName,
    url: existing.url || incoming.url
  }

  const mergedChildren = dedupeMenuItems([
    ...(Array.isArray(existing.menuSvoList) ? existing.menuSvoList : []),
    ...(Array.isArray(incoming.menuSvoList) ? incoming.menuSvoList : [])
  ])

  if (mergedChildren.length) {
    merged.menuSvoList = mergedChildren
  } else {
    delete merged.menuSvoList
  }

  return merged
}

function dedupeMenuItems(items = []) {
  const deduped = []

  items.forEach((item) => {
    const normalizedItem = { ...item }
    if (Array.isArray(normalizedItem.menuSvoList) && normalizedItem.menuSvoList.length) {
      normalizedItem.menuSvoList = dedupeMenuItems(normalizedItem.menuSvoList)
    } else {
      delete normalizedItem.menuSvoList
    }

    const key = menuDedupKey(normalizedItem)
    const existingIndex = deduped.findIndex((current) => menuDedupKey(current) === key)
    if (existingIndex >= 0) {
      deduped[existingIndex] = mergeMenuEntries(deduped[existingIndex], normalizedItem)
      return
    }

    deduped.push(normalizedItem)
  })

  return deduped
}

function normalizeMenuItem(item) {
  const normalized = { ...item }
  const menuId = Number(normalized.id || normalized.Id || 0)
  const url = normalizeUrl(normalized.url || normalized.Url)
  normalized.url = url

  const children = Array.isArray(normalized.menuSvoList)
    ? normalized.menuSvoList.map(normalizeMenuItem)
    : Array.isArray(normalized.children)
      ? normalized.children.map(normalizeMenuItem)
      : []

  if (children.length) {
    normalized.menuSvoList = dedupeMenuItems(children)
  } else {
    delete normalized.menuSvoList
  }

  const resolvedTitle = resolveMenuTitle(menuId, url, children)
  if (resolvedTitle) {
    normalized.menuName = resolvedTitle
  }

  return normalized
}

function resolveMenuTitle(menuId, url, children) {
  if (MENU_TITLE_BY_URL[url]) {
    return MENU_TITLE_BY_URL[url]
  }

  if (MENU_TITLE_BY_ID[menuId]) {
    return MENU_TITLE_BY_ID[menuId]
  }

  if (url === 'dashboard') {
    return TITLES.dashboard
  }

  const childUrls = children.map(child => normalizeUrl(child.url))
  for (const rule of ROOT_PREFIX_MAP) {
    if (rule.prefixes.some(prefix => childUrls.some(urlItem => urlItem === prefix || urlItem.startsWith(prefix)))) {
      return rule.title
    }
  }

  return ''
}

function appendChildren(menu, items) {
  const currentChildren = Array.isArray(menu.menuSvoList) ? [...menu.menuSvoList] : []

  items.forEach((item) => {
    const normalizedItem = normalizeMenuItem(item)
    if (!currentChildren.some(current => normalizeUrl(current.url) === normalizedItem.url)) {
      currentChildren.push(normalizedItem)
    }
  })

  if (currentChildren.length) {
    menu.menuSvoList = currentChildren
  }
}

function relocateKnownChildren(menuList) {
  menuList.forEach((menu) => {
    if (!Array.isArray(menu.menuSvoList) || !menu.menuSvoList.length) {
      return
    }

    const nextChildren = menu.menuSvoList.filter((child) => {
      const normalizedUrl = normalizeUrl(child.url)
      const targetRoot = ROOT_BY_CHILD_URL[normalizedUrl]
      return !targetRoot || targetRoot === menu.menuName
    })

    if (nextChildren.length) {
      menu.menuSvoList = nextChildren
      return
    }

    delete menu.menuSvoList
  })
}

function applyChildInjections(menuList) {
  menuList.forEach((menu) => {
    const items = ROOT_CHILDREN[menu.menuName]
    if (items) {
      appendChildren(menu, items)
    }
  })
}

function applyRootFallbacks(menuList) {
  ROOT_FALLBACKS.forEach((fallback) => {
    if (!menuList.some(item => item.menuName === fallback.menuName)) {
      menuList.push(normalizeMenuItem(fallback))
    }
  })
}

export function normalizePlatformMenu(menuItems = []) {
  const normalized = dedupeMenuItems(Array.isArray(menuItems) ? menuItems.map(normalizeMenuItem) : [])

  relocateKnownChildren(normalized)
  applyChildInjections(normalized)
  applyRootFallbacks(normalized)

  return normalized
    .map((item, index) => ({ ...item, __orderIndex: index }))
    .sort((a, b) => {
      const orderA = ROOT_MENU_ORDER[a.menuName] ?? 999
      const orderB = ROOT_MENU_ORDER[b.menuName] ?? 999
      if (orderA === orderB) {
        return a.__orderIndex - b.__orderIndex
      }
      return orderA - orderB
    })
    .map((item) => {
      const nextItem = { ...item }
      delete nextItem.__orderIndex
      return nextItem
    })
}
