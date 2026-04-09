<template>
  <div class="monitor-automation-page">
    <el-card shadow="hover" class="page-card">
      <template #header>
        <div class="header-row">
          <div>
            <div class="title">监控深化</div>
            <div class="subtitle">主机告警、数据库告警与 SSL 自动化统一工作台。</div>
          </div>
          <div class="header-actions">
            <el-button type="primary" @click="refreshAll">刷新</el-button>
            <el-button @click="scanHostAlerts">扫描主机</el-button>
            <el-button @click="scanDBAlerts">扫描数据库</el-button>
            <el-button @click="scanSSLDomains">扫描 SSL</el-button>
          </div>
        </div>
      </template>

      <div class="summary-grid">
        <div class="summary-card red"><div class="summary-label">开放事件</div><div class="summary-value">{{ overview.openEventCount }}</div></div>
        <div class="summary-card blue"><div class="summary-label">主机规则</div><div class="summary-value">{{ overview.hostRuleCount }}</div></div>
        <div class="summary-card teal"><div class="summary-label">数据库健康</div><div class="summary-value">{{ overview.databaseHealthyCount }}/{{ overview.databaseTotalCount }}</div></div>
        <div class="summary-card amber"><div class="summary-label">即将过期域名</div><div class="summary-value">{{ overview.expiringDomainCount }}</div></div>
        <div class="summary-card slate"><div class="summary-label">启用机器人</div><div class="summary-value">{{ overview.enabledRobotCount || 0 }}</div></div>
        <div class="summary-card wine"><div class="summary-label">通知失败</div><div class="summary-value">{{ overview.failedNotifyLogCount || 0 }}</div></div>
      </div>

      <div class="workbench-grid">
        <div class="workbench-panel">
          <div class="workbench-title">风险提示</div>
          <el-empty v-if="!overview.riskTips?.length" description="当前没有风险提示" :image-size="56" />
          <div v-else class="workbench-list">
            <div v-for="(item, index) in overview.riskTips" :key="`risk-${index}`" class="workbench-item risk">
              <div class="workbench-item-title">风险 {{ index + 1 }}</div>
              <div class="workbench-item-summary">{{ item }}</div>
            </div>
          </div>
        </div>

        <div class="workbench-panel">
          <div class="workbench-title">推荐动作</div>
          <div class="workbench-list">
            <div v-for="item in overview.recommendedActions || []" :key="item.title" class="workbench-item">
              <div>
                <div class="workbench-item-title">{{ item.title }}</div>
                <div class="workbench-item-summary">{{ item.summary }}</div>
              </div>
              <el-button type="primary" link @click="openWorkbenchPath(item.path)">去处理</el-button>
            </div>
          </div>
        </div>

        <div class="workbench-panel">
          <div class="workbench-title">最近事件</div>
          <el-empty v-if="!overview.recentEvents?.length" description="暂无事件" :image-size="56" />
          <div v-else class="workbench-list">
            <div v-for="item in overview.recentEvents" :key="`${item.title}-${item.occurredAt}`" class="workbench-item">
              <div>
                <div class="workbench-item-title">{{ item.title }}</div>
                <div class="workbench-item-summary">{{ item.summary || '-' }}</div>
              </div>
              <div class="workbench-item-meta">
                <el-tag :type="severityTagType(item.severity)" size="small">{{ item.severity || '-' }}</el-tag>
                <span>{{ item.occurredAt || '-' }}</span>
              </div>
            </div>
          </div>
        </div>

        <div class="workbench-panel">
          <div class="workbench-title">最近动作</div>
          <el-empty v-if="!overview.recentActions?.length" description="暂无动作" :image-size="56" />
          <div v-else class="workbench-list">
            <div v-for="item in overview.recentActions" :key="`${item.title}-${item.path}`" class="workbench-item">
              <div>
                <div class="workbench-item-title">{{ item.title }}</div>
                <div class="workbench-item-summary">{{ item.summary }}</div>
              </div>
              <el-button type="primary" link @click="openWorkbenchPath(item.path)">查看</el-button>
            </div>
          </div>
        </div>
      </div>

      <el-tabs v-model="activeTab">
        <el-tab-pane label="主机告警" name="host">
          <div class="section-header">
            <div class="section-title">主机规则</div>
            <el-button type="primary" size="small" @click="openHostRuleDialog()">新建主机规则</el-button>
          </div>
          <div class="template-grid">
            <div v-for="item in hostTemplates" :key="item.metricKey" class="template-card" @click="applyHostTemplate(item)">
              <div class="template-name">{{ item.name }}</div>
              <div class="template-detail">{{ item.description }}</div>
            </div>
          </div>
          <el-table :data="hostRules" stripe class="data-table">
            <el-table-column label="规则名称" prop="name" min-width="180" />
            <el-table-column label="主机" prop="hostName" min-width="180" />
            <el-table-column label="指标" prop="metricKey" width="140" />
            <el-table-column label="阈值" width="120"><template #default="{ row }">{{ row.operator }} {{ row.thresholdValue }}</template></el-table-column>
            <el-table-column label="等级" width="100"><template #default="{ row }"><el-tag :type="severityTagType(row.severity)">{{ row.severity }}</el-tag></template></el-table-column>
            <el-table-column label="操作" width="180" fixed="right">
              <template #default="{ row }">
                <div class="table-actions">
                  <el-button type="primary" link @click="openHostRuleDialog(row)">编辑</el-button>
                  <el-button :type="row.status === 1 ? 'warning' : 'success'" link @click="toggleHostRule(row)">{{ row.status === 1 ? '停用' : '启用' }}</el-button>
                  <el-button type="danger" link @click="deleteHostRule(row)">删除</el-button>
                </div>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <el-tab-pane label="数据库告警" name="db">
          <div class="section-header">
            <div class="section-title">数据库规则</div>
            <el-button type="primary" size="small" @click="openDBRuleDialog()">新建数据库规则</el-button>
          </div>
          <el-table :data="dbRules" stripe class="data-table">
            <el-table-column label="规则名称" prop="name" min-width="180" />
            <el-table-column label="数据库" prop="databaseName" min-width="180" />
            <el-table-column label="类型" width="120"><template #default="{ row }">{{ databaseTypeText(row.databaseType) }}</template></el-table-column>
            <el-table-column label="指标" prop="metricKey" width="140" />
            <el-table-column label="操作" width="180" fixed="right">
              <template #default="{ row }">
                <div class="table-actions">
                  <el-button type="primary" link @click="openDBRuleDialog(row)">编辑</el-button>
                  <el-button :type="row.status === 1 ? 'warning' : 'success'" link @click="toggleDBRule(row)">{{ row.status === 1 ? '停用' : '启用' }}</el-button>
                  <el-button type="danger" link @click="deleteDBRule(row)">删除</el-button>
                </div>
              </template>
            </el-table-column>
          </el-table>
          <el-table :data="dbSnapshots" stripe class="data-table">
            <el-table-column label="数据库" prop="databaseName" min-width="180" />
            <el-table-column label="地址" min-width="180"><template #default="{ row }">{{ row.host }}:{{ row.port }}</template></el-table-column>
            <el-table-column label="状态" width="100"><template #default="{ row }"><el-tag :type="row.available === 1 ? 'success' : 'danger'">{{ row.available === 1 ? '健康' : '异常' }}</el-tag></template></el-table-column>
            <el-table-column label="延迟(ms)" prop="latencyMs" width="120" />
            <el-table-column label="错误信息" prop="errorMsg" min-width="220" show-overflow-tooltip />
          </el-table>
        </el-tab-pane>

        <el-tab-pane label="SSL 自动化" name="ssl">
          <div class="section-header">
            <div class="section-title">域名巡检与证书部署</div>
            <div class="header-actions">
              <el-button type="primary" size="small" @click="openDomainDialog()">新增域名</el-button>
              <el-button size="small" @click="openDeployDialog()">手动部署证书</el-button>
            </div>
          </div>
          <el-table :data="domains" stripe class="data-table">
            <el-table-column label="域名" prop="domain" min-width="180" />
            <el-table-column label="存活" width="100"><template #default="{ row }"><el-tag :type="row.isAlive === 1 ? 'success' : 'danger'">{{ row.isAlive === 1 ? '正常' : '异常' }}</el-tag></template></el-table-column>
            <el-table-column label="状态" width="180">
              <template #default="{ row }">
                <div class="table-actions">
                  <el-tag :type="row.status === 1 ? 'success' : 'info'">{{ row.status === 1 ? '启用' : '停用' }}</el-tag>
                  <el-button :type="row.status === 1 ? 'warning' : 'success'" link @click="toggleDomain(row)">{{ row.status === 1 ? '停用' : '启用' }}</el-button>
                </div>
              </template>
            </el-table-column>
            <el-table-column label="HTTP" prop="statusCode" width="100" />
            <el-table-column label="操作" width="120" fixed="right">
              <template #default="{ row }">
                <div class="table-actions">
                  <el-button type="primary" link @click="openDomainDialog(row)">编辑</el-button>
                  <el-button type="danger" link @click="deleteDomain(row)">删除</el-button>
                </div>
              </template>
            </el-table-column>
            <el-table-column label="剩余天数" prop="sslDaysLeft" width="120" />
            <el-table-column label="错误信息" prop="errorMsg" min-width="220" show-overflow-tooltip />
          </el-table>
          <el-table :data="schedules" stripe class="data-table">
            <el-table-column label="Cron" prop="cronExpr" min-width="160" />
            <el-table-column label="过期阈值(天)" prop="expireAlertDays" width="120" />
            <el-table-column label="通知" width="120"><template #default="{ row }">{{ row.notifyEnabled ? robotNameById(row.notifyRobotId) : '关闭' }}</template></el-table-column>
            <el-table-column label="操作" width="120"><template #default="{ row }"><el-button type="primary" link @click="openScheduleDialog(row)">编辑</el-button></template></el-table-column>
          </el-table>
          <el-table :data="certs" stripe class="data-table">
            <el-table-column label="域名" prop="domain" min-width="180" />
            <el-table-column label="CA" prop="caProvider" width="120" />
            <el-table-column label="状态" width="100"><template #default="{ row }"><el-tag :type="certStatusTagType(row.status)">{{ certStatusText(row.status) }}</el-tag></template></el-table-column>
            <el-table-column label="剩余天数" prop="daysLeft" width="120" />
            <el-table-column label="到期时间" min-width="180"><template #default="{ row }">{{ formatTime(row.expireTime) }}</template></el-table-column>
          </el-table>
          <el-table :data="deployLogs" stripe class="data-table">
            <el-table-column label="域名" prop="domain" min-width="160" />
            <el-table-column label="目标主机" prop="hostName" min-width="160" />
            <el-table-column label="状态" width="100"><template #default="{ row }"><el-tag :type="deployStatusTagType(row.status)">{{ deployStatusText(row.status) }}</el-tag></template></el-table-column>
            <el-table-column label="失败原因" prop="errorMsg" min-width="220" show-overflow-tooltip />
          </el-table>
        </el-tab-pane>

        <el-tab-pane label="监控事件" name="events">
          <el-table :data="events" stripe class="data-table">
            <el-table-column label="资源类型" prop="resourceType" width="120" />
            <el-table-column label="资源名称" prop="resourceName" min-width="180" />
            <el-table-column label="标题" prop="title" min-width="220" show-overflow-tooltip />
            <el-table-column label="等级" width="100"><template #default="{ row }"><el-tag :type="severityTagType(row.severity)">{{ row.severity }}</el-tag></template></el-table-column>
            <el-table-column label="状态" width="100"><template #default="{ row }"><el-tag :type="row.status === 'open' ? 'danger' : 'success'">{{ row.status === 'open' ? '触发中' : '已恢复' }}</el-tag></template></el-table-column>
            <el-table-column label="摘要" prop="summary" min-width="260" show-overflow-tooltip />
          </el-table>
        </el-tab-pane>
      </el-tabs>
    </el-card>

    <el-dialog v-model="hostRuleDialogVisible" :title="hostRuleForm.id ? '编辑主机规则' : '新建主机规则'" width="620px">
      <el-form :model="hostRuleForm" label-width="110px">
        <el-form-item label="规则名称"><el-input v-model="hostRuleForm.name" /></el-form-item>
        <el-form-item label="主机"><el-select v-model="hostRuleForm.hostId" filterable style="width: 100%"><el-option v-for="item in hosts" :key="item.id" :label="`${item.hostName} (${item.sshIp || item.privateIp || item.publicIp || '-'})`" :value="item.id" /></el-select></el-form-item>
        <el-form-item label="指标"><el-select v-model="hostRuleForm.metricKey" style="width: 100%"><el-option v-for="item in hostTemplates" :key="item.metricKey" :label="item.name" :value="item.metricKey" /></el-select></el-form-item>
        <el-form-item label="比较方式"><el-select v-model="hostRuleForm.operator" style="width: 100%"><el-option v-for="item in operators" :key="item" :label="item" :value="item" /></el-select></el-form-item>
        <el-form-item label="阈值"><el-input-number v-model="hostRuleForm.thresholdValue" :min="0" :max="99999" style="width: 100%" /></el-form-item>
        <el-form-item label="等级"><el-select v-model="hostRuleForm.severity" style="width: 100%"><el-option v-for="item in severities" :key="item" :label="item" :value="item" /></el-select></el-form-item>
        <el-form-item label="通知机器人"><el-select v-model="hostRuleForm.notifyRobotIds" multiple filterable style="width: 100%"><el-option v-for="item in robots" :key="item.id" :label="item.name" :value="item.id" /></el-select></el-form-item>
        <el-form-item label="备注"><el-input v-model="hostRuleForm.remark" type="textarea" :rows="3" /></el-form-item>
      </el-form>
      <template #footer><el-button @click="hostRuleDialogVisible = false">取消</el-button><el-button type="primary" @click="submitHostRule">保存</el-button></template>
    </el-dialog>

    <el-dialog v-model="dbRuleDialogVisible" :title="dbRuleForm.id ? '编辑数据库规则' : '新建数据库规则'" width="620px">
      <el-form :model="dbRuleForm" label-width="110px">
        <el-form-item label="规则名称"><el-input v-model="dbRuleForm.name" /></el-form-item>
        <el-form-item label="数据库"><el-select v-model="dbRuleForm.databaseId" filterable style="width: 100%"><el-option v-for="item in databases" :key="item.id" :label="`${item.name} (${databaseTypeText(item.type)})`" :value="item.id" /></el-select></el-form-item>
        <el-form-item label="指标"><el-select v-model="dbRuleForm.metricKey" style="width: 100%"><el-option label="连通性" value="connectivity" /><el-option label="连接延迟(ms)" value="latency_ms" /></el-select></el-form-item>
        <el-form-item label="比较方式"><el-select v-model="dbRuleForm.operator" style="width: 100%"><el-option v-for="item in operators" :key="item" :label="item" :value="item" /></el-select></el-form-item>
        <el-form-item label="阈值"><el-input-number v-model="dbRuleForm.thresholdValue" :min="0" :max="99999" style="width: 100%" /></el-form-item>
        <el-form-item label="等级"><el-select v-model="dbRuleForm.severity" style="width: 100%"><el-option v-for="item in severities" :key="item" :label="item" :value="item" /></el-select></el-form-item>
        <el-form-item label="通知机器人"><el-select v-model="dbRuleForm.notifyRobotIds" multiple filterable style="width: 100%"><el-option v-for="item in robots" :key="item.id" :label="item.name" :value="item.id" /></el-select></el-form-item>
        <el-form-item label="备注"><el-input v-model="dbRuleForm.remark" type="textarea" :rows="3" /></el-form-item>
      </el-form>
      <template #footer><el-button @click="dbRuleDialogVisible = false">取消</el-button><el-button type="primary" @click="submitDBRule">保存</el-button></template>
    </el-dialog>

    <el-dialog v-model="domainDialogVisible" :title="domainForm.id ? '编辑域名' : '新增域名'" width="560px">
      <el-form :model="domainForm" label-width="110px">
        <el-form-item label="域名"><el-input v-model="domainForm.domain" placeholder="例如：deviops.cn" /></el-form-item>
        <el-form-item label="标签"><el-input v-model="domainForm.tags" /></el-form-item>
        <el-form-item label="备注"><el-input v-model="domainForm.remark" type="textarea" :rows="3" /></el-form-item>
      </el-form>
      <template #footer><el-button @click="domainDialogVisible = false">取消</el-button><el-button type="primary" @click="submitDomain">保存</el-button></template>
    </el-dialog>

    <el-dialog v-model="scheduleDialogVisible" title="编辑巡检调度" width="620px">
      <el-form :model="scheduleForm" label-width="120px">
        <el-form-item label="启用"><el-switch v-model="scheduleForm.enabled" /></el-form-item>
        <el-form-item label="Cron 表达式"><el-input v-model="scheduleForm.cronExpr" placeholder="例如：29 13 * * *" /></el-form-item>
        <el-form-item label="过期阈值(天)"><el-input-number v-model="scheduleForm.expireAlertDays" :min="1" :max="365" style="width: 100%" /></el-form-item>
        <el-form-item label="超时(ms)"><el-input-number v-model="scheduleForm.scanTimeoutMs" :min="1000" :max="60000" style="width: 100%" /></el-form-item>
        <el-form-item label="通知开关"><el-switch v-model="scheduleForm.notifyEnabled" /></el-form-item>
        <el-form-item label="通知机器人"><el-select v-model="scheduleForm.notifyRobotId" clearable style="width: 100%"><el-option v-for="item in robots" :key="item.id" :label="item.name" :value="item.id" /></el-select></el-form-item>
        <el-form-item label="自动部署"><el-switch v-model="scheduleForm.autoDeployEnabled" /></el-form-item>
        <el-form-item label="部署主机"><el-select v-model="scheduleForm.deployHostId" clearable filterable style="width: 100%"><el-option v-for="item in hosts" :key="item.id" :label="item.hostName" :value="item.id" /></el-select></el-form-item>
        <el-form-item label="部署路径"><el-input v-model="scheduleForm.deployPath" /></el-form-item>
        <el-form-item label="重载命令"><el-input v-model="scheduleForm.reloadCommand" /></el-form-item>
      </el-form>
      <template #footer><el-button @click="scheduleDialogVisible = false">取消</el-button><el-button type="primary" @click="submitSchedule">保存</el-button></template>
    </el-dialog>

    <el-dialog v-model="deployDialogVisible" title="手动部署证书" width="620px">
      <el-form :model="deployForm" label-width="120px">
        <el-form-item label="证书"><el-select v-model="deployForm.certId" filterable style="width: 100%"><el-option v-for="item in certs" :key="item.id" :label="`${item.domain} (${item.caProvider || item.certSource || '-'})`" :value="item.id" /></el-select></el-form-item>
        <el-form-item label="目标主机"><el-select v-model="deployForm.hostId" filterable style="width: 100%"><el-option v-for="item in hosts" :key="item.id" :label="item.hostName" :value="item.id" /></el-select></el-form-item>
        <el-form-item label="部署路径"><el-input v-model="deployForm.deployPath" placeholder="/etc/nginx/ssl" /></el-form-item>
        <el-form-item label="重载命令"><el-input v-model="deployForm.reloadCommand" placeholder="留空则使用默认 nginx reload" /></el-form-item>
      </el-form>
      <template #footer><el-button @click="deployDialogVisible = false">取消</el-button><el-button type="primary" @click="submitDeploy">执行部署</el-button></template>
    </el-dialog>
  </div>
