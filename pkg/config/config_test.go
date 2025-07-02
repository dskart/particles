package config

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Config struct {
	Foo    string  `yaml:"Foo" env:"FOO"`
	SubFoo SubFoo  `yaml:"SubFoo" env:"SUB_FOO"`
	Float  float64 `yaml:"Float" env:"FLOAT"`
}

type SubFoo struct {
	SubBar        int `yaml:"SubBar" env:"SUB_BAR"`
	PointerSubBar int `yaml:"PointerSubBar" env:"POINTER_SUB_BAR"`
}

func TestConfig_UnmarshalConfig(t *testing.T) {
	for name, tc := range map[string]struct {
		Prefix      string
		Environment map[string]string
		In          Config
		Expected    *Config
	}{
		"EnvOnly": {
			Prefix: "TEST",
			Environment: map[string]string{
				"TEST__FOO":                      "foo",
				"TEST__FLOAT":                    "2.1234123",
				"TEST__SUB_FOO__SUB_BAR":         "1",
				"TEST__SUB_FOO__POINTER_SUB_BAR": "2",
			},
			Expected: &Config{
				Foo:   "foo",
				Float: 2.1234123,
				SubFoo: SubFoo{
					SubBar:        1,
					PointerSubBar: 2,
				},
			},
		},
		"Override": {
			Prefix: "TEST",
			Environment: map[string]string{
				"TEST__FOO": "foo",
			},
			In: Config{
				Foo: "bar",
			},
			Expected: &Config{
				Foo: "foo",
			},
		},
		"NoPrefix": {
			Prefix: "",
			Environment: map[string]string{
				"FOO": "foo",
			},
			Expected: &Config{
				Foo: "foo",
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			config := tc.In
			_, err := unmarshalConfig(tc.Prefix, "__", reflect.ValueOf(&config), func(key string) (*string, error) {
				if v, ok := tc.Environment[key]; ok {
					return &v, nil
				}
				return nil, nil
			})
			assert.NoError(t, err)
			assert.Equal(t, tc.Expected, &config)
		})
	}
}
