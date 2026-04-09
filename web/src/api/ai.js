import request from '@/utils/request'

const AI_REQUEST_TIMEOUT = 90000

export default {
    getAIOverview() {
        return request({
            url: '/ai/overview',
            method: 'get',
            timeout: AI_REQUEST_TIMEOUT
        })
    },
    getAITemplates() {
        return request({
            url: '/ai/templates',
            method: 'get',
            timeout: AI_REQUEST_TIMEOUT
        })
    },
    suggestAIKnowledge(keyword) {
        return request({
            url: '/ai/knowledge/suggest',
            method: 'get',
            params: { keyword },
            timeout: AI_REQUEST_TIMEOUT
        })
    },
    renderAIPrompt(data) {
        return request({
            url: '/ai/render',
            method: 'post',
            data,
            timeout: AI_REQUEST_TIMEOUT
        })
    },
    diagnoseAI(data) {
        return request({
            url: '/ai/diagnose',
            method: 'post',
            data,
            timeout: AI_REQUEST_TIMEOUT
        })
    },
    chatAIAssistant(data) {
        return request({
            url: '/ai/assistant/chat',
            method: 'post',
            data,
            timeout: AI_REQUEST_TIMEOUT
        })
    },
    getAIHistory() {
        return request({
            url: '/ai/history',
            method: 'get',
            timeout: AI_REQUEST_TIMEOUT
        })
    },
    getAIHistoryDetail(sessionId) {
        return request({
            url: `/ai/history/${sessionId}`,
            method: 'get',
            timeout: AI_REQUEST_TIMEOUT
        })
    },
    getAIAssistantHistory() {
        return request({
            url: '/ai/assistant/history',
            method: 'get',
            timeout: AI_REQUEST_TIMEOUT
        })
    },
    getAIAssistantHistoryDetail(sessionId) {
        return request({
            url: `/ai/assistant/history/${sessionId}`,
            method: 'get',
            timeout: AI_REQUEST_TIMEOUT
        })
    },
    getAIAssistantTemplates() {
        return request({
            url: '/ai/assistant/templates',
            method: 'get',
            timeout: AI_REQUEST_TIMEOUT
        })
    },
    getAIAssistantReports() {
        return request({
            url: '/ai/assistant/reports',
            method: 'get',
            timeout: AI_REQUEST_TIMEOUT
        })
    },
    decideAIAssistantConfirmation(id, data) {
        return request({
            url: `/ai/assistant/confirm/${id}`,
            method: 'post',
            data,
            timeout: AI_REQUEST_TIMEOUT
        })
    }
}
