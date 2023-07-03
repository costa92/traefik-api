package config

type Provider interface {
	Name() string
	Config(*providerHelper, interface{}) ([]byte, error)
}

type providerHelper struct {
	configFile string
	log        Logger
}
