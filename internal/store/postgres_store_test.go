package store

import (
	"errors"
	"log"
	"os"
	"testing"
	"time"

	"github.com/agugliotta/dog-app-bff/internal/types"
)

func setupTestDB() *PostgresStore {
	// Usamos una variable de entorno específica para los tests, o la misma que en main si no hay.
	connStr := os.Getenv("TEST_DB_CONN_STRING")
	if connStr == "" {
		// Fallback para desarrollo local si TEST_DB_CONN_STRING no está configurada,
		// usa la misma que la app principal. En un CI/CD, TEST_DB_CONN_STRING sería obligatoria.
		connStr = os.Getenv("DB_CONN_STRING")
		if connStr == "" {
			log.Fatal("Las variables de entorno TEST_DB_CONN_STRING o DB_CONN_STRING no están configuradas.")
		}
	}

	store, err := NewPostgresStore(connStr)
	if err != nil {
		log.Fatalf("No se pudo conectar a la base de datos de prueba: %v", err)
	}

	// Opcional: Limpiar o sembrar la base de datos de prueba antes de los tests
	// (¡MUY importante para tests de integración reales!)
	// executeSQL(store.db, "DELETE FROM breeds") // Limpia la tabla
	// executeSQL(store.db, "INSERT INTO breeds (id, name, temperament, origin) VALUES ('test-breed', 'Test Breed', 'Test', 'Testland')") // Si quieres un dato específico para todos los tests.
	// executeSQL es una función auxiliar para ejecutar SQL simple.
	// func executeSQL(db *sql.DB, query string) {
	// 	_, err := db.Exec(query)
	// 	if err != nil {
	// 		log.Fatalf("Error al ejecutar SQL en la base de datos de prueba: %v", err)
	// 	}
	// }
	return store
}

// TestMain permite realizar configuraciones y limpiezas globales para los tests de este paquete.
func TestMain(m *testing.M) {
	// Puedes configurar tu base de datos de prueba aquí si es un set-up muy costoso.
	// Por simplicidad, setupTestDB se llama en cada test si es unitario,
	// o puedes usar una instancia global si los tests son independientes.

	// Normalmente aquí iniciarías un contenedor de base de datos específico para tests
	// o harías cualquier configuración de una sola vez.

	code := m.Run() // Ejecuta todos los tests en el paquete

	// Limpieza después de que todos los tests se hayan ejecutado
	// (ej. detener el contenedor de DB de prueba, si lo iniciaste aquí).

	os.Exit(code)
}

// TestGetBreeds verifica que podemos obtener razas desde el store de PostgreSQL.
func TestGetBreeds(t *testing.T) {
	// setupTestDB obtiene una conexión a la base de datos.
	// Para un test unitario, podríamos usar un mock o stub aquí.
	// Para un test de integración, usamos la DB real.
	// Dado que `NewPostgresStore` se conecta a una DB real, esto es un test de integración.
	store := setupTestDB()
	defer store.db.Close() // Cierra la conexión después de que termine el test.

	// Asegurémonos de que la tabla tenga al menos los datos base que insertamos.
	// Si estás ejecutando tests repetidamente sin limpiar la DB, es posible que los datos se dupliquen,
	// lo cual es una razón para usar una DB de test separada o limpiar antes de cada test.

	breeds, err := store.GetBreeds()
	if err != nil {
		t.Fatalf("GetBreeds falló: %v", err)
	}

	if len(breeds) == 0 {
		t.Errorf("GetBreeds devolvió 0 razas, esperaba al menos una.")
	}

	// Podemos verificar si una raza específica que esperamos está en la lista.
	foundGolden := false
	for _, breed := range breeds {
		if breed.Slug == "golden-retriever" && breed.Name == "Golden Retriever" {
			foundGolden = true
			break
		}
	}
	if !foundGolden {
		t.Errorf("No se encontró 'Golden Retriever' en las razas obtenidas.")
	}

	// Opcional: verificar la cantidad exacta si los datos son fijos para el test.
	// if len(breeds) != 5 {
	//     t.Errorf("Esperaba 5 razas, obtuve %d", len(breeds))
	// }
}

func TestGetBreedBySlug(t *testing.T) {
	store := setupTestDB()
	defer store.Close()

	slugToFind := "golden-retriever"
	breed, err := store.GetBreedBySlug(slugToFind)
	if err != nil {
		t.Fatalf("failed to get breed by slug: %v", err)
	}
	if breed == nil {
		t.Fatalf("breed with slug %s not found", slugToFind)
	}

	if breed.Slug != slugToFind {
		t.Errorf("expected slug %s, got %s", slugToFind, breed.Slug)
	}
	if breed.Name != "Golden Retriever" {
		t.Errorf("expected name %s, got %s", "Golden Retriever", breed.Name)
	}
}

func TestGetPets(t *testing.T) {
	store := setupTestDB()
	defer store.db.Close()

	pets, err := store.GetPets()
	if err != nil {
		t.Fatalf("GetPets failed: %v", err)
	}

	if len(pets) == 0 {
		t.Errorf("GetPets devolvió 0 mascotas, esperaba al menos una.")
	}
}

func TestPets(t *testing.T) {
	var id string
	store := setupTestDB()
	defer store.db.Close()

	t.Run("should create a new pet", func(t *testing.T) {
		breeds, _ := store.GetBreeds()
		newPet := types.Pet{ID: "", Name: "Fido", Birth: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), Breed: breeds[0]}
		petWithID, err := store.CreatePet(newPet.Name, newPet.Birth, newPet.Breed.ID)

		if err != nil && !errors.Is(err, ErrNotFound) {
			t.Errorf("error in creating new pet:'%v'", err)
		}

		id = petWithID.ID
	})

	t.Run("should return the new created pet by id", func(t *testing.T) {
		_, err := store.GetPetByID(id)

		if err != nil && !errors.Is(err, ErrNotFound) {
			t.Errorf("error new pet not found:'%v'", err)
		}
	})

	t.Run("should delete the new created pet by id", func(t *testing.T) {
		err := store.DeletePet(id)

		if err != nil && !errors.Is(err, ErrNotFound) {
			t.Errorf("error new pet not found:'%v'", err)
		}
	})

}
