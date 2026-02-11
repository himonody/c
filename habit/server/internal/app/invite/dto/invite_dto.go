package dto

type InviteInfoResponse struct {
	FriendCode string `json:"friend_code"`
	InviteURL  string `json:"invite_url"`
}

type InviteFriendsRequest struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}

type InviteFriendItem struct {
	Avatar   string `json:"avatar"`
	Nickname string `json:"nickname"`
	Username string `json:"username"`
}

type InviteFriendsResponse struct {
	List     []*InviteFriendItem `json:"list"`
	Total    int64               `json:"total"`
	Page     int                 `json:"page"`
	PageSize int                 `json:"pageSize"`
}
