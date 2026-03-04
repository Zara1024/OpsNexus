<template>
  <el-card>
    <template #header>活跃告警</template>
    <el-table :data="list" v-loading="loading">
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="severity" label="等级" width="90" />
      <el-table-column prop="status" label="状态" width="100" />
      <el-table-column prop="starts_at" label="开始时间" width="180" />
      <el-table-column prop="labels" label="标签" min-width="260" show-overflow-tooltip />
      <el-table-column label="操作" width="180" fixed="right">
        <template #default="scope">
          <el-button v-permission="'monitor:alert:ack'" link type="primary" @click="handleAck(scope.row.id)">确认</el-button>
          <el-button v-permission="'monitor:alert:silence'" link type="warning" @click="handleSilence(scope.row.id)">静默30m</el-button>
        </template>
      </el-table-column>
    </el-table>
  </el-card>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { ackAlertApi, listActiveAlertsApi, silenceAlertApi, type AlertEventItem } from '@/api/monitor'

const loading = ref(false)
const list = ref<AlertEventItem[]>([])

const loadData = async () => {
  loading.value = true
  try {
    const res = await listActiveAlertsApi()
    list.value = res.list || []
  } finally {
    loading.value = false
  }
}

const handleAck = async (id: number) => {
  await ackAlertApi(id)
  ElMessage.success('已确认')
  await loadData()
}

const handleSilence = async (id: number) => {
  await silenceAlertApi(id, 30)
  ElMessage.success('已静默30分钟')
  await loadData()
}

onMounted(loadData)
</script>
