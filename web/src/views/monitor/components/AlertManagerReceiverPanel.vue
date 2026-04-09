<template>
  <div class="alertmanager-receiver-panel">
    <div class="subpanel-header">
      <div>
        <div class="subpanel-title">AlertManager 接收器</div>
        <div class="subpanel-subtitle">读取真实 AlertManager 当前 receiver 配置，辅助核对告警路由与分发链路。</div>
      </div>
      <el-button size="small" icon="Refresh" @click="reload">刷新</el-button>
    </div>

    <el-alert
      v-if="loadError"
      :title="loadError"
      type="error"
      :closable="false"
      class="subpanel-alert"
    />

    <el-alert
      v-else-if="!alertManagerSources.length"
      title="暂无已启用的 AlertManager 告警源。"
      type="info"
      :closable="false"
      class="subpanel-alert"
    />

    <template v-else>
      <div class="selector-row">
        <span class="selector-label">联调源</span>
        <el-select v-model="selectedSourceId" style="width: 260px" filterable @change="handleSourceChange">
          <el-option
            v-for="source in alertManagerSources"
            :key="source.id"
            :label="`${source.name} (${source.apiBaseUrl || '-'})`"
            :value="source.id"
          />
        </el-select>
        <el-tag :type="statusTagType" effect="dark">{{ statusText }}</el-tag>
      </div>

      <el-table
        v-loading="loading"
        :data="receiverList"
        stripe
        class="data-table"
        empty-text="当前没有 receiver"
      >
        <el-table-column label="名称" prop="name" min-width="160" show-overflow-tooltip />
        <el-table-column label="生效状态" width="96">
          <template #default="{ row }">
            <el-tag :type="row.active ? 'success' : 'info'">
              {{ row.active ? 'active' : 'inactive' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="集成通道" min-width="220" show-overflow-tooltip>
          <template #default="{ row }">
            {{ formatIntegrations(row.integrations) }}
          </template>
        </el-table-column>
      </el-table>
    </template>
  </div>
</template>

<script>
export default {
  name: 'AlertManagerReceiverPanel',
  data() {
    return {
      loading: false,
      loadError: '',
      alertManagerSources: [],
      selectedSourceId: null,
      statusAvailable: false,
      receiverList: []
    }
  },
  computed: {
    statusText() {
      return this.statusAvailable ? '连接正常' : '待检测'
    },
    statusTagType() {
      return this.statusAvailable ? 'success' : 'info'
    }
  },
  methods: {
    async reload() {
      this.loadError = ''
      await this.loadSources()
      if (!this.selectedSourceId) return

      this.loading = true
      try {
        await this.fetchStatus()
        await this.fetchReceivers()
      } catch (error) {
        this.loadError = error.message || '加载 AlertManager 接收器失败'
      } finally {
        this.loading = false
      }
    },
    async loadSources() {
      const { data: res } = await this.$api.getMonitorAlertSourceList()
      if (res.code !== 200) {
        throw new Error(res.message || '获取告警源列表失败')
      }
      const sourceList = Array.isArray(res.data) ? res.data : []
      this.alertManagerSources = sourceList.filter(source => Number(source.type) === 4 && Number(source.status) === 1)
      if (!this.alertManagerSources.length) {
        this.selectedSourceId = null
        this.receiverList = []
        this.statusAvailable = false
        return
      }
      const exists = this.alertManagerSources.some(source => source.id === this.selectedSourceId)
      if (!exists) {
        this.selectedSourceId = this.alertManagerSources[0].id
      }
    },
    async fetchStatus() {
      const { data: res } = await this.$api.getMonitorAlertManagerStatus({ sourceId: this.selectedSourceId })
      if (res.code !== 200) {
        throw new Error(res.message || '获取 AlertManager 状态失败')
      }
      this.statusAvailable = !!res.data?.available
    },
    async fetchReceivers() {
      const { data: res } = await this.$api.getMonitorAlertManagerReceivers({ sourceId: this.selectedSourceId })
      if (res.code !== 200) {
        throw new Error(res.message || '获取 AlertManager 接收器失败')
      }
      this.receiverList = Array.isArray(res.data) ? res.data : []
    },
    handleSourceChange() {
      this.reload()
    },
    formatIntegrations(integrations) {
      if (!Array.isArray(integrations) || !integrations.length) return '-'
      return integrations.map(integration => integration.name).filter(Boolean).join(', ') || '-'
    }
  },
  mounted() {
    this.reload().catch(error => {
      this.loadError = error.message || '加载 AlertManager 接收器失败'
    })
  }
}
</script>

<style scoped>
.alertmanager-receiver-panel {
  margin-top: 20px;
  padding-top: 20px;
  border-top: 1px solid rgba(148, 163, 184, 0.14);
}

.subpanel-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 16px;
}

.subpanel-title {
  font-size: 16px;
  font-weight: 700;
  color: #0f172a;
}

.subpanel-subtitle {
  margin-top: 6px;
  font-size: 12px;
  color: #64748b;
}

.subpanel-alert {
  margin-top: 16px;
}

.selector-row {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 12px;
  margin: 16px 0;
}

.selector-label {
  font-size: 12px;
  color: #64748b;
}

@media (max-width: 768px) {
  .subpanel-header {
    flex-direction: column;
  }
}
</style>
