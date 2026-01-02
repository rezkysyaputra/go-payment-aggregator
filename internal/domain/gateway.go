package domain

type PaymentGateway interface {
	CreatePayment(tx *CreatePaymentRequest) (*PaymentResponse, error)
	CheckStatus(orderID string) (string, error)
}

type PaymentResponse struct {
	Token      string `json:"token"`
	PaymentURL string `json:"payment_url"`
}

type CreatePaymentRequest struct {
	OrderID       string   `json:"order_id"`
	Amount        int64    `json:"amount"`
	PaymentMethod string   `json:"payment_method"`
	ExpiryMinutes int64    `json:"expiry_minutes"`
	Customer      Customer `json:"customer"`
	Items         []Item   `json:"items"`
}
