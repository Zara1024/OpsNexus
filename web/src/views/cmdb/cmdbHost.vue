<template>
  <div class="cmdb-host-management">
    <el-card shadow="hover" class="host-card">
    <!-- 左右布局容器 -->
    <div class="host-management-container">
      <!-- 宸︿晶鍒嗙粍鏍?-->
      <CmdbGroup
        ref="cmdbGroup"
        :group-list="groupList"
        :expanded-keys="expandedKeys"
        @group-search="handleGroupSearch"
        @group-click="handleGroupClick"
        @node-expand="handleNodeExpand"
        @node-collapse="handleNodeCollapse"
        @collapse-all="handleCollapseAll"
        @expand-all="handleExpandAll"
        @create-group="handleCreateGroup"
        @update-group="handleUpdateGroup"
        @delete-group="handleDeleteGroup"
      />
      <!-- 右侧主机管理区域 -->
      <div class="host-table-section">
        <!-- 搜索表单 -->
        <div class="search-section">
          <el-form :inline="true" :model="queryParams" class="demo-form-inline">
            <el-form-item label="主机名称" prop="hostName">
              <el-input
                  size="small"
                  placeholder="请输入主机名称"
                  clearable
                  v-model="queryParams.hostName"
                  @keyup.enter="handleQuery"
                  style="width: 160px;"
              />
            </el-form-item>
            <el-form-item label="IP地址" prop="ip">
              <el-input
                  size="small"
                  placeholder="请输入IP地址"
                  clearable  
                  v-model="queryParams.ip"
                  @keyup.enter="handleQuery"
                  style="width: 120px;"
              />
            </el-form-item>
            <el-form-item label="主机状态" prop="status">
              <el-select size="small" placeholder="请选择状态" v-model="queryParams.status" style="width: 120px;">
                <el-option v-for="item in statusList" :key="item.value" :label="item.label" :value="item.value"></el-option>
              </el-select>
            </el-form-item>
            <!-- 操作按钮 -->
            <div class="action-section">
            <el-row :gutter="10" class="mb8" style="text-align: left">
              <el-col :span="24">
                <!-- 搜索按钮 - 蓝色 -->
                <el-button type="primary" size="small" @click="handleQuery" style="margin-right: 10px">
                  <el-icon><Search /></el-icon>
                  <span style="margin-left: 4px">搜索</span>
                </el-button>
                
                <!-- 重置按钮 - 黄色 -->
                <el-button type="warning" size="small" @click="resetQuery" style="margin-right: 10px">
                  <el-icon><Refresh /></el-icon>
                  <span style="margin-left: 4px">重置</span>
                </el-button>
                
                <!--新建主机 - 绿色-->
                <el-dropdown
                  ref="createDropdown"
                  @command="handleCreateCommand"
                  @visible-change="handleDropdownVisibleChange"
                  :hide-on-click="true"
                  trigger="click"
                  placement="bottom-start">
                  <el-button
                    type="success"
                    size="small"
                    style="margin-right: 10px"
                    v-authority="['cmdb:ecs:add']"
                    @click.stop="handleCreateClick">
                    <el-icon><Plus /></el-icon>
                    <span style="margin-left: 4px">新建</span>
                    <el-icon class="el-icon--right"><ArrowDown /></el-icon>
                  </el-button>
                  <template #dropdown>
                    <el-dropdown-menu>
                      <el-dropdown-item command="importHost"><el-icon color="#409EFC" :size="20"><Edit /></el-icon>导入主机</el-dropdown-item>
                      <el-dropdown-item command="excelImport"><el-icon color="#409EFC" :size="20"><Folder /></el-icon>Excel导入</el-dropdown-item>
                      <el-dropdown-item command="cloudHost"><el-icon color="#409EFC" :size="21"><MostlyCloudy /></el-icon>云主机</el-dropdown-item>
                    </el-dropdown-menu>
                  </template>
                </el-dropdown>
                
                <!-- 批量删除按钮 -->
                <el-tooltip
                  :disabled="hasSelectedHosts"
                  content="请先勾选需要删除的主机"
                  placement="top"
                >
                  <span class="delete-host-btn-wrapper">
                    <el-button
                      :key="hasSelectedHosts ? 'delete-enabled' : 'delete-disabled'"
                      size="small"
                      v-authority="['cmdb:ecs:delete']"
                      @click="handleBatchHostDelete"
                      :loading="batchActionLoading === 'delete'"
                      :disabled="!hasSelectedHosts"
                      class="delete-host-btn"
                      style="margin-left: 10px"
                    >
                      <el-icon><Delete /></el-icon>
                      <span style="margin-left: 4px">删除主机</span>
                    </el-button>
                  </span>
                </el-tooltip>
              </el-col>
            </el-row>
            </div>
          </el-form>
        </div>

        <!-- 主机表格 -->
        <CmdbHostTable
          ref="hostTable"
          :key="$route.fullPath"
          :host-list="hostList"
          :loading="loading"
          @selection-change="handleHostSelectionChange"
          @show-detail="showHostDetail"
          @edit-host="showEditHostDialog"
          @open-terminal="handleHostSSH"
          @open-audit="handleHostAudit"
          @sync-host="handleHostSync"
          @show-upload="showUploadDialog"
          @execute-command="executeCommand"
          @delete-host="handleHostDelete"
        />

        <div v-if="selectedHosts.length" class="batch-toolbar">
          <div class="batch-toolbar__summary">
            已选择 <strong>{{ selectedHosts.length }}</strong> 台主机
          </div>
          <div class="batch-toolbar__actions">
            <el-button
              size="small"
              type="primary"
              :loading="batchActionLoading === 'sync'"
              @click="handleBatchSync"
            >
              批量同步
            </el-button>
            <el-button
              size="small"
              :loading="batchActionLoading === 'connectivity'"
              @click="handleBatchConnectivityTest"
            >
              批量测试连通性
            </el-button>
            <el-button
              size="small"
              type="success"
              :loading="batchActionLoading === 'deploy-agent'"
              @click="handleBatchDeployAgent"
            >
              部署 Agent
            </el-button>
            <el-button
              size="small"
              type="danger"
              plain
              :loading="batchActionLoading === 'uninstall-agent'"
              @click="handleBatchUninstallAgent"
            >
              卸载 Agent
            </el-button>
          </div>
        </div>

        <!-- 分页 -->
        <div class="pagination-section">
          <el-pagination
              @size-change="handleSizeChange"
              @current-change="handleCurrentChange"
              :current-page="queryParams.pageNum"
              :page-sizes="[10, 50, 100, 500]"
              :page-size="queryParams.pageSize"
              layout="total, sizes, prev, pager, next, jumper"
              :total="total"
          ></el-pagination>
        </div>
      </div>
    </div>

    <!-- 鏂板涓绘満瀵硅瘽妗?-->
    <CreateHost
      :visible="addDialogVisible"
      :group-list="groupList"
      :auth-list="authList"
      @close="addDialogVisible = false"
      @submit="addHost"
      @refresh-auth-list="getAuthList"
    />

    <!-- 缂栬緫涓绘満瀵硅瘽妗?-->
    <EditHost
      :visible="editDialogVisible"
      :host-info="hostInfo"
      :group-list="groupList"
      :auth-list="authList"
      @close="editDialogVisible = false"
      @submit="editHost"
    />

    <!-- 导入云主机对话框 -->
    <CreateCloud v-model="cloudDialogVisible" @success="handleCloudImportSuccess" />

    <!-- Excel瀵煎叆瀵硅瘽妗?-->
    <CreateExcel v-model="ExcelDialogVisible" @success="handleExcelImportSuccess" />

    <!-- SSH缁堢瀵硅瘽妗?-->
    <HostSSH 
      v-if="sshDialogVisible"
      :visible="sshDialogVisible"
      :host-id="currentHostId"
      @update:visible="val => {
        sshDialogVisible = val
      }"
    />

    <!-- 鏂囦欢涓婁紶瀵硅瘽妗?-->
    <el-dialog title="文件上传" v-model="uploadDialogVisible" width="25%">
      <el-form :model="uploadForm" :rules="uploadRules" ref="uploadFormRef" label-width="100px">
        <el-form-item label="目标主机">
          <el-input v-model="currentUploadHost.hostName" disabled />
        </el-form-item>
        <el-form-item label="目标路径" prop="targetPath">
          <el-input v-model="uploadForm.targetPath" placeholder="请输入目标路径" />
        </el-form-item>
        <el-form-item label="上传文件" prop="file">
          <el-upload
            class="upload-demo"
            :auto-upload="false"
            :on-change="handleFileChange"
            :on-remove="handleFileRemove"
            :show-file-list="false"
          >
            <el-button type="primary">选择文件</el-button>
            <template #tip>
              <div class="el-upload__tip" v-if="uploadForm.file">
                已选择：{{ uploadForm.file.name }}
                <el-button
                  type="danger"
                  text
                  icon="Close"
                  circle
                  size="small"
                  @click.stop="handleFileRemove"
                  style="margin-left: 8px"
                />
              </div>
              <div v-if="false" class="el-upload__tip" style="color: #999; margin-top: 15px">
                提示：请上传小于 {{ hostUploadMaxLabel }} 的文件
              </div>
              <div v-if="false" class="el-upload__tip" style="color: #999; margin-top: 15px">
                提示：支持直接上传文件，实际耗时取决于网络与目标主机状态
              </div>
            </template>
          </el-upload>
          <div class="upload-dialog-hint">
            支持直接上传文件，实际耗时取决于网络与目标主机状态
          </div>
        </el-form-item>
        <el-progress 
          v-if="isUploading"
          :percentage="uploadProgress" 
          :status="uploadProgress === 100 ? 'success' : ''"
        />
      </el-form>
      <template #footer>
        <el-button @click="uploadDialogVisible = false">取消</el-button>
        <el-button 
          type="primary" 
          @click="handleUpload"
          :loading="isUploading"
          :disabled="!uploadForm.file || !uploadForm.hostId"
        >
          寮€濮嬩笂浼?        </el-button>
      </template>
    </el-dialog>


    <!-- 主机详情抽屉 -->
    <el-drawer
      v-model="detailDrawer"
      title="主机详情"
      direction="rtl"
      size="40%"
      :before-close="handleDetailClose">

      <!-- 浠〃鐩橀儴鍒?-->
      <div class="dashboard-section">
        <div class="gauge-container">
          <div ref="cpuGauge" class="gauge-item"></div>
          <div ref="memoryGauge" class="gauge-item"></div>
          <div ref="diskGauge" class="gauge-item"></div>
        </div>
      </div>

      <!-- 基本信息部分 -->
      <h3 style="margin: 5px 0 10px 0">基本信息</h3>
      <el-descriptions :column="1" border>
        <el-descriptions-item label="主机名称">{{ hostDetail.hostName }}</el-descriptions-item>
        <el-descriptions-item label="主机分组">{{ getGroupName(hostDetail.groupId) }}</el-descriptions-item>
        <el-descriptions-item label="设备类型">{{ getDeviceTypeLabel(hostDetail.deviceType) }}</el-descriptions-item>
        <el-descriptions-item label="管理地址">
          {{ formatAccessAddress(hostDetail) }}
        </el-descriptions-item>
        <el-descriptions-item v-if="!isWindowsHost(hostDetail)" label="认证类型">
          {{ getAuthTypeName(hostDetail.sshKeyId) }}
        </el-descriptions-item>
        <el-descriptions-item v-else label="RDP连接">
          {{ formatRDPAccess(hostDetail) }}
        </el-descriptions-item>
        <el-descriptions-item label="描述信息">{{ hostDetail.remark }}</el-descriptions-item>
      </el-descriptions>

      <div class="detail-section">
        <div class="detail-section__head">
          <h3>连接中心</h3>
          <el-tag :type="isWindowsHost(hostDetail) ? 'primary' : 'success'">
            {{ isWindowsHost(hostDetail) ? 'Windows' : 'Linux' }}
          </el-tag>
        </div>
        <div class="connection-center">
          <div
            v-for="entry in hostConnectionEntries"
            :key="entry.key"
            :class="['connection-card', { 'connection-card--disabled': !entry.available }]"
          >
            <div class="connection-card__header">
              <div class="connection-card__title">{{ entry.title }}</div>
              <el-tag size="small" :type="entry.available ? 'success' : 'info'">
                {{ entry.available ? '可用' : '不可用' }}
              </el-tag>
            </div>
            <div class="connection-card__body">
              {{ entry.detail || entry.reason || '可直接从当前主机页进入。' }}
            </div>
            <div class="connection-card__footer">
              <el-button
                v-if="entry.key === 'rdp'"
                size="small"
                :disabled="!entry.available"
                @click="handleConnectionEntryAction(entry)"
              >
                复制连接信息
              </el-button>
              <el-button
                v-else
                size="small"
                type="primary"
                :disabled="!entry.available"
                @click="handleConnectionEntryAction(entry)"
              >
                立即进入
              </el-button>
            </div>
          </div>
        </div>
      </div>

      <div class="detail-section">
        <div class="detail-section__head">
          <h3>终端审计</h3>
          <div class="detail-section__actions">
            <el-button link type="primary" @click="openHostAudit()">查看审计</el-button>
            <el-button
              link
              type="danger"
              :disabled="!recentRiskAuditSessions.length"
              @click="handleOpenRiskAudit"
            >
              鏈€杩戦闄╁懡浠?            </el-button>
          </div>
        </div>
        <el-skeleton v-if="recentAuditLoading" :rows="3" animated />
        <div v-else-if="recentAuditSessions.length" class="audit-summary-list">
          <div v-for="item in recentAuditSessions.slice(0, 3)" :key="item.sessionId" class="audit-summary-item">
            <div class="audit-summary-item__meta">
              <div class="audit-summary-item__title">{{ item.hostIp || item.hostName || hostDetail.hostName }}</div>
              <div class="audit-summary-item__time">{{ item.startTime || '-' }}</div>
            </div>
            <div class="audit-summary-item__command">{{ item.latestCommand || '暂无命令摘要' }}</div>
            <div class="audit-summary-item__footer">
              <el-tag size="small" :type="getAuditRiskTagType(item.riskLevel)">{{ getAuditRiskLabel(item.riskLevel) }}</el-tag>
              <el-button link type="primary" @click="openHostAudit({ sessionId: item.sessionId })">定位会话</el-button>
            </div>
          </div>
        </div>
        <el-empty v-else description="当前主机暂无终端审计记录" />
      </div>

      <!-- 扩展信息部分 -->
      <div style="margin: 20px 0 10px 0; display: flex; justify-content: space-between; align-items: center;">
        <h3 style="margin: 0">扩展信息</h3>
        <el-button 
          type="primary" 
          size="mini" 
          icon="Refresh"
          :loading="syncLoading"
          :disabled="!canSyncHost(hostDetail)"
          v-authority="['cmdb:ecs:rsync']"
          @click="handleHostSync()"
        >
          {{ syncLoading ? '鍚屾涓?..' : '鍚屾' }}
        </el-button>
      </div>
      <el-descriptions :column="1" border>
        <el-descriptions-item label="实例ID">{{ hostDetail.instanceId }}</el-descriptions-item>
        <el-descriptions-item label="实例名称">{{ hostDetail.name }}</el-descriptions-item>
        <el-descriptions-item label="操作系统">{{ hostDetail.os }}</el-descriptions-item>
        <el-descriptions-item label="CPU">{{ formatHostResourceDisplay(hostDetail.cpu, '核') }}</el-descriptions-item>
        <el-descriptions-item label="内存">{{ formatHostResourceDisplay(hostDetail.memory, 'G') }}</el-descriptions-item>
        <el-descriptions-item label="磁盘">{{ formatHostResourceDisplay(hostDetail.disk, 'GB') }}</el-descriptions-item>
        <el-descriptions-item label="内网IP">{{ hostDetail.privateIp }}</el-descriptions-item>
        <el-descriptions-item label="公网IP">{{ hostDetail.publicIp || '无' }}</el-descriptions-item>
        <el-descriptions-item v-if="isWindowsHost(hostDetail)" label="RDP地址/域名">{{ hostDetail.remoteDomain || '未配置' }}</el-descriptions-item>
        <el-descriptions-item v-if="isWindowsHost(hostDetail)" label="RDP端口">{{ hostDetail.rdpPort || 3389 }}</el-descriptions-item>
        <el-descriptions-item v-if="isWindowsHost(hostDetail)" label="RDP用户名">{{ hostDetail.rdpUsername || '未配置' }}</el-descriptions-item>
        <el-descriptions-item v-if="isWindowsHost(hostDetail)" label="RDP域">{{ hostDetail.rdpDomain || '未配置' }}</el-descriptions-item>
        <el-descriptions-item label="实例计费方式">{{ hostDetail.billingType }}</el-descriptions-item>
        <el-descriptions-item label="网络计费方式">{{ hostDetail.networkBillingType || '按流量计费' }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ hostDetail.createTime }}</el-descriptions-item>
        <el-descriptions-item label="到期时间">{{ hostDetail.expireTime || '无' }}</el-descriptions-item>
        <el-descriptions-item label="更新时间">{{ hostDetail.updateTime }}</el-descriptions-item>
      </el-descriptions>
    </el-drawer>

    <!-- 鍛戒护鎵ц瀵硅瘽妗?-->
    <el-dialog
      v-if="commandDialog"
      title="执行命令"
      v-model="commandDialog.visible"
      width="72vw"
      top="8vh"
      class="command-dialog"
      :before-close="() => commandDialog.visible = false"
    >
      <el-form class="command-dialog__form" label-width="92px">
        <el-form-item label="主机名称">
          <el-input v-model="commandDialog.hostName" disabled />
        </el-form-item>
        <el-form-item label="执行命令">
          <el-input
            type="textarea"
            :rows="3"
            v-model="commandDialog.command"
            placeholder="请输入要执行的命令"
            clearable
          />
        </el-form-item>

        <el-alert
          v-if="commandDialog.riskPresentation"
          :title="`风险等级：${commandDialog.riskPresentation.text}`"
          :type="commandDialog.riskPresentation.level >= 2 ? 'error' : commandDialog.riskPresentation.tagType"
          :description="commandDialog.riskPresentation.reason || '当前命令未命中风险说明。'"
          :closable="false"
          style="margin-bottom: 16px"
        />
        
        <el-form-item>
          <el-button 
            type="primary" 
            @click="submitCommand"
            :loading="commandDialog.loading || commandDialog.previewLoading"
          >
            {{ commandDialog.previewLoading ? '椋庨櫓妫€鏌ヤ腑...' : '鎵ц' }}
          </el-button>
          <el-tag 
            v-if="commandDialog.status"
            :type="commandDialog.status === '执行成功' ? 'success' : 'danger'"
            style="margin-left: 20px"
          >
            {{ commandDialog.status }}
          </el-tag>
        </el-form-item>

      </el-form>

      <div v-if="commandDialog.output" class="command-output">
        <div class="command-output__header">执行结果</div>
        <pre class="command-output__panel">{{ commandDialog.output }}</pre>
      </div>

      <template #footer>
        <el-button @click="commandDialog.visible = false">关闭</el-button>
      </template>
    </el-dialog>
  </el-card>
  </div>
