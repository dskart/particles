package config

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type RootConfig any

type UnmarshalConfigOptions struct {
	filePath  string
	seperator string
	prefix    string
	smClient  *secretsmanager.Client
}

func WithSecretManager(smClient *secretsmanager.Client) func(*UnmarshalConfigOptions) {
	return func(options *UnmarshalConfigOptions) {
		options.smClient = smClient
	}
}

// UnmarshalConfig populates config with values from environment variables and a yaml file.
// Look at `UnmarshalConfigFromEnv` for more information on environment unmarshalling.
func UnmarshalConfig(ctx context.Context, config any, opts ...func(*UnmarshalConfigOptions)) error {
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
