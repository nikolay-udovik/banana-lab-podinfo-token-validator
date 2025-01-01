package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// Config structure holds all configuration sections.
type Config struct {
	Log     LogConfig
	Podinfo PodinfoConfig
	Redis   RedisConfig
}

// LogConfig holds logging configuration.
type LogConfig struct {
	Level  string
	Format string
}

// PodinfoConfig holds configuration specific to the Podinfo app.
type PodinfoConfig struct {
	BaseURL           string `mapstructure:"base_url"`       // Base URL for Podinfo
	TokenEndpoint     string `mapstructure:"token_endpoint"` // Endpoint to generate a token
	TokenValidatePath string `mapstructure:"token_validate"` // Endpoint to validate the token
	CacheEndpoint     string `mapstructure:"cache_endpoint"` // Endpoint to store validation results in cache
}

// RedisConfig holds Redis connection details.
type RedisConfig struct {
	Host                string `mapstructure:"host"`                  // Redis hostname
	Port                int    `mapstructure:"port"`                  // Redis port
	Password            string `mapstructure:"password"`              // Redis password
	ValidationResultKey string `mapstructure:"validation_result_key"` // Cache key for validation results
}

// AppConfigInstance is the global instance of the configuration.
var AppConfigInstance *Config

// LoadConfig loads the configuration from a file and environment variables.
func LoadConfig(configFilePath string) (*Config, error) {
	if AppConfigInstance != nil {
		return AppConfigInstance, nil
	}

	// Set default config file path if not provided
	if configFilePath == "" {
		configFilePath = "./internal/config/config.yaml"
	}

	// Check if the config file exists
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("configuration file not found: %s", configFilePath)
	}

	// Set config file and automatic environment variable binding
	viper.SetConfigFile(configFilePath)
	viper.AutomaticEnv()

	// Read the config file
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	// Unmarshal config into the struct
	config := &Config{}
	if err := viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	// Bind environment variables
	viper.BindEnv("log.level", "LOG_LEVEL")
	viper.BindEnv("log.format", "LOG_FORMAT")

	// Podinfo configuration
	viper.BindEnv("podinfo.base_url", "PODINFO_BASE_URL")
	viper.BindEnv("podinfo.token_endpoint", "PODINFO_TOKEN_ENDPOINT")
	viper.BindEnv("podinfo.token_validate", "PODINFO_TOKEN_VALIDATE_PATH")
	viper.BindEnv("podinfo.cache_endpoint", "PODINFO_CACHE_ENDPOINT")

	// Redis configuration
	viper.BindEnv("redis.host", "REDIS_HOST")
	viper.BindEnv("redis.port", "REDIS_PORT")
	viper.BindEnv("redis.password", "REDIS_PASSWORD")
	viper.BindEnv("redis.validation_result_key", "REDIS_VALIDATION_RESULT_KEY")

	// Re-marshal the environment variable overrides into the config struct
	if err := viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("unable to decode into struct after binding env variables: %w", err)
	}

	AppConfigInstance = config
	return AppConfigInstance, nil
}
