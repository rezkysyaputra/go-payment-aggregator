package utils

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateApiKey(prefix string) (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	if prefix == "" {
		return hex.EncodeToString(b), nil
	}

	return prefix + "_" + hex.EncodeToString(b), nil
}
