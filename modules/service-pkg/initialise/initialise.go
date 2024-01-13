package initialise

import (
	"context"
	"log/slog"

	"github.com/pseudo-su/golang-temporal-service-template/modules/service-pkg/envconfig"
	"github.com/pseudo-su/golang-temporal-service-template/modules/service-pkg/initialise/logsetup"
)

type InitConfigInterface interface {
	LogConfig() envconfig.LogConfig
}

func ServiceWithConfig[T InitConfigInterface](ctx context.Context, cfg T) (T, error) {
	cfg, err := envconfig.ParseEnv(cfg)
	if err != nil {
		return cfg, err
	}

	logCfg := cfg.LogConfig()
	logger := logsetup.NewLogger(
		logsetup.WithLevelStr(logCfg.Level),
		logsetup.WithModeStr(logCfg.Mode),
	)
	slog.SetDefault(logger)

	return cfg, nil
}
