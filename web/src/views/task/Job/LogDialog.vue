<template>
  <div v-if="visible" class="custom-dialog log-dialog">
    <div class="dialog-mask"></div>
    <div class="dialog-container">
      <div class="dialog-header">
        <h3>{{ title }}</h3>
        <button class="dialog-close" @click="handleClose">×</button>
      </div>
      <div class="dialog-body" v-html="content"></div>
      <div class="dialog-footer">
        <button class="dialog-button" @click="handleClose">关闭</button>
      </div>
    </div>
  </div>
</template>

<script>
import { highlight } from '@/utils/highlight'

export default {
  name: 'LogDialog',
  props: {
    title: {
      type: String,
      default: '任务日志'
    }
  },
  data() {
    return {
      visible: false,
      content: ''
    }
  },
  methods: {
    show(content) {
      let processedContent = content
      if (typeof processedContent !== 'string') {
        processedContent = processedContent.logs || JSON.stringify(processedContent, null, 2)
      }
      processedContent = processedContent.replace(/\\n/g, '\n').replace(/\\r/g, '\r')
      const highlightedContent = highlight(processedContent, 'bash')
      this.content = `
        <div class="highlight-container">
          <pre class="dialog-code dialog-code--log">${highlightedContent}</pre>
        </div>
      `
      this.visible = true
    },
    handleClose() {
      this.visible = false
      this.$emit('close')
    }
  }
}
</script>

<style scoped>
.custom-dialog {
  position: fixed;
  inset: 0;
  z-index: 9999;
}

.dialog-mask {
  position: absolute;
  inset: 0;
  background: rgba(2, 6, 23, 0.72);
  backdrop-filter: blur(6px);
}

.dialog-container {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: min(1200px, 88vw);
  min-width: 800px;
  background: linear-gradient(180deg, rgba(30, 41, 59, 0.98), rgba(15, 23, 42, 0.98));
  border: 1px solid rgba(148, 163, 184, 0.24);
  border-radius: 20px;
  box-shadow: 0 28px 80px rgba(2, 8, 23, 0.56);
  color: rgba(226, 232, 240, 0.92);
}

.dialog-header,
.dialog-footer {
  padding: 18px 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.dialog-header {
  border-bottom: 1px solid rgba(148, 163, 184, 0.14);
}

.dialog-header h3 {
  margin: 0;
  font-size: 18px;
  color: #f8fbff;
}

.dialog-body {
  padding: 20px;
  max-height: 70vh;
  overflow: auto;
}

.dialog-footer {
  border-top: 1px solid rgba(148, 163, 184, 0.14);
  justify-content: flex-end;
}

.dialog-close,
.dialog-button {
  border: 1px solid rgba(148, 163, 184, 0.22);
  background: rgba(255, 255, 255, 0.05);
  color: rgba(226, 232, 240, 0.92);
  border-radius: 12px;
  cursor: pointer;
}

.dialog-close {
  width: 36px;
  height: 36px;
  font-size: 20px;
}

.dialog-button {
  padding: 8px 16px;
}

.dialog-close:hover,
.dialog-button:hover {
  background: rgba(59, 130, 246, 0.16);
  border-color: rgba(59, 130, 246, 0.3);
}

:deep(.dialog-code) {
  margin: 0;
  min-height: 320px;
  padding: 18px 20px;
  border-radius: 16px;
  border: 1px solid rgba(148, 163, 184, 0.14);
  background: rgba(4, 11, 23, 0.96);
  color: rgba(226, 232, 240, 0.92);
  white-space: pre-wrap;
  line-height: 1.7;
  overflow: auto;
}

:deep(.hljs) {
  background: transparent;
}
</style>
