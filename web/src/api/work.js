import request from '@/utils/request'

export default {
    getWorkOrderSummary() {
        return request({
            url: '/work-orders/summary',
            method: 'get'
        })
    },
    getWorkOrderList(params) {
        return request({
            url: '/work-orders',
            method: 'get',
            params
        })
    },
    getWorkOrderDetail(type, id) {
        return request({
            url: `/work-orders/${type}/${id}`,
            method: 'get'
        })
    },
    createScriptWorkOrder(data) {
        return request({
            url: '/work-orders/script',
            method: 'post',
            data
        })
    },
    approveScriptWorkOrder(id, data) {
        return request({
            url: `/work-orders/script/${id}/approve`,
            method: 'post',
            data
        })
    },
    rejectScriptWorkOrder(id, data) {
        return request({
            url: `/work-orders/script/${id}/reject`,
            method: 'post',
            data
        })
    },
    executeScriptWorkOrder(id, data) {
        return request({
            url: `/work-orders/script/${id}/execute`,
            method: 'post',
            data
        })
    }
}
