package app

type Config struct {
	MaxSessTime string    `yaml:"MaxSessTime" env:"MAX_SESS_TIME"`
	SimConfig   SimConfig `yaml:"SimConfig" env:"SIM_CONFIG"`
}

func NewConfig() Config {
	return Config{
		MaxSessTime: "5m",
		SimConfig:   NewSimConfig(),
	}
}

func (c *Config) Validate() error {
	if err := c.SimConfig.Validate(); err != nil {
		return err
	}

	return nil
}
