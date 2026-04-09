<template>
  <div class="namespaces-monitoring">
    <el-row :gutter="16" style="margin-bottom: 12px;">
      <el-col :span="8" v-for="namespace in namespaceData" :key="namespace.name">
        <div class="namespace-card" v-loading="namespacesLoading">
          <!-- 卡片头部 -->
          <div class="card-header">
            <div class="namespace-info-header">
              <div class="namespace-icon">
                <img src="@/assets/image/k8s.svg" alt="k8s" class="k8s-icon" />
                <div class="pulse-ring"></div>
              </div>
              <div class="namespace-details">
                <h4>{{ namespace.name }}</h4>
                <el-tag 
                  :type="getNamespaceStatusTag(namespace.status)" 
                  size="small"
                  :class="['status-tag', namespace.status.toLowerCase()]"
                >
                  <i class="status-dot"></i>
                  {{ getNamespaceStatusText(namespace.status) }}
                </el-tag>
              </div>
            </div>
            <div class="namespace-name">{{ namespace.name }}</div>
          </div>
          
          <div v-if="hasNamespaceMetrics(namespaceMetrics[namespace.name])" class="namespace-metrics">
            <!-- CPU 使用率 -->
            <div class="metric-item">
              <div class="metric-header">
                <img src="@/assets/image/cpu.svg" alt="CPU" class="metric-icon" />
                <span class="metric-label">CPU使用率</span>
                <span class="metric-value">{{ getNamespaceCpuUsageRate(namespaceMetrics[namespace.name]) }}%</span>
              </div>
              <div class="metric-progress">
                <div class="progress-bg">
                  <div 
                    class="progress-fill cpu-progress" 
                    :style="{ width: getNamespaceCpuUsageRate(namespaceMetrics[namespace.name]) + '%' }"
                  >
                    <div class="progress-glow"></div>
                  </div>
                </div>
              </div>
              <div class="metric-detail">
                {{ formatCpuValue(namespaceMetrics[namespace.name].totalUsage?.cpu) }} / {{ formatQuotaValue(namespaceMetrics[namespace.name].resourceQuota?.hard?.cpu, formatCpuValue) }}
              </div>
            </div>
            
            <!-- 内存使用率 -->
            <div class="metric-item">
              <div class="metric-header">
                <img src="@/assets/image/内存.svg" alt="Memory" class="metric-icon" />
                <span class="metric-label">内存使用率</span>
                <span class="metric-value">{{ getNamespaceMemoryUsageRate(namespaceMetrics[namespace.name]) }}%</span>
              </div>
              <div class="metric-progress">
                <div class="progress-bg">
                  <div 
                    class="progress-fill memory-progress" 
                    :style="{ width: getNamespaceMemoryUsageRate(namespaceMetrics[namespace.name]) + '%' }"
                  >
                    <div class="progress-glow"></div>
                  </div>
                </div>
              </div>
              <div class="metric-detail">
                {{ formatMemoryValue(namespaceMetrics[namespace.name].totalUsage?.memory) }} / {{ formatQuotaValue(namespaceMetrics[namespace.name].resourceQuota?.hard?.memory, formatMemoryValue) }}
              </div>
            </div>
            
            <!-- 系统信息 -->
            <div class="system-stats">
              <div class="stat-item">
                <div class="stat-icon pods-icon">
                  <i class="icon-dot"></i>
                </div>
                <div class="stat-content">
                  <span class="stat-label">Pod数量</span>
                  <span class="stat-value">{{ namespaceMetrics[namespace.name].podCount || 0 }}</span>
                </div>
              </div>
              <div class="stat-item">
                <div class="stat-icon services-icon">
                  <i class="icon-pulse"></i>
                </div>
                <div class="stat-content">
                  <span class="stat-label">服务数量</span>
                  <span class="stat-value">{{ namespaceMetrics[namespace.name].serviceCount || 0 }}</span>
                </div>
              </div>
            </div>
          </div>

          <div v-else class="no-metrics">
            <div class="loading-animation">
              <div class="loading-dot"></div>
              <div class="loading-dot"></div>
              <div class="loading-dot"></div>
            </div>
            <p class="no-data-text">{{ getNamespaceMetricsMessage(namespaceMetrics[namespace.name]) }}</p>
          </div>

          <!-- 操作按钮 -->
          <div class="namespace-actions">
            <button 
              class="action-btn primary-btn"
              @click="viewNamespaceDetail(namespace)"
            >
              <svg class="btn-icon" viewBox="0 0 24 24" width="16" height="16">
                <path d="M12,9A3,3 0 0,0 9,12A3,3 0 0,0 12,15A3,3 0 0,0 15,12A3,3 0 0,0 12,9M12,17A5,5 0 0,1 7,12A5,5 0 0,1 12,7A5,5 0 0,1 17,12A5,5 0 0,1 12,17M12,4.5C7,4.5 2.73,7.61 1,12C2.73,16.39 7,19.5 12,19.5C17,19.5 21.27,16.39 23,12C21.27,7.61 17,4.5 12,4.5Z" fill="currentColor"/>
              </svg>
              <span>详细信息</span>
              <div class="btn-glow"></div>
            </button>
            <button 
              class="action-btn secondary-btn"
              @click="refreshNamespaceMetrics(namespace.name)"
            >
              <svg class="btn-icon refresh-icon" viewBox="0 0 24 24" width="16" height="16">
                <path d="M17.65,6.35C16.2,4.9 14.21,4 12,4A8,8 0 0,0 4,12A8,8 0 0,0 12,20C15.73,20 18.84,17.45 19.73,14H17.65C16.83,16.33 14.61,18 12,18A6,6 0 0,1 6,12A6,6 0 0,1 12,6C13.66,6 15.14,6.69 16.22,7.78L13,11H20V4L17.65,6.35Z" fill="currentColor"/>
              </svg>
              <span>刷新</span>
            </button>
          </div>
        </div>
      </el-col>
    </el-row>

      <!-- 命名空间详情对话框 -->
      <el-dialog
        v-model="namespaceDetailDialogVisible"
        :title="`命名空间详情 - ${selectedNamespace?.name || ''}`"
        width="70%"
        class="namespace-detail-dialog"
        :before-close="closeNamespaceDetailDialog"
      >
        <div class="namespace-detail-content" v-if="selectedNamespace">
          <!-- 基本信息区域 -->
          <div class="detail-section">
            <h4 class="section-title">
              <svg class="section-icon" viewBox="0 0 24 24" width="18" height="18">
                <path d="M13,9H11V7H13M13,17H11V11H13M12,2A10,10 0 0,0 2,12A10,10 0 0,0 12,22A10,10 0 0,0 22,12A10,10 0 0,0 12,2Z" fill="currentColor"/>
              </svg>
              基本信息
            </h4>
            <div class="info-grid">
              <div class="info-item">
                <span class="info-label">命名空间名称</span>
                <span class="info-value">{{ selectedNamespace.name }}</span>
              </div>
              <div class="info-item">
                <span class="info-label">状态</span>
                <el-tag :type="getNamespaceStatusTag(selectedNamespace.status)" size="small">
                  {{ getNamespaceStatusText(selectedNamespace.status) }}
                </el-tag>
              </div>
              <div class="info-item">
                <span class="info-label">创建时间</span>
                <span class="info-value">{{ formatTimestamp(selectedNamespace.createdAt) }}</span>
              </div>
              <div class="info-item">
                <span class="info-label">标签数量</span>
                <span class="info-value">{{ Object.keys(selectedNamespace.labels || {}).length }}</span>
              </div>
            </div>
          </div>

          <!-- 资源使用情况 -->
          <div class="detail-section">
            <h4 class="section-title">
              <svg class="section-icon" viewBox="0 0 24 24" width="18" height="18">
                <path d="M19,3H5C3.9,3 3,3.9 3,5V19C3,20.1 3.9,21 5,21H19C20.1,21 21,20.1 21,19V5C21,3.9 20.1,3 19,3M19,19H5V5H19V19Z" fill="currentColor"/>
              </svg>
              资源使用情况
            </h4>
            <div class="resource-stats">
              <div class="resource-stat-item">
                <div class="stat-header">
                  <svg class="stat-icon cpu-icon" viewBox="0 0 24 24" width="20" height="20">
                    <path d="M17,17H7V7H17M21,11V9H19V7C19,5.89 18.1,5 17,5H15V3H13V5H11V3H9V5H7C5.89,5 5,5.89 5,7V9H3V11H5V13H3V15H5V17C5,18.1 5.89,19 7,19H9V21H11V19H13V21H15V19H17C18.1,19 19,18.1 19,17V15H21V13H19V11M16,8H8V16H16V8Z" fill="currentColor"/>
                  </svg>
                  <span class="stat-label">CPU使用率</span>
                </div>
                <div class="stat-value">{{ getNamespaceCpuUsageRate(selectedNamespace.metrics) }}%</div>
                <div class="custom-progress large">
                  <div 
                    class="progress-bar cpu-progress" 
                    :style="{ width: getNamespaceCpuUsageRate(selectedNamespace.metrics) + '%' }"
                  ></div>
                </div>
                <div class="stat-details">
                  <span>使用: {{ formatCpuValue(selectedNamespace.metrics?.totalUsage?.cpu) }}</span>
                  <span>限制: {{ formatQuotaValue(selectedNamespace.metrics?.resourceQuota?.hard?.cpu, formatCpuValue) }}</span>
                </div>
              </div>

              <div class="resource-stat-item">
                <div class="stat-header">
                  <svg class="stat-icon memory-icon" viewBox="0 0 24 24" width="20" height="20">
                    <path d="M17,7H22V9H19V17H17V19H15V17H13V19H11V17H9V19H7V17H5V9H2V7H7V5H9V7H11V5H13V7H15V5H17V7M17,9V15H7V9H17Z" fill="currentColor"/>
                  </svg>
                  <span class="stat-label">内存使用率</span>
                </div>
                <div class="stat-value">{{ getNamespaceMemoryUsageRate(selectedNamespace.metrics) }}%</div>
                <div class="custom-progress large">
                  <div 
                    class="progress-bar memory-progress" 
                    :style="{ width: getNamespaceMemoryUsageRate(selectedNamespace.metrics) + '%' }"
                  ></div>
                </div>
                <div class="stat-details">
                  <span>使用: {{ formatMemoryValue(selectedNamespace.metrics?.totalUsage?.memory) }}</span>
                  <span>限制: {{ formatQuotaValue(selectedNamespace.metrics?.resourceQuota?.hard?.memory, formatMemoryValue) }}</span>
                </div>
              </div>
            </div>
          </div>

          <!-- Pod列表 -->
          <div class="detail-section">
            <h4 class="section-title">
              <svg class="section-icon" viewBox="0 0 24 24" width="18" height="18">
                <path d="M19,5V7H15V5H19M9,5V11H5V5H9M19,13V19H15V13H19M9,17V19H5V17H9M21,3H13V9H21V3M11,3H3V13H11V3M21,11H13V21H21V11M11,15H3V21H11V15Z" fill="currentColor"/>
              </svg>
              Pod列表 ({{ namespacePods.length }})
            </h4>
            <div class="pod-list">
              <div class="pod-item" v-for="pod in namespacePods" :key="pod.name">
                <span class="pod-name">{{ pod.name }}</span>
                <el-tag :type="getPodStatusTag(pod.status)" size="small" class="pod-status">
                  {{ pod.status }}
                </el-tag>
                <span class="pod-resource">CPU: {{ formatCpuValue(pod.cpu) }}</span>
                <span class="pod-resource">Memory: {{ formatMemoryValue(pod.memory) }}</span>
              </div>
            </div>
          </div>
        </div>
      </el-dialog>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import k8sApi from '@/api/k8s'

