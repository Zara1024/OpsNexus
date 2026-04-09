import request from '@/utils/request'

// 应用管理 API
const appApi = {
  // 创建应用
  createApplication(data) {
    return request({
      url: '/apps',
      method: 'post',
      data
    })
  },

  // 获取应用列表
  getApplicationList(params) {
    return request({
      url: '/apps',
      method: 'get',
      params
    })
  },

  // 获取应用详情
  getApplicationDetail(id) {
    return request({
      url: `/apps/${id}`,
      method: 'get'
    })
  },

  // 更新应用
  updateApplication(id, data) {
    return request({
      url: `/apps/${id}`,
      method: 'put',
      data
    })
  },

  // 删除应用
  deleteApplication(id) {
    return request({
      url: `/apps/${id}`,
      method: 'delete'
    })
  },

  // 获取应用 Jenkins 环境列表
  getAppJenkinsEnvs(id) {
    return request({
      url: `/apps/${id}/jenkins-envs`,
      method: 'get'
    })
  },

  // 新增应用 Jenkins 环境
  addAppJenkinsEnv(id, data) {
    return request({
      url: `/apps/${id}/jenkins-envs`,
      method: 'post',
      data
    })
  },

  // 更新应用 Jenkins 环境
  updateAppJenkinsEnv(id, envId, data) {
    return request({
      url: `/apps/${id}/jenkins-envs/${envId}`,
      method: 'put',
      data
    })
  },

  // 删除应用 Jenkins 环境
  deleteAppJenkinsEnv(id, envId) {
    return request({
      url: `/apps/${id}/jenkins-envs/${envId}`,
      method: 'delete'
    })
  },

  // 获取Jenkins服务器列表
  getJenkinsServers() {
    return request({
      url: '/apps/jenkins-servers',
      method: 'get'
    })
  },

  // 搜索Jenkins任务
  searchJenkinsJobs(serverId, searchKey = '') {
    return request({
      url: `/jenkins/${serverId}/jobs/search`,
      method: 'get',
      params: {
        keyword: searchKey
      }
    })
  },

  // 验证Jenkins任务是否存在
  validateJenkinsJob(data) {
    return request({
      url: '/apps/jenkins-job/validate',
      method: 'post',
      data
    })
  },

  // 获取Jenkins任务参数
  getJenkinsJobParameters(serverId, jobName) {
    return request({
      url: `/jenkins/${serverId}/jobs/${jobName}/parameters`,
      method: 'get',
      timeout: 35000
    })
  },

  // 启动Jenkins任务(带参数)
  startJenkinsJob(serverId, jobName, data) {
    return request({
      url: `/jenkins/${serverId}/jobs/${jobName}/start`,
      method: 'post',
      data
    })
  },

  // 快速发布相关API

  // 获取可发布应用列表
  getDeployableApplications(params) {
    return request({
      url: '/apps/deployment/applications',
      method: 'get',
      params
    })
  },

  // 创建快速发布流程
  createQuickDeployment(data) {
    return request({
      url: '/apps/deployment/quick',
      method: 'post',
      data
    })
  },

  // 执行快速发布流程
  executeQuickDeployment(data) {
    return request({
      url: '/apps/deployment/execute',
      method: 'post',
      data
    })
  },

  // 获取发布详情
  getDeploymentDetail(id) {
    return request({
      url: `/apps/deployment/${id}`,
      method: 'get'
    })
  },

  // 获取发布列表
  getDeploymentList(params) {
    return request({
      url: '/apps/deployment/list',
      method: 'get',
      params
    })
  },

  // 业务组筛选服务树
  getServiceTree(params) {
    return request({
      url: '/apps/service-tree',
      method: 'get',
      params
    })
  },

  // 获取业务组选项数据（级联选择器用）
  getBusinessGroupOptions() {
    return request({
      url: '/apps/business-group-options',
      method: 'get'
    })
  },

  // 获取应用环境配置
  getEnvironment(params) {
    return request({
      url: '/apps/environment',
      method: 'get',
      params
    })
  },

  // === 新增的4个核心API ===

  // 1. 启动Jenkins任务(总任务，多任务，依次串行执行)
  // 执行发布任务已存在: executeQuickDeployment

  // 2. 实时获取每个任务的日志
  getTaskLog(taskId, start = 0) {
    return request({
      url: `/apps/deployment/tasks/${taskId}/log`,
      method: 'get',
      params: { start }
    })
  },

  // 3. 实时获取每个任务的构建状态
  getTaskStatus(taskId) {
    return request({
      url: `/apps/deployment/tasks/${taskId}/status`,
      method: 'get'
    })
  },

  // 4. 停止任务
  stopTask(taskId) {
    return request({
      url: `/apps/deployment/tasks/${taskId}/stop`,
      method: 'post'
    })
  },

  // 通过Jenkins API停止任务（备用方法）
  stopJenkinsTask(serverId, jobName, buildNumber) {
    return request({
      url: `/jenkins/${serverId}/jobs/${jobName}/builds/${buildNumber}/stop`,
      method: 'post'
    })
  },

  // 批量获取任务状态（优化性能）
  getBatchTaskStatus(taskIds) {
    return request({
      url: '/apps/deployment/tasks/batch-status',
      method: 'post',
      data: { task_ids: taskIds }
    })
  },

  // 删除发布记录
  deleteDeployment(id) {
    return request({
      url: `/apps/deployment/${id}`,
      method: 'delete'
    })
  }
}

export default appApi
