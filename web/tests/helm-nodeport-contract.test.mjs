import assert from 'node:assert/strict'
import fs from 'node:fs'

const servicesTemplate = fs.readFileSync(
  new URL('../../deploy/helm/opsnexus/templates/services.yaml', import.meta.url),
  'utf8'
)
const chartValues = fs.readFileSync(
  new URL('../../deploy/helm/opsnexus/values.yaml', import.meta.url),
  'utf8'
)
const currentSourceExample = fs.readFileSync(
  new URL('../../deploy/helm/values-current-source.example.yaml', import.meta.url),
  'utf8'
)

function testServicesTemplateSupportsExplicitNodePorts() {
  assert.match(
    servicesTemplate,
    /type:\s*\{\{\s*\.Values\.api\.service\.type\s*\}\}[\s\S]*nodePort:/s,
    'expected the chart service template to support an explicit api.service.nodePort override when using NodePort'
  )

  assert.match(
    servicesTemplate,
    /type:\s*\{\{\s*\.Values\.web\.service\.type\s*\}\}[\s\S]*nodePort:/s,
    'expected the chart service template to support an explicit web.service.nodePort override when using NodePort'
  )
}

function testChartValuesExposeOptionalNodePortFields() {
  assert.match(
    chartValues,
    /api:\s*[\s\S]*service:\s*[\s\S]*nodePort:\s*null/s,
    'expected values.yaml to expose api.service.nodePort as an optional field'
  )

  assert.match(
    chartValues,
    /web:\s*[\s\S]*service:\s*[\s\S]*nodePort:\s*null/s,
    'expected values.yaml to expose web.service.nodePort as an optional field'
  )
}

function testCurrentSourceExampleFixesWebNodePortTo30080() {
  assert.match(
    currentSourceExample,
    /web:\s*[\s\S]*service:\s*[\s\S]*type:\s*NodePort[\s\S]*nodePort:\s*30080/s,
    'expected the current-source Helm example to pin the web NodePort to 30080'
  )
}

async function main() {
  testServicesTemplateSupportsExplicitNodePorts()
  testChartValuesExposeOptionalNodePortFields()
  testCurrentSourceExampleFixesWebNodePortTo30080()
  console.log('helm nodeport contract tests passed')
}

main().catch((error) => {
  console.error(error)
  process.exitCode = 1
})
