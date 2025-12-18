package mock

import (
	"fmt"
	"go-payment-aggregator/internal/gateway"

	"github.com/google/uuid"
)

type MockGateway struct{}

func NewMockGateway() *MockGateway {
	return &MockGateway{}
}

func (g *MockGateway) CreateTransaction(orderID string, amount float64) (*gateway.PaymentResponse, error) {
	id := uuid.New().String()
	return &gateway.PaymentResponse{
		Token:       "mock-token-" + id,
		RedirectURL: "mock-redirect-url-" + id,
		ExternalID:  id,
	}, nil
}

func (g *MockGateway) VerifySignature(payload map[string]any) bool {
	return true
}

func (g *MockGateway) GetOrderID(payload map[string]any) (string, error) {
	// Di mock kita anggap fieldnya "order_id" juga
	orderID, ok := payload["order_id"].(string)
	if !ok {
		return "", fmt.Errorf("order_id not found")
	}
	return orderID, nil
}

func (g *MockGateway) GetStatus(payload map[string]any) (string, error) {
	// Simplifikasi: Kalau mock, assume selalu statusnya 'paid'
	// Atau bisa juga baca dari payload kalau mau testing status failed
	if s, ok := payload["mock_status"].(string); ok {
		return s, nil
	}
	return "paid", nil
}
