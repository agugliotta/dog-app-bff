package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"
	"time"

	"github.com/agugliotta/dog-app-bff/internal/store"
	"github.com/agugliotta/dog-app-bff/internal/types"
)

type PetStoreMock struct {
	pets   []types.Pet
	breeds []types.Breed
	fail   bool
}

func (m *PetStoreMock) GetPets() ([]types.Pet, error) {
	if m.fail {
		return nil, errors.New("store error")
	}
	return m.pets, nil
}

func (m *PetStoreMock) GetPetByID(id string) (*types.Pet, error) {
	if m.fail {
		return nil, errors.New("store error")
	}
	for _, p := range m.pets {
		if p.ID == id {
			return &p, nil
		}
	}
	return nil, store.ErrNotFound
}

func (m *PetStoreMock) CreatePet(name string, birth time.Time, breedID string) (*types.Pet, error) {
	if m.fail {
		return nil, errors.New("store error")
	}
	var breed types.Breed
	found := false
	for _, b := range m.breeds {
		if b.ID == breedID {
			breed = b
			found = true
			break
		}
	}
	if !found {
		return nil, store.ErrNotFound
	}
	pet := types.Pet{
		ID:    "new-pet-id",
		Name:  name,
		Birth: birth,
		Breed: breed,
	}
	m.pets = append(m.pets, pet)
	return &pet, nil
}

func (m *PetStoreMock) DeletePet(id string) error {
	index := -1
	for i, p := range m.pets {
		if p.ID == id {
			index = i
		}
	}
	if index == -1 {
		return store.ErrNotFound
	}
	m.pets = slices.Delete(m.pets, index, index+1)
	return nil
}

type BreedStoreMock struct {
	breeds []types.Breed
}

func (m *BreedStoreMock) GetBreedByID(id string) (*types.Breed, error) {
	for _, b := range m.breeds {
		if b.ID == id {
			return &b, nil
		}
	}
	return nil, store.ErrNotFound
}

func (m *BreedStoreMock) GetBreeds() ([]types.Breed, error) {
	return m.breeds, nil
}

func TestGetPetsHandler(t *testing.T) {
	breeds := []types.Breed{{ID: "b1", Name: "Breed1", Temperament: "T1", Origin: "O1"}}
	pets := []types.Pet{{ID: "p1", Name: "Fido", Birth: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), Breed: breeds[0]}}
	petStore := &PetStoreMock{pets: pets, breeds: breeds}
	breedStore := &BreedStoreMock{breeds: breeds}
	handler := NewPetHandler(petStore, breedStore)

	req, _ := http.NewRequest("GET", "/api/v1/pets", nil)
	rec := httptest.NewRecorder()
	handler.PetsHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
	if rec.Header().Get("Content-Type") != "application/json" {
		t.Errorf("expected application/json, got %s", rec.Header().Get("Content-Type"))
	}
	var got []types.Pet
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Errorf("error decoding response: %v", err)
	}
	if len(got) != 1 || got[0].ID != "p1" {
		t.Errorf("unexpected pets: %+v", got)
	}
}

func TestGetPetByIDHandler(t *testing.T) {
	breeds := []types.Breed{{ID: "b1", Name: "Breed1", Temperament: "T1", Origin: "O1"}}
	pets := []types.Pet{{ID: "p1", Name: "Fido", Birth: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), Breed: breeds[0]}}
	petStore := &PetStoreMock{pets: pets, breeds: breeds}
	breedStore := &BreedStoreMock{breeds: breeds}
	handler := NewPetHandler(petStore, breedStore)

	t.Run("found", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/pets/p1", nil)
		rec := httptest.NewRecorder()
		handler.GetPetByIDHandler(rec, req)
		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}
		var got types.Pet
		if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
			t.Errorf("error decoding: %v", err)
		}
		if got.ID != "p1" {
			t.Errorf("unexpected pet: %+v", got)
		}
	})

	t.Run("not found", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/pets/doesnotexist", nil)
		rec := httptest.NewRecorder()
		handler.GetPetByIDHandler(rec, req)
		if rec.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", rec.Code)
		}
		if rec.Body.String() != "Pet not found\n" {
			t.Errorf("unexpected body: %q", rec.Body.String())
		}
	})
}

func TestCreatePetHandler(t *testing.T) {
	breeds := []types.Breed{{ID: "b1", Name: "Breed1", Temperament: "T1", Origin: "O1"}}
	petStore := &PetStoreMock{breeds: breeds}
	breedStore := &BreedStoreMock{breeds: breeds}
	handler := NewPetHandler(petStore, breedStore)

	t.Run("success", func(t *testing.T) {
		reqBody := types.CreatePetRequest{
			Name:    "Fido",
			Birth:   "2020-01-01",
			BreedID: "b1",
		}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/api/v1/pets", bytes.NewReader(body))
		rec := httptest.NewRecorder()
		handler.PetsHandler(rec, req)
		if rec.Code != http.StatusCreated {
			t.Errorf("expected 201, got %d", rec.Code)
		}
		var got types.Pet
		if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
			t.Errorf("error decoding: %v", err)
		}
		if got.Name != "Fido" || got.Breed.ID != "b1" {
			t.Errorf("unexpected pet: %+v", got)
		}
	})

	t.Run("bad breed", func(t *testing.T) {
		reqBody := types.CreatePetRequest{
			Name:    "Fido",
			Birth:   "2020-01-01",
			BreedID: "doesnotexist",
		}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/api/v1/pets", bytes.NewReader(body))
		rec := httptest.NewRecorder()
		handler.PetsHandler(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rec.Code)
		}
		if rec.Body.String() != "Error at checking the breed\n" {
			t.Errorf("unexpected body: %q", rec.Body.String())
		}
	})

	t.Run("bad date", func(t *testing.T) {
		reqBody := types.CreatePetRequest{
			Name:    "Fido",
			Birth:   "not-a-date",
			BreedID: "b1",
		}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/api/v1/pets", bytes.NewReader(body))
		rec := httptest.NewRecorder()
		handler.PetsHandler(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rec.Code)
		}
		if rec.Body.String() != "Bad date of birth format. Use YYYY-MM-DD\n" {
			t.Errorf("unexpected body: %q", rec.Body.String())
		}
	})

	t.Run("bad json", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/api/v1/pets", bytes.NewReader([]byte("not-json")))
		rec := httptest.NewRecorder()
		handler.PetsHandler(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rec.Code)
		}
		if rec.Body.String() != "Error decoding the body of the request\n" {
			t.Errorf("unexpected body: %q", rec.Body.String())
		}
	})
}
