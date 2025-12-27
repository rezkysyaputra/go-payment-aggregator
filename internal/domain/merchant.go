package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// --- Enums ---
type MerchantStatus string

const (
	MerchantStatusActive    MerchantStatus = "ACTIVE"
	MerchantStatusSuspended MerchantStatus = "SUSPENDED"
	MerchantStatusInactive  MerchantStatus = "INACTIVE"
)

// --- Entity ---
type Merchant struct {
	ID          uuid.UUID      `json:"id"`
	Name        string         `json:"name"`
	Email       string         `json:"email"`
	ApiKey      string         `json:"api_key,omitempty"` // Plain text (hanya saat create)
	APIKeyHash  string         `json:"-"`                 // Hash (disimpan di DB)
	CallbackURL string         `json:"callback_url"`
	Status      MerchantStatus `json:"status"`
	Balance     int64          `json:"balance"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

// --- Repository Interface (Kontrak ke Database) ---
type MerchantRepository interface {
	Create(ctx context.Context, m *Merchant) error
	FindByApiKey(ctx context.Context, apiKey string) (*Merchant, error)
	FindByID(ctx context.Context, id uuid.UUID) (*Merchant, error)
}

// --- Usecase Interface (Kontrak ke Logic) ---
type MerchantUsecase interface {
	Register(ctx context.Context, req *RegisterMerchantRequest) (*Merchant, error)
}

// --- DTOs (Data Transfer Objects) ---
type RegisterMerchantRequest struct {
	Name        string
	Email       string
	CallbackURL string
}
