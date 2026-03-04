import { createRouter, createWebHistory } from 'vue-router';
import { useUserStore } from '@/stores/user';
const routes = [
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
            {
                path: 'k8s/cluster',
                name: 'K8sCluster',
                component: () => import('@/views/k8s/cluster/index.vue'),
                meta: { title: 'K8s集群', permission: 'k8s:cluster:list' },
            },
            {
                path: 'k8s/node',
                name: 'K8sNode',
                component: () => import('@/views/k8s/node/index.vue'),
                meta: { title: 'K8s节点', permission: 'k8s:node:list' },
            },
            {
                path: 'k8s/namespace',
                name: 'K8sNamespace',
                component: () => import('@/views/k8s/namespace/index.vue'),
                meta: { title: '命名空间', permission: 'k8s:namespace:list' },
            },
            {
                path: 'k8s/workload/deployment',
                name: 'K8sWorkloadDeployment',
                component: () => import('@/views/k8s/workload/deployment.vue'),
                meta: { title: 'Deployment', permission: 'k8s:workload:deployment:list' },
            },
            {
                path: 'k8s/workload/pod',
                name: 'K8sWorkloadPod',
                component: () => import('@/views/k8s/workload/pod.vue'),
                meta: { title: 'Pod管理', permission: 'k8s:workload:pod:list' },
            },
            {
                path: 'k8s/network/service',
                name: 'K8sNetworkService',
                component: () => import('@/views/k8s/network/service.vue'),
                meta: { title: 'Service', permission: 'k8s:network:service:list' },
            },
            {
                path: 'k8s/network/ingress',
                name: 'K8sNetworkIngress',
                component: () => import('@/views/k8s/network/ingress.vue'),
                meta: { title: 'Ingress', permission: 'k8s:network:ingress:list' },
            },
            {
                path: 'k8s/config/configmap',
                name: 'K8sConfigMap',
                component: () => import('@/views/k8s/config/configmap.vue'),
                meta: { title: 'ConfigMap', permission: 'k8s:config:configmap:list' },
            },
            {
                path: 'k8s/config/secret',
                name: 'K8sSecret',
                component: () => import('@/views/k8s/config/secret.vue'),
                meta: { title: 'Secret', permission: 'k8s:config:secret:list' },
            },
            {
                path: 'monitor/alert-rule',
                name: 'MonitorAlertRule',
                component: () => import('@/views/monitor/alert-rule/index.vue'),
                meta: { title: '告警规则', permission: 'monitor:rule:list' },
            },
            {
                path: 'monitor/alert',
                name: 'MonitorAlert',
                component: () => import('@/views/monitor/alert/index.vue'),
                meta: { title: '活跃告警', permission: 'monitor:alert:list' },
            },
            {
                path: 'monitor/alert/history',
                name: 'MonitorAlertHistory',
                component: () => import('@/views/monitor/alert/history.vue'),
                meta: { title: '历史告警', permission: 'monitor:alert:history' },
            },
            {
                path: 'monitor/channel',
                name: 'MonitorChannel',
                component: () => import('@/views/monitor/channel/index.vue'),
                meta: { title: '通知渠道', permission: 'monitor:channel:list' },
            },
            {
                path: 'audit/operation',
                name: 'AuditOperation',
                component: () => import('@/views/audit/operation/index.vue'),
                meta: { title: '操作日志', permission: 'audit:log:list' },
            },
            {
                path: 'audit/login',
                name: 'AuditLogin',
                component: () => import('@/views/audit/login/index.vue'),
                meta: { title: '登录日志', permission: 'audit:login:list' },
            },
        ],
    },
];
const router = createRouter({
    history: createWebHistory(),
    routes,
});
router.beforeEach(async (to, _from, next) => {
    const userStore = useUserStore();
    if (to.path === '/login') {
        if (userStore.isLogin) {
            next('/dashboard');
            return;
        }
        next();
        return;
    }
    if (!userStore.isLogin) {
        next('/login');
        return;
    }
    if (!userStore.userInfo) {
        try {
            await userStore.fetchUserInfo();
        }
        catch (_err) {
            userStore.reset();
            next('/login');
            return;
        }
    }
    const needPermission = to.meta.permission;
    const permissionList = userStore.permissions || [];
    const hasGlobalPermission = permissionList.includes('*:*:*') ||
        permissionList.includes('system:*:*') ||
        permissionList.includes('cmdb:*:*') ||
        permissionList.includes('k8s:*:*') ||
        permissionList.includes('monitor:*:*') ||
        permissionList.includes('audit:*:*') ||
        permissionList.includes('dashboard:*:*');
    if (needPermission && !hasGlobalPermission && !permissionList.includes(needPermission)) {
        next('/dashboard');
        return;
    }
    next();
});
export default router;
