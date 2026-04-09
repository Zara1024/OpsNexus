<!-- eslint-disable vue/multi-word-component-names -->
<script setup>
/* global defineOptions */

import { ref, reactive, onMounted, computed, watch, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'

defineOptions({
  name: 'AppApplication'
})
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Plus,
  Edit,
  Delete,
  Search,
  Refresh,
  Monitor,
  Promotion,
  Setting,
  SuccessFilled,
  WarningFilled,
  Close
} from '@element-plus/icons-vue'
import appApi from '@/api/app'
import cmdbApi from '@/api/cmdb'
import k8sApi from '@/api/k8s'
import systemApi from '@/api/system'
import {
  getGovernanceDisplayText,
  getGovernanceSuggestionPath,
  getGovernanceTagType as getGovernanceDisplayTagType
} from '@/utils/governance'

// 基础状态
const loading = ref(false)
const searchKeyword = ref('')
const router = useRouter()
const route = useRoute()

// 应用列表数据
const applicationList = ref([])
const currentApplication = ref({})
const detailLoading = ref(false)

const createEmptyOverviewSummary = () => ({
  total_applications: 0,
  active_applications: 0,
  inactive_applications: 0,
  abnormal_applications: 0,
  total_environments: 0,
  configured_environments: 0,
  total_deployments: 0,
  successful_deployments: 0,
  failed_deployments: 0
})
const overviewSummary = ref(createEmptyOverviewSummary())

// 对话框状态
const createDialogVisible = ref(false)
const editDialogVisible = ref(false)
const detailDialogVisible = ref(false)

// 分页参数
const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

// 表单数据
const applicationForm = ref({
  // 必填参数
  name: '',                    // 应用名称
  business_group_id: null,     // 业务组ID
  business_dept_id: null,      // 业务部门ID

  // 基本信息（可选）
  description: '',             // 应用介绍
  repo_url: '',               // 仓库地址

  // 负责人信息（可选）
  dev_owners: [],             // 研发负责人ID数组
  test_owners: [],            // 测试负责人ID数组
  ops_owners: [],             // 运维负责人ID数组

  // 技术信息（可选）
  programming_lang: '',       // 编程语言
  start_command: '',          // 启动命令
  stop_command: '',           // 停止命令
  health_api: '',             // 健康检查接口

  // 关联资源（可选）
  domains: [],                // 关联域名数组
  hosts: [],                  // 关联主机ID数组
  databases: [],              // 关联数据库ID数组
  other_res: {                // 其他资源配置
    rabbitmq: [],
    zookeeper: [],
    kafka: [],
    redis: [],
    other: []
  },

  // Jenkins环境配置（可选）
  jenkins_envs: []            // Jenkins环境配置数组
})


// 编程语言选项
const programmingLanguages = [
  'Java', 'Go', 'Python', 'Node.js', 'C++', 'C#', '.NET', 'PHP', 'Ruby', 'Rust'
]


// 业务组选项
const businessGroups = ref([])

// 业务部门选项
const businessDepts = ref([])

// 用户选项（负责人）
const userOptions = ref([])

// Jenkins环境配置
const jenkinsEnvDialogVisible = ref(false)
const jenkinsEnvs = ref([])
const jenkinsEnvForm = ref({
  env_name: '',
  jenkins_server_id: null,
  job_name: '',
  deploy_config: {
    cluster_id: null,
    namespace: '',
    workload_type: 'deployment',
    workload_name: '',
    release_governance_enabled: false,
    pre_check_mode: 'observe'
  }
})
const editingJenkinsEnv = ref(null)
const showJenkinsEnvForm = ref(false)
const currentAppId = ref(null)
const currentAppName = ref('')
const routePreset = reactive({
  appId: 0,
  env: '',
  focus: ''
})

// Jenkins服务器和任务选项
const jenkinsServers = ref([])
const jenkinsJobs = ref([])
const k8sClusters = ref([])
const jobSearchLoading = ref(false)
const jobValidationLoading = ref(false)
const jobValidationStatus = ref(null) // null, 'success', 'error'

// 过滤后的应用列表
const filteredApplicationList = computed(() => {
  const list = Array.isArray(applicationList.value) ? applicationList.value : []
  if (!searchKeyword.value) return list
  return list.filter(item =>
    item.name?.toLowerCase().includes(searchKeyword.value.toLowerCase()) ||
    item.code?.toLowerCase().includes(searchKeyword.value.toLowerCase()) ||
    item.programming_lang?.toLowerCase().includes(searchKeyword.value.toLowerCase())
  )
})

const applyApplicationTableLayout = () => {
  nextTick(() => {
    const hiddenColumns = new Set([2, 5, 6, 9])
    const widthMap = {
      0: 160,
      1: 150,
      3: 120,
      4: 130,
      7: 92,
      8: 92,
      10: 200
    }

    const syncMainTable = (tableElement) => {
      tableElement.style.width = '1160px'
      tableElement.style.maxWidth = '1160px'

      const cols = Array.from(tableElement.querySelectorAll('colgroup col'))
      cols.forEach((col, index) => {
        if (hiddenColumns.has(index)) {
          col.style.display = 'none'
          col.style.width = '0px'
          return
        }
        if (widthMap[index]) {
          col.style.width = `${widthMap[index]}px`
        }
      })

      tableElement.querySelectorAll('thead tr, tbody tr').forEach(row => {
        Array.from(row.children).forEach((cell, index) => {
          if (hiddenColumns.has(index)) {
            cell.style.display = 'none'
            return
          }
          cell.style.display = ''
          if (widthMap[index]) {
            cell.style.width = `${widthMap[index]}px`
            cell.style.maxWidth = `${widthMap[index]}px`
          }
        })
      })
    }

    document.querySelectorAll('.app-table .el-table__header-wrapper table, .app-table .el-table__body-wrapper table').forEach(syncMainTable)

    document.querySelectorAll('.app-table .el-table__fixed-right').forEach(wrapper => {
      wrapper.style.width = '200px'
      wrapper.style.right = '0px'
    })
    document.querySelectorAll('.app-table .el-table__fixed-right colgroup col, .app-table .el-table__fixed-right .el-table__cell').forEach(element => {
      element.style.width = '200px'
      element.style.maxWidth = '200px'
    })
  })
}

const currentFocusedEnv = computed(() => {
  const envName = String(routePreset.env || '').trim().toLowerCase()
  if (!envName) return null
  return (currentApplication.value?.jenkins_envs || []).find(item =>
    String(item.env_name || '').trim().toLowerCase() === envName
  ) || null
})

// 获取应用列表
const fetchApplicationList = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      pageSize: pagination.pageSize
    }

    if (searchKeyword.value) {
      params.keyword = searchKeyword.value
    }

    const response = await appApi.getApplicationList(params)

    if (response.data?.code === 200) {
      applicationList.value = response.data.data?.list || []
      pagination.total = response.data.data?.total || 0
      overviewSummary.value = response.data.data?.summary || createEmptyOverviewSummary()
      applyApplicationTableLayout()
    } else {
      ElMessage.error(response.data?.message || '获取应用列表失败')
    }
  } catch (error) {
    console.error('获取应用列表失败:', error)
    ElMessage.error('获取应用列表失败')
  } finally {
    loading.value = false
  }
}

// 创建应用
const handleCreateApp = () => {
  resetForm()
  createDialogVisible.value = true
}

// 编辑应用
const handleEditApp = (row) => {
  currentApplication.value = { ...row }
  Object.assign(applicationForm.value, row)
  editDialogVisible.value = true
}

// 查看应用详情
const handleViewApp = async (row) => {
  detailLoading.value = true
  detailDialogVisible.value = true
  currentApplication.value = { ...row }
  currentAppId.value = row.id
  currentAppName.value = row.name
  try {
    const response = await appApi.getApplicationDetail(row.id)
    if (response.data?.code === 200) {
      currentApplication.value = response.data.data || { ...row }
    } else {
      ElMessage.error(response.data?.message || '获取应用详情失败')
    }
  } catch (error) {
    console.error('获取应用详情失败:', error)
    ElMessage.error(error.response?.data?.message || '获取应用详情失败')
  } finally {
    detailLoading.value = false
  }
}


// 处理Jenkins环境配置
const handleJenkinsEnvConfig = async (row) => {
  currentAppId.value = row.id
  currentAppName.value = row.name
  jenkinsEnvDialogVisible.value = true
  await Promise.all([
    loadJenkinsEnvs(),
    loadJenkinsServers()
  ])
}

// 加载Jenkins环境列表
const loadJenkinsEnvs = async () => {
  try {
    const response = await appApi.getAppJenkinsEnvs(currentAppId.value)
    if (response.data?.code === 200) {
      jenkinsEnvs.value = response.data.data || []
    }
  } catch (error) {
    console.error('获取Jenkins环境列表失败:', error)
    ElMessage.error('获取Jenkins环境列表失败')
  }
}

// 添加Jenkins环境
const handleAddJenkinsEnv = () => {
  editingJenkinsEnv.value = null
  jenkinsEnvForm.value = {
    env_name: '',
    jenkins_server_id: null,
    job_name: '',
    deploy_config: {
      cluster_id: null,
      namespace: '',
      workload_type: 'deployment',
      workload_name: '',
      release_governance_enabled: false,
      pre_check_mode: 'observe'
    }
  }
  jenkinsJobs.value = []
  jobValidationStatus.value = null
  showJenkinsEnvForm.value = true
}

// 编辑Jenkins环境
const handleEditJenkinsEnv = (env) => {
  editingJenkinsEnv.value = env
  showJenkinsEnvForm.value = true
  jenkinsEnvForm.value = {
    env_name: env.env_name,
    jenkins_server_id: env.jenkins_server_id,
    job_name: env.job_name,
    deploy_config: {
      cluster_id: env.deploy_config?.cluster_id || null,
      namespace: env.deploy_config?.namespace || '',
      workload_type: env.deploy_config?.workload_type || 'deployment',
      workload_name: env.deploy_config?.workload_name || '',
      release_governance_enabled: Boolean(env.deploy_config?.release_governance_enabled),
      pre_check_mode: env.deploy_config?.pre_check_mode || 'observe'
    }
  }
  jobValidationStatus.value = null
}

