package model

import "time"

type WeightRecord struct {
	ID         string    `json:"id"`
	PetID      string    `json:"pet_id"`
	Weight     float64   `json:"weight"`
	RecordedAt string    `json:"recorded_at"`
	Note       *string   `json:"note"`
	CreatedAt  time.Time `json:"created_at"`
}

type CreateWeightRequest struct {
	Weight     float64  `json:"weight"`
	RecordedAt string   `json:"recorded_at"`
	Note       *string  `json:"note"`
}
