<template>
  <AIWorkspace
    eyebrow="AI Workspace"
    title="AI 诊断分析台"
    subtitle="复用 Prompt Template、知识库和历史会话，在统一工作区中生成终端、SQL、巡检、容量治理等多模型诊断分析。"
  >
    <template #actions>
      <el-button icon="Refresh" @click="refreshHistory">刷新历史</el-button>
    </template>

    <template #intro>
      <PageIntro
        title="工作方式"
        text="先选择诊断场景与目标，再补齐检索词和知识上下文，最后由 AI 智能运维助手调用模板与多模型能力生成可复盘、可落地的分析结果。"
      />
    </template>

    <template #rail>
      <SectionCard title="AI 智能运维助手" subtitle="诊断分析与知识联动的多模型工作台入口。">
        <div class="ai-diagnosis__positioning">{{ brand.assistantPositioning }}</div>
        <div class="page-chip-row">
          <span v-for="item in brand.aiModels || []" :key="item" class="platform-chip">{{ item }}</span>
        </div>
      </SectionCard>
      <SectionCard title="AI 运行状态" subtitle="展示当前模型链路、知识资产与诊断会话情况。">
        <div class="ai-diagnosis__runtime-status">
          <span :class="['ai-diagnosis__runtime-badge', `is-${runtimeStatus.status}`]">
            {{ runtimeStatus.statusText }}
          </span>
          <div class="ai-diagnosis__runtime-meta">{{ runtimeStatus.provider }} / {{ runtimeStatus.model }}</div>
          <div v-if="runtimeStatus.checkedAt" class="ai-diagnosis__runtime-meta">最近探测：{{ runtimeStatus.checkedAt }}</div>
          <div v-if="runtimeErrorDisplay" class="ai-diagnosis__runtime-meta ai-diagnosis__runtime-meta--warning">{{ runtimeErrorLabel }}：{{ runtimeErrorDisplay }}</div>
          <div class="ai-diagnosis__runtime-meta">推理强度：{{ runtimeStatus.reasoningEffort }}</div>
          <div class="ai-diagnosis__runtime-meta">知识条目：{{ overviewStats.knowledgeItems || 0 }}</div>
          <div class="ai-diagnosis__runtime-meta">Prompt 模板：{{ overviewStats.promptTemplates || templateOptions.length }}</div>
          <div class="ai-diagnosis__runtime-meta">诊断会话：{{ overviewStats.diagnosisSessions || historyList.length }}</div>
        </div>
      </SectionCard>
      <SectionCard title="领域入口" subtitle="按领域快速进入 AI 助手或诊断工作流。">
        <div class="ai-diagnosis__scene-list">
          <button
            v-for="item in workspaceDomains"
            :key="item.key"
            class="ai-diagnosis__scene"
            @click="openDomain(item)"
          >
            <div class="ai-diagnosis__scene-title">{{ item.label }}</div>
            <div class="ai-diagnosis__scene-desc">{{ item.description }}</div>
          </button>
        </div>
      </SectionCard>

      <SectionCard title="诊断场景" subtitle="按任务类型切换提示模板。">
        <div class="ai-diagnosis__scene-list">
          <button
            v-for="scene in sceneOptions"
            :key="scene.value"
            :class="['ai-diagnosis__scene', { active: formData.scene === scene.value }]"
            @click="handleSceneChange(scene.value)"
          >
            <div class="ai-diagnosis__scene-title">{{ scene.label }}</div>
            <div class="ai-diagnosis__scene-desc">{{ scene.description }}</div>
          </button>
        </div>
      </SectionCard>

      <SectionCard title="关联知识" subtitle="点击标签即可参与本次诊断。">
        <div class="page-chip-row">
          <el-check-tag
            v-for="item in knowledgeSuggestions"
            :key="item.id"
            :checked="selectedKnowledgeIds.includes(item.id)"
            @change="toggleKnowledge(item.id)"
          >
            {{ item.title }}
          </el-check-tag>
        </div>
        <EmptyState
          v-if="!knowledgeSuggestions.length"
          title="暂无推荐知识"
          description="切换场景或调整检索词后，会自动尝试召回关联知识。"
        />
      </SectionCard>

      <SectionCard title="最近会话" subtitle="用于回看诊断上下文与提示词演进。">
        <div v-if="historyList.length" class="ai-diagnosis__history">
          <button
            v-for="item in historyList.slice(0, 6)"
            :key="item.sessionId"
            class="ai-diagnosis__history-item"
            @click="showHistory(item.sessionId)"
          >
            <div class="ai-diagnosis__history-id">{{ item.sessionId }}</div>
            <div class="ai-diagnosis__history-meta">{{ item.latestTime || '-' }}</div>
          </button>
        </div>
        <EmptyState
          v-else
          title="还没有会话记录"
          description="完成一次诊断后，这里会沉淀最近的 AI 会话。"
        />
      </SectionCard>
    </template>

    <StatStrip :items="statItemsEnhanced" />

    <SectionCard title="诊断输入" subtitle="补齐目标、检索词和模板后即可生成诊断结果。">
      <el-alert
        v-if="routeContext.source === 'capacity-suggestion'"
        :title="`当前诊断来自容量建议联动：${routeContext.workloadName || '-'} / 快照 ${formData.targetId || '-'}`"
        type="info"
        :closable="false"
        class="page-alert"
      />

      <el-form :model="formData" label-width="120px" class="ai-diagnosis__form">
        <el-form-item label="诊断场景">
          <el-radio-group v-model="formData.scene" @change="handleSceneChange">
            <el-radio v-for="scene in sceneOptions" :key="scene.value" :label="scene.value">{{ scene.label }}</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item :label="currentSceneMeta.targetLabel">
          <el-input v-model="formData.targetId" :placeholder="currentSceneMeta.targetPlaceholder" />
        </el-form-item>
        <el-form-item label="知识检索词">
          <el-input v-model="formData.keyword" :placeholder="currentSceneMeta.keywordPlaceholder" />
        </el-form-item>
        <el-form-item label="模板">
          <el-select v-model="formData.templateName" style="width: 280px">
            <el-option v-for="item in templateOptions" :key="item.name" :label="item.name" :value="item.name" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="diagnosing" @click="runDiagnosis">生成分析</el-button>
          <el-button @click="loadKnowledgeSuggestions">刷新知识推荐</el-button>
        </el-form-item>
      </el-form>
    </SectionCard>

    <div class="platform-data-grid platform-data-grid--three">
      <SectionCard title="诊断报告" dense>
        <pre class="ai-diagnosis__code">{{ renderResult.report || '-' }}</pre>
      </SectionCard>
      <SectionCard title="系统提示" dense>
        <pre class="ai-diagnosis__code">{{ renderResult.systemPrompt || '-' }}</pre>
      </SectionCard>
      <SectionCard title="生成结果" dense>
        <pre class="ai-diagnosis__code">{{ renderResult.renderedPrompt || '-' }}</pre>
      </SectionCard>
    </div>

    <SectionCard title="最近 AI 会话" subtitle="可继续深入查看具体的角色消息与提示内容。">
      <el-table :data="historyList" size="small" stripe>
        <el-table-column label="会话 ID" prop="sessionId" min-width="180" />
        <el-table-column label="消息数" prop="messageCount" width="100" />
        <el-table-column label="最近时间" prop="latestTime" min-width="170" />
        <el-table-column label="操作" width="90">
          <template #default="{ row }">
            <el-button type="primary" link @click="showHistory(row.sessionId)">详情</el-button>
          </template>
        </el-table-column>
      </el-table>
    </SectionCard>

    <el-dialog v-model="historyVisible" title="AI 会话详情" width="900px" class="ops-dialog ops-overlay--lg">
      <div class="ai-diagnosis__history-detail">
        <div v-for="item in historyDetail" :key="item.id" class="ai-diagnosis__history-message">
          <div class="ai-diagnosis__history-role">{{ item.role }}</div>
          <pre class="ai-diagnosis__code">{{ item.message }}</pre>
        </div>
      </div>
    </el-dialog>
  </AIWorkspace>
