package handlers

import (
	"net/http"
	"path/filepath"

	"mastercom-service/internal/models"
	"mastercom-service/internal/services"
	"mastercom-service/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

type DocumentHandler struct {
	documentService *services.DocumentService
	logger          *logger.DatadogLogger
}

func NewDocumentHandler(documentService *services.DocumentService, logger *logger.DatadogLogger) *DocumentHandler {
	return &DocumentHandler{
		documentService: documentService,
		logger:          logger,
	}
}

// UploadDocument handles document upload
func (h *DocumentHandler) UploadDocument(c *gin.Context) {
	span := tracer.StartSpan("document.upload", tracer.ResourceName("UploadDocument"))
	defer span.Finish()

	file, err := c.FormFile("file")
	if err != nil {
		h.logger.ErrorWithSpan(span, "Failed to get uploaded file", logrus.Fields{"error": err.Error()})
		span.SetTag("error", true)
		span.SetTag("error.message", "No file uploaded")
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	caseID := c.PostForm("caseId")
	if caseID == "" {
		h.logger.ErrorWithSpan(span, "Missing case ID", logrus.Fields{})
		span.SetTag("error", true)
		span.SetTag("error.message", "Case ID is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Case ID is required"})
		return
	}

	description := c.PostForm("description")
	uploadedBy := c.PostForm("uploadedBy")

	span.SetTag("document.filename", file.Filename)
	span.SetTag("document.size", file.Size)
	span.SetTag("document.case_id", caseID)
	span.SetTag("document.uploaded_by", uploadedBy)

	// Read file content
	openedFile, err := file.Open()
	if err != nil {
		h.logger.ErrorWithSpan(span, "Failed to open uploaded file", logrus.Fields{"error": err.Error()})
		span.SetTag("error", true)
		span.SetTag("error.message", "Failed to process file")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process file"})
		return
	}
	defer openedFile.Close()

	// Read file content into byte slice
	content := make([]byte, file.Size)
	_, err = openedFile.Read(content)
	if err != nil {
		h.logger.ErrorWithSpan(span, "Failed to read file content", logrus.Fields{"error": err.Error()})
		span.SetTag("error", true)
		span.SetTag("error.message", "Failed to read file content")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file content"})
		return
	}

	// Create document
	document := models.NewDocument(
		caseID,
		file.Filename,
		filepath.Ext(file.Filename),
		content,
		uploadedBy,
		description,
	)

	if err := h.documentService.UploadDocument(document); err != nil {
		h.logger.ErrorWithSpan(span, "Failed to upload document", logrus.Fields{"error": err.Error()})
		span.SetTag("error", true)
		span.SetTag("error.message", "Failed to upload document")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload document"})
		return
	}

	h.logger.InfoWithSpan(span, "Document uploaded successfully", logrus.Fields{
		"documentId": document.ID,
		"filename": document.FileName,
		"fileSize": document.FileSize,
		"caseId": document.CaseID,
	})

	span.SetTag("document.id", document.ID)
	span.SetTag("document.file_type", document.FileType)

	c.JSON(http.StatusCreated, document)
}

// GetDocument handles retrieving a document
func (h *DocumentHandler) GetDocument(c *gin.Context) {
	documentID := c.Param("id")
	if documentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Document ID is required"})
		return
	}

	span := tracer.StartSpan("document.get", tracer.ResourceName("GetDocument"))
	defer span.Finish()

	span.SetTag("document.id", documentID)

	document, err := h.documentService.GetDocument(documentID)
	if err != nil {
		h.logger.ErrorWithSpan(span, "Failed to get document", logrus.Fields{
			"documentId": documentID,
			"error": err.Error(),
		})
		span.SetTag("error", true)
		span.SetTag("error.message", "Document not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
		return
	}

	h.logger.InfoWithSpan(span, "Document retrieved successfully", logrus.Fields{
		"documentId": documentID,
		"filename": document.FileName,
		"fileSize": document.FileSize,
	})

	c.JSON(http.StatusOK, document)
}

// DeleteDocument handles deleting a document
func (h *DocumentHandler) DeleteDocument(c *gin.Context) {
	documentID := c.Param("id")
	if documentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Document ID is required"})
		return
	}

	span := tracer.StartSpan("document.delete", tracer.ResourceName("DeleteDocument"))
	defer span.Finish()

	span.SetTag("document.id", documentID)

	if err := h.documentService.DeleteDocument(documentID); err != nil {
		h.logger.ErrorWithSpan(span, "Failed to delete document", logrus.Fields{
			"documentId": documentID,
			"error": err.Error(),
		})
		span.SetTag("error", true)
		span.SetTag("error.message", "Failed to delete document")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete document"})
		return
	}

	h.logger.InfoWithSpan(span, "Document deleted successfully", logrus.Fields{
		"documentId": documentID,
	})

	c.JSON(http.StatusOK, gin.H{"message": "Document deleted successfully"})
}

// Global handler functions for compatibility with main.go
var (
	documentService *services.DocumentService
	documentHandler *DocumentHandler
)

func InitDocumentHandlers(logger *logger.DatadogLogger) {
	documentService = services.NewDocumentService(logger)
	documentHandler = NewDocumentHandler(documentService, logger)
}

func UploadDocument(c *gin.Context) {
	if documentHandler == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Handler not initialized"})
		return
	}
	documentHandler.UploadDocument(c)
}

func GetDocument(c *gin.Context) {
	if documentHandler == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Handler not initialized"})
		return
	}
	documentHandler.GetDocument(c)
}

func DeleteDocument(c *gin.Context) {
	if documentHandler == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Handler not initialized"})
		return
	}
	documentHandler.DeleteDocument(c)
}
