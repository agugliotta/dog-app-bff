package utils

import (
	"github.com/gin-gonic/gin"
)

const (
	ErrInternalServer  = "Internal Server Error"
	ErrPetNotFound     = "Pet not found"
	ErrBreedNotFound   = "Breed not found"
	ErrInvalidBody     = "Invalid request body"
	ErrBadDateFormat   = "Bad date of birth format. Use YYYY-MM-DD"
	ErrPetNameRequired = "Pet name is required"
	ErrBreedIDRequired = "Breed ID is required"
)

func RespondError(c *gin.Context, status int, errMsg string, details ...string) {
	resp := gin.H{"error": errMsg}
	if len(details) > 0 {
		resp["message"] = details[0]
	}
	c.JSON(status, resp)
}
