<template>
  <div class="alert-notify-page">
    <el-card shadow="hover" class="page-card">
      <template #header>
        <div class="card-header">
          <div>
            <div class="title">告警推送</div>
            <div class="subtitle">查看真实推送日志，并直接维护通知机器人配置。</div>
          </div>
          <el-button type="primary" size="small" icon="Refresh" @click="refreshAll">刷新</el-button>
        </div>
      </template>

      <el-alert
        v-if="summary.totalNotifyRobots === 0"
        title="远端真环境当前未配置通知机器人，可先在左侧新建 webhook 机器人继续联调。"
        type="info"
        :closable="false"
        class="page-alert"
      />

      <el-row :gutter="16" class="summary-row">
        <el-col :xs="24" :sm="12" :lg="6">
          <div class="summary-card ink">
            <div class="summary-label">推送记录</div>
            <div class="summary-value">{{ summary.totalNotifyLogs }}</div>
            <div class="summary-desc">总推送明细行数</div>
          </div>
        </el-col>
        <el-col :xs="24" :sm="12" :lg="6">
          <div class="summary-card green">
            <div class="summary-label">推送成功</div>
            <div class="summary-value">{{ summary.successfulNotifyLogs }}</div>
            <div class="summary-desc">状态为 `success` 的发送记录</div>
          </div>
        </el-col>
        <el-col :xs="24" :sm="12" :lg="6">
          <div class="summary-card red">
            <div class="summary-label">推送失败</div>
            <div class="summary-value">{{ summary.failedNotifyLogs }}</div>
            <div class="summary-desc">需要重点排查的发送失败</div>
          </div>
        </el-col>
        <el-col :xs="24" :sm="12" :lg="6">
          <div class="summary-card teal">
            <div class="summary-label">当前机器人</div>
            <div class="summary-value">{{ summary.totalNotifyRobots }}</div>
            <div class="summary-desc">已配置的通知机器人数量</div>
          </div>
        </el-col>
      </el-row>

      <div class="content-grid">
        <section class="panel">
          <div class="panel-header">
            <div>
              <div class="panel-title">通知机器人</div>
              <div class="panel-subtitle">支持 webhook、飞书、钉钉、企业微信、邮件等通道的统一维护。</div>
            </div>
            <div class="panel-actions">
              <el-button type="primary" size="small" icon="Plus" @click="openCreateRobotDialog">新建机器人</el-button>
            </div>
          </div>

          <el-table
            v-if="robotList.length"
            :data="robotList"
            stripe
            class="data-table"
            empty-text="暂无通知机器人"
          >
            <el-table-column label="名称" prop="name" min-width="150" show-overflow-tooltip />
            <el-table-column label="类型" width="120">
              <template #default="{ row }">
                <el-tag :type="robotTypeTagType(row.type)">
                  {{ robotTypeText(row.type) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="状态" width="90">
              <template #default="{ row }">
                <el-tag :type="row.status === 1 ? 'success' : 'info'">
                  {{ row.status === 1 ? '启用' : '停用' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="目标地址" min-width="200" show-overflow-tooltip>
              <template #default="{ row }">
                {{ row.webhook || row.server || '-' }}
              </template>
            </el-table-column>
            <el-table-column label="操作" width="220" fixed="right">
              <template #default="{ row }">
                <div class="operation-group">
                  <el-button type="primary" link @click="openEditRobotDialog(row)">编辑</el-button>
                  <el-button type="success" link @click="openTestRobotDialog(row)">测试</el-button>
                  <el-button
                    :type="row.status === 1 ? 'warning' : 'success'"
                    link
                    @click="toggleRobotStatus(row)"
                  >
                    {{ row.status === 1 ? '停用' : '启用' }}
                  </el-button>
                  <el-button type="danger" link @click="handleDeleteRobot(row)">删除</el-button>
                </div>
              </template>
            </el-table-column>
          </el-table>

          <el-empty
            v-else
            description="远端真环境当前未配置机器人"
          />
          <alert-manager-receiver-panel ref="alertManagerReceiverPanel" />
        </section>

        <section class="panel">
          <div class="panel-header">
            <div>
              <div class="panel-title">推送日志</div>
              <div class="panel-subtitle">基于 `monitor_webhook_notify_log` 的真实发送明细。</div>
            </div>
          </div>

          <div class="search-section">
            <el-form :model="queryParams" :inline="true" class="search-form">
              <el-form-item label="关键词">
                <el-input
                  v-model="queryParams.keyword"
                  placeholder="搜索告警标题、来源、机器人名"
                  clearable
                  style="width: 260px"
                  @keyup.enter="handleQuery"
                />
              </el-form-item>
              <el-form-item label="状态">
                <el-select v-model="queryParams.status" placeholder="全部" clearable style="width: 140px">
                  <el-option label="成功" value="success" />
                  <el-option label="失败" value="failed" />
                </el-select>
              </el-form-item>
              <el-form-item label="类型">
                <el-select v-model="queryParams.robotType" placeholder="全部" clearable style="width: 160px">
                  <el-option label="飞书" value="feishu" />
                  <el-option label="钉钉" value="dingtalk" />
                  <el-option label="企业微信" value="wechat" />
                  <el-option label="Webhook" value="webhook" />
                  <el-option label="邮件" value="email" />
                  <el-option label="Teams" value="teams" />
                </el-select>
              </el-form-item>
              <el-form-item>
                <el-button type="primary" icon="Search" @click="handleQuery">搜索</el-button>
                <el-button icon="Refresh" @click="resetQuery">重置</el-button>
              </el-form-item>
            </el-form>
          </div>

          <el-table
            v-loading="loading"
            :data="notifyLogList"
            stripe
            class="data-table"
            empty-text="暂无推送日志"
          >
            <el-table-column label="推送时间" prop="createdAt" min-width="150" />
            <el-table-column label="告警标题" prop="alertTitle" min-width="220" show-overflow-tooltip />
            <el-table-column label="机器人" prop="robotName" min-width="140" show-overflow-tooltip />
            <el-table-column label="状态" width="100">
              <template #default="{ row }">
                <el-tag :type="row.status === 'success' ? 'success' : 'danger'" effect="dark">
                  {{ row.status === 'success' ? '成功' : '失败' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="错误信息" prop="errorMsg" min-width="160" show-overflow-tooltip />
          </el-table>

          <div class="pagination-section">
            <el-pagination
              :current-page="queryParams.pageNum"
              :page-size="queryParams.pageSize"
              :page-sizes="[10, 20, 50, 100]"
              :total="total"
              layout="total, sizes, prev, pager, next, jumper"
              @size-change="handleSizeChange"
              @current-change="handleCurrentChange"
            />
          </div>
        </section>
      </div>
    </el-card>

    <el-dialog
      v-model="robotDialogVisible"
      :title="robotDialogMode === 'create' ? '新建通知机器人' : '编辑通知机器人'"
      width="720px"
      destroy-on-close
    >
      <el-form
        ref="robotFormRef"
        :model="robotForm"
        :rules="robotRules"
        label-width="110px"
        class="robot-form"
      >
        <el-form-item label="名称" prop="name">
          <el-input v-model="robotForm.name" placeholder="例如：opsnexus-webhook-test" />
        </el-form-item>
        <el-form-item label="类型" prop="type">
          <el-select v-model="robotForm.type" placeholder="请选择类型" style="width: 100%">
            <el-option label="Webhook" value="webhook" />
            <el-option label="飞书" value="feishu" />
            <el-option label="钉钉" value="dingtalk" />
            <el-option label="企业微信" value="wechat" />
            <el-option label="Teams" value="teams" />
            <el-option label="邮件" value="email" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="robotForm.status">
            <el-radio :label="1">启用</el-radio>
            <el-radio :label="0">停用</el-radio>
          </el-radio-group>
        </el-form-item>

        <template v-if="robotForm.type === 'email'">
          <el-form-item label="SMTP 服务" prop="server">
            <el-input v-model="robotForm.server" placeholder="例如：smtp.qq.com" />
          </el-form-item>
          <el-form-item label="SMTP 端口" prop="port">
            <el-input-number v-model="robotForm.port" :min="1" controls-position="right" />
          </el-form-item>
          <el-form-item label="用户名" prop="username">
            <el-input v-model="robotForm.username" placeholder="请输入邮箱用户名" />
          </el-form-item>
          <el-form-item label="密码" prop="password">
            <el-input v-model="robotForm.password" show-password placeholder="请输入邮箱密码或授权码" />
          </el-form-item>
          <el-form-item label="发件昵称">
            <el-input v-model="robotForm.nickname" placeholder="可选" />
          </el-form-item>
          <el-form-item label="收件地址" prop="webhook">
            <el-input v-model="robotForm.webhook" placeholder="多个收件人可用逗号分隔" />
          </el-form-item>
        </template>

        <template v-else>
          <el-form-item label="Webhook 地址" prop="webhook">
            <el-input v-model="robotForm.webhook" placeholder="请输入 webhook 地址" />
          </el-form-item>
          <el-form-item label="签名/Token">
            <el-input v-model="robotForm.secret" placeholder="可选，按通道填写" />
          </el-form-item>
          <el-form-item label="请求方法">
            <el-select v-model="robotForm.method" placeholder="请选择方法" style="width: 100%">
              <el-option label="POST" value="POST" />
              <el-option label="PUT" value="PUT" />
              <el-option label="PATCH" value="PATCH" />
              <el-option label="GET" value="GET" />
            </el-select>
          </el-form-item>
          <el-form-item label="自定义 Header">
            <el-input
              v-model="robotForm.headers"
              type="textarea"
              :rows="3"
              placeholder='例如：{"Authorization":"Bearer xxx"}'
            />
          </el-form-item>
        </template>

        <el-form-item label="消息模板">
          <el-input
            v-model="robotForm.template"
            type="textarea"
            :rows="4"
            placeholder="可选，支持后续替换 title/content/level/source 等变量"
          />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="robotForm.remark" type="textarea" :rows="3" placeholder="记录用途或接入说明" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="robotDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="robotSubmitting" @click="submitRobotForm">保存</el-button>
      </template>
    </el-dialog>

    <el-dialog
      v-model="testDialogVisible"
      title="测试通知机器人"
      width="620px"
      destroy-on-close
    >
      <el-alert
        type="info"
        :closable="false"
        class="page-alert"
        title="测试消息会复用当前真实派发链路，并写入现有推送日志。"
      />
      <div class="test-robot-meta">
        <div><span>机器人：</span>{{ testForm.robotName || '-' }}</div>
        <div><span>类型：</span>{{ robotTypeText(testForm.robotType) }}</div>
      </div>
      <el-form label-width="110px" class="robot-form">
        <el-form-item label="来源">
          <el-input v-model="testForm.source" placeholder="默认：manual-test" />
        </el-form-item>
        <el-form-item label="标题">
          <el-input v-model="testForm.title" placeholder="留空时自动生成测试标题" />
        </el-form-item>
        <el-form-item label="级别">
          <el-select v-model="testForm.level" style="width: 100%">
            <el-option label="info" value="info" />
            <el-option label="warning" value="warning" />
            <el-option label="critical" value="critical" />
            <el-option label="P1" value="P1" />
            <el-option label="P2" value="P2" />
            <el-option label="P3" value="P3" />
          </el-select>
        </el-form-item>
        <el-form-item label="内容">
          <el-input
            v-model="testForm.content"
            type="textarea"
            :rows="4"
            placeholder="留空时自动生成测试说明"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="testDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="testSubmitting" @click="submitTestForm">发送测试</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import { ElMessage, ElMessageBox } from 'element-plus'
import AlertManagerReceiverPanel from './components/AlertManagerReceiverPanel.vue'

const createEmptyRobotForm = () => ({
  id: null,
  name: '',
  type: 'webhook',
  webhook: '',
  secret: '',
  status: 1,
  remark: '',
  server: '',
  port: 25,
  username: '',
  password: '',
  nickname: '',
  headers: '',
  method: 'POST',
  template: ''
})

const createEmptyTestForm = () => ({
  robotId: null,
  robotName: '',
  robotType: '',
  source: 'manual-test',
  title: '',
  level: 'info',
  content: ''
})

export default {
  name: 'AlertNotify',
  components: {
    AlertManagerReceiverPanel
  },
  data() {
    return {
      loading: false,
      total: 0,
      robotList: [],
      notifyLogList: [],
      robotDialogVisible: false,
      robotDialogMode: 'create',
      robotSubmitting: false,
      testDialogVisible: false,
      testSubmitting: false,
      robotForm: createEmptyRobotForm(),
      testForm: createEmptyTestForm(),
      robotRules: {
        name: [{ required: true, message: '请输入机器人名称', trigger: 'blur' }],
        type: [{ required: true, message: '请选择机器人类型', trigger: 'change' }],
        status: [{ required: true, message: '请选择状态', trigger: 'change' }],
        webhook: [{
          validator: (rule, value, callback) => {
            if (this.robotForm.type !== 'email' && !value) {
              callback(new Error('请输入 webhook 地址'))
              return
            }
            if (this.robotForm.type === 'email' && !value) {
              callback(new Error('请输入收件地址'))
              return
            }
            callback()
          },
          trigger: 'blur'
        }],
        server: [{
          validator: (rule, value, callback) => {
            if (this.robotForm.type === 'email' && !value) {
              callback(new Error('请输入 SMTP 服务地址'))
              return
            }
            callback()
          },
          trigger: 'blur'
        }]
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
      queryParams: {
        keyword: '',
        status: '',
        robotType: '',
        pageNum: 1,
        pageSize: 10
      }
    }
  },
  methods: {
    buildQueryParams() {
      const params = {
        pageNum: this.queryParams.pageNum,
        pageSize: this.queryParams.pageSize
      }
      if (this.queryParams.keyword) params.keyword = this.queryParams.keyword
      if (this.queryParams.status) params.status = this.queryParams.status
      if (this.queryParams.robotType) params.robotType = this.queryParams.robotType
      return params
    },
    async fetchSummary() {
      const { data: res } = await this.$api.getMonitorAlertSummary()
      if (res.code !== 200) {
        throw new Error(res.message || '获取告警推送摘要失败')
      }
      this.summary = res.data || this.summary
    },
    async fetchRobots() {
      const { data: res } = await this.$api.getMonitorAlertNotifyRobotList()
      if (res.code !== 200) {
        throw new Error(res.message || '获取通知机器人列表失败')
      }
      this.robotList = res.data || []
    },
    async fetchNotifyLogs() {
      const { data: res } = await this.$api.queryMonitorAlertNotifyLogList(this.buildQueryParams())
      if (res.code !== 200) {
        throw new Error(res.message || '获取推送日志失败')
      }
      const data = res.data || {}
      this.notifyLogList = data.list || []
      this.total = data.total || 0
    },
    async refreshAll() {
      this.loading = true
      try {
        await Promise.all([this.fetchSummary(), this.fetchRobots(), this.fetchNotifyLogs()])
        this.$refs.alertManagerReceiverPanel?.reload?.()
      } catch (error) {
        ElMessage.error(error.message || '刷新告警推送失败')
      } finally {
        this.loading = false
      }
    },
    handleQuery() {
      this.queryParams.pageNum = 1
      this.fetchNotifyLogs().catch(error => {
        ElMessage.error(error.message || '查询推送日志失败')
      })
    },
    resetQuery() {
      this.queryParams = {
        keyword: '',
        status: '',
        robotType: '',
        pageNum: 1,
        pageSize: 10
      }
      this.fetchNotifyLogs().catch(error => {
        ElMessage.error(error.message || '重置查询失败')
      })
    },
    handleSizeChange(size) {
      this.queryParams.pageSize = size
      this.fetchNotifyLogs().catch(error => {
        ElMessage.error(error.message || '刷新推送日志失败')
      })
    },
    handleCurrentChange(page) {
      this.queryParams.pageNum = page
      this.fetchNotifyLogs().catch(error => {
        ElMessage.error(error.message || '刷新推送日志失败')
      })
    },
    robotTypeText(type) {
      const textMap = {
        feishu: '飞书',
        dingtalk: '钉钉',
        wechat: '企业微信',
        email: '邮件',
        webhook: 'Webhook',
        teams: 'Teams'
      }
      return textMap[type] || type || '未知'
    },
    robotTypeTagType(type) {
      const typeMap = {
        feishu: 'success',
        dingtalk: 'warning',
        wechat: 'primary',
        email: 'info',
        webhook: 'danger',
        teams: ''
      }
      return typeMap[type] || 'info'
    },
    openCreateRobotDialog() {
      this.robotDialogMode = 'create'
      this.robotForm = createEmptyRobotForm()
      this.robotDialogVisible = true
    },
    openEditRobotDialog(row) {
      this.robotDialogMode = 'edit'
      this.robotForm = {
        id: row.id,
        name: row.name || '',
        type: row.type || 'webhook',
        webhook: row.webhook || '',
        secret: row.secret || '',
        status: Number(row.status ?? 1),
        remark: row.remark || '',
        server: row.server || '',
        port: Number(row.port || 25),
        username: row.username || '',
        password: row.password || '',
        nickname: row.nickname || '',
        headers: row.headers || '',
        method: row.method || 'POST',
        template: row.template || ''
      }
      this.robotDialogVisible = true
    },
    openTestRobotDialog(row) {
      this.testForm = {
        ...createEmptyTestForm(),
        robotId: row.id,
        robotName: row.name || '',
        robotType: row.type || ''
      }
      this.testDialogVisible = true
    },
    async submitRobotForm() {
      const valid = await this.$refs.robotFormRef.validate().catch(() => false)
      if (!valid) return

      this.robotSubmitting = true
      try {
        const payload = {
          name: this.robotForm.name,
          type: this.robotForm.type,
          webhook: this.robotForm.webhook,
          secret: this.robotForm.secret,
          status: this.robotForm.status,
          remark: this.robotForm.remark,
          server: this.robotForm.server,
          port: this.robotForm.port,
          username: this.robotForm.username,
          password: this.robotForm.password,
          nickname: this.robotForm.nickname,
          headers: this.robotForm.headers,
          method: this.robotForm.method,
          template: this.robotForm.template
        }

        if (this.robotDialogMode === 'create') {
          await this.$api.createMonitorAlertNotifyRobot(payload)
          ElMessage.success('通知机器人创建成功')
        } else {
          await this.$api.updateMonitorAlertNotifyRobot(this.robotForm.id, payload)
          ElMessage.success('通知机器人更新成功')
        }
        this.robotDialogVisible = false
        await Promise.all([this.fetchSummary(), this.fetchRobots()])
      } catch (error) {
        ElMessage.error(error?.response?.data?.message || error.message || '保存通知机器人失败')
      } finally {
        this.robotSubmitting = false
      }
    },
    async submitTestForm() {
      if (!this.testForm.robotId) {
        ElMessage.error('未选择通知机器人')
        return
      }

      this.testSubmitting = true
      try {
        const payload = {
          source: this.testForm.source,
          title: this.testForm.title,
          level: this.testForm.level,
          content: this.testForm.content
        }
        const { data: res } = await this.$api.testMonitorAlertNotifyRobot(this.testForm.robotId, payload)
        if (res.code !== 200) {
          throw new Error(res.message || '发送测试通知失败')
        }
        const summary = res.data || {}
        const statusText = summary.status === 'success' ? '成功' : summary.status === 'partial' ? '部分成功' : '失败'
        ElMessage.success(`测试发送完成：${statusText}（成功 ${summary.successCount || 0} / 失败 ${summary.failedCount || 0}）`)
        this.testDialogVisible = false
        this.queryParams.pageNum = 1
        await Promise.all([this.fetchSummary(), this.fetchNotifyLogs()])
      } catch (error) {
        ElMessage.error(error?.response?.data?.message || error.message || '发送测试通知失败')
      } finally {
        this.testSubmitting = false
      }
    },
    async toggleRobotStatus(row) {
      const targetStatus = row.status === 1 ? 0 : 1
      const actionText = targetStatus === 1 ? '启用' : '停用'
      try {
        await ElMessageBox.confirm(`确认${actionText}通知机器人“${row.name}”吗？`, '提示', {
          type: 'warning'
        })
        await this.$api.updateMonitorAlertNotifyRobotStatus(row.id, targetStatus)
        ElMessage.success(`通知机器人已${actionText}`)
        await Promise.all([this.fetchSummary(), this.fetchRobots()])
      } catch (error) {
        if (error !== 'cancel' && error !== 'close') {
          ElMessage.error(error?.response?.data?.message || error.message || `${actionText}通知机器人失败`)
        }
      }
    },
    async handleDeleteRobot(row) {
      try {
        await ElMessageBox.confirm(`确认删除通知机器人“${row.name}”吗？`, '警告', {
          type: 'warning'
        })
        await this.$api.deleteMonitorAlertNotifyRobot(row.id)
        ElMessage.success('通知机器人删除成功')
        await Promise.all([this.fetchSummary(), this.fetchRobots()])
      } catch (error) {
        if (error !== 'cancel' && error !== 'close') {
          ElMessage.error(error?.response?.data?.message || error.message || '删除通知机器人失败')
        }
      }
    }
  },
  created() {
    this.refreshAll()
  }
}
</script>

<style scoped>
.alert-notify-page {
  padding: 20px;
  min-height: 100vh;
  background:
    radial-gradient(circle at top left, rgba(20, 184, 166, 0.16), transparent 26%),
    radial-gradient(circle at top right, rgba(37, 99, 235, 0.14), transparent 24%),
    linear-gradient(145deg, #f4fffe 0%, #f7fbff 52%, #fbfffd 100%);
}

.page-card {
  border-radius: 20px;
  border: 1px solid rgba(148, 163, 184, 0.16);
  box-shadow: 0 18px 46px rgba(15, 23, 42, 0.08);
}

.card-header,
.panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
}

.panel-actions {
  display: flex;
  gap: 12px;
  align-items: center;
}

.title {
  font-size: 24px;
  font-weight: 700;
  color: #0f172a;
}

.subtitle,
.panel-subtitle {
  margin-top: 6px;
  font-size: 13px;
  color: #64748b;
}

.page-alert {
  margin-bottom: 18px;
}

.summary-row {
  margin-bottom: 20px;
}

.summary-card {
  min-height: 126px;
  padding: 18px;
  border-radius: 18px;
  color: #fff;
  box-shadow: 0 12px 28px rgba(15, 23, 42, 0.1);
}

.summary-card.ink {
  background: linear-gradient(135deg, #334155 0%, #1e293b 100%);
}

.summary-card.green {
  background: linear-gradient(135deg, #059669 0%, #047857 100%);
}

.summary-card.red {
  background: linear-gradient(135deg, #dc2626 0%, #b91c1c 100%);
}

.summary-card.teal {
  background: linear-gradient(135deg, #0f766e 0%, #115e59 100%);
}

.summary-label {
  font-size: 14px;
  opacity: 0.92;
}

.summary-value {
  margin-top: 12px;
  font-size: 34px;
  font-weight: 700;
  line-height: 1;
}

.summary-desc {
  margin-top: 14px;
  font-size: 12px;
  opacity: 0.86;
}

.content-grid {
  display: grid;
  grid-template-columns: minmax(280px, 0.82fr) minmax(0, 1.58fr);
  gap: 20px;
}

.panel {
  padding: 20px;
  border-radius: 18px;
  background: rgba(255, 255, 255, 0.84);
  border: 1px solid rgba(148, 163, 184, 0.12);
}

.panel-title {
  font-size: 18px;
  font-weight: 700;
  color: #0f172a;
}

.search-section {
  margin: 18px 0;
  padding: 18px 18px 4px;
  border-radius: 16px;
  background: linear-gradient(135deg, rgba(20, 184, 166, 0.08), rgba(37, 99, 235, 0.06));
  border: 1px solid rgba(148, 163, 184, 0.12);
}

.search-form {
  display: flex;
  flex-wrap: wrap;
}

.data-table {
  border-radius: 14px;
  overflow: hidden;
}

.data-table :deep(.el-table__header th) {
  background: #f8fafc !important;
  color: #0f172a;
  font-weight: 700;
}

.operation-group {
  display: flex;
  gap: 8px;
  align-items: center;
  flex-wrap: wrap;
}

.pagination-section {
  display: flex;
  justify-content: center;
  padding-top: 20px;
}

.test-robot-meta {
  margin: 14px 0 18px;
  padding: 12px 14px;
  border-radius: 12px;
  background: rgba(15, 118, 110, 0.06);
  color: #0f172a;
  display: grid;
  gap: 8px;
}

.test-robot-meta span {
  color: #64748b;
}

.el-button {
  border-radius: 10px;
}

.robot-form :deep(.el-input-number) {
  width: 100%;
}

@media (max-width: 1200px) {
  .content-grid {
    grid-template-columns: 1fr;
  }
}
</style>
