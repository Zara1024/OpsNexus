<template>
  <div class="table-section">
    <el-table
      v-loading="loading"
      :data="hostListWithMonitor"
      stripe
      style="width: 100%"
      class="host-table"
      @selection-change="handleSelectionChange"
    >
      <el-table-column type="selection" width="48" />
      <el-table-column label="名称" min-width="220" show-overflow-tooltip>
        <template #default="scope">
          <div class="host-name-cell" @mouseenter="showCopyIcon($event)" @mouseleave="hideCopyIcon">
            <template v-if="isWindowsHost(scope.row)">
              <el-icon style="font-size: 18px; color: #409EFF; flex-shrink: 0;"><Monitor /></el-icon>
            </template>
            <template v-else>
              <img
                src="@/assets/image/linux.svg"
                alt="linux"
                style="height: 20px; object-fit: contain; flex-shrink: 0;"
              />
            </template>
            <el-link type="primary" @click="$emit('show-detail', scope.row)">{{ scope.row.name || '-' }}</el-link>
            <el-icon
              class="copy-icon"
              @click="copyToClipboard(scope.row.name || '', '名称', $event)"
              style="display: none; margin-left: 5px; cursor: pointer; color: #409EFF;"
            >
              <DocumentCopy />
            </el-icon>
          </div>
        </template>
      </el-table-column>

      <el-table-column label="地址" min-width="180" show-overflow-tooltip>
        <template #default="scope">
          <div class="asset-cell" @mouseenter="showCopyIcon($event)" @mouseleave="hideCopyIcon">
            <span>{{ scope.row.address || '-' }}</span>
            <el-icon
              class="copy-icon"
              @click="copyToClipboard(scope.row.address || '', '地址', $event)"
              style="display: none; margin-left: 5px; cursor: pointer; color: #409EFF;"
            >
              <DocumentCopy />
            </el-icon>
          </div>
        </template>
      </el-table-column>

      <el-table-column label="账号" min-width="140" show-overflow-tooltip>
        <template #default="scope">
          <div class="asset-cell" @mouseenter="showCopyIcon($event)" @mouseleave="hideCopyIcon">
            <span>{{ scope.row.account || '-' }}</span>
            <el-icon
              class="copy-icon"
              @click="copyToClipboard(scope.row.account || '', '账号', $event)"
              style="display: none; margin-left: 5px; cursor: pointer; color: #409EFF;"
            >
              <DocumentCopy />
            </el-icon>
          </div>
        </template>
      </el-table-column>

      <el-table-column label="平台" min-width="130" align="center">
        <template #default="scope">
          <el-tag :type="isWindowsHost(scope.row) ? 'primary' : 'success'" effect="plain">
            {{ formatPlatform(scope.row) }}
          </el-tag>
        </template>
      </el-table-column>

      <el-table-column label="连接性" min-width="110" align="center">
        <template #default="scope">
          <el-tag :type="scope.row.connectivity?.type || 'info'">
            {{ scope.row.connectivity?.text || '-' }}
          </el-tag>
        </template>
      </el-table-column>

      <el-table-column label="操作" fixed="right" width="300" min-width="300">
        <template #default="scope">
          <div class="table-operation">
            <el-button type="primary" link v-authority="['cmdb:ecs:edit']" @click="$emit('edit-host', scope.row.id)">
              编辑
            </el-button>
            <el-button type="info" link v-authority="['cmdb:ecs:monitor']" @click="showMonitor(scope.row)">
              监控
            </el-button>
            <el-dropdown trigger="click" @command="handleActionCommand($event, scope.row)">
              <el-button type="primary" link class="table-operation__more">
                更多
                <el-icon><ArrowDown /></el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="detail">详情</el-dropdown-item>
                  <el-dropdown-item
                    command="terminal"
                    :disabled="!canUseSSH(scope.row)"
                    v-authority="['cmdb:ecs:terminal']"
                  >
                    终端
                  </el-dropdown-item>
                  <el-dropdown-item command="audit" :disabled="!scope.row?.id">审计</el-dropdown-item>
                  <el-dropdown-item
                    command="sync"
                    :disabled="!canSyncHost(scope.row)"
                    v-authority="['cmdb:ecs:rsync']"
                  >
                    同步
                  </el-dropdown-item>
                  <el-dropdown-item command="process">进程监控</el-dropdown-item>
                  <el-dropdown-item command="port">TCP 端口</el-dropdown-item>
                  <el-dropdown-item
                    command="upload"
                    :disabled="!canUseSSH(scope.row)"
                    v-authority="['cmdb:ecs:upload']"
                  >
                    文件上传
                  </el-dropdown-item>
                  <el-dropdown-item
                    command="command"
                    :disabled="!canUseSSH(scope.row)"
                    v-authority="['cmdb:ecs:shell']"
                  >
                    执行命令
                  </el-dropdown-item>
                  <el-dropdown-item command="delete" divided v-authority="['cmdb:ecs:delete']">
                    删除主机
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
        </template>
      </el-table-column>
    </el-table>

    <monitor-dialog
      v-if="showMonitorDialog"
      v-model="showMonitorDialog"
      :host-id="currentHostId"
      style="z-index: 2001"
    />

    <process-monitor-dialog
      v-if="showProcessDialog"
      v-model="showProcessDialog"
      :host-id="currentProcessHostId"
      style="z-index: 2002"
    />

    <tcp-port-monitor-dialog
      v-if="showTcpPortDialog"
      v-model="showTcpPortDialog"
      :host-id="currentTcpPortHostId"
      style="z-index: 2003"
    />
  </div>
