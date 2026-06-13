package repository

import (
	"cat-manager/backend/internal/model"
	"testing"
)

func TestPetRepo_CreateAndGet(t *testing.T) {
	db := newTestDB(t)
	repo := NewPetRepo(db)

	req := model.CreatePetRequest{
		Species:  "cat",
		Name:     "Mimi",
		Breed:    "British Shorthair",
		Gender:   "female",
		Birthday: "2022-03-15",
		Color:    "blue",
	}
	created, err := repo.Create(req)
	if err != nil {
		t.Fatalf("create pet: %v", err)
	}
	if created.ID == "" {
		t.Fatal("expected non-empty id")
	}
	if created.Name != "Mimi" || created.Species != "cat" {
		t.Fatalf("unexpected pet: %+v", created)
	}
	if created.AvatarURL != "" {
		t.Fatalf("expected empty avatar_url, got %q", created.AvatarURL)
	}

	got, err := repo.GetByID(created.ID)
	if err != nil {
		t.Fatalf("get pet: %v", err)
	}
	if got.Name != "Mimi" {
		t.Fatalf("got name %q", got.Name)
	}
}

func TestPetRepo_List(t *testing.T) {
	db := newTestDB(t)
	repo := NewPetRepo(db)

	for _, name := range []string{"A", "B", "C"} {
		if _, err := repo.Create(model.CreatePetRequest{
			Species: "dog", Name: name, Gender: "male", Birthday: "2023-01-01",
		}); err != nil {
			t.Fatalf("create %s: %v", name, err)
		}
	}

	pets, err := repo.List()
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(pets) != 3 {
		t.Fatalf("expected 3 pets, got %d", len(pets))
	}
}

func TestPetRepo_Update(t *testing.T) {
	db := newTestDB(t)
	repo := NewPetRepo(db)

	pet, err := repo.Create(model.CreatePetRequest{
		Species: "cat", Name: "Old", Gender: "male", Birthday: "2020-01-01",
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	newName := "New"
	passed := "2024-12-01"
	updated, err := repo.Update(pet.ID, model.UpdatePetRequest{
		Name:     &newName,
		PassedAt: &passed,
	})
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if updated.Name != "New" {
		t.Fatalf("expected name New, got %q", updated.Name)
	}
	if updated.PassedAt == nil || *updated.PassedAt != "2024-12-01" {
		t.Fatalf("expected passed_at set, got %+v", updated.PassedAt)
	}
}

func TestPetRepo_UpdateClearPassedAt(t *testing.T) {
	db := newTestDB(t)
	repo := NewPetRepo(db)

	passed := "2024-12-01"
	pet, err := repo.Create(model.CreatePetRequest{
		Species: "cat", Name: "Ghost", Gender: "male", Birthday: "2010-01-01",
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if _, err := repo.Update(pet.ID, model.UpdatePetRequest{PassedAt: &passed}); err != nil {
		t.Fatalf("set passed: %v", err)
	}

	// Empty string should clear passed_at to NULL
	empty := ""
	if _, err := repo.Update(pet.ID, model.UpdatePetRequest{PassedAt: &empty}); err != nil {
		t.Fatalf("clear passed: %v", err)
	}
	got, err := repo.GetByID(pet.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.PassedAt != nil {
		t.Fatalf("expected nil passed_at after clear, got %v", got.PassedAt)
	}
}

func TestPetRepo_DeleteCascadesWeights(t *testing.T) {
	db := newTestDB(t)
	petRepo := NewPetRepo(db)
	weightRepo := NewWeightRepo(db)

	pet, err := petRepo.Create(model.CreatePetRequest{
		Species: "cat", Name: "Casc", Gender: "male", Birthday: "2020-01-01",
	})
	if err != nil {
		t.Fatalf("create pet: %v", err)
	}
	if _, err := weightRepo.Create(pet.ID, model.CreateWeightRequest{
		Weight: 4.5, RecordedAt: "2024-01-01",
	}); err != nil {
		t.Fatalf("create weight: %v", err)
	}

	weights, _ := weightRepo.List(pet.ID)
	if len(weights) != 1 {
		t.Fatalf("expected 1 weight before delete, got %d", len(weights))
	}

	if err := petRepo.Delete(pet.ID); err != nil {
		t.Fatalf("delete pet: %v", err)
	}

	// Foreign key ON DELETE CASCADE should remove the weight record
	weights, _ = weightRepo.List(pet.ID)
	if len(weights) != 0 {
		t.Fatalf("expected 0 weights after cascade delete, got %d", len(weights))
	}
}
