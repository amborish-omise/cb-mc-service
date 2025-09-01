package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"mastercom-service/internal/models"
	"mastercom-service/internal/services"
	"mastercom-service/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	// Initialize services and handlers
	logger := logger.NewDatadogLogger()
	caseService := services.NewCaseService(logger)
	caseHandler := NewCaseHandler(caseService, logger)
	
	// Setup routes
	api := router.Group("/api/v6")
	cases := api.Group("/cases")
	{
		cases.POST("", caseHandler.CreateCase)
		cases.GET("", caseHandler.ListCases)
		cases.GET("/:id", caseHandler.GetCase)
		cases.PUT("/:id", caseHandler.UpdateCase)
		cases.DELETE("/:id", caseHandler.DeleteCase)
	}
	
	return router
}

func createMockCaseRequest() models.CreateCaseRequest {
	return models.CreateCaseRequest{
		CaseType:              "PRE_ARBITRATION",
		PrimaryAccountNumber:  "4111111111111111",
		TransactionAmount:     100.00,
		TransactionCurrency:   "USD",
		TransactionDate:       time.Now(),
		TransactionID:         "123456789",
		MerchantName:          "Test Merchant",
		MerchantCategoryCode:  "5411",
		ReasonCode:            "10.1",
		DisputeAmount:         100.00,
		DisputeCurrency:       "USD",
		FilingAs:              "ISSUER",
		FilingIca:             "123456",
		FiledAgainstIca:       "654321",
		FiledBy:               "Test User",
		FiledByContactName:    "John Doe",
		FiledByContactPhone:   "+1234567890",
		FiledByContactEmail:   "john.doe@example.com",
	}
}

func TestCreateCase_Success(t *testing.T) {
	router := setupTestRouter()
	
	req := createMockCaseRequest()
	reqBody, _ := json.Marshal(req)
	
	w := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/api/v6/cases", bytes.NewBuffer(reqBody))
	request.Header.Set("Content-Type", "application/json")
	
	router.ServeHTTP(w, request)
	
	assert.Equal(t, http.StatusCreated, w.Code)
	
	var response models.Case
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	
	assert.NotEmpty(t, response.ID)
	assert.Equal(t, req.CaseType, response.CaseType)
	assert.Equal(t, req.PrimaryAccountNumber, response.PrimaryAccountNumber)
	assert.Equal(t, req.TransactionAmount, response.TransactionAmount)
	assert.Equal(t, "PENDING", response.Status)
}

func TestCreateCase_InvalidRequest(t *testing.T) {
	router := setupTestRouter()
	
	// Missing required fields
	req := models.CreateCaseRequest{
		CaseType: "PRE_ARBITRATION",
		// Missing other required fields
	}
	reqBody, _ := json.Marshal(req)
	
	w := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/api/v6/cases", bytes.NewBuffer(reqBody))
	request.Header.Set("Content-Type", "application/json")
	
	router.ServeHTTP(w, request)
	
	assert.Equal(t, http.StatusBadRequest, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	
	assert.Contains(t, response, "error")
}

func TestCreateCase_InvalidJSON(t *testing.T) {
	router := setupTestRouter()
	
	// Invalid JSON
	reqBody := []byte(`{"invalid": json}`)
	
	w := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/api/v6/cases", bytes.NewBuffer(reqBody))
	request.Header.Set("Content-Type", "application/json")
	
	router.ServeHTTP(w, request)
	
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetCase_Success(t *testing.T) {
	router := setupTestRouter()
	
	// First create a case
	req := createMockCaseRequest()
	reqBody, _ := json.Marshal(req)
	
	w := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/api/v6/cases", bytes.NewBuffer(reqBody))
	request.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, request)
	
	var createdCase models.Case
	json.Unmarshal(w.Body.Bytes(), &createdCase)
	
	// Now get the case
	w = httptest.NewRecorder()
	request, _ = http.NewRequest("GET", "/api/v6/cases/"+createdCase.ID, nil)
	router.ServeHTTP(w, request)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response models.Case
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	
	assert.Equal(t, createdCase.ID, response.ID)
	assert.Equal(t, req.CaseType, response.CaseType)
}

func TestGetCase_NotFound(t *testing.T) {
	router := setupTestRouter()
	
	w := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/api/v6/cases/nonexistent-id", nil)
	router.ServeHTTP(w, request)
	
	assert.Equal(t, http.StatusNotFound, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	
	assert.Contains(t, response, "error")
}

func TestListCases_Success(t *testing.T) {
	router := setupTestRouter()
	
	// Create multiple cases
	for i := 0; i < 3; i++ {
		req := createMockCaseRequest()
		req.TransactionID = "123456789" + string(rune(i+'0'))
		reqBody, _ := json.Marshal(req)
		
		w := httptest.NewRecorder()
		request, _ := http.NewRequest("POST", "/api/v6/cases", bytes.NewBuffer(reqBody))
		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, request)
	}
	
	// List cases
	w := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/api/v6/cases", nil)
	router.ServeHTTP(w, request)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	
	assert.Contains(t, response, "cases")
	assert.Contains(t, response, "total")
	assert.Contains(t, response, "page")
	assert.Contains(t, response, "limit")
	
	cases := response["cases"].([]interface{})
	assert.Len(t, cases, 3)
}

func TestListCases_WithPagination(t *testing.T) {
	router := setupTestRouter()
	
	// Create multiple cases
	for i := 0; i < 5; i++ {
		req := createMockCaseRequest()
		req.TransactionID = "123456789" + string(rune(i+'0'))
		reqBody, _ := json.Marshal(req)
		
		w := httptest.NewRecorder()
		request, _ := http.NewRequest("POST", "/api/v6/cases", bytes.NewBuffer(reqBody))
		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, request)
	}
	
	// List cases with pagination
	w := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/api/v6/cases?page=1&limit=2", nil)
	router.ServeHTTP(w, request)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	
	cases := response["cases"].([]interface{})
	assert.Len(t, cases, 2)
	assert.Equal(t, float64(5), response["total"])
	assert.Equal(t, float64(1), response["page"])
	assert.Equal(t, float64(2), response["limit"])
}

