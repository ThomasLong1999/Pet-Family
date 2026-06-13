package model

type Dashboard struct {
	Pets     []PetSummary     `json:"pets"`
	Reminders []HealthReminder `json:"reminders"`
}

type PetSummary struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Species       string   `json:"species"`
	AvatarURL     string   `json:"avatar_url"`
	DominantColor string   `json:"dominant_color"`
	LatestWeight  *float64 `json:"latest_weight"`
	WeightTrend   *float64 `json:"weight_trend"` // difference from previous
	PassedAt      *string  `json:"passed_at"`
}
