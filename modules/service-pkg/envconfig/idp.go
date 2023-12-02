package envconfig

type IdpAegisConfig struct {
	ClientID     GsmSecret `yaml:"CLIENT_ID" json:"client_id"`
	ClientSecret GsmSecret `env:"CLIENT_SECRET" json:"client_secret"`
	Insecure     bool      `env:"INSECURE" json:"insecure"`
	Uri          UriConfig `envPrefix:"URI_" json:"uri"`
}

type IdpConfig struct {
	Mode  string         `env:"MODE" json:"mode"`
	Aegis IdpAegisConfig `envPrefix:"AEGIS_" json:"aegis"`
}
