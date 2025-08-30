package pets

import (
	"errors"
	"net/http"
	"time"

	"github.com/agugliotta/dog-app-bff/internal/handlers/utils"
	"github.com/agugliotta/dog-app-bff/internal/store"
	"github.com/agugliotta/dog-app-bff/internal/types"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type PetHandler struct {
	logger     *zap.Logger
	petStore   store.PetStore
	breedStore store.BreedStore
}

func NewPetHandler(logger *zap.Logger, ps store.PetStore, bs store.BreedStore) *PetHandler {
	return &PetHandler{
		logger:     logger,
		petStore:   ps,
		breedStore: bs,
	}
}

func (ph *PetHandler) GetPetsHandler(c *gin.Context) {
	pets, err := ph.petStore.GetPets()
	if err != nil {
		ph.logger.Error("Failed to get pets", zap.Error(err))
		utils.RespondError(c, http.StatusInternalServerError, utils.ErrInternalServer, err.Error())
		return
	}
	ph.logger.Info("Fetched pets", zap.Int("count", len(pets)))
	c.JSON(http.StatusOK, pets)
}

func (ph *PetHandler) GetPetByIDHandler(c *gin.Context) {
	id := c.Param("id")
	pet, err := ph.petStore.GetPetByID(id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			ph.logger.Warn("Pet not found", zap.String("pet_id", id))
			utils.RespondError(c, http.StatusNotFound, utils.ErrPetNotFound)
			return
		}
		ph.logger.Error("Failed to get pet by ID", zap.String("pet_id", id), zap.Error(err))
		utils.RespondError(c, http.StatusInternalServerError, utils.ErrInternalServer, err.Error())
		return
	}
	ph.logger.Info("Fetched pet", zap.String("pet_id", pet.ID))
	c.JSON(http.StatusOK, pet)
}

func (ph *PetHandler) CreatePetHandler(c *gin.Context) {
	var requestBody types.CreatePetRequest

	if err := c.BindJSON(&requestBody); err != nil {
		ph.logger.Warn("Invalid request body", zap.Error(err))
		utils.RespondError(c, http.StatusBadRequest, utils.ErrInvalidBody, err.Error())
		return
	}
	if requestBody.Name == "" {
		ph.logger.Warn("Pet name is required")
		utils.RespondError(c, http.StatusBadRequest, utils.ErrPetNameRequired)
		return
	}
	if requestBody.BreedID == "" {
		ph.logger.Warn("Breed ID is required")
		utils.RespondError(c, http.StatusBadRequest, utils.ErrBreedIDRequired)
		return
	}
	_, err := ph.breedStore.GetBreedByID(requestBody.BreedID)
	if err != nil {
		ph.logger.Warn("Breed not found", zap.String("breed_id", requestBody.BreedID), zap.Error(err))
		utils.RespondError(c, http.StatusBadRequest, utils.ErrBreedNotFound, err.Error())
		return
	}
	birth, err := time.Parse("2006-01-02", requestBody.Birth)
	if err != nil {
		ph.logger.Warn("Bad date format", zap.String("birth", requestBody.Birth), zap.Error(err))
		utils.RespondError(c, http.StatusBadRequest, utils.ErrBadDateFormat, err.Error())
		return
	}
	newPet, err := ph.petStore.CreatePet(requestBody.Name, birth, requestBody.BreedID)
	if err != nil {
		ph.logger.Error("Error creating pet in store", zap.Error(err))
		utils.RespondError(c, http.StatusInternalServerError, "Error creating pet", err.Error())
		return
	}
	ph.logger.Info("Created pet", zap.String("pet_id", newPet.ID), zap.String("name", newPet.Name))
	c.JSON(http.StatusCreated, newPet)
}

func (ph *PetHandler) DeletePetById(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		ph.logger.Warn("Pet ID is required")
		utils.RespondError(c, http.StatusBadRequest, "Pet ID is required")
		return
	}
	err := ph.petStore.DeletePet(id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			ph.logger.Warn("Pet not found for delete", zap.String("pet_id", id))
			utils.RespondError(c, http.StatusNotFound, utils.ErrPetNotFound, err.Error())
			return
		}
		ph.logger.Error("Error deleting pet", zap.String("pet_id", id), zap.Error(err))
		utils.RespondError(c, http.StatusInternalServerError, utils.ErrInternalServer, err.Error())
		return
	}
	ph.logger.Info("Deleted pet", zap.String("pet_id", id))
	c.JSON(http.StatusNoContent, nil)
}
