package config

import (
	"fmt"
	"strings"
)

func (cfg *TestsuiteEnvConfig) EnvName() string {
	if s, ok := strings.CutSuffix(cfg.Env, ".ext"); ok {
		return s
	}
	return cfg.Env
}

func (cfg *TestsuiteEnvConfig) TemporalGrpcUri() string {
	uri := NewConfigUri("", cfg.TemporalUriHost, cfg.TemporalUriPort, "")
	return uri.String()
}

func (cfg *TestsuiteEnvConfig) ServiceFrontdoorApiGrpcUri() string {
	host := fmt.Sprintf("%s%s", cfg.ServiceFrontdoorApiUriHostPrefix, cfg.ServiceFrontdoorApiUriHost)
	uri := NewConfigUri("", host, cfg.ServiceFrontdoorApiUriPort, "")
	return uri.String()
}

func (cfg *TestsuiteEnvConfig) ServiceFrontdoorHttpUri() string {
	host := fmt.Sprintf("%s%s", cfg.ServiceFrontdoorApiUriHostPrefix, cfg.ServiceFrontdoorApiUriHost)
	uri := NewConfigUri(cfg.ServiceFrontdoorApiUriHttpScheme, host, cfg.ServiceFrontdoorApiUriPort, "")
	return uri.String()
}
