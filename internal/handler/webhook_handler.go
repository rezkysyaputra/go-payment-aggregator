package handler

import (
	"go-payment-aggregator/internal/domain/transaction"
	"go-payment-aggregator/internal/helper"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type WebhookHandler struct {
	txService transaction.TransactionService
	log       *logrus.Logger
}

func NewWebhookHandler(s transaction.TransactionService, log *logrus.Logger) *WebhookHandler {
	return &WebhookHandler{
		txService: s,
		log:       log,
	}
}

func (h *WebhookHandler) Midtrans(c *gin.Context) {
	var payload transaction.NotificationPayload
	// bind JSON payload
	if err := c.ShouldBindJSON(&payload); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, false, "invalid payload")
		return
	}

	// process notification
	if err := h.txService.ProcessMidtransNotification(payload); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, false, "failed to process notification")
		return
	}

	helper.SuccessResponse(c, http.StatusOK, true, "notification processed", nil)
}
