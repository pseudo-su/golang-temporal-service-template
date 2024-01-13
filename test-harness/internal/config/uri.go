package config

import (
	"fmt"
	"strings"
)

type ConfigUri struct {
	Scheme string
	Host   string
	Port   string
	Path   string
}

func (uri *ConfigUri) String() string {
	if uri.Host == "" {
		return ""
	}
	var fullUri = uri.Host
	if uri.Scheme != "" {
		fullUri = fmt.Sprintf("%s://%s", uri.Scheme, fullUri)
	}
	if uri.Port != "" {
		fullUri = fmt.Sprintf("%s:%s", fullUri, uri.Port)
	}
	if uri.Path != "" {
		path := strings.TrimPrefix(uri.Path, "/")
		fullUri = fmt.Sprintf("%s/%s", fullUri, path)
	}
	return fullUri
}

func NewConfigUri(
	scheme string,
	host string,
	port string,
	path string,
) *ConfigUri {
	return &ConfigUri{
		Scheme: scheme,
		Host:   host,
		Port:   port,
		Path:   path,
	}
}
