package service

import (
	"errors"
	"time"

	"habit/internal/app/challenge/dto"
	"habit/internal/repo"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ChallengeStatService struct {
	repo   *repo.ChallengeStatRepository
	logger *zap.Logger
}

func NewChallengeStatService(repo *repo.ChallengeStatRepository, logger *zap.Logger) *ChallengeStatService {
	return &ChallengeStatService{repo: repo, logger: logger}
}

func (s *ChallengeStatService) GetTotalStat() (*dto.ChallengeTotalStatResponse, error) {
	stat, err := s.repo.GetTotalStat()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &dto.ChallengeTotalStatResponse{}, nil
		}
		s.logger.Error("Failed to get challenge total stat", zap.Error(err))
		return nil, errors.New("failed to get total stat")
	}

	resp := &dto.ChallengeTotalStatResponse{
		TotalUserCnt:       stat.TotalUserCnt,
		TotalJoinCnt:       stat.TotalJoinCnt,
		TotalSuccessCnt:    stat.TotalSuccessCnt,
		TotalFailCnt:       stat.TotalFailCnt,
		TotalJoinAmount:    stat.TotalJoinAmount.String(),
		TotalSuccessAmount: stat.TotalSuccessAmount.String(),
		TotalFailAmount:    stat.TotalFailAmount.String(),
		TotalPlatformBonus: stat.TotalPlatformBonus.String(),
		TotalPoolAmount:    stat.TotalPoolAmount.String(),
		UpdatedAt:          stat.UpdatedAt.Format(time.DateTime),
	}
	return resp, nil
}
