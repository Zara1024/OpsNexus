<template>
  <div class="cmdb-db-management">
    <el-card shadow="hover" class="db-card">
      <div class="db-management-container">
        <div class="db-table-section">
          <div class="search-section">
            <el-form :inline="true" :model="queryParams">
              <el-form-item label="名称">
                <el-input
                  v-model="queryParams.name"
                  size="small"
                  placeholder="请输入数据库资产名称"
                  clearable
                  style="width: 180px"
                  @keyup.enter="handleQuery"
                />
              </el-form-item>
              <el-form-item label="地址">
                <el-input
                  v-model="queryParams.address"
                  size="small"
                  placeholder="请输入数据库地址"
                  clearable
                  style="width: 200px"
                  @keyup.enter="handleQuery"
                />
              </el-form-item>
              <el-form-item label="平台">
                <el-select
                  v-model="queryParams.type"
                  size="small"
                  clearable
                  placeholder="全部平台"
                  style="width: 180px"
                >
                  <el-option
                    v-for="item in dbTypeOptions"
                    :key="item.value"
                    :label="item.label"
                    :value="item.value"
                  />
                </el-select>
              </el-form-item>
              <div class="action-section">
                <el-button type="primary" size="small" @click="handleQuery">
                  <el-icon><Search /></el-icon>
                  <span class="button-label">搜索</span>
                </el-button>
                <el-button type="warning" size="small" @click="resetQuery">
                  <el-icon><Refresh /></el-icon>
                  <span class="button-label">重置</span>
                </el-button>
                <el-button type="success" size="small" @click="showAddDialog">
                  <el-icon><Plus /></el-icon>
                  <span class="button-label">新建数据库资产</span>
                </el-button>
              </div>
            </el-form>
          </div>

          <div class="table-section">
            <el-table
              v-loading="loading"
              :data="pagedDatabaseRows"
              stripe
              style="width: 100%"
              class="db-table"
            >
              <el-table-column label="名称" min-width="240" show-overflow-tooltip>
                <template #default="scope">
                  <div class="db-name-cell">
                    <img
                      :src="getDbIcon(scope.row.type)"
                      :alt="getDbName(scope.row.type)"
                      class="db-icon"
                    />
                    <div class="db-name-content">
                      <el-link type="primary" @click="goToDetails(scope.row)">
                        {{ scope.row.name || '-' }}
                      </el-link>
                      <div class="db-name-subtitle">
                        默认库: {{ scope.row.defaultDatabase || '-' }}
                      </div>
                    </div>
                  </div>
                </template>
              </el-table-column>

              <el-table-column label="地址" min-width="220" show-overflow-tooltip>
                <template #default="scope">
                  <span>{{ scope.row.address || '-' }}</span>
                </template>
              </el-table-column>

              <el-table-column label="账号" min-width="160" show-overflow-tooltip>
                <template #default="scope">
                  <span>{{ scope.row.account || '-' }}</span>
                </template>
              </el-table-column>

              <el-table-column label="平台" min-width="140" align="center">
                <template #default="scope">
                  <el-tag :type="getDbTagType(scope.row.type)" effect="plain">
                    {{ scope.row.platform || '-' }}
                  </el-tag>
                </template>
              </el-table-column>

              <el-table-column label="连接性" min-width="120" align="center">
                <template #default="scope">
                  <el-tag :type="scope.row.connectivity?.type || 'info'">
                    {{ scope.row.connectivity?.text || '-' }}
                  </el-tag>
                </template>
              </el-table-column>

              <el-table-column label="操作" fixed="right" width="320" min-width="320">
                <template #default="scope">
                  <div class="table-operation">
                    <el-button type="primary" link @click="goToDetails(scope.row)">
                      数据库操作
                    </el-button>
                    <el-button type="info" link @click="goToWorkOrders(scope.row)">
                      SQL工单
                    </el-button>
                    <el-button type="warning" link @click="showEditDialog(scope.row)">
                      编辑
                    </el-button>
                    <el-button type="danger" link @click="handleDeleteDatabase(scope.row)">
                      删除
                    </el-button>
                  </div>
                </template>
              </el-table-column>
            </el-table>
          </div>

          <div class="pagination-section">
            <el-pagination
              :current-page="queryParams.pageNum"
              :page-size="queryParams.pageSize"
              :page-sizes="[10, 20, 50, 100]"
              layout="total, sizes, prev, pager, next, jumper"
              :total="displayTotal"
              @size-change="handleSizeChange"
              @current-change="handleCurrentChange"
            />
          </div>
        </div>
      </div>

      <el-dialog
        v-model="dialogVisible"
        :title="dialogTitle"
        width="52%"
        @close="handleDialogClose"
      >
        <el-form
          ref="formRef"
          :model="formData"
          :rules="formRules"
          label-width="100px"
        >
          <div class="asset-form-section">
            <div class="asset-form-section__title">基础信息</div>
            <el-row :gutter="16">
              <el-col :span="12">
                <el-form-item label="名称" prop="name">
                  <el-input v-model="formData.name" placeholder="请输入数据库资产名称" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="地址" prop="address">
                  <el-input v-model="formData.address" placeholder="例如 10.0.0.31:3306" />
                </el-form-item>
              </el-col>
            </el-row>
            <el-row :gutter="16">
              <el-col :span="12">
                <el-form-item label="平台" prop="platform">
                  <el-input v-model="formData.platform" placeholder="例如 MySQL / PostgreSQL" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="类型" prop="type">
                  <el-select
                    v-model="formData.type"
                    placeholder="请选择数据库类型"
                    style="width: 100%"
                    @change="handleTypeChange"
                  >
                    <el-option
                      v-for="item in dbTypeOptions"
                      :key="item.value"
                      :label="item.label"
                      :value="item.value"
                    />
                  </el-select>
                </el-form-item>
              </el-col>
            </el-row>
            <el-row :gutter="16">
              <el-col :span="12">
                <el-form-item label="分组/节点" prop="groupId">
                  <el-select
                    v-model="formData.groupId"
                    filterable
                    clearable
                    placeholder="请选择分组/节点"
                    style="width: 100%"
                  >
                    <el-option
                      v-for="group in selectableGroups"
                      :key="group.id"
                      :label="group.displayName"
                      :value="group.id"
                    />
                  </el-select>
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="默认数据库" prop="defaultDatabase">
                  <el-input v-model="formData.defaultDatabase" placeholder="例如 orders" />
                </el-form-item>
              </el-col>
            </el-row>
          </div>

          <div class="asset-form-section">
            <div class="asset-form-section__title">访问设置</div>
            <el-row :gutter="16">
              <el-col :span="12">
                <el-form-item label="协议组" prop="protocolGroup">
                  <el-input v-model="formData.protocolGroup" placeholder="例如 default / readonly" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="账号" prop="accountId">
                  <el-select
                    v-model="formData.accountId"
                    filterable
                    clearable
                    placeholder="请选择账号"
                    style="width: 100%"
                  >
                    <el-option
                      v-for="account in accountList"
                      :key="account.id"
                      :label="account.alias || account.name || String(account.id)"
                      :value="account.id"
                    />
                  </el-select>
                </el-form-item>
              </el-col>
            </el-row>
          </div>

          <div class="asset-form-section">
            <div class="asset-form-section__title">补充信息</div>
            <el-row :gutter="16">
              <el-col :span="12">
                <el-form-item label="标签" prop="tags">
                  <el-input v-model="formData.tags" placeholder="多个标签可用逗号分隔" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="激活中" prop="isActive">
                  <el-switch
                    v-model="formData.isActive"
                    inline-prompt
                    active-text="是"
                    inactive-text="否"
                  />
                </el-form-item>
              </el-col>
            </el-row>
            <el-row :gutter="16">
              <el-col :span="24">
                <el-form-item label="备注" prop="remark">
                  <el-input
                    v-model="formData.remark"
                    type="textarea"
                    :rows="3"
                    placeholder="补充数据库资产的用途或说明"
                  />
                </el-form-item>
              </el-col>
            </el-row>
          </div>
        </el-form>

        <template #footer>
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitForm">确定</el-button>
        </template>
      </el-dialog>
    </el-card>
  </div>
