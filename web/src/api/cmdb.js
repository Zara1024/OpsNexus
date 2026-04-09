import request from "@/utils/request"

export default {
    // 获取主机监控全量数据
    getHostMonitorData(params) {
        const { id, duration, start, end } = params
        let url = `monitor/hosts/${id}/all-metrics`
        const queryParams = {}
        
        if (duration) {
            queryParams.duration = duration
        } else if (start && end) {
            queryParams.start = start
            queryParams.end = end
        }
        
        return request({
            url,
            method: 'get',
            params: queryParams
        })
    },

    // 批量获取主机简要监控数据
    getHostsMonitorData(ids) {
        return request({
            url: 'monitor/hosts',
            method: 'get',
            params: { ids }
        })
    },

    // 获取主机TOP进程监控数据
    getHostTopProcesses(hostId) {
        return request({
            url: `monitor/hosts/${hostId}/top-processes`,
            method: 'get'
        })
    },

    // 获取主机TCP端口监听状态
    getHostTcpPorts(hostId) {
        return request({
            url: `monitor/hosts/${hostId}/ports`,
            method: 'get'
        })
    },    // 获取主机列表
    getHostList(params) {
        return request({
            url: 'cmdb/host/list',
            method: 'get',
            params
        })
    },
    getAllCmdbGroups() {
        return request({
            url: 'cmdb/grouplist',
            method: 'get'
        })
    },
    getCmdbGroupByName(name) {
        return request({
            url: 'cmdb/groupbyname',
            method: 'get',
            params: {
                name: name
            }
        })
    },
    createCmdbGroup(data) {
        return request({
            url: 'cmdb/groupadd',
            method: 'post',
            data: data
        })
    },
    updateCmdbGroup(data) {
        return request({
            url: 'cmdb/groupupdate',
            method: 'put',
            data: data
        })
    },
    deleteCmdbGroup(id) {
        const data = {
            id
        }
        return request({
            url: 'cmdb/groupdelete',
            method: 'delete',
            data: data
        })
    },

    // 主机管理
    createCmdbHost(data) {
        return request({
            url: 'cmdb/hostcreate',
            method: 'post',
            data: data
        })
    },
    updateCmdbHost(data) {
        return request({
            url: 'cmdb/hostupdate',
            method: 'put',
            data: data
        })
    },
    deleteCmdbHost(id) {
        return request({
            url: 'cmdb/hostdelete',
            method: 'delete',
            data: { id }
        })
    },
    getCmdbHostList(params) {
        // 确保参数名与后端API一致
        const queryParams = {
            page: params.page || 1,
            pageSize: params.pageSize || 10
        }
        
        // 添加搜索条件，确保参数名与后端API一致
        if (params.name) queryParams.name = params.name
        if (params.ip) queryParams.ip = params.ip
        if (params.status) queryParams.status = params.status
        if (params.groupId) queryParams.groupId = params.groupId

        return request({
            url: 'cmdb/hostlist',
            method: 'get',
            params: queryParams
        })
    },
    getHostAssetList(params = {}) {
        return this.getCmdbHostList(params)
    },
    getCmdbHostById(id) {
        return request({
            url: 'cmdb/hostinfo',
            method: 'get',
            params: { id }
        })
    },
    GetCmdbHostsByHostNameLike(hostName, params) {
        return request({
            url: '/cmdb/hostbyname',
            method: 'get',
            params: { 
                name: hostName,
                ...params
            }
        })
    },
    GetCmdbHostsByIP(ip, params) {
        return request({
            url: 'cmdb/hostbyip',
            method: 'get',
            params: { 
                ip,
                ...params
            }
        })
    },
    GetCmdbHostsByStatus(status, params) {
        return request({
            url: 'cmdb/hostbystatus',
            method: 'get',
            params: { 
                status,
                ...params
            }
        })
    },
    getCmdbHostsByGroupId(groupId, params) {
        return request({
            url: 'cmdb/hostgroup',
            method: 'get',
            params: {
                groupId,
                ...params
            }
        })
    },
    hostcloudcreatealiyun(data) {
        return request({
            url: 'cmdb/hostcloudcreatealiyun',
            method: 'post',
            data: data
        })
    },
    hostcloudcreatetencent(data) {
        return request({
            url: 'cmdb/hostcloudcreatetencent',
            method: 'post',
            data: data
        })
    },
    
    // 获取带主机列表的分组数据
    getGroupListWithHosts() {
        return request({
            url: 'cmdb/grouplistwithhosts',
            method: 'get'
        })
    },

    // WebSocket连接
    getHostSSHWebSocketUrl(hostId) {
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
        return `${protocol}//${window.location.host}/api/v1/cmdb/hostssh/connect/${hostId}`
    },

    // 文件上传到主机
    uploadFileToHost(hostId, formData, config = {}) {
        return request({
            url: `cmdb/hostssh/upload/${hostId}`,
            method: 'post',
            data: formData,
            headers: {
                'Content-Type': 'multipart/form-data'
            },
            ...config
        })
    },

    // SSH命令风险预检查
    previewHostCommandRisk(hostId, command) {
        return request({
            url: `cmdb/hostssh/command-risk/${hostId}`,
            method: 'get',
            params: {
                command
            }
        })
    },

    // SSH执行命令
    executeHostCommand(hostId, command, options = {}) {
        const params = {
            command
        }
        if (options.riskAck !== undefined) params.riskAck = options.riskAck
        if (options.confirmedRiskLevel !== undefined) params.confirmedRiskLevel = options.confirmedRiskLevel
        return request({
            url: `cmdb/hostssh/command/${hostId}`,
            method: 'get',
            params
        })
    },

    // 数据库管理
    createDatabase(data) {
        return request({
            url: 'cmdb/database',
            method: 'post',
            data: data
        })
    },
    updateDatabase(data) {
        return request({
            url: 'cmdb/database',
            method: 'put',
            data: data
        })
    },
    deleteDatabase(id) {
        return request({
            url: 'cmdb/database',
            method: 'delete',
            data: { id }
        })
    },
    getDatabase(id) {
        return request({
            url: 'cmdb/database/info',
            method: 'get',
            params: { id }
        })
    },
    listDatabases(params = {}) {
        return request({
            url: 'cmdb/databaselist',
            method: 'get',
            params: params
        })
    },
    getDatabaseAssetList(params = {}) {
        return this.listDatabases(params)
    },
    getDatabasesByName(name, params = {}) {
        return request({
            url: 'cmdb/database/byname',
            method: 'get',
            params: {
                name,
                ...params
            }
        })
    },
    getDatabasesByType(type, params = {}) {
        return request({
            url: 'cmdb/database/bytype',
            method: 'get',
            params: {
                type,
                ...params
            }
        })
    },
    getDatabaseAsset(id) {
        return this.getDatabase(id)
    },

    // SQL查询
    executeSelectSQL(data) {
        return request({
            url: 'cmdb/sql/select',
            method: 'post',
            data: data,
            headers: {
                'Content-Type': 'application/json'
            }
        })
    },
    // SQL插入
    executeInsertSQL(data) {
        return request({
            url: 'cmdb/sql',
            method: 'post',
            data: data,
            headers: {
                'Content-Type': 'application/json'
            }
        })
    },
    // SQL更新
    executeUpdateSQL(data) {
        return request({
            url: 'cmdb/sql',
            method: 'put',
            data: data,
            headers: {
                'Content-Type': 'application/json'
            }
        })
    },
    executeDeleteSQL(data) {
        return request({
            url: 'cmdb/sql',
            method: 'delete',
            data: data,
            headers: {
                'Content-Type': 'application/json'
            }
        })
    },
    executeRawSQL(data) {
        return request({
            url: 'cmdb/sql/execute',
            method: 'post',
            data: data,
            headers: {
                'Content-Type': 'application/json'
            }
        })
    },
    // 获取数据库列表
    executeDatabase(data) {
        return request({
            url: 'cmdb/sql/databaselist',
            method: 'post',
            data: data,
            headers: {
                'Content-Type': 'application/json'
            }
        })
    },

    // SQL日志管理
    GetCmdbSqlLogList(params) {
        return request({
            url: 'cmdb/sqlLog/list',
            method: 'get',
            params: params
        })
    },
    DeleteCmdbSqlLogById(id) {
        return request({
            url: 'cmdb/sqlLog/delete',
            method: 'delete',
            data: { id }
        })
    },
    BatchDeleteCmdbSqlLog(ids) {
        return request({
            url: 'cmdb/sqlLog/batch/delete',
            method: 'delete',
            data: { ids }
        })
    },
    CleanCmdbSqlLog() {
        return request({
            url: 'cmdb/sqlLog/clean',
            method: 'delete'
        })
    },
    getSQLWorkOrderSummary() {
        return request({
            url: 'cmdb/sql/work-orders/summary',
            method: 'get'
        })
    },
    getSQLWorkOrderList(params) {
        return request({
            url: 'cmdb/sql/work-orders',
            method: 'get',
            params
        })
    },
    getSQLWorkOrderDetail(id) {
        return request({
            url: `cmdb/sql/work-orders/${id}`,
            method: 'get'
        })
    },
    createSQLWorkOrder(data) {
        return request({
            url: 'cmdb/sql/work-orders',
            method: 'post',
            data
        })
    },
    approveSQLWorkOrder(id, data) {
        return request({
            url: `cmdb/sql/work-orders/${id}/approve`,
            method: 'post',
            data
        })
    },
    rejectSQLWorkOrder(id, data) {
        return request({
            url: `cmdb/sql/work-orders/${id}/reject`,
            method: 'post',
            data
        })
    },
    executeSQLWorkOrder(id, data) {
        return request({
            url: `cmdb/sql/work-orders/${id}/execute`,
            method: 'post',
            data
        })
    },

    // 从Excel导入主机
    ImportHostsFromExcel(data, config = {}) {
        return request({
            url: 'cmdb/hostimport',
            method: 'post',
            data: data,
            headers: {
                'Content-Type': 'multipart/form-data'
            },
            ...config
        })
    },

    // 下载主机导入模板
    DownloadHostTemplate() {
        return request({
            url: 'cmdb/hosttemplate',
            method: 'get',
            responseType: 'blob'
        })
    },

    // 同步主机配置信息
    syncHostConfig(id) {
        return request({
            url: 'cmdb/hostsync',
            method: 'post',
            data: {
                id: id
            },
            headers: {
                'Content-Type': 'application/json'
            }
        })
    },

    testHostConnectivity(hostIds) {
        const ids = Array.isArray(hostIds) ? hostIds : [hostIds]
        return request({
            url: 'cmdb/hostconnectivity',
            method: 'post',
            data: {
                hostIds: ids
            }
        })
    },

    // Agent管理
    // 部署agent到指定主机 (支持单个或批量)
    createDevice(data) {
        return request({
            url: 'cmdb/device',
            method: 'post',
            data
        })
    },
    updateDevice(data) {
        return request({
            url: 'cmdb/device',
            method: 'put',
            data
        })
    },
    deleteDevice(id) {
        return request({
            url: 'cmdb/device',
            method: 'delete',
            data: { id }
        })
    },
    listDevices(params = {}) {
        return request({
            url: 'cmdb/devicelist',
            method: 'get',
            params
        })
    },
    getDevice(id) {
        return request({
            url: 'cmdb/device/info',
            method: 'get',
            params: { id }
        })
    },
    testDeviceConnectivity(ids) {
        const deviceIds = Array.isArray(ids) ? ids : [ids]
        return request({
            url: 'cmdb/device/connectivity',
            method: 'post',
            data: { deviceIds }
        })
    },

    deployAgent(hostIds, version = '1.0.0') {
        // 确保hostIds是数组格式
        const ids = Array.isArray(hostIds) ? hostIds : [hostIds]
        return request({
            url: 'monitor/agent/deploy',
            method: 'post',
            data: {
                hostIds: ids,
                version: version
            }
        })
    },
    
    // 卸载指定主机的agent (支持单个或批量)
    uninstallAgent(hostIds) {
        // 确保hostIds是数组格式
        const ids = Array.isArray(hostIds) ? hostIds : [hostIds]
        return request({
            url: 'monitor/agent/uninstall',
            method: 'delete',
            data: {
                hostIds: ids
            }
        })
    },
    
    // 根据主机id获取agent状态
    getAgentStatus(id) {
        return request({
            url: `monitor/agent/status/${id}`,
            method: 'get'
        })
    },
    
    // 重启agent
    restartAgent(id) {
        return request({
            url: `monitor/agent/restart/${id}`,
            method: 'post'
        })
    },
    
    // 获取agent列表
    getAgentList(params) {
        return request({
            url: 'monitor/agent/list',
            method: 'get',
            params
        })
    },

    // 删除agent数据 (用于离线服务器)
    deleteAgent(id) {
        return request({
            url: `monitor/agent/delete/${id}`,
            method: 'delete'
        })
    }
}
