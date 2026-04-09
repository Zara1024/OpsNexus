import request from "@/utils/request"

export default {
    getMonitorAlertSummary() {
        return request({
            url: 'monitor/alerts/summary',
            method: 'get'
        })
    },
    queryMonitorAlertIncidentList(params) {
        return request({
            url: 'monitor/alerts/incidents',
            method: 'get',
            params
        })
    },
    queryMonitorAlertWebhookList(params) {
        return request({
            url: 'monitor/alerts/webhooks',
            method: 'get',
            params
        })
    },
    queryMonitorAlertNotifyLogList(params) {
        return request({
            url: 'monitor/alerts/notify/logs',
            method: 'get',
            params
        })
    },
    getMonitorAlertNotifyRobotList() {
        return request({
            url: 'monitor/alerts/notify/robots',
            method: 'get'
        })
    },
    createMonitorAlertNotifyRobot(data) {
        return request({
            url: 'monitor/alerts/notify/robots',
            method: 'post',
            data
        })
    },
    updateMonitorAlertNotifyRobot(id, data) {
        return request({
            url: `monitor/alerts/notify/robots/${id}`,
            method: 'put',
            data
        })
    },
    updateMonitorAlertNotifyRobotStatus(id, status) {
        return request({
            url: `monitor/alerts/notify/robots/${id}/status`,
            method: 'post',
            data: { status }
        })
    },
    testMonitorAlertNotifyRobot(id, data) {
        return request({
            url: `monitor/alerts/notify/robots/${id}/test`,
            method: 'post',
            data
        })
    },
    deleteMonitorAlertNotifyRobot(id) {
        return request({
            url: `monitor/alerts/notify/robots/${id}`,
            method: 'delete'
        })
    },
    getMonitorAlertSourceList() {
        return request({
            url: 'monitor/alerts/sources',
            method: 'get'
        })
    },
    createMonitorAlertSource(data) {
        return request({
            url: 'monitor/alerts/sources',
            method: 'post',
            data
        })
    },
    updateMonitorAlertSource(id, data) {
        return request({
            url: `monitor/alerts/sources/${id}`,
            method: 'put',
            data
        })
    },
    updateMonitorAlertSourceStatus(id, status) {
        return request({
            url: `monitor/alerts/sources/${id}/status`,
            method: 'post',
            data: { status }
        })
    },
    deleteMonitorAlertSource(id) {
        return request({
            url: `monitor/alerts/sources/${id}`,
            method: 'delete'
        })
    },
    getMonitorAlertManagerStatus(params) {
        return request({
            url: 'monitor/alerts/alertmanager/status',
            method: 'get',
            params
        })
    },
    getMonitorAlertManagerSilences(params) {
        return request({
            url: 'monitor/alerts/alertmanager/silences',
            method: 'get',
            params
        })
    },
    createMonitorAlertManagerSilence(data) {
        return request({
            url: 'monitor/alerts/alertmanager/silences',
            method: 'post',
            data
        })
    },
    deleteMonitorAlertManagerSilence(silenceId, params) {
        return request({
            url: `monitor/alerts/alertmanager/silences/${silenceId}`,
            method: 'delete',
            params
        })
    },
    getMonitorAlertManagerReceivers(params) {
        return request({
            url: 'monitor/alerts/alertmanager/receivers',
            method: 'get',
            params
        })
    },
    getMonitorAutomationOverview() {
        return request({
            url: 'monitor/automation/overview',
            method: 'get'
        })
    },
    getMonitorAutomationEvents(params) {
        return request({
            url: 'monitor/automation/events',
            method: 'get',
            params
        })
    },
    getMonitorHostAlertTemplates() {
        return request({
            url: 'monitor/automation/host-alert/templates',
            method: 'get'
        })
    },
    getMonitorHostAlertRules() {
        return request({
            url: 'monitor/automation/host-alert/rules',
            method: 'get'
        })
    },
    createMonitorHostAlertRule(data) {
        return request({
            url: 'monitor/automation/host-alert/rules',
            method: 'post',
            data
        })
    },
    updateMonitorHostAlertRule(id, data) {
        return request({
            url: `monitor/automation/host-alert/rules/${id}`,
            method: 'put',
            data
        })
    },
    updateMonitorHostAlertRuleStatus(id, status) {
        return request({
            url: `monitor/automation/host-alert/rules/${id}/status`,
            method: 'post',
            data: { status }
        })
    },
    deleteMonitorHostAlertRule(id) {
        return request({
            url: `monitor/automation/host-alert/rules/${id}`,
            method: 'delete'
        })
    },
    scanMonitorHostAlerts() {
        return request({
            url: 'monitor/automation/host-alert/scan',
            method: 'post'
        })
    },
    getMonitorDBAlertRules() {
        return request({
            url: 'monitor/automation/db-alert/rules',
            method: 'get'
        })
    },
    createMonitorDBAlertRule(data) {
        return request({
            url: 'monitor/automation/db-alert/rules',
            method: 'post',
            data
        })
    },
    updateMonitorDBAlertRule(id, data) {
        return request({
            url: `monitor/automation/db-alert/rules/${id}`,
            method: 'put',
            data
        })
    },
    updateMonitorDBAlertRuleStatus(id, status) {
        return request({
            url: `monitor/automation/db-alert/rules/${id}/status`,
            method: 'post',
            data: { status }
        })
    },
    deleteMonitorDBAlertRule(id) {
        return request({
            url: `monitor/automation/db-alert/rules/${id}`,
            method: 'delete'
        })
    },
    getMonitorDBHealthSnapshots() {
        return request({
            url: 'monitor/automation/db-alert/snapshots',
            method: 'get'
        })
    },
    scanMonitorDBAlerts() {
        return request({
            url: 'monitor/automation/db-alert/scan',
            method: 'post'
        })
    },
    getMonitorSSLDomains() {
        return request({
            url: 'monitor/automation/ssl/domains',
            method: 'get'
        })
    },
    createMonitorSSLDomain(data) {
        return request({
            url: 'monitor/automation/ssl/domains',
            method: 'post',
            data
        })
    },
    updateMonitorSSLDomain(id, data) {
        return request({
            url: `monitor/automation/ssl/domains/${id}`,
            method: 'put',
            data
        })
    },
    updateMonitorSSLDomainStatus(id, status) {
        return request({
            url: `monitor/automation/ssl/domains/${id}/status`,
            method: 'post',
            data: { status }
        })
    },
    deleteMonitorSSLDomain(id) {
        return request({
            url: `monitor/automation/ssl/domains/${id}`,
            method: 'delete'
        })
    },
    getMonitorSSLSchedules() {
        return request({
            url: 'monitor/automation/ssl/schedules',
            method: 'get'
        })
    },
    updateMonitorSSLSchedule(id, data) {
        return request({
            url: `monitor/automation/ssl/schedules/${id}`,
            method: 'put',
            data
        })
    },
    scanMonitorSSLDomains() {
        return request({
            url: 'monitor/automation/ssl/scan',
            method: 'post'
        })
    },
    getMonitorSSLCerts() {
        return request({
            url: 'monitor/automation/ssl/certs',
            method: 'get'
        })
    },
    getMonitorSSLDeployLogs() {
        return request({
            url: 'monitor/automation/ssl/deploy/logs',
            method: 'get'
        })
    },
    deployMonitorSSLCert(data) {
        return request({
            url: 'monitor/automation/ssl/deploy',
            method: 'post',
            data
        })
    }
}
