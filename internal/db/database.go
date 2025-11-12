package db

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	viper := viper.New()
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	host := viper.GetString("DATABASE_HOST")
	user := viper.GetString("DATABASE_USER")
	password := viper.GetString("DATABASE_PASSWORD")
	dbname := viper.GetString("DATABASE_NAME")
	port := viper.GetString("DATABASE_PORT")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC", host, user, password, dbname, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("DB connection failed: %v", err)
		return nil, err
	}
	return db, nil
}

// migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/payment_aggregator?sslmode=disable" up
