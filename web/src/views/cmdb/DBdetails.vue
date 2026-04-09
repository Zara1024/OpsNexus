<template>
  <div class="db-details-container">
    <el-card class="asset-card">
      <template #header>
        <div class="card-header">
          <div>
            <div class="card-title">数据库资产详情</div>
            <div class="card-subtitle">查看资产信息，并继续进入数据库操作与 SQL 工单流程。</div>
          </div>
          <div class="header-actions">
            <el-button size="small" @click="goWorkOrderCenter">SQL工单中心</el-button>
            <el-button type="primary" size="small" @click="getDBDetails">
              <el-icon><Refresh /></el-icon>
              <span style="margin-left: 4px">刷新资产</span>
            </el-button>
          </div>
        </div>
      </template>

      <el-descriptions :column="2" border>
        <el-descriptions-item label="名称">{{ dbInfo.name || '-' }}</el-descriptions-item>
        <el-descriptions-item label="地址">{{ dbInfo.address || '-' }}</el-descriptions-item>
        <el-descriptions-item label="平台">{{ dbInfo.platform || '-' }}</el-descriptions-item>
        <el-descriptions-item label="类型">
          <el-tag :type="getDatabaseTypeTag(dbInfo.type)">
            {{ getDatabaseTypeLabel(dbInfo.type) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="账号">{{ dbInfo.accountAlias || '-' }}</el-descriptions-item>
        <el-descriptions-item label="分组/节点">{{ dbInfo.groupName || '-' }}</el-descriptions-item>
        <el-descriptions-item label="默认数据库">{{ dbInfo.defaultDatabase || '-' }}</el-descriptions-item>
        <el-descriptions-item label="标签">{{ dbInfo.tags || '-' }}</el-descriptions-item>
        <el-descriptions-item label="备注">{{ dbInfo.remark || '-' }}</el-descriptions-item>
        <el-descriptions-item label="更新时间">{{ formatDate(dbInfo.updatedAt) }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ formatDate(dbInfo.createdAt) }}</el-descriptions-item>
      </el-descriptions>
    </el-card>

    <div class="sql-execute-container">
      <el-card>
        <template #header>
          <div class="card-header simple-header">
            <span class="card-title">数据库操作</span>
          </div>
        </template>

        <div class="sql-form-container">
          <el-form :model="sqlForm" label-width="100px">
            <el-form-item label="资产名称">
              <el-input v-model="dbInfo.name" disabled />
            </el-form-item>
            <el-form-item :label="databaseTargetLabel">
              <el-select
                v-model="sqlForm.databaseName"
                :placeholder="databaseTargetPlaceholder"
                filterable
                style="width: 100%"
                @focus="getDatabaseList"
              >
                <el-option
                  v-for="item in databaseList"
                  :key="item"
                  :label="item"
                  :value="item"
                />
              </el-select>
              <div v-if="dbConnectionInfo.host" class="connection-info">
                实例地址: {{ dbConnectionInfo.host }}<span v-if="dbConnectionInfo.port">:{{ dbConnectionInfo.port }}</span>
              </div>
            </el-form-item>
            <el-form-item :label="operationTypeLabel">
              <el-select
                v-model="sqlForm.sqlType"
                :placeholder="operationTypePlaceholder"
                style="width: 100%"
                @change="updateSQLPlaceholder"
              >
                <el-option
                  v-for="item in currentOperationTypeOptions"
                  :key="item.value"
                  :label="item.label"
                  :value="item.value"
                />
              </el-select>
            </el-form-item>
            <el-form-item :label="statementLabel">
              <el-input
                v-model="sqlForm.sql"
                type="textarea"
                :rows="6"
                :placeholder="sqlPlaceholder"
              />
            </el-form-item>
            <el-form-item label="工单标题">
              <el-input v-model="sqlForm.workOrderTitle" :placeholder="workOrderTitlePlaceholder" />
            </el-form-item>
            <el-form-item label="变更原因">
              <el-input
                v-model="sqlForm.workOrderReason"
                type="textarea"
                :rows="3"
                placeholder="说明本次变更的目的、影响范围和回退方案。"
              />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="executeSQL">{{ executeButtonText }}</el-button>
              <el-button type="warning" :disabled="isDirectOnlyOperation" @click="submitWorkOrder">
                提交工单
              </el-button>
              <el-button @click="goWorkOrderCenter">工单中心</el-button>
            </el-form-item>
          </el-form>
        </div>
      </el-card>

      <el-card>
        <template #header>
          <div class="card-header simple-header">
            <span class="card-title">执行结果</span>
          </div>
        </template>
        <div class="result-container">
          <pre>{{ executionResult }}</pre>
        </div>
      </el-card>
    </div>

    <el-card class="work-order-card">
      <template #header>
        <div class="card-header">
          <span class="card-title">最近 SQL 工单</span>
          <el-button size="small" type="primary" link @click="goWorkOrderCenter">查看全部</el-button>
        </div>
      </template>
      <el-table :data="recentWorkOrders" size="small" stripe empty-text="当前数据库资产暂无 SQL 工单">
        <el-table-column label="工单号" prop="orderNo" min-width="160" />
        <el-table-column label="标题" prop="title" min-width="180" show-overflow-tooltip />
        <el-table-column label="类型" prop="operationType" width="120" />
        <el-table-column label="风险" width="90">
          <template #default="{ row }">
            <el-tag :type="riskTagType(row.riskLevel)" effect="dark">{{ riskLabel(row.riskLevel) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="120">
          <template #default="{ row }">
            <el-tag :type="workOrderStatusTagType(row.status)" effect="dark">{{ workOrderStatusLabel(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="创建时间" prop="createTime" min-width="160" />
      </el-table>
    </el-card>
  </div>
</template>

<script>
import { computed, onMounted, ref } from 'vue'
import { Refresh } from '@element-plus/icons-vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import cmdbAPI from '@/api/cmdb'
import {
  getCmdbDatabasePlatformLabel,
  getCmdbDatabaseTypeLabel,
  getCmdbDatabaseTypeTag
} from '@/utils/cmdbPresentation.mjs'

const sqlTypeOptions = [
  { label: '查询', value: 'select' },
  { label: '插入', value: 'insert' },
  { label: '更新', value: 'update' },
  { label: '删除', value: 'delete' },
  { label: '原生 SQL', value: 'raw' }
]

const redisCommandModeOptions = [
  { label: '只读命令', value: 'read' },
  { label: '写入命令', value: 'write' },
  { label: '原生命令', value: 'raw' }
]

export default {
  name: 'DBdetails',
  components: {
    Refresh
  },
  setup() {
    const route = useRoute()
    const router = useRouter()

    const dbInfo = ref({
      name: '',
      address: '',
      platform: '',
      type: 1,
      accountAlias: '',
      groupName: '',
      defaultDatabase: '',
      tags: '',
      remark: '',
      createdAt: '',
      updatedAt: ''
    })
    const sqlForm = ref({
      sqlType: 'select',
      sql: '',
      databaseName: '',
      workOrderTitle: '',
      workOrderReason: ''
    })
    const executionResult = ref('请选择数据库并输入 SQL/命令后执行，或先提交 SQL 工单。')
    const databaseList = ref([])
    const recentWorkOrders = ref([])
    const dbConnectionInfo = ref({
      host: '',
      port: ''
    })
    const sqlPlaceholder = ref('示例: SELECT * FROM table_name WHERE id = 1;')

    const getDatabaseId = () => Number(route.params.id || route.query.id)
    const databaseType = computed(() => Number(dbInfo.value.type || 1))
    const isMySQLDatabaseType = computed(() => databaseType.value === 1)
    const isPostgreSQLDatabaseType = computed(() => databaseType.value === 2)
    const isRedisDatabaseType = computed(() => databaseType.value === 3)

    const currentOperationTypeOptions = computed(() => (
      isRedisDatabaseType.value ? redisCommandModeOptions : sqlTypeOptions
    ))

    const databaseTargetLabel = computed(() => (
      isRedisDatabaseType.value ? '选择逻辑库' : '选择数据库'
    ))

    const databaseTargetPlaceholder = computed(() => (
      isRedisDatabaseType.value ? '请选择 Redis 逻辑库' : '请选择要操作的数据库'
    ))

    const operationTypeLabel = computed(() => (
      isRedisDatabaseType.value ? '命令模式' : 'SQL 类型'
    ))

    const operationTypePlaceholder = computed(() => (
      isRedisDatabaseType.value ? '请选择命令模式' : '请选择 SQL 类型'
    ))

    const statementLabel = computed(() => (
      isRedisDatabaseType.value ? 'Redis 命令' : 'SQL 语句'
    ))

    const executeButtonText = computed(() => (
      isRedisDatabaseType.value ? '执行命令' : '执行'
    ))

    const workOrderTitlePlaceholder = computed(() => (
      isRedisDatabaseType.value ? '例如：修复缓存键过期策略' : '例如：修复订单库索引'
    ))

    const isDirectOnlyOperation = computed(() => (
      (isRedisDatabaseType.value && sqlForm.value.sqlType === 'read') ||
      (!isRedisDatabaseType.value && sqlForm.value.sqlType === 'select')
    ))

    const getDatabaseTypeLabel = (type) => getCmdbDatabaseTypeLabel(type)
    const getDatabaseTypeTag = (type) => getCmdbDatabaseTypeTag(type)

    const normalizeDatabaseList = (data) => {
      if (Array.isArray(data?.databases) && data.databases.length) {
        return data.databases
      }
      if (isRedisDatabaseType.value) {
        return ['0']
      }
      return []
    }

    const updateSQLPlaceholder = (type) => {
      if (isRedisDatabaseType.value) {
        switch (type) {
          case 'write':
            sqlPlaceholder.value = '示例: SET app:feature:toggle true'
            break
          case 'raw':
            sqlPlaceholder.value = '示例: HGETALL app:session:1001'
            break
          default:
            sqlPlaceholder.value = '示例: GET app:feature:toggle'
        }
        return
      }

      switch (type) {
        case 'insert':
          sqlPlaceholder.value = "示例: INSERT INTO table_name (field1) VALUES ('value');"
          break
        case 'update':
          sqlPlaceholder.value = "示例: UPDATE table_name SET field1 = 'value' WHERE id = 1;"
          break
        case 'delete':
          sqlPlaceholder.value = '示例: DELETE FROM table_name WHERE id = 1;'
          break
        case 'raw':
          sqlPlaceholder.value = isPostgreSQLDatabaseType.value
            ? '示例: ALTER TABLE public.orders ADD COLUMN note TEXT;'
            : '示例: SHOW INDEX FROM table_name;'
          break
        default:
          sqlPlaceholder.value = isPostgreSQLDatabaseType.value
            ? '示例: SELECT datname FROM pg_database ORDER BY datname;'
            : '示例: SELECT * FROM table_name WHERE id = 1;'
      }
    }

    const getDatabaseList = async () => {
      try {
        const databaseId = getDatabaseId()
        if (!databaseId) {
          ElMessage.error('数据库资产 ID 不能为空')
          return
        }

        const { data: res } = await cmdbAPI.executeDatabase({ databaseId })
        if (res.code !== 200) {
          throw new Error(res.message || '获取数据库列表失败')
        }

        databaseList.value = normalizeDatabaseList(res.data || {})
        dbConnectionInfo.value = {
          host: res.data?.host || '',
          port: res.data?.port || ''
        }

        if (!sqlForm.value.databaseName) {
          sqlForm.value.databaseName = dbInfo.value.defaultDatabase || databaseList.value[0] || ''
        }
      } catch (error) {
        console.error('获取数据库列表失败:', error)
        ElMessage.error(error.message || '获取数据库列表失败')
      }
    }

    const getDBDetails = async () => {
      try {
        const databaseId = getDatabaseId()
        if (!databaseId) {
          ElMessage.error('数据库资产 ID 不能为空')
          return
        }

        const { data: res } = await cmdbAPI.getDatabaseAsset(databaseId)
        if (res.code !== 200) {
          throw new Error(res.message || '获取数据库资产详情失败')
        }

        const data = res.data || {}
        dbInfo.value = {
          name: data.name || '',
          address: data.address || '',
          platform: getCmdbDatabasePlatformLabel(data.type, data.platform),
          type: Number(data.type || 1),
          accountAlias: data.accountAlias || data.accountName || (data.accountId ? `账号ID: ${data.accountId}` : ''),
          groupName: data.groupName || (data.groupId ? `分组ID: ${data.groupId}` : ''),
          defaultDatabase: data.defaultDatabase || '',
          tags: data.tags || '',
          remark: data.remark || data.description || '',
          createdAt: data.createdAt || '',
          updatedAt: data.updatedAt || ''
        }

        sqlForm.value.sqlType = isRedisDatabaseType.value ? 'read' : 'select'
        updateSQLPlaceholder(sqlForm.value.sqlType)

        if (!sqlForm.value.databaseName && dbInfo.value.defaultDatabase) {
          sqlForm.value.databaseName = dbInfo.value.defaultDatabase
        }
      } catch (error) {
        console.error('获取数据库资产详情失败:', error)
        ElMessage.error(error.message || '获取数据库资产详情失败')
      }
    }

    const buildExecutionRequest = () => ({
      databaseId: getDatabaseId(),
      databaseName: sqlForm.value.databaseName,
      sql: sqlForm.value.sql
    })

    const executeSQL = async () => {
      if (!sqlForm.value.sql) {
        ElMessage.warning(isRedisDatabaseType.value ? '请输入 Redis 命令' : '请输入 SQL 语句')
        return
      }
      if (!sqlForm.value.databaseName) {
        ElMessage.warning(isRedisDatabaseType.value ? '请选择 Redis 逻辑库' : '请选择数据库')
        return
      }

      const requestData = buildExecutionRequest()

      try {
        let res
        if (isRedisDatabaseType.value) {
          res = await cmdbAPI.executeRawSQL(requestData)
        } else {
          switch (sqlForm.value.sqlType) {
            case 'select':
              res = await cmdbAPI.executeSelectSQL(requestData)
              break
            case 'insert':
              res = await cmdbAPI.executeInsertSQL(requestData)
              break
            case 'update':
              res = await cmdbAPI.executeUpdateSQL(requestData)
              break
            case 'delete':
              res = await cmdbAPI.executeDeleteSQL(requestData)
              break
            case 'raw':
              res = await cmdbAPI.executeRawSQL(requestData)
              break
            default:
              ElMessage.warning('请选择有效的 SQL 类型')
              return
          }
        }

        if (res.data?.code === 200) {
          executionResult.value = JSON.stringify(res.data.data, null, 2)
          ElMessage.success(isRedisDatabaseType.value ? '命令执行成功' : 'SQL 执行成功')
        } else {
          executionResult.value = JSON.stringify(res.data, null, 2)
          ElMessage.error(res.data?.message || (isRedisDatabaseType.value ? '命令执行失败' : 'SQL 执行失败'))
        }
      } catch (error) {
        console.error('数据库操作执行失败:', error)
        executionResult.value = error.message
        ElMessage.error(`${isRedisDatabaseType.value ? '命令执行失败' : 'SQL 执行失败'}: ${error.message || '未知错误'}`)
      }
    }

    const loadRecentWorkOrders = async () => {
      try {
        const databaseId = getDatabaseId()
        if (!databaseId) {
          return
        }

        const { data: res } = await cmdbAPI.getSQLWorkOrderList({
          page: 1,
          pageSize: 5,
          databaseId
        })

        if (res.code === 200) {
          recentWorkOrders.value = res.data?.list || []
        }
      } catch (error) {
        console.error('获取最近 SQL 工单失败:', error)
      }
    }

    const submitWorkOrder = async () => {
      if (!sqlForm.value.sql) {
        ElMessage.warning(isRedisDatabaseType.value ? '请输入 Redis 命令' : '请输入 SQL 语句')
        return
      }
      if (!sqlForm.value.databaseName) {
        ElMessage.warning(isRedisDatabaseType.value ? '请选择 Redis 逻辑库' : '请选择数据库')
        return
      }
      if (isDirectOnlyOperation.value) {
        ElMessage.warning(
          isRedisDatabaseType.value
            ? '只读 Redis 命令建议直接执行，工单主要用于变更类命令'
            : 'SELECT 建议直接执行查询，SQL 工单主要用于变更类 SQL'
        )
        return
      }

      try {
        const { data: res } = await cmdbAPI.createSQLWorkOrder({
          databaseId: getDatabaseId(),
          databaseName: sqlForm.value.databaseName,
          title: sqlForm.value.workOrderTitle,
          reason: sqlForm.value.workOrderReason,
          sql: sqlForm.value.sql
        })

        if (res.code !== 200) {
          throw new Error(res.message || '提交 SQL 工单失败')
        }

        executionResult.value = JSON.stringify(res.data, null, 2)
        ElMessage.success(`SQL 工单提交成功: ${res.data?.orderNo || ''}`)
        sqlForm.value.workOrderTitle = ''
        sqlForm.value.workOrderReason = ''
        await loadRecentWorkOrders()
      } catch (error) {
        console.error('提交 SQL 工单失败:', error)
        ElMessage.error(error.message || '提交 SQL 工单失败')
      }
    }

    const goWorkOrderCenter = () => {
      const databaseId = getDatabaseId()
      router.push({
        path: '/cmdb/sql-work-orders',
        query: databaseId ? { databaseId } : {}
      })
    }

    const formatDate = (value) => {
      if (!value) {
        return '-'
      }
      const date = new Date(value)
      return date.toLocaleString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit',
        second: '2-digit'
      })
    }

    const riskLabel = (level) => ({ 0: '低', 1: '中', 2: '高' }[Number(level)] || '未知')
    const riskTagType = (level) => (Number(level) >= 2 ? 'danger' : (Number(level) === 1 ? 'warning' : 'success'))
    const workOrderStatusLabel = (status) => ({
      1: '待审批',
      2: '已批准',
      3: '已驳回',
      4: '执行中',
      5: '执行成功',
      6: '执行失败'
    }[status] || '未知')
    const workOrderStatusTagType = (status) => ({
      1: 'warning',
      2: 'primary',
      3: 'info',
      4: 'warning',
      5: 'success',
      6: 'danger'
    }[status] || 'info')

    onMounted(async () => {
      await getDBDetails()
      await getDatabaseList()
      await loadRecentWorkOrders()
    })

    return {
      currentOperationTypeOptions,
      databaseList,
      databaseTargetLabel,
      databaseTargetPlaceholder,
      dbConnectionInfo,
      dbInfo,
      executeButtonText,
      executeSQL,
      executionResult,
      formatDate,
      getDBDetails,
      getDatabaseList,
      getDatabaseTypeLabel,
      getDatabaseTypeTag,
      goWorkOrderCenter,
      isDirectOnlyOperation,
      isMySQLDatabaseType,
      isPostgreSQLDatabaseType,
      isRedisDatabaseType,
      operationTypeLabel,
      operationTypePlaceholder,
      redisCommandModeOptions,
      recentWorkOrders,
      riskLabel,
      riskTagType,
      sqlForm,
      sqlPlaceholder,
      statementLabel,
      submitWorkOrder,
      updateSQLPlaceholder,
      workOrderStatusLabel,
      workOrderStatusTagType,
      workOrderTitlePlaceholder
    }
  }
}
</script>

<style scoped>
.db-details-container {
  padding: 20px;
}

.asset-card,
.work-order-card {
  border-radius: 18px;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}

.simple-header {
  justify-content: flex-start;
}

.card-title {
  font-size: 18px;
  font-weight: 700;
}

.card-subtitle {
  margin-top: 4px;
  color: #64748b;
  font-size: 13px;
}

.header-actions {
  display: flex;
  gap: 10px;
}

.sql-execute-container {
  margin-top: 20px;
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
}

.sql-form-container {
  padding: 10px 0;
}

.connection-info {
  margin-top: 8px;
  color: #64748b;
  font-size: 12px;
}

.result-container {
  padding: 16px;
  min-height: 320px;
  max-height: 520px;
  overflow: auto;
  background-color: #020617;
  border-radius: 12px;
  color: #e2e8f0;
  font-family: Consolas, Monaco, monospace;
}

.work-order-card {
  margin-top: 20px;
}

pre {
  margin: 0;
  white-space: pre-wrap;
  word-break: break-word;
}

@media (max-width: 960px) {
  .sql-execute-container {
    grid-template-columns: 1fr;
  }

  .card-header {
    flex-direction: column;
    align-items: flex-start;
  }
}
</style>