</template>

<script>
import cmdbAPI from '@/api/cmdb'
import configAPI from '@/api/config'
import { createDatabaseAssetFormModel, mapDatabaseRowToAssetRow } from '@/utils/cmdbAssetPresentation.mjs'
import {
  getCmdbDatabasePlatformLabel,
  getCmdbDatabaseTypeIconKey,
  getCmdbDatabaseTypeLabel,
  getCmdbDatabaseTypeOptions,
  getCmdbDatabaseTypeTag
} from '@/utils/cmdbPresentation.mjs'
import { Plus, Refresh, Search } from '@element-plus/icons-vue'

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

function flattenSelectableGroups(groupList = []) {
  const groups = []

  const walk = (items) => {
    items.forEach((group) => {
      if (group.isDefault || !group.parentId) {
        groups.push({
          id: group.id,
          displayName: group.name
        })
      }
      if (Array.isArray(group.children) && group.children.length) {
        group.children.forEach((child) => {
          if (!child.children || !child.children.length) {
            groups.push({
              id: child.id,
              displayName: child.name
            })
          }
        })
      }
    })
  }

  walk(groupList)
  return groups
}

function trimString(value) {
  return String(value ?? '').trim()
}

function normalizeOptionalNumber(value, fallback = '') {
  if (value === undefined || value === null || value === '') {
    return fallback
  }
  const normalized = Number(value)
  return Number.isFinite(normalized) ? normalized : fallback
}

