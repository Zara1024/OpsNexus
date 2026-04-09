package service

import (
	"fmt"

	"dodevops-api/api/k8s/model"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func (s *K8sWorkloadServiceImpl) applyWorkloadUpdateRequest(workloadLabels *map[string]string, podTemplate *corev1.PodTemplateSpec, req *model.UpdateWorkloadRequest) error {
	if req == nil {
		return nil
	}

	if req.Labels != nil && len(req.Labels) > 0 {
		if *workloadLabels == nil {
			*workloadLabels = make(map[string]string)
		}
		for k, v := range req.Labels {
			(*workloadLabels)[k] = v
		}
	}

	if req.Template.Labels != nil && len(req.Template.Labels) > 0 {
		if podTemplate.Labels == nil {
			podTemplate.Labels = make(map[string]string)
		}
		for k, v := range req.Template.Labels {
			podTemplate.Labels[k] = v
		}
	}

	if req.Template.Containers != nil {
		updatedContainers, err := s.buildUpdatedContainersFromRequest(podTemplate.Spec.Containers, req.Template.Containers)
		if err != nil {
			return err
		}
		podTemplate.Spec.Containers = updatedContainers
	}

	if req.Template.Volumes != nil {
		podTemplate.Spec.Volumes = make([]corev1.Volume, 0, len(req.Template.Volumes))
		for _, volumeSpec := range req.Template.Volumes {
			volume := corev1.Volume{Name: volumeSpec.Name}

			switch volumeSpec.Type {
			case "EmptyDir":
				volume.VolumeSource = corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{}}
			case "HostPath":
				if config := volumeSpec.Config; config != nil {
					if path, ok := config["path"].(string); ok {
						volume.VolumeSource = corev1.VolumeSource{HostPath: &corev1.HostPathVolumeSource{Path: path}}
					}
				}
			case "ConfigMap":
				if config := volumeSpec.Config; config != nil {
					if name, ok := config["name"].(string); ok {
						volume.VolumeSource = corev1.VolumeSource{ConfigMap: &corev1.ConfigMapVolumeSource{
							LocalObjectReference: corev1.LocalObjectReference{Name: name},
						}}
					}
				}
			case "Secret":
				if config := volumeSpec.Config; config != nil {
					if secretName, ok := config["secretName"].(string); ok {
						volume.VolumeSource = corev1.VolumeSource{Secret: &corev1.SecretVolumeSource{SecretName: secretName}}
					}
				}
			case "PersistentVolumeClaim":
				if config := volumeSpec.Config; config != nil {
					if claimName, ok := config["claimName"].(string); ok {
						volume.VolumeSource = corev1.VolumeSource{PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
							ClaimName: claimName,
						}}
					}
				}
			}

			podTemplate.Spec.Volumes = append(podTemplate.Spec.Volumes, volume)
		}
	}

	if req.Template.NodeSelector != nil {
		podTemplate.Spec.NodeSelector = req.Template.NodeSelector
	}

	if req.Template.Tolerations != nil {
		podTemplate.Spec.Tolerations = make([]corev1.Toleration, 0, len(req.Template.Tolerations))
		for _, tolerationSpec := range req.Template.Tolerations {
			podTemplate.Spec.Tolerations = append(podTemplate.Spec.Tolerations, corev1.Toleration{
				Key:      tolerationSpec.Key,
				Operator: corev1.TolerationOperator(tolerationSpec.Operator),
				Value:    tolerationSpec.Value,
				Effect:   corev1.TaintEffect(tolerationSpec.Effect),
			})
		}
	}

	return nil
}

func (s *K8sWorkloadServiceImpl) buildUpdatedContainersFromRequest(existing []corev1.Container, specs []model.ContainerSpec) ([]corev1.Container, error) {
	existingContainers := make(map[string]corev1.Container, len(existing))
	for _, container := range existing {
		existingContainers[container.Name] = container
	}

	updatedContainers := make([]corev1.Container, 0, len(specs))
	for _, containerSpec := range specs {
		var container corev1.Container
		if existingContainer, exists := existingContainers[containerSpec.Name]; exists {
			container = existingContainer
		} else {
			container = corev1.Container{Name: containerSpec.Name}
		}

		if containerSpec.Image != "" {
			container.Image = containerSpec.Image
		}
		if container.Image == "" {
			return nil, fmt.Errorf("container %s missing image", containerSpec.Name)
		}

		if len(containerSpec.Ports) > 0 {
			container.Ports = make([]corev1.ContainerPort, 0, len(containerSpec.Ports))
			for _, portSpec := range containerSpec.Ports {
				protocol := corev1.Protocol(portSpec.Protocol)
				if protocol == "" {
					protocol = corev1.ProtocolTCP
				}
				container.Ports = append(container.Ports, corev1.ContainerPort{
					Name:          portSpec.Name,
					ContainerPort: portSpec.ContainerPort,
					Protocol:      protocol,
				})
			}
		}

		if containerSpec.Env != nil {
			container.Env = make([]corev1.EnvVar, 0, len(containerSpec.Env))
			for _, envVar := range containerSpec.Env {
				container.Env = append(container.Env, corev1.EnvVar{
					Name:  envVar.Name,
					Value: envVar.Value,
				})
			}
		}

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

		if containerSpec.VolumeMounts != nil {
			container.VolumeMounts = make([]corev1.VolumeMount, 0, len(containerSpec.VolumeMounts))
			for _, mountSpec := range containerSpec.VolumeMounts {
				container.VolumeMounts = append(container.VolumeMounts, corev1.VolumeMount{
					Name:      mountSpec.Name,
					MountPath: mountSpec.MountPath,
					ReadOnly:  mountSpec.ReadOnly,
				})
			}
		}

		if containerSpec.Command != nil {
			container.Command = containerSpec.Command
		}
		if containerSpec.Args != nil {
			container.Args = containerSpec.Args
		}

		updatedContainers = append(updatedContainers, container)
	}

	return updatedContainers, nil
}
