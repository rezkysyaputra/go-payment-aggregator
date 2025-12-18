package config

import (
	"go-payment-aggregator/internal/domain/merchant"
	"go-payment-aggregator/internal/domain/transaction"
	"go-payment-aggregator/internal/gateway"
	"go-payment-aggregator/internal/gateway/midtrans"
	"go-payment-aggregator/internal/gateway/mock"
	"go-payment-aggregator/internal/gateway/xendit"
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
	merchantRepo := merchant.NewMerchantRepository(config.DB, config.Log)
	transactionRepo := transaction.NewTransactionRepository(config.DB, config.Log)

	// setup gateways
	midtransServerKey := config.Config.GetString("midtrans.server_key")
	midtransGateway := midtrans.NewMidtransGateway(midtransServerKey)

	mockGateway := mock.NewMockGateway()

	xenditApiKey := config.Config.GetString("xendit.api_key")
	xenditCallbackToken := config.Config.GetString("xendit.callback_token")
	xenditGateway := xendit.NewXenditGateway(xenditApiKey, xenditCallbackToken)

	gateways := map[string]gateway.PaymentGateway{
		"midtrans": midtransGateway,
		"mock":     mockGateway,
		"xendit":   xenditGateway,
	}
	// setup redis
	redisClient := NewRedis(config.Config, config.Log)

	// setup services
	merchantService := merchant.NewMerchantService(merchantRepo, config.Log)
	transactionService := transaction.NewTransactionService(transactionRepo, gateways, config.Config, config.Log, redisClient)

	// setup handlers
	merchantHandler := handler.NewMerchantHandler(merchantService, config.Log)
	transactionHandler := handler.NewTransactionHandler(transactionService, config.Log)
	webhookHandler := handler.NewWebhookHandler(transactionService, config.Log)

	// setup router
	router.SetupRouter(config.App, *merchantHandler, *transactionHandler, merchantRepo, *webhookHandler, config.Log)
}
