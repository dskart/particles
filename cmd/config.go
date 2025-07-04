package cmd

import (
	"github.com/dskart/particles/api"
	"github.com/dskart/particles/app"
)

type Config struct {
	App app.Config `yaml:"App" env:"APP"`
	API api.Config `yaml:"API" env:"API"`
}

func NewConfig() Config {
	return Config{
		App: app.NewConfig(),
		API: api.NewConfig(),
	}
}
