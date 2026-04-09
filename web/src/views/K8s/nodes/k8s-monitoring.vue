<template>
  <div class="k8s-monitoring-dashboard">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-title">
        <h2>K8s 监控仪表板</h2>
        <p class="header-desc">实时监控集群资源使用情况</p>
      </div>
      <div class="header-controls">
        <el-select v-model="selectedClusterId" placeholder="选择集群" @change="handleClusterChange">
          <el-option 
            v-for="cluster in clusterList" 
            :key="cluster.id" 
            :label="cluster.name" 
            :value="cluster.id"
          />
        </el-select>
        <el-button :icon="Refresh" @click="refreshAllData" :loading="loading">刷新数据</el-button>
      </div>
    </div>

    <!-- 监控标签页 -->
    <el-tabs v-model="activeTab" @tab-change="handleTabChange">
      <!-- 节点监控 -->
      <el-tab-pane label="节点监控" name="nodes">
        <NodesMonitoring 
          :selected-cluster-id="selectedClusterId"
          ref="nodesMonitoringRef"
        />
      </el-tab-pane>

      <!-- 命名空间监控 -->
      <el-tab-pane label="命名空间监控" name="namespaces">
        <NamespaceMonitoring 
          :selected-cluster-id="selectedClusterId"
          ref="namespaceMonitoringRef"
        />
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import {
  Refresh
} from '@element-plus/icons-vue'
import k8sApi from '@/api/k8s'
import NodesMonitoring from './NodesMonitoring.vue'
import NamespaceMonitoring from './NamespaceMonitoring.vue'

// 响应式数据
const loading = ref(false)
const activeTab = ref('nodes')
const selectedClusterId = ref('')
const clusterList = ref([])

// 组件引用
const nodesMonitoringRef = ref(null)
const namespaceMonitoringRef = ref(null)

// 工具函数
// API调用函数
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
    } else {
      ElMessage.error(responseData.message || '获取集群列表失败')
    }
  } catch (error) {
    console.error('获取集群列表失败:', error)
    ElMessage.warning('无法获取集群列表，请检查后端服务')
  }
}


// 事件处理函数
const handleClusterChange = async () => {
  if (selectedClusterId.value) {
    await loadTabData()
    // 通知组件集群变更
    if (nodesMonitoringRef.value) {
      nodesMonitoringRef.value.handleClusterChange()
    }
    if (namespaceMonitoringRef.value) {
      namespaceMonitoringRef.value.handleClusterChange()
    }
  }
}

const handleTabChange = async (tabName) => {
  activeTab.value = tabName
  await loadTabData()
}

const loadTabData = async () => {
  if (!selectedClusterId.value) return
  
  if (activeTab.value === 'nodes') {
    // 节点监控由 NodesMonitoring 组件处理
    if (nodesMonitoringRef.value) {
      nodesMonitoringRef.value.refreshAllNodes()
    }
  } else if (activeTab.value === 'namespaces') {
    // 命名空间监控由 NamespaceMonitoring 组件处理
    if (namespaceMonitoringRef.value) {
      namespaceMonitoringRef.value.refreshAllData()
    }
  }
}

const refreshAllData = async () => {
  loading.value = true
  try {
    if (activeTab.value === 'nodes' && nodesMonitoringRef.value) {
      await nodesMonitoringRef.value.refreshAllNodes()
    } else if (activeTab.value === 'namespaces' && namespaceMonitoringRef.value) {
      await namespaceMonitoringRef.value.refreshAllData()
    } else {
      await loadTabData()
    }
    ElMessage.success('数据刷新成功')
  } catch (error) {
    ElMessage.error('数据刷新失败')
  } finally {
    loading.value = false
  }
}


// 组件挂载时初始化
onMounted(async () => {
  try {
    await fetchClusterList()
    if (selectedClusterId.value) {
      await loadTabData()
    }
  } catch (error) {
    console.error('页面初始化失败:', error)
  }
})
</script>

<style scoped>
.k8s-monitoring-dashboard {
  padding: 24px;
  background:
    radial-gradient(circle at top left, rgba(59, 130, 246, 0.14), transparent 30%),
    linear-gradient(180deg, rgba(15, 23, 42, 0.96), rgba(2, 6, 23, 0.92));
  min-height: 100vh;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  background: linear-gradient(180deg, rgba(30, 41, 59, 0.98), rgba(15, 23, 42, 0.94));
  padding: 20px 24px;
  border-radius: 18px;
  border: 1px solid var(--border-medium);
  box-shadow: var(--shadow-card);
}

.header-title h2 {
  margin: 0 0 4px 0;
  font-size: 24px;
  font-weight: 600;
  color: var(--text-primary);
}

.header-desc {
  margin: 0;
  color: var(--text-secondary);
  font-size: 14px;
}

.header-controls {
  display: flex;
  gap: 12px;
  align-items: center;
}

.header-controls .el-select {
  width: 200px;
}

.header-controls :deep(.el-select__wrapper) {
  background: var(--bg-input);
  border: 1px solid var(--border-subtle);
  box-shadow: none;
}

.header-controls :deep(.el-select__selected-item) {
  color: var(--text-primary);
}

.header-controls :deep(.el-button) {
  background: linear-gradient(135deg, var(--color-primary-strong), var(--color-primary));
  border-color: transparent;
  color: #fff;
  box-shadow: var(--shadow-glow);
}

.k8s-monitoring-dashboard :deep(.el-tabs__header) {
  margin-bottom: 16px;
}

.k8s-monitoring-dashboard :deep(.el-tabs__nav-wrap::after) {
  background: var(--border-subtle);
}

.k8s-monitoring-dashboard :deep(.el-tabs__item) {
  color: var(--text-muted);
}

.k8s-monitoring-dashboard :deep(.el-tabs__item.is-active) {
  color: #dbeafe;
}

.k8s-monitoring-dashboard :deep(.el-tabs__active-bar) {
  background: linear-gradient(135deg, var(--color-primary-strong), var(--color-primary));
}

</style>
