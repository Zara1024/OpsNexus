<template>
  <AIWorkspace
    eyebrow="AI Workspace"
    title="AI 智能运维助手工作台"
    subtitle="支持接入多种大模型，在统一工作区中完成对话管控、Agent 协作、知识检索、诊断分析与巡检辅助。"
  >
    <template #actions>
      <el-button type="primary" @click="startNewConversation">新建会话</el-button>
      <el-button @click="refreshHistory">刷新历史</el-button>
    </template>

    <template #intro>
      <PageIntro
        title="工作方式"
        text="直接输入主机 IP 或主机名即可发起查询，也可以从右侧模型接入、快捷提示和巡检模板开始，让诊断与处置共享一套上下文。"
      />
    </template>

    <template #rail>
      <aside class="assistant-sidebar">
        <div class="sidebar-card brand-card">
          <div class="brand-row">
            <div class="brand-logo-wrap">
              <img :src="brand.logo" :alt="brand.name" class="brand-logo">
            </div>
            <div>
              <div class="brand-eyebrow">AIOps Command Center</div>
              <div class="brand-title">{{ brand.name }}</div>
              <div class="brand-subtitle">{{ brand.slogan }}</div>
            </div>
          </div>
          <p class="brand-description">{{ brand.description }}</p>
          <p class="brand-positioning">{{ brand.assistantPositioning }}</p>
          <div class="model-chip-row">
            <span v-for="item in modelProviders" :key="item" class="model-chip">{{ item }}</span>
          </div>
        </div>

        <div class="sidebar-card runtime-card">
          <div class="section-title">AI 运行态</div>
          <div class="runtime-status-row">
            <span :class="['runtime-badge', `runtime-badge--${runtimeStatus.status}`]">
              {{ runtimeStatus.statusText }}
            </span>
            <span class="runtime-provider">{{ runtimeStatus.provider }} / {{ runtimeStatus.model }}</span>
          </div>
          <div v-if="runtimeStatus.checkedAt" class="runtime-meta">最近探测：{{ runtimeStatus.checkedAt }}</div>
          <div v-if="runtimeErrorDisplay" class="runtime-meta runtime-meta--warning">{{ runtimeErrorLabel }}：{{ runtimeErrorDisplay }}</div>
          <div class="runtime-meta">推理强度：{{ runtimeStatus.reasoningEffort }}</div>
          <div class="runtime-meta">待确认动作：{{ overviewStats.pendingConfirmations || 0 }}</div>
          <div class="runtime-meta">诊断会话：{{ overviewStats.diagnosisSessions || 0 }}</div>
        </div>

        <div class="sidebar-card domain-card">
          <div class="section-title">工作台入口</div>
          <button
            v-for="item in workspaceDomains"
            :key="item.key"
            class="domain-item"
            @click="handleDomainAction(item)"
          >
            <div class="domain-item__title">{{ item.label }}</div>
            <div class="domain-item__desc">{{ item.description }}</div>
          </button>
        </div>

        <div class="sidebar-card capability-card">
          <div class="section-title">能力场景</div>
          <button
            v-for="item in capabilityScenes"
            :key="item.name"
            class="capability-item capability-item--interactive"
            @click="handleCapabilityScene(item)"
          >
            <div class="capability-name">{{ item.name }}</div>
            <div class="capability-desc">{{ item.description }}</div>
          </button>
        </div>

        <div class="sidebar-card context-card">
          <div class="section-title">当前上下文</div>
          <div v-if="assistantContext" class="context-body">
            <div class="context-item">范围：{{ contextScopeText }}</div>
            <div class="context-item" v-if="assistantContext.currentHostName">主机：{{ assistantContext.currentHostName }}</div>
            <div class="context-item" v-if="assistantContext.currentClusterName">集群：{{ assistantContext.currentClusterName }}</div>
            <div class="context-item" v-if="assistantContext.currentNamespace">命名空间：{{ assistantContext.currentNamespace }}</div>
            <div class="context-item" v-if="assistantContext.currentWorkloadName">工作负载：{{ assistantContext.currentWorkloadName }}</div>
            <div class="context-summary">{{ assistantContext.summary || '当前还没有显式上下文。' }}</div>
          </div>
          <div v-else class="context-empty">当前还没有显式上下文。</div>
        </div>

        <div class="sidebar-card prompt-card">
          <div class="section-title">快捷提示</div>
          <button
            v-for="prompt in promptList"
            :key="prompt"
            class="prompt-chip"
            @click="sendMessage(prompt)"
          >
            {{ prompt }}
          </button>
        </div>

        <div class="sidebar-card history-card">
          <div class="history-header">
            <div class="section-title">最近会话</div>
            <el-button link type="primary" @click="startNewConversation">新建</el-button>
          </div>
          <div class="history-list">
            <button
              v-for="item in historyList"
              :key="item.sessionId"
              :class="['history-item', { active: item.sessionId === currentSessionId }]"
              @click="openHistorySession(item.sessionId)"
            >
              <div class="history-id">{{ item.sessionId }}</div>
              <div class="history-meta">{{ formatHistoryMeta(item) }}</div>
            </button>
          </div>
        </div>

        <div class="sidebar-card template-card">
          <div class="section-title">巡检模板</div>
          <button
            v-for="item in inspectionTemplates"
            :key="item.id"
            class="template-item"
            @click="sendMessage(`使用${item.name}为当前主机生成巡检报告`)"
          >
            <div class="template-name">{{ item.name }}</div>
            <div class="template-desc">{{ item.description }}</div>
          </button>
        </div>

        <div class="sidebar-card reports-card">
          <div class="section-title">报告归档</div>
          <div class="report-list">
            <button
              v-for="item in recentReports"
              :key="item.id"
              class="report-item"
              @click="sendMessage(`查看最近巡检报告列表`)"
            >
              <div class="report-name">{{ item.templateName || '巡检报告' }}</div>
              <div class="report-meta">{{ item.targetName }} · {{ item.createdAt }}</div>
            </button>
          </div>
        </div>
      </aside>
    </template>

    <StatStrip :items="assistantStatsEnhanced" />

    <div class="chat-shell">
      <div class="chat-stream">
        <article
          v-for="(item, index) in messages"
          :key="`${item.role}-${index}`"
          :class="['message-card', `message-card--${item.role}`]"
        >
          <div class="message-role">{{ item.role === 'assistant' ? 'AI 助手' : '你' }}</div>
          <div v-if="item.role === 'assistant' && item.data?.model" class="message-meta">
            <span class="meta-tag" :class="{ 'meta-tag--fallback': !item.data?.usedLlm }">
              {{ item.data?.usedLlm ? `LLM: ${item.data.model}` : `Fallback: ${item.data.model}` }}
            </span>
            <span v-if="item.data?.fallbackReason" class="meta-reason">{{ item.data.fallbackReason }}</span>
          </div>
          <div class="message-content">{{ item.content }}</div>

          <div v-if="item.data?.context" class="context-inline">
            <span>上下文：{{ item.data.context.summary || '已更新' }}</span>
          </div>

          <div v-if="item.data?.toolSteps?.length" class="steps-panel">
            <div class="panel-title">工具步骤</div>
            <div class="steps-list">
              <div v-for="step in item.data.toolSteps" :key="`${step.tool}-${step.summary}`" class="step-item">
                <div class="step-title">{{ step.tool }}</div>
                <div class="step-status">{{ step.status }}</div>
                <div class="step-summary">{{ step.summary }}</div>
              </div>
            </div>
          </div>

          <div v-if="item.data?.hostMatches?.length" class="host-grid">
            <div
              v-for="host in item.data.hostMatches"
              :key="host.id"
              class="host-card"
            >
              <div class="host-title-row">
                <div class="host-title">{{ host.hostName }}</div>
                <span class="host-status">{{ host.statusText }}</span>
              </div>
              <div class="host-meta">分组：{{ host.groupName || '未分组' }}</div>
              <div class="host-meta">管理地址：{{ host.sshIp || host.privateIp || host.publicIp || '-' }}</div>
              <div class="host-meta">系统：{{ host.os || '-' }}</div>
              <div class="host-meta">规格：{{ host.cpu || '-' }} / {{ host.memory || '-' }}</div>
            </div>
          </div>

          <div v-if="item.data?.commandResult" class="result-panel">
            <div class="panel-title">命令执行结果</div>
            <div class="panel-subtitle">
              {{ item.data.commandResult.hostName }} · {{ item.data.commandResult.command }}
            </div>
            <pre class="panel-block">{{ item.data.commandResult.output || '-' }}</pre>
          </div>

          <div v-if="item.data?.inspectionResult" class="result-panel">
            <div class="panel-title">巡检摘要</div>
            <div class="panel-subtitle">模板：{{ item.data.inspectionResult.templateName || '默认模板' }}</div>
            <div class="panel-subtitle">{{ item.data.inspectionResult.summary }}</div>
            <div class="inspection-grid">
              <div
                v-for="check in item.data.inspectionResult.checks"
                :key="`${check.name}-${check.command}`"
                class="inspection-item"
              >
                <div class="inspection-name">{{ check.name }}</div>
                <div class="inspection-status">{{ check.status === 'success' ? '成功' : '失败' }}</div>
                <div class="inspection-command">{{ check.command }}</div>
              </div>
            </div>
            <pre class="panel-block">{{ item.data.inspectionResult.report || '-' }}</pre>
          </div>

          <div v-if="item.data?.pendingConfirmation" class="result-panel confirmation-panel">
            <div class="panel-title">高风险确认</div>
            <div class="panel-subtitle">{{ item.data.pendingConfirmation.summary }}</div>
            <div class="action-row">
              <el-button size="small" type="danger" @click="decideConfirmation(item.data.pendingConfirmation.id, 'approve')">
                确认执行
              </el-button>
              <el-button size="small" @click="decideConfirmation(item.data.pendingConfirmation.id, 'cancel')">
                取消
              </el-button>
            </div>
          </div>

          <div v-if="item.data?.actions?.length" class="action-row">
            <el-button
              v-for="action in item.data.actions"
              :key="`${action.kind}-${action.label}`"
              size="small"
              @click="sendMessage(action.message)"
            >
              {{ action.label }}
            </el-button>
          </div>
        </article>
        <div ref="streamBottomRef"></div>
      </div>

      <div class="composer">
        <el-input
          v-model="draft"
          type="textarea"
          :rows="4"
          resize="none"
          placeholder="例如：查询主机 10.0.0.200；查看主机 10.0.0.200 的磁盘占用；为主机 10.0.0.200 生成巡检报告"
          @keydown.enter.exact.prevent="sendMessage()"
        />
        <div class="composer-actions">
          <span class="composer-tip">支持主机 IP 和主机名检索。</span>
          <el-button type="primary" :loading="sending" @click="sendMessage()">发送</el-button>
        </div>
      </div>
    </div>
  </AIWorkspace>