</template>

<script>
import { ElMessage } from 'element-plus'
import AIWorkspace from '@/components/platform/AIWorkspace.vue'
import EmptyState from '@/components/platform/EmptyState.vue'
import PageIntro from '@/components/platform/PageIntro.vue'
import SectionCard from '@/components/platform/SectionCard.vue'
import StatStrip from '@/components/platform/StatStrip.vue'
import { BRANDING } from '@/constants/branding'
import { getAIRuntimeErrorLabel, normalizeAIRuntimeErrorMessage } from '@/utils/aiRuntime'

export default {
  name: 'AIDiagnosis',
  components: {
    AIWorkspace,
    EmptyState,
    PageIntro,
    SectionCard,
    StatStrip
  },
  computed: {
    currentSceneMeta() {
      return this.getSceneMeta(this.formData.scene)
    },
    sceneOptions() {
      if (Array.isArray(this.overview?.diagnosisScenes) && this.overview.diagnosisScenes.length) {
        return this.overview.diagnosisScenes
      }
      return [
        { value: 'terminal_audit', label: '终端审计复盘', description: '适合终端操作审计和回放分析。' },
        { value: 'sql_work_order', label: 'SQL 工单分析', description: '适合 SQL 风险、回滚与审批场景。' },
        { value: 'inspection_report', label: '巡检报告分析', description: '适合巡检报告复盘与风险总结。' },
        { value: 'workload_capacity', label: '容量治理', description: '适合 HPA 与工作负载容量建议。' },
        { value: 'knowledge_search', label: '知识检索增强', description: '适合基于知识库条目做上下文分析。' }
      ]
    },
    statItems() {
      return [
        { label: '知识候选', value: this.knowledgeSuggestions.length, hint: '当前召回的知识条目', tone: 'primary' },
        { label: '已选知识', value: this.selectedKnowledgeIds.length, hint: '参与本次分析的知识上下文', tone: 'success' },
        { label: '模板数量', value: this.templateOptions.length, hint: '当前可选 Prompt Template', tone: 'warning' },
        { label: '历史会话', value: this.historyList.length, hint: '最近沉淀的 AI 会话', tone: 'neutral' }
      ]
    },
    statItemsEnhanced() {
      const items = this.statItems.map(item => ({ ...item }))
      if (items[0]) {
        items[0].hint = `知识库 ${this.overviewStats.knowledgeItems || 0} 条`
      }
      if (items[2]) {
        items[2].value = this.overviewStats.promptTemplates || this.templateOptions.length
        items[2].hint = `${this.runtimeStatus.provider} / ${this.runtimeStatus.model}`
      }
      if (items[3]) {
        items[3].value = this.overviewStats.diagnosisSessions || this.historyList.length
      }
      return items
    },
    overviewStats() {
      return this.overview?.stats || {}
    },
    runtimeStatus() {
      const runtime = this.overview?.runtime || {}
      return {
        status: runtime.status || 'fallback',
        statusText: runtime.statusText || '未接入实时 LLM',
        provider: runtime.provider || 'openai',
        model: runtime.model || 'gpt-5.4',
        reasoningEffort: runtime.reasoningEffort || 'medium',
        checkedAt: runtime.checkedAt || '',
        lastError: runtime.lastError || ''
      }
    },
    runtimeErrorLabel() {
      return getAIRuntimeErrorLabel(this.runtimeStatus.status)
    },
    runtimeErrorDisplay() {
      return normalizeAIRuntimeErrorMessage(this.runtimeStatus.lastError)
    },
    workspaceDomains() {
      return this.overview?.domains || []
    }
  },
  data() {
    return {
      brand: BRANDING,
      overview: null,
      diagnosing: false,
      historyVisible: false,
      historyList: [],
      historyDetail: [],
      knowledgeSuggestions: [],
      selectedKnowledgeIds: [],
      templateOptions: [],
      renderResult: {},
      routeContext: {
        source: '',
        namespace: '',
        workloadName: '',
        autoRun: false
      },
      formData: {
        scene: 'terminal_audit',
        targetId: '',
        keyword: 'terminal audit',
        templateName: 'terminal_audit_review'
      }
    }
  },
  methods: {
    getSceneMeta(scene) {
      const map = {
        terminal_audit: {
          targetLabel: '目标 ID',
          targetPlaceholder: '请输入终端会话 ID',
          keyword: 'terminal audit',
          keywordPlaceholder: '例如：terminal audit / kubectl exec',
          templateName: 'terminal_audit_review'
        },
        sql_work_order: {
          targetLabel: '目标 ID',
          targetPlaceholder: '请输入 SQL 工单 ID',
          keyword: 'sql rollback',
          keywordPlaceholder: '例如：sql rollback / ddl review',
          templateName: 'yaml_change_review'
        },
        alert_analysis: {
          targetLabel: '告警关键词',
          targetPlaceholder: '可选：输入告警关键词，留空则分析当前摘要',
          keyword: 'alert analysis',
          keywordPlaceholder: '例如：mysql / kubelet / 磁盘',
          templateName: 'incident_analysis'
        },
        deployment_review: {
          targetLabel: '发布 ID',
          targetPlaceholder: '请输入快速发布 ID',
          keyword: 'deployment review',
          keywordPlaceholder: '例如：rollback / 发布失败 / 环境治理',
          templateName: 'incident_analysis'
        },
        inspection_report: {
          targetLabel: '知识文章 ID',
          targetPlaceholder: '请输入巡检知识文章 ID',
          keyword: 'inspection report',
          keywordPlaceholder: '例如：inspection report / 巡检 runbook',
          templateName: 'incident_analysis'
        },
        workload_capacity: {
          targetLabel: '快照 ID',
          targetPlaceholder: '请输入容量建议历史快照 ID',
          keyword: 'hpa capacity autoscaling',
          keywordPlaceholder: '例如：demo-nginx hpa capacity autoscaling',
          templateName: 'incident_analysis'
        },
        knowledge_search: {
          targetLabel: '知识文章 ID',
          targetPlaceholder: '请输入知识文章 ID',
          keyword: 'knowledge context search',
          keywordPlaceholder: '例如：postgres rollback knowledge',
          templateName: 'incident_analysis'
        }
      }
      return map[scene] || map.terminal_audit
    },
    applySceneDefaults(scene, preserveTarget = false) {
      const meta = this.getSceneMeta(scene)
      this.formData.scene = scene
      if (!preserveTarget) {
        this.formData.targetId = ''
      }
      this.formData.keyword = meta.keyword
      this.formData.templateName = meta.templateName
      this.selectedKnowledgeIds = []
      this.renderResult = {}
    },
    handleSceneChange(scene) {
      this.routeContext = {
        source: '',
        namespace: '',
        workloadName: '',
        autoRun: false
      }
      this.applySceneDefaults(scene)
      this.loadKnowledgeSuggestions()
    },
    applyRouteQuery() {
      const query = this.$route?.query || {}
      const scene = String(query.scene || '').trim()
      if (!scene) return false
      const meta = this.getSceneMeta(scene)
      this.formData.scene = scene
      this.formData.targetId = String(query.targetId || '').trim()
      this.formData.keyword = String(query.keyword || meta.keyword).trim()
      this.formData.templateName = String(query.templateName || meta.templateName).trim()
      this.routeContext = {
        source: String(query.source || '').trim(),
        namespace: String(query.namespace || '').trim(),
        workloadName: String(query.workloadName || '').trim(),
        autoRun: String(query.autoRun || '') === '1'
      }
      return this.routeContext.autoRun && Boolean(this.formData.targetId)
    },
    isTargetRequired(scene) {
      return !['alert_analysis', 'knowledge_search'].includes(scene)
    },
    async loadOverview() {
      const { data: res } = await this.$api.getAIOverview()
      if (res.code === 200) this.overview = res.data || null
    },
    async loadTemplates() {
      const { data: res } = await this.$api.getAITemplates()
      if (res.code === 200) this.templateOptions = res.data || []
    },
    async loadKnowledgeSuggestions() {
      const { data: res } = await this.$api.suggestAIKnowledge(this.formData.keyword)
      if (res.code === 200) {
        this.knowledgeSuggestions = res.data || []
      }
    },
    toggleKnowledge(id) {
      if (this.selectedKnowledgeIds.includes(id)) {
        this.selectedKnowledgeIds = this.selectedKnowledgeIds.filter(item => item !== id)
      } else {
        this.selectedKnowledgeIds.push(id)
      }
    },
    async runDiagnosis() {
      if (this.isTargetRequired(this.formData.scene) && !this.formData.targetId) {
        return ElMessage.warning('请输入目标 ID')
      }
      this.diagnosing = true
      try {
        const { data: res } = await this.$api.diagnoseAI({
          scene: this.formData.scene,
          targetId: this.formData.targetId,
          keyword: this.formData.keyword,
          templateName: this.formData.templateName,
          knowledgeIds: this.selectedKnowledgeIds
        })
        if (res.code !== 200) throw new Error(res.message || '诊断失败')
        this.renderResult = res.data || {}
        if (Array.isArray(res.data?.knowledgeItems) && res.data.knowledgeItems.length) {
          this.knowledgeSuggestions = res.data.knowledgeItems
        }
        ElMessage.success('AI 诊断结果已生成')
        this.refreshHistory()
      } catch (error) {
        ElMessage.error(error.message || '诊断失败')
      } finally {
        this.diagnosing = false
      }
    },
    async refreshHistory() {
      const { data: res } = await this.$api.getAIHistory()
      if (res.code === 200) this.historyList = res.data || []
    },
    async showHistory(sessionId) {
      const { data: res } = await this.$api.getAIHistoryDetail(sessionId)
      if (res.code !== 200) return ElMessage.error(res.message || '获取历史失败')
      this.historyDetail = res.data || []
      this.historyVisible = true
    },
    openDomain(item) {
      if (item?.route) {
        this.$router.push(item.route)
      }
    }
  },
  watch: {
    '$route.query': {
      async handler() {
        const shouldAutoRun = this.applyRouteQuery()
        await this.loadKnowledgeSuggestions()
        if (shouldAutoRun) {
          await this.runDiagnosis()
        }
      }
    }
  },
  async mounted() {
    const shouldAutoRun = this.applyRouteQuery()
    await this.loadOverview()
    await this.loadTemplates()
    await this.loadKnowledgeSuggestions()
    await this.refreshHistory()
    if (shouldAutoRun) {
      await this.runDiagnosis()
    }
  }
}
</script>

