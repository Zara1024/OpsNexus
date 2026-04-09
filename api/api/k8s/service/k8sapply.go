package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"dodevops-api/api/k8s/dao"
	"dodevops-api/api/k8s/model"
	"dodevops-api/common/result"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	meta "k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	utilyaml "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/discovery"
	cacheddiscovery "k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
)

const defaultApplyFieldManager = "opsnexus"

type K8sApplyService struct {
	clusterDao *dao.KubeClusterDao
}

type manifestResourceClient interface {
	Get(ctx context.Context, name string, opts metav1.GetOptions, subresources ...string) (*unstructured.Unstructured, error)
	Create(ctx context.Context, obj *unstructured.Unstructured, opts metav1.CreateOptions, subresources ...string) (*unstructured.Unstructured, error)
	Update(ctx context.Context, obj *unstructured.Unstructured, opts metav1.UpdateOptions, subresources ...string) (*unstructured.Unstructured, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (*unstructured.Unstructured, error)
}

type applyManifestOptions struct {
	DryRun          bool
	ServerSideApply bool
	FieldManager    string
	ForceConflicts  bool
}

func NewK8sApplyService(db *gorm.DB) *K8sApplyService {
	return &K8sApplyService{
		clusterDao: dao.NewKubeClusterDao(db),
	}
}

func (s *K8sApplyService) ApplyManifest(c *gin.Context, clusterID uint, req *model.ApplyManifestRequest) {
	if strings.TrimSpace(req.YAMLContent) == "" {
		result.Failed(c, http.StatusBadRequest, "YAML内容不能为空")
		return
	}

	cluster, err := s.clusterDao.GetByID(clusterID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result.Failed(c, http.StatusNotFound, "集群不存在")
			return
		}
		result.Failed(c, http.StatusInternalServerError, "获取集群信息失败: "+err.Error())
		return
	}

	manifests, err := decodeManifestDocuments(req.YAMLContent)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "YAML解析失败: "+err.Error())
		return
	}
	if len(manifests) == 0 {
		result.Failed(c, http.StatusBadRequest, "YAML中没有可应用的资源对象")
		return
	}

	dynamicClient, mapper, err := s.createDynamicClient(cluster.Credential)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "连接K8s集群失败: "+err.Error())
		return
	}

	serverSideApply := isServerSideApplyEnabled(req)
	applyOpts := applyManifestOptions{
		DryRun:          req.DryRun,
		ServerSideApply: serverSideApply,
		FieldManager:    strings.TrimSpace(req.FieldManager),
		ForceConflicts:  req.ForceConflicts,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	results := make([]model.ApplyManifestItemResult, 0, len(manifests))
	for index, manifest := range manifests {
		cleanManifest := sanitizeManifestForApply(manifest)
		if err := validateManifestBasics(cleanManifest); err != nil {
			result.Failed(c, http.StatusBadRequest, fmt.Sprintf("第 %d 个资源无效: %v", index+1, err))
			return
		}

		gvk := schema.FromAPIVersionAndKind(cleanManifest.GetAPIVersion(), cleanManifest.GetKind())
		mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
		if err != nil {
			result.Failed(c, http.StatusBadRequest, fmt.Sprintf("第 %d 个资源无法映射到目标集群: %v", index+1, err))
			return
		}

		isNamespaced := mapping.Scope.Name() == meta.RESTScopeNameNamespace
		if err := prepareManifestNamespace(cleanManifest, isNamespaced, req.Namespace); err != nil {
			result.Failed(c, http.StatusBadRequest, fmt.Sprintf("第 %d 个资源命名空间校验失败: %v", index+1, err))
			return
		}

		if req.ValidateOnly {
			results = append(results, buildValidatedResult(cleanManifest))
			continue
		}

		resourceClient := dynamicClient.Resource(mapping.Resource)
		var applyClient manifestResourceClient = resourceClient
		if isNamespaced {
			applyClient = resourceClient.Namespace(cleanManifest.GetNamespace())
		}

		itemResult, err := applyManifestResource(ctx, applyClient, cleanManifest, applyOpts)
		if err != nil {
			result.Failed(c, statusCodeFromK8sError(err), fmt.Sprintf("应用第 %d 个资源失败: %v", index+1, err))
			return
		}
		results = append(results, *itemResult)
	}

	result.Success(c, model.ApplyManifestResponse{
		Success:         true,
		Message:         buildApplyManifestMessage(req, len(results)),
		DryRun:          req.DryRun,
		ValidateOnly:    req.ValidateOnly,
		ServerSideApply: serverSideApply,
		Results:         results,
	})
}

