package config

var MyEnvConfig EnvConfig
var MyFileConfig *FileConfig

type EventBroker int8

type EnvConfig struct {
	ConfigFile     string `env:"CONFIG_FILE"      envDefault:"./config.json"`
	Debug          bool   `env:"DEBUG"            envDefault:"false"`
	UseKubeCfg     bool   `env:"USE_KUBECONFIG"   envDefault:"true"`
	KubeConfigFile string `env:"KUBECONFIG"       envDefault:"${HOME}/.kube/config" envExpand:"true"`
}

type Registry struct {
	Type     string            `json:"type"`
	Metadata map[string]string `json:"metadata"`
}

type FileConfig struct {
	EventBroker struct {
		Type     string            `json:"type"`
		Url      string            `json:"connect_url"`
		Metadata map[string]string `json:"metadata"`
	} `json:"event_broker"`

	Registries map[string]Registry `json:"registries"`
	AppKey     string              `json:"app_key"`
}
