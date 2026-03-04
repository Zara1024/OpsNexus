import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import { useUserStore } from '@/stores/user'

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/login/index.vue'),
    meta: { title: '登录' },
  },
  {
    path: '/',
    component: () => import('@/components/Layout/AppLayout.vue'),
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/dashboard/index.vue'),
        meta: { title: '仪表盘' },
      },
      {
        path: 'system/user',
        name: 'SystemUser',
        component: () => import('@/views/system/user/index.vue'),
        meta: { title: '用户管理', permission: 'system:user:list' },
      },
      {
        path: 'system/role',
        name: 'SystemRole',
        component: () => import('@/views/system/role/index.vue'),
        meta: { title: '角色管理', permission: 'system:role:list' },
      },
      {
        path: 'system/menu',
        name: 'SystemMenu',
        component: () => import('@/views/system/menu/index.vue'),
        meta: { title: '菜单管理', permission: 'system:menu:list' },
      },
      {
        path: 'system/department',
        name: 'SystemDepartment',
        component: () => import('@/views/system/department/index.vue'),
        meta: { title: '部门管理', permission: 'system:department:list' },
      },
      {
        path: 'cmdb/host',
        name: 'CMDBHost',
        component: () => import('@/views/cmdb/host/index.vue'),
        meta: { title: '主机管理', permission: 'cmdb:host:list' },
      },
      {
        path: 'cmdb/group',
        name: 'CMDBGroup',
        component: () => import('@/views/cmdb/group/index.vue'),
        meta: { title: '分组管理', permission: 'cmdb:group:list' },
      },
      {
        path: 'cmdb/credential',
        name: 'CMDBCredential',
        component: () => import('@/views/cmdb/credential/index.vue'),
        meta: { title: '凭据管理', permission: 'cmdb:credential:list' },
      },
      {
        path: 'cmdb/terminal',
        name: 'CMDBTerminal',
        component: () => import('@/views/cmdb/terminal/index.vue'),
        meta: { title: 'Web终端', permission: 'cmdb:host:terminal' },
      },
    ],
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach(async (to, _from, next) => {
  const userStore = useUserStore()

  if (to.path === '/login') {
    if (userStore.isLogin) {
      next('/dashboard')
      return
    }
    next()
    return
  }

  if (!userStore.isLogin) {
    next('/login')
    return
  }

  if (!userStore.userInfo) {
    try {
      await userStore.fetchUserInfo()
    } catch (_err) {
      userStore.reset()
      next('/login')
      return
    }
  }

  const needPermission = to.meta.permission as string | undefined
  const permissionList = userStore.permissions || []
  const hasGlobalPermission = permissionList.includes('*:*:*') || permissionList.includes('system:*:*') || permissionList.includes('cmdb:*:*')
  if (needPermission && !hasGlobalPermission && !permissionList.includes(needPermission)) {
    next('/dashboard')
    return
  }

  next()
})

export default router