</template>

<script>
import { getCmdbDatabaseTypeLabel } from '@/utils/cmdbPresentation.mjs'

const createHostRuleForm = () => ({ id: null, name: '', hostId: null, metricKey: 'cpu_usage', operator: '>=', thresholdValue: 90, severity: 'P2', notifyRobotIds: [], remark: '' })
const createDBRuleForm = () => ({ id: null, name: '', databaseId: null, metricKey: 'connectivity', operator: '>=', thresholdValue: 1, severity: 'P2', notifyRobotIds: [], remark: '' })
const createDomainForm = () => ({ id: null, domain: '', tags: '', remark: '' })
const createScheduleForm = () => ({ id: null, enabled: false, cronExpr: '', expireAlertDays: 30, scanTimeoutMs: 8000, notifyEnabled: false, notifyRobotId: null, autoDeployEnabled: false, deployHostId: null, deployPath: '', reloadCommand: '' })
const createDeployForm = () => ({ certId: null, hostId: null, deployPath: '/etc/nginx/ssl', reloadCommand: '' })

export default {
  name: 'MonitorAutomation',
  data() {
    return {
      activeTab: 'host',
      overview: {
        openEventCount: 0,
        hostRuleCount: 0,
        databaseHealthyCount: 0,
        databaseTotalCount: 0,
        expiringDomainCount: 0,
        enabledRobotCount: 0,
        failedNotifyLogCount: 0,
        recentEvents: [],
        recentActions: [],
        riskTips: [],
        recommendedActions: []
      },
      hostTemplates: [],
      hostRules: [],
      dbRules: [],
      dbSnapshots: [],
      domains: [],
      schedules: [],
      certs: [],
      deployLogs: [],
      events: [],
      hosts: [],
      databases: [],
      robots: [],
      operators: ['>', '>=', '<', '<=', '==', '!='],
      severities: ['P1', 'P2', 'P3', 'P4'],
      hostRuleDialogVisible: false,
      dbRuleDialogVisible: false,
      domainDialogVisible: false,
      scheduleDialogVisible: false,
      deployDialogVisible: false,
      hostRuleForm: createHostRuleForm(),
      dbRuleForm: createDBRuleForm(),
      domainForm: createDomainForm(),
      scheduleForm: createScheduleForm(),
      deployForm: createDeployForm()
    }
  },
  methods: {
    applyRouteQuery() {
      const tab = String(this.$route?.query?.tab || '').trim()
      if (tab && ['host', 'db', 'ssl', 'events'].includes(tab)) {
        this.activeTab = tab
      }
    },
    openWorkbenchPath(path) {
      if (!path) return
      this.$router.push(path)
    },
    async refreshAll() {
      try {
        this.applyRouteQuery()
        const [overview, templates, hostRules, dbRules, dbSnapshots, domains, schedules, certs, deployLogs, events, hosts, databases, robots] = await Promise.all([
          this.$api.getMonitorAutomationOverview(),
          this.$api.getMonitorHostAlertTemplates(),
          this.$api.getMonitorHostAlertRules(),
          this.$api.getMonitorDBAlertRules(),
          this.$api.getMonitorDBHealthSnapshots(),
          this.$api.getMonitorSSLDomains(),
          this.$api.getMonitorSSLSchedules(),
          this.$api.getMonitorSSLCerts(),
          this.$api.getMonitorSSLDeployLogs(),
          this.$api.getMonitorAutomationEvents({ pageNum: 1, pageSize: 20 }),
          this.$api.getCmdbHostList({ page: 1, pageSize: 500 }),
          this.$api.listDatabases({ page: 1, pageSize: 500 }),
          this.$api.getMonitorAlertNotifyRobotList()
        ])
        this.overview = overview.data?.data || this.overview
        this.hostTemplates = templates.data?.data || []
        this.hostRules = hostRules.data?.data || []
        this.dbRules = dbRules.data?.data || []
        this.dbSnapshots = dbSnapshots.data?.data || []
        this.domains = domains.data?.data || []
        this.schedules = schedules.data?.data || []
        this.certs = certs.data?.data || []
        this.deployLogs = deployLogs.data?.data || []
        this.events = events.data?.data?.list || []
        this.hosts = hosts.data?.data?.list || []
        this.databases = databases.data?.data?.list || []
        this.robots = robots.data?.data || []
      } catch (error) {
        this.$message.error(error?.response?.data?.message || error.message || '刷新监控深化页面失败')
      }
    },
    ensureApiSuccess(response, fallbackMessage) {
      const code = response?.data?.code
      if (code !== 200) {
        throw new Error(response?.data?.message || fallbackMessage || '请求失败')
      }
      return response.data?.data
    },
    validateHostRuleForm(payload) {
      if (!payload.name) throw new Error('请输入规则名称')
      if (!payload.hostId) throw new Error('请选择主机')
      if (!payload.metricKey) throw new Error('请选择指标')
      if (!payload.severity) throw new Error('请选择等级')
    },
    validateDBRuleForm(payload) {
      if (!payload.name) throw new Error('请输入规则名称')
      if (!payload.databaseId) throw new Error('请选择数据库')
      if (!payload.metricKey) throw new Error('请选择指标')
      if (!payload.severity) throw new Error('请选择等级')
    },
    formatTime(value) { return value ? new Date(value).toLocaleString('zh-CN') : '-' },
    severityTagType(value) { return ({ P1: 'danger', P2: 'warning', P3: 'primary', P4: 'info' })[value] || 'info' },
    databaseTypeText(type) { return getCmdbDatabaseTypeLabel(type) || 'Unknown' },
    certStatusText(status) { return ({ 1: '申请中', 2: '验证中', 3: '已签发', 4: '已过期', 5: '失败' })[Number(status)] || '未知' },
    certStatusTagType(status) { return ({ 1: 'info', 2: 'warning', 3: 'success', 4: 'danger', 5: 'danger' })[Number(status)] || 'info' },
    deployStatusText(status) { return ({ 1: '执行中', 2: '成功', 3: '失败' })[Number(status)] || '未知' },
    deployStatusTagType(status) { return ({ 1: 'warning', 2: 'success', 3: 'danger' })[Number(status)] || 'info' },
    robotNameById(id) { const item = this.robots.find(robot => robot.id === id); return item ? item.name : '' },
    applyHostTemplate(item) { this.hostRuleForm = { ...createHostRuleForm(), name: item.name, metricKey: item.metricKey, operator: item.operator, thresholdValue: item.thresholdValue, severity: item.severity }; this.hostRuleDialogVisible = true },
    openHostRuleDialog(row) { this.hostRuleForm = row ? { ...createHostRuleForm(), ...row } : createHostRuleForm(); this.hostRuleDialogVisible = true },
    openDBRuleDialog(row) { this.dbRuleForm = row ? { ...createDBRuleForm(), ...row } : createDBRuleForm(); this.dbRuleDialogVisible = true },
    openDomainDialog(row) { this.domainForm = row ? { ...createDomainForm(), ...row } : createDomainForm(); this.domainDialogVisible = true },
    openScheduleDialog(row) { this.scheduleForm = row ? { ...createScheduleForm(), ...row } : createScheduleForm(); this.scheduleDialogVisible = true },
    openDeployDialog() { this.deployForm = createDeployForm(); this.deployDialogVisible = true },
    async submitHostRule() {
      try {
        const payload = { ...this.hostRuleForm }
        this.validateHostRuleForm(payload)
        const response = payload.id
          ? await this.$api.updateMonitorHostAlertRule(payload.id, payload)
          : await this.$api.createMonitorHostAlertRule(payload)
        this.ensureApiSuccess(response, '保存主机规则失败')
        this.$message.success('主机规则已保存')
        this.hostRuleDialogVisible = false
        this.refreshAll()
      } catch (error) { this.$message.error(error?.response?.data?.message || error.message || '保存主机规则失败') }
    },
    async submitDBRule() {
      try {
        const payload = { ...this.dbRuleForm }
        this.validateDBRuleForm(payload)
        const response = payload.id
          ? await this.$api.updateMonitorDBAlertRule(payload.id, payload)
          : await this.$api.createMonitorDBAlertRule(payload)
        this.ensureApiSuccess(response, '保存数据库规则失败')
        this.$message.success('数据库规则已保存')
        this.dbRuleDialogVisible = false
        this.refreshAll()
      } catch (error) { this.$message.error(error?.response?.data?.message || error.message || '保存数据库规则失败') }
    },
    async submitDomain() {
      try {
        const payload = { ...this.domainForm }
        if (payload.id) await this.$api.updateMonitorSSLDomain(payload.id, payload)
        else await this.$api.createMonitorSSLDomain(payload)
        this.$message.success('域名巡检对象已保存')
        this.domainDialogVisible = false
        this.refreshAll()
      } catch (error) { this.$message.error(error?.response?.data?.message || error.message || '保存域名失败') }
    },
    async submitSchedule() {
      try {
        await this.$api.updateMonitorSSLSchedule(this.scheduleForm.id, this.scheduleForm)
        this.$message.success('巡检调度已保存')
        this.scheduleDialogVisible = false
        this.refreshAll()
      } catch (error) { this.$message.error(error?.response?.data?.message || error.message || '保存巡检调度失败') }
    },
    async submitDeploy() {
      try {
        await this.$api.deployMonitorSSLCert(this.deployForm)
        this.$message.success('证书部署任务已执行')
        this.deployDialogVisible = false
        this.refreshAll()
      } catch (error) { this.$message.error(error?.response?.data?.message || error.message || '执行证书部署失败') }
    },
    async toggleHostRule(row) { try { await this.$api.updateMonitorHostAlertRuleStatus(row.id, row.status === 1 ? 0 : 1); this.refreshAll() } catch (error) { this.$message.error(error?.response?.data?.message || error.message || '更新主机规则状态失败') } },
    async toggleDBRule(row) { try { await this.$api.updateMonitorDBAlertRuleStatus(row.id, row.status === 1 ? 0 : 1); this.refreshAll() } catch (error) { this.$message.error(error?.response?.data?.message || error.message || '更新数据库规则状态失败') } },
    async toggleDomain(row) { try { await this.$api.updateMonitorSSLDomainStatus(row.id, row.status === 1 ? 0 : 1); this.refreshAll() } catch (error) { this.$message.error(error?.response?.data?.message || error.message || '更新域名状态失败') } },
    async deleteHostRule(row) { try { await this.$confirm(`确认删除主机规则“${row.name}”吗？`, '提示', { type: 'warning' }); await this.$api.deleteMonitorHostAlertRule(row.id); this.refreshAll() } catch (error) { if (error !== 'cancel' && error !== 'close') this.$message.error(error?.response?.data?.message || error.message || '删除主机规则失败') } },
    async deleteDBRule(row) { try { await this.$confirm(`确认删除数据库规则“${row.name}”吗？`, '提示', { type: 'warning' }); await this.$api.deleteMonitorDBAlertRule(row.id); this.refreshAll() } catch (error) { if (error !== 'cancel' && error !== 'close') this.$message.error(error?.response?.data?.message || error.message || '删除数据库规则失败') } },
    async deleteDomain(row) { try { await this.$confirm(`确认删除域名“${row.domain}”吗？`, '提示', { type: 'warning' }); await this.$api.deleteMonitorSSLDomain(row.id); this.refreshAll() } catch (error) { if (error !== 'cancel' && error !== 'close') this.$message.error(error?.response?.data?.message || error.message || '删除域名失败') } },
    async scanHostAlerts() { try { await this.$api.scanMonitorHostAlerts(); this.$message.success('主机告警扫描完成'); this.refreshAll() } catch (error) { this.$message.error(error?.response?.data?.message || error.message || '主机告警扫描失败') } },
    async scanDBAlerts() { try { await this.$api.scanMonitorDBAlerts(); this.$message.success('数据库告警扫描完成'); this.refreshAll() } catch (error) { this.$message.error(error?.response?.data?.message || error.message || '数据库告警扫描失败') } },
    async scanSSLDomains() { try { await this.$api.scanMonitorSSLDomains(); this.$message.success('SSL 域名扫描完成'); this.refreshAll() } catch (error) { this.$message.error(error?.response?.data?.message || error.message || 'SSL 域名扫描失败') } }
  },
  watch: {
    '$route.query': {
      handler() {
        this.applyRouteQuery()
      },
      deep: true
    }
  },
  created() {
    this.applyRouteQuery()
    this.refreshAll()
  }
}
</script>

