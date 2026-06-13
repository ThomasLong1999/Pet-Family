package service

import (
	"cat-manager/backend/internal/model"
	"cat-manager/backend/internal/repository"
)

type HealthService struct {
	repo *repository.HealthRepo
}

func NewHealthService(repo *repository.HealthRepo) *HealthService {
	return &HealthService{repo: repo}
}

func (s *HealthService) List(petID string, recordType string) ([]model.HealthRecord, error) {
	return s.repo.List(petID, recordType)
}

func (s *HealthService) Create(petID string, req model.CreateHealthRequest, reportURL string) (*model.HealthRecord, error) {
	return s.repo.Create(petID, req, reportURL)
}

func (s *HealthService) Update(recordID string, req model.UpdateHealthRequest) (*model.HealthRecord, error) {
	return s.repo.Update(recordID, req)
}

func (s *HealthService) Delete(petID string, recordID string) error {
	return s.repo.Delete(petID, recordID)
}
