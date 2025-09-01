package models

import (
	"time"

	"github.com/google/uuid"
)

// Case represents a MasterCom case
type Case struct {
	ID                    string    `json:"id" bson:"_id"`
	CaseType              string    `json:"caseType" validate:"required"`
	CaseTypeDescription   string    `json:"caseTypeDescription"`
	PrimaryAccountNumber  string    `json:"primaryAccountNumber" validate:"required"`
	TransactionAmount     float64   `json:"transactionAmount" validate:"required"`
	TransactionCurrency   string    `json:"transactionCurrency" validate:"required"`
	TransactionDate       time.Time `json:"transactionDate" validate:"required"`
	TransactionID         string    `json:"transactionId" validate:"required"`
	MerchantName          string    `json:"merchantName"`
	MerchantCategoryCode  string    `json:"merchantCategoryCode"`
	ReasonCode            string    `json:"reasonCode" validate:"required"`
	ReasonDescription     string    `json:"reasonDescription"`
	DisputeAmount         float64   `json:"disputeAmount"`
	DisputeCurrency       string    `json:"disputeCurrency"`
	FilingAs              string    `json:"filingAs" validate:"required"`
	FilingIca             string    `json:"filingIca" validate:"required"`
	FiledAgainstIca       string    `json:"filedAgainstIca" validate:"required"`
	FiledBy               string    `json:"filedBy"`
	FiledByContactName    string    `json:"filedByContactName"`
	FiledByContactPhone   string    `json:"filedByContactPhone"`
	FiledByContactEmail   string    `json:"filedByContactEmail"`
	Status                string    `json:"status"`
	CreatedAt             time.Time `json:"createdAt"`
	UpdatedAt             time.Time `json:"updatedAt"`
	Documents             []Document `json:"documents,omitempty"`
}

// CreateCaseRequest represents the request to create a new case
type CreateCaseRequest struct {
	CaseType              string    `json:"caseType" validate:"required"`
	PrimaryAccountNumber  string    `json:"primaryAccountNumber" validate:"required"`
	TransactionAmount     float64   `json:"transactionAmount" validate:"required"`
	TransactionCurrency   string    `json:"transactionCurrency" validate:"required"`
	TransactionDate       time.Time `json:"transactionDate" validate:"required"`
	TransactionID         string    `json:"transactionId" validate:"required"`
	MerchantName          string    `json:"merchantName"`
	MerchantCategoryCode  string    `json:"merchantCategoryCode"`
	ReasonCode            string    `json:"reasonCode" validate:"required"`
	DisputeAmount         float64   `json:"disputeAmount"`
	DisputeCurrency       string    `json:"disputeCurrency"`
	FilingAs              string    `json:"filingAs" validate:"required"`
	FilingIca             string    `json:"filingIca" validate:"required"`
	FiledAgainstIca       string    `json:"filedAgainstIca" validate:"required"`
	FiledBy               string    `json:"filedBy"`
	FiledByContactName    string    `json:"filedByContactName"`
	FiledByContactPhone   string    `json:"filedByContactPhone"`
	FiledByContactEmail   string    `json:"filedByContactEmail"`
}

// CaseResponse represents the response for case operations
type CaseResponse struct {
	ID                    string    `json:"id"`
	CaseType              string    `json:"caseType"`
	CaseTypeDescription   string    `json:"caseTypeDescription"`
	PrimaryAccountNumber  string    `json:"primaryAccountNumber"`
	TransactionAmount     float64   `json:"transactionAmount"`
	TransactionCurrency   string    `json:"transactionCurrency"`
	TransactionDate       time.Time `json:"transactionDate"`
	TransactionID         string    `json:"transactionId"`
	MerchantName          string    `json:"merchantName"`
	MerchantCategoryCode  string    `json:"merchantCategoryCode"`
	ReasonCode            string    `json:"reasonCode"`
	ReasonDescription     string    `json:"reasonDescription"`
	DisputeAmount         float64   `json:"disputeAmount"`
	DisputeCurrency       string    `json:"disputeCurrency"`
	FilingAs              string    `json:"filingAs"`
	FilingIca             string    `json:"filingIca"`
	FiledAgainstIca       string    `json:"filedAgainstIca"`
	FiledBy               string    `json:"filedBy"`
	FiledByContactName    string    `json:"filedByContactName"`
	FiledByContactPhone   string    `json:"filedByContactPhone"`
	FiledByContactEmail   string    `json:"filedByContactEmail"`
	Status                string    `json:"status"`
	CreatedAt             time.Time `json:"createdAt"`
	UpdatedAt             time.Time `json:"updatedAt"`
}

// NewCase creates a new case from a create request
func NewCase(req *CreateCaseRequest) *Case {
	now := time.Now()
	return &Case{
		ID:                    uuid.New().String(),
		CaseType:              req.CaseType,
		PrimaryAccountNumber:  req.PrimaryAccountNumber,
		TransactionAmount:     req.TransactionAmount,
		TransactionCurrency:   req.TransactionCurrency,
		TransactionDate:       req.TransactionDate,
		TransactionID:         req.TransactionID,
		MerchantName:          req.MerchantName,
		MerchantCategoryCode:  req.MerchantCategoryCode,
		ReasonCode:            req.ReasonCode,
		DisputeAmount:         req.DisputeAmount,
		DisputeCurrency:       req.DisputeCurrency,
		FilingAs:              req.FilingAs,
		FilingIca:             req.FilingIca,
		FiledAgainstIca:       req.FiledAgainstIca,
		FiledBy:               req.FiledBy,
		FiledByContactName:    req.FiledByContactName,
		FiledByContactPhone:   req.FiledByContactPhone,
		FiledByContactEmail:   req.FiledByContactEmail,
		Status:                "PENDING",
		CreatedAt:             now,
		UpdatedAt:             now,
		Documents:             []Document{},
	}
}
