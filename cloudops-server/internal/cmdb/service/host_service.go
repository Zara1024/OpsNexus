package service

import (
	"fmt"
	"net"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/Zara1024/OpsNexus/cloudops-server/internal/cmdb/model"
	"github.com/Zara1024/OpsNexus/cloudops-server/internal/cmdb/repository"
	apperrors "github.com/Zara1024/OpsNexus/cloudops-server/pkg/errors"
)

// HostService 提供主机资产管理能力。
type HostService struct {
	repo           *repository.HostRepository
	credentialRepo *repository.CredentialRepository
	credentialSvc  *CredentialService
}

// HostCreateRequest 定义创建主机请求。
type HostCreateRequest struct {
	Hostname      string  `json:"hostname" binding:"required,max=128"`
	InnerIP       string  `json:"inner_ip" binding:"required,max=64"`
	OuterIP       string  `json:"outer_ip" binding:"max=64"`
	OSType        string  `json:"os_type" binding:"max=64"`
	OSVersion     string  `json:"os_version" binding:"max=128"`
	Arch          string  `json:"arch" binding:"max=64"`
	CPUCores      int     `json:"cpu_cores"`
	MemoryGB      float64 `json:"memory_gb"`
	DiskGB        float64 `json:"disk_gb"`
	GroupID       int64   `json:"group_id"`
	CredentialID  int64   `json:"credential_id"`
	CloudProvider string  `json:"cloud_provider" binding:"max=64"`
	InstanceID    string  `json:"instance_id" binding:"max=128"`
	Region        string  `json:"region" binding:"max=64"`
	AgentStatus   string  `json:"agent_status" binding:"max=32"`
	Labels        string  `json:"labels"`
	Status        int16   `json:"status"`
}

// HostUpdateRequest 定义更新主机请求。
type HostUpdateRequest = HostCreateRequest

// HostBatchRequest 定义主机批量操作请求。
type HostBatchRequest struct {
	Action string  `json:"action" binding:"required,oneof=delete enable disable"`
	IDs    []int64 `json:"ids" binding:"required,min=1"`
}

// HostListResponse 定义主机分页结果。
type HostListResponse struct {
	List  []model.Host `json:"list"`
	Total int64        `json:"total"`
}

// NewHostService 创建主机服务。
func NewHostService(repo *repository.HostRepository, credentialRepo *repository.CredentialRepository, credentialSvc *CredentialService) *HostService {
	return &HostService{repo: repo, credentialRepo: credentialRepo, credentialSvc: credentialSvc}
}

// CreateHost 创建主机。
func (s *HostService) CreateHost(req *HostCreateRequest) error {
	status := req.Status
	if status == 0 {
		status = 1
	}
	agentStatus := strings.TrimSpace(req.AgentStatus)
	if agentStatus == "" {
		agentStatus = "unknown"
	}
	return s.repo.Create(&model.Host{
		Hostname:      strings.TrimSpace(req.Hostname),
		InnerIP:       strings.TrimSpace(req.InnerIP),
		OuterIP:       strings.TrimSpace(req.OuterIP),
		OSType:        strings.TrimSpace(req.OSType),
		OSVersion:     strings.TrimSpace(req.OSVersion),
		Arch:          strings.TrimSpace(req.Arch),
		CPUCores:      req.CPUCores,
		MemoryGB:      req.MemoryGB,
		DiskGB:        req.DiskGB,
		GroupID:       req.GroupID,
		CredentialID:  req.CredentialID,
		CloudProvider: strings.TrimSpace(req.CloudProvider),
		InstanceID:    strings.TrimSpace(req.InstanceID),
		Region:        strings.TrimSpace(req.Region),
		AgentStatus:   agentStatus,
		Labels:        strings.TrimSpace(req.Labels),
		Status:        status,
	})
}

