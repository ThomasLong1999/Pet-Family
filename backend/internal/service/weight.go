package service

import (
	"cat-manager/backend/internal/model"
	"cat-manager/backend/internal/repository"
)

type WeightService struct {
	repo *repository.WeightRepo
}

func NewWeightService(repo *repository.WeightRepo) *WeightService {
	return &WeightService{repo: repo}
}

func (s *WeightService) List(petID string) ([]model.WeightRecord, error) {
	return s.repo.List(petID)
}

func (s *WeightService) Create(petID string, req model.CreateWeightRequest) (*model.WeightRecord, error) {
	return s.repo.Create(petID, req)
}

func (s *WeightService) Delete(petID string, recordID string) error {
	return s.repo.Delete(petID, recordID)
}
