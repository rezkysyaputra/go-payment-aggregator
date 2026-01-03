package pkg

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
)

// GenerateApiKey generates a random API key with an optional prefix.
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

// HashKey hashes the given raw key using SHA-256 and returns the hexadecimal representation.
func HashKey256(rawKey string) string {
	hasher := sha256.New()
	hasher.Write([]byte(rawKey))

	return hex.EncodeToString(hasher.Sum(nil))
}

func HashKey512(rawKey string) string {
	hasher := sha512.New()
	hasher.Write([]byte(rawKey))

	return hex.EncodeToString(hasher.Sum(nil))
}
