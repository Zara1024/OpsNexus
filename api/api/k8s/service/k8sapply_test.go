package service

import (
	"context"
	"encoding/json"
	"testing"

	"dodevops-api/api/k8s/model"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
)

type stubManifestResourceClient struct {
	getObject    *unstructured.Unstructured
	getErr       error
	createObject *unstructured.Unstructured
	createErr    error
	updateObject *unstructured.Unstructured
	updateErr    error
	patchObject  *unstructured.Unstructured
	patchErr     error

	createCalled bool
	updateCalled bool
	patchCalled  bool

	createdInput  *unstructured.Unstructured
	updatedInput  *unstructured.Unstructured
	patchType     types.PatchType
	patchData     []byte
	createOptions metav1.CreateOptions
	updateOptions metav1.UpdateOptions
	patchOptions  metav1.PatchOptions
}

func (s *stubManifestResourceClient) Get(context.Context, string, metav1.GetOptions, ...string) (*unstructured.Unstructured, error) {
	return s.getObject, s.getErr
}

func (s *stubManifestResourceClient) Create(_ context.Context, obj *unstructured.Unstructured, opts metav1.CreateOptions, _ ...string) (*unstructured.Unstructured, error) {
	s.createCalled = true
	s.createdInput = obj.DeepCopy()
	s.createOptions = opts
	return s.createObject, s.createErr
}

func (s *stubManifestResourceClient) Update(_ context.Context, obj *unstructured.Unstructured, opts metav1.UpdateOptions, _ ...string) (*unstructured.Unstructured, error) {
	s.updateCalled = true
	s.updatedInput = obj.DeepCopy()
	s.updateOptions = opts
	return s.updateObject, s.updateErr
}

func (s *stubManifestResourceClient) Patch(_ context.Context, _ string, pt types.PatchType, data []byte, opts metav1.PatchOptions, _ ...string) (*unstructured.Unstructured, error) {
	s.patchCalled = true
	s.patchType = pt
	s.patchData = append([]byte(nil), data...)
	s.patchOptions = opts
	return s.patchObject, s.patchErr
}

func TestDecodeManifestDocumentsMultiDoc(t *testing.T) {
	yamlContent := `
apiVersion: v1
kind: ConfigMap
metadata:
  name: cm-a
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: deploy-a
  namespace: demo
`

	manifests, err := decodeManifestDocuments(yamlContent)
	if err != nil {
		t.Fatalf("decodeManifestDocuments returned error: %v", err)
	}
	if len(manifests) != 2 {
		t.Fatalf("expected 2 manifests, got %d", len(manifests))
	}
	if manifests[0].GetKind() != "ConfigMap" || manifests[0].GetName() != "cm-a" {
		t.Fatalf("unexpected first manifest: %s/%s", manifests[0].GetKind(), manifests[0].GetName())
	}
	if manifests[1].GetKind() != "Deployment" || manifests[1].GetNamespace() != "demo" {
		t.Fatalf("unexpected second manifest namespace: %s", manifests[1].GetNamespace())
	}
}

func TestPrepareManifestNamespace(t *testing.T) {
	obj := &unstructured.Unstructured{}
	obj.SetAPIVersion("apps/v1")
	obj.SetKind("Deployment")
	obj.SetName("demo")

	if err := prepareManifestNamespace(obj, true, "ns-a"); err != nil {
		t.Fatalf("prepareManifestNamespace returned error: %v", err)
	}
	if obj.GetNamespace() != "ns-a" {
		t.Fatalf("expected namespace ns-a, got %s", obj.GetNamespace())
	}

	clusterScoped := &unstructured.Unstructured{}
	clusterScoped.SetAPIVersion("v1")
	clusterScoped.SetKind("Namespace")
	clusterScoped.SetName("demo")
	clusterScoped.SetNamespace("unexpected")
	if err := prepareManifestNamespace(clusterScoped, false, ""); err == nil {
		t.Fatal("expected cluster-scoped resource namespace validation to fail")
	}
}

