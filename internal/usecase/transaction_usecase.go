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
	gateway         domain.PaymentGateway
	timeout         time.Duration
}

func NewTransactionUC(r domain.TransactionRepository, g domain.PaymentGateway, t time.Duration) domain.TransactionUC {
	return &TransactionUC{
		transactionRepo: r,
		gateway:         g,
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

	// save transaction to repository
	createdTransaction, err := u.transactionRepo.Create(ctx, transaction)
	if err != nil {
		return nil, err
	}

	// prepare payment request
	paymentRequest := &domain.CreatePaymentRequest{
		OrderID:       createdTransaction.OrderID,
		Amount:        createdTransaction.Amount,
		PaymentMethod: createdTransaction.PaymentMethod,
		ExpiryMinutes: int64(time.Until(transaction.ExpiredAt).Minutes()),
		Customer:      req.Customer,
		Items:         req.Items,
	}

	// process payment via gateway
	paymentResponse, err := u.gateway.CreatePayment(paymentRequest)
	if err != nil {
		return nil, err
	}

	rawJsonResponse, err := pkg.ToJSON(paymentResponse)
	if err != nil {
		return nil, err
	}
	// update transaction with payment details
	createdTransaction.PaymentURL = paymentResponse.PaymentURL
	createdTransaction.ExternalID = paymentResponse.Token
	createdTransaction.RawResponse = string(rawJsonResponse)

	// save updated transaction to repository
	updatedTransaction, err := u.transactionRepo.Update(ctx, createdTransaction)
	if err != nil {
		return nil, err
	}

	return updatedTransaction, nil
}

func (u *TransactionUC) Get(ctx context.Context, id uuid.UUID) (*domain.Transaction, error) {
	// create context with timeout
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	// get transaction from repository
	getTransaction, err := u.transactionRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return getTransaction, nil
}
