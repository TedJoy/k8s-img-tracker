package config

var MyEnvConfig EnvConfig
var MyFileConfig *FileConfig

type EventBroker int8

type EnvConfig struct {
	ConfigFile     string `env:"CONFIG_FILE"      envDefault:"./config.yml"`
	Debug          bool   `env:"DEBUG"            envDefault:"false"`
	UseKubeCfg     bool   `env:"USE_KUBECONFIG"   envDefault:"true"`
	KubeConfigFile string `env:"KUBECONFIG"       envDefault:"${HOME}/.kube/config" envExpand:"true"`
}

type Registry struct {
	Type     string `yaml:"type"`
	Metadata interface{}
}

type FileConfig struct {
	EventBroker struct {
		Type string `yaml:"type"`
		Url  string `yaml:"connect_url"`
	} `yaml:"event_broker"`

	Registries map[string]Registry `yaml:"registries"`
	AppKey     string              `yaml:"app_key"`
}
