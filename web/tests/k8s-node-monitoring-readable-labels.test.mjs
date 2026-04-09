import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('../src/views/K8s/nodes/NodesMonitoring.vue', import.meta.url), 'utf8')

function testMonitoringDashboardShowsReadablePrimaryLabels() {
  for (const label of [
    'CPU使用率',
    '内存使用率',
    'Pod数量',
    '更新时间',
    '详细信息',
    '刷新'
  ]) {
    assert.match(source, new RegExp(label), `expected NodesMonitoring.vue to include readable label: ${label}`)
  }
}

function testNodeDetailDialogShowsReadableSectionLabels() {
  for (const label of [
    '节点监控',
    '资源使用情况',
    '系统信息',
    '操作系统',
    '内核版本',
    '容器运行时',
    '运行的Pod',
    '暂无运行中的Pod'
  ]) {
    assert.match(source, new RegExp(label), `expected NodesMonitoring.vue to include readable detail label: ${label}`)
  }
}

function testStatusAndToastMessagesAreReadable() {
  for (const label of [
    '就绪',
    '未就绪',
    '调度禁用',
    '节点',
    '监控数据已刷新',
    '获取节点'
  ]) {
    assert.match(source, new RegExp(label), `expected NodesMonitoring.vue to include readable status/message text: ${label}`)
  }
}

async function main() {
  testMonitoringDashboardShowsReadablePrimaryLabels()
  testNodeDetailDialogShowsReadableSectionLabels()
  testStatusAndToastMessagesAreReadable()
  console.log('k8s node monitoring readable labels tests passed')
}

main().catch((error) => {
  console.error(error)
  process.exitCode = 1
})