func TestApplyManifestResourceUsesServerSideApply(t *testing.T) {
	client := &stubManifestResourceClient{
		getObject: &unstructured.Unstructured{Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "ConfigMap",
			"metadata": map[string]interface{}{
				"name":            "cm-a",
				"namespace":       "demo",
				"resourceVersion": "7",
			},
		}},
		patchObject: &unstructured.Unstructured{Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "ConfigMap",
			"metadata": map[string]interface{}{
				"name":            "cm-a",
				"namespace":       "demo",
				"resourceVersion": "8",
			},
		}},
	}

	obj := &unstructured.Unstructured{}
	obj.SetAPIVersion("v1")
	obj.SetKind("ConfigMap")
	obj.SetName("cm-a")
	obj.SetNamespace("demo")

	resultItem, err := applyManifestResource(context.Background(), client, obj, applyManifestOptions{
		DryRun:          true,
		ServerSideApply: true,
		ForceConflicts:  true,
	})
	if err != nil {
		t.Fatalf("applyManifestResource returned error: %v", err)
	}

	if !client.patchCalled {
		t.Fatal("expected Patch to be called")
	}
	if client.patchType != types.ApplyPatchType {
		t.Fatalf("expected ApplyPatchType, got %s", client.patchType)
	}
	if client.patchOptions.FieldManager != defaultApplyFieldManager {
		t.Fatalf("expected default field manager %s, got %s", defaultApplyFieldManager, client.patchOptions.FieldManager)
	}
	if len(client.patchOptions.DryRun) != 1 || client.patchOptions.DryRun[0] != metav1.DryRunAll {
		t.Fatalf("expected dry-run option %v, got %v", metav1.DryRunAll, client.patchOptions.DryRun)
	}
	if client.patchOptions.Force == nil || !*client.patchOptions.Force {
		t.Fatal("expected Force=true in patch options")
	}
	if resultItem.Operation != "configured" {
		t.Fatalf("expected configured operation, got %s", resultItem.Operation)
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(client.patchData, &payload); err != nil {
		t.Fatalf("patch payload is not valid JSON: %v", err)
	}
	if payload["kind"] != "ConfigMap" {
		t.Fatalf("unexpected patch payload kind: %v", payload["kind"])
	}
}

func TestApplyManifestResourceUpdatesExistingOnClientSide(t *testing.T) {
	client := &stubManifestResourceClient{
		getObject: &unstructured.Unstructured{Object: map[string]interface{}{
			"apiVersion": "apps/v1",
			"kind":       "Deployment",
			"metadata": map[string]interface{}{
				"name":            "deploy-a",
				"namespace":       "demo",
				"resourceVersion": "3",
			},
		}},
		updateObject: &unstructured.Unstructured{Object: map[string]interface{}{
			"apiVersion": "apps/v1",
			"kind":       "Deployment",
			"metadata": map[string]interface{}{
				"name":            "deploy-a",
				"namespace":       "demo",
				"resourceVersion": "4",
			},
		}},
	}

	obj := &unstructured.Unstructured{}
	obj.SetAPIVersion("apps/v1")
	obj.SetKind("Deployment")
	obj.SetName("deploy-a")
	obj.SetNamespace("demo")

	resultItem, err := applyManifestResource(context.Background(), client, obj, applyManifestOptions{})
	if err != nil {
		t.Fatalf("applyManifestResource returned error: %v", err)
	}

	if !client.updateCalled {
		t.Fatal("expected Update to be called")
	}
	if client.createCalled {
		t.Fatal("did not expect Create to be called")
	}
	if got := client.updatedInput.GetResourceVersion(); got != "3" {
		t.Fatalf("expected update input resourceVersion to be 3, got %s", got)
	}
	if resultItem.Operation != "updated" {
		t.Fatalf("expected updated operation, got %s", resultItem.Operation)
	}
}

func TestApplyManifestResourceCreatesWhenMissing(t *testing.T) {
	client := &stubManifestResourceClient{
		getErr: apierrors.NewNotFound(schema.GroupResource{Resource: "configmaps"}, "cm-b"),
		createObject: &unstructured.Unstructured{Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "ConfigMap",
			"metadata": map[string]interface{}{
				"name":            "cm-b",
				"namespace":       "demo",
				"resourceVersion": "1",
			},
		}},
	}

	obj := &unstructured.Unstructured{}
	obj.SetAPIVersion("v1")
	obj.SetKind("ConfigMap")
	obj.SetName("cm-b")
	obj.SetNamespace("demo")

	resultItem, err := applyManifestResource(context.Background(), client, obj, applyManifestOptions{
		DryRun: true,
	})
	if err != nil {
		t.Fatalf("applyManifestResource returned error: %v", err)
	}

	if !client.createCalled {
		t.Fatal("expected Create to be called")
	}
	if len(client.createOptions.DryRun) != 1 || client.createOptions.DryRun[0] != metav1.DryRunAll {
		t.Fatalf("expected create dry-run option, got %v", client.createOptions.DryRun)
	}
	if resultItem.Operation != "created" {
		t.Fatalf("expected created operation, got %s", resultItem.Operation)
	}
}

func TestBuildApplyManifestMessage(t *testing.T) {
	serverSideApply := true
	req := &model.ApplyManifestRequest{ServerSideApply: &serverSideApply}
	if got := buildApplyManifestMessage(req, 2); got == "" {
		t.Fatal("expected message to be non-empty")
	}
}

func TestIsServerSideApplyEnabledDefaultsToFalse(t *testing.T) {
	req := &model.ApplyManifestRequest{}
	if isServerSideApplyEnabled(req) {
		t.Fatal("expected server-side apply to be disabled by default")
	}

	enabled := true
	req.ServerSideApply = &enabled
	if !isServerSideApplyEnabled(req) {
		t.Fatal("expected explicit server-side apply to be enabled")
	}
}
