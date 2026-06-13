package model

import "time"

type Photo struct {
	ID        string    `json:"id"`
	PetID     string    `json:"pet_id"`
	URL       string    `json:"url"`
	AgeGroup  string    `json:"age_group"`
	Caption   *string   `json:"caption"`
	CreatedAt time.Time `json:"created_at"`
}
