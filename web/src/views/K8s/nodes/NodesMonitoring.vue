<template>
  <div class="nodes-monitoring">
    <el-row :gutter="20" style="margin-bottom: 20px;">
      <el-col :span="8" v-for="node in nodesList" :key="node.name">
        <div class="node-card" v-loading="nodesLoading">
          <!-- 卡片头部 -->
          <div class="card-header">
            <div class="node-info-header">
              <div class="node-icon">
                <img src="@/assets/image/k8s.svg" alt="k8s" class="k8s-icon" />
                <div class="pulse-ring"></div>
              </div>
              <div class="node-details">
                <h4>{{ node.name }}</h4>
                <el-tag 
                  :type="getNodeStatusTag(node.status)" 
                  size="small"
                  :class="['status-tag', node.status.toLowerCase()]"
                >
                  <i class="status-dot"></i>
                  {{ getNodeStatusText(node.status) }}
                </el-tag>
              </div>
            </div>
            <div class="node-ip">{{ node.name }}</div>
          </div>
          
          <div v-if="nodeMetrics[node.name]" class="node-metrics">
            <!-- CPU 使用率 -->
            <div class="metric-item">
              <div class="metric-header">
                <img src="@/assets/image/cpu.svg" alt="CPU" class="metric-icon" />
                <span class="metric-label">CPU使用率</span>
                <span class="metric-value">{{ getCpuUsageRate(nodeMetrics[node.name]) }}%</span>
              </div>
              <div class="metric-progress">
                <div class="progress-bg">
                  <div 
                    class="progress-fill cpu-progress" 
                    :style="{ width: getCpuUsageRate(nodeMetrics[node.name]) + '%' }"
                  >
                    <div class="progress-glow"></div>
                  </div>
                </div>
              </div>
              <div class="metric-detail">
                {{ formatCpuValue(nodeMetrics[node.name].usage?.cpu) }} / {{ formatCpuValue(nodeMetrics[node.name].capacity?.cpu) }}
              </div>
            </div>
            
            <!-- 内存使用率 -->
            <div class="metric-item">
              <div class="metric-header">
                <img src="@/assets/image/cpu.svg" alt="Memory" class="metric-icon" />
                <span class="metric-label">内存使用率</span>
                <span class="metric-value">{{ getMemoryUsageRate(nodeMetrics[node.name]) }}%</span>
              </div>
              <div class="metric-progress">
                <div class="progress-bg">
                  <div 
                    class="progress-fill memory-progress" 
                    :style="{ width: getMemoryUsageRate(nodeMetrics[node.name]) + '%' }"
                  >
                    <div class="progress-glow"></div>
                  </div>
                </div>
              </div>
              <div class="metric-detail">
                {{ formatMemoryValue(nodeMetrics[node.name].usage?.memory) }} / {{ formatMemoryValue(nodeMetrics[node.name].capacity?.memory) }}
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
                  <span class="stat-value">{{ getDisplayedPodCount(node) }}</span>
                </div>
              </div>
              <div class="stat-item">
                <div class="stat-icon time-icon">
                  <i class="icon-pulse"></i>
                </div>
                <div class="stat-content">
                  <span class="stat-label">更新时间</span>
                  <span class="stat-value">{{ formatTimestamp(nodeMetrics[node.name].timestamp) }}</span>
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
            <span class="no-data-text">正在获取监控数据...</span>
          </div>

          <!-- 操作按钮 -->
          <div class="node-actions">
            <button 
              class="action-btn primary-btn" 
              @click="viewNodeDetail(node)"
            >
              <View class="btn-icon" />
              <span>详细信息</span>
              <div class="btn-glow"></div>
            </button>
            <button 
              class="action-btn secondary-btn" 
              @click="refreshNodeMetrics(node.name)"
            >
              <Refresh class="btn-icon refresh-icon" />
              <span>刷新</span>
            </button>
          </div>
        </div>
      </el-col>
    </el-row>

    <!-- 节点详情弹框 -->
    <el-dialog
      v-model="nodeDetailVisible"
      :title="`节点监控 - ${selectedNode?.name}`"
      width="70%"
      class="node-detail-dialog"
    >
      <div v-if="selectedNodeMetrics" v-loading="nodeDetailLoading" class="node-detail-content">
        <el-row :gutter="16" style="margin-bottom: 16px;">
          <el-col :span="12">
            <div class="resource-card">
              <div class="card-header-detail">
                <img src="@/assets/image/cpu.svg" alt="Resource" class="header-icon" />
                <h3>资源使用情况</h3>
              </div>
              <div class="resource-metrics">
                <div class="metric-row-detail">
                  <div class="metric-info">
                    <img src="@/assets/image/cpu.svg" alt="CPU" class="metric-icon-detail" />
                    <span class="metric-name">CPU:</span>
                    <span class="metric-percentage">{{ getCpuUsageRate(selectedNodeMetrics) }}%</span>
                  </div>
                  <div class="progress-container">
                    <div class="progress-bg">
                      <div 
                        class="progress-fill cpu-progress" 
                        :style="{ width: getCpuUsageRate(selectedNodeMetrics) + '%' }"
                      >
                        <div class="progress-glow"></div>
                      </div>
                    </div>
                  </div>
                  <div class="metric-values">
                    {{ formatCpuValue(selectedNodeMetrics.usage?.cpu) }} / {{ formatCpuValue(selectedNodeMetrics.capacity?.cpu) }}
                  </div>
                </div>
                <div class="metric-row-detail">
                  <div class="metric-info">
                    <img src="@/assets/image/cpu.svg" alt="Memory" class="metric-icon-detail" />
                    <span class="metric-name">内存:</span>
                    <span class="metric-percentage">{{ getMemoryUsageRate(selectedNodeMetrics) }}%</span>
                  </div>
                  <div class="progress-container">
                    <div class="progress-bg">
                      <div 
                        class="progress-fill memory-progress" 
                        :style="{ width: getMemoryUsageRate(selectedNodeMetrics) + '%' }"
                      >
                        <div class="progress-glow"></div>
                      </div>
                    </div>
                  </div>
                  <div class="metric-values">
                    {{ formatMemoryValue(selectedNodeMetrics.usage?.memory) }} / {{ formatMemoryValue(selectedNodeMetrics.capacity?.memory) }}
                  </div>
                </div>
              </div>
            </div>
          </el-col>
          <el-col :span="12">
            <el-card>
              <template #header>系统信息</template>
              <div v-if="selectedNodeMetrics.systemInfo" class="system-info">
                <el-descriptions :column="1" size="small">
                  <el-descriptions-item label="操作系统">
                    {{ selectedNodeMetrics.systemInfo.osImage }}
                  </el-descriptions-item>
                  <el-descriptions-item label="内核版本">
                    {{ selectedNodeMetrics.systemInfo.kernelVersion }}
                  </el-descriptions-item>
                  <el-descriptions-item label="Kubelet版本">
                    {{ selectedNodeMetrics.systemInfo.kubeletVersion }}
                  </el-descriptions-item>
                  <el-descriptions-item label="容器运行时">
                    {{ selectedNodeMetrics.systemInfo.containerRuntimeVersion }}
                  </el-descriptions-item>
                </el-descriptions>
              </div>
            </el-card>
          </el-col>
        </el-row>

        <!-- Pod列表 -->
        <div class="pod-list-card">
          <div class="card-header-detail">
            <div class="pod-icon">
              <i class="icon-pod"></i>
            </div>
            <h3>运行的Pod ({{ getFormattedPodCount(selectedNodeMetrics) }})</h3>
          </div>
          <div class="pod-table-container">
            <div v-if="getFormattedPodList(selectedNodeMetrics).length > 0" class="pod-table">
              <div class="pod-header">
                <div class="pod-col pod-name">Pod名称</div>
                <div class="pod-col pod-namespace">命名空间</div>
                <div class="pod-col pod-cpu">CPU使用</div>
                <div class="pod-col pod-memory">内存使用</div>
                <div class="pod-col pod-cpu-rate">CPU使用率</div>
                <div class="pod-col pod-memory-rate">内存使用率</div>
              </div>
              <div class="pod-body">
                <div 
                  v-for="(pod, index) in getFormattedPodList(selectedNodeMetrics)" 
                  :key="index"
                  class="pod-row"
                >
                  <div class="pod-col pod-name">
                    <div class="pod-name-content">
                      <div class="pod-status-dot"></div>
                      {{ pod.podName }}
                    </div>
                  </div>
                  <div class="pod-col pod-namespace">
                    <span class="namespace-tag">{{ pod.namespace }}</span>
                  </div>
                  <div class="pod-col pod-cpu">
                    <span class="resource-value">{{ formatCpuValue(pod.usage?.cpu) }}</span>
                  </div>
                  <div class="pod-col pod-memory">
                    <span class="resource-value">{{ formatMemoryValue(pod.usage?.memory) }}</span>
                  </div>
                  <div class="pod-col pod-cpu-rate">
                    <div class="usage-rate-container">
                      <div class="mini-progress">
                        <div 
                          class="mini-progress-fill cpu-progress" 
                          :style="{ width: (pod.usageRate?.cpuRate || 0) + '%' }"
                        ></div>
                      </div>
                      <span class="rate-text">{{ Math.round(pod.usageRate?.cpuRate || 0) }}%</span>
                    </div>
                  </div>
                  <div class="pod-col pod-memory-rate">
                    <div class="usage-rate-container">
                      <div class="mini-progress">
                        <div 
                          class="mini-progress-fill memory-progress" 
                          :style="{ width: (pod.usageRate?.memoryRate || 0) + '%' }"
                        ></div>
                      </div>
                      <span class="rate-text">{{ Math.round(pod.usageRate?.memoryRate || 0) }}%</span>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <div v-else class="no-pods">
              <div class="no-pods-icon">
                <i class="icon-empty"></i>
              </div>
              <p class="no-pods-text">暂无运行中的Pod</p>
            </div>
          </div>
        </div>
      </div>
      
      <template #footer>
        <el-button @click="nodeDetailVisible = false">关闭</el-button>
        <el-button type="primary" @click="refreshNodeDetail">刷新数据</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>





