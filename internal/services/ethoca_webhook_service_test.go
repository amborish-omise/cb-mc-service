package services

import (
	"context"
	"testing"

	"mastercom-service/internal/models"
	"mastercom-service/pkg/logger"

	"github.com/stretchr/testify/assert"
)

func TestNewEthocaWebhookService(t *testing.T) {
	logger := logger.NewDatadogLogger()
	config := &models.WebhookConfig{
		Endpoint:   "/test",
		SecretKey:  "test-key",
		Timeout:    30,
		MaxRetries: 3,
		BatchSize:  25,
	}

	service := NewEthocaWebhookService(logger, config)

	assert.NotNil(t, service)
	assert.Equal(t, logger, service.logger)
	assert.Equal(t, config, service.config)
}

func TestProcessWebhook_Success(t *testing.T) {
	logger := logger.NewDatadogLogger()
	config := &models.WebhookConfig{
		Endpoint:   "/test",
		SecretKey:  "test-key",
		Timeout:    30,
		MaxRetries: 3,
		BatchSize:  25,
	}

	service := NewEthocaWebhookService(logger, config)

	webhook := &models.EthocaWebhook{
		Outcomes: []models.AlertOutcome{
			{
				AlertID:      "A4IM9K2MIYL9F2BPF9TWUIXTU",
				Outcome:      "STOPPED",
				RefundStatus: "NOT_REFUNDED",
				Refund: models.Refund{
					Amount: models.RefundAmount{
						Value:        100.00,
						CurrencyCode: "USD",
					},
					Timestamp: "2021-06-18T22:11:05+05:00",
				},
				AmountStopped: models.AmountStopped{
					Value:        100.00,
					CurrencyCode: "USD",
				},
			},
		},
	}

	ctx := context.Background()
	acknowledgment, err := service.ProcessWebhook(ctx, webhook)

	assert.NoError(t, err)
	assert.NotNil(t, acknowledgment)
	assert.Len(t, acknowledgment.OutcomeResponses, 1)
	assert.Equal(t, "A4IM9K2MIYL9F2BPF9TWUIXTU", acknowledgment.OutcomeResponses[0].AlertID)
	assert.Equal(t, "SUCCESS", acknowledgment.OutcomeResponses[0].Status)
}

func TestProcessWebhook_ValidationError(t *testing.T) {
	logger := logger.NewDatadogLogger()
	config := &models.WebhookConfig{
		Endpoint:   "/test",
		SecretKey:  "test-key",
		Timeout:    30,
		MaxRetries: 3,
		BatchSize:  25,
	}

	service := NewEthocaWebhookService(logger, config)

	// Invalid webhook - refund status is REFUNDED but amount is 0
	webhook := &models.EthocaWebhook{
		Outcomes: []models.AlertOutcome{
			{
				AlertID:      "A4IM9K2MIYL9F2BPF9TWUIXTU",
				Outcome:      "RESOLVED",
				RefundStatus: "REFUNDED",
				Refund: models.Refund{
					Amount: models.RefundAmount{
						Value:        0.00, // Invalid: refund amount must be > 0 when status is REFUNDED
						CurrencyCode: "USD",
					},
					Timestamp: "2021-06-18T22:11:05+05:00",
				},
				AmountStopped: models.AmountStopped{
					Value:        100.00,
					CurrencyCode: "USD",
				},
			},
		},
	}

	ctx := context.Background()
	acknowledgment, err := service.ProcessWebhook(ctx, webhook)

	assert.NoError(t, err) // Service handles errors gracefully
	assert.NotNil(t, acknowledgment)
	assert.Len(t, acknowledgment.OutcomeResponses, 1)
	assert.Equal(t, "A4IM9K2MIYL9F2BPF9TWUIXTU", acknowledgment.OutcomeResponses[0].AlertID)
	assert.Equal(t, "FAILURE", acknowledgment.OutcomeResponses[0].Status)
	assert.NotNil(t, acknowledgment.OutcomeResponses[0].Errors)
}

func TestProcessWebhook_EmptyOutcomes(t *testing.T) {
	logger := logger.NewDatadogLogger()
	config := &models.WebhookConfig{
		Endpoint:   "/test",
		SecretKey:  "test-key",
		Timeout:    30,
		MaxRetries: 3,
		BatchSize:  25,
	}

	service := NewEthocaWebhookService(logger, config)

	webhook := &models.EthocaWebhook{
		Outcomes: []models.AlertOutcome{}, // Empty outcomes
	}

	ctx := context.Background()
	acknowledgment, err := service.ProcessWebhook(ctx, webhook)

	assert.NoError(t, err)
	assert.NotNil(t, acknowledgment)
	assert.Len(t, acknowledgment.OutcomeResponses, 0)
}

