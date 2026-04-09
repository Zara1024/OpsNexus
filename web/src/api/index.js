/**
 * 基础 API 接口管理
 */

import request from "@/utils/request"
import systemAPI from './system'
import cmdbAPI from './cmdb'
import dashboardAPI from './dashboard'
import monitorAPI from './monitor'
import searchAPI from './search'
import knowledgeAPI from './knowledge'
import configAPI from './config'
import workAPI from './work'
import aiAPI from './ai'
import * as toolAPI from './tool'

export default {
    captcha() {
        return request({
            url: 'captcha',
            method: 'get'
        })
    },
    login(params) {
        return request({
            url: 'login',
            method: 'post',
            data: params
        })
    },
    ...systemAPI,
    ...cmdbAPI,
    ...dashboardAPI,
    ...monitorAPI,
    ...searchAPI,
    ...knowledgeAPI,
    ...configAPI,
    ...workAPI,
    ...aiAPI,
    ...toolAPI,

    system: systemAPI,
    cmdb: cmdbAPI,
    dashboard: dashboardAPI,
    monitor: monitorAPI,
    search: searchAPI,
    knowledge: knowledgeAPI,
    config: configAPI,
    work: workAPI,
    ai: aiAPI,
    tool: toolAPI
}
