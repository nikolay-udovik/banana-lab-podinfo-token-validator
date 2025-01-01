package redis

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"monkale.io/podinfo-token-validator/internal/config"
)

func TestRedisClient(t *testing.T) {
	// Load configuration
	cfg, err := config.LoadConfig("../config/config.yaml")
	assert.NoError(t, err, "Failed to load configuration")

	// Initialize Redis client using configuration
	rdb, err := NewRedisClient(cfg.Redis.Host, cfg.Redis.Port, cfg.Redis.Password)
	assert.NoError(t, err, "Failed to connect to Redis")
	defer func() {
		if rdb != nil {
			assert.NoError(t, rdb.Close(), "Failed to close Redis client")
		}
	}()

	// Key for testing
	testKey := "test_key"
	testValue := "test_value"

	// Set a value in Redis
	err = rdb.client.Set(ctx, testKey, testValue, 0).Err()
	assert.NoError(t, err, "Failed to set key in Redis")

	// Check if the key exists
	exists, err := rdb.KeyExists(testKey)
	assert.NoError(t, err, "Failed to check key existence")
	assert.True(t, exists, "Key does not exist in Redis")

	// Retrieve the value of the key
	val, err := rdb.GetValue(testKey)
	assert.NoError(t, err, "Failed to get key value")
	assert.Equal(t, testValue, val, "Value does not match expected")

	// Clean up: Delete the key
	err = rdb.client.Del(ctx, testKey).Err()
	assert.NoError(t, err, "Failed to delete test key from Redis")
}
