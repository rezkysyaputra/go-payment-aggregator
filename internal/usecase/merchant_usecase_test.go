package usecase_test

import (
	"context"
	"errors"
	"go-payment-aggregator/internal/domain"
	"go-payment-aggregator/internal/mocks"
	"go-payment-aggregator/internal/pkg"
	"go-payment-aggregator/internal/usecase"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMerchantUsecase_Register(t *testing.T) {
	reqUC := &domain.RegisterMerchantRequest{
		Name:        "Test Merchant",
		Email:       "test@example.com",
		CallbackURL: "https://example.com/callback",
	}

	returnedMerchant := &domain.Merchant{
		Name:        reqUC.Name,
		Email:       reqUC.Email,
		CallbackURL: reqUC.CallbackURL,
		Status:      domain.MerchantStatusActive,
		Balance:     0,
	}

	tests := []struct {
		name    string
		mock    func(repo *mocks.MockMerchantRepository)
		wantErr bool
	}{
		{
			name: "Success Register Merchant",
			mock: func(repo *mocks.MockMerchantRepository) {

				repo.On("Create", mock.Anything, mock.MatchedBy(func(m *domain.Merchant) bool {

					return m.Name == reqUC.Name &&
						m.Email == reqUC.Email &&
						m.CallbackURL == reqUC.CallbackURL &&
						m.Status == domain.MerchantStatusActive
				})).Return(returnedMerchant, nil)
			},
			wantErr: false,
		},
		{
			name: "Failed Repository Register Merchant",
			mock: func(repo *mocks.MockMerchantRepository) {
				repo.On("Create", mock.Anything, mock.MatchedBy(func(m *domain.Merchant) bool {
					return m.Name == reqUC.Name
				})).Return(nil, assert.AnError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.MockMerchantRepository)

			tt.mock(mockRepo)

			merchantUC := usecase.NewMerchantUC(mockRepo, time.Second*2)

			ctx := context.Background()
			res, err := merchantUC.Register(ctx, reqUC)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, returnedMerchant.Name, res.Name)
				assert.Equal(t, returnedMerchant.Status, res.Status)
				assert.NotEmpty(t, res.ApiKey)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestMerchantUsecase_GetProfile(t *testing.T) {
	merchantID := pkg.GenerateUUIDV7()
	returnedMerchant := &domain.Merchant{
		ID:          merchantID,
		Name:        "Merchant Test",
		Email:       "merchant@example.com",
		CallbackURL: "https://example.com/callback",
		ApiKey:      "merchant-api-key",
		Status:      domain.MerchantStatusActive,
		Balance:     1000,
	}

	tests := []struct {
		name    string
		mock    func(repo *mocks.MockMerchantRepository)
		wantErr bool
	}{
		{
			name: "Success Get Merchant Profile",
			mock: func(repo *mocks.MockMerchantRepository) {

				repo.On("FindByID", mock.Anything, merchantID).Return(returnedMerchant, nil)
			},
			wantErr: false,
		},
		{
			name: "Failed Get Merchant Profile - Not Found",
			mock: func(repo *mocks.MockMerchantRepository) {
				repo.On("FindByID", mock.Anything, merchantID).Return(nil, errors.New("Merchant Not Found"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.MockMerchantRepository)

			tt.mock(mockRepo)

			merchantUC := usecase.NewMerchantUC(mockRepo, time.Second*2)

			ctx := context.Background()
			res, err := merchantUC.GetProfile(ctx, merchantID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, returnedMerchant.ID, res.ID)
				assert.Equal(t, returnedMerchant.Name, res.Name)
				assert.Equal(t, returnedMerchant.Balance, res.Balance)
			}
		})
	}
}

func TestMerchantUsecase_UpdateProfile(t *testing.T) {
	reqUC := &domain.UpdateMerchantRequest{
		Name:        "New Merchant Name",
		CallbackURL: "https://new.com/callback",
	}

	merchantID := pkg.GenerateUUIDV7()
	existingMerchant := &domain.Merchant{
		ID:          merchantID,
		Name:        "New Merchant Name",
		Email:       "merchant@example.com",
		CallbackURL: "https://old.com/callback",
		ApiKey:      "merchant-api-key",
		Status:      domain.MerchantStatusActive,
		Balance:     0,
	}

	tests := []struct {
		name    string
		mock    func(repo *mocks.MockMerchantRepository)
		wantErr bool
	}{
		{
			name: "Success Update Merchant Profile",
			mock: func(repo *mocks.MockMerchantRepository) {
				repo.On("FindByID", mock.Anything, merchantID).Return(existingMerchant, nil)

				repo.On("Update", mock.Anything, mock.MatchedBy(func(m *domain.Merchant) bool {
					return m.Name == reqUC.Name && m.CallbackURL == reqUC.CallbackURL
				})).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "Failed Repository Update Merchant Profile - Not Found",
			mock: func(repo *mocks.MockMerchantRepository) {
				repo.On("FindByID", mock.Anything, merchantID).Return(nil, errors.New("Merchant Not Found"))
			},
			wantErr: true,
		},
		{
			name: "Failed Repository Update Merchant Profile - Bad Request",
			mock: func(repo *mocks.MockMerchantRepository) {
				repo.On("FindByID", mock.Anything, merchantID).Return(existingMerchant, nil)

				repo.On("Update", mock.Anything, mock.MatchedBy(func(m *domain.Merchant) bool {
					return m.Name == reqUC.Name && m.CallbackURL == reqUC.CallbackURL
				})).Return(errors.New("Bad Request"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.MockMerchantRepository)

			tt.mock(mockRepo)

			merchantUC := usecase.NewMerchantUC(mockRepo, time.Second*2)

			ctx := context.Background()
			res, err := merchantUC.UpdateProfile(ctx, merchantID, reqUC)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, existingMerchant.Name, res.Name)
				assert.NotEmpty(t, res.ApiKey)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestMerchantUsecase_ValidateApiKey(t *testing.T) {
	merchantID := pkg.GenerateUUIDV7()
	apiKey := pkg.GenerateApiKey("mch")
	apiKeyHash := pkg.HashKey256(apiKey)
	returnedMerchant := &domain.Merchant{
		ID:          merchantID,
		Name:        "Merchant Test",
		Email:       "merchant@example.com",
		CallbackURL: "https://example.com/callback",
		ApiKey:      apiKey,
		APIKeyHash:  apiKeyHash,
		Status:      domain.MerchantStatusActive,
		Balance:     1000,
	}

	tests := []struct {
		name    string
		mock    func(repo *mocks.MockMerchantRepository)
		wantErr bool
	}{
		{
			name: "Success Validate API Key",
			mock: func(repo *mocks.MockMerchantRepository) {
				repo.On("FindByApiKey", mock.Anything, apiKeyHash).Return(returnedMerchant, nil)
			},
			wantErr: false,
		},
		{
			name: "Failed Validate API Key - Merchant Not Found",
			mock: func(repo *mocks.MockMerchantRepository) {
				repo.On("FindByApiKey", mock.Anything, apiKeyHash).Return(nil, errors.New("Merchant Not Found"))
			},
			wantErr: true,
		},
		{
			name: "Failed Validate API Key - Inactive Merchant",
			mock: func(repo *mocks.MockMerchantRepository) {
				inactiveMerchant := *returnedMerchant
				inactiveMerchant.Status = domain.MerchantStatusInactive

				repo.On("FindByApiKey", mock.Anything, apiKeyHash).Return(&inactiveMerchant, nil)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.MockMerchantRepository)

			tt.mock(mockRepo)

			merchantUC := usecase.NewMerchantUC(mockRepo, time.Second*2)

			ctx := context.Background()
			res, err := merchantUC.ValidateApiKey(ctx, apiKey)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, returnedMerchant.ID, res.ID)
				assert.Equal(t, returnedMerchant.Name, res.Name)
				assert.Equal(t, returnedMerchant.Balance, res.Balance)
			}
		})
	}
}
