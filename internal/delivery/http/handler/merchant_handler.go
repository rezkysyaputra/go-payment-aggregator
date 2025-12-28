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

func (h *MerchantHandler) Get(c *gin.Context) {
	// get merchant from context
	merchantData, exists := c.Get("merchant")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "unauthorized", "merchant not found in context")
		return
	}

	// type assert merchant data
	merchant := merchantData.(*domain.Merchant)

	// get fresh merchant profile
	ctx := c.Request.Context()
	freshMerchant, err := h.merchantUC.GetProfile(ctx, merchant.ID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "error", "failed to get profile")
		return
	}

	// prepare response data
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

	// send success response
	response.Success(c, http.StatusOK, "success", "merchant profile retrieved successfully", data)
}

func (h *MerchantHandler) Update(c *gin.Context) {
	// get merchant from context
	merchantData, exists := c.Get("merchant")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "unauthorized", "merchant not found in context")
		return
	}

	// type assert merchant data
	merchant := merchantData.(*domain.Merchant)

	// bind JSON request
	var req domain.UpdateMerchantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "error", err.Error())
		return
	}

	// call usecase to update profile
	ctx := c.Request.Context()
	updateMerchant, err := h.merchantUC.UpdateProfile(ctx, merchant.ID, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "error", "failed to update profile")
		return
	}

	// prepare response data
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

	// send success response
	response.Success(c, http.StatusOK, "success", "merchant profile updated successfully", data)
}
