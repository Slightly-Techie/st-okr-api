package db

import (
	"fmt"
	"log"

	"github.com/Slightly-Techie/st-okr-api/config"
	"github.com/Slightly-Techie/st-okr-api/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func InitDB() (*gorm.DB, error) {
	// Create the DB Connection String from the config
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", config.ENV.DBUser, config.ENV.DBPassword, config.ENV.DBHost, config.ENV.DBPort, config.ENV.DBName)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: false,
		},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Run DB Migrations
	err = db.AutoMigrate(&models.User{}, &models.Company{})
	if err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Print("Connected to database and migrations applied")
	return db, nil
}
