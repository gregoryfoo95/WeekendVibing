package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"fithero-backend/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDatabase initializes the database connection using GORM
func InitDatabase() {
	// Database configuration from environment variables
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "fithero_user")
	dbPassword := getEnv("DB_PASSWORD", "fithero_password")
	dbName := getEnv("DB_NAME", "fithero")

	// Create connection string
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		dbHost, dbUser, dbPassword, dbName, dbPort)

	var err error
	
	// Try to connect to the database with retries
	for i := 0; i < 30; i++ {
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err == nil {
			// Test the connection
			sqlDB, err := DB.DB()
			if err == nil {
				err = sqlDB.Ping()
				if err == nil {
					break
				}
			}
		}
		log.Println("Waiting for database to be ready...")
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Successfully connected to database with GORM")

	// Auto-migrate the schema
	err = AutoMigrate()
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Database migration completed")
}

// AutoMigrate migrates all models
func AutoMigrate() error {
	return DB.AutoMigrate(
		&models.User{},
		&models.Task{},
		&models.DailyTask{},
		&models.Achievement{},
		&models.UserAchievement{},
	)
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}

// getEnv gets environment variable with fallback
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
} 