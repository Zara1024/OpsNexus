import ecskey from '@/views/configcenter/ecs-key.vue'
import accountauth from '@/views/configcenter/accountauth.vue'
import KeyManage from '@/views/configcenter/KeyManage'
import LdapSettings from '@/views/configcenter/LdapSettings.vue'

const routes = [
    {
        path: '/config/ecskey',
        component: ecskey,
        meta: { sTitle: '配置中心', tTitle: '主机凭据' }
    },
    {
        path: '/config/accountauth',
        alias: ['/system/config'],
        component: accountauth,
        meta: { sTitle: '配置中心', tTitle: '通用凭据' }
    },
    {
        path: '/config/keymanage',
        component: KeyManage,
        meta: { sTitle: '配置中心', tTitle: '密钥管理' }
    },
    {
        path: '/config/ldap',
        component: LdapSettings,
        meta: { sTitle: '配置中心', tTitle: 'LDAP 集成' }
    }
]

export default routes
