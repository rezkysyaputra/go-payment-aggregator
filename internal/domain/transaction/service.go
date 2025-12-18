package transaction

import (
	"context"
	"encoding/json"
	"fmt"
	"go-payment-aggregator/internal/gateway"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type TransactionService interface {
	CreateTransaction(MerchantID uuid.UUID, orderID string, amount float64, provider string) (*Transaction, error)
	ProcessNotification(provider string, payload NotificationPayload) error
	FindById(id uuid.UUID) (*Transaction, error)
	FindOrderById(orderID string) (*Transaction, error)
	UpdateStatusAndRaw(id uuid.UUID, status string, rawJSON string) error
}

type NotificationPayload map[string]any

type TransactionServiceImpl struct {
	TransactionRepository TransactionRepository
	Gateways              map[string]gateway.PaymentGateway
	Config                *viper.Viper
	Log                   *logrus.Logger
	Redis                 *redis.Client
}

func NewTransactionService(repo TransactionRepository, gateways map[string]gateway.PaymentGateway, config *viper.Viper, log *logrus.Logger, redis *redis.Client) TransactionService {
	return &TransactionServiceImpl{
		TransactionRepository: repo,
		Gateways:              gateways,
		Config:                config,
		Log:                   log,
		Redis:                 redis,
	}
}

func (s *TransactionServiceImpl) sendWebhookToMerchant(tx *Transaction) {
	if tx.Merchant == nil || tx.Merchant.CallbackURL == "" {
		return
	}

	payload := map[string]any{
		"transaction_id": tx.ID,
		"order_id":       tx.OrderID,
		"status":         tx.Status,
		"amount":         tx.Amount,
		"provider":       tx.Provider,
		"callback_url":   tx.Merchant.CallbackURL, // Include URL in payload for worker
	}

	body, _ := json.Marshal(payload)

	// Push to Redis Queue
	ctx := context.Background()
	// Use RPUSH to add to the tail of the queue
	err := s.Redis.RPush(ctx, "webhook_queue", body).Err()
	if err != nil {
		s.Log.Errorf("Failed to push webhook task to redis: %v", err)
		// Fallback: log error, maybe try direct http? For now just log.
		return
	}
	s.Log.Infof("Webhook task enqueued for Order ID: %s", tx.OrderID)
}

// CreateTransaction creates a new transaction
func (s *TransactionServiceImpl) CreateTransaction(merchantID uuid.UUID, orderID string, amount float64, provider string) (*Transaction, error) {
	// check duplicate order id
	if exists, _ := s.TransactionRepository.FindOrderById(orderID); exists != nil {
		return nil, fmt.Errorf("duplicate order_id: %s", orderID)
	}

	// create transaction model
	t := &Transaction{
		ID:         uuid.New(),
		MerchantID: merchantID,
		OrderID:    orderID,
		Amount:     amount,
		Provider:   provider,
		Status:     "pending",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// select gateway
	gw, ok := s.Gateways[provider]
	if !ok {
		return nil, fmt.Errorf("payment provider not found: %s", provider)
	}

	// call gateway implementation
	res, err := gw.CreateTransaction(orderID, amount)
	if err != nil {
		s.Log.Errorf("Error calling gateway CreateTransaction: %v", err)
		return nil, err
	}

	// set response data
	t.RedirectURL = res.RedirectURL
	t.ExternalRef = res.Token

	// save to database
	if err := s.TransactionRepository.Create(t); err != nil {
		s.Log.Errorf("Error saving transaction to database: %v", err)
		return nil, err
	}

	return t, nil
}

// ProcessNotification processes generic payment notification
func (s *TransactionServiceImpl) ProcessNotification(provider string, payload NotificationPayload) error {
	// select gateway
	gw, ok := s.Gateways[provider]
	if !ok {
		return fmt.Errorf("payment provider not found: %s", provider)
	}

	// 1. Verify Signature
	if !gw.VerifySignature(payload) {
		return fmt.Errorf("invalid signature")
	}

	// 2. Get Order ID
	orderID, err := gw.GetOrderID(payload)
	if err != nil {
		return fmt.Errorf("failed to get order id: %v", err)
	}

	// 3. Get Status
	status, err := gw.GetStatus(payload)
	if err != nil {
		return fmt.Errorf("failed to get status: %v", err)
	}

	// 4. Find Transaction
	tx, err := s.TransactionRepository.FindOrderById(orderID)
	if err != nil {
		return fmt.Errorf("transaction not found: %v", err)
	}

	// check if transaction is already in final status or same status
	if isFinalStatus(tx.Status) || status == tx.Status {
		return nil
	}

	// store raw payload as JSON
	rawB, _ := json.Marshal(payload)
	rawJSON := string(rawB)

	// 5. Update Transaction
	if err := s.TransactionRepository.UpdateStatusAndRaw(tx.ID, status, rawJSON); err != nil {
		return fmt.Errorf("failed to update transaction: %v", err)
	}

	// 6. Send Webhook to Merchant
	tx.Status = status // update status in struct for webhook payload
	s.sendWebhookToMerchant(tx)

	return nil
}

func isFinalStatus(status string) bool {
	return status == "paid" || status == "failed"
}

// FindById implements TransactionService.
func (s *TransactionServiceImpl) FindById(id uuid.UUID) (*Transaction, error) {
	return s.TransactionRepository.FindById(id)
}

// FindOrderById implements TransactionService.
func (s *TransactionServiceImpl) FindOrderById(orderID string) (*Transaction, error) {
	return s.TransactionRepository.FindOrderById(orderID)
}

// UpdateStatusAndRaw implements TransactionService.
func (s *TransactionServiceImpl) UpdateStatusAndRaw(id uuid.UUID, status string, rawJSON string) error {
	return s.TransactionRepository.UpdateStatusAndRaw(id, status, rawJSON)
}
