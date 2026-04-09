<template>
  <TablePage
    class="taskjob-management task-job-page"
    section-title="任务列表"
    section-subtitle="统一管理普通任务与定时任务，快速查看调度状态、执行频次和模板编排。"
    wide
  >
    <template #header>
      <PageHeader
        eyebrow="Task Center"
        title="任务作业"
        subtitle="围绕任务定义、执行主机和调度开关建立统一入口，减少跨页面排查与操作成本。"
      >
        <template #actions>
          <el-button :loading="loading" @click="loadTasks">
            <el-icon><Refresh /></el-icon>
            刷新列表
          </el-button>
          <el-button type="primary" v-authority="['task:job:add']" @click="handleCreate">
            <el-icon><Plus /></el-icon>
            新增任务
          </el-button>
        </template>
        <template #meta>
          <div class="page-chip-row">
            <span class="platform-chip">当前筛选 {{ filterModeLabel }}</span>
            <span class="platform-chip">任务总数 {{ total }}</span>
            <span class="platform-chip">定时任务 {{ scheduledCount }}</span>
          </div>
        </template>
        <template #intro>
          <PageIntro
            title="使用建议"
            text="先用名称、类型和状态缩小范围，再进入任务详情链路启动脚本；新增任务时优先补齐模板与主机，再配置定时表达式。"
          />
        </template>
      </PageHeader>
    </template>

    <template #toolbar>
      <PageToolbar>
        <el-form :inline="true" :model="queryParams" class="filter-cluster">
          <el-form-item class="filter-field" label="任务名称">
            <el-input
              v-model="queryParams.name"
              placeholder="请输入任务名称"
              clearable
              @keyup.enter="handleSearch"
            />
          </el-form-item>
          <el-form-item class="filter-field" label="任务类型">
            <el-select v-model="queryParams.type" placeholder="请选择任务类型" clearable>
              <el-option v-for="item in taskTypeOptions" :key="item.value" :label="item.label" :value="item.value" />
            </el-select>
          </el-form-item>
          <el-form-item class="filter-field" label="任务状态">
            <el-select v-model="queryParams.status" placeholder="请选择任务状态" clearable>
              <el-option v-for="item in taskStatusOptions" :key="item.value" :label="item.label" :value="item.value" />
            </el-select>
          </el-form-item>
        </el-form>
        <template #actions>
          <el-button type="primary" @click="handleSearch">
            <el-icon><Search /></el-icon>
            搜索
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
      :data="tasks"
      stripe
      class="task-job-table"
      empty-text="当前没有匹配的任务作业"
    >
      <el-table-column label="任务名称" min-width="180" show-overflow-tooltip>
        <template #default="{ row }">
          <el-button type="primary" link class="task-name-link" @click="handleOpenFlow(row)">
            {{ row.name || '未命名任务' }}
          </el-button>
        </template>
      </el-table-column>
      <el-table-column label="类型" width="120">
        <template #default="{ row }">
          <el-tag :type="taskTypeTagType(row.type)">{{ getTaskTypeName(row.type) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="状态" width="110">
        <template #default="{ row }">
          <el-tag :type="taskStatusTagType(row.status)" effect="dark">
            {{ getStatusText(row.status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="定时表达式" min-width="150" show-overflow-tooltip>
        <template #default="{ row }">
          <el-tag v-if="row.type === 2 && row.cronExpr" size="small" type="info">{{ row.cronExpr }}</el-tag>
          <span v-else class="muted">-</span>
        </template>
      </el-table-column>
      <el-table-column label="执行耗时" width="110">
        <template #default="{ row }">
          {{ formatDuration(row.duration) }}
        </template>
      </el-table-column>
      <el-table-column label="执行次数" width="100">
        <template #default="{ row }">
          {{ row.executeCount || 0 }}
        </template>
      </el-table-column>
      <el-table-column label="任务开关" width="110" align="center">
        <template #default="{ row }">
          <el-switch
            v-if="row.type === 2"
            v-model="row.isActive"
            inline-prompt
            active-text="开"
            inactive-text="关"
            :disabled="row.status === 1"
            @change="handleToggleTask(row)"
          />
          <span v-else class="muted">-</span>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="180" fixed="right">
        <template #default="{ row }">
          <div class="row-actions">
            <el-button type="primary" link v-authority="['task:job:start']" @click="handleOpenFlow(row)">
              启动
            </el-button>
            <el-button type="danger" link v-authority="['task:job:delete']" @click="handleDelete(row.id)">
              删除
            </el-button>
          </div>
        </template>
      </el-table-column>
    </el-table>

    <template #footer>
      <el-pagination
        :current-page="queryParams.page"
        :page-sizes="[10, 50, 100]"
        :page-size="queryParams.pageSize"
        layout="total, sizes, prev, pager, next, jumper"
        :total="total"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </template>
  </TablePage>

  <el-dialog
    v-model="formVisible"
    class="modern-dialog task-job-form-dialog"
    title="新增任务"
    width="860px"
    destroy-on-close
    @closed="formDialogClosed"
  >
    <el-form ref="taskFormRef" :model="currentTask" :rules="taskRules" label-width="108px" class="task-job-form">
      <div class="dialog-grid dialog-grid--two">
        <el-form-item label="任务名称" prop="name">
          <el-input v-model="currentTask.name" placeholder="请输入任务名称" />
        </el-form-item>
        <el-form-item label="任务类型" prop="type">
          <el-select v-model="currentTask.type" placeholder="请选择任务类型" style="width: 100%">
            <el-option v-for="item in taskTypeOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
      </div>

      <div v-if="currentTask.type === 2" class="schedule-note">
        <div class="schedule-note__label">调度提示</div>
        <div class="schedule-note__text">
          {{ nextExecutionTime ? `下次执行时间：${nextExecutionTime}` : '请输入合法 Cron 表达式后可预估最近一次执行时间。' }}
        </div>
      </div>

      <el-form-item v-if="currentTask.type === 2" label="定时表达式" prop="cronExpr">
        <el-input
          v-model="currentTask.cronExpr"
          placeholder="例如：0 */5 * * * *"
          @change="calculateNextExecution"
        />
      </el-form-item>

      <div class="platform-data-grid">
        <section class="platform-surface selection-surface">
          <div class="selection-surface__header">
            <div>
              <div class="selection-surface__title">任务模板</div>
              <div class="selection-surface__subtitle">支持多选，保存后会按选择顺序形成任务链路。</div>
            </div>
            <el-button type="primary" plain @click="handleAddTemplate">选择模板</el-button>
          </div>
          <el-form-item class="selection-form-item" prop="shell">
            <div class="selection-list">
              <div v-if="selectedTemplateDetails.length" class="selection-tag-list">
                <div v-for="template in selectedTemplateDetails" :key="template.id" class="selection-chip-card">
                  <div>
                    <div class="selection-chip-card__title">{{ template.name }}</div>
                    <div class="selection-chip-card__meta">{{ getTemplateTypeName(template.type) }}</div>
                  </div>
                  <el-button :icon="Close" text @click="removeTemplate(template.id)" />
                </div>
              </div>
              <div v-else class="selection-empty">尚未选择任务模板</div>
            </div>
          </el-form-item>
        </section>

        <section class="platform-surface selection-surface">
          <div class="selection-surface__header">
            <div>
              <div class="selection-surface__title">执行主机</div>
              <div class="selection-surface__subtitle">从 CMDB 主机列表中选择执行目标，可跨分组追加。</div>
            </div>
            <el-button type="primary" plain @click="handleAddHost">选择主机</el-button>
          </div>
          <el-form-item class="selection-form-item" prop="hostIDs">
            <div class="selection-list">
              <div v-if="selectedHosts.length" class="selection-tag-list">
                <div v-for="host in selectedHosts" :key="host.id" class="selection-chip-card">
                  <div>
                    <div class="selection-chip-card__title">{{ host.name }}</div>
                    <div class="selection-chip-card__meta">{{ host.privateIp || host.publicIp || '-' }}</div>
                  </div>
                  <el-button :icon="Close" text @click="removeHost(host.id)" />
                </div>
              </div>
              <div v-else class="selection-empty">尚未选择执行主机</div>
            </div>
          </el-form-item>
        </section>
      </div>

      <el-form-item label="备注信息">
        <el-input
          v-model="currentTask.remark"
          type="textarea"
          :rows="3"
          maxlength="300"
          show-word-limit
          placeholder="请输入备注信息（可选）"
        />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="formVisible = false">取消</el-button>
      <el-button type="primary" @click="handleSubmit">保存任务</el-button>
    </template>
  </el-dialog>

  <el-dialog
    v-model="showTemplateDialog"
    class="compact-dialog template-picker-dialog"
    title="选择任务模板"
    width="560px"
    append-to-body
  >
    <div class="picker-shell">
      <div class="picker-note">可多选模板，推荐先筛选模板类型，再组合执行顺序。</div>
      <el-select
        v-model="templateDraft"
        multiple
        filterable
        collapse-tags
        collapse-tags-tooltip
        placeholder="请选择任务模板"
        style="width: 100%"
      >
        <el-option
          v-for="template in templates"
          :key="template.id"
          :label="`${template.name} · ${getTemplateTypeName(template.type)}`"
          :value="template.id"
        />
      </el-select>
    </div>
    <template #footer>
      <el-button @click="showTemplateDialog = false">取消</el-button>
      <el-button type="primary" @click="applyTemplateSelection">确定</el-button>
    </template>
  </el-dialog>

  <CreateTaskHost
    v-model="showHostDialog"
    :selected-hosts="selectedHosts"
    @hosts-selected="handleHostSelected"
  />
  <JobFlow ref="jobFlowRef" />
</template>

<script setup>
import { computed, nextTick, onMounted, reactive, ref, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Close, Plus, Refresh, Search } from '@element-plus/icons-vue'
import PageHeader from '@/components/platform/PageHeader.vue'
import PageIntro from '@/components/platform/PageIntro.vue'
import PageToolbar from '@/components/platform/PageToolbar.vue'
import StatStrip from '@/components/platform/StatStrip.vue'
import TablePage from '@/components/platform/TablePage.vue'
import {
  CreateTask,
  DeleteTask,
  GetAllTemplates,
  GetNextExecutionTime,
  GetTasksByName,
  GetTasksByStatus,
  GetTasksByType,
  ListTasks,
  PauseScheduledTask,
  ResumeScheduledTask
} from '@/api/task'
import CreateTaskHost from './Job/CreateTaskHost.vue'
import JobFlow from './Job/JobFlow.vue'

const taskTypeOptions = [
  { label: '普通任务', value: 1 },
  { label: '定时任务', value: 2 }
]

const taskStatusOptions = [
  { label: '等待中', value: 1 },
  { label: '运行中', value: 2 },
  { label: '成功', value: 3 },
  { label: '异常', value: 4 },
  { label: '已暂停', value: 5 }
]

const createEmptyQuery = () => ({
  page: 1,
  pageSize: 10,
  name: '',
  type: '',
  status: ''
})

const createEmptyTask = () => ({
  name: '',
  type: 1,
  shell: '',
  hostIDs: '',
  cronExpr: '',
  remark: ''
})

const queryParams = reactive(createEmptyQuery())
const currentTask = reactive(createEmptyTask())
const tasks = ref([])
const templates = ref([])
const selectedTemplates = ref([])
const templateDraft = ref([])
const selectedHosts = ref([])
const loading = ref(false)
const formVisible = ref(false)
const showTemplateDialog = ref(false)
const showHostDialog = ref(false)
const total = ref(0)
const nextExecutionTime = ref('')
const taskFormRef = ref(null)
const jobFlowRef = ref(null)

const taskRules = reactive({
  name: [{ required: true, message: '请输入任务名称', trigger: 'blur' }],
  type: [{ required: true, message: '请选择任务类型', trigger: 'change' }],
  shell: [{
    validator: (_rule, _value, callback) => {
      if (selectedTemplates.value.length > 0) {
        callback()
        return
      }
      callback(new Error('请选择任务模板'))
    },
    trigger: 'change'
  }],
  hostIDs: [{
    validator: (_rule, _value, callback) => {
      if (selectedHosts.value.length > 0) {
        callback()
        return
      }
      callback(new Error('请选择执行主机'))
    },
    trigger: 'change'
  }],
  cronExpr: [{
    validator: (_rule, value, callback) => {
      if (currentTask.type === 2 && !String(value || '').trim()) {
        callback(new Error('请输入定时表达式'))
        return
      }
      callback()
    },
    trigger: 'blur'
  }]
})

const scheduledCount = computed(() => tasks.value.filter(item => item.type === 2).length)
const runningCount = computed(() => tasks.value.filter(item => item.status === 2).length)
const pausedCount = computed(() => tasks.value.filter(item => item.status === 5).length)
const selectedTemplateDetails = computed(() => selectedTemplates.value.map(id => {
  const matched = templates.value.find(item => item.id === id)
  return matched || { id, name: `模板 #${id}`, type: 1 }
}))

const filterModeLabel = computed(() => {
  const parts = []
  if (queryParams.name.trim()) {
    parts.push(`名称：${queryParams.name.trim()}`)
  }
  if (queryParams.type !== '') {
    parts.push(`类型：${getTaskTypeName(queryParams.type)}`)
  }
  if (queryParams.status !== '') {
    parts.push(`状态：${getStatusText(queryParams.status)}`)
  }
  return parts.length ? parts.join(' / ') : '全部任务'
})

const statItems = computed(() => [
  { label: '任务总数', value: total.value, hint: '当前筛选结果', tone: 'primary' },
  { label: '定时任务', value: scheduledCount.value, hint: '带 Cron 调度的任务', tone: 'warning' },
  { label: '运行中', value: runningCount.value, hint: '正在执行的任务', tone: 'success' },
  { label: '已暂停', value: pausedCount.value, hint: '可恢复的定时任务', tone: 'danger' }
])

const parseCollection = response => {
  const root = response?.data
  if (Array.isArray(root)) {
    return { list: root, total: root.length }
  }
  if (Array.isArray(root?.data)) {
    return { list: root.data, total: Number(root.total || root.data.length || 0) }
  }
  if (Array.isArray(root?.list)) {
    return { list: root.list, total: Number(root.total || root.list.length || 0) }
  }
  if (Array.isArray(root?.data?.list)) {
    return { list: root.data.list, total: Number(root.data.total || root.total || root.data.list.length || 0) }
  }
  if (Array.isArray(root?.data?.data)) {
    return { list: root.data.data, total: Number(root.data.total || root.total || root.data.data.length || 0) }
  }
  return { list: [], total: 0 }
}

const isSuccessResponse = response => {
  const code = response?.data?.code
  return code === undefined || code === 200
}

const normalizeTask = item => {
  const status = Number(item?.status ?? item?.Status ?? 1)
  return {
    id: item?.id ?? item?.ID ?? '',
    name: item?.name ?? item?.Name ?? '',
    type: Number(item?.type ?? item?.Type ?? 1),
    shell: item?.shell ?? item?.Shell ?? '',
    hostIDs: String(item?.host_ids ?? item?.hostIDs ?? item?.HostIDs ?? ''),
    cronExpr: item?.cron_expr ?? item?.cronExpr ?? item?.CronExpr ?? '',
    status,
    duration: Number(item?.duration ?? item?.Duration ?? 0),
    executeCount: Number(item?.execute_count ?? item?.executeCount ?? item?.ExecuteCount ?? 0),
    nextRunTime: item?.next_run_time ?? item?.nextRunTime ?? item?.NextRunTime ?? '',
    remark: item?.remark ?? item?.Remark ?? '',
    createdAt: item?.created_at ?? item?.createdAt ?? item?.CreatedAt ?? '',
    isActive: status === 2
  }
}

const normalizeTemplate = item => ({
  id: item?.id ?? item?.ID ?? '',
  name: item?.name ?? item?.Name ?? '',
  type: Number(item?.type ?? item?.Type ?? 1)
})

const normalizeHost = item => ({
  id: item?.id ?? item?.ID ?? '',
  name: item?.name ?? item?.Name ?? item?.hostName ?? '未命名主机',
  privateIp: item?.privateIp ?? item?.PrivateIp ?? item?.private_ip ?? '',
  publicIp: item?.publicIp ?? item?.PublicIp ?? item?.public_ip ?? '',
  remark: item?.remark ?? item?.Remark ?? ''
})

const getTaskTypeName = type => taskTypeOptions.find(item => item.value === Number(type))?.label || '未知类型'
const getTemplateTypeName = type => ({ 1: 'Shell', 2: 'Python', 3: 'Ansible' }[Number(type)] || 'Shell')
const getStatusText = status => taskStatusOptions.find(item => item.value === Number(status))?.label || '未知状态'
const taskTypeTagType = type => ({ 1: 'info', 2: 'warning', 3: 'success', 4: '' }[Number(type)] || 'info')
const taskStatusTagType = status => ({ 1: 'info', 2: 'warning', 3: 'success', 4: 'danger', 5: 'info' }[Number(status)] || 'info')

const formatTime = value => {
  if (!value) return '-'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return String(value)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
    hour12: false
  }).replace(/\//g, '-')
}

const formatDuration = value => `${Number(value || 0)}s`

const loadTasks = async () => {
  loading.value = true
  try {
    const hasName = Boolean(queryParams.name.trim())
    const hasType = queryParams.type !== ''
    const hasStatus = queryParams.status !== ''
    let response

    if (hasName && !hasType && !hasStatus) {
      response = await GetTasksByName({
        page: queryParams.page,
        pageNum: queryParams.page,
        size: queryParams.pageSize,
        pageSize: queryParams.pageSize,
        name: queryParams.name.trim()
      })
    } else if (!hasName && hasType && !hasStatus) {
      response = await GetTasksByType({
        page: queryParams.page,
        pageNum: queryParams.page,
        size: queryParams.pageSize,
        pageSize: queryParams.pageSize,
        type: queryParams.type
      })
    } else if (!hasName && !hasType && hasStatus) {
      response = await GetTasksByStatus({
        page: queryParams.page,
        pageNum: queryParams.page,
        size: queryParams.pageSize,
        pageSize: queryParams.pageSize,
        status: queryParams.status
      })
    } else {
      response = await ListTasks({
        page: queryParams.page,
        pageSize: queryParams.pageSize,
        name: hasName ? queryParams.name.trim() : undefined,
        type: hasType ? queryParams.type : undefined,
        status: hasStatus ? queryParams.status : undefined
      })
    }

    const { list, total: count } = parseCollection(response)
    tasks.value = list.map(normalizeTask)
    total.value = count
  } catch (error) {
    console.error('获取任务列表失败:', error)
    tasks.value = []
    total.value = 0
    ElMessage.error(error?.response?.data?.message || error?.message || '获取任务列表失败')
  } finally {
    loading.value = false
  }
}

const fetchTemplates = async () => {
  try {
    const response = await GetAllTemplates()
    const { list } = parseCollection(response)
    templates.value = list.map(normalizeTemplate)
  } catch (error) {
    console.error('获取模板列表失败:', error)
    templates.value = []
  }
}

const resetCurrentTask = () => {
  Object.assign(currentTask, createEmptyTask())
  selectedTemplates.value = []
  templateDraft.value = []
  selectedHosts.value = []
  nextExecutionTime.value = ''
}

const handleSearch = () => {
  queryParams.page = 1
  loadTasks()
}

const resetQuery = () => {
  Object.assign(queryParams, createEmptyQuery())
  loadTasks()
}

const handleSizeChange = value => {
  queryParams.pageSize = value
  queryParams.page = 1
  loadTasks()
}

const handleCurrentChange = value => {
  queryParams.page = value
  loadTasks()
}

const handleCreate = async () => {
  resetCurrentTask()
  formVisible.value = true
  await nextTick()
  taskFormRef.value?.clearValidate()
}

const calculateNextExecution = async () => {
  if (currentTask.type !== 2 || !currentTask.cronExpr.trim()) {
    nextExecutionTime.value = ''
    return
  }

  try {
    const response = await GetNextExecutionTime({ cron: currentTask.cronExpr.trim() })
    nextExecutionTime.value = response?.data?.data?.next_execution_time
      ? formatTime(response.data.data.next_execution_time)
      : '无法计算'
  } catch (error) {
    console.error('计算下次执行时间失败:', error)
    nextExecutionTime.value = '计算失败'
  }
}

const handleAddTemplate = async () => {
  if (!templates.value.length) {
    await fetchTemplates()
  }
  templateDraft.value = [...selectedTemplates.value]
  showTemplateDialog.value = true
}

const applyTemplateSelection = () => {
  selectedTemplates.value = [...templateDraft.value]
  currentTask.shell = selectedTemplates.value.join(',')
  showTemplateDialog.value = false
  taskFormRef.value?.clearValidate?.(['shell'])
}

const removeTemplate = templateId => {
  selectedTemplates.value = selectedTemplates.value.filter(id => id !== templateId)
  templateDraft.value = templateDraft.value.filter(id => id !== templateId)
  currentTask.shell = selectedTemplates.value.join(',')
  taskFormRef.value?.clearValidate?.(['shell'])
}

const handleAddHost = () => {
  showHostDialog.value = true
}

const handleHostSelected = hosts => {
  const normalizedHosts = Array.isArray(hosts) ? hosts.map(normalizeHost).filter(item => item.id) : []
  const existingIds = new Set(selectedHosts.value.map(item => item.id))
  selectedHosts.value = [
    ...selectedHosts.value,
    ...normalizedHosts.filter(item => !existingIds.has(item.id))
  ]
  currentTask.hostIDs = selectedHosts.value.map(item => item.id).join(',')
  taskFormRef.value?.clearValidate?.(['hostIDs'])
}

const removeHost = hostId => {
  selectedHosts.value = selectedHosts.value.filter(item => item.id !== hostId)
  currentTask.hostIDs = selectedHosts.value.map(item => item.id).join(',')
  taskFormRef.value?.clearValidate?.(['hostIDs'])
}

const handleSubmit = async () => {
  try {
    await taskFormRef.value?.validate()
    const payload = {
      name: currentTask.name.trim(),
      type: Number(currentTask.type),
      shell: selectedTemplates.value.join(','),
      host_ids: selectedHosts.value.map(item => item.id).join(','),
      remark: currentTask.remark.trim()
    }

    if (currentTask.type === 2) {
      payload.cron_expr = currentTask.cronExpr.trim()
    }

    const response = await CreateTask(payload)
    if (!isSuccessResponse(response)) {
      throw new Error(response?.data?.message || '创建任务失败')
    }

    ElMessage.success('任务创建成功')
    formVisible.value = false
    loadTasks()
  } catch (error) {
    if (error?.message) {
      ElMessage.error(error?.response?.data?.message || error?.message || '创建任务失败')
    }
  }
}

const formDialogClosed = () => {
  taskFormRef.value?.clearValidate()
  resetCurrentTask()
}

const handleDelete = async id => {
  try {
    await ElMessageBox.confirm('确定要删除该任务吗？', '删除任务', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    const response = await DeleteTask({ id })
    if (!isSuccessResponse(response)) {
      throw new Error(response?.data?.message || '删除任务失败')
    }
    ElMessage.success('删除成功')
    loadTasks()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(error?.response?.data?.message || error?.message || '删除任务失败')
    }
  }
}

const handleOpenFlow = async row => {
  if (!row?.id) {
    ElMessage.error('任务ID无效')
    return
  }
  try {
    await jobFlowRef.value?.showFlow(row.id)
  } catch (error) {
    ElMessage.error(error?.message || '打开任务流程失败')
  }
}

const handleToggleTask = async row => {
  const nextState = row.isActive
  const action = nextState ? '启动' : '暂停'
  const message = nextState
    ? '确定要启动该定时任务吗？启动后任务会恢复自动执行。'
    : '确定要暂停该定时任务吗？暂停后任务不会继续自动执行。'

  try {
    await ElMessageBox.confirm(message, `${action}任务`, {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: nextState ? 'info' : 'warning'
    })
    const response = nextState ? await ResumeScheduledTask(row.id) : await PauseScheduledTask(row.id)
    if (!isSuccessResponse(response)) {
      throw new Error(response?.data?.message || `${action}任务失败`)
    }
    row.status = nextState ? 2 : 5
    ElMessage.success(`任务${action}成功`)
    loadTasks()
  } catch (error) {
    row.isActive = !row.isActive
    if (error !== 'cancel') {
      ElMessage.error(error?.response?.data?.message || error?.message || `${action}任务失败`)
    }
  }
}

watch(
  () => currentTask.type,
  value => {
    if (Number(value) !== 2) {
      currentTask.cronExpr = ''
      nextExecutionTime.value = ''
      taskFormRef.value?.clearValidate?.(['cronExpr'])
    }
  }
)

onMounted(async () => {
  await Promise.all([loadTasks(), fetchTemplates()])
})
</script>

<style scoped>
.task-job-page :deep(.page-actions) {
  align-self: flex-start;
}

.task-name-link {
  font-weight: 700;
}

.task-job-form {
  display: grid;
  gap: 18px;
}

.dialog-grid {
  display: grid;
  gap: 16px;
}

.dialog-grid--two {
  grid-template-columns: repeat(2, minmax(0, 1fr));
}

.schedule-note {
  padding: 14px 16px;
  border-radius: 18px;
  border: 1px solid var(--border-subtle);
  background: rgba(255, 255, 255, 0.04);
  display: grid;
  gap: 4px;
}

.schedule-note__label {
  font-size: 12px;
  font-weight: 700;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  color: rgba(125, 211, 252, 0.88);
}

.schedule-note__text {
  color: var(--text-secondary);
  line-height: 1.6;
}

.selection-surface {
  padding: 16px 18px;
  display: grid;
  gap: 14px;
}

.selection-surface__header {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  align-items: flex-start;
}

.selection-surface__title {
  font-size: 16px;
  font-weight: 700;
  color: var(--text-primary);
}

.selection-surface__subtitle {
  margin-top: 4px;
  font-size: 12px;
  line-height: 1.6;
  color: var(--text-muted);
}

.selection-form-item {
  margin-bottom: 0;
}

.selection-list {
  width: 100%;
}

.selection-tag-list {
  display: grid;
  gap: 10px;
}

.selection-chip-card {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  align-items: center;
  padding: 12px 14px;
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid var(--border-subtle);
}

.selection-chip-card__title {
  color: var(--text-primary);
  font-weight: 600;
  line-height: 1.5;
}

.selection-chip-card__meta {
  margin-top: 4px;
  color: var(--text-muted);
  font-size: 12px;
}

.selection-empty {
  padding: 18px 16px;
  border: 1px dashed var(--border-medium);
  border-radius: 16px;
  color: var(--text-muted);
  background: rgba(255, 255, 255, 0.03);
}

.picker-shell {
  display: grid;
  gap: 14px;
}

.picker-note {
  color: var(--text-muted);
  font-size: 13px;
  line-height: 1.6;
}

@media (max-width: 900px) {
  .dialog-grid--two {
    grid-template-columns: 1fr;
  }

  .selection-surface__header {
    flex-direction: column;
  }
}
</style>