const props = defineProps({
  selectedClusterId: {
    type: [String, Number],
    required: true
  }
})

const namespaceDetailDialogVisible = ref(false)
const selectedNamespace = ref(null)
const namespacePods = ref([])
const namespacesLoading = ref(false)
const namespaceData = ref([])
const namespaceMetrics = ref({})

const getEmptyMetrics = (namespace = {}) => ({
  namespace: namespace.name || '',
  timestamp: '',
  podCount: namespace.resourceCount?.podCount || 0,
  runningPods: 0,
  dataAvailable: false,
  collectionStatus: 'unavailable',
  unavailableReason: '暂无可用监控样本',
  totalUsage: {
    cpu: '0',
    memory: '0'
  },
  resourceQuota: {
    hard: {
      cpu: '',
      memory: ''
    },
    used: {
      cpu: '',
      memory: ''
    }
  },
  usageRate: {
    cpuRate: 0,
    memoryRate: 0
  },
  podMetrics: [],
  serviceCount: namespace.resourceCount?.serviceCount || 0
})

const hasNamespaceMetrics = (metrics) => metrics?.dataAvailable === true

const getNamespaceMetricsMessage = (metrics) => {
  if (!metrics) {
    return '正在获取监控数据...'
  }
  return metrics.unavailableReason || '暂无可用监控样本'
}

