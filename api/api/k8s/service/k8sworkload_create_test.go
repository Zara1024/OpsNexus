package service

import (
	"testing"

	"dodevops-api/api/k8s/model"
	corev1 "k8s.io/api/core/v1"
)

func TestBuildDeploymentForCreateDefaultsSelectorAndReplicas(t *testing.T) {
	service := &K8sWorkloadServiceImpl{}

	deployment, err := service.buildDeploymentForCreate("default", &model.CreateDeploymentRequest{
		Name: "opsnexus-create-e2e",
		Template: model.PodTemplateSpec{
			Containers: []model.ContainerSpec{
				{
					Name:  "nginx",
					Image: "nginx:1.27",
				},
			},
		},
	})
	if err != nil {
		t.Fatalf("buildDeploymentForCreate returned error: %v", err)
	}

	if deployment.Name != "opsnexus-create-e2e" {
		t.Fatalf("unexpected deployment name: %s", deployment.Name)
	}
	if deployment.Namespace != "default" {
		t.Fatalf("unexpected namespace: %s", deployment.Namespace)
	}
	if deployment.Spec.Replicas == nil || *deployment.Spec.Replicas != 1 {
		t.Fatalf("expected default replicas to be 1, got %+v", deployment.Spec.Replicas)
	}
	if got := deployment.Spec.Selector.MatchLabels["app"]; got != "opsnexus-create-e2e" {
		t.Fatalf("expected selector app label, got %q", got)
	}
	if got := deployment.Spec.Template.Labels["app"]; got != "opsnexus-create-e2e" {
		t.Fatalf("expected template app label, got %q", got)
	}
	if len(deployment.Spec.Template.Spec.Containers) != 1 {
		t.Fatalf("expected one container, got %d", len(deployment.Spec.Template.Spec.Containers))
	}
	if deployment.Spec.Template.Spec.Containers[0].Image != "nginx:1.27" {
		t.Fatalf("unexpected image: %s", deployment.Spec.Template.Spec.Containers[0].Image)
	}
	if deployment.Spec.Strategy.Type == "" {
		t.Fatal("expected deployment strategy type to be populated")
	}
	if status := service.getDeploymentStatus(deployment); status != "Pending" {
		t.Fatalf("expected freshly created deployment status to be Pending, got %s", status)
	}
}

