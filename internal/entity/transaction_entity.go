package entity

import (
	"time"

	"github.com/google/uuid"
)

type TransactionStatus string

const (
	TransactionStatusPending TransactionStatus = "PENDING"
	TransactionStatusPaid    TransactionStatus = "PAID"
	TransactionStatusFailed  TransactionStatus = "FAILED"
	TransactionStatusExpired TransactionStatus = "EXPIRED"
)

type Transaction struct {
	ID         uuid.UUID `json:"id"`
	MerchantID uuid.UUID `json:"merchant_id"`

	OrderID    string `json:"order_id"`
	ExternalID string `json:"external_id"`

	Provider      string `json:"provider"`
	PaymentMethod string `json:"payment_method"`

	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`

	Status     TransactionStatus `json:"status"`
	PaymentURL string            `json:"payment_url"`

	RawResponse string `json:"-"`

	ExpiredAt time.Time `json:"expired_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