<style scoped>
.header-row,
.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
}

.header-actions,
.table-actions {
  display: flex;
  gap: 12px;
  align-items: center;
  flex-wrap: wrap;
}

.title {
  font-size: 24px;
  font-weight: 800;
  color: var(--text-primary);
}

.subtitle {
  color: var(--text-muted);
  font-size: 13px;
  margin-top: 6px;
}

.summary-grid,
.workbench-grid {
  display: grid;
  gap: 16px;
  margin-bottom: 22px;
}

.summary-grid {
  grid-template-columns: repeat(3, minmax(0, 1fr));
}

.summary-card,
.workbench-panel,
.template-card {
  border-radius: 20px;
  border: 1px solid var(--border-subtle);
  box-shadow: var(--shadow-card);
}

.summary-card {
  padding: 18px;
  color: #fff;
}

.summary-card.red { background: linear-gradient(135deg, rgba(251, 113, 133, 0.96), rgba(239, 68, 68, 0.92)); }
.summary-card.blue { background: linear-gradient(135deg, rgba(37, 99, 235, 0.96), rgba(20, 130, 255, 0.92)); }
.summary-card.teal { background: linear-gradient(135deg, rgba(13, 148, 136, 0.96), rgba(45, 212, 191, 0.88)); }
.summary-card.amber { background: linear-gradient(135deg, rgba(246, 173, 60, 0.96), rgba(245, 158, 11, 0.92)); }
.summary-card.slate { background: linear-gradient(135deg, rgba(71, 85, 105, 0.96), rgba(30, 41, 59, 0.92)); }
.summary-card.wine { background: linear-gradient(135deg, rgba(190, 24, 93, 0.96), rgba(157, 23, 77, 0.92)); }

