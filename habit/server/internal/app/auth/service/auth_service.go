package service

import (
	"context"
	"errors"
	"time"

	"habit/internal/app/auth/dto"
	"habit/internal/model"
	"habit/internal/repo"
	"habit/pkg/cache"
	"habit/pkg/database"
	"habit/pkg/utils"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AuthService struct {
	userRepo   *repo.UserRepository
	walletRepo *repo.WalletRepository
	logger     *zap.Logger
	redis      *redis.Client
}

func NewAuthService(userRepo *repo.UserRepository, walletRepo *repo.WalletRepository, logger *zap.Logger) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		walletRepo: walletRepo,
		logger:     logger,
		redis:      database.RedisClient,
	}
}

// storeTokenInRedis stores token in Redis for single sign-on
// Only one token per user is allowed (single sign-on)
func (s *AuthService) storeTokenInRedis(userID int64, token string) error {
	ctx := context.Background()
	tokenKey := cache.UserTokenKey(userID)

	// Delete old token if exists (enforce single sign-on)
	oldToken, err := s.redis.Get(ctx, tokenKey).Result()
	if err == nil && oldToken != "" && oldToken != token {
		// Blacklist old token
		oldTokenKey := cache.TokenBlacklistKey(oldToken)
		s.redis.Set(ctx, oldTokenKey, "1", cache.TokenBlacklistExpiration)
	}

	// Store new token with expiration
	return s.redis.Set(ctx, tokenKey, token, cache.UserTokenExpiration).Err()
}

// isTokenBlacklisted checks if token is in blacklist
func (s *AuthService) isTokenBlacklisted(token string) bool {
	ctx := context.Background()
	blacklistKey := cache.TokenBlacklistKey(token)
	_, err := s.redis.Get(ctx, blacklistKey).Result()
	return err == nil
}

// Register handles user registration
func (s *AuthService) Register(req *dto.RegisterRequest, registerIP string) (*dto.LoginResponse, error) {
	// Check if username already exists
	exists, err := s.userRepo.UsernameExists(req.Username)
	if err != nil {
		s.logger.Error("Failed to check username existence", zap.Error(err))
		return nil, errors.New("internal server error")
	}
	if exists {
		return nil, errors.New("username already exists")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		s.logger.Error("Failed to hash password", zap.Error(err))
		return nil, errors.New("internal server error")
	}

	// Handle friend code (invite code)
	var refID int
	if req.FriendCode != "" {
		inviter, err := s.userRepo.FindByRefCode(req.FriendCode)
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				s.logger.Error("Failed to find inviter", zap.Error(err))
			}
			// If friend code is invalid, just ignore it (set refID to 0)
			refID = 0
		} else {
			refID = int(inviter.ID)
		}
	}

	// Generate random nickname (6 characters)
	randomNickname := utils.GenerateRandomNickname()

	// Create user
	user := &model.AppUser{
		Username:     req.Username,
		Nickname:     randomNickname, // Auto-generated random nickname
		Pwd:          hashedPassword,
		RefID:        refID,
		FriendCode:   req.FriendCode,
		Status:       "1",
		OnlineStatus: "2", // Online after registration
		RegisterAt:   time.Now(),
		RegisterIP:   registerIP,
	}

	if err := s.userRepo.Create(user); err != nil {
		s.logger.Error("Failed to create user", zap.Error(err))
		return nil, errors.New("failed to create user")
	}

	// Generate reference code after user is created (we need the user ID)
	refCode := utils.GenerateRefCode(user.ID)
	user.RefCode = refCode

	// Update user with ref code
	if err := database.DB.Model(user).Update("ref_code", refCode).Error; err != nil {
		s.logger.Error("Failed to update ref code", zap.Error(err))
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		s.logger.Error("Failed to generate token", zap.Error(err))
		return nil, errors.New("failed to generate token")
	}

	// Store token in Redis (single sign-on)
	if err := s.storeTokenInRedis(user.ID, token); err != nil {
		s.logger.Error("Failed to store token in Redis", zap.Error(err))
	}

	return &dto.LoginResponse{
		Token: token,
		UserInfo: dto.UserInfo{
			ID:           user.ID,
			Username:     user.Username,
			Nickname:     user.Nickname,
			Avatar:       user.Avatar,
			RefCode:      user.RefCode,
			Status:       user.Status,
			OnlineStatus: user.OnlineStatus,
		},
	}, nil
}

// Login handles user login
func (s *AuthService) Login(req *dto.LoginRequest, loginIP string) (*dto.LoginResponse, error) {
	// Find user by username
	user, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid username or password")
		}
		s.logger.Error("Failed to find user", zap.Error(err))
		return nil, errors.New("internal server error")
	}

	// Check if user is active
	if user.Status != "1" {
		return nil, errors.New("account is disabled")
	}

	// Verify password
	if !utils.CheckPassword(req.Password, user.Pwd) {
		return nil, errors.New("invalid username or password")
	}

	// Update last login info
	if err := s.userRepo.UpdateLastLogin(user.ID, loginIP); err != nil {
		s.logger.Error("Failed to update last login", zap.Error(err))
	}

	// Update online status
	if err := s.userRepo.UpdateOnlineStatus(user.ID, "2"); err != nil {
		s.logger.Error("Failed to update online status", zap.Error(err))
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		s.logger.Error("Failed to generate token", zap.Error(err))
		return nil, errors.New("failed to generate token")
	}

	// Store token in Redis (single sign-on - will invalidate old token)
	if err := s.storeTokenInRedis(user.ID, token); err != nil {
		s.logger.Error("Failed to store token in Redis", zap.Error(err))
	}

	return &dto.LoginResponse{
		Token: token,
		UserInfo: dto.UserInfo{
			ID:           user.ID,
			Username:     user.Username,
			Nickname:     user.Nickname,
			Avatar:       user.Avatar,
			RefCode:      user.RefCode,
			Status:       user.Status,
			OnlineStatus: user.OnlineStatus,
		},
	}, nil
}

