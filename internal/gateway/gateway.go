package gateway

// PaymentResponse represents the response from the payment gateway
type PaymentResponse struct {
	Token       string
	RedirectURL string
	ExternalID  string
}

// PaymentGateway represents the payment gateway interface
type PaymentGateway interface {
	CreateTransaction(orderID string, amount float64) (*PaymentResponse, error)
	VerifySignature(payload map[string]any) bool
	GetStatus(payload map[string]any) (string, error)
	GetOrderID(payload map[string]any) (string, error)
}