</template>

<script>

import * as echarts from 'echarts'
import configApi from '@/api/config'
import cmdbApi from '@/api/cmdb'
import CreateCloud from './Host/CreateCloud.vue'
import HostSSH from './Host/SSH.vue'
import CreateExcel from './Host/CreateExcel.vue'
import CmdbGroup from './Host/CmdbGroup.vue'
import CmdbHostTable from './Host/CmdbHostTableCompact.vue'
import CreateHost from './Host/CreateHost.vue'
import EditHost from './Host/EditHost.vue'
import { createHostAssetFormModel } from '@/utils/cmdbAssetPresentation.mjs'
import { extractCommandOutput } from '@/utils/cmdbPresentation.mjs'
import {
  formatHostResourceValue,
  resolveHostSyncTarget
} from '@/utils/cmdbHostDetailPresentation.mjs'
import {
  buildHostAuditRouteQuery,
  buildHostConnectionEntries,
  getCommandRiskPresentation,
  summarizeBatchConnectivity
} from '@/utils/cmdbHostPhase1.mjs'

export default {
  components: {
    CreateCloud,
    HostSSH,
    CreateExcel,
    CmdbGroup,
    CmdbHostTable,
    CreateHost,
    EditHost
  },
  data() {
    return {
      ExcelDialogVisible: false,
      commandDialog: null,
      expandedKeys: [],
      statusList: [
        { value: 2, label: '未认证' },
        { value: 1, label: '认证成功' },
        { value: 3, label: '认证失败' }
      ],
      loading: false,
      queryParams: {
        pageNum: 1,
        pageSize: 10,
        hostName: '',
        ip: '',
        status: '',
        groupId: ''
      },
      hostList: [],
      total: 0,
      selectedHosts: [],
      batchActionLoading: '',
      recentAuditLoading: false,
      recentAuditSessions: [],
      addDialogVisible: false,
      editDialogVisible: false,
      cloudDialogVisible: false,
      groupList: [],
      defaultProps: {
        children: 'children',
        label: 'name'
      },
      currentGroupId: null,
      authList: [],
      addForm: createHostAssetFormModel(),
      addFormRules: {
        hostName: [{ required: true, message: '请输入主机名称', trigger: 'blur' }],
        ip: [{ required: true, message: '请输入IP地址', trigger: 'blur' }],
        port: [{ required: true, message: '请输入端口号', trigger: 'blur' }],
        username: [{ required: true, message: '请输入连接用户名', trigger: 'blur' }],
        authId: [{ required: true, message: '请选择认证凭据', trigger: 'change' }],
        groupId: [{ required: true, message: '请选择主机分组', trigger: 'change' }]
      },
      hostInfo: {},
      editFormRules: {
        hostName: [{ required: true, message: '请输入主机名称', trigger: 'blur' }],
        ip: [{ required: true, message: '请输入IP地址', trigger: 'blur' }],
        port: [{ required: true, message: '请输入端口号', trigger: 'blur' }],
        username: [{ required: true, message: '请输入连接用户名', trigger: 'blur' }],
        authId: [{ required: true, message: '请选择认证凭据', trigger: 'change' }],
        groupId: [{ required: true, message: '请选择主机分组', trigger: 'change' }]
      },
      // SSH 终端对话框相关
      sshDialogVisible: false,
      currentHostId: null,
      // 上传对话框相关
      uploadDialogVisible: false,
      uploadForm: {
        hostId: null,
        file: null,
        targetPath: '/tmp'
      },
      currentUploadHost: null,
      uploadRules: {
      file: [{ required: true, message: '请选择上传文件', trigger: 'change' }],
      targetPath: [{ 
        required: true, 
        message: '请输入目标路径',
        trigger: ['blur', 'change'],
        validator: (rule, value, callback) => {
          if (value === '/tmp' || (value && value.trim() !== '')) {
            callback()
          } else {
            callback(new Error('请输入目标路径'))
          }
        }
      }]
      },
      isUploading: false,
      uploadProgress: 0,
    
    // 主机详情相关
    detailDrawer: false,
    syncLoading: false,
    hostDetail: {
        hostName: '',
        groupId: '',
        privateIp: '',
        publicIp: '',
        sshIp: '',
        sshName: '',
        sshKeyId: '',
        sshPort: 22,
        remark: '',
        vendor: '',
        region: '',
        instanceId: '',
        os: '',
        status: 0,
        cpu: '',
        memory: '',
        disk: '',
        deviceType: 'linux',
        remoteDomain: '',
        rdpPort: 3389,
        rdpUsername: '',
        rdpPassword: '',
        rdpDomain: '',
        supportsSsh: true,
        billingType: '',
        createTime: '',
        expireTime: '',
        updateTime: '',
        name: '',
        cpuUsage: 0,
        memoryUsage: 0,
        diskUsage: 0
      },
    // ECharts 实例
    cpuChart: null,
    memoryChart: null,
      diskChart: null
    }
  },
  computed: {
    hasSelectedHosts() {
      return Array.isArray(this.selectedHosts) && this.selectedHosts.length > 0
    },
    hostConnectionEntries() {
      return buildHostConnectionEntries(this.hostDetail || {})
    },
    recentRiskAuditSessions() {
      return (this.recentAuditSessions || []).filter(item => Number(item?.riskLevel || 0) > 0)
    }
  },
  created() {
    this.getAllGroups()
    this.getAuthList()
    this.getHostList()
  },

    beforeRouteEnter(to, from, next) {
      next(vm => {
        // 立即获取主机列表
        vm.getHostList().then(() => {
          // 涓绘満鍒楄〃鍔犺浇瀹屾垚鍚庣珛鍗宠Е鍙戠洃鎺ф暟鎹姞杞?          vm.$refs.hostTable?.fetchMonitorData()
        })
      })
    },

    beforeRouteUpdate(to, from, next) {
      // 立即获取主机列表
      this.getHostList().then(() => {
          this.$refs.hostTable?.fetchMonitorData()
        next()
      })
  },
  methods: {
    formatHostResourceDisplay(value, unit) {
      return formatHostResourceValue(value, unit)
    },
    async getAllGroups() {
      const { data: res } = await this.$api.getAllCmdbGroups()
      if (res.code === 200) {
        this.groupList = res.data
        // 设置默认分组为业务组
        const businessGroup = this.groupList.find(group => group.name === '默认业务组')
        if (businessGroup) {
          this.addForm.groupId = businessGroup.id
        }
      }
    },
    handleHostSelectionChange(selection) {
      this.selectedHosts = selection || []
    },
    getTerminalTargetHost() {
      if (this.selectedHosts.length > 1) {
        this.$message.warning('终端连接一次只支持一台主机，请只保留一个目标')
        return null
      }
      if (this.selectedHosts.length === 1) {
        return this.selectedHosts[0]
      }
      return this.hostList[0] || null
    },
    formatBatchHostNames(hosts = [], limit = 6) {
      const names = hosts.map(host => host.hostName || `主机-${host.id}`)
      if (names.length <= limit) {
        return names.join('、')
      }
      return `${names.slice(0, limit).join('、')} 等 ${names.length} 台主机`
    },
    async confirmBatchAction(title, message, type = 'warning') {
      const confirmResult = await this.$confirm(message, title, {
        confirmButtonText: '确认',
        cancelButtonText: '取消',
        type
      }).catch(err => err)
      return confirmResult === 'confirm'
    },
    async showBatchResultDialog(title, lines) {
      await this.$alert(lines.join('<br/>'), title, {
        confirmButtonText: '确定',
        dangerouslyUseHTMLString: true
      })
    },
    getBatchRunnableHosts(predicate, emptyMessage) {
      const hosts = (this.selectedHosts || []).filter(predicate)
      if (!hosts.length) {
        this.$message.warning(emptyMessage)
        return []
      }
      return hosts
    },
    isHostSyncAccepted(responseData) {
      if (responseData?.code !== 200) {
        return false
      }
      return Number(responseData?.data?.status ?? 1) !== 3
    },
    async handleBatchHostDelete() {
      if (!this.selectedHosts.length) {
        this.$message.warning('请先勾选需要删除的主机')
        return
      }
      await this.deleteHosts(this.selectedHosts, { batch: true })
    },
    closeDeletedHostContexts(hostIds = []) {
      const deletedIdSet = new Set(
        (hostIds || [])
          .map(id => Number(id))
          .filter(id => Number.isFinite(id))
      )
      if (!deletedIdSet.size) {
        return
      }
      if (this.detailDrawer && deletedIdSet.has(Number(this.hostDetail?.id))) {
        this.handleDetailClose()
      }
      if (this.currentUploadHost?.id && deletedIdSet.has(Number(this.currentUploadHost.id))) {
        this.uploadDialogVisible = false
        this.currentUploadHost = null
        this.uploadForm = {
          hostId: null,
          file: null,
          targetPath: '/tmp'
        }
      }
      if (this.currentHostId && deletedIdSet.has(Number(this.currentHostId))) {
        this.sshDialogVisible = false
        this.currentHostId = null
        if (this.commandDialog) {
          this.commandDialog.visible = false
        }
      }
    },
    async deleteHosts(hosts, { batch = false } = {}) {
      const targetHosts = (hosts || []).filter(host => host?.id)
      if (!targetHosts.length) {
        this.$message.warning(batch ? '请先勾选需要删除的主机' : '未获取到可删除的主机')
        return
      }
      const confirmed = await this.confirmBatchAction(
        batch ? '确认批量删除主机' : '确认删除主机',
        batch
          ? `将删除 ${targetHosts.length} 台主机：${this.formatBatchHostNames(targetHosts)}`
          : `是否确认删除主机"${targetHosts[0].hostName || `主机-${targetHosts[0].id}`}"?`
      )
      if (!confirmed) {
        this.$message.info(batch ? '已取消批量删除' : '已取消删除')
        return
      }
      if (batch) {
        this.batchActionLoading = 'delete'
      }
      try {
        const results = await Promise.allSettled(
          targetHosts.map(host => this.$api.deleteCmdbHost(host.id))
        )
        const successHosts = []
        const failedHosts = []
        results.forEach((item, index) => {
          const host = targetHosts[index]
          const hostName = host.hostName || `主机-${host.id}`
          if (item.status === 'fulfilled' && item.value?.data?.code === 200) {
            successHosts.push({ id: host.id, hostName })
          } else {
            const reason = item.status === 'fulfilled'
              ? (item.value?.data?.message || '返回结果异常')
              : (item.reason?.message || '请求失败')
            failedHosts.push(`${hostName}：${reason}`)
          }
        })
        if (successHosts.length) {
          this.closeDeletedHostContexts(successHosts.map(host => host.id))
          await this.getHostList()
        }
        if (!batch && !failedHosts.length) {
          this.$message.success('删除成功')
          return
        }
        const lines = [
          `共处理 ${targetHosts.length} 台主机。`,
          `删除成功：${successHosts.length} 台`,
          `删除失败：${failedHosts.length} 台`
        ]
        if (successHosts.length) {
          lines.push(`已删除主机：${successHosts.slice(0, 8).map(host => host.hostName).join('、')}`)
        }
        if (failedHosts.length) {
          lines.push(`失败详情：${failedHosts.slice(0, 8).join('<br/>')}`)
        }
        await this.showBatchResultDialog(batch ? '批量删除主机结果' : '删除主机结果', lines)
      } catch (error) {
        this.$message.error(`${batch ? '批量删除主机失败' : '删除主机失败'}: ${error.message || '未知错误'}`)
      } finally {
        if (batch) {
          this.batchActionLoading = ''
        }
      }
    },
    async handleBatchSync() {
      const runnableHosts = this.getBatchRunnableHosts(host => this.canSyncHost(host), '请选择支持 SSH 同步的 Linux 主机')
      if (!runnableHosts.length) {
        return
      }
      const skippedHosts = (this.selectedHosts || []).filter(host => !this.canSyncHost(host))
      const confirmed = await this.confirmBatchAction(
        '确认批量同步',
        `将为 ${runnableHosts.length} 台主机触发配置同步：${this.formatBatchHostNames(runnableHosts)}`
      )
      if (!confirmed) {
        return
      }

      this.batchActionLoading = 'sync'
      try {
        const results = await Promise.allSettled(
          runnableHosts.map(host => this.$api.syncHostConfig(host.id))
        )
        const successHosts = []
        const failedHosts = []

        results.forEach((item, index) => {
          const host = runnableHosts[index]
          if (item.status === 'fulfilled' && this.isHostSyncAccepted(item.value?.data)) {
            successHosts.push(host.hostName)
          } else {
            const reason = item.status === 'fulfilled'
              ? (item.value?.data?.data?.message || item.value?.data?.message || '同步失败')
              : (item.reason?.message || '未知错误')
            failedHosts.push(`${host.hostName}: ${reason}`)
          }
        })

        const lines = [
          `已选择 ${this.selectedHosts.length} 台主机，本次同步 ${runnableHosts.length} 台可执行主机`,
          `同步成功：${successHosts.length} 台`,
          `同步失败：${failedHosts.length} 台`
        ]
        if (skippedHosts.length) {
          lines.push(`已跳过 ${skippedHosts.length} 台 Windows 或未配置 SSH 的主机`)
        }
        if (failedHosts.length) {
          lines.push(`失败明细：${failedHosts.slice(0, 8).join('<br/>')}`)
        }

        await this.showBatchResultDialog('批量同步结果', lines)
        this.getHostList()
      } catch (error) {
        this.$message.error(`批量同步失败: ${error.message || '未知错误'}`)
      } finally {
        this.batchActionLoading = ''
      }
    },
    async handleBatchConnectivityTest() {
      if (!this.selectedHosts.length) {
        this.$message.warning('请先选择需要测试连通性的主机')
        return
      }
      const confirmed = await this.confirmBatchAction(
        '确认批量测试',
        `将测试 ${this.selectedHosts.length} 台主机的连通性`
      )
      if (!confirmed) {
        return
      }

      this.batchActionLoading = 'connectivity'
      try {
        const { data: res } = await this.$api.testHostConnectivity(this.selectedHosts.map(host => host.id))
        if (res.code !== 200) {
          throw new Error(res.message || '连通性测试失败')
        }
        const summary = summarizeBatchConnectivity(this.selectedHosts, res.data?.items || [])
        const failedItems = summary.items.filter(item => item.status !== 'connected')
        const lines = [
          `测试总数：${summary.total} 台`,
          `连通成功：${summary.connected} 台`,
          `连通失败：${summary.disconnected} 台`
        ]
        if (failedItems.length) {
          lines.push(`失败明细：${failedItems.slice(0, 8).map(item => `${item.hostName}: ${item.reason}`).join('<br/>')}`)
        }
        await this.showBatchResultDialog('批量连通性测试结果', lines)
      } catch (error) {
        this.$message.error(`批量连通性测试失败: ${error.message || '未知错误'}`)
      } finally {
        this.batchActionLoading = ''
      }
    },
    async handleBatchDeployAgent() {
      const runnableHosts = this.getBatchRunnableHosts(
        host => !this.isWindowsHost(host) && this.canUseSSH(host),
        '请选择支持 SSH 的 Linux 主机进行 Agent 部署'
      )
      if (!runnableHosts.length) {
        return
      }
      const confirmed = await this.confirmBatchAction(
        '确认批量部署 Agent',
        `将为 ${runnableHosts.length} 台主机部署 Agent：${this.formatBatchHostNames(runnableHosts)}`,
        'info'
      )
      if (!confirmed) {
        return
      }

      this.batchActionLoading = 'deploy-agent'
      try {
        const { data: res } = await this.$api.deployAgent(runnableHosts.map(host => host.id))
        if (res.code !== 200) {
          throw new Error(res.message || '部署 Agent 失败')
        }
        await this.showBatchResultDialog('批量部署 Agent', [
          `已为 ${runnableHosts.length} 台主机发起 Agent 部署任务`,
          '请稍后刷新列表查看 Agent 安装状态'
        ])
        this.getHostList()
      } catch (error) {
        this.$message.error(`批量部署 Agent 失败: ${error.message || '未知错误'}`)
      } finally {
        this.batchActionLoading = ''
      }
    },
    async handleBatchUninstallAgent() {
      const runnableHosts = this.getBatchRunnableHosts(
        host => !this.isWindowsHost(host),
        '请选择 Linux 主机进行 Agent 卸载'
      )
      if (!runnableHosts.length) {
        return
      }
      const confirmed = await this.confirmBatchAction(
        '确认批量卸载 Agent',
        `将为 ${runnableHosts.length} 台主机卸载 Agent：${this.formatBatchHostNames(runnableHosts)}`
      )
      if (!confirmed) {
        return
      }

      this.batchActionLoading = 'uninstall-agent'
      try {
        const { data: res } = await this.$api.uninstallAgent(runnableHosts.map(host => host.id))
        if (res.code !== 200) {
          throw new Error(res.message || '卸载 Agent 失败')
        }
        await this.showBatchResultDialog('批量卸载 Agent', [
          `已为 ${runnableHosts.length} 台主机发起 Agent 卸载任务`,
          '请稍后刷新列表查看 Agent 状态'
        ])
        this.getHostList()
      } catch (error) {
        this.$message.error(`批量卸载 Agent 失败: ${error.message || '未知错误'}`)
      } finally {
        this.batchActionLoading = ''
      }
    },
    getAuditRiskLabel(level) {
      return { 0: '低风险', 1: '中风险', 2: '高风险' }[Number(level)] || '未知风险'
    },
    getAuditRiskTagType(level) {
      return { 0: 'success', 1: 'warning', 2: 'danger' }[Number(level)] || 'info'
    },
    async fetchRecentHostAudit(host) {
      if (!host?.id) {
        this.recentAuditSessions = []
        return
      }
      this.recentAuditLoading = true
      try {
        const { data: res } = await this.$api.queryTerminalAuditSessionList({
          hostId: host.id,
          hostKeyword: host.hostName || host.sshIp || '',
          pageNum: 1,
          pageSize: 5
        })
        if (res.code === 200) {
          this.recentAuditSessions = res.data?.list || []
          return
        }
        this.recentAuditSessions = []
      } catch (error) {
        console.error('获取主机终端审计摘要失败:', error)
        this.recentAuditSessions = []
      } finally {
        this.recentAuditLoading = false
      }
    },
    openHostAudit(extraQuery = {}, host = null) {
      const targetHost = host || this.hostDetail
      const query = buildHostAuditRouteQuery(targetHost, extraQuery)
      this.$router.push({
        path: '/monitor/recording',
        query
      })
    },
    handleHostAudit(row) {
      if (!row?.id) {
        this.$message.warning('请选择可审计的主机')
        return
      }
      this.openHostAudit({}, row)
    },
    handleOpenRiskAudit() {
      const riskLevels = this.recentRiskAuditSessions.map(item => Number(item?.riskLevel || 0))
      const highestRiskLevel = riskLevels.length ? Math.max(...riskLevels) : 1
      this.openHostAudit({ riskLevel: highestRiskLevel || 1 })
    },
    handleConnectionEntryAction(entry) {
      if (!entry?.available) {
        return
      }
      if (entry.key === 'rdp') {
        if (!navigator.clipboard?.writeText) {
          this.$message.warning('当前浏览器不支持自动复制，请手动复制连接信息')
          return
        }
        navigator.clipboard.writeText(entry.detail || this.formatRDPAccess(this.hostDetail))
          .then(() => {
            this.$message.success('RDP 连接信息已复制')
          })
          .catch(() => {
            this.$message.warning('复制失败，请手动复制连接信息')
          })
        return
      }
      if (entry.key === 'ssh') {
        this.$router.push({
          path: '/cmdb/ssh',
          query: {
            hostId: this.hostDetail.id
          }
        })
        return
      }
      if (entry.key === 'command') {
        this.executeCommand(this.hostDetail)
        return
      }
      if (entry.key === 'upload') {
        this.showUploadDialog(this.hostDetail)
        return
      }
      if (entry.key === 'sync') {
        this.handleHostSync()
      }
    },

    // 处理分组搜索
    async handleGroupSearch(searchText) {
      this.groupSearchText = searchText
      if (!this.groupSearchText) {
        this.expandedKeys = []
        return
      }
      
      try {
        const { data: res } = await this.$api.getCmdbGroupByName(this.groupSearchText)
        
        if (res.code === 200 && res.data) {
          
          // 获取CmdbGroup组件的树引用
          const cmdbGroupRef = this.$refs.cmdbGroup
          const tree = cmdbGroupRef ? cmdbGroupRef.$refs.groupTree : null
          if (!tree) {
            console.error('树组件引用不存在')
            return
          }



          const findAndExpandParent = (groups, targetId, path = []) => {
            for (const group of groups) {

              if (group.id === targetId) {

                return [...path, group.id]
              }
              if (group.children && group.children.length > 0) {
                const foundPath = findAndExpandParent(group.children, targetId, [...path, group.id])
                if (foundPath) {
                  return foundPath
                }
              }
            }
            return null
          }
          
          // 获取展开路径
          const expandPath = findAndExpandParent(this.groupList, res.data.id)
          
          if (expandPath) {
            // 设置展开的keys
            this.expandedKeys = expandPath.slice(0, -1)
            
            this.$nextTick(() => {
              tree.setCurrentKey(res.data.id)
              
              // 确保树组件已更新
              setTimeout(() => {

              }, 500)
            })
          } else {
            console.warn('未找到匹配的分组')
            this.$message.warning('未找到匹配的分组')
          }
        } else {
          console.warn('分组搜索接口未返回结果')
        }
      } catch (error) {
        console.error('搜索分组失败:', error)
        this.$message.error('搜索分组失败: ' + (error?.message || '未知错误'))
      }
    },

    // 获取认证凭据列表
    async getAuthList() {
      try {
        const response = await configApi.getEcsAuthList({
          page: 1,
          pageSize: 100  // 鑾峰彇璁よ瘉鍑嵁锛岀敤浜庝笅鎷夐€夋嫨
        })
        
        if (response && response.data) {
          const res = response.data
          
          if (res.code === 200) {
            this.authList = Array.isArray(res.data?.list) ? res.data.list : []
          } else {
            console.error('获取认证凭据失败:', res.message || '未知错误')
            this.$message.error(`获取认证凭据失败: ${res.message || '未知错误'}`)
          }
        } else {
          console.error('获取认证凭据返回了无效响应:', response)
          this.$message.error('获取认证凭据失败: 响应格式无效')
        }
      } catch (error) {
        console.error('获取认证凭据异常:', error)
        this.$message.error(`获取认证凭据异常: ${error?.message || '未知错误'}`)
        this.authList = []
      }
    },
    
    // 获取主机列表
    async getHostList() {
      this.loading = true
      try {
        let response
        const { hostName, ip, status, pageNum, pageSize } = this.queryParams
        
        // 构建分页参数
        const baseParams = {
          page: pageNum,
          pageSize: pageSize,
          _t: Date.now()
        }

        
        const hasHostName = Boolean(hostName)
        const hasIp = Boolean(ip)
        const hasStatus = status !== '' && status !== null && status !== undefined
        const activeFilterCount = [hasHostName, hasIp, hasStatus].filter(Boolean).length

        // 根据查询条件选择API调用
        if (activeFilterCount > 1) {
          response = await this.$api.getCmdbHostList({
            ...baseParams,
            name: hostName || undefined,
            ip: ip || undefined,
            status: hasStatus ? status : undefined
          })
        } else if (hostName && !ip && !hasStatus) {
          response = await this.$api.GetCmdbHostsByHostNameLike(hostName, baseParams)
        } else if (ip && !hostName && !hasStatus) {
          response = await this.$api.GetCmdbHostsByIP(ip, baseParams)
        } else if (hasStatus && !hostName && !ip) {
          response = await this.$api.GetCmdbHostsByStatus(status, baseParams)
        } else {
          response = await this.$api.getCmdbHostList(baseParams)
        }
        
        
        // 处理axios响应结构
        const axiosResponse = response?.data ? response : { data: response }
        
        if (!axiosResponse || typeof axiosResponse !== 'object') {
          throw new Error('API返回无效响应格式')
        }

        const res = axiosResponse.data
        if (!res || typeof res !== 'object') {
          throw new Error('响应数据为空')
        }

        // 妫€鏌ュ搷搴旂爜
        if (res.code === undefined || res.code !== 200) {
          throw new Error(res.message || '获取主机列表失败')
        }

        if (res.data === undefined) {
          throw new Error('响应缺少 data 字段')
        }

        // 处理响应数据 - 适配不同API返回格式
        if (Array.isArray(res.data)) {
          this.hostList = res.data
          this.total = res.data.length
        } else if (res.data?.list) {
          this.hostList = res.data.list
          this.total = res.data.total
          if (res.data.page) {
            this.queryParams.pageNum = res.data.page
          }
          if (res.data.pageSize) {
            this.queryParams.pageSize = res.data.pageSize
          }
        } else {
          // 其他情况
          this.hostList = []
          this.total = 0
        }

        this.selectedHosts = []

        this.$nextTick(() => {
          this.$refs.hostTable?.fetchMonitorData()
        })
        
      } catch (error) {
        console.error('获取主机列表失败:', {
          error: error?.message || '未知错误',
          stack: error?.stack || '无堆栈信息',
          queryParams: this.queryParams
        })
        this.$message.error(`获取主机列表失败: ${error?.message || '未知错误'}`)
        this.hostList = []
        this.total = 0
        this.selectedHosts = []
      } finally {
        this.loading = false
      }
    },
    
    // 处理分组选择变化
    handleGroupChange(value) {
      if (value && value.length > 0) {
        // 鍙栨渶鍚庝竴绾т綔涓洪€変腑鍒嗙粍ID
        this.addForm.groupId = value[value.length - 1]
        this.hostInfo.groupId = value[value.length - 1]
      } else {
        const defaultGroup = this.groupList.find(item => item.isDefault)
        if (defaultGroup) {
          this.addForm.groupId = defaultGroup.id
          this.hostInfo.groupId = defaultGroup.id
        }
      }
    },

    // 根据分组获取主机
    async getHostsByGroup(groupId) {
      this.loading = true
      this.queryParams.groupId = groupId
      try {
        const { data: res } = await this.$api.getCmdbHostsByGroupId(groupId, {
          page: this.queryParams.pageNum,
          pageSize: this.queryParams.pageSize
        })
        if (res.code === 200) {
          this.hostList = res.data || []
          this.total = res.data?.length || 0
          this.selectedHosts = []
        }
      } catch (error) {
        console.error('获取主机列表失败:', error)
        this.hostList = []
        this.total = 0
        this.selectedHosts = []
      } finally {
        this.loading = false
      }
    },
    
    // 点击分组节点
    handleGroupClick(node, element) {
      let groupId
      if (element && element.data && element.data.id) {
        groupId = element.data.id
      } else if (element && element.id) {
        groupId = element.id
      } else if (node && node.key) {
        groupId = node.key
      }
      
      if (!groupId) {
        this.$message.warning("无法获取分组ID")
        return
      }
      
      this.currentGroupId = groupId
      this.getHostsByGroup(groupId)
    },

    handleNodeExpand(data, node) {
      if (!this.expandedKeys.includes(node.key)) {
        this.expandedKeys.push(node.key)
      }
    },

    handleNodeCollapse(data, node) {
      this.expandedKeys = this.expandedKeys.filter(key => key !== node.key)
    },

    handleCollapseAll() {
      this.expandedKeys = []
    },

    handleExpandAll() {
      const allKeys = []
      const collectKeys = (nodes) => {
        nodes.forEach(node => {
          allKeys.push(node.id)
          if (node.children && node.children.length > 0) {
            collectKeys(node.children)
          }
        })
      }
      collectKeys(this.groupList)
      this.expandedKeys = allKeys
    },
    
    // 搜索按钮操作
    handleQuery() {
      this.queryParams.pageNum = 1
      this.getHostList()
    },
    
    // 重置按钮操作
    resetQuery() {
      this.queryParams = {
        pageNum: 1,
        pageSize: 10,
        hostName: '',
        ip: '',
        status: '',
        groupId: ''
      }
      this.currentGroupId = null
      this.getHostList()
    },
    
    // pageSize变化
    handleSizeChange(newSize) {
      this.queryParams.pageSize = newSize
      this.getHostList()
    },
    
    // pageNum变化
    handleCurrentChange(newPage) {
      this.queryParams.pageNum = newPage
      this.getHostList()
    },
    
    // 新增主机
    async addHost(requestData) {
      try {
        const { data: res } = await this.$api.createCmdbHost(requestData)

        if (res.code === 200) {
          this.$message.success(res.data?.msg || '新增主机成功')
          this.addDialogVisible = false

          await this.getHostList()

          if (requestData.deviceType !== 'windows' && requestData.sshKeyId) {
            // 绛夊緟3绉掕鍚庣鍚屾涓绘満鐘舵€佷俊鎭紝鐒跺悗鍐嶆鍒锋柊
            setTimeout(async () => {
              await this.getHostList()
              this.$message.success('主机信息同步完成')
            }, 3000)
          }
        } else if (res.code === 426) {
          this.$message.error(`璁よ瘉鍑嵁涓嶅瓨鍦?鍑嵁ID: ${requestData.sshKeyId})锛岃妫€鏌ュ悗閲嶈瘯`)
          // 刷新凭据列表
          await this.getAuthList()
        } else {
          this.$message.error(res.message || '新增主机失败')
        }
      } catch (error) {
        console.error('新增主机失败:', error)
        this.$message.error('新增主机失败: ' + error.message)
      }
    },
    
    async showEditHostDialog(id) {
      const { data: res } = await this.$api.getCmdbHostById(id)
      if (res.code === 200) {
        this.hostInfo = {
          id: res.data.id,
          hostName: res.data.hostName,
          groupId: res.data.groupId,
          remark: res.data.remark,
          ip: res.data.sshIp,
          port: res.data.sshPort,
          username: res.data.sshName,
          authId: res.data.sshKeyId,
          vendor: Number(res.data.vendor || 1),
          deviceType: res.data.deviceType || 'linux',
          remoteDomain: res.data.remoteDomain || '',
          rdpPort: Number(res.data.rdpPort || 3389),
          rdpUsername: res.data.rdpUsername || '',
          rdpPassword: res.data.rdpPassword || '',
          rdpDomain: res.data.rdpDomain || ''
        }
        this.editDialogVisible = true
      }
    },
    
    editDialogClosed() {
      this.$refs.editFormRef.resetFields()
    },
    
    // 编辑主机信息
    async editHost(requestData) {
      try {
        // 验证凭据是否存在
        const authExists = requestData.deviceType !== 'windows'
          ? this.authList.some(auth => auth.id === requestData.sshKeyId)
          : true
        if (!authExists) {
          this.$message.error('当前认证凭据不存在，请刷新凭据列表后重试')
          return false
        }

        const port = Number(requestData.sshPort)
        if (requestData.deviceType !== 'windows' && (isNaN(port) || port < 1 || port > 65535)) {
          this.$message.error('端口号必须在 1-65535 之间')
          return false
        }
        const { data: res } = await this.$api.updateCmdbHost(requestData)
        if (res.code === 200) {
          this.$message.success('修改主机成功')
          this.editDialogVisible = false
          this.getHostList()
          return true
        } else if (res.code === 426) {
          this.$message.error(`璁よ瘉鍑嵁涓嶅瓨鍦?鍑嵁ID: ${requestData.sshKeyId})锛岃妫€鏌ュ悗閲嶈瘯`)
          // 刷新凭据列表
          await this.getAuthList()
          return false
        } else {
          this.$message.error(res.message || '修改主机失败')
          return false
        }
      } catch (error) {
        console.error('修改主机失败:', error)
        this.$message.error('修改主机失败: ' + error.message)
        return false
      }
    },
    
    // 鑾峰彇鐘舵€佹枃鏈?
    // 根据分组ID获取分组名称
    getGroupName(groupId) {
      if (!groupId) return '未分组'
      const findGroup = (groups, id) => {
        for (const group of groups) {
          if (group.id === id) return group.name
          if (group.children && group.children.length > 0) {
            const found = findGroup(group.children, id)
            if (found) return found
          }
        }
        return null
      }
      return findGroup(this.groupList, groupId) || '未知分组'
    },

    // 根据认证凭据ID获取认证类型名称
    getAuthTypeName(authId) {
      if (!authId) return '未设置'
      const auth = this.authList.find(item => item.id === authId)
      if (!auth) return '未知认证'
      switch (Number(auth.authType)) {
        case 1:
          return '密码认证'
        case 2:
          return '密钥认证'
        case 3:
          return '公钥免认证'
        default:
          return '未知认证'
      }
    },
    normalizeDeviceType(deviceType) {
      return String(deviceType || 'linux').toLowerCase() === 'windows' ? 'windows' : 'linux'
    },
    isWindowsHost(host) {
      return this.normalizeDeviceType(host?.deviceType) === 'windows'
    },
    canUseSSH(host) {
      return Boolean(host?.supportsSsh || (host?.sshIp && host?.sshName && Number(host?.sshKeyId) > 0))
    },
    canSyncHost(host) {
      return !this.isWindowsHost(host) && this.canUseSSH(host)
    },
    getDeviceTypeLabel(deviceType) {
      return this.isWindowsHost({ deviceType }) ? 'Windows主机' : 'Linux主机'
    },
    formatAccessAddress(host) {
      if (this.isWindowsHost(host)) {
        return host.remoteDomain || host.sshIp || '-'
      }
      if (!host.sshIp) {
        return '-'
      }
      return `${host.sshName || 'root'}@${host.sshIp}:${host.sshPort || 22}`
    },
    formatRDPAccess(host) {
      const target = host.remoteDomain || host.sshIp || '-'
      const user = host.rdpUsername || 'administrator'
      const domain = host.rdpDomain ? `${host.rdpDomain}\\` : ''
      return `${domain}${user} @ ${target}:${host.rdpPort || 3389}`
    },

    // 显示主机详情
    async showHostDetail(row) {
      try {
        const { data: res } = await this.$api.getCmdbHostById(row.id)

        if (res.code === 200) {
          this.hostDetail = res.data
          this.detailDrawer = true

          await Promise.all([
            this.fetchHostMonitorData(row.id),
            this.fetchRecentHostAudit(res.data)
          ])

          this.$nextTick(() => {
            this.initGaugeCharts()
          })
        } else {
          console.error('获取主机详情失败:', res.message)
          this.$message.error(res.message || '获取主机详情失败')
        }
      } catch (error) {
        console.error('获取主机详情失败:', error)
        this.$message.error('获取主机详情失败: ' + error.message)
      }
    },

    // 获取主机监控数据
    async fetchHostMonitorData(hostId) {
      try {
        const { data: res } = await this.$api.getHostsMonitorData(hostId)

        if (res.code === 200 && res.data) {
          const monitorData = res.data[hostId]
          if (monitorData) {
            this.hostDetail.cpuUsage = parseFloat(monitorData.cpuUsage?.toFixed(2) || 0)
            this.hostDetail.memoryUsage = parseFloat(monitorData.memoryUsage?.toFixed(2) || 0)
            this.hostDetail.diskUsage = parseFloat(monitorData.diskUsage?.toFixed(2) || 0)
          }
        }
      } catch (error) {
        console.error('获取主机监控数据失败:', error)
        this.hostDetail.cpuUsage = 0
        this.hostDetail.memoryUsage = 0
        this.hostDetail.diskUsage = 0
      }
    },

    // 关闭详情抽屉
    handleDetailClose() {
      this.detailDrawer = false
      this.recentAuditSessions = []
      this.destroyGaugeCharts()
    },

    handleCloudImportSuccess() {
      this.cloudDialogVisible = false
      this.getHostList()
    },

    // 处理Excel导入成功
    handleExcelImportSuccess() {
      this.getHostList()
    },

    // 连接SSH终端
    handleHostSSH(row = null) {
      const selectedHost = row || this.getTerminalTargetHost()
      if (!selectedHost) {
        return
      }
      if (!this.canUseSSH(selectedHost)) {
        this.$message.warning('当前主机未配置可用的 SSH 连接信息')
        return
      }
      this.$router.push({
        path: '/cmdb/ssh',
        query: {
          hostId: selectedHost.id
        }
      })
    },

    // 文件选择处理
    handleFileChange(file) {
      this.uploadForm.file = file.raw
      this.$refs.uploadFormRef.validateField('file')
    },

    // 文件删除处理
    handleFileRemove() {
      this.uploadForm.file = null
      this.$refs.uploadFormRef.validateField('file')
    },

    showUploadDialog(row) {
      if (!this.canUseSSH(row)) {
        this.$message.warning('当前主机未配置可用的 SSH 连接信息')
        return
      }
      this.currentUploadHost = row
      this.uploadForm = {
        hostId: row.id,
        file: null,
        targetPath: '/tmp'
      }
      this.$nextTick(() => {
        this.$refs.uploadFormRef?.clearValidate('targetPath')
      })
      this.uploadDialogVisible = true
    },

    handleDropdownVisibleChange(visible) {
      console.log('涓嬫媺妗嗘樉绀虹姸鎬佸彉鍖?', visible)
      if (visible) {
        const hasPermission = this.checkPermission(['cmdb:ecs:add'])
        console.log('鏉冮檺妫€鏌ョ粨鏋?', hasPermission)

        if (!hasPermission) {
          this.$message.warning('您没有新建主机的权限')
          this.$nextTick(() => {
            this.$refs.createDropdown?.hide()
          })
          return false
        }
      }
    },

    // 处理新建按钮点击（现在只是一个占位符，真正的逻辑在visible-change中）
    handleCreateClick() {
      console.log('触发新建主机下拉菜单')
    },

    // 澶勭悊涓嬫媺妗嗛€夐」鐐瑰嚮
    handleCreateCommand(command) {
      console.log('选择了下拉框选项:', command)

      switch (command) {
        case 'importHost':
          this.addDialogVisible = true
          break
        case 'excelImport':
          this.ExcelDialogVisible = true
          break
        case 'cloudHost':
          this.cloudDialogVisible = true
          break
      }
    },

    checkPermission(permissions) {
      console.log('checkPermission被调用，权限列表:', permissions)

      // 临时返回true用于测试
      console.log('鏉冮檺妫€鏌ラ€氳繃')
      return true

      // TODO: 瀹炵幇鐪熸鐨勬潈闄愭鏌ラ€昏緫
      // 鍋囪鎮ㄦ湁鍏ㄥ眬鐨勬潈闄愭鏌ユ柟娉?      // if (this.$checkPermission) {
      //   return this.$checkPermission(permissions)
      // }

      // 鎴栬€呮鏌tore涓殑鏉冮檺
      // if (this.$store && this.$store.getters.permissions) {
      //   const userPermissions = this.$store.getters.permissions
      //   return permissions.some(permission => userPermissions.includes(permission))
      // }
    },

    // 执行命令
    async executeCommand(row) {
      try {
        if (!this.canUseSSH(row)) {
          this.$message.warning('当前主机未配置可用的 SSH 连接信息')
          return
        }
        // 初始化commandDialog对象
        this.commandDialog = {
          visible: true,
          loading: false,
          previewLoading: false,
          command: '',
          output: '',
          status: '',
          hostName: row.hostName,
          riskAssessment: null,
          riskPresentation: null
        }
        
        this.currentHostId = row.id
        
        this.$nextTick(() => {
          this.commandDialog.visible = true
        })
      } catch (error) {
        console.error('命令执行初始化失败:', error)
        this.$message.error('命令执行初始化失败: ' + error.message)
      }
    },

    // 执行命令提交
    async submitCommand() {
      if (!this.commandDialog.command) {
        console.warn('命令内容为空')
        this.$message.warning('请输入要执行的命令')
        return
      }

      try {
        this.commandDialog.previewLoading = true
        this.commandDialog.status = ''
        this.commandDialog.output = ''

        const { data: previewRes } = await this.$api.previewHostCommandRisk(
          this.currentHostId,
          this.commandDialog.command
        )
        if (!previewRes || previewRes.code !== 200) {
          throw new Error(previewRes?.message || '命令风险预检失败')
        }

        this.commandDialog.riskAssessment = previewRes.data || {}
        this.commandDialog.riskPresentation = getCommandRiskPresentation(previewRes.data || {})
        this.commandDialog.previewLoading = false

        const executeOptions = {}
        if (this.commandDialog.riskPresentation.requiresConfirmation) {
          const confirmResult = await this.$confirm(
            `${this.commandDialog.riskPresentation.text}：${this.commandDialog.riskPresentation.reason || '请确认风险后再继续执行'}`,
            '命令风险确认',
            {
              confirmButtonText: '继续执行',
              cancelButtonText: '取消',
              type: this.commandDialog.riskPresentation.level >= 2 ? 'error' : 'warning'
            }
          ).catch(err => err)

          if (confirmResult !== 'confirm') {
            this.commandDialog.status = '已取消'
            this.commandDialog.output = '已取消执行'
            return
          }

          executeOptions.riskAck = true
          executeOptions.confirmedRiskLevel = this.commandDialog.riskPresentation.level
        }

        this.commandDialog.loading = true
        const { data: res } = await this.$api.executeHostCommand(
          this.currentHostId,
          this.commandDialog.command,
          executeOptions
        )

        if (res && res.code === 409) {
          this.commandDialog.riskAssessment = res.data || this.commandDialog.riskAssessment
          this.commandDialog.riskPresentation = getCommandRiskPresentation(res.data || this.commandDialog.riskAssessment || {})
          this.commandDialog.status = '需重新确认'
          this.commandDialog.output = res.message || '服务端要求重新确认风险后再执行'
          return
        }

        if (res && res.code === 200) {
          this.commandDialog.status = '执行成功'
          this.commandDialog.output = extractCommandOutput(res.data?.output || res.data) || '命令执行成功但无输出'
        } else {
          console.warn('命令执行失败:', res?.message)
          this.commandDialog.status = '执行失败'
          this.commandDialog.output = res?.message || '未知错误'
        }
      } catch (error) {
        console.error('API请求异常:', error)
        this.commandDialog.status = '请求失败'
        this.commandDialog.output = error.message || 'API请求异常'
      } finally {
        this.commandDialog.previewLoading = false
        this.commandDialog.loading = false
      }
    },

    // 文件上传处理
    async handleUpload() {
      try {
        // 验证表单
        await this.$refs.uploadFormRef.validate()
        if (!this.uploadForm.file) {
          return this.$message.warning('请选择上传文件')
        }
        if (!this.uploadForm.targetPath?.trim()) {
          return this.$message.warning('请输入目标路径')
        }
        if (this.isUploading) {
          return this.$message.warning('已有文件正在上传，请等待完成')
        }

        this.isUploading = true
        this.uploadProgress = 0

        // 纭繚鐩爣璺緞鏈夊€硷紝浣跨敤榛樿璺緞'/tmp'濡傛灉涓虹┖
        const destPath = this.uploadForm.targetPath || '/tmp'
        
        const formData = new FormData()
        formData.append('file', this.uploadForm.file)
        formData.append('destPath', destPath)

        const config = {
          headers: {
            'Content-Type': 'multipart/form-data'
          },
          timeout: 0,
          onUploadProgress: progressEvent => {
            const percentCompleted = Math.round(
              (progressEvent.loaded * 100) / progressEvent.total
            )
            this.uploadProgress = percentCompleted
          }
        }


        const { data: res } = await this.$api.uploadFileToHost(
          this.uploadForm.hostId,
          formData,
          config
        )
        if (res.code === 200) {
          this.$message.success('文件上传成功')
          this.uploadDialogVisible = false
          this.resetUploadForm()
        } else {
          this.$message.error(res.message || '文件上传失败')
        }
      } catch (error) {
        console.error('文件上传失败:', error)
        this.$message.error('文件上传失败: ' + (error.message || '未知错误'))
      } finally {
        this.isUploading = false
      }
    },

    // 重置上传表单
    resetUploadForm() {
      this.uploadForm = {
        hostId: this.currentUploadHost?.id || null,
        file: null,
        targetPath: '/tmp'
      }
      this.uploadProgress = 0
      this.$nextTick(() => {
        this.$refs.uploadFormRef?.clearValidate('targetPath')
      })
    },

    // 删除主机
    async handleHostDelete(row) {
      await this.deleteHosts([row], { batch: false })
    },

    async handleHostSync(targetHost = null) {
      try {
        const syncTarget = resolveHostSyncTarget(targetHost, this.hostDetail)
        if (!syncTarget || !syncTarget.id) {
          this.$message.warning('未找到可同步的主机')
          return
        }
        if (!this.canSyncHost(syncTarget)) {
          this.$message.warning('仅支持已配置 SSH 的 Linux 主机同步')
          return
        }

        this.syncLoading = true
        const { data: res } = await this.$api.syncHostConfig(syncTarget.id)

        if (this.isHostSyncAccepted(res)) {
          this.$message.success(res.data?.message || '已提交主机配置同步，请稍后刷新状态')
          setTimeout(() => {
            this.getHostList()
            if (this.detailDrawer && this.hostDetail.id === syncTarget.id) {
              this.showHostDetail({ id: syncTarget.id })
            }
          }, 3000)
        } else {
          this.$message.error(res.data?.message || res.message || '主机同步失败')
        }
      } catch (error) {
        console.error('主机同步失败:', error)
        this.$message.error('主机同步失败: ' + (error.message || '未知错误'))
      } finally {
        this.syncLoading = false
      }
    },

    // 分组管理 - 创建分组
    async handleCreateGroup(groupData) {
      try {
        console.log('创建分组数据:', groupData)
        const { data: res } = await cmdbApi.createCmdbGroup(groupData)
        
        if (res.code === 200) {
          this.$message.success('创建分组成功')
          // 刷新分组列表
          await this.getAllGroups()
        } else {
          this.$message.error(res.message || '创建分组失败')
        }
      } catch (error) {
        console.error('创建分组失败:', error)
        this.$message.error('创建分组失败: ' + (error.response?.data?.message || error.message))
      }
    },

    // 分组管理 - 更新分组
    async handleUpdateGroup(groupData) {
      try {
        console.log('更新分组数据:', groupData)
        const { data: res } = await cmdbApi.updateCmdbGroup(groupData)
        
        if (res.code === 200) {
          this.$message.success('更新分组成功')
          // 刷新分组列表
          await this.getAllGroups()
        } else {
          this.$message.error(res.message || '更新分组失败')
        }
      } catch (error) {
        console.error('更新分组失败:', error)
        this.$message.error('更新分组失败: ' + (error.response?.data?.message || error.message))
      }
    },

    // 分组管理 - 删除分组
    async handleDeleteGroup(groupId) {
      try {
        console.log('删除分组ID:', groupId)
        const { data: res } = await cmdbApi.deleteCmdbGroup(groupId)

        if (res.code === 200) {
          this.$message.success('删除分组成功')
          // 刷新分组列表
          await this.getAllGroups()
          // 如果删除的是当前选中的分组，重置选择
          if (this.currentGroupId === groupId) {
            this.currentGroupId = null
            this.getHostList()
          }
        } else {
          this.$message.error(res.message || '删除分组失败')
        }
      } catch (error) {
        console.error('删除分组失败:', error)
        this.$message.error('删除分组失败: ' + (error.response?.data?.message || error.message))
      }
    },

    // 初始化仪表盘图表
    initGaugeCharts() {
      this.$nextTick(() => {
        if (this.$refs.cpuGauge) {
          this.cpuChart = echarts.init(this.$refs.cpuGauge)
          this.cpuChart.setOption(this.getGaugeOption(this.hostDetail.cpuUsage || 0, 'CPU'))
        }

        // 初始化内存仪表盘
        if (this.$refs.memoryGauge) {
          this.memoryChart = echarts.init(this.$refs.memoryGauge)
          this.memoryChart.setOption(this.getGaugeOption(this.hostDetail.memoryUsage || 0, '内存'))
        }

        // 初始化磁盘仪表盘
        if (this.$refs.diskGauge) {
          this.diskChart = echarts.init(this.$refs.diskGauge)
          this.diskChart.setOption(this.getGaugeOption(this.hostDetail.diskUsage || 0, '磁盘'))
        }
      })
    },

    getGaugeOption(value, name) {
      return {
        series: [
          {
            type: 'gauge',
            startAngle: 225,
            endAngle: -45,
            min: 0,
            max: 100,
            radius: '85%',
            center: ['50%', '60%'],
            splitNumber: 10,
            axisLine: {
              lineStyle: {
                width: 10,
                color: [
                  [0.3, '#67e0e3'],
                  [0.7, '#37a2da'],
                  [1, '#fd666d']
                ]
              }
            },
            pointer: {
              length: '60%',
              width: 4,
              itemStyle: {
                color: '#4169E1'
              }
            },
            axisTick: {
              distance: -10,
              length: 4,
              lineStyle: {
                color: '#333',
                width: 1
              }
            },
            splitLine: {
              distance: -10,
              length: 8,
              lineStyle: {
                color: '#333',
                width: 2
              }
            },
            axisLabel: {
              color: '#666',
              distance: 20,
              fontSize: 12,
              formatter: function(value) {
                return value
              }
            },
            detail: {
              valueAnimation: true,
              formatter: '{value}%',
              color: '#ff0000',
              fontSize: 20,
              fontWeight: 'bold',
              offsetCenter: [0, '70%']
            },
            title: {
              offsetCenter: [0, '50%'],
              fontSize: 14,
              color: '#333',
              fontWeight: 'bold'
            },
            data: [
              {
                value: value,
                name: name
              }
            ]
          }
        ]
      }
    },

    destroyGaugeCharts() {
      if (this.cpuChart) {
        this.cpuChart.dispose()
        this.cpuChart = null
      }
      if (this.memoryChart) {
        this.memoryChart.dispose()
        this.memoryChart = null
      }
      if (this.diskChart) {
        this.diskChart.dispose()
        this.diskChart = null
      }
    }

  }
}

