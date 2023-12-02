package envconfig

type CloudEventsConfig struct {
	ProjectId    string `env:"PROJECT_ID" json:"project_id"`
	TopicId      string `env:"TOPIC_ID" json:"topic_id"`
	InsecureHost bool   `env:"INSECURE_HOST" json:"insecure_host"`
}
