package handler

import (
	"context"
	"go-payment-aggregator/internal/domain/merchant"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MerchantHandler struct {
	service merchant.MerchantService
}

func NewMerchantHandler(s merchant.MerchantService) *MerchantHandler {
	return &MerchantHandler{
		service: s,
	}
}

type registerMerchantRequest struct {
	Name        string `json:"name" binding:"required"`
	CallbackURL string `json:"callback_url" binding:"required,url"`
}

func (h *MerchantHandler) Register(c *gin.Context) {
	var req registerMerchantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	m, err := h.service.Create(context.Background(), req.Name, req.CallbackURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create merchant"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"merchant_id": m.ID,
		"api_key":     m.ApiKey,
	})
}
