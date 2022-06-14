package config

import (
	"encoding/json"
	"os"

	"github.com/caarlos0/env"
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
	decoder := json.NewDecoder(f)
	err = decoder.Decode(&MyFileConfig)
	if err != nil {
		panic(err)
	}
}
