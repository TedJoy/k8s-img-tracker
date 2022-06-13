package core

import (
	"os/exec"
	"strings"

	"git2.gnt-global.com/jlab/gdeploy/img-tracker/config"
	"git2.gnt-global.com/jlab/gdeploy/img-tracker/pkg/logger"
)

func GetImgDigestPrivate(user, password, registry, image string) string {
	loginPrivateRegistry(user, password, registry)
	return GetImgDigestPublic(image)
}

func loginPrivateRegistry(user, password, registry string) {
	if strings.HasSuffix(registry, "amazonaws.com") {
	}
}

func GetImgDigestPublic(image string) string {
	out, err := exec.Command("skopeo", "inspect", "docker://"+image, "-f={{.Digest}}").Output()

	lg := logger.New("main", config.MyEnvConfig.Debug)

	if err != nil {
		lg.SugaredLogger.Panicf("error: %v", err)
	}

	return strings.TrimSpace(string(out))
}