</script>

<style scoped>
/* 馃帹 鐜颁唬鍖栫鎶€鎰熻璁￠鏍?- 浠跨収cmdbDB.vue */

.cmdb-host-management {
  padding: 0;
  min-height: auto;
  background: transparent;
}

.host-card {
  background: rgba(12, 24, 41, 0.92);
  backdrop-filter: blur(18px);
  border-radius: 24px;
  box-shadow: var(--shadow-card);
  border: 1px solid var(--border-subtle);
}

.host-management-container {
  display: flex;
  height: calc(100vh - 180px);
}

.group-tree-section {
  width: 280px;
  margin-right: 12px;
}

.host-table-section {
  flex: 1;
  overflow-x: auto;
  overflow-y: visible;
  min-width: 0; /* 允许flex容器压缩 */
}

/* 🔍 搜索区域样式 */
.search-section {
  margin-bottom: 20px;
  padding: 20px;
  background: rgba(255, 255, 255, 0.03);
  border-radius: 18px;
  border: 1px solid var(--border-subtle);
}

/* 🎯 操作按钮区域 */
.action-section {
  margin-top: 15px;
  margin-bottom: 20px;
  padding-left: 0;
}

.table-section {
  margin-bottom: 15px;
}

.pagination-section {
  text-align: right;
  margin-top: 20px;
}

