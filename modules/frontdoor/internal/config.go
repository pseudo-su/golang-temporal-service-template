package internal

import "github.com/pseudo-su/golang-temporal-service-template/modules/service-pkg/envconfig"

type FrontdoorConfig struct {
	App           envconfig.AppConfig           `envPrefix:"APP_"`
	Tcp           envconfig.TcpConfig           `envPrefix:"TCP_"`
	Log           envconfig.LogConfig           `envPrefix:"LOG_"`
	Idp           envconfig.IdpConfig           `envPrefix:"IDP_"`
	Temporal      envconfig.TemporalConfig      `envPrefix:"TEMPORAL_"`
	OpenTelemetry envconfig.OpenTelemetryConfig `envPrefix:"OPEN_TELEMETRY_"`
	Launchdarkly  envconfig.LaunchdarklyConfig  `envPrefix:"LAUNCHDARKLY_"`
}