const hasNamespaceUsageSamples = (metrics, namespace = {}) => {
  const totalUsage = metrics?.totalUsage || {}
  const usageValues = [totalUsage.cpu, totalUsage.memory]
  const hasUsageValue = usageValues.some((value) => {
    if (value === undefined || value === null || value === '') {
      return false
    }
    const normalized = String(value).trim()
    return normalized !== '' && normalized !== '0' && normalized !== '0m' && normalized !== '0Mi'
  })

  return Boolean(
    metrics?.timestamp ||
    hasUsageValue ||
    (Array.isArray(metrics?.podMetrics) && metrics.podMetrics.length > 0) ||
    namespace.resourceCount?.podCount > 0 ||
    namespace.resourceCount?.serviceCount > 0
  )
}

const normalizeNamespace = (namespace) => ({
  name: namespace.name,
  status: namespace.status || 'Active',
  createdAt: namespace.createdAt || '',
  labels: namespace.labels || {},
  resourceCount: {
    podCount: namespace.resourceCount?.podCount || 0,
    serviceCount: namespace.resourceCount?.serviceCount || 0
  }
})

const parseNamespaces = (responseData) => {
  const namespaces = responseData?.data?.namespaces || responseData?.data || []
  return Array.isArray(namespaces) ? namespaces : []
}

