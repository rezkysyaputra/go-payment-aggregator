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
	if err := h.txService.ProcessNotification("midtrans", payload); err != nil {
		h.log.Errorf("Error processing Midtrans notification: %v", err)
		helper.ErrorResponse(c, http.StatusBadRequest, false, "failed to process notification")
		return
	}

	helper.SuccessResponse(c, http.StatusOK, true, "notification processed", nil)
}

func (h *WebhookHandler) Mock(c *gin.Context) {
	var payload transaction.NotificationPayload
	// bind JSON payload
	if err := c.ShouldBindJSON(&payload); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, false, "invalid payload")
		return
	}

	// process notification
	if err := h.txService.ProcessNotification("mock", payload); err != nil {
		h.log.Errorf("Error processing Mock notification: %v", err)
		helper.ErrorResponse(c, http.StatusBadRequest, false, "failed to process notification")
		return
	}

	helper.SuccessResponse(c, http.StatusOK, true, "notification processed", nil)
}

func (h *WebhookHandler) Xendit(c *gin.Context) {
	var payload map[string]any
	// bind JSON payload
	if err := c.ShouldBindJSON(&payload); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, false, "invalid payload")
		return
	}

	// Xendit: Inject X-CALLBACK-TOKEN from header to payload for verification
	token := c.GetHeader("x-callback-token")
	payload["x-callback-token"] = token

	// process notification
	if err := h.txService.ProcessNotification("xendit", payload); err != nil {
		h.log.Errorf("Error processing Xendit notification: %v", err)
		helper.ErrorResponse(c, http.StatusBadRequest, false, "failed to process notification")
		return
	}

	helper.SuccessResponse(c, http.StatusOK, true, "notification processed", nil)
}
