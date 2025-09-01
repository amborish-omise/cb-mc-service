package services

import (
	"testing"
	"time"

	"mastercom-service/internal/models"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupCaseService() *CaseService {
	logger := logrus.New()
	return NewCaseService(logger)
}

func createMockCase() *models.Case {
	return &models.Case{
		ID:                    "test-case-id",
		CaseType:              "PRE_ARBITRATION",
		PrimaryAccountNumber:  "4111111111111111",
		TransactionAmount:     100.00,
		TransactionCurrency:   "USD",
		TransactionDate:       time.Now(),
		TransactionID:         "123456789",
		MerchantName:          "Test Merchant",
		ReasonCode:            "10.1",
		FilingAs:              "ISSUER",
		FilingIca:             "123456",
		FiledAgainstIca:       "654321",
		Status:                "PENDING",
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
	}
}

func TestCaseService_CreateCase(t *testing.T) {
	service := setupCaseService()
	caseObj := createMockCase()
	
	err := service.CreateCase(caseObj)
	require.NoError(t, err)
	
	// Verify case was created
	retrievedCase, err := service.GetCase(caseObj.ID)
	require.NoError(t, err)
	assert.Equal(t, caseObj.ID, retrievedCase.ID)
	assert.Equal(t, caseObj.CaseType, retrievedCase.CaseType)
}

func TestCaseService_CreateCase_Duplicate(t *testing.T) {
	service := setupCaseService()
	caseObj := createMockCase()
	
	// Create case first time
	err := service.CreateCase(caseObj)
	require.NoError(t, err)
	
	// Try to create same case again
	err = service.CreateCase(caseObj)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
}

func TestCaseService_GetCase(t *testing.T) {
	service := setupCaseService()
	caseObj := createMockCase()
	
	// Create case
	err := service.CreateCase(caseObj)
	require.NoError(t, err)
	
	// Get case
	retrievedCase, err := service.GetCase(caseObj.ID)
	require.NoError(t, err)
	assert.Equal(t, caseObj.ID, retrievedCase.ID)
	assert.Equal(t, caseObj.CaseType, retrievedCase.CaseType)
	assert.Equal(t, caseObj.PrimaryAccountNumber, retrievedCase.PrimaryAccountNumber)
}

func TestCaseService_GetCase_NotFound(t *testing.T) {
	service := setupCaseService()
	
	_, err := service.GetCase("nonexistent-id")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestCaseService_ListCases(t *testing.T) {
	service := setupCaseService()
	
	// Create multiple cases
	case1 := createMockCase()
	case1.ID = "case-1"
	case1.Status = "PENDING"
	
	case2 := createMockCase()
	case2.ID = "case-2"
	case2.Status = "RESOLVED"
	
	case3 := createMockCase()
	case3.ID = "case-3"
	case3.Status = "PENDING"
	
	err := service.CreateCase(case1)
	require.NoError(t, err)
	err = service.CreateCase(case2)
	require.NoError(t, err)
	err = service.CreateCase(case3)
	require.NoError(t, err)
	
	// List all cases
	cases, total, err := service.ListCases(1, 10, "")
	require.NoError(t, err)
	assert.Equal(t, 3, total)
	assert.Len(t, cases, 3)
	
	// List cases with status filter
	pendingCases, total, err := service.ListCases(1, 10, "PENDING")
	require.NoError(t, err)
	assert.Equal(t, 2, total)
	assert.Len(t, pendingCases, 2)
	
	for _, c := range pendingCases {
		assert.Equal(t, "PENDING", c.Status)
	}
}

func TestCaseService_ListCases_Pagination(t *testing.T) {
	service := setupCaseService()
	
	// Create 5 cases
	for i := 1; i <= 5; i++ {
		caseObj := createMockCase()
		caseObj.ID = "case-" + string(rune(i+'0'))
		err := service.CreateCase(caseObj)
		require.NoError(t, err)
	}
	
	// Test pagination - page 1, limit 2
	cases, total, err := service.ListCases(1, 2, "")
	require.NoError(t, err)
	assert.Equal(t, 5, total)
	assert.Len(t, cases, 2)
	
	// Test pagination - page 2, limit 2
	cases, total, err = service.ListCases(2, 2, "")
	require.NoError(t, err)
	assert.Equal(t, 5, total)
	assert.Len(t, cases, 2)
	
	// Test pagination - page 3, limit 2
	cases, total, err = service.ListCases(3, 2, "")
	require.NoError(t, err)
	assert.Equal(t, 5, total)
	assert.Len(t, cases, 1)
}

func TestCaseService_UpdateCase(t *testing.T) {
	service := setupCaseService()
	caseObj := createMockCase()
	
	// Create case
	err := service.CreateCase(caseObj)
	require.NoError(t, err)
	
	// Update case
	caseObj.MerchantName = "Updated Merchant"
	caseObj.Status = "RESOLVED"
	
	err = service.UpdateCase(caseObj)
	require.NoError(t, err)
	
	// Verify update
	retrievedCase, err := service.GetCase(caseObj.ID)
	require.NoError(t, err)
	assert.Equal(t, "Updated Merchant", retrievedCase.MerchantName)
	assert.Equal(t, "RESOLVED", retrievedCase.Status)
}

func TestCaseService_UpdateCase_NotFound(t *testing.T) {
	service := setupCaseService()
	caseObj := createMockCase()
	
	err := service.UpdateCase(caseObj)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestCaseService_DeleteCase(t *testing.T) {
	service := setupCaseService()
	caseObj := createMockCase()
	
	// Create case
	err := service.CreateCase(caseObj)
	require.NoError(t, err)
	
	// Delete case
	err = service.DeleteCase(caseObj.ID)
	require.NoError(t, err)
	
	// Verify deletion
	_, err = service.GetCase(caseObj.ID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestCaseService_DeleteCase_NotFound(t *testing.T) {
	service := setupCaseService()
	
	err := service.DeleteCase("nonexistent-id")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestCaseService_ConcurrentAccess(t *testing.T) {
	service := setupCaseService()
	
	// Test concurrent creation
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(id int) {
			caseObj := createMockCase()
			caseObj.ID = "concurrent-case-" + string(rune(id+'0'))
			err := service.CreateCase(caseObj)
			assert.NoError(t, err)
			done <- true
		}(i)
	}
	
	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
	
	// Verify all cases were created
	cases, total, err := service.ListCases(1, 20, "")
	require.NoError(t, err)
	assert.Equal(t, 10, total)
	assert.Len(t, cases, 10)
}
