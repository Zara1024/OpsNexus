import assert from 'node:assert/strict'

import { getKubectlTerminalUiState } from '../src/utils/kubectlTerminalPresentation.mjs'

function testDisconnectedEmptyTerminalShowsPlaceholderAndDisablesDisconnect() {
  const state = getKubectlTerminalUiState({
    isConnected: false,
    hasTerminalContent: false
  })

  assert.equal(state.showPlaceholder, true)
  assert.equal(state.canDisconnect, false)
}

function testConnectedTerminalHidesPlaceholderAndEnablesDisconnect() {
  const state = getKubectlTerminalUiState({
    isConnected: true,
    hasTerminalContent: true
  })

  assert.equal(state.showPlaceholder, false)
  assert.equal(state.canDisconnect, true)
}

function testQuickCommandOutputHidesPlaceholderAndEnablesDisconnectWithoutWebsocket() {
  const state = getKubectlTerminalUiState({
    isConnected: false,
    hasTerminalContent: true
  })

  assert.equal(state.showPlaceholder, false)
  assert.equal(state.canDisconnect, true)
}

async function main() {
  testDisconnectedEmptyTerminalShowsPlaceholderAndDisablesDisconnect()
  testConnectedTerminalHidesPlaceholderAndEnablesDisconnect()
  testQuickCommandOutputHidesPlaceholderAndEnablesDisconnectWithoutWebsocket()
  console.log('kubectl terminal ui state tests passed')
}

main().catch((error) => {
  console.error(error)
  process.exitCode = 1
})
