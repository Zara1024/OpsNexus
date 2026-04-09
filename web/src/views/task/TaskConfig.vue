<template>
  <TablePage
    class="ansible-config-management task-config-page"
    section-title="配置列表"
    section-subtitle="围绕 Inventory、变量与 CLI 参数建立统一的 Ansible 配置中心。"
    wide
  >
    <template #header>
      <PageHeader
        eyebrow="Task Center"
        title="Ansible 配置中心"
        subtitle="按配置类型统一维护执行上下文，减少任务创建时重复输入和跨页面切换。"
      >
        <template #actions>
          <el-button :loading="loading" @click="handleRefresh">
            <el-icon><Refresh /></el-icon>
            刷新列表
          </el-button>
          <el-button type="primary" @click="handleCreate">
            <el-icon><Plus /></el-icon>
            新增配置
          </el-button>
        </template>
        <template #meta>
          <div class="page-chip-row">
            <span class="platform-chip">当前分类 {{ currentTypeOption.label }}</span>
            <span class="platform-chip">配置总数 {{ total }}</span>
            <span class="platform-chip">关键字 {{ searchKeyword || '全部' }}</span>
          </div>
        </template>
        <template #intro>
          <PageIntro
            title="使用建议"
            text="先在对应分类下沉淀通用配置，再在任务中按需引用；Inventory 建议维护主机清单，变量类配置保持结构化格式。"
          />
        </template>
      </PageHeader>
    </template>

    <template #toolbar>
      <PageToolbar>
        <el-form :inline="true" class="filter-cluster">
          <el-form-item class="filter-field" label="搜索关键字">
            <el-input
              v-model="searchKeyword"
              placeholder="请输入配置名称"
              clearable
              @clear="handleSearch"
              @keyup.enter="handleSearch"
            />
          </el-form-item>
        </el-form>
        <template #actions>
          <el-button type="primary" @click="handleSearch">
            <el-icon><Search /></el-icon>
            搜索
          </el-button>
          <el-button @click="resetSearch">
            <el-icon><Refresh /></el-icon>
            重置
          </el-button>
        </template>
      </PageToolbar>
    </template>

    <template #stats>
      <StatStrip :items="statItems" />
    </template>

    <div class="config-tab-grid">
      <button
        v-for="tab in tabOptions"
        :key="tab.name"
        type="button"
        :class="['config-tab-card', { 'is-active': activeTab === tab.name }]"
        @click="handleTabChange(tab.name)"
      >
        <span class="config-tab-card__icon">
          <el-icon><component :is="tab.icon" /></el-icon>
        </span>
        <span class="config-tab-card__label">{{ tab.label }}</span>
        <span class="config-tab-card__hint">{{ tab.hint }}</span>
      </button>
    </div>

    <el-table
      v-loading="loading"
      :data="configList"
      stripe
      class="task-config-table"
      empty-text="当前没有匹配的 Ansible 配置"
    >
      <el-table-column label="配置名称" min-width="180" show-overflow-tooltip>
        <template #default="{ row }">
          <el-button type="primary" link class="config-name-link" @click="handleView(row)">
            {{ row.name || '未命名配置' }}
          </el-button>
        </template>
      </el-table-column>
      <el-table-column label="类型" width="140">
        <template #default="{ row }">
          <el-tag :type="typeTagType(row.type)">{{ getTypeName(row.type) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="内容预览" min-width="260" show-overflow-tooltip>
        <template #default="{ row }">
          <pre class="config-preview">{{ formatContent(row.content) }}</pre>
        </template>
      </el-table-column>
      <el-table-column label="创建时间" min-width="180">
        <template #default="{ row }">
          {{ formatDate(row.createdAt) }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="220" fixed="right">
        <template #default="{ row }">
          <div class="row-actions">
            <el-button type="primary" link @click="handleView(row)">查看</el-button>
            <el-button type="warning" link @click="handleEdit(row)">编辑</el-button>
            <el-button type="danger" link @click="handleDelete(row)">删除</el-button>
          </div>
        </template>
      </el-table-column>
    </el-table>

    <template #footer>
      <el-pagination
        :current-page="currentPage"
        :page-sizes="[10, 20, 50, 100]"
        :page-size="pageSize"
        layout="total, sizes, prev, pager, next, jumper"
        :total="total"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </template>
  </TablePage>

  <el-dialog
    v-model="dialogVisible"
    class="modern-dialog task-config-form-dialog"
    :title="dialogType === 'create' ? '新增配置' : '编辑配置'"
    width="820px"
    destroy-on-close
    :close-on-click-modal="false"
    @closed="handleDialogClosed"
  >
    <el-form ref="configFormRef" :model="currentConfig" :rules="formRules" label-width="104px" class="task-config-form">
      <div class="editor-hint">
        <div class="editor-hint__label">当前分类</div>
        <div class="editor-hint__text">{{ getTypeName(currentConfig.type) }}：{{ getContentTypeTip(currentConfig.type) }}</div>
      </div>

      <div class="dialog-grid dialog-grid--two">
        <el-form-item label="配置名称" prop="name">
          <el-input v-model="currentConfig.name" placeholder="请输入配置名称" />
        </el-form-item>
        <el-form-item label="配置类型">
          <el-input :model-value="getTypeName(currentConfig.type)" disabled />
        </el-form-item>
      </div>

      <el-form-item label="配置内容" prop="content">
        <el-input
          v-model="currentConfig.content"
          class="config-editor"
          type="textarea"
          :rows="14"
          placeholder="请输入配置内容"
        />
      </el-form-item>

      <el-form-item label="备注信息">
        <el-input
          v-model="currentConfig.remark"
          type="textarea"
          :rows="3"
          maxlength="300"
          show-word-limit
          placeholder="请输入备注信息（可选）"
        />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="dialogVisible = false">取消</el-button>
      <el-button type="primary" :loading="submitting" @click="handleSubmit">保存配置</el-button>
    </template>
  </el-dialog>

  <el-dialog
    v-model="viewDialogVisible"
    class="modern-dialog task-config-detail-dialog"
    title="配置详情"
    width="860px"
    destroy-on-close
    @closed="handleViewDialogClosed"
  >
    <div v-if="viewConfig" class="detail-shell">
      <div class="detail-grid">
        <div class="detail-card">
          <div class="detail-card__label">配置名称</div>
          <div class="detail-card__value">{{ viewConfig.name }}</div>
        </div>
        <div class="detail-card">
          <div class="detail-card__label">配置类型</div>
          <div class="detail-card__value">{{ getTypeName(viewConfig.type) }}</div>
        </div>
        <div class="detail-card">
          <div class="detail-card__label">创建时间</div>
          <div class="detail-card__value">{{ formatDate(viewConfig.createdAt) }}</div>
        </div>
        <div class="detail-card">
          <div class="detail-card__label">备注信息</div>
          <div class="detail-card__value">{{ viewConfig.remark || '暂无备注' }}</div>
        </div>
      </div>
      <pre class="detail-code">{{ viewConfig.content }}</pre>
    </div>
    <template #footer>
      <el-button @click="viewDialogVisible = false">关闭</el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  List,
  Operation,
  Plus,
  Refresh,
  Search,
  Tools,
  TopRight
} from '@element-plus/icons-vue'
import PageHeader from '@/components/platform/PageHeader.vue'
import PageIntro from '@/components/platform/PageIntro.vue'
import PageToolbar from '@/components/platform/PageToolbar.vue'
import StatStrip from '@/components/platform/StatStrip.vue'
import TablePage from '@/components/platform/TablePage.vue'
import {
  CreateAnsibleConfig,
  DeleteAnsibleConfig,
  GetAnsibleConfigById,
  GetAnsibleConfigList,
  UpdateAnsibleConfig
} from '@/api/task'

const CONFIG_TYPES = {
  INVENTORY: 1,
  GLOBAL_VARS: 2,
  EXTRA_VARS: 3,
  CLI_ARGS: 4
}

const tabOptions = [
  { name: 'inventory', type: CONFIG_TYPES.INVENTORY, label: '主机清单', hint: 'Inventory', icon: List },
  { name: 'globalVars', type: CONFIG_TYPES.GLOBAL_VARS, label: '全局变量', hint: 'Global Vars', icon: Operation },
  { name: 'extraVars', type: CONFIG_TYPES.EXTRA_VARS, label: '扩展变量', hint: 'Extra Vars', icon: TopRight },
  { name: 'cliArgs', type: CONFIG_TYPES.CLI_ARGS, label: '命令行参数', hint: 'CLI Args', icon: Tools }
]

const createEmptyConfig = () => ({
  id: null,
  name: '',
  content: '',
  remark: '',
  type: CONFIG_TYPES.INVENTORY,
  createdAt: ''
})

const loading = ref(false)
const submitting = ref(false)
const activeTab = ref('inventory')
const searchKeyword = ref('')
const configList = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)
const dialogVisible = ref(false)
const viewDialogVisible = ref(false)
const dialogType = ref('create')
const viewConfig = ref(null)
const configFormRef = ref(null)
const currentConfig = reactive(createEmptyConfig())

const formRules = reactive({
  name: [{ required: true, message: '请输入配置名称', trigger: 'blur' }],
  content: [{ required: true, message: '请输入配置内容', trigger: 'blur' }]
})

const currentTypeOption = computed(() => tabOptions.find(item => item.name === activeTab.value) || tabOptions[0])
const currentType = computed(() => currentTypeOption.value.type)

const statItems = computed(() => {
  const withRemarkCount = configList.value.filter(item => Boolean(item.remark)).length
  const multiLineCount = configList.value.filter(item => String(item.content || '').includes('\n')).length
  return [
    { label: '分类总数', value: total.value, hint: currentTypeOption.value.label, tone: 'primary' },
    { label: '本页条目', value: configList.value.length, hint: '当前分页结果', tone: 'success' },
    { label: '已写备注', value: withRemarkCount, hint: '便于后续复用', tone: 'warning' },
    { label: '多行内容', value: multiLineCount, hint: '结构化配置条目', tone: 'danger' }
  ]
})

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

const normalizeConfig = item => ({
  id: item?.id ?? item?.ID ?? null,
  name: item?.name ?? item?.Name ?? '',
  content: item?.content ?? item?.Content ?? '',
  remark: item?.remark ?? item?.Remark ?? '',
  createdAt: item?.createdAt ?? item?.CreatedAt ?? item?.created_at ?? '',
  type: Number(item?.type ?? item?.Type ?? currentType.value)
})

const getTypeName = type => {
  switch (Number(type)) {
    case CONFIG_TYPES.INVENTORY:
      return '主机清单'
    case CONFIG_TYPES.GLOBAL_VARS:
      return '全局变量'
    case CONFIG_TYPES.EXTRA_VARS:
      return '扩展变量'
    case CONFIG_TYPES.CLI_ARGS:
      return '命令行参数'
    default:
      return '未知类型'
  }
}

const typeTagType = type => ({
  [CONFIG_TYPES.INVENTORY]: 'success',
  [CONFIG_TYPES.GLOBAL_VARS]: 'info',
  [CONFIG_TYPES.EXTRA_VARS]: 'warning',
  [CONFIG_TYPES.CLI_ARGS]: 'danger'
}[Number(type)] || 'info')

const getContentTypeTip = type => {
  switch (Number(type)) {
    case CONFIG_TYPES.INVENTORY:
      return '支持 INI、JSON 或 YAML 格式的 Inventory 内容。'
    case CONFIG_TYPES.GLOBAL_VARS:
      return '建议使用 JSON、YAML 或 key=value 形式维护全局变量。'
    case CONFIG_TYPES.EXTRA_VARS:
      return '额外变量优先级更高，建议保持结构化格式便于覆盖。'
    case CONFIG_TYPES.CLI_ARGS:
      return '请输入合法的 ansible-playbook 命令行参数。'
    default:
      return ''
  }
}

const formatDate = value => {
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

const formatContent = value => {
  const content = String(value || '').replace(/\r/g, '')
  if (!content.trim()) return '暂无内容'
  const firstLine = content.split('\n').find(line => line.trim()) || content.trim()
  return firstLine.length > 110 ? `${firstLine.slice(0, 110)}...` : firstLine
}

const loadConfigList = async () => {
  loading.value = true
  try {
    const response = await GetAnsibleConfigList({
      page: currentPage.value,
      size: pageSize.value,
      type: currentType.value,
      name: searchKeyword.value.trim() || undefined
    })
    const { list, total: count } = parseCollection(response)
    configList.value = list.map(normalizeConfig)
    total.value = count
  } catch (error) {
    console.error('获取配置列表失败:', error)
    configList.value = []
    total.value = 0
    ElMessage.error(error?.response?.data?.message || error?.message || '获取配置列表失败')
  } finally {
    loading.value = false
  }
}

const resetCurrentConfig = () => {
  Object.assign(currentConfig, createEmptyConfig(), { type: currentType.value })
}

const handleTabChange = tabName => {
  if (activeTab.value === tabName) return
  activeTab.value = tabName
  currentPage.value = 1
  searchKeyword.value = ''
  loadConfigList()
}

const handleSearch = () => {
  currentPage.value = 1
  loadConfigList()
}

const resetSearch = () => {
  searchKeyword.value = ''
  currentPage.value = 1
  loadConfigList()
}

const handleRefresh = () => {
  loadConfigList()
}

const handleSizeChange = value => {
  pageSize.value = value
  currentPage.value = 1
  loadConfigList()
}

const handleCurrentChange = value => {
  currentPage.value = value
  loadConfigList()
}

const handleCreate = () => {
  dialogType.value = 'create'
  resetCurrentConfig()
  dialogVisible.value = true
}

const fillConfigDetail = async id => {
  const response = await GetAnsibleConfigById(id)
  if (!isSuccessResponse(response)) {
    throw new Error(response?.data?.message || '获取配置详情失败')
  }
  return normalizeConfig(response?.data?.data || {})
}

const handleEdit = async row => {
  dialogType.value = 'edit'
  Object.assign(currentConfig, normalizeConfig(row))
  try {
    const detail = await fillConfigDetail(currentConfig.id)
    Object.assign(currentConfig, detail)
  } catch (error) {
    console.error('获取配置详情失败:', error)
  }
  dialogVisible.value = true
}

const handleView = async row => {
  viewConfig.value = normalizeConfig(row)
  try {
    const detail = await fillConfigDetail(viewConfig.value.id)
    viewConfig.value = detail
  } catch (error) {
    console.error('获取配置详情失败:', error)
  }
  viewDialogVisible.value = true
}

const handleDelete = async row => {
  try {
    await ElMessageBox.confirm(`确定要删除配置“${row.name || row.Name || '未命名配置'}”吗？`, '删除配置', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    const response = await DeleteAnsibleConfig(row.id || row.ID)
    if (!isSuccessResponse(response)) {
      throw new Error(response?.data?.message || '删除配置失败')
    }
    ElMessage.success('删除成功')
    loadConfigList()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(error?.response?.data?.message || error?.message || '删除配置失败')
    }
  }
}

const handleSubmit = async () => {
  try {
    await configFormRef.value?.validate()
    submitting.value = true
    const payload = {
      name: currentConfig.name.trim(),
      content: currentConfig.content,
      remark: currentConfig.remark.trim(),
      type: Number(currentConfig.type)
    }
    const response = dialogType.value === 'create'
      ? await CreateAnsibleConfig(payload)
      : await UpdateAnsibleConfig({ ...payload, id: currentConfig.id })

    if (!isSuccessResponse(response)) {
      throw new Error(response?.data?.message || '保存配置失败')
    }

    ElMessage.success(dialogType.value === 'create' ? '配置创建成功' : '配置更新成功')
    dialogVisible.value = false
    loadConfigList()
  } catch (error) {
    if (error?.message) {
      ElMessage.error(error?.response?.data?.message || error?.message || '保存配置失败')
    }
  } finally {
    submitting.value = false
  }
}

const handleDialogClosed = () => {
  configFormRef.value?.clearValidate()
  resetCurrentConfig()
}

const handleViewDialogClosed = () => {
  viewConfig.value = null
}

onMounted(() => {
  resetCurrentConfig()
  loadConfigList()
})
</script>

<style scoped>
.task-config-page :deep(.page-actions) {
  align-self: flex-start;
}

.config-tab-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  gap: 12px;
  margin-bottom: 18px;
}

.config-tab-card {
  appearance: none;
  border: 1px solid var(--border-subtle);
  border-radius: 18px;
  background: rgba(255, 255, 255, 0.04);
  color: var(--text-secondary);
  padding: 16px 18px;
  display: grid;
  gap: 8px;
  text-align: left;
  cursor: pointer;
  transition: transform var(--transition-fast), border-color var(--transition-fast), background var(--transition-fast);
}

.config-tab-card:hover {
  transform: translateY(-1px);
  border-color: var(--border-medium);
}

.config-tab-card.is-active {
  border-color: var(--border-strong);
  background: linear-gradient(180deg, rgba(18, 47, 91, 0.9), rgba(15, 23, 42, 0.95));
  box-shadow: 0 18px 36px rgba(37, 99, 235, 0.16);
}

.config-tab-card__icon {
  width: 38px;
  height: 38px;
  border-radius: 12px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  background: rgba(59, 130, 246, 0.16);
  color: rgba(191, 219, 254, 0.94);
  font-size: 18px;
}

.config-tab-card__label {
  color: var(--text-primary);
  font-size: 15px;
  font-weight: 700;
}

.config-tab-card__hint {
  color: var(--text-muted);
  font-size: 12px;
}

.config-name-link {
  font-weight: 700;
}

.config-preview {
  margin: 0;
  color: var(--text-muted);
  font-family: 'Consolas', 'Courier New', monospace;
  font-size: 12px;
  line-height: 1.6;
  white-space: pre-wrap;
}

.task-config-form {
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

.editor-hint {
  padding: 14px 16px;
  border-radius: 18px;
  border: 1px solid var(--border-subtle);
  background: rgba(255, 255, 255, 0.04);
  display: grid;
  gap: 4px;
}

.editor-hint__label {
  font-size: 12px;
  font-weight: 700;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  color: rgba(125, 211, 252, 0.88);
}

.editor-hint__text {
  color: var(--text-secondary);
  line-height: 1.6;
}

.config-editor :deep(.el-textarea__inner) {
  font-family: 'Consolas', 'Courier New', monospace;
  line-height: 1.65;
}

.detail-shell {
  display: grid;
  gap: 18px;
}

.detail-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 14px;
}

.detail-card {
  padding: 16px 18px;
  border-radius: 18px;
  border: 1px solid var(--border-subtle);
  background: rgba(255, 255, 255, 0.04);
  display: grid;
  gap: 6px;
}

.detail-card__label {
  font-size: 12px;
  font-weight: 600;
  color: var(--text-muted);
}

.detail-card__value {
  color: var(--text-primary);
  font-size: 15px;
  line-height: 1.6;
  word-break: break-word;
}

.detail-code {
  margin: 0;
  min-height: 320px;
  max-height: 60vh;
  overflow: auto;
  padding: 18px 20px;
  border-radius: 18px;
  border: 1px solid var(--border-subtle);
  background: rgba(4, 11, 23, 0.94);
  color: var(--text-secondary);
  font-family: 'Consolas', 'Courier New', monospace;
  font-size: 13px;
  line-height: 1.7;
  white-space: pre-wrap;
}

@media (max-width: 900px) {
  .dialog-grid--two,
  .detail-grid {
    grid-template-columns: 1fr;
  }
}
</style>
