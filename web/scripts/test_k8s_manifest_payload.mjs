import assert from 'node:assert/strict'

import {
  buildConfigMapCreatePayload,
  buildSecretCreatePayload
} from '../src/utils/k8sManifest.mjs'

const configMapYaml = `apiVersion: v1
kind: ConfigMap
metadata:
  name: demo-config
  labels:
    app: demo
immutable: true
data:
  app.yaml: |
    app:
      name: demo
`

const configMapPayload = buildConfigMapCreatePayload(configMapYaml)

assert.deepEqual(configMapPayload, {
  name: 'demo-config',
  labels: {
    app: 'demo'
  },
  data: {
    'app.yaml': 'app:\n  name: demo\n'
  },
  binaryData: {},
  immutable: true
})

const secretYaml = `apiVersion: v1
kind: Secret
metadata:
  name: demo-secret
  labels:
    tier: backend
type: Opaque
data:
  username: YWRtaW4=
stringData:
  password: password
`

const secretPayload = buildSecretCreatePayload(secretYaml)

assert.deepEqual(secretPayload, {
  name: 'demo-secret',
  labels: {
    tier: 'backend'
  },
  type: 'Opaque',
  data: {
    username: 'YWRtaW4='
  },
  stringData: {
    password: 'password'
  },
  immutable: false
})

assert.throws(
  () => buildConfigMapCreatePayload('apiVersion: v1\nkind: ConfigMap\nmetadata: {}\n'),
  /metadata\.name/
)

console.log('k8s manifest payload tests passed')
