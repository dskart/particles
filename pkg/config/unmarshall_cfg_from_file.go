package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

func WithFilePath(filePath string) func(*UnmarshalConfigOptions) {
	return func(options *UnmarshalConfigOptions) {
		options.filePath = filePath
	}
}

// UnmarshalConfigFromFile populates config with values from a yaml file.
// It will by default look for a file named config.yml in the current working directory.
// You can override this by passing the WithFilePath option.
func UnmarshalConfigFromFile(config any, opts ...func(*UnmarshalConfigOptions)) error {
	options := UnmarshalConfigOptions{}
	for _, opt := range opts {
		opt(&options)
	}

	if options.filePath != "" {
		f, err := os.Open(options.filePath)
		if err != nil {
			return err
		}
		defer func() {
			if err := f.Close(); err != nil {
				fmt.Fprintf(os.Stderr, "error closing file: %v\n", err)
			}
		}()
		if err := yaml.NewDecoder(f).Decode(config); err != nil {
			return fmt.Errorf("failed to decode config: %w", err)
		}
	} else if f, err := os.Open("config.yaml"); err == nil {
		if err := decodeConfigFile(f, config); err != nil {
			return err
		}
	} else if f, err := os.Open("config.yml"); err == nil {
		if err := decodeConfigFile(f, config); err != nil {
			return err
		}
	}

	return nil
}

func decodeConfigFile(f *os.File, config any) error {
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "error closing file: %v\n", err)
		}
	}()
	if err := yaml.NewDecoder(f).Decode(config); err != nil {
		return fmt.Errorf("error decoding config: %w", err)
	}

	return nil
}
