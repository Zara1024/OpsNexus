# Docker Compose Runbook

`docker/` is the only production deployment entrypoint for OpsNexus.

## Files

```text
docker/
├── docker-compose.yml
├── .env.example
├── api/config.yaml
├── web/devops.conf
├── mysql/
├── redis/
├── prometheus/
├── pushgateway/
└── alertmanager/
```

## First Start

```bash
cd docker
cp .env.example .env
vi .env
docker compose up -d --build
```

Optional Alertmanager:

```bash
docker compose --profile alerting up -d --build
```

## Verify

```bash
docker compose ps
curl -I http://127.0.0.1:${WEB_PORT:-8080}/
curl -fsS http://127.0.0.1:${WEB_PORT:-8080}/api/v1/captcha
```

Runtime data is stored under `docker/mysql/data`, `docker/redis/data`,
`docker/prometheus/data`, `docker/pushgateway/data`, `docker/alertmanager/data`,
`docker/api/upload`, and `docker/api/logs`.

Do not commit `docker/.env`, SSH keys, uploaded files, database data, or logs.