.batch-toolbar {
  margin: 16px 0 20px;
  padding: 14px 16px;
  border: 1px solid rgba(59, 130, 246, 0.18);
  border-radius: 16px;
  background: rgba(15, 23, 42, 0.72);
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
  flex-wrap: wrap;
}

.batch-toolbar__summary {
  color: #e2e8f0;
  font-size: 14px;
}

.batch-toolbar__actions {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}

.font-weight-bold {
  font-weight: bold;
}

.table-operation {
  display: flex;
  justify-content: space-around;
}

/* 瀹屽叏绉婚櫎琛ㄥ崟鍒嗗壊绾?*/
.el-dialog .el-form-item {
  border-bottom: none !important;
  margin-bottom: 12px;
  padding-bottom: 0;
}

/* 绉婚櫎琛屽拰鍒椾箣闂寸殑鍒嗗壊绾?*/
.el-row {
  border-bottom: none !important;
}

.el-col {
  border-right: none !important;
  padding-right: 0 !important;
  margin-right: 0 !important;
}

/* 绉婚櫎鏈€鍚庝竴涓垪鐨勫彸杈硅窛 */
.el-col:last-child {
  padding-right: 0 !important;
  margin-right: 0 !important;
}

/* 🎨 按钮样式优化 */
.el-button {
  border-radius: 8px;
  font-weight: 500;
  transition: all 0.3s ease;
}

