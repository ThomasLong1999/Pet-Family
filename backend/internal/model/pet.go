package model

import "time"

type Pet struct {
	ID            string     `json:"id"`
	Species       string     `json:"species"` // cat, dog, hamster, rabbit
	Name          string     `json:"name"`
	Breed         string     `json:"breed"`
	Gender        string     `json:"gender"` // male, female
	Birthday      string     `json:"birthday"`
	Color         string     `json:"color"`
	AvatarURL     string     `json:"avatar_url"`
	DominantColor string     `json:"dominant_color"`
	AdoptedAt     *string    `json:"adopted_at"`
	PassedAt      *string    `json:"passed_at"`
	Note          *string    `json:"note"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type CreatePetRequest struct {
	Species   string  `json:"species"`
	Name      string  `json:"name"`
	Breed     string  `json:"breed"`
	Gender    string  `json:"gender"`
	Birthday  string  `json:"birthday"`
	Color     string  `json:"color"`
	AdoptedAt *string `json:"adopted_at"`
	Note      *string `json:"note"`
}

type UpdatePetRequest struct {
	Species   *string `json:"species"`
	Name      *string `json:"name"`
	Breed     *string `json:"breed"`
	Gender    *string `json:"gender"`
	Birthday  *string `json:"birthday"`
	Color     *string `json:"color"`
	AdoptedAt *string `json:"adopted_at"`
	PassedAt  *string `json:"passed_at"`
	Note      *string `json:"note"`
}
