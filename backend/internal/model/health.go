package model

import "time"

type HealthRecord struct {
	ID        string    `json:"id"`
	PetID     string    `json:"pet_id"`
	Type      string    `json:"type"` // vaccine, deworming, checkup
	Name      string    `json:"name"`
	Date      string    `json:"date"`
	NextDate  *string   `json:"next_date"`
	Note      *string   `json:"note"`
	ReportURL string    `json:"report_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateHealthRequest struct {
	Type     string  `json:"type"`
	Name     string  `json:"name"`
	Date     string  `json:"date"`
	NextDate *string `json:"next_date"`
	Note     *string `json:"note"`
}

type UpdateHealthRequest struct {
	Type     *string `json:"type"`
	Name     *string `json:"name"`
	Date     *string `json:"date"`
	NextDate *string `json:"next_date"`
	Note     *string `json:"note"`
}

type HealthReminder struct {
	PetID      string `json:"pet_id"`
	PetName    string `json:"pet_name"`
	RecordType string `json:"record_type"`
	Name       string `json:"name"`
	NextDate   string `json:"next_date"`
	DaysLeft   int    `json:"days_left"`
}
