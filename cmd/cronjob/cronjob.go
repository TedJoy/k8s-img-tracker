package main

import (
	"context"
	"fmt"

	"git2.gnt-global.com/jlab/gdeploy/img-tracker/config"
	"git2.gnt-global.com/jlab/gdeploy/img-tracker/pkg/k8s"
	"git2.gnt-global.com/jlab/gdeploy/img-tracker/pkg/logger"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
	// get all deploy/sts/ds/ksvc with label img-tracker=enabled
	// filter and get all labels starts with img-tracker
	// get image name using suffix after img-tracker::
	// get location of image: initContainers::name or containers::name

	// get current sha hash of image using skopeo
	// skopeo inspect docker://443533367748.dkr.ecr.ap-southeast-1.amazonaws.com/ai/graph-matching -f '{{.Digest}}' (get stdout)

	// patch image at location to new sha
	// fmt.Print(core.GetImgDigestPublic("443533367748.dkr.ecr.ap-southeast-1.amazonaws.com/ai/graph-matching"))
	lg := logger.New("main", config.MyEnvConfig.Debug)
	lg.SugaredLogger.Infof("%v", config.MyEnvConfig.UseKubeCfg)
	clientSet := k8s.GetClientSet()
	pods, err := clientSet.AppsV1().Deployments("").List(context.TODO(), metav1.ListOptions{LabelSelector: config.MyFileConfig.AppKey + "=enabled"})
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("There are %d deploys in the cluster\n", len(pods.Items))
}
