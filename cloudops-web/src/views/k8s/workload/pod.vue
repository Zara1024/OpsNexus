<template>
  <el-card>
    <template #header>
      <div class="toolbar">
        <span>Pod 管理</span>
        <div class="actions">
          <el-select v-model="clusterId" placeholder="集群" style="width: 220px" @change="handleClusterChange">
            <el-option v-for="c in clusters" :key="c.id" :label="c.name" :value="c.id" />
          </el-select>
          <el-select v-model="namespace" placeholder="命名空间" style="width: 200px" @change="loadData">
            <el-option v-for="ns in namespaces" :key="ns" :label="ns" :value="ns" />
          </el-select>
        </div>
      </div>
    </template>

    <el-table :data="list" v-loading="loading">
      <el-table-column prop="metadata.name" label="Pod" min-width="220" />
      <el-table-column prop="status.phase" label="状态" width="120" />
      <el-table-column prop="spec.nodeName" label="节点" min-width="180" />
      <el-table-column label="操作" width="220" fixed="right">
        <template #default="scope">
          <el-button v-permission="'k8s:workload:pod:log'" link type="primary" @click="showLogs(scope.row.metadata?.name)">日志</el-button>
          <el-button v-permission="'k8s:workload:pod:delete'" link type="danger" @click="remove(scope.row.metadata?.name)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="logVisible" :title="`Pod 日志: ${logPod}`" width="880px">
      <el-input v-model="logs" type="textarea" :rows="24" readonly />
    </el-dialog>
  </el-card>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { deletePodApi, getPodLogsApi, listK8sClustersApi, listK8sNamespacesApi, listPodsApi, type K8sClusterItem } from '@/api/k8s'

const loading = ref(false)
const clusters = ref<K8sClusterItem[]>([])
const namespaces = ref<string[]>([])
const clusterId = ref<number | undefined>()
const namespace = ref('default')
const list = ref<any[]>([])

const logVisible = ref(false)
const logPod = ref('')
const logs = ref('')

const loadClusters = async () => {
  const res = await listK8sClustersApi()
  clusters.value = res.list || []
  if (!clusterId.value && clusters.value.length) clusterId.value = clusters.value[0].id
}

const loadNamespaces = async () => {
  if (!clusterId.value) return
  const res = await listK8sNamespacesApi(clusterId.value)
  namespaces.value = res.list || []
  if (!namespaces.value.includes(namespace.value)) namespace.value = namespaces.value[0] || 'default'
}

const loadData = async () => {
  if (!clusterId.value) return
  loading.value = true
  try {
    const res = await listPodsApi(clusterId.value, namespace.value)
    list.value = res.list || []
  } finally {
    loading.value = false
  }
}

const handleClusterChange = async () => {
  await loadNamespaces()
  await loadData()
}

const showLogs = async (pod: string) => {
  if (!clusterId.value) return
  const res = await getPodLogsApi(clusterId.value, namespace.value, pod, { tail_lines: 300 })
  logPod.value = pod
  logs.value = res.logs || ''
  logVisible.value = true
}

const remove = async (pod: string) => {
  if (!clusterId.value) return
  await ElMessageBox.confirm(`确认删除 Pod ${pod} 吗？`, '提示', { type: 'warning' })
  await deletePodApi(clusterId.value, namespace.value, pod)
  ElMessage.success('删除成功')
  await loadData()
}

onMounted(async () => {
  await loadClusters()
  await loadNamespaces()
  await loadData()
})
</script>

<style scoped>
.toolbar { display: flex; justify-content: space-between; align-items: center; }
.actions { display: flex; gap: 8px; }
</style>
