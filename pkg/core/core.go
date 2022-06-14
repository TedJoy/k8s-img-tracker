package core

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
)

func GetCurrentImageHash(podSpec v1.PodSpec, containerType, containerName string) string {
	var containers []v1.Container = []v1.Container{}
	switch containerType {
	case "containers":
		containers = podSpec.Containers
	case "initContainers":
		containers = podSpec.InitContainers
	}

	for _, container := range containers {
		if container.Name == containerName {
			return container.Image
		}
	}

	return ""
}

func CreatePatchOp(podSpec v1.PodSpec, containerType, containerName, upToDateImageHash string) PatchOperation {
	var containers []v1.Container = []v1.Container{}
	switch containerType {
	case "containers":
		containers = podSpec.Containers
	case "initContainers":
		containers = podSpec.InitContainers
	}

	index := -1
	for i, container := range containers {
		if container.Name == containerName {
			index = i
			break
		}
	}

	return PatchOperation{
		Op:    "replace",
		Path:  fmt.Sprintf("/spec/template/spec/%s/%d/image", containerType, index),
		Value: upToDateImageHash,
	}
}
