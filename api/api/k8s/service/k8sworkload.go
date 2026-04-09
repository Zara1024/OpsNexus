package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	aiModel "dodevops-api/api/ai/model"
	"dodevops-api/api/k8s/dao"
	"dodevops-api/api/k8s/model"
	"dodevops-api/common/result"
	pkgdb "dodevops-api/pkg/db"
	"dodevops-api/pkg/jwt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	appsv1 "k8s.io/api/apps/v1"
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	metricsv1beta1 "k8s.io/metrics/pkg/apis/metrics/v1beta1"
	metricsclientset "k8s.io/metrics/pkg/client/clientset/versioned"
	"sigs.k8s.io/yaml"
)

// IK8sWorkloadService K8s工作负载服务接口
type IK8sWorkloadService interface {
	// 工作负载列表和详情
	GetWorkloads(c *gin.Context, clusterId uint, namespaceName string, workloadType string, name string)
	GetWorkloadDetail(c *gin.Context, clusterId uint, namespaceName string, workloadType string, workloadName string)

	// Deployment管理
	CreateDeployment(c *gin.Context, clusterId uint, namespaceName string, req *model.CreateDeploymentRequest)
	UpdateDeployment(c *gin.Context, clusterId uint, namespaceName string, deploymentName string, req *model.UpdateWorkloadRequest)
	DeleteDeployment(c *gin.Context, clusterId uint, namespaceName string, deploymentName string)
	ScaleDeployment(c *gin.Context, clusterId uint, namespaceName string, deploymentName string, req *model.ScaleWorkloadRequest)
	RestartDeployment(c *gin.Context, clusterId uint, namespaceName string, deploymentName string)
	GetWorkloadAutoscaling(c *gin.Context, clusterId uint, namespaceName string, workloadType string, workloadName string)
	UpsertWorkloadAutoscaling(c *gin.Context, clusterId uint, namespaceName string, workloadType string, workloadName string, req *model.UpsertWorkloadAutoscalingRequest)
	DeleteWorkloadAutoscaling(c *gin.Context, clusterId uint, namespaceName string, workloadType string, workloadName string)
	GetWorkloadCapacitySuggestion(c *gin.Context, clusterId uint, namespaceName string, workloadType string, workloadName string)
	GetWorkloadCapacitySuggestionHistory(c *gin.Context, clusterId uint, namespaceName string, workloadType string, workloadName string, query model.WorkloadCapacitySuggestionHistoryQuery)
	GetWorkloadGovernanceOverview(c *gin.Context, clusterId uint, namespaceName string, workloadType string, workloadName string)

	// StatefulSet管理
	CreateStatefulSet(c *gin.Context, clusterId uint, namespaceName string, req *model.CreateStatefulSetRequest)
	UpdateStatefulSet(c *gin.Context, clusterId uint, namespaceName string, statefulSetName string, req *model.UpdateWorkloadRequest)
	DeleteStatefulSet(c *gin.Context, clusterId uint, namespaceName string, statefulSetName string)
	ScaleStatefulSet(c *gin.Context, clusterId uint, namespaceName string, statefulSetName string, req *model.ScaleWorkloadRequest)
	RestartStatefulSet(c *gin.Context, clusterId uint, namespaceName string, statefulSetName string)

	// DaemonSet管理
	CreateDaemonSet(c *gin.Context, clusterId uint, namespaceName string, req *model.CreateDaemonSetRequest)
	UpdateDaemonSet(c *gin.Context, clusterId uint, namespaceName string, daemonSetName string, req *model.UpdateWorkloadRequest)
	DeleteDaemonSet(c *gin.Context, clusterId uint, namespaceName string, daemonSetName string)
	RestartDaemonSet(c *gin.Context, clusterId uint, namespaceName string, daemonSetName string)

	// Job管理
	CreateJob(c *gin.Context, clusterId uint, namespaceName string, req *model.CreateJobRequest)
	UpdateJob(c *gin.Context, clusterId uint, namespaceName string, jobName string, req *model.UpdateWorkloadRequest)
	DeleteJob(c *gin.Context, clusterId uint, namespaceName string, jobName string)

	// CronJob管理
	CreateCronJob(c *gin.Context, clusterId uint, namespaceName string, req *model.CreateCronJobRequest)
	UpdateCronJob(c *gin.Context, clusterId uint, namespaceName string, cronJobName string, req *model.UpdateWorkloadRequest)
	DeleteCronJob(c *gin.Context, clusterId uint, namespaceName string, cronJobName string)

	// Pod管理
	GetPods(c *gin.Context, clusterId uint, namespaceName string)
	GetPodDetail(c *gin.Context, clusterId uint, namespaceName string, podName string)
	DeletePod(c *gin.Context, clusterId uint, namespaceName string, podName string)
	GetPodLogs(c *gin.Context, clusterId uint, namespaceName string, podName string, containerName string)
	GetPodEvents(c *gin.Context, clusterId uint, namespaceName string, podName string)
	GetPodYaml(c *gin.Context, clusterId uint, namespaceName string, podName string)
	UpdatePodYaml(c *gin.Context, clusterId uint, namespaceName string, podName string, req *model.UpdatePodYAMLRequest)
	CreatePodFromYAML(c *gin.Context, clusterId uint, namespaceName string, req *model.CreatePodFromYAMLRequest)
	ValidateYAML(c *gin.Context, clusterId uint, req *model.ValidateYAMLRequest)

	// 通用工作负载YAML管理
	GetWorkloadYaml(c *gin.Context, clusterId uint, namespaceName string, workloadType string, workloadName string)
	UpdateWorkloadYaml(c *gin.Context, clusterId uint, namespaceName string, req *model.UpdateWorkloadYAMLRequest)

	// 工作负载Pod管理
	GetWorkloadPods(c *gin.Context, clusterId uint, namespaceName string, workloadType string, workloadName string)

	// 监控相关API
	GetPodMetrics(c *gin.Context, clusterId uint, namespaceName string, podName string)
	GetNodeMetrics(c *gin.Context, clusterId uint, nodeName string)
	GetNamespaceMetrics(c *gin.Context, clusterId uint, namespaceName string)

	// Deployment版本回滚管理
	GetDeploymentHistory(c *gin.Context, clusterId uint, namespaceName string, deploymentName string)
	GetDeploymentRevision(c *gin.Context, clusterId uint, namespaceName string, deploymentName string, revision int64)
	RollbackDeployment(c *gin.Context, clusterId uint, namespaceName string, deploymentName string, req *model.RollbackDeploymentRequest)
	PauseDeployment(c *gin.Context, clusterId uint, namespaceName string, deploymentName string)
	ResumeDeployment(c *gin.Context, clusterId uint, namespaceName string, deploymentName string)
	GetDeploymentRolloutStatus(c *gin.Context, clusterId uint, namespaceName string, deploymentName string)
}

// K8sWorkloadServiceImpl K8s工作负载服务实现
type K8sWorkloadServiceImpl struct {
	clusterDao *dao.KubeClusterDao
}

func NewK8sWorkloadService(db *gorm.DB) IK8sWorkloadService {
	return &K8sWorkloadServiceImpl{
		clusterDao: dao.NewKubeClusterDao(db),
	}
}

// GetWorkloads 获取工作负载列表
func (s *K8sWorkloadServiceImpl) GetWorkloads(c *gin.Context, clusterId uint, namespaceName string, workloadType string, name string) {
	// 获取集群信息
	cluster, err := s.clusterDao.GetByID(clusterId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result.Failed(c, http.StatusNotFound, "集群不存在")
			return
		}
		result.Failed(c, http.StatusInternalServerError, "获取集群信息失败: "+err.Error())
		return
	}

	// 创建K8s客户端
	clientset, err := s.createK8sClient(cluster.Credential)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "连接K8s集群失败: "+err.Error())
		return
	}

	var workloads []model.K8sWorkload

	// 根据工作负载类型获取不同的资源
	switch strings.ToLower(workloadType) {
	case "deployment", "deployments":
		workloads, err = s.getDeployments(clientset, namespaceName)
	case "statefulset", "statefulsets":
		workloads, err = s.getStatefulSets(clientset, namespaceName)
	case "daemonset", "daemonsets":
		workloads, err = s.getDaemonSets(clientset, namespaceName)
	case "job", "jobs":
		workloads, err = s.getJobs(clientset, namespaceName)
	case "cronjob", "cronjobs":
		workloads, err = s.getCronJobs(clientset, namespaceName)
	case "pod", "pods":
		// 获取独立的 Pod（不属于任何工作负载）
		workloads, err = s.getStandalonePods(clientset, namespaceName)
	case "all":
		// 获取所有类型的工作负载
		deployments, _ := s.getDeployments(clientset, namespaceName)
		statefulsets, _ := s.getStatefulSets(clientset, namespaceName)
		daemonsets, _ := s.getDaemonSets(clientset, namespaceName)
		jobs, _ := s.getJobs(clientset, namespaceName)
		cronjobs, _ := s.getCronJobs(clientset, namespaceName)
		standalonePods, _ := s.getStandalonePods(clientset, namespaceName)

		workloads = append(workloads, deployments...)
		workloads = append(workloads, statefulsets...)
		workloads = append(workloads, daemonsets...)
		workloads = append(workloads, jobs...)
		workloads = append(workloads, cronjobs...)
		workloads = append(workloads, standalonePods...)
	default:
		result.Failed(c, http.StatusBadRequest, "不支持的工作负载类型: "+workloadType)
		return
	}

	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "获取工作负载列表失败: "+err.Error())
		return
	}

	autoscalingMap, autoscalingErr := s.listWorkloadAutoscalingSummaries(clientset, namespaceName)
	if autoscalingErr == nil {
		for i := range workloads {
			s.attachWorkloadAutoscalingSummary(&workloads[i], autoscalingMap)
		}
	}

	workloads = filterWorkloadsByName(workloads, name)

	response := model.WorkloadListResponse{
		Workloads: workloads,
		Total:     len(workloads),
	}

	result.Success(c, response)
}

// GetWorkloadDetail 获取工作负载详情
func filterWorkloadsByName(workloads []model.K8sWorkload, name string) []model.K8sWorkload {
	name = strings.ToLower(strings.TrimSpace(name))
	if name == "" {
		return workloads
	}

	filtered := make([]model.K8sWorkload, 0, len(workloads))
	for _, workload := range workloads {
		if strings.Contains(strings.ToLower(workload.Name), name) {
			filtered = append(filtered, workload)
		}
	}

	return filtered
}

func (s *K8sWorkloadServiceImpl) GetWorkloadDetail(c *gin.Context, clusterId uint, namespaceName string, workloadType string, workloadName string) {
	// 获取集群信息
	cluster, err := s.clusterDao.GetByID(clusterId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result.Failed(c, http.StatusNotFound, "集群不存在")
			return
		}
		result.Failed(c, http.StatusInternalServerError, "获取集群信息失败: "+err.Error())
		return
	}

	// 创建K8s客户端
	clientset, err := s.createK8sClient(cluster.Credential)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "连接K8s集群失败: "+err.Error())
		return
	}

	var workloadDetail *model.K8sWorkloadDetail

	// 根据工作负载类型获取详情
	switch strings.ToLower(workloadType) {
	case "deployment", "deployments":
		workloadDetail, err = s.getDeploymentDetail(clientset, namespaceName, workloadName)
	case "statefulset", "statefulsets":
		workloadDetail, err = s.getStatefulSetDetail(clientset, namespaceName, workloadName)
	case "daemonset", "daemonsets":
		workloadDetail, err = s.getDaemonSetDetail(clientset, namespaceName, workloadName)
	case "job", "jobs":
		workloadDetail, err = s.getJobDetail(clientset, namespaceName, workloadName)
	case "cronjob", "cronjobs":
		workloadDetail, err = s.getCronJobDetail(clientset, namespaceName, workloadName)
	default:
		result.Failed(c, http.StatusBadRequest, "不支持的工作负载类型: "+workloadType)
		return
	}

	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "获取工作负载详情失败: "+err.Error())
		return
	}

	if workloadDetail != nil && s.isWorkloadAutoscalingSupported(workloadType) {
		autoscaling, autoscalingErr := s.getSingleWorkloadAutoscalingSummary(clientset, namespaceName, workloadType, workloadName)
		if autoscalingErr == nil {
			s.enrichAutoscalingSummaryWarnings(autoscaling, workloadDetail.Resources)
			workloadDetail.Autoscaling = autoscaling
		}
	}

	result.Success(c, workloadDetail)
}

// GetWorkloadAutoscaling 获取工作负载HPA详情
func (s *K8sWorkloadServiceImpl) GetWorkloadAutoscaling(c *gin.Context, clusterId uint, namespaceName string, workloadType string, workloadName string) {
	response := model.WorkloadAutoscalingResponse{
		Supported:    s.isWorkloadAutoscalingSupported(workloadType),
		Exists:       false,
		WorkloadType: workloadType,
		WorkloadName: workloadName,
		Namespace:    namespaceName,
	}
	if !response.Supported {
		result.Success(c, response)
		return
	}

	cluster, err := s.clusterDao.GetByID(clusterId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result.Failed(c, http.StatusNotFound, "集群不存在")
			return
		}
		result.Failed(c, http.StatusInternalServerError, "获取集群信息失败: "+err.Error())
		return
	}
	if cluster.Status != model.ClusterStatusRunning {
		result.Failed(c, http.StatusBadRequest, "集群状态异常，无法查询HPA")
		return
	}

	clientset, err := s.createK8sClient(cluster.Credential)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "连接K8s集群失败: "+err.Error())
		return
	}

	autoscaling, err := s.getSingleWorkloadAutoscalingSummary(clientset, namespaceName, workloadType, workloadName)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "获取工作负载HPA失败: "+err.Error())
		return
	}
	if autoscaling != nil {
		if resources, resourcesErr := s.getWorkloadResourcesForAutoscaling(clientset, namespaceName, workloadType, workloadName); resourcesErr == nil {
			s.enrichAutoscalingSummaryWarnings(autoscaling, resources)
		}
		response.Exists = true
		response.Autoscaling = autoscaling
	}

	result.Success(c, response)
}

// UpsertWorkloadAutoscaling 创建或更新工作负载HPA
func (s *K8sWorkloadServiceImpl) UpsertWorkloadAutoscaling(c *gin.Context, clusterId uint, namespaceName string, workloadType string, workloadName string, req *model.UpsertWorkloadAutoscalingRequest) {
	if !s.isWorkloadAutoscalingSupported(workloadType) {
		result.Failed(c, http.StatusBadRequest, "当前工作负载类型不支持HPA")
		return
	}
	if err := s.validateUpsertWorkloadAutoscalingRequest(req); err != nil {
		result.Failed(c, http.StatusBadRequest, "HPA配置不合法: "+err.Error())
		return
	}

	cluster, err := s.clusterDao.GetByID(clusterId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result.Failed(c, http.StatusNotFound, "集群不存在")
			return
		}
		result.Failed(c, http.StatusInternalServerError, "获取集群信息失败: "+err.Error())
		return
	}
	if cluster.Status != model.ClusterStatusRunning {
		result.Failed(c, http.StatusBadRequest, "集群状态异常，无法配置HPA")
		return
	}

	clientset, err := s.createK8sClient(cluster.Credential)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "连接K8s集群失败: "+err.Error())
		return
	}

	currentReplicas, err := s.getWorkloadCurrentReplicas(clientset, namespaceName, workloadType, workloadName)
	if err != nil {
		if apierrors.IsNotFound(err) {
			result.Failed(c, http.StatusNotFound, "工作负载不存在: "+err.Error())
			return
		}
		result.Failed(c, http.StatusInternalServerError, "获取工作负载失败: "+err.Error())
		return
	}

	workloadResources, err := s.getWorkloadResourcesForAutoscaling(clientset, namespaceName, workloadType, workloadName)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "获取工作负载资源信息失败: "+err.Error())
		return
	}

	metricSpecs, err := s.buildHPAMetricSpecs(req.Metrics)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "构建HPA指标失败: "+err.Error())
		return
	}

	behavior, err := s.buildHPABehavior(req.Behavior)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "HPA行为策略不合法: "+err.Error())
		return
	}

	hpaClient := clientset.AutoscalingV2().HorizontalPodAutoscalers(namespaceName)
	existingHPA, err := s.findHorizontalPodAutoscalerByTarget(clientset, namespaceName, workloadType, workloadName)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "查询现有HPA失败: "+err.Error())
		return
	}

	minReplicas := req.MinReplicas
	targetRef := autoscalingv2.CrossVersionObjectReference{
		APIVersion: "apps/v1",
		Kind:       s.getWorkloadAutoscalingTargetKind(workloadType),
		Name:       workloadName,
	}

	var action string
	var savedHPA *autoscalingv2.HorizontalPodAutoscaler
	if existingHPA != nil {
		action = "updated"
		existingHPA.Spec.ScaleTargetRef = targetRef
		existingHPA.Spec.MinReplicas = &minReplicas
		existingHPA.Spec.MaxReplicas = req.MaxReplicas
		existingHPA.Spec.Metrics = metricSpecs
		existingHPA.Spec.Behavior = behavior

		savedHPA, err = hpaClient.Update(context.TODO(), existingHPA, metav1.UpdateOptions{})
		if err != nil {
			result.Failed(c, http.StatusInternalServerError, "更新HPA失败: "+err.Error())
			return
		}
	} else {
		action = "created"
		hpa := &autoscalingv2.HorizontalPodAutoscaler{
			ObjectMeta: metav1.ObjectMeta{
				Name:      s.buildWorkloadAutoscalingName(workloadName),
				Namespace: namespaceName,
			},
			Spec: autoscalingv2.HorizontalPodAutoscalerSpec{
				ScaleTargetRef: targetRef,
				MinReplicas:    &minReplicas,
				MaxReplicas:    req.MaxReplicas,
				Metrics:        metricSpecs,
				Behavior:       behavior,
			},
		}

		savedHPA, err = hpaClient.Create(context.TODO(), hpa, metav1.CreateOptions{})
		if err != nil {
			result.Failed(c, http.StatusInternalServerError, "创建HPA失败: "+err.Error())
			return
		}
	}

	autoscalingSummary := s.summarizeHorizontalPodAutoscaler(savedHPA)
	s.enrichAutoscalingSummaryWarnings(autoscalingSummary, workloadResources)

	result.Success(c, map[string]interface{}{
		"action":          action,
		"namespace":       namespaceName,
		"workloadType":    workloadType,
		"workloadName":    workloadName,
		"currentReplicas": currentReplicas,
		"autoscaling":     autoscalingSummary,
		"warnings":        autoscalingSummary.Warnings,
		"message":         fmt.Sprintf("%s 的HPA配置已%s", workloadName, map[string]string{"created": "创建", "updated": "更新"}[action]),
	})
}

// DeleteWorkloadAutoscaling 删除工作负载HPA
func (s *K8sWorkloadServiceImpl) DeleteWorkloadAutoscaling(c *gin.Context, clusterId uint, namespaceName string, workloadType string, workloadName string) {
	if !s.isWorkloadAutoscalingSupported(workloadType) {
		result.Failed(c, http.StatusBadRequest, "当前工作负载类型不支持HPA")
		return
	}

	cluster, err := s.clusterDao.GetByID(clusterId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result.Failed(c, http.StatusNotFound, "集群不存在")
			return
		}
		result.Failed(c, http.StatusInternalServerError, "获取集群信息失败: "+err.Error())
		return
	}
	if cluster.Status != model.ClusterStatusRunning {
		result.Failed(c, http.StatusBadRequest, "集群状态异常，无法删除HPA")
		return
	}

	clientset, err := s.createK8sClient(cluster.Credential)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "连接K8s集群失败: "+err.Error())
		return
	}

	existingHPA, err := s.findHorizontalPodAutoscalerByTarget(clientset, namespaceName, workloadType, workloadName)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "查询现有HPA失败: "+err.Error())
		return
	}
	if existingHPA == nil {
		result.Success(c, map[string]interface{}{
			"deleted":      false,
			"namespace":    namespaceName,
			"workloadType": workloadType,
			"workloadName": workloadName,
			"message":      fmt.Sprintf("%s 当前未配置HPA", workloadName),
		})
		return
	}

	if err := clientset.AutoscalingV2().HorizontalPodAutoscalers(namespaceName).Delete(context.TODO(), existingHPA.Name, metav1.DeleteOptions{}); err != nil {
		result.Failed(c, http.StatusInternalServerError, "删除HPA失败: "+err.Error())
		return
	}

	result.Success(c, map[string]interface{}{
		"deleted":      true,
		"hpaName":      existingHPA.Name,
		"namespace":    namespaceName,
		"workloadType": workloadType,
		"workloadName": workloadName,
		"message":      fmt.Sprintf("%s 的HPA已删除", workloadName),
	})
}

// GetWorkloadCapacitySuggestion 获取工作负载容量建议
func (s *K8sWorkloadServiceImpl) GetWorkloadCapacitySuggestion(c *gin.Context, clusterId uint, namespaceName string, workloadType string, workloadName string) {
	cluster, err := s.clusterDao.GetByID(clusterId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result.Failed(c, http.StatusNotFound, "集群不存在")
			return
		}
		result.Failed(c, http.StatusInternalServerError, "获取集群信息失败: "+err.Error())
		return
	}
	if cluster.Status != model.ClusterStatusRunning {
		result.Failed(c, http.StatusBadRequest, "集群状态异常，无法生成容量建议")
		return
	}

	clientset, err := s.createK8sClient(cluster.Credential)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "连接K8s集群失败: "+err.Error())
		return
	}

	workloadDetail, err := s.getWorkloadDetailForSuggestion(clientset, namespaceName, workloadType, workloadName)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "获取工作负载详情失败: "+err.Error())
		return
	}

	var autoscalingSummary *model.WorkloadAutoscalingSummary
	if s.isWorkloadAutoscalingSupported(workloadType) {
		autoscalingSummary, _ = s.getSingleWorkloadAutoscalingSummary(clientset, namespaceName, workloadType, workloadName)
		if autoscalingSummary != nil {
			s.enrichAutoscalingSummaryWarnings(autoscalingSummary, workloadDetail.Resources)
		}
	}

	suggestion := s.buildWorkloadCapacitySuggestionPayload(clusterId, cluster.Name, namespaceName, workloadType, workloadName, workloadDetail, autoscalingSummary)

	historyID, generatedAt, generatedByID, generatedBy, err := s.persistWorkloadCapacitySuggestion(c, suggestion)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "保存容量建议历史失败: "+err.Error())
		return
	}
	suggestion.HistoryID = historyID
	suggestion.GeneratedAt = generatedAt
	suggestion.GeneratedByID = generatedByID
	suggestion.GeneratedBy = generatedBy

	result.Success(c, suggestion)
}

// GetWorkloadGovernanceOverview 获取工作负载统一治理工作台概览
func (s *K8sWorkloadServiceImpl) GetWorkloadGovernanceOverview(c *gin.Context, clusterId uint, namespaceName string, workloadType string, workloadName string) {
	cluster, err := s.clusterDao.GetByID(clusterId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result.Failed(c, http.StatusNotFound, "集群不存在")
			return
		}
		result.Failed(c, http.StatusInternalServerError, "获取集群信息失败: "+err.Error())
		return
	}
	if cluster.Status != model.ClusterStatusRunning {
		result.Failed(c, http.StatusBadRequest, "集群状态异常，无法生成治理工作台")
		return
	}

	clientset, err := s.createK8sClient(cluster.Credential)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "连接K8s集群失败: "+err.Error())
		return
	}

	workloadDetail, err := s.getWorkloadDetailForSuggestion(clientset, namespaceName, workloadType, workloadName)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "获取工作负载详情失败: "+err.Error())
		return
	}

	var autoscalingSummary *model.WorkloadAutoscalingSummary
	if s.isWorkloadAutoscalingSupported(workloadType) {
		autoscalingSummary, _ = s.getSingleWorkloadAutoscalingSummary(clientset, namespaceName, workloadType, workloadName)
		if autoscalingSummary != nil {
			s.enrichAutoscalingSummaryWarnings(autoscalingSummary, workloadDetail.Resources)
		}
	}

	currentSuggestion := s.buildWorkloadCapacitySuggestionPayload(clusterId, cluster.Name, namespaceName, workloadType, workloadName, workloadDetail, autoscalingSummary)
	latestSnapshot := s.loadLatestWorkloadCapacitySuggestionSnapshot(clusterId, namespaceName, workloadType, workloadName)
	blockingRules, warningRules := s.buildWorkloadGovernanceRules(workloadDetail, autoscalingSummary, currentSuggestion, latestSnapshot)

	overview := model.WorkloadGovernanceOverviewResponse{
		ClusterID:           clusterId,
		ClusterName:         cluster.Name,
		Namespace:           namespaceName,
		WorkloadType:        strings.ToLower(strings.TrimSpace(workloadType)),
		WorkloadName:        workloadName,
		RiskLevel:           firstNonEmpty(currentSuggestion.RiskLevel, "medium"),
		CurrentSuggestion:   currentSuggestion,
		LatestSnapshot:      latestSnapshot,
		CurrentAutoscaling:  autoscalingSummary,
		AlertSummary:        currentSuggestion.AlertSummary,
		WatchMetrics:        currentSuggestion.WatchMetrics,
		Extensions:          s.buildWorkloadGovernanceExtensions(),
		RelatedApplications: s.loadWorkloadGovernanceApplications(clusterId, namespaceName, workloadType, workloadName),
		AlertCenterPath:     currentSuggestion.AlertCenterPath,
		AIDiagnosisPath:     s.buildWorkloadGovernanceAIDiagnosisPath(latestSnapshot, workloadName, namespaceName),
		K8sWorkloadPath:     s.buildWorkloadGovernanceK8sPath(clusterId, namespaceName, workloadType, workloadName, "governance"),
		BlockingRules:       blockingRules,
		WarningRules:        warningRules,
	}
	if len(overview.BlockingRules) > 0 && strings.ToLower(overview.RiskLevel) != "high" {
		overview.RiskLevel = "high"
	}
	overview.Blocking = len(overview.BlockingRules) > 0
	if overview.Blocking {
		overview.BlockingReason = overview.BlockingRules[0].Message
	}
	overview.Warnings = s.collectWorkloadGovernanceWarningMessages(overview.WarningRules)
	overview.Stages = s.buildWorkloadGovernanceStages(overview)

	result.Success(c, overview)
}

// GetWorkloadCapacitySuggestionHistory 获取工作负载容量建议历史
func (s *K8sWorkloadServiceImpl) GetWorkloadCapacitySuggestionHistory(
	c *gin.Context,
	clusterId uint,
	namespaceName string,
	workloadType string,
	workloadName string,
	query model.WorkloadCapacitySuggestionHistoryQuery,
) {
	if query.PageSize < 1 {
		query.PageSize = 10
	}
	if query.PageSize > 50 {
		query.PageSize = 50
	}
	if query.PageNum < 1 {
		query.PageNum = 1
	}

	list, total, err := dao.ListWorkloadCapacitySuggestionHistory(
		clusterId,
		namespaceName,
		workloadType,
		workloadName,
		query.PageNum,
		query.PageSize,
	)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "获取容量建议历史失败: "+err.Error())
		return
	}

	records := make([]model.WorkloadCapacitySuggestionResponse, 0, len(list))
	for _, item := range list {
		records = append(records, s.convertWorkloadCapacitySuggestionHistory(item))
	}

	result.SuccessWithPage(c, records, total, query.PageNum, query.PageSize)
}

func (s *K8sWorkloadServiceImpl) isWorkloadAutoscalingSupported(workloadType string) bool {
	switch s.normalizeAutoscalingWorkloadType(workloadType) {
	case "deployment", "statefulset":
		return true
	default:
		return false
	}
}

func (s *K8sWorkloadServiceImpl) normalizeAutoscalingWorkloadType(workloadType string) string {
	switch strings.ToLower(strings.TrimSpace(workloadType)) {
	case "deployment", "deployments":
		return "deployment"
	case "statefulset", "statefulsets":
		return "statefulset"
	default:
		return strings.ToLower(strings.TrimSpace(workloadType))
	}
}

func (s *K8sWorkloadServiceImpl) getWorkloadAutoscalingTargetKind(workloadType string) string {
	switch s.normalizeAutoscalingWorkloadType(workloadType) {
	case "deployment":
		return "Deployment"
	case "statefulset":
		return "StatefulSet"
	default:
		return ""
	}
}

func (s *K8sWorkloadServiceImpl) getWorkloadAutoscalingTargetKey(workloadType string, workloadName string) string {
	normalizedType := s.normalizeAutoscalingWorkloadType(workloadType)
	name := strings.TrimSpace(workloadName)
	if normalizedType == "" || name == "" {
		return ""
	}
	return normalizedType + "/" + name
}

func (s *K8sWorkloadServiceImpl) buildWorkloadAutoscalingName(workloadName string) string {
	baseName := strings.ToLower(strings.TrimSpace(workloadName))
	if baseName == "" {
		return "workload-hpa"
	}
	baseName = strings.ReplaceAll(baseName, "_", "-")
	const suffix = "-hpa"
	if len(baseName)+len(suffix) <= 63 {
		return baseName + suffix
	}
	trimmed := strings.TrimRight(baseName[:63-len(suffix)], "-")
	if trimmed == "" {
		trimmed = "workload"
	}
	return trimmed + suffix
}

func (s *K8sWorkloadServiceImpl) attachWorkloadAutoscalingSummary(workload *model.K8sWorkload, autoscalingMap map[string]*model.WorkloadAutoscalingSummary) {
	if workload == nil || len(autoscalingMap) == 0 {
		return
	}
	if summary, ok := autoscalingMap[s.getWorkloadAutoscalingTargetKey(string(workload.Type), workload.Name)]; ok {
		s.enrichAutoscalingSummaryWarnings(summary, workload.Resources)
		workload.Autoscaling = summary
	}
}

func (s *K8sWorkloadServiceImpl) listWorkloadAutoscalingSummaries(clientset *kubernetes.Clientset, namespaceName string) (map[string]*model.WorkloadAutoscalingSummary, error) {
	hpaList, err := clientset.AutoscalingV2().HorizontalPodAutoscalers(namespaceName).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	summaryMap := make(map[string]*model.WorkloadAutoscalingSummary, len(hpaList.Items))
	for i := range hpaList.Items {
		hpa := &hpaList.Items[i]
		key := s.getWorkloadAutoscalingTargetKey(hpa.Spec.ScaleTargetRef.Kind, hpa.Spec.ScaleTargetRef.Name)
		if key == "" {
			continue
		}
		summaryMap[key] = s.summarizeHorizontalPodAutoscaler(hpa)
	}

	return summaryMap, nil
}

func (s *K8sWorkloadServiceImpl) getSingleWorkloadAutoscalingSummary(clientset *kubernetes.Clientset, namespaceName string, workloadType string, workloadName string) (*model.WorkloadAutoscalingSummary, error) {
	hpa, err := s.findHorizontalPodAutoscalerByTarget(clientset, namespaceName, workloadType, workloadName)
	if err != nil || hpa == nil {
		return nil, err
	}
	return s.summarizeHorizontalPodAutoscaler(hpa), nil
}

func (s *K8sWorkloadServiceImpl) findHorizontalPodAutoscalerByTarget(clientset *kubernetes.Clientset, namespaceName string, workloadType string, workloadName string) (*autoscalingv2.HorizontalPodAutoscaler, error) {
	targetKey := s.getWorkloadAutoscalingTargetKey(workloadType, workloadName)
	hpaList, err := clientset.AutoscalingV2().HorizontalPodAutoscalers(namespaceName).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for i := range hpaList.Items {
		hpa := &hpaList.Items[i]
		if s.getWorkloadAutoscalingTargetKey(hpa.Spec.ScaleTargetRef.Kind, hpa.Spec.ScaleTargetRef.Name) == targetKey {
			return hpa, nil
		}
	}

	return nil, nil
}

func (s *K8sWorkloadServiceImpl) summarizeHorizontalPodAutoscaler(hpa *autoscalingv2.HorizontalPodAutoscaler) *model.WorkloadAutoscalingSummary {
	if hpa == nil {
		return nil
	}

	minReplicas := int32(1)
	if hpa.Spec.MinReplicas != nil {
		minReplicas = *hpa.Spec.MinReplicas
	}

	summary := &model.WorkloadAutoscalingSummary{
		Enabled:         true,
		Name:            hpa.Name,
		MinReplicas:     minReplicas,
		MaxReplicas:     hpa.Spec.MaxReplicas,
		CurrentReplicas: hpa.Status.CurrentReplicas,
		DesiredReplicas: hpa.Status.DesiredReplicas,
		Metrics:         s.mergeHPAMetrics(hpa.Spec.Metrics, hpa.Status.CurrentMetrics),
		Conditions:      s.convertHPAConditions(hpa.Status.Conditions),
		Behavior:        s.convertHPABehaviorFromK8s(hpa.Spec.Behavior),
	}
	if hpa.Status.LastScaleTime != nil {
		summary.LastScaleTime = hpa.Status.LastScaleTime.Format(time.RFC3339)
	}

	return summary
}

func (s *K8sWorkloadServiceImpl) enrichAutoscalingSummaryWarnings(summary *model.WorkloadAutoscalingSummary, resources model.WorkloadResources) {
	if summary == nil {
		return
	}

	warningSet := make(map[string]struct{})
	addWarning := func(message string) {
		message = strings.TrimSpace(message)
		if message == "" {
			return
		}
		warningSet[message] = struct{}{}
	}

	for _, metric := range summary.Metrics {
		resourceName := strings.ToLower(strings.TrimSpace(metric.ResourceName))
		if metric.TargetType != string(autoscalingv2.UtilizationMetricType) {
			continue
		}

		switch resourceName {
		case "cpu":
			if strings.TrimSpace(resources.Requests.CPU) == "" {
				addWarning("当前工作负载未配置 CPU requests，使用利用率阈值可能无法按预期生效")
			}
		case "memory":
			if strings.TrimSpace(resources.Requests.Memory) == "" {
				addWarning("当前工作负载未配置内存 requests，使用利用率阈值可能无法按预期生效")
			}
		}
	}

	if summary.Enabled && summary.MinReplicas > 0 && summary.MinReplicas == summary.MaxReplicas {
		addWarning("最小副本数与最大副本数相同，HPA 不会产生实际扩缩容")
	}

	if summary.Behavior != nil {
		if summary.Behavior.ScaleUp != nil && summary.Behavior.ScaleUp.SelectPolicy == "Disabled" {
			addWarning("扩容策略已禁用，HPA 不会执行扩容动作")
		}
		if summary.Behavior.ScaleDown != nil && summary.Behavior.ScaleDown.SelectPolicy == "Disabled" {
			addWarning("缩容策略已禁用，HPA 不会执行缩容动作")
		}
	}

	if len(warningSet) == 0 {
		summary.Warnings = nil
		return
	}

	warnings := make([]string, 0, len(warningSet))
	for message := range warningSet {
		warnings = append(warnings, message)
	}
	sort.Strings(warnings)
	summary.Warnings = warnings
}

func (s *K8sWorkloadServiceImpl) mergeHPAMetrics(metricSpecs []autoscalingv2.MetricSpec, metricStatuses []autoscalingv2.MetricStatus) []model.WorkloadAutoscalingMetricStatus {
	if len(metricSpecs) == 0 {
		return nil
	}

	statuses := make([]model.WorkloadAutoscalingMetricStatus, 0, len(metricSpecs))
	for _, metricSpec := range metricSpecs {
		if metricSpec.Type != autoscalingv2.ResourceMetricSourceType || metricSpec.Resource == nil {
			continue
		}

		item := model.WorkloadAutoscalingMetricStatus{
			ResourceName: string(metricSpec.Resource.Name),
			TargetType:   string(metricSpec.Resource.Target.Type),
		}
		switch metricSpec.Resource.Target.Type {
		case autoscalingv2.UtilizationMetricType:
			if metricSpec.Resource.Target.AverageUtilization != nil {
				item.TargetValue = strconv.FormatInt(int64(*metricSpec.Resource.Target.AverageUtilization), 10)
			}
		case autoscalingv2.AverageValueMetricType:
			if metricSpec.Resource.Target.AverageValue != nil {
				item.TargetValue = metricSpec.Resource.Target.AverageValue.String()
			}
		case autoscalingv2.ValueMetricType:
			if metricSpec.Resource.Target.Value != nil {
				item.TargetValue = metricSpec.Resource.Target.Value.String()
			}
		}

		for _, metricStatus := range metricStatuses {
			if metricStatus.Type != autoscalingv2.ResourceMetricSourceType || metricStatus.Resource == nil {
				continue
			}
			if metricStatus.Resource.Name != metricSpec.Resource.Name {
				continue
			}

			if metricStatus.Resource.Current.AverageUtilization != nil {
				item.CurrentAverageUtilization = metricStatus.Resource.Current.AverageUtilization
			}
			if metricStatus.Resource.Current.AverageValue != nil {
				item.CurrentAverageValue = metricStatus.Resource.Current.AverageValue.String()
			} else if metricStatus.Resource.Current.Value != nil {
				item.CurrentAverageValue = metricStatus.Resource.Current.Value.String()
			}
			break
		}

		statuses = append(statuses, item)
	}

	return statuses
}

