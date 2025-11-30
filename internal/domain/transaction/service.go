package transaction

import (
	"encoding/json"
	"fmt"
	"go-payment-aggregator/internal/gateway/midtrans"
	"log"
	"time"

	"github.com/google/uuid"
)

type TransactionService interface {
	CreateTransaction(MerchantID uuid.UUID, orderID string, amount float64, provider string) (*Transaction, error)
	ProcessMidtransNotification(payload NotificationPayload) error
	FindById(id uuid.UUID) (*Transaction, error)
	FindOrderById(orderID string) (*Transaction, error)
	UpdateStatusAndRaw(id uuid.UUID, status string, rawJSON string) error
}

type NotificationPayload map[string]any

type TransactionServiceImpl struct {
	TransactionRepository TransactionRepository
	ServerKey             string
}

func NewTransactionService(repo TransactionRepository, serverKey string) TransactionService {
	return &TransactionServiceImpl{TransactionRepository: repo, ServerKey: serverKey}
}

// CreateTransaction creates a new transaction
func (s *TransactionServiceImpl) CreateTransaction(merchantID uuid.UUID, orderID string, amount float64, provider string) (*Transaction, error) {
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

	// call Midtrans snap API
	res, err := midtrans.CreateTransaction(s.ServerKey, orderID, amount)
	if err != nil {
		log.Printf("Error calling Midtrans CreateTransaction: %v", err)
		return nil, err
	}

	// set response data
	t.RedirectURL = res.RedirectURL
	t.ExternalRef = res.Token

	// save to database
	if err := s.TransactionRepository.Create(t); err != nil {
		log.Printf("Error saving transaction to database: %v", err)
		return nil, err
	}

	return t, nil
}

// ProcessMidtransNotification processes Midtrans payment notification
// - verifies the signature
// - find the transaction by order ID
// - updates the transaction status and raw response
func (s *TransactionServiceImpl) ProcessMidtransNotification(payload NotificationPayload) error {
	// extract necessary fields
	orderID, ok := payload["order_id"].(string)
	if !ok {
		return fmt.Errorf("order_id not found in payload")
	}

	statusCode, _ := payload["status_code"].(string)

	// gross_amount can be string or float64
	grossStr, err := midtrans.GrossAmountToString(payload["gross_amount"])
	if err != nil {
		return fmt.Errorf("invalid gross_amount: %v", err)
	}

	payloadSignature := ""
	if sig, exists := payload["signature_key"].(string); exists {
		payloadSignature = sig
	}

	// verify signature
	if !midtrans.VerifySignature(payloadSignature, orderID, statusCode, grossStr, s.ServerKey) {
		return fmt.Errorf("invalid signature")
	}

	// map Midtrans status to internal status
	midtransStatus, _ := payload["transaction_status"].(string)
	internalStatus := mapMidtransToInternalStatus(midtransStatus)

	// find transaction by order ID
	tx, err := s.TransactionRepository.FindOrderById(orderID)
	if err != nil {
		return fmt.Errorf("transaction not found: %v", err)
	}

	// check if transaction is already in final status or same status
	if isFinalStatus(tx.Status) || internalStatus == tx.Status {
		return nil
	}

	// store raw payload as JSON
	rawB, _ := json.Marshal(payload)
	rawJSON := string(rawB)

	// update transaction status and raw response
	if err := s.TransactionRepository.UpdateStatusAndRaw(tx.ID, internalStatus, rawJSON); err != nil {
		return fmt.Errorf("failed to update transaction: %v", err)
	}

	return nil
}

func mapMidtransToInternalStatus(midtransStatus string) string {
	switch midtransStatus {
	case "capture", "settlement":
		return "paid"
	case "deny", "expire", "cancel":
		return "failed"
	case "pending":
		return "pending"
	default:
		return "unknown"
	}
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
