package handler

import (
	"context"
	"go-payment-aggregator/internal/domain/merchant"
	"go-payment-aggregator/internal/helper"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type MerchantHandler struct {
	service merchant.MerchantService
	log     *logrus.Logger
}

func NewMerchantHandler(s merchant.MerchantService, log *logrus.Logger) *MerchantHandler {
	return &MerchantHandler{
		service: s,
		log:     log,
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
