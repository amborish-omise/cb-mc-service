package config

import (
	"strconv"

	"mastercom-service/internal/models"
)

// LoadEthocaConfig loads Ethoca webhook configuration from environment variables
func LoadEthocaConfig() *models.WebhookConfig {
	timeout, _ := strconv.Atoi(getEnv("ETHOCA_WEBHOOK_TIMEOUT", "30"))
	maxRetries, _ := strconv.Atoi(getEnv("ETHOCA_WEBHOOK_MAX_RETRIES", "3"))
	batchSize, _ := strconv.Atoi(getEnv("ETHOCA_WEBHOOK_BATCH_SIZE", "25"))

	return &models.WebhookConfig{
		Endpoint:   getEnv("ETHOCA_WEBHOOK_ENDPOINT", "/api/v6/webhooks/ethoca"),
		SecretKey:  getEnv("ETHOCA_WEBHOOK_SECRET_KEY", "your-secret-key-here"),
		Timeout:    timeout,
		MaxRetries: maxRetries,
		BatchSize:  batchSize,
	}
}
