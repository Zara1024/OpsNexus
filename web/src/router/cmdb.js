const routes = [
  {
    path: '/cmdb/ecs',
    component: () => import('../views/cmdb/cmdbHost.vue'),
    meta: { sTitle: '资产管理', tTitle: '主机' }
  },
  {
    path: '/cmdb/group',
    component: () => import('../views/cmdb/cmdbGroup.vue'),
    meta: { sTitle: '资产管理', tTitle: '业务分组' }
  },
  {
    path: '/cmdb/db',
    component: () => import('../views/cmdb/cmdbDB.vue'),
    meta: { sTitle: '资产管理', tTitle: '数据库' }
  },
  {
    path: '/cmdb/device',
    alias: ['/cmdb/switch'],
    component: () => import('../views/cmdb/cmdbDevice.vue'),
    meta: { sTitle: '资产管理', tTitle: '网络设备' }
  },
  {
    path: '/cmdb/ssh',
    component: () => import('../views/cmdb/Host/SSH.vue'),
    meta: { sTitle: '资产管理', tTitle: '终端登录' }
  },
  {
    path: '/cmdb/dbdetails',
    component: () => import('../views/cmdb/DBdetails.vue'),
    meta: { sTitle: '数据库', tTitle: '数据库操作' }
  },
  {
    path: '/cmdb/sql-work-orders',
    component: () => import('../views/cmdb/SQLWorkOrderCenter.vue'),
    meta: { sTitle: '工单中心', tTitle: 'SQL工单' }
  }
]

export default routes
