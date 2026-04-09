<template>
  <PageContainer wide class="alert-center-modern-page">
    <PageHeader
      eyebrow="Monitoring"
      title="告警中心"
      subtitle="聚合告警事件、推送链路与告警源配置，统一查看风险态势"
    >
      <template #actions>
        <el-button @click="refreshAll">
          <el-icon><Refresh /></el-icon>
          刷新数据
        </el-button>
        <el-button type="primary" @click="openCreateSourceDialog">
          <el-icon><Plus /></el-icon>
          新增告警源
        </el-button>
      </template>
      <template #meta>
        <div class="page-chip-row">
          <span class="platform-chip">最近告警 {{ summary.latestAlertTime || '-' }}</span>
          <span class="platform-chip">最近推送 {{ summary.latestNotifyTime || '-' }}</span>
          <span class="platform-chip">告警源 {{ summary.enabledAlertSources }}/{{ summary.totalAlertSources }}</span>
        </div>
      </template>
      <template #intro>
        <PageIntro
          title="工作说明"
          text="告警中心将事件列表、告警源配置与静默管理收敛在一个界面中，便于值班、研判和通知链路排查。"
        />
      </template>
    </PageHeader>

    <el-alert
      v-if="summary.totalAlertSources === 0"
      title="当前未接入告警源，统计数据与事件汇聚可能不完整"
      type="warning"
      :closable="false"
      class="page-alert"
    />

    <StatStrip :items="statItems" />

    <div class="alert-center-grid">
      <SectionCard
        title="告警事件列表"
        subtitle="聚合待处理、处理中与已恢复事件，支持按命名空间和工作负载筛选"
        class="alert-center-grid__main"
      >
        <el-alert
          v-if="routeContext.source === 'capacity-suggestion'"
          :title="`已按容量建议定位到 ${routeContext.namespace || '-'} / ${routeContext.workloadName || '-'}`"
          type="info"
          :closable="false"
          class="page-alert"
        />

        <PageToolbar class="alert-center-section-toolbar">
          <el-form :model="queryParams" :inline="true" class="filter-cluster">
            <el-form-item class="filter-field" label="关键词">
              <el-input
                v-model="queryParams.keyword"
                placeholder="搜索告警标题、描述或业务线"
                clearable
                @keyup.enter="handleQuery"
              />
            </el-form-item>
            <el-form-item class="filter-field" label="状态">
              <el-select v-model="queryParams.status" placeholder="全部状态" clearable>
                <el-option label="待处理" :value="1" />
                <el-option label="处理中" :value="2" />
                <el-option label="已恢复" :value="3" />
              </el-select>
            </el-form-item>
            <el-form-item class="filter-field" label="等级">
              <el-select v-model="queryParams.level" placeholder="全部等级" clearable>
                <el-option label="P1" value="P1" />
                <el-option label="P2" value="P2" />
                <el-option label="P3" value="P3" />
                <el-option label="P4" value="P4" />
              </el-select>
            </el-form-item>
            <el-form-item class="filter-field" label="命名空间">
              <el-input
                v-model="queryParams.namespace"
                placeholder="例如 opsnexus-apply-e2e"
                clearable
                @keyup.enter="handleQuery"
              />
            </el-form-item>
            <el-form-item class="filter-field" label="工作负载">
              <el-input
                v-model="queryParams.workloadName"
                placeholder="例如 demo-nginx"
                clearable
                @keyup.enter="handleQuery"
              />
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

        <el-table
          v-loading="loading"
          :data="incidentList"
          stripe
          class="alert-center-table"
          empty-text="暂无告警事件"
        >
          <el-table-column label="告警时间" prop="alertTime" min-width="170" />
          <el-table-column label="等级" width="90">
            <template #default="{ row }">
              <el-tag :type="levelTagType(row.alertLevel)" effect="dark">
                {{ row.alertLevel || '-' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="状态" width="110">
            <template #default="{ row }">
              <el-tag :type="incidentStatusTagType(row.status)" effect="dark">
                {{ incidentStatusText(row.status) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="业务线" prop="businessLine" min-width="130" show-overflow-tooltip />
          <el-table-column label="告警描述" prop="alertDesc" min-width="320" show-overflow-tooltip />
          <el-table-column label="处理人" prop="handler" width="120" />
          <el-table-column label="操作" width="120" fixed="right">
            <template #default="{ row }">
              <el-button v-if="row.detailUrl" type="primary" link @click="openDetailUrl(row.detailUrl)">
                查看详情
              </el-button>
              <span v-else class="muted">-</span>
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
      </SectionCard>

      <SectionCard
        title="告警源配置"
        subtitle="支持 AlertManager、Nightingale、Zabbix 等告警平台统一接入与维护"
        class="alert-center-grid__side"
      >
        <template #actions>
          <el-button type="primary" plain @click="openCreateSourceDialog">
            <el-icon><Plus /></el-icon>
            新增告警源
          </el-button>
        </template>

        <el-table
          v-if="sourceList.length"
          :data="sourceList"
          stripe
          class="alert-center-table alert-center-table--compact"
          empty-text="暂无告警源"
        >
          <el-table-column label="名称" prop="name" min-width="150" show-overflow-tooltip />
          <el-table-column label="类型" width="130">
            <template #default="{ row }">
              <el-tag :type="sourceTypeTagType(row.type)">
                {{ sourceTypeText(row.type) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="状态" width="90">
            <template #default="{ row }">
              <el-tag :type="Number(row.status) === 1 ? 'success' : 'info'">
                {{ Number(row.status) === 1 ? '启用' : '停用' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="接入地址" min-width="220" show-overflow-tooltip>
            <template #default="{ row }">
              {{ row.apiBaseUrl || '-' }}
            </template>
          </el-table-column>
          <el-table-column label="主机 ID" prop="hostId" width="90" />
          <el-table-column label="Key ID" prop="keyId" width="90" />
          <el-table-column label="操作" width="210" fixed="right">
            <template #default="{ row }">
              <div class="operation-group">
                <el-button type="primary" link @click="openEditSourceDialog(row)">编辑</el-button>
                <el-button :type="Number(row.status) === 1 ? 'warning' : 'primary'" link @click="toggleSourceStatus(row)">
                  {{ Number(row.status) === 1 ? '停用' : '启用' }}
                </el-button>
                <el-button type="danger" link @click="handleDeleteSource(row)">删除</el-button>
              </div>
            </template>
          </el-table-column>
        </el-table>
        <div v-else class="empty-state alert-center-empty-state">
          <div class="empty-state__halo"></div>
          <div class="empty-state__title">暂无告警源</div>
          <div class="empty-state__description">
            请先接入至少一个 AlertManager 告警源，用于汇聚事件与推送数据。
          </div>
          <div class="empty-state__actions">
            <el-button type="primary" @click="openCreateSourceDialog">立即新增告警源</el-button>
          </div>
        </div>

        <div class="alert-center-silence">
          <AlertManagerSilencePanel ref="alertManagerSilencePanel" />
        </div>
      </SectionCard>
    </div>

    <el-dialog
      v-model="sourceDialogVisible"
      :title="sourceDialogMode === 'create' ? '新增告警源' : '编辑告警源'"
      width="680px"
      destroy-on-close
    >
      <el-form
        ref="sourceFormRef"
        :model="sourceForm"
        :rules="sourceRules"
        label-width="110px"
        class="source-form"
      >
        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="名称" prop="name">
              <el-input v-model="sourceForm.name" placeholder="例如 opsnexus-alertmanager" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="类型" prop="type">
              <el-select v-model="sourceForm.type" placeholder="请选择告警源类型" style="width: 100%">
                <el-option label="FlashDuty" :value="1" />
                <el-option label="Zabbix" :value="2" />
                <el-option label="Nightingale" :value="3" />
                <el-option label="AlertManager" :value="4" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="状态" prop="status">
              <el-radio-group v-model="sourceForm.status">
                <el-radio :label="1">启用</el-radio>
                <el-radio :label="0">停用</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="认证信息">
              <el-input v-model="sourceForm.appKey" placeholder="选填，可填写 token / appKey" />
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item label="API 地址" prop="apiBaseUrl">
              <el-input v-model="sourceForm.apiBaseUrl" placeholder="例如 http://alertmanager:9093" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="主机 ID">
              <el-input-number v-model="sourceForm.hostId" :min="0" controls-position="right" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="Key ID">
              <el-input-number v-model="sourceForm.keyId" :min="0" controls-position="right" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item label="备注">
              <el-input v-model="sourceForm.remark" type="textarea" :rows="3" placeholder="补充说明告警源用途或备注" />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <template #footer>
        <el-button @click="sourceDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="sourceSubmitting" @click="submitSourceForm">保存</el-button>
      </template>
    </el-dialog>
  </PageContainer>
</template>

<script>
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Refresh, Search } from '@element-plus/icons-vue'
import PageContainer from '@/components/platform/PageContainer.vue'
import PageHeader from '@/components/platform/PageHeader.vue'
import PageIntro from '@/components/platform/PageIntro.vue'
import PageToolbar from '@/components/platform/PageToolbar.vue'
import SectionCard from '@/components/platform/SectionCard.vue'
import StatStrip from '@/components/platform/StatStrip.vue'
import AlertManagerSilencePanel from './components/AlertManagerSilencePanel.vue'

const createEmptySourceForm = () => ({
  id: null,
  name: '',
  type: 4,
  status: 1,
  apiBaseUrl: '',
  appKey: '',
  keyId: 0,
  hostId: 0,
  remark: ''
})

const createEmptyIncidentQuery = () => ({
  keyword: '',
  status: '',
  level: '',
  namespace: '',
  workloadName: '',
  pageNum: 1,
  pageSize: 10
})

const createEmptyRouteContext = () => ({
  source: '',
  namespace: '',
  workloadName: ''
})

const parseRouteStatus = value => {
  if (value === '' || value === null || value === undefined) return ''
  const parsed = Number(value)
  return Number.isNaN(parsed) ? '' : parsed
}

export default {
  name: 'AlertCenter',
  components: {
    AlertManagerSilencePanel,
    PageContainer,
    PageHeader,
    PageIntro,
    PageToolbar,
    SectionCard,
    StatStrip,
    Plus,
    Refresh,
    Search
  },
  data() {
    return {
      loading: false,
      total: 0,
      incidentList: [],
      sourceList: [],
      sourceDialogVisible: false,
      sourceDialogMode: 'create',
      sourceSubmitting: false,
      sourceForm: createEmptySourceForm(),
      sourceRules: {
        name: [{ required: true, message: '请输入告警源名称', trigger: 'blur' }],
        type: [{ required: true, message: '请选择告警源类型', trigger: 'change' }],
        status: [{ required: true, message: '请选择状态', trigger: 'change' }]
      },
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
      queryParams: createEmptyIncidentQuery(),
      routeContext: createEmptyRouteContext()
    }
  },
  computed: {
    statItems() {
      return [
        { label: '待处理告警', value: this.summary.openIncidents, hint: '需要优先处理', tone: 'danger' },
        { label: '处理中告警', value: this.summary.processingIncidents, hint: '正在跟进', tone: 'warning' },
        { label: '已恢复告警', value: this.summary.resolvedIncidents, hint: '已闭环恢复', tone: 'success' },
        { label: '高危推送', value: this.summary.criticalWebhookLogs, hint: 'critical / P1 / P2 webhook', tone: 'primary' }
      ]
    }
  },
  methods: {
    applyRouteQuery() {
      const routeQuery = this.$route?.query || {}
      this.routeContext = {
        source: String(routeQuery.source || '').trim(),
        namespace: String(routeQuery.namespace || '').trim(),
        workloadName: String(routeQuery.workloadName || '').trim()
      }

      this.queryParams = {
        ...this.queryParams,
        keyword: String(routeQuery.keyword || '').trim(),
        status: parseRouteStatus(routeQuery.status),
        level: String(routeQuery.level || '').trim(),
        namespace: this.routeContext.namespace,
        workloadName: this.routeContext.workloadName,
        pageNum: 1
      }
    },
    buildQueryParams() {
      const params = {
        pageNum: this.queryParams.pageNum,
        pageSize: this.queryParams.pageSize
      }
      if (this.queryParams.keyword) params.keyword = this.queryParams.keyword
      if (this.queryParams.status !== '' && this.queryParams.status !== null && this.queryParams.status !== undefined) {
        params.status = this.queryParams.status
      }
      if (this.queryParams.level) params.level = this.queryParams.level
      if (this.queryParams.namespace) params.namespace = this.queryParams.namespace
      if (this.queryParams.workloadName) params.workloadName = this.queryParams.workloadName
      return params
    },
    async fetchSummary() {
      const { data: res } = await this.$api.getMonitorAlertSummary()
      if (res.code !== 200) {
        throw new Error(res.message || '获取告警概览失败')
      }
      this.summary = res.data || this.summary
    },
    async fetchIncidents() {
      const { data: res } = await this.$api.queryMonitorAlertIncidentList(this.buildQueryParams())
      if (res.code !== 200) {
        throw new Error(res.message || '获取告警事件失败')
      }
      const data = res.data || {}
      this.incidentList = data.list || []
      this.total = data.total || 0
    },
    async fetchSources() {
      const { data: res } = await this.$api.getMonitorAlertSourceList()
      if (res.code !== 200) {
        throw new Error(res.message || '获取告警源失败')
      }
      this.sourceList = res.data || []
    },
    async refreshAll() {
      this.loading = true
      try {
        await Promise.all([this.fetchSummary(), this.fetchIncidents(), this.fetchSources()])
        this.$refs.alertManagerSilencePanel?.reload?.()
      } catch (error) {
        ElMessage.error(error.message || '刷新告警中心失败')
      } finally {
        this.loading = false
      }
    },
    handleQuery() {
      this.queryParams.pageNum = 1
      this.fetchIncidents().catch(error => {
        ElMessage.error(error.message || '获取告警事件失败')
      })
    },
    resetQuery() {
      this.queryParams = createEmptyIncidentQuery()
      this.routeContext = createEmptyRouteContext()
      this.$router.replace({ path: this.$route.path, query: {} }).catch(() => {})
      this.fetchIncidents().catch(error => {
        ElMessage.error(error.message || '获取告警事件失败')
      })
    },
    handleSizeChange(size) {
      this.queryParams.pageSize = size
      this.queryParams.pageNum = 1
      this.fetchIncidents().catch(error => {
        ElMessage.error(error.message || '获取告警事件失败')
      })
    },
    handleCurrentChange(page) {
      this.queryParams.pageNum = page
      this.fetchIncidents().catch(error => {
        ElMessage.error(error.message || '获取告警事件失败')
      })
    },
    levelTagType(level) {
      const typeMap = {
        P1: 'danger',
        P2: 'warning',
        P3: 'primary',
        P4: 'info'
      }
      return typeMap[level] || 'info'
    },
    incidentStatusText(status) {
      const textMap = {
        1: '待处理',
        2: '处理中',
        3: '已恢复'
      }
      return textMap[Number(status)] || '未知状态'
    },
    incidentStatusTagType(status) {
      const typeMap = {
        1: 'danger',
        2: 'warning',
        3: 'success'
      }
      return typeMap[Number(status)] || 'info'
    },
    sourceTypeText(type) {
      const textMap = {
        1: 'FlashDuty',
        2: 'Zabbix',
        3: 'Nightingale',
        4: 'AlertManager'
      }
      return textMap[Number(type)] || '未知类型'
    },
    sourceTypeTagType(type) {
      const tagMap = {
        1: 'danger',
        2: 'warning',
        3: 'success',
        4: 'primary'
      }
      return tagMap[Number(type)] || 'info'
    },
    openDetailUrl(url) {
      window.open(url, '_blank')
    },
    openCreateSourceDialog() {
      this.sourceDialogMode = 'create'
      this.sourceForm = createEmptySourceForm()
      this.sourceDialogVisible = true
    },
    openEditSourceDialog(row) {
      this.sourceDialogMode = 'edit'
      this.sourceForm = {
        id: row.id,
        name: row.name || '',
        type: Number(row.type || 4),
        status: Number(row.status ?? 1),
        apiBaseUrl: row.apiBaseUrl || '',
        appKey: row.appKey || '',
        keyId: Number(row.keyId || 0),
        hostId: Number(row.hostId || 0),
        remark: row.remark || ''
      }
      this.sourceDialogVisible = true
    },
    async submitSourceForm() {
      const valid = await this.$refs.sourceFormRef.validate().catch(() => false)
      if (!valid) return

      this.sourceSubmitting = true
      try {
        const payload = {
          name: this.sourceForm.name,
          type: this.sourceForm.type,
          status: this.sourceForm.status,
          apiBaseUrl: this.sourceForm.apiBaseUrl,
          appKey: this.sourceForm.appKey,
          keyId: this.sourceForm.keyId,
          hostId: this.sourceForm.hostId,
          remark: this.sourceForm.remark
        }
        if (this.sourceDialogMode === 'create') {
          await this.$api.createMonitorAlertSource(payload)
          ElMessage.success('告警源创建成功')
        } else {
          await this.$api.updateMonitorAlertSource(this.sourceForm.id, payload)
          ElMessage.success('告警源更新成功')
        }
        this.sourceDialogVisible = false
        await Promise.all([this.fetchSummary(), this.fetchSources()])
        this.$refs.alertManagerSilencePanel?.reload?.()
      } catch (error) {
        ElMessage.error(error?.response?.data?.message || error.message || '提交告警源失败')
      } finally {
        this.sourceSubmitting = false
      }
    },
    async toggleSourceStatus(row) {
      const targetStatus = Number(row.status) === 1 ? 0 : 1
      const actionText = targetStatus === 1 ? '启用' : '停用'
      try {
        await ElMessageBox.confirm(`确认${actionText}告警源 ${row.name} 吗`, '提示', {
          type: 'warning'
        })
        await this.$api.updateMonitorAlertSourceStatus(row.id, targetStatus)
        ElMessage.success(`告警源已${actionText}`)
        await Promise.all([this.fetchSummary(), this.fetchSources()])
        this.$refs.alertManagerSilencePanel?.reload?.()
      } catch (error) {
        if (error !== 'cancel' && error !== 'close') {
          ElMessage.error(error?.response?.data?.message || error.message || `${actionText}告警源失败`)
        }
      }
    },
    async handleDeleteSource(row) {
      try {
        await ElMessageBox.confirm(`确认删除告警源 ${row.name} 吗`, '提示', {
          type: 'warning'
        })
        await this.$api.deleteMonitorAlertSource(row.id)
        ElMessage.success('告警源删除成功')
        await Promise.all([this.fetchSummary(), this.fetchSources()])
        this.$refs.alertManagerSilencePanel?.reload?.()
      } catch (error) {
        if (error !== 'cancel' && error !== 'close') {
          ElMessage.error(error?.response?.data?.message || error.message || '删除告警源失败')
        }
      }
    }
  },
  watch: {
    '$route.query': {
      handler() {
        this.applyRouteQuery()
        this.fetchIncidents().catch(error => {
          ElMessage.error(error.message || '获取告警事件失败')
        })
      }
    }
  },
  created() {
    this.applyRouteQuery()
    this.refreshAll()
  }
}
</script>

<style scoped>
.alert-center-grid {
  display: grid;
  grid-template-columns: minmax(0, 1.6fr) minmax(360px, 1fr);
  gap: 16px;
  align-items: start;
}

.alert-center-grid__main,
.alert-center-grid__side {
  min-width: 0;
}

.alert-center-section-toolbar {
  margin-bottom: 18px;
  padding: 0;
  background: transparent;
  border: none;
  box-shadow: none;
}

.alert-center-empty-state {
  margin-bottom: 18px;
}

.alert-center-silence {
  margin-top: 18px;
}

.alert-center-modern-page :deep(.el-table .cell.el-tooltip) {
  white-space: normal;
  word-break: break-word;
  line-height: 1.5;
}

.source-form :deep(.el-input-number) {
  width: 100%;
}

@media (max-width: 1280px) {
  .alert-center-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .alert-center-modern-page :deep(.page-actions) {
    width: 100%;
    justify-content: flex-start;
  }
}
</style>