// 保存Jenkins环境
const handleSaveJenkinsEnv = async () => {
  try {
    // 验证必填字段
    if (!jenkinsEnvForm.value.env_name?.trim()) {
      ElMessage.warning('请填写环境名称')
      return
    }

    // 如果选择了Jenkins服务器和任务名称，先验证任务是否存在
    const hasServerAndJob = jenkinsEnvForm.value.jenkins_server_id &&
                           jenkinsEnvForm.value.job_name?.trim()

    if (hasServerAndJob && (!jobValidationStatus.value || !jobValidationStatus.value.exists)) {
      try {
        await ElMessageBox.confirm(
          `Jenkins任务 "${jenkinsEnvForm.value.job_name.trim()}" 在所选服务器中可能不存在，确定要保存吗？`,
          '任务验证提醒',
          {
            confirmButtonText: '确定保存',
            cancelButtonText: '取消',
            type: 'warning'
          }
        )
      } catch (cancelError) {
        return // 用户取消操作
      }
    }

    // 准备提交数据，确保符合API规范
    const data = {
      env_name: jenkinsEnvForm.value.env_name.trim(),
      jenkins_server_id: jenkinsEnvForm.value.jenkins_server_id || null,
      job_name: jenkinsEnvForm.value.job_name?.trim() || '',
      deploy_config: {
        cluster_id: Number(jenkinsEnvForm.value.deploy_config?.cluster_id || 0) || 0,
        namespace: jenkinsEnvForm.value.deploy_config?.namespace?.trim() || '',
        workload_type: jenkinsEnvForm.value.deploy_config?.workload_type || '',
        workload_name: jenkinsEnvForm.value.deploy_config?.workload_name?.trim() || '',
        release_governance_enabled: Boolean(jenkinsEnvForm.value.deploy_config?.release_governance_enabled),
        pre_check_mode: jenkinsEnvForm.value.deploy_config?.pre_check_mode || 'observe'
      }
    }

    // 确保如果没有选择服务器，任务名称也为空
    if (!data.jenkins_server_id) {
      data.job_name = ''
    }

    if (!data.deploy_config.release_governance_enabled) {
      data.deploy_config = {
        cluster_id: 0,
        namespace: '',
        workload_type: '',
        workload_name: '',
        release_governance_enabled: false,
        pre_check_mode: 'observe'
      }
    } else if (!data.deploy_config.cluster_id || !data.deploy_config.namespace || !data.deploy_config.workload_type || !data.deploy_config.workload_name) {
      ElMessage.warning('启用发布治理时，请完整填写 K8s 集群、命名空间、工作负载类型和工作负载名称')
      return
    }

    if (editingJenkinsEnv.value) {
      // 更新环境
      await appApi.updateAppJenkinsEnv(currentAppId.value, editingJenkinsEnv.value.id, data)
      ElMessage.success('Jenkins环境更新成功')
    } else {
      // 新增环境
      await appApi.addAppJenkinsEnv(currentAppId.value, data)
      ElMessage.success('Jenkins环境创建成功')
    }

    resetJenkinsEnvForm()
    await loadJenkinsEnvs()
  } catch (error) {
    console.error('保存Jenkins环境失败:', error)
    const errorMessage = error.response?.data?.message || '保存Jenkins环境失败'
    ElMessage.error(errorMessage)
  }
}

// 删除Jenkins环境
const handleDeleteJenkinsEnv = async (env) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除环境 "${env.env_name}" 吗？`,
      '删除确认',
      {
        confirmButtonText: '删除',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    await appApi.deleteAppJenkinsEnv(currentAppId.value, env.id)
    ElMessage.success('Jenkins环境删除成功')
    await loadJenkinsEnvs()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除Jenkins环境失败:', error)
      ElMessage.error(error.response?.data?.message || '删除Jenkins环境失败')
    }
  }
}

// 重置Jenkins环境表单
const resetJenkinsEnvForm = () => {
  editingJenkinsEnv.value = null
  showJenkinsEnvForm.value = false
  jenkinsEnvForm.value = {
    env_name: '',
    jenkins_server_id: null,
    job_name: '',
    deploy_config: {
      cluster_id: null,
      namespace: '',
      workload_type: 'deployment',
      workload_name: '',
      release_governance_enabled: false,
      pre_check_mode: 'observe'
    }
  }
  jenkinsJobs.value = []
  jobValidationStatus.value = null
}

// 加载Jenkins服务器列表
const loadJenkinsServers = async () => {
  try {
    const response = await appApi.getJenkinsServers()
    if (response.data?.code === 200) {
      jenkinsServers.value = response.data.data || []
    }
  } catch (error) {
    console.error('获取Jenkins服务器列表失败:', error)
    ElMessage.error('获取Jenkins服务器列表失败')
  }
}

// Jenkins服务器选择变化时，重置任务名称和搜索任务
const handleServerChange = async (serverId) => {
  jenkinsEnvForm.value.job_name = ''
  jobValidationStatus.value = null
  if (serverId) {
    await searchJenkinsJobs('')
  } else {
    jenkinsJobs.value = []
  }
}

// 搜索Jenkins任务
const searchJenkinsJobs = async (searchKey = '') => {
  if (!jenkinsEnvForm.value.jenkins_server_id) {
    jenkinsJobs.value = []
    return
  }

  jobSearchLoading.value = true
  try {
    const response = await appApi.searchJenkinsJobs(
      jenkinsEnvForm.value.jenkins_server_id,
      searchKey
    )
    if (response.data?.code === 200) {
      jenkinsJobs.value = response.data.data || []
    }
  } catch (error) {
    console.error('搜索Jenkins任务失败:', error)
    ElMessage.error('搜索Jenkins任务失败')
  } finally {
    jobSearchLoading.value = false
  }
}

// 任务名称输入变化时的防抖搜索
let searchTimeout = null
const handleJobNameInput = (value) => {
  clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    if (jenkinsEnvForm.value.jenkins_server_id) {
      searchJenkinsJobs(value)
    }
  }, 500)
}

// 验证Jenkins任务是否存在
const validateJenkinsJob = async () => {
  if (!jenkinsEnvForm.value.jenkins_server_id || !jenkinsEnvForm.value.job_name) {
    jobValidationStatus.value = null
    return
  }

  jobValidationLoading.value = true
  try {
    const response = await appApi.validateJenkinsJob({
      jenkins_server_id: jenkinsEnvForm.value.jenkins_server_id,
      job_name: jenkinsEnvForm.value.job_name
    })

    if (response.data?.code === 200) {
      const result = response.data.data
      jobValidationStatus.value = {
        exists: result.exists,
        message: result.message,
        job_url: result.job_url
      }

      if (!result.exists) {
        ElMessage.warning(`任务验证失败: ${result.message}`)
      }
      // 成功时不显示消息提示，只显示图标
    }
  } catch (error) {
    console.error('验证Jenkins任务失败:', error)
    ElMessage.error('验证Jenkins任务失败')
    jobValidationStatus.value = { exists: false, message: '验证失败' }
  } finally {
    jobValidationLoading.value = false
  }
}

// 删除应用
const handleDeleteApp = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确认删除应用 "${row.name}"？删除后将无法恢复。`,
      '删除确认',
      {
        confirmButtonText: '确认删除',
        cancelButtonText: '取消',
        type: 'warning',
        confirmButtonClass: 'el-button--danger'
      }
    )

    const response = await appApi.deleteApplication(row.id)

    if (response.data?.code === 200) {
      ElMessage.success('应用删除成功')
      fetchApplicationList()
    } else {
      ElMessage.error(response.data?.message || '删除应用失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除应用失败:', error)
      ElMessage.error(error.response?.data?.message || '删除应用失败')
    }
  }
}

// 保存应用
const handleSaveApp = async () => {
  try {
    // 验证必填字段
    if (!applicationForm.value.name) {
      ElMessage.warning('请填写应用名称')
      return
    }
    if (!applicationForm.value.business_group_id) {
      ElMessage.warning('请选择业务组')
      return
    }
    if (!applicationForm.value.business_dept_id) {
      ElMessage.warning('请选择业务部门')
      return
    }

    // 准备提交数据，过滤空值
    const submitData = {}

    // 必填字段
    submitData.name = applicationForm.value.name
    submitData.business_group_id = applicationForm.value.business_group_id
    submitData.business_dept_id = applicationForm.value.business_dept_id

    // 可选字段，只包含非空值
    if (applicationForm.value.description) submitData.description = applicationForm.value.description
    if (applicationForm.value.repo_url) submitData.repo_url = applicationForm.value.repo_url
    if (applicationForm.value.programming_lang) submitData.programming_lang = applicationForm.value.programming_lang
    if (applicationForm.value.start_command) submitData.start_command = applicationForm.value.start_command
    if (applicationForm.value.stop_command) submitData.stop_command = applicationForm.value.stop_command
    if (applicationForm.value.health_api) submitData.health_api = applicationForm.value.health_api

    // 负责人数组，只包含非空数组
    if (applicationForm.value.dev_owners?.length > 0) submitData.dev_owners = applicationForm.value.dev_owners
    if (applicationForm.value.test_owners?.length > 0) submitData.test_owners = applicationForm.value.test_owners
    if (applicationForm.value.ops_owners?.length > 0) submitData.ops_owners = applicationForm.value.ops_owners

    // 关联资源，只包含非空数组
    if (applicationForm.value.domains?.length > 0) submitData.domains = applicationForm.value.domains
    if (applicationForm.value.hosts?.length > 0) submitData.hosts = applicationForm.value.hosts
    if (applicationForm.value.databases?.length > 0) submitData.databases = applicationForm.value.databases

    // 其他资源，只包含非空配置
    const otherRes = {}
    if (applicationForm.value.other_res.rabbitmq?.length > 0) otherRes.rabbitmq = applicationForm.value.other_res.rabbitmq
    if (applicationForm.value.other_res.zookeeper?.length > 0) otherRes.zookeeper = applicationForm.value.other_res.zookeeper
    if (applicationForm.value.other_res.kafka?.length > 0) otherRes.kafka = applicationForm.value.other_res.kafka
    if (applicationForm.value.other_res.redis?.length > 0) otherRes.redis = applicationForm.value.other_res.redis
    if (applicationForm.value.other_res.other?.length > 0) otherRes.other = applicationForm.value.other_res.other
    if (Object.keys(otherRes).length > 0) submitData.other_res = otherRes

    // Jenkins环境配置，只包含非空数组
    if (applicationForm.value.jenkins_envs?.length > 0) submitData.jenkins_envs = applicationForm.value.jenkins_envs

    let response
    if (editDialogVisible.value) {
      // 更新应用
      response = await appApi.updateApplication(currentApplication.value.id, submitData)
    } else {
      // 创建应用
      response = await appApi.createApplication(submitData)
    }

    if (response.data?.code === 200) {
      ElMessage.success(editDialogVisible.value ? '应用更新成功' : '应用创建成功')
      createDialogVisible.value = false
      editDialogVisible.value = false
      fetchApplicationList()
    } else {
      ElMessage.error(response.data?.message || '操作失败')
    }
  } catch (error) {
    console.error('保存应用失败:', error)
    ElMessage.error(error.response?.data?.message || '保存应用失败')
  }
}


