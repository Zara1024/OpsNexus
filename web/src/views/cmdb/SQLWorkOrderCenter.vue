<template>
  <TablePage
    class="sql-work-order-page"
    :section-title="sectionTitle"
    :section-subtitle="sectionSubtitle"
    wide
  >
    <template #header>
      <PageHeader
        eyebrow="Database Asset Workflow"
        title="SQL 工单中心"
        :subtitle="pageSubtitle"
      >
        <template #actions>
          <el-button @click="refreshAll">刷新</el-button>
        </template>
        <template #intro>
          <PageIntro
            title="处理建议"
            text="优先查看风险摘要、受影响对象和回滚建议，再决定审批、驳回或执行，避免数据库资产操作链路碎片化。"
          />
        </template>
      </PageHeader>
    </template>

    <template #toolbar>
      <PageToolbar>
        <el-form :inline="true" class="filter-cluster">
          <el-form-item class="filter-field" label="关键字">
            <el-input
              v-model="queryParams.keyword"
              placeholder="搜索工单号、标题或数据库"
              clearable
              @keyup.enter="handleQuery"
            />
          </el-form-item>
          <el-form-item class="filter-field" label="状态">
            <el-select v-model="queryParams.status" placeholder="全部" clearable>
              <el-option
                v-for="item in statusOptions"
                :key="item.value"
                :label="item.label"
                :value="item.value"
              />
            </el-select>
          </el-form-item>
        </el-form>
        <template #actions>
          <el-button type="primary" @click="handleQuery">搜索</el-button>
          <el-button @click="resetQuery">重置</el-button>
        </template>
      </PageToolbar>
    </template>

    <template v-if="showStatStrip" #stats>
      <StatStrip :items="statItems" />
    </template>

    <el-table
      v-loading="loading"
      :data="orderList"
      stripe
      class="data-table"
      :empty-text="emptyText"
    >
      <el-table-column label="工单号" prop="orderNo" min-width="140" />
      <el-table-column label="标题" prop="title" min-width="200" show-overflow-tooltip />
      <el-table-column label="数据库" prop="databaseName" min-width="120" />
      <el-table-column label="类型" width="100">
        <template #default="{ row }">
          <el-tag>{{ row.operationType }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="风险" width="100">
        <template #default="{ row }">
          <el-tag :type="riskTagType(row.riskLevel)" effect="dark">{{ riskLabel(row.riskLevel) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="状态" width="120">
        <template #default="{ row }">
          <el-tag :type="statusTagType(row.status)" effect="dark">{{ statusLabel(row.status) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="180" fixed="right">
        <template #default="{ row }">
          <div class="row-actions">
            <el-button type="primary" link @click="showDetail(row.id)">详情</el-button>
            <el-button v-if="row.status === 1" type="success" link @click="approveOrder(row)">批准</el-button>
            <el-button v-if="row.status === 1" type="danger" link @click="rejectOrder(row)">驳回</el-button>
            <el-button v-if="row.status === 2" type="warning" link @click="executeOrder(row)">执行</el-button>
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

    <el-drawer
      v-model="detailVisible"
      size="70%"
      title="数据库 SQL 工单详情"
      class="sql-work-order-drawer"
    >
      <div v-loading="detailLoading">
        <template v-if="detailData.id">
          <el-descriptions :column="2" border>
            <el-descriptions-item label="工单号">{{ detailData.orderNo }}</el-descriptions-item>
            <el-descriptions-item label="状态">{{ statusLabel(detailData.status) }}</el-descriptions-item>
            <el-descriptions-item label="数据库">{{ detailData.databaseName }}</el-descriptions-item>
            <el-descriptions-item label="实例">{{ detailData.instanceHost || '-' }}</el-descriptions-item>
            <el-descriptions-item label="申请人">{{ detailData.applicantName || '-' }}</el-descriptions-item>
            <el-descriptions-item label="审批人">{{ detailData.approverName || '-' }}</el-descriptions-item>
            <el-descriptions-item label="风险">{{ riskLabel(detailData.riskLevel) }}</el-descriptions-item>
            <el-descriptions-item label="影响对象">{{ detailData.affectedTables || '-' }}</el-descriptions-item>
          </el-descriptions>

          <el-alert
            :title="detailData.riskSummary || '暂无风险提示'"
            :type="detailData.riskLevel >= 2 ? 'error' : 'warning'"
            :closable="false"
            class="detail-alert"
          />

          <div class="detail-block">
            <div class="block-title">SQL 内容</div>
            <pre class="code-block">{{ detailData.sqlContent || '-' }}</pre>
          </div>

          <div class="detail-block">
            <div class="block-title">回滚建议</div>
            <pre class="code-block">{{ detailData.rollbackSql || detailData.rollbackHint || '-' }}</pre>
          </div>

          <div v-if="detailData.backupPreview" class="detail-block">
            <div class="block-title">执行前快照</div>
            <pre class="code-block">{{ detailData.backupPreview }}</pre>
          </div>

          <div class="detail-block">
            <div class="block-title">执行结果</div>
            <pre class="code-block">{{ detailData.resultMessage || '-' }}</pre>
          </div>
        </template>
      </div>
    </el-drawer>
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
  name: 'SQLWorkOrderCenter',
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
      detailVisible: false,
      total: 0,
      summary: { total: 0, pending: 0, approved: 0, succeeded: 0, highRisk: 0 },
      orderList: [],
      detailData: {},
      queryParams: {
        page: 1,
        pageSize: 10,
        status: '',
        keyword: '',
        databaseId: Number(this.$route.query.databaseId || 0) || undefined
      },
      statusOptions: [
        { label: '待审批', value: 1 },
        { label: '已批准', value: 2 },
        { label: '已驳回', value: 3 },
        { label: '执行中', value: 4 },
        { label: '执行成功', value: 5 },
        { label: '执行失败', value: 6 }
      ]
    }
  },
  computed: {
    isDatabaseScoped() {
      return Boolean(this.queryParams.databaseId)
    },
    showStatStrip() {
      return !this.isDatabaseScoped
    },
    sectionTitle() {
      return this.isDatabaseScoped ? '指定数据库资产 SQL 工单' : '数据库 SQL 工单'
    },
    sectionSubtitle() {
      return this.isDatabaseScoped
        ? '当前列表已按指定数据库资产筛选，便于聚焦审批、执行与审计。'
        : '集中查看所有数据库资产的审批、执行与审计记录。'
    },
    pageSubtitle() {
      return this.isDatabaseScoped
        ? '当前视图已绑定到某个数据库资产，适合快速追踪该资产的变更链路。'
        : '把数据库资产相关的审批、执行、审计和回滚建议集中到一个工作台。'
    },
    emptyText() {
      return this.isDatabaseScoped ? '当前数据库资产暂无 SQL 工单' : '暂无 SQL 工单'
    },
    statItems() {
      return [
        { label: '工单总量', value: this.summary.total, hint: '当前可见范围内的 SQL 工单', tone: 'primary' },
        { label: '待审批', value: this.summary.pending, hint: '等待人工审批', tone: 'warning' },
        { label: '已批准', value: this.summary.approved, hint: '可进入执行阶段', tone: 'success' },
        { label: '执行成功', value: this.summary.succeeded, hint: '执行与审计已完成', tone: 'success' },
        { label: '高风险', value: this.summary.highRisk, hint: '需要重点复核的数据库变更', tone: 'danger' }
      ]
    }
  },
  methods: {
    buildParams() {
      const params = { page: this.queryParams.page, pageSize: this.queryParams.pageSize }
      if (this.queryParams.status) params.status = this.queryParams.status
      if (this.queryParams.keyword) params.keyword = this.queryParams.keyword
      if (this.queryParams.databaseId) params.databaseId = this.queryParams.databaseId
      return params
    },
    async fetchSummary() {
      const { data: res } = await this.$api.getSQLWorkOrderSummary()
      if (res.code !== 200) {
        throw new Error(res.message || '获取 SQL 工单汇总失败')
      }
      this.summary = res.data || this.summary
    },
    async fetchList() {
      const { data: res } = await this.$api.getSQLWorkOrderList(this.buildParams())
      if (res.code !== 200) {
        throw new Error(res.message || '获取 SQL 工单列表失败')
      }
      this.orderList = res.data?.list || []
      this.total = res.data?.total || 0
    },
    async refreshAll() {
      this.loading = true
      try {
        if (this.isDatabaseScoped) {
          await this.fetchList()
        } else {
          await Promise.all([this.fetchSummary(), this.fetchList()])
        }
      } catch (error) {
        ElMessage.error(error.message || '刷新 SQL 工单失败')
      } finally {
        this.loading = false
      }
    },
    async showDetail(id) {
      this.detailVisible = true
      this.detailLoading = true
      try {
        const { data: res } = await this.$api.getSQLWorkOrderDetail(id)
        if (res.code !== 200) {
          throw new Error(res.message || '获取 SQL 工单详情失败')
        }
        this.detailData = res.data || {}
      } catch (error) {
        ElMessage.error(error.message || '获取 SQL 工单详情失败')
      } finally {
        this.detailLoading = false
      }
    },
    async approveOrder(row) {
      const { value } = await ElMessageBox.prompt('审批意见', `批准工单 ${row.orderNo}`, {
        inputValue: '同意执行',
        confirmButtonText: '批准',
        cancelButtonText: '取消'
      }).catch(() => ({ value: null }))

      if (value === null) {
        return
      }

      const { data: res } = await this.$api.approveSQLWorkOrder(row.id, { comment: value })
      if (res.code !== 200) {
        return ElMessage.error(res.message || '审批失败')
      }

      ElMessage.success('工单已批准')
      this.refreshAll()
      if (this.detailVisible && this.detailData.id === row.id) {
        this.showDetail(row.id)
      }
    },
    async rejectOrder(row) {
      const { value } = await ElMessageBox.prompt('驳回原因', `驳回工单 ${row.orderNo}`, {
        inputValue: '风险未评估完成',
        confirmButtonText: '驳回',
        cancelButtonText: '取消'
      }).catch(() => ({ value: null }))

      if (value === null) {
        return
      }

      const { data: res } = await this.$api.rejectSQLWorkOrder(row.id, { comment: value })
      if (res.code !== 200) {
        return ElMessage.error(res.message || '驳回失败')
      }

      ElMessage.success('工单已驳回')
      this.refreshAll()
      if (this.detailVisible && this.detailData.id === row.id) {
        this.showDetail(row.id)
      }
    },
    async executeOrder(row) {
      const { value } = await ElMessageBox.prompt('执行备注', `执行工单 ${row.orderNo}`, {
        inputValue: '按审批意见执行',
        confirmButtonText: '执行',
        cancelButtonText: '取消'
      }).catch(() => ({ value: null }))

      if (value === null) {
        return
      }

      const { data: res } = await this.$api.executeSQLWorkOrder(row.id, { comment: value })
      if (res.code !== 200) {
        return ElMessage.error(res.message || '执行失败')
      }

      ElMessage.success('工单执行成功')
      this.refreshAll()
      this.showDetail(row.id)
    },
    handleQuery() {
      this.queryParams.page = 1
      this.refreshAll()
    },
    resetQuery() {
      this.queryParams = {
        page: 1,
        pageSize: 10,
        status: '',
        keyword: '',
        databaseId: this.queryParams.databaseId
      }
      this.refreshAll()
    },
    handleSizeChange(size) {
      this.queryParams.pageSize = size
      this.queryParams.page = 1
      this.fetchList().catch((error) => ElMessage.error(error.message || '刷新失败'))
    },
    handleCurrentChange(page) {
      this.queryParams.page = page
      this.fetchList().catch((error) => ElMessage.error(error.message || '刷新失败'))
    },
    statusLabel(status) {
      return {
        1: '待审批',
        2: '已批准',
        3: '已驳回',
        4: '执行中',
        5: '执行成功',
        6: '执行失败'
      }[status] || '未知'
    },
    statusTagType(status) {
      return {
        1: 'warning',
        2: 'primary',
        3: 'info',
        4: 'warning',
        5: 'success',
        6: 'danger'
      }[status] || 'info'
    },
    riskLabel(level) {
      return { 0: '低', 1: '中', 2: '高' }[Number(level)] || '未知'
    },
    riskTagType(level) {
      return Number(level) >= 2 ? 'danger' : (Number(level) === 1 ? 'warning' : 'success')
    }
  },
  created() {
    this.refreshAll()
  }
}
</script>

<style scoped>
.row-actions {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 10px;
}

.sql-work-order-page :deep(.section-card__footer) {
  display: flex;
  justify-content: center;
}

.sql-work-order-page :deep(.el-table .cell.el-tooltip) {
  white-space: normal;
  word-break: break-word;
  line-height: 1.5;
}

.detail-alert,
.detail-block {
  margin-top: 16px;
}

.block-title {
  margin-bottom: 8px;
  font-weight: 700;
  color: var(--text-primary);
}

.code-block {
  margin: 0;
  padding: 16px;
  border-radius: 16px;
  background: rgba(3, 9, 18, 0.96);
  border: 1px solid var(--border-subtle);
  color: #e2e8f0;
  white-space: pre-wrap;
  word-break: break-word;
  font-family: Consolas, Monaco, monospace;
  line-height: 1.6;
}

@media (max-width: 960px) {
  .row-actions {
    gap: 8px;
  }
}
</style>
