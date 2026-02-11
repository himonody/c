package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret []byte
var expireHours int
var refreshThresholdHours int

type Claims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// InitJWT initializes JWT configuration
func InitJWT(secret string, expire, refreshThreshold int) {
	jwtSecret = []byte(secret)
	expireHours = expire
	refreshThresholdHours = refreshThreshold
}

// GenerateToken generates a JWT token
func GenerateToken(userID int64, username string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(time.Duration(expireHours) * time.Hour)

	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			IssuedAt:  jwt.NewNumericDate(nowTime),
			NotBefore: jwt.NewNumericDate(nowTime),
			Issuer:    "habit",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ParseToken parses and validates a JWT token
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// ShouldRefreshToken checks if token should be refreshed based on expiry time
func ShouldRefreshToken(claims *Claims) bool {
	if claims == nil || claims.ExpiresAt == nil {
		return false
	}

	// Calculate time until expiration
	now := time.Now()
	expiresAt := claims.ExpiresAt.Time
	timeUntilExpiry := expiresAt.Sub(now)

	// Refresh if token expires within threshold
	thresholdDuration := time.Duration(refreshThresholdHours) * time.Hour
	return timeUntilExpiry <= thresholdDuration && timeUntilExpiry > 0
}

// GetTokenExpiration returns token expiration duration in seconds
func GetTokenExpiration() int64 {
	return int64(expireHours * 3600)
}
