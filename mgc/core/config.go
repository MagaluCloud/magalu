package core

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

type Config struct {
	path     string
	fileName string
}

const (
	CONFIG_NAME     = "config"
	CONFIG_TYPE     = "yaml"
	CONFIG_FOLDER   = ".mgc"
	CONFIG_FILE     = CONFIG_NAME + "." + CONFIG_TYPE
	FILE_PERMISSION = 0777 // TODO: investigate how to lower permission
)

var configKey contextKey = "magalu.cloud/core/Config"

func NewConfigContext(parent context.Context, config *Config) context.Context {
	return context.WithValue(parent, configKey, config)
}

func ConfigFromContext(ctx context.Context) *Config {
	c, _ := ctx.Value(configKey).(*Config)
	return c
}

func NewConfig() *Config {
	path, err := configFilePath()
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

func (c *Config) BindPFlag(key string, flag *pflag.Flag) error {
	return viper.BindPFlag(key, flag)
}

func (c *Config) IsSet(key string) bool {
	return viper.IsSet(key)
}

func (c *Config) Get(key string) any {
	return viper.Get(key)
}

func (c *Config) Set(key string, value interface{}) error {
	if err := os.MkdirAll(c.path, FILE_PERMISSION); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("error creating dir at %s: %w", c.path, err)
	}
	viper.Set(key, value)

	if err := viper.WriteConfigAs(path.Join(c.path, c.fileName)); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("error writing to config file: %w", err)
	}

	return nil
}

func (c *Config) Delete(key string) error {
	configMap := viper.AllSettings()
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

	if err = os.MkdirAll(c.path, FILE_PERMISSION); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("error creating dir at %s: %w", c.path, err)
	}

	if err = os.WriteFile(path.Join(c.path, c.fileName), encodedConfig, FILE_PERMISSION); err != nil {
		return fmt.Errorf("error writing to config file: %w", err)
	}

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return nil
}

func configFilePath() (string, error) {
	dir, err := os.UserHomeDir()
	if err == nil {
		return path.Join(dir, CONFIG_FOLDER), nil
	}
	dir, wdErr := os.Getwd()
	if dir == "" || wdErr != nil {
		return "", errors.Join(fmt.Errorf("Error: unable to access home and current working directory"), wdErr, err)
	}

	return path.Join(dir, CONFIG_FOLDER), nil
}
