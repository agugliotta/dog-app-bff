package breeds

import (
	"errors"
	"net/http"

	"github.com/agugliotta/dog-app-bff/internal/handlers/utils"
	"github.com/agugliotta/dog-app-bff/internal/store"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// BreedHandler contains dependencies for breed-related handlers.
type BreedHandler struct {
	logger     *zap.Logger
	breedStore store.BreedStore
}

// NewBreedHandler initializes a new BreedHandler with the given store and logger.
func NewBreedHandler(logger *zap.Logger, bs store.BreedStore) *BreedHandler {
	return &BreedHandler{
		logger:     logger,
		breedStore: bs,
	}
}

// GetBreedsHandler handles HTTP requests to retrieve the list of breeds.
func (h *BreedHandler) GetBreedsHandler(c *gin.Context) {
	breeds, err := h.breedStore.GetBreeds()
	if err != nil {
		h.logger.Error("Error getting breeds from store", zap.Error(err))
		utils.RespondError(c, http.StatusInternalServerError, utils.ErrInternalServer, err.Error())
		return
	}
	h.logger.Info("Fetched breeds", zap.Int("count", len(breeds)))
	c.JSON(http.StatusOK, breeds)
}

// GetBreedBySlugHandler handles HTTP requests to retrieve a breed by slug.
func (h *BreedHandler) GetBreedBySlugHandler(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		h.logger.Warn("Breed slug is required")
		utils.RespondError(c, http.StatusBadRequest, "Breed slug is required")
		return
	}
	breed, err := h.breedStore.GetBreedBySlug(slug)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			h.logger.Warn("Breed not found", zap.String("slug", slug))
			utils.RespondError(c, http.StatusNotFound, utils.ErrBreedNotFound, err.Error())
			return
		}
		h.logger.Error("Error getting breed by slug", zap.String("slug", slug), zap.Error(err))
		utils.RespondError(c, http.StatusInternalServerError, utils.ErrInternalServer, err.Error())
		return
	}
	h.logger.Info("Fetched breed", zap.String("slug", breed.Slug))
	c.JSON(http.StatusOK, breed)
}
