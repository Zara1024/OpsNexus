import Personal from '@/views/system/Personal.vue'
import Admin from '@/views/system/Admin.vue'
import Role from '@/views/system/Role.vue'
import Dept from '@/views/system/Dept.vue'
import Post from '@/views/system/Post.vue'
import Menu from '@/views/system/Menu.vue'
import LoginLog from '@/views/monitor/LoginLog.vue'
import Operator from '@/views/monitor/Operator.vue'
import DBLog from '@/views/monitor/DBLog.vue'
import Recording from '@/views/monitor/Recording.vue'
import AlertCenter from '@/views/monitor/Alarm-rules.vue'
import AlertNotify from '@/views/monitor/Alarm-notify.vue'
import AlertHistory from '@/views/monitor/alarm-history.vue'
import MonitorAutomation from '@/views/monitor/https.vue'

const routes = [
    {
        path: '/system/personal',
        component: Personal,
        meta: { sTitle: '个人中心', tTitle: '个人信息' }
    },
    {
        path: '/system/admin',
        component: Admin,
        meta: { sTitle: '系统管理', tTitle: '用户信息' }
    },
    {
        path: '/system/role',
        component: Role,
        meta: { sTitle: '系统管理', tTitle: '角色信息' }
    },
    {
        path: '/system/menu',
        component: Menu,
        meta: { sTitle: '系统管理', tTitle: '菜单信息' }
    },
    {
        path: '/system/dept',
        component: Dept,
        meta: { sTitle: '系统管理', tTitle: '部门信息' }
    },
    {
        path: '/system/post',
        component: Post,
        meta: { sTitle: '系统管理', tTitle: '岗位信息' }
    },
    {
        path: '/monitor/loginlog',
        component: LoginLog,
        meta: { sTitle: '操作审计', tTitle: '登录日志' }
    },
    {
        path: '/monitor/operator',
        component: Operator,
        meta: { sTitle: '操作审计', tTitle: '操作日志' }
    },
    {
        path: '/monitor/dblog',
        component: DBLog,
        meta: { sTitle: '操作审计', tTitle: '数据日志' }
    },
    {
        path: '/monitor/recording',
        component: Recording,
        meta: { sTitle: '操作审计', tTitle: '终端审计与命令审计' }
    },
    {
        path: '/monitor/alert-center',
        alias: ['/monitor/incident'],
        component: AlertCenter,
        meta: { sTitle: '监控告警', tTitle: '告警中心' }
    },
    {
        path: '/monitor/alert-notify',
        component: AlertNotify,
        meta: { sTitle: '监控告警', tTitle: '告警推送' }
    },
    {
        path: '/monitor/alert-history',
        component: AlertHistory,
        meta: { sTitle: '监控告警', tTitle: '告警历史' }
    },
    {
        path: '/monitor/https',
        alias: ['/monitor/domain'],
        component: MonitorAutomation,
        meta: { sTitle: '监控告警', tTitle: '监控深化' }
    }
]

export default routes
