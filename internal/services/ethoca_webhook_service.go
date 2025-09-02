package services

import (
	"context"
	"fmt"
	"time"

	"mastercom-service/internal/models"
	"mastercom-service/pkg/logger"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// EthocaWebhookService handles processing of Ethoca webhook events
type EthocaWebhookService struct {
	logger *logger.DatadogLogger
	config *models.WebhookConfig
}

// NewEthocaWebhookService creates a new webhook service instance
func NewEthocaWebhookService(logger *logger.DatadogLogger, config *models.WebhookConfig) *EthocaWebhookService {
	return &EthocaWebhookService{
		logger: logger,
		config: config,
	}
}

// ProcessWebhook processes incoming webhook data and returns acknowledgment
func (s *EthocaWebhookService) ProcessWebhook(ctx context.Context, webhook *models.EthocaWebhook) (*models.OutcomeAcknowledgement, error) {
	s.logger.Info("Processing Ethoca webhook", logrus.Fields{
		"outcomeCount": len(webhook.Outcomes),
		"requestId":    ctx.Value("requestId"),
	})

	var statusUpdates []models.StatusUpdate

	for _, outcome := range webhook.Outcomes {
		statusUpdate, err := s.processOutcome(ctx, &outcome)
		if err != nil {
			s.logger.Error("Failed to process outcome", logrus.Fields{
				"alertId": outcome.AlertID,
				"error":   err.Error(),
			})

			// Create error status update
			statusUpdate = models.StatusUpdate{
				AlertID: outcome.AlertID,
				Status:  "FAILURE",
				Errors: &models.Errors{
					Error: []models.Error{
						{
							Source:      stringPtr("Service"),
							ReasonCode:  stringPtr("PROCESSING_ERROR"),
							Description: stringPtr("Failed to process outcome"),
							Recoverable: boolPtr(true),
							Details:     stringPtr(err.Error()),
						},
					},
				},
			}
		}

		statusUpdates = append(statusUpdates, statusUpdate)
	}

	acknowledgment := &models.OutcomeAcknowledgement{
		OutcomeResponses: statusUpdates,
	}

	s.logger.Info("Webhook processing completed", logrus.Fields{
		"processedCount": len(statusUpdates),
		"requestId":      ctx.Value("requestId"),
	})

	return acknowledgment, nil
}

// processOutcome processes a single alert outcome
func (s *EthocaWebhookService) processOutcome(ctx context.Context, outcome *models.AlertOutcome) (models.StatusUpdate, error) {
	// Create webhook event for tracking
	webhookEvent := &models.WebhookEvent{
		ID:           uuid.New().String(),
		AlertID:      outcome.AlertID,
		Outcome:      outcome.Outcome,
		RefundStatus: outcome.RefundStatus,
		Amount:       outcome.Refund.Amount.Value,
		Currency:     outcome.Refund.Amount.CurrencyCode,
		Comments:     outcome.Comments,
		ProcessedAt:  time.Now(),
		Status:       "PROCESSING",
	}

	s.logger.Info("Processing alert outcome", logrus.Fields{
		"alertId":      outcome.AlertID,
		"outcome":      outcome.Outcome,
		"refundStatus": outcome.RefundStatus,
		"amount":       outcome.Refund.Amount.Value,
		"currency":     outcome.Refund.Amount.CurrencyCode,
		"eventId":      webhookEvent.ID,
	})

	// Validate outcome data
	if err := s.validateOutcome(outcome); err != nil {
		webhookEvent.Status = "FAILED"
		webhookEvent.ErrorMessage = stringPtr(err.Error())
		s.logger.Error("Outcome validation failed", logrus.Fields{
			"alertId": outcome.AlertID,
			"error":   err.Error(),
		})
		return models.StatusUpdate{}, err
	}

	// Process based on outcome type
	switch outcome.Outcome {
	case "STOPPED", "PARTIALLY_STOPPED":
		if err := s.processFraudOutcome(ctx, outcome); err != nil {
			webhookEvent.Status = "FAILED"
			webhookEvent.ErrorMessage = stringPtr(err.Error())
			return models.StatusUpdate{}, err
		}
	case "RESOLVED", "RESOLVED_PREVIOUSLY_REFUNDED":
		if err := s.processDisputeOutcome(ctx, outcome); err != nil {
			webhookEvent.Status = "FAILED"
			webhookEvent.ErrorMessage = stringPtr(err.Error())
			return models.StatusUpdate{}, err
		}
	default:
		if err := s.processOtherOutcome(ctx, outcome); err != nil {
			webhookEvent.Status = "FAILED"
			webhookEvent.ErrorMessage = stringPtr(err.Error())
			return models.StatusUpdate{}, err
		}
	}

	webhookEvent.Status = "SUCCESS"
	s.logger.Info("Outcome processed successfully", logrus.Fields{
		"alertId": outcome.AlertID,
		"eventId": webhookEvent.ID,
	})

	// TODO: Store webhook event in database for audit trail
	// s.storeWebhookEvent(ctx, webhookEvent)

	return models.StatusUpdate{
		AlertID: outcome.AlertID,
		Status:  "SUCCESS",
	}, nil
}

// validateOutcome validates the outcome data
func (s *EthocaWebhookService) validateOutcome(outcome *models.AlertOutcome) error {
	// Basic validation is handled by struct tags, but we can add business logic here
	if outcome.RefundStatus == "REFUNDED" && outcome.Refund.Amount.Value <= 0 {
		return fmt.Errorf("refund amount must be greater than 0 when refund status is REFUNDED")
	}

	if outcome.Outcome == "STOPPED" && outcome.AmountStopped.Value <= 0 {
		return fmt.Errorf("amount stopped must be greater than 0 when outcome is STOPPED")
	}

	return nil
}

// processFraudOutcome processes fraud-related outcomes
func (s *EthocaWebhookService) processFraudOutcome(ctx context.Context, outcome *models.AlertOutcome) error {
	s.logger.Info("Processing fraud outcome", logrus.Fields{
		"alertId":       outcome.AlertID,
		"outcome":       outcome.Outcome,
		"amountStopped": outcome.AmountStopped.Value,
		"currency":      outcome.AmountStopped.CurrencyCode,
	})

	// TODO: Implement fraud outcome processing logic
	// - Update fraud database
	// - Notify fraud team
	// - Update case status
	// - Generate reports

	return nil
}

// processDisputeOutcome processes customer dispute outcomes
func (s *EthocaWebhookService) processDisputeOutcome(ctx context.Context, outcome *models.AlertOutcome) error {
	s.logger.Info("Processing dispute outcome", logrus.Fields{
		"alertId":      outcome.AlertID,
		"outcome":      outcome.Outcome,
		"refundAmount": outcome.Refund.Amount.Value,
		"currency":     outcome.Refund.Amount.CurrencyCode,
	})

	// TODO: Implement dispute outcome processing logic
	// - Update dispute database
	// - Process refund if applicable
	// - Update customer records
	// - Generate settlement reports

	return nil
}

// processOtherOutcome processes other types of outcomes
func (s *EthocaWebhookService) processOtherOutcome(ctx context.Context, outcome *models.AlertOutcome) error {
	s.logger.Info("Processing other outcome", logrus.Fields{
		"alertId": outcome.AlertID,
		"outcome": outcome.Outcome,
	})

	// TODO: Implement other outcome processing logic
	// - Log for review
	// - Route to appropriate team
	// - Update case status

	return nil
}

// GetWebhookConfig returns the webhook configuration
func (s *EthocaWebhookService) GetWebhookConfig() *models.WebhookConfig {
	return s.config
}

// Helper functions for creating pointers
func stringPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}
