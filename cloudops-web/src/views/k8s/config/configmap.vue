<template>
  <el-card>
    <template #header>
      <div class="toolbar">
        <span>ConfigMap 管理</span>
        <div class="actions">
          <el-select v-model="clusterId" placeholder="集群" style="width: 220px" @change="handleClusterChange">
            <el-option v-for="c in clusters" :key="c.id" :label="c.name" :value="c.id" />
          </el-select>
          <el-select v-model="namespace" placeholder="命名空间" style="width: 200px" @change="loadData">
            <el-option v-for="ns in namespaces" :key="ns" :label="ns" :value="ns" />
          </el-select>
          <el-button v-permission="'k8s:config:configmap:save'" type="primary" @click="openCreate">新建</el-button>
        </div>
      </div>
    </template>

    <el-table :data="list" v-loading="loading">
      <el-table-column prop="metadata.name" label="名称" min-width="220" />
      <el-table-column label="键数量" width="120">
        <template #default="scope">{{ Object.keys(scope.row.data || {}).length }}</template>
      </el-table-column>
      <el-table-column label="操作" width="220" fixed="right">
        <template #default="scope">
          <el-button v-permission="'k8s:config:configmap:save'" link type="primary" @click="openEdit(scope.row)">编辑</el-button>
          <el-button v-permission="'k8s:config:configmap:delete'" link type="danger" @click="remove(scope.row.metadata?.name)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="dialogVisible" title="保存 ConfigMap" width="760px">
      <el-form :model="form" label-width="90px">
        <el-form-item label="名称"><el-input v-model="form.name" /></el-form-item>
        <el-form-item label="内容JSON"><el-input v-model="form.dataText" type="textarea" :rows="14" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submit">保存</el-button>
      </template>
    </el-dialog>
  </el-card>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { deleteConfigMapApi, listConfigMapsApi, listK8sClustersApi, listK8sNamespacesApi, saveConfigMapApi, type K8sClusterItem } from '@/api/k8s'

const loading = ref(false)
const clusters = ref<K8sClusterItem[]>([])
const namespaces = ref<string[]>([])
const clusterId = ref<number | undefined>()
const namespace = ref('default')
const list = ref<any[]>([])

const dialogVisible = ref(false)
const form = reactive({ name: '', dataText: '{\n  "key": "value"\n}' })

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
    const res = await listConfigMapsApi(clusterId.value, namespace.value)
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
  form.name = ''
  form.dataText = '{\n  "key": "value"\n}'
  dialogVisible.value = true
}

const openEdit = (row: any) => {
  form.name = row.metadata?.name || ''
  form.dataText = JSON.stringify(row.data || {}, null, 2)
  dialogVisible.value = true
}

const submit = async () => {
  if (!clusterId.value) return
  const data = JSON.parse(form.dataText || '{}')
  await saveConfigMapApi(clusterId.value, namespace.value, { metadata: { name: form.name, namespace: namespace.value }, data })
  ElMessage.success('保存成功')
  dialogVisible.value = false
  await loadData()
}

const remove = async (name: string) => {
  if (!clusterId.value) return
  await ElMessageBox.confirm(`确认删除 ${name} 吗？`, '提示', { type: 'warning' })
  await deleteConfigMapApi(clusterId.value, namespace.value, name)
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
