package envconfig

type TemporalConfig struct {
	Uri            GrpcUriConfig `envPrefix:"URI_" json:"uri"`
	Namespace      string        `env:"NAMESPACE" json:"namespace"`
	TaskQueue      string        `env:"TASK_QUEUE" json:"task_queue"`
	ServiceAccount string        `env:"SERVICE_ACCOUNT" json:"service_account"`
}
