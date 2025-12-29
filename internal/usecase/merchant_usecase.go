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
	// hash api key
	apiKeyHash := pkg.HashKey(apiKey)

	// find merchant by api key
	merchant, err := u.merchantRepo.FindByApiKey(ctx, apiKeyHash)
	if err != nil {
		return nil, err
	}

	// check if merchant is active
	if merchant.Status != domain.MerchantStatusActive {
		return nil, errors.New("merchant is not active")
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
	}

	// save to repository
	createdMerchant, err := u.merchantRepo.Create(ctx, merchant)
	if err != nil {
		return nil, err
	}

	// set merchant ApiKey for response
	createdMerchant.ApiKey = apiKey

	return createdMerchant, nil
}

func (u *merchantUC) GetProfile(ctx context.Context, id uuid.UUID) (*domain.Merchant, error) {
	// find merchant by ID
	merchant, err := u.merchantRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return merchant, nil
}

func (u *merchantUC) UpdateProfile(ctx context.Context, id uuid.UUID, req *domain.UpdateMerchantRequest) (*domain.Merchant, error) {
	// find merchant by ID
	merchant, err := u.merchantRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// update fields
	if req.Name != "" {
		merchant.Name = req.Name
	}

	if req.CallbackURL != "" {
		merchant.CallbackURL = req.CallbackURL
	}

	merchant.UpdatedAt = time.Now()

	// save updates to repository
	if err := u.merchantRepo.Update(ctx, merchant); err != nil {
		return nil, err
	}

	return merchant, nil
}

func (u *merchantUC) RegenerateApiKey(ctx context.Context, id uuid.UUID) (string, error) {
	// generate new ApiKey
	newApiKey, err := pkg.GenerateApiKey("mch")
	if err != nil {
		return "", err
	}

	// hash new ApiKey
	newApiKeyHash := pkg.HashKey(newApiKey)

	// update api key in repository
	if err := u.merchantRepo.RegenerateApiKey(ctx, id, newApiKeyHash); err != nil {
		return "", err
	}

	return newApiKey, nil
}
