package cmd

import (
	"github.com/dskart/particles/api"
	"github.com/dskart/particles/app"
)

type Config struct {
	App app.Config `yaml:"App"`
	API api.Config `yaml:"API"`
}
