package main

import (
	"git2.gnt-global.com/jlab/gdeploy/img-tracker/config"
	"git2.gnt-global.com/jlab/gdeploy/img-tracker/pkg/logger"
)

func main() {
	lg := logger.New("main", config.MyEnvConfig.Debug)
	lg.SugaredLogger.Infof("%v", config.MyEnvConfig.ConfigFile)
	lg.SugaredLogger.Infof("%v", config.MyFileConfig)
}
