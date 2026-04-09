import request from '@/utils/request'

export default {
    getKnowledgeList(params) {
        return request({
            url: '/knowledge/articles',
            method: 'get',
            params
        })
    },
    getKnowledgeDetail(id) {
        return request({
            url: `/knowledge/articles/${id}`,
            method: 'get'
        })
    },
    createKnowledgeArticle(data) {
        return request({
            url: '/knowledge/articles',
            method: 'post',
            data
        })
    },
    bootstrapKnowledgeArticles() {
        return request({
            url: '/knowledge/articles/bootstrap',
            method: 'post'
        })
    },
    updateKnowledgeArticle(id, data) {
        return request({
            url: `/knowledge/articles/${id}`,
            method: 'put',
            data
        })
    },
    deleteKnowledgeArticle(id) {
        return request({
            url: `/knowledge/articles/${id}`,
            method: 'delete'
        })
    }
}