func (s *K8sApplyService) createDynamicClient(kubeconfig string) (dynamic.Interface, meta.RESTMapper, error) {
	config, err := clientcmd.RESTConfigFromKubeConfig([]byte(kubeconfig))
	if err != nil {
		return nil, nil, err
	}

	discoveryClient, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		return nil, nil, err
	}

	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, nil, err
	}

	mapper := restmapper.NewDeferredDiscoveryRESTMapper(cacheddiscovery.NewMemCacheClient(discoveryClient))
	return dynamicClient, mapper, nil
}

func decodeManifestDocuments(yamlContent string) ([]*unstructured.Unstructured, error) {
	decoder := utilyaml.NewYAMLOrJSONDecoder(bytes.NewReader([]byte(yamlContent)), 4096)
	manifests := make([]*unstructured.Unstructured, 0)

	for {
		var raw map[string]interface{}
		if err := decoder.Decode(&raw); err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		if len(raw) == 0 {
			continue
		}

		items, err := expandManifestObject(&unstructured.Unstructured{Object: raw})
		if err != nil {
			return nil, err
		}
		manifests = append(manifests, items...)
	}

	return manifests, nil
}

func expandManifestObject(obj *unstructured.Unstructured) ([]*unstructured.Unstructured, error) {
	items, found, err := unstructured.NestedSlice(obj.Object, "items")
	if err != nil {
		return nil, err
	}
	if !found || !strings.HasSuffix(obj.GetKind(), "List") {
		return []*unstructured.Unstructured{obj}, nil
	}

	results := make([]*unstructured.Unstructured, 0, len(items))
	for _, item := range items {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("List 中包含无法识别的资源条目")
		}
		if _, exists := itemMap["apiVersion"]; !exists {
			itemMap["apiVersion"] = obj.GetAPIVersion()
		}
		results = append(results, &unstructured.Unstructured{Object: itemMap})
	}

	return results, nil
}

func sanitizeManifestForApply(manifest *unstructured.Unstructured) *unstructured.Unstructured {
	clean := manifest.DeepCopy()
	unstructured.RemoveNestedField(clean.Object, "metadata", "creationTimestamp")
	unstructured.RemoveNestedField(clean.Object, "metadata", "deletionTimestamp")
	unstructured.RemoveNestedField(clean.Object, "metadata", "generation")
	unstructured.RemoveNestedField(clean.Object, "metadata", "managedFields")
	unstructured.RemoveNestedField(clean.Object, "metadata", "resourceVersion")
	unstructured.RemoveNestedField(clean.Object, "metadata", "selfLink")
	unstructured.RemoveNestedField(clean.Object, "metadata", "uid")
	unstructured.RemoveNestedField(clean.Object, "status")
	return clean
}

func validateManifestBasics(manifest *unstructured.Unstructured) error {
	if manifest.GetAPIVersion() == "" {
		return fmt.Errorf("缺少 apiVersion")
	}
	if manifest.GetKind() == "" {
		return fmt.Errorf("缺少 kind")
	}
	if manifest.GetName() == "" {
		return fmt.Errorf("%s 缺少 metadata.name", manifest.GetKind())
	}
	return nil
}

func prepareManifestNamespace(manifest *unstructured.Unstructured, isNamespaced bool, defaultNamespace string) error {
	if !isNamespaced {
		if manifest.GetNamespace() != "" {
			return fmt.Errorf("%s 是集群级资源，不应指定命名空间", manifest.GetKind())
		}
		return nil
	}

	if manifest.GetNamespace() != "" {
		return nil
	}

	namespace := strings.TrimSpace(defaultNamespace)
	if namespace == "" {
		namespace = "default"
	}
	manifest.SetNamespace(namespace)
	return nil
}

