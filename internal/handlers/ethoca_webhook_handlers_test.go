package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"mastercom-service/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// MockEthocaWebhookService is a mock implementation of the webhook service
type MockEthocaWebhookService struct {
	mock.Mock
}

func (m *MockEthocaWebhookService) ProcessWebhook(ctx context.Context, webhook *models.EthocaWebhook) (*models.OutcomeAcknowledgement, error) {
	args := m.Called(ctx, webhook)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.OutcomeAcknowledgement), args.Error(1)
}

func (m *MockEthocaWebhookService) GetWebhookConfig() *models.WebhookConfig {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*models.WebhookConfig)
}

// MockDatadogLogger is a mock implementation of the DatadogLogger
type MockDatadogLogger struct {
	*logrus.Logger
	mock.Mock
}

func NewMockDatadogLogger() *MockDatadogLogger {
	return &MockDatadogLogger{
		Logger: logrus.New(),
	}
}

func (m *MockDatadogLogger) InfoWithSpan(span tracer.Span, message string, fields logrus.Fields) {
	m.Called(span, message, fields)
}

func (m *MockDatadogLogger) ErrorWithSpan(span tracer.Span, message string, fields logrus.Fields) {
	m.Called(span, message, fields)
}

func (m *MockDatadogLogger) SetLevel(level logrus.Level) {
	m.Called(level)
}

func (m *MockDatadogLogger) Info(msg string, fields logrus.Fields) {
	m.Called(msg, fields)
}

func (m *MockDatadogLogger) Error(msg string, fields logrus.Fields) {
	m.Called(msg, fields)
}

func setupTestHandler() (*EthocaWebhookHandler, *MockEthocaWebhookService, *MockDatadogLogger) {
	mockService := &MockEthocaWebhookService{}
	mockLogger := NewMockDatadogLogger()

	handler := NewEthocaWebhookHandler(mockService, mockLogger)

	return handler, mockService, mockLogger
}

func setupGinContext() (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	return c, w
}

func TestNewEthocaWebhookHandler(t *testing.T) {
	mockService := &MockEthocaWebhookService{}
	mockLogger := NewMockDatadogLogger()

	handler := NewEthocaWebhookHandler(mockService, mockLogger)

	assert.NotNil(t, handler)
	assert.Equal(t, mockService, handler.webhookService)
	assert.Equal(t, mockLogger, handler.logger)
}

func TestHandleEthocaWebhook_Success(t *testing.T) {
	// Set up the global handler for testing
	mockService := &MockEthocaWebhookService{}
	mockLogger := NewMockDatadogLogger()

	// Set the global handler
	ethocaWebhookHandler = NewEthocaWebhookHandler(mockService, mockLogger)

	c, w := setupGinContext()

	// Create test webhook payload
	webhook := models.EthocaWebhook{
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

	payload, _ := json.Marshal(webhook)

	// Set up request
	c.Request = httptest.NewRequest("POST", "/api/v6/webhooks/ethoca", bytes.NewBuffer(payload))
	c.Request.Header.Set("Content-Type", "application/json")

	// Set up mock expectations
	mockLogger.On("InfoWithSpan", mock.Anything, "Received Ethoca webhook request", mock.Anything).Return()
	mockLogger.On("InfoWithSpan", mock.Anything, "Webhook payload validated successfully", mock.Anything).Return()
	mockLogger.On("InfoWithSpan", mock.Anything, "Webhook processed successfully", mock.Anything).Return()

	acknowledgment := &models.OutcomeAcknowledgement{
		OutcomeResponses: []models.StatusUpdate{
			{
				AlertID: "A4IM9K2MIYL9F2BPF9TWUIXTU",
				Status:  "SUCCESS",
			},
		},
	}

	mockService.On("ProcessWebhook", mock.Anything, &webhook).Return(acknowledgment, nil)

	// Call the handler function
	HandleEthocaWebhook(c)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "SUCCESS", response["status"])
	assert.NotEmpty(t, w.Header().Get("X-Request-ID"))

	mockService.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestHandleEthocaWebhook_InvalidMethod(t *testing.T) {
	// Set up the global handler for testing
	mockService := &MockEthocaWebhookService{}
	mockLogger := NewMockDatadogLogger()

	// Set the global handler
	ethocaWebhookHandler = NewEthocaWebhookHandler(mockService, mockLogger)

	c, w := setupGinContext()

	// Set up GET request (invalid method)
	c.Request = httptest.NewRequest("GET", "/api/v6/webhooks/ethoca", nil)

	// Set up mock expectations
	mockLogger.On("InfoWithSpan", mock.Anything, "Received Ethoca webhook request", mock.Anything).Return()
	mockLogger.On("ErrorWithSpan", mock.Anything, "Invalid HTTP method for webhook", mock.Anything).Return()

	// Call the handler function
	HandleEthocaWebhook(c)

	// Assertions
	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "Method not allowed", response["error"])
	assert.Equal(t, "METHOD_NOT_ALLOWED", response["code"])

	mockLogger.AssertExpectations(t)
}

