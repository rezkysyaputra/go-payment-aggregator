package handler

import (
	"go-payment-aggregator/internal/domain"

	"github.com/gin-gonic/gin"
)

type MerchantHandler struct {
	merchantUsecase domain.MerchantUsecase
}

func NewMerchantHandler(usecase domain.MerchantUsecase) *MerchantHandler {
	return &MerchantHandler{
		merchantUsecase: usecase,
	}
}

func (h *MerchantHandler) Register(c *gin.Context) {
	var req domain.RegisterMerchantRequest

	// bind JSON request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"code":    400,
			"status":  "error",
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	// call usecase to register merchant
	ctx := c.Request.Context()
	merchant, err := h.merchantUsecase.Register(ctx, &req)
	if err != nil {
		c.JSON(500, gin.H{
			"code":    500,
			"status":  "error",
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	// respond with created merchant details
	c.JSON(201, gin.H{
		"code":    201,
		"status":  "success",
		"message": "merchant created successfully",
		"data": gin.H{
			"id":           merchant.ID,
			"name":         merchant.Name,
			"email":        merchant.Email,
			"status":       merchant.Status,
			"api_key":      merchant.ApiKey,
			"callback_url": merchant.CallbackURL,
		},
	})
}