func TestProcessWebhook_MultipleOutcomes(t *testing.T) {
	logger := logger.NewDatadogLogger()
	config := &models.WebhookConfig{
		Endpoint:   "/test",
		SecretKey:  "test-key",
		Timeout:    30,
		MaxRetries: 3,
		BatchSize:  25,
	}

	service := NewEthocaWebhookService(logger, config)

	webhook := &models.EthocaWebhook{
		Outcomes: []models.AlertOutcome{
			{
				AlertID:      "A4IM9K2MIYL9F2BPF9TWUIXTU",
				Outcome:      "STOPPED",
				RefundStatus: "NOT_REFUNDED",
				Refund: models.Refund{
					Amount: models.RefundAmount{
						Value:        100.00,
						CurrencyCode: "USD",
					},
					Timestamp: "2021-06-18T22:11:05+05:00",
				},
				AmountStopped: models.AmountStopped{
					Value:        100.00,
					CurrencyCode: "USD",
				},
			},
			{
				AlertID:      "B5JN0L3NJZM0G3CQG0UXVJYUV",
				Outcome:      "RESOLVED",
				RefundStatus: "REFUNDED",
				Refund: models.Refund{
					Amount: models.RefundAmount{
						Value:        50.00,
						CurrencyCode: "EUR",
					},
					Timestamp: "2021-06-18T22:11:05+05:00",
				},
				AmountStopped: models.AmountStopped{
					Value:        50.00,
					CurrencyCode: "EUR",
				},
			},
		},
	}

	ctx := context.Background()
	acknowledgment, err := service.ProcessWebhook(ctx, webhook)

	assert.NoError(t, err)
	assert.NotNil(t, acknowledgment)
	assert.Len(t, acknowledgment.OutcomeResponses, 2)

	// Check first outcome
	assert.Equal(t, "A4IM9K2MIYL9F2BPF9TWUIXTU", acknowledgment.OutcomeResponses[0].AlertID)
	assert.Equal(t, "SUCCESS", acknowledgment.OutcomeResponses[0].Status)

	// Check second outcome
	assert.Equal(t, "B5JN0L3NJZM0G3CQG0UXVJYUV", acknowledgment.OutcomeResponses[1].AlertID)
	assert.Equal(t, "SUCCESS", acknowledgment.OutcomeResponses[1].Status)
}

func TestProcessWebhook_DisputeOutcome(t *testing.T) {
	logger := logger.NewDatadogLogger()
	config := &models.WebhookConfig{
		Endpoint:   "/test",
		SecretKey:  "test-key",
		Timeout:    30,
		MaxRetries: 3,
		BatchSize:  25,
	}

	service := NewEthocaWebhookService(logger, config)

	webhook := &models.EthocaWebhook{
		Outcomes: []models.AlertOutcome{
			{
				AlertID:      "C6KO1M4OKAN1H4DRH1VYWKZWV",
				Outcome:      "RESOLVED",
				RefundStatus: "REFUNDED",
				Refund: models.Refund{
					Amount: models.RefundAmount{
						Value:        75.00,
						CurrencyCode: "GBP",
					},
					Timestamp: "2021-06-18T22:11:05+05:00",
				},
				AmountStopped: models.AmountStopped{
					Value:        75.00,
					CurrencyCode: "GBP",
				},
				Comments: stringPtr("Customer dispute resolved"),
			},
		},
	}

	ctx := context.Background()
	acknowledgment, err := service.ProcessWebhook(ctx, webhook)

	assert.NoError(t, err)
	assert.NotNil(t, acknowledgment)
	assert.Len(t, acknowledgment.OutcomeResponses, 1)
	assert.Equal(t, "C6KO1M4OKAN1H4DRH1VYWKZWV", acknowledgment.OutcomeResponses[0].AlertID)
	assert.Equal(t, "SUCCESS", acknowledgment.OutcomeResponses[0].Status)
}

func TestProcessWebhook_OtherOutcome(t *testing.T) {
	logger := logger.NewDatadogLogger()
	config := &models.WebhookConfig{
		Endpoint:   "/test",
		SecretKey:  "test-key",
		Timeout:    30,
		MaxRetries: 3,
		BatchSize:  25,
	}

	service := NewEthocaWebhookService(logger, config)

	webhook := &models.EthocaWebhook{
		Outcomes: []models.AlertOutcome{
			{
				AlertID:      "D7LP2N5PLBO2I5ESJ2WZXLAXW",
				Outcome:      "INVESTIGATING",
				RefundStatus: "PENDING",
				Refund: models.Refund{
					Amount: models.RefundAmount{
						Value:        25.00,
						CurrencyCode: "CAD",
					},
					Timestamp: "2021-06-18T22:11:05+05:00",
				},
				AmountStopped: models.AmountStopped{
					Value:        25.00,
					CurrencyCode: "CAD",
				},
				ActionTimestamp: stringPtr("2021-06-18T22:11:05+05:00"),
			},
		},
	}

	ctx := context.Background()
	acknowledgment, err := service.ProcessWebhook(ctx, webhook)

	assert.NoError(t, err)
	assert.NotNil(t, acknowledgment)
	assert.Len(t, acknowledgment.OutcomeResponses, 1)
	assert.Equal(t, "D7LP2N5PLBO2I5ESJ2WZXLAXW", acknowledgment.OutcomeResponses[0].AlertID)
	assert.Equal(t, "SUCCESS", acknowledgment.OutcomeResponses[0].Status)
}