func TestHandleEthocaWebhook_InvalidContentType(t *testing.T) {
	// Set up the global handler for testing
	mockService := &MockEthocaWebhookService{}
	mockLogger := NewMockDatadogLogger()

	// Set the global handler
	ethocaWebhookHandler = NewEthocaWebhookHandler(mockService, mockLogger)

	c, w := setupGinContext()

	// Set up request with invalid content type
	c.Request = httptest.NewRequest("POST", "/api/v6/webhooks/ethoca", nil)
	c.Request.Header.Set("Content-Type", "text/plain")

	// Set up mock expectations
	mockLogger.On("InfoWithSpan", mock.Anything, "Received Ethoca webhook request", mock.Anything).Return()
	mockLogger.On("ErrorWithSpan", mock.Anything, "Invalid content type for webhook", mock.Anything).Return()

	// Call the handler function
	HandleEthocaWebhook(c)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "Invalid content type. Expected application/json", response["error"])
	assert.Equal(t, "INVALID_CONTENT_TYPE", response["code"])

	mockLogger.AssertExpectations(t)
}

func TestHandleEthocaWebhook_InvalidJSON(t *testing.T) {
	// Set up the global handler for testing
	mockService := &MockEthocaWebhookService{}
	mockLogger := NewMockDatadogLogger()

	// Set the global handler
	ethocaWebhookHandler = NewEthocaWebhookHandler(mockService, mockLogger)

	c, w := setupGinContext()

	// Set up request with invalid JSON
	c.Request = httptest.NewRequest("POST", "/api/v6/webhooks/ethoca", bytes.NewBufferString("invalid json"))
	c.Request.Header.Set("Content-Type", "application/json")

	// Set up mock expectations
	mockLogger.On("InfoWithSpan", mock.Anything, "Received Ethoca webhook request", mock.Anything).Return()
	mockLogger.On("ErrorWithSpan", mock.Anything, "Failed to parse webhook payload", mock.Anything).Return()

	// Call the handler function
	HandleEthocaWebhook(c)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "Invalid JSON payload", response["error"])
	assert.Equal(t, "INVALID_JSON", response["code"])

	mockLogger.AssertExpectations(t)
}

func TestHandleEthocaWebhook_ServiceError(t *testing.T) {
	// Set up the global handler for testing
	mockService := &MockEthocaWebhookService{}
	mockLogger := NewMockDatadogLogger()

	// Set the global handler
	ethocaWebhookHandler = NewEthocaWebhookHandler(mockService, mockLogger)

	c, w := setupGinContext()

	// Create test webhook payload
	webhook := models.EthocaWebhook{
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

	payload, _ := json.Marshal(webhook)

	// Set up request
	c.Request = httptest.NewRequest("POST", "/api/v6/webhooks/ethoca", bytes.NewBuffer(payload))
	c.Request.Header.Set("Content-Type", "application/json")

	// Set up mock expectations
	mockLogger.On("InfoWithSpan", mock.Anything, "Received Ethoca webhook request", mock.Anything).Return()
	mockLogger.On("InfoWithSpan", mock.Anything, "Webhook payload validated successfully", mock.Anything).Return()
	mockLogger.On("ErrorWithSpan", mock.Anything, "Webhook processing failed", mock.Anything).Return()

	mockService.On("ProcessWebhook", mock.Anything, &webhook).Return(nil, assert.AnError)

	// Call the handler function
	HandleEthocaWebhook(c)

	// Assertions
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "Internal server error", response["error"])
	assert.Equal(t, "PROCESSING_ERROR", response["code"])

	mockService.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestGetWebhookHealth(t *testing.T) {
	// Set up the global handler for testing
	mockService := &MockEthocaWebhookService{}
	mockLogger := NewMockDatadogLogger()

	// Set the global handler
	ethocaWebhookHandler = NewEthocaWebhookHandler(mockService, mockLogger)

	c, w := setupGinContext()

	// Set up mock expectations
	config := &models.WebhookConfig{
		Endpoint:   "/api/v6/webhooks/ethoca",
		SecretKey:  "test-secret",
		Timeout:    30,
		MaxRetries: 3,
		BatchSize:  25,
	}

	mockService.On("GetWebhookConfig").Return(config)

	// Call the handler function
	GetWebhookHealth(c)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "healthy", response["status"])
	assert.Equal(t, "ethoca-webhook", response["service"])
	assert.Equal(t, "/api/v6/webhooks/ethoca", response["endpoint"])
	assert.Equal(t, float64(30), response["timeout"])
	assert.Equal(t, float64(3), response["maxRetries"])
	assert.Equal(t, float64(25), response["batchSize"])
	assert.NotNil(t, response["timestamp"])

	mockService.AssertExpectations(t)
}

