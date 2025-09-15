package main

import (
	"net/http"
	"time"

	"github.com/Slightly-Techie/st-okr-api/config"
	"github.com/Slightly-Techie/st-okr-api/db"
	"github.com/Slightly-Techie/st-okr-api/internal/logger"
	"github.com/Slightly-Techie/st-okr-api/internal/message"
	"github.com/Slightly-Techie/st-okr-api/internal/routes"
	auth "github.com/Slightly-Techie/st-okr-api/pkg"
	"github.com/Slightly-Techie/st-okr-api/provider"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func main() {
	// Initialize logger
	logger.InitGlobal()
	defer logger.Custom.Close()

	logger.Info("Starting ST OKR API server")

	database, err := db.InitDB()
	if err != nil {
		logger.Fatal("Failed to connect to database", "error", err)
	}
	logger.Info("Database connection established")

	config := config.ENV

	// initialize rabitmq
	logger.Info("Initializing RabbitMQ connection")
	var connected bool
	for retries := 0; retries < 5; retries++ {
		if err := message.TestRabbitMQConnection(config); err != nil {
			logger.Warn("RabbitMQ connection attempt failed", "attempt", retries+1, "error", err)
			time.Sleep(10 * time.Second)
		} else {
			connected = true
			break
		}
	}

	if !connected {
		logger.Fatal("Failed to connect to RabbitMQ after 5 attempts")
	}

	logger.Info("RabbitMQ connection established")

	validator := validator.New()
	auth.NewAuth()

	go message.ConsumeMessages()

	provider := provider.NewProvider(database, validator)

	router := routes.SetupRouter(provider)
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello to the SlightlyTechie OKR API!"})
	})

	logger.Info("Server starting", "port", "8080")
	router.Run(":8080")
}
