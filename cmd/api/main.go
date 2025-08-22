package main

import (
	"log"
	"os"

	"github.com/agugliotta/dog-app-bff/internal/handlers"
	"github.com/agugliotta/dog-app-bff/internal/store"
	"github.com/gin-gonic/gin"
)

func main() {
	connStr := os.Getenv("DB_CONN_STRING")
	if connStr == "" {
		log.Fatal("La variable de entorno DB_CONN_STRING no está configurada. Por favor, configúrala.")
	}

	pgStore, err := store.NewPostgresStore(connStr)
	if err != nil {
		log.Fatalf("Error al inicializar el store de PostgreSQL: %v", err)
	}
	defer pgStore.Close()

	router := gin.Default()
	handlers.RegisterRoutes(router, pgStore, pgStore)

	log.Printf("Servidor iniciando en :8080...")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("El servidor falló al iniciar: %v", err)
	}
}
