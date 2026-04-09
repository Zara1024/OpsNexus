<template>
  <TablePage
    class="ansible-history-page"
    section-title="任务历史"
    section-subtitle="查看 Ansible 任务执行记录、状态与步骤明细"
  >
    <template #header>
      <PageHeader
        eyebrow="Task Center"
        title="任务执行历史"
        subtitle="回溯 Ansible 任务每次运行结果、耗时与步骤日志"
      >
        <template #actions>
          <el-button @click="router.back()">
            <el-icon><ArrowLeft /></el-icon>
            返回任务
          </el-button>
          <el-button type="primary" :loading="loading" @click="fetchHistory">
            <el-icon><Refresh /></el-icon>
            刷新历史
          </el-button>
        </template>
        <template #meta>
          <div class="page-chip-row">
            <span class="platform-chip">任务名称 {{ taskName }}</span>
            <span class="platform-chip">任务 ID {{ taskId || '-' }}</span>
            <span class="platform-chip">总记录 {{ total }}</span>
          </div>
        </template>
        <template #intro>
          <PageIntro
            title="查看说明"
            text="支持按执行状态筛选，并下钻查看每次任务的阶段流转和步骤详情。"
          />
        </template>
      </PageHeader>
    </template>

    <template #toolbar>
      <PageToolbar>
        <el-form :inline="true" :model="queryParams" class="filter-cluster">
          <el-form-item class="filter-field" label="执行状态">
            <el-select
              v-model="queryParams.status"
              placeholder="全部状态"
              clearable
              @change="handleSearch"
            >
              <el-option label="待执行" :value="1" />
              <el-option label="执行中" :value="2" />
              <el-option label="成功" :value="3" />
              <el-option label="失败" :value="4" />
            </el-select>
          </el-form-item>
        </el-form>
        <template #actions>
          <el-button type="primary" @click="handleSearch">
            <el-icon><Search /></el-icon>
            查询
          </el-button>
          <el-button @click="resetQuery">
            <el-icon><Refresh /></el-icon>
            重置
          </el-button>
        </template>
      </PageToolbar>
    </template>

    <template #stats>
      <StatStrip :items="statItems" />
    </template>

    <el-table
      v-loading="loading"
      :data="historyList"
      stripe
      class="history-table"
      empty-text="暂无任务历史记录"
    >
      <el-table-column label="记录 ID" width="120">
        <template #default="{ row }">
          <span class="history-id">#{{ row.ID }}</span>
        </template>
      </el-table-column>
      <el-table-column label="状态" width="110">
        <template #default="{ row }">
          <el-tag :type="statusTagType(row.status)" effect="dark">
            {{ getStatusText(row.status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作人" min-width="130">
        <template #default="{ row }">
          {{ row.OperatorName || 'System' }}
        </template>
      </el-table-column>
      <el-table-column label="开始时间" min-width="180">
        <template #default="{ row }">
          {{ formatTime(row.CreatedAt) }}
        </template>
      </el-table-column>
      <el-table-column label="完成时间" min-width="180">
        <template #default="{ row }">
          {{ formatTime(row.FinishedAt) }}
        </template>
      </el-table-column>
      <el-table-column label="执行耗时" width="130">
        <template #default="{ row }">
          {{ formatDuration(row.TotalDuration) }}
        </template>
      </el-table-column>
      <el-table-column label="步骤数" width="120">
        <template #default="{ row }">
          {{ Array.isArray(row.Works) ? row.Works.length : 0 }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="120" fixed="right">
        <template #default="{ row }">
          <el-button type="primary" link @click="viewLog(row)">查看详情</el-button>
        </template>
      </el-table-column>
    </el-table>

    <template #footer>
      <el-pagination
        :current-page="queryParams.page"
        :page-sizes="[10, 20, 50]"
        :page-size="queryParams.pageSize"
        layout="total, sizes, prev, pager, next, jumper"
        :total="total"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </template>
  </TablePage>

  <el-dialog
    v-model="logDialogVisible"
    title="任务执行详情"
    width="82%"
    top="5vh"
    destroy-on-close
    class="history-detail-dialog"
  >
    <div class="history-detail-shell" v-loading="detailLoading">
      <div class="history-detail-summary">
        <div class="history-detail-card">
          <span class="history-detail-card__label">任务名称</span>
          <span class="history-detail-card__value">{{ taskName }}</span>
        </div>
        <div class="history-detail-card">
          <span class="history-detail-card__label">任务 ID</span>
          <span class="history-detail-card__value">{{ taskId || '-' }}</span>
        </div>
        <div class="history-detail-card">
          <span class="history-detail-card__label">记录 ID</span>
          <span class="history-detail-card__value">{{ currentHistoryId || '-' }}</span>
        </div>
      </div>

      <AnsibleFlowTemp
        v-if="historySteps.length > 0"
        :steps="historySteps"
        :task-id="taskId"
        :history-mode="true"
        :history-id="currentHistoryId"
      />
      <div v-else class="history-empty-state">
        <el-empty description="暂无执行步骤" />
      </div>
    </div>
  </el-dialog>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ArrowLeft, Refresh, Search } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import PageHeader from '@/components/platform/PageHeader.vue'
import PageIntro from '@/components/platform/PageIntro.vue'
import PageToolbar from '@/components/platform/PageToolbar.vue'
import StatStrip from '@/components/platform/StatStrip.vue'
import TablePage from '@/components/platform/TablePage.vue'
import AnsibleFlowTemp from './Job/AnsibleFlowTemp.vue'
import { GetAnsibleHistoryDetail, GetAnsibleTaskDetail, GetAnsibleTaskHistory } from '@/api/task'

const route = useRoute()
const router = useRouter()

const taskId = ref(route.query.id)
const taskName = ref(route.query.name || '未命名任务')
const loading = ref(false)
const detailLoading = ref(false)
const historyList = ref([])
const total = ref(0)
const logDialogVisible = ref(false)
const currentHistoryId = ref(null)
const historySteps = ref([])

const queryParams = reactive({
  id: taskId.value,
  page: 1,
  pageSize: 10,
  status: ''
})

const statItems = computed(() => {
  const rows = historyList.value
  const successCount = rows.filter(item => Number(item.status) === 3).length
  const failureCount = rows.filter(item => Number(item.status) === 4).length
  const runningCount = rows.filter(item => Number(item.status) === 2).length

  return [
    { label: '历史总数', value: total.value, hint: '当前筛选后的记录', tone: 'primary' },
    { label: '成功次数', value: successCount, hint: '已成功完成', tone: 'success' },
    { label: '失败次数', value: failureCount, hint: '执行失败次数', tone: 'danger' },
    { label: '执行中', value: runningCount, hint: '仍在执行的记录', tone: 'warning' }
  ]
})

const resolveFallbackStatus = works => {
  const statuses = Array.isArray(works) ? works.map(work => Number(work.Status ?? work.status ?? 0)) : []
  if (statuses.includes(4)) return 4
  if (statuses.includes(2)) return 2
  if (statuses.length > 0 && statuses.every(status => status === 3)) return 3
  return 1
}

const buildFallbackHistoryRecord = taskInfo => {
  const works = Array.isArray(taskInfo?.Works) ? taskInfo.Works : []
  const status = Number(route.query.status || resolveFallbackStatus(works))
  const createdAt = route.query.createdAt || ''
  const updatedAt = route.query.updatedAt || ''

  return {
    ID: Number(taskId.value || taskInfo?.ID || 0),
    status,
    OperatorName: 'System',
    CreatedAt: createdAt,
    FinishedAt: status === 3 || status === 4 ? updatedAt : '',
    TotalDuration: Number(route.query.totalDuration || 0),
    Works: works
  }
}

const fetchHistory = async () => {
  loading.value = true
  try {
    const res = await GetAnsibleTaskHistory(queryParams)
    const payload = res?.data?.data || {}
    historyList.value = Array.isArray(payload.data)
      ? payload.data
      : Array.isArray(payload.list)
        ? payload.list
        : []
    total.value = Number(payload.total || historyList.value.length || 0)
  } catch (error) {
    if (error?.response?.status === 404) {
      try {
        const taskRes = await GetAnsibleTaskDetail(taskId.value)
        const taskInfo = taskRes?.data?.data?.task_info
        const fallbackRecord = buildFallbackHistoryRecord(taskInfo)
        const matchesStatus = !queryParams.status || Number(queryParams.status) === Number(fallbackRecord.status)

        historyList.value = matchesStatus ? [fallbackRecord] : []
        total.value = historyList.value.length
        return
      } catch (fallbackError) {
        console.error('Ansible history fallback failed:', fallbackError)
      }
    }

    console.error('获取任务历史失败', error)
    ElMessage.error('获取任务历史失败')
    historyList.value = []
    total.value = 0
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  queryParams.page = 1
  fetchHistory()
}

const resetQuery = () => {
  queryParams.status = ''
  queryParams.page = 1
  fetchHistory()
}

const handleSizeChange = value => {
  queryParams.pageSize = value
  queryParams.page = 1
  fetchHistory()
}

const handleCurrentChange = value => {
  queryParams.page = value
  fetchHistory()
}

const viewLog = async row => {
  logDialogVisible.value = true
  detailLoading.value = true
  currentHistoryId.value = row.ID
  historySteps.value = []

  try {
    const historyDetailRes = await GetAnsibleHistoryDetail({
      id: taskId.value,
      historyId: row.ID
    })

    const workHistories = historyDetailRes?.data?.data?.WorkHistories
    if (Array.isArray(workHistories) && workHistories.length) {
      historySteps.value = workHistories.map(work => ({
        task_id: taskId.value,
        work_id: work.WorkID,
        entry_file_name: work.HostName,
        status: work.Status,
        duration: work.Duration,
        original_work: work
      }))
      return
    }

    if (Array.isArray(row.Works) && row.Works.length) {
      historySteps.value = row.Works.map(work => ({
        ...work,
        task_id: taskId.value,
        work_id: work.ID || work.work_id || work.WorkID,
        entry_file_name: work.EntryFileName || work.entry_file_name,
        status: work.Status !== undefined ? work.Status : (work.status !== undefined ? work.status : row.status),
        duration: work.Duration || work.duration
      }))
      return
    }

    const taskRes = await GetAnsibleTaskDetail(taskId.value)
    const staticWorks = taskRes?.data?.data?.task_info?.Works
    if (Array.isArray(staticWorks) && staticWorks.length) {
      historySteps.value = staticWorks.map(work => ({
        ...work,
        task_id: taskId.value,
        work_id: work.workid,
        entry_file_name: work.EntryFileName,
        status: row.status,
        duration: 0
      }))
      return
    }

    ElMessage.warning('未获取到可展示的执行步骤')
  } catch (error) {
    console.error('加载执行详情失败:', error)
    if (Array.isArray(row.Works) && row.Works.length) {
      historySteps.value = row.Works.map(work => ({
        ...work,
        task_id: taskId.value,
        work_id: work.ID || work.work_id || work.WorkID,
        entry_file_name: work.EntryFileName || work.entry_file_name,
        status: work.Status !== undefined ? work.Status : (work.status !== undefined ? work.status : row.status),
        duration: work.Duration || work.duration
      }))
    } else {
      ElMessage.error('加载执行详情失败')
    }
  } finally {
    detailLoading.value = false
  }
}

const getStatusText = status => {
  const map = { 1: '待执行', 2: '执行中', 3: '成功', 4: '失败' }
  return map[Number(status)] || '未知'
}

const statusTagType = status => {
  const map = { 1: 'info', 2: 'warning', 3: 'success', 4: 'danger' }
  return map[Number(status)] || 'info'
}

const formatTime = value => {
  if (!value || String(value).startsWith('0001-01-01')) return '-'
  try {
    return new Date(value).toLocaleString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit',
      hour12: false
    }).replace(/\//g, '-')
  } catch (error) {
    return value
  }
}

const formatDuration = seconds => {
  const totalSeconds = Number(seconds || 0)
  if (!totalSeconds) return '0秒'

  const hours = Math.floor(totalSeconds / 3600)
  const minutes = Math.floor((totalSeconds % 3600) / 60)
  const remainingSeconds = totalSeconds % 60
  const parts = []

  if (hours > 0) parts.push(`${hours}小时`)
  if (minutes > 0) parts.push(`${minutes}分钟`)
  if (remainingSeconds > 0 || parts.length === 0) parts.push(`${remainingSeconds}秒`)

  return parts.join('')
}

onMounted(() => {
  fetchHistory()
})
</script>

<style scoped>
.ansible-history-page :deep(.page-actions) {
  align-items: center;
}

.history-id {
  display: inline-flex;
  align-items: center;
  padding: 4px 10px;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.06);
  border: 1px solid var(--border-subtle);
  color: var(--text-primary);
  font-family: 'Consolas', 'Courier New', monospace;
  font-size: 12px;
  font-weight: 700;
}

.history-detail-shell {
  display: grid;
  gap: 18px;
  min-height: 360px;
}

.history-detail-summary {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  gap: 14px;
}

.history-detail-card {
  padding: 16px 18px;
  border-radius: 18px;
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid var(--border-subtle);
  display: grid;
  gap: 8px;
}

.history-detail-card__label {
  font-size: 12px;
  font-weight: 600;
  color: var(--text-muted);
}

.history-detail-card__value {
  color: var(--text-primary);
  font-size: 16px;
  font-weight: 700;
  word-break: break-word;
}

.history-empty-state {
  display: grid;
  min-height: 280px;
  place-items: center;
  border-radius: 20px;
  border: 1px dashed var(--border-medium);
  background: rgba(255, 255, 255, 0.03);
}

@media (max-width: 768px) {
  .ansible-history-page :deep(.page-actions) {
    width: 100%;
    justify-content: flex-start;
  }
}
</style>
