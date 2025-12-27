package route

import (
	"go-payment-aggregator/internal/delivery/http/handler"

	"github.com/gin-gonic/gin"
)

type RouteConfig struct {
	App             *gin.Engine
	MerchantHandler *handler.MerchantHandler
}

func (c *RouteConfig) Setup() {
	c.SetupRoutes()
	c.App.Use(gin.Logger())
	c.App.Use(gin.Recovery())
}

func (c *RouteConfig) SetupRoutes() {
	v1 := c.App.Group("/api/v1")
	{
		merchant := v1.Group("/merchants")
		{
			merchant.POST("", c.MerchantHandler.Register)
		}

	}
}
