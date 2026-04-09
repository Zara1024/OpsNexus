<template>
  <TablePage
    class="loginlog-table-page"
    section-title="登录记录"
    section-subtitle="按用户、状态和时间回溯登录行为，适合作为审计与故障排查入口。"
  >
    <template #header>
      <PageHeader
        eyebrow="Access Audit"
        title="登录日志"
        subtitle="统一查看登录来源、失败原因与访问时间，快速识别异常账号行为。"
      >
        <template #intro>
          <PageIntro
            title="排查建议"
            text="优先关注失败率较高的账号、异常时间段和陌生来源 IP，再决定是否清理或归档。"
          />
        </template>
      </PageHeader>
    </template>

    <template #toolbar>
      <PageToolbar>
        <el-form :inline="true" class="filter-cluster">
          <el-form-item class="filter-field" label="用户名称">
            <el-input v-model="queryParams.username" placeholder="请输入用户名称" clearable @keyup.enter="handleQuery" />
          </el-form-item>
          <el-form-item class="filter-field" label="登录状态">
            <el-select v-model="queryParams.loginStatus" placeholder="请选择状态" clearable>
              <el-option v-for="item in loginStatusList" :key="item.value" :label="item.label" :value="item.value" />
            </el-select>
          </el-form-item>
          <el-form-item class="filter-field" label="开始时间">
            <el-date-picker v-model="queryParams.beginTime" type="date" value-format="yyyy-MM-dd" clearable placeholder="请选择开始时间" />
          </el-form-item>
          <el-form-item class="filter-field" label="结束时间">
            <el-date-picker v-model="queryParams.endTime" type="date" value-format="yyyy-MM-dd" clearable placeholder="请选择结束时间" />
          </el-form-item>
        </el-form>
        <template #actions>
          <el-button type="primary" icon="Search" @click="handleQuery">搜索</el-button>
          <el-button icon="Refresh" @click="resetQuery">重置</el-button>
          <el-tooltip :disabled="!multiple" content="请先勾选要删除的登录日志" placement="top">
            <span class="batch-delete-wrapper">
              <el-button type="danger" plain icon="Delete" :disabled="multiple" @click="batchHandleDelete" v-authority="['monitor:loginLog:delete']">
                批量删除
              </el-button>
            </span>
          </el-tooltip>
          <el-button type="danger" plain icon="Delete" @click="handleClean" v-authority="['monitor:loginLog:clean']">
            清空日志
          </el-button>
        </template>
      </PageToolbar>
    </template>

    <template #stats>
      <StatStrip :items="statItems" />
    </template>

    <el-table
      v-loading="Loading"
      :data="sysLoginInfoList"
      stripe
      class="loginlog-table"
      @selection-change="handleSelectionChange"
    >
      <el-table-column type="selection" width="52" />
      <el-table-column label="用户账号" prop="username" min-width="120" />
      <el-table-column label="登录 IP" prop="ipAddress" min-width="140" />
      <el-table-column label="登录地点" prop="loginLocation" min-width="130" show-overflow-tooltip />
      <el-table-column label="登录状态" prop="loginStatus" width="110">
        <template #default="scope">
          <el-tag :type="scope.row.loginStatus === 1 ? 'success' : 'danger'">
            {{ scope.row.loginStatus === 1 ? '成功' : '失败' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="提示消息" prop="message" min-width="140" show-overflow-tooltip />
      <el-table-column label="访问时间" prop="loginTime" min-width="160" />
      <el-table-column label="操作" width="120" fixed="right">
        <template #default="scope">
          <div class="operation-buttons">
            <el-button
              type="danger"
              icon="Delete"
              size="small"
              circle
              @click="handleDelete(scope.row.id)"
              v-authority="['monitor:loginLog:delete']"
            />
          </div>
        </template>
      </el-table-column>
    </el-table>

    <template #footer>
      <el-pagination
        :current-page="queryParams.pageNum"
        :page-sizes="[10, 50, 100, 500, 1000]"
        :page-size="queryParams.pageSize"
        layout="total, sizes, prev, pager, next, jumper"
        :total="total"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </template>
  </TablePage>
</template>

<script>
import PageHeader from '@/components/platform/PageHeader.vue'
import PageIntro from '@/components/platform/PageIntro.vue'
import PageToolbar from '@/components/platform/PageToolbar.vue'
import StatStrip from '@/components/platform/StatStrip.vue'
import TablePage from '@/components/platform/TablePage.vue'

export default {
  name: 'LoginLog',
  components: {
    PageHeader,
    PageIntro,
    PageToolbar,
    StatStrip,
    TablePage
  },
  data() {
    return {
      queryParams: {},
      loginStatusList: [
        { value: '1', label: '成功' },
        { value: '2', label: '失败' }
      ],
      sysLoginInfoList: [],
      Loading: true,
      ids: [],
      single: true,
      multiple: true,
      total: 0
    }
  },
  computed: {
    statItems() {
      const successCount = this.sysLoginInfoList.filter(item => Number(item.loginStatus) === 1).length
      const failureCount = this.sysLoginInfoList.filter(item => Number(item.loginStatus) === 2).length
      return [
        { label: '日志总量', value: this.total, hint: '当前条件下的记录总数', tone: 'primary' },
        { label: '成功登录', value: successCount, hint: '当前页成功记录', tone: 'success' },
        { label: '失败登录', value: failureCount, hint: '当前页失败记录', tone: 'danger' },
        { label: '已选记录', value: this.ids.length, hint: '可执行批量删除', tone: 'warning' }
      ]
    }
  },
  methods: {
    handleSelectionChange(selection) {
      this.ids = selection.map(item => item.id)
      this.single = selection.length !== 1
      this.multiple = !selection.length
    },
    async getSysLoginInfoList() {
      this.Loading = true
      const { data: res } = await this.$api.querySysLoginInfoList(this.queryParams)
      if (res.code !== 200) {
        this.$message.error(res.message)
      } else {
        this.sysLoginInfoList = res.data.list
        this.total = res.data.total
        this.Loading = false
      }
    },
    handleQuery() {
      this.getSysLoginInfoList()
    },
    resetQuery() {
      this.queryParams = {}
      this.getSysLoginInfoList()
      this.$message.success('重置成功')
    },
    handleSizeChange(newSize) {
      this.queryParams.pageSize = newSize
      this.getSysLoginInfoList()
    },
    handleCurrentChange(newPage) {
      this.queryParams.pageNum = newPage
      this.getSysLoginInfoList()
    },
    async handleClean() {
      const confirmResult = await this.$confirm('是否清空登录日志？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).catch(err => err)
      if (confirmResult !== 'confirm') {
        return this.$message.info('已取消')
      }
      const { data: res } = await this.$api.cleanSysLoginInfo()
      if (res.code !== 200) {
        this.$message.error(res.message)
      } else {
        this.$message.success('清空成功')
        await this.getSysLoginInfoList()
      }
    },
    async handleDelete(id) {
      const confirmResult = await this.$confirm(`是否确认删除登录日志编号为 "${id}" 的数据项？`, '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).catch(err => err)
      if (confirmResult !== 'confirm') {
        return this.$message.info('已取消删除')
      }
      const { data: res } = await this.$api.deleteSysLoginInfo(id)
      if (res.code !== 200) {
        this.$message.error(res.message)
      } else {
        this.$message.success('删除成功')
        await this.getSysLoginInfoList()
      }
    },
    async batchHandleDelete() {
      const loginInfoIds = this.ids
      const confirmResult = await this.$confirm(`是否确认删除登录日志编号为 "${loginInfoIds}" 的数据项？`, '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).catch(err => err)
      if (confirmResult !== 'confirm') {
        return this.$message.info('已取消删除')
      }
      const { data: res } = await this.$api.batchDeleteSysLoginInfo(loginInfoIds)
      if (res.code !== 200) {
        this.$message.error(res.message)
      } else {
        this.$message.success('删除成功')
        await this.getSysLoginInfoList()
      }
    }
  },
  created() {
    this.getSysLoginInfoList()
  }
}
</script>

<style scoped>
.loginlog-table-page :deep(.el-form) {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

.loginlog-table-page :deep(.el-form-item) {
  margin-bottom: 0;
}
.batch-delete-wrapper { display: inline-flex; }
</style>
