package xendit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-payment-aggregator/internal/gateway"
	"net/http"
)

type XenditGateway struct {
	ApiKey        string
	CallbackToken string
}

func NewXenditGateway(apiKey, callbackToken string) *XenditGateway {
	return &XenditGateway{
		ApiKey:        apiKey,
		CallbackToken: callbackToken,
	}
}

type InvoiceRequest struct {
	ExternalID  string  `json:"external_id"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
}

type InvoiceResponse struct {
	ID         string `json:"id"`
	InvoiceURL string `json:"invoice_url"`
}

func (g *XenditGateway) CreateTransaction(orderID string, amount float64) (*gateway.PaymentResponse, error) {
	url := "https://api.xendit.co/invoices"

	// Create request body
	reqBody := InvoiceRequest{
		ExternalID:  orderID,
		Amount:      amount,
		Description: fmt.Sprintf("Order %s", orderID),
	}

	// Convert request body to JSON
	jsonData, _ := json.Marshal(reqBody)

	// Create HTTP request
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.SetBasicAuth(g.ApiKey, "")
	req.Header.Set("Content-Type", "application/json")

	// Send HTTP request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	// Check response status
	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("xendit error status: %s", res.Status)
	}

	// Parse response
	var xenditRes InvoiceResponse
	if err := json.NewDecoder(res.Body).Decode(&xenditRes); err != nil {
		return nil, err
	}

	return &gateway.PaymentResponse{
		Token:       xenditRes.ID,
		RedirectURL: xenditRes.InvoiceURL,
		ExternalID:  xenditRes.ID,
	}, nil
}

func (g *XenditGateway) VerifySignature(payload map[string]any) bool {
	// implement signature verification
	token, ok := payload["x-callback-token"].(string)
	if !ok {
		return false
	}

	return token == g.CallbackToken
}

func (g *XenditGateway) GetOrderID(payload map[string]any) (string, error) {
	// implement get order id
	if oid, ok := payload["external_id"].(string); ok {
		return oid, nil
	}
	// implement error handling
	return "", fmt.Errorf("external_id not found")
}

func (g *XenditGateway) GetStatus(payload map[string]any) (string, error) {

	status, ok := payload["status"].(string)
	if !ok {
		return "unknwon", nil
	}

	switch status {
	case "PAID", "SETTLED":
		return "paid", nil
	case "EXPIRED":
		return "failed", nil
	default:
		return "pending", nil
	}

}
