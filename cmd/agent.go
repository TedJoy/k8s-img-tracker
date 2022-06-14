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
	// get all deploy/sts/ds/ksvc with label img-tracker=enabled
	// filter and get all labels starts with img-tracker
	// get image name using suffix after img-tracker::
	// get location of image: initContainers::name or containers::name

	// get current sha hash of image using skopeo
	// skopeo inspect docker://443533367748.dkr.ecr.ap-southeast-1.amazonaws.com/ai/graph-matching -f '{{.Digest}}' (get stdout)

	// patch image at location to new sha
	logger.Logger.Infof("%v", config.MyEnvConfig.UseKubeCfg)
	logger.Logger.Infof("%v", config.MyFileConfig)
	k8sClientSet := k8s.GetClientSet()
	dynamicClientSet := k8s.GetDynamic()

	var wg sync.WaitGroup
	wg.Add(2)

	go core.ReconcileDeployment(k8sClientSet, time.Second, wg)
	go core.ReconcileKnativeService(dynamicClientSet, time.Second, wg)
	// go core.ReconcileKnativeService(clientSet)

	wg.Wait()
}

// digest := core.GetImgDigest("443533367748.dkr.ecr.ap-southeast-1.amazonaws.com/ai/graph-matching")

// logger.Logger.Info(digest)