const parsePods = (responseData) => {
  const pods = responseData?.data?.pods || responseData?.data || []
  return Array.isArray(pods) ? pods : []
}

const formatCpuValue = (value) => {
  if (!value) return '0m'

  if (typeof value === 'string') {
    if (value.includes('n')) {
      const nanoValue = parseInt(value.replace('n', ''), 10)
      const milliValue = Math.round(nanoValue / 1000000)
      return `${milliValue}m`
    }
    if (value.includes('m')) {
      return value
    }
    return `${parseInt(value, 10) * 1000}m`
  }

  return `${Math.round(value / 1000000)}m`
}

const formatMemoryValue = (value) => {
  if (!value) return '0Mi'

  if (typeof value === 'string') {
    if (value.includes('Ki')) {
      const kiValue = parseInt(value.replace('Ki', ''), 10)
      const miValue = Math.round(kiValue / 1024)
      return `${miValue}Mi`
    }
    if (value.includes('Mi')) {
      return value
    }
    if (value.includes('Gi')) {
      const giValue = parseFloat(value.replace('Gi', ''))
      return `${Math.round(giValue * 1024)}Mi`
    }
    if (value.includes('Ti')) {
      const tiValue = parseFloat(value.replace('Ti', ''))
      return `${Math.round(tiValue * 1024 * 1024)}Mi`
    }
  }

  const bytes = parseInt(value, 10)
  if (Number.isNaN(bytes)) {
    return String(value)
  }
  const miValue = Math.round(bytes / (1024 * 1024))
  return `${miValue}Mi`
}

const formatQuotaValue = (value, formatter) => {
  if (!value) return '-'
  return formatter(value)
}

const getNamespaceCpuUsageRate = (metrics) => {
  const rate = Number(metrics?.usageRate?.cpuRate)
  if (!Number.isFinite(rate)) return 0
  return Math.min(Math.max(Math.round(rate), 0), 100)
}

const getNamespaceMemoryUsageRate = (metrics) => {
  const rate = Number(metrics?.usageRate?.memoryRate)
  if (!Number.isFinite(rate)) return 0
  return Math.min(Math.max(Math.round(rate), 0), 100)
}

const formatTimestamp = (timestamp) => {
  if (!timestamp) return '-'

  try {
    const date = new Date(timestamp)
    return date.toLocaleString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit'
    })
  } catch (error) {
    return timestamp
  }
}

const fetchNamespaceList = async () => {
  const response = await k8sApi.getNamespaces(props.selectedClusterId)
  const responseData = response.data || response

  if (responseData.code !== 200 && !responseData.success) {
    throw new Error(responseData.message || '获取命名空间列表失败')
  }

  namespaceData.value = parseNamespaces(responseData).map(normalizeNamespace)
  return namespaceData.value
}

const fetchNamespaceMetrics = async (namespaceName, options = {}) => {
  const { silent = false } = options
  const namespace = namespaceData.value.find(item => item.name === namespaceName) || { name: namespaceName }

  try {
    const response = await k8sApi.getNamespaceMetrics(props.selectedClusterId, namespaceName)
    const responseData = response.data || response

    if (responseData.code !== 200 && !responseData.success) {
      throw new Error(responseData.message || `获取命名空间 ${namespaceName} 监控数据失败`)
    }

    const metrics = responseData.data || {}
    const hasSamples = hasNamespaceUsageSamples(metrics, namespace)
    const nextMetrics = {
      ...getEmptyMetrics(namespace),
      ...metrics,
      dataAvailable: metrics.dataAvailable ?? hasSamples,
      collectionStatus: metrics.collectionStatus || (hasSamples ? 'ready' : 'unavailable'),
      unavailableReason: metrics.unavailableReason || (hasSamples ? '' : '暂无可用监控样本'),
      podMetrics: Array.isArray(metrics.podMetrics) ? metrics.podMetrics : [],
      serviceCount: namespace.resourceCount?.serviceCount || 0
    }

    namespaceMetrics.value = {
      ...namespaceMetrics.value,
      [namespaceName]: nextMetrics
    }

    if (selectedNamespace.value?.name === namespaceName) {
      selectedNamespace.value = {
        ...selectedNamespace.value,
        metrics: nextMetrics
      }
    }
  } catch (error) {
    namespaceMetrics.value = {
      ...namespaceMetrics.value,
      [namespaceName]: getEmptyMetrics(namespace)
    }

    if (!silent) {
      throw error
    }
  }
}

