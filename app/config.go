package app

type Config struct {
	MaxNLogs int `yaml:"MaxNLogs"`
}

func (c *Config) Validate() error {
	return nil
}
