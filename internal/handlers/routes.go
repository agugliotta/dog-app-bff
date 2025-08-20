package handlers

import (
	"github.com/agugliotta/dog-app-bff/internal/store"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, bs store.BreedStore, ps store.PetStore) {
	breedHandler := NewBreedHandler(bs)
	petHandler := NewPetHandler(ps, bs)
	groupV2 := router.Group("/api/v1")

	groupV2.GET("/breeds", breedHandler.GetBreedsHandler)
	groupV2.GET("/breeds/:slug", breedHandler.GetBreedBySlugHandler)

	groupV2.GET("/pets", petHandler.GetPetsHandler)
	groupV2.GET("/pets/:id", petHandler.GetPetByIDHandler)
	groupV2.DELETE("/pets/:id", petHandler.DeletePetById)
	groupV2.POST("/pets", petHandler.CreatePetHandler)
}
