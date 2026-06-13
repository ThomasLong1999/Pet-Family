package repository

import (
	"cat-manager/backend/internal/model"
	"testing"
)

func seedPet(t *testing.T, repo *PetRepo) *model.Pet {
	t.Helper()
	pet, err := repo.Create(model.CreatePetRequest{
		Species: "cat", Name: "Test", Gender: "male", Birthday: "2020-01-01",
	})
	if err != nil {
		t.Fatalf("seed pet: %v", err)
	}
	return pet
}

func TestWeightRepo_CreateListDelete(t *testing.T) {
	db := newTestDB(t)
	petRepo := NewPetRepo(db)
	repo := NewWeightRepo(db)
	pet := seedPet(t, petRepo)

	rec, err := repo.Create(pet.ID, model.CreateWeightRequest{
		Weight: 4.2, RecordedAt: "2024-05-01",
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if rec.Weight != 4.2 {
		t.Fatalf("expected weight 4.2, got %f", rec.Weight)
	}

	records, err := repo.List(pet.ID)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(records))
	}

	if err := repo.Delete(pet.ID, rec.ID); err != nil {
		t.Fatalf("delete: %v", err)
	}
	records, _ = repo.List(pet.ID)
	if len(records) != 0 {
		t.Fatalf("expected 0 after delete, got %d", len(records))
	}
}

func TestWeightRepo_ListOrdersByDateDesc(t *testing.T) {
	db := newTestDB(t)
	petRepo := NewPetRepo(db)
	repo := NewWeightRepo(db)
	pet := seedPet(t, petRepo)

	for _, date := range []string{"2024-01-01", "2024-06-01", "2024-03-01"} {
		if _, err := repo.Create(pet.ID, model.CreateWeightRequest{
			Weight: 4.0, RecordedAt: date,
		}); err != nil {
			t.Fatalf("create %s: %v", date, err)
		}
	}

	records, _ := repo.List(pet.ID)
	if len(records) != 3 {
		t.Fatalf("expected 3, got %d", len(records))
	}
	// Most recent first
	if records[0].RecordedAt != "2024-06-01" {
		t.Fatalf("expected newest first, got %q", records[0].RecordedAt)
	}
	if records[2].RecordedAt != "2024-01-01" {
		t.Fatalf("expected oldest last, got %q", records[2].RecordedAt)
	}
}

func TestWeightRepo_GetLatestAndPrevious(t *testing.T) {
	db := newTestDB(t)
	petRepo := NewPetRepo(db)
	repo := NewWeightRepo(db)
	pet := seedPet(t, petRepo)

	if _, err := repo.Create(pet.ID, model.CreateWeightRequest{Weight: 3.0, RecordedAt: "2024-01-01"}); err != nil {
		t.Fatal(err)
	}
	if _, err := repo.Create(pet.ID, model.CreateWeightRequest{Weight: 5.0, RecordedAt: "2024-06-01"}); err != nil {
		t.Fatal(err)
	}

	latest, err := repo.GetLatest(pet.ID)
	if err != nil {
		t.Fatalf("get latest: %v", err)
	}
	if latest.Weight != 5.0 {
		t.Fatalf("expected latest 5.0, got %f", latest.Weight)
	}

	prev, err := repo.GetPrevious(pet.ID, "2024-06-01")
	if err != nil {
		t.Fatalf("get previous: %v", err)
	}
	if prev.Weight != 3.0 {
		t.Fatalf("expected previous 3.0, got %f", prev.Weight)
	}
}
