package main

import (
	"testing"

	"monkale.io/podinfo-token-validator/internal/config"
	"monkale.io/podinfo-token-validator/internal/logger"
	"monkale.io/podinfo-token-validator/internal/podinfo"
	"monkale.io/podinfo-token-validator/internal/redis"

	"github.com/stretchr/testify/assert"
)

func TestMainFlow(t *testing.T) {
	// Load configuration
	cfg, err := config.LoadConfig("../internal/config/config.yaml")
	if err != nil {
		t.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	err = logger.InitLogger(&cfg.Log)
	if err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	// Step 1: Generate a token
	t.Log("Generating token...")
	token, err := podinfo.GenerateToken(&cfg.Podinfo)
	if err != nil {
		t.Fatalf("Error generating token: %v", err)
	}
	assert.NotEmpty(t, token, "Token should not be empty")

	// Step 2: Validate the token
	t.Log("Validating token...")
	expiresAt, tokenName, err := podinfo.ValidateToken(&cfg.Podinfo, token)
	if err != nil {
		t.Fatalf("Error validating token: %v", err)
	}
	assert.NotEmpty(t, expiresAt, "ExpiresAt should not be empty")
	assert.NotEmpty(t, tokenName, "TokenName should not be empty")

	// Step 3: Cache validation result
	t.Log("Caching validation result...")
	err = podinfo.CacheValidationResult(&cfg.Podinfo, token, true)
	assert.NoError(t, err, "Caching validation result should not return an error")

	// Step 4: Verify cache in Redis
	t.Log("Verifying cached validation result in Redis...")
	rdb, err := redis.NewRedisClient(cfg.Redis.Host, cfg.Redis.Port, cfg.Redis.Password)
	if err != nil {
		t.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer rdb.Close()

	key := cfg.Redis.ValidationResultKey
	exists, err := rdb.KeyExists(key)
	if err != nil {
		t.Fatalf("Failed to check key existence: %v", err)
	}
	assert.True(t, exists, "Validation result key should exist in Redis")

	val, err := rdb.GetValue(key)
	if err != nil {
		t.Fatalf("Failed to retrieve value for key '%s': %v", key, err)
	}
	assert.NotEmpty(t, val, "Cached validation result should not be empty")
}
