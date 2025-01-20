package config

import (
	"bytes"
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/MagaluCloud/magalu/mgc/core"
	"github.com/MagaluCloud/magalu/mgc/core/profile_manager"

	"github.com/invopop/yaml"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

type KeyPair struct {
	KeyID     string
	KeySecret string
}
type tempConfig struct {
	configMap  map[string]interface{}
	keyPairMap map[string]KeyPair
}

type Config struct {
	pm         *profile_manager.ProfileManager
	viper      *viper.Viper
	tempConfig *tempConfig
}

const (
	CONFIG_FILE_TYPE = "yaml"
	CONFIG_FILE      = "cli.yaml"
	ENV_PREFIX       = "MGC"
)

// contextKey is an unexported type for keys defined in this package.
// This prevents collisions with keys defined in other packages.
type contextKey string

var configKey contextKey = "github.com/MagaluCloud/magalu/mgc/core/Config"

func NewContext(parent context.Context, config *Config) context.Context {
	return context.WithValue(parent, configKey, config)
}

func FromContext(ctx context.Context) *Config {
	c, _ := ctx.Value(configKey).(*Config)
	return c
}

func New(pm *profile_manager.ProfileManager) *Config {
	c := &Config{pm: pm}
	c.init()
	return c
}

func (c *Config) init() {
	v := viper.New()
	v.AutomaticEnv()
	v.SetEnvPrefix(ENV_PREFIX)
	v.SetConfigType(CONFIG_FILE_TYPE)

	c.viper = v
	_ = c.readFromFile()
}

func (c *Config) FilePath() string {
	return c.pm.Current().Dir()
}

func (c *Config) BuiltInConfigs() (map[string]*core.Schema, error) {
	loggerConfigSchema, err := loggerSchema()
	if err != nil {
		return nil, fmt.Errorf("unable to get logger config schema: %w", err)
	}

	logfilterSchema := logfilterSchema()
	defaultOutputSchema := defaultOutputSchema()

	configMap := map[string]*core.Schema{
		"logging":       loggerConfigSchema,
		"logfilter":     logfilterSchema,
		"defaultOutput": defaultOutputSchema,
	}

	return configMap, nil
}

func (c *Config) Get(key string, out any) error {
	val := reflect.ValueOf(out)
	if val.Kind() == reflect.Pointer && !val.Elem().IsValid() {
		return fmt.Errorf("result should not be nil pointer")
	}

	decodeHookFunc := mapstructure.ComposeDecodeHookFunc(
		mapstructure.StringToTimeDurationHookFunc(),
		mapstructure.StringToSliceHookFunc(","),
		mapstructure.TextUnmarshallerHookFunc(),
		stringUnmarshalHook,
	)

	if c.tempConfig != nil {
		if outVal, found := c.tempConfig.configMap[key]; found {

			decodeConfig := mapstructure.DecoderConfig{
				DecodeHook: decodeHookFunc,
				Result:     out,
			}

			decoder, err := mapstructure.NewDecoder(&decodeConfig)
			if err != nil {
				return fmt.Errorf("fail to create a config decoder: %s ", err.Error())
			}

			err = decoder.Decode(outVal)
			if err != nil {
				return fmt.Errorf("fail to run config decoder: %s", err.Error())
			}
			return nil
		}
	}

	return c.viper.UnmarshalKey(
		key,
		out,
		viper.DecodeHook(decodeHookFunc),
	)
}

func stringUnmarshalHook(f reflect.Value, t reflect.Value) (interface{}, error) {
	str, ok := f.Interface().(string)
	if !ok {
		return f.Interface(), nil
	}

	dereferenced := t
	derefKind := dereferenced.Kind()

	for derefKind == reflect.Pointer || derefKind == reflect.Interface {
		dereferenced = dereferenced.Elem()
		derefKind = dereferenced.Kind()
	}

	switch derefKind {
	case reflect.Struct, reflect.Array, reflect.Slice:
		target := dereferenced.Addr().Interface()
		err := yaml.Unmarshal([]byte(str), target)

		return target, err
	case reflect.Map:
		target := t.Interface()
		err := yaml.Unmarshal([]byte(str), target)

		return target, err
	case reflect.Invalid:
		// Try to decode string to any, it may work. If not, just return the value as-is
		target := t.Interface()
		err := yaml.Unmarshal([]byte(str), &target)
		if err != nil {
			return str, nil
		} else {
			return target, nil
		}
	default:
		return str, nil
	}

}

func marshalValueIfNeeded(value any) (any, error) {
	if value == nil {
		return value, nil
	}

	v := reflect.ValueOf(value)
	kind := v.Type().Kind()

	for kind == reflect.Pointer || kind == reflect.Interface {
		v = v.Elem()
		kind = v.Kind()
	}

	switch kind {
	case reflect.Map, reflect.Struct, reflect.Array, reflect.Slice:
		b, err := yaml.Marshal(v.Interface())
		if err != nil {
			return nil, err
		}
		return string(b), nil
	default:
		return value, nil
	}
}

func (c *Config) NewTempConfig() {
	if c.tempConfig == nil {
		c.tempConfig = &tempConfig{
			configMap:  make(map[string]interface{}),
			keyPairMap: make(map[string]KeyPair),
		}
	}
	//Default tmp configs
	c.tempConfig.configMap["region"] = "br-se1"
}

func (c *Config) SetTempConfig(key string, value interface{}) error {
	if value == nil {
		return nil
	}

	marshaled, err := marshalValueIfNeeded(value)
	if err != nil {
		return fmt.Errorf("unable to marshal config %s: %w", key, err)
	}

	if c.tempConfig == nil {
		c.NewTempConfig()
	}

	c.tempConfig.configMap[key] = marshaled

	return nil
}

func (c *Config) GetTempKeyPair(pairName string) *KeyPair {
	if c.tempConfig == nil {
		return nil
	}
	return &KeyPair{KeyID: c.tempConfig.keyPairMap[pairName].KeyID, KeySecret: c.tempConfig.keyPairMap[pairName].KeySecret}

}

func (c *Config) AddTempKeyPair(pairName, keyID, keySecret string) {
	if c.tempConfig == nil {
		c.NewTempConfig()
	}
	c.tempConfig.keyPairMap[pairName] = KeyPair{KeyID: keyID, KeySecret: keySecret}
}
func (c *Config) TempConfig() map[string]interface{} {
	result := make(map[string]interface{})
	for key, value := range c.tempConfig.configMap {
		result[key] = value
	}
	return result
}

func (c *Config) Set(key string, value interface{}) error {
	marshaled, err := marshalValueIfNeeded(value)
	if err != nil {
		return fmt.Errorf("unable to marshal config %s: %w", key, err)
	}

	c.viper.Set(key, marshaled)

	obj := c.viper.AllSettings()
	return c.saveToConfigFile(obj)
}

func (c *Config) Delete(key string) error {
	configMap := c.viper.AllSettings()
	key = strings.ToLower(key)
	if _, ok := configMap[key]; !ok {
		return nil
	}

	delete(configMap, key)

	if err := c.saveToConfigFile(configMap); err != nil {
		return err
	}

	return c.readFromFile()
}

func (c *Config) readFromFile() (err error) {
	data, err := c.pm.Current().Read(CONFIG_FILE)
	if err != nil {
		return err
	}

	return c.viper.ReadConfig(bytes.NewBuffer(data))
}

func (c *Config) saveToConfigFile(configMap map[string]interface{}) error {
	encodedConfig, err := yaml.Marshal(configMap)
	if err != nil {
		return err
	}

	if err = c.pm.Current().Write(CONFIG_FILE, encodedConfig); err != nil {
		return fmt.Errorf("error writing to config file: %w", err)
	}

	return nil
}
