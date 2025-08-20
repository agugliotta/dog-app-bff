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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pets)
}

func (ph *PetHandler) GetPetByIDHandler(c *gin.Context) {
	id := c.Param("id")
	pet, err := ph.petStore.GetPetByID(id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Pet not found"})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
	}

	c.JSON(http.StatusOK, pet)
}

func (ph *PetHandler) CreatePetHandler(c *gin.Context) {
	var requestBody types.CreatePetRequest

	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := ph.breedStore.GetBreedByID(requestBody.BreedID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error at checking the breed"})
		return
	}

	birth, err := time.Parse("2006-01-02", requestBody.Birth)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad date of birth format. Use YYYY-MM-DD"})
		return
	}

	newPet, err := ph.petStore.CreatePet(requestBody.Name, birth, requestBody.BreedID)
	if err != nil {
		log.Printf("Error creating pet in store: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating pet"})
		return
	}

	c.JSON(http.StatusCreated, newPet)
}

func (ph *PetHandler) DeletePetById(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Pet ID is required"})
		return
	}

	err := ph.petStore.DeletePet(id)

	if err != nil {
		// Correctly handle the two main error cases
		if errors.Is(err, store.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Pet not found", "message": err.Error()})
			return
		}
		// Catch all other server-side errors
		log.Printf("Error deleting pet with ID %s: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
