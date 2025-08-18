# Define variables para el proyecto
PROJECT_NAME := dog-app-bff
GO_APP_PORT := 8080

# --- Variables del Entorno Único (para Tests y Desarrollo Local) ---
DB_CONTAINER := $(PROJECT_NAME)-postgres-test
DB_PASSWORD := mysecretpassword
DB_NAME := dog_app_db_test
DB_PORT := 5432
DB_CONN_STRING := "host=localhost port=$(DB_PORT) user=postgres password=$(DB_PASSWORD) dbname=$(DB_NAME) sslmode=disable"

.PHONY: all clean run-local run build test test-unit test-integration db-setup-if-needed db-setup db-start db-stop

all: build run-local

# Limpia los binarios y módulos Go
clean:
	@echo "Limpiando binarios y módulos Go..."
	@rm -f ./bin/$(PROJECT_NAME)
	@go clean -modcache
	@echo "Limpieza de Go finalizada."

# --- Target principal para Ejecutar la App en Local ---
run-local: build db-setup-if-needed
	@echo "Ejecutando la aplicación Go en el entorno de desarrollo..."
	@DB_CONN_STRING=$(DB_CONN_STRING) go run ./cmd/api/main.go

# Compila la aplicación Go
build:
	@echo "Compilando la aplicación Go..."
	@go build -o ./bin/$(PROJECT_NAME) ./cmd/api/main.go
	@echo "Compilación finalizada. Binario en ./bin/$(PROJECT_NAME)"

# Ejecuta todos los tests (unidad e integración)
test: test-unit test-integration

# Ejecuta solo los tests unitarios (que no requieren DB)
test-unit:
	@echo "Ejecutando tests unitarios..."
	@go test -v ./internal/handlers/...
	@go test -v ./internal/types/...
	@echo "Tests unitarios finalizados."

# --- Target para Tests de Integración que Autogestiona la DB ---
test-integration: db-setup-if-needed
	@echo "Iniciando tests de integración con la DB gestionada automáticamente..."
	@trap 'make db-stop' EXIT; DB_CONN_STRING=$(DB_CONN_STRING) go test -v ./internal/store/...
	@echo "Tests de integración finalizados."

# --- Comandos relacionados con Docker y la Base de Datos ---

db-setup-if-needed:
	@echo "Comprobando si el contenedor de la base de datos está activo..."
	@if ! docker ps -f name=$(DB_CONTAINER) | grep -q $(DB_CONTAINER); then \
		echo "El contenedor no está activo. Iniciándolo y configurándolo..."; \
		make db-setup; \
	else \
		echo "El contenedor ya está activo. Omitiendo el inicio."; \
	fi

# Inicia y configura la DB
db-setup: db-stop db-start
	@echo "Configurando la base de datos..."
	@sleep 2
	@docker cp init.sql $(DB_CONTAINER):/init.sql
	@docker exec -i $(DB_CONTAINER) psql -U postgres -d $(DB_NAME) -v ON_ERROR_STOP=1 -f /init.sql
	@echo "Base de datos configurada."

# Inicia el contenedor de PostgreSQL
db-start:
	@echo "Iniciando contenedor Docker para PostgreSQL..."
	@docker run --name $(DB_CONTAINER) \
        -e POSTGRES_PASSWORD=$(DB_PASSWORD) \
        -e POSTGRES_DB=$(DB_NAME) \
        -p $(DB_PORT):5432 \
        -d postgres:latest
	@echo "Esperando que PostgreSQL esté listo..."
	@sleep 5

# Detiene el contenedor de PostgreSQL
db-stop:
	@echo "Deteniendo contenedor Docker de PostgreSQL..."
	@docker stop $(DB_CONTAINER) > /dev/null 2>&1 || true
	@docker rm $(DB_CONTAINER) > /dev/null 2>&1 || true
	@echo "Contenedor de PostgreSQL detenido y eliminado."
