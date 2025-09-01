package services

import (
	"testing"
	"time"

	"mastercom-service/internal/models"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupDocumentService() *DocumentService {
	logger := logrus.New()
	return NewDocumentService(logger)
}

func createMockDocument() *models.Document {
	return &models.Document{
		ID:          "test-doc-id",
		CaseID:      "test-case-id",
		FileName:    "test-document.pdf",
		FileType:    ".pdf",
		FileSize:    1024,
		Content:     []byte("Test document content"),
		UploadedBy:  "test-user",
		UploadedAt:  time.Now(),
		Description: "Test document description",
	}
}

func TestDocumentService_UploadDocument(t *testing.T) {
	service := setupDocumentService()
	document := createMockDocument()
	
	err := service.UploadDocument(document)
	require.NoError(t, err)
	
	// Verify document was uploaded
	retrievedDoc, err := service.GetDocument(document.ID)
	require.NoError(t, err)
	assert.Equal(t, document.ID, retrievedDoc.ID)
	assert.Equal(t, document.FileName, retrievedDoc.FileName)
	assert.Equal(t, document.FileSize, retrievedDoc.FileSize)
}

func TestDocumentService_UploadDocument_Duplicate(t *testing.T) {
	service := setupDocumentService()
	document := createMockDocument()
	
	// Upload document first time
	err := service.UploadDocument(document)
	require.NoError(t, err)
	
	// Try to upload same document again
	err = service.UploadDocument(document)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
}

func TestDocumentService_GetDocument(t *testing.T) {
	service := setupDocumentService()
	document := createMockDocument()
	
	// Upload document
	err := service.UploadDocument(document)
	require.NoError(t, err)
	
	// Get document
	retrievedDoc, err := service.GetDocument(document.ID)
	require.NoError(t, err)
	assert.Equal(t, document.ID, retrievedDoc.ID)
	assert.Equal(t, document.FileName, retrievedDoc.FileName)
	assert.Equal(t, document.FileType, retrievedDoc.FileType)
	assert.Equal(t, document.FileSize, retrievedDoc.FileSize)
	assert.Equal(t, document.Content, retrievedDoc.Content)
}

func TestDocumentService_GetDocument_NotFound(t *testing.T) {
	service := setupDocumentService()
	
	_, err := service.GetDocument("nonexistent-id")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestDocumentService_DeleteDocument(t *testing.T) {
	service := setupDocumentService()
	document := createMockDocument()
	
	// Upload document
	err := service.UploadDocument(document)
	require.NoError(t, err)
	
	// Delete document
	err = service.DeleteDocument(document.ID)
	require.NoError(t, err)
	
	// Verify deletion
	_, err = service.GetDocument(document.ID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestDocumentService_DeleteDocument_NotFound(t *testing.T) {
	service := setupDocumentService()
	
	err := service.DeleteDocument("nonexistent-id")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestDocumentService_GetDocumentsByCaseID(t *testing.T) {
	service := setupDocumentService()
	
	// Upload documents for different cases
	doc1 := createMockDocument()
	doc1.ID = "doc-1"
	doc1.CaseID = "case-1"
	doc1.FileName = "document1.pdf"
	
	doc2 := createMockDocument()
	doc2.ID = "doc-2"
	doc2.CaseID = "case-1"
	doc2.FileName = "document2.pdf"
	
	doc3 := createMockDocument()
	doc3.ID = "doc-3"
	doc3.CaseID = "case-2"
	doc3.FileName = "document3.pdf"
	
	err := service.UploadDocument(doc1)
	require.NoError(t, err)
	err = service.UploadDocument(doc2)
	require.NoError(t, err)
	err = service.UploadDocument(doc3)
	require.NoError(t, err)
	
	// Get documents for case-1
	documents, err := service.GetDocumentsByCaseID("case-1")
	require.NoError(t, err)
	assert.Len(t, documents, 2)
	
	// Verify correct documents
	fileNames := make(map[string]bool)
	for _, doc := range documents {
		fileNames[doc.FileName] = true
		assert.Equal(t, "case-1", doc.CaseID)
	}
	
	assert.True(t, fileNames["document1.pdf"])
	assert.True(t, fileNames["document2.pdf"])
	assert.False(t, fileNames["document3.pdf"])
	
	// Get documents for case-2
	documents, err = service.GetDocumentsByCaseID("case-2")
	require.NoError(t, err)
	assert.Len(t, documents, 1)
	assert.Equal(t, "document3.pdf", documents[0].FileName)
	
	// Get documents for non-existent case
	documents, err = service.GetDocumentsByCaseID("case-3")
	require.NoError(t, err)
	assert.Len(t, documents, 0)
}

func TestDocumentService_DifferentFileTypes(t *testing.T) {
	service := setupDocumentService()
	
	testCases := []struct {
		fileName string
		fileType string
		content  []byte
	}{
		{"document.pdf", ".pdf", []byte("PDF content")},
		{"image.jpg", ".jpg", []byte("Image content")},
		{"data.xlsx", ".xlsx", []byte("Excel content")},
		{"text.txt", ".txt", []byte("Text content")},
	}
	
	for _, tc := range testCases {
		t.Run(tc.fileName, func(t *testing.T) {
			document := createMockDocument()
			document.ID = "doc-" + tc.fileName
			document.FileName = tc.fileName
			document.FileType = tc.fileType
			document.Content = tc.content
			document.FileSize = int64(len(tc.content))
			
			err := service.UploadDocument(document)
			require.NoError(t, err)
			
			// Verify document
			retrievedDoc, err := service.GetDocument(document.ID)
			require.NoError(t, err)
			assert.Equal(t, tc.fileName, retrievedDoc.FileName)
			assert.Equal(t, tc.fileType, retrievedDoc.FileType)
			assert.Equal(t, tc.content, retrievedDoc.Content)
			assert.Equal(t, int64(len(tc.content)), retrievedDoc.FileSize)
		})
	}
}

func TestDocumentService_LargeFiles(t *testing.T) {
	service := setupDocumentService()
	
	// Create large content (1MB)
	largeContent := make([]byte, 1024*1024)
	for i := range largeContent {
		largeContent[i] = byte(i % 256)
	}
	
	document := createMockDocument()
	document.ID = "large-doc"
	document.FileName = "large-file.pdf"
	document.Content = largeContent
	document.FileSize = int64(len(largeContent))
	
	err := service.UploadDocument(document)
	require.NoError(t, err)
	
	// Verify large document
	retrievedDoc, err := service.GetDocument(document.ID)
	require.NoError(t, err)
	assert.Equal(t, int64(1024*1024), retrievedDoc.FileSize)
	assert.Equal(t, largeContent, retrievedDoc.Content)
}

func TestDocumentService_ConcurrentAccess(t *testing.T) {
	service := setupDocumentService()
	
	// Test concurrent uploads
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(id int) {
			document := createMockDocument()
			document.ID = "concurrent-doc-" + string(rune(id+'0'))
			document.FileName = "concurrent-doc-" + string(rune(id+'0')) + ".pdf"
			err := service.UploadDocument(document)
			assert.NoError(t, err)
			done <- true
		}(i)
	}
	
	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
	
	// Verify all documents were uploaded
	for i := 0; i < 10; i++ {
		docID := "concurrent-doc-" + string(rune(i+'0'))
		_, err := service.GetDocument(docID)
		assert.NoError(t, err)
	}
}
