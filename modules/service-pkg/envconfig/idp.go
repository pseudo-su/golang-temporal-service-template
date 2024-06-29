package envconfig

type IdpConfig struct {
	Client IdpClientConfig `envPrefix:"CLIENT_" json:"client"`
	Verify IdpVerifyConfig `envPrefix:"VERIFY_" json:"verify"`
}

type IdpClientConfig struct {
	Id         string         `env:"ID"`
	SecretMode CredentialMode `env:"SECRET_MODE"`
	Secret     Credential     `env:"SECRET"`
	Uri        UriConfig      `envPrefix:"URI_" json:"uri"`
	Enabled    bool           `env:"ENABLED" json:"enabled"`
}

type IdpVerifyConfig struct {
	Insecure bool `env:"INSECURE"`

	ProviderSystem   IdpProviderConfig `envPrefix:"PROVIDER_SYSTEM_"`
	ProviderCustomer IdpProviderConfig `envPrefix:"PROVIDER_CUSTOMER_"`
}

type IdpProviderConfig struct {
	Name          string `env:"NAME"`
	URL           string `env:"URL"`
	RefreshPeriod string `env:"REFRESH"`
	Enabled       bool   `env:"ENABLED"`
}
