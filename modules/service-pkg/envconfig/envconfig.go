package envconfig

import (
	"context"
	"log/slog"
	"os"

	"github.com/caarlos0/env/v10"
)

func ParseEnv[T any](cfg T) T {
	if err := env.Parse(cfg); err != nil {
		slog.ErrorContext(context.Background(), "unable to parse environment variables", slog.Any("error", err))
		os.Exit(1)
	}

	return cfg
}

func FetchSecrets[T any](cfg T) {
	// TODO: walk the cfg input and fetch the .Value for any GsmSecrets
}
