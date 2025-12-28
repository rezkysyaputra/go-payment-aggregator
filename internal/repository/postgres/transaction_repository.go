package postgres

import (
	"time"

	"github.com/google/uuid"
)

type TransactionModel struct {
	ID            uuid.UUID `gorm:"type:uuid;primary_key"`
	MerchantID    uuid.UUID `gorm:"type:uuid;not null"`
	OrderID       string    `gorm:"size:255;not null;unique"`
	Provider      string    `gorm:"size:255"`
	PaymentMethod string    `gorm:"size:255"`
	Amount        int64     `gorm:"not null"`
	Currency      string    `gorm:"size:10;not null;default:'IDR'"`
	Status        string    `gorm:"size:50;not null"`
	ExternalRef   string    `gorm:"size:255;not null"`
	RedirectURL   string    `gorm:"size:255"`
	RawResponse   []byte    `gorm:"type:jsonb"`
	ExpiredAt     time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (TransactionModel) TableName() string {
	return "transactions"
}