.el-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

/* 馃摑 杈撳叆妗嗗拰閫夋嫨鍣ㄦ牱寮?*/
.el-input :deep(.el-input__wrapper),
.el-select :deep(.el-input__wrapper) {
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid var(--border-subtle);
  border-radius: 12px;
  box-shadow: none;
  transition: all 0.3s ease;
}

.el-input :deep(.el-input__wrapper):hover,
.el-select :deep(.el-input__wrapper):hover {
  border-color: #c0c4cc;
}

.el-input :deep(.el-input__wrapper.is-focus),
.el-select :deep(.el-input__wrapper.is-focus) {
  border-color: #667eea;
  box-shadow: 0 0 0 2px rgba(103, 126, 234, 0.2);
  background: rgba(255, 255, 255, 1);
}

.el-input :deep(.el-input__inner),
.el-select :deep(.el-input__inner) {
  background: transparent;
  border: none;
  color: #2c3e50;
}

/* 馃敡 琛ㄥ崟椤规牱寮?*/
.search-section .el-form-item {
  margin-bottom: 0;
  margin-right: 16px;
}

.search-section .el-form-item__label {
  color: #606266;
  font-weight: 500;
}

/* 馃枼锔?缁堢鎸夐挳娓愬彉钃濊壊鏍峰紡 */
.delete-host-btn-wrapper {
  display: inline-flex;
}