// 重置表单
const resetForm = () => {
  applicationForm.value = {
    // 必填参数
    name: '',
    business_group_id: null,
    business_dept_id: null,

    // 基本信息（可选）
    description: '',
    repo_url: '',

    // 负责人信息（可选）
    dev_owners: [],
    test_owners: [],
    ops_owners: [],

    // 技术信息（可选）
    programming_lang: '',
    start_command: '',
    stop_command: '',
    health_api: '',

    // 关联资源（可选）
    domains: [],
    hosts: [],
    databases: [],
    other_res: {
      rabbitmq: [],
      zookeeper: [],
      kafka: [],
      redis: [],
      other: []
    },

    // Jenkins环境配置（可选）
    jenkins_envs: []
  }
}

// 刷新数据
const handleRefresh = () => {
  fetchApplicationList()
}

// 搜索数据
const handleSearch = () => {
  pagination.page = 1
  fetchApplicationList()
}

// 重置搜索
const resetSearch = () => {
  searchKeyword.value = ''
  pagination.page = 1
  fetchApplicationList()
}

// 快速发布
const handleQuickRelease = () => {
  router.push('/app/quick-release')
}

const openExternalLink = (link) => {
  if (!link) {
    ElMessage.warning('当前没有可打开的链接')
    return
  }
  window.open(link, '_blank', 'noopener,noreferrer')
}

const handleViewDeploymentDetail = (deployment) => {
  if (!deployment?.deployment_id) {
    ElMessage.warning('当前没有可查看的发布记录')
    return
  }
  router.push(`/app/quick-temp/${deployment.deployment_id}`)
}

const handleOpenGlobalSearch = (app) => {
  const keyword = app?.name || app?.code || ''
  router.push({
    path: '/search/global',
    query: keyword ? {
      keyword,
      types: 'application,alert'
    } : {}
  })
}

const handleOpenAlertCenter = () => {
  router.push('/monitor/alert-center')
}

const openInternalPath = (path) => {
  if (!path) {
    ElMessage.warning('当前没有可打开的联动页面')
    return
  }
  router.push(path)
}

const getGovernanceTagType = (governance) => getGovernanceDisplayTagType(governance)

const getGovernanceText = (governance) => getGovernanceDisplayText(governance)

const getGovernancePath = (governance, target) => getGovernanceSuggestionPath(governance, target)

const getGovernanceTargetText = (env) => {
  const config = env?.deploy_config || {}
  if (!config.cluster_id || !config.namespace || !config.workload_type || !config.workload_name) {
    return '-'
  }
  const clusterName = env?.release_governance?.cluster_name || `集群${config.cluster_id}`
  return `${clusterName} / ${config.namespace} / ${config.workload_type} / ${config.workload_name}`
}

const fetchK8sClusters = async () => {
  try {
    const response = await k8sApi.getClusterList()
    const responseData = response.data || response
    if (responseData.code === 200 || responseData.success) {
      const clusters = responseData.data?.list || responseData.data || []
      k8sClusters.value = clusters.map(cluster => ({
        id: cluster.id,
        name: cluster.name
      }))
    }
  } catch (error) {
    console.error('获取 K8s 集群列表失败:', error)
  }
}

// 分页变化
const handlePageChange = (page) => {
  pagination.page = page
  fetchApplicationList()
}

const handleSizeChange = (size) => {
  pagination.pageSize = size
  pagination.page = 1
  fetchApplicationList()
}

// 获取状态标签类型
const getStatusTagType = (status) => {
  const statusMap = {
    1: 'warning',
    2: 'success',
    3: 'danger',
    0: 'info'
  }
  return statusMap[status] || 'info'
}

// 获取状态文本
const getStatusText = (status) => {
  const statusMap = {
    1: '未激活',
    2: '已激活',
    3: '异常',
    0: '未知'
  }
  return statusMap[status] || '未知'
}

const getDeploymentStatusTagType = (status) => {
  const statusMap = {
    1: 'info',
    2: 'warning',
    3: 'success',
    4: 'danger',
    5: 'info'
  }
  return statusMap[status] || 'info'
}

const formatDateTime = (value) => {
  if (!value) return '-'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return value
  return date.toLocaleString()
}

const formatDuration = (seconds) => {
  if (!seconds && seconds !== 0) return '-'
  if (seconds < 60) return `${seconds}s`
  const minutes = Math.floor(seconds / 60)
  const remainSeconds = seconds % 60
  if (minutes < 60) return `${minutes}m ${remainSeconds}s`
  const hours = Math.floor(minutes / 60)
  const remainMinutes = minutes % 60
  return `${hours}h ${remainMinutes}m ${remainSeconds}s`
}

const getEnvironmentCoverageText = (row) => {
  const configured = row?.configured_environment_count || 0
  const total = row?.total_environment_count || 0
  return `${configured}/${total}`
}

const deploymentSuccessRate = computed(() => {
  const total = overviewSummary.value.total_deployments || 0
  if (!total) return 0
  return Math.round(((overviewSummary.value.successful_deployments || 0) / total) * 100)
})

// 获取业务组列表
const fetchBusinessGroups = async () => {
  try {
    const response = await cmdbApi.getAllCmdbGroups()
    if (response.data?.code === 200) {
      businessGroups.value = response.data.data || []
    }
  } catch (error) {
    console.error('获取业务组列表失败:', error)
  }
}

// 生成级联选择器的业务组数据：只能选择二级分组
const cascaderBusinessGroups = computed(() => {
  return businessGroups.value.map(group => ({
    value: group.id,
    label: group.name,
    leaf: false, // 一级分组不是叶子节点，不可选中
    children: group.children?.filter(child => !child.children || child.children.length === 0).map(child => ({
      value: child.id,
      label: child.name,
      leaf: true // 二级分组是叶子节点，可以选中
    })) || []
  }))
})

// 获取业务部门列表
const fetchBusinessDepts = async () => {
  try {
    // 使用 queryDeptList 获取部门数据
    const response = await systemApi.queryDeptList()
    if (response.data?.code === 200) {
      businessDepts.value = (response.data.data || []).filter(dept => dept.deptName)
    } else {
      // 如果第一个API失败，尝试第二个
      try {
        const response2 = await systemApi.querySysDeptVoList()
        if (response2.data?.code === 200) {
          businessDepts.value = (response2.data.data || []).filter(dept => dept.deptName)
        }
      } catch (error2) {
        console.error('获取业务部门列表失败:', error2)
      }
    }
  } catch (error) {
    console.error('获取业务部门列表失败:', error)
  }
}

// 获取用户列表（负责人选项）
const fetchUserOptions = async () => {
  try {
    const response = await systemApi.queryAdminList({
      pageNum: 1,
      pageSize: 1000, // 获取所有用户作为负责人选项
      status: '1' // 只获取启用状态的用户 (1启用，2停用)
    })
    if (response.data?.code === 200) {
      userOptions.value = response.data.data?.list?.map(user => ({
        id: user.id,
        name: user.nickname || user.username
      })).filter(user => user.name) || [] // 只保留有名称的用户
    }
  } catch (error) {
    console.error('获取用户列表失败:', error)
  }
}

// ID映射计算属性
const getDeptNameById = (deptId) => {
  if (!deptId) return '-'
  const dept = businessDepts.value.find(d => d.id === deptId)
  return dept ? dept.deptName : `部门ID: ${deptId}`
}

const getGroupNameById = (groupId) => {
  if (!groupId) return '-'

  // 在业务组层次结构中查找指定ID的分组
  for (const group of businessGroups.value) {
    if (group.children && group.children.length > 0) {
      const subGroup = group.children.find(sg => sg.id === groupId)
      if (subGroup) {
        return subGroup.name
      }
    }
  }

  return `业务组ID: ${groupId}`
}

const getUserNamesByIds = (userIds) => {
  if (!userIds || !Array.isArray(userIds) || userIds.length === 0) return '-'
  const userNames = userIds.map(id => {
    const user = userOptions.value.find(u => u.id === id)
    return user ? user.name : `用户ID: ${id}`
  })
  return userNames.join(', ')
}

const getServerNameById = (serverId) => {
  if (!serverId) return '未配置'
  const server = jenkinsServers.value.find(s => s.id === serverId)
  return server ? server.name : `服务器${serverId}`
}

const syncRoutePreset = () => {
  routePreset.appId = Number(route.query.appId || 0)
  routePreset.env = String(route.query.env || '').trim()
  routePreset.focus = String(route.query.focus || '').trim().toLowerCase()
}

const getEnvRowClassName = ({ row }) => {
  if (!currentFocusedEnv.value || routePreset.focus !== 'governance') return ''
  return String(row?.env_name || '').trim().toLowerCase() === String(currentFocusedEnv.value?.env_name || '').trim().toLowerCase()
    ? 'focused-env-row'
    : ''
}

