package config

import (
	"go-payment-aggregator/internal/delivery/http/handler"
	"go-payment-aggregator/internal/delivery/http/middleware"
	"go-payment-aggregator/internal/delivery/http/route"
	"go-payment-aggregator/internal/gateway"
	"go-payment-aggregator/internal/repository/postgres"
	"go-payment-aggregator/internal/usecase"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/midtrans/midtrans-go"
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

	// setup midtrans gateway
	var midtransEnv midtrans.EnvironmentType
	if config.Config.GetString("midtrans.environment") == "production" {
		midtransEnv = midtrans.Production
	} else {
		midtransEnv = midtrans.Sandbox
	}

	mtConfig := gateway.MidtransConfig{
		ServerKey: config.Config.GetString("midtrans.server_key"),
		Env:       midtransEnv,
	}

	midtransGateway := gateway.NewMidtransGateway(mtConfig)

	// setup repositories
	merchantRepository := postgres.NewMerchantRepository(config.DB)
	transactionRepository := postgres.NewTransactionRepository(config.DB)

	// setup usecases
	merchantUsecase := usecase.NewMerchantUC(merchantRepository, time.Second*2)                           // for now set hardcoded timeout to 2 seconds
	transactionUsecase := usecase.NewTransactionUC(transactionRepository, midtransGateway, time.Second*5) // for now set hardcoded timeout to 5 seconds

	// setup handlers
	merchantHandler := handler.NewMerchantHandler(merchantUsecase)
	transactionHandler := handler.NewTransactionHandler(transactionUsecase)

	// setup auth middleware
	authMiddleware := middleware.NewAuthMiddleware(merchantUsecase)

	// setup midtrans webhook handler
	midtransWebhookHandler := handler.NewMidtransWebhookHandler(transactionUsecase, config.Config.GetString("midtrans.server_key"))

	// router setup
	routeConfig := &route.RouteConfig{
		App:                    config.App,
		MerchantHandler:        merchantHandler,
		TransactionHandler:     transactionHandler,
		AuthMiddleware:         authMiddleware,
		MidtransWebhookHandler: midtransWebhookHandler,
	}

	routeConfig.Setup()
}
