<template>
  <el-card>
    <template #header>
      <div class="toolbar">
        <span>操作日志</span>
        <el-button v-permission="'audit:log:export'" @click="handleExport">导出CSV</el-button>
      </div>
    </template>

    <el-form :inline="true" class="search">
      <el-form-item>
        <el-input v-model="query.keyword" placeholder="用户/模块/路径" clearable />
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="loadData">查询</el-button>
      </el-form-item>
    </el-form>

    <el-table :data="list" v-loading="loading" size="small">
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="username" label="用户" width="120" />
      <el-table-column prop="module" label="模块" width="120" />
      <el-table-column prop="action" label="动作" width="120" />
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
import { reactive, ref, onMounted } from 'vue'
import { listAuditLogsApi, exportAuditLogsApi, type AuditLogItem } from '@/api/audit'

const loading = ref(false)
const total = ref(0)
const list = ref<AuditLogItem[]>([])
const query = reactive({ page: 1, page_size: 10, keyword: '' })

const loadData = async () => {
  loading.value = true
  try {
    const res = await listAuditLogsApi(query)
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

const handleExport = async () => {
  const blob = await exportAuditLogsApi(query)
  const url = window.URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = 'audit_logs.csv'
  a.click()
  window.URL.revokeObjectURL(url)
}

onMounted(loadData)
</script>

<style scoped>
.toolbar { display: flex; justify-content: space-between; align-items: center; }
.search { margin-bottom: 12px; }
.pager { margin-top: 12px; display: flex; justify-content: flex-end; }
</style>
