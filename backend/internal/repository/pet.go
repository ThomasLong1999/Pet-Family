package repository

import (
	"cat-manager/backend/internal/model"

	"github.com/google/uuid"
)

type PetRepo struct {
	db *DB
}

func NewPetRepo(db *DB) *PetRepo {
	return &PetRepo{db: db}
}

func (r *PetRepo) List() ([]model.Pet, error) {
	rows, err := r.db.Query("SELECT id, species, name, breed, gender, birthday, color, avatar_url, dominant_color, adopted_at, passed_at, note, created_at, updated_at FROM pets ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pets []model.Pet
	for rows.Next() {
		var p model.Pet
		if err := rows.Scan(&p.ID, &p.Species, &p.Name, &p.Breed, &p.Gender, &p.Birthday, &p.Color, &p.AvatarURL, &p.DominantColor, &p.AdoptedAt, &p.PassedAt, &p.Note, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		pets = append(pets, p)
	}
	return pets, rows.Err()
}

func (r *PetRepo) GetByID(id string) (*model.Pet, error) {
	var p model.Pet
	err := r.db.QueryRow("SELECT id, species, name, breed, gender, birthday, color, avatar_url, dominant_color, adopted_at, passed_at, note, created_at, updated_at FROM pets WHERE id = ?", id).
		Scan(&p.ID, &p.Species, &p.Name, &p.Breed, &p.Gender, &p.Birthday, &p.Color, &p.AvatarURL, &p.DominantColor, &p.AdoptedAt, &p.PassedAt, &p.Note, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PetRepo) Create(req model.CreatePetRequest) (*model.Pet, error) {
	id := uuid.New().String()

	// Use a transaction to handle the default timestamp
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`INSERT INTO pets (id, species, name, breed, gender, birthday, color, adopted_at, note, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, datetime('now'), datetime('now'))`,
		id, req.Species, req.Name, req.Breed, req.Gender, req.Birthday, req.Color, req.AdoptedAt, req.Note)
	if err != nil {
		return nil, err
	}

	var p model.Pet
	err = tx.QueryRow("SELECT id, species, name, breed, gender, birthday, color, avatar_url, dominant_color, adopted_at, passed_at, note, created_at, updated_at FROM pets WHERE id = ?", id).
		Scan(&p.ID, &p.Species, &p.Name, &p.Breed, &p.Gender, &p.Birthday, &p.Color, &p.AvatarURL, &p.DominantColor, &p.AdoptedAt, &p.PassedAt, &p.Note, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &p, nil
}

func (r *PetRepo) Update(id string, req model.UpdatePetRequest) (*model.Pet, error) {
	pet, err := r.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.Species != nil {
		pet.Species = *req.Species
	}
	if req.Name != nil {
		pet.Name = *req.Name
	}
	if req.Breed != nil {
		pet.Breed = *req.Breed
	}
	if req.Gender != nil {
		pet.Gender = *req.Gender
	}
	if req.Birthday != nil {
		pet.Birthday = *req.Birthday
	}
	if req.Color != nil {
		pet.Color = *req.Color
	}
	if req.AdoptedAt != nil {
		pet.AdoptedAt = req.AdoptedAt
	}
	if req.PassedAt != nil {
		if *req.PassedAt == "" {
			pet.PassedAt = nil
		} else {
			pet.PassedAt = req.PassedAt
		}
	}
	if req.Note != nil {
		pet.Note = req.Note
	}

	_, err = r.db.Exec(`UPDATE pets SET species=?, name=?, breed=?, gender=?, birthday=?, color=?, adopted_at=?, passed_at=?, note=?, updated_at=datetime('now') WHERE id=?`,
		pet.Species, pet.Name, pet.Breed, pet.Gender, pet.Birthday, pet.Color, pet.AdoptedAt, pet.PassedAt, pet.Note, id)
	if err != nil {
		return nil, err
	}

	return r.GetByID(id)
}

func (r *PetRepo) Delete(id string) error {
	_, err := r.db.Exec("DELETE FROM pets WHERE id = ?", id)
	return err
}

func (r *PetRepo) UpdateAvatar(id string, avatarURL string, dominantColor string) error {
	_, err := r.db.Exec("UPDATE pets SET avatar_url = ?, dominant_color = ?, updated_at = datetime('now') WHERE id = ?", avatarURL, dominantColor, id)
	return err
}