.delete-host-btn {
  background: linear-gradient(135deg, #ef4444 0%, #dc2626 100%) !important;
  border: none !important;
  color: white !important;
  box-shadow: 0 4px 12px rgba(239, 68, 68, 0.28);
}

.delete-host-btn:hover {
  background: linear-gradient(135deg, #dc2626 0%, #b91c1c 100%) !important;
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(220, 38, 38, 0.34);
}

/* 🎯 抽屉内容区域样式 - 减少顶部间距 */
.delete-host-btn.is-disabled,
.delete-host-btn.is-disabled:hover {
  background: rgba(148, 163, 184, 0.28) !important;
  border: none !important;
  color: rgba(255, 255, 255, 0.72) !important;
  box-shadow: none;
  transform: none;
}
.el-drawer :deep(.el-drawer__body) {
  padding-top: 10px;
}

/* 馃幆 浠〃鐩樻牱寮?*/
.dashboard-section {
  margin: 0;
  padding: 0;
}

.gauge-container {
  display: flex;
  justify-content: space-around;
  align-items: center;
  gap: 20px;
  margin-bottom: 5px;
}

.gauge-item {
  flex: 1;
  height: 180px;
  min-width: 0;
}

.detail-section {
  margin: 22px 0 0;
}

.detail-section__head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
}

.detail-section__head h3 {
  margin: 0;
}

.detail-section__actions {
  display: flex;
  align-items: center;
  gap: 10px;
}

.connection-center {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 12px;
}

.connection-card {
  padding: 16px;
  border-radius: 16px;
  border: 1px solid rgba(148, 163, 184, 0.16);
  background: rgba(15, 23, 42, 0.78);
  display: grid;
  gap: 12px;
}

.connection-card--disabled {
  background: rgba(30, 41, 59, 0.62);
}

.connection-card__header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 10px;
}

