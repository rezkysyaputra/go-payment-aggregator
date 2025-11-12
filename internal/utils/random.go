package utils

import (
	"crypto/rand"
	"math/big"
)

var base64char = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/")

func RandomBase64(n int) (string, error) {
	result := make([]byte, n)
	for i := 0; i < n; i++ {
		random, err := rand.Int(rand.Reader, big.NewInt(int64(len(base64char))))
		if err != nil {
			return "", err
		}
		result[i] = byte(base64char[random.Int64()])
	}
	return string(result), nil
}
