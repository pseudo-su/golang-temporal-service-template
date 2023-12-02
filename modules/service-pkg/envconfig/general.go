package envconfig

type LogConfig struct {
	Level string `env:"LEVEL" json:"level"`
	Mode  string `env:"MODE" json:"mode"`
}

type AppConfig struct {
	Name string `env:"NAME" json:"name"`
	Env  string `env:"ENV" json:"env"`
}

type TcpConfig struct {
	Port int `env:"PORT" json:"port"`
}