const buildNamespacePods = (pods, podMetrics = []) => {
  const metricsByName = new Map(
    podMetrics.map(metric => [metric.podName, metric])
  )

  return pods.map(pod => {
    const metric = metricsByName.get(pod.name)

    return {
      name: pod.name,
      status: pod.status || pod.phase || 'Unknown',
      cpu: metric?.usage?.cpu || '0',
      memory: metric?.usage?.memory || '0'
    }
  })
}

const fetchNamespacePods = async (namespaceName) => {
  const response = await k8sApi.getPodList(props.selectedClusterId, namespaceName)
  const responseData = response.data || response

  if (responseData.code !== 200 && !responseData.success) {
    throw new Error(responseData.message || `获取命名空间 ${namespaceName} Pod 列表失败`)
  }

  namespacePods.value = buildNamespacePods(
    parsePods(responseData),
    namespaceMetrics.value[namespaceName]?.podMetrics || []
  )
}

const refreshAllData = async () => {
  if (!props.selectedClusterId) {
    namespaceData.value = []
    namespaceMetrics.value = {}
    namespacePods.value = []
    selectedNamespace.value = null
    namespaceDetailDialogVisible.value = false
    return
  }

  try {
    namespacesLoading.value = true
    const namespaces = await fetchNamespaceList()
    await Promise.all(
      namespaces.map(namespace => fetchNamespaceMetrics(namespace.name, { silent: true }))
    )
  } catch (error) {
    console.error('获取命名空间监控数据失败:', error)
    namespaceData.value = []
    namespaceMetrics.value = {}
    ElMessage.error(error.message || '获取命名空间监控数据失败')
  } finally {
    namespacesLoading.value = false
  }
}

const refreshNamespaceMetrics = async (namespaceName) => {
  if (!props.selectedClusterId) return

  try {
    namespacesLoading.value = true
    await fetchNamespaceMetrics(namespaceName)
    if (selectedNamespace.value?.name === namespaceName) {
      await fetchNamespacePods(namespaceName)
    }
    ElMessage.success(`命名空间 ${namespaceName} 监控数据已刷新`)
  } catch (error) {
    console.error(`刷新命名空间 ${namespaceName} 监控数据失败:`, error)
    ElMessage.error(error.message || '刷新监控数据失败')
  } finally {
    namespacesLoading.value = false
  }
}

const getNamespaceStatusTag = (status) => {
  switch (status) {
    case 'Active':
      return 'success'
    case 'Terminating':
      return 'warning'
    case 'Inactive':
      return 'danger'
    default:
      return 'info'
  }
}

const getNamespaceStatusText = (status) => {
  switch (status) {
    case 'Active':
      return '活跃'
    case 'Terminating':
      return '终止中'
    case 'Inactive':
      return '非活跃'
    default:
      return '未知'
  }
}

const getPodStatusTag = (status) => {
  switch (status) {
    case 'Running':
      return 'success'
    case 'Pending':
      return 'warning'
    case 'Failed':
      return 'danger'
    case 'Succeeded':
      return 'info'
    default:
      return 'info'
  }
}

const viewNamespaceDetail = async (namespace) => {
  selectedNamespace.value = {
    ...namespace,
    metrics: namespaceMetrics.value[namespace.name] || getEmptyMetrics(namespace)
  }
  namespaceDetailDialogVisible.value = true

  try {
    namespacesLoading.value = true
    await fetchNamespacePods(namespace.name)
  } catch (error) {
    console.error(`获取命名空间 ${namespace.name} Pod 详情失败:`, error)
    namespacePods.value = []
    ElMessage.error(error.message || '获取命名空间 Pod 详情失败')
  } finally {
    namespacesLoading.value = false
  }
}

const closeNamespaceDetailDialog = () => {
  namespaceDetailDialogVisible.value = false
  selectedNamespace.value = null
  namespacePods.value = []
}

const handleClusterChange = async () => {
  await refreshAllData()
}

defineExpose({
  refreshAllData,
  handleClusterChange
})

watch(
  () => props.selectedClusterId,
  async () => {
    await refreshAllData()
  },
  { immediate: true }
)
</script>

<style scoped>
.namespaces-monitoring {
  padding: 0;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  min-height: 100vh;
  padding: 12px;
  border-radius: 16px;
}

/* 命名空间卡片 */
.namespace-card {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-radius: 16px;
  padding: 18px;
  box-shadow: 0 6px 24px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
  overflow: hidden;
}

.namespace-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 4px;
  background: linear-gradient(90deg, #4facfe 0%, #00f2fe 100%);
  border-radius: 20px 20px 0 0;
}

.namespace-card:hover {
  transform: translateY(-8px) scale(1.02);
  box-shadow: 0 16px 64px rgba(0, 0, 0, 0.15);
  border-color: rgba(79, 172, 254, 0.4);
}

/* 卡片头部 */
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 16px;
}

.namespace-info-header {
  display: flex;
  align-items: center;
  gap: 16px;
}

