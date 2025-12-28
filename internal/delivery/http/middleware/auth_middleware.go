package middleware

import (
	"go-payment-aggregator/internal/domain"
	"go-payment-aggregator/internal/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	merchantUC domain.MerchantUC
}

func NewAuthMiddleware(usecase domain.MerchantUC) *AuthMiddleware {
	return &AuthMiddleware{
		merchantUC: usecase,
	}
}

func (m *AuthMiddleware) RequireApiKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-KEY")

		if apiKey == "" {
			response.Error(c, http.StatusUnauthorized, "unauthorized", "missing api key")
			c.Abort()
			return
		}

		ctx := c.Request.Context()
		merchant, err := m.merchantUC.ValidateApiKey(ctx, apiKey)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "unauthorized", "invalid api key")
			c.Abort()
			return
		}

		c.Set("merchant", merchant)
		c.Next()
	}
}
