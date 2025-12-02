package test

import "go-payment-aggregator/internal/domain/merchant"

func ClearAll() {
	ClearMerchant()
}

func ClearMerchant() {
	err := db.Where("id is not null").Delete(&merchant.Merchant{}).Error
	if err != nil {
		log.Fatalf("failed to clear merchant: %v", err)
	}
}