</template>

<script setup>
import { computed, nextTick, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import aiAPI from '@/api/ai'
import { BRANDING } from '@/constants/branding'
import { getAIRuntimeErrorLabel, normalizeAIRuntimeErrorMessage } from '@/utils/aiRuntime'
import AIWorkspace from '@/components/platform/AIWorkspace.vue'
import PageIntro from '@/components/platform/PageIntro.vue'
import StatStrip from '@/components/platform/StatStrip.vue'

const route = useRoute()
const router = useRouter()
const brand = BRANDING
const modelProviders = computed(() => brand.aiModels || [])
const overview = ref(null)

const quickPrompts = [
  '查询主机 10.0.0.200',
  '查看主机 10.0.0.200 的磁盘占用',
  '在主机 10.0.0.200 执行 `free -m`',
  '为主机 10.0.0.200 生成巡检报告',
  '结合知识库总结最近一次巡检风险',
  '根据当前告警上下文给出处置建议'
]

const capabilityScenes = [
  {
    name: '对话控机',
    description: '主机定位、资产信息查看、只读命令执行和快捷联动。',
    prompt: '帮我定位一台主机，并展示可执行的只读检查动作'
  },
  {
    name: 'Agent 与工具编排',
    description: '以工具步骤形式串联查询、诊断、确认和执行上下文。',
    prompt: '用分步骤工具链帮我做一次故障排查'
  },
  {
    name: '知识检索增强',
    description: '结合知识库、工单和历史会话沉淀出可复用的答案结构。',
    prompt: '结合知识库和历史工单，帮我总结当前问题的处理方案'
  },
  {
    name: '智能巡检',
    description: '自动汇总磁盘、内存、端口和基础状态，生成巡检报告。',
    prompt: '为当前主机生成一份巡检报告'
  }
]

const welcomeMessage = {
  role: 'assistant',
  content: [
    '这里是 OpsNexus 的 AI 智能运维助手工作台，可在统一上下文中承载多模型问答、巡检和诊断协作。',
    '工作台支持接入 OpenAI、Claude、Gemini、DeepSeek、Qwen、Ollama 与本地模型。',
    '你可以直接告诉我主机 IP，或者像“主机 prod-api-01”这样的名称。'
  ].join('\n'),
  data: {
    actions: quickPrompts.map((prompt) => ({
      label: prompt,
      message: prompt,
      kind: 'prompt'
    }))
  }
}

const promptList = computed(() => {
  if (Array.isArray(overview.value?.quickPrompts) && overview.value.quickPrompts.length) {
    return overview.value.quickPrompts
  }
  return quickPrompts
})

const overviewStats = computed(() => overview.value?.stats || {})

const runtimeStatus = computed(() => {
  const runtime = overview.value?.runtime || {}
  return {
    status: runtime.status || 'fallback',
    statusText: runtime.statusText || '未接入实时 LLM',
    provider: runtime.provider || 'openai',
    model: runtime.model || 'gpt-5.4',
    reasoningEffort: runtime.reasoningEffort || 'medium',
    checkedAt: runtime.checkedAt || '',
    lastError: runtime.lastError || ''
  }
})

const runtimeErrorLabel = computed(() => getAIRuntimeErrorLabel(runtimeStatus.value.status))
const runtimeErrorDisplay = computed(() => normalizeAIRuntimeErrorMessage(runtimeStatus.value.lastError))

const workspaceDomains = computed(() => {
  if (Array.isArray(overview.value?.domains) && overview.value.domains.length) {
    return overview.value.domains
  }
  return []
})

const messages = ref([welcomeMessage])
const historyList = ref([])
const currentSessionId = ref('')
const draft = ref('')
const sending = ref(false)
const streamBottomRef = ref(null)
const inspectionTemplates = ref([])
const recentReports = ref([])
const assistantContext = ref(null)

const contextScopeText = computed(() => {
  const scope = assistantContext.value?.currentScope
  const labels = {
    host: '主机',
    workload: '工作负载',
    alert: '告警',
    workorder: '工单',
    deployment: '发布'
  }
  return labels[scope] || '未设定'
})

const assistantStats = computed(() => ([
  {
    label: '接入模型',
    value: modelProviders.value.length,
    hint: '支持多模型路由与接入扩展',
    tone: 'primary'
  },
  {
    label: '上下文状态',
    value: assistantContext.value ? 1 : 0,
    hint: assistantContext.value ? contextScopeText.value : '当前无显式上下文',
    tone: 'success'
  },
  {
    label: '巡检模板',
    value: inspectionTemplates.value.length,
    hint: '可直接发起巡检',
    tone: 'warning'
  },
  {
    label: '最近会话',
    value: historyList.value.length,
    hint: '可快速回看历史对话',
    tone: 'neutral'
  }
]))

const assistantStatsEnhanced = computed(() => {
  const items = assistantStats.value.map((item) => ({ ...item }))
  if (items[0]) {
    items[0].hint = `${runtimeStatus.value.provider} / ${runtimeStatus.value.model}`
  }
  if (items[2]) {
    items[2].value = overviewStats.value.inspectionTemplates || inspectionTemplates.value.length
  }
  if (items[3]) {
    items[3].value = overviewStats.value.assistantSessions || historyList.value.length
  }
  return items
})

const scrollToBottom = async () => {
  await nextTick()
  streamBottomRef.value?.scrollIntoView({ behavior: 'smooth', block: 'end' })
}

const refreshHistory = async () => {
  const { data: res } = await aiAPI.getAIAssistantHistory()
  if (res.code === 200) {
    historyList.value = res.data || []
  }
}

const loadTemplates = async () => {
  const { data: res } = await aiAPI.getAIAssistantTemplates()
  if (res.code === 200) {
    inspectionTemplates.value = res.data || []
  }
}

const loadReports = async () => {
  const { data: res } = await aiAPI.getAIAssistantReports()
  if (res.code === 200) {
    recentReports.value = res.data || []
  }
}

const loadOverview = async () => {
  const { data: res } = await aiAPI.getAIOverview()
  if (res.code === 200) {
    overview.value = res.data || null
  }
}

const startNewConversation = () => {
  currentSessionId.value = ''
  draft.value = ''
  assistantContext.value = null
  messages.value = [welcomeMessage]
  scrollToBottom()
}

const pushAssistantMessage = (payload) => {
  if (payload.context) {
    assistantContext.value = payload.context
  }
  if (Array.isArray(payload.availableTemplates) && payload.availableTemplates.length) {
    inspectionTemplates.value = payload.availableTemplates
  }
  if (Array.isArray(payload.recentReports)) {
    recentReports.value = payload.recentReports
  }
  messages.value.push({
    role: 'assistant',
    content: payload.assistantMessage || 'AI 助手本次没有返回内容。',
    data: payload
  })
}

const decideConfirmation = async (id, decision) => {
  if (!id) {
    return
  }
  try {
    const { data: res } = await aiAPI.decideAIAssistantConfirmation(id, { decision })
    if (res.code !== 200) {
      throw new Error(res.message || '确认失败')
    }
    messages.value.push({
      role: 'assistant',
      content: res.data?.resultSummary || `确认任务 ${decision === 'approve' ? '已执行' : '已取消'}`,
      data: {
        model: '',
        usedLlm: false
      }
    })
    await loadReports()
    scrollToBottom()
  } catch (error) {
    messages.value.push({
      role: 'assistant',
      content: error.message || '确认执行失败',
      data: {
        model: '',
        usedLlm: false
      }
    })
  }
}

const sendMessage = async (presetMessage = '') => {
  const content = String(presetMessage || draft.value).trim()
  if (!content || sending.value) {
    return
  }

  messages.value.push({
    role: 'user',
    content,
    data: null
  })
  draft.value = ''
  sending.value = true
  scrollToBottom()

  try {
    const { data: res } = await aiAPI.chatAIAssistant({
      sessionId: currentSessionId.value,
      message: content
    })
    if (res.code !== 200) {
      throw new Error(res.message || '发送失败')
    }
    const payload = res.data || {}
    currentSessionId.value = payload.sessionId || currentSessionId.value
    pushAssistantMessage(payload)
    await refreshHistory()
    scrollToBottom()
  } catch (error) {
    pushAssistantMessage({
      assistantMessage: error.message || 'AI 助手调用失败，请稍后重试。',
      actions: []
    })
  } finally {
    sending.value = false
  }
}

const openHistorySession = async (sessionId) => {
  const { data: res } = await aiAPI.getAIAssistantHistoryDetail(sessionId)
  if (res.code !== 200) {
    return
  }
  currentSessionId.value = sessionId
  const rows = res.data || []
  messages.value = rows.map((item) => ({
    role: item.role,
    content: item.message,
    data: null
  }))
  if (!messages.value.length) {
    messages.value = [welcomeMessage]
  }
  scrollToBottom()
}

const formatHistoryMeta = (item) => {
  const count = item.messageCount || 0
  const time = item.latestTime || ''
  return `${count} 条消息 · ${time}`
}

const handleDomainAction = async (item) => {
  if (item?.route && String(item.route).startsWith('/ai/diagnosis')) {
    await router.push(item.route)
    return
  }
  if (item?.prompt) {
    await sendMessage(item.prompt)
  }
}

const handleCapabilityScene = async (item) => {
  if (!item?.prompt) {
    return
  }
  await sendMessage(item.prompt)
}

const tryAutoPrompt = async () => {
  const prompt = String(route.query.prompt || '').trim()
  if (prompt) {
    await sendMessage(prompt)
  }
}

watch(
  () => route.query.prompt,
  async (value, oldValue) => {
    const nextPrompt = String(value || '').trim()
    const previousPrompt = String(oldValue || '').trim()
    if (nextPrompt && nextPrompt !== previousPrompt) {
      await sendMessage(nextPrompt)
    }
  }
)

onMounted(async () => {
  await loadOverview()
  await refreshHistory()
  await loadTemplates()
  await loadReports()
  await tryAutoPrompt()
  scrollToBottom()
})
</script>

<style scoped lang="scss">
.assistant-page {
  min-height: calc(100vh - 50px);
  padding: 20px;
  background:
    radial-gradient(circle at top left, rgba(56, 189, 248, 0.14), transparent 30%),
    radial-gradient(circle at top right, rgba(14, 165, 233, 0.14), transparent 26%),
    linear-gradient(145deg, #07111f 0%, #09172b 42%, #050b16 100%);
}

.assistant-grid {
  display: grid;
  grid-template-columns: 320px minmax(0, 1fr);
  gap: 20px;
}

.assistant-sidebar,
.assistant-main {
  min-width: 0;
}

.sidebar-card,
.assistant-hero,
.chat-shell {
  border: 1px solid rgba(148, 163, 184, 0.16);
  background: rgba(7, 16, 31, 0.72);
  box-shadow: 0 24px 60px rgba(2, 8, 23, 0.32);
  backdrop-filter: blur(18px);
}

.sidebar-card {
  border-radius: var(--radius-lg);
  padding: 18px;
}

.assistant-sidebar {
  display: grid;
  gap: 16px;
}

.brand-row {
  display: flex;
  align-items: center;
  gap: 14px;
}

.brand-logo-wrap {
  width: 56px;
  height: 56px;
  border-radius: 16px;
  display: grid;
  place-items: center;
  background: linear-gradient(135deg, rgba(14, 165, 233, 0.2), rgba(99, 102, 241, 0.16));
  border: 1px solid rgba(125, 211, 252, 0.24);
}

.brand-logo {
  width: 34px;
  height: 34px;
}

.brand-eyebrow,
.hero-eyebrow {
  color: rgba(125, 211, 252, 0.88);
  font-size: 12px;
  letter-spacing: 0.18em;
  text-transform: uppercase;
}

.brand-title {
  margin-top: 4px;
  color: #f8fafc;
  font-size: 24px;
  font-weight: 700;
}

.brand-subtitle {
  color: rgba(226, 232, 240, 0.76);
  font-size: 13px;
}

.brand-description,
.brand-positioning,
.capability-desc,
.hero-description,
.composer-tip,
.history-meta,
.panel-subtitle,
.host-meta,
.inspection-command {
  color: rgba(148, 163, 184, 0.88);
}

.brand-description {
  margin: 14px 0 0;
  line-height: 1.7;
}

.brand-positioning {
  margin: 12px 0 0;
  line-height: 1.7;
}

.runtime-status-row {
  display: grid;
  gap: 10px;
  margin-bottom: 10px;
}

.runtime-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: fit-content;
  padding: 6px 12px;
  border-radius: var(--control-radius-pill);
  font-size: 12px;
  font-weight: 700;
}

.runtime-badge--ready {
  background: rgba(34, 197, 94, 0.16);
  color: #86efac;
  border: 1px solid rgba(34, 197, 94, 0.28);
}

.runtime-badge--fallback {
  background: rgba(245, 158, 11, 0.16);
  color: #fcd34d;
  border: 1px solid rgba(245, 158, 11, 0.28);
}

.runtime-badge--degraded {
  background: rgba(248, 113, 113, 0.16);
  color: #fca5a5;
  border: 1px solid rgba(248, 113, 113, 0.28);
}

.runtime-provider {
  color: #f8fafc;
  font-size: 14px;
  font-weight: 700;
}

.runtime-meta {
  color: rgba(148, 163, 184, 0.88);
  font-size: 13px;
  line-height: 1.7;
}

.runtime-meta--warning {
  color: #fca5a5;
  padding: 10px 12px;
  border-radius: 14px;
  background: rgba(127, 29, 29, 0.18);
  border: 1px solid rgba(248, 113, 113, 0.22);
  word-break: break-word;
}

.section-title,
.capability-name,
.panel-title,
.host-title,
.inspection-name {
  color: #f8fafc;
  font-weight: 700;
}

.capability-card,
.runtime-card,
.domain-card,
.context-card,
.prompt-card,
.history-card,
.template-card,
.reports-card {
  display: grid;
  gap: 12px;
}

.domain-item {
  width: 100%;
  text-align: left;
  padding: 14px;
  border-radius: var(--action-card-radius);
  border: 1px solid rgba(148, 163, 184, 0.14);
  background: rgba(15, 23, 42, 0.52);
  color: #e2e8f0;
  cursor: pointer;
  transition: all 0.25s ease;
}

.domain-item:hover {
  border-color: rgba(56, 189, 248, 0.36);
  background: rgba(14, 165, 233, 0.14);
  transform: translateY(-1px);
}

.domain-item__title {
  color: #f8fafc;
  font-size: 14px;
  font-weight: 700;
}

.domain-item__desc {
  margin-top: 6px;
  font-size: 12px;
  line-height: 1.6;
  color: rgba(148, 163, 184, 0.88);
}

.model-chip-row {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.model-chip {
  display: inline-flex;
  align-items: center;
  min-height: 28px;
  padding: 4px 10px;
  border-radius: var(--control-radius-pill);
  border: 1px solid rgba(125, 211, 252, 0.18);
  background: rgba(14, 165, 233, 0.08);
  color: rgba(226, 232, 240, 0.94);
  font-size: 11px;
  font-weight: 700;
}

.context-body,
.report-list {
  display: grid;
  gap: 10px;
}

.context-item,
.context-summary,
.context-empty,
.template-desc,
.report-meta,
.step-summary {
  color: rgba(148, 163, 184, 0.88);
}

.template-item,
.report-item {
  width: 100%;
  text-align: left;
  border: 1px solid rgba(148, 163, 184, 0.14);
  background: rgba(15, 23, 42, 0.52);
  color: #e2e8f0;
  border-radius: var(--action-card-radius);
  padding: 12px 14px;
  cursor: pointer;
  transition: all 0.25s ease;
}

.template-item:hover,
.report-item:hover {
  border-color: rgba(56, 189, 248, 0.36);
  background: rgba(14, 165, 233, 0.14);
}

.template-name,
.report-name,
.step-title {
  color: #f8fafc;
  font-weight: 600;
}

.capability-item {
  border-radius: var(--action-card-radius);
  padding: 14px;
  border: 1px solid rgba(125, 211, 252, 0.12);
  background: rgba(15, 23, 42, 0.48);
}

.capability-item--interactive {
  width: 100%;
  text-align: left;
  cursor: pointer;
  transition: all 0.25s ease;
}

.capability-item--interactive:hover {
  border-color: rgba(56, 189, 248, 0.36);
  background: rgba(14, 165, 233, 0.14);
  transform: translateY(-1px);
}

.prompt-chip,
.history-item {
  border: 1px solid rgba(148, 163, 184, 0.14);
  background: rgba(15, 23, 42, 0.52);
  color: #e2e8f0;
  cursor: pointer;
  transition: all 0.25s ease;
}

.prompt-chip {
  width: 100%;
  text-align: left;
  border-radius: var(--action-card-radius);
  padding: 12px 14px;
}

.prompt-chip:hover,
.history-item:hover,
.history-item.active {
  border-color: rgba(56, 189, 248, 0.36);
  background: rgba(14, 165, 233, 0.14);
  transform: translateY(-1px);
}

.history-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.history-list {
  display: grid;
  gap: 10px;
}

.history-item {
  width: 100%;
  padding: 12px 14px;
  border-radius: var(--action-card-radius);
  text-align: left;
}

.history-id {
  color: #f8fafc;
  font-size: 13px;
  font-weight: 600;
}

.assistant-main {
  display: grid;
  gap: 16px;
}

.assistant-hero {
  border-radius: 26px;
  padding: 24px 28px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 20px;
}

.hero-title {
  margin: 8px 0 10px;
  color: #f8fafc;
  font-size: 32px;
  line-height: 1.2;
}

.hero-description {
  max-width: 760px;
  margin: 0;
  line-height: 1.8;
}

.hero-actions {
  display: flex;
  gap: 12px;
  flex-shrink: 0;
}

.chat-shell {
  border-radius: var(--radius-lg);
  padding: 20px;
  display: grid;
  gap: 18px;
  min-height: 720px;
}

.chat-stream {
  display: grid;
  gap: 14px;
  align-content: start;
  min-height: 0;
  overflow: auto;
  padding-right: 6px;
}

.message-card {
  max-width: 92%;
  border-radius: 20px;
  padding: 18px;
  border: 1px solid rgba(148, 163, 184, 0.12);
}

.message-card--assistant {
  justify-self: start;
  background: rgba(10, 20, 36, 0.88);
}

.message-card--user {
  justify-self: end;
  background: linear-gradient(135deg, rgba(14, 165, 233, 0.24), rgba(99, 102, 241, 0.22));
}

.message-role {
  margin-bottom: 10px;
  color: rgba(125, 211, 252, 0.92);
  font-size: 12px;
  font-weight: 700;
  letter-spacing: 0.12em;
  text-transform: uppercase;
}

.message-content,
.panel-block {
  color: #e2e8f0;
  white-space: pre-wrap;
  word-break: break-word;
  line-height: 1.8;
}

.context-inline {
  margin-top: 12px;
  color: rgba(125, 211, 252, 0.88);
  font-size: 13px;
}

.steps-panel {
  margin-top: 14px;
  padding: 14px;
  border-radius: var(--action-card-radius);
  border: 1px solid rgba(148, 163, 184, 0.14);
  background: rgba(15, 23, 42, 0.5);
}

.steps-list {
  display: grid;
  gap: 10px;
  margin-top: 10px;
}

.step-item {
  display: grid;
  grid-template-columns: 160px 90px 1fr;
  gap: 12px;
  align-items: start;
  padding: 10px 12px;
  border-radius: 14px;
  background: rgba(2, 6, 23, 0.52);
}

.step-status {
  color: #7dd3fc;
  font-size: 13px;
}

.message-meta {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 10px;
  flex-wrap: wrap;
}

.meta-tag {
  padding: 4px 10px;
  border-radius: var(--control-radius-pill);
  font-size: 12px;
  color: #7dd3fc;
  background: rgba(14, 165, 233, 0.16);
}

.meta-tag--fallback {
  color: #fbbf24;
  background: rgba(251, 191, 36, 0.14);
}

.meta-reason {
  color: rgba(148, 163, 184, 0.8);
  font-size: 12px;
}

.host-grid,
.inspection-grid {
  display: grid;
  gap: 12px;
  margin-top: 14px;
}

.host-grid {
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
}

.host-card,
.inspection-item {
  border-radius: var(--action-card-radius);
  padding: 14px;
  border: 1px solid rgba(148, 163, 184, 0.14);
  background: rgba(15, 23, 42, 0.56);
}

.host-title-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 8px;
}

.host-status,
.inspection-status {
  padding: 4px 10px;
  border-radius: var(--control-radius-pill);
  font-size: 12px;
  color: #7dd3fc;
  background: rgba(14, 165, 233, 0.16);
}

.result-panel {
  margin-top: 16px;
  padding: 16px;
  border-radius: var(--action-card-radius);
  border: 1px solid rgba(148, 163, 184, 0.14);
  background: rgba(2, 6, 23, 0.52);
}

.panel-block {
  margin: 12px 0 0;
  padding: 14px;
  border-radius: 16px;
  background: rgba(2, 6, 23, 0.88);
  font-family: Consolas, Monaco, monospace;
}

.action-row {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  margin-top: 16px;
}

.composer {
  padding-top: 8px;
  border-top: 1px solid rgba(148, 163, 184, 0.12);
}

.composer :deep(.el-textarea__inner) {
  min-height: 118px;
  border-radius: 18px;
  border: 1px solid rgba(148, 163, 184, 0.18);
  background: rgba(2, 6, 23, 0.68);
  color: #e2e8f0;
  box-shadow: none;
}

.composer :deep(.el-textarea__inner:focus) {
  border-color: rgba(56, 189, 248, 0.44);
}

.composer-actions {
  margin-top: 14px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

@media (max-width: 1280px) {
  .assistant-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 900px) {
  .assistant-page {
    padding: 12px;
  }

  .assistant-hero {
    padding: 20px;
    display: grid;
  }

  .hero-title {
    font-size: 26px;
  }

  .hero-actions,
  .composer-actions,
  .action-row {
    flex-wrap: wrap;
  }

  .step-item {
    grid-template-columns: 1fr;
  }

  .message-card {
    max-width: 100%;
  }
}
</style>
