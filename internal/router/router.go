package router

import (
	"go-payment-aggregator/internal/domain/merchant"
	"go-payment-aggregator/internal/handler"
	"go-payment-aggregator/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(app *gin.Engine, merchantHandler handler.MerchantHandler, transactionHandler handler.TransactionHandler, merchantRepo merchant.MerchantRepository, webhookHandler handler.WebhookHandler) {
	app.GET("/", func(ctx *gin.Context) {
		ctx.String(200, "Hello world")
	})

	v1 := app.Group("/v1")

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
		transactionGroup.GET("/order/:order_id", transactionHandler.GetByOrderId)
		transactionGroup.PUT("/:id", transactionHandler.UpdateStatusAndRaw)
	}
}
