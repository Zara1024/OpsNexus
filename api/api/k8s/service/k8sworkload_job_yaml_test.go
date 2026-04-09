package service

import (
	"errors"
	"testing"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func TestIsJobImmutableUpdateError(t *testing.T) {
	service := &K8sWorkloadServiceImpl{}

	immutableErr := apierrors.NewInvalid(
		schema.GroupKind{Group: "batch", Kind: "Job"},
		"demo-job",
		field.ErrorList{
			field.Invalid(field.NewPath("spec", "template"), "value", "field is immutable"),
		},
	)

	if !service.isJobImmutableUpdateError(immutableErr) {
		t.Fatal("expected immutable job update error to be detected")
	}

	validationErr := apierrors.NewInvalid(
		schema.GroupKind{Group: "batch", Kind: "Job"},
		"demo-job",
		field.ErrorList{
			field.Invalid(field.NewPath("metadata", "name"), "value", "must be a valid DNS label"),
		},
	)

	if service.isJobImmutableUpdateError(validationErr) {
		t.Fatal("expected non-immutable validation error to be ignored")
	}

	if service.isJobImmutableUpdateError(errors.New("plain error")) {
		t.Fatal("expected plain error to be ignored")
	}
}
