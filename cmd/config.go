package cmd

import "github.com/dskart/particles/app"

type Config struct {
	App app.Config `yaml:"App"`
}
