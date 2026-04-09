import Tools from '@/views/Tools/Tools'
import Agent from '@/views/Tools/Agent'
const routes = [
    {
        path: '/ops/tools',
        component: Tools,
        meta: {sTitle: '运维工具', tTitle: '工具市场'}
    },
    {
        path: '/ops/agent',
        component: Agent,
        meta: {sTitle: '运维工具', tTitle: 'Agent 列表'}
    },
]

export default routes
