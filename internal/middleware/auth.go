package middleware

import (
	"context"
	"go-payment-aggregator/internal/domain/merchant"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func APIKeyAuth(merchantRepo merchant.MerchantRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("x-api-key")
		if apiKey == "" {
			log.Println("API key missing in request")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized - API key missing"})
			c.Abort()
			return
		}

		m, err := merchantRepo.FindByApiKey(context.Background(), apiKey)
		if err != nil {
			log.Printf("Error fetching merchant by API key: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized - Invalid API key"})
			c.Abort()
			return
		}

		c.Set("merchant", m)
		c.Next()
	}
}
