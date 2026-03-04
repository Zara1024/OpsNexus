<template>
  <el-card>
    <template #header>
      <div class="toolbar">
        <span>Deployment 管理</span>
        <div class="actions">
          <el-select v-model="clusterId" placeholder="集群" style="width: 220px" @change="handleClusterChange">
            <el-option v-for="c in clusters" :key="c.id" :label="c.name" :value="c.id" />
          </el-select>
          <el-select v-model="namespace" placeholder="命名空间" style="width: 200px" @change="loadData">
            <el-option v-for="ns in namespaces" :key="ns" :label="ns" :value="ns" />
          </el-select>
          <el-button v-permission="'k8s:workload:deployment:create'" type="primary" @click="openCreate">新建</el-button>
        </div>
      </div>
    </template>

    <el-table :data="list" v-loading="loading">
      <el-table-column prop="metadata.name" label="名称" min-width="180" />
      <el-table-column label="副本" width="120">
        <template #default="scope">{{ scope.row.status?.readyReplicas || 0 }}/{{ scope.row.spec?.replicas || 0 }}</template>
      </el-table-column>
      <el-table-column prop="status.availableReplicas" label="可用副本" width="120" />
      <el-table-column label="操作" min-width="300" fixed="right">
        <template #default="scope">
          <el-button v-permission="'k8s:workload:deployment:restart'" link type="warning" @click="restart(scope.row.metadata?.name)">重启</el-button>
          <el-button v-permission="'k8s:workload:deployment:scale'" link type="success" @click="openScale(scope.row)">伸缩</el-button>
          <el-button v-permission="'k8s:workload:deployment:delete'" link type="danger" @click="remove(scope.row.metadata?.name)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="dialogVisible" title="创建 Deployment" width="760px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="名称"><el-input v-model="form.name" /></el-form-item>
        <el-form-item label="镜像"><el-input v-model="form.image" /></el-form-item>
        <el-form-item label="副本数"><el-input-number v-model="form.replicas" :min="0" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitCreate">创建</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="scaleVisible" title="伸缩 Deployment" width="420px">
      <el-form>
        <el-form-item label="副本数" label-width="80px"><el-input-number v-model="scaleReplicas" :min="0" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="scaleVisible = false">取消</el-button>
        <el-button type="primary" @click="submitScale">确认</el-button>
      </template>
    </el-dialog>
  </el-card>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  createDeploymentApi,
  deleteDeploymentApi,
  listDeploymentsApi,
  listK8sClustersApi,
  listK8sNamespacesApi,
  restartDeploymentApi,
  scaleDeploymentApi,
  type K8sClusterItem,
} from '@/api/k8s'

const loading = ref(false)
const clusters = ref<K8sClusterItem[]>([])
const namespaces = ref<string[]>([])
const clusterId = ref<number | undefined>(undefined)
const namespace = ref('default')
const list = ref<any[]>([])

const dialogVisible = ref(false)
const form = reactive({ name: '', image: 'nginx:latest', replicas: 1 })
const scaleVisible = ref(false)
const scaleName = ref('')
const scaleReplicas = ref(1)

const loadClusters = async () => {
  const res = await listK8sClustersApi()
  clusters.value = res.list || []
  if (!clusterId.value && clusters.value.length) clusterId.value = clusters.value[0].id
}

const loadNamespaces = async () => {
  if (!clusterId.value) return
  const res = await listK8sNamespacesApi(clusterId.value)
  namespaces.value = res.list || []
  if (!namespaces.value.includes(namespace.value)) {
    namespace.value = namespaces.value[0] || 'default'
  }
}

const loadData = async () => {
  if (!clusterId.value || !namespace.value) return
  loading.value = true
  try {
    const res = await listDeploymentsApi(clusterId.value, namespace.value)
    list.value = res.list || []
  } finally {
    loading.value = false
  }
}

const handleClusterChange = async () => {
  await loadNamespaces()
  await loadData()
}

const openCreate = () => {
  Object.assign(form, { name: '', image: 'nginx:latest', replicas: 1 })
  dialogVisible.value = true
}

const submitCreate = async () => {
  if (!clusterId.value) return
  const payload = {
    apiVersion: 'apps/v1',
    kind: 'Deployment',
    metadata: { name: form.name, namespace: namespace.value },
    spec: {
      replicas: form.replicas,
      selector: { matchLabels: { app: form.name } },
      template: {
        metadata: { labels: { app: form.name } },
        spec: { containers: [{ name: form.name, image: form.image, ports: [{ containerPort: 80 }] }] },
      },
    },
  }
  await createDeploymentApi(clusterId.value, namespace.value, payload)
  ElMessage.success('创建成功')
  dialogVisible.value = false
  await loadData()
}

const openScale = (row: any) => {
  scaleName.value = row.metadata?.name || ''
  scaleReplicas.value = row.spec?.replicas || 1
  scaleVisible.value = true
}

const submitScale = async () => {
  if (!clusterId.value || !scaleName.value) return
  await scaleDeploymentApi(clusterId.value, namespace.value, scaleName.value, scaleReplicas.value)
  ElMessage.success('伸缩成功')
  scaleVisible.value = false
  await loadData()
}

const restart = async (name: string) => {
  if (!clusterId.value) return
  await restartDeploymentApi(clusterId.value, namespace.value, name)
  ElMessage.success('重启成功')
  await loadData()
}

const remove = async (name: string) => {
  if (!clusterId.value) return
  await ElMessageBox.confirm(`确认删除 ${name} 吗？`, '提示', { type: 'warning' })
  await deleteDeploymentApi(clusterId.value, namespace.value, name)
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
