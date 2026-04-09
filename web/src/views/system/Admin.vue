<template>
  <PageContainer wide class="system-admin-page">
    <PageHeader
      eyebrow="Identity & Access"
      title="用户管理"
      subtitle="统一维护管理员账号、部门归属、岗位角色与启停状态"
    >
      <template #actions>
        <el-button @click="refreshAll">
          <el-icon><Refresh /></el-icon>
          刷新
        </el-button>
        <el-button type="primary" v-authority="['base:admin:add']" @click="openAddDialog">
          <el-icon><Plus /></el-icon>
          新增用户
        </el-button>
      </template>
      <template #meta>
        <div class="page-chip-row">
          <span class="platform-chip">部门 {{ deptCount }}</span>
          <span class="platform-chip">角色 {{ roleList.length }}</span>
          <span class="platform-chip">岗位 {{ postList.length }}</span>
          <span v-if="currentDeptLabel" class="platform-chip">当前部门 {{ currentDeptLabel }}</span>
        </div>
      </template>
      <template #intro>
        <PageIntro
          title="治理说明"
          text="左侧按部门定位账号范围，右侧统一维护用户、角色与岗位归属，所有改动均保持现有权限与登录链路不变。"
        />
      </template>
    </PageHeader>

    <StatStrip :items="statItems" />

    <div class="split-layout admin-split-layout">
      <SectionCard title="组织架构" subtitle="按部门筛选用户账号范围" class="admin-tree-card">
        <template #actions>
          <el-button v-if="currentDeptId" type="primary" link @click="clearDeptScope">清除范围</el-button>
        </template>

        <div class="admin-tree-summary page-chip-row">
          <span class="platform-chip">部门 {{ deptCount }}</span>
          <span class="platform-chip">用户 {{ tableTotal }}</span>
        </div>

        <el-tree
          v-if="deptList.length"
          :data="deptList"
          :props="defaultProps"
          node-key="id"
          default-expand-all
          :highlight-current="true"
          class="admin-tree"
          @node-click="handleDeptClick"
          @node-expand="handleNodeExpand"
          @node-collapse="handleNodeCollapse"
        >
          <template #default="{ node, data }">
            <span class="admin-tree-node">
              <el-icon class="admin-tree-node__icon">
                <component :is="!data.parentId && expandedKeys.includes(node.key) ? 'FolderOpened' : 'Folder'" />
              </el-icon>
              <el-tooltip :content="node.label" placement="top" :show-after="300">
                <span :class="['admin-tree-node__label', { 'admin-tree-node__label--root': !data.parentId }]">
                  {{ formatCompactLabel(node.label, 16, 6) }}
                </span>
              </el-tooltip>
            </span>
          </template>
        </el-tree>
        <div v-else class="empty-state admin-tree-empty">
          <div class="empty-state__halo"></div>
          <div class="empty-state__title">暂无部门数据</div>
          <div class="empty-state__description">请先维护组织架构后再按部门查看用户</div>
        </div>
      </SectionCard>

      <SectionCard
        :title="listSectionTitle"
        :subtitle="listSectionSubtitle"
        class="admin-table-card"
      >
        <PageToolbar class="admin-toolbar">
          <el-form :inline="true" :model="queryParams" class="filter-cluster">
            <el-form-item class="filter-field" label="用户名">
              <el-input
                v-model="queryParams.username"
                placeholder="请输入用户名"
                clearable
                @keyup.enter="handleQuery"
              />
            </el-form-item>
            <el-form-item class="filter-field" label="状态">
              <el-select v-model="queryParams.status" placeholder="全部状态" clearable>
                <el-option v-for="item in statusList" :key="item.value" :label="item.label" :value="item.value" />
              </el-select>
            </el-form-item>
          </el-form>
          <template #actions>
            <el-button type="primary" @click="handleQuery">
              <el-icon><Search /></el-icon>
              查询
            </el-button>
            <el-button @click="resetQuery">
              <el-icon><Refresh /></el-icon>
              重置
            </el-button>
          </template>
        </PageToolbar>

        <el-table
          v-loading="Loading"
          :data="displayAdminList"
          stripe
          class="admin-table"
          empty-text="暂无用户数据"
        >
          <el-table-column label="用户名" prop="username" min-width="150" />
          <el-table-column label="昵称" prop="nickname" min-width="140" />
          <el-table-column label="部门" prop="deptName" min-width="140" show-overflow-tooltip>
            <template #default="{ row }">
              {{ row.deptName || currentDeptLabel || '-' }}
            </template>
          </el-table-column>
          <el-table-column label="角色" min-width="130" show-overflow-tooltip>
            <template #default="{ row }">
              <el-tag type="primary">{{ row.roleName || '-' }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column label="状态" width="120">
            <template #default="{ row }">
              <el-switch
                v-model="row.status"
                :active-value="1"
                :inactive-value="2"
                @change="adminUpdateStatus(row)"
              />
            </template>
          </el-table-column>
          <el-table-column label="操作" width="180" fixed="right">
            <template #default="{ row }">
              <div class="row-actions">
                <el-button type="primary" link v-authority="['base:admin:edit']" @click="showEditAdminDialog(row.id)">
                  编辑
                </el-button>
                <el-button type="warning" link v-authority="['base:admin:reset']" @click="handleResetPwd(row)">
                  重置密码
                </el-button>
                <el-button type="danger" link v-authority="['base:admin:delete']" @click="handleAdminDelete(row)">
                  删除
                </el-button>
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
            :total="tableTotal"
            @size-change="handleSizeChange"
            @current-change="handleCurrentChange"
          />
        </template>
      </SectionCard>
    </div>

    <el-dialog v-model="addDialogVisible" title="新增用户" width="820px" @closed="addDialogClosed">
      <el-form ref="addFormRefForm" :model="addForm" :rules="addFormRules" label-width="92px" class="admin-form">
        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="用户名" prop="username">
              <el-input v-model="addForm.username" placeholder="请输入用户名" maxlength="30" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="密码" prop="password">
              <el-input v-model="addForm.password" placeholder="请输入密码" type="password" maxlength="20" show-password />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="昵称" prop="nickname">
              <el-input v-model="addForm.nickname" placeholder="请输入昵称" maxlength="30" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="所属部门" prop="deptId">
              <Treeselect v-model="addForm.deptId" :options="deptList" :show-count="true" placeholder="请选择所属部门" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="手机号" prop="phone">
              <el-input v-model="addForm.phone" placeholder="请输入手机号" maxlength="11" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="邮箱" prop="email">
              <el-input v-model="addForm.email" placeholder="请输入邮箱" maxlength="50" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="岗位" prop="postId">
              <el-select v-model="addForm.postId" placeholder="请选择岗位" style="width: 100%">
                <el-option v-for="item in postList" :key="item.id" :label="item.postName" :value="item.id" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="角色" prop="roleId">
              <el-select v-model="addForm.roleId" placeholder="请选择角色" style="width: 100%">
                <el-option v-for="item in roleList" :key="item.id" :label="item.roleName" :value="item.id" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item label="状态" prop="status">
              <el-radio-group v-model="addForm.status">
                <el-radio :label="1">启用</el-radio>
                <el-radio :label="2">停用</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item label="备注" prop="note">
              <el-input v-model="addForm.note" type="textarea" :rows="4" placeholder="请输入备注信息" />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <template #footer>
        <el-button @click="addDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="addAdmin">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="editDialogVisible" title="编辑用户" width="820px" @closed="editDialogClosed">
      <el-form ref="editFormRefForm" :model="adminInfo" :rules="editFormRules" label-width="92px" class="admin-form">
        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="用户名" prop="username">
              <el-input v-model="adminInfo.username" placeholder="请输入用户名" maxlength="30" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="昵称" prop="nickname">
              <el-input v-model="adminInfo.nickname" placeholder="请输入昵称" maxlength="30" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="所属部门" prop="deptId">
              <Treeselect v-model="adminInfo.deptId" :options="deptList" :show-count="true" placeholder="请选择所属部门" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="状态" prop="status">
              <el-radio-group v-model="adminInfo.status">
                <el-radio :label="1">启用</el-radio>
                <el-radio :label="2">停用</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="手机号" prop="phone">
              <el-input v-model="adminInfo.phone" placeholder="请输入手机号" maxlength="11" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="邮箱" prop="email">
              <el-input v-model="adminInfo.email" placeholder="请输入邮箱" maxlength="50" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="岗位" prop="postId">
              <el-select v-model="adminInfo.postId" placeholder="请选择岗位" style="width: 100%">
                <el-option v-for="item in postList" :key="item.id" :label="item.postName" :value="item.id" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="角色" prop="roleId">
              <el-select v-model="adminInfo.roleId" placeholder="请选择角色" style="width: 100%">
                <el-option v-for="item in roleList" :key="item.id" :label="item.roleName" :value="item.id" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item label="备注" prop="note">
              <el-input v-model="adminInfo.note" type="textarea" :rows="4" placeholder="请输入备注信息" />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <template #footer>
        <el-button @click="editDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="editAdminInfo">确定</el-button>
      </template>
    </el-dialog>
  </PageContainer>
</template>

<script>
import Treeselect from 'vue3-treeselect'
import 'vue3-treeselect/dist/vue3-treeselect.css'
import { Delete, Edit, Folder, FolderOpened, Key, Plus, Refresh, Search } from '@element-plus/icons-vue'
import PageContainer from '@/components/platform/PageContainer.vue'
import PageHeader from '@/components/platform/PageHeader.vue'
import PageIntro from '@/components/platform/PageIntro.vue'
import PageToolbar from '@/components/platform/PageToolbar.vue'
import SectionCard from '@/components/platform/SectionCard.vue'
import StatStrip from '@/components/platform/StatStrip.vue'

const createEmptyQuery = () => ({
  username: '',
  status: '',
  pageNum: 1,
  pageSize: 10
})

const createEmptyAdminForm = () => ({
  username: '',
  password: '',
  deptId: undefined,
  postId: undefined,
  roleId: undefined,
  email: '',
  nickname: '',
  status: 1,
  phone: '',
  note: ''
})

const flattenDeptNodes = nodes => {
  let count = 0
  ;(nodes || []).forEach(node => {
    count += 1
    if (Array.isArray(node.children) && node.children.length) {
      count += flattenDeptNodes(node.children)
    }
  })
  return count
}

const normalizeDeptUsers = rows => (rows || []).map(item => ({
  id: item.id,
  username: item.username,
  nickname: item.nickname,
  status: item.status,
  icon: item.icon,
  email: item.email,
  phone: item.phone,
  note: item.note,
  deptId: item.deptId || item.dept_id,
  deptName: item.dept_name || item.deptName,
  postId: item.postId || item.post_id,
  postName: item.post_name || item.postName,
  roleId: item.roleId || item.role_id,
  roleName: item.role_name || item.roleName,
  createTime: item.create_time || item.createTime
}))

export default {
  name: 'SystemAdmin',
  components: {
    Delete,
    Edit,
    Folder,
    FolderOpened,
    Key,
    PageContainer,
    PageHeader,
    PageIntro,
    PageToolbar,
    Plus,
    Refresh,
    Search,
    SectionCard,
    StatStrip,
    Treeselect
  },
  data() {
    return {
      expandedKeys: [],
      statusList: [
        { value: '1', label: '启用' },
        { value: '2', label: '停用' }
      ],
      Loading: false,
      queryParams: createEmptyQuery(),
      adminList: [],
      scopedDeptUsers: [],
      total: 0,
      addDialogVisible: false,
      deptList: [],
      roleList: [],
      postList: [],
      addForm: createEmptyAdminForm(),
      addFormRules: {
        deptId: [{ required: true, message: '请选择所属部门', trigger: 'change' }],
        postId: [{ required: true, message: '请选择岗位', trigger: 'change' }],
        roleId: [{ required: true, message: '请选择角色', trigger: 'change' }],
        username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
        password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
        status: [{ required: true, message: '请选择状态', trigger: 'change' }],
        email: [{ required: true, message: '请输入邮箱', trigger: 'blur' }],
        nickname: [{ required: true, message: '请输入昵称', trigger: 'blur' }],
        phone: [{ required: true, message: '请输入手机号', trigger: 'blur' }]
      },
      editDialogVisible: false,
      adminInfo: createEmptyAdminForm(),
      editFormRules: {
        deptId: [{ required: true, message: '请选择所属部门', trigger: 'change' }],
        postId: [{ required: true, message: '请选择岗位', trigger: 'change' }],
        roleId: [{ required: true, message: '请选择角色', trigger: 'change' }],
        username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
        status: [{ required: true, message: '请选择状态', trigger: 'change' }],
        email: [{ required: true, message: '请输入邮箱', trigger: 'blur' }],
        nickname: [{ required: true, message: '请输入昵称', trigger: 'blur' }],
        phone: [{ required: true, message: '请输入手机号', trigger: 'blur' }]
      },
      defaultProps: {
        children: 'children',
        label: 'label'
      },
      currentDeptId: null,
      currentDeptLabel: ''
    }
  },
  computed: {
    filteredDeptUsers() {
      const username = String(this.queryParams.username || '').trim().toLowerCase()
      const status = String(this.queryParams.status || '')
      return this.scopedDeptUsers.filter(item => {
        const matchUsername = !username || String(item.username || '').toLowerCase().includes(username)
        const matchStatus = !status || String(item.status) === status
        return matchUsername && matchStatus
      })
    },
    displayAdminList() {
      if (!this.currentDeptId) return this.adminList
      const start = (this.queryParams.pageNum - 1) * this.queryParams.pageSize
      return this.filteredDeptUsers.slice(start, start + this.queryParams.pageSize)
    },
    tableTotal() {
      return this.currentDeptId ? this.filteredDeptUsers.length : this.total
    },
    deptCount() {
      return flattenDeptNodes(this.deptList)
    },
    listSectionTitle() {
      return this.currentDeptLabel ? `${this.currentDeptLabel} 用户` : '用户列表'
    },
    listSectionSubtitle() {
      return this.currentDeptLabel
        ? `仅显示 ${this.currentDeptLabel} 下的账号与权限归属`
        : '支持按部门、状态与用户名筛选平台账号'
    },
    statItems() {
      const source = this.currentDeptId ? this.filteredDeptUsers : this.adminList
      const enabledCount = source.filter(item => Number(item.status) === 1).length
      const disabledCount = source.filter(item => Number(item.status) === 2).length
      const roleCoverage = new Set(source.map(item => item.roleName).filter(Boolean)).size

      return [
        { label: '账号总数', value: this.tableTotal, hint: this.currentDeptLabel ? '当前部门账号' : '平台账号总量', tone: 'primary' },
        { label: '启用账号', value: enabledCount, hint: '当前可登录账号', tone: 'success' },
        { label: '停用账号', value: disabledCount, hint: '当前停用账号', tone: 'danger' },
        { label: '角色覆盖', value: roleCoverage, hint: '已关联角色数', tone: 'warning' }
      ]
    }
  },
  created() {
    this.refreshAll()
  },
  methods: {
    formatCompactLabel(value, prefixLength = 16, suffixLength = 6) {
      const text = String(value || '')
      if (!text) return '-'
      if (text.length <= prefixLength + suffixLength) {
        return text
      }
      return `${text.slice(0, prefixLength)}...${text.slice(-suffixLength)}`
    },
    avatarText(value) {
      return String(value || '?').trim().slice(0, 1).toUpperCase()
    },
    async refreshAll() {
      await Promise.all([this.getDeptVoList(), this.getRoleVoList(), this.getPostVoList()])
      await this.refreshTableData(true)
    },
    async refreshTableData(preservePage = false) {
      if (!preservePage) {
        this.queryParams.pageNum = 1
      }
      if (this.currentDeptId) {
        await this.loadAdminsByDept(this.currentDeptId, this.currentDeptLabel, true)
      } else {
        await this.getAdminList()
      }
    },
    async getAdminList() {
      this.Loading = true
      try {
        const { data: res } = await this.$api.queryAdminList(this.queryParams)
        if (res.code !== 200) {
          this.$message.error(res.message)
          return
        }
        this.adminList = res.data.list || []
        this.total = res.data.total || 0
      } finally {
        this.Loading = false
      }
    },
    async loadAdminsByDept(deptId, deptLabel = '', preservePage = false) {
      this.Loading = true
      try {
        const { data: res } = await this.$api.deptUsers(deptId)
        if (res.code !== 200) {
          this.$message.error(res.message)
          return
        }
        this.currentDeptId = deptId
        this.currentDeptLabel = deptLabel || this.currentDeptLabel
        this.scopedDeptUsers = normalizeDeptUsers(res.data)
        if (!preservePage) {
          this.queryParams.pageNum = 1
        }
      } finally {
        this.Loading = false
      }
    },
    handleQuery() {
      this.queryParams.pageNum = 1
      if (!this.currentDeptId) {
        this.getAdminList()
      }
    },
    resetQuery() {
      const pageSize = this.queryParams.pageSize
      this.queryParams = {
        ...createEmptyQuery(),
        pageSize
      }
      if (!this.currentDeptId) {
        this.getAdminList()
      }
    },
    handleSizeChange(size) {
      this.queryParams.pageSize = size
      this.queryParams.pageNum = 1
      if (!this.currentDeptId) {
        this.getAdminList()
      }
    },
    handleCurrentChange(page) {
      this.queryParams.pageNum = page
      if (!this.currentDeptId) {
        this.getAdminList()
      }
    },
    async adminUpdateStatus(row) {
      const text = Number(row.status) === 1 ? '启用' : '停用'
      const confirmResult = await this.$confirm(`确认${text}账号 ${row.username} 吗`, '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).catch(err => err)

      if (confirmResult !== 'confirm') {
        await this.refreshTableData(true)
        return this.$message.info('已取消操作')
      }

      const { data: res } = await this.$api.updateAdminStatus(row.id, row.status)
      if (res.code !== 200) {
        this.$message.error(res.message)
      } else {
        this.$message.success(`${text}成功`)
      }
      await this.refreshTableData(true)
    },
    async getDeptVoList() {
      const { data: res } = await this.$api.querySysDeptVoList()
      if (res.code !== 200) {
        this.$message.error(res.message)
      } else {
        this.deptList = this.$handleTree.handleTree(res.data, 'id')
      }
    },
    async getRoleVoList() {
      const { data: res } = await this.$api.querySysRoleVoList()
      if (res.code !== 200) {
        this.$message.error(res.message)
      } else {
        this.roleList = res.data || []
      }
    },
    async getPostVoList() {
      const { data: res } = await this.$api.querySysPostVoList()
      if (res.code !== 200) {
        this.$message.error(res.message)
      } else {
        this.postList = res.data || []
      }
    },
    openAddDialog() {
      this.addForm = createEmptyAdminForm()
      this.addDialogVisible = true
    },
    addDialogClosed() {
      this.$refs.addFormRefForm?.resetFields()
      this.addForm = createEmptyAdminForm()
    },
    addAdmin() {
      this.$refs.addFormRefForm.validate(async valid => {
        if (!valid) return
        const { data: res } = await this.$api.addAdmin(this.addForm)
        if (res.code !== 200) {
          this.$message.error(res.message)
        } else {
          this.$message.success('新增用户成功')
          this.addDialogVisible = false
          await this.refreshTableData(true)
        }
      })
    },
    async showEditAdminDialog(id) {
      const { data: res } = await this.$api.adminInfo(id)
      if (res.code !== 200) {
        this.$message.error(res.message)
      } else {
        this.adminInfo = {
          ...createEmptyAdminForm(),
          ...res.data
        }
        this.editDialogVisible = true
      }
    },
    editDialogClosed() {
      this.$refs.editFormRefForm?.resetFields()
      this.adminInfo = createEmptyAdminForm()
    },
    editAdminInfo() {
      this.$refs.editFormRefForm.validate(async valid => {
        if (!valid) return
        const { data: res } = await this.$api.adminUpdate(this.adminInfo)
        if (res.code !== 200) {
          this.$message.error(res.message)
        } else {
          this.editDialogVisible = false
          await this.refreshTableData(true)
          this.$message.success('更新用户成功')
        }
      })
    },
    async handleAdminDelete(row) {
      const confirmResult = await this.$confirm(`确认删除用户 ${row.username} 吗`, '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).catch(err => err)
      if (confirmResult !== 'confirm') {
        return this.$message.info('已取消操作')
      }
      const { data: res } = await this.$api.deleteAdmin(row.id)
      if (res.code !== 200) {
        this.$message.error(res.message)
      } else {
        this.$message.success('删除成功')
        await this.refreshTableData(true)
      }
    },
    handleResetPwd(row) {
      this.$prompt(`请输入 ${row.username} 的新密码`, '重置密码', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        closeOnClickModal: false,
        inputPattern: /^.{5,20}$/,
        inputErrorMessage: '密码长度需为 5 到 20 位'
      }).then(({ value }) => {
        this.$api.resetPassword(row.id, value).then(() => {
          this.$message.success(`密码已重置，新密码为：${value}`)
        })
      }).catch(() => {})
    },
    handleNodeExpand(data, node) {
      if (!this.expandedKeys.includes(node.key)) {
        this.expandedKeys.push(node.key)
      }
    },
    handleNodeCollapse(data, node) {
      this.expandedKeys = this.expandedKeys.filter(key => key !== node.key)
    },
    handleDeptClick(data, node) {
      const deptId = data?.id || node?.key
      if (!deptId) {
        this.$message.warning('未获取到部门 ID')
        return
      }
      this.loadAdminsByDept(deptId, data?.label || node?.label || '')
    },
    clearDeptScope() {
      this.currentDeptId = null
      this.currentDeptLabel = ''
      this.scopedDeptUsers = []
      this.queryParams.pageNum = 1
      this.getAdminList()
    }
  }
}
</script>

<style scoped>
.system-admin-page :deep(.page-actions) {
  align-items: center;
}

.admin-split-layout {
  align-items: start;
}

.system-admin-page :deep(.admin-split-layout) {
  grid-template-columns: 260px minmax(0, 1fr);
}

.admin-tree-card,
.admin-table-card {
  min-width: 0;
}

.admin-tree-summary {
  margin-bottom: 16px;
}

.admin-toolbar {
  margin-bottom: 18px;
  padding: 0;
  background: transparent;
  border: none;
  box-shadow: none;
}

.admin-tree {
  padding-right: 4px;
}

.admin-tree-node {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  min-width: 0;
}

.admin-tree-node__icon {
  color: rgba(125, 211, 252, 0.9);
}

.admin-tree-node__label {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.admin-tree-node__label--root {
  font-weight: 700;
  color: var(--text-primary);
}

.admin-tree-empty {
  min-height: 300px;
}

.admin-avatar {
  background: linear-gradient(135deg, rgba(37, 99, 235, 0.8), rgba(14, 165, 233, 0.88));
  color: #fff;
  font-weight: 700;
}

@media (max-width: 1080px) {
  .admin-split-layout {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .system-admin-page :deep(.page-actions) {
    width: 100%;
    justify-content: flex-start;
  }
}
</style>
