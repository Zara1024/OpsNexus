<template>
  <el-card>
    <template #header>
      <div class="toolbar">
        <span>K8s 节点管理</span>
        <div class="actions">
          <el-select v-model="clusterId" placeholder="选择集群" style="width: 260px" @change="loadData">
            <el-option v-for="c in clusterOptions" :key="c.id" :label="c.name" :value="c.id" />
          </el-select>
          <el-button @click="loadData">刷新</el-button>
        </div>
      </div>
    </template>

    <el-table :data="list" v-loading="loading">
      <el-table-column prop="name" label="节点名" min-width="220" />
      <el-table-column label="角色" min-width="150">
        <template #default="scope">{{ (scope.row.roles || []).join(', ') }}</template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="120" />
      <el-table-column prop="pod_count" label="Pod数" width="90" />
      <el-table-column label="可调度" width="100">
        <template #default="scope">{{ scope.row.unschedulable ? '否' : '是' }}</template>
      </el-table-column>
      <el-table-column label="操作" width="260" fixed="right">
        <template #default="scope">
          <el-button v-permission="'k8s:node:cordon'" link type="warning" @click="handleCordon(scope.row.name)">封锁</el-button>
          <el-button v-permission="'k8s:node:uncordon'" link type="success" @click="handleUncordon(scope.row.name)">解封</el-button>
          <el-button v-permission="'k8s:node:drain'" link type="danger" @click="handleDrain(scope.row.name)">排水</el-button>
        </template>
      </el-table-column>
    </el-table>
  </el-card>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { cordonK8sNodeApi, drainK8sNodeApi, listK8sClustersApi, listK8sNodesApi, uncordonK8sNodeApi, type K8sClusterItem, type K8sNodeItem } from '@/api/k8s'

const loading = ref(false)
const list = ref<K8sNodeItem[]>([])
const clusterOptions = ref<K8sClusterItem[]>([])
const clusterId = ref<number | undefined>(undefined)

const loadClusters = async () => {
  const res = await listK8sClustersApi()
  clusterOptions.value = res.list || []
  if (!clusterId.value && clusterOptions.value.length) {
    clusterId.value = clusterOptions.value[0].id
  }
}

const loadData = async () => {
  if (!clusterId.value) return
  loading.value = true
  try {
    const res = await listK8sNodesApi(clusterId.value)
    list.value = res.list || []
  } finally {
    loading.value = false
  }
}

const handleCordon = async (name: string) => {
  if (!clusterId.value) return
  await cordonK8sNodeApi(clusterId.value, name)
  ElMessage.success('封锁成功')
  await loadData()
}

const handleUncordon = async (name: string) => {
  if (!clusterId.value) return
  await uncordonK8sNodeApi(clusterId.value, name)
  ElMessage.success('解封成功')
  await loadData()
}

const handleDrain = async (name: string) => {
  if (!clusterId.value) return
  await ElMessageBox.confirm(`确认对节点 ${name} 执行排水吗？`, '提示', { type: 'warning' })
  await drainK8sNodeApi(clusterId.value, name)
  ElMessage.success('排水任务已执行')
  await loadData()
}

onMounted(async () => {
  await loadClusters()
  await loadData()
})
</script>

<style scoped>
.toolbar { display: flex; justify-content: space-between; align-items: center; }
.actions { display: flex; gap: 8px; }
</style>
