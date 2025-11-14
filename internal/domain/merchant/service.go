package merchant

import (
	"context"
	"go-payment-aggregator/internal/utils"
	"time"

	"github.com/google/uuid"
)

type MerchantService interface {
	Create(ctx context.Context, name, callbackURL string) (*Merchant, error)
}

type MerchantServiceImpl struct {
	MerchantRepository MerchantRepository
}

func NewMerchantService(repo MerchantRepository) MerchantService {
	return &MerchantServiceImpl{
		MerchantRepository: repo,
	}
}

func (s *MerchantServiceImpl) Create(ctx context.Context, name, callbackURL string) (*Merchant, error) {
	// generate api key
	random, err := utils.RandomBase64(40)
	if err != nil {
		return nil, err
	}

	// prepend prefix
	apiKey := "mch_" + random

	// create merchant
	merchant := &Merchant{
		ID:          uuid.New(),
		Name:        name,
		ApiKey:      apiKey,
		CallbackURL: callbackURL,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.MerchantRepository.Create(ctx, merchant); err != nil {
		return nil, err
	}

	return merchant, nil
}
