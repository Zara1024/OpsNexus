<template>
  <el-card>
    <template #header>登录日志</template>
    <el-table :data="list" v-loading="loading" size="small">
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="username" label="用户" width="120" />
      <el-table-column prop="ip" label="IP" width="140" />
      <el-table-column prop="method" label="方法" width="90" />
      <el-table-column prop="path" label="路径" min-width="220" />
      <el-table-column prop="response_code" label="响应码" width="90" />
      <el-table-column prop="created_at" label="时间" width="180" />
    </el-table>

    <div class="pager">
      <el-pagination
        background
        layout="total, prev, pager, next"
        :total="total"
        :current-page="query.page"
        :page-size="query.page_size"
        @current-change="handlePageChange"
      />
    </div>
  </el-card>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { listLoginLogsApi, type AuditLogItem } from '@/api/audit'

const loading = ref(false)
const total = ref(0)
const list = ref<AuditLogItem[]>([])
const query = reactive({ page: 1, page_size: 10, keyword: '' })

const loadData = async () => {
  loading.value = true
  try {
    const res = await listLoginLogsApi(query)
    list.value = res.list || []
    total.value = res.total || 0
  } finally {
    loading.value = false
  }
}

const handlePageChange = (page: number) => {
  query.page = page
  loadData()
}

onMounted(loadData)
</script>

<style scoped>
.pager { margin-top: 12px; display: flex; justify-content: flex-end; }
</style>
