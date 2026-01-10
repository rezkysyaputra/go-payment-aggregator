package pkg

func MapXenditStatus(xenditStatus string) string {
	switch xenditStatus {
	case "PAID":
		return "PAID"
	case "PENDING":
		return "PENDING"
	case "EXPIRED", "CANCELLED":
		return "FAILED"
	default:
		return "PENDING"
	}
}

func VerifySignature(orderID, statusCode, grossAmount, serverKey, signatureKey string) bool {
	signatureString := orderID + statusCode + grossAmount + serverKey
	expectedSignature := HashKey512(signatureString)

	return expectedSignature == signatureKey
}
