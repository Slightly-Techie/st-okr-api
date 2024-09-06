package main

import (
	"log"

	"github.com/Slightly-Techie/st-okr-api/internal/config"
	"github.com/Slightly-Techie/st-okr-api/internal/database"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	// "github.com/go-playground/validator/v10"
)

func main() {
	db.InitDB()
	// validator := validator.New()

	engine := gin.Default()

	engine.Use(cors.Default())

	engine.GET("/", func(c *gin.Context) {
		c.String(200, "pong")
	})

	// v1:=engine.Group("/api/v1")

	if err := engine.Run(":" + config.ENV.ServerPort); err != nil {
		log.Panicf("error: %s", err)
	}

	log.Printf("server running on port: %s", config.ENV.ServerPort)
}
