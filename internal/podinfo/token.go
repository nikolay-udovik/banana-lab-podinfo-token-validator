package podinfo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"monkale.io/podinfo-token-validator/internal/config"
)

// GenerateToken calls the POST /token endpoint to generate a token.
func GenerateToken(cfg *config.PodinfoConfig) (string, error) {
	url := fmt.Sprintf("%s%s", cfg.BaseURL, cfg.TokenEndpoint)
	resp, err := http.Post(url, "application/json", nil)
	if err != nil {
		return "", fmt.Errorf("failed to call token endpoint: %w", err)
	}
	defer resp.Body.Close()

	var result map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to parse token response: %w", err)
	}
	return result["token"], nil
}

// ValidateToken calls the GET /token/validate endpoint to validate the token.
func ValidateToken(cfg *config.PodinfoConfig, token string) (string, string, error) {
	url := fmt.Sprintf("%s/token/validate", cfg.BaseURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("failed to call token validate endpoint: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("validation failed with status: %s", resp.Status)
	}

	var result struct {
		ExpiresAt string `json:"expires_at"`
		TokenName string `json:"token_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", "", fmt.Errorf("failed to parse validation response: %w", err)
	}
	return result.ExpiresAt, result.TokenName, nil
}

// CacheValidationResult calls POST/PUT /cache/validation_result to cache validation results.
func CacheValidationResult(cfg *config.PodinfoConfig, token string, valid bool) error {
	url := fmt.Sprintf("%s%s", cfg.BaseURL, cfg.CacheEndpoint)
	data := map[string]interface{}{
		"token": token,
		"valid": valid,
	}
	body, _ := json.Marshal(data)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to cache validation result: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("failed to cache validation result: %s", resp.Status)
	}

	return nil
}
