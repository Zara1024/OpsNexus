package service

import (
	"context"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	apperrors "github.com/Zara1024/OpsNexus/cloudops-server/pkg/errors"
)

// NodeService 提供节点管理能力。
type NodeService struct {
	clusterSvc *ClusterService
}

// NodeItem 定义节点展示结构。
type NodeItem struct {
	Name          string            `json:"name"`
	Roles         []string          `json:"roles"`
	Status        string            `json:"status"`
	PodCount      int               `json:"pod_count"`
	Labels        map[string]string `json:"labels"`
	Unschedulable bool              `json:"unschedulable"`
}

// NewNodeService 创建节点服务。
func NewNodeService(clusterSvc *ClusterService) *NodeService {
	return &NodeService{clusterSvc: clusterSvc}
}

// ListNodes 查询节点列表。
func (s *NodeService) ListNodes(clusterID int64) ([]NodeItem, error) {
	client, _, err := s.clusterSvc.GetClient(clusterID)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	nodes, err := client.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	pods, _ := client.CoreV1().Pods("").List(ctx, metav1.ListOptions{})
	podCountMap := map[string]int{}
	for _, p := range pods.Items {
		if p.Spec.NodeName != "" {
			podCountMap[p.Spec.NodeName]++
		}
	}
	res := make([]NodeItem, 0, len(nodes.Items))
	for _, n := range nodes.Items {
		item := NodeItem{
			Name:          n.Name,
			Roles:         parseNodeRoles(n.Labels),
			Status:        parseNodeStatus(n.Status.Conditions),
			PodCount:      podCountMap[n.Name],
			Labels:        n.Labels,
			Unschedulable: n.Spec.Unschedulable,
		}
		res = append(res, item)
	}
	return res, nil
}

// CordonNode 封锁节点。
func (s *NodeService) CordonNode(clusterID int64, nodeName string) error {
	return s.patchUnschedulable(clusterID, nodeName, true)
}

// UncordonNode 解封节点。
func (s *NodeService) UncordonNode(clusterID int64, nodeName string) error {
	return s.patchUnschedulable(clusterID, nodeName, false)
}

// DrainNode 简化实现排水逻辑。
func (s *NodeService) DrainNode(clusterID int64, nodeName string) error {
	if err := s.CordonNode(clusterID, nodeName); err != nil {
		return err
	}
	client, _, err := s.clusterSvc.GetClient(clusterID)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	pods, err := client.CoreV1().Pods("").List(ctx, metav1.ListOptions{FieldSelector: "spec.nodeName=" + nodeName})
	if err != nil {
		return err
	}
	for _, p := range pods.Items {
		if p.Namespace == "kube-system" {
			continue
		}
		if err = client.CoreV1().Pods(p.Namespace).Delete(ctx, p.Name, metav1.DeleteOptions{}); err != nil {
			return &apperrors.AppError{Code: 30007, HTTPCode: 400, Message: "节点排水失败"}
		}
	}
	return nil
}

func (s *NodeService) patchUnschedulable(clusterID int64, nodeName string, value bool) error {
	client, _, err := s.clusterSvc.GetClient(clusterID)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()
	patch := []byte("{\"spec\":{\"unschedulable\":" + map[bool]string{true: "true", false: "false"}[value] + "}}")
	_, err = client.CoreV1().Nodes().Patch(ctx, nodeName, types.StrategicMergePatchType, patch, metav1.PatchOptions{})
	return err
}

func parseNodeRoles(labels map[string]string) []string {
	roles := make([]string, 0)
	for k := range labels {
		if strings.HasPrefix(k, "node-role.kubernetes.io/") {
			roles = append(roles, strings.TrimPrefix(k, "node-role.kubernetes.io/"))
		}
	}
	if len(roles) == 0 {
		roles = append(roles, "worker")
	}
	return roles
}

func parseNodeStatus(conditions []corev1.NodeCondition) string {
	for _, c := range conditions {
		if c.Type == corev1.NodeReady {
			if c.Status == corev1.ConditionTrue {
				return "Ready"
			}
			return "NotReady"
		}
	}
	return "Unknown"
}