func applyManifestResource(ctx context.Context, client manifestResourceClient, manifest *unstructured.Unstructured, opts applyManifestOptions) (*model.ApplyManifestItemResult, error) {
	existing, err := client.Get(ctx, manifest.GetName(), metav1.GetOptions{})
	exists := err == nil
	if err != nil && !apierrors.IsNotFound(err) {
		return nil, err
	}

	itemResult := &model.ApplyManifestItemResult{
		APIVersion: manifest.GetAPIVersion(),
		Kind:       manifest.GetKind(),
		Name:       manifest.GetName(),
		Namespace:  manifest.GetNamespace(),
	}

	if opts.ServerSideApply {
		fieldManager := opts.FieldManager
		if fieldManager == "" {
			fieldManager = defaultApplyFieldManager
		}

		patchOptions := metav1.PatchOptions{
			FieldManager: fieldManager,
			DryRun:       dryRunOption(opts.DryRun),
		}
		if opts.ForceConflicts {
			force := true
			patchOptions.Force = &force
		}

		payload, err := json.Marshal(manifest.Object)
		if err != nil {
			return nil, err
		}

		applied, err := client.Patch(ctx, manifest.GetName(), types.ApplyPatchType, payload, patchOptions)
		if err != nil {
			return nil, err
		}

		if exists {
			itemResult.Operation = "configured"
		} else {
			itemResult.Operation = "created"
		}
		itemResult.ResourceVersion = applied.GetResourceVersion()
		return itemResult, nil
	}

	if exists {
		manifest.SetResourceVersion(existing.GetResourceVersion())
		updated, err := client.Update(ctx, manifest, metav1.UpdateOptions{DryRun: dryRunOption(opts.DryRun)})
		if err != nil {
			return nil, err
		}
		itemResult.Operation = "updated"
		itemResult.ResourceVersion = updated.GetResourceVersion()
		return itemResult, nil
	}

	created, err := client.Create(ctx, manifest, metav1.CreateOptions{DryRun: dryRunOption(opts.DryRun)})
	if err != nil {
		return nil, err
	}
	itemResult.Operation = "created"
	itemResult.ResourceVersion = created.GetResourceVersion()
	return itemResult, nil
}

func buildValidatedResult(manifest *unstructured.Unstructured) model.ApplyManifestItemResult {
	return model.ApplyManifestItemResult{
		APIVersion: manifest.GetAPIVersion(),
		Kind:       manifest.GetKind(),
		Name:       manifest.GetName(),
		Namespace:  manifest.GetNamespace(),
		Operation:  "validated",
	}
}

func buildApplyManifestMessage(req *model.ApplyManifestRequest, count int) string {
	switch {
	case req.ValidateOnly:
		return fmt.Sprintf("YAML 校验通过，共识别 %d 个资源", count)
	case req.DryRun:
		return fmt.Sprintf("DryRun 校验通过，共模拟应用 %d 个资源", count)
	case isServerSideApplyEnabled(req):
		return fmt.Sprintf("Server-Side Apply 成功，共应用 %d 个资源", count)
	default:
		return fmt.Sprintf("资源应用成功，共处理 %d 个资源", count)
	}
}

func isServerSideApplyEnabled(req *model.ApplyManifestRequest) bool {
	if req.ServerSideApply == nil {
		return false
	}
	return *req.ServerSideApply
}

func dryRunOption(enabled bool) []string {
	if !enabled {
		return nil
	}
	return []string{metav1.DryRunAll}
}

func statusCodeFromK8sError(err error) int {
	switch {
	case apierrors.IsBadRequest(err), apierrors.IsInvalid(err):
		return http.StatusBadRequest
	case apierrors.IsUnauthorized(err):
		return http.StatusUnauthorized
	case apierrors.IsForbidden(err):
		return http.StatusForbidden
	case apierrors.IsNotFound(err):
		return http.StatusNotFound
	case apierrors.IsConflict(err), apierrors.IsAlreadyExists(err):
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
