package merchant

import (
	"context"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type MerchantRepository interface {
	Create(ctx context.Context, merchant *Merchant) error
	FindByApiKey(ctx context.Context, apiKey string) (*Merchant, error)
}

type MerchantRepositoryImpl struct {
	db  *gorm.DB
	log *logrus.Logger
}

func NewMerchantRepository(db *gorm.DB, log *logrus.Logger) MerchantRepository {
	return &MerchantRepositoryImpl{
		db:  db,
		log: log,
	}
}

func (r MerchantRepositoryImpl) Create(ctx context.Context, merchant *Merchant) error {
	return r.db.WithContext(ctx).Create(merchant).Error
}

func (r MerchantRepositoryImpl) FindByApiKey(ctx context.Context, apiKey string) (*Merchant, error) {
	var merchant Merchant
	if err := r.db.WithContext(ctx).Where("api_key = ?", apiKey).First(&merchant).Error; err != nil {
		return nil, err
	}
	return &merchant, nil
}
