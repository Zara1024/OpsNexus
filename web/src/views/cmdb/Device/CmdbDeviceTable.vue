<template>
  <div class="table-section">
    <el-table
      v-loading="loading"
      :data="deviceList"
      stripe
      style="width: 100%"
      class="device-table"
    >
      <el-table-column label="名称" min-width="180" show-overflow-tooltip>
        <template #default="scope">
          <div class="device-name-cell">
            <el-link type="primary" @click="$emit('connect-device', scope.row)">
              {{ scope.row.name || '-' }}
            </el-link>
          </div>
        </template>
      </el-table-column>

      <el-table-column label="地址" min-width="150" show-overflow-tooltip>
        <template #default="scope">
          <div class="asset-cell" @mouseenter="showCopyIcon($event)" @mouseleave="hideCopyIcon($event)">
            <span>{{ scope.row.address || '-' }}</span>
            <el-icon
              class="copy-icon"
              @click="copyToClipboard(scope.row.address || '', '地址', $event)"
            >
              <DocumentCopy />
            </el-icon>
          </div>
        </template>
      </el-table-column>

      <el-table-column label="账号" min-width="130" show-overflow-tooltip>
        <template #default="scope">
          <div class="asset-cell" @mouseenter="showCopyIcon($event)" @mouseleave="hideCopyIcon($event)">
            <span>{{ scope.row.account || '-' }}</span>
            <el-icon
              class="copy-icon"
              @click="copyToClipboard(scope.row.account || '', '账号', $event)"
            >
              <DocumentCopy />
            </el-icon>
          </div>
        </template>
      </el-table-column>

      <el-table-column label="平台" min-width="110" align="center">
        <template #default="scope">
          <el-tag effect="plain" type="success">
            {{ scope.row.platform || '-' }}
          </el-tag>
        </template>
      </el-table-column>

      <el-table-column label="连接性" min-width="100" align="center">
        <template #default="scope">
          <el-tag :type="scope.row.connectivity?.type || 'info'">
            {{ scope.row.connectivity?.text || '-' }}
          </el-tag>
        </template>
      </el-table-column>

      <el-table-column label="操作" fixed="right" width="220" min-width="220">
        <template #default="scope">
          <div class="table-operation">
            <el-button type="primary" link @click="$emit('connect-device', scope.row)">
              连接
            </el-button>
            <el-button type="info" link @click="$emit('test-connectivity', scope.row)">
              测试连接
            </el-button>
            <el-button type="warning" link @click="$emit('edit-device', scope.row)">
              编辑
            </el-button>
            <el-button type="danger" link @click="$emit('delete-device', scope.row)">
              删除
            </el-button>
          </div>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script>
export default {
  name: 'CmdbDeviceTable',
  props: {
    deviceList: {
      type: Array,
      required: true
    },
    loading: {
      type: Boolean,
      default: false
    }
  },
  methods: {
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
      if (!text) {
        return
      }

      try {
        await navigator.clipboard.writeText(text)
        this.$message.success(`${type} 已复制: ${text}`)
        this.markCopied(event)
      } catch (error) {
        const textArea = document.createElement('textarea')
        textArea.value = text
        document.body.appendChild(textArea)
        textArea.focus()
        textArea.select()
        try {
          document.execCommand('copy')
          this.$message.success(`${type} 已复制: ${text}`)
          this.markCopied(event)
        } catch (fallbackError) {
          this.$message.error('复制失败，请手动复制')
        }
        document.body.removeChild(textArea)
      }
    },
    markCopied(event) {
      const icon = event?.target?.closest('.copy-icon')
      if (!icon) {
        return
      }
      icon.classList.add('copied')
      setTimeout(() => {
        icon.classList.remove('copied')
      }, 1000)
    }
  }
}
</script>

<style scoped>
.table-section {
  margin-bottom: 15px;
  width: 100%;
}

.device-table {
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
}

.device-table :deep(.el-table__header) {
  background: linear-gradient(135deg, #667eea, #764ba2);
}

.device-table :deep(.el-table__header th) {
  background: transparent !important;
  color: #2c3e50 !important;
  font-weight: 700 !important;
  border-bottom: none;
  text-shadow: 0 1px 2px rgba(255, 255, 255, 0.8);
}

.device-table :deep(.el-table__row) {
  transition: all 0.3s ease;
}

.device-table :deep(.el-table__row:hover) {
  background-color: rgba(103, 126, 234, 0.05) !important;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.device-name-cell,
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
  gap: 4px;
  min-width: 0;
  white-space: nowrap;
}

.table-operation .el-button {
  margin: 0;
}

.copy-icon {
  display: none;
  opacity: 0;
  transition: all 0.3s ease;
  font-size: 14px !important;
  padding: 2px;
  border-radius: 4px;
  cursor: pointer;
  color: #409eff;
}

.copy-icon:hover {
  background-color: rgba(64, 158, 255, 0.1);
  transform: scale(1.1);
}

.asset-cell:hover .copy-icon {
  opacity: 1;
  display: inline-block !important;
}

.copy-icon.copied {
  color: #67c23a !important;
  transform: scale(1.2);
}

.device-table :deep(.el-table__cell),
.device-table :deep(.cell) {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
</style>
