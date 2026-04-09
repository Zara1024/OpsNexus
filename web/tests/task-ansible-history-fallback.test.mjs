import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('../src/views/task/AnsibleTaskHistory.vue', import.meta.url), 'utf8')

function testHistoryPageFallsBackToTaskDetailWhenHistoryApiMissing() {
  assert.match(
    source,
    /const buildFallbackHistoryRecord = taskInfo => \{[\s\S]*Works:\s*works[\s\S]*\}/s,
    "expected the ansible history page to define a fallback record builder that preserves the task's Works list"
  )

  assert.match(
    source,
    /const fetchHistory = async \(\) => \{[\s\S]*const res = await GetAnsibleTaskHistory\(queryParams\)[\s\S]*catch \(error\) \{[\s\S]*error\?\.response\?\.status === 404[\s\S]*GetAnsibleTaskDetail\(taskId\.value\)[\s\S]*historyList\.value = matchesStatus \? \[fallbackRecord\] : \[\][\s\S]*return[\s\S]*ElMessage\.error\('获取任务历史失败'\)/s,
    'expected fetchHistory to fall back to current task detail when the ansible history API returns 404, instead of immediately treating the page as a hard failure'
  )
}

async function main() {
  testHistoryPageFallsBackToTaskDetailWhenHistoryApiMissing()
  console.log('task ansible history fallback tests passed')
}

main().catch((error) => {
  console.error(error)
  process.exitCode = 1
})
