<template>
  <TablePage
    class="accountauth-table-page"
    section-title="通用凭据列表"
    section-subtitle="统一管理数据库、CI/CD、监控和通用服务账号，支持快速查看与维护。"
  >
    <template #header>
      <PageHeader
        eyebrow="Credential Center"
        title="账号认证管理"
        subtitle="收敛数据库、Jenkins 与通用账号，统一提供检索、解密和生命周期治理。"
      >
        <template #actions>
          <el-button type="primary" icon="Plus" @click="showAddDialog" v-authority="['config:common:add']">
            新增账号
          </el-button>
        </template>
        <template #intro>
          <PageIntro
            title="运维建议"
            text="优先使用别名、服务地址和账号类型做归类，避免相同目标系统散落成多条重复凭据。"
          />
        </template>
      </PageHeader>
    </template>

    <template #toolbar>
      <PageToolbar>
        <el-form :inline="true" class="filter-cluster">
          <el-form-item class="filter-field" label="账号别名">
            <el-input v-model="queryParams.alias" placeholder="请输入账号别名" clearable @keyup.enter="handleQuery" />
          </el-form-item>
          <el-form-item class="filter-field" label="账号类型">
            <el-select v-model="queryParams.type" placeholder="请选择账号类型" clearable>
              <el-option label="Mysql" :value="1" />
              <el-option label="Postgre" :value="2" />
              <el-option label="Redis" :value="3" />
              <el-option label="Jenkins" :value="4" />
              <el-option label="Zabbix" :value="5" />
              <el-option label="通用账号" :value="6" />
            </el-select>
          </el-form-item>
        </el-form>
        <template #actions>
          <el-button type="primary" icon="Search" @click="handleQuery">搜索</el-button>
          <el-button icon="Refresh" @click="resetQuery">重置</el-button>
        </template>
      </PageToolbar>
    </template>

    <template #stats>
      <StatStrip :items="statItems" />
    </template>

    <el-table stripe v-loading="loading" :data="accountList" class="accountauth-table">
      <el-table-column label="账号别名" prop="alias" min-width="180">
        <template #default="scope">
          <div class="credential-cell">
            <img src="@/assets/image/账号.svg" class="credential-icon" alt="账号" />
            <span>{{ scope.row.alias }}</span>
          </div>
        </template>
      </el-table-column>
      <el-table-column label="服务地址" prop="host" min-width="150">
        <template #default="scope">
          <div class="credential-cell">
            <img src="@/assets/image/url.svg" class="credential-icon" alt="地址" />
            <span>{{ scope.row.host }}</span>
          </div>
        </template>
      </el-table-column>
      <el-table-column label="端口" prop="port" width="84">
        <template #default="scope">
          <div class="credential-cell">
            <img src="@/assets/image/端口.svg" class="credential-icon" alt="端口" />
            <span>{{ scope.row.port }}</span>
          </div>
        </template>
      </el-table-column>
      <el-table-column label="用户名" prop="name" min-width="120">
        <template #default="scope">
          <div class="credential-cell">
            <img src="@/assets/image/ren.svg" class="credential-icon" alt="用户" />
            <span>{{ scope.row.name }}</span>
          </div>
        </template>
      </el-table-column>
      <el-table-column label="账号类型" width="130">
        <template #default="scope">
          <el-tag :type="tagType(scope.row.type)">{{ typeText(scope.row.type) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="更新时间" prop="updatedAt" min-width="150" />
      <el-table-column label="操作" width="180" fixed="right">
        <template #default="scope">
          <div class="operation-buttons">
            <el-button size="small" v-authority="['config:common:edit']" type="primary" icon="Edit" circle @click="showEditDialog(scope.row)" />
            <el-button size="small" v-authority="['config:common:delete']" type="danger" icon="Delete" circle @click="handleDelete(scope.row)" />
            <el-button size="small" v-authority="['config:common:decrypt']" type="warning" icon="Key" circle @click="handleDecrypt(scope.row)" />
          </div>
        </template>
      </el-table-column>
    </el-table>

    <template #footer>
      <el-pagination
        :current-page="queryParams.pageNum"
        :page-sizes="[10, 20, 50, 100]"
        :page-size="queryParams.pageSize"
        layout="total, sizes, prev, pager, next, jumper"
        :total="total"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </template>

    <el-dialog :title="dialogTitle" v-model="dialogVisible" width="560px">
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="100px">
        <el-form-item label="账号别名" prop="alias">
          <el-input v-model="formData.alias" placeholder="请输入账号别名" />
        </el-form-item>
        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="服务地址" prop="host">
              <el-input v-model="formData.host" placeholder="192.168.1.1:3306 (无需协议)" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="服务类型" prop="type">
              <el-select v-model="formData.type" placeholder="请选择服务类型">
                <el-option label="Mysql" :value="1" />
                <el-option label="Postgre" :value="2" />
                <el-option label="Redis" :value="3" />
                <el-option label="Jenkins" :value="4" />
                <el-option label="Zabbix" :value="5" />
                <el-option label="通用账号" :value="6" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="用户名" prop="name">
              <el-input v-model="formData.name" placeholder="请输入用户名" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="密码" prop="password">
              <el-input v-model="formData.password" show-password placeholder="请输入密码" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="备注" prop="remark">
          <el-input v-model="formData.remark" type="textarea" :rows="3" placeholder="请输入备注信息" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitForm">确定</el-button>
      </template>
    </el-dialog>
  </TablePage>
</template>

<script>
import API from '@/api/config'
import PageHeader from '@/components/platform/PageHeader.vue'
import PageIntro from '@/components/platform/PageIntro.vue'
import PageToolbar from '@/components/platform/PageToolbar.vue'
import StatStrip from '@/components/platform/StatStrip.vue'
import TablePage from '@/components/platform/TablePage.vue'
import { formatAccountHostPortInput, parseAccountHostPortInput } from '@/utils/accountAuthEndpoint.mjs'

export default {
  name: 'AccountAuth',
  components: {
    PageHeader,
    PageIntro,
    PageToolbar,
    StatStrip,
    TablePage
  },
  data() {
    const validateHost = (_rule, value, callback) => {
      const parsedEndpoint = parseAccountHostPortInput(value)
      if (!parsedEndpoint.valid) {
        callback(new Error('请输入服务地址，格式为 IP或域名:端口'))
        return
      }
      callback()
    }

    return {
      queryParams: {
        alias: '',
        type: undefined,
        pageNum: 1,
        pageSize: 10
      },
      loading: false,
      accountList: [],
      total: 0,
      dialogVisible: false,
      dialogTitle: '',
      formData: {
        id: '',
        alias: '',
        host: '',
        name: '',
        password: '',
        type: undefined,
        remark: ''
      },
      formRules: {
        alias: [{ required: true, message: '请输入账号别名', trigger: 'blur' }],
        host: [
          { required: true, message: '请输入服务地址(格式: IP或域名:端口)', trigger: 'blur' },
          { validator: validateHost, trigger: ['blur', 'change'] }
        ],
        name: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
        password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
        type: [{ required: true, message: '请选择服务类型', trigger: 'change' }]
      }
    }
  },
  computed: {
    statItems() {
      const mysqlCount = this.accountList.filter(item => Number(item.type) === 1).length
      const jenkinsCount = this.accountList.filter(item => Number(item.type) === 4).length
      const genericCount = this.accountList.filter(item => Number(item.type) >= 5).length
      return [
        { label: '凭据总数', value: this.total, hint: '当前条件下的账号规模', tone: 'primary' },
        { label: 'Mysql', value: mysqlCount, hint: '数据库账号', tone: 'success' },
        { label: 'Jenkins', value: jenkinsCount, hint: '持续交付凭据', tone: 'warning' },
        { label: '通用账号', value: genericCount, hint: '监控与其他服务', tone: 'danger' }
      ]
    }
  },
  methods: {
    typeText(type) {
      return {
        1: 'Mysql',
        2: 'Postgre',
        3: 'Redis',
        4: 'Jenkins',
        5: 'Zabbix',
        6: '通用账号'
      }[Number(type)] || '未知类型'
    },
    tagType(type) {
      return {
        1: 'success',
        2: 'warning',
        3: 'danger',
        4: 'primary',
        5: 'info',
        6: 'info'
      }[Number(type)] || 'info'
    },
    async getList() {
      this.loading = true
      try {
        const { data: res } = await API.listAccountAuth({
          page: this.queryParams.pageNum,
          pageSize: this.queryParams.pageSize,
          alias: this.queryParams.alias || undefined,
          type: this.queryParams.type || undefined
        })
        if (res.code === 200) {
          this.accountList = res.data?.list || []
          this.total = res.data?.total || 0
        } else {
          this.$message.error(res.message || '获取账号列表失败')
        }
      } catch (error) {
        console.error('获取账号列表失败:', error)
        this.$message.error('获取账号列表失败')
      } finally {
        this.loading = false
      }
    },
    handleQuery() {
      this.queryParams.pageNum = 1
      this.getList()
    },
    resetQuery() {
      this.queryParams = {
        alias: '',
        type: undefined,
        pageNum: 1,
        pageSize: 10
      }
      this.getList()
    },
    handleSizeChange(val) {
      this.queryParams.pageSize = val
      this.queryParams.pageNum = 1
      this.getList()
    },
    handleCurrentChange(val) {
      this.queryParams.pageNum = val
      this.getList()
    },
    showAddDialog() {
      this.dialogTitle = '创建账号'
      this.$nextTick(() => {
        this.formData = {
          id: '',
          alias: '',
          host: '',
          name: '',
          password: '',
          type: undefined,
          remark: ''
        }
        this.dialogVisible = true
      })
    },
    showEditDialog(row) {
      this.dialogTitle = '修改账号'
      this.$nextTick(() => {
        this.formData = {
          id: row.id,
          alias: row.alias,
          host: formatAccountHostPortInput(row.host, row.port),
          name: row.name,
          password: '',
          type: row.type,
          remark: row.remark,
          createdAt: row.createdAt,
          updatedAt: row.updatedAt
        }
        this.dialogVisible = true
      })
    },
    async submitForm() {
      try {
        const valid = await this.$refs.formRef.validate().catch(() => false)
        if (!valid) return
        const parsedEndpoint = parseAccountHostPortInput(this.formData.host)
        if (!parsedEndpoint.valid) {
          this.$message.error('请输入服务地址，格式为 IP或域名:端口')
          return
        }
        const formData = {
          ...this.formData,
          host: parsedEndpoint.host,
          port: parsedEndpoint.port,
          type: Number(this.formData.type),
          id: this.formData.id ? Number(this.formData.id) : undefined,
          createdAt: this.formData.createdAt,
          updatedAt: this.formData.updatedAt
        }
        let res
        if (formData.id) {
          const updateData = { ...formData }
          delete updateData.createdAt
          delete updateData.updatedAt
          res = await API.updateAccountAuth({
            id: updateData.id,
            alias: updateData.alias,
            host: updateData.host,
            port: updateData.port,
            name: updateData.name,
            password: updateData.password,
            type: updateData.type,
            remark: updateData.remark
          })
        } else {
          const createData = { ...formData }
          delete createData.id
          res = await API.createAccountAuth(createData)
        }
        if (res.data.code === 200) {
          this.$message.success(formData.id ? '修改成功' : '创建成功')
          this.dialogVisible = false
          await this.getList()
        } else {
          this.$message.error(res.data.message || (formData.id ? '修改失败' : '创建失败'))
        }
      } catch (error) {
        console.error('操作失败:', error)
        this.$message.error(`操作失败: ${error.message}`)
      }
    },
    async handleDecrypt(row) {
      try {
        const res = await API.decryptPassword({}, {
          params: { id: row.id }
        })
        if (res.data.code === 200) {
          this.$alert(
            `<div><p>账号: ${row.alias}</p><p>密码: <span style="color: #9bd1ff; font-weight: 700; font-size: 18px;">${res.data.data.password}</span></p></div>`,
            '解密结果',
            {
              confirmButtonText: '确定',
              customClass: 'decrypt-result-alert',
              dangerouslyUseHTMLString: true
            }
          )
        } else {
          this.$message.error(res.data.message || '解密失败')
        }
      } catch (error) {
        console.error('解密失败:', error)
        this.$message.error(`解密失败: ${error.response?.data?.message || error.message}`)
      }
    },
    async handleDelete(row) {
      try {
        await this.$confirm(`确定删除账号 "${row.alias}"?`, '提示', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        })
        const res = await API.deleteAccountAuth(row.id)
        if (res.data.code === 200) {
          this.$message.success('删除成功')
          this.getList()
        } else {
          this.$message.error(res.data.message || '删除失败')
        }
      } catch (error) {
        if (error === 'cancel' || error === 'close') return
        console.error('删除失败:', error)
      }
    }
  },
  created() {
    this.getList()
  }
}
</script>

<style scoped>
.credential-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}

.credential-icon {
  width: 16px;
  height: 16px;
  object-fit: contain;
}
</style>
