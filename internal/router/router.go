package router

import (
	"go-payment-aggregator/internal/domain/merchant"
	"go-payment-aggregator/internal/handler"
	"go-payment-aggregator/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(merchantHandler handler.MerchantHandler, transactionHandler handler.TransactionHandler, merchantRepo merchant.MerchantRepository, webhookHandler handler.WebhookHandler) *gin.Engine {
	r := gin.Default()

	r.GET("/", func(ctx *gin.Context) {
		ctx.String(200, "Hello world")
	})

	v1 := r.Group("/v1")

	// Webhook group
	webhookGroup := v1.Group("/webhook")
	{
		webhookGroup.POST("/midtrans", webhookHandler.Midtrans)
	}

	// Merchant group
	merchantGroup := v1.Group("/merchant")
	{
		merchantGroup.POST("register", merchantHandler.Register)
	}

	// Transaction group
	transactionGroup := v1.Group("/transaction")
	transactionGroup.Use(middleware.APIKeyAuth(merchantRepo))
	{
		transactionGroup.POST("/", transactionHandler.Create)
		transactionGroup.GET("/:id", transactionHandler.GetById)
		transactionGroup.GET("/:order_id", transactionHandler.GetByOrderId)
		transactionGroup.PUT("/:id", transactionHandler.UpdateStatusAndRaw)
	}

	return r
}