.connection-card__title {
  color: #f8fafc;
  font-weight: 600;
}

.connection-card__body {
  min-height: 48px;
  color: #cbd5e1;
  font-size: 13px;
  line-height: 1.6;
}

.connection-card__footer {
  display: flex;
  justify-content: flex-end;
}

.audit-summary-list {
  display: grid;
  gap: 12px;
}

.audit-summary-item {
  padding: 14px 16px;
  border-radius: 14px;
  border: 1px solid rgba(148, 163, 184, 0.16);
  background: rgba(15, 23, 42, 0.72);
  display: grid;
  gap: 10px;
}

.audit-summary-item__meta,
.audit-summary-item__footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 10px;
}

.audit-summary-item__title {
  color: #f8fafc;
  font-weight: 600;
}

.audit-summary-item__time {
  color: #94a3b8;
  font-size: 12px;
}

.audit-summary-item__command {
  color: #cbd5e1;
  font-size: 13px;
  line-height: 1.6;
  word-break: break-word;
}

.command-dialog :deep(.el-dialog) {
  width: min(72vw, 1100px);
  max-width: 1100px;
}

.command-dialog :deep(.el-dialog__body) {
  display: grid;
  gap: 18px;
}

.command-dialog__form {
  padding: 18px;
  border-radius: 18px;
  border: 1px solid var(--border-subtle);
  background: rgba(255, 255, 255, 0.03);
}

.command-output {
  padding: 18px;
  border-radius: 18px;
  border: 1px solid var(--border-subtle);
  background: rgba(2, 6, 23, 0.72);
}

.command-output__header {
  margin-bottom: 12px;
  font-size: 13px;
  font-weight: 700;
  color: var(--text-muted);
  letter-spacing: 0.04em;
}

.command-output__panel {
  margin: 0;
  min-height: 320px;
  max-height: 420px;
  overflow: auto;
  padding: 18px;
  border-radius: 14px;
  border: 1px solid rgba(148, 163, 184, 0.18);
  background: linear-gradient(180deg, rgba(2, 6, 23, 0.96), rgba(15, 23, 42, 0.94));
  color: #e2e8f0;
  font-family: Consolas, Monaco, 'Courier New', monospace;
  font-size: 13px;
  line-height: 1.65;
  white-space: pre-wrap;
  word-break: break-word;
}

.upload-dialog-hint {
  margin-top: 12px;
  color: #94a3b8;
  font-size: 12px;
  line-height: 1.6;
}
</style>





