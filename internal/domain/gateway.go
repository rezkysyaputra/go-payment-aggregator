package domain

type PaymentGateway interface {
	CreatePayment(req *CreatePaymentRequest) (*PaymentResponse, error)
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
	Currency      string   `json:"currency"`
	ExpiryMinutes int32    `json:"expiry_minutes"`
	Customer      Customer `json:"customer"`
	Items         []Item   `json:"items"`
}
