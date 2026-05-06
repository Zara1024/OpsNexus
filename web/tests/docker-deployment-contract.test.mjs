import assert from 'node:assert/strict'
import fs from 'node:fs'

const composeSource = fs.readFileSync(new URL('../../docker/docker-compose.yml', import.meta.url), 'utf8')
const envExampleSource = fs.readFileSync(new URL('../../docker/.env.example', import.meta.url), 'utf8')
const apiDockerfileSource = fs.readFileSync(new URL('../../docker/api/Dockerfile', import.meta.url), 'utf8')
const webDockerfileSource = fs.readFileSync(new URL('../../docker/web/Dockerfile', import.meta.url), 'utf8')
const mysqlInitSource = fs.readFileSync(new URL('../../docker/mysql/init.sql', import.meta.url), 'utf8')
const readmeSource = fs.readFileSync(new URL('../../README.md', import.meta.url), 'utf8')

function testComposeBuildsApiAndWebFromRepoSource() {
  assert.match(
    composeSource,
    /devops-api:[\s\S]*build:\s*[\s\S]*context:\s*\.\.[\s\S]*dockerfile:\s*docker\/api\/Dockerfile/s,
    'expected docker compose to build devops-api from the repository source'
  )

  assert.match(
    composeSource,
    /devops-web:[\s\S]*build:\s*[\s\S]*context:\s*\.\.[\s\S]*dockerfile:\s*docker\/web\/Dockerfile/s,
    'expected docker compose to build devops-web from the repository source'
  )
}

function testComposeAllowsRemoteRootConnectionsWithoutPasswordDrift() {
  assert.match(
    composeSource,
    /MYSQL_ROOT_HOST:\s*\$\{MYSQL_ROOT_HOST:-%\}/,
    'expected mysql service to expose MYSQL_ROOT_HOST so root can connect from other containers'
  )

  assert.ok(
    !mysqlInitSource.includes('devops@2025'),
    'expected mysql init script to avoid hardcoded root passwords'
  )
}

function testComposeRequiresProductionSecretsFromEnvFile() {
  for (const key of [
    'MYSQL_ROOT_PASSWORD',
    'DB_PASSWORD',
    'REDIS_PASSWORD',
    'MONITOR_AGENT_HEARTBEAT_TOKEN',
    'MONITOR_WEBHOOK_TOKEN'
  ]) {
    assert.match(
      composeSource,
      new RegExp(`\\$\\{${key}:\\?set ${key} in docker/\\.env\\}`),
      `expected docker compose to require ${key} from docker/.env`
    )
  }

  assert.ok(
    !composeSource.includes('CHANGE_ME') && !envExampleSource.includes('CHANGE_ME'),
    'expected docker deployment files to avoid CHANGE_ME as a pseudo-secret'
  )
}

function testApiDockerImageIncludesSourceBuildAndTemplates() {
  assert.match(
    apiDockerfileSource,
    /COPY api\/go\.mod api\/go\.sum \.\/api\//,
    'expected api Dockerfile to copy Go module files from the repository source'
  )

  assert.match(
    apiDockerfileSource,
    /go build\s+-o\s+.*\/devops/s,
    'expected api Dockerfile to build the Go binary during image creation'
  )

  assert.match(
    apiDockerfileSource,
    /COPY --from=.*\/src\/api\/common\/templates .*\/app\/common\/templates/s,
    'expected api Dockerfile to copy common templates into the runtime image'
  )

  assert.match(
    apiDockerfileSource,
    /COPY --from=.*\/src\/api\/upload\/xlsl\/host\.xlsx .*\/app\/upload\/xlsl\/host\.xlsx/s,
    'expected api Dockerfile to include the CMDB host import template'
  )
}

function testWebDockerImageBuildsDistFromSource() {
  assert.match(
    webDockerfileSource,
    /COPY web\/package\.json web\/package-lock\.json \.\/web\//,
    'expected web Dockerfile to copy frontend package manifests from the repository source'
  )

  assert.match(
    webDockerfileSource,
    /npm (ci|install)[\s\S]*npm run build/s,
    'expected web Dockerfile to build the frontend dist during image creation'
  )

  assert.match(
    webDockerfileSource,
    /COPY --from=.*\/src\/web\/dist \/usr\/share\/nginx\/html/s,
    'expected web Dockerfile to copy built dist files into the runtime image'
  )
}

function testWebHealthcheckUsesIpv4Loopback() {
  assert.match(
    composeSource,
    /devops-web:[\s\S]*healthcheck:[\s\S]*http:\/\/127\.0\.0\.1\/healthz/s,
    'expected devops-web healthcheck to probe 127.0.0.1/healthz for reliable loopback checks'
  )
}

function testReadmeIsDockerOnlyForProductionDeployment() {
  assert.match(
    readmeSource,
    /仅支持 Docker Compose 部署/,
    'expected README to declare Docker Compose as the only production deployment mode'
  )

  assert.ok(
    !readmeSource.match(/deploy\/helm|deploy\/systemd|helm upgrade|kubectl apply/),
    'expected README to avoid old Helm, systemd, and raw Kubernetes deployment instructions'
  )
}

async function main() {
  testComposeBuildsApiAndWebFromRepoSource()
  testComposeAllowsRemoteRootConnectionsWithoutPasswordDrift()
  testComposeRequiresProductionSecretsFromEnvFile()
  testApiDockerImageIncludesSourceBuildAndTemplates()
  testWebDockerImageBuildsDistFromSource()
  testWebHealthcheckUsesIpv4Loopback()
  testReadmeIsDockerOnlyForProductionDeployment()
  console.log('docker deployment contract tests passed')
}

main().catch((error) => {
  console.error(error)
  process.exitCode = 1
})