func (s *K8sWorkloadServiceImpl) convertHPAConditions(conditions []autoscalingv2.HorizontalPodAutoscalerCondition) []model.WorkloadAutoscalingCondition {
	if len(conditions) == 0 {
		return nil
	}

	resultConditions := make([]model.WorkloadAutoscalingCondition, 0, len(conditions))
	for _, condition := range conditions {
		resultConditions = append(resultConditions, model.WorkloadAutoscalingCondition{
			Type:               string(condition.Type),
			Status:             string(condition.Status),
			LastTransitionTime: condition.LastTransitionTime.Format(time.RFC3339),
			Reason:             condition.Reason,
			Message:            condition.Message,
		})
	}
	return resultConditions
}

func (s *K8sWorkloadServiceImpl) convertHPABehaviorFromK8s(behavior *autoscalingv2.HorizontalPodAutoscalerBehavior) *model.WorkloadAutoscalingBehavior {
	if behavior == nil {
		return nil
	}

	return &model.WorkloadAutoscalingBehavior{
		ScaleUp:   s.convertHPABehaviorRuleFromK8s(behavior.ScaleUp),
		ScaleDown: s.convertHPABehaviorRuleFromK8s(behavior.ScaleDown),
	}
}

func (s *K8sWorkloadServiceImpl) convertHPABehaviorRuleFromK8s(rule *autoscalingv2.HPAScalingRules) *model.WorkloadAutoscalingBehaviorRule {
	if rule == nil {
		return nil
	}

	converted := &model.WorkloadAutoscalingBehaviorRule{
		StabilizationWindowSeconds: rule.StabilizationWindowSeconds,
		Policies:                   s.convertHPAScalingPoliciesFromK8s(rule.Policies),
	}
	if rule.SelectPolicy != nil {
		converted.SelectPolicy = string(*rule.SelectPolicy)
	}

	return converted
}

func (s *K8sWorkloadServiceImpl) convertHPAScalingPoliciesFromK8s(policies []autoscalingv2.HPAScalingPolicy) []model.WorkloadAutoscalingScalingPolicy {
	if len(policies) == 0 {
		return nil
	}

	resultPolicies := make([]model.WorkloadAutoscalingScalingPolicy, 0, len(policies))
	for _, policy := range policies {
		resultPolicies = append(resultPolicies, model.WorkloadAutoscalingScalingPolicy{
			Type:          string(policy.Type),
			Value:         policy.Value,
			PeriodSeconds: policy.PeriodSeconds,
		})
	}
	return resultPolicies
}

func (s *K8sWorkloadServiceImpl) buildHPAMetricSpecs(metrics []model.WorkloadAutoscalingMetricSpec) ([]autoscalingv2.MetricSpec, error) {
	resultSpecs := make([]autoscalingv2.MetricSpec, 0, len(metrics))
	for _, metric := range metrics {
		resourceName := strings.ToLower(strings.TrimSpace(metric.ResourceName))
		targetType := strings.ToLower(strings.TrimSpace(metric.TargetType))
		targetValue := strings.TrimSpace(metric.TargetValue)

		resourceMetric := &autoscalingv2.ResourceMetricSource{
			Name: corev1.ResourceName(resourceName),
		}
		switch targetType {
		case "utilization":
			utilization, err := strconv.Atoi(targetValue)
			if err != nil {
				return nil, fmt.Errorf("%s 利用率阈值必须为整数", resourceName)
			}
			averageUtilization := int32(utilization)
			resourceMetric.Target = autoscalingv2.MetricTarget{
				Type:               autoscalingv2.UtilizationMetricType,
				AverageUtilization: &averageUtilization,
			}
		case "averagevalue":
			quantity, err := resource.ParseQuantity(targetValue)
			if err != nil {
				return nil, fmt.Errorf("%s 平均值阈值格式无效: %w", resourceName, err)
			}
			resourceMetric.Target = autoscalingv2.MetricTarget{
				Type:         autoscalingv2.AverageValueMetricType,
				AverageValue: &quantity,
			}
		default:
			return nil, fmt.Errorf("暂不支持的目标类型: %s", metric.TargetType)
		}

		resultSpecs = append(resultSpecs, autoscalingv2.MetricSpec{
			Type:     autoscalingv2.ResourceMetricSourceType,
			Resource: resourceMetric,
		})
	}

	return resultSpecs, nil
}

func (s *K8sWorkloadServiceImpl) buildHPABehavior(behavior *model.WorkloadAutoscalingBehavior) (*autoscalingv2.HorizontalPodAutoscalerBehavior, error) {
	if behavior == nil {
		return nil, nil
	}

	return &autoscalingv2.HorizontalPodAutoscalerBehavior{
		ScaleUp:   s.buildHPABehaviorRule(behavior.ScaleUp),
		ScaleDown: s.buildHPABehaviorRule(behavior.ScaleDown),
	}, nil
}

func (s *K8sWorkloadServiceImpl) buildHPABehaviorRule(rule *model.WorkloadAutoscalingBehaviorRule) *autoscalingv2.HPAScalingRules {
	if rule == nil {
		return nil
	}

	converted := &autoscalingv2.HPAScalingRules{
		StabilizationWindowSeconds: rule.StabilizationWindowSeconds,
		Policies:                   s.buildHPAScalingPolicies(rule.Policies),
	}
	if selectPolicy := strings.TrimSpace(rule.SelectPolicy); selectPolicy != "" {
		policy := autoscalingv2.ScalingPolicySelect(selectPolicy)
		converted.SelectPolicy = &policy
	}

	return converted
}

func (s *K8sWorkloadServiceImpl) buildHPAScalingPolicies(policies []model.WorkloadAutoscalingScalingPolicy) []autoscalingv2.HPAScalingPolicy {
	if len(policies) == 0 {
		return nil
	}

	converted := make([]autoscalingv2.HPAScalingPolicy, 0, len(policies))
	for _, policy := range policies {
		converted = append(converted, autoscalingv2.HPAScalingPolicy{
			Type:          autoscalingv2.HPAScalingPolicyType(policy.Type),
			Value:         policy.Value,
			PeriodSeconds: policy.PeriodSeconds,
		})
	}
	return converted
}

func (s *K8sWorkloadServiceImpl) validateUpsertWorkloadAutoscalingRequest(req *model.UpsertWorkloadAutoscalingRequest) error {
	if req == nil {
		return fmt.Errorf("请求不能为空")
	}
	if req.MinReplicas <= 0 {
		return fmt.Errorf("最小副本数必须大于0")
	}
	if req.MaxReplicas < req.MinReplicas {
		return fmt.Errorf("最大副本数不能小于最小副本数")
	}
	if len(req.Metrics) == 0 {
		return fmt.Errorf("至少需要配置一个资源指标")
	}

	seenResources := make(map[string]struct{}, len(req.Metrics))
	for _, metric := range req.Metrics {
		resourceName := strings.ToLower(strings.TrimSpace(metric.ResourceName))
		targetType := strings.ToLower(strings.TrimSpace(metric.TargetType))
		targetValue := strings.TrimSpace(metric.TargetValue)

		if resourceName != "cpu" && resourceName != "memory" {
			return fmt.Errorf("资源指标仅支持 cpu 或 memory")
		}
		if _, exists := seenResources[resourceName]; exists {
			return fmt.Errorf("资源指标 %s 重复配置", resourceName)
		}
		seenResources[resourceName] = struct{}{}

		switch targetType {
		case "utilization":
			value, err := strconv.Atoi(targetValue)
			if err != nil || value <= 0 {
				return fmt.Errorf("%s 利用率阈值必须为正整数", resourceName)
			}
		case "averagevalue":
			if _, err := resource.ParseQuantity(targetValue); err != nil {
				return fmt.Errorf("%s 平均值阈值格式无效", resourceName)
			}
		default:
			return fmt.Errorf("指标目标类型仅支持 Utilization 或 AverageValue")
		}
	}

	return s.validateHPABehavior(req.Behavior)
}

func (s *K8sWorkloadServiceImpl) validateHPABehavior(behavior *model.WorkloadAutoscalingBehavior) error {
	if behavior == nil {
		return nil
	}

	if err := s.validateHPABehaviorRule("scaleUp", behavior.ScaleUp); err != nil {
		return err
	}
	if err := s.validateHPABehaviorRule("scaleDown", behavior.ScaleDown); err != nil {
		return err
	}
	return nil
}

func (s *K8sWorkloadServiceImpl) getWorkloadDetailForSuggestion(clientset *kubernetes.Clientset, namespaceName string, workloadType string, workloadName string) (*model.K8sWorkloadDetail, error) {
	switch strings.ToLower(strings.TrimSpace(workloadType)) {
	case "deployment", "deployments":
		return s.getDeploymentDetail(clientset, namespaceName, workloadName)
	case "statefulset", "statefulsets":
		return s.getStatefulSetDetail(clientset, namespaceName, workloadName)
	case "daemonset", "daemonsets":
		return s.getDaemonSetDetail(clientset, namespaceName, workloadName)
	case "job", "jobs":
		return s.getJobDetail(clientset, namespaceName, workloadName)
	case "cronjob", "cronjobs":
		return s.getCronJobDetail(clientset, namespaceName, workloadName)
	default:
		return nil, fmt.Errorf("不支持的工作负载类型: %s", workloadType)
	}
}

func (s *K8sWorkloadServiceImpl) getWorkloadResourcesForAutoscaling(clientset *kubernetes.Clientset, namespaceName string, workloadType string, workloadName string) (model.WorkloadResources, error) {
	switch s.normalizeAutoscalingWorkloadType(workloadType) {
	case "deployment":
		deployment, err := clientset.AppsV1().Deployments(namespaceName).Get(context.TODO(), workloadName, metav1.GetOptions{})
		if err != nil {
			return model.WorkloadResources{}, err
		}
		return s.extractPodResources(&deployment.Spec.Template.Spec), nil
	case "statefulset":
		statefulSet, err := clientset.AppsV1().StatefulSets(namespaceName).Get(context.TODO(), workloadName, metav1.GetOptions{})
		if err != nil {
			return model.WorkloadResources{}, err
		}
		return s.extractPodResources(&statefulSet.Spec.Template.Spec), nil
	default:
		return model.WorkloadResources{}, fmt.Errorf("当前工作负载类型不支持HPA")
	}
}

type workloadCapacityAlertEventRow struct {
	Title           string
	Severity        string
	Status          string
	Summary         string
	LastTriggeredAt string
}

type workloadCapacityIncidentRow struct {
	AlertDesc  string
	AlertLevel string
	Status     int
	AlertTime  string
}

type workloadGovernanceAppRow struct {
	AppID        uint
	AppName      string
	AppCode      string
	Environment  string
	DeployConfig string
}

type workloadGovernanceAppDeployConfig struct {
	ClusterID                uint   `json:"cluster_id"`
	Namespace                string `json:"namespace"`
	WorkloadType             string `json:"workload_type"`
	WorkloadName             string `json:"workload_name"`
	ReleaseGovernanceEnabled bool   `json:"release_governance_enabled"`
	PreCheckMode             string `json:"pre_check_mode"`
}

func (s *K8sWorkloadServiceImpl) buildWorkloadCapacitySuggestionPayload(
	clusterId uint,
	clusterName string,
	namespaceName string,
	workloadType string,
	workloadName string,
	workloadDetail *model.K8sWorkloadDetail,
	autoscalingSummary *model.WorkloadAutoscalingSummary,
) model.WorkloadCapacitySuggestionResponse {
	alertSummary := s.buildWorkloadCapacityAlertSummary(namespaceName, workloadName)
	riskLevel, actions, policy, watchMetrics, followUpWindow := s.buildWorkloadCapacityRecommendation(workloadDetail, autoscalingSummary, alertSummary)
	renderedPrompt, systemPrompt := s.renderWorkloadCapacityPrompt(clusterName, namespaceName, workloadType, workloadName, workloadDetail, autoscalingSummary, alertSummary, riskLevel)
	report := s.buildWorkloadCapacityReport(clusterName, namespaceName, workloadType, workloadName, workloadDetail, autoscalingSummary, alertSummary, riskLevel, actions, policy, watchMetrics, followUpWindow)
	alertCenterQuery := s.buildWorkloadCapacityAlertCenterQuery(namespaceName, workloadName)

	return model.WorkloadCapacitySuggestionResponse{
		ClusterID:          clusterId,
		ClusterName:        clusterName,
		Namespace:          namespaceName,
		WorkloadType:       strings.ToLower(strings.TrimSpace(workloadType)),
		WorkloadName:       workloadName,
		TemplateName:       "prediction_suggestion",
		RenderedPrompt:     renderedPrompt,
		SystemPrompt:       systemPrompt,
		Report:             report,
		UsedFallback:       true,
		AlertKeyword:       workloadName,
		AlertCenterPath:    s.buildWorkloadCapacityAlertCenterPath(alertCenterQuery),
		AlertCenterQuery:   alertCenterQuery,
		Autoscaling:        autoscalingSummary,
		AlertSummary:       alertSummary,
		RiskLevel:          riskLevel,
		RecommendedActions: actions,
		RecommendedPolicy:  policy,
		WatchMetrics:       watchMetrics,
		FollowUpWindow:     followUpWindow,
	}
}

func (s *K8sWorkloadServiceImpl) buildWorkloadCapacityAlertCenterQuery(namespaceName, workloadName string) model.WorkloadCapacityAlertCenterQuery {
	return model.WorkloadCapacityAlertCenterQuery{
		Keyword:      strings.TrimSpace(workloadName),
		Namespace:    strings.TrimSpace(namespaceName),
		WorkloadName: strings.TrimSpace(workloadName),
		Source:       "capacity-suggestion",
	}
}

func (s *K8sWorkloadServiceImpl) buildWorkloadCapacityAlertCenterPath(query model.WorkloadCapacityAlertCenterQuery) string {
	values := url.Values{}
	if strings.TrimSpace(query.Keyword) != "" {
		values.Set("keyword", strings.TrimSpace(query.Keyword))
	}
	if strings.TrimSpace(query.Namespace) != "" {
		values.Set("namespace", strings.TrimSpace(query.Namespace))
	}
	if strings.TrimSpace(query.WorkloadName) != "" {
		values.Set("workloadName", strings.TrimSpace(query.WorkloadName))
	}
	if strings.TrimSpace(query.Source) != "" {
		values.Set("source", strings.TrimSpace(query.Source))
	}

	encoded := values.Encode()
	if encoded == "" {
		return "/monitor/alert-center"
	}
	return "/monitor/alert-center?" + encoded
}

func (s *K8sWorkloadServiceImpl) persistWorkloadCapacitySuggestion(
	c *gin.Context,
	suggestion model.WorkloadCapacitySuggestionResponse,
) (uint, string, uint, string, error) {
	var (
		generatedByID uint
		generatedBy   string
	)
	if admin, err := jwt.GetAdmin(c); err == nil && admin != nil {
		generatedByID = admin.ID
		generatedBy = strings.TrimSpace(admin.Username)
	}

	entity := &model.WorkloadCapacitySuggestionHistoryEntity{
		ClusterID:          suggestion.ClusterID,
		ClusterName:        suggestion.ClusterName,
		NamespaceName:      suggestion.Namespace,
		WorkloadType:       strings.ToLower(strings.TrimSpace(suggestion.WorkloadType)),
		WorkloadName:       suggestion.WorkloadName,
		TemplateName:       suggestion.TemplateName,
		RenderedPrompt:     suggestion.RenderedPrompt,
		SystemPrompt:       suggestion.SystemPrompt,
		Report:             suggestion.Report,
		UsedFallback:       suggestion.UsedFallback,
		AlertKeyword:       suggestion.AlertKeyword,
		AlertCenterPath:    suggestion.AlertCenterPath,
		AlertCenterQuery:   s.stringifyJSON(suggestion.AlertCenterQuery),
		RiskLevel:          suggestion.RiskLevel,
		RecommendedActions: s.stringifyJSON(suggestion.RecommendedActions),
		RecommendedPolicy:  s.stringifyJSON(suggestion.RecommendedPolicy),
		WatchMetrics:       s.stringifyJSON(suggestion.WatchMetrics),
		FollowUpWindow:     suggestion.FollowUpWindow,
		Autoscaling:        s.stringifyJSON(suggestion.Autoscaling),
		AlertSummary:       s.stringifyJSON(suggestion.AlertSummary),
		GeneratedByID:      generatedByID,
		GeneratedBy:        generatedBy,
	}
	if err := dao.CreateWorkloadCapacitySuggestionHistory(entity); err != nil {
		return 0, "", 0, "", err
	}
	return entity.ID, entity.CreatedAt.Format("2006-01-02 15:04:05"), entity.GeneratedByID, entity.GeneratedBy, nil
}

func (s *K8sWorkloadServiceImpl) convertWorkloadCapacitySuggestionHistory(
	entity model.WorkloadCapacitySuggestionHistoryEntity,
) model.WorkloadCapacitySuggestionResponse {
	response := model.WorkloadCapacitySuggestionResponse{
		ClusterID:       entity.ClusterID,
		ClusterName:     entity.ClusterName,
		Namespace:       entity.NamespaceName,
		WorkloadType:    entity.WorkloadType,
		WorkloadName:    entity.WorkloadName,
		HistoryID:       entity.ID,
		GeneratedAt:     entity.CreatedAt.Format("2006-01-02 15:04:05"),
		GeneratedByID:   entity.GeneratedByID,
		GeneratedBy:     entity.GeneratedBy,
		TemplateName:    entity.TemplateName,
		RenderedPrompt:  entity.RenderedPrompt,
		SystemPrompt:    entity.SystemPrompt,
		Report:          entity.Report,
		UsedFallback:    entity.UsedFallback,
		AlertKeyword:    entity.AlertKeyword,
		AlertCenterPath: entity.AlertCenterPath,
		RiskLevel:       entity.RiskLevel,
		FollowUpWindow:  entity.FollowUpWindow,
	}
	s.unmarshalJSONString(entity.AlertCenterQuery, &response.AlertCenterQuery)
	s.unmarshalJSONString(entity.Autoscaling, &response.Autoscaling)
	s.unmarshalJSONString(entity.AlertSummary, &response.AlertSummary)
	s.unmarshalJSONString(entity.RecommendedActions, &response.RecommendedActions)
	s.unmarshalJSONString(entity.RecommendedPolicy, &response.RecommendedPolicy)
	s.unmarshalJSONString(entity.WatchMetrics, &response.WatchMetrics)
	return response
}

func (s *K8sWorkloadServiceImpl) loadLatestWorkloadCapacitySuggestionSnapshot(
	clusterId uint,
	namespaceName string,
	workloadType string,
	workloadName string,
) *model.WorkloadCapacitySuggestionResponse {
	var entity model.WorkloadCapacitySuggestionHistoryEntity
	if err := pkgdb.Db.Where(
		"cluster_id = ? AND namespace_name = ? AND workload_type = ? AND workload_name = ?",
		clusterId,
		strings.TrimSpace(namespaceName),
		strings.ToLower(strings.TrimSpace(workloadType)),
		strings.TrimSpace(workloadName),
	).Order("created_at DESC, id DESC").Take(&entity).Error; err != nil {
		return nil
	}
	response := s.convertWorkloadCapacitySuggestionHistory(entity)
	return &response
}

func (s *K8sWorkloadServiceImpl) buildWorkloadGovernanceAIDiagnosisPath(
	latestSnapshot *model.WorkloadCapacitySuggestionResponse,
	workloadName string,
	namespaceName string,
) string {
	if latestSnapshot == nil || latestSnapshot.HistoryID == 0 {
		return ""
	}

	values := url.Values{}
	values.Set("scene", "workload_capacity")
	values.Set("source", "workload-governance")
	values.Set("targetId", strconv.FormatUint(uint64(latestSnapshot.HistoryID), 10))
	values.Set("keyword", strings.TrimSpace(workloadName+" hpa capacity autoscaling"))
	values.Set("templateName", "incident_analysis")
	if strings.TrimSpace(namespaceName) != "" {
		values.Set("namespace", strings.TrimSpace(namespaceName))
	}
	if strings.TrimSpace(workloadName) != "" {
		values.Set("workloadName", strings.TrimSpace(workloadName))
	}
	values.Set("autoRun", "1")
	return "/ai/diagnosis?" + values.Encode()
}

func (s *K8sWorkloadServiceImpl) buildWorkloadGovernanceK8sPath(clusterId uint, namespaceName string, workloadType string, workloadName string, action string) string {
	values := url.Values{}
	values.Set("clusterId", strconv.FormatUint(uint64(clusterId), 10))
	values.Set("namespace", strings.TrimSpace(namespaceName))
	values.Set("type", strings.ToLower(strings.TrimSpace(workloadType)))
	values.Set("name", strings.TrimSpace(workloadName))
	if strings.TrimSpace(action) != "" {
		values.Set("action", strings.TrimSpace(action))
	}
	return "/k8s/workload?" + values.Encode()
}

func (s *K8sWorkloadServiceImpl) buildWorkloadGovernanceAppDetailPath(appId uint, environment string) string {
	if appId == 0 {
		return ""
	}
	values := url.Values{}
	values.Set("appId", strconv.FormatUint(uint64(appId), 10))
	if strings.TrimSpace(environment) != "" {
		values.Set("env", strings.TrimSpace(environment))
	}
	values.Set("focus", "governance")
	return "/app/application?" + values.Encode()
}

func (s *K8sWorkloadServiceImpl) normalizeWorkloadGovernanceAppDeployConfig(cfg workloadGovernanceAppDeployConfig) workloadGovernanceAppDeployConfig {
	cfg.Namespace = strings.TrimSpace(cfg.Namespace)
	cfg.WorkloadType = strings.ToLower(strings.TrimSpace(cfg.WorkloadType))
	cfg.WorkloadName = strings.TrimSpace(cfg.WorkloadName)
	cfg.PreCheckMode = strings.ToLower(strings.TrimSpace(cfg.PreCheckMode))
	if cfg.PreCheckMode == "" {
		cfg.PreCheckMode = "observe"
	}
	return cfg
}

func (s *K8sWorkloadServiceImpl) loadWorkloadGovernanceApplications(clusterId uint, namespaceName string, workloadType string, workloadName string) []model.WorkloadGovernanceAppLink {
	var rows []workloadGovernanceAppRow
	_ = pkgdb.Db.Table("app_jenkins_env env").
		Select(`
			COALESCE(env.app_id, 0) AS app_id,
			COALESCE(app.name, '') AS app_name,
			COALESCE(app.code, '') AS app_code,
			COALESCE(env.env_name, '') AS environment,
			COALESCE(CAST(env.deploy_config AS CHAR), '') AS deploy_config
		`).
		Joins("LEFT JOIN app_application app ON app.id = env.app_id").
		Where("env.deploy_config IS NOT NULL AND env.deploy_config <> ''").
		Scan(&rows).Error

	links := make([]model.WorkloadGovernanceAppLink, 0)
	normalizedType := strings.ToLower(strings.TrimSpace(workloadType))
	normalizedNamespace := strings.TrimSpace(namespaceName)
	normalizedName := strings.TrimSpace(workloadName)
	for _, row := range rows {
		var cfg workloadGovernanceAppDeployConfig
		if strings.TrimSpace(row.DeployConfig) == "" || strings.TrimSpace(row.DeployConfig) == "null" {
			continue
		}
		if err := json.Unmarshal([]byte(row.DeployConfig), &cfg); err != nil {
			continue
		}
		cfg = s.normalizeWorkloadGovernanceAppDeployConfig(cfg)
		if cfg.ClusterID != clusterId ||
			cfg.Namespace != normalizedNamespace ||
			cfg.WorkloadType != normalizedType ||
			cfg.WorkloadName != normalizedName {
			continue
		}
		links = append(links, model.WorkloadGovernanceAppLink{
			AppID:          row.AppID,
			AppName:        strings.TrimSpace(row.AppName),
			AppCode:        strings.TrimSpace(row.AppCode),
			Environment:    strings.TrimSpace(row.Environment),
			PreCheckMode:   cfg.PreCheckMode,
			ReleaseEnabled: cfg.ReleaseGovernanceEnabled,
			DetailPath:     s.buildWorkloadGovernanceAppDetailPath(row.AppID, row.Environment),
		})
	}
	sort.SliceStable(links, func(i, j int) bool {
		if links[i].AppName == links[j].AppName {
			return links[i].Environment < links[j].Environment
		}
		return links[i].AppName < links[j].AppName
	})
	return links
}

func (s *K8sWorkloadServiceImpl) buildWorkloadGovernanceRules(
	workloadDetail *model.K8sWorkloadDetail,
	autoscalingSummary *model.WorkloadAutoscalingSummary,
	currentSuggestion model.WorkloadCapacitySuggestionResponse,
	latestSnapshot *model.WorkloadCapacitySuggestionResponse,
) ([]model.WorkloadGovernanceRuleHit, []model.WorkloadGovernanceRuleHit) {
	blockingRules := make([]model.WorkloadGovernanceRuleHit, 0)
	warningRules := make([]model.WorkloadGovernanceRuleHit, 0)
	addRule := func(target *[]model.WorkloadGovernanceRuleHit, ruleType, code, title, message, source, action string) {
		*target = append(*target, model.WorkloadGovernanceRuleHit{
			Type:    ruleType,
			Code:    code,
			Title:   title,
			Message: message,
			Source:  source,
			Action:  action,
		})
	}

	if currentSuggestion.AlertSummary.CriticalCount > 0 {
		addRule(&blockingRules, "blocking", "critical_alerts", "存在高优先级告警", fmt.Sprintf("当前存在 %d 条高优先级告警或 incident，发布前应先处理。", currentSuggestion.AlertSummary.CriticalCount), "alert-center", "先处理高优先级告警")
	} else if currentSuggestion.AlertSummary.OpenEventCount > 0 {
		addRule(&blockingRules, "blocking", "open_alerts", "存在未恢复告警", fmt.Sprintf("当前存在 %d 条未恢复自动化事件，建议先完成告警闭环后再发布。", currentSuggestion.AlertSummary.OpenEventCount), "alert-center", "先清理未恢复告警")
	}

	if autoscalingSummary == nil {
		addRule(&warningRules, "warning", "hpa_missing", "尚未配置 HPA", "当前工作负载仍依赖手工扩缩容，建议先建立 HPA 基线。", "autoscaling", "先配置 HPA")
	} else {
		if autoscalingSummary.MaxReplicas > 0 && autoscalingSummary.DesiredReplicas >= autoscalingSummary.MaxReplicas {
			addRule(&blockingRules, "blocking", "hpa_ceiling", "HPA 已触顶", "当前期望副本已达到最大副本数，发布前建议先处理容量上限。", "autoscaling", "调整 HPA 最大副本或资源配置")
		}
		for index, warning := range autoscalingSummary.Warnings {
			if strings.TrimSpace(warning) == "" {
				continue
			}
			addRule(&warningRules, "warning", fmt.Sprintf("hpa_warning_%d", index+1), "HPA 治理提示", strings.TrimSpace(warning), "autoscaling", "修正 HPA 策略配置")
		}
	}

	if strings.TrimSpace(workloadDetail.Resources.Requests.CPU) == "" {
		addRule(&warningRules, "warning", "cpu_request_missing", "缺少 CPU requests", "未设置 CPU requests 会降低 HPA 判断和调度精度。", "workload", "补齐 CPU requests")
	}
	if strings.TrimSpace(workloadDetail.Resources.Requests.Memory) == "" {
		addRule(&warningRules, "warning", "memory_request_missing", "缺少内存 requests", "未设置内存 requests 会削弱容量建议和调度保护。", "workload", "补齐内存 requests")
	}
	if latestSnapshot == nil {
		addRule(&warningRules, "warning", "no_history_snapshot", "缺少历史快照", "当前还没有可复用的容量建议历史快照，AI 诊断与历史对比能力受限。", "capacity-history", "先生成容量建议快照")
	}

	return blockingRules, warningRules
}

func (s *K8sWorkloadServiceImpl) collectWorkloadGovernanceWarningMessages(rules []model.WorkloadGovernanceRuleHit) []string {
	resultSet := make(map[string]struct{})
	list := make([]string, 0, len(rules))
	for _, rule := range rules {
		message := strings.TrimSpace(rule.Message)
		if message == "" {
			continue
		}
		if _, exists := resultSet[message]; exists {
			continue
		}
		resultSet[message] = struct{}{}
		list = append(list, message)
	}
	return list
}

func (s *K8sWorkloadServiceImpl) buildWorkloadGovernanceStages(overview model.WorkloadGovernanceOverviewResponse) []model.WorkloadGovernanceStage {
	stages := make([]model.WorkloadGovernanceStage, 0, 4)

	preStatus := "ready"
	preSummary := "已具备发布前治理基线。"
	preActionText := "查看容量建议"
	preActionPath := s.buildWorkloadGovernanceK8sPath(overview.ClusterID, overview.Namespace, overview.WorkloadType, overview.WorkloadName, "capacity")
	if overview.Blocking {
		preStatus = "blocking"
		preSummary = firstNonEmpty(overview.BlockingReason, "存在阻断项，发布前需先处置。")
		preActionText = "处理告警阻断"
		preActionPath = overview.AlertCenterPath
	} else if len(overview.WarningRules) > 0 || strings.ToLower(overview.RiskLevel) == "medium" {
		preStatus = "warning"
		preSummary = "存在需关注项，建议在发布前完成一次治理核查。"
		if overview.CurrentAutoscaling == nil {
			preActionText = "先配置 HPA"
			preActionPath = s.buildWorkloadGovernanceK8sPath(overview.ClusterID, overview.Namespace, overview.WorkloadType, overview.WorkloadName, "autoscaling")
		}
	}
	stages = append(stages, model.WorkloadGovernanceStage{
		Key:               "pre_release",
		Title:             "发布前治理",
		Status:            preStatus,
		Summary:           preSummary,
		Blocking:          overview.Blocking,
		PrimaryActionText: preActionText,
		PrimaryActionPath: preActionPath,
		Items:             append([]string{}, overview.Warnings...),
	})

	autoscalingStatus := "ready"
	autoscalingSummary := "HPA 策略已接入，可持续用作发布前后弹性基线。"
	autoscalingItems := make([]string, 0)
	if overview.CurrentAutoscaling == nil {
		autoscalingStatus = "pending"
		autoscalingSummary = "当前尚未配置 HPA，建议先补齐自动扩缩容基线。"
		autoscalingItems = append(autoscalingItems, "未配置 HPA")
	} else if len(overview.CurrentAutoscaling.Warnings) > 0 {
		autoscalingStatus = "warning"
		autoscalingSummary = fmt.Sprintf("HPA 已接入，但存在 %d 条治理提示。", len(overview.CurrentAutoscaling.Warnings))
		autoscalingItems = append(autoscalingItems, overview.CurrentAutoscaling.Warnings...)
	}
	stages = append(stages, model.WorkloadGovernanceStage{
		Key:               "autoscaling",
		Title:             "弹性策略",
		Status:            autoscalingStatus,
		Summary:           autoscalingSummary,
		PrimaryActionText: "配置 HPA",
		PrimaryActionPath: s.buildWorkloadGovernanceK8sPath(overview.ClusterID, overview.Namespace, overview.WorkloadType, overview.WorkloadName, "autoscaling"),
		Items:             autoscalingItems,
	})

	capacityStatus := "ready"
	capacitySummary := "容量建议处于可观察状态。"
	if strings.ToLower(overview.RiskLevel) == "high" {
		capacityStatus = "blocking"
		capacitySummary = "容量风险较高，建议先消除阻断项后再发布。"
	} else if strings.ToLower(overview.RiskLevel) == "medium" {
		capacityStatus = "warning"
		capacitySummary = "当前存在中风险容量项，建议结合建议动作完成治理。"
	}
	capacityItems := append([]string{}, overview.CurrentSuggestion.WatchMetrics...)
	if len(capacityItems) > 4 {
		capacityItems = capacityItems[:4]
	}
	stages = append(stages, model.WorkloadGovernanceStage{
		Key:               "capacity",
		Title:             "容量建议",
		Status:            capacityStatus,
		Summary:           capacitySummary,
		Blocking:          capacityStatus == "blocking",
		PrimaryActionText: "查看容量建议",
		PrimaryActionPath: s.buildWorkloadGovernanceK8sPath(overview.ClusterID, overview.Namespace, overview.WorkloadType, overview.WorkloadName, "capacity"),
		Items:             capacityItems,
	})

	postStatus := "ready"
	postSummary := "发布后观察链路已具备。"
	postItems := make([]string, 0, len(overview.WatchMetrics))
	postItems = append(postItems, overview.WatchMetrics...)
	if overview.AlertSummary.OpenEventCount > 0 || overview.AlertSummary.IncidentCount > 0 {
		postStatus = "warning"
		postSummary = fmt.Sprintf("发布后需重点观察：未恢复事件 %d，关联 incident %d。", overview.AlertSummary.OpenEventCount, overview.AlertSummary.IncidentCount)
	}
	if len(postItems) > 4 {
		postItems = postItems[:4]
	}
	stages = append(stages, model.WorkloadGovernanceStage{
		Key:               "post_release",
		Title:             "发布后治理",
		Status:            postStatus,
		Summary:           postSummary,
		PrimaryActionText: "打开告警中心",
		PrimaryActionPath: overview.AlertCenterPath,
		Items:             postItems,
	})

	return stages
}

func (s *K8sWorkloadServiceImpl) buildWorkloadGovernanceExtensions() []model.WorkloadGovernanceExtension {
	return []model.WorkloadGovernanceExtension{
		{
			Key:            "keda",
			Title:          "KEDA 扩展位",
			Status:         "reserved",
			Summary:        "当前版本先保留 KEDA 兼容结构，后续可直接补齐触发器与冷却策略。",
			ReservedFields: []string{"enabled", "pollingInterval", "cooldownPeriod", "minReplicaCount", "maxReplicaCount", "triggers"},
		},
		{
			Key:            "vpa",
			Title:          "VPA 扩展位",
			Status:         "reserved",
			Summary:        "当前版本先保留 VPA 数据结构，后续可补资源推荐与 update policy。",
			ReservedFields: []string{"mode", "updatePolicy", "containerPolicies", "recommendation"},
		},
	}
}

func (s *K8sWorkloadServiceImpl) unmarshalJSONString(raw string, target interface{}) {
	if target == nil || strings.TrimSpace(raw) == "" || strings.TrimSpace(raw) == "null" {
		return
	}
	_ = json.Unmarshal([]byte(raw), target)
}

func (s *K8sWorkloadServiceImpl) buildWorkloadCapacityAlertSummary(namespaceName, workloadName string) model.WorkloadCapacityAlertSummary {
	keyword := strings.TrimSpace(workloadName)
	if strings.TrimSpace(namespaceName) != "" {
		keyword = keyword + " " + strings.TrimSpace(namespaceName)
	}
	likeWorkload := "%" + strings.TrimSpace(workloadName) + "%"
	likeNamespace := "%" + strings.TrimSpace(namespaceName) + "%"

	summary := model.WorkloadCapacityAlertSummary{
		RecentEvents:    []model.WorkloadCapacityAlertEvent{},
		RecentIncidents: []model.WorkloadCapacityIncident{},
	}

	var eventRows []workloadCapacityAlertEventRow
	_ = pkgdb.Db.Raw(`
		SELECT
			COALESCE(title, '') AS title,
			COALESCE(severity, '') AS severity,
			COALESCE(status, '') AS status,
			COALESCE(summary, '') AS summary,
			COALESCE(DATE_FORMAT(last_triggered_at, '%Y-%m-%d %H:%i:%s'), '') AS last_triggered_at
		FROM monitor_alert_event
		WHERE resource_name = ? OR resource_name LIKE ? OR labels LIKE ? OR labels LIKE ?
		ORDER BY last_triggered_at DESC, id DESC
		LIMIT 5
	`, workloadName, likeWorkload, likeWorkload, likeNamespace).Scan(&eventRows).Error

	for _, row := range eventRows {
		summary.RecentEvents = append(summary.RecentEvents, model.WorkloadCapacityAlertEvent{
			Title:           row.Title,
			Severity:        row.Severity,
			Status:          row.Status,
			Summary:         row.Summary,
			LastTriggeredAt: row.LastTriggeredAt,
		})
		if row.Status == "open" {
			summary.OpenEventCount++
		}
		if row.Status == "resolved" {
			summary.ResolvedEventCount++
		}
		if isCriticalAlertSeverity(row.Severity) {
			summary.CriticalCount++
		}
	}

	var incidentRows []workloadCapacityIncidentRow
	_ = pkgdb.Db.Raw(`
		SELECT
			COALESCE(alert_desc, '') AS alert_desc,
			COALESCE(alert_level, '') AS alert_level,
			COALESCE(status, 0) AS status,
			COALESCE(DATE_FORMAT(alert_time, '%Y-%m-%d %H:%i:%s'), '') AS alert_time
		FROM monitor_incident
		WHERE alert_desc LIKE ? OR alert_desc LIKE ?
		ORDER BY alert_time DESC, id DESC
		LIMIT 5
	`, likeWorkload, likeNamespace).Scan(&incidentRows).Error

	for _, row := range incidentRows {
		summary.RecentIncidents = append(summary.RecentIncidents, model.WorkloadCapacityIncident{
			AlertDesc:  row.AlertDesc,
			AlertLevel: row.AlertLevel,
			Status:     row.Status,
			AlertTime:  row.AlertTime,
		})
		summary.IncidentCount++
		if isCriticalAlertSeverity(row.AlertLevel) {
			summary.CriticalCount++
		}
	}

	return summary
}

func isCriticalAlertSeverity(severity string) bool {
	switch strings.ToLower(strings.TrimSpace(severity)) {
	case "critical", "p1", "p2", "high":
		return true
	default:
		return false
	}
}

