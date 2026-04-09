import assert from 'node:assert/strict'
import fs from 'node:fs'
import path from 'node:path'
import { fileURLToPath } from 'node:url'

const __filename = fileURLToPath(import.meta.url)
const __dirname = path.dirname(__filename)
const repoRoot = path.resolve(__dirname, '..', '..')
const readmePath = path.join(repoRoot, 'README.md')
const helmReadmePath = path.join(repoRoot, 'deploy', 'helm', 'README.md')
const exampleValuesPath = path.join(repoRoot, 'deploy', 'helm', 'values-current-source.example.yaml')

const readmeSource = fs.readFileSync(readmePath, 'utf8')
const helmReadmeSource = fs.readFileSync(helmReadmePath, 'utf8')

function testCurrentSourceExampleValuesFileExists() {
  assert.ok(
    fs.existsSync(exampleValuesPath),
    'expected deploy/helm/values-current-source.example.yaml to exist so new machines have a ready-to-copy current-source Helm override'
  )
}

function testRootReadmePointsToCurrentSourceExample() {
  assert.match(
    readmeSource,
    /deploy\/helm\/values-current-source\.example\.yaml/,
    'expected the root README to point new-machine Helm deployments at the current-source example values file'
  )
}

function testHelmReadmeDocumentsCurrentSourceExample() {
  assert.match(
    helmReadmeSource,
    /values-current-source\.example\.yaml/,
    'expected deploy/helm/README.md to document the current-source image override example'
  )
}

async function main() {
  testCurrentSourceExampleValuesFileExists()
  testRootReadmePointsToCurrentSourceExample()
  testHelmReadmeDocumentsCurrentSourceExample()
  console.log('helm current source values example documentation tests passed')
}

main().catch((error) => {
  console.error(error)
  process.exitCode = 1
})
