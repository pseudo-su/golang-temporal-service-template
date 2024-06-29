package envconfig

type CredentialMode string

const (
	CredentialModeEnvVars  CredentialMode = "env_var"
	CredentialModeGsmFetch CredentialMode = "gsm_fetch"
)

type Credential struct {
	Config string `json:"config"`
	Value  string `json:"-"`
}
