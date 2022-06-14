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
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

func WatchDeployment(clientset *kubernetes.Clientset) {
	factory := informers.NewSharedInformerFactoryWithOptions(clientset, time.Minute, informers.WithTweakListOptions(func(lo *metav1.ListOptions) {
		lo.LabelSelector = config.MyFileConfig.AppKey + "/deployment=true"
	}))

	deploymentInformer := factory.Apps().V1().Deployments()

	stopper := make(chan struct{})
	defer close(stopper)

	informer := deploymentInformer.Informer()
	defer runtime.HandleCrash()
	go factory.Start(stopper)

	if !cache.WaitForCacheSync(stopper, informer.HasSynced) {
		runtime.HandleError(fmt.Errorf("timed out waiting for caches to sync"))
		return
	}

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    onAdd,
		UpdateFunc: onUpdate,
		DeleteFunc: func(interface{}) { logger.Logger.Debug("delete not implemented") },
	})

	<-stopper
}

func ReconcileDeployment(clientSet kubernetes.Interface, d time.Duration, wg sync.WaitGroup) {
	defer wg.Done()
	for range time.Tick(d) {
		ReconcileDeploymentRun(clientSet)
	}
}

func ReconcileDeploymentRun(clientSet kubernetes.Interface) {
	objs, err := clientSet.AppsV1().Deployments("").List(context.TODO(), metav1.ListOptions{LabelSelector: config.MyFileConfig.AppKey + "/deployment=true"})
	if err != nil {
		panic(err.Error())
	}

	for _, item := range objs.Items {
		itemConfigRaw := item.Annotations[config.MyFileConfig.AppKey+"/config"]
		itemConfig := &map[string]string{}
		json.Unmarshal([]byte(itemConfigRaw), &itemConfig)
		for k, v := range *itemConfig {
			splits1 := strings.Split(k, "/")
			splits2 := strings.Split(v, ":")
			upToDateHash := GetImgDigest(v)
			upToDateImageHash := splits2[0] + "@" + upToDateHash
			currentImageHash := GetCurrentImageHash(item.Spec.Template.Spec, splits1[0], splits1[1])

			// logger.Logger.Infof("type: %s, name: %s, image: %s, image@hash: %s, current image@hash: %s", splits1[0], splits1[1], v, upToDateImageHash, currentImageHash)

			if currentImageHash != upToDateImageHash {
				patchOp := CreatePatchOp(item.Spec.Template.Spec, splits1[0], splits1[1], upToDateImageHash)
				logger.Logger.Debug(patchOp)
				patchOpBytes, err := json.Marshal([]PatchOperation{patchOp})
				if err != nil {
					logger.Logger.Panicw("error masharling patch op", "patchOp", patchOp, "patchOp", patchOp, "err", err)
				}
				_, err = clientSet.AppsV1().Deployments(item.Namespace).Patch(context.TODO(), item.Name, types.JSONPatchType, patchOpBytes, metav1.PatchOptions{})
				if err != nil {
					logger.Logger.Panicw("error patching", "patchOp", patchOp, "err", err)
				}
			}
		}
	}
}
