package config

import (
	"fmt"
	"path/filepath"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Use the correct absolute path to the config file for testing
	configFilePath, err := filepath.Abs("./config.yaml")
	if err != nil {
		t.Fatalf("Error finding config file path: %v", err)
	}

	// Test loading the config from the file
	cfg, err := LoadConfig(configFilePath)
	if err != nil {
		t.Fatalf("Error loading config: %v", err)
	}

	// Validate LogConfig
	fmt.Printf("Log Level: %s\n", cfg.Log.Level)
	fmt.Printf("Log Format: %s\n", cfg.Log.Format)

	// Validate PodinfoConfig
	fmt.Printf("Podinfo BaseURL: %s\n", cfg.Podinfo.BaseURL)
	fmt.Printf("Podinfo TokenEndpoint: %s\n", cfg.Podinfo.TokenEndpoint)
	fmt.Printf("Podinfo TokenValidatePath: %s\n", cfg.Podinfo.TokenValidatePath)
	fmt.Printf("Podinfo CacheEndpoint: %s\n", cfg.Podinfo.CacheEndpoint)

	// Validate RedisConfig
	fmt.Printf("Redis Host: %s\n", cfg.Redis.Host)
	fmt.Printf("Redis Port: %d\n", cfg.Redis.Port)
	fmt.Printf("Redis Password: %s\n", cfg.Redis.Password)
	fmt.Printf("Redis ValidationResultKey: %s\n", cfg.Redis.ValidationResultKey)
}
