# Define variables para facilitar la modificación
PROJECT_NAME := dog-app-bff
DOCKER_DB_CONTAINER := $(PROJECT_NAME)-postgres-test
DOCKER_DB_PASSWORD := mysecretpassword
DOCKER_DB_NAME := dog_app_db_test
DB_PORT := 5432
GO_APP_PORT := 8080

# String de conexión a la base de datos para tests (usando la DB de test)
TEST_DB_CONN_STRING := "host=localhost port=$(DB_PORT) user=postgres password=$(DOCKER_DB_PASSWORD) dbname=$(DOCKER_DB_NAME) sslmode=disable"

# .PHONY: all clean run build test test-integration test-unit db-start db-stop db-clean db-setup-test test-integration-auto-db # Puedes listar todos los targets, o solo los públicos
.PHONY: all clean run build test test-unit test-integration db-start db-stop db-clean db-setup-test

all: build run

# Limpia los binarios y módulos Go
clean:
	@echo "Limpiando binarios y módulos Go..."
	@rm -f ./bin/$(PROJECT_NAME)
	@go clean -modcache
	@echo "Limpieza de Go finalizada."

# Ejecuta la aplicación Go (para desarrollo, no tests)
run: build
	@echo "Ejecutando la aplicación Go..."
	@DB_CONN_STRING=$(TEST_DB_CONN_STRING) go run ./cmd/api/main.go

# Compila la aplicación Go
build:
	@echo "Compilando la aplicación Go..."
	@go build -o ./bin/$(PROJECT_NAME) ./cmd/api/main.go
	@echo "Compilación finalizada. Binario en ./bin/$(PROJECT_NAME)"

# Ejecuta todos los tests (unidad e integración, asume que la DB de test está corriendo)
# Ahora 'test-integration' se encargará de la DB automáticamente
test: test-unit test-integration

# Ejecuta solo los tests unitarios (que no requieren DB)
test-unit:
	@echo "Ejecutando tests unitarios..."
	@go test -v ./internal/handlers/...
	@go test -v ./internal/types/... # si tuvieras tests aquí
	@echo "Tests unitarios finalizados."

# --- Nuevo Target para Tests de Integración que Autogestiona la DB ---
# Este target ahora depende de db-setup-test para iniciar y configurar la DB.
# Y usa 'db-teardown-test' para limpiar al final, incluso si los tests fallan.
test-integration: db-setup-test
	@echo "Iniciando tests de integración con DB gestionada automáticamente..."
	@trap 'make db-stop' EXIT; TEST_DB_CONN_STRING=$(TEST_DB_CONN_STRING) go test -v ./internal/store/...
	@echo "Tests de integración finalizados."

# --- Comandos relacionados con Docker y la Base de Datos de Test ---

# Inicia el contenedor de PostgreSQL para tests
db-start:
	@echo "Iniciando contenedor Docker para PostgreSQL de test..."
	@docker run --name $(DOCKER_DB_CONTAINER) \
		-e POSTGRES_PASSWORD=$(DOCKER_DB_PASSWORD) \
		-e POSTGRES_DB=$(DOCKER_DB_NAME) \
		-p $(DB_PORT):5432 \
		-d postgres:latest
	@echo "Esperando que PostgreSQL esté listo (esto puede tomar un momento)..."
	# Pequeño hack para esperar a que la DB esté lista. En producción usarías una tool como wait-for-it.sh
	@sleep 5
	@echo "Contenedor de PostgreSQL de test iniciado."

# Detiene el contenedor de PostgreSQL para tests
db-stop:
	@echo "Deteniendo contenedor Docker de PostgreSQL de test..."
	@docker stop $(DOCKER_DB_CONTAINER) > /dev/null 2>&1 || true # Redirige salida para no ver "no existe"
	@docker rm $(DOCKER_DB_CONTAINER) > /dev/null 2>&1 || true
	@echo "Contenedor de PostgreSQL de test detenido y eliminado."

# Este target se asegura que la DB esté limpia y configurada para cada ejecución de test-integration
db-setup-test: db-stop db-start
	@echo "Configurando la base de datos de test..."
	@sleep 2
	@docker exec -i $(DOCKER_DB_CONTAINER) psql -U postgres -d $(DOCKER_DB_NAME) -c " \
		DROP TABLE IF EXISTS pets CASCADE; \
		DROP TABLE IF EXISTS breeds CASCADE; \
		CREATE TABLE breeds ( \
			id VARCHAR(255) PRIMARY KEY, \
			name VARCHAR(255) NOT NULL, \
			temperament TEXT, \
			origin VARCHAR(255) \
		); \
		CREATE TABLE pets ( \
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(), \
			name VARCHAR(255) NOT NULL, \
			birth DATE NOT NULL, \
			breed_id VARCHAR(255) NOT NULL REFERENCES breeds(id) \
		); \
		INSERT INTO breeds (id, name, temperament, origin) VALUES \
		('golden-retriever', 'Golden Retriever', 'Friendly, Intelligent, Devoted', 'Scotland'), \
		('german-shepherd', 'German Shepherd', 'Intelligent, Obedient, Courageous', 'Germany'), \
		('poodle', 'Poodle', 'Intelligent, Proud, Active', 'Germany/France'), \
		('labrador-retriever', 'Labrador Retriever', 'Outgoing, Even-tempered, Gentle', 'Canada'), \
		('bulldog', 'Bulldog', 'Docile, Willful, Friendly', 'England'); \
		INSERT INTO pets (name, birth, breed_id) VALUES \
		('Buddy', '2022-05-10', 'golden-retriever'), \
		('Max', '2023-01-20', 'german-shepherd'); \
	"
	@echo "Base de datos de test configurada."