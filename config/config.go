package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort         string
	DBHost             string
	DBPort             string
	DBUser             string
	DBPassword         string
	DBName             string
	JWTKey             string
	SMTPUsername       string
	SMTPPassword       string
	SMTPHost           string
	SMTPAddress        string
	SMTPPort           string
	GoogleClientID     string
	GoogleClientSecret string
	RabbitUser         string
	RabbitPassword     string
	RabbitHost         string
	RabbitPort         string
}

var ENV = initConfig()

func initConfig() Config {

	err := godotenv.Load()
	if err != nil {
		log.Printf("unable to load .env")
	}

	return Config{
		ServerPort:         getEnv("PORT", "8080"),
		DBPort:             getEnv("DB_PORT", "5432"),
		DBHost:             getEnv("DB_HOST", "localhost"),
		DBUser:             getEnv("DB_USER", "alexanderdomakyaareh"),
		DBPassword:         getEnv("DB_PASSWORD", "mypassword"),
		DBName:             getEnv("DB_NAME", "postgres"),
		JWTKey:             getEnv("JWT_KEY", "someJWTKey"),
		SMTPUsername:       getEnv("SMTP_USERNAME", "someEmail"),
		SMTPPassword:       getEnv("SMTP_PASSWORD", "somePassword"),
		SMTPHost:           getEnv("SMTP_HOST", "smtp.emailprovider.com"),
		SMTPAddress:        getEnv("SMTP_ADDR", "smtp.gmail.com:587"),
		SMTPPort:           getEnv("SMTP_PORT", "587"),
		GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", "some-client-id"),
		GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", "some-client-secret"),
		RabbitUser:         getEnv("RABBIT_USER", "guest"),
		RabbitPassword:     getEnv("RABBIT_PASSWORD", "guest"),
		RabbitHost:         getEnv("RABBIT_HOST", "localhost"),
		RabbitPort:         getEnv("RABBIT_PORT", "5672"),
	}
}

func getEnv(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
