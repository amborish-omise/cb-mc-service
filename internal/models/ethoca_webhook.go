package models

import "time"

// EthocaWebhook represents the main webhook payload structure
type EthocaWebhook struct {
	Outcomes []AlertOutcome `json:"outcomes" validate:"required,dive"`
}

// AlertOutcome represents a single alert outcome
type AlertOutcome struct {
	AlertID         string        `json:"alertId" validate:"required,min=25,max=25"`
	Outcome         string        `json:"outcome" validate:"required,min=5,max=30"`
	RefundStatus    string        `json:"refundStatus" validate:"required,min=8,max=12"`
	Refund          Refund        `json:"refund" validate:"required"`
	AmountStopped   AmountStopped `json:"amountStopped" validate:"required"`
	Comments        *string       `json:"comments,omitempty" validate:"omitempty,min=1,max=1024"`
	ActionTimestamp *string       `json:"actionTimestamp,omitempty" validate:"omitempty,min=10,max=25"`
}

// Refund represents refund information
type Refund struct {
	Amount                  RefundAmount `json:"amount" validate:"required"`
	Type                    *string      `json:"type,omitempty" validate:"omitempty,min=6,max=9"`
	Timestamp               string       `json:"timestamp" validate:"required,min=10,max=25"`
	TransactionID           *string      `json:"transactionId,omitempty" validate:"omitempty,min=1,max=64"`
	AcquirerReferenceNumber *string      `json:"acquirerReferenceNumber,omitempty" validate:"omitempty,min=1,max=24"`
}

// RefundAmount represents the refund amount and currency
type RefundAmount struct {
	Value        float64 `json:"value" validate:"required,min=1,max=999999"`
	CurrencyCode string  `json:"currencyCode" validate:"required,min=3,max=3"`
}

// AmountStopped represents the amount stopped due to fraud
type AmountStopped struct {
	Value        float64 `json:"value" validate:"required,min=1,max=999999"`
	CurrencyCode string  `json:"currencyCode" validate:"required,min=3,max=3"`
}

// OutcomeAcknowledgement represents the response to a webhook submission
type OutcomeAcknowledgement struct {
	OutcomeResponses []StatusUpdate `json:"outcomeResponses"`
}

// StatusUpdate represents the status of a processed outcome
type StatusUpdate struct {
	AlertID string  `json:"alertId" validate:"required"`
	Status  string  `json:"status" validate:"required"`
	Errors  *Errors `json:"errors,omitempty"`
}

// Errors represents error information
type Errors struct {
	Error []Error `json:"Error" validate:"required"`
}

// Error represents a single error
type Error struct {
	Source      *string `json:"Source,omitempty" validate:"omitempty,max=100"`
	ReasonCode  *string `json:"ReasonCode,omitempty" validate:"omitempty,max=100"`
	Description *string `json:"Description,omitempty" validate:"omitempty,max=1000"`
	Recoverable *bool   `json:"Recoverable,omitempty"`
	Details     *string `json:"Details,omitempty" validate:"omitempty,max=1000"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Errors Errors `json:"Errors" validate:"required"`
}

// WebhookEvent represents a processed webhook event for logging/tracking
type WebhookEvent struct {
	ID           string    `json:"id"`
	AlertID      string    `json:"alertId"`
	Outcome      string    `json:"outcome"`
	RefundStatus string    `json:"refundStatus"`
	Amount       float64   `json:"amount"`
	Currency     string    `json:"currency"`
	Comments     *string   `json:"comments,omitempty"`
	ProcessedAt  time.Time `json:"processedAt"`
	Status       string    `json:"status"`
	ErrorMessage *string   `json:"errorMessage,omitempty"`
}

// WebhookConfig represents configuration for the webhook endpoint
type WebhookConfig struct {
	Endpoint   string `json:"endpoint"`
	SecretKey  string `json:"secretKey"`
	Timeout    int    `json:"timeout"`
	MaxRetries int    `json:"maxRetries"`
	BatchSize  int    `json:"batchSize"`
}
