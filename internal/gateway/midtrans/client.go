package midtrans

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-payment-aggregator/internal/gateway"
	"log"
	"net/http"
)

type SnapRequest struct {
	TransactionDetails TransactionDetails `json:"transaction_details"`
}

type TransactionDetails struct {
	OrderID     string  `json:"order_id"`
	GrossAmount float64 `json:"gross_amount"`
}

type SnapResponse struct {
	Token       string `json:"token"`
	RedirectURL string `json:"redirect_url"`
}

type MidtransGateway struct {
	ServerKey string
}

func NewMidtransGateway(serverKey string) *MidtransGateway {
	return &MidtransGateway{
		ServerKey: serverKey,
	}
}

func (g *MidtransGateway) CreateTransaction(orderID string, amount float64) (*gateway.PaymentResponse, error) {
	if g.ServerKey == "" {
		log.Println("Midtrans server key is not configured")
		return nil, fmt.Errorf("midtrans server key is not configured")
	}

	url := "https://app.sandbox.midtrans.com/snap/v1/transactions"
	body := SnapRequest{}
	body.TransactionDetails.OrderID = orderID
	body.TransactionDetails.GrossAmount = amount

	jsonData, err := json.Marshal(body)
	if err != nil {
		log.Printf("Error marshaling JSON: %v", err)
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error creating HTTP request: %v", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(g.ServerKey, "")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Printf("Error making HTTP request to Midtrans: %v", err)
		return nil, err
	}

	defer res.Body.Close()

	var snapRes SnapResponse
	if err := json.NewDecoder(res.Body).Decode(&snapRes); err != nil {
		log.Printf("Error decoding Midtrans response: %v", err)
		return nil, err
	}

	return &gateway.PaymentResponse{
		Token:       snapRes.Token,
		RedirectURL: snapRes.RedirectURL,
	}, nil
}

func (g *MidtransGateway) VerifySignature(payload map[string]any) bool {
	orderID, ok := payload["order_id"].(string)
	if !ok {
		return false
	}

	statusCode, _ := payload["status_code"].(string)

	grossStr, err := GrossAmountToString(payload["gross_amount"])
	if err != nil {
		return false
	}

	payloadSignature := ""
	if sig, exists := payload["signature_key"].(string); exists {
		payloadSignature = sig
	}

	return VerifySignature(payloadSignature, orderID, statusCode, grossStr, g.ServerKey)
}

func (g *MidtransGateway) GetOrderID(payload map[string]any) (string, error) {
	orderID, ok := payload["order_id"].(string)
	if !ok {
		return "", fmt.Errorf("order_id not found")
	}
	return orderID, nil
}

func (g *MidtransGateway) GetStatus(payload map[string]any) (string, error) {
	status, ok := payload["transaction_status"].(string)
	if !ok {
		return "", fmt.Errorf("transaction_status not found")
	}

	switch status {
	case "capture", "settlement":
		return "paid", nil
	case "deny", "expire", "cancel":
		return "failed", nil
	case "pending":
		return "pending", nil
	default:
		return "unknown", nil
	}
}