.namespace-icon {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 48px;
  height: 48px;
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
  border-radius: 12px;
  box-shadow: 0 4px 16px rgba(79, 172, 254, 0.4);
}

.k8s-icon {
  width: 24px;
  height: 24px;
  filter: brightness(0) invert(1);
}

.pulse-ring {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 60px;
  height: 60px;
  border: 2px solid rgba(79, 172, 254, 0.6);
  border-radius: 50%;
  animation: pulse 2s infinite;
}

.namespace-details h4 {
  margin: 0 0 6px 0;
  font-size: 16px;
  font-weight: 600;
  color: #1a1a1a;
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
}

.status-tag {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 4px 12px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 500;
}

.status-tag.active {
  background: rgba(46, 204, 113, 0.1);
  color: #2ecc71;
  border: 1px solid rgba(46, 204, 113, 0.2);
}

.status-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: currentColor;
  animation: statusPulse 2s infinite;
}

.namespace-name {
  color: #606266;
  font-size: 12px;
  font-weight: 500;
  text-align: right;
  opacity: 0.8;
}

/* 监控指标 */
.namespace-metrics {
  margin-bottom: 16px;
}

.metric-item {
  margin-bottom: 12px;
  background: rgba(255, 255, 255, 0.6);
  border-radius: 10px;
  padding: 12px;
  border: 1px solid rgba(255, 255, 255, 0.3);
  transition: all 0.3s ease;
}

.metric-item:hover {
  background: rgba(255, 255, 255, 0.8);
  transform: translateX(4px);
}

.metric-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 8px;
}

.metric-header .metric-icon {
  width: 16px;
  height: 16px;
  filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.1));
  animation: iconFloat 3s ease-in-out infinite;
}

.metric-label {
  font-size: 12px;
  font-weight: 500;
  color: #333;
  margin-left: 6px;
  flex: 1;
}

