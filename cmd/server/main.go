package main

import (
	"go-payment-aggregator/internal/db"
	"go-payment-aggregator/internal/domain/merchant"
	"go-payment-aggregator/internal/handler"
	"go-payment-aggregator/internal/router"
	"log"

	"github.com/spf13/viper"
)

func main() {

	// connect DB
	db, err := db.InitDB()
	if err != nil {
		log.Fatalf("failed to connect ")
	}

	// dependency injection
	repo := merchant.NewMerchantRepository(db)
	service := merchant.NewMerchantService(repo)
	merchantHandler := handler.NewMerchantHandler(service)

	// setup router
	r := router.SetupRouter(*merchantHandler)

	port := viper.GetString("DATABASE_PORT")
	if port == "" {
		port = "8080"
	}

	err = r.Run(":" + port)
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}

	log.Println("Server running on port " + port)
}
