package store

import (
	"errors"
	"log"
	"os"
	"testing"
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
		if breed.ID == "golden-retriever" && breed.Name == "Golden Retriever" {
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

func TestGetBreedByID(t *testing.T) {
	store := setupTestDB()
	defer store.db.Close()

	t.Run("should return breed for existing ID", func(t *testing.T) {
		store := setupTestDB()
		// No deferred store.db.Close() here, it should be managed by TestMain or suite setup/teardown
		// For subtests, it's safer to ensure a clean slate, so if setupTestDB creates a new connection,
		// it's okay to close it. If it gets from a pool, don't close here.
		// Given your setupTestDB, it likely creates a new one, so defer close is fine here for isolation.
		defer store.db.Close()

		idToFind := "golden-retriever" // Asegúrate de que este ID esté en tu db-setup-test
		breed, err := store.GetBreedByID(idToFind)
		if err != nil {
			t.Fatalf("GetBreedByID falló para ID '%s': %v", idToFind, err)
		}

		if breed == nil {
			t.Fatalf("GetBreedByID devolvió nil para ID existente '%s'", idToFind)
		}
		if breed.ID != idToFind {
			t.Errorf("ID de raza incorrecto: esperado '%s', obtenido '%s'", idToFind, breed.ID)
		}
		if breed.Name != "Golden Retriever" {
			t.Errorf("Nombre de raza incorrecto: esperado 'Golden Retriever', obtenido '%s'", breed.Name)
		}
		// ... (más aserciones)
	})

	// Escenario 2: Raza no existente (verificando store.ErrNotFound)
	t.Run("should return ErrNotFound for non-existent ID", func(t *testing.T) {
		store := setupTestDB()
		defer store.db.Close()

		idToFind := "non-existent-breed-123" // ID que sabes que no está en la DB
		_, err := store.GetBreedByID(idToFind)

		if err == nil {
			t.Errorf("GetBreedByID debería haber devuelto un error para ID no existente '%s', pero devolvió nil", idToFind)
		}

		// ¡Aserción clave! Usar errors.Is para verificar el error sentinel
		if !errors.Is(err, ErrNotFound) { // Asegúrate de importar "errors" aquí si no está
			t.Errorf("Tipo de error incorrecto para ID no existente: esperado 'store.ErrNotFound', obtenido '%v'", err)
		}
	})

}