import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import {
  Refresh,
  View
} from '@element-plus/icons-vue'
import k8sApi from '@/api/k8s'

// Props
const props = defineProps({
  selectedClusterId: {
    type: [String, Number],
    required: true
  }
})

// 鑺傜偣鐩戞帶鐩稿叧
const nodesList = ref([])
const nodeMetrics = ref({})
const nodesLoading = ref(false)
const nodeDetailVisible = ref(false)
const nodeDetailLoading = ref(false)
const selectedNode = ref(null)
const selectedNodeMetrics = ref(null)

// 宸ュ叿鍑芥暟
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

const getNodeStatusText = (status) => {
  switch (status) {
    case 'Ready':
      return '就绪'
    case 'NotReady':
      return '未就绪'
    case 'SchedulingDisabled':
      return '调度禁用'
    default:
      return '未知'
  }
}

const formatCpuValue = (cpu) => {
  if (!cpu) return '0m'
  
  // 澶勭悊绾崇鍗曚綅 (n)
  if (cpu.includes('n')) {
    const value = parseInt(cpu.replace('n', ''))
    return Math.round(value / 1000000) + 'm'
  }
  
  // 澶勭悊姣牳鍗曚綅 (m)  
  if (cpu.includes('m')) {
    return cpu
  }
  
  // 澶勭悊鏁存牳鍗曚綅
  const value = parseFloat(cpu)
  if (!isNaN(value)) {
    return Math.round(value * 1000) + 'm'
  }
  
  return cpu
}

