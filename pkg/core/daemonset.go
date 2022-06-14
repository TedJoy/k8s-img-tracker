package core

import (
	"fmt"
	"time"

	"git2.gnt-global.com/jlab/gdeploy/img-tracker/config"
	"git2.gnt-global.com/jlab/gdeploy/img-tracker/pkg/logger"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

func WatchDaemonSet(clientset *kubernetes.Clientset) {
	factory := informers.NewSharedInformerFactoryWithOptions(clientset, time.Minute, informers.WithTweakListOptions(func(lo *metav1.ListOptions) {
		lo.LabelSelector = config.MyFileConfig.AppKey + "/daemonset=true"
	}))

	daemonsetInformer := factory.Apps().V1().DaemonSets()

	stopper := make(chan struct{})
	defer close(stopper)

	informer := daemonsetInformer.Informer()
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

func ReconcileDaemonSet() {

}
