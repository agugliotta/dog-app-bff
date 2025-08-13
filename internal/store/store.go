package store

import (
	"errors"
	"time"

	"github.com/agugliotta/dog-app-bff/internal/types"
)

var (
	ErrNotFound            = errors.New("not found")
	ErrForeignKeyViolation = errors.New("foreign key violation")
)

type BreedStore interface {
	GetBreedByID(id string) (*types.Breed, error)
	GetBreeds() ([]types.Breed, error)
}

type PetStore interface {
	GetPets() ([]types.Pet, error)
	GetPetByID(id string) (*types.Pet, error)
	CreatePet(name string, birth time.Time, breedID string) (*types.Pet, error)
}