.summary-label {
  font-size: 13px;
  opacity: 0.92;
}

.summary-value {
  margin-top: 12px;
  font-size: 32px;
  font-weight: 800;
}

.section-title,
.workbench-title,
.template-name {
  font-size: 18px;
  font-weight: 700;
  color: var(--text-primary);
}

.workbench-grid {
  grid-template-columns: repeat(2, minmax(0, 1fr));
}

.workbench-panel {
  padding: 18px;
  background: rgba(255, 255, 255, 0.03);
}

.workbench-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.workbench-item {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  align-items: center;
  padding: 12px 14px;
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.04);
}

.workbench-item.risk {
  align-items: flex-start;
  background: rgba(251, 113, 133, 0.12);
}

.workbench-item-title {
  font-size: 14px;
  font-weight: 700;
  color: var(--text-primary);
}

.workbench-item-summary {
  margin-top: 4px;
  font-size: 12px;
  line-height: 1.6;
  color: var(--text-muted);
}

.workbench-item-meta {
  display: flex;
  flex-direction: column;
  gap: 8px;
  align-items: flex-end;
  font-size: 12px;
  color: var(--text-muted);
}

.template-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(220px, 1fr));
  gap: 14px;
  margin: 18px 0;
}

.template-card {
  padding: 16px;
  cursor: pointer;
  background: rgba(255, 255, 255, 0.03);
}

.template-detail {
  margin-top: 8px;
  color: var(--text-muted);
  font-size: 13px;
  min-height: 38px;
}

.data-table {
  margin-top: 12px;
  border-radius: 16px;
  overflow: hidden;
}

@media (max-width: 1200px) {
  .summary-grid,
  .workbench-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 768px) {
  .summary-grid,
  .workbench-grid {
    grid-template-columns: 1fr;
  }

  .header-row,
  .section-header,
  .workbench-item {
    flex-direction: column;
    align-items: flex-start;
  }

  .workbench-item-meta {
    align-items: flex-start;
  }
}
</style>