function normalizeDatabaseForm(database = {}) {
  const base = createDatabaseAssetFormModel()
  return {
    ...base,
    ...database,
    id: normalizeOptionalNumber(database?.id, base.id),
    name: database?.name ?? base.name,
    address: database?.address ?? database?.endpoint ?? base.address,
    platform: database?.platform || getCmdbDatabasePlatformLabel(database?.type, ''),
    groupId: normalizeOptionalNumber(database?.groupId, base.groupId),
    defaultDatabase: database?.defaultDatabase ?? base.defaultDatabase,
    protocolGroup: database?.protocolGroup || base.protocolGroup,
    accountId: normalizeOptionalNumber(database?.accountId, base.accountId),
    tags: database?.tags ?? base.tags,
    isActive: database?.isActive ?? base.isActive,
    remark: database?.remark || database?.description || base.remark,
    description: database?.description || database?.remark || base.description,
    type: normalizeOptionalNumber(database?.type, base.type)
  }
}

export default {
  name: 'CmdbDB',
  components: {
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
        address: '',
        type: undefined
      },
      rawDatabaseList: [],
      total: 0,
      useClientPagination: false,
      accountList: [],
      groupList: [],
      dbTypeOptions: getCmdbDatabaseTypeOptions(),
      dialogVisible: false,
      dialogMode: 'create',
      formData: createDatabaseAssetFormModel(),
      formRules: {
        name: [{ required: true, message: '请输入名称', trigger: 'blur' }],
        address: [{ required: true, message: '请输入地址', trigger: 'blur' }],
        platform: [{ required: true, message: '请输入平台', trigger: 'blur' }],
        type: [{ required: true, message: '请选择类型', trigger: 'change' }],
        groupId: [{ required: true, message: '请选择分组/节点', trigger: 'change' }],
        accountId: [{ required: true, message: '请选择账号', trigger: 'change' }]
      }
    }
  },
  computed: {
    dialogTitle() {
      return this.dialogMode === 'edit' ? '编辑数据库资产' : '新建数据库资产'
    },
    selectableGroups() {
      return flattenSelectableGroups(this.groupList)
    },
    accountLabelById() {
      return this.accountList.reduce((result, account) => {
        result[String(account.id)] = account.alias || account.name || ''
        return result
      }, {})
    },
    hasBackendSearchHint() {
      return Boolean(trimString(this.queryParams.name) || this.queryParams.type)
    },
    databaseRows() {
      return this.rawDatabaseList.map((item) => mapDatabaseRowToAssetRow({
        ...item,
        accountName: this.accountLabelById[String(item.accountId)] || item.accountName
      }))
    },
    filteredDatabaseRows() {
      const nameKeyword = trimString(this.queryParams.name).toLowerCase()
      const addressKeyword = trimString(this.queryParams.address).toLowerCase()
      const selectedType = this.queryParams.type

      return this.databaseRows.filter((item) => {
        const matchName = nameKeyword
          ? String(item.name || '').toLowerCase().includes(nameKeyword)
          : true
        const matchAddress = addressKeyword
          ? String(item.address || '').toLowerCase().includes(addressKeyword)
          : true
        const matchType = selectedType
          ? Number(item.type) === Number(selectedType)
          : true
        return matchName && matchAddress && matchType
      })
    },
    pagedDatabaseRows() {
      if (!this.useClientPagination) {
        return this.filteredDatabaseRows
      }

      const start = (this.queryParams.pageNum - 1) * this.queryParams.pageSize
      const end = start + this.queryParams.pageSize
      return this.filteredDatabaseRows.slice(start, end)
    },
    displayTotal() {
      return this.useClientPagination ? this.filteredDatabaseRows.length : this.total
    }
  },
  watch: {
    filteredDatabaseRows() {
      if (this.useClientPagination) {
        this.ensureClientPageInRange()
      }
    }
  },
  async created() {
    await Promise.all([
      this.getGroupList(),
      this.getAccountList()
    ])
    await this.getDatabaseList()
  },
  methods: {
    isValidationFailure(error) {
      return Boolean(error)
        && typeof error === 'object'
        && Object.values(error).some((value) => Array.isArray(value))
    },
    getDbIcon(type) {
      const iconMap = {
        mysql: require('@/assets/image/mysql.svg'),
        postgresql: require('@/assets/image/PostgreSQL.svg'),
        redis: require('@/assets/image/redis.svg'),
        mongodb: require('@/assets/image/mongodb.svg'),
        elasticsearch: require('@/assets/image/Elasticsearch.svg'),
        dameng: require('@/assets/image/dameng.svg'),
        gaussdb: require('@/assets/image/gaussdb.svg'),
        oracle: require('@/assets/image/oracle.svg')
      }
      return iconMap[getCmdbDatabaseTypeIconKey(type)] || iconMap.mysql
    },
    getDbName(type) {
      return getCmdbDatabaseTypeLabel(type)
    },
    getDbTagType(type) {
      return getCmdbDatabaseTypeTag(type)
    },
    async getAccountList() {
      try {
        const { data: res } = await configAPI.listAccountAuth({
          page: 1,
          pageSize: 100
        })
        if (res.code === 200) {
          this.accountList = normalizeAccountResponse(res.data)
          return
        }
        this.$message.error(res.message || '获取账号列表失败')
      } catch (error) {
        console.error('获取账号列表失败:', error)
        this.$message.error(`获取账号列表失败: ${error.message || '未知错误'}`)
      }
    },
    async getGroupList() {
      try {
        const { data: res } = await cmdbAPI.getAllCmdbGroups()
        if (res.code === 200) {
          this.groupList = Array.isArray(res.data) ? res.data : []
          return
        }
        this.$message.error(res.message || '获取分组列表失败')
      } catch (error) {
        console.error('获取分组列表失败:', error)
        this.$message.error(`获取分组列表失败: ${error.message || '未知错误'}`)
      }
    },
    ensureClientPageInRange() {
      const maxPage = Math.max(1, Math.ceil(this.filteredDatabaseRows.length / this.queryParams.pageSize))
      if (this.queryParams.pageNum > maxPage) {
        this.queryParams.pageNum = maxPage
      }
    },
    async fetchDefaultDatabasePage() {
      const requestedPage = this.queryParams.pageNum
      const { data: res } = await cmdbAPI.listDatabases({
        page: requestedPage,
        pageSize: this.queryParams.pageSize
      })

      if (res.code !== 200) {
        throw new Error(res.message || '获取数据库资产列表失败')
      }

      this.rawDatabaseList = normalizeListResponse(res.data)
      this.total = normalizeTotalResponse(res.data) ?? this.rawDatabaseList.length
      this.useClientPagination = false

      const maxPage = Math.max(1, Math.ceil(this.total / this.queryParams.pageSize))
      if (requestedPage > maxPage) {
        this.queryParams.pageNum = maxPage
        if (maxPage !== requestedPage) {
          return this.fetchDefaultDatabasePage()
        }
      }
    },
    async fetchSearchDatabaseResults() {
      let response
      const nameKeyword = trimString(this.queryParams.name)
      const selectedType = this.queryParams.type

      if (nameKeyword) {
        response = await cmdbAPI.getDatabasesByName(nameKeyword)
      } else if (selectedType) {
        response = await cmdbAPI.getDatabasesByType(selectedType)
      } else {
        return this.fetchDefaultDatabasePage()
      }

      const res = response.data
      if (res.code !== 200) {
        throw new Error(res.message || '查询数据库资产失败')
      }

      this.rawDatabaseList = normalizeListResponse(res.data)
      this.total = this.rawDatabaseList.length
      this.useClientPagination = true
      this.ensureClientPageInRange()
    },
    async getDatabaseList() {
      this.loading = true
      try {
        if (this.hasBackendSearchHint) {
          await this.fetchSearchDatabaseResults()
        } else {
          await this.fetchDefaultDatabasePage()
        }
      } catch (error) {
        console.error('获取数据库资产列表失败:', error)
        this.$message.error(`获取数据库资产列表失败: ${error.message || '未知错误'}`)
        this.rawDatabaseList = []
        this.total = 0
      } finally {
        this.loading = false
      }
    },
    async handleQuery() {
      this.queryParams.pageNum = 1
      await this.getDatabaseList()
    },
    async resetQuery() {
      this.queryParams = {
        pageNum: 1,
        pageSize: 10,
        name: '',
        address: '',
        type: undefined
      }
      await this.getDatabaseList()
    },
    async handleSizeChange(size) {
      this.queryParams.pageSize = size
      this.queryParams.pageNum = 1
      if (this.useClientPagination) {
        this.ensureClientPageInRange()
        return
      }
      await this.getDatabaseList()
    },
    async handleCurrentChange(page) {
      this.queryParams.pageNum = page
      if (this.useClientPagination) {
        this.ensureClientPageInRange()
        return
      }
      await this.getDatabaseList()
    },
    handleTypeChange(type) {
      if (!trimString(this.formData.platform)) {
        this.formData.platform = getCmdbDatabaseTypeLabel(type)
      }
    },
    showAddDialog() {
      this.dialogMode = 'create'
      this.formData = createDatabaseAssetFormModel()
      this.dialogVisible = true
      this.$nextTick(() => {
        this.$refs.formRef?.clearValidate()
      })
    },
    async showEditDialog(row) {
      try {
        const { data: res } = await cmdbAPI.getDatabaseAsset(row.id)
        if (res.code === 200) {
          this.dialogMode = 'edit'
          this.formData = normalizeDatabaseForm(res.data || row)
          this.dialogVisible = true
          this.$nextTick(() => {
            this.$refs.formRef?.clearValidate()
          })
          return
        }
        this.$message.error(res.message || '获取数据库资产详情失败')
      } catch (error) {
        console.error('获取数据库资产详情失败:', error)
        this.$message.error(`获取数据库资产详情失败: ${error.message || '未知错误'}`)
      }
    },
    handleDialogClose() {
      this.formData = createDatabaseAssetFormModel()
    },
    async submitForm() {
      try {
        await this.$refs.formRef.validate()

        const remark = trimString(this.formData.remark)
        const payload = {
          name: trimString(this.formData.name),
          address: trimString(this.formData.address),
          platform: trimString(this.formData.platform) || getCmdbDatabaseTypeLabel(this.formData.type),
          groupId: Number(this.formData.groupId),
          defaultDatabase: trimString(this.formData.defaultDatabase),
          protocolGroup: trimString(this.formData.protocolGroup) || 'default',
          accountId: Number(this.formData.accountId),
          tags: trimString(this.formData.tags),
          isActive: Boolean(this.formData.isActive),
          description: remark,
          remark,
          type: Number(this.formData.type)
        }

        let res
        if (this.dialogMode === 'edit' && this.formData.id) {
          res = await cmdbAPI.updateDatabase({
            ...payload,
            id: Number(this.formData.id)
          })
        } else {
          res = await cmdbAPI.createDatabase(payload)
        }

        if (res.data?.code === 200) {
          this.$message.success(this.dialogMode === 'edit' ? '数据库资产更新成功' : '数据库资产创建成功')
          this.dialogVisible = false
          await this.getDatabaseList()
          return
        }

        this.$message.error(res.data?.message || '数据库资产保存失败')
      } catch (error) {
        if (this.isValidationFailure(error)) {
          return
        }
        this.$message.error(`数据库资产保存失败: ${error.message || '未知错误'}`)
      }
    },
    async handleDeleteDatabase(row) {
      const confirmResult = await this.$confirm(
        `是否确认删除数据库资产 "${row.name}"?`,
        '提示',
        {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        }
      ).catch((error) => error)

      if (confirmResult !== 'confirm') {
        return
      }

      try {
        const { data: res } = await cmdbAPI.deleteDatabase(row.id)
        if (res.code === 200) {
          this.$message.success('数据库资产删除成功')
          await this.getDatabaseList()
          return
        }
        this.$message.error(res.message || '数据库资产删除失败')
      } catch (error) {
        console.error('数据库资产删除失败:', error)
        this.$message.error(`数据库资产删除失败: ${error.message || '未知错误'}`)
      }
    },
    goToDetails(row) {
      this.$router.push({
        path: '/cmdb/dbdetails',
        query: { id: row.id }
      })
    },
    goToWorkOrders(row) {
      this.$router.push({
        path: '/cmdb/sql-work-orders',
        query: { databaseId: row.id }
      })
    }
  }
}
</script>

