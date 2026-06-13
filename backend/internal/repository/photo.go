package repository

import (
	"cat-manager/backend/internal/model"

	"github.com/google/uuid"
)

type PhotoRepo struct {
	db *DB
}

func NewPhotoRepo(db *DB) *PhotoRepo {
	return &PhotoRepo{db: db}
}

func (r *PhotoRepo) List(petID string) ([]model.Photo, error) {
	rows, err := r.db.Query("SELECT id, pet_id, url, age_group, caption, created_at FROM photos WHERE pet_id = ? ORDER BY age_group, created_at", petID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var photos []model.Photo
	for rows.Next() {
		var p model.Photo
		if err := rows.Scan(&p.ID, &p.PetID, &p.URL, &p.AgeGroup, &p.Caption, &p.CreatedAt); err != nil {
			return nil, err
		}
		photos = append(photos, p)
	}
	return photos, rows.Err()
}

func (r *PhotoRepo) Create(petID string, url string, ageGroup string, caption *string) (*model.Photo, error) {
	id := uuid.New().String()

	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`INSERT INTO photos (id, pet_id, url, age_group, caption, created_at) VALUES (?, ?, ?, ?, ?, datetime('now'))`,
		id, petID, url, ageGroup, caption)
	if err != nil {
		return nil, err
	}

	var p model.Photo
	err = tx.QueryRow("SELECT id, pet_id, url, age_group, caption, created_at FROM photos WHERE id = ?", id).
		Scan(&p.ID, &p.PetID, &p.URL, &p.AgeGroup, &p.Caption, &p.CreatedAt)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PhotoRepo) Delete(petID string, photoID string) error {
	_, err := r.db.Exec("DELETE FROM photos WHERE id = ? AND pet_id = ?", photoID, petID)
	return err
}

func (r *PhotoRepo) GetByID(photoID string) (*model.Photo, error) {
	var p model.Photo
	err := r.db.QueryRow("SELECT id, pet_id, url, age_group, caption, created_at FROM photos WHERE id = ?", photoID).
		Scan(&p.ID, &p.PetID, &p.URL, &p.AgeGroup, &p.Caption, &p.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}
