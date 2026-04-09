import WorkOrderList from '@/views/work/Work-Order-list.vue'
import WorkOrderApply from '@/views/work/Work-Order.vue'

const routes = [
    {
        path: '/work/orders',
        component: WorkOrderList,
        meta: { sTitle: '工单中心', tTitle: '工单列表' }
    },
    {
        path: '/work/apply',
        component: WorkOrderApply,
        meta: { sTitle: '工单中心', tTitle: '工单申请' }
    }
]

export default routes
