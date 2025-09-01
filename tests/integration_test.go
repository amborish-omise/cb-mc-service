package tests

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"mastercom-service/internal/handlers"
	"mastercom-service/internal/models"
	"mastercom-service/pkg/middleware"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupIntegrationTestServer() *gin.Engine {
	gin.SetMode(gin.TestMode)
	
	// Initialize logger
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	
	// Initialize handlers
	handlers.InitHandlers(logger)
	handlers.InitDocumentHandlers(logger)
	
	// Setup router
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.Logger(logger))
	router.Use(middleware.CORS())
	
	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "mastercom-service",
			"version": "v1.0.0",
		})
	})
	
	// API routes
	api := router.Group("/api/v6")
	{
		cases := api.Group("/cases")
		{
			cases.POST("", handlers.CreateCase)
			cases.GET("", handlers.ListCases)
			cases.GET("/:id", handlers.GetCase)
			cases.PUT("/:id", handlers.UpdateCase)
			cases.DELETE("/:id", handlers.DeleteCase)
		}
		
		documents := api.Group("/documents")
		{
			documents.POST("", handlers.UploadDocument)
			documents.GET("/:id", handlers.GetDocument)
			documents.DELETE("/:id", handlers.DeleteDocument)
		}
	}
	
	return router
}

func TestHealthEndpoint(t *testing.T) {
	router := setupIntegrationTestServer()
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	
	assert.Equal(t, "ok", response["status"])
	assert.Equal(t, "mastercom-service", response["service"])
	assert.Equal(t, "v1.0.0", response["version"])
}

func TestCompleteCaseWorkflow(t *testing.T) {
	router := setupIntegrationTestServer()
	
	// Step 1: Create a case
	caseReq := models.CreateCaseRequest{
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
	
	reqBody, _ := json.Marshal(caseReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v6/cases", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusCreated, w.Code)
	
	var createdCase models.Case
	err := json.Unmarshal(w.Body.Bytes(), &createdCase)
	require.NoError(t, err)
	assert.NotEmpty(t, createdCase.ID)
	assert.Equal(t, "PENDING", createdCase.Status)
	
	// Step 2: Upload a document for the case
	documentContent := "This is a test document for the case"
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	
	part, err := writer.CreateFormFile("file", "test-document.pdf")
	require.NoError(t, err)
	part.Write([]byte(documentContent))
	
	writer.WriteField("caseId", createdCase.ID)
	writer.WriteField("description", "Supporting documentation")
	writer.WriteField("uploadedBy", "test-user")
	writer.Close()
	
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v6/documents", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusCreated, w.Code)
	
	var uploadedDoc models.Document
	err = json.Unmarshal(w.Body.Bytes(), &uploadedDoc)
	require.NoError(t, err)
	assert.NotEmpty(t, uploadedDoc.ID)
	assert.Equal(t, createdCase.ID, uploadedDoc.CaseID)
	
	// Step 3: Get the case and verify it exists
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v6/cases/"+createdCase.ID, nil)
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var retrievedCase models.Case
	err = json.Unmarshal(w.Body.Bytes(), &retrievedCase)
	require.NoError(t, err)
	assert.Equal(t, createdCase.ID, retrievedCase.ID)
	assert.Equal(t, caseReq.CaseType, retrievedCase.CaseType)
	
	// Step 4: Update the case
	updateReq := caseReq
	updateReq.MerchantName = "Updated Merchant"
	updateReqBody, _ := json.Marshal(updateReq)
	
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/api/v6/cases/"+createdCase.ID, bytes.NewBuffer(updateReqBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var updatedCase models.Case
	err = json.Unmarshal(w.Body.Bytes(), &updatedCase)
	require.NoError(t, err)
	assert.Equal(t, "Updated Merchant", updatedCase.MerchantName)
	
	// Step 5: List cases and verify the case is there
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v6/cases", nil)
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var listResponse map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &listResponse)
	require.NoError(t, err)
	
	cases := listResponse["cases"].([]interface{})
	assert.Len(t, cases, 1)
	
	// Step 6: Get the document
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v6/documents/"+uploadedDoc.ID, nil)
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var retrievedDoc models.Document
	err = json.Unmarshal(w.Body.Bytes(), &retrievedDoc)
	require.NoError(t, err)
	assert.Equal(t, uploadedDoc.ID, retrievedDoc.ID)
	assert.Equal(t, "test-document.pdf", retrievedDoc.FileName)
	
	// Step 7: Delete the document
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/api/v6/documents/"+uploadedDoc.ID, nil)
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	// Step 8: Delete the case
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/api/v6/cases/"+createdCase.ID, nil)
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	// Step 9: Verify case is deleted
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v6/cases/"+createdCase.ID, nil)
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestErrorHandling(t *testing.T) {
	router := setupIntegrationTestServer()
	
	// Test invalid JSON
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v6/cases", strings.NewReader(`{"invalid": json}`))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusBadRequest, w.Code)
	
	// Test missing required fields
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v6/cases", strings.NewReader(`{"caseType": "PRE_ARBITRATION"}`))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusBadRequest, w.Code)
	
	// Test non-existent case
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v6/cases/nonexistent-id", nil)
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusNotFound, w.Code)
	
	// Test non-existent document
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v6/documents/nonexistent-id", nil)
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCORSHeaders(t *testing.T) {
	router := setupIntegrationTestServer()
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("OPTIONS", "/api/v6/cases", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
	assert.Contains(t, w.Header().Get("Access-Control-Allow-Headers"), "Content-Type")
	assert.Contains(t, w.Header().Get("Access-Control-Allow-Methods"), "POST")
}
