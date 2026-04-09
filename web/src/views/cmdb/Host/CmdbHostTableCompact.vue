<template>
  <div class="table-section">
    <el-table
      v-loading="loading"
      :data="hostListWithMonitor"
      row-key="id"
      stripe
      style="width: 100%"
      class="host-table host-table--compact"
      @selection-change="handleSelectionChange"
    >
      <el-table-column type="selection" width="48" :reserve-selection="true" />
      <el-table-column min-width="160" show-overflow-tooltip>
        <template #header>
          <span>主机名称</span>
        </template>
        <template v-slot="scope">
          <div class="host-name-cell" @mouseenter="showCopyIcon($event)" @mouseleave="hideCopyIcon">
            <template v-if="isWindowsHost(scope.row)">
              <el-icon class="host-os-icon"><Monitor /></el-icon>
            </template>
            <template v-else>
              <img
                src="@/assets/image/linux.svg"
                alt="linux"
                class="host-os-icon host-os-icon--linux"
              />
            </template>
            <div class="host-name-copy">
              <el-link type="primary" @click="$emit('show-detail', scope.row)">{{ scope.row.hostName }}</el-link>
              <div class="host-name-copy__meta">{{ isWindowsHost(scope.row) ? 'Windows' : 'Linux' }}</div>
            </div>
            <el-icon
              class="copy-icon"
              @click="copyToClipboard(scope.row.hostName, '主机名称', $event)"
              style="display: none; margin-left: 5px; cursor: pointer; color: #409EFF;"
            >
              <DocumentCopy />
            </el-icon>
          </div>
        </template>
      </el-table-column>

      <el-table-column min-width="150" show-overflow-tooltip>
        <template #header>
          <span>IP地址</span>
        </template>
        <template v-slot="scope">
          <div class="ip-cell" @mouseenter="showCopyIcon($event)" @mouseleave="hideCopyIcon">
            <div v-if="scope.row.publicIp" class="ip-row public-ip">
              <span class="ip-badge ip-badge--public">PUB</span>
              <span>{{ scope.row.publicIp || '无公网IP' }}</span>
              <el-icon
                class="copy-icon"
                @click="copyToClipboard(scope.row.publicIp, '公网IP', $event)"
                style="display: none; margin-left: 5px; cursor: pointer; color: #409EFF;"
              >
                <DocumentCopy />
              </el-icon>
            </div>
            <div v-if="scope.row.privateIp" class="ip-row private-ip">
              <span class="ip-badge ip-badge--private">PRI</span>
              <span>{{ scope.row.privateIp || '无内网IP' }}</span>
              <el-icon
                class="copy-icon"
                @click="copyToClipboard(scope.row.privateIp, '内网IP', $event)"
                style="display: none; margin-left: 5px; cursor: pointer; color: #67C23A;"
              >
                <DocumentCopy />
              </el-icon>
            </div>
          </div>
        </template>
      </el-table-column>

      <el-table-column min-width="170">
        <template #header>
          <span>资源概览</span>
        </template>
        <template v-slot="scope">
          <div class="resource-stack">
            <div class="resource-stack__item">
              <span class="resource-stack__label">CPU</span>
              <el-progress
                :percentage="scope.row.cpuUsage || 0"
                :color="getUsageColor(scope.row.cpuUsage)"
                :show-text="false"
                :stroke-width="6"
              />
              <span class="resource-stack__value">{{ formatUsage(scope.row.cpuUsage) }}</span>
            </div>
            <div class="resource-stack__item">
              <span class="resource-stack__label">内存</span>
              <el-progress
                :percentage="scope.row.memoryUsage || 0"
                :color="getUsageColor(scope.row.memoryUsage)"
                :show-text="false"
                :stroke-width="6"
              />
              <span class="resource-stack__value">{{ formatUsage(scope.row.memoryUsage) }}</span>
            </div>
            <div class="resource-stack__item">
              <span class="resource-stack__label">磁盘</span>
              <el-progress
                :percentage="scope.row.diskUsage || 0"
                :color="getUsageColor(scope.row.diskUsage)"
                :show-text="false"
                :stroke-width="6"
              />
              <span class="resource-stack__value">{{ formatUsage(scope.row.diskUsage) }}</span>
            </div>
          </div>
        </template>
      </el-table-column>

      <el-table-column min-width="100" show-overflow-tooltip>
        <template #header>
          <span>配置</span>
        </template>
        <template v-slot="scope">
          <div class="config-cell">
            <span class="config-pill">CFG</span>
            <span class="config-text">{{ formatConfigSpec(scope.row) }}</span>
          </div>
        </template>
      </el-table-column>

      <el-table-column min-width="100">
        <template #header>
          <span>状态</span>
        </template>
        <template v-slot="scope">
          <div class="status-stack">
            <div class="status-inline">
              <span :class="['status-dot', scope.row.isAlive ? 'is-online' : 'is-offline']"></span>
              <span class="status-text">{{ scope.row.isAlive ? '在线' : '离线' }}</span>
            </div>
            <span :class="['status-chip', `status-chip--${getStatusTagType(scope.row.status)}`]">
              {{ getStatusText(scope.row.status) }}
            </span>
          </div>
        </template>
      </el-table-column>

      <el-table-column min-width="170">
        <template #header>
          <span>操作</span>
        </template>
        <template v-slot="scope">
          <div class="table-operation">
            <el-button type="primary" link v-authority="['cmdb:ecs:edit']" @click="$emit('edit-host', scope.row.id)">
              编辑
            </el-button>
            <el-button type="info" link v-authority="['cmdb:ecs:monitor']" @click="showMonitor(scope.row)">
              监控
            </el-button>
            <el-button type="danger" link v-authority="['cmdb:ecs:delete']" @click="$emit('delete-host', scope.row)">
              删除主机
            </el-button>
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
import LegacyHostTable from './CmdbHostTable.vue'
import { formatHostResourceValue } from '@/utils/cmdbHostDetailPresentation.mjs'

