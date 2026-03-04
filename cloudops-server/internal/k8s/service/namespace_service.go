package service

import (
	"context"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NamespaceService 提供命名空间能力。
type NamespaceService struct {
	clusterSvc *ClusterService
}

// NewNamespaceService 创建命名空间服务。
func NewNamespaceService(clusterSvc *ClusterService) *NamespaceService {
	return &NamespaceService{clusterSvc: clusterSvc}
}

// ListNamespaces 查询命名空间列表。
func (s *NamespaceService) ListNamespaces(clusterID int64) ([]string, error) {
	client, _, err := s.clusterSvc.GetClient(clusterID)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()
	res, err := client.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	list := make([]string, 0, len(res.Items))
	for _, item := range res.Items {
		list = append(list, item.Name)
	}
	return list, nil
}
