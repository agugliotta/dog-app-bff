package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"path"

	"github.com/agugliotta/dog-app-bff/internal/store"
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
func (h *BreedHandler) GetBreedsHandler(w http.ResponseWriter, r *http.Request) {
	// Obtenemos las razas desde nuestro store.
	breeds, err := h.breedStore.GetBreeds()
	if err != nil {
		log.Printf("Error al obtener razas desde el store: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Configura el encabezado Content-Type para indicar que la respuesta es JSON.
	w.Header().Set("Content-Type", "application/json")

	// Codifica la slice de razas a JSON y la escribe en la respuesta HTTP.
	err = json.NewEncoder(w).Encode(breeds)
	if err != nil {
		log.Printf("Error al codificar razas a JSON: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (h *BreedHandler) GetBreedByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := path.Base(r.URL.Path)

	breed, err := h.breedStore.GetBreedByID(id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			http.Error(w, "Breed not found", http.StatusNotFound)
			return
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(breed)
	if err != nil {
		log.Printf("Error al codificar raza a JSON: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

}