func (s *K8sWorkloadServiceImpl) buildWorkloadCapacityRecommendation(
	workloadDetail *model.K8sWorkloadDetail,
	autoscalingSummary *model.WorkloadAutoscalingSummary,
	alertSummary model.WorkloadCapacityAlertSummary,
) (
	string,
	[]model.WorkloadCapacityRecommendation,
	model.WorkloadCapacityPolicy,
	[]string,
	string,
) {
	recommendations := make([]model.WorkloadCapacityRecommendation, 0, 5)
	addRecommendation := func(priority, action, reason, effect string) {
		recommendations = append(recommendations, model.WorkloadCapacityRecommendation{
			Priority:       priority,
			Action:         action,
			Reason:         reason,
			ExpectedEffect: effect,
		})
	}

	riskLevel := "low"
	followUpWindow := "24h"
	if alertSummary.CriticalCount > 0 || alertSummary.OpenEventCount > 0 {
		riskLevel = "high"
		followUpWindow = "2h"
		addRecommendation("P1", "优先排查关联告警与事件", "当前工作负载存在未恢复或高优先级告警", "避免容量问题继续放大为可用性故障")
	}

	if autoscalingSummary == nil {
		if riskLevel == "low" {
			riskLevel = "medium"
			followUpWindow = "8h"
		}
		addRecommendation("P1", "评估并启用 HPA", "当前工作负载仍依赖手工扩缩容", "在流量波动时自动调节副本，降低人工干预成本")
	} else {
		if len(autoscalingSummary.Warnings) > 0 && riskLevel == "low" {
			riskLevel = "medium"
			followUpWindow = "8h"
		}
		for _, warning := range autoscalingSummary.Warnings {
			addRecommendation("P1", warning, "当前 HPA 治理存在配置缺口", "降低误配置导致的失效或扩缩容冻结风险")
		}
		if autoscalingSummary.MaxReplicas > 0 && autoscalingSummary.DesiredReplicas >= autoscalingSummary.MaxReplicas {
			if riskLevel != "high" {
				riskLevel = "medium"
			}
			addRecommendation("P1", "评估提高 HPA 最大副本数或优化请求资源", "期望副本已触顶最大副本数", "为容量增长留出弹性空间，避免告警持续抖动")
		}
	}

	if strings.TrimSpace(workloadDetail.Resources.Requests.CPU) == "" {
		addRecommendation("P1", "补齐 CPU requests", "当前工作负载缺少 CPU requests，HPA 利用率和调度依据都不稳定", "提升 HPA 计算准确性并改善调度稳定性")
	}
	if strings.TrimSpace(workloadDetail.Resources.Requests.Memory) == "" {
		addRecommendation("P2", "补齐内存 requests", "当前工作负载缺少内存 requests，容量建议和调度保护不足", "降低内存竞争和 OOM 风险")
	}

	if len(recommendations) == 0 {
		addRecommendation("P2", "保持现有 HPA 策略并持续观察", "当前未发现明显容量风险或关联告警", "以最低变更代价维持稳定运行")
	}

	policy := model.WorkloadCapacityPolicy{
		Type: "none",
	}
	if autoscalingSummary != nil {
		policy.Type = "hpa"
		policy.MinReplicas = strconv.FormatInt(int64(autoscalingSummary.MinReplicas), 10)
		policy.MaxReplicas = strconv.FormatInt(int64(autoscalingSummary.MaxReplicas), 10)
		if len(autoscalingSummary.Metrics) > 0 {
			metric := autoscalingSummary.Metrics[0]
			policy.Metric = strings.ToLower(metric.ResourceName)
			policy.TargetUtilization = metric.TargetValue
		}
	} else if strings.TrimSpace(workloadDetail.Resources.Requests.CPU) != "" || strings.TrimSpace(workloadDetail.Resources.Requests.Memory) != "" {
		policy.Type = "manual-scale"
	}

	watchMetricSet := make(map[string]struct{})
	for _, metric := range []string{"pod 副本就绪率", "容器重启次数"} {
		watchMetricSet[metric] = struct{}{}
	}
	if autoscalingSummary != nil {
		for _, metric := range autoscalingSummary.Metrics {
			watchMetricSet[strings.ToUpper(metric.ResourceName)+" 当前值/目标值"] = struct{}{}
		}
	}
	if strings.TrimSpace(workloadDetail.Resources.Requests.CPU) != "" {
		watchMetricSet["CPU requests / 实际使用"] = struct{}{}
	}
	if strings.TrimSpace(workloadDetail.Resources.Requests.Memory) != "" {
		watchMetricSet["内存 requests / 实际使用"] = struct{}{}
	}
	if alertSummary.OpenEventCount > 0 || alertSummary.IncidentCount > 0 {
		watchMetricSet["关联告警数量与恢复时间"] = struct{}{}
	}
	watchMetrics := make([]string, 0, len(watchMetricSet))
	for metric := range watchMetricSet {
		watchMetrics = append(watchMetrics, metric)
	}
	sort.Strings(watchMetrics)

	if len(recommendations) > 5 {
		recommendations = recommendations[:5]
	}

	return riskLevel, recommendations, policy, watchMetrics, followUpWindow
}

func (s *K8sWorkloadServiceImpl) renderWorkloadCapacityPrompt(
	clusterName, namespaceName, workloadType, workloadName string,
	workloadDetail *model.K8sWorkloadDetail,
	autoscalingSummary *model.WorkloadAutoscalingSummary,
	alertSummary model.WorkloadCapacityAlertSummary,
	riskLevel string,
) (string, string) {
	var tmpl aiModel.PromptTemplate
	if err := pkgdb.Db.Where("name = ? AND enabled = 1", "prediction_suggestion").First(&tmpl).Error; err != nil {
		return "", ""
	}

	metricName := "replicas"
	metricType := "replica"
	currentValue := strconv.FormatInt(int64(workloadDetail.ReadyReplicas), 10)
	threshold := "-"
	scalingState := "未配置 HPA"
	if autoscalingSummary != nil {
		scalingState = s.stringifyJSON(autoscalingSummary)
		if len(autoscalingSummary.Metrics) > 0 {
			metric := autoscalingSummary.Metrics[0]
			metricName = strings.ToLower(metric.ResourceName)
			metricType = metric.TargetType
			if metric.CurrentAverageUtilization != nil {
				currentValue = strconv.FormatInt(int64(*metric.CurrentAverageUtilization), 10)
			} else if strings.TrimSpace(metric.CurrentAverageValue) != "" {
				currentValue = metric.CurrentAverageValue
			} else {
				currentValue = strconv.FormatInt(int64(s.effectiveAutoscalingReplicaValue(autoscalingSummary.CurrentReplicas, workloadDetail.Replicas)), 10)
			}
			threshold = metric.TargetValue
		}
	}

	variables := map[string]interface{}{
		"targetType":   "k8s-workload",
		"targetName":   workloadName,
		"clusterName":  clusterName,
		"namespace":    namespaceName,
		"resourceName": workloadName,
		"metricName":   metricName,
		"metricType":   metricType,
		"currentValue": currentValue,
		"threshold":    threshold,
		"trendType":    "current-alert-linked-snapshot",
		"changeRate":   "0",
		"daysLeft":     "0",
		"confidence":   riskLevel,
		"resourceSpec": s.stringifyJSON(workloadDetail.Resources),
		"scalingState": scalingState,
		"predictions":  s.stringifyJSON(alertSummary),
	}

	return s.renderTemplateString(tmpl.Template, variables), s.renderTemplateString(tmpl.SystemPrompt, variables)
}

func (s *K8sWorkloadServiceImpl) buildWorkloadCapacityReport(
	clusterName, namespaceName, workloadType, workloadName string,
	workloadDetail *model.K8sWorkloadDetail,
	autoscalingSummary *model.WorkloadAutoscalingSummary,
	alertSummary model.WorkloadCapacityAlertSummary,
	riskLevel string,
	recommendations []model.WorkloadCapacityRecommendation,
	policy model.WorkloadCapacityPolicy,
	watchMetrics []string,
	followUpWindow string,
) string {
	lines := []string{
		"# 工作负载容量建议报告",
		"",
		fmt.Sprintf("- 集群: %s", clusterName),
		fmt.Sprintf("- 命名空间: %s", namespaceName),
		fmt.Sprintf("- 工作负载: %s / %s", strings.ToLower(workloadType), workloadName),
		fmt.Sprintf("- 当前状态: %s", workloadDetail.Status),
		fmt.Sprintf("- 风险等级: %s", riskLevel),
		"",
		"## HPA 摘要",
	}

	if autoscalingSummary == nil {
		lines = append(lines, "- 当前未配置 HPA")
	} else {
		currentReplicas := s.effectiveAutoscalingReplicaValue(autoscalingSummary.CurrentReplicas, workloadDetail.Replicas)
		desiredReplicas := s.effectiveAutoscalingReplicaValue(autoscalingSummary.DesiredReplicas, workloadDetail.Replicas)
		lines = append(lines,
			fmt.Sprintf("- HPA: %s", autoscalingSummary.Name),
			fmt.Sprintf("- 副本范围: %d - %d", autoscalingSummary.MinReplicas, autoscalingSummary.MaxReplicas),
			fmt.Sprintf("- 当前/期望副本: %d/%d", currentReplicas, desiredReplicas),
		)
		for _, metric := range autoscalingSummary.Metrics {
			currentMetric := metric.CurrentAverageValue
			if metric.CurrentAverageUtilization != nil {
				currentMetric = fmt.Sprintf("%d%%", *metric.CurrentAverageUtilization)
			}
			lines = append(lines, fmt.Sprintf("- 指标 %s: 当前=%s, 目标=%s", strings.ToUpper(metric.ResourceName), firstNonEmpty(currentMetric, "-"), firstNonEmpty(metric.TargetValue, "-")))
		}
		if len(autoscalingSummary.Warnings) > 0 {
			lines = append(lines, "- 治理提示:")
			for _, warning := range autoscalingSummary.Warnings {
				lines = append(lines, fmt.Sprintf("  - %s", warning))
			}
		}
	}

	lines = append(lines,
		"",
		"## 关联告警摘要",
		fmt.Sprintf("- 未恢复自动化事件: %d", alertSummary.OpenEventCount),
		fmt.Sprintf("- 已恢复自动化事件: %d", alertSummary.ResolvedEventCount),
		fmt.Sprintf("- 关联事件工单: %d", alertSummary.IncidentCount),
		fmt.Sprintf("- 高优先级告警数: %d", alertSummary.CriticalCount),
	)
	if len(alertSummary.RecentEvents) > 0 {
		lines = append(lines, "- 最近事件:")
		for _, item := range alertSummary.RecentEvents {
			lines = append(lines, fmt.Sprintf("  - [%s/%s] %s", item.Severity, item.Status, firstNonEmpty(item.Title, item.Summary)))
		}
	} else {
		lines = append(lines, "- 最近事件: 暂无直接命中的自动化告警")
	}

	lines = append(lines, "", "## 建议动作")
	for _, item := range recommendations {
		lines = append(lines, fmt.Sprintf("- %s %s：%s；预期收益：%s", item.Priority, item.Action, item.Reason, item.ExpectedEffect))
	}

	lines = append(lines,
		"",
		"## 建议策略",
		fmt.Sprintf("- 类型: %s", policy.Type),
		fmt.Sprintf("- 最小副本: %s", firstNonEmpty(policy.MinReplicas, "-")),
		fmt.Sprintf("- 最大副本: %s", firstNonEmpty(policy.MaxReplicas, "-")),
		fmt.Sprintf("- 核心指标: %s", firstNonEmpty(policy.Metric, "-")),
		fmt.Sprintf("- 阈值: %s", firstNonEmpty(policy.TargetUtilization, "-")),
		"",
		"## 持续观察",
	)
	for _, metric := range watchMetrics {
		lines = append(lines, fmt.Sprintf("- %s", metric))
	}
	lines = append(lines, "", fmt.Sprintf("## 建议复查窗口\n- %s", followUpWindow))

	return strings.Join(lines, "\n")
}

func (s *K8sWorkloadServiceImpl) renderTemplateString(template string, variables map[string]interface{}) string {
	output := template
	for key, value := range variables {
		placeholder := "{{" + key + "}}"
		output = strings.ReplaceAll(output, placeholder, s.stringifyJSON(value))
	}
	return output
}

func (s *K8sWorkloadServiceImpl) stringifyJSON(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	default:
		data, _ := json.Marshal(v)
		return string(data)
	}
}

func (s *K8sWorkloadServiceImpl) effectiveAutoscalingReplicaValue(current, fallback int32) int32 {
	if current > 0 {
		return current
	}
	if fallback > 0 {
		return fallback
	}
	return 0
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func (s *K8sWorkloadServiceImpl) validateHPABehaviorRule(direction string, rule *model.WorkloadAutoscalingBehaviorRule) error {
	if rule == nil {
		return nil
	}
	if rule.StabilizationWindowSeconds != nil && *rule.StabilizationWindowSeconds < 0 {
		return fmt.Errorf("%s 稳定窗口不能为负数", direction)
	}
	if selectPolicy := strings.TrimSpace(rule.SelectPolicy); selectPolicy != "" {
		switch selectPolicy {
		case "Max", "Min", "Disabled":
		default:
			return fmt.Errorf("%s 选择策略仅支持 Max、Min、Disabled", direction)
		}
	}
	for _, policy := range rule.Policies {
		if policy.Type != "Pods" && policy.Type != "Percent" {
			return fmt.Errorf("%s 限速策略仅支持 Pods 或 Percent", direction)
		}
		if policy.Value <= 0 {
			return fmt.Errorf("%s 限速策略数值必须大于0", direction)
		}
		if policy.PeriodSeconds <= 0 {
			return fmt.Errorf("%s 限速策略周期必须大于0", direction)
		}
	}
	return nil
}

func (s *K8sWorkloadServiceImpl) getWorkloadCurrentReplicas(clientset *kubernetes.Clientset, namespaceName string, workloadType string, workloadName string) (int32, error) {
	switch s.normalizeAutoscalingWorkloadType(workloadType) {
	case "deployment":
		deployment, err := clientset.AppsV1().Deployments(namespaceName).Get(context.TODO(), workloadName, metav1.GetOptions{})
		if err != nil {
			return 0, err
		}
		if deployment.Spec.Replicas == nil {
			return 1, nil
		}
		return *deployment.Spec.Replicas, nil
	case "statefulset":
		statefulSet, err := clientset.AppsV1().StatefulSets(namespaceName).Get(context.TODO(), workloadName, metav1.GetOptions{})
		if err != nil {
			return 0, err
		}
		if statefulSet.Spec.Replicas == nil {
			return 1, nil
		}
		return *statefulSet.Spec.Replicas, nil
	default:
		return 0, fmt.Errorf("当前工作负载类型不支持HPA")
	}
}

// createK8sClient 创建K8s客户端
func (s *K8sWorkloadServiceImpl) createK8sClient(kubeconfig string) (*kubernetes.Clientset, error) {
	config, err := clientcmd.RESTConfigFromKubeConfig([]byte(kubeconfig))
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

// createMetricsClient 创建Metrics客户端
func (s *K8sWorkloadServiceImpl) createMetricsClient(kubeconfig string) (*metricsclientset.Clientset, error) {
	config, err := clientcmd.RESTConfigFromKubeConfig([]byte(kubeconfig))
	if err != nil {
		return nil, err
	}

	metricsClient, err := metricsclientset.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return metricsClient, nil
}

// getDeployments 获取Deployment列表
func (s *K8sWorkloadServiceImpl) getDeployments(clientset *kubernetes.Clientset, namespaceName string) ([]model.K8sWorkload, error) {
	deploymentList, err := clientset.AppsV1().Deployments(namespaceName).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var workloads []model.K8sWorkload
	for _, deployment := range deploymentList.Items {
		workload := s.convertDeploymentToWorkload(&deployment)
		workloads = append(workloads, workload)
	}

	return workloads, nil
}

// getStatefulSets 获取StatefulSet列表
func (s *K8sWorkloadServiceImpl) getStatefulSets(clientset *kubernetes.Clientset, namespaceName string) ([]model.K8sWorkload, error) {
	statefulSetList, err := clientset.AppsV1().StatefulSets(namespaceName).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var workloads []model.K8sWorkload
	for _, statefulSet := range statefulSetList.Items {
		workload := s.convertStatefulSetToWorkload(&statefulSet)
		workloads = append(workloads, workload)
	}

	return workloads, nil
}

// getDaemonSets 获取DaemonSet列表
func (s *K8sWorkloadServiceImpl) getDaemonSets(clientset *kubernetes.Clientset, namespaceName string) ([]model.K8sWorkload, error) {
	daemonSetList, err := clientset.AppsV1().DaemonSets(namespaceName).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var workloads []model.K8sWorkload
	for _, daemonSet := range daemonSetList.Items {
		workload := s.convertDaemonSetToWorkload(&daemonSet)
		workloads = append(workloads, workload)
	}

	return workloads, nil
}

// getJobs 获取Job列表
func (s *K8sWorkloadServiceImpl) getJobs(clientset *kubernetes.Clientset, namespaceName string) ([]model.K8sWorkload, error) {
	jobList, err := clientset.BatchV1().Jobs(namespaceName).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var workloads []model.K8sWorkload
	for _, job := range jobList.Items {
		workload := s.convertJobToWorkload(&job)
		workloads = append(workloads, workload)
	}

	return workloads, nil
}

// getCronJobs 获取CronJob列表
func (s *K8sWorkloadServiceImpl) getCronJobs(clientset *kubernetes.Clientset, namespaceName string) ([]model.K8sWorkload, error) {
	cronJobList, err := clientset.BatchV1().CronJobs(namespaceName).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var workloads []model.K8sWorkload
	for _, cronJob := range cronJobList.Items {
		workload := s.convertCronJobToWorkload(&cronJob)
		workloads = append(workloads, workload)
	}

	return workloads, nil
}

// getStandalonePods 获取独立的Pod列表（不被工作负载管理的Pod）
func (s *K8sWorkloadServiceImpl) getStandalonePods(clientset *kubernetes.Clientset, namespaceName string) ([]model.K8sWorkload, error) {
	// 获取所有Pod
	podList, err := clientset.CoreV1().Pods(namespaceName).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	// 获取所有工作负载，用于检查Pod是否属于某个工作负载
	deployments, _ := clientset.AppsV1().Deployments(namespaceName).List(context.TODO(), metav1.ListOptions{})
	statefulSets, _ := clientset.AppsV1().StatefulSets(namespaceName).List(context.TODO(), metav1.ListOptions{})
	daemonSets, _ := clientset.AppsV1().DaemonSets(namespaceName).List(context.TODO(), metav1.ListOptions{})
	jobs, _ := clientset.BatchV1().Jobs(namespaceName).List(context.TODO(), metav1.ListOptions{})

	// 创建一个映射来跟踪被工作负载管理的Pod
	managedPods := make(map[string]bool)

	// 检查所有工作负载的选择器
	for _, deployment := range deployments.Items {
		if deployment.Spec.Selector != nil && deployment.Spec.Selector.MatchLabels != nil {
			for _, pod := range podList.Items {
				if s.podMatchesSelector(pod.Labels, deployment.Spec.Selector.MatchLabels) {
					managedPods[pod.Name] = true
				}
			}
		}
	}

	for _, statefulSet := range statefulSets.Items {
		if statefulSet.Spec.Selector != nil && statefulSet.Spec.Selector.MatchLabels != nil {
			for _, pod := range podList.Items {
				if s.podMatchesSelector(pod.Labels, statefulSet.Spec.Selector.MatchLabels) {
					managedPods[pod.Name] = true
				}
			}
		}
	}

	for _, daemonSet := range daemonSets.Items {
		if daemonSet.Spec.Selector != nil && daemonSet.Spec.Selector.MatchLabels != nil {
			for _, pod := range podList.Items {
				if s.podMatchesSelector(pod.Labels, daemonSet.Spec.Selector.MatchLabels) {
					managedPods[pod.Name] = true
				}
			}
		}
	}

	for _, job := range jobs.Items {
		if job.Spec.Selector != nil && job.Spec.Selector.MatchLabels != nil {
			for _, pod := range podList.Items {
				if s.podMatchesSelector(pod.Labels, job.Spec.Selector.MatchLabels) {
					managedPods[pod.Name] = true
				}
			}
		}
	}

	// 收集独立的Pod
	var workloads []model.K8sWorkload
	for _, pod := range podList.Items {
		if !managedPods[pod.Name] {
			workload := s.convertPodToWorkload(&pod)
			workloads = append(workloads, workload)
		}
	}

	return workloads, nil
}

// podMatchesSelector 检查Pod标签是否匹配选择器
func (s *K8sWorkloadServiceImpl) podMatchesSelector(podLabels, selectorLabels map[string]string) bool {
	if selectorLabels == nil {
		return false
	}
	for key, value := range selectorLabels {
		if podLabels[key] != value {
			return false
		}
	}
	return true
}

// convertPodToWorkload 转换Pod为工作负载
func (s *K8sWorkloadServiceImpl) convertPodToWorkload(pod *corev1.Pod) model.K8sWorkload {
	return model.K8sWorkload{
		Name:      pod.Name,
		Type:      model.WorkloadTypePod,
		Namespace: pod.Namespace,
		Labels:    pod.Labels,
		Replicas:  1, // Pod总是只有1个副本
		ReadyReplicas: func() int32 {
			if s.getPodStatus(pod) == "Running" {
				return 1
			}
			return 0
		}(),
		Resources: s.extractPodResources(&pod.Spec),
		Images:    s.extractImages(&pod.Spec),
		Status:    s.getPodStatus(pod),
		CreatedAt: pod.CreationTimestamp.Format(time.RFC3339),
		UpdatedAt: pod.CreationTimestamp.Format(time.RFC3339), // Pod通常不会更新，使用创建时间
	}
}

// convertDeploymentToWorkload 转换Deployment为工作负载
func (s *K8sWorkloadServiceImpl) convertDeploymentToWorkload(deployment *appsv1.Deployment) model.K8sWorkload {
	return model.K8sWorkload{
		Name:          deployment.Name,
		Type:          model.WorkloadTypeDeployment,
		Namespace:     deployment.Namespace,
		Labels:        deployment.Labels,
		Replicas:      *deployment.Spec.Replicas,
		ReadyReplicas: deployment.Status.ReadyReplicas,
		Resources:     s.extractPodResources(&deployment.Spec.Template.Spec),
		Images:        s.extractImages(&deployment.Spec.Template.Spec),
		Status:        s.getDeploymentStatus(deployment),
		CreatedAt:     deployment.CreationTimestamp.Format(time.RFC3339),
		UpdatedAt:     s.getLastUpdateTime(deployment.Status.Conditions),
	}
}

// convertStatefulSetToWorkload 转换StatefulSet为工作负载
func (s *K8sWorkloadServiceImpl) convertStatefulSetToWorkload(statefulSet *appsv1.StatefulSet) model.K8sWorkload {
	return model.K8sWorkload{
		Name:          statefulSet.Name,
		Type:          model.WorkloadTypeStatefulSet,
		Namespace:     statefulSet.Namespace,
		Labels:        statefulSet.Labels,
		Replicas:      *statefulSet.Spec.Replicas,
		ReadyReplicas: statefulSet.Status.ReadyReplicas,
		Resources:     s.extractPodResources(&statefulSet.Spec.Template.Spec),
		Images:        s.extractImages(&statefulSet.Spec.Template.Spec),
		Status:        s.getStatefulSetStatus(statefulSet),
		CreatedAt:     statefulSet.CreationTimestamp.Format(time.RFC3339),
		UpdatedAt:     s.getStatefulSetLastUpdateTime(statefulSet.Status.Conditions),
	}
}

// convertDaemonSetToWorkload 转换DaemonSet为工作负载
func (s *K8sWorkloadServiceImpl) convertDaemonSetToWorkload(daemonSet *appsv1.DaemonSet) model.K8sWorkload {
	return model.K8sWorkload{
		Name:          daemonSet.Name,
		Type:          model.WorkloadTypeDaemonSet,
		Namespace:     daemonSet.Namespace,
		Labels:        daemonSet.Labels,
		Replicas:      daemonSet.Status.DesiredNumberScheduled,
		ReadyReplicas: daemonSet.Status.NumberReady,
		Resources:     s.extractPodResources(&daemonSet.Spec.Template.Spec),
		Images:        s.extractImages(&daemonSet.Spec.Template.Spec),
		Status:        s.getDaemonSetStatus(daemonSet),
		CreatedAt:     daemonSet.CreationTimestamp.Format(time.RFC3339),
		UpdatedAt:     s.getDaemonSetLastUpdateTime(daemonSet.Status.Conditions),
	}
}

// convertJobToWorkload 转换Job为工作负载
func (s *K8sWorkloadServiceImpl) convertJobToWorkload(job *batchv1.Job) model.K8sWorkload {
	var replicas int32 = 1
	if job.Spec.Completions != nil {
		replicas = *job.Spec.Completions
	}

	return model.K8sWorkload{
		Name:          job.Name,
		Type:          model.WorkloadTypeJob,
		Namespace:     job.Namespace,
		Labels:        job.Labels,
		Replicas:      replicas,
		ReadyReplicas: job.Status.Succeeded,
		Resources:     s.extractPodResources(&job.Spec.Template.Spec),
		Images:        s.extractImages(&job.Spec.Template.Spec),
		Status:        s.getJobStatus(job),
		CreatedAt:     job.CreationTimestamp.Format(time.RFC3339),
		UpdatedAt:     s.getJobLastUpdateTime(job.Status.Conditions),
	}
}

// convertCronJobToWorkload 转换CronJob为工作负载
func (s *K8sWorkloadServiceImpl) convertCronJobToWorkload(cronJob *batchv1.CronJob) model.K8sWorkload {
	var replicas int32 = 1
	if cronJob.Spec.JobTemplate.Spec.Completions != nil {
		replicas = *cronJob.Spec.JobTemplate.Spec.Completions
	}

	return model.K8sWorkload{
		Name:          cronJob.Name,
		Type:          model.WorkloadTypeCronJob,
		Namespace:     cronJob.Namespace,
		Labels:        cronJob.Labels,
		Replicas:      replicas,
		ReadyReplicas: 0, // CronJob没有常驻副本
		Resources:     s.extractPodResources(&cronJob.Spec.JobTemplate.Spec.Template.Spec),
		Images:        s.extractImages(&cronJob.Spec.JobTemplate.Spec.Template.Spec),
		Status:        s.getCronJobStatus(cronJob),
		CreatedAt:     cronJob.CreationTimestamp.Format(time.RFC3339),
		UpdatedAt:     s.getCronJobLastUpdateTime(cronJob),
	}
}

// extractPodResources 提取Pod资源配置
func (s *K8sWorkloadServiceImpl) extractPodResources(podSpec *corev1.PodSpec) model.WorkloadResources {
	var totalRequests, totalLimits model.ResourceSpec

	// 汇总所有容器的资源
	for _, container := range podSpec.Containers {
		if container.Resources.Requests != nil {
			if cpu, ok := container.Resources.Requests[corev1.ResourceCPU]; ok {
				totalRequests.CPU = s.addResourceQuantity(totalRequests.CPU, cpu.String())
			}
			if memory, ok := container.Resources.Requests[corev1.ResourceMemory]; ok {
				totalRequests.Memory = s.addResourceQuantity(totalRequests.Memory, memory.String())
			}
		}

		if container.Resources.Limits != nil {
			if cpu, ok := container.Resources.Limits[corev1.ResourceCPU]; ok {
				totalLimits.CPU = s.addResourceQuantity(totalLimits.CPU, cpu.String())
			}
			if memory, ok := container.Resources.Limits[corev1.ResourceMemory]; ok {
				totalLimits.Memory = s.addResourceQuantity(totalLimits.Memory, memory.String())
			}
		}
	}

	return model.WorkloadResources{
		Requests: totalRequests,
		Limits:   totalLimits,
	}
}

// extractImages 提取镜像列表
func (s *K8sWorkloadServiceImpl) extractImages(podSpec *corev1.PodSpec) []string {
	var images []string
	for _, container := range podSpec.Containers {
		images = append(images, container.Image)
	}
	return images
}

// addResourceQuantity 添加资源数量（简化处理）
func (s *K8sWorkloadServiceImpl) addResourceQuantity(current, add string) string {
	if current == "" {
		return add
	}
	if add == "" {
		return current
	}

	// 简化处理，实际应该解析并相加资源量
	currentQ, err1 := resource.ParseQuantity(current)
	addQ, err2 := resource.ParseQuantity(add)
	if err1 != nil || err2 != nil {
		return add // 解析失败时返回新值
	}

	currentQ.Add(addQ)
	return currentQ.String()
}

// getDeploymentStatus 获取Deployment状态
func (s *K8sWorkloadServiceImpl) getDeploymentStatus(deployment *appsv1.Deployment) string {
	desiredReplicas := int32(0)
	if deployment.Spec.Replicas != nil {
		desiredReplicas = *deployment.Spec.Replicas
	}

	if desiredReplicas == 0 {
		return "Stopped"
	}
	if deployment.Status.ReadyReplicas == desiredReplicas {
		return "Running"
	}
	if deployment.Status.ReadyReplicas > 0 {
		return "Partial"
	}
	return "Pending"
}

// getStatefulSetStatus 获取StatefulSet状态
func (s *K8sWorkloadServiceImpl) getStatefulSetStatus(statefulSet *appsv1.StatefulSet) string {
	desiredReplicas := int32(0)
	if statefulSet.Spec.Replicas != nil {
		desiredReplicas = *statefulSet.Spec.Replicas
	}

	if desiredReplicas == 0 {
		return "Stopped"
	}
	if statefulSet.Status.ReadyReplicas == desiredReplicas {
		return "Running"
	}
	if statefulSet.Status.ReadyReplicas > 0 {
		return "Partial"
	}
	return "Pending"
}

// getDaemonSetStatus 获取DaemonSet状态
func (s *K8sWorkloadServiceImpl) getDaemonSetStatus(daemonSet *appsv1.DaemonSet) string {
	if daemonSet.Status.DesiredNumberScheduled == 0 {
		return "Pending"
	}
	if daemonSet.Status.NumberReady == daemonSet.Status.DesiredNumberScheduled {
		return "Running"
	}
	if daemonSet.Status.NumberReady > 0 {
		return "Partial"
	}
	return "Pending"
}

// getJobStatus 获取Job状态
func (s *K8sWorkloadServiceImpl) getJobStatus(job *batchv1.Job) string {
	if job.Status.Succeeded > 0 {
		return "Completed"
	}
	if job.Status.Failed > 0 {
		return "Failed"
	}
	if job.Status.Active > 0 {
		return "Running"
	}
	return "Pending"
}

// getCronJobStatus 获取CronJob状态
func (s *K8sWorkloadServiceImpl) getCronJobStatus(cronJob *batchv1.CronJob) string {
	if cronJob.Spec.Suspend != nil && *cronJob.Spec.Suspend {
		return "Suspended"
	}
	if len(cronJob.Status.Active) > 0 {
		return "Running"
	}
	return "Active"
}

// getLastUpdateTime 获取最后更新时间
func (s *K8sWorkloadServiceImpl) getLastUpdateTime(conditions []appsv1.DeploymentCondition) string {
	if len(conditions) == 0 {
		return ""
	}

	var latestTime metav1.Time
	for _, condition := range conditions {
		if condition.LastUpdateTime.After(latestTime.Time) {
			latestTime = condition.LastUpdateTime
		}
	}

	return latestTime.Format(time.RFC3339)
}

// getJobLastUpdateTime 获取Job最后更新时间
func (s *K8sWorkloadServiceImpl) getJobLastUpdateTime(conditions []batchv1.JobCondition) string {
	if len(conditions) == 0 {
		return ""
	}

	var latestTime metav1.Time
	for _, condition := range conditions {
		if condition.LastTransitionTime.After(latestTime.Time) {
			latestTime = condition.LastTransitionTime
		}
	}

	return latestTime.Format(time.RFC3339)
}

// getStatefulSetLastUpdateTime 获取StatefulSet最后更新时间
func (s *K8sWorkloadServiceImpl) getStatefulSetLastUpdateTime(conditions []appsv1.StatefulSetCondition) string {
	if len(conditions) == 0 {
		return ""
	}

	var latestTime metav1.Time
	for _, condition := range conditions {
		if condition.LastTransitionTime.After(latestTime.Time) {
			latestTime = condition.LastTransitionTime
		}
	}

	return latestTime.Format(time.RFC3339)
}

// getDaemonSetLastUpdateTime 获取DaemonSet最后更新时间
func (s *K8sWorkloadServiceImpl) getDaemonSetLastUpdateTime(conditions []appsv1.DaemonSetCondition) string {
	if len(conditions) == 0 {
		return ""
	}

	var latestTime metav1.Time
	for _, condition := range conditions {
		if condition.LastTransitionTime.After(latestTime.Time) {
			latestTime = condition.LastTransitionTime
		}
	}

	return latestTime.Format(time.RFC3339)
}

// getCronJobLastUpdateTime 获取CronJob最后更新时间
func (s *K8sWorkloadServiceImpl) getCronJobLastUpdateTime(cronJob *batchv1.CronJob) string {
	if cronJob.Status.LastScheduleTime != nil {
		return cronJob.Status.LastScheduleTime.Format(time.RFC3339)
	}
	return cronJob.CreationTimestamp.Format(time.RFC3339)
}

// ===================== 工作负载详情获取方法 =====================

// getDeploymentDetail 获取Deployment详情
func (s *K8sWorkloadServiceImpl) getDeploymentDetail(clientset *kubernetes.Clientset, namespaceName, deploymentName string) (*model.K8sWorkloadDetail, error) {
	deployment, err := clientset.AppsV1().Deployments(namespaceName).Get(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	// 获取关联的Pod列表
	pods, err := s.getWorkloadPods(clientset, namespaceName, deployment.Spec.Selector)
	if err != nil {
		pods = []model.K8sPodInfo{} // 获取失败时返回空列表
	}

	// 获取相关事件
	events, err := s.getWorkloadEvents(clientset, namespaceName, "Deployment", deploymentName)
	if err != nil {
		events = []model.K8sEvent{} // 获取失败时返回空列表
	}

	workload := s.convertDeploymentToWorkload(deployment)

	return &model.K8sWorkloadDetail{
		K8sWorkload: workload,
		Pods:        pods,
		Conditions:  s.convertDeploymentConditions(deployment.Status.Conditions),
		Events:      events,
		Spec:        deployment.Spec,
	}, nil
}

// getStatefulSetDetail 获取StatefulSet详情
func (s *K8sWorkloadServiceImpl) getStatefulSetDetail(clientset *kubernetes.Clientset, namespaceName, statefulSetName string) (*model.K8sWorkloadDetail, error) {
	statefulSet, err := clientset.AppsV1().StatefulSets(namespaceName).Get(context.TODO(), statefulSetName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	pods, err := s.getWorkloadPods(clientset, namespaceName, statefulSet.Spec.Selector)
	if err != nil {
		pods = []model.K8sPodInfo{}
	}

	events, err := s.getWorkloadEvents(clientset, namespaceName, "StatefulSet", statefulSetName)
	if err != nil {
		events = []model.K8sEvent{}
	}

	workload := s.convertStatefulSetToWorkload(statefulSet)

	return &model.K8sWorkloadDetail{
		K8sWorkload: workload,
		Pods:        pods,
		Conditions:  s.convertStatefulSetConditions(statefulSet.Status.Conditions),
		Events:      events,
		Spec:        statefulSet.Spec,
	}, nil
}

// getDaemonSetDetail 获取DaemonSet详情
func (s *K8sWorkloadServiceImpl) getDaemonSetDetail(clientset *kubernetes.Clientset, namespaceName, daemonSetName string) (*model.K8sWorkloadDetail, error) {
	daemonSet, err := clientset.AppsV1().DaemonSets(namespaceName).Get(context.TODO(), daemonSetName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	pods, err := s.getWorkloadPods(clientset, namespaceName, daemonSet.Spec.Selector)
	if err != nil {
		pods = []model.K8sPodInfo{}
	}

	events, err := s.getWorkloadEvents(clientset, namespaceName, "DaemonSet", daemonSetName)
	if err != nil {
		events = []model.K8sEvent{}
	}

	workload := s.convertDaemonSetToWorkload(daemonSet)

	return &model.K8sWorkloadDetail{
		K8sWorkload: workload,
		Pods:        pods,
		Conditions:  s.convertDaemonSetConditions(daemonSet.Status.Conditions),
		Events:      events,
		Spec:        daemonSet.Spec,
	}, nil
}

// getJobDetail 获取Job详情
func (s *K8sWorkloadServiceImpl) getJobDetail(clientset *kubernetes.Clientset, namespaceName, jobName string) (*model.K8sWorkloadDetail, error) {
	job, err := clientset.BatchV1().Jobs(namespaceName).Get(context.TODO(), jobName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	pods, err := s.getWorkloadPods(clientset, namespaceName, job.Spec.Selector)
	if err != nil {
		pods = []model.K8sPodInfo{}
	}

	events, err := s.getWorkloadEvents(clientset, namespaceName, "Job", jobName)
	if err != nil {
		events = []model.K8sEvent{}
	}

	workload := s.convertJobToWorkload(job)

	return &model.K8sWorkloadDetail{
		K8sWorkload: workload,
		Pods:        pods,
		Conditions:  s.convertJobConditions(job.Status.Conditions),
		Events:      events,
		Spec:        job.Spec,
	}, nil
}

// getCronJobDetail 获取CronJob详情
func (s *K8sWorkloadServiceImpl) getCronJobDetail(clientset *kubernetes.Clientset, namespaceName, cronJobName string) (*model.K8sWorkloadDetail, error) {
	cronJob, err := clientset.BatchV1().CronJobs(namespaceName).Get(context.TODO(), cronJobName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	// CronJob的Pod是通过其创建的Job获取的
	pods := []model.K8sPodInfo{}

	events, err := s.getWorkloadEvents(clientset, namespaceName, "CronJob", cronJobName)
	if err != nil {
		events = []model.K8sEvent{}
	}

	workload := s.convertCronJobToWorkload(cronJob)

	return &model.K8sWorkloadDetail{
		K8sWorkload: workload,
		Pods:        pods,
		Conditions:  []model.WorkloadCondition{}, // CronJob没有标准的Conditions
		Events:      events,
		Spec:        cronJob.Spec,
	}, nil
}

// ===================== Pod管理方法 =====================

// GetPods 获取Pod列表
func (s *K8sWorkloadServiceImpl) GetPods(c *gin.Context, clusterId uint, namespaceName string) {
	cluster, err := s.clusterDao.GetByID(clusterId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result.Failed(c, http.StatusNotFound, "集群不存在")
			return
		}
		result.Failed(c, http.StatusInternalServerError, "获取集群信息失败: "+err.Error())
		return
	}

	clientset, err := s.createK8sClient(cluster.Credential)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "连接K8s集群失败: "+err.Error())
		return
	}

	podList, err := clientset.CoreV1().Pods(namespaceName).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "获取Pod列表失败: "+err.Error())
		return
	}

	var pods []model.K8sPodInfo
	for _, pod := range podList.Items {
		podInfo := s.convertToPodInfo(&pod)
		pods = append(pods, podInfo)
	}

	result.Success(c, map[string]interface{}{
		"pods":  pods,
		"total": len(pods),
	})
}

// GetPodDetail 获取Pod详情
func (s *K8sWorkloadServiceImpl) GetPodDetail(c *gin.Context, clusterId uint, namespaceName, podName string) {
	cluster, err := s.clusterDao.GetByID(clusterId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result.Failed(c, http.StatusNotFound, "集群不存在")
			return
		}
		result.Failed(c, http.StatusInternalServerError, "获取集群信息失败: "+err.Error())
		return
	}

	clientset, err := s.createK8sClient(cluster.Credential)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "连接K8s集群失败: "+err.Error())
		return
	}

	pod, err := clientset.CoreV1().Pods(namespaceName).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		result.Failed(c, http.StatusNotFound, "Pod不存在: "+err.Error())
		return
	}

	// 获取Pod事件
	events, err := s.getWorkloadEvents(clientset, namespaceName, "Pod", podName)
	if err != nil {
		events = []model.K8sEvent{}
	}

	podInfo := s.convertToPodInfo(pod)
	podDetail := model.K8sPodDetail{
		K8sPodInfo: podInfo,
		Conditions: s.convertPodConditions(pod.Status.Conditions),
		Events:     events,
		Volumes:    s.convertPodVolumes(pod),
		Spec:       pod.Spec,
	}

	result.Success(c, podDetail)
}

