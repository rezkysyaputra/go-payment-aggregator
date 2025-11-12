package router

import (
	"go-payment-aggregator/internal/handler"

	"github.com/gin-gonic/gin"
)

func SetupRouter(merchantHandler handler.MerchantHandler) *gin.Engine {
	r := gin.Default()

	r.GET("/", func(ctx *gin.Context) {
		ctx.String(200, "Hello world")
	})

	v1 := r.Group("/v1")

	merchantGroup := v1.Group("/merchant")
	{
		merchantGroup.POST("register", merchantHandler.Register)
	}

	return r
}
