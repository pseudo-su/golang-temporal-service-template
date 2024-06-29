package internal

import "github.com/pseudo-su/golang-temporal-service-template/modules/service-pkg/envconfig"

type WorkerConfig struct {
	App         envconfig.AppConfig         `envPrefix:"APP_"`
	Tcp         envconfig.TcpConfig         `envPrefix:"TCP_"`
	Log         envconfig.LogConfig         `envPrefix:"LOG_"`
	Idp         envconfig.IdpConfig         `envPrefix:"IDP_"`
	Temporal    envconfig.TemporalConfig    `envPrefix:"TEMPORAL_"`
	CloudEvents envconfig.CloudEventsConfig `envPrefix:"CLOUD_EVENTS_"`
}

func (mc *WorkerConfig) LogConfig() envconfig.LogConfig {
	return mc.Log
}
