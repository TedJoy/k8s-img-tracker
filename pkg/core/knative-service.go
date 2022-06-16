package core

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"git2.gnt-global.com/jlab/gdeploy/img-tracker/config"
	"git2.gnt-global.com/jlab/gdeploy/img-tracker/pkg/logger"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
)

// func WatchKnativeService(clientset *kubernetes.Clientset) {
// 	factory := informers.NewSharedInformerFactoryWithOptions(clientset, time.Minute, informers.WithTweakListOptions(func(lo *metav1.ListOptions) {
// 		lo.LabelSelector = config.MyFileConfig.AppKey + "/deployment=true"
// 	}))

// 	deploymentInformer := factory.Apps().V1().KnativeServices()

// 	stopper := make(chan struct{})
// 	defer close(stopper)

// 	informer := deploymentInformer.Informer()
// 	defer runtime.HandleCrash()
// 	go factory.Start(stopper)

// 	if !cache.WaitForCacheSync(stopper, informer.HasSynced) {
// 		runtime.HandleError(fmt.Errorf("timed out waiting for caches to sync"))
// 		return
// 	}

// 	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
// 		AddFunc:    onAdd,
// 		UpdateFunc: onUpdate,
// 		DeleteFunc: func(interface{}) { logger.Logger.Debug("delete not implemented") },
// 	})

// 	<-stopper
// }

var knativeServiceResource = schema.GroupVersionResource{
	Group:    "serving.knative.dev",
	Version:  "v1",
	Resource: "services",
}

func ReconcileKnativeService(clientSet dynamic.Interface, d time.Duration, wg sync.WaitGroup) {
	defer wg.Done()
	for range time.Tick(d) {
		ReconcileKnativeServiceRun(clientSet)
	}
}

func ReconcileKnativeServiceRun(clientSet dynamic.Interface) {
	imageHashMap := make(map[string]string)

	objs, err := clientSet.Resource(knativeServiceResource).List(context.TODO(), metav1.ListOptions{LabelSelector: config.MyFileConfig.AppKey + "/knative-service=true"})
	if err != nil {
		if err.Error() == "the server could not find the requested resource" {
			return
		} else {
			logger.Logger.Info(err)
		}
	}

	for _, item := range objs.Items {
		// TODO: no                                   need to patch when containers/name does not exist

		itemConfigRaw := item.GetAnnotations()[config.MyFileConfig.AppKey+"/config"]
		itemConfig := &map[string]string{}
		json.Unmarshal([]byte(itemConfigRaw), &itemConfig)

		for k, v := range *itemConfig {
			splits1 := strings.Split(k, "/")
			currentImageHash, err := GetCurrentImageHashKnative(item, splits1[0], splits1[1])

			if err != nil {
				logger.Logger.Info(err)
			}

			if currentImageHash == "" {
				continue
			}
			splits2 := strings.Split(v, ":")
			upToDateHash := GetImgDigest(v, imageHashMap)
			upToDateImageHash := splits2[0] + "@" + upToDateHash

			logger.Logger.Debugf("CurrentImageHashKnative: %s", currentImageHash)

			logger.Logger.Debugf("type: %s, name: %s, image: %s, image@hash: %s, current image@hash: %s", splits1[0], splits1[1], v, upToDateImageHash, currentImageHash)

			if currentImageHash != upToDateImageHash {
				patchOp := CreatePatchOpKnative(item, splits1[0], splits1[1], upToDateImageHash)
				logger.Logger.Debug(patchOp)
				patchOpBytes, err := json.Marshal([]PatchOperation{patchOp})
				if err != nil {
					logger.Logger.Infow("error masharling patch op", "patchOp", patchOp, "patchOp", patchOp, "err", err)
				}
				_, err = clientSet.Resource(knativeServiceResource).Namespace(item.GetNamespace()).Patch(context.TODO(), item.GetName(), types.JSONPatchType, patchOpBytes, metav1.PatchOptions{})
				if err != nil {
					logger.Logger.Infow("error patching", "patchOp", patchOp, "err", err)
				}
			}
		}
	}
}

func GetCurrentImageHashKnative(item unstructured.Unstructured, containerType, containerName string) (string, error) {
	containers, found, err := unstructured.NestedSlice(item.UnstructuredContent(), "spec", "template", "spec", containerType)
	if err != nil {
		logger.Logger.Debug(err)
	}

	logger.Logger.Debugf("Containers: %s, ContainerName: %s", containers, containerName)

	if !found {
		return "", nil
	}

	for _, container := range containers {
		containerMap := container.(map[string]interface{})
		if containerMap["name"] == containerName {
			return containerMap["image"].(string), nil
		}
	}

	return "", nil
}

func CreatePatchOpKnative(item unstructured.Unstructured, containerType, containerName, upToDateImageHash string) PatchOperation {
	containers, found, err := unstructured.NestedSlice(item.UnstructuredContent(), "spec", "template", "spec", containerType)
	if err != nil {
		logger.Logger.Debug(err)
	}

	logger.Logger.Debugf("Containers: %s, ContainerName: %s", containers, containerName)

	if !found {
		return PatchOperation{}
	}

	index := -1
	for i, container := range containers {
		containerMap := container.(map[string]interface{})
		if containerMap["name"] == containerName {
			index = i
		}
	}

	return PatchOperation{
		Op:    "replace",
		Path:  fmt.Sprintf("/spec/template/spec/%s/%d/image", containerType, index),
		Value: upToDateImageHash,
	}
}
