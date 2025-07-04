package api

import "fmt"

type Config struct {
	MaxNumSessions int    `yaml:"MaxNumSessions" env:"MAX_NUM_SESSIONS"`
	SSHHostKey     string `yaml:"SSHostKey" env:"SSH_HOST_KEY"`
}

func (c *Config) Validate() error {
	if c.MaxNumSessions <= 0 {
		return fmt.Errorf("MaxNumSessions '%d'needs to be bigger than 0", c.MaxNumSessions)
	}
	return nil
}
