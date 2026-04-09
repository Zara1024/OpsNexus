package service

import (
	"testing"

	"dodevops-api/api/k8s/model"
	autoscalingv2 "k8s.io/api/autoscaling/v2"
)

func TestEnrichAutoscalingSummaryWarnings(t *testing.T) {
	svc := &K8sWorkloadServiceImpl{}
	summary := &model.WorkloadAutoscalingSummary{
		Enabled:     true,
		MinReplicas: 1,
		MaxReplicas: 1,
		Metrics: []model.WorkloadAutoscalingMetricStatus{
			{
				ResourceName: "cpu",
				TargetType:   string(autoscalingv2.UtilizationMetricType),
				TargetValue:  "70",
			},
		},
		Behavior: &model.WorkloadAutoscalingBehavior{
			ScaleUp: &model.WorkloadAutoscalingBehaviorRule{
				SelectPolicy: "Disabled",
			},
		},
	}

	svc.enrichAutoscalingSummaryWarnings(summary, model.WorkloadResources{})

	if len(summary.Warnings) != 3 {
		t.Fatalf("expected 3 warnings, got %d: %#v", len(summary.Warnings), summary.Warnings)
	}
}

func TestValidateUpsertWorkloadAutoscalingRequestRejectsInvalidPolicyValue(t *testing.T) {
	svc := &K8sWorkloadServiceImpl{}
	req := &model.UpsertWorkloadAutoscalingRequest{
		MinReplicas: 1,
		MaxReplicas: 3,
		Metrics: []model.WorkloadAutoscalingMetricSpec{
			{
				ResourceName: "cpu",
				TargetType:   "AverageValue",
				TargetValue:  "100m",
			},
		},
		Behavior: &model.WorkloadAutoscalingBehavior{
			ScaleUp: &model.WorkloadAutoscalingBehaviorRule{
				Policies: []model.WorkloadAutoscalingScalingPolicy{
					{
						Type:          "Pods",
						Value:         0,
						PeriodSeconds: 15,
					},
				},
			},
		},
	}

	if err := svc.validateUpsertWorkloadAutoscalingRequest(req); err == nil {
		t.Fatalf("expected validation error for invalid policy value")
	}
}

func TestValidateUpsertWorkloadAutoscalingRequestAcceptsValidBehaviorPolicies(t *testing.T) {
	svc := &K8sWorkloadServiceImpl{}
	req := &model.UpsertWorkloadAutoscalingRequest{
		MinReplicas: 1,
		MaxReplicas: 4,
		Metrics: []model.WorkloadAutoscalingMetricSpec{
			{
				ResourceName: "cpu",
				TargetType:   "AverageValue",
				TargetValue:  "120m",
			},
			{
				ResourceName: "memory",
				TargetType:   "AverageValue",
				TargetValue:  "256Mi",
			},
		},
		Behavior: &model.WorkloadAutoscalingBehavior{
			ScaleUp: &model.WorkloadAutoscalingBehaviorRule{
				SelectPolicy: "Min",
				Policies: []model.WorkloadAutoscalingScalingPolicy{
					{Type: "Pods", Value: 2, PeriodSeconds: 15},
					{Type: "Percent", Value: 50, PeriodSeconds: 15},
				},
			},
			ScaleDown: &model.WorkloadAutoscalingBehaviorRule{
				SelectPolicy: "Max",
				Policies: []model.WorkloadAutoscalingScalingPolicy{
					{Type: "Pods", Value: 1, PeriodSeconds: 15},
				},
			},
		},
	}

	if err := svc.validateUpsertWorkloadAutoscalingRequest(req); err != nil {
		t.Fatalf("expected valid request, got error: %v", err)
	}
}