func TestBuildDeploymentForCreateBuildsResourcesVolumesAndStrategy(t *testing.T) {
	service := &K8sWorkloadServiceImpl{}

	deployment, err := service.buildDeploymentForCreate("opsnexus", &model.CreateDeploymentRequest{
		Name:     "opsnexus-web",
		Replicas: 3,
		Labels: map[string]string{
			"tier": "web",
		},
		Template: model.PodTemplateSpec{
			Labels: map[string]string{
				"app": "opsnexus-web",
			},
			NodeSelector: map[string]string{
				"kubernetes.io/os": "linux",
			},
			Tolerations: []model.Toleration{
				{
					Key:      "dedicated",
					Operator: "Equal",
					Value:    "ops",
					Effect:   "NoSchedule",
				},
			},
			Volumes: []model.VolumeSpec{
				{
					Name: "config",
					Type: "ConfigMap",
					Config: map[string]interface{}{
						"name": "opsnexus-config",
					},
				},
			},
			Containers: []model.ContainerSpec{
				{
					Name:  "web",
					Image: "nginx:1.27",
					Ports: []model.ContainerPort{
						{
							Name:          "http",
							ContainerPort: 80,
						},
					},
					Env: []model.EnvVar{
						{
							Name:  "APP_ENV",
							Value: "prod",
						},
					},
					Resources: model.WorkloadResources{
						Requests: model.ResourceSpec{
							CPU:    "100m",
							Memory: "128Mi",
						},
						Limits: model.ResourceSpec{
							CPU:    "200m",
							Memory: "256Mi",
						},
					},
					VolumeMounts: []model.VolumeMount{
						{
							Name:      "config",
							MountPath: "/etc/config",
						},
					},
					Command: []string{"nginx"},
					Args:    []string{"-g", "daemon off;"},
				},
			},
		},
		Strategy: model.DeploymentStrategy{
			Type: "RollingUpdate",
			RollingUpdate: model.RollingUpdateDeployment{
				MaxUnavailable: "25%",
				MaxSurge:       "1",
			},
		},
	})
	if err != nil {
		t.Fatalf("buildDeploymentForCreate returned error: %v", err)
	}

	if deployment.Spec.Replicas == nil || *deployment.Spec.Replicas != 3 {
		t.Fatalf("expected replicas to be 3, got %+v", deployment.Spec.Replicas)
	}
	if deployment.Spec.Template.Spec.NodeSelector["kubernetes.io/os"] != "linux" {
		t.Fatalf("unexpected node selector: %+v", deployment.Spec.Template.Spec.NodeSelector)
	}
	if len(deployment.Spec.Template.Spec.Tolerations) != 1 {
		t.Fatalf("expected one toleration, got %d", len(deployment.Spec.Template.Spec.Tolerations))
	}
	if len(deployment.Spec.Template.Spec.Volumes) != 1 {
		t.Fatalf("expected one volume, got %d", len(deployment.Spec.Template.Spec.Volumes))
	}
	if deployment.Spec.Template.Spec.Volumes[0].ConfigMap == nil || deployment.Spec.Template.Spec.Volumes[0].ConfigMap.Name != "opsnexus-config" {
		t.Fatalf("unexpected configmap volume: %+v", deployment.Spec.Template.Spec.Volumes[0].VolumeSource)
	}

	container := deployment.Spec.Template.Spec.Containers[0]
	if len(container.Ports) != 1 || container.Ports[0].Protocol != corev1.ProtocolTCP {
		t.Fatalf("expected default TCP port, got %+v", container.Ports)
	}
	if len(container.Env) != 1 || container.Env[0].Name != "APP_ENV" {
		t.Fatalf("unexpected env vars: %+v", container.Env)
	}
	if len(container.VolumeMounts) != 1 || container.VolumeMounts[0].MountPath != "/etc/config" {
		t.Fatalf("unexpected volume mounts: %+v", container.VolumeMounts)
	}
	if container.Resources.Requests.Cpu().MilliValue() != 100 {
		t.Fatalf("unexpected CPU request: %s", container.Resources.Requests.Cpu().String())
	}
	if container.Resources.Limits.Memory().String() != "256Mi" {
		t.Fatalf("unexpected memory limit: %s", container.Resources.Limits.Memory().String())
	}
	if deployment.Spec.Strategy.Type != "RollingUpdate" {
		t.Fatalf("unexpected strategy type: %s", deployment.Spec.Strategy.Type)
	}
	if deployment.Spec.Strategy.RollingUpdate == nil {
		t.Fatal("expected rolling update strategy to be populated")
	}
	if deployment.Spec.Strategy.RollingUpdate.MaxUnavailable == nil || deployment.Spec.Strategy.RollingUpdate.MaxUnavailable.String() != "25%" {
		t.Fatalf("unexpected maxUnavailable: %+v", deployment.Spec.Strategy.RollingUpdate.MaxUnavailable)
	}
	if deployment.Spec.Strategy.RollingUpdate.MaxSurge == nil || deployment.Spec.Strategy.RollingUpdate.MaxSurge.String() != "1" {
		t.Fatalf("unexpected maxSurge: %+v", deployment.Spec.Strategy.RollingUpdate.MaxSurge)
	}
}

func TestBuildDeploymentForCreateRejectsInvalidResourceQuantity(t *testing.T) {
	service := &K8sWorkloadServiceImpl{}

	_, err := service.buildDeploymentForCreate("default", &model.CreateDeploymentRequest{
		Name: "opsnexus-invalid",
		Template: model.PodTemplateSpec{
			Containers: []model.ContainerSpec{
				{
					Name:  "web",
					Image: "nginx:1.27",
					Resources: model.WorkloadResources{
						Requests: model.ResourceSpec{
							CPU: "bad-cpu",
						},
					},
				},
			},
		},
	})
	if err == nil {
		t.Fatal("expected invalid resource quantity to return error")
	}
}

func TestBuildStatefulSetForCreateBuildsRequiredFields(t *testing.T) {
	service := &K8sWorkloadServiceImpl{}

	statefulSet, err := service.buildStatefulSetForCreate("default", &model.CreateStatefulSetRequest{
		Name:        "opsnexus-statefulset",
		ServiceName: "opsnexus-headless",
		Replicas:    2,
		Template: model.PodTemplateSpec{
			Labels: map[string]string{
				"app": "opsnexus-statefulset",
			},
			Containers: []model.ContainerSpec{
				{
					Name:  "nginx",
					Image: "nginx:1.27",
				},
			},
		},
	})
	if err != nil {
		t.Fatalf("buildStatefulSetForCreate returned error: %v", err)
	}

	if statefulSet.Spec.ServiceName != "opsnexus-headless" {
		t.Fatalf("unexpected serviceName: %s", statefulSet.Spec.ServiceName)
	}
	if statefulSet.Spec.Replicas == nil || *statefulSet.Spec.Replicas != 2 {
		t.Fatalf("unexpected replicas: %+v", statefulSet.Spec.Replicas)
	}
	if got := statefulSet.Spec.Selector.MatchLabels["app"]; got != "opsnexus-statefulset" {
		t.Fatalf("unexpected selector label: %s", got)
	}
	if status := service.getStatefulSetStatus(statefulSet); status != "Pending" {
		t.Fatalf("expected freshly created statefulset status to be Pending, got %s", status)
	}
}

