package repository

import (
	"cat-manager/backend/internal/model"

	"github.com/google/uuid"
)

type WeightRepo struct {
	db *DB
}

func NewWeightRepo(db *DB) *WeightRepo {
	return &WeightRepo{db: db}
}

func (r *WeightRepo) List(petID string) ([]model.WeightRecord, error) {
	rows, err := r.db.Query("SELECT id, pet_id, weight, recorded_at, note, created_at FROM weight_records WHERE pet_id = ? ORDER BY recorded_at DESC", petID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []model.WeightRecord
	for rows.Next() {
		var w model.WeightRecord
		if err := rows.Scan(&w.ID, &w.PetID, &w.Weight, &w.RecordedAt, &w.Note, &w.CreatedAt); err != nil {
			return nil, err
		}
		records = append(records, w)
	}
	return records, rows.Err()
}

func (r *WeightRepo) Create(petID string, req model.CreateWeightRequest) (*model.WeightRecord, error) {
	id := uuid.New().String()

	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`INSERT INTO weight_records (id, pet_id, weight, recorded_at, note, created_at) VALUES (?, ?, ?, ?, ?, datetime('now'))`,
		id, petID, req.Weight, req.RecordedAt, req.Note)
	if err != nil {
		return nil, err
	}

	var w model.WeightRecord
	err = tx.QueryRow("SELECT id, pet_id, weight, recorded_at, note, created_at FROM weight_records WHERE id = ?", id).
		Scan(&w.ID, &w.PetID, &w.Weight, &w.RecordedAt, &w.Note, &w.CreatedAt)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return &w, nil
}

func (r *WeightRepo) Delete(petID string, recordID string) error {
	_, err := r.db.Exec("DELETE FROM weight_records WHERE id = ? AND pet_id = ?", recordID, petID)
	return err
}

func (r *WeightRepo) GetLatest(petID string) (*model.WeightRecord, error) {
	var w model.WeightRecord
	err := r.db.QueryRow("SELECT id, pet_id, weight, recorded_at, note, created_at FROM weight_records WHERE pet_id = ? ORDER BY recorded_at DESC LIMIT 1", petID).
		Scan(&w.ID, &w.PetID, &w.Weight, &w.RecordedAt, &w.Note, &w.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &w, nil
}

func (r *WeightRepo) GetPrevious(petID string, beforeDate string) (*model.WeightRecord, error) {
	var w model.WeightRecord
	err := r.db.QueryRow("SELECT id, pet_id, weight, recorded_at, note, created_at FROM weight_records WHERE pet_id = ? AND recorded_at < ? ORDER BY recorded_at DESC LIMIT 1", petID, beforeDate).
		Scan(&w.ID, &w.PetID, &w.Weight, &w.RecordedAt, &w.Note, &w.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &w, nil
}
