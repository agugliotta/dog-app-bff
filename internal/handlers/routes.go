package handlers

import (
	"strings"

	"github.com/agugliotta/dog-app-bff/internal/store"
	"github.com/gin-gonic/gin"
	ginprometheus "github.com/zsais/go-gin-prometheus"
	"go.uber.org/zap"
)

func RegisterRoutes(router *gin.Engine, logger *zap.Logger, bs store.BreedStore, ps store.PetStore) {
	breedHandler := NewBreedHandler(logger, bs)
	petHandler := NewPetHandler(logger, ps, bs)

	prometheus := ginprometheus.NewWithConfig(ginprometheus.Config{
		Subsystem: "gin",
	})
	prometheus.ReqCntURLLabelMappingFn = func(c *gin.Context) string {
		url := c.Request.URL.Path
		for _, p := range c.Params {
			if p.Key == "id" {
				url = strings.Replace(url, p.Value, ":id", 1)
				break
			} else if p.Key == "slug" {
				url = strings.Replace(url, p.Value, ":slug", 1)
				break
			}
		}
		return url
	}
	prometheus.Use(router)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	groupV2 := router.Group("/api/v1")

	groupV2.GET("/breeds", breedHandler.GetBreedsHandler)
	groupV2.GET("/breeds/:slug", breedHandler.GetBreedBySlugHandler)

	groupV2.GET("/pets", petHandler.GetPetsHandler)
	groupV2.GET("/pets/:id", petHandler.GetPetByIDHandler)
	groupV2.DELETE("/pets/:id", petHandler.DeletePetById)
	groupV2.POST("/pets", petHandler.CreatePetHandler)
}
