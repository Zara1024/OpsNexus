<template>
  <el-row :gutter="12">
    <el-col :span="16">
      <el-card>
        <template #header>
          <div class="toolbar">
            <span>Web 终端</span>
            <div>
              <el-input-number v-model="hostID" :min="1" controls-position="right" style="width: 140px; margin-right: 8px" />
              <el-button v-permission="'cmdb:host:terminal'" type="primary" @click="connect">连接</el-button>
              <el-button @click="disconnect">断开</el-button>
            </div>
          </div>
        </template>

        <div ref="terminalBox" class="terminal-output">{{ output }}</div>

        <div class="input-row">
          <el-input v-model="command" placeholder="输入命令并回车" @keyup.enter="sendCommand" />
          <el-button type="primary" @click="sendCommand">发送</el-button>
        </div>
      </el-card>
    </el-col>

    <el-col :span="8">
      <el-card>
        <template #header><span>SSH 审计记录</span></template>
        <el-table :data="records" size="small" height="520">
          <el-table-column prop="host_id" label="主机" width="72" />
          <el-table-column prop="session_id" label="会话ID" min-width="140" show-overflow-tooltip />
          <el-table-column prop="start_time" label="开始时间" min-width="160" />
        </el-table>
        <div class="pager">
          <el-pagination
            background
            small
            layout="prev, pager, next"
            :total="total"
            :current-page="query.page"
            :page-size="query.page_size"
            @current-change="handlePageChange"
          />
        </div>
      </el-card>
    </el-col>
  </el-row>
</template>

<script setup lang="ts">
import { nextTick, onBeforeUnmount, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { listCMDBSSHRecordsApi, type CMDBSSHRecordItem } from '@/api/cmdb'

const route = useRoute()
const terminalBox = ref<HTMLElement>()
const hostID = ref<number>(Number(route.query.host_id || 1))
const command = ref('')
const output = ref('')
const ws = ref<WebSocket | null>(null)

const records = ref<CMDBSSHRecordItem[]>([])
const total = ref(0)
const query = ref({ page: 1, page_size: 10 })

const appendOutput = async (text: string) => {
  output.value += text
  await nextTick()
  if (terminalBox.value) {
    terminalBox.value.scrollTop = terminalBox.value.scrollHeight
  }
}

const connect = () => {
  disconnect()
  const token = localStorage.getItem('opsnexus_access_token') || ''
  if (!token) {
    ElMessage.warning('请先登录')
    return
  }
  const protocol = location.protocol === 'https:' ? 'wss' : 'ws'
  const wsURL = `${protocol}://${location.host}/ws/cmdb/terminal/${hostID.value}?token=${encodeURIComponent(token)}`
  ws.value = new WebSocket(wsURL)

  ws.value.onopen = () => appendOutput(`\r\n[系统] 已连接主机 ${hostID.value}\r\n`)
  ws.value.onmessage = (ev) => appendOutput(String(ev.data || ''))
  ws.value.onerror = () => appendOutput('\r\n[系统] 连接异常\r\n')
  ws.value.onclose = () => appendOutput('\r\n[系统] 连接关闭\r\n')
}

const disconnect = () => {
  if (ws.value) {
    ws.value.close()
    ws.value = null
  }
}

const sendCommand = () => {
  if (!ws.value || ws.value.readyState !== WebSocket.OPEN) {
    ElMessage.warning('终端未连接')
    return
  }
  const text = `${command.value}\n`
  ws.value.send(text)
  command.value = ''
}

const loadRecords = async () => {
  const res = await listCMDBSSHRecordsApi(query.value)
  records.value = res.list || []
  total.value = res.total || 0
}

const handlePageChange = (page: number) => {
  query.value.page = page
  loadRecords()
}

onMounted(loadRecords)
onBeforeUnmount(disconnect)
</script>

<style scoped>
.toolbar { display: flex; justify-content: space-between; align-items: center; }
.terminal-output {
  height: 480px;
  overflow: auto;
  background: #0b1020;
  color: #c4f1d2;
  padding: 12px;
  border-radius: 8px;
  white-space: pre-wrap;
  font-family: Consolas, monospace;
}
.input-row { margin-top: 10px; display: grid; grid-template-columns: 1fr 80px; gap: 8px; }
.pager { margin-top: 12px; display: flex; justify-content: flex-end; }
</style>
