package config

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type EnvConfig struct {
	Foo        string
	PointerFoo *string   `env:"POINTER_FOO"`
	SubFoo     EnvSubFoo `env:"SUB_FOO"`
	Float      float64   `env:"FLOAT"`
}

type EnvSubFoo struct {
	SubBar        int `env:"SUB_BAR"`
	PointerSubBar *int
}

func toPointer[T any](v T) *T {
	return &v
}

func TestConfig_UnmarshalConfigFromEnv(t *testing.T) {
	ctx := context.Background()
	for name, tc := range map[string]struct {
		options     []func(*UnmarshalConfigOptions)
		environment map[string]string
		expected    EnvConfig
	}{
		"WithPrefix": {
			options: []func(*UnmarshalConfigOptions){WithPrefix("TEST")},
			environment: map[string]string{
				"TEST__FOO":                    "foo",
				"TEST__POINTER_FOO":            "pointer_foo",
				"TEST__FLOAT":                  "2.1234123",
				"TEST__SUB_FOO__SUB_BAR":       "1",
				"TEST__SUB_FOO__POINTERSUBBAR": "2",
			},
			expected: EnvConfig{
				Foo:        "foo",
				PointerFoo: toPointer("pointer_foo"),
				Float:      2.1234123,
				SubFoo: EnvSubFoo{
					SubBar:        1,
					PointerSubBar: toPointer(2),
				},
			},
		},
		"NoPrefix": {
			environment: map[string]string{
				"FOO": "foo",
			},
			expected: EnvConfig{
				Foo: "foo",
			},
		},
		"WithNilPointer": {
			environment: map[string]string{
				"FOO": "foo",
			},
			expected: EnvConfig{
				Foo:        "foo",
				PointerFoo: nil,
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			config := EnvConfig{}
			for k, v := range tc.environment {
				err := os.Setenv(k, v)
				require.NoError(t, err)
			}
			err := UnmarshalConfigFromEnv(ctx, &config, tc.options...)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, config)
		})
	}
}
