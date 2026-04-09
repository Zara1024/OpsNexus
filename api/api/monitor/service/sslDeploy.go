package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"path"
	"strconv"
	"strings"
	"time"

	cmdbDao "dodevops-api/api/cmdb/dao"
	cmdbModel "dodevops-api/api/cmdb/model"
	configModel "dodevops-api/api/configcenter/model"
	monitorDao "dodevops-api/api/monitor/dao"
	"dodevops-api/api/monitor/model"
	"dodevops-api/common"
	"dodevops-api/common/util"

	"github.com/pkg/sftp"
	"github.com/robfig/cron/v3"
	"golang.org/x/crypto/ssh"
)

func (s *MonitorAutomationServiceImpl) deploySSLCertificate(ctx context.Context, req model.MonitorSSLDeployRequest) (*model.MonitorSSLCertDeployLogEntity, error) {
	cert, err := monitorDao.GetMonitorSSLCertByID(req.CertID)
	if err != nil {
		return nil, fmt.Errorf("SSL 证书不存在: %w", err)
	}
	if strings.TrimSpace(cert.Cert) == "" || strings.TrimSpace(cert.PrivateKey) == "" {
		return nil, fmt.Errorf("SSL 证书或私钥内容为空")
	}

	host, err := cmdbDao.NewCmdbHostSSHDao().GetHostSSHInfo(req.HostID)
	if err != nil {
		return nil, fmt.Errorf("目标主机不存在: %w", err)
	}

	deployPath := strings.TrimSpace(req.DeployPath)
	if deployPath == "" {
		return nil, fmt.Errorf("部署路径不能为空")
	}
	reloadCommand := strings.TrimSpace(req.ReloadCommand)
	if reloadCommand == "" {
		reloadCommand = "nginx -s reload || systemctl reload nginx || systemctl reload openresty"
	}

	logLines := make([]string, 0, 16)
	backupFiles := make([]string, 0, 2)
	deployFiles := make([]string, 0, 2)
	appendSSLDeployLog(&logLines, "开始部署证书到主机: %s (%s)", hostDisplayName(*host), firstNonEmptyAutomation(host.SSHIP, host.PrivateIP, host.PublicIP))

	logEntry := &model.MonitorSSLCertDeployLogEntity{
		CertID:     cert.ID,
		Domain:     cert.Domain,
		HostID:     host.ID,
		HostName:   hostDisplayName(*host),
		DeployPath: deployPath,
		Status:     1,
		BackupFiles: "[]",
		DeployFiles: "[]",
		Logs:       "[]",
		ErrorMsg:   "",
	}
	if err = monitorDao.CreateMonitorSSLCertDeployLog(logEntry); err != nil {
		return nil, err
	}

	updateLogEntry := func(status int, errorMsg string) error {
		logEntry.Status = status
		logEntry.ErrorMsg = strings.TrimSpace(errorMsg)
		logEntry.BackupFiles = encodeStringList(backupFiles)
		logEntry.DeployFiles = encodeStringList(deployFiles)
		logEntry.Logs = encodeStringList(logLines)
		return monitorDao.UpdateMonitorSSLCertDeployLog(logEntry)
	}

	client, sftpClient, err := openHostSFTPClient(host)
	if err != nil {
		appendSSLDeployLog(&logLines, "SSH/SFTP 连接失败: %v", err)
		_ = updateLogEntry(3, err.Error())
		_, _ = s.triggerAutomationEvent(automationTriggerInput{
			RuleCategory: model.MonitorRuleCategorySSL,
			ResourceType: "ssl_cert",
			ResourceID:   cert.ID,
			ResourceName: cert.Domain,
			EventKey:     "deploy_failed",
			Title:        fmt.Sprintf("SSL 部署失败: %s", cert.Domain),
			Summary:      err.Error(),
			Detail:       strings.Join(logLines, "\n"),
			Severity:     "P2",
			Operator:     ">=",
			CurrentValue: 1,
		})
		return nil, err
	}
	defer client.Close()
	defer sftpClient.Close()

	pemFile := path.Join(deployPath, cert.Domain+".pem")
	keyFile := path.Join(deployPath, cert.Domain+".key")
	tmpPemFile := pemFile + ".tmp"
	tmpKeyFile := keyFile + ".tmp"

	if _, err = runRemoteSSHCommand(ctx, client, "mkdir -p "+shellQuote(deployPath)); err != nil {
		appendSSLDeployLog(&logLines, "创建部署目录失败: %v", err)
		_ = updateLogEntry(3, err.Error())
		return nil, err
	}
	appendSSLDeployLog(&logLines, "部署目录已就绪: %s", deployPath)

	timestamp := time.Now().Format("20060102150405")
	if ok, _ := remoteFileExists(ctx, client, pemFile); ok {
		backup := pemFile + "." + timestamp + ".bak"
		if _, err = runRemoteSSHCommand(ctx, client, "cp "+shellQuote(pemFile)+" "+shellQuote(backup)); err == nil {
			backupFiles = append(backupFiles, backup)
			appendSSLDeployLog(&logLines, "已备份证书文件: %s", backup)
		}
	}
	if ok, _ := remoteFileExists(ctx, client, keyFile); ok {
		backup := keyFile + "." + timestamp + ".bak"
		if _, err = runRemoteSSHCommand(ctx, client, "cp "+shellQuote(keyFile)+" "+shellQuote(backup)); err == nil {
			backupFiles = append(backupFiles, backup)
			appendSSLDeployLog(&logLines, "已备份密钥文件: %s", backup)
		}
	}

	if err = uploadRemoteFile(sftpClient, tmpPemFile, cert.Cert); err != nil {
		appendSSLDeployLog(&logLines, "上传证书文件失败: %v", err)
		_ = updateLogEntry(3, err.Error())
		return nil, err
	}
	if err = uploadRemoteFile(sftpClient, tmpKeyFile, cert.PrivateKey); err != nil {
		appendSSLDeployLog(&logLines, "上传私钥文件失败: %v", err)
		_ = updateLogEntry(3, err.Error())
		return nil, err
	}

	if _, err = runRemoteSSHCommand(ctx, client, "mv "+shellQuote(tmpPemFile)+" "+shellQuote(pemFile)); err != nil {
		appendSSLDeployLog(&logLines, "移动证书文件失败: %v", err)
		_ = updateLogEntry(3, err.Error())
		return nil, err
	}
	if _, err = runRemoteSSHCommand(ctx, client, "mv "+shellQuote(tmpKeyFile)+" "+shellQuote(keyFile)); err != nil {
		appendSSLDeployLog(&logLines, "移动私钥文件失败: %v", err)
		_ = updateLogEntry(3, err.Error())
		return nil, err
	}
	deployFiles = append(deployFiles, pemFile, keyFile)
	appendSSLDeployLog(&logLines, "证书文件已部署: %s", pemFile)
	appendSSLDeployLog(&logLines, "私钥文件已部署: %s", keyFile)

	if _, err = runRemoteSSHCommand(ctx, client, "chmod 600 "+shellQuote(pemFile)+" "+shellQuote(keyFile)); err != nil {
		appendSSLDeployLog(&logLines, "设置文件权限失败: %v", err)
		_ = updateLogEntry(3, err.Error())
		return nil, err
	}
	appendSSLDeployLog(&logLines, "已设置文件权限为 600")

	if output, err := runRemoteSSHCommand(ctx, client, reloadCommand); err != nil {
		appendSSLDeployLog(&logLines, "重载服务失败: %v", err)
		if strings.TrimSpace(output) != "" {
			appendSSLDeployLog(&logLines, "重载输出: %s", output)
		}
		_ = updateLogEntry(3, err.Error())
		_, _ = s.triggerAutomationEvent(automationTriggerInput{
			RuleCategory: model.MonitorRuleCategorySSL,
			ResourceType: "ssl_cert",
			ResourceID:   cert.ID,
			ResourceName: cert.Domain,
			EventKey:     "deploy_failed",
			Title:        fmt.Sprintf("SSL 部署失败: %s", cert.Domain),
			Summary:      err.Error(),
			Detail:       strings.Join(logLines, "\n"),
			Severity:     "P2",
			Operator:     ">=",
			CurrentValue: 1,
		})
		return nil, err
	} else {
		appendSSLDeployLog(&logLines, "服务重载成功")
		if strings.TrimSpace(output) != "" {
			appendSSLDeployLog(&logLines, "重载输出: %s", output)
		}
	}

	appendSSLDeployLog(&logLines, "证书部署完成")
	if err = updateLogEntry(2, ""); err != nil {
		return nil, err
	}
	_ = s.resolveAutomationEvent(automationResolveInput{
		Fingerprint:    buildAutomationFingerprint("ssl_cert", cert.ID, "deploy_failed"),
		RecoveryValue:  0,
		Solution:       fmt.Sprintf("%s 证书部署恢复正常", cert.Domain),
		RecoveryRemark: strings.Join(logLines, "\n"),
	})
	return logEntry, nil
}

