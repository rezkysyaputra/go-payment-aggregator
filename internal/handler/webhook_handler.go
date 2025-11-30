package handler

import (
	"go-payment-aggregator/internal/domain/transaction"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WebhookHandler struct {
	txService transaction.TransactionService
}

func NewWebhookHandler(s transaction.TransactionService) *WebhookHandler {
	return &WebhookHandler{
		txService: s,
	}
}

func (h *WebhookHandler) Midtrans(c *gin.Context) {
	var payload transaction.NotificationPayload
	// bind JSON payload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	// process notification
	if err := h.txService.ProcessMidtransNotification(payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to process notification"})
		return
	}

	c.JSON(http.StatusOK, "ok")
}
