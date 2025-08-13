package types

import "time"

type Breed struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Temperament string `json:"temperament"`
	Origin      string `json:"origin"`
}

type Pet struct {
	ID    string    `json:"id"`
	Name  string    `json:"name"`
	Birth time.Time `json:"birth"`
	Breed Breed     `json:"breed"`
}

type CreatePetRequest struct {
	Name    string `json:"name"`
	Birth   string `json:"birth"` // La fecha se envía como un string, luego la convertiremos a time.Time
	BreedID string `json:"breedId"`
}
