<template>
  <TablePage
    class="alert-history-table-page"
    section-title="Webhook 历史"
    section-subtitle="查看告警推送记录、状态与失败原因，辅助排查通知链路"
    wide
  >
    <template #header>
      <PageHeader
        eyebrow="Monitoring"
        title="Webhook 历史"
        subtitle="追踪 webhook 推送结果、来源分布与异常原因"
      >
        <template #actions>
          <el-button type="primary" :loading="loading" @click="refreshAll">
            <el-icon><Refresh /></el-icon>
            刷新历史
          </el-button>
        </template>
        <template #meta>
          <div class="page-chip-row">
            <span class="platform-chip">最近告警 {{ summary.latestAlertTime || '-' }}</span>
            <span class="platform-chip">已恢复 {{ summary.resolvedIncidents }}</span>
            <span class="platform-chip">历史总数 {{ total }}</span>
          </div>
        </template>
        <template #intro>
          <PageIntro
            title="排障说明"
            text="支持按来源、等级与推送状态筛选，快速定位失败 webhook 与告警上下文。"
          />
        </template>
      </PageHeader>
    </template>

    <template #toolbar>
      <PageToolbar>
        <el-form :model="queryParams" :inline="true" class="filter-cluster">
          <el-form-item class="filter-field" label="关键词">
            <el-input
              v-model="queryParams.keyword"
              placeholder="搜索标题、内容或错误信息"
              clearable
              @keyup.enter="handleQuery"
            />
          </el-form-item>
          <el-form-item class="filter-field" label="来源">
            <el-select v-model="queryParams.source" placeholder="全部来源" clearable>
              <el-option label="Prometheus" value="prometheus" />
              <el-option label="Zabbix" value="zabbix" />
              <el-option label="Deployment" value="deployment" />
            </el-select>
          </el-form-item>
          <el-form-item class="filter-field" label="等级">
            <el-select v-model="queryParams.level" placeholder="全部等级" clearable>
              <el-option label="critical" value="critical" />
              <el-option label="warning" value="warning" />
              <el-option label="info" value="info" />
              <el-option label="P1" value="P1" />
              <el-option label="P2" value="P2" />
              <el-option label="P3" value="P3" />
              <el-option label="P4" value="P4" />
            </el-select>
          </el-form-item>
          <el-form-item class="filter-field" label="推送状态">
            <el-select v-model="queryParams.status" placeholder="全部状态" clearable>
              <el-option label="success" value="success" />
              <el-option label="failed" value="failed" />
              <el-option label="partial" value="partial" />
            </el-select>
          </el-form-item>
        </el-form>
        <template #actions>
          <el-button type="primary" @click="handleQuery">
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
      :data="webhookList"
      stripe
      class="alert-history-table"
      empty-text="暂无 webhook 记录"
    >
      <el-table-column label="推送时间" prop="createdAt" min-width="170" />
      <el-table-column label="来源" width="130">
        <template #default="{ row }">
          <el-tag :type="sourceTagType(row.source)">
            {{ sourceLabel(row.source) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="等级" width="110">
        <template #default="{ row }">
          <el-tag :type="levelTagType(row.level)" effect="dark">
            {{ row.level || '-' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="标题" prop="title" min-width="260" show-overflow-tooltip />
      <el-table-column label="推送状态" width="120">
        <template #default="{ row }">
          <el-tag :type="historyStatusTagType(row.status)" effect="dark">
            {{ row.status || '-' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="推送结果" min-width="120">
        <template #default="{ row }">
          {{ row.successCount }}/{{ row.notifyCount }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="120" fixed="right">
        <template #default="{ row }">
          <el-button type="primary" link @click="showDetail(row)">查看详情</el-button>
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
  </TablePage>

  <el-drawer
    v-model="detailVisible"
    title="Webhook 详情"
    size="54%"
    :destroy-on-close="true"
  >
    <div v-if="activeLog" class="alert-history-detail">
      <div class="alert-history-detail__meta">
        <div class="detail-meta-card">
          <span class="detail-meta-card__label">告警来源</span>
          <span class="detail-meta-card__value">{{ sourceLabel(activeLog.source) }}</span>
        </div>
        <div class="detail-meta-card">
          <span class="detail-meta-card__label">等级</span>
          <span class="detail-meta-card__value">{{ activeLog.level || '-' }}</span>
        </div>
        <div class="detail-meta-card">
          <span class="detail-meta-card__label">推送状态</span>
          <span class="detail-meta-card__value">{{ activeLog.status || '-' }}</span>
        </div>
        <div class="detail-meta-card">
          <span class="detail-meta-card__label">推送结果</span>
          <span class="detail-meta-card__value">{{ activeLog.successCount }}/{{ activeLog.notifyCount }}</span>
        </div>
      </div>

      <section class="detail-block">
        <div class="detail-block__title">标题</div>
        <div class="detail-block__text">{{ activeLog.title || '-' }}</div>
      </section>

      <section class="detail-block">
        <div class="detail-block__title">内容</div>
        <pre class="detail-block__pre">{{ activeLog.content || '-' }}</pre>
      </section>

      <section class="detail-block detail-block--grid">
        <div>
          <div class="detail-block__title">标签</div>
          <pre class="detail-block__pre">{{ activeLog.tags || '-' }}</pre>
        </div>
        <div>
          <div class="detail-block__title">扩展信息</div>
          <pre class="detail-block__pre">{{ activeLog.extra || '-' }}</pre>
        </div>
      </section>

      <section class="detail-block">
        <div class="detail-block__title">错误信息</div>
        <pre class="detail-block__pre">{{ activeLog.errorMsg || '-' }}</pre>
      </section>
    </div>
    <el-empty v-else description="暂无详情数据" />
  </el-drawer>
</template>

<script>
import PageHeader from '@/components/platform/PageHeader.vue'
import PageIntro from '@/components/platform/PageIntro.vue'
import PageToolbar from '@/components/platform/PageToolbar.vue'
import StatStrip from '@/components/platform/StatStrip.vue'
import TablePage from '@/components/platform/TablePage.vue'

const createEmptyQuery = () => ({
  keyword: '',
  source: '',
  level: '',
  status: '',
  pageNum: 1,
  pageSize: 10
})

export default {
  name: 'AlertHistory',
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
      detailVisible: false,
      total: 0,
      webhookList: [],
      activeLog: null,
      summary: {
        totalIncidents: 0,
        openIncidents: 0,
        processingIncidents: 0,
        resolvedIncidents: 0,
        totalWebhookLogs: 0,
        criticalWebhookLogs: 0,
        totalNotifyLogs: 0,
        successfulNotifyLogs: 0,
        failedNotifyLogs: 0,
        totalNotifyRobots: 0,
        enabledNotifyRobots: 0,
        totalAlertSources: 0,
        enabledAlertSources: 0,
        latestAlertTime: '',
        latestNotifyTime: ''
      },
      queryParams: createEmptyQuery()
    }
  },
  computed: {
    statItems() {
      return [
        { label: 'Webhook 总量', value: this.summary.totalWebhookLogs, hint: '最近 webhook 记录', tone: 'primary' },
        { label: '高危推送', value: this.summary.criticalWebhookLogs, hint: 'critical / P1 / P2', tone: 'danger' },
        { label: '最近告警', value: this.summary.latestAlertTime || '-', hint: '最近一次 webhook', tone: 'warning' },
        { label: '已恢复', value: this.summary.resolvedIncidents, hint: '已恢复 incident', tone: 'success' }
      ]
    }
  },
  methods: {
    buildQueryParams() {
      const params = {
        pageNum: this.queryParams.pageNum,
        pageSize: this.queryParams.pageSize
      }
      if (this.queryParams.keyword) params.keyword = this.queryParams.keyword
      if (this.queryParams.source) params.source = this.queryParams.source
      if (this.queryParams.level) params.level = this.queryParams.level
      if (this.queryParams.status) params.status = this.queryParams.status
      return params
    },
    async fetchSummary() {
      const { data: res } = await this.$api.getMonitorAlertSummary()
      if (res.code !== 200) {
        throw new Error(res.message || '获取告警概览失败')
      }
      this.summary = res.data || this.summary
    },
    async fetchWebhookList() {
      const { data: res } = await this.$api.queryMonitorAlertWebhookList(this.buildQueryParams())
      if (res.code !== 200) {
        throw new Error(res.message || '获取 webhook 历史失败')
      }
      const data = res.data || {}
      this.webhookList = data.list || []
      this.total = data.total || 0
    },
    async refreshAll() {
      this.loading = true
      try {
        await Promise.all([this.fetchSummary(), this.fetchWebhookList()])
      } catch (error) {
        this.$message.error(error.message || '刷新 webhook 历史失败')
      } finally {
        this.loading = false
      }
    },
    handleQuery() {
      this.queryParams.pageNum = 1
      this.fetchWebhookList().catch(error => {
        this.$message.error(error.message || '获取 webhook 历史失败')
      })
    },
    resetQuery() {
      this.queryParams = createEmptyQuery()
      this.fetchWebhookList().catch(error => {
        this.$message.error(error.message || '获取 webhook 历史失败')
      })
    },
    handleSizeChange(size) {
      this.queryParams.pageSize = size
      this.queryParams.pageNum = 1
      this.fetchWebhookList().catch(error => {
        this.$message.error(error.message || '获取 webhook 历史失败')
      })
    },
    handleCurrentChange(page) {
      this.queryParams.pageNum = page
      this.fetchWebhookList().catch(error => {
        this.$message.error(error.message || '获取 webhook 历史失败')
      })
    },
    showDetail(row) {
      this.activeLog = row
      this.detailVisible = true
    },
    sourceLabel(source) {
      const textMap = {
        prometheus: 'Prometheus',
        zabbix: 'Zabbix',
        deployment: 'Deployment'
      }
      return textMap[source] || source || '未知来源'
    },
    sourceTagType(source) {
      const typeMap = {
        prometheus: 'success',
        zabbix: 'warning',
        deployment: 'info'
      }
      return typeMap[source] || 'info'
    },
    levelTagType(level) {
      const typeMap = {
        critical: 'danger',
        warning: 'warning',
        info: 'info',
        P1: 'danger',
        P2: 'warning',
        P3: 'primary',
        P4: 'info'
      }
      return typeMap[level] || 'info'
    },
    historyStatusTagType(status) {
      const typeMap = {
        success: 'success',
        failed: 'danger',
        partial: 'warning'
      }
      return typeMap[status] || 'info'
    }
  },
  created() {
    this.refreshAll()
  }
}
</script>

<style scoped>
.alert-history-table-page :deep(.page-actions) {
  align-items: center;
}

.alert-history-detail {
  display: grid;
  gap: 18px;
}

.alert-history-table-page :deep(.el-table .cell.el-tooltip) {
  white-space: normal;
  word-break: break-word;
  line-height: 1.5;
}

.alert-history-detail__meta {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  gap: 14px;
}

.detail-meta-card,
.detail-block {
  border-radius: 18px;
  border: 1px solid var(--border-subtle);
  background: rgba(255, 255, 255, 0.04);
}

.detail-meta-card {
  padding: 16px 18px;
  display: grid;
  gap: 8px;
}

.detail-meta-card__label,
.detail-block__title {
  font-size: 12px;
  font-weight: 600;
  color: var(--text-muted);
}

.detail-meta-card__value {
  font-size: 16px;
  font-weight: 700;
  color: var(--text-primary);
  word-break: break-word;
}

.detail-block {
  padding: 18px;
}

.detail-block--grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 16px;
}

.detail-block__text,
.detail-block__pre {
  margin-top: 10px;
  color: var(--text-secondary);
  white-space: pre-wrap;
  word-break: break-word;
}

.detail-block__pre {
  min-height: 96px;
  padding: 14px;
  border-radius: 14px;
  background: rgba(2, 6, 23, 0.44);
  border: 1px solid rgba(148, 163, 184, 0.12);
  font-family: 'Consolas', 'Courier New', monospace;
}

@media (max-width: 960px) {
  .detail-block--grid {
    grid-template-columns: 1fr;
  }
}
</style>
