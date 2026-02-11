package service

import (
	"errors"
	"time"

	"habit/internal/app/leaderboard/dto"
	"habit/internal/repo"

	"go.uber.org/zap"
)

type LeaderboardService struct {
	repo   *repo.LeaderboardRepository
	logger *zap.Logger
}

func NewLeaderboardService(repo *repo.LeaderboardRepository, logger *zap.Logger) *LeaderboardService {
	return &LeaderboardService{repo: repo, logger: logger}
}

func (s *LeaderboardService) List(req *dto.LeaderboardListRequest) (*dto.LeaderboardListResponse, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}
	if req.RankType <= 0 {
		return nil, errors.New("bad_request")
	}

	rankDate := time.Now()
	if req.RankDate != "" {
		t, err := time.Parse("2006-01-02", req.RankDate)
		if err != nil {
			return nil, errors.New("bad_request")
		}
		rankDate = t
	}

	records, total, err := s.repo.ListDaily(req.RankType, rankDate, req.Page, req.PageSize)
	if err != nil {
		s.logger.Error("Failed to list leaderboard", zap.Error(err))
		return nil, errors.New("failed to list leaderboard")
	}

	list := make([]*dto.LeaderboardItem, 0, len(records))
	for _, r := range records {
		list = append(list, &dto.LeaderboardItem{
			UserID: r.UserID,
			Value:  r.Value.String(),
			RankNo: r.RankNo,
		})
	}

	return &dto.LeaderboardListResponse{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

func (s *LeaderboardService) Total(req *dto.LeaderboardTotalRequest) (*dto.LeaderboardTotalResponse, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}
	if req.RankType <= 0 {
		return nil, errors.New("bad_request")
	}

	records, total, err := s.repo.ListTotalSum(req.RankType, req.Page, req.PageSize)
	if err != nil {
		s.logger.Error("Failed to list total leaderboard", zap.Error(err))
		return nil, errors.New("failed to list leaderboard")
	}

	offset := (req.Page - 1) * req.PageSize
	list := make([]*dto.LeaderboardItem, 0, len(records))
	for i, r := range records {
		list = append(list, &dto.LeaderboardItem{
			UserID: r.UserID,
			Value:  r.Value.String(),
			RankNo: offset + i + 1,
		})
	}

	return &dto.LeaderboardTotalResponse{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