// DeletePod 删除Pod
func (s *K8sWorkloadServiceImpl) DeletePod(c *gin.Context, clusterId uint, namespaceName, podName string) {
	cluster, err := s.clusterDao.GetByID(clusterId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result.Failed(c, http.StatusNotFound, "集群不存在")
			return
		}
		result.Failed(c, http.StatusInternalServerError, "获取集群信息失败: "+err.Error())
		return
	}

	clientset, err := s.createK8sClient(cluster.Credential)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "连接K8s集群失败: "+err.Error())
		return
	}

	err = clientset.CoreV1().Pods(namespaceName).Delete(context.TODO(), podName, metav1.DeleteOptions{})
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "删除Pod失败: "+err.Error())
		return
	}

	result.Success(c, "Pod删除成功")
}

// GetPodLogs 获取Pod日志
func (s *K8sWorkloadServiceImpl) GetPodLogs(c *gin.Context, clusterId uint, namespaceName, podName, containerName string) {
	cluster, err := s.clusterDao.GetByID(clusterId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result.Failed(c, http.StatusNotFound, "集群不存在")
			return
		}
		result.Failed(c, http.StatusInternalServerError, "获取集群信息失败: "+err.Error())
		return
	}

	clientset, err := s.createK8sClient(cluster.Credential)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "连接K8s集群失败: "+err.Error())
		return
	}

	// 构建日志请求选项
	logOptions := &corev1.PodLogOptions{
		Container: containerName,
		Follow:    false,
		TailLines: func() *int64 { lines := int64(500); return &lines }(), // 默认最后500行
	}

	// 获取日志
	req := clientset.CoreV1().Pods(namespaceName).GetLogs(podName, logOptions)
	podLogs, err := req.Stream(context.TODO())
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "获取Pod日志失败: "+err.Error())
		return
	}
	defer podLogs.Close()

	// 读取日志内容
	buf := make([]byte, 2048)
	var logs strings.Builder
	for {
		numBytes, err := podLogs.Read(buf)
		if numBytes == 0 {
			break
		}
		if err != nil {
			break
		}
		logs.Write(buf[:numBytes])
	}

	result.Success(c, map[string]interface{}{
		"logs": logs.String(),
	})
}

// GetPodEvents 获取Pod事件
func (s *K8sWorkloadServiceImpl) GetPodEvents(c *gin.Context, clusterId uint, namespaceName, podName string) {
	cluster, err := s.clusterDao.GetByID(clusterId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result.Failed(c, http.StatusNotFound, "集群不存在")
			return
		}
		result.Failed(c, http.StatusInternalServerError, "获取集群信息失败: "+err.Error())
		return
	}

	clientset, err := s.createK8sClient(cluster.Credential)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "连接K8s集群失败: "+err.Error())
		return
	}

	// 首先验证Pod是否存在
	_, err = clientset.CoreV1().Pods(namespaceName).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		result.Failed(c, http.StatusNotFound, "Pod不存在: "+err.Error())
		return
	}

	// 获取Pod相关的事件
	events, err := s.getWorkloadEvents(clientset, namespaceName, "Pod", podName)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "获取Pod事件失败: "+err.Error())
		return
	}

	result.Success(c, map[string]interface{}{
		"events": events,
		"total":  len(events),
	})
}

// GetPodYaml 获取Pod的YAML配置
func (s *K8sWorkloadServiceImpl) GetPodYaml(c *gin.Context, clusterId uint, namespaceName, podName string) {
	cluster, err := s.clusterDao.GetByID(clusterId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result.Failed(c, http.StatusNotFound, "集群不存在")
			return
		}
		result.Failed(c, http.StatusInternalServerError, "获取集群信息失败: "+err.Error())
		return
	}

	clientset, err := s.createK8sClient(cluster.Credential)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "连接K8s集群失败: "+err.Error())
		return
	}

	pod, err := clientset.CoreV1().Pods(namespaceName).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		result.Failed(c, http.StatusNotFound, "Pod不存在: "+err.Error())
		return
	}

	result.Success(c, map[string]interface{}{
		"yaml": pod,
	})
}

// UpdatePodYaml 更新Pod的YAML配置
func (s *K8sWorkloadServiceImpl) UpdatePodYaml(c *gin.Context, clusterId uint, namespaceName, podName string, req *model.UpdatePodYAMLRequest) {
	// 获取集群信息
	cluster, err := s.clusterDao.GetByID(clusterId)
	if err != nil {
		result.Failed(c, http.StatusNotFound, "集群不存在")
		return
	}

	// 创建k8s客户端
	config, err := clientcmd.RESTConfigFromKubeConfig([]byte(cluster.Credential))
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "解析集群配置失败: "+err.Error())
		return
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "创建k8s客户端失败: "+err.Error())
		return
	}

	// 首先校验YAML格式
	validationResult, err := s.validateYAMLContent(req.YAMLContent, "pod")
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "YAML校验失败: "+err.Error())
		return
	}

	// 如果只是校验格式，直接返回校验结果
	if req.ValidateOnly {
		result.Success(c, model.UpdatePodYAMLResponse{
			Success:          validationResult.Valid,
			PodName:          podName,
			Namespace:        namespaceName,
			Message:          "YAML校验完成",
			ValidationResult: validationResult,
		})
		return
	}

	// 解析YAML为Pod对象
	newPod, err := s.parseYAMLToPod(req.YAMLContent)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "解析YAML失败: "+err.Error())
		return
	}

	// 验证Pod名称是否匹配
	if newPod.Name != podName {
		result.Failed(c, http.StatusBadRequest, fmt.Sprintf("YAML中的Pod名称(%s)与URL参数不匹配(%s)", newPod.Name, podName))
		return
	}

	// 验证命名空间是否匹配
	if newPod.Namespace != "" && newPod.Namespace != namespaceName {
		result.Failed(c, http.StatusBadRequest, fmt.Sprintf("YAML中的命名空间(%s)与URL参数不匹配(%s)", newPod.Namespace, namespaceName))
		return
	}

	// 确保命名空间设置正确
	if newPod.Namespace == "" {
		newPod.Namespace = namespaceName
	}

	ctx := context.TODO()

	// 获取现有的Pod
	existingPod, err := clientset.CoreV1().Pods(namespaceName).Get(ctx, podName, metav1.GetOptions{})
	if err != nil {
		result.Failed(c, http.StatusNotFound, "Pod不存在: "+err.Error())
		return
	}

	// 准备响应对象
	response := &model.UpdatePodYAMLResponse{
		Success:   false,
		PodName:   podName,
		Namespace: namespaceName,
		Changes:   []string{},
		Warnings:  []string{},
	}

	// 如果是DryRun模式，只返回分析结果
	if req.DryRun {
		changes := s.analyzePodChanges(existingPod, newPod)
		response.Success = true
		response.Message = "DryRun模式：仅分析变更，未实际更新"
		response.UpdateStrategy = s.determinePodUpdateStrategy(existingPod, newPod)
		response.Changes = changes
		response.ValidationResult = validationResult

		if response.UpdateStrategy == "recreate" {
			response.Warnings = append(response.Warnings, "此更新需要删除并重新创建Pod，会导致服务中断")
		}

		result.Success(c, response)
		return
	}

	// 检查Pod是否可以原地更新
	updateStrategy := s.determinePodUpdateStrategy(existingPod, newPod)
	response.UpdateStrategy = updateStrategy

	switch updateStrategy {
	case "patch":
		// 尝试原地更新（只能更新部分字段）
		err = s.patchPod(clientset, ctx, existingPod, newPod)
		if err != nil {
			result.Failed(c, http.StatusInternalServerError, "原地更新Pod失败: "+err.Error())
			return
		}
		response.Success = true
		response.Message = "Pod原地更新成功"
		response.Changes = []string{"已更新可变字段（如labels、annotations等）"}

	case "recreate":
		// 需要删除重建
		if !req.Force {
			result.Failed(c, http.StatusBadRequest, "此更新需要删除并重新创建Pod，请设置force=true以确认操作")
			return
		}

		// 记录原Pod的状态信息
		response.Warnings = append(response.Warnings, "Pod已被删除并重新创建，可能导致数据丢失")

		// 删除原Pod
		err = clientset.CoreV1().Pods(namespaceName).Delete(ctx, podName, metav1.DeleteOptions{})
		if err != nil {
			result.Failed(c, http.StatusInternalServerError, "删除原Pod失败: "+err.Error())
			return
		}

		// 等待Pod被完全删除
		time.Sleep(time.Second * 2)

		// 创建新Pod
		// 清除一些系统字段
		newPod.ResourceVersion = ""
		newPod.UID = ""
		newPod.CreationTimestamp = metav1.Time{}
		newPod.Status = corev1.PodStatus{}

		_, err = clientset.CoreV1().Pods(namespaceName).Create(ctx, newPod, metav1.CreateOptions{})
		if err != nil {
			result.Failed(c, http.StatusInternalServerError, "重新创建Pod失败: "+err.Error())
			return
		}

		response.Success = true
		response.Message = "Pod已删除并重新创建"
		response.Changes = []string{"Pod已完全重新创建"}

	default:
		result.Failed(c, http.StatusBadRequest, "无法确定更新策略")
		return
	}

	result.Success(c, response)
}

// ===================== 辅助转换方法 =====================

// getWorkloadPods 根据选择器获取工作负载的Pod列表
func (s *K8sWorkloadServiceImpl) getWorkloadPods(clientset *kubernetes.Clientset, namespaceName string, selector *metav1.LabelSelector) ([]model.K8sPodInfo, error) {
	// 将LabelSelector转换为字符串
	selectorStr := ""
	if selector != nil && selector.MatchLabels != nil {
		var selectors []string
		for key, value := range selector.MatchLabels {
			selectors = append(selectors, fmt.Sprintf("%s=%s", key, value))
		}
		selectorStr = strings.Join(selectors, ",")
	}

	podList, err := clientset.CoreV1().Pods(namespaceName).List(context.TODO(), metav1.ListOptions{
		LabelSelector: selectorStr,
	})
	if err != nil {
		return nil, err
	}

	var pods []model.K8sPodInfo
	for _, pod := range podList.Items {
		podInfo := s.convertToPodInfo(&pod)
		pods = append(pods, podInfo)
	}

	return pods, nil
}

// getWorkloadEvents 获取工作负载相关的事件
func (s *K8sWorkloadServiceImpl) getWorkloadEvents(clientset *kubernetes.Clientset, namespaceName, kind, name string) ([]model.K8sEvent, error) {
	events, err := clientset.CoreV1().Events(namespaceName).List(context.TODO(), metav1.ListOptions{
		FieldSelector: fmt.Sprintf("involvedObject.kind=%s,involvedObject.name=%s", kind, name),
	})
	if err != nil {
		return nil, err
	}

	var k8sEvents []model.K8sEvent
	for _, event := range events.Items {
		k8sEvent := model.K8sEvent{
			Type:      event.Type,
			Reason:    event.Reason,
			Message:   event.Message,
			Source:    event.Source.Component,
			Count:     event.Count,
			FirstTime: event.FirstTimestamp.Format(time.RFC3339),
			LastTime:  event.LastTimestamp.Format(time.RFC3339),
		}
		k8sEvents = append(k8sEvents, k8sEvent)
	}

	return k8sEvents, nil
}

// convertToPodInfo 转换Pod为PodInfo
func (s *K8sWorkloadServiceImpl) convertToPodInfo(pod *corev1.Pod) model.K8sPodInfo {
	// 计算重启次数
	var restartCount int32
	for _, containerStatus := range pod.Status.ContainerStatuses {
		restartCount += containerStatus.RestartCount
	}

	// 计算运行时间
	runningTime := ""
	if pod.Status.StartTime != nil {
		runningTime = time.Since(pod.Status.StartTime.Time).String()
	}

	// 转换容器信息
	var containers []model.ContainerInfo
	for i, container := range pod.Spec.Containers {
		containerInfo := model.ContainerInfo{
			Name:      container.Name,
			Image:     container.Image,
			Resources: s.convertContainerResources(container.Resources),
			Ports:     s.convertContainerPorts(container.Ports),
			Env:       s.convertEnvVars(container.Env),
		}

		// 如果有容器状态信息，添加状态相关信息
		if i < len(pod.Status.ContainerStatuses) {
			containerStatus := pod.Status.ContainerStatuses[i]
			containerInfo.Ready = containerStatus.Ready
			containerInfo.RestartCount = containerStatus.RestartCount
			containerInfo.State = s.getContainerState(containerStatus.State)
		}

		containers = append(containers, containerInfo)
	}

	return model.K8sPodInfo{
		Name:         pod.Name,
		Status:       s.getPodStatus(pod),
		Phase:        string(pod.Status.Phase),
		RestartCount: restartCount,
		NodeName:     pod.Spec.NodeName,
		PodIP:        pod.Status.PodIP,
		HostIP:       pod.Status.HostIP,
		Resources:    s.extractPodResources(&pod.Spec),
		RunningTime:  runningTime,
		CreatedAt:    pod.CreationTimestamp.Format(time.RFC3339),
		Labels:       pod.Labels,
		Containers:   containers,
	}
}

// convertContainerResources 转换容器资源
func (s *K8sWorkloadServiceImpl) convertContainerResources(resources corev1.ResourceRequirements) model.WorkloadResources {
	var requests, limits model.ResourceSpec

	if resources.Requests != nil {
		if cpu, ok := resources.Requests[corev1.ResourceCPU]; ok {
			requests.CPU = cpu.String()
		}
		if memory, ok := resources.Requests[corev1.ResourceMemory]; ok {
			requests.Memory = memory.String()
		}
	}

	if resources.Limits != nil {
		if cpu, ok := resources.Limits[corev1.ResourceCPU]; ok {
			limits.CPU = cpu.String()
		}
		if memory, ok := resources.Limits[corev1.ResourceMemory]; ok {
			limits.Memory = memory.String()
		}
	}

	return model.WorkloadResources{
		Requests: requests,
		Limits:   limits,
	}
}

// convertContainerPorts 转换容器端口
func (s *K8sWorkloadServiceImpl) convertContainerPorts(ports []corev1.ContainerPort) []model.ContainerPort {
	var containerPorts []model.ContainerPort
	for _, port := range ports {
		containerPorts = append(containerPorts, model.ContainerPort{
			Name:          port.Name,
			ContainerPort: port.ContainerPort,
			Protocol:      string(port.Protocol),
		})
	}
	return containerPorts
}

// convertEnvVars 转换环境变量
func (s *K8sWorkloadServiceImpl) convertEnvVars(envVars []corev1.EnvVar) []model.EnvVar {
	var modelEnvVars []model.EnvVar
	for _, envVar := range envVars {
		value := envVar.Value
		// 如果是从其他来源引用的值，显示引用信息
		if envVar.ValueFrom != nil {
			if envVar.ValueFrom.SecretKeyRef != nil {
				value = fmt.Sprintf("Secret:%s/%s", envVar.ValueFrom.SecretKeyRef.Name, envVar.ValueFrom.SecretKeyRef.Key)
			} else if envVar.ValueFrom.ConfigMapKeyRef != nil {
				value = fmt.Sprintf("ConfigMap:%s/%s", envVar.ValueFrom.ConfigMapKeyRef.Name, envVar.ValueFrom.ConfigMapKeyRef.Key)
			}
		}
		modelEnvVars = append(modelEnvVars, model.EnvVar{
			Name:  envVar.Name,
			Value: value,
		})
	}
	return modelEnvVars
}

// getContainerState 获取容器状态
func (s *K8sWorkloadServiceImpl) getContainerState(state corev1.ContainerState) string {
	if state.Running != nil {
		return "Running"
	}
	if state.Waiting != nil {
		return "Waiting"
	}
	if state.Terminated != nil {
		return "Terminated"
	}
	return "Unknown"
}

// getPodStatus 获取Pod状态
func (s *K8sWorkloadServiceImpl) getPodStatus(pod *corev1.Pod) string {
	// 检查Pod是否正在删除
	if pod.DeletionTimestamp != nil {
		return "Terminating"
	}

	// 根据Phase和Conditions确定状态
	switch pod.Status.Phase {
	case corev1.PodRunning:
		// 检查所有容器是否准备就绪
		for _, condition := range pod.Status.Conditions {
			if condition.Type == corev1.PodReady && condition.Status == corev1.ConditionTrue {
				return "Running"
			}
		}
		return "NotReady"
	case corev1.PodSucceeded:
		return "Completed"
	case corev1.PodFailed:
		return "Failed"
	case corev1.PodPending:
		return "Pending"
	default:
		return string(pod.Status.Phase)
	}
}

// convertDeploymentConditions 转换Deployment条件
func (s *K8sWorkloadServiceImpl) convertDeploymentConditions(conditions []appsv1.DeploymentCondition) []model.WorkloadCondition {
	var workloadConditions []model.WorkloadCondition
	for _, condition := range conditions {
		workloadConditions = append(workloadConditions, model.WorkloadCondition{
			Type:               string(condition.Type),
			Status:             string(condition.Status),
			LastTransitionTime: condition.LastTransitionTime.Format(time.RFC3339),
			Reason:             condition.Reason,
			Message:            condition.Message,
		})
	}
	return workloadConditions
}

// convertStatefulSetConditions 转换StatefulSet条件
func (s *K8sWorkloadServiceImpl) convertStatefulSetConditions(conditions []appsv1.StatefulSetCondition) []model.WorkloadCondition {
	var workloadConditions []model.WorkloadCondition
	for _, condition := range conditions {
		workloadConditions = append(workloadConditions, model.WorkloadCondition{
			Type:               string(condition.Type),
			Status:             string(condition.Status),
			LastTransitionTime: condition.LastTransitionTime.Format(time.RFC3339),
			Reason:             condition.Reason,
			Message:            condition.Message,
		})
	}
	return workloadConditions
}

// convertDaemonSetConditions 转换DaemonSet条件
func (s *K8sWorkloadServiceImpl) convertDaemonSetConditions(conditions []appsv1.DaemonSetCondition) []model.WorkloadCondition {
	var workloadConditions []model.WorkloadCondition
	for _, condition := range conditions {
		workloadConditions = append(workloadConditions, model.WorkloadCondition{
			Type:               string(condition.Type),
			Status:             string(condition.Status),
			LastTransitionTime: condition.LastTransitionTime.Format(time.RFC3339),
			Reason:             condition.Reason,
			Message:            condition.Message,
		})
	}
	return workloadConditions
}

// convertJobConditions 转换Job条件
func (s *K8sWorkloadServiceImpl) convertJobConditions(conditions []batchv1.JobCondition) []model.WorkloadCondition {
	var workloadConditions []model.WorkloadCondition
	for _, condition := range conditions {
		workloadConditions = append(workloadConditions, model.WorkloadCondition{
			Type:               string(condition.Type),
			Status:             string(condition.Status),
			LastTransitionTime: condition.LastTransitionTime.Format(time.RFC3339),
			Reason:             condition.Reason,
			Message:            condition.Message,
		})
	}
	return workloadConditions
}

// convertPodConditions 转换Pod条件
func (s *K8sWorkloadServiceImpl) convertPodConditions(conditions []corev1.PodCondition) []model.PodCondition {
	var podConditions []model.PodCondition
	for _, condition := range conditions {
		podConditions = append(podConditions, model.PodCondition{
			Type:               string(condition.Type),
			Status:             string(condition.Status),
			LastTransitionTime: condition.LastTransitionTime.Format(time.RFC3339),
			Reason:             condition.Reason,
			Message:            condition.Message,
		})
	}
	return podConditions
}

// convertPodVolumes 转换Pod存储卷信息
func (s *K8sWorkloadServiceImpl) convertPodVolumes(pod *corev1.Pod) []model.VolumeInfo {
	var volumes []model.VolumeInfo

	for _, volume := range pod.Spec.Volumes {
		volumeInfo := model.VolumeInfo{
			Name: volume.Name,
			Type: s.getVolumeType(volume.VolumeSource),
		}

		// 找到对应的挂载信息
		for _, container := range pod.Spec.Containers {
			for _, mount := range container.VolumeMounts {
				if mount.Name == volume.Name {
					volumeInfo.MountPath = mount.MountPath
					volumeInfo.ReadOnly = mount.ReadOnly
					break
				}
			}
		}

		volumes = append(volumes, volumeInfo)
	}

	return volumes
}

// getVolumeType 获取存储卷类型
func (s *K8sWorkloadServiceImpl) getVolumeType(volumeSource corev1.VolumeSource) string {
	switch {
	case volumeSource.EmptyDir != nil:
		return "EmptyDir"
	case volumeSource.HostPath != nil:
		return "HostPath"
	case volumeSource.PersistentVolumeClaim != nil:
		return "PersistentVolumeClaim"
	case volumeSource.Secret != nil:
		return "Secret"
	case volumeSource.ConfigMap != nil:
		return "ConfigMap"
	case volumeSource.NFS != nil:
		return "NFS"
	default:
		return "Unknown"
	}
}

// ===================== 暂未实现的接口方法 =====================

// CreateDeployment 创建Deployment
func (s *K8sWorkloadServiceImpl) CreateDeployment(c *gin.Context, clusterId uint, namespaceName string, req *model.CreateDeploymentRequest) {
	cluster, err := s.clusterDao.GetByID(clusterId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result.Failed(c, http.StatusNotFound, "闆嗙兢涓嶅瓨鍦?")
			return
		}
		result.Failed(c, http.StatusInternalServerError, "鑾峰彇闆嗙兢淇℃伅澶辫触: "+err.Error())
		return
	}

	if cluster.Status != model.ClusterStatusRunning {
		result.Failed(c, http.StatusBadRequest, "闆嗙兢鐘舵€佸紓甯革紝鏃犳硶鍒涘缓Deployment")
		return
	}

	clientset, err := s.createK8sClient(cluster.Credential)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "杩炴帴K8s闆嗙兢澶辫触: "+err.Error())
		return
	}

	deployment, err := s.buildDeploymentForCreate(namespaceName, req)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "鏋勫缓Deployment澶辫触: "+err.Error())
		return
	}

	createdDeployment, err := clientset.AppsV1().Deployments(namespaceName).Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		switch {
		case apierrors.IsAlreadyExists(err):
			result.Failed(c, http.StatusConflict, "Deployment宸插瓨鍦? "+err.Error())
		case apierrors.IsInvalid(err), apierrors.IsBadRequest(err):
			result.Failed(c, http.StatusBadRequest, "Deployment閰嶇疆鏃犳晥: "+err.Error())
		default:
			result.Failed(c, http.StatusInternalServerError, "鍒涘缓Deployment澶辫触: "+err.Error())
		}
		return
	}

	result.Success(c, s.convertDeploymentToWorkload(createdDeployment))
}

func (s *K8sWorkloadServiceImpl) buildDeploymentForCreate(namespaceName string, req *model.CreateDeploymentRequest) (*appsv1.Deployment, error) {
	if req == nil {
		return nil, fmt.Errorf("璇锋眰鍙傛暟涓嶈兘涓虹┖")
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, fmt.Errorf("Deployment鍚嶇О涓嶈兘涓虹┖")
	}

	if len(req.Template.Containers) == 0 {
		return nil, fmt.Errorf("Pod妯℃澘鑷冲皯闇€瑕佷竴涓鍣?")
	}

	replicas := req.Replicas
	if replicas < 0 {
		return nil, fmt.Errorf("鍓湰鏁颁笉鑳戒负璐熸暟")
	}
	if replicas == 0 {
		replicas = 1
	}

	selectorLabels := s.buildWorkloadSelectorLabels(name, req.Labels, req.Template.Labels)
	templateLabels := s.mergeStringMaps(req.Labels, req.Template.Labels)
	templateLabels = s.mergeStringMaps(templateLabels, selectorLabels)
	deploymentLabels := s.mergeStringMaps(req.Labels, selectorLabels)

	podTemplate, err := s.buildPodTemplateForCreate(templateLabels, req.Template)
	if err != nil {
		return nil, err
	}

	strategy := s.buildDeploymentStrategy(req.Strategy)

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespaceName,
			Labels:    deploymentLabels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: selectorLabels,
			},
			Template: *podTemplate,
			Strategy: strategy,
		},
	}, nil
}

func (s *K8sWorkloadServiceImpl) buildStatefulSetForCreate(namespaceName string, req *model.CreateStatefulSetRequest) (*appsv1.StatefulSet, error) {
	if req == nil {
		return nil, fmt.Errorf("璇锋眰鍙傛暟涓嶈兘涓虹┖")
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, fmt.Errorf("StatefulSet鍚嶇О涓嶈兘涓虹┖")
	}

	serviceName := strings.TrimSpace(req.ServiceName)
	if serviceName == "" {
		return nil, fmt.Errorf("StatefulSet蹇呴』鎸囧畾serviceName")
	}

	if len(req.Template.Containers) == 0 {
		return nil, fmt.Errorf("Pod妯℃澘鑷冲皯闇€瑕佷竴涓鍣?")
	}

	replicas := req.Replicas
	if replicas < 0 {
		return nil, fmt.Errorf("鍓湰鏁颁笉鑳戒负璐熸暟")
	}
	if replicas == 0 {
		replicas = 1
	}

	selectorLabels := s.buildWorkloadSelectorLabels(name, req.Labels, req.Template.Labels)
	templateLabels := s.mergeStringMaps(req.Labels, req.Template.Labels)
	templateLabels = s.mergeStringMaps(templateLabels, selectorLabels)
	statefulSetLabels := s.mergeStringMaps(req.Labels, selectorLabels)

	podTemplate, err := s.buildPodTemplateForCreate(templateLabels, req.Template)
	if err != nil {
		return nil, err
	}

	return &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespaceName,
			Labels:    statefulSetLabels,
		},
		Spec: appsv1.StatefulSetSpec{
			ServiceName: serviceName,
			Replicas:    &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: selectorLabels,
			},
			Template: *podTemplate,
		},
	}, nil
}

func (s *K8sWorkloadServiceImpl) buildDaemonSetForCreate(namespaceName string, req *model.CreateDaemonSetRequest) (*appsv1.DaemonSet, error) {
	if req == nil {
		return nil, fmt.Errorf("璇锋眰鍙傛暟涓嶈兘涓虹┖")
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, fmt.Errorf("DaemonSet鍚嶇О涓嶈兘涓虹┖")
	}

	if len(req.Template.Containers) == 0 {
		return nil, fmt.Errorf("Pod妯℃澘鑷冲皯闇€瑕佷竴涓鍣?")
	}

	selectorLabels := s.buildWorkloadSelectorLabels(name, req.Labels, req.Template.Labels)
	templateLabels := s.mergeStringMaps(req.Labels, req.Template.Labels)
	templateLabels = s.mergeStringMaps(templateLabels, selectorLabels)
	daemonSetLabels := s.mergeStringMaps(req.Labels, selectorLabels)

	podTemplate, err := s.buildPodTemplateForCreate(templateLabels, req.Template)
	if err != nil {
		return nil, err
	}

	return &appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespaceName,
			Labels:    daemonSetLabels,
		},
		Spec: appsv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: selectorLabels,
			},
			Template: *podTemplate,
		},
	}, nil
}

func (s *K8sWorkloadServiceImpl) buildJobSpecForCreate(template model.PodTemplateSpec, completions, parallelism int32) (*batchv1.JobSpec, error) {
	templateLabels := s.copyStringMap(template.Labels)
	podTemplate, err := s.buildPodTemplateForCreate(templateLabels, template)
	if err != nil {
		return nil, err
	}

	podTemplate.Spec.RestartPolicy = corev1.RestartPolicyNever

	jobSpec := &batchv1.JobSpec{
		Template: *podTemplate,
	}
	if completions > 0 {
		jobSpec.Completions = &completions
	}
	if parallelism > 0 {
		jobSpec.Parallelism = &parallelism
	}

	return jobSpec, nil
}

func (s *K8sWorkloadServiceImpl) buildJobForCreate(namespaceName string, req *model.CreateJobRequest) (*batchv1.Job, error) {
	if req == nil {
		return nil, fmt.Errorf("璇锋眰鍙傛暟涓嶈兘涓虹┖")
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, fmt.Errorf("Job鍚嶇О涓嶈兘涓虹┖")
	}

	if len(req.Template.Containers) == 0 {
		return nil, fmt.Errorf("Pod妯℃澘鑷冲皯闇€瑕佷竴涓鍣?")
	}

	jobSpec, err := s.buildJobSpecForCreate(req.Template, req.Completions, req.Parallelism)
	if err != nil {
		return nil, err
	}

	jobLabels := s.mergeStringMaps(req.Labels, req.Template.Labels)

	return &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespaceName,
			Labels:    jobLabels,
		},
		Spec: *jobSpec,
	}, nil
}

func (s *K8sWorkloadServiceImpl) buildCronJobForCreate(namespaceName string, req *model.CreateCronJobRequest) (*batchv1.CronJob, error) {
	if req == nil {
		return nil, fmt.Errorf("璇锋眰鍙傛暟涓嶈兘涓虹┖")
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, fmt.Errorf("CronJob鍚嶇О涓嶈兘涓虹┖")
	}

	schedule := strings.TrimSpace(req.Schedule)
	if schedule == "" {
		return nil, fmt.Errorf("CronJob蹇呴』鎸囧畾schedule")
	}

	if len(req.JobTemplate.Template.Containers) == 0 {
		return nil, fmt.Errorf("JobTemplate涓殑Pod妯℃澘鑷冲皯闇€瑕佷竴涓鍣?")
	}

	jobSpec, err := s.buildJobSpecForCreate(req.JobTemplate.Template, req.JobTemplate.Completions, req.JobTemplate.Parallelism)
	if err != nil {
		return nil, err
	}

	jobLabels := s.mergeStringMaps(req.Labels, req.JobTemplate.Labels)

	return &batchv1.CronJob{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespaceName,
			Labels:    jobLabels,
		},
		Spec: batchv1.CronJobSpec{
			Schedule: schedule,
			JobTemplate: batchv1.JobTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: jobLabels,
				},
				Spec: *jobSpec,
			},
		},
	}, nil
}

func (s *K8sWorkloadServiceImpl) buildPodTemplateForCreate(defaultLabels map[string]string, template model.PodTemplateSpec) (*corev1.PodTemplateSpec, error) {
	containers, err := s.buildContainersForCreate(template.Containers)
	if err != nil {
		return nil, err
	}

	volumes, err := s.buildVolumesForCreate(template.Volumes)
	if err != nil {
		return nil, err
	}

	tolerations := s.buildTolerationsForCreate(template.Tolerations)

	return &corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Labels: s.copyStringMap(defaultLabels),
		},
		Spec: corev1.PodSpec{
			Containers:   containers,
			Volumes:      volumes,
			NodeSelector: s.copyStringMap(template.NodeSelector),
			Tolerations:  tolerations,
		},
	}, nil
}

func (s *K8sWorkloadServiceImpl) buildContainersForCreate(containerSpecs []model.ContainerSpec) ([]corev1.Container, error) {
	containers := make([]corev1.Container, 0, len(containerSpecs))
	for _, containerSpec := range containerSpecs {
		containerName := strings.TrimSpace(containerSpec.Name)
		if containerName == "" {
			return nil, fmt.Errorf("瀹瑰櫒鍚嶇О涓嶈兘涓虹┖")
		}

		image := strings.TrimSpace(containerSpec.Image)
		if image == "" {
			return nil, fmt.Errorf("瀹瑰櫒 '%s' 缂哄皯闀滃儚閰嶇疆", containerName)
		}

		resources, err := s.buildResourceRequirementsForCreate(containerSpec.Resources)
		if err != nil {
			return nil, fmt.Errorf("瀹瑰櫒 '%s' 璧勬簮閰嶇疆鏃犳晥: %w", containerName, err)
		}

		ports := make([]corev1.ContainerPort, 0, len(containerSpec.Ports))
		for _, portSpec := range containerSpec.Ports {
			protocol := corev1.ProtocolTCP
			if strings.TrimSpace(portSpec.Protocol) != "" {
				protocol = corev1.Protocol(strings.TrimSpace(portSpec.Protocol))
			}

			ports = append(ports, corev1.ContainerPort{
				Name:          portSpec.Name,
				ContainerPort: portSpec.ContainerPort,
				Protocol:      protocol,
			})
		}

		envVars := make([]corev1.EnvVar, 0, len(containerSpec.Env))
		for _, envVar := range containerSpec.Env {
			envVars = append(envVars, corev1.EnvVar{
				Name:  envVar.Name,
				Value: envVar.Value,
			})
		}

		volumeMounts := make([]corev1.VolumeMount, 0, len(containerSpec.VolumeMounts))
		for _, volumeMountSpec := range containerSpec.VolumeMounts {
			volumeMounts = append(volumeMounts, corev1.VolumeMount{
				Name:      volumeMountSpec.Name,
				MountPath: volumeMountSpec.MountPath,
				ReadOnly:  volumeMountSpec.ReadOnly,
			})
		}

		containers = append(containers, corev1.Container{
			Name:         containerName,
			Image:        image,
			Ports:        ports,
			Env:          envVars,
			Resources:    resources,
			VolumeMounts: volumeMounts,
			Command:      containerSpec.Command,
			Args:         containerSpec.Args,
		})
	}

	return containers, nil
}

func (s *K8sWorkloadServiceImpl) buildResourceRequirementsForCreate(workloadResources model.WorkloadResources) (corev1.ResourceRequirements, error) {
	resources := corev1.ResourceRequirements{}

	if workloadResources.Limits.CPU != "" || workloadResources.Limits.Memory != "" {
		resources.Limits = make(corev1.ResourceList)
		if workloadResources.Limits.CPU != "" {
			cpuQuantity, err := resource.ParseQuantity(workloadResources.Limits.CPU)
			if err != nil {
				return corev1.ResourceRequirements{}, fmt.Errorf("CPU limits '%s' 瑙ｆ瀽澶辫触: %w", workloadResources.Limits.CPU, err)
			}
			resources.Limits[corev1.ResourceCPU] = cpuQuantity
		}
		if workloadResources.Limits.Memory != "" {
			memoryQuantity, err := resource.ParseQuantity(workloadResources.Limits.Memory)
			if err != nil {
				return corev1.ResourceRequirements{}, fmt.Errorf("Memory limits '%s' 瑙ｆ瀽澶辫触: %w", workloadResources.Limits.Memory, err)
			}
			resources.Limits[corev1.ResourceMemory] = memoryQuantity
		}
	}

	if workloadResources.Requests.CPU != "" || workloadResources.Requests.Memory != "" {
		resources.Requests = make(corev1.ResourceList)
		if workloadResources.Requests.CPU != "" {
			cpuQuantity, err := resource.ParseQuantity(workloadResources.Requests.CPU)
			if err != nil {
				return corev1.ResourceRequirements{}, fmt.Errorf("CPU requests '%s' 瑙ｆ瀽澶辫触: %w", workloadResources.Requests.CPU, err)
			}
			resources.Requests[corev1.ResourceCPU] = cpuQuantity
		}
		if workloadResources.Requests.Memory != "" {
			memoryQuantity, err := resource.ParseQuantity(workloadResources.Requests.Memory)
			if err != nil {
				return corev1.ResourceRequirements{}, fmt.Errorf("Memory requests '%s' 瑙ｆ瀽澶辫触: %w", workloadResources.Requests.Memory, err)
			}
			resources.Requests[corev1.ResourceMemory] = memoryQuantity
		}
	}

	return resources, nil
}

