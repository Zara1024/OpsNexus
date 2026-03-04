package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// WorkloadService 提供工作负载能力。
type WorkloadService struct {
	clusterSvc *ClusterService
}

// DeploymentScaleRequest 定义伸缩请求。
type DeploymentScaleRequest struct {
	Replicas int32 `json:"replicas" binding:"required,min=0"`
}

// NewWorkloadService 创建工作负载服务。
func NewWorkloadService(clusterSvc *ClusterService) *WorkloadService {
	return &WorkloadService{clusterSvc: clusterSvc}
}

// ListDeployments 查询 Deployment 列表。
func (s *WorkloadService) ListDeployments(clusterID int64, namespace string) ([]appsv1.Deployment, error) {
	client, _, err := s.clusterSvc.GetClient(clusterID)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, err := client.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return res.Items, nil
}

// GetDeployment 查询 Deployment 详情。
func (s *WorkloadService) GetDeployment(clusterID int64, namespace, name string) (*appsv1.Deployment, error) {
	client, _, err := s.clusterSvc.GetClient(clusterID)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return client.AppsV1().Deployments(namespace).Get(ctx, name, metav1.GetOptions{})
}

// CreateDeployment 创建 Deployment。
func (s *WorkloadService) CreateDeployment(clusterID int64, namespace string, item *appsv1.Deployment) error {
	client, _, err := s.clusterSvc.GetClient(clusterID)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err = client.AppsV1().Deployments(namespace).Create(ctx, item, metav1.CreateOptions{})
	return err
}

// UpdateDeployment 更新 Deployment。
func (s *WorkloadService) UpdateDeployment(clusterID int64, namespace, name string, item *appsv1.Deployment) error {
	client, _, err := s.clusterSvc.GetClient(clusterID)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	old, err := client.AppsV1().Deployments(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	item.ResourceVersion = old.ResourceVersion
	_, err = client.AppsV1().Deployments(namespace).Update(ctx, item, metav1.UpdateOptions{})
	return err
}

// DeleteDeployment 删除 Deployment。
func (s *WorkloadService) DeleteDeployment(clusterID int64, namespace, name string) error {
	client, _, err := s.clusterSvc.GetClient(clusterID)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return client.AppsV1().Deployments(namespace).Delete(ctx, name, metav1.DeleteOptions{})
}

// ScaleDeployment 执行 Deployment 伸缩。
func (s *WorkloadService) ScaleDeployment(clusterID int64, namespace, name string, replicas int32) error {
	client, _, err := s.clusterSvc.GetClient(clusterID)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	scale, err := client.AppsV1().Deployments(namespace).GetScale(ctx, name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	scale.Spec.Replicas = replicas
	_, err = client.AppsV1().Deployments(namespace).UpdateScale(ctx, name, scale, metav1.UpdateOptions{})
	return err
}

// RestartDeployment 执行滚动重启。
func (s *WorkloadService) RestartDeployment(clusterID int64, namespace, name string) error {
	client, _, err := s.clusterSvc.GetClient(clusterID)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	item, err := client.AppsV1().Deployments(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	if item.Spec.Template.Annotations == nil {
		item.Spec.Template.Annotations = map[string]string{}
	}
	item.Spec.Template.Annotations["kubectl.kubernetes.io/restartedAt"] = time.Now().Format(time.RFC3339)
	_, err = client.AppsV1().Deployments(namespace).Update(ctx, item, metav1.UpdateOptions{})
	return err
}

// ListPods 查询 Pod 列表。
func (s *WorkloadService) ListPods(clusterID int64, namespace string) ([]corev1.Pod, error) {
	client, _, err := s.clusterSvc.GetClient(clusterID)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, err := client.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return res.Items, nil
}

// GetPod 查询 Pod 详情。
func (s *WorkloadService) GetPod(clusterID int64, namespace, pod string) (*corev1.Pod, error) {
	client, _, err := s.clusterSvc.GetClient(clusterID)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return client.CoreV1().Pods(namespace).Get(ctx, pod, metav1.GetOptions{})
}

// DeletePod 删除 Pod。
func (s *WorkloadService) DeletePod(clusterID int64, namespace, pod string) error {
	client, _, err := s.clusterSvc.GetClient(clusterID)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return client.CoreV1().Pods(namespace).Delete(ctx, pod, metav1.DeleteOptions{})
}

// GetPodLogs 获取 Pod 日志。
func (s *WorkloadService) GetPodLogs(clusterID int64, namespace, pod, container string, tailLines int64) (string, error) {
	client, _, err := s.clusterSvc.GetClient(clusterID)
	if err != nil {
		return "", err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()
	opts := &corev1.PodLogOptions{Container: container, TailLines: &tailLines, Timestamps: true}
	req := client.CoreV1().Pods(namespace).GetLogs(pod, opts)
	stream, err := req.Stream(ctx)
	if err != nil {
		return "", err
	}
	defer stream.Close()
	buf := bytes.NewBuffer(nil)
	_, err = io.Copy(buf, stream)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

// BuildDemoDeployment 创建演示 Deployment 结构。
func BuildDemoDeployment(name, namespace, image string, replicas int32) *appsv1.Deployment {
	labels := map[string]string{"app": name}
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{MatchLabels: labels},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: labels},
				Spec: corev1.PodSpec{Containers: []corev1.Container{{
					Name:           name,
					Image:          image,
					Ports:          []corev1.ContainerPort{{ContainerPort: 80}},
					ReadinessProbe: &corev1.Probe{ProbeHandler: corev1.ProbeHandler{HTTPGet: &corev1.HTTPGetAction{Path: "/", Port: intstr.FromInt(80)}}, InitialDelaySeconds: 5},
				}}},
			},
		},
	}
}

// BuildPodLogSSEEvent 组装 SSE 数据片段。
func BuildPodLogSSEEvent(line string) string {
	return fmt.Sprintf("data: %s\n\n", line)
}
