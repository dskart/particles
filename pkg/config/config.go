package config

import (
	"context"
)

type RootConfig any

type UnmarshalConfigOptions struct {
	filePath  string
	seperator string
	prefix    string
}

// UnmarshalConfig populates config with values from environment variables and a yaml file.
// Look at `UnmarshalConfigFromEnv` for more information on environment unmarshalling.
func UnmarshalConfig(ctx context.Context, prefix string, config any, opts ...func(*UnmarshalConfigOptions)) error {
	err := UnmarshalConfigFromFile(config, opts...)
	if err != nil {
		return err
	}

	err = UnmarshalConfigFromEnv(ctx, config, opts...)
	if err != nil {
		return err
	}

	return nil
}
