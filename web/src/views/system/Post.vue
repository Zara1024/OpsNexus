<template>
  <TablePage
    class="system-post-page"
    section-title="岗位管理"
    section-subtitle="统一维护岗位编码、状态与岗位说明"
  >
    <template #header>
      <PageHeader
        eyebrow="Identity & Access"
        title="岗位管理"
        subtitle="维护岗位名称、编码与启停状态，支撑用户组织关系配置"
      >
        <template #actions>
          <el-button @click="getPostList">
            <el-icon><Refresh /></el-icon>
            刷新
          </el-button>
          <el-button type="primary" v-authority="['base:post:add']" @click="openAddDialog">
            <el-icon><Plus /></el-icon>
            新增岗位
          </el-button>
        </template>
        <template #intro>
          <PageIntro
            title="岗位说明"
            text="岗位用于补充用户职责分工，建议与角色搭配使用，保持编码稳定、说明清晰。"
          />
        </template>
      </PageHeader>
    </template>

    <template #toolbar>
      <PageToolbar>
        <el-form :model="queryParams" :inline="true" class="filter-cluster">
          <el-form-item class="filter-field" label="岗位名称">
            <el-input v-model="queryParams.postName" placeholder="请输入岗位名称" clearable @keyup.enter="handleQuery" />
          </el-form-item>
          <el-form-item class="filter-field" label="状态">
            <el-select v-model="queryParams.postStatus" placeholder="全部状态" clearable>
              <el-option v-for="item in postStatusList" :key="item.value" :label="item.label" :value="item.value" />
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
          <el-button type="danger" plain :disabled="multiple" v-authority="['base:post:delete']" @click="batchHandleDelete">
            <el-icon><Delete /></el-icon>
            批量删除
          </el-button>
        </template>
      </PageToolbar>
    </template>

    <template #stats>
      <StatStrip :items="statItems" />
    </template>

    <el-table
      v-loading="loading"
      :data="postList"
      stripe
      class="post-table"
      empty-text="暂无岗位数据"
      @selection-change="handleSelectionChange"
    >
      <el-table-column type="selection" width="52" />
      <el-table-column label="岗位名称" prop="postName" min-width="150" />
      <el-table-column label="岗位编码" prop="postCode" min-width="150" />
      <el-table-column label="状态" width="120">
        <template #default="{ row }">
          <el-switch
            v-model="row.postStatus"
            :active-value="1"
            :inactive-value="2"
            @change="postUpdateStatus(row)"
          />
        </template>
      </el-table-column>
      <el-table-column label="创建时间" prop="createTime" min-width="170" />
      <el-table-column label="说明" prop="remark" min-width="220" show-overflow-tooltip />
      <el-table-column label="操作" width="180" fixed="right">
        <template #default="{ row }">
          <div class="row-actions">
            <el-button type="primary" link v-authority="['base:post:edit']" @click="handleUpdate(row.id)">编辑</el-button>
            <el-button type="danger" link v-authority="['base:post:delete']" @click="handleDelete(row.id)">删除</el-button>
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

  <el-dialog title="新增岗位" v-model="addPostDialogVisible" width="560px" @closed="addPostDialogClosed">
    <el-form ref="addPostFormRefForm" label-width="90px" :rules="addPostFormRules" :model="addPostForm" class="post-form">
      <el-form-item label="岗位名称" prop="postName">
        <el-input v-model="addPostForm.postName" placeholder="请输入岗位名称" />
      </el-form-item>
      <el-form-item label="岗位编码" prop="postCode">
        <el-input v-model="addPostForm.postCode" placeholder="请输入岗位编码" />
      </el-form-item>
      <el-form-item label="状态" prop="postStatus">
        <el-radio-group v-model="addPostForm.postStatus">
          <el-radio :label="1">启用</el-radio>
          <el-radio :label="2">停用</el-radio>
        </el-radio-group>
      </el-form-item>
      <el-form-item label="说明" prop="remark">
        <el-input v-model="addPostForm.remark" type="textarea" :rows="4" placeholder="请输入岗位说明" />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="addPostDialogVisible = false">取消</el-button>
      <el-button type="primary" @click="addPost">确定</el-button>
    </template>
  </el-dialog>

  <el-dialog title="编辑岗位" v-model="editPostDialogVisible" width="560px" @closed="editPostDialogClosed">
    <el-form ref="editPostFormRefForm" label-width="90px" :rules="editPostFormRules" :model="editPostForm" class="post-form">
      <el-form-item label="岗位名称" prop="postName">
        <el-input v-model="editPostForm.postName" placeholder="请输入岗位名称" />
      </el-form-item>
      <el-form-item label="岗位编码" prop="postCode">
        <el-input v-model="editPostForm.postCode" placeholder="请输入岗位编码" />
      </el-form-item>
      <el-form-item label="状态" prop="postStatus">
        <el-radio-group v-model="editPostForm.postStatus">
          <el-radio :label="1">启用</el-radio>
          <el-radio :label="2">停用</el-radio>
        </el-radio-group>
      </el-form-item>
      <el-form-item label="说明" prop="remark">
        <el-input v-model="editPostForm.remark" type="textarea" :rows="4" placeholder="请输入岗位说明" />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="editPostDialogVisible = false">取消</el-button>
      <el-button type="primary" @click="editPost">确定</el-button>
    </template>
  </el-dialog>