</template>

<script>
import MonitorDialog from './MonitorDialog.vue'
import ProcessMonitorDialog from './ProcessMonitorDialog.vue'
import TcpPortMonitorDialog from './TcpPortMonitorDialog.vue'
import { getHostMonitorUnavailableReason, hasHostMonitorData } from '@/utils/hostMonitorState.mjs'
import { mapHostRowToAssetRow } from '@/utils/cmdbAssetPresentation.mjs'

export default {
  name: 'CmdbHostTable',
  components: {
    MonitorDialog,
    ProcessMonitorDialog,
    TcpPortMonitorDialog
  },
  props: {
    hostList: {
      type: Array,
      required: true
    },
    loading: {
      type: Boolean,
      default: false
    }
  },
  data() {
    return {
      monitorData: {},
      monitorFetchInFlight: false,
      refreshInterval: null,
      refreshRate: 10000,
      showMonitorDialog: false,
      currentHostId: '',
      showProcessDialog: false,
      currentProcessHostId: '',
      showTcpPortDialog: false,
      currentTcpPortHostId: ''
    }
  },
  watch: {
    hostList: {
      immediate: true,
      deep: true,
      handler(newVal) {
        if (newVal && newVal.length > 0) {
          this.startRefresh()
        } else {
          this.stopRefresh()
        }
      }
    }
  },
  computed: {
    hostListWithMonitor() {
      return this.hostList.map((host) => {
        const monitor = this.monitorData[host.id] || {}
        const rawOnlineStatus = monitor?.onlineStatus
        const hasOnlineStatus = rawOnlineStatus !== undefined
          && rawOnlineStatus !== null
          && String(rawOnlineStatus).trim() !== ''
        const normalizedOnlineStatus = hasOnlineStatus ? Number(rawOnlineStatus) : undefined
        const isAlive = Number.isFinite(normalizedOnlineStatus)
          ? normalizedOnlineStatus === 0
          : undefined
        return mapHostRowToAssetRow({
          ...host,
          cpuUsage: monitor.cpuUsage,
          memoryUsage: monitor.memoryUsage,
          diskUsage: monitor.diskUsage,
          isAlive,
          monitorDataAvailable: hasHostMonitorData(monitor),
          monitorCollectionStatus: monitor.collectionStatus || 'unavailable',
          monitorUnavailableReason: getHostMonitorUnavailableReason(monitor)
        })
      })
    }
  },
  methods: {
    handleSelectionChange(selection) {
      this.$emit('selection-change', selection)
    },
    isWindowsHost(host) {
      return String(host?.deviceType || '').toLowerCase() === 'windows'
    },
    canUseSSH(host) {
      return Boolean(host?.supportsSsh || (host?.sshIp && host?.sshName && Number(host?.sshKeyId) > 0))
    },
    canSyncHost(host) {
      return !this.isWindowsHost(host) && this.canUseSSH(host)
    },
    formatPlatform(host) {
      const value = String(host?.platform || '').trim()
      if (!value) {
        return this.isWindowsHost(host) ? 'Windows' : 'Linux'
      }
      if (value.toLowerCase() === 'windows') return 'Windows'
      if (value.toLowerCase() === 'linux') return 'Linux'
      return value
    },
    handleActionCommand(command, row) {
      if (command === 'detail') {
        this.$emit('show-detail', row)
        return
      }
      if (command === 'process') {
        this.showProcessMonitor(row)
        return
      }
      if (command === 'port') {
        this.showTcpPortMonitor(row)
        return
      }
      if (command === 'terminal') {
        this.$emit('open-terminal', row)
        return
      }
      if (command === 'audit') {
        this.$emit('open-audit', row)
        return
      }
      if (command === 'sync') {
        this.$emit('sync-host', row)
        return
      }
      if (command === 'upload') {
        this.$emit('show-upload', row)
        return
      }
      if (command === 'command') {
        this.$emit('execute-command', row)
        return
      }
      if (command === 'delete') {
        this.$emit('delete-host', row)
      }
    },
    async fetchMonitorData() {
      if (this.monitorFetchInFlight) return
      if (!this.hostList || this.hostList.length === 0) return

      const validHosts = this.hostList.filter((host) => host?.id)
      if (validHosts.length === 0) return

      this.monitorFetchInFlight = true
      try {
        const ids = validHosts.map((host) => host.id).join(',')
        const monitorRes = await this.$api.getHostsMonitorData(ids)
        if (monitorRes?.data?.code === 200) {
          this.monitorData = {
            ...this.monitorData,
            ...monitorRes.data.data
          }
        } else {
          console.error('获取主机监控数据失败:', monitorRes?.data?.message || '未知错误')
        }
      } catch (error) {
        console.error('获取主机监控数据异常:', error)
      } finally {
        this.monitorFetchInFlight = false
      }
    },
    startRefresh() {
      if (this.refreshInterval) return
      this.fetchMonitorData()
      this.refreshInterval = setInterval(() => {
        this.fetchMonitorData()
      }, this.refreshRate)
    },
    stopRefresh() {
      if (this.refreshInterval) {
        clearInterval(this.refreshInterval)
        this.refreshInterval = null
      }
    },
    showMonitor(host) {
      this.currentHostId = host.id
      this.showMonitorDialog = true
    },
    showProcessMonitor(host) {
      this.currentProcessHostId = host.id
      this.showProcessDialog = true
    },
    showTcpPortMonitor(host) {
      this.currentTcpPortHostId = host.id
      this.showTcpPortDialog = true
    },
    showCopyIcon(event) {
      const icons = event.currentTarget.querySelectorAll('.copy-icon')
      icons.forEach((icon) => {
        icon.style.display = 'inline-block'
      })
    },
    hideCopyIcon(event) {
      const icons = event.currentTarget.querySelectorAll('.copy-icon')
      icons.forEach((icon) => {
        icon.style.display = 'none'
      })
    },
    async copyToClipboard(text, type, event) {
      if (!text) return

      try {
        await navigator.clipboard.writeText(text)
        this.$message.success(`${type} 已复制: ${text}`)
        if (event && event.target) {
          const icon = event.target.closest('.copy-icon')
          if (icon) {
            icon.classList.add('copied')
            setTimeout(() => {
              icon.classList.remove('copied')
            }, 1000)
          }
        }
      } catch (error) {
        const textArea = document.createElement('textarea')
        textArea.value = text
        document.body.appendChild(textArea)
        textArea.focus()
        textArea.select()
        try {
          document.execCommand('copy')
          this.$message.success(`${type} 已复制: ${text}`)
        } catch (fallbackError) {
          this.$message.error('复制失败，请手动复制')
        }
        document.body.removeChild(textArea)
      }
    }
  },
  mounted() {
    this.startRefresh()
  },
  beforeUnmount() {
    this.stopRefresh()
  },
  beforeRouteEnter(to, from, next) {
    next((vm) => {
      vm.stopRefresh()
      vm.startRefresh()
    })
  },
  beforeRouteUpdate(to, from, next) {
    this.stopRefresh()
    this.startRefresh()
    next()
  },
  activated() {
    this.stopRefresh()
    this.startRefresh()
  }
}
</script>

