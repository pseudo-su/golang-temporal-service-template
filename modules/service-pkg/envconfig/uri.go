package envconfig

import (
	"fmt"
	"strings"
)

type UriConfig struct {
	Scheme string `env:"SCHEME" json:"scheme"`
	Host   string `env:"HOST" json:"host"`
	Port   string `env:"PORT" json:"port"`
	Path   string `env:"PATH" json:"path"`
}

type GrpcUriConfig struct {
	Host string `env:"HOST"`
	Port string `env:"PORT"`
}

func (uri GrpcUriConfig) AsString() string {
	return buildUri("", uri.Host, uri.Port, "")
}

func buildUri(scheme, host, port, path string) string {
	if host == "" {
		return ""
	}
	var fullUri = host
	if scheme != "" {
		fullUri = fmt.Sprintf("%s://%s", scheme, fullUri)
	}
	if port != "" {
		fullUri = fmt.Sprintf("%s:%s", fullUri, port)
	}
	if path != "" {
		p := strings.TrimPrefix(path, "/")
		fullUri = fmt.Sprintf("%s/%s", fullUri, p)
	}
	return fullUri
}
