package config

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

func WithSeperator(seperator string) func(*UnmarshalConfigOptions) {
	return func(options *UnmarshalConfigOptions) {
		options.seperator = seperator
	}
}

func WithPrefix(prefix string) func(*UnmarshalConfigOptions) {
	return func(options *UnmarshalConfigOptions) {
		options.prefix = prefix
	}
}

// UnmarshalConfigFromEnv populates config with values from environment variables. The names of the
// environment variables read are formed by joining the config struct's field names as capital letters with a seperate `__`.
// You can also provide a custom value to unmarshall with the `env:"SOMETHING"` struct tag.
// For example consider this config struct:
//
//	type Config struct {
//		Foo    string `env:"FOO"`
//		SubFoo SubFoo `env:"SUB_FOO"`
//	}
//
//	type SubFoo struct {
//		SubBar        int  `env:"SUB_PAR"`
//		PointerSubBar *int `env:"POINTER_SUB_BAR"`
//	}
//
// The following environment variables will be read:
//
//	FOO
//	SUB_FOO__SUB_BAR
//	SUB_FOO__POINTER_BAR
//
// The user can also provied a prefix with `WithPrefix("MY_PREFIX")` that will prepend each env variable lookup with prefix like so:
//
//	TEST__FOO
//	TEST__SUB_FOO__SUB_BAR
//	TEST__SUB_FOO__POINTER_BAR
func UnmarshalConfigFromEnv(ctx context.Context, config any, opts ...func(*UnmarshalConfigOptions)) error {
	options := UnmarshalConfigOptions{
		seperator: "__",
	}
	for _, opt := range opts {
		opt(&options)
	}

	prefix := options.prefix
	seperator := options.seperator
	_, err := unmarshalConfig(prefix, seperator, reflect.ValueOf(config), func(key string) (*string, error) {
		v, ok := os.LookupEnv(key)
		if !ok {
			return nil, nil
		}
		return &v, nil
	})
	return err
}

func unmarshalConfig(prefix string, seperator string, v reflect.Value, lookup func(string) (*string, error)) (didUnmarshal bool, err error) {
	if v.Kind() == reflect.Ptr && !v.IsNil() {
		if env, err := lookup(prefix); err != nil {
			return false, err
		} else if env != nil {
			switch dest := v.Interface().(type) {
			case *bool:
				switch *env {
				case "false":
					*dest = false
				case "true":
					*dest = true
				default:
					return false, fmt.Errorf("boolean config must be \"true\" or \"false\"")
				}
			case *int:
				n, err := strconv.Atoi(*env)
				if err != nil {
					return false, fmt.Errorf("invalid value for integer config")
				}
				*dest = n
			case *float64:
				n, err := strconv.ParseFloat(*env, 64)
				if err != nil {
					return false, fmt.Errorf("invalid value for float64 config")
				}
				*dest = n
			case *string:
				*dest = *env
			case *[]byte:
				buf, err := base64.StdEncoding.DecodeString(*env)
				if err != nil {
					return false, fmt.Errorf("byte slice configs must be base64 encoded")
				}
				*dest = buf
			case *[]int:
				parts := strings.Split(*env, ",")
				intParts := make([]int, len(parts))
				for i, part := range parts {
					intParts[i], err = strconv.Atoi(strings.TrimSpace(part))
					if err != nil {
						return false, fmt.Errorf("invalid value for integer config")
					}
				}
				*dest = intParts
			case *[]string:
				parts := strings.Split(*env, ",")
				for i, part := range parts {
					parts[i] = strings.TrimSpace(part)
				}
				*dest = parts
			default:
				return false, fmt.Errorf("unsupported environment config type %T", v.Elem().Interface())
			}
			return true, nil
		}
	}

	if v.Kind() == reflect.Ptr {
		if !v.IsNil() {
			return unmarshalConfig(prefix, seperator, v.Elem(), lookup)
		}

		newV := reflect.New(v.Type().Elem())
		didUnmarshal, err := unmarshalConfig(prefix, seperator, newV, lookup)
		if err != nil {
			return false, err
		} else if didUnmarshal {
			v.Set(newV)
		}
		return didUnmarshal, nil
	}

	if v.Kind() == reflect.Struct {
		t := v.Type()
		for i := range v.NumField() {
			field := t.Field(i)
			nextPrefix, ok := field.Tag.Lookup("env")
			if !ok {
				nextPrefix = strings.ToUpper(field.Name)
			}

			if prefix != "" {
				nextPrefix = prefix + seperator + nextPrefix
			}

			var nextV reflect.Value
			if v.Field(i).Kind() == reflect.Ptr {
				nextV = v.Field(i)
			} else {
				nextV = v.Field(i).Addr()
			}

			didUnmarshalField, err := unmarshalConfig(nextPrefix, seperator, nextV, lookup)
			if err != nil {
				return false, err
			}
			if didUnmarshalField {
				didUnmarshal = true
			}
		}
	}

	return didUnmarshal, nil
}
