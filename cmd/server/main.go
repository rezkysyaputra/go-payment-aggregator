package main

import (
	"fmt"
	"go-payment-aggregator/internal/config"
	"go-payment-aggregator/internal/db"
	"go-payment-aggregator/internal/domain/merchant"
	"go-payment-aggregator/internal/domain/transaction"
	"go-payment-aggregator/internal/handler"
	"go-payment-aggregator/internal/router"
	"log"
)

func main() {

	// load config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	fmt.Println("✅ Config loaded:", cfg.AppEnv)

	// connect DB
	db, err := db.InitDB(cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// dependency injection
	repo := merchant.NewMerchantRepository(db)
	service := merchant.NewMerchantService(repo)
	merchantHandler := handler.NewMerchantHandler(service)

	transactionRepo := transaction.NewTransactionRepository(db)
	transactionService := transaction.NewTransactionService(transactionRepo, cfg.MidtransServerKey)
	transactionHandler := handler.NewTransactionHandler(transactionService)

	// webhook handler
	webhookHandler := handler.NewWebhookHandler(transactionService)

	// setup router
	r := router.SetupRouter(*merchantHandler, *transactionHandler, repo, *webhookHandler)

	err = r.Run(":" + cfg.Port)
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}

	log.Println("🚀 Server running on port " + cfg.Port)
}
