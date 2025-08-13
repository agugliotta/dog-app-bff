package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"path"
	"time"

	"github.com/agugliotta/dog-app-bff/internal/store"
	"github.com/agugliotta/dog-app-bff/internal/types"
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

func (ph *PetHandler) getPetsHandler(w http.ResponseWriter, r *http.Request) {
	pets, err := ph.petStore.GetPets()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(pets)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (ph *PetHandler) GetPetByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := path.Base(r.URL.Path)
	pet, err := ph.petStore.GetPetByID(id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			http.Error(w, "Pet not found", http.StatusNotFound)
			return
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(pet)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (ph *PetHandler) createPetHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var requestBody types.CreatePetRequest
	err := json.NewDecoder(r.Body).Decode(&requestBody)

	if err != nil {
		http.Error(w, "Error decoding the body of the request", http.StatusBadRequest)
		return
	}

	_, err = ph.breedStore.GetBreedByID(requestBody.BreedID)
	if err != nil {
		http.Error(w, "Error at checking the breed", http.StatusBadRequest)
		return
	}

	birth, err := time.Parse("2006-01-02", requestBody.Birth)
	if err != nil {
		http.Error(w, "Bad date of birth format. Use YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	newPet, err := ph.petStore.CreatePet(requestBody.Name, birth, requestBody.BreedID)
	if err != nil {
		log.Printf("Error creating pet in store: %v", err)
		http.Error(w, "Error creating pet", http.StatusInternalServerError)
		return
	}

	// 5. Enviar la respuesta exitosa.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // El código 201 es estándar para 'Created'.
	if err := json.NewEncoder(w).Encode(newPet); err != nil {
		log.Printf("Error encoding response for created pet: %v", err)
	}
}

func (ph *PetHandler) PetsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		ph.getPetsHandler(w, r)

	case http.MethodPost:
		ph.createPetHandler(w, r)

	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}
