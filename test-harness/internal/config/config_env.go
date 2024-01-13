package config

type TestsuiteEnvConfig struct {
	Env string `env:"ENV"`

	TemporalUriHost        string `env:"TEMPORAL_URI_HOST"`
	TemporalUriPort        string `env:"TEMPORAL_URI_PORT"`
	TemporalNamespace      string `env:"TEMPORAL_NAMESPACE"`
	TemporalWaiterTimeout  int64  `env:"TEMPORAL_WAITER_TIMEOUT"`
	TemporalDiscardLogs    bool   `env:"TEMPORAL_DISCARD_LOGS"`
	TemporalServiceAccount string `env:"TEMPORAL_SERVICE_ACCOUNT"`

	ServiceFrontdoorApiUriHost       string `env:"SERVICE_FRONTDOOR_API_URI_HOST"`
	ServiceFrontdoorApiUriHostPrefix string `env:"SERVICE_FRONTDOOR_API_URI_HOST_PREFIX"`
	ServiceFrontdoorApiUriPort       string `env:"SERVICE_FRONTDOOR_API_URI_PORT"`
	ServiceFrontdoorApiInsecure      bool   `env:"SERVICE_FRONTDOOR_API_INSECURE"`
	ServiceFrontdoorApiUriHttpScheme string `env:"SERVICE_FRONTDOOR_API_URI_HTTP_SCHEME"`

	SSLEnabled bool `env:"SSL_ENABLED"`
}