<style scoped>
.ai-diagnosis__scene-list,
.ai-diagnosis__history,
.ai-diagnosis__history-detail {
  display: grid;
  gap: 12px;
}

.ai-diagnosis__positioning {
  margin-bottom: 12px;
  font-size: 13px;
  line-height: 1.7;
  color: var(--text-muted);
}

.ai-diagnosis__runtime-status {
  display: grid;
  gap: 10px;
}

.ai-diagnosis__runtime-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: fit-content;
  padding: 6px 12px;
  border-radius: var(--control-radius-pill);
  font-size: 12px;
  font-weight: 700;
}

.ai-diagnosis__runtime-badge.is-ready {
  background: rgba(34, 197, 94, 0.16);
  color: #86efac;
  border: 1px solid rgba(34, 197, 94, 0.28);
}

.ai-diagnosis__runtime-badge.is-fallback {
  background: rgba(245, 158, 11, 0.16);
  color: #fcd34d;
  border: 1px solid rgba(245, 158, 11, 0.28);
}

.ai-diagnosis__runtime-badge.is-degraded {
  background: rgba(248, 113, 113, 0.16);
  color: #fca5a5;
  border: 1px solid rgba(248, 113, 113, 0.28);
}

.ai-diagnosis__runtime-meta {
  font-size: 13px;
  line-height: 1.7;
  color: var(--text-muted);
}

