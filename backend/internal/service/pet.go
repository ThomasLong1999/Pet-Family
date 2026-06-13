package service

import (
	"cat-manager/backend/internal/model"
	"cat-manager/backend/internal/repository"
)

type PetService struct {
	repo      *repository.PetRepo
	weightSvc *WeightService
}

func NewPetService(repo *repository.PetRepo, weightSvc *WeightService) *PetService {
	return &PetService{repo: repo, weightSvc: weightSvc}
}

func (s *PetService) List() ([]model.Pet, error) {
	return s.repo.List()
}

func (s *PetService) GetByID(id string) (*model.Pet, error) {
	return s.repo.GetByID(id)
}

func (s *PetService) Create(req model.CreatePetRequest) (*model.Pet, error) {
	return s.repo.Create(req)
}

func (s *PetService) Update(id string, req model.UpdatePetRequest) (*model.Pet, error) {
	return s.repo.Update(id, req)
}

func (s *PetService) Delete(id string) error {
	return s.repo.Delete(id)
}

func (s *PetService) UpdateAvatar(id string, avatarURL string, dominantColor string) error {
	return s.repo.UpdateAvatar(id, avatarURL, dominantColor)
}
