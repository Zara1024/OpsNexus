<script setup>
import { ref, reactive, onMounted, computed, watch, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import yaml from 'js-yaml'
import {
  Search,
  Refresh,
  Plus,
  Delete,
  Setting,
  Monitor,
  DataAnalysis,
  DocumentCopy,
  Document,
  Connection
} from '@element-plus/icons-vue'
import k8sApi from '@/api/k8s'

// Import modular components
import ClusterSelector from './pods/ClusterSelector.vue'
import NamespaceSelector from './pods/NamespaceSelector.vue'
import PodListDialog from './pods/PodListDialog.vue'
import PodEventsDialog from './pods/PodEventsDialog.vue'
import PodYamlDialog from './pods/PodYamlDialog.vue'
import PodConfigDialog from './pods/PodConfigDialog.vue'
import CreatePodDialog from './pods/CreatePodDialog.vue'

const router = useRouter()
const route = useRoute()

const loading = ref(false)
const activeTab = ref('')
const queryParams = reactive({
  name: '',
  type: '',
  namespace: 'default'
})

const tableData = ref([])
const selectedClusterId = ref('')
const clusterList = ref([])
const autoSwitchingNamespace = ref(false)
// clusterList 已移至 ClusterSelector 组件
// namespaceList 和 namespaceLoading 已移至 NamespaceSelector 组件

// 对话框状态
const podListDialogVisible = ref(false)
const podEventsDialogVisible = ref(false)
const podYamlDialogVisible = ref(false)
const podYamlDialogEditable = ref(false)
const workloadYamlDialogVisible = ref(false)
const logDialogVisible = ref(false)
const scaleDialogVisible = ref(false)
const autoscalingDialogVisible = ref(false)
const capacitySuggestionDialogVisible = ref(false)
const governanceDialogVisible = ref(false)
const workloadLabelsDialogVisible = ref(false)
const allImagesDialogVisible = ref(false)
const podConfigDialogVisible = ref(false)
const schedulingDialogVisible = ref(false)
const createPodDialogVisible = ref(false)

// 当前操作的工作负载或Pod
const currentWorkload = ref({})
const currentPod = ref({})
const currentPodForEvents = ref({})
const currentPodLogs = ref('')
const currentYaml = ref('')

// YAML编辑器引用
const yamlEditor = ref(null)
const createPodDialogRef = ref(null)

// 扩容缩容表单
const scaleForm = reactive({
  replicas: 1
})

const autoscalingTargetTypeOptions = [
  { label: '利用率', value: 'Utilization' },
  { label: '平均值', value: 'AverageValue' }
]

const autoscalingSelectPolicyOptions = [
  { label: '取最大', value: 'Max' },
  { label: '取最小', value: 'Min' },
  { label: '禁用', value: 'Disabled' }
]

const createDefaultAutoscalingForm = () => ({
  minReplicas: 1,
  maxReplicas: 3,
  cpuEnabled: true,
  cpuTargetType: 'AverageValue',
  cpuTargetValue: '100m',
  memoryEnabled: false,
  memoryTargetType: 'AverageValue',
  memoryTargetValue: '256Mi',
  scaleUpWindow: 0,
  scaleDownWindow: 300,
  scaleUpSelectPolicy: 'Max',
  scaleDownSelectPolicy: 'Max',
  scaleUpPodsEnabled: true,
  scaleUpPodsValue: 4,
  scaleUpPercentEnabled: true,
  scaleUpPercentValue: 100,
  scaleDownPodsEnabled: false,
  scaleDownPodsValue: 1,
  scaleDownPercentEnabled: true,
  scaleDownPercentValue: 100
})

const autoscalingForm = reactive(createDefaultAutoscalingForm())
const autoscalingLoading = ref(false)
const autoscalingSubmitting = ref(false)
const autoscalingExists = ref(false)
const capacitySuggestionLoading = ref(false)
const capacitySuggestionHistoryLoading = ref(false)
const governanceLoading = ref(false)
const capacitySuggestionHistory = ref([])
const capacitySuggestionHistoryPagination = reactive({
  pageNum: 1,
  pageSize: 5,
  total: 0
})

const createEmptyAlertSummary = () => ({
  openEventCount: 0,
  resolvedEventCount: 0,
  incidentCount: 0,
  criticalCount: 0,
  recentEvents: [],
  recentIncidents: []
})

const createEmptyCapacitySuggestion = () => ({
  historyId: 0,
  generatedAt: '',
  generatedById: 0,
  generatedBy: '',
  riskLevel: '',
  report: '',
  renderedPrompt: '',
  systemPrompt: '',
  alertKeyword: '',
  alertCenterPath: '',
  alertCenterQuery: {
    keyword: '',
    namespace: '',
    workloadName: '',
    source: ''
  },
  recommendedActions: [],
  recommendedPolicy: {},
  watchMetrics: [],
  alertSummary: createEmptyAlertSummary()
})

const cloneCapacitySuggestion = (payload = createEmptyCapacitySuggestion()) => JSON.parse(JSON.stringify(payload))

const createEmptyGovernanceOverview = () => ({
  clusterId: 0,
  clusterName: '',
  namespace: '',
  workloadType: '',
  workloadName: '',
  riskLevel: '',
  blocking: false,
  blockingReason: '',
  warnings: [],
  blockingRules: [],
  warningRules: [],
  stages: [],
  currentSuggestion: createEmptyCapacitySuggestion(),
  latestSnapshot: null,
  currentAutoscaling: null,
  alertSummary: createEmptyAlertSummary(),
  watchMetrics: [],
  extensions: [],
  relatedApplications: [],
  alertCenterPath: '',
  aiDiagnosisPath: '',
  k8sWorkloadPath: ''
})

const capacitySuggestion = ref(createEmptyCapacitySuggestion())
const latestCapacitySuggestion = ref(createEmptyCapacitySuggestion())
const governanceOverview = ref(createEmptyGovernanceOverview())
const capacitySuggestionViewingLatest = computed(() => {
  const latestId = Number(latestCapacitySuggestion.value?.historyId || 0)
  const currentId = Number(capacitySuggestion.value?.historyId || 0)
  if (!currentId) return true
  if (!latestId) return false
  return latestId === currentId
})

const schedulingForm = reactive({
  nodeSelectorText: '',
  nodeAffinity: 'none',
  podAntiAffinity: false
})

// 日志查看参数
const logParams = reactive({
  container: '',
  lines: 100,
  follow: false
})

// 获取集群列表
const fetchClusterList = async () => {
  try {
    const response = await k8sApi.getClusterList()
    const responseData = response.data || response
    
    if (responseData.code === 200 || responseData.success) {
      const clusters = responseData.data?.list || responseData.data || []
      clusterList.value = clusters.map(cluster => ({
        id: cluster.id,
        name: cluster.name,
        status: cluster.status
      }))
      
      if (clusterList.value.length > 0 && !selectedClusterId.value) {
        const onlineCluster = clusterList.value.find(cluster => cluster.status === 2)
        selectedClusterId.value = onlineCluster ? onlineCluster.id : clusterList.value[0].id
      }
      
      console.log('集群列表加载成功:', clusterList.value)
    } else {
      ElMessage.error(responseData.message || '获取集群列表失败')
    }
  } catch (error) {
    console.error('获取集群列表失败:', error)
    ElMessage.warning('无法获取集群列表，请检查后端服务')
  }
}

// namespaceRequestPromise 已移至 NamespaceSelector 组件

// fetchNamespaceList 已移至 NamespaceSelector 组件

// 处理标签页切换
const handleTabChange = (tabName) => {
  console.log('标签页切换到:', tabName)
  activeTab.value = tabName
  queryParams.type = tabName
  handleQuery()
}

const findFirstNamespaceWithWorkloads = async (clusterId, params = {}, currentNamespace = '') => {
  const response = await k8sApi.getNamespaces(clusterId)
  const responseData = response.data || response
  const namespaces = responseData.data?.namespaces || responseData.data || []
  const namespaceList = (Array.isArray(namespaces) ? namespaces : [])
    .map((item) => item.name || item)
    .filter(Boolean)
    .filter((name) => name !== currentNamespace)

  for (const namespaceName of namespaceList) {
    const workloadResponse = await k8sApi.getWorkloadList(clusterId, namespaceName, params)
    const workloadData = workloadResponse.data || workloadResponse
    const workloads = workloadData.data?.workloads || workloadData.data || []
    if (Array.isArray(workloads) && workloads.length > 0) {
      return namespaceName
    }
  }

  return ''
}

// 监听activeTab变化，同步到queryParams.type
watch(activeTab, (newType) => {
  queryParams.type = newType
})

// 查询工作负载列表
const handleQuery = async () => {
  const queryStartTime = Date.now()
  
  try {
    if (!selectedClusterId.value) {
      ElMessage.warning('请选择一个集群')
      return
    }
    
    if (!queryParams.namespace) {
      ElMessage.warning('请选择命名空间')
      return
    }
    
    console.log('🔍 开始查询工作负载:', {
      clusterId: selectedClusterId.value,
      namespace: queryParams.namespace,
      type: queryParams.type,
      name: queryParams.name
    })
    
    loading.value = true
    
    const params = {}
    if (queryParams.type) params.type = queryParams.type
    if (queryParams.name) params.name = queryParams.name
    
    const response = await k8sApi.getWorkloadList(selectedClusterId.value, queryParams.namespace, params)
    
    const responseData = response.data || response
    console.log('工作负载列表API响应:', responseData)
    
    if (responseData.code === 200 || responseData.success) {
      // 根据API响应，数据结构是 { data: { workloads: [...] } }
      const workloads = responseData.data?.workloads || responseData.data || []
      // 确保workloads是数组
      const workloadList = Array.isArray(workloads) ? workloads : []
      if (
        workloadList.length === 0 &&
        !queryParams.name &&
        !queryParams.type &&
        !autoSwitchingNamespace.value
      ) {
        const fallbackNamespace = await findFirstNamespaceWithWorkloads(
          selectedClusterId.value,
          params,
          queryParams.namespace
        )
        if (fallbackNamespace) {
          autoSwitchingNamespace.value = true
          queryParams.namespace = fallbackNamespace
          ElMessage.info(`已切换到存在工作负载的命名空间：${fallbackNamespace}`)
          await handleQuery()
          return
        }
      }
      tableData.value = workloadList.map(workload => ({
        id: workload.name,
        name: workload.name,
        type: workload.type?.toLowerCase() || workload.kind?.toLowerCase(),
        namespace: workload.namespace,
        replicas: `${workload.readyReplicas || 0}/${workload.replicas || 0}`,
        readyReplicas: workload.readyReplicas || 0,
        totalReplicas: workload.replicas || 0,
        images: workload.images || [],
        labels: workload.labels || {},
        status: workload.status || getWorkloadStatus(workload),
        age: formatAge(workload.createdAt),
        updateTime: workload.createdAt,
        updatedAt: workload.updatedAt,
        conditions: workload.conditions || [],
        resources: workload.resources || {
          cpu: { requests: '0', limits: '0' },
          memory: { requests: '0', limits: '0' }
        },
        autoscaling: workload.autoscaling || null,
        rawData: workload
      }))
      
      console.log('工作负载列表加载成功:', tableData.value)
    } else {
      const errorMsg = responseData.message || '获取工作负载列表失败'

      // 特殊处理资源不存在的错误
      if (errorMsg.includes('the server could not find the requested resource')) {
        if (queryParams.type === 'cronjobs') {
          ElMessage.warning('当前集群不支持CronJob资源，可能是Kubernetes版本过低')
        } else if (queryParams.type) {
          ElMessage.warning(`当前集群不支持${queryParams.type}资源类型`)
        } else {
          ElMessage.warning('请求的资源不存在，请检查集群配置')
        }
      } else {
        ElMessage.error(errorMsg)
      }

      tableData.value = []
    }
  } catch (error) {
    console.error('获取工作负载列表失败:', error)
    
    if (error.code === 'ERR_NETWORK' || 
        error.message?.includes('ERR_CONNECTION_REFUSED') ||
        error.message?.includes('Failed to fetch')) {
      ElMessage.warning('后端服务连接失败，请检查服务状态')
    } else if (error.response?.status === 401) {
      ElMessage.error('认证失败，请重新登录')
    } else if (error.response?.status === 403) {
      ElMessage.error('权限不足，请联系管理员')
    } else {
      console.warn('API调用异常，但可能数据已正确加载')
    }
    
    tableData.value = []
  } finally {
    autoSwitchingNamespace.value = false
    loading.value = false
    console.log('✅ 工作负载查询完成，耗时:', Date.now() - queryStartTime + 'ms')
  }
}

// 获取工作负载状态
const getWorkloadStatus = (workload) => {
  // 如果后端直接返回了状态，优先使用
  if (workload.status) return workload.status
  
  const replicas = workload.replicas || 0
  const readyReplicas = workload.readyReplicas || 0
  
  if (workload.type === 'job' || workload.kind === 'Job') {
    return workload.succeeded ? 'Completed' : 
           workload.failed ? 'Failed' : 'Running'
  }
  
  if (workload.type === 'cronjob' || workload.kind === 'CronJob') {
    return workload.lastScheduleTime ? 'Active' : 'Suspended'
  }
  
  if (replicas === 0) return 'Stopped'
  if (readyReplicas === 0) return 'Pending'
  if (readyReplicas < replicas) return 'Partial'
  return 'Running'
}

const formatAge = (createdTimestamp) => {
  if (!createdTimestamp) return 'Unknown'
  
  const now = new Date()
  const created = new Date(createdTimestamp)
  const diff = Math.floor((now - created) / 1000)
  
  if (diff < 60) return `${diff}s`
  if (diff < 3600) return `${Math.floor(diff / 60)}m`
  if (diff < 86400) return `${Math.floor(diff / 3600)}h`
  return `${Math.floor(diff / 86400)}d`
}

const formatDateTime = (timestamp) => {
  if (!timestamp) return '-'
  
  const date = new Date(timestamp)
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  
  return `${year}-${month}-${day} ${hours}:${minutes}`
}

const resetQuery = () => {
  queryParams.name = ''
  queryParams.type = ''
  handleQuery()
}

const applyWorkloadTableLayout = () => {
  nextTick(() => {
    const hiddenColumns = new Set([5, 7])
    const widthMap = {
      0: 200,
      1: 84,
      2: 108,
      3: 136,
      4: 160,
      6: 180,
      8: 140
    }

    const syncMainTable = (tableElement) => {
      tableElement.style.width = '980px'
      tableElement.style.maxWidth = '980px'

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

    document.querySelectorAll('.workloads-table .el-table__header-wrapper table, .workloads-table .el-table__body-wrapper table').forEach(syncMainTable)

    document.querySelectorAll('.workloads-table .el-table__fixed-right').forEach(wrapper => {
      wrapper.style.width = '220px'
      wrapper.style.right = '0px'
    })
    document.querySelectorAll('.workloads-table .el-table__fixed-right colgroup col, .workloads-table .el-table__fixed-right .el-table__cell').forEach(element => {
      element.style.width = '220px'
      element.style.maxWidth = '220px'
    })
  })
}

// 导航到监控仪表板
const navigateToMonitoring = () => {
  router.push('/k8s/monitoring')
}

// 创建工作负载相关函数
const showCreatePodDialog = () => {
  if (!selectedClusterId.value) {
    ElMessage.warning('请先选择集群')
    return
  }
  if (!queryParams.namespace) {
    ElMessage.warning('请先选择命名空间')
    return
  }

  createPodDialogVisible.value = true
}

// 处理YAML创建校验
const handlePodPreview = async (data) => {
  try {
    createPodDialogRef.value?.setLoading(true)

    // 使用validateYaml来校验YAML格式
    const response = await k8sApi.validateYaml(selectedClusterId.value, data.yamlContent)
    const responseData = response.data || response

    const result = {
      success: responseData.code === 200,
      message: responseData.message || (responseData.code === 200 ? '可以创建工作负载' : '创建预览失败'),
      details: responseData.data
    }

    createPodDialogRef.value?.setDryRunResult(result)

    if (result.success) {
      ElMessage.success(result.message)
    } else {
      ElMessage.error(result.message)
    }
  } catch (error) {
    const result = {
      success: false,
      message: error.message || '工作负载创建预览失败',
      details: null
    }
    createPodDialogRef.value?.setDryRunResult(result)
    ElMessage.error('工作负载创建预览失败: ' + (error.message || '网络错误'))
  } finally {
    createPodDialogRef.value?.setLoading(false)
  }
}

// 处理工作负载创建
const handlePodCreate = async (data) => {
  try {
    createPodDialogRef.value?.setLoading(true)

    // 使用createPodFromYaml来创建工作负载（该API支持多种资源类型）
    const response = await k8sApi.createPodFromYaml(selectedClusterId.value, queryParams.namespace, data)
    const responseData = response.data || response

    if (responseData.code === 200) {
      ElMessage.success('工作负载创建成功!')
      createPodDialogVisible.value = false
      handleQuery() // 刷新工作负载列表
    } else {
      ElMessage.error(responseData.message || '工作负载创建失败')
    }
  } catch (error) {
    console.error('工作负载创建失败:', error)
    ElMessage.error('工作负载创建失败: ' + (error.message || '网络错误'))
  } finally {
    createPodDialogRef.value?.setLoading(false)
  }
}


// 集群选择变化处理
const handleClusterChange = async () => {
  // 清空数据，NamespaceSelector 组件会自动处理命名空间列表加载
  tableData.value = []

  if (selectedClusterId.value && queryParams.namespace) {
    handleQuery()
  }
}

// 命名空间选择变化处理
const handleNamespaceChange = () => {
  if (selectedClusterId.value && queryParams.namespace) {
    handleQuery()
  } else {
    tableData.value = []
  }
}

// 导航到容器详情页面
const navigateToPodDetail = async (row) => {
  try {
    console.log('🔍 点击工作负载名称，跳转到Pod详情:', row)
    console.log('📊 工作负载详情:', {
      name: row.name,
      type: row.type,
      namespace: queryParams.namespace
    })

    // 使用新的专门API获取该工作负载下的Pod列表
    const response = await k8sApi.getWorkloadPods(
      selectedClusterId.value,
      queryParams.namespace,
      row.type.toLowerCase(),
      row.name
    )
    const responseData = response.data || response

    if (responseData.code === 200 && responseData.data && responseData.data.length > 0) {
      // 获取第一个Pod（如果有多个Pod，跳转到第一个）
      const firstPod = responseData.data[0]
      console.log('🎯 跳转到第一个Pod:', firstPod.name)

      router.push({
        path: `/k8s/pod/${selectedClusterId.value}/${queryParams.namespace}/${firstPod.name}`
      })
    } else {
      ElMessage.warning('该工作负载下暂无Pod或Pod信息获取失败')
    }
  } catch (error) {
    console.error('获取工作负载Pod信息失败:', error)
    ElMessage.error('获取工作负载Pod信息失败，请检查网络连接')
  }
}


// 查看Pod列表
const viewPodList = async (row) => {
  try {
    loading.value = true
    console.log('🔍 点击容器组数量，查看Pod列表:', row)
    console.log('📊 工作负载详情:', {
      name: row.name,
      type: row.type,
      namespace: queryParams.namespace
    })

    // 使用新的专门API获取该工作负载下的Pod列表
    const response = await k8sApi.getWorkloadPods(
      selectedClusterId.value,
      queryParams.namespace,
      row.type.toLowerCase(),
      row.name
    )

    const responseData = response.data || response
    if (responseData.code === 200 || responseData.success) {
      // 根据新API响应，数据直接在 data 数组中
      const pods = responseData.data || []
      console.log('📋 获取到的Pod列表:', pods.length, '个Pod')
      console.log('📋 Pod详细信息:', pods.map(p => ({ name: p.name, status: p.status })))

      currentWorkload.value = {
        ...row,
        pods: pods.map(pod => ({
          name: pod.name,
          status: pod.status || pod.phase || 'Unknown',
          restartCount: pod.restarts || pod.restartCount || 0,
          nodeName: pod.nodeName || 'Unknown',
          podIP: pod.podIP || 'Unknown',
          hostIP: pod.hostIP || 'Unknown',
          age: pod.age || formatAge(pod.createdAt),
          runningTime: pod.runningTime || '',
          containers: pod.containers || [],
          resources: pod.resources || {
            requests: { cpu: '', memory: '' },
            limits: { cpu: '', memory: '' }
          },
          labels: pod.labels || {},
          conditions: pod.conditions || [],
          rawData: pod
        }))
      }
      podListDialogVisible.value = true
    } else {
      ElMessage.error(responseData.message || '获取Pod列表失败')
    }
  } catch (error) {
    console.error('获取Pod列表失败:', error)
    ElMessage.error('获取Pod列表失败，请检查网络连接')
  } finally {
    loading.value = false
  }
}

// 查看Pod日志
const viewPodLogs = async (pod) => {
  try {
    currentPod.value = pod
    logParams.container = pod.containers?.[0]?.name || ''
    
    loading.value = true
    const response = await k8sApi.getPodLogs(selectedClusterId.value, queryParams.namespace, pod.name, {
      container: logParams.container,
      tailLines: logParams.lines
    })
    
    const responseData = response.data || response
    if (responseData.code === 200 || responseData.success) {
      currentPodLogs.value = responseData.data || '暂无日志'
      logDialogVisible.value = true
    } else {
      ElMessage.error(responseData.message || '获取Pod日志失败')
    }
  } catch (error) {
    console.error('获取Pod日志失败:', error)
    ElMessage.error('获取Pod日志失败，请检查网络连接')
  } finally {
    loading.value = false
  }
}

// 查看Pod事件
const viewPodEvents = (pod) => {
  currentPodForEvents.value = pod
  podEventsDialogVisible.value = true
}

// 查看YAML
const viewYaml = async (pod, editable = false) => {
  try {
    loading.value = true
    const response = await k8sApi.getPodYaml(selectedClusterId.value, queryParams.namespace, pod.name)
    
    const responseData = response.data || response
    if (responseData.code === 200 || responseData.success) {
      currentYaml.value = responseData.data || 'apiVersion: v1\nkind: Pod\nmetadata:\n  name: ' + pod.name
      currentPod.value = pod
      podYamlDialogEditable.value = editable
      podYamlDialogVisible.value = true
    } else {
      ElMessage.error(responseData.message || '获取Pod YAML失败')
    }
  } catch (error) {
    console.error('获取Pod YAML失败:', error)
    ElMessage.error('获取Pod YAML失败，请检查网络连接')
  } finally {
    loading.value = false
  }
}

// 重构Pod
const rebuildPod = async (pod) => {
  try {
    await ElMessageBox.confirm(
      `确定要重构Pod "${pod.name}" 吗？\n重构操作会删除当前Pod并自动创建新的Pod实例。`,
      '重构Pod确认',
      {
        confirmButtonText: '确定重构',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )

    const response = await k8sApi.deletePod(selectedClusterId.value, queryParams.namespace, pod.name)

    const responseData = response.data || response
    if (responseData.code === 200 || responseData.success) {
      ElMessage.success(`Pod ${pod.name} 重构成功，新的Pod实例将自动创建`)
      // 重新获取Pod列表
      if (podListDialogVisible.value) {
        const currentRow = currentWorkload.value
        await viewPodList(currentRow)
      }
    } else {
      ElMessage.error(responseData.message || '重构Pod失败')
    }
  } catch (error) {
    if (error === 'cancel') {
      ElMessage.info('已取消重构操作')
    } else {
      console.error('重构Pod失败:', error)
      ElMessage.error('重构Pod失败，请检查网络连接')
    }
  }
}

// 扩容缩容
const scaleWorkload = (row) => {
  if (!canScale(row)) {
    ElMessage.warning('该工作负载不支持扩缩容操作')
    return
  }
  
  currentWorkload.value = row
  scaleForm.replicas = row.totalReplicas || 1
  scaleDialogVisible.value = true
}

// 提交扩容缩容
const submitScale = async () => {
  try {
    const response = await k8sApi.scaleWorkload(
      selectedClusterId.value, 
      queryParams.namespace, 
      currentWorkload.value.type,
      currentWorkload.value.name, 
      { replicas: scaleForm.replicas }
    )
    
    const responseData = response.data || response
    if (responseData.code === 200 || responseData.success) {
      ElMessage.success(`${currentWorkload.value.name} 扩缩容成功`)
      scaleDialogVisible.value = false
      handleQuery()
    } else {
      ElMessage.error(responseData.message || '扩缩容失败')
    }
  } catch (error) {
    console.error('扩缩容失败:', error)
    ElMessage.error('扩缩容失败，请检查网络连接')
  }
}

const getResourceRequestValue = (workload, resourceName) => {
  const value = workload?.resources?.requests?.[resourceName]
  return value && value !== '0' ? value : ''
}

const hasResourceRequest = (workload, resourceName) => Boolean(getResourceRequestValue(workload, resourceName))

const getDefaultAutoscalingForm = (workload = {}) => {
  const totalReplicas = Number(workload.totalReplicas || workload.replicas || 1)
  const minReplicas = totalReplicas > 0 ? totalReplicas : 1
  const cpuUseUtilization = hasResourceRequest(workload, 'cpu')
  const memoryUseUtilization = hasResourceRequest(workload, 'memory')

  return {
    minReplicas,
    maxReplicas: Math.max(minReplicas + 1, 3),
    cpuEnabled: true,
    cpuTargetType: cpuUseUtilization ? 'Utilization' : 'AverageValue',
    cpuTargetValue: cpuUseUtilization ? '70' : '100m',
    memoryEnabled: false,
    memoryTargetType: memoryUseUtilization ? 'Utilization' : 'AverageValue',
    memoryTargetValue: memoryUseUtilization ? '80' : '256Mi',
    scaleUpWindow: 0,
    scaleDownWindow: 300,
    scaleUpSelectPolicy: 'Max',
    scaleDownSelectPolicy: 'Max',
    scaleUpPodsEnabled: true,
    scaleUpPodsValue: 4,
    scaleUpPercentEnabled: true,
    scaleUpPercentValue: 100,
    scaleDownPodsEnabled: false,
    scaleDownPodsValue: 1,
    scaleDownPercentEnabled: true,
    scaleDownPercentValue: 100
  }
}

const resetAutoscalingForm = (workload = currentWorkload.value) => {
  Object.assign(autoscalingForm, getDefaultAutoscalingForm(workload))
  autoscalingExists.value = false
}

const getBehaviorPolicyValue = (rule, type) => {
  const policy = (rule?.policies || []).find(item => item.type === type)
  return policy?.value ?? null
}

const fillAutoscalingForm = (autoscaling) => {
  const nextForm = getDefaultAutoscalingForm(currentWorkload.value)
  nextForm.cpuEnabled = false
  nextForm.memoryEnabled = false

  if (autoscaling?.minReplicas) nextForm.minReplicas = autoscaling.minReplicas
  if (autoscaling?.maxReplicas) nextForm.maxReplicas = autoscaling.maxReplicas

  ;(autoscaling?.metrics || []).forEach(metric => {
    const resourceName = String(metric?.resourceName || '').toLowerCase()
    if (!['cpu', 'memory'].includes(resourceName)) return

    nextForm[`${resourceName}Enabled`] = true
    nextForm[`${resourceName}TargetType`] = metric.targetType || nextForm[`${resourceName}TargetType`]
    nextForm[`${resourceName}TargetValue`] = metric.targetValue || nextForm[`${resourceName}TargetValue`]
  })

  if (autoscaling?.behavior?.scaleUp?.stabilizationWindowSeconds !== undefined && autoscaling?.behavior?.scaleUp?.stabilizationWindowSeconds !== null) {
    nextForm.scaleUpWindow = Number(autoscaling.behavior.scaleUp.stabilizationWindowSeconds)
  }
  if (autoscaling?.behavior?.scaleDown?.stabilizationWindowSeconds !== undefined && autoscaling?.behavior?.scaleDown?.stabilizationWindowSeconds !== null) {
    nextForm.scaleDownWindow = Number(autoscaling.behavior.scaleDown.stabilizationWindowSeconds)
  }
  if (autoscaling?.behavior?.scaleUp?.selectPolicy) {
    nextForm.scaleUpSelectPolicy = autoscaling.behavior.scaleUp.selectPolicy
  }
  if (autoscaling?.behavior?.scaleDown?.selectPolicy) {
    nextForm.scaleDownSelectPolicy = autoscaling.behavior.scaleDown.selectPolicy
  }

  const scaleUpPodsValue = getBehaviorPolicyValue(autoscaling?.behavior?.scaleUp, 'Pods')
  const scaleUpPercentValue = getBehaviorPolicyValue(autoscaling?.behavior?.scaleUp, 'Percent')
  const scaleDownPodsValue = getBehaviorPolicyValue(autoscaling?.behavior?.scaleDown, 'Pods')
  const scaleDownPercentValue = getBehaviorPolicyValue(autoscaling?.behavior?.scaleDown, 'Percent')

  nextForm.scaleUpPodsEnabled = scaleUpPodsValue !== null
  if (scaleUpPodsValue !== null) nextForm.scaleUpPodsValue = Number(scaleUpPodsValue)
  nextForm.scaleUpPercentEnabled = scaleUpPercentValue !== null
  if (scaleUpPercentValue !== null) nextForm.scaleUpPercentValue = Number(scaleUpPercentValue)
  nextForm.scaleDownPodsEnabled = scaleDownPodsValue !== null
  if (scaleDownPodsValue !== null) nextForm.scaleDownPodsValue = Number(scaleDownPodsValue)
  nextForm.scaleDownPercentEnabled = scaleDownPercentValue !== null
  if (scaleDownPercentValue !== null) nextForm.scaleDownPercentValue = Number(scaleDownPercentValue)

  Object.assign(autoscalingForm, nextForm)
}

const buildBehaviorPolicies = (direction) => {
  const prefix = direction === 'Up' ? 'scaleUp' : 'scaleDown'
  const policies = []

  if (autoscalingForm[`${prefix}PodsEnabled`]) {
    policies.push({
      type: 'Pods',
      value: Number(autoscalingForm[`${prefix}PodsValue`]),
      periodSeconds: 15
    })
  }

  if (autoscalingForm[`${prefix}PercentEnabled`]) {
    policies.push({
      type: 'Percent',
      value: Number(autoscalingForm[`${prefix}PercentValue`]),
      periodSeconds: 15
    })
  }

  return policies
}

const buildAutoscalingPayload = () => {
  const metrics = []
  if (autoscalingForm.cpuEnabled) {
    metrics.push({
      resourceName: 'cpu',
      targetType: autoscalingForm.cpuTargetType,
      targetValue: String(autoscalingForm.cpuTargetValue || '').trim()
    })
  }
  if (autoscalingForm.memoryEnabled) {
    metrics.push({
      resourceName: 'memory',
      targetType: autoscalingForm.memoryTargetType,
      targetValue: String(autoscalingForm.memoryTargetValue || '').trim()
    })
  }

  const behavior = {}
  if (autoscalingForm.scaleUpWindow !== null && autoscalingForm.scaleUpWindow !== '' && autoscalingForm.scaleUpWindow !== undefined) {
    behavior.scaleUp = {
      stabilizationWindowSeconds: Number(autoscalingForm.scaleUpWindow),
      selectPolicy: autoscalingForm.scaleUpSelectPolicy,
      policies: buildBehaviorPolicies('Up')
    }
  }
  if (autoscalingForm.scaleDownWindow !== null && autoscalingForm.scaleDownWindow !== '' && autoscalingForm.scaleDownWindow !== undefined) {
    behavior.scaleDown = {
      stabilizationWindowSeconds: Number(autoscalingForm.scaleDownWindow),
      selectPolicy: autoscalingForm.scaleDownSelectPolicy,
      policies: buildBehaviorPolicies('Down')
    }
  }

  return {
    minReplicas: Number(autoscalingForm.minReplicas),
    maxReplicas: Number(autoscalingForm.maxReplicas),
    metrics,
    ...(Object.keys(behavior).length ? { behavior } : {})
  }
}

const canConfigureAutoscaling = (workload) => canScale(workload)

const autoscalingHasInvalidPolicies = computed(() => {
  const values = []
  if (autoscalingForm.scaleUpPodsEnabled) values.push(Number(autoscalingForm.scaleUpPodsValue))
  if (autoscalingForm.scaleUpPercentEnabled) values.push(Number(autoscalingForm.scaleUpPercentValue))
  if (autoscalingForm.scaleDownPodsEnabled) values.push(Number(autoscalingForm.scaleDownPodsValue))
  if (autoscalingForm.scaleDownPercentEnabled) values.push(Number(autoscalingForm.scaleDownPercentValue))
  return values.some(value => !Number.isFinite(value) || value <= 0)
})

const autoscalingPersistedWarnings = computed(() => currentWorkload.value?.autoscaling?.warnings || [])

const autoscalingNeedsRequestsWarning = computed(() => {
  if (!autoscalingDialogVisible.value) return false
  return (
    (autoscalingForm.cpuEnabled && autoscalingForm.cpuTargetType === 'Utilization' && !hasResourceRequest(currentWorkload.value, 'cpu')) ||
    (autoscalingForm.memoryEnabled && autoscalingForm.memoryTargetType === 'Utilization' && !hasResourceRequest(currentWorkload.value, 'memory'))
  )
})

const formatAutoscalingTarget = (metric = {}) => {
  const resourceLabel = String(metric.resourceName || '').toLowerCase() === 'memory' ? '内存' : 'CPU'
  const targetType = metric.targetType === 'Utilization' ? '利用率' : '平均值'
  const suffix = metric.targetType === 'Utilization' ? '%' : ''
  return `${resourceLabel} ${targetType} ${metric.targetValue || '-'}${suffix}`
}

const formatAutoscalingCurrent = (metric = {}) => {
  if (metric.targetType === 'Utilization') {
    return metric.currentAverageUtilization !== undefined && metric.currentAverageUtilization !== null
      ? `${metric.currentAverageUtilization}%`
      : '-'
  }
  return metric.currentAverageValue || '-'
}

const capacityRiskTagType = (riskLevel) => {
  switch (String(riskLevel || '').toLowerCase()) {
    case 'high':
      return 'danger'
    case 'medium':
      return 'warning'
    case 'low':
      return 'success'
    default:
      return 'info'
  }
}

const capacityRiskText = (riskLevel) => {
  switch (String(riskLevel || '').toLowerCase()) {
    case 'high':
      return '高风险'
    case 'medium':
      return '需关注'
    case 'low':
      return '风险可控'
    default:
      return '-'
  }
}

const governanceStageTagType = (status) => {
  switch (String(status || '').toLowerCase()) {
    case 'blocking':
      return 'danger'
    case 'warning':
      return 'warning'
    case 'ready':
      return 'success'
    case 'pending':
      return 'info'
    default:
      return 'info'
  }
}

const openAutoscalingDialog = async (row) => {
  if (!canConfigureAutoscaling(row)) {
    ElMessage.warning('当前工作负载类型不支持 HPA 自动扩缩容')
    return
  }

  currentWorkload.value = row
  resetAutoscalingForm(row)
  autoscalingDialogVisible.value = true
  autoscalingLoading.value = true

  try {
    const response = await k8sApi.getWorkloadAutoscaling(
      selectedClusterId.value,
      queryParams.namespace,
      row.type,
      row.name
    )

    const responseData = response.data || response
    if (responseData.code !== 200 && !responseData.success) {
      throw new Error(responseData.message || '获取 HPA 配置失败')
    }

    const autoscalingData = responseData.data?.autoscaling || null
    autoscalingExists.value = Boolean(responseData.data?.exists && autoscalingData)
    if (autoscalingExists.value) {
      currentWorkload.value = {
        ...row,
        autoscaling: autoscalingData
      }
      fillAutoscalingForm(autoscalingData)
    }
  } catch (error) {
    autoscalingExists.value = false
    resetAutoscalingForm(row)
    console.error('获取 HPA 配置失败:', error)
    ElMessage.error(error.message || '获取 HPA 配置失败')
  } finally {
    autoscalingLoading.value = false
  }
}

const openCapacitySuggestionDialog = async (row) => {
  currentWorkload.value = row
  capacitySuggestionDialogVisible.value = true
  capacitySuggestionLoading.value = true
  capacitySuggestion.value = createEmptyCapacitySuggestion()
  latestCapacitySuggestion.value = createEmptyCapacitySuggestion()
  capacitySuggestionHistory.value = []
  capacitySuggestionHistoryPagination.pageNum = 1
  capacitySuggestionHistoryPagination.total = 0

  try {
    const response = await k8sApi.getWorkloadCapacitySuggestion(
      selectedClusterId.value,
      queryParams.namespace,
      row.type,
      row.name
    )
    const responseData = response.data || response
    if (responseData.code !== 200 && !responseData.success) {
      throw new Error(responseData.message || '获取容量建议失败')
    }
    latestCapacitySuggestion.value = cloneCapacitySuggestion(responseData.data || createEmptyCapacitySuggestion())
    capacitySuggestion.value = cloneCapacitySuggestion(latestCapacitySuggestion.value)
  } catch (error) {
    console.error('获取容量建议失败:', error)
    ElMessage.error(error.message || '获取容量建议失败')
  } finally {
    capacitySuggestionLoading.value = false
  }

  await fetchCapacitySuggestionHistory(row)
}

const openGovernanceDialog = async (row) => {
  currentWorkload.value = row
  governanceDialogVisible.value = true
  governanceLoading.value = true
  governanceOverview.value = createEmptyGovernanceOverview()

  try {
    const response = await k8sApi.getWorkloadGovernanceOverview(
      selectedClusterId.value,
      queryParams.namespace,
      row.type,
      row.name
    )
    const responseData = response.data || response
    if (responseData.code !== 200 && !responseData.success) {
      throw new Error(responseData.message || '获取治理工作台失败')
    }
    governanceOverview.value = {
      ...createEmptyGovernanceOverview(),
      ...(responseData.data || {})
    }
  } catch (error) {
    console.error('获取治理工作台失败:', error)
    ElMessage.error(error.message || '获取治理工作台失败')
  } finally {
    governanceLoading.value = false
  }
}

const fetchCapacitySuggestionHistory = async (row = currentWorkload.value) => {
  if (!row?.name || !row?.type) return

  capacitySuggestionHistoryLoading.value = true
  try {
    const response = await k8sApi.getWorkloadCapacitySuggestionHistory(
      selectedClusterId.value,
      queryParams.namespace,
      row.type,
      row.name,
      {
        pageNum: capacitySuggestionHistoryPagination.pageNum,
        pageSize: capacitySuggestionHistoryPagination.pageSize
      }
    )
    const responseData = response.data || response
    if (responseData.code !== 200 && !responseData.success) {
      throw new Error(responseData.message || '获取容量建议历史失败')
    }
    const pageData = responseData.data || {}
    capacitySuggestionHistory.value = pageData.list || []
    capacitySuggestionHistoryPagination.total = Number(pageData.total || 0)
  } catch (error) {
    console.error('获取容量建议历史失败:', error)
    ElMessage.error(error.message || '获取容量建议历史失败')
  } finally {
    capacitySuggestionHistoryLoading.value = false
  }
}

const viewCapacitySuggestionHistory = (record) => {
  capacitySuggestion.value = cloneCapacitySuggestion(record || createEmptyCapacitySuggestion())
}

const resetToLatestCapacitySuggestion = () => {
  capacitySuggestion.value = cloneCapacitySuggestion(latestCapacitySuggestion.value)
}

const handleCapacitySuggestionHistorySizeChange = async (size) => {
  capacitySuggestionHistoryPagination.pageSize = size
  capacitySuggestionHistoryPagination.pageNum = 1
  await fetchCapacitySuggestionHistory()
}

const handleCapacitySuggestionHistoryPageChange = async (page) => {
  capacitySuggestionHistoryPagination.pageNum = page
  await fetchCapacitySuggestionHistory()
}

const openAlertCenter = (suggestion = capacitySuggestion.value) => {
  if (suggestion?.alertCenterPath) {
    router.push(suggestion.alertCenterPath)
    return
  }

  const alertCenterQuery = suggestion?.alertCenterQuery || {}
  router.push({
    path: '/monitor/alert-center',
    query: {
      source: alertCenterQuery.source || 'capacity-suggestion',
      keyword: alertCenterQuery.keyword || suggestion?.alertKeyword || currentWorkload.value.name || '',
      namespace: alertCenterQuery.namespace || queryParams.namespace || '',
      workloadName: alertCenterQuery.workloadName || currentWorkload.value.name || ''
    }
  })
}

const openAIDiagnosis = (suggestion = capacitySuggestion.value) => {
  const historyId = Number(suggestion?.historyId || 0)
  if (!historyId) {
    ElMessage.warning('请先生成容量建议快照后再进入 AI 诊断')
    return
  }

  const workloadName = suggestion?.workloadName || currentWorkload.value.name || ''
  const namespaceName = suggestion?.namespace || queryParams.namespace || ''
  router.push({
    path: '/ai/diagnosis',
    query: {
      scene: 'workload_capacity',
      source: 'capacity-suggestion',
      targetId: String(historyId),
      keyword: `${workloadName} hpa capacity autoscaling`.trim(),
      templateName: 'incident_analysis',
      namespace: namespaceName,
      workloadName,
      autoRun: '1'
    }
  })
}

const openGovernanceAlertCenter = () => {
  const path = governanceOverview.value?.alertCenterPath || governanceOverview.value?.currentSuggestion?.alertCenterPath
  if (!path) {
    ElMessage.warning('当前没有可联动的告警中心路径')
    return
  }
  governanceDialogVisible.value = false
  router.push(path)
}

const openGovernanceAIDiagnosis = () => {
  if (governanceOverview.value?.aiDiagnosisPath) {
    governanceDialogVisible.value = false
    router.push(governanceOverview.value.aiDiagnosisPath)
    return
  }
  const latestSnapshot = governanceOverview.value?.latestSnapshot
  if (latestSnapshot?.historyId) {
    governanceDialogVisible.value = false
    openAIDiagnosis(latestSnapshot)
    return
  }
  ElMessage.warning('当前还没有历史容量快照，暂时无法进入 AI 诊断')
}

const openGovernanceCapacity = async () => {
  governanceDialogVisible.value = false
  await nextTick()
  await openCapacitySuggestionDialog(currentWorkload.value)
}

const openGovernanceAutoscaling = async () => {
  governanceDialogVisible.value = false
  await nextTick()
  await openAutoscalingDialog(currentWorkload.value)
}

const openGovernanceAppPath = (item) => {
  if (!item?.detailPath) {
    ElMessage.warning('当前关联应用没有可打开的入口')
    return
  }
  governanceDialogVisible.value = false
  router.push(item.detailPath)
}

const submitAutoscaling = async () => {
  const payload = buildAutoscalingPayload()
  if (payload.maxReplicas < payload.minReplicas) {
    ElMessage.warning('最大副本数不能小于最小副本数')
    return
  }
  if (!payload.metrics.length) {
    ElMessage.warning('请至少启用一个 CPU 或内存指标')
    return
  }
  if (payload.metrics.some(metric => !metric.targetValue)) {
    ElMessage.warning('请补全已启用指标的阈值')
    return
  }
  if (autoscalingHasInvalidPolicies.value) {
    ElMessage.warning('启用中的扩缩容限速策略数值必须大于 0')
    return
  }

  autoscalingSubmitting.value = true
  try {
    const response = await k8sApi.upsertWorkloadAutoscaling(
      selectedClusterId.value,
      queryParams.namespace,
      currentWorkload.value.type,
      currentWorkload.value.name,
      payload
    )

    const responseData = response.data || response
    if (responseData.code !== 200 && !responseData.success) {
      throw new Error(responseData.message || '保存 HPA 配置失败')
    }

    const autoscalingData = responseData.data?.autoscaling || null
    autoscalingExists.value = Boolean(autoscalingData)
    if (autoscalingData) {
      currentWorkload.value = {
        ...currentWorkload.value,
        autoscaling: autoscalingData
      }
    }

    ElMessage.success(responseData.data?.message || 'HPA 配置保存成功')
    if (autoscalingData?.warnings?.length) {
      ElMessage.warning(autoscalingData.warnings[0])
    }
    autoscalingDialogVisible.value = false
    handleQuery()
  } catch (error) {
    console.error('保存 HPA 配置失败:', error)
    ElMessage.error(error.message || '保存 HPA 配置失败')
  } finally {
    autoscalingSubmitting.value = false
  }
}

const deleteAutoscaling = async () => {
  if (!autoscalingExists.value) {
    ElMessage.info('当前工作负载尚未配置 HPA')
    return
  }

  try {
    await ElMessageBox.confirm(
      `确认删除 ${currentWorkload.value.name} 的 HPA 自动扩缩容配置吗？`,
      '删除 HPA',
      {
        confirmButtonText: '确认删除',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    autoscalingSubmitting.value = true
    const response = await k8sApi.deleteWorkloadAutoscaling(
      selectedClusterId.value,
      queryParams.namespace,
      currentWorkload.value.type,
      currentWorkload.value.name
    )

    const responseData = response.data || response
    if (responseData.code !== 200 && !responseData.success) {
      throw new Error(responseData.message || '删除 HPA 配置失败')
    }

    autoscalingExists.value = false
    currentWorkload.value = {
      ...currentWorkload.value,
      autoscaling: null
    }
    ElMessage.success(responseData.data?.message || 'HPA 配置删除成功')
    autoscalingDialogVisible.value = false
    handleQuery()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除 HPA 配置失败:', error)
      ElMessage.error(error.message || '删除 HPA 配置失败')
    }
  } finally {
    autoscalingSubmitting.value = false
  }
}

// 重启工作负载
const restartWorkload = async (row) => {
  if (!canRestart(row)) {
    ElMessage.warning('该工作负载不支持重启操作')
    return
  }
  
  try {
    await ElMessageBox.confirm(
      `确定要重启 ${row.type} "${row.name}" 吗？`,
      '重启确认',
      {
        confirmButtonText: '确定重启',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )
    
    const response = await k8sApi.restartWorkload(
      selectedClusterId.value,
      queryParams.namespace,
      row.type,
      row.name
    )
    
    const responseData = response.data || response
    if (responseData.code === 200 || responseData.success) {
      ElMessage.success(`${row.name} 重启成功`)
      handleQuery()
    } else {
      ElMessage.error(responseData.message || '重启失败')
    }
  } catch (error) {
    if (error === 'cancel') {
      ElMessage.info('已取消重启操作')
    } else {
      console.error('重启失败:', error)
      ElMessage.error('重启失败，请检查网络连接')
    }
  }
}

const getWorkloadTypeName = (type) => {
  const nameMap = {
    'deployment': 'Deployment',
    'statefulset': 'StatefulSet',
    'daemonset': 'DaemonSet',
    'job': 'Job',
    'cronjob': 'CronJob',
    'pod': 'Pod'
  }
  return nameMap[type] || type
}

const normalizeRouteWorkloadType = (type) => {
  const normalized = String(type || '').toLowerCase().trim()
  const typeMap = {
    deployment: 'deployment',
    deployments: 'deployment',
    statefulset: 'statefulset',
    statefulsets: 'statefulset',
    daemonset: 'daemonset',
    daemonsets: 'daemonset',
    job: 'job',
    jobs: 'job',
    cronjob: 'cronjob',
    cronjobs: 'cronjob'
  }
  return typeMap[normalized] || normalized
}

const applyRoutePreset = async () => {
  const clusterId = Number(route.query.clusterId || 0)
  const namespaceName = String(route.query.namespace || '').trim()
  const workloadName = String(route.query.name || '').trim()
  const workloadType = normalizeRouteWorkloadType(route.query.type)
  const action = String(route.query.action || '').trim()

  if (!clusterId && !namespaceName && !workloadName && !action) {
    return false
  }

  if (clusterId > 0) {
    selectedClusterId.value = clusterId
  }
  if (namespaceName) {
    queryParams.namespace = namespaceName
  }

  if (selectedClusterId.value && queryParams.namespace) {
    await handleQuery()
  }

  if (!workloadName || !workloadType) {
    return true
  }

  const targetRow = (tableData.value || []).find(row =>
    row.name === workloadName && normalizeRouteWorkloadType(row.type) === workloadType
  )
  if (!targetRow) {
    return true
  }

  if (action === 'capacity') {
    await openCapacitySuggestionDialog(targetRow)
  } else if (action === 'governance') {
    await openGovernanceDialog(targetRow)
  } else if (action === 'autoscaling') {
    await openAutoscalingDialog(targetRow)
  }

  return true
}

// 复制内容到剪贴板
const copyToClipboard = async (text, successMessage = '已复制到剪贴板') => {
  try {
    await navigator.clipboard.writeText(text)
    ElMessage.success(successMessage)
  } catch (error) {
    console.error('复制失败:', error)
    const textArea = document.createElement('textarea')
    textArea.value = text
    document.body.appendChild(textArea)
    textArea.select()
    try {
      document.execCommand('copy')
      ElMessage.success(successMessage)
    } catch (fallbackError) {
      ElMessage.error('复制失败，请手动复制')
    }
    document.body.removeChild(textArea)
  }
}

// 判断是否为系统命名空间或系统工作负载
const isSystemWorkload = (workload) => {
  // 系统命名空间
  const systemNamespaces = ['kube-system', 'kube-public', 'kube-node-lease', 'calico-system', 'tigera-operator']
  if (systemNamespaces.includes(workload.namespace)) {
    return true
  }
  
  // 系统工作负载名称前缀
  const systemPrefixes = ['kube-', 'calico-', 'coredns', 'metrics-server', 'node-local-dns', 'kubernetes-dashboard']
  if (systemPrefixes.some(prefix => workload.name.startsWith(prefix))) {
    return true
  }
  
  return false
}

// 判断是否可以扩缩容
const canScale = (workload) => {
  // 只有 Deployment 和 StatefulSet 支持扩缩容
  if (!['deployment', 'statefulset'].includes(workload.type)) {
    return false
  }
  
  // 系统工作负载不允许扩缩容
  return !isSystemWorkload(workload)
}

// 判断是否可以重启
const canRestart = (workload) => {
  // Deployment / StatefulSet / DaemonSet 支持重启
  if (!['deployment', 'statefulset', 'daemonset'].includes(workload.type)) {
    return false
  }
  
  // 系统工作负载不允许重启
  return !isSystemWorkload(workload)
}

// 判断是否可以删除
const canDelete = (workload) => {
  // 系统工作负载不允许删除
  return !isSystemWorkload(workload)
}

// 判断是否可以更新Pod配置
const canUpdatePodConfig = (workload) => {
  // 只有 Deployment 和 StatefulSet 支持Pod配置更新
  if (!['deployment', 'statefulset'].includes(workload.type)) {
    return false
  }
  
  // 系统工作负载不允许配置更新
  return !isSystemWorkload(workload)
}

// 判断是否可以更新调度
const canUpdateScheduling = (workload) => {
  // 只有 Deployment、StatefulSet 支持调度更新
  if (!['deployment', 'statefulset'].includes(workload.type)) {
    return false
  }
  
  // 系统工作负载不允许调度更新
  return !isSystemWorkload(workload)
}

// 判断是否可以编辑YAML
const canEditYaml = (workload) => {
  // 系统工作负载不允许编辑YAML
  return !isSystemWorkload(workload)
}

// 获取可见标签数量（排除系统标签）
const getVisibleLabelCount = (labels) => {
  if (!labels) return 0
  
  const systemLabelPrefixes = [
    'kubernetes.io/',
    'beta.kubernetes.io/',
    'node-role.kubernetes.io/',
    'node.kubernetes.io/',
    'app.kubernetes.io/managed-by',
    'pod-template-hash'
  ]
  
  return Object.keys(labels).filter(key => 
    !systemLabelPrefixes.some(prefix => key.startsWith(prefix))
  ).length
}

// 根据副本数获取Pod状态标签类型
const getPodStatusTagByReplicas = (readyReplicas, totalReplicas) => {
  if (totalReplicas === 0) return 'info'
  if (readyReplicas === 0) return 'danger'
  if (readyReplicas < totalReplicas) return 'warning'
  return 'success'
}

// 获取副本状态文本
const getReplicaStatusText = (readyReplicas, totalReplicas) => {
  if (totalReplicas === 0) return '已停止'
  if (readyReplicas === 0) return '启动中'
  if (readyReplicas < totalReplicas) return '部分就绪'
  return '全部就绪'
}

// 获取副本状态样式类
const getReplicaStatusClass = (readyReplicas, totalReplicas) => {
  if (totalReplicas === 0) return 'status-stopped'
  if (readyReplicas === 0) return 'status-starting'
  if (readyReplicas < totalReplicas) return 'status-partial'
  return 'status-ready'
}

// 查看工作负载标签
const viewWorkloadLabels = (row) => {
  currentWorkload.value = row
  workloadLabelsDialogVisible.value = true
}

// 查看所有镜像
const viewAllImages = (row) => {
  currentWorkload.value = row
  allImagesDialogVisible.value = true
}

// 更新Pod配置
const updatePodConfig = async (row) => {
  if (!canUpdatePodConfig(row)) {
    ElMessage.warning('该工作负载不支持Pod配置更新')
    return
  }

  await openPodConfigDialog(row)
}

// 更新调度
const updateScheduling = async (row) => {
  if (!canUpdateScheduling(row)) {
    ElMessage.warning('该工作负载不支持调度更新')
    return
  }

  await openSchedulingDialog(row)
}

// 编辑工作负载YAML
const editWorkloadYaml = async (row) => {
  if (!canEditYaml(row)) {
    ElMessage.warning('系统工作负载不允许编辑YAML')
    return
  }
  
  try {
    loading.value = true
    console.log('🔍 开始获取工作负载YAML...', row)
    
    // 首先尝试直接获取工作负载YAML
    try {
      const response = await k8sApi.getWorkloadYaml(selectedClusterId.value, queryParams.namespace, row.type, row.name)
      const responseData = response.data || response
      
      if (responseData.code === 200 || responseData.success) {
        currentWorkload.value = row
        // 检查返回的数据结构并正确处理
        let yamlContent = ''
        if (responseData.data && responseData.data.yamlContent) {
          // 新的API返回yamlContent字段，直接使用
          yamlContent = responseData.data.yamlContent
        } else if (responseData.data && responseData.data.yaml) {
          // 兼容旧的yaml字段，将对象转换为YAML字符串
          try {
            yamlContent = yaml.dump(responseData.data.yaml, { indent: 2, lineWidth: -1 })
          } catch (error) {
            console.error('YAML转换失败:', error)
            yamlContent = JSON.stringify(responseData.data.yaml, null, 2)
          }
        } else if (typeof responseData.data === 'string') {
          yamlContent = responseData.data
        } else {
          yamlContent = `apiVersion: apps/v1\nkind: ${row.type}\nmetadata:\n  name: ${row.name}`
        }
        currentYaml.value = yamlContent
        workloadYamlDialogVisible.value = true
        return
      }
    } catch (workloadError) {
      console.log('工作负载YAML获取失败，尝试通过Pod获取:', workloadError)
    }
    
    // 如果工作负载YAML获取失败，尝试获取该工作负载下的Pod列表
    console.log('🔍 通过获取工作负载Pod列表来获取Pod YAML...')
    const detailResponse = await k8sApi.getWorkloadPods(selectedClusterId.value, queryParams.namespace, row.type.toLowerCase(), row.name)
    const detailData = detailResponse.data || detailResponse

    if ((detailData.code === 200 || detailData.success) && detailData.data && detailData.data.length > 0) {
      // 获取第一个Pod的YAML
      const firstPod = detailData.data[0]
      console.log('🔍 获取Pod YAML:', firstPod.name)
      
      const podYamlResponse = await k8sApi.getPodYaml(selectedClusterId.value, queryParams.namespace, firstPod.name)
      const podYamlData = podYamlResponse.data || podYamlResponse
      
      if (podYamlData.code === 200) {
        currentWorkload.value = row
        // 检查返回的数据结构并正确处理
        let yamlContent = ''
        if (podYamlData.data && podYamlData.data.yamlContent) {
          // 新的API返回yamlContent字段，直接使用
          yamlContent = podYamlData.data.yamlContent
        } else if (podYamlData.data && podYamlData.data.yaml) {
          // 兼容旧的yaml字段，将对象转换为YAML字符串
          try {
            yamlContent = yaml.dump(podYamlData.data.yaml, { indent: 2, lineWidth: -1 })
          } catch (error) {
            console.error('YAML转换失败:', error)
            yamlContent = JSON.stringify(podYamlData.data.yaml, null, 2)
          }
        } else if (typeof podYamlData.data === 'string') {
          yamlContent = podYamlData.data
        } else {
          yamlContent = `# Pod YAML for ${firstPod.name}\n# 工作负载: ${row.name} (${row.type})\n`
        }
        currentYaml.value = yamlContent
        workloadYamlDialogVisible.value = true
      } else {
        throw new Error(podYamlData.message || 'Pod YAML获取失败')
      }
    } else {
      // 如果没有Pod，生成基础的工作负载YAML模板
      console.log('⚠️ 没有找到Pod，生成基础YAML模板')
      const templateYaml = generateWorkloadYamlTemplate(row)
      currentWorkload.value = row
      currentYaml.value = templateYaml
      workloadYamlDialogVisible.value = true
      ElMessage.warning('未找到实际YAML内容，显示基础模板')
    }
    
  } catch (error) {
    console.error('获取工作负载YAML失败:', error)
    ElMessage.error('获取工作负载YAML失败: ' + (error.message || '请检查网络连接'))
  } finally {
    loading.value = false
  }
}

// 生成工作负载YAML模板
const generateWorkloadYamlTemplate = (workload) => {
  const kind = workload.type.charAt(0).toUpperCase() + workload.type.slice(1)
  return `apiVersion: apps/v1
kind: ${kind}
metadata:
  name: ${workload.name}
  namespace: ${queryParams.namespace}
  labels:
    app: ${workload.name}
spec:
  replicas: ${workload.totalReplicas || 1}
  selector:
    matchLabels:
      app: ${workload.name}
  template:
    metadata:
      labels:
        app: ${workload.name}
    spec:
      containers:
      - name: ${workload.name}
        image: nginx:latest
        ports:
        - containerPort: 80
---
# 注意: 这是一个基础模板，请根据实际需求修改
# 工作负载类型: ${workload.type}
# 当前状态: ${workload.status}
# 副本数: ${workload.replicas}`
}

onMounted(async () => {
  try {
    console.log('🚀 开始加载k8s工作负载页面')
    const startTime = Date.now()
    
    // 加载集群列表
    console.log('📡 正在加载集群列表...')
    await fetchClusterList()
    console.log('✅ 集群列表加载完成，耗时:', Date.now() - startTime + 'ms')
    
    let routePresetApplied = false
    if (selectedClusterId.value) {
      console.log('🔄 开始并行加载命名空间和工作负载数据')

      // 命名空间加载已移至 NamespaceSelector 组件
      // 路由预设必须优先于首次查询，否则 default namespace 会覆盖 query preset。
      routePresetApplied = await applyRoutePreset()

      if (!routePresetApplied && queryParams.namespace) {
        console.log('📦 立即开始查询工作负载:', queryParams.namespace)
        handleQuery().catch(error => {
          console.error('工作负载初始查询失败:', error)
        })
      }
    } else {
      routePresetApplied = await applyRoutePreset()
    }

    console.log('🎉 页面初始化完成，总耗时:', Date.now() - startTime + 'ms')
  } catch (error) {
    console.error('页面初始化失败:', error)
  }
})

// 监听YAML对话框的打开状态，自动聚焦编辑器
watch(tableData, () => {
  applyWorkloadTableLayout()
}, { deep: true, flush: 'post' })

watch(workloadYamlDialogVisible, (newVal) => {
  if (newVal) {
    nextTick(() => {
      if (yamlEditor.value && yamlEditor.value.focus) {
        yamlEditor.value.focus()
      }
    })
  }
})


// 辅助函数
const formatCpu = (cpuStr) => {
  if (!cpuStr || cpuStr === '0' || cpuStr === '') return '-'
  return cpuStr
}

const formatMemory = (memoryStr) => {
  if (!memoryStr || memoryStr === '0' || memoryStr === '') return '-'

  if (memoryStr.endsWith('Ki')) {
    const kb = parseInt(memoryStr.replace('Ki', ''))
    if (kb < 1024) return memoryStr
    const mb = (kb / 1024).toFixed(1)
    return `${mb}Mi`
  }

  if (memoryStr.endsWith('Mi')) {
    const mb = parseInt(memoryStr.replace('Mi', ''))
    if (mb < 1024) return memoryStr
    const gb = (mb / 1024).toFixed(1)
    return `${gb}Gi`
  }

  if (memoryStr.endsWith('Gi')) {
    return memoryStr
  }

  // 如果没有单位，假设是字节
  const bytes = parseInt(memoryStr)
  if (!isNaN(bytes) && bytes > 0) {
    if (bytes < 1024) return `${bytes}B`
    if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)}Ki`
    if (bytes < 1024 * 1024 * 1024) return `${(bytes / (1024 * 1024)).toFixed(1)}Mi`
    return `${(bytes / (1024 * 1024 * 1024)).toFixed(1)}Gi`
  }

  return memoryStr || '-'
}

const getImageTag = (image) => {
  if (!image) return 'latest'
  const parts = image.split(':')
  return parts.length > 1 ? parts[parts.length - 1] : 'latest'
}

const getImageRegistry = (image) => {
  if (!image) return 'docker.io'
  const parts = image.split('/')
  if (parts.length === 1) return 'docker.io'
  if (parts[0].includes('.') || parts[0].includes(':')) {
    return parts[0]
  }
  return 'docker.io'
}

// 获取用户自定义标签
const getUserLabels = (labels) => {
  if (!labels) return {}

  const systemLabelPrefixes = [
    'kubernetes.io/',
    'beta.kubernetes.io/',
    'node-role.kubernetes.io/',
    'node.kubernetes.io/',
    'app.kubernetes.io/managed-by',
    'pod-template-hash'
  ]

  const userLabels = {}
  Object.entries(labels).forEach(([key, value]) => {
    if (!systemLabelPrefixes.some(prefix => key.startsWith(prefix))) {
      userLabels[key] = value
    }
  })

  return userLabels
}

// 获取系统标签
const getSystemLabels = (labels) => {
  if (!labels) return {}

  const systemLabelPrefixes = [
    'kubernetes.io/',
    'beta.kubernetes.io/',
    'node-role.kubernetes.io/',
    'node.kubernetes.io/',
    'app.kubernetes.io/managed-by',
    'pod-template-hash'
  ]

  const systemLabels = {}
  Object.entries(labels).forEach(([key, value]) => {
    if (systemLabelPrefixes.some(prefix => key.startsWith(prefix))) {
      systemLabels[key] = value
    }
  })

  return systemLabels
}

const extractYamlContent = (responseData, fallback = '') => {
  if (responseData?.data?.yamlContent) {
    return responseData.data.yamlContent
  }
  if (responseData?.data?.yaml) {
    return yaml.dump(responseData.data.yaml, { indent: 2, lineWidth: -1, noRefs: true })
  }
  if (typeof responseData?.data === 'string') {
    return responseData.data
  }
  return fallback
}

const ensureWorkloadTemplateSpec = (manifest) => {
  if (!manifest.spec) {
    manifest.spec = {}
  }
  if (!manifest.spec.template) {
    manifest.spec.template = {}
  }
  if (!manifest.spec.template.spec) {
    manifest.spec.template.spec = {}
  }
  return manifest.spec.template.spec
}

const ensurePrimaryWorkloadContainer = (manifest) => {
  const templateSpec = ensureWorkloadTemplateSpec(manifest)
  if (!Array.isArray(templateSpec.containers) || templateSpec.containers.length === 0) {
    throw new Error('未找到可更新的容器配置')
  }
  return templateSpec.containers[0]
}

const getPrimaryWorkloadContainer = (workloadDetail) => {
  return workloadDetail?.spec?.template?.spec?.containers?.[0]
    || workloadDetail?.pods?.[0]?.containers?.[0]
    || null
}

const formatNodeSelector = (nodeSelector = {}) => {
  return Object.entries(nodeSelector)
    .map(([key, value]) => `${key}=${value}`)
    .join(', ')
}

const parseNodeSelectorText = (text = '') => {
  const selector = {}
  const parts = text
    .split(',')
    .map(part => part.trim())
    .filter(Boolean)

  parts.forEach(part => {
    const separatorIndex = part.indexOf('=')
    if (separatorIndex <= 0 || separatorIndex === part.length - 1) {
      throw new Error(`节点选择器格式错误: ${part}`)
    }

    const key = part.slice(0, separatorIndex).trim()
    const value = part.slice(separatorIndex + 1).trim()
    if (!key || !value) {
      throw new Error(`节点选择器格式错误: ${part}`)
    }

    selector[key] = value
  })

  return selector
}

const detectNodeAffinityMode = (affinity = {}) => {
  if (affinity?.nodeAffinity?.requiredDuringSchedulingIgnoredDuringExecution?.nodeSelectorTerms?.length) {
    return 'required'
  }
  if (affinity?.nodeAffinity?.preferredDuringSchedulingIgnoredDuringExecution?.length) {
    return 'preferred'
  }
  return 'none'
}

const hasPodAntiAffinity = (affinity = {}) => {
  return Boolean(
    affinity?.podAntiAffinity?.requiredDuringSchedulingIgnoredDuringExecution?.length
    || affinity?.podAntiAffinity?.preferredDuringSchedulingIgnoredDuringExecution?.length
  )
}

const buildNodeAffinity = (nodeSelector, mode) => {
  if (mode === 'none' || Object.keys(nodeSelector).length === 0) {
    return null
  }

  const matchExpressions = Object.entries(nodeSelector).map(([key, value]) => ({
    key,
    operator: 'In',
    values: [value]
  }))

  if (mode === 'required') {
    return {
      requiredDuringSchedulingIgnoredDuringExecution: {
        nodeSelectorTerms: [{ matchExpressions }]
      }
    }
  }

  return {
    preferredDuringSchedulingIgnoredDuringExecution: [{
      weight: 100,
      preference: {
        matchExpressions
      }
    }]
  }
}

const buildPodAntiAffinity = (manifest) => {
  const labels = manifest?.spec?.template?.metadata?.labels
    || manifest?.spec?.selector?.matchLabels
    || {}

  if (Object.keys(labels).length === 0) {
    throw new Error('当前工作负载缺少标签，无法生成 Pod 反亲和规则')
  }

  return {
    preferredDuringSchedulingIgnoredDuringExecution: [{
      weight: 100,
      podAffinityTerm: {
        labelSelector: {
          matchLabels: labels
        },
        topologyKey: 'kubernetes.io/hostname'
      }
    }]
  }
}

const updateContainerResources = (container, updateData) => {
  const requests = {}
  const limits = {}

  if (updateData.cpuRequest) {
    requests.cpu = updateData.cpuRequest
  }
  if (updateData.memoryRequest) {
    requests.memory = updateData.memoryRequest
  }
  if (updateData.cpuLimit) {
    limits.cpu = updateData.cpuLimit
  }
  if (updateData.memoryLimit) {
    limits.memory = updateData.memoryLimit
  }

  if (Object.keys(requests).length === 0 && Object.keys(limits).length === 0) {
    delete container.resources
    return
  }

  const resources = { ...(container.resources || {}) }
  if (Object.keys(requests).length > 0) {
    resources.requests = requests
  } else {
    delete resources.requests
  }

  if (Object.keys(limits).length > 0) {
    resources.limits = limits
  } else {
    delete resources.limits
  }

  container.resources = resources
}

const loadWorkloadDetail = async (row) => {
  const response = await k8sApi.getWorkloadDetail(selectedClusterId.value, queryParams.namespace, row.type, row.name)
  const responseData = response.data || response

  if (responseData.code !== 200 && !responseData.success) {
    throw new Error(responseData.message || '获取工作负载详情失败')
  }

  return responseData.data || {}
}

const updateWorkloadViaYaml = async (row, patchFn) => {
  const response = await k8sApi.getWorkloadYaml(selectedClusterId.value, queryParams.namespace, row.type, row.name)
  const responseData = response.data || response
  const yamlContent = extractYamlContent(responseData)

  if (!yamlContent) {
    throw new Error('获取工作负载 YAML 失败')
  }

  const manifest = yaml.load(yamlContent)
  if (!manifest || typeof manifest !== 'object') {
    throw new Error('工作负载 YAML 内容无效')
  }

  const updatedManifest = patchFn(manifest) || manifest
  const nextYaml = yaml.dump(updatedManifest, { indent: 2, lineWidth: -1, noRefs: true })
  const updateResponse = await k8sApi.updateWorkloadYaml(
    selectedClusterId.value,
    queryParams.namespace,
    row.type,
    row.name,
    nextYaml
  )
  const updateResponseData = updateResponse.data || updateResponse

  if (updateResponseData.code !== 200 && !updateResponseData.success) {
    throw new Error(updateResponseData.message || '更新工作负载 YAML 失败')
  }

  return updateResponseData
}

const openPodConfigDialog = async (row) => {
  try {
    loading.value = true
    const detail = await loadWorkloadDetail(row)
    const primaryContainer = getPrimaryWorkloadContainer(detail)

    if (!primaryContainer) {
      throw new Error('未找到可配置的容器')
    }

    currentWorkload.value = {
      ...row,
      containerName: primaryContainer.name,
      image: primaryContainer.image,
      ports: primaryContainer.ports || [],
      envVars: (primaryContainer.env || []).map(env => ({
        name: env.name || '',
        value: env.value ?? ''
      })),
      resources: {
        requests: {
          cpu: primaryContainer.resources?.requests?.cpu || '',
          memory: primaryContainer.resources?.requests?.memory || ''
        },
        limits: {
          cpu: primaryContainer.resources?.limits?.cpu || '',
          memory: primaryContainer.resources?.limits?.memory || ''
        }
      }
    }
    podConfigDialogVisible.value = true
  } catch (error) {
    console.error('获取 Pod 配置详情失败:', error)
    ElMessage.error(`获取 Pod 配置详情失败: ${error.message || '请检查网络连接'}`)
  } finally {
    loading.value = false
  }
}

const openSchedulingDialog = async (row) => {
  try {
    loading.value = true
    const detail = await loadWorkloadDetail(row)
    const templateSpec = detail?.spec?.template?.spec || {}
    const affinity = templateSpec.affinity || {}

    currentWorkload.value = row
    schedulingForm.nodeSelectorText = formatNodeSelector(templateSpec.nodeSelector || {})
    schedulingForm.nodeAffinity = detectNodeAffinityMode(affinity)
    schedulingForm.podAntiAffinity = hasPodAntiAffinity(affinity)
    schedulingDialogVisible.value = true
  } catch (error) {
    console.error('获取调度配置详情失败:', error)
    ElMessage.error(`获取调度配置详情失败: ${error.message || '请检查网络连接'}`)
  } finally {
    loading.value = false
  }
}

// 处理Pod配置提交
const handlePodConfigSubmit = async (updateData) => {
  console.log('🔧 开始处理Pod配置更新:', {
    clusterId: selectedClusterId.value,
    namespace: queryParams.namespace,
    workloadName: currentWorkload.value.name,
    updateData: updateData
  })

  try {
    const responseData = await updateWorkloadViaYaml(currentWorkload.value, (manifest) => {
      const container = ensurePrimaryWorkloadContainer(manifest)
      updateContainerResources(container, updateData)

      if (updateData.envVars.length > 0) {
        container.env = updateData.envVars.map(env => ({
          name: env.name,
          value: env.value ?? ''
        }))
      } else {
        delete container.env
      }

      return manifest
    })

    console.log('📤 API响应:', responseData)
    ElMessage.success(`${currentWorkload.value.name} Pod配置更新成功`)
    podConfigDialogVisible.value = false
    handleQuery() // 刷新列表
  } catch (error) {
    console.error('💥 Pod配置更新异常:', error)
    ElMessage.error(`Pod配置更新失败: ${error.message || '请检查网络连接'}`)
  }
}

// 提交调度配置
const submitScheduling = async () => {
  try {
    const nodeSelector = parseNodeSelectorText(schedulingForm.nodeSelectorText)
    if (schedulingForm.nodeAffinity !== 'none' && Object.keys(nodeSelector).length === 0) {
      throw new Error('配置节点亲和性时请先填写节点选择器')
    }
    const nodeAffinity = buildNodeAffinity(nodeSelector, schedulingForm.nodeAffinity)

    await updateWorkloadViaYaml(currentWorkload.value, (manifest) => {
      const templateSpec = ensureWorkloadTemplateSpec(manifest)
      const affinity = { ...(templateSpec.affinity || {}) }

      if (Object.keys(nodeSelector).length > 0) {
        templateSpec.nodeSelector = nodeSelector
      } else {
        delete templateSpec.nodeSelector
      }

      if (nodeAffinity) {
        affinity.nodeAffinity = nodeAffinity
      } else {
        delete affinity.nodeAffinity
      }

      if (schedulingForm.podAntiAffinity) {
        affinity.podAntiAffinity = buildPodAntiAffinity(manifest)
      } else {
        delete affinity.podAntiAffinity
      }

      if (Object.keys(affinity).length > 0) {
        templateSpec.affinity = affinity
      } else {
        delete templateSpec.affinity
      }

      return manifest
    })

    ElMessage.success(`${currentWorkload.value.name} 调度配置更新成功`)
    schedulingDialogVisible.value = false
    handleQuery() // 刷新列表
  } catch (error) {
    console.error('调度配置更新失败:', error)
    ElMessage.error(`调度配置更新失败: ${error.message || '请检查网络连接'}`)
  }
}

// 从YAML内容中解析工作负载信息
const parseWorkloadFromYaml = (yamlContent) => {
  try {
    const yamlObj = yaml.load(yamlContent)
    if (!yamlObj || typeof yamlObj !== 'object') {
      throw new Error('无效的YAML格式')
    }

    const workloadType = yamlObj.kind?.toLowerCase()
    const workloadName = yamlObj.metadata?.name

    if (!workloadType || !workloadName) {
      throw new Error('YAML中缺少kind或metadata.name字段')
    }

    return { workloadType, workloadName }
  } catch (error) {
    console.error('解析YAML失败:', error)
    throw error
  }
}

// 处理YAML保存事件
const handleYamlSave = async (data) => {
  try {
    // 更新currentYaml内容
    currentYaml.value = data.yamlContent

    if (data.resourceType === 'Pod') {
      // 对于Pod，使用现有的Pod更新API
      const response = await k8sApi.updatePodYaml(
        selectedClusterId.value,
        queryParams.namespace,
        data.resourceName,
        data.yamlContent
      )

      const responseData = response.data || response
      if (responseData.code === 200 || responseData.success) {
        ElMessage.success(`${data.resourceName} YAML配置保存成功`)
        podYamlDialogVisible.value = false
        podYamlDialogEditable.value = false
      } else {
        throw new Error(responseData.message || '保存失败')
      }
    } else {
      // 对于工作负载，使用新的通用API
      // 从YAML内容中解析实际的工作负载名称和类型
      const { workloadType, workloadName } = parseWorkloadFromYaml(data.yamlContent)
      const yamlSaveOptions = {
        force: workloadType === 'job'
      }

      const response = await k8sApi.updateWorkloadYaml(
        selectedClusterId.value,
        queryParams.namespace,
        workloadType,
        workloadName,
        data.yamlContent,
        yamlSaveOptions
      )

      const responseData = response.data || response
      if (responseData.code === 200 || responseData.success) {
        ElMessage.success(`${workloadName} YAML配置保存成功`)
        workloadYamlDialogVisible.value = false
      } else {
        throw new Error(responseData.message || '保存失败')
      }
    }

    handleQuery() // 刷新列表
  } catch (error) {
    console.error('YAML配置保存失败:', error)
    ElMessage.error('YAML配置保存失败: ' + (error.message || '请检查网络连接'))
  }
}

const deleteWorkload = async (row) => {
  if (!canDelete(row)) {
    ElMessage.warning('系统工作负载不允许删除')
    return
  }

  try {
    await ElMessageBox.confirm(
      `确定要删除 ${getWorkloadTypeName(row.type)} "${row.name}" 吗？此操作不可恢复！`,
      '删除确认',
      {
        confirmButtonText: '确定删除',
        cancelButtonText: '取消',
        type: 'warning',
        dangerouslyUseHTMLString: false
      }
    )

    console.log('🗑️ 开始删除工作负载:', {
      clusterId: selectedClusterId.value,
      namespace: queryParams.namespace,
      workloadType: row.type.toLowerCase(),
      workloadName: row.name
    })

    let response

    // 根据类型选择不同的删除API
    switch (row.type.toLowerCase()) {
      case 'pod':
        // Pod使用专门的删除Pod API
        response = await k8sApi.deletePod(selectedClusterId.value, queryParams.namespace, row.name)
        break
      case 'deployment':
        // Deployment使用专门的删除Deployment API
        response = await k8sApi.deleteDeployment(selectedClusterId.value, queryParams.namespace, row.name)
        break
      default:
        // 其他工作负载使用通用工作负载API
        response = await k8sApi.deleteWorkloadByType(
          selectedClusterId.value,
          queryParams.namespace,
          row.type,
          row.name
        )
        break
    }

    console.log('📤 删除API响应:', response)

    const responseData = response.data || response
    if (responseData.code === 200 || responseData.success) {
      ElMessage.success(`${getWorkloadTypeName(row.type)} "${row.name}" 删除成功`)
      // 刷新工作负载列表
      await handleQuery()
    } else {
      ElMessage.error(responseData.message || `删除 ${getWorkloadTypeName(row.type)} 失败`)
    }
  } catch (error) {
    if (error === 'cancel') {
      ElMessage.info('已取消删除操作')
    } else {
      console.error('删除工作负载失败:', error)
      ElMessage.error(`删除失败: ${error.message || '请检查网络连接'}`)
    }
  }
}
</script>

<template>
  <div class="k8s-workloads-management">
    <el-card shadow="hover" class="workloads-card">
      <template #header>
        <div class="card-header">
          <span class="title">K8s 工作负载管理</span>
          <div class="header-actions">
            <ClusterSelector
              v-model="selectedClusterId"
              @change="handleClusterChange"
            />
            
            <NamespaceSelector
              v-model="queryParams.namespace"
              :cluster-id="selectedClusterId"
              @change="handleNamespaceChange"
            />
          </div>
        </div>
      </template>
      
      <!-- 搜索表单 -->
      <div class="search-section">
        <el-form :inline="true" :model="queryParams" class="search-form">
          <el-form-item label="工作负载名称">
            <el-input
              v-model="queryParams.name"
              placeholder="请输入名称"
              clearable
              size="small"
              style="width: 200px"
              @keyup.enter="handleQuery"
            />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" :icon="Search" size="small" @click="handleQuery">
              搜索
            </el-button>
            <el-button :icon="Refresh" size="small" @click="resetQuery">
              重置
            </el-button>
            <el-button :icon="Monitor" type="success" size="small" @click="navigateToMonitoring">
              监控仪表板
            </el-button>
            <el-button :icon="Plus" v-authority="['k8s:workload:add']" type="primary" size="small" @click="showCreatePodDialog">
              创建工作负载
            </el-button>
          </el-form-item>
        </el-form>

        <!-- 工作负载类型标签页 -->
        <div class="workload-type-section">
          <el-tabs v-model="activeTab" @tab-change="handleTabChange" class="workload-tabs">
            <el-tab-pane label="全部" name="">
              <template #label>
                <span class="tab-label">全部</span>
              </template>
            </el-tab-pane>
            <el-tab-pane label="Deployment" name="deployments">
              <template #label>
                <span class="tab-label">Deployment</span>
              </template>
            </el-tab-pane>
            <el-tab-pane label="StatefulSet" name="statefulsets">
              <template #label>
                <span class="tab-label">StatefulSet</span>
              </template>
            </el-tab-pane>
            <el-tab-pane label="DaemonSet" name="daemonsets">
              <template #label>
                <span class="tab-label">DaemonSet</span>
              </template>
            </el-tab-pane>
            <el-tab-pane label="Job" name="jobs">
              <template #label>
                <span class="tab-label">Job</span>
              </template>
            </el-tab-pane>
            <el-tab-pane label="CronJob" name="cronjobs">
              <template #label>
                <span class="tab-label">CronJob</span>
              </template>
            </el-tab-pane>
          </el-tabs>
        </div>
      </div>

      <!-- 工作负载列表表格 -->
      <div class="table-section">
        <el-table
          :data="tableData"
          v-loading="loading"
          stripe
          style="width: 100%"
          class="workloads-table"
        >
          <el-table-column prop="name" label="名称" min-width="200">
            <template #default="{ row }">
              <div class="workload-name-container">
                <img src="@/assets/image/k8s.svg" alt="k8s" class="k8s-icon" />
                <div class="workload-info">
                  <div 
                    class="workload-name clickable-name" 
                    @click="navigateToPodDetail(row)"
                  >
                    {{ row.name }}
                  </div>
                  <span
                    class="workload-type-label"
                  >
                    {{ getWorkloadTypeName(row.type) }}
                  </span>
                </div>
              </div>
            </template>
          </el-table-column>
          
          <el-table-column label="标签" min-width="100" align="center">
            <template #default="{ row }">
              <div class="label-container">
                <el-badge :value="getVisibleLabelCount(row.labels)" :max="99" class="label-badge">
                  <el-button
                    type="text"
                    size="small"
                    circle
                    @click="viewWorkloadLabels(row)"
                    class="label-icon-button"
                  >
                    <img src="@/assets/image/标签.svg" alt="标签" width="14" height="14" />
                  </el-button>
                </el-badge>
              </div>
            </template>
          </el-table-column>
          
          <el-table-column label="容器组数量" min-width="120" align="center">
            <template #default="{ row }">
              <div class="pod-status-container">
                <el-tag 
                  :type="getPodStatusTagByReplicas(row.readyReplicas, row.totalReplicas)" 
                  size="default"
                  class="pod-count-tag"
                  @click="viewPodList(row)"
                >
                  <el-icon class="pod-icon"><Monitor /></el-icon>
                  {{ row.replicas }}
                </el-tag>
                <div class="pod-status-text">
                  <span :class="getReplicaStatusClass(row.readyReplicas, row.totalReplicas)">
                    {{ getReplicaStatusText(row.readyReplicas, row.totalReplicas) }}
                  </span>
                </div>
              </div>
            </template>
          </el-table-column>
          
          <el-table-column label="Request/Limits" width="130">
            <template #default="{ row }">
              <div class="resource-info">
                <div class="resource-row">
                  <span class="resource-type">CPU:</span>
                  <span class="resource-values">
                    <span class="request-value">{{ formatCpu(row.resources?.requests?.cpu) }}</span>
                    <span class="separator">/</span>
                    <span class="limit-value">{{ formatCpu(row.resources?.limits?.cpu) }}</span>
                  </span>
                </div>
                <div class="resource-row">
                  <span class="resource-type">Mem:</span>
                  <span class="resource-values">
                    <span class="request-value">{{ formatMemory(row.resources?.requests?.memory) }}</span>
                    <span class="separator">/</span>
                    <span class="limit-value">{{ formatMemory(row.resources?.limits?.memory) }}</span>
                  </span>
                </div>
              </div>
            </template>
          </el-table-column>

          <el-table-column label="HPA" width="150">
            <template #default="{ row }">
              <div v-if="row.autoscaling?.enabled" class="autoscaling-cell">
                <div class="autoscaling-head">
                  <el-tag type="success" size="small">已启用</el-tag>
                  <span class="autoscaling-range">{{ row.autoscaling.minReplicas }} - {{ row.autoscaling.maxReplicas }}</span>
                </div>
                <div class="autoscaling-metrics">
                  <el-tag
                    v-for="metric in row.autoscaling.metrics || []"
                    :key="`${row.name}-${metric.resourceName}`"
                    type="info"
                    effect="plain"
                    size="small"
                    class="autoscaling-metric-tag"
                  >
                    {{ formatAutoscalingTarget(metric) }}
                  </el-tag>
                </div>
                <div class="autoscaling-status-line">
                  当前/期望副本: {{ row.autoscaling.currentReplicas || row.totalReplicas || 0 }}/{{ row.autoscaling.desiredReplicas || row.totalReplicas || 0 }}
                </div>
                <div v-if="row.autoscaling.warnings?.length" class="autoscaling-warning-line">
                  治理提示 {{ row.autoscaling.warnings.length }} 条
                </div>
              </div>
              <div v-else-if="canConfigureAutoscaling(row)" class="autoscaling-empty">
                <el-tag type="info" effect="plain" size="small">未配置</el-tag>
              </div>
              <span v-else class="no-update">-</span>
            </template>
          </el-table-column>

          <el-table-column label="治理工作台" min-width="240">
            <template #default="{ row }">
              <div class="governance-cell">
                <div class="governance-head">
                  <el-tag :type="row.autoscaling?.enabled ? 'success' : 'info'" size="small">
                    {{ row.autoscaling?.enabled ? 'HPA 已接入' : '待接入 HPA' }}
                  </el-tag>
                  <el-tag v-if="row.autoscaling?.warnings?.length" type="warning" size="small">
                    {{ row.autoscaling.warnings.length }} 条提示
                  </el-tag>
                </div>
                <div class="governance-summary-line">
                  风险、阻断、容量建议、告警中心、AI 诊断、发布联动统一收口
                </div>
                <el-button type="primary" link @click="openGovernanceDialog(row)">进入治理工作台</el-button>
              </div>
            </template>
          </el-table-column>
           
          <el-table-column label="镜像" min-width="250">
            <template #default="{ row }">
              <div class="images-list">
                <div
                  v-for="(image, index) in row.images.slice(0, 1)"
                  :key="index"
                  class="image-tag-wrapper"
                  @click="copyToClipboard(image, '镜像地址已复制')"
                >
                  <el-icon class="copy-icon"><DocumentCopy /></el-icon>
                  <span class="full-image-name">{{ image }}</span>
                </div>
                <el-button
                  v-if="row.images.length > 1"
                  type="text"
                  size="small"
                  class="more-images-btn"
                  @click="viewAllImages(row)"
                >
                  +{{ row.images.length - 1 }}个镜像
                </el-button>
              </div>
            </template>
          </el-table-column>
          
          <el-table-column label="更新时间" min-width="150">
            <template #default="{ row }">
              <div class="time-info">
                <span v-if="row.updatedAt" class="datetime-text">
                  {{ formatDateTime(row.updatedAt) }}
                </span>
                <span v-else class="no-update">-</span>
              </div>
            </template>
          </el-table-column>
          
          <el-table-column label="操作" width="220" fixed="right">
            <template #default="{ row }">
              <div class="operation-buttons">
                <el-tooltip content="治理工作台" placement="top">
                  <el-button
                    type="warning"
                    :icon="Connection"
                    size="small"
                    circle
                    @click="openGovernanceDialog(row)"
                  />
                </el-tooltip>

                <el-tooltip :content="canScale(row) ? '伸缩' : '不支持伸缩'" placement="top">
                  <el-button
                    type="primary"
                    size="small"
                    circle
                    v-authority="['k8s:workload:expandable']"
                    :disabled="!canScale(row)"
                    @click="scaleWorkload(row)"
                  >
                    <img src="@/assets/image/扩容.svg" alt="伸缩" width="16" height="16" style="filter: brightness(0) invert(1);" />
                  </el-button>
                </el-tooltip>

                <el-tooltip :content="canConfigureAutoscaling(row) ? 'HPA 自动扩缩容' : '当前类型不支持 HPA'" placement="top">
                  <el-button
                    :type="row.autoscaling?.enabled ? 'success' : 'info'"
                    :icon="Monitor"
                    size="small"
                    circle
                    v-authority="['k8s:workload:expandable']"
                    :disabled="!canConfigureAutoscaling(row)"
                    @click="openAutoscalingDialog(row)"
                  />
                </el-tooltip>

                <el-tooltip content="容量建议" placement="top">
                  <el-button
                    type="success"
                    :icon="DataAnalysis"
                    size="small"
                    circle
                    v-authority="['k8s:workload:expandable']"
                    @click="openCapacitySuggestionDialog(row)"
                  />
                </el-tooltip>
                 
                <el-tooltip :content="canRestart(row) ? '重构' : '不支持重构'" placement="top">
                  <el-button
                    type="warning"
                    size="small"
                    circle
                    v-authority="['k8s:workload:restart']"
                    :disabled="!canRestart(row)"
                    @click="restartWorkload(row)"
                  >
                    <img src="@/assets/image/重启.svg" alt="重启" width="14" height="14" style="filter: brightness(0) invert(1);" />
                  </el-button>
                </el-tooltip>
                
                <el-tooltip :content="canUpdatePodConfig(row) ? '更新Pod配置' : '系统资源不可配置'" placement="top">
                  <el-button
                    type="success"
                    :icon="Setting"
                    size="small"
                    circle
                    v-authority="['k8s:workload:resource']"
                    :disabled="!canUpdatePodConfig(row)"
                    @click="updatePodConfig(row)"
                  />
                </el-tooltip>
                
                <el-tooltip :content="canUpdateScheduling(row) ? '更新调度' : '不支持调度更新'" placement="top">
                  <el-button
                    type="info"
                    :icon="Monitor"
                    size="small"
                    circle
                    v-authority="['k8s:workload:dispatch']"
                    :disabled="!canUpdateScheduling(row)"
                    @click="updateScheduling(row)"
                  />
                </el-tooltip>
                
                <el-tooltip :content="canEditYaml(row) ? '编辑YAML' : '系统资源不可编辑'" placement="top">
                  <el-button
                    type="primary"
                    :icon="Document"
                    size="small"
                    circle
                    v-authority="['k8s:workload:edityaml']"
                    :disabled="!canEditYaml(row)"
                    @click="editWorkloadYaml(row)"
                  />
                </el-tooltip>
                
                <el-tooltip :content="canDelete(row) ? '删除' : '系统资源不可删除'" placement="top">
                  <el-button
                    type="danger"
                    :icon="Delete"
                    size="small"
                    circle
                    v-authority="['k8s:workload:delete']"
                    :disabled="!canDelete(row)"
                    @click="deleteWorkload(row)"
                  />
                </el-tooltip>
              </div>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </el-card>


    <!-- Pod列表对话框 -->
    <PodListDialog
      :visible="podListDialogVisible"
      :workload="currentWorkload"
      @update:visible="podListDialogVisible = $event"
      @close="podListDialogVisible = false"
      @view-logs="viewPodLogs"
      @view-yaml="viewYaml"
      @rebuild-pod="rebuildPod"
      @view-events="viewPodEvents"
    />

    <!-- Pod事件对话框 -->
    <PodEventsDialog
      :visible="podEventsDialogVisible"
      :cluster-id="selectedClusterId"
      :namespace-name="queryParams.namespace"
      :pod-name="currentPodForEvents.name || ''"
      @update:visible="podEventsDialogVisible = $event"
      @close="podEventsDialogVisible = false"
    />

    <!-- Pod日志对话框 -->
    <el-dialog
      v-model="logDialogVisible"
      :title="`Pod日志 - ${currentPod.name || ''}`"
      width="1000px"
      class="log-dialog"
    >
      <div class="log-controls">
        <el-form :inline="true" size="small">
          <el-form-item label="容器">
            <el-select v-model="logParams.container" style="width: 200px">
              <el-option
                v-for="container in currentPod.containers || []"
                :key="container.name"
                :label="container.name"
                :value="container.name"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="行数">
            <el-input-number v-model="logParams.lines" :min="10" :max="1000" />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="viewPodLogs(currentPod)">刷新日志</el-button>
            <el-button @click="copyToClipboard(currentPodLogs, '日志已复制')">复制日志</el-button>
          </el-form-item>
        </el-form>
      </div>
      <div class="log-content">
        <pre>{{ currentPodLogs }}</pre>
      </div>
    </el-dialog>

    <!-- Pod YAML查看对话框 -->
    <PodYamlDialog
      :visible="podYamlDialogVisible"
      :yaml-content="currentYaml"
      :resource-name="currentPod.name"
      :resource-type="'Pod'"
      :editable="podYamlDialogEditable"
      @update:visible="podYamlDialogVisible = $event"
      @close="podYamlDialogVisible = false; podYamlDialogEditable = false"
      @save="handleYamlSave"
    />

    <!-- 扩缩容对话框 -->
    <el-dialog
      v-model="scaleDialogVisible"
      :title="`扩缩容 - ${currentWorkload.name || ''}`"
      width="400px"
      class="scale-dialog"
    >
      <el-alert
        v-if="currentWorkload.autoscaling?.enabled"
        type="warning"
        :closable="false"
        show-icon
        class="scale-warning-alert"
        title="当前工作负载仍受 HPA 管理，手动修改副本数后可能被自动扩缩容策略重新调整。"
      />
      <el-alert
        v-if="currentWorkload.autoscaling?.warnings?.length"
        type="warning"
        :closable="false"
        show-icon
        class="scale-warning-alert"
      >
        <template #title>
          <div class="autoscaling-warning-list">
            <div
              v-for="(warning, index) in currentWorkload.autoscaling.warnings"
              :key="`scale-warning-${index}`"
            >
              {{ warning }}
            </div>
          </div>
        </template>
      </el-alert>
      <el-form :model="scaleForm" label-width="80px">
        <el-form-item label="副本数" required>
          <el-input-number
            v-model="scaleForm.replicas"
            :min="0"
            :max="100"
            style="width: 100%"
          />
          <div class="form-tip">当前副本数: {{ currentWorkload.totalReplicas }}</div>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="scaleDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitScale">确定</el-button>
        </span>
      </template>
    </el-dialog>

    <el-dialog
      v-model="autoscalingDialogVisible"
      :title="`HPA 自动扩缩容 - ${currentWorkload.name || ''}`"
      width="720px"
      class="autoscaling-dialog"
      @close="resetAutoscalingForm(currentWorkload)"
    >
      <div v-loading="autoscalingLoading">
        <el-alert
          v-if="autoscalingNeedsRequestsWarning"
          type="warning"
          :closable="false"
          show-icon
          class="autoscaling-tip"
          title="当前工作负载未配置对应资源 requests，若使用利用率阈值可能无法生效。建议先补 requests，或改用平均值阈值。"
        />

        <el-alert
          v-if="autoscalingPersistedWarnings.length"
          type="warning"
          :closable="false"
          show-icon
          class="autoscaling-tip"
        >
          <template #title>
            <div class="autoscaling-warning-list">
              <div
                v-for="(warning, index) in autoscalingPersistedWarnings"
                :key="`autoscaling-warning-${index}`"
              >
                {{ warning }}
              </div>
            </div>
          </template>
        </el-alert>

        <el-descriptions :column="2" border class="autoscaling-summary">
          <el-descriptions-item label="工作负载">
            {{ getWorkloadTypeName(currentWorkload.type) }} / {{ currentWorkload.name || '-' }}
          </el-descriptions-item>
          <el-descriptions-item label="HPA 状态">
            <el-tag :type="autoscalingExists ? 'success' : 'info'" size="small">
              {{ autoscalingExists ? '已配置' : '未配置' }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="当前副本">
            {{ currentWorkload.readyReplicas || 0 }}/{{ currentWorkload.totalReplicas || 0 }}
          </el-descriptions-item>
          <el-descriptions-item label="最近伸缩">
            {{ currentWorkload.autoscaling?.lastScaleTime ? formatDateTime(currentWorkload.autoscaling.lastScaleTime) : '-' }}
          </el-descriptions-item>
        </el-descriptions>

        <el-form :model="autoscalingForm" label-width="120px" class="autoscaling-form">
          <el-row :gutter="16">
            <el-col :span="12">
              <el-form-item label="最小副本">
                <el-input-number v-model="autoscalingForm.minReplicas" :min="1" :max="100" style="width: 100%" />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="最大副本">
                <el-input-number v-model="autoscalingForm.maxReplicas" :min="1" :max="200" style="width: 100%" />
              </el-form-item>
            </el-col>
          </el-row>

          <div class="autoscaling-metric-card">
            <div class="autoscaling-metric-header">
              <div class="autoscaling-metric-title">CPU 指标</div>
              <el-switch v-model="autoscalingForm.cpuEnabled" active-text="启用" inactive-text="关闭" />
            </div>
            <el-row :gutter="16" v-if="autoscalingForm.cpuEnabled">
              <el-col :span="10">
                <el-form-item label="目标类型" label-width="90px">
                  <el-select v-model="autoscalingForm.cpuTargetType" style="width: 100%">
                    <el-option
                      v-for="item in autoscalingTargetTypeOptions"
                      :key="`cpu-${item.value}`"
                      :label="item.label"
                      :value="item.value"
                    />
                  </el-select>
                </el-form-item>
              </el-col>
              <el-col :span="14">
                <el-form-item :label="autoscalingForm.cpuTargetType === 'Utilization' ? '目标值(%)' : '目标值'" label-width="90px">
                  <el-input
                    v-model="autoscalingForm.cpuTargetValue"
                    :placeholder="autoscalingForm.cpuTargetType === 'Utilization' ? '例如 70' : '例如 100m'"
                  />
                </el-form-item>
              </el-col>
            </el-row>
          </div>

          <div class="autoscaling-metric-card">
            <div class="autoscaling-metric-header">
              <div class="autoscaling-metric-title">内存指标</div>
              <el-switch v-model="autoscalingForm.memoryEnabled" active-text="启用" inactive-text="关闭" />
            </div>
            <el-row :gutter="16" v-if="autoscalingForm.memoryEnabled">
              <el-col :span="10">
                <el-form-item label="目标类型" label-width="90px">
                  <el-select v-model="autoscalingForm.memoryTargetType" style="width: 100%">
                    <el-option
                      v-for="item in autoscalingTargetTypeOptions"
                      :key="`memory-${item.value}`"
                      :label="item.label"
                      :value="item.value"
                    />
                  </el-select>
                </el-form-item>
              </el-col>
              <el-col :span="14">
                <el-form-item :label="autoscalingForm.memoryTargetType === 'Utilization' ? '目标值(%)' : '目标值'" label-width="90px">
                  <el-input
                    v-model="autoscalingForm.memoryTargetValue"
                    :placeholder="autoscalingForm.memoryTargetType === 'Utilization' ? '例如 80' : '例如 256Mi'"
                  />
                </el-form-item>
              </el-col>
            </el-row>
          </div>

          <div class="autoscaling-metric-card">
            <div class="autoscaling-metric-title">治理窗口</div>
            <el-row :gutter="16">
              <el-col :span="12">
                <el-form-item label="扩容稳定窗口">
                  <el-input-number v-model="autoscalingForm.scaleUpWindow" :min="0" :max="3600" style="width: 100%" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="缩容稳定窗口">
                  <el-input-number v-model="autoscalingForm.scaleDownWindow" :min="0" :max="3600" style="width: 100%" />
                </el-form-item>
              </el-col>
            </el-row>

            <el-divider content-position="left">扩容策略</el-divider>
            <el-row :gutter="16">
              <el-col :span="12">
                <el-form-item label="扩容选择策略">
                  <el-select v-model="autoscalingForm.scaleUpSelectPolicy" style="width: 100%">
                    <el-option
                      v-for="item in autoscalingSelectPolicyOptions"
                      :key="`up-${item.value}`"
                      :label="item.label"
                      :value="item.value"
                    />
                  </el-select>
                </el-form-item>
              </el-col>
            </el-row>
            <el-row :gutter="16">
              <el-col :span="12">
                <el-form-item label="Pods 限速">
                  <div class="autoscaling-policy-line">
                    <el-switch v-model="autoscalingForm.scaleUpPodsEnabled" />
                    <el-input-number
                      v-model="autoscalingForm.scaleUpPodsValue"
                      :min="1"
                      :max="100"
                      :disabled="!autoscalingForm.scaleUpPodsEnabled"
                      style="width: 140px"
                    />
                    <span class="autoscaling-policy-tip">每 15 秒最多增加 Pods</span>
                  </div>
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="Percent 限速">
                  <div class="autoscaling-policy-line">
                    <el-switch v-model="autoscalingForm.scaleUpPercentEnabled" />
                    <el-input-number
                      v-model="autoscalingForm.scaleUpPercentValue"
                      :min="1"
                      :max="500"
                      :disabled="!autoscalingForm.scaleUpPercentEnabled"
                      style="width: 140px"
                    />
                    <span class="autoscaling-policy-tip">每 15 秒最多增加百分比</span>
                  </div>
                </el-form-item>
              </el-col>
            </el-row>

            <el-divider content-position="left">缩容策略</el-divider>
            <el-row :gutter="16">
              <el-col :span="12">
                <el-form-item label="缩容选择策略">
                  <el-select v-model="autoscalingForm.scaleDownSelectPolicy" style="width: 100%">
                    <el-option
                      v-for="item in autoscalingSelectPolicyOptions"
                      :key="`down-${item.value}`"
                      :label="item.label"
                      :value="item.value"
                    />
                  </el-select>
                </el-form-item>
              </el-col>
            </el-row>
            <el-row :gutter="16">
              <el-col :span="12">
                <el-form-item label="Pods 限速">
                  <div class="autoscaling-policy-line">
                    <el-switch v-model="autoscalingForm.scaleDownPodsEnabled" />
                    <el-input-number
                      v-model="autoscalingForm.scaleDownPodsValue"
                      :min="1"
                      :max="100"
                      :disabled="!autoscalingForm.scaleDownPodsEnabled"
                      style="width: 140px"
                    />
                    <span class="autoscaling-policy-tip">每 15 秒最多减少 Pods</span>
                  </div>
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="Percent 限速">
                  <div class="autoscaling-policy-line">
                    <el-switch v-model="autoscalingForm.scaleDownPercentEnabled" />
                    <el-input-number
                      v-model="autoscalingForm.scaleDownPercentValue"
                      :min="1"
                      :max="500"
                      :disabled="!autoscalingForm.scaleDownPercentEnabled"
                      style="width: 140px"
                    />
                    <span class="autoscaling-policy-tip">每 15 秒最多减少百分比</span>
                  </div>
                </el-form-item>
              </el-col>
            </el-row>
          </div>
        </el-form>

        <div v-if="currentWorkload.autoscaling?.metrics?.length" class="autoscaling-runtime-panel">
          <div class="autoscaling-runtime-title">当前运行态</div>
          <div class="autoscaling-runtime-tags">
            <el-tag
              v-for="metric in currentWorkload.autoscaling.metrics"
              :key="`runtime-${metric.resourceName}`"
              type="success"
              effect="light"
            >
              {{ formatAutoscalingTarget(metric) }} | 当前 {{ formatAutoscalingCurrent(metric) }}
            </el-tag>
          </div>
        </div>
      </div>

      <template #footer>
        <span class="dialog-footer">
          <el-button v-if="autoscalingExists" type="danger" plain :loading="autoscalingSubmitting" @click="deleteAutoscaling">
            删除 HPA
          </el-button>
          <el-button @click="autoscalingDialogVisible = false">取消</el-button>
          <el-button type="primary" :loading="autoscalingSubmitting" @click="submitAutoscaling">
            {{ autoscalingExists ? '更新 HPA' : '创建 HPA' }}
          </el-button>
        </span>
      </template>
    </el-dialog>

    <el-dialog
      v-model="governanceDialogVisible"
      :title="`治理工作台 - ${currentWorkload.name || ''}`"
      width="960px"
      class="autoscaling-dialog"
    >
      <div v-loading="governanceLoading">
        <el-descriptions :column="2" border class="autoscaling-summary">
          <el-descriptions-item label="工作负载">
            {{ getWorkloadTypeName(currentWorkload.type) }} / {{ currentWorkload.name || '-' }}
          </el-descriptions-item>
          <el-descriptions-item label="风险等级">
            <el-tag :type="capacityRiskTagType(governanceOverview.riskLevel)" size="small">
              {{ capacityRiskText(governanceOverview.riskLevel) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="阻断状态">
            <el-tag :type="governanceOverview.blocking ? 'danger' : 'success'" size="small">
              {{ governanceOverview.blocking ? '存在阻断项' : '当前可放行' }}
            </el-tag>
            <span v-if="governanceOverview.blockingReason" class="governance-blocking-reason">
              {{ governanceOverview.blockingReason }}
            </span>
          </el-descriptions-item>
          <el-descriptions-item label="关联告警">
            未恢复 {{ governanceOverview.alertSummary?.openEventCount || 0 }} / Incident {{ governanceOverview.alertSummary?.incidentCount || 0 }}
          </el-descriptions-item>
          <el-descriptions-item label="当前策略">
            {{ governanceOverview.currentSuggestion?.recommendedPolicy?.type || '-' }}
          </el-descriptions-item>
          <el-descriptions-item label="历史快照">
            <span v-if="governanceOverview.latestSnapshot?.historyId">#{{ governanceOverview.latestSnapshot.historyId }} / {{ governanceOverview.latestSnapshot.generatedAt || '-' }}</span>
            <span v-else>暂无</span>
          </el-descriptions-item>
        </el-descriptions>

        <div class="governance-toolbar">
          <el-button type="primary" @click="openGovernanceAutoscaling">HPA 策略</el-button>
          <el-button type="success" @click="openGovernanceCapacity">容量建议</el-button>
          <el-button @click="openGovernanceAlertCenter">告警中心</el-button>
          <el-button @click="openGovernanceAIDiagnosis">AI 诊断</el-button>
        </div>

        <div class="governance-grid">
          <div class="autoscaling-metric-card">
            <div class="autoscaling-metric-title">治理阶段</div>
            <div class="governance-stage-list">
              <div
                v-for="item in governanceOverview.stages || []"
                :key="item.key"
                class="governance-stage-item"
              >
                <div class="governance-stage-head">
                  <span class="governance-stage-title">{{ item.title }}</span>
                  <el-tag :type="governanceStageTagType(item.status)" size="small">{{ item.status || 'unknown' }}</el-tag>
                </div>
                <div class="governance-stage-summary">{{ item.summary || '-' }}</div>
                <div v-if="item.items?.length" class="governance-stage-items">
                  <span v-for="(subItem, index) in item.items" :key="`${item.key}-${index}`">{{ subItem }}</span>
                </div>
                <el-button v-if="item.primaryActionPath" type="primary" link @click="router.push(item.primaryActionPath)">
                  {{ item.primaryActionText || '查看详情' }}
                </el-button>
              </div>
            </div>
          </div>

          <div class="autoscaling-metric-card">
            <div class="autoscaling-metric-title">阻断 / 提示规则</div>
            <el-empty
              v-if="!(governanceOverview.blockingRules?.length || governanceOverview.warningRules?.length)"
              description="当前没有阻断或提示规则"
              :image-size="56"
            />
            <div v-else class="governance-rule-list">
              <div
                v-for="rule in governanceOverview.blockingRules || []"
                :key="`blocking-${rule.code}`"
                class="governance-rule-item blocking"
              >
                <div class="governance-rule-head">
                  <el-tag type="danger" size="small">BLOCKING</el-tag>
                  <span>{{ rule.title }}</span>
                </div>
                <div class="governance-rule-text">{{ rule.message }}</div>
                <div v-if="rule.action" class="governance-rule-action">{{ rule.action }}</div>
              </div>
              <div
                v-for="rule in governanceOverview.warningRules || []"
                :key="`warning-${rule.code}`"
                class="governance-rule-item"
              >
                <div class="governance-rule-head">
                  <el-tag type="warning" size="small">WARNING</el-tag>
                  <span>{{ rule.title }}</span>
                </div>
                <div class="governance-rule-text">{{ rule.message }}</div>
                <div v-if="rule.action" class="governance-rule-action">{{ rule.action }}</div>
              </div>
            </div>
          </div>

          <div class="autoscaling-metric-card">
            <div class="autoscaling-metric-title">持续观察指标</div>
            <el-empty v-if="!governanceOverview.watchMetrics?.length" description="暂无观察指标" :image-size="56" />
            <div v-else class="autoscaling-runtime-tags">
              <el-tag v-for="item in governanceOverview.watchMetrics" :key="item" size="small" class="autoscaling-metric-tag">
                {{ item }}
              </el-tag>
            </div>
          </div>

          <div class="autoscaling-metric-card">
            <div class="autoscaling-metric-title">关联应用 / 发布入口</div>
            <el-empty v-if="!governanceOverview.relatedApplications?.length" description="当前没有关联应用映射" :image-size="56" />
            <div v-else class="governance-related-apps">
              <div
                v-for="item in governanceOverview.relatedApplications"
                :key="`${item.appId}-${item.environment}`"
                class="governance-app-item"
              >
                <div>
                  <div class="governance-app-name">{{ item.appName || '-' }} / {{ item.environment || '-' }}</div>
                  <div class="governance-app-meta">
                    {{ item.appCode || '-' }} · {{ item.preCheckMode || 'observe' }} · {{ item.releaseEnabled ? '治理启用' : '仅映射' }}
                  </div>
                </div>
                <el-button type="primary" link @click="openGovernanceAppPath(item)">应用治理入口</el-button>
              </div>
            </div>
          </div>

          <div class="autoscaling-metric-card">
            <div class="autoscaling-metric-title">扩展位</div>
            <div class="governance-extension-list">
              <div
                v-for="item in governanceOverview.extensions || []"
                :key="item.key"
                class="governance-extension-item"
              >
                <div class="governance-stage-head">
                  <span class="governance-stage-title">{{ item.title }}</span>
                  <el-tag type="info" size="small">{{ item.status || 'reserved' }}</el-tag>
                </div>
                <div class="governance-stage-summary">{{ item.summary }}</div>
                <div class="governance-stage-items">
                  <span v-for="field in item.reservedFields || []" :key="`${item.key}-${field}`">{{ field }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="governanceDialogVisible = false">关闭</el-button>
        </span>
      </template>
    </el-dialog>

    <el-dialog
      v-model="capacitySuggestionDialogVisible"
      :title="`容量建议 - ${currentWorkload.name || ''}`"
      width="860px"
      class="autoscaling-dialog"
    >
      <div v-loading="capacitySuggestionLoading">
        <el-descriptions :column="2" border class="autoscaling-summary">
          <el-descriptions-item label="工作负载">
            {{ getWorkloadTypeName(currentWorkload.type) }} / {{ currentWorkload.name || '-' }}
          </el-descriptions-item>
          <el-descriptions-item label="风险等级">
            <el-tag :type="capacityRiskTagType(capacitySuggestion.riskLevel)" size="small">
              {{ capacityRiskText(capacitySuggestion.riskLevel) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="当前快照">
            <div class="capacity-snapshot-line">
              <el-tag :type="capacitySuggestionViewingLatest ? 'success' : 'warning'" size="small">
                {{ capacitySuggestionViewingLatest ? '最新生成' : `历史快照 #${capacitySuggestion.historyId || '-'}` }}
              </el-tag>
              <el-button v-if="!capacitySuggestionViewingLatest" type="primary" link @click="resetToLatestCapacitySuggestion">
                查看最新
              </el-button>
            </div>
          </el-descriptions-item>
          <el-descriptions-item label="生成信息">
            {{ capacitySuggestion.generatedAt || '-' }} / {{ capacitySuggestion.generatedBy || '-' }}
          </el-descriptions-item>
          <el-descriptions-item label="关联告警">
            未恢复 {{ capacitySuggestion.alertSummary?.openEventCount || 0 }} / 工单 {{ capacitySuggestion.alertSummary?.incidentCount || 0 }}
          </el-descriptions-item>
          <el-descriptions-item label="告警中心">
            <el-button type="primary" link @click="openAlertCenter">打开告警中心</el-button>
          </el-descriptions-item>
          <el-descriptions-item label="AI 诊断">
            <el-button type="primary" link @click="openAIDiagnosis">打开 AI 诊断</el-button>
          </el-descriptions-item>
        </el-descriptions>

        <div class="autoscaling-metric-card">
          <div class="autoscaling-metric-title">建议动作</div>
          <el-empty v-if="!capacitySuggestion.recommendedActions?.length" description="暂无建议动作" :image-size="56" />
          <div v-else class="capacity-actions">
            <div
              v-for="(item, index) in capacitySuggestion.recommendedActions"
              :key="`capacity-action-${index}`"
              class="capacity-action-item"
            >
              <div class="capacity-action-head">
                <el-tag :type="item.priority === 'P1' ? 'danger' : item.priority === 'P2' ? 'warning' : 'info'" size="small">
                  {{ item.priority }}
                </el-tag>
                <span class="capacity-action-title">{{ item.action }}</span>
              </div>
              <div class="capacity-action-desc">{{ item.reason }}</div>
              <div class="capacity-action-effect">预期收益：{{ item.expectedEffect }}</div>
            </div>
          </div>
        </div>

        <div class="autoscaling-metric-card">
          <div class="autoscaling-metric-title">建议策略</div>
          <el-descriptions :column="2" border>
            <el-descriptions-item label="策略类型">{{ capacitySuggestion.recommendedPolicy?.type || '-' }}</el-descriptions-item>
            <el-descriptions-item label="副本范围">{{ capacitySuggestion.recommendedPolicy?.minReplicas || '-' }} - {{ capacitySuggestion.recommendedPolicy?.maxReplicas || '-' }}</el-descriptions-item>
            <el-descriptions-item label="核心指标">{{ capacitySuggestion.recommendedPolicy?.metric || '-' }}</el-descriptions-item>
            <el-descriptions-item label="阈值">{{ capacitySuggestion.recommendedPolicy?.targetUtilization || '-' }}</el-descriptions-item>
          </el-descriptions>
        </div>

        <div class="autoscaling-metric-card">
          <div class="autoscaling-metric-title">建议报告</div>
          <pre class="capacity-report-block">{{ capacitySuggestion.report || '-' }}</pre>
        </div>

        <div class="autoscaling-metric-card">
          <div class="autoscaling-metric-title">Prompt 渲染结果</div>
          <pre class="capacity-report-block">{{ capacitySuggestion.renderedPrompt || '-' }}</pre>
        </div>

        <div class="autoscaling-metric-card">
          <div class="autoscaling-metric-title">历史快照</div>
          <el-table
            v-loading="capacitySuggestionHistoryLoading"
            :data="capacitySuggestionHistory"
            size="small"
            stripe
            class="capacity-history-table"
            empty-text="暂无历史快照"
          >
            <el-table-column label="生成时间" prop="generatedAt" min-width="160" />
            <el-table-column label="生成人" prop="generatedBy" width="120" />
            <el-table-column label="风险" width="110">
              <template #default="{ row }">
                <el-tag :type="capacityRiskTagType(row.riskLevel)" size="small">
                  {{ capacityRiskText(row.riskLevel) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="关联告警" min-width="140">
              <template #default="{ row }">
                未恢复 {{ row.alertSummary?.openEventCount || 0 }} / 工单 {{ row.alertSummary?.incidentCount || 0 }}
              </template>
            </el-table-column>
            <el-table-column label="操作" width="220" fixed="right">
              <template #default="{ row }">
                <div class="capacity-history-actions">
                  <el-button type="primary" link @click="viewCapacitySuggestionHistory(row)">查看</el-button>
                  <el-button type="primary" link @click="openAlertCenter(row)">告警中心</el-button>
                  <el-button type="primary" link @click="openAIDiagnosis(row)">AI 诊断</el-button>
                </div>
              </template>
            </el-table-column>
          </el-table>
          <div class="capacity-history-pagination">
            <el-pagination
              small
              :current-page="capacitySuggestionHistoryPagination.pageNum"
              :page-size="capacitySuggestionHistoryPagination.pageSize"
              :page-sizes="[5, 10, 20]"
              :total="capacitySuggestionHistoryPagination.total"
              layout="total, sizes, prev, pager, next"
              @size-change="handleCapacitySuggestionHistorySizeChange"
              @current-change="handleCapacitySuggestionHistoryPageChange"
            />
          </div>
        </div>
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="capacitySuggestionDialogVisible = false">关闭</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 工作负载标签查看对话框 -->
    <el-dialog
      v-model="workloadLabelsDialogVisible"
      :title="`工作负载标签 - ${currentWorkload.name || ''}`"
      width="700px"
      class="workload-labels-view-dialog"
    >
      <div class="labels-view-content" v-if="currentWorkload.labels">
        <!-- 用户自定义标签 -->
        <div class="labels-section" v-if="Object.keys(getUserLabels(currentWorkload.labels)).length > 0">
          <h4>用户标签</h4>
          <div class="labels-list">
            <el-tag
              v-for="(value, key) in getUserLabels(currentWorkload.labels)"
              :key="key"
              type="primary"
              size="default"
              class="label-tag"
              @click="copyToClipboard(`${key}=${value}`, '标签信息已复制')"
            >
              <el-icon class="tag-icon"><DocumentCopy /></el-icon>
              {{ key }}={{ value }}
            </el-tag>
          </div>
        </div>
        
        <!-- 系统标签 -->
        <div class="labels-section" v-if="Object.keys(getSystemLabels(currentWorkload.labels)).length > 0">
          <h4>系统标签</h4>
          <div class="labels-list">
            <el-tag
              v-for="(value, key) in getSystemLabels(currentWorkload.labels)"
              :key="key"
              type="info"
              size="default"
              class="label-tag system-label"
              @click="copyToClipboard(`${key}=${value}`, '标签信息已复制')"
            >
              <el-icon class="tag-icon"><DocumentCopy /></el-icon>
              {{ key }}={{ value }}
            </el-tag>
          </div>
        </div>
        
        <!-- 没有标签的提示 -->
        <div v-if="!currentWorkload.labels || Object.keys(currentWorkload.labels).length === 0" class="no-labels">
          <el-empty description="该工作负载没有标签" :image-size="60" />
        </div>
      </div>
    </el-dialog>

    <!-- 所有镜像查看对话框 -->
    <el-dialog
      v-model="allImagesDialogVisible"
      :title="`镜像列表 - ${currentWorkload.name || ''}`"
      width="800px"
      class="all-images-dialog"
    >
      <div class="images-view-content" v-if="currentWorkload.images">
        <div class="images-section">
          <h4>容器镜像 ({{ currentWorkload.images?.length || 0 }}个)</h4>
          <div class="all-images-list">
            <el-card
              v-for="(image, index) in currentWorkload.images"
              :key="index"
              class="image-card"
              shadow="hover"
            >
              <div class="image-info">
                <div class="image-name">
                  <el-icon class="image-icon"><Connection /></el-icon>
                  <span class="full-image-name">{{ image }}</span>
                </div>
                <div class="image-actions">
                  <el-button
                    type="primary"
                    size="small"
                    :icon="DocumentCopy"
                    @click="copyToClipboard(image, '镜像地址已复制')"
                  >
                    复制
                  </el-button>
                </div>
              </div>
              <div class="image-details">
                <el-tag size="small" type="info">{{ getImageTag(image) }}</el-tag>
                <el-tag size="small" type="success">{{ getImageRegistry(image) }}</el-tag>
              </div>
            </el-card>
          </div>
        </div>
        
        <!-- 没有镜像的提示 -->
        <div v-if="!currentWorkload.images || currentWorkload.images.length === 0" class="no-images">
          <el-empty description="该工作负载没有镜像" :image-size="60" />
        </div>
      </div>
    </el-dialog>

    <!-- Pod配置更新对话框 -->
    <PodConfigDialog
      :visible="podConfigDialogVisible"
      :workload="currentWorkload"
      @update:visible="podConfigDialogVisible = $event"
      @close="podConfigDialogVisible = false"
      @submit="handlePodConfigSubmit"
    />

    <!-- 调度更新对话框 -->
    <el-dialog
      v-model="schedulingDialogVisible"
      :title="`更新调度 - ${currentWorkload.name || ''}`"
      width="500px"
      class="scheduling-dialog"
    >
      <el-form :model="schedulingForm" label-width="120px">
        <el-form-item label="节点选择器">
          <el-input
            v-model="schedulingForm.nodeSelectorText"
            placeholder="如: kubernetes.io/arch=amd64"
            style="width: 100%"
          />
          <div class="form-tip">格式: key=value，多个用逗号分隔</div>
        </el-form-item>
        <el-form-item label="节点亲和性">
          <el-select
            v-model="schedulingForm.nodeAffinity"
            placeholder="请选择节点亲和性"
            style="width: 100%"
          >
            <el-option label="无要求" value="none" />
            <el-option label="偏好调度" value="preferred" />
            <el-option label="必须调度" value="required" />
          </el-select>
        </el-form-item>
        <el-form-item label="Pod反亲和性">
          <el-switch
            v-model="schedulingForm.podAntiAffinity"
            active-text="启用"
            inactive-text="禁用"
          />
          <div class="form-tip">避免Pod调度到同一节点</div>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="schedulingDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitScheduling">更新调度</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 工作负载YAML编辑对话框 -->
    <PodYamlDialog
      :visible="workloadYamlDialogVisible"
      :yaml-content="currentYaml"
      :resource-name="currentWorkload.name"
      :resource-type="currentWorkload.type || 'Workload'"
      :editable="true"
      @update:visible="workloadYamlDialogVisible = $event"
      @close="workloadYamlDialogVisible = false"
      @save="handleYamlSave"
    />

    <!-- 创建工作负载对话框 -->
    <CreatePodDialog
      ref="createPodDialogRef"
      :visible="createPodDialogVisible"
      :cluster-id="selectedClusterId"
      :cluster-name="clusterList.find(c => c.id === selectedClusterId)?.name"
      :namespace="queryParams.namespace"
      @update:visible="createPodDialogVisible = $event"
      @close="createPodDialogVisible = false"
      @preview="handlePodPreview"
      @create="handlePodCreate"
    />
  </div>
</template>


<style scoped>
.k8s-workloads-management {
  padding: 20px;
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.workloads-card {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 16px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.title {
  font-size: 20px;
  font-weight: 600;
  color: #2c3e50;
  background: linear-gradient(45deg, #667eea, #764ba2);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.header-actions {
  display: flex;
  gap: 12px;
  align-items: center;
}

.search-section {
  margin-bottom: 24px;
  padding: 20px;
  background: rgba(103, 126, 234, 0.05);
  border-radius: 12px;
  border: 1px solid rgba(103, 126, 234, 0.1);
}

.search-form .el-form-item {
  margin-bottom: 0;
  margin-right: 16px;
}

/* 工作负载类型标签页样式 */
.workload-type-section {
  margin-top: 20px;
  padding-top: 20px;
  border-top: 1px solid rgba(103, 126, 234, 0.1);
}

.workload-tabs {
  margin: 0;
}

.workload-tabs :deep(.el-tabs__header) {
  margin-bottom: 0;
}

.workload-tabs :deep(.el-tabs__item) {
  font-weight: 500;
  color: #606266;
}

.workload-tabs :deep(.el-tabs__item.is-active) {
  color: #409EFF;
  font-weight: 600;
}

.search-form .el-form-item__label {
  color: #606266;
  font-weight: 500;
}

.table-section {
  margin-top: 20px;
}

.workloads-table {
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
}

.workloads-table :deep(.el-table__header table),
.workloads-table :deep(.el-table__body table) {
  width: 100% !important;
  table-layout: fixed !important;
}

.workloads-table :deep(col:nth-child(2)) {
  width: 84px !important;
}

.workloads-table :deep(col:nth-child(3)) {
  width: 108px !important;
}

.workloads-table :deep(col:nth-child(4)) {
  width: 136px !important;
}

.workloads-table :deep(col:nth-child(5)) {
  width: 160px !important;
}

.workloads-table :deep(col:nth-child(7)) {
  width: 180px !important;
}

.workloads-table :deep(col:nth-child(9)) {
  width: 140px !important;
}

.workloads-table :deep(col:nth-child(6)),
.workloads-table :deep(col:nth-child(8)),
.workloads-table :deep(.el-table__header-wrapper th:nth-child(6)),
.workloads-table :deep(.el-table__header-wrapper th:nth-child(8)),
.workloads-table :deep(.el-table__body-wrapper td:nth-child(6)),
.workloads-table :deep(.el-table__body-wrapper td:nth-child(8)) {
  display: none !important;
}

.workloads-table :deep(.el-table__fixed-right) {
  width: 220px !important;
}

.workloads-table :deep(.el-table__fixed-right colgroup col),
.workloads-table :deep(.el-table__fixed-right .el-table__cell) {
  width: 220px !important;
}

.workloads-table :deep(.el-table__header) {
  background: linear-gradient(135deg, #667eea, #764ba2);
}

.workloads-table :deep(.el-table__header th) {
  background: transparent !important;
  color: #2c3e50 !important;
  font-weight: 700 !important;
  border-bottom: none;
  text-shadow: 0 1px 2px rgba(255, 255, 255, 0.8);
}

.workloads-table :deep(.el-table__row) {
  transition: all 0.3s ease;
}

.workloads-table :deep(.el-table__row:hover) {
  background-color: rgba(103, 126, 234, 0.05) !important;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.workload-name-container {
  display: flex;
  align-items: center;
  gap: 10px;
  cursor: default;
}

.workload-name-container:hover {
  transform: none !important;
  background-color: transparent !important;
}

.workload-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.workload-name {
  font-weight: 600;
  color: #2c3e50;
  font-size: 14px;
}

.clickable-name {
  color: #409EFF !important;
  cursor: pointer;
  transition: all 0.2s ease;
  text-decoration: underline;
  text-decoration-color: transparent;
}

.clickable-name:hover {
  color: #337ECC !important;
  text-decoration-color: #409EFF;
  text-shadow: 0 1px 2px rgba(64, 158, 255, 0.2);
}

.workload-type-label {
  font-size: 12px;
  color: #E6A23C;
  font-weight: 500;
  pointer-events: none;
  user-select: none;
}

.workload-type-tag {
  font-size: 11px;
  height: 18px;
  line-height: 16px;
  padding: 0 6px;
}

.pod-name-container {
  display: flex;
  align-items: center;
  gap: 10px;
}

.k8s-icon {
  width: 20px;
  height: 20px;
  flex-shrink: 0;
  filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.1));
}

.workload-name-link {
  font-weight: 600;
  color: #667eea;
  text-decoration: none;
  transition: all 0.3s ease;
}

.workload-name-link:hover {
  color: #764ba2;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.pod-name {
  font-weight: 500;
  color: #2c3e50;
}

.resource-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.resource-row {
  display: flex;
  align-items: center;
  gap: 6px;
}

.resource-type {
  font-size: 12px;
  color: #909399;
  min-width: 35px;
  font-weight: 500;
}

.resource-values {
  display: flex;
  align-items: center;
  gap: 2px;
}

.request-value {
  font-size: 12px;
  color: #67c23a;
  font-weight: 500;
}

.separator {
  font-size: 12px;
  color: #dcdfe6;
  margin: 0 2px;
}

.limit-value {
  font-size: 12px;
  color: #e6a23c;
  font-weight: 500;
}

.resource-item {
  display: flex;
  align-items: center;
  gap: 4px;
}

.resource-label {
  font-size: 12px;
  color: #909399;
  min-width: 55px;
}

.resource-value {
  font-size: 12px;
  color: #606266;
  font-weight: 500;
}

.ip-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.ip-item {
  display: flex;
  align-items: center;
  gap: 4px;
}

.ip-label {
  font-size: 12px;
  color: #909399;
  min-width: 35px;
}

.ip-value {
  font-size: 12px;
  color: #606266;
  font-weight: 500;
}

.images-list {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
}

.image-tag {
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 4px;
  max-width: 150px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.image-tag:hover {
  transform: translateY(-2px);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}

.copy-icon {
  font-size: 10px;
}

.more-images {
  cursor: default;
}

/* Pod状态样式 */
.pod-status-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
}

.pod-count-tag {
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 6px 12px;
  font-weight: 600;
}

.pod-count-tag:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
}

.pod-icon {
  font-size: 14px;
}

.pod-status-text {
  font-size: 11px;
  line-height: 1.2;
}

.status-ready {
  color: #67c23a;
  font-weight: 500;
}

.status-partial {
  color: #e6a23c;
  font-weight: 500;
}

.status-starting {
  color: #f56c6c;
  font-weight: 500;
}

.status-stopped {
  color: #909399;
  font-weight: 500;
}

/* 标签容器样式 */
.label-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100%;
  padding: 4px 0;
}

.label-badge {
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

.label-icon-button {
  background: transparent;
  border: none;
  color: #606266;
  transition: all 0.3s ease;
}

.label-icon-button:hover {
  background: transparent;
  color: #409eff;
  transform: scale(1.1);
}

/* 镜像显示优化 */
.image-tag-wrapper {
  display: flex;
  align-items: center;
  gap: 6px;
  cursor: pointer;
  transition: all 0.3s ease;
  margin-bottom: 4px;
}

.image-tag-wrapper:hover {
  transform: translateY(-1px);
}

.image-tag-wrapper .full-image-name {
  font-family: 'Monaco', 'Courier New', monospace;
  font-size: 11px;
  color: #2c3e50;
  word-break: break-all;
  line-height: 1.4;
  white-space: normal;
}

.image-tag-wrapper .copy-icon {
  color: #666;
  font-size: 12px;
  flex-shrink: 0;
}

.more-images-btn {
  color: #409eff;
  font-size: 12px;
  padding: 2px 6px;
  margin-left: 4px;
}

.more-images-btn:hover {
  color: #66b1ff;
  background-color: rgba(64, 158, 255, 0.1);
}

/* 时间显示 */
.time-info {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2px;
}

.datetime-text {
  font-family: 'Monaco', 'Courier New', monospace;
  font-size: 12px;
  color: #2c3e50;
}

.no-update {
  color: #909399;
  font-size: 12px;
}

.running-time {
  font-size: 10px;
  color: #909399;
  line-height: 1.2;
}

.no-update {
  color: #c0c4cc;
  font-size: 12px;
}


.operation-buttons {
  display: flex;
  gap: 4px;
  justify-content: flex-start;
  flex-wrap: wrap;
}

.operation-buttons .el-button {
  width: 28px;
  height: 28px;
  padding: 0;
  transition: all 0.3s ease;
}

.operation-buttons .el-button:hover:not(.is-disabled) {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
}

.operation-buttons .el-button.is-disabled {
  cursor: not-allowed;
  opacity: 0.5;
  background-color: #f5f7fa !important;
  border-color: #e4e7ed !important;
  color: #c0c4cc !important;
}

.operation-buttons .el-button.is-disabled:hover {
  transform: none;
  box-shadow: none;
}

/* 对话框样式 */
.pod-list-dialog :deep(.el-dialog),
.log-dialog :deep(.el-dialog),
.yaml-dialog :deep(.el-dialog),
.scale-dialog :deep(.el-dialog),
.autoscaling-dialog :deep(.el-dialog),
.workload-labels-view-dialog :deep(.el-dialog),
.all-images-dialog :deep(.el-dialog) {
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
}

.pod-list-dialog :deep(.el-dialog__header),
.log-dialog :deep(.el-dialog__header),
.yaml-dialog :deep(.el-dialog__header),
.scale-dialog :deep(.el-dialog__header),
.autoscaling-dialog :deep(.el-dialog__header),
.workload-labels-view-dialog :deep(.el-dialog__header),
.all-images-dialog :deep(.el-dialog__header) {
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: white;
  border-top-left-radius: 16px;
  border-top-right-radius: 16px;
  padding: 20px 24px;
}

.pod-list-dialog :deep(.el-dialog__title),
.log-dialog :deep(.el-dialog__title),
.yaml-dialog :deep(.el-dialog__title),
.scale-dialog :deep(.el-dialog__title),
.autoscaling-dialog :deep(.el-dialog__title),
.workload-labels-view-dialog :deep(.el-dialog__title),
.all-images-dialog :deep(.el-dialog__title) {
  color: white;
  font-weight: 600;
}

.log-controls,
.yaml-controls {
  margin-bottom: 16px;
  padding: 16px;
  background: rgba(103, 126, 234, 0.05);
  border-radius: 8px;
}

.log-content,
.yaml-content {
  background: #2c3e50;
  color: #ecf0f1;
  padding: 16px;
  border-radius: 8px;
  max-height: 400px;
  overflow: auto;
}

.log-content pre,
.yaml-content pre {
  margin: 0;
  font-family: 'Courier New', monospace;
  font-size: 12px;
  line-height: 1.4;
  white-space: pre-wrap;
  word-wrap: break-word;
}

.form-tip {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
  line-height: 1.4;
}

.dialog-footer {
  display: flex;
  gap: 12px;
  justify-content: flex-end;
}

/* 集群选择样式 */
.cluster-option {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
}

.cluster-name {
  font-weight: 500;
  color: #2c3e50;
}

.cluster-status-tag {
  margin-left: 8px;
}

/* 命名空间选择样式 */
.namespace-option {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
}

.namespace-name {
  font-weight: 500;
  color: #2c3e50;
}

.namespace-status-tag {
  margin-left: 8px;
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

.el-loading-mask {
  background-color: rgba(103, 126, 234, 0.1);
  backdrop-filter: blur(4px);
}

.scale-warning-alert,
.autoscaling-tip {
  margin-bottom: 16px;
}

.autoscaling-cell {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.autoscaling-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}

.autoscaling-range,
.autoscaling-status-line {
  font-size: 12px;
  color: #5c6470;
}

.autoscaling-warning-line {
  font-size: 12px;
  color: #c27b11;
}

.autoscaling-empty {
  display: flex;
  align-items: center;
  min-height: 56px;
}

.autoscaling-metrics,
.autoscaling-runtime-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.autoscaling-metric-tag {
  max-width: 100%;
}

.autoscaling-summary {
  margin-bottom: 16px;
}

.autoscaling-form {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.autoscaling-metric-card {
  border: 1px solid rgba(103, 126, 234, 0.16);
  border-radius: 12px;
  padding: 16px;
  background: rgba(248, 250, 255, 0.9);
}

.autoscaling-metric-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 12px;
}

.autoscaling-metric-title,
.autoscaling-runtime-title {
  font-size: 14px;
  font-weight: 600;
  color: #2c3e50;
}

.autoscaling-runtime-panel {
  margin-top: 8px;
  padding: 16px;
  border-radius: 12px;
  background: rgba(245, 247, 255, 0.95);
  border: 1px dashed rgba(103, 126, 234, 0.28);
}

.autoscaling-warning-list {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.autoscaling-policy-line {
  display: flex;
  align-items: center;
  gap: 10px;
  width: 100%;
}

.autoscaling-policy-tip {
  font-size: 12px;
  color: #7a8291;
  line-height: 1.4;
}

.capacity-actions {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.capacity-action-item {
  border: 1px solid rgba(103, 126, 234, 0.14);
  border-radius: 10px;
  padding: 12px 14px;
  background: rgba(255, 255, 255, 0.92);
}

.capacity-action-head {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 6px;
}

.capacity-action-title {
  font-weight: 600;
  color: #243044;
}

.capacity-action-desc,
.capacity-action-effect {
  font-size: 13px;
  color: #5f6b7b;
  line-height: 1.5;
}

.capacity-action-effect {
  margin-top: 4px;
}

.capacity-snapshot-line,
.capacity-history-actions {
  display: flex;
  align-items: center;
  gap: 10px;
}

.capacity-history-table {
  margin-top: 12px;
}

.capacity-history-pagination {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}

.capacity-report-block {
  margin: 0;
  padding: 14px;
  border-radius: 12px;
  background: #0f172a;
  color: #e2e8f0;
  white-space: pre-wrap;
  word-break: break-word;
  font-family: Consolas, Monaco, monospace;
  max-height: 340px;
  overflow: auto;
}

.governance-cell {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.governance-head,
.governance-stage-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}

.governance-summary-line,
.governance-stage-summary,
.governance-rule-text,
.governance-rule-action,
.governance-app-meta,
.governance-blocking-reason {
  font-size: 12px;
  color: #5f6b7b;
  line-height: 1.5;
}

.governance-toolbar {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
  margin-bottom: 16px;
}

.governance-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 16px;
}

.governance-stage-list,
.governance-rule-list,
.governance-extension-list,
.governance-related-apps {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.governance-stage-item,
.governance-rule-item,
.governance-extension-item,
.governance-app-item {
  border: 1px solid rgba(103, 126, 234, 0.14);
  border-radius: 12px;
  padding: 12px 14px;
  background: rgba(255, 255, 255, 0.94);
}

.governance-rule-item.blocking {
  border-color: rgba(220, 38, 38, 0.2);
  background: rgba(254, 242, 242, 0.9);
}

.governance-stage-title,
.governance-app-name {
  font-weight: 600;
  color: #243044;
}

.governance-stage-items {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin: 8px 0 2px;
}

.governance-stage-items span {
  padding: 4px 8px;
  border-radius: 999px;
  background: rgba(15, 23, 42, 0.05);
  color: #475569;
  font-size: 12px;
}

.governance-rule-head {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 6px;
  font-weight: 600;
  color: #243044;
}

.governance-app-item {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  align-items: center;
}

/* 响应式设计 */
@media (max-width: 1200px) {
  .operation-buttons {
    gap: 4px;
  }

  .governance-grid {
    grid-template-columns: 1fr;
  }
  
  .operation-buttons .el-button {
    margin: 1px;
  }
  
  .header-actions .el-select {
    min-width: 180px;
  }
}

@media (max-width: 768px) {
  .k8s-workloads-management {
    padding: 10px;
  }
  
  .search-form {
    flex-direction: column;
  }
  
  .search-form .el-form-item {
    margin-right: 0;
    margin-bottom: 12px;
  }
  
  .operation-buttons {
    flex-direction: column;
    gap: 4px;
  }

  .governance-app-item {
    flex-direction: column;
    align-items: flex-start;
  }
  
  .workloads-table :deep(.el-table__row:hover) {
    transform: none;
  }
}

/* 标签和镜像对话框样式 */
.labels-view-content,
.images-view-content {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.no-labels,
.no-images {
  text-align: center;
  padding: 40px 20px;
  color: #909399;
}

/* 镜像对话框特殊样式 */
.all-images-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.image-card {
  border: 1px solid rgba(103, 126, 234, 0.2);
  border-radius: 8px;
  transition: all 0.3s ease;
}

.image-card:hover {
  border-color: #667eea;
  box-shadow: 0 4px 12px rgba(103, 126, 234, 0.15);
}

.image-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
}

.image-name {
  display: flex;
  align-items: center;
  gap: 8px;
  flex: 1;
  min-width: 0;
}

.image-icon {
  color: #667eea;
  font-size: 16px;
  flex-shrink: 0;
}

.full-image-name {
  font-family: 'Monaco', 'Courier New', monospace;
  font-size: 13px;
  color: #2c3e50;
  word-break: break-all;
  line-height: 1.4;
}

.image-actions {
  flex-shrink: 0;
}

.image-details {
  display: flex;
  gap: 8px;
  margin-top: 12px;
  padding-top: 12px;
  border-top: 1px solid #f0f0f0;
}

</style>

