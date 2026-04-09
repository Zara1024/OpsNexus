<template>
  <div class="kubectl-terminal-page">
    <el-card class="kubectl-terminal-card">
      <template #header>
        <div class="terminal-header">
          <div class="terminal-title">
            <el-icon><Monitor /></el-icon>
            <span>kubectl 终端</span>
            <el-tag size="small" type="info">集群 {{ clusterId }}</el-tag>
          </div>
          <div class="terminal-controls">
            <el-select
              v-model="selectedNamespace"
              placeholder="选择命名空间"
              filterable
              style="width: 220px"
              @change="handleNamespaceChange"
            >
              <el-option
                v-for="item in namespaceOptions"
                :key="item.name"
                :label="item.name"
                :value="item.name"
              />
            </el-select>
            <el-button type="primary" :loading="connecting" @click="connectTerminal">
              {{ isConnected ? '重新连接' : '连接终端' }}
            </el-button>
            <el-button :disabled="!uiState.canDisconnect" @click="disconnectTerminal">断开</el-button>
          </div>
        </div>
      </template>

      <div class="command-toolbar">
        <el-input
          v-model="quickCommand"
          placeholder="输入 kubectl 命令，例如：get pods -A"
          @keyup.enter="executeQuickCommand"
        >
          <template #prepend>kubectl</template>
        </el-input>
        <el-button type="success" :loading="executingCommand" @click="executeQuickCommand">
          执行命令
        </el-button>
      </div>

      <div class="terminal-wrapper">
        <div ref="terminalElement" class="terminal-surface" @click="focusTerminal"></div>
        <div v-if="uiState.showPlaceholder" class="terminal-placeholder">
          <p>选择命名空间后连接终端，或直接执行上方 kubectl 命令。</p>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { computed, nextTick, onMounted, onUnmounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Monitor } from '@element-plus/icons-vue'
import { Terminal as XTerm } from '@xterm/xterm'
import { FitAddon } from '@xterm/addon-fit'
import '@xterm/xterm/css/xterm.css'
import k8sApi from '@/api/k8s'
import storage from '@/utils/storage'
import { getKubectlTerminalUiState } from '@/utils/kubectlTerminalPresentation.mjs'

const route = useRoute()

const clusterId = computed(() => route.params.clusterId)
const namespaceOptions = ref([])
const selectedNamespace = ref('default')
const quickCommand = ref('')
const connecting = ref(false)
const executingCommand = ref(false)
const isConnected = ref(false)
const hasTerminalContent = ref(false)

const terminalElement = ref(null)
let terminal = null
let fitAddon = null
let websocket = null
let resizeObserver = null

const uiState = computed(() => getKubectlTerminalUiState({
  isConnected: isConnected.value,
  hasTerminalContent: hasTerminalContent.value
}))

const initTerminal = () => {
  if (!terminalElement.value) return

  terminal = new XTerm({
    cursorBlink: true,
    fontSize: 14,
    fontFamily: 'Monaco, Menlo, Consolas, monospace',
    theme: {
      background: '#0b1220',
      foreground: '#dbe7ff',
      cursor: '#67e8f9',
      selection: '#1d4ed880'
    },
    scrollback: 2000,
    convertEol: true
  })

  fitAddon = new FitAddon()
  terminal.loadAddon(fitAddon)
  terminal.open(terminalElement.value)
  fitAddon.fit()

  terminal.onData((data) => {
    if (websocket && websocket.readyState === WebSocket.OPEN) {
      websocket.send(JSON.stringify({
        operation: 'stdin',
        data
      }))
    }
  })

  resizeObserver = new ResizeObserver(() => {
    if (!fitAddon) return
    fitAddon.fit()
    if (websocket && websocket.readyState === WebSocket.OPEN) {
      websocket.send(JSON.stringify({
        operation: 'resize',
        data: {
          cols: terminal.cols,
          rows: terminal.rows
        }
      }))
    }
  })
  resizeObserver.observe(terminalElement.value)
}

const loadNamespaces = async () => {
  try {
    const response = await k8sApi.getNamespaces(clusterId.value)
    const responseData = response.data || response
    const list = responseData.data || []
    namespaceOptions.value = Array.isArray(list) ? list : []
    if (namespaceOptions.value.length > 0 && !namespaceOptions.value.find(item => item.name === selectedNamespace.value)) {
      selectedNamespace.value = namespaceOptions.value[0].name
    }
  } catch (error) {
    console.error('加载命名空间失败:', error)
    ElMessage.error('加载命名空间失败')
  }
}

