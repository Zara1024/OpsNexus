import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('../src/views/ai/AIAssistant.vue', import.meta.url), 'utf8')

function testAssistantCopyUsesHostIpExamples() {
  assert.match(
    source,
    /查询主机 10\.0\.0\.200/,
    'expected AI assistant examples to keep the host IP lookup prompt'
  )

  assert.match(
    source,
    /查看主机 10\.0\.0\.200 的磁盘占用/,
    'expected AI assistant examples to show disk lookup with a host IP'
  )

  assert.match(
    source,
    /为主机 10\.0\.0\.200 生成巡检报告/,
    'expected AI assistant examples to show inspection generation with a host IP'
  )
}

function testAssistantCopyStopsAdvertisingHostIds() {
  const forbiddenPhrases = [
    '主机 ID',
    '主机 12'
  ]

  for (const phrase of forbiddenPhrases) {
    assert.ok(
      !source.includes(phrase),
      `expected AI assistant copy to omit ${phrase}`
    )
  }
}

async function main() {
  testAssistantCopyUsesHostIpExamples()
  testAssistantCopyStopsAdvertisingHostIds()
  console.log('AI assistant host copy tests passed')
}

main().catch((error) => {
  console.error(error)
  process.exitCode = 1
})