func (s *K8sWorkloadServiceImpl) buildVolumesForCreate(volumeSpecs []model.VolumeSpec) ([]corev1.Volume, error) {
	volumes := make([]corev1.Volume, 0, len(volumeSpecs))
	for _, volumeSpec := range volumeSpecs {
		volumeName := strings.TrimSpace(volumeSpec.Name)
		if volumeName == "" {
			return nil, fmt.Errorf("Volume鍚嶇О涓嶈兘涓虹┖")
		}

		volume := corev1.Volume{Name: volumeName}
		config := volumeSpec.Config

		switch strings.ToLower(strings.TrimSpace(volumeSpec.Type)) {
		case "emptydir":
			volume.VolumeSource = corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{},
			}
		case "hostpath":
			path, _ := config["path"].(string)
			path = strings.TrimSpace(path)
			if path == "" {
				return nil, fmt.Errorf("Volume '%s' 缂哄皯 HostPath.path", volumeName)
			}
			volume.VolumeSource = corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: path,
				},
			}
		case "configmap":
			configMapName, _ := config["name"].(string)
			configMapName = strings.TrimSpace(configMapName)
			if configMapName == "" {
				return nil, fmt.Errorf("Volume '%s' 缂哄皯 ConfigMap.name", volumeName)
			}
			volume.VolumeSource = corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: configMapName,
					},
				},
			}
		case "secret":
			secretName, _ := config["secretName"].(string)
			secretName = strings.TrimSpace(secretName)
			if secretName == "" {
				return nil, fmt.Errorf("Volume '%s' 缂哄皯 Secret.secretName", volumeName)
			}
			volume.VolumeSource = corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: secretName,
				},
			}
		case "persistentvolumeclaim", "pvc":
			claimName, _ := config["claimName"].(string)
			claimName = strings.TrimSpace(claimName)
			if claimName == "" {
				return nil, fmt.Errorf("Volume '%s' 缂哄皯 PersistentVolumeClaim.claimName", volumeName)
			}
			volume.VolumeSource = corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: claimName,
				},
			}
		default:
			return nil, fmt.Errorf("Volume '%s' 涓嶆敮鎸佺殑绫诲瀷: %s", volumeName, volumeSpec.Type)
		}

		volumes = append(volumes, volume)
	}

	return volumes, nil
}

func (s *K8sWorkloadServiceImpl) buildTolerationsForCreate(tolerationSpecs []model.Toleration) []corev1.Toleration {
	tolerations := make([]corev1.Toleration, 0, len(tolerationSpecs))
	for _, tolerationSpec := range tolerationSpecs {
		operator := corev1.TolerationOperator(strings.TrimSpace(tolerationSpec.Operator))
		if operator == "" {
			operator = corev1.TolerationOpEqual
		}

		tolerations = append(tolerations, corev1.Toleration{
			Key:      tolerationSpec.Key,
			Operator: operator,
			Value:    tolerationSpec.Value,
			Effect:   corev1.TaintEffect(strings.TrimSpace(tolerationSpec.Effect)),
		})
	}
	return tolerations
}

func (s *K8sWorkloadServiceImpl) buildDeploymentStrategy(strategy model.DeploymentStrategy) appsv1.DeploymentStrategy {
	strategyType := strings.TrimSpace(strategy.Type)
	if strategyType == "" {
		strategyType = string(appsv1.RollingUpdateDeploymentStrategyType)
	}

	deploymentStrategy := appsv1.DeploymentStrategy{
		Type: appsv1.DeploymentStrategyType(strategyType),
	}

	if deploymentStrategy.Type != appsv1.RollingUpdateDeploymentStrategyType {
		return deploymentStrategy
	}

	rollingUpdate := &appsv1.RollingUpdateDeployment{}
	if strings.TrimSpace(strategy.RollingUpdate.MaxUnavailable) != "" {
		maxUnavailable := intstr.Parse(strings.TrimSpace(strategy.RollingUpdate.MaxUnavailable))
		rollingUpdate.MaxUnavailable = &maxUnavailable
	}
	if strings.TrimSpace(strategy.RollingUpdate.MaxSurge) != "" {
		maxSurge := intstr.Parse(strings.TrimSpace(strategy.RollingUpdate.MaxSurge))
		rollingUpdate.MaxSurge = &maxSurge
	}
	deploymentStrategy.RollingUpdate = rollingUpdate

	return deploymentStrategy
}

func (s *K8sWorkloadServiceImpl) buildWorkloadSelectorLabels(defaultName string, workloadLabels, templateLabels map[string]string) map[string]string {
	if len(templateLabels) > 0 {
		return s.copyStringMap(templateLabels)
	}
	if len(workloadLabels) > 0 {
		return s.copyStringMap(workloadLabels)
	}
	return map[string]string{
		"app": defaultName,
	}
}

func (s *K8sWorkloadServiceImpl) mergeStringMaps(base map[string]string, overlays ...map[string]string) map[string]string {
	merged := s.copyStringMap(base)
	if merged == nil {
		merged = make(map[string]string)
	}
	for _, overlay := range overlays {
		for key, value := range overlay {
			merged[key] = value
		}
	}
	return merged
}

func (s *K8sWorkloadServiceImpl) copyStringMap(input map[string]string) map[string]string {
	if len(input) == 0 {
		return nil
	}

	output := make(map[string]string, len(input))
	for key, value := range input {
		output[key] = value
	}
	return output
}

// UpdateDeployment 更新Deployment
func (s *K8sWorkloadServiceImpl) UpdateDeployment(c *gin.Context, clusterId uint, namespaceName string, deploymentName string, req *model.UpdateWorkloadRequest) {
	// 获取集群信息
	cluster, err := s.clusterDao.GetByID(clusterId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result.Failed(c, http.StatusNotFound, "集群不存在")
			return
		}
		result.Failed(c, http.StatusInternalServerError, "获取集群信息失败: "+err.Error())
		return
	}

	// 检查集群状态
	if cluster.Status != model.ClusterStatusRunning {
		result.Failed(c, http.StatusBadRequest, "集群状态异常，无法操作Deployment")
		return
	}

	// 创建K8s客户端
	clientset, err := s.createK8sClient(cluster.Credential)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "连接K8s集群失败: "+err.Error())
		return
	}

	// 获取当前的Deployment
	deployment, err := clientset.AppsV1().Deployments(namespaceName).Get(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		result.Failed(c, http.StatusNotFound, "Deployment不存在: "+err.Error())
		return
	}

	// 更新Deployment的标签
	if req.Labels != nil && len(req.Labels) > 0 {
		if deployment.Labels == nil {
			deployment.Labels = make(map[string]string)
		}
		for k, v := range req.Labels {
			deployment.Labels[k] = v
		}
	}

	// 更新Pod标签
	if req.Template.Labels != nil && len(req.Template.Labels) > 0 {
		if deployment.Spec.Template.Labels == nil {
			deployment.Spec.Template.Labels = make(map[string]string)
		}
		for k, v := range req.Template.Labels {
			deployment.Spec.Template.Labels[k] = v
		}
	}

	// 更新容器配置
	if req.Template.Containers != nil && len(req.Template.Containers) > 0 {
		// 创建一个容器名称到现有容器的映射，用于保留现有配置
		existingContainers := make(map[string]corev1.Container)
		for _, container := range deployment.Spec.Template.Spec.Containers {
			existingContainers[container.Name] = container
		}

		// 更新或添加新容器
		updatedContainers := make([]corev1.Container, 0, len(req.Template.Containers))
		for _, containerSpec := range req.Template.Containers {
			// 从现有容器开始，保留原有配置
			var container corev1.Container
			if existing, exists := existingContainers[containerSpec.Name]; exists {
				container = existing
			} else {
				container = corev1.Container{
					Name: containerSpec.Name,
				}
			}

			// 只有当新的镜像不为空时才更新镜像
			if containerSpec.Image != "" {
				container.Image = containerSpec.Image
			}

			// 确保容器有镜像（这是必需的）
			if container.Image == "" {
				result.Failed(c, http.StatusBadRequest, fmt.Sprintf("容器 '%s' 缺少镜像配置", containerSpec.Name))
				return
			}

			// 设置端口
			if len(containerSpec.Ports) > 0 {
				container.Ports = make([]corev1.ContainerPort, 0, len(containerSpec.Ports))
				for _, portSpec := range containerSpec.Ports {
					container.Ports = append(container.Ports, corev1.ContainerPort{
						Name:          portSpec.Name,
						ContainerPort: portSpec.ContainerPort,
						Protocol:      corev1.Protocol(portSpec.Protocol),
					})
				}
			}

			// 设置环境变量
			if len(containerSpec.Env) > 0 {
				container.Env = make([]corev1.EnvVar, 0, len(containerSpec.Env))
				for _, envVar := range containerSpec.Env {
					container.Env = append(container.Env, corev1.EnvVar{
						Name:  envVar.Name,
						Value: envVar.Value,
					})
				}
			}

			// 设置资源限制
			container.Resources = corev1.ResourceRequirements{}
			if containerSpec.Resources.Limits.CPU != "" || containerSpec.Resources.Limits.Memory != "" {
				container.Resources.Limits = make(corev1.ResourceList)
				if containerSpec.Resources.Limits.CPU != "" {
					container.Resources.Limits[corev1.ResourceCPU] = resource.MustParse(containerSpec.Resources.Limits.CPU)
				}
				if containerSpec.Resources.Limits.Memory != "" {
					container.Resources.Limits[corev1.ResourceMemory] = resource.MustParse(containerSpec.Resources.Limits.Memory)
				}
			}
			if containerSpec.Resources.Requests.CPU != "" || containerSpec.Resources.Requests.Memory != "" {
				container.Resources.Requests = make(corev1.ResourceList)
				if containerSpec.Resources.Requests.CPU != "" {
					container.Resources.Requests[corev1.ResourceCPU] = resource.MustParse(containerSpec.Resources.Requests.CPU)
				}
				if containerSpec.Resources.Requests.Memory != "" {
					container.Resources.Requests[corev1.ResourceMemory] = resource.MustParse(containerSpec.Resources.Requests.Memory)
				}
			}

			// 设置存储卷挂载
			if len(containerSpec.VolumeMounts) > 0 {
				container.VolumeMounts = make([]corev1.VolumeMount, 0, len(containerSpec.VolumeMounts))
				for _, volumeMountSpec := range containerSpec.VolumeMounts {
					container.VolumeMounts = append(container.VolumeMounts, corev1.VolumeMount{
						Name:      volumeMountSpec.Name,
						MountPath: volumeMountSpec.MountPath,
						ReadOnly:  volumeMountSpec.ReadOnly,
					})
				}
			}

			updatedContainers = append(updatedContainers, container)
		}

		// 设置更新后的容器列表
		deployment.Spec.Template.Spec.Containers = updatedContainers

	}

	// 更新存储卷
	if req.Template.Volumes != nil && len(req.Template.Volumes) > 0 {
		deployment.Spec.Template.Spec.Volumes = make([]corev1.Volume, 0, len(req.Template.Volumes))
		for _, volumeSpec := range req.Template.Volumes {
			volume := corev1.Volume{
				Name: volumeSpec.Name,
			}

			// 根据类型设置存储卷源
			switch volumeSpec.Type {
			case "EmptyDir":
				volume.VolumeSource = corev1.VolumeSource{
					EmptyDir: &corev1.EmptyDirVolumeSource{},
				}
			case "HostPath":
				if config := volumeSpec.Config; config != nil {
					if path, ok := config["path"].(string); ok {
						volume.VolumeSource = corev1.VolumeSource{
							HostPath: &corev1.HostPathVolumeSource{
								Path: path,
							},
						}
					}
				}
			case "ConfigMap":
				if config := volumeSpec.Config; config != nil {
					if name, ok := config["name"].(string); ok {
						volume.VolumeSource = corev1.VolumeSource{
							ConfigMap: &corev1.ConfigMapVolumeSource{
								LocalObjectReference: corev1.LocalObjectReference{
									Name: name,
								},
							},
						}
					}
				}
			case "Secret":
				if config := volumeSpec.Config; config != nil {
					if secretName, ok := config["secretName"].(string); ok {
						volume.VolumeSource = corev1.VolumeSource{
							Secret: &corev1.SecretVolumeSource{
								SecretName: secretName,
							},
						}
					}
				}
			case "PersistentVolumeClaim":
				if config := volumeSpec.Config; config != nil {
					if claimName, ok := config["claimName"].(string); ok {
						volume.VolumeSource = corev1.VolumeSource{
							PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
								ClaimName: claimName,
							},
						}
					}
				}
			}

			deployment.Spec.Template.Spec.Volumes = append(deployment.Spec.Template.Spec.Volumes, volume)
		}
	}

	// 更新节点选择器
	if req.Template.NodeSelector != nil && len(req.Template.NodeSelector) > 0 {
		deployment.Spec.Template.Spec.NodeSelector = req.Template.NodeSelector
	}

	// 更新容忍度
	if req.Template.Tolerations != nil && len(req.Template.Tolerations) > 0 {
		deployment.Spec.Template.Spec.Tolerations = make([]corev1.Toleration, 0, len(req.Template.Tolerations))
		for _, tolerationSpec := range req.Template.Tolerations {
			deployment.Spec.Template.Spec.Tolerations = append(deployment.Spec.Template.Spec.Tolerations, corev1.Toleration{
				Key:      tolerationSpec.Key,
				Operator: corev1.TolerationOperator(tolerationSpec.Operator),
				Value:    tolerationSpec.Value,
				Effect:   corev1.TaintEffect(tolerationSpec.Effect),
			})
		}
	}

	// 更新部署策略
	if req.Strategy != nil {
		if strategyMap, ok := req.Strategy.(map[string]interface{}); ok {
			if strategyType, exists := strategyMap["type"]; exists {
				if strategyTypeStr, ok := strategyType.(string); ok {
					deployment.Spec.Strategy.Type = appsv1.DeploymentStrategyType(strategyTypeStr)

					// 如果是滚动更新策略，设置滚动更新参数
					if strategyTypeStr == "RollingUpdate" {
						if rollingUpdate, exists := strategyMap["rollingUpdate"]; exists {
							if rollingUpdateMap, ok := rollingUpdate.(map[string]interface{}); ok {
								deployment.Spec.Strategy.RollingUpdate = &appsv1.RollingUpdateDeployment{}

								if maxUnavailable, exists := rollingUpdateMap["maxUnavailable"]; exists {
									if maxUnavailableStr, ok := maxUnavailable.(string); ok {
										deployment.Spec.Strategy.RollingUpdate.MaxUnavailable = &intstr.IntOrString{
											Type:   intstr.String,
											StrVal: maxUnavailableStr,
										}
									}
								}

								if maxSurge, exists := rollingUpdateMap["maxSurge"]; exists {
									if maxSurgeStr, ok := maxSurge.(string); ok {
										deployment.Spec.Strategy.RollingUpdate.MaxSurge = &intstr.IntOrString{
											Type:   intstr.String,
											StrVal: maxSurgeStr,
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}

	// 执行更新操作
	updatedDeployment, err := clientset.AppsV1().Deployments(namespaceName).Update(context.TODO(), deployment, metav1.UpdateOptions{})
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "更新Deployment失败: "+err.Error())
		return
	}

	// 构造响应
	response := map[string]interface{}{
		"deploymentName":    deploymentName,
		"namespace":         namespaceName,
		"resourceVersion":   updatedDeployment.ResourceVersion,
		"generation":        updatedDeployment.Generation,
		"replicas":          updatedDeployment.Status.Replicas,
		"readyReplicas":     updatedDeployment.Status.ReadyReplicas,
		"availableReplicas": updatedDeployment.Status.AvailableReplicas,
		"updateTime":        time.Now().Format(time.RFC3339),
		"message":           fmt.Sprintf("Deployment '%s' 更新成功", deploymentName),
	}

	result.Success(c, response)
}

// DeleteDeployment 删除Deployment
func (s *K8sWorkloadServiceImpl) DeleteDeployment(c *gin.Context, clusterId uint, namespaceName string, deploymentName string) {
	// 获取集群信息
	cluster, err := s.clusterDao.GetByID(clusterId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result.Failed(c, http.StatusNotFound, "集群不存在")
			return
		}
		result.Failed(c, http.StatusInternalServerError, "获取集群信息失败: "+err.Error())
		return
	}

	// 检查集群状态
	if cluster.Status != model.ClusterStatusRunning {
		result.Failed(c, http.StatusBadRequest, "集群状态异常，无法操作Deployment")
		return
	}

	// 创建K8s客户端
	clientset, err := s.createK8sClient(cluster.Credential)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "连接K8s集群失败: "+err.Error())
		return
	}

	// 先获取当前的Deployment，确认存在并收集删除前的信息
	deployment, err := clientset.AppsV1().Deployments(namespaceName).Get(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		result.Failed(c, http.StatusNotFound, "Deployment不存在: "+err.Error())
		return
	}

	// 记录删除前的信息
	currentReplicas := int32(0)
	if deployment.Spec.Replicas != nil {
		currentReplicas = *deployment.Spec.Replicas
	}
	creationTime := deployment.CreationTimestamp.Time

	// 执行删除操作
	deletePolicy := metav1.DeletePropagationForeground // 前台删除，确保依赖资源也被删除
	deleteOptions := metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}

	err = clientset.AppsV1().Deployments(namespaceName).Delete(context.TODO(), deploymentName, deleteOptions)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "删除Deployment失败: "+err.Error())
		return
	}

	// 构造响应
	response := map[string]interface{}{
		"deploymentName": deploymentName,
		"namespace":      namespaceName,
		"clusterId":      clusterId,
		"replicas":       currentReplicas,
		"creationTime":   creationTime.Format(time.RFC3339),
		"deletionTime":   time.Now().Format(time.RFC3339),
		"message":        fmt.Sprintf("Deployment '%s' 删除成功", deploymentName),
		"note":           "相关的Pod和ReplicaSet将随之被删除",
	}

	result.Success(c, response)
}

// ScaleDeployment 伸缩Deployment
func (s *K8sWorkloadServiceImpl) ScaleDeployment(c *gin.Context, clusterId uint, namespaceName string, deploymentName string, req *model.ScaleWorkloadRequest) {
	// 获取集群信息
	cluster, err := s.clusterDao.GetByID(clusterId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result.Failed(c, http.StatusNotFound, "集群不存在")
			return
		}
		result.Failed(c, http.StatusInternalServerError, "获取集群信息失败: "+err.Error())
		return
	}

	// 检查集群状态
	if cluster.Status != model.ClusterStatusRunning {
		result.Failed(c, http.StatusBadRequest, "集群状态异常，无法操作Deployment")
		return
	}

	// 创建K8s客户端
	clientset, err := s.createK8sClient(cluster.Credential)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "连接K8s集群失败: "+err.Error())
		return
	}

	// 验证副本数
	if req.Replicas < 0 {
		result.Failed(c, http.StatusBadRequest, "副本数不能为负数")
		return
	}

	// 获取当前的Deployment
	deployment, err := clientset.AppsV1().Deployments(namespaceName).Get(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		result.Failed(c, http.StatusNotFound, "Deployment不存在: "+err.Error())
		return
	}

	// 记录当前副本数
	currentReplicas := int32(0)
	if deployment.Spec.Replicas != nil {
		currentReplicas = *deployment.Spec.Replicas
	}
	autoscalingSummary, _ := s.getSingleWorkloadAutoscalingSummary(clientset, namespaceName, "deployment", deploymentName)
	if autoscalingSummary != nil {
		if resources, resourcesErr := s.getWorkloadResourcesForAutoscaling(clientset, namespaceName, "deployment", deploymentName); resourcesErr == nil {
			s.enrichAutoscalingSummaryWarnings(autoscalingSummary, resources)
		}
	}

	// 更新副本数
	deployment.Spec.Replicas = &req.Replicas

	// 执行伸缩操作
	updatedDeployment, err := clientset.AppsV1().Deployments(namespaceName).Update(context.TODO(), deployment, metav1.UpdateOptions{})
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "伸缩Deployment失败: "+err.Error())
		return
	}

	// 构造响应
	response := map[string]interface{}{
		"deploymentName":     deploymentName,
		"namespace":          namespaceName,
		"currentReplicas":    currentReplicas,
		"targetReplicas":     req.Replicas,
		"readyReplicas":      updatedDeployment.Status.ReadyReplicas,
		"availableReplicas":  updatedDeployment.Status.AvailableReplicas,
		"autoscalingManaged": autoscalingSummary != nil,
		"autoscaling":        autoscalingSummary,
		"warnings": func() []string {
			if autoscalingSummary == nil {
				return nil
			}
			return autoscalingSummary.Warnings
		}(),
		"scaleTime": time.Now().Format(time.RFC3339),
		"message":   fmt.Sprintf("Deployment '%s' 已成功伸缩，副本数从 %d 调整为 %d", deploymentName, currentReplicas, req.Replicas),
	}
	if autoscalingSummary != nil {
		response["warning"] = fmt.Sprintf("当前工作负载仍受 HPA(%s) 管理，手动副本调整可能被自动扩缩容策略再次覆盖", autoscalingSummary.Name)
	}

	result.Success(c, response)
}

// RestartDeployment 重启Deployment
func (s *K8sWorkloadServiceImpl) RestartDeployment(c *gin.Context, clusterId uint, namespaceName string, deploymentName string) {
	// 获取集群信息
	cluster, err := s.clusterDao.GetByID(clusterId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result.Failed(c, http.StatusNotFound, "集群不存在")
			return
		}
		result.Failed(c, http.StatusInternalServerError, "获取集群信息失败: "+err.Error())
		return
	}

	// 检查集群状态
	if cluster.Status != model.ClusterStatusRunning {
		result.Failed(c, http.StatusBadRequest, "集群状态异常，无法操作Deployment")
		return
	}

	// 创建K8s客户端
	clientset, err := s.createK8sClient(cluster.Credential)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "连接K8s集群失败: "+err.Error())
		return
	}

	// 获取当前的Deployment
	deployment, err := clientset.AppsV1().Deployments(namespaceName).Get(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		result.Failed(c, http.StatusNotFound, "Deployment不存在: "+err.Error())
		return
	}

	// 记录重启前的信息
	currentReplicas := int32(0)
	if deployment.Spec.Replicas != nil {
		currentReplicas = *deployment.Spec.Replicas
	}

	// 触发滚动重启：通过更新Pod模板的annotation来强制重新创建Pod
	restartTime := time.Now().Format(time.RFC3339)
	if deployment.Spec.Template.Annotations == nil {
		deployment.Spec.Template.Annotations = make(map[string]string)
	}
	deployment.Spec.Template.Annotations["kubectl.kubernetes.io/restartedAt"] = restartTime

	// 更新Deployment来触发重启
	updatedDeployment, err := clientset.AppsV1().Deployments(namespaceName).Update(context.TODO(), deployment, metav1.UpdateOptions{})
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "重启Deployment失败: "+err.Error())
		return
	}

	// 构造响应
	response := map[string]interface{}{
		"deploymentName":    deploymentName,
		"namespace":         namespaceName,
		"replicas":          currentReplicas,
		"restartTime":       restartTime,
		"readyReplicas":     updatedDeployment.Status.ReadyReplicas,
		"availableReplicas": updatedDeployment.Status.AvailableReplicas,
		"message":           fmt.Sprintf("Deployment '%s' 重启已触发，Pod将进行滚动重启", deploymentName),
		"note":              "重启过程可能需要几分钟时间，请耐心等待Pod重新创建完成",
	}

	result.Success(c, response)
}

// CreateStatefulSet 创建StatefulSet
func (s *K8sWorkloadServiceImpl) CreateStatefulSet(c *gin.Context, clusterId uint, namespaceName string, req *model.CreateStatefulSetRequest) {
	cluster, err := s.clusterDao.GetByID(clusterId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result.Failed(c, http.StatusNotFound, "闆嗙兢涓嶅瓨鍦?")
			return
		}
		result.Failed(c, http.StatusInternalServerError, "鑾峰彇闆嗙兢淇℃伅澶辫触: "+err.Error())
		return
	}

	if cluster.Status != model.ClusterStatusRunning {
		result.Failed(c, http.StatusBadRequest, "闆嗙兢鐘舵€佸紓甯革紝鏃犳硶鍒涘缓StatefulSet")
		return
	}

	clientset, err := s.createK8sClient(cluster.Credential)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "杩炴帴K8s闆嗙兢澶辫触: "+err.Error())
		return
	}

	statefulSet, err := s.buildStatefulSetForCreate(namespaceName, req)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "鏋勫缓StatefulSet澶辫触: "+err.Error())
		return
	}

	createdStatefulSet, err := clientset.AppsV1().StatefulSets(namespaceName).Create(context.TODO(), statefulSet, metav1.CreateOptions{})
	if err != nil {
		switch {
		case apierrors.IsAlreadyExists(err):
			result.Failed(c, http.StatusConflict, "StatefulSet宸插瓨鍦? "+err.Error())
		case apierrors.IsInvalid(err), apierrors.IsBadRequest(err):
			result.Failed(c, http.StatusBadRequest, "StatefulSet閰嶇疆鏃犳晥: "+err.Error())
		default:
			result.Failed(c, http.StatusInternalServerError, "鍒涘缓StatefulSet澶辫触: "+err.Error())
		}
		return
	}

	result.Success(c, s.convertStatefulSetToWorkload(createdStatefulSet))
}

// UpdateStatefulSet 更新StatefulSet
func (s *K8sWorkloadServiceImpl) UpdateStatefulSet(c *gin.Context, clusterId uint, namespaceName string, statefulSetName string, req *model.UpdateWorkloadRequest) {
	cluster, err := s.clusterDao.GetByID(clusterId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result.Failed(c, http.StatusNotFound, "闂嗗棛鍏㈡稉宥呯摠閸?")
			return
		}
		result.Failed(c, http.StatusInternalServerError, "閼惧嘲褰囬梿鍡欏參娣団剝浼呮径杈Е: "+err.Error())
		return
	}

	if cluster.Status != model.ClusterStatusRunning {
		result.Failed(c, http.StatusBadRequest, "闂嗗棛鍏㈤悩鑸碘偓浣哥磽鐢潻绱濋弮鐘崇《閹碘晝缂塖tatefulSet")
		return
	}

	clientset, err := s.createK8sClient(cluster.Credential)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "鏉╃偞甯碖8s闂嗗棛鍏㈡径杈Е: "+err.Error())
		return
	}

	statefulSet, err := clientset.AppsV1().StatefulSets(namespaceName).Get(context.TODO(), statefulSetName, metav1.GetOptions{})
	if err != nil {
		result.Failed(c, http.StatusNotFound, "StatefulSet娑撳秴鐡ㄩ崷? "+err.Error())
		return
	}

	if err := s.applyWorkloadUpdateRequest(&statefulSet.Labels, &statefulSet.Spec.Template, req); err != nil {
		result.Failed(c, http.StatusBadRequest, "StatefulSet鏇存柊鍙傛暟閿欒: "+err.Error())
		return
	}

	updatedStatefulSet, err := clientset.AppsV1().StatefulSets(namespaceName).Update(context.TODO(), statefulSet, metav1.UpdateOptions{})
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "鏇存柊StatefulSet澶辫触: "+err.Error())
		return
	}

	result.Success(c, map[string]interface{}{
		"statefulSetName": statefulSetName,
		"namespace":       namespaceName,
		"resourceVersion": updatedStatefulSet.ResourceVersion,
		"generation":      updatedStatefulSet.Generation,
		"replicas":        updatedStatefulSet.Status.Replicas,
		"readyReplicas":   updatedStatefulSet.Status.ReadyReplicas,
		"updateTime":      time.Now().Format(time.RFC3339),
		"message":         fmt.Sprintf("StatefulSet '%s' 鏇存柊鎴愬姛", statefulSetName),
	})
	return
	result.Failed(c, http.StatusNotImplemented, "StatefulSet更新功能暂未实现，敬请期待")
}

// DeleteStatefulSet 删除StatefulSet
func (s *K8sWorkloadServiceImpl) DeleteStatefulSet(c *gin.Context, clusterId uint, namespaceName string, statefulSetName string) {
	cluster, err := s.clusterDao.GetByID(clusterId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result.Failed(c, http.StatusNotFound, "闆嗙兢涓嶅瓨鍦?")
			return
		}
		result.Failed(c, http.StatusInternalServerError, "鑾峰彇闆嗙兢淇℃伅澶辫触: "+err.Error())
		return
	}

	if cluster.Status != model.ClusterStatusRunning {
		result.Failed(c, http.StatusBadRequest, "闆嗙兢鐘舵€佸紓甯革紝鏃犳硶鍒犻櫎StatefulSet")
		return
	}

	clientset, err := s.createK8sClient(cluster.Credential)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "杩炴帴K8s闆嗙兢澶辫触: "+err.Error())
		return
	}

	statefulSet, err := clientset.AppsV1().StatefulSets(namespaceName).Get(context.TODO(), statefulSetName, metav1.GetOptions{})
	if err != nil {
		result.Failed(c, http.StatusNotFound, "StatefulSet涓嶅瓨鍦? "+err.Error())
		return
	}

	currentReplicas := int32(0)
	if statefulSet.Spec.Replicas != nil {
		currentReplicas = *statefulSet.Spec.Replicas
	}

	deletePolicy := metav1.DeletePropagationForeground
	deleteOptions := metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}
	if err := clientset.AppsV1().StatefulSets(namespaceName).Delete(context.TODO(), statefulSetName, deleteOptions); err != nil {
		result.Failed(c, http.StatusInternalServerError, "鍒犻櫎StatefulSet澶辫触: "+err.Error())
		return
	}

	result.Success(c, map[string]interface{}{
		"statefulSetName": statefulSetName,
		"namespace":       namespaceName,
		"clusterId":       clusterId,
		"replicas":        currentReplicas,
		"serviceName":     statefulSet.Spec.ServiceName,
		"deletionTime":    time.Now().Format(time.RFC3339),
		"message":         fmt.Sprintf("StatefulSet '%s' 删除成功", statefulSetName),
	})
}

// ScaleStatefulSet 伸缩StatefulSet
func (s *K8sWorkloadServiceImpl) ScaleStatefulSet(c *gin.Context, clusterId uint, namespaceName string, statefulSetName string, req *model.ScaleWorkloadRequest) {
	cluster, err := s.clusterDao.GetByID(clusterId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result.Failed(c, http.StatusNotFound, "闆嗙兢涓嶅瓨鍦?")
			return
		}
		result.Failed(c, http.StatusInternalServerError, "鑾峰彇闆嗙兢淇℃伅澶辫触: "+err.Error())
		return
	}

	if cluster.Status != model.ClusterStatusRunning {
		result.Failed(c, http.StatusBadRequest, "闆嗙兢鐘舵€佸紓甯革紝鏃犳硶鎵╃缉StatefulSet")
		return
	}

	clientset, err := s.createK8sClient(cluster.Credential)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "杩炴帴K8s闆嗙兢澶辫触: "+err.Error())
		return
	}

	if req.Replicas < 0 {
		result.Failed(c, http.StatusBadRequest, "鍓湰鏁颁笉鑳戒负璐熸暟")
		return
	}

	statefulSet, err := clientset.AppsV1().StatefulSets(namespaceName).Get(context.TODO(), statefulSetName, metav1.GetOptions{})
	if err != nil {
		result.Failed(c, http.StatusNotFound, "StatefulSet涓嶅瓨鍦? "+err.Error())
		return
	}

	currentReplicas := int32(0)
	if statefulSet.Spec.Replicas != nil {
		currentReplicas = *statefulSet.Spec.Replicas
	}
	autoscalingSummary, _ := s.getSingleWorkloadAutoscalingSummary(clientset, namespaceName, "statefulset", statefulSetName)
	if autoscalingSummary != nil {
		if resources, resourcesErr := s.getWorkloadResourcesForAutoscaling(clientset, namespaceName, "statefulset", statefulSetName); resourcesErr == nil {
			s.enrichAutoscalingSummaryWarnings(autoscalingSummary, resources)
		}
	}

	statefulSet.Spec.Replicas = &req.Replicas
	updatedStatefulSet, err := clientset.AppsV1().StatefulSets(namespaceName).Update(context.TODO(), statefulSet, metav1.UpdateOptions{})
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "鎵╃缉StatefulSet澶辫触: "+err.Error())
		return
	}

	result.Success(c, map[string]interface{}{
		"statefulSetName":    statefulSetName,
		"namespace":          namespaceName,
		"currentReplicas":    currentReplicas,
		"targetReplicas":     req.Replicas,
		"readyReplicas":      updatedStatefulSet.Status.ReadyReplicas,
		"autoscalingManaged": autoscalingSummary != nil,
		"autoscaling":        autoscalingSummary,
		"warnings": func() []string {
			if autoscalingSummary == nil {
				return nil
			}
			return autoscalingSummary.Warnings
		}(),
		"scaleTime": time.Now().Format(time.RFC3339),
		"message":   fmt.Sprintf("StatefulSet '%s' 已成功伸缩，副本数从 %d 调整为 %d", statefulSetName, currentReplicas, req.Replicas),
		"warning": func() string {
			if autoscalingSummary == nil {
				return ""
			}
			return fmt.Sprintf("当前工作负载仍受 HPA(%s) 管理，手动副本调整可能被自动扩缩容策略再次覆盖", autoscalingSummary.Name)
		}(),
	})
}

// RestartStatefulSet 重启StatefulSet
func (s *K8sWorkloadServiceImpl) RestartStatefulSet(c *gin.Context, clusterId uint, namespaceName string, statefulSetName string) {
	cluster, err := s.clusterDao.GetByID(clusterId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result.Failed(c, http.StatusNotFound, "闆嗙兢涓嶅瓨鍦?")
			return
		}
		result.Failed(c, http.StatusInternalServerError, "鑾峰彇闆嗙兢淇℃伅澶辫触: "+err.Error())
		return
	}

	if cluster.Status != model.ClusterStatusRunning {
		result.Failed(c, http.StatusBadRequest, "闆嗙兢鐘舵€佸紓甯革紝鏃犳硶閲嶅惎StatefulSet")
		return
	}

	clientset, err := s.createK8sClient(cluster.Credential)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "杩炴帴K8s闆嗙兢澶辫触: "+err.Error())
		return
	}

	statefulSet, err := clientset.AppsV1().StatefulSets(namespaceName).Get(context.TODO(), statefulSetName, metav1.GetOptions{})
	if err != nil {
		result.Failed(c, http.StatusNotFound, "StatefulSet涓嶅瓨鍦? "+err.Error())
		return
	}

	restartTime := time.Now().Format(time.RFC3339)
	if statefulSet.Spec.Template.Annotations == nil {
		statefulSet.Spec.Template.Annotations = make(map[string]string)
	}
	statefulSet.Spec.Template.Annotations["kubectl.kubernetes.io/restartedAt"] = restartTime

	updatedStatefulSet, err := clientset.AppsV1().StatefulSets(namespaceName).Update(context.TODO(), statefulSet, metav1.UpdateOptions{})
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "閲嶅惎StatefulSet澶辫触: "+err.Error())
		return
	}

	result.Success(c, map[string]interface{}{
		"statefulSetName": statefulSetName,
		"namespace":       namespaceName,
		"restartTime":     restartTime,
		"readyReplicas":   updatedStatefulSet.Status.ReadyReplicas,
		"message":         fmt.Sprintf("StatefulSet '%s' 重启已触发，Pod 将进行滚动重启", statefulSetName),
	})
}

// CreateDaemonSet 创建DaemonSet
func (s *K8sWorkloadServiceImpl) CreateDaemonSet(c *gin.Context, clusterId uint, namespaceName string, req *model.CreateDaemonSetRequest) {
	cluster, err := s.clusterDao.GetByID(clusterId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result.Failed(c, http.StatusNotFound, "闆嗙兢涓嶅瓨鍦?")
			return
		}
		result.Failed(c, http.StatusInternalServerError, "鑾峰彇闆嗙兢淇℃伅澶辫触: "+err.Error())
		return
	}

	if cluster.Status != model.ClusterStatusRunning {
		result.Failed(c, http.StatusBadRequest, "闆嗙兢鐘舵€佸紓甯革紝鏃犳硶鍒涘缓DaemonSet")
		return
	}

	clientset, err := s.createK8sClient(cluster.Credential)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "杩炴帴K8s闆嗙兢澶辫触: "+err.Error())
		return
	}

	daemonSet, err := s.buildDaemonSetForCreate(namespaceName, req)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "鏋勫缓DaemonSet澶辫触: "+err.Error())
		return
	}

	createdDaemonSet, err := clientset.AppsV1().DaemonSets(namespaceName).Create(context.TODO(), daemonSet, metav1.CreateOptions{})
	if err != nil {
		switch {
		case apierrors.IsAlreadyExists(err):
			result.Failed(c, http.StatusConflict, "DaemonSet宸插瓨鍦? "+err.Error())
		case apierrors.IsInvalid(err), apierrors.IsBadRequest(err):
			result.Failed(c, http.StatusBadRequest, "DaemonSet閰嶇疆鏃犳晥: "+err.Error())
		default:
			result.Failed(c, http.StatusInternalServerError, "鍒涘缓DaemonSet澶辫触: "+err.Error())
		}
		return
	}

	result.Success(c, s.convertDaemonSetToWorkload(createdDaemonSet))
}

