package transaction

import (
	"go-payment-aggregator/internal/domain/merchant"
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID          uuid.UUID          `gorm:"type:uuid;default:uuid_generate_v4();primary_key;" json:"id"`
	MerchantID  uuid.UUID          `gorm:"type:uuid;not null" json:"merchant_id"`
	Merchant    *merchant.Merchant `gorm:"foreignKey:MerchantID" json:"merchant,omitempty"`
	OrderID     string             `gorm:"type:size:100;not null" json:"order_id"`
	Provider    string             `gorm:"type:size:100;not null" json:"provider"`
	Amount      float64            `gorm:"type:numeric;not null" json:"amount"`
	Status      string             `gorm:"type:size:50;not null" json:"status"`
	ExternalRef string             `gorm:"type:size:100;not null" json:"external_ref"`
	RedirectURL string             `gorm:"type:size:255;not null" json:"redirect_url"`
	RawResponse string             `gorm:"type:jsonb;not null" json:"raw_response"`
	CreatedAt   time.Time          `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time          `gorm:"autoUpdateTime" json:"updated_at"`
}
