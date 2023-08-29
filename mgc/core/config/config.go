package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"

	"magalu.cloud/core"
	"magalu.cloud/core/logger"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

// contextKey is an unexported type for keys defined in this package.
// This prevents collisions with keys defined in other packages.
type contextKey string

type Config struct {
	path     string
	fileName string
}

const (
	CONFIG_NAME = "cli"
	CONFIG_TYPE = "yaml"
	CONFIG_FILE = CONFIG_NAME + "." + CONFIG_TYPE
)

var configKey contextKey = "magalu.cloud/core/Config"

func NewContext(parent context.Context, config *Config) context.Context {
	return context.WithValue(parent, configKey, config)
}

func FromContext(ctx context.Context) *Config {
	c, _ := ctx.Value(configKey).(*Config)
	return c
}

func New() *Config {
	path, err := core.BuildMGCPath()
	if err != nil {
		// TODO: when it's done, use logger instead
		log.Println(err)
		return &Config{}
	}

	viper.SetConfigName(CONFIG_NAME)
	viper.SetConfigType(CONFIG_TYPE)
	viper.AddConfigPath(path)
	viper.AutomaticEnv()

	_ = viper.ReadInConfig()
	return &Config{path: path, fileName: CONFIG_FILE}
}

func (c *Config) BuiltInConfigs() (map[string]any, error) {
	loggerConfigSchema, err := logger.ConfigSchema()
	if err != nil {
		return nil, fmt.Errorf("unable to get logger config schema: %w", err)
	}

	configMap := map[string]any{
		"logging": loggerConfigSchema,
	}

	return configMap, nil
}

func (c *Config) Get(key string) any {
	return viper.Get(key)
}

func (c *Config) Set(key string, value interface{}) error {
	if err := os.MkdirAll(c.path, core.FILE_PERMISSION); err != nil {
		return fmt.Errorf("error creating dir at %s: %w", c.path, err)
	}
	viper.Set(key, value)

	if err := viper.WriteConfigAs(path.Join(c.path, c.fileName)); err != nil {
		return fmt.Errorf("error writing to config file: %w", err)
	}

	return nil
}

func (c *Config) Delete(key string) error {
	configMap := viper.AllSettings()
	if _, ok := configMap[key]; !ok {
		return nil
	}

	delete(configMap, key)

	if err := saveToConfigFile(c, configMap); err != nil {
		return err
	}

	return nil
}

func saveToConfigFile(c *Config, configMap map[string]interface{}) error {
	encodedConfig, err := yaml.Marshal(configMap)
	if err != nil {
		return err
	}

	if err = os.MkdirAll(c.path, core.FILE_PERMISSION); err != nil {
		return fmt.Errorf("error creating dir at %s: %w", c.path, err)
	}

	if err = os.WriteFile(path.Join(c.path, c.fileName), encodedConfig, core.FILE_PERMISSION); err != nil {
		return fmt.Errorf("error writing to config file: %w", err)
	}

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return nil
}