const formatMemoryValue = (memory) => {
  if (!memory) return '0Mi'
  
  // 澶勭悊Ki鍗曚綅
  if (memory.includes('Ki')) {
    const value = parseInt(memory.replace('Ki', ''))
    if (value >= 1024) {
      return Math.round(value / 1024) + 'Mi'
    }
    return memory
  }
  
  // 澶勭悊Mi鍗曚綅
  if (memory.includes('Mi')) {
    return memory
  }
  
  // 澶勭悊Gi鍗曚綅
  if (memory.includes('Gi')) {
    const value = parseFloat(memory.replace('Gi', ''))
    return Math.round(value * 1024) + 'Mi'
  }
  
  return memory
}

const getCpuUsageRate = (metrics) => {
  if (!metrics) return 0
  
  if (metrics.usageRate?.cpuRate > 0) {
    return Math.round(metrics.usageRate.cpuRate)
  }
  
  const usage = metrics.usage?.cpu
  const capacity = metrics.capacity?.cpu
  
  if (!usage || !capacity) return 0
  
  try {
    // 灏唘sage杞崲涓烘鏍?    let usageMillicores = 0
    if (usage.includes('n')) {
      usageMillicores = parseInt(usage.replace('n', '')) / 1000000
    } else if (usage.includes('m')) {
      usageMillicores = parseInt(usage.replace('m', ''))
    } else {
      usageMillicores = parseFloat(usage) * 1000
    }
    
    // 灏哻apacity杞崲涓烘鏍?    let capacityMillicores = 0
    if (capacity.includes('m')) {
      capacityMillicores = parseInt(capacity.replace('m', ''))
    } else if (capacity.includes('n')) {
      capacityMillicores = parseInt(capacity.replace('n', '')) / 1000000
    } else {
      capacityMillicores = parseFloat(capacity) * 1000
    }
    
    if (capacityMillicores > 0) {
      const rate = Math.round((usageMillicores / capacityMillicores) * 100)
      return Math.min(Math.max(rate, 0), 100)
    }
  } catch (error) {
    console.error('CPU使用鐜囪绠楅敊璇?', error)
  }
  
  return 0
}

