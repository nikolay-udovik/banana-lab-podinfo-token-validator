package main

import (
	"flag"
	"fmt"
	"os"

	"monkale.io/podinfo-token-validator/internal/config"
	"monkale.io/podinfo-token-validator/internal/logger"
	"monkale.io/podinfo-token-validator/internal/podinfo"
	"monkale.io/podinfo-token-validator/internal/redis"
)

func main() {
	// Define a CLI flag for the config file path
	configFilePath := flag.String("config", "", "Path to the configuration file")
	flag.Parse()

	// Load configuration
	cfg, err := config.LoadConfig(*configFilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	if err := logger.InitLogger(&cfg.Log); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	// Log the configuration file path
	logger.Sugar.Infow("Using configuration file", "path", *configFilePath)

	// Step 1: Generate a token
	logger.Sugar.Info("Generating token...")
	token, err := podinfo.GenerateToken(&cfg.Podinfo)
	if err != nil {
		logger.Sugar.Fatalf("Error generating token: %v", err)
	}
	logger.Sugar.Infof("Token generated successfully.")

	// Step 2: Validate the token
	logger.Sugar.Info("Validating token...")
	expiresAt, tokenName, err := podinfo.ValidateToken(&cfg.Podinfo, token)
	if err != nil {
		logger.Sugar.Fatalf("Error validating token: %v", err)
	}
	logger.Sugar.Infof("Token validated: Name=%s, ExpiresAt=%s", tokenName, expiresAt)

	// Step 3: Cache validation result
	logger.Sugar.Info("Caching validation result...")
	if err := podinfo.CacheValidationResult(&cfg.Podinfo, token, true); err != nil {
		logger.Sugar.Fatalf("Error caching validation result: %v", err)
	}
	logger.Sugar.Info("Validation result cached successfully.")

	// Step 4: Verify cache in Redis
	logger.Sugar.Info("Verifying cached validation result in Redis...")
	rdb, err := redis.NewRedisClient(cfg.Redis.Host, cfg.Redis.Port, cfg.Redis.Password)
	if err != nil {
		logger.Sugar.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer rdb.Close()

	// Check if the key exists
	key := cfg.Redis.ValidationResultKey
	exists, err := rdb.KeyExists(key)
	if err != nil {
		logger.Sugar.Fatalf("Failed to check key existence: %v", err)
	}

	if !exists {
		logger.Sugar.Fatalf("Validation result key '%s' does not exist in Redis.", key)
	}
	logger.Sugar.Infof("Key '%s' exists in Redis.", key)

	// Retrieve the value
	_, err = rdb.GetValue(key) // Use `=` instead of `:=`
	if err != nil {
		logger.Sugar.Fatalf("Failed to retrieve value for key '%s': %v", key, err)
	}
	logger.Sugar.Infof("Cached validation result retrieved successfully for key '%s'.", key)
}
