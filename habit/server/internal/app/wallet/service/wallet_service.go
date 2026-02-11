package service

import (
	"errors"
	"time"

	"habit/internal/app/wallet/dto"
	"habit/internal/model"
	"habit/internal/repo"

	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type WalletService struct {
	walletRepo *repo.WalletRepository
	logger     *zap.Logger
}

func NewWalletService(walletRepo *repo.WalletRepository, logger *zap.Logger) *WalletService {
	return &WalletService{
		walletRepo: walletRepo,
		logger:     logger,
	}
}

// GetWalletInfo 获取钱包信息（带缓存）
func (s *WalletService) GetWalletInfo(userID int64) (*dto.WalletInfo, error) {
	wallet, err := s.walletRepo.FindByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 钱包不存在，创建新钱包
			wallet = &model.AppUserWallet{
				UserID:    userID,
				PayStatus: "1",
				Balance:   decimal.NewFromFloat(0),
				Frozen:    decimal.NewFromFloat(0),
				TotalR:    decimal.NewFromFloat(0),
				TotalW:    decimal.NewFromFloat(0),
				TotalRe:   decimal.NewFromFloat(0),
				TotalI:    decimal.NewFromFloat(0),
			}
			if err := s.walletRepo.Create(wallet); err != nil {
				s.logger.Error("Failed to create wallet", zap.Error(err))
				return nil, errors.New("failed to create wallet")
			}
		} else {
			s.logger.Error("Failed to find wallet", zap.Error(err))
			return nil, errors.New("failed to get wallet")
		}
	}

	return &dto.WalletInfo{
		UserId:    wallet.UserID,
		PayStatus: wallet.PayStatus,
		Balance:   wallet.Balance.String(),   // 转为字符串
		Frozen:    wallet.Frozen.String(),    // 转为字符串
		TotalR:    wallet.TotalR.String(),    // 转为字符串
		TotalW:    wallet.TotalW.String(),    // 转为字符串
		TotalRe:   wallet.TotalRe.String(),   // 转为字符串
		TotalI:    wallet.TotalI.String(),    // 转为字符串
		CreatedAt: wallet.CreatedAt.Format(time.DateTime),
		UpdatedAt: wallet.UpdatedAt.Format(time.DateTime),
	}, nil
}

func (s *WalletService) SetWalletAddress(userID int64, address string) error {
	// Ensure wallet exists
	_, err := s.walletRepo.FindByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			wallet := &model.AppUserWallet{
				UserID:    userID,
				PayStatus: "1",
				Balance:   decimal.NewFromFloat(0),
				Frozen:    decimal.NewFromFloat(0),
				TotalR:    decimal.NewFromFloat(0),
				TotalW:    decimal.NewFromFloat(0),
				TotalRe:   decimal.NewFromFloat(0),
				TotalI:    decimal.NewFromFloat(0),
				Address:   address,
			}
			if err := s.walletRepo.Create(wallet); err != nil {
				s.logger.Error("Failed to create wallet", zap.Error(err))
				return errors.New("failed to create wallet")
			}
			return nil
		}
		s.logger.Error("Failed to find wallet", zap.Error(err))
		return errors.New("failed to get wallet")
	}

	if err := s.walletRepo.UpdateAddress(userID, address); err != nil {
		s.logger.Error("Failed to update wallet address", zap.Error(err))
		return errors.New("failed to update address")
	}

	return nil
}
