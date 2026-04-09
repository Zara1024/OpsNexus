<template>
  <div class="cmdb-device-management">
    <el-card shadow="hover" class="device-card">
      <div class="device-management-container">
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

        <div class="device-table-section">
          <div class="search-section">
            <el-form :inline="true" :model="queryParams">
              <el-form-item label="名称">
                <el-input
                  v-model="queryParams.name"
                  size="small"
                  placeholder="请输入设备名称"
                  clearable
                  style="width: 160px"
                  @keyup.enter="handleQuery"
                />
              </el-form-item>
              <el-form-item label="地址">
                <el-input
                  v-model="queryParams.address"
                  size="small"
                  placeholder="请输入设备地址"
                  clearable
                  style="width: 180px"
                  @keyup.enter="handleQuery"
                />
              </el-form-item>
              <div class="action-section">
                <el-button type="primary" size="small" @click="handleQuery">
                  <el-icon><Search /></el-icon>
                  <span style="margin-left: 4px">搜索</span>
                </el-button>
                <el-button type="warning" size="small" @click="resetQuery">
                  <el-icon><Refresh /></el-icon>
                  <span style="margin-left: 4px">重置</span>
                </el-button>
                <el-button type="success" size="small" @click="showAddDialog">
                  <el-icon><Plus /></el-icon>
                  <span style="margin-left: 4px">新建</span>
                </el-button>
              </div>
            </el-form>
          </div>

          <CmdbDeviceTable
            :device-list="pagedDevices"
            :loading="loading"
            @connect-device="handleConnectDevice"
            @test-connectivity="handleConnectivityTest"
            @edit-device="showEditDialog"
            @delete-device="handleDeleteDevice"
          />

          <div class="pagination-section">
            <el-pagination
              :current-page="queryParams.pageNum"
              :page-size="queryParams.pageSize"
              :page-sizes="[10, 20, 50, 100]"
              layout="total, sizes, prev, pager, next, jumper"
              :total="filteredDeviceRows.length"
              @size-change="handleSizeChange"
              @current-change="handleCurrentChange"
            />
          </div>
        </div>
      </div>

      <CreateDevice
        :visible="addDialogVisible"
        :group-list="groupList"
        :account-list="accountList"
        @close="addDialogVisible = false"
        @submit="createDevice"
      />

      <EditDevice
        :visible="editDialogVisible"
        :device-info="deviceInfo"
        :group-list="groupList"
        :account-list="accountList"
        @close="editDialogVisible = false"
        @submit="updateDevice"
      />
    </el-card>
  </div>
</template>

<script>
import cmdbAPI from '@/api/cmdb'
import configAPI from '@/api/config'
import { mapDeviceRowToAssetRow } from '@/utils/cmdbAssetPresentation.mjs'
import CmdbGroup from './Host/CmdbGroup.vue'
import CmdbDeviceTable from './Device/CmdbDeviceTable.vue'
import CreateDevice from './Device/CreateDevice.vue'
import EditDevice from './Device/EditDevice.vue'
import { Plus, Refresh, Search } from '@element-plus/icons-vue'

const DEVICE_FETCH_PAGE_SIZE = 100
const DEVICE_FETCH_MAX_PAGES = 1000
const SAFE_COMMAND_ADDRESS_PATTERN = /^(?!-)[A-Za-z0-9.:[\]-]+$/

function normalizeListResponse(payload) {
  if (Array.isArray(payload)) {
    return payload
  }
  if (Array.isArray(payload?.list)) {
    return payload.list
  }
  return []
}

function normalizeTotalResponse(payload) {
  const raw = Number(payload?.total)
  if (Number.isFinite(raw) && raw >= 0) {
    return raw
  }
  return null
}

function normalizeAccountResponse(payload) {
  if (Array.isArray(payload)) {
    return payload
  }
  if (Array.isArray(payload?.list)) {
    return payload.list
  }
  return []
}

