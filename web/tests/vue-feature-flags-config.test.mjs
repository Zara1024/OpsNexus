import assert from 'node:assert/strict'
import path from 'node:path'
import { pathToFileURL } from 'node:url'

async function loadWebpackConfig() {
  const configPath = path.resolve(process.cwd(), 'web', 'vue.config.js')
  const module = await import(pathToFileURL(configPath))
  return module.default || module
}

async function testHydrationMismatchFeatureFlagIsDefined() {
  const config = await loadWebpackConfig()
  const plugins = config.configureWebpack?.plugins || []
  const definePlugin = plugins.find((plugin) => plugin?.definitions?.__VUE_PROD_HYDRATION_MISMATCH_DETAILS__ !== undefined)

  assert.ok(
    definePlugin,
    'expected vue.config.js to define __VUE_PROD_HYDRATION_MISMATCH_DETAILS__ via webpack DefinePlugin'
  )

  assert.equal(
    definePlugin.definitions.__VUE_PROD_HYDRATION_MISMATCH_DETAILS__,
    'false',
    'expected __VUE_PROD_HYDRATION_MISMATCH_DETAILS__ to be disabled explicitly'
  )
}

async function main() {
  await testHydrationMismatchFeatureFlagIsDefined()
  console.log('vue feature flags config tests passed')
}

main().catch((error) => {
  console.error(error)
  process.exitCode = 1
})
