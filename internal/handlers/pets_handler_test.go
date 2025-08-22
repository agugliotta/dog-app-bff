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
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
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

func (m *BreedStoreMock) GetBreedBySlug(slug string) (*types.Breed, error) {
	for _, b := range m.breeds {
		if b.Slug == slug {
			return &b, nil
		}
	}
	return nil, store.ErrNotFound
}

func TestGetPetsHandler(t *testing.T) {
	breeds := []types.Breed{{ID: "b1", Name: "Breed1", Temperament: "T1", Origin: "O1"}}
	pets := []types.Pet{{ID: "p1", Name: "Fido", Birth: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), Breed: breeds[0]}}
	petStore := &PetStoreMock{pets: pets, breeds: breeds}
	breedStore := &BreedStoreMock{breeds: breeds}
	router := gin.Default()
	RegisterRoutes(router, zap.NewNop(), breedStore, petStore)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/pets", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	assert.Equal(t, ContentTypeExpected, w.Header().Get("Content-Type"))

	var got []types.Pet
	err := json.NewDecoder(w.Body).Decode(&got)
	assert.NoError(t, err, "error decoding response")
	assert.Len(t, got, 1, "unexpected number of pets")
	assert.Equal(t, "p1", got[0].ID, "unexpected pet ID")
}

func TestGetPetByIDHandler(t *testing.T) {
	breeds := []types.Breed{{ID: "b1", Name: "Breed1", Temperament: "T1", Origin: "O1"}}
	pets := []types.Pet{{ID: "p1", Name: "Fido", Birth: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), Breed: breeds[0]}}
	petStore := &PetStoreMock{pets: pets, breeds: breeds}
	breedStore := &BreedStoreMock{breeds: breeds}
	router := gin.Default()
	RegisterRoutes(router, zap.NewNop(), breedStore, petStore)

	t.Run("found", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/pets/p1", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var got types.Pet
		err := json.NewDecoder(w.Body).Decode(&got)
		assert.NoError(t, err, "error decoding response")
		assert.Equal(t, "p1", got.ID, "unexpected pet ID")
	})

	t.Run("not found", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/pets/doesnotexist", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var resp map[string]interface{}
		err := json.NewDecoder(w.Body).Decode(&resp)
		assert.NoError(t, err, "error decoding response")
		assert.Equal(t, ErrPetNotFound, resp["error"], "unexpected error message")
	})
}

func TestCreatePetHandler(t *testing.T) {
	breeds := []types.Breed{{ID: "b1", Name: "Breed1", Temperament: "T1", Origin: "O1"}}
	petStore := &PetStoreMock{breeds: breeds}
	breedStore := &BreedStoreMock{breeds: breeds}
	router := gin.Default()
	RegisterRoutes(router, zap.NewNop(), breedStore, petStore)

	t.Run("success", func(t *testing.T) {
		reqBody := types.CreatePetRequest{
			Name:    "Fido",
			Birth:   "2020-01-01",
			BreedID: "b1",
		}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/api/v1/pets", bytes.NewReader(body))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code, "expected status 201")
		var got types.Pet
		err := json.NewDecoder(w.Body).Decode(&got)
		assert.NoError(t, err, "error decoding response")
		assert.Equal(t, "Fido", got.Name, "unexpected pet name")
		assert.Equal(t, "b1", got.Breed.ID, "unexpected breed ID")
	})

	t.Run("bad breed", func(t *testing.T) {
		reqBody := types.CreatePetRequest{
			Name:    "Fido",
			Birth:   "2020-01-01",
			BreedID: "doesnotexist",
		}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/api/v1/pets", bytes.NewReader(body))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code, "expected status 400")
		var resp map[string]interface{}
		err := json.NewDecoder(w.Body).Decode(&resp)
		assert.NoError(t, err, "error decoding response")
		assert.Equal(t, ErrBreedNotFound, resp["error"], "unexpected error message")
	})

	t.Run("bad date", func(t *testing.T) {
		reqBody := types.CreatePetRequest{
			Name:    "Fido",
			Birth:   "not-a-date",
			BreedID: "b1",
		}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/api/v1/pets", bytes.NewReader(body))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code, "expected status 400")
		var resp map[string]interface{}
		err := json.NewDecoder(w.Body).Decode(&resp)
		assert.NoError(t, err, "error decoding response")
		assert.Equal(t, ErrBadDateFormat, resp["error"], "unexpected error message")
	})

	t.Run("bad json", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/api/v1/pets", bytes.NewReader([]byte("not-json")))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code, "expected status 400")
		var resp map[string]interface{}
		err := json.NewDecoder(w.Body).Decode(&resp)
		assert.NoError(t, err, "error decoding response")
		assert.Contains(t, resp, "error", "expected error field in response")
	})
}
