package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"
	"time"
)

// GenerateRefCode generates a unique reference code for user
func GenerateRefCode(userID int64) string {
	// Combine timestamp and user ID for uniqueness
	timestamp := time.Now().UnixNano()
	str := fmt.Sprintf("%d%d", userID, timestamp)

	// Generate random bytes
	b := make([]byte, 4)
	rand.Read(b)

	// Combine and encode
	code := fmt.Sprintf("%s%s", str, base64.URLEncoding.EncodeToString(b))

	// Clean up and make it uppercase
	code = strings.ReplaceAll(code, "-", "")
	code = strings.ReplaceAll(code, "_", "")
	code = strings.ToUpper(code)

	// Limit to 20 characters
	if len(code) > 20 {
		code = code[:20]
	}

	return code
}
