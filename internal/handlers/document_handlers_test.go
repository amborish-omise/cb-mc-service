package handlers

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"mastercom-service/internal/models"
	"mastercom-service/internal/services"
	"mastercom-service/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupDocumentTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	// Initialize services and handlers
	logger := logger.NewDatadogLogger()
	documentService := services.NewDocumentService(logger)
	documentHandler := NewDocumentHandler(documentService, logger)
	
	// Setup routes
	api := router.Group("/api/v6")
	documents := api.Group("/documents")
	{
		documents.POST("", documentHandler.UploadDocument)
		documents.GET("/:id", documentHandler.GetDocument)
		documents.DELETE("/:id", documentHandler.DeleteDocument)
	}
	
	return router
}

func createMockMultipartRequest(caseID, fileName, content, description, uploadedBy string) (*http.Request, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	
	// Add file
	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return nil, err
	}
	part.Write([]byte(content))
	
	// Add form fields
	writer.WriteField("caseId", caseID)
	writer.WriteField("description", description)
	writer.WriteField("uploadedBy", uploadedBy)
	
	writer.Close()
	
	req, err := http.NewRequest("POST", "/api/v6/documents", body)
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, nil
}

func TestUploadDocument_Success(t *testing.T) {
	router := setupDocumentTestRouter()
	
	// Create mock multipart request
	req, err := createMockMultipartRequest(
		"test-case-id",
		"test-document.pdf",
		"This is test document content",
		"Test document description",
		"test-user",
	)
	require.NoError(t, err)
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusCreated, w.Code)
	
	var response models.Document
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	
	assert.NotEmpty(t, response.ID)
	assert.Equal(t, "test-case-id", response.CaseID)
	assert.Equal(t, "test-document.pdf", response.FileName)
	assert.Equal(t, ".pdf", response.FileType)
	assert.Equal(t, int64(29), response.FileSize) // Length of "This is test document content"
	assert.Equal(t, "test-user", response.UploadedBy)
	assert.Equal(t, "Test document description", response.Description)
}

func TestUploadDocument_MissingFile(t *testing.T) {
	router := setupDocumentTestRouter()
	
	// Create request without file
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("caseId", "test-case-id")
	writer.WriteField("description", "Test document")
	writer.WriteField("uploadedBy", "test-user")
	writer.Close()
	
	req, err := http.NewRequest("POST", "/api/v6/documents", body)
	require.NoError(t, err)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusBadRequest, w.Code)
	
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	
	assert.Contains(t, response, "error")
}

func TestUploadDocument_MissingCaseID(t *testing.T) {
	router := setupDocumentTestRouter()
	
	// Create request without caseId
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	
	// Add file
	part, err := writer.CreateFormFile("file", "test-document.pdf")
	require.NoError(t, err)
	part.Write([]byte("Test content"))
	
	writer.WriteField("description", "Test document")
	writer.WriteField("uploadedBy", "test-user")
	writer.Close()
	
	req, err := http.NewRequest("POST", "/api/v6/documents", body)
	require.NoError(t, err)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusBadRequest, w.Code)
	
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	
	assert.Contains(t, response, "error")
}

func TestGetDocument_Success(t *testing.T) {
	router := setupDocumentTestRouter()
	
	// First upload a document
	req, err := createMockMultipartRequest(
		"test-case-id",
		"test-document.pdf",
		"This is test document content",
		"Test document description",
		"test-user",
	)
	require.NoError(t, err)
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	var uploadedDoc models.Document
	json.Unmarshal(w.Body.Bytes(), &uploadedDoc)
	
	// Now get the document
	w = httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/api/v6/documents/"+uploadedDoc.ID, nil)
	router.ServeHTTP(w, request)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response models.Document
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	
	assert.Equal(t, uploadedDoc.ID, response.ID)
	assert.Equal(t, uploadedDoc.FileName, response.FileName)
	assert.Equal(t, uploadedDoc.FileSize, response.FileSize)
}

func TestGetDocument_NotFound(t *testing.T) {
	router := setupDocumentTestRouter()
	
	w := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/api/v6/documents/nonexistent-id", nil)
	router.ServeHTTP(w, request)
	
	assert.Equal(t, http.StatusNotFound, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	
	assert.Contains(t, response, "error")
}

func TestDeleteDocument_Success(t *testing.T) {
	router := setupDocumentTestRouter()
	
	// First upload a document
	req, err := createMockMultipartRequest(
		"test-case-id",
		"test-document.pdf",
		"This is test document content",
		"Test document description",
		"test-user",
	)
	require.NoError(t, err)
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	var uploadedDoc models.Document
	json.Unmarshal(w.Body.Bytes(), &uploadedDoc)
	
	// Now delete the document
	w = httptest.NewRecorder()
	request, _ := http.NewRequest("DELETE", "/api/v6/documents/"+uploadedDoc.ID, nil)
	router.ServeHTTP(w, request)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	
	assert.Contains(t, response, "message")
	
	// Verify document is deleted
	w = httptest.NewRecorder()
	request, _ = http.NewRequest("GET", "/api/v6/documents/"+uploadedDoc.ID, nil)
	router.ServeHTTP(w, request)
	
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteDocument_NotFound(t *testing.T) {
	router := setupDocumentTestRouter()
	
	w := httptest.NewRecorder()
	request, _ := http.NewRequest("DELETE", "/api/v6/documents/nonexistent-id", nil)
	router.ServeHTTP(w, request)
	
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestUploadDocument_DifferentFileTypes(t *testing.T) {
	router := setupDocumentTestRouter()
	
	testCases := []struct {
		fileName string
		content  string
		expectedExt string
	}{
		{"document.pdf", "PDF content", ".pdf"},
		{"image.jpg", "Image content", ".jpg"},
		{"data.xlsx", "Excel content", ".xlsx"},
		{"text.txt", "Text content", ".txt"},
	}
	
	for _, tc := range testCases {
		t.Run(tc.fileName, func(t *testing.T) {
			req, err := createMockMultipartRequest(
				"test-case-id",
				tc.fileName,
				tc.content,
				"Test document",
				"test-user",
			)
			require.NoError(t, err)
			
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			
			assert.Equal(t, http.StatusCreated, w.Code)
			
			var response models.Document
			err = json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)
			
			assert.Equal(t, tc.expectedExt, response.FileType)
			assert.Equal(t, int64(len(tc.content)), response.FileSize)
		})
	}
}

func TestUploadDocument_LargeFile(t *testing.T) {
	router := setupDocumentTestRouter()
	
	// Create a large content (simulating large file)
	largeContent := strings.Repeat("A", 1024*1024) // 1MB
	
	req, err := createMockMultipartRequest(
		"test-case-id",
		"large-file.pdf",
		largeContent,
		"Large test document",
		"test-user",
	)
	require.NoError(t, err)
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusCreated, w.Code)
	
	var response models.Document
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	
	assert.Equal(t, int64(1024*1024), response.FileSize)
}
