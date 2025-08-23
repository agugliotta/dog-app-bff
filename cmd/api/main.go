package main

import (
	"os"

	"github.com/agugliotta/dog-app-bff/internal/handlers"
	"github.com/agugliotta/dog-app-bff/internal/store"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	connStr := os.Getenv("DB_CONN_STRING")
	if connStr == "" {
		logger.Fatal("La variable de entorno DB_CONN_STRING no está configurada. Por favor, configúrala.")
	}

	pgStore, err := store.NewPostgresStore(connStr)
	if err != nil {
		logger.Fatal("Error al inicializar el store de PostgreSQL", zap.Error(err))
	}
	defer pgStore.Close()

	router := gin.Default()
	handlers.RegisterRoutes(router, logger, pgStore, pgStore)

	logger.Info("Servidor iniciando", zap.String("address", ":8080"))
	if err := router.Run(":8080"); err != nil {
		logger.Fatal("El servidor falló al iniciar", zap.Error(err))
	}
}
