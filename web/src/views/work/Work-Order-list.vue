<template>
  <TablePage
    class="workorder-table-page"
    section-title="工单流转列表"
    section-subtitle="统一承载快速发布、脚本发布、服务上线与 SQL 工单，并提供 AI 与知识联动入口。"
    wide
  >
    <template #header>
      <PageHeader
        eyebrow="Change Center"
        title="工单中心"
        subtitle="把发布、变更、审批与执行联动到同一工作台，帮助团队优先处理真正有风险的事项。"
      >
        <template #actions>
          <el-button icon="Refresh" @click="refreshAll">刷新</el-button>
        </template>
        <template #intro>
          <PageIntro
            title="处理优先级"
            text="先看待处理和高风险工单，再结合 AI 诊断与知识库确认上下文，最后执行批准、驳回或执行动作。"
          />
        </template>
      </PageHeader>
    </template>

    <template #toolbar>
      <PageToolbar>
        <el-form :inline="true" class="filter-cluster">
          <el-form-item class="filter-field" label="关键字">
            <el-input v-model="queryParams.keyword" placeholder="搜索标题、申请人、应用" clearable @keyup.enter="handleQuery" />
          </el-form-item>
          <el-form-item class="filter-field" label="类型">
            <el-select v-model="queryParams.type" placeholder="全部" clearable>
              <el-option label="快速发布" value="quick" />
              <el-option label="脚本发布" value="script" />
              <el-option label="服务上线" value="service" />
              <el-option label="SQL 工单" value="sql" />
            </el-select>
          </el-form-item>
          <el-form-item class="filter-field" label="状态">
            <el-select v-model="queryParams.status" placeholder="全部" clearable>
              <el-option label="待处理/待发布" :value="1" />
              <el-option label="成功" :value="2" />
              <el-option label="失败" :value="3" />
              <el-option label="驳回/取消" :value="4" />
              <el-option label="已取消" :value="5" />
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

    <el-table v-loading="loading" :data="workOrderList" stripe class="data-table" empty-text="暂无工单数据">
      <el-table-column label="类型" width="110">
        <template #default="{ row }">
          <el-tag :type="typeTagType(row.type)">{{ row.typeLabel }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="标题" prop="title" min-width="220">
        <template #default="{ row }">
          <el-tooltip :content="row.title || '-'" placement="top" :show-after="300">
            <span>{{ formatCompactLabel(row.title, 16, 6) }}</span>
          </el-tooltip>
        </template>
      </el-table-column>
      <el-table-column label="应用" prop="appName" min-width="140" show-overflow-tooltip />
      <el-table-column label="申请人" prop="applicantName" width="120" />
      <el-table-column label="当前处理人" prop="currentHandler" min-width="110" show-overflow-tooltip />
      <el-table-column label="风险" width="100">
        <template #default="{ row }">
          <el-tag v-if="getOrderRiskText(row)" :type="riskTagType(getOrderRiskLevel(row))">{{ getOrderRiskText(row) }}</el-tag>
          <span v-else>-</span>
        </template>
      </el-table-column>
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="statusTagType(getOrderStatusText(row))">{{ getOrderStatusText(row) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="180" fixed="right">
        <template #default="{ row }">
          <div class="row-actions">
            <el-button type="primary" link @click="showDetail(row)">详情</el-button>
            <el-button v-if="row.canApprove" type="success" link @click="approveOrder(row)">批准</el-button>
            <el-button v-if="row.canReject" type="danger" link @click="rejectOrder(row)">驳回</el-button>
            <el-button v-if="row.canExecute" type="warning" link @click="executeOrder(row)">执行</el-button>
          </div>
        </template>
      </el-table-column>
    </el-table>

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
      :title="detailData.title || '工单详情'"
      width="880px"
      class="ops-dialog ops-overlay--lg ops-overlay--full workorder-detail-dialog"
    >
      <div class="ops-dialog__body workorder-detail-shell">
        <section v-if="detailData.basic" class="ops-dialog__section">
          <div class="ops-dialog__section-title">工单概览</div>
          <el-descriptions :column="2" border>
            <el-descriptions-item label="工单类型">{{ detailData.typeLabel }}</el-descriptions-item>
            <el-descriptions-item label="状态">{{ getOrderStatusText(detailData) }}</el-descriptions-item>
            <el-descriptions-item v-for="(value, key) in detailData.basic" :key="key" :label="key">
              {{ stringifyValue(value) }}
            </el-descriptions-item>
          </el-descriptions>
        </section>

        <section v-if="detailData.items && detailData.items.length" class="ops-dialog__section detail-items">
          <div class="ops-dialog__section-title">关联条目</div>
          <el-table :data="detailData.items" size="small" stripe max-height="320">
            <el-table-column
              v-for="column in detailColumns"
              :key="column"
              :prop="column"
              :label="column"
              min-width="120"
              show-overflow-tooltip
            />
          </el-table>
        </section>
      </div>
      <template #footer>
        <div v-if="detailData.type === 'script'" class="detail-action-row">
          <el-button v-if="detailData.canApprove" type="success" @click="approveOrder(detailData)">批准工单</el-button>
          <el-button v-if="detailData.canReject" type="danger" @click="rejectOrder(detailData)">驳回工单</el-button>
          <el-button v-if="detailData.canExecute" type="warning" @click="executeOrder(detailData)">执行工单</el-button>
        </div>
        <div v-else-if="detailData.type === 'sql'" class="detail-action-row">
          <el-button v-if="getOrderLink(detailData.basic, 'aiDiagnosisPath')" type="primary" @click="openInternalPath(getOrderLink(detailData.basic, 'aiDiagnosisPath'))">AI 诊断</el-button>
          <el-button v-if="getOrderLink(detailData.basic, 'knowledgePath')" type="warning" @click="openInternalPath(getOrderLink(detailData.basic, 'knowledgePath'))">关联知识</el-button>
          <el-button v-if="detailData.canApprove" type="success" @click="approveOrder(detailData)">批准工单</el-button>
          <el-button v-if="detailData.canReject" type="danger" @click="rejectOrder(detailData)">驳回工单</el-button>
          <el-button v-if="detailData.canExecute" type="warning" @click="executeOrder(detailData)">执行工单</el-button>
        </div>
      </template>
    </el-dialog>
  </TablePage>
</template>

<script>
import { ElMessage, ElMessageBox } from 'element-plus'
import PageHeader from '@/components/platform/PageHeader.vue'
import PageIntro from '@/components/platform/PageIntro.vue'
import PageToolbar from '@/components/platform/PageToolbar.vue'
import StatStrip from '@/components/platform/StatStrip.vue'
import TablePage from '@/components/platform/TablePage.vue'

export default {
  name: 'WorkOrderList',
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
      total: 0,
      workOrderList: [],
      detailVisible: false,
      detailData: {},
      queryParams: {
        page: 1,
        pageSize: 10,
        type: '',
        status: '',
        keyword: ''
      },
      summary: {
        total: 0,
        pending: 0,
        running: 0,
        success: 0,
        failed: 0,
        canceled: 0,
        quickDeploy: 0,
        scriptRelease: 0,
        serviceRelease: 0
      }
    }
  },
  computed: {
    detailColumns() {
      const first = this.detailData.items?.[0] || {}
      return Object.keys(first)
    },
    statItems() {
      return [
        { label: '工单总数', value: this.summary.total, hint: '当前条件下的工单规模', tone: 'primary' },
        { label: '待处理', value: this.summary.pending, hint: '需要审批或执行', tone: 'warning' },
        { label: '处理中', value: this.summary.running, hint: '仍在流转中的工单', tone: 'neutral' },
        { label: '成功', value: this.summary.success, hint: '闭环完成的工单', tone: 'success' },
        { label: '失败', value: this.summary.failed, hint: '需要复盘的工单', tone: 'danger' }
      ]
    }
  },
  methods: {
    formatCompactLabel(value, prefixLength = 16, suffixLength = 6) {
      const text = String(value || '')
      if (!text) return '-'
      if (text.length <= prefixLength + suffixLength) {
        return text
      }
      return `${text.slice(0, prefixLength)}...${text.slice(-suffixLength)}`
    },
    buildQueryParams() {
      const params = {
        page: this.queryParams.page,
        pageSize: this.queryParams.pageSize
      }
      if (this.queryParams.type) params.type = this.queryParams.type
      if (this.queryParams.status) params.status = this.queryParams.status
      if (this.queryParams.keyword) params.keyword = this.queryParams.keyword
      return params
    },
    async fetchSummary() {
      const { data: res } = await this.$api.getWorkOrderSummary()
      if (res.code !== 200) {
        throw new Error(res.message || '获取工单摘要失败')
      }
      this.summary = res.data || this.summary
    },
    async fetchList() {
      const { data: res } = await this.$api.getWorkOrderList(this.buildQueryParams())
      if (res.code !== 200) {
        throw new Error(res.message || '获取工单列表失败')
      }
      const payload = res.data || {}
      this.workOrderList = payload.list || []
      this.total = payload.total || 0
    },
    async refreshAll() {
      this.loading = true
      try {
        await Promise.all([this.fetchSummary(), this.fetchList()])
      } catch (error) {
        ElMessage.error(error.message || '刷新工单中心失败')
      } finally {
        this.loading = false
      }
    },
    handleQuery() {
      this.queryParams.page = 1
      this.fetchList().catch(error => ElMessage.error(error.message || '搜索工单失败'))
    },
    resetQuery() {
      this.queryParams = {
        page: 1,
        pageSize: 10,
        type: '',
        status: '',
        keyword: ''
      }
      this.refreshAll()
    },
    handleSizeChange(size) {
      this.queryParams.pageSize = size
      this.queryParams.page = 1
      this.fetchList().catch(error => ElMessage.error(error.message || '刷新工单失败'))
    },
    handleCurrentChange(page) {
      this.queryParams.page = page
      this.fetchList().catch(error => ElMessage.error(error.message || '刷新工单失败'))
    },
    async showDetail(row) {
      try {
        const { data: res } = await this.$api.getWorkOrderDetail(row.type, row.id)
        if (res.code !== 200) {
          throw new Error(res.message || '获取工单详情失败')
        }
        this.detailData = res.data || {}
        this.detailVisible = true
      } catch (error) {
        ElMessage.error(error.message || '获取工单详情失败')
      }
    },
    async approveOrder(row) {
      const action = await ElMessageBox.prompt(`请输入审批意见`, `批准 ${row.title || row.orderNo || '工单'}`, {
        confirmButtonText: '批准',
        cancelButtonText: '取消',
        inputValue: '同意执行'
      }).catch(() => null)
      if (!action) return
      const request = row.type === 'sql'
        ? this.$api.approveSQLWorkOrder(row.id, { comment: action.value })
        : this.$api.approveScriptWorkOrder(row.id, { comment: action.value })
      const { data: res } = await request
      if (res.code !== 200) {
        ElMessage.error(res.message || '工单批准失败')
        return
      }
      ElMessage.success('工单已批准')
      this.refreshAll()
      if (this.detailVisible) this.showDetail({ type: row.type || 'script', id: row.id })
    },
    async rejectOrder(row) {
      const action = await ElMessageBox.prompt(`请输入驳回原因`, `驳回 ${row.title || row.orderNo || '工单'}`, {
        confirmButtonText: '驳回',
        cancelButtonText: '取消',
        inputValue: '风险未确认'
      }).catch(() => null)
      if (!action) return
      const request = row.type === 'sql'
        ? this.$api.rejectSQLWorkOrder(row.id, { comment: action.value })
        : this.$api.rejectScriptWorkOrder(row.id, { comment: action.value })
      const { data: res } = await request
      if (res.code !== 200) {
        ElMessage.error(res.message || '工单驳回失败')
        return
      }
      ElMessage.success('工单已驳回')
      this.refreshAll()
      if (this.detailVisible) this.showDetail({ type: row.type || 'script', id: row.id })
    },
    async executeOrder(row) {
      const action = await ElMessageBox.prompt(`请输入执行备注`, `执行 ${row.title || row.orderNo || '工单'}`, {
        confirmButtonText: '执行',
        cancelButtonText: '取消',
        inputValue: '按审批意见执行'
      }).catch(() => null)
      if (!action) return
      const request = row.type === 'sql'
        ? this.$api.executeSQLWorkOrder(row.id, { comment: action.value })
        : this.$api.executeScriptWorkOrder(row.id, { comment: action.value })
      const { data: res } = await request
      if (res.code !== 200) {
        ElMessage.error(res.message || '工单执行失败')
        return
      }
      ElMessage.success('工单执行成功')
      this.refreshAll()
      this.showDetail({ type: row.type || 'script', id: row.id })
    },
    stringifyValue(value) {
      if (value === null || value === undefined || value === '') return '-'
      if (typeof value === 'object') return JSON.stringify(value)
      return String(value)
    },
    getRecordValue(record, ...keys) {
      if (!record || typeof record !== 'object') return ''
      for (const key of keys) {
        const value = record[key]
        if (value !== null && value !== undefined && value !== '') {
          return value
        }
      }
      return ''
    },
    getOrderRiskText(record) {
      return this.getRecordValue(record, 'riskText', 'risk_text')
    },
    getOrderRiskLevel(record) {
      return this.getRecordValue(record, 'riskLevel', 'risk_level')
    },
    getOrderStatusText(record) {
      return this.getRecordValue(record, 'statusText', 'status_text') || '未知'
    },
    getOrderDetailHint(record) {
      return this.getRecordValue(record, 'detailHint', 'detail_hint') || '-'
    },
    getOrderLink(record, field) {
      if (field === 'aiDiagnosisPath') {
        return this.getRecordValue(record, 'aiDiagnosisPath', 'ai_diagnosis_path')
      }
      if (field === 'knowledgePath') {
        return this.getRecordValue(record, 'knowledgePath', 'knowledge_path')
      }
      return ''
    },
    openInternalPath(path) {
      if (!path) {
        ElMessage.warning('当前没有可打开的联动入口')
        return
      }
      this.$router.push(path)
    },
    typeTagType(type) {
      const map = { quick: 'primary', script: 'warning', service: 'success', sql: 'danger' }
      return map[type] || 'info'
    },
    riskTagType(level) {
      return ({ high: 'danger', medium: 'warning', low: 'success' })[level] || 'info'
    },
    statusTagType(statusText) {
      statusText = String(statusText || '')
      if (statusText.includes('成功')) return 'success'
      if (statusText.includes('失败') || statusText.includes('驳回')) return 'danger'
      if (statusText.includes('取消')) return 'info'
      if (statusText.includes('中')) return 'warning'
      return 'primary'
    }
  },
  mounted() {
    this.refreshAll()
  }
}
</script>

<style scoped>
.workorder-detail-shell {
  min-height: 220px;
}

.detail-action-row {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}
</style>

