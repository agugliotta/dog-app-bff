package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/agugliotta/dog-app-bff/internal/store"
	"github.com/agugliotta/dog-app-bff/internal/types"
)

type StoreMock struct{}

func (sm *StoreMock) GetBreeds() ([]types.Breed, error) {
	return []types.Breed{
		{ID: "mock-breed-1", Name: "Mock Poodle", Temperament: "Mock Temp 1", Origin: "Mockland"},
		{ID: "mock-breed-2", Name: "Mock Bulldog", Temperament: "Mock Temp 2", Origin: "Mockland"},
	}, nil
}

// GetBreedByID implementa el método GetBreedByID de la interfaz BreedStore para el mock.
func (m *StoreMock) GetBreedByID(id string) (*types.Breed, error) {
	if id == "mock-breed-1" {
		return &types.Breed{ID: "mock-breed-1", Name: "Mock Poodle", Temperament: "Mock Temp 1", Origin: "Mockland"}, nil
	}
	if id == "non-existent-id" {
		return nil, store.ErrNotFound // ¡Ahora devuelve tu error sentinel!
	}
	// Podrías añadir un caso para simular un error interno de store también:
	// if id == "error-id" {
	//    return nil, errors.New("simulated internal store error")
	// }
	return nil, store.ErrNotFound // Default para IDs no definidos en el mock
}

func TestGetBreedsHandler(t *testing.T) {
	var breeds []types.Breed
	handler := NewBreedHandler(&StoreMock{})
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("GET", "/api/v1/breeds", nil)

	if err != nil {
		t.Fatalf("Error en el request")
	}

	handler.GetBreedsHandler(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Errorf("The Request was unsuccesful")
	}

	if recorder.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Wrong header: expected 'application/json', got '%s'", recorder.Header().Get("Content-Type"))
	}

	err = json.NewDecoder(recorder.Body).Decode(&breeds)
	if err != nil {
		t.Errorf("Error al decodificar la respuesta JSON: %v", err)
	}

	if len(breeds) != 2 {
		t.Errorf("Wrong lenght in the result array")
	}

	if breeds[0].ID != "mock-breed-1" {
		t.Errorf("Wrong id for first record, is %s and it should be %s", breeds[0].ID, "mock-breed-1")
	}
}

func TestGetBreedByIDHandler(t *testing.T) {
	t.Run("should return breed for existing ID", func(t *testing.T) {
		var breed types.Breed
		handler := NewBreedHandler(&StoreMock{})
		recorder := httptest.NewRecorder()
		request, err := http.NewRequest("GET", "/api/v1/breeds/mock-breed-1", nil)

		if err != nil {
			t.Fatalf("Error en el request")
		}

		handler.GetBreedByIDHandler(recorder, request)

		if recorder.Code != http.StatusOK {
			t.Errorf("The Request was unsuccesful")
		}

		if recorder.Header().Get("Content-Type") != "application/json" {
			t.Errorf("Wrong header: expected 'application/json', got '%s'", recorder.Header().Get("Content-Type"))
		}

		err = json.NewDecoder(recorder.Body).Decode(&breed)
		if err != nil {
			t.Errorf("Error al decodificar la respuesta JSON: %v", err)
		}

		expectedID := "mock-breed-1"
		expectedName := "Mock Poodle"
		if breed.ID != expectedID {
			t.Errorf("ID de raza incorrecto: esperado '%s', obtenido '%s'", expectedID, breed.ID)
		}
		if breed.Name != expectedName {
			t.Errorf("Nombre de raza incorrecto: esperado '%s', obtenido '%s'", expectedName, breed.Name)
		}
	})

	t.Run("should return 404 for non-existent ID", func(t *testing.T) {
		handler := NewBreedHandler(&StoreMock{})
		recorder := httptest.NewRecorder()
		request, err := http.NewRequest("GET", "/api/v1/breeds/non-existent-id", nil) // Este ID ahora causa ErrNotFound en el mock
		if err != nil {
			t.Fatalf("Error al crear la solicitud: %v", err)
		}

		handler.GetBreedByIDHandler(recorder, request)

		if recorder.Code != http.StatusNotFound { // Verifica el 404
			t.Errorf("Código de estado incorrecto para ID no existente: esperado %d, obtenido %d", http.StatusNotFound, recorder.Code)
		}
		expectedBody := "Breed not found\n" // Verifica el mensaje específico
		if recorder.Body.String() != expectedBody {
			t.Errorf("Cuerpo de respuesta de error incorrecto: esperado '%s', obtenido '%s'", expectedBody, recorder.Body.String())
		}
		if recorder.Header().Get("Content-Type") != "text/plain; charset=utf-8" { // http.Error devuelve text/plain por defecto
			t.Errorf("Encabezado Content-Type incorrecto: esperado 'text/plain; charset=utf-8', obtenido '%s'", recorder.Header().Get("Content-Type"))
		}
	})
}
