package main

import (
	"sync"
	"time"

	"git2.gnt-global.com/jlab/gdeploy/img-tracker/config"
	"git2.gnt-global.com/jlab/gdeploy/img-tracker/pkg/core"
	"git2.gnt-global.com/jlab/gdeploy/img-tracker/pkg/k8s"
	"git2.gnt-global.com/jlab/gdeploy/img-tracker/pkg/logger"
)

func init() {
	logger.Logger = logger.New("main", config.MyEnvConfig.Debug)
}

func main() {
	k8sClientSet := k8s.GetClientSet()
	dynamicClientSet := k8s.GetDynamic()

	var wg sync.WaitGroup
	wg.Add(4)

	go core.ReconcileDeployment(k8sClientSet, time.Second, wg)
	go core.ReconcileDaemonSet(k8sClientSet, time.Second, wg)
	go core.ReconcileStatefulSet(k8sClientSet, time.Second, wg)
	go core.ReconcileKnativeService(dynamicClientSet, time.Second, wg)

	wg.Wait()
}
