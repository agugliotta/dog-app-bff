package main

import (
	"log"
	"net/http"
	"os"

	"github.com/agugliotta/dog-app-bff/internal/handlers"
	"github.com/agugliotta/dog-app-bff/internal/store"
)

// APIServer representa nuestra aplicación de servidor HTTP.
// Contiene la dirección de escucha y una referencia a nuestro store de datos.
type APIServer struct {
	addr       string
	breedStore store.BreedStore // Nuestra interfaz de store, que será una instancia de PostgresStore
	petStore   store.PetStore
}

// NewAPIServer crea una nueva instancia de APIServer.
// Recibe la dirección en la que escuchará y la implementación del store a usar.
func NewAPIServer(addr string, bs store.BreedStore, ps store.PetStore) *APIServer {
	return &APIServer{
		addr:       addr,
		breedStore: bs,
		petStore:   ps,
	}
}

// Run inicia el servidor HTTP.
func (s *APIServer) Run() {
	// Inicializa el router estándar de Go.
	router := http.NewServeMux()

	// Registra todas nuestras rutas, pasando el router y el store.
	handlers.RegisterRoutes(router, s.breedStore, s.petStore)

	log.Printf("Servidor iniciando en %s...", s.addr)
	// Inicia el servidor HTTP, usando nuestro router para manejar las solicitudes.
	err := http.ListenAndServe(s.addr, router)
	if err != nil {
		log.Fatalf("El servidor falló al iniciar: %v", err)
	}
}

func main() {
	// 1. Obtener la cadena de conexión de PostgreSQL de una variable de entorno.
	connStr := os.Getenv("DB_CONN_STRING")
	if connStr == "" {
		log.Fatal("La variable de entorno DB_CONN_STRING no está configurada. Por favor, configúrala.")
	}

	// 2. Inicializar el store de PostgreSQL.
	// Esto establece la conexión a la base de datos.
	pgStore, err := store.NewPostgresStore(connStr)
	if err != nil {
		log.Fatalf("Error al inicializar el store de PostgreSQL: %v", err)
	}
	// Asegúrate de cerrar la conexión a la base de datos cuando la aplicación se detenga.
	defer pgStore.Close() // Esto se ejecutará cuando main() termine.

	// 3. Crear una nueva instancia de APIServer, inyectando el store de PostgreSQL.
	server := NewAPIServer(":8080", pgStore, pgStore)

	// 4. Iniciar el servidor.
	server.Run()
}
