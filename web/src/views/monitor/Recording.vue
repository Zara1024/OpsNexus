<template>
  <TablePage
    class="terminal-audit-page"
    section-title="审计会话列表"
    section-subtitle="统一展示 SSH / Pod / kubectl 会话，支持命令检索、录像回放、健康检测与原始日志下载。"
    wide
  >
    <template #header>
      <PageHeader
        eyebrow="Terminal Audit"
        title="终端录像审计"
        subtitle="先聚焦会话、风险和最新命令，再下钻回放和原始日志，避免高密度信息在首屏相互争抢注意力。"
      >
        <template #actions>
          <el-button type="primary" @click="refreshAll">刷新</el-button>
        </template>
        <template #intro>
          <PageIntro
            title="排查建议"
            text="优先关注“可回放录像”和“风险会话”，命令聚合历史会话保留统一提示；详情抽屉收窄后，命令审计和录像回放按区域分组，减少横向扫描成本。"
          />
        </template>
      </PageHeader>
    </template>

    <template #toolbar>
      <PageToolbar>
        <el-form :model="queryParams" :inline="true" class="filter-cluster">
          <el-form-item class="filter-field filter-field--wide" label="会话ID">
            <el-input v-model="queryParams.sessionId" clearable @keyup.enter="handleQuery" />
          </el-form-item>
          <el-form-item class="filter-field filter-field--wide" label="目标主机">
            <el-input v-model="queryParams.hostKeyword" clearable @keyup.enter="handleQuery" />
          </el-form-item>
          <el-form-item class="filter-field filter-field--wide" label="命令关键字">
            <el-input v-model="queryParams.keyword" clearable @keyup.enter="handleQuery" />
          </el-form-item>
          <el-form-item class="filter-field filter-field--compact" label="风险">
            <el-select v-model="queryParams.riskLevel" clearable>
              <el-option label="低风险" :value="0" />
              <el-option label="中风险" :value="1" />
              <el-option label="高风险" :value="2" />
            </el-select>
          </el-form-item>
          <el-form-item class="filter-field filter-field--compact" label="仅敏感命令">
            <el-switch v-model="queryParams.sensitiveOnly" />
          </el-form-item>
          <el-form-item class="filter-field filter-field--compact" label="开始时间">
            <el-date-picker
              v-model="queryParams.beginTime"
              type="datetime"
              placeholder="开始时间"
              value-format="YYYY-MM-DD HH:mm:ss"
            />
          </el-form-item>
          <el-form-item class="filter-field filter-field--compact" label="结束时间">
            <el-date-picker
              v-model="queryParams.endTime"
              type="datetime"
              placeholder="结束时间"
              value-format="YYYY-MM-DD HH:mm:ss"
            />
          </el-form-item>
        </el-form>
        <template #actions>
          <el-button type="primary" @click="handleQuery">搜索</el-button>
          <el-button @click="resetQuery">重置</el-button>
        </template>
      </PageToolbar>
    </template>

    <template #stats>
      <StatStrip :items="summaryItems" />
    </template>

    <el-alert
      v-if="summary.commandOnlySessions > 0"
      :title="`当前仍有 ${summary.commandOnlySessions} 个历史会话只有命令聚合数据，页面已明确标识不可回放状态。`"
      type="info"
      :closable="false"
      class="page-alert"
    />

    <el-table v-loading="loading" :data="sessionList" stripe class="audit-table" empty-text="暂无终端审计会话">
      <el-table-column label="会话ID" min-width="136">
        <template #default="{ row }">
          <el-tooltip :content="row.sessionId || '-'" placement="top">
            <span class="ops-cell-ellipsis ops-cell-ellipsis--strong">{{ formatCompactLabel(row.sessionId, 10, 6) }}</span>
          </el-tooltip>
        </template>
      </el-table-column>
      <el-table-column label="终端类型" width="80">
        <template #default="{ row }">
          <el-tag :type="sessionTypeTagType(row.sessionType)" effect="dark">{{ sessionTypeLabel(row.sessionType) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="数据来源" width="84">
        <template #default="{ row }">
          <el-tag :type="row.dataSource === 'recording' ? 'success' : 'info'" effect="dark">
            {{ row.dataSource === 'recording' ? '录制会话' : '命令聚合' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="录像状态" width="96">
        <template #default="{ row }">
          <el-tooltip :content="row.recordingWarning || row.recordingStateText || '-'" placement="top">
            <el-tag :type="recordingStateTagType(row.recordingState)" effect="dark">{{ row.recordingStateText || '-' }}</el-tag>
          </el-tooltip>
        </template>
      </el-table-column>
      <el-table-column label="目标" min-width="110" show-overflow-tooltip>
        <template #default="{ row }">{{ row.hostIp || row.hostName || '-' }}</template>
      </el-table-column>
      <el-table-column label="开始时间" prop="startTime" min-width="140" />
      <el-table-column label="命令数" prop="commandCount" width="70" />
      <el-table-column label="风险" width="80">
        <template #default="{ row }">
          <el-tag :type="riskTagType(row.riskLevel)" effect="dark">{{ riskLabel(row.riskLevel) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="最新命令" min-width="140">
        <template #default="{ row }">
          <el-tooltip :content="row.latestCommand || '-'" placement="top">
            <span class="ops-cell-ellipsis">{{ row.latestCommand || '-' }}</span>
          </el-tooltip>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="130" fixed="right">
        <template #default="{ row }">
          <div class="row-actions">
            <el-button type="primary" link @click="showDetail(row, 'commands')">详情</el-button>
            <el-button v-if="row.playbackAvailable" type="warning" link @click="showDetail(row, 'playback')">回放</el-button>
            <el-button v-if="canDownload(row)" type="success" link @click="downloadRecording(row)">下载</el-button>
          </div>
        </template>
      </el-table-column>
    </el-table>

    <template #footer>
      <el-pagination
        :current-page="queryParams.pageNum"
        :page-size="queryParams.pageSize"
        :page-sizes="[10, 20, 50, 100]"
        :total="total"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </template>

    <el-drawer
      v-model="detailVisible"
      size="920px"
      title="终端审计详情"
      class="ops-drawer ops-overlay--lg recording-detail-drawer"
      @closed="handleDetailClosed"
    >
      <div v-loading="detailLoading" class="ops-drawer__body recording-detail-drawer__body">
        <template v-if="detail.session">
          <section class="ops-drawer__section">
            <div class="ops-drawer__section-title">会话概览</div>
            <el-descriptions :column="2" border>
              <el-descriptions-item label="会话ID">{{ detail.session.sessionId }}</el-descriptions-item>
              <el-descriptions-item label="终端类型">{{ sessionTypeLabel(detail.session.sessionType) }}</el-descriptions-item>
              <el-descriptions-item label="录像状态">{{ detail.session.recordingStateText || '-' }}</el-descriptions-item>
              <el-descriptions-item label="持续时长">{{ formatDuration(detail.session.duration) }}</el-descriptions-item>
              <el-descriptions-item label="操作用户">{{ detail.session.username || detail.session.sshUser || '-' }}</el-descriptions-item>
              <el-descriptions-item label="目标">{{ detail.session.hostIp || detail.session.hostName || '-' }}</el-descriptions-item>
              <el-descriptions-item label="录制文件" :span="2">
                <el-tooltip :content="detail.session.filePath || '-'" placement="top">
                  <span class="ops-cell-ellipsis">{{ detail.session.filePath || '-' }}</span>
                </el-tooltip>
              </el-descriptions-item>
              <el-descriptions-item label="录制大小">
                {{ formatFileSize(detail.session.actualFileSize || detail.session.fileSize) }}
              </el-descriptions-item>
            </el-descriptions>
          </section>

          <section v-if="playback.health && playback.health.message" class="ops-drawer__section">
            <div class="ops-drawer__section-title">录制状态</div>
            <el-alert
              :title="playback.health.message"
              :type="playback.health.canPlayback ? 'success' : (detail.session.dataSource === 'recording' ? 'warning' : 'info')"
              :closable="false"
            />
          </section>

          <section class="ops-drawer__section">
            <div class="ops-drawer__section-title">快捷动作</div>
            <div class="detail-actions">
              <el-button v-if="playback.health && playback.health.canPlayback" type="warning" @click="activeDetailTab = 'playback'">打开回放</el-button>
              <el-button v-if="canDownload(detail.session)" type="success" @click="downloadRecording(detail.session)">下载原始日志</el-button>
            </div>
          </section>

          <section class="ops-drawer__section">
            <div class="ops-drawer__section-title">审计详情</div>
            <el-tabs v-model="activeDetailTab">
              <el-tab-pane label="命令审计" name="commands">
                <el-table :data="detail.commands" stripe max-height="520">
                  <el-table-column label="#" prop="sequence" width="70" />
                  <el-table-column label="执行时间" prop="executeTime" min-width="170" />
                  <el-table-column label="耗时秒" width="100">
                    <template #default="{ row }">{{ formatElapsed(row.elapsedSeconds) }}</template>
                  </el-table-column>
                  <el-table-column label="命令内容" prop="command" min-width="360" show-overflow-tooltip />
                  <el-table-column label="风险" width="88">
                    <template #default="{ row }">
                      <el-tag :type="riskTagType(row.riskLevel)" effect="dark">{{ row.isSensitive ? '敏感' : riskLabel(row.riskLevel) }}</el-tag>
                    </template>
                  </el-table-column>
                  <el-table-column label="风险原因" prop="riskReason" min-width="180" show-overflow-tooltip />
                </el-table>
              </el-tab-pane>

              <el-tab-pane label="录像回放" name="playback">
                <div v-loading="playbackLoading">
                  <template v-if="playback.health && playback.health.canPlayback">
                    <div class="playback-stats">
                      <el-tag type="info">事件 {{ playback.stats.totalEvents || 0 }}</el-tag>
                      <el-tag type="warning">输入 {{ playback.stats.inputEvents || 0 }}</el-tag>
                      <el-tag type="success">输出 {{ playback.stats.outputEvents || 0 }}</el-tag>
                      <el-tag>窗口变化 {{ playback.stats.resizeEvents || 0 }}</el-tag>
                    </div>

                    <div class="playback-toolbar">
                      <el-input v-model="playbackQuery.keyword" clearable placeholder="检索录像内容" style="width: 260px" @keyup.enter="handlePlaybackQuery" />
                      <el-select v-model="playbackQuery.eventType" clearable placeholder="全部事件" style="width: 140px">
                        <el-option label="输入" value="input" />
                        <el-option label="输出" value="output" />
                        <el-option label="窗口变化" value="resize" />
                        <el-option label="系统事件" value="system" />
                      </el-select>
                      <el-button type="primary" @click="handlePlaybackQuery">检索</el-button>
                      <el-button @click="resetPlaybackQuery">重置</el-button>
                      <el-button type="success" :disabled="!playback.events.length" @click="startPlayback">播放</el-button>
                      <el-button :disabled="!player.playing" @click="pausePlayback">暂停</el-button>
                      <el-button :disabled="!playback.events.length" @click="stepPlayback">单步</el-button>
                      <el-button :disabled="!playback.events.length" @click="resetPlayer">重置</el-button>
                    </div>

                    <div class="player-box">
                      <div class="player-meta">
                        <span>进度 {{ player.cursor }}/{{ playback.events.length }}</span>
                        <span>{{ playbackProgress }}%</span>
                      </div>
                      <el-progress :percentage="playbackProgress" :stroke-width="10" />
                      <div ref="playbackScreen" class="player-screen">{{ player.screen || '点击播放或单步开始回放当前页事件。' }}</div>
                    </div>

                    <el-table :data="playback.events" stripe max-height="360">
                      <el-table-column label="行号" prop="line" width="80" />
                      <el-table-column label="时间" prop="at" min-width="170" />
                      <el-table-column label="事件类型" width="90">
                        <template #default="{ row }">
                          <el-tag :type="playbackEventTagType(row.eventType)" effect="dark">{{ row.eventTypeText }}</el-tag>
                        </template>
                      </el-table-column>
                      <el-table-column label="相对秒" width="100">
                        <template #default="{ row }">{{ formatElapsed(row.relativeSeconds) }}</template>
                      </el-table-column>
                      <el-table-column label="内容" min-width="420">
                        <template #default="{ row }">
                          <div class="event-content" :class="{ matched: row.matched }">{{ row.content || '-' }}</div>
                        </template>
                      </el-table-column>
                    </el-table>

                    <div class="playback-pagination">
                      <el-pagination
                        :current-page="playbackQuery.pageNum"
                        :page-size="playbackQuery.pageSize"
                        :page-sizes="[50, 100, 200]"
                        :total="playback.total || 0"
                        layout="total, sizes, prev, pager, next, jumper"
                        @size-change="handlePlaybackSizeChange"
                        @current-change="handlePlaybackCurrentChange"
                      />
                    </div>
                  </template>
                  <el-empty v-else description="当前会话暂无可回放录像，可继续查看命令审计或下载原始日志。" />
                </div>
              </el-tab-pane>
            </el-tabs>
          </section>
        </template>
      </div>
    </el-drawer>
  </TablePage>
</template>

<script>
import storage from '@/utils/storage'
import PageHeader from '@/components/platform/PageHeader.vue'
import PageIntro from '@/components/platform/PageIntro.vue'
import PageToolbar from '@/components/platform/PageToolbar.vue'
import StatStrip from '@/components/platform/StatStrip.vue'
import TablePage from '@/components/platform/TablePage.vue'

const createEmptyPlayback = () => ({
  health: {
    state: '',
    stateText: '',
    message: '',
    canPlayback: false,
    canDownload: false,
    fileExists: false,
    fileReadable: false,
    isEmpty: false,
    recordedFileSize: 0,
    actualFileSize: 0,
    sizeMismatch: false
  },
  stats: {
    totalEvents: 0,
    inputEvents: 0,
    outputEvents: 0,
    resizeEvents: 0,
    systemEvents: 0,
    matchedEvents: 0
  },
  events: [],
  total: 0
})

const createDefaultQueryParams = () => ({
  sessionId: '',
  hostId: '',
  hostKeyword: '',
  keyword: '',
  riskLevel: '',
  sensitiveOnly: false,
  beginTime: '',
  endTime: '',
  pageNum: 1,
  pageSize: 10
})

export default {
  name: 'TerminalAuditRecording',
  components: {
    PageHeader,
    PageIntro,
    PageToolbar,
    StatStrip,
    TablePage
  },
  data() {
    return {
      loading: false,
      detailLoading: false,
      playbackLoading: false,
      detailVisible: false,
      activeDetailTab: 'commands',
      total: 0,
      sessionList: [],
      summary: {
        totalSessions: 0,
        recordedSessions: 0,
        commandOnlySessions: 0,
        totalCommands: 0,
        sensitiveCommands: 0,
        riskySessions: 0
      },
      detail: { session: null, commands: [] },
      playback: createEmptyPlayback(),
      playbackQuery: {
        keyword: '',
        eventType: '',
        pageNum: 1,
        pageSize: 100
      },
      player: {
        timer: null,
        playing: false,
        cursor: 0,
        screen: ''
      },
      queryParams: createDefaultQueryParams()
    }
  },
  watch: {
    '$route.query': {
      handler(query, oldQuery) {
        if (JSON.stringify(query || {}) === JSON.stringify(oldQuery || {})) {
          return
        }
        this.applyRouteQuery(query)
        this.refreshAll()
      }
    }
  },
  computed: {
    summaryItems() {
      return [
        { label: '会话总数', value: this.summary.totalSessions, hint: '当前审计范围内的会话规模', tone: 'primary' },
        { label: '可回放录像', value: this.summary.recordedSessions, hint: '支持直接进入录像回放', tone: 'success' },
        { label: '敏感命令', value: this.summary.sensitiveCommands, hint: '需要优先复核的命令轨迹', tone: 'warning' },
        { label: '风险会话', value: this.summary.riskySessions, hint: '建议先看详情与原始日志', tone: 'danger' }
      ]
    },
    playbackProgress() {
      const total = this.playback.events.length || 0
      return total ? Math.round((this.player.cursor / total) * 100) : 0
    }
  },
  methods: {
    formatCompactLabel(value, prefixLength = 10, suffixLength = 6) {
      const text = String(value || '')
      if (!text) return '-'
      if (text.length <= prefixLength + suffixLength) {
        return text
      }
      return `${text.slice(0, prefixLength)}...${text.slice(-suffixLength)}`
    },
    buildQueryParams() {
      const params = {
        pageNum: this.queryParams.pageNum,
        pageSize: this.queryParams.pageSize
      }
      if (this.queryParams.sessionId) params.sessionId = this.queryParams.sessionId
      if (this.queryParams.hostId) params.hostId = this.queryParams.hostId
      if (this.queryParams.hostKeyword) params.hostKeyword = this.queryParams.hostKeyword
      if (this.queryParams.keyword) params.keyword = this.queryParams.keyword
      if (this.queryParams.riskLevel !== '' && this.queryParams.riskLevel !== null && this.queryParams.riskLevel !== undefined) params.riskLevel = this.queryParams.riskLevel
      if (this.queryParams.sensitiveOnly) params.sensitiveOnly = true
      if (this.queryParams.beginTime) params.beginTime = this.queryParams.beginTime
      if (this.queryParams.endTime) params.endTime = this.queryParams.endTime
      return params
    },
    applyRouteQuery(query = {}) {
      const nextQuery = createDefaultQueryParams()
      nextQuery.sessionId = String(query.sessionId || '')
      nextQuery.hostId = String(query.hostId || '')
      nextQuery.hostKeyword = String(query.hostKeyword || query.hostIp || '')
      nextQuery.keyword = String(query.keyword || '')
      nextQuery.beginTime = String(query.beginTime || '')
      nextQuery.endTime = String(query.endTime || '')

      if (query.riskLevel !== undefined && query.riskLevel !== '') {
        const parsedRiskLevel = Number(query.riskLevel)
        nextQuery.riskLevel = Number.isFinite(parsedRiskLevel) ? parsedRiskLevel : ''
      }

      nextQuery.sensitiveOnly = query.sensitiveOnly === true || query.sensitiveOnly === 'true' || query.sensitiveOnly === '1'

      const pageNum = Number(query.pageNum)
      if (Number.isFinite(pageNum) && pageNum > 0) {
        nextQuery.pageNum = pageNum
      }

      const pageSize = Number(query.pageSize)
      if (Number.isFinite(pageSize) && pageSize > 0) {
        nextQuery.pageSize = pageSize
      }

      this.queryParams = nextQuery
    },
    buildRouteQuery() {
      const query = {}
      if (this.queryParams.sessionId) query.sessionId = String(this.queryParams.sessionId)
      if (this.queryParams.hostId) query.hostId = String(this.queryParams.hostId)
      if (this.queryParams.hostKeyword) query.hostKeyword = String(this.queryParams.hostKeyword)
      if (this.queryParams.keyword) query.keyword = String(this.queryParams.keyword)
      if (this.queryParams.riskLevel !== '' && this.queryParams.riskLevel !== null && this.queryParams.riskLevel !== undefined) {
        query.riskLevel = String(this.queryParams.riskLevel)
      }
      if (this.queryParams.sensitiveOnly) query.sensitiveOnly = 'true'
      if (this.queryParams.beginTime) query.beginTime = this.queryParams.beginTime
      if (this.queryParams.endTime) query.endTime = this.queryParams.endTime
      if (Number(this.queryParams.pageNum) > 1) query.pageNum = String(this.queryParams.pageNum)
      if (Number(this.queryParams.pageSize) !== 10) query.pageSize = String(this.queryParams.pageSize)
      return query
    },
    refreshByRouteQuery() {
      const nextQuery = this.buildRouteQuery()
      const currentQuery = this.$route.query || {}
      if (JSON.stringify(currentQuery) === JSON.stringify(nextQuery)) {
        this.refreshAll()
        return
      }
      this.$router.replace({ path: this.$route.path, query: nextQuery }).catch(() => {})
    },
    buildPlaybackQueryParams() {
      const params = {
        pageNum: this.playbackQuery.pageNum,
        pageSize: this.playbackQuery.pageSize
      }
      if (this.playbackQuery.keyword) params.keyword = this.playbackQuery.keyword
      if (this.playbackQuery.eventType) params.eventType = this.playbackQuery.eventType
      return params
    },
    async fetchSummary() {
      const { data: res } = await this.$api.getTerminalAuditSummary()
      if (res.code !== 200) throw new Error(res.message || '获取终端审计汇总失败')
      this.summary = res.data || this.summary
    },
    async fetchSessions() {
      const { data: res } = await this.$api.queryTerminalAuditSessionList(this.buildQueryParams())
      if (res.code !== 200) throw new Error(res.message || '获取终端审计会话列表失败')
      this.sessionList = res.data?.list || []
      this.total = res.data?.total || 0
    },
    async refreshAll() {
      this.loading = true
      try {
        await Promise.all([this.fetchSummary(), this.fetchSessions()])
      } catch (error) {
        this.$message.error(error.message || '刷新终端审计数据失败')
      } finally {
        this.loading = false
      }
    },
    handleQuery() {
      this.queryParams.pageNum = 1
      this.refreshByRouteQuery()
    },
    resetQuery() {
      this.queryParams = createDefaultQueryParams()
      this.refreshByRouteQuery()
    },
    handleSizeChange(size) {
      this.queryParams.pageSize = size
      this.queryParams.pageNum = 1
      this.refreshByRouteQuery()
    },
    handleCurrentChange(page) {
      this.queryParams.pageNum = page
      this.refreshByRouteQuery()
    },
    resetPlaybackState() {
      this.stopPlayback()
      this.playback = createEmptyPlayback()
      this.playbackQuery = { keyword: '', eventType: '', pageNum: 1, pageSize: 100 }
      this.player.cursor = 0
      this.player.screen = ''
    },
    async fetchPlayback(sessionId) {
      this.playbackLoading = true
      try {
        const { data: res } = await this.$api.getTerminalAuditSessionPlayback(sessionId, this.buildPlaybackQueryParams())
        if (res.code !== 200) throw new Error(res.message || '获取录像回放失败')
        this.playback = res.data || createEmptyPlayback()
        this.resetPlayer()
      } catch (error) {
        this.resetPlaybackState()
        this.$message.error(error.message || '获取录像回放失败')
      } finally {
        this.playbackLoading = false
      }
    },
    async showDetail(row, preferredTab = 'commands') {
      this.detailVisible = true
      this.detailLoading = true
      this.activeDetailTab = preferredTab
      this.detail = { session: null, commands: [] }
      this.resetPlaybackState()
      try {
        const { data: res } = await this.$api.getTerminalAuditSessionDetail(row.sessionId)
        if (res.code !== 200) throw new Error(res.message || '获取会话详情失败')
        this.detail = res.data || this.detail
        await this.fetchPlayback(row.sessionId)
        if (preferredTab === 'playback' && !(this.playback.health && this.playback.health.canPlayback)) {
          this.activeDetailTab = 'commands'
        }
      } catch (error) {
        this.$message.error(error.message || '获取会话详情失败')
      } finally {
        this.detailLoading = false
      }
    },
    handleDetailClosed() {
      this.resetPlaybackState()
      this.detail = { session: null, commands: [] }
      this.activeDetailTab = 'commands'
    },
    handlePlaybackQuery() {
      if (!this.detail.session) return
      this.playbackQuery.pageNum = 1
      this.fetchPlayback(this.detail.session.sessionId)
    },
    resetPlaybackQuery() {
      if (!this.detail.session) return
      this.playbackQuery = { keyword: '', eventType: '', pageNum: 1, pageSize: 100 }
      this.fetchPlayback(this.detail.session.sessionId)
    },
    handlePlaybackSizeChange(size) {
      if (!this.detail.session) return
      this.playbackQuery.pageSize = size
      this.playbackQuery.pageNum = 1
      this.fetchPlayback(this.detail.session.sessionId)
    },
    handlePlaybackCurrentChange(page) {
      if (!this.detail.session) return
      this.playbackQuery.pageNum = page
      this.fetchPlayback(this.detail.session.sessionId)
    },
    canDownload(row) {
      return row && row.dataSource === 'recording' && row.storageType === 1 && row.filePath
    },
    downloadRecording(row) {
      const token = storage.getItem('token')
      const tokenValue = typeof token === 'object' ? token.access_token || token.token : token
      window.open(this.$api.getTerminalAuditRecordingDownloadUrl(row.sessionId, tokenValue), '_blank')
    },
    renderPlaybackScreen(limit = this.player.cursor) {
      const parts = []
      for (let index = 0; index < limit; index++) {
        const event = this.playback.events[index]
        if (!event) continue
        if (event.eventType === 'input') {
          parts.push(`\n$ ${String(event.content || '').replace(/\s+$/, '')}\n`)
        } else if (event.eventType === 'output') {
          parts.push(event.content || '')
          if (event.content && !String(event.content).endsWith('\n')) parts.push('\n')
        } else if (event.eventType === 'resize') {
          parts.push(`\n[窗口变化] ${event.content}\n`)
        } else {
          parts.push(`\n[${event.eventTypeText}] ${event.content}\n`)
        }
      }
      return parts.join('').trimStart()
    },
    scrollPlaybackScreen() {
      this.$nextTick(() => {
        const screen = this.$refs.playbackScreen
        if (screen) screen.scrollTop = screen.scrollHeight
      })
    },
    startPlayback() {
      if (!this.playback.events.length) return
      this.stopPlayback()
      this.player.playing = true
      this.player.timer = setInterval(() => {
        if (this.player.cursor >= this.playback.events.length) {
          this.stopPlayback()
          return
        }
        this.player.cursor += 1
        this.player.screen = this.renderPlaybackScreen()
        this.scrollPlaybackScreen()
      }, 700)
    },
    pausePlayback() {
      this.stopPlayback()
    },
    stepPlayback() {
      if (!this.playback.events.length || this.player.cursor >= this.playback.events.length) return
      this.stopPlayback()
      this.player.cursor += 1
      this.player.screen = this.renderPlaybackScreen()
      this.scrollPlaybackScreen()
    },
    resetPlayer() {
      this.stopPlayback()
      this.player.cursor = 0
      this.player.screen = ''
    },
    stopPlayback() {
      if (this.player.timer) clearInterval(this.player.timer)
      this.player.timer = null
      this.player.playing = false
    },
    riskLabel(level) {
      return { 0: '低风险', 1: '中风险', 2: '高风险' }[Number(level)] || '低风险'
    },
    riskTagType(level) {
      return { 0: 'success', 1: 'warning', 2: 'danger' }[Number(level)] || 'info'
    },
    recordingStateTagType(state) {
      return { ready: 'success', command_only: 'info', empty: 'warning', missing: 'danger', read_error: 'danger', unsupported: 'info' }[state] || 'info'
    },
    sessionTypeLabel(type) {
      return { ssh: 'SSH', pod: 'Pod', kubectl: 'kubectl', unknown: '未知' }[type] || '未知'
    },
    sessionTypeTagType(type) {
      return { ssh: 'primary', pod: 'success', kubectl: 'warning' }[type] || 'info'
    },
    playbackEventTagType(type) {
      return { input: 'warning', output: 'success', resize: 'info', system: 'primary' }[type] || 'info'
    },
    formatDuration(seconds) {
      const value = Number(seconds || 0)
      if (!value) return '0秒'
      if (value < 60) return `${value}秒`
      const minutes = Math.floor(value / 60)
      const remainSeconds = value % 60
      return remainSeconds ? `${minutes}分 ${remainSeconds}秒` : `${minutes}分`
    },
    formatElapsed(seconds) {
      const value = Number(seconds || 0)
      return value || value === 0 ? value.toFixed(2) : '-'
    },
    formatFileSize(size) {
      const value = Number(size || 0)
      if (!value) return '0 B'
      if (value < 1024) return `${value} B`
      if (value < 1024 * 1024) return `${(value / 1024).toFixed(1)} KB`
      return `${(value / (1024 * 1024)).toFixed(2)} MB`
    }
  },
  created() {
    this.applyRouteQuery(this.$route.query)
    this.refreshAll()
  },
  beforeUnmount() {
    this.stopPlayback()
  }
}
</script>

<style scoped>
.terminal-audit-page :deep(.page-actions) {
  align-items: center;
}

.terminal-audit-page :deep(.filter-field .el-form-item__content),
.terminal-audit-page :deep(.filter-field .el-select),
.terminal-audit-page :deep(.filter-field .el-input),
.terminal-audit-page :deep(.filter-field .el-date-editor) {
  width: 100%;
}

.terminal-audit-page :deep(.audit-table .cell) {
  min-width: 0;
}

.row-actions,
.detail-actions,
.playback-toolbar,
.player-meta {
  display: flex;
  gap: 12px;
  align-items: center;
}

.recording-detail-drawer__body {
  min-height: 240px;
}

.recording-detail-drawer :deep(.el-descriptions__table) {
  width: 100%;
  table-layout: fixed;
}

.recording-detail-drawer :deep(.el-descriptions__cell) {
  word-break: break-word;
}

.playback-stats {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.playback-pagination {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}

.player-box {
  padding: 16px;
  border-radius: 16px;
  background: rgba(2, 6, 23, 0.72);
  border: 1px solid var(--border-subtle);
}

.player-meta {
  justify-content: space-between;
  color: #cbd5e1;
  font-size: 13px;
}

.player-screen {
  margin-top: 12px;
  min-height: 220px;
  max-height: 320px;
  overflow: auto;
  padding: 16px;
  border-radius: 12px;
  background: #020617;
  color: #d7f7ff;
  white-space: pre-wrap;
  word-break: break-word;
  font-family: Consolas, Monaco, monospace;
  line-height: 1.5;
}

.event-content {
  white-space: pre-wrap;
  word-break: break-word;
  line-height: 1.5;
}

.event-content.matched {
  background: rgba(250, 204, 21, 0.18);
  border-radius: 8px;
  padding: 6px 8px;
}

@media (max-width: 960px) {
  .row-actions,
  .detail-actions,
  .playback-toolbar,
  .player-meta {
    display: grid;
  }
}
</style>
