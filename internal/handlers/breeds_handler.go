package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/agugliotta/dog-app-bff/internal/store"
	"github.com/gin-gonic/gin"
)

// BreedHandler es un struct que contendrá las dependencias (como el store) necesarias para los handlers de razas.
type BreedHandler struct {
	breedStore store.BreedStore
}

// NewBreedHandler crea e inicializa un nuevo BreedHandler.
// Es un constructor que nos permite "inyectar" el store.
func NewBreedHandler(bs store.BreedStore) *BreedHandler {
	return &BreedHandler{
		breedStore: bs,
	}

}

// GetBreedsHandler maneja las solicitudes HTTP para obtener la lista de razas.
// Es un método en el BreedHandler, lo que nos da acceso a 'h.breedStore'.
func (h *BreedHandler) GetBreedsHandler(c *gin.Context) {
	breeds, err := h.breedStore.GetBreeds()
	if err != nil {
		log.Printf("Error al obtener razas desde el store: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, breeds)
}

func (h *BreedHandler) GetBreedBySlugHandler(c *gin.Context) {
	slug := c.Param("slug")

	breed, err := h.breedStore.GetBreedBySlug(slug)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Breed not found", "message": err.Error()})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, breed)
}
