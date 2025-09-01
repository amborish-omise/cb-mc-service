package services

import (
	"errors"
	"sync"
	"time"

	"mastercom-service/internal/models"
	"mastercom-service/pkg/logger"

	"github.com/sirupsen/logrus"
)

type CaseService struct {
	cases  map[string]*models.Case
	mutex  sync.RWMutex
	logger *logger.DatadogLogger
}

func NewCaseService(logger *logger.DatadogLogger) *CaseService {
	return &CaseService{
		cases:  make(map[string]*models.Case),
		logger: logger,
	}
}

func (s *CaseService) CreateCase(caseObj *models.Case) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Check if case already exists
	if _, exists := s.cases[caseObj.ID]; exists {
		return errors.New("case already exists")
	}

	// Store the case
	s.cases[caseObj.ID] = caseObj
	s.logger.Info("Case created successfully", logrus.Fields{"caseId": caseObj.ID})
	return nil
}

func (s *CaseService) GetCase(caseID string) (*models.Case, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	caseObj, exists := s.cases[caseID]
	if !exists {
		return nil, errors.New("case not found")
	}

	return caseObj, nil
}

func (s *CaseService) ListCases(page, limit int, status string) ([]*models.Case, int, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var filteredCases []*models.Case
	for _, caseObj := range s.cases {
		if status == "" || caseObj.Status == status {
			filteredCases = append(filteredCases, caseObj)
		}
	}

	total := len(filteredCases)

	// Simple pagination
	start := (page - 1) * limit
	end := start + limit
	if start >= total {
		return []*models.Case{}, total, nil
	}
	if end > total {
		end = total
	}

	return filteredCases[start:end], total, nil
}

func (s *CaseService) UpdateCase(caseObj *models.Case) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Check if case exists
	if _, exists := s.cases[caseObj.ID]; !exists {
		return errors.New("case not found")
	}

	// Update timestamp
	caseObj.UpdatedAt = time.Now()

	// Store the updated case
	s.cases[caseObj.ID] = caseObj
	s.logger.Info("Case updated successfully", logrus.Fields{"caseId": caseObj.ID})
	return nil
}

func (s *CaseService) DeleteCase(caseID string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Check if case exists
	if _, exists := s.cases[caseID]; !exists {
		return errors.New("case not found")
	}

	// Delete the case
	delete(s.cases, caseID)
	s.logger.Info("Case deleted successfully", logrus.Fields{"caseId": caseID})
	return nil
}
