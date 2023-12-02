package envconfig

type OpenTelemetryConfig struct {
	Uri GrpcUriConfig `envPrefix:"URI_" json:"uri"`
}
