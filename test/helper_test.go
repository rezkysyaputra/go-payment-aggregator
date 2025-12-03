package test

import (
	"go-payment-aggregator/internal/domain/merchant"
	"go-payment-aggregator/internal/domain/transaction"

	"github.com/google/uuid"
)

func ClearAll() {
	ClearMerchant()
	ClearTransaction()
}

func ClearMerchant() {
	err := db.Where("id is not null").Delete(&merchant.Merchant{}).Error
	if err != nil {
		log.Fatalf("failed to clear merchant: %v", err)
	}

}
func ClearTransaction() {
	err := db.Where("id is not null").Delete(&transaction.Transaction{}).Error
	if err != nil {
		log.Fatalf("failed to clear transaction: %v", err)
	}
}

func CreateMerchant() *merchant.Merchant {
	merchant := &merchant.Merchant{
		Name:        "Test Merchant",
		ApiKey:      "mch_test_api_key",
		CallbackURL: "http://localhost:8080",
	}
	err := db.Create(merchant).Error
	if err != nil {
		log.Fatalf("failed to create test merchant: %v", err)
	}

	return merchant
}

func CreateTransaction(merchantID uuid.UUID, amount float64, provider string) *transaction.Transaction {
	t := &transaction.Transaction{
		MerchantID: merchantID,
		OrderID:    "order_123",
		Amount:     amount,
		Provider:   provider}
	err := db.Create(t).Error
	if err != nil {
		log.Fatalf("failed to create test transaction: %v", err)
	}

	return t
}
