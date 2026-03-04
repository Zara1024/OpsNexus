<template>
  <el-card>
    <template #header>历史告警</template>
    <el-table :data="list" v-loading="loading">
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="severity" label="等级" width="90" />
      <el-table-column prop="status" label="状态" width="100" />
      <el-table-column prop="starts_at" label="开始时间" width="180" />
      <el-table-column prop="ends_at" label="恢复时间" width="180" />
      <el-table-column prop="labels" label="标签" min-width="260" show-overflow-tooltip />
    </el-table>
  </el-card>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { listAlertHistoryApi, type AlertEventItem } from '@/api/monitor'

const loading = ref(false)
const list = ref<AlertEventItem[]>([])

const loadData = async () => {
  loading.value = true
  try {
    const res = await listAlertHistoryApi()
    list.value = res.list || []
  } finally {
    loading.value = false
  }
}

onMounted(loadData)
</script>
