package entity

import (
	"time"

	"github.com/google/uuid"
)

type MerchantStatus string

const (
	MerchantStatusActive    MerchantStatus = "ACTIVE"
	MerchantStatusSuspended MerchantStatus = "SUSPENDED"
	MerchantStatusInactive  MerchantStatus = "INACTIVE"
)

type Merchant struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`

	ApiKey     string `json:"api_key"`
	APIKeyHash string `json:"-"`

	CallbackURL string         `json:"callback_url"`
	Status      MerchantStatus `json:"status"`

	Balance int64 `json:"balance"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
