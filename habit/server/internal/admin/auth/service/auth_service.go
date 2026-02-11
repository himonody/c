package service

import (
	"context"
	"errors"

	"habit/internal/admin/auth/dto"
	"habit/internal/repo"
	"habit/pkg/cache"
	"habit/pkg/database"
	"habit/pkg/utils"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AdminAuthService struct {
	adminRepo *repo.AdminRepository
	logger    *zap.Logger
	redis     *redis.Client
}

func NewAdminAuthService(adminRepo *repo.AdminRepository, logger *zap.Logger) *AdminAuthService {
	return &AdminAuthService{
		adminRepo: adminRepo,
		logger:    logger,
		redis:     database.RedisClient,
	}
}

// storeTokenInRedis stores admin token in Redis for single sign-on
func (s *AdminAuthService) storeTokenInRedis(adminID int64, token string) error {
	ctx := context.Background()
	tokenKey := cache.AdminTokenKey(adminID)

	// Delete old token if exists (enforce single sign-on)
	oldToken, err := s.redis.Get(ctx, tokenKey).Result()
	if err == nil && oldToken != "" && oldToken != token {
		// Blacklist old token
		oldTokenKey := cache.AdminTokenBlacklistKey(oldToken)
		s.redis.Set(ctx, oldTokenKey, "1", cache.TokenBlacklistExpiration)
	}

	// Store new token with expiration
	return s.redis.Set(ctx, tokenKey, token, cache.AdminTokenExpiration).Err()
}

// isTokenBlacklisted checks if admin token is in blacklist
func (s *AdminAuthService) isTokenBlacklisted(token string) bool {
	ctx := context.Background()
	blacklistKey := cache.AdminTokenBlacklistKey(token)
	_, err := s.redis.Get(ctx, blacklistKey).Result()
	return err == nil
}

// Login handles admin login
func (s *AdminAuthService) Login(req *dto.AdminLoginRequest, loginIP string) (*dto.AdminLoginResponse, error) {
	// Find admin by username
	admin, err := s.adminRepo.FindByUsername(req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid username or password")
		}
		s.logger.Error("Failed to find admin", zap.Error(err))
		return nil, errors.New("internal server error")
	}

	// Check if admin is active
	if admin.Status != 1 {
		return nil, errors.New("account is disabled")
	}

	// Verify password (using bcrypt)
	if !utils.CheckPassword(req.Password, admin.Password) {
		return nil, errors.New("invalid username or password")
	}

	// Generate JWT token
	token, err := utils.GenerateToken(admin.ID, admin.Username)
	if err != nil {
		s.logger.Error("Failed to generate token", zap.Error(err))
		return nil, errors.New("failed to generate token")
	}

	// Store token in Redis (single sign-on - will invalidate old token)
	if err := s.storeTokenInRedis(admin.ID, token); err != nil {
		s.logger.Error("Failed to store token in Redis", zap.Error(err))
	}

	return &dto.AdminLoginResponse{
		Token: token,
		AdminInfo: dto.AdminUserInfo{
			ID:       admin.ID,
			Username: admin.Username,
			NickName: admin.NickName,
			Role:     admin.Role,
			Status:   admin.Status,
		},
	}, nil
}

// Logout handles admin logout
func (s *AdminAuthService) Logout(adminID int64, token string) error {
	ctx := context.Background()

	// Add token to blacklist
	blacklistKey := cache.AdminTokenBlacklistKey(token)
	if err := s.redis.Set(ctx, blacklistKey, "1", cache.TokenBlacklistExpiration).Err(); err != nil {
		s.logger.Error("Failed to blacklist token", zap.Error(err))
	}

	// Delete token from Redis
	tokenKey := cache.AdminTokenKey(adminID)
	if err := s.redis.Del(ctx, tokenKey).Err(); err != nil {
		s.logger.Error("Failed to delete token from Redis", zap.Error(err))
		return err
	}

	return nil
}

// ValidateToken validates an admin token and returns admin ID and claims
func (s *AdminAuthService) ValidateToken(token string) (int64, *utils.Claims, error) {
	// Check if token is blacklisted
	if s.isTokenBlacklisted(token) {
		return 0, nil, errors.New("token has been revoked")
	}

	claims, err := utils.ParseToken(token)
	if err != nil {
		return 0, nil, err
	}

	// Check if token exists in Redis (single sign-on check)
	ctx := context.Background()
	tokenKey := cache.AdminTokenKey(claims.UserID)
	storedToken, err := s.redis.Get(ctx, tokenKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, nil, errors.New("token not found or expired")
		}
		s.logger.Error("Failed to get token from Redis", zap.Error(err))
		return 0, nil, errors.New("internal server error")
	}

	// Verify token matches (single sign-on)
	if storedToken != token {
		return 0, nil, errors.New("token has been replaced by a new login")
	}

	return claims.UserID, claims, nil
}

// RefreshToken generates a new admin token if needed
func (s *AdminAuthService) RefreshToken(adminID int64, username string, oldToken string) (string, bool, error) {
	// Generate new token
	newToken, err := utils.GenerateToken(adminID, username)
	if err != nil {
		s.logger.Error("Failed to generate new token", zap.Error(err))
		return "", false, err
	}

	// Store new token in Redis
	if err := s.storeTokenInRedis(adminID, newToken); err != nil {
		s.logger.Error("Failed to store new token in Redis", zap.Error(err))
		return "", false, err
	}

	// Blacklist old token
	ctx := context.Background()
	blacklistKey := cache.AdminTokenBlacklistKey(oldToken)
	if err := s.redis.Set(ctx, blacklistKey, "1", cache.TokenBlacklistExpiration).Err(); err != nil {
		s.logger.Error("Failed to blacklist old token", zap.Error(err))
	}

	return newToken, true, nil
}
