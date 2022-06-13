package config

import (
	"os"

	"github.com/caarlos0/env"
	"gopkg.in/yaml.v2"
)

func init() {
	MyEnvConfig = EnvConfig{}
	if err := env.Parse(&MyEnvConfig); err != nil {
		panic(err)
	}

	readConfigFile(MyEnvConfig.ConfigFile)
}

func readConfigFile(filepath string) {
	f, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	MyFileConfig = &FileConfig{}
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&MyFileConfig)
	if err != nil {
		panic(err)
	}
}
