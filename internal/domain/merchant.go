package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Enums for MerchantStatus
type MerchantStatus string

const (
	MerchantStatusActive    MerchantStatus = "ACTIVE"
	MerchantStatusSuspended MerchantStatus = "SUSPENDED"
	MerchantStatusInactive  MerchantStatus = "INACTIVE"
)

// Entities for Merchant
type Merchant struct {
	ID          uuid.UUID      `json:"id"`
	Name        string         `json:"name"`
	Email       string         `json:"email"`
	ApiKey      string         `json:"api_key,omitempty"`
	APIKeyHash  string         `json:"-"`
	CallbackURL string         `json:"callback_url"`
	Status      MerchantStatus `json:"status"`
	Balance     int64          `json:"balance"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

// Interfaces for Repository
type MerchantRepository interface {
	Create(ctx context.Context, m *Merchant) error
	Update(ctx context.Context, m *Merchant) error
	FindByApiKey(ctx context.Context, apiKey string) (*Merchant, error)
	FindByID(ctx context.Context, id uuid.UUID) (*Merchant, error)
}

// Interfaces for Usecase
type MerchantUC interface {
	Register(ctx context.Context, req *RegisterMerchantRequest) (*Merchant, error)
	GetProfile(ctx context.Context, id uuid.UUID) (*Merchant, error)
	UpdateProfile(ctx context.Context, id uuid.UUID, req *UpdateMerchantRequest) (*Merchant, error)
	ValidateApiKey(ctx context.Context, apiKey string) (*Merchant, error)
}

// Request structs
type RegisterMerchantRequest struct {
	Name        string
	Email       string
	CallbackURL string
}

// UpdateMerchantRequest struct
type UpdateMerchantRequest struct {
	Name        string
	CallbackURL string
}
