import request from "@/utils/request"

export default {
    globalSearch(params) {
        return request({
            url: 'search/global',
            method: 'get',
            params
        })
    }
}
