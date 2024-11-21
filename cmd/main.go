package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Slightly-Techie/st-okr-api/config"
	"github.com/Slightly-Techie/st-okr-api/db"
	"github.com/Slightly-Techie/st-okr-api/internal/message"
	"github.com/Slightly-Techie/st-okr-api/internal/routes"
	auth "github.com/Slightly-Techie/st-okr-api/pkg"
	"github.com/Slightly-Techie/st-okr-api/provider"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func main() {
	database, err := db.InitDB()
	if err != nil {
		log.Fatalf("could not connect to db: %v", err)
	}

	config := config.ENV

	// initialize rabitmq
	var connected bool
	for retries := 0; retries <5;retries ++{
		if err := message.TestRabbitMQConnection(config); err != nil {
			log.Printf("Attempt %d: %v", retries, err)
			time.Sleep(10 * time.Second)
		} else {
			connected = true
			break
		}
	}

	if !connected {
		log.Fatalf("could not connect to RabbitMQ")
	}

	log.Println("Connected to RabbitMQ")


	validator := validator.New()
	auth.NewAuth()

	go message.ConsumeMessages()

	provider := provider.NewProvider(database, validator)

	router := routes.SetupRouter(provider)
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello to the SlightlyTechie OKR API!"})
	})

	router.Run(":8080")
}
