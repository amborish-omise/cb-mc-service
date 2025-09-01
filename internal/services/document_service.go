package services

import (
	"errors"
	"sync"

	"mastercom-service/internal/models"
	"mastercom-service/pkg/logger"

	"github.com/sirupsen/logrus"
)

type DocumentService struct {
	documents map[string]*models.Document
	mutex     sync.RWMutex
	logger    *logger.DatadogLogger
}

func NewDocumentService(logger *logger.DatadogLogger) *DocumentService {
	return &DocumentService{
		documents: make(map[string]*models.Document),
		logger:    logger,
	}
}

func (s *DocumentService) UploadDocument(document *models.Document) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Check if document already exists
	if _, exists := s.documents[document.ID]; exists {
		return errors.New("document already exists")
	}

	// Store the document
	s.documents[document.ID] = document
	s.logger.Info("Document uploaded successfully", logrus.Fields{"documentId": document.ID})
	return nil
}

func (s *DocumentService) GetDocument(documentID string) (*models.Document, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	document, exists := s.documents[documentID]
	if !exists {
		return nil, errors.New("document not found")
	}

	return document, nil
}

func (s *DocumentService) DeleteDocument(documentID string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Check if document exists
	if _, exists := s.documents[documentID]; !exists {
		return errors.New("document not found")
	}

	// Delete the document
	delete(s.documents, documentID)
	s.logger.Info("Document deleted successfully", logrus.Fields{"documentId": documentID})
	return nil
}

func (s *DocumentService) GetDocumentsByCaseID(caseID string) ([]*models.Document, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var documents []*models.Document
	for _, document := range s.documents {
		if document.CaseID == caseID {
			documents = append(documents, document)
		}
	}

	return documents, nil
}
