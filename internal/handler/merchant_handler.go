package handler

import (
	"context"
	"go-payment-aggregator/internal/domain/merchant"
	"go-payment-aggregator/internal/helper"
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
		helper.ErrorResponse(c, http.StatusBadRequest, false, err.Error())
		return
	}

	m, err := h.service.Create(context.Background(), req.Name, req.CallbackURL)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, false, err.Error())
		return
	}

	helper.SuccessResponse(c, http.StatusOK, true, "merchant created", gin.H{
		"merchant_id": m.ID,
		"api_key":     m.ApiKey,
	})
}