// Logout handles user logout
func (s *AuthService) Logout(userID int64, token string) error {
	// Update online status to offline
	if err := s.userRepo.UpdateOnlineStatus(userID, "1"); err != nil {
		s.logger.Error("Failed to update online status", zap.Error(err))
		return err
	}

	ctx := context.Background()

	// Add token to blacklist
	blacklistKey := cache.TokenBlacklistKey(token)
	if err := s.redis.Set(ctx, blacklistKey, "1", cache.TokenBlacklistExpiration).Err(); err != nil {
		s.logger.Error("Failed to blacklist token", zap.Error(err))
	}

	// Delete token from Redis
	tokenKey := cache.UserTokenKey(userID)
	if err := s.redis.Del(ctx, tokenKey).Err(); err != nil {
		s.logger.Error("Failed to delete token from Redis", zap.Error(err))
		return err
	}

	return nil
}

// ValidateToken validates a token and returns user ID and claims
func (s *AuthService) ValidateToken(token string) (int64, *utils.Claims, error) {
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
	tokenKey := cache.UserTokenKey(claims.UserID)
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

// RefreshToken generates a new token if needed
func (s *AuthService) RefreshToken(userID int64, username string, oldToken string) (string, bool, error) {
	// Generate new token
	newToken, err := utils.GenerateToken(userID, username)
	if err != nil {
		s.logger.Error("Failed to generate new token", zap.Error(err))
		return "", false, err
	}

	// Store new token in Redis
	if err := s.storeTokenInRedis(userID, newToken); err != nil {
		s.logger.Error("Failed to store new token in Redis", zap.Error(err))
		return "", false, err
	}

	// Blacklist old token
	ctx := context.Background()
	blacklistKey := cache.TokenBlacklistKey(oldToken)
	if err := s.redis.Set(ctx, blacklistKey, "1", cache.TokenBlacklistExpiration).Err(); err != nil {
		s.logger.Error("Failed to blacklist old token", zap.Error(err))
	}

	return newToken, true, nil
}

// ChangePassword 修改密码
func (s *AuthService) ChangePassword(userID int64, req *dto.ChangePasswordRequest) error {
	// 获取用户信息
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		s.logger.Error("Failed to find user", zap.Error(err))
		return errors.New("user not found")
	}

	// 验证旧密码
	if !utils.CheckPassword(req.OldPassword, user.Pwd) {
		return errors.New("old password is incorrect")
	}

	// 哈希新密码
	newPasswordHash, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		s.logger.Error("Failed to hash password", zap.Error(err))
		return errors.New("internal server error")
	}

	// 更新密码
	if err := s.userRepo.UpdatePassword(userID, newPasswordHash); err != nil {
		s.logger.Error("Failed to update password", zap.Error(err))
		return errors.New("failed to update password")
	}

	return nil
}

// SetPayPassword 设置/修改支付密码
func (s *AuthService) SetPayPassword(userID int64, req *dto.SetPayPasswordRequest) error {
	// 获取钱包信息
	wallet, err := s.walletRepo.FindByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 钱包不存在，创建新钱包
			wallet = &model.AppUserWallet{
				UserID:    userID,
				PayStatus: "1",
			}
			if err := s.walletRepo.Create(wallet); err != nil {
				s.logger.Error("Failed to create wallet", zap.Error(err))
				return errors.New("failed to create wallet")
			}
		} else {
			s.logger.Error("Failed to find wallet", zap.Error(err))
			return errors.New("internal server error")
		}
	}

	// 如果已设置支付密码，需要验证旧密码
	if wallet.PayPwd != "" && req.OldPayPassword == "" {
		return errors.New("old pay password is required")
	}

	// 验证旧支付密码
	if wallet.PayPwd != "" {
		if !utils.CheckPassword(req.OldPayPassword, wallet.PayPwd) {
			return errors.New("old pay password is incorrect")
		}
	}

	// 哈希新支付密码
	payPasswordHash, err := utils.HashPassword(req.PayPassword)
	if err != nil {
		s.logger.Error("Failed to hash pay password", zap.Error(err))
		return errors.New("internal server error")
	}

	// 更新支付密码
	if err := s.walletRepo.UpdatePayPassword(userID, payPasswordHash); err != nil {
		s.logger.Error("Failed to update pay password", zap.Error(err))
		return errors.New("failed to update pay password")
	}

	return nil
}

// GetUserInfo 获取用户信息（带缓存）
func (s *AuthService) GetUserInfo(userID int64) (*dto.UserInfo, error) {
	// 从数据库查询用户信息（Repository 层已实现缓存）
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		s.logger.Error("Failed to find user", zap.Error(err))
		return nil, errors.New("user not found")
	}

	return &dto.UserInfo{
		ID:           user.ID,
		Username:     user.Username,
		Nickname:     user.Nickname,
		Avatar:       user.Avatar,
		RefCode:      user.RefCode,
		Status:       user.Status,
		OnlineStatus: user.OnlineStatus,
	}, nil
}

// UpdateProfile 更新用户资料
func (s *AuthService) UpdateProfile(userID int64, req *dto.UpdateProfileRequest) error {
	updates := make(map[string]interface{})

	if req.Nickname != "" {
		updates["nickname"] = req.Nickname
	}

	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}

	if len(updates) == 0 {
		return errors.New("no fields to update")
	}

	if err := s.userRepo.UpdateProfile(userID, updates); err != nil {
		s.logger.Error("Failed to update profile", zap.Error(err))
		return errors.New("failed to update profile")
	}

	return nil
}
