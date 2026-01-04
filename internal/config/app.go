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

func Bootstrap(b *BootstrapConfig) {
	var midtransEnv midtrans.EnvironmentType
	if b.Config.GetString("MIDTRANS_ENVIRONMENT") == "production" {
		midtransEnv = midtrans.Production
	} else {
		midtransEnv = midtrans.Sandbox
	}

	mtConfig := gateway.MidtransConfig{
		ServerKey: b.Config.GetString("MIDTRANS_SERVER_KEY"),
		Env:       midtransEnv,
	}

	midtransGateway := gateway.NewMidtransGateway(mtConfig)

	merchantRepository := postgres.NewMerchantRepository(b.DB)
	transactionRepository := postgres.NewTransactionRepository(b.DB)

	merchantUsecase := usecase.NewMerchantUC(merchantRepository, time.Second*2)
	transactionUsecase := usecase.NewTransactionUC(transactionRepository, midtransGateway, time.Second*time.Duration(b.Config.GetInt64("CONTEXT_TIMEOUT")))

	merchantHandler := handler.NewMerchantHandler(merchantUsecase)
	transactionHandler := handler.NewTransactionHandler(transactionUsecase)

	authMiddleware := middleware.NewAuthMiddleware(merchantUsecase)

	midtransWebhookHandler := handler.NewMidtransWebhookHandler(transactionUsecase, b.Config.GetString("MIDTRANS_SERVER_KEY"))

	routeConfig := &route.RouteConfig{
		App:                    b.App,
		MerchantHandler:        merchantHandler,
		TransactionHandler:     transactionHandler,
		AuthMiddleware:         authMiddleware,
		MidtransWebhookHandler: midtransWebhookHandler,
	}

	routeConfig.Setup()
}