const getMemoryUsageRate = (metrics) => {
  if (!metrics) return 0
  
  if (metrics.usageRate?.memoryRate > 0) {
    return Math.round(metrics.usageRate.memoryRate)
  }
  
  const usage = metrics.usage?.memory
  const capacity = metrics.capacity?.memory
  
  if (!usage || !capacity) return 0
  
  try {
    // 灏唘sage杞崲涓篕i
    let usageKi = 0
    if (usage.includes('Ki')) {
      usageKi = parseInt(usage.replace('Ki', ''))
    } else if (usage.includes('Mi')) {
      usageKi = parseInt(usage.replace('Mi', '')) * 1024
    } else if (usage.includes('Gi')) {
      usageKi = parseFloat(usage.replace('Gi', '')) * 1024 * 1024
    }
    
    // 灏哻apacity杞崲涓篕i
    let capacityKi = 0
    if (capacity.includes('Ki')) {
      capacityKi = parseInt(capacity.replace('Ki', ''))
    } else if (capacity.includes('Mi')) {
      capacityKi = parseInt(capacity.replace('Mi', '')) * 1024
    } else if (capacity.includes('Gi')) {
      capacityKi = parseFloat(capacity.replace('Gi', '')) * 1024 * 1024
    }
    
    if (capacityKi > 0) {
      const rate = Math.round((usageKi / capacityKi) * 100)
      return Math.min(Math.max(rate, 0), 100)
    }
  } catch (error) {
    console.error('内存使用鐜囪绠楅敊璇?', error)
  }
  
  return 0
}

// 鏍煎紡鍖朠od鏁伴噺
const getFormattedPodCount = (metrics) => {
  if (!metrics) return 0
  
  if (metrics.podMetrics && metrics.podMetrics.length > 0) {
    return metrics.podMetrics.length
  }
  
  if (metrics.podCount > 0) {
    return metrics.podCount
  }
  
  return 0
}

// 鑾峰彇鏍煎紡鍖栫殑Pod鍒楄〃
const getFormattedPodList = (metrics) => {
  if (!metrics) return []
  
  // 濡傛灉鏈夊疄闄呯殑Pod鏁版嵁锛岃繑鍥炴牸寮忓寲鍚庣殑鏁版嵁
  if (metrics.podMetrics && metrics.podMetrics.length > 0) {
    return metrics.podMetrics.map(pod => ({
      podName: pod.podName || pod.name || 'unknown-pod',
      namespace: pod.namespace || 'default',
      usage: {
        cpu: formatCpuValue(pod.usage?.cpu || '0m'),
        memory: formatMemoryValue(pod.usage?.memory || '0Mi')
      },
      usageRate: {
        cpuRate: pod.usageRate?.cpuRate || 0,
        memoryRate: pod.usageRate?.memoryRate || 0
      }
    }))
  }

  return []
}

const getDisplayedPodCount = (node) => {
  if (!node) return 0
  return nodeMetrics.value[node.name]?.podCount || node.podCount || 0
}

const getNodeStatusTag = (status) => {
  const tagMap = {
    'Ready': 'success',
    'NotReady': 'danger',
    'Unknown': 'warning',
    'SchedulingDisabled': 'info'
  }
  return tagMap[status] || 'info'
}

// API璋冪敤鍑芥暟
const fetchNodesList = async () => {
  if (!props.selectedClusterId) return
  
  try {
    nodesLoading.value = true
    const response = await k8sApi.getClusterNodes(props.selectedClusterId)
    const responseData = response.data || response
    
    if (responseData.code === 200 || responseData.success) {
      const nodes = responseData.data?.list || responseData.data || []
      nodesList.value = nodes.map(node => ({
        name: node.nodeName || node.name,
        status: node.status || 'Unknown',
        role: node.role || 'worker',
        podCount: node.podMetrics?.allocated || 0
      }))
      
      await Promise.all(nodesList.value.map(node => fetchNodeMetrics(node.name)))
    }
  } catch (error) {
    console.error('获取节点列表失败:', error)
    ElMessage.error('获取节点列表失败')
  } finally {
    nodesLoading.value = false
  }
}

const fetchNodeMetrics = async (nodeName) => {
  if (!props.selectedClusterId || !nodeName) return
  
  try {
    const response = await k8sApi.getNodeMetrics(props.selectedClusterId, nodeName)
    const responseData = response.data || response
    
    if (responseData.code === 200 || responseData.success) {
      const metrics = responseData.data || {}
      nodeMetrics.value[nodeName] = {
        usage: {
          cpu: metrics.usage?.cpu || '0m',
          memory: metrics.usage?.memory || '0Mi'
        },
        capacity: {
          cpu: metrics.capacity?.cpu || '0',
          memory: metrics.capacity?.memory || '0Gi'
        },
        usageRate: {
          cpuRate: parseFloat(metrics.cpuUsageRate) || 0,
          memoryRate: parseFloat(metrics.memoryUsageRate) || 0
        },
        podCount: metrics.podCount || nodesList.value.find(node => node.name === nodeName)?.podCount || 0,
        timestamp: metrics.timestamp || Date.now(),
        systemInfo: metrics.systemInfo,
        podMetrics: metrics.podMetrics || []
      }
    }
  } catch (error) {
    console.error(`获取节点 ${nodeName} 监控数据失败:`, error)
    nodeMetrics.value[nodeName] = {
      usage: { cpu: '0m', memory: '0Mi' },
      capacity: { cpu: '0', memory: '0Gi' },
      usageRate: { cpuRate: 0, memoryRate: 0 },
      podCount: 0,
      timestamp: Date.now()
    }
  }
}


