package envconfig

type CommandCentreConfig struct {
	Environment  string `env:"ENVIRONMENT" json:"environment"`
	EmulatorHost string `env:"EMULATOR_HOST" json:"emulator_host"`
	ProjectId    string `env:"PROJECT_ID" json:"project_id"`
	TopicId      string `env:"TOPIC_ID" json:"topic_id"`
}
