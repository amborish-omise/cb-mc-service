package models

import (
	"time"

	"github.com/google/uuid"
)

// Document represents a document attached to a case
type Document struct {
	ID          string    `json:"id" bson:"_id"`
	CaseID      string    `json:"caseId" validate:"required"`
	FileName    string    `json:"fileName" validate:"required"`
	FileType    string    `json:"fileType" validate:"required"`
	FileSize    int64     `json:"fileSize"`
	Content     []byte    `json:"content,omitempty"`
	UploadedBy  string    `json:"uploadedBy"`
	UploadedAt  time.Time `json:"uploadedAt"`
	Description string    `json:"description"`
}

// DocumentResponse represents the response for document operations
type DocumentResponse struct {
	ID          string    `json:"id"`
	CaseID      string    `json:"caseId"`
	FileName    string    `json:"fileName"`
	FileType    string    `json:"fileType"`
	FileSize    int64     `json:"fileSize"`
	UploadedBy  string    `json:"uploadedBy"`
	UploadedAt  time.Time `json:"uploadedAt"`
	Description string    `json:"description"`
}

// NewDocument creates a new document
func NewDocument(caseID, fileName, fileType string, content []byte, uploadedBy, description string) *Document {
	return &Document{
		ID:          uuid.New().String(),
		CaseID:      caseID,
		FileName:    fileName,
		FileType:    fileType,
		FileSize:    int64(len(content)),
		Content:     content,
		UploadedBy:  uploadedBy,
		UploadedAt:  time.Now(),
		Description: description,
	}
}