<style scoped>
.table-section {
  margin-bottom: 15px;
  width: 100%;
}

.host-table {
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
}

.host-table :deep(.el-table__header) {
  background: linear-gradient(135deg, #667eea, #764ba2);
}

.host-table :deep(.el-table__header th) {
  background: transparent !important;
  color: #2c3e50 !important;
  font-weight: 700 !important;
  border-bottom: none;
  text-shadow: 0 1px 2px rgba(255, 255, 255, 0.8);
}

.host-table :deep(.el-table__row) {
  transition: all 0.3s ease;
}

.host-table :deep(.el-table__row:hover) {
  background-color: rgba(103, 126, 234, 0.05) !important;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.host-name-cell,
.asset-cell {
  display: flex;
  align-items: center;
  gap: 8px;
  white-space: nowrap;
  overflow: hidden;
  position: relative;
}

.table-operation {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 6px;
  min-width: 0;
  white-space: nowrap;
}

.table-operation .el-button {
  margin: 0;
}

.table-operation__more {
  display: inline-flex;
  align-items: center;
  gap: 2px;
}

.copy-icon {
  opacity: 0;
  transition: all 0.3s ease;
  font-size: 14px !important;
  padding: 2px;
  border-radius: 4px;
}

.copy-icon:hover {
  background-color: rgba(64, 158, 255, 0.1);
  transform: scale(1.1);
}

.host-name-cell:hover .copy-icon,
.asset-cell:hover .copy-icon {
  opacity: 1;
  display: inline-block !important;
}

.copy-icon.copied {
  color: #67c23a !important;
  transform: scale(1.2);
}

.host-table :deep(.el-table__cell) {
  white-space: nowrap;
  overflow: hidden;
}

.host-table :deep(.cell) {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
</style>
