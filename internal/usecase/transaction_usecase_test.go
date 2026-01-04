package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"go-payment-aggregator/internal/domain"
	"go-payment-aggregator/internal/mocks"
	"go-payment-aggregator/internal/pkg"
	"go-payment-aggregator/internal/usecase"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTransactionUsecase_CreateTransaction(t *testing.T) {
	merchantID := pkg.GenerateUUIDV7()

	reqUC := &domain.CreateTransactionRequest{
		OrderID:       "ORDER-TEST-123",
		Amount:        100000,
		Currency:      "IDR",
		Provider:      "midtrans",
		PaymentMethod: "credit_card",
		Customer: domain.Customer{
			Name:  "user",
			Email: "user@example.com",
		},
		Items: []domain.Item{
			{Name: "Item 1", Quantity: 1, Price: 100000},
		},
	}

	reqGateway := &domain.CreatePaymentRequest{
		OrderID:       "ORDER-TEST-123",
		Amount:        100000,
		PaymentMethod: "credit_card",
		ExpiryMinutes: 2,
		Customer:      reqUC.Customer,
		Items:         reqUC.Items,
	}

	mockTransaction := &domain.Transaction{
		OrderID:       reqUC.OrderID,
		Amount:        reqUC.Amount,
		PaymentMethod: reqUC.PaymentMethod,
		Currency:      reqUC.Currency,
		Status:        domain.TransactionStatusPending,
	}

	paymentResponse := &domain.PaymentResponse{
		PaymentURL: "https://payment-gateway.com/pay/12345",
		Token:      "ext-12345",
	}

	matchGatewayRequest := mock.MatchedBy(func(req *domain.CreatePaymentRequest) bool {
		return req.OrderID == reqGateway.OrderID &&
			req.Amount == reqGateway.Amount &&
			req.PaymentMethod == reqGateway.PaymentMethod &&
			req.Customer == reqGateway.Customer &&
			(req.ExpiryMinutes == 1 || req.ExpiryMinutes == 2)
	})

	tests := []struct {
		name    string
		mock    func(repo *mocks.MockTransactionRepository, gateway *mocks.MockPaymentGateway)
		wantErr bool
	}{
		{
			name: "Success Create Transaction",
			mock: func(repo *mocks.MockTransactionRepository, gateway *mocks.MockPaymentGateway) {
				repo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Transaction")).
					Return(mockTransaction, nil)

				gateway.On("CreatePayment", matchGatewayRequest).
					Return(paymentResponse, nil)

				repo.On("Update", mock.Anything, mock.AnythingOfType("*domain.Transaction")).
					Return(mockTransaction, nil)
			},
			wantErr: false,
		},
		{
			name: "Failed Repository Create",
			mock: func(repo *mocks.MockTransactionRepository, gateway *mocks.MockPaymentGateway) {
				repo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Transaction")).
					Return(nil, errors.New("database error"))
			},
			wantErr: true,
		},
		{
			name: "Failed Payment Gateway",
			mock: func(repo *mocks.MockTransactionRepository, gateway *mocks.MockPaymentGateway) {
				repo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Transaction")).
					Return(mockTransaction, nil)

				gateway.On("CreatePayment", matchGatewayRequest).
					Return(nil, errors.New("gateway timeout"))
			},
			wantErr: true,
		},
		{
			name: "Failed Repository Update",
			mock: func(repo *mocks.MockTransactionRepository, gateway *mocks.MockPaymentGateway) {
				repo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Transaction")).
					Return(mockTransaction, nil)

				gateway.On("CreatePayment", matchGatewayRequest).
					Return(paymentResponse, nil)

				repo.On("Update", mock.Anything, mock.AnythingOfType("*domain.Transaction")).
					Return(nil, errors.New("database update error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.MockTransactionRepository)
			mockGateway := new(mocks.MockPaymentGateway)

			tt.mock(mockRepo, mockGateway)

			transactionUC := usecase.NewTransactionUC(mockRepo, mockGateway, time.Second*2)

			ctx := context.Background()
			res, err := transactionUC.Create(ctx, merchantID, reqUC)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, domain.TransactionStatusPending, res.Status)
				assert.Equal(t, paymentResponse.PaymentURL, res.PaymentURL)
			}

			mockRepo.AssertExpectations(t)
			mockGateway.AssertExpectations(t)
		})
	}
}

func TestTransactionUsecase_GetTransaction(t *testing.T) {
	transactionID := pkg.GenerateUUIDV7()
	mockTransaction := &domain.Transaction{
		ID:            transactionID,
		OrderID:       "ORDER-TEST-123",
		Amount:        100000,
		PaymentMethod: "credit_card",
		Currency:      "IDR",
		Status:        domain.TransactionStatusPending,
	}

	tests := []struct {
		name    string
		mock    func(repo *mocks.MockTransactionRepository)
		wantErr bool
	}{
		{
			name: "Success Get Transaction",
			mock: func(repo *mocks.MockTransactionRepository) {
				repo.On("Get", mock.Anything, transactionID).Return(mockTransaction, nil)
			},
			wantErr: false,
		},
		{
			name: "Failed Get Transaction - Not Found",
			mock: func(repo *mocks.MockTransactionRepository) {
				repo.On("Get", mock.Anything, transactionID).Return(nil, errors.New("Transaction not found"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.MockTransactionRepository)
			mockGateway := new(mocks.MockPaymentGateway)

			tt.mock(mockRepo)

			transactionUC := usecase.NewTransactionUC(mockRepo, mockGateway, time.Second*2)

			ctx := context.Background()
			res, err := transactionUC.Get(ctx, transactionID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, mockTransaction.ID, res.ID)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestTransactionUsecase_HandleNotification(t *testing.T) {
	orderID := "ORDER-TEST-123"
	mockTransaction := &domain.Transaction{
		ID:      pkg.GenerateUUIDV7(),
		OrderID: orderID,
		Amount:  100000,
		Status:  domain.TransactionStatusPending,
	}

	req := &domain.UpdateStatusRequest{
		OrderID: orderID,
		Status:  "PAID",
	}

	tests := []struct {
		name    string
		mock    func(repo *mocks.MockTransactionRepository)
		wantErr bool
	}{
		{
			name: "Success Handle Notification",
			mock: func(repo *mocks.MockTransactionRepository) {
				repo.On("FindByOrderID", mock.Anything, orderID).
					Return(mockTransaction, nil)

				repo.On("Update", mock.Anything, mock.MatchedBy(func(tx *domain.Transaction) bool {
					return tx.OrderID == orderID && tx.Status == domain.TransactionStatusPaid
				})).Return(mockTransaction, nil)
			},
			wantErr: false,
		},
		{
			name: "Transaction Not Found",
			mock: func(repo *mocks.MockTransactionRepository) {
				repo.On("FindByOrderID", mock.Anything, orderID).
					Return(nil, errors.New("transaction not found"))
			},
			wantErr: true,
		},
		{
			name: "Failed to Update Transaction",
			mock: func(repo *mocks.MockTransactionRepository) {
				repo.On("FindByOrderID", mock.Anything, orderID).
					Return(mockTransaction, nil)

				repo.On("Update", mock.Anything, mock.Anything).
					Return(nil, errors.New("database error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.MockTransactionRepository)
			mockGateway := new(mocks.MockPaymentGateway)

			tt.mock(mockRepo)

			transactionUC := usecase.NewTransactionUC(mockRepo, mockGateway, time.Second*2)

			ctx := context.Background()
			err := transactionUC.HandleNotification(ctx, req)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