function normalizeOpenableWebUrl(value) {
  const raw = String(value || '').trim()
  if (!raw) {
    return ''
  }

  try {
    const prepared = /^https?:\/\//i.test(raw) ? raw : `http://${raw}`
    const parsed = new URL(prepared)
    if (!['http:', 'https:'].includes(parsed.protocol)) {
      return ''
    }
    if (!parsed.hostname) {
      return ''
    }
    return parsed.toString()
  } catch (error) {
    return ''
  }
}

function sanitizeCommandAddress(value) {
  const raw = String(value || '').trim()
  if (!raw) {
    return ''
  }
  if (/\s/.test(raw)) {
    return ''
  }
  if (!SAFE_COMMAND_ADDRESS_PATTERN.test(raw)) {
    return ''
  }
  return raw
}

function extractGroupId(node, element) {
  if (element?.data?.id) return element.data.id
  if (element?.id) return element.id
  if (node?.key) return node.key
  return null
}

function openExternalWindow(url) {
  const popup = window.open(url, '_blank', 'noopener,noreferrer')
  if (!popup) {
    return false
  }

  try {
    popup.opener = null
  } catch (error) {
    // Ignore cross-window restrictions.
  }

  return !popup.closed
}

export default {
  name: 'CmdbDevice',
  components: {
    CmdbGroup,
    CmdbDeviceTable,
    CreateDevice,
    EditDevice,
    Plus,
    Refresh,
    Search
  },
  data() {
    return {
      loading: false,
      queryParams: {
        pageNum: 1,
        pageSize: 10,
        name: '',
        address: ''
      },
      groupList: [],
      expandedKeys: [],
      currentGroupId: null,
      rawDeviceList: [],
      accountList: [],
      addDialogVisible: false,
      editDialogVisible: false,
      deviceInfo: {}
    }
  },
  computed: {
    accountLabelById() {
      return this.accountList.reduce((result, account) => {
        result[String(account.id)] = account.alias || account.name || ''
        return result
      }, {})
    },
    deviceRows() {
      return this.rawDeviceList.map((device) => mapDeviceRowToAssetRow({
        ...device,
        accountName: this.accountLabelById[String(device.accountId)] || ''
      }))
    },
    filteredDeviceRows() {
      const keyword = String(this.queryParams.name || '').trim().toLowerCase()
      const addressKeyword = String(this.queryParams.address || '').trim().toLowerCase()

      return this.deviceRows.filter((device) => {
        const matchGroup = this.currentGroupId
          ? Number(device.groupId) === Number(this.currentGroupId)
          : true
        const matchName = keyword
          ? String(device.name || '').toLowerCase().includes(keyword)
          : true
        const matchAddress = addressKeyword
          ? String(device.address || '').toLowerCase().includes(addressKeyword)
          : true
        return matchGroup && matchName && matchAddress
      })
    },
    pagedDevices() {
      const start = (this.queryParams.pageNum - 1) * this.queryParams.pageSize
      const end = start + this.queryParams.pageSize
      return this.filteredDeviceRows.slice(start, end)
    }
  },
  async created() {
    await Promise.all([
      this.getAllGroups(),
      this.getAccountList()
    ])
    await this.getDeviceList()
  },
  methods: {
    async getAllGroups() {
      try {
        const { data: res } = await cmdbAPI.getAllCmdbGroups()
        if (res.code === 200) {
          this.groupList = Array.isArray(res.data) ? res.data : []
        } else {
          this.$message.error(res.message || '获取分组列表失败')
        }
      } catch (error) {
        console.error('获取分组列表失败:', error)
        this.$message.error(`获取分组列表失败: ${error.message || '未知错误'}`)
      }
    },
    async getAccountList() {
      try {
        const { data: res } = await configAPI.listAccountAuth({
          page: 1,
          pageSize: 100
        })
        if (res.code === 200) {
          this.accountList = normalizeAccountResponse(res.data)
        } else {
          this.$message.error(res.message || '获取账号列表失败')
        }
      } catch (error) {
        console.error('获取账号列表失败:', error)
        this.$message.error(`获取账号列表失败: ${error.message || '未知错误'}`)
      }
    },
    async fetchAllDevices() {
      const devices = []
      let page = 1
      let total = null
      let reachedMaxPages = false

      while (page <= DEVICE_FETCH_MAX_PAGES) {
        const { data: res } = await cmdbAPI.listDevices({
          page,
          pageSize: DEVICE_FETCH_PAGE_SIZE
        })

        if (res.code !== 200) {
          throw new Error(res.message || '获取设备列表失败')
        }

        const currentBatch = normalizeListResponse(res.data)
        total = normalizeTotalResponse(res.data)
        devices.push(...currentBatch)

        if (!currentBatch.length || currentBatch.length < DEVICE_FETCH_PAGE_SIZE) {
          break
        }

        if (total !== null && devices.length >= total) {
          break
        }

        page += 1
      }

      if (page > DEVICE_FETCH_MAX_PAGES) {
        reachedMaxPages = true
      }

      if (reachedMaxPages) {
        this.$message.warning(`设备列表抓取已达到最大页数 ${DEVICE_FETCH_MAX_PAGES}，当前结果可能不完整`)
      }

      return devices
    },
    async getDeviceList() {
      this.loading = true
      try {
        this.rawDeviceList = await this.fetchAllDevices()
      } catch (error) {
        console.error('获取设备列表失败:', error)
        this.$message.error(`获取设备列表失败: ${error.message || '未知错误'}`)
        this.rawDeviceList = []
      } finally {
        this.loading = false
      }
    },
    handleQuery() {
      this.queryParams.pageNum = 1
    },
    resetQuery() {
      this.queryParams = {
        pageNum: 1,
        pageSize: 10,
        name: '',
        address: ''
      }
      this.currentGroupId = null
    },
    handleSizeChange(size) {
      this.queryParams.pageSize = size
      this.queryParams.pageNum = 1
    },
    handleCurrentChange(page) {
      this.queryParams.pageNum = page
    },
    showAddDialog() {
      this.addDialogVisible = true
    },
    async createDevice(payload) {
      try {
        const { data: res } = await cmdbAPI.createDevice(payload)
        if (res.code === 200) {
          this.$message.success('设备创建成功')
          this.addDialogVisible = false
          await this.getDeviceList()
        } else {
          this.$message.error(res.message || '设备创建失败')
        }
      } catch (error) {
        console.error('设备创建失败:', error)
        this.$message.error(`设备创建失败: ${error.message || '未知错误'}`)
      }
    },
    async showEditDialog(row) {
      try {
        const { data: res } = await cmdbAPI.getDevice(row.id)
        if (res.code === 200) {
          this.deviceInfo = res.data || {}
          this.editDialogVisible = true
          return
        }
        this.$message.error(res.message || '获取设备详情失败')
      } catch (error) {
        console.error('获取设备详情失败:', error)
        this.$message.error(`获取设备详情失败: ${error.message || '未知错误'}`)
      }
    },
    async updateDevice(payload) {
      try {
        const { data: res } = await cmdbAPI.updateDevice(payload)
        if (res.code === 200) {
          this.$message.success('设备更新成功')
          this.editDialogVisible = false
          await this.getDeviceList()
        } else {
          this.$message.error(res.message || '设备更新失败')
        }
      } catch (error) {
        console.error('设备更新失败:', error)
        this.$message.error(`设备更新失败: ${error.message || '未知错误'}`)
      }
    },
    async handleDeleteDevice(row) {
      const confirmResult = await this.$confirm(
        `是否确认删除设备 "${row.name}"?`,
        '提示',
        {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        }
      ).catch((err) => err)

      if (confirmResult !== 'confirm') {
        return
      }

      try {
        const { data: res } = await cmdbAPI.deleteDevice(row.id)
        if (res.code === 200) {
          this.$message.success('设备删除成功')
          await this.getDeviceList()
        } else {
          this.$message.error(res.message || '设备删除失败')
        }
      } catch (error) {
        console.error('设备删除失败:', error)
        this.$message.error(`设备删除失败: ${error.message || '未知错误'}`)
      }
    },
    updateConnectivityState(deviceId, reachable) {
      this.rawDeviceList = this.rawDeviceList.map((device) => {
        if (Number(device.id) !== Number(deviceId)) {
          return device
        }
        return {
          ...device,
          reachable,
          connected: reachable,
          online: reachable
        }
      })
    },
    async handleConnectivityTest(row) {
      try {
        const { data: res } = await cmdbAPI.testDeviceConnectivity([row.id])
        if (res.code !== 200) {
          this.$message.error(res.message || '连接性测试失败')
          return
        }
        const item = res.data?.items?.[0]
        if (!item) {
          this.$message.warning('未返回设备连接性结果')
          return
        }

        this.updateConnectivityState(row.id, Boolean(item.reachable))
        if (item.reachable) {
          this.$message.success(`设备连接正常: ${item.displayAddress || item.address || row.address}`)
        } else {
          this.$message.warning(`设备不可达: ${item.reason || '未知原因'}`)
        }
      } catch (error) {
        console.error('设备连接性测试失败:', error)
        this.$message.error(`设备连接性测试失败: ${error.message || '未知错误'}`)
      }
    },
    async copyText(text, successMessage) {
      if (!text) {
        return false
      }

      try {
        await navigator.clipboard.writeText(text)
        this.$message.success(successMessage)
        return true
      } catch (error) {
        const textArea = document.createElement('textarea')
        textArea.value = text
        document.body.appendChild(textArea)
        textArea.focus()
        textArea.select()
        try {
          document.execCommand('copy')
          this.$message.success(successMessage)
          return true
        } catch (fallbackError) {
          this.$message.error('复制失败，请手动复制')
        } finally {
          document.body.removeChild(textArea)
        }
      }

      return false
    },
    async handleConnectDevice(row) {
      const protocolValues = String(row.protocolGroup || '')
        .split(',')
        .map((item) => item.trim().toLowerCase())
        .filter(Boolean)
      const protocols = new Set(protocolValues.length ? protocolValues : ['ssh'])
      const safeAddress = sanitizeCommandAddress(row.address)
      const openableWebUrl = normalizeOpenableWebUrl(row.webUrl)

      if (protocols.has('web') && openableWebUrl) {
        const opened = openExternalWindow(openableWebUrl)
        if (opened) {
          this.$message.success('已打开设备 Web 入口')
        } else {
          this.$message.warning('浏览器阻止了新窗口，请允许弹窗后重试')
        }
        return
      }

      if (protocols.has('ssh')) {
        if (!safeAddress) {
          this.$message.warning('设备地址格式不安全，无法生成 SSH 命令')
          return
        }
        const command = `ssh -p ${Number(row.sshPort || 22)} ${safeAddress}`
        await this.copyText(command, `SSH 连接命令已复制: ${command}`)
        return
      }

      if (protocols.has('telnet')) {
        if (!safeAddress) {
          this.$message.warning('设备地址格式不安全，无法生成 Telnet 命令')
          return
        }
        const command = `telnet ${safeAddress} ${Number(row.telnetPort || 23)}`
        await this.copyText(command, `Telnet 连接命令已复制: ${command}`)
        return
      }

      this.$message.warning('当前设备没有可用的连接入口，请检查协议组和访问字段')
    },
    async handleGroupSearch(searchText) {
      if (!searchText) {
        this.expandedKeys = []
        return
      }

      try {
        const { data: res } = await cmdbAPI.getCmdbGroupByName(searchText)
        if (res.code !== 200 || !res.data) {
          return
        }

        const cmdbGroupRef = this.$refs.cmdbGroup
        const tree = cmdbGroupRef ? cmdbGroupRef.$refs.groupTree : null
        if (!tree) {
          return
        }

        const findPath = (groups, targetId, path = []) => {
          for (const group of groups) {
            if (group.id === targetId) {
              return [...path, group.id]
            }
            if (Array.isArray(group.children) && group.children.length) {
              const foundPath = findPath(group.children, targetId, [...path, group.id])
              if (foundPath) {
                return foundPath
              }
            }
          }
          return null
        }

        const expandPath = findPath(this.groupList, res.data.id)
        if (expandPath) {
          this.expandedKeys = expandPath.slice(0, -1)
          this.$nextTick(() => {
            tree.setCurrentKey(res.data.id)
          })
        }
      } catch (error) {
        console.error('搜索分组失败:', error)
        this.$message.error(`搜索分组失败: ${error.message || '未知错误'}`)
      }
    },
    handleGroupClick(node, element) {
      const groupId = extractGroupId(node, element)
      if (!groupId) {
        this.$message.warning('无法获取分组 ID')
        return
      }
      this.currentGroupId = groupId
      this.queryParams.pageNum = 1
    },
    handleNodeExpand(data, node) {
      if (!this.expandedKeys.includes(node.key)) {
        this.expandedKeys.push(node.key)
      }
    },
    handleNodeCollapse(data, node) {
      this.expandedKeys = this.expandedKeys.filter((key) => key !== node.key)
    },
    handleCollapseAll() {
      this.expandedKeys = []
    },
    handleExpandAll() {
      const allKeys = []
      const collectKeys = (nodes) => {
        nodes.forEach((node) => {
          allKeys.push(node.id)
          if (node.children && node.children.length > 0) {
            collectKeys(node.children)
          }
        })
      }
      collectKeys(this.groupList)
      this.expandedKeys = allKeys
    },
    async handleCreateGroup(groupData) {
      try {
        const { data: res } = await cmdbAPI.createCmdbGroup(groupData)
        if (res.code === 200) {
          this.$message.success('创建分组成功')
          await this.getAllGroups()
        } else {
          this.$message.error(res.message || '创建分组失败')
        }
      } catch (error) {
        console.error('创建分组失败:', error)
        this.$message.error(`创建分组失败: ${error.message || '未知错误'}`)
      }
    },
    async handleUpdateGroup(groupData) {
      try {
        const { data: res } = await cmdbAPI.updateCmdbGroup(groupData)
        if (res.code === 200) {
          this.$message.success('更新分组成功')
          await this.getAllGroups()
        } else {
          this.$message.error(res.message || '更新分组失败')
        }
      } catch (error) {
        console.error('更新分组失败:', error)
        this.$message.error(`更新分组失败: ${error.message || '未知错误'}`)
      }
    },
    async handleDeleteGroup(groupId) {
      try {
        const { data: res } = await cmdbAPI.deleteCmdbGroup(groupId)
        if (res.code === 200) {
          this.$message.success('删除分组成功')
          await this.getAllGroups()
          if (Number(this.currentGroupId) === Number(groupId)) {
            this.currentGroupId = null
          }
        } else {
          this.$message.error(res.message || '删除分组失败')
        }
      } catch (error) {
        console.error('删除分组失败:', error)
        this.$message.error(`删除分组失败: ${error.message || '未知错误'}`)
      }
    }
  }
}
</script>

<style scoped>
.cmdb-device-management {
  padding: 0;
  min-height: auto;
  background: transparent;
}

.device-card {
  background: rgba(12, 24, 41, 0.92);
  backdrop-filter: blur(18px);
  border-radius: 24px;
  box-shadow: var(--shadow-card);
  border: 1px solid var(--border-subtle);
}

.device-management-container {
  display: flex;
  height: calc(100vh - 180px);
}

.device-table-section {
  flex: 1;
  overflow-x: auto;
  overflow-y: visible;
  min-width: 0;
}

.search-section {
  margin-bottom: 20px;
  padding: 20px;
  background: rgba(255, 255, 255, 0.03);
  border-radius: 18px;
  border: 1px solid var(--border-subtle);
}

.action-section {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

.pagination-section {
  text-align: right;
  margin-top: 20px;
}

.el-button {
  border-radius: 8px;
  font-weight: 500;
  transition: all 0.3s ease;
}

.el-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

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

.search-section .el-form-item {
  margin-bottom: 0;
  margin-right: 16px;
}

.search-section .el-form-item__label {
  color: #606266;
  font-weight: 500;
}
</style>