func openHostSFTPClient(host *cmdbModel.CmdbHost) (*ssh.Client, *sftp.Client, error) {
	authMethod, err := buildHostSSHAuthMethod(host)
	if err != nil {
		return nil, nil, err
	}
	client, err := ssh.Dial("tcp", net.JoinHostPort(host.SSHIP, strconv.Itoa(host.SSHPort)), &ssh.ClientConfig{
		User:            host.SSHName,
		Auth:            []ssh.AuthMethod{authMethod},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         20 * time.Second,
	})
	if err != nil {
		return nil, nil, err
	}
	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		client.Close()
		return nil, nil, err
	}
	return client, sftpClient, nil
}

func buildHostSSHAuthMethod(host *cmdbModel.CmdbHost) (ssh.AuthMethod, error) {
	if host.SSHKeyID == 0 {
		return nil, fmt.Errorf("主机未配置 SSH 凭据")
	}
	var auth configModel.EcsAuth
	if err := common.GetDB().Table("config_ecsauth").Where("id = ?", host.SSHKeyID).First(&auth).Error; err != nil {
		return nil, err
	}
	switch auth.Type {
	case 1:
		return ssh.Password(auth.Password), nil
	case 2:
		return util.NewSSHUtil().PublicKeyAuth(auth.PublicKey)
	case 3:
		if userKeyAuth, err := util.NewSSHUtil().UserKeyAuth(); err == nil {
			return userKeyAuth, nil
		}
		return util.NewSSHUtil().DefaultKeyAuth()
	default:
		return nil, fmt.Errorf("不支持的 SSH 凭据类型: %d", auth.Type)
	}
}