const connectTerminal = async () => {
  const token = storage.getItem('token')
  if (!token) {
    ElMessage.error('未找到登录凭证，请重新登录')
    return
  }

  disconnectTerminal()
  connecting.value = true

  try {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const wsUrl = `${protocol}//${window.location.host}/api/v1/k8s/cluster/${clusterId.value}/kubectl/terminal?namespace=${encodeURIComponent(selectedNamespace.value)}&token=${encodeURIComponent(token)}`
    websocket = new WebSocket(wsUrl)

    websocket.onopen = async () => {
      connecting.value = false
      isConnected.value = true
      hasTerminalContent.value = true
      terminal.clear()
      terminal.writeln('\x1B[1;32mOpsNexus kubectl terminal ready\x1B[0m')
      terminal.writeln(`\x1B[1;34mNamespace: ${selectedNamespace.value || 'current-context'}\x1B[0m`)
      terminal.writeln('')
      await nextTick()
      fitAddon?.fit()
      focusTerminal()
    }

    websocket.onmessage = (event) => {
      try {
        const message = JSON.parse(event.data)
        if ((message.operation === 'stdout' || message.operation === 'stderr') && message.data) {
          hasTerminalContent.value = true
          terminal.write(message.data)
        }
      } catch {
        hasTerminalContent.value = true
        terminal.write(event.data)
      }
    }

    websocket.onerror = (error) => {
      console.error('kubectl terminal websocket error:', error)
      ElMessage.error('kubectl 终端连接失败')
      connecting.value = false
      isConnected.value = false
    }

    websocket.onclose = () => {
      connecting.value = false
      isConnected.value = false
    }
  } catch (error) {
    console.error('连接 kubectl 终端失败:', error)
    ElMessage.error('连接 kubectl 终端失败')
    connecting.value = false
    isConnected.value = false
  }
}

const disconnectTerminal = () => {
  if (websocket) {
    websocket.close()
    websocket = null
  }
  isConnected.value = false
  hasTerminalContent.value = false
  terminal?.clear()
}

const executeQuickCommand = async () => {
  if (!quickCommand.value.trim()) {
    ElMessage.warning('请输入 kubectl 命令')
    return
  }

  try {
    executingCommand.value = true
    const response = await k8sApi.executeKubectl(clusterId.value, {
      command: quickCommand.value.trim(),
      namespace: selectedNamespace.value,
      timeoutSeconds: 60
    })
    const responseData = response.data || response
    const payload = responseData.data || {}

    if (responseData.code !== 200 && !responseData.success) {
      throw new Error(responseData.message || 'kubectl 执行失败')
    }

    const output = [
      `$ kubectl ${quickCommand.value.trim()}`,
      payload.stdout || '',
      payload.stderr || ''
    ].filter(Boolean).join('\r\n')

    hasTerminalContent.value = true
    terminal.writeln(output)
    terminal.writeln('')
    quickCommand.value = ''
  } catch (error) {
    console.error('kubectl 命令执行失败:', error)
    ElMessage.error(error.message || 'kubectl 命令执行失败')
  } finally {
    executingCommand.value = false
  }
}

const handleNamespaceChange = () => {
  if (isConnected.value) {
    connectTerminal()
  }
}

const focusTerminal = () => {
  terminal?.focus()
}

onMounted(async () => {
  initTerminal()
  await loadNamespaces()
})

onUnmounted(() => {
  disconnectTerminal()
  resizeObserver?.disconnect()
  terminal?.dispose()
})
</script>

<style scoped>
.kubectl-terminal-page {
  padding: 20px;
  min-height: 100vh;
  background: linear-gradient(135deg, #e2e8f0 0%, #f8fafc 100%);
}

.kubectl-terminal-card {
  border-radius: 16px;
}

.terminal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
}

.terminal-title {
  display: flex;
  align-items: center;
  gap: 10px;
  font-weight: 600;
}

.terminal-controls {
  display: flex;
  align-items: center;
  gap: 12px;
}

.command-toolbar {
  display: grid;
  grid-template-columns: 1fr auto;
  gap: 12px;
  margin-bottom: 16px;
}

.terminal-wrapper {
  position: relative;
}

.terminal-surface {
  height: 620px;
  border-radius: 12px;
  overflow: hidden;
  background: #0b1220;
  padding: 12px;
}

.terminal-placeholder {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #cbd5e1;
  pointer-events: none;
}

@media (max-width: 960px) {
  .terminal-header {
    flex-direction: column;
    align-items: flex-start;
  }

  .terminal-controls {
    width: 100%;
    flex-wrap: wrap;
  }

  .command-toolbar {
    grid-template-columns: 1fr;
  }

  .terminal-surface {
    height: 520px;
  }
}
</style>
