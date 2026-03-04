package service

import (
	"context"
	"time"

	v1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NetworkService 提供网络资源能力。
type NetworkService struct {
	clusterSvc *ClusterService
}

// NewNetworkService 创建网络服务。
func NewNetworkService(clusterSvc *ClusterService) *NetworkService {
	return &NetworkService{clusterSvc: clusterSvc}
}

// ListServices 查询 Service 列表。
func (s *NetworkService) ListServices(clusterID int64, namespace string) ([]v1.Service, error) {
	client, _, err := s.clusterSvc.GetClient(clusterID)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, err := client.CoreV1().Services(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return res.Items, nil
}

// ListIngresses 查询 Ingress 列表。
func (s *NetworkService) ListIngresses(clusterID int64, namespace string) ([]networkingv1.Ingress, error) {
	client, _, err := s.clusterSvc.GetClient(clusterID)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, err := client.NetworkingV1().Ingresses(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return res.Items, nil
}
