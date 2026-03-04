<template>
  <el-card>
    <template #header>
      <div class="toolbar">
        <span>命名空间</span>
        <el-select v-model="clusterId" placeholder="选择集群" style="width: 260px" @change="loadData">
          <el-option v-for="c in clusters" :key="c.id" :label="c.name" :value="c.id" />
        </el-select>
      </div>
    </template>

    <el-table :data="list" v-loading="loading">
      <el-table-column prop="name" label="命名空间" min-width="220" />
    </el-table>
  </el-card>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { listK8sClustersApi, listK8sNamespacesApi, type K8sClusterItem } from '@/api/k8s'

const loading = ref(false)
const clusters = ref<K8sClusterItem[]>([])
const clusterId = ref<number | undefined>()
const list = ref<{ name: string }[]>([])

const loadClusters = async () => {
  const res = await listK8sClustersApi()
  clusters.value = res.list || []
  if (!clusterId.value && clusters.value.length) {
    clusterId.value = clusters.value[0].id
  }
}

const loadData = async () => {
  if (!clusterId.value) return
  loading.value = true
  try {
    const res = await listK8sNamespacesApi(clusterId.value)
    list.value = (res.list || []).map((name) => ({ name }))
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  await loadClusters()
  await loadData()
})
</script>

<style scoped>
.toolbar { display: flex; justify-content: space-between; align-items: center; }
</style>