const applyRoutePreset = async () => {
  syncRoutePreset()
  if (!routePreset.appId) return

  let target = filteredApplicationList.value.find(item => Number(item.id) === routePreset.appId)
  if (!target) {
    try {
      const response = await appApi.getApplicationDetail(routePreset.appId)
      if (response.data?.code === 200) {
        target = response.data.data || null
      }
    } catch (error) {
      console.error('按路由预设加载应用失败:', error)
    }
  }

  if (!target) {
    ElMessage.warning(`未找到应用 #${routePreset.appId}`)
    return
  }
  await handleViewApp(target)
}

// 页面初始化
onMounted(() => {
  fetchApplicationList()
  fetchBusinessGroups()
  fetchBusinessDepts()
  fetchUserOptions()
  fetchK8sClusters()
  applyRoutePreset()
})

watch(filteredApplicationList, () => {
  applyApplicationTableLayout()
}, { deep: true, flush: 'post' })

watch(() => route.query, () => {
  applyRoutePreset()
}, { deep: true })
</script>

<template>
  <div class="application-management">
    <el-card shadow="hover" class="main-card">
      <template #header>
        <div class="card-header">
          <span class="title">应用管理</span>
          <div class="header-actions">
            <el-button type="primary" :icon="Plus"  v-authority="['app:application:add']" size="small" @click="handleCreateApp">
              创建应用
            </el-button>
          </div>
        </div>
      </template>

      <!-- 搜索卡片 -->
      <el-card shadow="hover" class="search-card">
        <div class="search-toolbar">
          <el-form :inline="true" class="search-form">
            <el-form-item>
              <el-input
                v-model="searchKeyword"
                placeholder="搜索应用名称、编码或技术栈"
                clearable
                size="small"
                style="width: 300px"
                :prefix-icon="Search"
                @keyup.enter="handleSearch"
              />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" :icon="Search" size="small" @click="handleSearch">
                搜索
              </el-button>
              <el-button :icon="Close" size="small" @click="resetSearch">
                重置
              </el-button>
              <el-button :icon="Refresh" size="small" @click="handleRefresh">
                刷新
              </el-button>
              <el-button type="success" :icon="Promotion" size="small" @click="handleQuickRelease">
                快速发布
              </el-button>
            </el-form-item>
          </el-form>
          <div class="toolbar-right">
            <span class="resource-count">共 {{ filteredApplicationList.length }} 个应用</span>
          </div>
        </div>
      </el-card>

      <div class="overview-grid">
        <div class="overview-card primary">
          <div class="overview-label">应用总数</div>
          <div class="overview-value">{{ overviewSummary.total_applications }}</div>
          <div class="overview-sub">当前筛选条件下的应用规模</div>
        </div>
        <div class="overview-card success">
          <div class="overview-label">已激活应用</div>
          <div class="overview-value">{{ overviewSummary.active_applications }}</div>
          <div class="overview-sub">未激活 {{ overviewSummary.inactive_applications }}，异常 {{ overviewSummary.abnormal_applications }}</div>
        </div>
        <div class="overview-card warning">
          <div class="overview-label">环境配置覆盖</div>
          <div class="overview-value">{{ overviewSummary.configured_environments }}/{{ overviewSummary.total_environments }}</div>
          <div class="overview-sub">Jenkins 服务器与任务名称都已配置的环境数</div>
        </div>
        <div class="overview-card danger">
          <div class="overview-label">发布成功率</div>
          <div class="overview-value">{{ deploymentSuccessRate }}%</div>
          <div class="overview-sub">成功 {{ overviewSummary.successful_deployments }}，失败 {{ overviewSummary.failed_deployments }}</div>
        </div>
      </div>

      <!-- 应用列表卡片 -->
      <el-card shadow="never" class="table-card">
        <el-table
          :data="filteredApplicationList"
          v-loading="loading"
          element-loading-text="加载中..."
          class="app-table"
          empty-text="暂无应用数据"
          stripe
        >
          <el-table-column prop="name" label="应用名称" min-width="150" fixed="left">
            <template #default="{ row }">
              <div class="app-name">
                <el-icon class="app-icon"><Monitor /></el-icon>
                <span class="app-name-link" @click="handleViewApp(row)">{{ row.name }}</span>
              </div>
            </template>
          </el-table-column>


          <el-table-column prop="business_group_id" label="业务线" width="150" show-overflow-tooltip>
            <template #default="{ row }">
              <el-tag size="small" type="primary">
                {{ row.business_group_name || getGroupNameById(row.business_group_id) }}
              </el-tag>
            </template>
          </el-table-column>

          <el-table-column label="环境配置" width="130">
            <template #default="{ row }">
              <div class="env-coverage-cell">
                <span class="coverage-value">{{ getEnvironmentCoverageText(row) }}</span>
                <span class="coverage-sub">已配置/总环境</span>
              </div>
            </template>
          </el-table-column>

          <el-table-column prop="dev_owners" label="研发负责人" min-width="150" show-overflow-tooltip>
            <template #default="{ row }">
              <span style="font-size: 12px;">{{ getUserNamesByIds(row.dev_owners) }}</span>
            </template>
          </el-table-column>

          <el-table-column prop="programming_lang" label="技术栈" width="100">
            <template #default="{ row }">
              <el-tag :type="row.programming_lang === 'Java' ? 'warning' : 'primary'" size="small">
                {{ row.programming_lang || '-' }}
              </el-tag>
            </template>
          </el-table-column>

          <el-table-column prop="status" label="状态" width="100">
            <template #default="{ row }">
              <el-tag :type="getStatusTagType(row.status)" size="small">
                {{ getStatusText(row.status) }}
              </el-tag>
            </template>
          </el-table-column>

          <el-table-column label="最近发布" width="190" show-overflow-tooltip>
            <template #default="{ row }">
              <div v-if="row.latest_deployment" class="latest-deployment-cell">
                <div class="latest-title" @click="handleViewDeploymentDetail(row.latest_deployment)">
                  {{ row.latest_deployment.title }}
                </div>
                <div class="latest-meta">
                  <el-tag :type="getDeploymentStatusTagType(row.latest_deployment.status)" size="small">
                    {{ row.latest_deployment.status_text }}
                  </el-tag>
                  <span>{{ row.latest_deployment.environment || '-' }}</span>
                </div>
                <div class="latest-meta muted">
                  <span>{{ row.latest_deployment.creator_name || '-' }}</span>
                  <span>{{ formatDateTime(row.latest_deployment.created_at) }}</span>
                </div>
              </div>
              <span v-else class="muted">暂无发布记录</span>
            </template>
          </el-table-column>

          <el-table-column label="操作" width="200" fixed="right">
            <template #default="{ row }">
              <div class="operation-buttons">
                <el-tooltip content="编辑" placement="top">
                  <el-button :icon="Edit" v-authority="['app:application:edit']" size="small" type="primary" circle @click="handleEditApp(row)" />
                </el-tooltip>
                <el-tooltip content="环境配置" placement="top">
                  <el-button :icon="Setting" v-authority="['app:application:env']" size="small" type="warning" circle @click="handleJenkinsEnvConfig(row)" />
                </el-tooltip>
                <el-tooltip content="删除" placement="top">
                  <el-button :icon="Delete" v-authority="['app:application:delete']" size="small" type="danger" circle @click="handleDeleteApp(row)" />
                </el-tooltip>
              </div>
            </template>
          </el-table-column>
        </el-table>

        <!-- 分页 -->
        <div class="pagination-wrapper">
          <el-pagination
            v-model:current-page="pagination.page"
            v-model:page-size="pagination.pageSize"
            :page-sizes="[10, 20, 50, 100]"
            :total="pagination.total"
            layout="total, sizes, prev, pager, next, jumper"
            @size-change="handleSizeChange"
            @current-change="handlePageChange"
          />
        </div>
      </el-card>
    </el-card>

    <!-- 创建应用对话框 -->
    <el-dialog
      v-model="createDialogVisible"
      title="创建应用"
      width="700px"
      class="create-dialog compact-dialog"
    >
      <el-form :model="applicationForm" label-width="90px" class="compact-form">
        <!-- 基础配置行 -->
        <el-form-item label="应用名称" required class="form-item-full">
          <el-input
            v-model="applicationForm.name"
            placeholder="请输入应用名称"
            size="small"
          />
        </el-form-item>

        <!-- 业务信息行 -->
        <div class="form-row">
          <el-form-item label="业务组" required class="form-item-half">
            <el-cascader
              v-model="applicationForm.business_group_id"
              :options="cascaderBusinessGroups"
              placeholder="请选择业务组"
              style="width: 100%"
              size="small"
              :props="{
                expandTrigger: 'hover',
                checkStrictly: true,
                emitPath: false,
                leaf: 'leaf'
              }"
            />
          </el-form-item>
          <el-form-item label="业务部门" required class="form-item-half">
            <el-select
              v-model="applicationForm.business_dept_id"
              placeholder="请选择业务部门"
              style="width: 100%"
              size="small"
              filterable
              clearable
            >
              <el-option
                v-for="dept in businessDepts"
                :key="dept.id"
                :label="dept.deptName"
                :value="dept.id"
              />
            </el-select>
          </el-form-item>
        </div>

        <!-- 描述和仓库地址 -->
        <el-form-item label="应用描述" class="form-item-full">
          <el-input
            v-model="applicationForm.description"
            type="textarea"
            :rows="2"
            placeholder="请输入应用描述"
            size="small"
          />
        </el-form-item>

        <el-form-item label="仓库地址" class="form-item-full">
          <el-input
            v-model="applicationForm.repo_url"
            placeholder="https://git.company.com/project.git"
            size="small"
          />
        </el-form-item>

        <!-- 技术配置行 -->
        <div class="form-row">
          <el-form-item label="编程语言" class="form-item-half">
            <el-select
              v-model="applicationForm.programming_lang"
              placeholder="请选择编程语言"
              style="width: 100%"
              size="small"
            >
              <el-option
                v-for="lang in programmingLanguages"
                :key="lang"
                :label="lang"
                :value="lang"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="健康检查" class="form-item-half">
            <el-input
              v-model="applicationForm.health_api"
              placeholder="/api/v1/health"
              size="small"
            />
          </el-form-item>
        </div>

        <!-- 命令配置行 -->
        <div class="form-row">
          <el-form-item label="启动命令" class="form-item-half">
            <el-input
              v-model="applicationForm.start_command"
              placeholder="./app --config=config.yaml"
              size="small"
            />
          </el-form-item>
          <el-form-item label="停止命令" class="form-item-half">
            <el-input
              v-model="applicationForm.stop_command"
              placeholder="kill -TERM $PID"
              size="small"
            />
          </el-form-item>
        </div>

        <!-- 负责人配置 -->
        <div class="form-row">
          <el-form-item label="研发负责人" class="form-item-half">
            <el-select
              v-model="applicationForm.dev_owners"
              placeholder="请选择研发负责人"
              style="width: 100%"
              size="small"
              multiple
              collapse-tags
              collapse-tags-tooltip
              filterable
              clearable
            >
              <el-option
                v-for="user in userOptions"
                :key="user.id"
                :label="user.name"
                :value="user.id"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="测试负责人" class="form-item-half">
            <el-select
              v-model="applicationForm.test_owners"
              placeholder="请选择测试负责人"
              style="width: 100%"
              size="small"
              multiple
              collapse-tags
              collapse-tags-tooltip
              filterable
              clearable
            >
              <el-option
                v-for="user in userOptions"
                :key="user.id"
                :label="user.name"
                :value="user.id"
              />
            </el-select>
          </el-form-item>
        </div>

        <el-form-item label="运维负责人" class="form-item-full">
          <el-select
            v-model="applicationForm.ops_owners"
            placeholder="请选择运维负责人"
            style="width: 100%"
            size="small"
            multiple
            collapse-tags
            collapse-tags-tooltip
            filterable
            clearable
          >
            <el-option
              v-for="user in userOptions"
              :key="user.id"
              :label="user.name || `用户${user.id}`"
              :value="user.id"
            />
          </el-select>
        </el-form-item>
      </el-form>

      <template #footer>
        <div class="dialog-footer">
          <el-button @click="createDialogVisible = false" size="small">取消</el-button>
          <el-button type="primary" @click="handleSaveApp" size="small">确认创建</el-button>
        </div>
      </template>
    </el-dialog>

    <!-- 编辑应用对话框 -->
    <el-dialog
      v-model="editDialogVisible"
      title="编辑应用"
      width="700px"
      class="edit-dialog compact-dialog"
    >
      <el-form :model="applicationForm" label-width="90px" class="compact-form">
        <!-- 基础配置行 -->
        <el-form-item label="应用名称" required class="form-item-full">
          <el-input
            v-model="applicationForm.name"
            placeholder="请输入应用名称"
            size="small"
          />
        </el-form-item>

        <!-- 业务信息行 -->
        <div class="form-row">
          <el-form-item label="业务组" required class="form-item-half">
            <el-cascader
              v-model="applicationForm.business_group_id"
              :options="cascaderBusinessGroups"
              placeholder="请选择业务组"
              style="width: 100%"
              size="small"
              :props="{
                expandTrigger: 'hover',
                checkStrictly: true,
                emitPath: false,
                leaf: 'leaf'
              }"
            />
          </el-form-item>
          <el-form-item label="业务部门" required class="form-item-half">
            <el-select
              v-model="applicationForm.business_dept_id"
              placeholder="请选择业务部门"
              style="width: 100%"
              size="small"
              filterable
              clearable
            >
              <el-option
                v-for="dept in businessDepts"
                :key="dept.id"
                :label="dept.deptName"
                :value="dept.id"
              />
            </el-select>
          </el-form-item>
        </div>

        <!-- 描述和仓库地址 -->
        <el-form-item label="应用描述" class="form-item-full">
          <el-input
            v-model="applicationForm.description"
            type="textarea"
            :rows="2"
            placeholder="请输入应用描述"
            size="small"
          />
        </el-form-item>

        <el-form-item label="仓库地址" class="form-item-full">
          <el-input
            v-model="applicationForm.repo_url"
            placeholder="https://git.company.com/project.git"
            size="small"
          />
        </el-form-item>

        <!-- 技术配置行 -->
        <div class="form-row">
          <el-form-item label="编程语言" class="form-item-half">
            <el-select
              v-model="applicationForm.programming_lang"
              placeholder="请选择编程语言"
              style="width: 100%"
              size="small"
            >
              <el-option
                v-for="lang in programmingLanguages"
                :key="lang"
                :label="lang"
                :value="lang"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="健康检查" class="form-item-half">
            <el-input
              v-model="applicationForm.health_api"
              placeholder="/api/v1/health"
              size="small"
            />
          </el-form-item>
        </div>

        <!-- 命令配置行 -->
        <div class="form-row">
          <el-form-item label="启动命令" class="form-item-half">
            <el-input
              v-model="applicationForm.start_command"
              placeholder="./app --config=config.yaml"
              size="small"
            />
          </el-form-item>
          <el-form-item label="停止命令" class="form-item-half">
            <el-input
              v-model="applicationForm.stop_command"
              placeholder="kill -TERM $PID"
              size="small"
            />
          </el-form-item>
        </div>

        <!-- 负责人配置 -->
        <div class="form-row">
          <el-form-item label="研发负责人" class="form-item-half">
            <el-select
              v-model="applicationForm.dev_owners"
              placeholder="请选择研发负责人"
              style="width: 100%"
              size="small"
              multiple
              collapse-tags
              collapse-tags-tooltip
              filterable
              clearable
            >
              <el-option
                v-for="user in userOptions"
                :key="user.id"
                :label="user.name"
                :value="user.id"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="测试负责人" class="form-item-half">
            <el-select
              v-model="applicationForm.test_owners"
              placeholder="请选择测试负责人"
              style="width: 100%"
              size="small"
              multiple
              collapse-tags
              collapse-tags-tooltip
              filterable
              clearable
            >
              <el-option
                v-for="user in userOptions"
                :key="user.id"
                :label="user.name"
                :value="user.id"
              />
            </el-select>
          </el-form-item>
        </div>

        <el-form-item label="运维负责人" class="form-item-full">
          <el-select
            v-model="applicationForm.ops_owners"
            placeholder="请选择运维负责人"
            style="width: 100%"
            size="small"
            multiple
            collapse-tags
            collapse-tags-tooltip
            filterable
            clearable
          >
            <el-option
              v-for="user in userOptions"
              :key="user.id"
              :label="user.name"
              :value="user.id"
            />
          </el-select>
        </el-form-item>
      </el-form>

      <template #footer>
        <div class="dialog-footer">
          <el-button @click="editDialogVisible = false" size="small">取消</el-button>
          <el-button type="primary" @click="handleSaveApp" size="small">确认更新</el-button>
        </div>
      </template>
    </el-dialog>

    <!-- 应用详情对话框 -->
    <el-dialog
      v-model="detailDialogVisible"
      title="应用详情"
      width="1100px"
      class="detail-dialog compact-dialog"
    >
      <div v-loading="detailLoading" class="detail-panel">
        <div class="detail-hero">
          <div>
            <div class="detail-title-row">
              <span class="detail-app-name">{{ currentApplication.name || '-' }}</span>
              <el-tag :type="getStatusTagType(currentApplication.status)" size="small">
                {{ getStatusText(currentApplication.status) }}
              </el-tag>
            </div>
            <div class="detail-app-meta">
              <span>编码：{{ currentApplication.code || '-' }}</span>
              <span>业务组：{{ currentApplication.business_group_name || getGroupNameById(currentApplication.business_group_id) }}</span>
              <span>部门：{{ currentApplication.business_dept_name || getDeptNameById(currentApplication.business_dept_id) }}</span>
            </div>
          </div>
          <div class="detail-actions-row">
            <el-button size="small" @click="handleOpenGlobalSearch(currentApplication)">全局搜索</el-button>
            <el-button size="small" @click="handleOpenAlertCenter">告警中心</el-button>
            <el-button size="small" type="primary" @click="handleQuickRelease">快速发布</el-button>
          </div>
        </div>

        <div class="detail-grid">
          <div class="detail-card">
            <div class="detail-card-title">基础信息</div>
            <div class="detail-item">
              <span class="detail-label">应用描述</span>
              <span class="detail-value">{{ currentApplication.description || '-' }}</span>
            </div>
            <div class="detail-item">
              <span class="detail-label">仓库地址</span>
              <span class="detail-value detail-link" @click="openExternalLink(currentApplication.repo_url)">
                {{ currentApplication.repo_url || '-' }}
              </span>
            </div>
            <div class="detail-item">
              <span class="detail-label">技术栈</span>
              <span class="detail-value">{{ currentApplication.programming_lang || '-' }}</span>
            </div>
            <div class="detail-item">
              <span class="detail-label">健康检查</span>
              <span class="detail-value">{{ currentApplication.health_api || '-' }}</span>
            </div>
            <div class="detail-item">
              <span class="detail-label">创建时间</span>
              <span class="detail-value">{{ formatDateTime(currentApplication.created_at) }}</span>
            </div>
          </div>

          <div class="detail-card">
            <div class="detail-card-title">责任人与发布</div>
            <div class="detail-item">
              <span class="detail-label">研发负责人</span>
              <span class="detail-value">{{ getUserNamesByIds(currentApplication.dev_owners) }}</span>
            </div>
            <div class="detail-item">
              <span class="detail-label">测试负责人</span>
              <span class="detail-value">{{ getUserNamesByIds(currentApplication.test_owners) }}</span>
            </div>
            <div class="detail-item">
              <span class="detail-label">运维负责人</span>
              <span class="detail-value">{{ getUserNamesByIds(currentApplication.ops_owners) }}</span>
            </div>
            <div class="detail-item">
              <span class="detail-label">环境覆盖</span>
              <span class="detail-value">{{ getEnvironmentCoverageText(currentApplication) }}</span>
            </div>
            <div class="detail-item">
              <span class="detail-label">发布统计</span>
              <span class="detail-value">
                {{ currentApplication.deployment_stats?.total || 0 }} 次
                / 成功 {{ currentApplication.deployment_stats?.success || 0 }}
                / 失败 {{ currentApplication.deployment_stats?.failed || 0 }}
              </span>
            </div>
          </div>
        </div>

        <el-alert
          v-if="routePreset.focus === 'governance' && currentFocusedEnv"
          type="info"
          :closable="false"
          class="governance-focus-alert"
        >
          <template #title>
            <div class="governance-focus-title">
              <span>当前已按治理入口打开：{{ currentFocusedEnv.env_name }} 环境</span>
              <el-tag :type="getGovernanceTagType(currentFocusedEnv.release_governance)" size="small">
                {{ getGovernanceText(currentFocusedEnv.release_governance) }}
              </el-tag>
            </div>
          </template>
            <div class="governance-focus-actions">
              <span class="governance-focus-target">{{ getGovernanceTargetText(currentFocusedEnv) }}</span>
              <div class="detail-table-actions">
              <el-button
                v-if="getGovernancePath(currentFocusedEnv.release_governance, 'k8sWorkload')"
                size="small"
                text
                type="primary"
                @click="openInternalPath(getGovernancePath(currentFocusedEnv.release_governance, 'k8sWorkload'))"
              >K8s 治理</el-button>
              <el-button
                v-if="getGovernancePath(currentFocusedEnv.release_governance, 'alertCenter')"
                size="small"
                text
                type="primary"
                @click="openInternalPath(getGovernancePath(currentFocusedEnv.release_governance, 'alertCenter'))"
              >告警中心</el-button>
              <el-button
                v-if="getGovernancePath(currentFocusedEnv.release_governance, 'aiDiagnosis')"
                size="small"
                text
                type="primary"
                @click="openInternalPath(getGovernancePath(currentFocusedEnv.release_governance, 'aiDiagnosis'))"
              >AI 诊断</el-button>
            </div>
          </div>
        </el-alert>

        <div class="detail-card full-width">
          <div class="detail-card-header">
            <span class="detail-card-title">环境配置</span>
            <el-button size="small" type="primary" @click="handleJenkinsEnvConfig(currentApplication)">管理环境</el-button>
          </div>
          <el-table :data="currentApplication.jenkins_envs || []" size="small" empty-text="暂无环境配置" :row-class-name="getEnvRowClassName">
            <el-table-column prop="env_name" label="环境" width="100" />
            <el-table-column prop="server_name" label="Jenkins服务器" min-width="180">
              <template #default="{ row }">
                <span>{{ row.server_name || getServerNameById(row.jenkins_server_id) }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="job_name" label="任务名称" min-width="180">
              <template #default="{ row }">
                <span>{{ row.job_name || '-' }}</span>
              </template>
            </el-table-column>
            <el-table-column label="任务地址" min-width="220" show-overflow-tooltip>
              <template #default="{ row }">
                <span
                  class="detail-link"
                  @click="openExternalLink(row.jenkins_job_url)"
                >
                  {{ row.jenkins_job_url || '-' }}
                </span>
              </template>
            </el-table-column>
            <el-table-column label="配置状态" width="110">
              <template #default="{ row }">
                <el-tag :type="row.is_configured ? 'success' : 'warning'" size="small">
                  {{ row.is_configured ? '已配置' : '待配置' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="治理目标" min-width="260" show-overflow-tooltip>
              <template #default="{ row }">
                <span>{{ getGovernanceTargetText(row) }}</span>
              </template>
            </el-table-column>
            <el-table-column label="治理状态" width="160">
              <template #default="{ row }">
                <el-tag :type="getGovernanceTagType(row.release_governance)" size="small">
                  {{ getGovernanceText(row.release_governance) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="治理联动" width="220" fixed="right">
              <template #default="{ row }">
                <div class="detail-table-actions">
                  <el-button
                    v-if="getGovernancePath(row.release_governance, 'k8sWorkload')"
                    size="small"
                    text
                    type="primary"
                    @click="openInternalPath(getGovernancePath(row.release_governance, 'k8sWorkload'))"
                  >K8s</el-button>
                  <el-button
                    v-if="getGovernancePath(row.release_governance, 'alertCenter')"
                    size="small"
                    text
                    type="primary"
                    @click="openInternalPath(getGovernancePath(row.release_governance, 'alertCenter'))"
                  >告警</el-button>
                  <el-button
                    v-if="getGovernancePath(row.release_governance, 'aiDiagnosis')"
                    size="small"
                    text
                    type="primary"
                    @click="openInternalPath(getGovernancePath(row.release_governance, 'aiDiagnosis'))"
                  >AI</el-button>
                </div>
              </template>
            </el-table-column>
          </el-table>
        </div>

        <div class="detail-card full-width">
          <div class="detail-card-header">
            <span class="detail-card-title">最近发布记录</span>
            <span class="detail-card-subtitle">与现有发布记录、日志和 Jenkins 构建入口联动</span>
          </div>
          <el-table :data="currentApplication.recent_deployments || []" size="small" empty-text="暂无发布记录">
            <el-table-column prop="title" label="发布标题" min-width="180" show-overflow-tooltip />
            <el-table-column prop="environment" label="环境" width="100" />
            <el-table-column prop="creator_name" label="创建人" width="120" />
            <el-table-column label="状态" width="100">
              <template #default="{ row }">
                <el-tag :type="getDeploymentStatusTagType(row.status)" size="small">
                  {{ row.status_text }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="build_number" label="构建号" width="100" />
            <el-table-column label="发布时间" width="180">
              <template #default="{ row }">
                {{ formatDateTime(row.created_at) }}
              </template>
            </el-table-column>
            <el-table-column label="耗时" width="100">
              <template #default="{ row }">
                {{ formatDuration(row.duration) }}
              </template>
            </el-table-column>
            <el-table-column label="操作" width="220" fixed="right">
              <template #default="{ row }">
                <div class="detail-table-actions">
                  <el-button size="small" text type="primary" @click="handleViewDeploymentDetail(row)">发布详情</el-button>
                  <el-button size="small" text type="primary" @click="openExternalLink(row.log_url)">日志</el-button>
                  <el-button size="small" text type="primary" @click="openExternalLink(row.jenkins_job_url)">Jenkins</el-button>
                </div>
              </template>
            </el-table-column>
          </el-table>
        </div>
      </div>
    </el-dialog>


    <!-- Jenkins环境配置对话框 -->
    <el-dialog
      v-model="jenkinsEnvDialogVisible"
      :title="`${currentAppName} - Jenkins环境配置`"
      width="900px"
      class="jenkins-env-dialog modern-dialog"
    >
      <el-card shadow="hover" class="env-main-card">
        <!-- 环境列表 -->
        <div class="env-list-section">
          <div class="section-header">
            <div class="header-left">
              <el-icon class="section-icon"><Monitor /></el-icon>
              <span class="section-title">环境列表</span>
            </div>
            <el-button v-authority="['app:application:envadd']" type="primary" size="small" @click="handleAddJenkinsEnv" class="add-btn">
              <el-icon><Plus /></el-icon>
              新增环境
            </el-button>
          </div>

          <el-table :data="jenkinsEnvs" size="small" class="env-table">
            <el-table-column prop="env_name" label="环境名称" width="120">
              <template #default="{ row }">
                <el-tag type="info" size="small" class="env-tag">{{ row.env_name }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="jenkins_server_id" label="Jenkins服务器" width="150">
              <template #default="{ row }">
                <span class="server-name">{{ row.server_name || getServerNameById(row.jenkins_server_id) }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="job_name" label="任务名称" min-width="150">
              <template #default="{ row }">
                <span class="job-name">{{ row.job_name || '未配置' }}</span>
              </template>
            </el-table-column>
            <el-table-column label="任务地址" min-width="180" show-overflow-tooltip>
              <template #default="{ row }">
                <span class="detail-link" @click="openExternalLink(row.jenkins_job_url)">
                  {{ row.jenkins_job_url || '-' }}
                </span>
              </template>
            </el-table-column>
            <el-table-column label="配置状态" width="100">
              <template #default="{ row }">
                <el-tag
                  :type="row.is_configured ? 'success' : 'warning'"
                  size="small"
                  class="status-tag"
                >
                  <el-icon>
                    <component :is="row.is_configured ? SuccessFilled : WarningFilled" />
                  </el-icon>
                  {{ row.is_configured ? '已配置' : '待配置' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="治理目标" width="190" show-overflow-tooltip>
              <template #default="{ row }">
                <span>{{ getGovernanceTargetText(row) }}</span>
              </template>
            </el-table-column>
            <el-table-column label="治理状态" width="140">
              <template #default="{ row }">
                <el-tag :type="getGovernanceTagType(row.release_governance)" size="small">
                  {{ getGovernanceText(row.release_governance) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="创建时间" width="160">
              <template #default="{ row }">
                <span class="create-time">{{ row.created_at ? new Date(row.created_at).toLocaleString() : '-' }}</span>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="140">
              <template #default="{ row }">
                <div class="table-actions">
                  <el-tooltip content="编辑环境" placement="top">
                    <el-button v-authority="['app:application:envedit']" size="small" type="primary" circle @click="handleEditJenkinsEnv(row)">
                      <el-icon><Edit /></el-icon>
                    </el-button>
                  </el-tooltip>
                  <el-tooltip content="删除环境" placement="top">
                    <el-button v-authority="['app:application:envdelete']"
 size="small" type="danger" circle @click="handleDeleteJenkinsEnv(row)">
                      <el-icon><Delete /></el-icon>
                    </el-button>
                  </el-tooltip>
                </div>
              </template>
            </el-table-column>
          </el-table>
        </div>

        <!-- 环境表单 -->
        <div class="env-form-section" v-if="showJenkinsEnvForm || editingJenkinsEnv !== null">
          <el-card shadow="hover" class="form-card">
            <div class="section-header">
              <div class="header-left">
                <el-icon class="section-icon"><Setting /></el-icon>
                <span class="section-title">
                  {{ editingJenkinsEnv ? '编辑环境' : '新增环境' }}
                </span>
              </div>
            </div>

            <el-form :model="jenkinsEnvForm" label-width="120px" size="small" class="modern-form">
              <div class="form-row">
                <el-form-item label="环境名称" required class="form-item-half">
                  <el-input
                    v-model="jenkinsEnvForm.env_name"
                    placeholder="请输入环境名称（如：prod、test、dev）"
                  />
                </el-form-item>
                <el-form-item label="Jenkins服务器" class="form-item-half">
                  <el-select
                    v-model="jenkinsEnvForm.jenkins_server_id"
                    placeholder="请选择Jenkins服务器"
                    style="width: 100%"
                    filterable
                    clearable
                    @change="handleServerChange"
                  >
                    <el-option
                      v-for="server in jenkinsServers"
                      :key="server.id"
                      :label="server.name"
                      :value="server.id"
                    />
                  </el-select>
                </el-form-item>
              </div>
              <el-form-item label="任务名称" class="form-item-full">
                <div class="job-name-input-group">
                  <el-select
                    v-model="jenkinsEnvForm.job_name"
                    placeholder="请选择或输入Jenkins任务名称"
                    style="width: 100%"
                    filterable
                    allow-create
                    clearable
                    :loading="jobSearchLoading"
                    @input="handleJobNameInput"
                    :disabled="!jenkinsEnvForm.jenkins_server_id"
                  >
                    <el-option
                      v-for="job in jenkinsJobs"
                      :key="job"
                      :label="job"
                      :value="job"
                    />
                  </el-select>
                  <el-button
                    type="primary"
                    :icon="Search"
                    :loading="jobValidationLoading"
                    @click="validateJenkinsJob"
                    :disabled="!jenkinsEnvForm.jenkins_server_id || !jenkinsEnvForm.job_name"
                    class="validate-btn"
                  >
                    验证任务
                  </el-button>
                </div>

                <!-- 验证状态显示 -->
                <div v-if="jobValidationStatus" class="validation-status">
                  <div class="validation-result">
                    <el-icon v-if="jobValidationStatus.exists" class="success-icon">
                      <SuccessFilled />
                    </el-icon>
                    <el-icon v-else class="error-icon">
                      <WarningFilled />
                    </el-icon>
                    <span v-if="!jobValidationStatus.exists" class="error-message">
                      {{ jobValidationStatus.message }}
                    </span>
                  </div>
                </div>
              </el-form-item>
              <el-divider content-position="left">发布治理联动</el-divider>
              <div class="form-row">
                <el-form-item label="启用治理" class="form-item-half">
                  <el-switch v-model="jenkinsEnvForm.deploy_config.release_governance_enabled" />
                </el-form-item>
                <el-form-item label="预检模式" class="form-item-half">
                  <el-select v-model="jenkinsEnvForm.deploy_config.pre_check_mode" style="width: 100%" :disabled="!jenkinsEnvForm.deploy_config.release_governance_enabled">
                    <el-option label="观察模式" value="observe" />
                    <el-option label="严格模式" value="strict" />
                  </el-select>
                </el-form-item>
              </div>
              <div class="form-row">
                <el-form-item label="K8s集群" class="form-item-half">
                  <el-select
                    v-model="jenkinsEnvForm.deploy_config.cluster_id"
                    placeholder="请选择集群"
                    style="width: 100%"
                    clearable
                    :disabled="!jenkinsEnvForm.deploy_config.release_governance_enabled"
                  >
                    <el-option
                      v-for="cluster in k8sClusters"
                      :key="cluster.id"
                      :label="cluster.name"
                      :value="cluster.id"
                    />
                  </el-select>
                </el-form-item>
                <el-form-item label="命名空间" class="form-item-half">
                  <el-input
                    v-model="jenkinsEnvForm.deploy_config.namespace"
                    placeholder="例如：opsnexus-apply-e2e"
                    :disabled="!jenkinsEnvForm.deploy_config.release_governance_enabled"
                  />
                </el-form-item>
              </div>
              <div class="form-row">
                <el-form-item label="工作负载类型" class="form-item-half">
                  <el-select
                    v-model="jenkinsEnvForm.deploy_config.workload_type"
                    style="width: 100%"
                    :disabled="!jenkinsEnvForm.deploy_config.release_governance_enabled"
                  >
                    <el-option label="Deployment" value="deployment" />
                    <el-option label="StatefulSet" value="statefulset" />
                  </el-select>
                </el-form-item>
                <el-form-item label="工作负载名称" class="form-item-half">
                  <el-input
                    v-model="jenkinsEnvForm.deploy_config.workload_name"
                    placeholder="例如：demo-nginx"
                    :disabled="!jenkinsEnvForm.deploy_config.release_governance_enabled"
                  />
                </el-form-item>
              </div>
              <el-form-item class="form-actions">
                <el-button type="primary" @click="handleSaveJenkinsEnv">
                  <el-icon><component :is="editingJenkinsEnv ? Edit : Plus" /></el-icon>
                  {{ editingJenkinsEnv ? '更新环境' : '创建环境' }}
                </el-button>
                <el-button @click="resetJenkinsEnvForm">取消</el-button>
              </el-form-item>
            </el-form>
          </el-card>
        </div>
      </el-card>

      <template #footer>
        <div class="dialog-footer">
          <el-button @click="jenkinsEnvDialogVisible = false" size="small">
            <el-icon><Close /></el-icon>
            关闭
          </el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.application-management {
  padding: 0;
  min-height: auto;
  background: transparent;
}

/* 主卡片样式 */
.main-card {
  background: rgba(12, 24, 41, 0.92);
  backdrop-filter: blur(18px);
  border-radius: 24px;
  box-shadow: var(--shadow-card);
  border: 1px solid var(--border-subtle);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.title {
  font-size: 20px;
  font-weight: 700;
  color: var(--text-primary);
}

.header-actions {
  display: flex;
  gap: 12px;
}

/* 搜索卡片样式 */
.search-card {
  background: rgba(255, 255, 255, 0.03);
  backdrop-filter: blur(12px);
  border-radius: 18px;
  box-shadow: none;
  border: 1px solid var(--border-subtle);
  margin-bottom: 20px;
}

.search-card :deep(.el-card__body) {
  padding: 20px;
}

.overview-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 16px;
  margin-bottom: 20px;
}

.overview-card {
  padding: 18px 20px;
  border-radius: 18px;
  color: #fff;
  min-height: 132px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  box-shadow: var(--shadow-card);
}

.overview-card.primary {
  background: linear-gradient(135deg, #2563eb, #1d4ed8);
}

.overview-card.success {
  background: linear-gradient(135deg, #059669, #047857);
}

.overview-card.warning {
  background: linear-gradient(135deg, #d97706, #b45309);
}

.overview-card.danger {
  background: linear-gradient(135deg, #dc2626, #b91c1c);
}

.overview-label {
  font-size: 13px;
  opacity: 0.88;
}

.overview-value {
  font-size: 28px;
  font-weight: 700;
  line-height: 1.2;
}

.overview-sub {
  font-size: 12px;
  opacity: 0.82;
}

/* 表格卡片样式 */
.table-card {
  background: rgba(255, 255, 255, 0.02);
  backdrop-filter: blur(12px);
  border-radius: 18px;
  border: 1px solid var(--border-subtle);
  box-shadow: none;
}

.table-card :deep(.el-card__body) {
  padding: 0;
}

/* 搜索工具栏样式 */
.search-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.search-form {
  display: flex;
  align-items: center;
  gap: 12px;
  margin: 0;
}

.search-form :deep(.el-form-item) {
  margin-bottom: 0;
  margin-right: 16px;
}

.search-form :deep(.el-form-item__label) {
  color: var(--text-secondary);
  font-weight: 600;
}

.toolbar-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.resource-count {
  font-size: 14px;
  color: var(--text-muted);
  font-weight: 600;
  white-space: nowrap;
}

/* 表格样式 */
.app-table {
  border-radius: 8px;
  overflow: hidden;
  border: none;
}

.app-table :deep(.el-table) {
  border: none;
}

.app-table :deep(.el-table__header) {
  background: transparent;
  border: none;
}

.app-table :deep(.el-table__row:hover) {
  background-color: rgba(53, 164, 255, 0.06);
}

.app-table :deep(.el-table__border) {
  border: none;
}

.app-table :deep(.el-table th) {
  border: none;
}

.app-table :deep(.el-table td) {
  border: none;
}

.app-table :deep(.el-table__body-wrapper) {
  border: none;
}

.app-name {
  display: flex;
  align-items: center;
  gap: 8px;
}

.app-icon {
  color: #409EFF;
  font-size: 16px;
}

.app-name-link {
  color: #9bd1ff;
  cursor: pointer;
  font-weight: 500;
  transition: color 0.3s;
}

.app-name-link:hover {
  color: #cde7ff;
  text-decoration: underline;
}

.env-coverage-cell {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.coverage-value {
  font-size: 15px;
  font-weight: 700;
  color: var(--text-primary);
}

.coverage-sub {
  font-size: 11px;
  color: #6b7280;
}

.latest-deployment-cell {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.latest-title {
  color: #2563eb;
  cursor: pointer;
  font-weight: 600;
}

.latest-title:hover {
  text-decoration: underline;
}

.latest-meta {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
  font-size: 12px;
  color: #374151;
}

.muted {
  color: #6b7280;
}

.operation-buttons {
  display: flex;
  gap: 4px;
  justify-content: flex-start;
  align-items: center;
  flex-wrap: nowrap;
}

.operation-buttons .el-button {
  padding: 0;
  min-width: 28px;
  width: 28px;
  height: 28px;
}

/* 分页样式 */
.pagination-wrapper {
  margin-top: 20px;
  display: flex;
  justify-content: center;
  padding: 20px;
  background: rgba(248, 249, 250, 0.8);
  border-radius: 0 0 12px 12px;
}

/* 对话框样式 */
.compact-dialog :deep(.el-dialog) {
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
}

.compact-dialog :deep(.el-dialog__header) {
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: white;
  border-top-left-radius: 16px;
  border-top-right-radius: 16px;
  padding: 20px 24px;
}

.compact-dialog :deep(.el-dialog__title) {
  color: white;
  font-weight: 600;
}

.compact-dialog :deep(.el-dialog__body) {
  padding: 15px 20px;
}

.compact-form {
  margin: 0;
}

.compact-form :deep(.el-form-item) {
  margin-bottom: 12px;
}

.form-row {
  display: flex;
  gap: 16px;
  margin-bottom: 12px;
}

.form-item-half {
  flex: 1;
  margin-bottom: 0 !important;
}

.form-item-full {
  width: 100%;
  margin-bottom: 0 !important;
}

.dialog-footer {
  display: flex;
  gap: 12px;
  justify-content: flex-end;
}

/* 详情对话框样式 */
.detail-panel {
  display: flex;
  flex-direction: column;
  gap: 18px;
}

.detail-hero {
  display: flex;
  justify-content: space-between;
  gap: 20px;
  padding: 20px 22px;
  border-radius: 16px;
  background: linear-gradient(135deg, rgba(37, 99, 235, 0.12), rgba(29, 78, 216, 0.05));
  border: 1px solid rgba(37, 99, 235, 0.14);
}

.detail-title-row {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 10px;
}

.detail-app-name {
  font-size: 24px;
  font-weight: 700;
  color: #111827;
}

.detail-app-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
  color: #4b5563;
  font-size: 13px;
}

.detail-actions-row {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  flex-wrap: wrap;
}

.detail-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 18px;
}

.detail-card {
  padding: 18px 20px;
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.92);
  border: 1px solid rgba(148, 163, 184, 0.2);
  box-shadow: 0 10px 30px rgba(15, 23, 42, 0.06);
}

.detail-card.full-width {
  width: 100%;
}

.detail-card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
  margin-bottom: 14px;
}

.detail-card-title {
  font-size: 16px;
  font-weight: 700;
  color: #111827;
  margin-bottom: 14px;
}

.detail-card-header .detail-card-title {
  margin-bottom: 0;
}

.detail-card-subtitle {
  font-size: 12px;
  color: #6b7280;
}

.detail-item {
  display: flex;
  gap: 12px;
  padding: 10px 0;
  border-bottom: 1px solid #eef2f7;
}

.detail-item:last-child {
  border-bottom: none;
}

.detail-label {
  font-weight: 600;
  color: #6b7280;
  min-width: 88px;
  flex-shrink: 0;
}

.detail-value {
  color: #111827;
  word-break: break-all;
  flex: 1;
}

.detail-link {
  color: #2563eb;
  cursor: pointer;
}

.detail-link:hover {
  text-decoration: underline;
}

.detail-table-actions {
  display: flex;
  align-items: center;
  gap: 6px;
  flex-wrap: wrap;
}

.governance-focus-alert {
  border-radius: 16px;
}

.governance-focus-title {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

.governance-focus-actions {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  flex-wrap: wrap;
}

.governance-focus-target {
  color: #475569;
  font-size: 13px;
}

:deep(.focused-env-row) {
  --el-table-tr-bg-color: rgba(59, 130, 246, 0.08);
}

/* 通用样式 */
.el-tag {
  font-weight: 500;
  border-radius: 8px;
  border: none;
}

.el-button {
  border-radius: 8px;
  font-weight: 500;
  transition: all 0.3s ease;
}

.el-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.el-input :deep(.el-input__wrapper),
.el-select :deep(.el-input__wrapper) {
  background: rgba(255, 255, 255, 0.8);
  border: 1px solid rgba(103, 126, 234, 0.2);
  border-radius: 8px;
  box-shadow: none;
  transition: all 0.3s ease;
}

.el-input :deep(.el-input__wrapper):hover,
.el-select :deep(.el-input__wrapper):hover {
  border-color: #c0c4cc;
}

.el-input :deep(.el-input__wrapper.is-focus),
.el-select :deep(.el-input__wrapper.is-focus) {
  border-color: #667eea;
  box-shadow: 0 0 0 2px rgba(103, 126, 234, 0.2);
  background: rgba(255, 255, 255, 1);
}

.el-input :deep(.el-input__inner),
.el-select :deep(.el-input__inner) {
  background: transparent;
  border: none;
  color: #2c3e50;
}

.el-textarea :deep(.el-textarea__inner) {
  border-radius: 8px;
  border: 1px solid rgba(103, 126, 234, 0.2);
  transition: all 0.3s ease;
}

.el-textarea :deep(.el-textarea__inner):focus {
  border-color: #667eea;
  box-shadow: 0 0 0 2px rgba(103, 126, 234, 0.2);
}

.el-loading-mask {
  background-color: rgba(103, 126, 234, 0.1);
  backdrop-filter: blur(4px);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .application-management {
    padding: 12px;
  }

  .overview-grid {
    grid-template-columns: 1fr;
  }

  .page-header {
    flex-direction: column;
    align-items: stretch;
    gap: 12px;
  }

  .header-right {
    justify-content: flex-end;
  }

  .search-toolbar {
    flex-direction: column;
    align-items: stretch;
    gap: 12px;
  }

  .toolbar-right {
    justify-content: center;
  }

  .search-form {
    flex-direction: column;
    align-items: stretch;
    gap: 8px;
  }

  .compact-dialog {
    width: 95% !important;
  }

  .form-row {
    flex-direction: column;
    gap: 8px;
  }

  .detail-hero {
    flex-direction: column;
  }

  .detail-grid {
    grid-template-columns: 1fr;
  }

  .operation-buttons {
    flex-direction: column;
    gap: 4px;
  }

  .app-table :deep(.el-table) {
    font-size: 12px;
  }

  .pagination-wrapper {
    padding: 16px 12px;
  }
}

@media (max-width: 480px) {
  .app-table :deep(.el-table__cell) {
    padding: 8px 4px;
  }

  .operation-buttons .el-button {
    padding: 4px;
    width: 28px;
    height: 28px;
  }
}

/* Jenkins环境配置对话框样式 */
.jenkins-env-dialog.modern-dialog :deep(.el-dialog__body) {
  padding: 0;
}

.jenkins-env-dialog .env-main-card {
  border: none;
  box-shadow: none;
  background: linear-gradient(135deg, #f8f9fa 0%, #ffffff 100%);
}

.jenkins-env-dialog .env-list-section {
  margin-bottom: 24px;
}

.jenkins-env-dialog .env-form-section {
  margin-top: 24px;
}

.jenkins-env-dialog .form-card {
  background: linear-gradient(135deg, #ffffff 0%, #f8f9fa 100%);
  border: 1px solid #e3e8f0;
}

.jenkins-env-dialog .section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding-bottom: 12px;
  border-bottom: 2px solid #e3e8f0;
}

.jenkins-env-dialog .header-left {
  display: flex;
  align-items: center;
  gap: 8px;
}

.jenkins-env-dialog .section-icon {
  color: #409EFF;
  font-size: 18px;
}

.jenkins-env-dialog .section-title {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

.jenkins-env-dialog .add-btn {
  border-radius: 20px;
  padding: 8px 16px;
  font-weight: 500;
}

.jenkins-env-dialog .env-table {
  border: 1px solid #e3e8f0;
  border-radius: 12px;
  overflow: hidden;
}

.jenkins-env-dialog .env-table :deep(.el-table__header) {
  background: linear-gradient(135deg, #409EFF 0%, #66b1ff 100%);
}

.jenkins-env-dialog .env-table :deep(.el-table__header th) {
  background: transparent;
  color: #2c3e50;
  font-weight: 600;
  border: none;
}

.jenkins-env-dialog .env-table :deep(.el-table__body tr:hover) {
  background-color: rgba(64, 158, 255, 0.1);
}

.jenkins-env-dialog .env-tag {
  background: linear-gradient(135deg, #409EFF 0%, #66b1ff 100%);
  color: white;
  border: none;
  font-weight: 500;
}

.jenkins-env-dialog .status-tag {
  display: flex;
  align-items: center;
  gap: 4px;
  font-weight: 500;
}

.jenkins-env-dialog .server-id,
.jenkins-env-dialog .job-name,
.jenkins-env-dialog .create-time {
  color: #606266;
  font-size: 13px;
}

.jenkins-env-dialog .table-actions {
  display: flex;
  gap: 8px;
  justify-content: center;
}

.jenkins-env-dialog .modern-form .form-row {
  display: flex;
  gap: 20px;
}

.jenkins-env-dialog .modern-form .form-item-half {
  flex: 1;
}

.jenkins-env-dialog .modern-form .form-item-full {
  width: 100%;
}

.jenkins-env-dialog .modern-form .form-actions {
  margin-top: 24px;
  padding-top: 20px;
  border-top: 1px solid #e3e8f0;
}

.jenkins-env-dialog .modern-form .el-form-item {
  margin-bottom: 20px;
}

.jenkins-env-dialog .dialog-footer {
  text-align: center;
  padding-top: 16px;
  border-top: 1px solid #e3e8f0;
}

/* 任务名称输入组样式 */
.jenkins-env-dialog .job-name-input-group {
  display: flex;
  gap: 8px;
  align-items: flex-start;
}

.jenkins-env-dialog .job-name-input-group .el-select {
  flex: 1;
}

.jenkins-env-dialog .validate-btn {
  flex-shrink: 0;
  padding: 8px 16px;
  border-radius: 6px;
  font-weight: 500;
}

.jenkins-env-dialog .validation-status {
  margin-top: 12px;
}

.jenkins-env-dialog .validation-result {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  border-radius: 6px;
  background-color: rgba(0, 0, 0, 0.02);
}

.jenkins-env-dialog .success-icon {
  color: #67c23a;
  font-size: 18px;
}

.jenkins-env-dialog .error-icon {
  color: #e6a23c;
  font-size: 18px;
}

.jenkins-env-dialog .error-message {
  color: #e6a23c;
  font-size: 14px;
  font-weight: 500;
}

.jenkins-env-dialog .server-name {
  color: #409EFF;
  font-weight: 500;
}
</style>

