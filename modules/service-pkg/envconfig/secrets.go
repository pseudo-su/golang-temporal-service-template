package envconfig

type GsmSecret struct {
	Gsm   string `env:"GSM" json:"gsm"`
	Value string `json:"-"`
}