.metric-value {
  font-size: 14px;
  font-weight: 700;
  background: linear-gradient(45deg, #4facfe, #00f2fe);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.metric-progress {
  margin-bottom: 6px;
}

.progress-bg {
  width: 100%;
  height: 6px;
  background: rgba(0, 0, 0, 0.1);
  border-radius: 3px;
  overflow: hidden;
  position: relative;
}

.progress-fill {
  height: 100%;
  border-radius: 3px;
  position: relative;
  transition: width 1s cubic-bezier(0.4, 0, 0.2, 1);
  overflow: hidden;
}

.cpu-progress {
  background: linear-gradient(90deg, #4facfe 0%, #00f2fe 100%);
}

.memory-progress {
  background: linear-gradient(90deg, #43e97b 0%, #38f9d7 100%);
}

.progress-glow {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(90deg, transparent 0%, rgba(255, 255, 255, 0.4) 50%, transparent 100%);
  animation: progressGlow 2s infinite;
}

.metric-detail {
  font-size: 12px;
  color: #666;
  font-family: 'Courier New', monospace;
}

/* 系统统计 */
.system-stats {
  display: flex;
  gap: 16px;
  margin-bottom: 20px;
}

.stat-item {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 12px;
  background: rgba(255, 255, 255, 0.6);
  padding: 12px;
  border-radius: 12px;
  border: 1px solid rgba(255, 255, 255, 0.3);
  transition: all 0.3s ease;
}

.stat-item:hover {
  background: rgba(255, 255, 255, 0.8);
  transform: scale(1.02);
}

.stat-icon {
  width: 32px;
  height: 32px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.pods-icon {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.services-icon {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
}

.icon-dot, .icon-pulse {
  width: 12px;
  height: 12px;
  background: white;
  border-radius: 50%;
  position: relative;
}

.icon-pulse {
  animation: iconPulse 2s infinite;
}

.stat-content {
  flex: 1;
}

.stat-label {
  display: block;
  font-size: 11px;
  color: #666;
  font-weight: 500;
  margin-bottom: 2px;
}

.stat-value {
  display: block;
  font-size: 16px;
  font-weight: 700;
  background: linear-gradient(45deg, #4facfe, #00f2fe);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

/* 无数据状态 */
.no-metrics {
  text-align: center;
  padding: 40px 20px;
  background: rgba(255, 255, 255, 0.3);
  border-radius: 12px;
  margin-bottom: 20px;
}

.loading-animation {
  display: flex;
  justify-content: center;
  gap: 4px;
  margin-bottom: 16px;
}

.loading-dot {
  width: 8px;
  height: 8px;
  background: #4facfe;
  border-radius: 50%;
  animation: loadingDot 1.4s infinite ease-in-out both;
}

.loading-dot:nth-child(1) { animation-delay: -0.32s; }
.loading-dot:nth-child(2) { animation-delay: -0.16s; }

.no-data-text {
  color: #666;
  font-size: 14px;
  margin: 0;
}

/* 操作按钮 */
.namespace-actions {
  display: flex;
  gap: 12px;
}

.action-btn {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 12px 16px;
  border: none;
  border-radius: 10px;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
  overflow: hidden;
}

.primary-btn {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
  color: white;
  box-shadow: 0 4px 16px rgba(79, 172, 254, 0.4);
}

.primary-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(79, 172, 254, 0.5);
}

.secondary-btn {
  background: rgba(255, 255, 255, 0.8);
  color: #4facfe;
  border: 1px solid rgba(79, 172, 254, 0.3);
}

.secondary-btn:hover {
  background: rgba(79, 172, 254, 0.1);
  border-color: rgba(79, 172, 254, 0.5);
  transform: translateY(-1px);
}

.btn-icon {
  width: 16px;
  height: 16px;
  color: inherit;
}

.refresh-icon {
  animation: none;
  transition: transform 0.3s ease;
}

.secondary-btn:hover .refresh-icon {
  animation: refreshSpin 0.8s ease-in-out;
}

.btn-glow {
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.4), transparent);
  transition: left 0.6s;
}

.primary-btn:hover .btn-glow {
  left: 100%;
}

/* 对话框样式 */
.namespace-detail-dialog :deep(.el-dialog) {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-radius: 16px;
  border: 1px solid rgba(255, 255, 255, 0.3);
}

.namespace-detail-dialog :deep(.el-dialog__header) {
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.1), rgba(118, 75, 162, 0.1));
  border-radius: 16px 16px 0 0;
  border-bottom: 1px solid rgba(255, 255, 255, 0.2);
  padding: 16px 20px;
}

.namespace-detail-dialog :deep(.el-dialog__title) {
  color: #333;
  font-weight: 600;
  font-size: 16px;
}

.namespace-detail-dialog :deep(.el-dialog__body) {
  padding: 16px 20px;
}

.namespace-detail-content {
  max-height: 60vh;
  overflow-y: auto;
}

.detail-section {
  margin-bottom: 16px;
  background: rgba(255, 255, 255, 0.8);
  border-radius: 10px;
  padding: 14px;
  border: 1px solid rgba(255, 255, 255, 0.3);
}

.section-title {
  display: flex;
  align-items: center;
  color: #333;
  font-size: 14px;
  font-weight: 600;
  margin-bottom: 12px;
  padding-bottom: 6px;
  border-bottom: 1px solid rgba(102, 126, 234, 0.2);
}

.section-icon {
  color: #667eea;
  margin-right: 6px;
  width: 16px;
  height: 16px;
}

/* 基本信息网格 */
.info-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 8px;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 3px;
}

.info-label {
  color: #666;
  font-size: 11px;
  font-weight: 500;
}

.info-value {
  color: #333;
  font-weight: 600;
  font-size: 12px;
}

/* 资源统计 */
.resource-stats {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 8px;
}

.resource-stat-item {
  background: rgba(255, 255, 255, 0.6);
  border-radius: 10px;
  padding: 14px;
  border: 1px solid rgba(255, 255, 255, 0.5);
}

.stat-header {
  display: flex;
  align-items: center;
  margin-bottom: 8px;
}

.stat-icon {
  margin-right: 6px;
  width: 16px;
  height: 16px;
}

.stat-label {
  color: #333;
  font-weight: 500;
  font-size: 12px;
}

.stat-value {
  color: #333;
  font-size: 20px;
  font-weight: 700;
  margin-bottom: 8px;
}

.stat-details {
  display: flex;
  justify-content: space-between;
  margin-top: 6px;
  font-size: 11px;
  color: #666;
}

/* Pod列表 */
.pod-list {
  max-height: 300px;
  overflow-y: auto;
}

.pod-item {
  display: flex;
  align-items: center;
  gap: 12px;
  background: rgba(255, 255, 255, 0.6);
  border-radius: 8px;
  padding: 8px 12px;
  margin-bottom: 6px;
  border: 1px solid rgba(255, 255, 255, 0.3);
  transition: all 0.3s ease;
  font-size: 12px;
}

.pod-item:hover {
  background: rgba(255, 255, 255, 0.8);
  transform: translateX(4px);
}

.pod-name {
  color: #333;
  font-weight: 500;
  font-size: 12px;
  min-width: 160px;
  font-family: 'Courier New', monospace;
}

.pod-status {
  width: fit-content;
  min-width: 70px;
}

.pod-resource {
  color: #666;
  font-size: 11px;
  font-family: 'Courier New', monospace;
  min-width: 80px;
}

.pod-resource:first-of-type {
  color: #4facfe;
}

.pod-resource:last-of-type {
  color: #43e97b;
}

/* 动画 */
@keyframes pulse {
  0% { opacity: 1; transform: translate(-50%, -50%) scale(0.8); }
  50% { opacity: 0.5; transform: translate(-50%, -50%) scale(1.1); }
  100% { opacity: 0; transform: translate(-50%, -50%) scale(1.4); }
}

@keyframes statusPulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

@keyframes iconFloat {
  0%, 100% { transform: translateY(0px); }
  50% { transform: translateY(-2px); }
}

@keyframes iconPulse {
  0%, 100% { opacity: 1; transform: scale(1); }
  50% { opacity: 0.7; transform: scale(1.1); }
}

@keyframes progressGlow {
  0% { transform: translateX(-100%); }
  100% { transform: translateX(100%); }
}

@keyframes loadingDot {
  0%, 80%, 100% { transform: scale(0); }
  40% { transform: scale(1); }
}

@keyframes refreshSpin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

/* 响应式设计 */
@media (max-width: 768px) {
  .namespaces-monitoring {
    padding: 12px;
  }
  
  .namespace-card {
    padding: 16px;
  }
  
  .namespace-actions {
    flex-direction: column;
    gap: 8px;
  }
  
  .system-stats {
    flex-direction: column;
    gap: 8px;
  }
}

/* Deep blue dark theme overrides */
.namespaces-monitoring {
  background: transparent;
  min-height: auto;
}

.namespace-card {
  background: linear-gradient(180deg, rgba(30, 41, 59, 0.96), rgba(15, 23, 42, 0.92));
  border: 1px solid var(--border-medium);
  box-shadow: var(--shadow-card);
}

.namespace-card::before {
  background: linear-gradient(90deg, rgba(59, 130, 246, 0), rgba(56, 189, 248, 0.88), rgba(59, 130, 246, 0));
}

.namespace-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 22px 44px rgba(2, 8, 23, 0.34);
  border-color: var(--border-strong);
}

