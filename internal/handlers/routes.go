package handlers

import (
	"net/http"

	"github.com/agugliotta/dog-app-bff/internal/store"
)

// RegisterRoutes es la función principal para registrar todos los handlers con el router HTTP.
// Recibe el http.ServeMux estándar y el store de la aplicación.
func RegisterRoutes(router *http.ServeMux, bs store.BreedStore, ps store.PetStore) {
	// Crea una instancia del handler de razas, inyectando el store.
	breedHandler := NewBreedHandler(bs)
	petHandler := NewPetHandler(ps, bs)

	// Registra el handler para la ruta /api/v1/breeds.
	// Como usamos http.ServeMux, no especificamos métodos aquí, se hará dentro del handler si es necesario.
	router.HandleFunc("/api/v1/breeds", breedHandler.GetBreedsHandler)

	// Ruta para obtener una raza por ID
	// La barra final es CRUCIAL para que ServeMux capture cualquier cosa después.
	// Nota: Si una solicitud es exactamente "/api/v1/breeds", GetBreedsHandler la maneja.
	// Si es "/api/v1/breeds/algo", GetBreedByIDHandler la maneja.
	router.HandleFunc("/api/v1/breeds/", breedHandler.GetBreedByIDHandler)

	router.HandleFunc("/api/v1/pets/", petHandler.GetPetByIDHandler)
	router.HandleFunc("/api/v1/pets", petHandler.PetsHandler)
}
