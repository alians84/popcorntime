package config

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var DB *gorm.DB

func ConnectPostgres() {
	dsn := buildDSN()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Логирование SQL
		NowFunc: func() time.Time {
			return time.Now().UTC() // UTC для всех временных меток
		},
		PrepareStmt: true, // Кэширование подготовленных запросов
	})

	if err != nil {
		log.Fatalf("Postgres connection error: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("DB instance error: %v", err)
	}

	// Настройки пула соединений
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	DB = db
	log.Println("Successfully connected to PostgreSQL")
}

func buildDSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", ""),
		getEnv("DB_NAME", "popcorntime"),
		getEnv("DB_PORT", "5432"),
	)
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
