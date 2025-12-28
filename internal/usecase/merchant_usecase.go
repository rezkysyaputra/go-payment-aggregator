package usecase

import (
	"context"
	"go-payment-aggregator/internal/domain"
	"go-payment-aggregator/internal/pkg"
	"time"
)

type merchantUC struct {
	merchantRepo domain.MerchantRepository
	timeout      time.Duration
}

func NewMerchantUC(r domain.MerchantRepository, t time.Duration) domain.MerchantUC {
	return &merchantUC{
		merchantRepo: r,
		timeout:      t,
	}
}

func (u *merchantUC) ValidateApiKey(ctx context.Context, apiKey string) (*domain.Merchant, error) {
	// hash api key
	apiKeyHash := pkg.HashKey(apiKey)

	// find merchant by api key
	merchant, err := u.merchantRepo.FindByApiKey(ctx, apiKeyHash)
	if err != nil {
		return nil, err
	}

	return merchant, nil
}

func (u *merchantUC) Register(c context.Context, req *domain.RegisterMerchantRequest) (*domain.Merchant, error) {
	// create context with timeout
	ctx, cancel := context.WithTimeout(c, u.timeout)
	defer cancel()

	// generate new UUIDV7
	id, err := pkg.GenerateUUIDV7()
	if err != nil {
		return nil, err
	}

	// generate ApiKey
	apiKey, err := pkg.GenerateApiKey("mch")
	if err != nil {
		return nil, err
	}

	// hash ApiKey
	apiKeyHash := pkg.HashKey(apiKey)

	// create merchant entity
	merchant := &domain.Merchant{
		ID:          id,
		Name:        req.Name,
		Email:       req.Email,
		ApiKey:      apiKey,
		APIKeyHash:  apiKeyHash,
		CallbackURL: req.CallbackURL,
		Status:      domain.MerchantStatusActive,
		Balance:     0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// save to repository
	if err := u.merchantRepo.Create(ctx, merchant); err != nil {
		return nil, err

	}

	return merchant, nil
}