<style scoped>
.cmdb-db-management {
  padding: 0;
  min-height: auto;
  background: transparent;
}

.db-card {
  background: rgba(12, 24, 41, 0.92);
  backdrop-filter: blur(18px);
  border-radius: 24px;
  box-shadow: var(--shadow-card);
  border: 1px solid var(--border-subtle);
}

.db-management-container {
  display: flex;
}

.db-table-section {
  flex: 1;
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

.button-label {
  margin-left: 4px;
}

.table-section {
  margin-bottom: 20px;
}

.db-table {
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
}

.db-table :deep(.el-table__header) {
  background: linear-gradient(135deg, #667eea, #764ba2);
}

.db-table :deep(.el-table__header th) {
  background: transparent !important;
  color: #2c3e50 !important;
  font-weight: 700 !important;
  border-bottom: none;
  text-shadow: 0 1px 2px rgba(255, 255, 255, 0.8);
}

.db-table :deep(.el-table__row) {
  transition: all 0.3s ease;
}

.db-table :deep(.el-table__row:hover) {
  background-color: rgba(103, 126, 234, 0.05) !important;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.db-name-cell {
  display: flex;
  align-items: center;
  gap: 10px;
}

.db-name-content {
  display: flex;
  flex-direction: column;
  min-width: 0;
}

.db-name-subtitle {
  margin-top: 2px;
  color: #94a3b8;
  font-size: 12px;
}

.db-icon {
  width: 20px;
  height: 20px;
  flex-shrink: 0;
}

.table-operation {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  flex-wrap: wrap;
  gap: 8px;
}

.pagination-section {
  text-align: right;
}

.asset-form-section {
  margin-bottom: 14px;
  padding: 12px 12px 4px;
  border-radius: 12px;
  border: 1px solid rgba(148, 163, 184, 0.18);
  background: rgba(15, 23, 42, 0.28);
}

.asset-form-section__title {
  margin-bottom: 10px;
  color: #e2e8f0;
  font-size: 13px;
  font-weight: 600;
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
.el-select :deep(.el-input__wrapper),
.el-textarea :deep(.el-textarea__inner) {
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid var(--border-subtle);
  border-radius: 12px;
  box-shadow: none;
}

.search-section .el-form-item {
  margin-bottom: 0;
  margin-right: 16px;
}

@media (max-width: 960px) {
  .table-operation {
    gap: 6px;
  }
}
</style>
