{{/*
Expand the name of the chart.
*/}}
{{- define "opsnexus.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
*/}}
{{- define "opsnexus.fullname" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
Create chart label.
*/}}
{{- define "opsnexus.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Common labels.
*/}}
{{- define "opsnexus.labels" -}}
helm.sh/chart: {{ include "opsnexus.chart" . }}
app.kubernetes.io/name: {{ include "opsnexus.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- with .Values.global.extraLabels }}
{{ toYaml . }}
{{- end }}
{{- end -}}

{{/*
Selector labels.
*/}}
{{- define "opsnexus.selectorLabels" -}}
app.kubernetes.io/name: {{ include "opsnexus.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}

{{/*
Component fullname.
*/}}
{{- define "opsnexus.componentFullname" -}}
{{- printf "%s-%s" (include "opsnexus.fullname" .root) .component | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Component labels.
*/}}
{{- define "opsnexus.componentLabels" -}}
{{ include "opsnexus.labels" .root }}
app.kubernetes.io/component: {{ .component }}
{{- end -}}

{{/*
Component selector labels.
*/}}
{{- define "opsnexus.componentSelectorLabels" -}}
{{ include "opsnexus.selectorLabels" .root }}
app.kubernetes.io/component: {{ .component }}
{{- end -}}

{{/*
Resolve service account name.
*/}}
{{- define "opsnexus.serviceAccountName" -}}
{{- if .Values.serviceAccount.create -}}
{{- default (include "opsnexus.fullname" .) .Values.serviceAccount.name -}}
{{- else -}}
{{- default "default" .Values.serviceAccount.name -}}
{{- end -}}
{{- end -}}

{{/*
Resolve image pull policy.
*/}}
{{- define "opsnexus.imagePullPolicy" -}}
{{- default .default .value -}}
{{- end -}}

{{/*
Resolve persistence storageClass.
*/}}
{{- define "opsnexus.storageClass" -}}
{{- $componentStorageClass := .componentStorageClass | default "" -}}
{{- if $componentStorageClass -}}
storageClassName: {{ $componentStorageClass | quote }}
{{- else if .root.Values.global.storageClass -}}
storageClassName: {{ .root.Values.global.storageClass | quote }}
{{- end -}}
{{- end -}}

{{/*
Resolve mysql secret name.
*/}}
{{- define "opsnexus.mysqlSecretName" -}}
{{- if .Values.mysql.auth.existingSecret -}}
{{- .Values.mysql.auth.existingSecret -}}
{{- else if .Values.externalServices.mysql.existingSecret -}}
{{- .Values.externalServices.mysql.existingSecret -}}
{{- else -}}
{{- printf "%s-mysql" (include "opsnexus.fullname" .) | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}

{{/*
Resolve redis secret name.
*/}}
{{- define "opsnexus.redisSecretName" -}}
{{- if .Values.redis.auth.existingSecret -}}
{{- .Values.redis.auth.existingSecret -}}
{{- else if .Values.externalServices.redis.existingSecret -}}
{{- .Values.externalServices.redis.existingSecret -}}
{{- else -}}
{{- printf "%s-redis" (include "opsnexus.fullname" .) | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}

{{/*
Resolve api secret name.
*/}}
{{- define "opsnexus.apiSecretName" -}}
{{- if .Values.api.existingSecret -}}
{{- .Values.api.existingSecret -}}
{{- else -}}
{{- printf "%s-api" (include "opsnexus.fullname" .) | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}

{{/*
Resolve API database host.
*/}}
{{- define "opsnexus.apiDbHost" -}}
{{- if .Values.mysql.enabled -}}
{{ include "opsnexus.componentFullname" (dict "root" . "component" "mysql") }}
{{- else -}}
{{ required "externalServices.mysql.host is required when mysql.enabled=false" .Values.externalServices.mysql.host }}
{{- end -}}
{{- end -}}

{{/*
Resolve API database port.
*/}}
{{- define "opsnexus.apiDbPort" -}}
{{- if .Values.mysql.enabled -}}
{{ .Values.mysql.service.port }}
{{- else -}}
{{ .Values.externalServices.mysql.port }}
{{- end -}}
{{- end -}}

{{/*
Resolve API database name.
*/}}
{{- define "opsnexus.apiDbName" -}}
{{- if .Values.mysql.enabled -}}
{{ .Values.mysql.auth.database }}
{{- else -}}
{{ .Values.externalServices.mysql.database }}
{{- end -}}
{{- end -}}

{{/*
Resolve API database username.
*/}}
{{- define "opsnexus.apiDbUser" -}}
{{- if .Values.mysql.enabled -}}
{{ .Values.mysql.auth.username }}
{{- else -}}
{{ .Values.externalServices.mysql.username }}
{{- end -}}
{{- end -}}

{{/*
Resolve API database password secret key.
*/}}
{{- define "opsnexus.apiDbPasswordSecretKey" -}}
{{- if .Values.mysql.enabled -}}
{{ .Values.mysql.auth.rootPasswordKey }}
{{- else -}}
{{ .Values.externalServices.mysql.passwordKey }}
{{- end -}}
{{- end -}}

{{/*
Resolve API redis address.
*/}}
{{- define "opsnexus.apiRedisAddress" -}}
{{- if .Values.redis.enabled -}}
{{ printf "%s:%v" (include "opsnexus.componentFullname" (dict "root" . "component" "redis")) .Values.redis.service.port }}
{{- else -}}
{{ required "externalServices.redis.address is required when redis.enabled=false" .Values.externalServices.redis.address }}
{{- end -}}
{{- end -}}

{{/*
Resolve API redis password secret key.
*/}}
{{- define "opsnexus.apiRedisPasswordSecretKey" -}}
{{- if .Values.redis.enabled -}}
{{ .Values.redis.auth.passwordKey }}
{{- else -}}
{{ .Values.externalServices.redis.passwordKey }}
{{- end -}}
{{- end -}}

{{/*
Resolve API Prometheus URL.
*/}}
{{- define "opsnexus.apiPrometheusURL" -}}
{{- if .Values.api.monitor.prometheusURL -}}
{{ .Values.api.monitor.prometheusURL }}
{{- else if .Values.prometheus.enabled -}}
{{ printf "http://%s:%v" (include "opsnexus.componentFullname" (dict "root" . "component" "prometheus")) .Values.prometheus.service.port }}
{{- else -}}
http://prometheus:9090
{{- end -}}
{{- end -}}

{{/*
Resolve API Pushgateway URL.
*/}}
{{- define "opsnexus.apiPushgatewayURL" -}}
{{- if .Values.api.monitor.pushgatewayURL -}}
{{ .Values.api.monitor.pushgatewayURL }}
{{- else if .Values.pushgateway.enabled -}}
{{ printf "http://%s:%v" (include "opsnexus.componentFullname" (dict "root" . "component" "pushgateway")) .Values.pushgateway.service.port }}
{{- else -}}
http://pushgateway:9091
{{- end -}}
{{- end -}}

{{/*
Resolve API heartbeat URL.
*/}}
{{- define "opsnexus.apiHeartbeatServerURL" -}}
{{- if .Values.api.monitor.agentHeartbeatServerURL -}}
{{ .Values.api.monitor.agentHeartbeatServerURL }}
{{- else -}}
{{ printf "http://%s:%v/api/v1/monitor/agent/heartbeat" (include "opsnexus.componentFullname" (dict "root" . "component" "api")) .Values.api.service.port }}
{{- end -}}
{{- end -}}

{{/*
Resolve Alertmanager target for Prometheus.
*/}}
{{- define "opsnexus.prometheusAlertmanagerTarget" -}}
{{- if .Values.prometheus.alerting.target -}}
{{ .Values.prometheus.alerting.target }}
{{- else -}}
{{ printf "%s:%v" (include "opsnexus.componentFullname" (dict "root" . "component" "alertmanager")) .Values.alertmanager.service.port }}
{{- end -}}
{{- end -}}
