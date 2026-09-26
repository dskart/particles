package api

import "fmt"

type Config struct {
	MaxNumSessions int    `yaml:"MaxNumSessions" env:"MAX_NUM_SESSIONS"`
	SSHHostKey     string `yaml:"SSHostKey" env:"SSH_HOST_KEY"`
	// PublicHost is the public hostname shown in the connection instructions on /.
	// Empty falls back to the request's Host header.
	PublicHost string `yaml:"PublicHost" env:"PUBLIC_HOST"`
}

func NewConfig() Config {
	// defaults
	return Config{
		MaxNumSessions: 10,
		SSHHostKey:     "",
		PublicHost:     "",
	}
}

func (c *Config) Validate() error {
	if c.MaxNumSessions <= 0 {
		return fmt.Errorf("MaxNumSessions '%d'needs to be bigger than 0", c.MaxNumSessions)
	}
	return nil
}
