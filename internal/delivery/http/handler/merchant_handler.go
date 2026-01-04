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

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "error", err.Error())
		return
	}

	ctx := c.Request.Context()
	merchant, err := h.merchantUC.Register(ctx, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "error", err.Error())
		return
	}

	data := response.RegisterMerchantResponse{
		ID:          merchant.ID.String(),
		Name:        merchant.Name,
		Email:       merchant.Email,
		Status:      string(merchant.Status),
		ApiKey:      merchant.ApiKey,
		CallbackURL: merchant.CallbackURL,
	}

	response.Success(c, http.StatusCreated, "success", "Merchant created successfully", data)
}

func (h *MerchantHandler) Get(c *gin.Context) {
	merchantData, exists := c.Get("merchant")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "unauthorized", "Merchant not found in context")
		return
	}

	merchant := merchantData.(*domain.Merchant)

	ctx := c.Request.Context()
	freshMerchant, err := h.merchantUC.GetProfile(ctx, merchant.ID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "error", "Failed to get profile")
		return
	}

	data := response.GetMerchantResponse{
		ID:          freshMerchant.ID.String(),
		Name:        freshMerchant.Name,
		Email:       freshMerchant.Email,
		Status:      string(freshMerchant.Status),
		Balance:     freshMerchant.Balance,
		CallbackURL: freshMerchant.CallbackURL,
		CreatedAt:   freshMerchant.CreatedAt,
		UpdatedAt:   freshMerchant.UpdatedAt,
	}

	response.Success(c, http.StatusOK, "success", "Merchant profile retrieved successfully", data)
}

func (h *MerchantHandler) Update(c *gin.Context) {
	merchantData, exists := c.Get("merchant")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "unauthorized", "Merchant not found in context")
		return
	}

	merchant := merchantData.(*domain.Merchant)

	var req domain.UpdateMerchantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "error", err.Error())
		return
	}

	ctx := c.Request.Context()
	updateMerchant, err := h.merchantUC.UpdateProfile(ctx, merchant.ID, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "error", "Failed to update profile")
		return
	}

	data := response.GetMerchantResponse{
		ID:          updateMerchant.ID.String(),
		Name:        updateMerchant.Name,
		Email:       updateMerchant.Email,
		Status:      string(updateMerchant.Status),
		Balance:     updateMerchant.Balance,
		CallbackURL: updateMerchant.CallbackURL,
		CreatedAt:   updateMerchant.CreatedAt,
		UpdatedAt:   updateMerchant.UpdatedAt,
	}

	response.Success(c, http.StatusOK, "success", "Merchant profile updated successfully", data)
}

func (h *MerchantHandler) RegenerateApiKey(c *gin.Context) {
	merchantData, exists := c.Get("merchant")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "unauthorized", "Merchant not found in context")
		return
	}

	merchant := merchantData.(*domain.Merchant)

	ctx := c.Request.Context()
	newApiKey, err := h.merchantUC.RegenerateApiKey(ctx, merchant.ID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "error", "Failed to regenerate API key")
		return
	}

	response.Success(c, http.StatusOK, "success", "API key regenerated successfully", &response.GenerateApiKeyResponse{
		ApiKey: newApiKey,
	})
}
