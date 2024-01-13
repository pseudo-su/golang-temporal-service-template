package internal

import "github.com/pseudo-su/golang-temporal-service-template/modules/service-pkg/envconfig"

type MocksConfig struct {
	App envconfig.AppConfig `envPrefix:"APP_"`
	Tcp envconfig.TcpConfig `envPrefix:"TCP_"`
	Log envconfig.LogConfig `envPrefix:"LOG_"`
}

func (mc *MocksConfig) LogConfig() envconfig.LogConfig {
	return mc.Log
}
