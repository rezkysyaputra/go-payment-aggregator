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
			name: "Success Get Merchant",
			mock: func(repo *mocks.MockMerchantRepository) {
				repo.On("FindByID", mock.Anything, merchantID).Return(returnedMerchant, nil)
			},
			wantErr: false,
		},
		{
			name: "Failed Get Merchant - Not Found",
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
