package route

import (
	"go-payment-aggregator/internal/delivery/http/handler"
	"go-payment-aggregator/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
)

type RouteConfig struct {
	App                    *gin.Engine
	MerchantHandler        *handler.MerchantHandler
	TransactionHandler     *handler.TransactionHandler
	AuthMiddleware         *middleware.AuthMiddleware
	MidtransWebhookHandler *handler.MidtransWebhookHandler
}

func (c *RouteConfig) Setup() {
	c.SetupRoutes()
	c.App.Use(gin.Logger())
	c.App.Use(gin.Recovery())
}

func (c *RouteConfig) SetupRoutes() {
	v1 := c.App.Group("/api/v1")
	{
		m := v1.Group("/merchants")
		{
			m.POST("", c.MerchantHandler.Register)
			m.GET("/profile", c.AuthMiddleware.RequireApiKey(), c.MerchantHandler.Get)
			m.PUT("/profile", c.AuthMiddleware.RequireApiKey(), c.MerchantHandler.Update)
			m.POST("/api-key/regenerate", c.AuthMiddleware.RequireApiKey(), c.MerchantHandler.RegenerateApiKey)
		}

		t := v1.Group("/transactions")
		{
			t.POST("", c.AuthMiddleware.RequireApiKey(), c.TransactionHandler.Create)
			t.GET("/:id", c.AuthMiddleware.RequireApiKey(), c.TransactionHandler.Get)
		}

		w := v1.Group("/webhooks")
		{
			w.POST("/midtrans", c.MidtransWebhookHandler.Handle)
		}
	}
}