func TestUpdateCase_Success(t *testing.T) {
	router := setupTestRouter()
	
	// First create a case
	req := createMockCaseRequest()
	reqBody, _ := json.Marshal(req)
	
	w := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/api/v6/cases", bytes.NewBuffer(reqBody))
	request.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, request)
	
	var createdCase models.Case
	json.Unmarshal(w.Body.Bytes(), &createdCase)
	
	// Update the case
	updateReq := createMockCaseRequest()
	updateReq.MerchantName = "Updated Merchant"
	updateReqBody, _ := json.Marshal(updateReq)
	
	w = httptest.NewRecorder()
	request, _ = http.NewRequest("PUT", "/api/v6/cases/"+createdCase.ID, bytes.NewBuffer(updateReqBody))
	request.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, request)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response models.Case
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	
	assert.Equal(t, createdCase.ID, response.ID)
	assert.Equal(t, "Updated Merchant", response.MerchantName)
}

func TestUpdateCase_NotFound(t *testing.T) {
	router := setupTestRouter()
	
	req := createMockCaseRequest()
	reqBody, _ := json.Marshal(req)
	
	w := httptest.NewRecorder()
	request, _ := http.NewRequest("PUT", "/api/v6/cases/nonexistent-id", bytes.NewBuffer(reqBody))
	request.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, request)
	
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDeleteCase_Success(t *testing.T) {
	router := setupTestRouter()
	
	// First create a case
	req := createMockCaseRequest()
	reqBody, _ := json.Marshal(req)
	
	w := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/api/v6/cases", bytes.NewBuffer(reqBody))
	request.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, request)
	
	var createdCase models.Case
	json.Unmarshal(w.Body.Bytes(), &createdCase)
	
	// Delete the case
	w = httptest.NewRecorder()
	request, _ = http.NewRequest("DELETE", "/api/v6/cases/"+createdCase.ID, nil)
	router.ServeHTTP(w, request)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	
	assert.Contains(t, response, "message")
	
	// Verify case is deleted
	w = httptest.NewRecorder()
	request, _ = http.NewRequest("GET", "/api/v6/cases/"+createdCase.ID, nil)
	router.ServeHTTP(w, request)
	
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteCase_NotFound(t *testing.T) {
	router := setupTestRouter()
	
	w := httptest.NewRecorder()
	request, _ := http.NewRequest("DELETE", "/api/v6/cases/nonexistent-id", nil)
	router.ServeHTTP(w, request)
	
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
