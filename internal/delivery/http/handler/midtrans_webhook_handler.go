package handler

import (
	"go-payment-aggregator/internal/domain"
	"go-payment-aggregator/internal/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MidtransWebhookHandler struct {
	transactionUC domain.TransactionUC
	ServerKey     string
}

func NewMidtransWebhookHandler(u domain.TransactionUC, serverKey string) *MidtransWebhookHandler {
	return &MidtransWebhookHandler{
		transactionUC: u,
		ServerKey:     serverKey,
	}
}

func (h *MidtransWebhookHandler) Handle(c *gin.Context) {
	var req MidtransWebhookRequest
	c.BindJSON(&req)

	isValidSignature := pkg.VerifySignature(req.OrderID, req.StatusCode, req.GrossAmount, h.ServerKey, req.SignatureKey)
	if !isValidSignature {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": "Invalid signature key"})
		return
	}

	domainReq := domain.UpdateStatusRequest{
		OrderID: req.OrderID,
		Status:  pkg.MapMidtransStatus(req.TransactionStatus, req.FraudStatus),
	}

	ctx := c.Request.Context()
	err := h.transactionUC.HandleNotification(ctx, &domainReq)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Notification processed"})
}

type MidtransWebhookRequest struct {
	OrderID           string `json:"order_id"`
	TransactionStatus string `json:"transaction_status"`
	FraudStatus       string `json:"fraud_status"`
	SignatureKey      string `json:"signature_key"`
	StatusCode        string `json:"status_code"`
	GrossAmount       string `json:"gross_amount"`
}
