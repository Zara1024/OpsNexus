<template>
  <TablePage
    class="knowledge-table-page"
    section-title="知识条目"
    section-subtitle="统一沉淀 SOP、FAQ、巡检和故障复盘，并支持快速联动 AI 诊断。"
    wide
  >
    <template #header>
      <PageHeader
        eyebrow="Knowledge Hub"
        title="知识库"
        subtitle="沉淀运维 SOP、FAQ、故障手册、巡检报告和变更经验，作为 AI 与工单协作的知识底座。"
      >
        <template #actions>
          <el-button icon="Refresh" @click="bootstrapArticles">初始化示例知识</el-button>
          <el-button type="primary" icon="Plus" @click="openCreateDialog">新建文章</el-button>
        </template>
        <template #intro>
          <PageIntro
            title="维护建议"
            text="优先录入高频故障和可复用 SOP，标题尽量可搜索，标签尽量体现系统、场景和处理动作。"
          />
        </template>
      </PageHeader>
    </template>

    <template #toolbar>
      <PageToolbar>
        <el-form :inline="true" class="filter-cluster">
          <el-form-item class="filter-field" label="关键字">
            <el-input v-model="queryParams.keyword" placeholder="搜索标题、内容、标签" clearable @keyup.enter="handleQuery" />
          </el-form-item>
          <el-form-item class="filter-field" label="类型">
            <el-select v-model="queryParams.type" placeholder="全部" clearable>
              <el-option label="runbook" value="runbook" />
              <el-option label="faq" value="faq" />
              <el-option label="incident" value="incident" />
              <el-option label="change" value="change" />
              <el-option label="inspection" value="inspection" />
            </el-select>
          </el-form-item>
          <el-form-item class="filter-field" label="分类">
            <el-input v-model="queryParams.category" placeholder="例如：K8s / 数据库" clearable />
          </el-form-item>
          <el-form-item class="filter-field" label="状态">
            <el-select v-model="queryParams.enabled" placeholder="全部" clearable>
              <el-option label="启用" :value="1" />
              <el-option label="停用" :value="0" />
            </el-select>
          </el-form-item>
        </el-form>
        <template #actions>
          <el-button type="primary" icon="Search" @click="handleQuery">搜索</el-button>
          <el-button icon="Refresh" @click="resetQuery">重置</el-button>
        </template>
      </PageToolbar>
    </template>

    <template #stats>
      <StatStrip :items="statItems" />
    </template>

    <div v-if="articleList.length">
      <el-table v-loading="loading" :data="articleList" stripe class="data-table" empty-text="暂无知识文章">
        <el-table-column label="标题" min-width="240">
          <template #default="{ row }">
            <div class="knowledge-title-cell">
              <el-tooltip :content="row.title || '-'" placement="top" :show-after="300">
                <span class="ops-cell-ellipsis ops-cell-ellipsis--strong">{{ formatCompactLabel(row.title, 18, 6) }}</span>
              </el-tooltip>
              <span class="ops-cell-ellipsis knowledge-title-cell__meta">{{ row.category || '未分类' }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="类型" width="102">
          <template #default="{ row }">
            <el-tag :type="typeTagType(row.type)">{{ row.type || '-' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="标签" min-width="160">
          <template #default="{ row }">
            <el-tooltip :content="row.tags || '-'" placement="top">
              <span class="ops-cell-ellipsis">{{ row.tags || '-' }}</span>
            </el-tooltip>
          </template>
        </el-table-column>
        <el-table-column label="评分" prop="score" width="78" />
        <el-table-column label="使用" prop="useCount" width="70" />
        <el-table-column label="状态" width="84">
          <template #default="{ row }">
            <el-tag :type="row.enabled === 1 ? 'success' : 'info'">{{ row.enabled === 1 ? '启用' : '停用' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="更新时间" prop="updateTime" min-width="128" />
        <el-table-column label="操作" width="184" fixed="right">
          <template #default="{ row }">
            <div class="row-actions">
              <el-button type="success" link @click="openAIDiagnosis(row)">AI 诊断</el-button>
              <el-button type="primary" link @click="showDetail(row)">详情</el-button>
              <el-button type="warning" link @click="openEditDialog(row)">编辑</el-button>
              <el-button type="danger" link @click="deleteArticle(row)">删除</el-button>
            </div>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <EmptyState
      v-else-if="!loading"
      title="当前知识库为空"
      description="可以先初始化内置 SOP 模板，再逐步沉淀团队经验，并把高频知识联动到 AI 与工单流程中。"
    >
      <template #actions>
        <el-button @click="bootstrapArticles">一键初始化示例知识</el-button>
        <el-button type="primary" @click="openCreateDialog">新建第一篇文章</el-button>
      </template>
    </EmptyState>

    <template #footer>
      <el-pagination
        :current-page="queryParams.page"
        :page-size="queryParams.pageSize"
        :page-sizes="[10, 20, 50]"
        :total="total"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </template>

    <el-dialog
      v-model="detailVisible"
      :title="detailData.title || '文章详情'"
      width="900px"
      class="ops-dialog ops-overlay--lg knowledge-detail-dialog"
    >
      <el-descriptions :column="2" border>
        <el-descriptions-item label="类型">{{ detailData.type || '-' }}</el-descriptions-item>
        <el-descriptions-item label="分类">{{ detailData.category || '-' }}</el-descriptions-item>
        <el-descriptions-item label="标签">{{ detailData.tags || '-' }}</el-descriptions-item>
        <el-descriptions-item label="评分">{{ detailData.score || '-' }}</el-descriptions-item>
        <el-descriptions-item label="使用次数">{{ detailData.useCount || 0 }}</el-descriptions-item>
      </el-descriptions>
      <div class="detail-content markdown-body" v-html="renderMarkdown(detailData.content || '-')"></div>
    </el-dialog>

    <el-dialog
      v-model="editorVisible"
      :title="editorMode === 'create' ? '新建知识文章' : '编辑知识文章'"
      width="880px"
      class="ops-dialog ops-overlay--lg knowledge-editor-dialog"
    >
      <el-form ref="formRef" :model="formData" :rules="rules" label-width="100px" class="knowledge-editor-form">
        <el-form-item label="标题" prop="title">
          <el-input v-model="formData.title" />
        </el-form-item>
        <el-row :gutter="16">
          <el-col :span="8">
            <el-form-item label="类型" prop="type">
              <el-select v-model="formData.type" style="width: 100%">
                <el-option label="runbook" value="runbook" />
                <el-option label="faq" value="faq" />
                <el-option label="incident" value="incident" />
                <el-option label="change" value="change" />
                <el-option label="inspection" value="inspection" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="分类">
              <el-input v-model="formData.category" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="状态">
              <el-radio-group v-model="formData.enabled">
                <el-radio :label="1">启用</el-radio>
                <el-radio :label="0">停用</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="关键字">
          <el-input v-model="formData.keywords" placeholder="用逗号分隔" />
        </el-form-item>
        <el-form-item label="标签">
          <el-input v-model="formData.tags" placeholder="例如：k8s,runbook,排障" />
        </el-form-item>
        <el-form-item label="评分">
          <el-input-number v-model="formData.score" :min="0.1" :max="1" :step="0.1" />
        </el-form-item>
        <el-form-item label="内容" prop="content">
          <el-input v-model="formData.content" type="textarea" :rows="12" placeholder="支持 Markdown / 纯文本内容" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editorVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="submitArticle">保存</el-button>
      </template>
    </el-dialog>
  </TablePage>
</template>

<script>
import { ElMessage, ElMessageBox } from 'element-plus'
import EmptyState from '@/components/platform/EmptyState.vue'
import PageHeader from '@/components/platform/PageHeader.vue'
import PageIntro from '@/components/platform/PageIntro.vue'
import PageToolbar from '@/components/platform/PageToolbar.vue'
import StatStrip from '@/components/platform/StatStrip.vue'
import TablePage from '@/components/platform/TablePage.vue'

const knowledgeSeeds = [
  {
    type: 'runbook',
    category: 'K8s',
    title: 'K8s 终端审计回放 SOP',
    keywords: 'terminal audit,pod,kubectl,runbook',
    tags: 'k8s,terminal,audit',
    score: 0.9,
    enabled: 1,
    content: `# K8s 终端审计回放 SOP

## 适用场景
- Pod Web terminal 异常回放
- kubectl terminal 误操作复盘
- 需要从 sys_session_recording 与 sys_command_audit 交叉核对时

## 操作步骤
1. 在终端录像审计页按会话 ID、命令关键字或风险等级筛选
2. 优先关注录像状态为可回放的会话
3. 先看命令审计，再切到录像回放核对输入输出顺序
4. 对 command_only 历史会话，只能走命令聚合复盘`
  },
  {
    type: 'change',
    category: '数据库',
    title: 'SQL 变更工单回滚检查清单',
    keywords: 'sql,rollback,change,workorder',
    tags: 'sql,rollback,change',
    score: 0.95,
    enabled: 1,
    content: `# SQL 变更工单回滚检查清单

## 提交前
- 明确变更目标库、schema、影响表
- 确认是否带 WHERE 或 LIMIT
- 评估是否需要先备份旧数据快照

## 审批时重点
- DELETE / TRUNCATE / DROP 默认按高风险处理
- UPDATE 必须确认影响范围与回滚方案
- 对 DDL 变更，优先要求结构备份或脚本回退方案`
  },
  {
    type: 'incident',
    category: '告警',
    title: 'AlertManager 告警分诊清单',
    keywords: 'alertmanager,alert,incident,triage',
    tags: 'alertmanager,alert,triage',
    score: 0.85,
    enabled: 1,
    content: `# AlertManager 告警分诊清单

## 一级判断
- 是否为重复告警
- 是否存在静默策略
- 是否影响生产链路

## 二级排查
- 先看告警中心事件详情
- 再查关联主机、K8s 节点、应用发布记录
- 若涉及终端操作，联动终端审计回放`
  },
  {
    type: 'faq',
    category: 'AI',
    title: 'AI 诊断上下文整理模板',
    keywords: 'ai,diagnosis,context,knowledge',
    tags: 'ai,diagnosis,prompt',
    score: 0.8,
    enabled: 1,
    content: `# AI 诊断上下文整理模板

## 最小上下文
- 发生时间
- 影响对象
- 当前现象
- 最近变更
- 已执行排查动作

## 推荐补充
- 告警详情
- 终端审计命令轨迹
- SQL 工单内容与执行结果
- K8s workload、pod、event 信息`
  },
  {
    type: 'inspection',
    category: '巡检',
    title: '周度巡检报告复盘模板',
    keywords: 'inspection,report,巡检,复盘',
    tags: 'inspection,runbook,ai',
    score: 0.88,
    enabled: 1,
    content: `# 周度巡检报告复盘模板

## 巡检范围
- 主机资源与可用性
- 数据库健康与连接情况
- SSL 证书剩余有效期
- 近期告警与发布变更

## 输出要求
1. 总结本周新增风险项
2. 标记需要立刻跟进的阻断问题
3. 识别可沉淀为 SOP / 知识的重复问题
4. 给出下周治理建议与责任人`
  }
]

const createEmptyForm = () => ({
  id: null,
  type: 'runbook',
  category: '',
  title: '',
  content: '',
  keywords: '',
  tags: '',
  score: 0.5,
  enabled: 1
})

export default {
  name: 'KnowledgeBase',
  components: {
    EmptyState,
    PageHeader,
    PageIntro,
    PageToolbar,
    StatStrip,
    TablePage
  },
  data() {
    return {
      loading: false,
      saving: false,
      total: 0,
      articleList: [],
      detailVisible: false,
      detailData: {},
      editorVisible: false,
      editorMode: 'create',
      formData: createEmptyForm(),
      queryParams: {
        page: 1,
        pageSize: 10,
        keyword: '',
        type: '',
        category: '',
        enabled: ''
      },
      rules: {
        title: [{ required: true, message: '请输入标题', trigger: 'blur' }],
        type: [{ required: true, message: '请选择类型', trigger: 'change' }],
        content: [{ required: true, message: '请输入内容', trigger: 'blur' }]
      }
    }
  },
  computed: {
    statItems() {
      const enabledCount = this.articleList.filter(item => item.enabled === 1).length
      const inspectionCount = this.articleList.filter(item => item.type === 'inspection').length
      const avgScore = this.articleList.length
        ? (this.articleList.reduce((sum, item) => sum + Number(item.score || 0), 0) / this.articleList.length).toFixed(2)
        : 0
      return [
        { label: '知识总量', value: this.total, hint: '当前条件下的知识文章', tone: 'primary' },
        { label: '启用文章', value: enabledCount, hint: '可被检索与联动', tone: 'success' },
        { label: '巡检类', value: inspectionCount, hint: '适合联动 AI 诊断', tone: 'warning' },
        { label: '平均评分', value: avgScore, hint: '当前页知识质量概览', tone: 'neutral' }
      ]
    }
  },
  methods: {
    formatCompactLabel(value, prefixLength = 18, suffixLength = 6) {
      const text = String(value || '')
      if (!text) return '-'
      if (text.length <= prefixLength + suffixLength) {
        return text
      }
      return `${text.slice(0, prefixLength)}...${text.slice(-suffixLength)}`
    },
    applyRouteQuery() {
      const query = this.$route?.query || {}
      if (!Object.keys(query).length) return
      if (query.keyword !== undefined) this.queryParams.keyword = String(query.keyword || '')
      if (query.type !== undefined) this.queryParams.type = String(query.type || '')
      if (query.category !== undefined) this.queryParams.category = String(query.category || '')
    },
    buildQueryParams() {
      const params = {
        page: this.queryParams.page,
        pageSize: this.queryParams.pageSize
      }
      if (this.queryParams.keyword) params.keyword = this.queryParams.keyword
      if (this.queryParams.type) params.type = this.queryParams.type
      if (this.queryParams.category) params.category = this.queryParams.category
      if (this.queryParams.enabled !== '' && this.queryParams.enabled !== null && this.queryParams.enabled !== undefined) {
        params.enabled = this.queryParams.enabled
      }
      return params
    },
    async fetchList() {
      this.loading = true
      try {
        const { data: res } = await this.$api.getKnowledgeList(this.buildQueryParams())
        if (res.code !== 200) {
          throw new Error(res.message || '获取知识库列表失败')
        }
        const payload = res.data || {}
        this.articleList = payload.list || []
        this.total = payload.total || 0
      } catch (error) {
        ElMessage.error(error.message || '获取知识库列表失败')
      } finally {
        this.loading = false
      }
    },
    handleQuery() {
      this.queryParams.page = 1
      this.fetchList()
    },
    resetQuery() {
      this.queryParams = {
        page: 1,
        pageSize: 10,
        keyword: '',
        type: '',
        category: '',
        enabled: ''
      }
      this.fetchList()
    },
    handleSizeChange(size) {
      this.queryParams.pageSize = size
      this.queryParams.page = 1
      this.fetchList()
    },
    handleCurrentChange(page) {
      this.queryParams.page = page
      this.fetchList()
    },
    async showDetail(row) {
      const { data: res } = await this.$api.getKnowledgeDetail(row.id)
      if (res.code !== 200) {
        ElMessage.error(res.message || '获取文章详情失败')
        return
      }
      this.detailData = res.data || {}
      this.detailVisible = true
    },
    async bootstrapArticles() {
      try {
        const { data: listRes } = await this.$api.getKnowledgeList({ page: 1, pageSize: 200 })
        if (listRes.code !== 200) {
          throw new Error(listRes.message || '读取知识库失败')
        }
        const existingTitles = new Set((listRes.data?.list || []).map(item => item.title))
        let inserted = 0
        for (const seed of knowledgeSeeds) {
          if (existingTitles.has(seed.title)) continue
          const { data: res } = await this.$api.createKnowledgeArticle(seed)
          if (res.code !== 200) {
            throw new Error(res.message || `初始化文章失败: ${seed.title}`)
          }
          inserted += 1
        }
        ElMessage.success(inserted > 0 ? `已初始化 ${inserted} 篇示例知识` : '示例知识已存在，无需重复初始化')
        this.fetchList()
      } catch (error) {
        ElMessage.error(error.message || '初始化知识库失败')
      }
    },
    openCreateDialog() {
      this.editorMode = 'create'
      this.formData = createEmptyForm()
      this.editorVisible = true
    },
    openEditDialog(row) {
      this.editorMode = 'edit'
      this.formData = { ...createEmptyForm(), ...row }
      this.editorVisible = true
    },
    openAIDiagnosis(row) {
      this.$router.push({
        path: '/ai/diagnosis',
        query: {
          scene: row.type === 'inspection' ? 'inspection_report' : 'knowledge_search',
          targetId: String(row.id || ''),
          keyword: row.keywords || row.title || '',
          templateName: 'incident_analysis',
          source: row.type === 'inspection' ? 'knowledge-inspection' : 'knowledge-base',
          autoRun: row.type === 'inspection' ? '1' : '0'
        }
      })
    },
    async submitArticle() {
      const valid = await this.$refs.formRef.validate().catch(() => false)
      if (!valid) return
      this.saving = true
      try {
        let res
        if (this.editorMode === 'create') {
          res = await this.$api.createKnowledgeArticle(this.formData)
        } else {
          res = await this.$api.updateKnowledgeArticle(this.formData.id, this.formData)
        }
        if (res.data.code !== 200) {
          throw new Error(res.data.message || '保存失败')
        }
        ElMessage.success(this.editorMode === 'create' ? '创建成功' : '更新成功')
        this.editorVisible = false
        this.fetchList()
      } catch (error) {
        ElMessage.error(error.message || '保存失败')
      } finally {
        this.saving = false
      }
    },
    async deleteArticle(row) {
      try {
        await ElMessageBox.confirm(`确认删除文章“${row.title}”吗？`, '提示', { type: 'warning' })
        const { data: res } = await this.$api.deleteKnowledgeArticle(row.id)
        if (res.code !== 200) {
          throw new Error(res.message || '删除失败')
        }
        ElMessage.success('删除成功')
        this.fetchList()
      } catch (error) {
        if (error !== 'cancel' && error !== 'close') {
          ElMessage.error(error.message || '删除失败')
        }
      }
    },
    typeTagType(type) {
      const map = { runbook: 'success', faq: 'primary', incident: 'danger', change: 'warning', inspection: 'info' }
      return map[type] || 'info'
    },
    renderMarkdown(content) {
      const escapeHtml = value => value
        .replace(/&/g, '&amp;')
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;')
      let html = escapeHtml(String(content || ''))
      html = html.replace(/```([\s\S]*?)```/g, (_, code) => `<pre><code>${code.trim()}</code></pre>`)
      html = html.replace(/^###\s+(.*)$/gm, '<h3>$1</h3>')
      html = html.replace(/^##\s+(.*)$/gm, '<h2>$1</h2>')
      html = html.replace(/^#\s+(.*)$/gm, '<h1>$1</h1>')
      html = html.replace(/^-\s+(.*)$/gm, '<li>$1</li>')
      html = html.replace(/(<li>.*<\/li>)/gs, '<ul>$1</ul>')
      html = html.replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>')
      html = html.replace(/`([^`]+)`/g, '<code>$1</code>')
      html = html.replace(/\n{2,}/g, '</p><p>')
      html = html.replace(/\n/g, '<br />')
      html = `<p>${html}</p>`
      html = html.replace(/<p><h/g, '<h').replace(/<\/h([1-3])><\/p>/g, '</h$1>')
      html = html.replace(/<p><ul>/g, '<ul>').replace(/<\/ul><\/p>/g, '</ul>')
      html = html.replace(/<p><pre>/g, '<pre>').replace(/<\/pre><\/p>/g, '</pre>')
      return html
    }
  },
  watch: {
    '$route.query': {
      handler() {
        this.applyRouteQuery()
        this.fetchList()
      },
      deep: true
    }
  },
  mounted() {
    this.applyRouteQuery()
    this.fetchList()
  }
}
</script>

<style scoped>
.knowledge-title-cell {
  display: grid;
  gap: 4px;
  min-width: 0;
}

.knowledge-title-cell__meta {
  font-size: 12px;
  color: var(--text-muted);
}

.knowledge-editor-form {
  display: grid;
  gap: 18px;
}

.knowledge-editor-form :deep(.el-input-number) {
  width: 140px;
}

.detail-content {
  margin-top: 20px;
  padding: 18px;
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.03);
  line-height: 1.75;
  color: var(--text-secondary);
}

.markdown-body :deep(h1),
.markdown-body :deep(h2),
.markdown-body :deep(h3) {
  color: var(--text-primary);
  margin: 12px 0;
}

.markdown-body :deep(pre) {
  padding: 14px;
  border-radius: 12px;
  background: rgba(3, 9, 18, 0.96);
  color: #e2e8f0;
  overflow: auto;
}

.markdown-body :deep(code) {
  font-family: Consolas, Monaco, monospace;
}
</style>