func uploadRemoteFile(client *sftp.Client, remotePath, content string) error {
	file, err := client.Create(remotePath)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err = file.Write([]byte(content)); err != nil {
		return err
	}
	return nil
}

func remoteFileExists(ctx context.Context, client *ssh.Client, remotePath string) (bool, error) {
	output, err := runRemoteSSHCommand(ctx, client, "[ -f "+shellQuote(remotePath)+" ] && echo yes || echo no")
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(output) == "yes", nil
}

func runRemoteSSHCommand(ctx context.Context, client *ssh.Client, command string) (string, error) {
	session, err := client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	type result struct {
		output []byte
		err    error
	}
	ch := make(chan result, 1)
	go func() {
		output, runErr := session.CombinedOutput(command)
		ch <- result{output: output, err: runErr}
	}()
	select {
	case <-ctx.Done():
		_ = session.Close()
		return "", ctx.Err()
	case res := <-ch:
		return strings.TrimSpace(string(res.output)), res.err
	}
}

func shellQuote(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "'\"'\"'") + "'"
}

func appendSSLDeployLog(lines *[]string, format string, args ...interface{}) {
	*lines = append(*lines, fmt.Sprintf("[%s] %s", time.Now().Format("2006-01-02 15:04:05"), fmt.Sprintf(format, args...)))
}

func encodeStringList(values []string) string {
	if len(values) == 0 {
		return "[]"
	}
	data, err := json.Marshal(values)
	if err != nil {
		return "[]"
	}
	return string(data)
}

func calculateNextCronRun(expr string, now time.Time) (time.Time, error) {
	schedule, err := cron.ParseStandard(expr)
	if err != nil {
		return time.Time{}, err
	}
	return schedule.Next(now), nil
}
