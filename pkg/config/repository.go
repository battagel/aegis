package config

type Configurator interface {
	GetConfig() (Config, error)
}
