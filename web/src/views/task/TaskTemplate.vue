<template>
  <TablePage
    class="tasktemplate-management task-template-page"
    section-title="模板列表"
    section-subtitle="统一维护 Shell、Python 与 Ansible 脚本模板，支持快速筛选、预览和编辑。"
    wide
  >
    <template #header>
      <PageHeader
        eyebrow="Task Center"
        title="任务模板"
        subtitle="把脚本资产统一沉淀到模板中心，便于复用、审核和在任务作业中快速编排。"
      >
        <template #actions>
          <el-button :loading="loading" @click="loadTemplates">
            <el-icon><Refresh /></el-icon>
            刷新列表
          </el-button>
          <el-button type="primary" v-authority="['task:template:add']" @click="handleCreate">
            <el-icon><Plus /></el-icon>
            新增模板
          </el-button>
        </template>
        <template #meta>
          <div class="page-chip-row">
            <span class="platform-chip">当前筛选 {{ filterSummary }}</span>
            <span class="platform-chip">模板总数 {{ total }}</span>
            <span class="platform-chip">Ansible 模板 {{ ansibleCount }}</span>
          </div>
        </template>
        <template #intro>
          <PageIntro
            title="使用建议"
            text="优先给模板补齐备注和类型；Shell/Python 适合命令脚本，Ansible 模板建议维护为 YAML/Playbook 结构。"
          />
        </template>
      </PageHeader>
    </template>

    <template #toolbar>
      <PageToolbar>
        <el-form :inline="true" :model="queryParams" class="filter-cluster">
          <el-form-item class="filter-field" label="模板名称">
            <el-input
              v-model="queryParams.name"
              placeholder="请输入模板名称"
              clearable
              @keyup.enter="handleSearch"
            />
          </el-form-item>
          <el-form-item class="filter-field" label="模板类型">
            <el-select v-model="queryParams.type" placeholder="请选择模板类型" clearable>
              <el-option v-for="item in templateTypeOptions" :key="item.value" :label="item.label" :value="item.value" />
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
      :data="templates"
      stripe
      class="task-template-table"
      empty-text="当前没有匹配的任务模板"
    >
      <el-table-column label="模板名称" min-width="180" show-overflow-tooltip>
        <template #default="{ row }">
          <span class="template-name">{{ row.name || '未命名模板' }}</span>
        </template>
      </el-table-column>
      <el-table-column label="类型" width="120">
        <template #default="{ row }">
          <el-tag :type="typeTagType(row.type)">{{ getTypeName(row.type) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="内容摘要" min-width="220" show-overflow-tooltip>
        <template #default="{ row }">
          <div class="template-preview">{{ formatPreview(row.content) }}</div>
        </template>
      </el-table-column>
      <el-table-column label="创建时间" min-width="170">
        <template #default="{ row }">
          {{ formatTime(row.createdAt) }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="180" fixed="right">
        <template #default="{ row }">
          <div class="row-actions">
            <el-button type="primary" link @click="handleViewContent(row)">查看</el-button>
            <el-button type="warning" link v-authority="['task:template:edit']" @click="handleEdit(row)">编辑</el-button>
            <el-button type="danger" link v-authority="['task:template:delete']" @click="handleDelete(row.id)">删除</el-button>
          </div>
        </template>
      </el-table-column>
    </el-table>

    <template #footer>
      <el-pagination
        :current-page="queryParams.pageNum"
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
    class="modern-dialog task-template-form-dialog"
    :title="currentTemplate.id ? '编辑模板' : '新增模板'"
    width="920px"
    destroy-on-close
    @closed="formDialogClosed"
  >
    <el-form ref="templateFormRef" :model="currentTemplate" :rules="templateRules" label-width="104px" class="task-template-form">
      <div class="dialog-grid dialog-grid--two">
        <el-form-item label="模板名称" prop="name">
          <el-input v-model="currentTemplate.name" placeholder="请输入模板名称" />
        </el-form-item>
        <el-form-item label="模板类型" prop="type">
          <el-select v-model="currentTemplate.type" placeholder="请选择模板类型" style="width: 100%">
            <el-option v-for="item in templateTypeOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
      </div>

      <div class="editor-hint">
        <div class="editor-hint__label">编辑建议</div>
        <div class="editor-hint__text">{{ editorHint }}</div>
      </div>

      <el-form-item label="模板内容" prop="content" class="template-editor-field">
        <CodeEditor
          v-model="currentTemplate.content"
          class="template-editor"
          :language="getLanguage(currentTemplate.type)"
          height="360px"
          placeholder="请输入模板内容"
        />
      </el-form-item>

      <el-form-item label="备注信息">
        <el-input
          v-model="currentTemplate.remark"
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
      <el-button type="primary" @click="handleSubmit">保存模板</el-button>
    </template>
  </el-dialog>

  <el-dialog
    v-model="contentDialogVisible"
    class="modern-dialog task-template-preview-dialog"
    title="脚本内容"
    width="920px"
    destroy-on-close
  >
    <div v-if="previewTemplate" class="preview-shell">
      <div class="page-chip-row">
        <span class="platform-chip">模板名称 {{ previewTemplate.name }}</span>
        <span class="platform-chip">模板类型 {{ getTypeName(previewTemplate.type) }}</span>
        <span class="platform-chip">创建人 {{ previewTemplate.createdBy || '-' }}</span>
      </div>
      <pre class="script-preview" v-html="scriptContent"></pre>
    </div>
    <template #footer>
      <el-button @click="contentDialogVisible = false">关闭</el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Refresh, Search } from '@element-plus/icons-vue'
import CodeEditor from '@/components/CodeEditor.vue'
import PageHeader from '@/components/platform/PageHeader.vue'
import PageIntro from '@/components/platform/PageIntro.vue'
import PageToolbar from '@/components/platform/PageToolbar.vue'
import StatStrip from '@/components/platform/StatStrip.vue'
import TablePage from '@/components/platform/TablePage.vue'
import { highlight } from '@/utils/highlight'
import {
  CreateTemplate,
  DeleteTemplate,
  GetAllTemplates,
  GetTemplateContent,
  GetTemplatesByName,
  GetTemplatesByType,
  UpdateTemplate
} from '@/api/task'

const templateTypeOptions = [
  { label: 'Shell', value: 1 },
  { label: 'Python', value: 2 },
  { label: 'Ansible', value: 3 }
]

const createEmptyQuery = () => ({
  pageNum: 1,
  pageSize: 10,
  name: '',
  type: ''
})

const createEmptyTemplate = () => ({
  id: null,
  name: '',
  type: 1,
  content: '',
  remark: ''
})

const queryParams = reactive(createEmptyQuery())
const currentTemplate = reactive(createEmptyTemplate())
const templates = ref([])
const loading = ref(false)
const formVisible = ref(false)
const contentDialogVisible = ref(false)
const total = ref(0)
const scriptContent = ref('')
const previewTemplate = ref(null)
const templateFormRef = ref(null)

const templateRules = reactive({
  name: [{ required: true, message: '请输入模板名称', trigger: 'blur' }],
  type: [{ required: true, message: '请选择模板类型', trigger: 'change' }],
  content: [{ required: true, message: '请输入模板内容', trigger: 'blur' }]
})

const shellCount = computed(() => templates.value.filter(item => item.type === 1).length)
const pythonCount = computed(() => templates.value.filter(item => item.type === 2).length)
const ansibleCount = computed(() => templates.value.filter(item => item.type === 3).length)

const filterSummary = computed(() => {
  const parts = []
  if (queryParams.name.trim()) {
    parts.push(`名称：${queryParams.name.trim()}`)
  }
  if (queryParams.type !== '') {
    parts.push(`类型：${getTypeName(queryParams.type)}`)
  }
  return parts.length ? parts.join(' / ') : '全部模板'
})

const statItems = computed(() => [
  { label: '模板总数', value: total.value, hint: '当前筛选结果', tone: 'primary' },
  { label: 'Shell', value: shellCount.value, hint: '命令脚本模板', tone: 'success' },
  { label: 'Python', value: pythonCount.value, hint: '自动化脚本模板', tone: 'warning' },
  { label: 'Ansible', value: ansibleCount.value, hint: 'Playbook / YAML 模板', tone: 'danger' }
])

const editorHint = computed(() => {
  switch (Number(currentTemplate.type)) {
    case 1:
      return '适合维护 Shell 命令、批处理脚本和轻量巡检逻辑。'
    case 2:
      return '适合维护 Python 自动化脚本，可直接复用函数与标准库。'
    case 3:
      return '适合维护 Playbook、Role 入口或 YAML 结构脚本，建议保持缩进规范。'
    default:
      return '请根据脚本语言选择正确的模板类型。'
  }
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

const normalizeTemplate = item => ({
  id: item?.id ?? item?.ID ?? null,
  name: item?.name ?? item?.Name ?? '',
  type: Number(item?.type ?? item?.Type ?? 1),
  content: item?.content ?? item?.Content ?? '',
  remark: item?.remark ?? item?.Remark ?? '',
  createdBy: item?.createdBy ?? item?.CreatedBy ?? '',
  createdAt: item?.createdAt ?? item?.CreatedAt ?? item?.created_at ?? '',
  updatedAt: item?.updatedAt ?? item?.UpdatedAt ?? item?.updated_at ?? ''
})

const getTypeName = type => templateTypeOptions.find(item => item.value === Number(type))?.label || 'Shell'
const getLanguage = type => ({ 1: 'bash', 2: 'python', 3: 'yaml' }[Number(type)] || 'bash')
const typeTagType = type => ({ 1: 'success', 2: 'warning', 3: 'info' }[Number(type)] || 'info')

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

const formatPreview = value => {
  const content = String(value || '').replace(/\r/g, '')
  if (!content.trim()) return '暂无内容摘要'
  const firstLine = content.split('\n').find(line => line.trim()) || content.trim()
  return firstLine.length > 80 ? `${firstLine.slice(0, 80)}...` : firstLine
}

const loadTemplates = async () => {
  loading.value = true
  try {
    const hasName = Boolean(queryParams.name.trim())
    const hasType = queryParams.type !== ''
    let response

    if (hasName && !hasType) {
      response = await GetTemplatesByName({
        name: queryParams.name.trim(),
        pageNum: queryParams.pageNum,
        pageSize: queryParams.pageSize
      })
    } else if (!hasName && hasType) {
      response = await GetTemplatesByType({
        type: queryParams.type,
        pageNum: queryParams.pageNum,
        pageSize: queryParams.pageSize
      })
    } else {
      response = await GetAllTemplates({
        pageNum: queryParams.pageNum,
        pageSize: queryParams.pageSize,
        name: hasName ? queryParams.name.trim() : undefined,
        type: hasType ? queryParams.type : undefined
      })
    }

    const { list, total: count } = parseCollection(response)
    templates.value = list.map(normalizeTemplate)
    total.value = count
  } catch (error) {
    console.error('获取模板列表失败:', error)
    templates.value = []
    total.value = 0
    ElMessage.error(error?.response?.data?.message || error?.message || '获取模板列表失败')
  } finally {
    loading.value = false
  }
}

const resetCurrentTemplate = () => {
  Object.assign(currentTemplate, createEmptyTemplate())
}

const handleSearch = () => {
  queryParams.pageNum = 1
  loadTemplates()
}

const resetQuery = () => {
  Object.assign(queryParams, createEmptyQuery())
  loadTemplates()
}

const handleSizeChange = value => {
  queryParams.pageSize = value
  queryParams.pageNum = 1
  loadTemplates()
}

const handleCurrentChange = value => {
  queryParams.pageNum = value
  loadTemplates()
}

const handleCreate = () => {
  resetCurrentTemplate()
  formVisible.value = true
}

const handleEdit = template => {
  Object.assign(currentTemplate, normalizeTemplate(template))
  formVisible.value = true
}

const handleDelete = async id => {
  try {
    await ElMessageBox.confirm('确定要删除该模板吗？', '删除模板', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    const response = await DeleteTemplate({ id })
    if (!isSuccessResponse(response)) {
      throw new Error(response?.data?.message || '删除模板失败')
    }
    ElMessage.success('删除成功')
    loadTemplates()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(error?.response?.data?.message || error?.message || '删除模板失败')
    }
  }
}

const handleViewContent = async template => {
  try {
    const normalizedTemplate = normalizeTemplate(template)
    const response = await GetTemplateContent({ id: normalizedTemplate.id })
    const rawContent = typeof response?.data === 'string'
      ? response.data
      : JSON.stringify(response?.data ?? '', null, 2)
    previewTemplate.value = normalizedTemplate
    scriptContent.value = highlight(rawContent, getLanguage(normalizedTemplate.type))
    contentDialogVisible.value = true
  } catch (error) {
    console.error('获取模板内容失败:', error)
    previewTemplate.value = normalizeTemplate(template)
    scriptContent.value = highlight('获取脚本内容失败', 'bash')
    contentDialogVisible.value = true
  }
}

const handleSubmit = async () => {
  try {
    await templateFormRef.value?.validate()
    const payload = {
      name: currentTemplate.name.trim(),
      type: Number(currentTemplate.type),
      content: currentTemplate.content,
      remark: currentTemplate.remark.trim()
    }

    const response = currentTemplate.id
      ? await UpdateTemplate({ ...payload, id: currentTemplate.id })
      : await CreateTemplate(payload)

    if (!isSuccessResponse(response)) {
      throw new Error(response?.data?.message || '保存模板失败')
    }

    ElMessage.success(currentTemplate.id ? '模板更新成功' : '模板创建成功')
    formVisible.value = false
    loadTemplates()
  } catch (error) {
    if (error?.message) {
      ElMessage.error(error?.response?.data?.message || error?.message || '保存模板失败')
    }
  }
}

const formDialogClosed = () => {
  templateFormRef.value?.clearValidate()
  resetCurrentTemplate()
}

onMounted(() => {
  loadTemplates()
})
</script>

<style scoped>
.task-template-page :deep(.page-actions) {
  align-self: flex-start;
}

.template-name {
  color: var(--text-primary);
  font-weight: 700;
}

.template-preview {
  color: var(--text-muted);
  line-height: 1.6;
}

.task-template-form {
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

.template-editor-field {
  margin-bottom: 0;
}

.template-editor :deep(.code-editor) {
  border-color: var(--border-medium);
  border-radius: 18px;
  overflow: hidden;
}

.template-editor :deep(.code-editor__highlight) {
  background: rgba(4, 11, 23, 0.94);
}

.preview-shell {
  display: grid;
  gap: 16px;
}

.script-preview {
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

.script-preview :deep(.hljs) {
  background: transparent;
}

@media (max-width: 900px) {
  .dialog-grid--two {
    grid-template-columns: 1fr;
  }
}
</style>