// UpdateHost 更新主机。
func (s *HostService) UpdateHost(id int64, req *HostUpdateRequest) error {
	if _, err := s.repo.GetByID(id); err != nil {
		return apperrors.ErrNotFound
	}
	return s.repo.Update(id, map[string]interface{}{
		"hostname":       strings.TrimSpace(req.Hostname),
		"inner_ip":       strings.TrimSpace(req.InnerIP),
		"outer_ip":       strings.TrimSpace(req.OuterIP),
		"os_type":        strings.TrimSpace(req.OSType),
		"os_version":     strings.TrimSpace(req.OSVersion),
		"arch":           strings.TrimSpace(req.Arch),
		"cpu_cores":      req.CPUCores,
		"memory_gb":      req.MemoryGB,
		"disk_gb":        req.DiskGB,
		"group_id":       req.GroupID,
		"credential_id":  req.CredentialID,
		"cloud_provider": strings.TrimSpace(req.CloudProvider),
		"instance_id":    strings.TrimSpace(req.InstanceID),
		"region":         strings.TrimSpace(req.Region),
		"agent_status":   strings.TrimSpace(req.AgentStatus),
		"labels":         strings.TrimSpace(req.Labels),
		"status":         req.Status,
	})
}

// DeleteHost 删除主机。
func (s *HostService) DeleteHost(id int64) error {
	return s.repo.SoftDelete(id)
}

// GetHost 查询主机详情。
func (s *HostService) GetHost(id int64) (*model.Host, error) {
	item, err := s.repo.GetByID(id)
	if err != nil {
		return nil, apperrors.ErrNotFound
	}
	return item, nil
}

// ListHosts 分页查询主机。
func (s *HostService) ListHosts(filter *repository.HostListFilter) (*HostListResponse, error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 10
	}
	list, total, err := s.repo.List(filter)
	if err != nil {
		return nil, err
	}
	return &HostListResponse{List: list, Total: total}, nil
}

// BatchOperateHosts 批量处理主机。
func (s *HostService) BatchOperateHosts(req *HostBatchRequest) error {
	switch req.Action {
	case "delete":
		return s.repo.BatchDelete(req.IDs)
	case "enable":
		return s.repo.BatchUpdateStatus(req.IDs, 1)
	case "disable":
		return s.repo.BatchUpdateStatus(req.IDs, 2)
	default:
		return &apperrors.AppError{Code: 20001, HTTPCode: 400, Message: "不支持的批量操作"}
	}
}

// TestSSHConnection 测试主机 SSH 连通性。
func (s *HostService) TestSSHConnection(hostID int64) error {
	host, err := s.repo.GetByID(hostID)
	if err != nil {
		return apperrors.ErrNotFound
	}
	credential, err := s.credentialRepo.GetByID(host.CredentialID)
	if err != nil {
		return &apperrors.AppError{Code: 20002, HTTPCode: 400, Message: "主机未配置有效凭据"}
	}
	username, password, privateKey, passphrase, err := s.credentialSvc.ResolvePlainCredential(credential.ID)
	if err != nil {
		return err
	}

	authMethods := make([]ssh.AuthMethod, 0)
	if credential.Type == "key" && strings.TrimSpace(privateKey) != "" {
		var signer ssh.Signer
		if strings.TrimSpace(passphrase) != "" {
			signer, err = ssh.ParsePrivateKeyWithPassphrase([]byte(privateKey), []byte(passphrase))
		} else {
			signer, err = ssh.ParsePrivateKey([]byte(privateKey))
		}
		if err != nil {
			return &apperrors.AppError{Code: 20003, HTTPCode: 400, Message: "私钥解析失败"}
		}
		authMethods = append(authMethods, ssh.PublicKeys(signer))
	} else if strings.TrimSpace(password) != "" {
		authMethods = append(authMethods, ssh.Password(password))
	}
	if len(authMethods) == 0 {
		return &apperrors.AppError{Code: 20004, HTTPCode: 400, Message: "凭据内容为空"}
	}

	client, err := ssh.Dial("tcp", net.JoinHostPort(host.InnerIP, "22"), &ssh.ClientConfig{
		User:            username,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         8 * time.Second,
	})
	if err != nil {
		return &apperrors.AppError{Code: 20005, HTTPCode: 400, Message: fmt.Sprintf("SSH连接失败: %v", err)}
	}
	_ = client.Close()
	return nil
}
