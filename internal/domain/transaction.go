package domain

import (
	"context"
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
	ID            uuid.UUID         `json:"id"`
	MerchantID    uuid.UUID         `json:"merchant_id"`
	OrderID       string            `json:"order_id"`
	ExternalID    string            `json:"external_id"`
	Provider      string            `json:"provider"`
	PaymentMethod string            `json:"payment_method"`
	Amount        int64             `json:"amount"`
	Currency      string            `json:"currency"`
	Status        TransactionStatus `json:"status"`
	PaymentURL    string            `json:"payment_url"`
	RawResponse   string            `json:"-"`
	ExpiredAt     time.Time         `json:"expired_at"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
}

type TransactionRepository interface {
	Create(ctx context.Context, tx *Transaction) (*Transaction, error)
	Update(ctx context.Context, tx *Transaction) (*Transaction, error)
	Get(ctx context.Context, id uuid.UUID) (*Transaction, error)
	FindByOrderID(ctx context.Context, orderID string) (*Transaction, error)
}

type TransactionUC interface {
	Create(ctx context.Context, merchantID uuid.UUID, req *CreateTransactionRequest) (*Transaction, error)
	Get(ctx context.Context, id uuid.UUID) (*Transaction, error)
	HandleNotification(ctx context.Context, req *UpdateStatusRequest) error
}

type Customer struct {
	Name  string `json:"name" validate:"required,min=3"`
	Email string `json:"email" validate:"required,email"`
}

type Item struct {
	Name     string `json:"name" validate:"required"`
	Quantity int32  `json:"quantity" validate:"required,min=1"`
	Price    int64  `json:"price" validate:"required,min=1"`
}

type CreateTransactionRequest struct {
	OrderID       string   `json:"order_id" validate:"required"`
	Amount        int64    `json:"amount" validate:"required,min=1"`
	Provider      string   `json:"provider" validate:"required,oneof=midtrans xendit stripe"`
	Currency      string   `json:"currency" validate:"required,len=3,uppercase"`
	PaymentMethod string   `json:"payment_method" validate:"required"`
	Customer      Customer `json:"customer" validate:"required"`
	Items         []Item   `json:"items" validate:"required,dive"`
}

type UpdateStatusRequest struct {
	OrderID string `json:"order_id" validate:"required"`
	Status  string `json:"status" validate:"required,oneof=PAID FAILED EXPIRED"`
}