// UpdateDaemonSet 更新DaemonSet
func (s *K8sWorkloadServiceImpl) UpdateDaemonSet(c *gin.Context, clusterId uint, namespaceName string, daemonSetName string, req *model.UpdateWorkloadRequest) {
	cluster, err := s.clusterDao.GetByID(clusterId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result.Failed(c, http.StatusNotFound, "闂嗗棛鍏㈡稉宥呯摠閸?")
			return
		}
		result.Failed(c, http.StatusInternalServerError, "閼惧嘲褰囬梿鍡欏參娣団剝浼呮径杈Е: "+err.Error())
		return
	}

	if cluster.Status != model.ClusterStatusRunning {
		result.Failed(c, http.StatusBadRequest, "闂嗗棛鍏㈤悩鑸碘偓浣哥磽鐢潻绱濋弮鐘崇《閸掓稑缂揇aemonSet")
		return
	}

	clientset, err := s.createK8sClient(cluster.Credential)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "鏉╃偞甯碖8s闂嗗棛鍏㈡径杈Е: "+err.Error())
		return
	}

	daemonSet, err := clientset.AppsV1().DaemonSets(namespaceName).Get(context.TODO(), daemonSetName, metav1.GetOptions{})
	if err != nil {
		result.Failed(c, http.StatusNotFound, "DaemonSet娑撳秴鐡ㄩ崷? "+err.Error())
		return
	}

	if err := s.applyWorkloadUpdateRequest(&daemonSet.Labels, &daemonSet.Spec.Template, req); err != nil {
		result.Failed(c, http.StatusBadRequest, "DaemonSet鏇存柊鍙傛暟閿欒: "+err.Error())
		return
	}

	updatedDaemonSet, err := clientset.AppsV1().DaemonSets(namespaceName).Update(context.TODO(), daemonSet, metav1.UpdateOptions{})
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "鏇存柊DaemonSet澶辫触: "+err.Error())
		return
	}

	result.Success(c, map[string]interface{}{
		"daemonSetName":   daemonSetName,
		"namespace":       namespaceName,
		"resourceVersion": updatedDaemonSet.ResourceVersion,
		"generation":      updatedDaemonSet.Generation,
		"desired":         updatedDaemonSet.Status.DesiredNumberScheduled,
		"ready":           updatedDaemonSet.Status.NumberReady,
		"updateTime":      time.Now().Format(time.RFC3339),
		"message":         fmt.Sprintf("DaemonSet '%s' 鏇存柊鎴愬姛", daemonSetName),
	})
	return
	result.Failed(c, http.StatusNotImplemented, "DaemonSet更新功能暂未实现，敬请期待")
}

// DeleteDaemonSet 删除DaemonSet
func (s *K8sWorkloadServiceImpl) DeleteDaemonSet(c *gin.Context, clusterId uint, namespaceName string, daemonSetName string) {
	cluster, err := s.clusterDao.GetByID(clusterId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result.Failed(c, http.StatusNotFound, "闆嗙兢涓嶅瓨鍦?")
			return
		}
		result.Failed(c, http.StatusInternalServerError, "鑾峰彇闆嗙兢淇℃伅澶辫触: "+err.Error())
		return
	}

	if cluster.Status != model.ClusterStatusRunning {
		result.Failed(c, http.StatusBadRequest, "闆嗙兢鐘舵€佸紓甯革紝鏃犳硶鍒犻櫎DaemonSet")
		return
	}

	clientset, err := s.createK8sClient(cluster.Credential)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "杩炴帴K8s闆嗙兢澶辫触: "+err.Error())
		return
	}

	daemonSet, err := clientset.AppsV1().DaemonSets(namespaceName).Get(context.TODO(), daemonSetName, metav1.GetOptions{})
	if err != nil {
		result.Failed(c, http.StatusNotFound, "DaemonSet涓嶅瓨鍦? "+err.Error())
		return
	}

	deletePolicy := metav1.DeletePropagationForeground
	deleteOptions := metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}
	if err := clientset.AppsV1().DaemonSets(namespaceName).Delete(context.TODO(), daemonSetName, deleteOptions); err != nil {
		result.Failed(c, http.StatusInternalServerError, "鍒犻櫎DaemonSet澶辫触: "+err.Error())
		return
	}

	result.Success(c, map[string]interface{}{
		"daemonSetName": daemonSetName,
		"namespace":     namespaceName,
		"clusterId":     clusterId,
		"selector":      daemonSet.Spec.Selector.MatchLabels,
		"deletionTime":  time.Now().Format(time.RFC3339),
		"message":       fmt.Sprintf("DaemonSet '%s' 删除成功", daemonSetName),
	})
}

// RestartDaemonSet 重启DaemonSet
func (s *K8sWorkloadServiceImpl) RestartDaemonSet(c *gin.Context, clusterId uint, namespaceName string, daemonSetName string) {
	cluster, err := s.clusterDao.GetByID(clusterId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result.Failed(c, http.StatusNotFound, "闆嗙兢涓嶅瓨鍦?")
			return
		}
		result.Failed(c, http.StatusInternalServerError, "鑾峰彇闆嗙兢淇℃伅澶辫触: "+err.Error())
		return
	}

	if cluster.Status != model.ClusterStatusRunning {
		result.Failed(c, http.StatusBadRequest, "闆嗙兢鐘舵€佸紓甯革紝鏃犳硶閲嶅惎DaemonSet")
		return
	}

	clientset, err := s.createK8sClient(cluster.Credential)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "杩炴帴K8s闆嗙兢澶辫触: "+err.Error())
		return
	}

	daemonSet, err := clientset.AppsV1().DaemonSets(namespaceName).Get(context.TODO(), daemonSetName, metav1.GetOptions{})
	if err != nil {
		result.Failed(c, http.StatusNotFound, "DaemonSet涓嶅瓨鍦? "+err.Error())
		return
	}

	restartTime := time.Now().Format(time.RFC3339)
	if daemonSet.Spec.Template.Annotations == nil {
		daemonSet.Spec.Template.Annotations = make(map[string]string)
	}
	daemonSet.Spec.Template.Annotations["kubectl.kubernetes.io/restartedAt"] = restartTime

	updatedDaemonSet, err := clientset.AppsV1().DaemonSets(namespaceName).Update(context.TODO(), daemonSet, metav1.UpdateOptions{})
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "閲嶅惎DaemonSet澶辫触: "+err.Error())
		return
	}

	result.Success(c, map[string]interface{}{
		"daemonSetName": daemonSetName,
		"namespace":     namespaceName,
		"restartTime":   restartTime,
		"readyReplicas": updatedDaemonSet.Status.NumberReady,
		"message":       fmt.Sprintf("DaemonSet '%s' 重启已触发，Pod 将进行滚动重启", daemonSetName),
	})
}

// CreateJob 创建Job
func (s *K8sWorkloadServiceImpl) CreateJob(c *gin.Context, clusterId uint, namespaceName string, req *model.CreateJobRequest) {
	cluster, err := s.clusterDao.GetByID(clusterId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result.Failed(c, http.StatusNotFound, "闆嗙兢涓嶅瓨鍦?")
			return
		}
		result.Failed(c, http.StatusInternalServerError, "鑾峰彇闆嗙兢淇℃伅澶辫触: "+err.Error())
		return
	}

	if cluster.Status != model.ClusterStatusRunning {
		result.Failed(c, http.StatusBadRequest, "闆嗙兢鐘舵€佸紓甯革紝鏃犳硶鍒涘缓Job")
		return
	}

	clientset, err := s.createK8sClient(cluster.Credential)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "杩炴帴K8s闆嗙兢澶辫触: "+err.Error())
		return
	}

	job, err := s.buildJobForCreate(namespaceName, req)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "鏋勫缓Job澶辫触: "+err.Error())
		return
	}

	createdJob, err := clientset.BatchV1().Jobs(namespaceName).Create(context.TODO(), job, metav1.CreateOptions{})
	if err != nil {
		switch {
		case apierrors.IsAlreadyExists(err):
			result.Failed(c, http.StatusConflict, "Job宸插瓨鍦? "+err.Error())
		case apierrors.IsInvalid(err), apierrors.IsBadRequest(err):
			result.Failed(c, http.StatusBadRequest, "Job閰嶇疆鏃犳晥: "+err.Error())
		default:
			result.Failed(c, http.StatusInternalServerError, "鍒涘缓Job澶辫触: "+err.Error())
		}
		return
	}

	result.Success(c, s.convertJobToWorkload(createdJob))
}

// DeleteJob 删除Job
func (s *K8sWorkloadServiceImpl) UpdateJob(c *gin.Context, clusterId uint, namespaceName string, jobName string, req *model.UpdateWorkloadRequest) {
	cluster, err := s.clusterDao.GetByID(clusterId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result.Failed(c, http.StatusNotFound, "闂嗗棛鍏㈡稉宥呯摠閸?")
			return
		}
		result.Failed(c, http.StatusInternalServerError, "閼惧嘲褰囬梿鍡欏參娣團剝浼呮径杈Е: "+err.Error())
		return
	}

	if cluster.Status != model.ClusterStatusRunning {
		result.Failed(c, http.StatusBadRequest, "闂嗗棛鍏㈤悩鑸碘偓浣哥磽鐢潻绱濋弮鐘崇《閸掓稑缂揓ob")
		return
	}

	clientset, err := s.createK8sClient(cluster.Credential)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "鏉╃偞甯碖8s闂嗗棛鍏㈡径杈Е: "+err.Error())
		return
	}

	job, err := clientset.BatchV1().Jobs(namespaceName).Get(context.TODO(), jobName, metav1.GetOptions{})
	if err != nil {
		result.Failed(c, http.StatusNotFound, "Job娑撳秴鐡ㄩ崷? "+err.Error())
		return
	}

	if err := s.applyWorkloadUpdateRequest(&job.Labels, &job.Spec.Template, req); err != nil {
		result.Failed(c, http.StatusBadRequest, "Job鏇存柊鍙傛暟閿欒: "+err.Error())
		return
	}

	updatedJob, err := clientset.BatchV1().Jobs(namespaceName).Update(context.TODO(), job, metav1.UpdateOptions{})
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "鏇存柊Job澶辫触: "+err.Error())
		return
	}

	result.Success(c, map[string]interface{}{
		"jobName":         jobName,
		"namespace":       namespaceName,
		"resourceVersion": updatedJob.ResourceVersion,
		"generation":      updatedJob.Generation,
		"active":          updatedJob.Status.Active,
		"succeeded":       updatedJob.Status.Succeeded,
		"failed":          updatedJob.Status.Failed,
		"updateTime":      time.Now().Format(time.RFC3339),
		"message":         fmt.Sprintf("Job '%s' 鏇存柊鎴愬姛", jobName),
	})
}

func (s *K8sWorkloadServiceImpl) DeleteJob(c *gin.Context, clusterId uint, namespaceName string, jobName string) {
	cluster, err := s.clusterDao.GetByID(clusterId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result.Failed(c, http.StatusNotFound, "闆嗙兢涓嶅瓨鍦?")
			return
		}
		result.Failed(c, http.StatusInternalServerError, "鑾峰彇闆嗙兢淇℃伅澶辫触: "+err.Error())
		return
	}

	if cluster.Status != model.ClusterStatusRunning {
		result.Failed(c, http.StatusBadRequest, "闆嗙兢鐘舵€佸紓甯革紝鏃犳硶鍒犻櫎Job")
		return
	}

	clientset, err := s.createK8sClient(cluster.Credential)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "杩炴帴K8s闆嗙兢澶辫触: "+err.Error())
		return
	}

	job, err := clientset.BatchV1().Jobs(namespaceName).Get(context.TODO(), jobName, metav1.GetOptions{})
	if err != nil {
		result.Failed(c, http.StatusNotFound, "Job涓嶅瓨鍦? "+err.Error())
		return
	}

	deletePolicy := metav1.DeletePropagationForeground
	deleteOptions := metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}
	if err := clientset.BatchV1().Jobs(namespaceName).Delete(context.TODO(), jobName, deleteOptions); err != nil {
		result.Failed(c, http.StatusInternalServerError, "鍒犻櫎Job澶辫触: "+err.Error())
		return
	}

	result.Success(c, map[string]interface{}{
		"jobName":      jobName,
		"namespace":    namespaceName,
		"clusterId":    clusterId,
		"selector":     job.Spec.Selector,
		"deletionTime": time.Now().Format(time.RFC3339),
		"message":      fmt.Sprintf("Job '%s' 删除成功", jobName),
	})
}

// CreateCronJob 创建CronJob
func (s *K8sWorkloadServiceImpl) CreateCronJob(c *gin.Context, clusterId uint, namespaceName string, req *model.CreateCronJobRequest) {
	cluster, err := s.clusterDao.GetByID(clusterId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result.Failed(c, http.StatusNotFound, "闆嗙兢涓嶅瓨鍦?")
			return
		}
		result.Failed(c, http.StatusInternalServerError, "鑾峰彇闆嗙兢淇℃伅澶辫触: "+err.Error())
		return
	}

	if cluster.Status != model.ClusterStatusRunning {
		result.Failed(c, http.StatusBadRequest, "闆嗙兢鐘舵€佸紓甯革紝鏃犳硶鍒涘缓CronJob")
		return
	}

	clientset, err := s.createK8sClient(cluster.Credential)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "杩炴帴K8s闆嗙兢澶辫触: "+err.Error())
		return
	}

	cronJob, err := s.buildCronJobForCreate(namespaceName, req)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "鏋勫缓CronJob澶辫触: "+err.Error())
		return
	}

	createdCronJob, err := clientset.BatchV1().CronJobs(namespaceName).Create(context.TODO(), cronJob, metav1.CreateOptions{})
	if err != nil {
		switch {
		case apierrors.IsAlreadyExists(err):
			result.Failed(c, http.StatusConflict, "CronJob宸插瓨鍦? "+err.Error())
		case apierrors.IsInvalid(err), apierrors.IsBadRequest(err):
			result.Failed(c, http.StatusBadRequest, "CronJob閰嶇疆鏃犳晥: "+err.Error())
		default:
			result.Failed(c, http.StatusInternalServerError, "鍒涘缓CronJob澶辫触: "+err.Error())
		}
		return
	}

	result.Success(c, s.convertCronJobToWorkload(createdCronJob))
}

// UpdateCronJob 更新CronJob
func (s *K8sWorkloadServiceImpl) UpdateCronJob(c *gin.Context, clusterId uint, namespaceName string, cronJobName string, req *model.UpdateWorkloadRequest) {
	cluster, err := s.clusterDao.GetByID(clusterId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result.Failed(c, http.StatusNotFound, "闂嗗棛鍏㈡稉宥呯摠閸?")
			return
		}
		result.Failed(c, http.StatusInternalServerError, "閼惧嘲褰囬梿鍡欏參娣團剝浼呮径杈Е: "+err.Error())
		return
	}

	if cluster.Status != model.ClusterStatusRunning {
		result.Failed(c, http.StatusBadRequest, "闂嗗棛鍏㈤悩鑸碘偓浣哥磽鐢潻绱濋弮鐘崇《閸掓稑缂揅ronJob")
		return
	}

	clientset, err := s.createK8sClient(cluster.Credential)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "鏉╃偞甯碖8s闂嗗棛鍏㈡径杈Е: "+err.Error())
		return
	}

	cronJob, err := clientset.BatchV1().CronJobs(namespaceName).Get(context.TODO(), cronJobName, metav1.GetOptions{})
	if err != nil {
		result.Failed(c, http.StatusNotFound, "CronJob娑撳秴鐡ㄩ崷? "+err.Error())
		return
	}

	if req.Labels != nil && len(req.Labels) > 0 {
		if cronJob.Labels == nil {
			cronJob.Labels = make(map[string]string)
		}
		for k, v := range req.Labels {
			cronJob.Labels[k] = v
		}
	}

	if err := s.applyWorkloadUpdateRequest(&cronJob.Spec.JobTemplate.Labels, &cronJob.Spec.JobTemplate.Spec.Template, req); err != nil {
		result.Failed(c, http.StatusBadRequest, "CronJob鏇存柊鍙傛暟閿欒: "+err.Error())
		return
	}

	updatedCronJob, err := clientset.BatchV1().CronJobs(namespaceName).Update(context.TODO(), cronJob, metav1.UpdateOptions{})
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "鏇存柊CronJob澶辫触: "+err.Error())
		return
	}

	result.Success(c, map[string]interface{}{
		"cronJobName":     cronJobName,
		"namespace":       namespaceName,
		"resourceVersion": updatedCronJob.ResourceVersion,
		"generation":      updatedCronJob.Generation,
		"schedule":        updatedCronJob.Spec.Schedule,
		"updateTime":      time.Now().Format(time.RFC3339),
		"message":         fmt.Sprintf("CronJob '%s' 鏇存柊鎴愬姛", cronJobName),
	})
	return
	result.Failed(c, http.StatusNotImplemented, "CronJob更新功能暂未实现，敬请期待")
}

// DeleteCronJob 删除CronJob
func (s *K8sWorkloadServiceImpl) DeleteCronJob(c *gin.Context, clusterId uint, namespaceName string, cronJobName string) {
	cluster, err := s.clusterDao.GetByID(clusterId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result.Failed(c, http.StatusNotFound, "闆嗙兢涓嶅瓨鍦?")
			return
		}
		result.Failed(c, http.StatusInternalServerError, "鑾峰彇闆嗙兢淇℃伅澶辫触: "+err.Error())
		return
	}

	if cluster.Status != model.ClusterStatusRunning {
		result.Failed(c, http.StatusBadRequest, "闆嗙兢鐘舵€佸紓甯革紝鏃犳硶鍒犻櫎CronJob")
		return
	}

	clientset, err := s.createK8sClient(cluster.Credential)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "杩炴帴K8s闆嗙兢澶辫触: "+err.Error())
		return
	}

	cronJob, err := clientset.BatchV1().CronJobs(namespaceName).Get(context.TODO(), cronJobName, metav1.GetOptions{})
	if err != nil {
		result.Failed(c, http.StatusNotFound, "CronJob涓嶅瓨鍦? "+err.Error())
		return
	}

	deletePolicy := metav1.DeletePropagationForeground
	deleteOptions := metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}
	if err := clientset.BatchV1().CronJobs(namespaceName).Delete(context.TODO(), cronJobName, deleteOptions); err != nil {
		result.Failed(c, http.StatusInternalServerError, "鍒犻櫎CronJob澶辫触: "+err.Error())
		return
	}

	result.Success(c, map[string]interface{}{
		"cronJobName":  cronJobName,
		"namespace":    namespaceName,
		"clusterId":    clusterId,
		"schedule":     cronJob.Spec.Schedule,
		"deletionTime": time.Now().Format(time.RFC3339),
		"message":      fmt.Sprintf("CronJob '%s' 删除成功", cronJobName),
	})
}

// ===================== 监控API实现 =====================

// GetPodMetrics 获取Pod监控指标
func (s *K8sWorkloadServiceImpl) GetPodMetrics(c *gin.Context, clusterId uint, namespaceName string, podName string) {
	// 获取集群信息
	cluster, err := s.clusterDao.GetByID(clusterId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result.Failed(c, http.StatusNotFound, "集群不存在")
			return
		}
		result.Failed(c, http.StatusInternalServerError, "获取集群信息失败: "+err.Error())
		return
	}

	// 创建K8s客户端
	clientset, err := s.createK8sClient(cluster.Credential)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "连接K8s集群失败: "+err.Error())
		return
	}

	// 创建Metrics客户端
	metricsClient, err := s.createMetricsClient(cluster.Credential)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "连接Metrics服务失败: "+err.Error())
		return
	}

	// 获取Pod基础信息
	pod, err := clientset.CoreV1().Pods(namespaceName).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		result.Failed(c, http.StatusNotFound, "Pod不存在: "+err.Error())
		return
	}

	// 获取Pod Metrics
	podMetrics, err := metricsClient.MetricsV1beta1().PodMetricses(namespaceName).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "获取Pod监控数据失败: "+err.Error())
		return
	}

	// 转换为响应结构
	podMetricsInfo := s.convertToPodMetricsInfo(pod, podMetrics)

	result.Success(c, podMetricsInfo)
}

// GetNodeMetrics 获取节点监控指标
func (s *K8sWorkloadServiceImpl) GetNodeMetrics(c *gin.Context, clusterId uint, nodeName string) {
	// 获取集群信息
	cluster, err := s.clusterDao.GetByID(clusterId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result.Failed(c, http.StatusNotFound, "集群不存在")
			return
		}
		result.Failed(c, http.StatusInternalServerError, "获取集群信息失败: "+err.Error())
		return
	}

	// 创建K8s客户端
	clientset, err := s.createK8sClient(cluster.Credential)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "连接K8s集群失败: "+err.Error())
		return
	}

	// 创建Metrics客户端
	metricsClient, err := s.createMetricsClient(cluster.Credential)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "连接Metrics服务失败: "+err.Error())
		return
	}

	// 获取节点基础信息
	node, err := clientset.CoreV1().Nodes().Get(context.TODO(), nodeName, metav1.GetOptions{})
	if err != nil {
		result.Failed(c, http.StatusNotFound, "节点不存在: "+err.Error())
		return
	}

	// 获取Node Metrics
	nodeMetrics, err := metricsClient.MetricsV1beta1().NodeMetricses().Get(context.TODO(), nodeName, metav1.GetOptions{})
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "获取节点监控数据失败: "+err.Error())
		return
	}

	// 获取节点上的所有Pod Metrics
	podMetricsList, err := metricsClient.MetricsV1beta1().PodMetricses("").List(context.TODO(), metav1.ListOptions{
		FieldSelector: fmt.Sprintf("spec.nodeName=%s", nodeName),
	})
	if err != nil {
		podMetricsList = &metricsv1beta1.PodMetricsList{} // 如果获取失败，使用空列表
	}

	// 转换为响应结构
	nodeMetricsInfo := s.convertToNodeMetricsInfo(node, nodeMetrics, podMetricsList, clientset)

	result.Success(c, nodeMetricsInfo)
}

// GetNamespaceMetrics 获取命名空间监控指标
func (s *K8sWorkloadServiceImpl) GetNamespaceMetrics(c *gin.Context, clusterId uint, namespaceName string) {
	// 获取集群信息
	cluster, err := s.clusterDao.GetByID(clusterId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result.Failed(c, http.StatusNotFound, "集群不存在")
			return
		}
		result.Failed(c, http.StatusInternalServerError, "获取集群信息失败: "+err.Error())
		return
	}

	// 创建K8s客户端
	clientset, err := s.createK8sClient(cluster.Credential)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "连接K8s集群失败: "+err.Error())
		return
	}

	// 创建Metrics客户端
	metricsClient, err := s.createMetricsClient(cluster.Credential)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "连接Metrics服务失败: "+err.Error())
		return
	}

	// 获取命名空间基础信息
	namespace, err := clientset.CoreV1().Namespaces().Get(context.TODO(), namespaceName, metav1.GetOptions{})
	if err != nil {
		result.Failed(c, http.StatusNotFound, "命名空间不存在: "+err.Error())
		return
	}

	// 获取命名空间下所有Pod的Metrics
	podMetricsList, err := metricsClient.MetricsV1beta1().PodMetricses(namespaceName).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "获取命名空间监控数据失败: "+err.Error())
		return
	}

	// 获取命名空间下所有Pod的基础信息
	podList, err := clientset.CoreV1().Pods(namespaceName).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "获取Pod列表失败: "+err.Error())
		return
	}

	// 转换为响应结构
	namespaceMetricsInfo := s.convertToNamespaceMetricsInfo(namespace, podMetricsList, podList, clientset)

	result.Success(c, namespaceMetricsInfo)
}

// ===================== 监控数据转换辅助方法 =====================

// convertToPodMetricsInfo 转换Pod监控信息
func (s *K8sWorkloadServiceImpl) convertToPodMetricsInfo(pod *corev1.Pod, podMetrics *metricsv1beta1.PodMetrics) model.PodMetricsInfo {
	// 计算总使用量
	var totalCPU, totalMemory resource.Quantity
	var containerMetrics []model.ContainerMetricsInfo

	// 遍历容器metrics
	for _, containerMetric := range podMetrics.Containers {
		cpuUsage := containerMetric.Usage[corev1.ResourceCPU]
		memoryUsage := containerMetric.Usage[corev1.ResourceMemory]

		totalCPU.Add(cpuUsage)
		totalMemory.Add(memoryUsage)

		// 获取对应的容器配置信息
		var containerRequests, containerLimits model.ResourceUsage
		var containerState string
		var restartCount int32

		// 从Pod规格中找到匹配的容器
		for _, container := range pod.Spec.Containers {
			if container.Name == containerMetric.Name {
				// 获取资源请求和限制
				if container.Resources.Requests != nil {
					if cpu, ok := container.Resources.Requests[corev1.ResourceCPU]; ok {
						containerRequests.CPU = cpu.String()
					}
					if memory, ok := container.Resources.Requests[corev1.ResourceMemory]; ok {
						containerRequests.Memory = memory.String()
					}
				}

				if container.Resources.Limits != nil {
					if cpu, ok := container.Resources.Limits[corev1.ResourceCPU]; ok {
						containerLimits.CPU = cpu.String()
					}
					if memory, ok := container.Resources.Limits[corev1.ResourceMemory]; ok {
						containerLimits.Memory = memory.String()
					}
				}
				break
			}
		}

		// 从Pod状态中获取容器状态和重启次数
		for _, containerStatus := range pod.Status.ContainerStatuses {
			if containerStatus.Name == containerMetric.Name {
				containerState = s.getContainerState(containerStatus.State)
				restartCount = containerStatus.RestartCount
				break
			}
		}

		// 计算使用率
		var cpuRate, memoryRate float64
		if containerRequests.CPU != "" {
			if requestCPU, err := resource.ParseQuantity(containerRequests.CPU); err == nil {
				if requestCPU.MilliValue() > 0 {
					cpuRate = float64(cpuUsage.MilliValue()) / float64(requestCPU.MilliValue()) * 100
				}
			}
		}
		if containerRequests.Memory != "" {
			if requestMemory, err := resource.ParseQuantity(containerRequests.Memory); err == nil {
				if requestMemory.Value() > 0 {
					memoryRate = float64(memoryUsage.Value()) / float64(requestMemory.Value()) * 100
				}
			}
		}

		containerMetrics = append(containerMetrics, model.ContainerMetricsInfo{
			Name: containerMetric.Name,
			Usage: model.ResourceUsage{
				CPU:    cpuUsage.String(),
				Memory: memoryUsage.String(),
			},
			Requests: containerRequests,
			Limits:   containerLimits,
			UsageRate: model.ResourceUsageRate{
				CPURate:    cpuRate,
				MemoryRate: memoryRate,
			},
			State:        containerState,
			RestartCount: restartCount,
		})
	}

	// 计算Pod级别的资源配额和使用率
	podRequests := s.extractPodResources(&pod.Spec)
	var podCPURate, podMemoryRate float64

	if podRequests.Requests.CPU != "" {
		if requestCPU, err := resource.ParseQuantity(podRequests.Requests.CPU); err == nil {
			if requestCPU.MilliValue() > 0 {
				podCPURate = float64(totalCPU.MilliValue()) / float64(requestCPU.MilliValue()) * 100
			}
		}
	}
	if podRequests.Requests.Memory != "" {
		if requestMemory, err := resource.ParseQuantity(podRequests.Requests.Memory); err == nil {
			if requestMemory.Value() > 0 {
				podMemoryRate = float64(totalMemory.Value()) / float64(requestMemory.Value()) * 100
			}
		}
	}

	return model.PodMetricsInfo{
		PodName:    pod.Name,
		Namespace:  pod.Namespace,
		NodeName:   pod.Spec.NodeName,
		Timestamp:  podMetrics.Timestamp.Format(time.RFC3339),
		Containers: containerMetrics,
		TotalUsage: model.ResourceUsage{
			CPU:    totalCPU.String(),
			Memory: totalMemory.String(),
		},
		ResourceQuota: model.PodResourceQuota{
			Requests: model.ResourceUsage{
				CPU:    podRequests.Requests.CPU,
				Memory: podRequests.Requests.Memory,
			},
			Limits: model.ResourceUsage{
				CPU:    podRequests.Limits.CPU,
				Memory: podRequests.Limits.Memory,
			},
		},
		UsageRate: model.ResourceUsageRate{
			CPURate:    podCPURate,
			MemoryRate: podMemoryRate,
		},
	}
}

// convertToNodeMetricsInfo 转换Node监控信息
func (s *K8sWorkloadServiceImpl) convertToNodeMetricsInfo(node *corev1.Node, nodeMetrics *metricsv1beta1.NodeMetrics, podMetricsList *metricsv1beta1.PodMetricsList, clientset *kubernetes.Clientset) model.NodeMetricsInfo {
	// 获取节点资源信息，安全地获取资源值
	var capacityCPU, capacityMemory string
	if cpu, ok := node.Status.Capacity[corev1.ResourceCPU]; ok {
		capacityCPU = cpu.String()
	}
	if memory, ok := node.Status.Capacity[corev1.ResourceMemory]; ok {
		capacityMemory = memory.String()
	}

	var allocatableCPU, allocatableMemory string
	if cpu, ok := node.Status.Allocatable[corev1.ResourceCPU]; ok {
		allocatableCPU = cpu.String()
	}
	if memory, ok := node.Status.Allocatable[corev1.ResourceMemory]; ok {
		allocatableMemory = memory.String()
	}

	var usageCPU, usageMemory string
	if cpu, ok := nodeMetrics.Usage[corev1.ResourceCPU]; ok {
		usageCPU = cpu.String()
	}
	if memory, ok := nodeMetrics.Usage[corev1.ResourceMemory]; ok {
		usageMemory = memory.String()
	}

	capacity := model.ResourceUsage{
		CPU:    capacityCPU,
		Memory: capacityMemory,
	}

	allocatable := model.ResourceUsage{
		CPU:    allocatableCPU,
		Memory: allocatableMemory,
	}

	usage := model.ResourceUsage{
		CPU:    usageCPU,
		Memory: usageMemory,
	}

	// 计算使用率
	var cpuRate, memoryRate float64
	if allocatableCPU, err := resource.ParseQuantity(allocatable.CPU); err == nil {
		if usedCPU, err := resource.ParseQuantity(usage.CPU); err == nil {
			if allocatableCPU.MilliValue() > 0 {
				cpuRate = float64(usedCPU.MilliValue()) / float64(allocatableCPU.MilliValue()) * 100
			}
		}
	}
	if allocatableMemory, err := resource.ParseQuantity(allocatable.Memory); err == nil {
		if usedMemory, err := resource.ParseQuantity(usage.Memory); err == nil {
			if allocatableMemory.Value() > 0 {
				memoryRate = float64(usedMemory.Value()) / float64(allocatableMemory.Value()) * 100
			}
		}
	}

	// 转换Pod监控摘要
	var podMetricsSummary []model.PodMetricsSummary
	for _, podMetric := range podMetricsList.Items {
		var totalCPU, totalMemory resource.Quantity
		for _, containerMetric := range podMetric.Containers {
			totalCPU.Add(containerMetric.Usage[corev1.ResourceCPU])
			totalMemory.Add(containerMetric.Usage[corev1.ResourceMemory])
		}

		podMetricsSummary = append(podMetricsSummary, model.PodMetricsSummary{
			PodName:   podMetric.Name,
			Namespace: podMetric.Namespace,
			Usage: model.ResourceUsage{
				CPU:    totalCPU.String(),
				Memory: totalMemory.String(),
			},
			UsageRate: model.ResourceUsageRate{
				CPURate:    0, // 此处可以计算Pod级别的使用率，但需要额外的Pod配置信息
				MemoryRate: 0,
			},
		})
	}

	return model.NodeMetricsInfo{
		NodeName:    node.Name,
		Timestamp:   nodeMetrics.Timestamp.Format(time.RFC3339),
		Usage:       usage,
		Capacity:    capacity,
		Allocatable: allocatable,
		UsageRate: model.ResourceUsageRate{
			CPURate:    cpuRate,
			MemoryRate: memoryRate,
		},
		PodCount:   len(podMetricsList.Items),
		PodMetrics: podMetricsSummary,
		SystemInfo: model.NodeSystemInfo{
			KernelVersion:           node.Status.NodeInfo.KernelVersion,
			OSImage:                 node.Status.NodeInfo.OSImage,
			ContainerRuntimeVersion: node.Status.NodeInfo.ContainerRuntimeVersion,
			KubeletVersion:          node.Status.NodeInfo.KubeletVersion,
			KubeProxyVersion:        node.Status.NodeInfo.KubeProxyVersion,
			Architecture:            node.Status.NodeInfo.Architecture,
		},
	}
}

// convertToNamespaceMetricsInfo 转换Namespace监控信息
func (s *K8sWorkloadServiceImpl) convertToNamespaceMetricsInfo(namespace *corev1.Namespace, podMetricsList *metricsv1beta1.PodMetricsList, podList *corev1.PodList, clientset *kubernetes.Clientset) model.NamespaceMetricsInfo {
	// 计算总使用量
	var totalCPU, totalMemory resource.Quantity
	var podMetricsSummary []model.PodMetricsSummary
	var runningPods int

	for _, podMetric := range podMetricsList.Items {
		var podCPU, podMemory resource.Quantity
		for _, containerMetric := range podMetric.Containers {
			podCPU.Add(containerMetric.Usage[corev1.ResourceCPU])
			podMemory.Add(containerMetric.Usage[corev1.ResourceMemory])
		}

		totalCPU.Add(podCPU)
		totalMemory.Add(podMemory)

		podMetricsSummary = append(podMetricsSummary, model.PodMetricsSummary{
			PodName:   podMetric.Name,
			Namespace: podMetric.Namespace,
			Usage: model.ResourceUsage{
				CPU:    podCPU.String(),
				Memory: podMemory.String(),
			},
			UsageRate: model.ResourceUsageRate{
				CPURate:    0, // 可以根据需要计算Pod使用率
				MemoryRate: 0,
			},
		})
	}

	// 统计运行中的Pod数量
	for _, pod := range podList.Items {
		if pod.Status.Phase == corev1.PodRunning {
			runningPods++
		}
	}

	// 尝试获取ResourceQuota（如果存在）
	var resourceQuota model.NamespaceResourceQuota
	if quotaList, err := clientset.CoreV1().ResourceQuotas(namespace.Name).List(context.TODO(), metav1.ListOptions{}); err == nil && len(quotaList.Items) > 0 {
		// 使用第一个ResourceQuota
		quota := quotaList.Items[0]
		if hardCPU, ok := quota.Status.Hard[corev1.ResourceCPU]; ok {
			resourceQuota.Hard.CPU = hardCPU.String()
		}
		if hardMemory, ok := quota.Status.Hard[corev1.ResourceMemory]; ok {
			resourceQuota.Hard.Memory = hardMemory.String()
		}
		if usedCPU, ok := quota.Status.Used[corev1.ResourceCPU]; ok {
			resourceQuota.Used.CPU = usedCPU.String()
		}
		if usedMemory, ok := quota.Status.Used[corev1.ResourceMemory]; ok {
			resourceQuota.Used.Memory = usedMemory.String()
		}
	}

	// 计算使用率
	var cpuRate, memoryRate float64
	if resourceQuota.Hard.CPU != "" && resourceQuota.Used.CPU != "" {
		if hardCPU, err := resource.ParseQuantity(resourceQuota.Hard.CPU); err == nil {
			if usedCPU, err := resource.ParseQuantity(resourceQuota.Used.CPU); err == nil {
				if hardCPU.MilliValue() > 0 {
					cpuRate = float64(usedCPU.MilliValue()) / float64(hardCPU.MilliValue()) * 100
				}
			}
		}
	}
	if resourceQuota.Hard.Memory != "" && resourceQuota.Used.Memory != "" {
		if hardMemory, err := resource.ParseQuantity(resourceQuota.Hard.Memory); err == nil {
			if usedMemory, err := resource.ParseQuantity(resourceQuota.Used.Memory); err == nil {
				if hardMemory.Value() > 0 {
					memoryRate = float64(usedMemory.Value()) / float64(hardMemory.Value()) * 100
				}
			}
		}
	}

	return model.NamespaceMetricsInfo{
		Namespace:   namespace.Name,
		Timestamp:   time.Now().Format(time.RFC3339),
		PodCount:    len(podList.Items),
		RunningPods: runningPods,
		TotalUsage: model.ResourceUsage{
			CPU:    totalCPU.String(),
			Memory: totalMemory.String(),
		},
		ResourceQuota: resourceQuota,
		UsageRate: model.ResourceUsageRate{
			CPURate:    cpuRate,
			MemoryRate: memoryRate,
		},
		PodMetrics: podMetricsSummary,
	}
}

