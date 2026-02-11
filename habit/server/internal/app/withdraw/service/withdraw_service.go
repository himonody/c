package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"habit/internal/app/withdraw/dto"
	"habit/internal/model"
	"habit/internal/repo"
	"habit/pkg/cache"
	"habit/pkg/database"
	"habit/pkg/utils"

	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type WithdrawService struct {
	db             *gorm.DB
	withdrawRepo   *repo.WithdrawRepository
	walletRepo     *repo.WalletRepository
	balanceLogRepo *repo.BalanceLogRepository
	logger         *zap.Logger
}

func NewWithdrawService(withdrawRepo *repo.WithdrawRepository, walletRepo *repo.WalletRepository, balanceLogRepo *repo.BalanceLogRepository, logger *zap.Logger) *WithdrawService {
	return &WithdrawService{
		db:             database.DB,
		withdrawRepo:   withdrawRepo,
		walletRepo:     walletRepo,
		balanceLogRepo: balanceLogRepo,
		logger:         logger,
	}
}

func (s *WithdrawService) GetWithdrawList(userID int64, req *dto.WithdrawListRequest) (*dto.WithdrawListResponse, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	records, total, err := s.withdrawRepo.ListByUserID(userID, req.Page, req.PageSize)
	if err != nil {
		s.logger.Error("Failed to get withdraw list", zap.Error(err))
		return nil, errors.New("failed to get withdraw list")
	}

	list := make([]*dto.WithdrawInfo, 0, len(records))
	for _, r := range records {
		reviewedAt := ""
		if r.ReviewedAt != nil {
			reviewedAt = r.ReviewedAt.Format(utils.MYDateLayout)
		}
		list = append(list, &dto.WithdrawInfo{
			BizID:        r.BizID,
			Amount:       r.Amount.String(),
			Fee:          r.Fee.String(),
			ActualAmount: r.ActualAmount.String(),
			Address:      r.Address,
			Status:       r.Status,
			RejectReason: r.RejectReason,
			TxHash:       r.TxHash,
			ReviewedAt:   reviewedAt,
			CreatedAt:    r.CreatedAt.Format(utils.MYDateLayout),
		})
	}

	return &dto.WithdrawListResponse{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

func (s *WithdrawService) ApplyWithdraw(userID int64, amount decimal.Decimal, address, applyIP, payPassword string) (string, error) {
	address = strings.TrimSpace(address)
	bizID := utils.GenerateWithdrawBizID(userID)
	feeRate := decimal.NewFromInt(3).Div(decimal.NewFromInt(100))

	ctx := context.Background()
	lockKey := fmt.Sprintf("lock:withdraw:apply:%d", userID)
	lock, locked, err := cache.AcquireLockWithRenew(ctx, database.RedisClient, lockKey, 5*time.Second)
	if err != nil {
		s.logger.Error("Failed to acquire withdraw lock", zap.Error(err))
		return "", errors.New("failed to apply withdraw")
	}
	if !locked {
		return "", errors.New("duplicate request")
	}
	defer func() {
		_ = lock.Unlock(ctx)
	}()

	err = s.db.Transaction(func(tx *gorm.DB) error {
		wallet, err := s.walletRepo.FindByUserIDForUpdate(tx, userID)
		if err != nil {
			return err
		}
		if wallet.PayPwd == "" {
			return errors.New("pay password not set")
		}
		if !utils.CheckPassword(payPassword, wallet.PayPwd) {
			return errors.New("invalid pay password")
		}

		fee := amount.Mul(feeRate).Round(8)
		actualAmount := amount.Sub(fee).Round(8)
		if actualAmount.LessThanOrEqual(decimal.Zero) {
			return errors.New("amount too small")
		}

		if wallet.Balance.LessThan(amount) {
			return errors.New("insufficient balance")
		}

		newBalance := wallet.Balance.Sub(amount)
		newFrozen := wallet.Frozen.Add(amount)
		if err := s.walletRepo.UpdateBalanceAndFrozen(tx, userID, newBalance, newFrozen); err != nil {
			return err
		}

		withdraw := &model.AppUserWithdraw{
			BizID:        bizID,
			UserID:       userID,
			Amount:       amount,
			Fee:          fee,
			ActualAmount: actualAmount,
			Address:      address,
			ApplyIP:      applyIP,
			Status:       1,
		}
		if err := s.withdrawRepo.Create(tx, withdraw); err != nil {
			return err
		}

		log := &model.AppUserBalanceLog{
			BizID:         bizID,
			UserID:        userID,
			Type:          1,
			Amount:        amount.Neg(),
			BeforeBalance: wallet.Balance,
			AfterBalance:  newBalance,
			BeforeFrozen:  wallet.Frozen,
			AfterFrozen:   newFrozen,
			Remark:        "提款申请",
			OperatorID:    0,
		}
		if err := s.balanceLogRepo.Create(tx, log); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		s.logger.Error("Failed to apply withdraw", zap.Error(err))
		if err.Error() == "insufficient balance" {
			return "", err
		}
		if err.Error() == "pay password not set" {
			return "", err
		}
		if err.Error() == "invalid pay password" {
			return "", err
		}
		if err.Error() == "amount too small" {
			return "", err
		}
		return "", errors.New("failed to apply withdraw")
	}

	return bizID, nil
}
