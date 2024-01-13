package envconfig

import (
	"github.com/caarlos0/env/v10"
)

func ParseEnv[T any](cfg T) (T, error) {
	if err := env.Parse(cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}

func FetchSecrets[T any](cfg T) {
	// TODO: walk the cfg input and fetch the .Value for any GsmSecrets
}
