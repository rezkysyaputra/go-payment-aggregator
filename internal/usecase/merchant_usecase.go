package usecase

import (
	"context"
	"errors"
	"go-payment-aggregator/internal/domain"
	"go-payment-aggregator/internal/pkg"
	"time"

	"github.com/google/uuid"
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
	apiKeyHash := pkg.HashKey256(apiKey)

	merchant, err := u.merchantRepo.FindByApiKey(ctx, apiKeyHash)
	if err != nil {
		return nil, err
	}

	if merchant.Status != domain.MerchantStatusActive {
		return nil, errors.New("merchant is not active")
	}

	return merchant, nil
}

func (u *merchantUC) Register(c context.Context, req *domain.RegisterMerchantRequest) (*domain.Merchant, error) {
	ctx, cancel := context.WithTimeout(c, u.timeout)
	defer cancel()

	id := pkg.GenerateUUIDV7()

	apiKey := pkg.GenerateApiKey("mch")

	apiKeyHash := pkg.HashKey256(apiKey)

	merchant := &domain.Merchant{
		ID:          id,
		Name:        req.Name,
		Email:       req.Email,
		ApiKey:      apiKey,
		APIKeyHash:  apiKeyHash,
		CallbackURL: req.CallbackURL,
		Status:      domain.MerchantStatusActive,
		Balance:     0,
	}

	createdMerchant, err := u.merchantRepo.Create(ctx, merchant)
	if err != nil {
		return nil, err
	}

	createdMerchant.ApiKey = apiKey

	return createdMerchant, nil
}

func (u *merchantUC) GetProfile(ctx context.Context, id uuid.UUID) (*domain.Merchant, error) {
	merchant, err := u.merchantRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return merchant, nil
}

func (u *merchantUC) UpdateProfile(ctx context.Context, id uuid.UUID, req *domain.UpdateMerchantRequest) (*domain.Merchant, error) {
	merchant, err := u.merchantRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		merchant.Name = req.Name
	}

	if req.CallbackURL != "" {
		merchant.CallbackURL = req.CallbackURL
	}

	merchant.UpdatedAt = time.Now()

	if err := u.merchantRepo.Update(ctx, merchant); err != nil {
		return nil, err
	}

	return merchant, nil
}

func (u *merchantUC) RegenerateApiKey(ctx context.Context, id uuid.UUID) (string, error) {
	newApiKey := pkg.GenerateApiKey("mch")

	newApiKeyHash := pkg.HashKey256(newApiKey)

	if err := u.merchantRepo.RegenerateApiKey(ctx, id, newApiKeyHash); err != nil {
		return "", err
	}

	return newApiKey, nil
}