func TestValidateOutcome_ValidOutcome(t *testing.T) {
	logger := logger.NewDatadogLogger()
	config := &models.WebhookConfig{
		Endpoint:   "/test",
		SecretKey:  "test-key",
		Timeout:    30,
		MaxRetries: 3,
		BatchSize:  25,
	}

	service := NewEthocaWebhookService(logger, config)

	outcome := &models.AlertOutcome{
		AlertID:      "A4IM9K2MIYL9F2BPF9TWUIXTU",
		Outcome:      "STOPPED",
		RefundStatus: "NOT_REFUNDED",
		Refund: models.Refund{
			Amount: models.RefundAmount{
				Value:        100.00,
				CurrencyCode: "USD",
			},
			Timestamp: "2021-06-18T22:11:05+05:00",
		},
		AmountStopped: models.AmountStopped{
			Value:        100.00,
			CurrencyCode: "USD",
		},
	}

	err := service.validateOutcome(outcome)
	assert.NoError(t, err)
}

func TestValidateOutcome_StoppedWithZeroAmount(t *testing.T) {
	logger := logger.NewDatadogLogger()
	config := &models.WebhookConfig{
		Endpoint:   "/test",
		SecretKey:  "test-key",
		Timeout:    30,
		MaxRetries: 3,
		BatchSize:  25,
	}

	service := NewEthocaWebhookService(logger, config)

	outcome := &models.AlertOutcome{
		AlertID:      "A4IM9K2MIYL9F2BPF9TWUIXTU",
		Outcome:      "STOPPED",
		RefundStatus: "NOT_REFUNDED",
		Refund: models.Refund{
			Amount: models.RefundAmount{
				Value:        100.00,
				CurrencyCode: "USD",
			},
			Timestamp: "2021-06-18T22:11:05+05:00",
		},
		AmountStopped: models.AmountStopped{
			Value:        0.00, // Invalid: amount stopped must be > 0 when outcome is STOPPED
			CurrencyCode: "USD",
		},
	}

	err := service.validateOutcome(outcome)
	assert.Error(t, err)
}

func TestValidateOutcome_ValidOutcomeWithNonStopped(t *testing.T) {
	logger := logger.NewDatadogLogger()
	config := &models.WebhookConfig{
		Endpoint:   "/test",
		SecretKey:  "test-key",
		Timeout:    30,
		MaxRetries: 3,
		BatchSize:  25,
	}

	service := NewEthocaWebhookService(logger, config)

	outcome := &models.AlertOutcome{
		AlertID:      "A4IM9K2MIYL9F2BPF9TWUIXTU",
		Outcome:      "RESOLVED", // Not STOPPED, so amount stopped validation won't apply
		RefundStatus: "REFUNDED",
		Refund: models.Refund{
			Amount: models.RefundAmount{
				Value:        100.00,
				CurrencyCode: "USD",
			},
			Timestamp: "2021-06-18T22:11:05+05:00",
		},
		AmountStopped: models.AmountStopped{
			Value:        0.00, // This should be valid since outcome is not STOPPED
			CurrencyCode: "USD",
		},
	}

	err := service.validateOutcome(outcome)
	assert.NoError(t, err)
}

func TestValidateOutcome_RefundedWithZeroAmount(t *testing.T) {
	logger := logger.NewDatadogLogger()
	config := &models.WebhookConfig{
		Endpoint:   "/test",
		SecretKey:  "test-key",
		Timeout:    30,
		MaxRetries: 3,
		BatchSize:  25,
	}

	service := NewEthocaWebhookService(logger, config)

	outcome := &models.AlertOutcome{
		AlertID:      "A4IM9K2MIYL9F2BPF9TWUIXTU",
		Outcome:      "RESOLVED",
		RefundStatus: "REFUNDED",
		Refund: models.Refund{
			Amount: models.RefundAmount{
				Value:        0.00, // Invalid: refund amount must be > 0 when status is REFUNDED
				CurrencyCode: "USD",
			},
			Timestamp: "2021-06-18T22:11:05+05:00",
		},
		AmountStopped: models.AmountStopped{
			Value:        100.00,
			CurrencyCode: "USD",
		},
	}

	err := service.validateOutcome(outcome)
	assert.Error(t, err)
}

func TestGetWebhookConfig(t *testing.T) {
	logger := logger.NewDatadogLogger()
	config := &models.WebhookConfig{
		Endpoint:   "/test",
		SecretKey:  "test-key",
		Timeout:    30,
		MaxRetries: 3,
		BatchSize:  25,
	}

	service := NewEthocaWebhookService(logger, config)

	retrievedConfig := service.GetWebhookConfig()

	assert.Equal(t, config, retrievedConfig)
	assert.Equal(t, "/test", retrievedConfig.Endpoint)
	assert.Equal(t, "test-key", retrievedConfig.SecretKey)
	assert.Equal(t, 30, retrievedConfig.Timeout)
	assert.Equal(t, 3, retrievedConfig.MaxRetries)
	assert.Equal(t, 25, retrievedConfig.BatchSize)
}
