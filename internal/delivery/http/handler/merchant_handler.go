package handler

import (
	"go-payment-aggregator/internal/domain"
	"go-payment-aggregator/internal/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MerchantHandler struct {
	merchantUC domain.MerchantUC
}

func NewMerchantHandler(usecase domain.MerchantUC) *MerchantHandler {
	return &MerchantHandler{
		merchantUC: usecase,
	}
}

func (h *MerchantHandler) Register(c *gin.Context) {
	var req domain.RegisterMerchantRequest

	// bind JSON request
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "error", err.Error())
		return
	}

	// call usecase to register merchant
	ctx := c.Request.Context()
	merchant, err := h.merchantUC.Register(ctx, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "error", err.Error())
		return
	}

	// prepare response data
	data := response.RegisterMerchantResponse{
		ID:          merchant.ID.String(),
		Name:        merchant.Name,
		Email:       merchant.Email,
		Status:      string(merchant.Status),
		ApiKey:      merchant.ApiKey,
		CallbackURL: merchant.CallbackURL,
	}

	// send success response
	response.Success(c, http.StatusCreated, "success", "merchant created successfully", data)
}
