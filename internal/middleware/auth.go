package middleware

import (
	"context"
	"go-payment-aggregator/internal/domain/merchant"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func APIKeyAuth(merchantRepo merchant.MerchantRepository, log *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("x-api-key")
		if apiKey == "" {
			log.Error("API key missing in request")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized - API key missing"})
			c.Abort()
			return
		}

		ctx := context.Background()

		m, err := merchantRepo.FindByApiKey(ctx, apiKey)
		if err != nil {
			log.Errorf("Error fetching merchant by API key: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized - Invalid API key"})
			c.Abort()
			return
		}

		c.Set("merchant", m)
		c.Next()
	}
}