const fetchNodeDetailPods = async (nodeName) => {
  if (!props.selectedClusterId || !nodeName) return null

  const response = await k8sApi.getNodeDetail(props.selectedClusterId, nodeName)
  const responseData = response.data || response
  if (responseData.code !== 200 && !responseData.success) {
    throw new Error(responseData.message || `获取节点 ${nodeName} 详情失败`)
  }

  const detail = responseData.data || {}
  const pods = Array.isArray(detail.pods) ? detail.pods : []
  const currentMetrics = nodeMetrics.value[nodeName] || {}
  const nextMetrics = {
    ...currentMetrics,
    podCount: pods.length || currentMetrics.podCount || nodesList.value.find(node => node.name === nodeName)?.podCount || 0,
    podMetrics: pods.map(pod => ({
      podName: pod.name || 'unknown-pod',
      namespace: pod.namespace || 'default',
      usage: {
        cpu: pod.cpuUsage || '0m',
        memory: pod.memUsage || '0Mi'
      },
      usageRate: {
        cpuRate: 0,
        memoryRate: 0
      }
    })),
    systemInfo: currentMetrics.systemInfo || {
      osImage: detail.configuration?.osImage || '',
      kernelVersion: detail.configuration?.kernelVersion || '',
      kubeletVersion: detail.configuration?.kubeletVersion || '',
      containerRuntimeVersion: detail.configuration?.containerRuntimeVersion || ''
    }
  }

  nodeMetrics.value[nodeName] = nextMetrics
  return nextMetrics
}

const refreshNodeMetrics = async (nodeName) => {
  await fetchNodeMetrics(nodeName)
  ElMessage.success(`节点 ${nodeName} 监控数据已刷新`)
}

const viewNodeDetail = async (node) => {
  selectedNode.value = node
  selectedNodeMetrics.value = nodeMetrics.value[node.name] || null
  nodeDetailVisible.value = true
  nodeDetailLoading.value = true

  try {
    await fetchNodeMetrics(node.name)
    selectedNodeMetrics.value = await fetchNodeDetailPods(node.name)
  } catch (error) {
    selectedNodeMetrics.value = nodeMetrics.value[node.name] || selectedNodeMetrics.value
    console.error(`获取节点 ${node.name} Pod 详情失败:`, error)
  } finally {
    nodeDetailLoading.value = false
  }
}

const refreshNodeDetail = async () => {
  if (selectedNode.value) {
    nodeDetailLoading.value = true
    try {
      await fetchNodeMetrics(selectedNode.value.name)
      selectedNodeMetrics.value = await fetchNodeDetailPods(selectedNode.value.name)
    } catch (error) {
      selectedNodeMetrics.value = nodeMetrics.value[selectedNode.value.name]
      console.error(`刷新节点 ${selectedNode.value.name} Pod 详情失败:`, error)
    } finally {
      nodeDetailLoading.value = false
    }
    ElMessage.success('监控数据已刷新')
  }
}

// 暴露给父组件的方法
const refreshAllNodes = async () => {
  await fetchNodesList()
}

// 鐩戝惉闆嗙兢鍙樺寲
const handleClusterChange = async () => {
  await fetchNodesList()
}

// 鏆撮湶鏂规硶缁欑埗缁勪欢
defineExpose({
  refreshAllNodes,
  handleClusterChange
})





</script>

<style scoped>
.nodes-monitoring {
  padding: 0;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  min-height: 100vh;
  padding: 20px;
  border-radius: 16px;
}

/* 鑺傜偣鍗＄墖 */
.node-card {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-radius: 20px;
  padding: 24px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
  overflow: hidden;
}

.node-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 4px;
  background: linear-gradient(90deg, #4facfe 0%, #00f2fe 100%);
  border-radius: 20px 20px 0 0;
}

.node-card:hover {
  transform: translateY(-8px) scale(1.02);
  box-shadow: 0 16px 64px rgba(0, 0, 0, 0.15);
  border-color: rgba(79, 172, 254, 0.4);
}

/* 鍗＄墖澶撮儴 */
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 20px;
}

.node-info-header {
  display: flex;
  align-items: center;
  gap: 16px;
}

