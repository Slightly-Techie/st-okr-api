package main

import (
	"log"
	"net/http"

	"github.com/Slightly-Techie/st-okr-api/db"
	"github.com/Slightly-Techie/st-okr-api/internal/routes"
	"github.com/Slightly-Techie/st-okr-api/provider"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func main() {
	database, err := db.InitDB()
	if err != nil {
		log.Fatalf("could not connect to db: %v", err)
	}

	validator := validator.New()

	provider := provider.NewProvider(database, validator)

	router := routes.SetupRouter(provider)
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello to the SlightlyTechie OKR API!"})
	})

	router.Run(":8080")
}
