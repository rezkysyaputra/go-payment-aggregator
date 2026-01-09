package usecase

import (
	"context"
	"errors"
	"go-payment-aggregator/internal/domain"
	"go-payment-aggregator/internal/pkg"
	"time"

	"github.com/google/uuid"
)

type TransactionUC struct {
	transactionRepo domain.TransactionRepository
	gateways        map[string]domain.PaymentGateway
	timeout         time.Duration
}

func NewTransactionUC(r domain.TransactionRepository, g map[string]domain.PaymentGateway, t time.Duration) domain.TransactionUC {
	return &TransactionUC{
		transactionRepo: r,
		gateways:        g,
		timeout:         t,
	}
}

func (u *TransactionUC) Create(ctx context.Context, merchantID uuid.UUID, req *domain.CreateTransactionRequest) (*domain.Transaction, error) {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	id := pkg.GenerateUUIDV7()

	expiryDuration := time.Minute * 2

	transaction := &domain.Transaction{
		ID:            id,
		MerchantID:    merchantID,
		OrderID:       req.OrderID,
		Provider:      req.Provider,
		Amount:        req.Amount,
		Currency:      req.Currency,
		Status:        domain.TransactionStatusPending,
		PaymentMethod: req.PaymentMethod,
		ExpiredAt:     time.Now().Add(expiryDuration),
	}

	createdTransaction, err := u.transactionRepo.Create(ctx, transaction)
	if err != nil {
		return nil, err
	}

	paymentRequest := &domain.CreatePaymentRequest{
		OrderID:       createdTransaction.OrderID,
		Amount:        createdTransaction.Amount,
		PaymentMethod: createdTransaction.PaymentMethod,
		ExpiryMinutes: int64(time.Until(transaction.ExpiredAt).Minutes()),
		Customer:      req.Customer,
		Items:         req.Items,
	}

	gateway, exists := u.gateways[req.Provider]
	if !exists {
		return nil, errors.New("payment provider not supported")
	}

	paymentResponse, err := gateway.CreatePayment(paymentRequest)
	if err != nil {
		return nil, err
	}

	rawJsonResponse := pkg.ToJSON(paymentResponse)

	createdTransaction.PaymentURL = paymentResponse.PaymentURL
	createdTransaction.ExternalID = paymentResponse.Token
	createdTransaction.RawResponse = string(rawJsonResponse)

	updatedTransaction, err := u.transactionRepo.Update(ctx, createdTransaction)
	if err != nil {
		return nil, err
	}

	return updatedTransaction, nil
}

func (u *TransactionUC) Get(ctx context.Context, id uuid.UUID) (*domain.Transaction, error) {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	getTransaction, err := u.transactionRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return getTransaction, nil
}

func (u *TransactionUC) HandleNotification(ctx context.Context, req *domain.UpdateStatusRequest) error {
	tx, err := u.transactionRepo.FindByOrderID(ctx, req.OrderID)
	if err != nil {
		return err
	}

	tx.Status = domain.TransactionStatus(req.Status)

	_, err = u.transactionRepo.Update(ctx, tx)
	if err != nil {
		return err
	}
	return nil
}
