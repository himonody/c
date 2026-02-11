package service

import (
	"errors"
	"habit/internal/app/invite/dto"
	"habit/internal/repo"

	"go.uber.org/zap"
)

type InviteService struct {
	repo   *repo.InviteRepository
	logger *zap.Logger
}

func NewInviteService(repo *repo.InviteRepository, logger *zap.Logger) *InviteService {
	return &InviteService{repo: repo, logger: logger}
}

func (s *InviteService) GetInviteInfo(userID int64) (*dto.InviteInfoResponse, error) {
	u, err := s.repo.FindByUserID(userID)
	if err != nil {
		s.logger.Error("Failed to find user", zap.Int64("userID", userID), zap.Error(err))
		return nil, errors.New("user_not_found")
	}

	return &dto.InviteInfoResponse{
		FriendCode: u.FriendCode,
		InviteURL:  "https://example.com/invite?code=" + u.FriendCode, // TODO: 从配置读取域名
	}, nil
}

func (s *InviteService) ListInvitedFriends(userID int, req *dto.InviteFriendsRequest) (*dto.InviteFriendsResponse, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	users, total, err := s.repo.ListInvitedUsers(userID, req.Page, req.PageSize)
	if err != nil {
		s.logger.Error("Failed to list invited friends", zap.Int("userID", userID), zap.Error(err))
		return nil, errors.New("failed_to_list_friends")
	}

	list := make([]*dto.InviteFriendItem, 0, len(users))
	for _, u := range users {
		list = append(list, &dto.InviteFriendItem{
			Avatar:   u.Avatar,
			Nickname: u.Nickname,
			Username: u.Username,
		})
	}

	return &dto.InviteFriendsResponse{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}
