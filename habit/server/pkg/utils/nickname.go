package utils

import (
	"crypto/rand"
	"math/big"
)

const (
	// Characters allowed in nickname: alphanumeric only
	nicknameChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

// GenerateRandomNickname generates a random 6-character nickname
func GenerateRandomNickname() string {
	result := make([]byte, 6)
	charsLen := big.NewInt(int64(len(nicknameChars)))

	for i := 0; i < 6; i++ {
		num, err := rand.Int(rand.Reader, charsLen)
		if err != nil {
			// Fallback to a simple pattern if crypto/rand fails
			result[i] = nicknameChars[i%len(nicknameChars)]
			continue
		}
		result[i] = nicknameChars[num.Int64()]
	}

	return string(result)
}
