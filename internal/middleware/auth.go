package middleware

import (
	"context"
	"go-payment-aggregator/internal/domain/merchant"
	"net/http"

	"github.com/gin-gonic/gin"
)

func APIKeyAuth(merchantRepo merchant.MerchantRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("x-api-key")
		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized - API key missing"})
			c.Abort()
			return
		}

		m, err := merchantRepo.FindByApiKey(context.Background(), apiKey)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized - Invalid API key"})
			c.Abort()
			return
		}

		c.Set("merchant", m)
		c.Next()
	}
}