.node-icon {
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
  z-index: 2;
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

@keyframes pulse {
  0% { opacity: 1; transform: translate(-50%, -50%) scale(0.8); }
  50% { opacity: 0.5; transform: translate(-50%, -50%) scale(1.1); }
  100% { opacity: 0; transform: translate(-50%, -50%) scale(1.4); }
}

.node-details h4 {
  margin: 0 0 8px 0;
  font-size: 18px;
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

@keyframes statusPulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

.node-ip {
  font-family: 'Courier New', monospace;
  font-size: 13px;
  color: #666;
  padding: 4px 8px;
  background: rgba(0, 0, 0, 0.05);
  border-radius: 6px;
  border: 1px solid rgba(0, 0, 0, 0.1);
}

/* 鎸囨爣鍖哄煙 */
.node-metrics {
  margin-bottom: 20px;
}

.metric-item {
  margin-bottom: 20px;
  background: rgba(255, 255, 255, 0.6);
  border-radius: 12px;
  padding: 16px;
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
  gap: 12px;
  margin-bottom: 12px;
}

.metric-icon {
  width: 20px;
  height: 20px;
  filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.1));
  animation: iconFloat 3s ease-in-out infinite;
}

@keyframes iconFloat {
  0%, 100% { transform: translateY(0px); }
  50% { transform: translateY(-2px); }
}

.metric-label {
  font-size: 14px;
  font-weight: 500;
  color: #333;
  flex: 1;
}

.metric-value {
  font-size: 16px;
  font-weight: 700;
  background: linear-gradient(45deg, #4facfe, #00f2fe);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

/* 鑷畾涔夎繘搴︽潯 */
.progress-bg {
  width: 100%;
  height: 8px;
  background: rgba(0, 0, 0, 0.1);
  border-radius: 4px;
  overflow: hidden;
  position: relative;
}

.progress-fill {
  height: 100%;
  border-radius: 4px;
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

@keyframes progressGlow {
  0% { transform: translateX(-100%); }
  100% { transform: translateX(100%); }
}

.metric-detail {
  font-size: 11px;
  color: #888;
  text-align: center;
  margin-top: 8px;
  font-family: 'Courier New', monospace;
}

/* 绯荤粺缁熻 */
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
  background: rgba(255, 255, 255, 0.4);
  padding: 12px;
  border-radius: 10px;
  border: 1px solid rgba(255, 255, 255, 0.3);
  transition: all 0.3s ease;
}

.stat-item:hover {
  background: rgba(255, 255, 255, 0.6);
  transform: scale(1.02);
}

.stat-icon {
  width: 32px;
  height: 32px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
}

.pods-icon {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.time-icon {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
}

.icon-dot {
  width: 8px;
  height: 8px;
  background: white;
  border-radius: 50%;
  animation: iconDot 1.5s infinite;
}

.icon-pulse {
  width: 10px;
  height: 10px;
  background: white;
  border-radius: 50%;
  animation: iconPulse 2s infinite;
}

@keyframes iconDot {
  0%, 100% { transform: scale(1); }
  50% { transform: scale(1.2); }
}

@keyframes iconPulse {
  0% { transform: scale(1); opacity: 1; }
  50% { transform: scale(1.3); opacity: 0.7; }
  100% { transform: scale(1); opacity: 1; }
}

.stat-content {
  flex: 1;
}

.stat-label {
  display: block;
  font-size: 11px;
  color: #666;
  margin-bottom: 2px;
}

.stat-value {
  display: block;
  font-size: 14px;
  font-weight: 600;
  color: #333;
}

/* 鏃犳暟鎹姸鎬?*/
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

@keyframes loadingDot {
  0%, 80%, 100% { transform: scale(0); }
  40% { transform: scale(1); }
}

.no-data-text {
  color: #666;
  font-size: 14px;
}

/* 鎿嶄綔鎸夐挳 */
.node-actions {
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
}

.refresh-icon {
  animation: none;
  transition: transform 0.3s ease;
}

.secondary-btn:hover .refresh-icon {
  animation: refreshSpin 0.8s ease-in-out;
}

@keyframes refreshSpin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
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

/* 瀵硅瘽妗嗘牱寮?*/
.node-detail-dialog :deep(.el-dialog) {
  border-radius: 20px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.2);
  max-height: 90vh;
  overflow: hidden;
}

.node-detail-dialog :deep(.el-dialog__header) {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
  color: white;
  border-top-left-radius: 20px;
  border-top-right-radius: 20px;
  padding: 24px;
}

.node-detail-dialog :deep(.el-dialog__title) {
  color: white;
  font-weight: 600;
  font-size: 18px;
}

.node-detail-dialog :deep(.el-dialog__body) {
  padding: 24px;
  max-height: 70vh;
  overflow-y: auto;
}

.node-detail-content {
  padding: 0;
}

/* 璧勬簮鍗＄墖 */
.resource-card {
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(15px);
  border-radius: 16px;
  padding: 20px;
  border: 1px solid rgba(255, 255, 255, 0.3);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
}

.card-header-detail {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 20px;
  padding-bottom: 12px;
  border-bottom: 1px solid rgba(0, 0, 0, 0.1);
}

.header-icon, .metric-icon-detail {
  width: 24px;
  height: 24px;
  filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.1));
}

