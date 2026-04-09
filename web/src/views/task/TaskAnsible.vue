<template>
  <TablePage
    class="task-ansible-page"
    section-title="任务列表"
    section-subtitle="统一管理手动与自动任务，支持快速筛选、编辑配置和进入执行历史。"
    wide
  >
    <template #header>
      <PageHeader
        eyebrow="Task Center"
        title="Ansible任务"
        subtitle="围绕任务定义、执行入口和周期调度建立统一视图，减少跨页面切换带来的治理成本。"
      >
        <template #actions>
          <el-button @click="fetchTasks">
            <el-icon><Refresh /></el-icon>
            刷新列表
          </el-button>
          <el-button type="primary" v-authority="['task:ansible:create']" @click="handleCreate">
            <el-icon><Plus /></el-icon>
            新增任务
          </el-button>
        </template>
        <template #meta>
          <div class="page-chip-row">
            <span class="platform-chip">当前搜索 {{ searchModeLabel }}</span>
            <span class="platform-chip">任务总数 {{ total }}</span>
            <span class="platform-chip">周期任务 {{ recurringCount }}</span>
          </div>
        </template>
        <template #intro>
          <PageIntro
            title="使用建议"
            text="优先通过名称或类型缩小范围；新增任务时先决定是否走配置中心，再维护主机分组、变量和执行入口。"
          />
        </template>
      </PageHeader>
    </template>

    <template #toolbar>
      <PageToolbar>
        <el-form :inline="true" :model="queryParams" class="filter-cluster">
          <el-form-item class="filter-field filter-field--mode" label="搜索模式">
            <el-radio-group v-model="searchMode" @change="handleSearchModeChange">
              <el-radio-button label="all">全部</el-radio-button>
              <el-radio-button label="name">按名称</el-radio-button>
              <el-radio-button label="type">按类型</el-radio-button>
            </el-radio-group>
          </el-form-item>
          <el-form-item v-if="searchMode !== 'type'" class="filter-field" label="任务名称">
            <el-input
              v-model="queryParams.name"
              placeholder="请输入任务名称"
              clearable
              @keyup.enter="searchTasks"
            />
          </el-form-item>
          <el-form-item v-if="searchMode !== 'name'" class="filter-field" label="任务类型">
            <el-select v-model="queryParams.type" placeholder="请选择任务类型" clearable @change="searchTasks">
              <el-option label="手动任务" :value="1" />
              <el-option label="自动任务" :value="2" />
            </el-select>
          </el-form-item>
        </el-form>
        <template #actions>
          <el-button type="primary" @click="searchTasks">
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
      class="task-ansible-table"
      empty-text="当前没有匹配的 Ansible 任务"
    >
      <el-table-column label="任务名称" min-width="180" show-overflow-tooltip>
        <template #default="{ row }">
          <el-button type="primary" link class="task-name-link" @click="goToHistory(row)">
            {{ row.name }}
          </el-button>
        </template>
      </el-table-column>
      <el-table-column label="类型" width="110">
        <template #default="{ row }">
          <el-tag :type="taskTypeTagType(row.type)">
            {{ getTaskTypeName(row.type) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="状态" width="110">
        <template #default="{ row }">
          <el-tag :type="statusTagType(row.status)" effect="dark">
            {{ getStatusText(row.status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="步骤数" width="90">
        <template #default="{ row }">
          {{ row.taskCount || 0 }}
        </template>
      </el-table-column>
      <el-table-column label="执行耗时" width="100">
        <template #default="{ row }">
          {{ formatDuration(row.totalDuration) }}
        </template>
      </el-table-column>
      <el-table-column label="Cron表达式" min-width="150" show-overflow-tooltip>
        <template #default="{ row }">
          <el-tag v-if="row.cron_expr" size="small" type="info">{{ row.cron_expr }}</el-tag>
          <span v-else class="muted">-</span>
        </template>
      </el-table-column>
      <el-table-column label="定时开关" width="90" align="center">
        <template #default="{ row }">
          <el-switch
            v-model="row.is_recurring"
            :active-value="1"
            :inactive-value="0"
            inline-prompt
            active-text="开"
            inactive-text="关"
            @change="handleRecurringChange(row)"
          />
        </template>
      </el-table-column>
      <el-table-column label="更新时间" min-width="140" prop="updatedAt" />
      <el-table-column label="操作" width="180" fixed="right">
        <template #default="{ row }">
          <div class="row-actions">
            <el-button type="primary" link v-authority="['task:ansible:start']" @click="showStartTaskDialog(row)">
              启动
            </el-button>
            <el-button type="warning" link @click="handleEdit(row)">编辑</el-button>
            <el-button type="danger" link v-authority="['task:ansible:delete']" @click="handleDelete(row.id)">
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
    :title="isEdit ? '编辑 Ansible 任务' : '新增 Ansible 任务'"
    width="860px"
    class="ops-dialog ops-overlay--full task-ansible-form-dialog"
    destroy-on-close
    @closed="formDialogClosed"
  >
    <el-form ref="ansibleFormRef" :model="currentTask" :rules="taskRules" label-width="118px" class="task-form ops-dialog__body">
      <el-row :gutter="16">
        <el-col :span="12">
          <el-form-item label="任务名称" prop="name">
            <el-input v-model="currentTask.name" placeholder="请输入任务名称" />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="任务类型" prop="type">
            <el-radio-group v-model="currentTask.type">
              <el-radio :label="1">手动任务</el-radio>
              <el-radio :label="2">自动任务</el-radio>
            </el-radio-group>
          </el-form-item>
        </el-col>
        <el-col :span="24">
          <el-form-item label="任务描述">
            <el-input
              v-model="currentTask.description"
              type="textarea"
              :rows="3"
              placeholder="请输入任务描述（可选）"
              maxlength="500"
              show-word-limit
            />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="使用配置中心">
            <el-switch v-model="currentTask.use_config" :active-value="1" :inactive-value="0" />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="周期性任务">
            <el-switch v-model="currentTask.is_recurring" :active-value="1" :inactive-value="0" />
          </el-form-item>
        </el-col>
        <el-col v-if="currentTask.is_recurring === 1" :span="24">
          <el-form-item label="Cron表达式" prop="cron_expr">
            <el-input v-model="currentTask.cron_expr" placeholder="例如：0 0 * * *（每天零点执行）" />
          </el-form-item>
        </el-col>
      </el-row>

      <div v-if="currentTask.use_config === 1" class="task-form-block">
        <div class="task-form-block__title">配置中心引用</div>
        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="Inventory配置">
              <el-select v-model="currentTask.inventory_config_id" placeholder="请选择 Inventory 配置" clearable filterable style="width: 100%">
                <el-option v-for="item in configOptions.inventory" :key="item.ID" :label="item.Name" :value="item.ID" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="全局变量配置">
              <el-select v-model="currentTask.global_vars_config_id" placeholder="请选择全局变量配置" clearable filterable style="width: 100%">
                <el-option v-for="item in configOptions.globalVars" :key="item.ID" :label="item.Name" :value="item.ID" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="额外变量配置">
              <el-select v-model="currentTask.extra_vars_config_id" placeholder="请选择额外变量配置" clearable filterable style="width: 100%">
                <el-option v-for="item in configOptions.extraVars" :key="item.ID" :label="item.Name" :value="item.ID" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="命令行参数配置">
              <el-select v-model="currentTask.cli_args_config_id" placeholder="请选择命令行参数配置" clearable filterable style="width: 100%">
                <el-option v-for="item in configOptions.cliArgs" :key="item.ID" :label="item.Name" :value="item.ID" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
      </div>

      <div v-else class="task-form-block">
        <div class="task-form-block__title">主机分组与变量</div>
        <el-form-item label="主机分组" prop="hostGroups">
          <div class="host-group-editor">
            <div class="host-group-editor__toolbar">
              <div class="host-group-editor__hint">建议按业务角色划分分组，例如 web、api、db。</div>
              <el-button type="primary" plain @click="handleAddHostGroup">添加主机分组</el-button>
            </div>
            <div v-if="hostGroupEntries.length" class="host-group-list">
              <div v-for="entry in hostGroupEntries" :key="entry[0]" class="host-group-card">
                <div class="host-group-card__main">
                  <div class="host-group-card__title">{{ entry[0] }}</div>
                  <div class="host-group-card__meta">{{ entry[1].length }} 台主机</div>
                </div>
                <div class="host-group-card__actions">
                  <el-button type="primary" link @click="showGroupHosts(entry[0])">
                    {{ expandedGroups.includes(entry[0]) ? '收起主机' : '查看主机' }}
                  </el-button>
                  <el-button type="danger" link @click="removeHostGroup(entry[0])">移除</el-button>
                </div>
                <div v-if="expandedGroups.includes(entry[0])" class="host-chip-list">
                  <el-tag v-for="hostId in entry[1]" :key="`${entry[0]}-${hostId}`" size="small" type="info">
                    ID: {{ hostId }}
                  </el-tag>
                </div>
              </div>
            </div>
            <div v-else class="empty-state host-group-empty-state">
              <div class="empty-state__halo"></div>
              <div class="empty-state__title">尚未配置主机分组</div>
              <div class="empty-state__description">点击“添加主机分组”后选择主机，即可为任务生成标准分组结构。</div>
            </div>
          </div>
        </el-form-item>

        <el-row v-if="currentTask.type === 2" :gutter="16">
          <el-col :span="24">
            <el-form-item label="Git仓库地址" prop="gitRepo">
              <el-input v-model="currentTask.gitRepo" placeholder="请输入 Git 仓库地址，例如 git@gitee.com:demo/ansible-playbook.git" />
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item label="Playbook路径">
              <el-select
                v-model="currentTask.playbook_paths"
                multiple
                filterable
                allow-create
                default-first-option
                placeholder="输入 Playbook 文件路径后回车"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="全局变量">
          <el-input
            v-model="currentTask.variables"
            type="textarea"
            :rows="3"
            placeholder='请输入 JSON 格式的全局变量，例如：{"name":"mysql-test","version":"5.7"}'
          />
        </el-form-item>
        <el-form-item label="额外变量">
          <el-input
            v-model="currentTask.extra_vars"
            type="textarea"
            :rows="3"
            placeholder="请输入 JSON 或 YAML 格式的额外变量"
          />
        </el-form-item>
        <el-form-item label="命令行参数">
          <el-input v-model="currentTask.cli_args" placeholder="请输入命令行参数" />
        </el-form-item>
      </div>

      <div v-if="currentTask.type === 1" class="task-form-block">
        <div class="task-form-block__title">手动任务附件</div>
        <el-form-item label="Playbook文件" prop="playbooks">
          <el-upload
            ref="playbookUpload"
            :auto-upload="false"
            :show-file-list="true"
            :limit="1"
            accept=".yml,.yaml"
            @change="handlePlaybookChange"
          >
            <el-button type="primary" plain>选择 Playbook</el-button>
            <template #tip>
              <div class="el-upload__tip">只能上传 YAML 文件，且不超过 10MB。</div>
            </template>
          </el-upload>
        </el-form-item>
        <el-form-item label="Roles压缩包">
          <el-upload
            ref="rolesUpload"
            :auto-upload="false"
            :show-file-list="true"
            :limit="1"
            accept=".zip"
            @change="handleRolesChange"
          >
            <el-button type="primary" plain>选择 Roles</el-button>
            <template #tip>
              <div class="el-upload__tip">只能上传 ZIP 文件，且不超过 50MB。</div>
            </template>
          </el-upload>
        </el-form-item>
      </div>
    </el-form>
    <template #footer>
      <el-button @click="formVisible = false">取消</el-button>
      <el-button type="primary" :loading="submitting" @click="handleSubmit">保存任务</el-button>
    </template>
  </el-dialog>

  <el-dialog
    v-model="hostGroupDialogVisible"
    title="配置主机分组"
    width="720px"
    class="ops-dialog ops-overlay--lg task-ansible-host-dialog"
    destroy-on-close
  >
    <div class="host-dialog-shell ops-dialog__body">
      <div class="host-dialog-toolbar">
        <el-input v-model="newGroupName" placeholder="请输入分组名称，例如 web、api、db" clearable />
        <el-button type="primary" :disabled="!newGroupName.trim()" @click="showHostSelector">选择主机</el-button>
      </div>

      <div v-if="hostGroupEntries.length" class="host-group-list">
        <div v-for="entry in hostGroupEntries" :key="entry[0]" class="host-group-card">
          <div class="host-group-card__main">
            <div class="host-group-card__title">{{ entry[0] }}</div>
            <div class="host-group-card__meta">{{ entry[1].length }} 台主机</div>
          </div>
          <div class="host-group-card__actions">
            <el-button type="primary" link @click="showGroupHosts(entry[0])">
              {{ expandedGroups.includes(entry[0]) ? '收起主机' : '查看主机' }}
            </el-button>
            <el-button type="danger" link @click="removeHostGroup(entry[0])">删除分组</el-button>
          </div>
          <div v-if="expandedGroups.includes(entry[0])" class="host-chip-list">
            <el-tag v-for="hostId in entry[1]" :key="`${entry[0]}-dialog-${hostId}`" size="small" type="info">
              ID: {{ hostId }}
            </el-tag>
          </div>
        </div>
      </div>
      <div v-else class="empty-state host-group-empty-state">
        <div class="empty-state__halo"></div>
        <div class="empty-state__icon">
          <el-icon><Connection /></el-icon>
        </div>
        <div class="empty-state__title">暂无主机分组</div>
        <div class="empty-state__description">先填写分组名称，再点击“选择主机”把目标主机加入分组。</div>
      </div>
    </div>
    <template #footer>
      <el-button @click="hostGroupDialogVisible = false">关闭</el-button>
    </template>
  </el-dialog>

  <CreateTaskHost v-model="showHostDialog" @hosts-selected="handleHostSelected" />
  <AnsibleJobFlow ref="ansibleJobFlowRef" />
</template>

<script setup>
import { computed, onMounted, reactive, ref, shallowRef } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Connection, Plus, Refresh, Search } from '@element-plus/icons-vue'
import PageHeader from '@/components/platform/PageHeader.vue'
import PageIntro from '@/components/platform/PageIntro.vue'
import PageToolbar from '@/components/platform/PageToolbar.vue'
import StatStrip from '@/components/platform/StatStrip.vue'
import TablePage from '@/components/platform/TablePage.vue'
import AnsibleJobFlow from './Job/AnsibleJobFlow.vue'
import CreateTaskHost from './Job/CreateTaskHost.vue'
import {
  CreateAnsibleTask,
  DeleteAnsibleTask,
  GetAnsibleConfigList,
  GetAnsibleTaskById,
  GetAnsibleTaskList,
  GetAnsibleTasksByName,
  GetAnsibleTasksByType,
  UpdateAnsibleTask
} from '@/api/task'
import { buildRecurringTogglePayload } from '@/utils/ansibleTaskPayload.mjs'

const createEmptyTask = () => ({
  name: '',
  description: '',
  type: 1,
  hostGroups: {},
  variables: '',
  extra_vars: '',
  cli_args: '',
  gitRepo: '',
  playbookFile: null,
  rolesFile: null,
  use_config: 0,
  inventory_config_id: null,
  global_vars_config_id: null,
  extra_vars_config_id: null,
  cli_args_config_id: null,
  is_recurring: 0,
  cron_expr: '',
  playbook_paths: [],
  view_id: null
})

const router = useRouter()
const tasks = ref([])
const loading = ref(false)
const formVisible = ref(false)
const searchMode = ref('all')
const ansibleFormRef = shallowRef(null)
const total = ref(0)
const submitting = ref(false)
const isValidationFailure = (error) => Boolean(error)
  && typeof error === 'object'
  && Object.values(error).some((value) => Array.isArray(value))
const hostGroupDialogVisible = ref(false)
const showHostDialog = ref(false)
const newGroupName = ref('')
const isEdit = ref(false)
const editId = ref(null)
const ansibleJobFlowRef = shallowRef(null)
const expandedGroups = ref([])

const queryParams = reactive({
  page: 1,
  pageSize: 10,
  name: '',
  type: null
})

const currentTask = ref(createEmptyTask())

const configOptions = reactive({
  inventory: [],
  globalVars: [],
  extraVars: [],
  cliArgs: []
})

const searchModeLabel = computed(() => {
  const map = { all: '全部', name: '按名称', type: '按类型' }
  return map[searchMode.value] || '全部'
})

const recurringCount = computed(() => tasks.value.filter(item => Number(item.is_recurring) === 1).length)
const hostGroupEntries = computed(() => Object.entries(currentTask.value.hostGroups || {}))

const statItems = computed(() => {
  const rows = tasks.value
  const manualCount = rows.filter(item => Number(item.type) === 1).length
  const autoCount = rows.filter(item => Number(item.type) === 2).length
  const runningCount = rows.filter(item => Number(item.status) === 2).length

  return [
    { label: '任务总量', value: total.value, hint: '当前查询结果', tone: 'primary' },
    { label: '手动任务', value: manualCount, hint: '当前页手动执行入口', tone: 'success' },
    { label: '自动任务', value: autoCount, hint: '当前页自动任务数量', tone: 'warning' },
    { label: '运行中', value: runningCount, hint: '当前页正在执行的任务', tone: 'danger' }
  ]
})

const taskRules = reactive({
  name: [{ required: true, message: '请输入任务名称', trigger: 'blur' }],
  type: [{ required: true, message: '请选择任务类型', trigger: 'change' }],
  hostGroups: [{
    validator: (_rule, _value, callback) => {
      if (currentTask.value.use_config === 0 && Object.keys(currentTask.value.hostGroups || {}).length === 0) {
        callback(new Error('请至少配置一个主机分组'))
      } else {
        callback()
      }
    },
    trigger: 'change'
  }],
  gitRepo: [{
    validator: (_rule, _value, callback) => {
      if (currentTask.value.use_config === 0 && Number(currentTask.value.type) === 2 && !currentTask.value.gitRepo) {
        callback(new Error('请输入 Git 仓库地址'))
      } else {
        callback()
      }
    },
    trigger: 'blur'
  }],
  cron_expr: [{
    validator: (_rule, _value, callback) => {
      if (Number(currentTask.value.is_recurring) === 1 && !String(currentTask.value.cron_expr || '').trim()) {
        callback(new Error('请输入 Cron 表达式'))
      } else {
        callback()
      }
    },
    trigger: 'blur'
  }]
})

const parseObjectMaybe = value => {
  if (!value) return {}
  if (typeof value === 'object' && !Array.isArray(value)) return value
  try {
    return JSON.parse(value)
  } catch (_error) {
    return {}
  }
}

const parseArrayMaybe = value => {
  if (!value) return []
  if (Array.isArray(value)) return value
  try {
    const parsed = JSON.parse(value)
    return Array.isArray(parsed) ? parsed : []
  } catch (_error) {
    return []
  }
}

function formatTime(value) {
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
  } catch (_error) {
    return value
  }
}

const normalizeTaskRow = item => ({
  id: item.ID ?? item.id,
  name: item.Name ?? item.name,
  description: item.Description ?? item.description,
  type: Number(item.Type ?? item.type ?? 1),
  gitRepo: item.GitRepo ?? item.gitRepo ?? '',
  hostGroups: parseObjectMaybe(item.HostGroups ?? item.hostGroups),
  allHostIDs: item.AllHostIDs ?? item.allHostIDs,
  variables: item.GlobalVars ?? item.variables ?? '',
  status: Number(item.status ?? item.Status ?? 0),
  errorMsg: item.ErrorMsg ?? item.errorMsg,
  taskCount: Number(item.TaskCount ?? item.taskCount ?? 0),
  totalDuration: Number(item.TotalDuration ?? item.totalDuration ?? 0),
  is_recurring: Number(item.IsRecurring ?? item.is_recurring ?? 0),
  cron_expr: item.CronExpr ?? item.cron_expr ?? '',
  createdAt: formatTime(item.CreatedAt ?? item.createdAt),
  updatedAt: formatTime(item.UpdatedAt ?? item.updatedAt),
  works: item.Works ?? item.works ?? []
})

const fetchConfigs = async () => {
  try {
    const types = [
      { type: 1, key: 'inventory' },
      { type: 2, key: 'globalVars' },
      { type: 3, key: 'extraVars' },
      { type: 4, key: 'cliArgs' }
    ]

    for (const item of types) {
      const res = await GetAnsibleConfigList({
        page: 1,
        size: 100,
        type: item.type
      })
      if (res.data && res.data.code === 200) {
        configOptions[item.key] = res.data.data?.list || []
      }
    }
  } catch (error) {
    console.error('获取配置列表失败:', error)
  }
}

const fetchTasks = async () => {
  loading.value = true
  try {
    let response

    if (searchMode.value === 'name' && queryParams.name && queryParams.name.trim()) {
      response = await GetAnsibleTasksByName({
        name: queryParams.name,
        page: queryParams.page,
        pageSize: queryParams.pageSize
      })
    } else if (searchMode.value === 'type' && queryParams.type !== null && queryParams.type !== undefined) {
      response = await GetAnsibleTasksByType({
        type: queryParams.type,
        page: queryParams.page,
        pageSize: queryParams.pageSize
      })
    } else if (searchMode.value === 'all' && queryParams.name && queryParams.name.trim()) {
      response = await GetAnsibleTasksByName({
        name: queryParams.name,
        page: queryParams.page,
        pageSize: queryParams.pageSize
      })
    } else if (searchMode.value === 'all' && queryParams.type !== null && queryParams.type !== undefined) {
      response = await GetAnsibleTasksByType({
        type: queryParams.type,
        page: queryParams.page,
        pageSize: queryParams.pageSize
      })
    } else {
      response = await GetAnsibleTaskList({
        page: queryParams.page,
        pageSize: queryParams.pageSize
      })
    }

    const responseData = response?.data || {}
    const payload = responseData.data ?? responseData
    let taskList = []

    if (Array.isArray(payload?.data)) {
      taskList = payload.data
    } else if (Array.isArray(payload?.list)) {
      taskList = payload.list
    } else if (Array.isArray(payload)) {
      taskList = payload
    }

    tasks.value = taskList.map(normalizeTaskRow)
    total.value = Number(payload?.total ?? responseData.total ?? taskList.length ?? 0)
  } catch (error) {
    console.error('获取 Ansible 任务列表失败:', error)
    ElMessage.error(`获取任务列表失败: ${error.message || '未知错误'}`)
    tasks.value = []
    total.value = 0
  } finally {
    loading.value = false
  }
}

const handleCreate = () => {
  isEdit.value = false
  editId.value = null
  currentTask.value = createEmptyTask()
  expandedGroups.value = []
  formVisible.value = true
}

const handleEdit = async row => {
  isEdit.value = true
  editId.value = row.id

  try {
    const res = await GetAnsibleTaskById(row.id)
    if (res?.data?.code !== 200) {
      throw new Error(res?.data?.message || '获取任务详情失败')
    }

    const data = res.data.data?.task_info || {}
    currentTask.value = {
      ...createEmptyTask(),
      name: data.Name || row.name || '',
      description: data.Description || row.description || '',
      type: Number(data.Type ?? row.type ?? 1),
      hostGroups: parseObjectMaybe(data.HostGroups || row.hostGroups),
      variables: data.GlobalVars || row.variables || '',
      extra_vars: data.ExtraVars || '',
      cli_args: data.CliArgs || '',
      gitRepo: data.GitRepo || row.gitRepo || '',
      use_config: Number(data.UseConfig || 0),
      inventory_config_id: data.InventoryConfigID || null,
      global_vars_config_id: data.GlobalVarsConfigID || null,
      extra_vars_config_id: data.ExtraVarsConfigID || null,
      cli_args_config_id: data.CliArgsConfigID || null,
      is_recurring: Number(data.IsRecurring ?? row.is_recurring ?? 0),
      cron_expr: data.CronExpr || row.cron_expr || '',
      playbook_paths: parseArrayMaybe(data.PlaybookPaths),
      view_id: data.ViewId || null
    }
    expandedGroups.value = []
    formVisible.value = true
  } catch (error) {
    console.error(error)
    ElMessage.error(error.message || '获取任务详情失败')
  }
}

const handleSubmit = async () => {
  try {
    await ansibleFormRef.value?.validate()
    submitting.value = true

    const formData = new FormData()
    if (isEdit.value && editId.value) {
      formData.append('id', String(editId.value))
    }

    formData.append('name', currentTask.value.name)
    if (currentTask.value.description) {
      formData.append('description', currentTask.value.description)
    }
    formData.append('type', String(currentTask.value.type))
    formData.append('use_config', currentTask.value.use_config ? '1' : '0')

    if (currentTask.value.use_config === 1) {
      if (currentTask.value.inventory_config_id) formData.append('inventory_config_id', currentTask.value.inventory_config_id)
      if (currentTask.value.global_vars_config_id) formData.append('global_vars_config_id', currentTask.value.global_vars_config_id)
      if (currentTask.value.extra_vars_config_id) formData.append('extra_vars_config_id', currentTask.value.extra_vars_config_id)
      if (currentTask.value.cli_args_config_id) formData.append('cli_args_config_id', currentTask.value.cli_args_config_id)
    } else {
      formData.append('hostGroups', JSON.stringify(currentTask.value.hostGroups || {}))
      if (currentTask.value.variables) formData.append('variables', currentTask.value.variables)
      if (currentTask.value.extra_vars) formData.append('extra_vars', currentTask.value.extra_vars)
      if (currentTask.value.cli_args) formData.append('cli_args', currentTask.value.cli_args)
    }

    formData.append('is_recurring', currentTask.value.is_recurring ? '1' : '0')
    if (currentTask.value.is_recurring && currentTask.value.cron_expr) {
      formData.append('cron_expr', currentTask.value.cron_expr)
    }

    if (Number(currentTask.value.type) === 1) {
      if (currentTask.value.playbookFile) {
        formData.append('playbooks', currentTask.value.playbookFile)
      }
      if (currentTask.value.rolesFile) {
        formData.append('roles', currentTask.value.rolesFile)
      }
    } else {
      formData.append('gitRepo', currentTask.value.gitRepo || '')
      if (Array.isArray(currentTask.value.playbook_paths) && currentTask.value.playbook_paths.length > 0) {
        formData.append('playbook_paths', JSON.stringify(currentTask.value.playbook_paths))
      }
    }

    if (currentTask.value.view_id) {
      formData.append('view_id', currentTask.value.view_id)
    }

    let response
    if (isEdit.value) {
      let variablesData = {}
      if (typeof currentTask.value.variables === 'string' && currentTask.value.variables.trim()) {
        try {
          variablesData = JSON.parse(currentTask.value.variables)
        } catch (_error) {
          ElMessage.error('全局变量格式错误，请输入有效的 JSON 字符串')
          return
        }
      } else {
        variablesData = currentTask.value.variables || {}
      }

      const updateData = {
        id: editId.value,
        name: currentTask.value.name,
        description: currentTask.value.description || '',
        type: currentTask.value.type,
        useConfig: currentTask.value.use_config,
        isRecurring: currentTask.value.is_recurring,
        cronExpr: currentTask.value.cron_expr || '',
        viewId: currentTask.value.view_id || 0,
        inventoryConfigId: currentTask.value.inventory_config_id || 0,
        globalVarsConfigId: currentTask.value.global_vars_config_id || 0,
        extraVarsConfigId: currentTask.value.extra_vars_config_id || 0,
        cliArgsConfigId: currentTask.value.cli_args_config_id || 0,
        variables: variablesData,
        extraVars: currentTask.value.extra_vars || '',
        cliArgs: currentTask.value.cli_args || '',
        hostGroups: parseObjectMaybe(currentTask.value.hostGroups),
        gitRepo: currentTask.value.gitRepo || '',
        playbookPaths: parseArrayMaybe(currentTask.value.playbook_paths)
      }
      response = await UpdateAnsibleTask(updateData)
    } else {
      response = await CreateAnsibleTask(formData)
    }

    if (response?.data?.code === 200) {
      ElMessage.success(isEdit.value ? 'Ansible任务更新成功' : 'Ansible任务创建成功')
      formVisible.value = false
      fetchTasks()
    } else {
      throw new Error(response?.data?.message || (isEdit.value ? '更新任务失败' : '创建任务失败'))
    }
  } catch (error) {
    if (isValidationFailure(error)) {
      return
    }
    console.error('保存 Ansible 任务失败:', error)
    ElMessage.error(`保存任务失败: ${error.message || '未知错误'}`)
  } finally {
    submitting.value = false
  }
}

const handleDelete = async id => {
  try {
    await ElMessageBox.confirm(
      '确定要删除该 Ansible 任务吗？删除后将同时删除相关子任务和执行记录。',
      '删除确认',
      {
        confirmButtonText: '确定删除',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    const response = await DeleteAnsibleTask(id)
    if (response?.data?.code === 200 || response?.status === 200) {
      ElMessage.success('任务删除成功')
      fetchTasks()
    } else {
      throw new Error(response?.data?.message || '删除任务失败')
    }
  } catch (error) {
    if (error !== 'cancel' && error !== 'close') {
      console.error('删除任务失败:', error)
      ElMessage.error(`删除失败: ${error.response?.data?.message || error.message || '未知错误'}`)
    }
  }
}

const handleAddHostGroup = () => {
  newGroupName.value = ''
  hostGroupDialogVisible.value = true
}

const showHostSelector = () => {
  if (!newGroupName.value.trim()) {
    ElMessage.warning('请先输入分组名称')
    return
  }
  showHostDialog.value = true
}

const handleHostSelected = hosts => {
  if (!newGroupName.value.trim()) {
    ElMessage.warning('请先输入分组名称')
    return
  }

  if (!hosts.length) {
    ElMessage.warning('请选择至少一台主机')
    return
  }

  if (currentTask.value.hostGroups[newGroupName.value]) {
    ElMessage.warning(`分组“${newGroupName.value}”已存在，请更换分组名称`)
    return
  }

  currentTask.value.hostGroups[newGroupName.value] = hosts.map(host => host.id)
  ElMessage.success(`已添加分组“${newGroupName.value}”，包含 ${hosts.length} 台主机`)
  newGroupName.value = ''
  showHostDialog.value = false
}

const removeHostGroup = groupName => {
  delete currentTask.value.hostGroups[groupName]
  expandedGroups.value = expandedGroups.value.filter(item => item !== groupName)
}

const handleRecurringChange = async row => {
  try {
    const statusText = Number(row.is_recurring) === 1 ? '开启' : '关闭'
    await UpdateAnsibleTask(buildRecurringTogglePayload(row))
    ElMessage.success(`已${statusText}任务“${row.name}”的定时调度`)
  } catch (error) {
    row.is_recurring = Number(row.is_recurring) === 1 ? 0 : 1
    console.error('更新定时状态失败:', error)
    ElMessage.error('更新状态失败')
  }
}

const showGroupHosts = groupName => {
  const index = expandedGroups.value.indexOf(groupName)
  if (index > -1) {
    expandedGroups.value.splice(index, 1)
  } else {
    expandedGroups.value.push(groupName)
  }
}

const handlePlaybookChange = file => {
  currentTask.value.playbookFile = file.raw
}

const handleRolesChange = file => {
  currentTask.value.rolesFile = file.raw
}

const searchTasks = () => {
  queryParams.page = 1
  fetchTasks()
}

const resetQuery = () => {
  searchMode.value = 'all'
  queryParams.name = ''
  queryParams.type = null
  queryParams.page = 1
  fetchTasks()
}

const handleSearchModeChange = mode => {
  if (mode === 'name') {
    queryParams.type = null
  } else if (mode === 'type') {
    queryParams.name = ''
  }
  queryParams.page = 1
  fetchTasks()
}

const handleSizeChange = value => {
  queryParams.pageSize = value
  queryParams.page = 1
  fetchTasks()
}

const handleCurrentChange = value => {
  queryParams.page = value
  fetchTasks()
}

const formDialogClosed = () => {
  ansibleFormRef.value?.clearValidate?.()
  currentTask.value = createEmptyTask()
  expandedGroups.value = []
  newGroupName.value = ''
  isEdit.value = false
  editId.value = null
}

const statusTagType = status => {
  const map = { 1: 'info', 2: 'warning', 3: 'success', 4: 'danger' }
  return map[Number(status)] || 'info'
}

const getStatusText = status => {
  const map = { 1: '等待中', 2: '运行中', 3: '成功', 4: '失败' }
  return map[Number(status)] || '未知'
}

const taskTypeTagType = type => (Number(type) === 1 ? 'success' : 'primary')

const getTaskTypeName = type => {
  const map = { 1: '手动任务', 2: '自动任务', 3: 'K8s任务' }
  return map[Number(type)] || '未知任务'
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

const showStartTaskDialog = async task => {
  try {
    await ansibleJobFlowRef.value?.showFlow?.(task.id)
  } catch (error) {
    console.error('显示任务流程失败:', error)
    ElMessage.error(`显示任务流程失败: ${error.message || '未知错误'}`)
  }
}

const goToHistory = row => {
  router.push({
    name: 'AnsibleTaskHistory',
    query: {
      id: row.id,
      name: row.name,
      status: row.status,
      createdAt: row.createdAt,
      updatedAt: row.updatedAt,
      totalDuration: row.totalDuration,
      taskCount: row.taskCount
    }
  })
}

onMounted(() => {
  fetchTasks()
  fetchConfigs()
})
</script>

<style scoped>
.task-ansible-page :deep(.page-actions) {
  align-items: center;
}

.task-name-link {
  padding: 0;
  font-weight: 700;
}

.filter-field--mode {
  flex: 1 1 320px;
}

.task-form {
  display: grid;
  gap: 18px;
}

.task-form-block {
  padding: 18px;
  border-radius: 20px;
  border: 1px solid var(--border-subtle);
  background: rgba(255, 255, 255, 0.04);
}

.task-form-block__title {
  margin-bottom: 16px;
  font-size: 14px;
  font-weight: 700;
  color: var(--text-primary);
}

.host-group-editor,
.host-dialog-shell {
  display: grid;
  gap: 14px;
  width: 100%;
}

.host-group-editor__toolbar,
.host-dialog-toolbar {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  align-items: center;
}

.host-group-editor__toolbar .el-button,
.host-dialog-toolbar .el-button {
  flex-shrink: 0;
}

.host-group-editor__toolbar :deep(.el-input),
.host-dialog-toolbar :deep(.el-input) {
  flex: 1 1 280px;
}

.host-group-editor__hint {
  flex: 1 1 260px;
  color: var(--text-muted);
  font-size: 13px;
  line-height: 1.6;
}

.host-group-list {
  display: grid;
  gap: 12px;
}

.host-group-card {
  padding: 16px 18px;
  border-radius: 18px;
  border: 1px solid var(--border-subtle);
  background: rgba(2, 6, 23, 0.34);
  display: grid;
  gap: 12px;
}

.host-group-card__main,
.host-group-card__actions {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
}

.host-group-card__title {
  font-size: 15px;
  font-weight: 700;
  color: var(--text-primary);
}

.host-group-card__meta {
  color: var(--text-muted);
  font-size: 12px;
}

.host-chip-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.host-group-empty-state {
  min-height: 220px;
}

@media (max-width: 768px) {
  .task-ansible-page :deep(.page-actions) {
    width: 100%;
    justify-content: flex-start;
  }
}
</style>
