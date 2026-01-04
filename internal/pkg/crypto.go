package pkg

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
)

func GenerateApiKey(prefix string) string {
	b := make([]byte, 32)
	rand.Read(b)

	if prefix == "" {
		return hex.EncodeToString(b)
	}

	return prefix + "_" + hex.EncodeToString(b)
}

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