export default {
  name: 'CmdbHostTableCompact',
  extends: LegacyHostTable,
  methods: {
    formatConfigSpec(row = {}) {
      const cpu = formatHostResourceValue(row.cpu, 'C')
      const memory = formatHostResourceValue(row.memory, 'G')
      return `${cpu} / ${memory}`
    },
    formatUsage(value) {
      const numeric = Number(value)
      if (!Number.isFinite(numeric)) {
        return '0%'
      }
      return `${numeric.toFixed(1)}%`
    },
    getUsageColor(value) {
      const numeric = Number(value)
      if (!Number.isFinite(numeric) || numeric < 60) {
        return '#22c55e'
      }
      if (numeric < 85) {
        return '#f59e0b'
      }
      return '#f87171'
    },
    getStatusTagType(status) {
      const normalized = Number(status)
      if (normalized === 1) {
        return 'success'
      }
      if (normalized === 2) {
        return 'warning'
      }
      if (normalized === 3) {
        return 'danger'
      }
      return 'info'
    },
    getStatusText(status) {
      const normalized = Number(status)
      if (normalized === 1) {
        return '认证成功'
      }
      if (normalized === 2) {
        return '未认证'
      }
      if (normalized === 3) {
        return '认证失败'
      }
      return '未知状态'
    }
  }
}
</script>

