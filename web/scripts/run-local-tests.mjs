import fs from 'node:fs'
import path from 'node:path'
import process from 'node:process'
import { spawnSync } from 'node:child_process'
import { fileURLToPath } from 'node:url'

const __filename = fileURLToPath(import.meta.url)
const __dirname = path.dirname(__filename)
const projectRoot = path.resolve(__dirname, '..')
const workspaceRoot = path.resolve(projectRoot, '..')
const testsDir = path.join(projectRoot, 'tests')
const nodeTestFile = path.join(projectRoot, 'src', 'utils', 'cmdbHostPhase1.test.mjs')

function runNode(args) {
  const result = spawnSync(process.execPath, args, {
    cwd: workspaceRoot,
    stdio: 'inherit'
  })

  if (result.status !== 0) {
    process.exit(result.status ?? 1)
  }
}

function getStandaloneTests() {
  return fs
    .readdirSync(testsDir, { withFileTypes: true })
    .filter((entry) => entry.isFile() && entry.name.endsWith('.test.mjs'))
    .map((entry) => path.join(testsDir, entry.name))
    .sort((left, right) => left.localeCompare(right))
}

const standaloneTests = getStandaloneTests()

for (const testFile of standaloneTests) {
  console.log(`\n[run-local-tests] ${path.relative(projectRoot, testFile)}`)
  runNode([testFile])
}

console.log(`\n[run-local-tests] ${path.relative(projectRoot, nodeTestFile)}`)
runNode(['--test', nodeTestFile])

console.log(`\n[run-local-tests] Completed ${standaloneTests.length + 1} test files`)
