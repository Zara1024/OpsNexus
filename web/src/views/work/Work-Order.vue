<template>
  <PageContainer tight class="work-apply-page">
    <PageHeader
      eyebrow="Work Order Center"
      title="工单申请"
      subtitle="先将脚本执行类变更统一收口到标准工单表单中，后续再持续吸纳快速发布、服务上线等更多入口。"
    >
      <template #actions>
        <el-button type="primary" @click="$router.push('/app/quick-release')">进入快速发布</el-button>
        <el-button @click="$router.push('/work/orders')">进入工单中心</el-button>
      </template>
      <template #intro>
        <PageIntro
          title="填写建议"
          text="工单标题尽量描述目标变更，变更原因写清影响范围与回退策略，脚本内容建议直接附上可执行命令。"
        />
      </template>
    </PageHeader>

    <SectionCard
      title="脚本执行工单"
      subtitle="当前阶段优先覆盖脚本执行类工单，适合临时运维、巡检修复和标准化发布脚本。"
    >
      <el-alert
        title="当前入口会直接创建脚本工单；如果你要做应用发布，请优先走快速发布页面。"
        type="info"
        :closable="false"
      />

      <el-form ref="formRef" :model="formData" :rules="rules" label-width="110px" class="apply-form">
        <el-form-item label="工单标题" prop="title">
          <el-input v-model="formData.title" placeholder="例如：修复巡检脚本的目录权限问题" />
        </el-form-item>

        <el-row :gutter="16" class="form-grid">
          <el-col :span="12">
            <el-form-item label="应用名称" prop="appName">
              <el-input v-model="formData.appName" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="应用编码" prop="appCode">
              <el-input v-model="formData.appCode" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="变更原因" prop="reason">
          <el-input v-model="formData.reason" type="textarea" :rows="4" placeholder="说明变更背景、影响范围和回退思路" />
        </el-form-item>

        <el-form-item label="执行目录" prop="executeDir">
          <el-input v-model="formData.executeDir" placeholder="/root/opsnexus-p0-acceptance" />
        </el-form-item>

        <el-form-item label="脚本内容" prop="scriptContent">
          <el-input
            v-model="formData.scriptContent"
            type="textarea"
            :rows="10"
            placeholder="例如：pwd&#10;ls -la&#10;echo work-order-ok"
          />
        </el-form-item>

        <el-form-item>
          <div class="form-actions">
            <el-button type="primary" :loading="submitting" @click="submitScriptWorkOrder">提交脚本工单</el-button>
            <el-button @click="resetForm">重置</el-button>
          </div>
        </el-form-item>
      </el-form>
    </SectionCard>
  </PageContainer>
</template>

<script>
import { ElMessage } from 'element-plus'
import PageContainer from '@/components/platform/PageContainer.vue'
import PageHeader from '@/components/platform/PageHeader.vue'
import PageIntro from '@/components/platform/PageIntro.vue'
import SectionCard from '@/components/platform/SectionCard.vue'

const createEmptyForm = () => ({
  title: '',
  reason: '',
  businessGroupId: 1,
  appId: 0,
  appName: 'opsnexus-local-script',
  appCode: 'opsnexus-local-script',
  executeDir: '/root/opsnexus-p0-acceptance',
  scriptContent: 'pwd\nls -la'
})

export default {
  name: 'WorkOrderApply',
  components: {
    PageContainer,
    PageHeader,
    PageIntro,
    SectionCard
  },
  data() {
    return {
      submitting: false,
      formData: createEmptyForm(),
      rules: {
        title: [{ required: true, message: '请输入工单标题', trigger: 'blur' }],
        reason: [{ required: true, message: '请输入变更原因', trigger: 'blur' }],
        appName: [{ required: true, message: '请输入应用名称', trigger: 'blur' }],
        appCode: [{ required: true, message: '请输入应用编码', trigger: 'blur' }],
        executeDir: [{ required: true, message: '请输入执行目录', trigger: 'blur' }],
        scriptContent: [{ required: true, message: '请输入脚本内容', trigger: 'blur' }]
      }
    }
  },
  methods: {
    async submitScriptWorkOrder() {
      const valid = await this.$refs.formRef.validate().catch(() => false)
      if (!valid) return
      this.submitting = true
      try {
        const { data: res } = await this.$api.createScriptWorkOrder(this.formData)
        if (res.code !== 200) {
          throw new Error(res.message || '创建脚本工单失败')
        }
        ElMessage.success('脚本工单已创建')
        this.$router.push('/work/orders')
      } catch (error) {
        ElMessage.error(error.message || '创建脚本工单失败')
      } finally {
        this.submitting = false
      }
    },
    resetForm() {
      this.formData = createEmptyForm()
    }
  }
}
</script>

<style scoped>
.work-apply-page :deep(.page-actions) {
  align-items: center;
}

.work-apply-page :deep(.section-card__body) {
  display: grid;
  gap: 20px;
}

.form-grid {
  margin-bottom: 6px;
}

.form-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

.apply-form :deep(.el-form-item) {
  margin-bottom: 22px;
}

.apply-form :deep(.el-form-item__label) {
  padding-bottom: 2px;
}

.apply-form :deep(.el-form-item__error) {
  margin-top: 6px;
}

.apply-form :deep(.el-form-item:last-child) {
  margin-bottom: 0;
}

.apply-form :deep(.el-textarea__inner) {
  min-height: 160px;
}

@media (max-width: 960px) {
  .form-actions {
    width: 100%;
  }
}
</style>