func TestGetWebhookStats(t *testing.T) {
	// Set up the global handler for testing
	mockService := &MockEthocaWebhookService{}
	mockLogger := NewMockDatadogLogger()

	// Set the global handler
	ethocaWebhookHandler = NewEthocaWebhookHandler(mockService, mockLogger)

	c, w := setupGinContext()

	// Call the handler function
	GetWebhookStats(c)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "ok", response["status"])
	assert.NotNil(t, response["stats"])
	assert.NotNil(t, response["timestamp"])

	stats := response["stats"].(map[string]interface{})
	assert.Equal(t, float64(0), stats["totalWebhooks"])
	assert.Equal(t, float64(0), stats["successfulWebhooks"])
	assert.Equal(t, float64(0), stats["failedWebhooks"])
	assert.Equal(t, "0ms", stats["averageProcessingTime"])
	assert.NotNil(t, stats["lastProcessedAt"])
}

func TestHandleEthocaWebhook_EmptyOutcomes(t *testing.T) {
	// Set up the global handler for testing
	mockService := &MockEthocaWebhookService{}
	mockLogger := NewMockDatadogLogger()

	// Set the global handler
	ethocaWebhookHandler = NewEthocaWebhookHandler(mockService, mockLogger)

	c, w := setupGinContext()

	// Create test webhook payload with empty outcomes
	webhook := models.EthocaWebhook{
		Outcomes: []models.AlertOutcome{},
	}

	payload, _ := json.Marshal(webhook)

	// Set up request
	c.Request = httptest.NewRequest("POST", "/api/v6/webhooks/ethoca", bytes.NewBuffer(payload))
	c.Request.Header.Set("Content-Type", "application/json")

	// Set up mock expectations
	mockLogger.On("InfoWithSpan", mock.Anything, "Received Ethoca webhook request", mock.Anything).Return()
	mockLogger.On("ErrorWithSpan", mock.Anything, "Webhook payload contains no outcomes", mock.Anything).Return()

	// Call the handler function
	HandleEthocaWebhook(c)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "No outcomes provided in webhook payload", response["error"])
	assert.Equal(t, "NO_OUTCOMES", response["code"])

	mockLogger.AssertExpectations(t)
}

func TestHandleEthocaWebhook_RequestIDGeneration(t *testing.T) {
	// Set up the global handler for testing
	mockService := &MockEthocaWebhookService{}
	mockLogger := NewMockDatadogLogger()

	// Set the global handler
	ethocaWebhookHandler = NewEthocaWebhookHandler(mockService, mockLogger)

	c, w := setupGinContext()

	// Create test webhook payload
	webhook := models.EthocaWebhook{
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

	payload, _ := json.Marshal(webhook)

	// Set up request
	c.Request = httptest.NewRequest("POST", "/api/v6/webhooks/ethoca", bytes.NewBuffer(payload))
	c.Request.Header.Set("Content-Type", "application/json")

	// Set up mock expectations
	mockLogger.On("InfoWithSpan", mock.Anything, "Received Ethoca webhook request", mock.Anything).Return()
	mockLogger.On("InfoWithSpan", mock.Anything, "Webhook payload validated successfully", mock.Anything).Return()
	mockLogger.On("InfoWithSpan", mock.Anything, "Webhook processed successfully", mock.Anything).Return()

	acknowledgment := &models.OutcomeAcknowledgement{
		OutcomeResponses: []models.StatusUpdate{
			{
				AlertID: "A4IM9K2MIYL9F2BPF9TWUIXTU",
				Status:  "SUCCESS",
			},
		},
	}

	mockService.On("ProcessWebhook", mock.Anything, &webhook).Return(acknowledgment, nil)

	// Call the handler function
	HandleEthocaWebhook(c)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	// Check that request ID is generated and included in headers
	requestID := w.Header().Get("X-Request-ID")
	assert.NotEmpty(t, requestID)
	assert.Len(t, requestID, 36) // UUID length

	mockService.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}
