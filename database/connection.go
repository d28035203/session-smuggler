package database

import (
	"fmt"
	"os"

	"github.com/d28035203/legendary-succotash/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Database wraps the GORM connection used by handlers.
type Database struct {
	DB *gorm.DB
}

// Connect opens a PostgreSQL connection and runs AutoMigrate.
func Connect() (*Database, error) {
	dsn := getDSN()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, fmt.Errorf("connection to database failed: %w", err)
	}

	if err := db.AutoMigrate(&models.User{}, &models.UserSessions{}); err != nil {
		return nil, fmt.Errorf("AutoMigrate failed: %w", err)
	}

	return &Database{DB: db}, nil
}

func getDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DBNAME"),
	)
}
