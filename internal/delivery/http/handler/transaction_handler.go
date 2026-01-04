package handler

import (
	"go-payment-aggregator/internal/domain"
	"go-payment-aggregator/internal/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TransactionHandler struct {
	transactionUC domain.TransactionUC
}

func NewTransactionHandler(u domain.TransactionUC) *TransactionHandler {
	return &TransactionHandler{
		transactionUC: u,
	}
}

func (h *TransactionHandler) Create(c *gin.Context) {
	merchantData, exists := c.Get("merchant")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "unauthorized", "Merchant not found in context")
		return
	}

	merchant := merchantData.(*domain.Merchant)

	var req domain.CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "error", err.Error())
		return
	}

	ctx := c.Request.Context()
	createdTransaction, err := h.transactionUC.Create(ctx, merchant.ID, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "error", err.Error())
		return
	}

	data := response.CreateTransactionResponse{
		ID:            createdTransaction.ID.String(),
		MerchantID:    createdTransaction.MerchantID.String(),
		OrderID:       createdTransaction.OrderID,
		Provider:      createdTransaction.Provider,
		Currency:      createdTransaction.Currency,
		Amount:        createdTransaction.Amount,
		Status:        string(createdTransaction.Status),
		PaymentMethod: createdTransaction.PaymentMethod,
		PaymentURL:    createdTransaction.PaymentURL,
		ExternalID:    createdTransaction.ExternalID,
		ExpiredAt:     createdTransaction.ExpiredAt,
		CreatedAt:     createdTransaction.CreatedAt,
		UpdatedAt:     createdTransaction.UpdatedAt,
	}

	response.Success(c, http.StatusCreated, "success", "Transaction created successfully", data)
}

func (h *TransactionHandler) Get(c *gin.Context) {
	transactionID := c.Param("id")
	if transactionID == "" {
		response.Error(c, http.StatusBadRequest, "error", "Transaction ID is required")
		return
	}

	id, err := uuid.Parse(transactionID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "error", "Invalid transaction ID")
		return
	}

	ctx := c.Request.Context()
	transaction, err := h.transactionUC.Get(ctx, id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "error", "Transaction not found")
		return
	}

	data := response.CreateTransactionResponse{
		ID:            transaction.ID.String(),
		MerchantID:    transaction.MerchantID.String(),
		OrderID:       transaction.OrderID,
		Provider:      transaction.Provider,
		Currency:      transaction.Currency,
		Amount:        transaction.Amount,
		Status:        string(transaction.Status),
		PaymentMethod: transaction.PaymentMethod,
		PaymentURL:    transaction.PaymentURL,
		ExternalID:    transaction.ExternalID,
		ExpiredAt:     transaction.ExpiredAt,
		CreatedAt:     transaction.CreatedAt,
		UpdatedAt:     transaction.UpdatedAt,
	}

	response.Success(c, http.StatusOK, "success", "Transaction retrieved successfully", data)
}