</template>

<script>
import { Delete, Plus, Refresh, Search } from '@element-plus/icons-vue'
import PageHeader from '@/components/platform/PageHeader.vue'
import PageIntro from '@/components/platform/PageIntro.vue'
import PageToolbar from '@/components/platform/PageToolbar.vue'
import StatStrip from '@/components/platform/StatStrip.vue'
import TablePage from '@/components/platform/TablePage.vue'

const createEmptyQuery = () => ({
  postName: '',
  postStatus: '',
  pageNum: 1,
  pageSize: 10
})

const createEmptyPost = () => ({
  postName: '',
  postCode: '',
  postStatus: 1,
  remark: ''
})

export default {
  name: 'SystemPost',
  components: {
    Delete,
    PageHeader,
    PageIntro,
    PageToolbar,
    Plus,
    Refresh,
    Search,
    StatStrip,
    TablePage
  },
  data() {
    return {
      queryParams: createEmptyQuery(),
      postStatusList: [
        { value: '1', label: '启用' },
        { value: '2', label: '停用' }
      ],
      loading: true,
      postList: [],
      total: 0,
      addPostDialogVisible: false,
      addPostFormRules: {
        postName: [{ required: true, message: '请输入岗位名称', trigger: 'blur' }],
        postCode: [{ required: true, message: '请输入岗位编码', trigger: 'blur' }],
        postStatus: [{ required: true, message: '请选择岗位状态', trigger: 'change' }]
      },
      addPostForm: createEmptyPost(),
      editPostDialogVisible: false,
      editPostForm: createEmptyPost(),
      editPostFormRules: {
        postName: [{ required: true, message: '请输入岗位名称', trigger: 'blur' }],
        postCode: [{ required: true, message: '请输入岗位编码', trigger: 'blur' }],
        postStatus: [{ required: true, message: '请选择岗位状态', trigger: 'change' }]
      },
      ids: [],
      multiple: true
    }
  },
  computed: {
    statItems() {
      const enabledCount = this.postList.filter(item => Number(item.postStatus) === 1).length
      const disabledCount = this.postList.filter(item => Number(item.postStatus) === 2).length
      return [
        { label: '岗位总数', value: this.total, hint: '当前岗位覆盖数量', tone: 'primary' },
        { label: '启用岗位', value: enabledCount, hint: '当前启用岗位', tone: 'success' },
        { label: '停用岗位', value: disabledCount, hint: '当前停用岗位', tone: 'danger' },
        { label: '已选岗位', value: this.ids.length, hint: '当前批量操作数量', tone: 'warning' }
      ]
    }
  },
  methods: {
    openAddDialog() {
      this.addPostForm = createEmptyPost()
      this.addPostDialogVisible = true
    },
    async getPostList() {
      this.loading = true
      const { data: res } = await this.$api.queryPostList(this.queryParams)
      if (res.code !== 200) {
        this.$message.error(res.message)
      } else {
        this.postList = res.data.list || []
        this.total = res.data.total || 0
        this.loading = false
      }
    },
    handleQuery() {
      this.queryParams.pageNum = 1
      this.getPostList()
    },
    resetQuery() {
      this.queryParams = createEmptyQuery()
      this.getPostList()
      this.$message.success('已重置筛选')
    },
    handleSizeChange(newSize) {
      this.queryParams.pageSize = newSize
      this.queryParams.pageNum = 1
      this.getPostList()
    },
    handleCurrentChange(newPage) {
      this.queryParams.pageNum = newPage
      this.getPostList()
    },
    async postUpdateStatus(row) {
      const text = Number(row.postStatus) === 2 ? '停用' : '启用'
      const confirmResult = await this.$confirm(`确认${text}岗位 ${row.postName} 吗`, '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).catch(err => err)
      if (confirmResult !== 'confirm') {
        await this.getPostList()
        return this.$message.info('已取消操作')
      }
      await this.$api.updatePostStatus(row.id, row.postStatus)
      this.$message.success(text + '成功')
      await this.getPostList()
    },
    addPostDialogClosed() {
      this.$refs.addPostFormRefForm?.resetFields()
      this.addPostForm = createEmptyPost()
    },
    addPost() {
      this.$refs.addPostFormRefForm.validate(async valid => {
        if (!valid) return
        const { data: res } = await this.$api.addPost(this.addPostForm)
        if (res.code !== 200) {
          this.$message.error(res.message)
        } else {
          this.$message.success('新增岗位成功')
          this.addPostDialogVisible = false
          await this.getPostList()
        }
      })
    },
    async handleUpdate(id) {
      const { data: res } = await this.$api.postInfo(id)
      if (res.code !== 200) {
        this.$message.error(res.message)
      } else {
        this.editPostForm = res.data || createEmptyPost()
        this.editPostDialogVisible = true
      }
    },
    editPostDialogClosed() {
      this.$refs.editPostFormRefForm?.resetFields()
      this.editPostForm = createEmptyPost()
    },
    editPost() {
      this.$refs.editPostFormRefForm.validate(async valid => {
        if (!valid) return
        const { data: res } = await this.$api.updatePost(this.editPostForm)
        if (res.code !== 200) {
          this.$message.error(res.message)
        } else {
          this.$message.success('更新岗位成功')
          this.editPostDialogVisible = false
          await this.getPostList()
        }
      })
    },
    async handleDelete(id) {
      const confirmResult = await this.$confirm(`确认删除岗位 ${id} 吗`, '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).catch(err => err)
      if (confirmResult !== 'confirm') {
        return this.$message.info('已取消操作')
      }
      const { data: res } = await this.$api.deleteSysPost(id)
      if (res.code !== 200) {
        this.$message.error(res.message)
      } else {
        this.$message.success('删除成功')
        await this.getPostList()
      }
    },
    handleSelectionChange(selection) {
      this.ids = selection.map(item => item.id)
      this.multiple = !selection.length
    },
    async batchHandleDelete() {
      const postIds = this.ids
      const confirmResult = await this.$confirm(`确认删除选中的岗位 ${postIds} 吗`, '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).catch(err => err)
      if (confirmResult !== 'confirm') {
        return this.$message.info('已取消操作')
      }
      const { data: res } = await this.$api.batchDeleteSysPost(postIds)
      if (res.code !== 200) {
        this.$message.error(res.message)
      } else {
        this.$message.success('删除成功')
        await this.getPostList()
      }
    }
  },
  created() {
    this.getPostList()
  }
}
</script>

<style scoped>
.system-post-page :deep(.page-actions) {
  align-items: center;
}

.post-form :deep(.el-textarea__inner) {
  min-height: 120px;
}

@media (max-width: 768px) {
  .system-post-page :deep(.page-actions) {
    width: 100%;
    justify-content: flex-start;
  }
}
</style>
