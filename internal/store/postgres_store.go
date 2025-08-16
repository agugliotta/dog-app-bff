package store

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"

	"github.com/agugliotta/dog-app-bff/internal/types"
)

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(connStr string) (*PostgresStore, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Conectado exitosamente a PostgreSQL!")
	return &PostgresStore{db: db}, nil
}

func (s *PostgresStore) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

// BREEDS
func (s *PostgresStore) GetBreeds() ([]types.Breed, error) {
	rows, err := s.db.Query("SELECT id, name, temperament, origin FROM breeds")
	if err != nil {
		return nil, fmt.Errorf("failed to query breeds: %w", err)
	}
	defer rows.Close()

	var breeds []types.Breed
	for rows.Next() {
		var breed types.Breed
		if err := rows.Scan(&breed.ID, &breed.Name, &breed.Temperament, &breed.Origin); err != nil {
			return nil, fmt.Errorf("failed to scan breed: %w", err)
		}
		breeds = append(breeds, breed)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return breeds, nil
}

func (s *PostgresStore) GetBreedByID(id string) (*types.Breed, error) {
	var breed types.Breed
	err := s.db.QueryRow("SELECT id, name, temperament, origin FROM breeds WHERE id=$1", id).Scan(&breed.ID, &breed.Name, &breed.Temperament, &breed.Origin)

	switch err { // El switch ya manejará los diferentes tipos de error de 'err'
	case sql.ErrNoRows:
		// Es crucial que el error para "no encontrado" sea específico
		return nil, ErrNotFound
	case nil: // Si err es nil, significa que todo fue exitoso
		return &breed, nil
	default: // Cualquier otro tipo de error de la base de datos
		return nil, fmt.Errorf("query error for breed ID %s: %w", id, err)
	}
}

// PETS
func (s *PostgresStore) GetPets() ([]types.Pet, error) {
	query := `
		SELECT
            p.id, p.name, p.birth,
            b.id AS breed_id, b.name AS breed_name, b.temperament AS breed_temperament, b.origin AS breed_origin
        FROM
            pets p
        JOIN
            breeds b ON p.breed_id = b.id
	`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query pets: %w", err)
	}
	defer rows.Close()

	var pets []types.Pet
	for rows.Next() {
		var pet types.Pet
		var breed types.Breed
		if err := rows.Scan(&pet.ID, &pet.Name, &pet.Birth, &breed.ID, &breed.Name, &breed.Temperament, &breed.Origin); err != nil {
			return nil, fmt.Errorf("failed to scan pet or breed: %w", err)
		}
		pet.Breed = breed
		pets = append(pets, pet)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return pets, nil
}

func (s *PostgresStore) GetPetByID(id string) (*types.Pet, error) {
	var pet types.Pet
	var breed types.Breed
	query := `
		SELECT
            p.id, p.name, p.birth,
            b.id AS breed_id, b.name AS breed_name, b.temperament AS breed_temperament, b.origin AS breed_origin
        FROM
            pets p
        JOIN
            breeds b ON p.breed_id = b.id
		WHERE
			p.id=$1
	`
	err := s.db.QueryRow(query, id).Scan(&pet.ID, &pet.Name, &pet.Birth, &breed.ID, &breed.Name, &breed.Temperament, &breed.Origin)
	switch err {
	case sql.ErrNoRows:
		return nil, ErrNotFound
	case nil:
		pet.Breed = breed
		return &pet, nil
	default:
		return nil, fmt.Errorf("query error for breed ID %s: %w", id, err)
	}
}

func (s *PostgresStore) CreatePet(name string, birth time.Time, breedID string) (*types.Pet, error) {
	// Paso 1: Validar si la raza existe. Reutilizamos el método GetBreedByID.
	breed, err := s.GetBreedByID(breedID)
	if err != nil {
		// Si GetBreedByID falla, propagamos ese error directamente.
		// Podría ser ErrNotFound, u otro error interno.
		return nil, fmt.Errorf("failed to get breed with ID %s: %w", breedID, err)
	}

	query := `
        INSERT INTO pets (name, birth, breed_id)
        VALUES ($1, $2, $3)
        RETURNING id
    `

	var newPetID string

	err = s.db.QueryRow(query, name, birth, breedID).Scan(&newPetID)
	if err != nil {
		// No uses log.Fatalf. Devuelve el error para que el llamador lo maneje.
		return nil, fmt.Errorf("failed to insert pet and get ID: %w", err)
	}
	newPet := &types.Pet{
		ID:    newPetID,
		Name:  name,
		Birth: birth,
		Breed: *breed,
	}

	return newPet, nil
}

func (s *PostgresStore) DeletePet(id string) error {
	_, err := s.db.Exec("DELETE FROM pets WHERE id=$1", id)

	if err != nil {
		return err
	}
	return nil
}
