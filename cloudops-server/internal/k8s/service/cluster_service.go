package service

import (
	"context"
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/Zara1024/OpsNexus/cloudops-server/internal/k8s/model"
	"github.com/Zara1024/OpsNexus/cloudops-server/internal/k8s/repository"
	"github.com/Zara1024/OpsNexus/cloudops-server/pkg/crypto"
	apperrors "github.com/Zara1024/OpsNexus/cloudops-server/pkg/errors"
)

// ClusterService 提供集群管理能力。
type ClusterService struct {
	repo *repository.ClusterRepository
	aes  *crypto.AESGCM
}

// ClusterCreateRequest 定义新增集群请求。
type ClusterCreateRequest struct {
	Name         string `json:"name" binding:"required,max=128"`
	APIServerURL string `json:"api_server_url" binding:"required,max=512"`
	Kubeconfig   string `json:"kubeconfig" binding:"required"`
	Description  string `json:"description" binding:"max=512"`
}

// ClusterUpdateRequest 定义更新集群请求。
type ClusterUpdateRequest struct {
	Name         string `json:"name" binding:"required,max=128"`
	APIServerURL string `json:"api_server_url" binding:"required,max=512"`
	Kubeconfig   string `json:"kubeconfig"`
	Description  string `json:"description" binding:"max=512"`
	Status       int16  `json:"status"`
}

// ClusterDetail 定义集群详情响应。
type ClusterDetail struct {
	model.Cluster
	Health string `json:"health"`
}

// NewClusterService 创建集群服务。
func NewClusterService(repo *repository.ClusterRepository, aes *crypto.AESGCM) *ClusterService {
	return &ClusterService{repo: repo, aes: aes}
}

// ListClusters 查询集群列表。
func (s *ClusterService) ListClusters() ([]model.Cluster, error) {
	return s.repo.List()
}

// CreateCluster 注册集群。
func (s *ClusterService) CreateCluster(req *ClusterCreateRequest, createdBy int64) error {
	cipher, err := s.aes.Encrypt(strings.TrimSpace(req.Kubeconfig))
	if err != nil {
		return err
	}
	item := &model.Cluster{
		Name:                strings.TrimSpace(req.Name),
		APIServerURL:        strings.TrimSpace(req.APIServerURL),
		KubeconfigEncrypted: cipher,
		Description:         strings.TrimSpace(req.Description),
		Status:              1,
		CreatedBy:           createdBy,
	}
	return s.repo.Create(item)
}

// UpdateCluster 更新集群。
func (s *ClusterService) UpdateCluster(id int64, req *ClusterUpdateRequest) error {
	if _, err := s.repo.GetByID(id); err != nil {
		return apperrors.ErrNotFound
	}
	updates := map[string]interface{}{
		"name":           strings.TrimSpace(req.Name),
		"api_server_url": strings.TrimSpace(req.APIServerURL),
		"description":    strings.TrimSpace(req.Description),
		"status":         req.Status,
	}
	if updates["status"].(int16) == 0 {
		updates["status"] = int16(1)
	}
	if strings.TrimSpace(req.Kubeconfig) != "" {
		cipher, err := s.aes.Encrypt(strings.TrimSpace(req.Kubeconfig))
		if err != nil {
			return err
		}
		updates["kubeconfig_encrypted"] = cipher
	}
	return s.repo.Update(id, updates)
}

// DeleteCluster 删除集群。
func (s *ClusterService) DeleteCluster(id int64) error {
	return s.repo.Delete(id)
}

// GetClusterDetail 获取集群详情和健康状态。
func (s *ClusterService) GetClusterDetail(id int64) (*ClusterDetail, error) {
	item, err := s.repo.GetByID(id)
	if err != nil {
		return nil, apperrors.ErrNotFound
	}
	health := "unknown"
	if err = s.TestCluster(id); err == nil {
		health = "healthy"
	} else {
		health = "unreachable"
	}
	return &ClusterDetail{Cluster: *item, Health: health}, nil
}

// TestCluster 测试集群连接。
func (s *ClusterService) TestCluster(id int64) error {
	client, _, err := s.GetClient(id)
	if err != nil {
		return &apperrors.AppError{Code: 30001, HTTPCode: 400, Message: "集群连接失败"}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()
	_, err = client.Discovery().ServerVersion()
	if err != nil {
		return &apperrors.AppError{Code: 30001, HTTPCode: 400, Message: "集群连接失败"}
	}
	_, err = client.CoreV1().Nodes().List(ctx, metav1.ListOptions{Limit: 1})
	if err != nil {
		return &apperrors.AppError{Code: 30002, HTTPCode: 400, Message: "集群鉴权失败"}
	}
	return nil
}

// GetClient 构建 Kubernetes 客户端。
func (s *ClusterService) GetClient(id int64) (*kubernetes.Clientset, *rest.Config, error) {
	item, err := s.repo.GetByID(id)
	if err != nil {
		return nil, nil, apperrors.ErrNotFound
	}
	plain, err := s.aes.Decrypt(item.KubeconfigEncrypted)
	if err != nil {
		return nil, nil, err
	}
	restConfig, err := clientcmd.RESTConfigFromKubeConfig([]byte(plain))
	if err != nil {
		return nil, nil, err
	}
	restConfig.Timeout = 12 * time.Second
	client, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, nil, err
	}
	return client, restConfig, nil
}