func TestBuildStatefulSetForCreateRequiresServiceName(t *testing.T) {
	service := &K8sWorkloadServiceImpl{}

	_, err := service.buildStatefulSetForCreate("default", &model.CreateStatefulSetRequest{
		Name: "opsnexus-statefulset",
		Template: model.PodTemplateSpec{
			Containers: []model.ContainerSpec{
				{
					Name:  "nginx",
					Image: "nginx:1.27",
				},
			},
		},
	})
	if err == nil {
		t.Fatal("expected empty serviceName to return error")
	}
}

func TestBuildDaemonSetForCreateBuildsRequiredFields(t *testing.T) {
	service := &K8sWorkloadServiceImpl{}

	daemonSet, err := service.buildDaemonSetForCreate("default", &model.CreateDaemonSetRequest{
		Name: "opsnexus-daemonset",
		Template: model.PodTemplateSpec{
			Labels: map[string]string{
				"app": "opsnexus-daemonset",
			},
			Containers: []model.ContainerSpec{
				{
					Name:  "nginx",
					Image: "nginx:1.27",
				},
			},
		},
	})
	if err != nil {
		t.Fatalf("buildDaemonSetForCreate returned error: %v", err)
	}

	if got := daemonSet.Spec.Selector.MatchLabels["app"]; got != "opsnexus-daemonset" {
		t.Fatalf("unexpected selector label: %s", got)
	}
	if status := service.getDaemonSetStatus(daemonSet); status != "Pending" {
		t.Fatalf("expected freshly created daemonset status to be Pending, got %s", status)
	}
}

func TestBuildJobForCreateDefaultsRestartPolicy(t *testing.T) {
	service := &K8sWorkloadServiceImpl{}

	job, err := service.buildJobForCreate("default", &model.CreateJobRequest{
		Name:        "opsnexus-job",
		Completions: 1,
		Parallelism: 1,
		Template: model.PodTemplateSpec{
			Labels: map[string]string{
				"app": "opsnexus-job",
			},
			Containers: []model.ContainerSpec{
				{
					Name:  "busybox",
					Image: "busybox:1.36",
				},
			},
		},
	})
	if err != nil {
		t.Fatalf("buildJobForCreate returned error: %v", err)
	}

	if job.Spec.Template.Spec.RestartPolicy != corev1.RestartPolicyNever {
		t.Fatalf("expected restartPolicy Never, got %s", job.Spec.Template.Spec.RestartPolicy)
	}
	if job.Spec.Completions == nil || *job.Spec.Completions != 1 {
		t.Fatalf("unexpected completions: %+v", job.Spec.Completions)
	}
}

func TestBuildCronJobForCreateBuildsScheduleAndJobTemplate(t *testing.T) {
	service := &K8sWorkloadServiceImpl{}

	cronJob, err := service.buildCronJobForCreate("default", &model.CreateCronJobRequest{
		Name:     "opsnexus-cronjob",
		Schedule: "*/5 * * * *",
		JobTemplate: model.JobTemplateSpec{
			Labels: map[string]string{
				"app": "opsnexus-cronjob",
			},
			Template: model.PodTemplateSpec{
				Labels: map[string]string{
					"app": "opsnexus-cronjob",
				},
				Containers: []model.ContainerSpec{
					{
						Name:  "busybox",
						Image: "busybox:1.36",
					},
				},
			},
		},
	})
	if err != nil {
		t.Fatalf("buildCronJobForCreate returned error: %v", err)
	}

	if cronJob.Spec.Schedule != "*/5 * * * *" {
		t.Fatalf("unexpected schedule: %s", cronJob.Spec.Schedule)
	}
	if cronJob.Spec.JobTemplate.Spec.Template.Spec.RestartPolicy != corev1.RestartPolicyNever {
		t.Fatalf("expected restartPolicy Never, got %s", cronJob.Spec.JobTemplate.Spec.Template.Spec.RestartPolicy)
	}
}
