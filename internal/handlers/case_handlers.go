package handlers

import (
	"net/http"
	"strconv"

	"mastercom-service/internal/models"
	"mastercom-service/internal/services"
	"mastercom-service/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"github.com/sirupsen/logrus"
)

type CaseHandler struct {
	caseService *services.CaseService
	validator   *validator.Validate
	logger      *logger.DatadogLogger
}

func NewCaseHandler(caseService *services.CaseService, logger *logger.DatadogLogger) *CaseHandler {
	return &CaseHandler{
		caseService: caseService,
		validator:   validator.New(),
		logger:      logger,
	}
}

// CreateCase handles the creation of a new case
func (h *CaseHandler) CreateCase(c *gin.Context) {
	span := tracer.StartSpan("case.create", tracer.ResourceName("CreateCase"))
	defer span.Finish()

	var req models.CreateCaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.ErrorWithSpan(span, "Failed to bind JSON request", logrus.Fields{"error": err.Error()})
		span.SetTag("error", true)
		span.SetTag("error.message", "Invalid request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Validate request
	if err := h.validator.Struct(req); err != nil {
		h.logger.ErrorWithSpan(span, "Validation failed", logrus.Fields{"error": err.Error()})
		span.SetTag("error", true)
		span.SetTag("error.message", "Validation failed")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
		return
	}

	// Create case
	caseObj := models.NewCase(&req)
	if err := h.caseService.CreateCase(caseObj); err != nil {
		h.logger.ErrorWithSpan(span, "Failed to create case", logrus.Fields{"error": err.Error()})
		span.SetTag("error", true)
		span.SetTag("error.message", "Failed to create case")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create case"})
		return
	}

	h.logger.InfoWithSpan(span, "Case created successfully", logrus.Fields{
		"caseId": caseObj.ID,
		"caseType": caseObj.CaseType,
		"filingAs": caseObj.FilingAs,
	})

	span.SetTag("case.id", caseObj.ID)
	span.SetTag("case.type", caseObj.CaseType)
	c.JSON(http.StatusCreated, caseObj)
}

// ListCases handles listing all cases with pagination
func (h *CaseHandler) ListCases(c *gin.Context) {
	span := tracer.StartSpan("case.list", tracer.ResourceName("ListCases"))
	defer span.Finish()

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	status := c.Query("status")

	span.SetTag("pagination.page", page)
	span.SetTag("pagination.limit", limit)
	span.SetTag("filter.status", status)

	cases, total, err := h.caseService.ListCases(page, limit, status)
	if err != nil {
		h.logger.ErrorWithSpan(span, "Failed to list cases", logrus.Fields{"error": err.Error()})
		span.SetTag("error", true)
		span.SetTag("error.message", "Failed to list cases")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list cases"})
		return
	}

	h.logger.InfoWithSpan(span, "Cases listed successfully", logrus.Fields{
		"total": total,
		"count": len(cases),
		"page": page,
		"limit": limit,
	})

	span.SetTag("cases.total", total)
	span.SetTag("cases.count", len(cases))

	c.JSON(http.StatusOK, gin.H{
		"cases": cases,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// GetCase handles retrieving a specific case by ID
func (h *CaseHandler) GetCase(c *gin.Context) {
	caseID := c.Param("id")
	if caseID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Case ID is required"})
		return
	}

	span := tracer.StartSpan("case.get", tracer.ResourceName("GetCase"))
	defer span.Finish()

	span.SetTag("case.id", caseID)

	caseObj, err := h.caseService.GetCase(caseID)
	if err != nil {
		h.logger.ErrorWithSpan(span, "Failed to get case", logrus.Fields{
			"caseId": caseID,
			"error": err.Error(),
		})
		span.SetTag("error", true)
		span.SetTag("error.message", "Case not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "Case not found"})
		return
	}

	h.logger.InfoWithSpan(span, "Case retrieved successfully", logrus.Fields{
		"caseId": caseID,
		"caseType": caseObj.CaseType,
	})

	c.JSON(http.StatusOK, caseObj)
}

// UpdateCase handles updating a case
func (h *CaseHandler) UpdateCase(c *gin.Context) {
	caseID := c.Param("id")
	if caseID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Case ID is required"})
		return
	}

	span := tracer.StartSpan("case.update", tracer.ResourceName("UpdateCase"))
	defer span.Finish()

	span.SetTag("case.id", caseID)

	var req models.CreateCaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.ErrorWithSpan(span, "Failed to bind JSON request", logrus.Fields{"error": err.Error()})
		span.SetTag("error", true)
		span.SetTag("error.message", "Invalid request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Validate request
	if err := h.validator.Struct(req); err != nil {
		h.logger.ErrorWithSpan(span, "Validation failed", logrus.Fields{"error": err.Error()})
		span.SetTag("error", true)
		span.SetTag("error.message", "Validation failed")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
		return
	}

	// Update case
	caseObj := models.NewCase(&req)
	caseObj.ID = caseID
	if err := h.caseService.UpdateCase(caseObj); err != nil {
		h.logger.ErrorWithSpan(span, "Failed to update case", logrus.Fields{
			"caseId": caseID,
			"error": err.Error(),
		})
		span.SetTag("error", true)
		span.SetTag("error.message", "Failed to update case")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update case"})
		return
	}

	h.logger.InfoWithSpan(span, "Case updated successfully", logrus.Fields{
		"caseId": caseID,
		"caseType": caseObj.CaseType,
	})

	c.JSON(http.StatusOK, caseObj)
}

// DeleteCase handles deleting a case
func (h *CaseHandler) DeleteCase(c *gin.Context) {
	caseID := c.Param("id")
	if caseID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Case ID is required"})
		return
	}

	span := tracer.StartSpan("case.delete", tracer.ResourceName("DeleteCase"))
	defer span.Finish()

	span.SetTag("case.id", caseID)

	if err := h.caseService.DeleteCase(caseID); err != nil {
		h.logger.ErrorWithSpan(span, "Failed to delete case", logrus.Fields{
			"caseId": caseID,
			"error": err.Error(),
		})
		span.SetTag("error", true)
		span.SetTag("error.message", "Failed to delete case")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete case"})
		return
	}

	h.logger.InfoWithSpan(span, "Case deleted successfully", logrus.Fields{
		"caseId": caseID,
	})

	c.JSON(http.StatusOK, gin.H{"message": "Case deleted successfully"})
}

// Global handler functions for compatibility with main.go
var (
	caseService *services.CaseService
	caseHandler *CaseHandler
)

func InitHandlers(logger *logger.DatadogLogger) {
	caseService = services.NewCaseService(logger)
	caseHandler = NewCaseHandler(caseService, logger)
}

func CreateCase(c *gin.Context) {
	if caseHandler == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Handler not initialized"})
		return
	}
	caseHandler.CreateCase(c)
}

func ListCases(c *gin.Context) {
	if caseHandler == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Handler not initialized"})
		return
	}
	caseHandler.ListCases(c)
}

func GetCase(c *gin.Context) {
	if caseHandler == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Handler not initialized"})
		return
	}
	caseHandler.GetCase(c)
}

func UpdateCase(c *gin.Context) {
	if caseHandler == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Handler not initialized"})
		return
	}
	caseHandler.UpdateCase(c)
}

func DeleteCase(c *gin.Context) {
	if caseHandler == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Handler not initialized"})
		return
	}
	caseHandler.DeleteCase(c)
}
