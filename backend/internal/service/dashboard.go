package service

import (
	"cat-manager/backend/internal/model"
	"cat-manager/backend/internal/repository"
)

type DashboardService struct {
	petRepo    *repository.PetRepo
	weightRepo *repository.WeightRepo
	healthRepo *repository.HealthRepo
}

func NewDashboardService(petRepo *repository.PetRepo, weightRepo *repository.WeightRepo, healthRepo *repository.HealthRepo) *DashboardService {
	return &DashboardService{petRepo: petRepo, weightRepo: weightRepo, healthRepo: healthRepo}
}

func (s *DashboardService) Get() (*model.Dashboard, error) {
	pets, err := s.petRepo.List()
	if err != nil {
		return nil, err
	}

	reminders, err := s.healthRepo.GetUpcomingReminders(7)
	if err != nil {
		return nil, err
	}

	summaries := make([]model.PetSummary, 0, len(pets))
	for _, p := range pets {
		summary := model.PetSummary{
			ID:            p.ID,
			Name:          p.Name,
			Species:       p.Species,
			AvatarURL:     p.AvatarURL,
			DominantColor: p.DominantColor,
			PassedAt:      p.PassedAt,
		}

		latest, err := s.weightRepo.GetLatest(p.ID)
		if err == nil {
			summary.LatestWeight = &latest.Weight

			prev, err := s.weightRepo.GetPrevious(p.ID, latest.RecordedAt)
			if err == nil {
				trend := latest.Weight - prev.Weight
				summary.WeightTrend = &trend
			}
		}

		summaries = append(summaries, summary)
	}

	return &model.Dashboard{
		Pets:      summaries,
		Reminders: reminders,
	}, nil
}
