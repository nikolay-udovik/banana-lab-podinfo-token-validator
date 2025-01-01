package redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type RedisClient struct {
	client *redis.Client
}

// NewRedisClient initializes a new Redis client.
func NewRedisClient(host string, port int, password string) (*RedisClient, error) {
	options := &redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: password, // Set empty string if no password
		DB:       0,        // Default DB
	}

	rdb := redis.NewClient(options)

	// Test connection
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisClient{client: rdb}, nil
}

// KeyExists checks if a key exists in Redis.
func (r *RedisClient) KeyExists(key string) (bool, error) {
	exists, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check key existence: %w", err)
	}
	return exists > 0, nil
}

// GetValue retrieves the value of a key from Redis.
func (r *RedisClient) GetValue(key string) (string, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("key %s does not exist", key)
	} else if err != nil {
		return "", fmt.Errorf("failed to get value for key %s: %w", key, err)
	}
	return val, nil
}

// Close closes the Redis client connection.
func (r *RedisClient) Close() error {
	return r.client.Close()
}
