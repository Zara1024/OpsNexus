export function getKubectlTerminalUiState({
  isConnected = false,
  hasTerminalContent = false
} = {}) {
  const connected = Boolean(isConnected)
  const contentReady = Boolean(hasTerminalContent)

  return {
    showPlaceholder: !connected && !contentReady,
    canDisconnect: connected || contentReady
  }
}
