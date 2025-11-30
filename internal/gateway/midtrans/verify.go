package midtrans

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"strings"
)

// ComputeSignature computes the SHA512 signature for Midtrans transaction verification
func ComputeSignature(orderId, statusCode, grossAmount, serverKey string) string {
	data := orderId + statusCode + grossAmount + serverKey
	// Compute SHA512 hash
	sum := sha512.Sum512([]byte(data))
	// Convert to hexadecimal string
	return hex.EncodeToString(sum[:])
}

// VerifySignature verifies the SHA512 signature from Midtrans notification
func VerifySignature(payloadSignature, orderId, statusCode, grossAmount, serverKey string) bool {
	if payloadSignature == "" {
		return false
	}

	computed := ComputeSignature(orderId, statusCode, grossAmount, serverKey)
	// Compare signatures in a case-insensitive manner
	return strings.EqualFold(payloadSignature, computed)
}

func GrossAmountToString(gross any) (string, error) {
	switch v := gross.(type) {
	case string:
		return v, nil
	case float64:
		if v == float64(int64(v)) {
			return fmt.Sprintf("%d", int64(v)), nil
		}
		return fmt.Sprintf("%v", v), nil
	case int:
		return fmt.Sprintf("%d", v), nil
	default:
		return "", fmt.Errorf("unsupported gross amount type: %T", gross)
	}
}