<style scoped>
.table-section {
  margin-bottom: 15px;
  width: 100%;
  padding: 0;
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

.host-table :deep(.el-table__header th .cell) {
  padding: 0 4px;
  font-size: 13px;
}

.host-table :deep(.el-table__row) {
  transition: all 0.3s ease;
}

.host-table :deep(.el-table__row:hover) {
  background-color: rgba(103, 126, 234, 0.05) !important;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.host-table :deep(.el-table__cell) {
  padding: 10px 8px;
}

.host-name-cell {
  display: flex;
  align-items: center;
  gap: 8px;
  white-space: nowrap;
  overflow: hidden;
  position: relative;
}

.host-os-icon {
  font-size: 18px;
  color: #409eff;
  flex-shrink: 0;
}

.host-os-icon--linux {
  height: 20px;
  object-fit: contain;
}

.host-name-copy {
  min-width: 0;
  display: grid;
  gap: 4px;
}

.host-name-copy__meta {
  font-size: 11px;
  color: #7a8291;
  line-height: 1;
}

.ip-cell {
  white-space: nowrap;
  overflow: hidden;
}

.ip-row {
  display: flex;
  align-items: center;
  gap: 4px;
  white-space: nowrap;
  font-size: 12px;
  line-height: 1.2;
  position: relative;
}

.ip-row.public-ip {
  color: #409eff;
  margin-bottom: 2px;
}

.ip-row.private-ip {
  color: #67c23a;
}

.ip-row span {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.ip-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 28px;
  padding: 1px 6px;
  border-radius: 999px;
  font-size: 10px;
  font-weight: 700;
  letter-spacing: 0.04em;
  flex-shrink: 0;
}

.ip-badge--public {
  background: rgba(64, 158, 255, 0.12);
  color: #409eff;
}

.ip-badge--private {
  background: rgba(103, 194, 58, 0.12);
  color: #67c23a;
}

.resource-stack {
  display: grid;
  gap: 6px;
}

.resource-stack__item {
  display: grid;
  grid-template-columns: 30px minmax(0, 1fr) 34px;
  gap: 6px;
  align-items: center;
}

.resource-stack__label,
.resource-stack__value {
  font-size: 11px;
  color: #94a3b8;
  line-height: 1;
}

.resource-stack__value {
  text-align: right;
  font-weight: 600;
  color: #e2e8f0;
}

.config-cell {
  display: flex;
  align-items: center;
  gap: 4px;
  min-width: 0;
  white-space: nowrap;
  overflow: hidden;
}

.config-pill {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 28px;
  padding: 1px 6px;
  border-radius: 999px;
  background: rgba(103, 126, 234, 0.12);
  color: #667eea;
  font-size: 10px;
  font-weight: 700;
  letter-spacing: 0.04em;
  flex-shrink: 0;
}

.config-text {
  min-width: 0;
  font-size: 11px;
  color: #cbd5e1;
  overflow: hidden;
  text-overflow: ellipsis;
}

.status-stack {
  display: grid;
  gap: 4px;
}

.status-inline {
  display: flex;
  align-items: center;
  min-width: 0;
  white-space: nowrap;
}

.status-dot {
  width: 8px;
  height: 8px;
  margin-right: 6px;
  border-radius: 50%;
  flex-shrink: 0;
}

.status-dot.is-online {
  background: #67c23a;
  box-shadow: 0 0 0 3px rgba(103, 194, 58, 0.16);
}

.status-dot.is-offline {
  background: #f56c6c;
  box-shadow: 0 0 0 3px rgba(245, 108, 108, 0.16);
}

.status-text {
  min-width: 0;
  font-size: 11px;
  font-weight: 600;
  color: #cbd5e1;
}

.status-chip {
  display: inline-flex;
  align-items: center;
  width: fit-content;
  max-width: 100%;
  padding: 1px 6px;
  border-radius: 999px;
  font-size: 10px;
  font-weight: 700;
  line-height: 1.2;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.status-chip--success {
  color: #86efac;
  background: rgba(34, 197, 94, 0.16);
}

.status-chip--danger {
  color: #fca5a5;
  background: rgba(248, 113, 113, 0.14);
}

.status-chip--warning {
  color: #fcd34d;
  background: rgba(245, 158, 11, 0.14);
}

.status-chip--info,
.status-chip--primary {
  color: #bfdbfe;
  background: rgba(59, 130, 246, 0.14);
}

.table-operation {
  display: flex;
  align-items: center;
  justify-content: flex-start;
  gap: 6px;
  min-width: 0;
  white-space: nowrap;
  width: 100%;
}

.table-operation .el-button {
  margin: 0;
  font-size: 12px;
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
.ip-row:hover .copy-icon {
  opacity: 1;
  display: inline-block !important;
}

.copy-icon.copied {
  color: #67c23a !important;
  transform: scale(1.2);
}

.host-table :deep(.cell) {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.host-table :deep(.el-progress-bar__outer) {
  background: rgba(148, 163, 184, 0.18);
}
</style>
