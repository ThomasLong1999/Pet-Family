package repository

import (
	"cat-manager/backend/internal/model"

	"github.com/google/uuid"
)

type HealthRepo struct {
	db *DB
}

func NewHealthRepo(db *DB) *HealthRepo {
	return &HealthRepo{db: db}
}

const healthCols = "id, pet_id, type, name, date, next_date, note, report_url, created_at, updated_at"

func scanHealth(scanner interface{ Scan(...interface{}) error }, h *model.HealthRecord) error {
	return scanner.Scan(&h.ID, &h.PetID, &h.Type, &h.Name, &h.Date, &h.NextDate, &h.Note, &h.ReportURL, &h.CreatedAt, &h.UpdatedAt)
}

func (r *HealthRepo) List(petID string, recordType string) ([]model.HealthRecord, error) {
	query := "SELECT " + healthCols + " FROM health_records WHERE pet_id = ?"
	args := []interface{}{petID}

	if recordType != "" {
		query += " AND type = ?"
		args = append(args, recordType)
	}
	query += " ORDER BY date DESC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []model.HealthRecord
	for rows.Next() {
		var h model.HealthRecord
		if err := scanHealth(rows, &h); err != nil {
			return nil, err
		}
		records = append(records, h)
	}
	return records, rows.Err()
}

func (r *HealthRepo) GetByID(recordID string) (*model.HealthRecord, error) {
	var h model.HealthRecord
	err := scanHealth(r.db.QueryRow("SELECT "+healthCols+" FROM health_records WHERE id = ?", recordID), &h)
	if err != nil {
		return nil, err
	}
	return &h, nil
}

func (r *HealthRepo) Create(petID string, req model.CreateHealthRequest, reportURL string) (*model.HealthRecord, error) {
	id := uuid.New().String()

	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`INSERT INTO health_records (id, pet_id, type, name, date, next_date, note, report_url, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, datetime('now'), datetime('now'))`,
		id, petID, req.Type, req.Name, req.Date, req.NextDate, req.Note, reportURL)
	if err != nil {
		return nil, err
	}

	var h model.HealthRecord
	err = scanHealth(tx.QueryRow("SELECT "+healthCols+" FROM health_records WHERE id = ?", id), &h)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return &h, nil
}

func (r *HealthRepo) Update(recordID string, req model.UpdateHealthRequest) (*model.HealthRecord, error) {
	record, err := r.GetByID(recordID)
	if err != nil {
		return nil, err
	}

	if req.Type != nil {
		record.Type = *req.Type
	}
	if req.Name != nil {
		record.Name = *req.Name
	}
	if req.Date != nil {
		record.Date = *req.Date
	}
	if req.NextDate != nil {
		record.NextDate = req.NextDate
	}
	if req.Note != nil {
		record.Note = req.Note
	}

	_, err = r.db.Exec(`UPDATE health_records SET type=?, name=?, date=?, next_date=?, note=?, updated_at=datetime('now') WHERE id=?`,
		record.Type, record.Name, record.Date, record.NextDate, record.Note, recordID)
	if err != nil {
		return nil, err
	}

	return r.GetByID(recordID)
}

func (r *HealthRepo) Delete(petID string, recordID string) error {
	_, err := r.db.Exec("DELETE FROM health_records WHERE id = ? AND pet_id = ?", recordID, petID)
	return err
}

func (r *HealthRepo) GetUpcomingReminders(days int) ([]model.HealthReminder, error) {
	query := `SELECT h.pet_id, p.name, h.type, h.name, h.next_date,
		CAST(julianday(h.next_date) - julianday('now') AS INTEGER) as days_left
		FROM health_records h
		JOIN pets p ON p.id = h.pet_id
		WHERE h.next_date IS NOT NULL
		AND julianday(h.next_date) - julianday('now') <= ?
		AND julianday(h.next_date) - julianday('now') >= 0
		ORDER BY h.next_date ASC`

	rows, err := r.db.Query(query, days)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reminders []model.HealthReminder
	for rows.Next() {
		var rem model.HealthReminder
		if err := rows.Scan(&rem.PetID, &rem.PetName, &rem.RecordType, &rem.Name, &rem.NextDate, &rem.DaysLeft); err != nil {
			return nil, err
		}
		reminders = append(reminders, rem)
	}
	return reminders, rows.Err()
}
