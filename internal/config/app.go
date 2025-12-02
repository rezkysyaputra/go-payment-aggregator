package config

import (
	"go-payment-aggregator/internal/domain/merchant"
	"go-payment-aggregator/internal/domain/transaction"
	"go-payment-aggregator/internal/handler"
	"go-payment-aggregator/internal/router"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB       *gorm.DB
	App      *gin.Engine
	Log      *logrus.Logger
	Config   *viper.Viper
	Validate *validator.Validate
}

func Bootstrap(config *BootstrapConfig) {
	// setup repositories
	merchantRepo := merchant.NewMerchantRepository(config.DB)
	transactionRepo := transaction.NewTransactionRepository(config.DB)

	// setup services
	merchantService := merchant.NewMerchantService(merchantRepo)
	transactionService := transaction.NewTransactionService(transactionRepo, config.Config)

	// setup handlers
	merchantHandler := handler.NewMerchantHandler(merchantService)
	transactionHandler := handler.NewTransactionHandler(transactionService)
	webhookHandler := handler.NewWebhookHandler(transactionService)

	// setup router
	router.SetupRouter(config.App, *merchantHandler, *transactionHandler, merchantRepo, *webhookHandler)
}
