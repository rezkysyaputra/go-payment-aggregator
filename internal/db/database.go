package db

import (
	"fmt"
	"go-payment-aggregator/internal/config"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC", cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("DB connection failed: %v", err)
		return nil, err
	}
	return db, nil
}

// migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/payment_aggregator?sslmode=disable" up
// migrate create -ext sql -dir migrations -seq create_users_table
