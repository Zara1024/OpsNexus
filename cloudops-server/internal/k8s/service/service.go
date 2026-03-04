package service

import (
	"github.com/Zara1024/OpsNexus/cloudops-server/internal/k8s/repository"
	"github.com/Zara1024/OpsNexus/cloudops-server/pkg/crypto"
	"gorm.io/gorm"
)

// Services 聚合 k8s 模块业务服务。
type Services struct {
	Cluster   *ClusterService
	Node      *NodeService
	Namespace *NamespaceService
	Workload  *WorkloadService
	Network   *NetworkService
	Config    *ConfigService
}

// New 创建 k8s 模块服务集合。
func New(db *gorm.DB, aes *crypto.AESGCM) *Services {
	repos := repository.New(db)
	clusterSvc := NewClusterService(repos.Cluster, aes)
	return &Services{
		Cluster:   clusterSvc,
		Node:      NewNodeService(clusterSvc),
		Namespace: NewNamespaceService(clusterSvc),
		Workload:  NewWorkloadService(clusterSvc),
		Network:   NewNetworkService(clusterSvc),
		Config:    NewConfigService(clusterSvc),
	}
}
