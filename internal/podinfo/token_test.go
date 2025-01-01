package podinfo

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"monkale.io/podinfo-token-validator/internal/config"
)

func TestGenerateToken(t *testing.T) {
	cfg, err := config.LoadConfig("../config/config.yaml")
	if err != nil {
		t.Fatalf("Failed to load configuration: %v", err)
	}

	token, err := GenerateToken(&cfg.Podinfo)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	assert.NotEmpty(t, token, "Token should not be empty")
}

func TestValidateToken(t *testing.T) {
	cfg, err := config.LoadConfig("../config/config.yaml")
	if err != nil {
		t.Fatalf("Failed to load configuration: %v", err)
	}

	token, err := GenerateToken(&cfg.Podinfo)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	expiresAt, tokenName, err := ValidateToken(&cfg.Podinfo, token)
	if err != nil {
		t.Fatalf("ValidateToken failed: %v", err)
	}

	assert.NotEmpty(t, expiresAt, "ExpiresAt should not be empty")
	assert.NotEmpty(t, tokenName, "TokenName should not be empty")
}

func TestCacheValidationResult(t *testing.T) {
	cfg, err := config.LoadConfig("../config/config.yaml")
	if err != nil {
		t.Fatalf("Failed to load configuration: %v", err)
	}

	token, err := GenerateToken(&cfg.Podinfo)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	err = CacheValidationResult(&cfg.Podinfo, token, true)
	assert.NoError(t, err, "CacheValidationResult should not return an error")
}
