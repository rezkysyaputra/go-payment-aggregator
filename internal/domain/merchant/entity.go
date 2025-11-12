package merchant

import (
	"time"

	"github.com/google/uuid"
)

type Merchant struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	Name        string    `gorm:"size:255;not null" json:"name"`
	ApiKey      string    `gorm:"size:255;not null;unique" json:"api_key"`
	CallbackURL string    `gorm:"size:255;not null" json:"callback_url"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