.card-header-detail h3 {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: #333;
}

.resource-metrics {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.metric-row-detail {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 16px;
  background: rgba(255, 255, 255, 0.5);
  border-radius: 12px;
  border: 1px solid rgba(255, 255, 255, 0.3);
  transition: all 0.3s ease;
}

.metric-row-detail:hover {
  background: rgba(255, 255, 255, 0.7);
  transform: translateX(4px);
}

.metric-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.metric-name {
  font-weight: 500;
  color: #333;
  flex: 1;
}

.metric-percentage {
  font-size: 18px;
  font-weight: 700;
  background: linear-gradient(45deg, #4facfe, #00f2fe);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.progress-container {
  margin: 4px 0;
}

.metric-values {
  font-size: 12px;
  color: #666;
  font-family: 'Courier New', monospace;
  text-align: center;
}

.system-info {
  padding: 0;
}

/* Pod鍒楄〃鍗＄墖 */
.pod-list-card {
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(15px);
  border-radius: 16px;
  padding: 20px;
  border: 1px solid rgba(255, 255, 255, 0.3);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  margin-top: 20px;
}

.pod-icon {
  width: 32px;
  height: 32px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.icon-pod {
  width: 12px;
  height: 12px;
  background: white;
  border-radius: 2px;
  position: relative;
}

.icon-pod::before {
  content: '';
  position: absolute;
  width: 8px;
  height: 8px;
  background: rgba(255, 255, 255, 0.7);
  border-radius: 1px;
  top: 2px;
  left: 2px;
}

.pod-table-container {
  margin-top: 16px;
}

.pod-table {
  background: rgba(255, 255, 255, 0.6);
  border-radius: 12px;
  overflow: hidden;
  border: 1px solid rgba(255, 255, 255, 0.3);
}

.pod-header {
  display: grid;
  grid-template-columns: 2fr 1fr 1fr 1fr 1fr 1fr;
  gap: 12px;
  padding: 12px 16px;
  background: rgba(79, 172, 254, 0.1);
  font-weight: 600;
  font-size: 13px;
  color: #333;
  border-bottom: 1px solid rgba(0, 0, 0, 0.1);
}

.pod-body {
  max-height: 300px;
  overflow-y: auto;
}

.pod-row {
  display: grid;
  grid-template-columns: 2fr 1fr 1fr 1fr 1fr 1fr;
  gap: 12px;
  padding: 12px 16px;
  border-bottom: 1px solid rgba(0, 0, 0, 0.05);
  transition: all 0.3s ease;
}

.pod-row:hover {
  background: rgba(79, 172, 254, 0.05);
}

.pod-row:last-child {
  border-bottom: none;
}

.pod-col {
  display: flex;
  align-items: center;
  font-size: 12px;
}

.pod-name-content {
  display: flex;
  align-items: center;
  gap: 8px;
}

.pod-status-dot {
  width: 6px;
  height: 6px;
  background: #2ecc71;
  border-radius: 50%;
  animation: statusPulse 2s infinite;
}

.namespace-tag {
  padding: 2px 8px;
  background: rgba(103, 126, 234, 0.1);
  color: #667eea;
  border-radius: 12px;
  font-size: 11px;
  border: 1px solid rgba(103, 126, 234, 0.2);
}

.resource-value {
  font-family: 'Courier New', monospace;
  color: #333;
  font-weight: 500;
}

.usage-rate-container {
  display: flex;
  align-items: center;
  gap: 8px;
  width: 100%;
}

.mini-progress {
  flex: 1;
  height: 4px;
  background: rgba(0, 0, 0, 0.1);
  border-radius: 2px;
  overflow: hidden;
}

.mini-progress-fill {
  height: 100%;
  border-radius: 2px;
  transition: width 0.8s ease;
}

.rate-text {
  font-size: 11px;
  font-weight: 600;
  color: #333;
  min-width: 35px;
  text-align: right;
}

.no-pods {
  text-align: center;
  padding: 40px 20px;
  color: #666;
}

.no-pods-icon {
  width: 48px;
  height: 48px;
  margin: 0 auto 16px;
  background: rgba(0, 0, 0, 0.1);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.icon-empty {
  width: 20px;
  height: 20px;
  background: rgba(0, 0, 0, 0.3);
  border-radius: 2px;
}

.no-pods-text {
  margin: 0;
  font-size: 14px;
  color: #888;
}

/* 鍝嶅簲寮忚璁?*/
@media (max-width: 768px) {
  .nodes-monitoring {
    padding: 12px;
  }
  
  .node-card {
    padding: 16px;
  }
  
  .node-actions {
    flex-direction: column;
    gap: 8px;
  }
  
  .system-stats {
    flex-direction: column;
    gap: 8px;
  }
  
  .metric-row {
    flex-direction: column;
    align-items: stretch;
    gap: 8px;
  }
  
  .metric-row span:last-child {
    text-align: left;
  }
}

/* Deep blue dark theme overrides */
.nodes-monitoring {
  background: transparent;
  min-height: auto;
}

.node-card,
.resource-card,
.pod-list-card {
  background: linear-gradient(180deg, rgba(30, 41, 59, 0.96), rgba(15, 23, 42, 0.92));
  border: 1px solid var(--border-medium);
  box-shadow: var(--shadow-card);
}

.node-card::before {
  background: linear-gradient(90deg, rgba(59, 130, 246, 0), rgba(56, 189, 248, 0.88), rgba(59, 130, 246, 0));
}

.node-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 22px 44px rgba(2, 8, 23, 0.34);
  border-color: var(--border-strong);
}

.node-icon {
  background: linear-gradient(135deg, var(--color-primary-strong), var(--color-accent));
  box-shadow: 0 10px 24px rgba(37, 99, 235, 0.28);
}

.pulse-ring {
  border-color: rgba(56, 189, 248, 0.4);
}

.node-details h4,
.card-header-detail h3,
.metric-name,
.resource-value,
.rate-text {
  color: var(--text-primary);
  text-shadow: none;
}

.node-ip,
.metric-detail,
.stat-label,
.no-data-text,
.metric-label,
.no-pods,
.no-pods-text,
.metric-values {
  color: var(--text-muted);
}

.node-ip,
.metric-item,
.stat-item,
.no-metrics,
.metric-row-detail,
.pod-row,
.pod-header,
.pod-table {
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid var(--border-subtle);
  box-shadow: none;
}

.metric-item:hover,
.stat-item:hover,
.metric-row-detail:hover,
.pod-row:hover {
  background: rgba(59, 130, 246, 0.08);
  border-color: rgba(59, 130, 246, 0.18);
  transform: none;
}

.metric-value,
.metric-percentage,
.stat-value {
  color: var(--text-secondary);
  background: none;
  -webkit-text-fill-color: currentColor;
}

.progress-bg,
.mini-progress {
  background: rgba(148, 163, 184, 0.14);
}

.cpu-progress {
  background: linear-gradient(90deg, var(--color-primary), var(--color-accent));
}

.memory-progress {
  background: linear-gradient(90deg, var(--color-teal), #5eead4);
}

.pods-icon,
.time-icon {
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

.pod-table {
  background: rgba(15, 23, 42, 0.92);
  border: 1px solid rgba(148, 163, 184, 0.22);
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.04);
}

.pod-header {
  background: rgba(30, 41, 59, 0.96);
  border-bottom: 1px solid rgba(148, 163, 184, 0.24);
  color: #f8fafc;
}

.pod-body {
  background: rgba(15, 23, 42, 0.88);
}

.pod-row {
  background: rgba(15, 23, 42, 0.72);
  border-bottom: 1px solid rgba(71, 85, 105, 0.42);
}

.pod-row:hover {
  background: rgba(30, 41, 59, 0.92);
  border-color: rgba(96, 165, 250, 0.24);
}

.pod-col {
  color: #e2e8f0;
}

.namespace-tag {
  background: rgba(59, 130, 246, 0.18);
  color: #dbeafe;
  border-color: rgba(96, 165, 250, 0.28);
}

.resource-value,
.rate-text {
  color: #f8fafc;
}

.node-detail-dialog :deep(.el-dialog) {
  background: linear-gradient(180deg, rgba(30, 41, 59, 0.98), rgba(15, 23, 42, 0.98));
  border: 1px solid var(--border-medium);
  box-shadow: var(--shadow-elevated);
}

.node-detail-dialog :deep(.el-dialog__header) {
  background: linear-gradient(180deg, rgba(51, 65, 85, 0.98), rgba(30, 41, 59, 0.98));
  border-bottom: 1px solid var(--border-subtle);
}

.node-detail-dialog :deep(.el-dialog__title) {
  color: var(--text-primary);
}

.node-detail-dialog :deep(.el-dialog__body),
.system-info,
.pod-table {
  color: var(--text-secondary);
}

.card-header-detail {
  border-bottom-color: var(--border-subtle);
}

.no-pods-icon {
  background: rgba(255, 255, 255, 0.08);
}

.icon-empty {
  background: rgba(148, 163, 184, 0.48);
}
</style>

