package handlers

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/agugliotta/dog-app-bff/internal/store"
	"github.com/agugliotta/dog-app-bff/internal/types"
	"github.com/gin-gonic/gin"
)

type PetHandler struct {
	petStore   store.PetStore
	breedStore store.BreedStore
}

func NewPetHandler(ps store.PetStore, bs store.BreedStore) *PetHandler {
	return &PetHandler{
		petStore:   ps,
		breedStore: bs,
	}
}

func (ph *PetHandler) GetPetsHandler(c *gin.Context) {
	pets, err := ph.petStore.GetPets()
	if err != nil {
		respondError(c, http.StatusInternalServerError, ErrInternalServer, err.Error())
		return
	}
	c.JSON(http.StatusOK, pets)
}

func (ph *PetHandler) GetPetByIDHandler(c *gin.Context) {
	id := c.Param("id")
	pet, err := ph.petStore.GetPetByID(id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			respondError(c, http.StatusNotFound, ErrPetNotFound)
			return
		}
		respondError(c, http.StatusInternalServerError, ErrInternalServer, err.Error())
		return
	}
	c.JSON(http.StatusOK, pet)
}

func (ph *PetHandler) CreatePetHandler(c *gin.Context) {
	var requestBody types.CreatePetRequest

	if err := c.BindJSON(&requestBody); err != nil {
		respondError(c, http.StatusBadRequest, ErrInvalidBody, err.Error())
		return
	}
	if requestBody.Name == "" {
		respondError(c, http.StatusBadRequest, ErrPetNameRequired)
		return
	}
	if requestBody.BreedID == "" {
		respondError(c, http.StatusBadRequest, ErrBreedIDRequired)
		return
	}
	_, err := ph.breedStore.GetBreedByID(requestBody.BreedID)
	if err != nil {
		respondError(c, http.StatusBadRequest, ErrBreedNotFound, err.Error())
		return
	}
	birth, err := time.Parse("2006-01-02", requestBody.Birth)
	if err != nil {
		respondError(c, http.StatusBadRequest, ErrBadDateFormat, err.Error())
		return
	}
	newPet, err := ph.petStore.CreatePet(requestBody.Name, birth, requestBody.BreedID)
	if err != nil {
		log.Printf("Error creating pet in store: %v", err)
		respondError(c, http.StatusInternalServerError, "Error creating pet", err.Error())
		return
	}
	c.JSON(http.StatusCreated, newPet)
}

func (ph *PetHandler) DeletePetById(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		respondError(c, http.StatusBadRequest, "Pet ID is required")
		return
	}
	err := ph.petStore.DeletePet(id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			respondError(c, http.StatusNotFound, ErrPetNotFound, err.Error())
			return
		}
		log.Printf("Error deleting pet with ID %s: %v", id, err)
		respondError(c, http.StatusInternalServerError, ErrInternalServer, err.Error())
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
