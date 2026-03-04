<template>
  <el-card>
    <template #header>
      <div class="toolbar">
        <span>K8s 集群管理</span>
        <el-button v-permission="'k8s:cluster:create'" type="primary" @click="openCreate">新增集群</el-button>
      </div>
    </template>

    <el-table :data="list" v-loading="loading">
      <el-table-column prop="name" label="名称" min-width="140" />
      <el-table-column prop="api_server_url" label="API Server" min-width="220" />
      <el-table-column prop="version" label="版本" width="120" />
      <el-table-column prop="node_count" label="节点" width="90" />
      <el-table-column prop="pod_count" label="Pod" width="90" />
      <el-table-column prop="health" label="健康状态" width="120" />
      <el-table-column label="状态" width="100">
        <template #default="scope">{{ scope.row.status === 1 ? '启用' : '禁用' }}</template>
      </el-table-column>
      <el-table-column label="操作" min-width="320" fixed="right">
        <template #default="scope">
          <el-button v-permission="'k8s:cluster:detail'" link type="primary" @click="openDetail(scope.row.id)">详情</el-button>
          <el-button v-permission="'k8s:cluster:update'" link type="primary" @click="openEdit(scope.row)">编辑</el-button>
          <el-button v-permission="'k8s:cluster:test'" link type="success" @click="handleTest(scope.row.id)">连接测试</el-button>
          <el-button v-permission="'k8s:cluster:delete'" link type="danger" @click="handleDelete(scope.row.id)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="dialogVisible" :title="form.id ? '编辑集群' : '新增集群'" width="720px">
      <el-form :model="form" label-width="110px">
        <el-form-item label="集群名称"><el-input v-model="form.name" /></el-form-item>
        <el-form-item label="API Server"><el-input v-model="form.api_server_url" /></el-form-item>
        <el-form-item label="描述"><el-input v-model="form.description" type="textarea" :rows="2" /></el-form-item>
        <el-form-item label="Kubeconfig"><el-input v-model="form.kubeconfig" type="textarea" :rows="8" placeholder="编辑时可留空" /></el-form-item>
        <el-form-item v-if="form.id" label="状态">
          <el-switch v-model="statusSwitch" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submit">保存</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="detailVisible" title="集群详情" width="640px">
      <el-descriptions v-if="detail" :column="2" border>
        <el-descriptions-item label="名称">{{ detail.name }}</el-descriptions-item>
        <el-descriptions-item label="健康状态">{{ detail.health || '-' }}</el-descriptions-item>
        <el-descriptions-item label="API Server">{{ detail.api_server_url }}</el-descriptions-item>
        <el-descriptions-item label="版本">{{ detail.version || '-' }}</el-descriptions-item>
        <el-descriptions-item label="节点数">{{ detail.node_count }}</el-descriptions-item>
        <el-descriptions-item label="Pod数">{{ detail.pod_count }}</el-descriptions-item>
        <el-descriptions-item label="描述" :span="2">{{ detail.description || '-' }}</el-descriptions-item>
      </el-descriptions>
    </el-dialog>
  </el-card>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  createK8sClusterApi,
  deleteK8sClusterApi,
  getK8sClusterApi,
  listK8sClustersApi,
  testK8sClusterApi,
  updateK8sClusterApi,
  type K8sClusterItem,
} from '@/api/k8s'

const loading = ref(false)
const list = ref<K8sClusterItem[]>([])
const dialogVisible = ref(false)
const detailVisible = ref(false)
const detail = ref<any>(null)

const form = reactive<any>({
  id: 0,
  name: '',
  api_server_url: '',
  description: '',
  kubeconfig: '',
  status: 1,
})

const statusSwitch = computed({
  get: () => form.status === 1,
  set: (v: boolean) => (form.status = v ? 1 : 2),
})

const loadData = async () => {
  loading.value = true
  try {
    const res = await listK8sClustersApi()
    list.value = res.list || []
  } finally {
    loading.value = false
  }
}

const openCreate = () => {
  Object.assign(form, { id: 0, name: '', api_server_url: '', description: '', kubeconfig: '', status: 1 })
  dialogVisible.value = true
}

const openEdit = (row: K8sClusterItem) => {
  Object.assign(form, { id: row.id, name: row.name, api_server_url: row.api_server_url, description: row.description, kubeconfig: '', status: row.status })
  dialogVisible.value = true
}

const openDetail = async (id: number) => {
  detail.value = await getK8sClusterApi(id)
  detailVisible.value = true
}

const submit = async () => {
  if (form.id) {
    await updateK8sClusterApi(form.id, form)
  } else {
    await createK8sClusterApi(form)
  }
  ElMessage.success('保存成功')
  dialogVisible.value = false
  await loadData()
}

const handleTest = async (id: number) => {
  await testK8sClusterApi(id)
  ElMessage.success('连接成功')
}

const handleDelete = async (id: number) => {
  await ElMessageBox.confirm('确认删除该集群吗？', '提示', { type: 'warning' })
  await deleteK8sClusterApi(id)
  ElMessage.success('删除成功')
  await loadData()
}

onMounted(loadData)
</script>

<style scoped>
.toolbar { display: flex; justify-content: space-between; align-items: center; }
</style>
