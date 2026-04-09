# AlertManager Integration Baseline

This repository already contains the P1 alert-center application code.
The remaining deployment gap is wiring Prometheus and AlertManager together in a way that works for Docker, Kubernetes, and future Helm packaging.

## Docker

The optional Docker Compose profile is `p1-alerting`.

```bash
cd docker
docker compose --profile p1-alerting up -d alertmanager prometheus
```

Relevant files:

- `docker/alertmanager/alertmanager.yml`
- `docker/prometheus/prometheus.yml`
- `docker/docker-compose.yml`
- `docker/.env`

## Kubernetes

Apply the optional AlertManager manifests before or together with Prometheus:

```bash
kubectl apply -f k8s/alertmanager/
kubectl apply -f k8s/prometheus/
```

Relevant files:

- `k8s/alertmanager/alertmanager-cm0-configmap.yaml`
- `k8s/alertmanager/alertmanager-deployment.yaml`
- `k8s/alertmanager/alertmanager-service.yaml`
- `k8s/prometheus/prometheus-cm0-configmap.yaml`

## Helm

The repository still does not include a full Helm chart.
For now, `deploy/helm/values-alerting.example.yaml` defines the minimum values contract that a future chart or umbrella chart should honor:

- `api.env.MONITOR_*`
- `prometheus.alerting.target`
- `alertmanager.enabled`
- `alertmanager.service.port`
- `alertmanager.config`

## Runtime validation checklist

1. Prometheus can resolve `alertmanager:9093`.
2. OpsNexus `monitor_alert_source` contains one enabled source of type `4`.
3. `GET /api/v1/monitor/alerts/alertmanager/status` returns `code=200`.
4. `GET /api/v1/monitor/alerts/alertmanager/receivers` returns at least one receiver.
5. `POST /api/v1/monitor/alerts/webhook` creates rows in:
   - `monitor_webhook_log`
   - `monitor_webhook_notify_log`
