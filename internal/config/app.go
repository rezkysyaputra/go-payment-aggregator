package config

import (
	"go-payment-aggregator/internal/delivery/http/handler"
	"go-payment-aggregator/internal/delivery/http/route"
	"go-payment-aggregator/internal/repository/postgres"
	"go-payment-aggregator/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
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
	Redis    *redis.Client
}

func Bootstrap(config *BootstrapConfig) {
	// setup repositories
	// merchantRepo := merchant.NewMerchantRepository(config.DB, config.Log)
	// transactionRepo := transaction.NewTransactionRepository(config.DB, config.Log)

	// // setup gateways
	// midtransServerKey := config.Config.GetString("midtrans.server_key")
	// midtransGateway := midtrans.NewMidtransGateway(midtransServerKey)

	// mockGateway := mock.NewMockGateway()

	// xenditApiKey := config.Config.GetString("xendit.api_key")
	// xenditCallbackToken := config.Config.GetString("xendit.callback_token")
	// xenditGateway := xendit.NewXenditGateway(xenditApiKey, xenditCallbackToken)

	// gateways := map[string]gateway.PaymentGateway{
	// 	"midtrans": midtransGateway,
	// 	"mock":     mockGateway,
	// 	"xendit":   xenditGateway,
	// }
	// // setup redis
	// redisClient := NewRedis(config.Config, config.Log)

	// // setup services
	// merchantService := merchant.NewMerchantService(merchantRepo, config.Log)
	// transactionService := transaction.NewTransactionService(transactionRepo, gateways, config.Config, config.Log, redisClient)

	// // setup handlers
	// merchantHandler := handler.NewMerchantHandler(merchantService, config.Log)
	// transactionHandler := handler.NewTransactionHandler(transactionService, config.Log)
	// webhookHandler := handler.NewWebhookHandler(transactionService, config.Log)

	// setup router
	// router.SetupRouter(config.App, *merchantHandler, *transactionHandler, merchantRepo, *webhookHandler, config.Log)

	// setup repositories
	merchantRepository := postgres.NewMerchantRepository(config.DB)
	// setup usecases
	merchantUsecase := usecase.NewMerchantUC(merchantRepository, config.Config.GetDuration("context_timeout"))
	// setup handlers
	merchantHandler := handler.NewMerchantHandler(merchantUsecase)

	// router setup
	routeConfig := &route.RouteConfig{
		App:             config.App,
		MerchantHandler: merchantHandler,
	}

	routeConfig.Setup()
}
