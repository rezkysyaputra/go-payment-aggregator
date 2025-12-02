package handler

import (
	"go-payment-aggregator/internal/domain/merchant"
	"go-payment-aggregator/internal/domain/transaction"
	"go-payment-aggregator/internal/helper"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type TransactionHandler struct {
	service transaction.TransactionService
	log     *logrus.Logger
}

func NewTransactionHandler(s transaction.TransactionService, log *logrus.Logger) *TransactionHandler {
	return &TransactionHandler{
		service: s,
		log:     log,
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
		h.log.Printf("Error binding JSON: %v", err)
		helper.ErrorResponse(c, http.StatusBadRequest, false, err.Error())
		return
	}

	m := c.MustGet("merchant").(*merchant.Merchant)

	// create transaction
	t, err := h.service.CreateTransaction(m.ID, req.OrderID, req.Amount, req.Provider)
	if err != nil {
		h.log.Printf("Error creating transaction: %v", err)
		helper.ErrorResponse(c, http.StatusInternalServerError, false, "failed to create transaction")
		return
	}

	helper.SuccessResponse(c, http.StatusCreated, true, "transaction created", gin.H{
		"transaction_id": t.ID,
		"redirect_url":   t.RedirectURL,
		"status":         t.Status,
	})
}

func (h *TransactionHandler) GetById(c *gin.Context) {
	id := c.Param("id")

	if _, err := uuid.Parse(id); err != nil {
		h.log.Printf("Invalid transaction ID: %v", err)
		helper.ErrorResponse(c, http.StatusBadRequest, false, "invalid transaction ID")
		return
	}

	tx, err := h.service.FindById(uuid.Must(uuid.Parse(id)))
	if err != nil {
		h.log.Printf("Error finding transaction: %v", err)
		helper.ErrorResponse(c, http.StatusNotFound, false, "transaction not found")
		return
	}

	helper.SuccessResponse(c, http.StatusOK, true, "transaction found", gin.H{
		"transaction_id": tx.ID,
		"redirect_url":   tx.RedirectURL,
		"status":         tx.Status,
	})
}

func (h *TransactionHandler) GetByOrderId(c *gin.Context) {
	orderId := c.Param("order_id")

	tx, err := h.service.FindOrderById(orderId)
	if err != nil {
		h.log.Printf("Error finding transaction: %v", err)
		helper.ErrorResponse(c, http.StatusNotFound, false, "transaction not found")
		return
	}

	helper.SuccessResponse(c, http.StatusOK, true, "transaction found", gin.H{
		"transaction_id": tx.ID,
		"redirect_url":   tx.RedirectURL,
		"status":         tx.Status,
	})
}

func (h *TransactionHandler) UpdateStatusAndRaw(c *gin.Context) {
	id := c.Param("id")

	if _, err := uuid.Parse(id); err != nil {
		h.log.Printf("Invalid transaction ID: %v", err)
		helper.ErrorResponse(c, http.StatusBadRequest, false, "invalid transaction ID")
		return
	}

	var req struct {
		Status  string `json:"status" binding:"required"`
		RawJSON string `json:"raw_json" binding:"required"`
	}
	// bind JSON request
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Printf("Error binding JSON: %v", err)
		helper.ErrorResponse(c, http.StatusBadRequest, false, err.Error())
		return
	}

	// update transaction status and raw response
	if err := h.service.UpdateStatusAndRaw(uuid.Must(uuid.Parse(id)), req.Status, req.RawJSON); err != nil {
		h.log.Printf("Error updating transaction: %v", err)
		helper.ErrorResponse(c, http.StatusInternalServerError, false, "failed to update transaction")
		return
	}

	helper.SuccessResponse(c, http.StatusOK, true, "transaction updated", nil)
}
