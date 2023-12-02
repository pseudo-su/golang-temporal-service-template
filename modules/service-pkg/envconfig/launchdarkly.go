package envconfig

type LaunchdarklySaasConfig struct {
	SdkKey GsmSecret `env:"SDK_KEY" json:"sdk_key"`
}

type LaunchdarklyLocalConfig struct {
	Filepath string `env:"FILEPATH" json:"filepath"`
}

type LaunchdarklyConfig struct {
	Mode            string                 `env:"MODE" json:"mode"`
	ClientTimeout   string                 `env:"CLIENT_TIMEOUT" json:"client_timeout"`
	Saas            LaunchdarklySaasConfig `envPrefix:"SAAS_" json:"saas"`
	LocalDatasource LaunchdarklySaasConfig `envPrefix:"LOCAL_DATASOURCE_" json:"local_datasource"`
	ProxyUri        UriConfig              `envPrefix:"PROXY_URI_" json:"proxy_uri"`
	RelayProxyUri   UriConfig              `envPrefix:"RELAY_PROXY_URI_" json:"relay_proxy_uri"`
}
