package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"

	"magalu.cloud/core"
	"magalu.cloud/core/logger"

	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

// contextKey is an unexported type for keys defined in this package.
// This prevents collisions with keys defined in other packages.
type contextKey string

type Config struct {
	path  string
	fs    afero.Fs
	viper *viper.Viper
}

const (
	CONFIG_NAME   = "cli"
	CONFIG_TYPE   = "yaml"
	CONFIG_FOLDER = ".config/mgc"
	CONFIG_FILE   = CONFIG_NAME + "." + CONFIG_TYPE
	ENV_PREFIX    = "MGC"
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
	dirname, err := core.BuildMGCPath()
	if err != nil {
		// TODO: when it's done, use logger instead
		log.Println(err)
	}

	c := &Config{}
	c.init(dirname, afero.NewOsFs())
	return c
}

func (c *Config) init(dirname string, fs afero.Fs) {
	v := viper.New()
	v.SetFs(fs)
	v.SetConfigName(CONFIG_NAME)
	v.SetConfigType(CONFIG_TYPE)
	v.AddConfigPath(dirname)
	v.AutomaticEnv()
	v.SetEnvPrefix(ENV_PREFIX)

	_ = v.ReadInConfig()

	c.path = path.Join(dirname, CONFIG_FILE)
	c.viper = v
	c.fs = fs
}

func (c *Config) BuiltInConfigs() (map[string]*core.Schema, error) {
	loggerConfigSchema, err := logger.ConfigSchema()
	if err != nil {
		return nil, fmt.Errorf("unable to get logger config schema: %w", err)
	}

	configMap := map[string]*core.Schema{
		"logging": loggerConfigSchema,
	}

	return configMap, nil
}

func (c *Config) Get(key string) any {
	return c.viper.Get(key)
}

func (c *Config) Set(key string, value interface{}) error {
	if err := os.MkdirAll(path.Dir(c.path), core.FILE_PERMISSION); err != nil {
		return fmt.Errorf("error creating dir at %s: %w", c.path, err)
	}
	c.viper.Set(key, value)

	if err := c.viper.WriteConfigAs(c.path); err != nil {
		return fmt.Errorf("error writing to config file: %w", err)
	}

	return nil
}

func (c *Config) Delete(key string) error {
	configMap := c.viper.AllSettings()
	if _, ok := configMap[key]; !ok {
		return nil
	}

	delete(configMap, key)

	if err := c.saveToConfigFile(configMap); err != nil {
		return err
	}

	return nil
}

func (c *Config) saveToConfigFile(configMap map[string]interface{}) error {
	encodedConfig, err := yaml.Marshal(configMap)
	if err != nil {
		return err
	}

	if err = c.fs.MkdirAll(path.Dir(c.path), core.FILE_PERMISSION); err != nil {
		return fmt.Errorf("error creating dir at %s: %w", c.path, err)
	}

	if err = afero.WriteFile(c.fs, c.path, encodedConfig, core.FILE_PERMISSION); err != nil {
		return fmt.Errorf("error writing to config file: %w", err)
	}

	if err := c.viper.ReadInConfig(); err != nil {
		return err
	}

	return nil
}