// ===================== YAML管理API实现 =====================

// CreatePodFromYAML 通过YAML创建Pod
func (s *K8sWorkloadServiceImpl) CreatePodFromYAML(c *gin.Context, clusterId uint, namespaceName string, req *model.CreatePodFromYAMLRequest) {
	// 获取集群信息
	cluster, err := s.clusterDao.GetByID(clusterId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result.Failed(c, http.StatusNotFound, "集群不存在")
			return
		}
		result.Failed(c, http.StatusInternalServerError, "获取集群信息失败: "+err.Error())
		return
	}

	// 首先校验YAML格式
	validationResult, err := s.validateYAMLContent(req.YAMLContent, "pod")
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "YAML校验失败: "+err.Error())
		return
	}

	// 如果只是校验格式，直接返回校验结果
	if req.ValidateOnly {
		result.Success(c, model.CreatePodFromYAMLResponse{
			Success:          validationResult.Valid,
			Message:          "YAML校验完成",
			ValidationResult: validationResult,
		})
		return
	}

	// 解析YAML为Kubernetes对象
	obj, err := s.parseYAMLToKubernetesObject(req.YAMLContent)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "解析YAML失败: "+err.Error())
		return
	}

	// 创建K8s客户端
	clientset, err := s.createK8sClient(cluster.Credential)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "连接K8s集群失败: "+err.Error())
		return
	}

	// 如果是DryRun，只进行校验不实际创建
	createOptions := metav1.CreateOptions{}
	if req.DryRun {
		createOptions.DryRun = []string{metav1.DryRunAll}
	}

	var resourceName, resourceKind string

	// 根据资源类型进行相应的创建操作
	switch resource := obj.(type) {
	case *corev1.Pod:
		// 验证Pod对象是否有必要的字段
		if resource.Name == "" && resource.GenerateName == "" {
			result.Failed(c, http.StatusBadRequest, "Pod必须指定name或generateName")
			return
		}

		// 验证Pod是否有容器定义
		if len(resource.Spec.Containers) == 0 {
			result.Failed(c, http.StatusBadRequest, "Pod必须至少包含一个容器")
			return
		}

		// 如果YAML中没有指定命名空间，使用URL参数中的命名空间
		if resource.Namespace == "" {
			resource.Namespace = namespaceName
		} else if resource.Namespace != namespaceName {
			result.Failed(c, http.StatusBadRequest, fmt.Sprintf("YAML中指定的命名空间(%s)与URL参数不匹配(%s)", resource.Namespace, namespaceName))
			return
		}

		// 创建Pod
		created, err := clientset.CoreV1().Pods(namespaceName).Create(context.TODO(), resource, createOptions)
		if err != nil {
			result.Failed(c, http.StatusInternalServerError, "创建Pod失败: "+err.Error())
			return
		}
		resourceName = created.Name
		resourceKind = "Pod"

	case *appsv1.Deployment:
		// 验证Deployment对象是否有必要的字段
		if resource.Name == "" && resource.GenerateName == "" {
			result.Failed(c, http.StatusBadRequest, "Deployment必须指定name或generateName")
			return
		}

		// 如果YAML中没有指定命名空间，使用URL参数中的命名空间
		if resource.Namespace == "" {
			resource.Namespace = namespaceName
		} else if resource.Namespace != namespaceName {
			result.Failed(c, http.StatusBadRequest, fmt.Sprintf("YAML中指定的命名空间(%s)与URL参数不匹配(%s)", resource.Namespace, namespaceName))
			return
		}

		// 创建Deployment
		created, err := clientset.AppsV1().Deployments(namespaceName).Create(context.TODO(), resource, createOptions)
		if err != nil {
			result.Failed(c, http.StatusInternalServerError, "创建Deployment失败: "+err.Error())
			return
		}
		resourceName = created.Name
		resourceKind = "Deployment"

	case *corev1.Service:
		// 如果YAML中没有指定命名空间，使用URL参数中的命名空间
		if resource.Namespace == "" {
			resource.Namespace = namespaceName
		} else if resource.Namespace != namespaceName {
			result.Failed(c, http.StatusBadRequest, fmt.Sprintf("YAML中指定的命名空间(%s)与URL参数不匹配(%s)", resource.Namespace, namespaceName))
			return
		}

		// 创建Service
		created, err := clientset.CoreV1().Services(namespaceName).Create(context.TODO(), resource, createOptions)
		if err != nil {
			result.Failed(c, http.StatusInternalServerError, "创建Service失败: "+err.Error())
			return
		}
		resourceName = created.Name
		resourceKind = "Service"

	case *corev1.ConfigMap:
		// 如果YAML中没有指定命名空间，使用URL参数中的命名空间
		if resource.Namespace == "" {
			resource.Namespace = namespaceName
		} else if resource.Namespace != namespaceName {
			result.Failed(c, http.StatusBadRequest, fmt.Sprintf("YAML中指定的命名空间(%s)与URL参数不匹配(%s)", resource.Namespace, namespaceName))
			return
		}

		// 创建ConfigMap
		created, err := clientset.CoreV1().ConfigMaps(namespaceName).Create(context.TODO(), resource, createOptions)
		if err != nil {
			result.Failed(c, http.StatusInternalServerError, "创建ConfigMap失败: "+err.Error())
			return
		}
		resourceName = created.Name
		resourceKind = "ConfigMap"

	case *networkingv1.Ingress:
		// 验证Ingress对象是否有必要的字段
		if resource.Name == "" && resource.GenerateName == "" {
			result.Failed(c, http.StatusBadRequest, "Ingress必须指定name或generateName")
			return
		}

		// 如果YAML中没有指定命名空间，使用URL参数中的命名空间
		if resource.Namespace == "" {
			resource.Namespace = namespaceName
		} else if resource.Namespace != namespaceName {
			result.Failed(c, http.StatusBadRequest, fmt.Sprintf("YAML中指定的命名空间(%s)与URL参数不匹配(%s)", resource.Namespace, namespaceName))
			return
		}

		// 创建Ingress
		created, err := clientset.NetworkingV1().Ingresses(namespaceName).Create(context.TODO(), resource, createOptions)
		if err != nil {
			result.Failed(c, http.StatusInternalServerError, "创建Ingress失败: "+err.Error())
			return
		}
		resourceName = created.Name
		resourceKind = "Ingress"

	case *corev1.PersistentVolume:
		// 验证PersistentVolume对象是否有必要的字段
		if resource.Name == "" && resource.GenerateName == "" {
			result.Failed(c, http.StatusBadRequest, "PersistentVolume必须指定name或generateName")
			return
		}

		// 验证PV必须有容量定义
		if resource.Spec.Capacity == nil || len(resource.Spec.Capacity) == 0 {
			result.Failed(c, http.StatusBadRequest, "PersistentVolume必须指定容量(capacity)")
			return
		}

		// 验证PV必须有访问模式
		if len(resource.Spec.AccessModes) == 0 {
			result.Failed(c, http.StatusBadRequest, "PersistentVolume必须指定访问模式(accessModes)")
			return
		}

		// 注意：PersistentVolume是集群级别资源，不属于任何命名空间
		// 如果YAML中指定了命名空间，应该给出警告
		if resource.Namespace != "" {
			result.Failed(c, http.StatusBadRequest, "PersistentVolume是集群级别资源，不应指定命名空间")
			return
		}

		// 创建PersistentVolume（注意PV是集群级别资源，不需要命名空间参数）
		created, err := clientset.CoreV1().PersistentVolumes().Create(context.TODO(), resource, createOptions)
		if err != nil {
			result.Failed(c, http.StatusInternalServerError, "创建PersistentVolume失败: "+err.Error())
			return
		}
		resourceName = created.Name
		resourceKind = "PersistentVolume"

	case *corev1.PersistentVolumeClaim:
		// 验证PersistentVolumeClaim对象是否有必要的字段
		if resource.Name == "" && resource.GenerateName == "" {
			result.Failed(c, http.StatusBadRequest, "PersistentVolumeClaim必须指定name或generateName")
			return
		}

		// 验证PVC必须有资源请求
		if resource.Spec.Resources.Requests == nil || len(resource.Spec.Resources.Requests) == 0 {
			result.Failed(c, http.StatusBadRequest, "PersistentVolumeClaim必须指定资源请求(resources.requests)")
			return
		}

		// 验证PVC必须有访问模式
		if len(resource.Spec.AccessModes) == 0 {
			result.Failed(c, http.StatusBadRequest, "PersistentVolumeClaim必须指定访问模式(accessModes)")
			return
		}

		// 如果YAML中没有指定命名空间，使用URL参数中的命名空间
		if resource.Namespace == "" {
			resource.Namespace = namespaceName
		} else if resource.Namespace != namespaceName {
			result.Failed(c, http.StatusBadRequest, fmt.Sprintf("YAML中指定的命名空间(%s)与URL参数不匹配(%s)", resource.Namespace, namespaceName))
			return
		}

		// 创建PersistentVolumeClaim
		created, err := clientset.CoreV1().PersistentVolumeClaims(namespaceName).Create(context.TODO(), resource, createOptions)
		if err != nil {
			result.Failed(c, http.StatusInternalServerError, "创建PersistentVolumeClaim失败: "+err.Error())
			return
		}
		resourceName = created.Name
		resourceKind = "PersistentVolumeClaim"

	default:
		result.Failed(c, http.StatusBadRequest, fmt.Sprintf("不支持的资源类型: %T", resource))
		return
	}

	// 构造响应
	response := model.CreatePodFromYAMLResponse{
		Success:   true,
		PodName:   resourceName, // 这里沿用原来的字段名，但实际上是资源名称
		Namespace: namespaceName,
		ParsedObject: map[string]interface{}{
			"kind":      resourceKind,
			"name":      resourceName,
			"namespace": namespaceName,
		},
	}

	if req.DryRun {
		response.Message = fmt.Sprintf("%s配置校验成功，未实际创建", resourceKind)
		response.ValidationResult = validationResult
	} else {
		response.Message = fmt.Sprintf("%s '%s' 创建成功", resourceKind, resourceName)
	}

	result.Success(c, response)
}

// ValidateYAML 校验YAML格式
func (s *K8sWorkloadServiceImpl) ValidateYAML(c *gin.Context, clusterId uint, req *model.ValidateYAMLRequest) {
	// 校验YAML内容
	validationResult, err := s.validateYAMLContent(req.YAMLContent, req.ResourceType)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "校验失败: "+err.Error())
		return
	}

	result.Success(c, validationResult)
}

// ===================== YAML处理辅助方法 =====================

// validateYAMLContent 校验YAML内容
func (s *K8sWorkloadServiceImpl) validateYAMLContent(yamlContent, resourceType string) (*model.ValidateYAMLResponse, error) {
	response := &model.ValidateYAMLResponse{
		Valid:       true,
		Errors:      []string{},
		Warnings:    []string{},
		Suggestions: []string{},
	}

	// 1. 基础YAML语法校验
	var parsedObject map[string]interface{}
	if err := yaml.Unmarshal([]byte(yamlContent), &parsedObject); err != nil {
		response.Valid = false
		response.Errors = append(response.Errors, fmt.Sprintf("YAML语法错误: %v", err))
		return response, nil
	}

	response.ParsedObject = parsedObject

	// 2. 检查必要的Kubernetes字段
	if apiVersion, ok := parsedObject["apiVersion"].(string); !ok || apiVersion == "" {
		response.Valid = false
		response.Errors = append(response.Errors, "缺少必需的字段: apiVersion")
	}

	if kind, ok := parsedObject["kind"].(string); !ok || kind == "" {
		response.Valid = false
		response.Errors = append(response.Errors, "缺少必需的字段: kind")
	} else {
		// 3. 如果指定了资源类型，检查是否匹配
		if resourceType != "" && strings.ToLower(kind) != strings.ToLower(resourceType) {
			response.Warnings = append(response.Warnings, fmt.Sprintf("资源类型不匹配，期望: %s, 实际: %s", resourceType, kind))
		}
	}

	if metadata, ok := parsedObject["metadata"].(map[string]interface{}); !ok {
		response.Valid = false
		response.Errors = append(response.Errors, "缺少必需的字段: metadata")
	} else {
		if name, ok := metadata["name"].(string); !ok || name == "" {
			response.Valid = false
			response.Errors = append(response.Errors, "缺少必需的字段: metadata.name")
		}
	}

	// 4. 针对Pod的特殊校验
	if strings.ToLower(resourceType) == "pod" || (parsedObject["kind"] != nil && strings.ToLower(parsedObject["kind"].(string)) == "pod") {
		if err := s.validatePodSpec(parsedObject, response); err != nil {
			return response, err
		}
	}

	// 5. 提供改进建议
	s.addYAMLSuggestions(parsedObject, response)

	return response, nil
}

// validatePodSpec 校验Pod规格
func (s *K8sWorkloadServiceImpl) validatePodSpec(parsedObject map[string]interface{}, response *model.ValidateYAMLResponse) error {
	spec, ok := parsedObject["spec"].(map[string]interface{})
	if !ok {
		response.Valid = false
		response.Errors = append(response.Errors, "缺少必需的字段: spec")
		return nil
	}

	// 检查containers字段
	containers, ok := spec["containers"].([]interface{})
	if !ok || len(containers) == 0 {
		response.Valid = false
		response.Errors = append(response.Errors, "缺少必需的字段: spec.containers")
		return nil
	}

	// 校验每个容器
	for i, container := range containers {
		containerMap, ok := container.(map[string]interface{})
		if !ok {
			response.Valid = false
			response.Errors = append(response.Errors, fmt.Sprintf("容器[%d]配置格式错误", i))
			continue
		}

		// 检查容器名称
		if name, ok := containerMap["name"].(string); !ok || name == "" {
			response.Valid = false
			response.Errors = append(response.Errors, fmt.Sprintf("容器[%d]缺少名称", i))
		}

		// 检查容器镜像
		if image, ok := containerMap["image"].(string); !ok || image == "" {
			response.Valid = false
			response.Errors = append(response.Errors, fmt.Sprintf("容器[%d]缺少镜像", i))
		} else {
			// 检查镜像标签
			if !strings.Contains(image, ":") {
				response.Warnings = append(response.Warnings, fmt.Sprintf("容器[%d]镜像未指定标签，将使用latest", i))
				response.Suggestions = append(response.Suggestions, fmt.Sprintf("建议为容器[%d]指定明确的镜像标签，如: %s:v1.0", i, image))
			}
		}

		// 检查资源配置
		if resources, ok := containerMap["resources"].(map[string]interface{}); ok {
			s.validateResourcesConfig(resources, i, response)
		} else {
			response.Warnings = append(response.Warnings, fmt.Sprintf("容器[%d]未配置资源限制", i))
			response.Suggestions = append(response.Suggestions, fmt.Sprintf("建议为容器[%d]配置CPU和内存的requests和limits", i))
		}
	}

	return nil
}

// validateResourcesConfig 校验资源配置
func (s *K8sWorkloadServiceImpl) validateResourcesConfig(resources map[string]interface{}, containerIndex int, response *model.ValidateYAMLResponse) {
	// 检查requests配置
	if requests, ok := resources["requests"].(map[string]interface{}); ok {
		if cpu, exists := requests["cpu"]; exists {
			if cpuStr, ok := cpu.(string); ok {
				if _, err := resource.ParseQuantity(cpuStr); err != nil {
					response.Errors = append(response.Errors, fmt.Sprintf("容器[%d] CPU requests格式错误: %s", containerIndex, cpuStr))
				}
			}
		}
		if memory, exists := requests["memory"]; exists {
			if memStr, ok := memory.(string); ok {
				if _, err := resource.ParseQuantity(memStr); err != nil {
					response.Errors = append(response.Errors, fmt.Sprintf("容器[%d] Memory requests格式错误: %s", containerIndex, memStr))
				}
			}
		}
	}

	// 检查limits配置
	if limits, ok := resources["limits"].(map[string]interface{}); ok {
		if cpu, exists := limits["cpu"]; exists {
			if cpuStr, ok := cpu.(string); ok {
				if _, err := resource.ParseQuantity(cpuStr); err != nil {
					response.Errors = append(response.Errors, fmt.Sprintf("容器[%d] CPU limits格式错误: %s", containerIndex, cpuStr))
				}
			}
		}
		if memory, exists := limits["memory"]; exists {
			if memStr, ok := memory.(string); ok {
				if _, err := resource.ParseQuantity(memStr); err != nil {
					response.Errors = append(response.Errors, fmt.Sprintf("容器[%d] Memory limits格式错误: %s", containerIndex, memStr))
				}
			}
		}
	}
}

// addYAMLSuggestions 添加YAML改进建议
func (s *K8sWorkloadServiceImpl) addYAMLSuggestions(parsedObject map[string]interface{}, response *model.ValidateYAMLResponse) {
	// 检查是否有labels
	if metadata, ok := parsedObject["metadata"].(map[string]interface{}); ok {
		if _, hasLabels := metadata["labels"]; !hasLabels {
			response.Suggestions = append(response.Suggestions, "建议添加标签(labels)以便更好地管理和选择Pod")
		}
	}

	// 检查是否有健康检查
	if spec, ok := parsedObject["spec"].(map[string]interface{}); ok {
		if containers, ok := spec["containers"].([]interface{}); ok {
			for i, container := range containers {
				if containerMap, ok := container.(map[string]interface{}); ok {
					if _, hasLiveness := containerMap["livenessProbe"]; !hasLiveness {
						response.Suggestions = append(response.Suggestions, fmt.Sprintf("建议为容器[%d]配置存活性探针(livenessProbe)", i))
					}
					if _, hasReadiness := containerMap["readinessProbe"]; !hasReadiness {
						response.Suggestions = append(response.Suggestions, fmt.Sprintf("建议为容器[%d]配置就绪性探针(readinessProbe)", i))
					}
				}
			}
		}
	}

	// 检查安全上下文
	if spec, ok := parsedObject["spec"].(map[string]interface{}); ok {
		if _, hasSecurityContext := spec["securityContext"]; !hasSecurityContext {
			response.Suggestions = append(response.Suggestions, "建议配置安全上下文(securityContext)以提高安全性")
		}
	}
}

// parseYAMLToPod 解析YAML为Pod对象
func (s *K8sWorkloadServiceImpl) parseYAMLToPod(yamlContent string) (*corev1.Pod, error) {
	// 创建解码器
	scheme := runtime.NewScheme()
	if err := corev1.AddToScheme(scheme); err != nil {
		return nil, fmt.Errorf("添加core/v1到scheme失败: %v", err)
	}
	// 添加apps/v1以支持Deployment等资源
	if err := appsv1.AddToScheme(scheme); err != nil {
		return nil, fmt.Errorf("添加apps/v1到scheme失败: %v", err)
	}

	decode := serializer.NewCodecFactory(scheme).UniversalDeserializer().Decode

	// 将YAML转换为JSON
	jsonData, err := yaml.YAMLToJSON([]byte(yamlContent))
	if err != nil {
		return nil, fmt.Errorf("YAML转JSON失败: %v", err)
	}

	// 解码为Kubernetes对象
	obj, _, err := decode(jsonData, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("解码Kubernetes对象失败: %v", err)
	}

	// 类型断言为Pod
	pod, ok := obj.(*corev1.Pod)
	if !ok {
		return nil, fmt.Errorf("YAML内容不是有效的Pod配置，实际类型: %T", obj)
	}

	return pod, nil
}

// parseYAMLToKubernetesObject 解析YAML为任意Kubernetes对象
func (s *K8sWorkloadServiceImpl) parseYAMLToKubernetesObject(yamlContent string) (runtime.Object, error) {
	// 创建解码器
	scheme := runtime.NewScheme()
	if err := corev1.AddToScheme(scheme); err != nil {
		return nil, fmt.Errorf("添加core/v1到scheme失败: %v", err)
	}
	// 添加apps/v1以支持Deployment等资源
	if err := appsv1.AddToScheme(scheme); err != nil {
		return nil, fmt.Errorf("添加apps/v1到scheme失败: %v", err)
	}
	// 添加batch/v1以支持Job等资源
	if err := batchv1.AddToScheme(scheme); err != nil {
		return nil, fmt.Errorf("添加batch/v1到scheme失败: %v", err)
	}
	// 添加networking/v1以支持Ingress等资源
	if err := networkingv1.AddToScheme(scheme); err != nil {
		return nil, fmt.Errorf("添加networking/v1到scheme失败: %v", err)
	}

	decode := serializer.NewCodecFactory(scheme).UniversalDeserializer().Decode

	// 处理API版本兼容性，特别是Ingress v1beta1到v1的转换
	processedContent, err := s.convertAPIVersionCompatibility(yamlContent)
	if err != nil {
		return nil, fmt.Errorf("API版本兼容性处理失败: %v", err)
	}

	// 将YAML转换为JSON
	jsonData, err := yaml.YAMLToJSON([]byte(processedContent))
	if err != nil {
		return nil, fmt.Errorf("YAML转JSON失败: %v", err)
	}

	// 解码为Kubernetes对象
	obj, _, err := decode(jsonData, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("解码Kubernetes对象失败: %v", err)
	}

	return obj, nil
}

// GetWorkloadPods 获取工作负载下的Pod列表
func (s *K8sWorkloadServiceImpl) GetWorkloadPods(c *gin.Context, clusterId uint, namespaceName string, workloadType string, workloadName string) {
	// 获取集群信息
	cluster, err := s.clusterDao.GetByID(clusterId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result.Failed(c, http.StatusNotFound, "集群不存在")
			return
		}
		result.Failed(c, http.StatusInternalServerError, "获取集群信息失败: "+err.Error())
		return
	}

	// 创建K8s客户端
	clientset, err := s.createK8sClient(cluster.Credential)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "连接K8s集群失败: "+err.Error())
		return
	}

	var pods []model.K8sPodInfo

	// 根据工作负载类型获取Pod列表
	switch strings.ToLower(workloadType) {
	case "deployment", "deployments":
		deployment, err := clientset.AppsV1().Deployments(namespaceName).Get(context.TODO(), workloadName, metav1.GetOptions{})
		if err != nil {
			result.Failed(c, http.StatusNotFound, "Deployment不存在: "+err.Error())
			return
		}
		pods, err = s.getWorkloadPods(clientset, namespaceName, deployment.Spec.Selector)
		if err != nil {
			result.Failed(c, http.StatusInternalServerError, "获取Pod列表失败: "+err.Error())
			return
		}

	case "statefulset", "statefulsets":
		statefulSet, err := clientset.AppsV1().StatefulSets(namespaceName).Get(context.TODO(), workloadName, metav1.GetOptions{})
		if err != nil {
			result.Failed(c, http.StatusNotFound, "StatefulSet不存在: "+err.Error())
			return
		}
		pods, err = s.getWorkloadPods(clientset, namespaceName, statefulSet.Spec.Selector)
		if err != nil {
			result.Failed(c, http.StatusInternalServerError, "获取Pod列表失败: "+err.Error())
			return
		}

	case "daemonset", "daemonsets":
		daemonSet, err := clientset.AppsV1().DaemonSets(namespaceName).Get(context.TODO(), workloadName, metav1.GetOptions{})
		if err != nil {
			result.Failed(c, http.StatusNotFound, "DaemonSet不存在: "+err.Error())
			return
		}
		pods, err = s.getWorkloadPods(clientset, namespaceName, daemonSet.Spec.Selector)
		if err != nil {
			result.Failed(c, http.StatusInternalServerError, "获取Pod列表失败: "+err.Error())
			return
		}

	case "job", "jobs":
		job, err := clientset.BatchV1().Jobs(namespaceName).Get(context.TODO(), workloadName, metav1.GetOptions{})
		if err != nil {
			result.Failed(c, http.StatusNotFound, "Job不存在: "+err.Error())
			return
		}
		pods, err = s.getWorkloadPods(clientset, namespaceName, job.Spec.Selector)
		if err != nil {
			result.Failed(c, http.StatusInternalServerError, "获取Pod列表失败: "+err.Error())
			return
		}

	case "cronjob", "cronjobs":
		// CronJob比较特殊，需要获取其创建的Job，然后获取Job的Pod
		jobs, err := clientset.BatchV1().Jobs(namespaceName).List(context.TODO(), metav1.ListOptions{
			LabelSelector: fmt.Sprintf("job-name contains %s", workloadName),
		})
		if err != nil {
			result.Failed(c, http.StatusInternalServerError, "获取CronJob的Job列表失败: "+err.Error())
			return
		}

		allPods := []model.K8sPodInfo{}
		for _, job := range jobs.Items {
			jobPods, err := s.getWorkloadPods(clientset, namespaceName, job.Spec.Selector)
			if err != nil {
				continue // 忽略获取失败的Job
			}
			allPods = append(allPods, jobPods...)
		}
		pods = allPods

	default:
		result.Failed(c, http.StatusBadRequest, "不支持的工作负载类型: "+workloadType)
		return
	}

	result.Success(c, pods)
}

// ===================== Deployment 版本回滚管理实现 =====================

// GetDeploymentHistory 获取Deployment版本历史
func (s *K8sWorkloadServiceImpl) GetDeploymentHistory(c *gin.Context, clusterId uint, namespaceName string, deploymentName string) {
	cluster, err := s.clusterDao.GetByID(clusterId)
	if err != nil {
		result.Failed(c, http.StatusNotFound, "集群不存在")
		return
	}

	clientset, err := s.createK8sClient(cluster.Credential)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "连接K8s集群失败: "+err.Error())
		return
	}

	// 获取Deployment
	deployment, err := clientset.AppsV1().Deployments(namespaceName).Get(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		result.Failed(c, http.StatusNotFound, "Deployment不存在: "+err.Error())
		return
	}

	// 获取关联的ReplicaSets
	replicaSets, err := clientset.AppsV1().ReplicaSets(namespaceName).List(context.TODO(), metav1.ListOptions{
		LabelSelector: metav1.FormatLabelSelector(deployment.Spec.Selector),
	})
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "获取ReplicaSet列表失败: "+err.Error())
		return
	}

	// 构建版本历史
	var revisions []model.DeploymentRevision
	currentRevision := deployment.Annotations["deployment.kubernetes.io/revision"]

	for _, rs := range replicaSets.Items {
		if rsRevision, exists := rs.Annotations["deployment.kubernetes.io/revision"]; exists {
			revision, _ := strconv.ParseInt(rsRevision, 10, 64)

			// 获取镜像列表
			var images []string
			for _, container := range rs.Spec.Template.Spec.Containers {
				images = append(images, container.Image)
			}

			status := "historical"
			if rsRevision == currentRevision {
				status = "current"
			}

			revisions = append(revisions, model.DeploymentRevision{
				Revision:     revision,
				CreationTime: rs.CreationTimestamp.Format(time.RFC3339),
				ChangeReason: rs.Annotations["kubernetes.io/change-cause"],
				Images:       images,
				Labels:       rs.Labels,
				Annotations:  rs.Annotations,
				Status:       status,
				ReplicasSummary: model.ReplicasSummary{
					Desired:   *rs.Spec.Replicas,
					Current:   rs.Status.Replicas,
					Updated:   rs.Status.Replicas,
					Ready:     rs.Status.ReadyReplicas,
					Available: rs.Status.AvailableReplicas,
				},
			})
		}
	}

	// 按版本号排序
	for i := 0; i < len(revisions)-1; i++ {
		for j := i + 1; j < len(revisions); j++ {
			if revisions[i].Revision < revisions[j].Revision {
				revisions[i], revisions[j] = revisions[j], revisions[i]
			}
		}
	}

	currentRev, _ := strconv.ParseInt(currentRevision, 10, 64)
	response := model.DeploymentRolloutHistoryResponse{
		DeploymentName:  deploymentName,
		Namespace:       namespaceName,
		CurrentRevision: currentRev,
		TotalRevisions:  len(revisions),
		Revisions:       revisions,
	}

	result.Success(c, response)
}

// GetDeploymentRevision 获取Deployment指定版本详情
func (s *K8sWorkloadServiceImpl) GetDeploymentRevision(c *gin.Context, clusterId uint, namespaceName string, deploymentName string, revision int64) {
	cluster, err := s.clusterDao.GetByID(clusterId)
	if err != nil {
		result.Failed(c, http.StatusNotFound, "集群不存在")
		return
	}

	clientset, err := s.createK8sClient(cluster.Credential)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "连接K8s集群失败: "+err.Error())
		return
	}

	// 获取Deployment
	deployment, err := clientset.AppsV1().Deployments(namespaceName).Get(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		result.Failed(c, http.StatusNotFound, "Deployment不存在: "+err.Error())
		return
	}

	// 查找指定版本的ReplicaSet
	replicaSets, err := clientset.AppsV1().ReplicaSets(namespaceName).List(context.TODO(), metav1.ListOptions{
		LabelSelector: metav1.FormatLabelSelector(deployment.Spec.Selector),
	})
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "获取ReplicaSet列表失败: "+err.Error())
		return
	}

	var targetRS *appsv1.ReplicaSet
	revisionStr := strconv.FormatInt(revision, 10)
	for _, rs := range replicaSets.Items {
		if rs.Annotations["deployment.kubernetes.io/revision"] == revisionStr {
			targetRS = &rs
			break
		}
	}

	if targetRS == nil {
		result.Failed(c, http.StatusNotFound, fmt.Sprintf("版本 %d 不存在", revision))
		return
	}

	// 获取镜像列表
	var images []string
	for _, container := range targetRS.Spec.Template.Spec.Containers {
		images = append(images, container.Image)
	}

	// 获取相关事件
	events, err := clientset.CoreV1().Events(namespaceName).List(context.TODO(), metav1.ListOptions{
		FieldSelector: fmt.Sprintf("involvedObject.name=%s", targetRS.Name),
	})
	if err != nil {
		// 忽略事件获取失败
	}

	var eventList []model.K8sEvent
	if events != nil {
		for _, event := range events.Items {
			eventList = append(eventList, model.K8sEvent{
				Type:      event.Type,
				Reason:    event.Reason,
				Message:   event.Message,
				Source:    event.Source.Component,
				Count:     event.Count,
				FirstTime: event.FirstTimestamp.Format(time.RFC3339),
				LastTime:  event.LastTimestamp.Format(time.RFC3339),
			})
		}
	}

	// 构建Pod模板
	podTemplate := model.PodTemplateSpec{
		Labels:     targetRS.Spec.Template.Labels,
		Containers: []model.ContainerSpec{},
	}

	for _, container := range targetRS.Spec.Template.Spec.Containers {
		containerSpec := model.ContainerSpec{
			Name:  container.Name,
			Image: container.Image,
			Resources: model.WorkloadResources{
				Requests: model.ResourceSpec{},
				Limits:   model.ResourceSpec{},
			},
		}

		if container.Resources.Requests != nil {
			if cpu, exists := container.Resources.Requests["cpu"]; exists && !cpu.IsZero() {
				containerSpec.Resources.Requests.CPU = cpu.String()
			}
			if memory, exists := container.Resources.Requests["memory"]; exists && !memory.IsZero() {
				containerSpec.Resources.Requests.Memory = memory.String()
			}
		}

		if container.Resources.Limits != nil {
			if cpu, exists := container.Resources.Limits["cpu"]; exists && !cpu.IsZero() {
				containerSpec.Resources.Limits.CPU = cpu.String()
			}
			if memory, exists := container.Resources.Limits["memory"]; exists && !memory.IsZero() {
				containerSpec.Resources.Limits.Memory = memory.String()
			}
		}

		podTemplate.Containers = append(podTemplate.Containers, containerSpec)
	}

	currentRevision := deployment.Annotations["deployment.kubernetes.io/revision"]
	status := "historical"
	if revisionStr == currentRevision {
		status = "current"
	}

	detail := model.DeploymentRevisionDetail{
		DeploymentRevision: model.DeploymentRevision{
			Revision:     revision,
			CreationTime: targetRS.CreationTimestamp.Format(time.RFC3339),
			ChangeReason: targetRS.Annotations["kubernetes.io/change-cause"],
			Images:       images,
			Labels:       targetRS.Labels,
			Annotations:  targetRS.Annotations,
			Status:       status,
			ReplicasSummary: model.ReplicasSummary{
				Desired:   *targetRS.Spec.Replicas,
				Current:   targetRS.Status.Replicas,
				Updated:   targetRS.Status.Replicas,
				Ready:     targetRS.Status.ReadyReplicas,
				Available: targetRS.Status.AvailableReplicas,
			},
		},
		PodTemplate: podTemplate,
		Strategy: model.DeploymentStrategy{
			Type: string(deployment.Spec.Strategy.Type),
		},
		Events: eventList,
	}

	if deployment.Spec.Strategy.RollingUpdate != nil {
		detail.Strategy.RollingUpdate = model.RollingUpdateDeployment{
			MaxUnavailable: deployment.Spec.Strategy.RollingUpdate.MaxUnavailable.String(),
			MaxSurge:       deployment.Spec.Strategy.RollingUpdate.MaxSurge.String(),
		}
	}

	result.Success(c, detail)
}

// RollbackDeployment 回滚Deployment到指定版本
func (s *K8sWorkloadServiceImpl) RollbackDeployment(c *gin.Context, clusterId uint, namespaceName string, deploymentName string, req *model.RollbackDeploymentRequest) {
	cluster, err := s.clusterDao.GetByID(clusterId)
	if err != nil {
		result.Failed(c, http.StatusNotFound, "集群不存在")
		return
	}

	clientset, err := s.createK8sClient(cluster.Credential)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "连接K8s集群失败: "+err.Error())
		return
	}

	// 获取Deployment
	deployment, err := clientset.AppsV1().Deployments(namespaceName).Get(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		result.Failed(c, http.StatusNotFound, "Deployment不存在: "+err.Error())
		return
	}

	fromRevision, _ := strconv.ParseInt(deployment.Annotations["deployment.kubernetes.io/revision"], 10, 64)
	targetRevision := req.ToRevision

	// 如果目标版本为0，回滚到上一个版本
	if targetRevision == 0 {
		replicaSets, err := clientset.AppsV1().ReplicaSets(namespaceName).List(context.TODO(), metav1.ListOptions{
			LabelSelector: metav1.FormatLabelSelector(deployment.Spec.Selector),
		})
		if err != nil {
			result.Failed(c, http.StatusInternalServerError, "获取ReplicaSet列表失败: "+err.Error())
			return
		}

		// 找到最近的历史版本
		var maxHistoricalRevision int64 = 0
		currentRevisionStr := deployment.Annotations["deployment.kubernetes.io/revision"]

		for _, rs := range replicaSets.Items {
			if rsRevision, exists := rs.Annotations["deployment.kubernetes.io/revision"]; exists && rsRevision != currentRevisionStr {
				revision, _ := strconv.ParseInt(rsRevision, 10, 64)
				if revision > maxHistoricalRevision {
					maxHistoricalRevision = revision
				}
			}
		}

		if maxHistoricalRevision == 0 {
			result.Failed(c, http.StatusBadRequest, "没有可回滚的历史版本")
			return
		}
		targetRevision = maxHistoricalRevision
	}

	// 查找目标版本的ReplicaSet
	replicaSets, err := clientset.AppsV1().ReplicaSets(namespaceName).List(context.TODO(), metav1.ListOptions{
		LabelSelector: metav1.FormatLabelSelector(deployment.Spec.Selector),
	})
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "获取ReplicaSet列表失败: "+err.Error())
		return
	}

	var targetRS *appsv1.ReplicaSet
	targetRevisionStr := strconv.FormatInt(targetRevision, 10)
	for _, rs := range replicaSets.Items {
		if rs.Annotations["deployment.kubernetes.io/revision"] == targetRevisionStr {
			targetRS = &rs
			break
		}
	}

	if targetRS == nil {
		result.Failed(c, http.StatusNotFound, fmt.Sprintf("目标版本 %d 不存在", targetRevision))
		return
	}

	// 执行回滚：将目标版本的模板应用到Deployment
	deployment.Spec.Template = targetRS.Spec.Template

	// 添加回滚原因注释
	if deployment.Annotations == nil {
		deployment.Annotations = make(map[string]string)
	}
	deployment.Annotations["kubernetes.io/change-cause"] = fmt.Sprintf("kubectl rollout undo deployment/%s --to-revision=%d", deploymentName, targetRevision)

	_, err = clientset.AppsV1().Deployments(namespaceName).Update(context.TODO(), deployment, metav1.UpdateOptions{})
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "回滚Deployment失败: "+err.Error())
		return
	}

	response := model.RollbackDeploymentResponse{
		Success:        true,
		Message:        fmt.Sprintf("Deployment '%s' 已成功回滚到版本 %d", deploymentName, targetRevision),
		FromRevision:   fromRevision,
		ToRevision:     targetRevision,
		DeploymentName: deploymentName,
		Namespace:      namespaceName,
		RolloutStatus:  "Progressing",
	}

	result.Success(c, response)
}