.ai-diagnosis__runtime-meta--warning {
  color: #fca5a5;
  padding: 10px 12px;
  border-radius: 14px;
  background: rgba(127, 29, 29, 0.18);
  border: 1px solid rgba(248, 113, 113, 0.22);
  word-break: break-word;
}

.ai-diagnosis__scene,
.ai-diagnosis__history-item {
  width: 100%;
  text-align: left;
  padding: 14px;
  border-radius: var(--action-card-radius);
  border: 1px solid var(--border-subtle);
  background: rgba(255, 255, 255, 0.03);
  color: var(--text-secondary);
  transition: all var(--transition-fast);
}

.ai-diagnosis__scene:hover,
.ai-diagnosis__history-item:hover,
.ai-diagnosis__scene.active {
  border-color: var(--border-strong);
  background: rgba(53, 164, 255, 0.08);
}

.ai-diagnosis__scene-title,
.ai-diagnosis__history-id,
.ai-diagnosis__history-role {
  font-size: 14px;
  font-weight: 700;
  color: var(--text-primary);
}

.ai-diagnosis__scene-desc,
.ai-diagnosis__history-meta {
  margin-top: 6px;
  font-size: 12px;
  line-height: 1.6;
  color: var(--text-muted);
}

.ai-diagnosis__form {
  max-width: 920px;
}

.ai-diagnosis__code {
  margin: 0;
  min-height: 260px;
  padding: 16px;
  border-radius: var(--action-card-radius);
  background: rgba(3, 9, 18, 0.94);
  color: #e2e8f0;
  white-space: pre-wrap;
  word-break: break-word;
  font-family: Consolas, Monaco, monospace;
}

.ai-diagnosis__history-message {
  display: grid;
  gap: 8px;
}
</style>
