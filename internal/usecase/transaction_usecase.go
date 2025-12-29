package usecase

import (
	"context"
	"go-payment-aggregator/internal/domain"
	"go-payment-aggregator/internal/pkg"
	"time"

	"github.com/google/uuid"
)

type TransactionUC struct {
	transactionRepo domain.TransactionRepository
	timeout         time.Duration
}

func NewTransactionUC(r domain.TransactionRepository, t time.Duration) domain.TransactionUC {
	return &TransactionUC{
		transactionRepo: r,
		timeout:         t,
	}
}

func (u *TransactionUC) Create(ctx context.Context, merchantID uuid.UUID, req *domain.CreateTransactionRequest) (*domain.Transaction, error) {
	// create context with timeout
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	// generate UUIDV7 for transaction ID
	id, err := pkg.GenerateUUIDV7()
	if err != nil {
		return nil, err
	}

	transaction := &domain.Transaction{
		ID:            id,
		MerchantID:    merchantID,
		OrderID:       req.OrderID,
		Provider:      req.Provider,
		Amount:        req.Amount,
		Currency:      req.Currency,
		PaymentMethod: req.PaymentMethod,
		Status:        domain.TransactionStatusPending,
	}

	// save transaction to repository
	createdTransaction, err := u.transactionRepo.Create(ctx, transaction)
	if err != nil {
		return nil, err
	}
	return createdTransaction, nil
}
