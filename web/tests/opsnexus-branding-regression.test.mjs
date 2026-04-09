import assert from 'node:assert/strict'
import fs from 'node:fs'

const webIndexSource = fs.readFileSync(new URL('../public/index.html', import.meta.url), 'utf8')
const apiMainSource = fs.readFileSync(new URL('../../api/main.go', import.meta.url), 'utf8')
const swaggerDocsGoSource = fs.readFileSync(new URL('../../api/docs/docs.go', import.meta.url), 'utf8')
const swaggerJsonSource = fs.readFileSync(new URL('../../api/docs/swagger.json', import.meta.url), 'utf8')
const swaggerYamlSource = fs.readFileSync(new URL('../../api/docs/swagger.yaml', import.meta.url), 'utf8')

function testHtmlTitleUsesOpsNexusBrand() {
  assert.match(
    webIndexSource,
    /<title>OpsNexus<\/title>/,
    'expected the built page title source to use OpsNexus'
  )

  assert.doesNotMatch(
    webIndexSource,
    /devops运维管理系统/,
    'expected the public page shell to stop shipping the legacy devops title'
  )
}

function testSwaggerAnnotationsUseOpsNexusBrand() {
  assert.match(
    apiMainSource,
    /@title OpsNexus/,
    'expected Swagger annotations to use OpsNexus as the API title'
  )

  assert.match(
    apiMainSource,
    /@description OpsNexus API documentation\./,
    'expected Swagger annotations to describe the API with OpsNexus branding'
  )
}

function testGeneratedSwaggerArtifactsUseOpsNexusBrand() {
  for (const [label, source] of [
    ['docs.go', swaggerDocsGoSource],
    ['swagger.json', swaggerJsonSource],
    ['swagger.yaml', swaggerYamlSource]
  ]) {
    assert.match(
      source,
      /OpsNexus/,
      `expected ${label} to contain the OpsNexus brand`
    )

    assert.doesNotMatch(
      source,
      /devops运维管理系统/,
      `expected ${label} to stop containing the legacy devops brand`
    )
  }
}

async function main() {
  testHtmlTitleUsesOpsNexusBrand()
  testSwaggerAnnotationsUseOpsNexusBrand()
  testGeneratedSwaggerArtifactsUseOpsNexusBrand()
  console.log('opsnexus branding regression tests passed')
}

main().catch((error) => {
  console.error(error)
  process.exitCode = 1
})