// PauseDeployment 暂停Deployment滚动更新
func (s *K8sWorkloadServiceImpl) PauseDeployment(c *gin.Context, clusterId uint, namespaceName string, deploymentName string) {
	cluster, err := s.clusterDao.GetByID(clusterId)
	if err != nil {
		result.Failed(c, http.StatusNotFound, "集群不存在")
		return
	}

	clientset, err := s.createK8sClient(cluster.Credential)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "连接K8s集群失败: "+err.Error())
		return
	}

	// 获取Deployment
	deployment, err := clientset.AppsV1().Deployments(namespaceName).Get(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		result.Failed(c, http.StatusNotFound, "Deployment不存在: "+err.Error())
		return
	}

	// 检查是否已经暂停
	if deployment.Spec.Paused {
		response := model.PauseDeploymentResponse{
			Success:        false,
			Message:        fmt.Sprintf("Deployment '%s' 已经处于暂停状态", deploymentName),
			DeploymentName: deploymentName,
			Namespace:      namespaceName,
			Status:         "Already Paused",
		}
		result.Success(c, response)
		return
	}

	// 暂停Deployment
	deployment.Spec.Paused = true
	_, err = clientset.AppsV1().Deployments(namespaceName).Update(context.TODO(), deployment, metav1.UpdateOptions{})
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "暂停Deployment失败: "+err.Error())
		return
	}

	response := model.PauseDeploymentResponse{
		Success:        true,
		Message:        fmt.Sprintf("Deployment '%s' 已成功暂停", deploymentName),
		DeploymentName: deploymentName,
		Namespace:      namespaceName,
		Status:         "Paused",
	}

	result.Success(c, response)
}

// ResumeDeployment 恢复Deployment滚动更新
func (s *K8sWorkloadServiceImpl) ResumeDeployment(c *gin.Context, clusterId uint, namespaceName string, deploymentName string) {
	cluster, err := s.clusterDao.GetByID(clusterId)
	if err != nil {
		result.Failed(c, http.StatusNotFound, "集群不存在")
		return
	}

	clientset, err := s.createK8sClient(cluster.Credential)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "连接K8s集群失败: "+err.Error())
		return
	}

	// 获取Deployment
	deployment, err := clientset.AppsV1().Deployments(namespaceName).Get(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		result.Failed(c, http.StatusNotFound, "Deployment不存在: "+err.Error())
		return
	}

	// 检查是否已经运行
	if !deployment.Spec.Paused {
		response := model.ResumeDeploymentResponse{
			Success:        false,
			Message:        fmt.Sprintf("Deployment '%s' 未处于暂停状态", deploymentName),
			DeploymentName: deploymentName,
			Namespace:      namespaceName,
			Status:         "Already Running",
		}
		result.Success(c, response)
		return
	}

	// 恢复Deployment
	deployment.Spec.Paused = false
	_, err = clientset.AppsV1().Deployments(namespaceName).Update(context.TODO(), deployment, metav1.UpdateOptions{})
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "恢复Deployment失败: "+err.Error())
		return
	}

	response := model.ResumeDeploymentResponse{
		Success:        true,
		Message:        fmt.Sprintf("Deployment '%s' 已成功恢复", deploymentName),
		DeploymentName: deploymentName,
		Namespace:      namespaceName,
		Status:         "Progressing",
	}

	result.Success(c, response)
}

// GetDeploymentRolloutStatus 获取Deployment滚动发布状态
func (s *K8sWorkloadServiceImpl) GetDeploymentRolloutStatus(c *gin.Context, clusterId uint, namespaceName string, deploymentName string) {
	cluster, err := s.clusterDao.GetByID(clusterId)
	if err != nil {
		result.Failed(c, http.StatusNotFound, "集群不存在")
		return
	}

	clientset, err := s.createK8sClient(cluster.Credential)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "连接K8s集群失败: "+err.Error())
		return
	}

	// 获取Deployment
	deployment, err := clientset.AppsV1().Deployments(namespaceName).Get(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		result.Failed(c, http.StatusNotFound, "Deployment不存在: "+err.Error())
		return
	}

	currentRevision, _ := strconv.ParseInt(deployment.Annotations["deployment.kubernetes.io/revision"], 10, 64)

	// 转换状态条件
	var conditions []model.WorkloadCondition
	for _, condition := range deployment.Status.Conditions {
		conditions = append(conditions, model.WorkloadCondition{
			Type:               string(condition.Type),
			Status:             string(condition.Status),
			LastTransitionTime: condition.LastTransitionTime.Format(time.RFC3339),
			Reason:             condition.Reason,
			Message:            condition.Message,
		})
	}

	// 判断滚动发布状态
	status := "Progressing"
	rolloutComplete := false

	if deployment.Spec.Paused {
		status = "Paused"
	} else if deployment.Status.UpdatedReplicas == *deployment.Spec.Replicas &&
		deployment.Status.ReadyReplicas == *deployment.Spec.Replicas &&
		deployment.Status.AvailableReplicas == *deployment.Spec.Replicas {
		status = "Complete"
		rolloutComplete = true
	} else {
		// 检查是否有失败条件
		for _, condition := range deployment.Status.Conditions {
			if condition.Type == appsv1.DeploymentProgressing && condition.Status == corev1.ConditionFalse {
				status = "Failed"
				break
			}
		}
	}

	progressDeadline := int32(600) // 默认10分钟
	if deployment.Spec.ProgressDeadlineSeconds != nil {
		progressDeadline = *deployment.Spec.ProgressDeadlineSeconds
	}

	response := model.DeploymentRolloutStatusResponse{
		DeploymentName:     deploymentName,
		Namespace:          namespaceName,
		CurrentRevision:    currentRevision,
		UpdatedReplicas:    deployment.Status.UpdatedReplicas,
		ReadyReplicas:      deployment.Status.ReadyReplicas,
		AvailableReplicas:  deployment.Status.AvailableReplicas,
		ObservedGeneration: deployment.Status.ObservedGeneration,
		Conditions:         conditions,
		Strategy: model.DeploymentStrategy{
			Type: string(deployment.Spec.Strategy.Type),
		},
		Paused:           deployment.Spec.Paused,
		ProgressDeadline: progressDeadline,
		RolloutComplete:  rolloutComplete,
		Status:           status,
	}

	if deployment.Spec.Strategy.RollingUpdate != nil {
		response.Strategy.RollingUpdate = model.RollingUpdateDeployment{
			MaxUnavailable: deployment.Spec.Strategy.RollingUpdate.MaxUnavailable.String(),
			MaxSurge:       deployment.Spec.Strategy.RollingUpdate.MaxSurge.String(),
		}
	}

	result.Success(c, response)
}

// convertAPIVersionCompatibility 处理API版本兼容性转换
func (s *K8sWorkloadServiceImpl) convertAPIVersionCompatibility(yamlContent string) (string, error) {
	// 解析YAML为通用对象
	var obj map[string]interface{}
	if err := yaml.Unmarshal([]byte(yamlContent), &obj); err != nil {
		return yamlContent, nil // 如果解析失败，返回原始内容
	}

	// 检查API版本
	apiVersion, ok := obj["apiVersion"].(string)
	if !ok {
		return yamlContent, nil // 如果没有apiVersion字段，直接返回
	}

	// 如果是Ingress v1beta1版本，需要转换到v1
	if apiVersion == "networking.k8s.io/v1beta1" {
		kind, ok := obj["kind"].(string)
		if ok && kind == "Ingress" {
			obj["apiVersion"] = "networking.k8s.io/v1"

			// 转换spec.rules[].http.paths[].backend格式
			if spec, ok := obj["spec"].(map[string]interface{}); ok {
				if rules, ok := spec["rules"].([]interface{}); ok {
					for _, rule := range rules {
						if ruleMap, ok := rule.(map[string]interface{}); ok {
							if http, ok := ruleMap["http"].(map[string]interface{}); ok {
								if paths, ok := http["paths"].([]interface{}); ok {
									for _, path := range paths {
										if pathMap, ok := path.(map[string]interface{}); ok {
											if backend, ok := pathMap["backend"].(map[string]interface{}); ok {
												// v1beta1格式: backend.serviceName, backend.servicePort
												// v1格式: backend.service.name, backend.service.port.number
												if serviceName, hasServiceName := backend["serviceName"]; hasServiceName {
													if servicePort, hasServicePort := backend["servicePort"]; hasServicePort {
														// 删除旧格式字段
														delete(backend, "serviceName")
														delete(backend, "servicePort")

														// 添加新格式字段
														backend["service"] = map[string]interface{}{
															"name": serviceName,
															"port": map[string]interface{}{
																"number": servicePort,
															},
														}
													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}

			// 重新序列化为YAML
			convertedBytes, err := yaml.Marshal(obj)
			if err != nil {
				return yamlContent, nil // 如果序列化失败，返回原始内容
			}
			return string(convertedBytes), nil
		}
	}

	return yamlContent, nil
}

// ===================== Pod更新辅助方法 =====================

// analyzePodChanges 分析Pod变更
func (s *K8sWorkloadServiceImpl) analyzePodChanges(existingPod, newPod *corev1.Pod) []string {
	changes := []string{}

	// 比较labels
	if !s.compareMaps(existingPod.Labels, newPod.Labels) {
		changes = append(changes, "Labels发生变更")
	}

	// 比较annotations
	if !s.compareMaps(existingPod.Annotations, newPod.Annotations) {
		changes = append(changes, "Annotations发生变更")
	}

	// 比较容器镜像
	if !s.compareContainerImages(existingPod.Spec.Containers, newPod.Spec.Containers) {
		changes = append(changes, "容器镜像发生变更")
	}

	// 比较容器端口
	if !s.compareContainerPorts(existingPod.Spec.Containers, newPod.Spec.Containers) {
		changes = append(changes, "容器端口配置发生变更")
	}

	// 比较环境变量
	if !s.compareContainerEnvs(existingPod.Spec.Containers, newPod.Spec.Containers) {
		changes = append(changes, "环境变量发生变更")
	}

	// 比较资源配置
	if !s.compareContainerResources(existingPod.Spec.Containers, newPod.Spec.Containers) {
		changes = append(changes, "资源配置发生变更")
	}

	if len(changes) == 0 {
		changes = append(changes, "未检测到显著变更")
	}

	return changes
}

// determinePodUpdateStrategy 确定Pod更新策略
func (s *K8sWorkloadServiceImpl) determinePodUpdateStrategy(existingPod, newPod *corev1.Pod) string {
	// 检查是否有不可变字段发生变化

	// 容器镜像变更需要重建
	if !s.compareContainerImages(existingPod.Spec.Containers, newPod.Spec.Containers) {
		return "recreate"
	}

	// 容器端口变更需要重建
	if !s.compareContainerPorts(existingPod.Spec.Containers, newPod.Spec.Containers) {
		return "recreate"
	}

	// 环境变量变更需要重建
	if !s.compareContainerEnvs(existingPod.Spec.Containers, newPod.Spec.Containers) {
		return "recreate"
	}

	// 资源配置变更需要重建
	if !s.compareContainerResources(existingPod.Spec.Containers, newPod.Spec.Containers) {
		return "recreate"
	}

	// 卷配置变更需要重建
	if !s.compareVolumes(existingPod.Spec.Volumes, newPod.Spec.Volumes) {
		return "recreate"
	}

	// 只有labels和annotations变更可以原地更新
	return "patch"
}

// patchPod 原地更新Pod
func (s *K8sWorkloadServiceImpl) patchPod(clientset *kubernetes.Clientset, ctx context.Context, existingPod, newPod *corev1.Pod) error {
	// 只更新可变字段
	patchPod := existingPod.DeepCopy()

	// 更新labels
	if patchPod.Labels == nil {
		patchPod.Labels = make(map[string]string)
	}
	for k, v := range newPod.Labels {
		patchPod.Labels[k] = v
	}

	// 更新annotations
	if patchPod.Annotations == nil {
		patchPod.Annotations = make(map[string]string)
	}
	for k, v := range newPod.Annotations {
		patchPod.Annotations[k] = v
	}

	// 执行更新
	_, err := clientset.CoreV1().Pods(existingPod.Namespace).Update(ctx, patchPod, metav1.UpdateOptions{})
	return err
}

// 辅助比较方法
func (s *K8sWorkloadServiceImpl) compareMaps(map1, map2 map[string]string) bool {
	if len(map1) != len(map2) {
		return false
	}
	for k, v := range map1 {
		if map2[k] != v {
			return false
		}
	}
	return true
}

func (s *K8sWorkloadServiceImpl) compareContainerImages(containers1, containers2 []corev1.Container) bool {
	if len(containers1) != len(containers2) {
		return false
	}
	for i, c1 := range containers1 {
		if i >= len(containers2) || c1.Image != containers2[i].Image {
			return false
		}
	}
	return true
}

func (s *K8sWorkloadServiceImpl) compareContainerPorts(containers1, containers2 []corev1.Container) bool {
	if len(containers1) != len(containers2) {
		return false
	}
	for i, c1 := range containers1 {
		if i >= len(containers2) {
			return false
		}
		c2 := containers2[i]
		if len(c1.Ports) != len(c2.Ports) {
			return false
		}
		for j, p1 := range c1.Ports {
			if j >= len(c2.Ports) || p1.ContainerPort != c2.Ports[j].ContainerPort || p1.Protocol != c2.Ports[j].Protocol {
				return false
			}
		}
	}
	return true
}

func (s *K8sWorkloadServiceImpl) compareContainerEnvs(containers1, containers2 []corev1.Container) bool {
	if len(containers1) != len(containers2) {
		return false
	}
	for i, c1 := range containers1 {
		if i >= len(containers2) {
			return false
		}
		c2 := containers2[i]
		if len(c1.Env) != len(c2.Env) {
			return false
		}
		for j, e1 := range c1.Env {
			if j >= len(c2.Env) || e1.Name != c2.Env[j].Name || e1.Value != c2.Env[j].Value {
				return false
			}
		}
	}
	return true
}

func (s *K8sWorkloadServiceImpl) compareContainerResources(containers1, containers2 []corev1.Container) bool {
	if len(containers1) != len(containers2) {
		return false
	}
	for i, c1 := range containers1 {
		if i >= len(containers2) {
			return false
		}
		c2 := containers2[i]

		// 比较requests
		if !s.compareResourceList(c1.Resources.Requests, c2.Resources.Requests) {
			return false
		}

		// 比较limits
		if !s.compareResourceList(c1.Resources.Limits, c2.Resources.Limits) {
			return false
		}
	}
	return true
}

func (s *K8sWorkloadServiceImpl) compareResourceList(rl1, rl2 corev1.ResourceList) bool {
	if len(rl1) != len(rl2) {
		return false
	}
	for k, v1 := range rl1 {
		v2, exists := rl2[k]
		if !exists || !v1.Equal(v2) {
			return false
		}
	}
	return true
}

func (s *K8sWorkloadServiceImpl) compareVolumes(volumes1, volumes2 []corev1.Volume) bool {
	if len(volumes1) != len(volumes2) {
		return false
	}
	for i, v1 := range volumes1 {
		if i >= len(volumes2) || v1.Name != volumes2[i].Name {
			return false
		}
	}
	return true
}

// ===================== 通用工作负载YAML管理API实现 =====================

// GetWorkloadYaml 获取工作负载的YAML配置
func (s *K8sWorkloadServiceImpl) GetWorkloadYaml(c *gin.Context, clusterId uint, namespaceName string, workloadType string, workloadName string) {
	// 获取集群信息
	cluster, err := s.clusterDao.GetByID(clusterId)
	if err != nil {
		result.Failed(c, http.StatusNotFound, "集群不存在")
		return
	}

	// 创建k8s客户端
	config, err := clientcmd.RESTConfigFromKubeConfig([]byte(cluster.Credential))
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "解析集群配置失败: "+err.Error())
		return
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "创建k8s客户端失败: "+err.Error())
		return
	}

	ctx := context.TODO()
	var yamlContent string

	// 根据工作负载类型获取对应的资源YAML
	switch strings.ToLower(workloadType) {
	case "deployment":
		deployment, err := clientset.AppsV1().Deployments(namespaceName).Get(ctx, workloadName, metav1.GetOptions{})
		if err != nil {
			result.Failed(c, http.StatusNotFound, "Deployment不存在: "+err.Error())
			return
		}
		// 清理系统字段，但保留APIVersion和Kind
		deployment.ManagedFields = nil
		deployment.ResourceVersion = ""
		deployment.UID = ""
		deployment.SelfLink = ""
		deployment.Generation = 0
		deployment.CreationTimestamp = metav1.Time{}
		deployment.Status = appsv1.DeploymentStatus{}
		// 确保APIVersion和Kind字段存在
		deployment.APIVersion = "apps/v1"
		deployment.Kind = "Deployment"

		yamlBytes, err := yaml.Marshal(deployment)
		if err != nil {
			result.Failed(c, http.StatusInternalServerError, "序列化YAML失败: "+err.Error())
			return
		}
		yamlContent = string(yamlBytes)

	case "statefulset":
		statefulset, err := clientset.AppsV1().StatefulSets(namespaceName).Get(ctx, workloadName, metav1.GetOptions{})
		if err != nil {
			result.Failed(c, http.StatusNotFound, "StatefulSet不存在: "+err.Error())
			return
		}
		// 清理系统字段，但保留APIVersion和Kind
		statefulset.ManagedFields = nil
		statefulset.ResourceVersion = ""
		statefulset.UID = ""
		statefulset.SelfLink = ""
		statefulset.Generation = 0
		statefulset.CreationTimestamp = metav1.Time{}
		statefulset.Status = appsv1.StatefulSetStatus{}
		// 确保APIVersion和Kind字段存在
		statefulset.APIVersion = "apps/v1"
		statefulset.Kind = "StatefulSet"

		yamlBytes, err := yaml.Marshal(statefulset)
		if err != nil {
			result.Failed(c, http.StatusInternalServerError, "序列化YAML失败: "+err.Error())
			return
		}
		yamlContent = string(yamlBytes)

	case "daemonset":
		daemonset, err := clientset.AppsV1().DaemonSets(namespaceName).Get(ctx, workloadName, metav1.GetOptions{})
		if err != nil {
			result.Failed(c, http.StatusNotFound, "DaemonSet不存在: "+err.Error())
			return
		}
		// 清理系统字段，但保留APIVersion和Kind
		daemonset.ManagedFields = nil
		daemonset.ResourceVersion = ""
		daemonset.UID = ""
		daemonset.SelfLink = ""
		daemonset.Generation = 0
		daemonset.CreationTimestamp = metav1.Time{}
		daemonset.Status = appsv1.DaemonSetStatus{}
		// 确保APIVersion和Kind字段存在
		daemonset.APIVersion = "apps/v1"
		daemonset.Kind = "DaemonSet"

		yamlBytes, err := yaml.Marshal(daemonset)
		if err != nil {
			result.Failed(c, http.StatusInternalServerError, "序列化YAML失败: "+err.Error())
			return
		}
		yamlContent = string(yamlBytes)

	case "job":
		job, err := clientset.BatchV1().Jobs(namespaceName).Get(ctx, workloadName, metav1.GetOptions{})
		if err != nil {
			result.Failed(c, http.StatusNotFound, "Job不存在: "+err.Error())
			return
		}
		// 清理系统字段，但保留APIVersion和Kind
		job.ManagedFields = nil
		job.ResourceVersion = ""
		job.UID = ""
		job.SelfLink = ""
		job.Generation = 0
		job.CreationTimestamp = metav1.Time{}
		job.Status = batchv1.JobStatus{}
		// 确保APIVersion和Kind字段存在
		job.APIVersion = "batch/v1"
		job.Kind = "Job"

		yamlBytes, err := yaml.Marshal(job)
		if err != nil {
			result.Failed(c, http.StatusInternalServerError, "序列化YAML失败: "+err.Error())
			return
		}
		yamlContent = string(yamlBytes)

	case "cronjob":
		cronjob, err := clientset.BatchV1().CronJobs(namespaceName).Get(ctx, workloadName, metav1.GetOptions{})
		if err != nil {
			result.Failed(c, http.StatusNotFound, "CronJob不存在: "+err.Error())
			return
		}
		// 清理系统字段，但保留APIVersion和Kind
		cronjob.ManagedFields = nil
		cronjob.ResourceVersion = ""
		cronjob.UID = ""
		cronjob.SelfLink = ""
		cronjob.Generation = 0
		cronjob.CreationTimestamp = metav1.Time{}
		cronjob.Status = batchv1.CronJobStatus{}
		// 确保APIVersion和Kind字段存在
		cronjob.APIVersion = "batch/v1"
		cronjob.Kind = "CronJob"

		yamlBytes, err := yaml.Marshal(cronjob)
		if err != nil {
			result.Failed(c, http.StatusInternalServerError, "序列化YAML失败: "+err.Error())
			return
		}
		yamlContent = string(yamlBytes)

	default:
		result.Failed(c, http.StatusBadRequest, "不支持的工作负载类型: "+workloadType+". 支持的类型: deployment,statefulset,daemonset,job,cronjob")
		return
	}

	response := model.GetWorkloadYAMLResponse{
		Success:      true,
		WorkloadType: workloadType,
		WorkloadName: workloadName,
		Namespace:    namespaceName,
		YAMLContent:  yamlContent,
		Message:      fmt.Sprintf("成功获取%s '%s'的YAML配置", workloadType, workloadName),
	}

	result.Success(c, response)
}

// UpdateWorkloadYaml 通用工作负载YAML更新
func (s *K8sWorkloadServiceImpl) UpdateWorkloadYaml(c *gin.Context, clusterId uint, namespaceName string, req *model.UpdateWorkloadYAMLRequest) {
	// 获取集群信息
	cluster, err := s.clusterDao.GetByID(clusterId)
	if err != nil {
		result.Failed(c, http.StatusNotFound, "集群不存在")
		return
	}

	// 创建k8s客户端
	config, err := clientcmd.RESTConfigFromKubeConfig([]byte(cluster.Credential))
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "解析集群配置失败: "+err.Error())
		return
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "创建k8s客户端失败: "+err.Error())
		return
	}

	// 首先校验YAML格式
	validationResult, err := s.validateYAMLContent(req.YAMLContent, req.WorkloadType)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "YAML校验失败: "+err.Error())
		return
	}

	// 如果只是校验格式，直接返回校验结果
	if req.ValidateOnly {
		result.Success(c, model.UpdateWorkloadYAMLResponse{
			Success:          validationResult.Valid,
			WorkloadType:     req.WorkloadType,
			WorkloadName:     req.WorkloadName,
			Namespace:        namespaceName,
			Message:          "YAML校验完成",
			ValidationResult: validationResult,
		})
		return
	}

	ctx := context.TODO()
	response := &model.UpdateWorkloadYAMLResponse{
		Success:      false,
		WorkloadType: req.WorkloadType,
		WorkloadName: req.WorkloadName,
		Namespace:    namespaceName,
		Changes:      []string{},
		Warnings:     []string{},
	}

	// 根据工作负载类型处理更新
	switch strings.ToLower(req.WorkloadType) {
	case "deployment":
		err = s.updateDeploymentYAML(clientset, ctx, namespaceName, req, response)
	case "statefulset":
		err = s.updateStatefulSetYAML(clientset, ctx, namespaceName, req, response)
	case "daemonset":
		err = s.updateDaemonSetYAML(clientset, ctx, namespaceName, req, response)
	case "job":
		err = s.updateJobYAML(clientset, ctx, namespaceName, req, response)
	case "cronjob":
		err = s.updateCronJobYAML(clientset, ctx, namespaceName, req, response)
	default:
		result.Failed(c, http.StatusBadRequest, "不支持的工作负载类型: "+req.WorkloadType+". 支持的类型: deployment,statefulset,daemonset,job,cronjob")
		return
	}

	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "更新工作负载失败: "+err.Error())
		return
	}

	if response.Success {
		response.AppliedAt = time.Now().Format("2006-01-02 15:04:05")
	}

	result.Success(c, response)
}

// updateDeploymentYAML 更新Deployment YAML
func (s *K8sWorkloadServiceImpl) updateDeploymentYAML(clientset *kubernetes.Clientset, ctx context.Context, namespaceName string, req *model.UpdateWorkloadYAMLRequest, response *model.UpdateWorkloadYAMLResponse) error {
	// 解析新的YAML
	newDeployment := &appsv1.Deployment{}
	if err := yaml.Unmarshal([]byte(req.YAMLContent), newDeployment); err != nil {
		return fmt.Errorf("解析Deployment YAML失败: %v", err)
	}

	// 验证名称匹配
	if newDeployment.Name != req.WorkloadName {
		return fmt.Errorf("YAML中的Deployment名称(%s)与请求参数不匹配(%s)", newDeployment.Name, req.WorkloadName)
	}

	// 确保命名空间正确
	if newDeployment.Namespace == "" {
		newDeployment.Namespace = namespaceName
	} else if newDeployment.Namespace != namespaceName {
		return fmt.Errorf("YAML中的命名空间(%s)与URL参数不匹配(%s)", newDeployment.Namespace, namespaceName)
	}

	// 获取现有的Deployment
	existingDeployment, err := clientset.AppsV1().Deployments(namespaceName).Get(ctx, req.WorkloadName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("Deployment不存在: %v", err)
	}

	// DryRun模式
	if req.DryRun {
		changes := s.analyzeDeploymentChanges(existingDeployment, newDeployment)
		response.Success = true
		response.Message = "DryRun模式：仅分析变更，未实际更新"
		response.UpdateStrategy = "rolling"
		response.Changes = changes
		return nil
	}

	// 保留重要的系统字段
	newDeployment.ResourceVersion = existingDeployment.ResourceVersion
	newDeployment.UID = existingDeployment.UID
	newDeployment.Generation = existingDeployment.Generation

	// 执行更新
	_, err = clientset.AppsV1().Deployments(namespaceName).Update(ctx, newDeployment, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("更新Deployment失败: %v", err)
	}

	response.Success = true
	response.Message = "Deployment更新成功"
	response.UpdateStrategy = "rolling"
	response.Changes = []string{"Deployment配置已更新，正在执行滚动更新"}

	return nil
}

// updateStatefulSetYAML 更新StatefulSet YAML
func (s *K8sWorkloadServiceImpl) updateStatefulSetYAML(clientset *kubernetes.Clientset, ctx context.Context, namespaceName string, req *model.UpdateWorkloadYAMLRequest, response *model.UpdateWorkloadYAMLResponse) error {
	// 解析新的YAML
	newStatefulSet := &appsv1.StatefulSet{}
	if err := yaml.Unmarshal([]byte(req.YAMLContent), newStatefulSet); err != nil {
		return fmt.Errorf("解析StatefulSet YAML失败: %v", err)
	}

	// 验证名称匹配
	if newStatefulSet.Name != req.WorkloadName {
		return fmt.Errorf("YAML中的StatefulSet名称(%s)与请求参数不匹配(%s)", newStatefulSet.Name, req.WorkloadName)
	}

	// 确保命名空间正确
	if newStatefulSet.Namespace == "" {
		newStatefulSet.Namespace = namespaceName
	} else if newStatefulSet.Namespace != namespaceName {
		return fmt.Errorf("YAML中的命名空间(%s)与URL参数不匹配(%s)", newStatefulSet.Namespace, namespaceName)
	}

	// 获取现有的StatefulSet
	existingStatefulSet, err := clientset.AppsV1().StatefulSets(namespaceName).Get(ctx, req.WorkloadName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("StatefulSet不存在: %v", err)
	}

	// DryRun模式
	if req.DryRun {
		response.Success = true
		response.Message = "DryRun模式：仅分析变更，未实际更新"
		response.UpdateStrategy = "rolling"
		response.Changes = []string{"StatefulSet配置变更"}
		return nil
	}

	// 保留重要的系统字段
	newStatefulSet.ResourceVersion = existingStatefulSet.ResourceVersion
	newStatefulSet.UID = existingStatefulSet.UID
	newStatefulSet.Generation = existingStatefulSet.Generation

	// 执行更新
	_, err = clientset.AppsV1().StatefulSets(namespaceName).Update(ctx, newStatefulSet, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("更新StatefulSet失败: %v", err)
	}

	response.Success = true
	response.Message = "StatefulSet更新成功"
	response.UpdateStrategy = "rolling"
	response.Changes = []string{"StatefulSet配置已更新"}

	return nil
}

// updateDaemonSetYAML 更新DaemonSet YAML
func (s *K8sWorkloadServiceImpl) updateDaemonSetYAML(clientset *kubernetes.Clientset, ctx context.Context, namespaceName string, req *model.UpdateWorkloadYAMLRequest, response *model.UpdateWorkloadYAMLResponse) error {
	// 解析新的YAML
	newDaemonSet := &appsv1.DaemonSet{}
	if err := yaml.Unmarshal([]byte(req.YAMLContent), newDaemonSet); err != nil {
		return fmt.Errorf("解析DaemonSet YAML失败: %v", err)
	}

	// 验证名称匹配
	if newDaemonSet.Name != req.WorkloadName {
		return fmt.Errorf("YAML中的DaemonSet名称(%s)与请求参数不匹配(%s)", newDaemonSet.Name, req.WorkloadName)
	}

	// 确保命名空间正确
	if newDaemonSet.Namespace == "" {
		newDaemonSet.Namespace = namespaceName
	} else if newDaemonSet.Namespace != namespaceName {
		return fmt.Errorf("YAML中的命名空间(%s)与URL参数不匹配(%s)", newDaemonSet.Namespace, namespaceName)
	}

	// 获取现有的DaemonSet
	existingDaemonSet, err := clientset.AppsV1().DaemonSets(namespaceName).Get(ctx, req.WorkloadName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("DaemonSet不存在: %v", err)
	}

	// DryRun模式
	if req.DryRun {
		response.Success = true
		response.Message = "DryRun模式：仅分析变更，未实际更新"
		response.UpdateStrategy = "rolling"
		response.Changes = []string{"DaemonSet配置变更"}
		return nil
	}

	// 保留重要的系统字段
	newDaemonSet.ResourceVersion = existingDaemonSet.ResourceVersion
	newDaemonSet.UID = existingDaemonSet.UID
	newDaemonSet.Generation = existingDaemonSet.Generation

	// 执行更新
	_, err = clientset.AppsV1().DaemonSets(namespaceName).Update(ctx, newDaemonSet, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("更新DaemonSet失败: %v", err)
	}

	response.Success = true
	response.Message = "DaemonSet更新成功"
	response.UpdateStrategy = "rolling"
	response.Changes = []string{"DaemonSet配置已更新"}

	return nil
}

// updateJobYAML 更新Job YAML
func (s *K8sWorkloadServiceImpl) updateJobYAML(clientset *kubernetes.Clientset, ctx context.Context, namespaceName string, req *model.UpdateWorkloadYAMLRequest, response *model.UpdateWorkloadYAMLResponse) error {
	newJob := &batchv1.Job{}
	if err := yaml.Unmarshal([]byte(req.YAMLContent), newJob); err != nil {
		return fmt.Errorf("解析Job YAML失败: %v", err)
	}

	if newJob.Name != req.WorkloadName {
		return fmt.Errorf("YAML中的Job名称(%s)与请求参数不匹配(%s)", newJob.Name, req.WorkloadName)
	}

	if newJob.Namespace == "" {
		newJob.Namespace = namespaceName
	} else if newJob.Namespace != namespaceName {
		return fmt.Errorf("YAML中的命名空间(%s)与URL参数不匹配(%s)", newJob.Namespace, namespaceName)
	}

	existingJob, err := clientset.BatchV1().Jobs(namespaceName).Get(ctx, req.WorkloadName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("Job不存在: %v", err)
	}

	if req.DryRun {
		response.Success = true
		response.Message = "DryRun模式：仅分析变更，未实际更新"
		response.UpdateStrategy = "update"
		response.Changes = []string{"Job配置变更"}
		return nil
	}

	newJob.ResourceVersion = existingJob.ResourceVersion
	newJob.UID = existingJob.UID
	newJob.Generation = existingJob.Generation

	_, err = clientset.BatchV1().Jobs(namespaceName).Update(ctx, newJob, metav1.UpdateOptions{})
	if err == nil {
		response.Success = true
		response.Message = "Job更新成功"
		response.UpdateStrategy = "update"
		response.Changes = []string{"Job配置已更新"}
		return nil
	}

	if !s.isJobImmutableUpdateError(err) {
		return fmt.Errorf("更新Job失败: %v", err)
	}

	if !req.Force {
		return fmt.Errorf("Job包含不可变字段变更，请设置force=true后删除并重建: %v", err)
	}

	deletePolicy := metav1.DeletePropagationForeground
	if err := clientset.BatchV1().Jobs(namespaceName).Delete(ctx, req.WorkloadName, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		return fmt.Errorf("删除旧Job失败: %v", err)
	}

	if err := s.waitForJobDeletion(ctx, clientset, namespaceName, req.WorkloadName); err != nil {
		return err
	}

	newJob.ResourceVersion = ""
	newJob.UID = ""
	newJob.Generation = 0
	newJob.ManagedFields = nil
	newJob.CreationTimestamp = metav1.Time{}
	newJob.Status = batchv1.JobStatus{}

	if _, err := clientset.BatchV1().Jobs(namespaceName).Create(ctx, newJob, metav1.CreateOptions{}); err != nil {
		return fmt.Errorf("重建Job失败: %v", err)
	}

	response.Success = true
	response.Message = "Job包含不可变字段变更，已删除后重建"
	response.UpdateStrategy = "recreate"
	response.Changes = []string{"Job配置已更新", "Job已按新配置重建"}
	response.Warnings = append(response.Warnings, "Job重建会重新触发执行，请关注运行结果")
	return nil
}

func (s *K8sWorkloadServiceImpl) isJobImmutableUpdateError(err error) bool {
	if err == nil {
		return false
	}

	message := strings.ToLower(err.Error())
	return apierrors.IsInvalid(err) &&
		(strings.Contains(message, "immutable") || strings.Contains(message, "spec.template"))
}

func (s *K8sWorkloadServiceImpl) waitForJobDeletion(ctx context.Context, clientset *kubernetes.Clientset, namespaceName, jobName string) error {
	deadline := time.Now().Add(15 * time.Second)
	for {
		_, err := clientset.BatchV1().Jobs(namespaceName).Get(ctx, jobName, metav1.GetOptions{})
		if apierrors.IsNotFound(err) {
			return nil
		}
		if err != nil {
			return fmt.Errorf("等待Job删除时查询失败: %v", err)
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("等待Job删除超时: %s", jobName)
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(500 * time.Millisecond):
		}
	}
}

// updateCronJobYAML 更新CronJob YAML
func (s *K8sWorkloadServiceImpl) updateCronJobYAML(clientset *kubernetes.Clientset, ctx context.Context, namespaceName string, req *model.UpdateWorkloadYAMLRequest, response *model.UpdateWorkloadYAMLResponse) error {
	// 解析新的YAML
	newCronJob := &batchv1.CronJob{}
	if err := yaml.Unmarshal([]byte(req.YAMLContent), newCronJob); err != nil {
		return fmt.Errorf("解析CronJob YAML失败: %v", err)
	}

	// 验证名称匹配
	if newCronJob.Name != req.WorkloadName {
		return fmt.Errorf("YAML中的CronJob名称(%s)与请求参数不匹配(%s)", newCronJob.Name, req.WorkloadName)
	}

	// 确保命名空间正确
	if newCronJob.Namespace == "" {
		newCronJob.Namespace = namespaceName
	} else if newCronJob.Namespace != namespaceName {
		return fmt.Errorf("YAML中的命名空间(%s)与URL参数不匹配(%s)", newCronJob.Namespace, namespaceName)
	}

	// 获取现有的CronJob
	existingCronJob, err := clientset.BatchV1().CronJobs(namespaceName).Get(ctx, req.WorkloadName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("CronJob不存在: %v", err)
	}

	// DryRun模式
	if req.DryRun {
		response.Success = true
		response.Message = "DryRun模式：仅分析变更，未实际更新"
		response.UpdateStrategy = "update"
		response.Changes = []string{"CronJob配置变更"}
		return nil
	}

	// 保留重要的系统字段
	newCronJob.ResourceVersion = existingCronJob.ResourceVersion
	newCronJob.UID = existingCronJob.UID
	newCronJob.Generation = existingCronJob.Generation

	// 执行更新
	_, err = clientset.BatchV1().CronJobs(namespaceName).Update(ctx, newCronJob, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("更新CronJob失败: %v", err)
	}

	response.Success = true
	response.Message = "CronJob更新成功"
	response.UpdateStrategy = "update"
	response.Changes = []string{"CronJob配置已更新"}

	return nil
}

// analyzeDeploymentChanges 分析Deployment变更
func (s *K8sWorkloadServiceImpl) analyzeDeploymentChanges(existing, new *appsv1.Deployment) []string {
	changes := []string{}

	// 比较副本数
	if existing.Spec.Replicas != nil && new.Spec.Replicas != nil && *existing.Spec.Replicas != *new.Spec.Replicas {
		changes = append(changes, fmt.Sprintf("副本数变更: %d -> %d", *existing.Spec.Replicas, *new.Spec.Replicas))
	}

	// 比较容器镜像
	if len(existing.Spec.Template.Spec.Containers) > 0 && len(new.Spec.Template.Spec.Containers) > 0 {
		for i, existingContainer := range existing.Spec.Template.Spec.Containers {
			if i < len(new.Spec.Template.Spec.Containers) {
				newContainer := new.Spec.Template.Spec.Containers[i]
				if existingContainer.Image != newContainer.Image {
					changes = append(changes, fmt.Sprintf("容器镜像变更 (%s): %s -> %s", existingContainer.Name, existingContainer.Image, newContainer.Image))
				}
			}
		}
	}

	// 比较labels
	if !s.compareMaps(existing.Spec.Template.Labels, new.Spec.Template.Labels) {
		changes = append(changes, "Pod模板Labels发生变更")
	}

	if len(changes) == 0 {
		changes = append(changes, "未检测到显著变更")
	}

	return changes
}
