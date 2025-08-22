package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/agugliotta/dog-app-bff/internal/store"
	"github.com/agugliotta/dog-app-bff/internal/types"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type StoreMock struct{}

func (sm *StoreMock) GetBreeds() ([]types.Breed, error) {
	return []types.Breed{
		{ID: "mock-breed-1", Name: "Mock Poodle", Temperament: "Mock Temp 1", Origin: "Mockland", Slug: "mock-poodle"},
		{ID: "mock-breed-2", Name: "Mock Bulldog", Temperament: "Mock Temp 2", Origin: "Mockland", Slug: "mock-bulldog"},
	}, nil
}

func (m *StoreMock) GetBreedByID(id string) (*types.Breed, error) {
	if id == "mock-breed-1" {
		return &types.Breed{ID: "mock-breed-1", Name: "Mock Poodle", Temperament: "Mock Temp 1", Origin: "Mockland", Slug: "mock-poodle"}, nil
	}
	if id == "non-existent-id" {
		return nil, store.ErrNotFound
	}
	return nil, store.ErrNotFound
}

func (m *StoreMock) GetBreedBySlug(slug string) (*types.Breed, error) {
	if slug == "mock-poodle" {
		return &types.Breed{ID: "mock-breed-1", Name: "Mock Poodle", Temperament: "Mock Temp 1", Origin: "Mockland", Slug: "mock-poodle"}, nil
	}
	if slug == "non-existent-slug" {
		return nil, store.ErrNotFound
	}
	return nil, store.ErrNotFound
}

func TestGetBreedsHandler(t *testing.T) {
	var breeds []types.Breed
	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/api/v1/breeds", nil)
	assert.NoError(t, err)
	router := gin.Default()
	RegisterRoutes(router, &StoreMock{}, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "The Request was unsuccessful")
	assert.Equal(t, ContentTypeExpected, w.Header().Get("Content-Type"), "Wrong header")

	err = json.NewDecoder(w.Body).Decode(&breeds)
	assert.NoError(t, err, "Error decoding JSON response")

	assert.Len(t, breeds, 2, "Wrong length in the result array")
	assert.Equal(t, "mock-breed-1", breeds[0].ID, "Wrong id for first record")
}

func TestGetBreedByIDHandler(t *testing.T) {
	t.Run("should return breed for existing Slug", func(t *testing.T) {
		var breed types.Breed
		w := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/api/v1/breeds/mock-poodle", nil)
		assert.NoError(t, err)
		router := gin.Default()
		RegisterRoutes(router, &StoreMock{}, nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "The Request was unsuccessful")
		assert.Equal(t, ContentTypeExpected, w.Header().Get("Content-Type"), "Wrong header")

		err = json.NewDecoder(w.Body).Decode(&breed)
		assert.NoError(t, err, "Error decoding JSON response")

		assert.Equal(t, "mock-poodle", breed.Slug, "Incorrect breed slug")
		assert.Equal(t, "Mock Poodle", breed.Name, "Incorrect breed name")
	})

	t.Run("should return 404 for non-existent ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/api/v1/breeds/non-existent-id", nil)
		assert.NoError(t, err)
		router := gin.Default()
		RegisterRoutes(router, &StoreMock{}, nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code, "Incorrect status code for non-existent ID")
		var resp map[string]interface{}
		err = json.NewDecoder(w.Body).Decode(&resp)
		assert.NoError(t, err, "error decoding response")
		assert.Equal(t, "Breed not found", resp["error"], "Incorrect error response body")
		assert.Equal(t, ContentTypeExpected, w.Header().Get("Content-Type"), "Incorrect Content-Type header")
	})
}
