package config

import (
	"fmt"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

// ColumnSpec defines a single column in the CSV schema.
type ColumnSpec struct {
	Name     string `mapstructure:"name"`
	Type     string `mapstructure:"type"`     // supported: string, int, float, bool
	Required bool   `mapstructure:"required"` // true if the column must be present
}

// Config holds the global configuration.
type Config struct {
	LogLevel string       `mapstructure:"log_level"`
	Schema   []ColumnSpec `mapstructure:"schema"`
}

var (
	once   sync.Once
	global *Config
)

// Load reads configuration from file and environment variables.
// It is safe to call multiple times; the config is loaded only once.
func Load() (*Config, error) {
	var err error
	once.Do(func() {
		v := viper.New()
		// Environment variable support
		v.SetEnvPrefix("CSV2JSON")
		v.AutomaticEnv()
		v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

		// Config file location
		configPath := v.GetString("CONFIG")
		if configPath == "" {
			configPath = "./config.yaml"
		}
		v.SetConfigFile(configPath)

		if readErr := v.ReadInConfig(); readErr != nil {
			err = fmt.Errorf("failed to read config file %s: %w", configPath, readErr)
			return
		}

		var cfg Config
		if decodeErr := v.Unmarshal(&cfg); decodeErr != nil {
			err = fmt.Errorf("failed to decode config: %w", decodeErr)
			return
		}
		global = &cfg
	})
	return global, err
}