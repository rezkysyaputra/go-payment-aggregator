package handler

import (
	"go-payment-aggregator/internal/domain/merchant"
	"go-payment-aggregator/internal/domain/transaction"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	service transaction.TransactionService
}

func NewTransactionHandler(s transaction.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		service: s,
	}
}

type createTransactionRequest struct {
	OrderID  string  `json:"order_id" binding:"required"`
	Amount   float64 `json:"amount" binding:"required,gt=0"`
	Provider string  `json:"provider" binding:"required,oneof=midtrans"`
}

func (h *TransactionHandler) Create(c *gin.Context) {
	var req createTransactionRequest
	// bind JSON request
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	m := c.MustGet("merchant").(*merchant.Merchant)

	// create transaction
	t, err := h.service.CreateTransaction(m.ID, req.OrderID, req.Amount, req.Provider)
	if err != nil {
		log.Printf("Error creating transaction: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create transaction"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"transaction_id": t.ID,
		"redirect_url":   t.RedirectURL,
		"status":         t.Status,
	})
}
