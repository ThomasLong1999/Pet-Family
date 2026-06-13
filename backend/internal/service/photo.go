package service

import (
	"cat-manager/backend/internal/model"
	"cat-manager/backend/internal/repository"
	"fmt"
	"math"
	"time"
)

type PhotoService struct {
	repo *repository.PhotoRepo
}

func NewPhotoService(repo *repository.PhotoRepo) *PhotoService {
	return &PhotoService{repo: repo}
}

func (s *PhotoService) List(petID string) ([]model.Photo, error) {
	return s.repo.List(petID)
}

// CreateWithAgeGroup creates a photo with an explicitly provided age group.
func (s *PhotoService) CreateWithAgeGroup(petID string, url string, ageGroup string, caption *string) (*model.Photo, error) {
	return s.repo.Create(petID, url, ageGroup, caption)
}

func (s *PhotoService) Delete(petID string, photoID string) error {
	return s.repo.Delete(petID, photoID)
}

// CalculateAgeGroup determines the age group label based on birthday.
// Exported so the handler can use it as fallback.
func CalculateAgeGroup(birthday string) string {
	bd, err := time.Parse("2006-01-02", birthday)
	if err != nil {
		return "1y"
	}

	now := time.Now()
	months := int(math.Floor(now.Sub(bd).Hours() / 24 / 30.44))

	if months <= 0 {
		return "1m"
	}
	if months <= 12 {
		return fmt.Sprintf("%dm", months)
	}
	years := months / 12
	if years > 20 {
		years = 20
	}
	return fmt.Sprintf("%dy", years)
}

// keep old unexported version for backward compat in tests if any
func calculateAgeGroup(birthday string) string {
	return CalculateAgeGroup(birthday)
}
