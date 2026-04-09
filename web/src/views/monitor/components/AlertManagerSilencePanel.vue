<template>
  <div class="alertmanager-silence-panel">
    <div class="subpanel-header">
      <div>
        <div class="subpanel-title">AlertManager 静默</div>
        <div class="subpanel-subtitle">基于真实 AlertManager v2 API 查看、创建和结束静默规则。</div>
      </div>
      <div class="subpanel-actions">
        <el-button size="small" icon="Refresh" @click="reload">刷新</el-button>
        <el-button
          type="primary"
          size="small"
          icon="Plus"
          :disabled="!selectedSourceId"
          @click="openCreateDialog"
        >
          新建静默
        </el-button>
      </div>
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
      title="暂无已启用的 AlertManager 告警源，可先在上方告警源管理中新增并启用。"
      type="info"
      :closable="false"
      class="subpanel-alert"
    />

    <template v-else>
      <div class="selector-row">
        <span class="selector-label">联调源</span>
        <el-select
          v-model="selectedSourceId"
          style="width: 260px"
          filterable
          @change="handleSourceChange"
        >
          <el-option
            v-for="source in alertManagerSources"
            :key="source.id"
            :label="`${source.name} (${source.apiBaseUrl || '-'})`"
            :value="source.id"
          />
        </el-select>
        <el-tag :type="statusTagType" effect="dark">{{ statusText }}</el-tag>
      </div>

      <div class="status-grid">
        <div class="status-card">
          <span class="status-label">连接地址</span>
          <span class="status-value monospace">{{ currentSource?.apiBaseUrl || '-' }}</span>
        </div>
        <div class="status-card">
          <span class="status-label">集群状态</span>
          <span class="status-value">{{ status.clusterStatus || '-' }}</span>
        </div>
        <div class="status-card">
          <span class="status-label">版本</span>
          <span class="status-value">{{ status.version || '-' }}</span>
        </div>
        <div class="status-card">
          <span class="status-label">运行时间</span>
          <span class="status-value">{{ status.uptime || '-' }}</span>
        </div>
      </div>

      <el-alert
        v-if="statusError"
        :title="statusError"
        type="warning"
        :closable="false"
        class="subpanel-alert"
      />

      <el-table
        v-loading="loading"
        :data="silenceList"
        stripe
        class="data-table"
        empty-text="当前没有静默规则"
      >
        <el-table-column label="状态" width="92">
          <template #default="{ row }">
            <el-tag :type="silenceStatusTagType(row.status)">
              {{ silenceStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="匹配条件" min-width="240" show-overflow-tooltip>
          <template #default="{ row }">
            {{ formatMatchers(row.matchers) }}
          </template>
        </el-table-column>
        <el-table-column label="时间范围" min-width="220">
          <template #default="{ row }">
            <div>{{ row.startsAt || '-' }}</div>
            <div class="muted">{{ row.endsAt || '-' }}</div>
          </template>
        </el-table-column>
        <el-table-column label="创建人" prop="createdBy" width="110" />
        <el-table-column label="说明" prop="comment" min-width="180" show-overflow-tooltip />
        <el-table-column label="操作" width="92" fixed="right">
          <template #default="{ row }">
            <el-button
              v-if="row.status !== 'expired'"
              type="danger"
              link
              @click="handleDeleteSilence(row)"
            >
              结束
            </el-button>
            <span v-else class="muted">-</span>
          </template>
        </el-table-column>
      </el-table>
    </template>

    <el-dialog
      v-model="dialogVisible"
      title="新建 AlertManager 静默"
      width="720px"
      destroy-on-close
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="formRules"
        label-width="100px"
      >
        <el-form-item label="告警源" prop="sourceId">
          <el-select v-model="form.sourceId" style="width: 100%">
            <el-option
              v-for="source in alertManagerSources"
              :key="source.id"
              :label="source.name"
              :value="source.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="时间范围" prop="timeRange">
          <el-date-picker
            v-model="form.timeRange"
            type="datetimerange"
            style="width: 100%"
            range-separator="至"
            start-placeholder="开始时间"
            end-placeholder="结束时间"
            value-format="YYYY-MM-DD HH:mm:ss"
          />
        </el-form-item>
        <el-form-item label="创建人">
          <el-input v-model="form.createdBy" placeholder="默认使用当前登录用户" />
        </el-form-item>
        <el-form-item label="静默说明" prop="comment">
          <el-input
            v-model="form.comment"
            type="textarea"
            :rows="3"
            placeholder="例如：opsnexus-e2e 调试期间屏蔽 Watchdog"
          />
        </el-form-item>
        <el-form-item label="匹配条件">
          <div class="matcher-list">
            <div
              v-for="(matcher, index) in form.matchers"
              :key="index"
              class="matcher-row"
            >
              <el-input v-model="matcher.name" placeholder="label 名称，如 alertname" />
              <el-input v-model="matcher.value" placeholder="匹配值，如 Watchdog" />
              <el-switch v-model="matcher.isRegex" inline-prompt active-text="正则" inactive-text="精确" />
              <el-button
                :disabled="form.matchers.length === 1"
                icon="Delete"
                @click="removeMatcher(index)"
              />
            </div>
            <el-button type="primary" link icon="Plus" @click="addMatcher">添加匹配条件</el-button>
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="submitForm">创建静默</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import { ElMessage, ElMessageBox } from 'element-plus'

const createMatcher = () => ({
  name: 'alertname',
  value: '',
  isRegex: false
})

const createEmptyForm = (sourceId = null) => ({
  sourceId,
  timeRange: [],
  createdBy: '',
  comment: '',
  matchers: [createMatcher()]
})

export default {
  name: 'AlertManagerSilencePanel',
  data() {
    return {
      loading: false,
      loadError: '',
      statusError: '',
      dialogVisible: false,
      submitting: false,
      alertManagerSources: [],
      selectedSourceId: null,
      status: {
        available: false,
        clusterStatus: '',
        version: '',
        uptime: ''
      },
      silenceList: [],
      form: createEmptyForm(),
      formRules: {
        sourceId: [{ required: true, message: '请选择告警源', trigger: 'change' }],
        comment: [{ required: true, message: '请输入静默说明', trigger: 'blur' }],
        timeRange: [{
          validator: (rule, value, callback) => {
            if (!Array.isArray(value) || value.length !== 2) {
              callback(new Error('请选择静默时间范围'))
              return
            }
            callback()
          },
          trigger: 'change'
        }]
      }
    }
  },
  computed: {
    currentSource() {
      return this.alertManagerSources.find(source => source.id === this.selectedSourceId) || null
    },
    statusText() {
      if (this.statusError) return '状态异常'
      return this.status.available ? '连接正常' : '待检测'
    },
    statusTagType() {
      if (this.statusError) return 'warning'
      return this.status.available ? 'success' : 'info'
    }
  },
  methods: {
    async reload() {
      this.loadError = ''
      this.statusError = ''
      await this.loadSources()
      if (!this.selectedSourceId) return

      this.loading = true
      const [statusResult, silenceResult] = await Promise.allSettled([
        this.fetchStatus(),
        this.fetchSilences()
      ])
      if (statusResult.status === 'rejected') {
        this.statusError = statusResult.reason?.message || '获取 AlertManager 状态失败'
      }
      if (silenceResult.status === 'rejected') {
        this.loadError = silenceResult.reason?.message || '获取 AlertManager 静默失败'
      }
      this.loading = false
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
        this.status = { available: false, clusterStatus: '', version: '', uptime: '' }
        this.silenceList = []
        return
      }
      const exists = this.alertManagerSources.some(source => source.id === this.selectedSourceId)
      if (!exists) {
        this.selectedSourceId = this.alertManagerSources[0].id
      }
      if (!this.form.sourceId) {
        this.form.sourceId = this.selectedSourceId
      }
    },
    async fetchStatus() {
      const { data: res } = await this.$api.getMonitorAlertManagerStatus({ sourceId: this.selectedSourceId })
      if (res.code !== 200) {
        throw new Error(res.message || '获取 AlertManager 状态失败')
      }
      this.status = res.data || this.status
    },
    async fetchSilences() {
      const { data: res } = await this.$api.getMonitorAlertManagerSilences({ sourceId: this.selectedSourceId })
      if (res.code !== 200) {
        throw new Error(res.message || '获取 AlertManager 静默失败')
      }
      this.silenceList = Array.isArray(res.data) ? res.data : []
    },
    handleSourceChange() {
      this.reload().catch(error => {
        this.loadError = error.message || '切换 AlertManager 告警源失败'
      })
    },
    openCreateDialog() {
      this.form = createEmptyForm(this.selectedSourceId)
      this.dialogVisible = true
    },
    addMatcher() {
      this.form.matchers.push(createMatcher())
    },
    removeMatcher(index) {
      if (this.form.matchers.length === 1) return
      this.form.matchers.splice(index, 1)
    },
    buildSubmitPayload() {
      const [startAt, endAt] = this.form.timeRange
      return {
        sourceId: this.form.sourceId,
        createdBy: this.form.createdBy,
        comment: this.form.comment,
        startsAt: this.toISOString(startAt),
        endsAt: this.toISOString(endAt),
        matchers: this.form.matchers.map(matcher => ({
          name: matcher.name,
          value: matcher.value,
          isRegex: !!matcher.isRegex
        }))
      }
    },
    toISOString(value) {
      if (!value) return ''
      return new Date(String(value).replace(' ', 'T')).toISOString()
    },
    validateMatchers() {
      const invalidMatcher = this.form.matchers.some(matcher => !matcher.name || !matcher.value)
      if (invalidMatcher) {
        throw new Error('请完整填写所有匹配条件')
      }
    },
    async submitForm() {
      const valid = await this.$refs.formRef.validate().catch(() => false)
      if (!valid) return

      try {
        this.validateMatchers()
      } catch (error) {
        ElMessage.error(error.message)
        return
      }

      this.submitting = true
      try {
        const { data: res } = await this.$api.createMonitorAlertManagerSilence(this.buildSubmitPayload())
        if (res.code !== 200) {
          throw new Error(res.message || '创建静默失败')
        }
        ElMessage.success('AlertManager 静默创建成功')
        this.dialogVisible = false
        this.selectedSourceId = this.form.sourceId
        await this.reload()
      } catch (error) {
        ElMessage.error(error?.response?.data?.message || error.message || '创建静默失败')
      } finally {
        this.submitting = false
      }
    },
    async handleDeleteSilence(row) {
      try {
        await ElMessageBox.confirm(`确认结束静默“${this.formatMatchers(row.matchers)}”吗？`, '提示', {
          type: 'warning'
        })
        const { data: res } = await this.$api.deleteMonitorAlertManagerSilence(row.id, {
          sourceId: row.sourceId || this.selectedSourceId
        })
        if (res.code !== 200) {
          throw new Error(res.message || '结束静默失败')
        }
        ElMessage.success('静默已结束')
        await this.reload()
      } catch (error) {
        if (error !== 'cancel' && error !== 'close') {
          ElMessage.error(error?.response?.data?.message || error.message || '结束静默失败')
        }
      }
    },
    formatMatchers(matchers) {
      if (!Array.isArray(matchers) || !matchers.length) return '-'
      return matchers
        .map(matcher => `${matcher.name}${matcher.isRegex ? '=~' : '='}${matcher.value}`)
        .join(', ')
    },
    silenceStatusText(status) {
      const textMap = {
        active: '生效中',
        pending: '待生效',
        expired: '已过期'
      }
      return textMap[status] || status || '-'
    },
    silenceStatusTagType(status) {
      const typeMap = {
        active: 'success',
        pending: 'warning',
        expired: 'info'
      }
      return typeMap[status] || 'info'
    }
  },
  mounted() {
    this.reload().catch(error => {
      this.loadError = error.message || '加载 AlertManager 静默失败'
    })
  }
}
</script>

<style scoped>
.alertmanager-silence-panel {
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
  color: #111827;
}

.subpanel-subtitle {
  margin-top: 6px;
  font-size: 12px;
  color: #6b7280;
}

.subpanel-actions {
  display: flex;
  gap: 10px;
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

.selector-label,
.status-label {
  font-size: 12px;
  color: #64748b;
}

.status-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
  margin-bottom: 16px;
}

.status-card {
  padding: 12px 14px;
  border-radius: 14px;
  background: linear-gradient(135deg, rgba(59, 130, 246, 0.08), rgba(248, 113, 113, 0.08));
  border: 1px solid rgba(148, 163, 184, 0.12);
}

.status-value {
  display: block;
  margin-top: 6px;
  color: #111827;
  font-weight: 600;
}

.monospace {
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 12px;
  word-break: break-all;
}

.matcher-list {
  width: 100%;
}

.matcher-row {
  display: grid;
  grid-template-columns: minmax(120px, 1fr) minmax(160px, 1.2fr) 96px 48px;
  gap: 10px;
  margin-bottom: 10px;
}

.muted {
  color: #94a3b8;
}

@media (max-width: 768px) {
  .subpanel-header {
    flex-direction: column;
  }

  .status-grid {
    grid-template-columns: 1fr;
  }

  .matcher-row {
    grid-template-columns: 1fr;
  }
}
</style>
