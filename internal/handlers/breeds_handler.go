package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/agugliotta/dog-app-bff/internal/store"
	"github.com/gin-gonic/gin"
)

// BreedHandler contains dependencies for breed-related handlers.
type BreedHandler struct {
	breedStore store.BreedStore
}

// NewBreedHandler initializes a new BreedHandler with the given store.
func NewBreedHandler(bs store.BreedStore) *BreedHandler {
	return &BreedHandler{
		breedStore: bs,
	}
}

// GetBreedsHandler handles HTTP requests to retrieve the list of breeds.
func (h *BreedHandler) GetBreedsHandler(c *gin.Context) {
	breeds, err := h.breedStore.GetBreeds()
	if err != nil {
		log.Printf("Error getting breeds from store: %v", err)
		respondError(c, http.StatusInternalServerError, ErrInternalServer, err.Error())
		return
	}
	c.JSON(http.StatusOK, breeds)
}

// GetBreedBySlugHandler handles HTTP requests to retrieve a breed by slug.
func (h *BreedHandler) GetBreedBySlugHandler(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		respondError(c, http.StatusBadRequest, "Breed slug is required")
		return
	}
	breed, err := h.breedStore.GetBreedBySlug(slug)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			respondError(c, http.StatusNotFound, ErrBreedNotFound, err.Error())
			return
		}
		respondError(c, http.StatusInternalServerError, ErrInternalServer, err.Error())
		return
	}
	c.JSON(http.StatusOK, breed)
}
