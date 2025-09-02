package handlers

import (
	"context"
	"net/http"
	"time"

	"mastercom-service/internal/config"
	"mastercom-service/internal/models"
	"mastercom-service/internal/services"
	"mastercom-service/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// contextKey is a custom type for context keys to avoid collisions
type contextKey string

const (
	requestIDKey contextKey = "requestId"
)

type EthocaWebhookHandler struct {
	webhookService *services.EthocaWebhookService
	logger         *logger.DatadogLogger
}

// NewEthocaWebhookHandler creates a new webhook handler instance
func NewEthocaWebhookHandler(webhookService *services.EthocaWebhookService, logger *logger.DatadogLogger) *EthocaWebhookHandler {
	return &EthocaWebhookHandler{
		webhookService: webhookService,
		logger:         logger,
	}
}

// InitEthocaWebhookHandlers initializes the webhook service and handlers
func InitEthocaWebhookHandlers(logger *logger.DatadogLogger) {
	// Initialize webhook configuration
	config := config.LoadEthocaConfig()

	webhookService := services.NewEthocaWebhookService(logger, config)
	ethocaWebhookHandler = NewEthocaWebhookHandler(webhookService, logger)
}

var ethocaWebhookHandler *EthocaWebhookHandler

// HandleEthocaWebhook processes incoming Ethoca webhook requests
func HandleEthocaWebhook(c *gin.Context) {
	span := tracer.StartSpan("ethoca.webhook.process", tracer.ResourceName("EthocaWebhook"))
	defer span.Finish()

	// Generate request ID for tracking
	requestID := uuid.New().String()
	ctx := c.Request.Context()
	ctx = context.WithValue(ctx, requestIDKey, requestID)

	// Add request ID to response headers
	c.Header("X-Request-ID", requestID)

	ethocaWebhookHandler.logger.InfoWithSpan(span, "Received Ethoca webhook request", logrus.Fields{
		"requestId": requestID,
		"method":    c.Request.Method,
		"path":      c.Request.URL.Path,
		"userAgent": c.Request.UserAgent(),
		"remoteIP":  c.ClientIP(),
	})

	// Validate request method
	if c.Request.Method != http.MethodPost {
		ethocaWebhookHandler.logger.ErrorWithSpan(span, "Invalid HTTP method for webhook", logrus.Fields{
			"method":    c.Request.Method,
			"requestId": requestID,
		})
		c.JSON(http.StatusMethodNotAllowed, gin.H{
			"error": "Method not allowed",
			"code":  "METHOD_NOT_ALLOWED",
		})
		return
	}

	// Validate content type
	contentType := c.GetHeader("Content-Type")
	if contentType != "application/json" {
		ethocaWebhookHandler.logger.ErrorWithSpan(span, "Invalid content type for webhook", logrus.Fields{
			"contentType": contentType,
			"requestId":   requestID,
		})
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid content type. Expected application/json",
			"code":  "INVALID_CONTENT_TYPE",
		})
		return
	}

	// Parse and validate webhook payload
	var webhook models.EthocaWebhook
	if err := c.ShouldBindJSON(&webhook); err != nil {
		ethocaWebhookHandler.logger.ErrorWithSpan(span, "Failed to parse webhook payload", logrus.Fields{
			"error":     err.Error(),
			"requestId": requestID,
		})
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid JSON payload",
			"code":    "INVALID_JSON",
			"details": err.Error(),
		})
		return
	}

	// Validate webhook structure
	if len(webhook.Outcomes) == 0 {
		ethocaWebhookHandler.logger.ErrorWithSpan(span, "Webhook payload contains no outcomes", logrus.Fields{
			"requestId": requestID,
		})
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No outcomes provided in webhook payload",
			"code":  "NO_OUTCOMES",
		})
		return
	}

	if len(webhook.Outcomes) > 25 {
		ethocaWebhookHandler.logger.ErrorWithSpan(span, "Webhook payload exceeds maximum outcomes limit", logrus.Fields{
			"outcomeCount": len(webhook.Outcomes),
			"requestId":    requestID,
		})
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Maximum of 25 outcomes allowed per webhook",
			"code":  "TOO_MANY_OUTCOMES",
		})
		return
	}

	// Process webhook
	startTime := time.Now()
	acknowledgment, err := ethocaWebhookHandler.webhookService.ProcessWebhook(ctx, &webhook)
	processingTime := time.Since(startTime)

	if err != nil {
		ethocaWebhookHandler.logger.ErrorWithSpan(span, "Failed to process webhook", logrus.Fields{
			"error":          err.Error(),
			"requestId":      requestID,
			"processingTime": processingTime.String(),
		})
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
			"code":  "PROCESSING_ERROR",
		})
		return
	}

	// Log successful processing
	ethocaWebhookHandler.logger.InfoWithSpan(span, "Webhook processed successfully", logrus.Fields{
		"requestId":      requestID,
		"processingTime": processingTime.String(),
		"outcomeCount":   len(acknowledgment.OutcomeResponses),
	})

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"status":    "SUCCESS",
		"outcomes":  acknowledgment.OutcomeResponses,
		"requestId": requestID,
	})
}

// GetWebhookHealth returns the health status of the webhook service
func GetWebhookHealth(c *gin.Context) {
	config := ethocaWebhookHandler.webhookService.GetWebhookConfig()

	c.JSON(http.StatusOK, gin.H{
		"status":     "healthy",
		"service":    "ethoca-webhook",
		"endpoint":   config.Endpoint,
		"timeout":    config.Timeout,
		"maxRetries": config.MaxRetries,
		"batchSize":  config.BatchSize,
		"timestamp":  time.Now().UTC().Format(time.RFC3339),
	})
}

// GetWebhookStats returns statistics about webhook processing
func GetWebhookStats(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"stats": gin.H{
			"totalWebhooks":         0,
			"successfulWebhooks":    0,
			"failedWebhooks":        0,
			"averageProcessingTime": "0ms",
			"lastProcessedAt":       nil,
		},
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}