.namespace-icon {
  background: linear-gradient(135deg, var(--color-primary-strong), var(--color-accent));
  box-shadow: 0 10px 24px rgba(37, 99, 235, 0.28);
}

.pulse-ring {
  border-color: rgba(56, 189, 248, 0.4);
}

.namespace-details h4,
.section-title,
.info-value,
.pod-name {
  color: var(--text-primary);
  text-shadow: none;
}

.namespace-name,
.metric-detail,
.stat-label,
.no-data-text,
.info-label,
.stat-details,
.pod-resource {
  color: var(--text-muted);
}

.metric-item,
.stat-item,
.no-metrics,
.detail-section,
.resource-stat-item,
.pod-item {
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid var(--border-subtle);
  box-shadow: none;
}

.metric-item:hover,
.stat-item:hover,
.resource-stat-item:hover,
.pod-item:hover {
  background: rgba(59, 130, 246, 0.08);
  border-color: rgba(59, 130, 246, 0.18);
  transform: none;
}

.metric-label,
.stat-value,
.metric-value {
  color: var(--text-secondary);
  background: none;
  -webkit-text-fill-color: currentColor;
}

.progress-bg,
.custom-progress {
  background: rgba(148, 163, 184, 0.14);
}

.cpu-progress {
  background: linear-gradient(90deg, var(--color-primary), var(--color-accent));
}

.memory-progress {
  background: linear-gradient(90deg, var(--color-teal), #5eead4);
}

.pods-icon,
.services-icon {
  background: linear-gradient(135deg, rgba(37, 99, 235, 0.24), rgba(56, 189, 248, 0.18));
}

.loading-dot {
  background: var(--color-accent);
}

.primary-btn {
  background: linear-gradient(135deg, var(--color-primary-strong), var(--color-primary));
  color: #fff;
  box-shadow: var(--shadow-glow);
}

.secondary-btn {
  background: rgba(15, 23, 42, 0.24);
  color: var(--text-secondary);
  border: 1px solid var(--border-medium);
}

.secondary-btn:hover {
  background: rgba(59, 130, 246, 0.08);
  border-color: var(--border-strong);
}

.namespace-detail-dialog :deep(.el-dialog) {
  background: linear-gradient(180deg, rgba(30, 41, 59, 0.98), rgba(15, 23, 42, 0.98));
  border: 1px solid var(--border-medium);
  box-shadow: var(--shadow-elevated);
}

.namespace-detail-dialog :deep(.el-dialog__header) {
  background: linear-gradient(180deg, rgba(51, 65, 85, 0.98), rgba(30, 41, 59, 0.98));
  border-bottom: 1px solid var(--border-subtle);
}

.namespace-detail-dialog :deep(.el-dialog__title) {
  color: var(--text-primary);
}

.namespace-detail-dialog :deep(.el-dialog__body),
.pod-list {
  color: var(--text-secondary);
}

.section-title {
  border-bottom-color: var(--border-subtle);
}

.section-icon {
  color: #93c5fd;
}
</style>
